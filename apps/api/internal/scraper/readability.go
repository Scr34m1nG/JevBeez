package scraper

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-shiori/go-readability"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
	"jevbeez/api/internal/security"
)

type LocalScraper struct {
	client *http.Client
}

func NewLocalScraper() *LocalScraper {
	return &LocalScraper{
		client: security.NewSafeHTTPClient(),
	}
}

func (s *LocalScraper) Scrape(ctx context.Context, targetURL string) (*model.Article, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("upstream server returned error status")
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		ct := strings.ToLower(contentType)
		if !strings.Contains(ct, "text/html") && !strings.Contains(ct, "application/xhtml+xml") {
			return nil, errors.New("unsupported content type: must be html")
		}
	}

	limitedReader := io.LimitReader(resp.Body, config.MaxScrapeBodyBytes)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	parsedURL, _ := url.Parse(targetURL)
	article, err := readability.FromReader(bytes.NewReader(bodyBytes), parsedURL)
	if err != nil {
		return nil, err
	}

	cleanText := strings.TrimSpace(article.TextContent)
	words := len(strings.Fields(cleanText))

	return &model.Article{
		Title:        strings.TrimSpace(article.Title),
		Content:      cleanText,
		URL:          targetURL,
		WordCount:    words,
		SourceDomain: parsedURL.Hostname(),
	}, nil
}
