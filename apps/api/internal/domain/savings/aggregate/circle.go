package aggregate

import (
	"errors"
	"time"

	"hustlex/internal/domain/savings/event"
	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrCircleNotRecruiting    = errors.New("circle is not recruiting members")
	ErrCircleFull             = errors.New("circle has reached maximum members")
	ErrAlreadyMember          = errors.New("user is already a member of this circle")
	ErrNotMember              = errors.New("user is not a member of this circle")
	ErrCircleNotActive        = errors.New("circle is not active")
	ErrCannotLeaveActiveCircle = errors.New("cannot leave circle after contributions started")
	ErrAdminCannotLeave       = errors.New("admin cannot leave the circle")
	ErrMinimumMembers         = errors.New("need at least 2 members to start")
	ErrAlreadyStarted         = errors.New("circle has already started")
	ErrNoPendingContribution  = errors.New("no pending contribution for this round")
	ErrNotAdmin               = errors.New("only admin can perform this action")
	ErrInvalidFrequency       = errors.New("invalid contribution frequency")
)

// CircleType represents the type of savings circle
type CircleType string

const (
	CircleTypeRotational  CircleType = "rotational"
	CircleTypeFixedTarget CircleType = "fixed_target"
	CircleTypeEmergency   CircleType = "emergency"
)

func (t CircleType) String() string {
	return string(t)
}

func (t CircleType) IsValid() bool {
	switch t {
	case CircleTypeRotational, CircleTypeFixedTarget, CircleTypeEmergency:
		return true
	}
	return false
}

// CircleStatus represents the status of a circle
type CircleStatus string

const (
	CircleStatusRecruiting CircleStatus = "recruiting"
	CircleStatusActive     CircleStatus = "active"
	CircleStatusCompleted  CircleStatus = "completed"
	CircleStatusCancelled  CircleStatus = "cancelled"
)

func (s CircleStatus) String() string {
	return string(s)
}

func (s CircleStatus) IsRecruiting() bool { return s == CircleStatusRecruiting }
func (s CircleStatus) IsActive() bool     { return s == CircleStatusActive }

// ContributionFrequency represents how often contributions are made
type ContributionFrequency string

const (
	FrequencyDaily    ContributionFrequency = "daily"
	FrequencyWeekly   ContributionFrequency = "weekly"
	FrequencyBiweekly ContributionFrequency = "biweekly"
	FrequencyMonthly  ContributionFrequency = "monthly"
)

func (f ContributionFrequency) NextDueDate(from time.Time) time.Time {
	switch f {
	case FrequencyDaily:
		return from.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return from.AddDate(0, 0, 7)
	case FrequencyBiweekly:
		return from.AddDate(0, 0, 14)
	case FrequencyMonthly:
		return from.AddDate(0, 1, 0)
	default:
		return from.AddDate(0, 0, 7)
	}
}

// MemberRole represents a member's role in the circle
type MemberRole string

