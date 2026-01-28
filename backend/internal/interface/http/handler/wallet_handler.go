package handler

import (
	"encoding/json"
	"net/http"

	"hustlex/internal/application/wallet/command"
	"hustlex/internal/application/wallet/handler"
	"hustlex/internal/application/wallet/query"
	"hustlex/internal/interface/http/middleware"
	"hustlex/internal/interface/http/response"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	depositHandler  *handler.DepositHandler
	withdrawHandler *handler.WithdrawHandler
	transferHandler *handler.TransferHandler
	queryHandler    *query.WalletQueryHandler
}

// NewWalletHandler creates a new wallet HTTP handler
func NewWalletHandler(
	depositHandler *handler.DepositHandler,
	withdrawHandler *handler.WithdrawHandler,
	transferHandler *handler.TransferHandler,
	queryHandler *query.WalletQueryHandler,
) *WalletHandler {
	return &WalletHandler{
		depositHandler:  depositHandler,
		withdrawHandler: withdrawHandler,
		transferHandler: transferHandler,
		queryHandler:    queryHandler,
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

	if req.Amount <= 0 {
		response.ValidationError(w, map[string]string{"amount": "amount must be positive"})
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
	if err != nil {
		response.BadRequest(w, err.Error())
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

	// Validate required fields
	errors := make(map[string]string)
	if req.Amount <= 0 {
		errors["amount"] = "amount must be positive"
	}
	if req.AccountNumber == "" {
		errors["account_number"] = "account number is required"
	}
	if req.BankCode == "" {
		errors["bank_code"] = "bank code is required"
	}
	if req.PIN == "" {
		errors["pin"] = "PIN is required"
	}

	if len(errors) > 0 {
		response.ValidationError(w, errors)
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
	if err != nil {
		response.BadRequest(w, err.Error())
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

	// Validate required fields
	errors := make(map[string]string)
	if req.Amount <= 0 {
		errors["amount"] = "amount must be positive"
	}
	if req.RecipientPhone == "" {
		errors["recipient_phone"] = "recipient phone is required"
	}
	if req.PIN == "" {
		errors["pin"] = "PIN is required"
	}

	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	if req.Currency == "" {
		req.Currency = "NGN"
	}

	result, err := h.transferHandler.Handle(r.Context(), command.Transfer{
		FromUserID:  userID.String(),
		ToUserPhone: req.RecipientPhone,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		PIN:         req.PIN,
		RequestedBy: userID.String(),
	})
	if err != nil {
		response.BadRequest(w, err.Error())
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
