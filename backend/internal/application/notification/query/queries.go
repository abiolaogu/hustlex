package query

import (
	"context"
	"time"

	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/notification/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// GetNotification retrieves a single notification
type GetNotification struct {
	NotificationID string
	UserID         string
}

// NotificationDTO represents notification data for API responses
type NotificationDTO struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Channel     string                 `json:"channel"`
	Title       string                 `json:"title"`
	Body        string                 `json:"body"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Priority    string                 `json:"priority"`
	Status      string                 `json:"status"`
	IsRead      bool                   `json:"is_read"`
	SentAt      *time.Time             `json:"sent_at,omitempty"`
	DeliveredAt *time.Time             `json:"delivered_at,omitempty"`
	ReadAt      *time.Time             `json:"read_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// GetMyNotifications retrieves user's notifications
type GetMyNotifications struct {
	UserID  string
	Type    string
	Channel string
	IsRead  *bool
	Page    int
	Limit   int
}

// NotificationListResult represents paginated notification results
type NotificationListResult struct {
	Notifications []NotificationDTO `json:"notifications"`
	Total         int64             `json:"total"`
	UnreadCount   int               `json:"unread_count"`
	Page          int               `json:"page"`
	Limit         int               `json:"limit"`
	TotalPages    int               `json:"total_pages"`
}

// GetUnreadNotifications retrieves unread notifications
type GetUnreadNotifications struct {
	UserID string
}

// GetUnreadCount returns the count of unread notifications
type GetUnreadCount struct {
	UserID string
}

// UnreadCountResult contains the unread count
type UnreadCountResult struct {
	Count int `json:"count"`
}

// GetNotificationPreferences retrieves user's notification preferences
type GetNotificationPreferences struct {
	UserID string
}

// PreferencesDTO represents notification preferences for API responses
type PreferencesDTO struct {
	SMSEnabled        bool   `json:"sms_enabled"`
	EmailEnabled      bool   `json:"email_enabled"`
	PushEnabled       bool   `json:"push_enabled"`
	InAppEnabled      bool   `json:"in_app_enabled"`
	TransactionAlerts bool   `json:"transaction_alerts"`
	GigNotifications  bool   `json:"gig_notifications"`
	CircleUpdates     bool   `json:"circle_updates"`
	LoanReminders     bool   `json:"loan_reminders"`
	Promotions        bool   `json:"promotions"`
	SecurityAlerts    bool   `json:"security_alerts"`
	QuietHoursEnabled bool   `json:"quiet_hours_enabled"`
	QuietHoursStart   string `json:"quiet_hours_start"`
	QuietHoursEnd     string `json:"quiet_hours_end"`
	DailyDigest       bool   `json:"daily_digest"`
	WeeklyReport      bool   `json:"weekly_report"`
}

// GetDeviceTokens retrieves user's device tokens
type GetDeviceTokens struct {
	UserID string
}

// DeviceTokenDTO represents a device token for API responses
type DeviceTokenDTO struct {
	ID       string    `json:"id"`
	Platform string    `json:"platform"`
	DeviceID string    `json:"device_id"`
	IsActive bool      `json:"is_active"`
	AddedAt  time.Time `json:"added_at"`
}

// NotificationQueryHandler handles notification queries
type NotificationQueryHandler struct {
	notificationRepo repository.NotificationRepository
	preferencesRepo  repository.PreferencesRepository
	deviceTokenRepo  repository.DeviceTokenRepository
	statsRepo        repository.NotificationStatisticsRepository
}

// NewNotificationQueryHandler creates a new query handler
func NewNotificationQueryHandler(
	notificationRepo repository.NotificationRepository,
	preferencesRepo repository.PreferencesRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
	statsRepo repository.NotificationStatisticsRepository,
) *NotificationQueryHandler {
	return &NotificationQueryHandler{
		notificationRepo: notificationRepo,
		preferencesRepo:  preferencesRepo,
		deviceTokenRepo:  deviceTokenRepo,
		statsRepo:        statsRepo,
	}
}

// HandleGetNotification retrieves a single notification
func (h *NotificationQueryHandler) HandleGetNotification(ctx context.Context, q GetNotification) (*NotificationDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	notification, err := h.notificationRepo.FindByID(ctx, q.NotificationID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if notification.UserID() != userID {
		return nil, repository.ErrNotificationNotFound
	}

	return notificationToDTO(notification), nil
}

// HandleGetMyNotifications retrieves user's notifications
func (h *NotificationQueryHandler) HandleGetMyNotifications(ctx context.Context, q GetMyNotifications) (*NotificationListResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 20
	}

	filter := repository.NotificationFilter{
		IsRead: q.IsRead,
		Offset: (q.Page - 1) * q.Limit,
		Limit:  q.Limit,
	}

	if q.Type != "" {
		t := aggregate.NotificationType(q.Type)
		filter.Type = &t
	}

	if q.Channel != "" {
		c := aggregate.Channel(q.Channel)
		filter.Channel = &c
	}

	notifications, total, err := h.notificationRepo.FindByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	// Get unread count
	unreadCount, _ := h.notificationRepo.CountUnread(ctx, userID)

	dtos := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = *notificationToDTO(n)
	}

	totalPages := int(total) / q.Limit
	if int(total)%q.Limit > 0 {
		totalPages++
	}

	return &NotificationListResult{
		Notifications: dtos,
		Total:         total,
		UnreadCount:   unreadCount,
		Page:          q.Page,
		Limit:         q.Limit,
		TotalPages:    totalPages,
	}, nil
}

