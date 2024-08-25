package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.Default()

	router.Use(app.ipWhitelistMiddleware())
	router.Use(app.rateLimitMiddleware())

	// for checking if api gateway service is healthy and running
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "api gateway healthy and running!"})
	})

	// authentication microservice endpoints
	router.POST("/authenticate", app.AuthHandler.Login())
	router.POST("/signup", app.AuthHandler.SignUp())
	router.POST("/logout", app.AuthHandler.Logout())

	// authentication (user) microservice endpoints
	router.GET("/users", app.authenticationMiddleware(), app.UserHandler.GetUsers())
	router.PATCH("/user", app.authenticationMiddleware(), app.UserHandler.UpdateUser())

	// inventory microservice endpoints
	inventoryServiceRouter := router.Group("/inventory")
	inventoryServiceRouter.Use(app.authenticationMiddleware())

	inventoryServiceRouter.GET("/products", app.InventoryHandler.GetProducts())
	inventoryServiceRouter.GET("/product/:id", app.InventoryHandler.GetProductByID())
	inventoryServiceRouter.POST("/product", app.InventoryHandler.CreateProduct())
	inventoryServiceRouter.PATCH("/product/:id", app.InventoryHandler.UpdateProduct())
	inventoryServiceRouter.DELETE("/product/:id", app.InventoryHandler.DeleteProduct())

	// order microservice
	orderServiceRouter := router.Group("")
	orderServiceRouter.Use(app.authenticationMiddleware())

	orderServiceRouter.GET("/orders", app.OrderHandler.GetOrders())
	orderServiceRouter.GET("/order/:id", app.OrderHandler.GetOrder())
	orderServiceRouter.POST("/order", app.OrderHandler.CreateOrder())

	return router
}
