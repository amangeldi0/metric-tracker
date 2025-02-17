package main

import (
	"io"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	go func() {
		err := run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/")
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: got %d, want 200 or 404", resp.StatusCode)
	}
}
