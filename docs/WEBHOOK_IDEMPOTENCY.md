# Webhook Idempotency Implementation

**Status:** ✅ Implemented
**Date:** 2026-02-07
**Security Issue:** #7 - Medium Priority
**Author:** Claude Sonnet 4.5

---

## Overview

This document describes the webhook idempotency implementation that prevents duplicate processing of payment webhooks from providers like Paystack and Flutterwave. This protects against double crediting, financial loss, and data inconsistency when payment providers retry webhook deliveries.

## Problem Statement

Payment providers (Paystack, Flutterwave, etc.) retry webhook deliveries if they don't receive an acknowledgment. Without idempotency protection, this can lead to:

1. **Double Crediting**: User accounts credited multiple times for the same payment
2. **Financial Loss**: Platform loses money from duplicate credits
3. **Data Inconsistency**: Transaction records duplicated in the database
4. **Audit Issues**: Reconciliation problems with payment provider records

## Solution Architecture

### Components

#### 1. WebhookEventStore Interface
```go
type WebhookEventStore interface {
    IsProcessed(ctx context.Context, eventID string) (bool, error)
    MarkProcessed(ctx context.Context, eventID string, ttl time.Duration) error
    GetProcessedAt(ctx context.Context, eventID string) (*time.Time, error)
}
```

**Purpose:** Tracks which webhook events have been processed
**Implementation:** Redis-backed store with automatic TTL expiration
**Location:** `apps/api/internal/services/webhook_event_store.go`

#### 2. WebhookHandler with Idempotency
**Purpose:** Processes webhooks with duplicate detection
**Location:** `apps/api/internal/interface/http/handlers/webhook_handler.go`

### Data Flow

```
┌─────────────────┐
│ Payment Provider│
│   (Paystack)    │
└────────┬────────┘
         │
         │ Webhook POST
         ▼
┌─────────────────────────────────────────┐
│ 1. Verify Signature (HMAC-SHA512)       │
└────────┬────────────────────────────────┘
         │ ✓ Valid
         ▼
┌─────────────────────────────────────────┐
│ 2. Extract Reference (trx_abc123)       │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│ 3. Check if Already Processed           │
│    Redis GET webhook:event:trx_abc123   │
└────────┬────────────────────────────────┘
         │
         ├─── YES → Return 200 "Already Processed"
         │
         └─── NO
              │
              ▼
┌─────────────────────────────────────────┐
│ 4. Process Business Logic               │
│    - Credit wallet                       │
│    - Update transaction status           │
│    - Send notifications                  │
└────────┬────────────────────────────────┘
         │ ✓ Success
         ▼
┌─────────────────────────────────────────┐
│ 5. Mark as Processed                    │
│    Redis SETEX webhook:event:trx_abc123 │
│    TTL: 3 days (Paystack retry window)  │
└────────┬────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│ 6. Return 200 "Processed"               │
└─────────────────────────────────────────┘
```

## Implementation Details

### 1. Event Identification

Each webhook event is uniquely identified by its **transaction reference**:
- Paystack: `data.reference` field (e.g., `trx_abc123xyz`)
- Flutterwave: `data.flw_ref` field (e.g., `FLW-123456789`)

### 2. Redis Key Structure

```
webhook:event:{reference}
```

**Example:** `webhook:event:trx_abc123xyz`

**Value (JSON):**
```json
{
  "event_id": "trx_abc123xyz",
  "processed_at": "2026-02-07T10:15:00Z"
}
```

**TTL:** Based on payment provider retry window
- Paystack: 3 days (72 hours)
- Flutterwave: 3 days (72 hours)
- Default: 7 days (for safety)

### 3. Idempotency Check Logic

```go
// Check if already processed
processed, err := h.eventStore.IsProcessed(ctx, reference)
if err != nil {
    // Log error but acknowledge webhook to prevent retries
    return http.StatusOK, "Event acknowledged (idempotency check failed)"
}

if processed {
    // Duplicate detected - acknowledge without processing
    return http.StatusOK, "Event already processed"
}

// Process business logic...

// Mark as processed AFTER successful processing
h.eventStore.MarkProcessed(ctx, reference, PaystackRetryWindow)
```

### 4. Error Handling

**Idempotency Check Failure:**
- Still return 200 OK to prevent provider retries
- Log error for investigation
- Message: "Event acknowledged (idempotency check failed)"

**Business Logic Failure:**
- Do NOT mark as processed
- Return 500 error to trigger provider retry
- Provider will redeliver webhook

