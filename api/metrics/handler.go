package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/http/response"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	logger  *zap.Logger
	storage *MemStorage
}

func NewHandler(logger *zap.Logger, storage *MemStorage) *Handler {
	return &Handler{
		logger:  logger,
		storage: storage,
	}
}

func (h *Handler) JSONUpdateHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			response.ServerErrorResponse(w, r, err)
			return
		}
	}(r.Body)

	var jsonMetric Metric

	err := dec.Decode(&jsonMetric)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	err = h.storage.Save(jsonMetric)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	err = h.storage.Retrieve(&jsonMetric)
	if err != nil {
		response.ServerErrorResponse(w, r, err)
		return
	}

	response.WriteJson(w, http.StatusOK, jsonMetric)
}

func (h *Handler) JSONGetHandler(w http.ResponseWriter, r *http.Request) {
	jsonDec := json.NewDecoder(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			response.ServerErrorResponse(w, r, err)
			return
		}
	}(r.Body)

	var jsonMetric Metric

	err := jsonDec.Decode(&jsonMetric)
	if err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	err = h.storage.Retrieve(&jsonMetric)
	if err != nil {
		h.logger.Debug("error retrieving value", zap.String("value", jsonMetric.ID), zap.Error(err))
		response.NotFoundResponse(w, r)
		return
	}

	response.WriteJson(w, http.StatusOK, jsonMetric)
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var metrics Metric

	metrics.MType = chi.URLParam(r, "metricType")
	metrics.ID = chi.URLParam(r, "metricName")
	value := chi.URLParam(r, "metricValue")

	if metrics.MType == "gauge" {

		gaugeValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			response.BadRequestResponse(w, r, err)
			return
		}
		metrics.Value = &gaugeValue

	} else if metrics.MType == "counter" {
		counterValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			response.BadRequestResponse(w, r, err)
			return
		}

		metrics.Delta = &counterValue

	} else {
		response.BadRequestResponse(w, r, errors.New("invalid metric type"))
		return
	}

	if err := h.storage.Save(metrics); err != nil {
		response.BadRequestResponse(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "text-plain")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	var metrics Metric
	var result string

	metrics.MType = chi.URLParam(r, "metricType")
	metrics.ID = chi.URLParam(r, "metricName")

	err := h.storage.Retrieve(&metrics)
	if err != nil {
		h.logger.Debug("error retrieving value", zap.Error(err))
		response.NotFoundResponse(w, r)
		return
	}

	if metrics.MType == TypeGauge {
		result = fmt.Sprintf("%g", *metrics.Value)
	} else if metrics.MType == TypeCounter {
		result = fmt.Sprintf("%d", *metrics.Delta)
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err = io.WriteString(w, result)

	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	metrics := h.storage.GetAllMetric()

	for k, v := range metrics {
		result += fmt.Sprintf("%s: %s\n", k, v)
	}

	w.Header().Set("Content-Type", "text/html")
	_, err := io.WriteString(w, result)

	if err != nil {
		response.ServerErrorResponse(w, r, err)
	}
}
