package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/models"
)

// OTPRepositoryImpl implements the OTPRepository interface using PostgreSQL
type OTPRepositoryImpl struct {
	db *gorm.DB
}

// NewOTPRepository creates a new instance of OTPRepositoryImpl
func NewOTPRepository(db *gorm.DB) repository.OTPRepository {
	return &OTPRepositoryImpl{db: db}
}

// Save persists an OTP code to the database
func (r *OTPRepositoryImpl) Save(ctx context.Context, otp *repository.OTPCode) error {
	model := r.toModel(otp)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	return nil
}

// FindLatestValid retrieves the latest valid OTP for a phone number and purpose
func (r *OTPRepositoryImpl) FindLatestValid(ctx context.Context, phone string, purpose string) (*repository.OTPCode, error) {
	var model models.OTPCode

	err := r.db.WithContext(ctx).
		Where("phone = ? AND purpose = ? AND is_used = ? AND expires_at > ?", phone, purpose, false, time.Now()).
		Order("created_at DESC").
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("OTP not found for phone %s", phone)
		}
		return nil, fmt.Errorf("failed to find OTP: %w", err)
	}

	return r.toDomain(&model), nil
}

// MarkUsed marks an OTP as used
func (r *OTPRepositoryImpl) MarkUsed(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_used":     true,
			"verified_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("OTP not found with id %s", id)
	}

	return nil
}

// IncrementAttempts increments the attempt count for an OTP
func (r *OTPRepositoryImpl) IncrementAttempts(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("id = ?", id).
		Update("attempts", gorm.Expr("attempts + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("failed to increment OTP attempts: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("OTP not found with id %s", id)
	}

	return nil
}

// DeleteExpired deletes all expired OTP codes
func (r *OTPRepositoryImpl) DeleteExpired(ctx context.Context) error {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.OTPCode{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete expired OTPs: %w", result.Error)
	}

	return nil
}

// DeleteUnused deletes unused OTPs for a phone/purpose combination
func (r *OTPRepositoryImpl) DeleteUnused(ctx context.Context, phone, purpose string) error {
	result := r.db.WithContext(ctx).
		Where("phone = ? AND purpose = ? AND is_used = ?", phone, purpose, false).
		Delete(&models.OTPCode{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete unused OTPs: %w", result.Error)
	}

	return nil
}

// toModel converts an OTPCode to a database model
func (r *OTPRepositoryImpl) toModel(otp *repository.OTPCode) *models.OTPCode {
	model := &models.OTPCode{
		BaseModel: models.BaseModel{
			ID:        otp.ID,
			CreatedAt: otp.CreatedAt,
		},
		Phone:     otp.Phone,
		Code:      otp.Code,
		Purpose:   otp.Purpose,
		ExpiresAt: otp.ExpiresAt,
		Attempts:  otp.Attempts,
		IsUsed:    otp.IsUsed,
	}

	return model
}

// toDomain converts a database model to an OTPCode
func (r *OTPRepositoryImpl) toDomain(model *models.OTPCode) *repository.OTPCode {
	return &repository.OTPCode{
		ID:        model.ID,
		Phone:     model.Phone,
		Code:      model.Code,
		Purpose:   model.Purpose,
		ExpiresAt: model.ExpiresAt,
		IsUsed:    model.IsUsed,
		Attempts:  model.Attempts,
		CreatedAt: model.CreatedAt,
	}
}
