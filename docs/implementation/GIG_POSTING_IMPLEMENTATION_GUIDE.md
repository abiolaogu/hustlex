# Gig Posting API Endpoint - Implementation Guide

**Created:** February 5, 2026
**Status:** Implementation Ready
**Priority:** HIGH (MVP Critical Path)
**Estimated Effort:** 3-5 days

---

## Executive Summary

This document provides a step-by-step guide to implement the gig posting API endpoint (`POST /api/gigs`), the first critical feature in the Core Gig Workflow. The domain layer and application layer are complete; this guide focuses on implementing the missing HTTP handler and PostgreSQL repository.

**Prerequisites Completed:**
- ✅ Domain aggregate (Gig) with business logic
- ✅ Application layer (CreateGig command & handler)
- ✅ Database schema (models.Gig ORM model)
- ✅ HTTP route defined in router.go

**Implementation Required:**
- ❌ HTTP handler for `POST /api/gigs`
- ❌ PostgreSQL repository implementation
- ❌ Integration tests

---

## Architecture Overview

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│ HTTP Layer (Interface)                               │
│ File: interface/http/handler/gig_handler.go         │
│ - Parse HTTP request                                 │
│ - Validate input                                     │
│ - Call application handler                           │
│ - Return HTTP response                               │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│ Application Layer                                    │
│ File: application/gig/handler/gig_handler.go        │
│ - HandleCreateGig() ✅ IMPLEMENTED                   │
│ - Orchestrate domain operations                      │
│ - Call repository                                    │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│ Domain Layer                                         │
│ File: domain/gig/aggregate/gig.go                   │
│ - NewGig() ✅ IMPLEMENTED                            │
│ - Business rules & validation                        │
│ - Generate domain events                             │
└──────────────────────┬──────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────┐
│ Infrastructure Layer                                 │
│ File: infrastructure/persistence/gig_repository.go  │
│ - GigRepository implementation ❌ TODO               │
│ - PostgreSQL persistence                             │
│ - Event publishing                                   │
└─────────────────────────────────────────────────────┘
```

---

## Step 1: Implement PostgreSQL Repository

**File to create:** `apps/api/internal/infrastructure/persistence/gig_repository.go`

### 1.1 Repository Structure

```go
package persistence

import (
    "context"
    "errors"
    "time"

    "gorm.io/gorm"
    "hustlex/internal/domain/gig/aggregate"
    "hustlex/internal/domain/gig/repository"
    "hustlex/internal/domain/shared/valueobject"
    "hustlex/internal/models"
)

type postgresGigRepository struct {
    db *gorm.DB
}

// NewPostgresGigRepository creates a new PostgreSQL gig repository
func NewPostgresGigRepository(db *gorm.DB) repository.GigRepository {
    return &postgresGigRepository{db: db}
}
```

### 1.2 Save Method Implementation

```go
// Save persists a gig aggregate
func (r *postgresGigRepository) Save(ctx context.Context, gig *aggregate.Gig) error {
    model := r.toModel(gig)

    // Use GORM's Save (upsert: create or update)
    if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
        return err
    }

    return nil
}

