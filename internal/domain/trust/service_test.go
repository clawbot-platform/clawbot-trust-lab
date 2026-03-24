package trust_test

import (
	"context"
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

func TestCreateArtifactStoresAndUsesMemory(t *testing.T) {
	artifacts := store.NewInMemoryTrustArtifactStore()
	client := memory.NewStub("http://127.0.0.1:8091")
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
}
