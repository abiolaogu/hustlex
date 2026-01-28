// Package jobs implements background task processing using Asynq.
// Handles scheduled and async tasks like savings contributions, loan processing, and notifications.
package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// =============================================================================
// Task Type Constants
// =============================================================================

const (
	// Savings Tasks
	TypeSavingsContributionReminder = "savings:contribution_reminder"
	TypeSavingsProcessContribution  = "savings:process_contribution"
	TypeSavingsProcessPayout        = "savings:process_payout"
	TypeSavingsCircleMatured        = "savings:circle_matured"

	// Loan Tasks
	TypeLoanPaymentReminder   = "loan:payment_reminder"
	TypeLoanProcessRepayment  = "loan:process_repayment"
	TypeLoanCheckDefault      = "loan:check_default"
	TypeLoanUpdateCreditScore = "loan:update_credit_score"

	// Notification Tasks
	TypeNotificationPush  = "notification:push"
	TypeNotificationSMS   = "notification:sms"
	TypeNotificationEmail = "notification:email"
	TypeNotificationBulk  = "notification:bulk"

	// Gig Tasks
	TypeGigDeadlineReminder      = "gig:deadline_reminder"
	TypeGigContractAutoComplete  = "gig:contract_auto_complete"
	TypeGigEscrowRelease         = "gig:escrow_release"
	TypeGigReviewReminder        = "gig:review_reminder"

	// User Tasks
	TypeUserCreditScoreRecalc = "user:credit_score_recalc"
	TypeUserActivityReward    = "user:activity_reward"
	TypeUserInactivityCheck   = "user:inactivity_check"

	// System Tasks
	TypeSystemCleanupExpiredOTPs    = "system:cleanup_expired_otps"
	TypeSystemCleanupExpiredTokens  = "system:cleanup_expired_tokens"
	TypeSystemDailyAnalytics        = "system:daily_analytics"
	TypeSystemWeeklyReport          = "system:weekly_report"
)

// =============================================================================
// Task Payloads
// =============================================================================

// SavingsContributionReminderPayload for contribution reminders
type SavingsContributionReminderPayload struct {
	CircleID     string    `json:"circle_id"`
	UserID       string    `json:"user_id"`
	Amount       int64     `json:"amount"`
	DueDate      time.Time `json:"due_date"`
	CircleName   string    `json:"circle_name"`
	ReminderType string    `json:"reminder_type"` // "day_before", "due_day", "overdue"
}

// SavingsProcessContributionPayload for processing auto-contributions
type SavingsProcessContributionPayload struct {
	CircleID  string `json:"circle_id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
}

// SavingsProcessPayoutPayload for processing circle payouts
type SavingsProcessPayoutPayload struct {
	CircleID    string `json:"circle_id"`
	RecipientID string `json:"recipient_id"`
	Amount      int64  `json:"amount"`
	PayoutRound int    `json:"payout_round"`
}

// LoanPaymentReminderPayload for loan payment reminders
type LoanPaymentReminderPayload struct {
	LoanID       string    `json:"loan_id"`
	UserID       string    `json:"user_id"`
	Amount       int64     `json:"amount"`
	DueDate      time.Time `json:"due_date"`
	DaysUntilDue int       `json:"days_until_due"`
}

// LoanCheckDefaultPayload for checking loan defaults
type LoanCheckDefaultPayload struct {
	LoanID        string    `json:"loan_id"`
	UserID        string    `json:"user_id"`
	DueDate       time.Time `json:"due_date"`
	OutstandingBal int64    `json:"outstanding_balance"`
}

// NotificationPushPayload for push notifications
type NotificationPushPayload struct {
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
	DeviceToken string            `json:"device_token,omitempty"`
}

// NotificationSMSPayload for SMS notifications
type NotificationSMSPayload struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	SenderID    string `json:"sender_id,omitempty"`
}

// NotificationEmailPayload for email notifications
type NotificationEmailPayload struct {
	UserID      string            `json:"user_id"`
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	Template    string            `json:"template,omitempty"`
	TemplateData map[string]interface{} `json:"template_data,omitempty"`
}

// BulkNotificationPayload for sending to multiple users
type BulkNotificationPayload struct {
	UserIDs []string `json:"user_ids"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Channel string   `json:"channel"` // "push", "sms", "email", "all"
}

