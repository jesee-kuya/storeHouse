# Middleware Package

This package provides comprehensive middleware functions for the expense tracker API built with Go and the Chi router. The middleware is designed to enhance security, performance, and functionality of HTTP requests.

## Features

- **Authentication & Authorization**: JWT-based auth with role-based access control
- **CORS Support**: Configurable cross-origin resource sharing
- **Security Headers**: Protection against common web vulnerabilities
- **Rate Limiting**: Protection against abuse and DoS attacks
- **Request Validation**: Input validation and sanitization
- **Logging**: Comprehensive request/response logging
- **Middleware Stack**: Easy-to-use middleware combinations

## Quick Start

```go
package main

import (
    "storeHouse/database"
    "storeHouse/hanlers"
    "storeHouse/middleware"
    "github.com/jmoiron/sqlx"
)

func main() {
    db := database.ConnectDB()
    defer db.Close()
    
    // Create middleware stack
    mw := middleware.DefaultMiddleware(db)
    
    // Create router with middleware
    router := hanlers.Router(db)
    
    // Apply middleware
    router = mw.ApplyAll(router)
    
    // Start server
    hanlers.StartServer(db)
}
```

## Middleware Types

### 1. Authentication Middleware (`auth.go`)

Handles user authentication and session management.

#### Usage Examples:

```go
// Require authentication for all requests
router.Use(middleware.AuthMiddleware(db))

// Optional authentication (works with or without auth)
router.Use(middleware.OptionalAuth)

// Get user from context
func handler(w http.ResponseWriter, r *http.Request) {
    user := middleware.GetUserFromContext(r)
    if user == nil {
        // User not authenticated
        return
    }
    // Use authenticated user info
}
```

### 2. Authorization Middleware (`authorization.go`)

Implements role-based access control (RBAC).

#### Available Roles:
- `Admin`: Full system access
- `Treasurer`: Financial operations access
- `Clerk`: Limited access

#### Usage Examples:

```go
// Require specific role
router.Route("/admin", func(r chi.Router) {
    r.Use(middleware.RequireAdmin)
    r.Get("/", adminHandler)
})

// Multiple allowed roles
router.Route("/reports", func(r chi.Router) {
    r.Use(middleware.RequireRole(middleware.RoleAdmin, middleware.RoleTreasurer))
    r.Get("/", reportsHandler)
})

// Resource owner check
router.Route("/profile", func(r chi.Router) {
    r.Use(middleware.ResourceOwner)
    r.Get("/{id}", getProfile)
})
```

### 3. CORS Middleware (`cors.go`)

Manages cross-origin requests.

#### Usage Examples:

```go
// Default CORS (development)
router.Use(middleware.SimpleCORS)

// Production CORS
router.Use(middleware.CORS(middleware.ProductionCORSConfig()))

// Custom CORS
config := &middleware.CORSConfig{
    AllowedOrigins: []string{"https://yourdomain.com"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowCredentials: true,
}
router.Use(middleware.CORS(config))
```

### 4. Security Middleware (`security.go`)

Adds security headers and protection.

#### Features:
- Content Security Policy (CSP)
- XSS Protection
- Clickjacking Prevention
- MIME Sniffing Protection
- HSTS (HTTP Strict Transport Security)

#### Usage Examples:

```go
// Add security headers
router.Use(middleware.SecurityHeaders)

// No caching for sensitive data
router.Use(middleware.NoCache)

// Cache control for public data
router.Use(middleware.CacheControl(3600)) // 1 hour

// Request size limit
router.Use(middleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB

// Content type validation
router.Use(middleware.ContentTypeValidation([]string{"application/json"}))

// IP whitelist
router.Use(middleware.IPWhitelist([]string{"192.168.1.1", "10.0.0.0/8"}))
```

### 5. Rate Limiting Middleware (`ratelimit.go`)

Prevents abuse and DoS attacks.

#### Usage Examples:

```go
// Simple rate limiting
router.Use(middleware.RateLimitByIP(100, time.Minute))

// User-based rate limiting
router.Use(middleware.RateLimitByUser(1000, time.Hour))

// Strict rate limiting for sensitive endpoints
router.Route("/admin", func(r chi.Router) {
    r.Use(middleware.StrictRateLimit)
    r.Get("/", adminHandler)
})
```

### 6. Validation Middleware (`validation.go`)

Validates incoming requests.

#### Usage Examples:

```go
// JSON validation
router.Use(middleware.ValidateJSON)

// Required fields
router.Use(middleware.ValidateRequiredFields("email", "password"))

// Email validation
router.Use(middleware.ValidateEmail)

// Amount validation
router.Use(middleware.ValidateAmount)

// Custom validation rules
rules := []middleware.ValidationRule{
    {
        Field:       "email",
        Required:    true,
        Pattern:     `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
        ErrorMessage: "Invalid email format",
    },
    {
        Field:     "password",
        Required:  true,
        MinLength: 8,
    },
}
router.Use(middleware.ValidateRequest(rules))

