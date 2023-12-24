package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/LeonLow97/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) gRPCAuthenticationHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Error dialing logger-service: %v", err)
			return
		}
		defer conn.Close()

		client := pb.NewAuthenticationServiceClient(conn)

		var auth AuthRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&auth); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// sending grpc request to grpc authenticate server
		resp, err := client.Authenticate(ctx, &pb.AuthRequest{
			Username: auth.Username,
			Password: auth.Password,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				log.Printf("Authenticate grpc received error code %d", int32(errorCode))
				switch int32(errorCode) {
				case 16:
					c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid Credentials. Please try again."})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status")
				return
			}
		}

		cookie := &http.Cookie{
			Name:     "ims-token-oiweqj",
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