// GigDeadlineReminderPayload for gig deadline reminders
type GigDeadlineReminderPayload struct {
	ContractID   string    `json:"contract_id"`
	GigID        string    `json:"gig_id"`
	FreelancerID string    `json:"freelancer_id"`
	ClientID     string    `json:"client_id"`
	Deadline     time.Time `json:"deadline"`
	HoursLeft    int       `json:"hours_left"`
}

// GigEscrowReleasePayload for releasing escrow funds
type GigEscrowReleasePayload struct {
	ContractID   string `json:"contract_id"`
	FreelancerID string `json:"freelancer_id"`
	Amount       int64  `json:"amount"`
	Reference    string `json:"reference"`
}

// UserCreditScoreRecalcPayload for credit score recalculation
type UserCreditScoreRecalcPayload struct {
	UserID  string `json:"user_id"`
	Trigger string `json:"trigger"` // "loan_repayment", "savings_completion", "scheduled", etc.
}

// =============================================================================
// Task Creation Functions
// =============================================================================

// NewSavingsContributionReminderTask creates a contribution reminder task
func NewSavingsContributionReminderTask(payload SavingsContributionReminderPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSavingsContributionReminder, data, asynq.MaxRetry(3)), nil
}

// NewSavingsProcessContributionTask creates a contribution processing task
func NewSavingsProcessContributionTask(payload SavingsProcessContributionPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSavingsProcessContribution, data, asynq.MaxRetry(5), asynq.Queue("critical")), nil
}

// NewSavingsProcessPayoutTask creates a payout processing task
func NewSavingsProcessPayoutTask(payload SavingsProcessPayoutPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSavingsProcessPayout, data, asynq.MaxRetry(5), asynq.Queue("critical")), nil
}

// NewLoanPaymentReminderTask creates a loan payment reminder task
func NewLoanPaymentReminderTask(payload LoanPaymentReminderPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeLoanPaymentReminder, data, asynq.MaxRetry(3)), nil
}

// NewLoanCheckDefaultTask creates a loan default check task
func NewLoanCheckDefaultTask(payload LoanCheckDefaultPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeLoanCheckDefault, data, asynq.MaxRetry(3), asynq.Queue("critical")), nil
}

// NewNotificationPushTask creates a push notification task
func NewNotificationPushTask(payload NotificationPushPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeNotificationPush, data, asynq.MaxRetry(3)), nil
}

// NewNotificationSMSTask creates an SMS notification task
func NewNotificationSMSTask(payload NotificationSMSPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeNotificationSMS, data, asynq.MaxRetry(3)), nil
}

// NewNotificationEmailTask creates an email notification task
func NewNotificationEmailTask(payload NotificationEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeNotificationEmail, data, asynq.MaxRetry(3)), nil
}

// NewBulkNotificationTask creates a bulk notification task
func NewBulkNotificationTask(payload BulkNotificationPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeNotificationBulk, data, asynq.MaxRetry(3)), nil
}

// NewGigDeadlineReminderTask creates a gig deadline reminder task
func NewGigDeadlineReminderTask(payload GigDeadlineReminderPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeGigDeadlineReminder, data, asynq.MaxRetry(3)), nil
}

// NewGigEscrowReleaseTask creates an escrow release task
func NewGigEscrowReleaseTask(payload GigEscrowReleasePayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeGigEscrowRelease, data, asynq.MaxRetry(5), asynq.Queue("critical")), nil
}

// NewUserCreditScoreRecalcTask creates a credit score recalculation task
func NewUserCreditScoreRecalcTask(payload UserCreditScoreRecalcPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeUserCreditScoreRecalc, data, asynq.MaxRetry(3)), nil
}

// =============================================================================
// Task Handler
// =============================================================================