// SaveWithEvents persists gig and publishes domain events
func (r *postgresGigRepository) SaveWithEvents(ctx context.Context, gig *aggregate.Gig) error {
    // Start transaction
    tx := r.db.WithContext(ctx).Begin()
    if tx.Error != nil {
        return tx.Error
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Save gig
    model := r.toModel(gig)
    if err := tx.Save(&model).Error; err != nil {
        tx.Rollback()
        return err
    }

    // TODO: Publish domain events to message queue
    // For MVP, events can be logged or skipped
    events := gig.DomainEvents()
    if len(events) > 0 {
        // Example: Store events in an event store table
        // Or publish to RabbitMQ for async processing
        _ = events // Placeholder for event publishing
    }

    // Clear events after publishing
    gig.ClearDomainEvents()

    return tx.Commit().Error
}
```

### 1.3 Domain to Model Mapping

```go
// toModel converts domain aggregate to ORM model
func (r *postgresGigRepository) toModel(gig *aggregate.Gig) *models.Gig {
    // Extract skill ID (might be nil)
    var skillIDStr string
    if skillID := gig.SkillID(); skillID != nil {
        skillIDStr = skillID.String()
    }

    // Extract deadline (might be nil)
    var deadline *time.Time
    if d := gig.Deadline(); d != nil {
        deadline = d
    }

    return &models.Gig{
        ID:           gig.ID().String(),
        ClientID:     gig.ClientID().String(),
        Title:        gig.Title(),
        Description:  gig.Description(),
        Category:     gig.Category(),
        SkillID:      skillIDStr,
        BudgetMin:    gig.Budget().Min().Amount(),
        BudgetMax:    gig.Budget().Max().Amount(),
        Currency:     string(gig.Currency()),
        DeliveryDays: gig.DeliveryDays(),
        Deadline:     deadline,
        IsRemote:     gig.IsRemote(),
        Location:     gig.Location(),
        Tags:         gig.Tags(),
        Attachments:  gig.Attachments(),
        Status:       gig.Status().String(),
        ViewCount:    gig.ViewCount(),
        CreatedAt:    gig.CreatedAt(),
        UpdatedAt:    time.Now(),
    }
}
```

### 1.4 Retrieval Methods

```go
// FindByID retrieves a gig by ID
func (r *postgresGigRepository) FindByID(ctx context.Context, id valueobject.GigID) (*aggregate.Gig, error) {
    var model models.Gig

    err := r.db.WithContext(ctx).
        Where("id = ?", id.String()).
        First(&model).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrGigNotFound
        }
        return nil, err
    }

    return r.toDomain(&model)
}

// FindByClientID retrieves gigs posted by a client
func (r *postgresGigRepository) FindByClientID(ctx context.Context, clientID valueobject.UserID) ([]*aggregate.Gig, error) {
    var models []models.Gig

    err := r.db.WithContext(ctx).
        Where("client_id = ?", clientID.String()).
        Order("created_at DESC").
        Find(&models).Error

    if err != nil {
        return nil, err
    }

    gigs := make([]*aggregate.Gig, 0, len(models))
    for _, model := range models {
        gig, err := r.toDomain(&model)
        if err != nil {
            // Log error but continue with other gigs
            continue
        }
        gigs = append(gigs, gig)
    }

    return gigs, nil
}
```

### 1.5 Model to Domain Mapping

```go
// toDomain converts ORM model to domain aggregate
func (r *postgresGigRepository) toDomain(model *models.Gig) (*aggregate.Gig, error) {
    gigID, err := valueobject.NewGigID(model.ID)
    if err != nil {
        return nil, err
    }

    clientID, err := valueobject.NewUserID(model.ClientID)
    if err != nil {
        return nil, err
    }

    // Reconstruct budget
    budgetMin, err := valueobject.NewMoney(model.BudgetMin, valueobject.Currency(model.Currency))
    if err != nil {
        return nil, err
    }
    budgetMax, err := valueobject.NewMoney(model.BudgetMax, valueobject.Currency(model.Currency))
    if err != nil {
        return nil, err
    }
    budget, err := aggregate.NewBudget(budgetMin, budgetMax)
    if err != nil {
        return nil, err
    }

    // Reconstruct gig aggregate
    gig, err := aggregate.NewGig(
        gigID,
        clientID,
        model.Title,
        model.Description,
        model.Category,
        budget,
        model.DeliveryDays,
        model.IsRemote,
    )
    if err != nil {
        return nil, err
    }

    // Set additional fields via Update method
    var skillID *valueobject.SkillID
    if model.SkillID != "" {
        sid, err := valueobject.NewSkillID(model.SkillID)
        if err == nil {
            skillID = &sid
        }
    }

    if err := gig.Update(
        model.Title,
        model.Description,
        model.Category,
        skillID,
        budget,
        model.DeliveryDays,
        model.Deadline,
        model.IsRemote,
        model.Location,
        model.Attachments,
        model.Tags,
    ); err != nil {
        return nil, err
    }

    return gig, nil
}
```

### 1.6 List with Filtering

```go
// List retrieves gigs with filters, sorting, and pagination
func (r *postgresGigRepository) List(ctx context.Context, filter repository.GigFilter) ([]*aggregate.Gig, int, error) {
    query := r.db.WithContext(ctx).Model(&models.Gig{})

    // Apply filters
    if filter.Category != nil {
        query = query.Where("category = ?", *filter.Category)
    }
    if filter.SkillID != nil {
        query = query.Where("skill_id = ?", filter.SkillID.String())
    }
    if filter.MinBudget != nil {
        query = query.Where("budget_min >= ?", *filter.MinBudget)
    }
    if filter.MaxBudget != nil {
        query = query.Where("budget_max <= ?", *filter.MaxBudget)
    }
    if filter.Location != nil {
        query = query.Where("location ILIKE ?", "%"+*filter.Location+"%")
    }
    if filter.Status != nil {
        query = query.Where("status = ?", *filter.Status)
    }
    if filter.IsRemote != nil {
        query = query.Where("is_remote = ?", *filter.IsRemote)
    }
    if filter.SearchQuery != nil {
        searchPattern := "%" + *filter.SearchQuery + "%"
        query = query.Where(
            "title ILIKE ? OR description ILIKE ?",
            searchPattern, searchPattern,
        )
    }

    // Count total before pagination
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply sorting
    sortBy := "created_at"
    sortOrder := "DESC"
    if filter.SortBy != nil {
        sortBy = *filter.SortBy
    }
    if filter.SortOrder != nil {
        sortOrder = *filter.SortOrder
    }
    query = query.Order(sortBy + " " + sortOrder)

    // Apply pagination
    page := 1
    limit := 20
    if filter.Page != nil {
        page = *filter.Page
    }
    if filter.Limit != nil {
        limit = *filter.Limit
    }
    offset := (page - 1) * limit
    query = query.Offset(offset).Limit(limit)

    // Execute query
    var models []models.Gig
    if err := query.Find(&models).Error; err != nil {
        return nil, 0, err
    }

    // Convert to domain aggregates
    gigs := make([]*aggregate.Gig, 0, len(models))
    for _, model := range models {
        gig, err := r.toDomain(&model)
        if err != nil {
            continue
        }
        gigs = append(gigs, gig)
    }

    return gigs, int(total), nil
}
```

### 1.7 Delete Method

```go
// Delete soft-deletes a gig
func (r *postgresGigRepository) Delete(ctx context.Context, id valueobject.GigID) error {
    result := r.db.WithContext(ctx).
        Model(&models.Gig{}).
        Where("id = ?", id.String()).
        Update("status", "deleted")

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return repository.ErrGigNotFound
    }

    return nil
}
```

---

## Step 2: Implement HTTP Handler

**File to create:** `apps/api/internal/interface/http/handler/gig_handler.go`

### 2.1 Handler Structure

```go
package handler