const (
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

// MemberStatus represents a member's status
type MemberStatus string

const (
	MemberStatusPending MemberStatus = "pending"
	MemberStatusActive  MemberStatus = "active"
	MemberStatusLeft    MemberStatus = "left"
	MemberStatusRemoved MemberStatus = "removed"
)

// ContributionStatus represents contribution status
type ContributionStatus string

const (
	ContributionPending ContributionStatus = "pending"
	ContributionPaid    ContributionStatus = "paid"
	ContributionOverdue ContributionStatus = "overdue"
	ContributionWaived  ContributionStatus = "waived"
)

// Member represents a circle member entity
type Member struct {
	id             valueobject.MemberID
	userID         valueobject.UserID
	position       int
	role           MemberRole
	status         MemberStatus
	totalContrib   int64
	missedPayments int
	hasReceived    bool
	joinedAt       time.Time
}

func NewMember(id valueobject.MemberID, userID valueobject.UserID, position int, role MemberRole) *Member {
	return &Member{
		id:             id,
		userID:         userID,
		position:       position,
		role:           role,
		status:         MemberStatusActive,
		totalContrib:   0,
		missedPayments: 0,
		hasReceived:    false,
		joinedAt:       time.Now().UTC(),
	}
}

func (m *Member) ID() valueobject.MemberID { return m.id }
func (m *Member) UserID() valueobject.UserID { return m.userID }
func (m *Member) Position() int { return m.position }
func (m *Member) Role() MemberRole { return m.role }
func (m *Member) Status() MemberStatus { return m.status }
func (m *Member) TotalContrib() int64 { return m.totalContrib }
func (m *Member) MissedPayments() int { return m.missedPayments }
func (m *Member) HasReceived() bool { return m.hasReceived }
func (m *Member) JoinedAt() time.Time { return m.joinedAt }
func (m *Member) IsAdmin() bool { return m.role == RoleAdmin }
func (m *Member) IsActive() bool { return m.status == MemberStatusActive }

func (m *Member) RecordContribution(amount int64) {
	m.totalContrib += amount
}

func (m *Member) RecordMissedPayment() {
	m.missedPayments++
}

func (m *Member) MarkReceived() {
	m.hasReceived = true
}

func (m *Member) ResetForNewCycle() {
	m.hasReceived = false
}

func (m *Member) Leave() {
	m.status = MemberStatusLeft
}

func (m *Member) UpdatePosition(newPosition int) {
	m.position = newPosition
}

// Contribution represents a scheduled contribution
type Contribution struct {
	id            valueobject.ContributionID
	memberID      valueobject.MemberID
	round         int
	amount        valueobject.Money
	dueDate       time.Time
	paidAt        *time.Time
	status        ContributionStatus
	transactionID *valueobject.TransactionID
	lateFee       int64
}

func NewContribution(id valueobject.ContributionID, memberID valueobject.MemberID, round int, amount valueobject.Money, dueDate time.Time) *Contribution {
	return &Contribution{
		id:       id,
		memberID: memberID,
		round:    round,
		amount:   amount,
		dueDate:  dueDate,
		status:   ContributionPending,
		lateFee:  0,
	}
}

func (c *Contribution) ID() valueobject.ContributionID { return c.id }
func (c *Contribution) MemberID() valueobject.MemberID { return c.memberID }
func (c *Contribution) Round() int { return c.round }
func (c *Contribution) Amount() valueobject.Money { return c.amount }
func (c *Contribution) DueDate() time.Time { return c.dueDate }
func (c *Contribution) PaidAt() *time.Time { return c.paidAt }
func (c *Contribution) Status() ContributionStatus { return c.status }
func (c *Contribution) TransactionID() *valueobject.TransactionID { return c.transactionID }
func (c *Contribution) LateFee() int64 { return c.lateFee }
func (c *Contribution) IsPending() bool { return c.status == ContributionPending }

func (c *Contribution) MarkPaid(transactionID valueobject.TransactionID, lateFee int64) {
	now := time.Now().UTC()
	c.status = ContributionPaid
	c.paidAt = &now
	c.transactionID = &transactionID
	c.lateFee = lateFee
}

func (c *Contribution) IsOverdue() bool {
	return c.status == ContributionPending && time.Now().After(c.dueDate)
}

// Circle is the aggregate root for savings circles
type Circle struct {
	sharedevent.AggregateRoot

	id              valueobject.CircleID
	name            string
	description     string
	circleType      CircleType
	contributionAmt valueobject.Money
	frequency       ContributionFrequency
	maxMembers      int
	totalRounds     int
	currentRound    int
	poolBalance     int64
	totalSaved      int64
	status          CircleStatus
	isPrivate       bool
	inviteCode      string
	rules           []string
	startDate       *time.Time
	nextPayoutDate  *time.Time
	members         []*Member
	contributions   []*Contribution
	createdBy       valueobject.UserID
	createdAt       time.Time
	updatedAt       time.Time
	version         int64
}

// NewCircle creates a new savings circle
func NewCircle(
	id valueobject.CircleID,
	creatorID valueobject.UserID,
	name string,
	description string,
	circleType CircleType,
	contributionAmt valueobject.Money,
	frequency ContributionFrequency,
	maxMembers int,
	totalRounds int,
	isPrivate bool,
	inviteCode string,
) (*Circle, error) {
	circle := &Circle{
		id:              id,
		name:            name,
		description:     description,
		circleType:      circleType,
		contributionAmt: contributionAmt,
		frequency:       frequency,
		maxMembers:      maxMembers,
		totalRounds:     totalRounds,
		currentRound:    0,
		poolBalance:     0,
		totalSaved:      0,
		status:          CircleStatusRecruiting,
		isPrivate:       isPrivate,
		inviteCode:      inviteCode,
		rules:           make([]string, 0),
		members:         make([]*Member, 0),
		contributions:   make([]*Contribution, 0),
		createdBy:       creatorID,
		createdAt:       time.Now().UTC(),
		updatedAt:       time.Now().UTC(),
		version:         1,
	}

	// Add creator as admin
	memberID := valueobject.GenerateMemberID()
	adminMember := NewMember(memberID, creatorID, 1, RoleAdmin)
	circle.members = append(circle.members, adminMember)

	circle.RecordEvent(event.NewCircleCreated(
		id.String(),
		creatorID.String(),
		name,
		circleType.String(),
		contributionAmt.Amount(),
		string(frequency),
		maxMembers,
	))

	return circle, nil
}

// Getters
func (c *Circle) ID() valueobject.CircleID { return c.id }
func (c *Circle) Name() string { return c.name }
func (c *Circle) Description() string { return c.description }
func (c *Circle) Type() CircleType { return c.circleType }
func (c *Circle) ContributionAmount() valueobject.Money { return c.contributionAmt }
func (c *Circle) Frequency() ContributionFrequency { return c.frequency }
func (c *Circle) MaxMembers() int { return c.maxMembers }
func (c *Circle) TotalRounds() int { return c.totalRounds }
func (c *Circle) CurrentRound() int { return c.currentRound }
func (c *Circle) PoolBalance() int64 { return c.poolBalance }
func (c *Circle) TotalSaved() int64 { return c.totalSaved }
func (c *Circle) Status() CircleStatus { return c.status }
func (c *Circle) IsPrivate() bool { return c.isPrivate }
func (c *Circle) InviteCode() string { return c.inviteCode }
func (c *Circle) Rules() []string { return c.rules }
func (c *Circle) StartDate() *time.Time { return c.startDate }
func (c *Circle) NextPayoutDate() *time.Time { return c.nextPayoutDate }
func (c *Circle) Members() []*Member { return c.members }
func (c *Circle) Contributions() []*Contribution { return c.contributions }
func (c *Circle) CreatedBy() valueobject.UserID { return c.createdBy }
func (c *Circle) CreatedAt() time.Time { return c.createdAt }
func (c *Circle) UpdatedAt() time.Time { return c.updatedAt }
func (c *Circle) Version() int64 { return c.version }
func (c *Circle) CurrentMembers() int { return len(c.activeMembers()) }
func (c *Circle) IsFull() bool { return c.CurrentMembers() >= c.maxMembers }

func (c *Circle) activeMembers() []*Member {
	active := make([]*Member, 0)
	for _, m := range c.members {
		if m.IsActive() {
			active = append(active, m)
		}
	}
	return active
}

// Business Methods

// SetRules sets the circle rules
func (c *Circle) SetRules(rules []string) {
	c.rules = rules
	c.updatedAt = time.Now().UTC()
}

// AddMember adds a new member to the circle
func (c *Circle) AddMember(userID valueobject.UserID) (*Member, error) {
	if !c.status.IsRecruiting() {
		return nil, ErrCircleNotRecruiting
	}

	if c.IsFull() {
		return nil, ErrCircleFull
	}

	// Check if already a member
	if c.FindMemberByUserID(userID) != nil {
		return nil, ErrAlreadyMember
	}

	position := c.CurrentMembers() + 1
	memberID := valueobject.GenerateMemberID()
	member := NewMember(memberID, userID, position, RoleMember)
	c.members = append(c.members, member)
	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewMemberJoined(
		c.id.String(),
		memberID.String(),
		userID.String(),
		position,
	))

	// Auto-start if full
	if c.IsFull() {
		_ = c.Start()
	}

	return member, nil
}

