package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"hustlex/internal/domain/shared/valueobject"
	walletservice "hustlex/internal/domain/wallet/service"
)

// Client configuration
type Config struct {
	SecretKey   string
	PublicKey   string
	BaseURL     string
	Timeout     time.Duration
	MaxRetries  int
}

// DefaultConfig returns default Paystack configuration
func DefaultConfig(secretKey, publicKey string) Config {
	return Config{
		SecretKey:  secretKey,
		PublicKey:  publicKey,
		BaseURL:    "https://api.paystack.co",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	}
}

// Client implements the PaymentProvider interface for Paystack
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new Paystack client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Paystack API Response structures
type paystackResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type initializeResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

type verifyResponse struct {
	ID              int64  `json:"id"`
	Status          string `json:"status"`
	Reference       string `json:"reference"`
	Amount          int64  `json:"amount"`
	Currency        string `json:"currency"`
	Channel         string `json:"channel"`
	GatewayResponse string `json:"gateway_response"`
	PaidAt          string `json:"paid_at"`
	Customer        struct {
		Email string `json:"email"`
	} `json:"customer"`
	Authorization struct {
		AuthorizationCode string `json:"authorization_code"`
		Bank              string `json:"bank"`
		Last4             string `json:"last4"`
		Channel           string `json:"channel"`
	} `json:"authorization"`
}

type transferRecipientResponse struct {
	RecipientCode string `json:"recipient_code"`
	Type          string `json:"type"`
	Name          string `json:"name"`
}

