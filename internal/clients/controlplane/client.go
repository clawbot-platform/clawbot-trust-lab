package controlplane

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type Client interface {
	Health(context.Context) error
	ListRuns(context.Context) ([]RunRef, error)
	CreateRun(context.Context, CreateRunRequest) (RunRef, error)
	ListPolicies(context.Context) ([]PolicyRef, error)
	CreatePolicy(context.Context, CreatePolicyRequest) (PolicyRef, error)
	RegisterBenchmarkMetadata(context.Context, benchmark.RegistrationRequest) (benchmark.RegistrationResult, error)
}

type RunRef struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PolicyRef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type CreateRunRequest struct {
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Status       string         `json:"status"`
	ScenarioType string         `json:"scenario_type"`
	MetadataJSON map[string]any `json:"metadata_json"`
}

type CreatePolicyRequest struct {
	Name        string         `json:"name"`
	Category    string         `json:"category"`
	Version     string         `json:"version"`
	Enabled     bool           `json:"enabled"`
	Description string         `json:"description"`
	RulesJSON   map[string]any `json:"rules_json"`
}

type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func New(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *HTTPClient) Health(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("control-plane health returned %d", res.StatusCode)
	}
	return nil
}

func (c *HTTPClient) ListRuns(ctx context.Context) ([]RunRef, error) {
	var response struct {
		Data []RunRef `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/runs", nil, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *HTTPClient) CreateRun(ctx context.Context, input CreateRunRequest) (RunRef, error) {
	var response struct {
		Data RunRef `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/runs", input, &response); err != nil {
		return RunRef{}, err
	}
	return response.Data, nil
}

func (c *HTTPClient) ListPolicies(ctx context.Context) ([]PolicyRef, error) {
	var response struct {
		Data []PolicyRef `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/policies", nil, &response); err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *HTTPClient) CreatePolicy(ctx context.Context, input CreatePolicyRequest) (PolicyRef, error) {
	var response struct {
		Data PolicyRef `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/policies", input, &response); err != nil {
		return PolicyRef{}, err
	}
	return response.Data, nil
}

func (c *HTTPClient) RegisterBenchmarkMetadata(_ context.Context, request benchmark.RegistrationRequest) (benchmark.RegistrationResult, error) {
	return benchmark.RegistrationResult{
		RegistrationID: "cp-stub-" + request.ScenarioPackID,
		Status:         "accepted_stub",
		RegisteredAt:   time.Now().UTC(),
	}, nil
}

func (c *HTTPClient) doJSON(ctx context.Context, method string, path string, request any, out any) error {
	var body *bytes.Reader
	if request == nil {
		body = bytes.NewReader(nil)
	} else {
		raw, err := json.Marshal(request)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("control-plane request failed with status %d", res.StatusCode)
	}

	if out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}
