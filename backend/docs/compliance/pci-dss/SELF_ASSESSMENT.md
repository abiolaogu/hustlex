# HustleX PCI DSS v4.0 Self-Assessment Questionnaire

**SAQ Type:** SAQ A-EP (E-commerce merchants using third-party payment processor)
**Assessment Date:** [Date]
**Assessor:** [Name/Internal]
**Next Assessment:** [Date + 1 year]

---

## Overview

HustleX processes payments through Paystack, a PCI DSS Level 1 certified payment service provider. As a result, HustleX does not directly store, process, or transmit cardholder data (CHD). This self-assessment documents controls that protect the payment integration and customer data.

### Scope Determination

| Component | In Scope | Justification |
|-----------|----------|---------------|
| Web/Mobile Application | Yes | Redirects to payment page |
| API Servers | Yes | Receives payment confirmations |
| Database | No | No CHD stored |
| Network Infrastructure | Yes | Hosts in-scope systems |
| Paystack Integration | Yes | Payment tokenization |

### PCI DSS Applicability

Since HustleX:
- Does NOT store card numbers (PAN)
- Does NOT process card numbers
- Uses tokenization via Paystack
- Receives only masked card info (last 4 digits)

**Applicable SAQ:** SAQ A-EP (partially), with many controls inherited from Paystack.

---

## Requirement 1: Install and Maintain Network Security Controls

### 1.1 Network Security Policies

| Control | Status | Evidence |
|---------|--------|----------|
| 1.1.1 Network security policies defined | Compliant | Information Security Policy |
| 1.1.2 Network diagram documented | Compliant | Architecture documentation |
| 1.1.3 Data flow diagram | Compliant | Payment flow documentation |

### 1.2 Network Security Controls

| Control | Status | Implementation |
|---------|--------|----------------|
| 1.2.1 Inbound/outbound traffic restricted | Compliant | Cloud security groups |
| 1.2.2 Configuration standards | Compliant | Terraform IaC |
| 1.2.3 DMZ for public-facing components | Compliant | Cloud VPC architecture |

### 1.3 Access Restriction

| Control | Status | Implementation |
|---------|--------|----------------|
| 1.3.1 Inbound traffic restricted | Compliant | Security group rules |
| 1.3.2 Outbound traffic restricted | Compliant | Egress rules defined |
| 1.3.3 Anti-spoofing measures | Compliant | Cloud provider controls |

**Notes:**
- Network security managed through AWS/GCP security groups and VPC configuration
- Infrastructure as Code (Terraform) ensures consistent configuration
- No direct internet access to database servers

---

## Requirement 2: Apply Secure Configurations

### 2.1 Secure Defaults

| Control | Status | Implementation |
|---------|--------|----------------|
| 2.1.1 Vendor defaults changed | Compliant | Automated provisioning |
| 2.1.2 Default accounts disabled/removed | Compliant | Configuration management |

### 2.2 Configuration Standards

| Control | Status | Implementation |
|---------|--------|----------------|
| 2.2.1 Configuration standards defined | Compliant | IaC templates |
| 2.2.2 Wireless settings (N/A) | N/A | Cloud infrastructure |
| 2.2.3 Primary functions separation | Compliant | Microservices architecture |
| 2.2.4 Necessary services only | Compliant | Minimal container images |
| 2.2.5 Insecure services secured | Compliant | TLS everywhere |
| 2.2.6 Security parameters configured | Compliant | Hardened configurations |
| 2.2.7 Non-console admin encrypted | Compliant | SSH/TLS only |

**Notes:**
- All infrastructure defined in Terraform with security baselines
- Container images are minimal (distroless where possible)
- Regular security scanning of configurations

---

## Requirement 3: Protect Stored Account Data

### 3.1 Data Retention

| Control | Status | Implementation |
|---------|--------|----------------|
| 3.1.1 Retention policies defined | Compliant | Data Retention Policy |
| 3.1.2 SAD not stored | Compliant | No SAD received/stored |

