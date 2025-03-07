package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuntimeMetrics(t *testing.T) {
	runtimeMetrics := NewRuntimeMetrics()

	assert.NoError(t, runtimeMetrics.Update())
}
