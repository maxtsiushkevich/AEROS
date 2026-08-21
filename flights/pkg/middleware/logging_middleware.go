package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next(w, r)

			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		}
	}
}
