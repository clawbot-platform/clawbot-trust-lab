package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/commerce"
	domainevents "clawbot-trust-lab/internal/domain/events"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	executionsvc "clawbot-trust-lab/internal/services/scenario"
)

type ScenarioService interface {
	ListPacks() []scenario.ScenarioPack
	GetPack(string) (scenario.ScenarioPack, error)
}

type ScenarioExecutionService interface {
	ListScenarios() []scenario.Scenario
	Execute(context.Context, string) (executionsvc.ExecutionResult, error)
}

type TrustService interface {
	CreateArtifact(context.Context, trust.CreateArtifactInput) (trust.TrustArtifact, error)
	ListArtifacts() []trust.TrustArtifact
	LoadMemoryContext(context.Context, string) (memory.LoadScenarioContextResponse, error)
}

type ReplayService interface {
	CreateCase(context.Context, replay.CreateCaseInput) (replay.ReplayCase, error)
	ListCases() []replay.ReplayCase
	SimilarCases(context.Context, string) (memory.FetchSimilarCasesResponse, error)
}

type BenchmarkService interface {
	RegisterRound(context.Context, benchmark.RegistrationRequest) (benchmark.RegistrationResult, error)
	Status() map[string]any
}

type CommerceService interface {
	ListOrders() []commerce.Order
	GetOrder(string) (commerce.Order, error)
}

type EventService interface {
	ListEvents() []domainevents.Record
}

type TrustDecisionService interface {
	ListDecisions() []trust.TrustDecision
	GetDecision(string) (trust.TrustDecision, error)
}

type TrustLabState struct {
	AppEnv          string
	ControlPlaneURL string
	ClawMemBaseURL  string
}

type TrustLabHandler struct {
	scenarios      ScenarioService
	execution      ScenarioExecutionService
	trust          TrustService
	replay         ReplayService
	benchmark      BenchmarkService
	commerce       CommerceService
	events         EventService
	trustDecisions TrustDecisionService
	state          TrustLabState
}

func NewTrustLabHandler(
	scenarios ScenarioService,
	execution ScenarioExecutionService,
	trustService TrustService,
	replayService ReplayService,
	benchmarkService BenchmarkService,
	commerceService CommerceService,
	eventService EventService,
	trustDecisionService TrustDecisionService,
	state TrustLabState,
) *TrustLabHandler {
	return &TrustLabHandler{
		scenarios:      scenarios,
		execution:      execution,
		trust:          trustService,
		replay:         replayService,
		benchmark:      benchmarkService,
		commerce:       commerceService,
		events:         eventService,
		trustDecisions: trustDecisionService,
		state:          state,
	}
}

func (h *TrustLabHandler) ScenarioTypes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"types": scenario.KnownTypes(),
		},
	})
}

func (h *TrustLabHandler) ListScenarios(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.execution.ListScenarios()})
}

func (h *TrustLabHandler) ExecuteScenario(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ScenarioID string `json:"scenario_id"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.execution.Execute(r.Context(), strings.TrimSpace(input.ScenarioID))
	if err != nil {
		var trustMemoryErr *trust.MemorySyncError
		var replayMemoryErr *replay.MemorySyncError
		if errors.As(err, &trustMemoryErr) || errors.As(err, &replayMemoryErr) || memory.IsDependencyFailure(err) {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": result})
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

func (h *TrustLabHandler) ReplayStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status":         "active",
		"memory_backend": "clawmem_http",
		"replay_cases":   len(h.replay.ListCases()),
		"phase":          "Phase 5",
	}
	if scenarioID := strings.TrimSpace(r.URL.Query().Get("scenario_id")); scenarioID != "" {
		response, err := h.replay.SimilarCases(r.Context(), scenarioID)
		if err != nil {
			data["memory_status"] = "degraded"
			data["memory_error"] = err.Error()
		} else {
			data["memory_status"] = "ok"
			data["similar_cases"] = response.Cases
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (h *TrustLabHandler) TrustStatus(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status":            "active",
		"control_plane_url": h.state.ControlPlaneURL,
		"clawmem_base_url":  h.state.ClawMemBaseURL,
		"memory_backend":    "clawmem_http",
		"artifact_count":    len(h.trust.ListArtifacts()),
		"trust_decisions":   len(h.trustDecisions.ListDecisions()),
		"artifact_families": []string{"trust_artifact", "mandate_artifact", "provenance_artifact"},
	}
	if scenarioID := strings.TrimSpace(r.URL.Query().Get("scenario_id")); scenarioID != "" {
		response, err := h.trust.LoadMemoryContext(r.Context(), scenarioID)
		if err != nil {
			data["memory_status"] = "degraded"
			data["memory_error"] = err.Error()
		} else {
			data["memory_status"] = "ok"
			data["memory_context"] = response.Context
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": data})
}

func (h *TrustLabHandler) BenchmarkStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.benchmark.Status()})
}

func (h *TrustLabHandler) CreateArtifact(w http.ResponseWriter, r *http.Request) {
	var input trust.CreateArtifactInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	artifact, err := h.trust.CreateArtifact(r.Context(), input)
	if err != nil {
		var memoryErr *trust.MemorySyncError
		if errors.As(err, &memoryErr) || memory.IsDependencyFailure(err) {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
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
	item, err := h.replay.CreateCase(r.Context(), input)
	if err != nil {
		var memoryErr *replay.MemorySyncError
		if errors.As(err, &memoryErr) || memory.IsDependencyFailure(err) {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": item})
}

func (h *TrustLabHandler) ListReplayCases(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.replay.ListCases()})
}

func (h *TrustLabHandler) ListOrders(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.commerce.ListOrders()})
}

func (h *TrustLabHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	item, err := h.commerce.GetOrder(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
}

func (h *TrustLabHandler) ListEvents(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.events.ListEvents()})
}

func (h *TrustLabHandler) ListTrustDecisions(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"data": h.trustDecisions.ListDecisions()})
}

func (h *TrustLabHandler) GetTrustDecision(w http.ResponseWriter, r *http.Request) {
	item, err := h.trustDecisions.GetDecision(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": item})
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
