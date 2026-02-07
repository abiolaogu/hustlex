package event

import (
	"time"
)

// WebhookEventID represents a unique webhook event identifier
type WebhookEventID string

// WebhookEvent represents a payment webhook event
type WebhookEvent struct {
	EventID     WebhookEventID
	Provider    string // "paystack", "flutterwave", etc.
	EventType   string // "charge.success", "transfer.success", etc.
	Reference   string // Payment reference
	ProcessedAt time.Time
	Payload     []byte // Raw webhook payload for audit
}

// NewWebhookEvent creates a new webhook event
func NewWebhookEvent(provider, eventType, reference string, payload []byte) *WebhookEvent {
	return &WebhookEvent{
		EventID:     WebhookEventID(reference), // Use reference as event ID for idempotency
		Provider:    provider,
		EventType:   eventType,
		Reference:   reference,
		ProcessedAt: time.Now().UTC(),
		Payload:     payload,
	}
}

// WebhookProcessed represents a domain event when a webhook is processed
type WebhookProcessed struct {
	EventID     WebhookEventID
	Provider    string
	EventType   string
	Reference   string
	ProcessedAt time.Time
}
