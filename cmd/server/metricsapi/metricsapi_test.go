package metricsapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	ms := New()

	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
	}{
		{
			name:       "Valid Counter Metric",
			method:     http.MethodPost,
			url:        "/update/counter/test_metric/10",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Valid Gauge Metric",
			method:     http.MethodPost,
			url:        "/update/gauge/test_metric/3.14",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Metric Type",
			method:     http.MethodPost,
			url:        "/update/invalid/test_metric/10",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid Metric Value",
			method:     http.MethodPost,
			url:        "/update/counter/test_metric/abc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Method Not Allowed",
			method:     http.MethodGet,
			url:        "/update/counter/test_metric/10",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "Invalid URL Path",
			method:     http.MethodPost,
			url:        "/update/counter/test_metric",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			recorder := httptest.NewRecorder()

			ms.UpdateHandler(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}

func TestMemStorage(t *testing.T) {
	ms := New()

	// Test SetCounterMetric
	err := ms.SetCounterMetric("counter_metric", 10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test GetCounterMetric
	counterValue, err := ms.GetCounterMetric("counter_metric")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if counterValue != 10 {
		t.Errorf("got %d, want 10", counterValue)
	}

	// Test Increment Counter Metric
	err = ms.SetCounterMetric("counter_metric", 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	counterValue, err = ms.GetCounterMetric("counter_metric")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if counterValue != 15 {
		t.Errorf("got %d, want 15", counterValue)
	}

	// Test SetGaugeMetric
	err = ms.SetGaugeMetric("gauge_metric", 3.14)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test GetGaugeMetric
	gaugeValue, err := ms.GetGaugeMetric("gauge_metric")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if gaugeValue != 3.14 {
		t.Errorf("got %f, want 3.14", gaugeValue)
	}

	// Test Error Cases
	err = ms.SetGaugeMetric("counter_metric", 2.71)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
