package service

import (
	"context"
	"time"
	"users/internal/domain/models"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, name string, email string, birthday time.Time, password string) (
		access *string,
		refresh *string,
		err error)
	ActivateUser(ctx context.Context, id uuid.UUID) error
	Profile(ctx context.Context, id uuid.UUID) (*models.User, error)
	ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
