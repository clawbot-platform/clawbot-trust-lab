package app

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/http/handlers"
	httpmw "clawbot-trust-lab/internal/http/middleware"
	"clawbot-trust-lab/internal/http/routes"
	"clawbot-trust-lab/internal/platform/bootstrap"
	"clawbot-trust-lab/internal/version"
)

func NewLogger(level string, writer io.Writer) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slogLevel}))
}

func Run(ctx context.Context, cfg config.Config, logger *slog.Logger) error {
	deps, err := bootstrap.Build(cfg)
	if err != nil {
		return err
	}

	system := handlers.NewSystemHandler(func(ctx context.Context) error {
		return bootstrap.Ready(ctx, deps)
	}, version.Current())
	trustLab := handlers.NewTrustLabHandler(deps.Scenarios, deps.Execution, deps.Trust, deps.Replay, deps.Benchmark, deps.Commerce, deps.Events, deps.TrustFlow, deps.Detection, handlers.TrustLabState{
		AppEnv:          cfg.AppEnv,
		ControlPlaneURL: cfg.ControlPlaneURL,
		ClawMemBaseURL:  cfg.ClawMemBaseURL,
	})

	server := &http.Server{
		Addr: cfg.ServiceAddress,
		Handler: routes.New(httpmw.RequestLogger(logger), routes.Services{
			System:   system,
			TrustLab: trustLab,
		}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
