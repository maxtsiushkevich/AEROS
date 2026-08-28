package service

import (
	"auth/internal/auth"
	"auth/internal/cache"
	"auth/internal/errors"
	"auth/internal/storage"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	storage            storage.AuthStorage
	revokedTokensCache cache.Cache
}

func CreateAuthService(storage storage.AuthStorage, cache cache.Cache) *AuthService {
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

func (s *AuthService) Logout(ctx context.Context, refresh string) error {
	claims := &auth.Claims{}
	_, err := jwt.ParseWithClaims(refresh, claims,
		func(token *jwt.Token) (interface{}, error) {
			return auth.JwtRefreshKey, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || claims.ExpiresAt == nil {
		return errors.CreateTokenError("failed to parse refresh token")
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return s.revokedTokensCache.Set(ctx, "revoked:refresh:"+refresh, []byte("1"), ttl)
}

func (s *AuthService) refreshTokens(ctx context.Context, userID *uuid.UUID, email *string, version *uint32) (access *string, refresh *string, err error) {
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
	_, err = s.revokedTokensCache.Get(ctx, "revoked:refresh:"+oldRefreshToken)
	if err == nil {
		return "", "", errors.CreateTokenError("refresh token was revoked")
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
		return "", "", errors.CreateTokenError("invalid refresh token")
	}

	user, err := s.storage.ReadByID(ctx, claims.Id)
	if err != nil || user.Version != claims.Version {
		return "", "", errors.CreateTokenError("invalid refresh token")
	}

	newAccess, newRefresh, err := s.refreshTokens(ctx, &user.ID, &user.Email, &user.Version)
	if err != nil {
		return "", "", err
	}

	return *newAccess, *newRefresh, nil
}
