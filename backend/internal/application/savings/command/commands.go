package command

import (
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

// CreateCircle creates a new savings circle
type CreateCircle struct {
	CreatorID       string
	Name            string
	Description     string
	Type            string // rotational, fixed_target, emergency
	ContributionAmt int64
	Currency        string
	Frequency       string // daily, weekly, biweekly, monthly
	MaxMembers      int
	TotalRounds     int
	StartDate       *time.Time
	IsPrivate       bool
	Rules           []string
}

// CreateCircleResult is the result of creating a circle
type CreateCircleResult struct {
	CircleID        string    `json:"circle_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	ContributionAmt int64     `json:"contribution_amount"`
	Currency        string    `json:"currency"`
	Frequency       string    `json:"frequency"`
	MaxMembers      int       `json:"max_members"`
	InviteCode      string    `json:"invite_code"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// JoinCircle joins a savings circle
type JoinCircle struct {
	CircleID string
	UserID   string
}

// JoinCircleByCode joins a circle using invite code
type JoinCircleByCode struct {
	InviteCode string
	UserID     string
}

// JoinCircleResult is the result of joining a circle
type JoinCircleResult struct {
	MemberID string `json:"member_id"`
	CircleID string `json:"circle_id"`
	Position int    `json:"position"`
	Status   string `json:"status"`
}

// LeaveCircle leaves a savings circle
type LeaveCircle struct {
	CircleID string
	UserID   string
}

// StartCircle manually starts a circle (admin only)
type StartCircle struct {
	CircleID string
	AdminID  string
}

// MakeContribution makes a contribution to the circle
type MakeContribution struct {
	CircleID      string
	UserID        string
	TransactionID string // from wallet
}

// MakeContributionResult is the result of making a contribution
type MakeContributionResult struct {
	ContributionID string    `json:"contribution_id"`
	CircleID       string    `json:"circle_id"`
	Round          int       `json:"round"`
	Amount         int64     `json:"amount"`
	LateFee        int64     `json:"late_fee"`
	TotalPaid      int64     `json:"total_paid"`
	PaidAt         time.Time `json:"paid_at"`
}

// TriggerPayout triggers a payout to the current recipient
type TriggerPayout struct {
	CircleID string
}

// Helper methods

func (c CreateCircle) GetCreatorID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.CreatorID)
}

func (c CreateCircle) GetContributionAmount() (valueobject.Money, error) {
	currency := valueobject.Currency(c.Currency)
	if c.Currency == "" {
		currency = valueobject.CurrencyNGN
	}
	return valueobject.NewMoney(c.ContributionAmt, currency)
}

func (c JoinCircle) GetCircleID() (valueobject.CircleID, error) {
	return valueobject.NewCircleID(c.CircleID)
}

func (c JoinCircle) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c LeaveCircle) GetCircleID() (valueobject.CircleID, error) {
	return valueobject.NewCircleID(c.CircleID)
}

func (c LeaveCircle) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c MakeContribution) GetCircleID() (valueobject.CircleID, error) {
	return valueobject.NewCircleID(c.CircleID)
}

func (c MakeContribution) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c MakeContribution) GetTransactionID() (valueobject.TransactionID, error) {
	return valueobject.NewTransactionID(c.TransactionID)
}
