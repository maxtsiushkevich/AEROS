package repository

import (
	"context"
	"time"
	"users/internal/domain/models"
	"users/internal/domain/outbox"

	"github.com/google/uuid"
)

type UserStorage interface {
	Create(ctx context.Context, u *models.User) error
	Read(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	UpdateEmailWithOutbox(ctx context.Context, id uuid.UUID, email string, event outbox.OutboxEvent) error
	ClaimOutboxBatch(ctx context.Context, limit int, lease time.Duration) ([]outbox.OutboxEvent, error)
	MarkOutboxDelivered(ctx context.Context, id uuid.UUID) error
	RetryOutboxEvent(ctx context.Context, id uuid.UUID, nextAttempt time.Time, lastError string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
