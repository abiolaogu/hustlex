package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestInMemoryCSRFStore_GenerateToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	ctx := context.Background()

	token, err := store.GenerateToken(ctx, "user123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	// Verify token is stored
	if !store.ValidateToken(ctx, "user123", token) {
		t.Fatal("Generated token should be valid")
	}
}

func TestInMemoryCSRFStore_ValidateToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	ctx := context.Background()

	token, _ := store.GenerateToken(ctx, "user123")

	tests := []struct {
		name     string
		userID   string
		token    string
		expected bool
	}{
		{
			name:     "Valid token for correct user",
			userID:   "user123",
			token:    token,
			expected: true,
		},
		{
			name:     "Valid token for wrong user",
			userID:   "user456",
			token:    token,
			expected: false,
		},
		{
			name:     "Invalid token",
			userID:   "user123",
			token:    "invalid-token",
			expected: false,
		},
		{
			name:     "Empty token",
			userID:   "user123",
			token:    "",
			expected: false,
		},
		{
			name:     "Empty user ID",
			userID:   "",
			token:    token,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := store.ValidateToken(ctx, tt.userID, tt.token)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInMemoryCSRFStore_RevokeToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	ctx := context.Background()

	token, _ := store.GenerateToken(ctx, "user123")

	// Verify token is valid before revocation
	if !store.ValidateToken(ctx, "user123", token) {
		t.Fatal("Token should be valid before revocation")
	}

	// Revoke token
	err := store.RevokeToken(ctx, token)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify token is invalid after revocation
	if store.ValidateToken(ctx, "user123", token) {
		t.Fatal("Token should be invalid after revocation")
	}
}

func TestInMemoryCSRFStore_TokenExpiry(t *testing.T) {
	config := DefaultCSRFConfig()
	config.TokenLifetime = 100 * time.Millisecond
	store := NewInMemoryCSRFStore(config)
	ctx := context.Background()

	token, _ := store.GenerateToken(ctx, "user123")

	// Token should be valid immediately
	if !store.ValidateToken(ctx, "user123", token) {
		t.Fatal("Token should be valid immediately after creation")
	}

	// Wait for token to expire
	time.Sleep(150 * time.Millisecond)

	// Token should be invalid after expiry
	if store.ValidateToken(ctx, "user123", token) {
		t.Fatal("Token should be invalid after expiry")
	}
}

func TestCSRFProtection_SafeMethods(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false // For testing

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	safeMethods := []string{http.MethodGet, http.MethodHead, http.MethodOptions}

	for _, method := range safeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/test", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", rr.Code)
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_NoAuth(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	unsafeMethods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range unsafeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/test", nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Should pass because no user is authenticated
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status 200 for unauthenticated request, got %d", rr.Code)
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_WithAuth_NoToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	unsafeMethods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range unsafeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/test", nil)
			// Simulate authenticated user
			ctx := context.WithValue(req.Context(), ContextKeyUserID, "user123")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Should fail because no CSRF token provided
			if rr.Code != http.StatusForbidden {
				t.Errorf("Expected status 403, got %d", rr.Code)
			}

			if !strings.Contains(rr.Body.String(), "CSRF") {
				t.Errorf("Expected CSRF error message, got: %s", rr.Body.String())
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_WithAuth_ValidToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Generate valid token
	ctx := context.Background()
	token, _ := store.GenerateToken(ctx, "user123")

	unsafeMethods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range unsafeMethods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/test", nil)
			// Simulate authenticated user
			ctx := context.WithValue(req.Context(), ContextKeyUserID, "user123")
			req = req.WithContext(ctx)
			// Add CSRF token header
			req.Header.Set(csrfHeaderName, token)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Should succeed with valid token
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d. Body: %s", rr.Code, rr.Body.String())
			}

			// Should receive a new token in response
			newToken := rr.Header().Get(csrfHeaderName)
			if newToken == "" {
				t.Error("Expected new CSRF token in response header")
			}
		})
	}
}

