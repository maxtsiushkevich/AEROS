package service

import "auth/internal/storage"

type AuthService struct {
	storage storage.AuthStorage
}

func CreateAuthService(storage storage.AuthStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}
