# ADR-012: OTP-Based User Authentication

## Status

Accepted

## Date

2024-01-15

## Context

HustleX targets the Nigerian market where:
- Phone ownership is near-universal (95%+ mobile penetration)
- Many users don't have email addresses
- Password fatigue is common (users forget passwords)
- SMS is reliable and trusted
- Financial apps require strong authentication

Requirements:
1. Low friction onboarding
2. Secure authentication
3. No password management burden
4. Works on basic feature phones (USSD fallback)
5. Supports 2FA for sensitive operations

## Decision

We chose **OTP (One-Time Password) via SMS** as the primary authentication method, with phone number as the primary identifier.

### Authentication Flow:

1. User enters phone number
2. System sends 6-digit OTP via SMS
3. User enters OTP to verify
4. System issues JWT tokens
5. Transaction PIN for sensitive operations

### Key Reasons:

1. **Universal Access**: Every user has a phone number; not all have email.

2. **User Experience**: No password to remember, quick login flow.

3. **Nigerian Market Fit**: SMS is trusted, WhatsApp/mobile culture dominant.

4. **Security**: OTP changes every time, no credential stuffing possible.

5. **Regulatory Alignment**: BVN/NIN verification links to phone numbers.

## Consequences

### Positive

- **Zero password issues**: No password resets, no weak passwords
- **Quick onboarding**: 2-step login vs traditional username/password
- **Familiar flow**: Users accustomed to OTP from banking apps
- **Multi-device friendly**: Same phone number works on any device
- **Strong authentication**: Time-limited codes, device-independent

### Negative

- **SMS cost**: Each login/verification costs ₦2-5 (0.5-1 cent USD)
- **SMS reliability**: Delivery delays during network congestion
- **SIM swap attacks**: Phone number takeover is possible
- **No offline login**: Requires SMS network availability
- **International roaming**: Users abroad may not receive SMS

### Neutral

