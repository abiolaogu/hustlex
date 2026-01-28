# 12. Risk Management

---

## Risk Management Framework

### Risk Governance Structure

```
┌─────────────────────────────────────────────────────────────────────┐
│                         BOARD OF DIRECTORS                           │
│                    (Ultimate Risk Oversight)                         │
└─────────────────────────────────────┬───────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────┐
│                              CEO                                     │
│                    (Enterprise Risk Owner)                           │
└─────────────────────────────────────┬───────────────────────────────┘
                                      │
        ┌─────────────────────────────┼─────────────────────────────┐
        │                             │                             │
        ▼                             ▼                             ▼
┌───────────────┐           ┌───────────────┐           ┌───────────────┐
│   CTO         │           │   CFO         │           │   COO         │
│               │           │               │           │               │
│ • Technology  │           │ • Financial   │           │ • Operational │
│ • Security    │           │ • Credit      │           │ • Fraud       │
│ • Data        │           │ • Compliance  │           │ • Reputational│
└───────────────┘           └───────────────┘           └───────────────┘
```

### Risk Categories

| Category | Description | Primary Owner |
|----------|-------------|---------------|
| **Strategic** | Market, competitive, business model risks | CEO |
| **Financial** | Credit, liquidity, FX risks | CFO |
| **Operational** | Process, fraud, vendor risks | COO |
| **Technology** | Security, availability, data risks | CTO |
| **Regulatory** | Compliance, licensing risks | CEO/CFO |
| **Reputational** | Brand, trust, PR risks | CEO/CMO |

---

## Risk Register

### Top 20 Risks

| # | Risk | Category | Likelihood | Impact | Score | Mitigation | Owner |
|---|------|----------|------------|--------|-------|------------|-------|
| 1 | Credit defaults exceed projections | Financial | High | Critical | **20** | Conservative underwriting, collections process | CFO |
| 2 | Regulatory shutdown/restrictions | Regulatory | Medium | Critical | **16** | Proactive engagement, sandbox, compliance | CEO |
| 3 | Data breach/cybersecurity incident | Technology | Medium | Critical | **16** | Security hardening, SOC2, pentesting | CTO |
| 4 | Fraud losses exceed threshold | Operational | High | High | **16** | ML fraud detection, limits, insurance | CTO |
| 5 | Funding gap (unable to raise) | Financial | Medium | Critical | **16** | Multiple investor tracks, revenue focus | CEO |
| 6 | Key person departure | Operational | Medium | High | **12** | ESOP, culture, succession planning | CEO |
| 7 | Payment partner failure | Operational | Low | Critical | **12** | Multiple providers, monitoring | CTO |
| 8 | Competitor replication | Strategic | High | Medium | **12** | Speed, moat building, customer lock-in | CEO |
| 9 | FX volatility impact | Financial | High | Medium | **12** | Natural hedge, USD reserves | CFO |
| 10 | Agent fraud/misconduct | Operational | High | Medium | **12** | KYC, monitoring, bonds, training | COO |
| 11 | Economic downturn | Strategic | Medium | High | **12** | Diversified segments, essential services | CEO |
| 12 | Technical platform outage | Technology | Medium | High | **12** | Multi-AZ, DR plan, monitoring | CTO |
| 13 | Slow user adoption | Strategic | Medium | High | **12** | Iteration, channel diversification | CMO |
| 14 | Interest rate changes (policy) | Financial | Medium | Medium | **9** | Flexible pricing, hedging | CFO |
| 15 | Talent acquisition challenges | Operational | Medium | Medium | **9** | Competitive comp, culture, remote | CEO |
| 16 | Ajo circle default cascade | Financial | Low | High | **8** | Circle insurance, limits, monitoring | CFO |
| 17 | Negative press/PR crisis | Reputational | Low | High | **8** | PR plan, transparency, crisis response | CEO |
| 18 | Third-party data breach | Technology | Low | High | **8** | Vendor assessment, contracts, monitoring | CTO |
| 19 | SMS/notification failure | Operational | Medium | Medium | **9** | Multiple providers, fallback channels | CTO |
| 20 | Legal/IP disputes | Regulatory | Low | Medium | **6** | Legal counsel, IP protection | CEO |

