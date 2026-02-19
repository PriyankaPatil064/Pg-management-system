package middleware

import (
	"context"
	"net/http"
	"pg-management-system/internal/handlers"
	"strings"
)

type contextKey string

const UserEmailKey contextKey = "user_email"
const UserClaimsKey contextKey = "user_claims"

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

		// Store email and full claims in context
		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
