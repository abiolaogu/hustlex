package query

import (
	"context"
	"time"

	"hustlex/internal/domain/credit/aggregate"
	"hustlex/internal/domain/credit/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// GetCreditScore retrieves a user's credit score
type GetCreditScore struct {
	UserID string
}

// CreditScoreDTO represents credit score data for API responses
type CreditScoreDTO struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Score              int       `json:"score"`
	Tier               string    `json:"tier"`
	MaxLoanAmount      int64     `json:"max_loan_amount"`
	InterestRate       float64   `json:"interest_rate"`
	GigCompletionScore int       `json:"gig_completion_score"`
	RatingScore        int       `json:"rating_score"`
	SavingsScore       int       `json:"savings_score"`
	AccountAgeScore    int       `json:"account_age_score"`
	VerificationScore  int       `json:"verification_score"`
	CommunityScore     int       `json:"community_score"`
	TotalGigsCompleted int       `json:"total_gigs_completed"`
	AverageRating      float64   `json:"average_rating"`
	LastCalculatedAt   time.Time `json:"last_calculated_at"`
	CreatedAt          time.Time `json:"created_at"`
}

// GetLoan retrieves a single loan
type GetLoan struct {
	LoanID string
	UserID string // for access check
}

// LoanDTO represents loan data for API responses
type LoanDTO struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Principal        int64      `json:"principal"`
	InterestRate     float64    `json:"interest_rate"`
	InterestAmount   int64      `json:"interest_amount"`
	TotalAmount      int64      `json:"total_amount"`
	AmountRepaid     int64      `json:"amount_repaid"`
	RemainingBalance int64      `json:"remaining_balance"`
	Currency         string     `json:"currency"`
	TenureMonths     int        `json:"tenure_months"`
	MonthlyPayment   int64      `json:"monthly_payment"`
	Status           string     `json:"status"`
	Purpose          string     `json:"purpose"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	DisbursedAt      *time.Time `json:"disbursed_at,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	IsOverdue        bool       `json:"is_overdue"`
	Repayments       []RepaymentDTO `json:"repayments,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// RepaymentDTO represents repayment data
type RepaymentDTO struct {
	ID            string    `json:"id"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	TransactionID string    `json:"transaction_id"`
	PaidAt        time.Time `json:"paid_at"`
}

// GetMyLoans retrieves user's loans
type GetMyLoans struct {
	UserID string
	Status string // optional filter
}

// GetLoanHistory retrieves loan history with pagination
type GetLoanHistory struct {
	UserID string
	Page   int
	Limit  int
}

