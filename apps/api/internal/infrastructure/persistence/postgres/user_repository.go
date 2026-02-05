package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// UserRepository implements repository.UserRepository for PostgreSQL
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// Save persists a user aggregate
func (r *UserRepository) Save(ctx context.Context, user *aggregate.User) error {
	query := `
		INSERT INTO users (
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
		ON CONFLICT (id) DO UPDATE SET
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			username = EXCLUDED.username,
			full_name = EXCLUDED.full_name,
			profile_image = EXCLUDED.profile_image,
			bio = EXCLUDED.bio,
			location = EXCLUDED.location,
			state = EXCLUDED.state,
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			is_verified = EXCLUDED.is_verified,
			status = EXCLUDED.status,
			tier = EXCLUDED.tier,
			last_login_at = EXCLUDED.last_login_at,
			updated_at = EXCLUDED.updated_at,
			version = users.version + 1
		WHERE users.version = $20
	`

	status := "active"
	if !user.IsActive() {
		status = "inactive"
	}

	var referredByID *string
	if rb := user.ReferredBy(); rb != nil {
		id := rb.String()
		referredByID = &id
	}

	result, err := r.db.ExecContext(ctx, query,
		user.ID().String(),
		user.Phone().String(),
		nullString(user.Email().String()),
		nullString(user.Username()),
		user.FullName().String(),
		nullString(user.ProfileImage()),
		nullString(user.Bio()),
		nullString(user.Location()),
		nullString(user.State()),
		user.DateOfBirth(),
		nullString(user.Gender()),
		user.IsVerified(),
		status,
		user.Tier().String(),
		user.ReferralCode(),
		referredByID,
		user.LastLoginAt(),
		user.CreatedAt(),
		user.UpdatedAt(),
		user.Version(),
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("optimistic locking error: user version mismatch")
	}

	// Save skills separately
	if err := r.saveUserSkills(ctx, user); err != nil {
		return fmt.Errorf("failed to save user skills: %w", err)
	}

	return nil
}

