package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"hustlex/internal/domain/wallet/event"
	"hustlex/internal/domain/wallet/repository"
	redisimpl "hustlex/internal/infrastructure/persistence/redis"
	"hustlex/internal/interface/http/response"
)

// WebhookHandler handles payment webhook callbacks
type WebhookHandler struct {
	eventStore    repository.WebhookEventStore
	webhookSecret string
	// walletService will be added when wallet service is implemented
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(eventStore repository.WebhookEventStore, webhookSecret string) *WebhookHandler {
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

// HandlePaystackWebhook handles Paystack payment webhooks with idempotency protection
func (h *WebhookHandler) HandlePaystackWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the raw body for signature verification
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to read body: %v", err)
		response.BadRequest(w, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	// Verify webhook signature - CRITICAL SECURITY CHECK
	signature := r.Header.Get("X-Paystack-Signature")
	if signature == "" {
		log.Printf("[SECURITY] Webhook received without signature")
		response.Unauthorized(w, "Missing webhook signature")
		return
	}

	if !h.verifyPaystackSignature(body, signature) {
		log.Printf("[SECURITY] Webhook signature verification failed")
		response.Unauthorized(w, "Invalid webhook signature")
		return
	}

	// Parse the webhook event
	var webhookEvent PaystackWebhookEvent
	if err := json.Unmarshal(body, &webhookEvent); err != nil {
		log.Printf("[WEBHOOK] Failed to parse event: %v", err)
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Extract reference from event data for idempotency check
	reference, err := h.extractReference(webhookEvent.Data)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to extract reference: %v", err)
		response.BadRequest(w, "Missing reference in webhook data")
		return
	}

	// Check idempotency - has this webhook already been processed?
	eventID := event.WebhookEventID(reference)
	processed, err := h.eventStore.IsProcessed(r.Context(), eventID)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to check idempotency: %v", err)
		response.InternalError(w, "Failed to process webhook")
		return
	}

	if processed {
		// This webhook has already been processed - acknowledge it but don't reprocess
		log.Printf("[WEBHOOK] Duplicate webhook detected: %s (already processed)", reference)
		response.Success(w, map[string]interface{}{
			"success": true,
			"message": "Webhook already processed",
		})
		return
	}

	// Create webhook event record
	webhookEventRecord := event.NewWebhookEvent("paystack", webhookEvent.Event, reference, body)

	// Mark as processed BEFORE actual processing to prevent race conditions
	// This ensures that if multiple identical webhooks arrive simultaneously,
	// only one will be processed (SetNX provides atomicity)
	err = h.eventStore.MarkProcessed(r.Context(), webhookEventRecord)
	if err != nil {
		if errors.Is(err, redisimpl.ErrWebhookAlreadyProcessed) {
			// Another goroutine/process beat us to it
			log.Printf("[WEBHOOK] Race condition detected: %s (processed by another instance)", reference)
			response.Success(w, map[string]interface{}{
				"success": true,
				"message": "Webhook already processed",
			})
			return
		}
		log.Printf("[WEBHOOK] Failed to mark webhook as processed: %v", err)
		response.InternalError(w, "Failed to process webhook")
		return
	}

	// Log the event for audit trail (without sensitive data)
	log.Printf("[WEBHOOK] Processing Paystack event: %s, reference: %s", webhookEvent.Event, reference)

	// Handle different event types
	switch webhookEvent.Event {
	case "charge.success":
		h.handleChargeSuccess(r.Context(), webhookEvent.Data)
	case "transfer.success":
		h.handleTransferSuccess(r.Context(), webhookEvent.Data)
	case "transfer.failed":
		h.handleTransferFailed(r.Context(), webhookEvent.Data)
	case "transfer.reversed":
		h.handleTransferReversed(r.Context(), webhookEvent.Data)
	default:
		// Acknowledge unknown events to prevent retries
		log.Printf("[WEBHOOK] Unhandled event type: %s", webhookEvent.Event)
	}

	response.Success(w, map[string]interface{}{
		"success": true,
		"message": "Webhook processed successfully",
	})
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

// extractReference extracts the payment reference from webhook data
func (h *WebhookHandler) extractReference(data json.RawMessage) (string, error) {
	var refData struct {
		Reference string `json:"reference"`
	}

	if err := json.Unmarshal(data, &refData); err != nil {
		return "", err
	}

	if refData.Reference == "" {
		return "", errors.New("reference not found in webhook data")
	}

	return refData.Reference, nil
}

// handleChargeSuccess processes successful payment charges
func (h *WebhookHandler) handleChargeSuccess(ctx context.Context, data json.RawMessage) {
	var chargeData PaystackChargeData
	if err := json.Unmarshal(data, &chargeData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse charge data: %v", err)
		return
	}

	// Verify the charge status
	if chargeData.Status != "success" {
		log.Printf("[WEBHOOK] Charge not successful: %s", chargeData.Reference)
		return
	}

	log.Printf("[WEBHOOK] Processing successful charge: %s, amount: %d",
		chargeData.Reference, chargeData.Amount)

	// TODO: Call wallet service to credit user's wallet
	// Example: h.walletService.ProcessDeposit(ctx, chargeData.Reference)
}

// handleTransferSuccess processes successful transfers (withdrawals)
func (h *WebhookHandler) handleTransferSuccess(ctx context.Context, data json.RawMessage) {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return
	}

	log.Printf("[WEBHOOK] Transfer successful: %s", transferData.Reference)

	// TODO: Update transaction status to completed
	// Example: h.walletService.CompleteWithdrawal(ctx, transferData.Reference)
}

// handleTransferFailed processes failed transfers
func (h *WebhookHandler) handleTransferFailed(ctx context.Context, data json.RawMessage) {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return
	}

	log.Printf("[WEBHOOK] Transfer failed: %s, reason: %s",
		transferData.Reference, transferData.Reason)

	// TODO: Reverse the withdrawal and refund user
	// Example: h.walletService.FailWithdrawal(ctx, transferData.Reference, transferData.Reason)
}

// handleTransferReversed processes reversed transfers
func (h *WebhookHandler) handleTransferReversed(ctx context.Context, data json.RawMessage) {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return
	}

	log.Printf("[WEBHOOK] Transfer reversed: %s", transferData.Reference)

	// TODO: Handle transfer reversal
	// Example: h.walletService.ReverseWithdrawal(ctx, transferData.Reference)
}
