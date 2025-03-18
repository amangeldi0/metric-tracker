package storage

import (
	"context"
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsUpdate struct {
	ID    string   `json:"id" binding:"required"`
	MType string   `json:"type" binding:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty" binding:"required_if=MType counter"`
	Value *float64 `json:"value,omitempty" binding:"required_if=MType gauge"`
}

type MetricsValue struct {
	ID    string   `json:"id" binding:"required"`
	MType string   `json:"type" binding:"required,oneof=counter gauge"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type Storage interface {
	SetGauge(string, float64) error
	AddCounter(string, int64) error

	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)

	GetAll() ([]MetricsValue, error)
	GetMiddleware() gin.HandlerFunc
	Ping(context.Context) error
	Close() error
}

func SetupStorage(log logger) (Storage, error) {
	if config.Config.DatabaseDSN != "" {
		db, err := pgxpool.New(context.Background(), config.Config.DatabaseDSN)

		if err != nil {
			return nil, err
		}

		return NewPsStorage(db, log)
	}

	if config.Config.FileStoragePath != "" {
		fs, err := NewFileStorage(log)
		if err != nil {
			return nil, err
		}

		if err = fs.Restore(); err != nil {
			return nil, err
		}
		fs.Start()

		return fs, nil
	}

	return NewMem(), nil
}
