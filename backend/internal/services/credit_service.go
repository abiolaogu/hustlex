package services

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"hustlex/internal/models"
)

// CreditService handles credit score calculations and loan management
type CreditService struct {
	db *gorm.DB
}

// NewCreditService creates a new credit service
func NewCreditService(db *gorm.DB) *CreditService {
	return &CreditService{db: db}
}

// Credit score constants
const (
	MaxCreditScore = 850
	MinCreditScore = 0

	// Component weights (must sum to 100)
	WeightGigCompletion    = 25 // 25%
	WeightGigRatings       = 20 // 20%
	WeightAjoRecord        = 20 // 20%
	WeightAccountAge       = 15 // 15%
	WeightVerificationLevel = 10 // 10%
	WeightCommunityStanding = 10 // 10%

	// Tier thresholds
	TierBronzeMax   = 300
	TierSilverMax   = 500
	TierGoldMax     = 700
	TierPlatinumMax = 850
)

// Credit errors
var (
	ErrCreditScoreNotFound  = errors.New("credit score not found")
	ErrLoanNotFound         = errors.New("loan not found")
	ErrLoanNotEligible      = errors.New("not eligible for loan")
	ErrActiveLoanExists     = errors.New("active loan already exists")
	ErrLoanAmountTooHigh    = errors.New("loan amount exceeds limit")
	ErrLoanAmountTooLow     = errors.New("loan amount below minimum")
	ErrLoanAlreadyRepaid    = errors.New("loan already fully repaid")
	ErrInvalidRepaymentAmount = errors.New("invalid repayment amount")
)

// Loan limits by tier (in Kobo)
var LoanLimits = map[models.UserTier]int64{
	models.TierBronze:   5000000,    // ₦50,000
	models.TierSilver:   20000000,   // ₦200,000
	models.TierGold:     50000000,   // ₦500,000
	models.TierPlatinum: 100000000,  // ₦1,000,000
}

// Interest rates by tier (monthly %)
var InterestRates = map[models.UserTier]float64{
	models.TierBronze:   5.0,  // 5% per month
	models.TierSilver:   4.0,  // 4% per month
	models.TierGold:     3.5,  // 3.5% per month
	models.TierPlatinum: 3.0,  // 3% per month
}

// GetCreditScore retrieves or creates a user's credit score
func (s *CreditService) GetCreditScore(ctx context.Context, userID uuid.UUID) (*models.CreditScore, error) {
	var score models.CreditScore
	err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&score).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new credit score
			score = models.CreditScore{
				UserID:              userID,
				Score:               0,
				GigCompletionScore:  0,
				GigRatingScore:      0,
				AjoRecordScore:      0,
				AccountAgeScore:     0,
				VerificationScore:   0,
				CommunityScore:      0,
				TotalGigsCompleted:  0,
				TotalGigsCancelled:  0,
				TotalAjoCompleted:   0,
				TotalAjoDefaulted:   0,
				TotalLoansTaken:     0,
				TotalLoansRepaid:    0,
			}
			if err := s.db.WithContext(ctx).Create(&score).Error; err != nil {
				return nil, err
			}
			return &score, nil
		}
		return nil, err
	}

	return &score, nil
}

