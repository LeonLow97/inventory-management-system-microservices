package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/observability/logging"
	"github.com/gin-gonic/gin"
)

func testLogger() ports.Logger {
	return logging.New(io.Discard)
}

func newTestContext(method, target string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(method, target, nil)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	return ctx, recorder
}

func TestRespondError(t *testing.T) {
	testCases := []struct {
		name       string
		err        error
		wantStatus int
		wantKind   coreerrors.Kind
		wantBody   string
	}{
		{name: "invalid argument", err: coreerrors.New(coreerrors.KindInvalidArgument, "bad input", nil), wantStatus: http.StatusBadRequest, wantKind: coreerrors.KindInvalidArgument, wantBody: `{"error":"bad input"}`},
		{name: "unauthenticated hides details", err: coreerrors.New(coreerrors.KindUnauthenticated, "detail", nil), wantStatus: http.StatusUnauthorized, wantKind: coreerrors.KindUnauthenticated, wantBody: `{"error":"Unauthorized"}`},
		{name: "permission denied", err: coreerrors.New(coreerrors.KindPermissionDenied, "forbidden", nil), wantStatus: http.StatusForbidden, wantKind: coreerrors.KindPermissionDenied, wantBody: `{"error":"forbidden"}`},
		{name: "not found", err: coreerrors.New(coreerrors.KindNotFound, "missing", nil), wantStatus: http.StatusNotFound, wantKind: coreerrors.KindNotFound, wantBody: `{"error":"missing"}`},
		{name: "conflict", err: coreerrors.New(coreerrors.KindConflict, "conflict", nil), wantStatus: http.StatusConflict, wantKind: coreerrors.KindConflict, wantBody: `{"error":"conflict"}`},
		{name: "already exists", err: coreerrors.New(coreerrors.KindAlreadyExists, "duplicate", nil), wantStatus: http.StatusConflict, wantKind: coreerrors.KindAlreadyExists, wantBody: `{"error":"duplicate"}`},
		{name: "failed precondition", err: coreerrors.New(coreerrors.KindFailedPrecondition, "blocked", nil), wantStatus: http.StatusPreconditionFailed, wantKind: coreerrors.KindFailedPrecondition, wantBody: `{"error":"blocked"}`},
		{name: "out of range", err: coreerrors.New(coreerrors.KindOutOfRange, "too large", nil), wantStatus: http.StatusBadRequest, wantKind: coreerrors.KindOutOfRange, wantBody: `{"error":"too large"}`},
		{name: "unavailable hides details", err: coreerrors.New(coreerrors.KindUnavailable, "private", nil), wantStatus: http.StatusServiceUnavailable, wantKind: coreerrors.KindUnavailable, wantBody: `{"error":"Service temporarily unavailable"}`},
		{name: "deadline exceeded hides details", err: coreerrors.New(coreerrors.KindDeadlineExceeded, "private", nil), wantStatus: http.StatusGatewayTimeout, wantKind: coreerrors.KindDeadlineExceeded, wantBody: `{"error":"Request timed out"}`},
		{name: "rate limited hides details", err: coreerrors.New(coreerrors.KindResourceExhausted, "private", nil), wantStatus: http.StatusTooManyRequests, wantKind: coreerrors.KindResourceExhausted, wantBody: `{"error":"Rate limit exceeded"}`},
		{name: "internal hides details", err: coreerrors.New(coreerrors.KindInternal, "private", nil), wantStatus: http.StatusInternalServerError, wantKind: coreerrors.KindInternal, wantBody: `{"error":"Internal server error"}`},
		{name: "unclassified error", err: errors.New("private"), wantStatus: http.StatusInternalServerError, wantKind: coreerrors.KindInternal, wantBody: `{"error":"Internal server error"}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, recorder := newTestContext(http.MethodGet, "/")
			handler := newBaseHandler(testLogger())
			gotKind := handler.RespondError(ctx, testCase.err)
			if recorder.Code != testCase.wantStatus {
				t.Errorf("HTTP status = %d, want %d", recorder.Code, testCase.wantStatus)
			}
			if gotKind != testCase.wantKind {
				t.Errorf("error kind = %q, want %q", gotKind, testCase.wantKind)
			}
			if body := strings.TrimSpace(recorder.Body.String()); body != testCase.wantBody {
				t.Errorf("response body = %s, want %s", body, testCase.wantBody)
			}
		})
	}
}

func TestGetQueryParam(t *testing.T) {
	testCases := []struct {
		name        string
		target      string
		wantDefault string
	}{
		{name: "query value", target: "/users?limit=25", wantDefault: "25"},
		{name: "missing value uses default", target: "/users", wantDefault: "10"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, _ := newTestContext(http.MethodGet, testCase.target)
			handler := newBaseHandler(testLogger())

			if got := handler.GetQueryParam(ctx, "limit", "10"); got != testCase.wantDefault {
				t.Errorf("GetQueryParam() = %q, want %q", got, testCase.wantDefault)
			}
		})
	}
}

func TestResponseHelpers(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		wantBody   string
		noContent  bool
	}{
		{name: "JSON envelope", statusCode: http.StatusOK, wantBody: `{"data":{"name":"Ada"}}`},
		{name: "no content", statusCode: http.StatusNoContent, noContent: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, recorder := newTestContext(http.MethodGet, "/")
			handler := newBaseHandler(testLogger())
			if testCase.noContent {
				handler.ResponseNoContent(ctx, testCase.statusCode)
			} else {
				handler.ResponseJSON(ctx, testCase.statusCode, gin.H{"name": "Ada"})
			}
			if got := ctx.Writer.Status(); got != testCase.statusCode {
				t.Errorf("status = %d, want %d", got, testCase.statusCode)
			}
			if body := strings.TrimSpace(recorder.Body.String()); body != testCase.wantBody {
				t.Errorf("body = %s, want %s", body, testCase.wantBody)
			}
		})
	}
}

type readinessCheckerStub struct {
	err        error
	gotContext bool
}

func (s *readinessCheckerStub) Check(ctx context.Context) error {
	_, s.gotContext = ctx.Deadline()
	return s.err
}

func TestHealthHandlers(t *testing.T) {
	testCases := []struct {
		name       string
		handler    string
		pingErr    error
		wantStatus int
		wantBody   string
	}{
		{name: "liveness", handler: "live", wantStatus: http.StatusOK, wantBody: `{"message":"api gateway is running","status":"alive"}`},
		{name: "readiness succeeds", handler: "ready", wantStatus: http.StatusOK, wantBody: `{"message":"all systems operational","status":"ready"}`},
		{name: "readiness reports dependency failure", handler: "ready", pingErr: errors.New("connection refused"), wantStatus: http.StatusServiceUnavailable, wantBody: `{"error":"connection refused","status":"unhealthy"}`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			checker := &readinessCheckerStub{err: testCase.pingErr}
			handler := NewHealthHandler(checker)
			ctx, recorder := newTestContext(http.MethodGet, "/health")
			if testCase.handler == "live" {
				handler.Liveness(ctx)
			} else {
				handler.Readiness(ctx)
				if !checker.gotContext {
					t.Error("readiness check did not receive a deadline context")
				}
			}
			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, testCase.wantStatus)
			}
			if body := strings.TrimSpace(recorder.Body.String()); body != testCase.wantBody {
				t.Errorf("body = %s, want %s", body, testCase.wantBody)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "valid password", password: "Aa123456!", wantErr: false},
		{name: "missing special character", password: "Aa123456", wantErr: true},
		{name: "password too short", password: "Aa1!", wantErr: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := struct {
				Password string `validate:"required,password_format,min=8,max=20"`
			}{Password: testCase.password}
			err := newBaseHandler(testLogger()).ValidateStruct(request)
			if (err != nil) != testCase.wantErr {
				t.Errorf("ValidateStruct() error = %v, wantErr %t", err, testCase.wantErr)
			}
		})
	}
}
