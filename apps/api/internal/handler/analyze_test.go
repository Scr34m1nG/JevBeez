package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
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
	app := fiber.New()
	cacheSvc := &cache.CacheService{}
	scraperSvc := scraper.NewScraperService("")
	searcher := &search.SearchManager{}
	evaluator := decision.NewEvaluator()

	handler := NewAnalyzeHandler(cacheSvc, scraperSvc, searcher, nil, evaluator)
	app.Post("/api/analyze", handler.Handle)
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
