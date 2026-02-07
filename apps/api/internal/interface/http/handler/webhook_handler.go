package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"hustlex/internal/config"
	"hustlex/internal/infrastructure"
)

// WebhookHandler handles payment webhook callbacks with idempotency
type WebhookHandler struct {
	eventStore infrastructure.WebhookEventStore
	config     *config.Config
	// walletService WalletService // TODO: Add when wallet service is available
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(
	eventStore infrastructure.WebhookEventStore,
	cfg *config.Config,
) *WebhookHandler {
	return &WebhookHandler{
		eventStore: eventStore,
		config:     cfg,
	}
}

// PaystackWebhookEvent represents a Paystack webhook event
type PaystackWebhookEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// PaystackChargeData represents charge success data
type PaystackChargeData struct {
	ID              int64  `json:"id"`
	Reference       string `json:"reference"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Channel         string `json:"channel"`
	GatewayResponse string `json:"gateway_response"`
	PaidAt          string `json:"paid_at"`
	Metadata        struct {
		UserID string `json:"user_id"`
	} `json:"metadata"`
}

// PaystackTransferData represents transfer event data
type PaystackTransferData struct {
	ID        int64  `json:"id"`
	Reference string `json:"reference"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
}

// HandlePaystackWebhook handles Paystack payment webhooks with idempotency protection
// @Summary Handle Paystack webhook
// @Description Process payment notifications from Paystack (with duplicate detection)
// @Tags Webhooks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /webhooks/paystack [post]
func (h *WebhookHandler) HandlePaystackWebhook(c *fiber.Ctx) error {
	ctx := c.Context()

	// Read the raw body for signature verification
	body := c.Body()

	// Verify webhook signature - CRITICAL SECURITY CHECK
	signature := c.Get("X-Paystack-Signature")
	if signature == "" {
		log.Printf("[SECURITY] Webhook received without signature from IP: %s", c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Missing webhook signature",
		})
	}

	// Verify the signature using HMAC-SHA512
	if !h.verifyPaystackSignature(body, signature) {
		log.Printf("[SECURITY] Webhook signature verification failed from IP: %s", c.IP())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid webhook signature",
		})
	}

	// Parse the webhook event
	var event PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("[WEBHOOK] Failed to parse event: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Extract event reference for idempotency key
	eventReference, err := h.extractEventReference(event)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to extract event reference: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid event data",
		})
	}

	// IDEMPOTENCY CHECK: Check if event already processed
	isProcessed, err := h.eventStore.IsProcessed(ctx, eventReference)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to check idempotency for %s: %v", eventReference, err)
		// Continue processing despite error (fail open)
	}

	if isProcessed {
		// Duplicate webhook detected - acknowledge without reprocessing
		log.Printf("[WEBHOOK] Duplicate webhook detected: %s (event: %s) - acknowledging", eventReference, event.Event)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event already processed",
		})
	}

	// Log the event for audit trail (without sensitive data)
	log.Printf("[WEBHOOK] Processing new Paystack event: %s, reference: %s", event.Event, eventReference)

	// Process the webhook based on event type
	var processErr error
	switch event.Event {
	case "charge.success":
		processErr = h.handleChargeSuccess(ctx, event.Data, eventReference)
	case "transfer.success":
		processErr = h.handleTransferSuccess(ctx, event.Data, eventReference)
	case "transfer.failed":
		processErr = h.handleTransferFailed(ctx, event.Data, eventReference)
	case "transfer.reversed":
		processErr = h.handleTransferReversed(ctx, event.Data, eventReference)
	default:
		// Acknowledge unknown events to prevent retries
		log.Printf("[WEBHOOK] Unhandled event type: %s", event.Event)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event acknowledged",
		})
	}

	if processErr != nil {
		log.Printf("[WEBHOOK] Failed to process event %s: %v", eventReference, processErr)
		// Don't mark as processed if processing failed - allow retry
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to process webhook",
		})
	}

	// IDEMPOTENCY: Mark event as processed (7 day TTL covers typical retry windows)
	if err := h.eventStore.MarkProcessed(ctx, eventReference, 7*24*time.Hour); err != nil {
		log.Printf("[WEBHOOK] Warning: Failed to mark event as processed %s: %v", eventReference, err)
		// Continue - event was processed successfully
	}

	log.Printf("[WEBHOOK] Successfully processed event: %s", eventReference)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Webhook processed",
	})
}

