# CSRF Protection Implementation

**Document Version:** 1.0
**Last Updated:** 2026-02-06
**Author:** Claude Sonnet 4.5
**Status:** Implemented

---

## Overview

This document describes the Cross-Site Request Forgery (CSRF) protection implementation for the HustleX API. CSRF protection is a critical security control that prevents attackers from tricking authenticated users into performing unwanted actions.

### What is CSRF?

Cross-Site Request Forgery (CSRF) is an attack that forces an authenticated user to execute unwanted actions on a web application. For example:
- An attacker could embed a hidden form on `evil.com` that transfers money from the victim's HustleX wallet
- When the victim visits `evil.com` while logged into HustleX, the form auto-submits
- Without CSRF protection, the browser sends the user's session cookies, and the transfer succeeds

### Implementation Approach

HustleX implements CSRF protection using the **Synchronizer Token Pattern** recommended by OWASP:
1. Server generates a unique, unpredictable CSRF token for each user session
2. Token is sent to client in a cookie (readable by JavaScript)
3. Client includes token in `X-CSRF-Token` header for state-changing requests
4. Server validates that the header token matches the expected token for the user
5. Tokens are rotated after each successful request for enhanced security

---

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────────┐
│                    CSRF Protection Flow                      │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. User authenticates → Auth middleware sets userID         │
│                                                               │
│  2. GET /api/profile                                         │
│     → CSRF middleware generates token                        │
│     → Sets cookie: csrf_token=abc123                         │
│     → Sets header: X-CSRF-Token: abc123                      │
│     → Stores token in memory: {userID: "user123", token: ...}│
│                                                               │
│  3. POST /api/wallet/transfer                                │
│     → Client sends header: X-CSRF-Token: abc123              │
│     → CSRF middleware validates token against userID         │
│     → If valid: allow request + rotate token                 │
│     → If invalid: return 403 Forbidden                       │
│                                                               │
│  4. Logout → RevokeCSRFToken() deletes token from store      │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

### Files

- **Implementation:** `apps/api/internal/interface/http/middleware/csrf.go`
- **Tests:** `apps/api/internal/interface/http/middleware/csrf_test.go`
- **Documentation:** `docs/CSRF_PROTECTION.md` (this file)

---

## Configuration

### Default Configuration

```go
CSRFConfig{
    TokenLifetime: 4 * time.Hour,        // Token expires after 4 hours
    CookiePath:    "/",                   // Cookie available for all paths
    CookieDomain:  "",                    // Current domain only
    Secure:        true,                  // HTTPS only (must be true in production)
    SameSite:      http.SameSiteStrictMode, // Strongest CSRF protection
    SkipPaths: []string{                  // Paths that don't need CSRF validation
        "/health",
        "/metrics",
        "/api/v1/auth/login",             // Login uses credentials, not session
        "/api/v1/auth/register",
        "/api/v1/webhooks/",              // Webhooks use signature verification
    },
}
```

### Environment-Specific Settings

**Development:**
```go
config := DefaultCSRFConfig()
config.Secure = false  // Allow HTTP (localhost)
config.SameSite = http.SameSiteLaxMode
```

**Production:**
```go
config := DefaultCSRFConfig()
config.Secure = true  // HTTPS only
config.SameSite = http.SameSiteStrictMode
config.CookieDomain = ".hustlex.ng"  // Allow subdomains
```

---

## Usage

### Backend Integration

#### 1. Initialize CSRF Store

```go
// In main.go or router setup
csrfConfig := middleware.DefaultCSRFConfig()
csrfConfig.Secure = true // Production

csrfStore := middleware.NewInMemoryCSRFStore(csrfConfig)
```

#### 2. Add to Middleware Chain

```go
// Apply CSRF protection after authentication
router.Use(
    middleware.EnhancedRequestID,
    middleware.EnhancedSecurityHeaders,
    middleware.Auth(jwtService),            // Must come before CSRF
    middleware.CSRFProtection(csrfStore, csrfConfig),
    middleware.RateLimiting(rateLimiter),
)
```

**Important:** CSRF middleware must come **after** authentication middleware because it needs the `userID` from context.

#### 3. Revoke on Logout

```go
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    // Revoke JWT token (existing)
    token := getTokenFromRequest(r)
    h.tokenBlacklist.BlacklistToken(r.Context(), token, expiresAt)

    // Revoke CSRF token (new)
    middleware.RevokeCSRFToken(w, r, h.csrfStore, h.csrfConfig)

    // Return success
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
```

### Frontend Integration (Mobile/Web)

#### Flutter (Mobile App)

