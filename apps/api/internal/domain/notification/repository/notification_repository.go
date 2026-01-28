package repository

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/shared/valueobject"
)

// Repository errors
var (
	ErrNotificationNotFound   = errors.New("notification not found")
	ErrPreferencesNotFound    = errors.New("notification preferences not found")
	ErrDeviceTokenNotFound    = errors.New("device token not found")
)

// NotificationRepository defines the interface for notification persistence
type NotificationRepository interface {
	// Save persists a notification
	Save(ctx context.Context, notification *aggregate.Notification) error

	// FindByID retrieves a notification by ID
	FindByID(ctx context.Context, id string) (*aggregate.Notification, error)

	// FindByUserID retrieves notifications for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID, filter NotificationFilter) ([]*aggregate.Notification, int64, error)

	// FindUnread retrieves unread notifications for a user
	FindUnread(ctx context.Context, userID valueobject.UserID) ([]*aggregate.Notification, error)

	// FindPending retrieves pending notifications for retry
	FindPending(ctx context.Context, limit int) ([]*aggregate.Notification, error)

	// FindFailed retrieves failed notifications for retry
	FindFailed(ctx context.Context, limit int) ([]*aggregate.Notification, error)

	// CountUnread returns the count of unread notifications
	CountUnread(ctx context.Context, userID valueobject.UserID) (int, error)

	// MarkAllRead marks all notifications as read for a user
	MarkAllRead(ctx context.Context, userID valueobject.UserID) error

	// Delete deletes a notification
	Delete(ctx context.Context, id string) error

	// DeleteOld deletes notifications older than a given date
	DeleteOld(ctx context.Context, before time.Time) (int64, error)
}

// NotificationFilter contains filter options
type NotificationFilter struct {
	Type     *aggregate.NotificationType
	Channel  *aggregate.Channel
	Status   *aggregate.DeliveryStatus
	IsRead   *bool
	FromDate *time.Time
	ToDate   *time.Time
	Offset   int
	Limit    int
}

// NotificationDTO represents notification data for API responses
type NotificationDTO struct {
	ID          string
	UserID      string
	Type        string
	Channel     string
	Title       string
	Body        string
	Data        map[string]interface{}
	Priority    string
	Status      string
	IsRead      bool
	SentAt      *time.Time
	DeliveredAt *time.Time
	ReadAt      *time.Time
	CreatedAt   time.Time
}

// PreferencesRepository defines the interface for notification preferences
type PreferencesRepository interface {
	// Save persists notification preferences
	Save(ctx context.Context, prefs *aggregate.NotificationPreferences) error

	// FindByUserID retrieves preferences for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID) (*aggregate.NotificationPreferences, error)

	// Update updates notification preferences
	Update(ctx context.Context, prefs *aggregate.NotificationPreferences) error
}

// DeviceTokenRepository manages push notification device tokens
type DeviceTokenRepository interface {
	// Save persists a device token
	Save(ctx context.Context, token *DeviceToken) error

	// FindByUserID retrieves device tokens for a user
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*DeviceToken, error)

	// FindByToken retrieves a device token
	FindByToken(ctx context.Context, token string) (*DeviceToken, error)

	// Delete deletes a device token
	Delete(ctx context.Context, token string) error

	// DeleteByUserID deletes all device tokens for a user
	DeleteByUserID(ctx context.Context, userID valueobject.UserID) error
}

// DeviceToken represents a push notification device token
type DeviceToken struct {
	ID        string
	UserID    string
	Token     string
	Platform  string // ios, android, web
	DeviceID  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NotificationStatisticsRepository provides notification analytics
type NotificationStatisticsRepository interface {
	// GetUserStats gets notification statistics for a user
	GetUserStats(ctx context.Context, userID valueobject.UserID) (*UserNotificationStats, error)

	// GetPlatformStats gets platform-wide notification statistics
	GetPlatformStats(ctx context.Context, period time.Duration) (*PlatformNotificationStats, error)
}

// UserNotificationStats contains user notification statistics
type UserNotificationStats struct {
	UserID            string
	TotalNotifications int
	UnreadCount       int
	LastNotificationAt *time.Time
	ByChannel         map[string]int
	ByType            map[string]int
}

// PlatformNotificationStats contains platform-wide statistics
type PlatformNotificationStats struct {
	TotalSent     int64
	TotalDelivered int64
	TotalFailed   int64
	DeliveryRate  float64
	ByChannel     map[string]int64
	ByType        map[string]int64
	ByStatus      map[string]int64
}
