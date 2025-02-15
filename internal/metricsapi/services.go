package metricsapi

import (
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"strconv"
)

func New(logger *zap.Logger) *MemStorage {
	return &MemStorage{
		data:   make(map[string]interface{}),
		logger: logger,
	}
}

func (ms *MemStorage) SetCounterMetric(key string, value int64) error {
	currentValue, exists := ms.data[key]

	if exists {
		if reflect.TypeOf(currentValue) == reflect.TypeOf(int64(0)) {
			ms.data[key] = currentValue.(int64) + value
			ms.logger.Info(fmt.Sprintf("metrics counter metric set to %d", value))
			return nil
		} else {
			ms.logger.Warn("metrics counter metric already exists", zap.String("key", key))
			return ErrCannotAssign
		}
	}
	ms.data[key] = value
	ms.logger.Info(fmt.Sprintf("metrics counter metric set to %d", value))

	return nil
}

func (ms *MemStorage) GetAll() map[string]interface{} {
	return ms.data
}

func (ms *MemStorage) SetGaugeMetric(key string, value float64) error {
	currentValue := ms.data[key]
	if reflect.TypeOf(currentValue) == reflect.TypeOf(int64(0)) {
		ms.logger.Warn("failed set gauge metric")
		return ErrCannotAssign
	}

	ms.logger.Info(fmt.Sprintf("metrics gauge metric set to %f", value))
	ms.data[key] = value
	return nil
}

func (ms *MemStorage) GetCounterMetric(key string) (int64, error) {
	v, ok := ms.data[key]
	if !ok {
		return 0, ErrNotFound
	}
	return v.(int64), nil
}

func (ms *MemStorage) GetGaugeMetric(key string) (float64, error) {
	v, ok := ms.data[key]
	if !ok {
		ms.logger.Warn("failed get gauge metric")
		return 0, ErrNotFound
	}
	return v.(float64), nil
}

func (ms *MemStorage) Save(metricType string, key string, value string) error {
	if metricType == TypeCounter {
		intMetric, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			ms.logger.Warn("got invalid counter metric", zap.String("key", key), zap.String("value", value))
			return ErrInvalidMetricValue
		}
		return ms.SetCounterMetric(key, intMetric)

	} else if metricType == TypeGauge {
		floatMetric, err := strconv.ParseFloat(value, 64)
		if err != nil {
			ms.logger.Warn("got invalid gauge metric", zap.String("key", key), zap.String("value", value))
			return ErrInvalidMetricValue
		}

		return ms.SetGaugeMetric(key, floatMetric)
	}

	ms.logger.Warn("unsupported metric type", zap.String("type", metricType))
	return ErrInvalidMetricType
}
