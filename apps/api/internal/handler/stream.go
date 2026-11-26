package handler

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/scraper"
	"jevbeez/api/internal/security"
)

func sendSSE(w *bufio.Writer, event string, data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b); err != nil {
		return err
	}
	return w.Flush()
}

func (h *AnalyzeHandler) StreamHandle(c *fiber.Ctx) error {
	clientIP := c.IP()

	// 1. Concurrency limit per IP (Anti-Slowloris: max 2 active streams per IP)
	val, _ := h.activeStreams.LoadOrStore(clientIP, new(int32))
	activeCount := val.(*int32)
	current := atomic.AddInt32(activeCount, 1)
	if current > 2 {
		atomic.AddInt32(activeCount, -1)
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Terlalu banyak permintaan streaming bersamaan dari IP ini",
		})
	}

	// 2. Extract parameters (supports query or body)
	rawURL := strings.TrimSpace(c.Query("url"))
	rawText := strings.TrimSpace(c.Query("text"))
	deepDive := c.Query("deep_dive") == "true"

	if rawURL == "" && rawText == "" && len(c.Body()) > 0 {
		var req model.AnalyzeRequest
		if err := c.BodyParser(&req); err == nil {
			rawURL = strings.TrimSpace(req.URL)
			rawText = strings.TrimSpace(req.Text)
			deepDive = req.DeepDive
		}
	}

	if rawURL == "" && rawText == "" {
		atomic.AddInt32(activeCount, -1)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Silakan masukkan tautan (URL) atau teks pesan yang ingin diverifikasi",
		})
	}

	// 3. Set SSE headers (with X-Accel-Buffering for Traefik/Nginx)
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer atomic.AddInt32(activeCount, -1)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if rawURL != "" {
			h.streamURL(ctx, w, clientIP, rawURL, deepDive)
			return
		}
		h.streamText(ctx, w, clientIP, rawText)
	}))

	return nil
}

