package main

import (
	"github.com/amangeldi0/metric-tracker/cmd/server/metric"
	"testing"
)

func TestUpdateGauge(t *testing.T) {
	ms := metric.NewMemStorage()

	ms.UpdateGauge(10.5, "testGauge")
	ms.UpdateGauge(20.2, "testGauge") // Обновляем значение

	if value, exists := ms.Gauges()["testGauge"]; !exists {
		t.Fatalf("Expected gauge 'testGauge' to exist")
	} else if value != 20.2 {
		t.Errorf("Expected gauge value 20.2, got %f", value)
	}
}

func TestAddCounter(t *testing.T) {
	ms := metric.NewMemStorage()

	// Добавляем новую counter метрику
	ms.AddCounter(5, "testCounter")
	ms.AddCounter(10, "testCounter") // Увеличиваем значение

	// Проверяем обновление значения
	if value, exists := ms.Counters()["testCounter"]; !exists {
		t.Fatalf("Expected counter 'testCounter' to exist")
	} else if value != 15 {
		t.Errorf("Expected counter value 15, got %d", value)
	}
}

func TestEmptyStorage(t *testing.T) {
	ms := metric.NewMemStorage()

	if len(ms.Gauges()) != 0 {
		t.Errorf("Expected no gauges, found %d", len(ms.Gauges()))
	}
	if len(ms.Counters()) != 0 {
		t.Errorf("Expected no counters, found %d", len(ms.Counters()))
	}
}
