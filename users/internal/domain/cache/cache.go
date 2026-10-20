package cache

import (
	"context"
	"time"
	"users/internal/domain/otp"
)

type Cache interface {
	Get(ctx context.Context, email string, purpose otp.OTPPurpose) ([]byte, error)
	Set(ctx context.Context, email string, purpose otp.OTPPurpose, code []byte, ttl time.Duration) error
	Delete(ctx context.Context, email string) error
}
