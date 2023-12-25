package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const AUTHENTICATION_SERVICE_URL = "authentication-service:8001"

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.Use(app.ipWhitelistMiddleware())
	router.Use(app.rateLimitMiddleware())

	// gRPC Communication with Authentication service
	authenticationHandlerGRPC := app.gRPCAuthenticationHandler(AUTHENTICATION_SERVICE_URL)
	signUpHandlerGRPC := app.gRPCSignUpHandler(AUTHENTICATION_SERVICE_URL)

	updateUserHandlerGRPC := app.grpcUpdateUserHandler(AUTHENTICATION_SERVICE_URL)
	getUsersHandlerGRPC := app.grpcGetUsersHandler(AUTHENTICATION_SERVICE_URL)

	// for pinging and testing the api gateway
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "api gateway healthy and running!"})
	})

	// setting up different paths to handle requests for each microservice
	router.POST("/authenticate", authenticationHandlerGRPC)
	router.POST("/signup", signUpHandlerGRPC)

	router.PATCH("/user", updateUserHandlerGRPC)
	router.GET("/users", getUsersHandlerGRPC)

	return router
}
