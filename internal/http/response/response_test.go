package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJson(t *testing.T) {
	rr := httptest.NewRecorder()
	data := envelope{"message": "hello"}

	WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	expectedBody := `{"message":"hello"}`
	body := strings.TrimSpace(rr.Body.String())

	if body != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, body)
	}
}

func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ErrorResponse(rr, req, http.StatusBadRequest, "invalid request")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	expectedBody := `{"error":"invalid request"}`
	body := strings.TrimSpace(rr.Body.String())

	if body != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, body)
	}
}

func TestNotFoundResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	NotFoundResponse(rr, nil)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}

	var resp envelope
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expectedMessage := "the required resource could not be found"
	if resp["error"] != expectedMessage {
		t.Errorf("expected error message %s, got %s", expectedMessage, resp["error"])
	}
}
