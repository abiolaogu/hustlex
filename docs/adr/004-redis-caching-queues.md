# ADR-004: Redis for Caching and Job Queues

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires:
1. **Session/Token Storage**: Fast access to JWT blacklists and refresh tokens
2. **Rate Limiting**: Per-user and per-endpoint request throttling
3. **OTP Management**: Temporary storage for one-time passwords with TTL
4. **Caching**: Reduce database load for frequently accessed data
5. **Job Queues**: Background processing for notifications, payments, scheduled tasks
6. **Real-time Features**: Pub/sub for future live notifications

We needed a fast, reliable in-memory data store that integrates well with Go.

## Decision

We chose **Redis 7+** as our caching layer and backing store for job queues.

### Key Reasons:

1. **Sub-millisecond Latency**: In-memory storage provides <1ms read/write operations.

2. **Rich Data Structures**:
   - Strings: OTP codes, JWT tokens
   - Hashes: User session data
   - Sorted Sets: Rate limiting sliding windows
   - Lists: Job queues
   - Pub/Sub: Real-time event distribution

3. **TTL Support**: Automatic expiration for OTPs, rate limit windows, and cache entries.

4. **Asynq Compatibility**: Native Redis backend for our chosen job queue library.

5. **Atomic Operations**: INCR, SETNX, and Lua scripts enable race-condition-free rate limiting.

6. **Persistence Options**: RDB snapshots and AOF logs prevent data loss.

## Consequences

### Positive

- **Performance**: 100,000+ ops/second per instance
- **Simplicity**: Single dependency for caching AND job queues
- **Reliability**: Battle-tested in high-scale production systems
- **Flexibility**: Multiple data structures for different use cases
- **Observability**: Built-in MONITOR, SLOWLOG, and INFO commands
- **Clustering**: Redis Cluster for horizontal scaling when needed

### Negative

- **Memory cost**: All data in RAM (can be expensive at scale)
- **Single-threaded**: Main event loop is single-threaded (mitigated with Redis 6+ I/O threads)
- **Data loss risk**: Without persistence, crashes lose in-memory data
- **Complexity**: Clustering adds operational overhead

### Neutral

- Requires separate infrastructure from PostgreSQL
- Lua scripting for complex atomic operations
- Connection pooling recommended for Go clients

## Use Cases in HustleX

### 1. OTP Storage

```
Key: otp:{phone}
Value: {code}
TTL: 5 minutes
```

### 2. Rate Limiting (Sliding Window)

```
Key: ratelimit:{user_id}:{endpoint}:{window_start}
Value: request_count
TTL: window_duration + buffer
```

### 3. JWT Blacklist

```
Key: jwt:blacklist:{token_jti}
Value: 1
TTL: token_remaining_lifetime
```

### 4. Job Queues (Asynq)

```
Queues:
- critical: Financial transactions (6 workers)
- default: Savings, loans, reminders (3 workers)
- low: Notifications, cleanup (1 worker)
```

### 5. Caching

```
Key: cache:user:{user_id}
Value: JSON user profile
TTL: 15 minutes

Key: cache:gigs:category:{category}
Value: JSON gig listings
TTL: 5 minutes
```

## Alternatives Considered

### Alternative 1: Memcached

**Pros**: Simple, fast, multi-threaded
**Cons**: Limited data structures (only strings), no persistence, no pub/sub

**Rejected because**: Need for rich data structures (sorted sets for rate limiting, lists for queues).

### Alternative 2: In-Process Cache (Go map)

**Pros**: No external dependency, fastest possible access
**Cons**: No persistence, doesn't scale across multiple instances, no TTL support

**Rejected because**: Horizontal scaling requires shared state across API instances.

### Alternative 3: PostgreSQL for Everything

**Pros**: Single database, ACID transactions
**Cons**: Slower than in-memory stores, polling for job queues, higher connection overhead

**Rejected because**: Performance requirements for rate limiting and caching.

### Alternative 4: RabbitMQ for Queues

**Pros**: Advanced routing, message acknowledgment, durable queues
**Cons**: Additional infrastructure, different protocol, overkill for our use case

**Rejected because**: Asynq with Redis provides sufficient functionality with simpler architecture.

## References

- [Redis Official Documentation](https://redis.io/documentation)
- [Asynq - Distributed Task Queue](https://github.com/hibiken/asynq)
- [Redis Rate Limiting Patterns](https://redis.io/commands/incr#pattern-rate-limiter)
- [go-redis Client](https://github.com/redis/go-redis)
