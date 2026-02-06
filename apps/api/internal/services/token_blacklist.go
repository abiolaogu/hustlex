package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklistService handles token revocation and blacklisting
type TokenBlacklistService struct {
	redis *redis.Client
}

// NewTokenBlacklistService creates a new token blacklist service
func NewTokenBlacklistService(redis *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{
		redis: redis,
	}
}

// BlacklistToken adds a token to the blacklist with expiration
// The token is stored with a hash to avoid storing the full token
// TTL is set to the remaining lifetime of the token
func (s *TokenBlacklistService) BlacklistToken(ctx context.Context, token string, expiresAt time.Time) error {
	// Hash the token to avoid storing full token in Redis
	tokenHash := hashToken(token)
	key := fmt.Sprintf("token:blacklist:%s", tokenHash)

	// Calculate TTL - only store until token would expire anyway
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Store with TTL - value is the blacklist timestamp
	return s.redis.Set(ctx, key, time.Now().Unix(), ttl).Err()
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (s *TokenBlacklistService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	tokenHash := hashToken(token)
	key := fmt.Sprintf("token:blacklist:%s", tokenHash)

	// Check if key exists
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}

	return result > 0, nil
}

// BlacklistAllUserTokens blacklists all tokens for a user
// This is useful when password/PIN changes or account suspension
// Since we can't enumerate all user tokens, we use a user-level blacklist marker
func (s *TokenBlacklistService) BlacklistAllUserTokens(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:tokens:revoked:%s", userID)

	// Store revocation timestamp - tokens issued before this time are invalid
	// Set a long TTL (e.g., 30 days) to cover all possible token lifetimes
	return s.redis.Set(ctx, key, time.Now().Unix(), 30*24*time.Hour).Err()
}

// IsUserTokenRevoked checks if all tokens for a user have been revoked
// Returns the revocation timestamp if revoked
func (s *TokenBlacklistService) IsUserTokenRevoked(ctx context.Context, userID string, tokenIssuedAt time.Time) (bool, error) {
	key := fmt.Sprintf("user:tokens:revoked:%s", userID)

	result, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		// No revocation marker exists
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check user revocation: %w", err)
	}

	// Parse revocation timestamp
	var revokedAt int64
	fmt.Sscanf(result, "%d", &revokedAt)

	// Token is revoked if it was issued before the revocation time
	return tokenIssuedAt.Before(time.Unix(revokedAt, 0)), nil
}

// ClearUserTokenRevocation removes the user-level token revocation
// This would be used if re-enabling an account
func (s *TokenBlacklistService) ClearUserTokenRevocation(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:tokens:revoked:%s", userID)
	return s.redis.Del(ctx, key).Err()
}

// GetBlacklistStats returns statistics about the blacklist
// Useful for monitoring and debugging
func (s *TokenBlacklistService) GetBlacklistStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count blacklisted tokens
	iter := s.redis.Scan(ctx, 0, "token:blacklist:*", 100).Iterator()
	count := int64(0)
	for iter.Next(ctx) {
		count++
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	stats["blacklisted_tokens"] = count

	// Count user revocations
	iter = s.redis.Scan(ctx, 0, "user:tokens:revoked:*", 100).Iterator()
	count = 0
	for iter.Next(ctx) {
		count++
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	stats["user_revocations"] = count

	return stats, nil
}

// hashToken creates a SHA-256 hash of the token
// This prevents storing the actual token in Redis (defense in depth)
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
