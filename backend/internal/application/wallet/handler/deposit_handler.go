package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hustlex/internal/application/wallet/command"
	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
	"hustlex/internal/domain/wallet/repository"
)

// PaymentGateway defines the interface for payment processing
// This is a PORT - infrastructure will provide the ADAPTER (e.g., Paystack)
type PaymentGateway interface {
	// InitiateDeposit initiates a deposit transaction
	InitiateDeposit(ctx context.Context, req InitiateDepositRequest) (*InitiateDepositResponse, error)
	// VerifyPayment verifies a payment by reference
	VerifyPayment(ctx context.Context, reference string) (*PaymentVerification, error)
}

type InitiateDepositRequest struct {
	UserID      string
	Email       string
	Amount      int64 // In kobo
	Currency    string
	Reference   string
	CallbackURL string
	Metadata    map[string]interface{}
}

type InitiateDepositResponse struct {
	Reference   string
	AccessCode  string
	PaymentURL  string
	ExpiresAt   time.Time
}

type PaymentVerification struct {
	Reference     string
	Status        string // success, failed, pending
	Amount        int64
	Currency      string
	Channel       string
	PaidAt        time.Time
	GatewayResponse string
}

// DepositHandler handles deposit commands
type DepositHandler struct {
	walletRepo     repository.WalletRepository
	paymentGateway PaymentGateway
}

// NewDepositHandler creates a new deposit handler
func NewDepositHandler(
	walletRepo repository.WalletRepository,
	paymentGateway PaymentGateway,
) *DepositHandler {
	return &DepositHandler{
		walletRepo:     walletRepo,
		paymentGateway: paymentGateway,
	}
}

// HandleInitiateDeposit initiates a deposit via payment gateway
func (h *DepositHandler) HandleInitiateDeposit(ctx context.Context, cmd command.Deposit) (*command.DepositResult, error) {
	// Validate amount
	amount, err := cmd.GetMoney()
	if err != nil {
		return nil, err
	}

	if !amount.IsPositive() {
		return nil, errors.New("amount must be positive")
	}

	// Minimum deposit check (₦100 = 10000 kobo)
	minDeposit := valueobject.MustNewMoney(10000, amount.Currency())
	if amount.LessThan(minDeposit) {
		return nil, fmt.Errorf("minimum deposit is %s", minDeposit.String())
	}

	// Verify wallet exists
	userID, err := valueobject.NewUserID(cmd.RequestedBy)
	if err != nil {
		return nil, err
	}

	wallet, err := h.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !wallet.IsActive() {
		return nil, aggregate.ErrWalletLocked
	}

	// Initiate payment via gateway
	gatewayResp, err := h.paymentGateway.InitiateDeposit(ctx, InitiateDepositRequest{
		UserID:    cmd.RequestedBy,
		Amount:    cmd.Amount,
		Currency:  cmd.Currency,
		Reference: cmd.Reference,
		Metadata: map[string]interface{}{
			"wallet_id": wallet.ID().String(),
			"channel":   cmd.Channel,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	return &command.DepositResult{
		Reference:   gatewayResp.Reference,
		PaymentURL:  gatewayResp.PaymentURL,
		AccessCode:  gatewayResp.AccessCode,
		ProcessedAt: time.Now().UTC(),
	}, nil
}

// HandleVerifyDeposit verifies and credits a deposit
func (h *DepositHandler) HandleVerifyDeposit(ctx context.Context, reference string) (*command.DepositResult, error) {
	// Verify with payment gateway
	verification, err := h.paymentGateway.VerifyPayment(ctx, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}

	if verification.Status != "success" {
		return nil, fmt.Errorf("payment not successful: %s", verification.Status)
	}

	// Get wallet from transaction metadata or reference
	// In a real implementation, you'd look this up from a pending transactions table
	// For now, we'll assume the reference contains the wallet ID or user lookup

	// This is a simplified flow - production would have proper idempotency
	return &command.DepositResult{
		Reference:   verification.Reference,
		ProcessedAt: verification.PaidAt,
	}, nil
}

// HandleCreditWallet directly credits a wallet (used after payment verification)
func (h *DepositHandler) HandleCreditWallet(ctx context.Context, cmd command.Deposit) (*command.DepositResult, error) {
	// Get the money amount
	amount, err := cmd.GetMoney()
	if err != nil {
		return nil, err
	}

	// Load wallet
	userID, err := valueobject.NewUserID(cmd.RequestedBy)
	if err != nil {
		return nil, err
	}

	wallet, err := h.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Credit the wallet
	err = wallet.Credit(amount, cmd.Source, cmd.Reference, cmd.Description)
	if err != nil {
		return nil, err
	}

	// Persist with events
	if err := h.walletRepo.SaveWithEvents(ctx, wallet); err != nil {
		return nil, err
	}

	return &command.DepositResult{
		TransactionID: cmd.Reference,
		Reference:     cmd.Reference,
		NewBalance:    wallet.AvailableBalance().Amount(),
		ProcessedAt:   time.Now().UTC(),
	}, nil
}
