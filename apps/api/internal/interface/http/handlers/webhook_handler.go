package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"hustlex/internal/services"
)

// WebhookHandler handles payment webhook callbacks with idempotency protection
type WebhookHandler struct {
	eventStore    services.WebhookEventStore
	webhookSecret string
}

// NewWebhookHandler creates a new webhook handler with idempotency support
func NewWebhookHandler(eventStore services.WebhookEventStore, webhookSecret string) *WebhookHandler {
	return &WebhookHandler{
		eventStore:    eventStore,
		webhookSecret: webhookSecret,
	}
}

// PaystackWebhookEvent represents a Paystack webhook event
type PaystackWebhookEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// PaystackChargeData represents charge success data
type PaystackChargeData struct {
	Reference       string `json:"reference"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Channel         string `json:"channel"`
	GatewayResponse string `json:"gateway_response"`
	Metadata        struct {
		UserID string `json:"user_id"`
	} `json:"metadata"`
}

// PaystackTransferData represents transfer event data
type PaystackTransferData struct {
	Reference string `json:"reference"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

// HandlePaystackWebhook handles Paystack payment webhooks with idempotency
func (h *WebhookHandler) HandlePaystackWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read the raw body for signature verification
	body := make([]byte, 0)
	if r.Body != nil {
		defer r.Body.Close()
		var err error
		body, err = http.ReadAll(r.Body)
		if err != nil {
			log.Printf("[WEBHOOK] Failed to read body: %v", err)
			h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
				"success": false,
				"error":   "Failed to read request body",
			})
			return
		}
	}

	// Verify webhook signature - CRITICAL SECURITY CHECK
	signature := r.Header.Get("X-Paystack-Signature")
	if signature == "" {
		log.Printf("[SECURITY] Webhook received without signature")
		h.respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Missing webhook signature",
		})
		return
	}

	// Verify the signature using HMAC-SHA512
	if !h.verifyPaystackSignature(body, signature) {
		log.Printf("[SECURITY] Webhook signature verification failed")
		h.respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "Invalid webhook signature",
		})
		return
	}

	// Parse the webhook event
	var event PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("[WEBHOOK] Failed to parse event: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Extract event reference for idempotency
	reference, err := h.extractReference(event)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to extract reference from event: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid event data",
		})
		return
	}

	// IDEMPOTENCY CHECK: Check if already processed
	processed, err := h.eventStore.IsProcessed(ctx, reference)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to check if event processed: %v", err)
		// Don't fail the webhook - acknowledge it to prevent retries
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Event acknowledged (idempotency check failed)",
		})
		return
	}

	if processed {
		// Event already processed - acknowledge duplicate
		log.Printf("[WEBHOOK] Duplicate event detected: %s (already processed)", reference)
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Event already processed",
		})
		return
	}

	// Log the event for audit trail (without sensitive data)
	log.Printf("[WEBHOOK] Processing new Paystack event: %s, reference: %s", event.Event, reference)

	// Handle different event types
	switch event.Event {
	case "charge.success":
		h.handleChargeSuccess(w, r, event.Data, reference)
	case "transfer.success":
		h.handleTransferSuccess(w, r, event.Data, reference)
	case "transfer.failed":
		h.handleTransferFailed(w, r, event.Data, reference)
	case "transfer.reversed":
		h.handleTransferReversed(w, r, event.Data, reference)
	default:
		// Acknowledge unknown events to prevent retries
		log.Printf("[WEBHOOK] Unhandled event type: %s", event.Event)
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Event acknowledged",
		})
	}
}

// extractReference extracts the transaction reference from the webhook event
func (h *WebhookHandler) extractReference(event PaystackWebhookEvent) (string, error) {
	var data struct {
		Reference string `json:"reference"`
	}

	if err := json.Unmarshal(event.Data, &data); err != nil {
		return "", err
	}

	if data.Reference == "" {
		return "", ErrMissingReference
	}

	return data.Reference, nil
}

