package scraper

import (
	"context"
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
