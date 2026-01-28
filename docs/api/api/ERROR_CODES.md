# HustleX API Error Codes Reference

This document provides a comprehensive reference for all error codes returned by the HustleX API.

## Error Response Format

All errors follow the [RFC 7807](https://tools.ietf.org/html/rfc7807) Problem Details specification:

```json
{
  "type": "https://api.hustlex.ng/errors/{error-type}",
  "title": "Human-readable error title",
  "status": 400,
  "detail": "Detailed explanation of the error",
  "instance": "/v1/endpoint/that/failed",
  "errorCode": "SPECIFIC_ERROR_CODE",
  "errors": [
    {
      "field": "fieldName",
      "code": "FIELD_ERROR_CODE",
      "message": "Field-specific error message"
    }
  ]
}
```

---

## HTTP Status Codes

| Status | Meaning | When Used |
|--------|---------|-----------|
| 200 | OK | Successful GET, PUT, PATCH requests |
| 201 | Created | Successful POST creating a new resource |
| 204 | No Content | Successful DELETE requests |
| 400 | Bad Request | Invalid request syntax or parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Valid auth but insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource state conflict |
| 422 | Unprocessable Entity | Validation errors |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server-side error |
| 503 | Service Unavailable | Maintenance or dependency failure |

---

## Error Categories

### Authentication Errors (AUTH_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `AUTH_INVALID_CREDENTIALS` | 401 | Phone number or password is incorrect | Verify credentials and retry |
| `AUTH_TOKEN_EXPIRED` | 401 | Access token has expired | Use refresh token to get new access token |
| `AUTH_TOKEN_INVALID` | 401 | Token is malformed or tampered | Re-authenticate with login |
| `AUTH_REFRESH_TOKEN_EXPIRED` | 401 | Refresh token has expired | Re-authenticate with login |
| `AUTH_SESSION_REVOKED` | 401 | Session was explicitly revoked | Re-authenticate with login |
| `AUTH_OTP_INVALID` | 400 | OTP code is incorrect | Request new OTP and retry |
| `AUTH_OTP_EXPIRED` | 400 | OTP has expired (5 min validity) | Request new OTP |
| `AUTH_OTP_MAX_ATTEMPTS` | 429 | Too many incorrect OTP attempts | Wait 15 minutes before retrying |
| `AUTH_PHONE_NOT_VERIFIED` | 403 | Phone number not verified | Complete OTP verification |
| `AUTH_ACCOUNT_LOCKED` | 403 | Account temporarily locked | Contact support or wait for unlock |
| `AUTH_ACCOUNT_SUSPENDED` | 403 | Account suspended for policy violation | Contact support |
| `AUTH_PIN_NOT_SET` | 403 | Transaction PIN not configured | Set PIN before transactions |
| `AUTH_PIN_INVALID` | 400 | Transaction PIN is incorrect | Verify PIN and retry |
| `AUTH_PIN_LOCKED` | 403 | PIN locked after multiple failures | Reset PIN or wait 30 minutes |

### Validation Errors (VALIDATION_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `VALIDATION_REQUIRED` | 422 | Required field is missing | Include the required field |
| `VALIDATION_INVALID_FORMAT` | 422 | Field format is invalid | Follow the expected format |
| `VALIDATION_MIN_LENGTH` | 422 | Value below minimum length | Provide longer value |
| `VALIDATION_MAX_LENGTH` | 422 | Value exceeds maximum length | Shorten the value |
| `VALIDATION_MIN_VALUE` | 422 | Numeric value below minimum | Increase the value |
| `VALIDATION_MAX_VALUE` | 422 | Numeric value exceeds maximum | Decrease the value |
| `VALIDATION_INVALID_PHONE` | 422 | Phone number format invalid | Use Nigerian format: +234XXXXXXXXXX |
| `VALIDATION_INVALID_EMAIL` | 422 | Email format invalid | Use valid email format |
| `VALIDATION_INVALID_BVN` | 422 | BVN format invalid | BVN must be 11 digits |
| `VALIDATION_INVALID_NIN` | 422 | NIN format invalid | NIN must be 11 digits |
| `VALIDATION_INVALID_NUBAN` | 422 | Account number invalid | Account must be 10 digits |
| `VALIDATION_WEAK_PASSWORD` | 422 | Password doesn't meet requirements | Min 8 chars, uppercase, lowercase, number |
| `VALIDATION_INVALID_AMOUNT` | 422 | Amount format or value invalid | Positive number, max 2 decimals |
| `VALIDATION_INVALID_PIN` | 422 | PIN format invalid | PIN must be 4 digits |

### Wallet Errors (WALLET_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `WALLET_INSUFFICIENT_BALANCE` | 400 | Not enough available balance | Fund wallet or reduce amount |
| `WALLET_BELOW_MINIMUM` | 400 | Amount below minimum transaction | Increase to minimum amount |
| `WALLET_ABOVE_MAXIMUM` | 400 | Amount exceeds transaction limit | Reduce amount or upgrade KYC |
| `WALLET_DAILY_LIMIT_EXCEEDED` | 400 | Daily transaction limit reached | Wait until tomorrow or upgrade |
| `WALLET_WITHDRAWAL_PENDING` | 409 | Another withdrawal in progress | Wait for pending to complete |
| `WALLET_TRANSFER_SELF` | 400 | Cannot transfer to self | Use different recipient |
| `WALLET_RECIPIENT_NOT_FOUND` | 404 | Transfer recipient doesn't exist | Verify recipient details |
| `WALLET_LOCKED` | 403 | Wallet is locked | Contact support |

### Bank Account Errors (BANK_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `BANK_ACCOUNT_NOT_FOUND` | 404 | Bank account doesn't exist | Add bank account first |
| `BANK_ACCOUNT_VERIFICATION_FAILED` | 400 | NUBAN verification failed | Verify account number and bank |
| `BANK_NAME_MISMATCH` | 400 | Account name doesn't match user | Use account in your name |
| `BANK_NOT_SUPPORTED` | 400 | Bank not supported for transfers | Use supported bank |
| `BANK_TRANSFER_FAILED` | 500 | Transfer to bank failed | Retry or contact support |
| `BANK_ACCOUNT_LIMIT_REACHED` | 400 | Maximum bank accounts reached | Remove an account first |

### Payment Errors (PAYMENT_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `PAYMENT_FAILED` | 400 | Payment processing failed | Retry or use different method |
| `PAYMENT_CANCELLED` | 400 | Payment was cancelled | Initiate new payment |
| `PAYMENT_EXPIRED` | 400 | Payment session expired | Initiate new payment |
| `PAYMENT_DUPLICATE` | 409 | Duplicate payment detected | Check transaction history |
| `PAYMENT_CARD_DECLINED` | 400 | Card was declined | Use different card |
| `PAYMENT_INVALID_REFERENCE` | 404 | Payment reference not found | Verify reference |

### Gig Errors (GIG_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `GIG_NOT_FOUND` | 404 | Gig doesn't exist | Verify gig ID |
| `GIG_NOT_OPEN` | 400 | Gig is not accepting proposals | Find another gig |
| `GIG_ALREADY_APPLIED` | 409 | Already submitted proposal | Edit existing proposal |
| `GIG_NOT_OWNER` | 403 | Not the gig owner | Only owner can modify |
| `GIG_CANNOT_CANCEL` | 400 | Gig cannot be cancelled | Gig has active contract |
| `GIG_BUDGET_INVALID` | 422 | Budget range invalid | Min must be less than max |

### Proposal Errors (PROPOSAL_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `PROPOSAL_NOT_FOUND` | 404 | Proposal doesn't exist | Verify proposal ID |
| `PROPOSAL_NOT_OWNER` | 403 | Not the proposal owner | Only owner can modify |
| `PROPOSAL_ALREADY_ACCEPTED` | 409 | Proposal already accepted | Cannot modify accepted proposal |
| `PROPOSAL_BID_INVALID` | 422 | Bid outside allowed range | Bid within gig budget range |
| `PROPOSAL_CANNOT_WITHDRAW` | 400 | Cannot withdraw proposal | Proposal already processed |

### Contract Errors (CONTRACT_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `CONTRACT_NOT_FOUND` | 404 | Contract doesn't exist | Verify contract ID |
| `CONTRACT_NOT_PARTY` | 403 | Not a contract party | Only client/freelancer can access |
| `CONTRACT_NOT_ACTIVE` | 400 | Contract not in active state | Contract may be completed/cancelled |
| `CONTRACT_ALREADY_EXISTS` | 409 | Gig already has a contract | Cannot create another contract |

### Milestone Errors (MILESTONE_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `MILESTONE_NOT_FOUND` | 404 | Milestone doesn't exist | Verify milestone ID |
| `MILESTONE_NOT_PENDING` | 400 | Milestone not in pending state | Cannot complete/approve |
| `MILESTONE_NOT_SUBMITTED` | 400 | Milestone not submitted for review | Freelancer must complete first |
| `MILESTONE_ALREADY_PAID` | 409 | Milestone already paid | Cannot approve twice |

### Savings Errors (SAVINGS_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `SAVINGS_CIRCLE_NOT_FOUND` | 404 | Savings circle doesn't exist | Verify circle ID |
| `SAVINGS_NOT_MEMBER` | 403 | Not a circle member | Join circle first |
| `SAVINGS_NOT_ADMIN` | 403 | Not circle admin | Only admin can perform action |
| `SAVINGS_CIRCLE_FULL` | 400 | Circle has maximum members | Wait or join different circle |
| `SAVINGS_ALREADY_MEMBER` | 409 | Already a circle member | Cannot join twice |
| `SAVINGS_CONTRIBUTION_NOT_DUE` | 400 | Contribution not due yet | Wait for contribution window |
| `SAVINGS_ALREADY_CONTRIBUTED` | 409 | Already contributed this round | Wait for next round |
| `SAVINGS_CIRCLE_STARTED` | 400 | Circle already started | Cannot modify started circle |
| `SAVINGS_INVITE_NOT_FOUND` | 404 | Invitation doesn't exist | Request new invite |
| `SAVINGS_INVITE_EXPIRED` | 400 | Invitation has expired | Request new invite |
| `SAVINGS_INVITE_ALREADY_USED` | 409 | Invitation already used | Request new invite |

### Credit Errors (CREDIT_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `CREDIT_NOT_ELIGIBLE` | 400 | Not eligible for loans | Improve credit score |
| `CREDIT_SCORE_UNAVAILABLE` | 400 | Credit score not calculated | Complete more transactions |
| `CREDIT_ACTIVE_LOAN_EXISTS` | 409 | Already have active loan | Repay existing loan first |
| `CREDIT_OFFER_NOT_FOUND` | 404 | Loan offer doesn't exist | Check available offers |
| `CREDIT_AMOUNT_INVALID` | 422 | Amount outside offer limits | Choose amount within limits |
| `CREDIT_TENOR_INVALID` | 422 | Tenor outside allowed range | Choose valid tenor |

### Loan Errors (LOAN_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `LOAN_NOT_FOUND` | 404 | Loan doesn't exist | Verify loan ID |
| `LOAN_NOT_ACTIVE` | 400 | Loan not in active state | Cannot make repayment |
| `LOAN_OVERPAYMENT` | 400 | Payment exceeds balance | Pay exact or less amount |
| `LOAN_REPAYMENT_FAILED` | 500 | Repayment processing failed | Retry or contact support |

### KYC Errors (KYC_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `KYC_BVN_VERIFICATION_FAILED` | 400 | BVN verification failed | Verify BVN details |
| `KYC_BVN_MISMATCH` | 400 | BVN doesn't match user details | Use your own BVN |
| `KYC_NIN_VERIFICATION_FAILED` | 400 | NIN verification failed | Verify NIN details |
| `KYC_DOCUMENT_INVALID` | 400 | ID document invalid or unreadable | Upload clear document |
| `KYC_DOCUMENT_EXPIRED` | 400 | ID document has expired | Upload valid document |
| `KYC_ALREADY_VERIFIED` | 409 | KYC level already verified | No action needed |
| `KYC_PENDING_REVIEW` | 400 | Previous submission under review | Wait for review completion |

### Rate Limit Errors (RATE_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests | Wait and retry (see Retry-After header) |
| `RATE_LIMIT_AUTH` | 429 | Auth rate limit exceeded | Wait 60 seconds |
| `RATE_LIMIT_TRANSACTION` | 429 | Transaction rate limit exceeded | Wait before next transaction |

### System Errors (SYSTEM_*)

| Code | HTTP Status | Description | Resolution |
|------|-------------|-------------|------------|
| `SYSTEM_ERROR` | 500 | Internal server error | Retry or contact support |
| `SYSTEM_MAINTENANCE` | 503 | System under maintenance | Check status page |
| `SYSTEM_DEPENDENCY_ERROR` | 503 | External service unavailable | Retry later |
| `SYSTEM_TIMEOUT` | 504 | Request timed out | Retry the request |

---

## Field-Specific Validation Codes

Used in the `errors[].code` field for 422 responses:

| Code | Description |
|------|-------------|
| `REQUIRED` | Field is required but missing |
| `INVALID_FORMAT` | Field format doesn't match expected pattern |
| `INVALID_TYPE` | Field type is incorrect |
| `MIN_LENGTH` | String shorter than minimum |
| `MAX_LENGTH` | String longer than maximum |
| `MIN_VALUE` | Number below minimum |
| `MAX_VALUE` | Number above maximum |
| `NOT_UNIQUE` | Value already exists |
| `INVALID_REFERENCE` | Referenced resource doesn't exist |
| `INVALID_ENUM` | Value not in allowed set |

---

## Error Handling Best Practices

### Client Implementation

```javascript
async function apiRequest(endpoint, options) {
  try {
    const response = await fetch(`${BASE_URL}${endpoint}`, options);

    if (!response.ok) {
      const error = await response.json();

      switch (error.errorCode) {
        case 'AUTH_TOKEN_EXPIRED':
          // Attempt token refresh
          await refreshToken();
          return apiRequest(endpoint, options);

        case 'RATE_LIMIT_EXCEEDED':
          // Wait and retry
          const retryAfter = response.headers.get('Retry-After');
          await sleep(parseInt(retryAfter) * 1000);
          return apiRequest(endpoint, options);

        case 'VALIDATION_REQUIRED':
        case 'VALIDATION_INVALID_FORMAT':
          // Show field-specific errors to user
          displayValidationErrors(error.errors);
          break;

        default:
          // Show generic error message
          displayError(error.detail);
      }

      throw new ApiError(error);
    }

    return response.json();
  } catch (e) {
    // Handle network errors
    throw new NetworkError(e.message);
  }
}
```

### Flutter Implementation

```dart
Future<T> apiRequest<T>(String endpoint, {required T Function(Map<String, dynamic>) parser}) async {
  final response = await http.get(Uri.parse('$baseUrl$endpoint'));

  if (response.statusCode >= 400) {
    final error = jsonDecode(response.body);

    switch (error['errorCode']) {
      case 'AUTH_TOKEN_EXPIRED':
        await refreshToken();
        return apiRequest(endpoint, parser: parser);

      case 'WALLET_INSUFFICIENT_BALANCE':
        throw InsufficientBalanceException(error['detail']);

      default:
        throw ApiException(error['errorCode'], error['detail']);
    }
  }

  return parser(jsonDecode(response.body));
}
```

---

## Support

If you encounter an error not documented here or need assistance:

- **Email**: api-support@hustlex.ng
- **Status Page**: https://status.hustlex.ng
- **Documentation**: https://docs.hustlex.ng
