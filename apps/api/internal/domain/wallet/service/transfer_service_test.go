package service

import (
	"context"
	"errors"
	"testing"

	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
	"hustlex/internal/domain/wallet/repository"
)

// mockWalletRepository is a test double for WalletRepository
type mockWalletRepository struct {
	wallets      map[string]*aggregate.Wallet
	findErr      error
	saveErr      error
	saveCallCount int
}

func newMockWalletRepo() *mockWalletRepository {
	return &mockWalletRepository{
		wallets: make(map[string]*aggregate.Wallet),
	}
}

func (m *mockWalletRepository) FindByID(ctx context.Context, id valueobject.WalletID) (*aggregate.Wallet, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	for _, w := range m.wallets {
		if w.ID().Equals(id) {
			return w, nil
		}
	}
	return nil, repository.ErrWalletNotFound
}

func (m *mockWalletRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.Wallet, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	w, ok := m.wallets[userID.String()]
	if !ok {
		return nil, repository.ErrWalletNotFound
	}
	return w, nil
}

func (m *mockWalletRepository) Exists(ctx context.Context, userID valueobject.UserID) (bool, error) {
	_, ok := m.wallets[userID.String()]
	return ok, nil
}

func (m *mockWalletRepository) Save(ctx context.Context, wallet *aggregate.Wallet) error {
	m.saveCallCount++
	if m.saveErr != nil {
		return m.saveErr
	}
	m.wallets[wallet.UserID().String()] = wallet
	return nil
}

func (m *mockWalletRepository) SaveWithEvents(ctx context.Context, wallet *aggregate.Wallet) error {
	return m.Save(ctx, wallet)
}

func (m *mockWalletRepository) addWallet(w *aggregate.Wallet) {
	m.wallets[w.UserID().String()] = w
}

// Helper to create a funded wallet
func createFundedWallet(userID valueobject.UserID, balance int64) *aggregate.Wallet {
	wallet := aggregate.NewWallet(userID, valueobject.NGN)
	if balance > 0 {
		credit := valueobject.MustNewMoney(balance, valueobject.NGN)
		wallet.Credit(credit, "test", "REF", "Initial")
	}
	wallet.ClearEvents()
	return wallet
}

// Helper to create a wallet with escrow
func createWalletWithEscrow(userID valueobject.UserID, available, escrow int64) *aggregate.Wallet {
	wallet := aggregate.NewWallet(userID, valueobject.NGN)
	total := available + escrow
	if total > 0 {
		credit := valueobject.MustNewMoney(total, valueobject.NGN)
		wallet.Credit(credit, "test", "REF", "Initial")
	}
	if escrow > 0 {
		escrowAmount := valueobject.MustNewMoney(escrow, valueobject.NGN)
		wallet.HoldInEscrow(escrowAmount, "REF", "Test escrow")
	}
	wallet.ClearEvents()
	return wallet
}

func TestTransferService_Transfer(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	senderID := valueobject.GenerateUserID()
	recipientID := valueobject.GenerateUserID()

	senderWallet := createFundedWallet(senderID, 10000) // ₦100
	recipientWallet := createFundedWallet(recipientID, 5000) // ₦50

	repo.addWallet(senderWallet)
	repo.addWallet(recipientWallet)

	service := NewTransferService(repo)

	result, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  senderID,
		ToUserID:    recipientID,
		Amount:      valueobject.MustNewMoney(3000, valueobject.NGN), // ₦30
		Description: "Test transfer",
		Reference:   "TRF123",
	})

	if err != nil {
		t.Fatalf("Transfer() unexpected error: %v", err)
	}

	// Sender: 10000 - 3000 (amount) - 1000 (fee) = 6000
	if result.SenderNewBalance.Amount() != 6000 {
		t.Errorf("Transfer() sender balance = %d, want 6000", result.SenderNewBalance.Amount())
	}

	// Recipient: 5000 + 3000 = 8000
	if result.RecipientNewBalance.Amount() != 8000 {
		t.Errorf("Transfer() recipient balance = %d, want 8000", result.RecipientNewBalance.Amount())
	}

	// Fee should be ₦10 (1000 kobo)
	if result.Fee.Amount() != 1000 {
		t.Errorf("Transfer() fee = %d, want 1000", result.Fee.Amount())
	}

	if result.Reference != "TRF123" {
		t.Errorf("Transfer() reference = %s, want TRF123", result.Reference)
	}

	// Verify repo save was called
	if repo.saveCallCount != 2 {
		t.Errorf("Transfer() save count = %d, want 2", repo.saveCallCount)
	}
}

