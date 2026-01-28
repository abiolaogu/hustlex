# HustleX - System Architecture Document

## Architecture Overview

HustleX follows a microservices architecture designed for high scalability, maintainability, and resilience. The system is built to handle 100,000+ concurrent users with sub-200ms response times.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ Flutter iOS  │  │Flutter Android│  │  Web Client  │  │ USSD Gateway │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              API GATEWAY                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │  Kong/Traefik: Rate Limiting, Auth, Load Balancing, SSL Termination  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           MICROSERVICES LAYER                                │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌──────────┐  │
│  │   User     │ │    Gig     │ │  Savings   │ │   Wallet   │ │  Credit  │  │
│  │  Service   │ │  Service   │ │  Service   │ │  Service   │ │ Service  │  │
│  │   (Go)     │ │   (Go)     │ │   (Go)     │ │   (Go)     │ │  (Go)    │  │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘ └──────────┘  │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │  Learning  │ │ Community  │ │Notification│ │  Payment   │              │
│  │  Service   │ │  Service   │ │  Service   │ │  Service   │              │
│  │   (Go)     │ │   (Go)     │ │   (Go)     │ │   (Go)     │              │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘              │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                            DATA LAYER                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │ PostgreSQL │ │   Redis    │ │Elasticsearch│ │    S3      │              │
│  │  (Primary) │ │  (Cache)   │ │  (Search)   │ │ (Storage)  │              │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘              │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MESSAGE QUEUE / EVENTS                               │
├─────────────────────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │     Apache Kafka / RabbitMQ: Event Sourcing, Async Processing        │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Technology Stack

### Backend (Go)
| Component | Technology | Justification |
|-----------|------------|---------------|
| Language | Go 1.22+ | High concurrency, low latency, efficient memory |
| Framework | Fiber/Gin | Fast HTTP router, middleware support |
| ORM | GORM | PostgreSQL support, migrations |
| Validation | go-playground/validator | Struct validation |
| Auth | JWT + OAuth2 | Industry standard, stateless |
| API Docs | Swagger/OpenAPI | Auto-generated documentation |

### Mobile (Flutter)
| Component | Technology | Justification |
|-----------|------------|---------------|
| Framework | Flutter 3.x | Cross-platform, single codebase |
| State Management | Riverpod | Scalable, testable state management |
| HTTP Client | Dio | Interceptors, retry logic |
| Local Storage | Hive | Fast, encrypted local DB |
| Push Notifications | Firebase Cloud Messaging | Reliable delivery |

### Data Layer
| Component | Technology | Justification |
|-----------|------------|---------------|
| Primary DB | PostgreSQL 15 | ACID, JSON support, robust |
| Cache | Redis 7 | Sub-ms latency, pub/sub |
| Search | Elasticsearch 8 | Full-text search, analytics |
| Object Storage | MinIO/S3 | Scalable file storage |
| Message Queue | Kafka | Event streaming, durability |

### Infrastructure
| Component | Technology | Justification |
|-----------|------------|---------------|
| Container | Docker | Consistent environments |
| Orchestration | Kubernetes | Auto-scaling, self-healing |
| CI/CD | GitHub Actions | Integrated, affordable |
| Monitoring | Prometheus + Grafana | Metrics, alerting |
| Logging | ELK Stack | Centralized logging |
| CDN | Cloudflare | DDoS protection, caching |

---

## Database Schema

### Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     users       │       │     skills      │       │   user_skills   │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ phone           │◄──────│ name            │──────►│ user_id (FK)    │
│ email           │       │ category        │       │ skill_id (FK)   │
│ full_name       │       │ description     │       │ proficiency     │
│ profile_image   │       │ icon            │       │ verified        │
│ bio             │       └─────────────────┘       └─────────────────┘
│ location        │
│ credit_score    │       ┌─────────────────┐       ┌─────────────────┐
│ tier            │       │      gigs       │       │  gig_proposals  │
│ created_at      │       ├─────────────────┤       ├─────────────────┤
│ updated_at      │       │ id (PK)         │       │ id (PK)         │
└─────────────────┘       │ client_id (FK)  │◄──────│ gig_id (FK)     │
        │                 │ title           │       │ hustler_id (FK) │
        │                 │ description     │       │ cover_letter    │
        │                 │ category        │       │ proposed_price  │
        ▼                 │ budget_min      │       │ delivery_days   │
┌─────────────────┐       │ budget_max      │       │ status          │
│    wallets      │       │ deadline        │       └─────────────────┘
├─────────────────┤       │ status          │
│ id (PK)         │       │ created_at      │       ┌─────────────────┐
│ user_id (FK)    │       └─────────────────┘       │   gig_reviews   │
│ balance         │               │                 ├─────────────────┤
│ escrow_balance  │               ▼                 │ id (PK)         │
│ currency        │       ┌─────────────────┐       │ gig_id (FK)     │
│ updated_at      │       │  gig_contracts  │       │ reviewer_id(FK) │
└─────────────────┘       ├─────────────────┤       │ rating          │
                          │ id (PK)         │       │ review_text     │
