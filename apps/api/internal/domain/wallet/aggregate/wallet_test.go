package aggregate

import (
	"testing"
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

func TestNewWallet(t *testing.T) {
	userID := valueobject.GenerateUserID()

	wallet := NewWallet(userID, valueobject.NGN)

	if wallet.ID().IsEmpty() {
		t.Error("NewWallet() should generate a wallet ID")
	}
	if !wallet.UserID().Equals(userID) {
		t.Error("NewWallet() should set the user ID")
	}
	if wallet.Currency() != valueobject.NGN {
		t.Errorf("NewWallet() currency = %s, want NGN", wallet.Currency())
	}
	if wallet.Status() != WalletStatusActive {
		t.Errorf("NewWallet() status = %s, want active", wallet.Status())
	}
	if !wallet.AvailableBalance().IsZero() {
		t.Error("NewWallet() available balance should be zero")
	}
	if !wallet.EscrowBalance().IsZero() {
		t.Error("NewWallet() escrow balance should be zero")
	}
	if !wallet.SavingsBalance().IsZero() {
		t.Error("NewWallet() savings balance should be zero")
	}
	if wallet.HasPIN() {
		t.Error("NewWallet() should not have PIN set")
	}

	// Should record WalletCreated event
	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("NewWallet() should record 1 event, got %d", len(events))
	}
}

func TestWallet_Credit(t *testing.T) {
	userID := valueobject.GenerateUserID()
	wallet := NewWallet(userID, valueobject.NGN)
	wallet.ClearEvents() // Clear creation event

	amount := valueobject.MustNewMoney(10000, valueobject.NGN) // ₦100

	err := wallet.Credit(amount, "deposit", "REF123", "Test deposit")
	if err != nil {
		t.Fatalf("Credit() unexpected error: %v", err)
	}

	if wallet.AvailableBalance().Amount() != 10000 {
		t.Errorf("Credit() available = %d, want 10000", wallet.AvailableBalance().Amount())
	}
	if wallet.LedgerBalance().Amount() != 10000 {
		t.Errorf("Credit() ledger = %d, want 10000", wallet.LedgerBalance().Amount())
	}

	// Should record WalletCredited event
	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("Credit() should record 1 event, got %d", len(events))
	}
}

func TestWallet_Credit_Errors(t *testing.T) {
	t.Run("zero amount", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		amount := valueobject.MustNewMoney(0, valueobject.NGN)

		err := wallet.Credit(amount, "deposit", "REF", "test")
		if err != ErrInvalidAmount {
			t.Errorf("Credit(0) error = %v, want ErrInvalidAmount", err)
		}
	})

	t.Run("currency mismatch", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		amount := valueobject.MustNewMoney(1000, valueobject.USD)

		err := wallet.Credit(amount, "deposit", "REF", "test")
		if err != ErrCurrencyMismatch {
			t.Errorf("Credit(USD) error = %v, want ErrCurrencyMismatch", err)
		}
	})

	t.Run("locked wallet", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		wallet.Lock("test")
		amount := valueobject.MustNewMoney(1000, valueobject.NGN)

		err := wallet.Credit(amount, "deposit", "REF", "test")
		if err != ErrWalletLocked {
			t.Errorf("Credit() on locked wallet error = %v, want ErrWalletLocked", err)
		}
	})
}

func TestWallet_Debit(t *testing.T) {
	userID := valueobject.GenerateUserID()
	wallet := NewWallet(userID, valueobject.NGN)

	// Credit first
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF1", "Initial")
	wallet.ClearEvents()

	amount := valueobject.MustNewMoney(3000, valueobject.NGN)
	fee := valueobject.MustNewMoney(100, valueobject.NGN)

	err := wallet.Debit(amount, "withdrawal", "REF2", "Test withdrawal", fee)
	if err != nil {
		t.Fatalf("Debit() unexpected error: %v", err)
	}

	// 10000 - 3000 - 100 = 6900
	if wallet.AvailableBalance().Amount() != 6900 {
		t.Errorf("Debit() available = %d, want 6900", wallet.AvailableBalance().Amount())
	}

	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("Debit() should record 1 event, got %d", len(events))
	}
}

