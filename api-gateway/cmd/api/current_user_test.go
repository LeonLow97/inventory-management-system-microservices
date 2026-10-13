package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/LeonLow97/internal/contextstore"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/usecases"
	"github.com/LeonLow97/internal/delivery/http/handlers"
	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/LeonLow97/internal/observability/logging"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type currentUserRouteService struct {
	usecases.User
	id int64
}

func (s *currentUserRouteService) GetCurrentUser(ctx context.Context) (*domain.User, error) {
	s.id, _ = contextstore.UserIDFromContext(ctx)
	return &domain.User{ID: s.id, Email: "ada@example.com", FirstName: "Ada", LastName: "Lovelace", Password: "must-not-leak", Active: true}, nil
}

func TestCurrentUserRouteUsesCookieIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{}
	cfg.AuthJWTToken.Name = "auth-token"
	cfg.AuthJWTToken.Secret = "test-secret"
	logger := logging.New(io.Discard)
	service := &currentUserRouteService{}
	app := application{cfg: cfg, HealthHandler: handlers.NewHealthHandler(routeReadinessStub{}), UserHandler: handlers.NewUserHandler(handlers.JWTConfig{}, service, logger)}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 42, "session_id": "session", "iss": "authify-service", "exp": time.Now().Add(time.Hour).Unix()}).SignedString([]byte(cfg.AuthJWTToken.Secret))
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/users/me?user_id=999", nil)
	request.AddCookie(&http.Cookie{Name: cfg.AuthJWTToken.Name, Value: token})
	recorder := httptest.NewRecorder()
	app.registerRoutes(logger).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || service.id != 42 {
		t.Fatalf("status = %d, identity = %d", recorder.Code, service.id)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, `"first_name":"Ada"`) || !strings.Contains(body, `"last_name":"Lovelace"`) || strings.Contains(body, "password") || strings.Contains(body, "token") {
		t.Fatalf("unsafe or incorrect response: %s", body)
	}
	if recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("profile must not be cached by HTTP caches")
	}
}