// RecalculateCreditScore recalculates a user's credit score based on all factors
func (s *CreditService) RecalculateCreditScore(ctx context.Context, userID uuid.UUID) (*models.CreditScore, error) {
	score, err := s.GetCreditScore(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get user for account age
	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Calculate component scores
	gigCompletionScore := s.calculateGigCompletionScore(ctx, userID)
	gigRatingScore := s.calculateGigRatingScore(ctx, userID)
	ajoRecordScore := s.calculateAjoRecordScore(ctx, userID)
	accountAgeScore := s.calculateAccountAgeScore(user.CreatedAt)
	verificationScore := s.calculateVerificationScore(ctx, userID)
	communityScore := s.calculateCommunityScore(ctx, userID)

	// Calculate weighted total
	totalScore := int(
		(float64(gigCompletionScore) * WeightGigCompletion / 100) +
		(float64(gigRatingScore) * WeightGigRatings / 100) +
		(float64(ajoRecordScore) * WeightAjoRecord / 100) +
		(float64(accountAgeScore) * WeightAccountAge / 100) +
		(float64(verificationScore) * WeightVerificationLevel / 100) +
		(float64(communityScore) * WeightCommunityStanding / 100),
	)

	// Cap score
	if totalScore > MaxCreditScore {
		totalScore = MaxCreditScore
	}
	if totalScore < MinCreditScore {
		totalScore = MinCreditScore
	}

	// Get activity stats
	var gigsCompleted, gigsCancelled int64
	s.db.WithContext(ctx).Model(&models.GigContract{}).
		Where("hustler_id = ? AND status = ?", userID, models.ContractCompleted).
		Count(&gigsCompleted)
	s.db.WithContext(ctx).Model(&models.GigContract{}).
		Where("hustler_id = ? AND status = ?", userID, models.ContractCancelled).
		Count(&gigsCancelled)

	var ajoCompleted, ajoDefaulted int64
	s.db.WithContext(ctx).Model(&models.Contribution{}).
		Joins("JOIN circle_members ON contributions.member_id = circle_members.id").
		Where("circle_members.user_id = ? AND contributions.status = ?", userID, "paid").
		Count(&ajoCompleted)
	s.db.WithContext(ctx).Model(&models.Contribution{}).
		Joins("JOIN circle_members ON contributions.member_id = circle_members.id").
		Where("circle_members.user_id = ? AND contributions.status = ?", userID, "overdue").
		Count(&ajoDefaulted)

	// Update score record
	score.Score = totalScore
	score.GigCompletionScore = gigCompletionScore
	score.GigRatingScore = gigRatingScore
	score.AjoRecordScore = ajoRecordScore
	score.AccountAgeScore = accountAgeScore
	score.VerificationScore = verificationScore
	score.CommunityScore = communityScore
	score.TotalGigsCompleted = int(gigsCompleted)
	score.TotalGigsCancelled = int(gigsCancelled)
	score.TotalAjoCompleted = int(ajoCompleted)
	score.TotalAjoDefaulted = int(ajoDefaulted)
	score.LastCalculatedAt = time.Now().UTC()

	if err := s.db.WithContext(ctx).Save(score).Error; err != nil {
		return nil, err
	}

	// Update user tier based on score
	if err := s.updateUserTier(ctx, userID, totalScore); err != nil {
		return nil, err
	}

	return score, nil
}

// calculateGigCompletionScore scores based on gig completion rate
func (s *CreditService) calculateGigCompletionScore(ctx context.Context, userID uuid.UUID) int {
	var completed, cancelled int64

	// As hustler
	s.db.WithContext(ctx).Model(&models.GigContract{}).
		Where("hustler_id = ? AND status = ?", userID, models.ContractCompleted).
		Count(&completed)
	s.db.WithContext(ctx).Model(&models.GigContract{}).
		Where("hustler_id = ? AND status = ?", userID, models.ContractCancelled).
		Count(&cancelled)

	total := completed + cancelled
	if total == 0 {
		return 0
	}

	completionRate := float64(completed) / float64(total)

	// Scale: 90%+ = max score, below 70% = low score
	if completionRate >= 0.95 {
		return MaxCreditScore
	} else if completionRate >= 0.90 {
		return 750
	} else if completionRate >= 0.85 {
		return 650
	} else if completionRate >= 0.80 {
		return 550
	} else if completionRate >= 0.70 {
		return 400
	} else if completionRate >= 0.50 {
		return 250
	}
	return 100
}

// calculateGigRatingScore scores based on average gig ratings
func (s *CreditService) calculateGigRatingScore(ctx context.Context, userID uuid.UUID) int {
	var avgRating float64
	var count int64

	err := s.db.WithContext(ctx).Model(&models.GigReview{}).
		Where("reviewee_id = ?", userID).
		Select("AVG(overall_rating)").
		Scan(&avgRating).Error

	if err != nil {
		return 0
	}

	s.db.WithContext(ctx).Model(&models.GigReview{}).
		Where("reviewee_id = ?", userID).
		Count(&count)

	if count == 0 {
		return 0
	}

	// Scale 1-5 rating to 0-850
	// 5.0 = 850, 4.5 = 750, 4.0 = 600, 3.5 = 450, 3.0 = 300, below 3 = scaled down
	scaledScore := int((avgRating - 1) / 4 * float64(MaxCreditScore))

	// Apply volume bonus (more reviews = more confidence)
	volumeMultiplier := math.Min(1.0, float64(count)/50) // Max bonus at 50 reviews
	scaledScore = int(float64(scaledScore) * (0.7 + 0.3*volumeMultiplier))

	return scaledScore
}

// calculateAjoRecordScore scores based on Ajo/Esusu participation
func (s *CreditService) calculateAjoRecordScore(ctx context.Context, userID uuid.UUID) int {
	var onTimeContributions, totalContributions int64

	// Get all member records for user
	var memberIDs []uuid.UUID
	s.db.WithContext(ctx).Model(&models.CircleMember{}).
		Where("user_id = ?", userID).
		Pluck("id", &memberIDs)

	if len(memberIDs) == 0 {
		return 0
	}

	// Count contributions
	s.db.WithContext(ctx).Model(&models.Contribution{}).
		Where("member_id IN ? AND status = ?", memberIDs, "paid").
		Count(&onTimeContributions)

	s.db.WithContext(ctx).Model(&models.Contribution{}).
		Where("member_id IN ?", memberIDs).
		Count(&totalContributions)

	if totalContributions == 0 {
		return 0
	}

	onTimeRate := float64(onTimeContributions) / float64(totalContributions)

	// Scale: 100% on time = max, deduct heavily for missed payments
	baseScore := int(onTimeRate * float64(MaxCreditScore))

	// Count circles completed
	var circlesCompleted int64
	s.db.WithContext(ctx).Model(&models.CircleMember{}).
		Joins("JOIN savings_circles ON circle_members.circle_id = savings_circles.id").
		Where("circle_members.user_id = ? AND savings_circles.status = ?", userID, "completed").
		Count(&circlesCompleted)

	// Bonus for completing full circles
	circleBonus := int(math.Min(100, float64(circlesCompleted)*20))

	return int(math.Min(float64(MaxCreditScore), float64(baseScore+circleBonus)))
}

// calculateAccountAgeScore scores based on how long user has been on platform
func (s *CreditService) calculateAccountAgeScore(createdAt time.Time) int {
	ageInDays := int(time.Since(createdAt).Hours() / 24)

	// Scale: 0-7 days = 0, 30 days = 200, 90 days = 400, 180 days = 600, 365+ days = 850
	if ageInDays < 7 {
		return 0
	} else if ageInDays < 30 {
		return int(float64(ageInDays-7) / 23 * 200)
	} else if ageInDays < 90 {
		return 200 + int(float64(ageInDays-30)/60*200)
	} else if ageInDays < 180 {
		return 400 + int(float64(ageInDays-90)/90*200)
	} else if ageInDays < 365 {
		return 600 + int(float64(ageInDays-180)/185*200)
	}
	return MaxCreditScore
}

// calculateVerificationScore scores based on verification level
func (s *CreditService) calculateVerificationScore(ctx context.Context, userID uuid.UUID) int {
	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return 0
	}

	score := 0

	// Phone verified (required for registration) - 200 points
	if user.PhoneVerified {
		score += 200
	}

	// Email verified - 150 points
	if user.EmailVerified {
		score += 150
	}

	// BVN verified - 300 points
	if user.BVNVerified {
		score += 300
	}

	// NIN verified - 200 points
	if user.NINVerified {
		score += 200
	}

	return score
}

