package handlers

import (
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/delivery/http/dto"
	"github.com/LeonLow97/internal/delivery/http/errors"
	"github.com/gin-gonic/gin"
)

// JWTConfig contains the cookie settings used for authentication tokens.
type JWTConfig struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
}

type UserHandler struct {
	BaseHandler
	jwtConfig   JWTConfig
	UserService usecases.User
}

func NewUserHandler(jwtConfig JWTConfig, userService usecases.User, logger ports.Logger) *UserHandler {
	return &UserHandler{
		BaseHandler: newBaseHandler(logger),
		jwtConfig:   jwtConfig,
		UserService: userService,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	user, err := h.UserService.GetCurrentUser(c.Request.Context())
	if err != nil {
		_ = h.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToLoginResponse(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	// Deserialize request body and validate request fields
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}
	if err := h.ValidateStruct(req); err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}

	ctx := c.Request.Context()
	clientIP := c.ClientIP()

	resp, token, err := h.UserService.Login(ctx, dto.FromLoginRequest(&req), clientIP)
	if err != nil {
		_ = h.RespondError(c, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.jwtConfig.Name,
		Value:    token,
		MaxAge:   h.jwtConfig.MaxAge,
		Path:     h.jwtConfig.Path,
		Domain:   h.jwtConfig.Domain,
		Secure:   h.jwtConfig.Secure,
		HttpOnly: h.jwtConfig.HTTPOnly,
		SameSite: http.SameSiteLaxMode,
	})
	c.JSON(http.StatusOK, dto.ToLoginResponse(resp))
}

func (h *UserHandler) SignUp(c *gin.Context) {
	// Deserialize request body and validate request fields
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}
	if err := h.ValidateStruct(req); err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}

	if err := h.UserService.SignUp(c.Request.Context(), dto.FromSignUpRequest(&req)); err != nil {
		_ = h.RespondError(c, err)
		return
	}
	h.ResponseNoContent(c, http.StatusCreated)
}

func (h *UserHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     h.jwtConfig.Name,
		Value:    "",
		MaxAge:   -1,
		Path:     h.jwtConfig.Path,
		Domain:   h.jwtConfig.Domain,
		Secure:   h.jwtConfig.Secure,
		HttpOnly: h.jwtConfig.HTTPOnly,
		SameSite: http.SameSiteLaxMode,
	})
	h.ResponseNoContent(c, http.StatusOK)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	limitStr := h.GetQueryParam(c, "limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}
	cursor := h.GetQueryParam(c, "cursor", "")

	domainResp, nextCursor, err := h.UserService.GetUsers(c.Request.Context(), int64(limit), cursor)
	if err != nil {
		_ = h.RespondError(c, err)
		return
	}
	h.ResponseJSON(c, http.StatusOK, dto.ToGetUsersResponse(domainResp, nextCursor))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.ErrBadRequest.APIError(c, nil, h.logger)
		return
	}
	if err := h.ValidateStruct(req); err != nil {
		apierror.ErrBadRequest.APIError(c, err, h.logger)
		return
	}

	if err := h.UserService.UpdateUser(c.Request.Context(), dto.FromUpdateUserRequest(&req)); err != nil {
		_ = h.RespondError(c, err)
		return
	}
	h.ResponseNoContent(c, http.StatusNoContent)
}