func TestWallet_Debit_Errors(t *testing.T) {
	t.Run("insufficient funds", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		credit := valueobject.MustNewMoney(5000, valueobject.NGN)
		wallet.Credit(credit, "deposit", "REF", "Initial")

		amount := valueobject.MustNewMoney(10000, valueobject.NGN)
		fee := valueobject.Zero(valueobject.NGN)

		err := wallet.Debit(amount, "withdrawal", "REF", "test", fee)
		if err != ErrInsufficientFunds {
			t.Errorf("Debit() error = %v, want ErrInsufficientFunds", err)
		}
	})

	t.Run("insufficient funds with fee", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		credit := valueobject.MustNewMoney(5000, valueobject.NGN)
		wallet.Credit(credit, "deposit", "REF", "Initial")

		amount := valueobject.MustNewMoney(4900, valueobject.NGN)
		fee := valueobject.MustNewMoney(200, valueobject.NGN) // 4900 + 200 > 5000

		err := wallet.Debit(amount, "withdrawal", "REF", "test", fee)
		if err != ErrInsufficientFunds {
			t.Errorf("Debit() with fee error = %v, want ErrInsufficientFunds", err)
		}
	})

	t.Run("zero amount", func(t *testing.T) {
		wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
		amount := valueobject.MustNewMoney(0, valueobject.NGN)
		fee := valueobject.Zero(valueobject.NGN)

		err := wallet.Debit(amount, "withdrawal", "REF", "test", fee)
		if err != ErrInvalidAmount {
			t.Errorf("Debit(0) error = %v, want ErrInvalidAmount", err)
		}
	})
}

func TestWallet_HoldInEscrow(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")
	wallet.ClearEvents()

	escrowAmount := valueobject.MustNewMoney(4000, valueobject.NGN)

	err := wallet.HoldInEscrow(escrowAmount, "CONTRACT123", "Contract escrow")
	if err != nil {
		t.Fatalf("HoldInEscrow() unexpected error: %v", err)
	}

	if wallet.AvailableBalance().Amount() != 6000 {
		t.Errorf("HoldInEscrow() available = %d, want 6000", wallet.AvailableBalance().Amount())
	}
	if wallet.EscrowBalance().Amount() != 4000 {
		t.Errorf("HoldInEscrow() escrow = %d, want 4000", wallet.EscrowBalance().Amount())
	}
	// Ledger should remain the same
	if wallet.LedgerBalance().Amount() != 10000 {
		t.Errorf("HoldInEscrow() ledger = %d, want 10000", wallet.LedgerBalance().Amount())
	}

	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("HoldInEscrow() should record 1 event, got %d", len(events))
	}
}

