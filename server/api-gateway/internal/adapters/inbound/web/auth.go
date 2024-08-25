package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	AuthService auth.Auth
}

func NewAuthHandler(authService auth.Auth) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("failed to bind JSON in Login handler: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		// validate http request json property values
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			log.Printf("validation failed for request in Login handler: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		user := domain.User{
			Username: req.Username,
			Password: req.Password,
		}

		domainResp, err := h.AuthService.Login(c, user)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				log.Printf("gRPC authentication error: code=%d, message=%s", errorCode, status.Message())

				switch errorCode {
				case codes.InvalidArgument: // Code 3
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  http.StatusBadRequest,
						"message": fmt.Sprintf("Bad Request: %s", status.Message()),
					})
				case codes.Unauthenticated: // Code 16
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  http.StatusUnauthorized,
						"message": "Invalid credentials. Please try again.",
					})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{
						"status":  http.StatusInternalServerError,
						"message": "Internal Server Error",
					})
				}
				return
			} else {
				log.Printf("Error retrieving gRPC status: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Internal Server Error",
				})
				return
			}
		}

		resp := &dto.LoginResponse{
			FirstName: domainResp.FirstName,
			LastName:  domainResp.LastName,
			Username:  domainResp.Username,
			Email:     domainResp.Email,
			Active:    int32(domainResp.Active),
			Admin:     int32(domainResp.Admin),
			Token:     domainResp.Token,
		}

		cookie := &http.Cookie{
			Name:     "ims-token",
			Value:    resp.Token,
			MaxAge:   3600,
			Path:     "/",
			Domain:   "localhost",
			Secure:   false,
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, gin.H{
			"first_name": resp.FirstName,
			"last_name":  resp.LastName,
			"username":   resp.Username,
			"email":      resp.Email,
			"active":     resp.Active,
			"admin":      resp.Admin,
		})
	}
}

func (h *AuthHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SignUpRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("failed to bind JSON in SignUp handler: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			log.Printf("validation failed for request in SignUp handler: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		user := domain.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Username:  req.Username,
			Password:  req.Password,
			Email:     req.Email,
		}

		if err := h.AuthService.SignUp(c, user); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				log.Printf("SignUp grpc received error code %d with err %v", int32(errorCode), err)
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 6:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Username %s has been taken.", req.Username)})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status")
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"status": "Created", "message": fmt.Sprintf("Successfully created user %s", req.Username)})
	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("ims-token")
		if err == http.ErrNoCookie {
			log.Println("No 'ims-token' cookie found.")
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Logged out successfully!"})
		} else if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
		}

		cookie.MaxAge = -1 // Invalidate the existing cookie by setting MaxAge to -1

		http.SetCookie(c.Writer, cookie) // Update the cookie in the response header to invalidate it

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Logged out successfully!"})
	}
}
