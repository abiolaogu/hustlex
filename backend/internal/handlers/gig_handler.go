package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hustlex/backend/internal/models"
	"github.com/hustlex/backend/internal/services"
)

// GigHandler handles gig-related HTTP requests
type GigHandler struct {
	gigService *services.GigService
}

// NewGigHandler creates a new gig handler
func NewGigHandler(gigService *services.GigService) *GigHandler {
	return &GigHandler{gigService: gigService}
}

// CreateGigRequest represents the request body for creating a gig
type CreateGigRequest struct {
	Title         string   `json:"title" validate:"required,min=10,max=200"`
	Description   string   `json:"description" validate:"required,min=50,max=5000"`
	CategoryID    string   `json:"category_id" validate:"required,uuid"`
	SkillIDs      []string `json:"skill_ids" validate:"required,min=1,max=5,dive,uuid"`
	BudgetMin     int64    `json:"budget_min" validate:"required,min=100000"` // Min ₦1,000 in kobo
	BudgetMax     int64    `json:"budget_max" validate:"required,gtefield=BudgetMin"`
	Deadline      string   `json:"deadline" validate:"required"` // ISO date
	Location      string   `json:"location" validate:"omitempty,max=100"`
	IsRemote      bool     `json:"is_remote"`
	Attachments   []string `json:"attachments" validate:"omitempty,max=10,dive,url"`
	Requirements  string   `json:"requirements" validate:"omitempty,max=2000"`
}

// UpdateGigRequest represents the request body for updating a gig
type UpdateGigRequest struct {
	Title        string   `json:"title" validate:"omitempty,min=10,max=200"`
	Description  string   `json:"description" validate:"omitempty,min=50,max=5000"`
	BudgetMin    int64    `json:"budget_min" validate:"omitempty,min=100000"`
	BudgetMax    int64    `json:"budget_max" validate:"omitempty"`
	Deadline     string   `json:"deadline" validate:"omitempty"`
	Location     string   `json:"location" validate:"omitempty,max=100"`
	IsRemote     *bool    `json:"is_remote"`
	Attachments  []string `json:"attachments" validate:"omitempty,max=10,dive,url"`
	Requirements string   `json:"requirements" validate:"omitempty,max=2000"`
}

// SubmitProposalRequest represents the request body for submitting a proposal
type SubmitProposalRequest struct {
	CoverLetter  string `json:"cover_letter" validate:"required,min=50,max=2000"`
	ProposedRate int64  `json:"proposed_rate" validate:"required,min=100000"` // Min ₦1,000 in kobo
	DeliveryDays int    `json:"delivery_days" validate:"required,min=1,max=365"`
	Attachments  []string `json:"attachments" validate:"omitempty,max=5,dive,url"`
}

// DeliverWorkRequest represents the request body for delivering work
type DeliverWorkRequest struct {
	Message     string   `json:"message" validate:"required,min=10,max=2000"`
	Attachments []string `json:"attachments" validate:"required,min=1,max=10,dive,url"`
}

// SubmitReviewRequest represents the request body for submitting a review
type SubmitReviewRequest struct {
	Rating              int    `json:"rating" validate:"required,min=1,max=5"`
	Comment             string `json:"comment" validate:"required,min=10,max=1000"`
	CommunicationRating int    `json:"communication_rating" validate:"required,min=1,max=5"`
	QualityRating       int    `json:"quality_rating" validate:"required,min=1,max=5"`
	TimelinessRating    int    `json:"timeliness_rating" validate:"required,min=1,max=5"`
}

// CreateGig handles gig creation
// @Summary Create a new gig
// @Description Post a new gig/job to the marketplace
// @Tags Gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateGigRequest true "Gig details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /gigs [post]
func (h *GigHandler) CreateGig(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req CreateGigRequest
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

	input := services.CreateGigInput{
		ClientID:     userID,
		Title:        req.Title,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		SkillIDs:     req.SkillIDs,
		BudgetMin:    req.BudgetMin,
		BudgetMax:    req.BudgetMax,
		Deadline:     req.Deadline,
		Location:     req.Location,
		IsRemote:     req.IsRemote,
		Attachments:  req.Attachments,
		Requirements: req.Requirements,
	}

	gig, err := h.gigService.CreateGig(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"gig": gig,
		},
	})
}

// GetGig handles fetching a single gig
// @Summary Get gig details
// @Description Get detailed information about a specific gig
// @Tags Gigs
// @Accept json
// @Produce json
// @Param id path string true "Gig ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /gigs/{id} [get]
func (h *GigHandler) GetGig(c *fiber.Ctx) error {
	gigID := c.Params("id")

	gig, err := h.gigService.GetGig(c.Context(), gigID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Gig not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"gig": gig,
		},
	})
}

