package handler

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"hustlex/internal/application/wallet/command"
	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
	"hustlex/internal/domain/wallet/repository"
	"hustlex/internal/domain/wallet/service"
)

const (
	// TransferFee is the flat fee for P2P transfers (₦10 = 1000 kobo)
	TransferFee int64 = 1000
	// MinTransfer is the minimum transfer amount (₦100 = 10000 kobo)
	MinTransfer int64 = 10000
	// MaxTransfer is the maximum transfer amount (₦500,000 = 50000000 kobo)
	MaxTransfer int64 = 50000000
)

// UserLookup defines the interface for looking up users
type UserLookup interface {
	FindByPhone(ctx context.Context, phone string) (*UserInfo, error)
}

type UserInfo struct {
	ID       string
	Phone    string
	FullName string
	HasWallet bool
}

// TransferHandler handles P2P transfer commands
type TransferHandler struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	transferService *service.TransferService
	userLookup      UserLookup
}

// NewTransferHandler creates a new transfer handler
func NewTransferHandler(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	transferService *service.TransferService,
	userLookup UserLookup,
) *TransferHandler {
	return &TransferHandler{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		transferService: transferService,
		userLookup:      userLookup,
	}
}

// Handle processes a P2P transfer request
func (h *TransferHandler) Handle(ctx context.Context, cmd command.Transfer) (*command.TransferResult, error) {
	// Validate amount
	amount, err := cmd.GetMoney()
	if err != nil {
		return nil, err
	}

	if amount.Amount() < MinTransfer {
		return nil, fmt.Errorf("minimum transfer is ₦%.2f", float64(MinTransfer)/100)
	}

	if amount.Amount() > MaxTransfer {
		return nil, fmt.Errorf("maximum transfer is ₦%.2f", float64(MaxTransfer)/100)
	}

	// Load sender wallet
	senderUserID, err := valueobject.NewUserID(cmd.FromUserID)
	if err != nil {
		return nil, err
	}

	senderWallet, err := h.walletRepo.FindByUserID(ctx, senderUserID)
	if err != nil {
		return nil, err
	}

	// Verify sender wallet is active
	if !senderWallet.IsActive() {
		return nil, aggregate.ErrWalletLocked
	}

	// Verify PIN
	if !senderWallet.HasPIN() {
		return nil, aggregate.ErrPINNotSet
	}

	if err := bcrypt.CompareHashAndPassword([]byte(senderWallet.PINHash()), []byte(cmd.PIN)); err != nil {
		senderWallet.RecordFailedPINAttempt(MaxPINAttempts)
		_ = h.walletRepo.Save(ctx, senderWallet)
		return nil, aggregate.ErrInvalidPIN
	}

	senderWallet.ResetPINAttempts()

	// Look up recipient by phone
	recipient, err := h.userLookup.FindByPhone(ctx, cmd.ToUserPhone)
	if err != nil {
		return nil, fmt.Errorf("recipient not found")
	}

	// Can't transfer to self
	if recipient.ID == cmd.FromUserID {
		return nil, service.ErrSameWallet
	}

	recipientUserID, err := valueobject.NewUserID(recipient.ID)
	if err != nil {
		return nil, err
	}

	// Generate reference
	reference := cmd.Reference
	if reference == "" {
		reference = generateReference("TRF")
	}

	// Calculate fee (for future use)
	_ = valueobject.MustNewMoney(TransferFee, amount.Currency())

	// Execute transfer via domain service
	transferResult, err := h.transferService.Transfer(ctx, service.TransferRequest{
		FromUserID:  senderUserID,
		ToUserID:    recipientUserID,
		Amount:      amount,
		Description: cmd.Description,
		Reference:   reference,
	})
	if err != nil {
		return nil, err
	}

	// Create transaction records
	// Sender transaction
	senderTx := &repository.Transaction{
		WalletID:       senderWallet.ID().String(),
		Type:           repository.TransactionTypeTransferOut,
		Amount:         amount.Amount(),
		Fee:            TransferFee,
		Currency:       string(amount.Currency()),
		BalanceAfter:   transferResult.SenderNewBalance.Amount(),
		Status:         repository.TransactionStatusCompleted,
		Reference:      reference,
		Description:    fmt.Sprintf("Transfer to %s", recipient.FullName),
		CounterpartyID: &recipient.ID,
	}
	_ = h.transactionRepo.Save(ctx, senderTx)

	// Recipient transaction
	recipientTx := &repository.Transaction{
		WalletID:       recipient.ID, // This should be wallet ID
		Type:           repository.TransactionTypeTransferIn,
		Amount:         amount.Amount(),
		Fee:            0,
		Currency:       string(amount.Currency()),
		BalanceAfter:   transferResult.RecipientNewBalance.Amount(),
		Status:         repository.TransactionStatusCompleted,
		Reference:      reference,
		Description:    fmt.Sprintf("Transfer from %s", cmd.FromUserID), // Should be sender name
		CounterpartyID: &cmd.FromUserID,
	}
	_ = h.transactionRepo.Save(ctx, recipientTx)

	return &command.TransferResult{
		TransactionID: senderTx.ID,
		Reference:     reference,
		RecipientName: recipient.FullName,
		Amount:        amount.Amount(),
		Fee:           TransferFee,
		NewBalance:    transferResult.SenderNewBalance.Amount(),
		ProcessedAt:   time.Now().UTC(),
	}, nil
}
