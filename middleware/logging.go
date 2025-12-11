package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// RequestLoggerConfig holds configuration for request logging
type RequestLoggerConfig struct {
	SkipPaths     []string
	IncludeBody   bool
	IncludeHeader bool
	LogFormat     string
}

// DefaultLoggerConfig returns default logging configuration
func DefaultLoggerConfig() *RequestLoggerConfig {
	return &RequestLoggerConfig{
		SkipPaths:     []string{"/health", "/metrics"},
		IncludeBody:   true,
		IncludeHeader: false,
		LogFormat:     "json", // or "text"
	}
}

// RequestLogger returns a middleware that logs HTTP requests
func RequestLogger(config *RequestLoggerConfig) func(next http.Handler) http.Handler {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for specified paths
			if contains(config.SkipPaths, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Capture request data
			start := time.Now()

			// Log request
			if config.LogFormat == "json" {
				logRequestJSON(r, config)
			} else {
				logRequestText(r)
			}

			// Create response writer wrapper
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				statusCode:     200,
				bytesWritten:   0,
			}

			// Capture request body if enabled and not too large
			var requestBody string
			if config.IncludeBody && r.Body != nil && r.ContentLength > 0 && r.ContentLength < 1024*1024 { // 1MB limit
				bodyBytes, _ := io.ReadAll(r.Body)
				requestBody = string(bodyBytes)
				// Reset body for the handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Call next handler
			next.ServeHTTP(lrw, r)

			// Calculate duration
			duration := time.Since(start)

			// Log response
			if config.LogFormat == "json" {
				logResponseJSON(lrw, r, requestBody, duration, config)
			} else {
				logResponseText(lrw, r, duration)
			}
		})
	}
}

// SimpleLogger returns a simple logging middleware
func SimpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("📝 %s %s - Started", r.Method, r.URL.Path)

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
			bytesWritten:   0,
		}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("✅ %s %s - Completed %d in %v", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}

// StructuredLogger returns a structured JSON logging middleware
func StructuredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := generateRequestID()

		// Add request ID to response headers
		w.Header().Set("X-Request-ID", requestID)

		// Log request start
		logEntry := map[string]interface{}{
			"timestamp":   start.Format(time.RFC3339),
			"request_id":  requestID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"query":       r.URL.RawQuery,
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.Header.Get("User-Agent"),
			"event":       "request_started",
		}

		if body, err := captureRequestBody(r); err == nil && len(body) > 0 {
			logEntry["request_body"] = body
		}

		logJSON(logEntry)

		// Create response writer wrapper
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
			bytesWritten:   0,
		}

		// Call next handler
		next.ServeHTTP(lrw, r)

		// Calculate duration and log response
		duration := time.Since(start)

		responseEntry := map[string]interface{}{
			"timestamp":     time.Now().Format(time.RFC3339),
			"request_id":    requestID,
			"method":        r.Method,
			"path":          r.URL.Path,
			"status_code":   lrw.statusCode,
			"bytes_written": lrw.bytesWritten,
			"duration_ms":   duration.Milliseconds(),
			"event":         "request_completed",
		}

		if lrw.statusCode >= 400 {
			responseEntry["error"] = true
		}

		logJSON(responseEntry)
	})
}

// SecurityLogger returns a middleware that logs security-related events
func SecurityLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
			bytesWritten:   0,
		}

		next.ServeHTTP(lrw, r)

		// Log security events
		if isSecurityEvent(lrw.statusCode, r) {
			securityLog := map[string]interface{}{
				"timestamp":   time.Now().Format(time.RFC3339),
				"event_type":  "security",
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": lrw.statusCode,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.Header.Get("User-Agent"),
				"referer":     r.Header.Get("Referer"),
			}

			logJSON(securityLog)
		}
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code and bytes written
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(p)
	lrw.bytesWritten += n
	return n, err
}

// Helper functions for logging
func logRequestJSON(r *http.Request, config *RequestLoggerConfig) {
	logEntry := map[string]interface{}{
		"timestamp":   time.Now().Format(time.RFC3339),
		"level":       "info",
		"event":       "request",
		"method":      r.Method,
		"path":        r.URL.Path,
		"query":       r.URL.RawQuery,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.Header.Get("User-Agent"),
	}

	if config.IncludeHeader {
		logEntry["headers"] = r.Header
	}

	logJSON(logEntry)
}

func logRequestText(r *http.Request) {
	log.Printf("📝 %s %s %s - %s", r.Method, r.URL.Path, r.URL.RawQuery, r.RemoteAddr)
}

func logResponseJSON(lrw *loggingResponseWriter, r *http.Request, requestBody string, duration time.Duration, config *RequestLoggerConfig) {
	logEntry := map[string]interface{}{
		"timestamp":     time.Now().Format(time.RFC3339),
		"level":         "info",
		"event":         "response",
		"method":        r.Method,
		"path":          r.URL.Path,
		"status_code":   lrw.statusCode,
		"bytes_written": lrw.bytesWritten,
		"duration_ms":   duration.Milliseconds(),
	}

	if lrw.statusCode >= 400 {
		logEntry["level"] = "error"
	}

	if config.IncludeBody && requestBody != "" {
		logEntry["request_body"] = requestBody
	}

	logJSON(logEntry)
}

func logResponseText(lrw *loggingResponseWriter, r *http.Request, duration time.Duration) {
	statusIcon := "✅"
	if lrw.statusCode >= 400 {
		statusIcon = "❌"
	}

	log.Printf("%s %s %s - %d %d bytes in %v", statusIcon, r.Method, r.URL.Path, lrw.statusCode, lrw.bytesWritten, duration)
}

func logJSON(data map[string]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}
	log.Println(string(jsonData))
}

func generateRequestID() string {
	return string(rune(time.Now().UnixNano()))
}

func captureRequestBody(r *http.Request) (string, error) {
	if r.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// Reset body for the handler
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}

func isSecurityEvent(statusCode int, r *http.Request) bool {
	// Log security events: unauthorized, forbidden, bad requests, etc.
	return statusCode == http.StatusUnauthorized ||
		statusCode == http.StatusForbidden ||
		statusCode == http.StatusBadRequest ||
		statusCode == http.StatusRequestEntityTooLarge ||
		r.URL.Path == "/admin"
}

