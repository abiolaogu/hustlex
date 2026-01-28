# 13. KPIs and Monitoring

---

## North Star Metric

### Definition

> **Monthly Active Transacting Users (MATU)**
>
> The number of unique users who complete at least one value-generating transaction (gig payment, savings contribution, loan repayment, or bill payment) in a calendar month.

### Why MATU?

| Reason | Explanation |
|--------|-------------|
| **Reflects engagement** | Users must actively use the platform, not just register |
| **Tied to revenue** | Every transaction generates revenue |
| **Leading indicator** | MATU growth predicts future revenue |
| **Actionable** | Teams can influence through product and marketing |
| **Comparable** | Industry-standard metric for fintechs |

### MATU Targets

| Period | Target | Notes |
|--------|--------|-------|
| Month 3 | 5,000 | Launch milestone |
| Month 6 | 15,000 | Product-market fit signal |
| Month 12 | 60,000 | Year 1 target |
| Month 24 | 300,000 | Year 2 target |
| Month 36 | 800,000 | Year 3 target |

---

## KPI Tree

### Revenue Decomposition

```
                              REVENUE
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                  │
              ▼                  ▼                  ▼
        GIG REVENUE      LENDING REVENUE    OTHER REVENUE
              │                  │                  │
       ┌──────┴──────┐    ┌──────┴──────┐    ┌──────┴──────┐
       │             │    │             │    │             │
       ▼             ▼    ▼             ▼    ▼             ▼
    Gig GMV    Take Rate  Loan Book  Interest  Bills    Float
       │             │        │      Rate       Savings   Subs
       │             │        │                 API
   ┌───┴───┐         │    ┌───┴───┐
   │       │         │    │       │
   ▼       ▼         │    ▼       ▼
# Gigs  Avg Gig      │  # Loans  Avg Loan
Posted  Value        │  Active   Size
                     │
              ┌──────┴──────┐
              │             │
              ▼             ▼
         Default Rate  Collection Rate
```

### User Funnel Decomposition

```
                           REGISTERED USERS
                                 │
                                 │ Activation Rate
                                 ▼
                           ACTIVATED USERS
                           (KYC + 1st Tx)
                                 │
                                 │ Retention Rate
                                 ▼
                          MONTHLY ACTIVE USERS
                                 │
                                 │ Transaction Rate
                                 ▼
                    MONTHLY ACTIVE TRANSACTING USERS
                                 │
                    ┌────────────┼────────────┐
                    │            │            │
                    ▼            ▼            ▼
               Gig Users   Savings Users  Credit Users
                    │            │            │
                    ▼            ▼            ▼
              Tx/User      Contrib/User   Loan/User
```

---

## KPI Framework

### Tier 1: North Star & Strategic KPIs

| KPI | Definition | Target Y1 | Owner |
|-----|------------|-----------|-------|
| **MATU** | Monthly Active Transacting Users | 60,000 | CEO |
| **Revenue** | Total monthly revenue | ₦4M | CFO |
| **Gross Margin** | (Revenue - COGS) / Revenue | 55% | CFO |
| **NPS** | Net Promoter Score | >40 | COO |

### Tier 2: Functional KPIs

#### Growth KPIs (CMO)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| New User Signups | Users completing registration | 12,500/mo | Daily |
| Activation Rate | % completing KYC + first transaction | 60% | Weekly |
| CAC | Total marketing spend / new users | ₦1,500 | Monthly |
| Organic % | % of new users from organic/referral | 40% | Monthly |
| Referral Rate | % of users who refer others | 30% | Monthly |

#### Engagement KPIs (Product)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| DAU/MAU | Daily active / Monthly active | 30% | Daily |
| Sessions per User | Avg app opens per user per week | 5 | Weekly |
| Transactions per User | Avg transactions per MATU | 8/mo | Monthly |
| Feature Adoption | % using each major feature | Varies | Monthly |
| App Rating | iOS/Android store rating | >4.5 | Weekly |

#### Retention KPIs (Product/Ops)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| M1 Retention | % active in month 2 | 50% | Monthly |
| M3 Retention | % active in month 4 | 35% | Monthly |
| M6 Retention | % active in month 7 | 25% | Monthly |
| Churn Rate | % becoming inactive | <5%/mo | Monthly |
| Resurrection Rate | % returning after 30+ days inactive | 10% | Monthly |

#### Gig KPIs (Product)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Gigs Posted | New gigs posted per month | 5,000 | Weekly |
| Gig Completion Rate | % of posted gigs completed | 70% | Weekly |
| Avg Gig Value | Average gig transaction size | ₦15,000 | Weekly |
| Time to First Proposal | Avg hours to first proposal | <24 hrs | Weekly |
| Dispute Rate | % of gigs with disputes | <3% | Weekly |

