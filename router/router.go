package router

import (
	"net/http"

	"telegram-bridge/handlers"
	"telegram-bridge/logger"
	"telegram-bridge/middleware"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(handler *handlers.Handler, authToken string, log logger.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	// Auth middleware
	auth := middleware.AuthMiddleware(authToken, log)

	// Protected routes
	mux.Handle("/send", auth(methodFilter(http.MethodPost, http.HandlerFunc(handler.SendMessage), log)))

	// Public routes
	mux.Handle("/health", methodFilter(http.MethodGet, http.HandlerFunc(handler.Health), log))

	return mux
}

// methodFilter creates a middleware that only allows specific HTTP methods
func methodFilter(allowedMethod string, next http.Handler, log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			log.Warn("Method not allowed",
				logger.String("path", r.URL.Path),
				logger.String("method", r.Method),
				logger.String("allowed_method", allowedMethod))
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
