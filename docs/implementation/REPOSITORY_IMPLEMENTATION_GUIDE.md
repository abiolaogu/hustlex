# Database Repository Implementation Guide

## Executive Summary

This document provides a detailed implementation guide for the PostgreSQL repositories required for HustleX MVP. This is identified as a **Tier 1 Launch Blocker** in the PRD (estimated effort: 15-20 days).

**Status**: Implementation guide created (February 5, 2026)
**Next Action**: Developer to implement repositories following this template

---

## Overview

The repository pattern provides a clean abstraction between the domain layer and the persistence layer, following Clean Architecture principles. Each repository implements interfaces defined in `internal/domain/*/repository/`.

### Repository Responsibilities

1. **Data Persistence**: Save and retrieve domain aggregates
2. **Query Operations**: Find, search, filter operations
3. **Transaction Management**: Ensure data consistency
4. **Event Publishing**: Publish domain events after persistence
5. **Error Handling**: Translate database errors to domain errors

---

## Architecture

```
apps/api/
└── internal/
    ├── domain/                      # Domain layer (existing)
    │   ├── identity/repository/     # Interfaces
    │   ├── wallet/repository/
    │   ├── gig/repository/
    │   └── savings/repository/
    │
    └── infrastructure/              # Infrastructure layer
        └── persistence/
            └── postgres/            # PostgreSQL implementations
                ├── db.go                    # Database connection & pooling
                ├── tx.go                    # Transaction management
                ├── errors.go                # Error mapping
                ├── user_repository.go       # User aggregate
                ├── wallet_repository.go     # Wallet aggregate
                ├── gig_repository.go        # Gig aggregate
                ├── circle_repository.go     # Circle aggregate
                ├── otp_repository.go        # OTP persistence
                ├── session_repository.go    # Session management
                └── mappers/                 # Domain ↔ DB mappers
                    ├── user_mapper.go
                    ├── wallet_mapper.go
                    ├── gig_mapper.go
                    └── circle_mapper.go
```

---

## Implementation Priority

### Phase 1: Foundation (Days 1-3)
1. **Database Connection** (`db.go`)
   - Connection pooling with pgxpool
   - Health checks
   - Configuration management

2. **Transaction Management** (`tx.go`)
   - Transaction wrapper
   - Rollback handling
   - Nested transaction support

3. **Error Mapping** (`errors.go`)
   - Map pgx errors to domain errors
   - Handle constraint violations
   - Deadlock detection

### Phase 2: Core Repositories (Days 4-10)
Priority order based on MVP dependencies:

1. **User Repository** (Days 4-5)
   - Required for: Authentication, all user operations
   - Dependencies: None
   - Complexity: Medium
   - Files: `user_repository.go`, `otp_repository.go`, `session_repository.go`

2. **Wallet Repository** (Days 6-7)
   - Required for: All financial operations
   - Dependencies: User
   - Complexity: High (transactions, locks, escrow)
   - Files: `wallet_repository.go`, `transaction_repository.go`

3. **Gig Repository** (Days 8-9)
   - Required for: Marketplace operations
   - Dependencies: User, Wallet
   - Complexity: High (proposals, contracts, search)
   - Files: `gig_repository.go`, `proposal_repository.go`

4. **Circle Repository** (Days 10)
   - Required for: Savings circles
   - Dependencies: User, Wallet
   - Complexity: Medium (contributions, rotations)
   - Files: `circle_repository.go`, `contribution_repository.go`

### Phase 3: Additional Repositories (Days 11-15)
- Skill Repository
- Referral Repository
- Notification Repository
- Credit Repository

### Phase 4: Testing & Optimization (Days 16-20)
- Unit tests for each repository
- Integration tests with test database
- Performance optimization (indexes, query tuning)
- Load testing

---

## Implementation Template: User Repository

### File: `internal/infrastructure/persistence/postgres/user_repository.go`

