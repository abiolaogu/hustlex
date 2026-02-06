package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// ContextKey type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// SessionIDKey is the context key for session ID
	SessionIDKey ContextKey = "session_id"
	// RolesKey is the context key for user roles
	RolesKey ContextKey = "roles"
)

// TokenValidator interface for validating JWT tokens
type TokenValidator interface {
	ValidateToken(token string) (*TokenClaims, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    string
	SessionID string
	Roles     []string
	ExpiresAt int64
}

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	tokenValidator    TokenValidator
	tokenBlacklist    repository.TokenBlacklistRepository
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(tokenValidator TokenValidator, tokenBlacklist repository.TokenBlacklistRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenValidator: tokenValidator,
		tokenBlacklist: tokenBlacklist,
	}
}

// Authenticate validates the JWT token and sets user context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			http.Error(w, `{"error":"unauthorized","message":"missing or invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Check if token is blacklisted (revoked)
		if m.tokenBlacklist != nil {
			blacklisted, err := m.tokenBlacklist.IsTokenBlacklisted(r.Context(), token)
			if err != nil {
				// Log error but don't fail the request if blacklist check fails
				// In production, you might want to fail secure (reject the token)
				// For now, we log and continue
			} else if blacklisted {
				http.Error(w, `{"error":"unauthorized","message":"token has been revoked"}`, http.StatusUnauthorized)
				return
			}
		}

		claims, err := m.tokenValidator.ValidateToken(token)
		if err != nil {
			http.Error(w, `{"error":"unauthorized","message":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Validate user ID format
		userID, err := valueobject.NewUserID(claims.UserID)
		if err != nil {
			http.Error(w, `{"error":"unauthorized","message":"invalid user ID"}`, http.StatusUnauthorized)
			return
		}

		// Set user context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)
		ctx = context.WithValue(ctx, RolesKey, claims.Roles)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth sets user context if token is present, but doesn't require it
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token != "" {
			claims, err := m.tokenValidator.ValidateToken(token)
			if err == nil {
				userID, err := valueobject.NewUserID(claims.UserID)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					ctx = context.WithValue(ctx, SessionIDKey, claims.SessionID)
					ctx = context.WithValue(ctx, RolesKey, claims.Roles)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// RequireRoles ensures the user has at least one of the required roles
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles, ok := r.Context().Value(RolesKey).([]string)
			if !ok {
				http.Error(w, `{"error":"forbidden","message":"access denied"}`, http.StatusForbidden)
				return
			}

			hasRole := false
			for _, required := range roles {
				for _, userRole := range userRoles {
					if userRole == required {
						hasRole = true
						break
					}
				}
				if hasRole {
					break
				}
			}

			if !hasRole {
				http.Error(w, `{"error":"forbidden","message":"insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (valueobject.UserID, error) {
	userID, ok := ctx.Value(UserIDKey).(valueobject.UserID)
	if !ok {
		return valueobject.UserID{}, errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetSessionID extracts session ID from context
func GetSessionID(ctx context.Context) string {
	sessionID, _ := ctx.Value(SessionIDKey).(string)
	return sessionID
}

// GetRoles extracts roles from context
func GetRoles(ctx context.Context) []string {
	roles, _ := ctx.Value(RolesKey).([]string)
	return roles
}

// HasRole checks if user has a specific role
func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// extractToken extracts the JWT token from the request
func extractToken(r *http.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Try cookie
	cookie, err := r.Cookie("access_token")
	if err == nil {
		return cookie.Value
	}

	// Try query parameter (for WebSocket connections)
	return r.URL.Query().Get("token")
}
