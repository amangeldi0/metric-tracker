package metricsupdater

import (
	"errors"
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/agent/config"
	"github.com/amangeldi0/metric-tracker/internal/agent/metrics"
	"github.com/amangeldi0/metric-tracker/internal/agent/models"
	"github.com/go-resty/resty/v2"
)

var (
	ErrorInvalidMetricValueType = errors.New("invalid metric value type")
)

type (
	Updater struct {
		client *resty.Client
		store  store
		log    logger
	}

	logger interface {
		Errorf(template string, args ...interface{})
	}

	store interface {
		GetRuntime() map[string]metrics.Metric
		GetPollCount() metrics.Metric
		GetRandomValue() metrics.Metric
	}
)

func New(client *resty.Client, store store, log logger) *Updater {
	return &Updater{
		client: client,
		store:  store,
		log:    log,
	}
}

func (u Updater) UpdateMetrics() {
	for k, v := range u.store.GetRuntime() {
		if err := u.updateMetric(k, v); err != nil {
			u.log.Errorf("Failed to update metric \"%s\": %s", k, err)
		}
	}

	if err := u.updateMetric("PollCount", u.store.GetPollCount()); err != nil {
		u.log.Errorf("Failed to update metric \"PollCount\": %s", err)
	}

	if err := u.updateMetric("RandomValue", u.store.GetRandomValue()); err != nil {
		u.log.Errorf("Failed to update metric \"PollCount\": %s", err)
	}
}

func (u Updater) updateMetric(name string, metric metrics.Metric) error {
	body, err := u.parseMetric(name, metric)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/update", config.Config.Address)
	_, err = u.client.R().
		SetBody(body).
		Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (u Updater) parseMetric(name string, metric metrics.Metric) (*models.Metrics, error) {
	var obj models.Metrics
	obj.ID = name

	switch metric.Value.(type) {
	case float64:
		if metric.Type != metrics.GaugeType {
			return nil, ErrorInvalidMetricValueType
		}
		value := metric.Value.(float64)

		obj.MType = string(metrics.GaugeType)
		obj.Value = &value
	case int64:
		if metric.Type != metrics.CounterType {
			return nil, ErrorInvalidMetricValueType
		}
		delta := metric.Value.(int64)

		obj.MType = string(metrics.CounterType)
		obj.Delta = &delta
	default:
		return nil, ErrorInvalidMetricValueType
	}

	return &obj, nil
}
