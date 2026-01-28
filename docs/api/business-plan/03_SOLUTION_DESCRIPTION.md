# 03. Solution Description

---

## Product Vision

**HustleX is the financial super app that turns your hustle into financial opportunity.**

We combine four integrated modules—Wallet, Gig Marketplace, Savings Circles, and Credit—into a single platform where every transaction builds your financial identity.

---

## Product Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           HUSTLEX MOBILE APP                             │
│                    (Flutter - iOS & Android)                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐            │
│  │  WALLET   │  │   GIGS    │  │  SAVINGS  │  │  CREDIT   │            │
│  │           │  │           │  │           │  │           │            │
│  │ • Balance │  │ • Browse  │  │ • Circles │  │ • Score   │            │
│  │ • Send    │  │ • Post    │  │ • Join    │  │ • Loans   │            │
│  │ • Receive │  │ • Propose │  │ • Create  │  │ • History │            │
│  │ • Bills   │  │ • Escrow  │  │ • Contrib │  │ • Tips    │            │
│  │ • History │  │ • Review  │  │ • Payout  │  │ • Apply   │            │
│  └───────────┘  └───────────┘  └───────────┘  └───────────┘            │
│                                                                          │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                    HUSTLESCORE DATA ENGINE                         │  │
│  │         (Alternative Credit Scoring from Platform Activity)        │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ REST API
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         HUSTLEX BACKEND (Go)                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │ Auth Service│  │Wallet Svc   │  │ Gig Service │  │Credit Service│    │
│  │             │  │             │  │             │  │             │    │
│  │ • OTP      │  │ • Balance   │  │ • Listings  │  │ • Scoring   │    │
│  │ • JWT      │  │ • Transfer  │  │ • Proposals │  │ • Loans     │    │
│  │ • PIN      │  │ • Escrow    │  │ • Contracts │  │ • Repayment │    │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘    │
│                                                                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                      │
│  │Savings Svc  │  │Notif Service│  │ Jobs/Queue  │                      │
│  │             │  │             │  │             │                      │
│  │ • Circles  │  │ • SMS       │  │ • Asynq     │                      │
│  │ • Contrib  │  │ • Push      │  │ • Cron      │                      │
│  │ • Payouts  │  │ • Email     │  │ • Batch     │                      │
│  └─────────────┘  └─────────────┘  └─────────────┘                      │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
              ┌──────────┐   ┌──────────┐   ┌──────────┐
              │PostgreSQL│   │  Redis   │   │ Paystack │
              │(Primary) │   │(Cache/Q) │   │(Payments)│
              └──────────┘   └──────────┘   └──────────┘
```

---

## Module 1: HustleWallet

### Overview

HustleWallet is the foundation of the platform—a fully-featured digital wallet that enables instant payments, bill payments, and serves as the financial hub for all platform activities.

### Core Features

| Feature | Description | User Benefit |
|---------|-------------|--------------|
| **Virtual Account** | Dedicated account number (via Paystack) | Receive transfers from any bank |
| **P2P Transfers** | Instant transfers to other HustleX users | Free, instant payments |
| **Bank Transfers** | Send to any Nigerian bank account | Low-fee withdrawals (₦50) |
| **QR Payments** | Scan to pay at merchants | Contactless, quick checkout |
| **Bill Payments** | Airtime, data, utilities, subscriptions | One-stop-shop convenience |
| **Transaction History** | Full audit trail with search/filter | Financial tracking |
| **Multiple Wallets** | Main, Savings, Escrow balances | Organized money management |

### Technical Specifications

| Spec | Detail |
|------|--------|
| Virtual Account Provider | Paystack (primary), Flutterwave (backup) |
| Transaction Limits | ₦100 min, ₦5M max per transaction |
| Daily Limit (unverified) | ₦300,000 |
| Daily Limit (verified) | ₦5,000,000 |
| Withdrawal Fee | ₦50 flat |
| P2P Fee | Free |
| Settlement Time | Instant (intra-platform), 5-30 min (external) |

### User Flow: Send Money

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Select     │────▶│  Enter      │────▶│  Confirm    │
│  Recipient  │     │  Amount     │     │  PIN        │
└─────────────┘     └─────────────┘     └─────────────┘
                                               │
                                               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Success    │◀────│  Process    │◀────│  Verify     │
│  Receipt    │     │  Transfer   │     │  Balance    │
└─────────────┘     └─────────────┘     └─────────────┘
```