// TaskHandler processes background tasks
type TaskHandler struct {
	db     *gorm.DB
	client *asynq.Client
	// Add service dependencies
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(db *gorm.DB, redisAddr string) *TaskHandler {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	return &TaskHandler{
		db:     db,
		client: client,
	}
}

// Close closes the task handler
func (h *TaskHandler) Close() error {
	return h.client.Close()
}

// EnqueueTask enqueues a task for processing
func (h *TaskHandler) EnqueueTask(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return h.client.EnqueueContext(ctx, task, opts...)
}

// EnqueueTaskAt enqueues a task to be processed at a specific time
func (h *TaskHandler) EnqueueTaskAt(ctx context.Context, task *asynq.Task, processAt time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessAt(processAt))
	return h.client.EnqueueContext(ctx, task, opts...)
}

// EnqueueTaskIn enqueues a task to be processed after a delay
func (h *TaskHandler) EnqueueTaskIn(ctx context.Context, task *asynq.Task, delay time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	opts = append(opts, asynq.ProcessIn(delay))
	return h.client.EnqueueContext(ctx, task, opts...)
}

// =============================================================================
// Task Processors
// =============================================================================

// HandleSavingsContributionReminder processes contribution reminder tasks
func (h *TaskHandler) HandleSavingsContributionReminder(ctx context.Context, t *asynq.Task) error {
	var payload SavingsContributionReminderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[SAVINGS] Sending contribution reminder to user %s for circle %s (₦%.2f due on %s)",
		payload.UserID, payload.CircleID, float64(payload.Amount)/100, payload.DueDate.Format("2006-01-02"))

	// Get user's device token and send push notification
	var message string
	switch payload.ReminderType {
	case "day_before":
		message = fmt.Sprintf("Reminder: Your ₦%.2f contribution to %s is due tomorrow!", 
			float64(payload.Amount)/100, payload.CircleName)
	case "due_day":
		message = fmt.Sprintf("Your ₦%.2f contribution to %s is due today!", 
			float64(payload.Amount)/100, payload.CircleName)
	case "overdue":
		message = fmt.Sprintf("Your contribution to %s is overdue. Please contribute to avoid penalties.", 
			payload.CircleName)
	}

	// Create push notification task
	pushPayload := NotificationPushPayload{
		UserID: payload.UserID,
		Title:  "Savings Reminder",
		Body:   message,
		Data: map[string]string{
			"type":      "savings_reminder",
			"circle_id": payload.CircleID,
		},
	}

	pushTask, _ := NewNotificationPushTask(pushPayload)
	if _, err := h.EnqueueTask(ctx, pushTask); err != nil {
		log.Printf("[SAVINGS] Failed to enqueue push notification: %v", err)
	}

	return nil
}

