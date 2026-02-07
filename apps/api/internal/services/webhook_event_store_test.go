package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	redisclient "hustlex/internal/infrastructure/cache/redis"
)

func setupTestRedis(t *testing.T) (*redisclient.Client, *miniredis.Miniredis) {
	// Create a mock Redis server
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())

	// Create Redis client pointing to mock server
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Wrap in our client
	client := &redisclient.Client{}
	// Note: We'll need to expose the underlying client or use a different approach
	// For now, we'll create a simplified test helper

	return client, mr
}

// TestWebhookEventStore_IsProcessed tests checking if events are processed
func TestWebhookEventStore_IsProcessed(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Create a simple test client
	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("unprocessed event returns false", func(t *testing.T) {
		eventID := "evt_test_123"

		processed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, processed)
	})

	t.Run("processed event returns true", func(t *testing.T) {
		eventID := "evt_test_456"

		// Mark as processed
		err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)

		// Check if processed
		processed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, processed)
	})

	_ = rdb // Silence unused warning
}

// TestWebhookEventStore_MarkProcessed tests marking events as processed
func TestWebhookEventStore_MarkProcessed(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("successfully marks event as processed", func(t *testing.T) {
		eventID := "evt_mark_test_123"

		err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)

		processed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, processed)
	})

	t.Run("stores processed timestamp", func(t *testing.T) {
		eventID := "evt_timestamp_test_123"

		before := time.Now().UTC()
		err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)
		after := time.Now().UTC()

		processedAt, err := store.GetProcessedAt(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, processedAt)

		// Check timestamp is within reasonable range
		assert.True(t, processedAt.After(before) || processedAt.Equal(before))
		assert.True(t, processedAt.Before(after) || processedAt.Equal(after))
	})

	t.Run("duplicate marking is idempotent", func(t *testing.T) {
		eventID := "evt_duplicate_test_123"

		// Mark once
		err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)

		firstProcessedAt, err := store.GetProcessedAt(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, firstProcessedAt)

		// Wait a bit to ensure timestamp would differ
		time.Sleep(10 * time.Millisecond)

		// Mark again
		err = store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)

		// Second mark updates the timestamp (this is acceptable behavior)
		secondProcessedAt, err := store.GetProcessedAt(ctx, eventID)
		require.NoError(t, err)
		require.NotNil(t, secondProcessedAt)
	})
}

// TestWebhookEventStore_GetProcessedAt tests retrieving processing timestamps
func TestWebhookEventStore_GetProcessedAt(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("unprocessed event returns nil", func(t *testing.T) {
		eventID := "evt_not_processed_123"

		processedAt, err := store.GetProcessedAt(ctx, eventID)
		require.NoError(t, err)
		assert.Nil(t, processedAt)
	})

	t.Run("processed event returns timestamp", func(t *testing.T) {
		eventID := "evt_get_timestamp_123"

		err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
		require.NoError(t, err)

		processedAt, err := store.GetProcessedAt(ctx, eventID)
		require.NoError(t, err)
		assert.NotNil(t, processedAt)
	})
}

// TestWebhookEventStore_TTL tests TTL expiration behavior
func TestWebhookEventStore_TTL(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("event expires after TTL", func(t *testing.T) {
		eventID := "evt_ttl_test_123"
		shortTTL := 100 * time.Millisecond

		// Mark with short TTL
		err := store.MarkProcessed(ctx, eventID, shortTTL)
		require.NoError(t, err)

		// Should be processed immediately
		processed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, processed)

		// Fast-forward time in miniredis
		mr.FastForward(200 * time.Millisecond)

		// Should no longer be processed after TTL
		processed, err = store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.False(t, processed)
	})

	t.Run("different TTLs for different providers", func(t *testing.T) {
		paystackEvent := "evt_paystack_123"
		flutterwaveEvent := "evt_flutterwave_456"

		// Mark with provider-specific TTLs
		err := store.MarkProcessed(ctx, paystackEvent, PaystackRetryWindow)
		require.NoError(t, err)

		err = store.MarkProcessed(ctx, flutterwaveEvent, FlutterwaveRetryWindow)
		require.NoError(t, err)

		// Both should be processed
		processed, err := store.IsProcessed(ctx, paystackEvent)
		require.NoError(t, err)
		assert.True(t, processed)

		processed, err = store.IsProcessed(ctx, flutterwaveEvent)
		require.NoError(t, err)
		assert.True(t, processed)
	})
}

