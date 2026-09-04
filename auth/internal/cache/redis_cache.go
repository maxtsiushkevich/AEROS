package cache

import (
	"auth/internal/config"
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRevokedTokenCache struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedis(cfg *config.Config, logger *slog.Logger) (*RedisRevokedTokenCache, error) {
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

	return &RedisRevokedTokenCache{
		client: client,
		logger: logger,
	}, nil
}

func (c *RedisRevokedTokenCache) Get(ctx context.Context, tkn string) ([]byte, error) {
	value, err := c.client.Get(ctx, "revoked:refresh:"+tkn).Bytes()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *RedisRevokedTokenCache) Set(ctx context.Context, tkn string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, "revoked:refresh:"+tkn, value, ttl).Err()
}

func (c *RedisRevokedTokenCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
