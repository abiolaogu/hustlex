package query

import (
	"context"
	"time"

	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/identity/service"
	"hustlex/internal/domain/shared/valueobject"
)

// GetUser retrieves a user by ID
type GetUser struct {
	UserID string
}

// GetUserByPhone retrieves a user by phone number
type GetUserByPhone struct {
	Phone string
}

// GetUserByUsername retrieves a user by username
type GetUserByUsername struct {
	Username string
}

// UserDTO represents a user for API responses
type UserDTO struct {
	ID           string     `json:"id"`
	Phone        string     `json:"phone"`
	Email        string     `json:"email,omitempty"`
	FullName     string     `json:"full_name"`
	Username     string     `json:"username,omitempty"`
	ProfileImage string     `json:"profile_image,omitempty"`
	Bio          string     `json:"bio,omitempty"`
	Location     string     `json:"location,omitempty"`
	State        string     `json:"state,omitempty"`
	Gender       string     `json:"gender,omitempty"`
	IsVerified   bool       `json:"is_verified"`
	IsActive     bool       `json:"is_active"`
	Tier         string     `json:"tier"`
	ReferralCode string     `json:"referral_code"`
	Skills       []SkillDTO `json:"skills,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// SkillDTO represents a skill for API responses
type SkillDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category,omitempty"`
	Proficiency string   `json:"proficiency"`
	YearsExp    int      `json:"years_experience"`
	IsVerified  bool     `json:"is_verified"`
	Portfolio   []string `json:"portfolio_urls,omitempty"`
}

// GetUserSkills retrieves a user's skills
type GetUserSkills struct {
	UserID string
}

// GetSkillCatalog retrieves the skill catalog
type GetSkillCatalog struct {
	Category string // optional filter
}

// SearchSkills searches skills by name
type SearchSkills struct {
	Query string
	Limit int
}

// SkillCatalogDTO represents a catalog skill
type SkillCatalogDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// GetReferralStats retrieves referral statistics
type GetReferralStats struct {
	UserID string
}

// ReferralStatsDTO represents referral statistics
type ReferralStatsDTO struct {
	UserID        string   `json:"user_id"`
	ReferralCode  string   `json:"referral_code"`
	TotalReferrals int64   `json:"total_referrals"`
	ReferredUsers []string `json:"referred_users,omitempty"` // user IDs
}

// UserQueryHandler handles user queries
type UserQueryHandler struct {
	userRepo      repository.UserRepository
	skillRepo     repository.SkillRepository
	userSkillRepo repository.UserSkillRepository
	referralRepo  repository.ReferralRepository
}

// NewUserQueryHandler creates a new query handler
func NewUserQueryHandler(
	userRepo repository.UserRepository,
	skillRepo repository.SkillRepository,
	userSkillRepo repository.UserSkillRepository,
	referralRepo repository.ReferralRepository,
) *UserQueryHandler {
	return &UserQueryHandler{
		userRepo:      userRepo,
		skillRepo:     skillRepo,
		userSkillRepo: userSkillRepo,
		referralRepo:  referralRepo,
	}
}

// HandleGetUser retrieves a user by ID
func (h *UserQueryHandler) HandleGetUser(ctx context.Context, q GetUser) (*UserDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	// Get user skills
	skills := make([]SkillDTO, 0)
	for _, s := range user.Skills() {
		skills = append(skills, SkillDTO{
			ID:          s.SkillID.String(),
			Name:        s.SkillName,
			Proficiency: string(s.Proficiency),
			YearsExp:    s.YearsExp,
			IsVerified:  s.IsVerified,
			Portfolio:   s.PortfolioURLs,
		})
	}

	return &UserDTO{
		ID:           user.ID().String(),
		Phone:        user.Phone().Masked(),
		Email:        user.Email().Masked(),
		FullName:     user.FullName().String(),
		Username:     user.Username(),
		ProfileImage: user.ProfileImage(),
		Bio:          user.Bio(),
		Location:     user.Location(),
		State:        user.State(),
		Gender:       user.Gender(),
		IsVerified:   user.IsVerified(),
		IsActive:     user.IsActive(),
		Tier:         user.Tier().String(),
		ReferralCode: user.ReferralCode(),
		Skills:       skills,
		CreatedAt:    user.CreatedAt(),
		LastLoginAt:  user.LastLoginAt(),
	}, nil
}

