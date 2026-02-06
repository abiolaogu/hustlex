# Webhook Idempotency Protection

**Status:** ✅ Implemented (2026-02-06)
**Security Issue:** #7 (Medium Priority)
**Version:** 1.0

---

## Overview

This document describes the webhook idempotency protection system implemented in HustleX to prevent duplicate processing of webhook events from payment providers (Paystack, Flutterwave, etc.).

## Problem Statement

Payment providers typically implement retry logic for webhook delivery to ensure reliability. However, this can lead to:

- **Double crediting** of user accounts
- **Financial loss** for the platform
- **Data inconsistency** in transaction records
- **Accounting discrepancies**

### Example Scenario (Without Idempotency)

```
Time    Event
-----   -----
10:00   Paystack sends webhook: "charge.success" for ₦10,000
10:00   Server processes webhook, credits user +₦10,000 (Balance: ₦10,000)
10:01   Network timeout, Paystack doesn't receive 200 OK
10:02   Paystack retries webhook: "charge.success" for ₦10,000
10:02   Server processes webhook AGAIN, credits user +₦10,000 (Balance: ₦20,000) ❌
```

**Result:** User receives ₦20,000 instead of ₦10,000. Platform loses ₦10,000.

---

## Solution: Idempotency Protection

We implement a Redis-backed idempotency tracking system that ensures each webhook event is processed exactly once, even if delivered multiple times.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Webhook Flow                            │
└─────────────────────────────────────────────────────────────┘

1. Paystack sends webhook
        ↓
2. Verify HMAC signature ✓
        ↓
3. Extract event ID (event_type:reference)
        ↓
4. Check Redis: IsProcessed(eventID)?
        ↓
   ┌────┴────┐
   │ YES     │ NO
   ↓         ↓
5a. Return  5b. MarkProcessed(eventID) [ATOMIC]
   200 OK       ↓
   "Already     ┌─────┴─────┐
   processed"   │ SUCCESS   │ FAILURE (race)
                ↓           ↓
             6. Process   Return 200 OK
                event     "Already processing"
                ↓
             7. Credit wallet / Update status
                ↓
             8. Return 200 OK
```

---

## Implementation Details

### 1. Domain Interface

**File:** `internal/domain/webhook/repository/event_store.go`

```go
type WebhookEventStore interface {
    // Check if event was already processed
    IsProcessed(ctx context.Context, eventID string) (bool, error)

    // Mark event as processed (atomic operation)
    MarkProcessed(ctx context.Context, eventID string, expiresIn time.Duration) error

    // Get when event was first processed (for auditing)
    GetProcessedAt(ctx context.Context, eventID string) (time.Time, bool, error)

    // Delete event record (for testing/cleanup)
    Delete(ctx context.Context, eventID string) error
}
```

### 2. Redis Implementation

**File:** `internal/infrastructure/webhook/redis_event_store.go`

#### Key Features:

- **Atomic Operations:** Uses Redis `SETNX` to ensure only one request marks an event as processed
- **Auto-Expiration:** Event records expire after 7 days to prevent unbounded growth
- **Metadata Storage:** Stores timestamp of first processing for audit trail
- **Race Condition Protection:** Concurrent requests for the same event handled gracefully

#### Redis Key Format:

```
webhook:event:{event_type}:{reference}
```

**Examples:**
- `webhook:event:charge.success:TRX_abc123`
- `webhook:event:transfer.failed:TRX_def456`

#### Storage Schema:

```json
{
  "event_id": "charge.success:TRX_abc123",
  "processed_at": "2026-02-06T14:30:00Z",
  "version": 1
}
```

### 3. Webhook Handler Integration

**File:** `internal/interface/http/handler/webhook_handler.go`

#### Processing Flow:

1. **Signature Verification** (security)
2. **Parse Event** (validation)
3. **Extract Event ID** (event_type + reference)
4. **Idempotency Check** (`IsProcessed`)
5. **Mark as Processed** (`MarkProcessed` - atomic)
6. **Process Event** (credit wallet, update status, etc.)
7. **Acknowledge** (return 200 OK)

#### Event ID Generation:

```go
// Format: "event_type:reference"
eventID := fmt.Sprintf("%s:%s", event.Event, data.Reference)
```

**Why include event type?**
The same `reference` might appear in different event types:
- `charge.success:TRX_abc123` (initial deposit)
- `transfer.failed:TRX_abc123` (failed withdrawal attempt for same transaction)

---

## Security Considerations

### 1. Signature Verification First

Always verify webhook signature BEFORE checking idempotency:

```go
// ✅ CORRECT ORDER
1. Verify HMAC signature
2. Check idempotency
3. Process event

// ❌ WRONG ORDER
1. Check idempotency
2. Process event
3. Verify signature (attacker can replay)
```

### 2. Constant-Time Comparison

Use `hmac.Equal()` for signature comparison to prevent timing attacks:

```go
// ✅ CORRECT
return hmac.Equal([]byte(expectedSignature), []byte(signature))