---

## Module 2: HustleGigs (Marketplace)

### Overview

HustleGigs is an escrow-protected gig marketplace where freelancers can find work and clients can hire with confidence. Every completed gig builds the worker's reputation and credit score.

### Core Features

| Feature | Description | User Benefit |
|---------|-------------|--------------|
| **Gig Posting** | Clients post jobs with budget, timeline, requirements | Easy talent discovery |
| **Gig Discovery** | Workers browse by category, location, budget | Find relevant work |
| **Proposals** | Workers submit proposals with quotes | Competitive marketplace |
| **Escrow Protection** | Client funds held until work approved | Guaranteed payment |
| **Milestone Payments** | Large jobs split into deliverable milestones | Cash flow for workers |
| **Reviews & Ratings** | Two-way rating system | Trust and reputation |
| **Skills Verification** | Verified badges for proven skills | Stand out from crowd |
| **Dispute Resolution** | Mediation for disagreements | Fair outcomes |

### Gig Categories

| Category | Sub-categories | Avg. Job Value |
|----------|---------------|----------------|
| **Home Services** | Plumbing, Electrical, Cleaning, AC repair | ₦15,000-50,000 |
| **Logistics** | Delivery, Moving, Errands | ₦2,000-20,000 |
| **Creative** | Design, Photography, Video | ₦20,000-200,000 |
| **Tech** | Development, IT support, Repair | ₦30,000-500,000 |
| **Beauty** | Hair, Makeup, Nails | ₦5,000-30,000 |
| **Events** | Catering, DJ, MC, Security | ₦20,000-100,000 |
| **Professional** | Writing, Consulting, Tutoring | ₦10,000-100,000 |
| **Automotive** | Mechanic, Car wash, Tow | ₦5,000-50,000 |

### Escrow Flow

```
CLIENT                          ESCROW                         WORKER
   │                              │                              │
   │  Post Gig + Fund Escrow     │                              │
   │─────────────────────────────▶│                              │
   │                              │                              │
   │                              │  Funds Held Securely         │
   │                              │                              │
   │                              │                Accept Gig    │
   │                              │◀─────────────────────────────│
   │                              │                              │
   │                              │  Work Begins                 │
   │                              │                              │
   │                              │          Submit Milestone    │
   │                              │◀─────────────────────────────│
   │                              │                              │
   │  Approve Milestone           │                              │
   │─────────────────────────────▶│                              │
   │                              │                              │
   │                              │  Release to Worker (- 1.5%)  │
   │                              │─────────────────────────────▶│
   │                              │                              │
   │                              │             ₦₦₦              │
   │                              │                              │
```

### Pricing

| Item | Fee | Who Pays |
|------|-----|----------|
| Gig Posting | Free | - |
| Escrow Fee | 1.5% of gig value | Deducted from payment |
| Withdrawal | ₦50 | Worker (if withdrawing to bank) |
| Instant Payout | +0.5% | Worker (optional) |
| Featured Listing | ₦500/week | Client (optional) |

---

## Module 3: HustleSave (Ajo/Esusu Circles)

### Overview

HustleSave digitizes the traditional Ajo/Esusu savings system—preserving the community trust and forced discipline while adding transparency, security, and automation.

### Core Features

| Feature | Description | User Benefit |
|---------|-------------|--------------|
| **Create Circle** | Start a savings circle, set rules | Be a community leader |
| **Join Circle** | Request to join existing circles | Access forced savings |
| **Auto-Contribution** | Automatic deductions on schedule | Never miss a payment |
| **Transparent Ledger** | All contributions visible to all members | Trust through transparency |
| **Payout Schedule** | Clear, fair rotation | Know when you'll receive |
| **Early Withdrawal** | Access funds early (with penalty) | Emergency flexibility |
| **Circle Admin Tools** | Manage members, reminders, disputes | Easy administration |
| **Circle Discovery** | Find public circles to join | Expand your network |

### Circle Types Supported

