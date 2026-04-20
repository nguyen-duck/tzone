package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheService(client *redis.Client, ttl time.Duration) *CacheService {
	if ttl <= 0 {
		ttl = 2 * time.Minute
	}
	return &CacheService{client: client, ttl: ttl}
}

func (s *CacheService) enabled() bool {
	return s != nil && s.client != nil
}

func (s *CacheService) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	if !s.enabled() {
		return false, nil
	}

	value, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return false, err
	}

	return true, nil
}

func (s *CacheService) SetJSON(ctx context.Context, key string, value any) error {
	if !s.enabled() {
		return nil
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, encoded, s.ttl).Err()
}

func (s *CacheService) DeleteByPrefixes(ctx context.Context, prefixes ...string) error {
	if !s.enabled() || len(prefixes) == 0 {
		return nil
	}

	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}

		pattern := fmt.Sprintf("%s*", prefix)
		iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
				return err
			}
		}
		if err := iter.Err(); err != nil {
			return err
		}
	}

	return nil
}
