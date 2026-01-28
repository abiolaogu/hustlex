package handler

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/application/gig/command"
	"hustlex/internal/domain/gig/aggregate"
	"hustlex/internal/domain/gig/repository"
	"hustlex/internal/domain/gig/service"
	"hustlex/internal/domain/shared/valueobject"
)

// GigHandler handles gig-related commands
type GigHandler struct {
	gigRepo      repository.GigRepository
	proposalRepo repository.ProposalRepository
}

// NewGigHandler creates a new gig handler
func NewGigHandler(
	gigRepo repository.GigRepository,
	proposalRepo repository.ProposalRepository,
) *GigHandler {
	return &GigHandler{
		gigRepo:      gigRepo,
		proposalRepo: proposalRepo,
	}
}

// HandleCreateGig creates a new gig
func (h *GigHandler) HandleCreateGig(ctx context.Context, cmd command.CreateGig) (*command.CreateGigResult, error) {
	clientID, err := cmd.GetClientID()
	if err != nil {
		return nil, errors.New("invalid client ID")
	}

	budgetMin, budgetMax, err := cmd.GetBudget()
	if err != nil {
		return nil, err
	}

	budget, err := aggregate.NewBudget(budgetMin, budgetMax)
	if err != nil {
		return nil, err
	}

	skillID, err := cmd.GetSkillID()
	if err != nil {
		return nil, errors.New("invalid skill ID")
	}

	// Create gig aggregate
	gigID := valueobject.GenerateGigID()
	gig, err := aggregate.NewGig(
		gigID,
		clientID,
		cmd.Title,
		cmd.Description,
		cmd.Category,
		budget,
		cmd.DeliveryDays,
		cmd.IsRemote,
	)
	if err != nil {
		return nil, err
	}

	// Update optional fields
	if err := gig.Update(
		cmd.Title,
		cmd.Description,
		cmd.Category,
		skillID,
		budget,
		cmd.DeliveryDays,
		cmd.Deadline,
		cmd.IsRemote,
		cmd.Location,
		cmd.Attachments,
		cmd.Tags,
	); err != nil {
		return nil, err
	}

	// Save gig with events
	if err := h.gigRepo.SaveWithEvents(ctx, gig); err != nil {
		return nil, err
	}

	return &command.CreateGigResult{
		GigID:        gig.ID().String(),
		Title:        gig.Title(),
		Category:     gig.Category(),
		BudgetMin:    gig.Budget().Min().Amount(),
		BudgetMax:    gig.Budget().Max().Amount(),
		Currency:     string(gig.Currency()),
		DeliveryDays: gig.DeliveryDays(),
		Status:       gig.Status().String(),
		CreatedAt:    gig.CreatedAt(),
	}, nil
}

// HandleUpdateGig updates a gig
func (h *GigHandler) HandleUpdateGig(ctx context.Context, cmd command.UpdateGig) (*command.CreateGigResult, error) {
	gigID, err := valueobject.NewGigID(cmd.GigID)
	if err != nil {
		return nil, errors.New("invalid gig ID")
	}

	clientID, err := valueobject.NewUserID(cmd.ClientID)
	if err != nil {
		return nil, errors.New("invalid client ID")
	}

	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return nil, service.ErrGigNotFound
	}

	// Verify ownership
	if !gig.ClientID().Equals(clientID) {
		return nil, service.ErrUnauthorized
	}

	// Prepare budget
	currency := gig.Currency()
	budgetMin, err := valueobject.NewMoney(cmd.BudgetMin, currency)
	if err != nil {
		return nil, err
	}
	budgetMax, err := valueobject.NewMoney(cmd.BudgetMax, currency)
	if err != nil {
		return nil, err
	}
	budget, err := aggregate.NewBudget(budgetMin, budgetMax)
	if err != nil {
		return nil, err
	}

	// Prepare skill ID
	var skillID *valueobject.SkillID
	if cmd.SkillID != "" {
		sid, err := valueobject.NewSkillID(cmd.SkillID)
		if err == nil {
			skillID = &sid
		}
	}

	// Update gig
	if err := gig.Update(
		cmd.Title,
		cmd.Description,
		cmd.Category,
		skillID,
		budget,
		cmd.DeliveryDays,
		cmd.Deadline,
		cmd.IsRemote,
		cmd.Location,
		cmd.Attachments,
		cmd.Tags,
	); err != nil {
		return nil, err
	}

	if err := h.gigRepo.SaveWithEvents(ctx, gig); err != nil {
		return nil, err
	}

	return &command.CreateGigResult{
		GigID:        gig.ID().String(),
		Title:        gig.Title(),
		Category:     gig.Category(),
		BudgetMin:    gig.Budget().Min().Amount(),
		BudgetMax:    gig.Budget().Max().Amount(),
		Currency:     string(gig.Currency()),
		DeliveryDays: gig.DeliveryDays(),
		Status:       gig.Status().String(),
		CreatedAt:    gig.CreatedAt(),
	}, nil
}

