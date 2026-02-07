package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hustlex/internal/domain/wallet/event"
	"hustlex/internal/domain/wallet/repository"
	"hustlex/internal/infrastructure/cache/redis"
)

// WebhookEventStore implements repository.WebhookEventStore using Redis
type WebhookEventStore struct {
	client            *redis.Client
	retentionPeriod   time.Duration // How long to keep processed events
	keyPrefix         string
}

// NewWebhookEventStore creates a new Redis-backed webhook event store
func NewWebhookEventStore(client *redis.Client, retentionPeriod time.Duration) repository.WebhookEventStore {
	if retentionPeriod == 0 {
		retentionPeriod = 30 * 24 * time.Hour // Default: 30 days
	}

	return &WebhookEventStore{
		client:          client,
		retentionPeriod: retentionPeriod,
		keyPrefix:       "webhook:event:",
	}
}

// IsProcessed checks if a webhook event has already been processed
func (s *WebhookEventStore) IsProcessed(ctx context.Context, eventID event.WebhookEventID) (bool, error) {
	key := s.eventKey(eventID)
	exists, err := s.client.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check if webhook event exists: %w", err)
	}
	return exists, nil
}

// MarkProcessed records that a webhook event has been processed
func (s *WebhookEventStore) MarkProcessed(ctx context.Context, webhookEvent *event.WebhookEvent) error {
	key := s.eventKey(webhookEvent.EventID)

	// Use SetNX to atomically set only if the key doesn't exist
	// This prevents race conditions where multiple webhook deliveries arrive simultaneously
	success, err := s.client.SetNX(ctx, key, webhookEvent, s.retentionPeriod)
	if err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}

	if !success {
		// Event was already processed (key already exists)
		return ErrWebhookAlreadyProcessed
	}

	return nil
}

// GetEvent retrieves a processed webhook event
func (s *WebhookEventStore) GetEvent(ctx context.Context, eventID event.WebhookEventID) (*event.WebhookEvent, error) {
	key := s.eventKey(eventID)

	var webhookEvent event.WebhookEvent
	err := s.client.Get(ctx, key, &webhookEvent)
	if err != nil {
		if errors.Is(err, redis.ErrCacheMiss) {
			return nil, ErrWebhookEventNotFound
		}
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}

	return &webhookEvent, nil
}

// CleanupExpired removes webhook events older than the retention period
// Note: Redis handles expiration automatically via TTL, so this is a no-op
// It's included for interface completeness and potential future use with other storage backends
func (s *WebhookEventStore) CleanupExpired(ctx context.Context, retentionPeriod time.Duration) error {
	// Redis handles expiration automatically via TTL
	// No cleanup needed
	return nil
}

// eventKey generates the Redis key for a webhook event
func (s *WebhookEventStore) eventKey(eventID event.WebhookEventID) string {
	return fmt.Sprintf("%s%s", s.keyPrefix, eventID)
}

// Errors
var (
	ErrWebhookAlreadyProcessed = errors.New("webhook event already processed")
	ErrWebhookEventNotFound    = errors.New("webhook event not found")
)

// StoredWebhookEvent represents the data structure stored in Redis
// This is used for JSON serialization
type StoredWebhookEvent struct {
	EventID     string    `json:"event_id"`
	Provider    string    `json:"provider"`
	EventType   string    `json:"event_type"`
	Reference   string    `json:"reference"`
	ProcessedAt time.Time `json:"processed_at"`
	Payload     string    `json:"payload"` // Base64 encoded
}

// MarshalJSON implements json.Marshaler for WebhookEvent
func (e *event.WebhookEvent) MarshalJSON() ([]byte, error) {
	stored := StoredWebhookEvent{
		EventID:     string(e.EventID),
		Provider:    e.Provider,
		EventType:   e.EventType,
		Reference:   e.Reference,
		ProcessedAt: e.ProcessedAt,
		Payload:     string(e.Payload),
	}
	return json.Marshal(stored)
}

// UnmarshalJSON implements json.Unmarshaler for WebhookEvent
func (e *event.WebhookEvent) UnmarshalJSON(data []byte) error {
	var stored StoredWebhookEvent
	if err := json.Unmarshal(data, &stored); err != nil {
		return err
	}

	e.EventID = event.WebhookEventID(stored.EventID)
	e.Provider = stored.Provider
	e.EventType = stored.EventType
	e.Reference = stored.Reference
	e.ProcessedAt = stored.ProcessedAt
	e.Payload = []byte(stored.Payload)

	return nil
}