// HandleSavingsProcessContribution processes automatic contributions
func (h *TaskHandler) HandleSavingsProcessContribution(ctx context.Context, t *asynq.Task) error {
	var payload SavingsProcessContributionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[SAVINGS] Processing auto-contribution for user %s in circle %s (₦%.2f)",
		payload.UserID, payload.CircleID, float64(payload.Amount)/100)

	// Start transaction
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Check user's wallet balance
	var wallet struct {
		ID      string
		Balance int64
	}
	if err := tx.Table("wallets").Where("user_id = ?", payload.UserID).First(&wallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("wallet not found: %w", err)
	}

	if wallet.Balance < payload.Amount {
		tx.Rollback()
		// Notify user of insufficient funds
		pushPayload := NotificationPushPayload{
			UserID: payload.UserID,
			Title:  "Auto-Save Failed",
			Body:   "Your auto-contribution couldn't be processed due to insufficient funds.",
			Data: map[string]string{
				"type":      "auto_save_failed",
				"circle_id": payload.CircleID,
			},
		}
		pushTask, _ := NewNotificationPushTask(pushPayload)
		h.EnqueueTask(ctx, pushTask)
		return fmt.Errorf("insufficient balance: have %d, need %d", wallet.Balance, payload.Amount)
	}

	// 2. Deduct from wallet
	if err := tx.Table("wallets").Where("id = ?", wallet.ID).
		Update("balance", gorm.Expr("balance - ?", payload.Amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deduct from wallet: %w", err)
	}

	// 3. Create wallet transaction
	walletTx := map[string]interface{}{
		"id":             generateUUID(),
		"wallet_id":      wallet.ID,
		"type":           "savings_contribution",
		"amount":         -payload.Amount,
		"balance_before": wallet.Balance,
		"balance_after":  wallet.Balance - payload.Amount,
		"reference":      payload.Reference,
		"status":         "completed",
		"description":    "Auto-contribution to savings circle",
		"created_at":     time.Now(),
	}
	if err := tx.Table("wallet_transactions").Create(walletTx).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// 4. Record contribution
	contribution := map[string]interface{}{
		"id":        generateUUID(),
		"circle_id": payload.CircleID,
		"user_id":   payload.UserID,
		"amount":    payload.Amount,
		"reference": payload.Reference,
		"status":    "completed",
		"paid_at":   time.Now(),
		"created_at": time.Now(),
	}
	if err := tx.Table("contributions").Create(contribution).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record contribution: %w", err)
	}

	// 5. Update circle total collected
	if err := tx.Table("savings_circles").Where("id = ?", payload.CircleID).
		Update("total_collected", gorm.Expr("total_collected + ?", payload.Amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update circle total: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 6. Notify user of successful contribution
	pushPayload := NotificationPushPayload{
		UserID: payload.UserID,
		Title:  "Auto-Save Successful",
		Body:   fmt.Sprintf("₦%.2f has been saved to your circle!", float64(payload.Amount)/100),
		Data: map[string]string{
			"type":      "auto_save_success",
			"circle_id": payload.CircleID,
		},
	}
	pushTask, _ := NewNotificationPushTask(pushPayload)
	h.EnqueueTask(ctx, pushTask)

	log.Printf("[SAVINGS] Auto-contribution completed successfully")
	return nil
}

// HandleSavingsProcessPayout processes rotational payouts
func (h *TaskHandler) HandleSavingsProcessPayout(ctx context.Context, t *asynq.Task) error {
	var payload SavingsProcessPayoutPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[SAVINGS] Processing payout for user %s from circle %s (₦%.2f)",
		payload.RecipientID, payload.CircleID, float64(payload.Amount)/100)

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get recipient's wallet
	var wallet struct {
		ID      string
		Balance int64
	}
	if err := tx.Table("wallets").Where("user_id = ?", payload.RecipientID).First(&wallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("recipient wallet not found: %w", err)
	}

	// 2. Credit recipient's wallet
	if err := tx.Table("wallets").Where("id = ?", wallet.ID).
		Update("balance", gorm.Expr("balance + ?", payload.Amount)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to credit wallet: %w", err)
	}

	// 3. Create wallet transaction
	reference := fmt.Sprintf("PAYOUT-%s-%d", payload.CircleID[:8], payload.PayoutRound)
	walletTx := map[string]interface{}{
		"id":             generateUUID(),
		"wallet_id":      wallet.ID,
		"type":           "savings_payout",
		"amount":         payload.Amount,
		"balance_before": wallet.Balance,
		"balance_after":  wallet.Balance + payload.Amount,
		"reference":      reference,
		"status":         "completed",
		"description":    fmt.Sprintf("Savings circle payout (Round %d)", payload.PayoutRound),
		"created_at":     time.Now(),
	}
	if err := tx.Table("wallet_transactions").Create(walletTx).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// 4. Update circle current round
	if err := tx.Table("savings_circles").Where("id = ?", payload.CircleID).
		Update("current_round", payload.PayoutRound).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update circle round: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 5. Notify recipient
	pushPayload := NotificationPushPayload{
		UserID: payload.RecipientID,
		Title:  "🎉 Payout Received!",
		Body:   fmt.Sprintf("You've received ₦%.2f from your savings circle!", float64(payload.Amount)/100),
		Data: map[string]string{
			"type":      "payout_received",
			"circle_id": payload.CircleID,
		},
	}
	pushTask, _ := NewNotificationPushTask(pushPayload)
	h.EnqueueTask(ctx, pushTask)

	log.Printf("[SAVINGS] Payout completed successfully")
	return nil
}