// HandleGetUnreadNotifications retrieves unread notifications
func (h *NotificationQueryHandler) HandleGetUnreadNotifications(ctx context.Context, q GetUnreadNotifications) ([]NotificationDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	notifications, err := h.notificationRepo.FindUnread(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = *notificationToDTO(n)
	}

	return dtos, nil
}

// HandleGetUnreadCount returns the count of unread notifications
func (h *NotificationQueryHandler) HandleGetUnreadCount(ctx context.Context, q GetUnreadCount) (*UnreadCountResult, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	count, err := h.notificationRepo.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UnreadCountResult{Count: count}, nil
}

// HandleGetNotificationPreferences retrieves user's notification preferences
func (h *NotificationQueryHandler) HandleGetNotificationPreferences(ctx context.Context, q GetNotificationPreferences) (*PreferencesDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	prefs, err := h.preferencesRepo.FindByUserID(ctx, userID)
	if err != nil {
		// Return defaults if not found
		prefs = aggregate.NewNotificationPreferences(userID)
	}

	return &PreferencesDTO{
		SMSEnabled:        prefs.SMSEnabled(),
		EmailEnabled:      prefs.EmailEnabled(),
		PushEnabled:       prefs.PushEnabled(),
		InAppEnabled:      prefs.InAppEnabled(),
		TransactionAlerts: prefs.TransactionAlerts(),
		GigNotifications:  prefs.GigNotifications(),
		CircleUpdates:     prefs.CircleUpdates(),
		LoanReminders:     prefs.LoanReminders(),
		Promotions:        prefs.Promotions(),
		SecurityAlerts:    prefs.SecurityAlerts(),
		QuietHoursEnabled: prefs.QuietHoursEnabled(),
		QuietHoursStart:   prefs.QuietHoursStart(),
		QuietHoursEnd:     prefs.QuietHoursEnd(),
		DailyDigest:       prefs.DailyDigest(),
		WeeklyReport:      prefs.WeeklyReport(),
	}, nil
}

// HandleGetDeviceTokens retrieves user's device tokens
func (h *NotificationQueryHandler) HandleGetDeviceTokens(ctx context.Context, q GetDeviceTokens) ([]DeviceTokenDTO, error) {
	userID, err := valueobject.NewUserID(q.UserID)
	if err != nil {
		return nil, err
	}

	tokens, err := h.deviceTokenRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]DeviceTokenDTO, len(tokens))
	for i, t := range tokens {
		dtos[i] = DeviceTokenDTO{
			ID:       t.ID,
			Platform: t.Platform,
			DeviceID: t.DeviceID,
			IsActive: t.IsActive,
			AddedAt:  t.CreatedAt,
		}
	}

	return dtos, nil
}

func notificationToDTO(n *aggregate.Notification) *NotificationDTO {
	return &NotificationDTO{
		ID:          n.ID(),
		Type:        n.Type().String(),
		Channel:     n.Channel().String(),
		Title:       n.Title(),
		Body:        n.Body(),
		Data:        n.Data(),
		Priority:    n.Priority().String(),
		Status:      n.Status().String(),
		IsRead:      n.IsRead(),
		SentAt:      n.SentAt(),
		DeliveredAt: n.DeliveredAt(),
		ReadAt:      n.ReadAt(),
		CreatedAt:   n.CreatedAt(),
	}
}
