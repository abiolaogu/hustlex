package termii

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/notification/service"
	"hustlex/internal/domain/shared/valueobject"
)

// Config holds Termii configuration
type Config struct {
	APIKey     string
	SenderID   string
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
}

// DefaultConfig returns default Termii configuration
func DefaultConfig(apiKey, senderID string) Config {
	return Config{
		APIKey:     apiKey,
		SenderID:   senderID,
		BaseURL:    "https://api.ng.termii.com/api",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// Client implements the SMSProvider interface for Termii
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new Termii client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Termii API Response structures
type termiiResponse struct {
	Code      string `json:"code"`
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Balance   float64 `json:"balance"`
	User      string  `json:"user"`
}

type otpResponse struct {
	PinID     string `json:"pinId"`
	To        string `json:"to"`
	SmsStatus string `json:"smsStatus"`
}

type verifyOTPResponse struct {
	PinID    string `json:"pinId"`
	Verified bool   `json:"verified"`
	Msisdn   string `json:"msisdn"`
}

type balanceResponse struct {
	User     string  `json:"user"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// Send sends a notification via SMS
func (c *Client) Send(ctx context.Context, notification *aggregate.Notification) (string, error) {
	if notification.Channel() != aggregate.ChannelSMS {
		return "", errors.New("termii only supports SMS channel")
	}

	payload := map[string]interface{}{
		"api_key": c.config.APIKey,
		"to":      notification.Body(), // Assuming body contains phone number for SMS
		"from":    c.config.SenderID,
		"sms":     notification.Title() + "\n" + notification.Body(),
		"type":    "plain",
		"channel": "generic",
	}

	body, err := c.doRequest(ctx, "POST", "/sms/send", payload)
	if err != nil {
		return "", err
	}

	var resp termiiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if resp.Code != "ok" {
		return "", errors.New(resp.Message)
	}

	return resp.MessageID, nil
}

// GetStatus checks delivery status of a message
func (c *Client) GetStatus(ctx context.Context, messageID string) (aggregate.DeliveryStatus, error) {
	// Termii doesn't have a direct status check endpoint for regular SMS
	// This would typically use webhooks for delivery reports
	return aggregate.StatusSent, nil
}

// SupportsChannel returns true if this provider supports the channel
func (c *Client) SupportsChannel(channel aggregate.Channel) bool {
	return channel == aggregate.ChannelSMS
}

// SendOTP sends an OTP via SMS
func (c *Client) SendOTP(ctx context.Context, phone valueobject.PhoneNumber, otp string, expiryMinutes int) (string, error) {
	if expiryMinutes <= 0 {
		expiryMinutes = 5
	}

	// If OTP is provided, use send_token endpoint
	// Otherwise, use Termii's built-in OTP generation
	if otp != "" {
		payload := map[string]interface{}{
			"api_key":          c.config.APIKey,
			"message_type":     "NUMERIC",
			"to":               phone.String(),
			"from":             c.config.SenderID,
			"channel":          "generic",
			"pin_attempts":     3,
			"pin_time_to_live": expiryMinutes,
			"pin_length":       len(otp),
			"pin_placeholder":  "< 1234 >",
			"message_text":     fmt.Sprintf("Your HustleX verification code is %s. Valid for %d minutes.", otp, expiryMinutes),
			"pin_type":         "NUMERIC",
		}

		body, err := c.doRequest(ctx, "POST", "/sms/otp/send", payload)
		if err != nil {
			return "", err
		}

		var resp otpResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return "", err
		}

		if resp.SmsStatus != "Message Sent" {
			return "", fmt.Errorf("OTP sending failed: %s", resp.SmsStatus)
		}

		return resp.PinID, nil
	}

	// Use Termii's auto-generate OTP
	payload := map[string]interface{}{
		"api_key":          c.config.APIKey,
		"message_type":     "NUMERIC",
		"to":               phone.String(),
		"from":             c.config.SenderID,
		"channel":          "generic",
		"pin_attempts":     3,
		"pin_time_to_live": expiryMinutes,
		"pin_length":       6,
		"pin_placeholder":  "< 1234 >",
		"message_text":     "Your HustleX verification code is < 1234 >. Valid for " + fmt.Sprintf("%d", expiryMinutes) + " minutes.",
		"pin_type":         "NUMERIC",
	}

	body, err := c.doRequest(ctx, "POST", "/sms/otp/send", payload)
	if err != nil {
		return "", err
	}

	var resp otpResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.PinID, nil
}

// VerifyOTP verifies an OTP (if using Termii's built-in verification)
func (c *Client) VerifyOTP(ctx context.Context, pinID, otp string) (bool, error) {
	payload := map[string]interface{}{
		"api_key": c.config.APIKey,
		"pin_id":  pinID,
		"pin":     otp,
	}

	body, err := c.doRequest(ctx, "POST", "/sms/otp/verify", payload)
	if err != nil {
		return false, err
	}

	var resp verifyOTPResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, err
	}

	return resp.Verified, nil
}

// SendBulk sends SMS to multiple recipients
func (c *Client) SendBulk(ctx context.Context, phones []valueobject.PhoneNumber, message string) ([]string, error) {
	phoneStrings := make([]string, len(phones))
	for i, p := range phones {
		phoneStrings[i] = p.String()
	}

	payload := map[string]interface{}{
		"api_key": c.config.APIKey,
		"to":      phoneStrings,
		"from":    c.config.SenderID,
		"sms":     message,
		"type":    "plain",
		"channel": "generic",
	}

	body, err := c.doRequest(ctx, "POST", "/sms/send/bulk", payload)
	if err != nil {
		return nil, err
	}

	var resp termiiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Code != "ok" {
		return nil, errors.New(resp.Message)
	}

	// For bulk SMS, Termii returns a single message ID
	return []string{resp.MessageID}, nil
}

// GetBalance returns the SMS balance
func (c *Client) GetBalance(ctx context.Context) (float64, error) {
	url := fmt.Sprintf("/get-balance?api_key=%s", c.config.APIKey)
	body, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	var resp balanceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, err
	}

	return resp.Balance, nil
}

// SendWithTemplate sends SMS using a predefined template
func (c *Client) SendWithTemplate(ctx context.Context, phone valueobject.PhoneNumber, templateID string, data map[string]interface{}) (string, error) {
	payload := map[string]interface{}{
		"api_key":     c.config.APIKey,
		"phone_number": phone.String(),
		"device_id":   templateID,
		"template_id": templateID,
		"data":        data,
	}

	body, err := c.doRequest(ctx, "POST", "/sms/template", payload)
	if err != nil {
		return "", err
	}

	var resp termiiResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	if resp.Code != "ok" {
		return "", errors.New(resp.Message)
	}

	return resp.MessageID, nil
}

// doRequest performs an HTTP request to Termii API
func (c *Client) doRequest(ctx context.Context, method, path string, payload interface{}) ([]byte, error) {
	url := c.config.BaseURL + path

	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 500 {
			lastErr = fmt.Errorf("termii server error: %d", resp.StatusCode)
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		return respBody, nil
	}

	return nil, lastErr
}

// Ensure Client implements SMSProvider interface
var _ service.SMSProvider = (*Client)(nil)
