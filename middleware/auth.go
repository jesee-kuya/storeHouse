package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// AuthContextKey is the key for storing user information in context
const AuthContextKey contextKey = "user"

// User represents the authenticated user in context
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// AuthMiddleware validates JWT tokens and sets user context
func AuthMiddleware(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			const tokenPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, tokenPrefix) {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, tokenPrefix)
			if token == "" {
				http.Error(w, "Token required", http.StatusUnauthorized)
				return
			}

			// For this example, we'll validate against database
			// In a real implementation, you'd verify JWT signature
			var user User
			err := db.Get(&user,
				"SELECT id, username, email, role FROM users WHERE id = $1 AND is_active = true",
				token)

			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), AuthContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(r *http.Request) *User {
	if user, ok := r.Context().Value(AuthContextKey).(User); ok {
		return &user
	}
	return nil
}

// RequireAuth middleware that ensures user is authenticated
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Authentication required",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// OptionalAuth middleware that works with or without authentication
func OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user != nil {
			ctx := context.WithValue(r.Context(), AuthContextKey, *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
