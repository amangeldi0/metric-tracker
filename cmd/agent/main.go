package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	metricsapi "github.com/amangeldi0/metric-tracker/api/metrics"
	"github.com/amangeldi0/metric-tracker/internal/config"
	logger2 "github.com/amangeldi0/metric-tracker/internal/lib/logger"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"sync"
	"syscall"
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
	pollInterval := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	client := &http.Client{
		Timeout: time.Minute,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	for {
		select {

		case <-pollInterval.C:
			metrics = updateMetrics()
			logger.Info("Updated metrics", zap.Any("metrics", metrics))

		case <-reportInterval.C:
			if err = reportMetrics(ctx, client, cfg.Address, metrics); err != nil {
				logger.Error("error sending metrics", zap.Error(err))
			} else {
				logger.Info("metrics send")
			}

		case <-ctx.Done():
			logger.Info("shutting down agent...")
			pollInterval.Stop()
			reportInterval.Stop()
			os.Exit(0)
		}
	}
}

func reportMetrics(ctx context.Context, client *http.Client, host string, metrics []metricsapi.Metric) error {

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

			endpoint := fmt.Sprintf("http://%s/update/", host)
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonMetric))

			if err != nil {
				errChan <- err
				return
			}

			req.Header.Set("Content-Type", "application/json")
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
