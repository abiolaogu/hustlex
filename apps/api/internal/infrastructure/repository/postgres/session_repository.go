package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/shared/valueobject"
	"github.com/abiolaogu/hustlex/apps/api/internal/infrastructure/cache"
)

// SessionRepositoryImpl implements the SessionRepository interface using Redis/Cache
type SessionRepositoryImpl struct {
	cache cache.Cache
}

// NewSessionRepository creates a new instance of SessionRepositoryImpl
func NewSessionRepository(cache cache.Cache) repository.SessionRepository {
	return &SessionRepositoryImpl{cache: cache}
}

// StoreRefreshToken stores a refresh token for a user
func (r *SessionRepositoryImpl) StoreRefreshToken(ctx context.Context, userID valueobject.UserID, token string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	if err := r.cache.Set(ctx, key, token, expiry); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

// GetRefreshToken retrieves a stored refresh token for a user
func (r *SessionRepositoryImpl) GetRefreshToken(ctx context.Context, userID valueobject.UserID) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	token, err := r.cache.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return token, nil
}

// DeleteRefreshToken removes a refresh token (logout)
func (r *SessionRepositoryImpl) DeleteRefreshToken(ctx context.Context, userID valueobject.UserID) error {
	key := fmt.Sprintf("refresh_token:%s", userID.String())

	if err := r.cache.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

// CheckOTPRateLimit checks and increments OTP rate limit
func (r *SessionRepositoryImpl) CheckOTPRateLimit(ctx context.Context, phone string, maxRequests int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("otp_rate_limit:%s", phone)

	// Try to get current count
	countStr, err := r.cache.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, this is the first request
		if err := r.cache.Set(ctx, key, "1", window); err != nil {
			return false, fmt.Errorf("failed to set rate limit: %w", err)
		}
		return true, nil
	}

	// Parse current count
	var count int
	if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil {
		return false, fmt.Errorf("failed to parse rate limit count: %w", err)
	}

	// Check if limit exceeded
	if count >= maxRequests {
		return false, nil
	}

	// Increment count
	newCount := count + 1
	if err := r.cache.Set(ctx, key, fmt.Sprintf("%d", newCount), window); err != nil {
		return false, fmt.Errorf("failed to update rate limit: %w", err)
	}

	return true, nil
}
