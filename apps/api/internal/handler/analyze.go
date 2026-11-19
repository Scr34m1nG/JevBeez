package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"jevbeez/api/internal/cache"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/decision"
	"jevbeez/api/internal/jev"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/scraper"
	"jevbeez/api/internal/search"
	"jevbeez/api/internal/security"
)

type AnalyzeHandler struct {
	cache     *cache.CacheService
	scraper   *scraper.ScraperService
	searcher  *search.SearchManager
	jevClient *jev.Client
	evaluator *decision.Evaluator
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
	var req model.AnalyzeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid, pastikan JSON memiliki field 'url'",
		})
	}

	normalizedURL, err := security.ValidateAndNormalizeURL(req.URL)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("URL tidak valid: %s", err.Error()),
		})
	}

	if cached, err := h.cache.Get(c.Context(), normalizedURL); err == nil && cached != nil {
		return c.JSON(cached)
	}

	clientIP := c.IP()
	if allowed, remaining, err := h.cache.CheckAndSetCooldown(c.Context(), clientIP, config.DefaultCooldownDuration); err == nil && !allowed {
		c.Set("Retry-After", fmt.Sprintf("%.0f", remaining.Seconds()))
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error":             fmt.Sprintf("Mohon tunggu %.0f detik sebelum melakukan pengecekan berikutnya", remaining.Seconds()),
			"cooldown_seconds":  int(remaining.Seconds()),
		})
	}

	article, err := h.scraper.Extract(c.Context(), normalizedURL)
	if err != nil {
		if errors.Is(err, scraper.ErrInsufficientContent) || errors.Is(err, scraper.ErrEmptyArticle) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if errors.Is(err, security.ErrPrivateIPBlocked) || errors.Is(err, security.ErrInvalidIPAddress) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Akses ke URL ini diblokir demi alasan keamanan",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Gagal membaca artikel dari URL tersebut, situs mungkin tidak dapat diakses",
		})
	}

	evidences, err := h.searcher.SearchGrounding(c.Context(), article)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan pencarian bukti sedang mengalami gangguan",
		})
	}

	judgments, err := h.jevClient.Evaluate(c.Context(), article, evidences)
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Layanan evaluasi fakta sedang mengalami gangguan",
		})
	}

	result := h.evaluator.Evaluate(article, evidences, judgments)

	cacheTTL := config.DefaultCacheTTL
	if result.Verdict == config.VerdictUnverified {
		cacheTTL = config.UnverifiedCacheTTL
	}
	_ = h.cache.Set(c.Context(), normalizedURL, result, cacheTTL)

	return c.JSON(result)
}