func TestWallet_HoldInEscrow_InsufficientFunds(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(5000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	escrowAmount := valueobject.MustNewMoney(10000, valueobject.NGN)

	err := wallet.HoldInEscrow(escrowAmount, "CONTRACT", "test")
	if err != ErrInsufficientFunds {
		t.Errorf("HoldInEscrow() error = %v, want ErrInsufficientFunds", err)
	}
}

func TestWallet_ReleaseFromEscrow_ToWallet(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	escrowAmount := valueobject.MustNewMoney(4000, valueobject.NGN)
	wallet.HoldInEscrow(escrowAmount, "CONTRACT", "escrow")
	wallet.ClearEvents()

	releaseAmount := valueobject.MustNewMoney(4000, valueobject.NGN)

	err := wallet.ReleaseFromEscrow(releaseAmount, "CONTRACT", true, "")
	if err != nil {
		t.Fatalf("ReleaseFromEscrow() unexpected error: %v", err)
	}

	// Funds return to available
	if wallet.AvailableBalance().Amount() != 10000 {
		t.Errorf("ReleaseFromEscrow(toWallet) available = %d, want 10000", wallet.AvailableBalance().Amount())
	}
	if wallet.EscrowBalance().Amount() != 0 {
		t.Errorf("ReleaseFromEscrow(toWallet) escrow = %d, want 0", wallet.EscrowBalance().Amount())
	}
	// Ledger unchanged
	if wallet.LedgerBalance().Amount() != 10000 {
		t.Errorf("ReleaseFromEscrow(toWallet) ledger = %d, want 10000", wallet.LedgerBalance().Amount())
	}
}

func TestWallet_ReleaseFromEscrow_ToExternal(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	escrowAmount := valueobject.MustNewMoney(4000, valueobject.NGN)
	wallet.HoldInEscrow(escrowAmount, "CONTRACT", "escrow")
	wallet.ClearEvents()

	releaseAmount := valueobject.MustNewMoney(4000, valueobject.NGN)
	recipientID := valueobject.GenerateUserID().String()

	err := wallet.ReleaseFromEscrow(releaseAmount, "CONTRACT", false, recipientID)
	if err != nil {
		t.Fatalf("ReleaseFromEscrow() unexpected error: %v", err)
	}

	// Available unchanged
	if wallet.AvailableBalance().Amount() != 6000 {
		t.Errorf("ReleaseFromEscrow(external) available = %d, want 6000", wallet.AvailableBalance().Amount())
	}
	if wallet.EscrowBalance().Amount() != 0 {
		t.Errorf("ReleaseFromEscrow(external) escrow = %d, want 0", wallet.EscrowBalance().Amount())
	}
	// Ledger reduced by released amount
	if wallet.LedgerBalance().Amount() != 6000 {
		t.Errorf("ReleaseFromEscrow(external) ledger = %d, want 6000", wallet.LedgerBalance().Amount())
	}
}

func TestWallet_ReleaseFromEscrow_InsufficientEscrow(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	escrowAmount := valueobject.MustNewMoney(4000, valueobject.NGN)
	wallet.HoldInEscrow(escrowAmount, "CONTRACT", "escrow")

	releaseAmount := valueobject.MustNewMoney(5000, valueobject.NGN)

	err := wallet.ReleaseFromEscrow(releaseAmount, "CONTRACT", true, "")
	if err != ErrInsufficientEscrow {
		t.Errorf("ReleaseFromEscrow() error = %v, want ErrInsufficientEscrow", err)
	}
}

func TestWallet_MoveToSavings(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	savingsAmount := valueobject.MustNewMoney(3000, valueobject.NGN)

	err := wallet.MoveToSavings(savingsAmount)
	if err != nil {
		t.Fatalf("MoveToSavings() unexpected error: %v", err)
	}

	if wallet.AvailableBalance().Amount() != 7000 {
		t.Errorf("MoveToSavings() available = %d, want 7000", wallet.AvailableBalance().Amount())
	}
	if wallet.SavingsBalance().Amount() != 3000 {
		t.Errorf("MoveToSavings() savings = %d, want 3000", wallet.SavingsBalance().Amount())
	}
}

func TestWallet_WithdrawFromSavings(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	savingsAmount := valueobject.MustNewMoney(5000, valueobject.NGN)
	wallet.MoveToSavings(savingsAmount)

	withdrawAmount := valueobject.MustNewMoney(2000, valueobject.NGN)

	err := wallet.WithdrawFromSavings(withdrawAmount)
	if err != nil {
		t.Fatalf("WithdrawFromSavings() unexpected error: %v", err)
	}

	if wallet.AvailableBalance().Amount() != 7000 {
		t.Errorf("WithdrawFromSavings() available = %d, want 7000", wallet.AvailableBalance().Amount())
	}
	if wallet.SavingsBalance().Amount() != 3000 {
		t.Errorf("WithdrawFromSavings() savings = %d, want 3000", wallet.SavingsBalance().Amount())
	}
}

func TestWallet_WithdrawFromSavings_Insufficient(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	savingsAmount := valueobject.MustNewMoney(3000, valueobject.NGN)
	wallet.MoveToSavings(savingsAmount)

	withdrawAmount := valueobject.MustNewMoney(5000, valueobject.NGN)

	err := wallet.WithdrawFromSavings(withdrawAmount)
	if err == nil {
		t.Error("WithdrawFromSavings() should return error for insufficient savings")
	}
}

func TestWallet_Lock_Unlock(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	wallet.ClearEvents()

	wallet.Lock("Security concern")

	if wallet.Status() != WalletStatusLocked {
		t.Errorf("Lock() status = %s, want locked", wallet.Status())
	}
	if !wallet.IsLocked() {
		t.Error("Lock() IsLocked should be true")
	}
	if wallet.IsActive() {
		t.Error("Lock() IsActive should be false")
	}

	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("Lock() should record 1 event, got %d", len(events))
	}

	wallet.ClearEvents()
	wallet.Unlock()

	if wallet.Status() != WalletStatusActive {
		t.Errorf("Unlock() status = %s, want active", wallet.Status())
	}
	if wallet.IsLocked() {
		t.Error("Unlock() IsLocked should be false")
	}
	if !wallet.IsActive() {
		t.Error("Unlock() IsActive should be true")
	}

	events = wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("Unlock() should record 1 event, got %d", len(events))
	}
}

func TestWallet_SetPIN(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)

	if wallet.HasPIN() {
		t.Error("new wallet should not have PIN")
	}

	wallet.ClearEvents()
	wallet.SetPIN("hashed_pin_123")

	if !wallet.HasPIN() {
		t.Error("SetPIN() HasPIN should be true")
	}
	if wallet.PINHash() != "hashed_pin_123" {
		t.Errorf("SetPIN() PINHash = %s, want hashed_pin_123", wallet.PINHash())
	}

	events := wallet.DomainEvents()
	if len(events) != 1 {
		t.Fatalf("SetPIN() should record 1 event, got %d", len(events))
	}
}

