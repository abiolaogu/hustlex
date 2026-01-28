# HustleX Vendor Management Policy

**Document ID:** HX-POL-007
**Version:** 1.0
**Effective Date:** [Date]
**Last Review:** [Date]
**Next Review:** [Date + 1 year]
**Owner:** Chief Information Security Officer (CISO)

---

## 1. Purpose

This Vendor Management Policy establishes requirements for selecting, assessing, and managing third-party vendors and service providers. It ensures that vendors handling HustleX data or providing critical services meet appropriate security and compliance standards.

## 2. Scope

This policy applies to:
- All third-party vendors and service providers
- Cloud service providers
- Payment processors and financial partners
- Software and SaaS providers
- Consultants and contractors with system access
- Data processors acting on behalf of HustleX

## 3. Vendor Classification

### 3.1 Risk Tiers

| Tier | Risk Level | Criteria | Assessment Frequency |
|------|------------|----------|---------------------|
| Tier 1 | Critical | Access to PII/financial data, critical operations | Annual + continuous monitoring |
| Tier 2 | High | Access to confidential data, important operations | Annual |
| Tier 3 | Medium | Limited data access, replaceable services | Every 2 years |
| Tier 4 | Low | No data access, commodity services | Initial only |

### 3.2 Classification Criteria

**Tier 1 (Critical):**
- Processes or stores customer financial data
- Processes or stores BVN/NIN
- Provides core payment processing
- Single point of failure for operations
- Direct database access

**Tier 2 (High):**
- Access to customer PII
- Provides security services
- Infrastructure provider
- Has administrative access

**Tier 3 (Medium):**
- Access to internal data only
- Provides development tools
- Marketing services
- Consulting services

**Tier 4 (Low):**
- No data access
- Commodity services
- Office supplies
- General contractors

## 4. Vendor Selection

### 4.1 Selection Process

1. **Identify Need** - Document business requirement
2. **Research** - Identify potential vendors
3. **Initial Screen** - Basic capability and reputation check
4. **Security Assessment** - Based on risk tier
5. **Due Diligence** - Financial stability, references
6. **Negotiation** - Contract and security terms
7. **Approval** - Based on tier and risk
8. **Onboarding** - Implement controls and access

### 4.2 Selection Criteria

| Criterion | Weight | Tier 1 | Tier 2-4 |
|-----------|--------|--------|----------|
| Security certifications | High | Required | Preferred |
| Financial stability | High | Required | Required |
| Technical capability | High | Required | Required |
| Regulatory compliance | High | Required | Based on scope |
| References | Medium | Required | Preferred |
| SLA commitments | Medium | Required | Required |
| Incident history | Medium | Required | Preferred |

### 4.3 Security Certifications

**Preferred Certifications:**
- SOC 2 Type II
- ISO 27001
- PCI DSS (for payment processors)
- ISO 27701 (for data processors)

**Minimum for Tier 1:**
- SOC 2 Type II or equivalent
- Penetration test results (annual)
- Security policy documentation

## 5. Security Assessment

### 5.1 Assessment Requirements by Tier

| Assessment | Tier 1 | Tier 2 | Tier 3 | Tier 4 |
|------------|--------|--------|--------|--------|
| Security questionnaire | Full | Full | Abbreviated | No |
| SOC 2 report review | Required | Required | If available | No |
| Penetration test results | Required | If available | No | No |
| On-site assessment | If needed | No | No | No |
| Policy review | Full | Summary | No | No |
| Technical assessment | Required | If needed | No | No |

### 5.2 Security Questionnaire

The security questionnaire covers:
- Information security governance
- Access control
- Data protection
- Network security
- Incident response
- Business continuity
- Compliance
- Subcontractor management

### 5.3 Assessment Findings

| Finding Severity | Tier 1 Action | Tier 2-3 Action |
|-----------------|---------------|-----------------|
| Critical | Block engagement until resolved | Block or accept with mitigation |
| High | Remediation plan required | Remediation or risk acceptance |
| Medium | Document and monitor | Document |
| Low | Note for improvement | Note |

## 6. Contractual Requirements

### 6.1 Required Contract Terms

**For All Vendors:**
- Confidentiality obligations
- Data protection requirements
- Right to audit (for Tier 1-2)
- Termination rights
- Insurance requirements
- Subcontractor restrictions

**For Tier 1-2 Vendors:**
- Security standards compliance
- Incident notification requirements
- Data location restrictions
- Personnel security requirements
- Business continuity requirements
- Liability provisions

### 6.2 Data Processing Agreement (DPA)

Required for vendors processing personal data:
- Processing purpose and scope
- Data subject categories
- Security measures
- Subprocessor requirements
- Data subject rights support
- Deletion/return requirements
- Audit rights

### 6.3 SLA Requirements

| Tier | Availability | Response Time | Support |
|------|--------------|---------------|---------|
| Tier 1 | 99.9% | < 1 hour critical | 24/7 |
| Tier 2 | 99.5% | < 4 hours critical | Business hours |
| Tier 3 | 99% | < 1 business day | Business hours |
| Tier 4 | Standard | Standard | Standard |

## 7. Ongoing Management

### 7.1 Monitoring Requirements

| Activity | Tier 1 | Tier 2 | Tier 3 | Tier 4 |
|----------|--------|--------|--------|--------|
| Performance review | Monthly | Quarterly | Annual | As needed |
| Security review | Annual + continuous | Annual | Every 2 years | Initial |
| Compliance check | Quarterly | Annual | Every 2 years | No |
| Access review | Quarterly | Semi-annual | Annual | N/A |
| Contract review | Annual | Annual | At renewal | At renewal |