```dart
class ApiClient {
  String? _csrfToken;

  Future<void> makeRequest(String method, String path, Map<String, dynamic> data) async {
    final headers = {
      'Authorization': 'Bearer $accessToken',
      'Content-Type': 'application/json',
    };

    // For state-changing requests, include CSRF token
    if (['POST', 'PUT', 'PATCH', 'DELETE'].contains(method)) {
      if (_csrfToken == null) {
        // Fetch CSRF token first with a GET request
        await _fetchCSRFToken();
      }
      headers['X-CSRF-Token'] = _csrfToken!;
    }

    final response = await http.request(
      method,
      Uri.parse('$baseUrl$path'),
      headers: headers,
      body: jsonEncode(data),
    );

    // Update CSRF token from response (token rotation)
    final newToken = response.headers['x-csrf-token'];
    if (newToken != null) {
      _csrfToken = newToken;
    }

    if (response.statusCode == 403 && response.body.contains('CSRF')) {
      // Token invalid, fetch new one and retry
      await _fetchCSRFToken();
      return makeRequest(method, path, data); // Retry once
    }

    return response;
  }

  Future<void> _fetchCSRFToken() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/v1/users/me'),
      headers: {'Authorization': 'Bearer $accessToken'},
    );
    _csrfToken = response.headers['x-csrf-token'];
  }
}
```

#### React (Admin Dashboard)

```javascript
class ApiService {
  constructor() {
    this.csrfToken = null;
  }

  async makeRequest(method, path, data) {
    const headers = {
      'Authorization': `Bearer ${getAccessToken()}`,
      'Content-Type': 'application/json',
    };

    // For state-changing requests, include CSRF token
    if (['POST', 'PUT', 'PATCH', 'DELETE'].includes(method)) {
      if (!this.csrfToken) {
        await this.fetchCSRFToken();
      }
      headers['X-CSRF-Token'] = this.csrfToken;
    }

    const response = await fetch(`${API_BASE_URL}${path}`, {
      method,
      headers,
      body: data ? JSON.stringify(data) : undefined,
      credentials: 'include', // Important for cookies
    });

    // Update CSRF token from response (token rotation)
    const newToken = response.headers.get('X-CSRF-Token');
    if (newToken) {
      this.csrfToken = newToken;
    }

    if (response.status === 403) {
      const error = await response.json();
      if (error.code === 'CSRF_TOKEN_INVALID') {
        // Token invalid, fetch new one and retry
        await this.fetchCSRFToken();
        return this.makeRequest(method, path, data); // Retry once
      }
    }

    return response;
  }

  async fetchCSRFToken() {
    const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
      headers: { 'Authorization': `Bearer ${getAccessToken()}` },
      credentials: 'include',
    });
    this.csrfToken = response.headers.get('X-CSRF-Token');
  }
}
```

---

## Security Considerations

### ✅ What This Protects Against

1. **Classic CSRF Attacks**: Prevents malicious sites from making authenticated requests
2. **Clickjacking + CSRF**: Combined with `X-Frame-Options: DENY`, prevents UI redressing attacks
3. **Confused Deputy**: Validates that requests originate from the legitimate client application

### ✅ Security Features

1. **Cryptographically Secure Tokens**: Uses UUIDv4 for unpredictable tokens
2. **Constant-Time Comparison**: Prevents timing attacks when validating tokens
3. **Token Rotation**: Issues new token after each request (reduces window of exposure)
4. **Token Expiration**: Tokens expire after 4 hours (configurable)
5. **SameSite=Strict**: Cookie policy prevents cross-site submission
6. **HTTPS Only**: Secure flag prevents token leakage over HTTP

### ⚠️ Limitations

1. **Does Not Protect Against XSS**: If attacker can execute JavaScript, they can read the CSRF token
   - **Mitigation**: Comprehensive XSS protection (input validation, CSP, output encoding)

2. **Does Not Protect Unauthenticated Endpoints**: CSRF protection only applies to authenticated requests
   - **Mitigation**: Rate limiting, CAPTCHA for sensitive public endpoints

3. **Subdomain Cookie Sharing**: If `CookieDomain=.hustlex.ng`, subdomains can read token
   - **Mitigation**: Only set domain for trusted subdomains, or use stricter domain policy

### 🔒 Best Practices

1. **Always Use HTTPS in Production**: Set `Secure: true`
2. **Use SameSite=Strict**: Strongest protection, but may break some legitimate cross-origin flows
3. **Rotate Tokens Frequently**: Current implementation rotates on every state-changing request
4. **Revoke on Logout**: Always call `RevokeCSRFToken()` when user logs out
5. **Monitor 403 Errors**: Spike in CSRF errors may indicate attack or client misconfiguration
6. **Never Log CSRF Tokens**: Tokens are secrets, don't log them in plain text

---

## Testing

### Unit Tests

Run the comprehensive test suite:

```bash
cd apps/api
go test -v ./internal/interface/http/middleware -run TestCSRF
```

