package commpress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecompresser(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.Write(body)
	})

	wrappedHandler := Decompresser(handler)

	var compressedBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedBuffer)
	_, err := gzipWriter.Write([]byte("test data"))
	assert.NoError(t, err)
	assert.NoError(t, gzipWriter.Close())

	req := httptest.NewRequest(http.MethodPost, "http://example.com", &compressedBuffer)
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(responseBody))
}

func TestCompressorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("compressed response"))
	})

	wrappedHandler := Decompresser(handler)

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))

	gzReader, err := gzip.NewReader(res.Body)
	assert.NoError(t, err)

	decompressedBody, err := io.ReadAll(gzReader)
	assert.NoError(t, err)
	assert.Equal(t, "compressed response", string(decompressedBody))
}
