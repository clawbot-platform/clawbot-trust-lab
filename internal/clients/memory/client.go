package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	Health(context.Context) error
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
	Tags         []string       `json:"tags"`
}

type FetchSimilarCasesRequest struct {
	ScenarioID string `json:"scenario_id"`
	Query      string `json:"query"`
}

type FetchSimilarCasesResponse struct {
	Cases []map[string]any `json:"cases"`
}

type StoreTrustArtifactRequest struct {
	ArtifactID     string         `json:"artifact_id"`
	ScenarioID     string         `json:"scenario_id"`
	Summary        string         `json:"summary"`
	ArtifactFamily string         `json:"artifact_family"`
	ArtifactType   string         `json:"artifact_type"`
	Metadata       map[string]any `json:"metadata"`
	Tags           []string       `json:"tags"`
}

type LoadScenarioContextRequest struct {
	ScenarioID string `json:"scenario_id"`
}

type LoadScenarioContextResponse struct {
	ScenarioID string         `json:"scenario_id"`
	Context    map[string]any `json:"context"`
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

type StatusError struct {
	Operation  string
	StatusCode int
	Message    string
}

func (e *StatusError) Error() string {
	message := strings.TrimSpace(e.Message)
	if message == "" {
		message = http.StatusText(e.StatusCode)
	}
	return fmt.Sprintf("%s failed with status %d: %s", e.Operation, e.StatusCode, message)
}

func New(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *HTTPClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return fmt.Errorf("build clawmem health request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call clawmem health endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return readStatusError("clawmem health check", resp)
	}
	return nil
}

func (c *HTTPClient) StoreReplayCase(ctx context.Context, request StoreReplayCaseRequest) error {
	payload := map[string]any{
		"scenario_id": request.ScenarioID,
		"source_id":   request.ReplayCaseID,
		"summary":     request.Summary,
		"metadata":    cloneMap(request.Metadata),
		"tags":        append([]string(nil), request.Tags...),
	}
	return c.post(ctx, "/api/v1/replay", payload, nil, "store replay memory")
}

func (c *HTTPClient) FetchSimilarCases(ctx context.Context, request FetchSimilarCasesRequest) (FetchSimilarCasesResponse, error) {
	values := url.Values{}
	if scenarioID := strings.TrimSpace(request.ScenarioID); scenarioID != "" {
		values.Set("scenario_id", scenarioID)
	}

	var response struct {
		Records []map[string]any `json:"records"`
		Total   int              `json:"total"`
	}
	if err := c.get(ctx, "/api/v1/replay", values, &response, "load replay memories"); err != nil {
		return FetchSimilarCasesResponse{}, err
	}

	filtered := make([]map[string]any, 0, len(response.Records))
	query := strings.TrimSpace(request.Query)
	for _, record := range response.Records {
		if query == "" {
			filtered = append(filtered, record)
			continue
		}
		summary, _ := record["outcome_summary"].(string)
		if strings.Contains(strings.ToLower(summary), strings.ToLower(query)) {
			filtered = append(filtered, record)
		}
	}

	return FetchSimilarCasesResponse{Cases: filtered}, nil
}

func (c *HTTPClient) StoreTrustArtifact(ctx context.Context, request StoreTrustArtifactRequest) error {
	metadata := cloneMap(request.Metadata)
	if strings.TrimSpace(request.ArtifactFamily) != "" {
		metadata["artifact_family"] = request.ArtifactFamily
	}
	if strings.TrimSpace(request.ArtifactType) != "" {
		metadata["artifact_type"] = request.ArtifactType
	}

	payload := map[string]any{
		"scenario_id":     request.ScenarioID,
		"source_id":       request.ArtifactID,
		"summary":         request.Summary,
		"artifact_family": request.ArtifactFamily,
		"artifact_type":   request.ArtifactType,
		"metadata":        metadata,
		"tags":            append([]string(nil), request.Tags...),
	}
	return c.post(ctx, "/api/v1/trust", payload, nil, "store trust memory")
}

func (c *HTTPClient) LoadScenarioContext(ctx context.Context, request LoadScenarioContextRequest) (LoadScenarioContextResponse, error) {
	values := url.Values{}
	values.Set("scenario_id", request.ScenarioID)

	var response struct {
		Records []map[string]any `json:"records"`
		Total   int              `json:"total"`
	}
	if err := c.get(ctx, "/api/v1/memories", values, &response, "load scenario context"); err != nil {
		return LoadScenarioContextResponse{}, err
	}

	return LoadScenarioContextResponse{
		ScenarioID: request.ScenarioID,
		Context: map[string]any{
			"memory_source": c.baseURL,
			"record_count":  response.Total,
			"records":       response.Records,
		},
	}, nil
}

func (c *HTTPClient) post(ctx context.Context, path string, payload any, out any, operation string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal %s payload: %w", operation, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build %s request: %w", operation, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s endpoint: %w", operation, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readStatusError(operation, resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s response: %w", operation, err)
	}
	return nil
}

func (c *HTTPClient) get(ctx context.Context, path string, values url.Values, out any, operation string) error {
	endpoint := c.baseURL + path
	if encoded := values.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build %s request: %w", operation, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s endpoint: %w", operation, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return readStatusError(operation, resp)
	}
	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s response: %w", operation, err)
	}
	return nil
}

func readStatusError(operation string, resp *http.Response) error {
	message := ""
	payload, err := io.ReadAll(resp.Body)
	if err == nil {
		var envelope struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(payload, &envelope) == nil {
			message = envelope.Error.Message
		}
		if message == "" {
			message = strings.TrimSpace(string(payload))
		}
	}
	return &StatusError{
		Operation:  operation,
		StatusCode: resp.StatusCode,
		Message:    message,
	}
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}

func IsDependencyFailure(err error) bool {
	var statusErr *StatusError
	if errors.As(err, &statusErr) {
		return true
	}
	return strings.Contains(err.Error(), "call clawmem")
}
