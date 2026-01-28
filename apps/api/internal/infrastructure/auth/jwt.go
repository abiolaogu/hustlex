package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hustlex/internal/interface/http/middleware"
)

// JWTValidator implements middleware.TokenValidator
type JWTValidator struct {
	secret []byte
	issuer string
}

// NewJWTValidator creates a new JWT validator
func NewJWTValidator(secret, issuer string) *JWTValidator {
	return &JWTValidator{
		secret: []byte(secret),
		issuer: issuer,
	}
}

// Claims represents custom JWT claims
type Claims struct {
	UserID    string   `json:"user_id"`
	SessionID string   `json:"session_id"`
	Roles     []string `json:"roles"`
	jwt.RegisteredClaims
}

// ValidateToken validates a JWT token and returns the claims
func (v *JWTValidator) ValidateToken(tokenString string) (*middleware.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check issuer if configured
	if v.issuer != "" && claims.Issuer != v.issuer {
		return nil, errors.New("invalid token issuer")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &middleware.TokenClaims{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		Roles:     claims.Roles,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

// GenerateToken creates a new JWT token
func (v *JWTValidator) GenerateToken(userID, sessionID string, roles []string, expiresIn time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		SessionID: sessionID,
		Roles:     roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    v.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(v.secret)
}
