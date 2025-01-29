package main

import (
	"flag"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cMux := chi.NewMux()
	ms := metricsapi.New()
	LoadConfig()

	if err := ParseConfig(); err != nil {
		panic(err)
	}

	cMux.Get("/", ms.GetAllHandler)
	cMux.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	cMux.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	flag.Parse()
	log.Printf("server started on %s", Config.Address)
	err := http.ListenAndServe(Config.Address, cMux)

	if err != nil {
		log.Panic(err)
	}
}
