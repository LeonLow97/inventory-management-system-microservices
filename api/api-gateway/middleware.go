package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// TODO: Store this in database
var allowedIPs = []string{
	"127.0.0.1",
	"192.168.65.1",
}

func (app *application) ipWhitelistMiddleware() gin.HandlerFunc {
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

// rateLimitMiddleware is the middleware for rate limiting client requests based on IP Address
func (app *application) rateLimitMiddleware() gin.HandlerFunc {
	// create a client type to store the last seen timing of the last request from client
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// create a rate limiter map for client IP Addresses
	var (
		// using another mutex so it does not block `clients` map during cleanup
		// in case there are many ip addresses in the map and the for loop takes a long time to execute
		cleanUpMu       sync.Mutex 
		ipRateLimiterMu sync.Mutex
		// key is client IP Address
		clients = make(map[string]*client)
	)

	// launching a goroutine in the background to remove entries from clients map
	// where lastSeen of last request made is more than 3 minutes, run this goroutine every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock before starting to clean `clients` map
			cleanUpMu.Lock()
			for ip, client := range clients {
				// if last seen is more than 3 minutes, remove this ip address from the `clients` map
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			cleanUpMu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		// get the IP Address of the client in the http request
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		ipRateLimiterMu.Lock()
		// assign a new rate limiter to the client ip address if it doesn't exist in `clients` map
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				// 1st argument (rate limit): 2 requests per second
				// 2nd argument (burst limit): allowing a burst up to 4 requests beyond the rate limit before further requests are limited
				limiter: rate.NewLimiter(2, 4),
			}
		}
		ipRateLimiterMu.Unlock()

		// update the last seen time of the client to now
		clients[ip].lastSeen = time.Now()

		// check if the request is allowed
		fmt.Println(clients[ip].limiter.Allow())
		fmt.Println(ip)
		if !clients[ip].limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
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
