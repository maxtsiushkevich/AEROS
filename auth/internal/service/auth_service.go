package service

import (
	"auth/internal/cache"
	"auth/internal/models"
	"auth/internal/storage"
	"auth/pkg/auth"
	"auth/pkg/errors"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	storage            storage.AuthStorage
	revokedTokensCache cache.RevokedTokenCache
}

func CreateAuthService(storage storage.AuthStorage, cache cache.RevokedTokenCache) *AuthService {
	return &AuthService{
		storage:            storage,
		revokedTokensCache: cache,
	}
}

func (s *AuthService) Login(ctx context.Context, email *string, password *string) (access *string, refresh *string, err error) {
	user, err := s.storage.Read(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	isValidPassword := auth.CheckPassword(*password, user.HashedPassword)
	if !isValidPassword {
		return nil, nil, errors.InvalidPasswordError
	}

	accessToken, err := auth.GenerateAccessToken(&user.ID, &user.Email, &user.Version)
	if err != nil {
		return nil, nil, errors.NewCreateTokenError("failed to create access token")
	}

	refreshToken, err := auth.GenerateRefreshToken(&user.ID, &user.Email, &user.Version)
	if err != nil {
		return nil, nil, errors.NewCreateTokenError("failed to create refresh token")
	}

	access = &accessToken
	refresh = &refreshToken

	return access, refresh, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return auth.JwtRefreshKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || claims.ExpiresAt == nil {
		return errors.NewCreateTokenError("failed to parse refresh token")
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}
	return s.revokedTokensCache.Set(ctx, refreshToken, []byte("revoked"), ttl)
}

func (s *AuthService) refreshTokens(ctx context.Context, userID *uuid.UUID, email *string, version *uint32) (access *string, refresh *string, err error) {
	newAccessToken, err := auth.GenerateAccessToken(userID, email, version)
	if err != nil {
		return nil, nil, errors.NewCreateTokenError("failed to create access token")
	}

	newRefreshToken, err := auth.GenerateRefreshToken(userID, email, version)
	if err != nil {
		return nil, nil, errors.NewCreateTokenError("failed to create refresh token")
	}

	access = &newAccessToken
	refresh = &newRefreshToken
	return access, refresh, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (access string, refresh string, err error) {
	_, err = s.revokedTokensCache.Get(ctx, oldRefreshToken)
	if err == nil {
		return "", "", errors.NewCreateTokenError("refresh token was revoked")
	}
	if err != redis.Nil {
		return "", "", errors.NewAuthError("CACHE_ERROR", fmt.Sprintf("failed to check refresh token: %v", err))
	}

	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(oldRefreshToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return auth.JwtRefreshKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !token.Valid || claims.Type != "refresh" {
		return "", "", errors.NewCreateTokenError("invalid refresh token")
	}

	user, err := s.storage.ReadByID(ctx, claims.Id)
	if err != nil || user.Version != claims.Version {
		return "", "", errors.NewCreateTokenError("invalid refresh token")
	}

	newAccess, newRefresh, err := s.refreshTokens(ctx, &user.ID, &user.Email, &user.Version)
	if err != nil {
		return "", "", err
	}

	return *newAccess, *newRefresh, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID *uuid.UUID, oldPassword *string, newPassword *string) (*models.UserAuth, error) {
	user, err := s.storage.ReadByID(ctx, *userID)
	if err != nil {
		return nil, err
	}

	// Check if the user password correct
	if !auth.CheckPassword(*oldPassword, user.HashedPassword) {
		return nil, errors.InvalidPasswordError
	}

	// Check if the new password is the same as the old password
	if auth.CheckPassword(*newPassword, user.HashedPassword) {
		return nil, errors.SamePasswordError
	}

	hashedPassword, err := auth.HashPassword(*newPassword)
	if err != nil {
		return nil, err
	}

	upd := models.UserAuthUpdate{
		ID:                *userID,
		NewHashedPassword: hashedPassword,
	}

	user, err = s.storage.Update(ctx, &upd)
	if err != nil {
		return nil, err
	}

	return user, nil
}