// extractEventReference extracts a unique reference from webhook event for idempotency
func (h *WebhookHandler) extractEventReference(event PaystackWebhookEvent) (string, error) {
	switch event.Event {
	case "charge.success":
		var data PaystackChargeData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return "", err
		}
		// Use Paystack's unique ID + reference for idempotency key
		return fmt.Sprintf("charge:%d:%s", data.ID, data.Reference), nil

	case "transfer.success", "transfer.failed", "transfer.reversed":
		var data PaystackTransferData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return "", err
		}
		// Use transfer ID + reference
		return fmt.Sprintf("transfer:%d:%s", data.ID, data.Reference), nil

	default:
		// For unknown events, use event type + hash of data
		dataHash := sha512.Sum512(event.Data)
		return fmt.Sprintf("%s:%x", event.Event, dataHash[:16]), nil
	}
}

// verifyPaystackSignature verifies the webhook signature using HMAC-SHA512
func (h *WebhookHandler) verifyPaystackSignature(payload []byte, signature string) bool {
	secret := h.config.Payment.WebhookSecret
	if secret == "" {
		log.Printf("[SECURITY] Webhook secret not configured - rejecting webhook")
		return false
	}

	// Create HMAC-SHA512 hash
	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write(payload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// handleChargeSuccess processes successful payment charges
func (h *WebhookHandler) handleChargeSuccess(ctx context.Context, data json.RawMessage, eventRef string) error {
	var chargeData PaystackChargeData
	if err := json.Unmarshal(data, &chargeData); err != nil {
		return fmt.Errorf("failed to parse charge data: %w", err)
	}

	// Verify the charge status
	if chargeData.Status != "success" {
		log.Printf("[WEBHOOK] Charge not successful: %s, status: %s", chargeData.Reference, chargeData.Status)
		return nil // Don't retry for non-success status
	}

	log.Printf("[WEBHOOK] Processing successful charge: %s, amount: %d %s, user: %s",
		chargeData.Reference, chargeData.Amount, chargeData.Currency, chargeData.Metadata.UserID)

	// TODO: Call wallet service to credit user's wallet
	// Example: h.walletService.ProcessDeposit(ctx, chargeData.Reference, chargeData.Amount, chargeData.Metadata.UserID)

	return nil
}

// handleTransferSuccess processes successful transfers (withdrawals)
func (h *WebhookHandler) handleTransferSuccess(ctx context.Context, data json.RawMessage, eventRef string) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		return fmt.Errorf("failed to parse transfer data: %w", err)
	}

	log.Printf("[WEBHOOK] Transfer successful: %s, amount: %d", transferData.Reference, transferData.Amount)

	// TODO: Update transaction status to completed
	// Example: h.walletService.CompleteWithdrawal(ctx, transferData.Reference)

	return nil
}

// handleTransferFailed processes failed transfers
func (h *WebhookHandler) handleTransferFailed(ctx context.Context, data json.RawMessage, eventRef string) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		return fmt.Errorf("failed to parse transfer data: %w", err)
	}

	log.Printf("[WEBHOOK] Transfer failed: %s, reason: %s", transferData.Reference, transferData.Reason)

	// TODO: Reverse the withdrawal and refund user
	// Example: h.walletService.FailWithdrawal(ctx, transferData.Reference, transferData.Reason)

	return nil
}

// handleTransferReversed processes reversed transfers
func (h *WebhookHandler) handleTransferReversed(ctx context.Context, data json.RawMessage, eventRef string) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		return fmt.Errorf("failed to parse transfer data: %w", err)
	}

	log.Printf("[WEBHOOK] Transfer reversed: %s", transferData.Reference)

	// TODO: Handle transfer reversal
	// Example: h.walletService.ReverseWithdrawal(ctx, transferData.Reference)

	return nil
}

// IPWhitelistMiddleware restricts webhook access to known Paystack IPs
func IPWhitelistMiddleware(allowedIPs []string) fiber.Handler {
	ipSet := make(map[string]bool)
	for _, ip := range allowedIPs {
		ipSet[ip] = true
	}

	return func(c *fiber.Ctx) error {
		clientIP := c.IP()

		// In production, verify against Paystack's IP whitelist
		if len(ipSet) > 0 && !ipSet[clientIP] {
			log.Printf("[SECURITY] Webhook request from unauthorized IP: %s", clientIP)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Forbidden",
			})
		}

		return c.Next()
	}
}

// PaystackWebhookIPs returns known Paystack webhook IPs
// Source: https://paystack.com/docs/payments/webhooks/#ip-whitelisting
func PaystackWebhookIPs() []string {
	return []string{
		"52.31.139.75",
		"52.49.173.169",
		"52.214.14.220",
	}
}
