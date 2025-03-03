package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/config"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services/auth"
	"github.com/LeonLow97/internal/pkg/apierror"
	"github.com/LeonLow97/internal/pkg/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	handler.Handler
	cfg         config.Config
	AuthService auth.Auth
}

func NewAuthHandler(cfg config.Config, authService auth.Auth) *AuthHandler {
	return &AuthHandler{
		Handler:     handler.NewHandler(),
		cfg:         cfg,
		AuthService: authService,
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}

		// Validate http request json property values
		if err := h.ValidateStruct(req); err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}

		user := domain.User{
			Email:    req.Email,
			Password: req.Password,
		}

		resp, err := h.AuthService.Login(c, user)
		if err != nil {
			status, ok := status.FromError(err)
			if !ok {
				apierror.ErrInternalServerError.APIError(c, err)
				return
			}

			var customErr *apierror.CustomError
			switch status.Code() {
			case codes.InvalidArgument:
				customErr = apierror.ErrBadRequest
			case codes.Unauthenticated:
				customErr = apierror.ErrUnauthorized
			default:
				customErr = apierror.ErrInternalServerError
			}
			customErr.APIError(c, err)
			return
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     h.cfg.AuthJWTToken.Name,
			Value:    resp.Token,
			MaxAge:   h.cfg.AuthJWTToken.MaxAge,
			Path:     h.cfg.AuthJWTToken.Path,
			Domain:   h.cfg.AuthJWTToken.Domain,
			Secure:   h.cfg.AuthJWTToken.Secure,
			HttpOnly: h.cfg.AuthJWTToken.HTTPOnly,
		})

		c.JSON(http.StatusOK, dto.LoginResponse{
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
			Email:     resp.Email,
			Active:    resp.Active,
			Admin:     resp.Admin,
			Token:     resp.Token,
		})
	}
}

func (h *AuthHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SignUpRequest
		if err := c.BindJSON(&req); err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}

		if err := h.ValidateStruct(req); err != nil {
			apierror.ErrBadRequest.APIError(c, err)
			return
		}

		user := domain.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Password:  req.Password,
			Email:     req.Email,
		}

		if err := h.AuthService.SignUp(c, user); err != nil {
			status, ok := status.FromError(err)
			if !ok {
				apierror.ErrInternalServerError.APIError(c, err)
				return
			}

			var customErr *apierror.CustomError
			switch status.Code() {
			case codes.InvalidArgument:
				customErr = apierror.ErrBadRequest
			case codes.AlreadyExists:
				customErr = apierror.ErrEmailAlreadyExists
			default:
				customErr = apierror.ErrInternalServerError
			}
			customErr.APIError(c, err)
			return
		}

		h.ResponseNoContent(c, http.StatusCreated)
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(h.cfg.AuthJWTToken.Name)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				log.Printf("No %s cookie found\n", h.cfg.AuthJWTToken.Name)
			default:
				log.Printf("failed to logout with error: %v\n", err)
			}
			h.ResponseNoContent(c, http.StatusOK)
			return
		}

		cookie.MaxAge = -1               // Invalidate the existing cookie by setting MaxAge to -1
		http.SetCookie(c.Writer, cookie) // Update the cookie in the response header to invalidate it
		h.ResponseNoContent(c, http.StatusOK)
	}
}
