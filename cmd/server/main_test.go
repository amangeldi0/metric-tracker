package main

import (
	"github.com/amangeldi0/metric-tracker/cmd/server/metric"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandlerInvalidMetricType(t *testing.T) {
	ms := metric.NewMemStorage()
	req := httptest.NewRequest(http.MethodPost, "/update/invalid/testMetric/42.5", nil)
	rr := httptest.NewRecorder()

	updateHandler(rr, req, ms)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rr.Code)
	}
}

func TestRunRoutes(t *testing.T) {
	mux := http.NewServeMux()
	ms := metric.NewMemStorage()

	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		updateHandler(w, r, ms)
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/update/gauge/testMetric/42.5", "text/plain", nil)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", resp.StatusCode)
	}

	if err = resp.Body.Close(); err != nil {
		t.Errorf("Failed to close response body: %v", err)
	}
}

func TestUpdateHandlerInvalidPath(t *testing.T) {
	ms := metric.NewMemStorage()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/testMetric", nil)
	rr := httptest.NewRecorder()

	updateHandler(rr, req, ms)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %v", rr.Code)
	}
}

func TestUpdateHandlerValidGauge(t *testing.T) {
	ms := metric.NewMemStorage()

	req := httptest.NewRequest(http.MethodPost, "/update/gauge/testMetric/42.5", nil)
	rr := httptest.NewRecorder()

	updateHandler(rr, req, ms)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}

	value, exists := (*ms)["testMetric"]
	if !exists {
		t.Errorf("Metric 'testMetric' not found in storage")
	}
	if value.Gauge != 42.5 {
		t.Errorf("Expected gauge value 42.5, got %v", value.Gauge)
	}
}

func TestUpdateHandlerValidCounter(t *testing.T) {
	ms := metric.NewMemStorage()

	req := httptest.NewRequest(http.MethodPost, "/update/counter/testCounter/10", nil)
	rr := httptest.NewRecorder()

	updateHandler(rr, req, ms)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}

	value, exists := (*ms)["testCounter"]
	if !exists {
		t.Errorf("Metric 'testCounter' not found in storage")
	}
	if value.Counter != 10 {
		t.Errorf("Expected counter value 10, got %v", value.Counter)
	}
}

func TestUpdateHandlerInvalidMethod(t *testing.T) {
	ms := metric.NewMemStorage()

	req := httptest.NewRequest(http.MethodGet, "/update/gauge/testMetric/42.5", nil)
	rr := httptest.NewRecorder()

	updateHandler(rr, req, ms)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %v, got %v", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestUpdateHandlerGaugeOverwrite(t *testing.T) {
	ms := metric.NewMemStorage()
	req1 := httptest.NewRequest(http.MethodPost, "/update/gauge/testMetric/42.5", nil)
	rr1 := httptest.NewRecorder()
	updateHandler(rr1, req1, ms)

	req2 := httptest.NewRequest(http.MethodPost, "/update/gauge/testMetric/99.9", nil)
	rr2 := httptest.NewRecorder()
	updateHandler(rr2, req2, ms)

	value := (*ms)["testMetric"]
	if value.Gauge != 99.9 {
		t.Errorf("Expected gauge value 99.9, got %v", value.Gauge)
	}
}
