package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"hustlex/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GigService handles gig-related operations
type GigService struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewGigService creates a new gig service
func NewGigService(db *gorm.DB, redis *redis.Client) *GigService {
	return &GigService{
		db:    db,
		redis: redis,
	}
}

// CreateGigInput represents input for creating a gig
type CreateGigInput struct {
	Title        string   `json:"title" validate:"required,min=10,max=200"`
	Description  string   `json:"description" validate:"required,min=50,max=5000"`
	Category     string   `json:"category" validate:"required"`
	SkillID      string   `json:"skill_id" validate:"omitempty,uuid"`
	BudgetMin    int64    `json:"budget_min" validate:"required,min=1000"` // min ₦10 (1000 kobo)
	BudgetMax    int64    `json:"budget_max" validate:"required,gtefield=BudgetMin"`
	DeliveryDays int      `json:"delivery_days" validate:"required,min=1,max=90"`
	Deadline     string   `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
	IsRemote     bool     `json:"is_remote"`
	Location     string   `json:"location" validate:"omitempty,max=100"`
	Tags         []string `json:"tags" validate:"omitempty,max=10,dive,max=30"`
	Attachments  []string `json:"attachments" validate:"omitempty,max=5,dive,url"`
}

// GigFilters represents filters for listing gigs
type GigFilters struct {
	Category    string   `query:"category"`
	SkillID     string   `query:"skill_id"`
	MinBudget   int64    `query:"min_budget"`
	MaxBudget   int64    `query:"max_budget"`
	IsRemote    *bool    `query:"is_remote"`
	Location    string   `query:"location"`
	Status      string   `query:"status"`
	Search      string   `query:"search"`
	SortBy      string   `query:"sort_by"` // newest, budget_high, budget_low, deadline
	Page        int      `query:"page"`
	Limit       int      `query:"limit"`
	ExcludeOwn  bool     `query:"exclude_own"`
}

// GigListResult represents paginated gig results
type GigListResult struct {
	Gigs       []models.Gig `json:"gigs"`
	Total      int64        `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

// CreateGig creates a new gig posting
func (s *GigService) CreateGig(ctx context.Context, clientID uuid.UUID, input *CreateGigInput) (*models.Gig, error) {
	// Validate budget
	if input.BudgetMin > input.BudgetMax {
		return nil, errors.New("minimum budget cannot exceed maximum budget")
	}

	gig := &models.Gig{
		ClientID:     clientID,
		Title:        input.Title,
		Description:  input.Description,
		Category:     input.Category,
		BudgetMin:    input.BudgetMin,
		BudgetMax:    input.BudgetMax,
		DeliveryDays: input.DeliveryDays,
		IsRemote:     input.IsRemote,
		Location:     input.Location,
		Status:       models.GigStatusOpen,
		Tags:         input.Tags,
		Attachments:  input.Attachments,
		Currency:     "NGN",
	}

	// Parse skill ID if provided
	if input.SkillID != "" {
		skillUUID, err := uuid.Parse(input.SkillID)
		if err == nil {
			gig.SkillID = &skillUUID
		}
	}

	// Parse deadline if provided
	if input.Deadline != "" {
		deadline, err := time.Parse("2006-01-02", input.Deadline)
		if err == nil {
			gig.Deadline = &deadline
		}
	}

	if err := s.db.Create(gig).Error; err != nil {
		return nil, fmt.Errorf("failed to create gig: %w", err)
	}

	// Load relationships
	s.db.Preload("Client").Preload("Skill").First(gig, gig.ID)

	// Invalidate cache
	s.invalidateGigCache(ctx)

	return gig, nil
}

// GetGig retrieves a single gig by ID
func (s *GigService) GetGig(ctx context.Context, gigID uuid.UUID) (*models.Gig, error) {
	var gig models.Gig
	
	err := s.db.
		Preload("Client").
		Preload("Skill").
		Preload("Proposals", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Proposals.Hustler").
		Preload("Contract").
		Preload("Contract.Hustler").
		First(&gig, gigID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gig not found")
		}
		return nil, err
	}

	// Increment view count
	s.db.Model(&gig).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	return &gig, nil
}

// ListGigs lists gigs with filters and pagination
func (s *GigService) ListGigs(ctx context.Context, userID *uuid.UUID, filters *GigFilters) (*GigListResult, error) {
	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 || filters.Limit > 50 {
		filters.Limit = 20
	}

	query := s.db.Model(&models.Gig{}).Preload("Client").Preload("Skill")

	// Apply filters
	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}

	if filters.SkillID != "" {
		skillUUID, _ := uuid.Parse(filters.SkillID)
		query = query.Where("skill_id = ?", skillUUID)
	}

	if filters.MinBudget > 0 {
		query = query.Where("budget_max >= ?", filters.MinBudget)
	}

	if filters.MaxBudget > 0 {
		query = query.Where("budget_min <= ?", filters.MaxBudget)
	}

	if filters.IsRemote != nil {
		query = query.Where("is_remote = ?", *filters.IsRemote)
	}

	if filters.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filters.Location+"%")
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	} else {
		// Default to open gigs
		query = query.Where("status = ?", models.GigStatusOpen)
	}

	if filters.Search != "" {
		query = query.Where(
			"to_tsvector('english', title || ' ' || description) @@ plainto_tsquery('english', ?)",
			filters.Search,
		)
	}

	// Exclude own gigs if requested
	if filters.ExcludeOwn && userID != nil {
		query = query.Where("client_id != ?", *userID)
	}

	// Count total
	var total int64
	query.Count(&total)

	// Apply sorting
	switch filters.SortBy {
	case "budget_high":
		query = query.Order("budget_max DESC")
	case "budget_low":
		query = query.Order("budget_min ASC")
	case "deadline":
		query = query.Order("deadline ASC NULLS LAST")
	case "popular":
		query = query.Order("view_count DESC")
	default: // newest
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)

	var gigs []models.Gig
	if err := query.Find(&gigs).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / filters.Limit
	if int(total)%filters.Limit > 0 {
		totalPages++
	}

	return &GigListResult{
		Gigs:       gigs,
		Total:      total,
		Page:       filters.Page,
		Limit:      filters.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateGig updates a gig
func (s *GigService) UpdateGig(ctx context.Context, gigID, clientID uuid.UUID, input *CreateGigInput) (*models.Gig, error) {
	var gig models.Gig
	if err := s.db.First(&gig, gigID).Error; err != nil {
		return nil, errors.New("gig not found")
	}

	// Verify ownership
	if gig.ClientID != clientID {
		return nil, errors.New("you can only update your own gigs")
	}

	// Cannot update if not open
	if gig.Status != models.GigStatusOpen {
		return nil, errors.New("cannot update gig that is already in progress")
	}

	// Update fields
	gig.Title = input.Title
	gig.Description = input.Description
	gig.Category = input.Category
	gig.BudgetMin = input.BudgetMin
	gig.BudgetMax = input.BudgetMax
	gig.DeliveryDays = input.DeliveryDays
	gig.IsRemote = input.IsRemote
	gig.Location = input.Location
	gig.Tags = input.Tags
	gig.Attachments = input.Attachments

	if input.SkillID != "" {
		skillUUID, _ := uuid.Parse(input.SkillID)
		gig.SkillID = &skillUUID
	}

	if input.Deadline != "" {
		deadline, _ := time.Parse("2006-01-02", input.Deadline)
		gig.Deadline = &deadline
	}

	if err := s.db.Save(&gig).Error; err != nil {
		return nil, err
	}

	s.invalidateGigCache(ctx)
	return &gig, nil
}

// DeleteGig deletes/cancels a gig
func (s *GigService) DeleteGig(ctx context.Context, gigID, clientID uuid.UUID) error {
	var gig models.Gig
	if err := s.db.First(&gig, gigID).Error; err != nil {
		return errors.New("gig not found")
	}

	if gig.ClientID != clientID {
		return errors.New("you can only delete your own gigs")
	}

	if gig.Status != models.GigStatusOpen {
		return errors.New("cannot delete gig that is already in progress")
	}

	gig.Status = models.GigStatusCancelled
	if err := s.db.Save(&gig).Error; err != nil {
		return err
	}

	s.invalidateGigCache(ctx)
	return nil
}

// SubmitProposalInput represents proposal submission input
type SubmitProposalInput struct {
	CoverLetter   string   `json:"cover_letter" validate:"required,min=50,max=2000"`
	ProposedPrice int64    `json:"proposed_price" validate:"required,min=1000"`
	DeliveryDays  int      `json:"delivery_days" validate:"required,min=1,max=90"`
	Attachments   []string `json:"attachments" validate:"omitempty,max=5,dive,url"`
}

// SubmitProposal submits a proposal for a gig
func (s *GigService) SubmitProposal(ctx context.Context, gigID, hustlerID uuid.UUID, input *SubmitProposalInput) (*models.GigProposal, error) {
	// Get gig
	var gig models.Gig
	if err := s.db.First(&gig, gigID).Error; err != nil {
		return nil, errors.New("gig not found")
	}

	// Validate gig status
	if gig.Status != models.GigStatusOpen {
		return nil, errors.New("this gig is no longer accepting proposals")
	}

	// Cannot propose to own gig
	if gig.ClientID == hustlerID {
		return nil, errors.New("you cannot submit a proposal to your own gig")
	}

	// Check for existing proposal
	var existingProposal models.GigProposal
	err := s.db.Where("gig_id = ? AND hustler_id = ?", gigID, hustlerID).First(&existingProposal).Error
	if err == nil {
		return nil, errors.New("you have already submitted a proposal for this gig")
	}

	// Validate proposed price is within budget
	if input.ProposedPrice < gig.BudgetMin || input.ProposedPrice > gig.BudgetMax {
		return nil, fmt.Errorf("proposed price must be between ₦%d and ₦%d", gig.BudgetMin/100, gig.BudgetMax/100)
	}

	proposal := &models.GigProposal{
		GigID:         gigID,
		HustlerID:     hustlerID,
		CoverLetter:   input.CoverLetter,
		ProposedPrice: input.ProposedPrice,
		DeliveryDays:  input.DeliveryDays,
		Status:        "pending",
		Attachments:   input.Attachments,
	}

	tx := s.db.Begin()

	if err := tx.Create(proposal).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to submit proposal: %w", err)
	}

	// Update proposal count on gig
	tx.Model(&gig).UpdateColumn("proposal_count", gorm.Expr("proposal_count + ?", 1))

	tx.Commit()

	// Load relationships
	s.db.Preload("Hustler").First(proposal, proposal.ID)

	return proposal, nil
}

