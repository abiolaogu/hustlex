package postgres

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
	"hustlex/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository implements the domain UserRepository interface using PostgreSQL
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Save persists a user aggregate
func (r *UserRepository) Save(ctx context.Context, user *aggregate.User) error {
	model := r.toModel(user)

	// Check if user exists
	var existing models.User
	err := r.db.WithContext(ctx).Where("id = ?", model.ID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new user
		return r.db.WithContext(ctx).Create(model).Error
	} else if err != nil {
		return err
	}

	// Update existing user
	return r.db.WithContext(ctx).Model(&existing).Updates(model).Error
}

// SaveWithEvents persists a user and publishes domain events
func (r *UserRepository) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
	// For now, just save the user
	// TODO: Implement event publishing with transaction
	return r.Save(ctx, user)
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		Where("id = ?", id.String()).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&model)
}

// FindByPhone retrieves a user by phone number
func (r *UserRepository) FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		Where("phone = ?", phone.String()).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&model)
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*aggregate.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		Where("email = ?", email.String()).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&model)
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		Where("username = ?", username).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&model)
}

// FindByReferralCode retrieves a user by their referral code
func (r *UserRepository) FindByReferralCode(ctx context.Context, code string) (*aggregate.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		Where("referral_code = ?", code).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return r.toDomain(&model)
}

// ExistsByPhone checks if a user with the phone exists
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("phone = ?", phone.String()).
		Count(&count).Error

	return count > 0, err
}

// ExistsByEmail checks if a user with the email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email.String()).
		Count(&count).Error

	return count > 0, err
}

// ExistsByUsername checks if a user with the username exists
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error

	return count > 0, err
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	uid, err := uuid.Parse(id.String())
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Where("id = ?", uid).
		Delete(&models.User{}).Error
}

// toModel converts domain aggregate to database model
func (r *UserRepository) toModel(user *aggregate.User) *models.User {
	uid, _ := uuid.Parse(user.ID().String())

	model := &models.User{
		BaseModel: models.BaseModel{
			ID:        uid,
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		},
		Phone:        user.Phone().String(),
		Email:        user.Email().String(),
		FullName:     user.FullName().String(),
		Username:     user.Username(),
		ProfileImage: user.ProfileImage(),
		Bio:          user.Bio(),
		Location:     user.Location(),
		State:        user.State(),
		Gender:       user.Gender(),
		IsVerified:   user.IsVerified(),
		IsActive:     user.IsActive(),
		Tier:         models.UserTier(user.Tier().String()),
		ReferralCode: user.ReferralCode(),
		LastLoginAt:  user.LastLoginAt(),
	}

	// Handle nullable fields
	if user.DateOfBirth() != nil {
		dob := *user.DateOfBirth()
		model.DateOfBirth = &dob
	}

	if user.ReferredBy() != nil {
		referredByID, _ := uuid.Parse(user.ReferredBy().String())
		model.ReferredBy = &referredByID
	}

	return model
}

// toDomain converts database model to domain aggregate
func (r *UserRepository) toDomain(model *models.User) (*aggregate.User, error) {
	// Create value objects
	userID, err := valueobject.NewUserID(model.ID.String())
	if err != nil {
		return nil, err
	}

	phone, err := valueobject.NewPhoneNumber(model.Phone)
	if err != nil {
		return nil, err
	}

	fullName, err := valueobject.NewFullName(model.FullName)
	if err != nil {
		return nil, err
	}

	email := valueobject.EmptyEmail()
	if model.Email != "" {
		email, err = valueobject.NewEmail(model.Email)
		if err != nil {
			return nil, err
		}
	}

	// Reconstruct user using domain factory
	user := aggregate.ReconstructUser(
		userID,
		phone,
		email,
		fullName,
		model.Username,
		model.ProfileImage,
		model.Bio,
		model.Location,
		model.State,
		model.DateOfBirth,
		model.Gender,
		model.IsVerified,
		model.IsActive,
		aggregate.UserTier(model.Tier),
		model.ReferralCode,
		r.convertReferredBy(model.ReferredBy),
		r.convertSkills(model.Skills),
		model.LastLoginAt,
		model.CreatedAt,
		model.UpdatedAt,
		1, // version
	)

	return user, nil
}

