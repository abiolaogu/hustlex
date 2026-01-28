// +build ignore

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"hustlex/internal/services"
)

// CreditHandler handles credit-related HTTP requests
type CreditHandler struct {
	creditService *services.CreditService
}

// NewCreditHandler creates a new credit handler
func NewCreditHandler(creditService *services.CreditService) *CreditHandler {
	return &CreditHandler{creditService: creditService}
}

// ApplyLoanRequest represents the request body for loan application
type ApplyLoanRequest struct {
	Amount  int64  `json:"amount" validate:"required,min=500000,max=50000000"` // ₦5,000 - ₦500,000 in kobo
	Purpose string `json:"purpose" validate:"required,min=10,max=500"`
	Tenure  int    `json:"tenure" validate:"required,min=7,max=90"` // Days
}

// RepayLoanRequest represents the request body for loan repayment
type RepayLoanRequest struct {
	Amount int64  `json:"amount" validate:"required,min=10000"` // Min ₦100 in kobo
	PIN    string `json:"pin" validate:"required,len=4,numeric"`
}

// GetCreditScore handles fetching user's credit score
// @Summary Get credit score
// @Description Get user's Hustle Credit score and breakdown
// @Tags Credit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /credit/score [get]
func (h *CreditHandler) GetCreditScore(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	score, err := h.creditService.GetCreditScore(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch credit score",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"credit_score": fiber.Map{
				"total_score":           score.TotalScore,
				"tier":                  score.Tier,
				"gig_completion_score":  score.GigCompletionScore,
				"rating_score":          score.RatingScore,
				"ajo_compliance_score":  score.AjoComplianceScore,
				"account_age_score":     score.AccountAgeScore,
				"verification_score":    score.VerificationScore,
				"community_score":       score.CommunityScore,
				"total_gigs_completed":  score.TotalGigsCompleted,
				"total_gigs_cancelled":  score.TotalGigsCancelled,
				"average_rating":        score.AverageRating,
				"ajo_contributions":     score.AjoContributions,
				"ajo_missed_payments":   score.AjoMissedPayments,
				"last_calculated_at":    score.LastCalculatedAt,
				"updated_at":            score.UpdatedAt,
			},
		},
	})
}

// GetCreditHistory handles fetching credit score history
// @Summary Get credit history
// @Description Get historical credit score changes
// @Tags Credit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /credit/history [get]
func (h *CreditHandler) GetCreditHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	history, total, err := h.creditService.GetCreditHistory(c.Context(), userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch credit history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"history": history,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetCreditTips handles fetching tips to improve credit score
// @Summary Get credit tips
// @Description Get personalized tips to improve credit score
// @Tags Credit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /credit/tips [get]
func (h *CreditHandler) GetCreditTips(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	tips, err := h.creditService.GetCreditImprovementTips(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch tips",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tips": tips,
		},
	})
}

// GetLoanEligibility handles checking loan eligibility
// @Summary Check loan eligibility
// @Description Check user's eligibility for micro-loans
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /credit/loans/eligibility [get]
func (h *CreditHandler) GetLoanEligibility(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	eligibility, err := h.creditService.CheckLoanEligibility(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to check eligibility",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"eligibility": eligibility,
		},
	})
}

// ApplyForLoan handles loan application
// @Summary Apply for loan
// @Description Apply for a micro-loan
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ApplyLoanRequest true "Loan details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /credit/loans/apply [post]
func (h *CreditHandler) ApplyForLoan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req ApplyLoanRequest
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

	input := services.LoanApplicationInput{
		UserID:       userID,
		Amount:       req.Amount,
		Purpose:      req.Purpose,
		TenureDays:   req.Tenure,
	}

	loan, err := h.creditService.ApplyForLoan(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Loan application submitted",
		"data": fiber.Map{
			"loan": loan,
		},
	})
}

