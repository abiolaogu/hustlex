package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"hustlex/internal/infrastructure/security/audit"
)

// Context keys for request metadata
type contextKey string

const (
	ContextKeyRequestID     contextKey = "request_id"
	ContextKeyCorrelationID contextKey = "correlation_id"
	ContextKeyIPAddress     contextKey = "ip_address"
	ContextKeyUserAgent     contextKey = "user_agent"
	ContextKeyRequestStart  contextKey = "request_start"
)

// EnhancedSecurityHeaders adds comprehensive security headers to responses (OWASP recommendations)
func EnhancedSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// XSS Protection (legacy, but still useful for older browsers)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"font-src 'self'; "+
				"connect-src 'self' https://api.paystack.co https://checkout.paystack.com; "+
				"frame-ancestors 'none'; "+
				"form-action 'self';")

		// HSTS (Strict Transport Security) - 1 year with preload
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Permissions Policy (Feature Policy replacement)
		w.Header().Set("Permissions-Policy",
			"geolocation=(), "+
				"microphone=(), "+
				"camera=(), "+
				"payment=(), "+
				"usb=(), "+
				"magnetometer=(), "+
				"gyroscope=(), "+
				"accelerometer=()")

		// Cache Control for API responses - prevent caching of sensitive data
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private, max-age=0")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}

		// Remove server identification headers
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")

		next.ServeHTTP(w, r)
	})
}

// RequestSanitizer validates and sanitizes incoming requests
func RequestSanitizer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Limit request body size (10MB for file uploads, 1MB for JSON)
		maxBodySize := int64(1 << 20) // 1MB default
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			maxBodySize = 10 << 20 // 10MB for file uploads
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

		// Validate Content-Type for POST/PUT/PATCH
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if contentType != "" &&
				!strings.HasPrefix(contentType, "application/json") &&
				!strings.HasPrefix(contentType, "multipart/form-data") &&
				!strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
				http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
				return
			}
		}

		// Block suspicious paths
		suspiciousPaths := []string{
			"/.env", "/.git", "/wp-admin", "/wp-login", "/phpinfo",
			"/.htaccess", "/config.php", "/admin.php", "/shell",
		}
		for _, path := range suspiciousPaths {
			if strings.Contains(strings.ToLower(r.URL.Path), path) {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// EnhancedRequestID adds a unique request ID and correlation ID to each request
func EnhancedRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists (from upstream proxy)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Check for correlation ID (for distributed tracing)
		correlationID := r.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = requestID
		}

		// Add to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyRequestID, requestID)
		ctx = context.WithValue(ctx, ContextKeyCorrelationID, correlationID)
		ctx = context.WithValue(ctx, ContextKeyIPAddress, getClientIP(r))
		ctx = context.WithValue(ctx, ContextKeyUserAgent, r.UserAgent())
		ctx = context.WithValue(ctx, ContextKeyRequestStart, time.Now())

		// Add to response headers
		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Correlation-ID", correlationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SecureCORS handles Cross-Origin Resource Sharing with security defaults
func SecureCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	originSet := make(map[string]bool)
	for _, origin := range allowedOrigins {
		originSet[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if origin != "" && (originSet["*"] || originSet[origin]) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers",
					"Accept, Authorization, Content-Type, X-Request-ID, X-Correlation-ID")
				w.Header().Set("Access-Control-Expose-Headers",
					"X-Request-ID, X-Correlation-ID, X-RateLimit-Remaining")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuditMiddleware logs all requests for compliance
func AuditMiddleware(logger audit.AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap response writer to capture status code
			wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Execute handler
			next.ServeHTTP(wrapped, r)

			// Skip health checks and static assets
			if strings.HasPrefix(r.URL.Path, "/health") ||
				strings.HasPrefix(r.URL.Path, "/static") ||
				strings.HasPrefix(r.URL.Path, "/favicon") {
				return
			}

			// Determine outcome
			outcome := audit.OutcomeSuccess
			if wrapped.statusCode >= 400 && wrapped.statusCode < 500 {
				outcome = audit.OutcomeFailure
			} else if wrapped.statusCode >= 500 {
				outcome = audit.OutcomeError
			}

			// Determine action
			action := audit.ActionRead
			switch r.Method {
			case http.MethodPost:
				action = audit.ActionCreate
			case http.MethodPut, http.MethodPatch:
				action = audit.ActionUpdate
			case http.MethodDelete:
				action = audit.ActionDelete
			}

			// Calculate duration
			var duration time.Duration
			if start, ok := r.Context().Value(ContextKeyRequestStart).(time.Time); ok {
				duration = time.Since(start)
			}

			// Log the access
			event := audit.AuditEvent{
				EventAction:    action,
				EventOutcome:   outcome,
				ActorIPAddress: getClientIP(r),
				ActorUserAgent: r.UserAgent(),
				TargetType:     "endpoint",
				TargetID:       r.URL.Path,
				Message:        r.Method + " " + r.URL.Path,
				Component:      "http",
				Metadata: map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"query":       r.URL.RawQuery,
					"status_code": wrapped.statusCode,
					"duration_ms": duration.Milliseconds(),
				},
			}

			// Don't block on audit logging
			go func() {
				ctx := context.Background()
				// Copy context values
				if reqID, ok := r.Context().Value(ContextKeyRequestID).(string); ok {
					ctx = context.WithValue(ctx, audit.ContextKeyRequestID, reqID)
				}
				if corrID, ok := r.Context().Value(ContextKeyCorrelationID).(string); ok {
					ctx = context.WithValue(ctx, audit.ContextKeyCorrelationID, corrID)
				}
				logger.LogAccess(ctx, event)
			}()
		})
	}
}

// statusResponseWriter wraps http.ResponseWriter to capture status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// getClientIP extracts the real client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

// RecoverPanic recovers from panics and returns 500 error
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic (you'd want to use proper logging here)
				// Don't expose internal error details
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SecureJSONResponse ensures JSON responses have proper content type
func SecureJSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap to set content type for JSON responses
		wrapped := &jsonResponseWriter{ResponseWriter: w}
		next.ServeHTTP(wrapped, r)
	})
}

type jsonResponseWriter struct {
	http.ResponseWriter
	headerWritten bool
}

func (w *jsonResponseWriter) Write(b []byte) (int, error) {
	if !w.headerWritten {
		// Check if response looks like JSON
		if len(b) > 0 && (b[0] == '{' || b[0] == '[') {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		w.headerWritten = true
	}
	return w.ResponseWriter.Write(b)
}

func (w *jsonResponseWriter) WriteHeader(code int) {
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(code)
}

// Chain combines multiple middleware into a single middleware
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// DefaultSecurityChain returns the default security middleware chain
func DefaultSecurityChain() func(http.Handler) http.Handler {
	return Chain(
		RecoverPanic,
		EnhancedRequestID,
		EnhancedSecurityHeaders,
		RequestSanitizer,
		SecureJSONResponse,
	)
}
