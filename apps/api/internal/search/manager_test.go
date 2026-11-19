package search

import (
	"context"
	"errors"
	"testing"

	"jevbeez/api/internal/model"
)

func TestSearchManager_BuildQuery(t *testing.T) {
	sm := &SearchManager{}

	t.Run("uses article title for regular news", func(t *testing.T) {
		art := &model.Article{
			Title:    "BMKG Klarifikasi Isu Gempa Megathrust - Detik News",
			Content:  "Ini adalah isi berita panjang lebih dari 80 kata...",
			IsSocial: false,
		}
		query := sm.buildQuery(art)
		expected := "BMKG Klarifikasi Isu Gempa Megathrust"
		if query != expected {
			t.Errorf("expected query '%s', got '%s'", expected, query)
		}
	})

	t.Run("uses caption text for social media post", func(t *testing.T) {
		art := &model.Article{
			Title:    "Postingan X (@berita)",
			Content:  "Waspada gempa susulan di pesisir selatan laut jawa hari ini. Mohon sebarkan informasi ini seluas-luasnya.",
			IsSocial: true,
			Platform: "twitter",
		}
		query := sm.buildQuery(art)
		expected := "Waspada gempa susulan di pesisir selatan laut jawa hari ini"
		if query != expected {
			t.Errorf("expected query '%s', got '%s'", expected, query)
		}
	})

	t.Run("falls back to title if social content is too short", func(t *testing.T) {
		art := &model.Article{
			Title:    "Gempa M5.0",
			Content:  "Test",
			IsSocial: true,
			Platform: "twitter",
		}
		query := sm.buildQuery(art)
		expected := "Gempa M5.0"
		if query != expected {
			t.Errorf("expected query '%s', got '%s'", expected, query)
		}
	})

	t.Run("falls back to content when article title is empty", func(t *testing.T) {
		art := &model.Article{
			Title:    "",
			Content:  "Pemerintah mengumumkan kenaikan insentif tenaga kesehatan mulai bulan depan. Kebijakan ini diharapkan...",
			IsSocial: false,
		}
		query := sm.buildQuery(art)
		expected := "Pemerintah mengumumkan kenaikan insentif tenaga kesehatan mulai bulan depan"
		if query != expected {
			t.Errorf("expected query '%s', got '%s'", expected, query)
		}
	})

	t.Run("truncates query safely with utf8 multibyte characters", func(t *testing.T) {
		// 130 characters containing 3-byte and 4-byte runes
		titleWithRunes := "🚨 Waspada gempa susulan berkekuatan magnitudo besar di wilayah pantai selatan jawa timur hari ini menurut bmkg dan peneliti geologi nasional"
		art := &model.Article{
			Title:    titleWithRunes,
			Content:  "Isi berita...",
			IsSocial: false,
		}
		query := sm.buildQuery(art)
		if len([]rune(query)) > 120 {
			t.Errorf("expected rune length <= 120, got %d", len([]rune(query)))
		}
		for i, r := range query {
			if r == '\ufffd' {
				t.Fatalf("found invalid UTF-8 replacement char at byte %d", i)
			}
		}
	})
}

type mockProvider struct {
	name    string
	results []model.Evidence
	err     error
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) Search(ctx context.Context, query string) ([]model.Evidence, error) {
	return m.results, m.err
}

func TestSearchManager_SearchGrounding_Fallback(t *testing.T) {
	ctx := context.Background()
	art := &model.Article{
		Title: "Judul Berita Valid",
	}

	t.Run("returns empty slice without error if one provider succeeds with 0 results", func(t *testing.T) {
		sm := &SearchManager{
			providers: []SearchProvider{
				&mockProvider{name: "failing", err: errors.New("rate limited")},
				&mockProvider{name: "empty", results: []model.Evidence{}, err: nil},
			},
		}
		res, err := sm.SearchGrounding(ctx, art)
		if err != nil {
			t.Fatalf("expected nil error, got: %v", err)
		}
		if len(res) != 0 {
			t.Errorf("expected 0 results, got %d", len(res))
		}
	})

	t.Run("returns error when all providers fail", func(t *testing.T) {
		sm := &SearchManager{
			providers: []SearchProvider{
				&mockProvider{name: "failing1", err: errors.New("err 1")},
				&mockProvider{name: "failing2", err: errors.New("err 2")},
			},
		}
		_, err := sm.SearchGrounding(ctx, art)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
