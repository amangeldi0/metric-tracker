package main

import (
	"flag"
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
	counter int64
)

func main() {
	var metrics []Metric

	cfg := config.New()

	defaultAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	sAddr := flag.String("a", defaultAddr, "input agent address ex: localhost:8080")

	pollInterval := flag.Duration("p", 2*time.Second, "input pollInterval")
	reportInterval := flag.Duration("r", 10*time.Second, "input pollInterval report interval")

	flag.Parse()

	go func() {
		for {
			metrics = updateMetrics()
			time.Sleep(*pollInterval)
		}
	}()
	for {
		err := reportMetrics(metrics, *sAddr)
		fmt.Println("Metrics reported:")
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(*reportInterval)
	}
}

func reportMetrics(metrics []Metric, url string) error {

	for _, m := range metrics {
		client := &http.Client{}
		endpoint := fmt.Sprintf("http://%s/%s/%s/%v", url, m.Type, m.Name, m.Value)

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
