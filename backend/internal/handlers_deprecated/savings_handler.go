// +build ignore

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"hustlex/internal/services"
)

// SavingsHandler handles savings circle-related HTTP requests
type SavingsHandler struct {
	savingsService *services.SavingsService
}

// NewSavingsHandler creates a new savings handler
func NewSavingsHandler(savingsService *services.SavingsService) *SavingsHandler {
	return &SavingsHandler{savingsService: savingsService}
}

// CreateCircleRequest represents the request body for creating a savings circle
type CreateCircleRequest struct {
	Name               string `json:"name" validate:"required,min=3,max=100"`
	Description        string `json:"description" validate:"omitempty,max=500"`
	Type               string `json:"type" validate:"required,oneof=rotational fixed_target emergency"`
	ContributionAmount int64  `json:"contribution_amount" validate:"required,min=100000"` // Min ₦1,000 in kobo
	Frequency          string `json:"frequency" validate:"required,oneof=daily weekly biweekly monthly"`
	MaxMembers         int    `json:"max_members" validate:"required,min=2,max=50"`
	StartDate          string `json:"start_date" validate:"required"` // ISO date
	TargetAmount       int64  `json:"target_amount" validate:"omitempty,min=100000"` // For fixed_target type
	IsPrivate          bool   `json:"is_private"`
}

// JoinCircleRequest represents the request body for joining a circle
type JoinCircleRequest struct {
	InviteCode string `json:"invite_code" validate:"required,len=8"`
}

// MakeContributionRequest represents the request body for making a contribution
type MakeContributionRequest struct {
	PIN string `json:"pin" validate:"required,len=4,numeric"`
}

// UpdateCircleRequest represents the request body for updating circle settings
type UpdateCircleRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
	IsPrivate   *bool  `json:"is_private"`
}

// CreateCircle handles creating a new savings circle
// @Summary Create savings circle
// @Description Create a new Ajo/Esusu savings circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCircleRequest true "Circle details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /savings/circles [post]
func (h *SavingsHandler) CreateCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req CreateCircleRequest
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

	// Validate type-specific requirements
	if req.Type == "fixed_target" && req.TargetAmount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Target amount is required for fixed target circles",
		})
	}

	input := services.CreateCircleInput{
		CreatorID:          userID,
		Name:               req.Name,
		Description:        req.Description,
		Type:               req.Type,
		ContributionAmount: req.ContributionAmount,
		Frequency:          req.Frequency,
		MaxMembers:         req.MaxMembers,
		StartDate:          req.StartDate,
		TargetAmount:       req.TargetAmount,
		IsPrivate:          req.IsPrivate,
	}

	circle, err := h.savingsService.CreateCircle(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Savings circle created successfully",
		"data": fiber.Map{
			"circle": circle,
		},
	})
}

// GetCircle handles fetching a single circle
// @Summary Get circle details
// @Description Get detailed information about a savings circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /savings/circles/{id} [get]
func (h *SavingsHandler) GetCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	circle, err := h.savingsService.GetCircle(c.Context(), circleID, userID)
	if err != nil {
		status := fiber.StatusNotFound
		if err.Error() == "not authorized to view this circle" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"circle": circle,
		},
	})
}

// ListPublicCircles handles listing public circles
// @Summary List public circles
// @Description Get paginated list of public savings circles to join
// @Tags Savings
// @Accept json
// @Produce json
// @Param type query string false "Filter by type (rotational, fixed_target, emergency)"
// @Param frequency query string false "Filter by frequency"
// @Param min_contribution query int false "Minimum contribution amount"
// @Param max_contribution query int false "Maximum contribution amount"
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /savings/circles/public [get]
func (h *SavingsHandler) ListPublicCircles(c *fiber.Ctx) error {
	filter := services.CircleFilter{
		Type:      c.Query("type"),
		Frequency: c.Query("frequency"),
		Search:    c.Query("search"),
	}

	if minContrib := c.QueryInt("min_contribution"); minContrib > 0 {
		filter.MinContribution = int64(minContrib)
	}

	if maxContrib := c.QueryInt("max_contribution"); maxContrib > 0 {
		filter.MaxContribution = int64(maxContrib)
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	circles, total, err := h.savingsService.ListPublicCircles(c.Context(), filter, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch circles",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"circles": circles,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetMyCircles handles fetching user's circles
// @Summary Get my circles
// @Description Get circles where user is a member
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (pending, active, completed, cancelled)"
// @Param role query string false "Filter by role (admin, member)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /savings/circles/my [get]
func (h *SavingsHandler) GetMyCircles(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	status := c.Query("status")
	role := c.Query("role")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	circles, total, err := h.savingsService.GetUserCircles(c.Context(), userID, status, role, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch circles",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"circles": circles,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// JoinCircle handles joining a circle with invite code
// @Summary Join a circle
// @Description Join a savings circle using an invite code
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body JoinCircleRequest true "Invite code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /savings/circles/join [post]
func (h *SavingsHandler) JoinCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req JoinCircleRequest
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

	membership, err := h.savingsService.JoinCircle(c.Context(), req.InviteCode, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully joined the circle",
		"data": fiber.Map{
			"membership": membership,
		},
	})
}

// JoinPublicCircle handles joining a public circle directly
// @Summary Join public circle
// @Description Join a public savings circle by ID
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /savings/circles/{id}/join [post]
func (h *SavingsHandler) JoinPublicCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	membership, err := h.savingsService.JoinPublicCircle(c.Context(), circleID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully joined the circle",
		"data": fiber.Map{
			"membership": membership,
		},
	})
}

// LeaveCircle handles leaving a circle
// @Summary Leave a circle
// @Description Leave a savings circle (only if no active contributions)
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /savings/circles/{id}/leave [post]
func (h *SavingsHandler) LeaveCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	if err := h.savingsService.LeaveCircle(c.Context(), circleID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Successfully left the circle",
	})
}

