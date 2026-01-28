package command

import (
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

// CreateGig creates a new gig posting
type CreateGig struct {
	ClientID     string
	Title        string
	Description  string
	Category     string
	SkillID      string
	BudgetMin    int64
	BudgetMax    int64
	Currency     string
	DeliveryDays int
	Deadline     *time.Time
	IsRemote     bool
	Location     string
	Tags         []string
	Attachments  []string
}

// CreateGigResult is the result of creating a gig
type CreateGigResult struct {
	GigID        string    `json:"gig_id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	BudgetMin    int64     `json:"budget_min"`
	BudgetMax    int64     `json:"budget_max"`
	Currency     string    `json:"currency"`
	DeliveryDays int       `json:"delivery_days"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateGig updates an existing gig
type UpdateGig struct {
	GigID        string
	ClientID     string // for ownership verification
	Title        string
	Description  string
	Category     string
	SkillID      string
	BudgetMin    int64
	BudgetMax    int64
	DeliveryDays int
	Deadline     *time.Time
	IsRemote     bool
	Location     string
	Tags         []string
	Attachments  []string
}

// CancelGig cancels a gig
type CancelGig struct {
	GigID    string
	ClientID string
	Reason   string
}

// SubmitProposal submits a proposal for a gig
type SubmitProposal struct {
	GigID         string
	HustlerID     string
	CoverLetter   string
	ProposedPrice int64
	Currency      string
	DeliveryDays  int
	Attachments   []string
}

// SubmitProposalResult is the result of submitting a proposal
type SubmitProposalResult struct {
	ProposalID    string    `json:"proposal_id"`
	GigID         string    `json:"gig_id"`
	ProposedPrice int64     `json:"proposed_price"`
	DeliveryDays  int       `json:"delivery_days"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// WithdrawProposal withdraws a proposal
type WithdrawProposal struct {
	ProposalID string
	GigID      string
	HustlerID  string
}

// AcceptProposal accepts a proposal and creates a contract
type AcceptProposal struct {
	GigID      string
	ProposalID string
	ClientID   string
}

// AcceptProposalResult is the result of accepting a proposal
type AcceptProposalResult struct {
	ContractID   string    `json:"contract_id"`
	GigID        string    `json:"gig_id"`
	HustlerID    string    `json:"hustler_id"`
	AgreedPrice  int64     `json:"agreed_price"`
	PlatformFee  int64     `json:"platform_fee"`
	DeliveryDays int       `json:"delivery_days"`
	DeadlineAt   time.Time `json:"deadline_at"`
	Status       string    `json:"status"`
}

// DeliverWork marks work as delivered
type DeliverWork struct {
	ContractID   string
	HustlerID    string
	Deliverables []string
}

// DeliverWorkResult is the result of delivering work
type DeliverWorkResult struct {
	ContractID  string    `json:"contract_id"`
	Status      string    `json:"status"`
	DeliveredAt time.Time `json:"delivered_at"`
}

// ApproveDelivery approves delivered work
type ApproveDelivery struct {
	ContractID string
	ClientID   string
	Notes      string
}

// ApproveDeliveryResult is the result of approving delivery
type ApproveDeliveryResult struct {
	ContractID  string    `json:"contract_id"`
	Status      string    `json:"status"`
	CompletedAt time.Time `json:"completed_at"`
	PaidAmount  int64     `json:"paid_amount"`
}

// DisputeContract raises a dispute on a contract
type DisputeContract struct {
	ContractID string
	UserID     string
	Reason     string
}

// CancelContract cancels a contract
type CancelContract struct {
	ContractID string
	UserID     string
	Reason     string
}

// SubmitReview submits a review for a completed contract
type SubmitReview struct {
	ContractID          string
	ReviewerID          string
	Rating              int
	ReviewText          string
	CommunicationRating int
	QualityRating       int
	TimelinessRating    int
}

// SubmitReviewResult is the result of submitting a review
type SubmitReviewResult struct {
	ReviewID   string    `json:"review_id"`
	ContractID string    `json:"contract_id"`
	Rating     int       `json:"rating"`
	CreatedAt  time.Time `json:"created_at"`
}

// Helper methods for validation

func (c CreateGig) GetClientID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.ClientID)
}

func (c CreateGig) GetBudget() (valueobject.Money, valueobject.Money, error) {
	currency := valueobject.Currency(c.Currency)
	if c.Currency == "" {
		currency = valueobject.CurrencyNGN
	}

	min, err := valueobject.NewMoney(c.BudgetMin, currency)
	if err != nil {
		return valueobject.Money{}, valueobject.Money{}, err
	}

	max, err := valueobject.NewMoney(c.BudgetMax, currency)
	if err != nil {
		return valueobject.Money{}, valueobject.Money{}, err
	}

	return min, max, nil
}

func (c CreateGig) GetSkillID() (*valueobject.SkillID, error) {
	if c.SkillID == "" {
		return nil, nil
	}
	skillID, err := valueobject.NewSkillID(c.SkillID)
	if err != nil {
		return nil, err
	}
	return &skillID, nil
}

func (c SubmitProposal) GetGigID() (valueobject.GigID, error) {
	return valueobject.NewGigID(c.GigID)
}

func (c SubmitProposal) GetHustlerID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.HustlerID)
}

func (c SubmitProposal) GetProposedPrice() (valueobject.Money, error) {
	currency := valueobject.Currency(c.Currency)
	if c.Currency == "" {
		currency = valueobject.CurrencyNGN
	}
	return valueobject.NewMoney(c.ProposedPrice, currency)
}

func (c AcceptProposal) GetGigID() (valueobject.GigID, error) {
	return valueobject.NewGigID(c.GigID)
}

func (c AcceptProposal) GetProposalID() (valueobject.ProposalID, error) {
	return valueobject.NewProposalID(c.ProposalID)
}

func (c AcceptProposal) GetClientID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.ClientID)
}

func (c DeliverWork) GetContractID() (valueobject.ContractID, error) {
	return valueobject.NewContractID(c.ContractID)
}

func (c DeliverWork) GetHustlerID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.HustlerID)
}

func (c ApproveDelivery) GetContractID() (valueobject.ContractID, error) {
	return valueobject.NewContractID(c.ContractID)
}

func (c ApproveDelivery) GetClientID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.ClientID)
}

func (c SubmitReview) GetContractID() (valueobject.ContractID, error) {
	return valueobject.NewContractID(c.ContractID)
}

func (c SubmitReview) GetReviewerID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.ReviewerID)
}
