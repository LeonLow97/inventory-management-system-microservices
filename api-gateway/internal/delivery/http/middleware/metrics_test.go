package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonLow97/internal/observability/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const route = "/metrics-middleware-test"
	beforeInFlight := testutil.ToFloat64(metrics.InFlightRequests)
	requestMetric := metrics.HTTPRequests.WithLabelValues(http.MethodGet, route, "201")
	beforeRequests := testutil.ToFloat64(requestMetric)
	router := gin.New()
	router.Use(MetricsMiddleware())
	router.GET(route, func(c *gin.Context) {
		if got, want := testutil.ToFloat64(metrics.InFlightRequests), beforeInFlight+1; got != want {
			t.Errorf("in-flight requests inside handler = %v, want %v", got, want)
		}
		c.Status(http.StatusCreated)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, route, nil))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := testutil.ToFloat64(metrics.InFlightRequests); got != beforeInFlight {
		t.Errorf("in-flight requests after handler = %v, want %v", got, beforeInFlight)
	}
	if got := testutil.ToFloat64(requestMetric); got != beforeRequests+1 {
		t.Errorf("request counter = %v, want %v", got, beforeRequests+1)
	}
}