type transferResponse struct {
	Reference   string `json:"reference"`
	Integration int64  `json:"integration"`
	Domain      string `json:"domain"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	TransferCode string `json:"transfer_code"`
}

type bankListResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Type string `json:"type"`
}

// InitializeTransaction starts a payment transaction
func (c *Client) InitializeTransaction(ctx context.Context, req walletservice.PaymentRequest) (*walletservice.PaymentInitResponse, error) {
	payload := map[string]interface{}{
		"email":     req.Email,
		"amount":    req.Amount.Amount(), // Amount in kobo
		"reference": req.Reference,
		"currency":  string(req.Amount.Currency()),
		"callback_url": req.CallbackURL,
		"metadata": map[string]interface{}{
			"user_id":   req.UserID,
			"wallet_id": req.WalletID,
		},
	}

	if len(req.Channels) > 0 {
		payload["channels"] = req.Channels
	}

	body, err := c.doRequest(ctx, "POST", "/transaction/initialize", payload)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var initResp initializeResponse
	if err := json.Unmarshal(dataBytes, &initResp); err != nil {
		return nil, err
	}

	return &walletservice.PaymentInitResponse{
		AuthorizationURL: initResp.AuthorizationURL,
		AccessCode:       initResp.AccessCode,
		Reference:        initResp.Reference,
	}, nil
}

// VerifyTransaction verifies a payment transaction
func (c *Client) VerifyTransaction(ctx context.Context, reference string) (*walletservice.PaymentVerifyResponse, error) {
	body, err := c.doRequest(ctx, "GET", "/transaction/verify/"+reference, nil)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var verifyResp verifyResponse
	if err := json.Unmarshal(dataBytes, &verifyResp); err != nil {
		return nil, err
	}

	amount, _ := valueobject.NewMoney(verifyResp.Amount, valueobject.Currency(verifyResp.Currency))

	return &walletservice.PaymentVerifyResponse{
		Reference:       verifyResp.Reference,
		Status:          verifyResp.Status,
		Amount:          amount,
		Channel:         verifyResp.Channel,
		GatewayResponse: verifyResp.GatewayResponse,
		PaidAt:          verifyResp.PaidAt,
		CustomerEmail:   verifyResp.Customer.Email,
		AuthCode:        verifyResp.Authorization.AuthorizationCode,
		Bank:            verifyResp.Authorization.Bank,
		Last4:           verifyResp.Authorization.Last4,
	}, nil
}

// ChargeAuthorization charges a saved card
func (c *Client) ChargeAuthorization(ctx context.Context, req walletservice.ChargeRequest) (*walletservice.PaymentVerifyResponse, error) {
	payload := map[string]interface{}{
		"email":              req.Email,
		"amount":             req.Amount.Amount(),
		"authorization_code": req.AuthorizationCode,
		"reference":          req.Reference,
		"currency":           string(req.Amount.Currency()),
	}

	body, err := c.doRequest(ctx, "POST", "/transaction/charge_authorization", payload)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var verifyResp verifyResponse
	if err := json.Unmarshal(dataBytes, &verifyResp); err != nil {
		return nil, err
	}

	amount, _ := valueobject.NewMoney(verifyResp.Amount, valueobject.Currency(verifyResp.Currency))

	return &walletservice.PaymentVerifyResponse{
		Reference:       verifyResp.Reference,
		Status:          verifyResp.Status,
		Amount:          amount,
		Channel:         verifyResp.Channel,
		GatewayResponse: verifyResp.GatewayResponse,
	}, nil
}

// InitiateTransfer initiates a bank transfer
func (c *Client) InitiateTransfer(ctx context.Context, req walletservice.TransferRequest) (*walletservice.TransferResponse, error) {
	// First create transfer recipient
	recipientPayload := map[string]interface{}{
		"type":           "nuban",
		"name":           req.AccountName,
		"account_number": req.AccountNumber,
		"bank_code":      req.BankCode,
		"currency":       string(req.Amount.Currency()),
	}

	body, err := c.doRequest(ctx, "POST", "/transferrecipient", recipientPayload)
	if err != nil {
		return nil, err
	}

	var recipientResp paystackResponse
	if err := json.Unmarshal(body, &recipientResp); err != nil {
		return nil, err
	}

	if !recipientResp.Status {
		return nil, errors.New(recipientResp.Message)
	}

	dataBytes, _ := json.Marshal(recipientResp.Data)
	var recipient transferRecipientResponse
	if err := json.Unmarshal(dataBytes, &recipient); err != nil {
		return nil, err
	}

	// Now initiate transfer
	transferPayload := map[string]interface{}{
		"source":    "balance",
		"amount":    req.Amount.Amount(),
		"recipient": recipient.RecipientCode,
		"reason":    req.Reason,
		"reference": req.Reference,
	}

	body, err = c.doRequest(ctx, "POST", "/transfer", transferPayload)
	if err != nil {
		return nil, err
	}

	var transferResp paystackResponse
	if err := json.Unmarshal(body, &transferResp); err != nil {
		return nil, err
	}

	if !transferResp.Status {
		return nil, errors.New(transferResp.Message)
	}

	dataBytes, _ = json.Marshal(transferResp.Data)
	var transfer transferResponse
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return nil, err
	}

	return &walletservice.TransferResponse{
		Reference:    transfer.Reference,
		TransferCode: transfer.TransferCode,
		Status:       transfer.Status,
		Amount:       req.Amount,
	}, nil
}

// VerifyTransfer verifies a transfer status
func (c *Client) VerifyTransfer(ctx context.Context, reference string) (*walletservice.TransferResponse, error) {
	body, err := c.doRequest(ctx, "GET", "/transfer/verify/"+reference, nil)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var transfer transferResponse
	if err := json.Unmarshal(dataBytes, &transfer); err != nil {
		return nil, err
	}

	amount, _ := valueobject.NewMoney(transfer.Amount, valueobject.Currency(transfer.Currency))

	return &walletservice.TransferResponse{
		Reference:    transfer.Reference,
		TransferCode: transfer.TransferCode,
		Status:       transfer.Status,
		Amount:       amount,
	}, nil
}

// ListBanks returns list of supported banks
func (c *Client) ListBanks(ctx context.Context, country string) ([]walletservice.Bank, error) {
	url := "/bank?country=" + country
	body, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var bankList []bankListResponse
	if err := json.Unmarshal(dataBytes, &bankList); err != nil {
		return nil, err
	}

	banks := make([]walletservice.Bank, len(bankList))
	for i, b := range bankList {
		banks[i] = walletservice.Bank{
			Name: b.Name,
			Code: b.Code,
			Type: b.Type,
		}
	}

	return banks, nil
}

// ResolveAccountNumber validates a bank account
func (c *Client) ResolveAccountNumber(ctx context.Context, accountNumber, bankCode string) (*walletservice.AccountInfo, error) {
	url := fmt.Sprintf("/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankCode)
	body, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var resp paystackResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, errors.New(resp.Message)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var accountInfo struct {
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
		BankID        int64  `json:"bank_id"`
	}
	if err := json.Unmarshal(dataBytes, &accountInfo); err != nil {
		return nil, err
	}

	return &walletservice.AccountInfo{
		AccountNumber: accountInfo.AccountNumber,
		AccountName:   accountInfo.AccountName,
		BankCode:      bankCode,
	}, nil
}

// doRequest performs an HTTP request to Paystack API
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

	req.Header.Set("Authorization", "Bearer "+c.config.SecretKey)
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
			lastErr = fmt.Errorf("paystack server error: %d", resp.StatusCode)
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		return respBody, nil
	}

	return nil, lastErr
}

// Ensure Client implements PaymentProvider interface
var _ walletservice.PaymentProvider = (*Client)(nil)
