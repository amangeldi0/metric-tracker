package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/config"
	liblogger "github.com/amangeldi0/metric-tracker/internal/lib/logger"
	"github.com/amangeldi0/metric-tracker/internal/metricsapi"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func run() error {
	logger, err := liblogger.New(zap.InfoLevel, "agent")
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	ms := metricsapi.New(logger)

	router := chi.NewRouter()
	sugar := logger.Sugar()

	router.Use(liblogger.WithLogging(sugar))

	router.Get("/", ms.GetAllHandler)
	router.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	logger.Info("server started", zap.String("address", cfg.Address))

	return srv.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
