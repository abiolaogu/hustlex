package aggregate

import (
	"testing"
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

func TestTierFromScore(t *testing.T) {
	tests := []struct {
		score int
		want  UserTier
	}{
		{0, TierBronze},
		{100, TierBronze},
		{399, TierBronze},
		{400, TierSilver},
		{500, TierSilver},
		{599, TierSilver},
		{600, TierGold},
		{700, TierGold},
		{749, TierGold},
		{750, TierPlatinum},
		{800, TierPlatinum},
		{850, TierPlatinum},
	}

	for _, tt := range tests {
		got := TierFromScore(tt.score)
		if got != tt.want {
			t.Errorf("TierFromScore(%d) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestUserTier_MaxLoanAmount(t *testing.T) {
	tests := []struct {
		tier UserTier
		want int64
	}{
		{TierBronze, 10000 * 100},   // ₦10,000
		{TierSilver, 50000 * 100},   // ₦50,000
		{TierGold, 200000 * 100},    // ₦200,000
		{TierPlatinum, 500000 * 100}, // ₦500,000
	}

	for _, tt := range tests {
		got := tt.tier.MaxLoanAmount()
		if got != tt.want {
			t.Errorf("%s.MaxLoanAmount() = %d, want %d", tt.tier, got, tt.want)
		}
	}
}

func TestUserTier_InterestRate(t *testing.T) {
	tests := []struct {
		tier UserTier
		want float64
	}{
		{TierBronze, 0.05},
		{TierSilver, 0.04},
		{TierGold, 0.03},
		{TierPlatinum, 0.02},
	}

	for _, tt := range tests {
		got := tt.tier.InterestRate()
		if got != tt.want {
			t.Errorf("%s.InterestRate() = %f, want %f", tt.tier, got, tt.want)
		}
	}
}

func TestNewCreditScore(t *testing.T) {
	userID := valueobject.GenerateUserID()

	cs := NewCreditScore(userID)

	if cs.ID() == "" {
		t.Error("NewCreditScore() should generate an ID")
	}
	if !cs.UserID().Equals(userID) {
		t.Error("NewCreditScore() should set userID")
	}
	if cs.Score() != 100 {
		t.Errorf("NewCreditScore() score = %d, want 100 (starting score)", cs.Score())
	}
	if cs.Tier() != TierBronze {
		t.Errorf("NewCreditScore() tier = %s, want bronze", cs.Tier())
	}
	if cs.Version() != 1 {
		t.Errorf("NewCreditScore() version = %d, want 1", cs.Version())
	}
}

func TestCreditScore_UpdateGigStats(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateGigStats(8, 10) // 80% completion rate

	if cs.TotalGigsCompleted() != 8 {
		t.Errorf("UpdateGigStats() completed = %d, want 8", cs.TotalGigsCompleted())
	}
	if cs.TotalGigsAccepted() != 10 {
		t.Errorf("UpdateGigStats() accepted = %d, want 10", cs.TotalGigsAccepted())
	}
	if cs.GigCompletionScore() != 80 {
		t.Errorf("UpdateGigStats() score = %d, want 80", cs.GigCompletionScore())
	}
}

func TestCreditScore_UpdateGigStats_PerfectCompletion(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateGigStats(10, 10) // 100% completion rate

	if cs.GigCompletionScore() != 100 {
		t.Errorf("UpdateGigStats() score = %d, want 100", cs.GigCompletionScore())
	}
}

func TestCreditScore_UpdateGigStats_ZeroAccepted(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateGigStats(0, 0) // No gigs

	if cs.GigCompletionScore() != 0 {
		t.Errorf("UpdateGigStats(0,0) score = %d, want 0", cs.GigCompletionScore())
	}
}

func TestCreditScore_UpdateRatingStats(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateRatingStats(4.5, 20) // 4.5 stars, 20 reviews

	if cs.AverageRating() != 4.5 {
		t.Errorf("UpdateRatingStats() rating = %f, want 4.5", cs.AverageRating())
	}
	if cs.TotalReviews() != 20 {
		t.Errorf("UpdateRatingStats() reviews = %d, want 20", cs.TotalReviews())
	}
	// 4.5 * 20 = 90
	if cs.RatingScore() != 90 {
		t.Errorf("UpdateRatingStats() score = %d, want 90", cs.RatingScore())
	}
}

func TestCreditScore_UpdateRatingStats_PerfectRating(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateRatingStats(5.0, 10)

	if cs.RatingScore() != 100 {
		t.Errorf("UpdateRatingStats(5.0) score = %d, want 100", cs.RatingScore())
	}
}

func TestCreditScore_UpdateSavingsStats(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	cs.UpdateSavingsStats(18, 20) // 90% on-time

	if cs.OnTimeContributions() != 18 {
		t.Errorf("UpdateSavingsStats() onTime = %d, want 18", cs.OnTimeContributions())
	}
	if cs.TotalContributions() != 20 {
		t.Errorf("UpdateSavingsStats() total = %d, want 20", cs.TotalContributions())
	}
	if cs.SavingsScore() != 90 {
		t.Errorf("UpdateSavingsStats() score = %d, want 90", cs.SavingsScore())
	}
}

func TestCreditScore_UpdateAccountAgeScore(t *testing.T) {
	tests := []struct {
		months    int
		wantScore int
	}{
		{0, 0},
		{1, 5},
		{6, 30},
		{12, 60},
		{20, 100}, // 100 is max
		{24, 100},
		{36, 100}, // Capped at 100
	}

	for _, tt := range tests {
		cs := NewCreditScore(valueobject.GenerateUserID())
		age := time.Duration(tt.months) * 30 * 24 * time.Hour

		cs.UpdateAccountAgeScore(age)

		if cs.AccountAgeScore() != tt.wantScore {
			t.Errorf("UpdateAccountAgeScore(%d months) = %d, want %d", tt.months, cs.AccountAgeScore(), tt.wantScore)
		}
	}
}

func TestCreditScore_UpdateVerificationScore(t *testing.T) {
	tests := []struct {
		phone, email, bvn, nin bool
		wantScore              int
	}{
		{false, false, false, false, 0},
		{true, false, false, false, 25},
		{true, true, false, false, 50},
		{true, true, true, false, 75},
		{true, true, true, true, 100},
		{false, false, true, true, 50},
	}

	for _, tt := range tests {
		cs := NewCreditScore(valueobject.GenerateUserID())

		cs.UpdateVerificationScore(tt.phone, tt.email, tt.bvn, tt.nin)

		if cs.VerificationScore() != tt.wantScore {
			t.Errorf("UpdateVerificationScore(%v,%v,%v,%v) = %d, want %d",
				tt.phone, tt.email, tt.bvn, tt.nin, cs.VerificationScore(), tt.wantScore)
		}
	}
}

func TestCreditScore_UpdateCommunityScore(t *testing.T) {
	tests := []struct {
		circles, referrals int
		wantScore          int
	}{
		{0, 0, 0},
		{1, 0, 10},
		{0, 2, 10},
		{2, 4, 40}, // 2*10 + 4*5 = 40
		{10, 10, 100}, // Capped at 100
		{20, 20, 100},
	}

	for _, tt := range tests {
		cs := NewCreditScore(valueobject.GenerateUserID())

		cs.UpdateCommunityScore(tt.circles, tt.referrals)

		if cs.CommunityScore() != tt.wantScore {
			t.Errorf("UpdateCommunityScore(%d,%d) = %d, want %d",
				tt.circles, tt.referrals, cs.CommunityScore(), tt.wantScore)
		}
	}
}

func TestCreditScore_Recalculate(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	// Set all components to 100
	cs.UpdateGigStats(10, 10)      // 100
	cs.UpdateRatingStats(5.0, 10)  // 100
	cs.UpdateSavingsStats(10, 10)  // 100
	cs.UpdateAccountAgeScore(24 * 30 * 24 * time.Hour) // 100
	cs.UpdateVerificationScore(true, true, true, true) // 100
	cs.UpdateCommunityScore(10, 10) // 100

	cs.Recalculate()

	// All 100 * 8.5 = 850 (max score)
	if cs.Score() != 850 {
		t.Errorf("Recalculate() with all 100s = %d, want 850", cs.Score())
	}
	if cs.Tier() != TierPlatinum {
		t.Errorf("Recalculate() tier = %s, want platinum", cs.Tier())
	}
}

func TestCreditScore_Recalculate_MixedScores(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	// Mixed scores
	cs.UpdateGigStats(8, 10)       // 80
	cs.UpdateRatingStats(4.0, 10)  // 80
	cs.UpdateSavingsStats(7, 10)   // 70
	cs.UpdateAccountAgeScore(6 * 30 * 24 * time.Hour) // 30
	cs.UpdateVerificationScore(true, true, false, false) // 50
	cs.UpdateCommunityScore(2, 2) // 30

	cs.Recalculate()

	// Weighted: 80*0.30 + 80*0.25 + 70*0.20 + 50*0.10 + 30*0.10 + 30*0.05
	// = 24 + 20 + 14 + 5 + 3 + 1.5 = 67.5
	// Scaled: 67.5 * 8.5 = 573.75 ≈ 573
	expectedScore := 573
	if cs.Score() != expectedScore {
		t.Errorf("Recalculate() with mixed scores = %d, want %d", cs.Score(), expectedScore)
	}
	if cs.Tier() != TierSilver {
		t.Errorf("Recalculate() tier = %s, want silver", cs.Tier())
	}
}

func TestCreditScore_Recalculate_ZeroScores(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	// All zero scores (default)
	cs.Recalculate()

	if cs.Score() != 0 {
		t.Errorf("Recalculate() with all zeros = %d, want 0", cs.Score())
	}
	if cs.Tier() != TierBronze {
		t.Errorf("Recalculate() tier = %s, want bronze", cs.Tier())
	}
}

func TestCreditScore_MaxLoanAmount(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	// Default tier is Bronze
	if cs.MaxLoanAmount() != TierBronze.MaxLoanAmount() {
		t.Errorf("MaxLoanAmount() = %d, want %d", cs.MaxLoanAmount(), TierBronze.MaxLoanAmount())
	}

	// After upgrading to Gold tier
	cs.UpdateGigStats(10, 10)
	cs.UpdateRatingStats(5.0, 10)
	cs.UpdateSavingsStats(10, 10)
	cs.UpdateAccountAgeScore(24 * 30 * 24 * time.Hour)
	cs.UpdateVerificationScore(true, true, false, false)
	cs.Recalculate()

	// Should be at least Gold tier now
	if cs.MaxLoanAmount() < TierGold.MaxLoanAmount() {
		t.Errorf("MaxLoanAmount() after upgrade = %d, want >= %d", cs.MaxLoanAmount(), TierGold.MaxLoanAmount())
	}
}

func TestCreditScore_InterestRate(t *testing.T) {
	cs := NewCreditScore(valueobject.GenerateUserID())

	// Default tier is Bronze
	if cs.InterestRate() != TierBronze.InterestRate() {
		t.Errorf("InterestRate() = %f, want %f", cs.InterestRate(), TierBronze.InterestRate())
	}
}

func TestReconstructCreditScore(t *testing.T) {
	userID := valueobject.GenerateUserID()
	now := time.Now().UTC()

	cs := ReconstructCreditScore(
		"cs-123",
		userID,
		650,
		TierGold,
		80, 90, 85, 50, 75, 40, // component scores
		15, 18, 4.5, 25, 20, 22, // stats
		now, now, now,
		5,
	)

	if cs.ID() != "cs-123" {
		t.Errorf("ID = %s, want cs-123", cs.ID())
	}
	if !cs.UserID().Equals(userID) {
		t.Error("UserID not set correctly")
	}
	if cs.Score() != 650 {
		t.Errorf("Score = %d, want 650", cs.Score())
	}
	if cs.Tier() != TierGold {
		t.Errorf("Tier = %s, want gold", cs.Tier())
	}
	if cs.GigCompletionScore() != 80 {
		t.Errorf("GigCompletionScore = %d, want 80", cs.GigCompletionScore())
	}
	if cs.RatingScore() != 90 {
		t.Errorf("RatingScore = %d, want 90", cs.RatingScore())
	}
	if cs.TotalGigsCompleted() != 15 {
		t.Errorf("TotalGigsCompleted = %d, want 15", cs.TotalGigsCompleted())
	}
	if cs.AverageRating() != 4.5 {
		t.Errorf("AverageRating = %f, want 4.5", cs.AverageRating())
	}
	if cs.Version() != 5 {
		t.Errorf("Version = %d, want 5", cs.Version())
	}
}

func TestUserTier_String(t *testing.T) {
	tests := []struct {
		tier UserTier
		want string
	}{
		{TierBronze, "bronze"},
		{TierSilver, "silver"},
		{TierGold, "gold"},
		{TierPlatinum, "platinum"},
	}

	for _, tt := range tests {
		if got := tt.tier.String(); got != tt.want {
			t.Errorf("%v.String() = %s, want %s", tt.tier, got, tt.want)
		}
	}
}
