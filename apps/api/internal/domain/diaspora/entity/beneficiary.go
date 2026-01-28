package entity

import (
	"time"

	"github.com/google/uuid"
)

// BeneficiaryType represents the type of beneficiary
type BeneficiaryType string

const (
	BeneficiaryTypeBank         BeneficiaryType = "bank"
	BeneficiaryTypeMobileWallet BeneficiaryType = "mobile_wallet"
	BeneficiaryTypeCash         BeneficiaryType = "cash_pickup"
)

// Relationship represents the relationship to the beneficiary
type Relationship string

const (
	RelationshipSelf     Relationship = "self"
	RelationshipFamily   Relationship = "family"
	RelationshipFriend   Relationship = "friend"
	RelationshipBusiness Relationship = "business"
	RelationshipOther    Relationship = "other"
)

// VerificationStatus represents beneficiary verification status
type VerificationStatus string

const (
	VerificationPending  VerificationStatus = "pending"
	VerificationVerified VerificationStatus = "verified"
	VerificationFailed   VerificationStatus = "failed"
	VerificationExpired  VerificationStatus = "expired"
)

// Beneficiary represents a remittance recipient
type Beneficiary struct {
	ID                uuid.UUID          `json:"id" db:"id"`
	UserID            uuid.UUID          `json:"user_id" db:"user_id"`
	LinkedWalletID    *uuid.UUID         `json:"linked_wallet_id,omitempty" db:"linked_wallet_id"`

	// Personal Information
	FirstName         string             `json:"first_name" db:"first_name"`
	LastName          string             `json:"last_name" db:"last_name"`
	MiddleName        string             `json:"middle_name,omitempty" db:"middle_name"`
	Email             string             `json:"email,omitempty" db:"email"`
	Phone             string             `json:"phone" db:"phone"`
	PhoneCountryCode  string             `json:"phone_country_code" db:"phone_country_code"`

	// Address
	Country           string             `json:"country" db:"country"`
	State             string             `json:"state,omitempty" db:"state"`
	City              string             `json:"city,omitempty" db:"city"`
	Address           string             `json:"address,omitempty" db:"address"`
	PostalCode        string             `json:"postal_code,omitempty" db:"postal_code"`

	// Payment Details
	Type              BeneficiaryType    `json:"type" db:"type"`
	Currency          string             `json:"currency" db:"currency"`
	BankCode          string             `json:"bank_code,omitempty" db:"bank_code"`
	BankName          string             `json:"bank_name,omitempty" db:"bank_name"`
	AccountNumber     string             `json:"account_number,omitempty" db:"account_number"`
	AccountName       string             `json:"account_name,omitempty" db:"account_name"`
	RoutingNumber     string             `json:"routing_number,omitempty" db:"routing_number"`
	SwiftCode         string             `json:"swift_code,omitempty" db:"swift_code"`
	IBAN              string             `json:"iban,omitempty" db:"iban"`
	MobileWalletProvider string          `json:"mobile_wallet_provider,omitempty" db:"mobile_wallet_provider"`
	MobileWalletNumber   string          `json:"mobile_wallet_number,omitempty" db:"mobile_wallet_number"`

	// Relationship & Verification
	Relationship      Relationship       `json:"relationship" db:"relationship"`
	VerificationStatus VerificationStatus `json:"verification_status" db:"verification_status"`
	VerifiedAt        *time.Time         `json:"verified_at,omitempty" db:"verified_at"`

	// Metadata
	Nickname          string             `json:"nickname,omitempty" db:"nickname"`
	Notes             string             `json:"notes,omitempty" db:"notes"`
	IsFavorite        bool               `json:"is_favorite" db:"is_favorite"`
	TransferCount     int                `json:"transfer_count" db:"transfer_count"`
	LastTransferAt    *time.Time         `json:"last_transfer_at,omitempty" db:"last_transfer_at"`

	// Audit
	IsActive          bool               `json:"is_active" db:"is_active"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time         `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewBeneficiary creates a new beneficiary
func NewBeneficiary(userID uuid.UUID, firstName, lastName, phone, country, currency string, beneficiaryType BeneficiaryType, relationship Relationship) *Beneficiary {
	now := time.Now()
	return &Beneficiary{
		ID:                 uuid.New(),
		UserID:             userID,
		FirstName:          firstName,
		LastName:           lastName,
		Phone:              phone,
		Country:            country,
		Currency:           currency,
		Type:               beneficiaryType,
		Relationship:       relationship,
		VerificationStatus: VerificationPending,
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// FullName returns the beneficiary's full name
func (b *Beneficiary) FullName() string {
	if b.MiddleName != "" {
		return b.FirstName + " " + b.MiddleName + " " + b.LastName
	}
	return b.FirstName + " " + b.LastName
}

// DisplayName returns the display name (nickname or full name)
func (b *Beneficiary) DisplayName() string {
	if b.Nickname != "" {
		return b.Nickname
	}
	return b.FullName()
}

// IsVerified checks if beneficiary is verified
func (b *Beneficiary) IsVerified() bool {
	return b.VerificationStatus == VerificationVerified
}

// CanReceiveTransfer checks if beneficiary can receive transfers
func (b *Beneficiary) CanReceiveTransfer() bool {
	return b.IsActive && b.IsVerified()
}

// SetBankDetails sets bank account details
func (b *Beneficiary) SetBankDetails(bankCode, bankName, accountNumber, accountName string) {
	b.Type = BeneficiaryTypeBank
	b.BankCode = bankCode
	b.BankName = bankName
	b.AccountNumber = accountNumber
	b.AccountName = accountName
	b.UpdatedAt = time.Now()
}

// SetMobileWalletDetails sets mobile wallet details
func (b *Beneficiary) SetMobileWalletDetails(provider, number string) {
	b.Type = BeneficiaryTypeMobileWallet
	b.MobileWalletProvider = provider
	b.MobileWalletNumber = number
	b.UpdatedAt = time.Now()
}

// Verify marks the beneficiary as verified
func (b *Beneficiary) Verify() {
	now := time.Now()
	b.VerificationStatus = VerificationVerified
	b.VerifiedAt = &now
	b.UpdatedAt = now
}

// IncrementTransferCount increments the transfer count
func (b *Beneficiary) IncrementTransferCount() {
	now := time.Now()
	b.TransferCount++
	b.LastTransferAt = &now
	b.UpdatedAt = now
}

// LinkWallet links a HustleX wallet to the beneficiary
func (b *Beneficiary) LinkWallet(walletID uuid.UUID) {
	b.LinkedWalletID = &walletID
	b.UpdatedAt = time.Now()
}
