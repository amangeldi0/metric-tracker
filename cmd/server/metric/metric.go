package metric

import "log"

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MemValue struct {
	Counter int64
	Gauge   float64
}

type MemStorage map[string]MemValue

func NewMemStorage() *MemStorage {
	storage := make(MemStorage)
	return &storage
}

func (m MemStorage) UpdateGauge(newValue float64, metricName string) {
	m[metricName] = MemValue{
		Gauge: newValue,
	}

	log.Println(metricName, newValue)
}

func (m MemStorage) AddCounter(value int64, metricName string) {

	c, ok := m[metricName]

	if !ok {
		m[metricName] = MemValue{
			Counter: value,
		}
	} else {
		m[metricName] = MemValue{
			Counter: c.Counter + value,
		}
	}

	log.Println(metricName, value)
}
