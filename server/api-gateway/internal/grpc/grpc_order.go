package grpcclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LeonLow97/models"
	pb "github.com/LeonLow97/proto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type OrderServiceClient interface {
	GRPCGetOrdersHandler() gin.HandlerFunc
	GRPCGetOrderHandler() gin.HandlerFunc
	GRPCCreateOrderHandler() gin.HandlerFunc
}

type orderGRPCClient struct {
	client pb.OrderServiceClient
}

func NewOrderGRPCClient(conn *grpc.ClientConn) OrderServiceClient {
	return &orderGRPCClient{
		client: pb.NewOrderServiceClient(conn),
	}
}

func (o orderGRPCClient) GRPCGetOrdersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var getOrdersRequest *pb.GetOrdersRequest

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

		getOrdersRequest = &pb.GetOrdersRequest{
			UserID: int64(userID),
		}

		resp, err := o.client.GetOrders(ctx, getOrdersRequest)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				return
			}
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
			"orders": responseData,
		})
	}
}

func (o orderGRPCClient) GRPCGetOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// retrieve product id from path param and convert to int32 type for grpc
		orderIDString := c.Param("id")
		orderID, err := strconv.ParseInt(orderIDString, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var getOrderRequest *pb.GetOrderRequest

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

		getOrderRequest = &pb.GetOrderRequest{
			UserID:  int64(userID),
			OrderID: orderID,
		}

		resp, err := o.client.GetOrder(ctx, getOrderRequest)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				return
			}
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
			"order":  responseData,
		})
	}
}

func (o orderGRPCClient) GRPCCreateOrderHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var req models.CreateOrderRequest
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

		var createOrderRequest *pb.CreateOrderRequest

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

		createOrderRequest = &pb.CreateOrderRequest{
			UserID:       int64(userID),
			CustomerName: req.CustomerName,
			ProductName:  req.ProductName,
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			Color:        req.Color,
			Size:         req.Size,
			Quantity:     req.Quantity,
			Description:  req.Description,
			Revenue:      req.Revenue,
			Cost:         req.Cost,
			Profit:       req.Profit,
			HasReviewed:  req.HasReviewed,
		}

		_, err = o.client.CreateOrder(ctx, createOrderRequest)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully created order!"})
	}
}
