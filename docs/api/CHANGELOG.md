# HustleX API Changelog

All notable changes to the HustleX API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this API adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- WebSocket support for real-time notifications
- Batch transaction endpoints
- Enhanced analytics endpoints

---

## [1.0.0] - 2024-01-15

### Added

#### Authentication
- `POST /auth/register` - User registration with phone number
- `POST /auth/login` - Login with phone and password
- `POST /auth/verify-otp` - OTP verification for 2FA
- `POST /auth/refresh` - Token refresh endpoint
- `POST /auth/logout` - Session invalidation
- `POST /auth/resend-otp` - Resend OTP with rate limiting
- `POST /auth/forgot-password` - Password reset initiation
- `POST /auth/reset-password` - Password reset completion
- `POST /auth/pin` - Set transaction PIN
- `PUT /auth/pin` - Change transaction PIN

#### Wallet
- `GET /wallet/balance` - Get wallet balance (available, pending, locked)
- `GET /wallet/transactions` - Transaction history with cursor pagination
- `GET /wallet/transactions/{id}` - Transaction details
- `POST /wallet/deposits` - Initiate deposit via Paystack
- `GET /wallet/deposits/{reference}/verify` - Verify deposit status
- `POST /wallet/withdrawals` - Initiate bank withdrawal
- `POST /wallet/transfers` - P2P transfer to HustleX users
- `GET /wallet/bank-accounts` - List linked bank accounts
- `POST /wallet/bank-accounts` - Add bank account with NUBAN verification
- `DELETE /wallet/bank-accounts/{id}` - Remove bank account
- `GET /wallet/banks` - List supported Nigerian banks
- `POST /wallet/banks/verify-account` - Verify bank account number

#### Gigs
- `GET /gigs` - Browse gigs with filters (category, budget, status)
- `GET /gigs/my-gigs` - User's gigs (as client or freelancer)
- `GET /gigs/{id}` - Gig details
- `POST /gigs` - Create new gig
- `PATCH /gigs/{id}` - Update gig
- `POST /gigs/{id}/cancel` - Cancel gig
- `GET /gigs/{id}/proposals` - Get proposals for a gig
- `POST /gigs/{id}/proposals` - Submit proposal
- `POST /gigs/{id}/proposals/{proposalId}/accept` - Accept proposal
- `GET /gigs/{id}/contract` - Get contract details
- `POST /gigs/{id}/milestones/{milestoneId}/complete` - Complete milestone
- `POST /gigs/{id}/milestones/{milestoneId}/approve` - Approve milestone
- `POST /gigs/{id}/reviews` - Submit review
- `GET /gigs/categories` - List categories and subcategories

#### Savings (Ajo/Esusu)
- `GET /savings/circles` - User's savings circles
- `GET /savings/circles/{id}` - Circle details with members
- `POST /savings/circles` - Create new circle
- `POST /savings/circles/{id}/invites` - Invite member
- `POST /savings/invites/{id}/accept` - Accept invitation
- `POST /savings/circles/{id}/contributions` - Make contribution
- `GET /savings/circles/{id}/contributions` - Contribution history
- `GET /savings/circles/{id}/payouts` - Payout schedule
- `GET /savings/stats` - User's savings statistics

#### Credit
- `GET /credit/score` - Credit score with factor breakdown
- `GET /credit/eligibility` - Loan eligibility check
- `GET /credit/offers` - Available loan products
- `POST /credit/loans` - Apply for loan
- `GET /credit/loans` - User's loan history
- `GET /credit/loans/{id}` - Loan details
- `POST /credit/loans/{id}/repayments` - Make repayment
- `GET /credit/loans/{id}/schedule` - Repayment schedule
- `GET /credit/stats` - Loan statistics

#### Profile & KYC
- `GET /profile` - Get user profile
- `PATCH /profile` - Update profile
- `POST /profile/avatar` - Upload avatar
- `GET /profile/kyc` - KYC status
- `POST /profile/kyc/bvn` - Submit BVN (Tier 1)
- `POST /profile/kyc/nin` - Submit NIN
- `POST /profile/kyc/document` - Upload ID document (Tier 2)

#### Notifications
- `GET /notifications` - Get notifications
- `POST /notifications/{id}/read` - Mark as read
- `POST /notifications/read-all` - Mark all as read
- `POST /notifications/devices` - Register push device
- `PUT /notifications/preferences` - Update preferences

#### Webhooks
- `POST /webhooks/paystack` - Paystack payment webhook

#### Health
- `GET /health` - Basic health check
- `GET /health/detailed` - Detailed component health

### Security
- JWT Bearer token authentication
- Transaction PIN for financial operations
- HMAC SHA512 webhook signature verification
- Rate limiting on all endpoints
- Idempotency key support for POST/PUT requests

### Standards
- OpenAPI 3.1 specification
- RFC 7807 Problem Details for errors
- Cursor-based pagination
- ISO 8601 timestamps
- NGN currency support

---

## Migration Guide

### Upgrading to v1.0.0

This is the initial release. No migration required.

---

## Deprecation Policy

- Deprecated endpoints will be marked with `X-Deprecated` header
- Deprecated endpoints remain functional for minimum 6 months
- Migration guides provided before deprecation
- Breaking changes only in major versions

---

## Versioning

The API uses URL versioning (`/v1/`, `/v2/`). Version is required in all requests.

- **Major version** (v1 → v2): Breaking changes
- **Minor updates**: New endpoints, new optional fields
- **Patch updates**: Bug fixes, documentation

---

## Rate Limits by Endpoint Category

| Category | Limit | Window |
|----------|-------|--------|
| Authentication | 10 requests | 1 minute |
| Transactions | 30 requests | 1 minute |
| Standard | 100 requests | 1 minute |
| Webhooks | 1000 requests | 1 minute |

---

## Support

- API Status: https://status.hustlex.ng
- Documentation: https://docs.hustlex.ng
- Support: api-support@hustlex.ng
