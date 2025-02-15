package main

import (
	"flag"
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/config"
	logger2 "github.com/amangeldi0/metric-tracker/internal/lib/logger"
	"github.com/amangeldi0/metric-tracker/internal/metricsapi"
	"go.uber.org/zap"
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

	logger, err := logger2.New(zap.InfoLevel, "agent")

	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.New()

	if err != nil {
		logger.Info("Failed to create config", zap.Error(err))
	}

	if len(flag.Args()) > 0 {
		log.Fatal("unknown flags", zap.Any("flags", flag.Args()))
	}

	var metrics []metricsapi.Metric
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second

	go func() {
		for {
			metrics = updateMetrics()
			time.Sleep(pollInterval)
		}
	}()
	for {
		err = reportMetrics(metrics, cfg.Address)
		logger.Info("reported metrichandlers", zap.Any("metrichandlers", metrics))

		if err != nil {
			logger.Error("failed to report metrichandlers", zap.Error(err))
		}

		time.Sleep(reportInterval)
	}
}

func reportMetrics(metrics []metricsapi.Metric, url string) error {

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

func updateMetrics() []metricsapi.Metric {
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
		value := msValue.FieldByName(metric)
		metrics = append(metrics, metricsapi.Metric{Name: field.Name, Type: "gauge", Value: value})
	}

	counter += 1
	metrics = append(metrics, metricsapi.Metric{Name: "RandomValue", Type: "gauge", Value: rand.Float64()})
	metrics = append(metrics, metricsapi.Metric{Name: "PollCounter", Type: "counter", Value: counter})

	return metrics
}
