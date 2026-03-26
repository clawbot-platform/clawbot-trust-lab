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

	args := os.Args[1:]
	if len(args) > 0 && args[0] == "report" {
		if err := app.RunReportCommand(ctx, cfg, logger, os.Stdout, args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "trust-lab report: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := app.Run(ctx, cfg, logger); err != nil {
		fmt.Fprintf(os.Stderr, "trust-lab: %v\n", err)
		os.Exit(1)
	}
}
