package service

import (
	"context"
	"errors"
	"time"

	"hustlex/internal/domain/identity/aggregate"
	"hustlex/internal/domain/identity/repository"
	"hustlex/internal/domain/shared/valueobject"
)

// Domain service errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this phone already exists")
	ErrInvalidOTP         = errors.New("invalid or expired OTP")
	ErrOTPExpired         = errors.New("OTP has expired")
	ErrTooManyAttempts    = errors.New("too many invalid attempts")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded, please try again later")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account is locked")
	ErrInvalidReferer     = errors.New("invalid referral code")
)

// OTPGenerator defines the interface for generating OTPs
type OTPGenerator interface {
	Generate(length int) (string, error)
}

// OTPSender defines the interface for sending OTPs (SMS)
type OTPSender interface {
	Send(ctx context.Context, phone, code, purpose string) error
}

// TokenGenerator defines the interface for generating auth tokens
type TokenGenerator interface {
	GenerateAccessToken(user *aggregate.User) (string, time.Time, error)
	GenerateRefreshToken(userID valueobject.UserID) (string, time.Time, error)
	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*RefreshTokenClaims, error)
}

// TokenClaims represents the claims in an access token
type TokenClaims struct {
	UserID    valueobject.UserID
	Phone     string
	Tier      string
	ExpiresAt time.Time
}

// RefreshTokenClaims represents the claims in a refresh token
type RefreshTokenClaims struct {
	UserID    valueobject.UserID
	ExpiresAt time.Time
}

// RegistrationService handles user registration domain logic
type RegistrationService struct {
	userRepo       repository.UserRepository
	otpRepo        repository.OTPRepository
	sessionRepo    repository.SessionRepository
	otpGenerator   OTPGenerator
	otpSender      OTPSender
	tokenGenerator TokenGenerator
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	sessionRepo repository.SessionRepository,
	otpGenerator OTPGenerator,
	otpSender OTPSender,
	tokenGenerator TokenGenerator,
) *RegistrationService {
	return &RegistrationService{
		userRepo:       userRepo,
		otpRepo:        otpRepo,
		sessionRepo:    sessionRepo,
		otpGenerator:   otpGenerator,
		otpSender:      otpSender,
		tokenGenerator: tokenGenerator,
	}
}

// RegistrationRequest represents a registration request
type RegistrationRequest struct {
	Phone        valueobject.PhoneNumber
	FullName     valueobject.FullName
	Email        valueobject.Email
	ReferralCode string
}

