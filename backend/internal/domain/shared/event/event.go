package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the base interface for all domain events
type DomainEvent interface {
	// EventID returns the unique identifier for this event instance
	EventID() string
	// EventType returns the type name of the event
	EventType() string
	// AggregateID returns the ID of the aggregate that produced this event
	AggregateID() string
	// AggregateType returns the type of aggregate
	AggregateType() string
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	// Version returns the event version for schema evolution
	Version() int
}

// BaseEvent provides common fields for all domain events
type BaseEvent struct {
	ID            string    `json:"event_id"`
	Type          string    `json:"event_type"`
	AggregateIDV  string    `json:"aggregate_id"`
	AggregateTypeV string   `json:"aggregate_type"`
	Occurred      time.Time `json:"occurred_at"`
	Ver           int       `json:"version"`
}

// NewBaseEvent creates a new base event
func NewBaseEvent(eventType, aggregateID, aggregateType string) BaseEvent {
	return BaseEvent{
		ID:             uuid.NewString(),
		Type:           eventType,
		AggregateIDV:   aggregateID,
		AggregateTypeV: aggregateType,
		Occurred:       time.Now().UTC(),
		Ver:            1,
	}
}

func (e BaseEvent) EventID() string       { return e.ID }
func (e BaseEvent) EventType() string     { return e.Type }
func (e BaseEvent) AggregateID() string   { return e.AggregateIDV }
func (e BaseEvent) AggregateType() string { return e.AggregateTypeV }
func (e BaseEvent) OccurredAt() time.Time { return e.Occurred }
func (e BaseEvent) Version() int          { return e.Ver }

// EventEnvelope wraps a domain event for transport/storage
type EventEnvelope struct {
	EventID       string          `json:"event_id"`
	EventType     string          `json:"event_type"`
	AggregateID   string          `json:"aggregate_id"`
	AggregateType string          `json:"aggregate_type"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Version       int             `json:"version"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      EventMetadata   `json:"metadata"`
}

// EventMetadata contains optional metadata for events
type EventMetadata struct {
	CorrelationID string            `json:"correlation_id,omitempty"`
	CausationID   string            `json:"causation_id,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	TraceID       string            `json:"trace_id,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}

// Wrap wraps a domain event in an envelope
func Wrap(event DomainEvent, metadata EventMetadata) (*EventEnvelope, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return &EventEnvelope{
		EventID:       event.EventID(),
		EventType:     event.EventType(),
		AggregateID:   event.AggregateID(),
		AggregateType: event.AggregateType(),
		OccurredAt:    event.OccurredAt(),
		Version:       event.Version(),
		Payload:       payload,
		Metadata:      metadata,
	}, nil
}

// AggregateRoot provides event collection capability for aggregates
type AggregateRoot struct {
	events []DomainEvent
}

// RecordEvent adds an event to be published after persistence
func (ar *AggregateRoot) RecordEvent(event DomainEvent) {
	ar.events = append(ar.events, event)
}

// DomainEvents returns and clears the collected events
func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	events := ar.events
	ar.events = make([]DomainEvent, 0)
	return events
}

// ClearEvents clears all collected events
func (ar *AggregateRoot) ClearEvents() {
	ar.events = make([]DomainEvent, 0)
}

// HasEvents returns true if there are pending events
func (ar *AggregateRoot) HasEvents() bool {
	return len(ar.events) > 0
}

// EventHandler is a function that handles a domain event
type EventHandler func(event DomainEvent) error

// EventBus defines the interface for publishing domain events
type EventBus interface {
	// Publish publishes an event to all subscribers
	Publish(event DomainEvent) error
	// PublishAll publishes multiple events
	PublishAll(events []DomainEvent) error
	// Subscribe registers a handler for an event type
	Subscribe(eventType string, handler EventHandler)
}

// EventStore defines the interface for persisting events
type EventStore interface {
	// Append appends events to the store
	Append(aggregateID string, events []DomainEvent, expectedVersion int64) error
	// Load loads all events for an aggregate
	Load(aggregateID string) ([]DomainEvent, error)
	// LoadFrom loads events from a specific version
	LoadFrom(aggregateID string, fromVersion int64) ([]DomainEvent, error)
}
