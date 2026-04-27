package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLCache struct {
	rdb *redis.Client
}

func NewURLCache(rdb *redis.Client) *URLCache {
	return &URLCache{rdb: rdb}
}

func (c *URLCache) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *URLCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}
func (c *URLCache) Del(ctx context.Context, key string) (int64, error) {
	return c.rdb.Del(ctx, key).Result()
}
