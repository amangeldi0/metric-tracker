package metricsapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (ms *MemStorage) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {

	var req Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.MType != TypeGauge && req.MType != TypeCounter {
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	if req.MType == TypeCounter {
		if req.Delta == nil {
			http.Error(w, "missing delta for counter", http.StatusBadRequest)
			return
		}
		err := ms.SetCounterMetric(req.ID, *req.Delta)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else if req.MType == TypeGauge {
		if req.Value == nil {
			http.Error(w, "missing value for gauge", http.StatusBadRequest)
			return
		}
		err := ms.SetGaugeMetric(req.ID, *req.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	_, v, err := ms.retrieve(req.MType, req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Metrics{
		ID:    req.ID,
		MType: req.MType,
	}
	if req.MType == TypeCounter {
		val := v.(int64)
		resp.Delta = &val
	} else if req.MType == TypeGauge {
		val := v.(float64)
		resp.Value = &val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ms *MemStorage) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	var req Metrics
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.MType != TypeGauge && req.MType != TypeCounter {
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	_, v, err := ms.retrieve(req.MType, req.ID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Metrics{
		ID:    req.ID,
		MType: req.MType,
	}
	if req.MType == TypeCounter {
		val := v.(int64)
		resp.Delta = &val
	} else if req.MType == TypeGauge {
		val := v.(float64)
		resp.Value = &val
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ms *MemStorage) GetAllHandler(w http.ResponseWriter, _ *http.Request) {
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
