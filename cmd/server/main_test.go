package main

import (
	liblogger "github.com/amangeldi0/metric-tracker/internal/lib/logger"
	"github.com/amangeldi0/metric-tracker/internal/metricsapi"
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

func mockRouter() http.Handler {
	logger := mockLogger()
	ms := metricsapi.New(logger)

	router := chi.NewRouter()
	sugar := logger.Sugar()
	router.Use(liblogger.WithLogging(sugar))

	router.Get("/", ms.GetAllHandler)
	router.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	return router
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := mockRouter()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateMetricHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/update/gauge/cpu_usage/45.7", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := mockRouter()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