// convertReferredBy converts UUID pointer to UserID pointer
func (r *UserRepository) convertReferredBy(referredBy *uuid.UUID) *valueobject.UserID {
	if referredBy == nil {
		return nil
	}

	userID, err := valueobject.NewUserID(referredBy.String())
	if err != nil {
		return nil
	}

	return &userID
}

// convertSkills converts model skills to domain skills
func (r *UserRepository) convertSkills(modelSkills []models.UserSkill) []aggregate.UserSkill {
	skills := make([]aggregate.UserSkill, 0, len(modelSkills))

	for _, ms := range modelSkills {
		skillID, err := valueobject.NewSkillID(ms.SkillID.String())
		if err != nil {
			continue
		}

		skill := aggregate.UserSkill{
			SkillID:       skillID,
			SkillName:     ms.Skill.Name,
			Proficiency:   aggregate.Proficiency(ms.Proficiency),
			YearsExp:      ms.YearsExp,
			IsVerified:    ms.IsVerified,
			PortfolioURLs: ms.PortfolioURLs,
			AddedAt:       ms.CreatedAt,
		}

		skills = append(skills, skill)
	}

	return skills
}

// OTPRepository implements OTP persistence using PostgreSQL
type OTPRepository struct {
	db *gorm.DB
}

// NewOTPRepository creates a new PostgreSQL OTP repository
func NewOTPRepository(db *gorm.DB) repository.OTPRepository {
	return &OTPRepository{db: db}
}

// Save persists an OTP code
func (r *OTPRepository) Save(ctx context.Context, otp *repository.OTPCode) error {
	model := &models.OTPCode{
		BaseModel: models.BaseModel{
			ID: uuid.MustParse(otp.ID),
		},
		Phone:     otp.Phone,
		Code:      otp.Code,
		Purpose:   otp.Purpose,
		ExpiresAt: otp.ExpiresAt,
		IsUsed:    otp.IsUsed,
		Attempts:  otp.Attempts,
		CreatedAt: otp.CreatedAt,
	}

	return r.db.WithContext(ctx).Create(model).Error
}

// FindLatestValid finds the latest valid (unused, unexpired) OTP
func (r *OTPRepository) FindLatestValid(ctx context.Context, phone string, purpose string) (*repository.OTPCode, error) {
	var model models.OTPCode
	err := r.db.WithContext(ctx).
		Where("phone = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
			phone, purpose, false, time.Now().UTC()).
		Order("created_at DESC").
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("OTP not found")
	} else if err != nil {
		return nil, err
	}

	return &repository.OTPCode{
		ID:        model.ID.String(),
		Phone:     model.Phone,
		Code:      model.Code,
		Purpose:   model.Purpose,
		ExpiresAt: model.ExpiresAt,
		IsUsed:    model.IsUsed,
		Attempts:  model.Attempts,
		CreatedAt: model.CreatedAt,
	}, nil
}

// MarkUsed marks an OTP as used
func (r *OTPRepository) MarkUsed(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("id = ?", uid).
		Update("is_used", true).Error
}

// IncrementAttempts increments the failed attempt counter
func (r *OTPRepository) IncrementAttempts(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).
		Model(&models.OTPCode{}).
		Where("id = ?", uid).
		UpdateColumn("attempts", gorm.Expr("attempts + ?", 1)).Error
}

// DeleteExpired removes expired OTPs (for cleanup)
func (r *OTPRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now().UTC()).
		Delete(&models.OTPCode{}).Error
}

// DeleteUnused deletes unused OTPs for a phone/purpose combination
func (r *OTPRepository) DeleteUnused(ctx context.Context, phone, purpose string) error {
	return r.db.WithContext(ctx).
		Where("phone = ? AND purpose = ? AND is_used = ?", phone, purpose, false).
		Delete(&models.OTPCode{}).Error
}
