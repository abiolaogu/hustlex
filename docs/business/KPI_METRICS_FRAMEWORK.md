# HustleX KPI & Metrics Framework

> Comprehensive Performance Measurement System

---

## Executive Summary

This framework defines the key performance indicators (KPIs) and metrics that drive HustleX's business success. It establishes measurement standards, targets, and accountability for all business functions.

**North Star Metric:** Weekly Active Transacting Users (WATU)

---

## 1. Metric Hierarchy

### 1.1 Metric Pyramid

```
                          ┌───────────────┐
                          │   NORTH STAR  │
                          │     WATU      │
                          └───────┬───────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
        ┌─────▼─────┐       ┌─────▼─────┐       ┌─────▼─────┐
        │  GROWTH   │       │ ENGAGEMENT│       │ REVENUE   │
        │  Metrics  │       │  Metrics  │       │  Metrics  │
        └─────┬─────┘       └─────┬─────┘       └─────┬─────┘
              │                   │                   │
    ┌─────────┼─────────┐ ┌───────┼───────┐ ┌─────────┼─────────┐
    │         │         │ │       │       │ │         │         │
┌───▼───┐ ┌───▼───┐ ┌───▼─▼─┐ ┌───▼───┐ ┌─▼─────┐ ┌───▼───┐ ┌───▼───┐
│ Acq.  │ │ Activ.│ │Retent.│ │Feature│ │Transac│ │ GMV   │ │ Unit  │
│Metrics│ │Metrics│ │Metrics│ │Usage  │ │tion   │ │       │ │ Econ  │
└───────┘ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘
```

### 1.2 North Star Metric

**Weekly Active Transacting Users (WATU)**

| Attribute | Definition |
|-----------|------------|
| Definition | Users who complete at least one financial transaction in a week |
| Calculation | Count of unique users with ≥1 transaction in trailing 7 days |
| Transactions Include | Gig payment, savings contribution, loan repayment, wallet transfer, withdrawal |
| Exclude | Profile views, browsing, messaging only |

**Why WATU?**
- Combines engagement AND monetization
- Leading indicator of revenue
- Reflects core product value
- Weekly timeframe balances signal vs. noise

**Targets:**
| Period | Target WATU | % of MAU |
|--------|-------------|----------|
| Month 3 | 15,000 | 30% |
| Month 6 | 60,000 | 35% |
| Month 12 | 200,000 | 40% |

---

## 2. Growth Metrics

### 2.1 Acquisition Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **New Signups** | Users completing registration | Count | 2,000/day |
| **App Downloads** | Total app installs | Count | 3,000/day |
| **Signup Conversion** | Downloads → Signups | Signups / Downloads | 65% |
| **CAC** | Cost to acquire user | Marketing Spend / New Users | ₦800 |
| **Organic %** | Non-paid signups | Organic Signups / Total | 40% |
| **Channel Attribution** | Signups by channel | Count by source | Tracked |

### 2.2 Activation Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **Activation Rate** | Signups completing key action | Activated / Signups | 50% |
| **Time to Activate** | Days to first transaction | Median days | < 3 days |
| **Profile Completion** | Users with full profile | Complete / Total | 70% |
| **KYC Verification** | Users BVN verified | Verified / Total | 60% |

**Activation Definition:**
A user is "activated" when they complete ONE of:
- First gig payment (received or sent)
- First savings contribution
- First loan disbursement
- First wallet funding ≥ ₦5,000

### 2.3 Retention Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **D1 Retention** | Return next day | D1 Active / D0 Signups | 60% |
| **D7 Retention** | Return within week | D7 Active / D0 Signups | 45% |
| **D30 Retention** | Return within month | D30 Active / D0 Signups | 30% |
| **D90 Retention** | Return within quarter | D90 Active / D0 Signups | 20% |
| **Monthly Churn** | Users not returning | (MAU_t-1 - Returning) / MAU_t-1 | < 8% |
| **Resurrection Rate** | Churned users returning | Resurrected / Churned | 15% |

### 2.4 Viral Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **K-Factor** | Viral coefficient | Invites × Conversion Rate | 1.5 |
| **Referral Rate** | Users who refer | Referrers / Total Users | 25% |
| **Invites Sent** | Circle/referral invites | Count | 10K/day |
| **Invite Conversion** | Invites → Signups | Signups / Invites | 20% |
| **Avg Referrals/User** | Referrals per referrer | Total Referred / Referrers | 3.2 |

---

## 3. Engagement Metrics

### 3.1 Activity Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **DAU** | Daily Active Users | Unique daily users | 150K |
| **WAU** | Weekly Active Users | Unique weekly users | 300K |
| **MAU** | Monthly Active Users | Unique monthly users | 500K |
| **DAU/MAU** | Stickiness ratio | DAU / MAU | 35% |
| **Sessions/User/Day** | Daily engagement | Total Sessions / DAU | 2.5 |
| **Session Duration** | Time in app | Avg minutes/session | 5 min |

### 3.2 Feature Adoption

