package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := New()

	if cfg.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Server.Port)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("expected host 'localhost', got %s", cfg.Server.Host)
	}

	if cfg.Server.Protocol != "http" {
		t.Errorf("expected protocol 'http', got %s", cfg.Server.Protocol)
	}
}
