package middleware

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// CSRF protection implementation following OWASP recommendations
// https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html

const (
	// ContextKeyCSRFToken is the context key for CSRF tokens
	ContextKeyCSRFToken contextKey = "csrf_token"

	// Default token lifetime (4 hours)
	defaultCSRFTokenLifetime = 4 * time.Hour

	// CSRF header name
	csrfHeaderName = "X-CSRF-Token"

	// CSRF cookie name
	csrfCookieName = "csrf_token"
)

// CSRFConfig holds CSRF protection configuration
type CSRFConfig struct {
	// TokenLifetime is how long a CSRF token is valid (default: 4 hours)
	TokenLifetime time.Duration

	// CookiePath is the path for the CSRF cookie (default: /)
	CookiePath string

	// CookieDomain is the domain for the CSRF cookie (optional)
	CookieDomain string

	// Secure sets the Secure flag on the cookie (should be true in production)
	Secure bool

	// SameSite sets the SameSite attribute (default: Strict)
	SameSite http.SameSite

	// SkipPaths are paths that should skip CSRF validation
	SkipPaths []string
}

// DefaultCSRFConfig returns the default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLifetime: defaultCSRFTokenLifetime,
		CookiePath:    "/",
		CookieDomain:  "",
		Secure:        true,
		SameSite:      http.SameSiteStrictMode,
		SkipPaths: []string{
			"/health",
			"/metrics",
			"/api/v1/auth/login",     // Login doesn't need CSRF (uses credentials)
			"/api/v1/auth/register",  // Register doesn't need CSRF
			"/api/v1/webhooks/",      // Webhooks use signature verification
		},
	}
}

// CSRFTokenStore manages CSRF tokens
type CSRFTokenStore interface {
	// GenerateToken generates a new CSRF token for the given user
	GenerateToken(ctx context.Context, userID string) (string, error)

	// ValidateToken validates a CSRF token for the given user
	ValidateToken(ctx context.Context, userID string, token string) bool

	// RevokeToken revokes a CSRF token
	RevokeToken(ctx context.Context, token string) error
}

// InMemoryCSRFStore is an in-memory implementation of CSRFTokenStore
type InMemoryCSRFStore struct {
	mu     sync.RWMutex
	tokens map[string]csrfToken
	config CSRFConfig
}

type csrfToken struct {
	userID    string
	token     string
	issuedAt  time.Time
	expiresAt time.Time
}

// NewInMemoryCSRFStore creates a new in-memory CSRF token store
func NewInMemoryCSRFStore(config CSRFConfig) *InMemoryCSRFStore {
	store := &InMemoryCSRFStore{
		tokens: make(map[string]csrfToken),
		config: config,
	}

	// Start cleanup goroutine
	go store.cleanupExpiredTokens()

	return store
}

// GenerateToken generates a new CSRF token
func (s *InMemoryCSRFStore) GenerateToken(ctx context.Context, userID string) (string, error) {
	token := generateSecureToken()

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.tokens[token] = csrfToken{
		userID:    userID,
		token:     token,
		issuedAt:  now,
		expiresAt: now.Add(s.config.TokenLifetime),
	}

	return token, nil
}

// ValidateToken validates a CSRF token
func (s *InMemoryCSRFStore) ValidateToken(ctx context.Context, userID string, token string) bool {
	if token == "" || userID == "" {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	storedToken, exists := s.tokens[token]
	if !exists {
		return false
	}

	// Check if token is expired
	if time.Now().After(storedToken.expiresAt) {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(storedToken.userID), []byte(userID)) == 1
}

// RevokeToken revokes a CSRF token
func (s *InMemoryCSRFStore) RevokeToken(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, token)
	return nil
}

