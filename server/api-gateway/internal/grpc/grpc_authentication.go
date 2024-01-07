package grpcclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LeonLow97/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/LeonLow97/proto"
)

type AuthenticationServiceClient interface {
	GRPCAuthenticateHandler() gin.HandlerFunc
	GRPCSignUpHandler() gin.HandlerFunc
	LogoutHandler(c *gin.Context)
}

type authenticationGRPCClient struct {
	client pb.AuthenticationServiceClient
}

func NewAuthenticationGRPCClient(conn *grpc.ClientConn) AuthenticationServiceClient {
	return &authenticationGRPCClient{
		client: pb.NewAuthenticationServiceClient(conn),
	}
}

func (a authenticationGRPCClient) GRPCAuthenticateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AuthRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		// validate http request json property values
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// sending grpc request to grpc authenticate server
		resp, err := a.client.Authenticate(ctx, &pb.AuthRequest{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				log.Printf("Authenticate grpc received error code %d with err %v", int32(errorCode), err)
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 16:
					c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Credentials. Please try again."})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				return
			}
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

func (a authenticationGRPCClient) GRPCSignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignUpRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		// validate http request json property values
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// sending grpc request to grpc authenticate server
		resp, err := a.client.SignUp(ctx, &pb.SignUpRequest{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Username:  req.Username,
			Password:  req.Password,
			Email:     req.Email,
		})
		if err != nil {
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
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"status": "Created", "message": fmt.Sprintf("Successfully created user %s", resp.Username)})
	}
}

func (a authenticationGRPCClient) LogoutHandler(c *gin.Context) {
	cookie, err := c.Request.Cookie("ims-token")
	if err == http.ErrNoCookie {
		log.Println("No 'ims-token' cookie found.")
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Logged out successfully!"})
		return
	} else if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
		return
	}

	cookie.MaxAge = -1 // Invalidate the existing cookie by setting MaxAge to -1

	http.SetCookie(c.Writer, cookie) // Update the cookie in the response header to invalidate it

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Logged out successfully!"})
}
