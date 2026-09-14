package middleware

import (
	"context"
	"net/http"
	"pkg/auth"
	"pkg/httperr"

	"github.com/gin-gonic/gin"
)

func GinAuthMiddleware(authorizer AuthorizationChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := auth.TokenFromRequest(c.Request)
		if tokenString == "" {
			httperr.Write(c.Writer, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		claims, err := auth.ParseAccessToken(tokenString)
		if err != nil {
			httperr.Write(c.Writer, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), claimsContextKey, claims)
		c.Request = c.Request.WithContext(ctx)

		if authorizer == nil {
			httperr.Write(c.Writer, http.StatusInternalServerError, "access checker not configured")
			c.Abort()
			return
		}

		act := methodToAction(c.Request.Method)

		obj := CasbinObjFromGin(c)

		ok, err := authorizer.IsAuthenticated(claims.Id.String(), obj, act)
		if err != nil {
			httperr.Write(c.Writer, http.StatusInternalServerError, "access check error")
			c.Abort()
			return
		}
		if !ok {
			httperr.Write(c.Writer, http.StatusForbidden, "access forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
