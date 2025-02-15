package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger, err := New(zapcore.InfoLevel, Server)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	logger, err = New(zapcore.DebugLevel, Agent)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestWithLogging(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	loggedHandler := WithLogging(sugar)(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.RemoteAddr = "127.0.0.1:1234"

	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "test response", rr.Body.String())
}

func TestLoggingResponseWriter(t *testing.T) {
	rr := httptest.NewRecorder()
	responseData := &responseData{}

	lrw := &loggingResponseWriter{
		ResponseWriter: rr,
		responseData:   responseData,
	}

	lrw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, responseData.status)

	data := []byte("Hello, World!")
	n, err := lrw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, len(data), responseData.size)
}
