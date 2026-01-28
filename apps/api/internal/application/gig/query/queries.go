package query

import (
	"context"
	"time"

	"hustlex/internal/domain/gig/aggregate"
	"hustlex/internal/domain/gig/repository"
	"hustlex/internal/domain/gig/service"
	"hustlex/internal/domain/shared/valueobject"
)

// GetGig retrieves a single gig
type GetGig struct {
	GigID string
}

// GetGigs retrieves gigs with filters
type GetGigs struct {
	Category      string
	SkillID       string
	MinBudget     int64
	MaxBudget     int64
	IsRemote      *bool
	Location      string
	Status        string
	SearchQuery   string
	ExcludeUserID string
	SortBy        string // newest, budget_high, budget_low, deadline, popular
	Page          int
	Limit         int
}

// GetMyGigs retrieves gigs posted by a user
type GetMyGigs struct {
	ClientID string
	Status   string
	Page     int
	Limit    int
}

// GigDTO represents a gig for API responses
type GigDTO struct {
	ID            string       `json:"id"`
	ClientID      string       `json:"client_id"`
	ClientName    string       `json:"client_name,omitempty"`
	Title         string       `json:"title"`
	Description   string       `json:"description"`
	Category      string       `json:"category"`
	SkillID       string       `json:"skill_id,omitempty"`
	SkillName     string       `json:"skill_name,omitempty"`
	BudgetMin     int64        `json:"budget_min"`
	BudgetMax     int64        `json:"budget_max"`
	Currency      string       `json:"currency"`
	DeliveryDays  int          `json:"delivery_days"`
	Deadline      *time.Time   `json:"deadline,omitempty"`
	IsRemote      bool         `json:"is_remote"`
	Location      string       `json:"location,omitempty"`
	Status        string       `json:"status"`
	ViewCount     int          `json:"view_count"`
	ProposalCount int          `json:"proposal_count"`
	IsFeatured    bool         `json:"is_featured"`
	Tags          []string     `json:"tags,omitempty"`
	Attachments   []string     `json:"attachments,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// GigListResult represents paginated gig results
type GigListResult struct {
	Gigs       []GigDTO `json:"gigs"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalPages int      `json:"total_pages"`
}

// GetGigProposals retrieves proposals for a gig
type GetGigProposals struct {
	GigID    string
	ClientID string // must be gig owner
}

// ProposalDTO represents a proposal for API responses
type ProposalDTO struct {
	ID              string    `json:"id"`
	GigID           string    `json:"gig_id"`
	HustlerID       string    `json:"hustler_id"`
	HustlerName     string    `json:"hustler_name,omitempty"`
	HustlerImage    string    `json:"hustler_image,omitempty"`
	HustlerRating   float64   `json:"hustler_rating,omitempty"`
	CoverLetter     string    `json:"cover_letter"`
	ProposedPrice   int64     `json:"proposed_price"`
	Currency        string    `json:"currency"`
	DeliveryDays    int       `json:"delivery_days"`
	Status          string    `json:"status"`
	Attachments     []string  `json:"attachments,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// GetMyProposals retrieves proposals submitted by a user
type GetMyProposals struct {
	HustlerID string
	Status    string
	Page      int
	Limit     int
}

// GetContract retrieves a contract
type GetContract struct {
	ContractID string
	UserID     string // must be party to contract
}

// GetMyContracts retrieves contracts for a user
type GetMyContracts struct {
	UserID string
	Role   string // client, hustler
	Status string
	Page   int
	Limit  int
}

// ContractDTO represents a contract for API responses
type ContractDTO struct {
	ID           string     `json:"id"`
	GigID        string     `json:"gig_id"`
	GigTitle     string     `json:"gig_title"`
	ClientID     string     `json:"client_id"`
	ClientName   string     `json:"client_name,omitempty"`
	HustlerID    string     `json:"hustler_id"`
	HustlerName  string     `json:"hustler_name,omitempty"`
	AgreedPrice  int64      `json:"agreed_price"`
	PlatformFee  int64      `json:"platform_fee"`
	NetPayout    int64      `json:"net_payout"`
	Currency     string     `json:"currency"`
	Status       string     `json:"status"`
	DeliveryDays int        `json:"delivery_days"`
	StartedAt    time.Time  `json:"started_at"`
	DeadlineAt   time.Time  `json:"deadline_at"`
	DeliveredAt  *time.Time `json:"delivered_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Deliverables []string   `json:"deliverables,omitempty"`
	HasReviewed  bool       `json:"has_reviewed"`
	IsOverdue    bool       `json:"is_overdue"`
}

