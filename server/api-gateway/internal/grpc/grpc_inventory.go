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
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/LeonLow97/proto"
)

type InventoryServiceClient interface {
	GRPCGetProductsHandler() gin.HandlerFunc
	GRPCGetProductByIDHandler() gin.HandlerFunc
	GRPCCreateProductHandler() gin.HandlerFunc
	GRPCUpdateProductHandler() gin.HandlerFunc
	GRPCDeleteProductHandler() gin.HandlerFunc
}

type inventoryGRPCClient struct {
	client pb.InventoryServiceClient
}

func NewInventoryGRPCClient(conn *grpc.ClientConn) InventoryServiceClient {
	return &inventoryGRPCClient{
		client: pb.NewInventoryServiceClient(conn),
	}
}

func (i inventoryGRPCClient) GRPCGetProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var getProductsRequest *pb.GetProductsRequest

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

		getProductsRequest = &pb.GetProductsRequest{
			UserID: int32(userID),
		}

		resp, err := i.client.GetProducts(ctx, getProductsRequest)
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
			"status":   http.StatusOK,
			"products": responseData,
		})
	}
}

func (i inventoryGRPCClient) GRPCGetProductByIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// retrieve product id from path param and convert to int32 type for grpc
		productIDString := c.Param("id")
		productID, err := strconv.ParseInt(productIDString, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		var getProductByIDRequest *pb.GetProductByIDRequest

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
		getProductByIDRequest = &pb.GetProductByIDRequest{
			UserID:    int32(userID),
			ProductID: int32(productID),
		}

		resp, err := i.client.GetProductByID(ctx, getProductByIDRequest)
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

		b, err := protojson.MarshalOptions{
			UseProtoNames: true,
		}.Marshal(resp)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		var responseData map[string]interface{}
		err = json.Unmarshal(b, &responseData)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"products": responseData,
		})
	}
}

func (i inventoryGRPCClient) GRPCCreateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var req models.CreateProductRequest
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

		var createProductRequest *pb.CreateProductRequest

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

		createProductRequest = &pb.CreateProductRequest{
			UserID:       int32(userID),
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			ProductName:  req.ProductName,
			Description:  req.Description,
			Size:         req.Size,
			Color:        req.Color,
			Quantity:     req.Quantity,
		}

		_, err = i.client.CreateProduct(ctx, createProductRequest)
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

		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": fmt.Sprintf("Successfully created %s", req.ProductName)})
	}
}

func (i inventoryGRPCClient) GRPCUpdateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// retrieve product id from path param and convert to int32 type for grpc
		productIDString := c.Param("id")
		productID, err := strconv.ParseInt(productIDString, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var req models.UpdateProductRequest
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

		var updateProductRequest *pb.UpdateProductRequest

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

		updateProductRequest = &pb.UpdateProductRequest{
			UserID:       int32(userID),
			ProductID:    int32(productID),
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			ProductName:  req.ProductName,
			Description:  req.Description,
			Size:         req.Size,
			Color:        req.Color,
			Quantity:     req.Quantity,
		}

		_, err = i.client.UpdateProduct(ctx, updateProductRequest)
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

		c.JSON(http.StatusNoContent, gin.H{"status": http.StatusNoContent, "message": "Updated Product!"})
	}
}

func (i inventoryGRPCClient) GRPCDeleteProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// retrieve product id from path param and convert to int32 type for grpc
		productIDString := c.Param("id")
		productID, err := strconv.ParseInt(productIDString, 10, 32)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var deleteProductRequest *pb.DeleteProductRequest

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

		deleteProductRequest = &pb.DeleteProductRequest{
			UserID:    int32(userID),
			ProductID: int32(productID),
		}

		_, err = i.client.DeleteProduct(ctx, deleteProductRequest)
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

		c.JSON(http.StatusNoContent, gin.H{"status": http.StatusNoContent, "message": "Deleted Product!"})
	}
}
