package handlers

import (
	"errors"
	"net/http"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseHandler struct {
	validator *validator.Validate
	logger    ports.Logger
}

func newBaseHandler(logger ports.Logger) BaseHandler {
	return BaseHandler{
		validator: validatorInstance(logger),
		logger:    logger,
	}
}

func (h *BaseHandler) RespondError(c *gin.Context, err error) coreerrors.Kind {
	var appErr *coreerrors.Error
	if !errors.As(err, &appErr) {
		h.logger.WithContext(c.Request.Context()).WithError(err).Error("unclassified application error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return coreerrors.KindInternal
	}

	statusCode := http.StatusInternalServerError
	message := "Internal server error"
	switch appErr.Kind {
	case coreerrors.KindInvalidArgument:
		statusCode, message = http.StatusBadRequest, appErr.Message
	case coreerrors.KindUnauthenticated:
		statusCode, message = http.StatusUnauthorized, "Unauthorized"
	case coreerrors.KindPermissionDenied:
		statusCode, message = http.StatusForbidden, appErr.Message
	case coreerrors.KindNotFound:
		statusCode, message = http.StatusNotFound, appErr.Message
	case coreerrors.KindAlreadyExists, coreerrors.KindConflict:
		statusCode, message = http.StatusConflict, appErr.Message
	case coreerrors.KindFailedPrecondition:
		statusCode, message = http.StatusPreconditionFailed, appErr.Message
	case coreerrors.KindOutOfRange:
		statusCode, message = http.StatusBadRequest, appErr.Message
	case coreerrors.KindUnavailable:
		statusCode, message = http.StatusServiceUnavailable, "Service temporarily unavailable"
	case coreerrors.KindDeadlineExceeded:
		statusCode, message = http.StatusGatewayTimeout, "Request timed out"
	case coreerrors.KindResourceExhausted:
		statusCode, message = http.StatusTooManyRequests, "Rate limit exceeded"
	case coreerrors.KindInternal:
		if appErr.Cause != nil {
			h.logger.WithContext(c.Request.Context()).WithError(appErr.Cause).Error("application internal error")
		}
	}

	c.JSON(statusCode, gin.H{"error": message})
	return appErr.Kind
}
