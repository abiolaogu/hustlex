package command

import (
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

// CreateWallet creates a new wallet for a user
type CreateWallet struct {
	UserID    string
	Currency  string
	RequestID string
}

// Deposit adds funds to a wallet
type Deposit struct {
	WalletID    string
	Amount      int64
	Currency    string
	Source      string // deposit, refund, etc.
	Reference   string
	Description string
	Channel     string // card, bank_transfer, ussd
	RequestedBy string
	RequestedAt time.Time
}

// Withdraw removes funds from a wallet to a bank account
type Withdraw struct {
	WalletID      string
	Amount        int64
	Currency      string
	BankCode      string
	AccountNumber string
	AccountName   string
	PIN           string
	Reference     string
	RequestedBy   string
	RequestedAt   time.Time
}

// Transfer moves funds between wallets
type Transfer struct {
	FromUserID     string
	ToUserPhone    string // Recipient identified by phone
	Amount         int64
	Currency       string
	Description    string
	PIN            string
	Reference      string
	RequestedBy    string
	RequestedAt    time.Time
}

// HoldEscrow moves funds to escrow for a gig
type HoldEscrow struct {
	UserID      string
	Amount      int64
	Currency    string
	ContractID  string
	Description string
	RequestedBy string
}

// ReleaseEscrow releases escrowed funds
type ReleaseEscrow struct {
	PayerUserID     string
	RecipientUserID string
	Amount          int64
	PlatformFee     int64
	Currency        string
	ContractID      string
	Description     string
	RequestedBy     string
}

// RefundEscrow returns escrowed funds to the payer
type RefundEscrow struct {
	UserID      string
	Amount      int64
	Currency    string
	ContractID  string
	Reason      string
	RequestedBy string
}

// SetTransactionPIN sets or updates the wallet PIN
type SetTransactionPIN struct {
	UserID     string
	PIN        string
	OldPIN     string // Required if changing existing PIN
	RequestedBy string
}

// LockWallet locks a wallet
type LockWallet struct {
	WalletID    string
	Reason      string
	RequestedBy string
}

// UnlockWallet unlocks a wallet
type UnlockWallet struct {
	WalletID    string
	RequestedBy string
}

// AddBankAccount adds a bank account for withdrawals
type AddBankAccount struct {
	UserID        string
	BankCode      string
	AccountNumber string
	AccountName   string
	IsDefault     bool
	RequestedBy   string
}

// RemoveBankAccount removes a saved bank account
type RemoveBankAccount struct {
	UserID      string
	AccountID   string
	RequestedBy string
}

// SetDefaultBankAccount sets a bank account as default
type SetDefaultBankAccount struct {
	UserID      string
	AccountID   string
	RequestedBy string
}

// Result types

// WalletResult is the result of wallet operations
type WalletResult struct {
	WalletID         string
	AvailableBalance int64
	EscrowBalance    int64
	SavingsBalance   int64
	Currency         string
}

// DepositResult is the result of a deposit
type DepositResult struct {
	TransactionID string
	Reference     string
	NewBalance    int64
	ProcessedAt   time.Time
	// For card payments
	PaymentURL    string
	AccessCode    string
}

// WithdrawResult is the result of a withdrawal
type WithdrawResult struct {
	TransactionID string
	Reference     string
	Status        string
	NewBalance    int64
	Fee           int64
	ProcessedAt   time.Time
}

// TransferResult is the result of a transfer
type TransferResult struct {
	TransactionID string
	Reference     string
	RecipientName string
	Amount        int64
	Fee           int64
	NewBalance    int64
	ProcessedAt   time.Time
}

// Money helper to convert command values to domain value objects
func (d Deposit) GetMoney() (valueobject.Money, error) {
	return valueobject.NewMoney(d.Amount, valueobject.Currency(d.Currency))
}

func (w Withdraw) GetMoney() (valueobject.Money, error) {
	return valueobject.NewMoney(w.Amount, valueobject.Currency(w.Currency))
}

func (t Transfer) GetMoney() (valueobject.Money, error) {
	return valueobject.NewMoney(t.Amount, valueobject.Currency(t.Currency))
}
