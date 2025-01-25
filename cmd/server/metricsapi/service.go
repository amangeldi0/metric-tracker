package metricsapi

import "strconv"

func (ms *MemStorage) Save(metricType string, key string, value string) error {
	if metricType == TypeCounter {
		intMetric, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return ErrInvalidMetricValue
		}
		return ms.SetCounterMetric(key, intMetric)

	} else if metricType == TypeGauge {
		floatMetric, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return ErrInvalidMetricValue
		}
		return ms.SetGaugeMetric(key, floatMetric)
	}
	return ErrInvalidMetricType
}
