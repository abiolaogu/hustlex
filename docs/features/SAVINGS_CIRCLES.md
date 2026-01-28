# Feature Documentation: Savings Circles (Ajo/Esusu)

> Complete Technical Documentation for the Savings Circle Feature

---

## 1. Feature Overview

### 1.1 Purpose

Savings Circles digitize the traditional Nigerian Ajo/Esusu group savings system, providing:

- **Security:** Funds held in platform escrow, not with individuals
- **Transparency:** All contributions and payouts tracked
- **Automation:** Reminders, auto-debit, scheduled payouts
- **Credit Building:** Regular contributions improve credit score

### 1.2 Circle Types

| Type | Description | Use Case |
|------|-------------|----------|
| **Rotational** | Members take turns receiving pool | Planned expenses, rotating need |
| **Fixed Target** | Each member saves toward goal | Emergency fund, general savings |

### 1.3 Key Concepts

| Term | Definition |
|------|------------|
| Circle | A savings group |
| Member | Participant in circle |
| Position | Turn order in rotational circles |
| Contribution | Regular payment to circle |
| Payout | Receiving the pooled funds |
| Round | One complete contribution cycle |

---

## 2. Data Models

### 2.1 SavingsCircle Model

```go
type SavingsCircle struct {
    ID                  uuid.UUID `gorm:"type:uuid;primaryKey"`
    CreatorID           uuid.UUID `gorm:"type:uuid;index"`
    Name                string    `gorm:"size:100;not null"`
    Description         string    `gorm:"size:500"`
    Type                string    `gorm:"size:20;not null"` // rotational, fixed_target
    ContributionAmount  float64   `gorm:"not null"`
    Frequency           string    `gorm:"size:20;not null"` // weekly, monthly
    MaxMembers          int       `gorm:"not null"`
    CurrentMemberCount  int       `gorm:"default:0"`
    TotalPool           float64   `gorm:"default:0"`
    CurrentRound        int       `gorm:"default:0"`
    Status              string    `gorm:"size:20;default:'forming'"` // forming, active, completed, dissolved
    StartDate           time.Time
    EndDate             *time.Time
    NextContributionDate time.Time
    NextPayoutDate      *time.Time
    InviteCode          string    `gorm:"size:10;uniqueIndex"`
    Rules               JSONB     `gorm:"type:jsonb"`
    CreatedAt           time.Time
    UpdatedAt           time.Time

    // Relations
    Creator       User            `gorm:"foreignKey:CreatorID"`
    Members       []CircleMember  `gorm:"foreignKey:CircleID"`
    Contributions []Contribution  `gorm:"foreignKey:CircleID"`
}
```

### 2.2 CircleMember Model

```go
type CircleMember struct {
    ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
    CircleID       uuid.UUID `gorm:"type:uuid;index"`
    UserID         uuid.UUID `gorm:"type:uuid;index"`
    Position       int       // Order for rotational payout
    Role           string    `gorm:"size:20;default:'member'"` // admin, member
    Status         string    `gorm:"size:20;default:'active'"` // active, removed, left
    ReceivedPayout bool      `gorm:"default:false"`
    TotalContributed float64 `gorm:"default:0"`
    AutoContribute bool      `gorm:"default:false"`
    JoinedAt       time.Time
    LeftAt         *time.Time

    // Relations
    Circle  SavingsCircle  `gorm:"foreignKey:CircleID"`
    User    User           `gorm:"foreignKey:UserID"`
}
```

### 2.3 Contribution Model

```go
type Contribution struct {
    ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
    CircleID      uuid.UUID `gorm:"type:uuid;index"`
    MemberID      uuid.UUID `gorm:"type:uuid;index"`
    Round         int       `gorm:"not null"`
    Amount        float64   `gorm:"not null"`
    DueDate       time.Time `gorm:"index"`
    PaidAt        *time.Time
    Status        string    `gorm:"size:20;default:'pending'"` // pending, paid, late, missed
    TransactionID *uuid.UUID
    CreatedAt     time.Time

    // Relations
    Circle CircleMember `gorm:"foreignKey:CircleID"`
    Member CircleMember `gorm:"foreignKey:MemberID"`
}
```

---

## 3. Business Logic

### 3.1 Circle Lifecycle

```
FORMING
  │
  ├─── [Members join]
  │
  └─── [Start date reached + min members] ───> ACTIVE
                                                  │
                                                  ├─── [All rounds complete] ───> COMPLETED
                                                  │
                                                  └─── [Admin dissolves] ───> DISSOLVED
```