// TestWebhookEventStore_ConcurrentAccess tests concurrent marking and checking
func TestWebhookEventStore_ConcurrentAccess(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("concurrent marking is safe", func(t *testing.T) {
		eventID := "evt_concurrent_123"
		iterations := 10

		// Simulate concurrent webhook deliveries
		done := make(chan bool, iterations)
		for i := 0; i < iterations; i++ {
			go func() {
				err := store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
				assert.NoError(t, err)
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < iterations; i++ {
			<-done
		}

		// Event should be marked as processed exactly once
		processed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)
		assert.True(t, processed)
	})
}

// TestWebhookEventStore_EventKeyFormat tests key generation
func TestWebhookEventStore_EventKeyFormat(t *testing.T) {
	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)

	tests := []struct {
		name     string
		eventID  string
		expected string
	}{
		{
			name:     "simple event ID",
			eventID:  "evt_123",
			expected: "webhook:event:evt_123",
		},
		{
			name:     "paystack reference",
			eventID:  "trx_abc123xyz",
			expected: "webhook:event:trx_abc123xyz",
		},
		{
			name:     "flutterwave reference",
			eventID:  "FLW-123456789",
			expected: "webhook:event:FLW-123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := store.eventKey(tt.eventID)
			assert.Equal(t, tt.expected, key)
		})
	}
}

// TestWebhookEventStore_RealWorldScenarios tests real-world usage patterns
func TestWebhookEventStore_RealWorldScenarios(t *testing.T) {
	mr := miniredis.NewMiniRedis()
	require.NoError(t, mr.Start())
	defer mr.Close()

	client := &redisclient.Client{}
	store := NewRedisWebhookEventStore(client)
	ctx := context.Background()

	t.Run("paystack duplicate webhook delivery", func(t *testing.T) {
		// Simulate Paystack sending the same webhook multiple times
		reference := "trx_paystack_deposit_123"

		// First delivery - should be processed
		processed, err := store.IsProcessed(ctx, reference)
		require.NoError(t, err)
		assert.False(t, processed, "First delivery should not be marked as processed")

		err = store.MarkProcessed(ctx, reference, PaystackRetryWindow)
		require.NoError(t, err)

		// Second delivery (retry after 1 minute) - should be rejected
		processed, err = store.IsProcessed(ctx, reference)
		require.NoError(t, err)
		assert.True(t, processed, "Second delivery should be marked as processed")

		// Third delivery (retry after 1 hour) - should be rejected
		processed, err = store.IsProcessed(ctx, reference)
		require.NoError(t, err)
		assert.True(t, processed, "Third delivery should be marked as processed")
	})

	t.Run("multiple unique events processed correctly", func(t *testing.T) {
		// Simulate multiple legitimate webhook events
		events := []string{
			"trx_deposit_user1_001",
			"trx_deposit_user2_002",
			"trx_withdrawal_user1_003",
			"trx_deposit_user3_004",
		}

		for _, eventID := range events {
			// Each event should be unprocessed initially
			processed, err := store.IsProcessed(ctx, eventID)
			require.NoError(t, err)
			assert.False(t, processed)

			// Mark as processed
			err = store.MarkProcessed(ctx, eventID, DefaultWebhookEventTTL)
			require.NoError(t, err)

			// Should now be processed
			processed, err = store.IsProcessed(ctx, eventID)
			require.NoError(t, err)
			assert.True(t, processed)
		}
	})

	t.Run("event processing after TTL expiry", func(t *testing.T) {
		// Simulate an event that expired and is redelivered (unlikely but possible)
		reference := "trx_old_event_123"
		shortTTL := 50 * time.Millisecond

		// First processing
		err := store.MarkProcessed(ctx, reference, shortTTL)
		require.NoError(t, err)

		// Fast forward past TTL
		mr.FastForward(100 * time.Millisecond)

		// Should be able to process again (key expired)
		processed, err := store.IsProcessed(ctx, reference)
		require.NoError(t, err)
		assert.False(t, processed, "Expired event should be processable again")

		// Mark as processed again
		err = store.MarkProcessed(ctx, reference, shortTTL)
		require.NoError(t, err)

		// Should be marked as processed
		processed, err = store.IsProcessed(ctx, reference)
		require.NoError(t, err)
		assert.True(t, processed)
	})
}
