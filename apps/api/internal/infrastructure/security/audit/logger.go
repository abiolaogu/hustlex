package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of audit event
type EventType string

const (
	EventTypeAccess         EventType = "ACCESS"
	EventTypeDataChange     EventType = "DATA_CHANGE"
	EventTypeAuthentication EventType = "AUTHENTICATION"
	EventTypeAuthorization  EventType = "AUTHORIZATION"
	EventTypeConfiguration  EventType = "CONFIGURATION"
	EventTypeSecurityAlert  EventType = "SECURITY_ALERT"
	EventTypeTransaction    EventType = "TRANSACTION"
)

// EventAction represents CRUD operations
type EventAction string

const (
	ActionCreate  EventAction = "C"
	ActionRead    EventAction = "R"
	ActionUpdate  EventAction = "U"
	ActionDelete  EventAction = "D"
	ActionExecute EventAction = "E"
)

// EventOutcome represents the result of an action
type EventOutcome string

const (
	OutcomeSuccess EventOutcome = "SUCCESS"
	OutcomeFailure EventOutcome = "FAILURE"
	OutcomeError   EventOutcome = "ERROR"
)

// Context keys for audit information
type contextKey string

const (
	ContextKeyCorrelationID contextKey = "correlation_id"
	ContextKeyRequestID     contextKey = "request_id"
	ContextKeyUserID        contextKey = "user_id"
	ContextKeySessionID     contextKey = "session_id"
	ContextKeyIPAddress     contextKey = "ip_address"
	ContextKeyUserAgent     contextKey = "user_agent"
)

