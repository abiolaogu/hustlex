# Feature Documentation: Gig Marketplace

> Complete Technical Documentation for the Gig Marketplace Feature

---

## 1. Feature Overview

### 1.1 Purpose

The Gig Marketplace enables Nigerian freelancers and clients to connect for short-term work opportunities. It provides:

- **For Clients:** Post gig requirements and find skilled freelancers
- **For Freelancers:** Discover work opportunities and earn income
- **For Platform:** Revenue through platform fees on completed gigs

### 1.2 Key Features

| Feature | Description |
|---------|-------------|
| Gig Posting | Clients post jobs with requirements and budget |
| Search & Filter | Find gigs by category, budget, skills |
| Proposals | Freelancers submit bids with cover letters |
| Contract Creation | Formalize agreement between parties |
| Escrow Payments | Secure funds until work completion |
| Delivery & Review | Submit work and receive ratings |

### 1.3 User Flows

```
CLIENT FLOW:
Post Gig → Review Proposals → Accept Proposal → Fund Escrow
    → Review Delivery → Approve/Request Revisions → Rate Freelancer

FREELANCER FLOW:
Browse Gigs → Submit Proposal → Contract Accepted
    → Complete Work → Submit Delivery → Receive Payment → Get Rated
```

---

## 2. Data Models

### 2.1 Gig Model

```go
type Gig struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
    ClientID        uuid.UUID `gorm:"type:uuid;index"`
    Title           string    `gorm:"size:200;not null"`
    Description     string    `gorm:"type:text;not null"`
    Category        string    `gorm:"size:50;index"`
    BudgetMin       float64   `gorm:"not null"`
    BudgetMax       float64   `gorm:"not null"`
    Deadline        time.Time
    Status          string    `gorm:"size:20;default:'open';index"` // open, in_progress, completed, cancelled
    Skills          []string  `gorm:"type:text[]"`
    Attachments     []string  `gorm:"type:text[]"`
    ViewsCount      int       `gorm:"default:0"`
    ProposalsCount  int       `gorm:"default:0"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt

    // Relations
    Client    User          `gorm:"foreignKey:ClientID"`
    Proposals []GigProposal `gorm:"foreignKey:GigID"`
    Contract  *GigContract  `gorm:"foreignKey:GigID"`
}
```

### 2.2 GigProposal Model

```go
type GigProposal struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
    GigID         uuid.UUID `gorm:"type:uuid;index"`
    FreelancerID  uuid.UUID `gorm:"type:uuid;index"`
    CoverLetter   string    `gorm:"type:text;not null"`
    ProposedPrice float64   `gorm:"not null"`
    DeliveryDays  int       `gorm:"not null"`
    Status        string    `gorm:"size:20;default:'pending'"` // pending, accepted, rejected, withdrawn
    Attachments   []string  `gorm:"type:text[]"`
    CreatedAt     time.Time
    UpdatedAt     time.Time

    // Relations
    Gig        Gig         `gorm:"foreignKey:GigID"`
    Freelancer User        `gorm:"foreignKey:FreelancerID"`
    Contract   *GigContract `gorm:"foreignKey:ProposalID"`
}
```

### 2.3 GigContract Model

```go
type GigContract struct {
    ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
    GigID            uuid.UUID `gorm:"type:uuid;index"`
    ProposalID       uuid.UUID `gorm:"type:uuid;index"`
    ClientID         uuid.UUID `gorm:"type:uuid;index"`
    FreelancerID     uuid.UUID `gorm:"type:uuid;index"`
    AgreedPrice      float64   `gorm:"not null"`
    PlatformFee      float64   `gorm:"not null"` // 10%
    FreelancerAmount float64   `gorm:"not null"` // 90%
    DeliveryDays     int       `gorm:"not null"`
    Deadline         time.Time
    Status           string    `gorm:"size:20;default:'in_progress'"` // in_progress, delivered, revision_requested, completed, disputed, cancelled
    DeliveredAt      *time.Time
    CompletedAt      *time.Time
    CreatedAt        time.Time
    UpdatedAt        time.Time

    // Relations
    Gig        Gig          `gorm:"foreignKey:GigID"`
    Proposal   GigProposal  `gorm:"foreignKey:ProposalID"`
    Client     User         `gorm:"foreignKey:ClientID"`
    Freelancer User         `gorm:"foreignKey:FreelancerID"`
    Deliveries []GigDelivery `gorm:"foreignKey:ContractID"`
    Reviews    []GigReview   `gorm:"foreignKey:ContractID"`
}
```

### 2.4 GigReview Model

```go
type GigReview struct {
    ID                uuid.UUID `gorm:"type:uuid;primaryKey"`
    ContractID        uuid.UUID `gorm:"type:uuid;index"`
    ReviewerID        uuid.UUID `gorm:"type:uuid;index"`
    RevieweeID        uuid.UUID `gorm:"type:uuid;index"`
    Rating            int       `gorm:"not null"` // 1-5
    ReviewText        string    `gorm:"type:text"`
    CommunicationRating int     // 1-5
    QualityRating       int     // 1-5
    TimelinessRating    int     // 1-5
    CreatedAt         time.Time

    // Relations
    Contract GigContract `gorm:"foreignKey:ContractID"`
    Reviewer User        `gorm:"foreignKey:ReviewerID"`
    Reviewee User        `gorm:"foreignKey:RevieweeID"`
}
```

---

## 3. API Endpoints

### 3.1 Gig Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/gigs` | List available gigs |
| POST | `/gigs` | Create new gig |
| GET | `/gigs/:id` | Get gig details |
| PUT | `/gigs/:id` | Update gig |
| DELETE | `/gigs/:id` | Cancel/delete gig |
| GET | `/gigs/my-gigs` | Get user's posted gigs |

