package main

import (
	"context"
	"flag"
	"github.com/amangeldi0/metric-tracker/internal/agent"
	"github.com/amangeldi0/metric-tracker/internal/config"
	logger2 "github.com/amangeldi0/metric-tracker/internal/logger"
	metricsapi "github.com/amangeldi0/metric-tracker/internal/metrics"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	counter int64
)

func main() {

	logger, err := logger2.New(zap.InfoLevel, "agent")

	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err := config.New()

	if err != nil {
		logger.Info("Failed to create config", zap.Error(err))
	}

	if len(flag.Args()) > 0 {
		log.Fatal("unknown flags", zap.Any("flags", flag.Args()))
	}

	var metrics []metricsapi.Metric
	pollInterval := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	reportInterval := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	client := &http.Client{
		Timeout: time.Minute,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	for {
		select {

		case <-pollInterval.C:
			metrics = agent.UpdateMetrics()
			logger.Info("Updated metrics", zap.Any("metrics", metrics))

		case <-reportInterval.C:
			if err = agent.ReportMetrics(ctx, client, cfg.Address, metrics); err != nil {
				logger.Error("error sending metrics", zap.Error(err))
			} else {
				logger.Info("metrics send")
			}

		case <-ctx.Done():
			logger.Info("shutting down agent...")
			pollInterval.Stop()
			reportInterval.Stop()
			os.Exit(0)
		}
	}
}
