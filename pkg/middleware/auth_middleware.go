package middleware

import (
	"context"
	"fmt"
	"net/http"
	"pkg/auth"
	"pkg/helpers"
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

func withClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, helpers.ClaimsContextKey, claims)
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

			r = r.WithContext(withClaims(r.Context(), claims))

			if authorizer == nil {
				httperr.Write(w, http.StatusInternalServerError, "access checker not configured")
				return
			}

			act := methodToAction(r.Method)
			obj := CasbinObjFromRequest(r)

			ok, err := authorizer.IsAuthenticated(claims.Id.String(), obj, act)

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
