package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	ReportInterval  int    `env:"REPORT_INTERVAL"`
	PollInterval    int    `env:"POLL_INTERVAL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

var (
	ReportInterval  = 10
	PollInterval    = 2
	DefaultAddress  = "localhost:8080"
	StoreInterval   = 300
	FileStoragePath = ""
	Restore         = false
)

func New() (*Config, error) {

	var config Config

	flag.StringVar(&config.Address, "a", DefaultAddress, "server address")
	flag.IntVar(&config.ReportInterval, "r", ReportInterval, "report interval")
	flag.IntVar(&config.PollInterval, "p", PollInterval, "poll interval")
	flag.StringVar(&config.FileStoragePath, "f", FileStoragePath, "file storage path")
	flag.BoolVar(&config.Restore, "r", Restore, "restore")
	flag.IntVar(&config.StoreInterval, "s", StoreInterval, "store interval")

	flag.Parse()

	err := env.Parse(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
