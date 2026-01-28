package handler

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hustlex/internal/application/wallet/command"
	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
	"hustlex/internal/domain/wallet/repository"
)

const (
	// WithdrawalFee is the flat fee for withdrawals (₦50 = 5000 kobo)
	WithdrawalFee int64 = 5000
	// MinWithdrawal is the minimum withdrawal amount (₦500 = 50000 kobo)
	MinWithdrawal int64 = 50000
	// MaxWithdrawal is the maximum withdrawal amount (₦1,000,000 = 100000000 kobo)
	MaxWithdrawal int64 = 100000000
	// MaxPINAttempts is the maximum number of PIN attempts before locking
	MaxPINAttempts int = 5
)

// TransferProvider defines the interface for bank transfers
// This is a PORT - infrastructure provides the ADAPTER (e.g., Paystack Transfer)
type TransferProvider interface {
	// InitiateTransfer initiates a bank transfer
	InitiateTransfer(ctx context.Context, req InitiateTransferRequest) (*InitiateTransferResponse, error)
	// VerifyAccountNumber verifies a bank account and returns the account name
	VerifyAccountNumber(ctx context.Context, bankCode, accountNumber string) (string, error)
	// GetBanks returns the list of supported banks
	GetBanks(ctx context.Context) ([]Bank, error)
}

type InitiateTransferRequest struct {
	Amount        int64
	Recipient     string // Recipient code or account details
	Reference     string
	Reason        string
	Currency      string
	BankCode      string
	AccountNumber string
	AccountName   string
}

type InitiateTransferResponse struct {
	Reference    string
	TransferCode string
	Status       string
	Message      string
}

type Bank struct {
	Code      string
	Name      string
	LongCode  string
	Country   string
	Currency  string
	IsActive  bool
}

// WithdrawHandler handles withdrawal commands
type WithdrawHandler struct {
	walletRepo       repository.WalletRepository
	transactionRepo  repository.TransactionRepository
	transferProvider TransferProvider
}

// NewWithdrawHandler creates a new withdrawal handler
func NewWithdrawHandler(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	transferProvider TransferProvider,
) *WithdrawHandler {
	return &WithdrawHandler{
		walletRepo:       walletRepo,
		transactionRepo:  transactionRepo,
		transferProvider: transferProvider,
	}
}

// Handle processes a withdrawal request
func (h *WithdrawHandler) Handle(ctx context.Context, cmd command.Withdraw) (*command.WithdrawResult, error) {
	// Validate amount
	amount, err := cmd.GetMoney()
	if err != nil {
		return nil, err
	}

	if amount.Amount() < MinWithdrawal {
		return nil, fmt.Errorf("minimum withdrawal is ₦%.2f", float64(MinWithdrawal)/100)
	}

	if amount.Amount() > MaxWithdrawal {
		return nil, fmt.Errorf("maximum withdrawal is ₦%.2f", float64(MaxWithdrawal)/100)
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

	// Verify wallet is active
	if !wallet.IsActive() {
		return nil, aggregate.ErrWalletLocked
	}

	// Verify PIN is set
	if !wallet.HasPIN() {
		return nil, aggregate.ErrPINNotSet
	}

	// Verify PIN
	if err := bcrypt.CompareHashAndPassword([]byte(wallet.PINHash()), []byte(cmd.PIN)); err != nil {
		wallet.RecordFailedPINAttempt(MaxPINAttempts)
		if saveErr := h.walletRepo.Save(ctx, wallet); saveErr != nil {
			// Log but don't fail
		}
		return nil, aggregate.ErrInvalidPIN
	}

	// Reset PIN attempts on success
	wallet.ResetPINAttempts()

	// Calculate fee
	fee := valueobject.MustNewMoney(WithdrawalFee, amount.Currency())

	// Generate reference if not provided
	reference := cmd.Reference
	if reference == "" {
		reference = generateReference("WTH")
	}

	// Debit the wallet (amount + fee)
	description := fmt.Sprintf("Withdrawal to %s - ****%s", cmd.BankCode, cmd.AccountNumber[len(cmd.AccountNumber)-4:])
	err = wallet.Debit(amount, "withdrawal", reference, description, fee)
	if err != nil {
		return nil, err
	}

	// Save wallet state
	if err := h.walletRepo.SaveWithEvents(ctx, wallet); err != nil {
		return nil, err
	}

	// Initiate bank transfer
	transferResp, err := h.transferProvider.InitiateTransfer(ctx, InitiateTransferRequest{
		Amount:        amount.Amount(),
		Reference:     reference,
		Reason:        "HustleX Withdrawal",
		Currency:      string(amount.Currency()),
		BankCode:      cmd.BankCode,
		AccountNumber: cmd.AccountNumber,
		AccountName:   cmd.AccountName,
	})
	if err != nil {
		// Transfer initiation failed - we need to refund
		// In production, this would be handled by a saga or compensation
		refundErr := wallet.Credit(amount.MustAdd(fee), "refund", reference+"-REFUND", "Withdrawal failed - refund")
		if refundErr == nil {
			_ = h.walletRepo.SaveWithEvents(ctx, wallet)
		}
		return nil, fmt.Errorf("failed to initiate transfer: %w", err)
	}

	// Create transaction record
	tx := &repository.Transaction{
		WalletID:      wallet.ID().String(),
		Type:          repository.TransactionTypeWithdrawal,
		Amount:        amount.Amount(),
		Fee:           WithdrawalFee,
		Currency:      string(amount.Currency()),
		BalanceAfter:  wallet.AvailableBalance().Amount(),
		Status:        repository.TransactionStatusPending,
		Reference:     reference,
		Description:   description,
		BankCode:      &cmd.BankCode,
		AccountNumber: &cmd.AccountNumber,
		AccountName:   &cmd.AccountName,
	}

	if err := h.transactionRepo.Save(ctx, tx); err != nil {
		// Log but don't fail - transfer is already initiated
	}

	return &command.WithdrawResult{
		TransactionID: tx.ID,
		Reference:     reference,
		Status:        transferResp.Status,
		NewBalance:    wallet.AvailableBalance().Amount(),
		Fee:           WithdrawalFee,
		ProcessedAt:   time.Now().UTC(),
	}, nil
}

// HandleVerifyAccount verifies a bank account
func (h *WithdrawHandler) HandleVerifyAccount(ctx context.Context, bankCode, accountNumber string) (string, error) {
	return h.transferProvider.VerifyAccountNumber(ctx, bankCode, accountNumber)
}

// HandleGetBanks returns the list of supported banks
func (h *WithdrawHandler) HandleGetBanks(ctx context.Context) ([]Bank, error) {
	return h.transferProvider.GetBanks(ctx)
}

func generateReference(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, n)
}
