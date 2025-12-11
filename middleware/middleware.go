package middleware

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// MiddlewareStack represents a collection of middleware functions
type MiddlewareStack struct {
	auth       func(http.Handler) http.Handler
	authOpt    func(http.Handler) http.Handler
	authorize  func(allowedRoles ...Role) func(next http.Handler) http.Handler
	cors       func(next http.Handler) http.Handler
	security   func(next http.Handler) http.Handler
	rateLimit  func(next http.Handler) http.Handler
	logging    func(next http.Handler) http.Handler
	validation func(next http.Handler) http.Handler
}

// MiddlewareOption functional option pattern for configuring middleware
type MiddlewareOption func(*MiddlewareStack)

// WithAuth configures authentication middleware
func WithAuth(db *sqlx.DB) MiddlewareOption {
	return func(ms *MiddlewareStack) {
		ms.auth = AuthMiddleware(db)
		ms.authOpt = OptionalAuth
	}
}

// WithAuthorization configures authorization middleware
func WithAuthorization() MiddlewareOption {
	return func(ms *MiddlewareStack) {
		ms.authorize = RequireRole
	}
}

// WithCORS configures CORS middleware
func WithCORS(config *CORSConfig) MiddlewareOption {
	return func(ms *MiddlewareStack) {
		ms.cors = CORS(config)
	}
}

// WithSecurity configures security middleware
func WithSecurity(secureHeaders bool) MiddlewareOption {
	return func(ms *MiddlewareStack) {
		if secureHeaders {
			ms.security = SecurityHeaders
		}
	}
}

// WithRateLimit configures rate limiting middleware
func WithRateLimit(limit int, window time.Duration) MiddlewareOption {
	return func(ms *MiddlewareStack) {
		ms.rateLimit = RateLimitByIP(limit, window)
	}
}

// WithLogging configures logging middleware
func WithLogging(config *RequestLoggerConfig) MiddlewareOption {
	return func(ms *MiddlewareStack) {
		if config != nil {
			ms.logging = RequestLogger(config)
		} else {
			ms.logging = SimpleLogger
		}
	}
}

// WithValidation configures validation middleware
func WithValidation() MiddlewareOption {
	return func(ms *MiddlewareStack) {
		ms.validation = ValidateJSON
	}
}

// NewMiddlewareStack creates a new middleware stack with default configuration
func NewMiddlewareStack(options ...MiddlewareOption) *MiddlewareStack {
	stack := &MiddlewareStack{}

	// Apply default middleware
	stack.cors = SimpleCORS
	stack.security = DevelopmentSecurity
	stack.logging = SimpleLogger
	stack.rateLimit = RateLimitByIP(100, time.Minute)

	// Apply provided options
	for _, option := range options {
		option(stack)
	}

	return stack
}

// ApplyAll applies all configured middleware to the given handler
func (ms *MiddlewareStack) ApplyAll(next http.Handler) http.Handler {
	handler := next

	// Apply middleware in reverse order (right to left)
	if ms.validation != nil {
		handler = ms.validation(handler)
	}

	if ms.rateLimit != nil {
		handler = ms.rateLimit(handler)
	}

	if ms.logging != nil {
		handler = ms.logging(handler)
	}

	if ms.security != nil {
		handler = ms.security(handler)
	}

	if ms.cors != nil {
		handler = ms.cors(handler)
	}

	if ms.auth != nil {
		handler = ms.auth(handler)
	}

	return handler
}

// ApplyAuth applies only authentication-related middleware
func (ms *MiddlewareStack) ApplyAuth(next http.Handler) http.Handler {
	handler := next

	if ms.authorize != nil {
		handler = ms.authorize()(handler)
	}

	if ms.authOpt != nil {
		handler = ms.authOpt(handler)
	}

	return handler
}

// ApplyPublic applies middleware for public endpoints (no auth required)
func (ms *MiddlewareStack) ApplyPublic(next http.Handler) http.Handler {
	handler := next

	if ms.validation != nil {
		handler = ms.validation(handler)
	}

	if ms.rateLimit != nil {
		handler = ms.rateLimit(handler)
	}

	if ms.logging != nil {
		handler = ms.logging(handler)
	}

	if ms.security != nil {
		handler = ms.security(handler)
	}

	if ms.cors != nil {
		handler = ms.cors(handler)
	}

	return handler
}