// RemoveMember allows a member to leave the circle
func (c *Circle) RemoveMember(userID valueobject.UserID) error {
	member := c.FindMemberByUserID(userID)
	if member == nil {
		return ErrNotMember
	}

	if c.status.IsActive() && c.currentRound > 0 {
		return ErrCannotLeaveActiveCircle
	}

	if member.IsAdmin() {
		return ErrAdminCannotLeave
	}

	member.Leave()
	c.reorderPositions()
	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewMemberLeft(
		c.id.String(),
		member.ID().String(),
		userID.String(),
	))

	return nil
}

func (c *Circle) reorderPositions() {
	position := 1
	for _, m := range c.members {
		if m.IsActive() {
			m.UpdatePosition(position)
			position++
		}
	}
}

// Start starts the circle
func (c *Circle) Start() error {
	if !c.status.IsRecruiting() {
		return ErrAlreadyStarted
	}

	if c.CurrentMembers() < 2 {
		return ErrMinimumMembers
	}

	now := time.Now().UTC()
	c.status = CircleStatusActive
	c.startDate = &now
	c.currentRound = 1
	c.updatedAt = now

	// Schedule first round contributions
	c.scheduleContributions()

	c.RecordEvent(event.NewCircleStarted(c.id.String()))

	return nil
}

