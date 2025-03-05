package middleware

import (
	"log"
	"net"

	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/gin-gonic/gin"
)

var adminPaths = map[string]struct{}{
	"/users": {},
}

// ipWhitelistMiddleware check the client's IP against a list of allowed IP addresses (whitelist)
func (m *Middleware) IPWhitelistingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Perform IP whitelisting check for admin endpoints
		if _, isAdminPath := adminPaths[c.Request.URL.Path]; !isAdminPath {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		if !isAllowedIP(clientIP, m.cfg.AdminWhitelistedIPs) {
			log.Printf("admin IP %s is not whitelisted\n", clientIP)
			apierror.ErrForbidden.APIError(c, nil)
			return
		}

		c.Next()
	}
}

func isAllowedIP(ip string, adminWhitelistedIPs []string) bool {
	allowed := false
	clientIP := net.ParseIP(ip)

	for _, allowedIP := range adminWhitelistedIPs {
		if clientIP.Equal(net.ParseIP(allowedIP)) {
			allowed = true
			break
		}
	}

	return allowed
}
