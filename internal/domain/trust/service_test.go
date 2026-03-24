package trust_test

import (
	"context"
	"errors"
	"testing"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/store"
)

type scenarioLookupStub struct{}

func (scenarioLookupStub) GetScenario(string) (scenario.Scenario, error) {
	return scenario.Scenario{
		ID:          "scenario-1",
		Name:        "Mandate review",
		Type:        scenario.ScenarioTypeMandateReview,
		Description: "test",
		PackID:      "starter-pack",
		Version:     "v1",
	}, nil
}

type memoryClientStub struct {
	storeTrustErr error
	storedTrust   []memory.StoreTrustArtifactRequest
	contextResp   memory.LoadScenarioContextResponse
}

func (m *memoryClientStub) Health(context.Context) error { return nil }
func (m *memoryClientStub) StoreReplayCase(context.Context, memory.StoreReplayCaseRequest) error {
	return nil
}
func (m *memoryClientStub) FetchSimilarCases(context.Context, memory.FetchSimilarCasesRequest) (memory.FetchSimilarCasesResponse, error) {
	return memory.FetchSimilarCasesResponse{}, nil
}
func (m *memoryClientStub) StoreTrustArtifact(_ context.Context, request memory.StoreTrustArtifactRequest) error {
	if m.storeTrustErr != nil {
		return m.storeTrustErr
	}
	m.storedTrust = append(m.storedTrust, request)
	return nil
}
func (m *memoryClientStub) LoadScenarioContext(context.Context, memory.LoadScenarioContextRequest) (memory.LoadScenarioContextResponse, error) {
	return m.contextResp, nil
}

func TestCreateArtifactStoresAndUsesMemory(t *testing.T) {
	artifacts := store.NewInMemoryTrustArtifactStore()
	client := &memoryClientStub{}
	service := trust.NewService(scenarioLookupStub{}, artifacts, client)

	artifact, err := service.CreateArtifact(context.Background(), trust.CreateArtifactInput{ScenarioID: "scenario-1"})
	if err != nil {
		t.Fatalf("CreateArtifact() error = %v", err)
	}

	if artifact.ID == "" {
		t.Fatal("expected artifact id")
	}
	if len(service.ListArtifacts()) != 1 {
		t.Fatalf("expected stored artifact, got %d", len(service.ListArtifacts()))
	}
	if len(client.storedTrust) != 1 {
		t.Fatalf("expected one clawmem write, got %d", len(client.storedTrust))
	}
}

func TestCreateArtifactReturnsMemorySyncError(t *testing.T) {
	artifacts := store.NewInMemoryTrustArtifactStore()
	client := &memoryClientStub{storeTrustErr: errors.New("clawmem unavailable")}
	service := trust.NewService(scenarioLookupStub{}, artifacts, client)

	_, err := service.CreateArtifact(context.Background(), trust.CreateArtifactInput{ScenarioID: "scenario-1"})
	if err == nil {
		t.Fatal("expected error")
	}
	var syncErr *trust.MemorySyncError
	if !errors.As(err, &syncErr) {
		t.Fatalf("expected MemorySyncError, got %T", err)
	}
	if len(service.ListArtifacts()) != 0 {
		t.Fatalf("expected no stored artifacts on clawmem failure, got %d", len(service.ListArtifacts()))
	}
}
