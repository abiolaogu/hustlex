package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMS      SMSConfig
	Payment  PaymentConfig
	Storage  StorageConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Environment  string
	AllowOrigins string
	RateLimit    int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	Issuer           string
}

// SMSConfig holds SMS gateway configuration
type SMSConfig struct {
	Provider  string // termii, africaistalking, twilio
	APIKey    string
	SecretKey string
	SenderID  string
}

// PaymentConfig holds payment gateway configuration
type PaymentConfig struct {
	Provider      string // paystack, flutterwave
	PublicKey     string
	SecretKey     string
	WebhookSecret string
}

// StorageConfig holds file storage configuration
type StorageConfig struct {
	Provider        string // s3, minio, local
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (for development)
	_ = godotenv.Load()

	environment := getEnv("ENVIRONMENT", "development")
	isProduction := environment == "production"

	// Validate required secrets in production
	jwtSecret := os.Getenv("JWT_SECRET")
	if isProduction && jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required in production")
	}
	if jwtSecret == "" {
		jwtSecret = "dev-only-secret-do-not-use-in-production"
	}

	// Validate CORS origins - never allow wildcard in production
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:8080")
	if isProduction && (corsOrigins == "*" || strings.Contains(corsOrigins, "*")) {
		return nil, fmt.Errorf("CORS_ORIGINS cannot be wildcard (*) in production")
	}

	// Validate database password in production
	dbPassword := os.Getenv("DB_PASSWORD")
	if isProduction && dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required in production")
	}
	if dbPassword == "" {
		dbPassword = "hustlex_password"
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Environment:  environment,
			AllowOrigins: corsOrigins,
			RateLimit:    getEnvInt("RATE_LIMIT", 100),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "hustlex"),
			Password:     dbPassword,
			DBName:       getEnv("DB_NAME", "hustlex"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:        jwtSecret,
			AccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			Issuer:        getEnv("JWT_ISSUER", "hustlex"),
		},
		SMS: SMSConfig{
			Provider:  getEnv("SMS_PROVIDER", "termii"),
			APIKey:    getEnv("SMS_API_KEY", ""),
			SecretKey: getEnv("SMS_SECRET_KEY", ""),
			SenderID:  getEnv("SMS_SENDER_ID", "HustleX"),
		},
		Payment: PaymentConfig{
			Provider:      getEnv("PAYMENT_PROVIDER", "paystack"),
			PublicKey:     getEnv("PAYMENT_PUBLIC_KEY", ""),
			SecretKey:     getEnv("PAYMENT_SECRET_KEY", ""),
			WebhookSecret: getEnv("PAYMENT_WEBHOOK_SECRET", ""),
		},
		Storage: StorageConfig{
			Provider:        getEnv("STORAGE_PROVIDER", "local"),
			Endpoint:        getEnv("STORAGE_ENDPOINT", ""),
			Bucket:          getEnv("STORAGE_BUCKET", "hustlex-files"),
			AccessKeyID:     getEnv("STORAGE_ACCESS_KEY", ""),
			SecretAccessKey: getEnv("STORAGE_SECRET_KEY", ""),
			Region:          getEnv("STORAGE_REGION", "us-east-1"),
		},
	}

	return cfg, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
