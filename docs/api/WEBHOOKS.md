# HustleX Webhook Documentation

This document describes how HustleX handles webhook notifications from external payment providers, specifically Paystack.

## Overview

Webhooks are HTTP callbacks that notify HustleX of events occurring in external systems. We currently integrate with:

- **Paystack** - Payment processing (deposits, withdrawals, transfers)

---

## Paystack Webhooks

### Endpoint

```
POST /v1/webhooks/paystack
```

### Security

All Paystack webhooks are verified using HMAC SHA512 signature validation:

1. Paystack sends the signature in the `X-Paystack-Signature` header
2. HustleX computes the expected signature using the secret key
3. Request is rejected if signatures don't match

**Signature Verification (Go):**

```go
func verifyPaystackSignature(payload []byte, signature string, secretKey string) bool {
    mac := hmac.New(sha512.New, []byte(secretKey))
    mac.Write(payload)
    expectedSignature := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

**Signature Verification (Node.js):**

```javascript
const crypto = require('crypto');

function verifyPaystackSignature(payload, signature, secretKey) {
  const expectedSignature = crypto
    .createHmac('sha512', secretKey)
    .update(payload)
    .digest('hex');
  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(expectedSignature)
  );
}
```

### Request Format

```http
POST /v1/webhooks/paystack HTTP/1.1
Host: api.hustlex.ng
Content-Type: application/json
X-Paystack-Signature: 74c1e6e0a6b8...

