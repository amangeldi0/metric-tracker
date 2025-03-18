package storage

import (
	"bytes"
	"encoding/json"
	"github.com/amangeldi0/metric-tracker/internal/server/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"os"
	"testing"
)

func getPointerFloat64(v float64) *float64 {
	return &v
}

func getPointerInt64(v int64) *int64 {
	return &v
}

func TestSuccessFileStorage(t *testing.T) {
	metrics := []MetricsValue{
		{
			ID:    "TestGauge",
			MType: string(GaugeType),
			Value: getPointerFloat64(100.5),
		},
		{
			ID:    "TestCounter",
			MType: string(CounterType),
			Delta: getPointerInt64(321),
		},
		{
			ID:    "TestFloat64SimilarInt64",
			MType: string(GaugeType),
			Value: getPointerFloat64(300.0),
		},
	}

	file, err := os.CreateTemp(t.TempDir(), "tests-file-storage-*.json")
	require.NoError(t, err)
	t.Setenv("FILE_STORAGE_PATH", file.Name())

	require.NoError(t, config.Parse())

	stor := NewMem()
	log := zaptest.NewLogger(t).Sugar()

	fStorage, err := NewFileStorage(log)
	require.NoError(t, err)

	for _, metric := range metrics {
		switch metric.MType {
		case string(GaugeType):
			stor.SetGauge(metric.ID, *metric.Value)
		case string(CounterType):
			stor.AddCounter(metric.ID, *metric.Delta)
		}
	}

	count, err := fStorage.update()

	require.NoError(t, err)
	require.Equal(t, count, len(metrics))

	stor = NewMem()

	fStorage, err = NewFileStorage(log)
	require.NoError(t, err)

	require.NoError(t, fStorage.Restore())

	all, _ := stor.GetAll()
	require.Equal(t, len(all), len(metrics))

	if err = os.Remove(file.Name()); err != nil {
		t.Logf("Не удалось удалить тестовый json-файл: %s", err)
	}
}

func TestNegativeFileStorage(t *testing.T) {
	metrics := []MetricsValue{
		{
			ID:    "TestGauge",
			MType: "InvalidType",
			Value: getPointerFloat64(100.5),
		},
	}

	file, err := os.CreateTemp(t.TempDir(), "tests-file-storage-*.json")
	require.NoError(t, err)
	t.Setenv("FILE_STORAGE_PATH", file.Name())

	jsonBytes, err := json.Marshal(&metrics)
	require.NoError(t, err)

	_, err = file.Write(jsonBytes)
	require.NoError(t, err)

	require.NoError(t, config.Parse())

	buf := new(bytes.Buffer)
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buf),
		zapcore.DebugLevel),
		zap.AddCaller(),
	)

	fStorage, err := NewFileStorage(log.Sugar())
	require.NoError(t, err)

	require.NoError(t, fStorage.Restore())
	require.Contains(t, buf.String(), "The metric couldn't be restored, it has an unknown type")

	if err = os.Remove(file.Name()); err != nil {
		t.Logf("Не удалось удалить тестовый json-файл: %s", err)
	}
}