// verifyPaystackSignature verifies the webhook signature using HMAC-SHA512
func (h *WebhookHandler) verifyPaystackSignature(payload []byte, signature string) bool {
	if h.webhookSecret == "" {
		log.Printf("[SECURITY] Webhook secret not configured")
		return false
	}

	// Create HMAC-SHA512 hash
	mac := hmac.New(sha512.New, []byte(h.webhookSecret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// handleChargeSuccess processes successful payment charges
func (h *WebhookHandler) handleChargeSuccess(w http.ResponseWriter, r *http.Request, data json.RawMessage, reference string) {
	ctx := r.Context()

	var chargeData PaystackChargeData
	if err := json.Unmarshal(data, &chargeData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse charge data: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid charge data",
		})
		return
	}

	// Verify the charge status
	if chargeData.Status != "success" {
		log.Printf("[WEBHOOK] Charge not successful: %s", reference)
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Charge status not successful, ignoring",
		})
		return
	}

	log.Printf("[WEBHOOK] Processing successful charge: %s, amount: %d", reference, chargeData.Amount)

	// TODO: Call wallet service to credit user's wallet
	// Example: h.walletService.ProcessDeposit(ctx, chargeData.Reference)

	// IMPORTANT: Mark as processed AFTER successful processing
	if err := h.eventStore.MarkProcessed(ctx, reference, services.PaystackRetryWindow); err != nil {
		log.Printf("[WEBHOOK] Failed to mark event as processed: %v", err)
		// Don't fail the webhook - the business logic succeeded
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Charge processed",
	})
}

// handleTransferSuccess processes successful transfers (withdrawals)
func (h *WebhookHandler) handleTransferSuccess(w http.ResponseWriter, r *http.Request, data json.RawMessage, reference string) {
	ctx := r.Context()

	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid transfer data",
		})
		return
	}

	log.Printf("[WEBHOOK] Transfer successful: %s", reference)

	// TODO: Update transaction status to completed
	// Example: h.walletService.CompleteWithdrawal(ctx, transferData.Reference)

	// Mark as processed
	if err := h.eventStore.MarkProcessed(ctx, reference, services.PaystackRetryWindow); err != nil {
		log.Printf("[WEBHOOK] Failed to mark event as processed: %v", err)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Transfer processed",
	})
}

// handleTransferFailed processes failed transfers
func (h *WebhookHandler) handleTransferFailed(w http.ResponseWriter, r *http.Request, data json.RawMessage, reference string) {
	ctx := r.Context()

	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid transfer data",
		})
		return
	}

	log.Printf("[WEBHOOK] Transfer failed: %s, reason: %s", reference, transferData.Reason)

	// TODO: Reverse the withdrawal and refund user
	// Example: h.walletService.FailWithdrawal(ctx, transferData.Reference, transferData.Reason)

	// Mark as processed
	if err := h.eventStore.MarkProcessed(ctx, reference, services.PaystackRetryWindow); err != nil {
		log.Printf("[WEBHOOK] Failed to mark event as processed: %v", err)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Transfer failure processed",
	})
}

// handleTransferReversed processes reversed transfers
func (h *WebhookHandler) handleTransferReversed(w http.ResponseWriter, r *http.Request, data json.RawMessage, reference string) {
	ctx := r.Context()

	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid transfer data",
		})
		return
	}

	log.Printf("[WEBHOOK] Transfer reversed: %s", reference)

	// TODO: Handle transfer reversal
	// Example: h.walletService.ReverseWithdrawal(ctx, transferData.Reference)

	// Mark as processed
	if err := h.eventStore.MarkProcessed(ctx, reference, services.PaystackRetryWindow); err != nil {
		log.Printf("[WEBHOOK] Failed to mark event as processed: %v", err)
	}

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Transfer reversal processed",
	})
}

// respondJSON sends a JSON response
func (h *WebhookHandler) respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Error definitions
var (
	ErrMissingReference = fmt.Errorf("missing transaction reference in webhook event")
)
