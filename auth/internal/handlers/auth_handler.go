package handlers

import (
	"auth/internal/service"
	"auth/internal/storage"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	storage   storage.AuthStorage
	service   *service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(storage storage.AuthStorage) (*AuthHandler, error) {
	return &AuthHandler{
		storage:   storage,
		service:   service.CreateAuthService(storage),
		validator: validator.New(),
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
