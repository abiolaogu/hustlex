package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// TokenBlacklistRepository implements token blacklisting using Redis
type TokenBlacklistRepository struct {
	client *Client
}

// NewTokenBlacklistRepository creates a new token blacklist repository
func NewTokenBlacklistRepository(client *Client) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{
		client: client,
	}
}

// BlacklistToken adds a token to the blacklist with expiration
func (r *TokenBlacklistRepository) BlacklistToken(ctx context.Context, token string, expiresAt time.Time) error {
	// Hash the token to avoid storing raw JWTs (security best practice)
	tokenHash := hashToken(token)
	key := fmt.Sprintf("token:blacklist:%s", tokenHash)

	// Calculate TTL until token expiration
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Store with TTL matching token expiration
	// Value doesn't matter - we just check existence
	return r.client.Set(ctx, key, true, ttl)
}

// IsTokenBlacklisted checks if a token is blacklisted
func (r *TokenBlacklistRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	tokenHash := hashToken(token)
	key := fmt.Sprintf("token:blacklist:%s", tokenHash)

	exists, err := r.client.Exists(ctx, key)
	if err != nil {
		// On error, fail secure - treat as blacklisted
		return true, err
	}

	return exists, nil
}

// CleanupExpiredTokens is a no-op for Redis since TTL handles expiration automatically
func (r *TokenBlacklistRepository) CleanupExpiredTokens(ctx context.Context) error {
	// Redis automatically removes expired keys, so no cleanup needed
	return nil
}

// hashToken creates a SHA-256 hash of the token for storage
// This prevents storing raw JWTs in Redis, adding a layer of security
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
