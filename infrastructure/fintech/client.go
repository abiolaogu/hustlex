// Package fintech provides a client wrapper for Global-FinTech services
package fintech

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config holds the fintech client configuration
type Config struct {
	BaseURL    string
	APIKey     string
	APISecret  string
	Timeout    time.Duration
	RetryCount int
}

// Client is the Global-FinTech API client
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new fintech client
func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = 3
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Currency represents supported currencies
type Currency string

const (
	CurrencyNGN Currency = "NGN"
	CurrencyUSD Currency = "USD"
	CurrencyGBP Currency = "GBP"
	CurrencyEUR Currency = "EUR"
	CurrencyCAD Currency = "CAD"
	CurrencyGHS Currency = "GHS"
	CurrencyKES Currency = "KES"
)

// FXQuote represents a foreign exchange quote
type FXQuote struct {
	SourceCurrency Currency  `json:"source_currency"`
	TargetCurrency Currency  `json:"target_currency"`
	Rate           float64   `json:"rate"`
	InverseRate    float64   `json:"inverse_rate"`
	Spread         float64   `json:"spread"`
	ValidUntil     time.Time `json:"valid_until"`
	QuoteID        string    `json:"quote_id"`
}

// RemittanceRequest represents a remittance transfer request
type RemittanceRequest struct {
	SourceCurrency   Currency `json:"source_currency"`
	TargetCurrency   Currency `json:"target_currency"`
	SourceAmount     float64  `json:"source_amount,omitempty"`
	TargetAmount     float64  `json:"target_amount,omitempty"`
	BeneficiaryID    string   `json:"beneficiary_id"`
	SenderID         string   `json:"sender_id"`
	Purpose          string   `json:"purpose"`
	QuoteID          string   `json:"quote_id,omitempty"`
	Reference        string   `json:"reference"`
	CallbackURL      string   `json:"callback_url,omitempty"`
}

// RemittanceResponse represents a remittance transfer response
type RemittanceResponse struct {
	TransactionID   string    `json:"transaction_id"`
	Status          string    `json:"status"`
	SourceAmount    float64   `json:"source_amount"`
	TargetAmount    float64   `json:"target_amount"`
	Fee             float64   `json:"fee"`
	Rate            float64   `json:"rate"`
	Reference       string    `json:"reference"`
	EstimatedArrival time.Time `json:"estimated_arrival"`
	CreatedAt       time.Time `json:"created_at"`
}

// Beneficiary represents a remittance beneficiary
type Beneficiary struct {
	ID            string            `json:"id"`
	UserID        string            `json:"user_id"`
	Name          string            `json:"name"`
	Email         string            `json:"email,omitempty"`
	Phone         string            `json:"phone"`
	Country       string            `json:"country"`
	Currency      Currency          `json:"currency"`
	BankCode      string            `json:"bank_code,omitempty"`
	AccountNumber string            `json:"account_number,omitempty"`
	MobileWallet  string            `json:"mobile_wallet,omitempty"`
	Relationship  string            `json:"relationship"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Verified      bool              `json:"verified"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// GetFXQuote retrieves a foreign exchange quote
func (c *Client) GetFXQuote(ctx context.Context, source, target Currency, amount float64) (*FXQuote, error) {
	url := fmt.Sprintf("%s/api/v1/fx/quote?source=%s&target=%s&amount=%.2f",
		c.config.BaseURL, source, target, amount)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var quote FXQuote
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &quote, nil
}

// CreateRemittance initiates a remittance transfer
func (c *Client) CreateRemittance(ctx context.Context, req RemittanceRequest) (*RemittanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/remittances", c.config.BaseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url,
		json.RawMessage(body).Reader())
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result RemittanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// GetRemittanceStatus retrieves the status of a remittance
func (c *Client) GetRemittanceStatus(ctx context.Context, transactionID string) (*RemittanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/remittances/%s", c.config.BaseURL, transactionID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result RemittanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// CreateBeneficiary creates a new beneficiary
func (c *Client) CreateBeneficiary(ctx context.Context, b Beneficiary) (*Beneficiary, error) {
	url := fmt.Sprintf("%s/api/v1/beneficiaries", c.config.BaseURL)

	body, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url,
		json.RawMessage(body).Reader())
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result Beneficiary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// ListBeneficiaries retrieves all beneficiaries for a user
func (c *Client) ListBeneficiaries(ctx context.Context, userID string) ([]Beneficiary, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/beneficiaries", c.config.BaseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result []Beneficiary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

// VerifyBeneficiary triggers verification for a beneficiary
func (c *Client) VerifyBeneficiary(ctx context.Context, beneficiaryID string) error {
	url := fmt.Sprintf("%s/api/v1/beneficiaries/%s/verify", c.config.BaseURL, beneficiaryID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// setHeaders sets common headers for API requests
func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("X-API-Key", c.config.APIKey)
	req.Header.Set("X-API-Secret", c.config.APISecret)
	req.Header.Set("Accept", "application/json")
}

// Corridor represents a supported remittance corridor
type Corridor struct {
	Source      Currency `json:"source"`
	Target      Currency `json:"target"`
	MinAmount   float64  `json:"min_amount"`
	MaxAmount   float64  `json:"max_amount"`
	Fee         float64  `json:"fee"`
	FeePercent  float64  `json:"fee_percent"`
	SpreadBps   int      `json:"spread_bps"` // Basis points
	DeliveryETA string   `json:"delivery_eta"`
}

// GetCorridors retrieves available remittance corridors
func (c *Client) GetCorridors(ctx context.Context) ([]Corridor, error) {
	url := fmt.Sprintf("%s/api/v1/corridors", c.config.BaseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result []Corridor
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}
