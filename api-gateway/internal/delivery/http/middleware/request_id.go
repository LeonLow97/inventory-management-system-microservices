package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/LeonLow97/internal/contextstore"
	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

// RequestID creates a request ID, stores it in the request context, and returns it to the caller.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, err := newRequestID()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx := contextstore.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestIDHeader, requestID)
		c.Next()
	}
}

func newRequestID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}

	// Mark the random ID as a UUID v4 so it has a familiar, standard representation.
	id[6] = (id[6] & 0x0f) | 0x40
	id[8] = (id[8] & 0x3f) | 0x80
	encoded := hex.EncodeToString(id[:])
	return encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:], nil
}
