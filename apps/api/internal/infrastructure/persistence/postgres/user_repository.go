package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"

	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// UserRepository is a PostgreSQL implementation of the UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &UserRepository{db: db}
}

// userRow represents the database row structure
type userRow struct {
	ID              string
	Phone           string
	Email           sql.NullString
	FirstName       string
	LastName        string
	Username        sql.NullString
	ProfileImage    sql.NullString
	Bio             sql.NullString
	Location        sql.NullString
	State           sql.NullString
	DateOfBirth     sql.NullTime
	Gender          sql.NullString
	IsVerified      bool
	IsActive        bool
	Tier            string
	ReferralCode    string
	ReferredBy      sql.NullString
	LastLoginAt     sql.NullTime
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int64
}

// Save persists a user aggregate
func (r *UserRepository) Save(ctx context.Context, user *aggregate.User) error {
	query := `
		INSERT INTO users (
			id, phone, email, status, role, phone_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (id) DO UPDATE SET
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			phone_verified = EXCLUDED.phone_verified,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	email := sql.NullString{Valid: false}
	if user.Email().String() != "" {
		email = sql.NullString{String: user.Email().String(), Valid: true}
	}

	status := "active"
	if !user.IsActive() {
		status = "inactive"
	} else if !user.IsVerified() {
		status = "pending_verification"
	}

	var userID string
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.ID().String(),
		user.Phone().String(),
		email,
		status,
		"consumer", // default role
		user.IsVerified(),
		user.CreatedAt(),
		user.UpdatedAt(),
	).Scan(&userID)

	if err != nil {
		return err
	}

	// Save profile data
	profileQuery := `
		INSERT INTO profiles (
			user_id, first_name, last_name, display_name, avatar_url, bio,
			city, state, date_of_birth, gender, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		ON CONFLICT (user_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			bio = EXCLUDED.bio,
			city = EXCLUDED.city,
			state = EXCLUDED.state,
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			updated_at = EXCLUDED.updated_at
	`

	firstName, lastName := user.FullName().FirstName(), user.FullName().LastName()

	_, err = r.db.ExecContext(
		ctx,
		profileQuery,
		user.ID().String(),
		firstName,
		lastName,
		sql.NullString{String: user.Username(), Valid: user.Username() != ""},
		sql.NullString{String: user.ProfileImage(), Valid: user.ProfileImage() != ""},
		sql.NullString{String: user.Bio(), Valid: user.Bio() != ""},
		sql.NullString{String: user.Location(), Valid: user.Location() != ""},
		sql.NullString{String: user.State(), Valid: user.State() != ""},
		sql.NullTime{Time: *user.DateOfBirth(), Valid: user.DateOfBirth() != nil},
		sql.NullString{String: user.Gender(), Valid: user.Gender() != ""},
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	if err != nil {
		return err
	}

	// Save tier information
	tierQuery := `
		INSERT INTO user_tiers (
			user_id, tier, updated_at
		) VALUES (
			$1, $2, $3
		)
		ON CONFLICT (user_id) DO UPDATE SET
			tier = EXCLUDED.tier,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.ExecContext(
		ctx,
		tierQuery,
		user.ID().String(),
		user.Tier().String(),
		user.UpdatedAt(),
	)

	if err != nil {
		return err
	}

	// Save referral info
	if user.ReferredBy() != nil {
		referralQuery := `
			INSERT INTO referrals (
				referrer_id, referred_id, created_at
			) VALUES (
				$1, $2, $3
			)
			ON CONFLICT (referred_id) DO NOTHING
		`

		_, err = r.db.ExecContext(
			ctx,
			referralQuery,
			user.ReferredBy().String(),
			user.ID().String(),
			user.CreatedAt(),
		)

		if err != nil {
			return err
		}
	}

	// Save skills
	if len(user.Skills()) > 0 {
		// First, delete existing skills
		_, err = r.db.ExecContext(
			ctx,
			"DELETE FROM user_skills WHERE user_id = $1",
			user.ID().String(),
		)
		if err != nil {
			return err
		}

		// Then insert all skills
		for _, skill := range user.Skills() {
			skillQuery := `
				INSERT INTO user_skills (
					user_id, skill_id, proficiency, years_experience,
					is_verified, portfolio_urls, created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8
				)
			`

			portfolioJSON, _ := json.Marshal(skill.PortfolioURLs)

			_, err = r.db.ExecContext(
				ctx,
				skillQuery,
				user.ID().String(),
				skill.SkillID.String(),
				string(skill.Proficiency),
				skill.YearsExp,
				skill.IsVerified,
				portfolioJSON,
				skill.AddedAt,
				time.Now().UTC(),
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SaveWithEvents persists a user and publishes domain events
func (r *UserRepository) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Save the user using a transaction-aware version
	err = r.saveInTx(ctx, tx, user)
	if err != nil {
		return err
	}

	// TODO: Publish domain events to event bus
	// This would integrate with the messaging infrastructure
	// For now, we'll just clear the events after save
	user.ClearEvents()

	// Commit transaction
	return tx.Commit()
}

// saveInTx is a helper that saves within an existing transaction
func (r *UserRepository) saveInTx(ctx context.Context, tx *sql.Tx, user *aggregate.User) error {
	// Similar to Save() but using tx instead of db
	// This is a simplified version - in production, you'd refactor Save() to use this
	query := `
		INSERT INTO users (
			id, phone, email, status, role, phone_verified, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (id) DO UPDATE SET
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			phone_verified = EXCLUDED.phone_verified,
			updated_at = EXCLUDED.updated_at
	`

	email := sql.NullString{Valid: false}
	if user.Email().String() != "" {
		email = sql.NullString{String: user.Email().String(), Valid: true}
	}

	status := "active"
	if !user.IsActive() {
		status = "inactive"
	} else if !user.IsVerified() {
		status = "pending_verification"
	}

	_, err := tx.ExecContext(
		ctx,
		query,
		user.ID().String(),
		user.Phone().String(),
		email,
		status,
		"consumer",
		user.IsVerified(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	return err
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error) {
	query := `
		SELECT
			u.id, u.phone, u.email,
			p.first_name, p.last_name, p.display_name, p.avatar_url, p.bio,
			p.city, p.state, p.date_of_birth, p.gender,
			u.phone_verified, u.status = 'active' as is_active,
			COALESCE(t.tier, 'bronze') as tier,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profiles p ON u.id = p.user_id
		LEFT JOIN user_tiers t ON u.id = t.user_id
		WHERE u.id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id.String())

	var (
		userID        string
		phone         string
		email         sql.NullString
		firstName     string
		lastName      string
		displayName   sql.NullString
		avatarURL     sql.NullString
		bio           sql.NullString
		city          sql.NullString
		state         sql.NullString
		dateOfBirth   sql.NullTime
		gender        sql.NullString
		phoneVerified bool
		isActive      bool
		tier          string
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := row.Scan(
		&userID, &phone, &email,
		&firstName, &lastName, &displayName, &avatarURL, &bio,
		&city, &state, &dateOfBirth, &gender,
		&phoneVerified, &isActive, &tier,
		&createdAt, &updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	// Reconstruct value objects
	phoneVO, err := valueobject.NewPhoneNumber(phone)
	if err != nil {
		return nil, err
	}

	fullName, err := valueobject.NewFullName(firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Reconstruct user aggregate
	// Note: This requires adding a Reconstitute method to the User aggregate
	user, err := aggregate.NewUser(
		valueobject.NewUserID(userID),
		phoneVO,
		fullName,
		"", // referral code - would need to fetch separately
	)
	if err != nil {
		return nil, err
	}

	// Set additional fields
	if email.Valid {
		emailVO, _ := valueobject.NewEmail(email.String)
		user.SetEmail(emailVO)
	}

	// Note: In a complete implementation, you would:
	// 1. Add a RehydrateUser factory method to the aggregate
	// 2. Fetch and attach skills
	// 3. Fetch referral information
	// 4. Set all other fields properly

	return user, nil
}

// FindByPhone retrieves a user by phone number
func (r *UserRepository) FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error) {
	query := `
		SELECT id
		FROM users
		WHERE phone = $1
	`

	var userID string
	err := r.db.QueryRowContext(ctx, query, phone.String()).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, valueobject.NewUserID(userID))
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*aggregate.User, error) {
	query := `
		SELECT id
		FROM users
		WHERE email = $1
	`

	var userID string
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, valueobject.NewUserID(userID))
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	query := `
		SELECT u.id
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE p.display_name = $1
	`

	var userID string
	err := r.db.QueryRowContext(ctx, query, username).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, valueobject.NewUserID(userID))
}

// FindByReferralCode retrieves a user by their referral code
func (r *UserRepository) FindByReferralCode(ctx context.Context, code string) (*aggregate.User, error) {
	query := `
		SELECT u.id
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		WHERE p.referral_code = $1
	`

	var userID string
	err := r.db.QueryRowContext(ctx, query, code).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, valueobject.NewUserID(userID))
}

// ExistsByPhone checks if a user with the phone exists
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, phone.String()).Scan(&exists)
	return exists, err
}

// ExistsByEmail checks if a user with the email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email.String()).Scan(&exists)
	return exists, err
}

// ExistsByUsername checks if a user with the username exists
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM profiles WHERE display_name = $1
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	return exists, err
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(ctx context.Context, id valueobject.UserID) error {
	query := `
		UPDATE users
		SET status = 'inactive',
			updated_at = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, time.Now().UTC(), id.String())
	return err
}
