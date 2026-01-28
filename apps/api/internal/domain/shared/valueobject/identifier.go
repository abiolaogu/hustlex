package valueobject

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

var (
	ErrInvalidID = errors.New("invalid identifier")
)

// UserID represents a unique user identifier
type UserID struct {
	value string
}

// NewUserID creates a new UserID from a string
func NewUserID(id string) (UserID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return UserID{}, ErrInvalidID
	}
	return UserID{value: id}, nil
}

// GenerateUserID creates a new random UserID
func GenerateUserID() UserID {
	return UserID{value: uuid.NewString()}
}

// MustNewUserID creates a UserID or panics
func MustNewUserID(id string) UserID {
	uid, err := NewUserID(id)
	if err != nil {
		panic(err)
	}
	return uid
}

// String returns the string representation
func (id UserID) String() string {
	return id.value
}

// IsEmpty returns true if the ID is empty
func (id UserID) IsEmpty() bool {
	return id.value == ""
}

// Equals checks equality with another UserID
func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

// WalletID represents a unique wallet identifier
type WalletID struct {
	value string
}

// NewWalletID creates a new WalletID from a string
func NewWalletID(id string) (WalletID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return WalletID{}, ErrInvalidID
	}
	return WalletID{value: id}, nil
}

// GenerateWalletID creates a new random WalletID
func GenerateWalletID() WalletID {
	return WalletID{value: uuid.NewString()}
}

func (id WalletID) String() string  { return id.value }
func (id WalletID) IsEmpty() bool   { return id.value == "" }
func (id WalletID) Equals(other WalletID) bool { return id.value == other.value }

// GigID represents a unique gig identifier
type GigID struct {
	value string
}

func NewGigID(id string) (GigID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return GigID{}, ErrInvalidID
	}
	return GigID{value: id}, nil
}

func GenerateGigID() GigID {
	return GigID{value: uuid.NewString()}
}

func (id GigID) String() string { return id.value }
func (id GigID) IsEmpty() bool  { return id.value == "" }
func (id GigID) Equals(other GigID) bool { return id.value == other.value }

// ContractID represents a unique contract identifier
type ContractID struct {
	value string
}

func NewContractID(id string) (ContractID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ContractID{}, ErrInvalidID
	}
	return ContractID{value: id}, nil
}

func GenerateContractID() ContractID {
	return ContractID{value: uuid.NewString()}
}

func (id ContractID) String() string { return id.value }
func (id ContractID) IsEmpty() bool  { return id.value == "" }
func (id ContractID) Equals(other ContractID) bool { return id.value == other.value }

// CircleID represents a unique savings circle identifier
type CircleID struct {
	value string
}

func NewCircleID(id string) (CircleID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return CircleID{}, ErrInvalidID
	}
	return CircleID{value: id}, nil
}

func GenerateCircleID() CircleID {
	return CircleID{value: uuid.NewString()}
}

func (id CircleID) String() string { return id.value }
func (id CircleID) IsEmpty() bool  { return id.value == "" }
func (id CircleID) Equals(other CircleID) bool { return id.value == other.value }

// LoanID represents a unique loan identifier
type LoanID struct {
	value string
}

func NewLoanID(id string) (LoanID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return LoanID{}, ErrInvalidID
	}
	return LoanID{value: id}, nil
}

func GenerateLoanID() LoanID {
	return LoanID{value: uuid.NewString()}
}

func (id LoanID) String() string { return id.value }
func (id LoanID) IsEmpty() bool  { return id.value == "" }
func (id LoanID) Equals(other LoanID) bool { return id.value == other.value }

// TransactionID represents a unique transaction identifier
type TransactionID struct {
	value string
}

func NewTransactionID(id string) (TransactionID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return TransactionID{}, ErrInvalidID
	}
	return TransactionID{value: id}, nil
}

func GenerateTransactionID() TransactionID {
	return TransactionID{value: uuid.NewString()}
}

func (id TransactionID) String() string { return id.value }
func (id TransactionID) IsEmpty() bool  { return id.value == "" }
func (id TransactionID) Equals(other TransactionID) bool { return id.value == other.value }

// Reference represents a business reference (e.g., payment reference)
type Reference struct {
	value string
}

var referencePattern = regexp.MustCompile(`^[A-Z]{2,5}[0-9]{14,20}$`)

func NewReference(value string) (Reference, error) {
	if value == "" {
		return Reference{}, errors.New("reference cannot be empty")
	}
	return Reference{value: value}, nil
}

func (r Reference) String() string { return r.value }
func (r Reference) IsEmpty() bool  { return r.value == "" }
func (r Reference) Equals(other Reference) bool { return r.value == other.value }

// SkillID represents a unique skill identifier
type SkillID struct {
	value string
}

func NewSkillID(id string) (SkillID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return SkillID{}, ErrInvalidID
	}
	return SkillID{value: id}, nil
}

func GenerateSkillID() SkillID {
	return SkillID{value: uuid.NewString()}
}

func MustNewSkillID(id string) SkillID {
	sid, err := NewSkillID(id)
	if err != nil {
		panic(err)
	}
	return sid
}

func (id SkillID) String() string { return id.value }
func (id SkillID) IsEmpty() bool  { return id.value == "" }
func (id SkillID) Equals(other SkillID) bool { return id.value == other.value }

// ProposalID represents a unique gig proposal identifier
type ProposalID struct {
	value string
}

func NewProposalID(id string) (ProposalID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ProposalID{}, ErrInvalidID
	}
	return ProposalID{value: id}, nil
}

func GenerateProposalID() ProposalID {
	return ProposalID{value: uuid.NewString()}
}

func (id ProposalID) String() string { return id.value }
func (id ProposalID) IsEmpty() bool  { return id.value == "" }
func (id ProposalID) Equals(other ProposalID) bool { return id.value == other.value }

// MemberID represents a unique circle member identifier
type MemberID struct {
	value string
}

func NewMemberID(id string) (MemberID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return MemberID{}, ErrInvalidID
	}
	return MemberID{value: id}, nil
}

func GenerateMemberID() MemberID {
	return MemberID{value: uuid.NewString()}
}

func (id MemberID) String() string { return id.value }
func (id MemberID) IsEmpty() bool  { return id.value == "" }
func (id MemberID) Equals(other MemberID) bool { return id.value == other.value }

// ContributionID represents a unique contribution identifier
type ContributionID struct {
	value string
}

func NewContributionID(id string) (ContributionID, error) {
	if _, err := uuid.Parse(id); err != nil {
		return ContributionID{}, ErrInvalidID
	}
	return ContributionID{value: id}, nil
}

func GenerateContributionID() ContributionID {
	return ContributionID{value: uuid.NewString()}
}

func (id ContributionID) String() string { return id.value }
func (id ContributionID) IsEmpty() bool  { return id.value == "" }
func (id ContributionID) Equals(other ContributionID) bool { return id.value == other.value }
