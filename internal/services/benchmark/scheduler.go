package benchmark

import (
	"context"
	"time"

	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
)

type SchedulerConfig struct {
	Enabled        bool
	ScenarioFamily string
	Interval       time.Duration
	MaxRuns        int
	DryRun         bool
}

const defaultSchedulerScenarioFamily = "commerce"

func (s *Service) ConfigureScheduler(cfg SchedulerConfig) {
	s.schedulerMu.Lock()
	defer s.schedulerMu.Unlock()

	s.schedulerConfig = cfg
	s.schedulerStatus.Enabled = cfg.Enabled
	s.schedulerStatus.ScenarioFamily = cfg.ScenarioFamily
	s.schedulerStatus.Interval = cfg.Interval.String()
	s.schedulerStatus.MaxRuns = cfg.MaxRuns
	s.schedulerStatus.DryRun = cfg.DryRun
}

func (s *Service) StartScheduler(ctx context.Context) {
	s.schedulerMu.Lock()
	cfg := s.schedulerConfig
	if !schedulerRunnable(cfg, s.schedulerStatus.Running) {
		s.schedulerMu.Unlock()
		return
	}
	s.schedulerStatus.Running = true
	s.schedulerStatus.NextRunAt = s.now().Add(cfg.Interval)
	s.schedulerMu.Unlock()

	go func() {
		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				s.stopScheduler()
				return
			case <-ticker.C:
				if cfg.MaxRuns > 0 && s.SchedulerStatus().ExecutedRuns >= cfg.MaxRuns {
					s.stopScheduler()
					return
				}
				_, _ = s.RunRound(ctx, domainbenchmark.RunInput{ScenarioFamily: cfg.ScenarioFamily})
				s.schedulerMu.Lock()
				s.schedulerStatus.NextRunAt = s.now().Add(cfg.Interval)
				s.schedulerMu.Unlock()
			}
		}
	}()
}

func (s *Service) RunScheduled(ctx context.Context, input domainbenchmark.SchedulerControlInput) ([]domainbenchmark.BenchmarkRound, error) {
	interval, err := resolveSchedulerInterval(s.schedulerConfig, input)
	if err != nil {
		return nil, err
	}
	scenarioFamily := resolveSchedulerScenarioFamily(s.schedulerConfig, input)
	maxRuns := resolveSchedulerMaxRuns(input)

	s.schedulerMu.Lock()
	s.schedulerStatus.Enabled = true
	s.schedulerStatus.Running = true
	s.schedulerStatus.Interval = interval.String()
	s.schedulerStatus.MaxRuns = maxRuns
	s.schedulerStatus.DryRun = input.DryRun
	s.schedulerStatus.ScenarioFamily = scenarioFamily
	s.schedulerMu.Unlock()

	items := make([]domainbenchmark.BenchmarkRound, 0, maxRuns)
	for i := 0; i < maxRuns; i++ {
		round, err := s.RunRound(ctx, domainbenchmark.RunInput{ScenarioFamily: scenarioFamily})
		if err != nil {
			s.stopScheduler()
			return nil, err
		}
		items = append(items, round)
		if i < maxRuns-1 {
			select {
			case <-ctx.Done():
				s.stopScheduler()
				return items, ctx.Err()
			case <-time.After(interval):
			}
		}
	}

	s.stopScheduler()
	return items, nil
}

func (s *Service) SchedulerStatus() domainbenchmark.SchedulerStatus {
	s.schedulerMu.RLock()
	defer s.schedulerMu.RUnlock()
	return s.schedulerStatus
}

func (s *Service) recordSchedulerExecution(roundID string, startedAt time.Time) {
	s.schedulerMu.Lock()
	defer s.schedulerMu.Unlock()
	s.schedulerStatus.ExecutedRuns++
	s.schedulerStatus.LastRoundID = roundID
	s.schedulerStatus.LastStartedAt = startedAt
	if s.schedulerConfig.Interval > 0 && s.schedulerStatus.Running {
		s.schedulerStatus.NextRunAt = startedAt.Add(s.schedulerConfig.Interval)
	}
}

func schedulerRunnable(cfg SchedulerConfig, running bool) bool {
	return cfg.Enabled && cfg.Interval > 0 && cfg.MaxRuns != 0 && !running
}

func resolveSchedulerInterval(cfg SchedulerConfig, input domainbenchmark.SchedulerControlInput) (time.Duration, error) {
	interval := cfg.Interval
	if input.Interval != "" {
		parsed, err := time.ParseDuration(input.Interval)
		if err != nil {
			return 0, err
		}
		interval = parsed
	}
	if interval <= 0 {
		interval = time.Second
	}
	return interval, nil
}

func resolveSchedulerScenarioFamily(cfg SchedulerConfig, input domainbenchmark.SchedulerControlInput) string {
	if input.ScenarioFamily != "" {
		return input.ScenarioFamily
	}
	if cfg.ScenarioFamily != "" {
		return cfg.ScenarioFamily
	}
	return defaultSchedulerScenarioFamily
}

func resolveSchedulerMaxRuns(input domainbenchmark.SchedulerControlInput) int {
	if input.MaxRuns > 0 {
		return input.MaxRuns
	}
	return 1
}

func (s *Service) stopScheduler() {
	s.schedulerMu.Lock()
	defer s.schedulerMu.Unlock()
	s.schedulerStatus.Running = false
	s.schedulerStatus.NextRunAt = time.Time{}
}