### 3.2 Rotational Circle Logic

```
Round 1:
├── All members contribute
├── Position 1 receives pool (₦100,000)
└── Next round scheduled

Round 2:
├── All members contribute
├── Position 2 receives pool (₦100,000)
└── Next round scheduled

... continues until all positions received

Round N (final):
├── All members contribute
├── Position N receives pool
└── Circle status → COMPLETED
```

### 3.3 Contribution Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    Contribution Process                          │
│                                                                  │
│  [Due Date - 3 days] ─── Send reminder notification             │
│                                                                  │
│  [Due Date - 1 day] ──── Send urgent reminder                   │
│                                                                  │
│  [Due Date] ─────────┬── Auto-contribute (if enabled)           │
│                      │   └── Debit wallet → Credit circle pool  │
│                      │                                           │
│                      └── Manual contribution expected            │
│                                                                  │
│  [Due Date + 1] ─────── Status: LATE, notify admin              │
│                                                                  │
│  [Due Date + 3] ─────── Final warning                           │
│                                                                  │
│  [Due Date + 7] ─────── Status: MISSED, consider removal        │
└─────────────────────────────────────────────────────────────────┘
```

### 3.4 Payout Processing

```go
func (s *SavingsService) ProcessPayout(circleID string) error {
    circle, _ := s.GetCircle(circleID)

    if circle.Type != "rotational" {
        return ErrInvalidCircleType
    }

    // Find recipient (member with current round position who hasn't received)
    recipient, err := s.GetPayoutRecipient(circleID, circle.CurrentRound)
    if err != nil {
        return err
    }

    // Verify all contributions collected
    pendingCount := s.GetPendingContributionCount(circleID, circle.CurrentRound)
    if pendingCount > 0 {
        return ErrContributionsPending
    }

    // Calculate payout amount
    payoutAmount := circle.ContributionAmount * float64(circle.CurrentMemberCount)

    return s.db.Transaction(func(tx *gorm.DB) error {
        // Credit recipient wallet from savings pool
        if err := s.walletService.CreditFromSavings(
            recipient.UserID.String(),
            payoutAmount,
            "savings_payout",
            fmt.Sprintf("Payout from %s (Round %d)", circle.Name, circle.CurrentRound),
        ); err != nil {
            return err
        }

        // Mark member as received
        tx.Model(&recipient).Update("received_payout", true)

        // Update circle
        circle.CurrentRound++
        circle.TotalPool = 0 // Reset for next round

        if circle.CurrentRound > circle.MaxMembers {
            circle.Status = "completed"
        } else {
            circle.NextPayoutDate = calculateNextPayoutDate(circle)
            circle.NextContributionDate = calculateNextContributionDate(circle)
        }

        tx.Save(&circle)

        // Notify recipient
        s.notificationService.Send(recipient.UserID.String(), Notification{
            Type:  "savings_payout",
            Title: "Payout Received!",
            Body:  fmt.Sprintf("You received ₦%s from %s", formatCurrency(payoutAmount), circle.Name),
        })

        // Notify all members
        s.notifyCircleMembers(circleID, Notification{
            Type:  "circle_update",
            Title: "Payout Completed",
            Body:  fmt.Sprintf("Round %d payout sent to %s", circle.CurrentRound-1, recipient.User.FullName),
        })

        return nil
    })
}
```

---

## 4. API Endpoints

### 4.1 Circle Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/savings/circles` | List circles (public + joined) |
| POST | `/savings/circles` | Create new circle |
| GET | `/savings/circles/:id` | Get circle details |
| PUT | `/savings/circles/:id` | Update circle (admin) |
| DELETE | `/savings/circles/:id` | Dissolve circle (admin) |

### 4.2 Membership Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/savings/circles/:id/join` | Join circle |
| POST | `/savings/circles/:id/leave` | Leave circle |
| GET | `/savings/circles/:id/members` | List members |
| PUT | `/savings/circles/:id/members/:mid` | Update member (admin) |
| DELETE | `/savings/circles/:id/members/:mid` | Remove member (admin) |

### 4.3 Contribution Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/savings/circles/:id/contribute` | Make contribution |
| GET | `/savings/circles/:id/contributions` | List contributions |
| GET | `/savings/my-contributions` | User's all contributions |

---

## 5. Background Jobs

