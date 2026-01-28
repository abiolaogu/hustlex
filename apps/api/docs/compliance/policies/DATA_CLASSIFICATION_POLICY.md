# HustleX Data Classification Policy

**Document ID:** HX-POL-003
**Version:** 1.0
**Effective Date:** [Date]
**Last Review:** [Date]
**Next Review:** [Date + 1 year]
**Owner:** Chief Information Security Officer (CISO)

---

## 1. Purpose

This Data Classification Policy establishes a framework for categorizing data based on sensitivity and business value. It ensures appropriate protection measures are applied consistently across all data types handled by HustleX.

## 2. Scope

This policy applies to:
- All data created, collected, stored, or processed by HustleX
- All employees, contractors, and third parties handling HustleX data
- All systems and applications that process, store, or transmit data
- Data in all formats: electronic, paper, verbal

## 3. Classification Levels

### 3.1 Highly Confidential (Level 4)

**Definition:** Data that, if disclosed, could cause severe harm to HustleX, its customers, or partners. This includes data protected by regulatory requirements.

**Examples:**
- Bank Verification Numbers (BVN)
- National Identification Numbers (NIN)
- Transaction PINs and passwords
- Encryption keys and secrets
- Customer financial account details
- Payment card data (if applicable)
- API keys and authentication tokens

**Handling Requirements:**
| Control | Requirement |
|---------|-------------|
| Encryption at Rest | AES-256-GCM required |
| Encryption in Transit | TLS 1.3 required |
| Access Control | Need-to-know, MFA required |
| Logging | Full audit trail required |
| Storage | Encrypted database only |
| Transmission | Encrypted channels only |
| Retention | Per regulatory requirements |
| Disposal | Cryptographic erasure |

### 3.2 Confidential (Level 3)

**Definition:** Data that is sensitive to HustleX operations and could cause significant harm if disclosed to unauthorized parties.

**Examples:**
- Customer personal information (name, address, phone)
- Transaction history and account balances
- Internal financial reports
- Employee personal data
- Business strategies and plans
- Source code and technical documentation
- Vendor contracts and pricing

**Handling Requirements:**
| Control | Requirement |
|---------|-------------|
| Encryption at Rest | Required for PII |
| Encryption in Transit | TLS 1.2+ required |
| Access Control | Role-based access |
| Logging | Access logging required |
| Storage | Secured systems only |
| Transmission | Encrypted preferred |
| Retention | Per retention schedule |
| Disposal | Secure deletion |

### 3.3 Internal (Level 2)

**Definition:** Data intended for internal use that could cause minor harm if disclosed externally.

**Examples:**
- Internal policies and procedures
- Non-sensitive meeting notes
- Internal announcements
- Project documentation
- Training materials
- General business correspondence

**Handling Requirements:**
| Control | Requirement |
|---------|-------------|
| Encryption at Rest | Optional |
| Encryption in Transit | TLS recommended |
| Access Control | Employee access |
| Logging | Standard logging |
| Storage | Internal systems |
| Transmission | Internal channels |
| Retention | Business need |
| Disposal | Standard deletion |

### 3.4 Public (Level 1)

**Definition:** Data that is intended for public disclosure or has no confidentiality requirements.

**Examples:**
- Marketing materials
- Public website content
- Press releases
- Public APIs documentation
- Job postings

**Handling Requirements:**
| Control | Requirement |
|---------|-------------|
| Encryption at Rest | Not required |
| Encryption in Transit | TLS for integrity |
| Access Control | Open access |
| Logging | Optional |
| Storage | Any approved system |
| Transmission | Any channel |
| Retention | Business need |
| Disposal | Standard deletion |

## 4. Classification Matrix

### 4.1 Customer Data Classification

| Data Element | Classification | Regulatory Basis |
|-------------|----------------|------------------|
| BVN | Highly Confidential | CBN Guidelines |
| NIN | Highly Confidential | NIMC Act |
| Transaction PIN | Highly Confidential | PCI DSS |
| Password Hash | Highly Confidential | Security Best Practice |
| Full Name | Confidential | NDPR |
| Email Address | Confidential | NDPR |
| Phone Number | Confidential | NDPR |
| Date of Birth | Confidential | NDPR |
| Account Balance | Confidential | Financial Privacy |
| Transaction History | Confidential | Financial Privacy |
| Profile Photo | Internal | Business Data |

### 4.2 System Data Classification

| Data Element | Classification | Notes |
|-------------|----------------|-------|
| Encryption Keys | Highly Confidential | Hardware security module |
| API Secrets | Highly Confidential | Key vault storage |
| Database Credentials | Highly Confidential | Secret management |
| Audit Logs | Confidential | Compliance requirement |
| Application Logs | Internal | May contain PII |
| Metrics Data | Internal | Operational data |
| Configuration Files | Internal | Infrastructure data |

## 5. Labeling Requirements

### 5.1 Electronic Documents

- Highly Confidential: Header/footer marking "HIGHLY CONFIDENTIAL"
- Confidential: Header/footer marking "CONFIDENTIAL"
- Internal: Header/footer marking "INTERNAL USE ONLY"
- Public: No marking required

