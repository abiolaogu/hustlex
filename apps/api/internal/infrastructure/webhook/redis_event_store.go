package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hustlex/internal/domain/webhook/repository"
	redisclient "hustlex/internal/infrastructure/cache/redis"

	"github.com/redis/go-redis/v9"
)

// RedisEventStore implements WebhookEventStore using Redis
type RedisEventStore struct {
	client *redisclient.Client
}

// NewRedisEventStore creates a new Redis-backed webhook event store
func NewRedisEventStore(client *redisclient.Client) repository.WebhookEventStore {
	return &RedisEventStore{
		client: client,
	}
}

// webhookEventKey generates the Redis key for a webhook event
func webhookEventKey(eventID string) string {
	return fmt.Sprintf("webhook:event:%s", eventID)
}

// eventRecord stores metadata about a processed webhook event
type eventRecord struct {
	EventID     string    `json:"event_id"`
	ProcessedAt time.Time `json:"processed_at"`
	Version     int       `json:"version"` // For future schema changes
}

// IsProcessed checks if a webhook event has already been processed
func (s *RedisEventStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	if eventID == "" {
		return false, errors.New("eventID cannot be empty")
	}

	key := webhookEventKey(eventID)
	exists, err := s.client.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to check event existence: %w", err)
	}

	return exists, nil
}

// MarkProcessed marks a webhook event as processed
func (s *RedisEventStore) MarkProcessed(ctx context.Context, eventID string, expiresIn time.Duration) error {
	if eventID == "" {
		return errors.New("eventID cannot be empty")
	}

	if expiresIn <= 0 {
		return errors.New("expiresIn must be positive")
	}

	key := webhookEventKey(eventID)
	record := eventRecord{
		EventID:     eventID,
		ProcessedAt: time.Now().UTC(),
		Version:     1,
	}

	// Use SetNX to ensure atomicity - only set if key doesn't exist
	// This prevents race conditions where multiple webhook deliveries arrive simultaneously
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal event record: %w", err)
	}

	// SetNX returns false if key already exists
	set, err := s.client.SetNX(ctx, key, data, expiresIn)
	if err != nil {
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}

	if !set {
		// Key already exists - event was processed by another request
		return ErrEventAlreadyProcessed
	}

	return nil
}

// GetProcessedAt retrieves when an event was first processed
func (s *RedisEventStore) GetProcessedAt(ctx context.Context, eventID string) (time.Time, bool, error) {
	if eventID == "" {
		return time.Time{}, false, errors.New("eventID cannot be empty")
	}

	key := webhookEventKey(eventID)
	var record eventRecord

	err := s.client.Get(ctx, key, &record)
	if err != nil {
		if errors.Is(err, redisclient.ErrCacheMiss) {
			return time.Time{}, false, nil
		}
		return time.Time{}, false, fmt.Errorf("failed to get event record: %w", err)
	}

	return record.ProcessedAt, true, nil
}

// Delete removes an event record
func (s *RedisEventStore) Delete(ctx context.Context, eventID string) error {
	if eventID == "" {
		return errors.New("eventID cannot be empty")
	}

	key := webhookEventKey(eventID)
	return s.client.Delete(ctx, key)
}

// Errors
var (
	// ErrEventAlreadyProcessed is returned when attempting to mark an event that's already been processed
	ErrEventAlreadyProcessed = errors.New("webhook event already processed")
)
