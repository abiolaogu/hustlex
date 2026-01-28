package aggregate

import (
	"errors"
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrInvalidScoreComponent = errors.New("score component must be between 0 and 100")
	ErrScoreOutOfRange       = errors.New("credit score must be between 0 and 850")
)

// UserTier represents the credit tier
type UserTier string

const (
	TierBronze   UserTier = "bronze"
	TierSilver   UserTier = "silver"
	TierGold     UserTier = "gold"
	TierPlatinum UserTier = "platinum"
)

func (t UserTier) String() string {
	return string(t)
}

// TierFromScore determines tier based on credit score
func TierFromScore(score int) UserTier {
	switch {
	case score >= 750:
		return TierPlatinum
	case score >= 600:
		return TierGold
	case score >= 400:
		return TierSilver
	default:
		return TierBronze
	}
}

// TierLoanLimit returns the maximum loan amount for a tier (in kobo)
func (t UserTier) MaxLoanAmount() int64 {
	switch t {
	case TierPlatinum:
		return 500000 * 100 // ₦500,000
	case TierGold:
		return 200000 * 100 // ₦200,000
	case TierSilver:
		return 50000 * 100 // ₦50,000
	case TierBronze:
		return 10000 * 100 // ₦10,000
	default:
		return 0
	}
}

// TierInterestRate returns the monthly interest rate for a tier
func (t UserTier) InterestRate() float64 {
	switch t {
	case TierPlatinum:
		return 0.02 // 2%
	case TierGold:
		return 0.03 // 3%
	case TierSilver:
		return 0.04 // 4%
	case TierBronze:
		return 0.05 // 5%
	default:
		return 0.05
	}
}

// CreditScore is the aggregate root for user credit scores
type CreditScore struct {
	sharedevent.AggregateRoot

	id                   string
	userID               valueobject.UserID
	score                int // 0-850
	tier                 UserTier

	// Score components (0-100 each)
	gigCompletionScore   int
	ratingScore          int
	savingsScore         int
	accountAgeScore      int
	verificationScore    int
	communityScore       int

	// Stats
	totalGigsCompleted   int
	totalGigsAccepted    int
	averageRating        float64
	totalReviews         int
	onTimeContributions  int
	totalContributions   int

	lastCalculatedAt     time.Time
	createdAt            time.Time
	updatedAt            time.Time
	version              int64
}

// NewCreditScore creates a new credit score for a user
func NewCreditScore(userID valueobject.UserID) *CreditScore {
	cs := &CreditScore{
		id:               valueobject.GenerateUserID().String(),
		userID:           userID,
		score:            100, // Starting score
		tier:             TierBronze,
		lastCalculatedAt: time.Now().UTC(),
		createdAt:        time.Now().UTC(),
		updatedAt:        time.Now().UTC(),
		version:          1,
	}
	return cs
}

// ReconstructCreditScore reconstructs from persistence
func ReconstructCreditScore(
	id string,
	userID valueobject.UserID,
	score int,
	tier UserTier,
	gigCompletionScore int,
	ratingScore int,
	savingsScore int,
	accountAgeScore int,
	verificationScore int,
	communityScore int,
	totalGigsCompleted int,
	totalGigsAccepted int,
	averageRating float64,
	totalReviews int,
	onTimeContributions int,
	totalContributions int,
	lastCalculatedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	version int64,
) *CreditScore {
	return &CreditScore{
		id:                  id,
		userID:              userID,
		score:               score,
		tier:                tier,
		gigCompletionScore:  gigCompletionScore,
		ratingScore:         ratingScore,
		savingsScore:        savingsScore,
		accountAgeScore:     accountAgeScore,
		verificationScore:   verificationScore,
		communityScore:      communityScore,
		totalGigsCompleted:  totalGigsCompleted,
		totalGigsAccepted:   totalGigsAccepted,
		averageRating:       averageRating,
		totalReviews:        totalReviews,
		onTimeContributions: onTimeContributions,
		totalContributions:  totalContributions,
		lastCalculatedAt:    lastCalculatedAt,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
		version:             version,
	}
}

