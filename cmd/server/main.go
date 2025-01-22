package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metric"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	mux := http.NewServeMux()

	ms := metric.NewMemStorage()

	mux.HandleFunc("POST /update/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		updateHandler(w, r, ms)
	})

	if err := run(mux); err != nil {
		panic(err)
	}

}

func run(handler http.Handler) error {
	return http.ListenAndServe(`:8080`, handler)
}

func updateHandler(w http.ResponseWriter, r *http.Request, metricStorage *metric.MemStorage) {
	w.Header().Set("Content-Type", "text/plain")

	pathSlice := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	if len(pathSlice) < 4 {
		http.Error(w, "Invalid specified the path", http.StatusNotFound)
		return
	}

	fmt.Println(pathSlice)

	metricType := pathSlice[1]
	metricName := pathSlice[2]
	metricValue := pathSlice[3]

	if metricType != metric.CounterType && metricType != metric.GaugeType {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	if metricType == metric.CounterType {
		count, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		metricStorage.AddCounter(count, metricName)

	}

	if metricType == metric.GaugeType {
		gauge, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			return
		}

		metricStorage.UpdateGauge(gauge, metricName)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Metric updated successfully")); err != nil {
		panic(err)
	}
}