| Type | Contribution | Frequency | Duration | Payout |
|------|--------------|-----------|----------|--------|
| **Daily Ajo** | ₦500-5,000 | Daily | 20-30 days | Rotating daily |
| **Weekly Ajo** | ₦2,000-20,000 | Weekly | 10-20 weeks | Rotating weekly |
| **Monthly Esusu** | ₦10,000-100,000 | Monthly | 6-12 months | Rotating monthly |
| **Target Savings** | Flexible | Flexible | Goal-based | At goal completion |
| **Challenge Savings** | Fixed | Daily/Weekly | 30-90 days | At challenge end |

### Circle Creation Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CREATE A CIRCLE                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. BASIC INFO                                                       │
│     ┌─────────────────────────────────────────────────────────────┐ │
│     │ Circle Name: [My Savings Circle                          ]  │ │
│     │ Description: [Monthly savings for my neighborhood        ]  │ │
│     │ Privacy:     (●) Invite Only  ( ) Public                   │ │
│     └─────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  2. CONTRIBUTION SETTINGS                                            │
│     ┌─────────────────────────────────────────────────────────────┐ │
│     │ Amount:      ₦ [10,000      ]  per [Month ▼]               │ │
│     │ Start Date:  [March 1, 2026                               ] │ │
│     │ Duration:    [12] months/cycles                            │ │
│     └─────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  3. PAYOUT RULES                                                     │
│     ┌─────────────────────────────────────────────────────────────┐ │
│     │ Payout Order:  (●) Random  ( ) Bid  ( ) Fixed Rotation     │ │
│     │ Payout Day:    [1st] of each month                         │ │
│     │ Early Withdraw: [✓] Allow (10% penalty)                    │ │
│     └─────────────────────────────────────────────────────────────┘ │
│                                                                      │
│  4. INVITE MEMBERS                                                   │
│     ┌─────────────────────────────────────────────────────────────┐ │
│     │ [+] Add from contacts  [+] Share invite link               │ │
│     │                                                             │ │
│     │ Invited (3/12):                                             │ │
│     │ • Adaeze O. ✓ Accepted                                      │ │
│     │ • Chidi N.  ⏳ Pending                                       │ │
│     │ • Funke A.  ⏳ Pending                                       │ │
│     └─────────────────────────────────────────────────────────────┘ │
│                                                                      │
│                    [ Create Circle →]                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Pricing

| Item | Fee | Notes |
|------|-----|-------|
| Create Circle | Free | Unlimited circles |
| Join Circle | Free | Subject to approval |
| Per-Cycle Fee | ₦100-500 | Based on contribution amount |
| Early Withdrawal | 10% penalty | Returned to circle pool |
| Circle Insurance | 1% of pool | Optional protection |

---

## Module 4: HustleCredit

### Overview

HustleCredit provides instant micro-loans based on your HustleScore—an alternative credit score built from your platform activity. No collateral, no salary slip, no traditional credit history required.

### Core Features

| Feature | Description | User Benefit |
|---------|-------------|--------------|
| **HustleScore** | Credit score from platform activity | Build credit through behavior |
| **Score Dashboard** | See your score and factors | Understand your creditworthiness |
| **Score Tips** | Recommendations to improve score | Actionable improvement path |
| **Instant Loans** | Apply and receive funds in minutes | No paperwork, no waiting |
| **Flexible Terms** | 7-90 day repayment options | Match your cash flow |
| **Credit Builder** | Small loans to build history | Start your credit journey |
| **Repayment Options** | Auto-debit, manual, early payoff | Convenient repayment |
| **Credit History** | Full borrowing history | Track your progress |

### HustleScore Algorithm

**Score Range:** 100 - 850 (similar to FICO for familiarity)

**Scoring Factors:**

| Factor | Weight | Data Points |
|--------|--------|-------------|
| **Gig Performance** | 25% | Completion rate, on-time delivery, client ratings |
| **Savings Behavior** | 25% | Ajo contribution consistency, circle completion |
| **Platform Activity** | 20% | App usage, transaction frequency, tenure |
| **Payment Behavior** | 20% | Bill payments on time, loan repayments |
| **Social Trust** | 10% | Circle memberships, referrals, verification level |

**Score Tiers:**

| Score | Tier | Loan Eligibility | Interest Rate |
|-------|------|------------------|---------------|
| 100-299 | Building | Credit Builder only | 5% monthly |
| 300-499 | Bronze | Up to ₦50,000 | 4.5% monthly |
| 500-649 | Silver | Up to ₦150,000 | 4% monthly |
| 650-749 | Gold | Up to ₦300,000 | 3.5% monthly |
| 750-850 | Platinum | Up to ₦500,000 | 3% monthly |

