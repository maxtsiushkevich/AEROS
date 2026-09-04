package middleware

import (
	"auth/internal/storage"
	"auth/pkg/auth"
	"auth/rbac"
	"context"
	"fmt"
	"net/http"

	"github.com/maxtsiushkevich/AEROS/pkg/httperr"
)

func methodToAction(method string) rbac.Action {
	switch method {
	case http.MethodGet:
		return rbac.Read
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return rbac.Write
	case http.MethodDelete:
		return rbac.Delete
	default:
		return rbac.Read
	}
}

type contextKey string

const claimsContextKey contextKey = "auth_claims"

func WithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*auth.Claims)
	return claims, ok
}

// Usage ClaimsFromContext
// claims, ok := middleware.ClaimsFromContext(r.Context())
// if !ok {
// 	httperr.Write(w, http.StatusUnauthorized, "missing claims")
// 	return
// }
// userID := claims.Id
// email := claims.Email
// version := claims.Version

func AuthMiddleware(rbacService rbac.AuthorizationService, authStorage storage.AuthStorage) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := auth.TokenFromRequest(r)
			if tokenString == "" {
				httperr.Write(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			// there is check if token in blacklist - if yes, return 401

			claims, err := auth.ParseAccessToken(tokenString)
			if err != nil {
				fmt.Println("JWT parse error:", err)
				httperr.Write(w, http.StatusUnauthorized, "invalid token")
				return
			}

			user, err := authStorage.ReadByID(r.Context(), claims.Id)
			if err != nil || user.Version != claims.Version {
				httperr.Write(w, http.StatusUnauthorized, "invalid token")
				return
			}

			// write claims to context for further use in handlers
			r = r.WithContext(WithClaims(r.Context(), claims))

			act := methodToAction(r.Method)
			ok, err := rbacService.IsAuthenticated(claims.Id.String(), r.URL.Path, act)
			if err != nil {
				httperr.Write(w, http.StatusInternalServerError, "access check error")
				return
			}

			if !ok {
				httperr.Write(w, http.StatusForbidden, "access forbidden")
				return
			}

			next(w, r)
		}
	}
}
