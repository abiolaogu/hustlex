package router

import (
	"net/http"
	"time"

	"hustlex/internal/infrastructure/security/audit"
	"hustlex/internal/infrastructure/security/ratelimit"
	"hustlex/internal/interface/http/handler"
	"hustlex/internal/interface/http/middleware"
)

// Config holds router configuration
type Config struct {
	AllowedOrigins   []string
	AllowCredentials bool
	RequestTimeout   time.Duration
	RateLimiter      middleware.RateLimiter
	AuditLogger      audit.AuditLogger
	AuthRateLimiter  ratelimit.RateLimiter
	TxnRateLimiter   ratelimit.RateLimiter
	OTPRateLimiter   ratelimit.RateLimiter
	PINRateLimiter   ratelimit.RateLimiter
}

// Handlers holds all HTTP handlers
type Handlers struct {
	Wallet       *handler.WalletHandler
	// Auth         *handler.AuthHandler
	// Gig          *handler.GigHandler
	// Circle       *handler.CircleHandler
	// Credit       *handler.CreditHandler
	// Notification *handler.NotificationHandler
}

// Router sets up all application routes
type Router struct {
	mux      *http.ServeMux
	config   Config
	handlers Handlers
	auth     *middleware.AuthMiddleware
}

// NewRouter creates a new router
func NewRouter(config Config, handlers Handlers, auth *middleware.AuthMiddleware) *Router {
	return &Router{
		mux:      http.NewServeMux(),
		config:   config,
		handlers: handlers,
		auth:     auth,
	}
}

// Setup configures all routes
func (r *Router) Setup() http.Handler {
	// Apply global middleware
	var handler http.Handler = r.mux

	// Recovery from panics (enhanced with proper error handling)
	handler = middleware.RecoverPanic(handler)

	// Logging
	handler = middleware.Logger(handler)

	// Enhanced Request ID with correlation ID support
	handler = middleware.EnhancedRequestID(handler)

	// Request sanitization (blocks suspicious paths, validates content-type)
	handler = middleware.RequestSanitizer(handler)

	// Enhanced security headers (OWASP recommended)
	handler = middleware.EnhancedSecurityHeaders(handler)

	// Secure CORS with origin validation
	handler = middleware.SecureCORS(r.config.AllowedOrigins)(handler)

	// Secure JSON response handling
	handler = middleware.SecureJSONResponse(handler)

	// Content type
	handler = middleware.ContentType(handler)

	// Request timeout
	if r.config.RequestTimeout > 0 {
		handler = middleware.Timeout(r.config.RequestTimeout)(handler)
	}

	// Audit logging for compliance (if configured)
	if r.config.AuditLogger != nil {
		handler = middleware.AuditMiddleware(r.config.AuditLogger)(handler)
	}

	// Setup routes
	r.setupHealthRoutes()
	r.setupAuthRoutes()
	r.setupWalletRoutes()
	r.setupGigRoutes()
	r.setupCircleRoutes()
	r.setupCreditRoutes()
	r.setupNotificationRoutes()
	r.setupAdminRoutes()

	return handler
}

