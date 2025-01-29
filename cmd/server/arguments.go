package main

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

var Config struct {
	Address string `env:"ADDRESS"`
}

func LoadConfig() {
	flag.StringVar(&Config.Address, "a", "localhost:8080", "server address")
}

func ParseConfig() error {
	flag.Parse()

	return env.Parse(&Config)
}