### 7.2 Continuous Monitoring (Tier 1)

- Security rating services (if available)
- Public breach notifications
- Regulatory actions
- Financial news
- Service status monitoring

### 7.3 Periodic Reviews

**Annual Review (Tier 1-2):**
- Updated SOC 2 report
- Updated security questionnaire
- Performance against SLAs
- Incident review
- Contract compliance

**Bi-annual Review (Tier 3):**
- Security questionnaire update
- Performance review
- Contract compliance

## 8. Access Management

### 8.1 Access Requirements

| Access Type | Approval Required | Review Frequency |
|-------------|-------------------|------------------|
| Production data | Security + Business owner | Quarterly |
| Production systems | Security + Engineering | Quarterly |
| Test data (sanitized) | Business owner | Annual |
| Test systems | Engineering lead | Annual |
| Network access | Security | Quarterly |

### 8.2 Access Controls

- Unique accounts for all vendor personnel
- Multi-factor authentication required
- VPN for remote access
- Activity logging
- Time-limited access where possible

### 8.3 Access Termination

Upon contract end:
- Disable all accounts immediately
- Remove network access
- Verify data deletion/return
- Revoke API keys/tokens
- Update documentation

## 9. Incident Management

### 9.1 Notification Requirements

| Incident Type | Notification | Timeline |
|---------------|--------------|----------|
| Security breach | HustleX Security Team | Immediate (< 2 hours) |
| Data breach | CISO + Legal | Immediate (< 2 hours) |
| Service outage | Operations Team | < 15 minutes |
| Security vulnerability | Security Team | < 24 hours |
| Compliance issue | Compliance Team | < 48 hours |

### 9.2 Incident Response

Vendors must:
- Have documented incident response procedures
- Provide root cause analysis
- Implement corrective actions
- Support HustleX investigations
- Preserve evidence as required

## 10. Compliance Requirements

### 10.1 Regulatory Alignment

| Regulation | Vendor Requirement |
|------------|-------------------|
| NDPR | DPA in place, data subject support |
| CBN Guidelines | Documented controls, audit support |
| PCI DSS | AOC/SAQ for payment processors |

### 10.2 Audit Rights

For Tier 1-2 vendors:
- Annual attestation of compliance
- Right to audit with reasonable notice
- Third-party audit acceptance
- Penetration test result sharing

## 11. Subcontractor Management

### 11.1 Requirements

Vendors must:
- Disclose all subcontractors
- Flow down security requirements
- Obtain approval for new subcontractors
- Remain liable for subcontractor actions

### 11.2 Subcontractor Assessment

Tier 1 vendor subcontractors:
- Security questionnaire
- Contract requirements
- Ongoing monitoring

## 12. Termination and Exit

### 12.1 Exit Planning

All Tier 1-2 contracts must include:
- Transition assistance period
- Data migration support
- Knowledge transfer requirements
- Minimum notice period

### 12.2 Exit Checklist

- [ ] Notify vendor of termination
- [ ] Initiate data export/return
- [ ] Disable access accounts
- [ ] Revoke API credentials
- [ ] Verify data deletion
- [ ] Obtain deletion certificate
- [ ] Update documentation
- [ ] Archive contract and records

### 12.3 Data Handling

At termination:
- Return all HustleX data
- Delete all copies securely
- Provide written confirmation
- Retain only as legally required

## 13. Roles and Responsibilities

### 13.1 Business Owner

- Identify vendor need
- Initiate selection process
- Manage vendor relationship
- Monitor performance

### 13.2 Security Team

- Conduct security assessments
- Review SOC reports
- Approve security controls
- Monitor security posture

### 13.3 Legal/Compliance

- Review and negotiate contracts
- Ensure regulatory compliance
- Manage DPAs
- Handle disputes

### 13.4 Procurement

- Manage vendor database
- Coordinate assessments
- Track contract terms
- Facilitate renewals

---

## Appendix A: Current Critical Vendors (Tier 1)

| Vendor | Service | Security Cert | Last Review |
|--------|---------|---------------|-------------|
| Paystack | Payment processing | PCI DSS L1 | [Date] |
| AWS/GCP | Cloud infrastructure | SOC 2 Type II | [Date] |
| [Database] | Managed database | SOC 2 Type II | [Date] |

## Appendix B: Security Questionnaire Topics

1. Information Security Governance
2. Risk Management
3. Asset Management
4. Access Control
5. Cryptography
6. Physical Security
7. Operations Security
8. Communications Security
9. System Development
10. Supplier Relationships
11. Incident Management
12. Business Continuity
13. Compliance

## Appendix C: Vendor Assessment Workflow

```
New Vendor Request
        │
        ▼
  Classify Risk Tier
        │
   ┌────┴────┬────────┬────────┐
   ▼         ▼        ▼        ▼
Tier 1    Tier 2   Tier 3   Tier 4
   │         │        │        │
   ▼         ▼        ▼        ▼
Full      Full     Abbrev.   Basic
Assessment Assess.  Quest.   Check
   │         │        │        │
   ▼         ▼        ▼        ▼
Security  Security  Business  Business
+ Legal   + Legal   Approval  Approval
Approval  Approval
   │         │        │        │
   └────┬────┴────────┴────────┘
        ▼
   Contract & Onboard
```

## Appendix D: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [Date] | [Author] | Initial release |
