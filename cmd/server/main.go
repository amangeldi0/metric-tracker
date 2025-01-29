package main

import (
	"flag"
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type EnvConfig struct {
	Address string `env:"ADDRESS"`
}

func main() {
	cMux := chi.NewMux()
	cfg := config.New()
	ms := metricsapi.New()
	var envCfg EnvConfig

	err := env.Parse(&envCfg)

	if err != nil {
		log.Printf("failed to parse env vars: %v", err)
	} else {
		log.Printf("parsed env vars: %+v", envCfg)
	}

	var addr string
	defaultAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	if envCfg.Address == "" {
		addr = *flag.String("a", defaultAddr, "input address ex: localhost:8080")
	} else {
		addr = envCfg.Address
	}

	cMux.Get("/", ms.GetAllHandler)
	cMux.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	cMux.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	flag.Parse()
	log.Printf("server started on %s", addr)
	err = http.ListenAndServe(addr, cMux)

	if err != nil {
		log.Panic(err)
	}
}