// ContractListResult represents paginated contract results
type ContractListResult struct {
	Contracts  []ContractDTO `json:"contracts"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

// GetReviews retrieves reviews for a user
type GetReviews struct {
	UserID string
	Page   int
	Limit  int
}

// ReviewDTO represents a review for API responses
type ReviewDTO struct {
	ID                  string    `json:"id"`
	ContractID          string    `json:"contract_id"`
	GigTitle            string    `json:"gig_title,omitempty"`
	ReviewerID          string    `json:"reviewer_id"`
	ReviewerName        string    `json:"reviewer_name,omitempty"`
	ReviewerImage       string    `json:"reviewer_image,omitempty"`
	Rating              int       `json:"rating"`
	ReviewText          string    `json:"review_text,omitempty"`
	CommunicationRating int       `json:"communication_rating,omitempty"`
	QualityRating       int       `json:"quality_rating,omitempty"`
	TimelinessRating    int       `json:"timeliness_rating,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

// ReviewListResult represents paginated review results
type ReviewListResult struct {
	Reviews       []ReviewDTO `json:"reviews"`
	Total         int64       `json:"total"`
	AverageRating float64     `json:"average_rating"`
	Page          int         `json:"page"`
	Limit         int         `json:"limit"`
	TotalPages    int         `json:"total_pages"`
}

// GigQueryHandler handles gig queries
type GigQueryHandler struct {
	gigRepo      repository.GigRepository
	proposalRepo repository.ProposalRepository
	contractRepo repository.ContractRepository
	reviewRepo   repository.ReviewRepository
}

// NewGigQueryHandler creates a new query handler
func NewGigQueryHandler(
	gigRepo repository.GigRepository,
	proposalRepo repository.ProposalRepository,
	contractRepo repository.ContractRepository,
	reviewRepo repository.ReviewRepository,
) *GigQueryHandler {
	return &GigQueryHandler{
		gigRepo:      gigRepo,
		proposalRepo: proposalRepo,
		contractRepo: contractRepo,
		reviewRepo:   reviewRepo,
	}
}

// HandleGetGig retrieves a single gig
func (h *GigQueryHandler) HandleGetGig(ctx context.Context, q GetGig) (*GigDTO, error) {
	gigID, err := valueobject.NewGigID(q.GigID)
	if err != nil {
		return nil, err
	}

	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return nil, service.ErrGigNotFound
	}

	// Increment view count
	gig.IncrementViewCount()
	_ = h.gigRepo.Save(ctx, gig)

	return gigToDTO(gig), nil
}

