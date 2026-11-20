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
		{"https://www.threads.com/@user/post/Dxyz456", config.PlatformThreads, true},
		{"https://www.threads.com/share/D7VmU4Wtk", config.PlatformThreads, true},
		{"https://threads.com/t/abc123", config.PlatformThreads, true},
		{"https://facebook.com/story.php?story_fbid=123", config.PlatformFacebook, true},
		{"https://fb.watch/abc123xyz/", config.PlatformFacebook, true},
		{"https://www.facebook.com/share/v/19z", config.PlatformFacebook, true},
		{"https://fb.com/share/v/123", config.PlatformFacebook, true},
		{"https://m.facebook.com/watch/?v=123", config.PlatformFacebook, true},
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

func TestSocialScraper_FacebookLoginWall(t *testing.T) {
	mockHTML := `<!DOCTYPE html>
<html>
<head>
	<meta property="og:title" content="Masuk atau Daftar untuk Melihat" />
	<meta property="og:description" content="Lihat postingan, foto, dan lainnya di Facebook." />
</head>
<body></body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	_, err := scraper.scrapeOpenGraph(context.Background(), server.URL, config.PlatformFacebook)
	if err == nil {
		t.Fatalf("expected error for facebook login wall, got nil")
	}
	if err != ErrSocialLoginRequired {
		t.Errorf("expected ErrSocialLoginRequired, got: %v", err)
	}
}

func TestSocialScraper_UnavailablePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	_, err := scraper.scrapeOpenGraph(context.Background(), server.URL, config.PlatformFacebook)
	if err != ErrSocialContentUnavailable {
		t.Errorf("expected ErrSocialContentUnavailable, got: %v", err)
	}
}

func TestIsSocialLoginWall(t *testing.T) {
	cases := []struct {
		title string
		desc  string
		want  bool
	}{
		{"Masuk atau Daftar untuk Melihat", "Lihat postingan, foto, dan lainnya di Facebook.", true},
		{"Log in or sign up to view", "See posts, photos and more on Facebook.", true},
		{"Login • Instagram", "Welcome back to Instagram.", true},
		{"Berita Hoax Vaksin", "Ini adalah caption lengkap mengenai klaim vaksin tertentu.", false},
	}

	for _, tc := range cases {
		got := IsSocialLoginWall(tc.title, tc.desc)
		if got != tc.want {
			t.Errorf("IsSocialLoginWall(%q, %q) = %v, want %v", tc.title, tc.desc, got, tc.want)
		}
	}
}

func TestSocialScraper_ParseThreadsCaption(t *testing.T) {
	mockHTML := `<!DOCTYPE html>
<html>
<head>
	<meta property="og:title" content="Rupnath Shaw | Data Analyst (&#064;iamrupnath) on Threads" />
	<meta property="og:description" content="Gemini 4 Sodium beating Claude Opus 5.5 and GPT-6 Astra. Google just dropped 4 Gemini 4 models." />
</head>
<body></body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(mockHTML))
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	article, err := scraper.scrapeOpenGraph(context.Background(), server.URL, config.PlatformThreads)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(article.Title, "@iamrupnath") {
		t.Errorf("expected title to have decoded '@iamrupnath', got: %s", article.Title)
	}
	if !strings.Contains(article.Content, "Gemini 4 Sodium") {
		t.Errorf("expected content to have post text, got: %s", article.Content)
	}
}

func TestSocialScraper_TwitterFxFallback(t *testing.T) {
	fxJSON := `{
		"code": 200,
		"message": "OK",
		"tweet": {
			"text": "Klaim bantuan bansos tunai tahap 4 resmi dibuka untuk seluruh warga.",
			"author": {
				"name": "Kemensos RI",
				"screen_name": "kemensosri"
			},
			"quote": {
				"text": "Waspada hoaks pendaftaran bansos via tautan telegram tidak resmi.",
				"author": {
					"name": "Saber Hoaks",
					"screen_name": "saberhoaks"
				}
			}
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fxJSON))
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	scraper.fxTwitterEndpoint = server.URL
	article, err := scraper.scrapeTwitterFx(context.Background(), server.URL+"/kemensosri/status/123456")
	if err != nil {
		t.Fatalf("unexpected error from scrapeTwitterFx: %v", err)
	}

	if !strings.Contains(article.Title, "@kemensosri") {
		t.Errorf("expected title to contain @kemensosri, got: %s", article.Title)
	}
	if !strings.Contains(article.Content, "Klaim bantuan bansos") {
		t.Errorf("expected content to contain main tweet text, got: %s", article.Content)
	}
	if !strings.Contains(article.Content, "[Kutipan @saberhoaks]") {
		t.Errorf("expected content to contain quote tweet text, got: %s", article.Content)
	}
}

func TestSocialScraper_TwitterFxPrivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"code": 401, "message": "Private account"}`))
	}))
	defer server.Close()

	scraper := NewSocialScraperWithClient(server.Client())
	scraper.fxTwitterEndpoint = server.URL
	_, err := scraper.scrapeTwitterFx(context.Background(), server.URL+"/privateuser/status/123456")
	if err != ErrSocialLoginRequired {
		t.Errorf("expected ErrSocialLoginRequired, got: %v", err)
	}
}
