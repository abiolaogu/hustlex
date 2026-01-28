package service

import (
	"context"
	"errors"

	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrProviderUnavailable = errors.New("notification provider is unavailable")
	ErrInvalidRecipient    = errors.New("invalid notification recipient")
	ErrRateLimited         = errors.New("notification rate limit exceeded")
)

// NotificationProvider defines the interface for external notification providers
type NotificationProvider interface {
	// Send sends a notification through this provider
	Send(ctx context.Context, notification *aggregate.Notification) (providerID string, err error)

	// GetStatus checks the delivery status from the provider
	GetStatus(ctx context.Context, providerID string) (aggregate.DeliveryStatus, error)

	// SupportsChannel returns true if the provider supports this channel
	SupportsChannel(channel aggregate.Channel) bool
}

// SMSProvider defines the interface for SMS providers (e.g., Termii)
type SMSProvider interface {
	NotificationProvider

	// SendOTP sends an OTP via SMS
	SendOTP(ctx context.Context, phone valueobject.PhoneNumber, otp string, expiryMinutes int) (messageID string, err error)

	// SendBulk sends SMS to multiple recipients
	SendBulk(ctx context.Context, phones []valueobject.PhoneNumber, message string) ([]string, error)

	// GetBalance returns the SMS balance
	GetBalance(ctx context.Context) (float64, error)
}

// EmailProvider defines the interface for email providers
type EmailProvider interface {
	NotificationProvider

	// SendTemplate sends a templated email
	SendTemplate(ctx context.Context, email valueobject.Email, templateID string, data map[string]interface{}) (messageID string, err error)

	// SendBulk sends emails to multiple recipients
	SendBulk(ctx context.Context, emails []valueobject.Email, subject, body string) ([]string, error)
}

// PushProvider defines the interface for push notification providers
type PushProvider interface {
	NotificationProvider

	// SendToDevice sends a push notification to a specific device
	SendToDevice(ctx context.Context, deviceToken string, notification *aggregate.Notification) (messageID string, err error)

	// SendToUser sends a push notification to all user's devices
	SendToUser(ctx context.Context, userID valueobject.UserID, notification *aggregate.Notification) ([]string, error)

	// SendToTopic sends a push notification to a topic
	SendToTopic(ctx context.Context, topic string, notification *aggregate.Notification) (messageID string, err error)
}

// NotificationService orchestrates notification delivery across channels
type NotificationService struct {
	smsProvider   SMSProvider
	emailProvider EmailProvider
	pushProvider  PushProvider
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	smsProvider SMSProvider,
	emailProvider EmailProvider,
	pushProvider PushProvider,
) *NotificationService {
	return &NotificationService{
		smsProvider:   smsProvider,
		emailProvider: emailProvider,
		pushProvider:  pushProvider,
	}
}

// Send sends a notification through the appropriate channel
func (s *NotificationService) Send(ctx context.Context, notification *aggregate.Notification) error {
	var provider NotificationProvider

	switch notification.Channel() {
	case aggregate.ChannelSMS:
		provider = s.smsProvider
	case aggregate.ChannelEmail:
		provider = s.emailProvider
	case aggregate.ChannelPush:
		provider = s.pushProvider
	case aggregate.ChannelInApp:
		// In-app notifications are just stored, no external provider needed
		return notification.MarkSent("in_app")
	default:
		return aggregate.ErrInvalidChannel
	}

	if provider == nil {
		return ErrProviderUnavailable
	}

	providerID, err := provider.Send(ctx, notification)
	if err != nil {
		notification.MarkFailed(err.Error())
		return err
	}

	return notification.MarkSent(providerID)
}

// SendOTP sends an OTP via SMS
func (s *NotificationService) SendOTP(ctx context.Context, phone valueobject.PhoneNumber, otp string, expiryMinutes int) (string, error) {
	if s.smsProvider == nil {
		return "", ErrProviderUnavailable
	}
	return s.smsProvider.SendOTP(ctx, phone, otp, expiryMinutes)
}

// GetDeliveryStatus checks the delivery status of a notification
func (s *NotificationService) GetDeliveryStatus(ctx context.Context, notification *aggregate.Notification) (aggregate.DeliveryStatus, error) {
	if notification.ProviderID() == "" {
		return notification.Status(), nil
	}

	var provider NotificationProvider

	switch notification.Channel() {
	case aggregate.ChannelSMS:
		provider = s.smsProvider
	case aggregate.ChannelEmail:
		provider = s.emailProvider
	case aggregate.ChannelPush:
		provider = s.pushProvider
	case aggregate.ChannelInApp:
		return notification.Status(), nil
	}

	if provider == nil {
		return notification.Status(), nil
	}

	return provider.GetStatus(ctx, notification.ProviderID())
}

// NotificationBuilder helps construct notifications
type NotificationBuilder struct {
	userID   valueobject.UserID
	nType    aggregate.NotificationType
	channel  aggregate.Channel
	title    string
	body     string
	data     map[string]interface{}
	priority aggregate.Priority
}

// NewNotificationBuilder creates a new notification builder
func NewNotificationBuilder(userID valueobject.UserID) *NotificationBuilder {
	return &NotificationBuilder{
		userID:   userID,
		priority: aggregate.PriorityNormal,
		data:     make(map[string]interface{}),
	}
}

// WithType sets the notification type
func (b *NotificationBuilder) WithType(nType aggregate.NotificationType) *NotificationBuilder {
	b.nType = nType
	return b
}

// WithChannel sets the delivery channel
func (b *NotificationBuilder) WithChannel(channel aggregate.Channel) *NotificationBuilder {
	b.channel = channel
	return b
}

// WithTitle sets the title
func (b *NotificationBuilder) WithTitle(title string) *NotificationBuilder {
	b.title = title
	return b
}

// WithBody sets the body
func (b *NotificationBuilder) WithBody(body string) *NotificationBuilder {
	b.body = body
	return b
}

// WithData sets additional data
func (b *NotificationBuilder) WithData(key string, value interface{}) *NotificationBuilder {
	b.data[key] = value
	return b
}

// WithPriority sets the priority
func (b *NotificationBuilder) WithPriority(priority aggregate.Priority) *NotificationBuilder {
	b.priority = priority
	return b
}

// Build creates the notification
func (b *NotificationBuilder) Build() (*aggregate.Notification, error) {
	id := valueobject.GenerateUserID().String()
	notification, err := aggregate.NewNotification(
		id,
		b.userID,
		b.nType,
		b.channel,
		b.title,
		b.body,
		b.priority,
	)
	if err != nil {
		return nil, err
	}

	if len(b.data) > 0 {
		notification.SetData(b.data)
	}

	return notification, nil
}
