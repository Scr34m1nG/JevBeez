package scraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/security"
)

var (
	ErrUnsupportedSocialPlatform = errors.New("platform media sosial tidak didukung")
	ErrSocialCaptionNotFound     = errors.New("caption atau teks postingan tidak ditemukan")
)

type SocialScraper struct {
	client *http.Client
}

func NewSocialScraper() *SocialScraper {
	return &SocialScraper{
		client: security.NewSafeHTTPClient(),
	}
}

func NewSocialScraperWithClient(client *http.Client) *SocialScraper {
	return &SocialScraper{
		client: client,
	}
}

func DetectSocialPlatform(rawURL string) (string, bool) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host == "" {
		return "", false
	}

	host := strings.ToLower(parsed.Hostname())
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}

	switch {
	case host == "twitter.com" || host == "x.com" || strings.HasSuffix(host, ".twitter.com") || strings.HasSuffix(host, ".x.com"):
		return config.PlatformTwitter, true
	case host == "tiktok.com" || strings.HasSuffix(host, ".tiktok.com"):
		return config.PlatformTikTok, true
	case host == "youtube.com" || host == "youtu.be" || strings.HasSuffix(host, ".youtube.com"):
		return config.PlatformYouTube, true
	case host == "instagram.com" || strings.HasSuffix(host, ".instagram.com"):
		return config.PlatformInstagram, true
	case host == "threads.net" || strings.HasSuffix(host, ".threads.net"):
		return config.PlatformThreads, true
	case host == "facebook.com" || host == "fb.watch" || strings.HasSuffix(host, ".facebook.com"):
		return config.PlatformFacebook, true
	default:
		return "", false
	}
}

func (s *SocialScraper) Scrape(ctx context.Context, targetURL string, platform string) (*model.Article, error) {
	switch platform {
	case config.PlatformTwitter:
		return s.scrapeTwitterOEmbed(ctx, targetURL)
	case config.PlatformYouTube:
		return s.scrapeYouTubeOEmbed(ctx, targetURL)
	case config.PlatformTikTok, config.PlatformInstagram, config.PlatformThreads, config.PlatformFacebook:
		return s.scrapeOpenGraph(ctx, targetURL, platform)
	default:
		return nil, ErrUnsupportedSocialPlatform
	}
}

type twitterOEmbedResp struct {
	HTML       string `json:"html"`
	AuthorName string `json:"author_name"`
}

func (s *SocialScraper) scrapeTwitterOEmbed(ctx context.Context, targetURL string) (*model.Article, error) {
	apiURL := fmt.Sprintf("%s?url=%s", config.TwitterOEmbedEndpoint, url.QueryEscape(targetURL))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitter oembed error with status: %d", resp.StatusCode)
	}

	var res twitterOEmbedResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	content := ExtractParagraphText(res.HTML)
	if strings.TrimSpace(content) == "" {
		return nil, ErrSocialCaptionNotFound
	}

	title := "Postingan X / Twitter"
	if res.AuthorName != "" {
		title = fmt.Sprintf("Postingan X (@%s)", res.AuthorName)
	}

	parsedURL, _ := url.Parse(targetURL)
	words := len(strings.Fields(content))

	return &model.Article{
		Title:        title,
		Content:      content,
		URL:          targetURL,
		WordCount:    words,
		SourceDomain: parsedURL.Hostname(),
		IsSocial:     true,
		Platform:     config.PlatformTwitter,
	}, nil
}

type youTubeOEmbedResp struct {
	Title      string `json:"title"`
	AuthorName string `json:"author_name"`
}

func (s *SocialScraper) scrapeYouTubeOEmbed(ctx context.Context, targetURL string) (*model.Article, error) {
	apiURL := fmt.Sprintf("%s?url=%s&format=json", config.YouTubeOEmbedEndpoint, url.QueryEscape(targetURL))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube oembed error with status: %d", resp.StatusCode)
	}

	var res youTubeOEmbedResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if strings.TrimSpace(res.Title) == "" {
		return nil, ErrSocialCaptionNotFound
	}

	content := res.Title
	if res.AuthorName != "" {
		content = fmt.Sprintf("%s - Video oleh %s", res.Title, res.AuthorName)
	}

	parsedURL, _ := url.Parse(targetURL)
	words := len(strings.Fields(content))

	return &model.Article{
		Title:        res.Title,
		Content:      content,
		URL:          targetURL,
		WordCount:    words,
		SourceDomain: parsedURL.Hostname(),
		IsSocial:     true,
		Platform:     config.PlatformYouTube,
	}, nil
}

func (s *SocialScraper) scrapeOpenGraph(ctx context.Context, targetURL string, platform string) (*model.Article, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream social page returned status %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, config.MaxScrapeBodyBytes)
	tags, err := ExtractMetaTags(limitedReader)
	if err != nil {
		return nil, err
	}

	description := tags["og:description"]
	if description == "" {
		description = tags["twitter:description"]
	}
	if description == "" {
		description = tags["description"]
	}

	rawTitle := tags["og:title"]
	if rawTitle == "" {
		rawTitle = tags["twitter:title"]
	}

	title := fmt.Sprintf("Postingan %s", strings.Title(platform))
	if rawTitle != "" {
		if idx := strings.Index(rawTitle, " on Instagram:"); idx != -1 {
			title = fmt.Sprintf("Postingan Instagram (%s)", strings.TrimSpace(rawTitle[:idx]))
		} else if idx := strings.Index(rawTitle, " on Threads:"); idx != -1 {
			title = fmt.Sprintf("Postingan Threads (%s)", strings.TrimSpace(rawTitle[:idx]))
		} else {
			title = rawTitle
		}
	}

	cleanContent := extractCleanCaption(description)
	if cleanContent == "" && rawTitle != "" {
		cleanContent = extractCleanCaption(rawTitle)
	}
	if cleanContent == "" {
		cleanContent = description
	}
	if cleanContent == "" {
		cleanContent = rawTitle
	}

	cleanContent = strings.TrimSpace(cleanContent)
	if cleanContent == "" {
		return nil, ErrSocialCaptionNotFound
	}

	parsedURL, _ := url.Parse(targetURL)
	words := len(strings.Fields(cleanContent))

	return &model.Article{
		Title:        title,
		Content:      cleanContent,
		URL:          targetURL,
		WordCount:    words,
		SourceDomain: parsedURL.Hostname(),
		IsSocial:     true,
		Platform:     platform,
	}, nil
}

func extractCleanCaption(desc string) string {
	desc = strings.TrimSpace(desc)
	for _, marker := range []string{": \"", ": \u201c", ": '"} {
		if idx := strings.Index(desc, marker); idx != -1 {
			candidate := strings.TrimSpace(desc[idx+len(marker):])
			// Strip trailing quotes or quote-dot combinations
			candidate = strings.TrimSuffix(candidate, "\".")
			candidate = strings.TrimSuffix(candidate, "\u201d.")
			candidate = strings.TrimRight(candidate, "\"\u201d\u201c'")
			candidate = strings.TrimSpace(candidate)
			if len(candidate) > 0 {
				return candidate
			}
		}
	}
	return desc
}