// calculateCommunityScore scores based on community engagement
func (s *CreditService) calculateCommunityScore(ctx context.Context, userID uuid.UUID) int {
	score := 0

	// Referrals made (10 points each, max 200)
	var referralCount int64
	s.db.WithContext(ctx).Model(&models.User{}).
		Where("referred_by = ?", userID).
		Count(&referralCount)
	score += int(math.Min(200, float64(referralCount*10)))

	// Savings circles created (50 points each, max 200)
	var circlesCreated int64
	s.db.WithContext(ctx).Model(&models.SavingsCircle{}).
		Where("creator_id = ?", userID).
		Count(&circlesCreated)
	score += int(math.Min(200, float64(circlesCreated*50)))

	// Gigs posted (20 points each, max 200)
	var gigsPosted int64
	s.db.WithContext(ctx).Model(&models.Gig{}).
		Where("client_id = ?", userID).
		Count(&gigsPosted)
	score += int(math.Min(200, float64(gigsPosted*20)))

	// Reviews given (10 points each, max 150)
	var reviewsGiven int64
	s.db.WithContext(ctx).Model(&models.GigReview{}).
		Where("reviewer_id = ?", userID).
		Count(&reviewsGiven)
	score += int(math.Min(150, float64(reviewsGiven*10)))

	return int(math.Min(float64(MaxCreditScore), float64(score)))
}

