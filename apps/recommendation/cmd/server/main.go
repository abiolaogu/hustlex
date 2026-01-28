package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	recommendation "recommendation"
	"recommendation/api"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting HustleX Recommendation Service",
		"version", Version,
		"commit", Commit,
		"build_time", BuildTime,
	)

	// Load configuration
	config := loadConfig()

	// Connect to PostgreSQL
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("Connected to PostgreSQL")

	// Connect to Redis
	redisOpts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		logger.Error("Failed to parse Redis URL", "error", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, continuing without cache", "error", err)
	} else {
		logger.Info("Connected to Redis/DragonflyDB")
	}

	// Create recommendation engine
	engineConfig := recommendation.DefaultConfig()
	engine, err := recommendation.NewEngine(dbPool, redisClient, engineConfig)
	if err != nil {
		logger.Error("Failed to create recommendation engine", "error", err)
		os.Exit(1)
	}
	logger.Info("Recommendation engine initialized")

	// Create API server
	server := api.NewServer(engine, logger)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      server,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	go func() {
		logger.Info("Server listening", "port", config.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited")
}

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	Environment string
}

func loadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://hustlex:hustlex_dev_password@localhost:5432/hustlex?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/1"),
		Environment: getEnv("ENV", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
