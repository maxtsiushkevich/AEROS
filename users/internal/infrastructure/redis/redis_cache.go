package redis

import (
	"context"
	"log/slog"
	"time"
	"users/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisOTPCache struct {
	client *redis.Client
	logger *slog.Logger
	topic  string
}

func NewRedisOTPCache(cfg *config.Config, logger *slog.Logger) (*RedisOTPCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
		Protocol: cfg.Redis.Protocol,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisOTPCache{
		client: client,
		logger: logger,
		topic:  cfg.Redis.Topic,
	}, nil
}

func (c *RedisOTPCache) Shutdown() error {
	err := c.client.Close()
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	c.logger.Info("Redis cache connection closed")
	return nil
}

func (c *RedisOTPCache) Get(ctx context.Context, email string) ([]byte, error) {
	value, err := c.client.Get(ctx, c.topic+":"+email).Bytes()
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (c *RedisOTPCache) Set(ctx context.Context, email string, code []byte, ttl time.Duration) error {
	return c.client.Set(ctx, c.topic+":"+email, code, ttl).Err()
}

func (c *RedisOTPCache) Delete(ctx context.Context, email string) error {
	return c.client.Del(ctx, c.topic+":"+email).Err()
}
