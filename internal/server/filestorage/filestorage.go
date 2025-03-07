package filestorage

import (
	"encoding/json"
	"errors"
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
	"time"
)

type (
	Storage struct {
		file *os.File

		encoder *json.Encoder
		decoder *json.Decoder

		storage store
		log     logger
	}

	logger interface {
		Infof(template string, args ...interface{})
		Errorf(template string, args ...interface{})
	}

	store interface {
		SetGauge(name string, value float64)
		AddCounter(name string, value int64)
		GetAll() []models.MetricsValue
	}
)

func New(storage store, log logger) (*Storage, error) {
	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Storage{
		file: file,

		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),

		storage: storage,
		log:     log,
	}, nil
}

func (u *Storage) Close() error {
	return u.file.Close()
}

func (u *Storage) Restore() error {
	if err := u.file.Sync(); err != nil {
		return err
	}

	var metrics []models.MetricsValue
	if err := u.decoder.Decode(&metrics); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	var errs int
	for _, metric := range metrics {
		switch metric.MType {
		case string(storage.GaugeType):
			u.storage.SetGauge(metric.ID, *metric.Value)
		case string(storage.CounterType):
			u.storage.AddCounter(metric.ID, *metric.Delta)
		default:
			errs++
			u.log.Errorf("The metric couldn't be restored, it has an unknown type: %+v", metrics)
		}
	}

	u.log.Infof("Successfully retrieved metrics (%d) from the file.", len(metrics)-errs)
	return nil
}

func (u *Storage) Start() {
	storeInterval := config.Config.StoreInterval
	if storeInterval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(storeInterval))
		for range ticker.C {
			if count, err := u.update(); err != nil {
				u.log.Errorf("Failed to save metrics to file: %s", err)
			} else {
				u.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
			}
		}
	}()
}

func (u *Storage) update() (int, error) {
	if err := u.file.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := u.file.Seek(0, 0); err != nil {
		return 0, err
	}

	metrics := u.storage.GetAll()
	if err := u.encoder.Encode(&metrics); err != nil {
		return 0, err
	}

	return len(metrics), nil
}

func (u *Storage) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		storeInterval := config.Config.StoreInterval
		if storeInterval > 0 {
			return
		} else if !strings.Contains(ctx.FullPath(), "/update") {
			return
		}

		ctx.Next()

		if count, err := u.update(); err != nil {
			u.log.Errorf("Failed to save metrics to file: %s", err)
		} else {
			u.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
		}
	}
}
