package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// BookingStatus represents the status of a diaspora booking
type BookingStatus string

const (
	BookingStatusPending      BookingStatus = "pending"
	BookingStatusQuoted       BookingStatus = "quoted"
	BookingStatusConfirmed    BookingStatus = "confirmed"
	BookingStatusPaid         BookingStatus = "paid"
	BookingStatusProcessing   BookingStatus = "processing"
	BookingStatusCompleted    BookingStatus = "completed"
	BookingStatusCancelled    BookingStatus = "cancelled"
	BookingStatusRefunded     BookingStatus = "refunded"
	BookingStatusFailed       BookingStatus = "failed"
)

// BookingType represents the type of diaspora booking
type BookingType string

const (
	BookingTypeService     BookingType = "service"      // Booking a service provider
	BookingTypeRemittance  BookingType = "remittance"   // Money transfer
	BookingTypeBillPayment BookingType = "bill_payment" // Pay bills in Nigeria
	BookingTypeAirtime     BookingType = "airtime"      // Buy airtime for family
	BookingTypeGiftCard    BookingType = "gift_card"    // Send gift cards
)

// DiasporaBooking represents a booking made by diaspora users
type DiasporaBooking struct {
	ID                  uuid.UUID       `json:"id" db:"id"`
	UserID              uuid.UUID       `json:"user_id" db:"user_id"`
	BeneficiaryID       *uuid.UUID      `json:"beneficiary_id,omitempty" db:"beneficiary_id"`
	ServiceProviderID   *uuid.UUID      `json:"service_provider_id,omitempty" db:"service_provider_id"`

	// Booking Details
	Type                BookingType     `json:"type" db:"type"`
	Status              BookingStatus   `json:"status" db:"status"`
	Reference           string          `json:"reference" db:"reference"`
	Description         string          `json:"description" db:"description"`

	// Currency & Amounts (Dual Currency Support)
	SourceCurrency      string          `json:"source_currency" db:"source_currency"`
	TargetCurrency      string          `json:"target_currency" db:"target_currency"`
	SourceAmount        decimal.Decimal `json:"source_amount" db:"source_amount"`
	TargetAmount        decimal.Decimal `json:"target_amount" db:"target_amount"`
	Fee                 decimal.Decimal `json:"fee" db:"fee"`
	FeeCurrency         string          `json:"fee_currency" db:"fee_currency"`
	TotalSourceAmount   decimal.Decimal `json:"total_source_amount" db:"total_source_amount"`

	// FX Rate Locking
	FXRate              decimal.Decimal `json:"fx_rate" db:"fx_rate"`
	FXQuoteID           string          `json:"fx_quote_id,omitempty" db:"fx_quote_id"`
	FXRateLockedAt      *time.Time      `json:"fx_rate_locked_at,omitempty" db:"fx_rate_locked_at"`
	FXRateExpiresAt     *time.Time      `json:"fx_rate_expires_at,omitempty" db:"fx_rate_expires_at"`
	FXRateLocked        bool            `json:"fx_rate_locked" db:"fx_rate_locked"`

	// Service Details (for service bookings)
	ServiceDate         *time.Time      `json:"service_date,omitempty" db:"service_date"`
	ServiceAddress      string          `json:"service_address,omitempty" db:"service_address"`
	ServiceCity         string          `json:"service_city,omitempty" db:"service_city"`
	ServiceState        string          `json:"service_state,omitempty" db:"service_state"`
	ServiceNotes        string          `json:"service_notes,omitempty" db:"service_notes"`

	// Verification
	RequiresVerification bool           `json:"requires_verification" db:"requires_verification"`
	VerificationCode    string          `json:"verification_code,omitempty" db:"verification_code"`
	VerifiedAt          *time.Time      `json:"verified_at,omitempty" db:"verified_at"`
	VerifiedByPhone     string          `json:"verified_by_phone,omitempty" db:"verified_by_phone"`

	// Payment
	PaymentMethod       string          `json:"payment_method,omitempty" db:"payment_method"`
	PaymentReference    string          `json:"payment_reference,omitempty" db:"payment_reference"`
	PaidAt              *time.Time      `json:"paid_at,omitempty" db:"paid_at"`

	// Fulfillment
	FulfilledAt         *time.Time      `json:"fulfilled_at,omitempty" db:"fulfilled_at"`
	FulfillmentNotes    string          `json:"fulfillment_notes,omitempty" db:"fulfillment_notes"`
	ProofOfDelivery     string          `json:"proof_of_delivery,omitempty" db:"proof_of_delivery"`

	// Metadata
	Metadata            map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// Audit
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
	CancelledAt         *time.Time      `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CancellationReason  string          `json:"cancellation_reason,omitempty" db:"cancellation_reason"`
}

// NewDiasporaBooking creates a new diaspora booking
func NewDiasporaBooking(userID uuid.UUID, bookingType BookingType, sourceCurrency, targetCurrency string) *DiasporaBooking {
	now := time.Now()
	return &DiasporaBooking{
		ID:             uuid.New(),
		UserID:         userID,
		Type:           bookingType,
		Status:         BookingStatusPending,
		Reference:      generateBookingReference(),
		SourceCurrency: sourceCurrency,
		TargetCurrency: targetCurrency,
		FeeCurrency:    sourceCurrency,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// generateBookingReference generates a unique booking reference
func generateBookingReference() string {
	return "DBK" + time.Now().Format("20060102") + uuid.New().String()[:8]
}

// SetAmounts sets the booking amounts
func (b *DiasporaBooking) SetAmounts(sourceAmount, targetAmount, fee decimal.Decimal) {
	b.SourceAmount = sourceAmount
	b.TargetAmount = targetAmount
	b.Fee = fee
	b.TotalSourceAmount = sourceAmount.Add(fee)
	b.UpdatedAt = time.Now()
}

// LockFXRate locks the FX rate for the booking
func (b *DiasporaBooking) LockFXRate(rate decimal.Decimal, quoteID string, expiresAt time.Time) {
	now := time.Now()
	b.FXRate = rate
	b.FXQuoteID = quoteID
	b.FXRateLockedAt = &now
	b.FXRateExpiresAt = &expiresAt
	b.FXRateLocked = true
	b.Status = BookingStatusQuoted
	b.UpdatedAt = now
}

// IsFXRateValid checks if the locked FX rate is still valid
func (b *DiasporaBooking) IsFXRateValid() bool {
	if !b.FXRateLocked || b.FXRateExpiresAt == nil {
		return false
	}
	return time.Now().Before(*b.FXRateExpiresAt)
}

// Confirm confirms the booking
func (b *DiasporaBooking) Confirm() error {
	if b.Status != BookingStatusQuoted && b.Status != BookingStatusPending {
		return errors.New("booking cannot be confirmed in current status")
	}
	b.Status = BookingStatusConfirmed
	b.UpdatedAt = time.Now()
	return nil
}

// MarkPaid marks the booking as paid
func (b *DiasporaBooking) MarkPaid(paymentMethod, paymentReference string) error {
	if b.Status != BookingStatusConfirmed {
		return errors.New("booking must be confirmed before payment")
	}
	now := time.Now()
	b.Status = BookingStatusPaid
	b.PaymentMethod = paymentMethod
	b.PaymentReference = paymentReference
	b.PaidAt = &now
	b.UpdatedAt = now
	return nil
}

// StartProcessing moves booking to processing state
func (b *DiasporaBooking) StartProcessing() error {
	if b.Status != BookingStatusPaid {
		return errors.New("booking must be paid before processing")
	}
	b.Status = BookingStatusProcessing
	b.UpdatedAt = time.Now()
	return nil
}

// Complete marks the booking as completed
func (b *DiasporaBooking) Complete(notes, proofOfDelivery string) error {
	if b.Status != BookingStatusProcessing {
		return errors.New("booking must be processing before completion")
	}
	now := time.Now()
	b.Status = BookingStatusCompleted
	b.FulfilledAt = &now
	b.FulfillmentNotes = notes
	b.ProofOfDelivery = proofOfDelivery
	b.UpdatedAt = now
	return nil
}

// Cancel cancels the booking
func (b *DiasporaBooking) Cancel(reason string) error {
	if b.Status == BookingStatusCompleted || b.Status == BookingStatusRefunded {
		return errors.New("cannot cancel completed or refunded booking")
	}
	now := time.Now()
	b.Status = BookingStatusCancelled
	b.CancelledAt = &now
	b.CancellationReason = reason
	b.UpdatedAt = now
	return nil
}

// SetVerification sets up verification requirements
func (b *DiasporaBooking) SetVerification(code string) {
	b.RequiresVerification = true
	b.VerificationCode = code
	b.UpdatedAt = time.Now()
}

// Verify verifies the booking with beneficiary
func (b *DiasporaBooking) Verify(phone string) error {
	if !b.RequiresVerification {
		return errors.New("booking does not require verification")
	}
	now := time.Now()
	b.VerifiedAt = &now
	b.VerifiedByPhone = phone
	b.UpdatedAt = now
	return nil
}

// SetServiceDetails sets service booking details
func (b *DiasporaBooking) SetServiceDetails(serviceDate time.Time, address, city, state, notes string) {
	b.ServiceDate = &serviceDate
	b.ServiceAddress = address
	b.ServiceCity = city
	b.ServiceState = state
	b.ServiceNotes = notes
	b.UpdatedAt = time.Now()
}

// CanRefund checks if booking can be refunded
func (b *DiasporaBooking) CanRefund() bool {
	return b.Status == BookingStatusPaid || b.Status == BookingStatusCancelled
}