### 3.2 Proposal Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/gigs/:id/proposals` | List proposals (client only) |
| POST | `/gigs/:id/proposals` | Submit proposal |
| PUT | `/gigs/:id/proposals/:pid` | Update proposal |
| DELETE | `/gigs/:id/proposals/:pid` | Withdraw proposal |
| POST | `/gigs/:id/proposals/:pid/accept` | Accept proposal |
| POST | `/gigs/:id/proposals/:pid/reject` | Reject proposal |

### 3.3 Contract Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/contracts` | List user's contracts |
| GET | `/contracts/:id` | Get contract details |
| POST | `/contracts/:id/deliver` | Submit delivery |
| POST | `/contracts/:id/approve` | Approve delivery |
| POST | `/contracts/:id/revision` | Request revision |
| POST | `/contracts/:id/dispute` | Open dispute |

### 3.4 Review Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/contracts/:id/review` | Submit review |
| GET | `/users/:id/reviews` | Get user reviews |

---

## 4. Business Logic

### 4.1 Gig Status Flow

```
OPEN
  │
  ├─── [Client cancels] ───> CANCELLED
  │
  └─── [Proposal accepted] ───> IN_PROGRESS
                                    │
                                    ├─── [Cancelled] ───> CANCELLED
                                    │
                                    └─── [Work delivered & approved] ───> COMPLETED
```

### 4.2 Contract Status Flow

```
IN_PROGRESS
  │
  ├─── [Freelancer delivers] ───> DELIVERED
  │                                   │
  │                                   ├─── [Client approves] ───> COMPLETED
  │                                   │
  │                                   └─── [Client requests revision] ───> REVISION_REQUESTED
  │                                                                              │
  │                                                                              └─── [Freelancer resubmits] ───> DELIVERED
  │
  ├─── [Dispute opened] ───> DISPUTED
  │
  └─── [Mutual cancellation] ───> CANCELLED
```

### 4.3 Escrow Logic

```
When proposal is accepted:
1. Client funds escrow (full agreed price)
2. Funds deducted from client wallet
3. Funds held in escrow_balance

When delivery is approved:
1. Calculate platform fee (10%)
2. Freelancer amount = agreed price - platform fee
3. Credit freelancer wallet
4. Clear escrow

On dispute resolution:
1. Admin reviews case
2. Determines refund split (0-100%)
3. Refund client portion
4. Release freelancer portion
```

### 4.4 Platform Fee Structure

```
Agreed Price: ₦100,000
Platform Fee: ₦10,000 (10%)
Freelancer Receives: ₦90,000

Fee Distribution:
- Platform revenue: ₦10,000
```

---

## 5. Implementation Details

### 5.1 Create Gig Handler

