package scraper

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"jevbeez/api/internal/config"
)

func TestDetectSocialPlatform(t *testing.T) {
	tests := []struct {
		url          string
		wantPlatform string
		wantSocial   bool
	}{
		{"https://x.com/elonmusk/status/123456789", config.PlatformTwitter, true},
		{"https://twitter.com/jokowi/status/987654321", config.PlatformTwitter, true},
		{"https://www.tiktok.com/@creator/video/12345", config.PlatformTikTok, true},
		{"https://youtu.be/dQw4w9WgXcQ", config.PlatformYouTube, true},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", config.PlatformYouTube, true},
		{"https://www.instagram.com/p/Cxyz123/", config.PlatformInstagram, true},
		{"https://www.threads.net/@user/post/Dxyz456", config.PlatformThreads, true},
		{"https://facebook.com/story.php?story_fbid=123", config.PlatformFacebook, true},
		{"https://fb.watch/abc123xyz/", config.PlatformFacebook, true},
		{"https://news.detik.com/berita/123-berita-nasional", "", false},
		{"https://kompas.com/tren/read/123", "", false},
		{"not-a-valid-url", "", false},
	}

	for _, tc := range tests {
		platform, isSocial := DetectSocialPlatform(tc.url)
		if isSocial != tc.wantSocial {
			t.Errorf("url '%s': expected isSocial %v, got %v", tc.url, tc.wantSocial, isSocial)
		}
		if platform != tc.wantPlatform {
			t.Errorf("url '%s': expected platform '%s', got '%s'", tc.url, tc.wantPlatform, platform)
		}
	}
}

func TestSocialScraper_ParseInstagramCaption(t *testing.T) {
	mockHTML := `<!DOCTYPE html>
<html>
<head>
	<meta property="og:title" content="TV NASIONAL on Instagram: &#34;Dividio ini terungkap kisah nyata tentang proses vaksinasi memicu alergi.&#34;" />
	<meta property="og:description" content="61K likes, 4,582 comments - nasional.tv5 on October 10, 2026: &#34;Dividio ini terungkap kisah nyata tentang proses vaksinasi memicu alergi.&#34;. " />
</head>
<body></body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)" {
			// Jika bukan social bot UA, Instagram mengembalikan halaman kosong
			w.Write([]byte("<html><head><title>Instagram</title></head><body></body></html>"))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	article, err := scraper.scrapeOpenGraph(context.Background(), server.URL, config.PlatformInstagram)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedCaption := "Dividio ini terungkap kisah nyata tentang proses vaksinasi memicu alergi."
	if !strings.Contains(article.Content, expectedCaption) {
		t.Errorf("expected content to contain '%s', got '%s'", expectedCaption, article.Content)
	}
	if strings.Contains(article.Content, "61K likes") {
		t.Errorf("expected likes prefix to be stripped from caption, got '%s'", article.Content)
	}
	if !strings.Contains(article.Title, "TV NASIONAL") {
		t.Errorf("expected title to identify author, got '%s'", article.Title)
	}
	if !article.IsSocial {
		t.Errorf("expected article.IsSocial to be true")
	}
}
