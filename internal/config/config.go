package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	OllamaBaseUrl    string
	NotionToken      string
	NotionDatabaseId string
	RequestTimeout   time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		OllamaBaseUrl:    getEnv("OLLAMA_BASE_URL", "http://localhost:11434"),
		NotionToken:      getEnv("NOTION_TOKEN", ""),
		NotionDatabaseId: getEnv("NOTION_DATABASE_ID", ""),
		RequestTimeout:   getDurationEnv("REQUEST_TIMEOUT", 30*time.Minute),
		ReadTimeout:      getDurationEnv("READ_TIMEOUT", 10*time.Minute),
		WriteTimeout:     getDurationEnv("WRITE_TIMEOUT", 10*time.Minute),
		IdleTimeout:      getDurationEnv("IDLE_TIMEOUT", 60*time.Minute),
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
