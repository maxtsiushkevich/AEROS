package cache

import (
	"context"
	"time"
	"users/internal/domain/otp"
)

type Cache interface {
	Get(ctx context.Context, key string, purpose otp.OTPPurpose) ([]byte, error)
	Set(ctx context.Context, key string, purpose otp.OTPPurpose, payload []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string, purpose otp.OTPPurpose) error
}
