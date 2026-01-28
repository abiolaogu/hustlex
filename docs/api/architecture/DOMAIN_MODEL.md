# HustleX Domain Model

## Overview

This document describes the Domain-Driven Design (DDD) architecture for HustleX, a Nigerian fintech super app combining gig marketplace, social savings (Ajo/Esusu), and alternative credit scoring.

---

## Bounded Contexts

### 1. Identity Context

**Purpose:** User identity, authentication, authorization, and profile management

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| User | A person registered on the platform |
| Credential | Authentication method (phone, email, PIN) |
| Session | Active authentication state with tokens |
| Permission | Authorization to perform a specific action |
| Role | Collection of permissions assigned to a user |
| Profile | User's personal and business information |
| KYC | Know Your Customer verification status |
| Tier | User level determining limits and features |

**Aggregates:**

```
UserAggregate (root: User)
├── User (entity)
│   ├── id: UserID
│   ├── phone: PhoneNumber
│   ├── email: Email (optional)
│   ├── fullName: FullName
│   ├── tier: UserTier
│   ├── status: UserStatus
│   └── createdAt: Timestamp
├── Credential (entity)
│   ├── type: CredentialType (phone_otp, pin, biometric)
│   ├── value: HashedValue
│   └── lastUsedAt: Timestamp
├── Profile (value object)
│   ├── avatar: URL
│   ├── bio: string
│   ├── location: Location
│   └── skills: []Skill
└── KYCRecord (entity)
    ├── bvn: BVN (encrypted)
    ├── idType: IDType
    ├── idNumber: string (encrypted)
    ├── verificationStatus: KYCStatus
    └── verifiedAt: Timestamp
```

**Domain Events:**
- `UserRegistered`
- `UserVerified`
- `UserTierUpgraded`
- `KYCCompleted`
- `UserSuspended`
- `UserReactivated`

---

### 2. Wallet Context

**Purpose:** Digital money management, transactions, and ledger

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| Wallet | Digital money container for a user |
| Balance | Current available funds in smallest unit (kobo) |
| Escrow | Funds held pending release (gig payments) |
| Transaction | Recorded money movement |
| Ledger | Immutable log of all financial changes |
| Transfer | Movement of funds between wallets |
| Deposit | Adding external funds to wallet |
| Withdrawal | Removing funds to external bank |

**Aggregates:**

```
WalletAggregate (root: Wallet)
├── Wallet (entity)
│   ├── id: WalletID
│   ├── userID: UserID
│   ├── availableBalance: Money
│   ├── escrowBalance: Money
│   ├── savingsBalance: Money
│   ├── ledgerBalance: Money
│   ├── currency: Currency
│   ├── status: WalletStatus
│   └── version: int64 (optimistic locking)
└── TransactionPIN (value object)
    ├── hash: HashedPIN
    ├── attempts: int
    └── lockedUntil: Timestamp

LedgerAggregate (root: Ledger)
├── Ledger (entity)
│   └── walletID: WalletID
└── LedgerEntry (entity)
    ├── id: EntryID
    ├── type: EntryType (credit, debit)
    ├── amount: Money
    ├── balanceAfter: Money
    ├── reference: TransactionReference
    ├── description: string
    └── createdAt: Timestamp
```

**Domain Events:**
- `WalletCreated`
- `WalletCredited`
- `WalletDebited`
- `FundsHeldInEscrow`
- `FundsReleasedFromEscrow`
- `WalletLocked`
- `WalletUnlocked`
- `WithdrawalInitiated`
- `WithdrawalCompleted`
- `WithdrawalFailed`

---

### 3. Gig Context

**Purpose:** Marketplace for work opportunities

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| Gig | Work opportunity posted by a client |
| Client | User who posts gigs and pays for work |
| Freelancer | User who completes gigs for payment |
| Proposal | Freelancer's offer to complete a gig |
| Contract | Agreement between client and freelancer |
| Milestone | Deliverable checkpoint with payment |
| Escrow | Funds held during gig execution |
| Review | Post-completion feedback and rating |
| Dispute | Conflict requiring resolution |

**Aggregates:**

