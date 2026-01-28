# HustleX Password Policy

**Document ID:** HX-POL-008
**Version:** 1.0
**Effective Date:** [Date]
**Last Review:** [Date]
**Next Review:** [Date + 1 year]
**Owner:** Chief Information Security Officer (CISO)

---

## 1. Purpose

This Password Policy establishes requirements for creating, managing, and protecting passwords and authentication credentials. It ensures that authentication mechanisms adequately protect HustleX systems and customer data.

## 2. Scope

This policy applies to:
- All employees, contractors, and third parties with system access
- All HustleX systems and applications
- Customer-facing authentication
- API and service account credentials
- Administrative and privileged accounts

## 3. Password Requirements

### 3.1 User Account Passwords

| Requirement | Standard | Notes |
|-------------|----------|-------|
| Minimum length | 12 characters | 14+ recommended |
| Complexity | 3 of 4 categories | Upper, lower, number, special |
| Maximum age | 90 days | 365 days for customers |
| Minimum age | 1 day | Prevent rapid cycling |
| History | Last 12 passwords | Prevent reuse |
| Lockout threshold | 5 failed attempts | 15-minute lockout |

### 3.2 Character Categories

| Category | Characters | Required |
|----------|------------|----------|
| Uppercase | A-Z | At least 1 |
| Lowercase | a-z | At least 1 |
| Numbers | 0-9 | At least 1 |
| Special | !@#$%^&*()-_+=[] | Recommended |

### 3.3 Privileged Account Passwords

| Requirement | Standard |
|-------------|----------|
| Minimum length | 16 characters |
| Complexity | All 4 categories |
| Maximum age | 30 days |
| History | Last 24 passwords |
| MFA | Required |
| Storage | Password vault only |

### 3.4 Service Account Passwords

| Requirement | Standard |
|-------------|----------|
| Minimum length | 32 characters |
| Complexity | Randomly generated |
| Maximum age | 90 days |
| Storage | Secret management system |
| Rotation | Automated where possible |

## 4. Customer Authentication

### 4.1 Customer Password Requirements

| Requirement | Standard |
|-------------|----------|
| Minimum length | 8 characters |
| Complexity | 2 of 4 categories minimum |
| Maximum age | 365 days (recommended, not forced) |
| Lockout threshold | 5 failed attempts |
| Lockout duration | Progressive (5, 15, 60 minutes) |

### 4.2 Transaction PIN

| Requirement | Standard |
|-------------|----------|
| Length | 4-6 digits |
| Attempt limit | 3 attempts |
| Lockout duration | 15 minutes, then escalating |
| Storage | Argon2id hash |
| Change frequency | User-initiated |

### 4.3 Implementation Reference

```go
// internal/infrastructure/security/crypto/encryption.go
func HashPassword(password string) (string, error) // Argon2id
func HashPIN(pin string) (string, error) // Argon2id
func VerifyPassword(password, encodedHash string) bool
func VerifyPIN(pin, encodedHash string) bool
```

## 5. Multi-Factor Authentication

### 5.1 MFA Requirements

| Account Type | MFA Required | Methods |
|--------------|--------------|---------|
| Admin/Privileged | Always | TOTP, Hardware key |
| Employee | Always | TOTP, SMS (backup) |
| Customer (standard) | Risk-based | OTP via SMS/Email |
| Customer (high-value) | Always | OTP, Biometric |
| API/Service | N/A | API keys + IP restriction |

### 5.2 MFA Triggers for Customers

MFA required when:
- Logging in from new device
- Logging in from new location
- Transaction above threshold (NGN 50,000)
- Changing sensitive settings
- Adding new payment methods

### 5.3 OTP Specifications

| Parameter | Value |
|-----------|-------|
| Length | 6 digits |
| Validity | 5 minutes |
| Attempts | 3 maximum |
| Rate limit | 3 per 5 minutes |
| Channel | SMS or Email |

## 6. Password Storage

### 6.1 Hashing Requirements

| Algorithm | Use Case | Parameters |
|-----------|----------|------------|
| Argon2id | User passwords | memory=64MB, iterations=3, parallelism=4 |
| Argon2id | Transaction PINs | memory=64MB, iterations=3, parallelism=4 |
| PBKDF2-SHA256 | Legacy migration only | iterations=310,000 |

### 6.2 Implementation Standards

```go
// Argon2id parameters (OWASP recommendations)
const (
    Argon2Memory     = 64 * 1024 // 64 MB
    Argon2Iterations = 3
    Argon2Parallelism = 4
    Argon2SaltLength  = 16
    Argon2KeyLength   = 32
)
```

### 6.3 Prohibited Practices

- NEVER store passwords in plain text
- NEVER use MD5 or SHA1 for passwords
- NEVER use the same salt for all passwords
- NEVER log password values
- NEVER transmit passwords in URLs

## 7. Password Management

### 7.1 Password Managers

**Approved Password Managers:**
- 1Password (Business)
- Bitwarden
- LastPass (Business)

**Requirements:**
- Master password meets policy
- MFA enabled
- Company-managed for business accounts
- Regular backup of vault

### 7.2 Password Vault (Infrastructure)

| Requirement | Standard |
|-------------|----------|
| Solution | HashiCorp Vault or equivalent |
| Access | Role-based, MFA required |
| Audit | All access logged |
| Rotation | Automated for supported systems |
| Backup | Encrypted, offsite |

