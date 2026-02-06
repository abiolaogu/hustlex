package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/aggregate"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/event"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/shared/valueobject"
	"github.com/abiolaogu/hustlex/apps/api/internal/models"
)

// UserRepositoryImpl implements the UserRepository interface using PostgreSQL
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save persists the user aggregate to the database
func (r *UserRepositoryImpl) Save(ctx context.Context, user *aggregate.User) error {
	model := r.toModel(user)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// SaveWithEvents persists the user aggregate and publishes domain events
func (r *UserRepositoryImpl) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Save the user
	model := r.toModel(user)
	if err := tx.Save(model).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save user: %w", err)
	}

	// TODO: Publish domain events to message queue
	// For now, we'll just clear the events after save
	// In production, this should publish to RabbitMQ or similar
	events := user.GetEvents()
	if len(events) > 0 {
		// Log events for debugging
		for _, evt := range events {
			switch e := evt.(type) {
			case *event.UserRegisteredEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserLoggedInEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserProfileUpdatedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserVerifiedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserTierUpgradedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserDeactivatedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.UserReactivatedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.SkillAddedEvent:
				// TODO: Publish to event bus
				_ = e
			case *event.SkillRemovedEvent:
				// TODO: Publish to event bus
				_ = e
			}
		}
		user.ClearEvents()
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error) {
	var model models.User

	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		First(&model, "id = ?", id.String()).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with id %s", id.String())
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return r.toDomain(&model)
}

// FindByPhone retrieves a user by their phone number
func (r *UserRepositoryImpl) FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error) {
	var model models.User

	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		First(&model, "phone = ?", phone.String()).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with phone %s", phone.String())
		}
		return nil, fmt.Errorf("failed to find user by phone: %w", err)
	}

	return r.toDomain(&model)
}

// FindByEmail retrieves a user by their email
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email valueobject.Email) (*aggregate.User, error) {
	var model models.User

	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		First(&model, "email = ?", email.String()).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with email %s", email.String())
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.toDomain(&model)
}

// FindByUsername retrieves a user by their username
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	var model models.User

	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		First(&model, "username = ?", username).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with username %s", username)
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return r.toDomain(&model)
}

// FindByReferralCode retrieves a user by their referral code
func (r *UserRepositoryImpl) FindByReferralCode(ctx context.Context, code string) (*aggregate.User, error) {
	var model models.User

	err := r.db.WithContext(ctx).
		Preload("Skills").
		Preload("Skills.Skill").
		First(&model, "referral_code = ?", code).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with referral code %s", code)
		}
		return nil, fmt.Errorf("failed to find user by referral code: %w", err)
	}

	return r.toDomain(&model)
}

// ExistsByPhone checks if a user exists with the given phone number
func (r *UserRepositoryImpl) ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("phone = ?", phone.String()).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by phone: %w", err)
	}

	return count > 0, nil
}

// ExistsByEmail checks if a user exists with the given email
func (r *UserRepositoryImpl) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email.String()).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by email: %w", err)
	}

	return count > 0, nil
}

// ExistsByUsername checks if a user exists with the given username
func (r *UserRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check user existence by username: %w", err)
	}

	return count > 0, nil
}

// Delete soft deletes a user by their ID
func (r *UserRepositoryImpl) Delete(ctx context.Context, id valueobject.UserID) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id.String())

	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found with id %s", id.String())
	}

	return nil
}