// updateUserTier updates the user's tier based on their credit score
func (s *CreditService) updateUserTier(ctx context.Context, userID uuid.UUID, score int) error {
	var tier models.UserTier

	switch {
	case score > TierGoldMax:
		tier = models.TierPlatinum
	case score > TierSilverMax:
		tier = models.TierGold
	case score > TierBronzeMax:
		tier = models.TierSilver
	default:
		tier = models.TierBronze
	}

	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", userID).
		Update("tier", tier).Error
}

// GetLoanEligibility checks if a user is eligible for a loan
func (s *CreditService) GetLoanEligibility(ctx context.Context, userID uuid.UUID) (*LoanEligibility, error) {
	creditScore, err := s.RecalculateCreditScore(ctx, userID)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Check for active loans
	var activeLoanCount int64
	s.db.WithContext(ctx).Model(&models.Loan{}).
		Where("borrower_id = ? AND status IN ?", userID, []string{"active", "overdue"}).
		Count(&activeLoanCount)

	eligibility := &LoanEligibility{
		CreditScore:    creditScore.Score,
		Tier:           user.Tier,
		MaxLoanAmount:  LoanLimits[user.Tier],
		InterestRate:   InterestRates[user.Tier],
		HasActiveLoan:  activeLoanCount > 0,
		IsEligible:     creditScore.Score >= TierBronzeMax && activeLoanCount == 0,
		MinCreditScore: TierBronzeMax,
	}

	if !eligibility.IsEligible {
		if activeLoanCount > 0 {
			eligibility.IneligibilityReason = "You have an active loan that must be repaid first"
		} else {
			eligibility.IneligibilityReason = "Your credit score is below the minimum required"
		}
	}

	return eligibility, nil
}

// LoanEligibility represents loan eligibility status
type LoanEligibility struct {
	CreditScore         int             `json:"credit_score"`
	Tier                models.UserTier `json:"tier"`
	MaxLoanAmount       int64           `json:"max_loan_amount"`
	InterestRate        float64         `json:"interest_rate"`
	HasActiveLoan       bool            `json:"has_active_loan"`
	IsEligible          bool            `json:"is_eligible"`
	MinCreditScore      int             `json:"min_credit_score"`
	IneligibilityReason string          `json:"ineligibility_reason,omitempty"`
}

// RequestLoanInput represents a loan request
type RequestLoanInput struct {
	UserID        uuid.UUID
	AmountKobo    int64
	DurationDays  int
	Purpose       string
}

