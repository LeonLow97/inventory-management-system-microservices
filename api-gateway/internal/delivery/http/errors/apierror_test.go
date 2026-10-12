package apierror

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonLow97/internal/observability/logging"
	"github.com/gin-gonic/gin"
)

func TestCustomErrorGetHTTPStatus(t *testing.T) {
	testCases := []struct {
		name       string
		code       string
		wantStatus int
		wantBody   string
	}{
		{name: "bad request", code: "400001", wantStatus: http.StatusBadRequest, wantBody: "Bad Request"},
		{name: "unauthorized", code: "401001", wantStatus: http.StatusUnauthorized, wantBody: "Unauthorized"},
		{name: "forbidden", code: "403001", wantStatus: http.StatusForbidden, wantBody: "Forbidden"},
		{name: "not found", code: "404001", wantStatus: http.StatusNotFound, wantBody: "Not Found"},
		{name: "conflict", code: "409001", wantStatus: http.StatusConflict, wantBody: "Conflict"},
		{name: "too many requests", code: "429001", wantStatus: http.StatusTooManyRequests, wantBody: "Too Many Requests"},
		{name: "internal server error", code: "500001", wantStatus: http.StatusInternalServerError, wantBody: "Internal Server Error"},
		{name: "unknown status", code: "418001", wantStatus: http.StatusInternalServerError, wantBody: "Internal Server Error"},
		{name: "non-numeric status", code: "abc001", wantStatus: http.StatusInternalServerError, wantBody: "Internal Server Error"},
		{name: "short code", code: "40", wantStatus: http.StatusInternalServerError, wantBody: "Internal Server Error"},
		{name: "empty code", code: "", wantStatus: http.StatusInternalServerError, wantBody: "Internal Server Error"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			customErr := &CustomError{Code: testCase.code}
			gotStatus, gotBody := customErr.getHTTPStatus()
			if gotStatus != testCase.wantStatus || gotBody != testCase.wantBody {
				t.Errorf("getHTTPStatus() = (%d, %q), want (%d, %q)", gotStatus, gotBody, testCase.wantStatus, testCase.wantBody)
			}
		})
	}
}

func TestCustomErrorAPIError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name       string
		customErr  *CustomError
		cause      error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "returns generic HTTP message",
			customErr:  &CustomError{Code: "404001", Message: "private detail"},
			cause:      errors.New("database detail"),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"message":"Not Found"}`,
		},
		{
			name:       "malformed code falls back safely",
			customErr:  &CustomError{Code: "x"},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"Internal Server Error"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			testCase.customErr.APIError(ctx, testCase.cause, logging.New(io.Discard))

			if recorder.Code != testCase.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, testCase.wantStatus)
			}
			if gotBody := recorder.Body.String(); gotBody != testCase.wantBody {
				t.Errorf("body = %s, want %s", gotBody, testCase.wantBody)
			}
		})
	}
}
