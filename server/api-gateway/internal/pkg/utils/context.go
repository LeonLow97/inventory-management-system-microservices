package utils

import (
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrMissingUserIDInGinContext = errors.New("missing userID in gin context")

// GetUserIDFromContext extracts and converts userID from the Gin context.
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, found := c.Get("userID")
	if !found {
		log.Println("UserID not found in JWT token claims")
		return 0, ErrMissingUserIDInGinContext
	}

	userIDStr, ok := userID.(string)
	if !ok {
		log.Println("UserID in context is not a string")
		return 0, ErrMissingUserIDInGinContext
	}

	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("Failed to convert userID to int in request context:", err)
		return 0, err
	}

	return userIDInt, nil
}
