package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonLow97/internal/delivery/http/handlers"
	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/LeonLow97/internal/observability/logging"
	"github.com/gin-gonic/gin"
)

type routeReadinessStub struct{}

func (routeReadinessStub) Check(_ context.Context) error { return nil }

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := logging.New(io.Discard)
	cfg := &config.Config{}
	cfg.AuthJWTToken.Name = "auth-token"
	cfg.AuthJWTToken.Secret = "route-test-secret"
	app := application{
		cfg:           cfg,
		HealthHandler: handlers.NewHealthHandler(routeReadinessStub{}),
		UserHandler: handlers.NewUserHandler(handlers.JWTConfig{
			Name:     cfg.AuthJWTToken.Name,
			MaxAge:   cfg.AuthJWTToken.MaxAge,
			Path:     cfg.AuthJWTToken.Path,
			Domain:   cfg.AuthJWTToken.Domain,
			Secure:   cfg.AuthJWTToken.Secure,
			HTTPOnly: cfg.AuthJWTToken.HTTPOnly,
		}, nil, logger),
	}
	router := app.registerRoutes(logger)

	testCases := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{name: "liveness route", method: http.MethodGet, target: "/health/liveness", wantStatus: http.StatusOK},
		{name: "readiness route", method: http.MethodGet, target: "/health/readiness", wantStatus: http.StatusOK},
		{name: "login route is registered", method: http.MethodPost, target: "/api/v1/login", wantStatus: http.StatusBadRequest},
		{name: "users route requires authentication", method: http.MethodGet, target: "/api/v1/users", wantStatus: http.StatusUnauthorized},
		{name: "update route requires authentication", method: http.MethodPatch, target: "/api/v1/users/me", wantStatus: http.StatusUnauthorized},
		{name: "unknown route remains not found", method: http.MethodGet, target: "/missing", wantStatus: http.StatusNotFound},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(testCase.method, testCase.target, nil)
			router.ServeHTTP(recorder, request)
			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", recorder.Code, testCase.wantStatus, recorder.Body.String())
			}
		})
	}
}
