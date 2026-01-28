package handler

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"hustlex/internal/application/savings/command"
	"hustlex/internal/domain/savings/aggregate"
	"hustlex/internal/domain/savings/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrCircleNotFound = errors.New("circle not found")
	ErrUnauthorized   = errors.New("unauthorized to perform this action")
)

// CircleHandler handles circle-related commands
type CircleHandler struct {
	circleRepo repository.CircleRepository
}

// NewCircleHandler creates a new circle handler
func NewCircleHandler(circleRepo repository.CircleRepository) *CircleHandler {
	return &CircleHandler{circleRepo: circleRepo}
}

// HandleCreateCircle creates a new savings circle
func (h *CircleHandler) HandleCreateCircle(ctx context.Context, cmd command.CreateCircle) (*command.CreateCircleResult, error) {
	creatorID, err := cmd.GetCreatorID()
	if err != nil {
		return nil, errors.New("invalid creator ID")
	}

	contributionAmt, err := cmd.GetContributionAmount()
	if err != nil {
		return nil, err
	}

	// Generate invite code
	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, errors.New("failed to generate invite code")
	}

	circleType := aggregate.CircleType(cmd.Type)
	if !circleType.IsValid() {
		return nil, errors.New("invalid circle type")
	}

	frequency := aggregate.ContributionFrequency(cmd.Frequency)

	circleID := valueobject.GenerateCircleID()
	circle, err := aggregate.NewCircle(
		circleID,
		creatorID,
		cmd.Name,
		cmd.Description,
		circleType,
		contributionAmt,
		frequency,
		cmd.MaxMembers,
		cmd.TotalRounds,
		cmd.IsPrivate,
		inviteCode,
	)
	if err != nil {
		return nil, err
	}

	if len(cmd.Rules) > 0 {
		circle.SetRules(cmd.Rules)
	}

	if err := h.circleRepo.SaveWithEvents(ctx, circle); err != nil {
		return nil, err
	}

	return &command.CreateCircleResult{
		CircleID:        circle.ID().String(),
		Name:            circle.Name(),
		Type:            circle.Type().String(),
		ContributionAmt: circle.ContributionAmount().Amount(),
		Currency:        string(circle.ContributionAmount().Currency()),
		Frequency:       string(circle.Frequency()),
		MaxMembers:      circle.MaxMembers(),
		InviteCode:      circle.InviteCode(),
		Status:          circle.Status().String(),
		CreatedAt:       circle.CreatedAt(),
	}, nil
}

// HandleJoinCircle adds a member to a circle
func (h *CircleHandler) HandleJoinCircle(ctx context.Context, cmd command.JoinCircle) (*command.JoinCircleResult, error) {
	circleID, err := cmd.GetCircleID()
	if err != nil {
		return nil, errors.New("invalid circle ID")
	}

	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	circle, err := h.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return nil, ErrCircleNotFound
	}

	member, err := circle.AddMember(userID)
	if err != nil {
		return nil, err
	}

	if err := h.circleRepo.SaveWithEvents(ctx, circle); err != nil {
		return nil, err
	}

	return &command.JoinCircleResult{
		MemberID: member.ID().String(),
		CircleID: circle.ID().String(),
		Position: member.Position(),
		Status:   string(member.Status()),
	}, nil
}

// HandleJoinCircleByCode joins a circle using invite code
func (h *CircleHandler) HandleJoinCircleByCode(ctx context.Context, cmd command.JoinCircleByCode) (*command.JoinCircleResult, error) {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	circle, err := h.circleRepo.FindByInviteCode(ctx, cmd.InviteCode)
	if err != nil {
		return nil, errors.New("invalid invite code")
	}

	member, err := circle.AddMember(userID)
	if err != nil {
		return nil, err
	}

	if err := h.circleRepo.SaveWithEvents(ctx, circle); err != nil {
		return nil, err
	}

	return &command.JoinCircleResult{
		MemberID: member.ID().String(),
		CircleID: circle.ID().String(),
		Position: member.Position(),
		Status:   string(member.Status()),
	}, nil
}

// HandleLeaveCircle removes a member from a circle
func (h *CircleHandler) HandleLeaveCircle(ctx context.Context, cmd command.LeaveCircle) error {
	circleID, err := cmd.GetCircleID()
	if err != nil {
		return errors.New("invalid circle ID")
	}

	userID, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	circle, err := h.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return ErrCircleNotFound
	}

	if err := circle.RemoveMember(userID); err != nil {
		return err
	}

	return h.circleRepo.SaveWithEvents(ctx, circle)
}

// HandleStartCircle manually starts a circle
func (h *CircleHandler) HandleStartCircle(ctx context.Context, cmd command.StartCircle) error {
	circleID, err := valueobject.NewCircleID(cmd.CircleID)
	if err != nil {
		return errors.New("invalid circle ID")
	}

	adminID, err := valueobject.NewUserID(cmd.AdminID)
	if err != nil {
		return errors.New("invalid admin ID")
	}

	circle, err := h.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return ErrCircleNotFound
	}

	// Verify admin
	if !circle.IsAdmin(adminID) {
		return ErrUnauthorized
	}

	if err := circle.Start(); err != nil {
		return err
	}

	return h.circleRepo.SaveWithEvents(ctx, circle)
}

// HandleMakeContribution records a contribution
func (h *CircleHandler) HandleMakeContribution(ctx context.Context, cmd command.MakeContribution) (*command.MakeContributionResult, error) {
	circleID, err := cmd.GetCircleID()
	if err != nil {
		return nil, errors.New("invalid circle ID")
	}

	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	transactionID, err := cmd.GetTransactionID()
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	circle, err := h.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return nil, ErrCircleNotFound
	}

	member := circle.FindMemberByUserID(userID)
	if member == nil {
		return nil, aggregate.ErrNotMember
	}

	contribution, err := circle.RecordContribution(member.ID(), transactionID)
	if err != nil {
		return nil, err
	}

	if err := h.circleRepo.SaveWithEvents(ctx, circle); err != nil {
		return nil, err
	}

	return &command.MakeContributionResult{
		ContributionID: contribution.ID().String(),
		CircleID:       circle.ID().String(),
		Round:          contribution.Round(),
		Amount:         contribution.Amount().Amount(),
		LateFee:        contribution.LateFee(),
		TotalPaid:      contribution.Amount().Amount() + contribution.LateFee(),
		PaidAt:         *contribution.PaidAt(),
	}, nil
}

func generateInviteCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, 8)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
