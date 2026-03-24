package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("CONTROL_PLANE_BASE_URL", "http://127.0.0.1:8080")
	t.Setenv("MEMORY_BASE_URL", "http://127.0.0.1:8091")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ServiceAddress != "127.0.0.1:8090" {
		t.Fatalf("unexpected ServiceAddress: %s", cfg.ServiceAddress)
	}
	if cfg.ControlPlaneTimeout.String() != "5s" {
		t.Fatalf("unexpected ControlPlaneTimeout: %s", cfg.ControlPlaneTimeout)
	}
	if cfg.ScenarioPacksDir != "./configs/scenario-packs" {
		t.Fatalf("unexpected ScenarioPacksDir: %s", cfg.ScenarioPacksDir)
	}
}

func TestLoadRequiresURLs(t *testing.T) {
	if _, err := Load(); err == nil {
		t.Fatal("expected missing URL validation error")
	}
}
