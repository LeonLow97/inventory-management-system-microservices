package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/gin-gonic/gin"
)

type userServiceStub struct {
	loginUser   *domain.User
	loginToken  string
	loginErr    error
	loginReq    domain.User
	loginCalls  int
	signUpErr   error
	signUpReq   domain.User
	signUpCalls int
	users       []domain.User
	nextCursor  string
	getUsersErr error
	getUsersReq struct {
		limit  int64
		cursor string
	}
	getUsersCalls int
	updateErr     error
	updateReq     domain.User
	updateCalls   int
}

func (s *userServiceStub) Login(_ context.Context, req domain.User, _ string) (*domain.User, string, error) {
	s.loginCalls++
	s.loginReq = req
	return s.loginUser, s.loginToken, s.loginErr
}

func (s *userServiceStub) SignUp(_ context.Context, req domain.User) error {
	s.signUpCalls++
	s.signUpReq = req
	return s.signUpErr
}

func (s *userServiceStub) GetUsers(_ context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	s.getUsersCalls++
	s.getUsersReq.limit = limit
	s.getUsersReq.cursor = cursor
	return s.users, s.nextCursor, s.getUsersErr
}

func (s *userServiceStub) UpdateUser(_ context.Context, req domain.User) error {
	s.updateCalls++
	s.updateReq = req
	return s.updateErr
}

func newUserHandlerForTest(service *userServiceStub) *UserHandler {
	cfg := JWTConfig{
		Name:     "auth-token",
		MaxAge:   3600,
		Path:     "/",
		HTTPOnly: true,
	}
	return NewUserHandler(cfg, service, testLogger())
}

func sendUserRequest(handler gin.HandlerFunc, method, target, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	router.Handle(method, request.URL.Path, handler)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestUserHandlerLogin(t *testing.T) {
	testCases := []struct {
		name           string
		body           string
		serviceErr     error
		wantStatus     int
		wantLoginCalls int
	}{
		{name: "successful login", body: `{"email":"ada@example.com","password":"Secret1!"}`, wantStatus: http.StatusOK, wantLoginCalls: 1},
		{name: "malformed JSON", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "invalid request fields", body: `{"email":"invalid","password":"weak"}`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"email":"ada@example.com","password":"Secret1!"}`, serviceErr: coreerrors.New(coreerrors.KindUnauthenticated, "bad credentials", nil), wantStatus: http.StatusUnauthorized, wantLoginCalls: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{loginUser: &domain.User{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Active: true, Admin: true}, loginToken: "signed-token", loginErr: testCase.serviceErr}
			handler := newUserHandlerForTest(service)
			recorder := sendUserRequest(handler.Login, http.MethodPost, "/login", testCase.body)

			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", recorder.Code, testCase.wantStatus, recorder.Body.String())
			}
			if service.loginCalls != testCase.wantLoginCalls {
				t.Errorf("login calls = %d, want %d", service.loginCalls, testCase.wantLoginCalls)
			}
			if testCase.wantStatus == http.StatusOK {
				cookies := recorder.Result().Cookies()
				if len(cookies) != 1 || cookies[0].Name != "auth-token" || cookies[0].Value != "signed-token" || !cookies[0].HttpOnly {
					t.Errorf("login cookie = %#v, want configured auth cookie", cookies)
				}
				if strings.Contains(recorder.Body.String(), "active") {
					t.Errorf("login response must not expose Active: %s", recorder.Body.String())
				}
			}
		})
	}
}

func TestUserHandlerSignUp(t *testing.T) {
	testCases := []struct {
		name       string
		body       string
		serviceErr error
		wantStatus int
		wantCalls  int
	}{
		{name: "successful sign up", body: `{"first_name":"Ada","last_name":"Lovelace","email":"ada@example.com","password":"Secret1!"}`, wantStatus: http.StatusCreated, wantCalls: 1},
		{name: "invalid request", body: `{"first_name":"","email":"bad"}`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"first_name":"Ada","last_name":"Lovelace","email":"ada@example.com","password":"Secret1!"}`, serviceErr: coreerrors.New(coreerrors.KindAlreadyExists, "duplicate", nil), wantStatus: http.StatusConflict, wantCalls: 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{signUpErr: testCase.serviceErr}
			handler := newUserHandlerForTest(service)
			recorder := sendUserRequest(handler.SignUp, http.MethodPost, "/signup", testCase.body)
			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", recorder.Code, testCase.wantStatus, recorder.Body.String())
			}
			if service.signUpCalls != testCase.wantCalls {
				t.Errorf("signup calls = %d, want %d", service.signUpCalls, testCase.wantCalls)
			}
			if testCase.wantCalls == 1 && service.signUpReq.Email != "ada@example.com" {
				t.Errorf("signup request email = %q, want ada@example.com", service.signUpReq.Email)
			}
		})
	}
}

