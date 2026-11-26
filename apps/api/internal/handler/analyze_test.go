package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"jevbeez/api/internal/cache"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/decision"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/scraper"
	"jevbeez/api/internal/search"
)

func setupTestApp() (*fiber.App, *AnalyzeHandler) {
	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		TrustedProxies:          []string{"127.0.0.1", "0.0.0.0", "192.168.0.0/16"},
	})
	cacheSvc := &cache.CacheService{}
	scraperSvc := scraper.NewScraperService("")
	searcher := &search.SearchManager{}
	evaluator := decision.NewEvaluator()

	handler := NewAnalyzeHandler(cacheSvc, scraperSvc, searcher, nil, evaluator)
	app.Post("/api/analyze", handler.Handle)
	app.Get("/api/analyze/stream", handler.StreamHandle)
	app.Get("/api/check/:hash", handler.GetCheckByHash)

	return app, handler
}

func TestAnalyzeHandler_PhishingZeroContactFastVerdict(t *testing.T) {
	app, _ := setupTestApp()

	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL:      "https://bca.co.id.login-promo.xyz/auth",
		DeepDive: false,
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !result.IsThreat {
		t.Errorf("expected IsThreat to be true")
	}
	if result.Verdict != config.VerdictPhishing {
		t.Errorf("expected verdict %s, got %s", config.VerdictPhishing, result.Verdict)
	}
	if result.TargetBrand != "bca" {
		t.Errorf("expected target brand 'bca', got %s", result.TargetBrand)
	}
	if result.IsDeepDive {
		t.Errorf("expected IsDeepDive to be false")
	}
	if result.ShareHash == "" {
		t.Errorf("expected non-empty ShareHash")
	}
}

func TestAnalyzeHandler_MalwareZeroContactFastVerdict(t *testing.T) {
	app, _ := setupTestApp()

	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL:      "https://undangan-pernikahan-digital.xyz/surat_undangan.apk",
		DeepDive: false,
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !result.IsThreat {
		t.Errorf("expected IsThreat to be true")
	}
	if result.Verdict != config.VerdictMalware {
		t.Errorf("expected verdict %s, got %s", config.VerdictMalware, result.Verdict)
	}
}

func TestAnalyzeHandler_DeepDiveMode(t *testing.T) {
	app, _ := setupTestApp()

	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL:      "https://bca.co.id.login-promo.xyz/auth",
		DeepDive: true,
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !result.IsThreat {
		t.Errorf("expected IsThreat to be true")
	}
	if !result.IsDeepDive {
		t.Errorf("expected IsDeepDive to be true")
	}
}

func TestAnalyzeHandler_EmptyRequest(t *testing.T) {
	app, _ := setupTestApp()

	payload, _ := json.Marshal(model.AnalyzeRequest{})
	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 for empty request, got %d", resp.StatusCode)
	}
}

func TestAnalyzeHandler_TextModeTooShort(t *testing.T) {
	app, _ := setupTestApp()

	payload, _ := json.Marshal(model.AnalyzeRequest{
		Text: "Pendek banget",
	})
	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 for short text, got %d", resp.StatusCode)
	}
}

func TestAnalyzeHandler_GetCheckByHash_InvalidFormat(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/check/invalid_char_!@#$", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid hash characters, got %d", resp.StatusCode)
	}
}

func TestAnalyzeHandler_GetCheckByHash_NotFound(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/check/abcdef1234567890", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Fatalf("expected status 404 for missing hash in nil cache, got %d", resp.StatusCode)
	}
}

func TestAnalyzeHandler_WrappedScamLinkReturnsScam(t *testing.T) {
	app, _ := setupTestApp()
	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL: "https://www.google.com/pagead/aclk?sa=L&adurl=https%3A%2F%2Fvk.com%2Faway.php%3Fto%3Dhttps%3A%2F%2Fbaf9105c4584fac832929f959728944d.cavemin.com%2Fsttodm%2F",
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if result.Verdict != config.VerdictScam {
		t.Fatalf("expected SCAM, got %s", result.Verdict)
	}
	if len(result.RedirectHops) == 0 {
		t.Error("expected redirect chain recorded")
	}
}

func TestAnalyzeHandler_WrappedEcommerceReturnsAds(t *testing.T) {
	app, _ := setupTestApp()
	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL: "https://www.google.com/pagead/aclk?sa=L&adurl=https%3A%2F%2Fwww.ebay.com%2Fitm%2F123",
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if result.Verdict != config.VerdictAds {
		t.Fatalf("expected ADS, got %s", result.Verdict)
	}
	if result.IsThreat {
		t.Error("ADS must not be a threat")
	}
	if result.ResolvedURL != "https://www.ebay.com/itm/123" {
		t.Errorf("expected resolved ebay url, got %s", result.ResolvedURL)
	}
}

func TestAnalyzeHandler_WrappedBrandSpoofStaysPhishing(t *testing.T) {
	app, _ := setupTestApp()
	payload, _ := json.Marshal(model.AnalyzeRequest{
		URL: "https://www.google.com/pagead/aclk?sa=L&adurl=https%3A%2F%2Fbca.co.id.login-promo.xyz%2Fauth",
	})

	req := httptest.NewRequest("POST", "/api/analyze", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result model.VerificationResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if result.Verdict != config.VerdictPhishing {
		t.Fatalf("expected PHISHING, got %s", result.Verdict)
	}
}

func TestAnalyzeHandler_StreamPhishingFastVerdict(t *testing.T) {
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/api/analyze/stream?url=https://bca.co.id.login-promo.xyz/auth", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/event-stream") {
		t.Errorf("expected text/event-stream, got %s", ct)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	if !strings.Contains(bodyStr, "event: progress") {
		t.Errorf("expected event: progress in stream, got %s", bodyStr)
	}
	if !strings.Contains(bodyStr, "event: verdict") {
		t.Errorf("expected event: verdict in stream, got %s", bodyStr)
	}
	if !strings.Contains(bodyStr, config.VerdictPhishing) {
		t.Errorf("expected PHISHING in stream, got %s", bodyStr)
	}
}

func TestAnalyzeHandler_StreamConcurrencyLimit(t *testing.T) {
	app, handler := setupTestApp()

	// Simulate 2 active streams by directly incrementing atomic counter for client IP
	clientIP := "192.168.1.100"
	val, _ := handler.activeStreams.LoadOrStore(clientIP, new(int32))
	activeCount := val.(*int32)
	activeCountStore := activeCount
	*activeCountStore = 2

	req := httptest.NewRequest("GET", "/api/analyze/stream?url=https://news.com/art", nil)
	req.Header.Set("X-Forwarded-For", clientIP)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests when 2 streams active, got %d", resp.StatusCode)
	}
}