// RequestLoan creates a new loan request
func (s *CreditService) RequestLoan(ctx context.Context, input RequestLoanInput) (*models.Loan, error) {
	eligibility, err := s.GetLoanEligibility(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	if !eligibility.IsEligible {
		return nil, ErrLoanNotEligible
	}

	if input.AmountKobo > eligibility.MaxLoanAmount {
		return nil, ErrLoanAmountTooHigh
	}

	minLoan := int64(1000000) // ₦10,000 minimum
	if input.AmountKobo < minLoan {
		return nil, ErrLoanAmountTooLow
	}

	// Validate duration (7-90 days)
	if input.DurationDays < 7 {
		input.DurationDays = 7
	}
	if input.DurationDays > 90 {
		input.DurationDays = 90
	}

	// Calculate interest (simple interest)
	monthlyInterestRate := eligibility.InterestRate / 100
	durationMonths := float64(input.DurationDays) / 30
	interestKobo := int64(float64(input.AmountKobo) * monthlyInterestRate * durationMonths)
	totalRepayment := input.AmountKobo + interestKobo

	dueDate := time.Now().UTC().AddDate(0, 0, input.DurationDays)

	loan := &models.Loan{
		BorrowerID:       input.UserID,
		PrincipalKobo:    input.AmountKobo,
		InterestRatePercent: eligibility.InterestRate,
		InterestKobo:     interestKobo,
		TotalRepaymentKobo: totalRepayment,
		OutstandingKobo:  totalRepayment,
		DurationDays:     input.DurationDays,
		DueDate:          dueDate,
		Status:           "pending",
		Purpose:          input.Purpose,
	}

	if err := s.db.WithContext(ctx).Create(loan).Error; err != nil {
		return nil, err
	}

	return loan, nil
}

// ApproveLoan approves a pending loan and disburses funds
func (s *CreditService) ApproveLoan(ctx context.Context, loanID uuid.UUID, walletService *WalletService) (*models.Loan, error) {
	var loan models.Loan
	if err := s.db.WithContext(ctx).Where("id = ?", loanID).First(&loan).Error; err != nil {
		return nil, ErrLoanNotFound
	}

	if loan.Status != "pending" {
		return nil, errors.New("loan is not pending")
	}

	// Update loan status
	loan.Status = "active"
	loan.DisbursedAt = timePtr(time.Now().UTC())

	if err := s.db.WithContext(ctx).Save(&loan).Error; err != nil {
		return nil, err
	}

	// Disburse to wallet
	_, err := walletService.Deposit(ctx, DepositInput{
		UserID:     loan.BorrowerID,
		AmountKobo: loan.PrincipalKobo,
		Reference:  "LOAN-" + loan.ID.String()[:8],
		Channel:    "loan_disbursement",
	})

	if err != nil {
		// Rollback loan status
		loan.Status = "pending"
		loan.DisbursedAt = nil
		s.db.WithContext(ctx).Save(&loan)
		return nil, err
	}

	// Update credit score stats
	s.db.WithContext(ctx).Model(&models.CreditScore{}).
		Where("user_id = ?", loan.BorrowerID).
		Update("total_loans_taken", gorm.Expr("total_loans_taken + 1"))

	return &loan, nil
}

// RepayLoan processes a loan repayment
func (s *CreditService) RepayLoan(ctx context.Context, loanID uuid.UUID, amountKobo int64, walletService *WalletService, pin string) (*models.LoanRepayment, error) {
	if amountKobo <= 0 {
		return nil, ErrInvalidRepaymentAmount
	}

	var loan models.Loan
	if err := s.db.WithContext(ctx).Where("id = ?", loanID).First(&loan).Error; err != nil {
		return nil, ErrLoanNotFound
	}

	if loan.Status == "repaid" {
		return nil, ErrLoanAlreadyRepaid
	}

	if amountKobo > loan.OutstandingKobo {
		amountKobo = loan.OutstandingKobo
	}

	// Get user wallet and verify they have sufficient balance
	wallet, err := walletService.GetWallet(ctx, loan.BorrowerID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amountKobo {
		return nil, ErrInsufficientBalance
	}

	var repayment *models.LoanRepayment

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create repayment record
		repayment = &models.LoanRepayment{
			LoanID:     loan.ID,
			AmountKobo: amountKobo,
			Status:     "completed",
		}

		if err := tx.Create(repayment).Error; err != nil {
			return err
		}

		// Update loan outstanding
		loan.OutstandingKobo -= amountKobo
		if loan.OutstandingKobo <= 0 {
			loan.OutstandingKobo = 0
			loan.Status = "repaid"
			loan.RepaidAt = timePtr(time.Now().UTC())

			// Update credit score stats
			tx.Model(&models.CreditScore{}).
				Where("user_id = ?", loan.BorrowerID).
				Update("total_loans_repaid", gorm.Expr("total_loans_repaid + 1"))
		}

		return tx.Save(&loan).Error
	})

	if err != nil {
		return nil, err
	}

	// Deduct from wallet
	_, err = walletService.Withdraw(ctx, WithdrawalInput{
		UserID:        loan.BorrowerID,
		AmountKobo:    amountKobo,
		PIN:           pin,
		BankCode:      "INTERNAL",
		AccountNumber: "LOAN_REPAYMENT",
		AccountName:   "Loan Repayment",
	})

	if err != nil {
		// Rollback repayment record
		s.db.WithContext(ctx).Delete(repayment)
		return nil, err
	}

	// Recalculate credit score
	go s.RecalculateCreditScore(context.Background(), loan.BorrowerID)

	return repayment, nil
}

