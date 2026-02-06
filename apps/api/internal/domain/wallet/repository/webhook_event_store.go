package repository

import (
	"context"
	"time"

	"hustlex/internal/domain/wallet/event"
)

// WebhookEventStore tracks processed webhook events for idempotency
type WebhookEventStore interface {
	// IsProcessed checks if a webhook event has already been processed
	IsProcessed(ctx context.Context, eventID event.WebhookEventID) (bool, error)

	// MarkProcessed records that a webhook event has been processed
	// Returns an error if the event was already processed
	MarkProcessed(ctx context.Context, webhookEvent *event.WebhookEvent) error

	// GetEvent retrieves a processed webhook event
	GetEvent(ctx context.Context, eventID event.WebhookEventID) (*event.WebhookEvent, error)

	// CleanupExpired removes webhook events older than the retention period
	// This is used for housekeeping to prevent unbounded growth
	CleanupExpired(ctx context.Context, retentionPeriod time.Duration) error
}
