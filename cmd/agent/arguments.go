package main

import (
	"flag"
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/caarlos0/env/v6"
	"time"
)

type EnvConfig struct {
	Addr           string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func getStartArguments() (string, time.Duration, time.Duration) {
	cfg := config.New()

	defaultAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	var envCfg EnvConfig
	err := env.Parse(&envCfg)

	sAddr := flag.String("a", defaultAddr, "input agent address ex: localhost:8080")
	pollIntervalSeconds := flag.Int("p", 2, "input pollInterval in seconds")
	reportIntervalSeconds := flag.Int("r", 10, "input reportInterval in seconds")

	if err != nil {
		return *sAddr,
			time.Duration(*pollIntervalSeconds) * time.Second,
			time.Duration(*reportIntervalSeconds) * time.Second
	}

	return envCfg.Addr, envCfg.PollInterval, envCfg.ReportInterval

}
