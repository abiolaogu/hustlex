package service

import (
	"context"
	"errors"

	"hustlex/internal/domain/gig/aggregate"
	"hustlex/internal/domain/gig/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// Domain errors
var (
	ErrGigNotFound         = errors.New("gig not found")
	ErrContractNotFound    = errors.New("contract not found")
	ErrUnauthorized        = errors.New("unauthorized to perform this action")
	ErrEscrowFailed        = errors.New("failed to hold funds in escrow")
	ErrPaymentReleaseFailed = errors.New("failed to release payment")
)

// EscrowService defines the interface for escrow operations
// This is a PORT - wallet bounded context provides the ADAPTER
type EscrowService interface {
	// HoldFunds holds funds in escrow for a contract
	HoldFunds(ctx context.Context, userID valueobject.UserID, contractID valueobject.ContractID, amount valueobject.Money, description string) error

	// ReleaseFunds releases funds from escrow to the recipient
	ReleaseFunds(ctx context.Context, payerID, recipientID valueobject.UserID, contractID valueobject.ContractID, amount, platformFee valueobject.Money) error

	// RefundFunds refunds escrowed funds to the original holder
	RefundFunds(ctx context.Context, userID valueobject.UserID, contractID valueobject.ContractID, amount valueobject.Money, reason string) error
}

// ContractService handles contract-related domain operations
type ContractService struct {
	gigRepo      repository.GigRepository
	contractRepo repository.ContractRepository
	escrowSvc    EscrowService
}

// NewContractService creates a new contract service
func NewContractService(
	gigRepo repository.GigRepository,
	contractRepo repository.ContractRepository,
	escrowSvc EscrowService,
) *ContractService {
	return &ContractService{
		gigRepo:      gigRepo,
		contractRepo: contractRepo,
		escrowSvc:    escrowSvc,
	}
}

// AcceptProposalRequest contains the data needed to accept a proposal
type AcceptProposalRequest struct {
	GigID      valueobject.GigID
	ProposalID valueobject.ProposalID
	ClientID   valueobject.UserID // must be the gig owner
}

// AcceptProposalResult contains the result of accepting a proposal
type AcceptProposalResult struct {
	Contract *aggregate.Contract
	Gig      *aggregate.Gig
}

// AcceptProposal accepts a proposal and creates a contract
// This is a domain service because it orchestrates multiple aggregates
func (s *ContractService) AcceptProposal(ctx context.Context, req AcceptProposalRequest) (*AcceptProposalResult, error) {
	// Load gig
	gig, err := s.gigRepo.FindByID(ctx, req.GigID)
	if err != nil {
		return nil, ErrGigNotFound
	}

	// Verify ownership
	if !gig.ClientID().Equals(req.ClientID) {
		return nil, ErrUnauthorized
	}

	// Accept proposal in gig aggregate
	contractID := valueobject.GenerateContractID()
	acceptedData, err := gig.AcceptProposal(req.ProposalID, contractID)
	if err != nil {
		return nil, err
	}

	// Create contract aggregate
	contract, err := aggregate.NewContract(acceptedData, req.ProposalID)
	if err != nil {
		return nil, err
	}

	// Hold funds in escrow
	err = s.escrowSvc.HoldFunds(
		ctx,
		req.ClientID,
		contractID,
		acceptedData.AgreedPrice,
		"Escrow for gig: "+gig.Title(),
	)
	if err != nil {
		return nil, ErrEscrowFailed
	}

	// Save gig with events
	if err := s.gigRepo.SaveWithEvents(ctx, gig); err != nil {
		// TODO: Release escrow on failure (compensating transaction)
		return nil, err
	}

	// Save contract with events
	if err := s.contractRepo.SaveWithEvents(ctx, contract); err != nil {
		// TODO: Release escrow and revert gig status (compensating transaction)
		return nil, err
	}

	return &AcceptProposalResult{
		Contract: contract,
		Gig:      gig,
	}, nil
}

// CompleteContractRequest contains data to complete a contract
type CompleteContractRequest struct {
	ContractID valueobject.ContractID
	ClientID   valueobject.UserID
	Notes      string
}

// CompleteContract approves delivery and releases payment
func (s *ContractService) CompleteContract(ctx context.Context, req CompleteContractRequest) (*aggregate.Contract, error) {
	// Load contract
	contract, err := s.contractRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, ErrContractNotFound
	}

	// Approve work
	if err := contract.Approve(req.ClientID, req.Notes); err != nil {
		return nil, err
	}

	// Release payment from escrow
	err = s.escrowSvc.ReleaseFunds(
		ctx,
		contract.ClientID(),
		contract.HustlerID(),
		contract.ID(),
		contract.AgreedPrice(),
		contract.PlatformFee(),
	)
	if err != nil {
		return nil, ErrPaymentReleaseFailed
	}

	// Load and update gig status
	gig, err := s.gigRepo.FindByID(ctx, contract.GigID())
	if err == nil {
		gig.MarkCompleted()
		_ = s.gigRepo.SaveWithEvents(ctx, gig)
	}

	// Save contract with events
	if err := s.contractRepo.SaveWithEvents(ctx, contract); err != nil {
		// TODO: Handle failed save (escrow already released - log for manual review)
		return nil, err
	}

	return contract, nil
}

