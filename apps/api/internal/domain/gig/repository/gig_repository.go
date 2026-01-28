package repository

import (
	"context"
	"time"

	"hustlex/internal/domain/gig/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// GigRepository defines the interface for gig persistence
type GigRepository interface {
	// Save persists a gig aggregate
	Save(ctx context.Context, gig *aggregate.Gig) error

	// SaveWithEvents persists a gig and publishes domain events
	SaveWithEvents(ctx context.Context, gig *aggregate.Gig) error

	// FindByID retrieves a gig by ID
	FindByID(ctx context.Context, id valueobject.GigID) (*aggregate.Gig, error)

	// FindByClientID retrieves gigs posted by a client
	FindByClientID(ctx context.Context, clientID valueobject.UserID, filter GigFilter) ([]*aggregate.Gig, int64, error)

	// List retrieves gigs with filters
	List(ctx context.Context, filter GigFilter) ([]*aggregate.Gig, int64, error)

	// Delete soft-deletes a gig
	Delete(ctx context.Context, id valueobject.GigID) error
}

// GigFilter contains filter options for listing gigs
type GigFilter struct {
	Category    string
	SkillID     *valueobject.SkillID
	MinBudget   int64
	MaxBudget   int64
	IsRemote    *bool
	Location    string
	Status      *aggregate.GigStatus
	SearchQuery string
	ExcludeUserID *valueobject.UserID
	SortBy      string // newest, budget_high, budget_low, deadline, popular
	Offset      int
	Limit       int
}

// ProposalRepository defines the interface for proposal persistence
type ProposalRepository interface {
	// Save persists a proposal
	Save(ctx context.Context, gigID valueobject.GigID, proposal *aggregate.Proposal) error

	// FindByID retrieves a proposal by ID
	FindByID(ctx context.Context, id valueobject.ProposalID) (*aggregate.Proposal, error)

	// FindByGigID retrieves all proposals for a gig
	FindByGigID(ctx context.Context, gigID valueobject.GigID) ([]*aggregate.Proposal, error)

	// FindByHustlerID retrieves proposals submitted by a hustler
	FindByHustlerID(ctx context.Context, hustlerID valueobject.UserID, status *aggregate.ProposalStatus, offset, limit int) ([]*ProposalWithGig, int64, error)

	// Delete removes a proposal
	Delete(ctx context.Context, id valueobject.ProposalID) error
}

// ProposalWithGig contains proposal with associated gig info
type ProposalWithGig struct {
	Proposal *aggregate.Proposal
	GigID    valueobject.GigID
	GigTitle string
	GigStatus aggregate.GigStatus
	ClientID valueobject.UserID
}

// ContractRepository defines the interface for contract persistence
type ContractRepository interface {
	// Save persists a contract aggregate
	Save(ctx context.Context, contract *aggregate.Contract) error

	// SaveWithEvents persists a contract and publishes domain events
	SaveWithEvents(ctx context.Context, contract *aggregate.Contract) error

	// FindByID retrieves a contract by ID
	FindByID(ctx context.Context, id valueobject.ContractID) (*aggregate.Contract, error)

	// FindByGigID retrieves contract for a gig
	FindByGigID(ctx context.Context, gigID valueobject.GigID) (*aggregate.Contract, error)

	// FindByUserID retrieves contracts for a user (as client or hustler)
	FindByUserID(ctx context.Context, userID valueobject.UserID, role string, status *aggregate.ContractStatus, offset, limit int) ([]*ContractDTO, int64, error)

	// FindActiveByHustlerID retrieves active contracts for a hustler
	FindActiveByHustlerID(ctx context.Context, hustlerID valueobject.UserID) ([]*aggregate.Contract, error)
}

// ContractDTO represents contract data for listings
type ContractDTO struct {
	ID           string
	GigID        string
	GigTitle     string
	ClientID     string
	ClientName   string
	HustlerID    string
	HustlerName  string
	AgreedPrice  int64
	PlatformFee  int64
	Currency     string
	Status       string
	DeliveryDays int
	StartedAt    time.Time
	DeadlineAt   time.Time
	DeliveredAt  *time.Time
	CompletedAt  *time.Time
	HasReviewed  bool
}

// ReviewRepository defines the interface for review persistence
type ReviewRepository interface {
	// Save persists a review
	Save(ctx context.Context, contractID valueobject.ContractID, review *aggregate.Review) error

	// FindByContractID retrieves reviews for a contract
	FindByContractID(ctx context.Context, contractID valueobject.ContractID) ([]*aggregate.Review, error)

	// FindByRevieweeID retrieves reviews for a user (reviewee)
	FindByRevieweeID(ctx context.Context, userID valueobject.UserID, offset, limit int) ([]*ReviewDTO, int64, error)

	// GetAverageRating gets the average rating for a user
	GetAverageRating(ctx context.Context, userID valueobject.UserID) (float64, int, error)
}

// ReviewDTO represents review data for API responses
type ReviewDTO struct {
	ID                  string
	ContractID          string
	GigTitle            string
	ReviewerID          string
	ReviewerName        string
	ReviewerProfileImage string
	Rating              int
	ReviewText          string
	CommunicationRating int
	QualityRating       int
	TimelinessRating    int
	CreatedAt           time.Time
}

// GigSearchRepository defines the interface for gig search
type GigSearchRepository interface {
	// Search performs full-text search on gigs
	Search(ctx context.Context, query string, filter GigFilter) ([]*GigSearchResult, int64, error)

	// IndexGig indexes a gig for search
	IndexGig(ctx context.Context, gig *aggregate.Gig) error

	// RemoveGig removes a gig from search index
	RemoveGig(ctx context.Context, gigID valueobject.GigID) error
}

// GigSearchResult represents search result
type GigSearchResult struct {
	GigID       string
	Title       string
	Description string
	Category    string
	BudgetMin   int64
	BudgetMax   int64
	Currency    string
	IsRemote    bool
	Location    string
	ClientName  string
	Score       float64
}

// GigStatisticsRepository defines the interface for gig statistics
type GigStatisticsRepository interface {
	// GetUserGigStats gets gig-related statistics for a user
	GetUserGigStats(ctx context.Context, userID valueobject.UserID) (*UserGigStats, error)

	// GetPlatformGigStats gets platform-wide gig statistics
	GetPlatformGigStats(ctx context.Context) (*PlatformGigStats, error)
}

// UserGigStats contains gig statistics for a user
type UserGigStats struct {
	UserID               string
	TotalGigsPosted      int64
	TotalGigsInProgress  int64
	TotalGigsCompleted   int64
	TotalProposalsSent   int64
	TotalProposalsAccepted int64
	TotalContractsAsHustler int64
	TotalContractsCompleted int64
	TotalEarnings        int64
	AverageRating        float64
	TotalReviews         int
}

// PlatformGigStats contains platform-wide statistics
type PlatformGigStats struct {
	TotalOpenGigs       int64
	TotalCompletedGigs  int64
	TotalActiveContracts int64
	TotalTransactionValue int64
	AverageGigValue     int64
}