// toModel converts a User aggregate to a database model
func (r *UserRepositoryImpl) toModel(user *aggregate.User) *models.User {
	model := &models.User{
		BaseModel: models.BaseModel{
			ID:        user.ID().String(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		},
		Phone:        user.Phone().String(),
		FullName:     user.FullName().String(),
		Username:     user.Username(),
		ProfileImage: user.ProfileImage(),
		Bio:          user.Bio(),
		Location:     user.Location(),
		State:        user.State(),
		IsVerified:   user.IsVerified(),
		IsActive:     user.IsActive(),
		Tier:         models.UserTier(user.Tier().String()),
		ReferralCode: user.ReferralCode(),
	}

	// Handle optional email
	if !user.Email().IsEmpty() {
		email := user.Email().String()
		model.Email = email
	}

	// Handle optional date of birth
	if !user.DateOfBirth().IsZero() {
		dob := user.DateOfBirth()
		model.DateOfBirth = &dob
	}

	// Handle optional gender
	if user.Gender() != "" {
		model.Gender = user.Gender()
	}

	// Handle optional referred by
	if user.ReferredBy() != nil {
		referredBy := user.ReferredBy().String()
		model.ReferredBy = &referredBy
	}

	// Handle optional last login
	if user.LastLoginAt() != nil && !user.LastLoginAt().IsZero() {
		lastLogin := *user.LastLoginAt()
		model.LastLoginAt = &lastLogin
	}

	// Convert user skills
	if len(user.Skills()) > 0 {
		model.Skills = make([]models.UserSkill, 0, len(user.Skills()))
		for _, skill := range user.Skills() {
			userSkill := models.UserSkill{
				BaseModel: models.BaseModel{
					CreatedAt: skill.AddedAt,
					UpdatedAt: skill.AddedAt,
				},
				UserID:       user.ID().String(),
				SkillID:      skill.SkillID.String(),
				Proficiency:  string(skill.Proficiency),
				YearsExp:     skill.YearsExp,
				IsVerified:   skill.IsVerified,
				PortfolioURLs: skill.PortfolioURLs,
			}
			model.Skills = append(model.Skills, userSkill)
		}
	}

	return model
}

// toDomain converts a database model to a User aggregate
func (r *UserRepositoryImpl) toDomain(model *models.User) (*aggregate.User, error) {
	// Parse required value objects
	id, err := valueobject.NewUserID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	phone, err := valueobject.NewPhoneNumber(model.Phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	fullName, err := valueobject.NewFullName(model.FullName)
	if err != nil {
		return nil, fmt.Errorf("invalid full name: %w", err)
	}

	// Parse optional email
	var email valueobject.Email
	if model.Email != "" {
		email, err = valueobject.NewEmail(model.Email)
		if err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
	}

	// Parse tier
	tier := parseUserTier(string(model.Tier))

	// Parse optional referred by
	var referredBy *valueobject.UserID
	if model.ReferredBy != nil && *model.ReferredBy != "" {
		rb, err := valueobject.NewUserID(*model.ReferredBy)
		if err != nil {
			return nil, fmt.Errorf("invalid referred by id: %w", err)
		}
		referredBy = &rb
	}

	// Reconstruct skills
	var skills []aggregate.UserSkill
	if len(model.Skills) > 0 {
		skills = make([]aggregate.UserSkill, 0, len(model.Skills))
		for _, userSkillModel := range model.Skills {
			// Parse skill ID
			skillID, err := valueobject.NewSkillID(userSkillModel.SkillID)
			if err != nil {
				return nil, fmt.Errorf("invalid skill id: %w", err)
			}

			// Parse proficiency
			proficiency := parseProficiency(userSkillModel.Proficiency)

			// Get skill name from preloaded data (if available)
			skillName := ""
			if userSkillModel.Skill.ID != "" {
				skillName = userSkillModel.Skill.Name
			}

			skill := aggregate.UserSkill{
				SkillID:       skillID,
				SkillName:     skillName,
				Proficiency:   proficiency,
				YearsExp:      userSkillModel.YearsExp,
				IsVerified:    userSkillModel.IsVerified,
				PortfolioURLs: userSkillModel.PortfolioURLs,
				AddedAt:       userSkillModel.CreatedAt,
			}
			skills = append(skills, skill)
		}
	}

	// Get version (assuming it's stored in BaseModel or defaults to 1)
	version := int64(1)

	// Reconstruct the user aggregate
	user := aggregate.ReconstructUser(
		id,
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
		tier,
		model.ReferralCode,
		referredBy,
		skills,
		model.LastLoginAt,
		model.CreatedAt,
		model.UpdatedAt,
		version,
	)

	return user, nil
}

// parseUserTier converts a string to UserTier
func parseUserTier(tier string) aggregate.UserTier {
	switch tier {
	case "bronze":
		return aggregate.TierBronze
	case "silver":
		return aggregate.TierSilver
	case "gold":
		return aggregate.TierGold
	case "platinum":
		return aggregate.TierPlatinum
	default:
		return aggregate.TierBronze // Default to bronze if invalid
	}
}

// parseProficiency converts a string to Proficiency
func parseProficiency(proficiency string) aggregate.Proficiency {
	switch proficiency {
	case "beginner":
		return aggregate.ProficiencyBeginner
	case "intermediate":
		return aggregate.ProficiencyIntermediate
	case "expert":
		return aggregate.ProficiencyExpert
	default:
		return aggregate.ProficiencyIntermediate // Default to intermediate if invalid
	}
}
