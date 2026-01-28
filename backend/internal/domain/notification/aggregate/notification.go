package aggregate

import (
	"errors"
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrInvalidNotificationType = errors.New("invalid notification type")
	ErrInvalidChannel          = errors.New("invalid notification channel")
	ErrNotificationAlreadySent = errors.New("notification already sent")
	ErrNotificationFailed      = errors.New("notification delivery failed")
)

// NotificationType represents the type of notification
type NotificationType string

const (
	TypeOTP             NotificationType = "otp"
	TypeTransaction     NotificationType = "transaction"
	TypeGigUpdate       NotificationType = "gig_update"
	TypeContractUpdate  NotificationType = "contract_update"
	TypePaymentReceived NotificationType = "payment_received"
	TypePaymentSent     NotificationType = "payment_sent"
	TypeCircleUpdate    NotificationType = "circle_update"
	TypeContribution    NotificationType = "contribution"
	TypePayout          NotificationType = "payout"
	TypeLoanUpdate      NotificationType = "loan_update"
	TypeCreditUpdate    NotificationType = "credit_update"
	TypePromotion       NotificationType = "promotion"
	TypeReminder        NotificationType = "reminder"
	TypeSystem          NotificationType = "system"
)

func (t NotificationType) String() string {
	return string(t)
}

func (t NotificationType) IsValid() bool {
	switch t {
	case TypeOTP, TypeTransaction, TypeGigUpdate, TypeContractUpdate,
		TypePaymentReceived, TypePaymentSent, TypeCircleUpdate,
		TypeContribution, TypePayout, TypeLoanUpdate, TypeCreditUpdate,
		TypePromotion, TypeReminder, TypeSystem:
		return true
	}
	return false
}

// Channel represents the delivery channel
type Channel string

const (
	ChannelSMS   Channel = "sms"
	ChannelEmail Channel = "email"
	ChannelPush  Channel = "push"
	ChannelInApp Channel = "in_app"
)

func (c Channel) String() string {
	return string(c)
}

func (c Channel) IsValid() bool {
	switch c {
	case ChannelSMS, ChannelEmail, ChannelPush, ChannelInApp:
		return true
	}
	return false
}

// DeliveryStatus represents the delivery status
type DeliveryStatus string

const (
	StatusPending   DeliveryStatus = "pending"
	StatusSent      DeliveryStatus = "sent"
	StatusDelivered DeliveryStatus = "delivered"
	StatusFailed    DeliveryStatus = "failed"
	StatusRead      DeliveryStatus = "read"
)

func (s DeliveryStatus) String() string {
	return string(s)
}

// Priority represents notification priority
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

func (p Priority) String() string {
	return string(p)
}

// Notification is the aggregate root for notifications
type Notification struct {
	sharedevent.AggregateRoot

	id           string
	userID       valueobject.UserID
	nType        NotificationType
	channel      Channel
	title        string
	body         string
	data         map[string]interface{} // Additional payload data
	priority     Priority
	status       DeliveryStatus
	providerID   string     // External provider message ID
	errorMessage string     // Error message if failed
	sentAt       *time.Time
	deliveredAt  *time.Time
	readAt       *time.Time
	expiresAt    *time.Time
	retryCount   int
	maxRetries   int
	createdAt    time.Time
	updatedAt    time.Time
}

// NewNotification creates a new notification
func NewNotification(
	id string,
	userID valueobject.UserID,
	nType NotificationType,
	channel Channel,
	title string,
	body string,
	priority Priority,
) (*Notification, error) {
	if !nType.IsValid() {
		return nil, ErrInvalidNotificationType
	}
	if !channel.IsValid() {
		return nil, ErrInvalidChannel
	}

	return &Notification{
		id:         id,
		userID:     userID,
		nType:      nType,
		channel:    channel,
		title:      title,
		body:       body,
		data:       make(map[string]interface{}),
		priority:   priority,
		status:     StatusPending,
		retryCount: 0,
		maxRetries: 3,
		createdAt:  time.Now().UTC(),
		updatedAt:  time.Now().UTC(),
	}, nil
}