- Requires SMS gateway integration (Termii, Africa's Talking)
- Rate limiting needed to prevent abuse
- OTP expiration window balance (security vs UX)

## Implementation Details

### OTP Configuration

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Code Length | 6 digits | Balance security and usability |
| Expiration | 5 minutes | Short enough for security, long enough for UX |
| Max Attempts | 5 per code | Prevent brute force |
| Resend Cooldown | 60 seconds | Prevent SMS flooding |
| Daily Limit | 10 per phone | Prevent abuse |

### Database Model

```go
type OTPCode struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
    Phone       string    `gorm:"index;size:20"`
    Code        string    `gorm:"size:6"`
    Purpose     string    `gorm:"size:50"` // login, register, pin_reset, transaction
    Attempts    int       `gorm:"default:0"`
    IsVerified  bool      `gorm:"default:false"`
    ExpiresAt   time.Time `gorm:"index"`
    CreatedAt   time.Time
}
```

### Request Flow

```
┌──────────┐      ┌──────────┐      ┌──────────┐      ┌──────────┐
│  Mobile  │      │   API    │      │  Redis   │      │   SMS    │
│   App    │      │  Server  │      │  Cache   │      │ Gateway  │
└────┬─────┘      └────┬─────┘      └────┬─────┘      └────┬─────┘
     │                 │                 │                 │
     │ POST /otp/request                 │                 │
     │ {phone: "+234..."}                │                 │
     │────────────────>│                 │                 │
     │                 │                 │                 │
     │                 │ Check rate limit│                 │
     │                 │────────────────>│                 │
     │                 │<────────────────│                 │
     │                 │                 │                 │
     │                 │ Generate OTP    │                 │
     │                 │ Store in Redis  │                 │
     │                 │────────────────>│                 │
     │                 │                 │                 │
     │                 │ Send SMS        │                 │
     │                 │────────────────────────────────────>
     │                 │                 │                 │
     │ 200 OK          │                 │                 │
     │ {expires_in: 300}                 │                 │
     │<────────────────│                 │                 │
     │                 │                 │                 │
     │ POST /otp/verify                  │                 │
     │ {phone, code}   │                 │                 │
     │────────────────>│                 │                 │
     │                 │                 │                 │
     │                 │ Verify OTP      │                 │
     │                 │────────────────>│                 │
     │                 │<────────────────│                 │
     │                 │                 │                 │
     │ 200 OK          │ Delete OTP      │                 │
     │ {access_token,  │────────────────>│                 │
     │  refresh_token} │                 │                 │
     │<────────────────│                 │                 │
```

### OTP Generation

```go
func GenerateOTP(phone string, purpose string) error {
    // Rate limiting check
    key := fmt.Sprintf("otp:ratelimit:%s", phone)
    count, _ := redis.Incr(ctx, key).Result()
    if count == 1 {
        redis.Expire(ctx, key, 24*time.Hour)
    }
    if count > 10 {
        return ErrRateLimitExceeded
    }

    // Generate secure random 6-digit code
    code := generateSecureCode(6)

    // Store OTP with expiration
    otpKey := fmt.Sprintf("otp:%s:%s", phone, purpose)
    data := OTPData{
        Code:      code,
        Attempts:  0,
        CreatedAt: time.Now(),
    }
    redis.Set(ctx, otpKey, data, 5*time.Minute)

    // Send via SMS gateway
    return smsService.Send(phone, fmt.Sprintf(
        "Your HustleX code is %s. Valid for 5 minutes. Do not share.",
        code,
    ))
}

func generateSecureCode(length int) string {
    const digits = "0123456789"
    b := make([]byte, length)
    rand.Read(b)
    for i := range b {
        b[i] = digits[int(b[i])%len(digits)]
    }
    return string(b)
}
```

### OTP Verification

```go
func VerifyOTP(phone string, code string, purpose string) (*User, error) {
    otpKey := fmt.Sprintf("otp:%s:%s", phone, purpose)

    // Get OTP data
    var data OTPData
    err := redis.Get(ctx, otpKey).Scan(&data)
    if err == redis.Nil {
        return nil, ErrOTPExpired
    }

    // Check attempts
    if data.Attempts >= 5 {
        redis.Del(ctx, otpKey) // Delete after max attempts
        return nil, ErrMaxAttemptsExceeded
    }

    // Verify code
    if data.Code != code {
        data.Attempts++
        redis.Set(ctx, otpKey, data, redis.KeepTTL)
        return nil, ErrInvalidOTP
    }

    // Success - delete OTP
    redis.Del(ctx, otpKey)

    // Get or create user
    user, err := userService.FindOrCreateByPhone(phone)
    if err != nil {
        return nil, err
    }

    return user, nil
}
```

### SMS Gateway Integration

```go
// Termii SMS Gateway
type TermiiService struct {
    apiKey   string
    senderID string
    baseURL  string
}

func (s *TermiiService) Send(phone string, message string) error {
    payload := map[string]interface{}{
        "api_key": s.apiKey,
        "to":      phone,
        "from":    s.senderID,
        "sms":     message,
        "type":    "plain",
        "channel": "generic",
    }

    resp, err := http.Post(s.baseURL+"/sms/send", "application/json",
        bytes.NewBuffer(jsonEncode(payload)))

    if err != nil || resp.StatusCode != 200 {
        // Queue for retry
        return asynq.EnqueueSMSRetry(phone, message)
    }

    return nil
}
```

## Security Measures

### 1. Rate Limiting

```go
// Per-phone limits
- 10 OTP requests per day
- 60 second cooldown between requests
- 5 verification attempts per OTP

// Per-IP limits (prevent enumeration)
- 100 OTP requests per hour per IP
- Increasing delay after failures
```

### 2. Brute Force Prevention

```go
// Exponential backoff after failures
attempts := getFailedAttempts(phone)
if attempts > 0 {
    delay := time.Duration(math.Pow(2, float64(attempts))) * time.Second
    if delay > maxDelay {
        delay = maxDelay
    }
    time.Sleep(delay)
}
```

### 3. SIM Swap Detection

```go
// Monitor for suspicious patterns
- Multiple failed OTPs followed by success
- Login from new device after phone inactivity
- Require additional verification (PIN) for high-risk actions
```

### 4. Transaction PIN

```go
// Secondary factor for sensitive operations
type TransactionPIN struct {
    UserID    uuid.UUID
    PINHash   string // bcrypt hashed
    FailedAt  []time.Time
    LockedAt  *time.Time
}

// Required for: transfers > ₦10,000, withdrawals, loan applications
func VerifyTransactionPIN(userID uuid.UUID, pin string) error {
    // Check lock status
    // Verify bcrypt hash
    // Track failures
}
```

## Alternatives Considered

### Alternative 1: Email + Password

**Pros**: No SMS cost, standard approach, offline login
**Cons**: Many Nigerian users lack email, password fatigue, support burden

**Rejected because**: Target users may not have email addresses.

### Alternative 2: Social Login (Google/Facebook)

**Pros**: No credential management, trusted providers
**Cons**: Requires social accounts, privacy concerns, dependency

**Rejected because**: Not universal in target demographic.

### Alternative 3: TOTP (Authenticator Apps)

**Pros**: No SMS cost, offline capable, more secure
**Cons**: Requires setup, authenticator app needed, complex for users

**Rejected because**: Too complex for target users; can be added as optional 2FA later.

### Alternative 4: WhatsApp OTP

**Pros**: Free delivery, higher delivery rates, rich messaging
**Cons**: Requires WhatsApp Business API, user must have WhatsApp

**Decision**: Add as secondary channel alongside SMS for cost optimization.

## References

- [NIST Digital Identity Guidelines](https://pages.nist.gov/800-63-3/)
- [Termii API Documentation](https://developers.termii.com/)
- [Africa's Talking SMS API](https://africastalking.com/sms)
- [SIM Swap Fraud Prevention](https://www.gsma.com/security/sim-swap-fraud/)
