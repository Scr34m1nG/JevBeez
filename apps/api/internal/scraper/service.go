package scraper

import (
	"context"
	"errors"
	"strings"

	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

var (
	ErrEmptyArticle        = errors.New("halaman web tidak memiliki judul atau konten yang dapat dibaca")
	ErrInsufficientContent = errors.New("konten artikel terlalu pendek atau dilindungi paywall")
)

type ScraperService struct {
	firecrawl *FirecrawlScraper
	local     *LocalScraper
	social    *SocialScraper
}

func NewScraperService(firecrawlKey string) *ScraperService {
	return &ScraperService{
		firecrawl: NewFirecrawlScraper(firecrawlKey),
		local:     NewLocalScraper(),
		social:    NewSocialScraper(),
	}
}

func (s *ScraperService) Extract(ctx context.Context, targetURL string) (*model.Article, error) {
	var article *model.Article
	var err error

	platform, isSocial := DetectSocialPlatform(targetURL)
	if isSocial {
		article, err = s.social.Scrape(ctx, targetURL, platform)
		if article != nil && err == nil && s.validateArticle(article) == nil {
			return article, nil
		}
	}

	if s.firecrawl.apiKey != "" {
		article, err = s.firecrawl.Scrape(ctx, targetURL)
	}

	if article == nil || err != nil {
		article, err = s.local.Scrape(ctx, targetURL)
		if err != nil {
			return nil, err
		}
	}

	if isSocial && article != nil {
		article.IsSocial = true
		article.Platform = platform
	}

	if err := s.validateArticle(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *ScraperService) validateArticle(article *model.Article) error {
	if article == nil || article.Content == "" {
		return ErrEmptyArticle
	}

	lowerTitle := strings.ToLower(article.Title)
	if strings.Contains(lowerTitle, "404 not found") ||
		strings.Contains(lowerTitle, "page not found") ||
		strings.Contains(lowerTitle, "access denied") ||
		strings.Contains(lowerTitle, "attention required") {
		return ErrEmptyArticle
	}

	minWords := config.MinArticleWordCount
	if article.IsSocial {
		minWords = config.MinSocialWordCount
	}

	if article.WordCount < minWords {
		return ErrInsufficientContent
	}

	return nil
}
