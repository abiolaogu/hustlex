package handlers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hustlex/internal/services"
)

// MockWebhookEventStore is a mock implementation of WebhookEventStore for testing
type MockWebhookEventStore struct {
	processed map[string]bool
	times     map[string]time.Time
}

func NewMockWebhookEventStore() *MockWebhookEventStore {
	return &MockWebhookEventStore{
		processed: make(map[string]bool),
		times:     make(map[string]time.Time),
	}
}

func (m *MockWebhookEventStore) IsProcessed(ctx context.Context, eventID string) (bool, error) {
	return m.processed[eventID], nil
}

func (m *MockWebhookEventStore) MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error {
	m.processed[eventID] = true
	m.times[eventID] = time.Now()
	return nil
}

func (m *MockWebhookEventStore) GetProcessedAt(ctx context.Context, eventID string) (*time.Time, error) {
	if t, ok := m.times[eventID]; ok {
		return &t, nil
	}
	return nil, nil
}

func TestWebhookHandler_HandlePaystackWebhook(t *testing.T) {
	webhookSecret := "test_secret_key"
	eventStore := NewMockWebhookEventStore()
	handler := NewWebhookHandler(eventStore, webhookSecret)

	t.Run("rejects webhook without signature", func(t *testing.T) {
		body := []byte(`{"event":"charge.success","data":{"reference":"trx_123"}}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"], "signature")
	})

	t.Run("rejects webhook with invalid signature", func(t *testing.T) {
		body := []byte(`{"event":"charge.success","data":{"reference":"trx_123"}}`)
		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", "invalid_signature")
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"], "Invalid")
	})

	t.Run("processes valid charge.success webhook", func(t *testing.T) {
		payload := map[string]interface{}{
			"event": "charge.success",
			"data": map[string]interface{}{
				"reference": "trx_test_001",
				"amount":    5000000, // 50,000 NGN in kobo
				"currency":  "NGN",
				"status":    "success",
				"metadata": map[string]interface{}{
					"user_id": "user_123",
				},
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		// Verify event was marked as processed
		processed, _ := eventStore.IsProcessed(context.Background(), "trx_test_001")
		assert.True(t, processed)
	})

	t.Run("rejects duplicate webhook delivery", func(t *testing.T) {
		reference := "trx_duplicate_test"

		// Mark as already processed
		eventStore.MarkProcessed(context.Background(), reference, time.Hour)

		payload := map[string]interface{}{
			"event": "charge.success",
			"data": map[string]interface{}{
				"reference": reference,
				"amount":    1000000,
				"currency":  "NGN",
				"status":    "success",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Contains(t, response["message"], "already processed")
	})

	t.Run("handles transfer.success webhook", func(t *testing.T) {
		payload := map[string]interface{}{
			"event": "transfer.success",
			"data": map[string]interface{}{
				"reference": "trx_transfer_001",
				"amount":    2000000,
				"status":    "success",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		// Verify event was marked as processed
		processed, _ := eventStore.IsProcessed(context.Background(), "trx_transfer_001")
		assert.True(t, processed)
	})

	t.Run("handles transfer.failed webhook", func(t *testing.T) {
		payload := map[string]interface{}{
			"event": "transfer.failed",
			"data": map[string]interface{}{
				"reference": "trx_transfer_failed_001",
				"amount":    1500000,
				"status":    "failed",
				"reason":    "Insufficient balance",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		// Verify event was marked as processed
		processed, _ := eventStore.IsProcessed(context.Background(), "trx_transfer_failed_001")
		assert.True(t, processed)
	})

	t.Run("acknowledges unknown event types", func(t *testing.T) {
		payload := map[string]interface{}{
			"event": "unknown.event",
			"data": map[string]interface{}{
				"reference": "trx_unknown_001",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Contains(t, response["message"], "acknowledged")
	})

	t.Run("rejects malformed JSON", func(t *testing.T) {
		body := []byte(`{invalid json}`)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("rejects event without reference", func(t *testing.T) {
		payload := map[string]interface{}{
			"event": "charge.success",
			"data": map[string]interface{}{
				// Missing reference field
				"amount":   1000000,
				"currency": "NGN",
				"status":   "success",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		req := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req.Header.Set("X-Paystack-Signature", signature)
		w := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWebhookHandler_VerifySignature(t *testing.T) {
	webhookSecret := "test_webhook_secret"
	eventStore := NewMockWebhookEventStore()
	handler := NewWebhookHandler(eventStore, webhookSecret)

	t.Run("valid signature passes verification", func(t *testing.T) {
		payload := []byte(`{"event":"test","data":{"reference":"test_123"}}`)
		signature := generateSignature(payload, webhookSecret)

		valid := handler.verifyPaystackSignature(payload, signature)
		assert.True(t, valid)
	})

	t.Run("invalid signature fails verification", func(t *testing.T) {
		payload := []byte(`{"event":"test","data":{"reference":"test_123"}}`)
		invalidSignature := "wrong_signature"

		valid := handler.verifyPaystackSignature(payload, invalidSignature)
		assert.False(t, valid)
	})

	t.Run("tampered payload fails verification", func(t *testing.T) {
		originalPayload := []byte(`{"event":"test","amount":1000}`)
		signature := generateSignature(originalPayload, webhookSecret)

		tamperedPayload := []byte(`{"event":"test","amount":9999}`)

		valid := handler.verifyPaystackSignature(tamperedPayload, signature)
		assert.False(t, valid)
	})

	t.Run("empty secret fails verification", func(t *testing.T) {
		emptySecretHandler := NewWebhookHandler(eventStore, "")
		payload := []byte(`{"event":"test"}`)
		signature := "any_signature"

		valid := emptySecretHandler.verifyPaystackSignature(payload, signature)
		assert.False(t, valid)
	})
}

func TestWebhookHandler_ExtractReference(t *testing.T) {
	eventStore := NewMockWebhookEventStore()
	handler := NewWebhookHandler(eventStore, "secret")

	tests := []struct {
		name        string
		event       PaystackWebhookEvent
		expected    string
		shouldError bool
	}{
		{
			name: "extracts reference from charge event",
			event: PaystackWebhookEvent{
				Event: "charge.success",
				Data:  json.RawMessage(`{"reference":"trx_abc123","amount":1000}`),
			},
			expected:    "trx_abc123",
			shouldError: false,
		},
		{
			name: "extracts reference from transfer event",
			event: PaystackWebhookEvent{
				Event: "transfer.success",
				Data:  json.RawMessage(`{"reference":"trx_xyz789","status":"success"}`),
			},
			expected:    "trx_xyz789",
			shouldError: false,
		},
		{
			name: "fails when reference is missing",
			event: PaystackWebhookEvent{
				Event: "charge.success",
				Data:  json.RawMessage(`{"amount":1000}`),
			},
			expected:    "",
			shouldError: true,
		},
		{
			name: "fails with malformed data",
			event: PaystackWebhookEvent{
				Event: "charge.success",
				Data:  json.RawMessage(`{invalid json}`),
			},
			expected:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reference, err := handler.extractReference(tt.event)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, reference)
			}
		})
	}
}

func TestWebhookHandler_RealWorldScenario(t *testing.T) {
	webhookSecret := "sk_test_abc123"
	eventStore := NewMockWebhookEventStore()
	handler := NewWebhookHandler(eventStore, webhookSecret)

	t.Run("simulate Paystack retry behavior", func(t *testing.T) {
		reference := "trx_real_world_001"

		// First delivery
		payload := map[string]interface{}{
			"event": "charge.success",
			"data": map[string]interface{}{
				"reference": reference,
				"amount":    5000000,
				"currency":  "NGN",
				"status":    "success",
			},
		}

		body, _ := json.Marshal(payload)
		signature := generateSignature(body, webhookSecret)

		// First attempt - should succeed
		req1 := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req1.Header.Set("X-Paystack-Signature", signature)
		w1 := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w1, req1)
		assert.Equal(t, http.StatusOK, w1.Code)

		var resp1 map[string]interface{}
		json.Unmarshal(w1.Body.Bytes(), &resp1)
		assert.True(t, resp1["success"].(bool))
		assert.NotContains(t, resp1["message"], "already processed")

		// Second attempt (retry after 1 minute) - should be rejected
		req2 := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req2.Header.Set("X-Paystack-Signature", signature)
		w2 := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w2, req2)
		assert.Equal(t, http.StatusOK, w2.Code)

		var resp2 map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &resp2)
		assert.True(t, resp2["success"].(bool))
		assert.Contains(t, resp2["message"], "already processed")

		// Third attempt (retry after 1 hour) - should be rejected
		req3 := httptest.NewRequest(http.MethodPost, "/webhooks/paystack", bytes.NewReader(body))
		req3.Header.Set("X-Paystack-Signature", signature)
		w3 := httptest.NewRecorder()

		handler.HandlePaystackWebhook(w3, req3)
		assert.Equal(t, http.StatusOK, w3.Code)

		var resp3 map[string]interface{}
		json.Unmarshal(w3.Body.Bytes(), &resp3)
		assert.True(t, resp3["success"].(bool))
		assert.Contains(t, resp3["message"], "already processed")
	})
}

// generateSignature creates an HMAC-SHA512 signature for testing
func generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
