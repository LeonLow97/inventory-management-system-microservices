package middleware

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (m *Middleware) RateLimitingMiddleware() gin.HandlerFunc {
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
			apierror.ErrInternalServerError.APIError(c)
			return
		}

		ipRateLimiterMu.Lock()
		log.Println("Chaewon!!")
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
		if !clients[ip].limiter.Allow() {
			apierror.ErrTooManyRequests.APIError(c)
			return
		}

		c.Next()
	}
}
