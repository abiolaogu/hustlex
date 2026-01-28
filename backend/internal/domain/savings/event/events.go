package event

import (
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
)

const (
	AggregateTypeCircle = "Circle"
)

// CircleCreated is emitted when a savings circle is created
type CircleCreated struct {
	sharedevent.BaseEvent
	CircleID        string `json:"circle_id"`
	CreatorID       string `json:"creator_id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	ContributionAmt int64  `json:"contribution_amount"`
	Frequency       string `json:"frequency"`
	MaxMembers      int    `json:"max_members"`
}

func NewCircleCreated(circleID, creatorID, name, circleType string, contributionAmt int64, frequency string, maxMembers int) *CircleCreated {
	return &CircleCreated{
		BaseEvent: sharedevent.NewBaseEvent(
			"CircleCreated",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:        circleID,
		CreatorID:       creatorID,
		Name:            name,
		Type:            circleType,
		ContributionAmt: contributionAmt,
		Frequency:       frequency,
		MaxMembers:      maxMembers,
	}
}

// MemberJoined is emitted when a member joins a circle
type MemberJoined struct {
	sharedevent.BaseEvent
	CircleID string `json:"circle_id"`
	MemberID string `json:"member_id"`
	UserID   string `json:"user_id"`
	Position int    `json:"position"`
}

func NewMemberJoined(circleID, memberID, userID string, position int) *MemberJoined {
	return &MemberJoined{
		BaseEvent: sharedevent.NewBaseEvent(
			"MemberJoined",
			circleID,
			AggregateTypeCircle,
		),
		CircleID: circleID,
		MemberID: memberID,
		UserID:   userID,
		Position: position,
	}
}

// MemberLeft is emitted when a member leaves a circle
type MemberLeft struct {
	sharedevent.BaseEvent
	CircleID string `json:"circle_id"`
	MemberID string `json:"member_id"`
	UserID   string `json:"user_id"`
}

func NewMemberLeft(circleID, memberID, userID string) *MemberLeft {
	return &MemberLeft{
		BaseEvent: sharedevent.NewBaseEvent(
			"MemberLeft",
			circleID,
			AggregateTypeCircle,
		),
		CircleID: circleID,
		MemberID: memberID,
		UserID:   userID,
	}
}

// CircleStarted is emitted when a circle starts accepting contributions
type CircleStarted struct {
	sharedevent.BaseEvent
	CircleID   string    `json:"circle_id"`
	StartedAt  time.Time `json:"started_at"`
	FirstRound int       `json:"first_round"`
}

func NewCircleStarted(circleID string) *CircleStarted {
	return &CircleStarted{
		BaseEvent: sharedevent.NewBaseEvent(
			"CircleStarted",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:   circleID,
		StartedAt:  time.Now().UTC(),
		FirstRound: 1,
	}
}

// ContributionMade is emitted when a member makes a contribution
type ContributionMade struct {
	sharedevent.BaseEvent
	CircleID       string    `json:"circle_id"`
	ContributionID string    `json:"contribution_id"`
	MemberID       string    `json:"member_id"`
	Round          int       `json:"round"`
	Amount         int64     `json:"amount"`
	LateFee        int64     `json:"late_fee"`
	PaidAt         time.Time `json:"paid_at"`
}

func NewContributionMade(circleID, contributionID, memberID string, round int, amount, lateFee int64) *ContributionMade {
	return &ContributionMade{
		BaseEvent: sharedevent.NewBaseEvent(
			"ContributionMade",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:       circleID,
		ContributionID: contributionID,
		MemberID:       memberID,
		Round:          round,
		Amount:         amount,
		LateFee:        lateFee,
		PaidAt:         time.Now().UTC(),
	}
}

// PayoutTriggered is emitted when a payout is made to a member
type PayoutTriggered struct {
	sharedevent.BaseEvent
	CircleID    string `json:"circle_id"`
	RecipientID string `json:"recipient_id"`
	UserID      string `json:"user_id"`
	Round       int    `json:"round"`
	Amount      int64  `json:"amount"`
}

func NewPayoutTriggered(circleID, recipientID, userID string, round int, amount int64) *PayoutTriggered {
	return &PayoutTriggered{
		BaseEvent: sharedevent.NewBaseEvent(
			"PayoutTriggered",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:    circleID,
		RecipientID: recipientID,
		UserID:      userID,
		Round:       round,
		Amount:      amount,
	}
}

// RoundCompleted is emitted when a round is completed
type RoundCompleted struct {
	sharedevent.BaseEvent
	CircleID     string `json:"circle_id"`
	Round        int    `json:"round"`
	TotalCollected int64 `json:"total_collected"`
	NextRound    int    `json:"next_round"`
}

func NewRoundCompleted(circleID string, round int, totalCollected int64) *RoundCompleted {
	return &RoundCompleted{
		BaseEvent: sharedevent.NewBaseEvent(
			"RoundCompleted",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:       circleID,
		Round:          round,
		TotalCollected: totalCollected,
		NextRound:      round + 1,
	}
}

// CircleCompleted is emitted when all rounds are finished
type CircleCompleted struct {
	sharedevent.BaseEvent
	CircleID    string `json:"circle_id"`
	TotalRounds int    `json:"total_rounds"`
	TotalSaved  int64  `json:"total_saved"`
}

func NewCircleCompleted(circleID string, totalRounds int, totalSaved int64) *CircleCompleted {
	return &CircleCompleted{
		BaseEvent: sharedevent.NewBaseEvent(
			"CircleCompleted",
			circleID,
			AggregateTypeCircle,
		),
		CircleID:    circleID,
		TotalRounds: totalRounds,
		TotalSaved:  totalSaved,
	}
}
