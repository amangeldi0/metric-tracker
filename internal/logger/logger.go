package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type AppType string

const (
	Server AppType = "server"
	Agent  AppType = "agent"
)

func New(level zapcore.Level, app AppType) (*zap.Logger, error) {

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout", "logfile"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()

	if err != nil {
		return &zap.Logger{}, err
	}

	if app != Server && app != Agent {
		return logger, nil
	}

	logger = logger.With(
		zap.String("agent", string(app)),
	)

	defer logger.Sync()

	return logger, nil
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(p)
	lrw.responseData.size += size
	return size, err
}

// WriteHeader перехватывает статус ответа
func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.responseData.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{}
			lrw := &loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)
			logger.Infow("HTTP request",
				"uri", r.RequestURI,
				"method", r.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}
