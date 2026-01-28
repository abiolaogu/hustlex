package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"hustlex/internal/models"
)

// NotificationService handles all notification operations
type NotificationService struct {
	db          *gorm.DB
	smsGateway  SMSGateway
	pushGateway PushGateway
}

// SMSGateway interface for SMS providers (Termii, Africa's Talking)
type SMSGateway interface {
	SendSMS(to, message string) error
	SendOTP(to, code string) error
}

// PushGateway interface for push notifications (Firebase)
type PushGateway interface {
	SendPush(token string, title, body string, data map[string]string) error
	SendToTopic(topic string, title, body string, data map[string]string) error
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB, sms SMSGateway, push PushGateway) *NotificationService {
	return &NotificationService{
		db:          db,
		smsGateway:  sms,
		pushGateway: push,
	}
}

// Notification errors
var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidChannel       = errors.New("invalid notification channel")
)

// NotificationChannel represents delivery channel
type NotificationChannel string

const (
	ChannelInApp NotificationChannel = "in_app"
	ChannelPush  NotificationChannel = "push"
	ChannelSMS   NotificationChannel = "sms"
	ChannelEmail NotificationChannel = "email"
)

// NotificationType defines notification categories
type NotificationType string

const (
	// Gig notifications
	NotifyGigNew           NotificationType = "gig_new"
	NotifyGigProposal      NotificationType = "gig_proposal"
	NotifyGigAccepted      NotificationType = "gig_accepted"
	NotifyGigDelivery      NotificationType = "gig_delivery"
	NotifyGigCompleted     NotificationType = "gig_completed"
	NotifyGigReview        NotificationType = "gig_review"
	NotifyGigCancelled     NotificationType = "gig_cancelled"

	// Savings notifications
	NotifyCircleInvite      NotificationType = "circle_invite"
	NotifyCircleJoined      NotificationType = "circle_joined"
	NotifyCircleStarted     NotificationType = "circle_started"
	NotifyContributionDue   NotificationType = "contribution_due"
	NotifyContributionPaid  NotificationType = "contribution_paid"
	NotifyContributionMissed NotificationType = "contribution_missed"
	NotifyPayoutReceived    NotificationType = "payout_received"
	NotifyCircleCompleted   NotificationType = "circle_completed"

	// Wallet notifications
	NotifyDeposit           NotificationType = "deposit"
	NotifyWithdrawal        NotificationType = "withdrawal"
	NotifyTransferSent      NotificationType = "transfer_sent"
	NotifyTransferReceived  NotificationType = "transfer_received"
	NotifyPaymentReceived   NotificationType = "payment_received"

	// Credit notifications
	NotifyCreditScoreUpdate NotificationType = "credit_score_update"
	NotifyTierUpgrade       NotificationType = "tier_upgrade"
	NotifyLoanApproved      NotificationType = "loan_approved"
	NotifyLoanDue           NotificationType = "loan_due"
	NotifyLoanOverdue       NotificationType = "loan_overdue"

	// Account notifications
	NotifyWelcome           NotificationType = "welcome"
	NotifyVerification      NotificationType = "verification"
	NotifySecurityAlert     NotificationType = "security_alert"

	// Marketing/Engagement
	NotifyPromotion         NotificationType = "promotion"
	NotifyReminder          NotificationType = "reminder"
)

// SendNotificationInput represents notification request
type SendNotificationInput struct {
	UserID      uuid.UUID
	Type        NotificationType
	Title       string
	Message     string
	Data        map[string]interface{}
	Channels    []NotificationChannel
	ScheduledAt *time.Time
}

// SendNotification creates and sends a notification
func (s *NotificationService) SendNotification(ctx context.Context, input SendNotificationInput) (*models.Notification, error) {
	// Create notification record
	notification := &models.Notification{
		UserID:  input.UserID,
		Type:    string(input.Type),
		Title:   input.Title,
		Body:    input.Message,
		IsRead:  false,
	}

	if err := s.db.WithContext(ctx).Create(notification).Error; err != nil {
		return nil, err
	}

	// Get user for delivery
	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ?", input.UserID).First(&user).Error; err != nil {
		return notification, nil // Still return notification even if user not found
	}

	// Send via requested channels
	for _, channel := range input.Channels {
		switch channel {
		case ChannelPush:
			// Push notifications disabled - no FCM token on User model
			// TODO: Implement push notification device token storage
		case ChannelSMS:
			if user.Phone != "" && s.smsGateway != nil {
				go s.smsGateway.SendSMS(user.Phone, input.Message)
			}
		}
	}

	return notification, nil
}