func (h *AnalyzeHandler) streamURL(ctx context.Context, w *bufio.Writer, clientIP string, rawURL string, deepDive bool) {
	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "start",
		"message": "Memeriksa tautan...",
	})

	if isExploit, reason, detail := security.DetectExploit(rawURL); isExploit {
		_ = h.cache.BlockIP(ctx, clientIP, reason, rawURL, config.DefaultBlockDuration)
		_ = sendSSE(w, "error", fiber.Map{
			"error":   fmt.Sprintf("Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena terdeteksi aktivitas eksploitasi (%s)", detail),
			"blocked": true,
			"reason":  reason,
		})
		return
	}

	normalizedURL, err := security.ValidateAndNormalizeURL(rawURL)
	if err != nil {
		errMsg := fmt.Sprintf("URL tidak valid: %s", err.Error())
		if errors.Is(err, security.ErrSelfURL) {
			errMsg = "Tautan yang Anda masukkan adalah situs JevBeez ini sendiri. Silakan coba periksa tautan berita dari web lain."
		}
		_ = sendSSE(w, "error", fiber.Map{"error": errMsg})
		return
	}

	cacheKey := normalizedURL
	if deepDive {
		cacheKey = normalizedURL + ":deep"
	}

	hSum := sha256.Sum256([]byte(cacheKey))
	shareHash := hex.EncodeToString(hSum[:])

	if cached, cErr := h.cache.GetByHash(ctx, shareHash); cErr == nil && cached != nil {
		if cached.ShareHash == "" {
			cached.ShareHash = shareHash
		}
		_ = sendSSE(w, "progress", fiber.Map{
			"phase":   "cached",
			"message": "Memuat hasil terverifikasi dari cache...",
		})
		_ = sendSSE(w, "verdict", cached)
		return
	}

	if allowed, remaining, cdErr := h.cache.CheckAndSetCooldown(ctx, clientIP, config.DefaultCooldownDuration); cdErr == nil && !allowed {
		_ = sendSSE(w, "error", fiber.Map{
			"error":            fmt.Sprintf("Mohon tunggu %.0f detik sebelum melakukan pengecekan berikutnya", remaining.Seconds()),
			"cooldown_seconds": int(remaining.Seconds()),
		})
		return
	}

	targetURL := normalizedURL
	var redirectHops []string
	if unwrap := security.UnwrapRedirectChain(normalizedURL); unwrap.Depth > 0 {
		targetURL = unwrap.FinalURL
		redirectHops = append(redirectHops, unwrap.Chain...)
	}

	if security.IsShortener(targetURL) {
		unrolled, hops, sErr := security.ResolveShortener(ctx, targetURL, config.MaxShortenerHops)
		if sErr == nil && unrolled != "" {
			targetURL = unrolled
			redirectHops = append(redirectHops, hops...)
		}
	}

	threatCheck := security.CheckURLHeuristics(targetURL)
	if !threatCheck.IsThreat && targetURL != normalizedURL {
		threatCheck = security.CheckURLHeuristics(normalizedURL)
	}

	if !threatCheck.IsThreat && security.IsEcommerceURL(targetURL) {
		destDomain := targetURL
		if parsed, pErr := url.Parse(targetURL); pErr == nil {
			destDomain = parsed.Hostname()
		}
		result := h.evaluator.EvaluateAds(rawURL, targetURL, destDomain, redirectHops)
		result.ShareHash = shareHash
		_ = h.cache.SetByHash(ctx, shareHash, result, config.DefaultCacheTTL)
		_ = h.cache.SetShare(ctx, shareHash, result)
		_ = sendSSE(w, "verdict", result)
		return
	}

	if threatCheck.IsThreat {
		_ = sendSSE(w, "progress", fiber.Map{
			"phase":   "threat_detected",
			"message": "Terdeteksi potensi ancaman keamanan...",
		})

		if !deepDive {
			threatInput := &model.ThreatEvaluationInput{
				RawURL:       rawURL,
				ResolvedURL:  targetURL,
				RedirectHops: redirectHops,
				ThreatType:   threatCheck.ThreatType,
				Confidence:   threatCheck.Confidence,
				Reasons:      threatCheck.Reasons,
				TargetBrand:  threatCheck.TargetBrand,
				IsDeepDive:   false,
			}
			result := h.evaluator.EvaluateThreat(threatInput)
			result.ShareHash = shareHash
			_ = h.cache.SetByHash(ctx, shareHash, result, config.DefaultCacheTTL)
			_ = h.cache.SetShare(ctx, shareHash, result)
			_ = sendSSE(w, "verdict", result)
			return
		}

		formThreats, pageTitle, _ := scraper.InspectFormThreats(ctx, targetURL)
		confidence := threatCheck.Confidence
		if len(formThreats) > 0 && confidence < 95.0 {
			confidence = 96.0
		}
		threatInput := &model.ThreatEvaluationInput{
			RawURL:       rawURL,
			ResolvedURL:  targetURL,
			RedirectHops: redirectHops,
			ThreatType:   threatCheck.ThreatType,
			Confidence:   confidence,
			Reasons:      threatCheck.Reasons,
			TargetBrand:  threatCheck.TargetBrand,
			FormThreats:  formThreats,
			PageTitle:    pageTitle,
			IsDeepDive:   true,
		}
		result := h.evaluator.EvaluateThreat(threatInput)
		result.ShareHash = shareHash
		_ = h.cache.SetByHash(ctx, shareHash, result, config.DefaultCacheTTL)
		_ = h.cache.SetShare(ctx, shareHash, result)
		_ = sendSSE(w, "verdict", result)
		return
	}

	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "extracting",
		"message": "Mengekstrak konten artikel...",
	})

	article, err := h.scraper.Extract(ctx, targetURL)
	if err != nil {
		if errors.Is(err, security.ErrPrivateIPBlocked) || errors.Is(err, security.ErrInvalidIPAddress) {
			_ = h.cache.BlockIP(ctx, clientIP, security.ReasonSSRFAttempt, normalizedURL, config.DefaultBlockDuration)
			_ = sendSSE(w, "error", fiber.Map{
				"error":   "Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena mencoba mengakses alamat IP privat/internal (SSRF)",
				"blocked": true,
				"reason":  security.ReasonSSRFAttempt,
			})
			return
		}

		errMsg := err.Error()
		if errMsg == "" || errors.Is(err, scraper.ErrEmptyArticle) {
			errMsg = "Gagal membaca konten dari URL tersebut, situs mungkin tidak dapat diakses atau dibatasi"
		}
		_ = sendSSE(w, "error", fiber.Map{"error": errMsg})
		return
	}

	_ = sendSSE(w, "metadata", fiber.Map{
		"title":         article.Title,
		"source_domain": article.SourceDomain,
		"word_count":    article.WordCount,
	})

	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "searching",
		"message": "Mencari bukti dan fakta pembanding...",
	})

	evidences, err := h.searcher.SearchGrounding(ctx, article)
	if err != nil {
		_ = sendSSE(w, "error", fiber.Map{"error": "Layanan pencarian bukti sedang mengalami gangguan"})
		return
	}

	_ = sendSSE(w, "evidence", fiber.Map{"evidences": evidences})

	if len(evidences) == 0 {
		result := h.evaluator.Evaluate(article, evidences, &model.JevJudgments{})
		result.ShareHash = shareHash
		_ = h.cache.SetByHash(ctx, shareHash, result, config.UnverifiedCacheTTL)
		_ = h.cache.SetShare(ctx, shareHash, result)
		_ = sendSSE(w, "verdict", result)
		return
	}

	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "evaluating",
		"message": "Menganalisis kredibilitas klaim dengan AI Decision Model...",
	})

	judgments, err := h.jevClient.Evaluate(ctx, article, evidences)
	if err != nil {
		_ = sendSSE(w, "error", fiber.Map{"error": "Layanan evaluasi fakta sedang mengalami gangguan"})
		return
	}

	result := h.evaluator.Evaluate(article, evidences, judgments)
	result.ShareHash = shareHash

	cacheTTL := config.DefaultCacheTTL
	if result.Verdict == config.VerdictUnverified {
		cacheTTL = config.UnverifiedCacheTTL
	}
	_ = h.cache.SetByHash(ctx, shareHash, result, cacheTTL)
	_ = h.cache.SetShare(ctx, shareHash, result)

	_ = sendSSE(w, "verdict", result)
}