import (
    "encoding/json"
    "net/http"
    "time"

    "hustlex/internal/application/gig/command"
    gighandler "hustlex/internal/application/gig/handler"
    "hustlex/internal/application/gig/query"
    "hustlex/internal/infrastructure/security/audit"
    "hustlex/internal/infrastructure/security/validation"
    "hustlex/internal/interface/http/middleware"
    "hustlex/internal/interface/http/response"
)

// GigHandler handles gig-related HTTP requests
type GigHandler struct {
    gigHandler      *gighandler.GigHandler
    contractHandler *gighandler.ContractHandler
    queryHandler    *query.GigQueryHandler
    auditLogger     audit.AuditLogger
}

// NewGigHandler creates a new HTTP gig handler
func NewGigHandler(
    gigHandler *gighandler.GigHandler,
    contractHandler *gighandler.ContractHandler,
    queryHandler *query.GigQueryHandler,
    auditLogger audit.AuditLogger,
) *GigHandler {
    return &GigHandler{
        gigHandler:      gigHandler,
        contractHandler: contractHandler,
        queryHandler:    queryHandler,
        auditLogger:     auditLogger,
    }
}
```

### 2.2 CreateGig HTTP Method

```go
// CreateGig handles POST /api/gigs
func (h *GigHandler) CreateGig(w http.ResponseWriter, r *http.Request) {
    // 1. Extract user ID from auth context
    userID, err := middleware.GetUserID(r.Context())
    if err != nil {
        response.Unauthorized(w, "unauthorized")
        return
    }

    // 2. Parse request body
    var req struct {
        Title        string    `json:"title"`
        Description  string    `json:"description"`
        Category     string    `json:"category"`
        SkillID      string    `json:"skill_id"`
        BudgetMin    int64     `json:"budget_min"`
        BudgetMax    int64     `json:"budget_max"`
        Currency     string    `json:"currency"`
        DeliveryDays int       `json:"delivery_days"`
        Deadline     *time.Time `json:"deadline"`
        IsRemote     bool      `json:"is_remote"`
        Location     string    `json:"location"`
        Tags         []string  `json:"tags"`
        Attachments  []string  `json:"attachments"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "invalid request body")
        return
    }

    // 3. Input validation
    v := validation.NewValidator()
    v.Required("title", req.Title).
        Required("description", req.Description).
        Required("category", req.Category).
        Positive("budget_min", req.BudgetMin).
        Positive("budget_max", req.BudgetMax).
        Min("budget_min", req.BudgetMin, 100).           // Minimum 100 kobo (1 Naira)
        Min("delivery_days", int64(req.DeliveryDays), 1). // At least 1 day
        Max("delivery_days", int64(req.DeliveryDays), 365). // Max 1 year
        SafeString("title", req.Title).
        SafeString("description", req.Description).
        SafeString("location", req.Location)

    // Validate budget range
    if req.BudgetMax < req.BudgetMin {
        v.AddError("budget_max", "budget_max must be greater than or equal to budget_min")
    }

    // Validate category
    validCategories := []string{
        "graphic_design", "content_writing", "digital_marketing",
        "video_editing", "web_development", "virtual_assistance",
        "photography", "event_planning", "tutoring", "beauty",
        "fashion", "delivery", "home_services", "accounting",
        "legal", "business_consulting", "translation",
    }
    v.OneOf("category", req.Category, validCategories)

    if v.HasErrors() {
        response.ValidationError(w, v.Errors().Errors)
        return
    }

    // Default currency to NGN
    if req.Currency == "" {
        req.Currency = "NGN"
    }

    // 4. Create command
    cmd := command.CreateGig{
        ClientID:     userID.String(),
        Title:        req.Title,
        Description:  req.Description,
        Category:     req.Category,
        SkillID:      req.SkillID,
        BudgetMin:    req.BudgetMin,
        BudgetMax:    req.BudgetMax,
        Currency:     req.Currency,
        DeliveryDays: req.DeliveryDays,
        Deadline:     req.Deadline,
        IsRemote:     req.IsRemote,
        Location:     req.Location,
        Tags:         req.Tags,
        Attachments:  req.Attachments,
    }

    // 5. Execute command via application handler
    result, err := h.gigHandler.HandleCreateGig(r.Context(), cmd)

    // 6. Audit logging
    if h.auditLogger != nil {
        outcome := audit.OutcomeSuccess
        message := "Gig created successfully"
        if err != nil {
            outcome = audit.OutcomeFailure
            message = "Gig creation failed"
        }
        h.auditLogger.LogTransaction(r.Context(), audit.AuditEvent{
            EventAction:    audit.ActionCreate,
            EventOutcome:   outcome,
            ActorUserID:    userID.String(),
            ActorIPAddress: getClientIP(r),
            ActorUserAgent: r.UserAgent(),
            TargetType:     "gig",
            TargetID:       "",
            Message:        message,
            Component:      "gig_handler",
            Metadata: map[string]interface{}{
                "category":      req.Category,
                "budget_min":    req.BudgetMin,
                "budget_max":    req.BudgetMax,
                "delivery_days": req.DeliveryDays,
            },
        })
    }

    if err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    // 7. Return success response
    response.Success(w, result)
}
```

### 2.3 Utility Functions

```go
// getClientIP extracts the real client IP from request headers
func getClientIP(r *http.Request) string {
    // Check X-Forwarded-For header (set by proxies)
    if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
        for i := 0; i < len(xff); i++ {
            if xff[i] == ',' {
                return xff[:i]
            }
        }
        return xff
    }

    // Check X-Real-IP header
    if xri := r.Header.Get("X-Real-IP"); xri != "" {
        return xri
    }

    // Fall back to RemoteAddr
    addr := r.RemoteAddr
    for i := len(addr) - 1; i >= 0; i-- {
        if addr[i] == ':' {
            return addr[:i]
        }
    }
    return addr
}
```

---

## Step 3: Wire Up in Router

**File to modify:** `apps/api/internal/interface/http/router/router.go`

### 3.1 Initialize Handler

Find the section where handlers are initialized (around line 80-100) and add:

```go
// Initialize Gig Repository
gigRepo := persistence.NewPostgresGigRepository(db)
proposalRepo := persistence.NewPostgresProposalRepository(db)  // TODO: Implement

// Initialize Application Handlers
gigAppHandler := gighandler.NewGigHandler(gigRepo, proposalRepo)
gigQueryHandler := query.NewGigQueryHandler(gigRepo, proposalRepo, nil)  // TODO: Add search repo

// Initialize HTTP Handler
gigHandler := handler.NewGigHandler(
    gigAppHandler,
    nil,  // Contract handler - TODO
    gigQueryHandler,
    auditLogger,
)
```

### 3.2 Update Route

Replace the `notImplemented` call with the actual handler (around line 173):

```go
// Before:
// r.Post("/gigs", notImplemented)

// After:
r.Post("/gigs", gigHandler.CreateGig)
```

---

## Step 4: Testing

### 4.1 Unit Test for Repository

**File:** `apps/api/internal/infrastructure/persistence/gig_repository_test.go`

```go
package persistence_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "hustlex/internal/domain/gig/aggregate"
    "hustlex/internal/domain/shared/valueobject"
    "hustlex/internal/infrastructure/persistence"
    "hustlex/internal/models"
    // ... test database setup imports
)

func TestGigRepository_Save(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer teardownTestDB(t, db)

    repo := persistence.NewPostgresGigRepository(db)

    // Create test gig
    gigID := valueobject.GenerateGigID()
    clientID, _ := valueobject.NewUserID("client-123")
    budgetMin, _ := valueobject.NewMoney(10000, valueobject.NGN)
    budgetMax, _ := valueobject.NewMoney(50000, valueobject.NGN)
    budget, _ := aggregate.NewBudget(budgetMin, budgetMax)

    gig, err := aggregate.NewGig(
        gigID,
        clientID,
        "Need a logo designed",
        "Looking for a modern logo for my startup",
        "graphic_design",
        budget,
        7,
        true,
    )
    require.NoError(t, err)

    // Test Save
    err = repo.Save(context.Background(), gig)
    assert.NoError(t, err)

    // Verify saved
    retrieved, err := repo.FindByID(context.Background(), gigID)
    assert.NoError(t, err)
    assert.Equal(t, gig.Title(), retrieved.Title())
    assert.Equal(t, gig.Category(), retrieved.Category())
}
```

### 4.2 Integration Test for HTTP Endpoint

**File:** `apps/api/internal/interface/http/handler/gig_handler_test.go`

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    // ... test setup imports
)

func TestCreateGig_Success(t *testing.T) {
    // Setup test server
    server := setupTestServer(t)
    defer server.Close()

    // Prepare request
    reqBody := map[string]interface{}{
        "title":         "Need a logo designed",
        "description":   "Looking for a modern logo for my startup",
        "category":      "graphic_design",
        "budget_min":    10000,
        "budget_max":    50000,
        "currency":      "NGN",
        "delivery_days": 7,
        "is_remote":     true,
    }

    body, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/gigs", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestJWT(t))

    // Execute request
    w := httptest.NewRecorder()
    server.ServeHTTP(w, req)

    // Assert response
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.NewDecoder(w.Body).Decode(&response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response["gig_id"])
    assert.Equal(t, "Need a logo designed", response["title"])
}

func TestCreateGig_ValidationError(t *testing.T) {
    server := setupTestServer(t)
    defer server.Close()

    // Missing required fields
    reqBody := map[string]interface{}{
        "title": "Test Gig",
        // Missing description, category, budget, etc.
    }

    body, _ := json.Marshal(reqBody)
    req := httptest.NewRequest("POST", "/api/gigs", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+getTestJWT(t))

    w := httptest.NewRecorder()
    server.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

---

## Step 5: Manual Testing with cURL

### 5.1 Start the Server

```bash
cd apps/api
go run cmd/server/main.go
```

### 5.2 Create a Test Gig

```bash
curl -X POST http://localhost:8081/api/gigs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Need a professional logo designed",
    "description": "Looking for a modern, minimalist logo for my tech startup. Must include both text and icon.",
    "category": "graphic_design",
    "skill_id": "skill-123",
    "budget_min": 15000,
    "budget_max": 50000,
    "currency": "NGN",
    "delivery_days": 7,
    "is_remote": true,
    "tags": ["logo", "branding", "modern"],
    "attachments": ["https://example.com/brief.pdf"]
  }'
```

### 5.3 Expected Response

```json
{
  "gig_id": "gig-abc123",
  "title": "Need a professional logo designed",
  "category": "graphic_design",
  "budget_min": 15000,
  "budget_max": 50000,
  "currency": "NGN",
  "delivery_days": 7,
  "status": "open",
  "created_at": "2026-02-05T22:30:00Z"
}
```

---

## Implementation Checklist

### Phase 1: Repository (Day 1-2)
- [ ] Create `infrastructure/persistence/gig_repository.go`
- [ ] Implement `Save()` and `SaveWithEvents()`
- [ ] Implement `FindByID()` and `FindByClientID()`
- [ ] Implement `List()` with filters
- [ ] Implement domain-to-model and model-to-domain mappers
- [ ] Write unit tests for repository

### Phase 2: HTTP Handler (Day 2-3)
- [ ] Create `interface/http/handler/gig_handler.go`
- [ ] Implement `CreateGig()` HTTP method
- [ ] Add request validation
- [ ] Add audit logging
- [ ] Wire up handler in router
- [ ] Write integration tests

### Phase 3: Testing (Day 3-4)
- [ ] Run unit tests
- [ ] Run integration tests
- [ ] Manual testing with cURL/Postman
- [ ] Test error cases (validation, auth, database errors)
- [ ] Performance testing (measure response time)

### Phase 4: Documentation (Day 4-5)
- [ ] Update API documentation (Swagger)
- [ ] Add code comments
- [ ] Update PRD implementation status
- [ ] Create deployment checklist

---

## Dependencies & Prerequisites

### Required Packages

```go
// Already in go.mod:
gorm.io/gorm
gorm.io/driver/postgres
github.com/go-chi/chi/v5
```

### Database Migration

Ensure the `gigs` table exists (should already be defined in Hasura migrations):

```sql
CREATE TABLE gigs (
    id VARCHAR(255) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    skill_id VARCHAR(255),
    budget_min BIGINT NOT NULL,
    budget_max BIGINT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'NGN',
    delivery_days INT NOT NULL,
    deadline TIMESTAMP,
    is_remote BOOLEAN NOT NULL DEFAULT true,
    location VARCHAR(255),
    tags TEXT[],
    attachments TEXT[],
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    view_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_gigs_client_id ON gigs(client_id);
CREATE INDEX idx_gigs_category ON gigs(category);
CREATE INDEX idx_gigs_status ON gigs(status);
CREATE INDEX idx_gigs_created_at ON gigs(created_at DESC);
```

---

## Error Handling

### Common Errors

| Error | HTTP Status | Response |
|-------|-------------|----------|
| Missing auth token | 401 | `{"error": "unauthorized"}` |
| Invalid request body | 400 | `{"error": "invalid request body"}` |
| Validation errors | 400 | `{"errors": {"field": "error message"}}` |
| Database error | 500 | `{"error": "internal server error"}` |
| Gig not found | 404 | `{"error": "gig not found"}` |

---

## Performance Considerations

1. **Database Indexes:** Ensure indexes exist on `client_id`, `category`, `status`, `created_at`
2. **Connection Pooling:** GORM handles this, but verify pool size in production
3. **Response Time Target:** < 200ms (p95) as per PRD
4. **Pagination:** Default limit is 20, max is 100
5. **Caching:** Consider caching gig listings in DragonflyDB for popular queries

---

## Security Considerations

1. **Authentication:** All gig posting requires valid JWT token
2. **Authorization:** Only gig owner can update/delete their gigs
3. **Input Validation:** Prevent SQL injection via parameterized queries (GORM handles this)
4. **XSS Prevention:** Sanitize user input (title, description, tags)
5. **Rate Limiting:** Consider rate limiting gig creation (e.g., 10 gigs/hour per user)
6. **Audit Logging:** All gig operations are logged with user ID, IP, timestamp

---

## Next Steps After Gig Posting

Once this endpoint is complete, implement in order:

1. **GET /api/gigs** - List gigs with filters
2. **GET /api/gigs/{id}** - Get single gig details
3. **POST /api/gigs/{id}/proposals** - Submit proposal
4. **POST /api/proposals/{id}/accept** - Accept proposal → create contract
5. **POST /api/contracts/{id}/deliver** - Submit work delivery
6. **POST /api/contracts/{id}/accept** - Approve delivery → release payment

---

## Support & Questions

- **Technical Lead:** Review code before merging
- **Database Admin:** Verify migrations and indexes
- **QA Team:** Execute test plan in staging environment
- **DevOps:** Deploy to staging first, monitor logs

---

**Document Version:** 1.0
**Last Updated:** February 5, 2026
**Next Review:** Upon implementation completion
