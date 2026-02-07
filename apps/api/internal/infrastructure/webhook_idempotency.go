package infrastructure

import (
	"context"
	"fmt"
	"time"

	"hustlex/internal/infrastructure/cache/redis"
)

// WebhookEventStore tracks processed webhook events to prevent duplicate processing
type WebhookEventStore interface {
	// IsProcessed checks if a webhook event has already been processed
	IsProcessed(ctx context.Context, eventID string) (bool, error)

	// MarkProcessed marks a webhook event as processed
	// TTL should be set long enough to cover typical retry windows (e.g., 7 days)
	MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error

	// GetProcessedAt returns when an event was processed
	GetProcessedAt(ctx context.Context, eventID string) (time.Time, error)
}

// RedisWebhookEventStore implements WebhookEventStore using Redis
type RedisWebhookEventStore struct {
	client *redis.Client
}

// NewRedisWebhookEventStore creates a new Redis-backed webhook event store
func NewRedisWebhookEventStore(client *redis.Client) *RedisWebhookEventStore {
	return &RedisWebhookEventStore{
		client: client,
	}
}

// IsProcessed checks if a webhook event has already been processed
func (s *RedisWebhookEventStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	key := webhookEventKey(eventID)
	exists, err := s.client.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check if webhook event processed: %w", err)
	}
	return exists, nil
}

// MarkProcessed marks a webhook event as processed with a TTL
func (s *RedisWebhookEventStore) MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error {
	key := webhookEventKey(eventID)

	data := map[string]interface{}{
		"event_id":     eventID,
		"processed_at": time.Now().UTC(),
	}

	if err := s.client.Set(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("failed to mark webhook event as processed: %w", err)
	}

	return nil
}

// GetProcessedAt returns when an event was processed
func (s *RedisWebhookEventStore) GetProcessedAt(ctx context.Context, eventID string) (time.Time, error) {
	key := webhookEventKey(eventID)

	var data struct {
		EventID     string    `json:"event_id"`
		ProcessedAt time.Time `json:"processed_at"`
	}

	if err := s.client.Get(ctx, key, &data); err != nil {
		if err == redis.ErrCacheMiss {
			return time.Time{}, ErrWebhookEventNotProcessed
		}
		return time.Time{}, fmt.Errorf("failed to get webhook event processed time: %w", err)
	}

	return data.ProcessedAt, nil
}

// webhookEventKey generates a Redis key for webhook event tracking
func webhookEventKey(eventID string) string {
	return fmt.Sprintf("webhook:event:%s", eventID)
}

// InMemoryWebhookEventStore implements WebhookEventStore using in-memory storage
// Used for testing or as a fallback when Redis is unavailable
type InMemoryWebhookEventStore struct {
	events map[string]time.Time
}

// NewInMemoryWebhookEventStore creates a new in-memory webhook event store
func NewInMemoryWebhookEventStore() *InMemoryWebhookEventStore {
	return &InMemoryWebhookEventStore{
		events: make(map[string]time.Time),
	}
}

// IsProcessed checks if a webhook event has already been processed
func (s *InMemoryWebhookEventStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	_, exists := s.events[eventID]
	return exists, nil
}

// MarkProcessed marks a webhook event as processed
// Note: TTL is not enforced in memory implementation (for simplicity in testing)
func (s *InMemoryWebhookEventStore) MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error {
	s.events[eventID] = time.Now().UTC()
	return nil
}

// GetProcessedAt returns when an event was processed
func (s *InMemoryWebhookEventStore) GetProcessedAt(ctx context.Context, eventID string) (time.Time, error) {
	processedAt, exists := s.events[eventID]
	if !exists {
		return time.Time{}, ErrWebhookEventNotProcessed
	}
	return processedAt, nil
}

// Error definitions
var (
	ErrWebhookEventNotProcessed = fmt.Errorf("webhook event not processed")
)
