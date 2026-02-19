package middleware

import (
	"context"
	"net/http"
	"pg-management-system/internal/handlers"
)

const UserRoleKey contextKey = "user_role"

// RBAC returns a middleware that checks if the user has one of the allowed roles
func RBAC(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The token validation and claims extraction should already be done by an Auth middleware
			// and stored in the context. However, if not, we can extract it here if the token is in the header.

			// For now, let's assume we extract it from the token in the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Assuming format "Bearer <token>"
			tokenString := authHeader
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}

			claims, err := handlers.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userRole := claims.Role

			isAllowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// Add role to context for further use if needed
			ctx := context.WithValue(r.Context(), UserRoleKey, userRole)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
