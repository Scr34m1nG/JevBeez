package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

type FirecrawlScraper struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

type firecrawlScrapeRequest struct {
	URL     string   `json:"url"`
	Formats []string `json:"formats"`
}

type firecrawlScrapeResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Markdown string `json:"markdown"`
		Metadata struct {
			Title string `json:"title"`
		} `json:"metadata"`
	} `json:"data"`
	Error string `json:"error"`
}

func NewFirecrawlScraper(apiKey string) *FirecrawlScraper {
	return &FirecrawlScraper{
		apiKey:   apiKey,
		client:   config.NewTrustedHTTPClient(),
		endpoint: config.FirecrawlScrapeEndpoint,
	}
}

func (s *FirecrawlScraper) Scrape(ctx context.Context, targetURL string) (*model.Article, error) {
	if s.apiKey == "" {
		return nil, errors.New("firecrawl api key not configured")
	}

	payload := firecrawlScrapeRequest{
		URL:     targetURL,
		Formats: []string{"markdown"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	endpoint := s.endpoint
	if endpoint == "" {
		endpoint = config.FirecrawlScrapeEndpoint
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("firecrawl scrape failed with status code %d", resp.StatusCode)
	}

	var res firecrawlScrapeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("firecrawl scrape error: %s", res.Error)
	}

	cleanText := strings.TrimSpace(res.Data.Markdown)
	words := len(strings.Fields(cleanText))
	parsedURL, _ := url.Parse(targetURL)

	return &model.Article{
		Title:        strings.TrimSpace(res.Data.Metadata.Title),
		Content:      cleanText,
		URL:          targetURL,
		WordCount:    words,
		SourceDomain: parsedURL.Hostname(),
	}, nil
}
