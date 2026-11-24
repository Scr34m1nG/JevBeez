package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsShortener(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://s.id/promo-nikah", true},
		{"http://bit.ly/undangan-apk", true},
		{"https://tinyurl.com/secret", true},
		{"https://link.id/bca-promo", true},
		{"https://kompas.com/news/123", false},
		{"https://detik.com/read/2026", false},
		{"invalid-url", false},
	}

	for _, tt := range tests {
		res := IsShortener(tt.url)
		if res != tt.expected {
			t.Errorf("IsShortener(%s) = %v; expected %v", tt.url, res, tt.expected)
		}
	}
}

func TestResolveShortener_NotShortener(t *testing.T) {
	ctx := context.Background()
	url := "https://example.com/normal-article"
	resolved, hops, err := ResolveShortener(ctx, url, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != url {
		t.Errorf("expected %s, got %s", url, resolved)
	}
	if len(hops) != 0 {
		t.Errorf("expected 0 hops, got %d", len(hops))
	}
}

func TestResolveShortener_RedirectChain(t *testing.T) {
	// Server mock for final destination
	finalServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Final Destination"))
	}))
	defer finalServer.Close()

	// Server mock for hop 1
	hop1Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, finalServer.URL, http.StatusFound)
	}))
	defer hop1Server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Custom resolver using local test server
	resolved, hops, err := resolveShortenerWithClient(ctx, hop1Server.URL, 5, hop1Server.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != finalServer.URL {
		t.Errorf("expected %s, got %s", finalServer.URL, resolved)
	}
	if len(hops) != 1 {
		t.Errorf("expected 1 hop, got %d", len(hops))
	}
}

func TestResolveShortener_MaxHopsExceeded(t *testing.T) {
	var loopServer *httptest.Server
	loopServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, loopServer.URL+"/next", http.StatusFound)
	}))
	defer loopServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, _, err := resolveShortenerWithClient(ctx, loopServer.URL, 2, loopServer.Client())
	if err == nil {
		t.Fatalf("expected error for max hops exceeded, got nil")
	}
}
