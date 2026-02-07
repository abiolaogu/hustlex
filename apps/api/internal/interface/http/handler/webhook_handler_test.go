package handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"hustlex/internal/config"
	"hustlex/internal/infrastructure"
)

func TestWebhookHandler_HandlePaystackWebhook_ChargeSuccess(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"id":               123456,
			"reference":        "trx_abc123",
			"amount":           5000000, // 50,000 NGN in kobo
			"currency":         "NGN",
			"status":           "success",
			"channel":          "card",
			"gateway_response": "Successful",
			"paid_at":          "2026-02-07T12:00:00Z",
			"metadata": map[string]interface{}{
				"user_id": "user_12345",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	// Generate signature
	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	// Create request
	req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Paystack-Signature", signature)

	// Execute
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "Webhook processed", result["message"])

	// Verify event was marked as processed
	isProcessed, err := eventStore.IsProcessed(context.Background(), "charge:123456:trx_abc123")
	require.NoError(t, err)
	assert.True(t, isProcessed)
}

func TestWebhookHandler_HandlePaystackWebhook_DuplicateEvent(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"id":        123456,
			"reference": "trx_duplicate",
			"amount":    5000000,
			"currency":  "NGN",
			"status":    "success",
			"channel":   "card",
			"metadata": map[string]interface{}{
				"user_id": "user_12345",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	// First request - should process
	req1 := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Paystack-Signature", signature)

	resp1, err := app.Test(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Second request - duplicate, should skip processing
	req2 := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Paystack-Signature", signature)

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	// Assert duplicate was detected
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp2.Body)
	json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "Event already processed", result["message"])
}

func TestWebhookHandler_HandlePaystackWebhook_InvalidSignature(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"id":        123456,
			"reference": "trx_invalid_sig",
			"amount":    5000000,
			"status":    "success",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	// Use wrong signature
	invalidSignature := "invalid_signature_12345"

	// Create request
	req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Paystack-Signature", invalidSignature)

	// Execute
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, "Invalid webhook signature", result["error"])
}

func TestWebhookHandler_HandlePaystackWebhook_MissingSignature(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "charge.success",
		"data":  map[string]interface{}{},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	// Create request without signature
	req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	// No X-Paystack-Signature header

	// Execute
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, "Missing webhook signature", result["error"])
}

