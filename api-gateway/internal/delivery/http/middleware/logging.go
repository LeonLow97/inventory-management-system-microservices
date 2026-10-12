package middleware

import (
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/LeonLow97/internal/contextstore"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/gin-gonic/gin"
)

// RequestLogger records each completed HTTP request through the application logger.
func RequestLogger(logger ports.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		fields := []ports.LogField{
			{Key: "request_id", Value: requestIDFromContext(c)},
			{Key: "method", Value: c.Request.Method},
			{Key: "route", Value: c.FullPath()},
			{Key: "status", Value: c.Writer.Status()},
			{Key: "client_ip", Value: c.ClientIP()},
			{Key: "duration_ms", Value: time.Since(startedAt).Milliseconds()},
		}
		requestLogger := logger.WithContext(c.Request.Context())

		switch {
		case c.Writer.Status() >= http.StatusInternalServerError:
			requestLogger.Error("HTTP request completed", fields...)
		case c.Writer.Status() >= http.StatusBadRequest:
			requestLogger.Warn("HTTP request completed", fields...)
		default:
			requestLogger.Info("HTTP request completed", fields...)
		}
	}
}

// Recovery logs panics through the application logger and returns a generic 500 response.
func Recovery(logger ports.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		panicLogger := logger.WithContext(c.Request.Context())
		if panicErr, ok := recovered.(error); ok {
			panicLogger = panicLogger.WithError(panicErr)
		}
		panicLogger.Error("panic recovered while handling HTTP request",
			ports.LogField{Key: "request_id", Value: requestIDFromContext(c)},
			ports.LogField{Key: "panic", Value: recovered},
			ports.LogField{Key: "stack", Value: string(debug.Stack())},
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	})
}

func requestIDFromContext(c *gin.Context) string {
	requestID, _ := contextstore.RequestIDFromContext(c.Request.Context())
	return requestID
}
