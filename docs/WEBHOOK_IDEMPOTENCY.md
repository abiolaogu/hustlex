# Webhook Idempotency Implementation

**Date:** 2026-02-07
**Security Issue:** #7 (MEDIUM Priority)
**Status:** ✅ IMPLEMENTED

---

## Overview

This document describes the webhook idempotency mechanism implemented for HustleX to prevent duplicate processing of payment webhook events from Paystack and other payment providers.

### Problem Statement

Payment providers like Paystack may retry webhook delivery multiple times due to:
- Network timeouts
- Temporary service unavailability
- HTTP 5xx errors from our endpoint
- Delayed acknowledgment (200 OK) responses

Without idempotency protection, duplicate webhook events can cause:
- **Double-crediting** of user wallets (financial loss)
- **Data inconsistencies** in transaction records
- **Duplicate notifications** to users
- **Incorrect audit trails**

---

## Solution Architecture

### Components

1. **WebhookEventStore Interface** (`internal/infrastructure/webhook_idempotency.go`)
   - Tracks processed webhook events
   - Provides idempotency guarantees

2. **RedisWebhookEventStore** (Production)
   - Redis-backed implementation
   - Distributed idempotency across multiple app instances
   - Automatic TTL-based cleanup (7 days)

3. **InMemoryWebhookEventStore** (Testing/Fallback)
   - In-memory implementation for tests
   - Fallback for single-instance deployments

4. **Enhanced Webhook Handler** (`internal/interface/http/handler/webhook_handler.go`)
   - Integrates idempotency checks
   - Handles duplicate detection gracefully

### Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                    Paystack Webhook Delivery                    │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
           ┌──────────────────────────────┐
           │  1. Verify Signature (HMAC)  │
           │     (Security Layer)         │
           └──────────┬───────────────────┘
                      │
                      ▼
           ┌──────────────────────────────┐
           │  2. Parse Webhook Event      │
           │     Extract Reference ID     │
           └──────────┬───────────────────┘
                      │
                      ▼
           ┌──────────────────────────────┐
           │  3. Check if Processed       │
           │     eventStore.IsProcessed() │
           └──────────┬───────────────────┘
                      │
          ┌───────────┴───────────┐
          │                       │
      Already              First Time
      Processed            Receiving
          │                       │
          ▼                       ▼
  ┌───────────────┐   ┌────────────────────────┐
  │ Return 200 OK │   │  4. Process Event      │
  │ Skip Business │   │     - Credit wallet    │
  │ Logic         │   │     - Update status    │
  │               │   │     - Send notification│
  └───────────────┘   └───────────┬────────────┘
                                  │
                                  ▼
                      ┌────────────────────────┐
                      │  5. Mark as Processed  │
                      │     eventStore.Mark()  │
                      │     TTL: 7 days        │
                      └───────────┬────────────┘
                                  │
                                  ▼
                      ┌────────────────────────┐
                      │  6. Return 200 OK      │
                      └────────────────────────┘
```

---

## Implementation Details

### 1. Idempotency Key Generation

Each webhook event is assigned a unique idempotency key:

```go
// For charge.success events
idempotencyKey = "charge:{paystack_id}:{reference}"
// Example: "charge:123456:trx_abc123"

// For transfer events
idempotencyKey = "transfer:{paystack_id}:{reference}"
// Example: "transfer:789012:trx_xyz456"

// For unknown events
idempotencyKey = "{event_type}:{sha512_hash}"
// Example: "custom.event:a1b2c3d4..."
```

**Why this format?**
- **Paystack ID**: Unique identifier from payment provider (prevents conflicts)
- **Reference**: Our internal transaction reference (additional uniqueness)
- **Event Type**: Distinguishes between charge/transfer/other events

### 2. Redis Storage

Events are stored in Redis with a 7-day TTL:

```
Key:   webhook:event:charge:123456:trx_abc123
Value: {
  "event_id": "charge:123456:trx_abc123",
  "processed_at": "2026-02-07T12:00:00Z"
}
TTL:   7 days (604,800 seconds)
```

**Why 7 days?**
- Paystack retries for up to 3 days (worst case)
- Extra buffer for edge cases and debugging
- Balances memory usage vs. safety

### 3. Race Condition Handling

Redis `SET NX` (Set if Not Exists) provides atomic check-and-set:

```go
// Pseudo-code for atomic idempotency
success := redis.SetNX(key, value, ttl)
if !success {
    // Key already exists - duplicate detected
    return "Event already processed"
}
// Key was set - proceed with processing
```

**Note**: Current implementation uses separate `Exists()` + `Set()` calls. For production under high load, consider upgrading to atomic `SetNX()`.

---

## API Usage

### Basic Usage

```go
import (
    "hustlex/internal/infrastructure"
    "hustlex/internal/infrastructure/cache/redis"
)

