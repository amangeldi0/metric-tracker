package metricsapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func mockLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func mockStorage() *MemStorage {
	return &MemStorage{
		data:   make(map[string]interface{}),
		logger: mockLogger(),
	}
}

func TestUpdateHandler(t *testing.T) {
	ms := mockStorage()
	router := chi.NewRouter()
	router.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	req := httptest.NewRequest("POST", "/update/counter/testCounter/10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetAllHandler(t *testing.T) {
	ms := mockStorage()
	err := ms.SetCounterMetric("testCounter", 50)
	assert.NoError(t, err)

	err = ms.SetGaugeMetric("testGauge", 3.14)
	assert.NoError(t, err)

	router := chi.NewRouter()
	router.Get("/", ms.GetAllHandler)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	body := rr.Body.String()
	assert.Contains(t, body, "testCounter: 50")
	assert.Contains(t, body, "testGauge: 3.14")
}

func TestSave(t *testing.T) {
	ms := mockStorage()

	err := ms.Save("counter", "testCounter", "10")
	assert.NoError(t, err)
	value, _ := ms.GetCounterMetric("testCounter")
	assert.Equal(t, int64(10), value)

}
