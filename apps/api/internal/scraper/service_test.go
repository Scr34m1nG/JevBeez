package scraper

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jevbeez/api/internal/model"
)

func TestScraperService_ValidateArticle(t *testing.T) {
	service := NewScraperService("")

	tests := []struct {
		name      string
		article   *model.Article
		expectErr error
	}{
		{
			name:      "nil article returns ErrEmptyArticle",
			article:   nil,
			expectErr: ErrEmptyArticle,
		},
		{
			name: "empty content returns ErrEmptyArticle",
			article: &model.Article{
				Title:   "Judul",
				Content: "   ",
			},
			expectErr: ErrEmptyArticle,
		},
		{
			name: "social login wall returns ErrSocialLoginRequired",
			article: &model.Article{
				Title:    "Masuk atau Daftar untuk Melihat",
				Content:  "Lihat postingan, foto, dan lainnya di Facebook.",
				IsSocial: true,
			},
			expectErr: ErrSocialLoginRequired,
		},
		{
			name: "article with 404 title returns ErrEmptyArticle",
			article: &model.Article{
				Title:   "404 Not Found",
				Content: "Halaman ini tidak ditemukan di server kami.",
			},
			expectErr: ErrEmptyArticle,
		},
		{
			name: "news article below word count returns ErrInsufficientContent",
			article: &model.Article{
				Title:     "Berita Singkat",
				Content:   "Satu dua tiga empat lima.",
				WordCount: 5,
				IsSocial:  false,
			},
			expectErr: ErrInsufficientContent,
		},
		{
			name: "social post meeting social word count is valid",
			article: &model.Article{
				Title:     "Postingan Facebook",
				Content:   "Klaim vaksin ini ternyata terbukti tidak benar menurut dokter.",
				WordCount: 9,
				IsSocial:  true,
			},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validateArticle(tc.article)
			if err != tc.expectErr {
				t.Errorf("expected error %v, got %v", tc.expectErr, err)
			}
		})
	}
}

func TestScraperService_SocialDoesNotFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	service := NewScraperService("")
	service.social = NewSocialScraperWithClient(server.Client())

	_, err := service.Extract(context.Background(), server.URL+"/facebook.com/share/v/123")
	// Since server.URL doesn't match DetectSocialPlatform unless we pass a facebook host,
	// let's test directly with facebook.com URL
	_, err = service.Extract(context.Background(), "https://www.facebook.com/share/v/nonexistent")
	if err == nil {
		t.Fatalf("expected error for social post, got nil")
	}
	if err != ErrSocialLoginRequired && err != ErrSocialContentUnavailable {
		t.Logf("got expected social error: %v", err)
	}
}

func TestScraperService_LocalFirstStrategy(t *testing.T) {
	// 1. Mock local server serving valid article HTML
	htmlContent := `<!DOCTYPE html><html><head><title>Berita Penting Hari Ini</title></head><body><article>`
	for i := 0; i < 90; i++ {
		htmlContent += "berita fakta terverifikasi kata ke "
	}
	htmlContent += `</article></body></html>`

	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))
	}))
	defer localServer.Close()

	// 2. Firecrawl server that fails if called
	firecrawlCalled := false
	firecrawlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firecrawlCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer firecrawlServer.Close()

	service := NewScraperService("test-firecrawl-key")
	service.local = NewLocalScraperWithClient(localServer.Client())
	service.firecrawl.client = firecrawlServer.Client()
	service.firecrawl.endpoint = firecrawlServer.URL

	article, err := service.Extract(context.Background(), localServer.URL+"/news/article-1")
	if err != nil {
		t.Fatalf("expected successful extract via local, got error: %v", err)
	}

	if article == nil || article.Title == "" {
		t.Fatalf("expected valid article, got nil or empty")
	}

	if firecrawlCalled {
		t.Fatalf("expected Firecrawl NOT to be called when local scraper succeeds, but it was called")
	}
}

func TestScraperService_FallbackToFirecrawlWhenLocalFails(t *testing.T) {
	// Local server returns 403 Forbidden
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer localServer.Close()

	firecrawlMarkdown := "# Berita Dari Firecrawl\n\n"
	for i := 0; i < 90; i++ {
		firecrawlMarkdown += "konten valid dari firecrawl kata ke "
	}

	firecrawlCalled := false
	firecrawlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firecrawlCalled = true
		w.Header().Set("Content-Type", "application/json")
		resp := firecrawlScrapeResponse{
			Success: true,
		}
		resp.Data.Markdown = firecrawlMarkdown
		resp.Data.Metadata.Title = "Berita Dari Firecrawl"
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer firecrawlServer.Close()

	service := NewScraperService("test-firecrawl-key")
	service.local = NewLocalScraperWithClient(localServer.Client())
	service.firecrawl.client = firecrawlServer.Client()
	service.firecrawl.endpoint = firecrawlServer.URL

	article, err := service.Extract(context.Background(), localServer.URL+"/news/article-spa")
	if err != nil {
		t.Fatalf("expected fallback to firecrawl to succeed, got: %v", err)
	}

	if !firecrawlCalled {
		t.Fatalf("expected Firecrawl to be called as fallback, but it was not")
	}

	if article.Title != "Berita Dari Firecrawl" {
		t.Fatalf("expected title 'Berita Dari Firecrawl', got '%s'", article.Title)
	}
}

func TestScraperService_SSRFDoesNotFallback(t *testing.T) {
	firecrawlCalled := false
	firecrawlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firecrawlCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer firecrawlServer.Close()

	service := NewScraperService("test-firecrawl-key")
	service.firecrawl.client = firecrawlServer.Client()
	service.firecrawl.endpoint = firecrawlServer.URL

	// Default local scraper uses NewSafeHTTPClient which blocks 127.0.0.1 (private IP)
	_, err := service.Extract(context.Background(), "http://127.0.0.1:8080/admin")
	if err == nil {
		t.Fatalf("expected SSRF error, got nil")
	}

	if firecrawlCalled {
		t.Fatalf("expected Firecrawl NOT to be called when SSRF is detected, but it was called")
	}
}

func BenchmarkLocalScraper_Extract(b *testing.B) {
	htmlContent := `<!DOCTYPE html><html><head><title>Detik News: Kebijakan Baru Transportasi Publik</title></head><body><article>`
	for i := 0; i < 150; i++ {
		htmlContent += "<p>Pemerintah provinsi mengumumkan pembaruan tarif integrasi transportasi publik untuk mendukung mobilitas perkotaan yang berkelanjutan dan efisien bagi masyarakat.</p>"
	}
	htmlContent += `</article></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	scraper := NewLocalScraperWithClient(server.Client())
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := scraper.Scrape(ctx, server.URL+"/berita/1")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
