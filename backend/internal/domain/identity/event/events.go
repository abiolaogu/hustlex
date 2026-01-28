package event

import (
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
)

const (
	AggregateTypeUser = "User"
)

// UserRegistered is emitted when a new user registers
type UserRegistered struct {
	sharedevent.BaseEvent
	UserID       string `json:"user_id"`
	Phone        string `json:"phone"`
	FullName     string `json:"full_name"`
	Email        string `json:"email,omitempty"`
	ReferralCode string `json:"referral_code"`
	ReferredBy   string `json:"referred_by,omitempty"`
}

func NewUserRegistered(userID, phone, fullName, email, referralCode, referredBy string) *UserRegistered {
	return &UserRegistered{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserRegistered",
			userID,
			AggregateTypeUser,
		),
		UserID:       userID,
		Phone:        phone,
		FullName:     fullName,
		Email:        email,
		ReferralCode: referralCode,
		ReferredBy:   referredBy,
	}
}

// UserLoggedIn is emitted when a user successfully logs in
type UserLoggedIn struct {
	sharedevent.BaseEvent
	UserID    string    `json:"user_id"`
	Phone     string    `json:"phone"`
	LoginTime time.Time `json:"login_time"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

func NewUserLoggedIn(userID, phone, ipAddress, userAgent string) *UserLoggedIn {
	return &UserLoggedIn{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserLoggedIn",
			userID,
			AggregateTypeUser,
		),
		UserID:    userID,
		Phone:     phone,
		LoginTime: time.Now().UTC(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}

// UserProfileUpdated is emitted when user updates their profile
type UserProfileUpdated struct {
	sharedevent.BaseEvent
	UserID        string            `json:"user_id"`
	UpdatedFields map[string]string `json:"updated_fields"`
}

func NewUserProfileUpdated(userID string, updatedFields map[string]string) *UserProfileUpdated {
	return &UserProfileUpdated{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserProfileUpdated",
			userID,
			AggregateTypeUser,
		),
		UserID:        userID,
		UpdatedFields: updatedFields,
	}
}

// UserVerified is emitted when user completes verification (KYC)
type UserVerified struct {
	sharedevent.BaseEvent
	UserID           string `json:"user_id"`
	VerificationType string `json:"verification_type"` // phone, email, bvn, nin
	VerifiedAt       time.Time `json:"verified_at"`
}

func NewUserVerified(userID, verificationType string) *UserVerified {
	return &UserVerified{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserVerified",
			userID,
			AggregateTypeUser,
		),
		UserID:           userID,
		VerificationType: verificationType,
		VerifiedAt:       time.Now().UTC(),
	}
}

// UserTierUpgraded is emitted when user's credit tier is upgraded
type UserTierUpgraded struct {
	sharedevent.BaseEvent
	UserID      string `json:"user_id"`
	PreviousTier string `json:"previous_tier"`
	NewTier      string `json:"new_tier"`
	Reason       string `json:"reason"`
}

func NewUserTierUpgraded(userID, previousTier, newTier, reason string) *UserTierUpgraded {
	return &UserTierUpgraded{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserTierUpgraded",
			userID,
			AggregateTypeUser,
		),
		UserID:       userID,
		PreviousTier: previousTier,
		NewTier:      newTier,
		Reason:       reason,
	}
}

// UserDeactivated is emitted when user account is deactivated
type UserDeactivated struct {
	sharedevent.BaseEvent
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
	By     string `json:"by"` // user_id who deactivated (self or admin)
}

func NewUserDeactivated(userID, reason, by string) *UserDeactivated {
	return &UserDeactivated{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserDeactivated",
			userID,
			AggregateTypeUser,
		),
		UserID: userID,
		Reason: reason,
		By:     by,
	}
}

// UserReactivated is emitted when user account is reactivated
type UserReactivated struct {
	sharedevent.BaseEvent
	UserID string `json:"user_id"`
	By     string `json:"by"`
}

func NewUserReactivated(userID, by string) *UserReactivated {
	return &UserReactivated{
		BaseEvent: sharedevent.NewBaseEvent(
			"UserReactivated",
			userID,
			AggregateTypeUser,
		),
		UserID: userID,
		By:     by,
	}
}

// OTPGenerated is emitted when an OTP is generated
type OTPGenerated struct {
	sharedevent.BaseEvent
	Phone     string    `json:"phone"`
	Purpose   string    `json:"purpose"` // login, register, reset_pin
	ExpiresAt time.Time `json:"expires_at"`
}

func NewOTPGenerated(phone, purpose string, expiresAt time.Time) *OTPGenerated {
	return &OTPGenerated{
		BaseEvent: sharedevent.NewBaseEvent(
			"OTPGenerated",
			phone,
			"OTP",
		),
		Phone:     phone,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
	}
}

// OTPVerified is emitted when an OTP is successfully verified
type OTPVerified struct {
	sharedevent.BaseEvent
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
}

func NewOTPVerified(phone, purpose string) *OTPVerified {
	return &OTPVerified{
		BaseEvent: sharedevent.NewBaseEvent(
			"OTPVerified",
			phone,
			"OTP",
		),
		Phone:   phone,
		Purpose: purpose,
	}
}

// SkillAdded is emitted when a user adds a skill to their profile
type SkillAdded struct {
	sharedevent.BaseEvent
	UserID      string `json:"user_id"`
	SkillID     string `json:"skill_id"`
	SkillName   string `json:"skill_name"`
	Proficiency string `json:"proficiency"`
}

func NewSkillAdded(userID, skillID, skillName, proficiency string) *SkillAdded {
	return &SkillAdded{
		BaseEvent: sharedevent.NewBaseEvent(
			"SkillAdded",
			userID,
			AggregateTypeUser,
		),
		UserID:      userID,
		SkillID:     skillID,
		SkillName:   skillName,
		Proficiency: proficiency,
	}
}

// SkillRemoved is emitted when a user removes a skill
type SkillRemoved struct {
	sharedevent.BaseEvent
	UserID  string `json:"user_id"`
	SkillID string `json:"skill_id"`
}

func NewSkillRemoved(userID, skillID string) *SkillRemoved {
	return &SkillRemoved{
		BaseEvent: sharedevent.NewBaseEvent(
			"SkillRemoved",
			userID,
			AggregateTypeUser,
		),
		UserID:  userID,
		SkillID: skillID,
	}
}
