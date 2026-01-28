package repository

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/domain/credit/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// Repository errors
var (
	ErrLoanNotFound        = errors.New("loan not found")
	ErrCreditScoreNotFound = errors.New("credit score not found")
)

// CreditScoreRepository defines the interface for credit score persistence
type CreditScoreRepository interface {
	// Save persists a credit score
	Save(ctx context.Context, score *aggregate.CreditScore) error

	// FindByUserID retrieves credit score for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.CreditScore, error)

	// FindByID retrieves a credit score by ID
	FindByID(ctx context.Context, id string) (*aggregate.CreditScore, error)

	// UpdateStats updates specific stats for a user's credit score
	UpdateStats(ctx context.Context, userID valueobject.UserID, stats CreditStatsUpdate) error
}

// CreditStatsUpdate contains stats to update
type CreditStatsUpdate struct {
	GigsCompleted       *int
	GigsAccepted        *int
	AverageRating       *float64
	TotalReviews        *int
	OnTimeContributions *int
	TotalContributions  *int
}

// LoanRepository defines the interface for loan persistence
type LoanRepository interface {
	// Save persists a loan aggregate
	Save(ctx context.Context, loan *aggregate.Loan) error

	// SaveWithEvents persists a loan and publishes domain events
	SaveWithEvents(ctx context.Context, loan *aggregate.Loan) error

	// FindByID retrieves a loan by ID
	FindByID(ctx context.Context, id valueobject.LoanID) (*aggregate.Loan, error)

	// FindByUserID retrieves loans for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID, status *aggregate.LoanStatus) ([]*aggregate.Loan, error)

	// FindActiveByUserID checks if user has an active loan
	FindActiveByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.Loan, error)

	// FindOverdue retrieves overdue loans
	FindOverdue(ctx context.Context) ([]*LoanDTO, error)

	// List retrieves loans with filters
	List(ctx context.Context, filter LoanFilter) ([]*LoanDTO, int64, error)
}

// LoanFilter contains filter options for listing loans
type LoanFilter struct {
	Status  *aggregate.LoanStatus
	MinAmount int64
	MaxAmount int64
	Overdue   bool
	Offset    int
	Limit     int
}

// LoanDTO represents loan data for API responses
type LoanDTO struct {
	ID             string
	UserID         string
	UserName       string
	Principal      int64
	InterestRate   float64
	InterestAmount int64
	TotalAmount    int64
	AmountRepaid   int64
	RemainingBalance int64
	Currency       string
	TenureMonths   int
	Status         string
	Purpose        string
	ApprovedAt     *time.Time
	DisbursedAt    *time.Time
	DueDate        *time.Time
	CompletedAt    *time.Time
	IsOverdue      bool
	CreatedAt      time.Time
}

// RepaymentRepository defines the interface for repayment persistence
type RepaymentRepository interface {
	// Save persists a repayment
	Save(ctx context.Context, loanID valueobject.LoanID, repayment *aggregate.Repayment) error

	// FindByLoanID retrieves repayments for a loan
	FindByLoanID(ctx context.Context, loanID valueobject.LoanID) ([]*RepaymentDTO, error)
}

// RepaymentDTO represents repayment data
type RepaymentDTO struct {
	ID            string
	LoanID        string
	Amount        int64
	Currency      string
	TransactionID string
	PaidAt        time.Time
}

// CreditStatisticsRepository defines the interface for credit statistics
type CreditStatisticsRepository interface {
	// GetLoanStats gets loan statistics for a user
	GetUserLoanStats(ctx context.Context, userID valueobject.UserID) (*UserLoanStats, error)

	// GetPlatformLoanStats gets platform-wide loan statistics
	GetPlatformLoanStats(ctx context.Context) (*PlatformLoanStats, error)
}

// UserLoanStats contains loan statistics for a user
type UserLoanStats struct {
	UserID           string
	TotalLoans       int
	ActiveLoans      int
	CompletedLoans   int
	DefaultedLoans   int
	TotalBorrowed    int64
	TotalRepaid      int64
	CurrentOutstanding int64
}

// PlatformLoanStats contains platform-wide statistics
type PlatformLoanStats struct {
	TotalLoansIssued    int64
	TotalAmountDisbursed int64
	TotalAmountRepaid   int64
	OutstandingBalance  int64
	DefaultRate         float64
	AverageInterestRate float64
}
