package config

import (
	"os"
)

type Config struct {
	Port         string
	ConnectorUrl string
	DatabaseURL  string
	LogLevel     string
}

func Load() *Config {
	cfg := &Config{
		Port:         getEnvOrDefault("PORT", "8080"),
		ConnectorUrl: getEnvOrDefault("KAFKA_CONNECT_URL", "http://192.168.49.2:31994"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
