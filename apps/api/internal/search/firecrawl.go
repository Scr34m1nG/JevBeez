package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

type FirecrawlSearchProvider struct {
	apiKey string
	client *http.Client
}

type firecrawlSearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

type firecrawlSearchResultItem struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Markdown    string `json:"markdown"`
}

type firecrawlSearchResponse struct {
	Success bool                        `json:"success"`
	Data    []firecrawlSearchResultItem `json:"data"`
	Error   string                      `json:"error"`
}

func NewFirecrawlSearchProvider(apiKey string) *FirecrawlSearchProvider {
	return &FirecrawlSearchProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: config.HTTPClientTimeout},
	}
}

func (p *FirecrawlSearchProvider) Name() string {
	return "firecrawl"
}

func (p *FirecrawlSearchProvider) Search(ctx context.Context, query string) ([]model.Evidence, error) {
	payload := firecrawlSearchRequest{
		Query: query,
		Limit: 5,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.FirecrawlSearchEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("firecrawl search returned status %d", resp.StatusCode)
	}

	var res firecrawlSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("firecrawl search error: %s", res.Error)
	}

	var evidences []model.Evidence
	for _, item := range res.Data {
		domain := ""
		if u, err := url.Parse(item.URL); err == nil {
			domain = u.Hostname()
		}
		snippet := item.Description
		if snippet == "" {
			snippet = item.Markdown
		}
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		evidences = append(evidences, model.Evidence{
			Source:  domain,
			Title:   item.Title,
			Snippet: snippet,
			URL:     item.URL,
		})
	}

	return evidences, nil
}
