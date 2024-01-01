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

func (app *application) gRPCGetProductsHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newInventoryGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var getProductsRequest *pb.GetProductsRequest

		// check if userID exists in the Gin context
		if userID, found := c.Get("userID"); found {
			id, err := strconv.Atoi(userID.(string))
			if err != nil {
				log.Println("Failed to convert userID to int in request context:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}

			getProductsRequest = &pb.GetProductsRequest{
				UserID: int32(id),
			}
		} else {
			log.Println("UserID not found in jwt token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		resp, err := grpcClient.GetProducts(ctx, getProductsRequest)
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
			"status":   http.StatusOK,
			"products": responseData,
		})
	}
}

func (app *application) gRPCGetProductByIDHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newInventoryGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

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
		if userID, found := c.Get("userID"); found {
			id, err := strconv.Atoi(userID.(string))
			if err != nil {
				log.Println("Failed to convert userID to int in request context:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}

			getProductByIDRequest = &pb.GetProductByIDRequest{
				UserID:    int32(id),
				ProductID: int32(productID),
			}
		} else {
			log.Println("UserID not found in jwt token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		resp, err := grpcClient.GetProductByID(ctx, getProductByIDRequest)
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

func (app *application) gRPCCreateProductHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newInventoryGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		var req models.CreateProductRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		var createProductRequest *pb.CreateProductRequest
		// check if userID exists in the Gin context
		if userID, found := c.Get("userID"); found {
			id, err := strconv.Atoi(userID.(string))
			if err != nil {
				log.Println("Failed to convert userID to int in request context:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}

			createProductRequest = &pb.CreateProductRequest{
				UserID:       int32(id),
				BrandName:    req.BrandName,
				CategoryName: req.CategoryName,
				ProductName:  req.ProductName,
				Description:  req.Description,
				Size:         req.Size,
				Color:        req.Color,
				Quantity:     req.Quantity,
			}
		} else {
			log.Println("UserID not found in jwt token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		_, err = grpcClient.CreateProduct(ctx, createProductRequest)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Brand Name or Category Name does not exist. Please try again."})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			}
		} else {
			log.Println("Unable to retrieve error status", err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "Created", "message": fmt.Sprintf("Successfully created user %s", req.ProductName)})
	}
}

func (app *application) gRPCAuthenticationHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newAuthenticationGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		var req models.AuthRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// sending grpc request to grpc authenticate server
		resp, err := grpcClient.Authenticate(ctx, &pb.AuthRequest{
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

func (app *application) gRPCSignUpHandler(urlString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grpcClient, _, err := newAuthenticationGRPCClient(urlString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		var req models.SignUpRequest
		// decode JSON request into struct (HTTP/1.1)
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// sending grpc request to grpc authenticate server
		resp, err := grpcClient.SignUp(ctx, &pb.SignUpRequest{
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

func (app application) logoutHandler(c *gin.Context) {
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

func newAuthenticationGRPCClient(urlString string) (pb.AuthenticationServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing authentication-service: %v", err)
	}

	client := pb.NewAuthenticationServiceClient(conn)
	return client, conn, nil
}

func newUsersGRPCClient(urlString string) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing users-service: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return client, conn, nil
}

func newInventoryGRPCClient(urlString string) (pb.InventoryServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(urlString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing inventory-service: %v", err)
	}

	client := pb.NewInventoryServiceClient(conn)
	return client, conn, nil
}
