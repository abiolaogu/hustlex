# HustleX Security Documentation

> Comprehensive Security Guide for HustleX Platform

---

## Table of Contents

1. [Security Overview](#1-security-overview)
2. [Authentication & Authorization](#2-authentication--authorization)
3. [Data Protection](#3-data-protection)
4. [Application Security](#4-application-security)
5. [Infrastructure Security](#5-infrastructure-security)
6. [Compliance & Regulations](#6-compliance--regulations)
7. [Incident Response](#7-incident-response)
8. [Security Best Practices](#8-security-best-practices)
9. [Security Checklist](#9-security-checklist)

---

## 1. Security Overview

### 1.1 Security Principles

HustleX follows these core security principles:

| Principle | Implementation |
|-----------|----------------|
| **Defense in Depth** | Multiple security layers |
| **Least Privilege** | Minimal access by default |
| **Secure by Default** | Security built into design |
| **Fail Secure** | Secure state on errors |
| **Audit Everything** | Complete activity logging |

### 1.2 Security Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Security Layers                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Application Layer                     │    │
│  │  • Input validation   • SQL injection prevention        │    │
│  │  • XSS prevention     • CSRF protection                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  Authentication Layer                    │    │
│  │  • OTP verification   • JWT tokens                      │    │
│  │  • Transaction PIN    • Biometric auth                  │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Transport Layer                       │    │
│  │  • TLS 1.3           • Certificate pinning              │    │
│  │  • HSTS              • Secure headers                   │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  Infrastructure Layer                    │    │
│  │  • Firewall          • Network isolation                │    │
│  │  • WAF               • DDoS protection                  │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                      Data Layer                          │    │
│  │  • Encryption at rest • Encryption in transit           │    │
│  │  • Key management     • Backup encryption               │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### 1.3 Threat Model

**Primary Assets:**
- User financial data (wallet balances, transactions)
- Personal information (PII)
- Authentication credentials
- Payment information

**Primary Threats:**
| Threat | Risk Level | Mitigation |
|--------|------------|------------|
| Account takeover | High | MFA, OTP, device tracking |
| Payment fraud | High | PIN, escrow, limits |
| Data breach | High | Encryption, access control |
| DDoS | Medium | Rate limiting, CDN |
| Insider threat | Medium | Audit logs, RBAC |

---

## 2. Authentication & Authorization

### 2.1 Authentication Flow

```
User                    App                     Backend
  │                      │                         │
  │──── Phone Number ────>│                         │
  │                      │──── OTP Request ────────>│
  │                      │<──── OTP Generated ──────│
  │                      │                         │
  │<──── Enter OTP ───────│                         │
  │──── OTP Code ────────>│                         │
  │                      │──── Verify OTP ─────────>│
  │                      │<──── JWT Tokens ─────────│
  │                      │                         │
  │<──── Authenticated ───│                         │
```

### 2.2 OTP Security

**Configuration:**
```yaml
otp:
  length: 6
  expiry: 300  # 5 minutes
  max_attempts: 5
  rate_limit:
    per_phone: 10/day
    per_ip: 100/hour
  cooldown: 60  # seconds between requests
```

**Implementation:**
```go
// Secure OTP generation
func generateSecureOTP(length int) string {
    const digits = "0123456789"
    b := make([]byte, length)
    if _, err := rand.Read(b); err != nil {
        panic(err)  // Cryptographic failure is critical
    }
    for i := range b {
        b[i] = digits[int(b[i])%len(digits)]
    }
    return string(b)
}

// OTP verification with timing attack prevention
func verifyOTP(stored, provided string) bool {
    return subtle.ConstantTimeCompare(
        []byte(stored),
        []byte(provided),
    ) == 1
}
```

### 2.3 JWT Token Security

**Token Structure:**
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_uuid",
    "phone": "+2348012345678",
    "tier": "silver",
    "type": "access",
    "iat": 1704067200,
    "exp": 1704068100,
    "jti": "unique_token_id"
  }
}
```

**Security Measures:**
| Measure | Implementation |
|---------|----------------|
| Short expiry | Access: 15min, Refresh: 7 days |
| Token rotation | New refresh token on each use |
| Blacklisting | Redis-based immediate revocation |
| JTI tracking | Prevent token reuse |

### 2.4 Transaction PIN

**Requirements:**
- 6 digits
- Not sequential (123456)
- Not repeating (111111)
- Different from phone number
- Changed every 90 days (recommended)

**Storage:**
```go
// PIN hashing with bcrypt
func hashPIN(pin string) (string, error) {
    // Cost of 12 provides good security/performance balance
    hash, err := bcrypt.GenerateFromPassword([]byte(pin), 12)
    return string(hash), err
}

func verifyPIN(hash, pin string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))
    return err == nil
}
```

### 2.5 Authorization Matrix

| Resource | Bronze | Silver | Gold | Admin |
|----------|--------|--------|------|-------|
| View wallet | ✅ | ✅ | ✅ | ✅ |
| Transfer (own) | ✅ | ✅ | ✅ | ✅ |
| Transfer limit | ₦50K | ₦200K | ₦500K | N/A |
| Create gigs | ✅ | ✅ | ✅ | ✅ |
| Apply for loan | ❌ | ✅ | ✅ | N/A |
| Admin panel | ❌ | ❌ | ❌ | ✅ |

---

## 3. Data Protection

### 3.1 Data Classification

| Classification | Examples | Protection |
|---------------|----------|------------|
| **Critical** | PIN hash, payment tokens | Encrypted, limited access |
| **Sensitive** | BVN, account numbers | Encrypted, audit logged |
| **Personal** | Name, email, phone | Encrypted at rest |
| **Internal** | Gig content, messages | Standard protection |
| **Public** | Public gig listings | None required |

### 3.2 Encryption Standards

**At Rest:**
- Database: AES-256 encryption
- File storage: Server-side encryption (S3 SSE)
- Backups: Encrypted with separate keys

**In Transit:**
- TLS 1.3 minimum
- Strong cipher suites only
- Certificate pinning in mobile app

**Key Management:**
```yaml
# Key hierarchy
master_key:
  storage: AWS KMS / HashiCorp Vault
  rotation: Annual

data_encryption_keys:
  derivation: HKDF from master
  rotation: Quarterly

application_keys:
  jwt_secret: 256-bit random
  rotation: On breach detection
```

### 3.3 Data Retention

| Data Type | Retention Period | Deletion Method |
|-----------|-----------------|-----------------|
| Transaction records | 7 years | Archive then delete |
| User profiles | Account lifetime + 1 year | Soft delete, then hard |
| OTP codes | 5 minutes | Auto-expire (Redis TTL) |
| Session tokens | 7 days | Auto-expire |
| Audit logs | 3 years | Archive |
| Support tickets | 2 years | Archive |

### 3.4 PII Handling

**Collected PII:**
- Phone number (identifier)
- Full name
- Email (optional)
- BVN/NIN (KYC)
- Bank account details
- Transaction history

**Protection Measures:**
- Minimize collection
- Encrypt sensitive fields
- Access logging
- Right to deletion (GDPR compliance)
- Data export on request

---

## 4. Application Security

### 4.1 Input Validation

**Backend Validation:**
```go
type TransferRequest struct {
    RecipientPhone string  `json:"recipient_phone" validate:"required,e164"`
    Amount         float64 `json:"amount" validate:"required,gt=0,lte=1000000"`
    Note           string  `json:"note" validate:"max=200"`
    PIN            string  `json:"pin" validate:"required,len=6,numeric"`
}

func (h *Handler) Transfer(c *fiber.Ctx) error {
    var req TransferRequest
    if err := c.BodyParser(&req); err != nil {
        return fiber.ErrBadRequest
    }

    // Validate struct
    if err := validate.Struct(req); err != nil {
        return c.Status(400).JSON(formatValidationErrors(err))
    }

    // Additional business validation
    if req.Amount < config.MinTransferAmount {
        return c.Status(400).JSON(ErrorResponse("Amount too low"))
    }

    // Process transfer...
}
```

### 4.2 SQL Injection Prevention

**Use Parameterized Queries:**
```go
// CORRECT: Parameterized query
db.Where("phone = ?", phone).First(&user)

// INCORRECT: String concatenation
db.Raw("SELECT * FROM users WHERE phone = '" + phone + "'")
```

### 4.3 XSS Prevention

**Mobile App:**
- React Native/Flutter auto-escape by default
- Never use `dangerouslySetInnerHTML` or equivalent
- Sanitize user-generated content before display

**API Response:**
- Set `Content-Type: application/json`
- Escape HTML in JSON strings

### 4.4 Rate Limiting

**Configuration:**
```go
rateLimits := map[string]RateLimit{
    "otp_request":     {Max: 5, Window: 15 * time.Minute},
    "login":           {Max: 5, Window: 5 * time.Minute},
    "transfer":        {Max: 30, Window: time.Hour},
    "api_general":     {Max: 100, Window: time.Minute},
}
```

**Implementation:**
```go
func RateLimitMiddleware(key string, limit RateLimit) fiber.Handler {
    return func(c *fiber.Ctx) error {
        identifier := c.IP() + ":" + key

        count, _ := redis.Incr(ctx, identifier).Result()
        if count == 1 {
            redis.Expire(ctx, identifier, limit.Window)
        }

        if count > int64(limit.Max) {
            return c.Status(429).JSON(fiber.Map{
                "error": "Rate limit exceeded",
                "retry_after": redis.TTL(ctx, identifier).Val().Seconds(),
            })
        }

        return c.Next()
    }
}
```

### 4.5 Security Headers

```go
app.Use(func(c *fiber.Ctx) error {
    c.Set("X-Content-Type-Options", "nosniff")
    c.Set("X-Frame-Options", "DENY")
    c.Set("X-XSS-Protection", "1; mode=block")
    c.Set("Content-Security-Policy", "default-src 'self'")
    c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
    return c.Next()
})
```

---

## 5. Infrastructure Security

### 5.1 Network Security

**Architecture:**
```
Internet
    │
    ▼
┌────────────────┐
│   CloudFlare   │  ← DDoS protection, WAF
│      CDN       │
└───────┬────────┘
        │
        ▼
┌────────────────┐
│ Load Balancer  │  ← TLS termination
└───────┬────────┘
        │
        ▼
┌────────────────────────────────────┐
│           Private VPC              │
│  ┌──────────┐  ┌──────────────┐   │
│  │   API    │  │   Workers    │   │
│  │ Servers  │  │              │   │
│  └────┬─────┘  └──────────────┘   │
│       │                           │
│  ┌────▼─────┐  ┌──────────────┐   │
│  │PostgreSQL│  │    Redis     │   │
│  │(Private) │  │  (Private)   │   │
│  └──────────┘  └──────────────┘   │
└────────────────────────────────────┘
```

**Security Groups:**
```yaml
# API servers
api_sg:
  inbound:
    - port: 443, source: load_balancer_sg
  outbound:
    - port: 5432, destination: database_sg
    - port: 6379, destination: redis_sg
    - port: 443, destination: 0.0.0.0/0  # External APIs

# Database
database_sg:
  inbound:
    - port: 5432, source: api_sg
  outbound: []  # No egress
```

### 5.2 Container Security

**Dockerfile Best Practices:**
```dockerfile
# Use specific version, not latest
FROM golang:1.22-alpine AS builder

# Don't run as root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Minimize attack surface
FROM alpine:3.19
RUN apk --no-cache add ca-certificates

# Copy only what's needed
COPY --from=builder /app/bin /app/bin

# Run as non-root
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --spider http://localhost:8080/health || exit 1
```

### 5.3 Secret Management

**Environment Variables (Development):**
```bash
# .env (never commit!)
JWT_SECRET=your-256-bit-secret
DB_PASSWORD=your-db-password
PAYSTACK_SECRET_KEY=sk_test_xxx
```

**Production Secrets:**
```yaml
# Kubernetes secrets (base64 encoded)
apiVersion: v1
kind: Secret
metadata:
  name: hustlex-secrets
type: Opaque
data:
  JWT_SECRET: <base64>
  DB_PASSWORD: <base64>
  PAYSTACK_SECRET_KEY: <base64>

# Or use external secrets manager
# AWS Secrets Manager / HashiCorp Vault
```

### 5.4 Logging & Monitoring

**Security Logging:**
```go
// Log security events
type SecurityEvent struct {
    Timestamp time.Time `json:"timestamp"`
    EventType string    `json:"event_type"`
    UserID    string    `json:"user_id,omitempty"`
    IP        string    `json:"ip"`
    UserAgent string    `json:"user_agent"`
    Success   bool      `json:"success"`
    Details   string    `json:"details,omitempty"`
}

// Events to log
- Authentication attempts (success/failure)
- Authorization failures
- Sensitive data access
- Configuration changes
- Admin actions
- Rate limit triggers
```

**Alerts:**
| Event | Threshold | Action |
|-------|-----------|--------|
| Failed logins | 10/min same IP | Block IP |
| Rate limit hits | 100/min | Alert ops |
| Auth failures | 5/min same user | Lock account |
| Large withdrawal | > ₦500,000 | Manual review |

---

## 6. Compliance & Regulations

### 6.1 Regulatory Requirements

**Nigeria:**
- CBN Guidelines for Payment Service Providers
- NDPR (Nigeria Data Protection Regulation)
- KYC/AML requirements

**International:**
- PCI-DSS (payment card handling)
- GDPR (if serving EU users)

### 6.2 KYC/AML Compliance

**KYC Tiers:**
| Tier | Requirements | Limits |
|------|--------------|--------|
| Basic | Phone verification | ₦50K/day |
| Standard | BVN verification | ₦200K/day |
| Enhanced | BVN + NIN + Address | ₦500K/day |
| Premium | Full KYC + Review | ₦1M/day |

**AML Monitoring:**
```go
// Suspicious activity detection
type SuspiciousActivityRule struct {
    Name        string
    Condition   func(tx Transaction) bool
    Action      func(tx Transaction) error
}

rules := []SuspiciousActivityRule{
    {
        Name: "Large single transaction",
        Condition: func(tx Transaction) bool {
            return tx.Amount > 500000
        },
        Action: flagForReview,
    },
    {
        Name: "Rapid multiple transactions",
        Condition: func(tx Transaction) bool {
            count := getTransactionCount(tx.UserID, 1*time.Hour)
            return count > 20
        },
        Action: flagForReview,
    },
    {
        Name: "New account large activity",
        Condition: func(tx Transaction) bool {
            user := getUser(tx.UserID)
            return user.CreatedAt.After(time.Now().AddDate(0, 0, -7)) &&
                   tx.Amount > 100000
        },
        Action: flagForReview,
    },
}
```

### 6.3 Data Protection Compliance

**NDPR Requirements:**
- Lawful basis for processing
- Data minimization
- Storage limitation
- Security measures
- Data subject rights
- Breach notification (72 hours)

**Implementation:**
- Privacy policy displayed
- Consent tracking
- Data export functionality
- Deletion requests honored
- Breach response plan

---

## 7. Incident Response

### 7.1 Incident Classification

| Severity | Description | Response Time | Examples |
|----------|-------------|---------------|----------|
| P1 Critical | Service down, data breach | 15 min | Payment system down, credentials leaked |
| P2 High | Major feature broken | 1 hour | Auth issues, payment delays |
| P3 Medium | Feature degraded | 4 hours | Slow performance, minor bugs |
| P4 Low | Minor issues | 24 hours | UI bugs, typos |

### 7.2 Incident Response Process

```
┌─────────────┐
│   Detect    │  ← Monitoring, user reports, alerts
└──────┬──────┘
       ▼
┌─────────────┐
│   Assess    │  ← Classify severity, identify scope
└──────┬──────┘
       ▼
┌─────────────┐
│   Contain   │  ← Stop the bleeding, isolate
└──────┬──────┘
       ▼
┌─────────────┐
│  Eradicate  │  ← Fix the root cause
└──────┬──────┘
       ▼
┌─────────────┐
│   Recover   │  ← Restore services
└──────┬──────┘
       ▼
┌─────────────┐
│   Review    │  ← Post-mortem, improve
└─────────────┘
```

### 7.3 Security Incident Playbooks

**Account Compromise:**
1. Disable affected account
2. Invalidate all sessions
3. Review recent activity
4. Contact user via verified channel
5. Reset credentials
6. Document timeline

**Data Breach:**
1. Identify scope of breach
2. Preserve evidence
3. Notify security team
4. Assess notification requirements
5. Notify affected users (if required)
6. File regulatory reports
7. Conduct post-mortem

### 7.4 Communication Templates

**User Notification (Security):**
```
Subject: Important Security Notice for Your HustleX Account

Dear [Name],

We detected unusual activity on your account on [date]. As a precaution,
we have [action taken].

What happened: [brief description]

What we've done: [actions]

What you should do:
1. Change your PIN
2. Review recent transactions
3. Contact us if you notice anything suspicious

If you have questions, contact support@hustlex.app

The HustleX Security Team
```

---

## 8. Security Best Practices

### 8.1 For Developers

**Code Review Checklist:**
- [ ] No hardcoded secrets
- [ ] Input validation on all endpoints
- [ ] Parameterized queries only
- [ ] Proper error handling (no stack traces to users)
- [ ] Logging without sensitive data
- [ ] Rate limiting applied
- [ ] Authentication required where needed
- [ ] Authorization checks in place

### 8.2 For Operations

**Daily:**
- [ ] Review security alerts
- [ ] Check system health
- [ ] Review failed login attempts

**Weekly:**
- [ ] Review access logs
- [ ] Check for unusual patterns
- [ ] Update security patches

**Monthly:**
- [ ] Review user permissions
- [ ] Rotate non-critical secrets
- [ ] Security training updates

### 8.3 For Users

**Security Tips (in-app):**
- Use a strong, unique PIN
- Enable biometric authentication
- Don't share OTP codes
- Verify transactions before confirming
- Report suspicious activity immediately

---

## 9. Security Checklist

### Pre-Launch Checklist

**Authentication:**
- [ ] OTP rate limiting configured
- [ ] JWT expiration set correctly
- [ ] PIN complexity enforced
- [ ] Session management working
- [ ] Logout invalidates tokens

**Data Protection:**
- [ ] Database encrypted at rest
- [ ] TLS configured correctly
- [ ] Sensitive data encrypted
- [ ] PII handling compliant
- [ ] Backup encryption enabled

**Application:**
- [ ] Input validation on all endpoints
- [ ] SQL injection prevented
- [ ] XSS prevented
- [ ] CSRF protection (if web)
- [ ] Rate limiting configured
- [ ] Security headers set

**Infrastructure:**
- [ ] Firewall rules configured
- [ ] Network isolation in place
- [ ] Secrets in secure storage
- [ ] Logging enabled
- [ ] Monitoring configured

**Compliance:**
- [ ] Privacy policy published
- [ ] Terms of service published
- [ ] KYC process implemented
- [ ] AML monitoring active
- [ ] Data retention policy defined

### Ongoing Security Tasks

**Weekly:**
- [ ] Review security alerts
- [ ] Check for dependency vulnerabilities
- [ ] Review access logs

**Monthly:**
- [ ] Security patch review
- [ ] Access permission audit
- [ ] Backup restoration test

**Quarterly:**
- [ ] Penetration testing
- [ ] Security training
- [ ] Policy review

**Annually:**
- [ ] Full security audit
- [ ] Compliance assessment
- [ ] Incident response drill

---

## Contact Information

**Security Team:**
- Email: security@hustlex.app
- Urgent: +234 XXX XXX XXXX

**Bug Bounty:**
- Report vulnerabilities to: security@hustlex.app
- Responsible disclosure appreciated

---

*This document is confidential and for internal use only.*

**Version 1.0 | Last Updated: January 2024**
