package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/version"
)

type scenarioCatalogStub struct{}

func (scenarioCatalogStub) ScenarioTypes() []scenario.ScenarioType {
	return []scenario.ScenarioType{scenario.ScenarioTypeMandateReview}
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

func TestTrustLabHandlerScenarioTypes(t *testing.T) {
	handler := NewTrustLabHandler(scenarioCatalogStub{}, TrustLabState{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios/types", nil)
	recorder := httptest.NewRecorder()

	handler.ScenarioTypes(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response map[string]map[string][]scenario.ScenarioType
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if len(response["data"]["types"]) != 1 {
		t.Fatalf("unexpected scenario type response: %#v", response)
	}
}
