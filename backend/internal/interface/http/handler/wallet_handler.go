package handler

import (
	"encoding/json"
	"net/http"

	"hustlex/internal/application/wallet/command"
	"hustlex/internal/application/wallet/handler"
	"hustlex/internal/application/wallet/query"
	"hustlex/internal/infrastructure/security/audit"
	"hustlex/internal/infrastructure/security/validation"
	"hustlex/internal/interface/http/middleware"
	"hustlex/internal/interface/http/response"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	depositHandler  *handler.DepositHandler
	withdrawHandler *handler.WithdrawHandler
	transferHandler *handler.TransferHandler
	queryHandler    *query.WalletQueryHandler
	auditLogger     audit.AuditLogger
}

// NewWalletHandler creates a new wallet HTTP handler
func NewWalletHandler(
	depositHandler *handler.DepositHandler,
	withdrawHandler *handler.WithdrawHandler,
	transferHandler *handler.TransferHandler,
	queryHandler *query.WalletQueryHandler,
	auditLogger audit.AuditLogger,
) *WalletHandler {
	return &WalletHandler{
		depositHandler:  depositHandler,
		withdrawHandler: withdrawHandler,
		transferHandler: transferHandler,
		queryHandler:    queryHandler,
		auditLogger:     auditLogger,
	}
}

// GetWallet handles GET /api/wallet
func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	wallet, err := h.queryHandler.HandleGetWallet(r.Context(), query.GetWallet{
		UserID: userID.String(),
	})
	if err != nil {
		response.NotFound(w, "wallet not found")
		return
	}

	response.Success(w, wallet)
}

// GetBalance handles GET /api/wallet/balance
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	// Use GetWallet as there's no separate balance endpoint
	wallet, err := h.queryHandler.HandleGetWallet(r.Context(), query.GetWallet{
		UserID: userID.String(),
	})
	if err != nil {
		response.NotFound(w, "wallet not found")
		return
	}

	response.Success(w, map[string]interface{}{
		"available_balance": wallet.AvailableBalance,
		"escrow_balance":    wallet.EscrowBalance,
		"savings_balance":   wallet.SavingsBalance,
		"total_balance":     wallet.TotalBalance,
		"currency":          wallet.Currency,
	})
}

