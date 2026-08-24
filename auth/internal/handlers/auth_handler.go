package handlers

import (
	"auth/internal/service"
	"auth/internal/storage"
	"net/http"
)

type AuthHandler struct {
	storage storage.AuthStorage
	service *service.AuthService
}

func NewAuthHandler(storage storage.AuthStorage) (*AuthHandler, error) {
	return &AuthHandler{
		storage: storage,
		service: service.CreateAuthService(storage),
	}, nil
}

func (h *AuthHandler) HandleCreateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
