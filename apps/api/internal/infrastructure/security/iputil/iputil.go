package iputil

import (
	"net"
	"net/http"
	"strings"
)

// IPExtractor provides secure IP address extraction with trusted proxy validation
type IPExtractor struct {
	trustedProxies []*net.IPNet
}

// NewIPExtractor creates a new IPExtractor with trusted proxy CIDR ranges
func NewIPExtractor(trustedProxyCIDRs []string) (*IPExtractor, error) {
	var trustedProxies []*net.IPNet

	for _, cidr := range trustedProxyCIDRs {
		if cidr == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		trustedProxies = append(trustedProxies, ipNet)
	}

	return &IPExtractor{
		trustedProxies: trustedProxies,
	}, nil
}

// isTrustedProxy checks if an IP address is within trusted proxy ranges
func (e *IPExtractor) isTrustedProxy(ip net.IP) bool {
	for _, ipNet := range e.trustedProxies {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// GetClientIP extracts the real client IP from an HTTP request
// It validates X-Forwarded-For header only when the request comes from a trusted proxy
// This prevents IP spoofing attacks by untrusted clients
func (e *IPExtractor) GetClientIP(r *http.Request) string {
	// Parse the immediate client address (direct connection)
	remoteAddr := r.RemoteAddr
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		remoteAddr = host
	}

	remoteIP := net.ParseIP(remoteAddr)
	if remoteIP == nil {
		// If we can't parse, fall back to RemoteAddr
		return remoteAddr
	}

	// Only trust X-Forwarded-For if request comes from a trusted proxy
	if e.isTrustedProxy(remoteIP) {
		// Check X-Forwarded-For header (standard for proxied requests)
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
			// The leftmost IP is the original client
			ips := strings.Split(xff, ",")
			clientIP := strings.TrimSpace(ips[0])
			// Validate it's a valid IP
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}

		// Check X-Real-IP header (alternative used by some proxies)
		if xri := r.Header.Get("X-Real-IP"); xri != "" {
			clientIP := strings.TrimSpace(xri)
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// If not from trusted proxy or no valid forwarded headers, use direct IP
	return remoteAddr
}