// CancelContractRequest contains data to cancel a contract
type CancelContractRequest struct {
	ContractID valueobject.ContractID
	UserID     valueobject.UserID
	Reason     string
}

// CancelContract cancels a contract and processes refund
func (s *ContractService) CancelContract(ctx context.Context, req CancelContractRequest) (*aggregate.Contract, error) {
	contract, err := s.contractRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, ErrContractNotFound
	}

	// Calculate refund amount based on who cancelled and contract state
	var refundAmount int64
	if contract.IsActive() {
		// Full refund if not yet delivered
		refundAmount = contract.AgreedPrice().Amount()
	} else if contract.IsDelivered() {
		// Partial refund or dispute resolution needed
		// For now, just mark as cancelled - admin will handle
		refundAmount = 0
	}

	// Cancel contract
	if err := contract.Cancel(req.UserID, req.Reason, refundAmount); err != nil {
		return nil, err
	}

	// Process refund if applicable
	if refundAmount > 0 {
		refundMoney := valueobject.MustNewMoney(refundAmount, contract.AgreedPrice().Currency())
		_ = s.escrowSvc.RefundFunds(
			ctx,
			contract.ClientID(),
			contract.ID(),
			refundMoney,
			"Contract cancelled: "+req.Reason,
		)
	}

	// Save contract
	if err := s.contractRepo.SaveWithEvents(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// ReviewService handles review-related operations
type ReviewService struct {
	contractRepo repository.ContractRepository
	reviewRepo   repository.ReviewRepository
}

// NewReviewService creates a new review service
func NewReviewService(
	contractRepo repository.ContractRepository,
	reviewRepo repository.ReviewRepository,
) *ReviewService {
	return &ReviewService{
		contractRepo: contractRepo,
		reviewRepo:   reviewRepo,
	}
}

// SubmitReviewRequest contains data to submit a review
type SubmitReviewRequest struct {
	ContractID          valueobject.ContractID
	ReviewerID          valueobject.UserID
	Rating              int
	ReviewText          string
	CommunicationRating int
	QualityRating       int
	TimelinessRating    int
}

// SubmitReview submits a review for a completed contract
func (s *ReviewService) SubmitReview(ctx context.Context, req SubmitReviewRequest) (*aggregate.Review, error) {
	contract, err := s.contractRepo.FindByID(ctx, req.ContractID)
	if err != nil {
		return nil, ErrContractNotFound
	}

	// Determine reviewee
	var revieweeID valueobject.UserID
	if contract.ClientID().Equals(req.ReviewerID) {
		// Client reviewing hustler
		revieweeID = contract.HustlerID()
	} else if contract.HustlerID().Equals(req.ReviewerID) {
		// Hustler reviewing client
		revieweeID = contract.ClientID()
	} else {
		return nil, ErrUnauthorized
	}

	// Create review
	reviewID := valueobject.GenerateTransactionID().String() // Using TransactionID generator for simplicity
	review, err := aggregate.NewReview(reviewID, req.ReviewerID, revieweeID, req.Rating, req.ReviewText)
	if err != nil {
		return nil, err
	}

	// Set detailed ratings
	review.SetDetailedRatings(req.CommunicationRating, req.QualityRating, req.TimelinessRating)

	// Add review to contract
	if err := contract.AddReview(review); err != nil {
		return nil, err
	}

	// Save
	if err := s.contractRepo.SaveWithEvents(ctx, contract); err != nil {
		return nil, err
	}

	return review, nil
}
