package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LeonLow97/models"
	pb "github.com/LeonLow97/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func newUsersGRPCClient(urlString string) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing users-service: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return client, conn, nil
}

func (app *application) grpcGetUsersHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newUsersGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		resp, err := grpcClient.GetUsers(ctx, &empty.Empty{})
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

func (app *application) grpcUpdateUserHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newUsersGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
		}

		var req models.UpdateUserRequest
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		var userIDContext int
		// check if userID exists in the Gin context
		if userID, found := c.Get("userID"); found {
			userIDContext, err = strconv.Atoi(userID.(string))
			if err != nil {
				log.Println("Failed to convert userID to int in request context:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		} else {
			log.Println("UserID not found in jwt token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, err = grpcClient.UpdateUser(ctx, &pb.UpdateUserRequest{
			UserID:    int64(userIDContext),
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