```go
package postgres

import (
    "context"
    "errors"
    "fmt"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "hustlex/internal/domain/identity/aggregate"
    "hustlex/internal/domain/identity/repository"
    "hustlex/internal/domain/shared/valueobject"
)

type UserRepository struct {
    pool *pgxpool.Pool
    eventPublisher EventPublisher // To publish domain events
}

func NewUserRepository(pool *pgxpool.Pool, eventPublisher EventPublisher) *UserRepository {
    return &UserRepository{
        pool: pool,
        eventPublisher: eventPublisher,
    }
}

// Save persists a user aggregate
func (r *UserRepository) Save(ctx context.Context, user *aggregate.User) error {
    query := `
        INSERT INTO users (
            id, phone, email, username, full_name,
            profile_image, bio, location, state, date_of_birth, gender,
            is_verified, is_active, tier, referral_code, referred_by,
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
            is_active = EXCLUDED.is_active,
            tier = EXCLUDED.tier,
            last_login_at = EXCLUDED.last_login_at,
            updated_at = EXCLUDED.updated_at,
            version = EXCLUDED.version
        WHERE users.version = $20 - 1
    `

    // Extract values from aggregate
    id := user.ID().String()
    phone := user.Phone().String()
    email := ""
    if user.Email() != nil {
        email = user.Email().String()
    }

    result, err := r.pool.Exec(ctx, query,
        id, phone, email, user.Username(), user.FullName().String(),
        user.ProfileImage(), user.Bio(), user.Location(), user.State(),
        user.DateOfBirth(), user.Gender(),
        user.IsVerified(), user.IsActive(), user.Tier().String(),
        user.ReferralCode(), user.ReferredBy(),
        user.LastLoginAt(), user.CreatedAt(), user.UpdatedAt(), user.Version(),
    )

    if err != nil {
        return mapError(err)
    }

    if result.RowsAffected() == 0 {
        return repository.ErrOptimisticLock
    }

    // Save user skills
    if err := r.saveUserSkills(ctx, user); err != nil {
        return err
    }

    return nil
}

// SaveWithEvents persists user and publishes domain events
func (r *UserRepository) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Save user in transaction
    if err := r.saveInTx(ctx, tx, user); err != nil {
        return err
    }

    // Publish events (within transaction or after commit depending on requirements)
    events := user.DomainEvents()
    for _, event := range events {
        if err := r.eventPublisher.Publish(ctx, event); err != nil {
            return err
        }
    }

    user.ClearEvents()

    return tx.Commit(ctx)
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error) {
    query := `
        SELECT
            id, phone, email, username, full_name,
            profile_image, bio, location, state, date_of_birth, gender,
            is_verified, is_active, tier, referral_code, referred_by,
            last_login_at, created_at, updated_at, version
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `

    var user aggregate.User
    var emailStr, usernameStr, profileImage, bio, location, state, gender *string
    var dateOfBirth, lastLoginAt *time.Time
    var referredBy *string

    err := r.pool.QueryRow(ctx, query, id.String()).Scan(
        &user.id, &user.phone, &emailStr, &usernameStr, &user.fullName,
        &profileImage, &bio, &location, &state, &dateOfBirth, &gender,
        &user.isVerified, &user.isActive, &user.tier, &user.referralCode, &referredBy,
        &lastLoginAt, &user.createdAt, &user.updatedAt, &user.version,
    )

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, repository.ErrUserNotFound
        }
        return nil, mapError(err)
    }

    // Load user skills
    skills, err := r.loadUserSkills(ctx, id)
    if err != nil {
        return nil, err
    }
    user.skills = skills

    return &user, nil
}

// FindByPhone retrieves a user by phone number
func (r *UserRepository) FindByPhone(ctx context.Context, phone valueobject.PhoneNumber) (*aggregate.User, error) {
    // Similar implementation to FindByID
    // ...
}

// ExistsByPhone checks if a user with the phone exists
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone valueobject.PhoneNumber) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1 AND deleted_at IS NULL)`

    var exists bool
    err := r.pool.QueryRow(ctx, query, phone.String()).Scan(&exists)
    if err != nil {
        return false, mapError(err)
    }

    return exists, nil
}

