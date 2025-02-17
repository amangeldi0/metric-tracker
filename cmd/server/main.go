package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/api/metrics"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/amangeldi0/metric-tracker/internal/http/response"
	liblogger "github.com/amangeldi0/metric-tracker/internal/lib/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func run() error {
	logger, err := liblogger.New(zap.InfoLevel, "server")
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	router := chi.NewRouter()
	sugar := logger.Sugar()

	router.Use(liblogger.WithLogging(sugar))

	router.MethodNotAllowed(response.MethodNotAllowedResponse)
	router.NotFound(response.NotFoundResponse)

	ms := metrics.NewStorage()

	metricHandler := metrics.NewHandler(logger, ms)

	router.Get("/", metricHandler.GetAllHandler)
	router.Get("/value/{metricType}/{metricName}", metricHandler.GetHandler)
	router.Post("/value/", metricHandler.JSONGetHandler)
	router.Post("/update/{metricType}/{metricName}/{metricValue}", metricHandler.UpdateHandler)
	router.Post("/update/", metricHandler.JSONUpdateHandler)

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
