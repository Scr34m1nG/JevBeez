package search

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"jevbeez/api/internal/model"
)

var ErrNoSearchProvider = errors.New("tidak ada search provider aktif yang tersedia")

type SearchManager struct {
	providers []SearchProvider
	counter   uint64
}

func NewSearchManager(tavilyKey, firecrawlKey string) (*SearchManager, error) {
	var active []SearchProvider

	if tavilyKey != "" {
		active = append(active, NewTavilyProvider(tavilyKey))
	}
	if firecrawlKey != "" {
		active = append(active, NewFirecrawlSearchProvider(firecrawlKey))
	}

	if len(active) == 0 {
		return nil, ErrNoSearchProvider
	}

	return &SearchManager{
		providers: active,
	}, nil
}

func (m *SearchManager) SearchGrounding(ctx context.Context, article *model.Article) ([]model.Evidence, error) {
	query := m.buildQuery(article)
	if query == "" {
		return nil, errors.New("query pencarian tidak dapat dibentuk")
	}

	n := len(m.providers)
	startIdx := int(atomic.AddUint64(&m.counter, 1) % uint64(n))

	var lastErr error
	var successfulProviders int
	for i := 0; i < n; i++ {
		idx := (startIdx + i) % n
		provider := m.providers[idx]

		results, err := provider.Search(ctx, query)
		if err == nil {
			successfulProviders++
			if len(results) > 0 {
				return results, nil
			}
		} else {
			lastErr = err
		}
	}

	if successfulProviders == 0 && lastErr != nil {
		return nil, fmt.Errorf("semua search provider gagal: %w", lastErr)
	}

	return []model.Evidence{}, nil
}

func (m *SearchManager) buildQuery(article *model.Article) string {
	if (article.IsSocial || strings.TrimSpace(article.Title) == "") && strings.TrimSpace(article.Content) != "" {
		content := strings.TrimSpace(article.Content)
		for _, delim := range []string{"\n", ".", "!"} {
			if idx := strings.Index(content, delim); idx > 15 {
				content = content[:idx]
				break
			}
		}
		content = strings.TrimSpace(content)
		runes := []rune(content)
		if len(runes) > 120 {
			content = string(runes[:120])
		}
		if len(content) >= 10 {
			return content
		}
	}

	title := strings.TrimSpace(article.Title)
	for _, separator := range []string{" - ", " | ", " – "} {
		if idx := strings.LastIndex(title, separator); idx > 20 {
			title = title[:idx]
		}
	}
	title = strings.TrimSpace(title)
	runes := []rune(title)
	if len(runes) > 120 {
		title = string(runes[:120])
	}
	return title
}
