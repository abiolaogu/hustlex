package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/shared/valueobject"
	"github.com/abiolaogu/hustlex/apps/api/internal/models"
)

// UserSkillRepositoryImpl implements the UserSkillRepository interface using PostgreSQL
type UserSkillRepositoryImpl struct {
	db *gorm.DB
}

// NewUserSkillRepository creates a new instance of UserSkillRepositoryImpl
func NewUserSkillRepository(db *gorm.DB) repository.UserSkillRepository {
	return &UserSkillRepositoryImpl{db: db}
}

// Save persists a user skill to the database
func (r *UserSkillRepositoryImpl) Save(ctx context.Context, userSkill *repository.UserSkillRecord) error {
	model := r.toModel(userSkill)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to save user skill: %w", err)
	}

	return nil
}

// Delete removes a user skill
func (r *UserSkillRepositoryImpl) Delete(ctx context.Context, userID valueobject.UserID, skillID valueobject.SkillID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND skill_id = ?", userID.String(), skillID.String()).
		Delete(&models.UserSkill{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete user skill: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user skill not found")
	}

	return nil
}

// FindByUserID retrieves all skills for a user
func (r *UserSkillRepositoryImpl) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*repository.UserSkillRecord, error) {
	var models []models.UserSkill

	err := r.db.WithContext(ctx).
		Preload("Skill").
		Where("user_id = ?", userID.String()).
		Order("created_at DESC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find user skills: %w", err)
	}

	userSkills := make([]*repository.UserSkillRecord, 0, len(models))
	for _, model := range models {
		userSkills = append(userSkills, r.toDomain(&model))
	}

	return userSkills, nil
}

// FindBySkillID retrieves all users with a specific skill
func (r *UserSkillRepositoryImpl) FindBySkillID(ctx context.Context, skillID valueobject.SkillID, limit, offset int) ([]*repository.UserSkillRecord, error) {
	var models []models.UserSkill

	err := r.db.WithContext(ctx).
		Preload("Skill").
		Preload("User").
		Where("skill_id = ?", skillID.String()).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find users by skill: %w", err)
	}

	userSkills := make([]*repository.UserSkillRecord, 0, len(models))
	for _, model := range models {
		userSkills = append(userSkills, r.toDomain(&model))
	}

	return userSkills, nil
}

// CountBySkillID counts users with a specific skill
func (r *UserSkillRepositoryImpl) CountBySkillID(ctx context.Context, skillID valueobject.SkillID) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.UserSkill{}).
		Where("skill_id = ?", skillID.String()).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count users by skill: %w", err)
	}

	return count, nil
}

// toModel converts a UserSkillRecord to a database model
func (r *UserSkillRepositoryImpl) toModel(userSkill *repository.UserSkillRecord) *models.UserSkill {
	return &models.UserSkill{
		BaseModel: models.BaseModel{
			ID:        userSkill.ID,
			CreatedAt: userSkill.CreatedAt,
			UpdatedAt: userSkill.UpdatedAt,
		},
		UserID:        userSkill.UserID,
		SkillID:       userSkill.SkillID,
		Proficiency:   userSkill.Proficiency,
		YearsExp:      userSkill.YearsExp,
		IsVerified:    userSkill.IsVerified,
		PortfolioURLs: userSkill.PortfolioURLs,
	}
}

// toDomain converts a database model to a UserSkillRecord
func (r *UserSkillRepositoryImpl) toDomain(model *models.UserSkill) *repository.UserSkillRecord {
	skillName := ""
	if model.Skill.ID != "" {
		skillName = model.Skill.Name
	}

	return &repository.UserSkillRecord{
		ID:            model.ID,
		UserID:        model.UserID,
		SkillID:       model.SkillID,
		SkillName:     skillName,
		Proficiency:   model.Proficiency,
		YearsExp:      model.YearsExp,
		IsVerified:    model.IsVerified,
		PortfolioURLs: model.PortfolioURLs,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