### Loan Products

| Product | Amount | Term | Rate | Purpose |
|---------|--------|------|------|---------|
| **Credit Builder** | ₦5,000-20,000 | 7-14 days | 5%/mo | Build credit history |
| **Emergency** | ₦10,000-50,000 | 7-30 days | 4-5%/mo | Urgent needs |
| **Standard** | ₦50,000-200,000 | 30-90 days | 3-4%/mo | Business, personal |
| **Premium** | ₦200,000-500,000 | 60-180 days | 3%/mo | Large expenses |

### Loan Application Flow

```
START ──▶ Check HustleScore ──▶ See Eligible Amount
                                       │
                              ┌────────┴────────┐
                              ▼                 ▼
                         Eligible          Not Eligible
                              │                 │
                              ▼                 ▼
                       Select Amount      Show Tips to
                       & Term             Improve Score
                              │
                              ▼
                       Review Terms
                       (Rate, Total Due)
                              │
                              ▼
                       Confirm & Accept
                              │
                              ▼
                       Disbursement
                       (Instant to Wallet)
                              │
                              ▼
                       Repayment Schedule
                       Set Up
```

---

## HustleScore Data Engine

### Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    DATA COLLECTION LAYER                         │
├─────────────────────────────────────────────────────────────────┤
│  Gig Data    │  Savings Data  │  Transaction  │  Behavioral    │
│              │                │  Data         │  Data          │
│  • Jobs done │  • Contrib %   │  • Frequency  │  • App opens   │
│  • Ratings   │  • Circles     │  • Volume     │  • Features    │
│  • On-time % │  • Completion  │  • Patterns   │  • Tenure      │
└──────┬───────┴───────┬────────┴───────┬───────┴───────┬────────┘
       │               │                │               │
       └───────────────┴────────────────┴───────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    SCORING ENGINE                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Raw Data ──▶ Feature Engineering ──▶ Model Scoring ──▶ Score   │
│                                                                  │
│  Models:                                                         │
│  • Logistic Regression (baseline)                               │
│  • Random Forest (production)                                   │
│  • XGBoost (evaluation)                                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    DECISION ENGINE                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Score + Rules ──▶ Loan Amount ──▶ Interest Rate ──▶ Decision  │
│                                                                  │
│  Rules:                                                          │
│  • Minimum platform tenure (30 days)                            │
│  • Minimum transactions (10)                                    │
│  • No current default                                           │
│  • Debt-to-income check                                         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Data Privacy & Security

| Measure | Implementation |
|---------|----------------|
| Data Minimization | Only collect what's needed for scoring |
| Encryption | AES-256 at rest, TLS 1.3 in transit |
| Access Control | Role-based, audit logged |
| User Control | Users can view all data collected |
| Deletion Rights | Full data deletion on account close |
| NDPR Compliance | Registered, DPO appointed |

---

## Competitive Differentiators

### Moat Analysis (Moat Map)

```
                        SUSTAINABLE ADVANTAGE
                               HIGH
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
         │    DATA MOAT         │    NETWORK EFFECTS   │
         │    ★★★★★             │    ★★★★☆             │
         │                      │                      │
         │  - Proprietary       │  - More users =      │
         │    credit data       │    more gigs =       │
         │  - Platform          │    more data =       │
         │    behavior data     │    better credit     │
         │  - Years to          │  - Cross-side        │
         │    replicate         │    effects           │
         │                      │                      │
   LOW ──┼──────────────────────┼──────────────────────┼── HIGH
BREADTH  │                      │                      │  BREADTH
         │                      │                      │
         │    SWITCHING COSTS   │    REGULATORY MOAT   │
         │    ★★★★☆             │    ★★★☆☆             │
         │                      │                      │
         │  - Credit history    │  - CBN MFB license   │
         │    locked in         │  - Limited licenses  │
         │  - Savings history   │    available         │
         │  - Gig reputation    │  - Compliance        │
         │  - Social circles    │    expertise         │
         │                      │                      │
         └──────────────────────┼──────────────────────┘
                                │
                               LOW
                        SUSTAINABLE ADVANTAGE
```

### Feature Comparison vs. Competitors