// HandleGetGigs retrieves gigs with filters
func (h *GigQueryHandler) HandleGetGigs(ctx context.Context, q GetGigs) (*GigListResult, error) {
	// Set defaults
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	filter := repository.GigFilter{
		Category:    q.Category,
		MinBudget:   q.MinBudget,
		MaxBudget:   q.MaxBudget,
		IsRemote:    q.IsRemote,
		Location:    q.Location,
		SearchQuery: q.SearchQuery,
		SortBy:      q.SortBy,
		Offset:      (q.Page - 1) * q.Limit,
		Limit:       q.Limit,
	}

	if q.SkillID != "" {
		skillID, err := valueobject.NewSkillID(q.SkillID)
		if err == nil {
			filter.SkillID = &skillID
		}
	}

	if q.Status != "" {
		status := aggregate.GigStatus(q.Status)
		filter.Status = &status
	}

	if q.ExcludeUserID != "" {
		userID, err := valueobject.NewUserID(q.ExcludeUserID)
		if err == nil {
			filter.ExcludeUserID = &userID
		}
	}

	gigs, total, err := h.gigRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]GigDTO, len(gigs))
	for i, gig := range gigs {
		dtos[i] = *gigToDTO(gig)
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &GigListResult{
		Gigs:       dtos,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// HandleGetMyGigs retrieves gigs posted by a user
func (h *GigQueryHandler) HandleGetMyGigs(ctx context.Context, q GetMyGigs) (*GigListResult, error) {
	clientID, err := valueobject.NewUserID(q.ClientID)
	if err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	filter := repository.GigFilter{
		Offset: (q.Page - 1) * q.Limit,
		Limit:  q.Limit,
	}

	if q.Status != "" {
		status := aggregate.GigStatus(q.Status)
		filter.Status = &status
	}

	gigs, total, err := h.gigRepo.FindByClientID(ctx, clientID, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]GigDTO, len(gigs))
	for i, gig := range gigs {
		dtos[i] = *gigToDTO(gig)
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &GigListResult{
		Gigs:       dtos,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// HandleGetGigProposals retrieves proposals for a gig
func (h *GigQueryHandler) HandleGetGigProposals(ctx context.Context, q GetGigProposals) ([]ProposalDTO, error) {
	gigID, err := valueobject.NewGigID(q.GigID)
	if err != nil {
		return nil, err
	}

	clientID, err := valueobject.NewUserID(q.ClientID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return nil, service.ErrGigNotFound
	}

	if !gig.ClientID().Equals(clientID) {
		return nil, service.ErrUnauthorized
	}

	proposals := gig.Proposals()
	dtos := make([]ProposalDTO, len(proposals))
	for i, p := range proposals {
		dtos[i] = ProposalDTO{
			ID:            p.ID().String(),
			GigID:         gigID.String(),
			HustlerID:     p.HustlerID().String(),
			CoverLetter:   p.CoverLetter(),
			ProposedPrice: p.ProposedPrice().Amount(),
			Currency:      string(p.ProposedPrice().Currency()),
			DeliveryDays:  p.DeliveryDays(),
			Status:        string(p.Status()),
			Attachments:   p.Attachments(),
			CreatedAt:     p.CreatedAt(),
		}
	}

	return dtos, nil
}

// HandleGetContract retrieves a contract
func (h *GigQueryHandler) HandleGetContract(ctx context.Context, q GetContract) (*ContractDTO, error) {
	contractID, err := valueobject.NewContractID(q.ContractID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	contract, err := h.contractRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, service.ErrContractNotFound
	}

	// Verify user is party to contract
	if !contract.IsParty(userID) {
		return nil, service.ErrUnauthorized
	}

	return contractToDTO(contract, userID), nil
}

// HandleGetMyContracts retrieves contracts for a user
func (h *GigQueryHandler) HandleGetMyContracts(ctx context.Context, q GetMyContracts) (*ContractListResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	var status *aggregate.ContractStatus
	if q.Status != "" {
		s := aggregate.ContractStatus(q.Status)
		status = &s
	}

	contracts, total, err := h.contractRepo.FindByUserID(
		ctx,
		userID,
		q.Role,
		status,
		(q.Page-1)*q.Limit,
		q.Limit,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]ContractDTO, len(contracts))
	for i, c := range contracts {
		dtos[i] = ContractDTO{
			ID:           c.ID,
			GigID:        c.GigID,
			GigTitle:     c.GigTitle,
			ClientID:     c.ClientID,
			ClientName:   c.ClientName,
			HustlerID:    c.HustlerID,
			HustlerName:  c.HustlerName,
			AgreedPrice:  c.AgreedPrice,
			PlatformFee:  c.PlatformFee,
			NetPayout:    c.AgreedPrice - c.PlatformFee,
			Currency:     c.Currency,
			Status:       c.Status,
			DeliveryDays: c.DeliveryDays,
			StartedAt:    c.StartedAt,
			DeadlineAt:   c.DeadlineAt,
			DeliveredAt:  c.DeliveredAt,
			CompletedAt:  c.CompletedAt,
			HasReviewed:  c.HasReviewed,
		}
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &ContractListResult{
		Contracts:  dtos,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// HandleGetReviews retrieves reviews for a user
func (h *GigQueryHandler) HandleGetReviews(ctx context.Context, q GetReviews) (*ReviewListResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	reviews, total, err := h.reviewRepo.FindByRevieweeID(ctx, userID, (q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		return nil, err
	}

	avgRating, _, err := h.reviewRepo.GetAverageRating(ctx, userID)
	if err != nil {
		avgRating = 0
	}

	dtos := make([]ReviewDTO, len(reviews))
	for i, r := range reviews {
		dtos[i] = ReviewDTO{
			ID:                  r.ID,
			ContractID:          r.ContractID,
			GigTitle:            r.GigTitle,
			ReviewerID:          r.ReviewerID,
			ReviewerName:        r.ReviewerName,
			ReviewerImage:       r.ReviewerProfileImage,
			Rating:              r.Rating,
			ReviewText:          r.ReviewText,
			CommunicationRating: r.CommunicationRating,
			QualityRating:       r.QualityRating,
			TimelinessRating:    r.TimelinessRating,
			CreatedAt:           r.CreatedAt,
		}
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &ReviewListResult{
		Reviews:       dtos,
		Total:         total,
		AverageRating: avgRating,
		Page:          q.Page,
		Limit:         q.Limit,
		TotalPages:    totalPages,
	}, nil
}

// Helper functions

func gigToDTO(gig *aggregate.Gig) *GigDTO {
	dto := &GigDTO{
		ID:            gig.ID().String(),
		ClientID:      gig.ClientID().String(),
		Title:         gig.Title(),
		Description:   gig.Description(),
		Category:      gig.Category(),
		BudgetMin:     gig.Budget().Min().Amount(),
		BudgetMax:     gig.Budget().Max().Amount(),
		Currency:      string(gig.Currency()),
		DeliveryDays:  gig.DeliveryDays(),
		Deadline:      gig.Deadline(),
		IsRemote:      gig.IsRemote(),
		Location:      gig.Location(),
		Status:        gig.Status().String(),
		ViewCount:     gig.ViewCount(),
		ProposalCount: gig.ProposalCount(),
		IsFeatured:    gig.IsFeatured(),
		Tags:          gig.Tags(),
		Attachments:   gig.Attachments(),
		CreatedAt:     gig.CreatedAt(),
		UpdatedAt:     gig.UpdatedAt(),
	}

	if gig.SkillID() != nil {
		dto.SkillID = gig.SkillID().String()
	}

	return dto
}

func contractToDTO(contract *aggregate.Contract, viewerID valueobject.UserID) *ContractDTO {
	return &ContractDTO{
		ID:           contract.ID().String(),
		GigID:        contract.GigID().String(),
		ClientID:     contract.ClientID().String(),
		HustlerID:    contract.HustlerID().String(),
		AgreedPrice:  contract.AgreedPrice().Amount(),
		PlatformFee:  contract.PlatformFee().Amount(),
		NetPayout:    contract.NetPayoutAmount().Amount(),
		Currency:     string(contract.AgreedPrice().Currency()),
		Status:       contract.Status().String(),
		DeliveryDays: contract.DeliveryDays(),
		StartedAt:    contract.StartedAt(),
		DeadlineAt:   contract.DeadlineAt(),
		DeliveredAt:  contract.DeliveredAt(),
		CompletedAt:  contract.CompletedAt(),
		Deliverables: contract.Deliverables(),
		HasReviewed:  contract.HasReviewFrom(viewerID),
		IsOverdue:    contract.IsOverdue(),
	}
}
