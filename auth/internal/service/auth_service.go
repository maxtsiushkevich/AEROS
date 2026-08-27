package service

import (
	"auth/internal/auth"
	"auth/internal/cache"
	"auth/internal/errors"
	"auth/internal/storage"
	"context"

	"github.com/google/uuid"
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

func (s *AuthService) Login(ctx context.Context, email *string, password *string) (access *string, refresh *string, err error) {
	user, err := s.storage.Read(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	isValidPassword := auth.CheckPassword(*password, user.HashedPassword)
	if !isValidPassword {
		return nil, nil, errors.InvalidPasswordError("incorrect password")
	}

	accessToken, err := auth.GenerateAccessToken(&user.ID, &user.Email, &user.Version)
	if err != nil {
		return nil, nil, errors.CreateTokenError("failed to create access token")
	}

	refreshToken, err := auth.GenerateRefreshToken(&user.ID, &user.Email, &user.Version)
	if err != nil {
		return nil, nil, errors.CreateTokenError("failed to create refresh token")
	}

	access = &accessToken
	refresh = &refreshToken

	return access, refresh, nil
}

func (s *AuthService) Logout(ctx context.Context, access string, refresh string) {
	// Add tokens to blacklist (Redis or Memcached)
}

func (s *AuthService) RefreshTokens(ctx context.Context, userID *uuid.UUID, email *string, version *uint32) (access *string, refresh *string, err error) {
	newAccessToken, err := auth.GenerateAccessToken(userID, email, version)
	if err != nil {
		return nil, nil, errors.CreateTokenError("failed to create access token")
	}

	newRefreshToken, err := auth.GenerateRefreshToken(userID, email, version)
	if err != nil {
		return nil, nil, errors.CreateTokenError("failed to create refresh token")
	}

	access = &newAccessToken
	refresh = &newRefreshToken
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (access string, refresh string, err error) {
	return "", "", nil
}
