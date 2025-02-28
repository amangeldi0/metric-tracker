package commpress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Decompresser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = gz
			defer gz.Close()
		}

		acceptedEncodings := r.Header.Get("Accept-Encoding")
		canCompressResponse := strings.Contains(acceptedEncodings, "gzip")

		if canCompressResponse {
			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()
			w = gzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
			w.Header().Set("Content-Encoding", "gzip")
		}

		next.ServeHTTP(w, r)
	})
}