func TestWallet_RecordFailedPINAttempt(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	wallet.SetPIN("hash")

	// Should not lock before max attempts
	wallet.RecordFailedPINAttempt(3)
	if wallet.PINAttempts() != 1 {
		t.Errorf("PINAttempts() = %d, want 1", wallet.PINAttempts())
	}
	if wallet.IsLocked() {
		t.Error("should not be locked after 1 attempt")
	}

	wallet.RecordFailedPINAttempt(3)
	if wallet.PINAttempts() != 2 {
		t.Errorf("PINAttempts() = %d, want 2", wallet.PINAttempts())
	}
	if wallet.IsLocked() {
		t.Error("should not be locked after 2 attempts")
	}

	// Third attempt should lock
	wallet.RecordFailedPINAttempt(3)
	if !wallet.IsLocked() {
		t.Error("should be locked after 3 attempts")
	}
}

func TestWallet_ResetPINAttempts(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	wallet.SetPIN("hash")
	wallet.RecordFailedPINAttempt(5)
	wallet.RecordFailedPINAttempt(5)

	if wallet.PINAttempts() != 2 {
		t.Errorf("PINAttempts() = %d, want 2", wallet.PINAttempts())
	}

	wallet.ResetPINAttempts()

	if wallet.PINAttempts() != 0 {
		t.Errorf("ResetPINAttempts() = %d, want 0", wallet.PINAttempts())
	}
}

func TestWallet_TotalBalance(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	credit := valueobject.MustNewMoney(10000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "Initial")

	escrowAmount := valueobject.MustNewMoney(3000, valueobject.NGN)
	wallet.HoldInEscrow(escrowAmount, "CONTRACT", "escrow")

	savingsAmount := valueobject.MustNewMoney(2000, valueobject.NGN)
	wallet.MoveToSavings(savingsAmount)

	// available: 5000, escrow: 3000, savings: 2000 = total 10000
	total := wallet.TotalBalance()
	if total.Amount() != 10000 {
		t.Errorf("TotalBalance() = %d, want 10000", total.Amount())
	}
}

func TestWallet_Version(t *testing.T) {
	wallet := NewWallet(valueobject.GenerateUserID(), valueobject.NGN)
	initialVersion := wallet.Version()

	credit := valueobject.MustNewMoney(1000, valueobject.NGN)
	wallet.Credit(credit, "deposit", "REF", "test")

	if wallet.Version() <= initialVersion {
		t.Error("Version should increment after Credit")
	}
}

func TestReconstitute(t *testing.T) {
	walletID := valueobject.GenerateWalletID()
	userID := valueobject.GenerateUserID()
	available := valueobject.MustNewMoney(5000, valueobject.NGN)
	escrow := valueobject.MustNewMoney(2000, valueobject.NGN)
	savings := valueobject.MustNewMoney(1000, valueobject.NGN)
	ledger := valueobject.MustNewMoney(8000, valueobject.NGN)

	wallet := Reconstitute(
		walletID,
		userID,
		available,
		escrow,
		savings,
		ledger,
		valueobject.NGN,
		WalletStatusActive,
		"pin_hash",
		0,
		timeNow(),
		timeNow(),
		5,
	)

	if !wallet.ID().Equals(walletID) {
		t.Error("Reconstitute() should set ID")
	}
	if !wallet.UserID().Equals(userID) {
		t.Error("Reconstitute() should set UserID")
	}
	if wallet.AvailableBalance().Amount() != 5000 {
		t.Errorf("Reconstitute() available = %d, want 5000", wallet.AvailableBalance().Amount())
	}
	if wallet.EscrowBalance().Amount() != 2000 {
		t.Errorf("Reconstitute() escrow = %d, want 2000", wallet.EscrowBalance().Amount())
	}
	if wallet.SavingsBalance().Amount() != 1000 {
		t.Errorf("Reconstitute() savings = %d, want 1000", wallet.SavingsBalance().Amount())
	}
	if wallet.Version() != 5 {
		t.Errorf("Reconstitute() version = %d, want 5", wallet.Version())
	}

	// Reconstitute should not record any events
	if len(wallet.DomainEvents()) != 0 {
		t.Error("Reconstitute() should not record events")
	}
}

// Helper
func timeNow() time.Time {
	return time.Now().UTC()
}
