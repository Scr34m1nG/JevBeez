package cache

import (
	"sync"
	"time"

	"jevbeez/api/internal/model"
)

type cacheEntry struct {
	value     *model.VerificationResult
	expiresAt time.Time
	updatedAt time.Time
}

type MemoryCache struct {
	mu          sync.RWMutex
	items       map[string]cacheEntry
	maxCapacity int
	defaultTTL  time.Duration
}

func NewMemoryCache(maxCapacity int, defaultTTL time.Duration) *MemoryCache {
	if maxCapacity <= 0 {
		maxCapacity = 500
	}
	if defaultTTL <= 0 {
		defaultTTL = 3 * time.Minute
	}
	return &MemoryCache{
		items:       make(map[string]cacheEntry, maxCapacity),
		maxCapacity: maxCapacity,
		defaultTTL:  defaultTTL,
	}
}

func (m *MemoryCache) Get(key string) *model.VerificationResult {
	if m == nil {
		return nil
	}

	m.mu.RLock()
	entry, exists := m.items[key]
	m.mu.RUnlock()

	if !exists {
		return nil
	}

	if time.Now().After(entry.expiresAt) {
		m.mu.Lock()
		delete(m.items, key)
		m.mu.Unlock()
		return nil
	}

	res := *entry.value
	res.IsCached = true
	return &res
}

func (m *MemoryCache) Set(key string, val *model.VerificationResult, ttl time.Duration) {
	if m == nil || val == nil {
		return
	}

	if ttl <= 0 {
		ttl = m.defaultTTL
	}

	now := time.Now()
	entry := cacheEntry{
		value:     val,
		expiresAt: now.Add(ttl),
		updatedAt: now,
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// If capacity reached and key is new, evict
	if len(m.items) >= m.maxCapacity {
		if _, exists := m.items[key]; !exists {
			m.evictOne(now)
		}
	}

	m.items[key] = entry
}

func (m *MemoryCache) evictOne(now time.Time) {
	// First pass: remove any expired key
	for k, v := range m.items {
		if now.After(v.expiresAt) {
			delete(m.items, k)
			return
		}
	}

	// Second pass: evict oldest updated key
	var oldestKey string
	var oldestTime time.Time
	first := true

	for k, v := range m.items {
		if first || v.updatedAt.Before(oldestTime) {
			oldestKey = k
			oldestTime = v.updatedAt
			first = false
		}
	}

	if oldestKey != "" {
		delete(m.items, oldestKey)
	}
}

func (m *MemoryCache) Len() int {
	if m == nil {
		return 0
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.items)
}
