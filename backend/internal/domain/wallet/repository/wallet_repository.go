package repository

import (
	"context"

	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
)

// WalletRepository defines the interface for wallet persistence
// This is a PORT in hexagonal architecture - domain defines what it needs,
// infrastructure provides the implementation
type WalletRepository interface {
	// FindByID retrieves a wallet by its unique identifier
	FindByID(ctx context.Context, id valueobject.WalletID) (*aggregate.Wallet, error)

	// FindByUserID retrieves a wallet by the owning user's ID
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.Wallet, error)

	// Exists checks if a wallet exists for the given user
	Exists(ctx context.Context, userID valueobject.UserID) (bool, error)

	// Save persists the wallet aggregate
	// Uses optimistic locking via the version field
	// Returns an error if the version has changed (concurrent modification)
	Save(ctx context.Context, wallet *aggregate.Wallet) error

	// SaveWithEvents saves the wallet and publishes domain events atomically
	// This ensures events are only published if persistence succeeds
	SaveWithEvents(ctx context.Context, wallet *aggregate.Wallet) error
}

// TransactionRepository defines the interface for transaction persistence
type TransactionRepository interface {
	// FindByReference retrieves a transaction by its reference
	FindByReference(ctx context.Context, reference string) (*Transaction, error)

	// FindByWalletID retrieves transactions for a wallet with pagination
	FindByWalletID(ctx context.Context, walletID valueobject.WalletID, filter TransactionFilter) ([]Transaction, int64, error)

	// Save persists a transaction record
	Save(ctx context.Context, tx *Transaction) error
}

// Transaction represents a persisted transaction record
type Transaction struct {
	ID               string
	WalletID         string
	Type             TransactionType
	Amount           int64
	Fee              int64
	Currency         string
	BalanceAfter     int64
	Status           TransactionStatus
	Reference        string
	Description      string
	CounterpartyID   *string
	BankCode         *string
	AccountNumber    *string
	AccountName      *string
	PaymentChannel   *string
	FailureReason    *string
	Metadata         map[string]interface{}
	CreatedAt        string
	UpdatedAt        string
}

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypeDeposit           TransactionType = "deposit"
	TransactionTypeWithdrawal        TransactionType = "withdrawal"
	TransactionTypeTransferOut       TransactionType = "transfer_out"
	TransactionTypeTransferIn        TransactionType = "transfer_in"
	TransactionTypeEscrowHold        TransactionType = "escrow_hold"
	TransactionTypeEscrowRelease     TransactionType = "escrow_release"
	TransactionTypeGigPayment        TransactionType = "gig_payment"
	TransactionTypeSavingsDeposit    TransactionType = "savings_deposit"
	TransactionTypeSavingsWithdrawal TransactionType = "savings_withdrawal"
	TransactionTypeContribution      TransactionType = "contribution"
	TransactionTypePayout            TransactionType = "payout"
	TransactionTypeLoanDisbursement  TransactionType = "loan_disbursement"
	TransactionTypeLoanRepayment     TransactionType = "loan_repayment"
	TransactionTypeRefund            TransactionType = "refund"
	TransactionTypeFee               TransactionType = "fee"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// TransactionFilter defines filters for querying transactions
type TransactionFilter struct {
	Type      *TransactionType
	Status    *TransactionStatus
	StartDate *string
	EndDate   *string
	MinAmount *int64
	MaxAmount *int64
	Offset    int
	Limit     int
}

// BankAccountRepository defines the interface for bank account persistence
type BankAccountRepository interface {
	// FindByID retrieves a bank account by ID
	FindByID(ctx context.Context, id string) (*BankAccount, error)

	// FindByUserID retrieves all bank accounts for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]BankAccount, error)

	// FindDefault retrieves the default bank account for a user
	FindDefault(ctx context.Context, userID valueobject.UserID) (*BankAccount, error)

	// Save persists a bank account
	Save(ctx context.Context, account *BankAccount) error

	// Delete removes a bank account
	Delete(ctx context.Context, id string) error

	// SetDefault sets a bank account as default
	SetDefault(ctx context.Context, userID valueobject.UserID, accountID string) error
}

// BankAccount represents a saved bank account
type BankAccount struct {
	ID            string
	UserID        string
	BankCode      string
	BankName      string
	AccountNumber string
	AccountName   string
	IsDefault     bool
	IsVerified    bool
	CreatedAt     string
	UpdatedAt     string
}
