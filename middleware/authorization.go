package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Role represents user roles in the system
type Role string

const (
	RoleAdmin     Role = "Admin"
	RoleTreasurer Role = "Treasurer"
	RoleClerk     Role = "Clerk"
)

// RequireRole middleware ensures the authenticated user has the required role
func RequireRole(allowedRoles ...Role) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			// Check if user's role is in the allowed roles
			userRole := Role(user.Role)
			if !containsRole(allowedRoles, userRole) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{
					"error":          "Insufficient permissions",
					"required_roles": rolesToString(allowedRoles),
					"user_role":      string(userRole),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware that only allows Admin role
func RequireAdmin(next http.Handler) http.Handler {
	return RequireRole(RoleAdmin)(next)
}

// RequireAdminOrTreasurer middleware that allows Admin or Treasurer roles
func RequireAdminOrTreasurer(next http.Handler) http.Handler {
	return RequireRole(RoleAdmin, RoleTreasurer)(next)
}

// RequireAnyRole middleware that allows any authenticated role
func RequireAnyRole(next http.Handler) http.Handler {
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

// ResourceOwner middleware that checks if the authenticated user owns the resource
// or has Admin role
func ResourceOwner(next http.Handler) http.Handler {
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

		// Admin can access any resource
		if Role(user.Role) == RoleAdmin {
			next.ServeHTTP(w, r)
			return
		}

		// Get resource ID from URL parameter (assumes {id} parameter)
		resourceID := chi.URLParam(r, "id")
		if resourceID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Resource ID required",
			})
			return
		}

		// Check if user owns the resource
		if resourceID != user.ID {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Access denied: can only access your own resources",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RoleBasedRoute middleware for chi router that applies role-based access control
func RoleBasedRoute(pattern string, roles ...Role) func(next http.Handler) http.Handler {
	return RequireRole(roles...)
}

// containsRole checks if a role is in the allowed roles slice
func containsRole(roles []Role, role Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// rolesToString converts roles slice to comma-separated string
func rolesToString(roles []Role) string {
	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = string(role)
	}
	return strings.Join(roleStrings, ", ")
}
