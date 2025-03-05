package middleware

import (
	"github.com/gin-gonic/gin"
)

func (m *Middleware) RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// // Retrieve user ID from request context
		// userID, err := contextstore.UserIDFromContext(c)
		// if err != nil {
		// 	switch {
		// 	case errors.Is(err, contextstore.ErrUserIDNotInContext):
		// 		apierror.ErrUnauthorized.APIError(c, err)
		// 	default:
		// 		apierror.ErrInternalServerError.APIError(c, err)
		// 	}
		// 	return
		// }

		// // Dynamically retrieve bucket name based on request http method
		// // TODO: Refactor this section to also include authentication bucket
		// bucketName := m.cfg.RateLimiting.DistributedLocks.Global
		// if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut ||
		// 	c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodDelete {
		// 	bucketName = m.cfg.RateLimiting.DistributedLocks.Write
		// } else if c.Request.Method == http.MethodGet {
		// 	bucketName = m.cfg.RateLimiting.DistributedLocks.Read
		// }

		// // Try to acquire a lock before modifying the token count
		// if !m.appCache.AcquireLock(c, userID, bucketName) {
		// 	log.Println("could not acquire distributed lock, try again later") // TODO: remove log in future to avoid flooding logs
		// 	apierror.ErrTooManyRequests.APIError(c, nil)
		// 	return
		// }
		// defer m.appCache.ReleaseLock(c, userID, bucketName)

		// // Retrieve current token count
		// key := m.appCache.UserTokenBucketKey(userID, bucketName)
		// currentTokens, err := m.appCache.RedisClient.Get(c, key).Int()
		// if err != nil && err != redis.Nil {
		// 	log.Println("failed to fetch token counts with error:", err)
		// 	apierror.ErrInternalServerError.APIError(c, err)
		// 	return
		// }

		// // Check if there are available tokens
		// if currentTokens == 0 {
		// 	log.Println("No tokens available")
		// 	apierror.ErrTooManyRequests.APIError(c, nil)
		// 	return
		// }
		// if currentTokens > 0 {
		// 	m.appCache.RedisClient.Decr(c, key)
		// }

		c.Next()
	}
}
