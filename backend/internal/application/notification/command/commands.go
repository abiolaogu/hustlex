package command

import (
	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// SendNotification represents a command to send a notification
type SendNotification struct {
	UserID   string
	Type     string
	Channel  string
	Title    string
	Body     string
	Data     map[string]interface{}
	Priority string
}

func (c SendNotification) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

func (c SendNotification) GetType() aggregate.NotificationType {
	return aggregate.NotificationType(c.Type)
}

func (c SendNotification) GetChannel() aggregate.Channel {
	return aggregate.Channel(c.Channel)
}

func (c SendNotification) GetPriority() aggregate.Priority {
	if c.Priority == "" {
		return aggregate.PriorityNormal
	}
	return aggregate.Priority(c.Priority)
}

// SendNotificationResult contains the result
type SendNotificationResult struct {
	NotificationID string
	Status         string
	ProviderID     string
}

// SendOTP represents a command to send an OTP
type SendOTP struct {
	Phone         string
	ExpiryMinutes int
}

func (c SendOTP) GetPhone() (valueobject.PhoneNumber, error) {
	return valueobject.NewPhoneNumber(c.Phone)
}

// SendOTPResult contains the OTP result
type SendOTPResult struct {
	MessageID string
	ExpiresAt string
}

// SendBulkNotification sends notifications to multiple users
type SendBulkNotification struct {
	UserIDs  []string
	Type     string
	Channel  string
	Title    string
	Body     string
	Data     map[string]interface{}
	Priority string
}

// SendBulkNotificationResult contains bulk send results
type SendBulkNotificationResult struct {
	Sent   int
	Failed int
	Errors []string
}

// MarkNotificationRead marks a notification as read
type MarkNotificationRead struct {
	NotificationID string
	UserID         string
}

func (c MarkNotificationRead) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// MarkAllNotificationsRead marks all notifications as read
type MarkAllNotificationsRead struct {
	UserID string
}

func (c MarkAllNotificationsRead) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// DeleteNotification deletes a notification
type DeleteNotification struct {
	NotificationID string
	UserID         string
}

func (c DeleteNotification) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// RegisterDeviceToken registers a push notification device token
type RegisterDeviceToken struct {
	UserID   string
	Token    string
	Platform string // ios, android, web
	DeviceID string
}

func (c RegisterDeviceToken) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// RemoveDeviceToken removes a device token
type RemoveDeviceToken struct {
	UserID string
	Token  string
}

func (c RemoveDeviceToken) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// UpdateNotificationPreferences updates user notification preferences
type UpdateNotificationPreferences struct {
	UserID string

	// Channel preferences
	SMSEnabled   *bool
	EmailEnabled *bool
	PushEnabled  *bool
	InAppEnabled *bool

	// Type preferences
	TransactionAlerts *bool
	GigNotifications  *bool
	CircleUpdates     *bool
	LoanReminders     *bool
	Promotions        *bool
	SecurityAlerts    *bool

	// Quiet hours
	QuietHoursEnabled *bool
	QuietHoursStart   *string
	QuietHoursEnd     *string

	// Digest
	DailyDigest  *bool
	WeeklyReport *bool
}

func (c UpdateNotificationPreferences) GetUserID() (valueobject.UserID, error) {
	return valueobject.NewUserID(c.UserID)
}

// PreferencesResult contains the updated preferences
type PreferencesResult struct {
	SMSEnabled        bool
	EmailEnabled      bool
	PushEnabled       bool
	InAppEnabled      bool
	QuietHoursEnabled bool
	QuietHoursStart   string
	QuietHoursEnd     string
}

// RetryFailedNotifications retries failed notifications
type RetryFailedNotifications struct {
	Limit int
}

// RetryResult contains retry results
type RetryResult struct {
	Attempted int
	Succeeded int
	Failed    int
}
