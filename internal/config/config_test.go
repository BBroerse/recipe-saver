package config_test

import (
	"recipe-processor/internal/config"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("REQUEST_TIMEOUT", "")
	t.Setenv("READ_TIMEOUT", "")
	t.Setenv("WRITE_TIMEOUT", "")
	t.Setenv("IDLE_TIMEOUT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("OLLAMA_BASE_URL", "")
	t.Setenv("NOTION_TOKEN", "")
	t.Setenv("NOTION_DATABASE_ID", "")

	cfg := config.Load()

	if cfg.Environment != "development" {
		t.Errorf("expected default Environment=development, got %s", cfg.Environment)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected default Port=8080, got %s", cfg.Port)
	}
	if cfg.RequestTimeout != 30*time.Minute {
		t.Errorf("expected default RequestTimeout=30m, got %v", cfg.RequestTimeout)
	}
	if cfg.ReadTimeout != 10*time.Minute {
		t.Errorf("expected default ReadTimeout=10m, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Minute {
		t.Errorf("expected default WriteTimeout=10m, got %v", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 60*time.Minute {
		t.Errorf("expected default IdleTimeout=60m, got %v", cfg.IdleTimeout)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default LogLevel=info, got %s", cfg.LogLevel)
	}
	if cfg.OllamaBaseUrl != "http://ollama:11434" {
		t.Errorf("expected default OllamaBaseUrl=http://ollama:11434, got %s", cfg.OllamaBaseUrl)
	}
	if cfg.NotionToken != "" {
		t.Errorf("expected default NotionToken empty, got %s", cfg.NotionToken)
	}
	if cfg.NotionDatabaseId != "" {
		t.Errorf("expected default NotionDatabaseId empty, got %s", cfg.NotionDatabaseId)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("PORT", "9000")
	t.Setenv("REQUEST_TIMEOUT", "120")
	t.Setenv("READ_TIMEOUT", "30")
	t.Setenv("WRITE_TIMEOUT", "45")
	t.Setenv("IDLE_TIMEOUT", "300")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("OLLAMA_BASE_URL", "http://localhost:1234")
	t.Setenv("NOTION_TOKEN", "xyz")
	t.Setenv("NOTION_DATABASE_ID", "abc")

	cfg := config.Load()

	if cfg.Environment != "production" {
		t.Errorf("expected ENV=production, got %s", cfg.Environment)
	}
	if cfg.Port != "9000" {
		t.Errorf("expected PORT=9000, got %s", cfg.Port)
	}
	if cfg.RequestTimeout != 120*time.Second {
		t.Errorf("expected RequestTimeout=120s, got %v", cfg.RequestTimeout)
	}
	if cfg.ReadTimeout != 30*time.Second {
		t.Errorf("expected ReadTimeout=30s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 45*time.Second {
		t.Errorf("expected WriteTimeout=45s, got %v", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 300*time.Second {
		t.Errorf("expected IdleTimeout=300s, got %v", cfg.IdleTimeout)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %s", cfg.LogLevel)
	}
	if cfg.OllamaBaseUrl != "http://localhost:1234" {
		t.Errorf("expected OllamaBaseUrl override, got %s", cfg.OllamaBaseUrl)
	}
	if cfg.NotionToken != "xyz" {
		t.Errorf("expected NotionToken=xyz, got %s", cfg.NotionToken)
	}
	if cfg.NotionDatabaseId != "abc" {
		t.Errorf("expected NotionDatabaseId=abc, got %s", cfg.NotionDatabaseId)
	}
}

func TestLoad_InvalidDurationFallback(t *testing.T) {
	t.Setenv("REQUEST_TIMEOUT", "not-a-number")
	t.Setenv("READ_TIMEOUT", "oops")
	t.Setenv("WRITE_TIMEOUT", "-12s")
	t.Setenv("IDLE_TIMEOUT", "123abc")

	cfg := config.Load()

	if cfg.RequestTimeout != 30*time.Minute {
		t.Errorf("expected default RequestTimeout=30m on invalid input, got %v", cfg.RequestTimeout)
	}
	if cfg.ReadTimeout != 10*time.Minute {
		t.Errorf("expected default ReadTimeout=10m on invalid input, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 10*time.Minute {
		t.Errorf("expected default WriteTimeout=10m on invalid input, got %v", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 60*time.Minute {
		t.Errorf("expected default IdleTimeout=60m on invalid input, got %v", cfg.IdleTimeout)
	}
}
