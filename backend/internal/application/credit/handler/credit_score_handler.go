package handler

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/application/credit/command"
	"hustlex/internal/domain/credit/aggregate"
	"hustlex/internal/domain/credit/repository"
)

// CreditScoreHandler handles credit score commands
type CreditScoreHandler struct {
	creditScoreRepo repository.CreditScoreRepository
}

// NewCreditScoreHandler creates a new credit score handler
func NewCreditScoreHandler(creditScoreRepo repository.CreditScoreRepository) *CreditScoreHandler {
	return &CreditScoreHandler{
		creditScoreRepo: creditScoreRepo,
	}
}

// HandleInitializeCreditScore creates a credit score for a new user
func (h *CreditScoreHandler) HandleInitializeCreditScore(ctx context.Context, cmd command.InitializeCreditScore) (*command.RecalculateCreditScoreResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check if already exists
	existing, _ := h.creditScoreRepo.FindByUserID(ctx, userID)
	if existing != nil {
		return &command.RecalculateCreditScoreResult{
			UserID:        existing.UserID().String(),
			Score:         existing.Score(),
			Tier:          existing.Tier().String(),
			MaxLoanAmount: existing.MaxLoanAmount(),
			InterestRate:  existing.InterestRate(),
		}, nil
	}

	// Create new credit score
	creditScore := aggregate.NewCreditScore(userID)

	if err := h.creditScoreRepo.Save(ctx, creditScore); err != nil {
		return nil, err
	}

	return &command.RecalculateCreditScoreResult{
		UserID:        creditScore.UserID().String(),
		Score:         creditScore.Score(),
		Tier:          creditScore.Tier().String(),
		MaxLoanAmount: creditScore.MaxLoanAmount(),
		InterestRate:  creditScore.InterestRate(),
	}, nil
}

// HandleUpdateCreditStats updates credit score components
func (h *CreditScoreHandler) HandleUpdateCreditStats(ctx context.Context, cmd command.UpdateCreditStats) (*command.RecalculateCreditScoreResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	creditScore, err := h.creditScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrCreditScoreNotFound
	}

	// Update gig stats
	if cmd.GigsCompleted != nil && cmd.GigsAccepted != nil {
		creditScore.UpdateGigStats(*cmd.GigsCompleted, *cmd.GigsAccepted)
	}

	// Update rating stats
	if cmd.AverageRating != nil && cmd.TotalReviews != nil {
		creditScore.UpdateRatingStats(*cmd.AverageRating, *cmd.TotalReviews)
	}

	// Update savings stats
	if cmd.OnTimeContributions != nil && cmd.TotalContributions != nil {
		creditScore.UpdateSavingsStats(*cmd.OnTimeContributions, *cmd.TotalContributions)
	}

	// Update account age
	if cmd.AccountAgeMonths != nil {
		duration := time.Duration(*cmd.AccountAgeMonths) * 30 * 24 * time.Hour
		creditScore.UpdateAccountAgeScore(duration)
	}

	// Update verification score
	if cmd.HasPhoneVerified != nil || cmd.HasEmailVerified != nil || cmd.HasBVNVerified != nil || cmd.HasNINVerified != nil {
		hasPhone := cmd.HasPhoneVerified != nil && *cmd.HasPhoneVerified
		hasEmail := cmd.HasEmailVerified != nil && *cmd.HasEmailVerified
		hasBVN := cmd.HasBVNVerified != nil && *cmd.HasBVNVerified
		hasNIN := cmd.HasNINVerified != nil && *cmd.HasNINVerified
		creditScore.UpdateVerificationScore(hasPhone, hasEmail, hasBVN, hasNIN)
	}

	// Update community score
	if cmd.CirclesJoined != nil || cmd.Referrals != nil {
		circles := 0
		referrals := 0
		if cmd.CirclesJoined != nil {
			circles = *cmd.CirclesJoined
		}
		if cmd.Referrals != nil {
			referrals = *cmd.Referrals
		}
		creditScore.UpdateCommunityScore(circles, referrals)
	}

	// Recalculate overall score
	creditScore.Recalculate()

	if err := h.creditScoreRepo.Save(ctx, creditScore); err != nil {
		return nil, err
	}

	return &command.RecalculateCreditScoreResult{
		UserID:        creditScore.UserID().String(),
		Score:         creditScore.Score(),
		Tier:          creditScore.Tier().String(),
		MaxLoanAmount: creditScore.MaxLoanAmount(),
		InterestRate:  creditScore.InterestRate(),
	}, nil
}

// HandleRecalculateCreditScore recalculates a user's credit score
func (h *CreditScoreHandler) HandleRecalculateCreditScore(ctx context.Context, cmd command.RecalculateCreditScore) (*command.RecalculateCreditScoreResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	creditScore, err := h.creditScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrCreditScoreNotFound
	}

	// Recalculate
	creditScore.Recalculate()

	if err := h.creditScoreRepo.Save(ctx, creditScore); err != nil {
		return nil, err
	}

	return &command.RecalculateCreditScoreResult{
		UserID:        creditScore.UserID().String(),
		Score:         creditScore.Score(),
		Tier:          creditScore.Tier().String(),
		MaxLoanAmount: creditScore.MaxLoanAmount(),
		InterestRate:  creditScore.InterestRate(),
	}, nil
}
