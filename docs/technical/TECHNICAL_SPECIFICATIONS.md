# HustleX Technical Specifications

> Version 1.0 | Last Updated: January 2024

## Table of Contents

1. [System Overview](#system-overview)
2. [Architecture](#architecture)
3. [Technology Stack](#technology-stack)
4. [Data Models](#data-models)
5. [API Specifications](#api-specifications)
6. [Security Specifications](#security-specifications)
7. [Performance Requirements](#performance-requirements)
8. [Integration Specifications](#integration-specifications)
9. [Infrastructure](#infrastructure)

---

## 1. System Overview

### 1.1 Purpose

HustleX is a comprehensive financial super-app designed for the Nigerian gig economy, combining:
- **Gig Marketplace**: Connect freelancers with clients for short-term work
- **Savings Circles (Ajo/Esusu)**: Digitized traditional social savings
- **Alternative Credit Scoring**: Build credit through platform activity
- **Digital Wallet**: Secure money management and transfers

### 1.2 Scope

| Component | Description |
|-----------|-------------|
| Backend API | RESTful Go API server |
| Mobile App | Cross-platform Flutter application (iOS/Android) |
| Background Jobs | Async task processing for payments, notifications |
| Infrastructure | Containerized deployment on Kubernetes |

### 1.3 Target Users

| User Type | Description |
|-----------|-------------|
| Gig Workers | Freelancers seeking short-term work |
| Gig Clients | Businesses/individuals posting jobs |
| Savers | Users participating in savings circles |
| Credit Seekers | Users building credit for microloans |

---

## 2. Architecture

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENTS                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   iOS App    │  │ Android App  │  │  Admin Web   │          │
│  │   (Flutter)  │  │   (Flutter)  │  │  (Future)    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼─────────────────┼─────────────────┼───────────────────┘
          │                 │                 │
          └────────────────┬┴─────────────────┘
                           │ HTTPS
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                      LOAD BALANCER                               │
│                   (Nginx / Cloud LB)                            │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                       API LAYER                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Go Fiber API                          │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐    │   │
│  │  │   Auth   │ │  Wallet  │ │   Gigs   │ │ Savings  │    │   │
│  │  │ Handler  │ │ Handler  │ │ Handler  │ │ Handler  │    │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘    │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐                 │   │
│  │  │  Credit  │ │ Profile  │ │  Common  │                 │   │
│  │  │ Handler  │ │ Handler  │ │ Handler  │                 │   │
│  │  └──────────┘ └──────────┘ └──────────┘                 │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────┬───────────────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  PostgreSQL  │  │    Redis     │  │    Asynq     │
│   Database   │  │    Cache     │  │  Job Queue   │
│              │  │              │  │              │
│  - Users     │  │  - Sessions  │  │  - Payments  │
│  - Wallets   │  │  - OTPs      │  │  - Notifs    │
│  - Gigs      │  │  - Rate Lim  │  │  - Scheduled │
│  - Savings   │  │  - Cache     │  │              │
│  - Credit    │  │              │  │              │
└──────────────┘  └──────────────┘  └──────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                   EXTERNAL SERVICES                              │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │ Paystack │ │  Termii  │ │ Firebase │ │    S3    │           │
│  │ Payments │ │   SMS    │ │   FCM    │ │ Storage  │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      BACKEND API                             │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                    HTTP Layer                        │    │
│  │  ┌─────────┐  ┌────────────┐  ┌─────────────────┐   │    │
│  │  │ Router  │─>│ Middleware │─>│    Handlers     │   │    │
│  │  │ (Fiber) │  │ (Auth/Rate)│  │ (REST Endpoints)│   │    │
│  │  └─────────┘  └────────────┘  └────────┬────────┘   │    │
│  └────────────────────────────────────────┼────────────┘    │
│                                           │                  │
│  ┌────────────────────────────────────────▼────────────┐    │
│  │                   Service Layer                      │    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐         │    │
│  │  │   Auth    │ │  Wallet   │ │   Gigs    │         │    │
│  │  │  Service  │ │  Service  │ │  Service  │         │    │
│  │  └───────────┘ └───────────┘ └───────────┘         │    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐         │    │
│  │  │  Savings  │ │  Credit   │ │  Payment  │         │    │
│  │  │  Service  │ │  Service  │ │  Service  │         │    │
│  │  └───────────┘ └───────────┘ └───────────┘         │    │
│  └────────────────────────────────────────┬────────────┘    │
│                                           │                  │
│  ┌────────────────────────────────────────▼────────────┐    │
│  │                    Data Layer                        │    │
│  │  ┌───────────┐ ┌───────────┐ ┌───────────┐         │    │
│  │  │   GORM    │ │  go-redis │ │   Asynq   │         │    │
│  │  │  (ORM)    │ │  Client   │ │  Client   │         │    │
│  │  └───────────┘ └───────────┘ └───────────┘         │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### 2.3 Mobile App Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     FLUTTER APP                              │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                 Presentation Layer                   │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐             │    │
│  │  │ Screens │  │ Widgets │  │Providers│             │    │
│  │  │  (UI)   │  │(Reusable)│ │(Riverpod)│            │    │
│  │  └────┬────┘  └────┬────┘  └────┬────┘             │    │
│  └───────┼────────────┼────────────┼───────────────────┘    │
│          │            │            │                         │
│          └────────────┴─────┬──────┘                         │
│                             │                                │
│  ┌──────────────────────────▼──────────────────────────┐    │
│  │                  Domain Layer                        │    │
│  │  ┌───────────────────────────────────────────────┐  │    │
│  │  │              Repositories                      │  │    │
│  │  │  (Data Access Abstraction)                    │  │    │
│  │  └───────────────────────────────────────────────┘  │    │
│  └──────────────────────────┬──────────────────────────┘    │
│                             │                                │
│  ┌──────────────────────────▼──────────────────────────┐    │
│  │                   Data Layer                         │    │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐          │    │
│  │  │   API    │  │  Local   │  │  Secure  │          │    │
│  │  │  Client  │  │ Storage  │  │ Storage  │          │    │
│  │  │  (Dio)   │  │  (Hive)  │  │(Keychain)│          │    │
│  │  └──────────┘  └──────────┘  └──────────┘          │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Technology Stack

### 3.1 Backend

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Language | Go | 1.21+ | Primary backend language |
| Framework | Fiber | 2.x | HTTP framework |
| ORM | GORM | 1.25+ | Database ORM |
| Database | PostgreSQL | 16 | Primary data store |
| Cache | Redis | 7 | Caching, sessions, queues |
| Job Queue | Asynq | 0.24+ | Background job processing |
| Auth | golang-jwt | 5.x | JWT token handling |

### 3.2 Mobile

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Framework | Flutter | 3.16+ | Cross-platform UI |
| Language | Dart | 3.x | Mobile development |
| State | Riverpod | 2.4+ | State management |
| Navigation | GoRouter | 13.x | Declarative routing |
| HTTP | Dio | 5.x | HTTP client |
| Storage | Hive | 2.x | Local database |
| Secure Storage | flutter_secure_storage | 9.x | Credentials |
| Payments | flutter_paystack | 1.x | Payment integration |

### 3.3 Infrastructure

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Container | Docker | 24.x | Containerization |
| Orchestration | Kubernetes | 1.28+ | Container orchestration |
| CI/CD | GitHub Actions | - | Continuous integration |
| Monitoring | Prometheus | 2.x | Metrics collection |
| Logging | ELK Stack | 8.x | Log aggregation |

---

## 4. Data Models

### 4.1 Entity Relationship Diagram

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│     User     │     │    Wallet    │     │ Transaction  │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │──┐  │ id (PK)      │──┐  │ id (PK)      │
│ phone        │  │  │ user_id (FK) │◄─┤  │ wallet_id(FK)│◄─┐
│ email        │  │  │ balance      │  │  │ type         │  │
│ full_name    │  │  │ escrow_bal   │  │  │ amount       │  │
│ tier         │  │  │ savings_bal  │  │  │ reference    │  │
│ bvn_verified │  │  │ currency     │  │  │ status       │  │
│ pin_hash     │  │  └──────────────┘  │  │ description  │  │
└──────────────┘  │                     │  └──────────────┘  │
       │          │                     │                    │
       │          └─────────────────────┴────────────────────┘
       │
       ├─────────────────────────────────────────────────────┐
       │                                                     │
       ▼                                                     │
┌──────────────┐     ┌──────────────┐     ┌──────────────┐  │
│     Gig      │     │ GigProposal  │     │ GigContract  │  │
├──────────────┤     ├──────────────┤     ├──────────────┤  │
│ id (PK)      │──┐  │ id (PK)      │──┐  │ id (PK)      │  │
│ client_id(FK)│◄─┤  │ gig_id (FK)  │◄─┤  │ proposal_id  │◄─┤
│ title        │  │  │ freelancer_id│  │  │ agreed_price │  │
│ description  │  │  │ price        │  │  │ platform_fee │  │
│ category     │  │  │ cover_letter │  │  │ delivery_days│  │
│ budget_min   │  │  │ status       │  │  │ status       │  │
│ budget_max   │  │  └──────────────┘  │  └──────────────┘  │
│ status       │  │                     │                    │
└──────────────┘  │                     │                    │
                  └─────────────────────┴────────────────────┘
       │
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│SavingsCircle │     │CircleMember  │     │ Contribution │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │──┐  │ id (PK)      │──┐  │ id (PK)      │
│ creator_id   │◄─┤  │ circle_id(FK)│◄─┤  │ circle_id(FK)│◄─┐
│ name         │  │  │ user_id (FK) │  │  │ member_id(FK)│  │
│ type         │  │  │ position     │  │  │ amount       │  │
│ contribution │  │  │ role         │  │  │ due_date     │  │
│ frequency    │  │  │ status       │  │  │ status       │  │
│ start_date   │  │  └──────────────┘  │  └──────────────┘  │
└──────────────┘  │                     │                    │
                  └─────────────────────┴────────────────────┘
       │
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ CreditScore  │     │    Loan      │     │LoanRepayment │
├──────────────┤     ├──────────────┤     ├──────────────┤
│ id (PK)      │     │ id (PK)      │──┐  │ id (PK)      │
│ user_id (FK) │◄────│ user_id (FK) │◄─┤  │ loan_id (FK) │◄─┐
│ score        │     │ amount       │  │  │ amount       │  │
│ tier         │     │ interest     │  │  │ tx_id        │  │
│ payment_hist │     │ tenure       │  │  │ paid_at      │  │
│ savings_hist │     │ status       │  │  └──────────────┘  │
│ gig_hist     │     │ due_date     │  │                    │
└──────────────┘     └──────────────┘  └────────────────────┘
```

### 4.2 Core Models

#### User Model

```go
type User struct {
    ID              uuid.UUID       `json:"id"`
    Phone           string          `json:"phone"`
    Email           *string         `json:"email,omitempty"`
    FullName        string          `json:"full_name"`
    ProfilePhoto    *string         `json:"profile_photo,omitempty"`
    DateOfBirth     *time.Time      `json:"date_of_birth,omitempty"`
    Tier            string          `json:"tier"` // bronze, silver, gold, platinum
    BVNVerified     bool            `json:"bvn_verified"`
    NINVerified     bool            `json:"nin_verified"`
    PINHash         string          `json:"-"`
    ReferralCode    string          `json:"referral_code"`
    ReferredBy      *uuid.UUID      `json:"referred_by,omitempty"`
    DeviceTokens    []string        `json:"-"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    DeletedAt       *time.Time      `json:"-"`
}
```

#### Wallet Model

```go
type Wallet struct {
    ID             uuid.UUID  `json:"id"`
    UserID         uuid.UUID  `json:"user_id"`
    Balance        float64    `json:"balance"`
    EscrowBalance  float64    `json:"escrow_balance"`
    SavingsBalance float64    `json:"savings_balance"`
    Currency       string     `json:"currency"` // NGN
    DailyLimit     float64    `json:"daily_limit"`
    MonthlyLimit   float64    `json:"monthly_limit"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}
```

#### Transaction Model

```go
type Transaction struct {
    ID          uuid.UUID  `json:"id"`
    WalletID    uuid.UUID  `json:"wallet_id"`
    Type        string     `json:"type"` // credit, debit
    Category    string     `json:"category"` // deposit, withdrawal, transfer, gig_payment, savings, loan
    Amount      float64    `json:"amount"`
    Fee         float64    `json:"fee"`
    Reference   string     `json:"reference"`
    Description string     `json:"description"`
    Status      string     `json:"status"` // pending, completed, failed
    Metadata    JSONB      `json:"metadata,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
}
```

### 4.3 Database Indexes

```sql
-- Performance indexes
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_tier ON users(tier);
CREATE INDEX idx_transactions_wallet_date ON transactions(wallet_id, created_at DESC);
CREATE INDEX idx_gigs_category_status ON gigs(category, status);
CREATE INDEX idx_gigs_client_status ON gigs(client_id, status);
CREATE INDEX idx_contributions_circle_due ON contributions(circle_id, due_date);
CREATE INDEX idx_loans_user_status ON loans(user_id, status);

-- Full-text search
CREATE INDEX idx_gigs_search ON gigs USING gin(
    to_tsvector('english', title || ' ' || description)
);
```

---

## 5. API Specifications

### 5.1 API Standards

| Standard | Value |
|----------|-------|
| Protocol | HTTPS |
| Format | JSON |
| Authentication | Bearer JWT |
| Versioning | URL path (v1) |
| Rate Limiting | Per-endpoint |
| Pagination | Cursor-based |

### 5.2 Response Format

#### Success Response

```json
{
    "success": true,
    "data": {
        // Response payload
    },
    "meta": {
        "request_id": "uuid",
        "timestamp": "2024-01-15T10:30:00Z"
    }
}
```

#### Error Response

```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid input",
        "details": [
            {
                "field": "amount",
                "message": "Amount must be positive"
            }
        ]
    },
    "meta": {
        "request_id": "uuid",
        "timestamp": "2024-01-15T10:30:00Z"
    }
}
```

### 5.3 Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 400 | Invalid input data |
| UNAUTHORIZED | 401 | Missing/invalid auth |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| RATE_LIMITED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |
| INSUFFICIENT_FUNDS | 400 | Wallet balance too low |
| OTP_EXPIRED | 400 | OTP has expired |
| INVALID_PIN | 400 | Transaction PIN incorrect |

### 5.4 Rate Limits

| Endpoint Category | Limit | Window |
|-------------------|-------|--------|
| OTP Request | 5 | 15 minutes |
| Login/Auth | 5 | 5 minutes |
| Transactions | 30 | 1 hour |
| General API | 100 | 1 minute |

---

## 6. Security Specifications

### 6.1 Authentication

| Mechanism | Details |
|-----------|---------|
| Primary | OTP via SMS |
| Tokens | JWT (Access: 15min, Refresh: 7 days) |
| Transaction | 6-digit PIN (bcrypt hashed) |
| Biometric | Device-level (iOS Face ID, Android Fingerprint) |

### 6.2 Data Protection

| Data Type | Protection |
|-----------|------------|
| Passwords/PINs | bcrypt (cost 12) |
| JWT Secret | 256-bit random |
| API Keys | Environment variables |
| PII | Encrypted at rest (AES-256) |
| Data in Transit | TLS 1.3 |

### 6.3 Security Headers

```
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

---

## 7. Performance Requirements

### 7.1 Response Time Targets

| Operation | Target (p95) | Maximum |
|-----------|--------------|---------|
| Auth (OTP) | 200ms | 500ms |
| Read (GET) | 100ms | 300ms |
| Write (POST) | 200ms | 500ms |
| Payment | 500ms | 2000ms |
| Search | 300ms | 1000ms |

### 7.2 Throughput Targets

| Metric | Target |
|--------|--------|
| API RPS | 10,000 |
| Concurrent Users | 50,000 |
| Daily Transactions | 1,000,000 |
| Background Jobs/sec | 1,000 |

### 7.3 Availability Targets

| Metric | Target |
|--------|--------|
| Uptime | 99.9% |
| Recovery Time | < 5 minutes |
| Data Durability | 99.999999% |

---

## 8. Integration Specifications

### 8.1 Payment Gateway (Paystack)

| Feature | Endpoint |
|---------|----------|
| Initialize | POST /transaction/initialize |
| Verify | GET /transaction/verify/:reference |
| Transfer | POST /transfer |
| Recipient | POST /transferrecipient |

### 8.2 SMS Gateway (Termii)

| Feature | Endpoint |
|---------|----------|
| Send SMS | POST /api/sms/send |
| Verify Number | POST /api/insight/number/query |
| Delivery Report | Webhook |

### 8.3 Push Notifications (Firebase)

| Feature | Endpoint |
|---------|----------|
| Send | POST /v1/projects/{project}/messages:send |
| Topic Subscribe | POST /iid/v1:batchAdd |
| Topic Unsubscribe | POST /iid/v1:batchRemove |

---

## 9. Infrastructure

### 9.1 Production Environment

| Component | Specification |
|-----------|---------------|
| API Servers | 3x (2 vCPU, 4GB RAM) |
| Database | 1x Primary + 1x Replica (4 vCPU, 16GB RAM) |
| Redis | 1x (2 vCPU, 4GB RAM) |
| Worker | 2x (2 vCPU, 4GB RAM) |
| Load Balancer | Managed (Cloud LB) |

### 9.2 Scaling Triggers

| Metric | Scale Up | Scale Down |
|--------|----------|------------|
| CPU | > 70% | < 30% |
| Memory | > 80% | < 40% |
| Request Queue | > 100 | < 20 |

### 9.3 Backup Strategy

| Data | Frequency | Retention |
|------|-----------|-----------|
| Database | Hourly | 30 days |
| Transaction Logs | Real-time | 1 year |
| User Data | Daily | 90 days |
| Configurations | On change | Indefinite |

---

## Appendix A: Glossary

| Term | Definition |
|------|------------|
| Ajo | Traditional Nigerian rotating savings |
| Esusu | Yoruba term for thrift contribution |
| Gig | Short-term freelance job |
| Escrow | Funds held in trust until conditions met |
| BVN | Bank Verification Number (Nigerian ID) |
| NIN | National Identification Number |

## Appendix B: References

- [Go Programming Language](https://golang.org/)
- [Flutter Documentation](https://flutter.dev/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Paystack API Reference](https://paystack.com/docs/api/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
