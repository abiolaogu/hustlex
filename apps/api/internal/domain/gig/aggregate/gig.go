package aggregate

import (
	"errors"
	"time"

	"hustlex/internal/domain/gig/event"
	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrGigNotOpen           = errors.New("gig is not accepting proposals")
	ErrGigAlreadyInProgress = errors.New("gig is already in progress")
	ErrCannotUpdateGig      = errors.New("cannot update gig in current status")
	ErrNotGigOwner          = errors.New("only gig owner can perform this action")
	ErrCannotProposeSelf    = errors.New("cannot submit proposal to own gig")
	ErrAlreadyProposed      = errors.New("already submitted a proposal for this gig")
	ErrProposalNotFound     = errors.New("proposal not found")
	ErrProposalNotPending   = errors.New("proposal is no longer pending")
	ErrInvalidBudget        = errors.New("invalid budget range")
	ErrPriceBelowBudget     = errors.New("proposed price is below minimum budget")
	ErrPriceAboveBudget     = errors.New("proposed price exceeds maximum budget")
)

// GigStatus represents the status of a gig
type GigStatus string

const (
	GigStatusOpen       GigStatus = "open"
	GigStatusInProgress GigStatus = "in_progress"
	GigStatusCompleted  GigStatus = "completed"
	GigStatusCancelled  GigStatus = "cancelled"
	GigStatusDisputed   GigStatus = "disputed"
)

func (s GigStatus) String() string {
	return string(s)
}

func (s GigStatus) IsOpen() bool {
	return s == GigStatusOpen
}

// ProposalStatus represents the status of a proposal
type ProposalStatus string

const (
	ProposalStatusPending   ProposalStatus = "pending"
	ProposalStatusAccepted  ProposalStatus = "accepted"
	ProposalStatusRejected  ProposalStatus = "rejected"
	ProposalStatusWithdrawn ProposalStatus = "withdrawn"
)

// Budget represents a gig's budget range
type Budget struct {
	min      valueobject.Money
	max      valueobject.Money
}

func NewBudget(min, max valueobject.Money) (Budget, error) {
	if min.GreaterThan(max) {
		return Budget{}, ErrInvalidBudget
	}
	return Budget{min: min, max: max}, nil
}

func (b Budget) Min() valueobject.Money { return b.min }
func (b Budget) Max() valueobject.Money { return b.max }

func (b Budget) Contains(amount valueobject.Money) bool {
	return !amount.LessThan(b.min) && !amount.GreaterThan(b.max)
}

// Proposal represents a proposal entity within the Gig aggregate
type Proposal struct {
	id            valueobject.ProposalID
	hustlerID     valueobject.UserID
	coverLetter   string
	proposedPrice valueobject.Money
	deliveryDays  int
	status        ProposalStatus
	attachments   []string
	createdAt     time.Time
	updatedAt     time.Time
}

func NewProposal(
	id valueobject.ProposalID,
	hustlerID valueobject.UserID,
	coverLetter string,
	proposedPrice valueobject.Money,
	deliveryDays int,
	attachments []string,
) *Proposal {
	return &Proposal{
		id:            id,
		hustlerID:     hustlerID,
		coverLetter:   coverLetter,
		proposedPrice: proposedPrice,
		deliveryDays:  deliveryDays,
		status:        ProposalStatusPending,
		attachments:   attachments,
		createdAt:     time.Now().UTC(),
		updatedAt:     time.Now().UTC(),
	}
}

func (p *Proposal) ID() valueobject.ProposalID   { return p.id }
func (p *Proposal) HustlerID() valueobject.UserID { return p.hustlerID }
func (p *Proposal) CoverLetter() string           { return p.coverLetter }
func (p *Proposal) ProposedPrice() valueobject.Money { return p.proposedPrice }
func (p *Proposal) DeliveryDays() int             { return p.deliveryDays }
func (p *Proposal) Status() ProposalStatus        { return p.status }
func (p *Proposal) Attachments() []string         { return p.attachments }
func (p *Proposal) CreatedAt() time.Time          { return p.createdAt }
func (p *Proposal) UpdatedAt() time.Time          { return p.updatedAt }
func (p *Proposal) IsPending() bool               { return p.status == ProposalStatusPending }