// RegistrationResult represents the result of a successful registration
type RegistrationResult struct {
	User         *aggregate.User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// Register registers a new user after OTP verification
func (s *RegistrationService) Register(ctx context.Context, req RegistrationRequest) (*RegistrationResult, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Generate referral code for the new user
	referralCode, _ := s.otpGenerator.Generate(8)

	// Create user aggregate
	userID := valueobject.GenerateUserID()
	user, err := aggregate.NewUser(userID, req.Phone, req.FullName, referralCode)
	if err != nil {
		return nil, err
	}

	// Set email if provided
	if !req.Email.IsEmpty() {
		user.SetEmail(req.Email)
	}

	// Handle referral if provided
	if req.ReferralCode != "" {
		referrer, err := s.userRepo.FindByReferralCode(ctx, req.ReferralCode)
		if err == nil && referrer != nil {
			user.SetReferrer(referrer.ID())
		}
		// Silently ignore invalid referral codes
	}

	// Persist user with events
	if err := s.userRepo.SaveWithEvents(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, expiresAt, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiry, err := s.tokenGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := s.sessionRepo.StoreRefreshToken(ctx, user.ID(), refreshToken, time.Until(refreshExpiry)); err != nil {
		return nil, err
	}

	return &RegistrationResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// AuthenticationService handles authentication domain logic
type AuthenticationService struct {
	userRepo       repository.UserRepository
	otpRepo        repository.OTPRepository
	sessionRepo    repository.SessionRepository
	otpGenerator   OTPGenerator
	otpSender      OTPSender
	tokenGenerator TokenGenerator
}

// NewAuthenticationService creates a new authentication service
func NewAuthenticationService(
	userRepo repository.UserRepository,
	otpRepo repository.OTPRepository,
	sessionRepo repository.SessionRepository,
	otpGenerator OTPGenerator,
	otpSender OTPSender,
	tokenGenerator TokenGenerator,
) *AuthenticationService {
	return &AuthenticationService{
		userRepo:       userRepo,
		otpRepo:        otpRepo,
		sessionRepo:    sessionRepo,
		otpGenerator:   otpGenerator,
		otpSender:      otpSender,
		tokenGenerator: tokenGenerator,
	}
}

// SendOTPRequest represents a request to send OTP
type SendOTPRequest struct {
	Phone   valueobject.PhoneNumber
	Purpose string // login, register, reset_pin
}

// SendOTP generates and sends an OTP
func (s *AuthenticationService) SendOTP(ctx context.Context, req SendOTPRequest) error {
	// Check rate limit
	allowed, err := s.sessionRepo.CheckOTPRateLimit(ctx, req.Phone.String(), 5, time.Hour)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrRateLimitExceeded
	}

	// Delete any existing unused OTPs
	_ = s.otpRepo.DeleteUnused(ctx, req.Phone.String(), req.Purpose)

	// Generate OTP
	code, err := s.otpGenerator.Generate(6)
	if err != nil {
		return err
	}

	// Create OTP record
	otp := &repository.OTPCode{
		Phone:     req.Phone.String(),
		Code:      code,
		Purpose:   req.Purpose,
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
		IsUsed:    false,
		Attempts:  0,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.otpRepo.Save(ctx, otp); err != nil {
		return err
	}

	// Send OTP via SMS
	return s.otpSender.Send(ctx, req.Phone.String(), code, req.Purpose)
}

// VerifyOTPRequest represents a request to verify OTP
type VerifyOTPRequest struct {
	Phone   valueobject.PhoneNumber
	Code    string
	Purpose string
}

// LoginResult represents the result of a successful login
type LoginResult struct {
	User         *aggregate.User
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// VerifyOTPAndLogin verifies OTP and logs in the user
func (s *AuthenticationService) VerifyOTPAndLogin(ctx context.Context, req VerifyOTPRequest, ipAddress, userAgent string) (*LoginResult, error) {
	// Find valid OTP
	otp, err := s.otpRepo.FindLatestValid(ctx, req.Phone.String(), req.Purpose)
	if err != nil {
		return nil, ErrInvalidOTP
	}

	if otp == nil {
		return nil, ErrInvalidOTP
	}

	// Check if expired
	if time.Now().After(otp.ExpiresAt) {
		return nil, ErrOTPExpired
	}

	// Check attempts
	if otp.Attempts >= 3 {
		return nil, ErrTooManyAttempts
	}

	// Verify code (should use constant-time comparison)
	if otp.Code != req.Code {
		_ = s.otpRepo.IncrementAttempts(ctx, otp.ID)
		return nil, ErrInvalidOTP
	}

	// Mark OTP as used
	_ = s.otpRepo.MarkUsed(ctx, otp.ID)

	// Find user
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive() {
		return nil, ErrAccountLocked
	}

	// Record login
	user.RecordLogin(ipAddress, userAgent)
	if err := s.userRepo.SaveWithEvents(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, expiresAt, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpiry, err := s.tokenGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, err
	}

	// Store refresh token
	if err := s.sessionRepo.StoreRefreshToken(ctx, user.ID(), refreshToken, time.Until(refreshExpiry)); err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshTokens refreshes access and refresh tokens
func (s *AuthenticationService) RefreshTokens(ctx context.Context, refreshToken string) (*LoginResult, error) {
	// Validate refresh token
	claims, err := s.tokenGenerator.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify token exists in storage
	storedToken, err := s.sessionRepo.GetRefreshToken(ctx, claims.UserID)
	if err != nil || storedToken != refreshToken {
		return nil, ErrInvalidCredentials
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive() {
		return nil, ErrAccountLocked
	}

	// Delete old refresh token
	_ = s.sessionRepo.DeleteRefreshToken(ctx, claims.UserID)

	// Generate new tokens
	accessToken, expiresAt, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, refreshExpiry, err := s.tokenGenerator.GenerateRefreshToken(user.ID())
	if err != nil {
		return nil, err
	}

	// Store new refresh token
	if err := s.sessionRepo.StoreRefreshToken(ctx, user.ID(), newRefreshToken, time.Until(refreshExpiry)); err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalidates the user's session
func (s *AuthenticationService) Logout(ctx context.Context, userID valueobject.UserID) error {
	return s.sessionRepo.DeleteRefreshToken(ctx, userID)
}
