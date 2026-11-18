package search

import (
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
}