#### Savings KPIs (Product)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Active Circles | Circles with activity this month | 2,000 | Weekly |
| Contribution Rate | % of expected contributions made | 90% | Weekly |
| Circle Completion Rate | % of circles completing full cycle | 85% | Monthly |
| Avg Circle Size | Members per active circle | 10 | Monthly |
| Circle GMV | Total monthly circle contributions | ₦100M | Weekly |

#### Credit KPIs (CFO/Risk)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Loans Disbursed | Number of new loans | 2,000/mo | Daily |
| Loan Book Size | Outstanding loan principal | ₦20M | Daily |
| Avg Loan Size | Average disbursement | ₦20,000 | Weekly |
| Approval Rate | % of applications approved | 40% | Weekly |
| Default Rate (30 days) | % of loans 30+ days overdue | <8% | Daily |
| Default Rate (90 days) | % of loans 90+ days overdue | <5% | Weekly |
| Collection Rate | % of due payments collected | 92% | Daily |
| Net Interest Margin | (Interest income - losses) / loan book | 25% | Monthly |

#### Operations KPIs (COO)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Support Response Time | Avg time to first response | <2 hrs | Real-time |
| Support Resolution Time | Avg time to resolution | <24 hrs | Daily |
| CSAT | Customer satisfaction score | >4.5/5 | Weekly |
| Ticket Volume | Support tickets per 1K MAU | <50 | Weekly |
| Agent Productivity | Signups per agent per month | 50 | Weekly |

#### Technology KPIs (CTO)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Uptime | Platform availability | 99.9% | Real-time |
| API Response Time | 95th percentile latency | <200ms | Real-time |
| Error Rate | % of requests with errors | <0.1% | Real-time |
| Deployment Frequency | Production deploys per week | 5+ | Weekly |
| Incident Count (P1/P2) | Critical/high incidents | 0 | Weekly |

#### Financial KPIs (CFO)

| KPI | Definition | Target | Frequency |
|-----|------------|--------|-----------|
| Monthly Burn | Cash spent per month | <$40K | Monthly |
| Runway | Months of cash remaining | >12 | Monthly |
| Revenue Growth | MoM revenue growth | 15%+ | Monthly |
| Unit Economics | LTV/CAC ratio | >3x | Quarterly |
| Gross Margin | Gross profit / revenue | 55%+ | Monthly |

---

## Dashboard Framework

### Executive Dashboard

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      HUSTLEX EXECUTIVE DASHBOARD                             │
│                           January 2026                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  NORTH STAR                    │  REVENUE                                    │
│  ┌─────────────────────────┐   │  ┌─────────────────────────┐               │
│  │ MATU: 5,234             │   │  │ Revenue: ₦4.2M         │               │
│  │ Target: 5,000  ✓        │   │  │ Target: ₦4M  ✓         │               │
│  │ MoM: +23%               │   │  │ MoM: +18%              │               │
│  └─────────────────────────┘   │  └─────────────────────────┘               │
│                                │                                             │
│  USER FUNNEL                   │  FINANCIAL HEALTH                           │
│  ┌─────────────────────────┐   │  ┌─────────────────────────┐               │
│  │ Registered: 25,000      │   │  │ Burn: $35K             │               │
│  │ Activated: 15,000 (60%) │   │  │ Runway: 14 months      │               │
│  │ MAU: 10,000 (40%)       │   │  │ Gross Margin: 54%      │               │
│  │ MATU: 5,234 (21%)       │   │  │ LTV/CAC: 2.1x          │               │
│  └─────────────────────────┘   │  └─────────────────────────┘               │
│                                │                                             │
│  PRODUCT HEALTH                │  RISK INDICATORS                            │
│  ┌─────────────────────────┐   │  ┌─────────────────────────┐               │
│  │ NPS: 42                 │   │  │ Default Rate: 6.5%     │               │
│  │ App Rating: 4.3         │   │  │ Fraud Rate: 0.3%       │               │
│  │ DAU/MAU: 28%            │   │  │ Uptime: 99.95%         │               │
│  │ Retention M1: 52%       │   │  │ P1 Incidents: 0        │               │
│  └─────────────────────────┘   │  └─────────────────────────┘               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Product Dashboard

| Section | Metrics Shown |
|---------|---------------|
| Funnel | Daily signups, activations, first transactions |
| Engagement | DAU, sessions, feature usage heatmap |
| Retention | Cohort curves, churn analysis |
| Features | Gig, savings, credit metrics |

