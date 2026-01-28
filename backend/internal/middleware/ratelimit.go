package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// RateLimiter implements rate limiting using Redis
type RateLimiter struct {
	redis    *redis.Client
	requests int           // Max requests
	window   time.Duration // Time window
	prefix   string        // Key prefix
}

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Redis    *redis.Client
	Requests int           // Default: 100
	Window   time.Duration // Default: 1 minute
	Prefix   string        // Default: "ratelimit"
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	if cfg.Requests == 0 {
		cfg.Requests = 100
	}
	if cfg.Window == 0 {
		cfg.Window = time.Minute
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "ratelimit"
	}

	return &RateLimiter{
		redis:    cfg.Redis,
		requests: cfg.Requests,
		window:   cfg.Window,
		prefix:   cfg.Prefix,
	}
}

// Limit returns a Fiber middleware for rate limiting
func (rl *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get identifier (IP or user ID)
		identifier := c.IP()
		if userID, ok := c.Locals("userID").(string); ok {
			identifier = userID
		}

		key := fmt.Sprintf("%s:%s", rl.prefix, identifier)

		// Check and increment using Redis
		ctx := context.Background()

		// Use sliding window algorithm
		now := time.Now().UnixMilli()
		windowStart := now - rl.window.Milliseconds()

		// Remove old entries
		rl.redis.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

		// Count current requests
		count, err := rl.redis.ZCard(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow the request (fail open)
			return c.Next()
		}

		if count >= int64(rl.requests) {
			// Get time until reset
			oldestEntry, err := rl.redis.ZRange(ctx, key, 0, 0).Result()
			if err == nil && len(oldestEntry) > 0 {
				oldestTime, _ := strconv.ParseInt(oldestEntry[0], 10, 64)
				resetTime := time.UnixMilli(oldestTime + rl.window.Milliseconds())

				c.Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
				c.Set("X-RateLimit-Remaining", "0")
				c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
				c.Set("Retry-After", strconv.FormatInt(int64(time.Until(resetTime).Seconds()), 10))
			}

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
			})
		}

		// Add current request
		rl.redis.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: strconv.FormatInt(now, 10),
		})
		rl.redis.Expire(ctx, key, rl.window)

		// Set headers
		remaining := rl.requests - int(count) - 1
		c.Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.window).Unix(), 10))

		return c.Next()
	}
}

// StrictLimit creates a stricter rate limiter for sensitive endpoints
func StrictLimit(redis *redis.Client) fiber.Handler {
	limiter := NewRateLimiter(RateLimiterConfig{
		Redis:    redis,
		Requests: 10,            // 10 requests
		Window:   time.Minute,   // per minute
		Prefix:   "ratelimit:strict",
	})
	return limiter.Limit()
}

// OTPLimit creates a rate limiter for OTP requests
func OTPLimit(redis *redis.Client) fiber.Handler {
	limiter := NewRateLimiter(RateLimiterConfig{
		Redis:    redis,
		Requests: 5,                 // 5 OTP requests
		Window:   15 * time.Minute,  // per 15 minutes
		Prefix:   "ratelimit:otp",
	})
	return limiter.Limit()
}

// LoginLimit creates a rate limiter for login attempts
func LoginLimit(redis *redis.Client) fiber.Handler {
	limiter := NewRateLimiter(RateLimiterConfig{
		Redis:    redis,
		Requests: 5,              // 5 login attempts
		Window:   5 * time.Minute, // per 5 minutes
		Prefix:   "ratelimit:login",
	})
	return limiter.Limit()
}

// TransactionLimit creates a rate limiter for financial transactions
func TransactionLimit(redis *redis.Client) fiber.Handler {
	limiter := NewRateLimiter(RateLimiterConfig{
		Redis:    redis,
		Requests: 30,             // 30 transactions
		Window:   time.Hour,      // per hour
		Prefix:   "ratelimit:transaction",
	})
	return limiter.Limit()
}
