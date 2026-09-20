package service

import (
	"time"
	"users/internal/domain/models"

	"github.com/google/uuid"
)

type UsersService interface {
	RegisterUser(name string, email string, birthday time.Time) (*models.User, error)
	ActivateUser(id uuid.UUID) error
	Profile(id uuid.UUID) (*models.User, error)
}