| Feature | Metric | Target Adoption |
|---------|--------|-----------------|
| **Gigs** | Users browsing gigs/week | 60% of MAU |
| **Gigs** | Users with active gig | 20% of MAU |
| **Savings** | Users in savings circle | 40% of MAU |
| **Savings** | Active contributors/month | 80% of circle members |
| **Credit** | Users with credit score | 50% of 90-day users |
| **Credit** | Users with active loan | 15% of eligible |
| **Wallet** | Wallet funding/month | 70% of MAU |

### 3.3 Content/Notification Metrics

| Metric | Definition | Target |
|--------|------------|--------|
| Push Notification CTR | Opens / Sent | 8% |
| Email Open Rate | Opens / Delivered | 25% |
| SMS Delivery Rate | Delivered / Sent | 98% |
| In-App Message CTR | Clicks / Impressions | 15% |

---

## 4. Revenue Metrics

### 4.1 Transaction Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **GMV** | Gross Merchandise Value | Total transaction value | ₦5B/month |
| **Transaction Count** | Total transactions | Count | 1M/month |
| **Avg Transaction** | Average value | GMV / Transactions | ₦5,000 |
| **Transactions/User** | Per-user frequency | Transactions / MAU | 8/month |

### 4.2 Revenue Metrics

| Metric | Definition | Formula | Target |
|--------|------------|---------|--------|
| **Gross Revenue** | Total revenue | Sum of all revenue | ₦500M/month |
| **Net Revenue** | After payment costs | Gross - Processing Fees | ₦480M/month |
| **ARPU** | Avg Revenue Per User | Revenue / MAU | ₦1,000/month |
| **ARPPU** | Avg Revenue Per Paying User | Revenue / Paying Users | ₦2,500/month |
| **Take Rate** | Platform revenue % | Revenue / GMV | 8% |

### 4.3 Revenue by Product

| Product | Metrics | Target Mix |
|---------|---------|------------|
| **Gigs** | Fee revenue, GMV | 35% |
| **Savings** | Circle fees, late fees | 10% |
| **Credit** | Interest, origination | 45% |
| **Wallet** | Funding, withdrawal, bills | 10% |

---

## 5. Unit Economics

### 5.1 Customer Lifetime Value (LTV)

| Component | Formula | Value |
|-----------|---------|-------|
| Avg Monthly Revenue | ARPU | ₦1,000 |
| Gross Margin | (Revenue - COGS) / Revenue | 73% |
| Avg Lifespan | 1 / Monthly Churn | 40 months |
| **LTV** | ARPU × Margin × Lifespan | **₦29,200** |

### 5.2 Customer Acquisition Cost (CAC)

| Component | Formula | Value |
|-----------|---------|-------|
| Marketing Spend | Monthly budget | ₦25M |
| New Users | Monthly acquisitions | 40,000 |
| **CAC** | Spend / New Users | **₦625** |

### 5.3 Key Ratios

| Ratio | Formula | Target | Current |
|-------|---------|--------|---------|
| **LTV/CAC** | LTV / CAC | > 3x | 46x |
| **Payback Period** | CAC / (ARPU × Margin) | < 12 months | 0.9 months |
| **Contribution Margin** | (Revenue - Variable Costs) / Revenue | > 40% | 45% |

---

## 6. Product-Specific Metrics

### 6.1 Gig Marketplace

| Metric | Definition | Target |
|--------|------------|--------|
| **Gigs Posted** | New gigs/day | 500 |
| **Gig Fill Rate** | Gigs with accepted proposal | 70% |
| **Avg Time to Hire** | Days to accept proposal | 3 days |
| **Proposal Rate** | Proposals per gig | 5 |
| **Completion Rate** | Completed / Started gigs | 90% |
| **Client Repeat Rate** | Clients posting 2+ gigs | 40% |
| **Freelancer Repeat Rate** | Freelancers with 2+ gigs | 60% |
| **Dispute Rate** | Disputes / Completed gigs | < 2% |
| **Platform Fee Take** | Revenue / GMV | 10% |

### 6.2 Savings Circles

| Metric | Definition | Target |
|--------|------------|--------|
| **Circles Created** | New circles/week | 500 |
| **Avg Circle Size** | Members per circle | 8 |
| **Circle Completion** | Circles completing cycle | 85% |
| **Contribution Rate** | On-time contributions | 95% |
| **Late Contribution %** | Late / Total contributions | < 10% |
| **Default Rate** | Members defaulting | < 3% |
| **Payout Satisfaction** | Payout delivered on time | 98% |
| **Circle Renewal** | Circles starting new cycle | 60% |

### 6.3 Credit Products

| Metric | Definition | Target |
|--------|------------|--------|
| **Loan Applications** | Applications/day | 1,000 |
| **Approval Rate** | Approved / Applied | 60% |
| **Disbursement Rate** | Disbursed / Approved | 95% |
| **Avg Loan Size** | Average disbursement | ₦40,000 |
| **Avg Tenure** | Loan duration | 4 months |
| **Interest Rate** | Monthly interest | 5% |
| **NPL Rate (30+)** | 30+ days past due / Book | < 5% |
| **NPL Rate (90+)** | 90+ days past due / Book | < 2% |
| **Collections Rate** | Collected / Due | 95% |
| **Write-off Rate** | Written off / Book | < 1.5% |

