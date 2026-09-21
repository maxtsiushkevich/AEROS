package repository

import (
	"context"
	"users/internal/domain/models"

	"github.com/google/uuid"
)

type UserStorage interface {
	Create(ctx context.Context, u *models.User) error
	Read(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