// Helper: Save user skills
func (r *UserRepository) saveUserSkills(ctx context.Context, user *aggregate.User) error {
    // Delete existing skills
    deleteQuery := `DELETE FROM user_skills WHERE user_id = $1`
    if _, err := r.pool.Exec(ctx, deleteQuery, user.ID().String()); err != nil {
        return err
    }

    // Insert current skills
    insertQuery := `
        INSERT INTO user_skills (
            user_id, skill_id, skill_name, proficiency,
            years_exp, is_verified, portfolio_urls, added_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

    for _, skill := range user.Skills() {
        _, err := r.pool.Exec(ctx, insertQuery,
            user.ID().String(),
            skill.SkillID.String(),
            skill.SkillName,
            skill.Proficiency,
            skill.YearsExp,
            skill.IsVerified,
            skill.PortfolioURLs,
            skill.AddedAt,
        )
        if err != nil {
            return err
        }
    }

    return nil
}

// Helper: Load user skills
func (r *UserRepository) loadUserSkills(ctx context.Context, userID valueobject.UserID) ([]aggregate.UserSkill, error) {
    query := `
        SELECT skill_id, skill_name, proficiency, years_exp,
               is_verified, portfolio_urls, added_at
        FROM user_skills
        WHERE user_id = $1
        ORDER BY added_at DESC
    `

    rows, err := r.pool.Query(ctx, query, userID.String())
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var skills []aggregate.UserSkill
    for rows.Next() {
        var skill aggregate.UserSkill
        err := rows.Scan(
            &skill.SkillID,
            &skill.SkillName,
            &skill.Proficiency,
            &skill.YearsExp,
            &skill.IsVerified,
            &skill.PortfolioURLs,
            &skill.AddedAt,
        )
        if err != nil {
            return nil, err
        }
        skills = append(skills, skill)
    }

    return skills, nil
}
```

---

## Database Connection Setup

### File: `internal/infrastructure/persistence/postgres/db.go`

```go
package postgres

import (
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
    Host     string
    Port     int
    Database string
    Username string
    Password string
    SSLMode  string

    MaxConns          int32
    MinConns          int32
    MaxConnLifetime   time.Duration
    MaxConnIdleTime   time.Duration
    HealthCheckPeriod time.Duration
}

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
        cfg.SSLMode,
    )

    poolConfig, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("unable to parse connection string: %w", err)
    }

    // Connection pool settings
    poolConfig.MaxConns = cfg.MaxConns
    poolConfig.MinConns = cfg.MinConns
    poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
    poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
    poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

    // Verify connection
    if err := pool.Ping(ctx); err != nil {
        pool.Close()
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }

    return pool, nil
}

func DefaultConfig() Config {
    return Config{
        Host:              "localhost",
        Port:              5432,
        Database:          "hustlex",
        Username:          "postgres",
        Password:          "postgres",
        SSLMode:           "disable",
        MaxConns:          25,
        MinConns:          5,
        MaxConnLifetime:   time.Hour,
        MaxConnIdleTime:   time.Minute * 30,
        HealthCheckPeriod: time.Minute,
    }
}
```

### File: `internal/infrastructure/persistence/postgres/errors.go`

```go
package postgres

import (
    "errors"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"

    "hustlex/internal/domain/identity/repository"
)

func mapError(err error) error {
    if err == nil {
        return nil
    }

    // No rows found
    if errors.Is(err, pgx.ErrNoRows) {
        return repository.ErrNotFound
    }

    // PostgreSQL errors
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505": // unique_violation
            return repository.ErrAlreadyExists
        case "23503": // foreign_key_violation
            return repository.ErrForeignKeyViolation
        case "23514": // check_violation
            return repository.ErrCheckViolation
        case "40001": // serialization_failure
            return repository.ErrSerializationFailure
        case "40P01": // deadlock_detected
            return repository.ErrDeadlock
        }
    }

    return err
}
```

---

## Transaction Management

### File: `internal/infrastructure/persistence/postgres/tx.go`

```go
package postgres

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TxFunc func(context.Context, pgx.Tx) error

// WithTransaction executes a function within a database transaction
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            panic(p)
        }
    }()

    if err := fn(ctx, tx); err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return fmt.Errorf("rollback transaction: %w (original error: %v)", rbErr, err)
        }
        return err
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }

    return nil
}
```

---

## Testing Strategy

### Unit Tests
- Mock pgxpool.Pool using interfaces
- Test error mapping
- Test data transformation (aggregate ↔ DB)

### Integration Tests
- Use test containers (testcontainers-go)
- Run against real PostgreSQL
- Test transaction rollback scenarios
- Test concurrent operations
- Test optimistic locking

### Example Test Setup

```go
package postgres_test

import (
    "context"
    "testing"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
    ctx := context.Background()

    postgresContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:16-alpine"),
        postgres.WithDatabase("hustlex_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)

    t.Cleanup(func() {
        _ = postgresContainer.Terminate(ctx)
    })

    connStr, err := postgresContainer.ConnectionString(ctx)
    require.NoError(t, err)

    pool, err := pgxpool.New(ctx, connStr)
    require.NoError(t, err)

    // Run migrations
    runMigrations(t, pool)

    return pool
}

func TestUserRepository_Save(t *testing.T) {
    pool := setupTestDB(t)
    defer pool.Close()

    repo := postgres.NewUserRepository(pool, &mockEventPublisher{})

    // Create test user
    user, err := aggregate.NewUser(
        valueobject.NewUserID(),
        valueobject.MustNewPhoneNumber("+2348012345678"),
        valueobject.MustNewFullName("John", "Doe"),
        "ABC123",
    )
    require.NoError(t, err)

    // Save user
    err = repo.Save(context.Background(), user)
    require.NoError(t, err)

    // Retrieve and verify
    retrieved, err := repo.FindByID(context.Background(), user.ID())
    require.NoError(t, err)
    assert.Equal(t, user.ID(), retrieved.ID())
    assert.Equal(t, user.Phone(), retrieved.Phone())
}
```

---

## Implementation Checklist

### Prerequisites
- [ ] Install pgx driver: `go get github.com/jackc/pgx/v5`
- [ ] Install pgxpool: `go get github.com/jackc/pgx/v5/pgxpool`
- [ ] Review database schema in `backend/hasura/migrations/default/1_init/up.sql`
- [ ] Ensure PostgreSQL is running (Docker or local)

### Foundation
- [ ] Create `internal/infrastructure/persistence/postgres/` directory
- [ ] Implement `db.go` (connection pooling)
- [ ] Implement `tx.go` (transaction management)
- [ ] Implement `errors.go` (error mapping)
- [ ] Create database configuration loader

### User Repository
- [ ] Implement `user_repository.go`
- [ ] Implement `otp_repository.go`
- [ ] Implement `session_repository.go`
- [ ] Implement `skill_repository.go`
- [ ] Implement `referral_repository.go`
- [ ] Write unit tests
- [ ] Write integration tests

### Wallet Repository
- [ ] Implement `wallet_repository.go`
- [ ] Implement `transaction_repository.go`
- [ ] Handle escrow operations
- [ ] Handle multi-currency balances
- [ ] Implement transaction locking (SELECT FOR UPDATE)
- [ ] Write unit tests
- [ ] Write integration tests

### Gig Repository
- [ ] Implement `gig_repository.go`
- [ ] Implement `proposal_repository.go`
- [ ] Implement `contract_repository.go`
- [ ] Implement search functionality
- [ ] Implement filtering and pagination
- [ ] Write unit tests
- [ ] Write integration tests

### Circle Repository
- [ ] Implement `circle_repository.go`
- [ ] Implement `circle_member_repository.go`
- [ ] Implement `contribution_repository.go`
- [ ] Handle rotation logic
- [ ] Write unit tests
- [ ] Write integration tests

### Integration & Testing
- [ ] Setup test containers
- [ ] Write integration tests for all repositories
- [ ] Test transaction rollback scenarios
- [ ] Test optimistic locking
- [ ] Test concurrent operations
- [ ] Performance testing with large datasets

### Documentation
- [ ] Document repository patterns
- [ ] Document transaction patterns
- [ ] Document error handling
- [ ] Add code examples
- [ ] Update API documentation

---

## Performance Considerations

### Connection Pooling
- **MaxConns**: 25 (adjust based on load)
- **MinConns**: 5 (keep warm connections)
- **MaxConnLifetime**: 1 hour
- **MaxConnIdleTime**: 30 minutes

### Query Optimization
1. **Use Indexes**: Ensure all foreign keys and frequently queried fields are indexed
2. **Batch Operations**: Use batch inserts for multiple records
3. **SELECT Only What You Need**: Avoid `SELECT *`
4. **Use Prepared Statements**: pgx automatically uses prepared statements
5. **Connection Reuse**: Always use connection pool, never create individual connections

### Locking Strategy
- **Optimistic Locking**: Use version field for user updates
- **Pessimistic Locking**: Use `SELECT FOR UPDATE` for wallet operations
- **Row-Level Locking**: PostgreSQL provides fine-grained locking

---

## Dependencies

### Required Go Packages
```bash
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
go get github.com/stretchr/testify
```

### Database
- PostgreSQL 16+
- Connection string format: `postgres://username:password@host:port/database?sslmode=disable`

---

## Next Steps

1. **Review this guide** with the team
2. **Setup development environment** (PostgreSQL, Go dependencies)
3. **Start with User Repository** (template for others)
4. **Write tests as you go** (TDD approach)
5. **Iterate on pattern** before implementing other repositories
6. **Code review** after each repository
7. **Integration testing** with real database
8. **Performance profiling** with production-like data

---

## Estimated Timeline

| Task | Effort | Priority |
|------|--------|----------|
| Foundation (db, tx, errors) | 2 days | P0 |
| User Repository + Tests | 3 days | P0 |
| Wallet Repository + Tests | 3 days | P0 |
| Gig Repository + Tests | 3 days | P0 |
| Circle Repository + Tests | 2 days | P1 |
| Additional Repositories | 4 days | P1 |
| Integration Testing | 2 days | P0 |
| Performance Optimization | 1 day | P1 |
| **Total** | **20 days** | |

---

## References

- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

---

*Document created: February 5, 2026*
*Author: Claude (Autonomous Factory)*
*Status: Ready for implementation*