| Feature | HustleX | OPay | PalmPay | Piggyvest | Carbon |
|---------|---------|------|---------|-----------|--------|
| Digital Wallet | ✅ | ✅ | ✅ | ❌ | ✅ |
| P2P Payments | ✅ | ✅ | ✅ | ❌ | ✅ |
| Gig Marketplace | ✅ | ❌ | ❌ | ❌ | ❌ |
| Escrow Protection | ✅ | ❌ | ❌ | ❌ | ❌ |
| Ajo/Esusu Circles | ✅ | ❌ | ❌ | ✅* | ❌ |
| Alternative Credit | ✅ | ❌ | ❌ | ❌ | ❌ |
| Micro-loans | ✅ | ✅ | ✅ | ❌ | ✅ |
| Credit Building | ✅ | ❌ | ❌ | ❌ | ❌ |
| Skills Verification | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Integration Level** | **High** | Low | Low | Low | Low |

*Piggyvest has "Safelock" but not true Ajo circles

---

## Kano Model Analysis

### Feature Classification

```
                    SATISFACTION
                        HIGH
                          │
                          │        DELIGHTERS
                          │        ★ HustleScore tips
                          │        ★ Credit builder loans
                          │        ★ Circle discovery
                          │        ★ Skills badges
              ────────────┼────────────────────
             ABSENT       │                   PRESENT
                          │
       MUST-HAVES         │        PERFORMANCE
       ✓ Secure wallet    │        ◆ Fast transfers
       ✓ Working app      │        ◆ Low fees
       ✓ Customer support │        ◆ High loan amounts
       ✓ Transaction      │        ◆ Quick loan approval
         history          │
                          │
                         LOW
                    SATISFACTION
```

### Feature Prioritization

| Priority | Feature | Type | Rationale |
|----------|---------|------|-----------|
| P0 | Wallet & transfers | Must-have | Table stakes |
| P0 | Escrow for gigs | Performance | Core differentiator |
| P0 | Savings circles | Performance | Core differentiator |
| P1 | Micro-loans | Performance | Revenue driver |
| P1 | HustleScore | Delighter | Unique value prop |
| P2 | Skills verification | Delighter | Trust building |
| P2 | Credit builder | Delighter | Long-term engagement |
| P3 | Circle discovery | Delighter | Network effects |

---

## Product Roadmap

### MVP (Month 1-3)

**Goal:** Launch core wallet + gigs + savings

| Module | Features Included |
|--------|-------------------|
| Wallet | Virtual account, P2P, withdrawals, history |
| Gigs | Post/browse gigs, proposals, escrow, reviews |
| Savings | Create/join circles, contributions, payouts |
| Auth | OTP login, PIN for transactions |

### V1.0 (Month 4-6)

**Goal:** Add credit, improve UX

| Module | Features Added |
|--------|----------------|
| Credit | HustleScore display, credit builder loans |
| Wallet | Bill payments, QR codes |
| Gigs | Skills verification, categories expansion |
| Savings | Auto-contribution, early withdrawal |

### V1.5 (Month 7-9)

**Goal:** Scale credit, add engagement

| Module | Features Added |
|--------|----------------|
| Credit | Full loan products, score tips |
| Gigs | Milestone payments, dispute resolution |
| Savings | Circle insurance, target savings |
| Platform | Referral program, notifications |

### V2.0 (Month 10-12)

**Goal:** B2B features, advanced credit

| Module | Features Added |
|--------|----------------|
| Credit | Higher limits, longer terms |
| Gigs | Business accounts, bulk posting |
| Savings | Corporate circles, employee savings |
| Platform | API for partners, analytics dashboard |

### Future Roadmap (2027+)

| Feature | Timeline | Description |
|---------|----------|-------------|
| Insurance | Q2 2027 | Life, health, device insurance |
| Investments | Q3 2027 | Mutual funds, treasury bills |
| BNPL | Q4 2027 | Buy-now-pay-later at merchants |
| B2B Payroll | Q1 2028 | Workforce payment management |
| Remittance | Q2 2028 | International money transfer |

---

**Previous Section:** [02_PROBLEM_ANALYSIS.md](./02_PROBLEM_ANALYSIS.md)
**Next Section:** [04_MARKET_ANALYSIS.md](./04_MARKET_ANALYSIS.md)