// ApplySecure applies middleware for admin/sensitive endpoints
func (ms *MiddlewareStack) ApplySecure(next http.Handler) http.Handler {
	handler := next

	// Apply in order
	if ms.validation != nil {
		handler = ms.validation(handler)
	}

	if ms.rateLimit != nil {
		handler = RateLimitByIP(10, time.Minute)(handler) // Stricter rate limiting
	}

	if ms.logging != nil {
		handler = SecurityLogger(handler) // Security-focused logging
	}

	if ms.security != nil {
		handler = SecureHeaders(handler)
	}

	if ms.cors != nil {
		handler = ms.cors(handler)
	}

	if ms.auth != nil {
		handler = ms.auth(handler)
	}

	return handler
}

// Common middleware presets

// DefaultMiddleware returns a standard middleware stack for general API use
func DefaultMiddleware(db *sqlx.DB) *MiddlewareStack {
	return NewMiddlewareStack(
		WithAuth(db),
		WithAuthorization(),
		WithCORS(DefaultCORSConfig()),
		WithSecurity(true),
		WithRateLimit(100, time.Minute),
		WithLogging(DefaultLoggerConfig()),
		WithValidation(),
	)
}

// PublicMiddleware returns middleware for public endpoints (no authentication)
func PublicMiddleware() *MiddlewareStack {
	return NewMiddlewareStack(
		WithCORS(DefaultCORSConfig()),
		WithSecurity(true),
		WithRateLimit(1000, time.Hour), // More permissive for public endpoints
		WithLogging(DefaultLoggerConfig()),
		WithValidation(),
	)
}

// AdminMiddleware returns middleware for admin endpoints
func AdminMiddleware(db *sqlx.DB) *MiddlewareStack {
	return NewMiddlewareStack(
		WithAuth(db),
		WithAuthorization(),
		WithCORS(ProductionCORSConfig()),
		WithSecurity(true),
		WithRateLimit(10, time.Minute), // Very restrictive
		WithLogging(DefaultLoggerConfig()),
		WithValidation(),
	)
}

// DevelopmentMiddleware returns middleware suitable for development
func DevelopmentMiddleware(db *sqlx.DB) *MiddlewareStack {
	return NewMiddlewareStack(
		WithAuth(db),
		WithAuthorization(),
		WithCORS(nil),                  // Will use default SimpleCORS
		WithSecurity(false),            // Less restrictive for dev
		WithRateLimit(1000, time.Hour), // Very permissive for dev
		WithLogging(nil),               // Will use SimpleLogger
		WithValidation(),
	)
}

// ProductionMiddleware returns middleware suitable for production
func ProductionMiddleware(db *sqlx.DB) *MiddlewareStack {
	return NewMiddlewareStack(
		WithAuth(db),
		WithAuthorization(),
		WithCORS(ProductionCORSConfig()),
		WithSecurity(true),
		WithRateLimit(50, time.Minute),
		WithLogging(DefaultLoggerConfig()),
		WithValidation(),
	)
}

// Helper functions for common middleware combinations

// AdminOnly applies admin-only middleware
func AdminOnly(db *sqlx.DB, next http.Handler) http.Handler {
	return RequireAdmin(AuthMiddleware(db)(next))
}

// TreasurerOrAdmin applies middleware for treasurer or admin roles
func TreasurerOrAdmin(db *sqlx.DB, next http.Handler) http.Handler {
	return RequireAdminOrTreasurer(AuthMiddleware(db)(next))
}

// Authenticated applies authentication middleware
func Authenticated(db *sqlx.DB, next http.Handler) http.Handler {
	return RequireAuth(AuthMiddleware(db)(next))
}

// Public applies public middleware (no authentication)
func Public(next http.Handler) http.Handler {
	return SimpleCORS(SecurityHeaders(SimpleLogger(next)))
}
