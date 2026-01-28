# ADR-011: UUID Primary Keys for Database Models

## Status

Accepted

## Date

2024-01-15

## Context

HustleX needs to choose a primary key strategy for database tables. The system:
- May scale to multiple database instances
- Exposes IDs in URLs and API responses
- Handles sensitive financial records
- May require data migration between environments
- Will have high insert rates (transactions, notifications)

Options considered:
1. Auto-incrementing integers (SERIAL/BIGSERIAL)
2. UUIDs (v4 random, v7 time-ordered)
3. ULIDs (Universally Unique Lexicographically Sortable Identifier)
4. Snowflake IDs (Twitter-style)

## Decision

We chose **UUID v4** as the primary key type for all database models.

### Implementation:

```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    CreatedAt time.Time      `gorm:"autoCreateTime"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### Key Reasons:

1. **Distributed Generation**: IDs can be generated client-side or across multiple services without coordination.

2. **Security**: Sequential IDs leak information (user count, transaction volume); UUIDs are unpredictable.

3. **Merge Safety**: Data from different environments can be merged without ID conflicts.

4. **API Safety**: UUIDs in URLs don't expose business metrics (e.g., `/users/123` reveals ~123 users).

5. **PostgreSQL Native**: `gen_random_uuid()` function and UUID type are built-in.

## Consequences

### Positive

- **No coordination**: Generate IDs anywhere without database roundtrip
- **Privacy**: Competitors can't estimate transaction volume from IDs
- **Idempotency**: Pre-generate IDs for idempotent API requests
- **Flexibility**: Easy data migration, environment merging
- **Indexing**: PostgreSQL handles UUID indexes efficiently

### Negative

- **Size**: 16 bytes vs 4 bytes for INT (4x storage for ID columns)
- **Index size**: Larger B-tree indexes
- **Human readability**: `550e8400-e29b-41d4-a716-446655440000` vs `12345`
- **Fragmentation**: Random UUIDs cause index fragmentation (mitigated by v7)

### Neutral

- Requires UUID library for generation
- Some tools expect integer IDs
- Pagination slightly more complex (cursor-based preferred)

## UUID Types Comparison

| Type | Format | Sortable | Unique | Use Case |
|------|--------|----------|--------|----------|
| UUID v4 | Random | No | Yes | General purpose |
| UUID v7 | Time-ordered | Yes | Yes | Time-series data |
| ULID | Time-ordered | Yes | Yes | Log entries |
| Snowflake | Time + worker | Yes | Yes | High-scale Twitter-style |

### Why UUID v4 over v7?

- v4 is more widely supported
- Our queries primarily filter by user_id + date, not by primary key order
- v7 exposes creation time in ID (privacy concern for some use cases)
- PostgreSQL's `gen_random_uuid()` generates v4 natively

## Database Configuration

### PostgreSQL Setup

```sql
-- Enable uuid-ossp extension (if needed for specific UUID versions)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Use gen_random_uuid() for v4 (built into PostgreSQL 13+)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Index example
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
```

### GORM Configuration

```go
import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Phone     string    `gorm:"uniqueIndex;size:20"`
    Email     string    `gorm:"size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate hook to generate UUID if not set
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}
```

### API Response

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "+2348012345678",
  "wallet": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "balance": 50000.00
  }
}
```

## Performance Considerations

### Index Size Impact

| Records | INT PK Index | UUID PK Index | Difference |
|---------|--------------|---------------|------------|
| 1M | ~30 MB | ~50 MB | +67% |
| 10M | ~300 MB | ~500 MB | +67% |
| 100M | ~3 GB | ~5 GB | +67% |

### Mitigation Strategies

1. **Adequate RAM**: Ensure indexes fit in memory
2. **Selective Indexing**: Only index frequently queried columns
3. **Composite Indexes**: Use (user_id, created_at) instead of (id)
4. **BRIN Indexes**: For time-series tables, use Block Range Indexes

### Cursor Pagination

```go
// Instead of offset-based pagination
// GET /transactions?page=5&limit=20

// Use cursor-based pagination
// GET /transactions?cursor=2024-01-15T10:30:00Z&limit=20

func GetTransactions(userID uuid.UUID, cursor time.Time, limit int) ([]Transaction, error) {
    return db.Where("user_id = ? AND created_at < ?", userID, cursor).
        Order("created_at DESC").
        Limit(limit).
        Find(&transactions)
}
```

## Alternatives Considered

### Alternative 1: Auto-Increment Integer

**Pros**: Small size, fast, human-readable, natural ordering
**Cons**: Exposes record count, requires database for ID generation, merge conflicts

**Rejected because**: Security concerns (enumeration attacks) and distributed generation needs.

### Alternative 2: ULID

**Pros**: Lexicographically sortable, URL-safe, smaller than UUID string
**Cons**: Less ecosystem support, requires external library

**Rejected because**: UUID has better PostgreSQL native support.

### Alternative 3: Snowflake ID

**Pros**: Time-ordered, compact (64-bit), high throughput
**Cons**: Requires worker ID coordination, reveals creation time, custom implementation

**Rejected because**: Coordination complexity outweighs benefits at our scale.

### Alternative 4: Hybrid (UUID externally, INT internally)

**Pros**: Best of both worlds
**Cons**: Dual columns, mapping complexity, migration overhead

**Rejected because**: Added complexity not justified.

## References

- [PostgreSQL UUID Type](https://www.postgresql.org/docs/current/datatype-uuid.html)
- [UUID vs Auto-Increment Performance](https://www.percona.com/blog/uuid-vs-int-performance-in-mysql/)
- [UUID v7 Specification](https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format)
- [Google UUID Library](https://github.com/google/uuid)