// ListGigs handles listing gigs with filters
// @Summary List gigs
// @Description Get paginated list of gigs with optional filters
// @Tags Gigs
// @Accept json
// @Produce json
// @Param category query string false "Category ID"
// @Param skill query string false "Skill ID"
// @Param min_budget query int false "Minimum budget in kobo"
// @Param max_budget query int false "Maximum budget in kobo"
// @Param location query string false "Location filter"
// @Param remote query bool false "Remote only"
// @Param search query string false "Search term"
// @Param sort query string false "Sort by (newest, budget_high, budget_low, deadline, popular)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /gigs [get]
func (h *GigHandler) ListGigs(c *fiber.Ctx) error {
	filter := services.GigFilter{
		CategoryID: c.Query("category"),
		SkillID:    c.Query("skill"),
		Location:   c.Query("location"),
		Search:     c.Query("search"),
		Sort:       c.Query("sort", "newest"),
	}

	if minBudget := c.Query("min_budget"); minBudget != "" {
		if val, err := strconv.ParseInt(minBudget, 10, 64); err == nil {
			filter.MinBudget = val
		}
	}

	if maxBudget := c.Query("max_budget"); maxBudget != "" {
		if val, err := strconv.ParseInt(maxBudget, 10, 64); err == nil {
			filter.MaxBudget = val
		}
	}

	if remote := c.Query("remote"); remote == "true" {
		filter.IsRemote = true
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 100 {
		limit = 100
	}

	gigs, total, err := h.gigService.ListGigs(c.Context(), filter, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch gigs",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"gigs": gigs,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// UpdateGig handles gig updates
// @Summary Update a gig
// @Description Update gig details (only by owner, only if status is open)
// @Tags Gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Gig ID"
// @Param request body UpdateGigRequest true "Update details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /gigs/{id} [put]
func (h *GigHandler) UpdateGig(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	gigID := c.Params("id")

	var req UpdateGigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.BudgetMin > 0 {
		updates["budget_min"] = req.BudgetMin
	}
	if req.BudgetMax > 0 {
		updates["budget_max"] = req.BudgetMax
	}
	if req.Deadline != "" {
		updates["deadline"] = req.Deadline
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.IsRemote != nil {
		updates["is_remote"] = *req.IsRemote
	}
	if req.Attachments != nil {
		updates["attachments"] = req.Attachments
	}
	if req.Requirements != "" {
		updates["requirements"] = req.Requirements
	}

	gig, err := h.gigService.UpdateGig(c.Context(), gigID, userID, updates)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "not authorized to update this gig" {
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
			"gig": gig,
		},
	})
}

// DeleteGig handles gig deletion
// @Summary Delete a gig
// @Description Delete a gig (only by owner, only if status is open)
// @Tags Gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Gig ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /gigs/{id} [delete]
func (h *GigHandler) DeleteGig(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	gigID := c.Params("id")

	if err := h.gigService.DeleteGig(c.Context(), gigID, userID); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "not authorized to delete this gig" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Gig deleted successfully",
	})
}

// GetMyGigs handles fetching user's own gigs
// @Summary Get my gigs
// @Description Get gigs created by the authenticated user
// @Tags Gigs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /gigs/my [get]
func (h *GigHandler) GetMyGigs(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	gigs, total, err := h.gigService.GetUserGigs(c.Context(), userID, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch gigs",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"gigs": gigs,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// SubmitProposal handles proposal submission
// @Summary Submit a proposal
// @Description Submit a proposal for a gig
// @Tags Proposals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Gig ID"
// @Param request body SubmitProposalRequest true "Proposal details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /gigs/{id}/proposals [post]
func (h *GigHandler) SubmitProposal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	gigID := c.Params("id")

	var req SubmitProposalRequest
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

	input := services.SubmitProposalInput{
		GigID:        gigID,
		HustlerID:    userID,
		CoverLetter:  req.CoverLetter,
		ProposedRate: req.ProposedRate,
		DeliveryDays: req.DeliveryDays,
		Attachments:  req.Attachments,
	}

	proposal, err := h.gigService.SubmitProposal(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"proposal": proposal,
		},
	})
}

// GetGigProposals handles fetching proposals for a gig
// @Summary Get gig proposals
// @Description Get all proposals for a gig (only gig owner)
// @Tags Proposals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Gig ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /gigs/{id}/proposals [get]
func (h *GigHandler) GetGigProposals(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	gigID := c.Params("id")

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	proposals, total, err := h.gigService.GetGigProposals(c.Context(), gigID, userID, page, limit)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "not authorized to view proposals" {
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
			"proposals": proposals,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetMyProposals handles fetching user's submitted proposals
// @Summary Get my proposals
// @Description Get proposals submitted by the authenticated user
// @Tags Proposals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /proposals/my [get]
func (h *GigHandler) GetMyProposals(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	proposals, total, err := h.gigService.GetUserProposals(c.Context(), userID, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch proposals",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"proposals": proposals,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// AcceptProposal handles accepting a proposal
// @Summary Accept a proposal
// @Description Accept a proposal and create a contract (funds moved to escrow)
// @Tags Proposals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Proposal ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /proposals/{id}/accept [post]
func (h *GigHandler) AcceptProposal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	proposalID := c.Params("id")

	contract, err := h.gigService.AcceptProposal(c.Context(), proposalID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Proposal accepted, contract created",
		"data": fiber.Map{
			"contract": contract,
		},
	})
}

// RejectProposal handles rejecting a proposal
// @Summary Reject a proposal
// @Description Reject a proposal
// @Tags Proposals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Proposal ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /proposals/{id}/reject [post]
func (h *GigHandler) RejectProposal(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	proposalID := c.Params("id")

	if err := h.gigService.RejectProposal(c.Context(), proposalID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Proposal rejected",
	})
}

// GetContract handles fetching a contract
// @Summary Get contract details
// @Description Get detailed information about a contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /contracts/{id} [get]
func (h *GigHandler) GetContract(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	contractID := c.Params("id")

	contract, err := h.gigService.GetContract(c.Context(), contractID, userID)
	if err != nil {
		status := fiber.StatusNotFound
		if err.Error() == "not authorized to view this contract" {
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
			"contract": contract,
		},
	})
}

