package valueobject

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount    = errors.New("invalid amount: must be non-negative")
	ErrCurrencyMismatch = errors.New("currency mismatch")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Currency represents a supported currency
type Currency string

const (
	NGN Currency = "NGN" // Nigerian Naira
	USD Currency = "USD" // US Dollar
	GHS Currency = "GHS" // Ghanaian Cedi
)

// IsValid checks if the currency is supported
func (c Currency) IsValid() bool {
	switch c {
	case NGN, USD, GHS:
		return true
	default:
		return false
	}
}

// SmallestUnit returns the name of the smallest unit
func (c Currency) SmallestUnit() string {
	switch c {
	case NGN:
		return "kobo"
	case USD:
		return "cents"
	case GHS:
		return "pesewas"
	default:
		return "units"
	}
}

// Money is an immutable value object representing monetary value
// Amount is stored in the smallest unit (kobo for NGN, cents for USD)
type Money struct {
	amount   int64
	currency Currency
}

// NewMoney creates a new Money value object
func NewMoney(amount int64, currency Currency) (Money, error) {
	if amount < 0 {
		return Money{}, ErrInvalidAmount
	}
	if !currency.IsValid() {
		return Money{}, fmt.Errorf("unsupported currency: %s", currency)
	}
	return Money{amount: amount, currency: currency}, nil
}

// MustNewMoney creates a Money or panics - use only for known-valid values
func MustNewMoney(amount int64, currency Currency) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// Zero returns a zero-value Money for the given currency
func Zero(currency Currency) Money {
	return Money{amount: 0, currency: currency}
}

// Amount returns the amount in smallest unit
func (m Money) Amount() int64 {
	return m.amount
}

// Currency returns the currency
func (m Money) Currency() Currency {
	return m.currency
}

// Add returns a new Money with the sum of both amounts
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return NewMoney(m.amount+other.amount, m.currency)
}

// MustAdd adds or panics - use only when currency match is guaranteed
func (m Money) MustAdd(other Money) Money {
	result, err := m.Add(other)
	if err != nil {
		panic(err)
	}
	return result
}

// Subtract returns a new Money with the difference
func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	if m.amount < other.amount {
		return Money{}, ErrInsufficientFunds
	}
	return NewMoney(m.amount-other.amount, m.currency)
}

// MustSubtract subtracts or panics
func (m Money) MustSubtract(other Money) Money {
	result, err := m.Subtract(other)
	if err != nil {
		panic(err)
	}
	return result
}

// Multiply returns a new Money multiplied by a factor
func (m Money) Multiply(factor float64) (Money, error) {
	if factor < 0 {
		return Money{}, errors.New("factor must be non-negative")
	}
	newAmount := int64(float64(m.amount) * factor)
	return NewMoney(newAmount, m.currency)
}

// Percentage returns a percentage of the money
func (m Money) Percentage(percent float64) (Money, error) {
	return m.Multiply(percent / 100)
}

// IsZero returns true if the amount is zero
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsPositive returns true if the amount is greater than zero
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// GreaterThan returns true if this amount is greater than other
func (m Money) GreaterThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}
	return m.amount > other.amount
}

// GreaterThanOrEqual returns true if this amount is >= other
func (m Money) GreaterThanOrEqual(other Money) bool {
	if m.currency != other.currency {
		return false
	}
	return m.amount >= other.amount
}

// LessThan returns true if this amount is less than other
func (m Money) LessThan(other Money) bool {
	if m.currency != other.currency {
		return false
	}
	return m.amount < other.amount
}

// Equals returns true if both amount and currency match
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns a human-readable representation
func (m Money) String() string {
	major := m.amount / 100
	minor := m.amount % 100

	var symbol string
	switch m.currency {
	case NGN:
		symbol = "₦"
	case USD:
		symbol = "$"
	case GHS:
		symbol = "GH₵"
	default:
		symbol = string(m.currency) + " "
	}

	return fmt.Sprintf("%s%d.%02d", symbol, major, minor)
}

// Format returns a formatted string with thousands separator
func (m Money) Format() string {
	major := m.amount / 100
	minor := m.amount % 100

	var symbol string
	switch m.currency {
	case NGN:
		symbol = "₦"
	case USD:
		symbol = "$"
	case GHS:
		symbol = "GH₵"
	default:
		symbol = string(m.currency) + " "
	}

	// Add thousands separator
	majorStr := formatWithCommas(major)
	return fmt.Sprintf("%s%s.%02d", symbol, majorStr, minor)
}

func formatWithCommas(n int64) string {
	if n < 0 {
		return "-" + formatWithCommas(-n)
	}
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}
	return formatWithCommas(n/1000) + "," + fmt.Sprintf("%03d", n%1000)
}
