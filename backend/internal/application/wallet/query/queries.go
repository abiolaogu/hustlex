package query

import (
	"context"
	"time"

	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/repository"
)

// GetWallet retrieves wallet information
type GetWallet struct {
	UserID string
}

// GetWalletResult is the result of GetWallet query
type GetWalletResult struct {
	WalletID         string    `json:"wallet_id"`
	UserID           string    `json:"user_id"`
	AvailableBalance int64     `json:"available_balance"`
	EscrowBalance    int64     `json:"escrow_balance"`
	SavingsBalance   int64     `json:"savings_balance"`
	TotalBalance     int64     `json:"total_balance"`
	Currency         string    `json:"currency"`
	Status           string    `json:"status"`
	HasPIN           bool      `json:"has_pin"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// GetTransactions retrieves transaction history
type GetTransactions struct {
	UserID    string
	Type      *string
	Status    *string
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

// GetTransactionsResult is the result of GetTransactions query
type GetTransactionsResult struct {
	Transactions []TransactionDTO `json:"transactions"`
	Total        int64            `json:"total"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
	TotalPages   int64            `json:"total_pages"`
}

// TransactionDTO represents a transaction for API responses
type TransactionDTO struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	Amount        int64                  `json:"amount"`
	Fee           int64                  `json:"fee"`
	Currency      string                 `json:"currency"`
	BalanceAfter  int64                  `json:"balance_after"`
	Status        string                 `json:"status"`
	Reference     string                 `json:"reference"`
	Description   string                 `json:"description"`
	Counterparty  *string                `json:"counterparty,omitempty"`
	BankCode      *string                `json:"bank_code,omitempty"`
	AccountNumber *string                `json:"account_number,omitempty"`
	AccountName   *string                `json:"account_name,omitempty"`
	FailureReason *string                `json:"failure_reason,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt     string                 `json:"created_at"`
}

// GetTransaction retrieves a single transaction
type GetTransaction struct {
	UserID    string
	Reference string
}

// GetBankAccounts retrieves saved bank accounts
type GetBankAccounts struct {
	UserID string
}

// BankAccountDTO represents a bank account for API responses
type BankAccountDTO struct {
	ID            string `json:"id"`
	BankCode      string `json:"bank_code"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	IsDefault     bool   `json:"is_default"`
	IsVerified    bool   `json:"is_verified"`
}

// GetWalletStats retrieves wallet statistics
type GetWalletStats struct {
	UserID string
	Period string // week, month, year, all
}

// WalletStatsResult contains wallet statistics
type WalletStatsResult struct {
	Period           string `json:"period"`
	TotalInflow      int64  `json:"total_inflow"`
	TotalOutflow     int64  `json:"total_outflow"`
	TotalDeposits    int64  `json:"total_deposits"`
	TotalWithdrawals int64  `json:"total_withdrawals"`
	TotalTransfersIn int64  `json:"total_transfers_in"`
	TotalTransfersOut int64 `json:"total_transfers_out"`
	TotalGigEarnings int64  `json:"total_gig_earnings"`
	TransactionCount int    `json:"transaction_count"`
}

// WalletQueryHandler handles wallet queries
type WalletQueryHandler struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	bankAccountRepo repository.BankAccountRepository
}

// NewWalletQueryHandler creates a new query handler
func NewWalletQueryHandler(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	bankAccountRepo repository.BankAccountRepository,
) *WalletQueryHandler {
	return &WalletQueryHandler{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		bankAccountRepo: bankAccountRepo,
	}
}

// HandleGetWallet retrieves wallet information
func (h *WalletQueryHandler) HandleGetWallet(ctx context.Context, q GetWallet) (*GetWalletResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	wallet, err := h.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetWalletResult{
		WalletID:         wallet.ID().String(),
		UserID:           wallet.UserID().String(),
		AvailableBalance: wallet.AvailableBalance().Amount(),
		EscrowBalance:    wallet.EscrowBalance().Amount(),
		SavingsBalance:   wallet.SavingsBalance().Amount(),
		TotalBalance:     wallet.TotalBalance().Amount(),
		Currency:         string(wallet.Currency()),
		Status:           string(wallet.Status()),
		HasPIN:           wallet.HasPIN(),
		CreatedAt:        wallet.CreatedAt(),
		UpdatedAt:        wallet.UpdatedAt(),
	}, nil
}

// HandleGetTransactions retrieves transaction history
func (h *WalletQueryHandler) HandleGetTransactions(ctx context.Context, q GetTransactions) (*GetTransactionsResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	wallet, err := h.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Build filter
	filter := repository.TransactionFilter{
		Offset: (q.Page - 1) * q.Limit,
		Limit:  q.Limit,
	}

	if q.Type != nil {
		txType := repository.TransactionType(*q.Type)
		filter.Type = &txType
	}

	if q.Status != nil {
		txStatus := repository.TransactionStatus(*q.Status)
		filter.Status = &txStatus
	}

	// Get transactions
	transactions, total, err := h.transactionRepo.FindByWalletID(ctx, wallet.ID(), filter)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	dtos := make([]TransactionDTO, len(transactions))
	for i, tx := range transactions {
		dtos[i] = TransactionDTO{
			ID:            tx.ID,
			Type:          string(tx.Type),
			Amount:        tx.Amount,
			Fee:           tx.Fee,
			Currency:      tx.Currency,
			BalanceAfter:  tx.BalanceAfter,
			Status:        string(tx.Status),
			Reference:     tx.Reference,
			Description:   tx.Description,
			Counterparty:  tx.CounterpartyID,
			BankCode:      tx.BankCode,
			AccountNumber: tx.AccountNumber,
			AccountName:   tx.AccountName,
			FailureReason: tx.FailureReason,
			Metadata:      tx.Metadata,
			CreatedAt:     tx.CreatedAt,
		}
	}

	totalPages := total / int64(q.Limit)
	if total%int64(q.Limit) > 0 {
		totalPages++
	}

	return &GetTransactionsResult{
		Transactions: dtos,
		Total:        total,
		Page:         q.Page,
		Limit:        q.Limit,
		TotalPages:   totalPages,
	}, nil
}

// HandleGetBankAccounts retrieves saved bank accounts
func (h *WalletQueryHandler) HandleGetBankAccounts(ctx context.Context, q GetBankAccounts) ([]BankAccountDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	accounts, err := h.bankAccountRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]BankAccountDTO, len(accounts))
	for i, acc := range accounts {
		dtos[i] = BankAccountDTO{
			ID:            acc.ID,
			BankCode:      acc.BankCode,
			BankName:      acc.BankName,
			AccountNumber: acc.AccountNumber,
			AccountName:   acc.AccountName,
			IsDefault:     acc.IsDefault,
			IsVerified:    acc.IsVerified,
		}
	}

	return dtos, nil
}