```
GigAggregate (root: Gig)
├── Gig (entity)
│   ├── id: GigID
│   ├── clientID: UserID
│   ├── title: string
│   ├── description: string
│   ├── category: Category
│   ├── budget: Money
│   ├── deadline: Timestamp
│   ├── status: GigStatus
│   └── requirements: []Requirement
├── Milestone (entity)
│   ├── id: MilestoneID
│   ├── title: string
│   ├── amount: Money
│   ├── dueDate: Timestamp
│   └── status: MilestoneStatus
└── Requirement (value object)
    ├── skill: Skill
    └── level: SkillLevel

ContractAggregate (root: Contract)
├── Contract (entity)
│   ├── id: ContractID
│   ├── gigID: GigID
│   ├── clientID: UserID
│   ├── freelancerID: UserID
│   ├── agreedAmount: Money
│   ├── status: ContractStatus
│   └── startedAt: Timestamp
├── ContractMilestone (entity)
│   ├── milestoneID: MilestoneID
│   ├── status: MilestoneStatus
│   ├── deliveredAt: Timestamp
│   └── approvedAt: Timestamp
└── Payment (entity)
    ├── id: PaymentID
    ├── amount: Money
    ├── status: PaymentStatus
    └── releasedAt: Timestamp

ProposalAggregate (root: Proposal)
├── Proposal (entity)
│   ├── id: ProposalID
│   ├── gigID: GigID
│   ├── freelancerID: UserID
│   ├── coverLetter: string
│   ├── proposedAmount: Money
│   ├── estimatedDuration: Duration
│   ├── status: ProposalStatus
│   └── submittedAt: Timestamp
└── Attachment (value object)
    ├── url: URL
    └── type: AttachmentType
```

**Domain Events:**
- `GigPosted`
- `GigUpdated`
- `GigCancelled`
- `ProposalSubmitted`
- `ProposalAccepted`
- `ProposalRejected`
- `ContractCreated`
- `MilestoneCompleted`
- `MilestoneApproved`
- `PaymentReleased`
- `GigCompleted`
- `ReviewSubmitted`
- `DisputeOpened`
- `DisputeResolved`

---

### 4. Savings Context

**Purpose:** Ajo/Esusu social savings circles

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| Circle | Savings group with rotating payouts |
| Member | Participant in a savings circle |
| Contribution | Periodic payment to the circle pool |
| Payout | Distribution of collected funds to a member |
| Round | Complete cycle of all contributions |
| Position | Order in the payout queue |
| Cycle | Duration between contributions |
| Pool | Total collected contributions |

**Aggregates:**

```
CircleAggregate (root: Circle)
├── Circle (entity)
│   ├── id: CircleID
│   ├── name: string
│   ├── description: string
│   ├── creatorID: UserID
│   ├── contributionAmount: Money
│   ├── frequency: Frequency (daily, weekly, monthly)
│   ├── maxMembers: int
│   ├── currentRound: int
│   ├── status: CircleStatus
│   └── startDate: Timestamp
├── Member (entity)
│   ├── userID: UserID
│   ├── position: int
│   ├── joinedAt: Timestamp
│   ├── status: MemberStatus
│   └── contributionsMade: int
├── Round (entity)
│   ├── number: int
│   ├── recipientID: UserID
│   ├── totalCollected: Money
│   ├── status: RoundStatus
│   └── completedAt: Timestamp
└── ContributionSchedule (value object)
    ├── nextDueDate: Timestamp
    └── reminderSent: bool

ContributionAggregate (root: Contribution)
└── Contribution (entity)
    ├── id: ContributionID
    ├── circleID: CircleID
    ├── memberID: UserID
    ├── roundNumber: int
    ├── amount: Money
    ├── status: ContributionStatus
    ├── dueDate: Timestamp
    └── paidAt: Timestamp
```

**Domain Events:**
- `CircleCreated`
- `MemberJoined`
- `MemberLeft`
- `ContributionScheduled`
- `ContributionReceived`
- `ContributionMissed`
- `RoundStarted`
- `PayoutDistributed`
- `RoundCompleted`
- `CircleCompleted`

---

### 5. Credit Context

**Purpose:** Credit scoring and lending

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| CreditProfile | User's creditworthiness record |
| CreditScore | Numerical assessment (0-1000) |
| ScoreFactor | Component contributing to score |
| Loan | Borrowed funds with repayment terms |
| Principal | Original borrowed amount |
| Interest | Cost of borrowing |
| Repayment | Payment toward loan balance |
| Default | Failure to repay as agreed |
| Collection | Recovery of overdue amounts |

