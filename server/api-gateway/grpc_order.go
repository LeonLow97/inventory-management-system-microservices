package main

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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func newOrderGRPCClient(urlString string) (pb.OrderServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing order-service: %v", err)
	}

	client := pb.NewOrderServiceClient(conn)
	return client, conn, nil
}

func (app *application) gRPCGetOrdersHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newOrderGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var getOrdersRequest *pb.GetOrdersRequest

		// check if userID exists in the Gin context
		userID, err := app.retrieveUserIDFromToken(c)
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

		resp, err := grpcClient.GetOrders(ctx, getOrdersRequest)
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

func (app *application) gRPCGetOrderHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newOrderGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

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
		userID, err := app.retrieveUserIDFromToken(c)
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

		resp, err := grpcClient.GetOrder(ctx, getOrderRequest)
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

func (app *application) gRPCCreateOrderHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newOrderGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

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
		if err = validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		var createOrderRequest *pb.CreateOrderRequest

		// check if userID exists in the Gin context
		userID, err := app.retrieveUserIDFromToken(c)
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

		_, err = grpcClient.CreateOrder(ctx, createOrderRequest)
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
