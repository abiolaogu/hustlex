package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"hustlex/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SavingsService handles savings circle (Ajo/Esusu) operations
type SavingsService struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewSavingsService creates a new savings service
func NewSavingsService(db *gorm.DB, redis *redis.Client) *SavingsService {
	return &SavingsService{
		db:    db,
		redis: redis,
	}
}

// CreateCircleInput represents input for creating a savings circle
type CreateCircleInput struct {
	Name            string   `json:"name" validate:"required,min=3,max=100"`
	Description     string   `json:"description" validate:"omitempty,max=500"`
	Type            string   `json:"type" validate:"required,oneof=rotational fixed_target emergency"`
	ContributionAmt int64    `json:"contribution_amount" validate:"required,min=100"` // min ₦1 (100 kobo)
	Frequency       string   `json:"frequency" validate:"required,oneof=daily weekly biweekly monthly"`
	MaxMembers      int      `json:"max_members" validate:"required,min=2,max=30"`
	TotalRounds     int      `json:"total_rounds" validate:"required,min=1"`
	StartDate       string   `json:"start_date" validate:"omitempty"`
	IsPrivate       bool     `json:"is_private"`
	Rules           []string `json:"rules" validate:"omitempty,max=10,dive,max=200"`
}

// CircleFilters represents filters for listing circles
type CircleFilters struct {
	Type      string `query:"type"`
	Status    string `query:"status"`
	MinAmount int64  `query:"min_amount"`
	MaxAmount int64  `query:"max_amount"`
	Frequency string `query:"frequency"`
	Search    string `query:"search"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}

// CircleListResult represents paginated circle results
type CircleListResult struct {
	Circles    []models.SavingsCircle `json:"circles"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// CreateCircle creates a new savings circle
func (s *SavingsService) CreateCircle(ctx context.Context, creatorID uuid.UUID, input *CreateCircleInput) (*models.SavingsCircle, error) {
	// Generate invite code
	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invite code: %w", err)
	}

	circle := &models.SavingsCircle{
		Name:            input.Name,
		Description:     input.Description,
		Type:            models.CircleType(input.Type),
		ContributionAmt: input.ContributionAmt,
		Currency:        "NGN",
		Frequency:       input.Frequency,
		MaxMembers:      input.MaxMembers,
		TotalRounds:     input.TotalRounds,
		CreatedBy:       creatorID,
		Status:          "recruiting",
		IsPrivate:       input.IsPrivate,
		InviteCode:      inviteCode,
		Rules:           input.Rules,
		CurrentMembers:  1, // Creator is the first member
	}

	// Parse start date if provided
	if input.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", input.StartDate)
		if err == nil {
			circle.StartDate = &startDate
		}
	}

	tx := s.db.Begin()

	if err := tx.Create(circle).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create circle: %w", err)
	}

	// Add creator as admin member
	member := &models.CircleMember{
		CircleID:     circle.ID,
		UserID:       creatorID,
		Position:     1, // First position (can be randomized later)
		Role:         "admin",
		Status:       "active",
		TotalContrib: 0,
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	tx.Commit()

	// Load relationships
	s.db.Preload("Creator").Preload("Members").Preload("Members.User").First(circle, circle.ID)

	return circle, nil
}

// GetCircle retrieves a single circle by ID
func (s *SavingsService) GetCircle(ctx context.Context, circleID uuid.UUID, userID uuid.UUID) (*models.SavingsCircle, error) {
	var circle models.SavingsCircle

	err := s.db.
		Preload("Creator").
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Members.User").
		First(&circle, circleID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("circle not found")
		}
		return nil, err
	}

	// Check if user is a member (for private circles)
	if circle.IsPrivate {
		isMember := false
		for _, m := range circle.Members {
			if m.UserID == userID {
				isMember = true
				break
			}
		}
		if !isMember {
			return nil, errors.New("you don't have access to this circle")
		}
	}

	return &circle, nil
}

