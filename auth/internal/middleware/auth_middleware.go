package middleware

import (
	"auth/internal/auth"
	"auth/rbac"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func tokenFromRequest(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}

func parseAccessToken(tokenString string) (*auth.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.JwtAccessKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	if claims.Type != "access" {
		return nil, fmt.Errorf("token type is not access")
	}
	return claims, nil
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

func AuthMiddleware(rbacService rbac.AuthorizationService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := tokenFromRequest(r)
			if tokenString == "" {
				httperr.Write(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			claims, err := parseAccessToken(tokenString)
			if err != nil {
				fmt.Println("JWT parse error:", err)
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