func (p *Proposal) Accept() {
	p.status = ProposalStatusAccepted
	p.updatedAt = time.Now().UTC()
}

func (p *Proposal) Reject() {
	p.status = ProposalStatusRejected
	p.updatedAt = time.Now().UTC()
}

func (p *Proposal) Withdraw() {
	p.status = ProposalStatusWithdrawn
	p.updatedAt = time.Now().UTC()
}

// Gig is the aggregate root for gig postings
type Gig struct {
	sharedevent.AggregateRoot

	id            valueobject.GigID
	clientID      valueobject.UserID
	title         string
	description   string
	category      string
	skillID       *valueobject.SkillID
	budget        Budget
	currency      valueobject.Currency
	deliveryDays  int
	deadline      *time.Time
	isRemote      bool
	location      string
	status        GigStatus
	viewCount     int
	isFeatured    bool
	attachments   []string
	tags          []string
	proposals     []*Proposal
	acceptedProposalID *valueobject.ProposalID
	createdAt     time.Time
	updatedAt     time.Time
	version       int64
}

// NewGig creates a new Gig aggregate
func NewGig(
	id valueobject.GigID,
	clientID valueobject.UserID,
	title string,
	description string,
	category string,
	budget Budget,
	deliveryDays int,
	isRemote bool,
) (*Gig, error) {
	gig := &Gig{
		id:           id,
		clientID:     clientID,
		title:        title,
		description:  description,
		category:     category,
		budget:       budget,
		currency:     valueobject.NGN,
		deliveryDays: deliveryDays,
		isRemote:     isRemote,
		status:       GigStatusOpen,
		viewCount:    0,
		isFeatured:   false,
		proposals:    make([]*Proposal, 0),
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
		version:      1,
	}

	gig.RecordEvent(event.NewGigPosted(
		id.String(),
		clientID.String(),
		title,
		category,
		budget.Min().Amount(),
		budget.Max().Amount(),
		isRemote,
	))

	return gig, nil
}

// ReconstructGig reconstructs a gig from persistence
func ReconstructGig(
	id valueobject.GigID,
	clientID valueobject.UserID,
	title string,
	description string,
	category string,
	skillID *valueobject.SkillID,
	budget Budget,
	currency valueobject.Currency,
	deliveryDays int,
	deadline *time.Time,
	isRemote bool,
	location string,
	status GigStatus,
	viewCount int,
	isFeatured bool,
	attachments []string,
	tags []string,
	proposals []*Proposal,
	acceptedProposalID *valueobject.ProposalID,
	createdAt time.Time,
	updatedAt time.Time,
	version int64,
) *Gig {
	return &Gig{
		id:                 id,
		clientID:           clientID,
		title:              title,
		description:        description,
		category:           category,
		skillID:            skillID,
		budget:             budget,
		currency:           currency,
		deliveryDays:       deliveryDays,
		deadline:           deadline,
		isRemote:           isRemote,
		location:           location,
		status:             status,
		viewCount:          viewCount,
		isFeatured:         isFeatured,
		attachments:        attachments,
		tags:               tags,
		proposals:          proposals,
		acceptedProposalID: acceptedProposalID,
		createdAt:          createdAt,
		updatedAt:          updatedAt,
		version:            version,
	}
}

