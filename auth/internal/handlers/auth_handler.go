package handlers

import (
	"auth/internal/errors"
	"auth/internal/service"
	"auth/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	storage   storage.AuthStorage
	service   *service.AuthService
	validator *validator.Validate
	logger    *slog.Logger
}

func NewAuthHandler(storage storage.AuthStorage, logger *slog.Logger) (*AuthHandler, error) {
	if err := storage.Open(); err != nil {
		logger.Error("Failed to open storage", "err", err)
		return nil, errors.NewAuthError("STORAGE_INIT_FAILED", "failed to initialize storage")
	}

	return &AuthHandler{
		storage:   storage,
		service:   service.CreateAuthService(storage),
		validator: validator.New(),
		logger:    logger,
	}, nil
}

func (h *AuthHandler) HandleRefreshTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
func (h *AuthHandler) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandler) HandleChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
