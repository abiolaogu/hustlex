package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, mr
}

func TestNewTokenBlacklistService(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	assert.NotNil(t, service)
	assert.NotNil(t, service.redis)
}

func TestBlacklistToken(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("blacklist valid token", func(t *testing.T) {
		token := "test.jwt.token"
		expiresAt := time.Now().Add(15 * time.Minute)

		err := service.BlacklistToken(ctx, token, expiresAt)
		assert.NoError(t, err)

		// Verify token is blacklisted
		isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
	})

	t.Run("blacklist expired token should not error", func(t *testing.T) {
		token := "expired.jwt.token"
		expiresAt := time.Now().Add(-1 * time.Hour) // Already expired

		err := service.BlacklistToken(ctx, token, expiresAt)
		assert.NoError(t, err)

		// Expired token should not be in blacklist
		isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	t.Run("blacklist token with short TTL", func(t *testing.T) {
		token := "short.lived.token"
		expiresAt := time.Now().Add(2 * time.Second)

		err := service.BlacklistToken(ctx, token, expiresAt)
		assert.NoError(t, err)

		// Should be blacklisted immediately
		isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)

		// Fast forward time in miniredis
		mr.FastForward(3 * time.Second)

		// Should no longer be blacklisted after TTL
		isBlacklisted, err = service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})
}

func TestIsTokenBlacklisted(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("non-blacklisted token", func(t *testing.T) {
		token := "valid.jwt.token"

		isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.False(t, isBlacklisted)
	})

	t.Run("blacklisted token", func(t *testing.T) {
		token := "blacklisted.jwt.token"
		expiresAt := time.Now().Add(15 * time.Minute)

		// Blacklist the token
		err := service.BlacklistToken(ctx, token, expiresAt)
		require.NoError(t, err)

		// Check it's blacklisted
		isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
	})

	t.Run("same token different instances should match", func(t *testing.T) {
		token := "test.jwt.token"
		expiresAt := time.Now().Add(15 * time.Minute)

		// Blacklist with one service instance
		service1 := NewTokenBlacklistService(client)
		err := service1.BlacklistToken(ctx, token, expiresAt)
		require.NoError(t, err)

		// Check with different service instance
		service2 := NewTokenBlacklistService(client)
		isBlacklisted, err := service2.IsTokenBlacklisted(ctx, token)
		assert.NoError(t, err)
		assert.True(t, isBlacklisted)
	})
}

func TestBlacklistAllUserTokens(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("revoke all tokens for user", func(t *testing.T) {
		userID := "123e4567-e89b-12d3-a456-426614174000"

		err := service.BlacklistAllUserTokens(ctx, userID)
		assert.NoError(t, err)

		// Old token (issued before revocation) should be revoked
		oldTokenIssuedAt := time.Now().Add(-1 * time.Hour)
		isRevoked, err := service.IsUserTokenRevoked(ctx, userID, oldTokenIssuedAt)
		assert.NoError(t, err)
		assert.True(t, isRevoked)

		// New token (issued after revocation) should not be revoked
		newTokenIssuedAt := time.Now().Add(1 * time.Second)
		isRevoked, err = service.IsUserTokenRevoked(ctx, userID, newTokenIssuedAt)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
	})

	t.Run("different users should not affect each other", func(t *testing.T) {
		userID1 := "user-1"
		userID2 := "user-2"

		// Revoke user 1's tokens
		err := service.BlacklistAllUserTokens(ctx, userID1)
		require.NoError(t, err)

		tokenIssuedAt := time.Now().Add(-1 * time.Hour)

		// User 1 tokens should be revoked
		isRevoked, err := service.IsUserTokenRevoked(ctx, userID1, tokenIssuedAt)
		assert.NoError(t, err)
		assert.True(t, isRevoked)

		// User 2 tokens should NOT be revoked
		isRevoked, err = service.IsUserTokenRevoked(ctx, userID2, tokenIssuedAt)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
	})
}

func TestIsUserTokenRevoked(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("user with no revocation", func(t *testing.T) {
		userID := "no-revoke-user"
		tokenIssuedAt := time.Now()

		isRevoked, err := service.IsUserTokenRevoked(ctx, userID, tokenIssuedAt)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
	})

	t.Run("token issued before revocation", func(t *testing.T) {
		userID := "revoked-user"

		// Issue token
		tokenIssuedAt := time.Now()
		time.Sleep(10 * time.Millisecond)

		// Revoke all tokens
		err := service.BlacklistAllUserTokens(ctx, userID)
		require.NoError(t, err)

		// Check token issued before revocation
		isRevoked, err := service.IsUserTokenRevoked(ctx, userID, tokenIssuedAt)
		assert.NoError(t, err)
		assert.True(t, isRevoked)
	})

	t.Run("token issued after revocation", func(t *testing.T) {
		userID := "revoked-then-new-token"

		// Revoke all tokens
		err := service.BlacklistAllUserTokens(ctx, userID)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		// Issue new token after revocation
		tokenIssuedAt := time.Now()

		// New token should not be revoked
		isRevoked, err := service.IsUserTokenRevoked(ctx, userID, tokenIssuedAt)
		assert.NoError(t, err)
		assert.False(t, isRevoked)
	})
}

func TestClearUserTokenRevocation(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	userID := "clear-revoke-user"
	tokenIssuedAt := time.Now().Add(-1 * time.Hour)

	// Revoke all tokens
	err := service.BlacklistAllUserTokens(ctx, userID)
	require.NoError(t, err)

	// Verify revocation
	isRevoked, err := service.IsUserTokenRevoked(ctx, userID, tokenIssuedAt)
	assert.NoError(t, err)
	assert.True(t, isRevoked)

	// Clear revocation
	err = service.ClearUserTokenRevocation(ctx, userID)
	assert.NoError(t, err)

	// Verify revocation is cleared
	isRevoked, err = service.IsUserTokenRevoked(ctx, userID, tokenIssuedAt)
	assert.NoError(t, err)
	assert.False(t, isRevoked)
}

func TestGetBlacklistStats(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("empty blacklist", func(t *testing.T) {
		stats, err := service.GetBlacklistStats(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), stats["blacklisted_tokens"])
		assert.Equal(t, int64(0), stats["user_revocations"])
	})

	t.Run("with blacklisted tokens and user revocations", func(t *testing.T) {
		// Blacklist 3 tokens
		for i := 0; i < 3; i++ {
			token := "token-" + string(rune(i))
			expiresAt := time.Now().Add(15 * time.Minute)
			err := service.BlacklistToken(ctx, token, expiresAt)
			require.NoError(t, err)
		}

		// Revoke tokens for 2 users
		for i := 0; i < 2; i++ {
			userID := "user-" + string(rune(i))
			err := service.BlacklistAllUserTokens(ctx, userID)
			require.NoError(t, err)
		}

		stats, err := service.GetBlacklistStats(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), stats["blacklisted_tokens"])
		assert.Equal(t, int64(2), stats["user_revocations"])
	})
}

