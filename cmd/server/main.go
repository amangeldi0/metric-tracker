package main

import (
	"github.com/amangeldi0/metric-tracker/cmd/server/metric"
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	ms := metric.NewMemStorage()

	mux.HandleFunc("/update/", ms.UpdateHandler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Panic(err)
	}
}
