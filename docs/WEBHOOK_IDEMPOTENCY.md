# Webhook Idempotency Protection

**Implementation Date:** 2026-02-06
**Status:** Implemented
**Security Issue:** #7 (Medium Severity)

## Overview

This document describes the webhook idempotency protection mechanism implemented to prevent duplicate processing of payment webhooks from Paystack and other payment providers.

## Problem Statement

Payment providers may retry webhook deliveries for various reasons:
- Network failures
- Timeout errors
- Provider-side retry logic
- Load balancer retries

Without idempotency protection, duplicate webhooks could cause:
- **Double crediting** of user accounts (financial loss)
- **Data inconsistency** in transaction records
- **Audit trail corruption**

## Solution Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                    WEBHOOK REQUEST                           │
│  POST /api/webhook/paystack                                  │
│  X-Paystack-Signature: <hmac-sha512>                        │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│              WebhookHandler                                  │
│  1. Verify HMAC-SHA512 signature                            │
│  2. Extract payment reference                                │
│  3. Check idempotency (IsProcessed?)                        │
│  4. Mark as processed (MarkProcessed - atomic SetNX)        │
│  5. Process webhook event                                    │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│           WebhookEventStore (Redis)                          │
│  Key: webhook:event:{reference}                              │
│  Value: {eventID, provider, type, timestamp, payload}       │
│  TTL: 30 days (configurable)                                 │
└─────────────────────────────────────────────────────────────┘
```

### Domain Model

**Location:** `internal/domain/wallet/event/webhook.go`

```go
type WebhookEvent struct {
    EventID     WebhookEventID  // Unique ID (payment reference)
    Provider    string          // "paystack", "flutterwave"
    EventType   string          // "charge.success", etc.
    Reference   string          // Payment reference
    ProcessedAt time.Time       // When processed
    Payload     []byte          // Raw webhook for audit
}
```

### Repository Interface

**Location:** `internal/domain/wallet/repository/webhook_event_store.go`

```go
type WebhookEventStore interface {
    IsProcessed(ctx, eventID) (bool, error)
    MarkProcessed(ctx, webhookEvent) error
    GetEvent(ctx, eventID) (*WebhookEvent, error)
    CleanupExpired(ctx, retentionPeriod) error
}
```

### Redis Implementation

**Location:** `internal/infrastructure/persistence/redis/webhook_event_store.go`

Uses Redis `SETNX` (Set if Not eXists) for atomic idempotency checks:

```go
func (s *WebhookEventStore) MarkProcessed(ctx, webhookEvent) error {
    // SetNX provides atomicity - only succeeds if key doesn't exist
    success, err := s.client.SetNX(ctx, key, webhookEvent, 30*24*time.Hour)
    if !success {
        return ErrWebhookAlreadyProcessed
    }
    return nil
}
```

### HTTP Handler

**Location:** `internal/interface/http/handler/webhook_handler.go`

Processing flow:

1. **Signature Verification** (HMAC-SHA512)
2. **Extract Reference** from webhook payload
3. **Idempotency Check** (IsProcessed?)
   - If processed → Return 200 OK (acknowledge duplicate)
4. **Atomic Mark** (MarkProcessed via SetNX)
   - If race condition detected → Return 200 OK
5. **Process Event** (credit wallet, update status, etc.)
6. **Return 200 OK**

## Security Features

### 1. Signature Verification

All webhooks must include a valid HMAC-SHA512 signature:

```go
func verifyPaystackSignature(payload []byte, signature string) bool {
    mac := hmac.New(sha512.New, []byte(webhookSecret))
    mac.Write(payload)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))

    // Constant-time comparison prevents timing attacks
    return hmac.Equal([]byte(expectedSignature), []byte(signature))
}
```

### 2. Atomic Idempotency

Redis `SETNX` ensures only one instance processes each webhook:

```
Time    Instance A              Instance B
----    ---------------------   ---------------------
T0      Receive webhook ref123
T1      IsProcessed? → No
T2      MarkProcessed (SetNX)   Receive webhook ref123
T3      ✓ Success               IsProcessed? → Yes
T4      Process event           ✓ Return 200 (dup)
T5      ✓ Return 200
```

### 3. Audit Trail

All processed webhooks are stored with:
- Full payload (for forensics)
- Processing timestamp
- Event metadata

Retained for 30 days (configurable).

## Testing

### Unit Tests

**Location:** `internal/interface/http/handler/webhook_handler_test.go`

Tests cover:
- ✅ Signature verification (valid, invalid, missing)
- ✅ Idempotency (duplicate detection)
- ✅ Malformed payloads
- ✅ Different event types
- ✅ Race conditions (concurrent requests)

**Location:** `internal/infrastructure/persistence/redis/webhook_event_store_test.go`

Tests cover:
- ✅ Mark processed (new vs. duplicate)
- ✅ IsProcessed checks
- ✅ JSON serialization
- ✅ Concurrent processing

### Integration Tests

For production testing:

```bash
# Test duplicate webhook delivery
curl -X POST http://localhost:8081/api/webhook/paystack \
  -H "X-Paystack-Signature: <valid-signature>" \
  -d '{"event":"charge.success","data":{"reference":"test_ref","status":"success"}}'

