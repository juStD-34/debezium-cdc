package config

import (
	"os"
)

type Config struct {
	Port            string
	KafkaConnectURL string
	DatabaseURL     string
	LogLevel        string
}

func Load(configFile, port, kafkaURL string) *Config {
	cfg := &Config{
		Port:            getEnvOrDefault("PORT", port),
		KafkaConnectURL: getEnvOrDefault("KAFKA_CONNECT_URL", kafkaURL),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
