package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, email string) ([]byte, error)
	Set(ctx context.Context, email string, code []byte, ttl time.Duration) error
	Delete(ctx context.Context, email string) error
}