func (c *Circle) scheduleContributions() {
	dueDate := c.frequency.NextDueDate(time.Now().UTC())

	for _, member := range c.activeMembers() {
		contribID := valueobject.GenerateContributionID()
		contribution := NewContribution(
			contribID,
			member.ID(),
			c.currentRound,
			c.contributionAmt,
			dueDate,
		)
		c.contributions = append(c.contributions, contribution)
	}

	c.nextPayoutDate = &dueDate
}

// RecordContribution records a member's contribution
func (c *Circle) RecordContribution(memberID valueobject.MemberID, transactionID valueobject.TransactionID) (*Contribution, error) {
	if !c.status.IsActive() {
		return nil, ErrCircleNotActive
	}

	// Find pending contribution
	var contribution *Contribution
	for _, cont := range c.contributions {
		if cont.MemberID().Equals(memberID) && cont.Round() == c.currentRound && cont.IsPending() {
			contribution = cont
			break
		}
	}

	if contribution == nil {
		return nil, ErrNoPendingContribution
	}

	// Calculate late fee (5% if overdue)
	var lateFee int64 = 0
	if contribution.IsOverdue() {
		lateFee = c.contributionAmt.Amount() / 20
	}

	contribution.MarkPaid(transactionID, lateFee)

	// Update pool balance
	c.poolBalance += c.contributionAmt.Amount() + lateFee
	c.totalSaved += c.contributionAmt.Amount() + lateFee

	// Update member stats
	member := c.FindMemberByID(memberID)
	if member != nil {
		member.RecordContribution(c.contributionAmt.Amount() + lateFee)
	}

	c.updatedAt = time.Now().UTC()

	c.RecordEvent(event.NewContributionMade(
		c.id.String(),
		contribution.ID().String(),
		memberID.String(),
		c.currentRound,
		c.contributionAmt.Amount(),
		lateFee,
	))

	// Check if round is complete
	if c.isRoundComplete() {
		c.completeRound()
	}

	return contribution, nil
}

func (c *Circle) isRoundComplete() bool {
	for _, cont := range c.contributions {
		if cont.Round() == c.currentRound && cont.IsPending() {
			return false
		}
	}
	return true
}

func (c *Circle) completeRound() {
	// Find recipient for this round
	recipient := c.FindMemberByPosition(c.currentRound)
	if recipient != nil {
		recipient.MarkReceived()

		c.RecordEvent(event.NewPayoutTriggered(
			c.id.String(),
			recipient.ID().String(),
			recipient.UserID().String(),
			c.currentRound,
			c.poolBalance,
		))
	}

	c.RecordEvent(event.NewRoundCompleted(
		c.id.String(),
		c.currentRound,
		c.poolBalance,
	))

	// Reset pool and advance round
	c.poolBalance = 0
	c.currentRound++

	if c.currentRound > c.totalRounds {
		c.status = CircleStatusCompleted
		c.RecordEvent(event.NewCircleCompleted(c.id.String(), c.totalRounds, c.totalSaved))
	} else {
		c.scheduleContributions()
	}

	c.updatedAt = time.Now().UTC()
}

// Helper methods

func (c *Circle) FindMemberByUserID(userID valueobject.UserID) *Member {
	for _, m := range c.members {
		if m.UserID().Equals(userID) && m.IsActive() {
			return m
		}
	}
	return nil
}

func (c *Circle) FindMemberByID(memberID valueobject.MemberID) *Member {
	for _, m := range c.members {
		if m.ID().Equals(memberID) {
			return m
		}
	}
	return nil
}

func (c *Circle) FindMemberByPosition(position int) *Member {
	for _, m := range c.members {
		if m.Position() == position && m.IsActive() {
			return m
		}
	}
	return nil
}

func (c *Circle) IsMember(userID valueobject.UserID) bool {
	return c.FindMemberByUserID(userID) != nil
}

func (c *Circle) IsAdmin(userID valueobject.UserID) bool {
	member := c.FindMemberByUserID(userID)
	return member != nil && member.IsAdmin()
}

func (c *Circle) GetPendingContributions(memberID valueobject.MemberID) []*Contribution {
	pending := make([]*Contribution, 0)
	for _, cont := range c.contributions {
		if cont.MemberID().Equals(memberID) && cont.IsPending() {
			pending = append(pending, cont)
		}
	}
	return pending
}
