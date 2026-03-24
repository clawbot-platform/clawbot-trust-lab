package routes

import (
	"net/http"

	"clawbot-trust-lab/internal/http/handlers"
	"clawbot-trust-lab/internal/http/middleware"
)

type Services struct {
	System   *handlers.SystemHandler
	TrustLab *handlers.TrustLabHandler
}

func New(loggerMiddleware func(http.Handler) http.Handler, services Services) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", services.System.Health)
	mux.HandleFunc("/readyz", services.System.Ready)
	mux.HandleFunc("/version", services.System.Version)

	mux.HandleFunc("/api/v1/scenarios/types", services.TrustLab.ScenarioTypes)
	mux.HandleFunc("/api/v1/replay/status", services.TrustLab.ReplayStatus)
	mux.HandleFunc("/api/v1/trust/status", services.TrustLab.TrustStatus)
	mux.HandleFunc("/api/v1/benchmark/status", services.TrustLab.BenchmarkStatus)

	return loggerMiddleware(mux)
}

func LoggerChain(loggerMiddleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return loggerMiddleware
}

var _ = middleware.RequestLogger
