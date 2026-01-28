package middleware

import (
	"github.com/gofiber/fiber/v2"
	"hustlex/internal/config"
)

// SecurityHeaders adds security headers to all responses
// Addresses OWASP security requirements for web applications
func SecurityHeaders(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-Frame-Options: Prevents clickjacking attacks
		// DENY = page cannot be displayed in a frame
		c.Set("X-Frame-Options", "DENY")

		// X-Content-Type-Options: Prevents MIME type sniffing
		// nosniff = browser must use declared Content-Type
		c.Set("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection: Legacy XSS filter (for older browsers)
		// 1; mode=block = enable XSS filter, block page if attack detected
		c.Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy: Controls how much referrer information is sent
		// strict-origin-when-cross-origin = send full URL for same-origin, origin only for cross-origin
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content-Security-Policy: Prevents XSS, clickjacking, and other code injection attacks
		// This is a restrictive policy suitable for an API
		csp := "default-src 'none'; frame-ancestors 'none'; form-action 'none'"
		c.Set("Content-Security-Policy", csp)

		// Permissions-Policy: Controls browser features (formerly Feature-Policy)
		// Disable access to sensitive features
		c.Set("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")

		// Cache-Control: Prevent caching of sensitive API responses
		// Only apply to authenticated endpoints
		if c.Get("Authorization") != "" {
			c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
			c.Set("Pragma", "no-cache")
			c.Set("Expires", "0")
		}

		// In production, enforce HTTPS
		if cfg.IsProduction() {
			// Strict-Transport-Security: Enforce HTTPS
			// max-age=31536000 = 1 year
			// includeSubDomains = apply to all subdomains
			// preload = allow inclusion in browser HSTS preload lists
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		return c.Next()
	}
}

// SecureAPIHeaders adds headers specifically for API responses
func SecureAPIHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ensure responses are treated as API data, not rendered
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")

		// API-specific CSP
		c.Set("Content-Security-Policy", "default-src 'none'")

		return c.Next()
	}
}

// RequestIDMiddleware adds a unique request ID for tracing
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for existing request ID (from load balancer, etc.)
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = c.Get("X-Correlation-ID")
		}
		if requestID == "" {
			// Generate a new request ID
			requestID = fiber.MIMETextHTML // This will be replaced with actual UUID generation
		}

		// Set request ID in response headers for debugging
		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		return c.Next()
	}
}

// CORSConfig returns a configured CORS middleware
// This should be used instead of allowing all origins
func CORSConfig(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		allowedOrigins := cfg.Server.AllowOrigins

		// Check if origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			c.Set("Access-Control-Allow-Origin", origin)
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-Request-ID")
			c.Set("Access-Control-Max-Age", "86400") // 24 hours
		}

		// Handle preflight requests
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin, allowedOrigins string) bool {
	if origin == "" {
		return false
	}

	// Parse comma-separated origins
	for _, allowed := range splitOrigins(allowedOrigins) {
		if allowed == origin {
			return true
		}
	}

	return false
}

// splitOrigins splits comma-separated origins into a slice
func splitOrigins(origins string) []string {
	var result []string
	current := ""
	for _, char := range origins {
		if char == ',' {
			if trimmed := trimSpace(current); trimmed != "" {
				result = append(result, trimmed)
			}
			current = ""
		} else {
			current += string(char)
		}
	}
	if trimmed := trimSpace(current); trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

// trimSpace removes leading and trailing whitespace
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// NoSniff adds X-Content-Type-Options header
func NoSniff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		return c.Next()
	}
}

// FrameGuard adds X-Frame-Options header
func FrameGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Frame-Options", "DENY")
		return c.Next()
	}
}