### Growth Dashboard

| Section | Metrics Shown |
|---------|---------------|
| Acquisition | Signups by channel, CAC by channel |
| Activation | Funnel conversion, drop-off points |
| Campaigns | Active campaigns, performance |
| Referrals | Referral rate, viral coefficient |

### Risk Dashboard

| Section | Metrics Shown |
|---------|---------------|
| Credit | Default rates, PAR, collection rates |
| Fraud | Fraud rate, detection rate, cases |
| Security | Vulnerabilities, incidents, alerts |
| Compliance | Open issues, deadlines |

---

## Reporting Cadence

### Daily Standup Metrics

| Metric | Owner | Action Threshold |
|--------|-------|------------------|
| Signups yesterday | Growth | <80% of daily target |
| Transactions yesterday | Product | <80% of daily avg |
| Support tickets open | Ops | >100 tickets |
| System errors | Eng | Any P1/P2 |
| Default rate (rolling 7d) | Risk | >10% |

### Weekly Business Review

**Attendees:** Leadership team
**Duration:** 60 minutes
**Day:** Monday

**Agenda:**
1. North star and Tier 1 KPIs (10 min)
2. Funnel and growth review (15 min)
3. Product and engagement (10 min)
4. Credit and risk (10 min)
5. Blockers and decisions (15 min)

### Monthly Board Report

**Sections:**
1. Executive summary (1 page)
2. KPI dashboard (2 pages)
3. Financial update (2 pages)
4. Risk and compliance (1 page)
5. Strategic initiatives (1 page)
6. Ask/decisions needed (1 page)

### Quarterly Business Review

**Attendees:** Full team + board
**Duration:** Half day

**Agenda:**
1. Quarter in review
2. OKR scoring
3. Financial deep dive
4. Strategic discussion
5. Next quarter OKRs
6. Team recognition

---

## Alert Framework

### Alert Severity Levels

| Level | Response | Examples |
|-------|----------|----------|
| **Critical** | Immediate action required | Platform down, data breach |
| **High** | Action within 1 hour | Major feature broken, fraud spike |
| **Medium** | Action within 24 hours | KPI significantly below target |
| **Low** | Review in next standup | Minor metrics deviation |

### Automated Alerts

| Alert | Condition | Channel | Owner |
|-------|-----------|---------|-------|
| Platform down | Uptime <99% for 5 min | PagerDuty, Slack #incidents | CTO |
| Error spike | Error rate >1% for 10 min | Slack #engineering | CTO |
| Default spike | Daily default rate >15% | Slack #risk, Email | CFO |
| Fraud alert | Fraud score >0.9 | Slack #fraud | CTO |
| Support SLA breach | Response >4 hours | Slack #support | COO |
| Signup drop | Daily signups <50% of avg | Slack #growth | CMO |
| Cash alert | Runway <6 months | Email to CEO, CFO | CFO |

---

## Data Infrastructure

### Data Stack

```
┌─────────────────────────────────────────────────────────────────────┐
│                        DATA SOURCES                                  │
├────────────┬────────────┬────────────┬────────────┬─────────────────┤
│ PostgreSQL │   Redis    │  Paystack  │   Termii   │  App Analytics  │
│  (Primary) │  (Events)  │ (Payments) │   (SMS)    │   (Amplitude)   │
└─────┬──────┴─────┬──────┴─────┬──────┴─────┬──────┴────────┬────────┘
      │            │            │            │               │
      └────────────┴────────────┴────────────┴───────────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    ETL Pipeline     │
                    │  (Airbyte / Fivetran)│
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │   Data Warehouse    │
                    │  (BigQuery / Redshift)│
                    └──────────┬──────────┘
                               │
              ┌────────────────┼────────────────┐
              │                │                │
              ▼                ▼                ▼
     ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
     │   Grafana   │  │   Metabase  │  │    dbt      │
     │ (Real-time) │  │(BI/Reports) │  │ (Transform) │
     └─────────────┘  └─────────────┘  └─────────────┘
```

### Data Governance

| Principle | Implementation |
|-----------|----------------|
| Single source of truth | Warehouse is authoritative |
| Data ownership | Each metric has defined owner |
| Data quality | Automated tests, alerts on anomalies |
| Access control | Role-based, need-to-know |
| Audit trail | All queries logged |

---

**Previous Section:** [12_RISK_MANAGEMENT.md](./12_RISK_MANAGEMENT.md)
**Next Section:** [14_TEAM_AND_GOVERNANCE.md](./14_TEAM_AND_GOVERNANCE.md)
