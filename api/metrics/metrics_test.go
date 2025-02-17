package metrics

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONUpdateHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	storage := NewStorage()
	h := NewHandler(logger, storage)

	metric := Metric{
		ID:    "test_metric",
		MType: TypeGauge,
		Value: func(v float64) *float64 { return &v }(10.5),
	}

	body, _ := json.Marshal(metric)
	r := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.JSONUpdateHandler(w, r)
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Errorf("failed to close response body: %s", err)
		}
	}(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestJSONGetHandler(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	storage := NewStorage()
	h := NewHandler(logger, storage)

	metric := Metric{
		ID:    "test_metric",
		MType: TypeCounter,
		Delta: func(v int64) *int64 { return &v }(5),
	}

	err := storage.Save(metric)
	if err != nil {
		t.Errorf("failed to save metric: %s", err)
	}

	body, _ := json.Marshal(metric)
	r := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.JSONGetHandler(w, r)
	res := w.Result()
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			t.Errorf("failed to close response body: %s", err)
		}
	}(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	responseBody, _ := io.ReadAll(res.Body)
	var returnedMetric Metric
	err = json.Unmarshal(responseBody, &returnedMetric)
	if err != nil {
		t.Errorf("failed to unmarshal response: %s", err)
	}
	assert.Equal(t, *metric.Delta, *returnedMetric.Delta)
}
