package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/version"
)

type scenarioServiceStub struct{}

func (scenarioServiceStub) ListPacks() []scenario.ScenarioPack {
	return []scenario.ScenarioPack{{ID: "starter-pack", Name: "Starter Pack"}}
}
func (scenarioServiceStub) GetPack(string) (scenario.ScenarioPack, error) {
	return scenario.ScenarioPack{ID: "starter-pack", Name: "Starter Pack"}, nil
}

type trustServiceStub struct {
	items []trust.TrustArtifact
}

func (s *trustServiceStub) CreateArtifact(_ context.Context, input trust.CreateArtifactInput) (trust.TrustArtifact, error) {
	item := trust.TrustArtifact{ID: "ta-" + input.ScenarioID, SourceScenarioID: input.ScenarioID}
	s.items = append(s.items, item)
	return item, nil
}
func (s *trustServiceStub) ListArtifacts() []trust.TrustArtifact { return s.items }

type replayServiceStub struct {
	items []replay.ReplayCase
}

func (s *replayServiceStub) CreateCase(input replay.CreateCaseInput) (replay.ReplayCase, error) {
	item := replay.ReplayCase{ID: "rc-1", ScenarioID: input.ScenarioID}
	s.items = append(s.items, item)
	return item, nil
}
func (s *replayServiceStub) ListCases() []replay.ReplayCase { return s.items }

type benchmarkServiceStub struct{}

func (benchmarkServiceStub) RegisterRound(_ context.Context, _ benchmark.RegistrationRequest) (benchmark.RegistrationResult, error) {
	return benchmark.RegistrationResult{RegistrationID: "bench-1", Status: "accepted_stub", RegisteredAt: time.Now().UTC()}, nil
}
func (benchmarkServiceStub) Status() map[string]any {
	return map[string]any{"registrations": 1, "last_status": "accepted_stub"}
}

func TestSystemHandlerHealth(t *testing.T) {
	handler := NewSystemHandler(func(context.Context) error { return nil }, version.Current())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerListPacks(t *testing.T) {
	handler := NewTrustLabHandler(scenarioServiceStub{}, &trustServiceStub{}, &replayServiceStub{}, benchmarkServiceStub{}, TrustLabState{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios/packs", nil)
	recorder := httptest.NewRecorder()

	handler.ListPacks(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response map[string][]scenario.ScenarioPack
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(response["data"]) != 1 {
		t.Fatalf("unexpected scenario pack response: %#v", response)
	}
}

func TestTrustLabHandlerCreateArtifact(t *testing.T) {
	trustStub := &trustServiceStub{}
	handler := NewTrustLabHandler(scenarioServiceStub{}, trustStub, &replayServiceStub{}, benchmarkServiceStub{}, TrustLabState{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/trust/artifacts", bytes.NewBufferString(`{"scenario_id":"starter-mandate-review"}`))
	recorder := httptest.NewRecorder()

	handler.CreateArtifact(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
	if len(trustStub.items) != 1 {
		t.Fatalf("expected artifact creation, got %d", len(trustStub.items))
	}
}

func TestTrustLabHandlerCreateReplayCase(t *testing.T) {
	replayStub := &replayServiceStub{}
	handler := NewTrustLabHandler(scenarioServiceStub{}, &trustServiceStub{}, replayStub, benchmarkServiceStub{}, TrustLabState{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/replay/cases", bytes.NewBufferString(`{"scenario_id":"starter-mandate-review","outcome_summary":"ok"}`))
	recorder := httptest.NewRecorder()

	handler.CreateReplayCase(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
	if len(replayStub.items) != 1 {
		t.Fatalf("expected replay creation, got %d", len(replayStub.items))
	}
}

func TestTrustLabHandlerRegisterBenchmarkRound(t *testing.T) {
	handler := NewTrustLabHandler(scenarioServiceStub{}, &trustServiceStub{}, &replayServiceStub{}, benchmarkServiceStub{}, TrustLabState{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/benchmark/rounds/register", bytes.NewBufferString(`{"scenario_pack_id":"starter-pack","scenario_pack_version":"v1","replay_case_refs":["rc-1"]}`))
	recorder := httptest.NewRecorder()

	handler.RegisterBenchmarkRound(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}