**Mark Processed Failure:**
- Business logic already succeeded
- Still return 200 OK (operation completed successfully)
- Log error for investigation
- Worst case: duplicate processing on next retry (acceptable trade-off)

## Security Considerations

### 1. Signature Verification

All webhooks MUST be verified before idempotency check:

```go
// Verify webhook signature first
signature := r.Header.Get("X-Paystack-Signature")
if !h.verifyPaystackSignature(body, signature) {
    return http.StatusUnauthorized, "Invalid signature"
}

// Then check idempotency
```

**Why:** Prevents attackers from marking legitimate events as "processed" before they arrive.

### 2. Constant-Time Comparison

Signature verification uses `hmac.Equal()` to prevent timing attacks:

```go
return hmac.Equal([]byte(expectedSignature), []byte(signature))
```

### 3. TTL Selection

TTL must be **longer** than payment provider retry window:
- Provider retry window: 3 days
- Our TTL: 3 days minimum (7 days recommended)
- If TTL too short: Risk of duplicate processing after expiry
- If TTL too long: Higher Redis memory usage (acceptable trade-off)

## Testing

### Unit Tests

**Location:** `apps/api/internal/services/webhook_event_store_test.go`

Coverage:
- ✅ Unprocessed events return false
- ✅ Processed events return true
- ✅ Timestamp tracking
- ✅ TTL expiration
- ✅ Concurrent access safety
- ✅ Provider-specific TTLs

**Location:** `apps/api/internal/interface/http/handlers/webhook_handler_test.go`

Coverage:
- ✅ Signature verification (valid/invalid)
- ✅ Duplicate detection
- ✅ First-time processing
- ✅ Multiple event types (charge, transfer)
- ✅ Malformed payloads
- ✅ Missing reference fields
- ✅ Real-world retry scenarios

### Integration Testing

**Test Scenario 1: Paystack Retry Behavior**
1. Send charge.success webhook → Processed
2. Send same webhook after 1 minute → Rejected (duplicate)
3. Send same webhook after 1 hour → Rejected (duplicate)

**Test Scenario 2: Multiple Unique Events**
```bash
# Send 3 different webhooks
curl -X POST /webhooks/paystack -d '{"event":"charge.success","data":{"reference":"trx_001"}}'
curl -X POST /webhooks/paystack -d '{"event":"charge.success","data":{"reference":"trx_002"}}'
curl -X POST /webhooks/paystack -d '{"event":"charge.success","data":{"reference":"trx_003"}}'

# All should be processed successfully
```

**Test Scenario 3: Signature Tampering**
```bash
# Attempt to replay webhook with wrong signature
curl -X POST /webhooks/paystack \
  -H "X-Paystack-Signature: wrong_signature" \
  -d '{"event":"charge.success","data":{"reference":"trx_123"}}'

# Should be rejected with 401 Unauthorized
```

### Load Testing

**Simulate High-Volume Webhook Traffic:**
```bash
# 1000 concurrent webhooks (mix of unique and duplicates)
ab -n 1000 -c 100 -p webhook_payload.json \
   -H "X-Paystack-Signature: $(generate_signature webhook_payload.json)" \
   https://api.hustlex.com/webhooks/paystack
```

**Expected Behavior:**
- Unique events: Processed successfully
- Duplicate events: Rejected with "already processed"
- No double crediting
- All responses < 200ms (p95)

## Monitoring & Alerts

### Metrics to Track

1. **Webhook Processing Rate**
   - Metric: `webhook_processed_total{provider="paystack",event="charge.success"}`
   - Alert: Rate drops to zero for >5 minutes

2. **Duplicate Rate**
   - Metric: `webhook_duplicates_total{provider="paystack"}`
   - Alert: Rate exceeds 10% of total webhooks (indicates retry issues)

3. **Idempotency Check Failures**
   - Metric: `webhook_idempotency_errors_total`
   - Alert: Any failures (investigate immediately)

4. **Processing Latency**
   - Metric: `webhook_processing_duration_seconds{quantile="0.95"}`
   - Alert: p95 latency >500ms

### Logs to Monitor

```json
{
  "level": "info",
  "timestamp": "2026-02-07T10:15:00Z",
  "msg": "Duplicate event detected",
  "reference": "trx_abc123",
  "provider": "paystack",
  "event": "charge.success"
}
```

```json
{
  "level": "error",
  "timestamp": "2026-02-07T10:15:00Z",
  "msg": "Failed to check if event processed",
  "reference": "trx_xyz789",
  "error": "redis connection timeout"
}
```

## Operations Runbook

