package service

import (
	"context"
	"errors"
	domainerr "users/internal/domain/errors"
	"users/internal/domain/repository"
)

type UserService struct {
	storage repository.UserStorage
}

func NewUserService(storage repository.UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (s *UserService) CheckEmailAvailability(ctx context.Context, email string) error {
	user, err := s.storage.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, domainerr.ErrUserNotFound) {
		return err
	}
	if user != nil {
		return domainerr.ErrUserAlreadyExists
	}

	return nil
}