// LoanListResult represents paginated loan results
type LoanListResult struct {
	Loans      []LoanDTO `json:"loans"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

// GetOverdueLoans retrieves overdue loans (admin)
type GetOverdueLoans struct {
	AdminID string
}

// GetLoanStats retrieves loan statistics for a user
type GetLoanStats struct {
	UserID string
}

// LoanStatsDTO represents loan statistics
type LoanStatsDTO struct {
	TotalLoans         int   `json:"total_loans"`
	ActiveLoans        int   `json:"active_loans"`
	CompletedLoans     int   `json:"completed_loans"`
	DefaultedLoans     int   `json:"defaulted_loans"`
	TotalBorrowed      int64 `json:"total_borrowed"`
	TotalRepaid        int64 `json:"total_repaid"`
	CurrentOutstanding int64 `json:"current_outstanding"`
}

// GetPlatformLoanStats retrieves platform-wide loan statistics (admin)
type GetPlatformLoanStats struct {
	AdminID string
}

// PlatformLoanStatsDTO represents platform loan statistics
type PlatformLoanStatsDTO struct {
	TotalLoansIssued     int64   `json:"total_loans_issued"`
	TotalAmountDisbursed int64   `json:"total_amount_disbursed"`
	TotalAmountRepaid    int64   `json:"total_amount_repaid"`
	OutstandingBalance   int64   `json:"outstanding_balance"`
	DefaultRate          float64 `json:"default_rate"`
	AverageInterestRate  float64 `json:"average_interest_rate"`
}

// AdminLoanFilter for listing loans with filters
type AdminLoanFilter struct {
	Status    string
	MinAmount int64
	MaxAmount int64
	Overdue   bool
	Page      int
	Limit     int
}

// CreditQueryHandler handles credit-related queries
type CreditQueryHandler struct {
	creditScoreRepo repository.CreditScoreRepository
	loanRepo        repository.LoanRepository
	repaymentRepo   repository.RepaymentRepository
	statsRepo       repository.CreditStatisticsRepository
}

// NewCreditQueryHandler creates a new query handler
func NewCreditQueryHandler(
	creditScoreRepo repository.CreditScoreRepository,
	loanRepo repository.LoanRepository,
	repaymentRepo repository.RepaymentRepository,
	statsRepo repository.CreditStatisticsRepository,
) *CreditQueryHandler {
	return &CreditQueryHandler{
		creditScoreRepo: creditScoreRepo,
		loanRepo:        loanRepo,
		repaymentRepo:   repaymentRepo,
		statsRepo:       statsRepo,
	}
}

// HandleGetCreditScore retrieves a user's credit score
func (h *CreditQueryHandler) HandleGetCreditScore(ctx context.Context, q GetCreditScore) (*CreditScoreDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	creditScore, err := h.creditScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return creditScoreToDTO(creditScore), nil
}

// HandleGetLoan retrieves a single loan
func (h *CreditQueryHandler) HandleGetLoan(ctx context.Context, q GetLoan) (*LoanDTO, error) {
	loanID, err := valueobject.NewLoanID(q.LoanID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if loan.UserID() != userID {
		return nil, repository.ErrLoanNotFound
	}

	return loanToDTO(loan), nil
}

// HandleGetMyLoans retrieves user's loans
func (h *CreditQueryHandler) HandleGetMyLoans(ctx context.Context, q GetMyLoans) ([]LoanDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	var status *aggregate.LoanStatus
	if q.Status != "" {
		s := aggregate.LoanStatus(q.Status)
		status = &s
	}

	loans, err := h.loanRepo.FindByUserID(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	dtos := make([]LoanDTO, len(loans))
	for i, loan := range loans {
		dtos[i] = *loanToDTO(loan)
	}

	return dtos, nil
}

// HandleGetLoanStats retrieves loan statistics for a user
func (h *CreditQueryHandler) HandleGetLoanStats(ctx context.Context, q GetLoanStats) (*LoanStatsDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	stats, err := h.statsRepo.GetUserLoanStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &LoanStatsDTO{
		TotalLoans:         stats.TotalLoans,
		ActiveLoans:        stats.ActiveLoans,
		CompletedLoans:     stats.CompletedLoans,
		DefaultedLoans:     stats.DefaultedLoans,
		TotalBorrowed:      stats.TotalBorrowed,
		TotalRepaid:        stats.TotalRepaid,
		CurrentOutstanding: stats.CurrentOutstanding,
	}, nil
}

// HandleGetPlatformLoanStats retrieves platform-wide loan statistics
func (h *CreditQueryHandler) HandleGetPlatformLoanStats(ctx context.Context, q GetPlatformLoanStats) (*PlatformLoanStatsDTO, error) {
	stats, err := h.statsRepo.GetPlatformLoanStats(ctx)
	if err != nil {
		return nil, err
	}

	return &PlatformLoanStatsDTO{
		TotalLoansIssued:     stats.TotalLoansIssued,
		TotalAmountDisbursed: stats.TotalAmountDisbursed,
		TotalAmountRepaid:    stats.TotalAmountRepaid,
		OutstandingBalance:   stats.OutstandingBalance,
		DefaultRate:          stats.DefaultRate,
		AverageInterestRate:  stats.AverageInterestRate,
	}, nil
}

// HandleGetOverdueLoans retrieves overdue loans
func (h *CreditQueryHandler) HandleGetOverdueLoans(ctx context.Context, q GetOverdueLoans) ([]repository.LoanDTO, error) {
	return h.loanRepo.FindOverdue(ctx)
}

func creditScoreToDTO(cs *aggregate.CreditScore) *CreditScoreDTO {
	return &CreditScoreDTO{
		ID:                 cs.ID(),
		UserID:             cs.UserID().String(),
		Score:              cs.Score(),
		Tier:               cs.Tier().String(),
		MaxLoanAmount:      cs.MaxLoanAmount(),
		InterestRate:       cs.InterestRate(),
		GigCompletionScore: cs.GigCompletionScore(),
		RatingScore:        cs.RatingScore(),
		SavingsScore:       cs.SavingsScore(),
		AccountAgeScore:    cs.AccountAgeScore(),
		VerificationScore:  cs.VerificationScore(),
		CommunityScore:     cs.CommunityScore(),
		TotalGigsCompleted: cs.TotalGigsCompleted(),
		AverageRating:      cs.AverageRating(),
		LastCalculatedAt:   cs.LastCalculatedAt(),
		CreatedAt:          cs.CreatedAt(),
	}
}

func loanToDTO(loan *aggregate.Loan) *LoanDTO {
	dto := &LoanDTO{
		ID:               loan.ID().String(),
		UserID:           loan.UserID().String(),
		Principal:        loan.Principal().Amount(),
		InterestRate:     loan.InterestRate(),
		InterestAmount:   loan.InterestAmount().Amount(),
		TotalAmount:      loan.TotalAmount().Amount(),
		AmountRepaid:     loan.AmountRepaid().Amount(),
		RemainingBalance: loan.RemainingBalance().Amount(),
		Currency:         string(loan.Principal().Currency()),
		TenureMonths:     loan.TenureMonths(),
		MonthlyPayment:   loan.MonthlyPayment().Amount(),
		Status:           loan.Status().String(),
		Purpose:          loan.Purpose(),
		ApprovedAt:       loan.ApprovedAt(),
		DisbursedAt:      loan.DisbursedAt(),
		DueDate:          loan.DueDate(),
		CompletedAt:      loan.CompletedAt(),
		IsOverdue:        loan.IsOverdue(),
		CreatedAt:        loan.CreatedAt(),
	}

	// Add repayments
	repayments := make([]RepaymentDTO, len(loan.Repayments()))
	for i, r := range loan.Repayments() {
		repayments[i] = RepaymentDTO{
			ID:            r.ID(),
			Amount:        r.Amount().Amount(),
			Currency:      string(r.Amount().Currency()),
			TransactionID: r.TransactionID().String(),
			PaidAt:        r.PaidAt(),
		}
	}
	dto.Repayments = repayments

	return dto
}

// Error sentinel
var ErrLoanNotFound = repository.ErrLoanNotFound
