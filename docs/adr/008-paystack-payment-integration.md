# ADR-008: Paystack as Primary Payment Gateway

## Status

Accepted

## Date

2024-01-15

## Context

HustleX requires payment processing for:
- Wallet funding (bank transfer, card payments)
- Escrow payments for gig contracts
- Savings circle contributions
- Loan disbursements and repayments
- Withdrawal to bank accounts

Requirements:
1. Support Nigerian banks and payment methods
2. Card payments (Visa, Mastercard, Verve)
3. Bank transfers (direct debit, USSD)
4. Instant settlement or predictable settlement cycles
5. Webhook support for payment notifications
6. Split payments for marketplace fees
7. Regulatory compliance (CBN licensed)

## Decision

We chose **Paystack** as our primary payment gateway with **Flutterwave** as secondary/backup.

### Key Reasons:

1. **Nigerian Market Leader**: Dominant in Nigeria, trusted by users, excellent local bank support.

2. **Comprehensive API**: Cards, bank transfers, USSD, QR codes, mobile money.

3. **Split Payments**: Native support for marketplace commission deduction (10% platform fee).

4. **Instant Verification**: BVN/NIN verification for KYC compliance.

5. **Reliable Webhooks**: Real-time payment notifications with retry logic.

6. **Developer Experience**: Excellent documentation, SDKs, test environment.

7. **Regulatory Compliance**: CBN licensed, PCI-DSS compliant.

## Consequences

### Positive

- **Trust**: Users familiar with Paystack checkout
- **Coverage**: 99%+ Nigerian banks supported
- **Speed**: Instant card payments, same-day settlements
- **Features**: Transfers API for withdrawals, subscriptions for recurring payments
- **Support**: 24/7 merchant support, dedicated account manager
- **Mobile SDK**: Native Flutter package available

### Negative

- **Fees**: 1.5% + ₦100 (capped at ₦2000) per transaction
- **Vendor lock-in**: Paystack-specific API, migration effort if switching
- **Settlement delays**: T+1 for some payment methods
- **International limitations**: Primarily Nigeria-focused

### Neutral

- Requires business verification for live mode
- Webhook endpoint security responsibility
- Refund processing time varies

## Implementation Details

### Payment Channels

| Channel | Use Case | Settlement |
|---------|----------|------------|
| Card | Quick deposits, subscriptions | Instant |
| Bank Transfer | Large deposits, lower fees | T+0 to T+1 |
| USSD | Users without smartphones | T+0 |
| QR Code | In-person transactions | Instant |

### Integration Architecture

```
Mobile App
    │
    ▼
Paystack Flutter SDK (initialize transaction)
    │
    ▼
Paystack Checkout (card entry, bank selection)
    │
    ▼
Paystack API (process payment)
    │
    ├──────────────────────────────┐
    ▼                              ▼
Webhook (payment.success)    Mobile callback
    │
    ▼
HustleX Backend (verify + credit wallet)
```

### Wallet Funding Flow

```go
// 1. Initialize transaction
func InitializeDeposit(userID string, amount float64) (*PaystackInit, error) {
    ref := generateReference()

    resp, err := paystack.Transaction.Initialize(&paystack.TransactionRequest{
        Amount:      int(amount * 100), // Kobo
        Email:       user.Email,
        Reference:   ref,
        CallbackURL: "https://api.hustlex.app/webhooks/paystack",
        Metadata: map[string]interface{}{
            "user_id":     userID,
            "type":        "wallet_deposit",
            "custom_fields": []map[string]string{
                {"display_name": "User ID", "value": userID},
            },
        },
    })

    return resp, err
}

// 2. Webhook handler
func HandlePaystackWebhook(c *fiber.Ctx) error {
    // Verify signature
    signature := c.Get("X-Paystack-Signature")
    if !verifyWebhookSignature(c.Body(), signature) {
        return c.SendStatus(401)
    }

    var event PaystackEvent
    json.Unmarshal(c.Body(), &event)

    switch event.Event {
    case "charge.success":
        return handleChargeSuccess(event.Data)
    case "transfer.success":
        return handleTransferSuccess(event.Data)
    case "transfer.failed":
        return handleTransferFailed(event.Data)
    }

    return c.SendStatus(200)
}

// 3. Credit wallet
func handleChargeSuccess(data ChargeData) error {
    // Verify transaction with Paystack
    verified, err := paystack.Transaction.Verify(data.Reference)
    if err != nil || verified.Status != "success" {
        return err
    }

    // Credit user wallet (idempotent)
    return walletService.CreditWallet(
        data.Metadata.UserID,
        float64(data.Amount) / 100,
        data.Reference,
        "Wallet deposit via Paystack",
    )
}
```

