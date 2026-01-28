package valueobject

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name      string
		amount    int64
		currency  Currency
		wantErr   bool
		errString string
	}{
		{
			name:     "valid NGN money",
			amount:   10000,
			currency: NGN,
			wantErr:  false,
		},
		{
			name:     "valid USD money",
			amount:   5000,
			currency: USD,
			wantErr:  false,
		},
		{
			name:     "zero amount is valid",
			amount:   0,
			currency: NGN,
			wantErr:  false,
		},
		{
			name:      "negative amount is invalid",
			amount:    -100,
			currency:  NGN,
			wantErr:   true,
			errString: "invalid amount",
		},
		{
			name:      "unsupported currency",
			amount:    100,
			currency:  Currency("EUR"),
			wantErr:   true,
			errString: "unsupported currency",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := NewMoney(tt.amount, tt.currency)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewMoney() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewMoney() unexpected error: %v", err)
				return
			}
			if money.Amount() != tt.amount {
				t.Errorf("Amount() = %d, want %d", money.Amount(), tt.amount)
			}
			if money.Currency() != tt.currency {
				t.Errorf("Currency() = %s, want %s", money.Currency(), tt.currency)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name       string
		money1     Money
		money2     Money
		wantAmount int64
		wantErr    bool
	}{
		{
			name:       "add same currency",
			money1:     MustNewMoney(10000, NGN),
			money2:     MustNewMoney(5000, NGN),
			wantAmount: 15000,
			wantErr:    false,
		},
		{
			name:       "add zero",
			money1:     MustNewMoney(10000, NGN),
			money2:     MustNewMoney(0, NGN),
			wantAmount: 10000,
			wantErr:    false,
		},
		{
			name:    "currency mismatch",
			money1:  MustNewMoney(10000, NGN),
			money2:  MustNewMoney(5000, USD),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Add(tt.money2)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Add() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Add() unexpected error: %v", err)
				return
			}
			if result.Amount() != tt.wantAmount {
				t.Errorf("Add() = %d, want %d", result.Amount(), tt.wantAmount)
			}
		})
	}
}

func TestMoney_Subtract(t *testing.T) {
	tests := []struct {
		name       string
		money1     Money
		money2     Money
		wantAmount int64
		wantErr    error
	}{
		{
			name:       "subtract valid",
			money1:     MustNewMoney(10000, NGN),
			money2:     MustNewMoney(3000, NGN),
			wantAmount: 7000,
			wantErr:    nil,
		},
		{
			name:       "subtract to zero",
			money1:     MustNewMoney(10000, NGN),
			money2:     MustNewMoney(10000, NGN),
			wantAmount: 0,
			wantErr:    nil,
		},
		{
			name:    "insufficient funds",
			money1:  MustNewMoney(5000, NGN),
			money2:  MustNewMoney(10000, NGN),
			wantErr: ErrInsufficientFunds,
		},
		{
			name:    "currency mismatch",
			money1:  MustNewMoney(10000, NGN),
			money2:  MustNewMoney(5000, USD),
			wantErr: ErrCurrencyMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money1.Subtract(tt.money2)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("Subtract() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Subtract() unexpected error: %v", err)
				return
			}
			if result.Amount() != tt.wantAmount {
				t.Errorf("Subtract() = %d, want %d", result.Amount(), tt.wantAmount)
			}
		})
	}
}

func TestMoney_Multiply(t *testing.T) {
	tests := []struct {
		name       string
		money      Money
		factor     float64
		wantAmount int64
		wantErr    bool
	}{
		{
			name:       "multiply by 2",
			money:      MustNewMoney(10000, NGN),
			factor:     2.0,
			wantAmount: 20000,
			wantErr:    false,
		},
		{
			name:       "multiply by 0.5",
			money:      MustNewMoney(10000, NGN),
			factor:     0.5,
			wantAmount: 5000,
			wantErr:    false,
		},
		{
			name:       "multiply by 0",
			money:      MustNewMoney(10000, NGN),
			factor:     0,
			wantAmount: 0,
			wantErr:    false,
		},
		{
			name:    "negative factor",
			money:   MustNewMoney(10000, NGN),
			factor:  -1.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.money.Multiply(tt.factor)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Multiply() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Multiply() unexpected error: %v", err)
				return
			}
			if result.Amount() != tt.wantAmount {
				t.Errorf("Multiply() = %d, want %d", result.Amount(), tt.wantAmount)
			}
		})
	}
}

