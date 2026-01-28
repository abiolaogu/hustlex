package query

import (
	"context"
	"time"

	"hustlex/internal/domain/savings/aggregate"
	"hustlex/internal/domain/savings/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// GetCircle retrieves a single circle
type GetCircle struct {
	CircleID string
	UserID   string // for access check
}

// GetCircleByCode retrieves a circle by invite code
type GetCircleByCode struct {
	InviteCode string
}

// GetCircles retrieves circles with filters
type GetCircles struct {
	Type      string
	Status    string
	MinAmount int64
	MaxAmount int64
	Frequency string
	Search    string
	Page      int
	Limit     int
}

// GetMyCircles retrieves circles user is a member of
type GetMyCircles struct {
	UserID string
	Status string
}

// CircleDTO represents a circle for API responses
type CircleDTO struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	Type            string       `json:"type"`
	ContributionAmt int64        `json:"contribution_amount"`
	Currency        string       `json:"currency"`
	Frequency       string       `json:"frequency"`
	MaxMembers      int          `json:"max_members"`
	CurrentMembers  int          `json:"current_members"`
	TotalRounds     int          `json:"total_rounds"`
	CurrentRound    int          `json:"current_round"`
	PoolBalance     int64        `json:"pool_balance"`
	TotalSaved      int64        `json:"total_saved"`
	Status          string       `json:"status"`
	IsPrivate       bool         `json:"is_private"`
	InviteCode      string       `json:"invite_code,omitempty"`
	Rules           []string     `json:"rules,omitempty"`
	StartDate       *time.Time   `json:"start_date,omitempty"`
	NextPayoutDate  *time.Time   `json:"next_payout_date,omitempty"`
	CreatorID       string       `json:"creator_id"`
	CreatorName     string       `json:"creator_name,omitempty"`
	Members         []MemberDTO  `json:"members,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
}

// MemberDTO represents a circle member
type MemberDTO struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	UserName       string    `json:"user_name,omitempty"`
	UserImage      string    `json:"user_image,omitempty"`
	Position       int       `json:"position"`
	Role           string    `json:"role"`
	Status         string    `json:"status"`
	TotalContrib   int64     `json:"total_contributed"`
	MissedPayments int       `json:"missed_payments"`
	HasReceived    bool      `json:"has_received"`
	JoinedAt       time.Time `json:"joined_at"`
}

// CircleListResult represents paginated circle results
type CircleListResult struct {
	Circles    []CircleDTO `json:"circles"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// GetCircleMembers retrieves members of a circle
type GetCircleMembers struct {
	CircleID string
}

// GetCircleContributions retrieves contributions for a circle
type GetCircleContributions struct {
	CircleID string
	Round    *int
}

// ContributionDTO represents a contribution
type ContributionDTO struct {
	ID         string     `json:"id"`
	MemberID   string     `json:"member_id"`
	UserID     string     `json:"user_id"`
	UserName   string     `json:"user_name,omitempty"`
	Round      int        `json:"round"`
	Amount     int64      `json:"amount"`
	Currency   string     `json:"currency"`
	DueDate    time.Time  `json:"due_date"`
	PaidAt     *time.Time `json:"paid_at,omitempty"`
	Status     string     `json:"status"`
	LateFee    int64      `json:"late_fee"`
	IsOverdue  bool       `json:"is_overdue"`
}

// GetMyContributions retrieves user's contributions
type GetMyContributions struct {
	UserID   string
	CircleID string
}

// GetUserSavingsStats retrieves user's savings statistics
type GetUserSavingsStats struct {
	UserID string
}

// UserSavingsStatsDTO represents savings statistics
type UserSavingsStatsDTO struct {
	TotalCircles     int     `json:"total_circles"`
	ActiveCircles    int     `json:"active_circles"`
	CompletedCircles int     `json:"completed_circles"`
	TotalContributed int64   `json:"total_contributed"`
	TotalReceived    int64   `json:"total_received"`
	MissedPayments   int     `json:"missed_payments"`
	OnTimeRate       float64 `json:"on_time_rate"`
}

// CircleQueryHandler handles circle queries
type CircleQueryHandler struct {
	circleRepo       repository.CircleRepository
	memberRepo       repository.MemberRepository
	contributionRepo repository.ContributionRepository
	statsRepo        repository.SavingsStatisticsRepository
}

// NewCircleQueryHandler creates a new query handler
func NewCircleQueryHandler(
	circleRepo repository.CircleRepository,
	memberRepo repository.MemberRepository,
	contributionRepo repository.ContributionRepository,
	statsRepo repository.SavingsStatisticsRepository,
) *CircleQueryHandler {
	return &CircleQueryHandler{
		circleRepo:       circleRepo,
		memberRepo:       memberRepo,
		contributionRepo: contributionRepo,
		statsRepo:        statsRepo,
	}
}

// HandleGetCircle retrieves a single circle
func (h *CircleQueryHandler) HandleGetCircle(ctx context.Context, q GetCircle) (*CircleDTO, error) {
	circleID, err := valueobject.NewCircleID(q.CircleID)
	if err != nil {
		return nil, err
	}

	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	circle, err := h.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return nil, err
	}

	// Check access for private circles
	if circle.IsPrivate() && !circle.IsMember(userID) {
		return nil, aggregate.ErrNotMember
	}

	return circleToDTO(circle, circle.IsMember(userID)), nil
}