### 3.2 SAD Storage

| Control | Status | Implementation |
|---------|--------|----------------|
| 3.2.1 CVV not stored | N/A | Not received |
| 3.2.2 PIN not stored | N/A | Not received |

### 3.3 PAN Display

| Control | Status | Implementation |
|---------|--------|----------------|
| 3.3.1 PAN masked | N/A | Only last 4 digits received |
| 3.3.2 PAN unmasking restricted | N/A | Full PAN not available |
| 3.3.3 PAN copied securely | N/A | No full PAN |

### 3.4 PAN Storage

| Control | Status | Implementation |
|---------|--------|----------------|
| 3.4.1 PAN rendered unreadable | N/A | No PAN stored |
| 3.4.2 Disk encryption | Compliant | AES-256 at rest |

### 3.5 Cryptographic Keys

| Control | Status | Implementation |
|---------|--------|----------------|
| 3.5.1 Key management documented | Compliant | Encryption procedures |
| 3.5.2 Keys stored securely | Compliant | Cloud KMS |

**Notes:**
- HustleX does not receive or store full card numbers (PAN)
- Only last 4 digits and card type received from Paystack
- Transaction tokens are stored instead of card data
- All sensitive data encrypted with AES-256-GCM

---

## Requirement 4: Protect Cardholder Data During Transmission

### 4.1 Strong Cryptography

| Control | Status | Implementation |
|---------|--------|----------------|
| 4.1.1 Strong cryptography for CHD transmission | Compliant | TLS 1.3 |
| 4.1.2 Certificates managed properly | Compliant | Automated certificate management |

### 4.2 PAN Transmission

| Control | Status | Implementation |
|---------|--------|----------------|
| 4.2.1 PAN secured during transmission | N/A | No PAN transmitted |
| 4.2.2 PAN not sent via end-user messaging | Compliant | Not applicable |

**Notes:**
- All external communications over HTTPS (TLS 1.3)
- HSTS enabled with 1-year max-age
- Certificate transparency monitoring
- No full PAN transmitted through HustleX systems

---

## Requirement 5: Protect Systems Against Malware

### 5.1 Anti-Malware Solutions

| Control | Status | Implementation |
|---------|--------|----------------|
| 5.1.1 Anti-malware deployed | Compliant | Cloud-native protection |
| 5.1.2 Detection mechanisms current | Compliant | Automated updates |

### 5.2 Malware Prevention

| Control | Status | Implementation |
|---------|--------|----------------|
| 5.2.1 Periodic scans | Compliant | Continuous monitoring |
| 5.2.2 Anti-malware mechanisms enabled | Compliant | Cannot be disabled |
| 5.2.3 Audit logs generated | Compliant | Cloud security logs |

### 5.3 Anti-Phishing

| Control | Status | Implementation |
|---------|--------|----------------|
| 5.3.1 Phishing protection | Compliant | Email filtering |
| 5.3.2 User awareness | Compliant | Security training |
| 5.3.3 Phishing simulations | Compliant | Quarterly tests |

**Notes:**
- Container-based infrastructure limits malware attack surface
- Immutable infrastructure prevents persistent malware
- Cloud provider security scanning enabled

---

## Requirement 6: Develop and Maintain Secure Systems

### 6.1 Security Patches

| Control | Status | Implementation |
|---------|--------|----------------|
| 6.1.1 Vulnerability process | Compliant | Vulnerability management |
| 6.1.2 Vulnerability identification | Compliant | Automated scanning |
| 6.1.3 Patches installed timely | Compliant | Critical: 72 hours |

### 6.2 Secure Development

| Control | Status | Implementation |
|---------|--------|----------------|
| 6.2.1 Secure SDLC | Compliant | SSDLC documentation |
| 6.2.2 Training | Compliant | Annual secure coding |
| 6.2.3 Code review | Compliant | PR requirements |
| 6.2.4 Vulnerable code prevented | Compliant | SAST scanning |

### 6.3 Security Testing