// HandleCancelGig cancels a gig
func (h *GigHandler) HandleCancelGig(ctx context.Context, cmd command.CancelGig) error {
	gigID, err := valueobject.NewGigID(cmd.GigID)
	if err != nil {
		return errors.New("invalid gig ID")
	}

	clientID, err := valueobject.NewUserID(cmd.ClientID)
	if err != nil {
		return errors.New("invalid client ID")
	}

	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return service.ErrGigNotFound
	}

	// Verify ownership
	if !gig.ClientID().Equals(clientID) {
		return service.ErrUnauthorized
	}

	if err := gig.Cancel(cmd.Reason); err != nil {
		return err
	}

	return h.gigRepo.SaveWithEvents(ctx, gig)
}

// HandleSubmitProposal submits a proposal for a gig
func (h *GigHandler) HandleSubmitProposal(ctx context.Context, cmd command.SubmitProposal) (*command.SubmitProposalResult, error) {
	gigID, err := cmd.GetGigID()
	if err != nil {
		return nil, errors.New("invalid gig ID")
	}

	hustlerID, err := cmd.GetHustlerID()
	if err != nil {
		return nil, errors.New("invalid hustler ID")
	}

	proposedPrice, err := cmd.GetProposedPrice()
	if err != nil {
		return nil, err
	}

	// Load gig
	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return nil, service.ErrGigNotFound
	}

	// Create proposal
	proposalID := valueobject.GenerateProposalID()
	proposal := aggregate.NewProposal(
		proposalID,
		hustlerID,
		cmd.CoverLetter,
		proposedPrice,
		cmd.DeliveryDays,
		cmd.Attachments,
	)

	// Submit proposal to gig
	if err := gig.SubmitProposal(proposal); err != nil {
		return nil, err
	}

	// Save gig with events
	if err := h.gigRepo.SaveWithEvents(ctx, gig); err != nil {
		return nil, err
	}

	return &command.SubmitProposalResult{
		ProposalID:    proposal.ID().String(),
		GigID:         gigID.String(),
		ProposedPrice: proposal.ProposedPrice().Amount(),
		DeliveryDays:  proposal.DeliveryDays(),
		Status:        string(proposal.Status()),
		CreatedAt:     proposal.CreatedAt(),
	}, nil
}

// HandleWithdrawProposal withdraws a proposal
func (h *GigHandler) HandleWithdrawProposal(ctx context.Context, cmd command.WithdrawProposal) error {
	gigID, err := valueobject.NewGigID(cmd.GigID)
	if err != nil {
		return errors.New("invalid gig ID")
	}

	proposalID, err := valueobject.NewProposalID(cmd.ProposalID)
	if err != nil {
		return errors.New("invalid proposal ID")
	}

	hustlerID, err := valueobject.NewUserID(cmd.HustlerID)
	if err != nil {
		return errors.New("invalid hustler ID")
	}

	gig, err := h.gigRepo.FindByID(ctx, gigID)
	if err != nil {
		return service.ErrGigNotFound
	}

	if err := gig.WithdrawProposal(proposalID, hustlerID); err != nil {
		return err
	}

	return h.gigRepo.SaveWithEvents(ctx, gig)
}

// ContractHandler handles contract-related commands
type ContractHandler struct {
	contractSvc *service.ContractService
	reviewSvc   *service.ReviewService
	contractRepo repository.ContractRepository
}

// NewContractHandler creates a new contract handler
func NewContractHandler(
	contractSvc *service.ContractService,
	reviewSvc *service.ReviewService,
	contractRepo repository.ContractRepository,
) *ContractHandler {
	return &ContractHandler{
		contractSvc:  contractSvc,
		reviewSvc:    reviewSvc,
		contractRepo: contractRepo,
	}
}

// HandleAcceptProposal accepts a proposal and creates a contract
func (h *ContractHandler) HandleAcceptProposal(ctx context.Context, cmd command.AcceptProposal) (*command.AcceptProposalResult, error) {
	gigID, err := cmd.GetGigID()
	if err != nil {
		return nil, errors.New("invalid gig ID")
	}

	proposalID, err := cmd.GetProposalID()
	if err != nil {
		return nil, errors.New("invalid proposal ID")
	}

	clientID, err := cmd.GetClientID()
	if err != nil {
		return nil, errors.New("invalid client ID")
	}

	result, err := h.contractSvc.AcceptProposal(ctx, service.AcceptProposalRequest{
		GigID:      gigID,
		ProposalID: proposalID,
		ClientID:   clientID,
	})
	if err != nil {
		return nil, err
	}

	return &command.AcceptProposalResult{
		ContractID:   result.Contract.ID().String(),
		GigID:        result.Contract.GigID().String(),
		HustlerID:    result.Contract.HustlerID().String(),
		AgreedPrice:  result.Contract.AgreedPrice().Amount(),
		PlatformFee:  result.Contract.PlatformFee().Amount(),
		DeliveryDays: result.Contract.DeliveryDays(),
		DeadlineAt:   result.Contract.DeadlineAt(),
		Status:       result.Contract.Status().String(),
	}, nil
}