**Aggregates:**

```
CreditProfileAggregate (root: CreditProfile)
├── CreditProfile (entity)
│   ├── id: ProfileID
│   ├── userID: UserID
│   ├── currentScore: CreditScore
│   ├── tier: CreditTier
│   ├── maxLoanAmount: Money
│   ├── lastCalculatedAt: Timestamp
│   └── factors: []ScoreFactor
├── ScoreHistory (entity)
│   ├── score: CreditScore
│   ├── calculatedAt: Timestamp
│   └── reason: string
└── ScoreFactor (value object)
    ├── name: string (gig_completion, savings_consistency, repayment_history)
    ├── weight: float
    └── value: float

LoanAggregate (root: Loan)
├── Loan (entity)
│   ├── id: LoanID
│   ├── userID: UserID
│   ├── principal: Money
│   ├── interestRate: Percentage
│   ├── totalDue: Money
│   ├── outstandingBalance: Money
│   ├── term: Duration
│   ├── status: LoanStatus
│   ├── disbursedAt: Timestamp
│   └── dueDate: Timestamp
├── RepaymentSchedule (entity)
│   ├── installments: []Installment
│   └── frequency: Frequency
├── Installment (value object)
│   ├── number: int
│   ├── amount: Money
│   ├── dueDate: Timestamp
│   └── status: InstallmentStatus
└── Repayment (entity)
    ├── id: RepaymentID
    ├── amount: Money
    ├── principalPortion: Money
    ├── interestPortion: Money
    ├── paidAt: Timestamp
    └── method: PaymentMethod
```

**Domain Events:**
- `CreditScoreCalculated`
- `CreditTierChanged`
- `LoanApplicationSubmitted`
- `LoanApproved`
- `LoanRejected`
- `LoanDisbursed`
- `RepaymentReceived`
- `RepaymentMissed`
- `LoanDefaulted`
- `LoanFullyRepaid`
- `CollectionInitiated`

---

### 6. Notification Context

**Purpose:** User communication across channels

**Ubiquitous Language:**
| Term | Definition |
|------|------------|
| Notification | Message to be delivered to user |
| Channel | Delivery method (push, SMS, email, in-app) |
| Template | Predefined notification format |
| Preference | User's channel settings |
| Delivery | Attempt to send notification |
| Read | User has seen the notification |

**Aggregates:**

```
NotificationAggregate (root: Notification)
├── Notification (entity)
│   ├── id: NotificationID
│   ├── userID: UserID
│   ├── type: NotificationType
│   ├── title: string
│   ├── body: string
│   ├── data: map[string]any
│   ├── priority: Priority
│   ├── status: NotificationStatus
│   ├── readAt: Timestamp
│   └── createdAt: Timestamp
└── Delivery (entity)
    ├── channel: Channel
    ├── status: DeliveryStatus
    ├── attemptedAt: Timestamp
    └── error: string

PreferenceAggregate (root: Preference)
└── Preference (entity)
    ├── userID: UserID
    ├── channels: map[NotificationType][]Channel
    ├── quietHours: TimeRange
    └── updatedAt: Timestamp
```

**Domain Events:**
- `NotificationCreated`
- `NotificationSent`
- `NotificationDelivered`
- `NotificationFailed`
- `NotificationRead`
- `PreferenceUpdated`

---

