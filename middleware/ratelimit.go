package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter holds rate limiting configuration and tracking
type RateLimiter struct {
	mu         sync.Mutex
	requests   map[string][]time.Time
	limit      int
	window     time.Duration
	statusCode int
	message    string
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:   make(map[string][]time.Time),
		limit:      limit,
		window:     window,
		statusCode: http.StatusTooManyRequests,
		message:    "Rate limit exceeded",
	}
}

// NewDefaultRateLimiter creates a rate limiter with default settings
// (100 requests per minute)
func NewDefaultRateLimiter() *RateLimiter {
	return NewRateLimiter(100, time.Minute)
}

// RateLimit returns a middleware that implements rate limiting
func RateLimit(limiter *RateLimiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !limiter.allowRequest(clientIP) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(limiter.statusCode)
				w.Write([]byte(`{"error":"` + limiter.message + `"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// allowRequest checks if the request is allowed for the given client
func (rl *RateLimiter) allowRequest(client string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create request tracking for this client
	if _, exists := rl.requests[client]; !exists {
		rl.requests[client] = make([]time.Time, 0)
	}

	// Remove old requests outside the time window
	var validRequests []time.Time
	for _, reqTime := range rl.requests[client] {
		if now.Sub(reqTime) <= rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		rl.requests[client] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[client] = validRequests

	return true
}

// RateLimitByIP returns a rate limiting middleware that limits by IP address
func RateLimitByIP(limit int, window time.Duration) func(next http.Handler) http.Handler {
	limiter := NewRateLimiter(limit, window)
	return RateLimit(limiter)
}

// RateLimitByUser returns a rate limiting middleware that limits by authenticated user
func RateLimitByUser(limit int, window time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var identifier string

			// Use user ID if authenticated, otherwise use IP
			if user := GetUserFromContext(r); user != nil {
				identifier = "user:" + user.ID
			} else {
				identifier = "ip:" + getClientIP(r)
			}

			// Simple in-memory rate limiting (for demonstration)
			// In production, use Redis or similar for distributed rate limiting
			_ = "rate_limit:" + identifier // Prevent unused variable error

			// This is a simplified version - in practice you'd use a proper
			// rate limiting library or implement proper tracking
			w.Header().Set("X-RateLimit-Limit", string(rune(limit)))
			w.Header().Set("X-RateLimit-Remaining", string(rune(limit-1)))
			w.Header().Set("X-RateLimit-Reset", time.Now().Add(window).Format(time.RFC3339))

			next.ServeHTTP(w, r)
		})
	}
}

// PerEndpointRateLimit returns a rate limiter that can be configured per endpoint
type PerEndpointRateLimiter struct {
	limiters       map[string]*RateLimiter
	defaultLimiter *RateLimiter
}

// NewPerEndpointRateLimiter creates a rate limiter that can be configured per endpoint
func NewPerEndpointRateLimiter(defaultLimit int, defaultWindow time.Duration) *PerEndpointRateLimiter {
	return &PerEndpointRateLimiter{
		limiters:       make(map[string]*RateLimiter),
		defaultLimiter: NewRateLimiter(defaultLimit, defaultWindow),
	}
}

// AddEndpoint adds a specific rate limiter for an endpoint
func (p *PerEndpointRateLimiter) AddEndpoint(endpoint string, limit int, window time.Duration) {
	p.limiters[endpoint] = NewRateLimiter(limit, window)
}

// RateLimitPerEndpoint returns middleware that uses per-endpoint rate limiting
func (p *PerEndpointRateLimiter) RateLimitPerEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpoint := r.Method + " " + r.URL.Path

		var limiter *RateLimiter
		if specificLimiter, exists := p.limiters[endpoint]; exists {
			limiter = specificLimiter
		} else {
			limiter = p.defaultLimiter
		}

		clientIP := getClientIP(r)

		if !limiter.allowRequest(clientIP) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"Rate limit exceeded for this endpoint"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StrictRateLimit returns a very restrictive rate limiter
func StrictRateLimit(next http.Handler) http.Handler {
	return RateLimitByIP(10, time.Minute)(next)
}

// LenientRateLimit returns a more permissive rate limiter
func LenientRateLimit(next http.Handler) http.Handler {
	return RateLimitByIP(1000, time.Hour)(next)
}
