package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"hustlex/internal/config"
	"hustlex/internal/infrastructure/auth"
	"hustlex/internal/infrastructure/security/audit"
	"hustlex/internal/infrastructure/security/ratelimit"
	"hustlex/internal/interface/http/middleware"
	"hustlex/internal/interface/http/router"

	"github.com/redis/go-redis/v9"
)

// Build-time variables
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	printBanner()
	log.Printf("Starting HustleX API %s (commit: %s, built: %s)", Version, Commit, BuildTime)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Redis client
	redisPort, _ := strconv.Atoi(cfg.Redis.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, redisPort),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = redisClient.Ping(ctx).Err()
	cancel()

	var useRedis bool
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (continuing without cache)", err)
		useRedis = false
	} else {
		defer redisClient.Close()
		log.Println("Connected to Redis/DragonflyDB")
		useRedis = true
	}

	// Initialize audit logger
	auditLogger := audit.NewInMemoryAuditLogger("hustlex-api")

	// Initialize rate limiters
	var authLimiter, txnLimiter, otpLimiter, pinLimiter ratelimit.RateLimiter
	if useRedis {
		authLimiter = ratelimit.NewRedisRateLimiter(redisClient, ratelimit.RateLimitAuth, "auth")
		txnLimiter = ratelimit.NewRedisRateLimiter(redisClient, ratelimit.RateLimitTransaction, "txn")
		otpLimiter = ratelimit.NewRedisRateLimiter(redisClient, ratelimit.RateLimitOTP, "otp")
		pinLimiter = ratelimit.NewRedisRateLimiter(redisClient, ratelimit.RateLimitPIN, "pin")
	} else {
		authLimiter = ratelimit.NewInMemoryRateLimiter(ratelimit.RateLimitAuth, "auth")
		txnLimiter = ratelimit.NewInMemoryRateLimiter(ratelimit.RateLimitTransaction, "txn")
		otpLimiter = ratelimit.NewInMemoryRateLimiter(ratelimit.RateLimitOTP, "otp")
		pinLimiter = ratelimit.NewInMemoryRateLimiter(ratelimit.RateLimitPIN, "pin")
	}

	// Initialize JWT validator and auth middleware
	jwtValidator := auth.NewJWTValidator(cfg.JWT.Secret, cfg.JWT.Issuer)
	authMiddleware := middleware.NewAuthMiddleware(jwtValidator)

	// Create router config
	routerConfig := router.Config{
		AllowedOrigins:  parseOrigins(cfg.Server.AllowOrigins),
		RequestTimeout:  cfg.Server.ReadTimeout,
		AuditLogger:     auditLogger,
		AuthRateLimiter: authLimiter,
		TxnRateLimiter:  txnLimiter,
		OTPRateLimiter:  otpLimiter,
		PINRateLimiter:  pinLimiter,
	}

	// Initialize handlers (nil for now - endpoints will return not implemented)
	handlers := router.Handlers{
		Wallet: nil, // Will be initialized when database is connected
	}

	r := router.NewRouter(routerConfig, handlers, authMiddleware)
	httpHandler := r.Setup()

	// Create server
	port := cfg.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		log.Printf("Environment: %s", cfg.Server.Environment)
		log.Printf("Health check: http://localhost:%s/health", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func parseOrigins(origins string) []string {
	if origins == "" {
		return []string{"http://localhost:3000"}
	}
	var result []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func printBanner() {
	banner := `
  _   _           _   _     __  __
 | | | |_   _ ___| |_| | ___\ \/ /
 | |_| | | | / __| __| |/ _ \\  /
 |  _  | |_| \__ \ |_| |  __//  \
 |_| |_|\__,_|___/\__|_|\___/_/\_\

`
	fmt.Print(banner)
}
