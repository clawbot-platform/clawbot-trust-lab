package memory

import (
	"context"
)

type Client interface {
	StoreReplayCase(context.Context, StoreReplayCaseRequest) error
	FetchSimilarCases(context.Context, FetchSimilarCasesRequest) (FetchSimilarCasesResponse, error)
	StoreTrustArtifact(context.Context, StoreTrustArtifactRequest) error
	LoadScenarioContext(context.Context, LoadScenarioContextRequest) (LoadScenarioContextResponse, error)
}

type StoreReplayCaseRequest struct {
	ReplayCaseID string         `json:"replay_case_id"`
	ScenarioID   string         `json:"scenario_id"`
	Summary      string         `json:"summary"`
	Metadata     map[string]any `json:"metadata"`
}

type FetchSimilarCasesRequest struct {
	ScenarioID string `json:"scenario_id"`
	Query      string `json:"query"`
}

type FetchSimilarCasesResponse struct {
	Cases []map[string]any `json:"cases"`
}

type StoreTrustArtifactRequest struct {
	ArtifactID string         `json:"artifact_id"`
	ScenarioID string         `json:"scenario_id"`
	Summary    string         `json:"summary"`
	Metadata   map[string]any `json:"metadata"`
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
	return FetchSimilarCasesResponse{Cases: []map[string]any{}}, nil
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
