package metricsapi

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (ms *MemStorage) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if err := ms.Save(metricType, metricName, metricValue); err != nil {
		if errors.Is(err, ErrInvalidMetricType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrInvalidMetricValue) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.Header().Add("Content-Type", "text-plain")
	w.WriteHeader(http.StatusOK)
}

func (ms *MemStorage) GetHandler(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricType != TypeGauge && metricType != TypeCounter {
		http.Error(w, "invalid metric type", http.StatusBadRequest)
	}

	var res string

	_, v, err := ms.retrieve(metricType, metricName)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if metricType == TypeCounter {
		res = fmt.Sprintf("%d", v)
	} else if metricType == TypeGauge {
		res = fmt.Sprintf("%g", v)
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err = io.WriteString(w, res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ms *MemStorage) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	var result string
	metrics := ms.GetAll()

	for key, v := range metrics {
		result += fmt.Sprintf("%v: %v\n", key, v)
	}

	w.Header().Set("Content-Type", "text/html")

	_, err := io.WriteString(w, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ms *MemStorage) retrieve(metricType, metricName string) (string, interface{}, error) {
	if metricType == TypeCounter {
		value, err := ms.GetCounterMetric(metricName)
		return metricType, value, err
	} else if metricType == TypeGauge {

		value, err := ms.GetGaugeMetric(metricName)
		return metricType, value, err
	}
	return "", nil, ErrInvalidMetricType
}
