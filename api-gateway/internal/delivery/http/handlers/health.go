package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	readinessChecker ReadinessChecker
}

type ReadinessChecker interface {
	Check(ctx context.Context) error
}

func NewHealthHandler(readinessChecker ReadinessChecker) *HealthHandler {
	return &HealthHandler{
		readinessChecker: readinessChecker,
	}
}

// Liveness just checks if the app process is running.
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "alive",
		"message": "api gateway is running",
	})
}

// Readiness checks external dependencies
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.readinessChecker.Check(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"message": "all systems operational",
	})
}
