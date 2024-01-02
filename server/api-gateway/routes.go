package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const AUTHENTICATION_SERVICE_URL = "authentication-service:8001"
const INVENTORY_SERVICE_URL = "inventory-service:8002"

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.Use(app.ipWhitelistMiddleware())
	router.Use(app.rateLimitMiddleware())

	// gRPC Communication with Authentication service
	authenticationHandlerGRPC := app.gRPCAuthenticationHandler(AUTHENTICATION_SERVICE_URL)
	signUpHandlerGRPC := app.gRPCSignUpHandler(AUTHENTICATION_SERVICE_URL)

	updateUserHandlerGRPC := app.grpcUpdateUserHandler(AUTHENTICATION_SERVICE_URL)
	getUsersHandlerGRPC := app.grpcGetUsersHandler(AUTHENTICATION_SERVICE_URL)

	getProductsHandlerGRPC := app.gRPCGetProductsHandler(INVENTORY_SERVICE_URL)
	getProductByIDHandlerGRPC := app.gRPCGetProductByIDHandler(INVENTORY_SERVICE_URL)
	createProductHandlerGRPC := app.gRPCCreateProductHandler(INVENTORY_SERVICE_URL)
	updateProductHandlerGRPC := app.gRPCUpdateProductHandler(INVENTORY_SERVICE_URL)

	// for pinging and testing the api gateway
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "api gateway healthy and running!"})
	})

	// setting up different paths to handle requests for each microservice
	router.POST("/authenticate", authenticationHandlerGRPC)
	router.POST("/signup", signUpHandlerGRPC)
	router.POST("/logout", app.logoutHandler)

	// authentication microservice (protected resource)
	router.PATCH("/user", app.authenticationMiddleware(), updateUserHandlerGRPC)
	router.GET("/users", app.authenticationMiddleware(), getUsersHandlerGRPC)

	// inventory microservice (protected resource)
	inventoryServiceEndpoint := router.Group("/inventory")
	inventoryServiceEndpoint.Use(app.authenticationMiddleware()) // apply authentication (JWT Token) to inventory microservice

	inventoryServiceEndpoint.GET("/products", getProductsHandlerGRPC)
	inventoryServiceEndpoint.GET("/product/:id", getProductByIDHandlerGRPC)
	inventoryServiceEndpoint.POST("/product", createProductHandlerGRPC)
	inventoryServiceEndpoint.PATCH("/product/:id", updateProductHandlerGRPC)

	return router
}