### 5.1 Contribution Reminder Job

```go
const TaskContributionReminder = "savings:contribution_reminder"

type ContributionReminderPayload struct {
    CircleID     string    `json:"circle_id"`
    MemberID     string    `json:"member_id"`
    DueDate      time.Time `json:"due_date"`
    Amount       float64   `json:"amount"`
    ReminderNum  int       `json:"reminder_num"` // 1, 2, 3
}

func (p *JobProcessor) HandleContributionReminder(ctx context.Context, t *asynq.Task) error {
    var payload ContributionReminderPayload
    json.Unmarshal(t.Payload(), &payload)

    member, _ := p.db.GetCircleMember(payload.MemberID)
    circle, _ := p.db.GetCircle(payload.CircleID)

    // Check if already paid
    if p.isContributionPaid(payload.CircleID, payload.MemberID, circle.CurrentRound) {
        return nil
    }

    messages := map[int]string{
        1: fmt.Sprintf("Reminder: ₦%s due to %s in 3 days", formatCurrency(payload.Amount), circle.Name),
        2: fmt.Sprintf("Urgent: ₦%s due to %s tomorrow", formatCurrency(payload.Amount), circle.Name),
        3: fmt.Sprintf("Final: ₦%s due to %s TODAY", formatCurrency(payload.Amount), circle.Name),
    }

    p.notificationService.Send(member.UserID.String(), Notification{
        Type:  "contribution_reminder",
        Title: "Savings Reminder",
        Body:  messages[payload.ReminderNum],
    })

    return nil
}
```

### 5.2 Auto-Contribution Job

```go
const TaskProcessContribution = "savings:process_contribution"

func (p *JobProcessor) HandleProcessContribution(ctx context.Context, t *asynq.Task) error {
    var payload ProcessContributionPayload
    json.Unmarshal(t.Payload(), &payload)

    member, _ := p.db.GetCircleMember(payload.MemberID)

    // Check if auto-contribute enabled
    if !member.AutoContribute {
        return nil
    }

    // Check if already paid
    if p.isContributionPaid(payload.CircleID, payload.MemberID, payload.Round) {
        return nil
    }

    circle, _ := p.db.GetCircle(payload.CircleID)

    // Attempt to debit wallet
    err := p.walletService.DebitToSavings(
        member.UserID.String(),
        circle.ContributionAmount,
        fmt.Sprintf("Auto-contribution to %s", circle.Name),
    )

    if err != nil {
        // Notify user of failed auto-debit
        p.notificationService.Send(member.UserID.String(), Notification{
            Type:  "auto_contribute_failed",
            Title: "Auto-Contribution Failed",
            Body:  fmt.Sprintf("Insufficient balance for %s contribution", circle.Name),
        })
        return nil // Don't retry, manual contribution needed
    }

    // Record contribution
    contribution := &Contribution{
        CircleID: uuid.MustParse(payload.CircleID),
        MemberID: member.ID,
        Round:    payload.Round,
        Amount:   circle.ContributionAmount,
        DueDate:  payload.DueDate,
        PaidAt:   timePtr(time.Now()),
        Status:   "paid",
    }
    p.db.Create(contribution)

    // Update circle pool
    p.db.Model(&SavingsCircle{}).
        Where("id = ?", payload.CircleID).
        Update("total_pool", gorm.Expr("total_pool + ?", circle.ContributionAmount))

    // Update member total
    p.db.Model(&member).
        Update("total_contributed", gorm.Expr("total_contributed + ?", circle.ContributionAmount))

    return nil
}
```

### 5.3 Payout Processing Job

```go
const TaskProcessPayout = "savings:process_payout"

func (p *JobProcessor) HandleProcessPayout(ctx context.Context, t *asynq.Task) error {
    var payload PayoutPayload
    json.Unmarshal(t.Payload(), &payload)

    // Idempotency
    if !p.redis.SetNX(ctx, "payout:"+payload.CircleID+":"+payload.Round, "1", 24*time.Hour).Val() {
        return nil
    }

    return p.savingsService.ProcessPayout(payload.CircleID)
}
```

---

## 6. Credit Score Integration

### 6.1 Contribution Impact

| Activity | Credit Score Impact |
|----------|-------------------|
| Join first circle | +10 points |
| On-time contribution | +5 points |
| Late contribution (1-2 days) | 0 points |
| Late contribution (3+ days) | -10 points |
| Missed contribution | -25 points |
| Complete full cycle | +50 points |
| Create/admin circle | +25 points |

