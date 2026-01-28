package handler

import (
	"context"
	"errors"

	"hustlex/internal/application/identity/command"
	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/identity/service"
	"hustlex/internal/domain/shared/valueobject"
)

// AuthHandler handles authentication commands
type AuthHandler struct {
	userRepo            repository.UserRepository
	authService         *service.AuthenticationService
	registrationService *service.RegistrationService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userRepo repository.UserRepository,
	authService *service.AuthenticationService,
	registrationService *service.RegistrationService,
) *AuthHandler {
	return &AuthHandler{
		userRepo:            userRepo,
		authService:         authService,
		registrationService: registrationService,
	}
}

// HandleSendOTP sends an OTP to the specified phone number
func (h *AuthHandler) HandleSendOTP(ctx context.Context, cmd command.SendOTP) error {
	phone, err := cmd.GetPhone()
	if err != nil {
		return err
	}

	// Validate purpose
	if cmd.Purpose != "login" && cmd.Purpose != "register" && cmd.Purpose != "reset_pin" {
		return errors.New("invalid OTP purpose")
	}

	// For login, verify user exists
	if cmd.Purpose == "login" {
		exists, err := h.userRepo.ExistsByPhone(ctx, phone)
		if err != nil {
			return err
		}
		if !exists {
			return service.ErrUserNotFound
		}
	}

	// For register, verify user doesn't exist
	if cmd.Purpose == "register" {
		exists, err := h.userRepo.ExistsByPhone(ctx, phone)
		if err != nil {
			return err
		}
		if exists {
			return service.ErrUserAlreadyExists
		}
	}

	return h.authService.SendOTP(ctx, service.SendOTPRequest{
		Phone:   phone,
		Purpose: cmd.Purpose,
	})
}

