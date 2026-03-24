package memory

import (
	"context"
	"testing"
)

func TestStubClientLoadScenarioContext(t *testing.T) {
	client := NewStub("http://127.0.0.1:8091")

	response, err := client.LoadScenarioContext(context.Background(), LoadScenarioContextRequest{ScenarioID: "scenario-1"})
	if err != nil {
		t.Fatalf("LoadScenarioContext() error = %v", err)
	}

	if response.ScenarioID != "scenario-1" {
		t.Fatalf("unexpected ScenarioID: %s", response.ScenarioID)
	}
}
