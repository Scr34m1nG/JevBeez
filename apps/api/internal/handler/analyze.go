package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"jevbeez/api/internal/cache"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/decision"
	"jevbeez/api/internal/jev"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/scraper"
	"jevbeez/api/internal/search"
	"jevbeez/api/internal/security"
	"sync"
)

type AnalyzeHandler struct {
	cache         *cache.CacheService
	scraper       *scraper.ScraperService
	searcher      *search.SearchManager
	jevClient     *jev.Client
	evaluator     *decision.Evaluator
	activeStreams sync.Map
}

func NewAnalyzeHandler(
	c *cache.CacheService,
	sc *scraper.ScraperService,
	sm *search.SearchManager,
	jc *jev.Client,
	ev *decision.Evaluator,
) *AnalyzeHandler {
	return &AnalyzeHandler{
		cache:     c,
		scraper:   sc,
		searcher:  sm,
		jevClient: jc,
		evaluator: ev,
	}
}

func (h *AnalyzeHandler) Handle(c *fiber.Ctx) error {
	clientIP := c.IP()

	var req model.AnalyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid, pastikan JSON memiliki field 'url' atau 'text'",
		})
	}

	trimmedURL := strings.TrimSpace(req.URL)
	trimmedText := strings.TrimSpace(req.Text)

	if trimmedURL == "" && trimmedText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Silakan masukkan tautan (URL) atau teks pesan yang ingin diverifikasi",
		})
	}

	if trimmedURL != "" {
		return h.handleURL(c, clientIP, req, trimmedURL)
	}

	return h.handleText(c, clientIP, trimmedText)
}