// HandleVerifyOTP verifies an OTP and logs in the user
func (h *AuthHandler) HandleVerifyOTP(ctx context.Context, cmd command.VerifyOTP) (*command.LoginResult, error) {
	phone, err := cmd.GetPhone()
	if err != nil {
		return nil, err
	}

	// Validate code format
	if len(cmd.Code) != 6 {
		return nil, errors.New("OTP must be 6 digits")
	}

	result, err := h.authService.VerifyOTPAndLogin(ctx, service.VerifyOTPRequest{
		Phone:   phone,
		Code:    cmd.Code,
		Purpose: cmd.Purpose,
	}, cmd.IPAddress, cmd.UserAgent)

	if err != nil {
		return nil, err
	}

	return &command.LoginResult{
		UserID:       result.User.ID().String(),
		Phone:        result.User.Phone().String(),
		FullName:     result.User.FullName().String(),
		Email:        result.User.Email().String(),
		Tier:         result.User.Tier().String(),
		IsVerified:   result.User.IsVerified(),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

// HandleRegister registers a new user
func (h *AuthHandler) HandleRegister(ctx context.Context, cmd command.RegisterUser) (*command.RegisterResult, error) {
	phone, err := cmd.GetPhone()
	if err != nil {
		return nil, err
	}

	fullName, err := cmd.GetFullName()
	if err != nil {
		return nil, err
	}

	email, err := cmd.GetEmail()
	if err != nil {
		return nil, err
	}

	result, err := h.registrationService.Register(ctx, service.RegistrationRequest{
		Phone:        phone,
		FullName:     fullName,
		Email:        email,
		ReferralCode: cmd.ReferralCode,
	})

	if err != nil {
		return nil, err
	}

	return &command.RegisterResult{
		UserID:       result.User.ID().String(),
		Phone:        result.User.Phone().String(),
		FullName:     result.User.FullName().String(),
		Email:        result.User.Email().String(),
		Tier:         result.User.Tier().String(),
		ReferralCode: result.User.ReferralCode(),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

// HandleRefreshToken refreshes authentication tokens
func (h *AuthHandler) HandleRefreshToken(ctx context.Context, cmd command.RefreshToken) (*command.LoginResult, error) {
	if cmd.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	result, err := h.authService.RefreshTokens(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &command.LoginResult{
		UserID:       result.User.ID().String(),
		Phone:        result.User.Phone().String(),
		FullName:     result.User.FullName().String(),
		Email:        result.User.Email().String(),
		Tier:         result.User.Tier().String(),
		IsVerified:   result.User.IsVerified(),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresAt:    result.ExpiresAt,
	}, nil
}

// HandleLogout logs out a user
func (h *AuthHandler) HandleLogout(ctx context.Context, cmd command.Logout) error {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	return h.authService.Logout(ctx, userID)
}

// ProfileHandler handles profile-related commands
type ProfileHandler struct {
	userRepo  repository.UserRepository
	skillRepo repository.SkillRepository
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(
	userRepo repository.UserRepository,
	skillRepo repository.SkillRepository,
) *ProfileHandler {
	return &ProfileHandler{
		userRepo:  userRepo,
		skillRepo: skillRepo,
	}
}

// HandleUpdateProfile updates a user's profile
func (h *ProfileHandler) HandleUpdateProfile(ctx context.Context, cmd command.UpdateProfile) (*command.UpdateProfileResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	// Update profile
	err = user.UpdateProfile(
		cmd.Username,
		cmd.Bio,
		cmd.Location,
		cmd.State,
		cmd.DateOfBirth,
		cmd.Gender,
		cmd.ProfileImage,
	)
	if err != nil {
		return nil, err
	}

	// Save changes
	if err := h.userRepo.SaveWithEvents(ctx, user); err != nil {
		return nil, err
	}

	return &command.UpdateProfileResult{
		UserID:       user.ID().String(),
		Phone:        user.Phone().String(),
		Email:        user.Email().String(),
		FullName:     user.FullName().String(),
		Username:     user.Username(),
		Bio:          user.Bio(),
		Location:     user.Location(),
		State:        user.State(),
		ProfileImage: user.ProfileImage(),
		UpdatedAt:    user.UpdatedAt(),
	}, nil
}

// HandleAddSkill adds a skill to a user's profile
func (h *ProfileHandler) HandleAddSkill(ctx context.Context, cmd command.AddUserSkill) (*command.AddSkillResult, error) {
	userID, err := cmd.GetUserID()
	if err != nil {
		return nil, err
	}

	skillID, err := cmd.GetSkillID()
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	// Get skill info
	skill, err := h.skillRepo.FindByID(ctx, skillID)
	if err != nil {
		return nil, errors.New("skill not found")
	}

	// Add skill to user
	proficiency := aggregate.Proficiency(cmd.Proficiency)
	if err := user.AddSkill(skillID, skill.Name, proficiency, cmd.YearsExp); err != nil {
		return nil, err
	}

	// Save changes
	if err := h.userRepo.SaveWithEvents(ctx, user); err != nil {
		return nil, err
	}

	return &command.AddSkillResult{
		UserID:    user.ID().String(),
		SkillID:   skillID.String(),
		SkillName: skill.Name,
	}, nil
}

// HandleRemoveSkill removes a skill from a user's profile
func (h *ProfileHandler) HandleRemoveSkill(ctx context.Context, cmd command.RemoveUserSkill) error {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	skillID, err := valueobject.NewSkillID(cmd.SkillID)
	if err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return service.ErrUserNotFound
	}

	if err := user.RemoveSkill(skillID); err != nil {
		return err
	}

	return h.userRepo.SaveWithEvents(ctx, user)
}

// AdminHandler handles admin-related commands
type AdminHandler struct {
	userRepo repository.UserRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(userRepo repository.UserRepository) *AdminHandler {
	return &AdminHandler{userRepo: userRepo}
}

// HandleDeactivateUser deactivates a user account
func (h *AdminHandler) HandleDeactivateUser(ctx context.Context, cmd command.DeactivateUser) error {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return service.ErrUserNotFound
	}

	if err := user.Deactivate(cmd.Reason, cmd.RequestedBy); err != nil {
		return err
	}

	return h.userRepo.SaveWithEvents(ctx, user)
}

// HandleReactivateUser reactivates a user account
func (h *AdminHandler) HandleReactivateUser(ctx context.Context, cmd command.ReactivateUser) error {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return service.ErrUserNotFound
	}

	if err := user.Reactivate(cmd.RequestedBy); err != nil {
		return err
	}

	return h.userRepo.SaveWithEvents(ctx, user)
}

// HandleUpgradeUserTier upgrades a user's tier
func (h *AdminHandler) HandleUpgradeUserTier(ctx context.Context, cmd command.UpgradeUserTier) error {
	userID, err := valueobject.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return service.ErrUserNotFound
	}

	newTier := aggregate.UserTier(cmd.NewTier)
	if err := user.UpgradeTier(newTier, cmd.Reason); err != nil {
		return err
	}

	return h.userRepo.SaveWithEvents(ctx, user)
}
