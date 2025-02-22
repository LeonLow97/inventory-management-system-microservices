package apierror

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// Generic errors
	ErrUnauthorized = &CustomError{
		Code:    "401000",
		Message: "Unauthorized",
	}
	ErrTooManyRequests = &CustomError{
		Code:    "429000",
		Message: "Too Many Requests",
	}
	ErrInternalServerError = &CustomError{
		Code:    "500000",
		Message: "Internal Server Error",
	}
)

// CustomError represents a structured error with an HTTP status code
type CustomError struct {
	Code    string
	Message string
}

// Error implements the error interface
func (e *CustomError) Error() string {
	return fmt.Sprintf("Error Code: %s | Message: %s", e.Code, e.Message)
}

// getHTTPStatus extracts the HTTP status code from the custom error code
// and returns the corresponding HTTP status message.
func (e *CustomError) getHTTPStatus() (int, string) {
	// Extract the first 3 digits of the error code (e.g., "401001" -> 401)
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
	case http.StatusInternalServerError:
		return http.StatusInternalServerError, "Internal Server Error"
	default:
		log.Printf("An expected error code: %d\n", statusCode)
		return http.StatusInternalServerError, "Internal Server Error"
	}
}

// APIError logs the original error and returns a JSON response
func (e *CustomError) APIError(c *gin.Context) {
	log.Printf("HTTP Error: %s", e.Message)

	httpStatus, genericMessage := e.getHTTPStatus()
	c.AbortWithStatusJSON(httpStatus, gin.H{"message": genericMessage})
}
