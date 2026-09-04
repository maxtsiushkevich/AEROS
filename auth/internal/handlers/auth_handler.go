package handlers

import (
	"auth/internal/cache"
	"auth/internal/dto"
	"auth/internal/service"
	"auth/internal/storage"
	authErrors "auth/pkg/errors"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/maxtsiushkevich/AEROS/pkg/httperr"
)

type AuthHandler struct {
	storage  storage.AuthStorage
	cache    cache.RevokedTokenCache
	service  *service.AuthService
	validate *validator.Validate
	logger   *slog.Logger
}

func NewAuthHandler(storage storage.AuthStorage, logger *slog.Logger, cache cache.RevokedTokenCache) (*AuthHandler, error) {
	return &AuthHandler{
		storage:  storage,
		cache:    cache,
		service:  service.CreateAuthService(storage, cache),
		validate: validator.New(),
		logger:   logger,
	}, nil
}

func (h *AuthHandler) HandleRefreshTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("refresh_token")
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Refresh token missed")
			return
		}

		access, refresh, err := h.service.Refresh(r.Context(), cookie.Value)
		if err != nil {
			status := http.StatusUnauthorized
			if authErr, ok := err.(*authErrors.AuthError); ok && authErr.Code == "CACHE_ERROR" {
				status = http.StatusInternalServerError
			}
			httperr.Write(w, status, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refresh,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"access_token": access}); err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to write response")
		}
	}
}

func (h *AuthHandler) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		var authRequest dto.AuthRequest

		if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		authRequest.Normalize()

		err := h.validate.Struct(authRequest)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		access, refresh, err := h.service.Login(ctx, &authRequest.Email, &authRequest.Password)
		if err != nil {
			httperr.Write(w, http.StatusUnauthorized, err.Error())
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    *refresh,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"access_token": *access})
	}
}

func (h *AuthHandler) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Logout complete"})
			return
		}

		err = h.service.Logout(r.Context(), refreshToken.Value)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to logout")
			h.logger.Error("Error", "err", err)
		}

		http.SetCookie(w, &http.Cookie{
			Name: "refresh_token", Value: "", Path: "/",
			HttpOnly: true, MaxAge: -1, SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Logout complete"})
	}
}

func (h *AuthHandler) HandleChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandler) HandleSecure() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