| Control | Status | Implementation |
|---------|--------|----------------|
| 6.3.1 Security testing pre-production | Compliant | CI/CD security gates |
| 6.3.2 Production/test data separated | Compliant | Separate environments |
| 6.3.3 Test data secured | Compliant | Anonymized data |

### 6.4 Public-Facing Applications

| Control | Status | Implementation |
|---------|--------|----------------|
| 6.4.1 XSS prevention | Compliant | Output encoding, CSP |
| 6.4.2 Input validation | Compliant | Server-side validation |
| 6.4.3 WAF deployed | Compliant | Cloud WAF |

### 6.5 Change Management

| Control | Status | Implementation |
|---------|--------|----------------|
| 6.5.1 Change control procedures | Compliant | Change Management Policy |
| 6.5.2 Significant changes reviewed | Compliant | Security review |
| 6.5.3 Rollback procedures | Compliant | Documented rollback |
| 6.5.4 Production access restricted | Compliant | Separate accounts |
| 6.5.5 Test/dev data separate | Compliant | Environment isolation |
| 6.5.6 Production data in test | Compliant | Not used / anonymized |

**Notes:**
- OWASP Top 10 vulnerabilities addressed in secure coding training
- Automated SAST/DAST in CI/CD pipeline
- Annual penetration testing by third party
- Code references: `internal/infrastructure/security/validation/`

---

## Requirement 7: Restrict Access

### 7.1 Access Control

| Control | Status | Implementation |
|---------|--------|----------------|
| 7.1.1 Access control policies | Compliant | Access Control Policy |
| 7.1.2 Business need-to-know | Compliant | RBAC implementation |
| 7.1.3 Access control systems | Compliant | IAM + Application RBAC |
| 7.1.4 Segregation of duties | Compliant | Role separation |

### 7.2 User Access Management

| Control | Status | Implementation |
|---------|--------|----------------|
| 7.2.1 Appropriate access levels | Compliant | Role-based access |
| 7.2.2 Privileges reviewed | Compliant | Quarterly reviews |
| 7.2.3 Unused accounts disabled | Compliant | 90-day inactivity |
| 7.2.4 Vendor accounts managed | Compliant | Time-limited access |

**Notes:**
- Role-based access control implemented in application
- Code reference: `internal/interface/http/middleware/auth.go`
- Quarterly access reviews documented
- Principle of least privilege enforced

---

## Requirement 8: Identify Users and Authenticate Access

### 8.1 User Identification

| Control | Status | Implementation |
|---------|--------|----------------|
| 8.1.1 Unique user IDs | Compliant | UUID-based identification |
| 8.1.2 Shared accounts prohibited | Compliant | Policy enforced |

### 8.2 Strong Authentication

| Control | Status | Implementation |
|---------|--------|----------------|
| 8.2.1 Appropriate authentication | Compliant | Password + MFA |
| 8.2.2 Strong cryptography for credentials | Compliant | Argon2id hashing |
| 8.2.3 MFA for CDE access | Compliant | Required for admin |

### 8.3 Authentication Policies

| Control | Status | Implementation |
|---------|--------|----------------|
| 8.3.1 Password complexity | Compliant | 12+ chars, complexity |
| 8.3.2 Password changes | Compliant | 90-day rotation |
| 8.3.3 Password history | Compliant | 12 passwords |
| 8.3.4 Account lockout | Compliant | 5 attempts |
| 8.3.5 MFA for remote access | Compliant | Required |
| 8.3.6 MFA for non-console admin | Compliant | Required |

### 8.4 Multi-Factor Authentication

| Control | Status | Implementation |
|---------|--------|----------------|
| 8.4.1 MFA for CDE access | Compliant | TOTP + backup codes |
| 8.4.2 MFA for remote access | Compliant | VPN + MFA |
| 8.4.3 MFA factors independent | Compliant | Separate channels |

### 8.5 System and Service Accounts