func TestTransferService_Transfer_SameWallet(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	userID := valueobject.GenerateUserID()
	wallet := createFundedWallet(userID, 10000)
	repo.addWallet(wallet)

	service := NewTransferService(repo)

	_, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  userID,
		ToUserID:    userID, // Same user
		Amount:      valueobject.MustNewMoney(3000, valueobject.NGN),
		Description: "Test",
		Reference:   "REF",
	})

	if err != ErrSameWallet {
		t.Errorf("Transfer() to self error = %v, want ErrSameWallet", err)
	}
}

func TestTransferService_Transfer_SenderNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	service := NewTransferService(repo)

	_, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  valueobject.GenerateUserID(), // Not in repo
		ToUserID:    valueobject.GenerateUserID(),
		Amount:      valueobject.MustNewMoney(3000, valueobject.NGN),
		Description: "Test",
		Reference:   "REF",
	})

	if err != repository.ErrWalletNotFound {
		t.Errorf("Transfer() sender not found error = %v, want ErrWalletNotFound", err)
	}
}

func TestTransferService_Transfer_RecipientNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	senderID := valueobject.GenerateUserID()
	wallet := createFundedWallet(senderID, 10000)
	repo.addWallet(wallet)

	service := NewTransferService(repo)

	_, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  senderID,
		ToUserID:    valueobject.GenerateUserID(), // Not in repo
		Amount:      valueobject.MustNewMoney(3000, valueobject.NGN),
		Description: "Test",
		Reference:   "REF",
	})

	if err != ErrRecipientNotFound {
		t.Errorf("Transfer() recipient not found error = %v, want ErrRecipientNotFound", err)
	}
}

func TestTransferService_Transfer_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	senderID := valueobject.GenerateUserID()
	recipientID := valueobject.GenerateUserID()

	senderWallet := createFundedWallet(senderID, 2000) // Only ₦20
	recipientWallet := createFundedWallet(recipientID, 5000)

	repo.addWallet(senderWallet)
	repo.addWallet(recipientWallet)

	service := NewTransferService(repo)

	_, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  senderID,
		ToUserID:    recipientID,
		Amount:      valueobject.MustNewMoney(5000, valueobject.NGN), // ₦50 > ₦20
		Description: "Test",
		Reference:   "REF",
	})

	if err != aggregate.ErrInsufficientFunds {
		t.Errorf("Transfer() insufficient funds error = %v, want ErrInsufficientFunds", err)
	}
}

func TestTransferService_Transfer_SaveError(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	senderID := valueobject.GenerateUserID()
	recipientID := valueobject.GenerateUserID()

	senderWallet := createFundedWallet(senderID, 10000)
	recipientWallet := createFundedWallet(recipientID, 5000)

	repo.addWallet(senderWallet)
	repo.addWallet(recipientWallet)
	repo.saveErr = errors.New("database error")

	service := NewTransferService(repo)

	_, err := service.Transfer(ctx, TransferRequest{
		FromUserID:  senderID,
		ToUserID:    recipientID,
		Amount:      valueobject.MustNewMoney(3000, valueobject.NGN),
		Description: "Test",
		Reference:   "REF",
	})

	if err == nil {
		t.Error("Transfer() should return error on save failure")
	}
}

// EscrowService Tests

