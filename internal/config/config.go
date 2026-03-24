package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	AppEnv              string
	ServiceAddress      string
	LogLevel            string
	ShutdownTimeout     time.Duration
	ControlPlaneURL     string
	ControlPlaneTimeout time.Duration
	MemoryURL           string
	MemoryTimeout       time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:          envOrDefault("APP_ENV", "development"),
		ServiceAddress:  envOrDefault("SERVICE_ADDRESS", "127.0.0.1:8090"),
		LogLevel:        envOrDefault("LOG_LEVEL", "info"),
		ControlPlaneURL: strings.TrimSpace(os.Getenv("CONTROL_PLANE_BASE_URL")),
		MemoryURL:       strings.TrimSpace(os.Getenv("MEMORY_BASE_URL")),
	}

	var err error
	cfg.ShutdownTimeout, err = time.ParseDuration(envOrDefault("SHUTDOWN_TIMEOUT", "10s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.ControlPlaneTimeout, err = time.ParseDuration(envOrDefault("CONTROL_PLANE_TIMEOUT", "5s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse CONTROL_PLANE_TIMEOUT: %w", err)
	}
	cfg.MemoryTimeout, err = time.ParseDuration(envOrDefault("MEMORY_TIMEOUT", "5s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse MEMORY_TIMEOUT: %w", err)
	}

	if cfg.ControlPlaneURL == "" {
		return Config{}, fmt.Errorf("CONTROL_PLANE_BASE_URL is required")
	}
	if cfg.MemoryURL == "" {
		return Config{}, fmt.Errorf("MEMORY_BASE_URL is required")
	}
	if cfg.ServiceAddress == "" {
		return Config{}, fmt.Errorf("SERVICE_ADDRESS is required")
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}
