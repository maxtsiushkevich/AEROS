package cache

import (
	"auth/internal/config"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedis(cfg *config.Config, logger *slog.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Database.Password,
		DB:       cfg.Redis.Db,
		Protocol: cfg.Redis.Protocol,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		logger: logger,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
