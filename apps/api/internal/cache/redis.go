package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/model"
)

type CacheService struct {
	client *redis.Client
	l1     *MemoryCache
}

func NewRedisService(addr, password string, db int) (*CacheService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), config.DialerTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", addr, err)
	}

	return &CacheService{
		client: client,
		l1:     NewMemoryCache(500, 3*time.Minute),
	}, nil
}

func NewMemoryOnlyService(l1 *MemoryCache) *CacheService {
	if l1 == nil {
		l1 = NewMemoryCache(500, 3*time.Minute)
	}
	return &CacheService{
		l1: l1,
	}
}

func (s *CacheService) Get(ctx context.Context, normalizedURL string) (*model.VerificationResult, error) {
	h := sha256.Sum256([]byte(normalizedURL))
	return s.GetByHash(ctx, hex.EncodeToString(h[:]))
}

func (s *CacheService) Set(ctx context.Context, normalizedURL string, result *model.VerificationResult, ttl time.Duration) error {
	h := sha256.Sum256([]byte(normalizedURL))
	return s.SetByHash(ctx, hex.EncodeToString(h[:]), result, ttl)
}

func (s *CacheService) GetByHash(ctx context.Context, hash string) (*model.VerificationResult, error) {
	if s == nil {
		return nil, nil
	}

	// 1. Check L1 Memory Cache first (< 1us)
	if s.l1 != nil {
		if val := s.l1.Get("cache:" + hash); val != nil {
			return val, nil
		}
	}

	if s.client == nil {
		return nil, nil
	}

	// 2. Check L2 Redis Cache
	key := fmt.Sprintf("jevbeez:cache:%s", hash)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var res model.VerificationResult
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, err
	}
	res.IsCached = true

	// Backfill L1 Memory cache
	if s.l1 != nil {
		s.l1.Set("cache:"+hash, &res, 3*time.Minute)
	}

	return &res, nil
}

func (s *CacheService) SetByHash(ctx context.Context, hash string, result *model.VerificationResult, ttl time.Duration) error {
	if s == nil || result == nil {
		return nil
	}

	// Save to L1 Memory cache
	if s.l1 != nil {
		memTTL := ttl
		if memTTL > 5*time.Minute {
			memTTL = 5 * time.Minute
		}
		s.l1.Set("cache:"+hash, result, memTTL)
	}

	if s.client == nil {
		return nil
	}

	key := fmt.Sprintf("jevbeez:cache:%s", hash)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *CacheService) GetShare(ctx context.Context, hash string) (*model.VerificationResult, error) {
	if s == nil {
		return nil, nil
	}

	if s.l1 != nil {
		if val := s.l1.Get("share:" + hash); val != nil {
			return val, nil
		}
	}

	if s.client == nil {
		return nil, nil
	}

	key := fmt.Sprintf("jevbeez:share:%s", hash)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var res model.VerificationResult
	if err := json.Unmarshal([]byte(val), &res); err != nil {
		return nil, err
	}
	res.IsCached = true

	if s.l1 != nil {
		s.l1.Set("share:"+hash, &res, 3*time.Minute)
	}

	return &res, nil
}

func (s *CacheService) SetShare(ctx context.Context, hash string, result *model.VerificationResult) error {
	if s == nil || result == nil {
		return nil
	}

	if s.l1 != nil {
		s.l1.Set("share:"+hash, result, 5*time.Minute)
	}

	if s.client == nil {
		return nil
	}

	key := fmt.Sprintf("jevbeez:share:%s", hash)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return s.client.SetNX(ctx, key, data, config.ShareCacheTTL).Err()
}

func (s *CacheService) CheckAndSetCooldown(ctx context.Context, ip string, duration time.Duration) (bool, time.Duration, error) {
	if s == nil || s.client == nil {
		return true, 0, nil
	}
	key := fmt.Sprintf("jevbeez:cooldown:%s", ip)
	ok, err := s.client.SetNX(ctx, key, "1", duration).Result()
	if err != nil {
		return false, 0, err
	}

	if !ok {
		ttl, err := s.client.TTL(ctx, key).Result()
		if err != nil {
			return false, duration, nil
		}
		return false, ttl, nil
	}

	return true, 0, nil
}

func (s *CacheService) BlockIP(ctx context.Context, ip, reason, payload string, duration time.Duration) error {
	if s == nil || s.client == nil || ip == "" {
		return nil
	}

	if len(payload) > 100 {
		payload = payload[:100] + "..."
	}

	now := time.Now()
	expiresAt := now.Add(duration)

	info := model.BlockedIPInfo{
		IP:        ip,
		Reason:    reason,
		Payload:   payload,
		BlockedAt: now.Format(time.RFC3339),
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("jevbeez:blocked:%s", ip)
	return s.client.Set(ctx, key, data, duration).Err()
}

func (s *CacheService) IsIPBlocked(ctx context.Context, ip string) (bool, *model.BlockedIPInfo, error) {
	if s == nil || s.client == nil || ip == "" {
		return false, nil, nil
	}

	key := fmt.Sprintf("jevbeez:blocked:%s", ip)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil, nil
		}
		return false, nil, err
	}

	var info model.BlockedIPInfo
	if err := json.Unmarshal([]byte(val), &info); err != nil {
		return true, &model.BlockedIPInfo{IP: ip, Reason: "BLOCKED"}, nil
	}

	return true, &info, nil
}

func (s *CacheService) Close() error {
	return s.client.Close()
}
