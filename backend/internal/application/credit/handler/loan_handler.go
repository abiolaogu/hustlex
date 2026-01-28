package handler

import (
	"context"
	"errors"

	"hustlex/internal/application/credit/command"
	"hustlex/internal/domain/credit/aggregate"
	"hustlex/internal/domain/credit/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrLoanNotFound       = errors.New("loan not found")
	ErrCreditScoreNotFound = errors.New("credit score not found")
	ErrUnauthorized       = errors.New("unauthorized to perform this action")
	ErrInsufficientCredit = errors.New("insufficient credit score for this loan")
)

// LoanHandler handles loan-related commands
type LoanHandler struct {
	loanRepo       repository.LoanRepository
	creditScoreRepo repository.CreditScoreRepository
}

// NewLoanHandler creates a new loan handler
func NewLoanHandler(
	loanRepo repository.LoanRepository,
	creditScoreRepo repository.CreditScoreRepository,
) *LoanHandler {
	return &LoanHandler{
		loanRepo:       loanRepo,
		creditScoreRepo: creditScoreRepo,
	}
}

// HandleApplyForLoan processes a loan application
func (h *LoanHandler) HandleApplyForLoan(ctx context.Context, cmd command.ApplyForLoan) (*command.ApplyForLoanResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	amount, err := cmd.GetAmount()
	if err != nil {
		return nil, err
	}

	// Check for existing active loan
	existingLoan, _ := h.loanRepo.FindActiveByUserID(ctx, userID)
	if existingLoan != nil {
		return nil, aggregate.ErrActiveLoanExists
	}

	// Get user's credit score
	creditScore, err := h.creditScoreRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, ErrCreditScoreNotFound
	}

	// Check if loan amount is within allowed limit
	maxAllowed := creditScore.MaxLoanAmount()
	if amount.Amount() > maxAllowed {
		return nil, aggregate.ErrLoanExceedsLimit
	}

	// Get interest rate based on tier
	interestRate := creditScore.InterestRate()

	// Create loan
	loanID := valueobject.GenerateLoanID()
	loan, err := aggregate.NewLoan(
		loanID,
		userID,
		amount,
		interestRate,
		cmd.TenureMonths,
		cmd.Purpose,
		maxAllowed,
	)
	if err != nil {
		return nil, err
	}

	// Save loan
	if err := h.loanRepo.SaveWithEvents(ctx, loan); err != nil {
		return nil, err
	}

	return &command.ApplyForLoanResult{
		LoanID:         loan.ID().String(),
		Principal:      loan.Principal().Amount(),
		InterestRate:   loan.InterestRate(),
		InterestAmount: loan.InterestAmount().Amount(),
		TotalAmount:    loan.TotalAmount().Amount(),
		TenureMonths:   loan.TenureMonths(),
		MonthlyPayment: loan.MonthlyPayment().Amount(),
		Status:         loan.Status().String(),
	}, nil
}

// HandleApproveLoan approves a pending loan
func (h *LoanHandler) HandleApproveLoan(ctx context.Context, cmd command.ApproveLoan) error {
	loanID, err := cmd.GetLoanID()
	if err != nil {
		return errors.New("invalid loan ID")
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return ErrLoanNotFound
	}

	if err := loan.Approve(); err != nil {
		return err
	}

	return h.loanRepo.SaveWithEvents(ctx, loan)
}

// HandleRejectLoan rejects a pending loan
func (h *LoanHandler) HandleRejectLoan(ctx context.Context, cmd command.RejectLoan) error {
	loanID, err := cmd.GetLoanID()
	if err != nil {
		return errors.New("invalid loan ID")
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return ErrLoanNotFound
	}

	if err := loan.Reject(cmd.Reason); err != nil {
		return err
	}

	return h.loanRepo.SaveWithEvents(ctx, loan)
}

// HandleDisburseLoan disburses an approved loan
func (h *LoanHandler) HandleDisburseLoan(ctx context.Context, cmd command.DisburseLoan) (*command.DisburseLoanResult, error) {
	loanID, err := cmd.GetLoanID()
	if err != nil {
		return nil, errors.New("invalid loan ID")
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return nil, ErrLoanNotFound
	}

	if err := loan.Disburse(); err != nil {
		return nil, err
	}

	if err := h.loanRepo.SaveWithEvents(ctx, loan); err != nil {
		return nil, err
	}

	return &command.DisburseLoanResult{
		LoanID:      loan.ID().String(),
		Amount:      loan.Principal().Amount(),
		DisbursedAt: loan.DisbursedAt().Format("2006-01-02T15:04:05Z07:00"),
		DueDate:     loan.DueDate().Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// HandleRecordRepayment records a loan repayment
func (h *LoanHandler) HandleRecordRepayment(ctx context.Context, cmd command.RecordRepayment) (*command.RecordRepaymentResult, error) {
	loanID, err := cmd.GetLoanID()
	if err != nil {
		return nil, errors.New("invalid loan ID")
	}

	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	amount, err := cmd.GetAmount()
	if err != nil {
		return nil, err
	}

	transactionID, err := cmd.GetTransactionID()
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return nil, ErrLoanNotFound
	}

	// Verify ownership
	if loan.UserID() != userID {
		return nil, ErrUnauthorized
	}

	// Generate repayment ID
	repaymentID := valueobject.GenerateUserID().String() // Using UserID generator for unique ID

	if err := loan.RecordRepayment(repaymentID, amount, transactionID); err != nil {
		return nil, err
	}

	if err := h.loanRepo.SaveWithEvents(ctx, loan); err != nil {
		return nil, err
	}

	return &command.RecordRepaymentResult{
		RepaymentID:      repaymentID,
		LoanID:           loan.ID().String(),
		AmountPaid:       amount.Amount(),
		RemainingBalance: loan.RemainingBalance().Amount(),
		IsFullyRepaid:    loan.IsFullyRepaid(),
		Status:           loan.Status().String(),
	}, nil
}

// HandleMarkLoanDefaulted marks a loan as defaulted
func (h *LoanHandler) HandleMarkLoanDefaulted(ctx context.Context, cmd command.MarkLoanDefaulted) error {
	loanID, err := cmd.GetLoanID()
	if err != nil {
		return errors.New("invalid loan ID")
	}

	loan, err := h.loanRepo.FindByID(ctx, loanID)
	if err != nil {
		return ErrLoanNotFound
	}

	if err := loan.MarkDefaulted(); err != nil {
		return err
	}

	return h.loanRepo.SaveWithEvents(ctx, loan)
}
