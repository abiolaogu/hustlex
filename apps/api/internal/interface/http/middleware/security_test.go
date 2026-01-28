package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnhancedSecurityHeaders(t *testing.T) {
	handler := EnhancedSecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	tests := []struct {
		header string
		want   string
	}{
		{"X-Frame-Options", "DENY"},
		{"X-XSS-Protection", "1; mode=block"},
		{"X-Content-Type-Options", "nosniff"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
		{"Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload"},
	}

	for _, tt := range tests {
		got := rr.Header().Get(tt.header)
		if got != tt.want {
			t.Errorf("Header %s = %q, want %q", tt.header, got, tt.want)
		}
	}

	// Check CSP is set
	csp := rr.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("Content-Security-Policy header should be set")
	}

	// Check Permissions-Policy is set
	pp := rr.Header().Get("Permissions-Policy")
	if pp == "" {
		t.Error("Permissions-Policy header should be set")
	}

	// Check cache control for API paths
	cc := rr.Header().Get("Cache-Control")
	if !strings.Contains(cc, "no-store") {
		t.Error("Cache-Control should contain no-store for API paths")
	}
}

func TestRequestSanitizer(t *testing.T) {
	handler := RequestSanitizer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("blocks suspicious paths", func(t *testing.T) {
		suspiciousPaths := []string{
			"/.env",
			"/.git/config",
			"/wp-admin/",
			"/wp-login.php",
			"/phpinfo.php",
		}

		for _, path := range suspiciousPaths {
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Errorf("Path %s should return 404, got %d", path, rr.Code)
			}
		}
	})

	t.Run("allows valid paths", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Valid path should return 200, got %d", rr.Code)
		}
	})

	t.Run("validates content type for POST", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/users", nil)
		req.Header.Set("Content-Type", "text/xml")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnsupportedMediaType {
			t.Errorf("Invalid content type should return 415, got %d", rr.Code)
		}
	})

	t.Run("allows JSON content type", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/users", nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("JSON content type should be allowed, got %d", rr.Code)
		}
	})
}

func TestEnhancedRequestID(t *testing.T) {
	handler := EnhancedRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check context values
		reqID := r.Context().Value(ContextKeyRequestID)
		if reqID == nil {
			t.Error("Request ID should be in context")
		}

		corrID := r.Context().Value(ContextKeyCorrelationID)
		if corrID == nil {
			t.Error("Correlation ID should be in context")
		}

		w.WriteHeader(http.StatusOK)
	}))

	t.Run("generates request ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("X-Request-ID") == "" {
			t.Error("X-Request-ID header should be set")
		}
	})

	t.Run("uses existing request ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", "existing-id")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("X-Request-ID") != "existing-id" {
			t.Error("Should use existing X-Request-ID")
		}
	})

	t.Run("uses existing correlation ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Correlation-ID", "trace-123")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("X-Correlation-ID") != "trace-123" {
			t.Error("Should use existing X-Correlation-ID")
		}
	})
}

func TestSecureCORS(t *testing.T) {
	allowedOrigins := []string{"https://hustlex.com", "https://app.hustlex.com"}
	handler := SecureCORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("allows configured origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://hustlex.com")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("Access-Control-Allow-Origin") != "https://hustlex.com" {
			t.Error("Should allow configured origin")
		}
	})

	t.Run("handles preflight", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "https://hustlex.com")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Preflight should return 204, got %d", rr.Code)
		}
	})

	t.Run("does not allow unauthorized origins", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://evil.com")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Error("Should not allow unauthorized origin")
		}
	})
}

func TestRecoverPanic(t *testing.T) {
	handler := RecoverPanic(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	// Should not panic
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Panic should return 500, got %d", rr.Code)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		xff        string
		xri        string
		want       string
	}{
		{
			name:       "from X-Forwarded-For",
			remoteAddr: "127.0.0.1:12345",
			xff:        "203.0.113.195, 70.41.3.18",
			want:       "203.0.113.195",
		},
		{
			name:       "from X-Real-IP",
			remoteAddr: "127.0.0.1:12345",
			xri:        "198.51.100.178",
			want:       "198.51.100.178",
		},
		{
			name:       "from RemoteAddr",
			remoteAddr: "192.168.1.1:12345",
			want:       "192.168.1.1",
		},
		{
			name:       "XFF takes precedence",
			remoteAddr: "127.0.0.1:12345",
			xff:        "203.0.113.195",
			xri:        "198.51.100.178",
			want:       "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xff != "" {
				req.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xri != "" {
				req.Header.Set("X-Real-IP", tt.xri)
			}

			got := getClientIP(req)
			if got != tt.want {
				t.Errorf("getClientIP() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestChain(t *testing.T) {
	var order []string

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	}

	handler := Chain(m1, m2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("Chain execution order length = %d, want %d", len(order), len(expected))
	}

	for i, v := range expected {
		if order[i] != v {
			t.Errorf("Chain execution order[%d] = %q, want %q", i, order[i], v)
		}
	}
}

func TestDefaultSecurityChain(t *testing.T) {
	handler := DefaultSecurityChain()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// Should have security headers
	if rr.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("Security headers should be applied")
	}

	// Should have request ID
	if rr.Header().Get("X-Request-ID") == "" {
		t.Error("Request ID should be generated")
	}
}

func TestStatusResponseWriter(t *testing.T) {
	t.Run("captures status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		wrapped.WriteHeader(http.StatusNotFound)

		if wrapped.statusCode != http.StatusNotFound {
			t.Errorf("statusCode = %d, want 404", wrapped.statusCode)
		}
	})

	t.Run("defaults to 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Write without calling WriteHeader
		wrapped.Write([]byte("test"))

		if wrapped.statusCode != http.StatusOK {
			t.Errorf("statusCode = %d, want 200", wrapped.statusCode)
		}
	})
}
