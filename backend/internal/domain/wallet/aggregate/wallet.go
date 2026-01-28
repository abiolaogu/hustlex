package aggregate

import (
	"errors"
	"time"

	"hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
	walletEvent "hustlex/internal/domain/wallet/event"
)

// Domain errors
var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInsufficientEscrow   = errors.New("insufficient escrow balance")
	ErrWalletLocked         = errors.New("wallet is locked")
	ErrInvalidAmount        = errors.New("amount must be positive")
	ErrCurrencyMismatch     = errors.New("currency mismatch")
	ErrPINNotSet            = errors.New("transaction PIN not set")
	ErrInvalidPIN           = errors.New("invalid transaction PIN")
	ErrWalletSuspended      = errors.New("wallet is suspended")
	ErrDailyLimitExceeded   = errors.New("daily transaction limit exceeded")
)

// WalletStatus represents the current state of a wallet
type WalletStatus string

const (
	WalletStatusActive    WalletStatus = "active"
	WalletStatusLocked    WalletStatus = "locked"
	WalletStatusSuspended WalletStatus = "suspended"
)

// Wallet is the aggregate root for wallet operations
type Wallet struct {
	event.AggregateRoot

	id               valueobject.WalletID
	userID           valueobject.UserID
	availableBalance valueobject.Money
	escrowBalance    valueobject.Money
	savingsBalance   valueobject.Money
	ledgerBalance    valueobject.Money
	currency         valueobject.Currency
	status           WalletStatus
	pinHash          string
	pinAttempts      int
	createdAt        time.Time
	updatedAt        time.Time
	version          int64
}

// NewWallet creates a new wallet for a user
func NewWallet(userID valueobject.UserID, currency valueobject.Currency) *Wallet {
	now := time.Now().UTC()
	walletID := valueobject.GenerateWalletID()

	wallet := &Wallet{
		id:               walletID,
		userID:           userID,
		availableBalance: valueobject.Zero(currency),
		escrowBalance:    valueobject.Zero(currency),
		savingsBalance:   valueobject.Zero(currency),
		ledgerBalance:    valueobject.Zero(currency),
		currency:         currency,
		status:           WalletStatusActive,
		createdAt:        now,
		updatedAt:        now,
		version:          1,
	}

	wallet.RecordEvent(walletEvent.NewWalletCreated(
		walletID.String(),
		userID.String(),
		string(currency),
	))

	return wallet
}

// Reconstitute recreates a wallet from persistence (no events recorded)
func Reconstitute(
	id valueobject.WalletID,
	userID valueobject.UserID,
	availableBalance, escrowBalance, savingsBalance, ledgerBalance valueobject.Money,
	currency valueobject.Currency,
	status WalletStatus,
	pinHash string,
	pinAttempts int,
	createdAt, updatedAt time.Time,
	version int64,
) *Wallet {
	return &Wallet{
		id:               id,
		userID:           userID,
		availableBalance: availableBalance,
		escrowBalance:    escrowBalance,
		savingsBalance:   savingsBalance,
		ledgerBalance:    ledgerBalance,
		currency:         currency,
		status:           status,
		pinHash:          pinHash,
		pinAttempts:      pinAttempts,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		version:          version,
	}
}

// Credit adds funds to the wallet
func (w *Wallet) Credit(amount valueobject.Money, source, reference, description string) error {
	if err := w.validateActive(); err != nil {
		return err
	}

	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if !amount.IsPositive() {
		return ErrInvalidAmount
	}

	newAvailable, err := w.availableBalance.Add(amount)
	if err != nil {
		return err
	}

	newLedger, err := w.ledgerBalance.Add(amount)
	if err != nil {
		return err
	}

	w.availableBalance = newAvailable
	w.ledgerBalance = newLedger
	w.touch()

	w.RecordEvent(walletEvent.NewWalletCredited(
		w.id.String(),
		w.userID.String(),
		amount.Amount(),
		string(amount.Currency()),
		source,
		reference,
		description,
		newAvailable.Amount(),
	))

	return nil
}

// Debit removes funds from the wallet
func (w *Wallet) Debit(amount valueobject.Money, destination, reference, description string, fee valueobject.Money) error {
	if err := w.validateActive(); err != nil {
		return err
	}

	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if !amount.IsPositive() {
		return ErrInvalidAmount
	}

	totalDebit := amount.MustAdd(fee)

	if w.availableBalance.LessThan(totalDebit) {
		return ErrInsufficientFunds
	}

	newAvailable := w.availableBalance.MustSubtract(totalDebit)
	newLedger := w.ledgerBalance.MustSubtract(totalDebit)

	w.availableBalance = newAvailable
	w.ledgerBalance = newLedger
	w.touch()

	w.RecordEvent(walletEvent.NewWalletDebited(
		w.id.String(),
		w.userID.String(),
		amount.Amount(),
		string(amount.Currency()),
		destination,
		reference,
		description,
		fee.Amount(),
		newAvailable.Amount(),
	))

	return nil
}

// HoldInEscrow moves funds from available to escrow
func (w *Wallet) HoldInEscrow(amount valueobject.Money, reference, reason string) error {
	if err := w.validateActive(); err != nil {
		return err
	}

	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if !amount.IsPositive() {
		return ErrInvalidAmount
	}

	if w.availableBalance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	newAvailable := w.availableBalance.MustSubtract(amount)
	newEscrow := w.escrowBalance.MustAdd(amount)

	w.availableBalance = newAvailable
	w.escrowBalance = newEscrow
	w.touch()

	w.RecordEvent(walletEvent.NewFundsHeldInEscrow(
		w.id.String(),
		w.userID.String(),
		amount.Amount(),
		reference,
		reason,
		newAvailable.Amount(),
		newEscrow.Amount(),
	))

	return nil
}