// GetProposals gets proposals for a gig (client only)
func (s *GigService) GetProposals(ctx context.Context, gigID, clientID uuid.UUID) ([]models.GigProposal, error) {
	var gig models.Gig
	if err := s.db.First(&gig, gigID).Error; err != nil {
		return nil, errors.New("gig not found")
	}

	if gig.ClientID != clientID {
		return nil, errors.New("you can only view proposals for your own gigs")
	}

	var proposals []models.GigProposal
	err := s.db.
		Preload("Hustler").
		Preload("Hustler.Skills").
		Preload("Hustler.Skills.Skill").
		Preload("Hustler.CreditScore").
		Where("gig_id = ?", gigID).
		Order("created_at DESC").
		Find(&proposals).Error

	return proposals, err
}

// AcceptProposal accepts a proposal and creates a contract
func (s *GigService) AcceptProposal(ctx context.Context, gigID, proposalID, clientID uuid.UUID) (*models.GigContract, error) {
	var gig models.Gig
	if err := s.db.First(&gig, gigID).Error; err != nil {
		return nil, errors.New("gig not found")
	}

	if gig.ClientID != clientID {
		return nil, errors.New("you can only accept proposals for your own gigs")
	}

	if gig.Status != models.GigStatusOpen {
		return nil, errors.New("this gig is no longer accepting proposals")
	}

	var proposal models.GigProposal
	if err := s.db.First(&proposal, proposalID).Error; err != nil {
		return nil, errors.New("proposal not found")
	}

	if proposal.GigID != gigID {
		return nil, errors.New("proposal does not belong to this gig")
	}

	if proposal.Status != "pending" {
		return nil, errors.New("this proposal is no longer available")
	}

	// Calculate platform fee (10%)
	platformFee := proposal.ProposedPrice / 10

	// Calculate deadline
	deadlineAt := time.Now().AddDate(0, 0, proposal.DeliveryDays)

	tx := s.db.Begin()

	// Create contract
	contract := &models.GigContract{
		GigID:        gigID,
		HustlerID:    proposal.HustlerID,
		ProposalID:   proposal.ID,
		AgreedPrice:  proposal.ProposedPrice,
		PlatformFee:  platformFee,
		DeliveryDays: proposal.DeliveryDays,
		Status:       models.ContractStatusActive,
		DeadlineAt:   deadlineAt,
	}

	if err := tx.Create(contract).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	// Update gig status
	tx.Model(&gig).Update("status", models.GigStatusInProgress)

	// Update accepted proposal status
	tx.Model(&proposal).Update("status", "accepted")

	// Reject other proposals
	tx.Model(&models.GigProposal{}).
		Where("gig_id = ? AND id != ? AND status = ?", gigID, proposalID, "pending").
		Update("status", "rejected")

	// TODO: Create escrow transaction (hold client's funds)
	// This would involve the WalletService

	tx.Commit()

	// Load relationships
	s.db.Preload("Gig").Preload("Hustler").Preload("Proposal").First(contract, contract.ID)

	return contract, nil
}