// GetUserLoans retrieves all loans for a user
func (s *CreditService) GetUserLoans(ctx context.Context, userID uuid.UUID) ([]models.Loan, error) {
	var loans []models.Loan
	err := s.db.WithContext(ctx).
		Where("borrower_id = ?", userID).
		Order("created_at DESC").
		Find(&loans).Error

	return loans, err
}

// GetLoan retrieves a specific loan
func (s *CreditService) GetLoan(ctx context.Context, loanID uuid.UUID) (*models.Loan, error) {
	var loan models.Loan
	err := s.db.WithContext(ctx).
		Preload("Repayments").
		Where("id = ?", loanID).
		First(&loan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}

	return &loan, nil
}

// GetCreditHistory retrieves credit score history (for future implementation)
func (s *CreditService) GetCreditHistory(ctx context.Context, userID uuid.UUID) (*CreditHistory, error) {
	score, err := s.GetCreditScore(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get loan history summary
	var totalLoans, repaidLoans, defaultedLoans int64
	s.db.WithContext(ctx).Model(&models.Loan{}).Where("borrower_id = ?", userID).Count(&totalLoans)
	s.db.WithContext(ctx).Model(&models.Loan{}).Where("borrower_id = ? AND status = ?", userID, "repaid").Count(&repaidLoans)
	s.db.WithContext(ctx).Model(&models.Loan{}).Where("borrower_id = ? AND status = ?", userID, "defaulted").Count(&defaultedLoans)

	return &CreditHistory{
		CurrentScore:       score.Score,
		ComponentScores: map[string]int{
			"gig_completion":    score.GigCompletionScore,
			"gig_ratings":       score.GigRatingScore,
			"ajo_record":        score.AjoRecordScore,
			"account_age":       score.AccountAgeScore,
			"verification":      score.VerificationScore,
			"community":         score.CommunityScore,
		},
		Stats: map[string]int{
			"gigs_completed":    score.TotalGigsCompleted,
			"gigs_cancelled":    score.TotalGigsCancelled,
			"ajo_completed":     score.TotalAjoCompleted,
			"ajo_defaulted":     score.TotalAjoDefaulted,
			"loans_taken":       score.TotalLoansTaken,
			"loans_repaid":      score.TotalLoansRepaid,
		},
		LoanSummary: LoanSummary{
			TotalLoans:    int(totalLoans),
			RepaidLoans:   int(repaidLoans),
			DefaultedLoans: int(defaultedLoans),
			ActiveLoans:   int(totalLoans - repaidLoans - defaultedLoans),
		},
		LastCalculatedAt: score.LastCalculatedAt,
	}, nil
}

// CreditHistory represents complete credit history
type CreditHistory struct {
	CurrentScore     int            `json:"current_score"`
	ComponentScores  map[string]int `json:"component_scores"`
	Stats            map[string]int `json:"stats"`
	LoanSummary      LoanSummary    `json:"loan_summary"`
	LastCalculatedAt time.Time      `json:"last_calculated_at"`
}

// LoanSummary represents loan activity summary
type LoanSummary struct {
	TotalLoans     int `json:"total_loans"`
	RepaidLoans    int `json:"repaid_loans"`
	DefaultedLoans int `json:"defaulted_loans"`
	ActiveLoans    int `json:"active_loans"`
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}
