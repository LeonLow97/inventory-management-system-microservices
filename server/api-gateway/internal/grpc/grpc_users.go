package grpcclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LeonLow97/models"
	pb "github.com/LeonLow97/proto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceClient interface {
	GRPCGetUsersHandler() gin.HandlerFunc
	GRPCUpdateUserHandler() gin.HandlerFunc
}

type userGRPCClient struct {
	client pb.UserServiceClient
}

func NewUserGRPCClient(conn *grpc.ClientConn) UserServiceClient {
	return &userGRPCClient{
		client: pb.NewUserServiceClient(conn),
	}
}

func (u userGRPCClient) GRPCGetUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := u.client.GetUsers(ctx, &empty.Empty{})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		// convert protocol buffers message `resp` into JSON byte slice `b` while
		// specifying the use of protocol buffers field names for JSON serialization
		b, err := protojson.MarshalOptions{
			UseProtoNames: true,
		}.Marshal(resp)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		// Unmarshal the gRPC response into a map[string]interface{} to avoid double-escaping
		var responseData map[string]interface{}
		err = json.Unmarshal(b, &responseData)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"users":  responseData,
		})
	}
}

func (u userGRPCClient) GRPCUpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UpdateUserRequest
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

		// check if userID exists in the Gin context
		userID, err := retrieveUserIDFromToken(c)
		switch {
		case errors.Is(err, ErrMissingUserIDInJWTToken):
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err = u.client.UpdateUser(ctx, &pb.UpdateUserRequest{
			UserID:    int64(userID),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Password:  req.Password,
			Email:     req.Email,
		})
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					log.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
				default:
					log.Println(err)
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			}
		}

		c.JSON(http.StatusNoContent, gin.H{})
	}
}
