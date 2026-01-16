package utils

import (
	"net/http"
	"strings"
)

// GetClientIP retrieves the client's real IP address, handling reverse proxies.
func GetClientIP(r *http.Request) string {
	// First, check for the X-Forwarded-For header (common when behind proxies/load balancers)
	// The format is: X-Forwarded-For: client-ip, proxy1, proxy2
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Split the IPs, and return the first one (the real client IP)
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // The first IP is usually the client's real IP
	}

	// Fallback to RemoteAddr if the X-Forwarded-For header is not present
	return strings.Split(r.RemoteAddr, ":")[0] // RemoteAddr contains IP:port, so we split on ":"
}
