package repository

import (
	"context"
	"time"

	"hustlex/internal/domain/savings/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// CircleRepository defines the interface for circle persistence
type CircleRepository interface {
	// Save persists a circle aggregate
	Save(ctx context.Context, circle *aggregate.Circle) error

	// SaveWithEvents persists a circle and publishes domain events
	SaveWithEvents(ctx context.Context, circle *aggregate.Circle) error

	// FindByID retrieves a circle by ID
	FindByID(ctx context.Context, id valueobject.CircleID) (*aggregate.Circle, error)

	// FindByInviteCode retrieves a circle by invite code
	FindByInviteCode(ctx context.Context, code string) (*aggregate.Circle, error)

	// FindByUserID retrieves circles a user is a member of
	FindByUserID(ctx context.Context, userID valueobject.UserID, status *aggregate.CircleStatus) ([]*aggregate.Circle, error)

	// List retrieves circles with filters
	List(ctx context.Context, filter CircleFilter) ([]*CircleDTO, int64, error)

	// Delete soft-deletes a circle
	Delete(ctx context.Context, id valueobject.CircleID) error
}

// CircleFilter contains filter options for listing circles
type CircleFilter struct {
	Type       *aggregate.CircleType
	Status     *aggregate.CircleStatus
	MinAmount  int64
	MaxAmount  int64
	Frequency  *aggregate.ContributionFrequency
	Search     string
	IsPublic   bool
	Offset     int
	Limit      int
}

// CircleDTO represents circle data for listings
type CircleDTO struct {
	ID              string
	Name            string
	Description     string
	Type            string
	ContributionAmt int64
	Currency        string
	Frequency       string
	MaxMembers      int
	CurrentMembers  int
	TotalRounds     int
	CurrentRound    int
	Status          string
	IsPrivate       bool
	CreatorID       string
	CreatorName     string
	StartDate       *time.Time
	NextPayoutDate  *time.Time
	CreatedAt       time.Time
}

// MemberRepository defines the interface for member persistence
type MemberRepository interface {
	// Save persists a member
	Save(ctx context.Context, circleID valueobject.CircleID, member *aggregate.Member) error

	// FindByID retrieves a member by ID
	FindByID(ctx context.Context, id valueobject.MemberID) (*aggregate.Member, error)

	// FindByCircleID retrieves all members of a circle
	FindByCircleID(ctx context.Context, circleID valueobject.CircleID) ([]*MemberDTO, error)

	// Delete removes a member
	Delete(ctx context.Context, id valueobject.MemberID) error
}

// MemberDTO represents member data for API responses
type MemberDTO struct {
	ID             string
	UserID         string
	UserName       string
	UserImage      string
	Position       int
	Role           string
	Status         string
	TotalContrib   int64
	MissedPayments int
	HasReceived    bool
	JoinedAt       time.Time
}

// ContributionRepository defines the interface for contribution persistence
type ContributionRepository interface {
	// Save persists a contribution
	Save(ctx context.Context, circleID valueobject.CircleID, contribution *aggregate.Contribution) error

	// FindByID retrieves a contribution by ID
	FindByID(ctx context.Context, id valueobject.ContributionID) (*aggregate.Contribution, error)

	// FindByCircleAndRound retrieves contributions for a specific round
	FindByCircleAndRound(ctx context.Context, circleID valueobject.CircleID, round int) ([]*ContributionDTO, error)

	// FindByMember retrieves contributions for a member
	FindByMember(ctx context.Context, memberID valueobject.MemberID) ([]*ContributionDTO, error)

	// FindPending retrieves pending contributions
	FindPending(ctx context.Context, circleID valueobject.CircleID) ([]*ContributionDTO, error)

	// FindOverdue retrieves overdue contributions
	FindOverdue(ctx context.Context) ([]*ContributionDTO, error)
}

// ContributionDTO represents contribution data for API responses
type ContributionDTO struct {
	ID            string
	CircleID      string
	CircleName    string
	MemberID      string
	UserID        string
	UserName      string
	Round         int
	Amount        int64
	Currency      string
	DueDate       time.Time
	PaidAt        *time.Time
	Status        string
	TransactionID *string
	LateFee       int64
}

// PayoutRepository defines the interface for payout tracking
type PayoutRepository interface {
	// RecordPayout records a payout made to a member
	RecordPayout(ctx context.Context, payout *Payout) error

	// FindByCircle retrieves payouts for a circle
	FindByCircle(ctx context.Context, circleID valueobject.CircleID) ([]*Payout, error)

	// FindByUser retrieves payouts received by a user
	FindByUser(ctx context.Context, userID valueobject.UserID) ([]*Payout, error)
}

// Payout represents a payout record
type Payout struct {
	ID            string
	CircleID      string
	CircleName    string
	MemberID      string
	UserID        string
	Round         int
	Amount        int64
	Currency      string
	TransactionID string
	PaidAt        time.Time
}

// SavingsStatisticsRepository defines the interface for savings statistics
type SavingsStatisticsRepository interface {
	// GetUserStats gets savings statistics for a user
	GetUserStats(ctx context.Context, userID valueobject.UserID) (*UserSavingsStats, error)

	// GetCircleStats gets statistics for a circle
	GetCircleStats(ctx context.Context, circleID valueobject.CircleID) (*CircleStats, error)
}

// UserSavingsStats contains savings statistics for a user
type UserSavingsStats struct {
	UserID           string
	TotalCircles     int
	ActiveCircles    int
	CompletedCircles int
	TotalContributed int64
	TotalReceived    int64
	MissedPayments   int
	OnTimeRate       float64
}

// CircleStats contains statistics for a circle
type CircleStats struct {
	CircleID        string
	TotalMembers    int
	ActiveMembers   int
	TotalContributed int64
	TotalPaidOut    int64
	CompletionRate  float64
	AverageLateFee  int64
}
