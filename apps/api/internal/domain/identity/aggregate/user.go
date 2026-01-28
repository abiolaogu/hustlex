package aggregate

import (
	"errors"
	"time"

	"hustlex/internal/domain/identity/event"
	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrUserNotActive     = errors.New("user account is not active")
	ErrUserAlreadyExists = errors.New("user with this phone already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
	ErrSkillAlreadyAdded = errors.New("skill already added to profile")
	ErrSkillNotFound     = errors.New("skill not found in profile")
	ErrCannotDeactivate  = errors.New("cannot deactivate user with active obligations")
)

// UserTier represents the credit tier of a user
type UserTier string

const (
	TierBronze   UserTier = "bronze"
	TierSilver   UserTier = "silver"
	TierGold     UserTier = "gold"
	TierPlatinum UserTier = "platinum"
)

func (t UserTier) String() string {
	return string(t)
}

func (t UserTier) IsValid() bool {
	switch t {
	case TierBronze, TierSilver, TierGold, TierPlatinum:
		return true
	}
	return false
}

// NextTier returns the next tier up, or the same if already at max
func (t UserTier) NextTier() UserTier {
	switch t {
	case TierBronze:
		return TierSilver
	case TierSilver:
		return TierGold
	case TierGold:
		return TierPlatinum
	default:
		return t
	}
}

// Proficiency represents skill proficiency level
type Proficiency string

const (
	ProficiencyBeginner     Proficiency = "beginner"
	ProficiencyIntermediate Proficiency = "intermediate"
	ProficiencyExpert       Proficiency = "expert"
)

// UserSkill represents a skill with proficiency
type UserSkill struct {
	SkillID       valueobject.SkillID
	SkillName     string
	Proficiency   Proficiency
	YearsExp      int
	IsVerified    bool
	PortfolioURLs []string
	AddedAt       time.Time
}

// User is the aggregate root for user identity
type User struct {
	sharedevent.AggregateRoot

	// Identity
	id           valueobject.UserID
	phone        valueobject.PhoneNumber
	email        valueobject.Email
	fullName     valueobject.FullName
	username     string

	// Profile
	profileImage string
	bio          string
	location     string
	state        string
	dateOfBirth  *time.Time
	gender       string

	// Status
	isVerified   bool
	isActive     bool
	tier         UserTier

	// Referral
	referralCode string
	referredBy   *valueobject.UserID

	// Skills
	skills       []UserSkill

	// Timestamps
	lastLoginAt  *time.Time
	createdAt    time.Time
	updatedAt    time.Time

	// Optimistic locking
	version      int64
}

// NewUser creates a new user aggregate
func NewUser(
	id valueobject.UserID,
	phone valueobject.PhoneNumber,
	fullName valueobject.FullName,
	referralCode string,
) (*User, error) {
	user := &User{
		id:           id,
		phone:        phone,
		fullName:     fullName,
		tier:         TierBronze,
		isActive:     true,
		isVerified:   false,
		referralCode: referralCode,
		skills:       make([]UserSkill, 0),
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
		version:      1,
	}

	user.RecordEvent(event.NewUserRegistered(
		id.String(),
		phone.String(),
		fullName.String(),
		"",
		referralCode,
		"",
	))

	return user, nil
}

// ReconstructUser reconstructs a user from persistence (no events emitted)
func ReconstructUser(
	id valueobject.UserID,
	phone valueobject.PhoneNumber,
	email valueobject.Email,
	fullName valueobject.FullName,
	username string,
	profileImage string,
	bio string,
	location string,
	state string,
	dateOfBirth *time.Time,
	gender string,
	isVerified bool,
	isActive bool,
	tier UserTier,
	referralCode string,
	referredBy *valueobject.UserID,
	skills []UserSkill,
	lastLoginAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
	version int64,
) *User {
	return &User{
		id:           id,
		phone:        phone,
		email:        email,
		fullName:     fullName,
		username:     username,
		profileImage: profileImage,
		bio:          bio,
		location:     location,
		state:        state,
		dateOfBirth:  dateOfBirth,
		gender:       gender,
		isVerified:   isVerified,
		isActive:     isActive,
		tier:         tier,
		referralCode: referralCode,
		referredBy:   referredBy,
		skills:       skills,
		lastLoginAt:  lastLoginAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
		version:      version,
	}
}

// Getters (read-only access to state)

func (u *User) ID() valueobject.UserID           { return u.id }
func (u *User) Phone() valueobject.PhoneNumber   { return u.phone }
func (u *User) Email() valueobject.Email         { return u.email }
func (u *User) FullName() valueobject.FullName   { return u.fullName }
func (u *User) Username() string                 { return u.username }
func (u *User) ProfileImage() string             { return u.profileImage }
func (u *User) Bio() string                      { return u.bio }
func (u *User) Location() string                 { return u.location }
func (u *User) State() string                    { return u.state }
func (u *User) DateOfBirth() *time.Time          { return u.dateOfBirth }
func (u *User) Gender() string                   { return u.gender }
func (u *User) IsVerified() bool                 { return u.isVerified }
func (u *User) IsActive() bool                   { return u.isActive }
func (u *User) Tier() UserTier                   { return u.tier }
func (u *User) ReferralCode() string             { return u.referralCode }
func (u *User) ReferredBy() *valueobject.UserID  { return u.referredBy }
func (u *User) Skills() []UserSkill              { return u.skills }
func (u *User) LastLoginAt() *time.Time          { return u.lastLoginAt }
func (u *User) CreatedAt() time.Time             { return u.createdAt }
func (u *User) UpdatedAt() time.Time             { return u.updatedAt }
func (u *User) Version() int64                   { return u.version }

// Business Methods

// SetEmail sets the user's email
func (u *User) SetEmail(email valueobject.Email) {
	u.email = email
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewUserProfileUpdated(u.id.String(), map[string]string{
		"email": email.String(),
	}))
}