// HandleLoanPaymentReminder sends loan payment reminders
func (h *TaskHandler) HandleLoanPaymentReminder(ctx context.Context, t *asynq.Task) error {
	var payload LoanPaymentReminderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[LOAN] Sending payment reminder to user %s for loan %s (₦%.2f due in %d days)",
		payload.UserID, payload.LoanID, float64(payload.Amount)/100, payload.DaysUntilDue)

	var message string
	if payload.DaysUntilDue == 0 {
		message = fmt.Sprintf("Your loan payment of ₦%.2f is due today!", float64(payload.Amount)/100)
	} else if payload.DaysUntilDue < 0 {
		message = fmt.Sprintf("Your loan payment of ₦%.2f is %d days overdue!", 
			float64(payload.Amount)/100, -payload.DaysUntilDue)
	} else {
		message = fmt.Sprintf("Reminder: ₦%.2f loan payment due in %d days", 
			float64(payload.Amount)/100, payload.DaysUntilDue)
	}

	pushPayload := NotificationPushPayload{
		UserID: payload.UserID,
		Title:  "Loan Payment Reminder",
		Body:   message,
		Data: map[string]string{
			"type":    "loan_reminder",
			"loan_id": payload.LoanID,
		},
	}

	pushTask, _ := NewNotificationPushTask(pushPayload)
	h.EnqueueTask(ctx, pushTask)

	// Also send SMS for overdue or due today
	if payload.DaysUntilDue <= 0 {
		var user struct {
			Phone string
		}
		if err := h.db.Table("users").Select("phone").Where("id = ?", payload.UserID).First(&user).Error; err == nil {
			smsPayload := NotificationSMSPayload{
				UserID:      payload.UserID,
				PhoneNumber: user.Phone,
				Message:     message + " Pay now to avoid credit score impact.",
			}
			smsTask, _ := NewNotificationSMSTask(smsPayload)
			h.EnqueueTask(ctx, smsTask)
		}
	}

	return nil
}

