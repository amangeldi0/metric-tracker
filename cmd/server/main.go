package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"github.com/amangeldi0/metric-tracker/internal/config"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	cfg := config.New()
	ms := metricsapi.New()

	mux.HandleFunc("/update/", ms.UpdateHandler)

	sAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	err := http.ListenAndServe(sAddr, mux)

	if err != nil {
		log.Panic(err)
	}

	log.Printf("server started on %s", sAddr)
}
