package app

import (
	"github.com/amangeldi0/metric-tracker/internal/agent/config"
	"github.com/amangeldi0/metric-tracker/internal/agent/metrics"
	metricsupdater "github.com/amangeldi0/metric-tracker/internal/agent/metrics_updater"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) {
		if err = logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)

	sugarLogger := logger.Sugar()

	config.Load()
	if err = config.Parse(); err != nil {
		sugarLogger.Panicf("Failed loading config: %s", err)
	}

	client := resty.New()
	store := metrics.NewRuntimeMetrics()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(config.Config.PollInterval))

			if err = store.Update(); err != nil {
				sugarLogger.Panicf("Failed to update metrics: %s", err)
			}
		}
	}()

	updater := metricsupdater.New(client, store, sugarLogger)
	for {
		time.Sleep(time.Second * time.Duration(config.Config.ReportInterval))
		updater.UpdateMetrics()
	}
}
