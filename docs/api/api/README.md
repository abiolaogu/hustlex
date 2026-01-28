# HustleX API Documentation

Complete API documentation for the HustleX platform - Nigerian gig economy financial services.

## Overview

HustleX provides APIs for:
- **Authentication** - User registration, login, 2FA, PIN management
- **Wallet** - Balance, transactions, deposits, withdrawals, transfers
- **Gigs** - Marketplace for freelance opportunities
- **Savings** - Ajo/Esusu rotating savings circles
- **Credit** - Credit scoring and loan management
- **Profile** - User profile and KYC verification

## Documentation Files

| File | Description |
|------|-------------|
| [`openapi.yaml`](./openapi.yaml) | OpenAPI 3.1 specification |
| [`CHANGELOG.md`](./CHANGELOG.md) | API version history |
| [`ERROR_CODES.md`](./ERROR_CODES.md) | Error codes reference |
| [`WEBHOOKS.md`](./WEBHOOKS.md) | Webhook integration guide |

## Postman Collection

Import these files into Postman for API testing:

- [`postman/HustleX.postman_collection.json`](./postman/HustleX.postman_collection.json) - API collection
- [`postman/HustleX.postman_environment.json`](./postman/HustleX.postman_environment.json) - Environment variables

## Quick Start

### Base URL

```
Production: https://api.hustlex.ng/v1
Staging:    https://staging-api.hustlex.ng/v1
Local:      http://localhost:8080/v1
```

### Authentication

Most endpoints require Bearer token authentication:

```bash
curl -X GET https://api.hustlex.ng/v1/wallet/balance \
  -H "Authorization: Bearer <access_token>"
```

### Rate Limits

| Category | Limit | Window |
|----------|-------|--------|
| Authentication | 10 requests | 1 minute |
| Transactions | 30 requests | 1 minute |
| Standard | 100 requests | 1 minute |

Rate limit headers are included in all responses:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining
- `X-RateLimit-Reset`: Unix timestamp when limit resets

## Generate Documentation

Use the generation script to validate spec and generate artifacts:

```bash
# Validate OpenAPI spec
../backend/scripts/generate-openapi.sh validate

# Generate all (bundle, SDKs, HTML docs)
../backend/scripts/generate-openapi.sh all
```

## API Standards

- **Versioning**: URL-based (`/v1/`)
- **Pagination**: Cursor-based
- **Errors**: RFC 7807 Problem Details
- **Timestamps**: ISO 8601 (UTC)
- **Currency**: NGN (Nigerian Naira)
- **Authentication**: JWT Bearer tokens
- **Idempotency**: `X-Idempotency-Key` header for POST/PUT

## Support

- **Status Page**: https://status.hustlex.ng
- **API Support**: api-support@hustlex.ng
- **Documentation**: https://docs.hustlex.ng
