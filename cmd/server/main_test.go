package main

import (
	"fmt"
	"github.com/amangeldi0/metric-tracker/cmd/server/metricsapi"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMainServer(t *testing.T) {
	// Создаем новую конфигурацию
	ms := metricsapi.New()

	// Создаем HTTP-сервер с тестовым маршрутом
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", ms.UpdateHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Проверяем, что сервер доступен по адресу
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/update/", server.URL))

	if err != nil {
		t.Fatalf("Сервер не запустился: %s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Ожидался код %d, но получили %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}