## Context Map

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          IDENTITY CONTEXT                                │
│                      (Upstream - Shared Kernel)                          │
│                                                                          │
│   Provides: UserID, authentication, authorization                        │
└─────────────────────────────────────────────────────────────────────────┘
                    │                    │                    │
                    │ U/D                │ U/D                │ U/D
                    ▼                    ▼                    ▼
    ┌───────────────────┐    ┌───────────────────┐    ┌───────────────────┐
    │   WALLET CONTEXT  │    │    GIG CONTEXT    │    │  CREDIT CONTEXT   │
    │                   │    │                   │    │                   │
    │ Provides:         │    │ Provides:         │    │ Provides:         │
    │ - Money transfers │    │ - Work tracking   │    │ - Credit scoring  │
    │ - Balance mgmt    │    │ - Payments        │    │ - Loan mgmt       │
    │ - Escrow          │    │ - Reviews         │    │                   │
    └─────────┬─────────┘    └─────────┬─────────┘    └─────────┬─────────┘
              │                        │                        │
              │ Partner               │ Partner                │ Partner
              │                        │                        │
              ▼                        ▼                        ▼
    ┌─────────────────────────────────────────────────────────────────────┐
    │                        SAVINGS CONTEXT                               │
    │                                                                      │
    │   Consumes: Wallet (for contributions/payouts)                       │
    │             Identity (for members)                                   │
    │   Publishes: ContributionMade, PayoutDistributed                     │
    └─────────────────────────────────────────────────────────────────────┘
                                      │
                                      │ U/D
                                      ▼
    ┌─────────────────────────────────────────────────────────────────────┐
    │                      NOTIFICATION CONTEXT                            │
    │                        (Downstream)                                  │
    │                                                                      │
    │   Consumes: Events from all contexts                                 │
    │   Provides: Multi-channel notification delivery                      │
    └─────────────────────────────────────────────────────────────────────┘
```

### Integration Patterns

| Upstream | Downstream | Pattern |
|----------|------------|---------|
| Identity | Wallet | Shared Kernel (UserID) |
| Identity | Gigs | Shared Kernel (UserID) |
| Identity | Credit | Shared Kernel (UserID) |
| Wallet | Gigs | Customer-Supplier (escrow) |
| Wallet | Savings | Customer-Supplier (contributions) |
| Wallet | Credit | Customer-Supplier (disbursement) |
| Gigs | Credit | Published Language (completion events) |
| Savings | Credit | Published Language (consistency events) |
| All | Notification | Anti-Corruption Layer |

---

## Cross-Context Domain Events

### Event Catalog

```go
// ============================================================================
// Identity Context Events
// ============================================================================

