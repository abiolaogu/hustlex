package handlers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"hustlex/internal/config"
	"hustlex/internal/services"
)

// WebhookHandler handles payment webhook callbacks
type WebhookHandler struct {
	walletService *services.WalletService
	config        *config.Config
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(walletService *services.WalletService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		walletService: walletService,
		config:        cfg,
	}
}

// PaystackWebhookEvent represents a Paystack webhook event
type PaystackWebhookEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

// PaystackChargeData represents charge success data
type PaystackChargeData struct {
	Reference     string `json:"reference"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	Channel       string `json:"channel"`
	GatewayResponse string `json:"gateway_response"`
	Metadata      struct {
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

// HandlePaystackWebhook handles Paystack payment webhooks
// @Summary Handle Paystack webhook
// @Description Process payment notifications from Paystack
// @Tags Webhooks
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /webhooks/paystack [post]
func (h *WebhookHandler) HandlePaystackWebhook(c *fiber.Ctx) error {
	// Read the raw body for signature verification
	body := c.Body()

	// Verify webhook signature - CRITICAL SECURITY CHECK
	signature := c.Get("X-Paystack-Signature")
	if signature == "" {
		log.Printf("[SECURITY] Webhook received without signature")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Missing webhook signature",
		})
	}

	// Verify the signature using HMAC-SHA512
	if !h.verifyPaystackSignature(body, signature) {
		log.Printf("[SECURITY] Webhook signature verification failed")
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

	// Log the event for audit trail (without sensitive data)
	log.Printf("[WEBHOOK] Received Paystack event: %s", event.Event)

	// Handle different event types
	switch event.Event {
	case "charge.success":
		return h.handleChargeSuccess(c, event.Data)
	case "transfer.success":
		return h.handleTransferSuccess(c, event.Data)
	case "transfer.failed":
		return h.handleTransferFailed(c, event.Data)
	case "transfer.reversed":
		return h.handleTransferReversed(c, event.Data)
	default:
		// Acknowledge unknown events to prevent retries
		log.Printf("[WEBHOOK] Unhandled event type: %s", event.Event)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Event acknowledged",
		})
	}
}

// verifyPaystackSignature verifies the webhook signature using HMAC-SHA512
func (h *WebhookHandler) verifyPaystackSignature(payload []byte, signature string) bool {
	secret := h.config.Payment.WebhookSecret
	if secret == "" {
		log.Printf("[SECURITY] Webhook secret not configured")
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
func (h *WebhookHandler) handleChargeSuccess(c *fiber.Ctx, data json.RawMessage) error {
	var chargeData PaystackChargeData
	if err := json.Unmarshal(data, &chargeData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse charge data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid charge data",
		})
	}

	// Verify the charge status
	if chargeData.Status != "success" {
		log.Printf("[WEBHOOK] Charge not successful: %s", chargeData.Reference)
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Charge status not successful, ignoring",
		})
	}

	log.Printf("[WEBHOOK] Processing successful charge: %s, amount: %d",
		chargeData.Reference, chargeData.Amount)

	// TODO: Call wallet service to credit user's wallet
	// Example: h.walletService.ProcessDeposit(c.Context(), chargeData.Reference)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Charge processed",
	})
}

// handleTransferSuccess processes successful transfers (withdrawals)
func (h *WebhookHandler) handleTransferSuccess(c *fiber.Ctx, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer successful: %s", transferData.Reference)

	// TODO: Update transaction status to completed
	// Example: h.walletService.CompleteWithdrawal(c.Context(), transferData.Reference)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer processed",
	})
}

// handleTransferFailed processes failed transfers
func (h *WebhookHandler) handleTransferFailed(c *fiber.Ctx, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer failed: %s, reason: %s",
		transferData.Reference, transferData.Reason)

	// TODO: Reverse the withdrawal and refund user
	// Example: h.walletService.FailWithdrawal(c.Context(), transferData.Reference, transferData.Reason)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer failure processed",
	})
}

// handleTransferReversed processes reversed transfers
func (h *WebhookHandler) handleTransferReversed(c *fiber.Ctx, data json.RawMessage) error {
	var transferData PaystackTransferData
	if err := json.Unmarshal(data, &transferData); err != nil {
		log.Printf("[WEBHOOK] Failed to parse transfer data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid transfer data",
		})
	}

	log.Printf("[WEBHOOK] Transfer reversed: %s", transferData.Reference)

	// TODO: Handle transfer reversal
	// Example: h.walletService.ReverseWithdrawal(c.Context(), transferData.Reference)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer reversal processed",
	})
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
		// Paystack IPs: 52.31.139.75, 52.49.173.169, 52.214.14.220
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
func PaystackWebhookIPs() []string {
	return []string{
		"52.31.139.75",
		"52.49.173.169",
		"52.214.14.220",
	}
}
