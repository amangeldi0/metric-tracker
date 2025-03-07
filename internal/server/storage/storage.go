package storage

import (
	"errors"
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"strings"
)

type Mem struct {
	gauge   map[string]float64
	counter map[string]int64
}

type MetricType string

var (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"

	ErrInvalidGaugeName   = errors.New("invalid gauge name")
	ErrInvalidCounterName = errors.New("invalid counter name")
)

func NewMem() Mem {
	return Mem{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *Mem) GetGauge(name string) (float64, error) {
	name = m.normalizeName(name)

	value, ok := m.gauge[name]
	if !ok {
		return 0, ErrInvalidGaugeName
	}

	return value, nil
}

func (m *Mem) SetGauge(name string, value float64) {
	name = m.normalizeName(name)
	m.gauge[name] = value
}

func (m *Mem) GetCounter(name string) (int64, error) {
	name = m.normalizeName(name)

	value, ok := m.counter[name]
	if !ok {
		return 0, ErrInvalidCounterName
	}

	return value, nil
}

func (m *Mem) AddCounter(name string, value int64) {
	name = m.normalizeName(name)

	if _, err := m.GetCounter(name); err != nil {
		m.counter[name] = value
	} else {
		m.counter[name] += value
	}
}

func (m *Mem) GetAll() []models.MetricsValue {
	var values []models.MetricsValue

	for k, v := range m.gauge {
		values = append(values, models.MetricsValue{
			ID:    k,
			MType: string(GaugeType),
			Value: &v,
		})
	}

	for k, v := range m.counter {
		values = append(values, models.MetricsValue{
			ID:    k,
			MType: string(CounterType),
			Delta: &v,
		})
	}

	return values
}

func (m *Mem) normalizeName(name string) string {
	return strings.TrimSpace(name)
}
