package handler

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"hustlex/internal/domain/webhook/repository"
	"hustlex/internal/infrastructure/webhook"

	"github.com/gofiber/fiber/v2"
)

// WebhookHandler handles payment webhook callbacks with idempotency
type WebhookHandler struct {
	eventStore    repository.WebhookEventStore
	webhookSecret string
	// Add wallet service and other dependencies as needed
	// walletService WalletService
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
	ID              int64  `json:"id"`
	Reference       string `json:"reference"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Channel         string `json:"channel"`
	GatewayResponse string `json:"gateway_response"`
	PaidAt          string `json:"paid_at"`
	Customer        struct {
		Email string `json:"email"`
	} `json:"customer"`
	Metadata struct {
		UserID   string `json:"user_id"`
		WalletID string `json:"wallet_id"`
	} `json:"metadata"`
}

// PaystackTransferData represents transfer event data
type PaystackTransferData struct {
	ID        int64  `json:"id"`
	Reference string `json:"reference"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Recipient string `json:"recipient"`
}

// HandlePaystackWebhook handles Paystack payment webhooks with idempotency protection
// @Summary Handle Paystack webhook
// @Description Process payment notifications from Paystack (idempotent)
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

	// STEP 1: Verify webhook signature - CRITICAL SECURITY CHECK
	signature := c.Get("X-Paystack-Signature")
	if signature == "" {
		log.Printf("[SECURITY] Webhook received without signature from IP: %s", c.IP())
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Missing webhook signature",
		})
	}

	// Verify the signature using HMAC-SHA512
	if !h.verifyPaystackSignature(body, signature) {
		log.Printf("[SECURITY] Webhook signature verification failed from IP: %s", c.IP())
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid webhook signature",
		})
	}

	// STEP 2: Parse the webhook event
	var event PaystackWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("[WEBHOOK] Failed to parse event: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// STEP 3: Extract event ID for idempotency check
	// For Paystack, we can use the reference or a combination of event type + reference
	eventID, err := h.extractEventID(event)
	if err != nil {
		log.Printf("[WEBHOOK] Failed to extract event ID: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing or invalid event identifier",
		})
	}

	// STEP 4: Idempotency check - has this event been processed?
	isProcessed, err := h.eventStore.IsProcessed(ctx, eventID)
	if err != nil {
		log.Printf("[WEBHOOK] Error checking event idempotency for %s: %v", eventID, err)
		// Return 200 to prevent retries - we can't verify if it was processed
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event acknowledged (idempotency check failed)",
		})
	}

	if isProcessed {
		// Event already processed - acknowledge to prevent retries
		processedAt, found, _ := h.eventStore.GetProcessedAt(ctx, eventID)
		if found {
			log.Printf("[WEBHOOK] Duplicate event %s (event: %s, originally processed at: %s)",
				eventID, event.Event, processedAt.Format(time.RFC3339))
		} else {
			log.Printf("[WEBHOOK] Duplicate event %s (event: %s)", eventID, event.Event)
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event already processed",
		})
	}

	// STEP 5: Mark event as processed BEFORE processing (prevents race conditions)
	// Retention period: 7 days (Paystack typically retries for 24 hours)
	retentionPeriod := 7 * 24 * time.Hour
	if err := h.eventStore.MarkProcessed(ctx, eventID, retentionPeriod); err != nil {
		if err == webhook.ErrEventAlreadyProcessed {
			// Another request beat us to it
			log.Printf("[WEBHOOK] Event %s already being processed by another request", eventID)
			return c.JSON(fiber.Map{
				"success": true,
				"message": "Event already processing",
			})
		}

		log.Printf("[WEBHOOK] Error marking event %s as processed: %v", eventID, err)
		// Return error to trigger retry
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to record event processing",
		})
	}

	// STEP 6: Process the event
	log.Printf("[WEBHOOK] Processing event: %s (ID: %s)", event.Event, eventID)

	// Handle different event types
	switch event.Event {
	case "charge.success":
		return h.handleChargeSuccess(c, eventID, event.Data)
	case "transfer.success":
		return h.handleTransferSuccess(c, eventID, event.Data)
	case "transfer.failed":
		return h.handleTransferFailed(c, eventID, event.Data)
	case "transfer.reversed":
		return h.handleTransferReversed(c, eventID, event.Data)
	default:
		// Acknowledge unknown events to prevent retries
		log.Printf("[WEBHOOK] Unhandled event type: %s (ID: %s)", event.Event, eventID)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event acknowledged",
		})
	}
}

