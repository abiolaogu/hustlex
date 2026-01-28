package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SupportedCurrency represents supported currencies for multi-currency wallets
type SupportedCurrency string

const (
	CurrencyNGN SupportedCurrency = "NGN" // Nigerian Naira
	CurrencyGBP SupportedCurrency = "GBP" // British Pound
	CurrencyUSD SupportedCurrency = "USD" // US Dollar
	CurrencyEUR SupportedCurrency = "EUR" // Euro
	CurrencyCAD SupportedCurrency = "CAD" // Canadian Dollar
	CurrencyGHS SupportedCurrency = "GHS" // Ghanaian Cedi
	CurrencyKES SupportedCurrency = "KES" // Kenyan Shilling
)

// CurrencyBalance represents a balance in a specific currency
type CurrencyBalance struct {
	Currency         SupportedCurrency `json:"currency" db:"currency"`
	AvailableBalance decimal.Decimal   `json:"available_balance" db:"available_balance"`
	PendingBalance   decimal.Decimal   `json:"pending_balance" db:"pending_balance"`
	LockedBalance    decimal.Decimal   `json:"locked_balance" db:"locked_balance"`
	TotalBalance     decimal.Decimal   `json:"total_balance" db:"total_balance"`
	UpdatedAt        time.Time         `json:"updated_at" db:"updated_at"`
}

