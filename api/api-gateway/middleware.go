package main

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

func IPWhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !isAllowedIP(clientIP) {
			// TODO: Might need to fix this error handling way
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
