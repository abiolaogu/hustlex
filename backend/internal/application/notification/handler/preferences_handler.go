package handler

import (
	"context"
	"errors"

	"hustlex/internal/application/notification/command"
	"hustlex/internal/domain/notification/aggregate"
	"hustlex/internal/domain/notification/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// PreferencesHandler handles notification preferences commands
type PreferencesHandler struct {
	preferencesRepo repository.PreferencesRepository
}

// NewPreferencesHandler creates a new preferences handler
func NewPreferencesHandler(preferencesRepo repository.PreferencesRepository) *PreferencesHandler {
	return &PreferencesHandler{
		preferencesRepo: preferencesRepo,
	}
}

// HandleUpdateNotificationPreferences updates user notification preferences
func (h *PreferencesHandler) HandleUpdateNotificationPreferences(ctx context.Context, cmd command.UpdateNotificationPreferences) (*command.PreferencesResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Get existing preferences or create new
	prefs, err := h.preferencesRepo.FindByUserID(ctx, userID)
	if err != nil {
		// Create default preferences
		prefs = aggregate.NewNotificationPreferences(userID)
	}

	// Update channel preferences if provided
	if cmd.SMSEnabled != nil || cmd.EmailEnabled != nil || cmd.PushEnabled != nil || cmd.InAppEnabled != nil {
		sms := prefs.SMSEnabled()
		email := prefs.EmailEnabled()
		push := prefs.PushEnabled()
		inApp := prefs.InAppEnabled()

		if cmd.SMSEnabled != nil {
			sms = *cmd.SMSEnabled
		}
		if cmd.EmailEnabled != nil {
			email = *cmd.EmailEnabled
		}
		if cmd.PushEnabled != nil {
			push = *cmd.PushEnabled
		}
		if cmd.InAppEnabled != nil {
			inApp = *cmd.InAppEnabled
		}

		prefs.SetChannelPreferences(sms, email, push, inApp)
	}

	// Update type preferences if provided
	if cmd.TransactionAlerts != nil || cmd.GigNotifications != nil || cmd.CircleUpdates != nil ||
		cmd.LoanReminders != nil || cmd.Promotions != nil || cmd.SecurityAlerts != nil {
		transactions := prefs.TransactionAlerts()
		gigs := prefs.GigNotifications()
		circles := prefs.CircleUpdates()
		loans := prefs.LoanReminders()
		promos := prefs.Promotions()
		security := prefs.SecurityAlerts()

		if cmd.TransactionAlerts != nil {
			transactions = *cmd.TransactionAlerts
		}
		if cmd.GigNotifications != nil {
			gigs = *cmd.GigNotifications
		}
		if cmd.CircleUpdates != nil {
			circles = *cmd.CircleUpdates
		}
		if cmd.LoanReminders != nil {
			loans = *cmd.LoanReminders
		}
		if cmd.Promotions != nil {
			promos = *cmd.Promotions
		}
		if cmd.SecurityAlerts != nil {
			security = *cmd.SecurityAlerts
		}

		prefs.SetTypePreferences(transactions, gigs, circles, loans, promos, security)
	}

	// Update quiet hours if provided
	if cmd.QuietHoursEnabled != nil || cmd.QuietHoursStart != nil || cmd.QuietHoursEnd != nil {
		enabled := prefs.QuietHoursEnabled()
		start := prefs.QuietHoursStart()
		end := prefs.QuietHoursEnd()

		if cmd.QuietHoursEnabled != nil {
			enabled = *cmd.QuietHoursEnabled
		}
		if cmd.QuietHoursStart != nil {
			start = *cmd.QuietHoursStart
		}
		if cmd.QuietHoursEnd != nil {
			end = *cmd.QuietHoursEnd
		}

		prefs.SetQuietHours(enabled, start, end)
	}

	// Update digest preferences if provided
	if cmd.DailyDigest != nil || cmd.WeeklyReport != nil {
		daily := prefs.DailyDigest()
		weekly := prefs.WeeklyReport()

		if cmd.DailyDigest != nil {
			daily = *cmd.DailyDigest
		}
		if cmd.WeeklyReport != nil {
			weekly = *cmd.WeeklyReport
		}

		prefs.SetDigestPreferences(daily, weekly)
	}

	// Save preferences
	if err := h.preferencesRepo.Save(ctx, prefs); err != nil {
		return nil, err
	}

	return &command.PreferencesResult{
		SMSEnabled:        prefs.SMSEnabled(),
		EmailEnabled:      prefs.EmailEnabled(),
		PushEnabled:       prefs.PushEnabled(),
		InAppEnabled:      prefs.InAppEnabled(),
		QuietHoursEnabled: prefs.QuietHoursEnabled(),
		QuietHoursStart:   prefs.QuietHoursStart(),
		QuietHoursEnd:     prefs.QuietHoursEnd(),
	}, nil
}

// HandleInitializePreferences initializes preferences for a new user
func (h *PreferencesHandler) HandleInitializePreferences(ctx context.Context, userIDStr string) error {
	userID, err := valueobject.NewUserID(userIDStr)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Check if already exists
	existing, _ := h.preferencesRepo.FindByUserID(ctx, userID)
	if existing != nil {
		return nil
	}

	prefs := aggregate.NewNotificationPreferences(userID)
	return h.preferencesRepo.Save(ctx, prefs)
}
