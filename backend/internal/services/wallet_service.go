package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"hustlex/internal/models"
)

// WalletService handles all wallet operations
type WalletService struct {
	db *gorm.DB
}

// NewWalletService creates a new wallet service
func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

// Common errors
var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInsufficientEscrow   = errors.New("insufficient escrow balance")
	ErrWalletLocked         = errors.New("wallet is locked")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrInvalidPIN           = errors.New("invalid transaction PIN")
	ErrPINNotSet            = errors.New("transaction PIN not set")
	ErrSameWallet           = errors.New("cannot transfer to same wallet")
	ErrDailyLimitExceeded   = errors.New("daily transaction limit exceeded")
	ErrMinimumNotMet        = errors.New("minimum amount not met")
	ErrMaximumExceeded      = errors.New("maximum amount exceeded")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

// Transaction limits (in Kobo - 100 Kobo = 1 Naira)
const (
	MinDepositKobo      int64 = 10000      // ₦100 minimum deposit
	MaxDepositKobo      int64 = 500000000  // ₦5,000,000 maximum deposit
	MinWithdrawalKobo   int64 = 50000      // ₦500 minimum withdrawal
	MaxWithdrawalKobo   int64 = 100000000  // ₦1,000,000 maximum withdrawal
	MinTransferKobo     int64 = 10000      // ₦100 minimum transfer
	MaxTransferKobo     int64 = 50000000   // ₦500,000 maximum transfer
	DailyWithdrawLimit  int64 = 500000000  // ₦5,000,000 daily withdrawal limit
	DailyTransferLimit  int64 = 200000000  // ₦2,000,000 daily transfer limit
	WithdrawalFeeKobo   int64 = 5000       // ₦50 withdrawal fee
	TransferFeeKobo     int64 = 1000       // ₦10 transfer fee
)

// GetWallet retrieves a user's wallet
func (s *WalletService) GetWallet(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&wallet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	return &wallet, nil
}

// GetOrCreateWallet gets existing wallet or creates a new one
func (s *WalletService) GetOrCreateWallet(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	wallet, err := s.GetWallet(ctx, userID)
	if err == nil {
		return wallet, nil
	}
	if !errors.Is(err, ErrWalletNotFound) {
		return nil, err
	}

	// Create new wallet
	wallet = &models.Wallet{
		UserID:         userID,
		Balance:        0,
		EscrowBalance:  0,
		SavingsBalance: 0,
		Currency:       "NGN",
		IsLocked:       false,
	}

	if err := s.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWalletByID retrieves a wallet by its ID
func (s *WalletService) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := s.db.WithContext(ctx).
		Where("id = ?", walletID).
		First(&wallet).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}

	return &wallet, nil
}

// DepositInput represents deposit request
type DepositInput struct {
	UserID    uuid.UUID
	AmountKobo int64
	Reference string // Payment gateway reference
	Channel   string // card, bank_transfer, ussd
	Metadata  map[string]interface{}
}

