package main

import (
	"context"
	"encoding/json"
	metricsapi "github.com/amangeldi0/metric-tracker/api/metrics"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestUpdateMetrics(t *testing.T) {
	metrics := updateMetrics()

	assert.NotEmpty(t, metrics, "Метрики не должны быть пустыми")

	var foundPollCount, foundRandomValue bool

	for _, m := range metrics {
		if m.ID == "PollCount" {
			foundPollCount = true
			assert.Equal(t, "counter", m.MType, "PollCount должен быть типа counter")
			assert.NotNil(t, m.Delta, "PollCount должен иметь Delta значение")
		}
		if m.ID == "RandomValue" {
			foundRandomValue = true
			assert.Equal(t, "gauge", m.MType, "RandomValue должен быть типа gauge")
			assert.NotNil(t, m.Value, "RandomValue должен иметь Value значение")
		}
	}

	assert.True(t, foundPollCount, "PollCount не найден среди метрик")
	assert.True(t, foundRandomValue, "RandomValue не найден среди метрик")
}

func TestReportMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Метод запроса должен быть POST")
		assert.Equal(t, "/update/", r.URL.Path, "Некорректный путь запроса")

		var metric metricsapi.Metric
		err := json.NewDecoder(r.Body).Decode(&metric)
		assert.NoError(t, err, "Ошибка при разборе JSON тела запроса")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &http.Client{Timeout: 10 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	metrics := []metricsapi.Metric{
		{ID: "TestMetric", MType: "gauge", Value: new(float64)},
	}

	err := reportMetrics(ctx, client, server.Listener.Addr().String(), metrics)
	assert.NoError(t, err, "Ошибка при отправке метрик")
}