// GetLoans handles fetching user's loans
// @Summary Get loans
// @Description Get user's loan history
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (pending, approved, disbursed, repaid, defaulted, rejected)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /credit/loans [get]
func (h *CreditHandler) GetLoans(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	loans, total, err := h.creditService.GetUserLoans(c.Context(), userID, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch loans",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"loans": loans,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetLoan handles fetching a single loan
// @Summary Get loan details
// @Description Get detailed information about a specific loan
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /credit/loans/{id} [get]
func (h *CreditHandler) GetLoan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	loanID := c.Params("id")

	loan, repayments, err := h.creditService.GetLoan(c.Context(), loanID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Loan not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"loan":       loan,
			"repayments": repayments,
		},
	})
}

// RepayLoan handles loan repayment
// @Summary Repay loan
// @Description Make a loan repayment
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Loan ID"
// @Param request body RepayLoanRequest true "Repayment details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /credit/loans/{id}/repay [post]
func (h *CreditHandler) RepayLoan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	loanID := c.Params("id")

	var req RepayLoanRequest
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

	repayment, err := h.creditService.RepayLoan(c.Context(), loanID, userID, req.Amount, req.PIN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Repayment successful",
		"data": fiber.Map{
			"repayment": repayment,
		},
	})
}

// GetLoanSchedule handles fetching repayment schedule
// @Summary Get repayment schedule
// @Description Get the repayment schedule for a loan
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Loan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /credit/loans/{id}/schedule [get]
func (h *CreditHandler) GetLoanSchedule(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	loanID := c.Params("id")

	schedule, err := h.creditService.GetRepaymentSchedule(c.Context(), loanID, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"schedule": schedule,
		},
	})
}

// GetActiveLoan handles fetching user's active loan
// @Summary Get active loan
// @Description Get user's currently active loan if any
// @Tags Loans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /credit/loans/active [get]
func (h *CreditHandler) GetActiveLoan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	loan, err := h.creditService.GetActiveLoan(c.Context(), userID)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"loan": nil,
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"loan": loan,
		},
	})
}

// GetTierBenefits handles fetching tier benefits
// @Summary Get tier benefits
// @Description Get benefits for each credit tier
// @Tags Credit
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /credit/tiers [get]
func (h *CreditHandler) GetTierBenefits(c *fiber.Ctx) error {
	tiers := []fiber.Map{
		{
			"tier":       "bronze",
			"min_score":  0,
			"max_score":  300,
			"benefits": []string{
				"Access to basic gigs",
				"Join up to 2 savings circles",
				"Standard platform fees",
			},
			"loan_limit":    0,
			"interest_rate": 0,
		},
		{
			"tier":       "silver",
			"min_score":  301,
			"max_score":  500,
			"benefits": []string{
				"Priority in gig matching",
				"Join up to 5 savings circles",
				"5% discount on platform fees",
				"Access to micro-loans up to ₦10,000",
			},
			"loan_limit":    1000000, // ₦10,000 in kobo
			"interest_rate": 5.0,
		},
		{
			"tier":       "gold",
			"min_score":  501,
			"max_score":  700,
			"benefits": []string{
				"Featured profile badge",
				"Join up to 10 savings circles",
				"10% discount on platform fees",
				"Access to loans up to ₦100,000",
				"Lower interest rates (4%)",
				"Early access to new features",
			},
			"loan_limit":    10000000, // ₦100,000 in kobo
			"interest_rate": 4.0,
		},
		{
			"tier":       "platinum",
			"min_score":  701,
			"max_score":  850,
			"benefits": []string{
				"Premium badge and visibility",
				"Unlimited savings circles",
				"15% discount on platform fees",
				"Access to loans up to ₦500,000",
				"Lowest interest rates (3%)",
				"Priority customer support",
				"Create premium savings circles",
				"Invite to exclusive events",
			},
			"loan_limit":    50000000, // ₦500,000 in kobo
			"interest_rate": 3.0,
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tiers": tiers,
		},
	})
}

// RecalculateCreditScore handles manual credit score recalculation
// @Summary Recalculate credit score
// @Description Trigger a manual recalculation of credit score
// @Tags Credit
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /credit/recalculate [post]
func (h *CreditHandler) RecalculateCreditScore(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	score, err := h.creditService.RecalculateScore(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to recalculate credit score",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Credit score recalculated",
		"data": fiber.Map{
			"credit_score": fiber.Map{
				"total_score": score.TotalScore,
				"tier":        score.Tier,
			},
		},
	})
}