// Deposit adds funds to a wallet
func (s *WalletService) Deposit(ctx context.Context, input DepositInput) (*models.Transaction, error) {
	if input.AmountKobo <= 0 {
		return nil, ErrInvalidAmount
	}
	if input.AmountKobo < MinDepositKobo {
		return nil, fmt.Errorf("%w: minimum deposit is ₦%.2f", ErrMinimumNotMet, float64(MinDepositKobo)/100)
	}
	if input.AmountKobo > MaxDepositKobo {
		return nil, fmt.Errorf("%w: maximum deposit is ₦%.2f", ErrMaximumExceeded, float64(MaxDepositKobo)/100)
	}

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get wallet (create if not exists)
		var wallet models.Wallet
		err := tx.Where("user_id = ?", input.UserID).First(&wallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				wallet = models.Wallet{
					UserID:   input.UserID,
					Currency: "NGN",
				}
				if err := tx.Create(&wallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		if wallet.IsLocked {
			return ErrWalletLocked
		}

		// Update balance
		wallet.Balance += input.AmountKobo
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// Create transaction record
		transaction = &models.Transaction{
			WalletID:        wallet.ID,
			Type:            models.TransactionDeposit,
			AmountKobo:      input.AmountKobo,
			FeeKobo:         0,
			BalanceAfterKobo: wallet.Balance,
			Status:          models.TransactionCompleted,
			Reference:       input.Reference,
			Description:     fmt.Sprintf("Deposit via %s", input.Channel),
			PaymentChannel:  input.Channel,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// WithdrawalInput represents withdrawal request
type WithdrawalInput struct {
	UserID       uuid.UUID
	AmountKobo   int64
	PIN          string
	BankCode     string
	AccountNumber string
	AccountName  string
}

// Withdraw removes funds from a wallet
func (s *WalletService) Withdraw(ctx context.Context, input WithdrawalInput) (*models.Transaction, error) {
	if input.AmountKobo <= 0 {
		return nil, ErrInvalidAmount
	}
	if input.AmountKobo < MinWithdrawalKobo {
		return nil, fmt.Errorf("%w: minimum withdrawal is ₦%.2f", ErrMinimumNotMet, float64(MinWithdrawalKobo)/100)
	}
	if input.AmountKobo > MaxWithdrawalKobo {
		return nil, fmt.Errorf("%w: maximum withdrawal is ₦%.2f", ErrMaximumExceeded, float64(MaxWithdrawalKobo)/100)
	}

	totalAmount := input.AmountKobo + WithdrawalFeeKobo

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get wallet with lock
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", input.UserID).
			First(&wallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}

		if wallet.IsLocked {
			return ErrWalletLocked
		}

		// Verify PIN
		if wallet.TransactionPIN == "" {
			return ErrPINNotSet
		}
		if err := verifyPIN(input.PIN, wallet.TransactionPIN); err != nil {
			// Increment failed attempts
			wallet.PINAttempts++
			if wallet.PINAttempts >= 5 {
				wallet.IsLocked = true
			}
			tx.Save(&wallet)
			return ErrInvalidPIN
		}

		// Reset PIN attempts on success
		wallet.PINAttempts = 0

		// Check daily limit
		dailyTotal, err := s.getDailyWithdrawalTotal(ctx, tx, wallet.ID)
		if err != nil {
			return err
		}
		if dailyTotal+input.AmountKobo > DailyWithdrawLimit {
			return ErrDailyLimitExceeded
		}

		// Check balance
		if wallet.Balance < totalAmount {
			return ErrInsufficientBalance
		}

		// Deduct from balance
		wallet.Balance -= totalAmount
		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// Generate reference
		ref := generateReference("WTH")

		// Create transaction record (pending until bank transfer completes)
		transaction = &models.Transaction{
			WalletID:        wallet.ID,
			Type:            models.TransactionWithdrawal,
			AmountKobo:      input.AmountKobo,
			FeeKobo:         WithdrawalFeeKobo,
			BalanceAfterKobo: wallet.Balance,
			Status:          models.TransactionPending,
			Reference:       ref,
			Description:     fmt.Sprintf("Withdrawal to %s - %s", input.BankCode, maskAccountNumber(input.AccountNumber)),
			BankCode:        input.BankCode,
			AccountNumber:   input.AccountNumber,
			AccountName:     input.AccountName,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// TransferInput represents P2P transfer request
type TransferInput struct {
	FromUserID  uuid.UUID
	ToUserID    uuid.UUID
	AmountKobo  int64
	PIN         string
	Description string
}

// Transfer moves funds between wallets (P2P)
func (s *WalletService) Transfer(ctx context.Context, input TransferInput) (*models.Transaction, error) {
	if input.FromUserID == input.ToUserID {
		return nil, ErrSameWallet
	}
	if input.AmountKobo <= 0 {
		return nil, ErrInvalidAmount
	}
	if input.AmountKobo < MinTransferKobo {
		return nil, fmt.Errorf("%w: minimum transfer is ₦%.2f", ErrMinimumNotMet, float64(MinTransferKobo)/100)
	}
	if input.AmountKobo > MaxTransferKobo {
		return nil, fmt.Errorf("%w: maximum transfer is ₦%.2f", ErrMaximumExceeded, float64(MaxTransferKobo)/100)
	}

	totalDebit := input.AmountKobo + TransferFeeKobo
	var senderTx, receiverTx *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get sender wallet with lock
		var senderWallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", input.FromUserID).
			First(&senderWallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrWalletNotFound
			}
			return err
		}

		if senderWallet.IsLocked {
			return ErrWalletLocked
		}

		// Verify PIN
		if senderWallet.TransactionPIN == "" {
			return ErrPINNotSet
		}
		if err := verifyPIN(input.PIN, senderWallet.TransactionPIN); err != nil {
			senderWallet.PINAttempts++
			if senderWallet.PINAttempts >= 5 {
				senderWallet.IsLocked = true
			}
			tx.Save(&senderWallet)
			return ErrInvalidPIN
		}
		senderWallet.PINAttempts = 0

		// Check daily limit
		dailyTotal, err := s.getDailyTransferTotal(ctx, tx, senderWallet.ID)
		if err != nil {
			return err
		}
		if dailyTotal+input.AmountKobo > DailyTransferLimit {
			return ErrDailyLimitExceeded
		}

		// Check balance
		if senderWallet.Balance < totalDebit {
			return ErrInsufficientBalance
		}

		// Get receiver wallet
		var receiverWallet models.Wallet
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", input.ToUserID).
			First(&receiverWallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create receiver wallet if not exists
				receiverWallet = models.Wallet{
					UserID:   input.ToUserID,
					Currency: "NGN",
				}
				if err := tx.Create(&receiverWallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Execute transfer
		senderWallet.Balance -= totalDebit
		receiverWallet.Balance += input.AmountKobo

		if err := tx.Save(&senderWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&receiverWallet).Error; err != nil {
			return err
		}

		// Generate reference
		ref := generateReference("TRF")

		// Get receiver info
		var receiver models.User
		tx.Select("first_name", "last_name", "phone_number").Where("id = ?", input.ToUserID).First(&receiver)

		// Create sender transaction
		description := input.Description
		if description == "" {
			description = fmt.Sprintf("Transfer to %s %s", receiver.FirstName, receiver.LastName)
		}
		senderTx = &models.Transaction{
			WalletID:         senderWallet.ID,
			Type:             models.TransactionTransfer,
			AmountKobo:       input.AmountKobo,
			FeeKobo:          TransferFeeKobo,
			BalanceAfterKobo: senderWallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        ref,
			Description:      description,
			CounterpartyID:   &input.ToUserID,
		}
		if err := tx.Create(senderTx).Error; err != nil {
			return err
		}

		// Get sender info
		var sender models.User
		tx.Select("first_name", "last_name", "phone_number").Where("id = ?", input.FromUserID).First(&sender)

		// Create receiver transaction
		receiverTx = &models.Transaction{
			WalletID:         receiverWallet.ID,
			Type:             models.TransactionReceived,
			AmountKobo:       input.AmountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: receiverWallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        ref,
			Description:      fmt.Sprintf("Received from %s %s", sender.FirstName, sender.LastName),
			CounterpartyID:   &input.FromUserID,
		}
		return tx.Create(receiverTx).Error
	})

	if err != nil {
		return nil, err
	}

	return senderTx, nil
}

// HoldEscrow moves funds from main balance to escrow
func (s *WalletService) HoldEscrow(ctx context.Context, userID uuid.UUID, amountKobo int64, reference string, description string) (*models.Transaction, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidAmount
	}

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error
		if err != nil {
			return ErrWalletNotFound
		}

		if wallet.IsLocked {
			return ErrWalletLocked
		}

		if wallet.Balance < amountKobo {
			return ErrInsufficientBalance
		}

		// Move to escrow
		wallet.Balance -= amountKobo
		wallet.EscrowBalance += amountKobo

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		// Record transaction
		transaction = &models.Transaction{
			WalletID:         wallet.ID,
			Type:             models.TransactionEscrowHold,
			AmountKobo:       amountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: wallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        reference,
			Description:      description,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// ReleaseEscrow releases escrow funds to recipient
func (s *WalletService) ReleaseEscrow(ctx context.Context, fromUserID, toUserID uuid.UUID, amountKobo int64, platformFeeKobo int64, reference string, description string) (*models.Transaction, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidAmount
	}

	var recipientTx *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Get payer wallet
		var payerWallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", fromUserID).
			First(&payerWallet).Error
		if err != nil {
			return ErrWalletNotFound
		}

		if payerWallet.EscrowBalance < amountKobo {
			return ErrInsufficientEscrow
		}

		// Get recipient wallet
		var recipientWallet models.Wallet
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", toUserID).
			First(&recipientWallet).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				recipientWallet = models.Wallet{
					UserID:   toUserID,
					Currency: "NGN",
				}
				if err := tx.Create(&recipientWallet).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Calculate amounts
		recipientAmount := amountKobo - platformFeeKobo

		// Update balances
		payerWallet.EscrowBalance -= amountKobo
		recipientWallet.Balance += recipientAmount

		if err := tx.Save(&payerWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&recipientWallet).Error; err != nil {
			return err
		}

		// Record payer escrow release
		payerTx := &models.Transaction{
			WalletID:         payerWallet.ID,
			Type:             models.TransactionEscrowRelease,
			AmountKobo:       amountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: payerWallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        reference,
			Description:      description,
			CounterpartyID:   &toUserID,
		}
		if err := tx.Create(payerTx).Error; err != nil {
			return err
		}

		// Record recipient payment
		recipientTx = &models.Transaction{
			WalletID:         recipientWallet.ID,
			Type:             models.TransactionGigPayment,
			AmountKobo:       recipientAmount,
			FeeKobo:          platformFeeKobo,
			BalanceAfterKobo: recipientWallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        reference,
			Description:      description,
			CounterpartyID:   &fromUserID,
		}
		return tx.Create(recipientTx).Error
	})

	if err != nil {
		return nil, err
	}

	return recipientTx, nil
}

// RefundEscrow returns escrowed funds to payer
func (s *WalletService) RefundEscrow(ctx context.Context, userID uuid.UUID, amountKobo int64, reference string, description string) (*models.Transaction, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidAmount
	}

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error
		if err != nil {
			return ErrWalletNotFound
		}

		if wallet.EscrowBalance < amountKobo {
			return ErrInsufficientEscrow
		}

		// Refund from escrow to main balance
		wallet.EscrowBalance -= amountKobo
		wallet.Balance += amountKobo

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction = &models.Transaction{
			WalletID:         wallet.ID,
			Type:             models.TransactionRefund,
			AmountKobo:       amountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: wallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        reference,
			Description:      description,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// MoveToSavings moves funds from main balance to savings
func (s *WalletService) MoveToSavings(ctx context.Context, userID uuid.UUID, amountKobo int64, circleID *uuid.UUID, description string) (*models.Transaction, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidAmount
	}

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error
		if err != nil {
			return ErrWalletNotFound
		}

		if wallet.IsLocked {
			return ErrWalletLocked
		}

		if wallet.Balance < amountKobo {
			return ErrInsufficientBalance
		}

		// Move to savings
		wallet.Balance -= amountKobo
		wallet.SavingsBalance += amountKobo

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		ref := generateReference("SAV")

		transaction = &models.Transaction{
			WalletID:         wallet.ID,
			Type:             models.TransactionSavingsDeposit,
			AmountKobo:       amountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: wallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        ref,
			Description:      description,
			CircleID:         circleID,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// WithdrawFromSavings moves funds from savings back to main balance
func (s *WalletService) WithdrawFromSavings(ctx context.Context, userID uuid.UUID, amountKobo int64, pin string, description string) (*models.Transaction, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidAmount
	}

	var transaction *models.Transaction

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&wallet).Error
		if err != nil {
			return ErrWalletNotFound
		}

		if wallet.IsLocked {
			return ErrWalletLocked
		}

		// Verify PIN
		if wallet.TransactionPIN == "" {
			return ErrPINNotSet
		}
		if err := verifyPIN(pin, wallet.TransactionPIN); err != nil {
			wallet.PINAttempts++
			if wallet.PINAttempts >= 5 {
				wallet.IsLocked = true
			}
			tx.Save(&wallet)
			return ErrInvalidPIN
		}
		wallet.PINAttempts = 0

		if wallet.SavingsBalance < amountKobo {
			return fmt.Errorf("insufficient savings balance")
		}

		// Move from savings to main
		wallet.SavingsBalance -= amountKobo
		wallet.Balance += amountKobo

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		ref := generateReference("SWD")

		transaction = &models.Transaction{
			WalletID:         wallet.ID,
			Type:             models.TransactionSavingsWithdrawal,
			AmountKobo:       amountKobo,
			FeeKobo:          0,
			BalanceAfterKobo: wallet.Balance,
			Status:           models.TransactionCompleted,
			Reference:        ref,
			Description:      description,
		}

		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactions retrieves transaction history for a wallet
func (s *WalletService) GetTransactions(ctx context.Context, userID uuid.UUID, filter TransactionFilter) ([]models.Transaction, int64, error) {
	wallet, err := s.GetWallet(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	query := s.db.WithContext(ctx).Model(&models.Transaction{}).Where("wallet_id = ?", wallet.ID)

	// Apply filters
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if !filter.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	var transactions []models.Transaction
	err = query.
		Order("created_at DESC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&transactions).Error

	return transactions, total, err
}

// TransactionFilter defines transaction query filters
type TransactionFilter struct {
	Type      models.TransactionType
	Status    models.TransactionStatus
	StartDate time.Time
	EndDate   time.Time
	Offset    int
	Limit     int
}

// GetTransaction retrieves a single transaction by reference
func (s *WalletService) GetTransaction(ctx context.Context, userID uuid.UUID, reference string) (*models.Transaction, error) {
	wallet, err := s.GetWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	var transaction models.Transaction
	err = s.db.WithContext(ctx).
		Where("wallet_id = ? AND reference = ?", wallet.ID, reference).
		First(&transaction).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	return &transaction, nil
}

// GetWalletSummary retrieves wallet balances and recent activity summary
func (s *WalletService) GetWalletSummary(ctx context.Context, userID uuid.UUID) (*WalletSummary, error) {
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get transaction stats for the month
	startOfMonth := time.Now().UTC().Truncate(24 * time.Hour).AddDate(0, 0, -time.Now().Day()+1)

	var monthlyInflow, monthlyOutflow int64

	// Calculate inflow (deposits + received + gig payments)
	s.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("wallet_id = ? AND created_at >= ? AND status = ? AND type IN ?",
			wallet.ID, startOfMonth, models.TransactionCompleted,
			[]models.TransactionType{models.TransactionDeposit, models.TransactionReceived, models.TransactionGigPayment}).
		Select("COALESCE(SUM(amount_kobo), 0)").
		Scan(&monthlyInflow)

	// Calculate outflow (withdrawals + transfers + escrow holds)
	s.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("wallet_id = ? AND created_at >= ? AND status = ? AND type IN ?",
			wallet.ID, startOfMonth, models.TransactionCompleted,
			[]models.TransactionType{models.TransactionWithdrawal, models.TransactionTransfer, models.TransactionEscrowHold}).
		Select("COALESCE(SUM(amount_kobo + fee_kobo), 0)").
		Scan(&monthlyOutflow)

	return &WalletSummary{
		WalletID:       wallet.ID,
		Balance:        wallet.Balance,
		EscrowBalance:  wallet.EscrowBalance,
		SavingsBalance: wallet.SavingsBalance,
		TotalBalance:   wallet.Balance + wallet.EscrowBalance + wallet.SavingsBalance,
		Currency:       wallet.Currency,
		IsLocked:       wallet.IsLocked,
		MonthlyInflow:  monthlyInflow,
		MonthlyOutflow: monthlyOutflow,
	}, nil
}

// WalletSummary represents wallet overview data
type WalletSummary struct {
	WalletID       uuid.UUID `json:"wallet_id"`
	Balance        int64     `json:"balance"`
	EscrowBalance  int64     `json:"escrow_balance"`
	SavingsBalance int64     `json:"savings_balance"`
	TotalBalance   int64     `json:"total_balance"`
	Currency       string    `json:"currency"`
	IsLocked       bool      `json:"is_locked"`
	MonthlyInflow  int64     `json:"monthly_inflow"`
	MonthlyOutflow int64     `json:"monthly_outflow"`
}

// Helper functions

func (s *WalletService) getDailyWithdrawalTotal(ctx context.Context, tx *gorm.DB, walletID uuid.UUID) (int64, error) {
	startOfDay := time.Now().UTC().Truncate(24 * time.Hour)

	var total int64
	err := tx.Model(&models.Transaction{}).
		Where("wallet_id = ? AND type = ? AND created_at >= ? AND status != ?",
			walletID, models.TransactionWithdrawal, startOfDay, models.TransactionFailed).
		Select("COALESCE(SUM(amount_kobo), 0)").
		Scan(&total).Error

	return total, err
}

func (s *WalletService) getDailyTransferTotal(ctx context.Context, tx *gorm.DB, walletID uuid.UUID) (int64, error) {
	startOfDay := time.Now().UTC().Truncate(24 * time.Hour)

	var total int64
	err := tx.Model(&models.Transaction{}).
		Where("wallet_id = ? AND type = ? AND created_at >= ? AND status != ?",
			walletID, models.TransactionTransfer, startOfDay, models.TransactionFailed).
		Select("COALESCE(SUM(amount_kobo), 0)").
		Scan(&total).Error

	return total, err
}

func generateReference(prefix string) string {
	timestamp := time.Now().Format("20060102150405")
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%s%s%06d", prefix, timestamp, n)
}

func maskAccountNumber(accountNo string) string {
	if len(accountNo) <= 4 {
		return "****"
	}
	return "****" + accountNo[len(accountNo)-4:]
}

// verifyPIN securely compares the provided PIN against the bcrypt hash
func verifyPIN(providedPIN, storedHash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPIN)); err != nil {
		return ErrInvalidPIN
	}
	return nil
}