// ❌ WRONG (timing attack vulnerable)
return expectedSignature == signature
```

### 3. Event Retention Period

Events are retained for **7 days** (configurable):

```go
retentionPeriod := 7 * 24 * time.Hour
```

**Why 7 days?**
- Paystack typically retries for **24 hours**
- 7 days provides safety margin for extreme edge cases
- Prevents unbounded Redis memory growth

---

## Testing

### Unit Tests

**File:** `internal/infrastructure/webhook/redis_event_store_test.go`

Tests cover:
- ✅ Basic idempotency (mark, check, retrieve)
- ✅ Race conditions (concurrent marking)
- ✅ Expiration (TTL verification)
- ✅ Edge cases (empty IDs, zero expiration, etc.)
- ✅ Error handling

### Integration Tests

#### Test Case 1: First Processing

```bash
# Send webhook
curl -X POST http://localhost:8081/webhooks/paystack \
  -H "X-Paystack-Signature: <valid_signature>" \
  -d '{
    "event": "charge.success",
    "data": {
      "reference": "TRX_test_001",
      "amount": 1000000,
      "status": "success"
    }
  }'

# Expected: 200 OK, event processed, wallet credited
```

#### Test Case 2: Duplicate Processing

```bash
# Send SAME webhook again
curl -X POST http://localhost:8081/webhooks/paystack \
  -H "X-Paystack-Signature: <valid_signature>" \
  -d '{
    "event": "charge.success",
    "data": {
      "reference": "TRX_test_001",
      "amount": 1000000,
      "status": "success"
    }
  }'

# Expected: 200 OK, "Event already processed", wallet NOT credited again
```

#### Test Case 3: Concurrent Requests

```bash
# Send 10 identical webhooks simultaneously
for i in {1..10}; do
  curl -X POST http://localhost:8081/webhooks/paystack \
    -H "X-Paystack-Signature: <valid_signature>" \
    -d '{...}' &
done

# Expected: Only ONE processes successfully, 9 get "already processed/processing"
```

---

## Monitoring & Observability

### Metrics to Track

1. **Duplicate Event Rate**
   ```
   webhook_duplicate_events_total{provider="paystack", event_type="charge.success"}
   ```

2. **Processing Time**
   ```
   webhook_processing_duration_seconds{provider="paystack", event_type="charge.success"}
   ```

3. **Idempotency Check Failures**
   ```
   webhook_idempotency_errors_total{provider="paystack", error="redis_unavailable"}
   ```

### Log Messages

```
[WEBHOOK] Duplicate event charge.success:TRX_abc123 (originally processed at: 2026-02-06T14:30:00Z)
[WEBHOOK] Event charge.success:TRX_abc123 already being processed by another request
[WEBHOOK] Processing event: charge.success (ID: charge.success:TRX_abc123)
[WEBHOOK] Successfully processed charge event charge.success:TRX_abc123
```

---

## Edge Cases & Handling

### 1. Redis Unavailable

**Scenario:** Redis is down when webhook arrives.

**Handling:**
```go
isProcessed, err := h.eventStore.IsProcessed(ctx, eventID)
if err != nil {
    log.Printf("[WEBHOOK] Error checking idempotency: %v", err)
    // Return 200 to prevent retries - we can't verify
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Event acknowledged (idempotency check failed)",
    })
}
```

**Consequence:** Possible duplicate processing if Redis is down AND provider retries immediately. Acceptable trade-off to prevent blocking critical payments.

### 2. Event ID Extraction Fails

**Scenario:** Webhook missing `reference` field.

**Handling:**
```go
eventID, err := h.extractEventID(event)
if err != nil {
    log.Printf("[WEBHOOK] Failed to extract event ID: %v", err)
    return c.Status(http.StatusBadRequest).JSON(fiber.Map{
        "success": false,
        "error":   "Missing event identifier",
    })
}
```

**Consequence:** Return 400 Bad Request. Provider will NOT retry (400 = permanent failure).

### 3. Race Condition

**Scenario:** Two webhook requests arrive simultaneously.

**Handling:**
```go
err := h.eventStore.MarkProcessed(ctx, eventID, retentionPeriod)
if err == webhook.ErrEventAlreadyProcessed {
    // Another request beat us to it
    return c.JSON(fiber.Map{
        "success": true,
        "message": "Event already processing",
    })
}
```

**Consequence:** One succeeds, one gets "already processing". Both return 200 OK (no retries).

### 4. Event Expires from Redis

**Scenario:** Webhook replayed after 7+ days.

**Handling:** Event will be reprocessed as if new.

**Mitigation:**
- Payment providers don't retry for 7+ days
- Business logic should validate transaction dates
- Database constraints should prevent duplicate transactions by reference

---

## Configuration

### Environment Variables

```bash
# Redis configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secret
REDIS_DB=0
REDIS_POOL_SIZE=10

# Webhook configuration
PAYSTACK_WEBHOOK_SECRET=sk_test_xxxxxxxxxxxxx
WEBHOOK_EVENT_RETENTION_DAYS=7
```

### Initialization

```go
// Create Redis client
redisClient, err := redis.NewClient(redis.Config{
    Host:     os.Getenv("REDIS_HOST"),
    Port:     6379,
    Password: os.Getenv("REDIS_PASSWORD"),
    DB:       0,
    PoolSize: 10,
})