# Send again - should be acknowledged but not reprocessed
curl -X POST http://localhost:8081/api/webhook/paystack \
  -H "X-Paystack-Signature: <valid-signature>" \
  -d '{"event":"charge.success","data":{"reference":"test_ref","status":"success"}}'
```

## Configuration

### Environment Variables

```bash
# Webhook secret from payment provider
PAYSTACK_WEBHOOK_SECRET=sk_test_xxxxxxxxxxxxx

# Redis configuration (for idempotency store)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Webhook event retention period (optional, default: 30 days)
WEBHOOK_RETENTION_DAYS=30
```

### Redis Requirements

- **Persistence:** Enable RDB or AOF for durability
- **Memory:** ~1KB per webhook event
- **Expiration:** Events auto-expire after retention period
- **High Availability:** Consider Redis Sentinel or Cluster for production

## Deployment Checklist

- [ ] Configure `PAYSTACK_WEBHOOK_SECRET` in environment
- [ ] Verify Redis is running and accessible
- [ ] Test webhook endpoint with Paystack test webhooks
- [ ] Monitor Redis memory usage
- [ ] Set up alerts for webhook processing failures
- [ ] Document webhook endpoint for payment provider

## Monitoring

### Metrics to Track

1. **Duplicate Rate:**
   ```
   duplicate_webhooks / total_webhooks
   ```
   Expected: 1-5% (depends on provider retry logic)

2. **Processing Time:**
   ```
   webhook_processing_duration_ms (p50, p95, p99)
   ```
   Expected: <50ms (idempotency check + Redis write)

3. **Error Rate:**
   ```
   failed_webhooks / total_webhooks
   ```
   Expected: <0.1%

### Log Messages

```
[WEBHOOK] Processing Paystack event: charge.success, reference: ref123
[WEBHOOK] Duplicate webhook detected: ref123 (already processed)
[WEBHOOK] Race condition detected: ref456 (processed by another instance)
[SECURITY] Webhook signature verification failed
```

## Error Handling

| Scenario | HTTP Status | Action |
|----------|-------------|--------|
| Invalid signature | 401 Unauthorized | Reject webhook |
| Malformed payload | 400 Bad Request | Reject webhook |
| Duplicate webhook | 200 OK | Acknowledge (don't reprocess) |
| Redis error | 500 Internal Error | Retry (provider will resend) |
| Processing error | 200 OK | Log error but acknowledge |

**Important:** Always return 200 OK for processed webhooks (even duplicates) to prevent unnecessary retries.

## Future Improvements

1. **Dead Letter Queue:** Store failed webhooks for manual retry
2. **Webhook Dashboard:** Admin UI to view/replay webhooks
3. **Multi-Provider Support:** Extend to Flutterwave, Stripe, etc.
4. **Webhook Versioning:** Handle API version changes
5. **Performance:** Batch processing for high volume

## References

- **Security Audit Report:** `docs/SECURITY_AUDIT_REPORT.md` (Issue #7)
- **PRD Section 10.3:** Next Recommended Tasks
- **Paystack Webhooks:** https://paystack.com/docs/payments/webhooks
- **Idempotency Best Practices:** https://stripe.com/docs/webhooks/best-practices

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2026-02-06 | 1.0 | Initial implementation |

---

**Implemented By:** Claude Sonnet 4.5
**Reviewed By:** _Pending_
**Approved By:** _Pending_