// HandleLoanCheckDefault checks for defaulted loans
func (h *TaskHandler) HandleLoanCheckDefault(ctx context.Context, t *asynq.Task) error {
	var payload LoanCheckDefaultPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[LOAN] Checking default status for loan %s", payload.LoanID)

	// Calculate days overdue
	daysOverdue := int(time.Since(payload.DueDate).Hours() / 24)

	// Grace period is 3 days
	if daysOverdue <= 3 {
		log.Printf("[LOAN] Loan %s is within grace period (%d days overdue)", payload.LoanID, daysOverdue)
		return nil
	}

	// Mark as defaulted after grace period
	if err := h.db.Table("loans").Where("id = ? AND status = ?", payload.LoanID, "active").
		Updates(map[string]interface{}{
			"status":      "defaulted",
			"default_date": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("failed to update loan status: %w", err)
	}

	// Impact credit score significantly
	creditPayload := UserCreditScoreRecalcPayload{
		UserID:  payload.UserID,
		Trigger: "loan_default",
	}
	creditTask, _ := NewUserCreditScoreRecalcTask(creditPayload)
	h.EnqueueTask(ctx, creditTask)

	// Notify user
	pushPayload := NotificationPushPayload{
		UserID: payload.UserID,
		Title:  "⚠️ Loan Defaulted",
		Body:   "Your loan has been marked as defaulted. This impacts your credit score. Contact support for assistance.",
		Data: map[string]string{
			"type":    "loan_defaulted",
			"loan_id": payload.LoanID,
		},
	}
	pushTask, _ := NewNotificationPushTask(pushPayload)
	h.EnqueueTask(ctx, pushTask)

	log.Printf("[LOAN] Loan %s marked as defaulted", payload.LoanID)
	return nil
}

// HandleNotificationPush sends push notifications
func (h *TaskHandler) HandleNotificationPush(ctx context.Context, t *asynq.Task) error {
	var payload NotificationPushPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[NOTIFICATION] Sending push to user %s: %s", payload.UserID, payload.Title)

	// Get user's device tokens
	var devices []struct {
		Token    string
		Platform string
	}
	if err := h.db.Table("user_devices").
		Select("device_token as token, platform").
		Where("user_id = ? AND is_active = ?", payload.UserID, true).
		Find(&devices).Error; err != nil {
		return fmt.Errorf("failed to get device tokens: %w", err)
	}

	if len(devices) == 0 {
		log.Printf("[NOTIFICATION] No active devices for user %s", payload.UserID)
		return nil
	}

	// Send to each device (FCM handles both iOS and Android)
	for _, device := range devices {
		log.Printf("[NOTIFICATION] Sending to device: %s (%s)", device.Token[:20], device.Platform)
		// TODO: Integrate with FCM SDK
		// fcm.Send(device.Token, payload.Title, payload.Body, payload.Data)
	}

	// Save notification to database
	notification := map[string]interface{}{
		"id":         generateUUID(),
		"user_id":    payload.UserID,
		"title":      payload.Title,
		"body":       payload.Body,
		"type":       payload.Data["type"],
		"data":       payload.Data,
		"read":       false,
		"created_at": time.Now(),
	}
	h.db.Table("notifications").Create(notification)

	return nil
}

// HandleNotificationSMS sends SMS notifications
func (h *TaskHandler) HandleNotificationSMS(ctx context.Context, t *asynq.Task) error {
	var payload NotificationSMSPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[NOTIFICATION] Sending SMS to %s: %s", payload.PhoneNumber, payload.Message[:50])

	// TODO: Integrate with Termii/Africa's Talking SDK
	// sms.Send(payload.PhoneNumber, payload.Message, payload.SenderID)

	return nil
}

// HandleNotificationEmail sends email notifications
func (h *TaskHandler) HandleNotificationEmail(ctx context.Context, t *asynq.Task) error {
	var payload NotificationEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[NOTIFICATION] Sending email to %s: %s", payload.To, payload.Subject)

	// TODO: Integrate with SendGrid SDK
	// email.Send(payload.To, payload.Subject, payload.Body, payload.Template, payload.TemplateData)

	return nil
}

// HandleGigEscrowRelease releases escrow funds to freelancer
func (h *TaskHandler) HandleGigEscrowRelease(ctx context.Context, t *asynq.Task) error {
	var payload GigEscrowReleasePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[GIG] Releasing escrow for contract %s to freelancer %s (₦%.2f)",
		payload.ContractID, payload.FreelancerID, float64(payload.Amount)/100)

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get freelancer's wallet
	var wallet struct {
		ID            string
		Balance       int64
		EscrowBalance int64
	}
	if err := tx.Table("wallets").Where("user_id = ?", payload.FreelancerID).First(&wallet).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("wallet not found: %w", err)
	}

	// Calculate platform fee (10%)
	platformFee := payload.Amount / 10
	netAmount := payload.Amount - platformFee

	// 2. Credit freelancer's main balance
	if err := tx.Table("wallets").Where("id = ?", wallet.ID).
		Updates(map[string]interface{}{
			"balance":        gorm.Expr("balance + ?", netAmount),
			"escrow_balance": gorm.Expr("escrow_balance - ?", payload.Amount),
		}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	// 3. Create wallet transaction
	walletTx := map[string]interface{}{
		"id":             generateUUID(),
		"wallet_id":      wallet.ID,
		"type":           "escrow_release",
		"amount":         netAmount,
		"balance_before": wallet.Balance,
		"balance_after":  wallet.Balance + netAmount,
		"reference":      payload.Reference,
		"status":         "completed",
		"description":    fmt.Sprintf("Escrow release (Fee: ₦%.2f)", float64(platformFee)/100),
		"metadata":       map[string]interface{}{"platform_fee": platformFee, "contract_id": payload.ContractID},
		"created_at":     time.Now(),
	}
	if err := tx.Table("wallet_transactions").Create(walletTx).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// 4. Update contract status
	if err := tx.Table("contracts").Where("id = ?", payload.ContractID).
		Update("payment_released", true).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update contract: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 5. Notify freelancer
	pushPayload := NotificationPushPayload{
		UserID: payload.FreelancerID,
		Title:  "💰 Payment Received!",
		Body:   fmt.Sprintf("₦%.2f has been released to your wallet!", float64(netAmount)/100),
		Data: map[string]string{
			"type":        "payment_received",
			"contract_id": payload.ContractID,
		},
	}
	pushTask, _ := NewNotificationPushTask(pushPayload)
	h.EnqueueTask(ctx, pushTask)

	log.Printf("[GIG] Escrow released successfully")
	return nil
}

