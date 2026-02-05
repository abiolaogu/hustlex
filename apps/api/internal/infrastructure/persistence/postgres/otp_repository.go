package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"hustlex/internal/domain/identity/repository"
)

// OTPRepository is a PostgreSQL implementation of the OTPRepository interface
type OTPRepository struct {
	db *sql.DB
}

// NewOTPRepository creates a new PostgreSQL OTP repository
func NewOTPRepository(db *sql.DB) repository.OTPRepository {
	return &OTPRepository{db: db}
}

// Save persists an OTP code
func (r *OTPRepository) Save(ctx context.Context, otp *repository.OTPCode) error {
	query := `
		INSERT INTO otps (
			id, phone, code_hash, purpose, expires_at, is_used, attempts, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		otp.ID,
		otp.Phone,
		otp.Code, // In production, this should be hashed
		otp.Purpose,
		otp.ExpiresAt,
		otp.IsUsed,
		otp.Attempts,
		otp.CreatedAt,
	)

	return err
}

// FindLatestValid finds the latest valid (unused, unexpired) OTP
func (r *OTPRepository) FindLatestValid(ctx context.Context, phone string, purpose string) (*repository.OTPCode, error) {
	query := `
		SELECT id, phone, code_hash, purpose, expires_at, is_used, attempts, created_at
		FROM otps
		WHERE phone = $1
			AND purpose = $2
			AND is_used = false
			AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	var otp repository.OTPCode
	err := r.db.QueryRowContext(ctx, query, phone, purpose).Scan(
		&otp.ID,
		&otp.Phone,
		&otp.Code,
		&otp.Purpose,
		&otp.ExpiresAt,
		&otp.IsUsed,
		&otp.Attempts,
		&otp.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no valid OTP found")
	}
	if err != nil {
		return nil, err
	}

	return &otp, nil
}

// MarkUsed marks an OTP as used
func (r *OTPRepository) MarkUsed(ctx context.Context, id string) error {
	query := `
		UPDATE otps
		SET is_used = true
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("OTP not found")
	}

	return nil
}

// IncrementAttempts increments the failed attempt counter
func (r *OTPRepository) IncrementAttempts(ctx context.Context, id string) error {
	query := `
		UPDATE otps
		SET attempts = attempts + 1
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteExpired removes expired OTPs (for cleanup)
func (r *OTPRepository) DeleteExpired(ctx context.Context) error {
	query := `
		DELETE FROM otps
		WHERE expires_at < NOW()
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// DeleteUnused deletes unused OTPs for a phone/purpose combination
func (r *OTPRepository) DeleteUnused(ctx context.Context, phone, purpose string) error {
	query := `
		DELETE FROM otps
		WHERE phone = $1
			AND purpose = $2
			AND is_used = false
	`

	_, err := r.db.ExecContext(ctx, query, phone, purpose)
	return err
}
