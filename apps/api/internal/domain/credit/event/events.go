package event

import (
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
)

// CreditScoreInitialized is raised when a credit score is created for a new user
type CreditScoreInitialized struct {
	sharedevent.BaseEvent
	UserID       string
	InitialScore int
	Tier         string
}

func NewCreditScoreInitialized(userID string, score int, tier string) CreditScoreInitialized {
	return CreditScoreInitialized{
		BaseEvent:    sharedevent.NewBaseEvent("credit.score.initialized"),
		UserID:       userID,
		InitialScore: score,
		Tier:         tier,
	}
}

// CreditScoreRecalculated is raised when a credit score is recalculated
type CreditScoreRecalculated struct {
	sharedevent.BaseEvent
	UserID       string
	OldScore     int
	NewScore     int
	OldTier      string
	NewTier      string
	TierChanged  bool
}

func NewCreditScoreRecalculated(userID string, oldScore, newScore int, oldTier, newTier string) CreditScoreRecalculated {
	return CreditScoreRecalculated{
		BaseEvent:   sharedevent.NewBaseEvent("credit.score.recalculated"),
		UserID:      userID,
		OldScore:    oldScore,
		NewScore:    newScore,
		OldTier:     oldTier,
		NewTier:     newTier,
		TierChanged: oldTier != newTier,
	}
}

// TierUpgraded is raised when a user's tier increases
type TierUpgraded struct {
	sharedevent.BaseEvent
	UserID         string
	OldTier        string
	NewTier        string
	NewMaxLoan     int64
	NewInterestRate float64
}

func NewTierUpgraded(userID, oldTier, newTier string, maxLoan int64, interestRate float64) TierUpgraded {
	return TierUpgraded{
		BaseEvent:       sharedevent.NewBaseEvent("credit.tier.upgraded"),
		UserID:          userID,
		OldTier:         oldTier,
		NewTier:         newTier,
		NewMaxLoan:      maxLoan,
		NewInterestRate: interestRate,
	}
}

// LoanApplied is raised when a user applies for a loan
type LoanApplied struct {
	sharedevent.BaseEvent
	LoanID       string
	UserID       string
	Principal    int64
	Currency     string
	InterestRate float64
	TenureMonths int
	Purpose      string
}

func NewLoanApplied(loanID, userID string, principal int64, currency string, rate float64, tenure int, purpose string) LoanApplied {
	return LoanApplied{
		BaseEvent:    sharedevent.NewBaseEvent("credit.loan.applied"),
		LoanID:       loanID,
		UserID:       userID,
		Principal:    principal,
		Currency:     currency,
		InterestRate: rate,
		TenureMonths: tenure,
		Purpose:      purpose,
	}
}

// LoanApproved is raised when a loan is approved
type LoanApproved struct {
	sharedevent.BaseEvent
	LoanID     string
	UserID     string
	Principal  int64
	Currency   string
	ApprovedAt time.Time
}

func NewLoanApproved(loanID, userID string, principal int64, currency string, approvedAt time.Time) LoanApproved {
	return LoanApproved{
		BaseEvent:  sharedevent.NewBaseEvent("credit.loan.approved"),
		LoanID:     loanID,
		UserID:     userID,
		Principal:  principal,
		Currency:   currency,
		ApprovedAt: approvedAt,
	}
}

// LoanRejected is raised when a loan is rejected
type LoanRejected struct {
	sharedevent.BaseEvent
	LoanID string
	UserID string
	Reason string
}

func NewLoanRejected(loanID, userID, reason string) LoanRejected {
	return LoanRejected{
		BaseEvent: sharedevent.NewBaseEvent("credit.loan.rejected"),
		LoanID:    loanID,
		UserID:    userID,
		Reason:    reason,
	}
}

// LoanDisbursed is raised when a loan is disbursed
type LoanDisbursed struct {
	sharedevent.BaseEvent
	LoanID      string
	UserID      string
	Amount      int64
	Currency    string
	DisbursedAt time.Time
	DueDate     time.Time
}

func NewLoanDisbursed(loanID, userID string, amount int64, currency string, disbursedAt, dueDate time.Time) LoanDisbursed {
	return LoanDisbursed{
		BaseEvent:   sharedevent.NewBaseEvent("credit.loan.disbursed"),
		LoanID:      loanID,
		UserID:      userID,
		Amount:      amount,
		Currency:    currency,
		DisbursedAt: disbursedAt,
		DueDate:     dueDate,
	}
}

// RepaymentRecorded is raised when a loan repayment is made
type RepaymentRecorded struct {
	sharedevent.BaseEvent
	LoanID           string
	UserID           string
	RepaymentID      string
	Amount           int64
	Currency         string
	RemainingBalance int64
	TransactionID    string
	PaidAt           time.Time
}

func NewRepaymentRecorded(loanID, userID, repaymentID string, amount int64, currency string, remaining int64, txnID string, paidAt time.Time) RepaymentRecorded {
	return RepaymentRecorded{
		BaseEvent:        sharedevent.NewBaseEvent("credit.repayment.recorded"),
		LoanID:           loanID,
		UserID:           userID,
		RepaymentID:      repaymentID,
		Amount:           amount,
		Currency:         currency,
		RemainingBalance: remaining,
		TransactionID:    txnID,
		PaidAt:           paidAt,
	}
}

// LoanCompleted is raised when a loan is fully repaid
type LoanCompleted struct {
	sharedevent.BaseEvent
	LoanID        string
	UserID        string
	TotalAmount   int64
	TotalRepaid   int64
	Currency      string
	CompletedAt   time.Time
	DaysToComplete int // Days from disbursement to completion
}

func NewLoanCompleted(loanID, userID string, totalAmount, totalRepaid int64, currency string, completedAt time.Time, days int) LoanCompleted {
	return LoanCompleted{
		BaseEvent:      sharedevent.NewBaseEvent("credit.loan.completed"),
		LoanID:         loanID,
		UserID:         userID,
		TotalAmount:    totalAmount,
		TotalRepaid:    totalRepaid,
		Currency:       currency,
		CompletedAt:    completedAt,
		DaysToComplete: days,
	}
}

// LoanDefaulted is raised when a loan is marked as defaulted
type LoanDefaulted struct {
	sharedevent.BaseEvent
	LoanID           string
	UserID           string
	OutstandingAmount int64
	Currency         string
	DaysOverdue      int
}

func NewLoanDefaulted(loanID, userID string, outstanding int64, currency string, daysOverdue int) LoanDefaulted {
	return LoanDefaulted{
		BaseEvent:         sharedevent.NewBaseEvent("credit.loan.defaulted"),
		LoanID:            loanID,
		UserID:            userID,
		OutstandingAmount: outstanding,
		Currency:          currency,
		DaysOverdue:       daysOverdue,
	}
}

// LoanOverdue is raised when a loan becomes overdue
type LoanOverdue struct {
	sharedevent.BaseEvent
	LoanID           string
	UserID           string
	OutstandingAmount int64
	Currency         string
	DueDate          time.Time
	DaysOverdue      int
}

func NewLoanOverdue(loanID, userID string, outstanding int64, currency string, dueDate time.Time, days int) LoanOverdue {
	return LoanOverdue{
		BaseEvent:         sharedevent.NewBaseEvent("credit.loan.overdue"),
		LoanID:            loanID,
		UserID:            userID,
		OutstandingAmount: outstanding,
		Currency:          currency,
		DueDate:           dueDate,
		DaysOverdue:       days,
	}
}
