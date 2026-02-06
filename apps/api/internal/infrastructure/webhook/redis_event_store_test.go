package webhook

import (
	"context"
	"testing"
	"time"

	"hustlex/internal/infrastructure/cache/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *redis.Client {
	// Use test Redis configuration
	config := redis.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       15, // Use separate DB for tests
		PoolSize: 10,
	}

	client, err := redis.NewClient(config)
	require.NoError(t, err, "Failed to connect to test Redis")

	return client
}

func cleanupTestRedis(t *testing.T, client *redis.Client, eventIDs []string) {
	ctx := context.Background()
	for _, eventID := range eventIDs {
		key := webhookEventKey(eventID)
		_ = client.Delete(ctx, key)
	}
}

func TestRedisEventStore_IsProcessed(t *testing.T) {
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	testCases := []struct {
		name      string
		eventID   string
		setup     func()
		wantExist bool
		wantErr   bool
	}{
		{
			name:      "event not processed",
			eventID:   "evt_unprocessed_001",
			setup:     func() {},
			wantExist: false,
			wantErr:   false,
		},
		{
			name:    "event already processed",
			eventID: "evt_processed_001",
			setup: func() {
				_ = store.MarkProcessed(ctx, "evt_processed_001", 1*time.Hour)
			},
			wantExist: true,
			wantErr:   false,
		},
		{
			name:      "empty event ID",
			eventID:   "",
			setup:     func() {},
			wantExist: false,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer cleanupTestRedis(t, client, []string{tc.eventID})

			exists, err := store.IsProcessed(ctx, tc.eventID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantExist, exists)
			}
		})
	}
}

func TestRedisEventStore_MarkProcessed(t *testing.T) {
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	testCases := []struct {
		name      string
		eventID   string
		expiresIn time.Duration
		setup     func()
		wantErr   error
	}{
		{
			name:      "mark new event as processed",
			eventID:   "evt_new_001",
			expiresIn: 1 * time.Hour,
			setup:     func() {},
			wantErr:   nil,
		},
		{
			name:      "mark already processed event (idempotency check)",
			eventID:   "evt_duplicate_001",
			expiresIn: 1 * time.Hour,
			setup: func() {
				_ = store.MarkProcessed(ctx, "evt_duplicate_001", 1*time.Hour)
			},
			wantErr: ErrEventAlreadyProcessed,
		},
		{
			name:      "empty event ID",
			eventID:   "",
			expiresIn: 1 * time.Hour,
			setup:     func() {},
			wantErr:   assert.AnError,
		},
		{
			name:      "zero expiration",
			eventID:   "evt_zero_expiry",
			expiresIn: 0,
			setup:     func() {},
			wantErr:   assert.AnError,
		},
		{
			name:      "negative expiration",
			eventID:   "evt_negative_expiry",
			expiresIn: -1 * time.Hour,
			setup:     func() {},
			wantErr:   assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer cleanupTestRedis(t, client, []string{tc.eventID})

			err := store.MarkProcessed(ctx, tc.eventID, tc.expiresIn)

			if tc.wantErr != nil {
				assert.Error(t, err)
				if tc.wantErr != assert.AnError {
					assert.ErrorIs(t, err, tc.wantErr)
				}
			} else {
				assert.NoError(t, err)

				// Verify event is now marked as processed
				exists, err := store.IsProcessed(ctx, tc.eventID)
				assert.NoError(t, err)
				assert.True(t, exists)
			}
		})
	}
}

