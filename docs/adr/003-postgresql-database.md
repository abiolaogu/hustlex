# ADR-003: PostgreSQL as Primary Database

## Status

Accepted

## Date

2024-01-15

## Context

HustleX handles sensitive financial data including:
- User wallets and balances
- Transaction histories
- Savings circle contributions
- Credit scores and loan records
- Escrow accounts

We needed a database that:
1. Ensures ACID compliance for financial transactions
2. Supports complex queries for analytics and reporting
3. Handles high concurrent read/write operations
4. Provides robust data integrity constraints
5. Scales with growing user base
6. Has strong ecosystem support in Go

## Decision

We chose **PostgreSQL 15+** as our primary relational database.

### Key Reasons:

1. **ACID Compliance**: Full transaction support ensures financial data integrity (critical for wallet operations and escrow).

2. **Advanced Features**:
   - JSONB columns for flexible schema (notifications, metadata)
   - Full-text search for gig listings
   - UUID support for distributed-system-friendly primary keys
   - Array types for multi-value fields

3. **Performance**:
   - Excellent query optimizer
   - Partial indexes for query acceleration
   - Connection pooling support (PgBouncer compatible)
   - Parallel query execution

4. **Reliability**:
   - Point-in-time recovery (PITR)
   - Streaming replication
   - Proven in banking/financial systems

5. **Go Ecosystem**: GORM provides excellent PostgreSQL support with migrations, hooks, and query building.

## Consequences

### Positive

- **Data integrity**: Foreign keys, constraints, and transactions prevent invalid states
- **Flexible querying**: Complex JOINs for reports (gig performance, savings analytics)
- **JSONB support**: Store notification payloads, API responses without schema changes
- **Full-text search**: Native search on gig titles/descriptions without Elasticsearch
- **Scalability**: Read replicas, connection pooling, and table partitioning options
- **Cost-effective**: Open source with excellent cloud support (AWS RDS, Google Cloud SQL)

### Negative

- **Operational complexity**: Requires DBA knowledge for optimization and maintenance
- **Vertical scaling limits**: Horizontal sharding is complex (may need Citus extension)
- **Connection overhead**: Each connection consumes ~10MB RAM
- **Learning curve**: Advanced features (CTEs, window functions) require SQL expertise

### Neutral

- SQL query language (standard, transferable skill)
- Requires index tuning for optimal performance
- Schema migrations need careful planning for zero-downtime deploys

## Alternatives Considered

### Alternative 1: MySQL/MariaDB

**Pros**: Widely deployed, simple setup, large community
**Cons**: Weaker ACID guarantees in some configurations, limited JSON support, less advanced query optimizer

**Rejected because**: PostgreSQL's stronger transactional guarantees and advanced features better suit financial applications.

### Alternative 2: MongoDB

**Pros**: Flexible schema, horizontal scaling, document model
**Cons**: No ACID transactions (pre-4.0), eventual consistency, complex joins

**Rejected because**: Financial transactions require strong consistency guarantees that document databases don't natively provide.

### Alternative 3: CockroachDB

**Pros**: Distributed SQL, automatic scaling, PostgreSQL compatible
**Cons**: Operational complexity, higher cost, smaller community

**Rejected because**: Overkill for initial launch; can migrate later if needed.

### Alternative 4: SQLite

**Pros**: Zero configuration, embedded, fast for reads
**Cons**: Single-writer limitation, no network access, not suitable for concurrent web apps

**Rejected because**: Doesn't support concurrent writes required for multi-user financial platform.

## Implementation Details

### Connection Configuration

```go
// Connection pooling settings
MaxOpenConns: 25      // Maximum open connections
MaxIdleConns: 5       // Idle connections kept alive
MaxLifetime: 5min     // Connection recycling
```

### Key Indexes

```sql
-- Composite indexes for query optimization
CREATE INDEX idx_transactions_user_date ON transactions(user_id, created_at DESC);
CREATE INDEX idx_gigs_category_status ON gigs(category, status);
CREATE INDEX idx_contributions_circle_date ON contributions(circle_id, due_date);

-- Full-text search
CREATE INDEX idx_gigs_search ON gigs USING gin(to_tsvector('english', title || ' ' || description));
```

## References

- [PostgreSQL Official Documentation](https://www.postgresql.org/docs/)
- [GORM PostgreSQL Driver](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)
- [PostgreSQL ACID Compliance](https://www.postgresql.org/docs/current/mvcc.html)
- [PgBouncer Connection Pooling](https://www.pgbouncer.org/)