### 6.2 Savings Consistency Score

```go
func calculateSavingsConsistencyScore(userID string) int {
    contributions := getLastYearContributions(userID)

    totalExpected := len(contributions)
    onTime := 0
    late := 0
    missed := 0

    for _, c := range contributions {
        switch c.Status {
        case "paid":
            if c.PaidAt.Before(c.DueDate.Add(24 * time.Hour)) {
                onTime++
            } else {
                late++
            }
        case "missed":
            missed++
        }
    }

    // Score calculation (max 200)
    if totalExpected == 0 {
        return 100 // Neutral
    }

    onTimeRatio := float64(onTime) / float64(totalExpected)
    lateRatio := float64(late) / float64(totalExpected)
    missedRatio := float64(missed) / float64(totalExpected)

    score := int(onTimeRatio*200 - lateRatio*50 - missedRatio*100)
    return max(0, min(200, score))
}
```

---

## 7. Mobile Implementation

### 7.1 Feature Structure

```
features/savings/
├── data/
│   ├── models/
│   │   ├── circle_model.dart
│   │   ├── member_model.dart
│   │   └── contribution_model.dart
│   ├── repositories/
│   │   └── savings_repository.dart
└── presentation/
    ├── providers/
    │   └── savings_provider.dart
    ├── screens/
    │   ├── circles_list_screen.dart
    │   ├── circle_details_screen.dart
    │   ├── create_circle_screen.dart
    │   ├── join_circle_screen.dart
    │   └── contribute_screen.dart
    └── widgets/
        ├── circle_card.dart
        ├── member_list.dart
        └── contribution_chart.dart
```

### 7.2 Circle Details UI

```dart
class CircleDetailsScreen extends ConsumerWidget {
  final String circleId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final circleAsync = ref.watch(circleDetailsProvider(circleId));

    return Scaffold(
      appBar: AppBar(title: Text('Circle Details')),
      body: circleAsync.when(
        data: (circle) => SingleChildScrollView(
          child: Column(
            children: [
              // Header with pool amount
              CirclePoolHeader(
                name: circle.name,
                totalPool: circle.totalPool,
                type: circle.type,
              ),

              // Progress indicator
              if (circle.type == 'rotational')
                RoundProgressIndicator(
                  currentRound: circle.currentRound,
                  totalRounds: circle.maxMembers,
                ),

              // Next dates
              NextDatesCard(
                contributionDate: circle.nextContributionDate,
                payoutDate: circle.nextPayoutDate,
              ),

              // Members list
              MembersList(members: circle.members),

              // Recent contributions
              ContributionHistory(contributions: circle.contributions),

              // Action button
              ContributeButton(
                amount: circle.contributionAmount,
                onPressed: () => _showContributeDialog(context, ref, circle),
              ),
            ],
          ),
        ),
        loading: () => CircleDetailsSkeleton(),
        error: (e, _) => ErrorWidget(error: e),
      ),
    );
  }
}
```

---

## 8. Security Considerations

1. **Fund Security:** All funds held in platform escrow
2. **Position Fairness:** Random or voluntary position assignment
3. **Admin Controls:** Limited admin powers, audited actions
4. **Contribution Verification:** PIN required for manual contributions
5. **Member Verification:** KYC required for high-value circles
6. **Dispute Resolution:** Admin and support intervention available

---

## 9. Notifications

| Event | Recipients | Message |
|-------|------------|---------|
| Circle created | Creator | "Circle created! Share invite code" |
| Member joined | Admin, new member | "New member joined" |
| Contribution reminder | Member | "₦X due in Y days" |
| Contribution received | Member, admin | "Contribution recorded" |
| Auto-debit failed | Member | "Auto-contribution failed" |
| Payout sent | Recipient | "You received ₦X" |
| Round completed | All | "Round X complete" |
| Circle completed | All | "Circle completed!" |
| Member left/removed | All | "Member left circle" |

---

## 10. Metrics

| Metric | Description |
|--------|-------------|
| Active circles | Circles in 'active' status |
| Total pool value | Sum of all circle pools |
| Contribution rate | On-time contributions % |
| Completion rate | Circles that finish vs dissolve |
| Average circle size | Mean members per circle |
| Most popular frequency | Weekly vs monthly |
| Default rate | Members removed for non-payment |

---

*Feature Version 1.0 | Last Updated: January 2024*