// Pagination
router.Use(middleware.PaginationMiddleware)
```

### 7. Logging Middleware (`logging.go`)

Comprehensive request and response logging.

#### Usage Examples:

```go
// Simple logging
router.Use(middleware.SimpleLogger)

// Structured JSON logging
router.Use(middleware.StructuredLogger)

// Custom logging configuration
config := &middleware.RequestLoggerConfig{
    SkipPaths:     []string{"/health"},
    IncludeBody:   false,
    IncludeHeader: false,
    LogFormat:     "json",
}
router.Use(middleware.RequestLogger(config))

// Security-focused logging
router.Use(middleware.SecurityLogger)
```

## Middleware Stack

The package provides pre-configured middleware stacks for different use cases:

### Default Stack
```go
mw := middleware.DefaultMiddleware(db)
router = mw.ApplyAll(router)
```

### Public Endpoints (No Auth)
```go
mw := middleware.PublicMiddleware()
router = mw.ApplyPublic(router)
```

### Admin Endpoints
```go
mw := middleware.AdminMiddleware(db)
router = mw.ApplySecure(router)
```

### Development
```go
mw := middleware.DevelopmentMiddleware(db)
router = mw.ApplyAll(router)
```

### Production
```go
mw := middleware.ProductionMiddleware(db)
router = mw.ApplyAll(router)
```

## Custom Middleware Combinations

Create custom middleware combinations:

```go
// Custom stack
stack := middleware.NewMiddlewareStack(
    middleware.WithAuth(db),
    middleware.WithAuthorization(),
    middleware.WithCORS(config),
    middleware.WithSecurity(true),
    middleware.WithRateLimit(50, time.Minute),
    middleware.WithLogging(nil), // Uses SimpleLogger
)

// Apply to specific routes
stack.ApplyAuth(handler)    // Auth + Authorization
stack.ApplyPublic(handler)  // Public access
stack.ApplySecure(handler)  // Admin/Sensitive
```

## Environment-Specific Configuration

### Development
```go
// Disable some security features for development
router.Use(middleware.DevelopmentSecurity)
router.Use(middleware.SimpleCORS)
router.Use(middleware.LenientRateLimit)
```

### Production
```go
// Full security in production
router.Use(middleware.ProductionSecurity)
router.Use(middleware.SecureCORS)
router.Use(middleware.RateLimitByIP(50, time.Minute))
```

## Best Practices

### 1. Order Matters
Apply middleware in the correct order:

```go
router.Use(
    middleware.CORS,           // CORS first
    middleware.SecurityHeaders, // Security headers
    middleware.RateLimit,      // Rate limiting
    middleware.AuthMiddleware, // Authentication
    middleware.SimpleLogger,   // Logging last
)
```

### 2. Endpoint-Specific Middleware
Apply stricter middleware to sensitive endpoints:

```go
// Public endpoint
router.Get("/health", middleware.Public(handler))

// Protected endpoint
router.Get("/transactions", 
    middleware.Authenticated(db, 
        middleware.RequireRole(middleware.RoleTreasurer, middleware.RoleAdmin)(handler)))

// Admin endpoint
router.Get("/admin/users", 
    middleware.AdminOnly(db, handler))
```

### 3. Environment Variables
Use environment variables for configuration:

```go
func getRateLimit() (int, time.Duration) {
    limit := 100 // default
    if val := os.Getenv("RATE_LIMIT"); val != "" {
        if l, err := strconv.Atoi(val); err == nil {
            limit = l
        }
    }
    return limit, time.Minute
}

router.Use(middleware.RateLimitByIP(getRateLimit()))
```

## Error Handling

All middleware returns appropriate HTTP status codes and JSON error responses:

```json
{
    "error": "Authentication required"
}
```

```json
{
    "error": "Insufficient permissions",
    "required_roles": "Admin",
    "user_role": "Clerk"
}
```

```json
{
    "error": "Rate limit exceeded"
}
```

## Performance Considerations

1. **Logging**: Disable request body logging in production for sensitive data
2. **Rate Limiting**: Use Redis or similar for distributed rate limiting
3. **Validation**: Cache validation rules when possible
4. **Security Headers**: Minimal performance impact, always enabled

## Security Notes

1. Always use HTTPS in production
2. Implement proper JWT token validation (current implementation uses database lookup)
3. Configure CORS properly for your domain
4. Use environment variables for sensitive configuration
5. Monitor rate limiting and security logs

## Integration with Router

Here's how to integrate with the existing router:

```go
func Router(db *sqlx.DB) *chi.Mux {
    router := chi.NewRouter()
    
    // Create middleware stack
    mw := middleware.DefaultMiddleware(db)
    
    // Apply global middleware
    router.Use(mw.ApplyAll)
    
    // Or apply specific middleware
    router.Use(middleware.SimpleLogger)
    router.Use(middleware.SecurityHeaders)
    
    // ... rest of your routes
}
```

This middleware package provides a comprehensive foundation for building secure, scalable, and maintainable APIs with proper authentication, authorization, and security measures.