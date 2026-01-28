# ADR-007: JWT-Based Authentication Strategy

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires secure authentication for:
- Mobile app users (primary)
- Future web dashboard users
- API integrations

Requirements:
1. Stateless authentication (scalable across multiple API instances)
2. Support for mobile app sessions
3. Secure token refresh mechanism
4. Revocation capability for compromised tokens
5. Integration with Nigerian phone numbers (primary identifier)

## Decision

We chose **JWT (JSON Web Tokens)** with a dual-token strategy (access + refresh tokens).

### Token Strategy:

1. **Access Token**: Short-lived (15 minutes), used for API authentication
2. **Refresh Token**: Long-lived (7 days), used to obtain new access tokens
3. **Redis Blacklist**: For immediate token revocation

### Key Reasons:

1. **Stateless**: No server-side session storage required (scales horizontally).

2. **Mobile-Friendly**: Tokens can be stored in secure storage on device.

3. **Self-Contained**: User claims embedded in token (user_id, tier, permissions).

4. **Industry Standard**: Well-understood, extensive library support.

5. **Flexible Expiration**: Short access tokens limit exposure; long refresh tokens improve UX.

## Consequences

### Positive

- **Scalability**: Any API instance can validate tokens without database lookup
- **Performance**: No session lookup on every request
- **Flexibility**: Claims can include custom data (user tier, permissions)
- **Cross-platform**: Same token works for mobile, web, integrations
- **Standard tooling**: Libraries available in all languages

### Negative

- **Token size**: JWTs larger than simple session IDs (~500 bytes)
- **Revocation complexity**: Requires blacklist for immediate revocation
- **Key management**: Secret rotation requires coordination
- **No built-in refresh**: Must implement refresh logic manually

### Neutral

- Requires HTTPS for secure transmission
- Token storage responsibility on client
- Clock skew handling needed for distributed systems

## Implementation Details

### Token Structure

```go
// Access Token Claims
type AccessTokenClaims struct {
    jwt.RegisteredClaims
    UserID    string `json:"user_id"`
    Phone     string `json:"phone"`
    Tier      string `json:"tier"`      // bronze, silver, gold, platinum
    TokenType string `json:"token_type"` // "access"
}

// Refresh Token Claims
type RefreshTokenClaims struct {
    jwt.RegisteredClaims
    UserID    string `json:"user_id"`
    TokenType string `json:"token_type"` // "refresh"
    Family    string `json:"family"`     // Token family for rotation
}
```

### Token Lifetimes

| Token Type | Lifetime | Storage |
|------------|----------|---------|
| Access Token | 15 minutes | Memory (mobile) |
| Refresh Token | 7 days | Secure Storage (mobile) |

### Authentication Flow

```
1. User requests OTP
   POST /api/v1/auth/otp/request
   Body: { "phone": "+2348012345678" }

2. User verifies OTP
   POST /api/v1/auth/otp/verify
   Body: { "phone": "+2348012345678", "code": "123456" }
   Response: { "access_token": "...", "refresh_token": "...", "user": {...} }

3. User accesses protected resource
   GET /api/v1/wallet
   Header: Authorization: Bearer <access_token>

4. Access token expires, refresh
   POST /api/v1/auth/refresh
   Body: { "refresh_token": "..." }
   Response: { "access_token": "...", "refresh_token": "..." }

5. User logs out
   POST /api/v1/auth/logout
   (Blacklists both tokens)
```

### Token Refresh Strategy

```go
// Refresh token rotation (prevents replay attacks)
func RefreshTokens(refreshToken string) (*TokenPair, error) {
    claims, err := ValidateRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }

    // Check if token family is blacklisted (possible theft)
    if IsTokenFamilyBlacklisted(claims.Family) {
        // Revoke all tokens for this user
        RevokeAllUserTokens(claims.UserID)
        return nil, ErrTokenCompromised
    }

    // Blacklist old refresh token
    BlacklistToken(refreshToken, claims.ExpiresAt)

    // Issue new token pair with same family
    return GenerateTokenPair(claims.UserID, claims.Family)
}
```

### Token Revocation (Redis Blacklist)

```go
// Blacklist token until expiry
func BlacklistToken(token string, expiresAt time.Time) error {
    ttl := time.Until(expiresAt)
    return redis.Set(ctx,
        "jwt:blacklist:"+getTokenJTI(token),
        "1",
        ttl,
    ).Err()
}

// Check if token is blacklisted
func IsTokenBlacklisted(token string) bool {
    exists, _ := redis.Exists(ctx, "jwt:blacklist:"+getTokenJTI(token)).Result()
    return exists > 0
}
```

### Security Measures

1. **HMAC-SHA256 Signing**: Industry-standard algorithm
2. **Unique JTI**: Each token has unique identifier for revocation
3. **Token Binding**: Refresh tokens bound to device fingerprint (future)
4. **Family Tracking**: Detect refresh token reuse attacks
5. **Secure Storage**: Mobile app uses flutter_secure_storage

## Alternatives Considered

### Alternative 1: Session-Based Authentication

**Pros**: Simple revocation, smaller tokens, server-controlled
**Cons**: Requires session store, doesn't scale horizontally without shared storage

**Rejected because**: Requires Redis/database lookup on every request; JWT is more scalable.

### Alternative 2: OAuth 2.0 with External Provider

**Pros**: Delegated authentication, social login support
**Cons**: Complex setup, external dependency, may not support Nigerian phone auth

**Rejected because**: OTP via phone is primary auth method; OAuth adds unnecessary complexity.

### Alternative 3: API Keys

**Pros**: Simple, no expiration management
**Cons**: No user context, revocation requires database update, not suitable for mobile apps

**Rejected because**: Need user-specific authentication with claims and expiration.

### Alternative 4: Paseto (Platform-Agnostic Security Tokens)

**Pros**: Modern, no algorithm confusion, stateless verification
**Cons**: Smaller ecosystem, less library support, team unfamiliarity

**Rejected because**: JWT has broader tooling support and team familiarity.

## References

- [RFC 7519 - JSON Web Token](https://tools.ietf.org/html/rfc7519)
- [JWT Best Practices](https://auth0.com/blog/a-look-at-the-latest-draft-for-jwt-bcp/)
- [golang-jwt Library](https://github.com/golang-jwt/jwt)
- [Refresh Token Rotation](https://auth0.com/docs/secure/tokens/refresh-tokens/refresh-token-rotation)
