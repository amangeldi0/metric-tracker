package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

var (
	ReportInterval = 10
	PollInterval   = 2
	DefaultAddress = "localhost:8080"
)

func New() (*Config, error) {

	var config Config

	flag.StringVar(&config.Address, "a", DefaultAddress, "server address")
	flag.IntVar(&config.ReportInterval, "r", ReportInterval, "report interval")
	flag.IntVar(&config.PollInterval, "p", PollInterval, "poll interval")

	flag.Parse()

	err := env.Parse(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
