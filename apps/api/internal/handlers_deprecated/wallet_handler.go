// +build ignore

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"hustlex/internal/services"
)

// WalletHandler handles wallet-related HTTP requests
type WalletHandler struct {
	walletService *services.WalletService
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

// InitiateDepositRequest represents the request body for initiating a deposit
type InitiateDepositRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=100000"` // Min ₦1,000 in kobo
	PaymentMethod string `json:"payment_method" validate:"required,oneof=card bank_transfer ussd"`
}

// WithdrawRequest represents the request body for withdrawal
type WithdrawRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=100000"` // Min ₦1,000 in kobo
	BankCode    string `json:"bank_code" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required,len=10,numeric"`
	PIN         string `json:"pin" validate:"required,len=4,numeric"`
}

// TransferRequest represents the request body for P2P transfer
type TransferRequest struct {
	RecipientPhone string `json:"recipient_phone" validate:"required,e164"`
	Amount         int64  `json:"amount" validate:"required,min=10000"` // Min ₦100 in kobo
	Note           string `json:"note" validate:"omitempty,max=200"`
	PIN            string `json:"pin" validate:"required,len=4,numeric"`
}

// AddBankAccountRequest represents the request body for adding a bank account
type AddBankAccountRequest struct {
	BankCode      string `json:"bank_code" validate:"required"`
	AccountNumber string `json:"account_number" validate:"required,len=10,numeric"`
	AccountName   string `json:"account_name" validate:"required,min=3,max=100"`
	IsDefault     bool   `json:"is_default"`
}

// GetWallet handles fetching user's wallet
// @Summary Get wallet
// @Description Get user's wallet balance and details
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /wallet [get]
func (h *WalletHandler) GetWallet(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	wallet, err := h.walletService.GetWallet(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch wallet",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"wallet": fiber.Map{
				"balance":         wallet.Balance,
				"escrow_balance":  wallet.EscrowBalance,
				"savings_balance": wallet.SavingsBalance,
				"total_earned":    wallet.TotalEarned,
				"total_withdrawn": wallet.TotalWithdrawn,
				"currency":        wallet.Currency,
				"is_locked":       wallet.IsLocked,
				"has_pin":         wallet.PINHash != "",
				"updated_at":      wallet.UpdatedAt,
			},
		},
	})
}

// InitiateDeposit handles deposit initiation
// @Summary Initiate deposit
// @Description Start a deposit transaction via payment gateway
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body InitiateDepositRequest true "Deposit details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/deposit [post]
func (h *WalletHandler) InitiateDeposit(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req InitiateDepositRequest
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

	paymentData, err := h.walletService.InitiateDeposit(c.Context(), userID, req.Amount, req.PaymentMethod)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Deposit initiated",
		"data":    paymentData,
	})
}

// VerifyDeposit handles deposit verification callback
// @Summary Verify deposit
// @Description Verify a deposit payment (webhook handler)
// @Tags Wallet
// @Accept json
// @Produce json
// @Param reference path string true "Payment reference"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/deposit/verify/{reference} [post]
func (h *WalletHandler) VerifyDeposit(c *fiber.Ctx) error {
	reference := c.Params("reference")

	transaction, err := h.walletService.VerifyDeposit(c.Context(), reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Deposit verified",
		"data": fiber.Map{
			"transaction": transaction,
		},
	})
}

// Withdraw handles withdrawal request
// @Summary Withdraw funds
// @Description Withdraw funds to a bank account
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body WithdrawRequest true "Withdrawal details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/withdraw [post]
func (h *WalletHandler) Withdraw(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req WithdrawRequest
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

	input := services.WithdrawInput{
		UserID:        userID,
		Amount:        req.Amount,
		BankCode:      req.BankCode,
		AccountNumber: req.AccountNumber,
		PIN:           req.PIN,
	}

	transaction, err := h.walletService.Withdraw(c.Context(), input)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "wallet is locked" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Withdrawal initiated",
		"data": fiber.Map{
			"transaction": transaction,
		},
	})
}

// Transfer handles P2P transfer
// @Summary Transfer funds
// @Description Transfer funds to another HustleX user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransferRequest true "Transfer details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/transfer [post]
func (h *WalletHandler) Transfer(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req TransferRequest
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

	input := services.TransferInput{
		SenderID:       userID,
		RecipientPhone: req.RecipientPhone,
		Amount:         req.Amount,
		Note:           req.Note,
		PIN:            req.PIN,
	}

	transaction, err := h.walletService.Transfer(c.Context(), input)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "wallet is locked" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transfer successful",
		"data": fiber.Map{
			"transaction": transaction,
		},
	})
}

