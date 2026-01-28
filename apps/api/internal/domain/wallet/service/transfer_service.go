package service

import (
	"context"
	"errors"

	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/repository"
)

// Domain errors for transfer operations
var (
	ErrSameWallet          = errors.New("cannot transfer to the same wallet")
	ErrRecipientNotFound   = errors.New("recipient wallet not found")
	ErrTransferLimitExceeded = errors.New("transfer limit exceeded")
)

// TransferService handles P2P transfer operations
// This is a domain service because transfers span two aggregates
type TransferService struct {
	walletRepo repository.WalletRepository
}

// NewTransferService creates a new transfer service
func NewTransferService(walletRepo repository.WalletRepository) *TransferService {
	return &TransferService{walletRepo: walletRepo}
}

// TransferRequest represents a P2P transfer request
type TransferRequest struct {
	FromUserID  valueobject.UserID
	ToUserID    valueobject.UserID
	Amount      valueobject.Money
	Description string
	Reference   string
}

// TransferResult represents the result of a transfer
type TransferResult struct {
	SenderNewBalance    valueobject.Money
	RecipientNewBalance valueobject.Money
	Reference           string
	Fee                 valueobject.Money
}

// Transfer executes a P2P transfer between two wallets
// Note: In a production system, this would need to handle distributed transactions
// or use the saga pattern for proper consistency across aggregates
func (s *TransferService) Transfer(ctx context.Context, req TransferRequest) (*TransferResult, error) {
	// Validate not transferring to self
	if req.FromUserID.Equals(req.ToUserID) {
		return nil, ErrSameWallet
	}

	// Load sender wallet
	senderWallet, err := s.walletRepo.FindByUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}

	// Load recipient wallet
	recipientWallet, err := s.walletRepo.FindByUserID(ctx, req.ToUserID)
	if err != nil {
		return nil, ErrRecipientNotFound
	}

	// Calculate fee (e.g., ₦10 flat fee for transfers)
	fee := valueobject.MustNewMoney(1000, req.Amount.Currency()) // 1000 kobo = ₦10

	// Debit sender
	err = senderWallet.Debit(req.Amount, "transfer_out", req.Reference, req.Description, fee)
	if err != nil {
		return nil, err
	}

	// Credit recipient (no fee on receiving end)
	err = recipientWallet.Credit(req.Amount, "transfer_in", req.Reference, req.Description)
	if err != nil {
		// In a real system, we'd need to rollback or compensate
		return nil, err
	}

	// Save both wallets
	// In a production system, these would need to be in a distributed transaction
	// or use eventual consistency with compensation
	if err := s.walletRepo.SaveWithEvents(ctx, senderWallet); err != nil {
		return nil, err
	}

	if err := s.walletRepo.SaveWithEvents(ctx, recipientWallet); err != nil {
		// Would need compensation here
		return nil, err
	}

	return &TransferResult{
		SenderNewBalance:    senderWallet.AvailableBalance(),
		RecipientNewBalance: recipientWallet.AvailableBalance(),
		Reference:           req.Reference,
		Fee:                 fee,
	}, nil
}

// EscrowService handles escrow operations for gig payments
type EscrowService struct {
	walletRepo repository.WalletRepository
}

// NewEscrowService creates a new escrow service
func NewEscrowService(walletRepo repository.WalletRepository) *EscrowService {
	return &EscrowService{walletRepo: walletRepo}
}

// EscrowReleaseRequest represents a request to release escrowed funds
type EscrowReleaseRequest struct {
	PayerUserID     valueobject.UserID
	RecipientUserID valueobject.UserID
	Amount          valueobject.Money
	PlatformFee     valueobject.Money
	Reference       string
	Description     string
}

// ReleaseEscrowToRecipient releases escrowed funds to a recipient
func (s *EscrowService) ReleaseEscrowToRecipient(ctx context.Context, req EscrowReleaseRequest) error {
	// Load payer wallet
	payerWallet, err := s.walletRepo.FindByUserID(ctx, req.PayerUserID)
	if err != nil {
		return err
	}

	// Load recipient wallet
	recipientWallet, err := s.walletRepo.FindByUserID(ctx, req.RecipientUserID)
	if err != nil {
		return err
	}

	// Release from payer's escrow (not back to their wallet)
	err = payerWallet.ReleaseFromEscrow(req.Amount, req.Reference, false, req.RecipientUserID.String())
	if err != nil {
		return err
	}

	// Credit recipient with amount minus platform fee
	netAmount, err := req.Amount.Subtract(req.PlatformFee)
	if err != nil {
		return err
	}

	err = recipientWallet.Credit(netAmount, "gig_payment", req.Reference, req.Description)
	if err != nil {
		return err
	}

	// Save both wallets
	if err := s.walletRepo.SaveWithEvents(ctx, payerWallet); err != nil {
		return err
	}

	if err := s.walletRepo.SaveWithEvents(ctx, recipientWallet); err != nil {
		return err
	}

	return nil
}

// RefundEscrow returns escrowed funds to the payer
func (s *EscrowService) RefundEscrow(ctx context.Context, userID valueobject.UserID, amount valueobject.Money, reference, reason string) error {
	wallet, err := s.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Release back to the wallet owner
	err = wallet.ReleaseFromEscrow(amount, reference, true, "")
	if err != nil {
		return err
	}

	return s.walletRepo.SaveWithEvents(ctx, wallet)
}