// GetTransactions handles GET /api/wallet/transactions
func (h *WalletHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	// Parse query parameters
	q := r.URL.Query()
	page := parseIntQuery(q.Get("page"), 1)
	limit := parseIntQuery(q.Get("limit"), 20)

	// Convert string to pointer if not empty
	var typeFilter *string
	var statusFilter *string
	if t := q.Get("type"); t != "" {
		typeFilter = &t
	}
	if s := q.Get("status"); s != "" {
		statusFilter = &s
	}

	result, err := h.queryHandler.HandleGetTransactions(r.Context(), query.GetTransactions{
		UserID: userID.String(),
		Type:   typeFilter,
		Status: statusFilter,
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Paginated(w, result.Transactions, result.Page, result.Limit, result.Total)
}

// InitiateDeposit handles POST /api/wallet/deposit/initiate
func (h *WalletHandler) InitiateDeposit(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	var req struct {
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Reference   string `json:"reference"`
		Description string `json:"description"`
		Channel     string `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Input validation
	v := validation.NewValidator()
	v.Positive("amount", req.Amount).
		Min("amount", req.Amount, 100).           // Minimum 100 kobo (1 Naira)
		Max("amount", req.Amount, 100000000).     // Maximum 1M Naira
		OneOf("currency", req.Currency, []string{"NGN", "USD", ""}).
		SafeString("description", req.Description).
		SafeString("reference", req.Reference)

	if v.HasErrors() {
		response.ValidationError(w, v.Errors().Errors)
		return
	}

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	// Get wallet first to get wallet ID
	wallet, err := h.queryHandler.HandleGetWallet(r.Context(), query.GetWallet{
		UserID: userID.String(),
	})
	if err != nil {
		response.NotFound(w, "wallet not found")
		return
	}

	result, err := h.depositHandler.HandleInitiateDeposit(r.Context(), command.Deposit{
		WalletID:    wallet.WalletID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Source:      "deposit",
		Reference:   req.Reference,
		Description: req.Description,
		Channel:     req.Channel,
		RequestedBy: userID.String(),
	})

	// Audit log the deposit attempt
	if h.auditLogger != nil {
		outcome := audit.OutcomeSuccess
		message := "Deposit initiated successfully"
		if err != nil {
			outcome = audit.OutcomeFailure
			message = "Deposit initiation failed"
		}
		h.auditLogger.LogTransaction(r.Context(), audit.AuditEvent{
			EventAction:    audit.ActionCreate,
			EventOutcome:   outcome,
			ActorUserID:    userID.String(),
			ActorIPAddress: getClientIP(r),
			ActorUserAgent: r.UserAgent(),
			TargetType:     "wallet",
			TargetID:       wallet.WalletID,
			Message:        message,
			Component:      "wallet_handler",
			Metadata: map[string]interface{}{
				"amount":    req.Amount,
				"currency":  req.Currency,
				"reference": req.Reference,
				"channel":   req.Channel,
			},
		})
	}

	if err != nil {
		response.BadRequest(w, "deposit initiation failed")
		return
	}

	response.Success(w, result)
}

// VerifyDeposit handles POST /api/wallet/deposit/verify
func (h *WalletHandler) VerifyDeposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Reference string `json:"reference"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Reference == "" {
		response.ValidationError(w, map[string]string{"reference": "reference is required"})
		return
	}

	result, err := h.depositHandler.HandleVerifyDeposit(r.Context(), req.Reference)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, result)
}

// InitiateWithdraw handles POST /api/wallet/withdraw
func (h *WalletHandler) InitiateWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	var req struct {
		Amount        int64  `json:"amount"`
		Currency      string `json:"currency"`
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
		AccountName   string `json:"account_name"`
		PIN           string `json:"pin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Input validation with security checks
	v := validation.NewValidator()
	v.Required("account_number", req.AccountNumber).
		Required("bank_code", req.BankCode).
		Required("pin", req.PIN).
		Positive("amount", req.Amount).
		Min("amount", req.Amount, 100).           // Minimum 100 kobo
		Max("amount", req.Amount, 100000000).     // Maximum 1M Naira
		AccountNumber("account_number", req.AccountNumber).
		BankCode("bank_code", req.BankCode).
		PIN("pin", req.PIN).
		SafeString("account_name", req.AccountName)

	if v.HasErrors() {
		response.ValidationError(w, v.Errors().Errors)
		return
	}

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	// Get wallet first to get wallet ID
	wallet, err := h.queryHandler.HandleGetWallet(r.Context(), query.GetWallet{
		UserID: userID.String(),
	})
	if err != nil {
		response.NotFound(w, "wallet not found")
		return
	}

	result, err := h.withdrawHandler.Handle(r.Context(), command.Withdraw{
		WalletID:      wallet.WalletID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		AccountNumber: req.AccountNumber,
		BankCode:      req.BankCode,
		AccountName:   req.AccountName,
		PIN:           req.PIN,
		RequestedBy:   userID.String(),
	})

	// Audit log the withdrawal attempt (sensitive operation)
	if h.auditLogger != nil {
		outcome := audit.OutcomeSuccess
		message := "Withdrawal initiated successfully"
		if err != nil {
			outcome = audit.OutcomeFailure
			message = "Withdrawal initiation failed"
		}
		h.auditLogger.LogTransaction(r.Context(), audit.AuditEvent{
			EventAction:    audit.ActionCreate,
			EventOutcome:   outcome,
			ActorUserID:    userID.String(),
			ActorIPAddress: getClientIP(r),
			ActorUserAgent: r.UserAgent(),
			TargetType:     "wallet",
			TargetID:       wallet.WalletID,
			Message:        message,
			Component:      "wallet_handler",
			Metadata: map[string]interface{}{
				"amount":            req.Amount,
				"currency":          req.Currency,
				"bank_code":         req.BankCode,
				"account_number_masked": maskAccountNumber(req.AccountNumber),
			},
		})
	}

	if err != nil {
		// Don't expose internal error details
		response.BadRequest(w, "withdrawal request failed")
		return
	}

	response.Success(w, result)
}

// Transfer handles POST /api/wallet/transfer
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	var req struct {
		RecipientPhone string `json:"recipient_phone"`
		Amount         int64  `json:"amount"`
		Currency       string `json:"currency"`
		Description    string `json:"description"`
		PIN            string `json:"pin"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	// Input validation with security checks
	v := validation.NewValidator()
	v.Required("recipient_phone", req.RecipientPhone).
		Required("pin", req.PIN).
		Positive("amount", req.Amount).
		Min("amount", req.Amount, 100).           // Minimum 100 kobo
		Max("amount", req.Amount, 100000000).     // Maximum 1M Naira
		Phone("recipient_phone", req.RecipientPhone).
		PIN("pin", req.PIN).
		SafeString("description", req.Description)

	if v.HasErrors() {
		response.ValidationError(w, v.Errors().Errors)
		return
	}

	// Normalize phone number
	normalizedPhone := validation.NormalizePhone(req.RecipientPhone)

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	result, err := h.transferHandler.Handle(r.Context(), command.Transfer{
		FromUserID:  userID.String(),
		ToUserPhone: normalizedPhone,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		PIN:         req.PIN,
		RequestedBy: userID.String(),
	})

	// Audit log the transfer attempt (sensitive operation)
	if h.auditLogger != nil {
		outcome := audit.OutcomeSuccess
		message := "Transfer completed successfully"
		if err != nil {
			outcome = audit.OutcomeFailure
			message = "Transfer failed"
		}
		h.auditLogger.LogTransaction(r.Context(), audit.AuditEvent{
			EventAction:    audit.ActionCreate,
			EventOutcome:   outcome,
			ActorUserID:    userID.String(),
			ActorIPAddress: getClientIP(r),
			ActorUserAgent: r.UserAgent(),
			TargetType:     "transfer",
			TargetID:       "", // Transaction ID would be in result
			Message:        message,
			Component:      "wallet_handler",
			Metadata: map[string]interface{}{
				"amount":               req.Amount,
				"currency":             req.Currency,
				"recipient_phone_masked": maskPhone(normalizedPhone),
			},
		})
	}

	if err != nil {
		// Don't expose internal error details
		response.BadRequest(w, "transfer failed")
		return
	}

	response.Success(w, result)
}

// GetBanks handles GET /api/wallet/banks
func (h *WalletHandler) GetBanks(w http.ResponseWriter, r *http.Request) {
	// Get banks via withdraw handler
	banks, err := h.withdrawHandler.HandleGetBanks(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Success(w, banks)
}

// ResolveAccount handles POST /api/wallet/resolve-account
func (h *WalletHandler) ResolveAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.AccountNumber == "" || req.BankCode == "" {
		response.ValidationError(w, map[string]string{
			"account_number": "account number is required",
			"bank_code":      "bank code is required",
		})
		return
	}

	accountName, err := h.withdrawHandler.HandleVerifyAccount(r.Context(), req.BankCode, req.AccountNumber)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, map[string]string{
		"account_name":   accountName,
		"account_number": req.AccountNumber,
		"bank_code":      req.BankCode,
	})
}

// GetBankAccounts handles GET /api/wallet/bank-accounts
func (h *WalletHandler) GetBankAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized")
		return
	}

	accounts, err := h.queryHandler.HandleGetBankAccounts(r.Context(), query.GetBankAccounts{
		UserID: userID.String(),
	})
	if err != nil {
		response.InternalError(w)
		return
	}

	response.Success(w, accounts)
}

// parseIntQuery parses an integer from a query string with a default value
func parseIntQuery(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var val int
	if _, err := jsonParseInt(s, &val); err != nil {
		return defaultVal
	}
	return val
}

func jsonParseInt(s string, v *int) (bool, error) {
	return true, json.Unmarshal([]byte(s), v)
}

// getClientIP extracts the real client IP from request headers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	addr := r.RemoteAddr
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

// maskAccountNumber masks an account number for audit logging
// Shows only last 4 digits: ****1234
func maskAccountNumber(accountNumber string) string {
	if len(accountNumber) <= 4 {
		return "****"
	}
	return "****" + accountNumber[len(accountNumber)-4:]
}

// maskPhone masks a phone number for audit logging
// Shows only last 4 digits: +234****1234
func maskPhone(phone string) string {
	if len(phone) <= 4 {
		return "****"
	}
	if len(phone) > 8 {
		return phone[:4] + "****" + phone[len(phone)-4:]
	}
	return "****" + phone[len(phone)-4:]
}