// Setup (in main.go or initialization)
redisClient, _ := redis.NewClient(redisConfig)
eventStore := infrastructure.NewRedisWebhookEventStore(redisClient)

handler := handler.NewWebhookHandler(eventStore, config)

// Register route
app.Post("/webhooks/paystack", handler.HandlePaystackWebhook)
```

### Manual Idempotency Check (if needed)

```go
ctx := context.Background()
eventID := "charge:123456:trx_abc123"

// Check if already processed
isProcessed, err := eventStore.IsProcessed(ctx, eventID)
if err != nil {
    log.Printf("Failed to check idempotency: %v", err)
}

if isProcessed {
    log.Printf("Duplicate event detected: %s", eventID)
    return
}

// Process event...
processWebhookEvent(event)

// Mark as processed (7 day TTL)
err = eventStore.MarkProcessed(ctx, eventID, 7*24*time.Hour)
if err != nil {
    log.Printf("Failed to mark event as processed: %v", err)
}
```

---

## Testing

### Unit Tests

Run tests for idempotency store:

```bash
cd apps/api
go test -v ./internal/infrastructure -run TestInMemoryWebhookEventStore
```

### Integration Tests

Run webhook handler tests:

```bash
go test -v ./internal/interface/http/handler -run TestWebhookHandler
```

### Test Coverage

```bash
go test -cover ./internal/infrastructure
go test -cover ./internal/interface/http/handler
```

Expected coverage: >90%

### Manual Testing with curl

```bash
# Generate signature
SECRET="your_webhook_secret"
PAYLOAD='{"event":"charge.success","data":{"id":123,"reference":"test_ref","amount":5000000,"status":"success"}}'
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha512 -hmac "$SECRET" | awk '{print $2}')

# Send first webhook (should process)
curl -X POST http://localhost:8081/webhooks/paystack \
  -H "Content-Type: application/json" \
  -H "X-Paystack-Signature: $SIGNATURE" \
  -d "$PAYLOAD"

# Send duplicate webhook (should skip)
curl -X POST http://localhost:8081/webhooks/paystack \
  -H "Content-Type: application/json" \
  -H "X-Paystack-Signature: $SIGNATURE" \
  -d "$PAYLOAD"
```

Expected output:
1. First request: `{"success":true,"message":"Webhook processed"}`
2. Second request: `{"success":true,"message":"Event already processed"}`

---

## Security Considerations

### 1. Signature Verification (Always First)

Idempotency checks happen **after** signature verification:

```go
// ✅ CORRECT ORDER
1. Verify HMAC signature
2. Parse event data
3. Check idempotency
4. Process event
```

This prevents attackers from:
- Replaying captured webhook payloads
- Polluting the idempotency store with fake events

### 2. Timing Attack Prevention

Signature comparison uses constant-time equality:

```go
hmac.Equal([]byte(expected), []byte(received))
```

### 3. IP Whitelisting (Optional)

Additional layer of defense:

```go
app.Use("/webhooks/paystack",
    handler.IPWhitelistMiddleware(handler.PaystackWebhookIPs()))