// HandleGetCircles retrieves circles with filters
func (h *CircleQueryHandler) HandleGetCircles(ctx context.Context, q GetCircles) (*CircleListResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	filter := repository.CircleFilter{
		MinAmount: q.MinAmount,
		MaxAmount: q.MaxAmount,
		Search:    q.Search,
		IsPublic:  true,
		Offset:    (q.Page - 1) * q.Limit,
		Limit:     q.Limit,
	}

	if q.Type != "" {
		t := aggregate.CircleType(q.Type)
		filter.Type = &t
	}

	if q.Status != "" {
		s := aggregate.CircleStatus(q.Status)
		filter.Status = &s
	} else {
		s := aggregate.CircleStatusRecruiting
		filter.Status = &s
	}

	if q.Frequency != "" {
		f := aggregate.ContributionFrequency(q.Frequency)
		filter.Frequency = &f
	}

	circles, total, err := h.circleRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]CircleDTO, len(circles))
	for i, c := range circles {
		dtos[i] = CircleDTO{
			ID:              c.ID,
			Name:            c.Name,
			Description:     c.Description,
			Type:            c.Type,
			ContributionAmt: c.ContributionAmt,
			Currency:        c.Currency,
			Frequency:       c.Frequency,
			MaxMembers:      c.MaxMembers,
			CurrentMembers:  c.CurrentMembers,
			TotalRounds:     c.TotalRounds,
			CurrentRound:    c.CurrentRound,
			Status:          c.Status,
			IsPrivate:       c.IsPrivate,
			CreatorID:       c.CreatorID,
			CreatorName:     c.CreatorName,
			StartDate:       c.StartDate,
			NextPayoutDate:  c.NextPayoutDate,
			CreatedAt:       c.CreatedAt,
		}
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &CircleListResult{
		Circles:    dtos,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// HandleGetMyCircles retrieves circles user is a member of
func (h *CircleQueryHandler) HandleGetMyCircles(ctx context.Context, q GetMyCircles) ([]CircleDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	var status *aggregate.CircleStatus
	if q.Status != "" {
		s := aggregate.CircleStatus(q.Status)
		status = &s
	}

	circles, err := h.circleRepo.FindByUserID(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	dtos := make([]CircleDTO, len(circles))
	for i, c := range circles {
		dtos[i] = *circleToDTO(c, true)
	}

	return dtos, nil
}

// HandleGetCircleContributions retrieves contributions for a circle
func (h *CircleQueryHandler) HandleGetCircleContributions(ctx context.Context, q GetCircleContributions) ([]ContributionDTO, error) {
	circleID, err := valueobject.NewCircleID(q.CircleID)
	if err != nil {
		return nil, err
	}

	var contributions []*repository.ContributionDTO
	if q.Round != nil {
		contributions, err = h.contributionRepo.FindByCircleAndRound(ctx, circleID, *q.Round)
	} else {
		contributions, err = h.contributionRepo.FindPending(ctx, circleID)
	}

	if err != nil {
		return nil, err
	}

	dtos := make([]ContributionDTO, len(contributions))
	for i, c := range contributions {
		dtos[i] = ContributionDTO{
			ID:        c.ID,
			MemberID:  c.MemberID,
			UserID:    c.UserID,
			UserName:  c.UserName,
			Round:     c.Round,
			Amount:    c.Amount,
			Currency:  c.Currency,
			DueDate:   c.DueDate,
			PaidAt:    c.PaidAt,
			Status:    c.Status,
			LateFee:   c.LateFee,
			IsOverdue: c.Status == "pending" && time.Now().After(c.DueDate),
		}
	}

	return dtos, nil
}

// HandleGetUserSavingsStats retrieves user's savings statistics
func (h *CircleQueryHandler) HandleGetUserSavingsStats(ctx context.Context, q GetUserSavingsStats) (*UserSavingsStatsDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	stats, err := h.statsRepo.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserSavingsStatsDTO{
		TotalCircles:     stats.TotalCircles,
		ActiveCircles:    stats.ActiveCircles,
		CompletedCircles: stats.CompletedCircles,
		TotalContributed: stats.TotalContributed,
		TotalReceived:    stats.TotalReceived,
		MissedPayments:   stats.MissedPayments,
		OnTimeRate:       stats.OnTimeRate,
	}, nil
}

func circleToDTO(circle *aggregate.Circle, includeSensitive bool) *CircleDTO {
	dto := &CircleDTO{
		ID:              circle.ID().String(),
		Name:            circle.Name(),
		Description:     circle.Description(),
		Type:            circle.Type().String(),
		ContributionAmt: circle.ContributionAmount().Amount(),
		Currency:        string(circle.ContributionAmount().Currency()),
		Frequency:       string(circle.Frequency()),
		MaxMembers:      circle.MaxMembers(),
		CurrentMembers:  circle.CurrentMembers(),
		TotalRounds:     circle.TotalRounds(),
		CurrentRound:    circle.CurrentRound(),
		PoolBalance:     circle.PoolBalance(),
		TotalSaved:      circle.TotalSaved(),
		Status:          circle.Status().String(),
		IsPrivate:       circle.IsPrivate(),
		Rules:           circle.Rules(),
		StartDate:       circle.StartDate(),
		NextPayoutDate:  circle.NextPayoutDate(),
		CreatorID:       circle.CreatedBy().String(),
		CreatedAt:       circle.CreatedAt(),
	}

	if includeSensitive {
		dto.InviteCode = circle.InviteCode()

		members := make([]MemberDTO, 0)
		for _, m := range circle.Members() {
			if m.IsActive() {
				members = append(members, MemberDTO{
					ID:             m.ID().String(),
					UserID:         m.UserID().String(),
					Position:       m.Position(),
					Role:           string(m.Role()),
					Status:         string(m.Status()),
					TotalContrib:   m.TotalContrib(),
					MissedPayments: m.MissedPayments(),
					HasReceived:    m.HasReceived(),
					JoinedAt:       m.JoinedAt(),
				})
			}
		}
		dto.Members = members
	}

	return dto
}