func (h *AnalyzeHandler) streamText(ctx context.Context, w *bufio.Writer, clientIP string, rawText string) {
	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "start",
		"message": "Membersihkan teks pesan...",
	})

	cleanedText, searchQuery := security.CleanBroadcastText(rawText)
	if cleanedText == "" || searchQuery == "" {
		_ = sendSSE(w, "error", fiber.Map{
			"error": "Teks pesan terlalu pendek atau tidak memuat klaim yang cukup setelah dibersihkan (minimal 25 karakter dan 4 kata)",
		})
		return
	}

	hSum := sha256.Sum256([]byte(cleanedText))
	shareHash := hex.EncodeToString(hSum[:])

	if cached, err := h.cache.GetByHash(ctx, shareHash); err == nil && cached != nil {
		if cached.ShareHash == "" {
			cached.ShareHash = shareHash
		}
		_ = sendSSE(w, "progress", fiber.Map{
			"phase":   "cached",
			"message": "Memuat hasil terverifikasi dari cache...",
		})
		_ = sendSSE(w, "verdict", cached)
		return
	}

	if allowed, remaining, err := h.cache.CheckAndSetCooldown(ctx, clientIP, config.DefaultCooldownDuration); err == nil && !allowed {
		_ = sendSSE(w, "error", fiber.Map{
			"error":            fmt.Sprintf("Mohon tunggu %.0f detik sebelum melakukan pengecekan berikutnya", remaining.Seconds()),
			"cooldown_seconds": int(remaining.Seconds()),
		})
		return
	}

	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "searching",
		"message": "Mencari bukti terkait klaim pesan...",
	})

	evidences, err := h.searcher.SearchQuery(ctx, searchQuery)
	if err != nil {
		_ = sendSSE(w, "error", fiber.Map{"error": "Layanan pencarian bukti sedang mengalami gangguan"})
		return
	}

	_ = sendSSE(w, "evidence", fiber.Map{"evidences": evidences})

	pseudoArticle := &model.Article{
		Title:        cleanedText,
		Content:      cleanedText,
		SourceDomain: "Pesan Berantai / WhatsApp",
		URL:          "",
		WordCount:    len(strings.Fields(cleanedText)),
	}

	if len(evidences) == 0 {
		result := h.evaluator.Evaluate(pseudoArticle, evidences, &model.JevJudgments{})
		result.IsTextMode = true
		result.ShareHash = shareHash
		_ = h.cache.SetByHash(ctx, shareHash, result, config.UnverifiedCacheTTL)
		_ = h.cache.SetShare(ctx, shareHash, result)
		_ = sendSSE(w, "verdict", result)
		return
	}

	_ = sendSSE(w, "progress", fiber.Map{
		"phase":   "evaluating",
		"message": "Menganalisis kredibilitas pesan dengan AI...",
	})

	judgments, err := h.jevClient.Evaluate(ctx, pseudoArticle, evidences)
	if err != nil {
		_ = sendSSE(w, "error", fiber.Map{"error": "Layanan evaluasi fakta sedang mengalami gangguan"})
		return
	}

	result := h.evaluator.Evaluate(pseudoArticle, evidences, judgments)
	result.IsTextMode = true
	result.ShareHash = shareHash

	cacheTTL := config.DefaultCacheTTL
	if result.Verdict == config.VerdictUnverified {
		cacheTTL = config.UnverifiedCacheTTL
	}
	_ = h.cache.SetByHash(ctx, shareHash, result, cacheTTL)
	_ = h.cache.SetShare(ctx, shareHash, result)

	_ = sendSSE(w, "verdict", result)
}