// Getters
func (g *Gig) ID() valueobject.GigID           { return g.id }
func (g *Gig) ClientID() valueobject.UserID    { return g.clientID }
func (g *Gig) Title() string                   { return g.title }
func (g *Gig) Description() string             { return g.description }
func (g *Gig) Category() string                { return g.category }
func (g *Gig) SkillID() *valueobject.SkillID   { return g.skillID }
func (g *Gig) Budget() Budget                  { return g.budget }
func (g *Gig) Currency() valueobject.Currency  { return g.currency }
func (g *Gig) DeliveryDays() int               { return g.deliveryDays }
func (g *Gig) Deadline() *time.Time            { return g.deadline }
func (g *Gig) IsRemote() bool                  { return g.isRemote }
func (g *Gig) Location() string                { return g.location }
func (g *Gig) Status() GigStatus               { return g.status }
func (g *Gig) ViewCount() int                  { return g.viewCount }
func (g *Gig) IsFeatured() bool                { return g.isFeatured }
func (g *Gig) Attachments() []string           { return g.attachments }
func (g *Gig) Tags() []string                  { return g.tags }
func (g *Gig) Proposals() []*Proposal          { return g.proposals }
func (g *Gig) ProposalCount() int              { return len(g.proposals) }
func (g *Gig) AcceptedProposalID() *valueobject.ProposalID { return g.acceptedProposalID }
func (g *Gig) CreatedAt() time.Time            { return g.createdAt }
func (g *Gig) UpdatedAt() time.Time            { return g.updatedAt }
func (g *Gig) Version() int64                  { return g.version }
func (g *Gig) IsOpen() bool                    { return g.status.IsOpen() }

// Business Methods

// Update updates the gig details (only allowed when open)
func (g *Gig) Update(
	title string,
	description string,
	category string,
	skillID *valueobject.SkillID,
	budget Budget,
	deliveryDays int,
	deadline *time.Time,
	isRemote bool,
	location string,
	attachments []string,
	tags []string,
) error {
	if !g.status.IsOpen() {
		return ErrCannotUpdateGig
	}

	updatedFields := make(map[string]string)

	if title != g.title {
		g.title = title
		updatedFields["title"] = title
	}
	if description != g.description {
		g.description = description
		updatedFields["description"] = description
	}
	if category != g.category {
		g.category = category
		updatedFields["category"] = category
	}
	g.skillID = skillID
	g.budget = budget
	g.deliveryDays = deliveryDays
	g.deadline = deadline
	g.isRemote = isRemote
	g.location = location
	g.attachments = attachments
	g.tags = tags
	g.updatedAt = time.Now().UTC()

	if len(updatedFields) > 0 {
		g.RecordEvent(event.NewGigUpdated(g.id.String(), updatedFields))
	}

	return nil
}

// SetFeatured marks the gig as featured
func (g *Gig) SetFeatured(featured bool) {
	g.isFeatured = featured
	g.updatedAt = time.Now().UTC()
}

// IncrementViewCount increments the view counter
func (g *Gig) IncrementViewCount() {
	g.viewCount++
}

// Cancel cancels the gig
func (g *Gig) Cancel(reason string) error {
	if !g.status.IsOpen() {
		return ErrGigAlreadyInProgress
	}

	g.status = GigStatusCancelled
	g.updatedAt = time.Now().UTC()

	g.RecordEvent(event.NewGigCancelled(g.id.String(), g.clientID.String(), reason))

	return nil
}

// SubmitProposal adds a proposal to the gig
func (g *Gig) SubmitProposal(proposal *Proposal) error {
	if !g.status.IsOpen() {
		return ErrGigNotOpen
	}

	if proposal.HustlerID().Equals(g.clientID) {
		return ErrCannotProposeSelf
	}

	// Check for existing proposal from same hustler
	for _, p := range g.proposals {
		if p.HustlerID().Equals(proposal.HustlerID()) {
			return ErrAlreadyProposed
		}
	}

	// Validate proposed price is within budget
	if !g.budget.Contains(proposal.ProposedPrice()) {
		if proposal.ProposedPrice().LessThan(g.budget.Min()) {
			return ErrPriceBelowBudget
		}
		return ErrPriceAboveBudget
	}

	g.proposals = append(g.proposals, proposal)
	g.updatedAt = time.Now().UTC()

	g.RecordEvent(event.NewProposalSubmitted(
		proposal.ID().String(),
		g.id.String(),
		proposal.HustlerID().String(),
		proposal.ProposedPrice().Amount(),
		proposal.DeliveryDays(),
	))

	return nil
}