// GetTransactions handles fetching transaction history
// @Summary Get transactions
// @Description Get paginated transaction history
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Filter by type (deposit, withdrawal, transfer_in, transfer_out, gig_payment, gig_earning, escrow_hold, escrow_release, contribution, payout, loan_disbursement, loan_repayment)"
// @Param status query string false "Filter by status (pending, completed, failed, cancelled)"
// @Param start_date query string false "Start date (ISO format)"
// @Param end_date query string false "End date (ISO format)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /wallet/transactions [get]
func (h *WalletHandler) GetTransactions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	filter := services.TransactionFilter{
		Type:      c.Query("type"),
		Status:    c.Query("status"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	transactions, total, err := h.walletService.GetTransactions(c.Context(), userID, filter, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch transactions",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"transactions": transactions,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetTransaction handles fetching a single transaction
// @Summary Get transaction details
// @Description Get detailed information about a specific transaction
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wallet/transactions/{id} [get]
func (h *WalletHandler) GetTransaction(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	transactionID := c.Params("id")

	transaction, err := h.walletService.GetTransaction(c.Context(), transactionID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Transaction not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"transaction": transaction,
		},
	})
}

// AddBankAccount handles adding a bank account
// @Summary Add bank account
// @Description Add a bank account for withdrawals
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AddBankAccountRequest true "Bank account details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/bank-accounts [post]
func (h *WalletHandler) AddBankAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req AddBankAccountRequest
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

	input := services.AddBankAccountInput{
		UserID:        userID,
		BankCode:      req.BankCode,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		IsDefault:     req.IsDefault,
	}

	account, err := h.walletService.AddBankAccount(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Bank account added successfully",
		"data": fiber.Map{
			"bank_account": account,
		},
	})
}

// GetBankAccounts handles fetching user's bank accounts
// @Summary Get bank accounts
// @Description Get user's saved bank accounts
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /wallet/bank-accounts [get]
func (h *WalletHandler) GetBankAccounts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	accounts, err := h.walletService.GetBankAccounts(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch bank accounts",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"bank_accounts": accounts,
		},
	})
}

// DeleteBankAccount handles deleting a bank account
// @Summary Delete bank account
// @Description Remove a saved bank account
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bank account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wallet/bank-accounts/{id} [delete]
func (h *WalletHandler) DeleteBankAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	accountID := c.Params("id")

	if err := h.walletService.DeleteBankAccount(c.Context(), accountID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Bank account deleted",
	})
}

// SetDefaultBankAccount handles setting default bank account
// @Summary Set default bank account
// @Description Set a bank account as default for withdrawals
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bank account ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /wallet/bank-accounts/{id}/default [post]
func (h *WalletHandler) SetDefaultBankAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	accountID := c.Params("id")

	if err := h.walletService.SetDefaultBankAccount(c.Context(), accountID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Default bank account updated",
	})
}

// GetBanks handles fetching list of supported banks
// @Summary Get banks
// @Description Get list of supported banks for withdrawals
// @Tags Wallet
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /wallet/banks [get]
func (h *WalletHandler) GetBanks(c *fiber.Ctx) error {
	banks, err := h.walletService.GetSupportedBanks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch banks",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"banks": banks,
		},
	})
}

// VerifyBankAccount handles bank account verification
// @Summary Verify bank account
// @Description Verify a bank account number and get account name
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param bank_code query string true "Bank code"
// @Param account_number query string true "Account number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/verify-account [get]
func (h *WalletHandler) VerifyBankAccount(c *fiber.Ctx) error {
	bankCode := c.Query("bank_code")
	accountNumber := c.Query("account_number")

	if bankCode == "" || accountNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Bank code and account number are required",
		})
	}

	accountName, err := h.walletService.VerifyBankAccount(c.Context(), bankCode, accountNumber)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"account_name": accountName,
		},
	})
}

// GetWalletStats handles fetching wallet statistics
// @Summary Get wallet stats
// @Description Get wallet statistics and analytics
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param period query string false "Period (week, month, year, all)"
// @Success 200 {object} map[string]interface{}
// @Router /wallet/stats [get]
func (h *WalletHandler) GetWalletStats(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	period := c.Query("period", "month")

	stats, err := h.walletService.GetWalletStats(c.Context(), userID, period)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch stats",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"stats": stats,
		},
	})
}
