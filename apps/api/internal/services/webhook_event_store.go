package services

import (
	"context"
	"fmt"
	"time"

	"hustlex/internal/infrastructure/cache/redis"
)

// WebhookEventStore manages idempotency for webhook events to prevent duplicate processing
type WebhookEventStore interface {
	// IsProcessed checks if a webhook event has already been processed
	IsProcessed(ctx context.Context, eventID string) (bool, error)

	// MarkProcessed marks a webhook event as processed with a TTL
	// The TTL should be longer than the payment provider's retry window
	MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error

	// GetProcessedAt returns when the event was processed (if it was)
	GetProcessedAt(ctx context.Context, eventID string) (*time.Time, error)
}

// RedisWebhookEventStore implements WebhookEventStore using Redis
type RedisWebhookEventStore struct {
	redis *redis.Client
}

// NewRedisWebhookEventStore creates a new Redis-backed webhook event store
func NewRedisWebhookEventStore(redisClient *redis.Client) *RedisWebhookEventStore {
	return &RedisWebhookEventStore{
		redis: redisClient,
	}
}

// IsProcessed checks if a webhook event has already been processed
func (s *RedisWebhookEventStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	key := s.eventKey(eventID)
	return s.redis.Exists(ctx, key)
}

// MarkProcessed marks a webhook event as processed with a TTL
func (s *RedisWebhookEventStore) MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error {
	key := s.eventKey(eventID)
	data := map[string]interface{}{
		"event_id":     eventID,
		"processed_at": time.Now().UTC(),
	}
	return s.redis.Set(ctx, key, data, ttl)
}

// GetProcessedAt returns when the event was processed (if it was)
func (s *RedisWebhookEventStore) GetProcessedAt(ctx context.Context, eventID string) (*time.Time, error) {
	key := s.eventKey(eventID)

	var data struct {
		ProcessedAt time.Time `json:"processed_at"`
	}

	if err := s.redis.Get(ctx, key, &data); err != nil {
		if err == redis.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	}

	return &data.ProcessedAt, nil
}

// eventKey generates a Redis key for a webhook event
func (s *RedisWebhookEventStore) eventKey(eventID string) string {
	return fmt.Sprintf("webhook:event:%s", eventID)
}

// Default TTL constants for different webhook event types
const (
	// DefaultWebhookEventTTL is the default TTL for webhook events (7 days)
	// This should be longer than the payment provider's retry window
	DefaultWebhookEventTTL = 7 * 24 * time.Hour

	// PaystackRetryWindow is Paystack's webhook retry window (typically 3 days)
	PaystackRetryWindow = 3 * 24 * time.Hour

	// FlutterwaveRetryWindow is Flutterwave's webhook retry window (typically 3 days)
	FlutterwaveRetryWindow = 3 * 24 * time.Hour
)