```

Paystack webhook IPs (as of 2026):
- 52.31.139.75
- 52.49.173.169
- 52.214.14.220

---

## Monitoring & Observability

### Metrics to Track

1. **Duplicate Detection Rate**
   ```
   webhook_duplicate_count{event_type="charge.success"}
   ```

2. **Idempotency Store Errors**
   ```
   webhook_idempotency_error_count{operation="check"}
   ```

3. **Processing Time**
   ```
   webhook_processing_duration_seconds{event_type="charge.success"}
   ```

### Log Messages

```
[WEBHOOK] Received Paystack event: charge.success
[WEBHOOK] Processing new event: charge:123456:trx_abc123
[WEBHOOK] Successfully processed event: charge:123456:trx_abc123
[WEBHOOK] Duplicate webhook detected: charge:123456:trx_abc123 - acknowledging
[WEBHOOK] Failed to check idempotency for charge:123456:trx_abc123: redis timeout
```

### Alerting

Set up alerts for:
- High duplicate detection rate (>50% of events)
- Idempotency store connection failures
- Webhook processing failures

---

## Performance Considerations

### Redis Performance

- **Latency**: ~1-2ms for `EXISTS` + `SET` operations
- **Throughput**: 10,000+ operations/second
- **Memory**: ~100 bytes per stored event

### Capacity Planning

Estimate for 1 million webhooks/month:
- Events stored: ~300,000 (within 7-day window)
- Memory usage: ~30 MB
- Redis queries: 2M reads + 1M writes per month

### Scalability

- Redis supports horizontal scaling (Redis Cluster)
- Stateless webhook handlers (scale horizontally)
- No single point of failure with Redis Sentinel

---

## Troubleshooting

### Issue: Duplicate Events Still Processed

**Symptoms**: User wallet credited twice for same transaction

**Possible Causes**:
1. Redis connection failure (falling back to in-memory store across instances)
2. Race condition (two requests processed simultaneously before marking)
3. TTL expired (event older than 7 days)

**Solutions**:
1. Check Redis connectivity and logs
2. Implement atomic `SetNX` for high-traffic scenarios
3. Increase TTL if Paystack retries are delayed

### Issue: Events Rejected as Duplicates (False Positives)

**Symptoms**: Legitimate webhook events skipped

**Possible Causes**:
1. Event ID generation collision
2. Redis key not expiring properly
3. Manual testing with same event ID

**Solutions**:
1. Review idempotency key format (ensure uniqueness)
2. Verify Redis TTL configuration
3. Use unique transaction references in tests

### Issue: High Idempotency Store Errors

**Symptoms**: Logs show frequent Redis connection errors

**Possible Causes**:
1. Redis server overloaded or down
2. Network connectivity issues
3. Redis memory limit reached

**Solutions**:
1. Scale Redis (vertical or horizontal)
2. Check network and firewall rules
3. Increase Redis `maxmemory` or enable eviction policy

---

## Migration Guide

### Migrating from Deprecated Handler

If upgrading from `internal/handlers_deprecated/webhook_handler.go`:

1. **Update imports**:
   ```go
   // Old
   import "hustlex/internal/handlers"

   // New
   import "hustlex/internal/interface/http/handler"
   import "hustlex/internal/infrastructure"
   ```

2. **Initialize event store**:
   ```go
   redisClient, _ := redis.NewClient(config.Redis)
   eventStore := infrastructure.NewRedisWebhookEventStore(redisClient)
   ```

3. **Create new handler**:
   ```go
   // Old
   webhookHandler := handlers.NewWebhookHandler(walletService, config)

   // New
   webhookHandler := handler.NewWebhookHandler(eventStore, config)
   ```

4. **No route changes needed** - same endpoint path

---

## Future Enhancements

### 1. Atomic SetNX Operation

Replace separate `Exists()` + `Set()` with atomic `SetNX()`:

```go
success, err := redisClient.SetNX(ctx, key, value, ttl)
if err != nil {
    return err
}
if !success {
    return ErrAlreadyProcessed
}
```

### 2. Event Replay Capability

Add admin endpoint to replay events (for recovery):

```go
POST /admin/webhooks/replay
{
  "event_id": "charge:123456:trx_abc123",
  "force": true
}
```

### 3. Multi-Provider Support

Extend to support other payment providers:
- Flutterwave webhooks
- Interswitch webhooks
- Stripe webhooks (for diaspora services)

### 4. Webhook Queue

Add RabbitMQ queue for async processing:
```
Webhook Endpoint → Queue → Worker Pool → Process Event
```

Benefits:
- Better fault tolerance
- Rate limiting protection
- Backpressure handling

---

## Compliance & Audit

### PCI DSS Compliance

✅ **Requirement 6.5.10**: Protection against replay attacks
- Idempotency prevents duplicate financial transactions

### NDPR Compliance

✅ **Data Protection**: Webhook events contain minimal PII
- Only user IDs stored in idempotency tracking
- 7-day TTL ensures timely data deletion

### Audit Trail

All webhook events logged with:
- Event type
- Event ID
- Processing status (new/duplicate)
- Timestamp
- Source IP (for security analysis)

---

## References

- **Paystack Webhooks Documentation**: https://paystack.com/docs/payments/webhooks/
- **RFC 7234 (HTTP Caching)**: Idempotency concepts
- **OWASP Webhook Security**: https://cheatsheetseries.owasp.org/cheatsheets/Webhook_Security_Cheat_Sheet.html
- **HustleX Security Audit Report**: `docs/SECURITY_AUDIT_REPORT.md` (Issue #7)

---

## Changelog

| Date       | Version | Changes                                    | Author           |
|------------|---------|--------------------------------------------|------------------|
| 2026-02-07 | 1.0     | Initial implementation of webhook idempotency | Claude Sonnet 4.5 |

---

**Document Owner**: Backend Team
**Reviewers**: Security Team, DevOps
**Next Review**: 2026-03-07 (1 month post-deployment)

---

*This implementation resolves Security Audit Issue #7 (MEDIUM Priority) from the HustleX Security Audit Report.*
