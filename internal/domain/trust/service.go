package trust

import (
	"context"
	"fmt"
	"strings"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/scenario"
)

type ScenarioLookup interface {
	GetScenario(string) (scenario.Scenario, error)
}

type Service struct {
	scenarios ScenarioLookup
	store     ArtifactStore
	memory    memory.Client
}

type ArtifactStore interface {
	Create(TrustArtifact) error
	List() []TrustArtifact
}

type CreateArtifactInput struct {
	ScenarioID string `json:"scenario_id"`
}

func NewService(scenarios ScenarioLookup, artifactStore ArtifactStore, memoryClient memory.Client) *Service {
	return &Service{scenarios: scenarios, store: artifactStore, memory: memoryClient}
}

func (s *Service) CreateArtifact(ctx context.Context, input CreateArtifactInput) (TrustArtifact, error) {
	if strings.TrimSpace(input.ScenarioID) == "" {
		return TrustArtifact{}, fmt.Errorf("scenario_id is required")
	}

	item, err := s.scenarios.GetScenario(input.ScenarioID)
	if err != nil {
		return TrustArtifact{}, err
	}

	artifact := TrustArtifact{
		ID:               "ta-" + item.ID,
		ArtifactFamily:   "trust",
		ArtifactType:     deriveArtifactType(item),
		SourceScenarioID: item.ID,
		Summary:          "Trust artifact created for scenario " + item.Name,
		Metadata: map[string]any{
			"scenario_pack_id": item.PackID,
			"scenario_type":    item.Type,
			"tags":             item.Tags,
			"trust_signals":    item.TrustSignals,
		},
		PolicyDecision: PolicyDecisionRef{
			PolicyID:      "phase-3-placeholder",
			PolicyVersion: "v1",
			Outcome:       "scaffolded",
		},
		CreatedAt: time.Now().UTC(),
	}

	if item.Type == scenario.ScenarioTypeMandateReview {
		artifact.Mandate = &MandateArtifact{
			Source:      "scenario-pack",
			Title:       item.Name,
			Description: item.Description,
		}
	} else {
		artifact.Provenance = &ProvenanceArtifact{
			SourceRepo: item.PackID,
			Revision:   item.Version,
			RecordedBy: "trust-lab-phase-3",
		}
	}

	if err := s.store.Create(artifact); err != nil {
		return TrustArtifact{}, err
	}
	if err := s.memory.StoreTrustArtifact(ctx, memory.StoreTrustArtifactRequest{
		ArtifactID: artifact.ID,
		ScenarioID: artifact.SourceScenarioID,
		Summary:    artifact.Summary,
		Metadata:   artifact.Metadata,
	}); err != nil {
		return TrustArtifact{}, err
	}

	return artifact, nil
}

func (s *Service) ListArtifacts() []TrustArtifact {
	return s.store.List()
}

func deriveArtifactType(item scenario.Scenario) string {
	if item.Type == scenario.ScenarioTypeMandateReview {
		return "mandate_artifact"
	}
	return "provenance_artifact"
}
