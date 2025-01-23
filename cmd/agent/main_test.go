package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var baseURL = "http://127.0.0.1:8080/update"

func TestGetGaugeMetricMaps(t *testing.T) {
	metrics := GetGaugeMetricMaps()

	if len(metrics) == 0 {
		t.Fatal("Expected non-empty metrics map")
	}

	expectedKeys := []string{"Alloc", "HeapAlloc", "NumGC", "RandomValue"}
	for _, key := range expectedKeys {
		if _, exists := metrics[key]; !exists {
			t.Errorf("Expected key %s to be present in metrics", key)
		}
	}
}

func TestPostMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if !strings.HasPrefix(r.URL.Path, "/update/") {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	metrics := map[string]float64{
		"Alloc":       1234.56,
		"RandomValue": 0.789,
	}
	pollCount := int64(42)

	baseURL = server.URL
	PostMetrics(metrics, pollCount)
}

func TestPostMetricsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	metrics := map[string]float64{
		"Alloc": 1234.56,
	}
	pollCount := int64(42)

	baseURL = server.URL
	PostMetrics(metrics, pollCount)
}