### Escrow for Gig Payments

```go
// Client funds escrow when accepting proposal
func FundEscrow(contractID string, amount float64) error {
    contract, _ := db.GetContract(contractID)

    // Create subaccount for split payment
    // 90% to escrow, 10% platform fee
    return paystack.Transaction.Initialize(&paystack.TransactionRequest{
        Amount:       int(amount * 100),
        Email:        client.Email,
        Reference:    "escrow-" + contractID,
        SubAccount:   config.EscrowSubaccountCode,
        Bearer:       "account", // Platform bears fees
        TransactionCharge: int(amount * 0.10 * 100), // 10% platform fee
    })
}

// Release escrow to freelancer on completion
func ReleaseEscrow(contractID string) error {
    contract, _ := db.GetContract(contractID)

    // Transfer to freelancer's bank account
    return paystack.Transfer.Initiate(&paystack.TransferRequest{
        Source:    "balance",
        Amount:    int(contract.FreelancerAmount * 100),
        Recipient: freelancer.PaystackRecipientCode,
        Reason:    "Payment for gig: " + contract.GigTitle,
        Reference: "release-" + contractID,
    })
}
```

### Withdrawal to Bank

```go
// Create transfer recipient (one-time setup)
func CreateBankRecipient(userID string, bankCode string, accountNumber string) error {
    resp, err := paystack.TransferRecipient.Create(&paystack.RecipientRequest{
        Type:          "nuban",
        Name:          user.FullName,
        AccountNumber: accountNumber,
        BankCode:      bankCode,
        Currency:      "NGN",
    })

    // Store recipient code for future transfers
    return db.UpdateUser(userID, map[string]interface{}{
        "paystack_recipient_code": resp.RecipientCode,
    })
}

// Process withdrawal
func ProcessWithdrawal(userID string, amount float64) error {
    user, _ := db.GetUser(userID)

    // Deduct from wallet first
    if err := walletService.DebitWallet(userID, amount, "withdrawal"); err != nil {
        return err
    }

    // Initiate transfer
    return paystack.Transfer.Initiate(&paystack.TransferRequest{
        Source:    "balance",
        Amount:    int(amount * 100),
        Recipient: user.PaystackRecipientCode,
        Reason:    "HustleX wallet withdrawal",
        Reference: generateWithdrawalReference(userID),
    })
}
```

## Fee Structure

| Transaction Type | Paystack Fee | Platform Fee | User Pays |
|-----------------|--------------|--------------|-----------|
| Card Deposit | 1.5% + ₦100 (max ₦2000) | 0% | Paystack fee |
| Bank Transfer Deposit | Free | 0% | Free |
| Escrow (Client) | 1.5% + ₦100 | 10% of gig value | Both |
| Withdrawal | ₦10 - ₦50 | 0% | Flat fee |

## Alternatives Considered

### Alternative 1: Flutterwave

**Pros**: International coverage, competitive rates, good API
**Cons**: Less trusted in Nigeria, occasional settlement delays

**Decision**: Keep as secondary gateway for redundancy and international payments.

### Alternative 2: Interswitch/Quickteller

**Pros**: Legacy player, bank partnerships
**Cons**: Outdated API, poor developer experience, slow integration

**Rejected because**: Developer experience and API quality significantly inferior.

### Alternative 3: Direct Bank Integration (NIBSS)

**Pros**: Lower fees, direct settlement
**Cons**: Complex integration, requires banking license, per-bank integration

**Rejected because**: Regulatory and technical complexity too high for current stage.

### Alternative 4: Stripe

**Pros**: Excellent API, global coverage
**Cons**: Limited Nigerian bank support, no USSD, higher fees for local payments

**Rejected because**: Paystack better suited for Nigeria-focused operations.

## References

- [Paystack API Documentation](https://paystack.com/docs/api/)
- [Paystack Flutter SDK](https://pub.dev/packages/flutter_paystack)
- [Paystack Webhooks Guide](https://paystack.com/docs/payments/webhooks/)
- [CBN Payment Service Providers](https://www.cbn.gov.ng/paymentsystem/psp.asp)