{
  "event": "charge.success",
  "data": {
    "id": 123456789,
    "domain": "live",
    "status": "success",
    "reference": "DEP-2024-abc123",
    "amount": 1000000,
    "currency": "NGN",
    "channel": "card",
    "customer": {
      "id": 12345,
      "email": "user@example.com"
    },
    "paid_at": "2024-01-15T10:35:00.000Z"
  }
}
```

### Response

HustleX always responds with `200 OK` to acknowledge receipt:

```json
{
  "status": "success",
  "message": "Webhook processed"
}
```

> **Note**: Returning non-2xx status codes causes Paystack to retry the webhook. Only return errors for temporary failures.

---

## Supported Paystack Events

### Deposits (Fund Wallet)

#### `charge.success`

Triggered when a card/bank payment succeeds.

**Payload:**
```json
{
  "event": "charge.success",
  "data": {
    "id": 123456789,
    "domain": "live",
    "status": "success",
    "reference": "DEP-2024-abc123",
    "amount": 1000000,
    "currency": "NGN",
    "channel": "card",
    "metadata": {
      "custom_fields": [
        {
          "display_name": "User ID",
          "variable_name": "user_id",
          "value": "usr_abc123"
        },
        {
          "display_name": "Transaction Type",
          "variable_name": "type",
          "value": "deposit"
        }
      ]
    },
    "authorization": {
      "authorization_code": "AUTH_xxx",
      "bin": "408408",
      "last4": "4081",
      "exp_month": "12",
      "exp_year": "2025",
      "channel": "card",
      "card_type": "visa",
      "bank": "GTBank",
      "country_code": "NG",
      "brand": "visa",
      "reusable": true,
      "signature": "SIG_xxx"
    },
    "customer": {
      "id": 12345,
      "first_name": "Adebayo",
      "last_name": "Okonkwo",
      "email": "user@example.com",
      "phone": "+2348012345678",
      "customer_code": "CUS_xxx"
    },
    "paid_at": "2024-01-15T10:35:00.000Z",
    "created_at": "2024-01-15T10:30:00.000Z"
  }
}
```

**HustleX Action:**
1. Verify reference matches pending deposit
2. Credit user's wallet with amount (minus fees if applicable)
3. Update transaction status to `completed`
4. Send push notification to user
5. Save card authorization if `saveCard` was requested

#### `charge.failed`

Triggered when a payment fails.

**Payload:**
```json
{
  "event": "charge.failed",
  "data": {
    "id": 123456789,
    "reference": "DEP-2024-abc123",
    "status": "failed",
    "gateway_response": "Declined",
    "message": "Card declined"
  }
}
```

**HustleX Action:**
1. Update transaction status to `failed`
2. Send notification with failure reason

---

### Withdrawals (Bank Transfer)

#### `transfer.success`

Triggered when a bank transfer succeeds.

**Payload:**
```json
{
  "event": "transfer.success",
  "data": {
    "amount": 5000000,
    "currency": "NGN",
    "domain": "live",
    "id": 123456789,
    "reference": "WTH-2024-xyz789",
    "source": "balance",
    "source_details": null,
    "reason": "Withdrawal to GTBank",
    "status": "success",
    "transfer_code": "TRF_xxx",
    "recipient": {
      "domain": "live",
      "type": "nuban",
      "currency": "NGN",
      "name": "ADEBAYO OKONKWO",
      "details": {
        "account_number": "0123456789",
        "account_name": "ADEBAYO OKONKWO",
        "bank_code": "058",
        "bank_name": "GTBank"
      }
    },
    "transferred_at": "2024-01-15T10:45:00.000Z"
  }
}
```

**HustleX Action:**
1. Update withdrawal status to `completed`
2. Send success notification to user
3. Update transaction record with transfer details

#### `transfer.failed`

Triggered when a bank transfer fails.

**Payload:**
```json
{
  "event": "transfer.failed",
  "data": {
    "id": 123456789,
    "reference": "WTH-2024-xyz789",
    "status": "failed",
    "reason": "Could not resolve bank account"
  }
}
```

**HustleX Action:**
1. Reverse the debit from user's wallet
2. Update withdrawal status to `failed`
3. Send notification with failure reason

#### `transfer.reversed`

Triggered when a transfer is reversed (e.g., account not found).

**Payload:**
```json
{
  "event": "transfer.reversed",
  "data": {
    "id": 123456789,
    "reference": "WTH-2024-xyz789",
    "status": "reversed",
    "reason": "Account not found"
  }
}
```

**HustleX Action:**
1. Credit back the withdrawal amount to wallet
2. Update withdrawal status to `reversed`
3. Send notification to user

---

### Recurring Payments

#### `subscription.create`

Triggered when a subscription is created.

**Payload:**
```json
{
  "event": "subscription.create",
  "data": {
    "domain": "live",
    "status": "active",
    "subscription_code": "SUB_xxx",
    "amount": 5000000,
    "cron_expression": "0 0 1 * *",
    "next_payment_date": "2024-02-01T00:00:00.000Z",
    "plan": {
      "name": "Monthly Savings",
      "plan_code": "PLN_xxx",
      "amount": 5000000,
      "interval": "monthly"
    },
    "customer": {
      "email": "user@example.com",
      "customer_code": "CUS_xxx"
    }
  }
}
```

#### `invoice.payment_failed`

Triggered when a recurring payment fails.

**HustleX Action:**
1. Mark contribution as missed
2. Notify user of failed payment
3. Retry logic if configured

---

## Event Processing

### Idempotency

All webhook events are processed idempotently using the event ID:

1. Store processed event IDs in database
2. Check if event already processed before handling
3. Skip duplicate events

```sql
CREATE TABLE processed_webhooks (
    event_id VARCHAR(255) PRIMARY KEY,
    event_type VARCHAR(100),
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Processing Flow

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Paystack  │────>│  API Gateway │────>│   Handler   │
└─────────────┘     └──────────────┘     └─────────────┘
                           │                    │
                    Verify Signature      Check Idempotency
                           │                    │
                           ▼                    ▼
                    ┌──────────────┐     ┌─────────────┐
                    │    Reject    │     │   Process   │
                    │  Invalid Sig │     │    Event    │
                    └──────────────┘     └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │   Update    │
                                        │  Database   │
                                        └─────────────┘
                                               │
                                               ▼
                                        ┌─────────────┐
                                        │    Send     │
                                        │Notification │
                                        └─────────────┘
```

### Retry Handling

Paystack retries failed webhooks with exponential backoff:

| Attempt | Delay |
|---------|-------|
| 1 | Immediate |
| 2 | 10 minutes |
| 3 | 30 minutes |
| 4 | 1 hour |
| 5 | 3 hours |
| 6 | 6 hours |
| 7 | 12 hours |
| 8 | 24 hours |

After 8 failed attempts, the webhook is marked as failed.

---

## Reference Naming Conventions

| Type | Pattern | Example |
|------|---------|---------|
| Deposit | `DEP-{YYYY}-{random}` | `DEP-2024-abc123` |
| Withdrawal | `WTH-{YYYY}-{random}` | `WTH-2024-xyz789` |
| Transfer | `TRF-{YYYY}-{random}` | `TRF-2024-def456` |
| Savings | `SAV-{circleId}-{random}` | `SAV-circle123-ghi789` |
| Loan Repayment | `LRP-{loanId}-{random}` | `LRP-loan456-jkl012` |

---

## Testing Webhooks

### Local Development

Use ngrok to expose local server:

```bash
ngrok http 8080
```

Configure the ngrok URL in Paystack dashboard.

### Test Events

Paystack provides test mode with simulated events. Use test API keys and card numbers:

| Card | Behavior |
|------|----------|
| 4084 0840 8408 4081 | Success |
| 4084 0840 8408 4099 | Declined |
| 5060 6666 6666 6666 | PIN required |

### Webhook Simulator

Simulate webhooks locally:

```bash
curl -X POST http://localhost:8080/v1/webhooks/paystack \
  -H "Content-Type: application/json" \
  -H "X-Paystack-Signature: $(echo -n '{"event":"charge.success"...}' | openssl dgst -sha512 -hmac 'your_secret_key')" \
  -d '{"event":"charge.success","data":{...}}'
```

---

## Monitoring

### Metrics

Track these webhook metrics:

- `webhook_received_total` - Total webhooks received
- `webhook_processed_total` - Successfully processed
- `webhook_failed_total` - Processing failures
- `webhook_duplicate_total` - Duplicate events
- `webhook_invalid_signature_total` - Invalid signatures
- `webhook_processing_duration` - Processing time

### Alerts

Configure alerts for:

- High failure rate (>5% in 5 minutes)
- Processing latency (>5 seconds)
- Invalid signature spike
- Missing expected events

---

## Troubleshooting

### Common Issues

**Signature Verification Failed**
- Ensure using raw request body (not parsed JSON)
- Check secret key is correct
- Verify no middleware modifying body

**Duplicate Processing**
- Implement idempotency using event ID
- Store processed events in database

**Missing Events**
- Check Paystack dashboard for delivery logs
- Verify webhook URL is accessible
- Check firewall/security rules

### Debug Mode

Enable webhook debug logging:

```go
// config/webhooks.go
type WebhookConfig struct {
    DebugMode bool `env:"WEBHOOK_DEBUG" default:"false"`
}
```

---

## Support

- Paystack Documentation: https://paystack.com/docs/webhooks
- HustleX API Support: api-support@hustlex.ng
- Status Page: https://status.hustlex.ng