type UserRegistered struct {
    EventID     string    `json:"event_id"`
    UserID      string    `json:"user_id"`
    Phone       string    `json:"phone"`
    Email       string    `json:"email,omitempty"`
    FullName    string    `json:"full_name"`
    ReferredBy  string    `json:"referred_by,omitempty"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type UserVerified struct {
    EventID     string    `json:"event_id"`
    UserID      string    `json:"user_id"`
    Method      string    `json:"method"` // bvn, id_card, selfie
    OccurredAt  time.Time `json:"occurred_at"`
}

type UserTierUpgraded struct {
    EventID     string    `json:"event_id"`
    UserID      string    `json:"user_id"`
    OldTier     string    `json:"old_tier"`
    NewTier     string    `json:"new_tier"`
    Reason      string    `json:"reason"`
    OccurredAt  time.Time `json:"occurred_at"`
}

// ============================================================================
// Wallet Context Events
// ============================================================================

type WalletCredited struct {
    EventID       string    `json:"event_id"`
    WalletID      string    `json:"wallet_id"`
    UserID        string    `json:"user_id"`
    Amount        int64     `json:"amount"`
    Currency      string    `json:"currency"`
    Source        string    `json:"source"` // deposit, gig_payment, payout, refund
    Reference     string    `json:"reference"`
    NewBalance    int64     `json:"new_balance"`
    OccurredAt    time.Time `json:"occurred_at"`
}

type WalletDebited struct {
    EventID       string    `json:"event_id"`
    WalletID      string    `json:"wallet_id"`
    UserID        string    `json:"user_id"`
    Amount        int64     `json:"amount"`
    Currency      string    `json:"currency"`
    Destination   string    `json:"destination"` // withdrawal, transfer, escrow, contribution
    Reference     string    `json:"reference"`
    NewBalance    int64     `json:"new_balance"`
    OccurredAt    time.Time `json:"occurred_at"`
}

type FundsHeldInEscrow struct {
    EventID     string    `json:"event_id"`
    WalletID    string    `json:"wallet_id"`
    Amount      int64     `json:"amount"`
    Reference   string    `json:"reference"` // contract_id or gig_id
    Reason      string    `json:"reason"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type FundsReleasedFromEscrow struct {
    EventID      string    `json:"event_id"`
    WalletID     string    `json:"wallet_id"`
    Amount       int64     `json:"amount"`
    Reference    string    `json:"reference"`
    RecipientID  string    `json:"recipient_id"`
    OccurredAt   time.Time `json:"occurred_at"`
}

// ============================================================================
// Gig Context Events
// ============================================================================

type GigPosted struct {
    EventID     string    `json:"event_id"`
    GigID       string    `json:"gig_id"`
    ClientID    string    `json:"client_id"`
    Title       string    `json:"title"`
    Category    string    `json:"category"`
    Budget      int64     `json:"budget"`
    Currency    string    `json:"currency"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type ContractCreated struct {
    EventID       string    `json:"event_id"`
    ContractID    string    `json:"contract_id"`
    GigID         string    `json:"gig_id"`
    ClientID      string    `json:"client_id"`
    FreelancerID  string    `json:"freelancer_id"`
    Amount        int64     `json:"amount"`
    Currency      string    `json:"currency"`
    OccurredAt    time.Time `json:"occurred_at"`
}

type GigCompleted struct {
    EventID       string    `json:"event_id"`
    GigID         string    `json:"gig_id"`
    ContractID    string    `json:"contract_id"`
    ClientID      string    `json:"client_id"`
    FreelancerID  string    `json:"freelancer_id"`
    Amount        int64     `json:"amount"`
    Rating        float64   `json:"rating,omitempty"`
    OccurredAt    time.Time `json:"occurred_at"`
}

type GigPaymentReleased struct {
    EventID       string    `json:"event_id"`
    ContractID    string    `json:"contract_id"`
    FreelancerID  string    `json:"freelancer_id"`
    GrossAmount   int64     `json:"gross_amount"`
    PlatformFee   int64     `json:"platform_fee"`
    NetAmount     int64     `json:"net_amount"`
    OccurredAt    time.Time `json:"occurred_at"`
}

// ============================================================================
// Savings Context Events
// ============================================================================

type CircleCreated struct {
    EventID            string    `json:"event_id"`
    CircleID           string    `json:"circle_id"`
    CreatorID          string    `json:"creator_id"`
    Name               string    `json:"name"`
    ContributionAmount int64     `json:"contribution_amount"`
    Frequency          string    `json:"frequency"`
    MaxMembers         int       `json:"max_members"`
    OccurredAt         time.Time `json:"occurred_at"`
}

type MemberJoined struct {
    EventID     string    `json:"event_id"`
    CircleID    string    `json:"circle_id"`
    UserID      string    `json:"user_id"`
    Position    int       `json:"position"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type ContributionMade struct {
    EventID       string    `json:"event_id"`
    CircleID      string    `json:"circle_id"`
    MemberID      string    `json:"member_id"`
    Amount        int64     `json:"amount"`
    RoundNumber   int       `json:"round_number"`
    OccurredAt    time.Time `json:"occurred_at"`
}

type PayoutDistributed struct {
    EventID       string    `json:"event_id"`
    CircleID      string    `json:"circle_id"`
    RecipientID   string    `json:"recipient_id"`
    Amount        int64     `json:"amount"`
    RoundNumber   int       `json:"round_number"`
    OccurredAt    time.Time `json:"occurred_at"`
}

// ============================================================================
// Credit Context Events
// ============================================================================

type CreditScoreCalculated struct {
    EventID     string    `json:"event_id"`
    UserID      string    `json:"user_id"`
    OldScore    int       `json:"old_score"`
    NewScore    int       `json:"new_score"`
    Factors     []string  `json:"factors"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type LoanDisbursed struct {
    EventID     string    `json:"event_id"`
    LoanID      string    `json:"loan_id"`
    UserID      string    `json:"user_id"`
    Principal   int64     `json:"principal"`
    Interest    int64     `json:"interest"`
    TotalDue    int64     `json:"total_due"`
    DueDate     time.Time `json:"due_date"`
    OccurredAt  time.Time `json:"occurred_at"`
}

type LoanRepaymentReceived struct {
    EventID           string    `json:"event_id"`
    LoanID            string    `json:"loan_id"`
    UserID            string    `json:"user_id"`
    Amount            int64     `json:"amount"`
    OutstandingBalance int64    `json:"outstanding_balance"`
    IsFullyRepaid     bool      `json:"is_fully_repaid"`
    OccurredAt        time.Time `json:"occurred_at"`
}

type LoanDefaulted struct {
    EventID           string    `json:"event_id"`
    LoanID            string    `json:"loan_id"`
    UserID            string    `json:"user_id"`
    OutstandingAmount int64     `json:"outstanding_amount"`
    DaysOverdue       int       `json:"days_overdue"`
    OccurredAt        time.Time `json:"occurred_at"`
}
```

---

## Shared Kernel

### Value Objects

| Value Object | Used By | Purpose |
|--------------|---------|---------|
| `Money` | All | Monetary amounts with currency |
| `UserID` | All | Reference to user across contexts |
| `PhoneNumber` | Identity, Notification | Nigerian phone format |
| `Email` | Identity, Notification | Email address |
| `Timestamp` | All | UTC timestamps |
| `Percentage` | Credit, Gigs | Percentage values |

### Specifications

| Specification | Context | Purpose |
|---------------|---------|---------|
| `SufficientBalance` | Wallet | Check if wallet has enough funds |
| `ValidLoanAmount` | Credit | Check loan within allowed range |
| `EligibleForPayout` | Savings | Check member can receive payout |
| `CanAcceptProposal` | Gigs | Check gig can accept proposals |

---

## Anti-Corruption Layers

### External Payment Gateway (Paystack)

```
┌─────────────────────────────────────────────────────────────────┐
│                       WALLET CONTEXT                             │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              PaymentGatewayPort (interface)              │    │
│  │  - InitiateDeposit(userID, amount) -> PaymentReference  │    │
│  │  - VerifyPayment(reference) -> PaymentResult            │    │
│  │  - InitiateWithdrawal(bankAccount, amount) -> Result    │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ Adapter
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                  INFRASTRUCTURE LAYER                            │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │           PaystackAdapter (implements port)              │    │
│  │  - Translates domain concepts to Paystack API           │    │
│  │  - Handles webhook signature verification               │    │
│  │  - Maps Paystack responses to domain events             │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      PAYSTACK API                                │
└─────────────────────────────────────────────────────────────────┘
```

### External SMS Gateway (Termii)

```
┌─────────────────────────────────────────────────────────────────┐
│                   NOTIFICATION CONTEXT                           │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                SMSGatewayPort (interface)                │    │
│  │  - SendOTP(phone, code) -> DeliveryResult               │    │
│  │  - SendNotification(phone, message) -> DeliveryResult   │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ Adapter
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│           TermiiAdapter / AfricasTalkingAdapter                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Implementation Guidelines

### 1. Aggregate Design Rules

- **Single Aggregate per Transaction**: Never modify multiple aggregates in one transaction
- **Reference by ID**: Aggregates reference others by ID, never by object reference
- **Small Aggregates**: Keep aggregates small for better concurrency
- **Eventual Consistency**: Accept eventual consistency between aggregates

### 2. Event Sourcing Considerations

While not implementing full event sourcing, follow these patterns:
- Record domain events within aggregates
- Publish events after successful persistence
- Use events for cross-context communication
- Consider event store for audit trail

### 3. Repository Contracts

```go
type Repository[T any] interface {
    FindByID(ctx context.Context, id string) (T, error)
    Save(ctx context.Context, entity T) error
    SaveWithEvents(ctx context.Context, entity T) error
}
```

### 4. Application Service Pattern

```go
type CommandHandler[C any, R any] interface {
    Handle(ctx context.Context, cmd C) (R, error)
}

type QueryHandler[Q any, R any] interface {
    Handle(ctx context.Context, query Q) (R, error)
}
```

---

## Migration Strategy

### Phase 1: Foundation
1. Create shared kernel value objects
2. Set up domain event infrastructure
3. Implement Wallet context (highest transaction volume)

### Phase 2: Core Contexts
4. Implement Identity context
5. Implement Gig context
6. Wire up event publishing

### Phase 3: Supporting Contexts
7. Implement Savings context
8. Implement Credit context
9. Implement Notification context

### Phase 4: Integration
10. Implement anti-corruption layers
11. Migrate existing data
12. Deprecate old code paths

---

**Document Version:** 1.0.0
**Last Updated:** January 2026
**Author:** Architecture Team
