package handlers

import (
	"net/http"

	"clawbot-trust-lab/internal/domain/scenario"
)

type ScenarioCatalog interface {
	ScenarioTypes() []scenario.ScenarioType
}

type TrustLabState struct {
	AppEnv          string
	ControlPlaneURL string
	MemoryURL       string
}

type TrustLabHandler struct {
	scenarios ScenarioCatalog
	state     TrustLabState
}

func NewTrustLabHandler(scenarios ScenarioCatalog, state TrustLabState) *TrustLabHandler {
	return &TrustLabHandler{scenarios: scenarios, state: state}
}

func (h *TrustLabHandler) ScenarioTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"types": h.scenarios.ScenarioTypes(),
		},
	})
}

func (h *TrustLabHandler) ReplayStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"status":        "scaffolded",
			"memory_client": "defined",
			"phase":         "Phase 2",
		},
	})
}

func (h *TrustLabHandler) TrustStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"status":            "scaffolded",
			"control_plane_url": h.state.ControlPlaneURL,
			"memory_contract":   "defined",
			"artifact_families": []string{"trust_artifact", "mandate_artifact", "provenance_artifact"},
		},
	})
}

func (h *TrustLabHandler) BenchmarkStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"status":                   "scaffolded",
			"benchmark_registration":   "control-plane client ready",
			"future_execution_runtime": "later phase",
		},
	})
}
