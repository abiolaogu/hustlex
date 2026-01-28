package aggregate

import (
	"errors"
	"time"

	"hustlex/internal/domain/gig/event"
	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Contract errors
var (
	ErrContractNotActive    = errors.New("contract is not active")
	ErrContractNotDelivered = errors.New("contract has not been delivered")
	ErrNotContractParty     = errors.New("not a party to this contract")
	ErrAlreadyReviewed      = errors.New("already submitted a review for this contract")
	ErrCannotDeliver        = errors.New("cannot deliver in current state")
	ErrCannotApprove        = errors.New("cannot approve in current state")
	ErrInvalidRating        = errors.New("rating must be between 1 and 5")
)

// ContractStatus represents the status of a contract
type ContractStatus string

const (
	ContractStatusActive    ContractStatus = "active"
	ContractStatusDelivered ContractStatus = "delivered"
	ContractStatusCompleted ContractStatus = "completed"
	ContractStatusDisputed  ContractStatus = "disputed"
	ContractStatusCancelled ContractStatus = "cancelled"
)

func (s ContractStatus) String() string {
	return string(s)
}

func (s ContractStatus) IsActive() bool {
	return s == ContractStatusActive
}

func (s ContractStatus) IsDelivered() bool {
	return s == ContractStatusDelivered
}

func (s ContractStatus) IsCompleted() bool {
	return s == ContractStatusCompleted
}

// Review represents a review for a completed contract
type Review struct {
	id                  string
	reviewerID          valueobject.UserID
	revieweeID          valueobject.UserID
	rating              int
	reviewText          string
	communicationRating int
	qualityRating       int
	timelinessRating    int
	isPublic            bool
	createdAt           time.Time
}

func NewReview(
	id string,
	reviewerID valueobject.UserID,
	revieweeID valueobject.UserID,
	rating int,
	reviewText string,
) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	return &Review{
		id:         id,
		reviewerID: reviewerID,
		revieweeID: revieweeID,
		rating:     rating,
		reviewText: reviewText,
		isPublic:   true,
		createdAt:  time.Now().UTC(),
	}, nil
}

func (r *Review) ID() string                   { return r.id }
func (r *Review) ReviewerID() valueobject.UserID { return r.reviewerID }
func (r *Review) RevieweeID() valueobject.UserID { return r.revieweeID }
func (r *Review) Rating() int                  { return r.rating }
func (r *Review) ReviewText() string           { return r.reviewText }
func (r *Review) CommunicationRating() int     { return r.communicationRating }
func (r *Review) QualityRating() int           { return r.qualityRating }
func (r *Review) TimelinessRating() int        { return r.timelinessRating }
func (r *Review) IsPublic() bool               { return r.isPublic }
func (r *Review) CreatedAt() time.Time         { return r.createdAt }

func (r *Review) SetDetailedRatings(communication, quality, timeliness int) {
	if communication >= 1 && communication <= 5 {
		r.communicationRating = communication
	}
	if quality >= 1 && quality <= 5 {
		r.qualityRating = quality
	}
	if timeliness >= 1 && timeliness <= 5 {
		r.timelinessRating = timeliness
	}
}

// Contract is the aggregate root for gig contracts
type Contract struct {
	sharedevent.AggregateRoot

	id           valueobject.ContractID
	gigID        valueobject.GigID
	proposalID   valueobject.ProposalID
	clientID     valueobject.UserID
	hustlerID    valueobject.UserID
	agreedPrice  valueobject.Money
	platformFee  valueobject.Money
	deliveryDays int
	status       ContractStatus
	startedAt    time.Time
	deadlineAt   time.Time
	deliveredAt  *time.Time
	completedAt  *time.Time
	deliverables []string
	clientNotes  string
	reviews      []*Review
	createdAt    time.Time
	updatedAt    time.Time
	version      int64
}

