package cache

import (
	"context"
	"testing"
	"time"

	"jevbeez/api/internal/model"
)

func TestMemoryCache_GetSet(t *testing.T) {
	mem := NewMemoryCache(10, 5*time.Minute)

	res := &model.VerificationResult{
		Verdict:      "VALID",
		ArticleTitle: "Uji coba klaim fakta",
	}

	mem.Set("hash1", res, 5*time.Minute)

	cached := mem.Get("hash1")
	if cached == nil {
		t.Fatalf("expected item to be in cache, got nil")
	}
	if cached.Verdict != "VALID" || cached.ArticleTitle != "Uji coba klaim fakta" {
		t.Errorf("unexpected cached data: %+v", cached)
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	mem := NewMemoryCache(10, 50*time.Millisecond)

	res := &model.VerificationResult{
		Verdict: "HOAX",
	}

	mem.Set("hash_exp", res, 50*time.Millisecond)

	if mem.Get("hash_exp") == nil {
		t.Fatalf("expected item to be present immediately")
	}

	time.Sleep(60 * time.Millisecond)

	if mem.Get("hash_exp") != nil {
		t.Fatalf("expected item to be expired, but got non-nil")
	}
}

func TestMemoryCache_CapacityEviction(t *testing.T) {
	// Max capacity = 3
	mem := NewMemoryCache(3, 5*time.Minute)

	res1 := &model.VerificationResult{ArticleTitle: "1"}
	res2 := &model.VerificationResult{ArticleTitle: "2"}
	res3 := &model.VerificationResult{ArticleTitle: "3"}
	res4 := &model.VerificationResult{ArticleTitle: "4"}

	mem.Set("k1", res1, 5*time.Minute)
	time.Sleep(2 * time.Millisecond)
	mem.Set("k2", res2, 5*time.Minute)
	time.Sleep(2 * time.Millisecond)
	mem.Set("k3", res3, 5*time.Minute)
	time.Sleep(2 * time.Millisecond)

	// Inserting 4th item should evict oldest (k1)
	mem.Set("k4", res4, 5*time.Minute)

	if mem.Get("k1") != nil {
		t.Errorf("expected k1 to be evicted when capacity reached")
	}
	if mem.Get("k4") == nil {
		t.Errorf("expected k4 to be present")
	}
	if mem.Len() > 3 {
		t.Errorf("expected cache size <= 3, got %d", mem.Len())
	}
}

func TestCacheService_L1Shielding(t *testing.T) {
	svc := NewMemoryOnlyService(nil)

	res := &model.VerificationResult{
		Verdict:      "HOAX",
		ArticleTitle: "Judul Berita Palsu",
	}

	ctx := context.Background()
	err := svc.SetByHash(ctx, "abc123hash", res, 5*time.Minute)
	if err != nil {
		t.Fatalf("unexpected error setting cache: %v", err)
	}

	cached, err := svc.GetByHash(ctx, "abc123hash")
	if err != nil {
		t.Fatalf("unexpected error getting cache: %v", err)
	}
	if cached == nil {
		t.Fatalf("expected cached result, got nil")
	}
	if !cached.IsCached {
		t.Errorf("expected IsCached to be true")
	}
	if cached.ArticleTitle != "Judul Berita Palsu" {
		t.Errorf("expected 'Judul Berita Palsu', got '%s'", cached.ArticleTitle)
	}
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	mem := NewMemoryCache(500, 5*time.Minute)
	res := &model.VerificationResult{Verdict: "VALID", ArticleTitle: "Judul"}
	mem.Set("test_key", res, 5*time.Minute)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = mem.Get("test_key")
	}
}

func BenchmarkMemoryCache_Set(b *testing.B) {
	mem := NewMemoryCache(500, 5*time.Minute)
	res := &model.VerificationResult{Verdict: "VALID", ArticleTitle: "Judul"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		mem.Set("test_key", res, 5*time.Minute)
	}
}

func BenchmarkCacheService_GetByHash_L1Hit(b *testing.B) {
	svc := NewMemoryOnlyService(nil)
	res := &model.VerificationResult{Verdict: "VALID", ArticleTitle: "Judul"}
	ctx := context.Background()
	_ = svc.SetByHash(ctx, "hash_benchmark", res, 5*time.Minute)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetByHash(ctx, "hash_benchmark")
	}
}
