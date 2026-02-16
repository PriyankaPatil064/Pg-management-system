package middleware

import (
	"context"
	"net/http"
	"pg-management-system/internal/handlers"
	"strings"
)

type contextKey string

const UserEmailKey contextKey = "user_email"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := handlers.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token. Please login again.", http.StatusUnauthorized)
			return
		}

		// Add email to context for potential use in handlers
		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
