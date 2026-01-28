# HustleX API Reference

> Version 1.0 | Base URL: `https://api.hustlex.app/api/v1`

## Table of Contents

1. [Authentication](#authentication)
2. [Wallet](#wallet)
3. [Gigs](#gigs)
4. [Savings](#savings)
5. [Credit](#credit)
6. [Profile](#profile)
7. [Notifications](#notifications)
8. [Webhooks](#webhooks)

---

## Overview

### Base URL

```
Production: https://api.hustlex.app/api/v1
Staging:    https://staging-api.hustlex.app/api/v1
```

### Authentication

All endpoints (except auth) require a Bearer token:

```
Authorization: Bearer <access_token>
```

### Response Format

```json
{
    "success": true,
    "data": { },
    "meta": {
        "request_id": "550e8400-e29b-41d4-a716-446655440000",
        "timestamp": "2024-01-15T10:30:00Z"
    }
}
```

### Error Format

```json
{
    "success": false,
    "error": {
        "code": "ERROR_CODE",
        "message": "Human readable message",
        "details": []
    }
}
```

---

## 1. Authentication

### Request OTP

Send a one-time password to a phone number.

```http
POST /auth/otp/request
```

**Request Body:**

```json
{
    "phone": "+2348012345678"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "expires_in": 300,
        "resend_cooldown": 60
    }
}
```

**Rate Limit:** 5 requests per 15 minutes

---

### Verify OTP

Verify the OTP and receive authentication tokens.

```http
POST /auth/otp/verify
```

**Request Body:**

```json
{
    "phone": "+2348012345678",
    "code": "123456"
}
```

**Response (Existing User):**

```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 900,
        "user": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "phone": "+2348012345678",
            "full_name": "John Doe",
            "tier": "silver",
            "has_pin": true
        },
        "is_new_user": false
    }
}
```

**Response (New User):**

```json
{
    "success": true,
    "data": {
        "access_token": "...",
        "refresh_token": "...",
        "is_new_user": true,
        "registration_token": "temp_token_for_registration"
    }
}
```

---

### Register User

Complete registration for new users.

```http
POST /auth/register
```

**Headers:**

```
Authorization: Bearer <registration_token>
```

**Request Body:**

```json
{
    "full_name": "John Doe",
    "email": "john@example.com",
    "date_of_birth": "1990-05-15",
    "referral_code": "ABC123"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "user": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "phone": "+2348012345678",
            "full_name": "John Doe",
            "email": "john@example.com",
            "tier": "bronze",
            "referral_code": "JOH123"
        },
        "wallet": {
            "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
            "balance": 0,
            "currency": "NGN"
        }
    }
}
```

---

### Set Transaction PIN

Set or update the transaction PIN.

```http
POST /auth/pin/set
```

**Request Body:**

```json
{
    "pin": "123456",
    "confirm_pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "message": "PIN set successfully"
    }
}
```

---

### Refresh Token

Get a new access token using refresh token.

```http
POST /auth/refresh
```

**Request Body:**

```json
{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
        "expires_in": 900
    }
}
```

---

### Logout

Invalidate current tokens.

```http
POST /auth/logout
```

**Response:**

```json
{
    "success": true,
    "data": {
        "message": "Logged out successfully"
    }
}
```

---

## 2. Wallet

### Get Wallet

Get current user's wallet details.

```http
GET /wallet
```

**Response:**

```json
{
    "success": true,
    "data": {
        "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
        "balance": 50000.00,
        "escrow_balance": 10000.00,
        "savings_balance": 25000.00,
        "currency": "NGN",
        "daily_limit": 500000.00,
        "monthly_limit": 5000000.00,
        "daily_spent": 15000.00,
        "monthly_spent": 150000.00
    }
}
```

---

### Get Transactions

Get transaction history with pagination.

```http
GET /wallet/transactions?cursor=2024-01-15T10:00:00Z&limit=20&type=credit
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| cursor | datetime | Pagination cursor (created_at) |
| limit | integer | Items per page (max 50) |
| type | string | Filter by type (credit, debit) |
| category | string | Filter by category |

**Response:**

```json
{
    "success": true,
    "data": {
        "transactions": [
            {
                "id": "123e4567-e89b-12d3-a456-426614174000",
                "type": "credit",
                "category": "deposit",
                "amount": 10000.00,
                "fee": 0,
                "reference": "TXN_20240115_001",
                "description": "Wallet deposit via Paystack",
                "status": "completed",
                "created_at": "2024-01-15T10:30:00Z"
            }
        ],
        "pagination": {
            "has_more": true,
            "next_cursor": "2024-01-14T15:20:00Z"
        }
    }
}
```

---

### Initiate Deposit

Start a wallet deposit via Paystack.

```http
POST /wallet/deposit
```

**Request Body:**

```json
{
    "amount": 10000.00,
    "channel": "card"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "authorization_url": "https://checkout.paystack.com/xyz123",
        "reference": "DEP_20240115_001",
        "access_code": "xyz123"
    }
}
```

---

### Transfer to User

Transfer funds to another HustleX user.

```http
POST /wallet/transfer
```

**Request Body:**

```json
{
    "recipient_phone": "+2348098765432",
    "amount": 5000.00,
    "note": "Payment for services",
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "transaction": {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "type": "debit",
            "category": "transfer",
            "amount": 5000.00,
            "fee": 0,
            "reference": "TRF_20240115_001",
            "status": "completed"
        },
        "recipient": {
            "name": "Jane Doe",
            "phone": "+2348098765432"
        },
        "new_balance": 45000.00
    }
}
```

---

### Withdraw to Bank

Withdraw funds to a bank account.

```http
POST /wallet/withdraw
```

**Request Body:**

```json
{
    "amount": 20000.00,
    "bank_code": "058",
    "account_number": "0123456789",
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "transaction": {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "type": "debit",
            "category": "withdrawal",
            "amount": 20000.00,
            "fee": 25.00,
            "status": "pending"
        },
        "estimated_arrival": "2024-01-15T12:00:00Z"
    }
}
```

---

### Get Banks

Get list of supported banks.

```http
GET /wallet/banks
```

**Response:**

```json
{
    "success": true,
    "data": {
        "banks": [
            {
                "code": "058",
                "name": "GTBank",
                "logo": "https://cdn.hustlex.app/banks/gtbank.png"
            },
            {
                "code": "044",
                "name": "Access Bank",
                "logo": "https://cdn.hustlex.app/banks/access.png"
            }
        ]
    }
}
```

---

### Verify Bank Account

Verify bank account details before withdrawal.

```http
POST /wallet/verify-account
```

**Request Body:**

```json
{
    "bank_code": "058",
    "account_number": "0123456789"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "account_name": "JOHN DOE",
        "account_number": "0123456789",
        "bank_name": "GTBank"
    }
}
```

---

## 3. Gigs

### List Gigs

Get available gigs with filters.

```http
GET /gigs?category=technology&status=open&min_budget=5000&max_budget=50000
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| category | string | Filter by category |
| status | string | Filter by status (open, in_progress, completed) |
| min_budget | number | Minimum budget |
| max_budget | number | Maximum budget |
| search | string | Search in title/description |
| cursor | string | Pagination cursor |
| limit | integer | Items per page |

**Response:**

```json
{
    "success": true,
    "data": {
        "gigs": [
            {
                "id": "123e4567-e89b-12d3-a456-426614174000",
                "title": "Build Mobile App UI",
                "description": "Need a skilled Flutter developer...",
                "category": "technology",
                "budget_min": 50000.00,
                "budget_max": 100000.00,
                "deadline": "2024-02-15",
                "status": "open",
                "proposals_count": 5,
                "client": {
                    "id": "abc123",
                    "name": "TechCorp Ltd",
                    "rating": 4.8,
                    "gigs_posted": 15
                },
                "skills": ["flutter", "dart", "mobile"],
                "created_at": "2024-01-10T08:00:00Z"
            }
        ],
        "pagination": {
            "has_more": true,
            "next_cursor": "abc123"
        }
    }
}
```

---

### Create Gig

Post a new gig (as client).

```http
POST /gigs
```

**Request Body:**

```json
{
    "title": "Build Mobile App UI",
    "description": "Need a skilled Flutter developer to build UI screens for our e-commerce app. Must have experience with complex animations and responsive design.",
    "category": "technology",
    "budget_min": 50000.00,
    "budget_max": 100000.00,
    "deadline": "2024-02-15",
    "skills": ["flutter", "dart", "mobile"],
    "attachments": ["https://cdn.hustlex.app/files/spec.pdf"]
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "gig": {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "title": "Build Mobile App UI",
            "status": "open",
            "created_at": "2024-01-15T10:30:00Z"
        }
    }
}
```

---

### Get Gig Details

Get full details of a specific gig.

```http
GET /gigs/{id}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "gig": {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "title": "Build Mobile App UI",
            "description": "...",
            "category": "technology",
            "budget_min": 50000.00,
            "budget_max": 100000.00,
            "deadline": "2024-02-15",
            "status": "open",
            "skills": ["flutter", "dart", "mobile"],
            "attachments": ["https://cdn.hustlex.app/files/spec.pdf"],
            "client": {
                "id": "abc123",
                "name": "TechCorp Ltd",
                "profile_photo": "https://cdn.hustlex.app/users/abc123.jpg",
                "rating": 4.8,
                "gigs_posted": 15,
                "member_since": "2023-06-01"
            },
            "proposals_count": 5,
            "views_count": 45,
            "created_at": "2024-01-10T08:00:00Z"
        },
        "user_proposal": null
    }
}
```

---

### Submit Proposal

Submit a proposal for a gig (as freelancer).

```http
POST /gigs/{id}/proposals
```

**Request Body:**

```json
{
    "cover_letter": "I am an experienced Flutter developer with 5+ years of experience. I have worked on similar e-commerce apps and can deliver high-quality UI with smooth animations.",
    "proposed_price": 75000.00,
    "delivery_days": 14,
    "attachments": ["https://cdn.hustlex.app/files/portfolio.pdf"]
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "proposal": {
            "id": "456e7890-f12c-34d5-b678-901234567890",
            "status": "pending",
            "proposed_price": 75000.00,
            "delivery_days": 14,
            "created_at": "2024-01-15T11:00:00Z"
        }
    }
}
```

---

### Get Gig Proposals

Get all proposals for a gig (client only).

```http
GET /gigs/{id}/proposals
```

**Response:**

```json
{
    "success": true,
    "data": {
        "proposals": [
            {
                "id": "456e7890-f12c-34d5-b678-901234567890",
                "freelancer": {
                    "id": "xyz789",
                    "name": "Jane Smith",
                    "profile_photo": "...",
                    "rating": 4.9,
                    "completed_gigs": 25,
                    "skills": ["flutter", "dart"]
                },
                "cover_letter": "...",
                "proposed_price": 75000.00,
                "delivery_days": 14,
                "status": "pending",
                "created_at": "2024-01-15T11:00:00Z"
            }
        ]
    }
}
```

---

### Accept Proposal

Accept a proposal and create contract (client only).

```http
POST /gigs/{gig_id}/proposals/{proposal_id}/accept
```

**Request Body:**

```json
{
    "agreed_price": 70000.00,
    "delivery_days": 14,
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "contract": {
            "id": "789e0123-a45b-67c8-d901-234567890abc",
            "gig_id": "123e4567-e89b-12d3-a456-426614174000",
            "agreed_price": 70000.00,
            "platform_fee": 7000.00,
            "freelancer_amount": 63000.00,
            "delivery_days": 14,
            "deadline": "2024-01-29T11:00:00Z",
            "status": "in_progress"
        },
        "escrow": {
            "amount": 70000.00,
            "status": "funded"
        }
    }
}
```

---

### Submit Deliverable

Submit work for a contract (freelancer only).

```http
POST /gigs/{gig_id}/contracts/{contract_id}/deliver
```

**Request Body:**

```json
{
    "message": "Hi, I have completed the work. Please find the deliverables attached.",
    "attachments": [
        "https://cdn.hustlex.app/files/delivery.zip",
        "https://cdn.hustlex.app/files/documentation.pdf"
    ]
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "delivery": {
            "id": "del123",
            "status": "pending_review",
            "submitted_at": "2024-01-25T14:00:00Z"
        }
    }
}
```

---

### Approve Delivery

Approve delivery and release escrow (client only).

```http
POST /gigs/{gig_id}/contracts/{contract_id}/approve
```

**Request Body:**

```json
{
    "rating": 5,
    "review": "Excellent work! Delivered on time with great attention to detail.",
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "contract": {
            "id": "789e0123-a45b-67c8-d901-234567890abc",
            "status": "completed"
        },
        "payment": {
            "amount": 63000.00,
            "status": "released"
        }
    }
}
```

---

## 4. Savings

### List Savings Circles

Get available and joined savings circles.

```http
GET /savings/circles?type=rotational&status=active
```

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| type | string | rotational, fixed_target |
| status | string | forming, active, completed |
| joined | boolean | Only user's circles |

**Response:**

```json
{
    "success": true,
    "data": {
        "circles": [
            {
                "id": "cir123",
                "name": "Tech Savers Club",
                "type": "rotational",
                "contribution_amount": 10000.00,
                "frequency": "weekly",
                "member_count": 8,
                "max_members": 10,
                "total_pool": 80000.00,
                "status": "active",
                "next_contribution_date": "2024-01-20",
                "creator": {
                    "id": "user123",
                    "name": "John Doe"
                },
                "created_at": "2024-01-01T00:00:00Z"
            }
        ]
    }
}
```

---

### Create Savings Circle

Create a new savings circle.

```http
POST /savings/circles
```

**Request Body:**

```json
{
    "name": "Tech Savers Club",
    "type": "rotational",
    "contribution_amount": 10000.00,
    "frequency": "weekly",
    "max_members": 10,
    "start_date": "2024-02-01",
    "description": "A savings circle for tech professionals"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "circle": {
            "id": "cir456",
            "name": "Tech Savers Club",
            "invite_code": "TSC2024",
            "status": "forming"
        }
    }
}
```

---

### Join Circle

Join an existing savings circle.

```http
POST /savings/circles/{id}/join
```

**Request Body:**

```json
{
    "invite_code": "TSC2024",
    "preferred_position": 5
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "membership": {
            "id": "mem789",
            "circle_id": "cir456",
            "position": 5,
            "role": "member",
            "status": "active"
        }
    }
}
```

---

### Make Contribution

Make a contribution to a savings circle.

```http
POST /savings/circles/{id}/contribute
```

**Request Body:**

```json
{
    "amount": 10000.00,
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "contribution": {
            "id": "con123",
            "amount": 10000.00,
            "status": "completed"
        },
        "circle_progress": {
            "contributions_this_round": 7,
            "total_collected": 70000.00,
            "next_payout_to": "Jane Doe",
            "next_payout_date": "2024-01-21"
        }
    }
}
```

---

### Get Circle Details

Get full details of a savings circle.

```http
GET /savings/circles/{id}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "circle": {
            "id": "cir456",
            "name": "Tech Savers Club",
            "type": "rotational",
            "contribution_amount": 10000.00,
            "frequency": "weekly",
            "status": "active",
            "current_round": 3,
            "total_rounds": 10
        },
        "members": [
            {
                "id": "mem001",
                "user": {
                    "id": "user123",
                    "name": "John Doe",
                    "profile_photo": "..."
                },
                "position": 1,
                "role": "admin",
                "received_payout": true,
                "contribution_status": "paid"
            }
        ],
        "contributions": [
            {
                "id": "con001",
                "round": 3,
                "member_name": "John Doe",
                "amount": 10000.00,
                "status": "completed",
                "paid_at": "2024-01-15T10:00:00Z"
            }
        ],
        "upcoming_payouts": [
            {
                "round": 4,
                "recipient": "Jane Doe",
                "amount": 100000.00,
                "date": "2024-01-22"
            }
        ]
    }
}
```

---

## 5. Credit

### Get Credit Score

Get user's current credit score and breakdown.

```http
GET /credit/score
```

**Response:**

```json
{
    "success": true,
    "data": {
        "credit_score": {
            "id": "cs123",
            "score": 720,
            "tier": "good",
            "max_score": 850,
            "last_updated": "2024-01-15T00:00:00Z"
        },
        "breakdown": {
            "payment_history": {
                "score": 180,
                "max": 200,
                "description": "On-time payments"
            },
            "savings_consistency": {
                "score": 150,
                "max": 200,
                "description": "Regular savings"
            },
            "gig_performance": {
                "score": 170,
                "max": 200,
                "description": "Completed gigs"
            },
            "account_age": {
                "score": 100,
                "max": 150,
                "description": "Account longevity"
            },
            "platform_activity": {
                "score": 120,
                "max": 100,
                "description": "Regular usage"
            }
        },
        "eligible_loan_amount": 100000.00,
        "improvement_tips": [
            "Complete 2 more gigs to improve your score",
            "Maintain your savings streak for 30 more days"
        ]
    }
}
```

---

### Get Credit History

Get credit score history over time.

```http
GET /credit/history?period=6m
```

**Response:**

```json
{
    "success": true,
    "data": {
        "history": [
            {
                "date": "2024-01-01",
                "score": 700
            },
            {
                "date": "2024-01-15",
                "score": 720
            }
        ],
        "trend": "improving"
    }
}
```

---

### List Loans

Get user's loan history.

```http
GET /credit/loans?status=active
```

**Response:**

```json
{
    "success": true,
    "data": {
        "loans": [
            {
                "id": "loan123",
                "amount": 50000.00,
                "interest_rate": 5.0,
                "tenure_months": 3,
                "monthly_payment": 17500.00,
                "total_repayment": 52500.00,
                "amount_paid": 17500.00,
                "amount_remaining": 35000.00,
                "status": "active",
                "next_payment_date": "2024-02-15",
                "disbursed_at": "2024-01-15T10:00:00Z"
            }
        ]
    }
}
```

---

### Apply for Loan

Apply for a new loan.

```http
POST /credit/loans/apply
```

**Request Body:**

```json
{
    "amount": 50000.00,
    "tenure_months": 3,
    "purpose": "business_expansion"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "application": {
            "id": "app123",
            "status": "approved",
            "loan": {
                "id": "loan456",
                "amount": 50000.00,
                "interest_rate": 5.0,
                "tenure_months": 3,
                "monthly_payment": 17500.00,
                "total_repayment": 52500.00
            },
            "disbursement_time": "instant"
        }
    }
}
```

---

### Repay Loan

Make a loan repayment.

```http
POST /credit/loans/{id}/repay
```

**Request Body:**

```json
{
    "amount": 17500.00,
    "pin": "123456"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "repayment": {
            "id": "rep123",
            "amount": 17500.00,
            "status": "completed"
        },
        "loan": {
            "amount_remaining": 17500.00,
            "next_payment_date": "2024-03-15",
            "status": "active"
        }
    }
}
```

---

## 6. Profile

### Get Profile

Get current user's profile.

```http
GET /profile
```

**Response:**

```json
{
    "success": true,
    "data": {
        "user": {
            "id": "user123",
            "phone": "+2348012345678",
            "email": "john@example.com",
            "full_name": "John Doe",
            "profile_photo": "https://cdn.hustlex.app/users/user123.jpg",
            "date_of_birth": "1990-05-15",
            "tier": "silver",
            "bvn_verified": true,
            "nin_verified": false,
            "referral_code": "JOH123",
            "created_at": "2023-06-01T00:00:00Z"
        },
        "stats": {
            "gigs_completed": 15,
            "gigs_posted": 5,
            "savings_circles": 2,
            "total_earned": 500000.00,
            "total_saved": 150000.00
        }
    }
}
```

---

### Update Profile

Update user profile information.

```http
PATCH /profile
```

**Request Body:**

```json
{
    "full_name": "John D. Doe",
    "email": "johndoe@example.com",
    "profile_photo": "base64_encoded_image_or_url"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "user": {
            "id": "user123",
            "full_name": "John D. Doe",
            "email": "johndoe@example.com"
        }
    }
}
```

---

### Get Bank Accounts

Get user's saved bank accounts.

```http
GET /profile/bank-accounts
```

**Response:**

```json
{
    "success": true,
    "data": {
        "accounts": [
            {
                "id": "ba123",
                "bank_name": "GTBank",
                "bank_code": "058",
                "account_number": "0123456789",
                "account_name": "JOHN DOE",
                "is_default": true
            }
        ]
    }
}
```

---

### Add Bank Account

Add a new bank account.

```http
POST /profile/bank-accounts
```

**Request Body:**

```json
{
    "bank_code": "044",
    "account_number": "0987654321"
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "account": {
            "id": "ba456",
            "bank_name": "Access Bank",
            "account_number": "0987654321",
            "account_name": "JOHN DOE"
        }
    }
}
```

---

## 7. Notifications

### Get Notifications

Get user notifications.

```http
GET /notifications?unread_only=true
```

**Response:**

```json
{
    "success": true,
    "data": {
        "notifications": [
            {
                "id": "not123",
                "type": "payment",
                "title": "Payment Received",
                "body": "You received ₦50,000 from TechCorp Ltd",
                "data": {
                    "transaction_id": "tx123",
                    "amount": 50000
                },
                "is_read": false,
                "created_at": "2024-01-15T10:30:00Z"
            }
        ],
        "unread_count": 5
    }
}
```

---

### Mark as Read

Mark notifications as read.

```http
POST /notifications/read
```

**Request Body:**

```json
{
    "notification_ids": ["not123", "not124"]
}
```

**Response:**

```json
{
    "success": true,
    "data": {
        "updated_count": 2
    }
}
```

---

### Update Push Settings

Update notification preferences.

```http
PUT /notifications/settings
```

**Request Body:**

```json
{
    "push_enabled": true,
    "email_enabled": false,
    "sms_enabled": true,
    "categories": {
        "payments": true,
        "gigs": true,
        "savings": true,
        "promotions": false
    }
}
```

---

## 8. Webhooks

### Paystack Webhook

Endpoint for Paystack payment notifications.

```http
POST /webhooks/paystack
```

**Headers:**

```
X-Paystack-Signature: sha512_hash
```

**Event Types:**

| Event | Description |
|-------|-------------|
| charge.success | Payment successful |
| transfer.success | Withdrawal successful |
| transfer.failed | Withdrawal failed |

---

### Termii Webhook

Endpoint for SMS delivery reports.

```http
POST /webhooks/termii
```

**Event Types:**

| Status | Description |
|--------|-------------|
| delivered | SMS delivered |
| failed | SMS failed |
| expired | SMS expired |

---

## Appendix: Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 422 | Validation Error |
| 429 | Rate Limited |
| 500 | Internal Error |

## Appendix: Gig Categories

| Code | Name |
|------|------|
| technology | Technology & IT |
| design | Design & Creative |
| writing | Writing & Content |
| marketing | Marketing & Sales |
| business | Business Services |
| lifestyle | Lifestyle & Events |
| education | Education & Training |
| other | Other |

## Appendix: Transaction Categories

| Code | Description |
|------|-------------|
| deposit | Wallet funding |
| withdrawal | Bank withdrawal |
| transfer_in | Received transfer |
| transfer_out | Sent transfer |
| gig_payment | Gig earnings |
| gig_escrow | Escrow funding |
| savings_contribution | Circle contribution |
| savings_payout | Circle payout |
| loan_disbursement | Loan received |
| loan_repayment | Loan payment |