### Risk Scoring Matrix

```
                           IMPACT
                    Low   Medium   High   Critical
                     │      │       │        │
         High    │   6  │   9   │  12   │   16   │
                 ├──────┼───────┼───────┼────────┤
LIKELIHOOD Medium│   4  │   6   │   9   │   12   │
                 ├──────┼───────┼───────┼────────┤
         Low     │   2  │   4   │   6   │    8   │
                 └──────┴───────┴───────┴────────┘

         Risk Score: 1-4 = Low, 5-9 = Medium, 10-16 = High, 17-25 = Critical
```

---

## Detailed Risk Analysis

### Risk 1: Credit Defaults Exceed Projections

**Description:** Loan defaults exceed the projected 5-8% range, threatening profitability and capital.

| Aspect | Detail |
|--------|--------|
| **Likelihood** | High (3/4) |
| **Impact** | Critical (5/5) |
| **Risk Score** | 20 (Critical) |
| **Risk Appetite** | Mitigate |

**Root Causes:**
- Economic downturn affecting borrowers
- Inadequate underwriting models
- Fraud in loan applications
- Insufficient collections process
- Segment concentration

**Mitigation Measures:**

| Measure | Description | Status |
|---------|-------------|--------|
| Conservative underwriting | Start with small limits, increase with history | Planned |
| Platform data scoring | Use gig/savings data for credit decisions | In development |
| Real-time monitoring | Daily default rate tracking, alerts | Planned |
| Collections process | Multi-channel, ethical collections | Planned |
| Loan loss reserves | 10% provision from day one | Planned |
| Portfolio diversification | Spread across segments, geographies | Planned |
| Credit insurance | Explore credit guarantee schemes | Research |

**Key Metrics:**
- Default rate (target: <8%)
- Days past due distribution
- Recovery rate
- Portfolio at risk (PAR30)

**Escalation Triggers:**
- Default rate >10% for 2 consecutive months
- PAR30 >15%
- Recovery rate <60%

---

### Risk 2: Regulatory Shutdown/Restrictions

**Description:** CBN or other regulator restricts or shuts down operations.

| Aspect | Detail |
|--------|--------|
| **Likelihood** | Medium (2/4) |
| **Impact** | Critical (5/5) |
| **Risk Score** | 16 (High) |
| **Risk Appetite** | Mitigate |

**Root Causes:**
- Operating without required licenses
- Non-compliance with regulations
- Industry-wide crackdowns
- Political interference
- Complaints from competitors

**Mitigation Measures:**

| Measure | Description | Status |
|---------|-------------|--------|
| CBN sandbox application | Enter regulatory sandbox early | In progress |
| Licensed partner use | Use Paystack (licensed) for payments | Active |
| Compliance officer | Dedicated compliance function | Planned (M6) |
| Regulatory engagement | Regular dialogue with CBN | Ongoing |
| Legal counsel | Retained fintech-specialist lawyers | Active |
| Industry association | Join fintech association | Planned |
| Documentation | Maintain compliance records | Ongoing |

**Key Metrics:**
- Regulatory filings on time (100%)
- Open regulatory issues (0)
- Sandbox status

**Escalation Triggers:**
- Any regulatory inquiry
- Industry-wide announcements
- License delays >60 days

---

### Risk 3: Data Breach/Cybersecurity Incident

**Description:** Unauthorized access to user data or systems.

| Aspect | Detail |
|--------|--------|
| **Likelihood** | Medium (2/4) |
| **Impact** | Critical (5/5) |
| **Risk Score** | 16 (High) |
| **Risk Appetite** | Mitigate |

**Root Causes:**
- Software vulnerabilities
- Weak authentication
- Social engineering
- Insider threats
- Third-party breaches

**Mitigation Measures:**

| Measure | Description | Status |
|---------|-------------|--------|
| Security hardening | Address identified vulnerabilities | In progress |
| Penetration testing | Quarterly pentests | Planned |
| SOC2 compliance | Achieve Type I certification | Planned (Y2) |
| Encryption | AES-256 at rest, TLS 1.3 in transit | Active |
| Access controls | RBAC, MFA for admin | Active |
| Security monitoring | Real-time alerts, SIEM | Planned |
| Incident response plan | Documented and tested | Draft |
| Cyber insurance | $500K coverage | Planned |

