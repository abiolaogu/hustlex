package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"hustlex/internal/config"
	"hustlex/internal/services"
)

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Phone  string    `json:"phone"`
	Tier   string    `json:"tier"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config, redis *redis.Client) fiber.Handler {
	blacklist := services.NewTokenBlacklistService(redis)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Missing authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid authorization format",
			})
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Invalid token claims",
			})
		}

		// Check expiry
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Token expired",
			})
		}

		// Check if token is blacklisted
		ctx := context.Background()
		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, tokenString)
		if err != nil {
			// Log error but don't fail - fail open for availability
			// In production, consider failing closed for security
		}
		if isBlacklisted {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Token has been revoked",
			})
		}

		// Check if all user tokens have been revoked
		if claims.IssuedAt != nil {
			isRevoked, err := blacklist.IsUserTokenRevoked(ctx, claims.UserID.String(), claims.IssuedAt.Time)
			if err != nil {
				// Log error but don't fail
			}
			if isRevoked {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"error":   "Token has been revoked due to security action",
				})
			}
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("phone", claims.Phone)
		c.Locals("tier", claims.Tier)

		return c.Next()
	}
}

// OptionalAuthMiddleware checks for token but doesn't require it
func OptionalAuthMiddleware(cfg *config.Config, redis *redis.Client) fiber.Handler {
	blacklist := services.NewTokenBlacklistService(redis)

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			return c.Next()
		}

		if claims, ok := token.Claims.(*Claims); ok {
			// Check blacklist even for optional auth
			ctx := context.Background()
			isBlacklisted, _ := blacklist.IsTokenBlacklisted(ctx, tokenString)
			if isBlacklisted {
				return c.Next() // Just skip auth, don't fail
			}

			// Check user-level revocation
			if claims.IssuedAt != nil {
				isRevoked, _ := blacklist.IsUserTokenRevoked(ctx, claims.UserID.String(), claims.IssuedAt.Time)
				if isRevoked {
					return c.Next() // Just skip auth, don't fail
				}
			}

			c.Locals("userID", claims.UserID)
			c.Locals("phone", claims.Phone)
			c.Locals("tier", claims.Tier)
		}

		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	return userID, nil
}

// RequireTier middleware checks if user has required tier or higher
func RequireTier(requiredTiers ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tier, ok := c.Locals("tier").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "User tier not found",
			})
		}

		for _, required := range requiredTiers {
			if tier == required {
				return c.Next()
			}
		}

		// Check tier hierarchy
		tierOrder := map[string]int{
			"bronze":   1,
			"silver":   2,
			"gold":     3,
			"platinum": 4,
		}

		userTierLevel := tierOrder[tier]
		for _, required := range requiredTiers {
			if userTierLevel >= tierOrder[required] {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Insufficient tier level",
		})
	}
}
