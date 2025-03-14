package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/amangeldi0/metric-tracker/internal/server/models"
	"github.com/amangeldi0/metric-tracker/internal/server/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateContentType(t *testing.T) {
	type args struct {
		contentTypeRequest string
		contentTypeNeed    string
		withoutContentType bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Positive text/plain",
			args: args{
				contentTypeRequest: "text/plain",
				contentTypeNeed:    "text/plain",
				withoutContentType: true,
			},
			want: true,
		},
		{
			name: "Positive json/application",
			args: args{
				contentTypeRequest: "json/application",
				contentTypeNeed:    "json/application",
				withoutContentType: true,
			},
			want: true,
		},
		{
			name: "Positive without request content-type",
			args: args{
				contentTypeRequest: "",
				contentTypeNeed:    "text/plain",
				withoutContentType: true,
			},
			want: true,
		},
		{
			name: "Negative without request content-type",
			args: args{
				contentTypeRequest: "",
				contentTypeNeed:    "json/application",
				withoutContentType: false,
			},
			want: false,
		},
		{
			name: "Negative json/application & text/plain",
			args: args{
				contentTypeRequest: "json/application",
				contentTypeNeed:    "text/plain",
				withoutContentType: true,
			},
			want: false,
		},
	}

	gin.SetMode(gin.ReleaseMode)
	bh := NewBase(nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.args.contentTypeRequest != "" {
				ctx.Request.Header.Set("Content-Type", tt.args.contentTypeRequest)
			}

			assert.Equal(t, tt.want, bh.validateContentType(ctx, tt.args.contentTypeNeed, tt.args.withoutContentType))
		})
	}
}

func getRandomFloat64() *float64 {
	v := rand.Float64()
	return &v
}

func getRandomInt64() *int64 {
	v := rand.Int63()
	return &v
}

func TestValidateAndShouldBindJSON(t *testing.T) {
	type args struct {
		obj         models.MetricsUpdate
		body        models.MetricsUpdate
		withoutBody bool
	}
	type want struct {
		errorResponse *models.ErrorResponse
		statusCode    int
		wantErr       bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Positive gauge",
			args: args{
				obj: models.MetricsUpdate{},
				body: models.MetricsUpdate{
					ID:    "TestGauge",
					MType: string(storage.GaugeType),
					Value: getRandomFloat64(),
				},
			},
			want: want{
				errorResponse: nil,
				statusCode:    0,
				wantErr:       false,
			},
		},
		{
			name: "Positive counter",
			args: args{
				obj: models.MetricsUpdate{},
				body: models.MetricsUpdate{
					ID:    "TestCounter",
					MType: string(storage.CounterType),
					Delta: getRandomInt64(),
				},
			},
			want: want{
				errorResponse: nil,
				statusCode:    0,
				wantErr:       false,
			},
		},
		{
			name: "Negative without params",
			args: args{
				obj:  models.MetricsUpdate{},
				body: models.MetricsUpdate{},
			},
			want: want{
				errorResponse: &models.ErrorResponse{Error: "Field validation for \"ID\" failed on the 'required' tag."},
				statusCode:    http.StatusBadRequest,
				wantErr:       true,
			},
		},
		{
			name: "Negative invalid type",
			args: args{
				obj: models.MetricsUpdate{},
				body: models.MetricsUpdate{
					ID:    "NTest",
					MType: "heh",
				},
			},
			want: want{
				errorResponse: &models.ErrorResponse{Error: "Field validation for \"MType\" failed on the 'oneof=counter gauge' tag."},
				statusCode:    http.StatusBadRequest,
				wantErr:       true,
			},
		},
		{
			name: "Negative without value (gauge)",
			args: args{
				obj: models.MetricsUpdate{},
				body: models.MetricsUpdate{
					ID:    "NTest",
					MType: string(storage.GaugeType),
					Delta: getRandomInt64(),
				},
			},
			want: want{
				errorResponse: &models.ErrorResponse{Error: "Field validation for \"Value\" failed on the 'required_if=MType gauge' tag."},
				statusCode:    http.StatusBadRequest,
				wantErr:       true,
			},
		},
		{
			name: "Negative without value (counter)",
			args: args{
				obj: models.MetricsUpdate{},
				body: models.MetricsUpdate{
					ID:    "NTest",
					MType: string(storage.CounterType),
					Value: getRandomFloat64(),
				},
			},
			want: want{
				errorResponse: &models.ErrorResponse{Error: "Field validation for \"Delta\" failed on the 'required_if=MType counter' tag."},
				statusCode:    http.StatusBadRequest,
				wantErr:       true,
			},
		},
		{
			name: "Negative without request body",
			args: args{
				withoutBody: true,
			},
			want: want{
				errorResponse: &models.ErrorResponse{Error: "Request body not provided."},
				statusCode:    http.StatusBadRequest,
				wantErr:       true,
			},
		},
	}

	gin.SetMode(gin.ReleaseMode)
	bh := NewBase(nil, zaptest.NewLogger(t).Sugar(), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body io.Reader
			if !tt.args.withoutBody {
				jsonBytes, err := json.Marshal(&tt.args.body)
				require.NoError(t, err)

				body = bytes.NewReader(jsonBytes)
			}

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/", body)

			response, statusCode, err := bh.validateAndShouldBindJSON(ctx, &tt.args.obj)

			assert.Equal(t, response, tt.want.errorResponse)
			assert.Equal(t, statusCode, tt.want.statusCode)
			if tt.want.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.obj, tt.args.body)
			}
		})
	}
}
