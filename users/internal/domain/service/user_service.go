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
	AdminActivateUser(ctx context.Context, id uuid.UUID) error
	ActivateUserByCode(ctx context.Context, id uuid.UUID, code string) error
	Profile(ctx context.Context, id uuid.UUID) (*models.User, error)
	RequestActivationCode(ctx context.Context, id uuid.UUID, email string) error
	ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
