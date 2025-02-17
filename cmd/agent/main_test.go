package main

import (
	"github.com/amangeldi0/metric-tracker/internal/metricsapi"
	"reflect"
	"runtime"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	metrics := updateMetrics()

	if len(metrics) == 0 {
		t.Errorf("expected non-empty metrichandlers, got %d", len(metrics))
	}

	found := false
	for _, metric := range metrics {
		if metric.ID == "RandomValue" && metric.MType == "gauge" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected metric RandomValue of type gauge to be present")
	}
}

func TestGaugeMetricsIntegrity(t *testing.T) {
	var MemStats runtime.MemStats
	runtime.ReadMemStats(&MemStats)
	msValue := reflect.ValueOf(MemStats)
	msType := msValue.Type()

	for _, metric := range metricsapi.GaugeMetrics {
		_, ok := msType.FieldByName(metric)
		if !ok {
			if metric == "RandomValue" {
				continue
			}
			t.Errorf("metric %s not found in runtime.MemStats", metric)
		}
	}
}
