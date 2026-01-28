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

	balance, err := h.queryHandler.HandleGetBalance(r.Context(), query.GetBalance{
		UserID: userID.String(),
	})
	if err != nil {
		response.NotFound(w, "wallet not found")
		return
	}

	response.Success(w, balance)
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

	result, err := h.queryHandler.HandleGetTransactions(r.Context(), query.GetTransactions{
		UserID: userID.String(),
		Type:   q.Get("type"),
		Status: q.Get("status"),
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
		Amount      int64    `json:"amount"`
		Currency    string   `json:"currency"`
		CallbackURL string   `json:"callback_url"`
		Channels    []string `json:"channels"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Amount <= 0 {
		response.ValidationError(w, map[string]string{"amount": "amount must be positive"})
		return
	}

	result, err := h.depositHandler.HandleInitiateDeposit(r.Context(), command.InitiateDeposit{
		UserID:      userID.String(),
		Amount:      req.Amount,
		Currency:    req.Currency,
		CallbackURL: req.CallbackURL,
		Channels:    req.Channels,
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

	result, err := h.depositHandler.HandleVerifyDeposit(r.Context(), command.VerifyDeposit{
		Reference: req.Reference,
	})
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
		Reason        string `json:"reason"`
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

	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	result, err := h.withdrawHandler.HandleWithdraw(r.Context(), command.Withdraw{
		UserID:        userID.String(),
		Amount:        req.Amount,
		Currency:      req.Currency,
		AccountNumber: req.AccountNumber,
		BankCode:      req.BankCode,
		AccountName:   req.AccountName,
		Reason:        req.Reason,
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
		RecipientID string `json:"recipient_id"`
		Amount      int64  `json:"amount"`
		Currency    string `json:"currency"`
		Description string `json:"description"`
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
	if req.RecipientID == "" {
		errors["recipient_id"] = "recipient ID is required"
	}

	if len(errors) > 0 {
		response.ValidationError(w, errors)
		return
	}

	result, err := h.transferHandler.HandleTransfer(r.Context(), command.Transfer{
		FromUserID:  userID.String(),
		ToUserID:    req.RecipientID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
	})
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, result)
}

// GetBanks handles GET /api/wallet/banks
func (h *WalletHandler) GetBanks(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		country = "nigeria"
	}

	banks, err := h.queryHandler.HandleGetBanks(r.Context(), query.GetBanks{
		Country: country,
	})
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

	account, err := h.queryHandler.HandleResolveAccount(r.Context(), query.ResolveAccount{
		AccountNumber: req.AccountNumber,
		BankCode:      req.BankCode,
	})
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, account)
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
