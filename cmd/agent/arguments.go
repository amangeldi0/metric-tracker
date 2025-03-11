package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

var Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func LoadConfig() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
	flag.IntVar(&Config.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&Config.PollInterval, "p", 2, "poll interval")
}

func ParseConfig() error {
	flag.Parse()

	return env.Parse(&Config)
}
