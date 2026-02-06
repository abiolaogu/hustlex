package iputil

import (
	"net/http"
	"testing"
)

func TestNewIPExtractor(t *testing.T) {
	tests := []struct {
		name        string
		cidrs       []string
		expectError bool
	}{
		{
			name:        "Valid CIDR ranges",
			cidrs:       []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
			expectError: false,
		},
		{
			name:        "Empty CIDR list",
			cidrs:       []string{},
			expectError: false,
		},
		{
			name:        "Invalid CIDR format",
			cidrs:       []string{"10.0.0.0/8", "invalid-cidr"},
			expectError: true,
		},
		{
			name:        "CIDR with empty strings",
			cidrs:       []string{"10.0.0.0/8", "", "192.168.0.0/16"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := NewIPExtractor(tt.cidrs)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectError && extractor == nil {
				t.Error("Expected valid extractor but got nil")
			}
		})
	}
}

func TestGetClientIP_FromTrustedProxy(t *testing.T) {
	// Setup extractor with common private network ranges
	extractor, err := NewIPExtractor([]string{
		"10.0.0.0/8",      // Private network
		"172.16.0.0/12",   // Private network
		"192.168.0.0/16",  // Private network
	})
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		xffHeader  string
		xriHeader  string
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For from trusted proxy (10.0.0.0/8)",
			remoteAddr: "10.0.0.5:12345",
			xffHeader:  "203.0.113.1",
			expectedIP: "203.0.113.1",
		},
		{
			name:       "X-Forwarded-For from trusted proxy (172.16.0.0/12)",
			remoteAddr: "172.16.0.1:54321",
			xffHeader:  "198.51.100.42",
			expectedIP: "198.51.100.42",
		},
		{
			name:       "X-Forwarded-For from trusted proxy (192.168.0.0/16)",
			remoteAddr: "192.168.1.1:8080",
			xffHeader:  "192.0.2.100",
			expectedIP: "192.0.2.100",
		},
		{
			name:       "Multiple IPs in X-Forwarded-For (takes leftmost)",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "203.0.113.5, 192.0.2.10, 198.51.100.20",
			expectedIP: "203.0.113.5",
		},
		{
			name:       "X-Real-IP from trusted proxy when no X-Forwarded-For",
			remoteAddr: "10.0.0.1:1234",
			xriHeader:  "203.0.113.10",
			expectedIP: "203.0.113.10",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "203.0.113.20",
			xriHeader:  "203.0.113.30",
			expectedIP: "203.0.113.20",
		},
		{
			name:       "IPv6 in X-Forwarded-For from trusted proxy",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "2001:db8::1",
			expectedIP: "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}
			if tt.xffHeader != "" {
				req.Header.Set("X-Forwarded-For", tt.xffHeader)
			}
			if tt.xriHeader != "" {
				req.Header.Set("X-Real-IP", tt.xriHeader)
			}

			ip := extractor.GetClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

func TestGetClientIP_FromUntrustedSource(t *testing.T) {
	// Setup extractor with specific trusted proxy range
	extractor, err := NewIPExtractor([]string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		xffHeader  string
		xriHeader  string
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For from untrusted source (ignored)",
			remoteAddr: "203.0.113.50:12345",
			xffHeader:  "198.51.100.1",
			expectedIP: "203.0.113.50", // Should use RemoteAddr, not XFF
		},
		{
			name:       "X-Real-IP from untrusted source (ignored)",
			remoteAddr: "203.0.113.60:54321",
			xriHeader:  "198.51.100.2",
			expectedIP: "203.0.113.60", // Should use RemoteAddr, not XRI
		},
		{
			name:       "Both headers from untrusted source (ignored)",
			remoteAddr: "203.0.113.70:8080",
			xffHeader:  "198.51.100.3",
			xriHeader:  "198.51.100.4",
			expectedIP: "203.0.113.70", // Should use RemoteAddr
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}
			if tt.xffHeader != "" {
				req.Header.Set("X-Forwarded-For", tt.xffHeader)
			}
			if tt.xriHeader != "" {
				req.Header.Set("X-Real-IP", tt.xriHeader)
			}

			ip := extractor.GetClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

func TestGetClientIP_InvalidHeaders(t *testing.T) {
	extractor, err := NewIPExtractor([]string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		xffHeader  string
		expectedIP string
	}{
		{
			name:       "Invalid IP in X-Forwarded-For",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "not-an-ip-address",
			expectedIP: "10.0.0.1", // Should fall back to RemoteAddr
		},
		{
			name:       "Empty X-Forwarded-For",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "",
			expectedIP: "10.0.0.1",
		},
		{
			name:       "Malformed X-Forwarded-For",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  ";;;invalid;;;",
			expectedIP: "10.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}
			if tt.xffHeader != "" {
				req.Header.Set("X-Forwarded-For", tt.xffHeader)
			}

			ip := extractor.GetClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

func TestGetClientIP_NoProxy(t *testing.T) {
	extractor, err := NewIPExtractor([]string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		expectedIP string
	}{
		{
			name:       "Direct connection with port",
			remoteAddr: "203.0.113.1:12345",
			expectedIP: "203.0.113.1",
		},
		{
			name:       "Direct connection without port",
			remoteAddr: "203.0.113.2",
			expectedIP: "203.0.113.2",
		},
		{
			name:       "IPv6 direct connection",
			remoteAddr: "[2001:db8::1]:8080",
			expectedIP: "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}

			ip := extractor.GetClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

func TestGetClientIP_EdgeCases(t *testing.T) {
	extractor, err := NewIPExtractor([]string{"10.0.0.0/8"})
	if err != nil {
		t.Fatalf("Failed to create extractor: %v", err)
	}

	tests := []struct {
		name       string
		remoteAddr string
		xffHeader  string
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For with spaces",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "  203.0.113.1  , 192.0.2.1  ",
			expectedIP: "203.0.113.1",
		},
		{
			name:       "X-Forwarded-For with mixed valid/invalid IPs",
			remoteAddr: "10.0.0.1:1234",
			xffHeader:  "203.0.113.1, invalid, 192.0.2.1",
			expectedIP: "203.0.113.1", // Takes first valid IP
		},
		{
			name:       "Localhost RemoteAddr",
			remoteAddr: "127.0.0.1:1234",
			xffHeader:  "203.0.113.1",
			expectedIP: "127.0.0.1", // localhost not in trusted range
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				RemoteAddr: tt.remoteAddr,
				Header:     make(http.Header),
			}
			if tt.xffHeader != "" {
				req.Header.Set("X-Forwarded-For", tt.xffHeader)
			}

			ip := extractor.GetClientIP(req)
			if ip != tt.expectedIP {
				t.Errorf("Expected IP %s, got %s", tt.expectedIP, ip)
			}
		})
	}
}

// Benchmark tests
func BenchmarkGetClientIP_DirectConnection(b *testing.B) {
	extractor, _ := NewIPExtractor([]string{"10.0.0.0/8"})
	req := &http.Request{
		RemoteAddr: "203.0.113.1:12345",
		Header:     make(http.Header),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractor.GetClientIP(req)
	}
}

func BenchmarkGetClientIP_TrustedProxy(b *testing.B) {
	extractor, _ := NewIPExtractor([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"})
	req := &http.Request{
		RemoteAddr: "10.0.0.1:12345",
		Header:     make(http.Header),
	}
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 192.0.2.1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractor.GetClientIP(req)
	}
}