func TestRedisEventStore_GetProcessedAt(t *testing.T) {
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	testCases := []struct {
		name       string
		eventID    string
		setup      func() time.Time
		wantExists bool
		wantErr    bool
	}{
		{
			name:    "get processed event timestamp",
			eventID: "evt_timestamp_001",
			setup: func() time.Time {
				beforeTime := time.Now().UTC()
				_ = store.MarkProcessed(ctx, "evt_timestamp_001", 1*time.Hour)
				return beforeTime
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name:       "get unprocessed event",
			eventID:    "evt_never_processed",
			setup:      func() time.Time { return time.Time{} },
			wantExists: false,
			wantErr:    false,
		},
		{
			name:       "empty event ID",
			eventID:    "",
			setup:      func() time.Time { return time.Time{} },
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			beforeTime := tc.setup()
			defer cleanupTestRedis(t, client, []string{tc.eventID})

			processedAt, exists, err := store.GetProcessedAt(ctx, tc.eventID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantExists, exists)

				if tc.wantExists {
					// Verify timestamp is reasonable
					assert.False(t, processedAt.IsZero())
					assert.True(t, processedAt.After(beforeTime) || processedAt.Equal(beforeTime))
					assert.True(t, processedAt.Before(time.Now().UTC().Add(1*time.Second)))
				} else {
					assert.True(t, processedAt.IsZero())
				}
			}
		})
	}
}

func TestRedisEventStore_Delete(t *testing.T) {
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	testCases := []struct {
		name    string
		eventID string
		setup   func()
		wantErr bool
	}{
		{
			name:    "delete existing event",
			eventID: "evt_delete_001",
			setup: func() {
				_ = store.MarkProcessed(ctx, "evt_delete_001", 1*time.Hour)
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent event (idempotent)",
			eventID: "evt_never_existed",
			setup:   func() {},
			wantErr: false,
		},
		{
			name:    "empty event ID",
			eventID: "",
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			defer cleanupTestRedis(t, client, []string{tc.eventID})

			err := store.Delete(ctx, tc.eventID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify event no longer exists
				exists, err := store.IsProcessed(ctx, tc.eventID)
				assert.NoError(t, err)
				assert.False(t, exists)
			}
		})
	}
}

func TestRedisEventStore_RaceCondition(t *testing.T) {
	// Test that concurrent MarkProcessed calls for the same event
	// only allow one to succeed (SetNX atomicity)
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	eventID := "evt_race_001"
	defer cleanupTestRedis(t, client, []string{eventID})

	// Launch 10 concurrent attempts to mark the same event
	results := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			err := store.MarkProcessed(ctx, eventID, 1*time.Hour)
			results <- err
		}()
	}

	// Collect results
	successCount := 0
	alreadyProcessedCount := 0
	for i := 0; i < 10; i++ {
		err := <-results
		if err == nil {
			successCount++
		} else if err == ErrEventAlreadyProcessed {
			alreadyProcessedCount++
		} else {
			t.Errorf("Unexpected error: %v", err)
		}
	}

	// Exactly one should succeed, the rest should get ErrEventAlreadyProcessed
	assert.Equal(t, 1, successCount, "Exactly one goroutine should succeed")
	assert.Equal(t, 9, alreadyProcessedCount, "Nine goroutines should get already processed error")
}

func TestRedisEventStore_Expiration(t *testing.T) {
	// Test that events expire after the specified duration
	client := setupTestRedis(t)
	store := NewRedisEventStore(client)
	ctx := context.Background()

	eventID := "evt_expire_001"
	defer cleanupTestRedis(t, client, []string{eventID})

	// Mark event with 2 second expiration
	err := store.MarkProcessed(ctx, eventID, 2*time.Second)
	require.NoError(t, err)

	// Verify it exists immediately
	exists, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, exists)

	// Wait for expiration (with buffer)
	time.Sleep(3 * time.Second)

	// Verify it no longer exists
	exists, err = store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.False(t, exists, "Event should have expired")
}

func TestRedisEventStore_WebhookKeyFormat(t *testing.T) {
	// Verify the key format is correct
	eventID := "evt_test_123"
	expectedKey := "webhook:event:evt_test_123"

	actualKey := webhookEventKey(eventID)
	assert.Equal(t, expectedKey, actualKey)
}
