package command

import (
	"hustlex/internal/domain/shared/valueobject"
)

// ApplyForLoan represents a loan application command
type ApplyForLoan struct {
	UserID       string
	Amount       int64
	Currency     string
	TenureMonths int
	Purpose      string
}

func (c ApplyForLoan) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c ApplyForLoan) GetAmount() (valueobject.Money, error) {
	return valueobject.NewMoney(c.Amount, valueobject.Currency(c.Currency))
}

// ApplyForLoanResult contains the result of loan application
type ApplyForLoanResult struct {
	LoanID         string
	Principal      int64
	InterestRate   float64
	InterestAmount int64
	TotalAmount    int64
	TenureMonths   int
	MonthlyPayment int64
	Status         string
}

// ApproveLoan represents a loan approval command
type ApproveLoan struct {
	LoanID  string
	AdminID string
}

func (c ApproveLoan) GetLoanID() (valueobject.LoanID, error) {
	return valueobject.NewLoanID(c.LoanID)
}

// RejectLoan represents a loan rejection command
type RejectLoan struct {
	LoanID  string
	AdminID string
	Reason  string
}

func (c RejectLoan) GetLoanID() (valueobject.LoanID, error) {
	return valueobject.NewLoanID(c.LoanID)
}

// DisburseLoan represents loan disbursement command
type DisburseLoan struct {
	LoanID        string
	AdminID       string
	TransactionID string
}

func (c DisburseLoan) GetLoanID() (valueobject.LoanID, error) {
	return valueobject.NewLoanID(c.LoanID)
}

func (c DisburseLoan) GetTransactionID() (valueobject.TransactionID, error) {
	return valueobject.NewTransactionID(c.TransactionID)
}

// DisburseLoanResult contains disbursement details
type DisburseLoanResult struct {
	LoanID      string
	Amount      int64
	DisbursedAt string
	DueDate     string
}

// RecordRepayment represents a loan repayment command
type RecordRepayment struct {
	LoanID        string
	UserID        string
	Amount        int64
	Currency      string
	TransactionID string
}

func (c RecordRepayment) GetLoanID() (valueobject.LoanID, error) {
	return valueobject.NewLoanID(c.LoanID)
}

func (c RecordRepayment) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c RecordRepayment) GetAmount() (valueobject.Money, error) {
	return valueobject.NewMoney(c.Amount, valueobject.Currency(c.Currency))
}

func (c RecordRepayment) GetTransactionID() (valueobject.TransactionID, error) {
	return valueobject.NewTransactionID(c.TransactionID)
}

// RecordRepaymentResult contains repayment details
type RecordRepaymentResult struct {
	RepaymentID      string
	LoanID           string
	AmountPaid       int64
	RemainingBalance int64
	IsFullyRepaid    bool
	Status           string
}

// MarkLoanDefaulted marks a loan as defaulted
type MarkLoanDefaulted struct {
	LoanID  string
	AdminID string
}

func (c MarkLoanDefaulted) GetLoanID() (valueobject.LoanID, error) {
	return valueobject.NewLoanID(c.LoanID)
}

// RecalculateCreditScore recalculates a user's credit score
type RecalculateCreditScore struct {
	UserID string
}

func (c RecalculateCreditScore) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// RecalculateCreditScoreResult contains the new credit score
type RecalculateCreditScoreResult struct {
	UserID        string
	Score         int
	Tier          string
	MaxLoanAmount int64
	InterestRate  float64
}

// UpdateCreditStats updates credit score components
type UpdateCreditStats struct {
	UserID              string
	GigsCompleted       *int
	GigsAccepted        *int
	AverageRating       *float64
	TotalReviews        *int
	OnTimeContributions *int
	TotalContributions  *int
	AccountAgeMonths    *int
	HasPhoneVerified    *bool
	HasEmailVerified    *bool
	HasBVNVerified      *bool
	HasNINVerified      *bool
	CirclesJoined       *int
	Referrals           *int
}

func (c UpdateCreditStats) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// InitializeCreditScore creates a credit score for a new user
type InitializeCreditScore struct {
	UserID string
}

func (c InitializeCreditScore) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}
