package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"

	sharedevent "hustlex/internal/domain/shared/event"
)

// EventHandler is a function that handles domain events
type EventHandler func(ctx context.Context, event sharedevent.DomainEvent) error

// InMemoryEventBus is a simple in-memory event bus for development/testing
type InMemoryEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
	async    bool
}

// NewInMemoryEventBus creates a new in-memory event bus
func NewInMemoryEventBus(async bool) *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]EventHandler),
		async:    async,
	}
}

// Subscribe registers a handler for an event type
func (b *InMemoryEventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// SubscribeAll registers a handler for all events
func (b *InMemoryEventBus) SubscribeAll(handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers["*"] = append(b.handlers["*"], handler)
}

// Publish publishes events to all registered handlers
func (b *InMemoryEventBus) Publish(ctx context.Context, events []interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, event := range events {
		domainEvent, ok := event.(sharedevent.DomainEvent)
		if !ok {
			continue
		}

		eventType := domainEvent.EventType()

		// Get specific handlers
		handlers := b.handlers[eventType]

		// Add wildcard handlers
		handlers = append(handlers, b.handlers["*"]...)

		for _, handler := range handlers {
			if b.async {
				go func(h EventHandler, e sharedevent.DomainEvent) {
					if err := h(ctx, e); err != nil {
						log.Printf("Error handling event %s: %v", e.EventType(), err)
					}
				}(handler, domainEvent)
			} else {
				if err := handler(ctx, domainEvent); err != nil {
					log.Printf("Error handling event %s: %v", eventType, err)
				}
			}
		}
	}

	return nil
}

// RedisEventBus uses Redis pub/sub for event distribution
type RedisEventBus struct {
	redisClient RedisPublisher
	handlers    map[string][]EventHandler
	mu          sync.RWMutex
	subscribing bool
}

// RedisPublisher interface for Redis publish operations
type RedisPublisher interface {
	Publish(ctx context.Context, channel string, message interface{}) error
}

// NewRedisEventBus creates a new Redis-based event bus
func NewRedisEventBus(redisClient RedisPublisher) *RedisEventBus {
	return &RedisEventBus{
		redisClient: redisClient,
		handlers:    make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for an event type
func (b *RedisEventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish publishes events to Redis
func (b *RedisEventBus) Publish(ctx context.Context, events []interface{}) error {
	for _, event := range events {
		domainEvent, ok := event.(sharedevent.DomainEvent)
		if !ok {
			continue
		}

		eventType := domainEvent.EventType()
		channel := fmt.Sprintf("events:%s", eventType)

		// Create event envelope
		envelope := EventEnvelope{
			EventType:   eventType,
			EventID:     domainEvent.EventID(),
			OccurredAt:  domainEvent.OccurredAt(),
			Payload:     event,
			PayloadType: reflect.TypeOf(event).String(),
		}

		data, err := json.Marshal(envelope)
		if err != nil {
			return err
		}

		if err := b.redisClient.Publish(ctx, channel, data); err != nil {
			return err
		}

		// Also publish to all-events channel
		if err := b.redisClient.Publish(ctx, "events:*", data); err != nil {
			return err
		}
	}

	return nil
}

// EventEnvelope wraps events for transport
type EventEnvelope struct {
	EventType   string      `json:"event_type"`
	EventID     string      `json:"event_id"`
	OccurredAt  interface{} `json:"occurred_at"`
	Payload     interface{} `json:"payload"`
	PayloadType string      `json:"payload_type"`
}

// OutboxEventBus implements the transactional outbox pattern
type OutboxEventBus struct {
	db          OutboxStore
	publisher   EventPublisher
	pollInterval int // in seconds
}

// OutboxStore interface for outbox persistence
type OutboxStore interface {
	SaveEvent(ctx context.Context, event OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkEventPublished(ctx context.Context, eventID string) error
	MarkEventFailed(ctx context.Context, eventID string, err string) error
}

// EventPublisher interface for publishing events
type EventPublisher interface {
	Publish(ctx context.Context, events []interface{}) error
}

// OutboxEvent represents an event in the outbox
type OutboxEvent struct {
	ID          string
	EventType   string
	Payload     []byte
	Status      string // pending, published, failed
	Retries     int
	CreatedAt   interface{}
	PublishedAt interface{}
	Error       string
}

// NewOutboxEventBus creates a new outbox-based event bus
func NewOutboxEventBus(db OutboxStore, publisher EventPublisher) *OutboxEventBus {
	return &OutboxEventBus{
		db:           db,
		publisher:    publisher,
		pollInterval: 5,
	}
}

// Publish saves events to the outbox for later publishing
func (b *OutboxEventBus) Publish(ctx context.Context, events []interface{}) error {
	for _, event := range events {
		domainEvent, ok := event.(sharedevent.DomainEvent)
		if !ok {
			continue
		}

		payload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		outboxEvent := OutboxEvent{
			ID:        domainEvent.EventID(),
			EventType: domainEvent.EventType(),
			Payload:   payload,
			Status:    "pending",
			CreatedAt: domainEvent.OccurredAt(),
		}

		if err := b.db.SaveEvent(ctx, outboxEvent); err != nil {
			return err
		}
	}

	return nil
}

// ProcessOutbox polls the outbox and publishes pending events
func (b *OutboxEventBus) ProcessOutbox(ctx context.Context) error {
	events, err := b.db.GetPendingEvents(ctx, 100)
	if err != nil {
		return err
	}

	for _, event := range events {
		var payload interface{}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			if markErr := b.db.MarkEventFailed(ctx, event.ID, err.Error()); markErr != nil {
				log.Printf("Failed to mark event as failed: %v", markErr)
			}
			continue
		}

		if err := b.publisher.Publish(ctx, []interface{}{payload}); err != nil {
			if markErr := b.db.MarkEventFailed(ctx, event.ID, err.Error()); markErr != nil {
				log.Printf("Failed to mark event as failed: %v", markErr)
			}
			continue
		}

		if err := b.db.MarkEventPublished(ctx, event.ID); err != nil {
			log.Printf("Failed to mark event as published: %v", err)
		}
	}

	return nil
}

// Ensure implementations satisfy interfaces
var _ EventPublisher = (*InMemoryEventBus)(nil)
var _ EventPublisher = (*RedisEventBus)(nil)
var _ EventPublisher = (*OutboxEventBus)(nil)