// DeliverContract marks a contract as delivered
func (s *GigService) DeliverContract(ctx context.Context, contractID, hustlerID uuid.UUID, deliverables []string) (*models.GigContract, error) {
	var contract models.GigContract
	if err := s.db.Preload("Gig").First(&contract, contractID).Error; err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.HustlerID != hustlerID {
		return nil, errors.New("you can only deliver your own contracts")
	}

	if contract.Status != models.ContractStatusActive {
		return nil, errors.New("this contract cannot be delivered")
	}

	now := time.Now()
	contract.Status = models.ContractStatusDelivered
	contract.DeliveredAt = &now
	contract.Deliverables = deliverables

	if err := s.db.Save(&contract).Error; err != nil {
		return nil, err
	}

	// TODO: Send notification to client

	return &contract, nil
}

// ApproveContract approves delivery and releases payment
func (s *GigService) ApproveContract(ctx context.Context, contractID, clientID uuid.UUID, notes string) (*models.GigContract, error) {
	var contract models.GigContract
	if err := s.db.Preload("Gig").First(&contract, contractID).Error; err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.Gig.ClientID != clientID {
		return nil, errors.New("you can only approve your own contracts")
	}

	if contract.Status != models.ContractStatusDelivered {
		return nil, errors.New("this contract has not been delivered yet")
	}

	now := time.Now()
	contract.Status = models.ContractStatusCompleted
	contract.CompletedAt = &now
	contract.ClientNotes = notes

	tx := s.db.Begin()

	if err := tx.Save(&contract).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Update gig status
	tx.Model(&models.Gig{}).Where("id = ?", contract.GigID).Update("status", models.GigStatusCompleted)

	// TODO: Release payment from escrow to hustler
	// This would involve the WalletService

	// TODO: Update credit scores for both parties

	tx.Commit()

	return &contract, nil
}