func TestCSRFProtection_UnsafeMethods_WithAuth_InvalidToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	// Simulate authenticated user
	ctx := context.WithValue(req.Context(), ContextKeyUserID, "user123")
	req = req.WithContext(ctx)
	// Add invalid CSRF token
	req.Header.Set(csrfHeaderName, "invalid-token")

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should fail with invalid token
	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rr.Code)
	}
}

func TestCSRFProtection_SkipPaths(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	skipPaths := []string{
		"/health",
		"/metrics",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/webhooks/paystack",
	}

	for _, path := range skipPaths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			// Simulate authenticated user (but no CSRF token)
			ctx := context.WithValue(req.Context(), ContextKeyUserID, "user123")
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			// Should succeed because path is in skip list
			if rr.Code != http.StatusOK {
				t.Errorf("Expected status 200 for skipped path %s, got %d", path, rr.Code)
			}
		})
	}
}

func TestCSRFProtection_TokenRotation(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	handler := CSRFProtection(store, config)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Generate initial token
	ctx := context.Background()
	token1, _ := store.GenerateToken(ctx, "user123")

	// Make first request
	req1 := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	ctx1 := context.WithValue(req1.Context(), ContextKeyUserID, "user123")
	req1 = req1.WithContext(ctx1)
	req1.Header.Set(csrfHeaderName, token1)

	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)

	if rr1.Code != http.StatusOK {
		t.Fatalf("First request failed: %d", rr1.Code)
	}

	// Get new token from response
	token2 := rr1.Header().Get(csrfHeaderName)
	if token2 == "" {
		t.Fatal("Expected new token in response")
	}

	if token2 == token1 {
		t.Error("Expected token to rotate (new token different from old)")
	}

	// Verify new token is valid
	if !store.ValidateToken(ctx, "user123", token2) {
		t.Error("New token should be valid")
	}
}

func TestRevokeCSRFToken(t *testing.T) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	config := DefaultCSRFConfig()
	config.Secure = false

	// Generate token
	ctx := context.Background()
	token, _ := store.GenerateToken(ctx, "user123")

	// Create request with CSRF cookie
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  csrfCookieName,
		Value: token,
	})

	rr := httptest.NewRecorder()

	// Revoke token
	err := RevokeCSRFToken(rr, req, store, config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify token is revoked
	if store.ValidateToken(ctx, "user123", token) {
		t.Error("Token should be revoked")
	}

	// Verify cookie is deleted
	cookies := rr.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == csrfCookieName {
			found = true
			if cookie.MaxAge != -1 {
				t.Error("Expected cookie to be marked for deletion (MaxAge = -1)")
			}
			if cookie.Value != "" {
				t.Error("Expected cookie value to be empty")
			}
		}
	}
	if !found {
		t.Error("Expected CSRF cookie in response")
	}
}

func TestGetCSRFToken(t *testing.T) {
	token := "test-token"
	ctx := context.WithValue(context.Background(), ContextKeyCSRFToken, token)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req = req.WithContext(ctx)

	result := GetCSRFToken(req)
	if result != token {
		t.Errorf("Expected %s, got %s", token, result)
	}

	// Test with no token in context
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	result2 := GetCSRFToken(req2)
	if result2 != "" {
		t.Errorf("Expected empty string, got %s", result2)
	}
}

func TestIsSafeMethod(t *testing.T) {
	tests := []struct {
		method   string
		expected bool
	}{
		{http.MethodGet, true},
		{http.MethodHead, true},
		{http.MethodOptions, true},
		{http.MethodTrace, true},
		{http.MethodPost, false},
		{http.MethodPut, false},
		{http.MethodPatch, false},
		{http.MethodDelete, false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := isSafeMethod(tt.method)
			if result != tt.expected {
				t.Errorf("For method %s, expected %v, got %v", tt.method, tt.expected, result)
			}
		})
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.GenerateToken(ctx, "user123")
	}
}

func BenchmarkValidateToken(b *testing.B) {
	store := NewInMemoryCSRFStore(DefaultCSRFConfig())
	ctx := context.Background()
	token, _ := store.GenerateToken(ctx, "user123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.ValidateToken(ctx, "user123", token)
	}
}