| Control | Status | Implementation |
|---------|--------|----------------|
| 8.5.1 Service accounts managed | Compliant | Vault managed |
| 8.5.2 Service account passwords | Compliant | Complex, rotated |

**Notes:**
- Password hashing: Argon2id (OWASP recommended)
- MFA via TOTP for admin, OTP for customer transactions
- Code reference: `internal/infrastructure/security/crypto/encryption.go`

---

## Requirement 9: Restrict Physical Access

### 9.1 Physical Access Controls

| Control | Status | Implementation |
|---------|--------|----------------|
| 9.1.1 Facility entry controls | Inherited | Cloud provider (AWS/GCP SOC 2) |
| 9.1.2 Visitor management | Inherited | Cloud provider |

**Notes:**
- All infrastructure in cloud (AWS/GCP)
- Physical security inherited from cloud provider
- Cloud provider maintains SOC 2 Type II certification
- No on-premises cardholder data environment

---

## Requirement 10: Log and Monitor Access

### 10.1 Audit Logging

| Control | Status | Implementation |
|---------|--------|----------------|
| 10.1.1 Audit logs enabled | Compliant | Application + infrastructure |
| 10.1.2 Audit logs for CDE | Compliant | All access logged |

### 10.2 Audit Log Content

| Control | Status | Implementation |
|---------|--------|----------------|
| 10.2.1 User access logged | Compliant | Audit logger |
| 10.2.2 Actions logged | Compliant | All CRUD operations |
| 10.2.3 Event details captured | Compliant | Who, what, when, where |

### 10.3 Log Protection

| Control | Status | Implementation |
|---------|--------|----------------|
| 10.3.1 Read access restricted | Compliant | Admin only |
| 10.3.2 Modification prevented | Compliant | Append-only |
| 10.3.3 Backup of logs | Compliant | Offsite storage |
| 10.3.4 Logs retained | Compliant | 1 year online, 7 years archive |

### 10.4 Log Review

| Control | Status | Implementation |
|---------|--------|----------------|
| 10.4.1 Daily log review | Compliant | Automated alerts |
| 10.4.2 Security events reviewed | Compliant | SIEM integration |
| 10.4.3 Exceptions investigated | Compliant | Incident process |

### 10.5 Time Synchronization

| Control | Status | Implementation |
|---------|--------|----------------|
| 10.5.1 Time synchronization | Compliant | NTP configuration |
| 10.5.2 Time data protected | Compliant | Restricted access |

**Notes:**
- Comprehensive audit logging implemented
- Code reference: `internal/infrastructure/security/audit/logger.go`
- Logs stored in PostgreSQL `audit_logs` table
- Retention: 7 years for compliance

---

## Requirement 11: Test Security Regularly

### 11.1 Security Testing

| Control | Status | Implementation |
|---------|--------|----------------|
| 11.1.1 Testing processes defined | Compliant | Security testing procedures |
| 11.1.2 Wireless AP detection | N/A | Cloud infrastructure |

### 11.2 Vulnerability Scanning

| Control | Status | Implementation |
|---------|--------|----------------|
| 11.2.1 Internal vulnerability scans | Compliant | Quarterly |
| 11.2.2 External vulnerability scans | Compliant | Quarterly (ASV) |
| 11.2.3 Remediation procedures | Compliant | Documented |

### 11.3 Penetration Testing

| Control | Status | Implementation |
|---------|--------|----------------|
| 11.3.1 Penetration testing | Compliant | Annual external |
| 11.3.2 Internal pen test | Compliant | Annual |
| 11.3.3 Findings remediated | Compliant | Per severity |

### 11.4 Intrusion Detection

| Control | Status | Implementation |
|---------|--------|----------------|
| 11.4.1 Intrusion detection | Compliant | Cloud IDS/IPS |
| 11.4.2 Detection mechanisms | Compliant | Network + host |
| 11.4.3 Alerts generated | Compliant | Real-time |

### 11.5 Change Detection

| Control | Status | Implementation |
|---------|--------|----------------|
| 11.5.1 Change detection mechanism | Compliant | FIM, Git |
| 11.5.2 Critical file monitoring | Compliant | Automated alerts |

