package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// Client wraps the Redis client with application-specific methods
type Client struct {
	rdb *redis.Client
}

// NewClient creates a new Redis client
func NewClient(config Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}

// Cache operations

// Set stores a value with expiration
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetNX sets a value only if the key doesn't exist (for distributed locks)
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return c.rdb.SetNX(ctx, key, data, expiration).Result()
}

// Expire sets expiration on a key
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.rdb.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// OTP operations

// StoreOTP stores an OTP with expiration
func (c *Client) StoreOTP(ctx context.Context, phone, otp string, expiration time.Duration) error {
	key := fmt.Sprintf("otp:%s", phone)
	data := map[string]interface{}{
		"otp":        otp,
		"attempts":   0,
		"created_at": time.Now().UTC(),
	}
	return c.Set(ctx, key, data, expiration)
}

// GetOTP retrieves an OTP
func (c *Client) GetOTP(ctx context.Context, phone string) (string, int, error) {
	key := fmt.Sprintf("otp:%s", phone)
	var data struct {
		OTP      string `json:"otp"`
		Attempts int    `json:"attempts"`
	}

	if err := c.Get(ctx, key, &data); err != nil {
		return "", 0, err
	}

	return data.OTP, data.Attempts, nil
}

// IncrementOTPAttempts increments the attempt counter
func (c *Client) IncrementOTPAttempts(ctx context.Context, phone string) error {
	key := fmt.Sprintf("otp:%s", phone)
	return c.rdb.HIncrBy(ctx, key, "attempts", 1).Err()
}

// DeleteOTP removes an OTP
func (c *Client) DeleteOTP(ctx context.Context, phone string) error {
	key := fmt.Sprintf("otp:%s", phone)
	return c.Delete(ctx, key)
}

// Session operations

// StoreSession stores a session
func (c *Client) StoreSession(ctx context.Context, sessionID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Set(ctx, key, data, expiration)
}

// GetSession retrieves a session
func (c *Client) GetSession(ctx context.Context, sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Get(ctx, key, dest)
}

// RefreshSession extends session expiration
func (c *Client) RefreshSession(ctx context.Context, sessionID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Expire(ctx, key, expiration)
}

// DeleteSession removes a session
func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Delete(ctx, key)
}

// Rate limiting

// RateLimit checks and increments rate limit counter
func (c *Client) RateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, int, error) {
	rateLimitKey := fmt.Sprintf("ratelimit:%s", key)

	pipe := c.rdb.Pipeline()
	incr := pipe.Incr(ctx, rateLimitKey)
	pipe.Expire(ctx, rateLimitKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := int(incr.Val())
	return count <= limit, count, nil
}

// Distributed locking

// AcquireLock attempts to acquire a distributed lock
func (c *Client) AcquireLock(ctx context.Context, lockName string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("lock:%s", lockName)
	return c.SetNX(ctx, key, time.Now().UTC(), ttl)
}

// ReleaseLock releases a distributed lock
func (c *Client) ReleaseLock(ctx context.Context, lockName string) error {
	key := fmt.Sprintf("lock:%s", lockName)
	return c.Delete(ctx, key)
}

// Pub/Sub

// Publish publishes a message to a channel
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.rdb.Publish(ctx, channel, data).Err()
}

// Subscribe subscribes to a channel
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.rdb.Subscribe(ctx, channels...)
}

// Cache key helpers

// UserCacheKey generates a cache key for user data
func UserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

// WalletCacheKey generates a cache key for wallet data
func WalletCacheKey(walletID string) string {
	return fmt.Sprintf("wallet:%s", walletID)
}

// CreditScoreCacheKey generates a cache key for credit score
func CreditScoreCacheKey(userID string) string {
	return fmt.Sprintf("credit:%s", userID)
}

// Error definitions
var (
	ErrCacheMiss = errors.New("cache miss")
)
