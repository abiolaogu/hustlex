package postgres

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/abiolaogu/hustlex/apps/api/internal/domain/identity/repository"
	"github.com/abiolaogu/hustlex/apps/api/internal/domain/shared/valueobject"
	"github.com/abiolaogu/hustlex/apps/api/internal/models"
)

// SkillRepositoryImpl implements the SkillRepository interface using PostgreSQL
type SkillRepositoryImpl struct {
	db *gorm.DB
}

// NewSkillRepository creates a new instance of SkillRepositoryImpl
func NewSkillRepository(db *gorm.DB) repository.SkillRepository {
	return &SkillRepositoryImpl{db: db}
}

// Save persists a skill to the database
func (r *SkillRepositoryImpl) Save(ctx context.Context, skill *repository.Skill) error {
	model := r.toModel(skill)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return fmt.Errorf("failed to save skill: %w", err)
	}

	return nil
}

// FindByID retrieves a skill by its ID
func (r *SkillRepositoryImpl) FindByID(ctx context.Context, id valueobject.SkillID) (*repository.Skill, error) {
	var model models.Skill

	err := r.db.WithContext(ctx).First(&model, "id = ?", id.String()).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("skill not found with id %s", id.String())
		}
		return nil, fmt.Errorf("failed to find skill by id: %w", err)
	}

	return r.toDomain(&model), nil
}

// FindByCategory retrieves skills by category
func (r *SkillRepositoryImpl) FindByCategory(ctx context.Context, category string) ([]*repository.Skill, error) {
	var models []models.Skill

	err := r.db.WithContext(ctx).
		Where("category = ? AND is_active = ?", category, true).
		Order("name ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find skills by category: %w", err)
	}

	skills := make([]*repository.Skill, 0, len(models))
	for _, model := range models {
		skills = append(skills, r.toDomain(&model))
	}

	return skills, nil
}

// FindAll retrieves all active skills
func (r *SkillRepositoryImpl) FindAll(ctx context.Context) ([]*repository.Skill, error) {
	var models []models.Skill

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("category ASC, name ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find all skills: %w", err)
	}

	skills := make([]*repository.Skill, 0, len(models))
	for _, model := range models {
		skills = append(skills, r.toDomain(&model))
	}

	return skills, nil
}

// Search searches skills by name
func (r *SkillRepositoryImpl) Search(ctx context.Context, query string, limit int) ([]*repository.Skill, error) {
	var models []models.Skill

	err := r.db.WithContext(ctx).
		Where("is_active = ? AND (name ILIKE ? OR description ILIKE ?)", true, "%"+query+"%", "%"+query+"%").
		Order("name ASC").
		Limit(limit).
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search skills: %w", err)
	}

	skills := make([]*repository.Skill, 0, len(models))
	for _, model := range models {
		skills = append(skills, r.toDomain(&model))
	}

	return skills, nil
}

// toModel converts a Skill to a database model
func (r *SkillRepositoryImpl) toModel(skill *repository.Skill) *models.Skill {
	return &models.Skill{
		BaseModel: models.BaseModel{
			ID:        skill.ID,
			CreatedAt: skill.CreatedAt,
			UpdatedAt: skill.UpdatedAt,
		},
		Name:        skill.Name,
		Category:    skill.Category,
		Description: skill.Description,
		Icon:        skill.Icon,
		IsActive:    skill.IsActive,
	}
}

// toDomain converts a database model to a Skill
func (r *SkillRepositoryImpl) toDomain(model *models.Skill) *repository.Skill {
	return &repository.Skill{
		ID:          model.ID,
		Name:        model.Name,
		Category:    model.Category,
		Description: model.Description,
		Icon:        model.Icon,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
