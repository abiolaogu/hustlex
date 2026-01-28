package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RemittanceStatus represents the status of a remittance transfer
type RemittanceStatus string

const (
	RemittanceStatusPending     RemittanceStatus = "pending"
	RemittanceStatusQuoted      RemittanceStatus = "quoted"
	RemittanceStatusInitiated   RemittanceStatus = "initiated"
	RemittanceStatusProcessing  RemittanceStatus = "processing"
	RemittanceStatusInTransit   RemittanceStatus = "in_transit"
	RemittanceStatusDelivered   RemittanceStatus = "delivered"
	RemittanceStatusCompleted   RemittanceStatus = "completed"
	RemittanceStatusFailed      RemittanceStatus = "failed"
	RemittanceStatusCancelled   RemittanceStatus = "cancelled"
	RemittanceStatusRefunded    RemittanceStatus = "refunded"
	RemittanceStatusOnHold      RemittanceStatus = "on_hold"
)

// RemittancePurpose represents the purpose of the transfer
type RemittancePurpose string

const (
	PurposeFamilySupport   RemittancePurpose = "family_support"
	PurposeEducation       RemittancePurpose = "education"
	PurposeMedical         RemittancePurpose = "medical"
	PurposeInvestment      RemittancePurpose = "investment"
	PurposePropertyPurchase RemittancePurpose = "property_purchase"
	PurposeBusiness        RemittancePurpose = "business"
	PurposeGift            RemittancePurpose = "gift"
	PurposeSalary          RemittancePurpose = "salary"
	PurposeOther           RemittancePurpose = "other"
)

// DeliveryMethod represents how the money will be delivered
type DeliveryMethod string

const (
	DeliveryBankTransfer   DeliveryMethod = "bank_transfer"
	DeliveryMobileWallet   DeliveryMethod = "mobile_wallet"
	DeliveryCashPickup     DeliveryMethod = "cash_pickup"
	DeliveryHomeDelivery   DeliveryMethod = "home_delivery"
	DeliveryHustleXWallet  DeliveryMethod = "hustlex_wallet"
)

// RecurrenceType represents the recurrence pattern
type RecurrenceType string

const (
	RecurrenceNone      RecurrenceType = "none"
	RecurrenceWeekly    RecurrenceType = "weekly"
	RecurrenceBiWeekly  RecurrenceType = "bi_weekly"
	RecurrenceMonthly   RecurrenceType = "monthly"
	RecurrenceQuarterly RecurrenceType = "quarterly"
)

