package main

import (
	"net/http"

	grpcclient "github.com/LeonLow97/internal/grpc"
	"github.com/gin-gonic/gin"
)

const (
	AUTHENTICATION_SERVICE_URL = "authentication-service:8001"
	INVENTORY_SERVICE_URL      = "inventory-service:8002"
	ORDER_SERVICE_URL          = "order-service:8003"
)

func (app *application) routes(grpcClient *grpcClientConn) *gin.Engine {
	router := gin.Default()

	router.Use(app.ipWhitelistMiddleware())
	router.Use(app.rateLimitMiddleware())

	// for pinging and testing the api gateway
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "api gateway healthy and running!"})
	})

	// gRPC Communication with Authentication service
	authenticationServiceClient := grpcclient.NewAuthenticationGRPCClient(grpcClient.authConn)

	authenticationHandlerGRPC := authenticationServiceClient.GRPCAuthenticateHandler()
	signUpHandlerGRPC := authenticationServiceClient.GRPCSignUpHandler()

	// setting up different paths to handle requests for each microservice
	router.POST("/authenticate", authenticationHandlerGRPC)
	router.POST("/signup", signUpHandlerGRPC)
	router.POST("/logout", authenticationServiceClient.LogoutHandler)

	// authentication microservice
	userServiceClient := grpcclient.NewUserGRPCClient(grpcClient.authConn)
	updateUserHandlerGRPC := userServiceClient.GRPCGetUsersHandler()
	getUsersHandlerGRPC := userServiceClient.GRPCGetUsersHandler()

	router.PATCH("/user", app.authenticationMiddleware(), updateUserHandlerGRPC)
	router.GET("/users", app.authenticationMiddleware(), getUsersHandlerGRPC)

	// inventory microservice
	inventoryServiceClient := grpcclient.NewInventoryGRPCClient(grpcClient.inventoryConn)

	getProductsHandlerGRPC := inventoryServiceClient.GRPCGetProductsHandler()
	getProductByIDHandlerGRPC := inventoryServiceClient.GRPCGetProductByIDHandler()
	createProductHandlerGRPC := inventoryServiceClient.GRPCCreateProductHandler()
	updateProductHandlerGRPC := inventoryServiceClient.GRPCUpdateProductHandler()
	deleteProductHandlerGRPC := inventoryServiceClient.GRPCDeleteProductHandler()

	inventoryServiceEndpoint := router.Group("/inventory")
	inventoryServiceEndpoint.Use(app.authenticationMiddleware())

	inventoryServiceEndpoint.GET("/products", getProductsHandlerGRPC)
	inventoryServiceEndpoint.GET("/product/:id", getProductByIDHandlerGRPC)
	inventoryServiceEndpoint.POST("/product", createProductHandlerGRPC)
	inventoryServiceEndpoint.PATCH("/product/:id", updateProductHandlerGRPC)
	inventoryServiceEndpoint.DELETE("/product/:id", deleteProductHandlerGRPC)

	// order microservice
	orderServiceClient := grpcclient.NewOrderGRPCClient(grpcClient.orderConn)

	getOrdersHandlerGRPC := orderServiceClient.GRPCGetOrdersHandler()
	getOrderHandlerGRPC := orderServiceClient.GRPCGetOrderHandler()
	createOrderHandlerGRPC := orderServiceClient.GRPCCreateOrderHandler()

	orderServiceEndpoint := router.Group("")
	orderServiceEndpoint.Use(app.authenticationMiddleware())

	orderServiceEndpoint.GET("/orders", getOrdersHandlerGRPC)
	orderServiceEndpoint.GET("/order/:id", getOrderHandlerGRPC)
	orderServiceEndpoint.POST("/order", createOrderHandlerGRPC)

	return router
}
