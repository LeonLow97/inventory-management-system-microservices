package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.Use(app.ipWhitelistMiddleware())
	router.Use(app.rateLimitMiddleware())

	// gRPC Communication with Authentication service
	// define handlers for each microservice, HTTP requests will be forwarded to the microservices
	// authenticationHandler := app.handler("http://authentication-service:8001/login")
	authenticationHandlerGRPC := app.gRPCAuthenticationHandler("authentication-service:8001")
	// signUpHandler := app.handler("http://authentication-service:8001/signup")

	// updateUserHandler := app.handler("http://authentication-service:8001/user")
	// getUsersHandler := app.handler("http://authentication-service:8001/users")

	// for pinging and testing the api gateway
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "api gateway healthy and running!"})
	})

	// setting up different paths to handle requests for each microservice
	router.POST("/authenticate", authenticationHandlerGRPC)
	// router.POST("/signup", signUpHandler)

	// router.PATCH("/user", updateUserHandler)
	// router.GET("/users", getUsersHandler)

	return router
}
