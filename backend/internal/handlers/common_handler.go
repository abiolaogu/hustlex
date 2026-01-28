package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hustlex/backend/internal/services"
)

// Package-level validator instance
var validate = validator.New()

// validateStruct validates a struct using go-playground/validator
func validateStruct(s interface{}) error {
	return validate.Struct(s)
}

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// RegisterDeviceRequest represents the request body for device registration
type RegisterDeviceRequest struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=ios android web"`
	DeviceID string `json:"device_id" validate:"required"`
}

// UpdatePreferencesRequest represents the request body for updating notification preferences
type UpdatePreferencesRequest struct {
	PushEnabled     *bool `json:"push_enabled"`
	SMSEnabled      *bool `json:"sms_enabled"`
	EmailEnabled    *bool `json:"email_enabled"`
	GigAlerts       *bool `json:"gig_alerts"`
	SavingsReminders *bool `json:"savings_reminders"`
	PaymentAlerts   *bool `json:"payment_alerts"`
	MarketingEmails *bool `json:"marketing_emails"`
}

// GetNotifications handles fetching user's notifications
// @Summary Get notifications
// @Description Get paginated list of notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by type"
// @Param read query bool false "Filter by read status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	notificationType := c.Query("type")
	readFilter := c.Query("read")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	var readPtr *bool
	if readFilter == "true" {
		t := true
		readPtr = &t
	} else if readFilter == "false" {
		f := false
		readPtr = &f
	}

	notifications, total, err := h.notificationService.GetNotifications(c.Context(), userID, notificationType, readPtr, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch notifications",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"notifications": notifications,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetUnreadCount handles fetching unread notification count
// @Summary Get unread count
// @Description Get count of unread notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	count, err := h.notificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch unread count",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"unread_count": count,
		},
	})
}

// MarkAsRead handles marking a notification as read
// @Summary Mark as read
// @Description Mark a notification as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	notificationID := c.Params("id")

	if err := h.notificationService.MarkAsRead(c.Context(), notificationID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Notification not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead handles marking all notifications as read
// @Summary Mark all as read
// @Description Mark all notifications as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /notifications/read-all [post]
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.notificationService.MarkAllAsRead(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to mark notifications as read",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "All notifications marked as read",
	})
}

// DeleteNotification handles deleting a notification
// @Summary Delete notification
// @Description Delete a notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /notifications/{id} [delete]
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	notificationID := c.Params("id")

	if err := h.notificationService.DeleteNotification(c.Context(), notificationID, userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Notification not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification deleted",
	})
}

// RegisterDevice handles device registration for push notifications
// @Summary Register device
// @Description Register a device for push notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RegisterDeviceRequest true "Device details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notifications/devices [post]
func (h *NotificationHandler) RegisterDevice(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req RegisterDeviceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := validateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if err := h.notificationService.RegisterDevice(c.Context(), userID, req.Token, req.Platform, req.DeviceID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Device registered for notifications",
	})
}

// UnregisterDevice handles device unregistration
// @Summary Unregister device
// @Description Unregister a device from push notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param device_id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Router /notifications/devices/{device_id} [delete]
func (h *NotificationHandler) UnregisterDevice(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	deviceID := c.Params("device_id")

	if err := h.notificationService.UnregisterDevice(c.Context(), userID, deviceID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Device unregistered",
	})
}

// GetPreferences handles fetching notification preferences
// @Summary Get notification preferences
// @Description Get user's notification preferences
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /notifications/preferences [get]
func (h *NotificationHandler) GetPreferences(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	prefs, err := h.notificationService.GetPreferences(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"preferences": prefs,
		},
	})
}

// UpdatePreferences handles updating notification preferences
// @Summary Update notification preferences
// @Description Update user's notification preferences
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdatePreferencesRequest true "Preference updates"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /notifications/preferences [put]
func (h *NotificationHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req UpdatePreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	updates := make(map[string]interface{})
	if req.PushEnabled != nil {
		updates["push_enabled"] = *req.PushEnabled
	}
	if req.SMSEnabled != nil {
		updates["sms_enabled"] = *req.SMSEnabled
	}
	if req.EmailEnabled != nil {
		updates["email_enabled"] = *req.EmailEnabled
	}
	if req.GigAlerts != nil {
		updates["gig_alerts"] = *req.GigAlerts
	}
	if req.SavingsReminders != nil {
		updates["savings_reminders"] = *req.SavingsReminders
	}
	if req.PaymentAlerts != nil {
		updates["payment_alerts"] = *req.PaymentAlerts
	}
	if req.MarketingEmails != nil {
		updates["marketing_emails"] = *req.MarketingEmails
	}

	prefs, err := h.notificationService.UpdatePreferences(c.Context(), userID, updates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Preferences updated",
		"data": fiber.Map{
			"preferences": prefs,
		},
	})
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	authService *services.AuthService
}

// NewUserHandler creates a new user handler
func NewUserHandler(authService *services.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// UpdateProfileRequest represents the request body for profile update
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     string `json:"email" validate:"omitempty,email"`
	Bio       string `json:"bio" validate:"omitempty,max=500"`
	Location  string `json:"location" validate:"omitempty,max=100"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

// UpdateProfile handles profile updates
// @Summary Update profile
// @Description Update user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "Profile updates"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := validateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	user, err := h.authService.UpdateProfile(c.Context(), userID, updates)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": sanitizeUser(user),
		},
	})
}

// GetPublicProfile handles fetching a user's public profile
// @Summary Get public profile
// @Description Get a user's public profile
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *UserHandler) GetPublicProfile(c *fiber.Ctx) error {
	targetUserID := c.Params("id")

	user, err := h.authService.GetUserByID(c.Context(), targetUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	// Return only public information
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": fiber.Map{
				"id":          user.ID,
				"first_name":  user.FirstName,
				"last_name":   user.LastName,
				"avatar_url":  user.AvatarURL,
				"bio":         user.Bio,
				"location":    user.Location,
				"tier":        user.Tier,
				"is_verified": user.IsVerified,
				"created_at":  user.CreatedAt,
			},
		},
	})
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles health check endpoint
// @Summary Health check
// @Description Check API health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "HustleX API is healthy",
		"version": "1.0.0",
	})
}

// ReadyCheck handles readiness check endpoint
// @Summary Readiness check
// @Description Check if API is ready to accept traffic
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) ReadyCheck(c *fiber.Ctx) error {
	// In production, this would check database connectivity, etc.
	return c.JSON(fiber.Map{
		"success": true,
		"ready":   true,
	})
}
