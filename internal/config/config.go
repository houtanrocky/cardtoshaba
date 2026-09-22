package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort          string
	ZarinHubBaseURL     string
	ZarinHubAPIUsername string
	ZarinHubAPIPassword string
	ZarinHubAppName     string
	ZarinHubTimeout     time.Duration
	DatabaseURL         string
	LogLevel            string
	Environment         string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "no .env file found, using system environment variables\n")
	}

	cfg := &Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		ZarinHubBaseURL:     getEnv("ZARINHUB_BASE_URL", ""),
		ZarinHubAPIUsername: getEnv("ZARINHUB_API_USERNAME", ""),
		ZarinHubAPIPassword: getEnv("ZARINHUB_API_PASSWORD", ""),
		ZarinHubAppName:     getEnv("ZARINHUB_APPLICATION_NAME", ""),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		Environment:         getEnv("ENVIRONMENT", "development"),
	}

	// Timeout
	timeoutStr := getEnv("ZARINHUB_TIMEOUT", "15s")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid ZARINHUB_TIMEOUT: %w", err)
	}
	cfg.ZarinHubTimeout = timeout

	// Validating required envs
	if cfg.ZarinHubBaseURL == "" {
		return nil, fmt.Errorf("ZARINHUB_BASE_URL is required")
	}
	if cfg.ZarinHubAPIUsername == "" {
		return nil, fmt.Errorf("ZARINHUB_API_USERNAME is required")
	}
	if cfg.ZarinHubAPIPassword == "" {
		return nil, fmt.Errorf("ZARINHUB_API_PASSWORD is required")
	}
	if cfg.ZarinHubAppName == "" {
		return nil, fmt.Errorf("ZARINHUB_APPLICATION_NAME is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
