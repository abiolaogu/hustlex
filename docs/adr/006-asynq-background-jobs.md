# ADR-006: Asynq for Background Job Processing

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires reliable background job processing for:

1. **Financial Operations**:
   - Scheduled savings contributions (auto-debit)
   - Loan repayment processing
   - Escrow release after gig completion
   - Payment webhook handling

2. **Notifications**:
   - Push notifications (FCM)
   - SMS alerts (payment reminders, OTPs)
   - Email notifications

3. **Scheduled Tasks**:
   - Credit score recalculation (weekly)
   - Savings circle payout scheduling
   - OTP cleanup (hourly)
   - Gig deadline reminders

Requirements:
- Reliable delivery (at-least-once)
- Retry with backoff
- Scheduled/delayed jobs
- Priority queues
- Monitoring dashboard
- Go-native integration

## Decision

We chose **Asynq** as our background job processing library.

### Key Reasons:

1. **Redis-Based**: Uses our existing Redis infrastructure, no additional dependencies.

2. **Go-Native**: Written in Go, idiomatic API, strong typing for task payloads.

3. **Reliable**: Guaranteed at-least-once execution with configurable retries.

4. **Flexible Scheduling**: Supports immediate, delayed, and cron-scheduled tasks.

5. **Priority Queues**: Multiple queue priorities with configurable concurrency.

6. **Built-in Monitoring**: Asynqmon provides web UI for job monitoring.

## Consequences

### Positive

- **Simplicity**: Single dependency (Redis) for job queue
- **Performance**: Handles thousands of jobs/second
- **Reliability**: Automatic retries, dead letter queue
- **Observability**: Built-in metrics and web dashboard
- **Flexibility**: Immediate, scheduled, and periodic tasks
- **Type safety**: Go structs for task payloads

### Negative

- **Single point of failure**: Redis outage affects all jobs
- **At-least-once semantics**: Tasks may run multiple times (need idempotency)
- **Limited routing**: No complex routing rules (vs RabbitMQ)
- **Go-only**: Can't process jobs from other languages

### Neutral

- Requires careful error handling for idempotency
- Redis persistence configuration affects job durability
- Monitoring dashboard requires separate deployment

## Implementation Details

### Queue Configuration

```go
// Queue priorities and concurrency
srv := asynq.NewServer(
    redisOpt,
    asynq.Config{
        Queues: map[string]int{
            "critical": 6,  // Financial transactions
            "default":  3,  // Savings, loans
            "low":      1,  // Notifications
        },
        RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
            return time.Duration(n*n) * time.Minute // Exponential backoff
        },
    },
)
```

### Task Types

| Task | Queue | Max Retries | Description |
|------|-------|-------------|-------------|
| `savings:contribution_reminder` | default | 3 | Due date reminders |
| `savings:process_contribution` | critical | 5 | Auto-debit contributions |
| `savings:process_payout` | critical | 5 | Rotational payouts |
| `gig:escrow_release` | critical | 5 | Payment release |
| `loan:process_repayment` | critical | 5 | Loan payments |
| `notification:push` | low | 3 | Firebase push |
| `notification:sms` | default | 3 | SMS delivery |
| `system:cleanup_expired_otps` | low | 1 | Hourly cleanup |

### Task Payload Example

```go
// Task definition
type ContributionReminderPayload struct {
    CircleID       string    `json:"circle_id"`
    UserID         string    `json:"user_id"`
    Amount         float64   `json:"amount"`
    DueDate        time.Time `json:"due_date"`
    ReminderNumber int       `json:"reminder_number"`
}

// Enqueue task
task, err := asynq.NewTask(
    "savings:contribution_reminder",
    payload,
    asynq.Queue("default"),
    asynq.ProcessIn(24*time.Hour),  // Delayed execution
    asynq.MaxRetry(3),
    asynq.Timeout(30*time.Second),
)
client.Enqueue(task)
```

### Idempotency Pattern

```go
func HandleEscrowRelease(ctx context.Context, t *asynq.Task) error {
    var p EscrowReleasePayload
    json.Unmarshal(t.Payload(), &p)

    // Idempotency check
    processed, err := redis.SetNX(ctx,
        "escrow:processed:"+p.ContractID,
        "1",
        24*time.Hour,
    ).Result()

    if !processed {
        return nil // Already processed, skip
    }

    // Process escrow release...
    return nil
}
```

### Scheduled Jobs (Cron)

```go
scheduler := asynq.NewScheduler(redisOpt, nil)

// Weekly credit score recalculation
scheduler.Register("0 0 * * 0", // Every Sunday at midnight
    asynq.NewTask("system:credit_score_batch", nil),
    asynq.Queue("low"),
)

// Hourly OTP cleanup
scheduler.Register("0 * * * *", // Every hour
    asynq.NewTask("system:cleanup_expired_otps", nil),
    asynq.Queue("low"),
)
```

## Alternatives Considered

### Alternative 1: RabbitMQ

**Pros**: Advanced routing, message acknowledgment, multiple protocols, clustering
**Cons**: Additional infrastructure, complex setup, overkill for our use case

**Rejected because**: Adds operational complexity; Redis-based solution meets our needs.

### Alternative 2: AWS SQS

**Pros**: Managed service, auto-scaling, high availability
**Cons**: Cloud lock-in, cost at scale, latency for Nigeria region, external dependency

**Rejected because**: Prefer self-hosted solution for cost control and latency.

### Alternative 3: PostgreSQL-based queue (pgq, que)

**Pros**: Uses existing database, transactional consistency
**Cons**: Polling overhead, limited throughput, database load

**Rejected because**: Don't want to add load to primary database.

### Alternative 4: Temporal/Cadence

**Pros**: Workflow orchestration, strong consistency, fault tolerance
**Cons**: Complex setup, significant learning curve, heavyweight

**Rejected because**: Over-engineered for our current requirements.

## References

- [Asynq GitHub Repository](https://github.com/hibiken/asynq)
- [Asynq Documentation](https://github.com/hibiken/asynq/wiki)
- [Asynqmon - Web UI](https://github.com/hibiken/asynqmon)
- [Building Reliable Background Jobs](https://brandur.org/job-drain)