// ReconstructNotification reconstructs from persistence
func ReconstructNotification(
	id string,
	userID valueobject.UserID,
	nType NotificationType,
	channel Channel,
	title string,
	body string,
	data map[string]interface{},
	priority Priority,
	status DeliveryStatus,
	providerID string,
	errorMessage string,
	sentAt *time.Time,
	deliveredAt *time.Time,
	readAt *time.Time,
	expiresAt *time.Time,
	retryCount int,
	maxRetries int,
	createdAt time.Time,
	updatedAt time.Time,
) *Notification {
	return &Notification{
		id:           id,
		userID:       userID,
		nType:        nType,
		channel:      channel,
		title:        title,
		body:         body,
		data:         data,
		priority:     priority,
		status:       status,
		providerID:   providerID,
		errorMessage: errorMessage,
		sentAt:       sentAt,
		deliveredAt:  deliveredAt,
		readAt:       readAt,
		expiresAt:    expiresAt,
		retryCount:   retryCount,
		maxRetries:   maxRetries,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (n *Notification) ID() string                       { return n.id }
func (n *Notification) UserID() valueobject.UserID       { return n.userID }
func (n *Notification) Type() NotificationType           { return n.nType }
func (n *Notification) Channel() Channel                 { return n.channel }
func (n *Notification) Title() string                    { return n.title }
func (n *Notification) Body() string                     { return n.body }
func (n *Notification) Data() map[string]interface{}     { return n.data }
func (n *Notification) Priority() Priority               { return n.priority }
func (n *Notification) Status() DeliveryStatus           { return n.status }
func (n *Notification) ProviderID() string               { return n.providerID }
func (n *Notification) ErrorMessage() string             { return n.errorMessage }
func (n *Notification) SentAt() *time.Time               { return n.sentAt }
func (n *Notification) DeliveredAt() *time.Time          { return n.deliveredAt }
func (n *Notification) ReadAt() *time.Time               { return n.readAt }
func (n *Notification) ExpiresAt() *time.Time            { return n.expiresAt }
func (n *Notification) RetryCount() int                  { return n.retryCount }
func (n *Notification) MaxRetries() int                  { return n.maxRetries }
func (n *Notification) CreatedAt() time.Time             { return n.createdAt }
func (n *Notification) UpdatedAt() time.Time             { return n.updatedAt }

// Business Methods

// SetData sets additional payload data
func (n *Notification) SetData(data map[string]interface{}) {
	n.data = data
	n.updatedAt = time.Now().UTC()
}

// AddData adds a single data field
func (n *Notification) AddData(key string, value interface{}) {
	if n.data == nil {
		n.data = make(map[string]interface{})
	}
	n.data[key] = value
	n.updatedAt = time.Now().UTC()
}

// SetExpiry sets when the notification expires
func (n *Notification) SetExpiry(expiresAt time.Time) {
	n.expiresAt = &expiresAt
	n.updatedAt = time.Now().UTC()
}

// MarkSent marks the notification as sent
func (n *Notification) MarkSent(providerID string) error {
	if n.status == StatusSent || n.status == StatusDelivered {
		return ErrNotificationAlreadySent
	}

	now := time.Now().UTC()
	n.status = StatusSent
	n.providerID = providerID
	n.sentAt = &now
	n.updatedAt = now

	return nil
}

// MarkDelivered marks the notification as delivered
func (n *Notification) MarkDelivered() {
	now := time.Now().UTC()
	n.status = StatusDelivered
	n.deliveredAt = &now
	n.updatedAt = now
}

// MarkRead marks the notification as read
func (n *Notification) MarkRead() {
	if n.readAt != nil {
		return
	}
	now := time.Now().UTC()
	n.status = StatusRead
	n.readAt = &now
	n.updatedAt = now
}

// MarkFailed marks the notification as failed
func (n *Notification) MarkFailed(errorMsg string) {
	n.status = StatusFailed
	n.errorMessage = errorMsg
	n.retryCount++
	n.updatedAt = time.Now().UTC()
}

// CanRetry checks if the notification can be retried
func (n *Notification) CanRetry() bool {
	return n.status == StatusFailed && n.retryCount < n.maxRetries
}

// ResetForRetry resets the notification for a retry attempt
func (n *Notification) ResetForRetry() {
	n.status = StatusPending
	n.errorMessage = ""
	n.updatedAt = time.Now().UTC()
}

// IsExpired checks if the notification has expired
func (n *Notification) IsExpired() bool {
	if n.expiresAt == nil {
		return false
	}
	return time.Now().After(*n.expiresAt)
}

// IsRead checks if the notification has been read
func (n *Notification) IsRead() bool {
	return n.readAt != nil
}