// ReleaseFromEscrow releases escrowed funds
func (w *Wallet) ReleaseFromEscrow(amount valueobject.Money, reference string, toWallet bool, recipientID string) error {
	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if !amount.IsPositive() {
		return ErrInvalidAmount
	}

	if w.escrowBalance.LessThan(amount) {
		return ErrInsufficientEscrow
	}

	newEscrow := w.escrowBalance.MustSubtract(amount)
	w.escrowBalance = newEscrow

	if toWallet {
		// Return to the same wallet
		newAvailable := w.availableBalance.MustAdd(amount)
		w.availableBalance = newAvailable
	} else {
		// Released to another party - reduce ledger balance
		newLedger := w.ledgerBalance.MustSubtract(amount)
		w.ledgerBalance = newLedger
	}

	w.touch()

	w.RecordEvent(walletEvent.NewFundsReleasedFromEscrow(
		w.id.String(),
		w.userID.String(),
		amount.Amount(),
		reference,
		recipientID,
		toWallet,
		newEscrow.Amount(),
	))

	return nil
}

// MoveToSavings moves funds from available to savings balance
func (w *Wallet) MoveToSavings(amount valueobject.Money) error {
	if err := w.validateActive(); err != nil {
		return err
	}

	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if w.availableBalance.LessThan(amount) {
		return ErrInsufficientFunds
	}

	w.availableBalance = w.availableBalance.MustSubtract(amount)
	w.savingsBalance = w.savingsBalance.MustAdd(amount)
	w.touch()

	return nil
}

// WithdrawFromSavings moves funds from savings to available
func (w *Wallet) WithdrawFromSavings(amount valueobject.Money) error {
	if err := w.validateActive(); err != nil {
		return err
	}

	if err := w.validateCurrency(amount); err != nil {
		return err
	}

	if w.savingsBalance.LessThan(amount) {
		return errors.New("insufficient savings balance")
	}

	w.savingsBalance = w.savingsBalance.MustSubtract(amount)
	w.availableBalance = w.availableBalance.MustAdd(amount)
	w.touch()

	return nil
}

// Lock prevents any transactions on the wallet
func (w *Wallet) Lock(reason string) {
	w.status = WalletStatusLocked
	w.touch()

	w.RecordEvent(walletEvent.NewWalletLocked(
		w.id.String(),
		w.userID.String(),
		reason,
	))
}

// Unlock allows transactions again
func (w *Wallet) Unlock() {
	w.status = WalletStatusActive
	w.pinAttempts = 0
	w.touch()

	w.RecordEvent(walletEvent.NewWalletUnlocked(
		w.id.String(),
		w.userID.String(),
	))
}

// SetPIN sets the transaction PIN hash
func (w *Wallet) SetPIN(pinHash string) {
	w.pinHash = pinHash
	w.pinAttempts = 0
	w.touch()

	w.RecordEvent(walletEvent.NewTransactionPINSet(
		w.id.String(),
		w.userID.String(),
	))
}

// RecordFailedPINAttempt records a failed PIN verification
func (w *Wallet) RecordFailedPINAttempt(maxAttempts int) {
	w.pinAttempts++
	if w.pinAttempts >= maxAttempts {
		w.Lock("Too many failed PIN attempts")
	}
	w.touch()
}

// ResetPINAttempts resets the PIN attempt counter
func (w *Wallet) ResetPINAttempts() {
	w.pinAttempts = 0
	w.touch()
}

// Getters
func (w *Wallet) ID() valueobject.WalletID              { return w.id }
func (w *Wallet) UserID() valueobject.UserID            { return w.userID }
func (w *Wallet) AvailableBalance() valueobject.Money   { return w.availableBalance }
func (w *Wallet) EscrowBalance() valueobject.Money      { return w.escrowBalance }
func (w *Wallet) SavingsBalance() valueobject.Money     { return w.savingsBalance }
func (w *Wallet) LedgerBalance() valueobject.Money      { return w.ledgerBalance }
func (w *Wallet) TotalBalance() valueobject.Money {
	return w.availableBalance.MustAdd(w.escrowBalance).MustAdd(w.savingsBalance)
}
func (w *Wallet) Currency() valueobject.Currency { return w.currency }
func (w *Wallet) Status() WalletStatus           { return w.status }
func (w *Wallet) PINHash() string                { return w.pinHash }
func (w *Wallet) HasPIN() bool                   { return w.pinHash != "" }
func (w *Wallet) PINAttempts() int               { return w.pinAttempts }
func (w *Wallet) CreatedAt() time.Time           { return w.createdAt }
func (w *Wallet) UpdatedAt() time.Time           { return w.updatedAt }
func (w *Wallet) Version() int64                 { return w.version }
func (w *Wallet) IsActive() bool                 { return w.status == WalletStatusActive }
func (w *Wallet) IsLocked() bool                 { return w.status == WalletStatusLocked }

// Private helpers
func (w *Wallet) touch() {
	w.updatedAt = time.Now().UTC()
	w.version++
}

func (w *Wallet) validateActive() error {
	switch w.status {
	case WalletStatusLocked:
		return ErrWalletLocked
	case WalletStatusSuspended:
		return ErrWalletSuspended
	}
	return nil
}

func (w *Wallet) validateCurrency(amount valueobject.Money) error {
	if amount.Currency() != w.currency {
		return ErrCurrencyMismatch
	}
	return nil
}