**Notes:**
- Quarterly vulnerability scans (internal + external ASV)
- Annual penetration test by third party
- Findings tracked to remediation

---

## Requirement 12: Support Information Security

### 12.1 Security Policies

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.1.1 Security policy established | Compliant | Information Security Policy |
| 12.1.2 Policy reviewed annually | Compliant | Annual review process |
| 12.1.3 Security roles defined | Compliant | RACI matrix |
| 12.1.4 Responsibility assigned | Compliant | CISO appointed |

### 12.2 Acceptable Use

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.2.1 Acceptable use policies | Compliant | AUP documented |

### 12.3 Risk Assessment

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.3.1 Risk assessment performed | Compliant | Annual |
| 12.3.2 Risk assessment documented | Compliant | Risk register |
| 12.3.3 Risk assessment reviewed | Compliant | Quarterly review |

### 12.4 PCI DSS Compliance

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.4.1 Compliance responsibility | Compliant | CISO assigned |
| 12.4.2 Compliance program | Compliant | Annual SAQ |

### 12.5 PCI DSS Scope

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.5.1 Scope documented | Compliant | This document |
| 12.5.2 Scope validated | Compliant | Annual review |

### 12.6 Security Awareness

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.6.1 Security awareness program | Compliant | Annual training |
| 12.6.2 Personnel acknowledgment | Compliant | Training records |
| 12.6.3 Phishing awareness | Compliant | Quarterly tests |

### 12.7 Personnel Screening

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.7.1 Background checks | Compliant | Pre-employment |

### 12.8 Third-Party Management

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.8.1 Third-parties identified | Compliant | Vendor inventory |
| 12.8.2 Third-party agreements | Compliant | Security terms |
| 12.8.3 Third-party due diligence | Compliant | Annual review |
| 12.8.4 Third-party compliance | Compliant | Paystack PCI DSS L1 |
| 12.8.5 Third-party responsibilities | Compliant | Documented |

### 12.9 Service Provider Compliance

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.9.1 Service provider agreements | Compliant | Paystack agreement |
| 12.9.2 Service provider monitoring | Compliant | AOC review |

### 12.10 Incident Response

| Control | Status | Implementation |
|---------|--------|----------------|
| 12.10.1 Incident response plan | Compliant | IR Policy |
| 12.10.2 Incident response procedures | Compliant | Documented |
| 12.10.3 Incident response testing | Compliant | Annual exercise |
| 12.10.4 Personnel trained | Compliant | Training records |
| 12.10.5 Security alerts monitored | Compliant | 24/7 |
| 12.10.6 Lessons learned | Compliant | Post-incident review |

---

## Attestation

This self-assessment questionnaire has been completed to the best of our knowledge. HustleX is committed to maintaining PCI DSS compliance and protecting cardholder data.

| Role | Name | Signature | Date |
|------|------|-----------|------|
| CISO | [Name] | | |
| CTO | [Name] | | |
| CEO | [Name] | | |

---

## Appendix A: Third-Party Service Providers

| Provider | Service | PCI DSS Status | Evidence |
|----------|---------|----------------|----------|
| Paystack | Payment processing | Level 1 certified | AOC on file |
| AWS/GCP | Infrastructure | SOC 2 Type II | Report on file |

## Appendix B: Network Diagram

[Reference: infrastructure/diagrams/network-architecture.png]

## Appendix C: Data Flow Diagram

```
Customer Browser
       │
       ▼ (HTTPS)
  Load Balancer
       │
       ▼ (TLS)
   API Server ────────────────▶ Paystack API
       │                            │
       │ (No PAN)                   │ (Tokenized)
       │                            │
       ▼                            ▼
   Database                   Card Network
  (Tokens only)
```

## Appendix D: Compensating Controls

No compensating controls required - all applicable requirements are met directly.

## Appendix E: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [Date] | [Author] | Initial assessment |
