package metric

import "testing"

func TestMemoryStorageUpdate(t *testing.T) {
	ms := NewMemStorage()
	key := "memoryusage"

	// Обновляем значение gauge
	ms.UpdateGauge(4.5, key)

	// Проверяем, что ключ существует
	value, exists := (*ms)[key]
	if !exists {
		t.Errorf("Key %s does not exist in MemStorage", key)
	}

	// Проверяем, что значение gauge обновлено
	if value.Gauge != 4.5 {
		t.Errorf("Expected gauge value 4.5, got %v", value.Gauge)
	}

	// Перезаписываем значение gauge
	ms.UpdateGauge(7.8, key)
	value = (*ms)[key]
	if value.Gauge != 7.8 {
		t.Errorf("Expected gauge value 7.8 after update, got %v", value.Gauge)
	}
}

func TestMemoryStorageAdd(t *testing.T) {
	ms := NewMemStorage()
	key := "pollcount"

	ms.AddCounter(1, key)

	value, exists := (*ms)[key]
	if !exists {
		t.Errorf("Key %s does not exist in MemStorage", key)
	}

	if value.Counter != 1 {
		t.Errorf("Expected counter value 1, got %v", value.Counter)
	}

	ms.AddCounter(1, key)
	value = (*ms)[key]
	if value.Counter != 2 {
		t.Errorf("Expected counter value 2, got %v", value.Counter)
	}

	ms.AddCounter(5, key)
	value = (*ms)[key]
	if value.Counter != 7 {
		t.Errorf("Expected counter value 7, got %v", value.Counter)
	}
}
