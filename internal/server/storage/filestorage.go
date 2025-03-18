package storage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"strings"
	"time"
)

type FileStorage struct {
	file *os.File

	encoder *json.Encoder
	decoder *json.Decoder

	log logger
	*Mem
}

func NewFileStorage(log logger) (*FileStorage, error) {
	file, err := os.OpenFile(config.Config.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	mem := NewMem()

	fs := &FileStorage{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
		log:     log,
		Mem:     mem,
	}

	return fs, nil
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}

func (s *FileStorage) Restore() error {
	if err := s.file.Sync(); err != nil {
		return err
	}

	var metrics []MetricsValue
	if err := s.decoder.Decode(&metrics); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	var errs int
	for _, metric := range metrics {
		switch metric.MType {
		case string(GaugeType):
			s.Mem.SetGauge(metric.ID, *metric.Value)
		case string(CounterType):
			s.Mem.AddCounter(metric.ID, *metric.Delta)
		default:
			errs++
			s.log.Errorf("The metric couldn't be restored, it has an unknown type: %+v", metrics)
		}
	}

	s.log.Infof("Successfully retrieved metrics (%d) from the file.", len(metrics)-errs)
	return nil
}

func (s *FileStorage) Start() {
	storeInterval := config.Config.StoreInterval
	if storeInterval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(storeInterval))
		for range ticker.C {
			if count, err := s.update(); err != nil {
				s.log.Errorf("Failed to save metrics to file: %s", err)
			} else {
				s.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
			}
		}
	}()
}

func (s *FileStorage) update() (int, error) {
	if err := s.file.Truncate(0); err != nil {
		return 0, err
	}

	if _, err := s.file.Seek(0, 0); err != nil {
		return 0, err
	}

	metrics, _ := s.GetAll()
	if err := s.encoder.Encode(&metrics); err != nil {
		return 0, err
	}

	return len(metrics), nil
}

func (s *FileStorage) GetMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		storeInterval := config.Config.StoreInterval
		if storeInterval > 0 {
			return
		} else if !strings.Contains(ctx.FullPath(), "/update") {
			return
		}

		ctx.Next()

		if count, err := s.update(); err != nil {
			s.log.Errorf("Failed to save metrics to file: %s", err)
		} else {
			s.log.Infof("Metrics (%d) are successfully synchronized and written to file.", count)
		}
	}
}

func (s *FileStorage) Ping(_ context.Context) error {
	_, err := s.file.Stat()
	if os.IsNotExist(err) {
		return err
	}

	return nil
}
