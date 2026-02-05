package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// SessionRepository is a Redis implementation of the SessionRepository interface
type SessionRepository struct {
	client *redis.Client
}

// NewSessionRepository creates a new Redis session repository
func NewSessionRepository(client *redis.Client) repository.SessionRepository {
	return &SessionRepository{client: client}
}

// StoreRefreshToken stores a refresh token with expiry
func (r *SessionRepository) StoreRefreshToken(
	ctx context.Context,
	userID valueobject.UserID,
	token string,
	expiry time.Duration,
) error {
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	return r.client.Set(ctx, key, token, expiry).Err()
}

// GetRefreshToken retrieves a stored refresh token
func (r *SessionRepository) GetRefreshToken(
	ctx context.Context,
	userID valueobject.UserID,
) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	token, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	}
	return token, err
}

// DeleteRefreshToken removes a refresh token (logout)
func (r *SessionRepository) DeleteRefreshToken(
	ctx context.Context,
	userID valueobject.UserID,
) error {
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	return r.client.Del(ctx, key).Err()
}

// CheckOTPRateLimit checks and increments OTP rate limit
func (r *SessionRepository) CheckOTPRateLimit(
	ctx context.Context,
	phone string,
	maxRequests int,
	window time.Duration,
) (bool, error) {
	key := fmt.Sprintf("otp_rate_limit:%s", phone)

	// Get current count
	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		// First request in this window
		count = 0
	} else if err != nil {
		return false, err
	}

	// Check if limit exceeded
	if count >= maxRequests {
		return false, nil // Rate limit exceeded
	}

	// Increment counter
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil // Within rate limit
}
