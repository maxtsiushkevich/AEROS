package middleware

import (
	"auth/rbac"
	"net/http"
	"strings"
)

func tokenFromRequest(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

func AuthMiddleware(rbacService rbac.AuthorizationService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
		}
	}
}