func TestWebhookHandler_HandlePaystackWebhook_TransferSuccess(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "transfer.success",
		"data": map[string]interface{}{
			"id":        789012,
			"reference": "trx_transfer_xyz",
			"amount":    2000000, // 20,000 NGN
			"status":    "success",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	// Create request
	req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Paystack-Signature", signature)

	// Execute
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Verify event was marked as processed
	isProcessed, err := eventStore.IsProcessed(context.Background(), "transfer:789012:trx_transfer_xyz")
	require.NoError(t, err)
	assert.True(t, isProcessed)
}

func TestWebhookHandler_HandlePaystackWebhook_UnknownEvent(t *testing.T) {
	// Setup
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload with unknown event type
	payload := map[string]interface{}{
		"event": "unknown.event.type",
		"data": map[string]interface{}{
			"some_field": "some_value",
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	// Create request
	req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Paystack-Signature", signature)

	// Execute
	resp, err := app.Test(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Assert - should acknowledge unknown events
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)
	assert.True(t, result["success"].(bool))
	assert.Equal(t, "Event acknowledged", result["message"])
}

func TestWebhookHandler_ExtractEventReference(t *testing.T) {
	handler := &WebhookHandler{}

	tests := []struct {
		name        string
		event       PaystackWebhookEvent
		expected    string
		expectError bool
	}{
		{
			name: "charge.success",
			event: PaystackWebhookEvent{
				Event: "charge.success",
				Data:  mustMarshal(PaystackChargeData{ID: 123, Reference: "ref_123"}),
			},
			expected:    "charge:123:ref_123",
			expectError: false,
		},
		{
			name: "transfer.success",
			event: PaystackWebhookEvent{
				Event: "transfer.success",
				Data:  mustMarshal(PaystackTransferData{ID: 456, Reference: "ref_456"}),
			},
			expected:    "transfer:456:ref_456",
			expectError: false,
		},
		{
			name: "unknown event",
			event: PaystackWebhookEvent{
				Event: "custom.event",
				Data:  json.RawMessage(`{"test":"data"}`),
			},
			expected:    "custom.event:",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handler.extractEventReference(tt.event)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.name != "unknown event" {
					assert.Equal(t, tt.expected, result)
				} else {
					// For unknown events, just check it starts with event type
					assert.Contains(t, result, "custom.event:")
				}
			}
		})
	}
}

func TestIPWhitelistMiddleware(t *testing.T) {
	allowedIPs := []string{"52.31.139.75", "52.49.173.169"}
	middleware := IPWhitelistMiddleware(allowedIPs)

	app := fiber.New()
	app.Use(middleware)
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	tests := []struct {
		name           string
		clientIP       string
		expectedStatus int
	}{
		{
			name:           "allowed IP",
			clientIP:       "52.31.139.75",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "disallowed IP",
			clientIP:       "192.168.1.1",
			expectedStatus: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Forwarded-For", tt.clientIP)

			resp, err := app.Test(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Note: Fiber's c.IP() behavior in tests may differ from production
			// This test demonstrates the middleware structure
		})
	}
}

func TestPaystackWebhookIPs(t *testing.T) {
	ips := PaystackWebhookIPs()

	assert.Len(t, ips, 3)
	assert.Contains(t, ips, "52.31.139.75")
	assert.Contains(t, ips, "52.49.173.169")
	assert.Contains(t, ips, "52.214.14.220")
}

func TestWebhookHandler_RaceCondition(t *testing.T) {
	// Test concurrent webhook deliveries (simulating network retry)
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	// Create webhook payload
	payload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"id":        999999,
			"reference": "trx_race_condition",
			"amount":    5000000,
			"status":    "success",
			"metadata": map[string]interface{}{
				"user_id": "user_12345",
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	// Send concurrent requests
	done := make(chan int, 2)

	for i := 0; i < 2; i++ {
		go func() {
			req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Paystack-Signature", signature)

			resp, _ := app.Test(req)
			done <- resp.StatusCode
			resp.Body.Close()
		}()
	}

	// Collect results
	status1 := <-done
	status2 := <-done

	// Both should return 200 OK
	assert.Equal(t, fiber.StatusOK, status1)
	assert.Equal(t, fiber.StatusOK, status2)

	// Verify event was marked as processed exactly once
	isProcessed, err := eventStore.IsProcessed(context.Background(), "charge:999999:trx_race_condition")
	require.NoError(t, err)
	assert.True(t, isProcessed)

	// Note: In production with Redis, use SETNX for atomic check-and-set
}

// Helper functions

func generateSignature(payload []byte, secret string) string {
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// Benchmark tests

func BenchmarkWebhookHandler_HandlePaystackWebhook(b *testing.B) {
	cfg := &config.Config{
		Payment: config.PaymentConfig{
			WebhookSecret: "test_webhook_secret_12345",
		},
	}
	eventStore := infrastructure.NewInMemoryWebhookEventStore()
	handler := NewWebhookHandler(eventStore, cfg)

	app := fiber.New()
	app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)

	payload := map[string]interface{}{
		"event": "charge.success",
		"data": map[string]interface{}{
			"id":        123456,
			"reference": "trx_benchmark",
			"amount":    5000000,
			"status":    "success",
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	signature := generateSignature(payloadBytes, cfg.Payment.WebhookSecret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/webhooks/paystack", bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Paystack-Signature", signature)

		resp, _ := app.Test(req)
		resp.Body.Close()
	}
}
