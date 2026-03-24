package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"clawbot-trust-lab/internal/app"
	"clawbot-trust-lab/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	logger := app.NewLogger(cfg.LogLevel, os.Stdout)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, cfg, logger); err != nil {
		fmt.Fprintf(os.Stderr, "trust-lab: %v\n", err)
		os.Exit(1)
	}
}
