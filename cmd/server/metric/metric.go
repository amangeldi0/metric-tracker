package metric

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type MemStorage struct {
	mu       sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
	}
}

func (ms *MemStorage) Gauges() map[string]float64 {
	return ms.gauges
}

func (ms *MemStorage) Counters() map[string]int64 {
	return ms.counters
}

func (ms *MemStorage) UpdateGauge(newValue float64, metricName string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.gauges[metricName] = newValue
	log.Printf("# Server - Updated Gauge - %s: %f", metricName, newValue)
}

func (ms *MemStorage) AddCounter(value int64, metricName string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.counters[metricName] += value
	log.Printf("# Server - Updated Counter - %s: %d", metricName, ms.counters[metricName])
}

func (ms *MemStorage) GetMetrics() map[string]map[string]interface{} {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	metrics := map[string]map[string]interface{}{
		"counters": {},
		"gauges":   {},
	}

	for name, value := range ms.counters {
		metrics["counters"][name] = value
	}
	for name, value := range ms.gauges {
		metrics["gauges"][name] = value
	}

	return metrics
}

func (ms *MemStorage) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Ошибка: Метод не поддерживается")
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	pathSlice := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathSlice) < 4 {
		http.Error(w, "Invalid specified the path", http.StatusNotFound)
		log.Println("Ошибка: Неполный путь в запросе")
		return
	}

	metricType := pathSlice[1]
	metricName := pathSlice[2]
	metricValue := pathSlice[3]

	log.Printf("Обработка метрики: тип=%s, имя=%s, значение=%s", metricType, metricName, metricValue)

	if metricType != CounterType && metricType != GaugeType {
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		log.Println("Ошибка: Неправильный тип метрики")
		return
	}

	if metricType == CounterType {
		count, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			log.Printf("Ошибка: Невалидное значение для counter: %s", metricValue)
			return
		}
		ms.AddCounter(count, metricName)
		log.Printf("Успешно обновлен counter: %s -> %d", metricName, count)
	}

	if metricType == GaugeType {
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid metric value", http.StatusBadRequest)
			log.Printf("Ошибка: Невалидное значение для gauge: %s", metricValue)
			return
		}
		ms.UpdateGauge(value, metricName)
		log.Printf("Успешно обновлен gauge: %s -> %f", metricName, value)
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Запрос обработан успешно")
}
