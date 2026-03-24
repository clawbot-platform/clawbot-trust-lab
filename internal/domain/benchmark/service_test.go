package benchmark

import (
	"context"
	"testing"
	"time"
)

type controlPlaneStub struct{}

func (controlPlaneStub) RegisterBenchmarkMetadata(_ context.Context, _ RegistrationRequest) (RegistrationResult, error) {
	return RegistrationResult{
		RegistrationID: "bench-1",
		Status:         "accepted_stub",
		RegisteredAt:   time.Now().UTC(),
	}, nil
}

func TestRegisterRound(t *testing.T) {
	service := NewService(controlPlaneStub{})

	result, err := service.RegisterRound(context.Background(), RegistrationRequest{
		ScenarioPackID:      "starter-pack",
		ScenarioPackVersion: "v1",
		ReplayCaseRefs:      []string{"rc-1"},
	})
	if err != nil {
		t.Fatalf("RegisterRound() error = %v", err)
	}

	if result.RegistrationID == "" {
		t.Fatal("expected registration id")
	}
	if service.Status()["registrations"] != 1 {
		t.Fatalf("unexpected status: %#v", service.Status())
	}
}