// SetReferrer sets who referred this user
func (u *User) SetReferrer(referrerID valueobject.UserID) {
	u.referredBy = &referrerID
	u.updatedAt = time.Now().UTC()
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(
	username string,
	bio string,
	location string,
	state string,
	dateOfBirth *time.Time,
	gender string,
	profileImage string,
) error {
	if !u.isActive {
		return ErrUserNotActive
	}

	updatedFields := make(map[string]string)

	if username != "" && username != u.username {
		u.username = username
		updatedFields["username"] = username
	}
	if bio != "" && bio != u.bio {
		u.bio = bio
		updatedFields["bio"] = bio
	}
	if location != "" && location != u.location {
		u.location = location
		updatedFields["location"] = location
	}
	if state != "" && state != u.state {
		u.state = state
		updatedFields["state"] = state
	}
	if dateOfBirth != nil {
		u.dateOfBirth = dateOfBirth
		updatedFields["date_of_birth"] = dateOfBirth.Format(time.RFC3339)
	}
	if gender != "" && gender != u.gender {
		u.gender = gender
		updatedFields["gender"] = gender
	}
	if profileImage != "" && profileImage != u.profileImage {
		u.profileImage = profileImage
		updatedFields["profile_image"] = profileImage
	}

	if len(updatedFields) > 0 {
		u.updatedAt = time.Now().UTC()
		u.RecordEvent(event.NewUserProfileUpdated(u.id.String(), updatedFields))
	}

	return nil
}

// RecordLogin records a successful login
func (u *User) RecordLogin(ipAddress, userAgent string) {
	now := time.Now().UTC()
	u.lastLoginAt = &now
	u.updatedAt = now

	u.RecordEvent(event.NewUserLoggedIn(u.id.String(), u.phone.String(), ipAddress, userAgent))
}

// MarkVerified marks the user as verified for a specific type
func (u *User) MarkVerified(verificationType string) {
	u.isVerified = true
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewUserVerified(u.id.String(), verificationType))
}

// UpgradeTier upgrades the user's tier
func (u *User) UpgradeTier(newTier UserTier, reason string) error {
	if !newTier.IsValid() {
		return ErrInvalidUserData
	}

	if !u.isActive {
		return ErrUserNotActive
	}

	previousTier := u.tier
	u.tier = newTier
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewUserTierUpgraded(u.id.String(), previousTier.String(), newTier.String(), reason))

	return nil
}

// AddSkill adds a skill to the user's profile
func (u *User) AddSkill(skillID valueobject.SkillID, skillName string, proficiency Proficiency, yearsExp int) error {
	if !u.isActive {
		return ErrUserNotActive
	}

	// Check if skill already exists
	for _, s := range u.skills {
		if s.SkillID == skillID {
			return ErrSkillAlreadyAdded
		}
	}

	skill := UserSkill{
		SkillID:       skillID,
		SkillName:     skillName,
		Proficiency:   proficiency,
		YearsExp:      yearsExp,
		IsVerified:    false,
		PortfolioURLs: make([]string, 0),
		AddedAt:       time.Now().UTC(),
	}

	u.skills = append(u.skills, skill)
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewSkillAdded(u.id.String(), skillID.String(), skillName, string(proficiency)))

	return nil
}

// RemoveSkill removes a skill from the user's profile
func (u *User) RemoveSkill(skillID valueobject.SkillID) error {
	if !u.isActive {
		return ErrUserNotActive
	}

	found := false
	newSkills := make([]UserSkill, 0, len(u.skills))
	for _, s := range u.skills {
		if s.SkillID != skillID {
			newSkills = append(newSkills, s)
		} else {
			found = true
		}
	}

	if !found {
		return ErrSkillNotFound
	}

	u.skills = newSkills
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewSkillRemoved(u.id.String(), skillID.String()))

	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate(reason, deactivatedBy string) error {
	if !u.isActive {
		return nil // Already deactivated
	}

	u.isActive = false
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewUserDeactivated(u.id.String(), reason, deactivatedBy))

	return nil
}

// Reactivate reactivates the user account
func (u *User) Reactivate(reactivatedBy string) error {
	if u.isActive {
		return nil // Already active
	}

	u.isActive = true
	u.updatedAt = time.Now().UTC()

	u.RecordEvent(event.NewUserReactivated(u.id.String(), reactivatedBy))

	return nil
}

// HasSkill checks if user has a specific skill
func (u *User) HasSkill(skillID valueobject.SkillID) bool {
	for _, s := range u.skills {
		if s.SkillID == skillID {
			return true
		}
	}
	return false
}

// GetSkill returns a specific skill if user has it
func (u *User) GetSkill(skillID valueobject.SkillID) (*UserSkill, bool) {
	for _, s := range u.skills {
		if s.SkillID == skillID {
			return &s, true
		}
	}
	return nil, false
}
