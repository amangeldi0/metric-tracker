package metricsapi

import (
	"errors"
	"go.uber.org/zap"
)

var (
	GaugeMetrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"RandomValue",
	}
)

type Metric struct {
	Name  string
	Type  string
	Value interface{}
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

var (
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrNotFound           = errors.New("value not found")
	ErrCannotAssign       = errors.New("cannot assign value, key is already in use by another metric type")
)

const (
	TypeCounter = "counter"
	TypeGauge   = "gauge"
)

type MemStorage struct {
	data   map[string]interface{}
	logger *zap.Logger
}