func TestEscrowService_ReleaseEscrowToRecipient(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	payerID := valueobject.GenerateUserID()
	recipientID := valueobject.GenerateUserID()

	// Payer has 5000 available and 10000 in escrow
	payerWallet := createWalletWithEscrow(payerID, 5000, 10000)
	recipientWallet := createFundedWallet(recipientID, 2000)

	repo.addWallet(payerWallet)
	repo.addWallet(recipientWallet)

	service := NewEscrowService(repo)

	err := service.ReleaseEscrowToRecipient(ctx, EscrowReleaseRequest{
		PayerUserID:     payerID,
		RecipientUserID: recipientID,
		Amount:          valueobject.MustNewMoney(10000, valueobject.NGN), // Full escrow
		PlatformFee:     valueobject.MustNewMoney(1000, valueobject.NGN), // 10% fee
		Reference:       "ESC123",
		Description:     "Gig payment",
	})

	if err != nil {
		t.Fatalf("ReleaseEscrowToRecipient() unexpected error: %v", err)
	}

	// Payer escrow should be 0
	if payerWallet.EscrowBalance().Amount() != 0 {
		t.Errorf("Payer escrow = %d, want 0", payerWallet.EscrowBalance().Amount())
	}

	// Recipient gets 10000 - 1000 (fee) = 9000, total = 2000 + 9000 = 11000
	if recipientWallet.AvailableBalance().Amount() != 11000 {
		t.Errorf("Recipient balance = %d, want 11000", recipientWallet.AvailableBalance().Amount())
	}
}

func TestEscrowService_ReleaseEscrowToRecipient_InsufficientEscrow(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	payerID := valueobject.GenerateUserID()
	recipientID := valueobject.GenerateUserID()

	payerWallet := createWalletWithEscrow(payerID, 5000, 5000) // Only 5000 in escrow
	recipientWallet := createFundedWallet(recipientID, 2000)

	repo.addWallet(payerWallet)
	repo.addWallet(recipientWallet)

	service := NewEscrowService(repo)

	err := service.ReleaseEscrowToRecipient(ctx, EscrowReleaseRequest{
		PayerUserID:     payerID,
		RecipientUserID: recipientID,
		Amount:          valueobject.MustNewMoney(10000, valueobject.NGN), // More than escrow
		PlatformFee:     valueobject.MustNewMoney(1000, valueobject.NGN),
		Reference:       "ESC123",
		Description:     "Gig payment",
	})

	if err != aggregate.ErrInsufficientEscrow {
		t.Errorf("ReleaseEscrowToRecipient() error = %v, want ErrInsufficientEscrow", err)
	}
}

func TestEscrowService_RefundEscrow(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	userID := valueobject.GenerateUserID()
	wallet := createWalletWithEscrow(userID, 5000, 10000)
	repo.addWallet(wallet)

	service := NewEscrowService(repo)

	err := service.RefundEscrow(ctx, userID, valueobject.MustNewMoney(10000, valueobject.NGN), "REF123", "Contract cancelled")

	if err != nil {
		t.Fatalf("RefundEscrow() unexpected error: %v", err)
	}

	// Escrow should be 0
	if wallet.EscrowBalance().Amount() != 0 {
		t.Errorf("RefundEscrow() escrow = %d, want 0", wallet.EscrowBalance().Amount())
	}

	// Available should be 5000 + 10000 = 15000
	if wallet.AvailableBalance().Amount() != 15000 {
		t.Errorf("RefundEscrow() available = %d, want 15000", wallet.AvailableBalance().Amount())
	}
}

func TestEscrowService_RefundEscrow_InsufficientEscrow(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	userID := valueobject.GenerateUserID()
	wallet := createWalletWithEscrow(userID, 5000, 3000)
	repo.addWallet(wallet)

	service := NewEscrowService(repo)

	err := service.RefundEscrow(ctx, userID, valueobject.MustNewMoney(5000, valueobject.NGN), "REF", "reason")

	if err != aggregate.ErrInsufficientEscrow {
		t.Errorf("RefundEscrow() error = %v, want ErrInsufficientEscrow", err)
	}
}

func TestEscrowService_RefundEscrow_WalletNotFound(t *testing.T) {
	ctx := context.Background()
	repo := newMockWalletRepo()

	service := NewEscrowService(repo)

	err := service.RefundEscrow(ctx, valueobject.GenerateUserID(), valueobject.MustNewMoney(1000, valueobject.NGN), "REF", "reason")

	if err != repository.ErrWalletNotFound {
		t.Errorf("RefundEscrow() wallet not found error = %v, want ErrWalletNotFound", err)
	}
}