func (h *AnalyzeHandler) handleURL(c *fiber.Ctx, clientIP string, req model.AnalyzeRequest, rawURL string) error {
	if isExploit, reason, detail := security.DetectExploit(rawURL); isExploit {
		_ = h.cache.BlockIP(c.Context(), clientIP, reason, rawURL, config.DefaultBlockDuration)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   fmt.Sprintf("Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena terdeteksi aktivitas eksploitasi (%s)", detail),
			"blocked": true,
			"reason":  reason,
		})
	}

	normalizedURL, err := security.ValidateAndNormalizeURL(rawURL)
	if err != nil {
		if errors.Is(err, security.ErrSelfURL) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tautan yang Anda masukkan adalah situs JevBeez ini sendiri. Silakan coba periksa tautan berita dari web lain.",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("URL tidak valid: %s", err.Error()),
		})
	}

	cacheKey := normalizedURL
	if req.DeepDive {
		cacheKey = normalizedURL + ":deep"
	}

	hSum := sha256.Sum256([]byte(cacheKey))
	shareHash := hex.EncodeToString(hSum[:])

	if cached, err := h.cache.GetByHash(c.Context(), shareHash); err == nil && cached != nil {
		if cached.ShareHash == "" {
			cached.ShareHash = shareHash
		}
		return c.JSON(cached)
	}

	if allowed, remaining, err := h.cache.CheckAndSetCooldown(c.Context(), clientIP, config.DefaultCooldownDuration); err == nil && !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", remaining.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":            fmt.Sprintf("Mohon tunggu %.0f detik sebelum melakukan pengecekan berikutnya", remaining.Seconds()),
			"cooldown_seconds": int(remaining.Seconds()),
		})
	}

	targetURL := normalizedURL
	var redirectHops []string
	if unwrap := security.UnwrapRedirectChain(normalizedURL); unwrap.Depth > 0 {
		targetURL = unwrap.FinalURL
		redirectHops = append(redirectHops, unwrap.Chain...)
	}

	if security.IsShortener(targetURL) {
		unrolled, hops, sErr := security.ResolveShortener(c.Context(), targetURL, config.MaxShortenerHops)
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
		_ = h.cache.SetByHash(c.Context(), shareHash, result, config.DefaultCacheTTL)
		_ = h.cache.SetShare(c.Context(), shareHash, result)
		return c.JSON(result)
	}

	if threatCheck.IsThreat {
		if !req.DeepDive {
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
			_ = h.cache.SetByHash(c.Context(), shareHash, result, config.DefaultCacheTTL)
			_ = h.cache.SetShare(c.Context(), shareHash, result)
			return c.JSON(result)
		}

		formThreats, pageTitle, _ := scraper.InspectFormThreats(c.Context(), targetURL)
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
		_ = h.cache.SetByHash(c.Context(), shareHash, result, config.DefaultCacheTTL)
		_ = h.cache.SetShare(c.Context(), shareHash, result)
		return c.JSON(result)
	}

	article, err := h.scraper.Extract(c.Context(), targetURL)
	if err != nil {
		if errors.Is(err, security.ErrPrivateIPBlocked) || errors.Is(err, security.ErrInvalidIPAddress) {
			_ = h.cache.BlockIP(c.Context(), clientIP, security.ReasonSSRFAttempt, normalizedURL, config.DefaultBlockDuration)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Akses ditolak: Alamat IP Anda diblokir selama 24 jam karena mencoba mengakses alamat IP privat/internal (SSRF)",
				"blocked": true,
				"reason":  security.ReasonSSRFAttempt,
			})
		}

		if len(redirectHops) > 0 && errors.Is(err, scraper.ErrEmptyArticle) {
			result := h.evaluator.EvaluateThreat(&model.ThreatEvaluationInput{
				RawURL:       rawURL,
				ResolvedURL:  targetURL,
				RedirectHops: redirectHops,
				ThreatType:   config.VerdictScam,
				Confidence:   80.0,
				Reasons:      []string{"Tautan dibungkus redirector dan halaman tujuan tidak memiliki konten yang dapat dibaca (cloaking)."},
			})
			result.ShareHash = shareHash
			_ = h.cache.SetByHash(c.Context(), shareHash, result, config.DefaultCacheTTL)
			_ = h.cache.SetShare(c.Context(), shareHash, result)
			return c.JSON(result)
		}

		errMsg := err.Error()
		if errMsg == "" || errors.Is(err, scraper.ErrEmptyArticle) {
			errMsg = "Gagal membaca konten dari URL tersebut, situs mungkin tidak dapat diakses atau dibatasi"
		}

		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	evidences, err := h.searcher.SearchGrounding(c.Context(), article)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan pencarian bukti sedang mengalami gangguan",
		})
	}

	if len(evidences) == 0 {
		result := h.evaluator.Evaluate(article, evidences, &model.JevJudgments{})
		result.ShareHash = shareHash
		_ = h.cache.SetByHash(c.Context(), shareHash, result, config.UnverifiedCacheTTL)
		_ = h.cache.SetShare(c.Context(), shareHash, result)
		return c.JSON(result)
	}

	judgments, err := h.jevClient.Evaluate(c.Context(), article, evidences)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan evaluasi fakta sedang mengalami gangguan",
		})
	}

	result := h.evaluator.Evaluate(article, evidences, judgments)
	result.ShareHash = shareHash

	cacheTTL := config.DefaultCacheTTL
	if result.Verdict == config.VerdictUnverified {
		cacheTTL = config.UnverifiedCacheTTL
	}
	_ = h.cache.SetByHash(c.Context(), shareHash, result, cacheTTL)
	_ = h.cache.SetShare(c.Context(), shareHash, result)

	return c.JSON(result)
}

