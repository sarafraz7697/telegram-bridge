package middleware

import (
	"net/http"

	"telegram-bridge/logger"
)

// AuthMiddleware creates a middleware that validates the Authorization header
func AuthMiddleware(expectedToken string, log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken != "Bearer "+expectedToken {
				log.Warn("Unauthorized access attempt",
					logger.String("path", r.URL.Path),
					logger.String("method", r.Method),
					logger.String("remote_addr", r.RemoteAddr))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