func TestMoney_Percentage(t *testing.T) {
	money := MustNewMoney(10000, NGN) // ₦100.00

	result, err := money.Percentage(10) // 10%
	if err != nil {
		t.Fatalf("Percentage() unexpected error: %v", err)
	}

	expected := int64(1000) // ₦10.00
	if result.Amount() != expected {
		t.Errorf("Percentage(10) = %d, want %d", result.Amount(), expected)
	}
}

func TestMoney_Comparisons(t *testing.T) {
	money1 := MustNewMoney(10000, NGN)
	money2 := MustNewMoney(5000, NGN)
	money3 := MustNewMoney(10000, NGN)
	moneyUSD := MustNewMoney(10000, USD)

	t.Run("GreaterThan", func(t *testing.T) {
		if !money1.GreaterThan(money2) {
			t.Error("10000 should be greater than 5000")
		}
		if money2.GreaterThan(money1) {
			t.Error("5000 should not be greater than 10000")
		}
		if money1.GreaterThan(moneyUSD) {
			t.Error("different currencies should return false")
		}
	})

	t.Run("GreaterThanOrEqual", func(t *testing.T) {
		if !money1.GreaterThanOrEqual(money2) {
			t.Error("10000 should be >= 5000")
		}
		if !money1.GreaterThanOrEqual(money3) {
			t.Error("10000 should be >= 10000")
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		if !money2.LessThan(money1) {
			t.Error("5000 should be less than 10000")
		}
		if money1.LessThan(money2) {
			t.Error("10000 should not be less than 5000")
		}
	})

	t.Run("Equals", func(t *testing.T) {
		if !money1.Equals(money3) {
			t.Error("same amount and currency should be equal")
		}
		if money1.Equals(money2) {
			t.Error("different amounts should not be equal")
		}
		if money1.Equals(moneyUSD) {
			t.Error("different currencies should not be equal")
		}
	})
}

func TestMoney_IsZero(t *testing.T) {
	zero := MustNewMoney(0, NGN)
	nonZero := MustNewMoney(100, NGN)

	if !zero.IsZero() {
		t.Error("zero amount should be zero")
	}
	if nonZero.IsZero() {
		t.Error("non-zero amount should not be zero")
	}
}

func TestMoney_IsPositive(t *testing.T) {
	zero := MustNewMoney(0, NGN)
	positive := MustNewMoney(100, NGN)

	if zero.IsPositive() {
		t.Error("zero should not be positive")
	}
	if !positive.IsPositive() {
		t.Error("100 should be positive")
	}
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		money    Money
		expected string
	}{
		{MustNewMoney(10000, NGN), "₦100.00"},
		{MustNewMoney(10050, NGN), "₦100.50"},
		{MustNewMoney(5000, USD), "$50.00"},
		{MustNewMoney(0, NGN), "₦0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.money.String(); got != tt.expected {
				t.Errorf("String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestMoney_Format(t *testing.T) {
	tests := []struct {
		money    Money
		expected string
	}{
		{MustNewMoney(100000000, NGN), "₦1,000,000.00"},
		{MustNewMoney(10050, NGN), "₦100.50"},
		{MustNewMoney(123456789, USD), "$1,234,567.89"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.money.Format(); got != tt.expected {
				t.Errorf("Format() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestCurrency_IsValid(t *testing.T) {
	tests := []struct {
		currency Currency
		valid    bool
	}{
		{NGN, true},
		{USD, true},
		{GHS, true},
		{Currency("EUR"), false},
		{Currency(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.currency), func(t *testing.T) {
			if got := tt.currency.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestCurrency_SmallestUnit(t *testing.T) {
	tests := []struct {
		currency Currency
		unit     string
	}{
		{NGN, "kobo"},
		{USD, "cents"},
		{GHS, "pesewas"},
		{Currency("EUR"), "units"},
	}

	for _, tt := range tests {
		t.Run(string(tt.currency), func(t *testing.T) {
			if got := tt.currency.SmallestUnit(); got != tt.unit {
				t.Errorf("SmallestUnit() = %s, want %s", got, tt.unit)
			}
		})
	}
}

func TestMustNewMoney_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustNewMoney should panic on invalid input")
		}
	}()

	MustNewMoney(-100, NGN)
}

func TestZero(t *testing.T) {
	zero := Zero(NGN)

	if zero.Amount() != 0 {
		t.Errorf("Zero() amount = %d, want 0", zero.Amount())
	}
	if zero.Currency() != NGN {
		t.Errorf("Zero() currency = %s, want NGN", zero.Currency())
	}
}
