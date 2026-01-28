package aggregate

import (
	"time"

	"hustlex/internal/domain/shared/valueobject"
)

// NotificationPreferences holds user notification settings
type NotificationPreferences struct {
	id        string
	userID    valueobject.UserID

	// Channel preferences
	smsEnabled   bool
	emailEnabled bool
	pushEnabled  bool
	inAppEnabled bool

	// Type preferences
	transactionAlerts   bool
	gigNotifications    bool
	circleUpdates       bool
	loanReminders       bool
	promotions          bool
	securityAlerts      bool

	// Quiet hours
	quietHoursEnabled   bool
	quietHoursStart     string // e.g., "22:00"
	quietHoursEnd       string // e.g., "07:00"

	// Frequency settings
	dailyDigest         bool
	weeklyReport        bool

	createdAt time.Time
	updatedAt time.Time
}

// NewNotificationPreferences creates default preferences for a user
func NewNotificationPreferences(userID valueobject.UserID) *NotificationPreferences {
	return &NotificationPreferences{
		id:                valueobject.GenerateUserID().String(),
		userID:            userID,
		smsEnabled:        true,
		emailEnabled:      true,
		pushEnabled:       true,
		inAppEnabled:      true,
		transactionAlerts: true,
		gigNotifications:  true,
		circleUpdates:     true,
		loanReminders:     true,
		promotions:        false, // Default off for marketing
		securityAlerts:    true,
		quietHoursEnabled: false,
		quietHoursStart:   "22:00",
		quietHoursEnd:     "07:00",
		dailyDigest:       false,
		weeklyReport:      false,
		createdAt:         time.Now().UTC(),
		updatedAt:         time.Now().UTC(),
	}
}

