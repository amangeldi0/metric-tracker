package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	ms := metricsapi.New()

	r := chi.NewRouter()
	r.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/update/counter/test_metric/10", server.URL))

	if err != nil {
		t.Fatalf("Сервер не запустился: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Ожидался код %d, но получили %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
