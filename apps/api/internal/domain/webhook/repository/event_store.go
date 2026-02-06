package repository

import (
	"context"
	"time"
)

// WebhookEventStore manages webhook event idempotency tracking
// to prevent duplicate processing of webhook events
type WebhookEventStore interface {
	// IsProcessed checks if a webhook event has already been processed
	// Returns true if the event was previously processed, false otherwise
	IsProcessed(ctx context.Context, eventID string) (bool, error)

	// MarkProcessed marks a webhook event as processed
	// The event record will be retained for the specified duration to prevent replays
	MarkProcessed(ctx context.Context, eventID string, expiresIn time.Duration) error

	// GetProcessedAt retrieves when an event was first processed
	// Returns the timestamp and true if the event was processed, or zero time and false otherwise
	GetProcessedAt(ctx context.Context, eventID string) (time.Time, bool, error)

	// Delete removes an event record (used for testing or manual cleanup)
	Delete(ctx context.Context, eventID string) error
}
