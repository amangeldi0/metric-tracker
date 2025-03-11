package metricsapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUpdateHandler(t *testing.T) {
	ms := New()
	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

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

			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}
}
