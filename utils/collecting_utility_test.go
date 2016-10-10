package utils

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http/httptest"
	"testing"
	"time"
	"strings"
)

const (
	componentName = "SomeTestComponent"
	gaugeName     = "SomeTestName"
	gaugeHelp = "Some help message"
)

var gauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: gaugeName,
		Help: gaugeHelp,
	},
)

func testCollector() error {
	gauge.Set(42.25)
	return nil
}

func TestE2E(t *testing.T) {
	RegisterMetrics(componentName, gauge)
	done := EnableMetricsCollecting(5*time.Millisecond, testCollector)
	time.Sleep(75 * time.Millisecond)
	done <- struct{}{}

	req := httptest.NewRequest("GET", "http://example.com/metrics", nil)
	w := httptest.NewRecorder()
	handler := GetHandler()

	handler.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Wrong error code, got: %d, expected: 200", w.Code)
	}
	resp := w.Body.String()
	t.Log(resp)

	validateOnCount(t, resp, gaugeName, 3, "No registered metric were in output")
	validateOnCount(t, resp, gaugeName + " 42.25", 1, "No correct value for desired metric were set")
	validateOnCount(t, resp, gaugeHelp, 1, "No help message were in output")
	validateOnCount(t, resp, "tap_metrics_collecting_duration_nanoseconds", 7,
		"Some 'tap_metrics_collecting_duration_nanoseconds' are missing")
	validateOnCount(t, resp, componentName, 8, "Some '" +componentName +"' are missing")
	validateOnCount(t, resp, "tap_metrics_collecting_count", 5, "Some 'tap_metrics_collecting_count' are missing")
}

func validateOnCount(t *testing.T, resp, substr string, expectedCount int, errMsg string) {
	count := strings.Count(resp, substr)
	if  count != expectedCount {
		t.Errorf(errMsg + "(got: %d, expected %d)", count, expectedCount)
	}
}