// HandleUserCreditScoreRecalc recalculates user's credit score
func (h *TaskHandler) HandleUserCreditScoreRecalc(ctx context.Context, t *asynq.Task) error {
	var payload UserCreditScoreRecalcPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Printf("[CREDIT] Recalculating credit score for user %s (trigger: %s)", payload.UserID, payload.Trigger)

	// Fetch all relevant data for scoring
	var stats struct {
		TotalLoans         int
		RepaidOnTime       int
		DefaultedLoans     int
		CompletedSavings   int
		CompletedGigs      int
		AccountAgeMonths   int
		WalletBalance      int64
		TotalEarnings      int64
	}

	// Loan history
	h.db.Table("loans").Where("user_id = ?", payload.UserID).
		Select("COUNT(*) as total_loans, SUM(CASE WHEN status = 'repaid' THEN 1 ELSE 0 END) as repaid_on_time, "+
			"SUM(CASE WHEN status = 'defaulted' THEN 1 ELSE 0 END) as defaulted_loans").
		Scan(&stats)

	// Savings participation
	h.db.Table("circle_memberships").Joins("JOIN savings_circles ON savings_circles.id = circle_memberships.circle_id").
		Where("circle_memberships.user_id = ? AND savings_circles.status = ?", payload.UserID, "completed").
		Count((*int64)(&stats.CompletedSavings))

	// Gig completion
	h.db.Table("contracts").Where("freelancer_id = ? AND status = ?", payload.UserID, "completed").
		Count((*int64)(&stats.CompletedGigs))

	// Calculate score components
	paymentHistory := calculatePaymentHistoryScore(stats.TotalLoans, stats.RepaidOnTime, stats.DefaultedLoans)
	savingsHistory := calculateSavingsScore(stats.CompletedSavings)
	gigPerformance := calculateGigScore(stats.CompletedGigs)
	// ... more components

	totalScore := (paymentHistory*35 + savingsHistory*25 + gigPerformance*20) / 100 // weighted

	// Ensure score is within bounds
	if totalScore < 300 {
		totalScore = 300
	} else if totalScore > 850 {
		totalScore = 850
	}

	// Update credit score
	h.db.Table("credit_scores").Where("user_id = ?", payload.UserID).
		Updates(map[string]interface{}{
			"total_score":     totalScore,
			"payment_history": paymentHistory,
			"savings_history": savingsHistory,
			"gig_performance": gigPerformance,
			"updated_at":      time.Now(),
		})

	// Record history
	h.db.Table("credit_score_history").Create(map[string]interface{}{
		"id":          generateUUID(),
		"user_id":     payload.UserID,
		"score":       totalScore,
		"change":      0, // Calculate from previous
		"reason":      payload.Trigger,
		"recorded_at": time.Now(),
	})

	log.Printf("[CREDIT] Credit score updated to %d", totalScore)
	return nil
}

// =============================================================================
// Worker Server
// =============================================================================

// WorkerServer manages the background job worker
type WorkerServer struct {
	server  *asynq.Server
	mux     *asynq.ServeMux
	handler *TaskHandler
}

// NewWorkerServer creates a new worker server
func NewWorkerServer(redisAddr string, db *gorm.DB, concurrency int) *WorkerServer {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("[WORKER] Task %s failed: %v", task.Type(), err)
			}),
		},
	)

	handler := NewTaskHandler(db, redisAddr)
	mux := asynq.NewServeMux()

	// Register handlers
	mux.HandleFunc(TypeSavingsContributionReminder, handler.HandleSavingsContributionReminder)
	mux.HandleFunc(TypeSavingsProcessContribution, handler.HandleSavingsProcessContribution)
	mux.HandleFunc(TypeSavingsProcessPayout, handler.HandleSavingsProcessPayout)
	mux.HandleFunc(TypeLoanPaymentReminder, handler.HandleLoanPaymentReminder)
	mux.HandleFunc(TypeLoanCheckDefault, handler.HandleLoanCheckDefault)
	mux.HandleFunc(TypeNotificationPush, handler.HandleNotificationPush)
	mux.HandleFunc(TypeNotificationSMS, handler.HandleNotificationSMS)
	mux.HandleFunc(TypeNotificationEmail, handler.HandleNotificationEmail)
	mux.HandleFunc(TypeGigEscrowRelease, handler.HandleGigEscrowRelease)
	mux.HandleFunc(TypeUserCreditScoreRecalc, handler.HandleUserCreditScoreRecalc)

	return &WorkerServer{
		server:  server,
		mux:     mux,
		handler: handler,
	}
}

