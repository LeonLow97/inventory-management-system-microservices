package metrics

import (
	"sync"
	"testing"
)

var initializeMetricsTestOnce sync.Once

func TestInitMetricsRegistersApplicationAndRuntimeMetrics(t *testing.T) {
	HTTPRequests.WithLabelValues("GET", "/metrics-init-test", "200").Inc()
	HTTPReqLatency.WithLabelValues("GET", "/metrics-init-test").Observe(0.01)
	initializeMetricsTestOnce.Do(InitMetrics)

	families, err := Registry.Gather()
	if err != nil {
		t.Fatalf("Registry.Gather() error = %v", err)
	}
	found := make(map[string]bool, len(families))
	for _, family := range families {
		found[family.GetName()] = true
	}

	for _, metricName := range []string{
		"http_requests_total",
		"http_request_duration_seconds",
		"http_inflight_requests",
		"go_goroutines",
		"process_open_fds",
	} {
		if !found[metricName] {
			t.Errorf("registered metrics do not include %q", metricName)
		}
	}
}
