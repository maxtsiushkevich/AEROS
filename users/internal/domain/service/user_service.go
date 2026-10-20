package service

import (
	"context"
	"time"
	"users/internal/domain/models"
	"users/internal/domain/otp"

	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, name string, email string, birthday time.Time, password string) (
		access *string,
		refresh *string,
		err error)
	ActivateUserByAdmin(ctx context.Context, id uuid.UUID) error
	ActivateUserByCode(ctx context.Context, id uuid.UUID, code string) error
	Profile(ctx context.Context, id uuid.UUID) (*models.User, error)
	RequestCode(ctx context.Context, id uuid.UUID, email string, purpose otp.OTPPurpose) error
	// ChangeEmail(ctx context.Context, id uuid.UUID, newEmail string) error
	// DeleteUser(ctx context.Context, id uuid.UUID) error
}
