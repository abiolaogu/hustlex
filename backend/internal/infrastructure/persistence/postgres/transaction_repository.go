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

// TransactionRepository implements repository.TransactionRepository using PostgreSQL
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new PostgreSQL transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Save persists a transaction
func (r *TransactionRepository) Save(ctx context.Context, txn *aggregate.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, wallet_id, type, amount, currency, balance_after,
			reference, description, metadata, status,
			counterparty_wallet_id, counterparty_name,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	var counterpartyWalletID, counterpartyName sql.NullString
	if txn.CounterpartyWalletID() != "" {
		counterpartyWalletID = sql.NullString{String: txn.CounterpartyWalletID(), Valid: true}
	}
	if txn.CounterpartyName() != "" {
		counterpartyName = sql.NullString{String: txn.CounterpartyName(), Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		txn.ID().String(),
		txn.WalletID().String(),
		txn.Type().String(),
		txn.Amount().Amount(),
		string(txn.Amount().Currency()),
		txn.BalanceAfter().Amount(),
		txn.Reference(),
		txn.Description(),
		txn.Metadata(),
		txn.Status().String(),
		counterpartyWalletID,
		counterpartyName,
		txn.CreatedAt(),
	)

	return err
}

// FindByID retrieves a transaction by ID
func (r *TransactionRepository) FindByID(ctx context.Context, id valueobject.TransactionID) (*aggregate.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, currency, balance_after,
			   reference, description, metadata, status,
			   counterparty_wallet_id, counterparty_name, created_at
		FROM transactions WHERE id = $1
	`

	var (
		txnID               string
		walletID            string
		txnType             string
		amount              int64
		currency            string
		balanceAfter        int64
		reference           string
		description         string
		metadata            string
		status              string
		counterpartyWallet  sql.NullString
		counterpartyName    sql.NullString
		createdAt           time.Time
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&txnID, &walletID, &txnType, &amount, &currency, &balanceAfter,
		&reference, &description, &metadata, &status,
		&counterpartyWallet, &counterpartyName, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrTransactionNotFound
		}
		return nil, err
	}

	return reconstructTransaction(
		txnID, walletID, txnType, amount, currency, balanceAfter,
		reference, description, metadata, status,
		counterpartyWallet.String, counterpartyName.String, createdAt,
	)
}

// FindByWalletID retrieves transactions for a wallet
func (r *TransactionRepository) FindByWalletID(ctx context.Context, walletID valueobject.WalletID, filter repository.TransactionFilter) ([]*repository.TransactionDTO, int64, error) {
	// Build query with filters
	baseQuery := `FROM transactions WHERE wallet_id = $1`
	args := []interface{}{walletID.String()}
	argPos := 2

	if filter.Type != nil {
		baseQuery += ` AND type = $` + string(rune('0'+argPos))
		args = append(args, filter.Type.String())
		argPos++
	}

	if filter.Status != nil {
		baseQuery += ` AND status = $` + string(rune('0'+argPos))
		args = append(args, filter.Status.String())
		argPos++
	}

	if filter.FromDate != nil {
		baseQuery += ` AND created_at >= $` + string(rune('0'+argPos))
		args = append(args, *filter.FromDate)
		argPos++
	}

	if filter.ToDate != nil {
		baseQuery += ` AND created_at <= $` + string(rune('0'+argPos))
		args = append(args, *filter.ToDate)
		argPos++
	}

	// Get total count
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get transactions with pagination
	selectQuery := `
		SELECT id, wallet_id, type, amount, currency, balance_after,
			   reference, description, status, counterparty_name, created_at
		` + baseQuery + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argPos)) + ` OFFSET $` + string(rune('0'+argPos+1))

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*repository.TransactionDTO
	for rows.Next() {
		var dto repository.TransactionDTO
		var counterpartyName sql.NullString

		err := rows.Scan(
			&dto.ID, &dto.WalletID, &dto.Type, &dto.Amount, &dto.Currency,
			&dto.BalanceAfter, &dto.Reference, &dto.Description, &dto.Status,
			&counterpartyName, &dto.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		dto.CounterpartyName = counterpartyName.String
		transactions = append(transactions, &dto)
	}

	return transactions, total, nil
}

// FindByReference retrieves a transaction by reference
func (r *TransactionRepository) FindByReference(ctx context.Context, reference string) (*aggregate.Transaction, error) {
	query := `
		SELECT id, wallet_id, type, amount, currency, balance_after,
			   reference, description, metadata, status,
			   counterparty_wallet_id, counterparty_name, created_at
		FROM transactions WHERE reference = $1
	`

	var (
		txnID               string
		walletID            string
		txnType             string
		amount              int64
		currency            string
		balanceAfter        int64
		ref                 string
		description         string
		metadata            string
		status              string
		counterpartyWallet  sql.NullString
		counterpartyName    sql.NullString
		createdAt           time.Time
	)

	err := r.db.QueryRowContext(ctx, query, reference).Scan(
		&txnID, &walletID, &txnType, &amount, &currency, &balanceAfter,
		&ref, &description, &metadata, &status,
		&counterpartyWallet, &counterpartyName, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrTransactionNotFound
		}
		return nil, err
	}

	return reconstructTransaction(
		txnID, walletID, txnType, amount, currency, balanceAfter,
		ref, description, metadata, status,
		counterpartyWallet.String, counterpartyName.String, createdAt,
	)
}

func reconstructTransaction(
	id, walletID, txnType string,
	amount int64, currency string,
	balanceAfter int64,
	reference, description, metadata, status string,
	counterpartyWalletID, counterpartyName string,
	createdAt time.Time,
) (*aggregate.Transaction, error) {
	tid, err := valueobject.NewTransactionID(id)
	if err != nil {
		return nil, err
	}

	wid, err := valueobject.NewWalletID(walletID)
	if err != nil {
		return nil, err
	}

	amt, err := valueobject.NewMoney(amount, valueobject.Currency(currency))
	if err != nil {
		return nil, err
	}

	balAfter, err := valueobject.NewMoney(balanceAfter, valueobject.Currency(currency))
	if err != nil {
		return nil, err
	}

	return aggregate.ReconstructTransaction(
		tid, wid,
		aggregate.TransactionType(txnType),
		amt, balAfter,
		reference, description, metadata,
		aggregate.TransactionStatus(status),
		counterpartyWalletID, counterpartyName,
		createdAt,
	), nil
}

// Ensure TransactionRepository implements the interface
var _ repository.TransactionRepository = (*TransactionRepository)(nil)
