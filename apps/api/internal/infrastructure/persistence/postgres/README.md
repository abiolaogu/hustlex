# PostgreSQL Repository Implementation

This package contains PostgreSQL implementations of domain repositories following the Repository Pattern and Clean Architecture principles.

## Architecture

```
Domain Layer (interfaces)
    ↓
Infrastructure Layer (implementations)
    ↓
PostgreSQL Database
```

## Pattern Overview

### Repository Interface (Domain Layer)
Located in `internal/domain/{module}/repository/`

Example: `internal/domain/identity/repository/user_repository.go`

```go
type UserRepository interface {
    Save(ctx context.Context, user *aggregate.User) error
    FindByID(ctx context.Context, id valueobject.UserID) (*aggregate.User, error)
    // ... other methods
}
```

### Repository Implementation (Infrastructure Layer)
Located in `internal/infrastructure/persistence/postgres/`

Example: `internal/infrastructure/persistence/postgres/user_repository.go`

```go
type UserRepository struct {
    db *DB
}

func NewUserRepository(db *DB) repository.UserRepository {
    return &UserRepository{db: db}
}
```

## Key Features

### 1. Optimistic Locking
All aggregates use version-based optimistic locking to prevent concurrent update conflicts:

```go
UPDATE users SET
    ... fields ...,
    version = version + 1
WHERE id = $1 AND version = $2
```

If `version` doesn't match, the update fails, indicating a concurrent modification.

### 2. Soft Deletes
Entities are soft-deleted using a `deleted_at` timestamp:

```go
UPDATE users SET deleted_at = NOW() WHERE id = $1
```

All queries filter out soft-deleted records:

```go
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL
```

### 3. Transaction Support
The `DB.WithTransaction()` helper manages database transactions:

```go
err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
    // Multiple operations within transaction
    return nil
})
```

### 4. Event Sourcing (Planned)
The `SaveWithEvents()` method is prepared for event publishing:

```go
func (r *UserRepository) SaveWithEvents(ctx context.Context, user *aggregate.User) error {
    return r.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Save aggregate
        if err := r.Save(ctx, user); err != nil {
            return err
        }

        // Publish domain events (TODO: integrate event bus)
        // events := user.DomainEvents()
        // for _, event := range events {
        //     eventBus.Publish(ctx, event)
        // }

        return nil
    })
}
```

## Database Configuration

Connection settings are managed through `Config`:

```go
cfg := postgres.Config{
    Host:            "localhost",
    Port:            5432,
    User:            "hustlex",
    Password:        "secret",
    DBName:          "hustlex_db",
    SSLMode:         "disable",
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: 5 * time.Minute,
}

db, err := postgres.NewDB(cfg)
```

## Schema Mapping

### User Aggregate Mapping

| Aggregate Field | Database Column | Type | Notes |
|----------------|-----------------|------|-------|
| `id` | `id` | UUID | Primary key |
| `phone` | `phone` | VARCHAR(20) | Unique, required |
| `email` | `email` | VARCHAR(255) | Unique, optional |
| `fullName` | `full_name` | VARCHAR(200) | Required |
| `tier` | `tier` | VARCHAR(20) | bronze/silver/gold/platinum |
| `isActive` | `status` | user_status ENUM | active/inactive/suspended |
| `skills` | `user_skills` table | One-to-many | Separate table |
| `version` | `version` | BIGINT | Optimistic locking |

### Skills Relationship

User skills are stored in a separate `user_skills` table with a one-to-many relationship:

```sql
CREATE TABLE user_skills (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    skill_id UUID,
    skill_name VARCHAR(100),
    proficiency VARCHAR(20),
    years_exp INT,
    is_verified BOOLEAN,
    portfolio_urls JSONB,
    added_at TIMESTAMPTZ
);
```

## Implementation Checklist

When implementing a new repository, ensure:

- [ ] Interface defined in domain layer (`internal/domain/{module}/repository/`)
- [ ] Implementation in infrastructure layer (`internal/infrastructure/persistence/postgres/`)
- [ ] Constructor function `New{Entity}Repository(db *DB)`
- [ ] All interface methods implemented
- [ ] Optimistic locking on updates
- [ ] Soft delete support
- [ ] Transaction support for complex operations
- [ ] Proper error wrapping with context
- [ ] Null handling for optional fields
- [ ] Related entities (one-to-many, many-to-many) handled correctly

## Testing

Repository tests should:

1. Use a test database or in-memory PostgreSQL (pgx/testdb)
2. Test all CRUD operations
3. Test concurrent updates (optimistic locking)
4. Test error cases (not found, constraint violations)
5. Test transaction rollback scenarios

Example:

```go
func TestUserRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := NewUserRepository(db)

    user, _ := aggregate.NewUser(...)
    err := repo.Save(context.Background(), user)

    assert.NoError(t, err)
}
```

## Next Steps

The following repositories need implementation following the User repository pattern:

1. **Wallet Repository** (`wallet_repository.go`)
   - Maps: `Wallet` aggregate → `wallets` table
   - Handles: Multi-currency balances, escrow accounts

2. **Gig Repository** (`gig_repository.go`)
   - Maps: `Gig` aggregate → `gigs` table
   - Handles: Proposals, contracts, reviews

3. **Circle Repository** (`circle_repository.go`)
   - Maps: `SavingsCircle` aggregate → `savings_circles` table
   - Handles: Members, contributions, payouts

4. **Notification Repository** (`notification_repository.go`)
   - Maps: `Notification` aggregate → `notifications` table
   - Handles: Delivery status, retries

5. **Credit Repository** (`credit_repository.go`)
   - Maps: `CreditScore`, `Loan` → `credit_scores`, `loans` tables
   - Handles: Score calculation, loan lifecycle

## Dependencies

- `github.com/lib/pq` - PostgreSQL driver
- Domain layer interfaces (no external dependencies)

## Resources

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
