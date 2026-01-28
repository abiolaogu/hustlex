package handler

import (
	"context"
	"errors"

	"hustlex/internal/application/notification/command"
	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/notification/repository"
	"hustlex/internal/domain/notification/service"
	"hustlex/internal/domain/shared/valueobject"
)

// Errors
var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrUnauthorized         = errors.New("unauthorized to perform this action")
)

// NotificationHandler handles notification commands
type NotificationHandler struct {
	notificationRepo repository.NotificationRepository
	preferencesRepo  repository.PreferencesRepository
	deviceTokenRepo  repository.DeviceTokenRepository
	notificationSvc  *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(
	notificationRepo repository.NotificationRepository,
	preferencesRepo repository.PreferencesRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
	notificationSvc *service.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationRepo: notificationRepo,
		preferencesRepo:  preferencesRepo,
		deviceTokenRepo:  deviceTokenRepo,
		notificationSvc:  notificationSvc,
	}
}

// HandleSendNotification sends a notification
func (h *NotificationHandler) HandleSendNotification(ctx context.Context, cmd command.SendNotification) (*command.SendNotificationResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Check user preferences
	prefs, _ := h.preferencesRepo.FindByUserID(ctx, userID)
	if prefs != nil {
		// Check if channel and type are enabled
		channel := cmd.GetChannel()
		nType := cmd.GetType()

		if !prefs.IsChannelEnabled(channel) {
			return &command.SendNotificationResult{
				NotificationID: "",
				Status:         "skipped_channel_disabled",
			}, nil
		}

		if !prefs.IsTypeEnabled(nType) {
			return &command.SendNotificationResult{
				NotificationID: "",
				Status:         "skipped_type_disabled",
			}, nil
		}

		// Check quiet hours
		priority := cmd.GetPriority()
		if !prefs.ShouldSendNow(priority) {
			// Queue for later instead of sending now
			// For now, we'll still send but mark as delayed
		}
	}

	// Build notification
	builder := service.NewNotificationBuilder(userID).
		WithType(cmd.GetType()).
		WithChannel(cmd.GetChannel()).
		WithTitle(cmd.Title).
		WithBody(cmd.Body).
		WithPriority(cmd.GetPriority())

	for key, value := range cmd.Data {
		builder.WithData(key, value)
	}

	notification, err := builder.Build()
	if err != nil {
		return nil, err
	}

	// Send notification
	if err := h.notificationSvc.Send(ctx, notification); err != nil {
		// Save failed notification for retry
		_ = h.notificationRepo.Save(ctx, notification)
		return &command.SendNotificationResult{
			NotificationID: notification.ID(),
			Status:         notification.Status().String(),
		}, err
	}

	// Save successful notification
	if err := h.notificationRepo.Save(ctx, notification); err != nil {
		return nil, err
	}

	return &command.SendNotificationResult{
		NotificationID: notification.ID(),
		Status:         notification.Status().String(),
		ProviderID:     notification.ProviderID(),
	}, nil
}

// HandleSendOTP sends an OTP via SMS
func (h *NotificationHandler) HandleSendOTP(ctx context.Context, cmd command.SendOTP) (*command.SendOTPResult, error) {
	phone, err := cmd.GetPhone()
	if err != nil {
		return nil, errors.New("invalid phone number")
	}

	expiryMinutes := cmd.ExpiryMinutes
	if expiryMinutes <= 0 {
		expiryMinutes = 5
	}

	// Generate OTP (in real implementation, this would be done by the caller)
	// For now, we just call the provider
	messageID, err := h.notificationSvc.SendOTP(ctx, phone, "", expiryMinutes)
	if err != nil {
		return nil, err
	}

	return &command.SendOTPResult{
		MessageID: messageID,
	}, nil
}

// HandleMarkNotificationRead marks a notification as read
func (h *NotificationHandler) HandleMarkNotificationRead(ctx context.Context, cmd command.MarkNotificationRead) error {
	userID, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	notification, err := h.notificationRepo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return ErrNotificationNotFound
	}

	// Verify ownership
	if notification.UserID() != userID {
		return ErrUnauthorized
	}

	notification.MarkRead()

	return h.notificationRepo.Save(ctx, notification)
}

// HandleMarkAllNotificationsRead marks all notifications as read
func (h *NotificationHandler) HandleMarkAllNotificationsRead(ctx context.Context, cmd command.MarkAllNotificationsRead) error {
	userID, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	return h.notificationRepo.MarkAllRead(ctx, userID)
}

// HandleDeleteNotification deletes a notification
func (h *NotificationHandler) HandleDeleteNotification(ctx context.Context, cmd command.DeleteNotification) error {
	userID, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	notification, err := h.notificationRepo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return ErrNotificationNotFound
	}

	// Verify ownership
	if notification.UserID() != userID {
		return ErrUnauthorized
	}

	return h.notificationRepo.Delete(ctx, cmd.NotificationID)
}

// HandleRegisterDeviceToken registers a push notification device token
func (h *NotificationHandler) HandleRegisterDeviceToken(ctx context.Context, cmd command.RegisterDeviceToken) error {
	userID, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if token already exists
	existing, _ := h.deviceTokenRepo.FindByToken(ctx, cmd.Token)
	if existing != nil {
		// Update if different user or inactive
		if existing.UserID != userID.String() || !existing.IsActive {
			existing.UserID = userID.String()
			existing.IsActive = true
			return h.deviceTokenRepo.Save(ctx, existing)
		}
		return nil
	}

	token := &repository.DeviceToken{
		ID:       valueobject.GenerateUserID().String(),
		UserID:   userID.String(),
		Token:    cmd.Token,
		Platform: cmd.Platform,
		DeviceID: cmd.DeviceID,
		IsActive: true,
	}

	return h.deviceTokenRepo.Save(ctx, token)
}

// HandleRemoveDeviceToken removes a device token
func (h *NotificationHandler) HandleRemoveDeviceToken(ctx context.Context, cmd command.RemoveDeviceToken) error {
	_, err := cmd.GetUserID()
	if err != nil {
		return errors.New("invalid user ID")
	}

	return h.deviceTokenRepo.Delete(ctx, cmd.Token)
}

// HandleRetryFailedNotifications retries failed notifications
func (h *NotificationHandler) HandleRetryFailedNotifications(ctx context.Context, cmd command.RetryFailedNotifications) (*command.RetryResult, error) {
	limit := cmd.Limit
	if limit <= 0 {
		limit = 100
	}

	failed, err := h.notificationRepo.FindFailed(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := &command.RetryResult{
		Attempted: len(failed),
	}

	for _, notification := range failed {
		if !notification.CanRetry() {
			continue
		}

		notification.ResetForRetry()

		if err := h.notificationSvc.Send(ctx, notification); err != nil {
			result.Failed++
		} else {
			result.Succeeded++
		}

		_ = h.notificationRepo.Save(ctx, notification)
	}

	return result, nil
}