// AuditEvent represents a single audit log entry
// Follows RFC 3881 / DICOM / IHE audit message format
type AuditEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	EventType EventType `json:"event_type"`
	// C=Create, R=Read, U=Update, D=Delete, E=Execute
	EventAction  EventAction  `json:"event_action"`
	EventOutcome EventOutcome `json:"event_outcome"`

	// Actor (who)
	ActorUserID    string `json:"actor_user_id,omitempty"`
	ActorUserName  string `json:"actor_user_name,omitempty"`
	ActorIPAddress string `json:"actor_ip_address,omitempty"`
	ActorUserAgent string `json:"actor_user_agent,omitempty"`
	ActorSessionID string `json:"actor_session_id,omitempty"`

	// Target (what)
	TargetType string `json:"target_type,omitempty"`
	TargetID   string `json:"target_id,omitempty"`
	TargetName string `json:"target_name,omitempty"`

	// Context
	CorrelationID string `json:"correlation_id,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
	Service       string `json:"service"`
	Component     string `json:"component,omitempty"`

	// Details
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// For data changes - careful not to log PII
	OldValue string `json:"old_value,omitempty"` // Redacted/hashed
	NewValue string `json:"new_value,omitempty"` // Redacted/hashed
}

// AuditFilter defines filters for querying audit logs
type AuditFilter struct {
	StartTime   time.Time
	EndTime     time.Time
	EventTypes  []EventType
	ActorUserID string
	TargetType  string
	TargetID    string
	Outcome     EventOutcome
	Limit       int
	Offset      int
}

// AuditLogger interface for audit logging
type AuditLogger interface {
	Log(ctx context.Context, event AuditEvent) error
	LogAccess(ctx context.Context, event AuditEvent) error
	LogDataChange(ctx context.Context, event AuditEvent) error
	LogAuthentication(ctx context.Context, event AuditEvent) error
	LogAuthorization(ctx context.Context, event AuditEvent) error
	LogSecurityAlert(ctx context.Context, event AuditEvent) error
	LogTransaction(ctx context.Context, event AuditEvent) error
	Query(ctx context.Context, filter AuditFilter) ([]AuditEvent, int64, error)
}

// PostgresAuditLogger implements AuditLogger with PostgreSQL
type PostgresAuditLogger struct {
	db      *sql.DB
	service string
}

// NewPostgresAuditLogger creates a new PostgreSQL-backed audit logger
func NewPostgresAuditLogger(db *sql.DB, service string) *PostgresAuditLogger {
	return &PostgresAuditLogger{
		db:      db,
		service: service,
	}
}

// enrichFromContext enriches the event with context information
func (l *PostgresAuditLogger) enrichFromContext(ctx context.Context, event *AuditEvent) {
	event.ID = uuid.NewString()
	event.Timestamp = time.Now().UTC()
	event.Service = l.service

	// Extract correlation ID from context
	if corrID, ok := ctx.Value(ContextKeyCorrelationID).(string); ok {
		event.CorrelationID = corrID
	}
	if reqID, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		event.RequestID = reqID
	}
	if event.ActorUserID == "" {
		if userID, ok := ctx.Value(ContextKeyUserID).(string); ok {
			event.ActorUserID = userID
		}
	}
	if event.ActorSessionID == "" {
		if sessionID, ok := ctx.Value(ContextKeySessionID).(string); ok {
			event.ActorSessionID = sessionID
		}
	}
	if event.ActorIPAddress == "" {
		if ip, ok := ctx.Value(ContextKeyIPAddress).(string); ok {
			event.ActorIPAddress = ip
		}
	}
	if event.ActorUserAgent == "" {
		if ua, ok := ctx.Value(ContextKeyUserAgent).(string); ok {
			event.ActorUserAgent = ua
		}
	}
}

// Log writes an audit event to the database
func (l *PostgresAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	l.enrichFromContext(ctx, &event)

	// Serialize metadata
	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO audit_logs (
			id, timestamp, event_type, event_action, event_outcome,
			actor_user_id, actor_user_name, actor_ip_address, actor_user_agent, actor_session_id,
			target_type, target_id, target_name,
			correlation_id, request_id, service, component,
			message, metadata, old_value, new_value, created_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	_, err = l.db.ExecContext(ctx, query,
		event.ID, event.Timestamp, event.EventType, event.EventAction, event.EventOutcome,
		nullString(event.ActorUserID), nullString(event.ActorUserName),
		nullString(event.ActorIPAddress), nullString(event.ActorUserAgent), nullString(event.ActorSessionID),
		nullString(event.TargetType), nullString(event.TargetID), nullString(event.TargetName),
		nullString(event.CorrelationID), nullString(event.RequestID), event.Service, nullString(event.Component),
		event.Message, metadataJSON, nullString(event.OldValue), nullString(event.NewValue),
		event.Timestamp.Format("2006-01-02"),
	)

	return err
}

// LogAccess logs an access event
func (l *PostgresAuditLogger) LogAccess(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAccess
	return l.Log(ctx, event)
}

// LogDataChange logs a data change event
func (l *PostgresAuditLogger) LogDataChange(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeDataChange
	return l.Log(ctx, event)
}

// LogAuthentication logs an authentication event
func (l *PostgresAuditLogger) LogAuthentication(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAuthentication
	return l.Log(ctx, event)
}

// LogAuthorization logs an authorization event
func (l *PostgresAuditLogger) LogAuthorization(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAuthorization
	return l.Log(ctx, event)
}

// LogSecurityAlert logs a security alert event
func (l *PostgresAuditLogger) LogSecurityAlert(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeSecurityAlert
	// TODO: Implement alerting integration (PagerDuty, Slack, etc.)
	return l.Log(ctx, event)
}

// LogTransaction logs a financial transaction event
func (l *PostgresAuditLogger) LogTransaction(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeTransaction
	return l.Log(ctx, event)
}

// Query retrieves audit events based on filters
func (l *PostgresAuditLogger) Query(ctx context.Context, filter AuditFilter) ([]AuditEvent, int64, error) {
	// Build query with filters
	query := `
		SELECT
			id, timestamp, event_type, event_action, event_outcome,
			actor_user_id, actor_user_name, actor_ip_address, actor_user_agent, actor_session_id,
			target_type, target_id, target_name,
			correlation_id, request_id, service, component,
			message, metadata, old_value, new_value
		FROM audit_logs
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`

	var args []interface{}
	argIndex := 1

	if !filter.StartTime.IsZero() {
		query += ` AND timestamp >= $` + itoa(argIndex)
		countQuery += ` AND timestamp >= $` + itoa(argIndex)
		args = append(args, filter.StartTime)
		argIndex++
	}
	if !filter.EndTime.IsZero() {
		query += ` AND timestamp <= $` + itoa(argIndex)
		countQuery += ` AND timestamp <= $` + itoa(argIndex)
		args = append(args, filter.EndTime)
		argIndex++
	}
	if len(filter.EventTypes) > 0 {
		query += ` AND event_type = ANY($` + itoa(argIndex) + `)`
		countQuery += ` AND event_type = ANY($` + itoa(argIndex) + `)`
		args = append(args, filter.EventTypes)
		argIndex++
	}
	if filter.ActorUserID != "" {
		query += ` AND actor_user_id = $` + itoa(argIndex)
		countQuery += ` AND actor_user_id = $` + itoa(argIndex)
		args = append(args, filter.ActorUserID)
		argIndex++
	}
	if filter.TargetType != "" {
		query += ` AND target_type = $` + itoa(argIndex)
		countQuery += ` AND target_type = $` + itoa(argIndex)
		args = append(args, filter.TargetType)
		argIndex++
	}
	if filter.TargetID != "" {
		query += ` AND target_id = $` + itoa(argIndex)
		countQuery += ` AND target_id = $` + itoa(argIndex)
		args = append(args, filter.TargetID)
		argIndex++
	}
	if filter.Outcome != "" {
		query += ` AND event_outcome = $` + itoa(argIndex)
		countQuery += ` AND event_outcome = $` + itoa(argIndex)
		args = append(args, filter.Outcome)
		argIndex++
	}

	// Get total count
	var total int64
	err := l.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Add ordering and pagination
	query += ` ORDER BY timestamp DESC`
	if filter.Limit > 0 {
		query += ` LIMIT $` + itoa(argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}
	if filter.Offset > 0 {
		query += ` OFFSET $` + itoa(argIndex)
		args = append(args, filter.Offset)
	}

	rows, err := l.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []AuditEvent
	for rows.Next() {
		var event AuditEvent
		var metadataJSON []byte
		var actorUserID, actorUserName, actorIP, actorUA, actorSession sql.NullString
		var targetType, targetID, targetName sql.NullString
		var corrID, reqID, component, oldVal, newVal sql.NullString

		err := rows.Scan(
			&event.ID, &event.Timestamp, &event.EventType, &event.EventAction, &event.EventOutcome,
			&actorUserID, &actorUserName, &actorIP, &actorUA, &actorSession,
			&targetType, &targetID, &targetName,
			&corrID, &reqID, &event.Service, &component,
			&event.Message, &metadataJSON, &oldVal, &newVal,
		)
		if err != nil {
			return nil, 0, err
		}

		event.ActorUserID = actorUserID.String
		event.ActorUserName = actorUserName.String
		event.ActorIPAddress = actorIP.String
		event.ActorUserAgent = actorUA.String
		event.ActorSessionID = actorSession.String
		event.TargetType = targetType.String
		event.TargetID = targetID.String
		event.TargetName = targetName.String
		event.CorrelationID = corrID.String
		event.RequestID = reqID.String
		event.Component = component.String
		event.OldValue = oldVal.String
		event.NewValue = newVal.String

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &event.Metadata)
		}

		events = append(events, event)
	}

	return events, total, rows.Err()
}

// Helper functions

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return itoa(i/10) + string(rune('0'+i%10))
}

// InMemoryAuditLogger is a simple in-memory implementation for testing
type InMemoryAuditLogger struct {
	events  []AuditEvent
	service string
}

// NewInMemoryAuditLogger creates a new in-memory audit logger
func NewInMemoryAuditLogger(service string) *InMemoryAuditLogger {
	return &InMemoryAuditLogger{
		events:  make([]AuditEvent, 0),
		service: service,
	}
}

func (l *InMemoryAuditLogger) Log(ctx context.Context, event AuditEvent) error {
	event.ID = uuid.NewString()
	event.Timestamp = time.Now().UTC()
	event.Service = l.service
	l.events = append(l.events, event)
	return nil
}

func (l *InMemoryAuditLogger) LogAccess(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAccess
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) LogDataChange(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeDataChange
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) LogAuthentication(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAuthentication
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) LogAuthorization(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeAuthorization
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) LogSecurityAlert(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeSecurityAlert
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) LogTransaction(ctx context.Context, event AuditEvent) error {
	event.EventType = EventTypeTransaction
	return l.Log(ctx, event)
}

func (l *InMemoryAuditLogger) Query(ctx context.Context, filter AuditFilter) ([]AuditEvent, int64, error) {
	var filtered []AuditEvent
	for _, e := range l.events {
		if !filter.StartTime.IsZero() && e.Timestamp.Before(filter.StartTime) {
			continue
		}
		if !filter.EndTime.IsZero() && e.Timestamp.After(filter.EndTime) {
			continue
		}
		if filter.ActorUserID != "" && e.ActorUserID != filter.ActorUserID {
			continue
		}
		if filter.TargetType != "" && e.TargetType != filter.TargetType {
			continue
		}
		if filter.TargetID != "" && e.TargetID != filter.TargetID {
			continue
		}
		if filter.Outcome != "" && e.EventOutcome != filter.Outcome {
			continue
		}
		filtered = append(filtered, e)
	}

	total := int64(len(filtered))

	// Apply pagination
	if filter.Offset > 0 && filter.Offset < len(filtered) {
		filtered = filtered[filter.Offset:]
	}
	if filter.Limit > 0 && filter.Limit < len(filtered) {
		filtered = filtered[:filter.Limit]
	}

	return filtered, total, nil
}

// Events returns all logged events (for testing)
func (l *InMemoryAuditLogger) Events() []AuditEvent {
	return l.events
}

// Clear clears all logged events (for testing)
func (l *InMemoryAuditLogger) Clear() {
	l.events = make([]AuditEvent, 0)
}
