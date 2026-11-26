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

type TavilyProvider struct {
	apiKey string
	client *http.Client
}

type tavilyRequest struct {
	APIKey      string `json:"api_key"`
	Query       string `json:"query"`
	Topic       string `json:"topic,omitempty"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type tavilyResultItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type tavilyResponse struct {
	Results []tavilyResultItem `json:"results"`
}

func NewTavilyProvider(apiKey string) *TavilyProvider {
	return &TavilyProvider{
		apiKey: apiKey,
		client: config.NewTrustedHTTPClient(),
	}
}

func (p *TavilyProvider) Name() string {
	return "tavily"
}

func (p *TavilyProvider) Search(ctx context.Context, query string) ([]model.Evidence, error) {
	payload := tavilyRequest{
		APIKey:      p.apiKey,
		Query:       query,
		Topic:       "news",
		SearchDepth: "basic",
		MaxResults:  5,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.TavilyEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily search returned status %d", resp.StatusCode)
	}

	var res tavilyResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	var evidences []model.Evidence
	for _, item := range res.Results {
		domain := ""
		if u, err := url.Parse(item.URL); err == nil {
			domain = u.Hostname()
		}
		evidences = append(evidences, model.Evidence{
			Source:  domain,
			Title:   item.Title,
			Snippet: item.Content,
			URL:     item.URL,
		})
	}

	return evidences, nil
}
