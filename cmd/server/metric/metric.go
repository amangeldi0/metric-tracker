package metric

import "log"

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MemValue struct {
	counter int64
	gauge   float64
}

type MemStorage map[string]MemValue

func NewMemStorage() *MemStorage {
	storage := make(MemStorage)
	return &storage
}

func (m MemStorage) UpdateGauge(newValue float64, metricName string) {
	m[metricName] = MemValue{
		gauge: newValue,
	}

	log.Println(metricName, newValue)
}

func (m MemStorage) AddCounter(value int64, metricName string) {

	c, ok := m[metricName]

	if !ok {
		m[metricName] = MemValue{
			counter: value,
		}
	} else {
		m[metricName] = MemValue{
			counter: c.counter + value,
		}
	}

	log.Println(metricName, value)
}