**Key Metrics:**
- Critical vulnerabilities (0)
- Time to patch (critical: <24 hours)
- Security audit findings

**Escalation Triggers:**
- Any unauthorized access detected
- Critical vulnerability discovered
- Third-party breach affecting us

---

### Risk 4: Fraud Losses Exceed Threshold

**Description:** Fraud from users, agents, or external actors exceeds acceptable levels.

| Aspect | Detail |
|--------|--------|
| **Likelihood** | High (3/4) |
| **Impact** | High (4/5) |
| **Risk Score** | 16 (High) |
| **Risk Appetite** | Mitigate |

**Fraud Types:**
- Identity fraud (fake accounts)
- Transaction fraud (unauthorized)
- Loan fraud (false applications)
- Agent fraud (collusion, theft)
- Internal fraud (employee)

**Mitigation Measures:**

| Measure | Description | Status |
|---------|-------------|--------|
| KYC verification | BVN, ID, selfie verification | Active |
| Device fingerprinting | Track device patterns | Planned |
| Transaction monitoring | Real-time fraud scoring | In development |
| Velocity limits | Transaction limits by tier | Active |
| Agent bonds | Require agent deposits | Planned |
| Segregation of duties | Separate approval/execution | Active |
| Fraud investigation | Dedicated fraud team | Planned (M6) |
| Fraud insurance | Coverage for losses | Planned |

**Key Metrics:**
- Fraud rate (target: <0.5% of GMV)
- Fraud detection rate (>95%)
- False positive rate (<5%)

**Escalation Triggers:**
- Fraud rate >1% for any month
- New fraud pattern detected
- Agent fraud incident

---

## Compliance Requirements

### Regulatory Landscape

| Regulator | Requirement | Status | Timeline |
|-----------|-------------|--------|----------|
| **CBN** | Payment Service Provider (PSP) | Via Paystack | Current |
| **CBN** | Regulatory Sandbox | Application ready | Q1 2026 |
| **CBN** | Mobile Money Operator (MMO) | Evaluate | Q4 2026 |
| **CBN** | Microfinance Bank (MFB) | Application | Q2 2027 |
| **NDPR** | Data Protection Registration | In progress | Q1 2026 |
| **NDPR** | Data Protection Impact Assessment | Planned | Q2 2026 |
| **FCCPC** | Consumer Protection Compliance | Ongoing | Current |
| **EFCC** | AML/CFT Compliance | Ongoing | Current |

### AML/CFT Program

| Component | Description | Status |
|-----------|-------------|--------|
| KYC Policy | Tiered verification requirements | Active |
| Transaction Monitoring | Automated suspicious activity detection | In development |
| SAR Filing | Process for suspicious activity reports | Documented |
| Sanctions Screening | Check against sanctions lists | Via partner |
| Record Keeping | 5-year transaction retention | Active |
| Training | Annual AML training for staff | Planned |

### Data Protection (NDPR)

| Requirement | Status |
|-------------|--------|
| Registration with NITDA | In progress |
| Data Protection Policy | Draft |
| Privacy Notice in App | Draft |
| Consent Management | In development |
| Data Processing Records | Planned |
| Data Protection Officer | To be appointed |
| Data Breach Procedure | Draft |

---

## Insurance Coverage

### Current/Planned Coverage

| Insurance Type | Coverage | Premium Est. | Status |
|----------------|----------|--------------|--------|
| **Professional Indemnity** | ₦100M | ₦2M/year | Planned |
| **Cyber Liability** | ₦50M | ₦1.5M/year | Planned |
| **Directors & Officers** | ₦50M | ₦1M/year | Planned |
| **Crime/Fidelity** | ₦20M | ₦500K/year | Planned |
| **Business Interruption** | ₦20M | ₦400K/year | Planned |
| **General Liability** | ₦10M | ₦200K/year | Planned |