// SaveWithEvents persists a user and publishes domain events
func (r *UserRepository) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
	return r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
		// Save user
		if err := r.Save(ctx, user); err != nil {
			return err
		}

		// TODO: Publish events to event bus
		// events := user.DomainEvents()
		// for _, event := range events {
		//     if err := eventBus.Publish(ctx, event); err != nil {
		//         return fmt.Errorf("failed to publish event: %w", err)
		//     }
		// }
		// user.ClearEvents()

		return nil
	})
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error) {
	query := `
		SELECT
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var (
		idStr, phoneStr, fullNameStr, tierStr, referralCode string
		emailStr, username, profileImage, bio               sql.NullString
		location, state, gender                             sql.NullString
		dateOfBirth, lastLoginAt                            sql.NullTime
		referredByStr                                       sql.NullString
		isVerified                                          bool
		status                                              string
		createdAt, updatedAt                                time.Time
		version                                             int64
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &phoneStr, &emailStr, &username, &fullNameStr,
		&profileImage, &bio, &location, &state,
		&dateOfBirth, &gender, &isVerified, &status,
		&tierStr, &referralCode, &referredByStr,
		&lastLoginAt, &createdAt, &updatedAt, &version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.reconstructUser(
		idStr, phoneStr, emailStr.String, fullNameStr, username.String,
		profileImage.String, bio.String, location.String, state.String,
		dateOfBirth, gender.String, isVerified, status == "active",
		tierStr, referralCode, referredByStr.String,
		lastLoginAt, createdAt, updatedAt, version,
	)
}

// FindByPhone retrieves a user by phone number
func (r *UserRepository) FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error) {
	query := `
		SELECT
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		FROM users
		WHERE phone = $1 AND deleted_at IS NULL
	`

	var (
		idStr, phoneStr, fullNameStr, tierStr, referralCode string
		emailStr, username, profileImage, bio               sql.NullString
		location, state, gender                             sql.NullString
		dateOfBirth, lastLoginAt                            sql.NullTime
		referredByStr                                       sql.NullString
		isVerified                                          bool
		status                                              string
		createdAt, updatedAt                                time.Time
		version                                             int64
	)

	err := r.db.QueryRowContext(ctx, query, phone.String()).Scan(
		&idStr, &phoneStr, &emailStr, &username, &fullNameStr,
		&profileImage, &bio, &location, &state,
		&dateOfBirth, &gender, &isVerified, &status,
		&tierStr, &referralCode, &referredByStr,
		&lastLoginAt, &createdAt, &updatedAt, &version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.reconstructUser(
		idStr, phoneStr, emailStr.String, fullNameStr, username.String,
		profileImage.String, bio.String, location.String, state.String,
		dateOfBirth, gender.String, isVerified, status == "active",
		tierStr, referralCode, referredByStr.String,
		lastLoginAt, createdAt, updatedAt, version,
	)
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*aggregate.User, error) {
	query := `
		SELECT
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var (
		idStr, phoneStr, fullNameStr, tierStr, referralCode string
		emailStr, username, profileImage, bio               sql.NullString
		location, state, gender                             sql.NullString
		dateOfBirth, lastLoginAt                            sql.NullTime
		referredByStr                                       sql.NullString
		isVerified                                          bool
		status                                              string
		createdAt, updatedAt                                time.Time
		version                                             int64
	)

	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(
		&idStr, &phoneStr, &emailStr, &username, &fullNameStr,
		&profileImage, &bio, &location, &state,
		&dateOfBirth, &gender, &isVerified, &status,
		&tierStr, &referralCode, &referredByStr,
		&lastLoginAt, &createdAt, &updatedAt, &version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.reconstructUser(
		idStr, phoneStr, emailStr.String, fullNameStr, username.String,
		profileImage.String, bio.String, location.String, state.String,
		dateOfBirth, gender.String, isVerified, status == "active",
		tierStr, referralCode, referredByStr.String,
		lastLoginAt, createdAt, updatedAt, version,
	)
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	query := `
		SELECT
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`

	var (
		idStr, phoneStr, fullNameStr, tierStr, referralCode string
		emailStr, usernameStr, profileImage, bio            sql.NullString
		location, state, gender                             sql.NullString
		dateOfBirth, lastLoginAt                            sql.NullTime
		referredByStr                                       sql.NullString
		isVerified                                          bool
		status                                              string
		createdAt, updatedAt                                time.Time
		version                                             int64
	)

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&idStr, &phoneStr, &emailStr, &usernameStr, &fullNameStr,
		&profileImage, &bio, &location, &state,
		&dateOfBirth, &gender, &isVerified, &status,
		&tierStr, &referralCode, &referredByStr,
		&lastLoginAt, &createdAt, &updatedAt, &version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.reconstructUser(
		idStr, phoneStr, emailStr.String, fullNameStr, usernameStr.String,
		profileImage.String, bio.String, location.String, state.String,
		dateOfBirth, gender.String, isVerified, status == "active",
		tierStr, referralCode, referredByStr.String,
		lastLoginAt, createdAt, updatedAt, version,
	)
}

// FindByReferralCode retrieves a user by their referral code
func (r *UserRepository) FindByReferralCode(ctx context.Context, code string) (*aggregate.User, error) {
	query := `
		SELECT
			id, phone, email, username, full_name,
			profile_image, bio, location, state,
			date_of_birth, gender, is_verified, status,
			tier, referral_code, referred_by,
			last_login_at, created_at, updated_at, version
		FROM users
		WHERE referral_code = $1 AND deleted_at IS NULL
	`

	var (
		idStr, phoneStr, fullNameStr, tierStr, referralCode string
		emailStr, username, profileImage, bio               sql.NullString
		location, state, gender                             sql.NullString
		dateOfBirth, lastLoginAt                            sql.NullTime
		referredByStr                                       sql.NullString
		isVerified                                          bool
		status                                              string
		createdAt, updatedAt                                time.Time
		version                                             int64
	)

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&idStr, &phoneStr, &emailStr, &username, &fullNameStr,
		&profileImage, &bio, &location, &state,
		&dateOfBirth, &gender, &isVerified, &status,
		&tierStr, &referralCode, &referredByStr,
		&lastLoginAt, &createdAt, &updatedAt, &version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.reconstructUser(
		idStr, phoneStr, emailStr.String, fullNameStr, username.String,
		profileImage.String, bio.String, location.String, state.String,
		dateOfBirth, gender.String, isVerified, status == "active",
		tierStr, referralCode, referredByStr.String,
		lastLoginAt, createdAt, updatedAt, version,
	)
}

// ExistsByPhone checks if a user with the phone exists
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, phone.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}

	return exists, nil
}

