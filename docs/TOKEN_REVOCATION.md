# Token Revocation Mechanism

## Overview

HustleX implements a comprehensive token revocation mechanism to address **Security Issue #1** from the security audit. This system allows immediate invalidation of JWT access tokens, protecting against token theft and unauthorized access.

## Problem Statement

Previously, JWT tokens remained valid until expiration, even after logout. This created security risks:
- Stolen tokens could be used until expiry (typically 15 minutes)
- Compromised accounts couldn't be immediately locked
- No way to invalidate sessions across devices
- Password/PIN changes didn't revoke existing tokens

## Solution Architecture

### Two-Level Revocation System

#### 1. Token-Level Blacklist
Individual tokens can be revoked (e.g., on logout):
```
Redis Key: token:blacklist:<sha256_hash_of_token>
TTL: Remaining token lifetime
Value: Revocation timestamp
```

#### 2. User-Level Revocation
All tokens for a user can be revoked (e.g., on password change):
```
Redis Key: user:tokens:revoked:<user_id>
TTL: 30 days (covers maximum token lifetime)
Value: Revocation timestamp
```

Tokens issued before the revocation timestamp are considered invalid.

## Implementation Details

### TokenBlacklistService

Located at: `apps/api/internal/services/token_blacklist.go`

**Core Methods:**
- `BlacklistToken(token, expiresAt)` - Blacklist a specific token
- `IsTokenBlacklisted(token)` - Check if a token is blacklisted
- `BlacklistAllUserTokens(userID)` - Revoke all tokens for a user
- `IsUserTokenRevoked(userID, tokenIssuedAt)` - Check if user's tokens are revoked
- `ClearUserTokenRevocation(userID)` - Re-enable an account
- `GetBlacklistStats()` - Monitoring and debugging

**Security Features:**
- Tokens are hashed (SHA-256) before storage - full tokens never stored in Redis
- Automatic TTL expiration - blacklist entries auto-delete when token would expire anyway
- Fail-open design - if Redis is unavailable, requests proceed (availability over security)
- Constant-time hash comparison to prevent timing attacks

### Authentication Middleware Integration

**Fiber Middleware** (`apps/api/internal/middleware/auth.go`):
- `AuthMiddleware` - Now checks blacklist after JWT validation
- `OptionalAuthMiddleware` - Checks blacklist but doesn't fail request

**Validation Flow:**
1. Extract token from Authorization header
2. Validate JWT signature and expiry
3. ✨ **NEW:** Check if token is blacklisted
4. ✨ **NEW:** Check if all user tokens are revoked
5. Allow request or return 401 Unauthorized

### AuthService Updates

**Logout Enhancement** (`apps/api/internal/services/auth_service.go`):
```go
func (s *AuthService) Logout(ctx, userID, accessToken) error {
    // Delete refresh token (existing behavior)
    s.redis.Del(ctx, refreshKey)

    // NEW: Blacklist access token
    s.blacklist.BlacklistToken(ctx, accessToken, expiresAt)
}
```

**New Method - Revoke All Tokens:**
```go
func (s *AuthService) RevokeAllUserTokens(ctx, userID) error {
    s.redis.Del(ctx, refreshKey)
    s.blacklist.BlacklistAllUserTokens(ctx, userID)
}
```

**Token Validation Enhancement:**
```go
func (s *AuthService) ValidateAccessToken(token) (*Claims, error) {
    // Validate JWT signature
    claims := parseJWT(token)

    // NEW: Check blacklist
    if blacklisted, _ := s.blacklist.IsTokenBlacklisted(token) {
        return ErrTokenRevoked
    }

    // NEW: Check user revocation
    if revoked, _ := s.blacklist.IsUserTokenRevoked(userID, issuedAt) {
        return ErrTokenRevokedSecurityAction
    }

    return claims
}
```

## Usage Examples

### Logout (Single Device)

```go
// User logs out - revoke their current access token
err := authService.Logout(ctx, userID, accessToken)
```

**What happens:**
1. Refresh token deleted from Redis
2. Access token added to blacklist (until its natural expiry)
3. Subsequent API calls with that token return 401

### Password/PIN Change (All Devices)

```go
// User changes password - revoke ALL their tokens
err := authService.RevokeAllUserTokens(ctx, userID)

// Then generate new tokens for current session
newTokens, _ := authService.GenerateTokens(user)
```

**What happens:**
1. All refresh tokens deleted
2. User-level revocation marker set
3. All tokens issued before this moment are invalid
4. User must re-authenticate on all devices

### Account Suspension (Admin Action)

```go
// Admin suspends user account
err := authService.RevokeAllUserTokens(ctx, suspendedUserID)

// Update user status
db.Model(&user).Update("is_active", false)
```

### Re-enabling Account

