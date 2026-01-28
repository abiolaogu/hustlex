package event

import (
	"time"

	"hustlex/internal/domain/shared/event"
)

const AggregateTypeWallet = "Wallet"

// WalletCreated is raised when a new wallet is created
type WalletCreated struct {
	event.BaseEvent
	WalletID  string `json:"wallet_id"`
	UserID    string `json:"user_id"`
	Currency  string `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

func NewWalletCreated(walletID, userID, currency string) WalletCreated {
	return WalletCreated{
		BaseEvent: event.NewBaseEvent("WalletCreated", walletID, AggregateTypeWallet),
		WalletID:  walletID,
		UserID:    userID,
		Currency:  currency,
		CreatedAt: time.Now().UTC(),
	}
}

// WalletCredited is raised when funds are added to a wallet
type WalletCredited struct {
	event.BaseEvent
	WalletID     string    `json:"wallet_id"`
	UserID       string    `json:"user_id"`
	Amount       int64     `json:"amount"`
	Currency     string    `json:"currency"`
	Source       string    `json:"source"` // deposit, gig_payment, payout, refund, transfer_in
	Reference    string    `json:"reference"`
	Description  string    `json:"description"`
	NewBalance   int64     `json:"new_balance"`
	TransactedAt time.Time `json:"transacted_at"`
}

func NewWalletCredited(walletID, userID string, amount int64, currency, source, reference, description string, newBalance int64) WalletCredited {
	return WalletCredited{
		BaseEvent:    event.NewBaseEvent("WalletCredited", walletID, AggregateTypeWallet),
		WalletID:     walletID,
		UserID:       userID,
		Amount:       amount,
		Currency:     currency,
		Source:       source,
		Reference:    reference,
		Description:  description,
		NewBalance:   newBalance,
		TransactedAt: time.Now().UTC(),
	}
}

// WalletDebited is raised when funds are removed from a wallet
type WalletDebited struct {
	event.BaseEvent
	WalletID     string    `json:"wallet_id"`
	UserID       string    `json:"user_id"`
	Amount       int64     `json:"amount"`
	Currency     string    `json:"currency"`
	Destination  string    `json:"destination"` // withdrawal, transfer_out, escrow, contribution, loan_repayment
	Reference    string    `json:"reference"`
	Description  string    `json:"description"`
	Fee          int64     `json:"fee"`
	NewBalance   int64     `json:"new_balance"`
	TransactedAt time.Time `json:"transacted_at"`
}

func NewWalletDebited(walletID, userID string, amount int64, currency, destination, reference, description string, fee, newBalance int64) WalletDebited {
	return WalletDebited{
		BaseEvent:    event.NewBaseEvent("WalletDebited", walletID, AggregateTypeWallet),
		WalletID:     walletID,
		UserID:       userID,
		Amount:       amount,
		Currency:     currency,
		Destination:  destination,
		Reference:    reference,
		Description:  description,
		Fee:          fee,
		NewBalance:   newBalance,
		TransactedAt: time.Now().UTC(),
	}
}

// FundsHeldInEscrow is raised when funds are moved to escrow
type FundsHeldInEscrow struct {
	event.BaseEvent
	WalletID       string    `json:"wallet_id"`
	UserID         string    `json:"user_id"`
	Amount         int64     `json:"amount"`
	Reference      string    `json:"reference"` // contract_id or gig_id
	Reason         string    `json:"reason"`
	NewAvailable   int64     `json:"new_available_balance"`
	NewEscrow      int64     `json:"new_escrow_balance"`
	HeldAt         time.Time `json:"held_at"`
}

func NewFundsHeldInEscrow(walletID, userID string, amount int64, reference, reason string, newAvailable, newEscrow int64) FundsHeldInEscrow {
	return FundsHeldInEscrow{
		BaseEvent:    event.NewBaseEvent("FundsHeldInEscrow", walletID, AggregateTypeWallet),
		WalletID:     walletID,
		UserID:       userID,
		Amount:       amount,
		Reference:    reference,
		Reason:       reason,
		NewAvailable: newAvailable,
		NewEscrow:    newEscrow,
		HeldAt:       time.Now().UTC(),
	}
}

// FundsReleasedFromEscrow is raised when escrowed funds are released
type FundsReleasedFromEscrow struct {
	event.BaseEvent
	WalletID     string    `json:"wallet_id"`
	UserID       string    `json:"user_id"`
	Amount       int64     `json:"amount"`
	Reference    string    `json:"reference"`
	RecipientID  string    `json:"recipient_id,omitempty"` // If released to another user
	ToWallet     bool      `json:"to_wallet"`              // True if returned to same wallet
	NewEscrow    int64     `json:"new_escrow_balance"`
	ReleasedAt   time.Time `json:"released_at"`
}

func NewFundsReleasedFromEscrow(walletID, userID string, amount int64, reference, recipientID string, toWallet bool, newEscrow int64) FundsReleasedFromEscrow {
	return FundsReleasedFromEscrow{
		BaseEvent:   event.NewBaseEvent("FundsReleasedFromEscrow", walletID, AggregateTypeWallet),
		WalletID:    walletID,
		UserID:      userID,
		Amount:      amount,
		Reference:   reference,
		RecipientID: recipientID,
		ToWallet:    toWallet,
		NewEscrow:   newEscrow,
		ReleasedAt:  time.Now().UTC(),
	}
}

// WalletLocked is raised when a wallet is locked
type WalletLocked struct {
	event.BaseEvent
	WalletID string    `json:"wallet_id"`
	UserID   string    `json:"user_id"`
	Reason   string    `json:"reason"`
	LockedAt time.Time `json:"locked_at"`
}

func NewWalletLocked(walletID, userID, reason string) WalletLocked {
	return WalletLocked{
		BaseEvent: event.NewBaseEvent("WalletLocked", walletID, AggregateTypeWallet),
		WalletID:  walletID,
		UserID:    userID,
		Reason:    reason,
		LockedAt:  time.Now().UTC(),
	}
}

// WalletUnlocked is raised when a wallet is unlocked
type WalletUnlocked struct {
	event.BaseEvent
	WalletID   string    `json:"wallet_id"`
	UserID     string    `json:"user_id"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

func NewWalletUnlocked(walletID, userID string) WalletUnlocked {
	return WalletUnlocked{
		BaseEvent:  event.NewBaseEvent("WalletUnlocked", walletID, AggregateTypeWallet),
		WalletID:   walletID,
		UserID:     userID,
		UnlockedAt: time.Now().UTC(),
	}
}

// TransactionPINSet is raised when a PIN is set or changed
type TransactionPINSet struct {
	event.BaseEvent
	WalletID string    `json:"wallet_id"`
	UserID   string    `json:"user_id"`
	SetAt    time.Time `json:"set_at"`
}

func NewTransactionPINSet(walletID, userID string) TransactionPINSet {
	return TransactionPINSet{
		BaseEvent: event.NewBaseEvent("TransactionPINSet", walletID, AggregateTypeWallet),
		WalletID:  walletID,
		UserID:    userID,
		SetAt:     time.Now().UTC(),
	}
}

// WithdrawalInitiated is raised when a withdrawal is requested
type WithdrawalInitiated struct {
	event.BaseEvent
	WalletID      string    `json:"wallet_id"`
	UserID        string    `json:"user_id"`
	Amount        int64     `json:"amount"`
	Fee           int64     `json:"fee"`
	Reference     string    `json:"reference"`
	BankCode      string    `json:"bank_code"`
	AccountNumber string    `json:"account_number"`
	AccountName   string    `json:"account_name"`
	InitiatedAt   time.Time `json:"initiated_at"`
}

func NewWithdrawalInitiated(walletID, userID string, amount, fee int64, reference, bankCode, accountNumber, accountName string) WithdrawalInitiated {
	return WithdrawalInitiated{
		BaseEvent:     event.NewBaseEvent("WithdrawalInitiated", walletID, AggregateTypeWallet),
		WalletID:      walletID,
		UserID:        userID,
		Amount:        amount,
		Fee:           fee,
		Reference:     reference,
		BankCode:      bankCode,
		AccountNumber: accountNumber,
		AccountName:   accountName,
		InitiatedAt:   time.Now().UTC(),
	}
}

// WithdrawalCompleted is raised when a withdrawal succeeds
type WithdrawalCompleted struct {
	event.BaseEvent
	WalletID    string    `json:"wallet_id"`
	Reference   string    `json:"reference"`
	CompletedAt time.Time `json:"completed_at"`
}

func NewWithdrawalCompleted(walletID, reference string) WithdrawalCompleted {
	return WithdrawalCompleted{
		BaseEvent:   event.NewBaseEvent("WithdrawalCompleted", walletID, AggregateTypeWallet),
		WalletID:    walletID,
		Reference:   reference,
		CompletedAt: time.Now().UTC(),
	}
}

// WithdrawalFailed is raised when a withdrawal fails
type WithdrawalFailed struct {
	event.BaseEvent
	WalletID   string    `json:"wallet_id"`
	Reference  string    `json:"reference"`
	Reason     string    `json:"reason"`
	Amount     int64     `json:"amount"` // Amount to be refunded
	FailedAt   time.Time `json:"failed_at"`
}

func NewWithdrawalFailed(walletID, reference, reason string, amount int64) WithdrawalFailed {
	return WithdrawalFailed{
		BaseEvent: event.NewBaseEvent("WithdrawalFailed", walletID, AggregateTypeWallet),
		WalletID:  walletID,
		Reference: reference,
		Reason:    reason,
		Amount:    amount,
		FailedAt:  time.Now().UTC(),
	}
}