// ReconstructPreferences reconstructs from persistence
func ReconstructPreferences(
	id string,
	userID valueobject.UserID,
	smsEnabled bool,
	emailEnabled bool,
	pushEnabled bool,
	inAppEnabled bool,
	transactionAlerts bool,
	gigNotifications bool,
	circleUpdates bool,
	loanReminders bool,
	promotions bool,
	securityAlerts bool,
	quietHoursEnabled bool,
	quietHoursStart string,
	quietHoursEnd string,
	dailyDigest bool,
	weeklyReport bool,
	createdAt time.Time,
	updatedAt time.Time,
) *NotificationPreferences {
	return &NotificationPreferences{
		id:                id,
		userID:            userID,
		smsEnabled:        smsEnabled,
		emailEnabled:      emailEnabled,
		pushEnabled:       pushEnabled,
		inAppEnabled:      inAppEnabled,
		transactionAlerts: transactionAlerts,
		gigNotifications:  gigNotifications,
		circleUpdates:     circleUpdates,
		loanReminders:     loanReminders,
		promotions:        promotions,
		securityAlerts:    securityAlerts,
		quietHoursEnabled: quietHoursEnabled,
		quietHoursStart:   quietHoursStart,
		quietHoursEnd:     quietHoursEnd,
		dailyDigest:       dailyDigest,
		weeklyReport:      weeklyReport,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

// Getters
func (p *NotificationPreferences) ID() string                 { return p.id }
func (p *NotificationPreferences) UserID() valueobject.UserID { return p.userID }
func (p *NotificationPreferences) SMSEnabled() bool           { return p.smsEnabled }
func (p *NotificationPreferences) EmailEnabled() bool         { return p.emailEnabled }
func (p *NotificationPreferences) PushEnabled() bool          { return p.pushEnabled }
func (p *NotificationPreferences) InAppEnabled() bool         { return p.inAppEnabled }
func (p *NotificationPreferences) TransactionAlerts() bool    { return p.transactionAlerts }
func (p *NotificationPreferences) GigNotifications() bool     { return p.gigNotifications }
func (p *NotificationPreferences) CircleUpdates() bool        { return p.circleUpdates }
func (p *NotificationPreferences) LoanReminders() bool        { return p.loanReminders }
func (p *NotificationPreferences) Promotions() bool           { return p.promotions }
func (p *NotificationPreferences) SecurityAlerts() bool       { return p.securityAlerts }
func (p *NotificationPreferences) QuietHoursEnabled() bool    { return p.quietHoursEnabled }
func (p *NotificationPreferences) QuietHoursStart() string    { return p.quietHoursStart }
func (p *NotificationPreferences) QuietHoursEnd() string      { return p.quietHoursEnd }
func (p *NotificationPreferences) DailyDigest() bool          { return p.dailyDigest }
func (p *NotificationPreferences) WeeklyReport() bool         { return p.weeklyReport }
func (p *NotificationPreferences) CreatedAt() time.Time       { return p.createdAt }
func (p *NotificationPreferences) UpdatedAt() time.Time       { return p.updatedAt }

// Business Methods

// SetChannelPreferences updates channel preferences
func (p *NotificationPreferences) SetChannelPreferences(sms, email, push, inApp bool) {
	p.smsEnabled = sms
	p.emailEnabled = email
	p.pushEnabled = push
	p.inAppEnabled = inApp
	p.updatedAt = time.Now().UTC()
}

// SetTypePreferences updates notification type preferences
func (p *NotificationPreferences) SetTypePreferences(transactions, gigs, circles, loans, promos, security bool) {
	p.transactionAlerts = transactions
	p.gigNotifications = gigs
	p.circleUpdates = circles
	p.loanReminders = loans
	p.promotions = promos
	p.securityAlerts = security
	p.updatedAt = time.Now().UTC()
}

// SetQuietHours sets quiet hours
func (p *NotificationPreferences) SetQuietHours(enabled bool, start, end string) {
	p.quietHoursEnabled = enabled
	p.quietHoursStart = start
	p.quietHoursEnd = end
	p.updatedAt = time.Now().UTC()
}

// SetDigestPreferences sets digest preferences
func (p *NotificationPreferences) SetDigestPreferences(daily, weekly bool) {
	p.dailyDigest = daily
	p.weeklyReport = weekly
	p.updatedAt = time.Now().UTC()
}

// IsChannelEnabled checks if a specific channel is enabled
func (p *NotificationPreferences) IsChannelEnabled(channel Channel) bool {
	switch channel {
	case ChannelSMS:
		return p.smsEnabled
	case ChannelEmail:
		return p.emailEnabled
	case ChannelPush:
		return p.pushEnabled
	case ChannelInApp:
		return p.inAppEnabled
	default:
		return false
	}
}

// IsTypeEnabled checks if a notification type is enabled
func (p *NotificationPreferences) IsTypeEnabled(nType NotificationType) bool {
	switch nType {
	case TypeTransaction, TypePaymentReceived, TypePaymentSent:
		return p.transactionAlerts
	case TypeGigUpdate, TypeContractUpdate:
		return p.gigNotifications
	case TypeCircleUpdate, TypeContribution, TypePayout:
		return p.circleUpdates
	case TypeLoanUpdate, TypeCreditUpdate:
		return p.loanReminders
	case TypePromotion:
		return p.promotions
	case TypeOTP, TypeSystem:
		return true // Always enabled for security
	case TypeReminder:
		return p.securityAlerts
	default:
		return true
	}
}

// ShouldSendNow checks if a notification should be sent now based on quiet hours
func (p *NotificationPreferences) ShouldSendNow(priority Priority) bool {
	// Urgent notifications always go through
	if priority == PriorityUrgent || priority == PriorityHigh {
		return true
	}

	if !p.quietHoursEnabled {
		return true
	}

	// Parse quiet hours and check current time
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	currentTime := hour*60 + minute

	// Parse start time
	var startHour, startMin int
	_, _ = parseTime(p.quietHoursStart, &startHour, &startMin)
	startTime := startHour*60 + startMin

	// Parse end time
	var endHour, endMin int
	_, _ = parseTime(p.quietHoursEnd, &endHour, &endMin)
	endTime := endHour*60 + endMin

	// Check if current time is in quiet hours
	if startTime < endTime {
		// Same day quiet hours (e.g., 14:00 - 16:00)
		if currentTime >= startTime && currentTime < endTime {
			return false
		}
	} else {
		// Overnight quiet hours (e.g., 22:00 - 07:00)
		if currentTime >= startTime || currentTime < endTime {
			return false
		}
	}

	return true
}

func parseTime(timeStr string, hour, min *int) (bool, error) {
	if len(timeStr) != 5 {
		return false, nil
	}
	// Simple parsing for HH:MM format
	*hour = int(timeStr[0]-'0')*10 + int(timeStr[1]-'0')
	*min = int(timeStr[3]-'0')*10 + int(timeStr[4]-'0')
	return true, nil
}