```go
func (h *GigHandler) Create(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)

    var req CreateGigRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse("Invalid request"))
    }

    // Validate
    if err := validate.Struct(req); err != nil {
        return c.Status(400).JSON(formatValidationErrors(err))
    }

    // Create gig
    gig, err := h.service.CreateGig(userID, req)
    if err != nil {
        return c.Status(500).JSON(ErrorResponse(err.Error()))
    }

    return c.Status(201).JSON(fiber.Map{
        "success": true,
        "data":    fiber.Map{"gig": gig},
    })
}
```

### 5.2 Accept Proposal Service

```go
func (s *GigService) AcceptProposal(clientID, proposalID string, agreedPrice float64, pin string) (*GigContract, error) {
    // Verify PIN
    if err := s.authService.VerifyPIN(clientID, pin); err != nil {
        return nil, ErrInvalidPIN
    }

    // Get proposal
    proposal, err := s.GetProposal(proposalID)
    if err != nil {
        return nil, err
    }

    // Verify client owns gig
    if proposal.Gig.ClientID.String() != clientID {
        return nil, ErrUnauthorized
    }

    // Calculate fees
    platformFee := agreedPrice * 0.10
    freelancerAmount := agreedPrice - platformFee

    // Start transaction
    return s.db.Transaction(func(tx *gorm.DB) (*GigContract, error) {
        // Fund escrow from client wallet
        if err := s.walletService.DebitToEscrow(clientID, agreedPrice); err != nil {
            return nil, ErrInsufficientFunds
        }

        // Create contract
        contract := &GigContract{
            GigID:            proposal.GigID,
            ProposalID:       proposal.ID,
            ClientID:         proposal.Gig.ClientID,
            FreelancerID:     proposal.FreelancerID,
            AgreedPrice:      agreedPrice,
            PlatformFee:      platformFee,
            FreelancerAmount: freelancerAmount,
            DeliveryDays:     proposal.DeliveryDays,
            Deadline:         time.Now().AddDate(0, 0, proposal.DeliveryDays),
            Status:           "in_progress",
        }

        if err := tx.Create(contract).Error; err != nil {
            return nil, err
        }

        // Update proposal status
        tx.Model(&proposal).Update("status", "accepted")

        // Update gig status
        tx.Model(&proposal.Gig).Update("status", "in_progress")

        // Reject other proposals
        tx.Model(&GigProposal{}).
            Where("gig_id = ? AND id != ?", proposal.GigID, proposal.ID).
            Update("status", "rejected")

        // Notify freelancer
        s.notificationService.Send(proposal.FreelancerID.String(), Notification{
            Type:  "proposal_accepted",
            Title: "Proposal Accepted!",
            Body:  fmt.Sprintf("Your proposal for '%s' was accepted", proposal.Gig.Title),
        })

        return contract, nil
    })
}
```

### 5.3 Escrow Release Job

```go
func (p *JobProcessor) HandleEscrowRelease(ctx context.Context, t *asynq.Task) error {
    var payload EscrowReleasePayload
    json.Unmarshal(t.Payload(), &payload)

    // Idempotency check
    if p.redis.SetNX(ctx, "escrow:released:"+payload.ContractID, "1", 24*time.Hour).Val() == false {
        return nil // Already processed
    }

    contract, _ := p.db.GetContract(payload.ContractID)

    // Credit freelancer
    if err := p.walletService.CreditFromEscrow(
        contract.FreelancerID.String(),
        contract.FreelancerAmount,
        "gig_payment",
        fmt.Sprintf("Payment for gig: %s", contract.Gig.Title),
    ); err != nil {
        return err
    }

    // Record platform revenue
    p.recordRevenue(contract.PlatformFee, "gig_fee", contract.ID.String())

    // Notify freelancer
    p.notificationService.Send(contract.FreelancerID.String(), Notification{
        Type:  "payment_received",
        Title: "Payment Received!",
        Body:  fmt.Sprintf("You received ₦%s for your gig", formatCurrency(contract.FreelancerAmount)),
    })

    return nil
}
```

---

## 6. Mobile Implementation

### 6.1 Gigs Feature Structure

