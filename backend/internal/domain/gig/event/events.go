package event

import (
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
)

const (
	AggregateTypeGig      = "Gig"
	AggregateTypeContract = "Contract"
)

// GigPosted is emitted when a new gig is posted
type GigPosted struct {
	sharedevent.BaseEvent
	GigID       string `json:"gig_id"`
	ClientID    string `json:"client_id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	BudgetMin   int64  `json:"budget_min"`
	BudgetMax   int64  `json:"budget_max"`
	IsRemote    bool   `json:"is_remote"`
}

func NewGigPosted(gigID, clientID, title, category string, budgetMin, budgetMax int64, isRemote bool) *GigPosted {
	return &GigPosted{
		BaseEvent: sharedevent.NewBaseEvent(
			"GigPosted",
			gigID,
			AggregateTypeGig,
		),
		GigID:     gigID,
		ClientID:  clientID,
		Title:     title,
		Category:  category,
		BudgetMin: budgetMin,
		BudgetMax: budgetMax,
		IsRemote:  isRemote,
	}
}

// GigUpdated is emitted when a gig is updated
type GigUpdated struct {
	sharedevent.BaseEvent
	GigID         string            `json:"gig_id"`
	UpdatedFields map[string]string `json:"updated_fields"`
}

func NewGigUpdated(gigID string, updatedFields map[string]string) *GigUpdated {
	return &GigUpdated{
		BaseEvent: sharedevent.NewBaseEvent(
			"GigUpdated",
			gigID,
			AggregateTypeGig,
		),
		GigID:         gigID,
		UpdatedFields: updatedFields,
	}
}

// GigCancelled is emitted when a gig is cancelled
type GigCancelled struct {
	sharedevent.BaseEvent
	GigID    string `json:"gig_id"`
	ClientID string `json:"client_id"`
	Reason   string `json:"reason,omitempty"`
}

func NewGigCancelled(gigID, clientID, reason string) *GigCancelled {
	return &GigCancelled{
		BaseEvent: sharedevent.NewBaseEvent(
			"GigCancelled",
			gigID,
			AggregateTypeGig,
		),
		GigID:    gigID,
		ClientID: clientID,
		Reason:   reason,
	}
}

// ProposalSubmitted is emitted when a proposal is submitted
type ProposalSubmitted struct {
	sharedevent.BaseEvent
	ProposalID    string `json:"proposal_id"`
	GigID         string `json:"gig_id"`
	HustlerID     string `json:"hustler_id"`
	ProposedPrice int64  `json:"proposed_price"`
	DeliveryDays  int    `json:"delivery_days"`
}

func NewProposalSubmitted(proposalID, gigID, hustlerID string, proposedPrice int64, deliveryDays int) *ProposalSubmitted {
	return &ProposalSubmitted{
		BaseEvent: sharedevent.NewBaseEvent(
			"ProposalSubmitted",
			gigID,
			AggregateTypeGig,
		),
		ProposalID:    proposalID,
		GigID:         gigID,
		HustlerID:     hustlerID,
		ProposedPrice: proposedPrice,
		DeliveryDays:  deliveryDays,
	}
}

// ProposalWithdrawn is emitted when a proposal is withdrawn
type ProposalWithdrawn struct {
	sharedevent.BaseEvent
	ProposalID string `json:"proposal_id"`
	GigID      string `json:"gig_id"`
	HustlerID  string `json:"hustler_id"`
}

func NewProposalWithdrawn(proposalID, gigID, hustlerID string) *ProposalWithdrawn {
	return &ProposalWithdrawn{
		BaseEvent: sharedevent.NewBaseEvent(
			"ProposalWithdrawn",
			gigID,
			AggregateTypeGig,
		),
		ProposalID: proposalID,
		GigID:      gigID,
		HustlerID:  hustlerID,
	}
}

// ProposalAccepted is emitted when a proposal is accepted
type ProposalAccepted struct {
	sharedevent.BaseEvent
	ProposalID string `json:"proposal_id"`
	GigID      string `json:"gig_id"`
	ContractID string `json:"contract_id"`
	HustlerID  string `json:"hustler_id"`
	ClientID   string `json:"client_id"`
	Amount     int64  `json:"amount"`
}

func NewProposalAccepted(proposalID, gigID, contractID, hustlerID, clientID string, amount int64) *ProposalAccepted {
	return &ProposalAccepted{
		BaseEvent: sharedevent.NewBaseEvent(
			"ProposalAccepted",
			gigID,
			AggregateTypeGig,
		),
		ProposalID: proposalID,
		GigID:      gigID,
		ContractID: contractID,
		HustlerID:  hustlerID,
		ClientID:   clientID,
		Amount:     amount,
	}
}

// ContractCreated is emitted when a contract is created from an accepted proposal
type ContractCreated struct {
	sharedevent.BaseEvent
	ContractID   string    `json:"contract_id"`
	GigID        string    `json:"gig_id"`
	ClientID     string    `json:"client_id"`
	HustlerID    string    `json:"hustler_id"`
	AgreedPrice  int64     `json:"agreed_price"`
	PlatformFee  int64     `json:"platform_fee"`
	DeliveryDays int       `json:"delivery_days"`
	DeadlineAt   time.Time `json:"deadline_at"`
}

func NewContractCreated(contractID, gigID, clientID, hustlerID string, agreedPrice, platformFee int64, deliveryDays int, deadlineAt time.Time) *ContractCreated {
	return &ContractCreated{
		BaseEvent: sharedevent.NewBaseEvent(
			"ContractCreated",
			contractID,
			AggregateTypeContract,
		),
		ContractID:   contractID,
		GigID:        gigID,
		ClientID:     clientID,
		HustlerID:    hustlerID,
		AgreedPrice:  agreedPrice,
		PlatformFee:  platformFee,
		DeliveryDays: deliveryDays,
		DeadlineAt:   deadlineAt,
	}
}

// WorkDelivered is emitted when work is delivered on a contract
type WorkDelivered struct {
	sharedevent.BaseEvent
	ContractID   string    `json:"contract_id"`
	HustlerID    string    `json:"hustler_id"`
	DeliveredAt  time.Time `json:"delivered_at"`
	Deliverables []string  `json:"deliverables,omitempty"`
}

func NewWorkDelivered(contractID, hustlerID string, deliverables []string) *WorkDelivered {
	return &WorkDelivered{
		BaseEvent: sharedevent.NewBaseEvent(
			"WorkDelivered",
			contractID,
			AggregateTypeContract,
		),
		ContractID:   contractID,
		HustlerID:    hustlerID,
		DeliveredAt:  time.Now().UTC(),
		Deliverables: deliverables,
	}
}

// WorkApproved is emitted when client approves delivered work
type WorkApproved struct {
	sharedevent.BaseEvent
	ContractID  string    `json:"contract_id"`
	ClientID    string    `json:"client_id"`
	HustlerID   string    `json:"hustler_id"`
	ApprovedAt  time.Time `json:"approved_at"`
	AmountToRelease int64 `json:"amount_to_release"`
}

func NewWorkApproved(contractID, clientID, hustlerID string, amountToRelease int64) *WorkApproved {
	return &WorkApproved{
		BaseEvent: sharedevent.NewBaseEvent(
			"WorkApproved",
			contractID,
			AggregateTypeContract,
		),
		ContractID:      contractID,
		ClientID:        clientID,
		HustlerID:       hustlerID,
		ApprovedAt:      time.Now().UTC(),
		AmountToRelease: amountToRelease,
	}
}

// ContractDisputed is emitted when a contract is disputed
type ContractDisputed struct {
	sharedevent.BaseEvent
	ContractID string `json:"contract_id"`
	DisputedBy string `json:"disputed_by"` // user_id
	Reason     string `json:"reason"`
}

func NewContractDisputed(contractID, disputedBy, reason string) *ContractDisputed {
	return &ContractDisputed{
		BaseEvent: sharedevent.NewBaseEvent(
			"ContractDisputed",
			contractID,
			AggregateTypeContract,
		),
		ContractID: contractID,
		DisputedBy: disputedBy,
		Reason:     reason,
	}
}

// ContractCancelled is emitted when a contract is cancelled
type ContractCancelled struct {
	sharedevent.BaseEvent
	ContractID  string `json:"contract_id"`
	CancelledBy string `json:"cancelled_by"`
	Reason      string `json:"reason"`
	RefundAmount int64 `json:"refund_amount,omitempty"`
}

func NewContractCancelled(contractID, cancelledBy, reason string, refundAmount int64) *ContractCancelled {
	return &ContractCancelled{
		BaseEvent: sharedevent.NewBaseEvent(
			"ContractCancelled",
			contractID,
			AggregateTypeContract,
		),
		ContractID:   contractID,
		CancelledBy:  cancelledBy,
		Reason:       reason,
		RefundAmount: refundAmount,
	}
}

// ReviewSubmitted is emitted when a review is submitted
type ReviewSubmitted struct {
	sharedevent.BaseEvent
	ReviewID   string `json:"review_id"`
	ContractID string `json:"contract_id"`
	ReviewerID string `json:"reviewer_id"`
	RevieweeID string `json:"reviewee_id"`
	Rating     int    `json:"rating"`
}

func NewReviewSubmitted(reviewID, contractID, reviewerID, revieweeID string, rating int) *ReviewSubmitted {
	return &ReviewSubmitted{
		BaseEvent: sharedevent.NewBaseEvent(
			"ReviewSubmitted",
			contractID,
			AggregateTypeContract,
		),
		ReviewID:   reviewID,
		ContractID: contractID,
		ReviewerID: reviewerID,
		RevieweeID: revieweeID,
		Rating:     rating,
	}
}
