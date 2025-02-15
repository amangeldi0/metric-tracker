package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg, err := New()

	assert.NoError(t, err)

	assert.Equal(t, "localhost:8080", cfg.Address, "want localhost:8080 got", cfg.Address)
}
