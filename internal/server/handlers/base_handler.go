package handlers

import (
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
)

type (
	BaseHandler struct {
		storage storage.Storage
		log     logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}
)

func NewBase(storage storage.Storage, log logger) *BaseHandler {
	return &BaseHandler{storage: storage, log: log}
}
