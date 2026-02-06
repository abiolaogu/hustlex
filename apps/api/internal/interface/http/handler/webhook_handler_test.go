package handler_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hustlex/internal/domain/wallet/event"
	"hustlex/internal/domain/wallet/repository"
	"hustlex/internal/interface/http/handler"
	redisimpl "hustlex/internal/infrastructure/persistence/redis"
)

// MockWebhookEventStore is a mock implementation for testing
type MockWebhookEventStore struct {
	processedEvents map[event.WebhookEventID]bool
	markError       error
}

func NewMockWebhookEventStore() *MockWebhookEventStore {
	return &MockWebhookEventStore{
		processedEvents: make(map[event.WebhookEventID]bool),
	}
}

func (m *MockWebhookEventStore) IsProcessed(ctx context.Context, eventID event.WebhookEventID) (bool, error) {
	return m.processedEvents[eventID], nil
}

func (m *MockWebhookEventStore) MarkProcessed(ctx context.Context, webhookEvent *event.WebhookEvent) error {
	if m.markError != nil {
		return m.markError
	}
	if m.processedEvents[webhookEvent.EventID] {
		return redisimpl.ErrWebhookAlreadyProcessed
	}
	m.processedEvents[webhookEvent.EventID] = true
	return nil
}

func (m *MockWebhookEventStore) GetEvent(ctx context.Context, eventID event.WebhookEventID) (*event.WebhookEvent, error) {
	return nil, nil
}

func (m *MockWebhookEventStore) CleanupExpired(ctx context.Context, retentionPeriod time.Duration) error {
	return nil
}

var _ repository.WebhookEventStore = (*MockWebhookEventStore)(nil)

func TestWebhookHandler_SignatureVerification(t *testing.T) {
	secret := "test_secret"
	eventStore := NewMockWebhookEventStore()
	webhookHandler := handler.NewWebhookHandler(eventStore, secret)

	tests := []struct {
		name           string
		body           string
		secret         string
		wantStatusCode int
	}{
		{
			name:           "valid signature",
			body:           `{"event":"charge.success","data":{"reference":"ref123","status":"success"}}`,
			secret:         secret,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "invalid signature",
			body:           `{"event":"charge.success","data":{"reference":"ref456","status":"success"}}`,
			secret:         "wrong_secret",
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "missing signature",
			body:           `{"event":"charge.success","data":{"reference":"ref789","status":"success"}}`,
			secret:         "",
			wantStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/webhook/paystack", bytes.NewBufferString(tt.body))

			// Generate signature if secret is provided
			if tt.secret != "" {
				mac := hmac.New(sha512.New, []byte(tt.secret))
				mac.Write([]byte(tt.body))
				signature := hex.EncodeToString(mac.Sum(nil))
				req.Header.Set("X-Paystack-Signature", signature)
			}

			rec := httptest.NewRecorder()
			webhookHandler.HandlePaystackWebhook(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("HandlePaystackWebhook() status = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebhookHandler_Idempotency(t *testing.T) {
	secret := "test_secret"
	eventStore := NewMockWebhookEventStore()
	webhookHandler := handler.NewWebhookHandler(eventStore, secret)

	body := `{"event":"charge.success","data":{"reference":"ref_idempotent","status":"success","amount":1000}}`

	// Create valid signature
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(body))
	signature := hex.EncodeToString(mac.Sum(nil))

	// First request - should process successfully
	t.Run("first request processes successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/webhook/paystack", bytes.NewBufferString(body))
		req.Header.Set("X-Paystack-Signature", signature)

		rec := httptest.NewRecorder()
		webhookHandler.HandlePaystackWebhook(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("First request status = %v, want %v", rec.Code, http.StatusOK)
		}
	})

	// Second request - should detect duplicate and acknowledge
	t.Run("duplicate request is acknowledged", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/webhook/paystack", bytes.NewBufferString(body))
		req.Header.Set("X-Paystack-Signature", signature)

		rec := httptest.NewRecorder()
		webhookHandler.HandlePaystackWebhook(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Duplicate request status = %v, want %v", rec.Code, http.StatusOK)
		}

		// Verify response mentions it was already processed
		// In production, you'd parse the JSON and check the message
	})
}

func TestWebhookHandler_InvalidPayload(t *testing.T) {
	secret := "test_secret"
	eventStore := NewMockWebhookEventStore()
	webhookHandler := handler.NewWebhookHandler(eventStore, secret)

	tests := []struct {
		name           string
		body           string
		wantStatusCode int
	}{
		{
			name:           "malformed JSON",
			body:           `{"event":"charge.success","data":{invalid json}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "missing reference",
			body:           `{"event":"charge.success","data":{"status":"success"}}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/webhook/paystack", bytes.NewBufferString(tt.body))

			// Generate valid signature for the body
			mac := hmac.New(sha512.New, []byte(secret))
			mac.Write([]byte(tt.body))
			signature := hex.EncodeToString(mac.Sum(nil))
			req.Header.Set("X-Paystack-Signature", signature)

			rec := httptest.NewRecorder()
			webhookHandler.HandlePaystackWebhook(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("HandlePaystackWebhook() status = %v, want %v", rec.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebhookHandler_EventTypes(t *testing.T) {
	secret := "test_secret"
	eventStore := NewMockWebhookEventStore()
	webhookHandler := handler.NewWebhookHandler(eventStore, secret)

	tests := []struct {
		name      string
		eventType string
		reference string
	}{
		{
			name:      "charge.success",
			eventType: "charge.success",
			reference: "charge_ref",
		},
		{
			name:      "transfer.success",
			eventType: "transfer.success",
			reference: "transfer_ref",
		},
		{
			name:      "transfer.failed",
			eventType: "transfer.failed",
			reference: "failed_ref",
		},
		{
			name:      "transfer.reversed",
			eventType: "transfer.reversed",
			reference: "reversed_ref",
		},
		{
			name:      "unknown event",
			eventType: "unknown.event",
			reference: "unknown_ref",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := `{"event":"` + tt.eventType + `","data":{"reference":"` + tt.reference + `","status":"success"}}`

			req := httptest.NewRequest(http.MethodPost, "/api/webhook/paystack", bytes.NewBufferString(body))

			// Generate valid signature
			mac := hmac.New(sha512.New, []byte(secret))
			mac.Write([]byte(body))
			signature := hex.EncodeToString(mac.Sum(nil))
			req.Header.Set("X-Paystack-Signature", signature)

			rec := httptest.NewRecorder()
			webhookHandler.HandlePaystackWebhook(rec, req)

			// All events should be acknowledged successfully
			if rec.Code != http.StatusOK {
				t.Errorf("HandlePaystackWebhook() status = %v, want %v", rec.Code, http.StatusOK)
			}

			// Verify event was marked as processed
			eventID := event.WebhookEventID(tt.reference)
			processed, _ := eventStore.IsProcessed(context.Background(), eventID)
			if !processed {
				t.Errorf("Event %s was not marked as processed", tt.reference)
			}
		})
	}
}
