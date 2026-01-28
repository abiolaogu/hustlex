package valueobject

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUserID(t *testing.T) {
	validUUID := uuid.NewString()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			id:      validUUID,
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			id:      "invalid-uuid",
			wantErr: true,
		},
		{
			name:    "empty string",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := NewUserID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Error("NewUserID() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewUserID() unexpected error: %v", err)
				return
			}
			if userID.String() != tt.id {
				t.Errorf("String() = %s, want %s", userID.String(), tt.id)
			}
		})
	}
}

func TestGenerateUserID(t *testing.T) {
	id1 := GenerateUserID()
	id2 := GenerateUserID()

	if id1.IsEmpty() {
		t.Error("GenerateUserID() should not return empty")
	}
	if id1.Equals(id2) {
		t.Error("GenerateUserID() should generate unique IDs")
	}

	// Verify it's a valid UUID
	if _, err := uuid.Parse(id1.String()); err != nil {
		t.Errorf("GenerateUserID() should generate valid UUID, got %s", id1.String())
	}
}

func TestUserID_Equals(t *testing.T) {
	validUUID := uuid.NewString()
	id1, _ := NewUserID(validUUID)
	id2, _ := NewUserID(validUUID)
	id3 := GenerateUserID()

	if !id1.Equals(id2) {
		t.Error("same UUIDs should be equal")
	}
	if id1.Equals(id3) {
		t.Error("different UUIDs should not be equal")
	}
}

func TestUserID_IsEmpty(t *testing.T) {
	var emptyID UserID
	validID := GenerateUserID()

	if !emptyID.IsEmpty() {
		t.Error("zero-value UserID should be empty")
	}
	if validID.IsEmpty() {
		t.Error("generated UserID should not be empty")
	}
}

func TestMustNewUserID(t *testing.T) {
	validUUID := uuid.NewString()

	t.Run("valid UUID", func(t *testing.T) {
		id := MustNewUserID(validUUID)
		if id.String() != validUUID {
			t.Errorf("MustNewUserID() = %s, want %s", id.String(), validUUID)
		}
	})

	t.Run("invalid UUID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustNewUserID should panic on invalid input")
			}
		}()
		MustNewUserID("invalid")
	})
}

func TestNewWalletID(t *testing.T) {
	validUUID := uuid.NewString()

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid UUID",
			id:      validUUID,
			wantErr: false,
		},
		{
			name:    "invalid UUID",
			id:      "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			walletID, err := NewWalletID(tt.id)
			if tt.wantErr {
				if err == nil {
					t.Error("NewWalletID() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewWalletID() unexpected error: %v", err)
				return
			}
			if walletID.String() != tt.id {
				t.Errorf("String() = %s, want %s", walletID.String(), tt.id)
			}
		})
	}
}

func TestGenerateWalletID(t *testing.T) {
	id := GenerateWalletID()

	if id.IsEmpty() {
		t.Error("GenerateWalletID() should not return empty")
	}

	if _, err := uuid.Parse(id.String()); err != nil {
		t.Errorf("GenerateWalletID() should generate valid UUID, got %s", id.String())
	}
}

func TestNewGigID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewGigID(validUUID)
	if err != nil {
		t.Fatalf("NewGigID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}

	_, err = NewGigID("invalid")
	if err != ErrInvalidID {
		t.Errorf("NewGigID() with invalid input should return ErrInvalidID, got %v", err)
	}
}

func TestNewContractID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewContractID(validUUID)
	if err != nil {
		t.Fatalf("NewContractID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestNewCircleID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewCircleID(validUUID)
	if err != nil {
		t.Fatalf("NewCircleID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestNewLoanID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewLoanID(validUUID)
	if err != nil {
		t.Fatalf("NewLoanID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestNewTransactionID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewTransactionID(validUUID)
	if err != nil {
		t.Fatalf("NewTransactionID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestNewReference(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid reference",
			value:   "TRF20240101120000123456",
			wantErr: false,
		},
		{
			name:    "another valid reference",
			value:   "DEP-12345-67890",
			wantErr: false,
		},
		{
			name:    "empty reference",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewReference(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Error("NewReference() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewReference() unexpected error: %v", err)
				return
			}
			if ref.String() != tt.value {
				t.Errorf("String() = %s, want %s", ref.String(), tt.value)
			}
		})
	}
}

func TestReference_Equals(t *testing.T) {
	ref1, _ := NewReference("REF123")
	ref2, _ := NewReference("REF123")
	ref3, _ := NewReference("REF456")

	if !ref1.Equals(ref2) {
		t.Error("same references should be equal")
	}
	if ref1.Equals(ref3) {
		t.Error("different references should not be equal")
	}
}

func TestNewProposalID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewProposalID(validUUID)
	if err != nil {
		t.Fatalf("NewProposalID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}

	_, err = NewProposalID("invalid")
	if err != ErrInvalidID {
		t.Errorf("NewProposalID() with invalid input should return ErrInvalidID")
	}
}

func TestNewSkillID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewSkillID(validUUID)
	if err != nil {
		t.Fatalf("NewSkillID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestMustNewSkillID(t *testing.T) {
	validUUID := uuid.NewString()

	t.Run("valid UUID", func(t *testing.T) {
		id := MustNewSkillID(validUUID)
		if id.String() != validUUID {
			t.Errorf("MustNewSkillID() = %s, want %s", id.String(), validUUID)
		}
	})

	t.Run("invalid UUID panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustNewSkillID should panic on invalid input")
			}
		}()
		MustNewSkillID("invalid")
	})
}

func TestNewMemberID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewMemberID(validUUID)
	if err != nil {
		t.Fatalf("NewMemberID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

func TestNewContributionID(t *testing.T) {
	validUUID := uuid.NewString()

	id, err := NewContributionID(validUUID)
	if err != nil {
		t.Fatalf("NewContributionID() unexpected error: %v", err)
	}
	if id.String() != validUUID {
		t.Errorf("String() = %s, want %s", id.String(), validUUID)
	}
}

// Test all Generate* functions create unique IDs
func TestGenerateFunctions_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)

	// Generate multiple IDs of each type and verify uniqueness
	for i := 0; i < 100; i++ {
		userID := GenerateUserID().String()
		if ids[userID] {
			t.Errorf("GenerateUserID() produced duplicate: %s", userID)
		}
		ids[userID] = true

		walletID := GenerateWalletID().String()
		if ids[walletID] {
			t.Errorf("GenerateWalletID() produced duplicate: %s", walletID)
		}
		ids[walletID] = true

		gigID := GenerateGigID().String()
		if ids[gigID] {
			t.Errorf("GenerateGigID() produced duplicate: %s", gigID)
		}
		ids[gigID] = true
	}
}