// Remittance represents a money transfer
type Remittance struct {
	ID                  uuid.UUID         `json:"id" db:"id"`
	UserID              uuid.UUID         `json:"user_id" db:"user_id"`
	BeneficiaryID       uuid.UUID         `json:"beneficiary_id" db:"beneficiary_id"`
	SourceWalletID      *uuid.UUID        `json:"source_wallet_id,omitempty" db:"source_wallet_id"`

	// Reference & Status
	Reference           string            `json:"reference" db:"reference"`
	ExternalReference   string            `json:"external_reference,omitempty" db:"external_reference"`
	Status              RemittanceStatus  `json:"status" db:"status"`
	StatusMessage       string            `json:"status_message,omitempty" db:"status_message"`

	// Currencies & Amounts
	SourceCurrency      string            `json:"source_currency" db:"source_currency"`
	TargetCurrency      string            `json:"target_currency" db:"target_currency"`
	SourceAmount        decimal.Decimal   `json:"source_amount" db:"source_amount"`
	TargetAmount        decimal.Decimal   `json:"target_amount" db:"target_amount"`
	FXRate              decimal.Decimal   `json:"fx_rate" db:"fx_rate"`
	FXQuoteID           string            `json:"fx_quote_id,omitempty" db:"fx_quote_id"`

	// Fees
	TransferFee         decimal.Decimal   `json:"transfer_fee" db:"transfer_fee"`
	FXFee               decimal.Decimal   `json:"fx_fee" db:"fx_fee"`
	TotalFee            decimal.Decimal   `json:"total_fee" db:"total_fee"`
	TotalSourceAmount   decimal.Decimal   `json:"total_source_amount" db:"total_source_amount"`

	// Purpose & Delivery
	Purpose             RemittancePurpose `json:"purpose" db:"purpose"`
	PurposeDescription  string            `json:"purpose_description,omitempty" db:"purpose_description"`
	DeliveryMethod      DeliveryMethod    `json:"delivery_method" db:"delivery_method"`

	// Recurring Transfer
	IsRecurring         bool              `json:"is_recurring" db:"is_recurring"`
	RecurrenceType      RecurrenceType    `json:"recurrence_type" db:"recurrence_type"`
	RecurrenceStartDate *time.Time        `json:"recurrence_start_date,omitempty" db:"recurrence_start_date"`
	RecurrenceEndDate   *time.Time        `json:"recurrence_end_date,omitempty" db:"recurrence_end_date"`
	NextRecurrenceDate  *time.Time        `json:"next_recurrence_date,omitempty" db:"next_recurrence_date"`
	ParentRemittanceID  *uuid.UUID        `json:"parent_remittance_id,omitempty" db:"parent_remittance_id"`

	// Tracking
	EstimatedDelivery   *time.Time        `json:"estimated_delivery,omitempty" db:"estimated_delivery"`
	DeliveredAt         *time.Time        `json:"delivered_at,omitempty" db:"delivered_at"`
	CompletedAt         *time.Time        `json:"completed_at,omitempty" db:"completed_at"`

	// Source Payment
	PaymentMethod       string            `json:"payment_method,omitempty" db:"payment_method"`
	PaymentReference    string            `json:"payment_reference,omitempty" db:"payment_reference"`
	PaymentProvider     string            `json:"payment_provider,omitempty" db:"payment_provider"`
	PaidAt              *time.Time        `json:"paid_at,omitempty" db:"paid_at"`

	// Compliance
	ComplianceStatus    string            `json:"compliance_status" db:"compliance_status"`
	ComplianceNotes     string            `json:"compliance_notes,omitempty" db:"compliance_notes"`
	AMLChecked          bool              `json:"aml_checked" db:"aml_checked"`
	AMLCheckedAt        *time.Time        `json:"aml_checked_at,omitempty" db:"aml_checked_at"`

	// Notifications
	SenderNotified      bool              `json:"sender_notified" db:"sender_notified"`
	BeneficiaryNotified bool              `json:"beneficiary_notified" db:"beneficiary_notified"`

	// Audit
	CreatedAt           time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at" db:"updated_at"`
	CancelledAt         *time.Time        `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CancellationReason  string            `json:"cancellation_reason,omitempty" db:"cancellation_reason"`
}

// NewRemittance creates a new remittance
func NewRemittance(userID, beneficiaryID uuid.UUID, sourceCurrency, targetCurrency string, purpose RemittancePurpose, deliveryMethod DeliveryMethod) *Remittance {
	now := time.Now()
	return &Remittance{
		ID:               uuid.New(),
		UserID:           userID,
		BeneficiaryID:    beneficiaryID,
		Reference:        generateRemittanceReference(),
		Status:           RemittanceStatusPending,
		SourceCurrency:   sourceCurrency,
		TargetCurrency:   targetCurrency,
		Purpose:          purpose,
		DeliveryMethod:   deliveryMethod,
		RecurrenceType:   RecurrenceNone,
		ComplianceStatus: "pending",
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

// generateRemittanceReference generates a unique remittance reference
func generateRemittanceReference() string {
	return "RMT" + time.Now().Format("20060102") + uuid.New().String()[:8]
}

// SetAmounts sets the remittance amounts
func (r *Remittance) SetAmounts(sourceAmount, targetAmount, fxRate, transferFee, fxFee decimal.Decimal) {
	r.SourceAmount = sourceAmount
	r.TargetAmount = targetAmount
	r.FXRate = fxRate
	r.TransferFee = transferFee
	r.FXFee = fxFee
	r.TotalFee = transferFee.Add(fxFee)
	r.TotalSourceAmount = sourceAmount.Add(r.TotalFee)
	r.UpdatedAt = time.Now()
}

// SetQuote sets the FX quote for the remittance
func (r *Remittance) SetQuote(quoteID string, rate decimal.Decimal, estimatedDelivery time.Time) {
	r.FXQuoteID = quoteID
	r.FXRate = rate
	r.EstimatedDelivery = &estimatedDelivery
	r.Status = RemittanceStatusQuoted
	r.UpdatedAt = time.Now()
}

// Initiate initiates the remittance transfer
func (r *Remittance) Initiate(paymentMethod, paymentReference, paymentProvider string) error {
	if r.Status != RemittanceStatusQuoted && r.Status != RemittanceStatusPending {
		return errors.New("remittance cannot be initiated in current status")
	}
	now := time.Now()
	r.Status = RemittanceStatusInitiated
	r.PaymentMethod = paymentMethod
	r.PaymentReference = paymentReference
	r.PaymentProvider = paymentProvider
	r.PaidAt = &now
	r.UpdatedAt = now
	return nil
}

// StartProcessing moves the remittance to processing
func (r *Remittance) StartProcessing(externalRef string) error {
	if r.Status != RemittanceStatusInitiated {
		return errors.New("remittance must be initiated before processing")
	}
	r.Status = RemittanceStatusProcessing
	r.ExternalReference = externalRef
	r.UpdatedAt = time.Now()
	return nil
}

// MarkInTransit marks the remittance as in transit
func (r *Remittance) MarkInTransit(statusMessage string) error {
	if r.Status != RemittanceStatusProcessing {
		return errors.New("remittance must be processing before transit")
	}
	r.Status = RemittanceStatusInTransit
	r.StatusMessage = statusMessage
	r.UpdatedAt = time.Now()
	return nil
}

// MarkDelivered marks the remittance as delivered
func (r *Remittance) MarkDelivered() error {
	if r.Status != RemittanceStatusInTransit && r.Status != RemittanceStatusProcessing {
		return errors.New("remittance must be in transit or processing")
	}
	now := time.Now()
	r.Status = RemittanceStatusDelivered
	r.DeliveredAt = &now
	r.UpdatedAt = now
	return nil
}

// Complete marks the remittance as completed
func (r *Remittance) Complete() error {
	if r.Status != RemittanceStatusDelivered {
		return errors.New("remittance must be delivered before completion")
	}
	now := time.Now()
	r.Status = RemittanceStatusCompleted
	r.CompletedAt = &now
	r.UpdatedAt = now
	return nil
}

// Fail marks the remittance as failed
func (r *Remittance) Fail(reason string) {
	r.Status = RemittanceStatusFailed
	r.StatusMessage = reason
	r.UpdatedAt = time.Now()
}

// Cancel cancels the remittance
func (r *Remittance) Cancel(reason string) error {
	if r.Status == RemittanceStatusCompleted || r.Status == RemittanceStatusDelivered {
		return errors.New("cannot cancel completed or delivered remittance")
	}
	now := time.Now()
	r.Status = RemittanceStatusCancelled
	r.CancelledAt = &now
	r.CancellationReason = reason
	r.UpdatedAt = now
	return nil
}

// Hold puts the remittance on hold
func (r *Remittance) Hold(reason string) {
	r.Status = RemittanceStatusOnHold
	r.StatusMessage = reason
	r.UpdatedAt = time.Now()
}

// PassAMLCheck marks AML check as passed
func (r *Remittance) PassAMLCheck() {
	now := time.Now()
	r.AMLChecked = true
	r.AMLCheckedAt = &now
	r.ComplianceStatus = "passed"
	r.UpdatedAt = now
}

// SetupRecurring sets up recurring transfer
func (r *Remittance) SetupRecurring(recurrenceType RecurrenceType, startDate, endDate time.Time) {
	r.IsRecurring = true
	r.RecurrenceType = recurrenceType
	r.RecurrenceStartDate = &startDate
	r.RecurrenceEndDate = &endDate
	r.NextRecurrenceDate = r.calculateNextRecurrence(startDate)
	r.UpdatedAt = time.Now()
}

// calculateNextRecurrence calculates the next recurrence date
func (r *Remittance) calculateNextRecurrence(from time.Time) *time.Time {
	var next time.Time
	switch r.RecurrenceType {
	case RecurrenceWeekly:
		next = from.AddDate(0, 0, 7)
	case RecurrenceBiWeekly:
		next = from.AddDate(0, 0, 14)
	case RecurrenceMonthly:
		next = from.AddDate(0, 1, 0)
	case RecurrenceQuarterly:
		next = from.AddDate(0, 3, 0)
	default:
		return nil
	}
	return &next
}

// CanRefund checks if remittance can be refunded
func (r *Remittance) CanRefund() bool {
	return r.Status == RemittanceStatusFailed ||
		r.Status == RemittanceStatusCancelled ||
		r.Status == RemittanceStatusOnHold
}

// NotifySender marks sender as notified
func (r *Remittance) NotifySender() {
	r.SenderNotified = true
	r.UpdatedAt = time.Now()
}

// NotifyBeneficiary marks beneficiary as notified
func (r *Remittance) NotifyBeneficiary() {
	r.BeneficiaryNotified = true
	r.UpdatedAt = time.Now()
}
