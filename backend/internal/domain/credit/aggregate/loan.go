package aggregate

import (
	"errors"
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Loan errors
var (
	ErrLoanExceedsLimit    = errors.New("loan amount exceeds maximum allowed for your tier")
	ErrInvalidLoanAmount   = errors.New("invalid loan amount")
	ErrInvalidTenure       = errors.New("tenure must be between 1 and 12 months")
	ErrLoanNotApproved     = errors.New("loan is not approved")
	ErrLoanNotDisbursed    = errors.New("loan has not been disbursed")
	ErrLoanAlreadyPaid     = errors.New("loan has already been fully repaid")
	ErrRepaymentExceeds    = errors.New("repayment amount exceeds remaining balance")
	ErrActiveLoanExists    = errors.New("user already has an active loan")
)

// LoanStatus represents the status of a loan
type LoanStatus string

const (
	LoanStatusPending   LoanStatus = "pending"
	LoanStatusApproved  LoanStatus = "approved"
	LoanStatusDisbursed LoanStatus = "disbursed"
	LoanStatusRepaying  LoanStatus = "repaying"
	LoanStatusCompleted LoanStatus = "completed"
	LoanStatusDefaulted LoanStatus = "defaulted"
	LoanStatusRejected  LoanStatus = "rejected"
)

func (s LoanStatus) String() string {
	return string(s)
}

// Repayment represents a loan repayment
type Repayment struct {
	id            string
	amount        valueobject.Money
	transactionID valueobject.TransactionID
	paidAt        time.Time
}

func NewRepayment(id string, amount valueobject.Money, transactionID valueobject.TransactionID) *Repayment {
	return &Repayment{
		id:            id,
		amount:        amount,
		transactionID: transactionID,
		paidAt:        time.Now().UTC(),
	}
}

func (r *Repayment) ID() string                         { return r.id }
func (r *Repayment) Amount() valueobject.Money          { return r.amount }
func (r *Repayment) TransactionID() valueobject.TransactionID { return r.transactionID }
func (r *Repayment) PaidAt() time.Time                  { return r.paidAt }

// Loan is the aggregate root for microloans
type Loan struct {
	sharedevent.AggregateRoot

	id             valueobject.LoanID
	userID         valueobject.UserID
	principal      valueobject.Money
	interestRate   float64 // Monthly rate
	interestAmount valueobject.Money
	totalAmount    valueobject.Money
	amountRepaid   valueobject.Money
	tenureMonths   int
	status         LoanStatus
	purpose        string
	approvedAt     *time.Time
	disbursedAt    *time.Time
	dueDate        *time.Time
	completedAt    *time.Time
	repayments     []*Repayment
	createdAt      time.Time
	updatedAt      time.Time
	version        int64
}

// NewLoan creates a new loan application
func NewLoan(
	id valueobject.LoanID,
	userID valueobject.UserID,
	principal valueobject.Money,
	interestRate float64,
	tenureMonths int,
	purpose string,
	maxAllowed int64,
) (*Loan, error) {
	if principal.Amount() <= 0 {
		return nil, ErrInvalidLoanAmount
	}

	if principal.Amount() > maxAllowed {
		return nil, ErrLoanExceedsLimit
	}

	if tenureMonths < 1 || tenureMonths > 12 {
		return nil, ErrInvalidTenure
	}

	// Calculate interest
	interestAmount := int64(float64(principal.Amount()) * interestRate * float64(tenureMonths))
	interest, _ := valueobject.NewMoney(interestAmount, principal.Currency())
	totalAmount := principal.MustAdd(interest)

	loan := &Loan{
		id:             id,
		userID:         userID,
		principal:      principal,
		interestRate:   interestRate,
		interestAmount: interest,
		totalAmount:    totalAmount,
		amountRepaid:   valueobject.MustNewMoney(0, principal.Currency()),
		tenureMonths:   tenureMonths,
		status:         LoanStatusPending,
		purpose:        purpose,
		repayments:     make([]*Repayment, 0),
		createdAt:      time.Now().UTC(),
		updatedAt:      time.Now().UTC(),
		version:        1,
	}

	return loan, nil
}

// ReconstructLoan reconstructs from persistence
func ReconstructLoan(
	id valueobject.LoanID,
	userID valueobject.UserID,
	principal valueobject.Money,
	interestRate float64,
	interestAmount valueobject.Money,
	totalAmount valueobject.Money,
	amountRepaid valueobject.Money,
	tenureMonths int,
	status LoanStatus,
	purpose string,
	approvedAt *time.Time,
	disbursedAt *time.Time,
	dueDate *time.Time,
	completedAt *time.Time,
	repayments []*Repayment,
	createdAt time.Time,
	updatedAt time.Time,
	version int64,
) *Loan {
	return &Loan{
		id:             id,
		userID:         userID,
		principal:      principal,
		interestRate:   interestRate,
		interestAmount: interestAmount,
		totalAmount:    totalAmount,
		amountRepaid:   amountRepaid,
		tenureMonths:   tenureMonths,
		status:         status,
		purpose:        purpose,
		approvedAt:     approvedAt,
		disbursedAt:    disbursedAt,
		dueDate:        dueDate,
		completedAt:    completedAt,
		repayments:     repayments,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		version:        version,
	}
}

// Getters
func (l *Loan) ID() valueobject.LoanID       { return l.id }
func (l *Loan) UserID() valueobject.UserID   { return l.userID }
func (l *Loan) Principal() valueobject.Money { return l.principal }
func (l *Loan) InterestRate() float64        { return l.interestRate }
func (l *Loan) InterestAmount() valueobject.Money { return l.interestAmount }
func (l *Loan) TotalAmount() valueobject.Money { return l.totalAmount }
func (l *Loan) AmountRepaid() valueobject.Money { return l.amountRepaid }
func (l *Loan) TenureMonths() int            { return l.tenureMonths }
func (l *Loan) Status() LoanStatus           { return l.status }
func (l *Loan) Purpose() string              { return l.purpose }
func (l *Loan) ApprovedAt() *time.Time       { return l.approvedAt }
func (l *Loan) DisbursedAt() *time.Time      { return l.disbursedAt }
func (l *Loan) DueDate() *time.Time          { return l.dueDate }
func (l *Loan) CompletedAt() *time.Time      { return l.completedAt }
func (l *Loan) Repayments() []*Repayment     { return l.repayments }
func (l *Loan) CreatedAt() time.Time         { return l.createdAt }
func (l *Loan) UpdatedAt() time.Time         { return l.updatedAt }
func (l *Loan) Version() int64               { return l.version }

// RemainingBalance returns the amount still owed
func (l *Loan) RemainingBalance() valueobject.Money {
	return l.totalAmount.MustSubtract(l.amountRepaid)
}

// IsFullyRepaid checks if the loan is fully repaid
func (l *Loan) IsFullyRepaid() bool {
	return l.amountRepaid.Amount() >= l.totalAmount.Amount()
}

// IsOverdue checks if the loan is past due
func (l *Loan) IsOverdue() bool {
	if l.dueDate == nil {
		return false
	}
	return time.Now().After(*l.dueDate) && !l.IsFullyRepaid() && l.status == LoanStatusRepaying
}

// Business Methods

// Approve approves the loan application
func (l *Loan) Approve() error {
	if l.status != LoanStatusPending {
		return errors.New("can only approve pending loans")
	}

	now := time.Now().UTC()
	l.status = LoanStatusApproved
	l.approvedAt = &now
	l.updatedAt = now

	return nil
}

// Reject rejects the loan application
func (l *Loan) Reject(reason string) error {
	if l.status != LoanStatusPending {
		return errors.New("can only reject pending loans")
	}

	l.status = LoanStatusRejected
	l.updatedAt = time.Now().UTC()

	return nil
}

// Disburse marks the loan as disbursed
func (l *Loan) Disburse() error {
	if l.status != LoanStatusApproved {
		return ErrLoanNotApproved
	}

	now := time.Now().UTC()
	dueDate := now.AddDate(0, l.tenureMonths, 0)

	l.status = LoanStatusDisbursed
	l.disbursedAt = &now
	l.dueDate = &dueDate
	l.updatedAt = now

	return nil
}

// StartRepayment transitions to repaying status after first payment
func (l *Loan) StartRepayment() {
	if l.status == LoanStatusDisbursed {
		l.status = LoanStatusRepaying
		l.updatedAt = time.Now().UTC()
	}
}

// RecordRepayment records a loan repayment
func (l *Loan) RecordRepayment(repaymentID string, amount valueobject.Money, transactionID valueobject.TransactionID) error {
	if l.status != LoanStatusDisbursed && l.status != LoanStatusRepaying {
		return ErrLoanNotDisbursed
	}

	if l.IsFullyRepaid() {
		return ErrLoanAlreadyPaid
	}

	remaining := l.RemainingBalance()
	if amount.GreaterThan(remaining) {
		return ErrRepaymentExceeds
	}

	repayment := NewRepayment(repaymentID, amount, transactionID)
	l.repayments = append(l.repayments, repayment)
	l.amountRepaid = l.amountRepaid.MustAdd(amount)

	if l.status == LoanStatusDisbursed {
		l.status = LoanStatusRepaying
	}

	// Check if fully repaid
	if l.IsFullyRepaid() {
		now := time.Now().UTC()
		l.status = LoanStatusCompleted
		l.completedAt = &now
	}

	l.updatedAt = time.Now().UTC()

	return nil
}

// MarkDefaulted marks the loan as defaulted
func (l *Loan) MarkDefaulted() error {
	if l.status != LoanStatusRepaying {
		return errors.New("can only default loans in repayment")
	}

	l.status = LoanStatusDefaulted
	l.updatedAt = time.Now().UTC()

	return nil
}

// MonthlyPayment calculates the monthly payment amount
func (l *Loan) MonthlyPayment() valueobject.Money {
	if l.tenureMonths <= 0 {
		return l.totalAmount
	}
	monthlyAmount := l.totalAmount.Amount() / int64(l.tenureMonths)
	return valueobject.MustNewMoney(monthlyAmount, l.totalAmount.Currency())
}
