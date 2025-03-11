package metrics

import (
	"errors"
	"math/rand"
	"runtime"
)

type (
	MetricType string
	Metric     struct {
		Type  MetricType
		Value interface{}
	}
)

var (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"

	ErrorInvalidPoolCount = errors.New("invalid poll count")
)

type RuntimeMetrics struct {
	runtime     map[string]Metric
	pollCount   Metric
	randomValue Metric
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{
		runtime:     make(map[string]Metric),
		pollCount:   Metric{Type: CounterType, Value: int64(0)},
		randomValue: Metric{Type: GaugeType, Value: float64(0)},
	}
}

func (m *RuntimeMetrics) Update() error {
	var runtimeMetrics runtime.MemStats
	runtime.ReadMemStats(&runtimeMetrics)

	m.runtime["Alloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Alloc)}
	m.runtime["BuckHashSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.BuckHashSys)}
	m.runtime["Frees"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Frees)}
	m.runtime["GCCPUFraction"] = Metric{Type: GaugeType, Value: runtimeMetrics.GCCPUFraction}
	m.runtime["GCSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.GCSys)}
	m.runtime["HeapAlloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapAlloc)}
	m.runtime["HeapIdle"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapIdle)}
	m.runtime["HeapInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapInuse)}
	m.runtime["HeapObjects"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapObjects)}
	m.runtime["HeapReleased"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapReleased)}
	m.runtime["HeapSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.HeapSys)}
	m.runtime["LastGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.LastGC)}
	m.runtime["Lookups"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Lookups)}
	m.runtime["MCacheInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MCacheInuse)}
	m.runtime["MCacheSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MCacheSys)}
	m.runtime["MSpanInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MSpanInuse)}
	m.runtime["MSpanSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.MSpanSys)}
	m.runtime["Mallocs"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Mallocs)}
	m.runtime["NextGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.NextGC)}
	m.runtime["NumForcedGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.NumForcedGC)}
	m.runtime["NumGC"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.NumGC)}
	m.runtime["OtherSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.OtherSys)}
	m.runtime["PauseTotalNs"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.PauseTotalNs)}
	m.runtime["StackInuse"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.StackInuse)}
	m.runtime["StackSys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.StackSys)}
	m.runtime["Sys"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.Sys)}
	m.runtime["TotalAlloc"] = Metric{Type: GaugeType, Value: float64(runtimeMetrics.TotalAlloc)}

	pollCount, ok := m.pollCount.Value.(int64)
	if !ok {
		return ErrorInvalidPoolCount
	}
	m.pollCount.Value = pollCount + 1

	m.randomValue.Value = rand.Float64()
	return nil
}

func (m *RuntimeMetrics) GetRuntime() map[string]Metric {
	return m.runtime
}

func (m *RuntimeMetrics) GetPollCount() Metric {
	return m.pollCount
}

func (m *RuntimeMetrics) GetRandomValue() Metric {
	return m.randomValue
}
