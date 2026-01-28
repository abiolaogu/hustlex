package repository

import (
	"context"
	"time"

	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// UserRepository defines the interface for user persistence
// This is a PORT - infrastructure provides the ADAPTER
type UserRepository interface {
	// Save persists a user aggregate
	Save(ctx context.Context, user *aggregate.User) error

	// SaveWithEvents persists a user and publishes domain events
	SaveWithEvents(ctx context.Context, user *aggregate.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error)

	// FindByPhone retrieves a user by phone number
	FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email valueobject.Email) (*aggregate.User, error)

	// FindByUsername retrieves a user by username
	FindByUsername(ctx context.Context, username string) (*aggregate.User, error)

	// FindByReferralCode retrieves a user by their referral code
	FindByReferralCode(ctx context.Context, code string) (*aggregate.User, error)

	// ExistsByPhone checks if a user with the phone exists
	ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error)

	// ExistsByEmail checks if a user with the email exists
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)

	// ExistsByUsername checks if a user with the username exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// Delete soft-deletes a user
	Delete(ctx context.Context, id valueobject.UserID) error
}

// OTPRepository defines the interface for OTP persistence
type OTPRepository interface {
	// Save persists an OTP code
	Save(ctx context.Context, otp *OTPCode) error

	// FindLatestValid finds the latest valid (unused, unexpired) OTP
	FindLatestValid(ctx context.Context, phone string, purpose string) (*OTPCode, error)

	// MarkUsed marks an OTP as used
	MarkUsed(ctx context.Context, id string) error

	// IncrementAttempts increments the failed attempt counter
	IncrementAttempts(ctx context.Context, id string) error

	// DeleteExpired removes expired OTPs (for cleanup)
	DeleteExpired(ctx context.Context) error

	// DeleteUnused deletes unused OTPs for a phone/purpose combination
	DeleteUnused(ctx context.Context, phone, purpose string) error
}

// OTPCode represents an OTP for authentication
type OTPCode struct {
	ID        string
	Phone     string
	Code      string // Should be hashed in production
	Purpose   string // login, register, reset_pin
	ExpiresAt time.Time
	IsUsed    bool
	Attempts  int
	CreatedAt time.Time
}

// SkillRepository defines the interface for skill catalog persistence
type SkillRepository interface {
	// Save persists a skill
	Save(ctx context.Context, skill *Skill) error

	// FindByID retrieves a skill by ID
	FindByID(ctx context.Context, id valueobject.SkillID) (*Skill, error)

	// FindByCategory retrieves skills by category
	FindByCategory(ctx context.Context, category string) ([]*Skill, error)

	// FindAll retrieves all active skills
	FindAll(ctx context.Context) ([]*Skill, error)

	// Search searches skills by name
	Search(ctx context.Context, query string, limit int) ([]*Skill, error)
}

// Skill represents a skill in the catalog
type Skill struct {
	ID          string
	Name        string
	Category    string
	Description string
	Icon        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserSkillRepository defines the interface for user skill persistence
type UserSkillRepository interface {
	// Save persists a user skill
	Save(ctx context.Context, userSkill *UserSkillRecord) error

	// Delete removes a user skill
	Delete(ctx context.Context, userID valueobject.UserID, skillID valueobject.SkillID) error

	// FindByUserID retrieves all skills for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*UserSkillRecord, error)

	// FindBySkillID retrieves all users with a specific skill
	FindBySkillID(ctx context.Context, skillID valueobject.SkillID, limit, offset int) ([]*UserSkillRecord, error)

	// CountBySkillID counts users with a specific skill
	CountBySkillID(ctx context.Context, skillID valueobject.SkillID) (int64, error)
}

// UserSkillRecord represents a persisted user skill
type UserSkillRecord struct {
	ID            string
	UserID        string
	SkillID       string
	SkillName     string
	Proficiency   string
	YearsExp      int
	IsVerified    bool
	PortfolioURLs []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// SessionRepository defines the interface for session management
type SessionRepository interface {
	// StoreRefreshToken stores a refresh token
	StoreRefreshToken(ctx context.Context, userID valueobject.UserID, token string, expiry time.Duration) error

	// GetRefreshToken retrieves a stored refresh token
	GetRefreshToken(ctx context.Context, userID valueobject.UserID) (string, error)

	// DeleteRefreshToken removes a refresh token (logout)
	DeleteRefreshToken(ctx context.Context, userID valueobject.UserID) error

	// CheckOTPRateLimit checks and increments OTP rate limit
	CheckOTPRateLimit(ctx context.Context, phone string, maxRequests int, window time.Duration) (bool, error)
}

// ReferralRepository defines the interface for referral tracking
type ReferralRepository interface {
	// RecordReferral records that one user referred another
	RecordReferral(ctx context.Context, referrerID, referredID valueobject.UserID) error

	// GetReferralCount gets the count of successful referrals for a user
	GetReferralCount(ctx context.Context, userID valueobject.UserID) (int64, error)

	// GetReferrals gets the list of users referred by a user
	GetReferrals(ctx context.Context, userID valueobject.UserID, limit, offset int) ([]valueobject.UserID, error)
}