```
features/gigs/
├── data/
│   ├── models/
│   │   ├── gig_model.dart
│   │   ├── proposal_model.dart
│   │   └── contract_model.dart
│   ├── repositories/
│   │   └── gigs_repository.dart
│   └── services/
│       └── gigs_service.dart
└── presentation/
    ├── providers/
    │   ├── gigs_provider.dart
    │   └── contracts_provider.dart
    ├── screens/
    │   ├── gigs_list_screen.dart
    │   ├── gig_details_screen.dart
    │   ├── create_gig_screen.dart
    │   ├── submit_proposal_screen.dart
    │   └── contract_screen.dart
    └── widgets/
        ├── gig_card.dart
        └── proposal_card.dart
```

### 6.2 Gigs Provider

```dart
@riverpod
class GigsNotifier extends _$GigsNotifier {
  @override
  FutureOr<GigsState> build() async {
    final gigs = await ref.watch(gigsRepositoryProvider).getGigs();
    return GigsState(gigs: gigs, filter: GigFilter());
  }

  Future<void> applyFilter(GigFilter filter) async {
    state = const AsyncValue.loading();
    state = await AsyncValue.guard(() async {
      final gigs = await ref.read(gigsRepositoryProvider).getGigs(
        category: filter.category,
        minBudget: filter.minBudget,
        maxBudget: filter.maxBudget,
      );
      return GigsState(gigs: gigs, filter: filter);
    });
  }

  Future<void> createGig(CreateGigRequest request) async {
    final gig = await ref.read(gigsRepositoryProvider).createGig(request);
    state = state.whenData((s) => s.copyWith(
      gigs: [gig, ...s.gigs],
    ));
  }
}
```

---

## 7. Search & Discovery

### 7.1 Full-Text Search

```sql
-- Create search index
CREATE INDEX idx_gigs_search ON gigs
USING gin(to_tsvector('english', title || ' ' || description));

-- Search query
SELECT *
FROM gigs
WHERE to_tsvector('english', title || ' ' || description)
      @@ plainto_tsquery('english', 'flutter mobile development')
  AND status = 'open'
ORDER BY ts_rank(
  to_tsvector('english', title || ' ' || description),
  plainto_tsquery('english', 'flutter mobile development')
) DESC;
```

### 7.2 Ranking Algorithm

```go
func rankGigs(gigs []Gig, query string) []Gig {
    for i := range gigs {
        gigs[i].Rank = calculateRank(gigs[i], query)
    }
    sort.Slice(gigs, func(i, j int) bool {
        return gigs[i].Rank > gigs[j].Rank
    })
    return gigs
}

func calculateRank(gig Gig, query string) float64 {
    var rank float64

    // Text relevance (from PostgreSQL ts_rank)
    rank += gig.TextRank * 10

    // Recency boost
    daysSinceCreated := time.Since(gig.CreatedAt).Hours() / 24
    rank += math.Max(0, 5-daysSinceCreated*0.1)

    // Client rating boost
    rank += gig.Client.Rating * 2

    return rank
}
```

---

## 8. Notifications

| Event | Recipient | Notification |
|-------|-----------|--------------|
| New proposal | Client | "New proposal on your gig" |
| Proposal accepted | Freelancer | "Your proposal was accepted!" |
| Proposal rejected | Freelancer | "Your proposal was not selected" |
| Work submitted | Client | "Delivery received - please review" |
| Revision requested | Freelancer | "Revision requested on your delivery" |
| Payment released | Freelancer | "You received ₦X for your gig" |
| Deadline approaching | Freelancer | "2 days left to deliver" |
| Review received | Both | "You received a new review" |

---

## 9. Security Considerations

1. **Authorization:** Users can only access their own gigs/proposals/contracts
2. **Escrow Protection:** Funds are locked and only released on approval
3. **Rate Limiting:** Prevent proposal spam
4. **Content Moderation:** Flag inappropriate gig content
5. **Dispute Resolution:** Admin intervention for conflicts

---

## 10. Metrics & Analytics

| Metric | Description |
|--------|-------------|
| Gigs posted | Daily/weekly/monthly count |
| Proposals per gig | Average proposals received |
| Acceptance rate | % of proposals accepted |
| Completion rate | % of contracts completed |
| Dispute rate | % of contracts disputed |
| Average gig value | Mean contract value |
| Platform revenue | Total fees collected |
| Time to hire | Average days from post to acceptance |

---

*Feature Version 1.0 | Last Updated: January 2024*
