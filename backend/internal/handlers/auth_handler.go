package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hustlex/backend/internal/services"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RequestOTPRequest represents the request body for OTP request
type RequestOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
}

// VerifyOTPRequest represents the request body for OTP verification
type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required,e164"`
	Code  string `json:"code" validate:"required,len=6"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Phone        string `json:"phone" validate:"required,e164"`
	Code         string `json:"code" validate:"required,len=6"`
	FirstName    string `json:"first_name" validate:"required,min=2,max=50"`
	LastName     string `json:"last_name" validate:"required,min=2,max=50"`
	Email        string `json:"email" validate:"omitempty,email"`
	ReferralCode string `json:"referral_code" validate:"omitempty,len=8"`
}

// RefreshTokenRequest represents the request body for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// SetPINRequest represents the request body for setting transaction PIN
type SetPINRequest struct {
	PIN string `json:"pin" validate:"required,len=4,numeric"`
}

// VerifyPINRequest represents the request body for PIN verification
type VerifyPINRequest struct {
	PIN string `json:"pin" validate:"required,len=4,numeric"`
}

// ChangePINRequest represents the request body for changing PIN
type ChangePINRequest struct {
	OldPIN string `json:"old_pin" validate:"required,len=4,numeric"`
	NewPIN string `json:"new_pin" validate:"required,len=4,numeric"`
}

// RequestOTP handles OTP request
// @Summary Request OTP
// @Description Send OTP to phone number for authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "Phone number"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 429 {object} map[string]interface{}
// @Router /auth/otp/request [post]
func (h *AuthHandler) RequestOTP(c *fiber.Ctx) error {
	var req RequestOTPRequest
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

	code, err := h.authService.RequestOTP(c.Context(), req.Phone)
	if err != nil {
		if err.Error() == "rate limit exceeded, please wait before requesting another OTP" {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to send OTP",
		})
	}

	response := fiber.Map{
		"success": true,
		"message": "OTP sent successfully",
	}

	// In development mode, return the OTP for testing
	if code != "" {
		response["otp"] = code // Remove in production
	}

	return c.JSON(response)
}

// VerifyOTP handles OTP verification for existing users
// @Summary Verify OTP
// @Description Verify OTP and login existing user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Phone and OTP"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/otp/verify [post]
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPRequest
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

	user, tokens, err := h.authService.VerifyOTPAndLogin(c.Context(), req.Phone, req.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":   sanitizeUser(user),
			"tokens": tokens,
		},
	})
}

// Register handles new user registration
// @Summary Register new user
// @Description Register a new user with verified OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
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

	input := services.RegisterInput{
		Phone:        req.Phone,
		OTP:          req.Code,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		ReferralCode: req.ReferralCode,
	}

	user, tokens, err := h.authService.Register(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":   sanitizeUser(user),
			"tokens": tokens,
		},
	})
}

// RefreshToken handles token refresh
// @Summary Refresh tokens
// @Description Get new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	tokens, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"tokens": tokens,
		},
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.authService.Logout(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}

// SetTransactionPIN handles setting the transaction PIN
// @Summary Set transaction PIN
// @Description Set a 4-digit transaction PIN for wallet operations
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SetPINRequest true "PIN"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/pin/set [post]
func (h *AuthHandler) SetTransactionPIN(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req SetPINRequest
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

	if err := h.authService.SetTransactionPIN(c.Context(), userID, req.PIN); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transaction PIN set successfully",
	})
}

// VerifyTransactionPIN handles PIN verification
// @Summary Verify transaction PIN
// @Description Verify the user's transaction PIN
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body VerifyPINRequest true "PIN"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/pin/verify [post]
func (h *AuthHandler) VerifyTransactionPIN(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req VerifyPINRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	valid, err := h.authService.VerifyTransactionPIN(c.Context(), userID, req.PIN)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "wallet is locked due to too many failed attempts" {
			status = fiber.StatusForbidden
		}
		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if !valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid PIN",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "PIN verified successfully",
	})
}

// ChangeTransactionPIN handles PIN change
// @Summary Change transaction PIN
// @Description Change the user's transaction PIN
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePINRequest true "Old and new PIN"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/pin/change [post]
func (h *AuthHandler) ChangeTransactionPIN(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req ChangePINRequest
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

	if err := h.authService.ChangeTransactionPIN(c.Context(), userID, req.OldPIN, req.NewPIN); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Transaction PIN changed successfully",
	})
}

// GetMe handles getting current user profile
// @Summary Get current user
// @Description Get the authenticated user's profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.authService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": sanitizeUser(user),
		},
	})
}
