package main

import (
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/delivery/http/handlers"
	"github.com/LeonLow97/internal/delivery/http/middleware"
	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/LeonLow97/internal/observability/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type application struct {
	cfg           *config.Config
	HealthHandler *handlers.HealthHandler
	UserHandler   *handlers.UserHandler
}

func (app *application) registerRoutes(logger ports.Logger) *gin.Engine {
	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.RequestLogger(logger),
		middleware.Recovery(logger),
		middleware.MetricsMiddleware(),
	)

	// Health Check endpoint
	healthGroup := router.Group("/health")
	healthGroup.GET("/liveness", app.HealthHandler.Liveness)
	healthGroup.GET("/readiness", app.HealthHandler.Readiness)

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})))

	// User API routes
	userGroup := router.Group("/api/v1")
	userGroup.POST("/login", app.UserHandler.Login)
	userGroup.POST("/signup", app.UserHandler.SignUp)
	userGroup.POST("/logout", app.UserHandler.Logout)

	authenticatedUserGroup := userGroup.Group("")
	authenticatedUserGroup.Use(middleware.Authentication(
		app.cfg.AuthJWTToken.Name,
		app.cfg.AuthJWTToken.Secret,
	))
	authenticatedUserGroup.GET("/users", app.UserHandler.GetUsers)
	authenticatedUserGroup.PATCH("/users/me", app.UserHandler.UpdateUser)

	return router
}