// SubmitReview submits a review for a completed contract
func (s *GigService) SubmitReview(ctx context.Context, contractID, reviewerID uuid.UUID, rating int, reviewText string, detailedRatings map[string]int) (*models.GigReview, error) {
	var contract models.GigContract
	if err := s.db.Preload("Gig").First(&contract, contractID).Error; err != nil {
		return nil, errors.New("contract not found")
	}

	if contract.Status != models.ContractStatusCompleted {
		return nil, errors.New("can only review completed contracts")
	}

	// Determine reviewer and reviewee
	var revieweeID uuid.UUID
	if contract.Gig.ClientID == reviewerID {
		// Client reviewing hustler
		revieweeID = contract.HustlerID
	} else if contract.HustlerID == reviewerID {
		// Hustler reviewing client
		revieweeID = contract.Gig.ClientID
	} else {
		return nil, errors.New("you are not part of this contract")
	}

	// Check for existing review
	var existing models.GigReview
	err := s.db.Where("contract_id = ? AND reviewer_id = ?", contractID, reviewerID).First(&existing).Error
	if err == nil {
		return nil, errors.New("you have already reviewed this contract")
	}

	// Validate rating
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	review := &models.GigReview{
		ContractID: contractID,
		ReviewerID: reviewerID,
		RevieweeID: revieweeID,
		Rating:     rating,
		ReviewText: reviewText,
		IsPublic:   true,
	}

	// Set detailed ratings if provided
	if commRating, ok := detailedRatings["communication"]; ok && commRating >= 1 && commRating <= 5 {
		review.CommunicationRating = commRating
	}
	if qualRating, ok := detailedRatings["quality"]; ok && qualRating >= 1 && qualRating <= 5 {
		review.QualityRating = qualRating
	}
	if timeRating, ok := detailedRatings["timeliness"]; ok && timeRating >= 1 && timeRating <= 5 {
		review.TimelinessRating = timeRating
	}

	tx := s.db.Begin()

	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to submit review: %w", err)
	}

	// Update reviewee's credit score
	tx.Exec(`
		UPDATE credit_scores 
		SET 
			total_reviews = total_reviews + 1,
			average_rating = (average_rating * total_reviews + ?) / (total_reviews + 1),
			updated_at = NOW()
		WHERE user_id = ?
	`, rating, revieweeID)

	tx.Commit()

	// Load relationships
	s.db.Preload("Reviewer").Preload("Reviewee").First(review, review.ID)

	return review, nil
}

// GetMyContracts gets contracts for a user (as client or hustler)
func (s *GigService) GetMyContracts(ctx context.Context, userID uuid.UUID, role string, status string) ([]models.GigContract, error) {
	query := s.db.
		Preload("Gig").
		Preload("Gig.Client").
		Preload("Hustler").
		Preload("Review")

	if role == "client" {
		query = query.Joins("JOIN gigs ON gigs.id = gig_contracts.gig_id").
			Where("gigs.client_id = ?", userID)
	} else {
		query = query.Where("hustler_id = ?", userID)
	}

	if status != "" {
		query = query.Where("gig_contracts.status = ?", status)
	}

	var contracts []models.GigContract
	err := query.Order("gig_contracts.created_at DESC").Find(&contracts).Error

	return contracts, err
}

func (s *GigService) invalidateGigCache(ctx context.Context) {
	s.redis.Del(ctx, "gigs:list:*")
}