// GetCircleMembers handles fetching circle members
// @Summary Get circle members
// @Description Get list of members in a savings circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id}/members [get]
func (h *SavingsHandler) GetCircleMembers(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	members, err := h.savingsService.GetCircleMembers(c.Context(), circleID, userID)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "not authorized to view members" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"members": members,
		},
	})
}

// MakeContribution handles making a contribution
// @Summary Make contribution
// @Description Make a contribution to a savings circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Param request body MakeContributionRequest true "Transaction PIN"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /savings/circles/{id}/contribute [post]
func (h *SavingsHandler) MakeContribution(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	var req MakeContributionRequest
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

	contribution, err := h.savingsService.MakeContribution(c.Context(), circleID, userID, req.PIN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Contribution made successfully",
		"data": fiber.Map{
			"contribution": contribution,
		},
	})
}

// GetContributionHistory handles fetching contribution history
// @Summary Get contribution history
// @Description Get contribution history for a circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Param round query int false "Filter by round number"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id}/contributions [get]
func (h *SavingsHandler) GetContributionHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	round := c.QueryInt("round", 0)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	contributions, total, err := h.savingsService.GetContributionHistory(c.Context(), circleID, userID, round, page, limit)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "not authorized to view contributions" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"contributions": contributions,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetMyContributions handles fetching user's own contributions
// @Summary Get my contributions
// @Description Get authenticated user's contributions in a circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Router /savings/circles/{id}/my-contributions [get]
func (h *SavingsHandler) GetMyContributions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	contributions, stats, err := h.savingsService.GetUserContributions(c.Context(), circleID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"contributions": contributions,
			"stats":         stats,
		},
	})
}

// GetCircleLeaderboard handles fetching circle leaderboard
// @Summary Get circle leaderboard
// @Description Get contribution leaderboard for a circle
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id}/leaderboard [get]
func (h *SavingsHandler) GetCircleLeaderboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	leaderboard, err := h.savingsService.GetCircleLeaderboard(c.Context(), circleID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"leaderboard": leaderboard,
		},
	})
}

// GetPendingContributions handles fetching pending contributions
// @Summary Get pending contributions
// @Description Get user's pending contributions across all circles
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /savings/contributions/pending [get]
func (h *SavingsHandler) GetPendingContributions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	contributions, err := h.savingsService.GetPendingContributions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch pending contributions",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"contributions": contributions,
		},
	})
}

// StartCircle handles starting a circle (admin only)
// @Summary Start circle
// @Description Start an active savings cycle (admin only)
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id}/start [post]
func (h *SavingsHandler) StartCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	if err := h.savingsService.StartCircle(c.Context(), circleID, userID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "only admin can start the circle" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Circle started successfully",
	})
}

// UpdateCircle handles updating circle settings
// @Summary Update circle
// @Description Update circle settings (admin only)
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Param request body UpdateCircleRequest true "Update details"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id} [put]
func (h *SavingsHandler) UpdateCircle(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	var req UpdateCircleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}

	circle, err := h.savingsService.UpdateCircle(c.Context(), circleID, userID, updates)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "only admin can update the circle" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"circle": circle,
		},
	})
}

// RegenerateInviteCode handles regenerating invite code
// @Summary Regenerate invite code
// @Description Generate a new invite code for the circle (admin only)
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Circle ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /savings/circles/{id}/regenerate-code [post]
func (h *SavingsHandler) RegenerateInviteCode(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	circleID := c.Params("id")

	newCode, err := h.savingsService.RegenerateInviteCode(c.Context(), circleID, userID)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "only admin can regenerate invite code" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"invite_code": newCode,
		},
	})
}

// GetSavingsSummary handles fetching user's savings summary
// @Summary Get savings summary
// @Description Get summary of user's savings activity across all circles
// @Tags Savings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /savings/summary [get]
func (h *SavingsHandler) GetSavingsSummary(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	summary, err := h.savingsService.GetUserSavingsSummary(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch savings summary",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"summary": summary,
		},
	})
}
