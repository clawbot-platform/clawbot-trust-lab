package memory

import (
	"context"

	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/trust"
)

type Client interface {
	StoreReplayCase(context.Context, StoreReplayCaseRequest) error
	FetchSimilarCases(context.Context, FetchSimilarCasesRequest) (FetchSimilarCasesResponse, error)
	StoreTrustArtifact(context.Context, StoreTrustArtifactRequest) error
	LoadScenarioContext(context.Context, LoadScenarioContextRequest) (LoadScenarioContextResponse, error)
}

type StoreReplayCaseRequest struct {
	Case replay.ReplayCase `json:"case"`
}

type FetchSimilarCasesRequest struct {
	ScenarioID string `json:"scenario_id"`
	Query      string `json:"query"`
}

type FetchSimilarCasesResponse struct {
	Cases []replay.ReplayCase `json:"cases"`
}

type StoreTrustArtifactRequest struct {
	Artifact trust.TrustArtifact `json:"artifact"`
}

type LoadScenarioContextRequest struct {
	ScenarioID string `json:"scenario_id"`
}

type LoadScenarioContextResponse struct {
	ScenarioID string         `json:"scenario_id"`
	Context    map[string]any `json:"context"`
}

type StubClient struct {
	baseURL string
}

func NewStub(baseURL string) *StubClient {
	return &StubClient{baseURL: baseURL}
}

func (c *StubClient) BaseURL() string {
	return c.baseURL
}

func (c *StubClient) StoreReplayCase(context.Context, StoreReplayCaseRequest) error {
	return nil
}

func (c *StubClient) FetchSimilarCases(context.Context, FetchSimilarCasesRequest) (FetchSimilarCasesResponse, error) {
	return FetchSimilarCasesResponse{Cases: []replay.ReplayCase{}}, nil
}

func (c *StubClient) StoreTrustArtifact(context.Context, StoreTrustArtifactRequest) error {
	return nil
}

func (c *StubClient) LoadScenarioContext(_ context.Context, request LoadScenarioContextRequest) (LoadScenarioContextResponse, error) {
	return LoadScenarioContextResponse{
		ScenarioID: request.ScenarioID,
		Context: map[string]any{
			"source": "stub-memory-client",
			"note":   "Phase 2 defines the clawmem contract without implementing storage internals.",
		},
	}, nil
}