// WithdrawProposal withdraws a proposal
func (g *Gig) WithdrawProposal(proposalID valueobject.ProposalID, hustlerID valueobject.UserID) error {
	proposal := g.FindProposal(proposalID)
	if proposal == nil {
		return ErrProposalNotFound
	}

	if !proposal.HustlerID().Equals(hustlerID) {
		return ErrNotGigOwner
	}

	if !proposal.IsPending() {
		return ErrProposalNotPending
	}

	proposal.Withdraw()
	g.updatedAt = time.Now().UTC()

	g.RecordEvent(event.NewProposalWithdrawn(
		proposalID.String(),
		g.id.String(),
		hustlerID.String(),
	))

	return nil
}

// AcceptProposal accepts a proposal (returns contract creation data)
func (g *Gig) AcceptProposal(proposalID valueobject.ProposalID, contractID valueobject.ContractID) (*AcceptedProposalData, error) {
	if !g.status.IsOpen() {
		return nil, ErrGigNotOpen
	}

	proposal := g.FindProposal(proposalID)
	if proposal == nil {
		return nil, ErrProposalNotFound
	}

	if !proposal.IsPending() {
		return nil, ErrProposalNotPending
	}

	// Accept this proposal
	proposal.Accept()
	g.acceptedProposalID = &proposalID

	// Reject all other pending proposals
	for _, p := range g.proposals {
		if !p.ID().Equals(proposalID) && p.IsPending() {
			p.Reject()
		}
	}

	// Update gig status
	g.status = GigStatusInProgress
	g.updatedAt = time.Now().UTC()

	// Calculate platform fee (10%)
	platformFee := proposal.ProposedPrice().Amount() / 10
	deadlineAt := time.Now().UTC().AddDate(0, 0, proposal.DeliveryDays())

	g.RecordEvent(event.NewProposalAccepted(
		proposalID.String(),
		g.id.String(),
		contractID.String(),
		proposal.HustlerID().String(),
		g.clientID.String(),
		proposal.ProposedPrice().Amount(),
	))

	return &AcceptedProposalData{
		ContractID:    contractID,
		GigID:         g.id,
		ClientID:      g.clientID,
		HustlerID:     proposal.HustlerID(),
		AgreedPrice:   proposal.ProposedPrice(),
		PlatformFee:   platformFee,
		DeliveryDays:  proposal.DeliveryDays(),
		DeadlineAt:    deadlineAt,
	}, nil
}

// FindProposal finds a proposal by ID
func (g *Gig) FindProposal(proposalID valueobject.ProposalID) *Proposal {
	for _, p := range g.proposals {
		if p.ID().Equals(proposalID) {
			return p
		}
	}
	return nil
}

// FindProposalByHustler finds a proposal by hustler ID
func (g *Gig) FindProposalByHustler(hustlerID valueobject.UserID) *Proposal {
	for _, p := range g.proposals {
		if p.HustlerID().Equals(hustlerID) {
			return p
		}
	}
	return nil
}

// MarkCompleted marks the gig as completed
func (g *Gig) MarkCompleted() {
	g.status = GigStatusCompleted
	g.updatedAt = time.Now().UTC()
}

// MarkDisputed marks the gig as disputed
func (g *Gig) MarkDisputed() {
	g.status = GigStatusDisputed
	g.updatedAt = time.Now().UTC()
}

// AcceptedProposalData contains data needed to create a contract
type AcceptedProposalData struct {
	ContractID   valueobject.ContractID
	GigID        valueobject.GigID
	ClientID     valueobject.UserID
	HustlerID    valueobject.UserID
	AgreedPrice  valueobject.Money
	PlatformFee  int64
	DeliveryDays int
	DeadlineAt   time.Time
}
