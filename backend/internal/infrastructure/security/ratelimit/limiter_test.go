package ratelimit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInMemoryRateLimiter_Allow(t *testing.T) {
	config := RateConfig{Requests: 5, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(ctx, "user-1")
		if err != nil {
			t.Fatalf("Allow() error: %v", err)
		}
		if !allowed {
			t.Errorf("Allow() request %d should be allowed", i+1)
		}
	}

	// 6th request should be denied
	allowed, err := limiter.Allow(ctx, "user-1")
	if err != nil {
		t.Fatalf("Allow() error: %v", err)
	}
	if allowed {
		t.Error("Allow() 6th request should be denied")
	}
}

func TestInMemoryRateLimiter_AllowN(t *testing.T) {
	config := RateConfig{Requests: 10, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// Request 5 at once
	allowed, err := limiter.AllowN(ctx, "user-1", 5)
	if err != nil {
		t.Fatalf("AllowN() error: %v", err)
	}
	if !allowed {
		t.Error("AllowN(5) should be allowed")
	}

	// Request 5 more
	allowed, err = limiter.AllowN(ctx, "user-1", 5)
	if err != nil {
		t.Fatalf("AllowN() error: %v", err)
	}
	if !allowed {
		t.Error("AllowN(5) second batch should be allowed")
	}

	// Request 1 more should be denied
	allowed, err = limiter.AllowN(ctx, "user-1", 1)
	if err != nil {
		t.Fatalf("AllowN() error: %v", err)
	}
	if allowed {
		t.Error("AllowN(1) should be denied after limit")
	}
}

func TestInMemoryRateLimiter_DifferentKeys(t *testing.T) {
	config := RateConfig{Requests: 2, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// User 1 uses their limit
	limiter.Allow(ctx, "user-1")
	limiter.Allow(ctx, "user-1")

	// User 1 should be denied
	allowed, _ := limiter.Allow(ctx, "user-1")
	if allowed {
		t.Error("User 1 should be rate limited")
	}

	// User 2 should still be allowed
	allowed, _ = limiter.Allow(ctx, "user-2")
	if !allowed {
		t.Error("User 2 should not be affected by User 1's limit")
	}
}

func TestInMemoryRateLimiter_Reset(t *testing.T) {
	config := RateConfig{Requests: 2, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// Use up the limit
	limiter.Allow(ctx, "user-1")
	limiter.Allow(ctx, "user-1")

	// Should be denied
	allowed, _ := limiter.Allow(ctx, "user-1")
	if allowed {
		t.Error("Should be rate limited")
	}

	// Reset
	err := limiter.Reset(ctx, "user-1")
	if err != nil {
		t.Fatalf("Reset() error: %v", err)
	}

	// Should be allowed again
	allowed, _ = limiter.Allow(ctx, "user-1")
	if !allowed {
		t.Error("Should be allowed after reset")
	}
}

func TestInMemoryRateLimiter_Remaining(t *testing.T) {
	config := RateConfig{Requests: 5, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// Initial remaining
	remaining, err := limiter.Remaining(ctx, "user-1")
	if err != nil {
		t.Fatalf("Remaining() error: %v", err)
	}
	if remaining != 5 {
		t.Errorf("Remaining() = %d, want 5", remaining)
	}

	// After 2 requests
	limiter.Allow(ctx, "user-1")
	limiter.Allow(ctx, "user-1")

	remaining, _ = limiter.Remaining(ctx, "user-1")
	if remaining != 3 {
		t.Errorf("Remaining() after 2 requests = %d, want 3", remaining)
	}
}

func TestInMemoryRateLimiter_WindowExpiry(t *testing.T) {
	config := RateConfig{Requests: 2, Window: 100 * time.Millisecond}
	limiter := NewInMemoryRateLimiter(config, "test")
	ctx := context.Background()

	// Use up limit
	limiter.Allow(ctx, "user-1")
	limiter.Allow(ctx, "user-1")

	// Should be denied
	allowed, _ := limiter.Allow(ctx, "user-1")
	if allowed {
		t.Error("Should be rate limited")
	}

	// Wait for window to expire
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	allowed, _ = limiter.Allow(ctx, "user-1")
	if !allowed {
		t.Error("Should be allowed after window expires")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	config := RateConfig{Requests: 2, Window: time.Minute}
	limiter := NewInMemoryRateLimiter(config, "test")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimitMiddleware(limiter, IPKeyFunc)
	wrappedHandler := middleware(handler)

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		rr := httptest.NewRecorder()

		wrappedHandler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Request %d status = %d, want 200", i+1, rr.Code)
		}
	}

	// 3rd request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Rate limited request status = %d, want 429", rr.Code)
	}

	retryAfter := rr.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Error("Rate limited response should have Retry-After header")
	}
}

func TestIPKeyFunc(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		xri        string
		want       string
	}{
		{
			name:       "X-Forwarded-For",
			remoteAddr: "127.0.0.1:12345",
			xff:        "203.0.113.195",
			want:       "203.0.113.195",
		},
		{
			name:       "X-Real-IP",
			remoteAddr: "127.0.0.1:12345",
			xri:        "198.51.100.178",
			want:       "198.51.100.178",
		},
		{
			name:       "RemoteAddr fallback",
			remoteAddr: "192.168.1.1:12345",
			want:       "192.168.1.1:12345",
		},
		{
			name:       "XFF takes precedence",
			remoteAddr: "127.0.0.1:12345",
			xff:        "203.0.113.195",
			xri:        "198.51.100.178",
			want:       "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xri != "" {
				req.Header.Set("X-Real-IP", tt.xri)
			}

			got := IPKeyFunc(req)
			if got != tt.want {
				t.Errorf("IPKeyFunc() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCompositeKeyFunc(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/wallet/transfer", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	got := CompositeKeyFunc(req)
	want := "192.168.1.1:12345:POST:/api/wallet/transfer"

	if got != want {
		t.Errorf("CompositeKeyFunc() = %q, want %q", got, want)
	}
}

func TestUserKeyFunc(t *testing.T) {
	type ctxKey string
	const userIDKey ctxKey = "user_id"

	keyFunc := UserKeyFunc(userIDKey)

	t.Run("with user ID in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		ctx := context.WithValue(req.Context(), userIDKey, "user-123")
		req = req.WithContext(ctx)

		got := keyFunc(req)
		if got != "user-123" {
			t.Errorf("UserKeyFunc() = %q, want user-123", got)
		}
	})

	t.Run("without user ID in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:12345"

		got := keyFunc(req)
		if got != "192.168.1.1:12345" {
			t.Errorf("UserKeyFunc() = %q, want IP address", got)
		}
	})
}

func TestRateConfigs(t *testing.T) {
	tests := []struct {
		name   string
		config RateConfig
	}{
		{"Default", RateLimitDefault},
		{"Auth", RateLimitAuth},
		{"OTP", RateLimitOTP},
		{"Transaction", RateLimitTransaction},
		{"APIHeavy", RateLimitAPIHeavy},
		{"PIN", RateLimitPIN},
		{"BVN", RateLimitBVN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Requests <= 0 {
				t.Error("Requests should be positive")
			}
			if tt.config.Window <= 0 {
				t.Error("Window should be positive")
			}
		})
	}
}