┌─────────────────┐       │ gig_id (FK)     │       │ created_at      │
│ savings_circles │       │ hustler_id (FK) │       └─────────────────┘
├─────────────────┤       │ agreed_price    │
│ id (PK)         │       │ status          │
│ name            │       │ started_at      │
│ type            │       │ completed_at    │
│ contribution_amt│       └─────────────────┘
│ frequency       │
│ total_members   │       ┌─────────────────┐       ┌─────────────────┐
│ current_round   │       │ circle_members  │       │  contributions  │
│ created_by (FK) │       ├─────────────────┤       ├─────────────────┤
│ status          │       │ id (PK)         │       │ id (PK)         │
│ created_at      │◄──────│ circle_id (FK)  │       │ member_id (FK)  │
└─────────────────┘       │ user_id (FK)    │──────►│ amount          │
                          │ position        │       │ due_date        │
                          │ joined_at       │       │ paid_at         │
                          │ status          │       │ status          │
                          └─────────────────┘       └─────────────────┘

┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│  transactions   │       │  credit_scores  │       │     courses     │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │       │ id (PK)         │       │ id (PK)         │
│ wallet_id (FK)  │       │ user_id (FK)    │       │ title           │
│ type            │       │ score           │       │ description     │
│ amount          │       │ tier            │       │ skill_id (FK)   │
│ reference       │       │ gig_completion  │       │ duration_mins   │
│ status          │       │ rating_avg      │       │ difficulty      │
│ metadata        │       │ savings_record  │       │ modules (JSON)  │
│ created_at      │       │ account_age     │       │ created_at      │
└─────────────────┘       │ updated_at      │       └─────────────────┘
                          └─────────────────┘
```

---

## API Design

### RESTful Endpoints

#### Authentication
```
POST   /api/v1/auth/register          # Register new user
POST   /api/v1/auth/login             # Login with phone/OTP
POST   /api/v1/auth/verify-otp        # Verify OTP
POST   /api/v1/auth/refresh           # Refresh access token
POST   /api/v1/auth/logout            # Logout (invalidate token)
```

#### Users
```
GET    /api/v1/users/me               # Get current user profile
PUT    /api/v1/users/me               # Update profile
GET    /api/v1/users/:id              # Get user public profile
POST   /api/v1/users/me/skills        # Add skill to profile
DELETE /api/v1/users/me/skills/:id    # Remove skill
```

#### Gigs
```
GET    /api/v1/gigs                   # List gigs (with filters)
POST   /api/v1/gigs                   # Create new gig
GET    /api/v1/gigs/:id               # Get gig details
PUT    /api/v1/gigs/:id               # Update gig
DELETE /api/v1/gigs/:id               # Delete gig
POST   /api/v1/gigs/:id/proposals     # Submit proposal
GET    /api/v1/gigs/:id/proposals     # List proposals (client only)
POST   /api/v1/gigs/:id/accept/:pid   # Accept proposal
```

#### Contracts
```
GET    /api/v1/contracts              # List user's contracts
GET    /api/v1/contracts/:id          # Get contract details
POST   /api/v1/contracts/:id/deliver  # Submit deliverables
POST   /api/v1/contracts/:id/approve  # Approve & release payment
POST   /api/v1/contracts/:id/dispute  # Open dispute
POST   /api/v1/contracts/:id/review   # Submit review
```

#### Savings Circles
```
GET    /api/v1/circles                # List user's circles
POST   /api/v1/circles                # Create circle
GET    /api/v1/circles/:id            # Get circle details
POST   /api/v1/circles/:id/join       # Join circle
POST   /api/v1/circles/:id/contribute # Make contribution
GET    /api/v1/circles/:id/members    # List members
POST   /api/v1/circles/:id/payout     # Trigger payout (admin)
```

#### Wallet
```
GET    /api/v1/wallet                 # Get wallet balance
GET    /api/v1/wallet/transactions    # Transaction history
POST   /api/v1/wallet/deposit         # Initialize deposit
POST   /api/v1/wallet/withdraw        # Request withdrawal
POST   /api/v1/wallet/transfer        # P2P transfer
```

#### Credit
```
GET    /api/v1/credit/score           # Get credit score
GET    /api/v1/credit/history         # Score change history
GET    /api/v1/credit/loans           # List loans
POST   /api/v1/credit/loans/apply     # Apply for loan
POST   /api/v1/credit/loans/:id/repay # Repay loan
```

#### Learning
```
GET    /api/v1/courses                # List courses
GET    /api/v1/courses/:id            # Get course details
POST   /api/v1/courses/:id/enroll     # Enroll in course
POST   /api/v1/courses/:id/progress   # Update progress
GET    /api/v1/courses/:id/certificate # Get certificate
```

---

## Security Architecture

### Authentication Flow
```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  User   │───►│   SMS   │───►│   OTP   │───►│  Token  │
│  Login  │    │ Gateway │    │ Verify  │    │  Issue  │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
     │                                            │
     │         ┌──────────────────────────────────┘
     │         ▼
     │    ┌─────────────────┐
     │    │  Access Token   │ (15 min expiry)
     │    │  Refresh Token  │ (30 day expiry)
     │    └─────────────────┘
     │              │
     │              ▼
     │    ┌─────────────────┐
     └───►│ Protected APIs  │ (Bearer token)
          └─────────────────┘
