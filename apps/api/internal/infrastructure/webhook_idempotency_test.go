package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryWebhookEventStore_MarkProcessed(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_test_12345"

	// Mark event as processed
	err := store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	require.NoError(t, err)

	// Verify event is marked as processed
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Event should be marked as processed")
}

func TestInMemoryWebhookEventStore_IsProcessed_NotProcessed(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_test_67890"

	// Check unprocessed event
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.False(t, isProcessed, "Event should not be marked as processed")
}

func TestInMemoryWebhookEventStore_GetProcessedAt(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_test_timestamp"
	beforeMark := time.Now().UTC()

	// Mark event as processed
	err := store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	require.NoError(t, err)

	afterMark := time.Now().UTC()

	// Get processed timestamp
	processedAt, err := store.GetProcessedAt(ctx, eventID)
	require.NoError(t, err)

	// Verify timestamp is within expected range
	assert.True(t, processedAt.After(beforeMark) || processedAt.Equal(beforeMark),
		"Processed timestamp should be after or equal to before mark time")
	assert.True(t, processedAt.Before(afterMark) || processedAt.Equal(afterMark),
		"Processed timestamp should be before or equal to after mark time")
}

func TestInMemoryWebhookEventStore_GetProcessedAt_NotProcessed(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_test_not_found"

	// Get processed timestamp for unprocessed event
	_, err := store.GetProcessedAt(ctx, eventID)
	assert.ErrorIs(t, err, ErrWebhookEventNotProcessed,
		"Should return ErrWebhookEventNotProcessed for unprocessed event")
}

func TestInMemoryWebhookEventStore_PreventDuplicateProcessing(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_test_duplicate"

	// Simulate webhook handler checking for duplicate
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.False(t, isProcessed, "First check should show event not processed")

	// Process the event
	err = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	require.NoError(t, err)

	// Simulate duplicate webhook delivery (retry)
	isProcessed, err = store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Second check should show event already processed")

	// Handler should skip processing and return 200 OK
}

func TestInMemoryWebhookEventStore_MultipleEvents(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	events := []string{
		"evt_charge_12345",
		"evt_charge_67890",
		"evt_transfer_11111",
		"evt_transfer_22222",
	}

	// Process first two events
	for i := 0; i < 2; i++ {
		err := store.MarkProcessed(ctx, events[i], 7*24*time.Hour)
		require.NoError(t, err)
	}

	// Verify only first two are processed
	for i, eventID := range events {
		isProcessed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)

		if i < 2 {
			assert.True(t, isProcessed, "Event %s should be processed", eventID)
		} else {
			assert.False(t, isProcessed, "Event %s should not be processed", eventID)
		}
	}
}

func TestWebhookEventKey(t *testing.T) {
	tests := []struct {
		name     string
		eventID  string
		expected string
	}{
		{
			name:     "Paystack event",
			eventID:  "evt_paystack_12345",
			expected: "webhook:event:evt_paystack_12345",
		},
		{
			name:     "Transfer event",
			eventID:  "trx_transfer_67890",
			expected: "webhook:event:trx_transfer_67890",
		},
		{
			name:     "Special characters",
			eventID:  "evt-test_123:456",
			expected: "webhook:event:evt-test_123:456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := webhookEventKey(tt.eventID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInMemoryWebhookEventStore_ConcurrentAccess(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_concurrent_test"

	// Simulate concurrent webhook deliveries
	done := make(chan bool, 2)

	// First goroutine
	go func() {
		isProcessed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)

		if !isProcessed {
			time.Sleep(10 * time.Millisecond) // Simulate processing time
			_ = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
		}
		done <- true
	}()

	// Second goroutine (race condition)
	go func() {
		time.Sleep(5 * time.Millisecond) // Start slightly after first

		isProcessed, err := store.IsProcessed(ctx, eventID)
		require.NoError(t, err)

		if !isProcessed {
			_ = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Verify event is marked as processed
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed)

	// Note: This test demonstrates the race condition issue with in-memory store.
	// Production code should use Redis with SETNX for atomic check-and-set.
}

func TestWebhookEventStore_Integration(t *testing.T) {
	// Integration test simulating real webhook flow
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	// Scenario: Paystack sends charge.success webhook
	eventID := "evt_charge_success_abc123"

	// Step 1: Webhook received, check if already processed
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.False(t, isProcessed, "New webhook should not be processed")

	// Step 2: Process the webhook (credit user wallet)
	// ... business logic here ...

	// Step 3: Mark as processed to prevent duplicate crediting
	err = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	require.NoError(t, err)

	// Step 4: Paystack retries webhook (network issue on their end)
	isProcessed, err = store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed, "Duplicate webhook should be detected")

	// Step 5: Handler returns 200 OK without reprocessing
	// This prevents double-crediting the user's wallet

	// Step 6: Verify processed timestamp
	processedAt, err := store.GetProcessedAt(ctx, eventID)
	require.NoError(t, err)
	assert.False(t, processedAt.IsZero(), "Processed timestamp should be set")
}

func TestWebhookEventStore_TTL(t *testing.T) {
	// Test TTL behavior (note: in-memory store doesn't enforce TTL)
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	eventID := "evt_ttl_test"

	// Mark with short TTL
	err := store.MarkProcessed(ctx, eventID, 1*time.Millisecond)
	require.NoError(t, err)

	// Event should still be marked (in-memory doesn't expire)
	isProcessed, err := store.IsProcessed(ctx, eventID)
	require.NoError(t, err)
	assert.True(t, isProcessed)

	// Note: Redis implementation will honor TTL automatically
}

func TestWebhookEventStore_ErrorHandling(t *testing.T) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	// Test empty event ID
	err := store.MarkProcessed(ctx, "", 7*24*time.Hour)
	assert.NoError(t, err, "Should allow empty event ID (edge case)")

	isProcessed, err := store.IsProcessed(ctx, "")
	require.NoError(t, err)
	assert.True(t, isProcessed)
}

// Benchmark tests

func BenchmarkInMemoryWebhookEventStore_MarkProcessed(b *testing.B) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventID := fmt.Sprintf("evt_bench_%d", i)
		_ = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	}
}

func BenchmarkInMemoryWebhookEventStore_IsProcessed(b *testing.B) {
	store := NewInMemoryWebhookEventStore()
	ctx := context.Background()

	// Pre-populate with 1000 events
	for i := 0; i < 1000; i++ {
		eventID := fmt.Sprintf("evt_bench_%d", i)
		_ = store.MarkProcessed(ctx, eventID, 7*24*time.Hour)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventID := fmt.Sprintf("evt_bench_%d", i%1000)
		_, _ = store.IsProcessed(ctx, eventID)
	}
}

// Note: We need to import fmt for benchmarks
import "fmt"
