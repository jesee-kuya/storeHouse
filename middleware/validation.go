package middleware

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ValidationRule represents a single validation rule
type ValidationRule struct {
	Field        string
	Required     bool
	MinLength    int
	MaxLength    int
	Pattern      string
	CustomFunc   func(string) bool
	ErrorMessage string
}

// ValidateRequest validates incoming requests based on rules
func ValidateRequest(rules []ValidationRule) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "DELETE" {
				// For GET and DELETE, validate query parameters
				errors := validateQueryParams(r.URL.Query(), rules)
				if len(errors) > 0 {
					sendValidationError(w, errors)
					return
				}
			} else {
				// For POST, PUT, PATCH, validate JSON body
				if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
					http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
					return
				}

				var data map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					http.Error(w, "Invalid JSON format", http.StatusBadRequest)
					return
				}

				errors := validateJSONBody(data, rules)
				if len(errors) > 0 {
					sendValidationError(w, errors)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateJSON validates that the request body is valid JSON
func ValidateJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Body == nil {
				http.Error(w, "Request body is required", http.StatusBadRequest)
				return
			}

			// Check content type
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}

			// Try to decode JSON
			var data map[string]interface{}
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&data); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			// Re-encode to reset the body for the next handler
			r.Body = &rewindableBody{data: data}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateRequiredFields validates that required fields are present
func ValidateRequiredFields(fields ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" || r.Method == "DELETE" {
				// Validate query parameters
				for _, field := range fields {
					if r.URL.Query().Get(field) == "" {
						http.Error(w, "Field '"+field+"' is required", http.StatusBadRequest)
						return
					}
				}
			} else {
				// Validate JSON body
				var data map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
					http.Error(w, "Invalid JSON format", http.StatusBadRequest)
					return
				}

				for _, field := range fields {
					if _, exists := data[field]; !exists || data[field] == nil || data[field] == "" {
						http.Error(w, "Field '"+field+"' is required", http.StatusBadRequest)
						return
					}
				}

				// Re-encode to reset the body
				r.Body = &rewindableBody{data: data}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ValidateEmail validates email format
func ValidateEmail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

			if email, ok := data["email"].(string); ok && !emailRegex.MatchString(email) {
				http.Error(w, "Invalid email format", http.StatusBadRequest)
				return
			}

			// Re-encode to reset the body
			r.Body = &rewindableBody{data: data}
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateAmount validates monetary amounts
func ValidateAmount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			if amount, ok := data["amount"].(string); ok {
				if _, err := strconv.ParseFloat(amount, 64); err != nil {
					http.Error(w, "Amount must be a valid number", http.StatusBadRequest)
					return
				}
			}

			// Re-encode to reset the body
			r.Body = &rewindableBody{data: data}
		}

		next.ServeHTTP(w, r)
	})
}

// PaginationMiddleware adds pagination validation and parameters
func PaginationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get pagination parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")
		pageStr := r.URL.Query().Get("page")

		// Set defaults
		limit := 10
		offset := 0

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 || l > 100 {
				http.Error(w, "Limit must be between 1 and 100", http.StatusBadRequest)
				return
			} else {
				limit = l
			}
		}

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
				http.Error(w, "Page must be a positive integer", http.StatusBadRequest)
				return
			} else {
				offset = (p - 1) * limit
			}
		} else if offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err != nil || o < 0 {
				http.Error(w, "Offset must be a non-negative integer", http.StatusBadRequest)
				return
			} else {
				offset = o
			}
		}

		// Add pagination info to request context
		// In a real implementation, you'd use context values
		w.Header().Set("X-Page-Limit", string(rune(limit)))
		w.Header().Set("X-Page-Offset", string(rune(offset)))

		next.ServeHTTP(w, r)
	})
}

// validateQueryParams validates query parameters
func validateQueryParams(query map[string][]string, rules []ValidationRule) []string {
	var errors []string

	for _, rule := range rules {
		values := query[rule.Field]
		if len(values) == 0 {
			if rule.Required {
				errors = append(errors, "Field '"+rule.Field+"' is required")
			}
			continue
		}

		value := values[0]
		ruleErrors := validateField(value, rule)
		errors = append(errors, ruleErrors...)
	}

	return errors
}

// validateJSONBody validates JSON body fields
func validateJSONBody(data map[string]interface{}, rules []ValidationRule) []string {
	var errors []string

	for _, rule := range rules {
		value, exists := data[rule.Field]

		if !exists || value == nil {
			if rule.Required {
				errors = append(errors, "Field '"+rule.Field+"' is required")
			}
			continue
		}

		valueStr, ok := value.(string)
		if !ok {
			errors = append(errors, "Field '"+rule.Field+"' must be a string")
			continue
		}

		ruleErrors := validateField(valueStr, rule)
		errors = append(errors, ruleErrors...)
	}

	return errors
}

// validateField validates a single field against rules
func validateField(value string, rule ValidationRule) []string {
	var errors []string

	if rule.MinLength > 0 && len(value) < rule.MinLength {
		errors = append(errors, "Field '"+rule.Field+"' must be at least "+string(rune(rule.MinLength))+" characters")
	}

	if rule.MaxLength > 0 && len(value) > rule.MaxLength {
		errors = append(errors, "Field '"+rule.Field+"' must be at most "+string(rune(rule.MaxLength))+" characters")
	}

	if rule.Pattern != "" {
		regex := regexp.MustCompile(rule.Pattern)
		if !regex.MatchString(value) {
			errorMsg := rule.ErrorMessage
			if errorMsg == "" {
				errorMsg = "Field '" + rule.Field + "' format is invalid"
			}
			errors = append(errors, errorMsg)
		}
	}

	if rule.CustomFunc != nil && !rule.CustomFunc(value) {
		errorMsg := rule.ErrorMessage
		if errorMsg == "" {
			errorMsg = "Field '" + rule.Field + "' is invalid"
		}
		errors = append(errors, errorMsg)
	}

	return errors
}

// sendValidationError sends validation errors as JSON response
func sendValidationError(w http.ResponseWriter, errors []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := map[string]interface{}{
		"error":   "Validation failed",
		"details": errors,
	}

	json.NewEncoder(w).Encode(response)
}

// rewindableBody is a wrapper to allow re-reading the request body
type rewindableBody struct {
	data map[string]interface{}
}

func (rb *rewindableBody) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rb *rewindableBody) Close() error {
	return nil
}