```

### Security Measures
1. **Transport**: TLS 1.3 for all communications
2. **Authentication**: JWT with RS256 signing
3. **Authorization**: Role-based access control (RBAC)
4. **Rate Limiting**: 100 req/min per IP, 1000 req/min per user
5. **Input Validation**: Strict schema validation on all inputs
6. **SQL Injection**: Parameterized queries via ORM
7. **XSS Prevention**: Content-Type enforcement, output encoding
8. **CSRF Protection**: SameSite cookies, CSRF tokens
9. **Data Encryption**: AES-256 for sensitive data at rest
10. **Audit Logging**: All sensitive operations logged

---

## Scalability Design

### Horizontal Scaling
- Stateless services allow unlimited horizontal scaling
- Database read replicas for read-heavy operations
- Redis cluster for session/cache distribution
- Kubernetes HPA for auto-scaling based on CPU/memory

### Caching Strategy
```
┌─────────────────┐
│   Client App    │
└────────┬────────┘
         │ (Cache-Control headers)
         ▼
┌─────────────────┐
│      CDN        │ ◄── Static assets, images
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   API Gateway   │ ◄── Response caching (5-60s TTL)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Redis       │ ◄── Session, hot data (5min-24hr TTL)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   PostgreSQL    │ ◄── Source of truth
└─────────────────┘
```

### Load Estimates
| Metric | Day 1 | Month 6 | Year 1 |
|--------|-------|---------|--------|
| DAU | 1,000 | 50,000 | 200,000 |
| Peak RPS | 50 | 2,500 | 10,000 |
| DB Size | 1 GB | 50 GB | 500 GB |
| Redis Memory | 256 MB | 4 GB | 16 GB |

---

## Deployment Architecture

### Kubernetes Cluster
```yaml
# Production namespace
- user-service: 3 replicas, 2 CPU, 4Gi RAM
- gig-service: 3 replicas, 2 CPU, 4Gi RAM
- savings-service: 3 replicas, 2 CPU, 4Gi RAM
- wallet-service: 3 replicas (critical), 4 CPU, 8Gi RAM
- notification-service: 2 replicas, 1 CPU, 2Gi RAM
- payment-service: 3 replicas (critical), 4 CPU, 8Gi RAM

# Data tier
- PostgreSQL: Primary + 2 read replicas
- Redis: 6-node cluster (3 primary, 3 replica)
- Elasticsearch: 3-node cluster
```

### CI/CD Pipeline
```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│  Push   │───►│  Build  │───►│  Test   │───►│ Deploy  │───►│ Monitor │
│         │    │ Docker  │    │  Unit   │    │ Staging │    │ Metrics │
└─────────┘    └─────────┘    │  E2E    │    │  Prod   │    └─────────┘
                              └─────────┘    └─────────┘
```

---

## Monitoring & Observability

### Metrics (Prometheus)
- Request rate, latency, error rate (RED)
- Database connection pool stats
- Redis hit/miss ratio
- Kafka consumer lag
- Custom business metrics (gigs created, transactions)

### Logging (ELK)
- Structured JSON logs
- Request correlation IDs
- PII masking in logs
- 30-day retention

### Alerting
| Condition | Severity | Action |
|-----------|----------|--------|
| Error rate > 5% | Critical | PagerDuty + Slack |
| p95 latency > 500ms | Warning | Slack |
| DB CPU > 80% | Warning | Slack |
| Wallet service down | Critical | PagerDuty + SMS |
| Fraud score > threshold | Critical | Auto-block + alert |

---

## Disaster Recovery

### Backup Strategy
- PostgreSQL: Continuous WAL archiving + daily snapshots
- Redis: RDB snapshots every 15 min
- S3: Cross-region replication

### Recovery Objectives
- RPO (Recovery Point Objective): 5 minutes
- RTO (Recovery Time Objective): 30 minutes

### Failover Plan
1. Database failover to read replica (automatic)
2. Service failover to secondary region (manual)
3. DNS switch via Cloudflare (< 5 min)

---

*Document Version: 1.0*  
*Last Updated: January 2026*