### 7.3 Shared Credentials

- Avoid shared credentials where possible
- If required, store in approved vault only
- Audit access to shared credentials
- Change immediately when team changes
- Document all shared credential usage

## 8. Prohibited Passwords

### 8.1 Blocklist

Passwords must not be:
- Dictionary words
- Username or email variations
- Company name variations
- Common patterns (123456, qwerty, password)
- Previously breached passwords
- Personal information (DOB, phone, name)

### 8.2 Breach Detection

- Check new passwords against HaveIBeenPwned
- Block passwords found in breaches
- Notify users of compromised credentials

### 8.3 Common Password List

Maintain blocklist of:
- Top 10,000 common passwords
- Company-specific variations
- Industry-specific patterns

## 9. Password Reset

### 9.1 Self-Service Reset

| Step | Requirement |
|------|-------------|
| Identity verification | Email or SMS OTP |
| Reset link validity | 15 minutes |
| Link usage | Single use |
| New password | Must meet policy |
| Notification | Confirm via alternate channel |

### 9.2 Administrative Reset

| Scenario | Process |
|----------|---------|
| Employee | IT ticket, manager approval |
| Customer | Support verification, OTP |
| Emergency | On-call approval, documented |

### 9.3 Forced Reset

Force password change for:
- Suspected compromise
- Security incident
- Policy non-compliance
- Extended inactivity (180+ days)

## 10. Session Management

### 10.1 Session Parameters

| Parameter | Value |
|-----------|-------|
| Session timeout (idle) | 30 minutes |
| Session maximum | 24 hours |
| Concurrent sessions | 3 maximum |
| Session token length | 256 bits |
| Token storage | HTTP-only, Secure cookies |

### 10.2 Session Termination

Sessions must end when:
- User logs out
- Idle timeout reached
- Maximum duration reached
- Password changed
- Account disabled
- Suspicious activity detected

## 11. API Authentication

### 11.1 API Key Requirements

| Requirement | Standard |
|-------------|----------|
| Key length | 256 bits minimum |
| Key format | Random, cryptographically secure |
| Rotation | 90 days recommended |
| Storage | Never in code, use vault |
| Transmission | HTTPS only, in header |

### 11.2 JWT Token Requirements

| Parameter | Value |
|-----------|-------|
| Algorithm | RS256 or ES256 |
| Access token expiry | 15 minutes |
| Refresh token expiry | 7 days |
| Refresh rotation | On each use |
| Claims | Minimal, no sensitive data |

### 11.3 Service-to-Service

| Requirement | Standard |
|-------------|----------|
| Authentication | Mutual TLS or JWT |
| Key management | Automated rotation |
| Access control | Least privilege |
| Audit | All requests logged |

## 12. Compliance

### 12.1 Regulatory Alignment

| Standard | Requirement | Compliance |
|----------|-------------|------------|
| PCI DSS 8.3.1 | Strong passwords | Met |
| PCI DSS 8.3.6 | Password history | Met |
| NIST 800-63B | Modern guidelines | Met |
| OWASP | Argon2id hashing | Met |

### 12.2 Audit Requirements

Log and retain:
- Authentication attempts (success/failure)
- Password changes
- Account lockouts
- Reset requests
- MFA events

Retention: 7 years

## 13. User Responsibilities

### 13.1 Do's

- Use unique passwords for each account
- Use a password manager
- Enable MFA where available
- Report suspicious activity
- Lock workstation when away

### 13.2 Don'ts

- Share passwords with anyone
- Write passwords down
- Use personal passwords for work
- Use work passwords for personal
- Send passwords via email/chat
- Store passwords in plain text

## 14. Enforcement

### 14.1 Technical Controls

- Password complexity enforced at creation
- Password history enforced
- Account lockout automated
- Session timeout enforced
- MFA required per policy

### 14.2 Policy Violations

| Violation | Consequence |
|-----------|-------------|
| Weak password | Forced reset |
| Password sharing | Warning, then suspension |
| Bypassing controls | Suspension, investigation |
| Repeated violations | Termination consideration |

---

## Appendix A: Password Strength Meter

| Score | Strength | Requirements Met |
|-------|----------|------------------|
| 0-20 | Weak | Rejected |
| 21-40 | Fair | Rejected |
| 41-60 | Good | Minimum acceptable |
| 61-80 | Strong | Recommended |
| 81-100 | Excellent | Ideal |

**Scoring Factors:**
- Length (longer = stronger)
- Character variety
- Uncommon patterns
- Not in breach database

## Appendix B: Password Generation Guidelines

**Strong Password Examples (Pattern):**
- 3-4 random words + numbers + symbol
- Passphrase with substitutions
- Generated by password manager (preferred)

**Example Patterns:**
- `Correct-Horse-Battery-Staple-42!`
- `Mountain$River7Sunset#Cloud`
- `[Generated: Xk9#mP2$vL4@nQ7&]`

## Appendix C: Implementation Checklist

### Backend
- [ ] Argon2id hashing implemented
- [ ] Password policy enforcement
- [ ] Breach database check
- [ ] Account lockout mechanism
- [ ] Password history tracking
- [ ] Secure reset process
- [ ] Audit logging

### Frontend
- [ ] Password strength meter
- [ ] Policy requirements display
- [ ] Clear error messages
- [ ] Secure password input
- [ ] MFA integration

## Appendix D: Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [Date] | [Author] | Initial release |
