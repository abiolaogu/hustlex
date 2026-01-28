package event

import (
	"time"

	sharedevent "hustlex/internal/domain/shared/event"
)

// NotificationCreated is raised when a notification is created
type NotificationCreated struct {
	sharedevent.BaseEvent
	NotificationID string
	UserID         string
	Type           string
	Channel        string
	Priority       string
}

func NewNotificationCreated(notificationID, userID, nType, channel, priority string) NotificationCreated {
	return NotificationCreated{
		BaseEvent:      sharedevent.NewBaseEvent("notification.created"),
		NotificationID: notificationID,
		UserID:         userID,
		Type:           nType,
		Channel:        channel,
		Priority:       priority,
	}
}

// NotificationSent is raised when a notification is sent
type NotificationSent struct {
	sharedevent.BaseEvent
	NotificationID string
	UserID         string
	Channel        string
	ProviderID     string
	SentAt         time.Time
}

func NewNotificationSent(notificationID, userID, channel, providerID string, sentAt time.Time) NotificationSent {
	return NotificationSent{
		BaseEvent:      sharedevent.NewBaseEvent("notification.sent"),
		NotificationID: notificationID,
		UserID:         userID,
		Channel:        channel,
		ProviderID:     providerID,
		SentAt:         sentAt,
	}
}

// NotificationDelivered is raised when a notification is confirmed delivered
type NotificationDelivered struct {
	sharedevent.BaseEvent
	NotificationID string
	UserID         string
	Channel        string
	DeliveredAt    time.Time
}

func NewNotificationDelivered(notificationID, userID, channel string, deliveredAt time.Time) NotificationDelivered {
	return NotificationDelivered{
		BaseEvent:      sharedevent.NewBaseEvent("notification.delivered"),
		NotificationID: notificationID,
		UserID:         userID,
		Channel:        channel,
		DeliveredAt:    deliveredAt,
	}
}

// NotificationFailed is raised when a notification fails to send
type NotificationFailed struct {
	sharedevent.BaseEvent
	NotificationID string
	UserID         string
	Channel        string
	Error          string
	RetryCount     int
	CanRetry       bool
}

func NewNotificationFailed(notificationID, userID, channel, errMsg string, retryCount int, canRetry bool) NotificationFailed {
	return NotificationFailed{
		BaseEvent:      sharedevent.NewBaseEvent("notification.failed"),
		NotificationID: notificationID,
		UserID:         userID,
		Channel:        channel,
		Error:          errMsg,
		RetryCount:     retryCount,
		CanRetry:       canRetry,
	}
}

// NotificationRead is raised when a notification is read
type NotificationRead struct {
	sharedevent.BaseEvent
	NotificationID string
	UserID         string
	ReadAt         time.Time
}

func NewNotificationRead(notificationID, userID string, readAt time.Time) NotificationRead {
	return NotificationRead{
		BaseEvent:      sharedevent.NewBaseEvent("notification.read"),
		NotificationID: notificationID,
		UserID:         userID,
		ReadAt:         readAt,
	}
}

// OTPSent is raised when an OTP is sent
type OTPSent struct {
	sharedevent.BaseEvent
	UserID      string
	Phone       string
	MessageID   string
	ExpiresAt   time.Time
}

func NewOTPSent(userID, phone, messageID string, expiresAt time.Time) OTPSent {
	return OTPSent{
		BaseEvent: sharedevent.NewBaseEvent("notification.otp.sent"),
		UserID:    userID,
		Phone:     phone,
		MessageID: messageID,
		ExpiresAt: expiresAt,
	}
}

// DeviceTokenRegistered is raised when a push notification token is registered
type DeviceTokenRegistered struct {
	sharedevent.BaseEvent
	UserID    string
	DeviceID  string
	Platform  string
}

func NewDeviceTokenRegistered(userID, deviceID, platform string) DeviceTokenRegistered {
	return DeviceTokenRegistered{
		BaseEvent: sharedevent.NewBaseEvent("notification.device.registered"),
		UserID:    userID,
		DeviceID:  deviceID,
		Platform:  platform,
	}
}

// DeviceTokenRemoved is raised when a push notification token is removed
type DeviceTokenRemoved struct {
	sharedevent.BaseEvent
	UserID    string
	DeviceID  string
}

func NewDeviceTokenRemoved(userID, deviceID string) DeviceTokenRemoved {
	return DeviceTokenRemoved{
		BaseEvent: sharedevent.NewBaseEvent("notification.device.removed"),
		UserID:    userID,
		DeviceID:  deviceID,
	}
}

// PreferencesUpdated is raised when notification preferences are updated
type PreferencesUpdated struct {
	sharedevent.BaseEvent
	UserID  string
	Changes map[string]interface{}
}

func NewPreferencesUpdated(userID string, changes map[string]interface{}) PreferencesUpdated {
	return PreferencesUpdated{
		BaseEvent: sharedevent.NewBaseEvent("notification.preferences.updated"),
		UserID:    userID,
		Changes:   changes,
	}
}