// cleanupExpiredTokens periodically removes expired tokens
func (s *InMemoryCSRFStore) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, t := range s.tokens {
			if now.After(t.expiresAt) {
				delete(s.tokens, token)
			}
		}
		s.mu.Unlock()
	}
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() string {
	// Use UUID v4 for cryptographically secure random token
	return base64.URLEncoding.EncodeToString([]byte(uuid.NewString()))
}

// CSRFProtection returns a middleware that provides CSRF protection
func CSRFProtection(store CSRFTokenStore, config CSRFConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF validation for configured paths
			for _, skipPath := range config.SkipPaths {
				if strings.HasPrefix(r.URL.Path, skipPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Safe methods don't require CSRF protection
			if isSafeMethod(r.Method) {
				// For GET requests, generate and set a new CSRF token if user is authenticated
				if r.Method == http.MethodGet {
					userID := getUserIDFromContext(r.Context())
					if userID != "" {
						if err := issueCSRFToken(w, store, userID, config); err != nil {
							// Log error but don't fail the request
							// In production, you'd use proper logging
						}
					}
				}
				next.ServeHTTP(w, r)
				return
			}

			// For state-changing methods (POST, PUT, DELETE, PATCH), validate CSRF token
			userID := getUserIDFromContext(r.Context())
			if userID == "" {
				// Not authenticated - CSRF protection doesn't apply
				next.ServeHTTP(w, r)
				return
			}

			// Extract CSRF token from header
			token := r.Header.Get(csrfHeaderName)
			if token == "" {
				// Try form field as fallback
				token = r.FormValue("csrf_token")
			}

			// Validate token
			if !store.ValidateToken(r.Context(), userID, token) {
				http.Error(w, `{"error":"Invalid or missing CSRF token","code":"CSRF_TOKEN_INVALID"}`,
					http.StatusForbidden)
				return
			}

			// Token is valid, add to context for handlers to use if needed
			ctx := context.WithValue(r.Context(), ContextKeyCSRFToken, token)

			// Generate a new token for the next request (token rotation)
			if err := issueCSRFToken(w, store, userID, config); err != nil {
				// Log error but don't fail the request
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isSafeMethod returns true if the HTTP method is considered safe (no state change)
func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// issueCSRFToken generates a new CSRF token and sets it in a cookie
func issueCSRFToken(w http.ResponseWriter, store CSRFTokenStore, userID string, config CSRFConfig) error {
	token, err := store.GenerateToken(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	// Set CSRF token in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     config.CookiePath,
		Domain:   config.CookieDomain,
		Expires:  time.Now().Add(config.TokenLifetime),
		MaxAge:   int(config.TokenLifetime.Seconds()),
		Secure:   config.Secure,
		HttpOnly: false, // Must be false so JavaScript can read it
		SameSite: config.SameSite,
	})

	// Also set in response header for single-page applications
	w.Header().Set(csrfHeaderName, token)

	return nil
}

// getUserIDFromContext extracts the user ID from the request context
// This assumes the auth middleware has already set the user ID
func getUserIDFromContext(ctx context.Context) string {
	// Try to get from ContextKeyUserID (set by auth middleware)
	if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetCSRFToken extracts the CSRF token from the request
// This is a helper function for handlers that need to access the validated token
func GetCSRFToken(r *http.Request) string {
	if token, ok := r.Context().Value(ContextKeyCSRFToken).(string); ok {
		return token
	}
	return ""
}

// RevokeCSRFToken revokes the CSRF token for the current user
// This should be called on logout or when the user's session is invalidated
func RevokeCSRFToken(w http.ResponseWriter, r *http.Request, store CSRFTokenStore, config CSRFConfig) error {
	// Get token from cookie
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil {
		return nil // No token to revoke
	}

	// Revoke from store
	if err := store.RevokeToken(r.Context(), cookie.Value); err != nil {
		return err
	}

	// Delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookieName,
		Value:    "",
		Path:     config.CookiePath,
		Domain:   config.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   config.Secure,
		HttpOnly: false,
		SameSite: config.SameSite,
	})

	return nil
}
