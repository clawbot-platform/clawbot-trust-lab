package trust

import (
	"context"
	"errors"
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
	ScenarioID string    `json:"scenario_id"`
	ArtifactID string    `json:"artifact_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

type MemorySyncError struct {
	Err error
}

func (e *MemorySyncError) Error() string {
	return "clawmem trust artifact write failed: " + e.Err.Error()
}

func (e *MemorySyncError) Unwrap() error {
	return e.Err
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
		ArtifactFamily:   deriveArtifactFamily(item),
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
	if strings.TrimSpace(input.ArtifactID) != "" {
		artifact.ID = strings.TrimSpace(input.ArtifactID)
	}
	if !input.CreatedAt.IsZero() {
		artifact.CreatedAt = input.CreatedAt.UTC()
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

	if err := s.memory.StoreTrustArtifact(ctx, memory.StoreTrustArtifactRequest{
		ArtifactID:     artifact.ID,
		ScenarioID:     artifact.SourceScenarioID,
		Summary:        artifact.Summary,
		ArtifactFamily: artifact.ArtifactFamily,
		ArtifactType:   artifact.ArtifactType,
		Metadata:       artifact.Metadata,
		Tags:           append([]string{"trust-artifact", string(item.Type)}, item.Tags...),
	}); err != nil {
		return TrustArtifact{}, &MemorySyncError{Err: err}
	}
	if err := s.store.Create(artifact); err != nil {
		return TrustArtifact{}, err
	}

	return artifact, nil
}

func (s *Service) ListArtifacts() []TrustArtifact {
	return s.store.List()
}

func (s *Service) LoadMemoryContext(ctx context.Context, scenarioID string) (memory.LoadScenarioContextResponse, error) {
	if strings.TrimSpace(scenarioID) == "" {
		return memory.LoadScenarioContextResponse{}, errors.New("scenario_id is required")
	}
	return s.memory.LoadScenarioContext(ctx, memory.LoadScenarioContextRequest{ScenarioID: scenarioID})
}

func deriveArtifactType(item scenario.Scenario) string {
	if item.Type == scenario.ScenarioTypeMandateReview {
		return "mandate_artifact"
	}
	return "provenance_artifact"
}

func deriveArtifactFamily(item scenario.Scenario) string {
	if item.Type == scenario.ScenarioTypeMandateReview {
		return "mandate"
	}
	return "provenance"
}
