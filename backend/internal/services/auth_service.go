package services

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hustlex/backend/internal/config"
	"github.com/hustlex/backend/internal/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.JWTConfig
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, redis *redis.Client, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		db:     db,
		redis:  redis,
		config: cfg,
	}
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Phone  string    `json:"phone"`
	Tier   string    `json:"tier"`
	jwt.RegisteredClaims
}

// RefreshClaims represents refresh token claims
type RefreshClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

// RegisterInput represents registration input
type RegisterInput struct {
	Phone    string `json:"phone" validate:"required,e164"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Referral string `json:"referral_code" validate:"omitempty,len=8"`
}

// LoginInput represents login input
type LoginInput struct {
	Phone string `json:"phone" validate:"required,e164"`
}

// VerifyOTPInput represents OTP verification input
type VerifyOTPInput struct {
	Phone   string `json:"phone" validate:"required,e164"`
	Code    string `json:"code" validate:"required,len=6"`
	Purpose string `json:"purpose" validate:"required,oneof=login register reset_pin"`
}

// SendOTP generates and "sends" an OTP (in production, integrate with SMS gateway)
func (s *AuthService) SendOTP(ctx context.Context, phone, purpose string) error {
	// Check rate limiting
	rateLimitKey := fmt.Sprintf("otp_rate:%s", phone)
	count, _ := s.redis.Incr(ctx, rateLimitKey).Result()
	if count == 1 {
		s.redis.Expire(ctx, rateLimitKey, 1*time.Hour)
	}
	if count > 5 {
		return errors.New("too many OTP requests, please try again later")
	}

	// Generate 6-digit OTP
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Store OTP in database
	otp := &models.OTPCode{
		Phone:     phone,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
	}

	// Delete existing unused OTPs for this phone/purpose
	s.db.Where("phone = ? AND purpose = ? AND is_used = ?", phone, purpose, false).Delete(&models.OTPCode{})

	if err := s.db.Create(otp).Error; err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	// Send OTP via SMS gateway (Termii, Africa's Talking, etc.)
	// In development mode only, log to a secure audit log (never to stdout)
	if os.Getenv("ENVIRONMENT") == "development" {
		log.Printf("[DEV-ONLY] OTP generated for phone ending in %s", phone[len(phone)-4:])
	}

	// TODO: Integrate SMS gateway here
	// Example: s.smsService.SendOTP(phone, code, purpose)

	return nil
}

// VerifyOTP verifies an OTP code
func (s *AuthService) VerifyOTP(ctx context.Context, input *VerifyOTPInput) (*models.User, error) {
	var otp models.OTPCode
	
	err := s.db.Where(
		"phone = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
		input.Phone, input.Purpose, false, time.Now(),
	).Order("created_at DESC").First(&otp).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired OTP")
		}
		return nil, err
	}

	// Check attempts
	if otp.Attempts >= 3 {
		return nil, errors.New("too many invalid attempts, please request a new OTP")
	}

	// Verify code using constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(otp.Code), []byte(input.Code)) != 1 {
		s.db.Model(&otp).Update("attempts", otp.Attempts+1)
		return nil, errors.New("invalid OTP code")
	}

	// Mark as used
	s.db.Model(&otp).Update("is_used", true)

	// Get or create user
	var user models.User
	err = s.db.Where("phone = ?", input.Phone).First(&user).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if input.Purpose == "login" {
			return nil, errors.New("user not found, please register first")
		}
		// Will be created in Register flow
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	// Update last login
	s.db.Model(&user).Update("last_login_at", time.Now())

	return &user, nil
}

// Register creates a new user
func (s *AuthService) Register(ctx context.Context, input *RegisterInput) (*models.User, error) {
	// Check if user already exists
	var existing models.User
	if err := s.db.Where("phone = ?", input.Phone).First(&existing).Error; err == nil {
		return nil, errors.New("user with this phone already exists")
	}

	// Generate referral code
	referralCode, _ := generateReferralCode()

	// Create user
	user := &models.User{
		Phone:        input.Phone,
		FullName:     input.FullName,
		Email:        input.Email,
		ReferralCode: referralCode,
		Tier:         models.TierBronze,
		IsActive:     true,
	}

	// Handle referral
	if input.Referral != "" {
		var referrer models.User
		if err := s.db.Where("referral_code = ?", input.Referral).First(&referrer).Error; err == nil {
			user.ReferredBy = &referrer.ID
		}
	}

	// Start transaction
	tx := s.db.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create wallet for user
	wallet := &models.Wallet{
		UserID:   user.ID,
		Balance:  0,
		Currency: "NGN",
	}
	if err := tx.Create(wallet).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Create credit score record
	creditScore := &models.CreditScore{
		UserID:           user.ID,
		Score:            100, // Starting score
		Tier:             models.TierBronze,
		LastCalculatedAt: time.Now(),
	}
	if err := tx.Create(creditScore).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create credit score: %w", err)
	}

	tx.Commit()

	// Load relationships
	s.db.Preload("Wallet").Preload("CreditScore").First(&user, user.ID)

	return user, nil
}

// GenerateTokens generates access and refresh tokens
func (s *AuthService) GenerateTokens(user *models.User) (*TokenPair, error) {
	now := time.Now()
	accessExpiry := now.Add(s.config.AccessExpiry)
	refreshExpiry := now.Add(s.config.RefreshExpiry)

	// Access token claims
	accessClaims := Claims{
		UserID: user.ID,
		Phone:  user.Phone,
		Tier:   string(user.Tier),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
			Subject:   user.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token claims
	refreshClaims := RefreshClaims{
		UserID:    user.ID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    s.config.Issuer,
			Subject:   user.ID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	// Store refresh token in Redis for validation
	ctx := context.Background()
	refreshKey := fmt.Sprintf("refresh:%s", user.ID.String())
	s.redis.Set(ctx, refreshKey, refreshTokenString, s.config.RefreshExpiry)

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry,
		TokenType:    "Bearer",
	}, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshTokens refreshes tokens using a refresh token
func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenString string) (*TokenPair, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid || claims.TokenType != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	// Verify token exists in Redis
	refreshKey := fmt.Sprintf("refresh:%s", claims.UserID.String())
	storedToken, err := s.redis.Get(ctx, refreshKey).Result()
	if err != nil || storedToken != refreshTokenString {
		return nil, errors.New("refresh token revoked or expired")
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Delete old refresh token
	s.redis.Del(ctx, refreshKey)

	// Generate new tokens
	return s.GenerateTokens(&user)
}

// Logout invalidates refresh token
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	refreshKey := fmt.Sprintf("refresh:%s", userID.String())
	return s.redis.Del(ctx, refreshKey).Err()
}

// SetTransactionPIN sets/updates user's transaction PIN
func (s *AuthService) SetTransactionPIN(userID uuid.UUID, pin string) error {
	if len(pin) != 4 && len(pin) != 6 {
		return errors.New("PIN must be 4 or 6 digits")
	}

	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&models.Wallet{}).Where("user_id = ?", userID).Update("pin", string(hashedPIN)).Error
}

// VerifyTransactionPIN verifies user's transaction PIN
func (s *AuthService) VerifyTransactionPIN(userID uuid.UUID, pin string) error {
	var wallet models.Wallet
	if err := s.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return err
	}

	if wallet.Pin == "" {
		return errors.New("transaction PIN not set")
	}

	if wallet.PinAttempts >= 5 {
		return errors.New("wallet locked due to too many invalid attempts")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(wallet.Pin), []byte(pin)); err != nil {
		s.db.Model(&wallet).Update("pin_attempts", wallet.PinAttempts+1)
		return errors.New("invalid PIN")
	}

	// Reset attempts on success
	s.db.Model(&wallet).Update("pin_attempts", 0)
	return nil
}

// Helper functions

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}

func generateReferralCode() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 8)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
