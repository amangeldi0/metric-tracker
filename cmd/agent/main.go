package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
)

func main() {

	var pollCount atomic.Int64

	collector := GetGaugeMetricMaps()

	go func() {
		for {
			collector = GetGaugeMetricMaps()
			pollCount.Add(1)
			time.Sleep(pollInterval)
		}
	}()

	for {
		PostMetrics(collector, pollCount.Load())
		time.Sleep(reportInterval)
	}
}

func PostMetrics(metrics map[string]float64, pollCount int64) {
	baseURL := "http://127.0.0.1:8080/update"

	client := &http.Client{}

	for k, v := range metrics {
		s := strconv.FormatFloat(v, 'f', -1, 64)
		url := baseURL + "/" + "gauge" + "/" + strings.ToLower(k) + "/" + s

		r, err := client.Post(url, "text/plain", nil)
		if err != nil {
			log.Printf("PostMetrics: Error posting metrics: %s", err.Error())
			continue
		}

		if r.StatusCode != http.StatusOK {
			log.Printf("PostMetrics: Server returned non-OK status: %d", r.StatusCode)
		}

		fmt.Printf(" - %s - %f \n", k, v)

		if err := r.Body.Close(); err != nil {
			log.Printf("PostMetrics: Error closing body: %s", err.Error())
		}
	}

	url := baseURL + "/" + "counter" + "/" + "pollcount" + "/" + strconv.FormatInt(pollCount, 10)

	r, err := client.Post(url, "text/plain", nil)

	if err != nil {
		log.Printf("PostMetrics: Error posting metrics: %s", err.Error())
	}

	if r.StatusCode != http.StatusOK {
		log.Printf("PostMetrics: Server returned non-OK status: %d", r.StatusCode)
	}

	fmt.Println()

	log.Println("PostMetrics: Success")

	if err := r.Body.Close(); err != nil {
		log.Printf("PostMetrics: Error closing body: %s", err.Error())
	}
}
func GetGaugeMetricMaps() map[string]float64 {

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": memStats.GCCPUFraction,
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}
	return metrics
}