// Start starts the worker server
func (w *WorkerServer) Start() error {
	log.Println("[WORKER] Starting background job worker...")
	return w.server.Start(w.mux)
}

// Shutdown gracefully shuts down the worker server
func (w *WorkerServer) Shutdown() {
	log.Println("[WORKER] Shutting down background job worker...")
	w.server.Shutdown()
	w.handler.Close()
}

// =============================================================================
// Scheduler
// =============================================================================

// Scheduler manages scheduled/cron jobs
type Scheduler struct {
	scheduler *asynq.Scheduler
	handler   *TaskHandler
}

// NewScheduler creates a new scheduler
func NewScheduler(redisAddr string, db *gorm.DB, location *time.Location) *Scheduler {
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: redisAddr},
		&asynq.SchedulerOpts{
			Location: location,
		},
	)

	handler := NewTaskHandler(db, redisAddr)

	return &Scheduler{
		scheduler: scheduler,
		handler:   handler,
	}
}

// RegisterScheduledJobs registers all scheduled/cron jobs
func (s *Scheduler) RegisterScheduledJobs() error {
	// Daily contribution reminders at 9 AM WAT
	if _, err := s.scheduler.Register("0 9 * * *", asynq.NewTask(
		"scheduler:daily_contribution_reminders", nil,
	)); err != nil {
		return fmt.Errorf("failed to register contribution reminders: %w", err)
	}

	// Daily loan payment reminders at 8 AM WAT
	if _, err := s.scheduler.Register("0 8 * * *", asynq.NewTask(
		"scheduler:daily_loan_reminders", nil,
	)); err != nil {
		return fmt.Errorf("failed to register loan reminders: %w", err)
	}

	// Check for loan defaults at midnight
	if _, err := s.scheduler.Register("0 0 * * *", asynq.NewTask(
		"scheduler:check_loan_defaults", nil,
	)); err != nil {
		return fmt.Errorf("failed to register loan default check: %w", err)
	}

	// Weekly credit score recalculation on Sundays at midnight
	if _, err := s.scheduler.Register("0 0 * * 0", asynq.NewTask(
		"scheduler:weekly_credit_recalc", nil,
	)); err != nil {
		return fmt.Errorf("failed to register credit recalc: %w", err)
	}

	// Cleanup expired OTPs every hour
	if _, err := s.scheduler.Register("0 * * * *", asynq.NewTask(
		TypeSystemCleanupExpiredOTPs, nil,
	)); err != nil {
		return fmt.Errorf("failed to register OTP cleanup: %w", err)
	}

	log.Println("[SCHEDULER] Registered all scheduled jobs")
	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start() error {
	log.Println("[SCHEDULER] Starting job scheduler...")
	return s.scheduler.Start()
}

// Shutdown gracefully shuts down the scheduler
func (s *Scheduler) Shutdown() {
	log.Println("[SCHEDULER] Shutting down job scheduler...")
	s.scheduler.Shutdown()
	s.handler.Close()
}

// =============================================================================
// Helper Functions
// =============================================================================

func generateUUID() string {
	// Simple UUID generation - in production use google/uuid
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func calculatePaymentHistoryScore(total, onTime, defaulted int) int {
	if total == 0 {
		return 650 // Neutral score for no history
	}
	
	// Base score
	score := 500
	
	// Add points for on-time payments
	score += (onTime * 50)
	if score > 850 {
		score = 850
	}
	
	// Subtract heavily for defaults
	score -= (defaulted * 150)
	if score < 300 {
		score = 300
	}
	
	return score
}

func calculateSavingsScore(completedCircles int) int {
	// Each completed savings circle adds to score
	score := 500 + (completedCircles * 30)
	if score > 850 {
		score = 850
	}
	return score
}

func calculateGigScore(completedGigs int) int {
	// Each completed gig contributes to score
	score := 500 + (completedGigs * 20)
	if score > 850 {
		score = 850
	}
	return score
}
