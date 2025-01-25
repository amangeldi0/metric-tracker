package metric

import (
	"testing"
)

func TestUpdateGauge(t *testing.T) {
	ms := NewMemStorage()

	ms.UpdateGauge(10.5, "testGauge")
	ms.UpdateGauge(20.2, "testGauge")

	if value, exists := ms.gauges["testGauge"]; !exists {
		t.Fatalf("Expected gauge 'testGauge' to exist")
	} else if value != 20.2 {
		t.Errorf("Expected gauge value 20.2, got %f", value)
	}
}

func TestAddCounter(t *testing.T) {
	ms := NewMemStorage()

	ms.AddCounter(5, "testCounter")
	ms.AddCounter(10, "testCounter")

	if value, exists := ms.counters["testCounter"]; !exists {
		t.Fatalf("Expected counter 'testCounter' to exist")
	} else if value != 15 {
		t.Errorf("Expected counter value 15, got %d", value)
	}
}

func TestEmptyStorage(t *testing.T) {
	ms := NewMemStorage()

	if len(ms.gauges) != 0 {
		t.Errorf("Expected no gauges, found %d", len(ms.gauges))
	}
	if len(ms.counters) != 0 {
		t.Errorf("Expected no counters, found %d", len(ms.counters))
	}
}