// SendBulkNotification sends notification to multiple users
func (s *NotificationService) SendBulkNotification(ctx context.Context, userIDs []uuid.UUID, notifType NotificationType, title, message string, data map[string]interface{}) error {
	for _, userID := range userIDs {
		_, err := s.SendNotification(ctx, SendNotificationInput{
			UserID:   userID,
			Type:     notifType,
			Title:    title,
			Message:  message,
			Data:     data,
			Channels: []NotificationChannel{ChannelInApp, ChannelPush},
		})
		if err != nil {
			// Log error but continue with others
			continue
		}
	}
	return nil
}

// GetUserNotifications retrieves notifications for a user
func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit, offset int) ([]models.Notification, int64, error) {
	query := s.db.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit == 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var notifications []models.Notification
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	return notifications, total, err
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

// MarkAllAsRead marks all user notifications as read
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

// GetUnreadCount returns count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error

	return count, err
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

// DeleteOldNotifications removes notifications older than given duration
func (s *NotificationService) DeleteOldNotifications(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan)

	result := s.db.WithContext(ctx).
		Where("created_at < ? AND is_read = ?", cutoff, true).
		Delete(&models.Notification{})

	return result.RowsAffected, result.Error
}

// UpdateFCMToken updates user's push notification token
func (s *NotificationService) UpdateFCMToken(ctx context.Context, userID uuid.UUID, token string) error {
	return s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("fcm_token", token).Error
}

// === Convenience notification methods ===

// NotifyGigProposalReceived notifies client about new proposal
func (s *NotificationService) NotifyGigProposalReceived(ctx context.Context, clientID uuid.UUID, gigTitle, hustlerName string, proposalID uuid.UUID) error {
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   clientID,
		Type:     NotifyGigProposal,
		Title:    "New Proposal Received",
		Message:  fmt.Sprintf("%s submitted a proposal for \"%s\"", hustlerName, gigTitle),
		Data:     map[string]interface{}{"proposal_id": proposalID.String()},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush},
	})
	return err
}

// NotifyProposalAccepted notifies hustler that proposal was accepted
func (s *NotificationService) NotifyProposalAccepted(ctx context.Context, hustlerID uuid.UUID, gigTitle string, contractID uuid.UUID) error {
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   hustlerID,
		Type:     NotifyGigAccepted,
		Title:    "Proposal Accepted! 🎉",
		Message:  fmt.Sprintf("Your proposal for \"%s\" has been accepted. Time to get to work!", gigTitle),
		Data:     map[string]interface{}{"contract_id": contractID.String()},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush, ChannelSMS},
	})
	return err
}

// NotifyWorkDelivered notifies client that work has been delivered
func (s *NotificationService) NotifyWorkDelivered(ctx context.Context, clientID uuid.UUID, gigTitle string, contractID uuid.UUID) error {
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   clientID,
		Type:     NotifyGigDelivery,
		Title:    "Work Delivered",
		Message:  fmt.Sprintf("The work for \"%s\" has been delivered. Please review and approve.", gigTitle),
		Data:     map[string]interface{}{"contract_id": contractID.String()},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush},
	})
	return err
}

// NotifyPaymentReceived notifies hustler about payment
func (s *NotificationService) NotifyPaymentReceived(ctx context.Context, hustlerID uuid.UUID, amountKobo int64, gigTitle string) error {
	amountNaira := float64(amountKobo) / 100
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   hustlerID,
		Type:     NotifyPaymentReceived,
		Title:    "Payment Received 💰",
		Message:  fmt.Sprintf("You received ₦%.2f for \"%s\"", amountNaira, gigTitle),
		Data:     map[string]interface{}{"amount": amountKobo},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush, ChannelSMS},
	})
	return err
}

// NotifyContributionReminder reminds user about upcoming contribution
func (s *NotificationService) NotifyContributionReminder(ctx context.Context, userID uuid.UUID, circleName string, amountKobo int64, dueDate time.Time) error {
	amountNaira := float64(amountKobo) / 100
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyContributionDue,
		Title:    "Contribution Due Soon",
		Message:  fmt.Sprintf("Your ₦%.2f contribution to \"%s\" is due on %s", amountNaira, circleName, dueDate.Format("Jan 2")),
		Data:     map[string]interface{}{"circle_name": circleName, "amount": amountKobo},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush},
	})
	return err
}

// NotifyCirclePayoutReceived notifies user about Ajo payout
func (s *NotificationService) NotifyCirclePayoutReceived(ctx context.Context, userID uuid.UUID, circleName string, amountKobo int64) error {
	amountNaira := float64(amountKobo) / 100
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyPayoutReceived,
		Title:    "Payout Received! 🎉",
		Message:  fmt.Sprintf("You received your ₦%.2f payout from \"%s\"", amountNaira, circleName),
		Data:     map[string]interface{}{"circle_name": circleName, "amount": amountKobo},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush, ChannelSMS},
	})
	return err
}

