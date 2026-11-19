package jev

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"jevbeez/api/internal/model"
)

func TestClient_Evaluate(t *testing.T) {
	ctx := context.Background()
	article := &model.Article{
		Title:   "Test Article",
		Content: "Test content...",
	}

	t.Run("returns error when API response contains error field", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"error": "quota exceeded"}`))
		}))
		defer ts.Close()

		client := &Client{
			apiKey: "test-key",
			client: ts.Client(),
		}
		// override endpoint via test or request
		_, err := client.evaluateWithURL(ctx, article, nil, ts.URL)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error when answers are empty", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"model": "jev-latest", "answers": {}}`))
		}))
		defer ts.Close()

		client := &Client{
			apiKey: "test-key",
			client: ts.Client(),
		}
		_, err := client.evaluateWithURL(ctx, article, nil, ts.URL)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
