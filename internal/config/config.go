package config

import (
	"fmt"
	"os"
	"strconv"
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
	ClawMemBaseURL      string
	ClawMemTimeout      time.Duration
	ScenarioPacksDir    string
	ReplayArchiveDir    string
	ReportsDir          string
	BenchmarkScheduler  SchedulerConfig
}

type SchedulerConfig struct {
	Enabled        bool
	ScenarioFamily string
	Interval       time.Duration
	MaxRuns        int
	DryRun         bool
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:           envOrDefault("APP_ENV", "development"),
		ServiceAddress:   envOrDefault("SERVICE_ADDRESS", "127.0.0.1:8090"),
		LogLevel:         envOrDefault("LOG_LEVEL", "info"),
		ControlPlaneURL:  strings.TrimSpace(os.Getenv("CONTROL_PLANE_BASE_URL")),
		ClawMemBaseURL:   envOrDefaultCompat("CLAWMEM_BASE_URL", "MEMORY_BASE_URL", "http://127.0.0.1:8088"),
		ScenarioPacksDir: envOrDefault("SCENARIO_PACKS_DIR", "./configs/scenario-packs"),
		ReplayArchiveDir: envOrDefault("REPLAY_ARCHIVE_DIR", "./var/replay-archive"),
		ReportsDir:       envOrDefault("REPORTS_DIR", "./reports"),
		BenchmarkScheduler: SchedulerConfig{
			Enabled:        strings.EqualFold(envOrDefault("BENCHMARK_SCHEDULER_ENABLED", "false"), "true"),
			ScenarioFamily: envOrDefault("BENCHMARK_SCHEDULER_SCENARIO_FAMILY", "commerce"),
			MaxRuns:        intEnvOrDefault("BENCHMARK_SCHEDULER_MAX_RUNS", 7),
			DryRun:         strings.EqualFold(envOrDefault("BENCHMARK_SCHEDULER_DRY_RUN", "false"), "true"),
		},
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
	cfg.ClawMemTimeout, err = time.ParseDuration(envOrDefaultCompat("CLAWMEM_TIMEOUT", "MEMORY_TIMEOUT", "5s"))
	if err != nil {
		return Config{}, fmt.Errorf("parse CLAWMEM_TIMEOUT: %w", err)
	}
	cfg.BenchmarkScheduler.Interval, err = time.ParseDuration(envOrDefault("BENCHMARK_SCHEDULER_INTERVAL", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("parse BENCHMARK_SCHEDULER_INTERVAL: %w", err)
	}

	if cfg.ControlPlaneURL == "" {
		return Config{}, fmt.Errorf("CONTROL_PLANE_BASE_URL is required")
	}
	if cfg.ClawMemBaseURL == "" {
		return Config{}, fmt.Errorf("CLAWMEM_BASE_URL is required")
	}
	if cfg.ServiceAddress == "" {
		return Config{}, fmt.Errorf("SERVICE_ADDRESS is required")
	}
	if cfg.ScenarioPacksDir == "" {
		return Config{}, fmt.Errorf("SCENARIO_PACKS_DIR is required")
	}
	if cfg.ReplayArchiveDir == "" {
		return Config{}, fmt.Errorf("REPLAY_ARCHIVE_DIR is required")
	}
	if cfg.ReportsDir == "" {
		return Config{}, fmt.Errorf("REPORTS_DIR is required")
	}
	if cfg.BenchmarkScheduler.MaxRuns < 0 {
		return Config{}, fmt.Errorf("BENCHMARK_SCHEDULER_MAX_RUNS must be >= 0")
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func envOrDefaultCompat(primary string, legacy string, fallback string) string {
	if value, ok := os.LookupEnv(primary); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	if value, ok := os.LookupEnv(legacy); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func intEnvOrDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