func TestHashToken(t *testing.T) {
	t.Run("same token produces same hash", func(t *testing.T) {
		token := "test.jwt.token"
		hash1 := hashToken(token)
		hash2 := hashToken(token)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("different tokens produce different hashes", func(t *testing.T) {
		token1 := "token1"
		token2 := "token2"
		hash1 := hashToken(token1)
		hash2 := hashToken(token2)
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("hash is deterministic", func(t *testing.T) {
		token := "deterministic.test.token"
		hash := hashToken(token)
		assert.Len(t, hash, 64) // SHA-256 produces 32 bytes = 64 hex chars
	})
}

func TestConcurrentBlacklistOperations(t *testing.T) {
	client, mr := setupTestRedis(t)
	defer mr.Close()

	service := NewTokenBlacklistService(client)
	ctx := context.Background()

	t.Run("concurrent blacklist and check operations", func(t *testing.T) {
		done := make(chan bool)
		numGoroutines := 10

		// Concurrent blacklist operations
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				token := "concurrent-token-" + string(rune(id))
				expiresAt := time.Now().Add(15 * time.Minute)

				err := service.BlacklistToken(ctx, token, expiresAt)
				assert.NoError(t, err)

				isBlacklisted, err := service.IsTokenBlacklisted(ctx, token)
				assert.NoError(t, err)
				assert.True(t, isBlacklisted)

				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}
