package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func Connect(redisURL string) (*redis.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	url := strings.TrimSpace(redisURL)
	if url == "" {
		cancel()
		return nil, nil, nil, fmt.Errorf("redis url is empty")
	}

	var client *redis.Client
	if strings.HasPrefix(url, "redis://") || strings.HasPrefix(url, "rediss://") {
		opt, err := redis.ParseURL(url)
		if err != nil {
			cancel()
			return nil, nil, nil, fmt.Errorf("failed to parse redis url: %w", err)
		}
		client = redis.NewClient(opt)
	} else {
		client = redis.NewClient(&redis.Options{Addr: url})
	}

	if err := client.Ping(ctx).Err(); err != nil {
		cancel()
		_ = client.Close()
		return nil, nil, nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, ctx, cancel, nil
}

func Close(client *redis.Client) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
