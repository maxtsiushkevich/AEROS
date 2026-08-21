package handlers

import "net/http"

type AuthHandler struct {
}

func NewAuthHandler() (*AuthHandler, error) {
	return nil, nil
}

func (h *AuthHandler) HandleCreateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
