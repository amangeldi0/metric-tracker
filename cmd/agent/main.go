package main

import (
	"bytes"
	"encoding/json"
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

	var metrics []metricsapi.Metrics
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

func reportMetrics(metrics []metricsapi.Metrics, url string) error {

	for _, m := range metrics {
		client := &http.Client{}
		endpoint := fmt.Sprintf("http://%s/%s/%s/%v", url, m.MType, m.ID, m.Value)

		data, err := json.Marshal(m)

		if err != nil {
			return fmt.Errorf("failed to marshal metric: %w", err)
		}

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(data))
		req.Header.Add("Content-Type", "application/json")

		if err != nil {
			return err
		}

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

func updateMetrics() []metricsapi.Metrics {
	var metrics []metricsapi.Metrics
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

		if !value.IsValid() || value.Kind() != reflect.Float64 {
			continue
		}

		val := value.Float()

		metrics = append(metrics, metricsapi.Metrics{ID: field.Name, MType: "gauge", Value: &val})
	}

	counter += 1
	randValue := rand.Float64()
	metrics = append(metrics, metricsapi.Metrics{ID: "RandomValue", MType: "gauge", Value: &randValue})
	metrics = append(metrics, metricsapi.Metrics{ID: "PollCounter", MType: "counter", Delta: &counter})

	return metrics
}
