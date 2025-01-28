package main

import (
	"flag"
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {

	cMux := chi.NewMux()
	cfg := config.New()
	ms := metricsapi.New()

	cMux.Get("/", ms.GetAllHandler)
	cMux.Get("/value/{metricType}/{metricName}", ms.GetHandler)
	cMux.Post("/update/{metricType}/{metricName}/{metricValue}", ms.UpdateHandler)

	defaultAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	sAddr := flag.String("a", defaultAddr, "input  address ex: localhost:8080")

	flag.Parse()
	log.Printf("server started on %s", *sAddr)
	err := http.ListenAndServe(*sAddr, cMux)

	if err != nil {
		log.Panic(err)
	}
}
