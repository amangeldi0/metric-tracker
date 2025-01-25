package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/go-resty/resty/v2"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

var (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	counter        int64
)

type MyAPIError struct {
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	var metrics []Metric
	go func() {
		for {
			metrics = updateMetrics()
			time.Sleep(pollInterval)
		}
	}()
	for {
		err := reportMetrics(metrics)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(reportInterval)
	}
}

func reportMetrics(metrics []Metric) error {

	client := resty.New()
	cfg := config.New()

	var responseErr MyAPIError

	url := fmt.Sprintf("%s://%s:%d/metrics", cfg.Server.Protocol, cfg.Server.Host, cfg.Server.Port)

	for _, m := range metrics {
		endpoint := fmt.Sprintf("%s/%s/%s/%v", url, m.Type, m.Name, m.Value)

		res, err := client.R().
			SetError(&responseErr).
			Post(endpoint)

		if err != nil {
			fmt.Println(responseErr)
			return err
		}

		res.Header().Add("Content-Type", "text/plain")

	}
	return nil
}

func updateMetrics() []Metric {
	var metrics []Metric
	var MemStats runtime.MemStats

	runtime.ReadMemStats(&MemStats)
	msValue := reflect.ValueOf(MemStats)
	msType := msValue.Type()

	for _, metric := range GaugeMetrics {
		field, ok := msType.FieldByName(metric)
		if !ok {
			continue
		}
		value := msValue.FieldByName(metric)
		metrics = append(metrics, Metric{Name: field.Name, Type: "gauge", Value: value})
	}

	counter += 1
	metrics = append(metrics, Metric{Name: "RandomValue", Type: "gauge", Value: rand.Float64()})
	metrics = append(metrics, Metric{Name: "PollCounter", Type: "counter", Value: counter})

	return metrics
}
