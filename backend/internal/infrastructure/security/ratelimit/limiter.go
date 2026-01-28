package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter interface for rate limiting implementations
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	AllowN(ctx context.Context, key string, n int) (bool, error)
	Reset(ctx context.Context, key string) error
	Remaining(ctx context.Context, key string) (int, error)
}

// RateConfig defines rate limiting configuration
type RateConfig struct {
	Requests int           // Max requests allowed
	Window   time.Duration // Time window
}

// Common rate limit configurations
var (
	RateLimitDefault     = RateConfig{Requests: 100, Window: time.Minute}
	RateLimitAuth        = RateConfig{Requests: 5, Window: time.Minute}       // Auth attempts
	RateLimitOTP         = RateConfig{Requests: 3, Window: 5 * time.Minute}   // OTP requests
	RateLimitTransaction = RateConfig{Requests: 10, Window: time.Minute}      // Transactions
	RateLimitAPIHeavy    = RateConfig{Requests: 20, Window: time.Minute}      // Heavy API calls
	RateLimitPIN         = RateConfig{Requests: 3, Window: 15 * time.Minute}  // PIN attempts
	RateLimitBVN         = RateConfig{Requests: 5, Window: 24 * time.Hour}    // BVN verification
)

// RedisRateLimiter implements RateLimiter using Redis with sliding window
type RedisRateLimiter struct {
	client *redis.Client
	config RateConfig
	prefix string
}

// NewRedisRateLimiter creates a new Redis-backed rate limiter
func NewRedisRateLimiter(client *redis.Client, config RateConfig, prefix string) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		config: config,
		prefix: prefix,
	}
}

func (r *RedisRateLimiter) key(identifier string) string {
	return fmt.Sprintf("ratelimit:%s:%s", r.prefix, identifier)
}

// Allow checks if a single request is allowed
func (r *RedisRateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
	return r.AllowN(ctx, identifier, 1)
}

// AllowN checks if N requests are allowed using sliding window algorithm
func (r *RedisRateLimiter) AllowN(ctx context.Context, identifier string, n int) (bool, error) {
	key := r.key(identifier)

	now := time.Now().UnixMilli()
	windowStart := now - r.config.Window.Milliseconds()

	pipe := r.client.Pipeline()

	// Remove expired entries
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))

	// Count current entries
	countCmd := pipe.ZCard(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := countCmd.Val()

	if int(count)+n > r.config.Requests {
		return false, nil
	}

	// Add new entries
	members := make([]redis.Z, n)
	for i := 0; i < n; i++ {
		members[i] = redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d-%d", now, i),
		}
	}

	pipe = r.client.Pipeline()
	pipe.ZAdd(ctx, key, members...)
	pipe.Expire(ctx, key, r.config.Window)

	_, err = pipe.Exec(ctx)
	return err == nil, err
}

// Reset clears the rate limit for an identifier
func (r *RedisRateLimiter) Reset(ctx context.Context, identifier string) error {
	return r.client.Del(ctx, r.key(identifier)).Err()
}

// Remaining returns the number of remaining requests allowed
func (r *RedisRateLimiter) Remaining(ctx context.Context, identifier string) (int, error) {
	key := r.key(identifier)

	now := time.Now().UnixMilli()
	windowStart := now - r.config.Window.Milliseconds()

	// Remove expired and count
	pipe := r.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	countCmd := pipe.ZCard(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	count := int(countCmd.Val())
	remaining := r.config.Requests - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// InMemoryRateLimiter implements RateLimiter using in-memory storage
// Suitable for single-instance deployments or testing
type InMemoryRateLimiter struct {
	mu      sync.RWMutex
	entries map[string][]time.Time
	config  RateConfig
	prefix  string
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter(config RateConfig, prefix string) *InMemoryRateLimiter {
	limiter := &InMemoryRateLimiter{
		entries: make(map[string][]time.Time),
		config:  config,
		prefix:  prefix,
	}
	// Start cleanup goroutine
	go limiter.cleanup()
	return limiter
}

func (r *InMemoryRateLimiter) key(identifier string) string {
	return fmt.Sprintf("%s:%s", r.prefix, identifier)
}

// cleanup periodically removes expired entries
func (r *InMemoryRateLimiter) cleanup() {
	ticker := time.NewTicker(r.config.Window)
	defer ticker.Stop()

	for range ticker.C {
		r.mu.Lock()
		cutoff := time.Now().Add(-r.config.Window)
		for key, times := range r.entries {
			var valid []time.Time
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(r.entries, key)
			} else {
				r.entries[key] = valid
			}
		}
		r.mu.Unlock()
	}
}

// Allow checks if a single request is allowed
func (r *InMemoryRateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
	return r.AllowN(ctx, identifier, 1)
}

// AllowN checks if N requests are allowed
func (r *InMemoryRateLimiter) AllowN(ctx context.Context, identifier string, n int) (bool, error) {
	key := r.key(identifier)
	now := time.Now()
	cutoff := now.Add(-r.config.Window)

	r.mu.Lock()
	defer r.mu.Unlock()

	// Get existing entries and filter expired
	times := r.entries[key]
	var valid []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	// Check if allowed
	if len(valid)+n > r.config.Requests {
		r.entries[key] = valid
		return false, nil
	}

	// Add new entries
	for i := 0; i < n; i++ {
		valid = append(valid, now)
	}
	r.entries[key] = valid

	return true, nil
}

// Reset clears the rate limit for an identifier
func (r *InMemoryRateLimiter) Reset(ctx context.Context, identifier string) error {
	key := r.key(identifier)

	r.mu.Lock()
	delete(r.entries, key)
	r.mu.Unlock()

	return nil
}

// Remaining returns the number of remaining requests allowed
func (r *InMemoryRateLimiter) Remaining(ctx context.Context, identifier string) (int, error) {
	key := r.key(identifier)
	now := time.Now()
	cutoff := now.Add(-r.config.Window)

	r.mu.RLock()
	times := r.entries[key]
	r.mu.RUnlock()

	count := 0
	for _, t := range times {
		if t.After(cutoff) {
			count++
		}
	}

	remaining := r.config.Requests - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// RateLimitMiddleware creates HTTP middleware for rate limiting
func RateLimitMiddleware(limiter RateLimiter, keyFunc func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)

			allowed, err := limiter.Allow(r.Context(), key)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				w.Header().Set("Retry-After", "60")
				w.Header().Set("X-RateLimit-Remaining", "0")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add remaining count to response header
			if remaining, err := limiter.Remaining(r.Context(), key); err == nil {
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IPKeyFunc extracts IP address from request for rate limiting
func IPKeyFunc(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// UserKeyFunc extracts user ID from context for rate limiting
func UserKeyFunc(userIDKey interface{}) func(r *http.Request) string {
	return func(r *http.Request) string {
		if userID, ok := r.Context().Value(userIDKey).(string); ok {
			return userID
		}
		return IPKeyFunc(r)
	}
}

// CompositeKeyFunc combines IP and endpoint for rate limiting
func CompositeKeyFunc(r *http.Request) string {
	return fmt.Sprintf("%s:%s:%s", IPKeyFunc(r), r.Method, r.URL.Path)
}
