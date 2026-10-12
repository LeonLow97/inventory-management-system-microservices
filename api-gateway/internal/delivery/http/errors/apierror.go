package apierror

import (
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/core/ports"
	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest = &CustomError{
		Code:    "400000",
		Message: "Bad Request",
	}
)

// CustomError represents a structured error with an HTTP status code
type CustomError struct {
	Code    string
	Message string
}

// getHTTPStatus extracts the HTTP status code from the custom error code
// and returns the corresponding HTTP status message.
func (e *CustomError) getHTTPStatus() (int, string) {
	// Extract the first 3 digits of the error code (e.g., "401001" -> 401)
	if len(e.Code) < 3 {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	statusCode, err := strconv.Atoi(e.Code[:3])
	if err != nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}

	// Map HTTP status codes to generic messages
	switch statusCode {
	case http.StatusBadRequest:
		return http.StatusBadRequest, "Bad Request"
	case http.StatusUnauthorized:
		return http.StatusUnauthorized, "Unauthorized"
	case http.StatusForbidden:
		return http.StatusForbidden, "Forbidden"
	case http.StatusNotFound:
		return http.StatusNotFound, "Not Found"
	case http.StatusTooManyRequests:
		return http.StatusTooManyRequests, "Too Many Requests"
	case http.StatusConflict:
		return http.StatusConflict, "Conflict"
	case http.StatusInternalServerError:
		return http.StatusInternalServerError, "Internal Server Error"
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// APIError logs the original error and returns a JSON response.
func (e *CustomError) APIError(c *gin.Context, err error, logger ports.Logger) {
	fields := []ports.LogField{{Key: "message", Value: e.Message}}
	requestLogger := logger.WithContext(c.Request.Context())
	if err != nil {
		requestLogger = requestLogger.WithError(err)
	}
	requestLogger.Warn("HTTP request rejected", fields...)

	httpStatus, genericMessage := e.getHTTPStatus()
	c.AbortWithStatusJSON(httpStatus, gin.H{"message": genericMessage})
}
