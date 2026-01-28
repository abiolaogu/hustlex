# HustleX Admin Manual

> Complete Guide for Platform Administrators

---

## Table of Contents

1. [Admin Overview](#1-admin-overview)
2. [Dashboard & Analytics](#2-dashboard--analytics)
3. [User Management](#3-user-management)
4. [Transaction Management](#4-transaction-management)
5. [Gig Management](#5-gig-management)
6. [Savings Circle Management](#6-savings-circle-management)
7. [Loan Management](#7-loan-management)
8. [Support & Disputes](#8-support--disputes)
9. [System Configuration](#9-system-configuration)
10. [Reports & Exports](#10-reports--exports)
11. [Security & Audit](#11-security--audit)
12. [Emergency Procedures](#12-emergency-procedures)

---

## 1. Admin Overview

### 1.1 Admin Roles

| Role | Permissions |
|------|-------------|
| **Super Admin** | Full system access, user management, configuration |
| **Finance Admin** | Transactions, withdrawals, reconciliation |
| **Support Admin** | User support, disputes, KYC verification |
| **Content Admin** | Gig moderation, reports, user flags |
| **Viewer** | Read-only access for monitoring |

### 1.2 Accessing Admin Panel

```
Production: https://admin.hustlex.app
Staging: https://admin-staging.hustlex.app
```

**Login Requirements:**
- Corporate email (@hustlex.app)
- 2FA enabled (Google Authenticator)
- VPN connection (for sensitive operations)

### 1.3 Admin Dashboard Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  HustleX Admin                    [Search] [Notifications] [👤] │
├──────────┬──────────────────────────────────────────────────────┤
│          │                                                       │
│ Dashboard│  Key Metrics (Today)                                 │
│ Users    │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│ Trans.   │  │ Users   │ │ Trans.  │ │ Volume  │ │ Active  │    │
│ Gigs     │  │ +1,234  │ │ 45,678  │ │₦50.2M   │ │ 12,456  │    │
│ Savings  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘    │
│ Credit   │                                                       │
│ Support  │  [Transaction Graph]     [User Growth Chart]         │
│ Config   │                                                       │
│ Reports  │  Recent Activity         Pending Actions             │
│          │  ├─ New user signup      ├─ 23 KYC reviews          │
│          │  ├─ Large withdrawal     ├─ 5 dispute cases         │
│          │  └─ Flagged transaction  └─ 12 loan applications    │
└──────────┴──────────────────────────────────────────────────────┘
```

---

## 2. Dashboard & Analytics

### 2.1 Key Performance Indicators

**Real-time Metrics:**
| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| Active Users | Users online now | < 100 (off-peak hours) |
| Transactions/min | Current throughput | > 1000 (investigate) |
| API Response Time | p95 latency | > 500ms (alert) |
| Error Rate | Failed requests | > 1% (critical) |

**Daily Metrics:**
| Metric | Description |
|--------|-------------|
| New Signups | Users registered today |
| DAU | Daily Active Users |
| Transaction Volume | Total ₦ processed |
| Gigs Posted | New gigs created |
| Savings Contributions | Circle payments |
| Loans Disbursed | New loans issued |

### 2.2 Analytics Reports

**User Analytics:**
- Signup funnel conversion
- User retention (D1, D7, D30)
- User segmentation by tier
- Geographic distribution

**Financial Analytics:**
- Revenue breakdown (gig fees, loan interest)
- Payment method distribution
- Withdrawal patterns
- Fraud indicators

**Engagement Analytics:**
- Feature usage heatmaps
- Session duration
- Most active hours
- App version distribution

### 2.3 Custom Reports

1. Navigate to **Reports** → **Custom Report**
2. Select metrics and dimensions
3. Set date range
4. Add filters
5. Generate or schedule

---

## 3. User Management

### 3.1 Viewing Users

**Search Users:**
- By phone number
- By email
- By name
- By user ID

**User Profile View:**
```
User: John Doe                              Status: Active
──────────────────────────────────────────────────────────
Phone: +2348012345678        Joined: Jan 15, 2024
Email: john@example.com      Last Active: 2 hours ago
Tier: Silver                 KYC: BVN Verified

Wallet Balance: ₦50,000.00
Escrow Balance: ₦10,000.00
Savings Balance: ₦25,000.00

Credit Score: 720 (Good)
Active Loans: 1 (₦35,000 remaining)

Gigs Posted: 5 | Completed: 12 | Rating: 4.8
Savings Circles: 2 active

[View Transactions] [View Gigs] [View Loans] [Flag User] [Suspend]
```

### 3.2 KYC Verification

**Pending Verifications Queue:**
1. Go to **Users** → **KYC Pending**
2. Review submitted documents
3. Cross-check with BVN/NIN records
4. Actions:
   - **Approve**: Verification successful
   - **Reject**: Document issues (specify reason)
   - **Request More**: Need additional documents

**Verification Checklist:**
- [ ] Name matches official records
- [ ] Date of birth matches
- [ ] Photo is clear and recent
- [ ] Document is not expired
- [ ] No signs of tampering

### 3.3 User Actions

| Action | Description | Requires |
|--------|-------------|----------|
| **Suspend** | Temporary disable account | Support Admin |
| **Ban** | Permanent account closure | Super Admin |
| **Reset PIN** | Clear PIN (user re-sets) | Support Admin |
| **Upgrade Tier** | Manually increase limits | Finance Admin |
| **Flag** | Mark for monitoring | Any Admin |
| **Unfreeze Wallet** | Remove restriction | Finance Admin |

### 3.4 User Tiers Management

**Tier Criteria:**
| Tier | Requirements | Limits |
|------|--------------|--------|
| Bronze | Phone verified | ₦50K daily |
| Silver | BVN verified | ₦200K daily |
| Gold | NIN + Address | ₦500K daily |
| Platinum | Full KYC + Review | ₦1M daily |

**Manual Tier Override:**
1. Go to user profile
2. Click **Edit Tier**
3. Select new tier
4. Provide justification
5. Submit (logged for audit)

---

## 4. Transaction Management

### 4.1 Transaction Monitoring

**Real-time Transaction Feed:**
```
Time       | User          | Type       | Amount    | Status
───────────┼───────────────┼────────────┼───────────┼─────────
10:45:32   | +2348012...   | Deposit    | ₦50,000   | Completed
10:45:30   | +2349087...   | Transfer   | ₦25,000   | Completed
10:45:28   | +2348045...   | Withdrawal | ₦100,000  | Pending ⚠️
10:45:25   | +2347012...   | Gig Payment| ₦30,000   | Completed
```

**Flagged Transactions:**
- Large amounts (> ₦500,000)
- Unusual patterns
- New user high activity
- Multiple failed attempts

### 4.2 Withdrawal Processing

**Pending Withdrawals Queue:**
1. Go to **Transactions** → **Pending Withdrawals**
2. Review each request:
   - User verification status
   - Transaction history
   - Account balance
   - Destination account
3. Actions:
   - **Approve**: Process withdrawal
   - **Hold**: Request verification
   - **Reject**: Decline with reason

**High-Value Withdrawal Protocol:**
For withdrawals > ₦500,000:
1. Verify user identity via phone
2. Confirm source of funds
3. Check AML indicators
4. Dual approval required

### 4.3 Refund Processing

**Initiating Refund:**
1. Locate original transaction
2. Click **Issue Refund**
3. Enter refund amount (full or partial)
4. Provide reason
5. Submit for approval

**Refund Approval:**
- < ₦10,000: Auto-approved
- ₦10,000 - ₦100,000: Single admin
- > ₦100,000: Dual approval required

### 4.4 Transaction Investigation

**Investigation Steps:**
1. Pull complete transaction trail
2. Review user communication logs
3. Check device/IP information
4. Contact user if needed
5. Document findings
6. Take appropriate action

---

## 5. Gig Management

### 5.1 Gig Moderation Queue

**Review New Gigs:**
1. Go to **Gigs** → **Pending Review**
2. Check for:
   - Prohibited content
   - Realistic budget
   - Clear requirements
   - No contact info in description
3. Actions:
   - **Approve**: Publish gig
   - **Reject**: Remove with reason
   - **Edit**: Modify and approve

### 5.2 Prohibited Content

**Auto-flagged Content:**
- External contact information
- Payment outside platform
- Illegal services
- Discriminatory requirements
- Copyright infringement

**Manual Review Triggers:**
- User reports
- Keyword matches
- Unusual activity patterns

### 5.3 Dispute Resolution

**Dispute Queue:**
1. Go to **Gigs** → **Disputes**
2. Review case details:
   - Original gig requirements
   - Submitted deliverables
   - Communication history
   - Both parties' claims
3. Resolution options:
   - **Favor Client**: Full refund
   - **Favor Freelancer**: Release payment
   - **Split**: Partial to each
   - **Escalate**: Need legal review

**Dispute SLA:**
- Initial response: 24 hours
- Resolution target: 72 hours
- Maximum: 7 days

---

## 6. Savings Circle Management

### 6.1 Circle Monitoring

**Active Circles Dashboard:**
```
Circle Name      | Type       | Members | Pool      | Health
─────────────────┼────────────┼─────────┼───────────┼────────
Tech Savers      | Rotational | 10/10   | ₦1,000,000| ✅ Good
Writers United   | Rotational | 8/12    | ₦400,000  | ✅ Good
Lagos Hustlers   | Fixed      | 15/15   | ₦750,000  | ⚠️ 2 late
Emergency Fund   | Fixed      | 6/10    | ₦180,000  | ❌ Stalled
```

### 6.2 Intervention Actions

**Late Contribution Handling:**
1. System sends automatic reminders
2. After 3 days late: Admin notification
3. Admin contacts member
4. Options:
   - Extend deadline
   - Remove member (refund eligible)
   - Circle dissolution (if critical)

**Circle Health Issues:**
| Issue | Action |
|-------|--------|
| Multiple late members | Send group reminder, contact admin creator |
| Payout failure | Investigate wallet, manual trigger if needed |
| Inactive circle | Contact creator, archive if abandoned |
| Dispute among members | Mediate, escalate if needed |

### 6.3 Manual Payout Trigger

When automatic payout fails:
1. Go to circle details
2. Click **Manual Payout**
3. Verify recipient is correct
4. Confirm payout amount
5. Process and document reason

---

## 7. Loan Management

### 7.1 Loan Applications

**Application Review (for flagged cases):**
Most loans are auto-approved based on credit score. Manual review for:
- Edge cases (score near threshold)
- Large amounts
- Previous default history
- Unusual patterns

**Review Criteria:**
- Credit score trend
- Income indicators (gig earnings)
- Existing obligations
- Account behavior

### 7.2 Default Management

**Default Timeline:**
| Days Overdue | Action |
|--------------|--------|
| 1-3 | Grace period, reminders |
| 4-7 | Late fee applied, escalated reminders |
| 8-14 | Phone contact attempt |
| 15-30 | Formal notice, credit score impact |
| 31+ | Collections process, account restriction |

**Collections Actions:**
1. Restrict new loans
2. Limit withdrawals
3. Apply to savings balance (if authorized)
4. External collections referral

### 7.3 Loan Restructuring

For users facing difficulty:
1. Review request
2. Options:
   - Payment holiday (1 month)
   - Extended tenure
   - Reduced installment (extended term)
3. Document agreement
4. Update loan terms

---

## 8. Support & Disputes

### 8.1 Support Ticket Queue

**Ticket Categories:**
| Priority | Category | SLA |
|----------|----------|-----|
| Critical | Account access, fraud | 1 hour |
| High | Failed transactions, disputes | 4 hours |
| Medium | Feature questions, complaints | 24 hours |
| Low | General inquiries | 48 hours |

### 8.2 Common Issues & Resolutions

**Deposit Not Credited:**
1. Verify payment on Paystack dashboard
2. Check webhook logs
3. If confirmed, credit manually
4. Document and investigate webhook failure

**Withdrawal Failed:**
1. Check bank account validity
2. Review Paystack transfer status
3. Re-initiate if bank issue
4. Refund to wallet if bank rejected

**OTP Not Received:**
1. Verify phone number format
2. Check SMS gateway status
3. Review delivery report
4. Offer alternative verification if repeated

### 8.3 Escalation Matrix

| Level | Handled By | Response Time |
|-------|------------|---------------|
| L1 | Support Agent | Immediate |
| L2 | Senior Support | 2 hours |
| L3 | Technical Team | 4 hours |
| L4 | Management | Same day |

---

## 9. System Configuration

### 9.1 Feature Flags

**Managing Features:**
1. Go to **Config** → **Feature Flags**
2. Toggle features on/off:
   - `loans_enabled`: Enable/disable loan applications
   - `savings_new_circles`: Allow new circle creation
   - `gigs_posting`: Enable gig posting
   - `withdrawals_enabled`: Process withdrawals
3. Changes take effect immediately

### 9.2 Rate Limits

**Configurable Limits:**
| Endpoint | Default | Adjustable |
|----------|---------|------------|
| OTP requests | 5/15min | Yes |
| Login attempts | 5/5min | Yes |
| Transaction rate | 30/hour | Yes |
| API general | 100/min | Yes |

### 9.3 Notification Templates

**Managing Templates:**
1. Go to **Config** → **Notifications**
2. Select template type (SMS, Push, Email)
3. Edit message content
4. Use variables: `{{user_name}}`, `{{amount}}`, etc.
5. Preview and save

### 9.4 Fee Configuration

**Adjustable Fees:**
| Fee Type | Default | Location |
|----------|---------|----------|
| Gig platform fee | 10% | Config → Fees |
| Loan interest | 5%/month | Config → Loans |
| Withdrawal fee | ₦10-50 | Config → Fees |
| Late payment penalty | 2% | Config → Loans |

---

## 10. Reports & Exports

### 10.1 Standard Reports

**Daily Reports (Auto-generated):**
- Transaction summary
- New user registrations
- Loan disbursements
- Support ticket metrics

**Weekly Reports:**
- Revenue analysis
- User growth
- Gig marketplace health
- Savings circle status

**Monthly Reports:**
- Financial statements
- KPI dashboard
- Compliance report
- System performance

### 10.2 Data Exports

**Export Options:**
- CSV for spreadsheets
- JSON for developers
- PDF for presentations

**Available Exports:**
- User list (with filters)
- Transaction history
- Gig analytics
- Loan portfolio

### 10.3 Scheduled Reports

1. Go to **Reports** → **Scheduled**
2. Create new schedule
3. Select report type
4. Set frequency (daily/weekly/monthly)
5. Add recipients
6. Reports delivered via email

---

## 11. Security & Audit

### 11.1 Audit Logs

All admin actions are logged:
- Timestamp
- Admin user
- Action taken
- Target (user/transaction)
- IP address
- Outcome

**Viewing Audit Logs:**
1. Go to **Security** → **Audit Logs**
2. Filter by date, admin, action type
3. Export for compliance

### 11.2 Access Reviews

**Monthly Access Review:**
- Review all admin accounts
- Verify role appropriateness
- Remove inactive accounts
- Document review

### 11.3 Suspicious Activity Monitoring

**Automated Alerts:**
- Multiple failed logins
- Unusual admin access patterns
- High-value transactions
- Velocity checks exceeded

**Manual Investigation:**
1. Review alert details
2. Check user/admin activity
3. Contact if legitimate
4. Take protective action if fraud

### 11.4 Compliance Reports

**Regulatory Requirements:**
- AML transaction reports
- Suspicious activity reports (SAR)
- Customer due diligence records
- Data protection compliance

---

## 12. Emergency Procedures

### 12.1 System Outage

**If API is down:**
1. Check monitoring dashboards
2. Notify engineering on-call
3. Update status page
4. Communicate with users (social media)
5. Document timeline

### 12.2 Security Breach

**If breach suspected:**
1. Do not discuss on unsecured channels
2. Contact Security Lead immediately
3. Preserve evidence (don't delete logs)
4. Follow incident response playbook
5. Notify legal/compliance if required

### 12.3 Mass Fraud Event

**If organized fraud detected:**
1. Enable emergency fraud filters
2. Pause high-risk operations (if needed)
3. Identify affected accounts
4. Coordinate response team
5. User communication plan

### 12.4 Payment Gateway Failure

**If Paystack is down:**
1. Monitor Paystack status page
2. Enable maintenance mode for deposits
3. Queue withdrawal requests
4. Communicate delays to users
5. Process backlog when restored

---

## Contact Information

**Internal Contacts:**
| Role | Name | Phone | Email |
|------|------|-------|-------|
| Engineering Lead | - | - | eng@hustlex.app |
| Security Lead | - | - | security@hustlex.app |
| Operations Lead | - | - | ops@hustlex.app |
| Legal | - | - | legal@hustlex.app |

**External Contacts:**
| Service | Support |
|---------|---------|
| Paystack | support@paystack.com |
| Termii | support@termii.com |
| AWS Support | AWS Console |

---

*This manual is confidential and for internal use only.*

**Version 1.0 | Last Updated: January 2024**