**Test Coverage:**
- ✅ Token generation and validation
- ✅ Token expiration
- ✅ Token revocation
- ✅ Safe methods (GET, HEAD, OPTIONS) bypass CSRF
- ✅ Unsafe methods (POST, PUT, DELETE, PATCH) require CSRF token
- ✅ Unauthenticated requests bypass CSRF
- ✅ Skip paths configuration
- ✅ Token rotation after successful request
- ✅ Invalid token rejection
- ✅ Constant-time comparison

### Integration Testing

```bash
# Test full flow with authenticated user
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone":"+2348012345678","password":"test123"}' \
  -c cookies.txt

# GET request receives CSRF token
curl -X GET http://localhost:8081/api/v1/users/me \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -b cookies.txt \
  -v  # Check X-CSRF-Token header in response

# POST request with CSRF token succeeds
curl -X POST http://localhost:8081/api/v1/wallet/transfer \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to":"user456","amount":1000}' \
  -b cookies.txt

# POST request without CSRF token fails
curl -X POST http://localhost:8081/api/v1/wallet/transfer \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to":"user456","amount":1000}' \
  -b cookies.txt
# Expected: 403 Forbidden with CSRF error
```

### Load Testing

Verify performance impact of CSRF validation:

```bash
# Benchmark token generation
go test -bench=BenchmarkGenerateToken -benchmem ./internal/interface/http/middleware

# Benchmark token validation
go test -bench=BenchmarkValidateToken -benchmem ./internal/interface/http/middleware
```

---

## Monitoring & Troubleshooting

### Metrics to Track

1. **CSRF Validation Failures**: Count of 403 responses with CSRF error
2. **Token Generation Rate**: Tokens generated per minute
3. **Token Store Size**: Number of active tokens in memory
4. **Token Expiry**: Tokens expired per cleanup cycle

### Common Issues

#### Issue: Client receives 403 "Invalid CSRF token"

**Causes:**
1. Client not including `X-CSRF-Token` header
2. Token expired (4-hour lifetime)
3. Token revoked (user logged out)
4. Token for wrong user (session hijacking attempt?)

**Resolution:**
- Check client code includes header for POST/PUT/DELETE/PATCH
- Implement token refresh flow (make GET request to fetch new token)
- Verify token is updated after each response (token rotation)

#### Issue: CSRF token not in response headers

**Causes:**
1. User not authenticated (CSRF only applies to authenticated requests)
2. Request to skipped path (e.g., `/health`, `/login`)
3. Safe method (GET) but not authenticated

**Resolution:**
- Ensure auth middleware runs before CSRF middleware
- Check `SkipPaths` configuration
- Verify `userID` is in request context

#### Issue: Mobile app performance degraded

**Cause:** Fetching CSRF token adds extra request

**Resolution:**
- Cache token in memory and reuse until 403 received
- Implement automatic token refresh on 403
- Consider longer token lifetime (8 hours) if security allows

---

## Compliance & Standards

### OWASP Compliance

This implementation follows OWASP recommendations:
- ✅ **Synchronizer Token Pattern**: Primary defense mechanism
- ✅ **SameSite Cookie Attribute**: Defense in depth
- ✅ **Custom Header**: Additional validation layer
- ✅ **Token Rotation**: Reduces exposure window
- ✅ **HTTPS Only**: Prevents token interception

**Reference:** [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

### Regulatory Requirements

- **CBN Requirements**: CSRF protection is required for financial operations
- **PCI DSS**: CSRF protection contributes to Requirement 6.5.9 (improper authentication)
- **NDPR**: CSRF protection helps prevent unauthorized data modification

---

## Roadmap & Future Enhancements

### Phase 2 (Post-Launch)

1. **Redis-Backed Token Store**: Replace in-memory store with Redis for multi-instance deployments
2. **Per-Session Token Limits**: Limit active CSRF tokens per user (detect session hijacking)
3. **Anomaly Detection**: Alert on unusual CSRF failure patterns
4. **Token Binding**: Bind CSRF token to specific device/IP (optional, may break legitimate use cases)

### Phase 3 (Advanced)

1. **Origin Header Validation**: Additional check for Origin/Referer headers
2. **Double Submit Cookie Pattern**: Alternative to token store for stateless deployments
3. **Encrypted Tokens**: Encrypt token payload (contains userID) for defense in depth

---

## References

1. **OWASP CSRF Prevention Cheat Sheet**
   https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html

2. **RFC 6265: HTTP State Management Mechanism (Cookies)**
   https://datatracker.ietf.org/doc/html/rfc6265

3. **SameSite Cookie Attribute**
   https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite

4. **CSRF and Mobile Applications**
   https://owasp.org/www-community/attacks/csrf

---

## Changelog

| Version | Date       | Author           | Changes                             |
|---------|------------|------------------|-------------------------------------|
| 1.0     | 2026-02-06 | Claude Sonnet 4.5 | Initial implementation and documentation |

---

**Document Owner:** Security Team
**Next Review:** 2026-03-06
**Classification:** Internal - Security Implementation

---

*For questions or clarification, contact security@hustlex.ng*