// HandleDeliverWork marks work as delivered
func (h *ContractHandler) HandleDeliverWork(ctx context.Context, cmd command.DeliverWork) (*command.DeliverWorkResult, error) {
	contractID, err := cmd.GetContractID()
	if err != nil {
		return nil, errors.New("invalid contract ID")
	}

	hustlerID, err := cmd.GetHustlerID()
	if err != nil {
		return nil, errors.New("invalid hustler ID")
	}

	contract, err := h.contractRepo.FindByID(ctx, contractID)
	if err != nil {
		return nil, service.ErrContractNotFound
	}

	if err := contract.Deliver(hustlerID, cmd.Deliverables); err != nil {
		return nil, err
	}

	if err := h.contractRepo.SaveWithEvents(ctx, contract); err != nil {
		return nil, err
	}

	return &command.DeliverWorkResult{
		ContractID:  contract.ID().String(),
		Status:      contract.Status().String(),
		DeliveredAt: *contract.DeliveredAt(),
	}, nil
}

// HandleApproveDelivery approves delivered work
func (h *ContractHandler) HandleApproveDelivery(ctx context.Context, cmd command.ApproveDelivery) (*command.ApproveDeliveryResult, error) {
	contractID, err := cmd.GetContractID()
	if err != nil {
		return nil, errors.New("invalid contract ID")
	}

	clientID, err := cmd.GetClientID()
	if err != nil {
		return nil, errors.New("invalid client ID")
	}

	contract, err := h.contractSvc.CompleteContract(ctx, service.CompleteContractRequest{
		ContractID: contractID,
		ClientID:   clientID,
		Notes:      cmd.Notes,
	})
	if err != nil {
		return nil, err
	}

	return &command.ApproveDeliveryResult{
		ContractID:  contract.ID().String(),
		Status:      contract.Status().String(),
		CompletedAt: *contract.CompletedAt(),
		PaidAmount:  contract.NetPayoutAmount().Amount(),
	}, nil
}

// HandleDisputeContract raises a dispute
func (h *ContractHandler) HandleDisputeContract(ctx context.Context, cmd command.DisputeContract) error {
	contractID, err := valueobject.NewContractID(cmd.ContractID)
	if err != nil {
		return errors.New("invalid contract ID")
	}

	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	contract, err := h.contractRepo.FindByID(ctx, contractID)
	if err != nil {
		return service.ErrContractNotFound
	}

	if err := contract.Dispute(userID, cmd.Reason); err != nil {
		return err
	}

	return h.contractRepo.SaveWithEvents(ctx, contract)
}

// HandleCancelContract cancels a contract
func (h *ContractHandler) HandleCancelContract(ctx context.Context, cmd command.CancelContract) error {
	contractID, err := valueobject.NewContractID(cmd.ContractID)
	if err != nil {
		return errors.New("invalid contract ID")
	}

	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	_, err = h.contractSvc.CancelContract(ctx, service.CancelContractRequest{
		ContractID: contractID,
		UserID:     userID,
		Reason:     cmd.Reason,
	})

	return err
}

// HandleSubmitReview submits a review
func (h *ContractHandler) HandleSubmitReview(ctx context.Context, cmd command.SubmitReview) (*command.SubmitReviewResult, error) {
	contractID, err := cmd.GetContractID()
	if err != nil {
		return nil, errors.New("invalid contract ID")
	}

	reviewerID, err := cmd.GetReviewerID()
	if err != nil {
		return nil, errors.New("invalid reviewer ID")
	}

	review, err := h.reviewSvc.SubmitReview(ctx, service.SubmitReviewRequest{
		ContractID:          contractID,
		ReviewerID:          reviewerID,
		Rating:              cmd.Rating,
		ReviewText:          cmd.ReviewText,
		CommunicationRating: cmd.CommunicationRating,
		QualityRating:       cmd.QualityRating,
		TimelinessRating:    cmd.TimelinessRating,
	})
	if err != nil {
		return nil, err
	}

	return &command.SubmitReviewResult{
		ReviewID:   review.ID(),
		ContractID: contractID.String(),
		Rating:     review.Rating(),
		CreatedAt:  time.Now().UTC(),
	}, nil
}