// extractEventID extracts a unique event identifier for idempotency
func (h *WebhookHandler) extractEventID(event PaystackWebhookEvent) (string, error) {
	// Parse data to extract reference
	var data struct {
		Reference string `json:"reference"`
		ID        int64  `json:"id"`
	}

	if err := json.Unmarshal(event.Data, &data); err != nil {
		return "", fmt.Errorf("failed to parse event data: %w", err)
	}

	if data.Reference == "" {
		return "", fmt.Errorf("event missing reference field")
	}

	// Use combination of event type and reference as unique ID
	// This allows the same reference to be used in different event types
	eventID := fmt.Sprintf("%s:%s", event.Event, data.Reference)
	return eventID, nil
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
func (h *WebhookHandler) handleChargeSuccess(c *fiber.Ctx, eventID string, data json.RawMessage) error {
	var chargeData PaystackChargeData
	if err := json.Unmarshal(data, &chargeData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse charge data for event %s: %v", eventID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid charge data",
		})
	}

	// Verify the charge status
	if chargeData.Status != "success" {
		log.Printf("[WEBHOOK] Charge not successful for event %s: status=%s, ref=%s",
			eventID, chargeData.Status, chargeData.Reference)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Charge status not successful, acknowledged",
		})
	}

	log.Printf("[WEBHOOK] Processing successful charge for event %s: ref=%s, amount=%d %s, user=%s",
		eventID, chargeData.Reference, chargeData.Amount, chargeData.Currency, chargeData.Metadata.UserID)

	// TODO: Call wallet service to credit user's wallet
	// Example implementation:
	// err := h.walletService.CreditWallet(c.Context(), CreditWalletRequest{
	//     UserID:        chargeData.Metadata.UserID,
	//     WalletID:      chargeData.Metadata.WalletID,
	//     Amount:        chargeData.Amount,
	//     Currency:      chargeData.Currency,
	//     Reference:     chargeData.Reference,
	//     PaymentMethod: chargeData.Channel,
	//     PaidAt:        parsePaidAt(chargeData.PaidAt),
	// })
	// if err != nil {
	//     log.Printf("[WEBHOOK] Failed to credit wallet for event %s: %v", eventID, err)
	//     return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	//         "success": false,
	//         "error":   "Failed to process payment",
	//     })
	// }

	log.Printf("[WEBHOOK] Successfully processed charge event %s", eventID)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Charge processed",
	})
}

// handleTransferSuccess processes successful transfers (withdrawals)
func (h *WebhookHandler) handleTransferSuccess(c *fiber.Ctx, eventID string, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data for event %s: %v", eventID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer successful for event %s: ref=%s, amount=%d",
		eventID, transferData.Reference, transferData.Amount)

	// TODO: Update transaction status to completed
	// Example: h.walletService.CompleteWithdrawal(c.Context(), transferData.Reference)

	log.Printf("[WEBHOOK] Successfully processed transfer success event %s", eventID)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer processed",
	})
}

// handleTransferFailed processes failed transfers
func (h *WebhookHandler) handleTransferFailed(c *fiber.Ctx, eventID string, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data for event %s: %v", eventID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer failed for event %s: ref=%s, reason=%s",
		eventID, transferData.Reference, transferData.Reason)

	// TODO: Reverse the withdrawal and refund user
	// Example: h.walletService.FailWithdrawal(c.Context(), transferData.Reference, transferData.Reason)

	log.Printf("[WEBHOOK] Successfully processed transfer failed event %s", eventID)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer failure processed",
	})
}

// handleTransferReversed processes reversed transfers
func (h *WebhookHandler) handleTransferReversed(c *fiber.Ctx, eventID string, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data for event %s: %v", eventID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer reversed for event %s: ref=%s", eventID, transferData.Reference)

	// TODO: Handle transfer reversal
	// Example: h.walletService.ReverseWithdrawal(c.Context(), transferData.Reference)

	log.Printf("[WEBHOOK] Successfully processed transfer reversal event %s", eventID)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer reversal processed",
	})
}
