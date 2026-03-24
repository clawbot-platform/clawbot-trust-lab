package routes

import (
	"net/http"

	"clawbot-trust-lab/internal/http/handlers"
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

	mux.HandleFunc("/api/v1/scenarios", services.TrustLab.ListScenarios)
	mux.HandleFunc("/api/v1/scenarios/execute", services.TrustLab.ExecuteScenario)
	mux.HandleFunc("/api/v1/scenarios/types", services.TrustLab.ScenarioTypes)
	mux.HandleFunc("/api/v1/scenarios/packs", services.TrustLab.ListPacks)
	mux.HandleFunc("/api/v1/scenarios/packs/{id}", services.TrustLab.GetPack)
	mux.HandleFunc("/api/v1/detection/evaluate", services.TrustLab.EvaluateDetection)
	mux.HandleFunc("/api/v1/detection/results", services.TrustLab.ListDetectionResults)
	mux.HandleFunc("/api/v1/detection/results/{id}", services.TrustLab.GetDetectionResult)
	mux.HandleFunc("/api/v1/detection/rules", services.TrustLab.ListDetectionRules)
	mux.HandleFunc("/api/v1/detection/summary", services.TrustLab.DetectionSummary)
	mux.HandleFunc("/api/v1/orders", services.TrustLab.ListOrders)
	mux.HandleFunc("/api/v1/orders/{id}", services.TrustLab.GetOrder)
	mux.HandleFunc("/api/v1/events", services.TrustLab.ListEvents)
	mux.HandleFunc("/api/v1/replay/status", services.TrustLab.ReplayStatus)
	mux.HandleFunc("/api/v1/replay/cases", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			services.TrustLab.ListReplayCases(w, r)
		case http.MethodPost:
			services.TrustLab.CreateReplayCase(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/v1/trust/status", services.TrustLab.TrustStatus)
	mux.HandleFunc("/api/v1/trust/artifacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			services.TrustLab.ListArtifacts(w, r)
		case http.MethodPost:
			services.TrustLab.CreateArtifact(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/v1/benchmark/status", services.TrustLab.BenchmarkStatus)
	mux.HandleFunc("/api/v1/trust/decisions", services.TrustLab.ListTrustDecisions)
	mux.HandleFunc("/api/v1/trust/decisions/{id}", services.TrustLab.GetTrustDecision)
	mux.HandleFunc("/api/v1/benchmark/rounds/register", services.TrustLab.RegisterBenchmarkRound)
	mux.HandleFunc("/api/v1/benchmark/rounds/status", services.TrustLab.BenchmarkRoundStatus)

	return loggerMiddleware(mux)
}