// NewContract creates a new contract from accepted proposal data
func NewContract(data *AcceptedProposalData, proposalID valueobject.ProposalID) (*Contract, error) {
	platformFee, err := valueobject.NewMoney(data.PlatformFee, data.AgreedPrice.Currency())
	if err != nil {
		return nil, err
	}

	contract := &Contract{
		id:           data.ContractID,
		gigID:        data.GigID,
		proposalID:   proposalID,
		clientID:     data.ClientID,
		hustlerID:    data.HustlerID,
		agreedPrice:  data.AgreedPrice,
		platformFee:  platformFee,
		deliveryDays: data.DeliveryDays,
		status:       ContractStatusActive,
		startedAt:    time.Now().UTC(),
		deadlineAt:   data.DeadlineAt,
		deliverables: make([]string, 0),
		reviews:      make([]*Review, 0),
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
		version:      1,
	}

	contract.RecordEvent(event.NewContractCreated(
		data.ContractID.String(),
		data.GigID.String(),
		data.ClientID.String(),
		data.HustlerID.String(),
		data.AgreedPrice.Amount(),
		data.PlatformFee,
		data.DeliveryDays,
		data.DeadlineAt,
	))

	return contract, nil
}

// ReconstructContract reconstructs a contract from persistence
func ReconstructContract(
	id valueobject.ContractID,
	gigID valueobject.GigID,
	proposalID valueobject.ProposalID,
	clientID valueobject.UserID,
	hustlerID valueobject.UserID,
	agreedPrice valueobject.Money,
	platformFee valueobject.Money,
	deliveryDays int,
	status ContractStatus,
	startedAt time.Time,
	deadlineAt time.Time,
	deliveredAt *time.Time,
	completedAt *time.Time,
	deliverables []string,
	clientNotes string,
	reviews []*Review,
	createdAt time.Time,
	updatedAt time.Time,
	version int64,
) *Contract {
	return &Contract{
		id:           id,
		gigID:        gigID,
		proposalID:   proposalID,
		clientID:     clientID,
		hustlerID:    hustlerID,
		agreedPrice:  agreedPrice,
		platformFee:  platformFee,
		deliveryDays: deliveryDays,
		status:       status,
		startedAt:    startedAt,
		deadlineAt:   deadlineAt,
		deliveredAt:  deliveredAt,
		completedAt:  completedAt,
		deliverables: deliverables,
		clientNotes:  clientNotes,
		reviews:      reviews,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		version:      version,
	}
}

// Getters
func (c *Contract) ID() valueobject.ContractID    { return c.id }
func (c *Contract) GigID() valueobject.GigID      { return c.gigID }
func (c *Contract) ProposalID() valueobject.ProposalID { return c.proposalID }
func (c *Contract) ClientID() valueobject.UserID  { return c.clientID }
func (c *Contract) HustlerID() valueobject.UserID { return c.hustlerID }
func (c *Contract) AgreedPrice() valueobject.Money { return c.agreedPrice }
func (c *Contract) PlatformFee() valueobject.Money { return c.platformFee }
func (c *Contract) DeliveryDays() int             { return c.deliveryDays }
func (c *Contract) Status() ContractStatus        { return c.status }
func (c *Contract) StartedAt() time.Time          { return c.startedAt }
func (c *Contract) DeadlineAt() time.Time         { return c.deadlineAt }
func (c *Contract) DeliveredAt() *time.Time       { return c.deliveredAt }
func (c *Contract) CompletedAt() *time.Time       { return c.completedAt }
func (c *Contract) Deliverables() []string        { return c.deliverables }
func (c *Contract) ClientNotes() string           { return c.clientNotes }
func (c *Contract) Reviews() []*Review            { return c.reviews }
func (c *Contract) CreatedAt() time.Time          { return c.createdAt }
func (c *Contract) UpdatedAt() time.Time          { return c.updatedAt }
func (c *Contract) Version() int64                { return c.version }

func (c *Contract) IsActive() bool    { return c.status.IsActive() }
func (c *Contract) IsDelivered() bool { return c.status.IsDelivered() }
func (c *Contract) IsCompleted() bool { return c.status.IsCompleted() }

