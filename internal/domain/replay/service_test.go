package replay_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/platform/store"
)

type memoryClientStub struct {
	storeReplayErr error
	storedReplay   []memory.StoreReplayCaseRequest
}

func (m *memoryClientStub) Health(context.Context) error { return nil }
func (m *memoryClientStub) StoreReplayCase(_ context.Context, request memory.StoreReplayCaseRequest) error {
	if m.storeReplayErr != nil {
		return m.storeReplayErr
	}
	m.storedReplay = append(m.storedReplay, request)
	return nil
}
func (m *memoryClientStub) FetchSimilarCases(context.Context, memory.FetchSimilarCasesRequest) (memory.FetchSimilarCasesResponse, error) {
	return memory.FetchSimilarCasesResponse{}, nil
}
func (m *memoryClientStub) StoreTrustArtifact(context.Context, memory.StoreTrustArtifactRequest) error {
	return nil
}
func (m *memoryClientStub) LoadScenarioContext(context.Context, memory.LoadScenarioContextRequest) (memory.LoadScenarioContextResponse, error) {
	return memory.LoadScenarioContextResponse{}, nil
}

func TestCreateCaseWritesReplayArchive(t *testing.T) {
	replayStore, err := store.NewFileReplayStore(filepath.Join(t.TempDir(), "archive"))
	if err != nil {
		t.Fatalf("NewFileReplayStore() error = %v", err)
	}
	client := &memoryClientStub{}
	service := replay.NewService(replayStore, client)

	item, err := service.CreateCase(context.Background(), replay.CreateCaseInput{
		ScenarioID:              "scenario-1",
		TrustArtifactRefs:       []string{"ta-scenario-1"},
		BenchmarkRoundRef:       "bench-1",
		OutcomeSummary:          "baseline replay outcome",
		PromotionRecommendation: "promote",
		PromotionReason:         "matches expectations",
	})
	if err != nil {
		t.Fatalf("CreateCase() error = %v", err)
	}

	if item.ID == "" {
		t.Fatal("expected replay case id")
	}
	if len(service.ListCases()) != 1 {
		t.Fatalf("expected 1 replay case, got %d", len(service.ListCases()))
	}
	if len(client.storedReplay) != 1 {
		t.Fatalf("expected one clawmem write, got %d", len(client.storedReplay))
	}
}

func TestCreateCaseReturnsMemorySyncError(t *testing.T) {
	replayStore, err := store.NewFileReplayStore(filepath.Join(t.TempDir(), "archive"))
	if err != nil {
		t.Fatalf("NewFileReplayStore() error = %v", err)
	}
	service := replay.NewService(replayStore, &memoryClientStub{storeReplayErr: errors.New("clawmem unavailable")})

	_, err = service.CreateCase(context.Background(), replay.CreateCaseInput{
		ScenarioID:     "scenario-1",
		OutcomeSummary: "baseline replay outcome",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var syncErr *replay.MemorySyncError
	if !errors.As(err, &syncErr) {
		t.Fatalf("expected MemorySyncError, got %T", err)
	}
	if len(service.ListCases()) != 0 {
		t.Fatalf("expected no replay case persisted on clawmem failure, got %d", len(service.ListCases()))
	}
}