// GetCircleByInviteCode retrieves a circle by its invite code
func (s *SavingsService) GetCircleByInviteCode(ctx context.Context, inviteCode string) (*models.SavingsCircle, error) {
	var circle models.SavingsCircle

	err := s.db.
		Preload("Creator").
		Where("invite_code = ?", inviteCode).
		First(&circle).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid invite code")
		}
		return nil, err
	}

	return &circle, nil
}

// ListCircles lists available circles with filters
func (s *SavingsService) ListCircles(ctx context.Context, filters *CircleFilters) (*CircleListResult, error) {
	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 50 {
		filters.Limit = 20
	}

	query := s.db.Model(&models.SavingsCircle{}).
		Preload("Creator").
		Where("is_private = ?", false) // Only show public circles

	// Apply filters
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	} else {
		query = query.Where("status = ?", "recruiting") // Default to recruiting
	}

	if filters.MinAmount > 0 {
		query = query.Where("contribution_amt >= ?", filters.MinAmount)
	}

	if filters.MaxAmount > 0 {
		query = query.Where("contribution_amt <= ?", filters.MaxAmount)
	}

	if filters.Frequency != "" {
		query = query.Where("frequency = ?", filters.Frequency)
	}

	if filters.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit).Order("created_at DESC")

	var circles []models.SavingsCircle
	if err := query.Find(&circles).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		totalPages++
	}

	return &CircleListResult{
		Circles:    circles,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetMyCircles gets circles the user is a member of
func (s *SavingsService) GetMyCircles(ctx context.Context, userID uuid.UUID, status string) ([]models.SavingsCircle, error) {
	query := s.db.
		Joins("JOIN circle_members ON circle_members.circle_id = savings_circles.id").
		Where("circle_members.user_id = ? AND circle_members.status = ?", userID, "active").
		Preload("Creator").
		Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Preload("Members.User")

	if status != "" {
		query = query.Where("savings_circles.status = ?", status)
	}

	var circles []models.SavingsCircle
	err := query.Order("savings_circles.created_at DESC").Find(&circles).Error

	return circles, err
}

// JoinCircle adds a user to a savings circle
func (s *SavingsService) JoinCircle(ctx context.Context, circleID, userID uuid.UUID) (*models.CircleMember, error) {
	var circle models.SavingsCircle
	if err := s.db.First(&circle, circleID).Error; err != nil {
		return nil, errors.New("circle not found")
	}

	// Check if circle is accepting new members
	if circle.Status != "recruiting" {
		return nil, errors.New("this circle is no longer accepting new members")
	}

	// Check if circle is full
	if circle.CurrentMembers >= circle.MaxMembers {
		return nil, errors.New("this circle is full")
	}

	// Check if user is already a member
	var existing models.CircleMember
	err := s.db.Where("circle_id = ? AND user_id = ?", circleID, userID).First(&existing).Error
	if err == nil {
		if existing.Status == "active" {
			return nil, errors.New("you are already a member of this circle")
		}
		// Rejoin if previously left
		existing.Status = "active"
		s.db.Save(&existing)
		s.db.Model(&circle).UpdateColumn("current_members", gorm.Expr("current_members + ?", 1))
		return &existing, nil
	}

	// Determine position (next available)
	position := circle.CurrentMembers + 1

	tx := s.db.Begin()

	member := &models.CircleMember{
		CircleID:     circleID,
		UserID:       userID,
		Position:     position,
		Role:         "member",
		Status:       "active",
		TotalContrib: 0,
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to join circle: %w", err)
	}

	// Update member count
	tx.Model(&circle).UpdateColumn("current_members", gorm.Expr("current_members + ?", 1))

	// If circle is now full, start it
	if circle.CurrentMembers+1 >= circle.MaxMembers && circle.Status == "recruiting" {
		tx.Model(&circle).Updates(map[string]interface{}{
			"status":     "active",
			"start_date": time.Now(),
		})
		// Schedule first contributions
		s.scheduleContributions(tx, &circle)
	}

	tx.Commit()

	// Load user relationship
	s.db.Preload("User").First(member, member.ID)

	return member, nil
}

// LeaveCircle removes a user from a savings circle
func (s *SavingsService) LeaveCircle(ctx context.Context, circleID, userID uuid.UUID) error {
	var circle models.SavingsCircle
	if err := s.db.First(&circle, circleID).Error; err != nil {
		return errors.New("circle not found")
	}

	// Cannot leave active circles that have started
	if circle.Status == "active" && circle.CurrentRound > 0 {
		return errors.New("cannot leave an active circle after contributions have started")
	}

	var member models.CircleMember
	err := s.db.Where("circle_id = ? AND user_id = ? AND status = ?", circleID, userID, "active").First(&member).Error
	if err != nil {
		return errors.New("you are not a member of this circle")
	}

	// Admin cannot leave (must delete or transfer)
	if member.Role == "admin" {
		return errors.New("admin cannot leave the circle. Please transfer ownership first or delete the circle")
	}

	tx := s.db.Begin()

	// Update member status
	tx.Model(&member).Update("status", "left")

	// Update member count
	tx.Model(&circle).UpdateColumn("current_members", gorm.Expr("current_members - ?", 1))

	// Reorder positions
	tx.Exec(`
		UPDATE circle_members 
		SET position = position - 1 
		WHERE circle_id = ? AND position > ? AND status = 'active'
	`, circleID, member.Position)

	tx.Commit()

	return nil
}

// MakeContribution records a contribution to the circle
func (s *SavingsService) MakeContribution(ctx context.Context, circleID, userID uuid.UUID, transactionID uuid.UUID) (*models.Contribution, error) {
	var circle models.SavingsCircle
	if err := s.db.First(&circle, circleID).Error; err != nil {
		return nil, errors.New("circle not found")
	}

	if circle.Status != "active" {
		return nil, errors.New("circle is not active")
	}

	var member models.CircleMember
	err := s.db.Where("circle_id = ? AND user_id = ? AND status = ?", circleID, userID, "active").First(&member).Error
	if err != nil {
		return nil, errors.New("you are not a member of this circle")
	}

	// Find pending contribution for current round
	var contribution models.Contribution
	err = s.db.Where(
		"circle_id = ? AND member_id = ? AND round = ? AND status = ?",
		circleID, member.ID, circle.CurrentRound, models.ContributionStatusPending,
	).First(&contribution).Error

	if err != nil {
		return nil, errors.New("no pending contribution found for this round")
	}

	now := time.Now()
	lateFee := int64(0)

	// Calculate late fee if past due date
	if now.After(contribution.DueDate) {
		// 5% late fee
		lateFee = circle.ContributionAmt / 20
	}

	tx := s.db.Begin()

	// Update contribution
	contribution.Status = models.ContributionStatusPaid
	contribution.PaidAt = &now
	contribution.TransactionID = &transactionID
	contribution.LateFee = lateFee

	if err := tx.Save(&contribution).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update member's total contribution
	tx.Model(&member).UpdateColumn("total_contrib", gorm.Expr("total_contrib + ?", circle.ContributionAmt+lateFee))

	// Update circle's pool balance
	tx.Model(&circle).UpdateColumn("pool_balance", gorm.Expr("pool_balance + ?", circle.ContributionAmt+lateFee))
	tx.Model(&circle).UpdateColumn("total_saved", gorm.Expr("total_saved + ?", circle.ContributionAmt+lateFee))

	// Check if all contributions for this round are complete
	var pendingCount int64
	tx.Model(&models.Contribution{}).Where(
		"circle_id = ? AND round = ? AND status = ?",
		circleID, circle.CurrentRound, models.ContributionStatusPending,
	).Count(&pendingCount)

	if pendingCount == 0 {
		// All contributions received, trigger payout
		s.triggerPayout(tx, &circle)
	}

	tx.Commit()

	return &contribution, nil
}

// triggerPayout distributes the pool to the current position holder
func (s *SavingsService) triggerPayout(tx *gorm.DB, circle *models.SavingsCircle) error {
	// Find the member for current round
	var recipient models.CircleMember
	err := tx.Where(
		"circle_id = ? AND position = ? AND status = ?",
		circle.ID, circle.CurrentRound, "active",
	).First(&recipient).Error

	if err != nil {
		return fmt.Errorf("recipient not found for position %d", circle.CurrentRound)
	}

	// Mark recipient as received
	tx.Model(&recipient).Update("has_received", true)

	// TODO: Create payout transaction through WalletService
	// Transfer circle.PoolBalance to recipient's wallet

	// Reset pool and advance round
	tx.Model(circle).Updates(map[string]interface{}{
		"pool_balance":  0,
		"current_round": circle.CurrentRound + 1,
	})

	// Calculate next payout date
	nextPayoutDate := calculateNextPayoutDate(circle.Frequency, time.Now())
	tx.Model(circle).Update("next_payout_date", nextPayoutDate)

	// Check if all rounds complete
	if circle.CurrentRound >= circle.TotalRounds {
		tx.Model(circle).Update("status", "completed")
	} else {
		// Schedule next round contributions
		s.scheduleContributions(tx, circle)
	}

	return nil
}

// scheduleContributions creates contribution records for the next round
func (s *SavingsService) scheduleContributions(tx *gorm.DB, circle *models.SavingsCircle) {
	round := circle.CurrentRound
	if round == 0 {
		round = 1
	}

	dueDate := calculateDueDate(circle.Frequency, time.Now())

	// Get active members
	var members []models.CircleMember
	tx.Where("circle_id = ? AND status = ?", circle.ID, "active").Find(&members)

	for _, member := range members {
		contribution := &models.Contribution{
			CircleID: circle.ID,
			MemberID: member.ID,
			Round:    round,
			Amount:   circle.ContributionAmt,
			DueDate:  dueDate,
			Status:   models.ContributionStatusPending,
		}
		tx.Create(contribution)
	}
}

// GetContributions gets contributions for a circle member
func (s *SavingsService) GetContributions(ctx context.Context, circleID, userID uuid.UUID) ([]models.Contribution, error) {
	var member models.CircleMember
	err := s.db.Where("circle_id = ? AND user_id = ?", circleID, userID).First(&member).Error
	if err != nil {
		return nil, errors.New("you are not a member of this circle")
	}

	var contributions []models.Contribution
	err = s.db.Where("circle_id = ? AND member_id = ?", circleID, member.ID).
		Order("round ASC").
		Find(&contributions).Error

	return contributions, err
}

// GetCircleLeaderboard gets contribution statistics for all members
func (s *SavingsService) GetCircleLeaderboard(ctx context.Context, circleID uuid.UUID) ([]models.CircleMember, error) {
	var members []models.CircleMember
	err := s.db.
		Preload("User").
		Where("circle_id = ? AND status = ?", circleID, "active").
		Order("total_contrib DESC, missed_payments ASC").
		Find(&members).Error

	return members, err
}

// StartCircle manually starts a circle (admin only)
func (s *SavingsService) StartCircle(ctx context.Context, circleID, adminID uuid.UUID) error {
	var circle models.SavingsCircle
	if err := s.db.First(&circle, circleID).Error; err != nil {
		return errors.New("circle not found")
	}

	if circle.CreatedBy != adminID {
		return errors.New("only the admin can start the circle")
	}

	if circle.Status != "recruiting" {
		return errors.New("circle has already started")
	}

	if circle.CurrentMembers < 2 {
		return errors.New("need at least 2 members to start")
	}

	tx := s.db.Begin()

	now := time.Now()
	tx.Model(&circle).Updates(map[string]interface{}{
		"status":       "active",
		"start_date":   now,
		"current_round": 1,
	})

	// Refresh circle
	tx.First(&circle, circleID)

	// Schedule contributions
	s.scheduleContributions(tx, &circle)

	tx.Commit()

	return nil
}

// Helper functions

func generateInviteCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Removed confusing chars
	result := make([]byte, 8)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func calculateDueDate(frequency string, from time.Time) time.Time {
	switch frequency {
	case "daily":
		return from.AddDate(0, 0, 1)
	case "weekly":
		return from.AddDate(0, 0, 7)
	case "biweekly":
		return from.AddDate(0, 0, 14)
	case "monthly":
		return from.AddDate(0, 1, 0)
	default:
		return from.AddDate(0, 0, 7)
	}
}

func calculateNextPayoutDate(frequency string, from time.Time) time.Time {
	return calculateDueDate(frequency, from)
}
