package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Environment
	Environment string
	Port        string

	// Server timeouts
	RequestTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration

	// Logging
	LogLevel string

	// LLM
	OllamaBaseUrl string

	// Notion tokens
	NotionToken      string
	NotionDatabaseId string
}

func Load() *Config {
	return &Config{
		Environment:      getEnv("ENV", "development"),
		Port:             getEnv("PORT", "8080"),
		RequestTimeout:   getDurationEnv("REQUEST_TIMEOUT", 30*time.Minute),
		ReadTimeout:      getDurationEnv("READ_TIMEOUT", 10*time.Minute),
		WriteTimeout:     getDurationEnv("WRITE_TIMEOUT", 10*time.Minute),
		IdleTimeout:      getDurationEnv("IDLE_TIMEOUT", 60*time.Minute),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		OllamaBaseUrl:    getEnv("OLLAMA_BASE_URL", "http://ollama:11434"),
		NotionToken:      getEnv("NOTION_TOKEN", ""),
		NotionDatabaseId: getEnv("NOTION_DATABASE_ID", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}

	return defaultValue
}
