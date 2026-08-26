package service

import (
	"auth/internal/storage"

	"auth/internal/cache"
)

type AuthService struct {
	storage            storage.AuthStorage
	revokedTokensCache cache.Cache
}

func CreateAuthService(storage storage.AuthStorage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (s *AuthService) Login(email *string, password *string) (access string, refresh string, err error) {
	return "", "", nil
}

func (s *AuthService) Logout(access string, refresh string) {
	// Add tokens to blacklist (Redis or Memcached)
}

func (s *AuthService) Refresh(oldRefreshToken string) (access string, refresh string, err error) {
	return "", "", nil
}
