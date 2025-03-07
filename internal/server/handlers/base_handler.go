package handlers

import (
	"github.com/amangeldi0/metric-tracker/internal/server/models"
)

type (
	BaseHandler struct {
		storage stor
		log     logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}

	stor interface {
		GetGauge(name string) (float64, error)
		SetGauge(name string, value float64)
		GetCounter(name string) (int64, error)
		AddCounter(name string, value int64)
		GetAll() []models.MetricsValue
	}
)

func NewBase(storage stor, log logger) *BaseHandler {
	return &BaseHandler{storage: storage, log: log}
}
