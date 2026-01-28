package audit

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryAuditLogger_Log(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:  ActionCreate,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-123",
		TargetType:   "wallet",
		TargetID:     "wallet-456",
		Message:      "Created wallet",
	}

	err := logger.Log(ctx, event)
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	events := logger.Events()
	if len(events) != 1 {
		t.Fatalf("Events() count = %d, want 1", len(events))
	}

	logged := events[0]
	if logged.ID == "" {
		t.Error("Log() should generate ID")
	}
	if logged.Timestamp.IsZero() {
		t.Error("Log() should set timestamp")
	}
	if logged.Service != "test-service" {
		t.Errorf("Log() service = %s, want test-service", logged.Service)
	}
	if logged.ActorUserID != "user-123" {
		t.Errorf("Log() actor = %s, want user-123", logged.ActorUserID)
	}
}

func TestInMemoryAuditLogger_LogAccess(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:  ActionRead,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-123",
		TargetType:   "wallet",
		TargetID:     "wallet-456",
		Message:      "Viewed wallet balance",
	}

	err := logger.LogAccess(ctx, event)
	if err != nil {
		t.Fatalf("LogAccess() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeAccess {
		t.Errorf("LogAccess() type = %s, want ACCESS", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_LogDataChange(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:  ActionUpdate,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-123",
		TargetType:   "user",
		TargetID:     "user-123",
		Message:      "Updated profile",
		OldValue:     "old-name",
		NewValue:     "new-name",
	}

	err := logger.LogDataChange(ctx, event)
	if err != nil {
		t.Fatalf("LogDataChange() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeDataChange {
		t.Errorf("LogDataChange() type = %s, want DATA_CHANGE", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_LogAuthentication(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:    ActionExecute,
		EventOutcome:   OutcomeSuccess,
		ActorUserID:    "user-123",
		ActorIPAddress: "192.168.1.1",
		Message:        "User logged in",
	}

	err := logger.LogAuthentication(ctx, event)
	if err != nil {
		t.Fatalf("LogAuthentication() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeAuthentication {
		t.Errorf("LogAuthentication() type = %s, want AUTHENTICATION", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_LogAuthorization(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:  ActionExecute,
		EventOutcome: OutcomeFailure,
		ActorUserID:  "user-123",
		TargetType:   "admin_panel",
		Message:      "Unauthorized access attempt",
	}

	err := logger.LogAuthorization(ctx, event)
	if err != nil {
		t.Fatalf("LogAuthorization() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeAuthorization {
		t.Errorf("LogAuthorization() type = %s, want AUTHORIZATION", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_LogSecurityAlert(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:    ActionExecute,
		EventOutcome:   OutcomeFailure,
		ActorIPAddress: "192.168.1.1",
		Message:        "Multiple failed login attempts detected",
		Metadata: map[string]interface{}{
			"attempt_count": 5,
		},
	}

	err := logger.LogSecurityAlert(ctx, event)
	if err != nil {
		t.Fatalf("LogSecurityAlert() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeSecurityAlert {
		t.Errorf("LogSecurityAlert() type = %s, want SECURITY_ALERT", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_LogTransaction(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	event := AuditEvent{
		EventAction:  ActionCreate,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-123",
		TargetType:   "transaction",
		TargetID:     "txn-456",
		Message:      "Deposit completed",
		Metadata: map[string]interface{}{
			"amount":   10000,
			"currency": "NGN",
		},
	}

	err := logger.LogTransaction(ctx, event)
	if err != nil {
		t.Fatalf("LogTransaction() error: %v", err)
	}

	events := logger.Events()
	if events[0].EventType != EventTypeTransaction {
		t.Errorf("LogTransaction() type = %s, want TRANSACTION", events[0].EventType)
	}
}

func TestInMemoryAuditLogger_Query(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	// Add some events
	logger.LogAccess(ctx, AuditEvent{
		EventAction:  ActionRead,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-1",
		TargetType:   "wallet",
		TargetID:     "wallet-1",
		Message:      "Event 1",
	})
	logger.LogAuthentication(ctx, AuditEvent{
		EventAction:  ActionExecute,
		EventOutcome: OutcomeSuccess,
		ActorUserID:  "user-2",
		Message:      "Event 2",
	})
	logger.LogSecurityAlert(ctx, AuditEvent{
		EventAction:  ActionExecute,
		EventOutcome: OutcomeFailure,
		ActorUserID:  "user-1",
		Message:      "Event 3",
	})

	t.Run("filter by actor", func(t *testing.T) {
		events, total, err := logger.Query(ctx, AuditFilter{
			ActorUserID: "user-1",
		})
		if err != nil {
			t.Fatalf("Query() error: %v", err)
		}
		if total != 2 {
			t.Errorf("Query() total = %d, want 2", total)
		}
		if len(events) != 2 {
			t.Errorf("Query() count = %d, want 2", len(events))
		}
	})

	t.Run("filter by outcome", func(t *testing.T) {
		events, total, err := logger.Query(ctx, AuditFilter{
			Outcome: OutcomeFailure,
		})
		if err != nil {
			t.Fatalf("Query() error: %v", err)
		}
		if total != 1 {
			t.Errorf("Query() total = %d, want 1", total)
		}
		if len(events) != 1 {
			t.Errorf("Query() count = %d, want 1", len(events))
		}
	})

	t.Run("filter by target", func(t *testing.T) {
		events, _, err := logger.Query(ctx, AuditFilter{
			TargetType: "wallet",
			TargetID:   "wallet-1",
		})
		if err != nil {
			t.Fatalf("Query() error: %v", err)
		}
		if len(events) != 1 {
			t.Errorf("Query() count = %d, want 1", len(events))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		events, total, err := logger.Query(ctx, AuditFilter{
			Limit:  2,
			Offset: 0,
		})
		if err != nil {
			t.Fatalf("Query() error: %v", err)
		}
		if total != 3 {
			t.Errorf("Query() total = %d, want 3", total)
		}
		if len(events) != 2 {
			t.Errorf("Query() count = %d, want 2", len(events))
		}
	})
}

func TestInMemoryAuditLogger_Clear(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	logger.LogAccess(ctx, AuditEvent{Message: "test"})
	logger.LogAccess(ctx, AuditEvent{Message: "test"})

	if len(logger.Events()) != 2 {
		t.Fatal("should have 2 events")
	}

	logger.Clear()

	if len(logger.Events()) != 0 {
		t.Error("Clear() should remove all events")
	}
}

func TestContextEnrichment(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	// Event without actor info
	event := AuditEvent{
		EventAction:  ActionRead,
		EventOutcome: OutcomeSuccess,
		Message:      "Test event",
	}

	// Log the event using the logger
	err := logger.Log(ctx, event)
	if err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	// Verify the context keys are defined correctly
	if ContextKeyCorrelationID != "correlation_id" {
		t.Error("ContextKeyCorrelationID should be defined")
	}
	if ContextKeyRequestID != "request_id" {
		t.Error("ContextKeyRequestID should be defined")
	}
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		eventType EventType
		want      string
	}{
		{EventTypeAccess, "ACCESS"},
		{EventTypeDataChange, "DATA_CHANGE"},
		{EventTypeAuthentication, "AUTHENTICATION"},
		{EventTypeAuthorization, "AUTHORIZATION"},
		{EventTypeConfiguration, "CONFIGURATION"},
		{EventTypeSecurityAlert, "SECURITY_ALERT"},
		{EventTypeTransaction, "TRANSACTION"},
	}

	for _, tt := range tests {
		if string(tt.eventType) != tt.want {
			t.Errorf("EventType = %s, want %s", tt.eventType, tt.want)
		}
	}
}

func TestEventActions(t *testing.T) {
	tests := []struct {
		action EventAction
		want   string
	}{
		{ActionCreate, "C"},
		{ActionRead, "R"},
		{ActionUpdate, "U"},
		{ActionDelete, "D"},
		{ActionExecute, "E"},
	}

	for _, tt := range tests {
		if string(tt.action) != tt.want {
			t.Errorf("EventAction = %s, want %s", tt.action, tt.want)
		}
	}
}

func TestEventOutcomes(t *testing.T) {
	tests := []struct {
		outcome EventOutcome
		want    string
	}{
		{OutcomeSuccess, "SUCCESS"},
		{OutcomeFailure, "FAILURE"},
		{OutcomeError, "ERROR"},
	}

	for _, tt := range tests {
		if string(tt.outcome) != tt.want {
			t.Errorf("EventOutcome = %s, want %s", tt.outcome, tt.want)
		}
	}
}

func TestNullString(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"", false},
		{"value", true},
	}

	for _, tt := range tests {
		result := nullString(tt.input)
		if result.Valid != tt.valid {
			t.Errorf("nullString(%q).Valid = %v, want %v", tt.input, result.Valid, tt.valid)
		}
		if tt.valid && result.String != tt.input {
			t.Errorf("nullString(%q).String = %q, want %q", tt.input, result.String, tt.input)
		}
	}
}

func TestItoa(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{9, "9"},
		{10, "10"},
		{15, "15"},
		{99, "99"},
		{123, "123"},
	}

	for _, tt := range tests {
		got := itoa(tt.input)
		if got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestAuditFilter_TimeRange(t *testing.T) {
	logger := NewInMemoryAuditLogger("test-service")
	ctx := context.Background()

	now := time.Now()

	logger.Log(ctx, AuditEvent{Message: "old event"})
	time.Sleep(10 * time.Millisecond)
	logger.Log(ctx, AuditEvent{Message: "new event"})

	events, _, _ := logger.Query(ctx, AuditFilter{
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now.Add(1 * time.Hour),
	})

	if len(events) != 2 {
		t.Errorf("Query() with time range count = %d, want 2", len(events))
	}
}