// MultiCurrencyWallet extends the base wallet with multi-currency support
type MultiCurrencyWallet struct {
	ID                uuid.UUID                        `json:"id" db:"id"`
	UserID            uuid.UUID                        `json:"user_id" db:"user_id"`
	BaseWalletID      uuid.UUID                        `json:"base_wallet_id" db:"base_wallet_id"`
	PrimaryCurrency   SupportedCurrency                `json:"primary_currency" db:"primary_currency"`
	Balances          map[SupportedCurrency]*CurrencyBalance `json:"balances" db:"-"`
	IsActive          bool                             `json:"is_active" db:"is_active"`
	CreatedAt         time.Time                        `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time                        `json:"updated_at" db:"updated_at"`
}

// NewMultiCurrencyWallet creates a new multi-currency wallet
func NewMultiCurrencyWallet(userID, baseWalletID uuid.UUID, primaryCurrency SupportedCurrency) *MultiCurrencyWallet {
	now := time.Now()
	wallet := &MultiCurrencyWallet{
		ID:              uuid.New(),
		UserID:          userID,
		BaseWalletID:    baseWalletID,
		PrimaryCurrency: primaryCurrency,
		Balances:        make(map[SupportedCurrency]*CurrencyBalance),
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Initialize primary currency balance
	wallet.Balances[primaryCurrency] = &CurrencyBalance{
		Currency:         primaryCurrency,
		AvailableBalance: decimal.Zero,
		PendingBalance:   decimal.Zero,
		LockedBalance:    decimal.Zero,
		TotalBalance:     decimal.Zero,
		UpdatedAt:        now,
	}

	return wallet
}

// EnableCurrency enables a new currency for the wallet
func (w *MultiCurrencyWallet) EnableCurrency(currency SupportedCurrency) error {
	if _, exists := w.Balances[currency]; exists {
		return errors.New("currency already enabled")
	}

	now := time.Now()
	w.Balances[currency] = &CurrencyBalance{
		Currency:         currency,
		AvailableBalance: decimal.Zero,
		PendingBalance:   decimal.Zero,
		LockedBalance:    decimal.Zero,
		TotalBalance:     decimal.Zero,
		UpdatedAt:        now,
	}
	w.UpdatedAt = now
	return nil
}

// GetBalance returns the balance for a specific currency
func (w *MultiCurrencyWallet) GetBalance(currency SupportedCurrency) (*CurrencyBalance, error) {
	balance, exists := w.Balances[currency]
	if !exists {
		return nil, errors.New("currency not enabled for this wallet")
	}
	return balance, nil
}

// Credit adds funds to a specific currency balance
func (w *MultiCurrencyWallet) Credit(currency SupportedCurrency, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		// Auto-enable currency if not exists
		if err := w.EnableCurrency(currency); err != nil {
			return err
		}
		balance = w.Balances[currency]
	}

	now := time.Now()
	balance.AvailableBalance = balance.AvailableBalance.Add(amount)
	balance.TotalBalance = balance.AvailableBalance.Add(balance.PendingBalance).Add(balance.LockedBalance)
	balance.UpdatedAt = now
	w.UpdatedAt = now
	return nil
}

// Debit removes funds from a specific currency balance
func (w *MultiCurrencyWallet) Debit(currency SupportedCurrency, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		return err
	}

	if balance.AvailableBalance.LessThan(amount) {
		return errors.New("insufficient funds")
	}

	now := time.Now()
	balance.AvailableBalance = balance.AvailableBalance.Sub(amount)
	balance.TotalBalance = balance.AvailableBalance.Add(balance.PendingBalance).Add(balance.LockedBalance)
	balance.UpdatedAt = now
	w.UpdatedAt = now
	return nil
}

// Lock locks funds for a pending transaction
func (w *MultiCurrencyWallet) Lock(currency SupportedCurrency, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		return err
	}

	if balance.AvailableBalance.LessThan(amount) {
		return errors.New("insufficient funds to lock")
	}

	now := time.Now()
	balance.AvailableBalance = balance.AvailableBalance.Sub(amount)
	balance.LockedBalance = balance.LockedBalance.Add(amount)
	balance.UpdatedAt = now
	w.UpdatedAt = now
	return nil
}

// Unlock releases locked funds
func (w *MultiCurrencyWallet) Unlock(currency SupportedCurrency, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		return err
	}

	if balance.LockedBalance.LessThan(amount) {
		return errors.New("insufficient locked funds")
	}

	now := time.Now()
	balance.LockedBalance = balance.LockedBalance.Sub(amount)
	balance.AvailableBalance = balance.AvailableBalance.Add(amount)
	balance.UpdatedAt = now
	w.UpdatedAt = now
	return nil
}

// DebitLocked debits from locked balance (for completed transactions)
func (w *MultiCurrencyWallet) DebitLocked(currency SupportedCurrency, amount decimal.Decimal) error {
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}

	balance, err := w.GetBalance(currency)
	if err != nil {
		return err
	}

	if balance.LockedBalance.LessThan(amount) {
		return errors.New("insufficient locked funds")
	}

	now := time.Now()
	balance.LockedBalance = balance.LockedBalance.Sub(amount)
	balance.TotalBalance = balance.AvailableBalance.Add(balance.PendingBalance).Add(balance.LockedBalance)
	balance.UpdatedAt = now
	w.UpdatedAt = now
	return nil
}

// GetEnabledCurrencies returns all enabled currencies
func (w *MultiCurrencyWallet) GetEnabledCurrencies() []SupportedCurrency {
	currencies := make([]SupportedCurrency, 0, len(w.Balances))
	for currency := range w.Balances {
		currencies = append(currencies, currency)
	}
	return currencies
}

// GetTotalBalanceInPrimary calculates total balance in primary currency
// Requires FX rates to be passed in
func (w *MultiCurrencyWallet) GetTotalBalanceInPrimary(rates map[SupportedCurrency]decimal.Decimal) decimal.Decimal {
	total := decimal.Zero

	for currency, balance := range w.Balances {
		if currency == w.PrimaryCurrency {
			total = total.Add(balance.TotalBalance)
		} else if rate, ok := rates[currency]; ok {
			// Convert to primary currency
			converted := balance.TotalBalance.Mul(rate)
			total = total.Add(converted)
		}
	}

	return total
}

// AllSupportedCurrencies returns all supported currencies
func AllSupportedCurrencies() []SupportedCurrency {
	return []SupportedCurrency{
		CurrencyNGN,
		CurrencyGBP,
		CurrencyUSD,
		CurrencyEUR,
		CurrencyCAD,
		CurrencyGHS,
		CurrencyKES,
	}
}

// GetCurrencySymbol returns the symbol for a currency
func GetCurrencySymbol(currency SupportedCurrency) string {
	symbols := map[SupportedCurrency]string{
		CurrencyNGN: "₦",
		CurrencyGBP: "£",
		CurrencyUSD: "$",
		CurrencyEUR: "€",
		CurrencyCAD: "C$",
		CurrencyGHS: "₵",
		CurrencyKES: "KSh",
	}
	if symbol, ok := symbols[currency]; ok {
		return symbol
	}
	return string(currency)
}

// GetCurrencyName returns the full name for a currency
func GetCurrencyName(currency SupportedCurrency) string {
	names := map[SupportedCurrency]string{
		CurrencyNGN: "Nigerian Naira",
		CurrencyGBP: "British Pound",
		CurrencyUSD: "US Dollar",
		CurrencyEUR: "Euro",
		CurrencyCAD: "Canadian Dollar",
		CurrencyGHS: "Ghanaian Cedi",
		CurrencyKES: "Kenyan Shilling",
	}
	if name, ok := names[currency]; ok {
		return name
	}
	return string(currency)
}