// Getters
func (cs *CreditScore) ID() string                 { return cs.id }
func (cs *CreditScore) UserID() valueobject.UserID { return cs.userID }
func (cs *CreditScore) Score() int                 { return cs.score }
func (cs *CreditScore) Tier() UserTier             { return cs.tier }
func (cs *CreditScore) GigCompletionScore() int    { return cs.gigCompletionScore }
func (cs *CreditScore) RatingScore() int           { return cs.ratingScore }
func (cs *CreditScore) SavingsScore() int          { return cs.savingsScore }
func (cs *CreditScore) AccountAgeScore() int       { return cs.accountAgeScore }
func (cs *CreditScore) VerificationScore() int     { return cs.verificationScore }
func (cs *CreditScore) CommunityScore() int        { return cs.communityScore }
func (cs *CreditScore) TotalGigsCompleted() int    { return cs.totalGigsCompleted }
func (cs *CreditScore) TotalGigsAccepted() int     { return cs.totalGigsAccepted }
func (cs *CreditScore) AverageRating() float64     { return cs.averageRating }
func (cs *CreditScore) TotalReviews() int          { return cs.totalReviews }
func (cs *CreditScore) OnTimeContributions() int   { return cs.onTimeContributions }
func (cs *CreditScore) TotalContributions() int    { return cs.totalContributions }
func (cs *CreditScore) LastCalculatedAt() time.Time { return cs.lastCalculatedAt }
func (cs *CreditScore) CreatedAt() time.Time       { return cs.createdAt }
func (cs *CreditScore) UpdatedAt() time.Time       { return cs.updatedAt }
func (cs *CreditScore) Version() int64             { return cs.version }

// MaxLoanAmount returns the maximum loan this user can take
func (cs *CreditScore) MaxLoanAmount() int64 {
	return cs.tier.MaxLoanAmount()
}

// InterestRate returns the monthly interest rate for this user
func (cs *CreditScore) InterestRate() float64 {
	return cs.tier.InterestRate()
}

// Business Methods

// UpdateGigStats updates gig-related statistics
func (cs *CreditScore) UpdateGigStats(completed, accepted int) {
	cs.totalGigsCompleted = completed
	cs.totalGigsAccepted = accepted

	// Calculate gig completion score
	if accepted > 0 {
		completionRate := float64(completed) / float64(accepted)
		cs.gigCompletionScore = int(completionRate * 100)
		if cs.gigCompletionScore > 100 {
			cs.gigCompletionScore = 100
		}
	}

	cs.updatedAt = time.Now().UTC()
}

// UpdateRatingStats updates rating-related statistics
func (cs *CreditScore) UpdateRatingStats(averageRating float64, totalReviews int) {
	cs.averageRating = averageRating
	cs.totalReviews = totalReviews

	// Calculate rating score (5-star scale to 0-100)
	cs.ratingScore = int(averageRating * 20)
	if cs.ratingScore > 100 {
		cs.ratingScore = 100
	}

	cs.updatedAt = time.Now().UTC()
}

// UpdateSavingsStats updates savings-related statistics
func (cs *CreditScore) UpdateSavingsStats(onTime, total int) {
	cs.onTimeContributions = onTime
	cs.totalContributions = total

	// Calculate savings score
	if total > 0 {
		onTimeRate := float64(onTime) / float64(total)
		cs.savingsScore = int(onTimeRate * 100)
	}

	cs.updatedAt = time.Now().UTC()
}

// UpdateAccountAgeScore updates score based on account age
func (cs *CreditScore) UpdateAccountAgeScore(accountAge time.Duration) {
	months := int(accountAge.Hours() / 24 / 30)
	// Max score at 24 months
	cs.accountAgeScore = months * 5
	if cs.accountAgeScore > 100 {
		cs.accountAgeScore = 100
	}

	cs.updatedAt = time.Now().UTC()
}

// UpdateVerificationScore updates score based on verifications
func (cs *CreditScore) UpdateVerificationScore(hasPhone, hasEmail, hasBVN, hasNIN bool) {
	score := 0
	if hasPhone {
		score += 25
	}
	if hasEmail {
		score += 25
	}
	if hasBVN {
		score += 25
	}
	if hasNIN {
		score += 25
	}
	cs.verificationScore = score

	cs.updatedAt = time.Now().UTC()
}

// UpdateCommunityScore updates community participation score
func (cs *CreditScore) UpdateCommunityScore(circlesJoined, referrals int) {
	score := circlesJoined * 10 + referrals * 5
	if score > 100 {
		score = 100
	}
	cs.communityScore = score

	cs.updatedAt = time.Now().UTC()
}

// Recalculate recalculates the overall credit score
func (cs *CreditScore) Recalculate() {
	// Weighted average of components
	// Weights: gig=30%, rating=25%, savings=20%, verification=10%, age=10%, community=5%
	weightedScore := float64(cs.gigCompletionScore)*0.30 +
		float64(cs.ratingScore)*0.25 +
		float64(cs.savingsScore)*0.20 +
		float64(cs.verificationScore)*0.10 +
		float64(cs.accountAgeScore)*0.10 +
		float64(cs.communityScore)*0.05

	// Scale to 0-850
	cs.score = int(weightedScore * 8.5)
	if cs.score > 850 {
		cs.score = 850
	}
	if cs.score < 0 {
		cs.score = 0
	}

	// Update tier
	cs.tier = TierFromScore(cs.score)

	cs.lastCalculatedAt = time.Now().UTC()
	cs.updatedAt = time.Now().UTC()
}