func TestUserHandlerLogout(t *testing.T) {
	handler := newUserHandlerForTest(&userServiceStub{})
	recorder := sendUserRequest(handler.Logout, http.MethodPost, "/logout", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "auth-token" || cookies[0].MaxAge >= 0 || cookies[0].Value != "" {
		t.Errorf("logout cookie = %#v, want expired auth cookie", cookies)
	}
}

func TestUserHandlerGetUsers(t *testing.T) {
	testCases := []struct {
		name       string
		target     string
		serviceErr error
		wantStatus int
		wantCalls  int
	}{
		{name: "defaults pagination and returns users", target: "/users", wantStatus: http.StatusOK, wantCalls: 1},
		{name: "passes pagination query", target: "/users?limit=25&cursor=next", wantStatus: http.StatusOK, wantCalls: 1},
		{name: "invalid limit", target: "/users?limit=many", wantStatus: http.StatusBadRequest},
		{name: "service error", target: "/users", serviceErr: coreerrors.New(coreerrors.KindUnavailable, "private", nil), wantStatus: http.StatusServiceUnavailable, wantCalls: 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{users: []domain.User{{ID: 7, FirstName: "Ada", Email: "ada@example.com"}}, nextCursor: "next-page", getUsersErr: testCase.serviceErr}
			handler := newUserHandlerForTest(service)
			recorder := sendUserRequest(handler.GetUsers, http.MethodGet, testCase.target, "")
			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", recorder.Code, testCase.wantStatus, recorder.Body.String())
			}
			if service.getUsersCalls != testCase.wantCalls {
				t.Errorf("get users calls = %d, want %d", service.getUsersCalls, testCase.wantCalls)
			}
			if testCase.wantCalls == 1 {
				wantLimit, wantCursor := int64(10), ""
				if testCase.target == "/users?limit=25&cursor=next" {
					wantLimit, wantCursor = 25, "next"
				}
				if service.getUsersReq.limit != wantLimit || service.getUsersReq.cursor != wantCursor {
					t.Errorf("pagination = (%d, %q), want (%d, %q)", service.getUsersReq.limit, service.getUsersReq.cursor, wantLimit, wantCursor)
				}
			}
		})
	}
}

func TestUserHandlerUpdateUser(t *testing.T) {
	testCases := []struct {
		name       string
		body       string
		serviceErr error
		wantStatus int
		wantCalls  int
	}{
		{name: "successful update", body: `{"first_name":"Ada","password":"Secret1!"}`, wantStatus: http.StatusNoContent, wantCalls: 1},
		{name: "malformed JSON", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "invalid field", body: `{"first_name":"A very long first name that exceeds fifty characters in total"}`, wantStatus: http.StatusBadRequest},
		{name: "service error", body: `{"last_name":"Lovelace"}`, serviceErr: coreerrors.New(coreerrors.KindNotFound, "missing", nil), wantStatus: http.StatusNotFound, wantCalls: 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{updateErr: testCase.serviceErr}
			handler := newUserHandlerForTest(service)
			recorder := sendUserRequest(handler.UpdateUser, http.MethodPatch, "/users/me", testCase.body)
			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d; body=%s", recorder.Code, testCase.wantStatus, recorder.Body.String())
			}
			if service.updateCalls != testCase.wantCalls {
				t.Errorf("update calls = %d, want %d", service.updateCalls, testCase.wantCalls)
			}
			if testCase.name == "successful update" && (service.updateReq.FirstName != "Ada" || service.updateReq.Password != "Secret1!") {
				t.Errorf("update request = %#v", service.updateReq)
			}
		})
	}
}