// setupHealthRoutes configures health check routes
func (r *Router) setupHealthRoutes() {
	r.mux.HandleFunc("GET /health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.mux.HandleFunc("GET /ready", func(w http.ResponseWriter, req *http.Request) {
		// TODO: Add database and service health checks
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})
}

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes() {
	// Public routes (no auth required) - with rate limiting
	// OTP endpoints have strict rate limits to prevent abuse
	r.mux.HandleFunc("POST /api/auth/otp/send", r.rateLimitedPublicHandler(r.config.OTPRateLimiter, notImplemented))
	r.mux.HandleFunc("POST /api/auth/otp/verify", r.rateLimitedPublicHandler(r.config.OTPRateLimiter, notImplemented))

	// Auth endpoints with standard auth rate limits
	r.mux.HandleFunc("POST /api/auth/register", r.rateLimitedPublicHandler(r.config.AuthRateLimiter, notImplemented))
	r.mux.HandleFunc("POST /api/auth/login", r.rateLimitedPublicHandler(r.config.AuthRateLimiter, notImplemented))
	r.mux.HandleFunc("POST /api/auth/refresh", r.rateLimitedPublicHandler(r.config.AuthRateLimiter, notImplemented))

	// Protected routes
	r.mux.HandleFunc("POST /api/auth/logout", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/auth/me", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/auth/profile", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/auth/password", r.protectedHandler(notImplemented))
}

// setupWalletRoutes configures wallet routes
func (r *Router) setupWalletRoutes() {
	if r.handlers.Wallet == nil {
		return
	}

	// Wallet routes (protected) - read operations
	r.mux.HandleFunc("GET /api/wallet", r.protectedHandler(r.handlers.Wallet.GetWallet))
	r.mux.HandleFunc("GET /api/wallet/balance", r.protectedHandler(r.handlers.Wallet.GetBalance))
	r.mux.HandleFunc("GET /api/wallet/transactions", r.protectedHandler(r.handlers.Wallet.GetTransactions))

	// Deposit routes - with transaction rate limiting
	r.mux.HandleFunc("POST /api/wallet/deposit/initiate", r.rateLimitedProtectedHandler(r.config.TxnRateLimiter, r.handlers.Wallet.InitiateDeposit))
	r.mux.HandleFunc("POST /api/wallet/deposit/verify", r.protectedHandler(r.handlers.Wallet.VerifyDeposit))

	// Withdrawal routes - with transaction rate limiting
	r.mux.HandleFunc("POST /api/wallet/withdraw", r.rateLimitedProtectedHandler(r.config.TxnRateLimiter, r.handlers.Wallet.InitiateWithdraw))

	// Transfer routes - with transaction rate limiting and PIN rate limiting
	r.mux.HandleFunc("POST /api/wallet/transfer", r.rateLimitedProtectedHandler(r.config.TxnRateLimiter, r.handlers.Wallet.Transfer))

	// Bank routes
	r.mux.HandleFunc("GET /api/wallet/banks", r.protectedHandler(r.handlers.Wallet.GetBanks))
	r.mux.HandleFunc("POST /api/wallet/resolve-account", r.protectedHandler(r.handlers.Wallet.ResolveAccount))

	// Webhook (public but validated via signature)
	r.mux.HandleFunc("POST /api/webhook/paystack", r.publicHandler(notImplemented))
}

// setupGigRoutes configures gig routes
func (r *Router) setupGigRoutes() {
	// Public gig listing
	r.mux.HandleFunc("GET /api/gigs", r.optionalAuthHandler(notImplemented))
	r.mux.HandleFunc("GET /api/gigs/{id}", r.optionalAuthHandler(notImplemented))

	// Protected gig routes
	r.mux.HandleFunc("POST /api/gigs", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/gigs/{id}", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("DELETE /api/gigs/{id}", r.protectedHandler(notImplemented))

	// Proposals
	r.mux.HandleFunc("POST /api/gigs/{id}/proposals", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/gigs/{id}/proposals", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/proposals/{id}/accept", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/proposals/{id}/reject", r.protectedHandler(notImplemented))

	// Contracts
	r.mux.HandleFunc("GET /api/contracts", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/contracts/{id}", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/contracts/{id}/deliver", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/contracts/{id}/accept", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/contracts/{id}/dispute", r.protectedHandler(notImplemented))

	// Reviews
	r.mux.HandleFunc("POST /api/contracts/{id}/review", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/users/{id}/reviews", r.publicHandler(notImplemented))

	// My gigs and contracts
	r.mux.HandleFunc("GET /api/me/gigs", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/me/contracts", r.protectedHandler(notImplemented))
}

// setupCircleRoutes configures savings circle routes
func (r *Router) setupCircleRoutes() {
	// Public circle listing
	r.mux.HandleFunc("GET /api/circles", r.optionalAuthHandler(notImplemented))
	r.mux.HandleFunc("GET /api/circles/{id}", r.optionalAuthHandler(notImplemented))

	// Protected circle routes
	r.mux.HandleFunc("POST /api/circles", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/circles/{id}", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/circles/{id}/join", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/circles/{id}/leave", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/circles/{id}/start", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/circles/join-by-code", r.protectedHandler(notImplemented))

	// Contributions
	r.mux.HandleFunc("POST /api/circles/{id}/contribute", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/circles/{id}/contributions", r.protectedHandler(notImplemented))

	// My circles
	r.mux.HandleFunc("GET /api/me/circles", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/me/circles/stats", r.protectedHandler(notImplemented))
}

// setupCreditRoutes configures credit and loan routes
func (r *Router) setupCreditRoutes() {
	// Credit score
	r.mux.HandleFunc("GET /api/credit/score", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/credit/recalculate", r.protectedHandler(notImplemented))

	// Loans
	r.mux.HandleFunc("GET /api/loans", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/loans/{id}", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/loans/apply", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("POST /api/loans/{id}/repay", r.protectedHandler(notImplemented))

	// Loan stats
	r.mux.HandleFunc("GET /api/me/loan-stats", r.protectedHandler(notImplemented))
}

// setupNotificationRoutes configures notification routes
func (r *Router) setupNotificationRoutes() {
	// Notifications
	r.mux.HandleFunc("GET /api/notifications", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("GET /api/notifications/unread-count", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/notifications/{id}/read", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/notifications/mark-all-read", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("DELETE /api/notifications/{id}", r.protectedHandler(notImplemented))

	// Preferences
	r.mux.HandleFunc("GET /api/notifications/preferences", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("PUT /api/notifications/preferences", r.protectedHandler(notImplemented))

	// Device tokens
	r.mux.HandleFunc("POST /api/notifications/device-token", r.protectedHandler(notImplemented))
	r.mux.HandleFunc("DELETE /api/notifications/device-token", r.protectedHandler(notImplemented))
}

// setupAdminRoutes configures admin routes
func (r *Router) setupAdminRoutes() {
	// Admin requires authentication and admin role
	adminMiddleware := func(h http.HandlerFunc) http.HandlerFunc {
		return r.protectedHandler(func(w http.ResponseWriter, req *http.Request) {
			if !middleware.HasRole(req.Context(), "admin") {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			h(w, req)
		})
	}

	// User management
	r.mux.HandleFunc("GET /api/admin/users", adminMiddleware(notImplemented))
	r.mux.HandleFunc("GET /api/admin/users/{id}", adminMiddleware(notImplemented))
	r.mux.HandleFunc("PUT /api/admin/users/{id}/status", adminMiddleware(notImplemented))

	// Loan management
	r.mux.HandleFunc("GET /api/admin/loans", adminMiddleware(notImplemented))
	r.mux.HandleFunc("GET /api/admin/loans/overdue", adminMiddleware(notImplemented))
	r.mux.HandleFunc("POST /api/admin/loans/{id}/approve", adminMiddleware(notImplemented))
	r.mux.HandleFunc("POST /api/admin/loans/{id}/reject", adminMiddleware(notImplemented))
	r.mux.HandleFunc("POST /api/admin/loans/{id}/disburse", adminMiddleware(notImplemented))
	r.mux.HandleFunc("POST /api/admin/loans/{id}/default", adminMiddleware(notImplemented))

	// Circle management
	r.mux.HandleFunc("GET /api/admin/circles", adminMiddleware(notImplemented))

	// Statistics
	r.mux.HandleFunc("GET /api/admin/stats/overview", adminMiddleware(notImplemented))
	r.mux.HandleFunc("GET /api/admin/stats/loans", adminMiddleware(notImplemented))
	r.mux.HandleFunc("GET /api/admin/stats/transactions", adminMiddleware(notImplemented))
}

// Helper methods for applying middleware

func (r *Router) publicHandler(h http.HandlerFunc) http.HandlerFunc {
	return h
}

func (r *Router) protectedHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.auth.Authenticate(http.HandlerFunc(h)).ServeHTTP(w, req)
	}
}

func (r *Router) optionalAuthHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.auth.OptionalAuth(http.HandlerFunc(h)).ServeHTTP(w, req)
	}
}

// rateLimitedPublicHandler applies rate limiting to public endpoints
func (r *Router) rateLimitedPublicHandler(limiter ratelimit.RateLimiter, h http.HandlerFunc) http.HandlerFunc {
	if limiter == nil {
		return h
	}
	return func(w http.ResponseWriter, req *http.Request) {
		ratelimit.RateLimitMiddleware(limiter, ratelimit.IPKeyFunc)(http.HandlerFunc(h)).ServeHTTP(w, req)
	}
}

// rateLimitedProtectedHandler applies rate limiting to protected endpoints
func (r *Router) rateLimitedProtectedHandler(limiter ratelimit.RateLimiter, h http.HandlerFunc) http.HandlerFunc {
	if limiter == nil {
		return r.protectedHandler(h)
	}
	return func(w http.ResponseWriter, req *http.Request) {
		// First authenticate, then apply rate limit
		authenticated := func(w http.ResponseWriter, req *http.Request) {
			ratelimit.RateLimitMiddleware(limiter, ratelimit.IPKeyFunc)(http.HandlerFunc(h)).ServeHTTP(w, req)
		}
		r.auth.Authenticate(http.HandlerFunc(authenticated)).ServeHTTP(w, req)
	}
}

// notImplemented is a placeholder for unimplemented handlers
func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"not_implemented","message":"this endpoint is not yet implemented"}`))
}
