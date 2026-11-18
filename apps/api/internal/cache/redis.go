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

	return &CacheService{client: client}, nil
}

func (s *CacheService) Get(ctx context.Context, normalizedURL string) (*model.VerificationResult, error) {
	key := s.formatURLKey(normalizedURL)
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
	return &res, nil
}

func (s *CacheService) Set(ctx context.Context, normalizedURL string, result *model.VerificationResult, ttl time.Duration) error {
	key := s.formatURLKey(normalizedURL)
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *CacheService) CheckAndSetCooldown(ctx context.Context, ip string, duration time.Duration) (bool, time.Duration, error) {
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

func (s *CacheService) Close() error {
	return s.client.Close()
}

func (s *CacheService) formatURLKey(normalizedURL string) string {
	h := sha256.Sum256([]byte(normalizedURL))
	return fmt.Sprintf("jevbeez:cache:%s", hex.EncodeToString(h[:]))
}
