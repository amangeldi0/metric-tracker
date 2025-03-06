package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	metricsapi "github.com/amangeldi0/metric-tracker/internal/metrics"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
)

var counter int64

func ReportMetrics(ctx context.Context, client *http.Client, host string, metrics []metricsapi.Metric) error {

	var wg sync.WaitGroup
	errChan := make(chan error, len(metrics))

	for _, m := range metrics {
		m := m
		wg.Add(1)

		go func(m metricsapi.Metric) {
			defer wg.Done()

			jsonMetric, err := json.Marshal(m)
			if err != nil {
				errChan <- err
				return
			}

			cJSONMetric, err := compressData(jsonMetric)
			if err != nil {
				errChan <- err
				return
			}

			endpoint := fmt.Sprintf("http://%s/update/", host)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, cJSONMetric)

			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Accept-Encoding", "gzip")

			// Устанавливаем Content-Length, чтобы сервер корректно обработал запрос
			if cJSONMetric != nil {
				req.ContentLength = int64(cJSONMetric.Size())
			}

			res, err := client.Do(req)
			if err != nil {
				errChan <- err
				return
			}

			err = res.Body.Close()
			if err != nil {
				errChan <- err
				return
			}
		}(m)

	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error

	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors: %v", len(errors), errors)
	}

	return nil
}

func UpdateMetrics() []metricsapi.Metric {
	var metrics []metricsapi.Metric
	var MemStats runtime.MemStats

	runtime.ReadMemStats(&MemStats)

	msValue := reflect.ValueOf(MemStats)
	msType := msValue.Type()

	for _, metric := range metricsapi.GaugeMetrics {
		field, ok := msType.FieldByName(metric)
		if !ok {
			continue
		}

		var value float64

		switch msValue.FieldByName(metric).Interface().(type) {
		case uint64:
			value = float64(msValue.FieldByName(metric).Interface().(uint64))
		case uint32:
			value = float64(msValue.FieldByName(metric).Interface().(uint32))
		case float64:
			value = msValue.FieldByName(metric).Interface().(float64)
		default:
			return nil

		}

		metrics = append(metrics, metricsapi.Metric{ID: field.Name, MType: "gauge", Value: &value})
	}

	counter += 1

	randValue := rand.Float64()
	metrics = append(metrics, metricsapi.Metric{ID: "RandomValue", MType: "gauge", Value: &randValue})
	metrics = append(metrics, metricsapi.Metric{ID: "PollCount", MType: "counter", Delta: &counter})

	return metrics
}

func compressData(data []byte) (*bytes.Reader, error) {
	b := new(bytes.Buffer)
	w, err := gzip.NewWriterLevel(b, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b.Bytes()), nil
}