// ExistsByEmail checks if a user with the email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// ExistsByUsername checks if a user with the username exists
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check username existence: %w", err)
	}

	return exists, nil
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	query := `UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Helper methods

func (r *UserRepository) saveUserSkills(ctx context.Context, user *aggregate.User) error {
	// Delete existing skills
	deleteQuery := `DELETE FROM user_skills WHERE user_id = $1`
	if _, err := r.db.ExecContext(ctx, deleteQuery, user.ID().String()); err != nil {
		return fmt.Errorf("failed to delete existing skills: %w", err)
	}

	// Insert current skills
	if len(user.Skills()) == 0 {
		return nil
	}

	insertQuery := `
		INSERT INTO user_skills (
			user_id, skill_id, skill_name, proficiency,
			years_exp, is_verified, portfolio_urls, added_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, skill := range user.Skills() {
		portfolioJSON, _ := json.Marshal(skill.PortfolioURLs)

		_, err := r.db.ExecContext(ctx, insertQuery,
			user.ID().String(),
			skill.SkillID.String(),
			skill.SkillName,
			skill.Proficiency,
			skill.YearsExp,
			skill.IsVerified,
			portfolioJSON,
			skill.AddedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert skill: %w", err)
		}
	}

	return nil
}

func (r *UserRepository) loadUserSkills(ctx context.Context, userID string) ([]aggregate.UserSkill, error) {
	query := `
		SELECT skill_id, skill_name, proficiency, years_exp,
		       is_verified, portfolio_urls, added_at
		FROM user_skills
		WHERE user_id = $1
		ORDER BY added_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load skills: %w", err)
	}
	defer rows.Close()

	skills := make([]aggregate.UserSkill, 0)
	for rows.Next() {
		var (
			skillIDStr   string
			skillName    string
			proficiency  string
			yearsExp     int
			isVerified   bool
			portfolioStr string
			addedAt      time.Time
		)

		err := rows.Scan(&skillIDStr, &skillName, &proficiency, &yearsExp,
			&isVerified, &portfolioStr, &addedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		var portfolioURLs []string
		if portfolioStr != "" {
			if err := json.Unmarshal([]byte(portfolioStr), &portfolioURLs); err != nil {
				portfolioURLs = []string{}
			}
		}

		skillID, _ := valueobject.NewSkillID(skillIDStr)
		skills = append(skills, aggregate.UserSkill{
			SkillID:       skillID,
			SkillName:     skillName,
			Proficiency:   aggregate.Proficiency(proficiency),
			YearsExp:      yearsExp,
			IsVerified:    isVerified,
			PortfolioURLs: portfolioURLs,
			AddedAt:       addedAt,
		})
	}

	return skills, nil
}

func (r *UserRepository) reconstructUser(
	idStr, phoneStr, emailStr, fullNameStr, username,
	profileImage, bio, location, state string,
	dateOfBirth sql.NullTime, gender string, isVerified, isActive bool,
	tierStr, referralCode, referredByStr string,
	lastLoginAt sql.NullTime, createdAt, updatedAt time.Time, version int64,
) (*aggregate.User, error) {
	// Load skills
	skills, err := r.loadUserSkills(context.Background(), idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load skills: %w", err)
	}

	// Convert value objects
	id, _ := valueobject.NewUserID(idStr)
	phone, _ := valueobject.NewPhoneNumber(phoneStr)
	email, _ := valueobject.NewEmail(emailStr)
	fullName, _ := valueobject.NewFullName(fullNameStr)

	var dateOfBirthPtr *time.Time
	if dateOfBirth.Valid {
		dateOfBirthPtr = &dateOfBirth.Time
	}

	var lastLoginPtr *time.Time
	if lastLoginAt.Valid {
		lastLoginPtr = &lastLoginAt.Time
	}

	var referredBy *valueobject.UserID
	if referredByStr != "" {
		rb, _ := valueobject.NewUserID(referredByStr)
		referredBy = &rb
	}

	tier := aggregate.UserTier(tierStr)

	return aggregate.ReconstructUser(
		id, phone, email, fullName, username,
		profileImage, bio, location, state,
		dateOfBirthPtr, gender, isVerified, isActive,
		tier, referralCode, referredBy, skills,
		lastLoginPtr, createdAt, updatedAt, version,
	), nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
