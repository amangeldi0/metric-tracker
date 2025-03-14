package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name           string
		argHTTPMethod  string
		wantStatusCode int
	}{
		{
			name:           "Positive POST",
			argHTTPMethod:  http.MethodPost,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Positive DELETE",
			argHTTPMethod:  http.MethodDelete,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	r := setupRouter(nil, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.argHTTPMethod, "/", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		httpMethod string
		path       string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "Positive gauge",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/123.0",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Positive gauge (with big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/123.123456789123456789",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative gauge (without id)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative gauge (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge//123.45",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative gauge (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "gauge/test/123.0",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative gauge (invalid value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "gauge/test/invalid_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},

		{
			name: "Positive counter",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/123",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Negative counter (without id)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative counter (without params)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter//10",
			},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "Negative counter (invalid http method)",
			args: args{
				httpMethod: http.MethodGet,
				path:       "counter/test/123",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (big value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/9223372036854775808",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (text value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/text_value",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Negative counter (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				path:       "counter/test/12345.6789",
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.httpMethod, fmt.Sprintf("/update/%s", tt.args.path), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.wantStatusCode)
		})
	}
}

func getPointerFloat64(v float64) *float64 {
	return &v
}

func getPointerInt64(v int64) *int64 {
	return &v
}

func TestUpdateByBody(t *testing.T) {
	type args struct {
		httpMethod string
		body       any
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive gauge",
			args: args{
				httpMethod: http.MethodPost,
				body: models.MetricsUpdate{
					ID:    "test",
					MType: string(storage.GaugeType),
					Value: getPointerFloat64(123.0),
				},
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "{\"id\":\"test\",\"type\":\"gauge\",\"value\":123}",
			},
		},
		{
			name: "Positive gauge (with big value)",
			args: args{
				httpMethod: http.MethodPost,
				body: models.MetricsUpdate{
					ID:    "test",
					MType: string(storage.GaugeType),
					Value: getPointerFloat64(123.123456789123456789),
				},
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "{\"id\":\"test\",\"type\":\"gauge\",\"value\":123.12345678912345}",
			},
		},
		{
			name: "Negative gauge (invalid http method)",
			args: args{
				httpMethod: http.MethodDelete,
				body: models.MetricsUpdate{
					ID:    "test",
					MType: string(storage.GaugeType),
					Value: getPointerFloat64(123.0),
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "",
			},
		},
		{
			name: "Negative gauge (invalid value)",
			args: args{
				httpMethod: http.MethodPost,
				body: map[string]any{
					"id":    "test",
					"type":  "gauge",
					"value": "invalid_value",
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Field value \\\"value\\\" must be float64.\"}",
			},
		},
		{
			name: "Positive counter",
			args: args{
				httpMethod: http.MethodPost,
				body: models.MetricsUpdate{
					ID:    "Test1",
					MType: "counter",
					Delta: getPointerInt64(123),
				},
			},
			want: want{
				statusCode: http.StatusOK,
				body:       "{\"id\":\"Test1\",\"type\":\"counter\",\"delta\":123}",
			},
		},
		{
			name: "Negative counter (invalid http method)",
			args: args{
				httpMethod: http.MethodDelete,
				body: models.MetricsUpdate{
					ID:    "Test2",
					MType: "counter",
					Delta: getPointerInt64(123),
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "",
			},
		},
		{
			name: "Negative counter (text value)",
			args: args{
				httpMethod: http.MethodPost,
				body: map[string]any{
					"id":    "Test3",
					"type":  "counter",
					"delta": "invalid_value",
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			},
		},
		{
			name: "Negative counter (float64 value)",
			args: args{
				httpMethod: http.MethodPost,
				body: map[string]any{
					"id":    "Test3",
					"type":  "counter",
					"delta": 12345.6789,
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Field value \\\"delta\\\" must be int64.\"}",
			},
		},
		{
			name: "Negative (without params)",
			args: args{
				httpMethod: http.MethodPost,
				body:       models.MetricsUpdate{},
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
			},
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(&tt.args.body)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(tt.args.httpMethod, "/update/", bytes.NewReader(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), tt.want.body)
		})
	}
}

func TestValue(t *testing.T) {
	type args struct {
		name       string
		metricType storage.MetricType
		value      interface{}
	}
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive (gauge)",
			args: args{
				name:       "TestGauge",
				metricType: storage.GaugeType,
				value:      10.5,
			},
			want: want{
				body:       "10.5",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Positive (counter)",
			args: args{
				name:       "TestCounter",
				metricType: storage.CounterType,
				value:      int64(10500),
			},
			want: want{
				body:       "10500",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Negative (gauge)",
			args: args{
				metricType: storage.GaugeType,
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative (counter)",
			args: args{
				metricType: storage.CounterType,
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "Negative (invalid metric type)",
			args: args{
				metricType: storage.MetricType("invalid"),
			},
			want: want{
				body:       "",
				statusCode: http.StatusNotFound,
			},
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.name != "" {
				switch tt.args.metricType {
				case storage.CounterType:
					memStorage.AddCounter(tt.args.name, tt.args.value.(int64))
				case storage.GaugeType:
					memStorage.SetGauge(tt.args.name, tt.args.value.(float64))
				}
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/value/%s/%s/", tt.args.metricType, tt.args.name), nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), tt.want.body)
		})
	}
}

func TestValueByBody(t *testing.T) {
	type args struct {
		name       string
		metricType storage.MetricType
		value      interface{}
	}
	type want struct {
		body       string
		statusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive (gauge)",
			args: args{
				name:       "TestGauge",
				metricType: storage.GaugeType,
				value:      10.5,
			},
			want: want{
				body:       "{\"id\":\"TestGauge\",\"type\":\"gauge\",\"value\":10.5}",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Positive (counter)",
			args: args{
				name:       "TestCounter",
				metricType: storage.CounterType,
				value:      int64(10500),
			},
			want: want{
				body:       "{\"id\":\"TestCounter\",\"type\":\"counter\",\"delta\":10500}",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Negative (gauge)",
			args: args{
				metricType: storage.GaugeType,
			},
			want: want{
				body:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "Negative (counter)",
			args: args{
				metricType: storage.CounterType,
			},
			want: want{
				body:       "{\"error\":\"Field validation for \\\"ID\\\" failed on the 'required' tag.\"}",
				statusCode: http.StatusBadRequest,
			},
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.name != "" {
				switch tt.args.metricType {
				case storage.CounterType:
					memStorage.AddCounter(tt.args.name, tt.args.value.(int64))
				case storage.GaugeType:
					memStorage.SetGauge(tt.args.name, tt.args.value.(float64))
				}
			}

			jsonBytes, err := json.Marshal(&models.MetricsValue{
				ID:    tt.args.name,
				MType: string(tt.args.metricType),
			})
			require.NoError(t, err)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, "/value/", bytes.NewReader(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, res.StatusCode, tt.want.statusCode)
			assert.Equal(t, string(body), tt.want.body)
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		httpMethod string
	}
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive",
			args: args{
				httpMethod: http.MethodGet,
			},
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/html; charset=utf-8",
			},
		},
		{
			name: "Negative (POST)",
			args: args{
				httpMethod: http.MethodPost,
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	memStorage := storage.NewMem()
	r := setupRouter(&memStorage, nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.args.httpMethod, "/", nil)

			r.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
