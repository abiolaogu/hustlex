package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/shared/valueobject"
	"github.com/abiolaogu/hustlex/apps/api/internal/models"
)

// ReferralRepositoryImpl implements the ReferralRepository interface using PostgreSQL
type ReferralRepositoryImpl struct {
	db *gorm.DB
}

// NewReferralRepository creates a new instance of ReferralRepositoryImpl
func NewReferralRepository(db *gorm.DB) repository.ReferralRepository {
	return &ReferralRepositoryImpl{db: db}
}

// RecordReferral records that one user referred another
func (r *ReferralRepositoryImpl) RecordReferral(ctx context.Context, referrerID, referredID valueobject.UserID) error {
	// Update the referred user's ReferredBy field
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", referredID.String()).
		Update("referred_by", referrerID.String())

	if result.Error != nil {
		return fmt.Errorf("failed to record referral: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("referred user not found")
	}

	return nil
}

// GetReferralCount gets the count of successful referrals for a user
func (r *ReferralRepositoryImpl) GetReferralCount(ctx context.Context, userID valueobject.UserID) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("referred_by = ?", userID.String()).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to get referral count: %w", err)
	}

	return count, nil
}

// GetReferrals gets the list of users referred by a user
func (r *ReferralRepositoryImpl) GetReferrals(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]valueobject.UserID, error) {
	var users []models.User

	err := r.db.WithContext(ctx).
		Select("id").
		Where("referred_by = ?", userID.String()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get referrals: %w", err)
	}

	referredIDs := make([]valueobject.UserID, 0, len(users))
	for _, user := range users {
		id, err := valueobject.NewUserID(user.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		referredIDs = append(referredIDs, id)
	}

	return referredIDs, nil
}
