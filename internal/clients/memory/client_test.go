package memory

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, time.Second)
	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("Health() error = %v", err)
	}
}

func TestStoreTrustArtifact(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/trust" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("Decode() error = %v", err)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := New(server.URL, time.Second)
	err := client.StoreTrustArtifact(context.Background(), StoreTrustArtifactRequest{
		ArtifactID:     "ta-1",
		ScenarioID:     "scenario-1",
		Summary:        "artifact summary",
		ArtifactFamily: "mandate",
		ArtifactType:   "mandate_artifact",
		Metadata:       map[string]any{"source": "test"},
		Tags:           []string{"trust"},
	})
	if err != nil {
		t.Fatalf("StoreTrustArtifact() error = %v", err)
	}

	if got["source_id"] != "ta-1" {
		t.Fatalf("unexpected source_id: %#v", got["source_id"])
	}
}

func TestLoadScenarioContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memories" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("scenario_id") != "scenario-1" {
			t.Fatalf("unexpected scenario id: %s", r.URL.Query().Get("scenario_id"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"records": []map[string]any{{"id": "mem-1", "summary": "summary"}},
			"total":   1,
		})
	}))
	defer server.Close()

	client := New(server.URL, time.Second)
	response, err := client.LoadScenarioContext(context.Background(), LoadScenarioContextRequest{ScenarioID: "scenario-1"})
	if err != nil {
		t.Fatalf("LoadScenarioContext() error = %v", err)
	}

	if response.ScenarioID != "scenario-1" {
		t.Fatalf("unexpected ScenarioID: %s", response.ScenarioID)
	}
	if response.Context["record_count"] != float64(1) && response.Context["record_count"] != 1 {
		t.Fatalf("unexpected record_count: %#v", response.Context["record_count"])
	}
}

func TestStoreReplayCasePropagatesStatusErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{"message": "upstream unavailable"},
		})
	}))
	defer server.Close()

	client := New(server.URL, time.Second)
	err := client.StoreReplayCase(context.Background(), StoreReplayCaseRequest{
		ReplayCaseID: "rc-1",
		ScenarioID:   "scenario-1",
		Summary:      "summary",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