// IsParty checks if the user is a party to this contract
func (c *Contract) IsParty(userID valueobject.UserID) bool {
	return c.clientID.Equals(userID) || c.hustlerID.Equals(userID)
}

// IsOverdue checks if the contract is past its deadline
func (c *Contract) IsOverdue() bool {
	return time.Now().After(c.deadlineAt) && c.status.IsActive()
}

// NetPayoutAmount returns the amount the hustler will receive (agreed - platform fee)
func (c *Contract) NetPayoutAmount() valueobject.Money {
	return c.agreedPrice.MustSubtract(c.platformFee)
}

// Business Methods

// Deliver marks the contract as delivered
func (c *Contract) Deliver(hustlerID valueobject.UserID, deliverables []string) error {
	if !c.hustlerID.Equals(hustlerID) {
		return ErrNotContractParty
	}

	if !c.status.IsActive() {
		return ErrCannotDeliver
	}

	now := time.Now().UTC()
	c.status = ContractStatusDelivered
	c.deliveredAt = &now
	c.deliverables = deliverables
	c.updatedAt = now

	c.RecordEvent(event.NewWorkDelivered(c.id.String(), hustlerID.String(), deliverables))

	return nil
}

// Approve approves the delivered work and completes the contract
func (c *Contract) Approve(clientID valueobject.UserID, notes string) error {
	if !c.clientID.Equals(clientID) {
		return ErrNotContractParty
	}

	if !c.status.IsDelivered() {
		return ErrCannotApprove
	}

	now := time.Now().UTC()
	c.status = ContractStatusCompleted
	c.completedAt = &now
	c.clientNotes = notes
	c.updatedAt = now

	c.RecordEvent(event.NewWorkApproved(
		c.id.String(),
		clientID.String(),
		c.hustlerID.String(),
		c.NetPayoutAmount().Amount(),
	))

	return nil
}

// Dispute marks the contract as disputed
func (c *Contract) Dispute(userID valueobject.UserID, reason string) error {
	if !c.IsParty(userID) {
		return ErrNotContractParty
	}

	if c.status.IsCompleted() {
		return errors.New("cannot dispute completed contract")
	}

	c.status = ContractStatusDisputed
	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewContractDisputed(c.id.String(), userID.String(), reason))

	return nil
}

// Cancel cancels the contract
func (c *Contract) Cancel(userID valueobject.UserID, reason string, refundAmount int64) error {
	if !c.IsParty(userID) {
		return ErrNotContractParty
	}

	if c.status.IsCompleted() {
		return errors.New("cannot cancel completed contract")
	}

	c.status = ContractStatusCancelled
	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewContractCancelled(c.id.String(), userID.String(), reason, refundAmount))

	return nil
}

// AddReview adds a review to the contract
func (c *Contract) AddReview(review *Review) error {
	if !c.status.IsCompleted() {
		return errors.New("can only review completed contracts")
	}

	// Check if user has already reviewed
	for _, r := range c.reviews {
		if r.ReviewerID().Equals(review.ReviewerID()) {
			return ErrAlreadyReviewed
		}
	}

	// Validate reviewer is a party to the contract
	if !c.IsParty(review.ReviewerID()) {
		return ErrNotContractParty
	}

	c.reviews = append(c.reviews, review)
	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewReviewSubmitted(
		review.ID(),
		c.id.String(),
		review.ReviewerID().String(),
		review.RevieweeID().String(),
		review.Rating(),
	))

	return nil
}

// HasReviewFrom checks if a user has already submitted a review
func (c *Contract) HasReviewFrom(userID valueobject.UserID) bool {
	for _, r := range c.reviews {
		if r.ReviewerID().Equals(userID) {
			return true
		}
	}
	return false
}

// GetReviewBy returns the review submitted by a specific user
func (c *Contract) GetReviewBy(userID valueobject.UserID) *Review {
	for _, r := range c.reviews {
		if r.ReviewerID().Equals(userID) {
			return r
		}
	}
	return nil
}
