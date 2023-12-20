package main

import (
	"github.com/gin-gonic/gin"
)

func (app *Config) routes() *gin.Engine {
	router := gin.Default()

	// middleware for IP whitelisting using go gin-gonic
	router.Use(IPWhitelistMiddleware())

	// define handlers for each microservice, HTTP requests will be forwarded to the microservices
	authenticationHandler := app.handler("http://authentication-service:8001/login")
	signUpHandler := app.handler("http://authentication-service:8001/signup")

	// setting up different paths to handle requests for each microservice
	router.POST("/authenticate", authenticationHandler)
	router.POST("/signup", signUpHandler)

	return router
}
