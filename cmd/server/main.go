package main

import (
	"flag"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cMux := chi.NewMux()
	ms := metricsapi.New()
	cfg, err := config.New()

	if err != nil {
		log.Fatal(err)
	}

	cMux.Get("/", ms.GetAllHandler)
	cMux.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	cMux.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	flag.Parse()
	log.Printf("server started on %s", cfg.Address)
	err = http.ListenAndServe(cfg.Address, cMux)

	if err != nil {
		log.Panic(err)
	}
}