// Create event store
eventStore := webhook.NewRedisEventStore(redisClient)

// Create webhook handler
webhookHandler := handler.NewWebhookHandler(
    eventStore,
    os.Getenv("PAYSTACK_WEBHOOK_SECRET"),
)
```

---

## Performance Considerations

### Redis Memory Usage

Each event record: ~150 bytes (JSON metadata)

**Calculation:**
- 100,000 webhooks/day
- 7-day retention
- Total events: 700,000
- Memory: 700,000 × 150 bytes = **105 MB**

**Conclusion:** Negligible memory footprint.

### Latency Impact

Additional operations per webhook:
1. `EXISTS` check: ~1ms
2. `SETNX` operation: ~1ms

**Total overhead:** ~2ms per webhook (minimal)

---

## Compliance & Audit

### PCI DSS Requirement 10.3

✅ **Audit Trail:** `GetProcessedAt()` provides timestamp of first processing

### SOC 2 Type II

✅ **Change Tracking:** All webhook processing events logged with correlation IDs

### Financial Reconciliation

Query processed events for date range:
```go
// Get all processed events for a date
events := redis.Keys(ctx, "webhook:event:*")
// Filter by timestamp in eventRecord
```

---

## Troubleshooting

### Issue: Webhooks not processing

**Checklist:**
1. Check Redis connectivity: `redis-cli PING`
2. Verify webhook signature secret matches Paystack dashboard
3. Check logs for signature verification failures
4. Verify IP whitelist (if using)

### Issue: Events marked as duplicate incorrectly

**Investigation:**
```bash
# Check event in Redis
redis-cli GET "webhook:event:charge.success:TRX_abc123"

# Check TTL
redis-cli TTL "webhook:event:charge.success:TRX_abc123"
```

**Solution:**
```bash
# Manually delete event record (use with caution)
redis-cli DEL "webhook:event:charge.success:TRX_abc123"
```

### Issue: High duplicate event rate

**Possible causes:**
1. Provider retry logic too aggressive (normal)
2. Network issues causing timeouts
3. Application responding slowly (>30s)

**Solution:** Optimize webhook processing time to <5 seconds.

---

## Migration Guide

### Existing Webhook Implementation

If you have an existing webhook handler without idempotency:

**Step 1:** Add event store dependency

```go
type WebhookHandler struct {
    walletService *services.WalletService
+   eventStore    repository.WebhookEventStore
}
```

**Step 2:** Add idempotency checks

```go
func (h *WebhookHandler) HandlePaystackWebhook(c *fiber.Ctx) error {
    // ... existing signature verification ...

+   // Extract event ID
+   eventID := fmt.Sprintf("%s:%s", event.Event, data.Reference)
+
+   // Check if already processed
+   if isProcessed, _ := h.eventStore.IsProcessed(ctx, eventID); isProcessed {
+       return c.JSON(fiber.Map{"success": true, "message": "Already processed"})
+   }
+
+   // Mark as processed
+   if err := h.eventStore.MarkProcessed(ctx, eventID, 7*24*time.Hour); err != nil {
+       // Handle error...
+   }

    // ... existing event processing ...
}
```

**Step 3:** Test with duplicate webhooks

**Step 4:** Monitor for duplicate event logs

---

## Comparison with Alternatives

### Alternative 1: Database-Based Idempotency

**Pros:**
- Persistent storage
- SQL queries for reconciliation

**Cons:**
- Slower (50-100ms vs 1-2ms)
- Database load increases
- Requires database schema changes

### Alternative 2: Idempotency Keys in Webhook Payload

**Pros:**
- Provider-native solution

**Cons:**
- Not all providers support it
- Requires provider configuration
- Less control over retention

**Decision:** Redis-based solution provides best balance of speed, reliability, and control.

---

## Future Enhancements

### Phase 2 (Post-Launch)

1. **Persistent Backup:** Sync event records to PostgreSQL for long-term audit
2. **Webhook Replay:** Admin UI to replay failed webhooks
3. **Dead Letter Queue:** Capture events that fail processing for manual review
4. **Advanced Monitoring:** Grafana dashboard with duplicate rate, processing time, error rate

---

## References

- [Security Audit Report](./SECURITY_AUDIT_REPORT.md) - Issue #7
- [PRD Section 10.3](./PRD.md#103-next-recommended-tasks) - Task 4
- [Paystack Webhook Documentation](https://paystack.com/docs/payments/webhooks/)
- [Redis SETNX Documentation](https://redis.io/commands/setnx/)

---

## Changelog

- **v1.0 (2026-02-06):** Initial implementation
  - Redis-backed event store
  - Webhook handler integration
  - Comprehensive test suite
  - Documentation

---

**Document Owner:** Engineering Team
**Last Updated:** 2026-02-06
**Next Review:** 2026-03-06

---

*This feature resolves **Security Issue #7** from the Security Audit Report.*
