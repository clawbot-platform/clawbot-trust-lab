package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "starter-pack.json")
	content := `{
  "id": "starter-pack",
  "name": "Starter Pack",
  "version": "v1",
  "description": "test pack",
  "scenarios": [
    {
      "id": "scenario-1",
      "code": "H1",
      "name": "Mandate Review",
      "version": "v1",
      "scenario_type": "mandate_review",
      "family": "commerce",
      "set_role": "stable",
      "description": "test scenario",
      "actors": ["reviewer"],
      "trust_signals": ["signal-a"],
      "expected_outcomes": ["artifact"],
      "tags": ["starter"],
      "feature_model": {
        "tier_a": ["psp_amount"],
        "tier_b": ["historical_refund_count"],
        "tier_c": ["mandate_status"]
      }
    }
  ]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	loader := New(dir)
	packs, err := loader.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() error = %v", err)
	}

	if len(packs) != 1 {
		t.Fatalf("expected 1 pack, got %d", len(packs))
	}
	if packs[0].Scenarios[0].ID != "scenario-1" {
		t.Fatalf("unexpected scenario payload: %#v", packs[0].Scenarios[0])
	}
	if packs[0].Scenarios[0].Code != "H1" {
		t.Fatalf("expected code H1, got %#v", packs[0].Scenarios[0])
	}
	if packs[0].Scenarios[0].SetRole != "stable" {
		t.Fatalf("expected stable set role, got %#v", packs[0].Scenarios[0].SetRole)
	}
	if len(packs[0].Scenarios[0].FeatureModel.TierB) != 1 {
		t.Fatalf("expected tier model to load, got %#v", packs[0].Scenarios[0].FeatureModel)
	}
}