func (h *AnalyzeHandler) handleText(c *fiber.Ctx, clientIP string, rawText string) error {
	cleanedText, searchQuery := security.CleanBroadcastText(rawText)
	if cleanedText == "" || searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Teks pesan terlalu pendek atau tidak memuat klaim yang cukup setelah dibersihkan (minimal 25 karakter dan 4 kata)",
		})
	}

	hSum := sha256.Sum256([]byte(cleanedText))
	shareHash := hex.EncodeToString(hSum[:])

	if cached, err := h.cache.GetByHash(c.Context(), shareHash); err == nil && cached != nil {
		if cached.ShareHash == "" {
			cached.ShareHash = shareHash
		}
		return c.JSON(cached)
	}

	if allowed, remaining, err := h.cache.CheckAndSetCooldown(c.Context(), clientIP, config.DefaultCooldownDuration); err == nil && !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", remaining.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":            fmt.Sprintf("Mohon tunggu %.0f detik sebelum melakukan pengecekan berikutnya", remaining.Seconds()),
			"cooldown_seconds": int(remaining.Seconds()),
		})
	}

	evidences, err := h.searcher.SearchQuery(c.Context(), searchQuery)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan pencarian bukti sedang mengalami gangguan",
		})
	}

	title := searchQuery
	if len([]rune(title)) > 80 {
		title = string([]rune(title)[:80]) + "..."
	}

	if len(evidences) == 0 {
		result := &model.VerificationResult{
			Verdict:      config.VerdictUnverified,
			Confidence:   50.0,
			ArticleTitle: title,
			ArticleURL:   "",
			SourceDomain: "Input Teks Langsung",
			Summary:      "Tidak ditemukan bukti rujukan atau klarifikasi di sumber kredibel mengenai isi klaim konten ini.",
			RedFlags:     []string{"Klaim konten tidak memiliki rekam jejak rujukan di sumber kredibel."},
			Evidences:    []model.Evidence{},
			ShareHash:    shareHash,
			IsTextMode:   true,
			CleanedText:  cleanedText,
			IsSocial:     false,
			Platform:     "Teks / Konten",
		}

		_ = h.cache.SetByHash(c.Context(), shareHash, result, config.UnverifiedCacheTTL)
		_ = h.cache.SetShare(c.Context(), shareHash, result)
		return c.JSON(result)
	}

	if len(evidences) > 3 {
		evidences = evidences[:3]
	}
	for i := range evidences {
		runes := []rune(evidences[i].Snippet)
		if len(runes) > 200 {
			evidences[i].Snippet = string(runes[:200]) + "..."
		}
	}

	article := &model.Article{
		Title:        title,
		Content:      cleanedText,
		URL:          "",
		SourceDomain: "Input Teks Langsung",
		WordCount:    len(strings.Fields(cleanedText)),
		IsSocial:     false,
		Platform:     "Teks / Konten",
	}

	judgments, err := h.jevClient.Evaluate(c.Context(), article, evidences)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan evaluasi konten sedang mengalami gangguan",
		})
	}

	result := h.evaluator.Evaluate(article, evidences, judgments)
	result.ShareHash = shareHash
	result.IsTextMode = true
	result.CleanedText = cleanedText
	result.SourceDomain = "Input Teks Langsung"

	cacheTTL := config.DefaultCacheTTL
	if result.Verdict == config.VerdictUnverified {
		cacheTTL = config.UnverifiedCacheTTL
	}
	_ = h.cache.SetByHash(c.Context(), shareHash, result, cacheTTL)
	_ = h.cache.SetShare(c.Context(), shareHash, result)

	return c.JSON(result)
}

func (h *AnalyzeHandler) GetCheckByHash(c *fiber.Ctx) error {
	hash := strings.TrimSpace(c.Params("hash"))
	if hash == "" || len(hash) < 8 || len(hash) > 64 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hash verifikasi tidak valid",
		})
	}

	for _, r := range hash {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Hash verifikasi memiliki format yang tidak valid",
			})
		}
	}

	result, err := h.cache.GetShare(c.Context(), hash)
	if err != nil || result == nil {
		result, err = h.cache.GetByHash(c.Context(), hash)
	}

	if err != nil || result == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Data hasil verifikasi tidak ditemukan atau sudah kedaluwarsa",
		})
	}

	return c.JSON(result)
}