// GetMyContracts handles fetching user's contracts
// @Summary Get my contracts
// @Description Get contracts where user is client or hustler
// @Tags Contracts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role query string false "Filter by role (client, hustler)"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /contracts/my [get]
func (h *GigHandler) GetMyContracts(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	role := c.Query("role")
	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	contracts, total, err := h.gigService.GetUserContracts(c.Context(), userID, role, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch contracts",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"contracts": contracts,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// DeliverWork handles work delivery
// @Summary Deliver work
// @Description Submit completed work for a contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Param request body DeliverWorkRequest true "Delivery details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /contracts/{id}/deliver [post]
func (h *GigHandler) DeliverWork(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	contractID := c.Params("id")

	var req DeliverWorkRequest
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

	if err := h.gigService.DeliverWork(c.Context(), contractID, userID, req.Message, req.Attachments); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Work delivered successfully",
	})
}

// ApproveDelivery handles delivery approval
// @Summary Approve delivery
// @Description Approve delivered work and release payment
// @Tags Contracts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /contracts/{id}/approve [post]
func (h *GigHandler) ApproveDelivery(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	contractID := c.Params("id")

	if err := h.gigService.ApproveDelivery(c.Context(), contractID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Delivery approved, payment released",
	})
}

// RequestRevision handles revision request
// @Summary Request revision
// @Description Request revisions on delivered work
// @Tags Contracts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Param request body map[string]string true "Revision request message"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /contracts/{id}/revision [post]
func (h *GigHandler) RequestRevision(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	contractID := c.Params("id")

	var req struct {
		Message string `json:"message" validate:"required,min=10,max=1000"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	if err := h.gigService.RequestRevision(c.Context(), contractID, userID, req.Message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Revision requested",
	})
}

// SubmitReview handles review submission
// @Summary Submit review
// @Description Submit a review for a completed contract
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contract ID"
// @Param request body SubmitReviewRequest true "Review details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /contracts/{id}/review [post]
func (h *GigHandler) SubmitReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	contractID := c.Params("id")

	var req SubmitReviewRequest
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

	input := services.SubmitReviewInput{
		ContractID:          contractID,
		ReviewerID:          userID,
		Rating:              req.Rating,
		Comment:             req.Comment,
		CommunicationRating: req.CommunicationRating,
		QualityRating:       req.QualityRating,
		TimelinessRating:    req.TimelinessRating,
	}

	review, err := h.gigService.SubmitReview(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"review": review,
		},
	})
}

// GetUserReviews handles fetching reviews for a user
// @Summary Get user reviews
// @Description Get all reviews for a specific user
// @Tags Reviews
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /users/{id}/reviews [get]
func (h *GigHandler) GetUserReviews(c *fiber.Ctx) error {
	targetUserID := c.Params("id")

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	reviews, stats, total, err := h.gigService.GetUserReviews(c.Context(), targetUserID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch reviews",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"reviews": reviews,
			"stats":   stats,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetCategories handles fetching skill categories
// @Summary Get categories
// @Description Get all skill categories
// @Tags Skills
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /categories [get]
func (h *GigHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.gigService.GetCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch categories",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"categories": categories,
		},
	})
}

// GetSkills handles fetching skills
// @Summary Get skills
// @Description Get skills optionally filtered by category
// @Tags Skills
// @Accept json
// @Produce json
// @Param category query string false "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /skills [get]
func (h *GigHandler) GetSkills(c *fiber.Ctx) error {
	categoryID := c.Query("category")

	skills, err := h.gigService.GetSkills(c.Context(), categoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to fetch skills",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"skills": skills,
		},
	})
}

// Helper function to sanitize user for response
func sanitizeUser(user *models.User) fiber.Map {
	return fiber.Map{
		"id":            user.ID,
		"phone":         user.Phone,
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"avatar_url":    user.AvatarURL,
		"bio":           user.Bio,
		"location":      user.Location,
		"tier":          user.Tier,
		"referral_code": user.ReferralCode,
		"is_verified":   user.IsVerified,
		"created_at":    user.CreatedAt,
	}
}
