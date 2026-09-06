package middleware

import (
	"context"
	"fmt"
	"net/http"
	"pkg/auth"
	"pkg/httperr"
)

func methodToAction(method string) string {
	switch method {
	case http.MethodGet:
		return "read"
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		return "write"
	case http.MethodDelete:
		return "delete"
	default:
		return "read"
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

type AuthorizationChecker interface {
	IsAuthenticated(sub string, obj string, act string) (bool, error)
}

func AuthMiddleware(authorizer AuthorizationChecker) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := auth.TokenFromRequest(r)
			if tokenString == "" {
				httperr.Write(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			claims, err := auth.ParseAccessToken(tokenString)
			if err != nil {
				fmt.Println("JWT parse error:", err)
				httperr.Write(w, http.StatusUnauthorized, "invalid token")
				return
			}

			r = r.WithContext(WithClaims(r.Context(), claims))

			if authorizer == nil {
				httperr.Write(w, http.StatusInternalServerError, "access checker not configured")
				return
			}

			act := methodToAction(r.Method)
			ok, err := authorizer.IsAuthenticated(claims.Id.String(), r.URL.Path, act)
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
