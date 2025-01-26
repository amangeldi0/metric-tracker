package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

var (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	counter        int64
)

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
		fmt.Println("Metrics reported:")
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(reportInterval)
	}
}

func reportMetrics(metrics []Metric) error {

	cfg := config.New()

	url := fmt.Sprintf("%s://%s:%d/metrics", cfg.Server.Protocol, cfg.Server.Host, cfg.Server.Port)

	for _, m := range metrics {
		client := &http.Client{}
		endpoint := fmt.Sprintf("%s/%s/%s/%v", url, m.Type, m.Name, m.Value)

		req, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			return err
		}

		req.Header.Add("Content-Type", "text/plain")
		res, err := client.Do(req)
		if err != nil {
			return err
		}

		err = res.Body.Close()
		if err != nil {
			return err
		}
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