// HandleGetUserByPhone retrieves a user by phone
func (h *UserQueryHandler) HandleGetUserByPhone(ctx context.Context, q GetUserByPhone) (*UserDTO, error) {
	phone, err := valueobject.NewPhoneNumber(q.Phone)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	return &UserDTO{
		ID:           user.ID().String(),
		Phone:        user.Phone().Masked(),
		FullName:     user.FullName().String(),
		ProfileImage: user.ProfileImage(),
		IsVerified:   user.IsVerified(),
		IsActive:     user.IsActive(),
	}, nil
}

// HandleGetUserByUsername retrieves a user by username
func (h *UserQueryHandler) HandleGetUserByUsername(ctx context.Context, q GetUserByUsername) (*UserDTO, error) {
	user, err := h.userRepo.FindByUsername(ctx, q.Username)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	skills := make([]SkillDTO, 0)
	for _, s := range user.Skills() {
		skills = append(skills, SkillDTO{
			ID:          s.SkillID.String(),
			Name:        s.SkillName,
			Proficiency: string(s.Proficiency),
			YearsExp:    s.YearsExp,
			IsVerified:  s.IsVerified,
		})
	}

	return &UserDTO{
		ID:           user.ID().String(),
		Phone:        user.Phone().Masked(),
		FullName:     user.FullName().String(),
		Username:     user.Username(),
		ProfileImage: user.ProfileImage(),
		Bio:          user.Bio(),
		Location:     user.Location(),
		IsVerified:   user.IsVerified(),
		Tier:         user.Tier().String(),
		Skills:       skills,
		CreatedAt:    user.CreatedAt(),
	}, nil
}

// HandleGetSkillCatalog retrieves the skill catalog
func (h *UserQueryHandler) HandleGetSkillCatalog(ctx context.Context, q GetSkillCatalog) ([]SkillCatalogDTO, error) {
	var skills []*repository.Skill
	var err error

	if q.Category != "" {
		skills, err = h.skillRepo.FindByCategory(ctx, q.Category)
	} else {
		skills, err = h.skillRepo.FindAll(ctx)
	}

	if err != nil {
		return nil, err
	}

	dtos := make([]SkillCatalogDTO, len(skills))
	for i, s := range skills {
		dtos[i] = SkillCatalogDTO{
			ID:          s.ID,
			Name:        s.Name,
			Category:    s.Category,
			Description: s.Description,
			Icon:        s.Icon,
		}
	}

	return dtos, nil
}

// HandleSearchSkills searches skills by name
func (h *UserQueryHandler) HandleSearchSkills(ctx context.Context, q SearchSkills) ([]SkillCatalogDTO, error) {
	limit := q.Limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	skills, err := h.skillRepo.Search(ctx, q.Query, limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]SkillCatalogDTO, len(skills))
	for i, s := range skills {
		dtos[i] = SkillCatalogDTO{
			ID:          s.ID,
			Name:        s.Name,
			Category:    s.Category,
			Description: s.Description,
		}
	}

	return dtos, nil
}

// HandleGetReferralStats retrieves referral statistics
func (h *UserQueryHandler) HandleGetReferralStats(ctx context.Context, q GetReferralStats) (*ReferralStatsDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	count, err := h.referralRepo.GetReferralCount(ctx, userID)
	if err != nil {
		count = 0
	}

	return &ReferralStatsDTO{
		UserID:        user.ID().String(),
		ReferralCode:  user.ReferralCode(),
		TotalReferrals: count,
	}, nil
}

// CheckUserExists checks if a user exists by phone
type CheckUserExists struct {
	Phone string
}

// HandleCheckUserExists checks if a user exists
func (h *UserQueryHandler) HandleCheckUserExists(ctx context.Context, q CheckUserExists) (bool, error) {
	phone, err := valueobject.NewPhoneNumber(q.Phone)
	if err != nil {
		return false, err
	}

	return h.userRepo.ExistsByPhone(ctx, phone)
}