```go
// Admin re-enables account
err := blacklistService.ClearUserTokenRevocation(ctx, userID)
db.Model(&user).Update("is_active", true)
```

## Security Considerations

### Defense in Depth
- Tokens hashed before storage (even if Redis is compromised, tokens aren't leaked)
- Separate refresh token deletion (belt and suspenders)
- Middleware validation at multiple layers

### Performance
- Redis lookup adds ~1-2ms per request (negligible)
- Hash computation is fast (SHA-256)
- TTL-based auto-expiry prevents unbounded growth

### Availability vs Security Trade-off
**Current Implementation: Fail Open**
- If Redis is unavailable, requests proceed
- Rationale: Platform availability > temporary security gap
- Alternative: Fail closed (more secure, less available)

**To change to fail-closed:**
```go
isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, token)
if err != nil {
    return fiber.StatusServiceUnavailable // Fail closed
}
```

### Token Lifetime Strategy
- Access tokens: 15 minutes (short-lived)
- Refresh tokens: 7 days
- Shorter access token lifetime reduces revocation urgency
- Balance between security (shorter) and UX (fewer refreshes)

## Monitoring and Operations

### Blacklist Statistics

```go
stats, _ := blacklistService.GetBlacklistStats(ctx)
// Returns:
// {
//   "blacklisted_tokens": 1234,
//   "user_revocations": 56
// }
```

**Recommended Monitoring:**
- Alert if `blacklisted_tokens` grows rapidly (potential attack)
- Track `user_revocations` for abuse detection
- Monitor Redis memory usage for blacklist keys

### Redis Key Patterns

```bash
# View all blacklisted tokens
redis-cli KEYS "token:blacklist:*"

# View all user revocations
redis-cli KEYS "user:tokens:revoked:*"

# Check specific token (you need the hash)
redis-cli GET "token:blacklist:<hash>"

# Check user revocation
redis-cli GET "user:tokens:revoked:<user_id>"
```

### Cleanup and Maintenance

**No manual cleanup required!** Redis TTL handles expiration:
- Blacklisted tokens auto-delete when they would expire anyway
- User revocations expire after 30 days
- Total memory impact: ~100 bytes × active blacklist size

**Estimated Memory Usage:**
- 10,000 blacklisted tokens ≈ 1MB Redis memory
- 1,000 user revocations ≈ 100KB Redis memory

## Testing

Comprehensive test suite at: `apps/api/internal/services/token_blacklist_test.go`

**Test Coverage:**
- Token blacklisting and checking
- TTL expiration behavior
- User-level revocation
- Concurrent operations
- Hash determinism
- Edge cases (expired tokens, non-existent users)

**Run tests:**
```bash
cd apps/api
go test ./internal/services -v -run TestTokenBlacklist
```

## Migration Guide

### For Developers

**Before (Logout):**
```go
func LogoutHandler(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uuid.UUID)
    authService.Logout(ctx, userID)
    return c.JSON(fiber.Map{"message": "Logged out"})
}
```

**After (Logout):**
```go
func LogoutHandler(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uuid.UUID)

    // Extract access token
    authHeader := c.Get("Authorization")
    token := strings.Split(authHeader, " ")[1]

    authService.Logout(ctx, userID, token)
    return c.JSON(fiber.Map{"message": "Logged out"})
}
```

**Middleware Setup:**
```go
// Before
app.Use(middleware.AuthMiddleware(cfg))

// After
app.Use(middleware.AuthMiddleware(cfg, redisClient))
```

### For Operations

**No infrastructure changes required!**
- Uses existing Redis instance
- No new services or databases
- Backward compatible (old tokens continue to work)

## Future Enhancements

### Potential Improvements
1. **Admin Dashboard:**
   - View active blacklist entries
   - Manually revoke tokens
   - Audit log of revocations

2. **Refresh Token Rotation:**
   - Invalidate old refresh token on use
   - Detect token replay attacks

3. **Device Management:**
   - Track active sessions per device
   - Selective revocation by device
   - "Logout everywhere except this device"

4. **Anomaly Detection:**
   - Detect suspicious token usage patterns
   - Auto-revoke on anomaly detection
   - Notify user of suspicious activity

5. **Distributed Rate Limiting:**
   - Rate limit revocation attempts
   - Prevent DoS via mass revocation

## References

- [Security Audit Report](./SECURITY_AUDIT_REPORT.md) - Issue #1
- [PRD Section 10.3](./PRD.md) - Security Hardening Task
- [OWASP JWT Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)
- [RFC 7519 - JSON Web Tokens](https://datatracker.ietf.org/doc/html/rfc7519)

---

**Status:** ✅ Implemented
**Security Issue:** #1 - Token Revocation Mechanism
**Severity:** HIGH
**Date Completed:** 2026-02-06
**Author:** Claude Sonnet 4.5