### 5.2 Database Fields

Fields containing classified data must be:
- Documented in data dictionary
- Tagged with classification level
- Protected per classification requirements
- Subject to masking in non-production environments

### 5.3 Email Communications

- Use classification markers in subject line for Confidential and above
- Enable encryption for Highly Confidential content
- Verify recipients before sending sensitive data

## 6. Data Handling Procedures

### 6.1 Collection

| Classification | Requirements |
|---------------|--------------|
| Highly Confidential | Legal basis, consent, minimize collection |
| Confidential | Business justification, consent where applicable |
| Internal | Business purpose |
| Public | No restrictions |

### 6.2 Processing

| Classification | Requirements |
|---------------|--------------|
| Highly Confidential | Encrypted processing, audit logging |
| Confidential | Access controls, logging |
| Internal | Standard controls |
| Public | No special requirements |

### 6.3 Storage

| Classification | Requirements |
|---------------|--------------|
| Highly Confidential | Encrypted, segregated, access-controlled |
| Confidential | Encrypted, access-controlled |
| Internal | Standard security controls |
| Public | Standard storage |

### 6.4 Transmission

| Classification | Requirements |
|---------------|--------------|
| Highly Confidential | TLS 1.3, encrypted payload |
| Confidential | TLS 1.2+, verify recipient |
| Internal | Encrypted channel preferred |
| Public | Standard channels |

### 6.5 Disposal

| Classification | Method |
|---------------|--------|
| Highly Confidential | Cryptographic erasure, certificate of destruction |
| Confidential | Secure deletion, DOD 5220.22-M or equivalent |
| Internal | Standard secure delete |
| Public | Standard deletion |

## 7. Data Masking Standards

### 7.1 Masking Formats

| Data Type | Display Format | Storage |
|-----------|---------------|---------|
| BVN | 123****890 | Encrypted |
| NIN | 12345****890 | Encrypted |
| Phone | +234****1234 | Full number, access controlled |
| Email | a****@domain.com | Full email, access controlled |
| Account Number | ****1234 | Full number, access controlled |
| Card Number (PAN) | ****-****-****-1234 | Not stored (tokenized) |

### 7.2 Implementation Reference

```go
// internal/infrastructure/security/crypto/encryption.go
func MaskPII(value string, showLast int) string
func MaskEmail(email string) string
func MaskBVN(bvn string) string
func MaskNIN(nin string) string
```

## 8. Roles and Responsibilities

### 8.1 Data Owners

- Assign classification to data under their ownership
- Define access requirements
- Approve access requests
- Review access periodically

### 8.2 Data Custodians

- Implement security controls per classification
- Maintain encryption and access controls
- Monitor for unauthorized access
- Report security incidents

### 8.3 Data Users

- Handle data per classification requirements
- Report suspected data breaches
- Complete classification training
- Follow handling procedures

## 9. Compliance Mapping

### 9.1 NDPR Alignment

| NDPR Principle | Data Classification Support |
|----------------|---------------------------|
| Lawfulness | Classification defines legal basis requirements |
| Purpose Limitation | Classification controls usage |
| Data Minimization | Highly Confidential requires justification |
| Accuracy | Classification includes validation requirements |
| Storage Limitation | Classification defines retention |
| Security | Classification defines protection measures |

### 9.2 PCI DSS Alignment (if applicable)

| PCI DSS Requirement | Data Classification Support |
|---------------------|---------------------------|
| Req 3: Protect stored data | Encryption requirements by classification |
| Req 4: Encrypt transmission | TLS requirements by classification |
| Req 7: Restrict access | Access control by classification |
| Req 9: Physical security | Handling requirements |

## 10. Training and Awareness

All employees must:
- Complete data classification training upon hire
- Receive annual refresher training
- Acknowledge understanding of this policy
- Know how to classify data they handle

## 11. Exceptions

Exceptions require:
- Business justification
- Risk assessment
- Compensating controls
- CISO approval
- Documentation with expiration

## 12. Enforcement

Violations may result in:
- Disciplinary action
- Access revocation
- Termination
- Legal action
- Regulatory reporting

---

## Appendix A: Classification Decision Tree

```
Is the data regulated (BVN, NIN, PCI)?
├── Yes → HIGHLY CONFIDENTIAL
└── No → Does disclosure cause significant harm?
    ├── Yes → Is it PII or financial data?
    │   ├── Yes → CONFIDENTIAL
    │   └── No → CONFIDENTIAL or INTERNAL (case-by-case)
    └── No → Is it for internal use?
        ├── Yes → INTERNAL
        └── No → PUBLIC
```

## Appendix B: Quick Reference Card

| If handling... | Classification | Key Requirement |
|----------------|---------------|-----------------|
| BVN, NIN, PINs | Highly Confidential | Must be encrypted |
| Customer names, emails | Confidential | Access controlled |
| Internal docs | Internal | Employee access only |
| Marketing content | Public | No restrictions |

## Appendix C: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [Date] | [Author] | Initial release |
