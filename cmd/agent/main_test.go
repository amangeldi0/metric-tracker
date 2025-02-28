package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	metricsapi "github.com/amangeldi0/metric-tracker/api/metrics"
	"github.com/stretchr/testify/assert"
	"io"
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
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Должен быть JSON Content-Type")
		assert.Equal(t, "gzip", r.Header.Get("Content-Encoding"), "Должен быть gzip Content-Encoding")

		gz, err := gzip.NewReader(r.Body)
		assert.NoError(t, err, "Ошибка при разжатии gzip")
		defer gz.Close()

		var metric metricsapi.Metric
		err = json.NewDecoder(gz).Decode(&metric)
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

func TestReportMetrics_Error(t *testing.T) {
	client := &http.Client{Timeout: 2 * time.Second}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	metrics := []metricsapi.Metric{
		{ID: "TestMetric", MType: "gauge", Value: new(float64)},
	}

	err := reportMetrics(ctx, client, "http://localhost:9999", metrics)
	assert.Error(t, err, "Должна возникнуть ошибка при отправке метрик на недоступный сервер")
}

func TestCompressData(t *testing.T) {
	data := []byte(`{"id":"TestMetric","type":"gauge","value":123.45}`)

	compressed, err := compressData(data)
	assert.NoError(t, err, "Ошибка при сжатии данных")

	gz, err := gzip.NewReader(compressed)
	assert.NoError(t, err, "Ошибка при создании gzip-ридера")
	defer gz.Close()

	uncompressedData, err := io.ReadAll(gz)
	assert.NoError(t, err, "Ошибка при чтении разжатых данных")

	assert.Equal(t, data, uncompressedData, "Разжатые данные не совпадают с исходными")
}