### 6.4 Wallet

| Metric | Definition | Target |
|--------|------------|--------|
| **Wallet Funding** | Deposits/day | ₦100M |
| **Avg Balance** | Average wallet balance | ₦15,000 |
| **Withdrawals** | Withdrawals/day | ₦80M |
| **P2P Transfers** | Transfers between users | ₦50M/day |
| **Bill Payments** | Bills paid/day | 5,000 |
| **Float Balance** | Total platform float | ₦500M |

---

## 7. Operational Metrics

### 7.1 Platform Performance

| Metric | Definition | Target |
|--------|------------|--------|
| **Uptime** | System availability | 99.9% |
| **API Latency (p95)** | 95th percentile response | < 200ms |
| **Error Rate** | Failed requests | < 0.1% |
| **App Crash Rate** | Crashes / Sessions | < 0.1% |
| **Load Time** | App launch to usable | < 3s |

### 7.2 Support Metrics

| Metric | Definition | Target |
|--------|------------|--------|
| **Ticket Volume** | Support tickets/day | < 500 |
| **Tickets/User** | Monthly tickets per MAU | < 0.5% |
| **First Response** | Time to first response | < 2 hours |
| **Resolution Time** | Time to resolve | < 24 hours |
| **CSAT** | Customer satisfaction | > 4.5/5 |
| **NPS** | Net Promoter Score | > 50 |

---

## 8. Dashboard & Reporting

### 8.1 Executive Dashboard (Daily)

| Section | Metrics |
|---------|---------|
| **Hero** | WATU, DAU, Revenue (24h) |
| **Growth** | New signups, CAC, Activation |
| **Engagement** | DAU/MAU, Sessions, Feature usage |
| **Revenue** | GMV, Revenue, ARPU |
| **Health** | NPL, Uptime, NPS |

### 8.2 Weekly Business Review

| Section | Metrics | Owner |
|---------|---------|-------|
| Growth | User growth, CAC, channels | Growth |
| Product | Feature adoption, retention | Product |
| Revenue | GMV, revenue, margins | Finance |
| Credit | NPL, disbursements, collections | Risk |
| Operations | Support, uptime, issues | Ops |

### 8.3 Monthly Board Report

| Section | Metrics |
|---------|---------|
| Financial | P&L, cash position, burn |
| Growth | MAU, WATU, LTV/CAC |
| Product | Key feature metrics, roadmap |
| Risk | Credit performance, fraud |
| Team | Headcount, key hires |

---

## 9. Targets by Period

### 9.1 Year 1 Quarterly Targets

| Metric | Q1 | Q2 | Q3 | Q4 |
|--------|-----|-----|-----|-----|
| MAU | 50K | 150K | 300K | 500K |
| WATU | 15K | 50K | 110K | 200K |
| GMV | ₦500M | ₦2B | ₦5B | ₦10B |
| Revenue | ₦50M | ₦200M | ₦500M | ₦1B |
| CAC | ₦1,000 | ₦800 | ₦700 | ₦600 |
| D30 Retention | 25% | 28% | 30% | 35% |
| NPL Rate | 10% | 8% | 6% | 5% |

### 9.2 Year 1-5 Annual Targets

| Metric | Y1 | Y2 | Y3 | Y4 | Y5 |
|--------|-----|-----|-----|-----|-----|
| MAU | 500K | 2M | 5M | 10M | 18M |
| WATU | 200K | 800K | 2.5M | 5M | 9M |
| GMV | ₦15B | ₦75B | ₦250B | ₦600B | ₦1.2T |
| Revenue | ₦1.5B | ₦6B | ₦24B | ₦60B | ₦120B |
| Net Income | (₦200M) | ₦500M | ₦5B | ₦15B | ₦35B |

---

## 10. Metric Ownership

### 10.1 RACI Matrix

| Metric Category | Responsible | Accountable | Consulted | Informed |
|-----------------|-------------|-------------|-----------|----------|
| Growth/Acquisition | Growth Lead | CMO | Product | Exec |
| Engagement | Product Lead | CPO | Growth | Exec |
| Revenue | Finance Lead | CFO | Product | Board |
| Credit | Risk Lead | CRO | Finance | Exec |
| Operations | Ops Lead | COO | Tech | Exec |
| Technology | Tech Lead | CTO | Product | Exec |

### 10.2 Review Cadence

| Level | Frequency | Participants | Focus |
|-------|-----------|--------------|-------|
| Team | Daily | Team | Tactical metrics |
| Department | Weekly | Leads | Functional metrics |
| Leadership | Weekly | C-suite | Business metrics |
| Board | Monthly | Board + CEO | Strategic metrics |

---

**Document Control**

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | January 2025 | BillyRonks Analytics | Initial version |

---

*© 2025 BillyRonks Global Limited. All rights reserved.*
