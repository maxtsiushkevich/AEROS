package helpers

import (
	"context"
	"pkg/auth"

	"github.com/gin-gonic/gin"
)

type contextKey string

const ClaimsContextKey contextKey = "auth_claims"

func ClaimsFromContextGin(c *gin.Context) (*auth.Claims, bool) {
	claims, ok := c.Request.Context().Value(ClaimsContextKey).(*auth.Claims)
	return claims, ok
}

func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*auth.Claims)
	return claims, ok
}