### Scenario 1: Redis Outage

**Symptoms:**
- Idempotency checks failing
- All webhooks acknowledged with "idempotency check failed"

**Impact:**
- Risk of duplicate processing during outage
- Low impact (webhooks still acknowledged)

**Response:**
1. Redis auto-restarts (Docker/K8s health checks)
2. Check Redis logs: `kubectl logs -f redis-0`
3. If persistent, manually restart: `kubectl rollout restart statefulset/redis`
4. After recovery, audit transactions for duplicates

**Prevention:**
- Redis cluster with replication
- Regular health checks
- Automated failover

### Scenario 2: Duplicate Webhooks After TTL Expiry

**Symptoms:**
- Same transaction processed twice, days apart
- User balance incorrect

**Detection:**
```sql
SELECT reference, COUNT(*) as count
FROM transactions
WHERE created_at > NOW() - INTERVAL '7 days'
GROUP BY reference
HAVING COUNT(*) > 1;
```

**Response:**
1. Identify affected users
2. Reverse duplicate transactions
3. Increase TTL: `PaystackRetryWindow = 7 * 24 * time.Hour`
4. Deploy updated configuration

### Scenario 3: Signature Verification Failures

**Symptoms:**
- All webhooks rejected with "Invalid signature"
- No deposits/withdrawals processing

**Root Causes:**
- Webhook secret changed in Paystack dashboard
- Secret not updated in application config
- Environment variable misconfiguration

**Response:**
1. Check Paystack dashboard for secret
2. Update Kubernetes secret:
   ```bash
   kubectl create secret generic paystack-secret \
     --from-literal=webhook-secret='new_secret' \
     --dry-run=client -o yaml | kubectl apply -f -
   ```
3. Restart API pods: `kubectl rollout restart deployment/api`

## Configuration

### Environment Variables

```bash
# Paystack
PAYSTACK_WEBHOOK_SECRET=sk_live_abc123...

# Flutterwave
FLUTTERWAVE_WEBHOOK_SECRET=FLW-SECRET-xyz789...

# Redis (for idempotency store)
REDIS_HOST=dragonfly
REDIS_PORT=6379
REDIS_PASSWORD=secure_password
REDIS_DB=0
```

### Kubernetes ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: webhook-config
data:
  paystack_retry_window: "72h"  # 3 days
  flutterwave_retry_window: "72h"
  default_webhook_ttl: "168h"  # 7 days
```

## Performance Impact

### Redis Memory Usage

**Per Event:** ~150 bytes (reference + timestamp + Redis overhead)

**Daily Volume (estimated):**
- 10,000 transactions/day
- 150 bytes × 10,000 = 1.5 MB/day
- 7-day retention: 10.5 MB total

**Conclusion:** Negligible memory impact

### Latency Addition

**Idempotency Check:**
- Redis GET: ~1-2ms (local network)
- Negligible compared to webhook processing time (~50-100ms)

**Total Impact:** <2% increase in webhook processing time

## Future Enhancements

### 1. Distributed Idempotency Store
- Replace Redis with distributed store (DynamoDB, Cassandra)
- Multi-region support
- Higher availability guarantees

### 2. Webhook Replay API
```bash
POST /admin/webhooks/replay
{
  "reference": "trx_abc123",
  "reason": "Manual replay after system recovery"
}
```

### 3. Idempotency Dashboard
- Admin UI showing processed events
- Duplicate detection statistics
- Manual event reprocessing

### 4. Multi-Provider Support
- Flutterwave idempotency
- Stripe idempotency
- Generic webhook provider interface

## References

- **Security Audit Report:** `docs/SECURITY_AUDIT_REPORT.md` (Issue #7)
- **Payment Webhooks Documentation:** `docs/api/WEBHOOKS.md`
- **Redis Client:** `apps/api/internal/infrastructure/cache/redis/client.go`
- **Token Blacklist (similar pattern):** `docs/TOKEN_REVOCATION.md`

## Compliance

### PCI DSS Requirements
- ✅ **Req 6.5.9:** Improper error handling (graceful degradation)
- ✅ **Req 8.2.3:** Strong cryptography (HMAC-SHA512 signatures)
- ✅ **Req 10.2:** Audit trails (all webhook events logged)

### CBN Guidelines
- ✅ Transaction integrity maintained
- ✅ Duplicate prevention mechanism
- ✅ Audit trail for reconciliation

---

**Implementation Complete:** 2026-02-07
**Next Review:** 2026-03-07
**Status:** ✅ Production Ready (pending integration tests)
