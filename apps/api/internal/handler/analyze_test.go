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
