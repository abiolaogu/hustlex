package command

import (
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

// SendOTP sends an OTP to a phone number
type SendOTP struct {
	Phone   string
	Purpose string // login, register, reset_pin
}

// VerifyOTP verifies an OTP code
type VerifyOTP struct {
	Phone     string
	Code      string
	Purpose   string
	IPAddress string
	UserAgent string
}

// RegisterUser registers a new user
type RegisterUser struct {
	Phone        string
	FullName     string
	Email        string
	ReferralCode string
}

// RegisterResult is the result of user registration
type RegisterResult struct {
	UserID       string    `json:"user_id"`
	Phone        string    `json:"phone"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email,omitempty"`
	Tier         string    `json:"tier"`
	ReferralCode string    `json:"referral_code"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// LoginResult is the result of a login
type LoginResult struct {
	UserID       string    `json:"user_id"`
	Phone        string    `json:"phone"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email,omitempty"`
	Tier         string    `json:"tier"`
	IsVerified   bool      `json:"is_verified"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshToken refreshes authentication tokens
type RefreshToken struct {
	RefreshToken string
}

// Logout logs out a user
type Logout struct {
	UserID string
}

// UpdateProfile updates user profile information
type UpdateProfile struct {
	UserID       string
	Username     string
	Bio          string
	Location     string
	State        string
	DateOfBirth  *time.Time
	Gender       string
	ProfileImage string
}

// UpdateProfileResult is the result of profile update
type UpdateProfileResult struct {
	UserID       string    `json:"user_id"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email,omitempty"`
	FullName     string    `json:"full_name"`
	Username     string    `json:"username,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Location     string    `json:"location,omitempty"`
	State        string    `json:"state,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AddUserSkill adds a skill to user's profile
type AddUserSkill struct {
	UserID      string
	SkillID     string
	Proficiency string // beginner, intermediate, expert
	YearsExp    int
}

// AddSkillResult is the result of adding a skill
type AddSkillResult struct {
	UserID    string `json:"user_id"`
	SkillID   string `json:"skill_id"`
	SkillName string `json:"skill_name"`
}

// RemoveUserSkill removes a skill from user's profile
type RemoveUserSkill struct {
	UserID  string
	SkillID string
}

// DeactivateUser deactivates a user account
type DeactivateUser struct {
	UserID      string
	Reason      string
	RequestedBy string
}

// ReactivateUser reactivates a user account
type ReactivateUser struct {
	UserID      string
	RequestedBy string
}

// UpgradeUserTier upgrades a user's credit tier
type UpgradeUserTier struct {
	UserID      string
	NewTier     string
	Reason      string
	RequestedBy string
}

// Helper methods

// GetPhone returns the phone as a value object
func (c RegisterUser) GetPhone() (valueobject.PhoneNumber, error) {
	return valueobject.NewPhoneNumber(c.Phone)
}

// GetFullName returns the full name as a value object
func (c RegisterUser) GetFullName() (valueobject.FullName, error) {
	return valueobject.NewFullName(c.FullName)
}

// GetEmail returns the email as a value object
func (c RegisterUser) GetEmail() (valueobject.Email, error) {
	if c.Email == "" {
		return valueobject.Email{}, nil
	}
	return valueobject.NewEmail(c.Email)
}

// GetPhone returns the phone as a value object
func (c VerifyOTP) GetPhone() (valueobject.PhoneNumber, error) {
	return valueobject.NewPhoneNumber(c.Phone)
}

// GetPhone returns the phone as a value object
func (c SendOTP) GetPhone() (valueobject.PhoneNumber, error) {
	return valueobject.NewPhoneNumber(c.Phone)
}

// GetUserID returns the user ID as a value object
func (c UpdateProfile) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// GetUserID returns the user ID as a value object
func (c AddUserSkill) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// GetSkillID returns the skill ID as a value object
func (c AddUserSkill) GetSkillID() (valueobject.SkillID, error) {
	return valueobject.NewSkillID(c.SkillID)
}
