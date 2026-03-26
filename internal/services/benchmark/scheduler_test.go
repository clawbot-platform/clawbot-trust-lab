package benchmark

import (
	"context"
	"errors"
	"testing"
	"time"

	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
)

func TestConfigureSchedulerUpdatesStatus(t *testing.T) {
	service := NewService(registrarStub{}, nil, nil, nil, nil, nil)

	service.ConfigureScheduler(SchedulerConfig{
		Enabled:        true,
		ScenarioFamily: "commerce",
		Interval:       5 * time.Second,
		MaxRuns:        3,
		DryRun:         true,
	})

	status := service.SchedulerStatus()
	if !status.Enabled || status.ScenarioFamily != "commerce" || status.Interval != "5s" || status.MaxRuns != 3 || !status.DryRun {
		t.Fatalf("unexpected scheduler status after configure: %#v", status)
	}
}

func TestStartSchedulerGuardConditions(t *testing.T) {
	tests := []struct {
		name   string
		cfg    SchedulerConfig
		status domainbenchmark.SchedulerStatus
	}{
		{name: "disabled", cfg: SchedulerConfig{Enabled: false, Interval: time.Millisecond, MaxRuns: 1}},
		{name: "invalid interval", cfg: SchedulerConfig{Enabled: true, Interval: 0, MaxRuns: 1}},
		{name: "zero max runs", cfg: SchedulerConfig{Enabled: true, Interval: time.Millisecond, MaxRuns: 0}},
		{name: "already running", cfg: SchedulerConfig{Enabled: true, Interval: time.Millisecond, MaxRuns: 1}, status: domainbenchmark.SchedulerStatus{Running: true}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewService(registrarStub{}, nil, nil, nil, nil, nil)
			service.schedulerConfig = tc.cfg
			service.schedulerStatus = tc.status

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			service.StartScheduler(ctx)

			if service.SchedulerStatus().Running != tc.status.Running {
				t.Fatalf("expected scheduler running=%t, got %#v", tc.status.Running, service.SchedulerStatus())
			}
		})
	}
}

func TestStartSchedulerStopsOnContextCancel(t *testing.T) {
	service := newRoundService(t)
	service.now = func() time.Time { return time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC) }
	service.ConfigureScheduler(SchedulerConfig{
		Enabled:        true,
		ScenarioFamily: "commerce",
		Interval:       50 * time.Millisecond,
		MaxRuns:        2,
	})

	ctx, cancel := context.WithCancel(context.Background())
	service.StartScheduler(ctx)

	status := service.SchedulerStatus()
	if !status.Running || status.NextRunAt.IsZero() {
		t.Fatalf("expected scheduler to start immediately, got %#v", status)
	}

	cancel()
	waitForCondition(t, time.Second, func() bool { return !service.SchedulerStatus().Running })

	status = service.SchedulerStatus()
	if status.Running || !status.NextRunAt.IsZero() {
		t.Fatalf("expected scheduler to stop and clear next run after cancel, got %#v", status)
	}
}

func TestStartSchedulerExecutesUntilMaxRuns(t *testing.T) {
	service := newRoundService(t)
	counter := 0
	service.now = func() time.Time {
		base := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
		value := base.Add(time.Duration(counter) * time.Minute)
		counter++
		return value
	}
	service.ConfigureScheduler(SchedulerConfig{
		Enabled:        true,
		ScenarioFamily: "commerce",
		Interval:       2 * time.Millisecond,
		MaxRuns:        1,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	service.StartScheduler(ctx)

	waitForCondition(t, 2*time.Second, func() bool { return service.SchedulerStatus().ExecutedRuns >= 1 })
	waitForCondition(t, 2*time.Second, func() bool { return !service.SchedulerStatus().Running })

	status := service.SchedulerStatus()
	if status.ExecutedRuns < 1 || status.LastRoundID == "" || status.LastStartedAt.IsZero() {
		t.Fatalf("expected scheduler execution to be recorded, got %#v", status)
	}
}

func TestRunScheduledResolveDefaultsAndErrors(t *testing.T) {
	cfg := SchedulerConfig{ScenarioFamily: "commerce", Interval: 5 * time.Second}

	if _, err := resolveSchedulerInterval(cfg, domainbenchmark.SchedulerControlInput{Interval: "bad"}); err == nil {
		t.Fatal("expected invalid scheduler interval to fail")
	}
	if interval, err := resolveSchedulerInterval(SchedulerConfig{}, domainbenchmark.SchedulerControlInput{}); err != nil || interval != time.Second {
		t.Fatalf("expected zero interval to default to 1s, got %s err=%v", interval, err)
	}
	if got := resolveSchedulerScenarioFamily(SchedulerConfig{}, domainbenchmark.SchedulerControlInput{}); got != defaultSchedulerScenarioFamily {
		t.Fatalf("expected default scenario family %s, got %s", defaultSchedulerScenarioFamily, got)
	}
	if got := resolveSchedulerScenarioFamily(cfg, domainbenchmark.SchedulerControlInput{}); got != "commerce" {
		t.Fatalf("expected config scenario family to be used, got %s", got)
	}
	if got := resolveSchedulerScenarioFamily(cfg, domainbenchmark.SchedulerControlInput{ScenarioFamily: "custom"}); got != "custom" {
		t.Fatalf("expected explicit scenario family to win, got %s", got)
	}
	if got := resolveSchedulerMaxRuns(domainbenchmark.SchedulerControlInput{}); got != 1 {
		t.Fatalf("expected max runs default 1, got %d", got)
	}
}

func TestRunScheduledStopsOnContextCancellationBetweenRuns(t *testing.T) {
	service := newRoundService(t)
	counter := 0
	service.now = func() time.Time {
		base := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
		value := base.Add(time.Duration(counter) * time.Minute)
		counter++
		return value
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	items, err := service.RunScheduled(ctx, domainbenchmark.SchedulerControlInput{
		ScenarioFamily: "commerce",
		Interval:       "20ms",
		MaxRuns:        2,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one completed run before cancellation, got %d", len(items))
	}
	if service.SchedulerStatus().Running {
		t.Fatalf("expected scheduler status to stop after cancellation, got %#v", service.SchedulerStatus())
	}
}

func TestRecordSchedulerExecutionUpdatesStatus(t *testing.T) {
	service := NewService(registrarStub{}, nil, nil, nil, nil, nil)
	service.schedulerConfig.Interval = 10 * time.Second
	service.schedulerStatus.Running = true

	startedAt := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
	service.recordSchedulerExecution("round-1", startedAt)

	status := service.SchedulerStatus()
	if status.ExecutedRuns != 1 || status.LastRoundID != "round-1" || !status.LastStartedAt.Equal(startedAt) {
		t.Fatalf("unexpected scheduler execution status %#v", status)
	}
	if !status.NextRunAt.Equal(startedAt.Add(10 * time.Second)) {
		t.Fatalf("expected next run to follow scheduler interval, got %#v", status)
	}

	service.schedulerStatus.Running = false
	service.schedulerStatus.NextRunAt = time.Time{}
	service.recordSchedulerExecution("round-2", startedAt.Add(time.Minute))
	if !service.SchedulerStatus().NextRunAt.IsZero() {
		t.Fatalf("expected stopped scheduler not to set next run, got %#v", service.SchedulerStatus())
	}
}

func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met before timeout")
}
