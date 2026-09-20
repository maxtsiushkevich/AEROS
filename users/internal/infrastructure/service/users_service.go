package service

import (
	"time"
	"users/internal/domain/models"
	"users/internal/domain/storage"

	"github.com/google/uuid"
)

type UsersService struct {
	storage storage.UsersStorage
}

func NewUsersService(storage storage.UsersStorage) *UsersService {
	return &UsersService{
		storage: storage,
	}
}

func (s *UsersService) RegisterUser(name string, email string, birthday time.Time) (*models.User, error) {
	return nil, nil
}

func (s *UsersService) ActivateUser(id uuid.UUID) error {
	return nil
}

func (s *UsersService) Profile(id uuid.UUID) (*models.User, error) {
	return nil, nil
}
