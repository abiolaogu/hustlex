package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/domain/wallet/aggregate"
	"hustlex/internal/domain/wallet/repository"
)

// WalletRepository implements repository.WalletRepository using PostgreSQL
type WalletRepository struct {
	db        *sql.DB
	eventBus  EventPublisher
}

// EventPublisher interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, events []interface{}) error
}

// NewWalletRepository creates a new PostgreSQL wallet repository
func NewWalletRepository(db *sql.DB, eventBus EventPublisher) *WalletRepository {
	return &WalletRepository{
		db:       db,
		eventBus: eventBus,
	}
}

// Save persists a wallet
func (r *WalletRepository) Save(ctx context.Context, wallet *aggregate.Wallet) error {
	query := `
		INSERT INTO wallets (
			id, user_id, balance, currency, escrow_balance,
			status, created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			balance = EXCLUDED.balance,
			escrow_balance = EXCLUDED.escrow_balance,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at,
			version = wallets.version + 1
		WHERE wallets.version = $9 - 1
	`

	result, err := r.db.ExecContext(ctx, query,
		wallet.ID().String(),
		wallet.UserID().String(),
		wallet.Balance().Amount(),
		string(wallet.Balance().Currency()),
		wallet.EscrowBalance().Amount(),
		wallet.Status().String(),
		wallet.CreatedAt(),
		wallet.UpdatedAt(),
		wallet.Version(),
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrConcurrentModification
	}

	return nil
}

// SaveWithEvents persists a wallet and publishes domain events
func (r *WalletRepository) SaveWithEvents(ctx context.Context, wallet *aggregate.Wallet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save wallet
	query := `
		INSERT INTO wallets (
			id, user_id, balance, currency, escrow_balance,
			status, created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			balance = EXCLUDED.balance,
			escrow_balance = EXCLUDED.escrow_balance,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at,
			version = wallets.version + 1
		WHERE wallets.version = $9 - 1
	`

	result, err := tx.ExecContext(ctx, query,
		wallet.ID().String(),
		wallet.UserID().String(),
		wallet.Balance().Amount(),
		string(wallet.Balance().Currency()),
		wallet.EscrowBalance().Amount(),
		wallet.Status().String(),
		wallet.CreatedAt(),
		wallet.UpdatedAt(),
		wallet.Version(),
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repository.ErrConcurrentModification
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// Publish events after successful commit
	events := wallet.DomainEvents()
	if len(events) > 0 && r.eventBus != nil {
		eventInterfaces := make([]interface{}, len(events))
		for i, e := range events {
			eventInterfaces[i] = e
		}
		if err := r.eventBus.Publish(ctx, eventInterfaces); err != nil {
			// Log error but don't fail the operation
			// Events can be retried via outbox pattern
		}
		wallet.ClearDomainEvents()
	}

	return nil
}

// FindByID retrieves a wallet by ID
func (r *WalletRepository) FindByID(ctx context.Context, id valueobject.WalletID) (*aggregate.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, escrow_balance,
			   status, created_at, updated_at, version
		FROM wallets WHERE id = $1
	`

	var (
		walletID      string
		userID        string
		balance       int64
		currency      string
		escrowBalance int64
		status        string
		createdAt     time.Time
		updatedAt     time.Time
		version       int64
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&walletID, &userID, &balance, &currency, &escrowBalance,
		&status, &createdAt, &updatedAt, &version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrWalletNotFound
		}
		return nil, err
	}

	return reconstructWallet(
		walletID, userID, balance, currency, escrowBalance,
		status, createdAt, updatedAt, version,
	)
}

// FindByUserID retrieves a wallet by user ID
func (r *WalletRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.Wallet, error) {
	query := `
		SELECT id, user_id, balance, currency, escrow_balance,
			   status, created_at, updated_at, version
		FROM wallets WHERE user_id = $1
	`

	var (
		walletID      string
		uid           string
		balance       int64
		currency      string
		escrowBalance int64
		status        string
		createdAt     time.Time
		updatedAt     time.Time
		version       int64
	)

	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(
		&walletID, &uid, &balance, &currency, &escrowBalance,
		&status, &createdAt, &updatedAt, &version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrWalletNotFound
		}
		return nil, err
	}

	return reconstructWallet(
		walletID, uid, balance, currency, escrowBalance,
		status, createdAt, updatedAt, version,
	)
}

// GetBalance returns wallet balance
func (r *WalletRepository) GetBalance(ctx context.Context, walletID valueobject.WalletID) (valueobject.Money, error) {
	query := `SELECT balance, currency FROM wallets WHERE id = $1`

	var balance int64
	var currency string

	err := r.db.QueryRowContext(ctx, query, walletID.String()).Scan(&balance, &currency)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return valueobject.Money{}, repository.ErrWalletNotFound
		}
		return valueobject.Money{}, err
	}

	return valueobject.NewMoney(balance, valueobject.Currency(currency))
}

func reconstructWallet(
	id, userID string,
	balance int64, currency string,
	escrowBalance int64,
	status string,
	createdAt, updatedAt time.Time,
	version int64,
) (*aggregate.Wallet, error) {
	wid, err := valueobject.NewWalletID(id)
	if err != nil {
		return nil, err
	}

	uid, err := valueobject.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	bal, err := valueobject.NewMoney(balance, valueobject.Currency(currency))
	if err != nil {
		return nil, err
	}

	escrow, err := valueobject.NewMoney(escrowBalance, valueobject.Currency(currency))
	if err != nil {
		return nil, err
	}

	return aggregate.ReconstructWallet(
		wid, uid, bal, escrow,
		aggregate.WalletStatus(status),
		createdAt, updatedAt, version,
	), nil
}

// Ensure WalletRepository implements the interface
var _ repository.WalletRepository = (*WalletRepository)(nil)
