package handlers

import (
	"context"
	"net/http"

	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
)

type ScenarioService interface {
	ListPacks() []scenario.ScenarioPack
	GetPack(string) (scenario.ScenarioPack, error)
}

type TrustService interface {
	CreateArtifact(context.Context, trust.CreateArtifactInput) (trust.TrustArtifact, error)
	ListArtifacts() []trust.TrustArtifact
}

type ReplayService interface {
	CreateCase(replay.CreateCaseInput) (replay.ReplayCase, error)
	ListCases() []replay.ReplayCase
}

type BenchmarkService interface {
	RegisterRound(context.Context, benchmark.RegistrationRequest) (benchmark.RegistrationResult, error)
	Status() map[string]any
}

type TrustLabState struct {
	AppEnv          string
	ControlPlaneURL string
	MemoryURL       string
}

type TrustLabHandler struct {
	scenarios ScenarioService
	trust     TrustService
	replay    ReplayService
	benchmark BenchmarkService
	state     TrustLabState
}

func NewTrustLabHandler(scenarios ScenarioService, trustService TrustService, replayService ReplayService, benchmarkService BenchmarkService, state TrustLabState) *TrustLabHandler {
	return &TrustLabHandler{
		scenarios: scenarios,
		trust:     trustService,
		replay:    replayService,
		benchmark: benchmarkService,
		state:     state,
	}
}

func (h *TrustLabHandler) ScenarioTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"types": scenario.KnownTypes(),
		},
	})
}

func (h *TrustLabHandler) ListPacks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.scenarios.ListPacks()})
}

func (h *TrustLabHandler) GetPack(w http.ResponseWriter, r *http.Request) {
	pack, err := h.scenarios.GetPack(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": pack})
}

func (h *TrustLabHandler) ReplayStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"status":        "active",
			"memory_client": "defined",
			"replay_cases":  len(h.replay.ListCases()),
			"phase":         "Phase 3",
		},
	})
}

func (h *TrustLabHandler) TrustStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"status":            "active",
			"control_plane_url": h.state.ControlPlaneURL,
			"memory_contract":   "defined",
			"artifact_count":    len(h.trust.ListArtifacts()),
			"artifact_families": []string{"trust_artifact", "mandate_artifact", "provenance_artifact"},
		},
	})
}

func (h *TrustLabHandler) BenchmarkStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": h.benchmark.Status(),
	})
}

func (h *TrustLabHandler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	var input trust.CreateArtifactInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	artifact, err := h.trust.CreateArtifact(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": artifact})
}

func (h *TrustLabHandler) ListArtifacts(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.trust.ListArtifacts()})
}

func (h *TrustLabHandler) CreateReplayCase(w http.ResponseWriter, r *http.Request) {
	var input replay.CreateCaseInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.replay.CreateCase(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": item})
}

func (h *TrustLabHandler) ListReplayCases(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.replay.ListCases()})
}

func (h *TrustLabHandler) RegisterBenchmarkRound(w http.ResponseWriter, r *http.Request) {
	var input benchmark.RegistrationRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.benchmark.RegisterRound(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": result})
}

func (h *TrustLabHandler) BenchmarkRoundStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.benchmark.Status()})
}
