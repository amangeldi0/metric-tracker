package storage

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

type Mem struct {
	gauge   map[string]float64
	counter map[string]int64

	mx sync.Mutex
}

type MetricType string

var (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"

	ErrInvalidGaugeName   = errors.New("invalid gauge name")
	ErrInvalidCounterName = errors.New("invalid counter name")
)

func NewMem() *Mem {
	return &Mem{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *Mem) GetGauge(name string) (float64, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	name = m.normalizeName(name)

	value, ok := m.gauge[name]
	if !ok {
		return 0, ErrInvalidGaugeName
	}

	return value, nil
}

func (m *Mem) SetGauge(name string, value float64) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	name = m.normalizeName(name)
	m.gauge[name] = value

	return nil
}

func (m *Mem) GetCounter(name string) (int64, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	name = m.normalizeName(name)
	value, ok := m.counter[name]

	if !ok {
		return 0, ErrInvalidCounterName
	}

	return value, nil
}

func (m *Mem) AddCounter(name string, value int64) error {
	name = m.normalizeName(name)
	_, err := m.GetCounter(name)

	m.mx.Lock()
	defer m.mx.Unlock()

	if err != nil {
		m.counter[name] = value
	} else {
		m.counter[name] += value
	}

	return nil
}

func (m *Mem) GetAll() ([]MetricsValue, error) {
	var values []MetricsValue

	for k, v := range m.gauge {
		values = append(values, MetricsValue{
			ID:    k,
			MType: string(GaugeType),
			Value: &v,
		})
	}

	for k, v := range m.counter {
		values = append(values, MetricsValue{
			ID:    k,
			MType: string(CounterType),
			Delta: &v,
		})
	}

	return values, nil
}

func (m *Mem) Ping(_ context.Context) error {
	return nil
}

func (m *Mem) GetMiddleware() gin.HandlerFunc {
	return func(_ *gin.Context) {}
}

func (m *Mem) normalizeName(name string) string {
	return strings.TrimSpace(name)
}

func (m *Mem) Close() error {
	*m = Mem{}
	return nil
}
