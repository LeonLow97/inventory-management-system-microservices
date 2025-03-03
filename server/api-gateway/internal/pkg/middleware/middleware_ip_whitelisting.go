package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: Store this in database
var allowedIPs = []string{
	"127.0.0.1",
	"192.168.65.1",
}

// ipWhitelistMiddleware check the client's IP against a list of allowed IP addresses (whitelist)
func (m *Middleware) IPWhitelistingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow requests to the /healthcheck endpoint regardless of IP for k8s probing
		if c.Request.URL.Path == "/healthcheck" {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		if !isAllowedIP(clientIP) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

func isAllowedIP(ip string) bool {
	allowed := false
	clientIP := net.ParseIP(ip)

	for _, allowedIP := range allowedIPs {
		if clientIP.Equal(net.ParseIP(allowedIP)) {
			allowed = true
			break
		}
	}

	return allowed
}
