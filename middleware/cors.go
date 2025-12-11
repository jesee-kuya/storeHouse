package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORSConfig holds CORS configuration options
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	MaxAge           int
	AllowCredentials bool
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"X-Total-Count",
		},
		MaxAge:           86400, // 24 hours
		AllowCredentials: false,
	}
}

// ProductionCORSConfig returns CORS configuration for production
func ProductionCORSConfig() *CORSConfig {
	var allowedOrigins []string

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		allowedOrigins = strings.Split(origins, ",")
	} else {
		allowedOrigins = []string{"https://yourdomain.com"}
	}

	return &CORSConfig{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"X-Total-Count",
		},
		MaxAge:           86400,
		AllowCredentials: true,
	}
}

// CORS returns a CORS middleware handler
func CORS(config *CORSConfig) func(next http.Handler) http.Handler {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if len(config.AllowedOrigins) > 0 && !contains(config.AllowedOrigins, "*") && !contains(config.AllowedOrigins, origin) {
				next.ServeHTTP(w, r)
				return
			}

			// Set CORS headers
			if origin != "" && len(config.AllowedOrigins) > 0 {
				if contains(config.AllowedOrigins, "*") {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if len(config.ExposedHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			}

			if config.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
			}

			if len(config.AllowedMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			}

			if len(config.AllowedHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SimpleCORS returns a simple CORS middleware with permissive settings
// This is useful for development
func SimpleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow common headers
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Requested-With")

		// Allow common methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SecureCORS returns CORS middleware with more restrictive settings
// This is suitable for production environments
func SecureCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only allow specific origins in production
		allowedOrigins := []string{
			"http://localhost:3000", // React dev server
			"http://localhost:8080", // Flutter web
			"https://yourdomain.com",
		}

		if contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Don't allow credentials for security
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

		// Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