**Total Annual Premium Estimate: ₦5.6M (~$3,500)**

### Insurance Gaps

| Gap | Risk | Plan to Address |
|-----|------|-----------------|
| Credit guarantee | Default losses | Explore DFI programs |
| Key person | Founder loss | Evaluate after Series A |
| Political risk | Regulatory action | Monitor availability |

---

## Business Continuity Planning

### Business Continuity Objectives

| Metric | Target |
|--------|--------|
| Recovery Time Objective (RTO) | 4 hours |
| Recovery Point Objective (RPO) | 1 hour |
| Maximum Tolerable Downtime | 24 hours |

### Disaster Recovery Plan

| Scenario | Impact | Recovery Plan |
|----------|--------|---------------|
| **Data center outage** | Platform unavailable | Failover to secondary region |
| **Database corruption** | Data loss | Restore from backup (<1 hour RPO) |
| **Payment partner down** | Transactions fail | Switch to backup provider |
| **DDoS attack** | Service degradation | CloudFlare protection, scaling |
| **Office inaccessible** | Team disruption | Remote work (already distributed) |
| **Key person unavailable** | Decision delay | Deputy assignments, documentation |

### Critical Systems

| System | RTO | Backup Strategy |
|--------|-----|-----------------|
| Production Database | 1 hour | Multi-AZ, daily backups |
| Application Servers | 15 min | Auto-scaling, multi-region |
| Payment Gateway | N/A | Multiple providers |
| Customer Data | 1 hour | Encrypted backups, offsite |
| Source Code | N/A | GitHub (distributed) |

---

## Risk Monitoring

### Risk Dashboard

| Risk Area | Metric | Threshold | Frequency |
|-----------|--------|-----------|-----------|
| **Credit** | Default rate | <8% | Daily |
| **Credit** | PAR30 | <15% | Weekly |
| **Fraud** | Fraud rate | <0.5% | Daily |
| **Fraud** | Detection rate | >95% | Weekly |
| **Security** | Critical vulns | 0 | Real-time |
| **Security** | Incidents | 0 | Real-time |
| **Ops** | Uptime | >99.9% | Real-time |
| **Ops** | Error rate | <0.1% | Real-time |
| **Compliance** | Open issues | 0 | Weekly |
| **Financial** | Runway | >12 months | Monthly |

### Reporting Cadence

| Report | Audience | Frequency | Owner |
|--------|----------|-----------|-------|
| Risk Dashboard | Leadership | Real-time | CTO |
| Risk Summary | CEO | Weekly | CFO |
| Board Risk Report | Board | Quarterly | CEO |
| Compliance Report | Board | Quarterly | Compliance |
| Audit Report | Board | Annual | External |

---

## Incident Response

### Incident Classification

| Severity | Definition | Response Time | Escalation |
|----------|------------|---------------|------------|
| **P1 - Critical** | Platform down, data breach, >₦1M fraud | 15 minutes | CEO immediately |
| **P2 - High** | Major feature down, security vulnerability | 1 hour | CTO within 1 hour |
| **P3 - Medium** | Minor feature issues, isolated fraud | 4 hours | Team lead |
| **P4 - Low** | Cosmetic issues, minor bugs | 24 hours | Engineer |

### Incident Response Process

```
DETECT ──▶ CLASSIFY ──▶ RESPOND ──▶ RESOLVE ──▶ REVIEW
   │           │           │           │           │
   │           │           │           │           │
Monitoring  Severity    Assemble    Fix issue   Post-mortem
Alerts      Assessment  Team        Communicate Document
User Report             Contain     Verify      Improve
```

### Communication Plan

| Stakeholder | P1 Notification | P2 Notification |
|-------------|-----------------|-----------------|
| CEO | Immediate | Within 1 hour |
| Board | Within 4 hours | Next board meeting |
| Regulators | As required | As required |
| Users | Within 24 hours | If affected |
| Press | Only if needed | No |

---

**Previous Section:** [11_FINANCIAL_MODEL.md](./11_FINANCIAL_MODEL.md)
**Next Section:** [13_KPIs_AND_MONITORING.md](./13_KPIs_AND_MONITORING.md)
