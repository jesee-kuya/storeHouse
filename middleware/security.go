package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// SecurityHeaders returns a middleware that adds security headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Enable HSTS (HTTP Strict Transport Security)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"frame-ancestors 'none';"
		w.Header().Set("Content-Security-Policy", csp)

		// Feature Policy / Permissions Policy
		w.Header().Set("Permissions-Policy",
			"geolocation=(), "+
				"microphone=(), "+
				"camera=(), "+
				"payment=(), "+
				"usb=()")

		next.ServeHTTP(w, r)
	})
}

// NoCache returns a middleware that prevents caching
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Surrogate-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// CacheControl returns a middleware that sets cache control headers
func CacheControl(maxAge int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
			next.ServeHTTP(w, r)
		})
	}
}

// RequestSizeLimit returns a middleware that limits request size
func RequestSizeLimit(maxBytes int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				http.Error(w, "Request too large", http.StatusRequestEntityTooLarge)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ContentTypeValidation returns a middleware that validates content type
func ContentTypeValidation(allowedTypes []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
				contentType := r.Header.Get("Content-Type")
				if contentType != "" {
					// Extract base content type (before semicolon)
					baseType := strings.Split(contentType, ";")[0]
					baseType = strings.TrimSpace(baseType)

					if !containsString(allowedTypes, baseType) {
						http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
						return
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// IPWhitelist returns a middleware that allows only whitelisted IP addresses
func IPWhitelist(allowedIPs []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !containsString(allowedIPs, clientIP) && !containsString(allowedIPs, "*") {
				http.Error(w, "Access denied from this IP", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecureHeaders combines multiple security middleware
func SecureHeaders(next http.Handler) http.Handler {
	return SecurityHeaders(NoCache(next))
}

// DevelopmentSecurity returns security middleware suitable for development
func DevelopmentSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Basic security headers for development
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Allow framing for development tools
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")

		next.ServeHTTP(w, r)
	})
}

// ProductionSecurity returns security middleware suitable for production
func ProductionSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Full security headers for production
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy for production
		csp := "default-src 'self'; " +
			"script-src 'self'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self'; " +
			"connect-src 'self'; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"frame-ancestors 'none';"
		w.Header().Set("Content-Security-Policy", csp)

		next.ServeHTTP(w, r)
	})
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (when behind proxy)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Get the first IP in the list
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for X-Real-IP header
	xrip := r.Header.Get("X-Real-IP")
	if xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fall back to remote address
	return r.RemoteAddr
}

// containsString checks if a slice contains a specific string
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