// NotifyTierUpgraded notifies user about tier upgrade
func (s *NotificationService) NotifyTierUpgraded(ctx context.Context, userID uuid.UUID, newTier string, creditScore int) error {
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyTierUpgrade,
		Title:    "Tier Upgraded! 🚀",
		Message:  fmt.Sprintf("Congratulations! You've been upgraded to %s tier with a credit score of %d", newTier, creditScore),
		Data:     map[string]interface{}{"tier": newTier, "credit_score": creditScore},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush},
	})
	return err
}

// NotifyLoanApproved notifies user about loan approval
func (s *NotificationService) NotifyLoanApproved(ctx context.Context, userID uuid.UUID, amountKobo int64, loanID uuid.UUID) error {
	amountNaira := float64(amountKobo) / 100
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyLoanApproved,
		Title:    "Loan Approved! 💸",
		Message:  fmt.Sprintf("Your loan of ₦%.2f has been approved and credited to your wallet", amountNaira),
		Data:     map[string]interface{}{"loan_id": loanID.String(), "amount": amountKobo},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush, ChannelSMS},
	})
	return err
}

// NotifyTransferReceived notifies user about incoming transfer
func (s *NotificationService) NotifyTransferReceived(ctx context.Context, userID uuid.UUID, amountKobo int64, senderName string) error {
	amountNaira := float64(amountKobo) / 100
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyTransferReceived,
		Title:    "Money Received 💵",
		Message:  fmt.Sprintf("You received ₦%.2f from %s", amountNaira, senderName),
		Data:     map[string]interface{}{"amount": amountKobo, "sender": senderName},
		Channels: []NotificationChannel{ChannelInApp, ChannelPush},
	})
	return err
}

// NotifyWelcome sends welcome notification to new users
func (s *NotificationService) NotifyWelcome(ctx context.Context, userID uuid.UUID, firstName string) error {
	_, err := s.SendNotification(ctx, SendNotificationInput{
		UserID:   userID,
		Type:     NotifyWelcome,
		Title:    "Welcome to HustleX! 🎉",
		Message:  fmt.Sprintf("Hey %s! Start your hustle journey by completing your profile and posting your first gig.", firstName),
		Channels: []NotificationChannel{ChannelInApp},
	})
	return err
}

// === Stub implementations for gateways ===

// TermiiSMSGateway implements SMSGateway using Termii
type TermiiSMSGateway struct {
	apiKey    string
	senderID  string
	baseURL   string
}

// NewTermiiSMSGateway creates a Termii SMS gateway
func NewTermiiSMSGateway(apiKey, senderID string) *TermiiSMSGateway {
	return &TermiiSMSGateway{
		apiKey:   apiKey,
		senderID: senderID,
		baseURL:  "https://api.ng.termii.com/api",
	}
}

// SendSMS sends an SMS via Termii
func (g *TermiiSMSGateway) SendSMS(to, message string) error {
	// TODO: Implement actual Termii API call
	// POST to /sms/send with payload:
	// {
	//   "to": to,
	//   "from": g.senderID,
	//   "sms": message,
	//   "type": "plain",
	//   "channel": "generic",
	//   "api_key": g.apiKey
	// }
	fmt.Printf("[SMS] To: %s, Message: %s\n", to, message)
	return nil
}

// SendOTP sends an OTP via Termii
func (g *TermiiSMSGateway) SendOTP(to, code string) error {
	message := fmt.Sprintf("Your HustleX verification code is %s. Valid for 10 minutes.", code)
	return g.SendSMS(to, message)
}

// FirebasePushGateway implements PushGateway using Firebase
type FirebasePushGateway struct {
	projectID string
	// credentials would go here
}

// NewFirebasePushGateway creates a Firebase push gateway
func NewFirebasePushGateway(projectID string) *FirebasePushGateway {
	return &FirebasePushGateway{
		projectID: projectID,
	}
}

// SendPush sends a push notification via Firebase
func (g *FirebasePushGateway) SendPush(token string, title, body string, data map[string]string) error {
	// TODO: Implement actual Firebase Cloud Messaging call
	fmt.Printf("[PUSH] Token: %s..., Title: %s, Body: %s\n", token[:20], title, body)
	return nil
}

// SendToTopic sends a push notification to a topic
func (g *FirebasePushGateway) SendToTopic(topic string, title, body string, data map[string]string) error {
	// TODO: Implement actual Firebase Cloud Messaging call
	fmt.Printf("[PUSH-TOPIC] Topic: %s, Title: %s, Body: %s\n", topic, title, body)
	return nil
}
