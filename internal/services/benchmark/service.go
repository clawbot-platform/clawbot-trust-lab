package benchmark

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	detectionsvc "clawbot-trust-lab/internal/services/detection"
	scenariosvc "clawbot-trust-lab/internal/services/scenario"
	"clawbot-trust-lab/internal/version"
)

type Registrar interface {
	RegisterRound(context.Context, domainbenchmark.RegistrationRequest) (domainbenchmark.RegistrationResult, error)
	Status() map[string]any
}

type ScenarioExecutor interface {
	ListScenarios() []domainscenario.Scenario
	Execute(context.Context, string) (scenariosvc.ExecutionResult, error)
}

type DetectionEvaluator interface {
	Evaluate(context.Context, detectionsvc.EvaluateInput) (detectionmodel.DetectionResult, error)
}

type ReplayReader interface {
	ListCases() []domainreplay.ReplayCase
}

type RoundStore interface {
	Put(domainbenchmark.BenchmarkRound)
	List() []domainbenchmark.BenchmarkRound
	Get(string) (domainbenchmark.BenchmarkRound, error)
	Latest() (domainbenchmark.BenchmarkRound, bool)
}

type ReportWriter interface {
	Generate(domainbenchmark.BenchmarkRound) (domainbenchmark.ReportIndex, error)
}

type Service struct {
	registrar       Registrar
	scenario        ScenarioExecutor
	detection       DetectionEvaluator
	replay          ReplayReader
	store           RoundStore
	reporter        ReportWriter
	now             func() time.Time
	schedulerMu     sync.RWMutex
	schedulerConfig SchedulerConfig
	schedulerStatus domainbenchmark.SchedulerStatus
}

func NewService(registrar Registrar, scenario ScenarioExecutor, detection DetectionEvaluator, replay ReplayReader, store RoundStore, reporter ReportWriter) *Service {
	return &Service{
		registrar: registrar,
		scenario:  scenario,
		detection: detection,
		replay:    replay,
		store:     store,
		reporter:  reporter,
		now:       func() time.Time { return time.Now().UTC() },
		schedulerStatus: domainbenchmark.SchedulerStatus{
			ScenarioFamily: "commerce",
			Interval:       "24h",
			MaxRuns:        7,
		},
	}
}

func (s *Service) RegisterRound(ctx context.Context, request domainbenchmark.RegistrationRequest) (domainbenchmark.RegistrationResult, error) {
	return s.registrar.RegisterRound(ctx, request)
}

func (s *Service) Status() map[string]any {
	status := s.registrar.Status()
	status["rounds"] = len(s.store.List())
	status["scheduler"] = s.SchedulerStatus()
	return status
}

func (s *Service) RunRound(ctx context.Context, input domainbenchmark.RunInput) (domainbenchmark.BenchmarkRound, error) {
	scenarioFamily := strings.TrimSpace(input.ScenarioFamily)
	if scenarioFamily == "" {
		scenarioFamily = "commerce"
	}
	if scenarioFamily != "commerce" {
		return domainbenchmark.BenchmarkRound{}, fmt.Errorf("scenario family %s is not supported in Phase 9", scenarioFamily)
	}

	startedAt := s.now()
	roundID := "round-" + startedAt.Format("20060102150405")
	previousRound, hasPrevious := s.store.Latest()

	round := domainbenchmark.BenchmarkRound{
		ID:              roundID,
		ScenarioFamily:  scenarioFamily,
		DetectorVersion: version.Current().Version,
		StartedAt:       startedAt,
		RoundStatus:     domainbenchmark.RoundStatusRunning,
	}

	var stableResults []domainbenchmark.ScenarioResult
	var livingResults []domainbenchmark.ScenarioResult
	var replayResults []domainbenchmark.ScenarioResult
	var promotions []domainbenchmark.PromotionDecision

	stableScenarios := s.stableScenarios()
	livingVariants := s.livingVariants()
	priorPromotions := promotedCases(previousRound)

	for _, item := range stableScenarios {
		execution, detection, err := s.executeAndEvaluate(ctx, item.ID)
		if err != nil {
			return domainbenchmark.BenchmarkRound{}, err
		}
		result := scenarioResultForStable(item, execution, detection)
		stableResults = append(stableResults, result)
		round.StableScenarioRefs = append(round.StableScenarioRefs, item.ID)
	}

	for _, variant := range livingVariants {
		execution, detection, err := s.executeAndEvaluate(ctx, variant.ScenarioID)
		if err != nil {
			return domainbenchmark.BenchmarkRound{}, err
		}
		result := scenarioResultForVariant(variant, execution, detection)
		promotion := decidePromotion(roundID, variant, result, previousRound)
		result.PromotedToReplay = promotion.Promoted
		livingResults = append(livingResults, result)
		round.ChallengerVariantRefs = append(round.ChallengerVariantRefs, variant.ID)
		if promotion.Promoted {
			promotions = append(promotions, promotion)
			if replayRef := firstRef(result.ReplayCaseRefs); replayRef != "" {
				round.ReplayCaseRefs = append(round.ReplayCaseRefs, replayRef)
			}
		}
	}

	replayPassCount := 0
	for _, promoted := range priorPromotions {
		execution, detection, err := s.executeAndEvaluate(ctx, promoted.ScenarioID)
		if err != nil {
			return domainbenchmark.BenchmarkRound{}, err
		}
		result := scenarioResultForReplay(promoted, execution, detection)
		if result.Passed {
			replayPassCount++
		} else {
			promotions = append(promotions, domainbenchmark.PromotionDecision{
				ID:                  promotionID(roundID, promoted.ScenarioID, "regression"),
				RoundID:             roundID,
				ScenarioID:          promoted.ScenarioID,
				ChallengerVariantID: promoted.ChallengerVariantID,
				PromotionReason:     domainbenchmark.PromotionReasonMeaningfulRegression,
				Rationale:           "Previously promoted replay case regressed below its expected detection floor.",
				DetectionResultRef:  result.DetectionResultRef,
				ReplayCaseRef:       firstRef(result.ReplayCaseRefs),
				ScenarioResultRef:   result.ID,
				Promoted:            true,
				CreatedAt:           s.now(),
			})
		}
		replayResults = append(replayResults, result)
	}

	round.ScenarioResults = append(round.ScenarioResults, stableResults...)
	round.ScenarioResults = append(round.ScenarioResults, livingResults...)
	round.ScenarioResults = append(round.ScenarioResults, replayResults...)
	round.PromotionResults = dedupePromotions(promotions)
	round.StableSet = stableSetResult(stableResults)
	round.LivingSet = livingSetResult(livingResults)
	round.Delta = detectionDelta(round.ScenarioResults, previousRound, hasPrevious)
	round.ReplayCaseRefs = uniqueStrings(round.ReplayCaseRefs)
	round.Recommendations = generateRecommendations(roundID, round.ScenarioResults, round.PromotionResults)
	round.Summary = roundSummary(round, replayResults, replayPassCount, round.PromotionResults)
	round.CompletedAt = s.now()
	round.RoundStatus = domainbenchmark.RoundStatusCompleted

	reports, err := s.reporter.Generate(round)
	if err != nil {
		return domainbenchmark.BenchmarkRound{}, err
	}
	round.Reports = reports
	round.ReportDir = reports.Directory
	s.store.Put(round)
	s.recordSchedulerExecution(round.ID, startedAt)

	return round, nil
}

func (s *Service) ListRounds() []domainbenchmark.BenchmarkRound {
	return s.store.List()
}

func (s *Service) GetRound(id string) (domainbenchmark.BenchmarkRound, error) {
	return s.store.Get(id)
}

func (s *Service) GetRoundSummary(id string) (domainbenchmark.RoundSummary, error) {
	round, err := s.store.Get(id)
	if err != nil {
		return domainbenchmark.RoundSummary{}, err
	}
	return round.Summary, nil
}

func (s *Service) GetRoundPromotions(id string) ([]domainbenchmark.PromotionDecision, error) {
	round, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return append([]domainbenchmark.PromotionDecision(nil), round.PromotionResults...), nil
}

func (s *Service) GetRoundDelta(id string) ([]domainbenchmark.DetectionDelta, error) {
	round, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return append([]domainbenchmark.DetectionDelta(nil), round.Delta...), nil
}

func (s *Service) GetRoundReports(id string) (domainbenchmark.ReportIndex, error) {
	round, err := s.store.Get(id)
	if err != nil {
		return domainbenchmark.ReportIndex{}, err
	}
	return round.Reports, nil
}

func (s *Service) ListRecommendations() []domainbenchmark.Recommendation {
	rounds := s.store.List()
	items := make([]domainbenchmark.Recommendation, 0)
	for _, round := range rounds {
		items = append(items, round.Recommendations...)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].LinkedRoundID == items[j].LinkedRoundID {
			leftPriority := recommendationPriorityRank(items[i].Priority)
			rightPriority := recommendationPriorityRank(items[j].Priority)
			if leftPriority == rightPriority {
				return items[i].ID > items[j].ID
			}
			return leftPriority > rightPriority
		}
		return items[i].LinkedRoundID > items[j].LinkedRoundID
	})
	return items
}

func DeriveRecommendations(round domainbenchmark.BenchmarkRound) []domainbenchmark.Recommendation {
	if len(round.Recommendations) > 0 {
		return append([]domainbenchmark.Recommendation(nil), round.Recommendations...)
	}
	return generateRecommendations(round.ID, round.ScenarioResults, round.PromotionResults)
}

func EnsureProductionBridgeSummary(round *domainbenchmark.BenchmarkRound) {
	if round == nil {
		return
	}

	round.Recommendations = DeriveRecommendations(*round)
	if round.Summary.RoundID == "" {
		round.Summary.RoundID = round.ID
	}
	if round.Summary.ScenarioFamily == "" {
		round.Summary.ScenarioFamily = round.ScenarioFamily
	}
	if round.ScenarioFamily == "" {
		round.ScenarioFamily = round.Summary.ScenarioFamily
	}
	if round.Summary.EvaluationMode == "" {
		round.Summary.EvaluationMode = "shadow"
	}
	if round.Summary.BlockingMode == "" {
		round.Summary.BlockingMode = "recommendation_only"
	}
	if round.Summary.ExistingControlNote == "" {
		round.Summary.ExistingControlNote = defaultExistingControlNote()
	}
	if round.Summary.RecommendedFollowUp == "" {
		round.Summary.RecommendedFollowUp = recommendedFollowUp(round.PromotionResults, round.Recommendations)
	}
	if round.Summary.Recommendations == 0 && len(round.Recommendations) > 0 {
		round.Summary.Recommendations = len(round.Recommendations)
	}
}

func (s *Service) GetRecommendation(id string) (domainbenchmark.Recommendation, error) {
	for _, item := range s.ListRecommendations() {
		if item.ID == id {
			return item, nil
		}
	}
	return domainbenchmark.Recommendation{}, fmt.Errorf("recommendation %s not found", id)
}

func (s *Service) LongRunSummary() domainbenchmark.LongRunSummary {
	rounds := s.store.List()
	summary := domainbenchmark.LongRunSummary{
		RoundsExecuted:       len(rounds),
		RecommendationCounts: map[domainbenchmark.RecommendationType]int{},
	}

	patternCounts := map[string]int{}
	for _, round := range rounds {
		summary.PromotionsOverTime = append(summary.PromotionsOverTime, domainbenchmark.MetricPoint{
			RoundID: round.ID,
			Value:   float64(round.Summary.PromotionCount),
		})
		summary.ReplayPassRateOverTime = append(summary.ReplayPassRateOverTime, domainbenchmark.MetricPoint{
			RoundID: round.ID,
			Value:   round.Summary.ReplayPassRate,
		})
		if round.Summary.RobustnessOutcome == domainbenchmark.RobustnessOutcomeNewBlindSpotDiscovered {
			summary.NewBlindSpots++
		}
		if round.Summary.RobustnessOutcome == domainbenchmark.RobustnessOutcomeRegressed {
			summary.RegressionsObserved++
		}
		for _, rec := range round.Recommendations {
			summary.RecommendationCounts[rec.Type]++
		}
		for _, item := range round.Summary.ImportantFindings {
			patternCounts[item]++
		}
	}
	for pattern, count := range patternCounts {
		if count > 1 {
			summary.TopRecurringPatterns = append(summary.TopRecurringPatterns, pattern)
		}
	}
	sort.Strings(summary.TopRecurringPatterns)
	return summary
}

func (s *Service) executeAndEvaluate(ctx context.Context, scenarioID string) (scenariosvc.ExecutionResult, detectionmodel.DetectionResult, error) {
	execution, err := s.scenario.Execute(ctx, scenarioID)
	if err != nil {
		return scenariosvc.ExecutionResult{}, detectionmodel.DetectionResult{}, err
	}
	detection, err := s.detection.Evaluate(ctx, detectionsvc.EvaluateInput{ScenarioID: scenarioID})
	if err != nil {
		return scenariosvc.ExecutionResult{}, detectionmodel.DetectionResult{}, err
	}
	return execution, detection, nil
}

func (s *Service) stableScenarios() []domainscenario.Scenario {
	items := make([]domainscenario.Scenario, 0)
	stableIDs := map[string]struct{}{
		"commerce-h1-direct-human-purchase":                   {},
		"commerce-h2-human-refund-valid-history":              {},
		"commerce-a1-agent-assisted-purchase-valid-controls":  {},
		"commerce-a2-fully-delegated-replenishment-purchase":  {},
		"commerce-a3-agent-assisted-refund-approval-evidence": {},
		"commerce-s1-refund-weak-authorization":               {},
		"commerce-s4-repeated-agent-refund-attempts":          {},
	}
	for _, item := range s.scenario.ListScenarios() {
		if _, ok := stableIDs[item.ID]; ok {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *Service) livingVariants() []domainbenchmark.ChallengerVariant {
	variants := []domainbenchmark.ChallengerVariant{
		{
			ID:                     "variant-v1-weakened-provenance",
			ScenarioID:             "commerce-v1-weakened-provenance",
			Title:                  "Weakened provenance",
			Description:            "Delegated purchase keeps provenance attached but materially weakens its confidence.",
			ChangeSet:              []string{"weakened_provenance", "low_confidence_context"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-v2-expired-mandate",
			ScenarioID:             "commerce-v2-expired-inactive-mandate",
			Title:                  "Expired mandate",
			Description:            "Delegated purchase proceeds with expired authority.",
			ChangeSet:              []string{"expired_mandate"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-v3-approval-removed",
			ScenarioID:             "commerce-v3-approval-removed",
			Title:                  "Approval removed",
			Description:            "Agent-driven refund proceeds without approval evidence.",
			ChangeSet:              []string{"approval_removed", "agent_refund"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationStepUp,
		},
		{
			ID:                     "variant-v4-actor-switch",
			ScenarioID:             "commerce-v4-actor-switch-human-to-agent",
			Title:                  "Actor switch from human to agent",
			Description:            "A sensitive refund path switches from human to agent without stronger controls.",
			ChangeSet:              []string{"actor_switch_to_agent"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-v5-repeat-attempt-escalation",
			ScenarioID:             "commerce-v5-repeat-attempt-escalation",
			Title:                  "Repeat attempt escalation",
			Description:            "A repeated refund attempt pattern escalates above the baseline threshold.",
			ChangeSet:              []string{"repeat_attempt_escalation"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationStepUp,
		},
		{
			ID:                     "variant-v6-merchant-scope-drift",
			ScenarioID:             "commerce-v6-merchant-scope-drift",
			Title:                  "Merchant scope drift",
			Description:            "Delegated purchase drifts into a new merchant or category scope.",
			ChangeSet:              []string{"merchant_scope_drift", "category_scope_drift"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-v7-high-value-delegated-purchase",
			ScenarioID:             "commerce-v7-high-value-delegated-purchase",
			Title:                  "High-value delegated purchase",
			Description:            "Delegated purchase materially exceeds the buyer's prior spend baseline.",
			ChangeSet:              []string{"high_value_delegated_purchase"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-s2-weak-provenance",
			ScenarioID:             "commerce-s2-delegated-purchase-weak-provenance",
			Title:                  "Suspicious delegated purchase with weak provenance",
			Description:            "Delegated purchase remains accepted even though provenance evidence is materially weak.",
			ChangeSet:              []string{"weak_provenance", "delegated_purchase"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-s3-approval-removed-after-authorization",
			ScenarioID:             "commerce-s3-approval-removed-after-authorization",
			Title:                  "Approval removed after authorization",
			Description:            "Refund keeps initial authority but approval disappears before action execution.",
			ChangeSet:              []string{"approval_removed", "refund_after_authorization"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationStepUp,
		},
		{
			ID:                     "variant-s5-scope-drift",
			ScenarioID:             "commerce-s5-merchant-scope-drift-delegated-action",
			Title:                  "Merchant or category scope drift",
			Description:            "Delegated purchase drifts outside prior merchant and category history.",
			ChangeSet:              []string{"merchant_scope_drift", "delegated_purchase"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
	}
	return variants
}

func scenarioResultForStable(item domainscenario.Scenario, execution scenariosvc.ExecutionResult, detection detectionmodel.DetectionResult) domainbenchmark.ScenarioResult {
	expected := expectedStableStatus(item.ID)
	result := baseScenarioResult(item.ID, "commerce", domainbenchmark.ScenarioSetStable, "", execution, detection)
	result.ExpectedMinimumStatus = expected
	result.Passed = meetsMinimumStatus(detection.Status, expected)
	if result.Passed {
		result.Notes = append(result.Notes, "stable baseline met its expected detector outcome")
	} else {
		result.Notes = append(result.Notes, "stable baseline deviated from its expected detector outcome")
	}
	return result
}

func scenarioResultForVariant(variant domainbenchmark.ChallengerVariant, execution scenariosvc.ExecutionResult, detection detectionmodel.DetectionResult) domainbenchmark.ScenarioResult {
	result := baseScenarioResult(variant.ScenarioID, "commerce", domainbenchmark.ScenarioSetLiving, variant.ID, execution, detection)
	result.ExpectedMinimumStatus = variant.ExpectedMinimumStatus
	result.Passed = meetsMinimumStatus(detection.Status, variant.ExpectedMinimumStatus)
	if result.Passed {
		result.Notes = append(result.Notes, "challenger variant met the expected minimum detector posture")
	} else {
		result.Notes = append(result.Notes, "challenger variant exposed a detector weakness")
	}
	return result
}

func scenarioResultForReplay(promoted domainbenchmark.PromotionDecision, execution scenariosvc.ExecutionResult, detection detectionmodel.DetectionResult) domainbenchmark.ScenarioResult {
	result := baseScenarioResult(promoted.ScenarioID, "commerce", domainbenchmark.ScenarioSetReplay, promoted.ChallengerVariantID, execution, detection)
	expected := expectedReplayStatus(promoted.PromotionReason)
	result.ExpectedMinimumStatus = expected
	result.Passed = meetsMinimumStatus(detection.Status, expected)
	if result.Passed {
		result.Notes = append(result.Notes, "replay regression case preserved the expected detector floor")
	} else {
		result.Notes = append(result.Notes, "replay regression case fell below the expected detector floor")
	}
	return result
}

func baseScenarioResult(scenarioID string, scenarioFamily string, setKind domainbenchmark.ScenarioSetKind, variantID string, execution scenariosvc.ExecutionResult, detection detectionmodel.DetectionResult) domainbenchmark.ScenarioResult {
	result := domainbenchmark.ScenarioResult{
		ID:                   scenarioResultID(setKind, scenarioID),
		ScenarioID:           scenarioID,
		ScenarioFamily:       scenarioFamily,
		SetKind:              setKind,
		ChallengerVariantID:  variantID,
		ExecutionStatus:      "completed",
		OrderRefs:            append([]string(nil), execution.Entities.OrderRefs...),
		RefundRefs:           append([]string(nil), execution.Entities.RefundRefs...),
		TrustDecisionRefs:    trustDecisionRefs(execution.TrustDecisions),
		ReplayCaseRefs:       append([]string(nil), execution.ReplayCaseRefs...),
		MemoryRecordRefs:     memoryRefs(execution.MemoryWrites),
		DetectionResultRef:   detection.ID,
		FinalDetectionStatus: detection.Status,
		FinalRecommendation:  detection.Recommendation,
		TriggeredRuleIDs:     append([]string(nil), detection.ReasonCodes...),
	}
	if contextData, ok := detection.Metadata["context"].(detectionmodel.DetectionContext); ok && contextData.TierProfile.TierCUsed {
		result.Notes = append(result.Notes, "tier_c_used")
	}
	return result
}

func stableSetResult(items []domainbenchmark.ScenarioResult) domainbenchmark.StableSetResult {
	result := domainbenchmark.StableSetResult{TotalCount: len(items)}
	for _, item := range items {
		result.ScenarioRefs = append(result.ScenarioRefs, item.ScenarioID)
		result.DetectionResultRefs = append(result.DetectionResultRefs, item.DetectionResultRef)
		if item.Passed {
			result.PassedCount++
		}
	}
	return result
}

func livingSetResult(items []domainbenchmark.ScenarioResult) domainbenchmark.LivingSetResult {
	result := domainbenchmark.LivingSetResult{TotalCount: len(items)}
	for _, item := range items {
		if item.ChallengerVariantID != "" {
			result.VariantRefs = append(result.VariantRefs, item.ChallengerVariantID)
		}
		result.DetectionResultRefs = append(result.DetectionResultRefs, item.DetectionResultRef)
		if item.Passed {
			result.CaughtCount++
		}
		if item.PromotedToReplay {
			result.PromotionCount++
		}
	}
	return result
}

func detectionDelta(current []domainbenchmark.ScenarioResult, previous domainbenchmark.BenchmarkRound, hasPrevious bool) []domainbenchmark.DetectionDelta {
	if !hasPrevious {
		return nil
	}

	previousByKey := map[string]domainbenchmark.ScenarioResult{}
	for _, item := range previous.ScenarioResults {
		previousByKey[deltaKey(item)] = item
	}

	out := make([]domainbenchmark.DetectionDelta, 0)
	for _, item := range current {
		prior, ok := previousByKey[deltaKey(item)]
		if !ok {
			continue
		}
		out = append(out, domainbenchmark.DetectionDelta{
			ScenarioID:            item.ScenarioID,
			SetKind:               item.SetKind,
			PreviousRoundID:       previous.ID,
			PreviousStatus:        prior.FinalDetectionStatus,
			CurrentStatus:         item.FinalDetectionStatus,
			ScoreDelta:            statusScore(item.FinalDetectionStatus) - statusScore(prior.FinalDetectionStatus),
			NewlyTriggeredRules:   diff(item.TriggeredRuleIDs, prior.TriggeredRuleIDs),
			ClearedRules:          diff(prior.TriggeredRuleIDs, item.TriggeredRuleIDs),
			RecommendationChanged: prior.FinalRecommendation != item.FinalRecommendation,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].ScenarioID == out[j].ScenarioID {
			return out[i].SetKind < out[j].SetKind
		}
		return out[i].ScenarioID < out[j].ScenarioID
	})
	return out
}

func roundSummary(round domainbenchmark.BenchmarkRound, replayResults []domainbenchmark.ScenarioResult, replayPassCount int, promotions []domainbenchmark.PromotionDecision) domainbenchmark.RoundSummary {
	replayPassRate := 1.0
	if len(replayResults) > 0 {
		replayPassRate = float64(replayPassCount) / float64(len(replayResults))
	}

	findings := make([]string, 0)
	for _, item := range promotions {
		findings = append(findings, fmt.Sprintf("%s promoted because %s.", item.ScenarioID, item.Rationale))
	}
	if len(replayResults) == 0 {
		findings = append(findings, "No previously promoted replay cases were available for regression retest.")
	} else if replayPassRate < 1 {
		findings = append(findings, fmt.Sprintf("Replay regression pass rate fell to %.2f.", replayPassRate))
	}

	outcome := domainbenchmark.RobustnessOutcomeMixed
	switch {
	case len(promotions) > 0:
		outcome = domainbenchmark.RobustnessOutcomeNewBlindSpotDiscovered
	case len(replayResults) > 0 && replayPassRate < 1:
		outcome = domainbenchmark.RobustnessOutcomeRegressed
	case round.StableSet.PassedCount == round.StableSet.TotalCount && round.LivingSet.CaughtCount == round.LivingSet.TotalCount:
		outcome = domainbenchmark.RobustnessOutcomeImproved
	}

	return domainbenchmark.RoundSummary{
		RoundID:             round.ID,
		ScenarioFamily:      round.ScenarioFamily,
		StableScenarioCount: len(round.StableScenarioRefs),
		ChallengerCount:     len(round.ChallengerVariantRefs),
		ReplayRetestCount:   len(replayResults),
		PromotionCount:      len(promotions),
		ReplayPassRate:      replayPassRate,
		RobustnessOutcome:   outcome,
		ImportantFindings:   findings,
		EvaluationMode:      "shadow",
		BlockingMode:        "recommendation_only",
		ExistingControlNote: defaultExistingControlNote(),
		RecommendedFollowUp: recommendedFollowUp(promotions, round.Recommendations),
		Recommendations:     len(round.Recommendations),
		TierCUsageCount:     tierCUsageCount(round.ScenarioResults),
	}
}

func defaultExistingControlNote() string {
	return "Run this harness as a sidecar evaluator beside the incumbent fraud stack, existing PSP controls, and current review queue so teams can compare recommendations before any production policy change."
}

func decidePromotion(roundID string, variant domainbenchmark.ChallengerVariant, result domainbenchmark.ScenarioResult, previous domainbenchmark.BenchmarkRound) domainbenchmark.PromotionDecision {
	actualRank := statusRank(result.FinalDetectionStatus)
	expectedRank := statusRank(variant.ExpectedMinimumStatus)
	reason := domainbenchmark.PromotionReason("")
	rationale := ""

	switch {
	case result.FinalDetectionStatus == detectionmodel.DetectionStatusClean && expectedRank > statusRank(detectionmodel.DetectionStatusClean):
		reason = domainbenchmark.PromotionReasonDetectorMiss
		rationale = "Suspicious challenger behavior evaluated as clean."
	case actualRank < expectedRank:
		reason = domainbenchmark.PromotionReasonSuspiciousBehaviorTooLow
		rationale = "Challenger behavior scored below its expected minimum detector posture."
	case len(result.TriggeredRuleIDs) == 0 && actualRank == expectedRank && isNovelVariant(variant.ID, previous):
		reason = domainbenchmark.PromotionReasonNewTrustGapPattern
		rationale = "The round surfaced a new trust-gap pattern that was not present in prior promotions."
	}

	promoted := reason != ""
	return domainbenchmark.PromotionDecision{
		ID:                  promotionID(roundID, variant.ScenarioID, string(reason)),
		RoundID:             roundID,
		ScenarioID:          variant.ScenarioID,
		ChallengerVariantID: variant.ID,
		PromotionReason:     reason,
		Rationale:           rationale,
		DetectionResultRef:  result.DetectionResultRef,
		ReplayCaseRef:       firstRef(result.ReplayCaseRefs),
		ScenarioResultRef:   result.ID,
		Promoted:            promoted,
		CreatedAt:           time.Now().UTC(),
	}
}

func promotedCases(round domainbenchmark.BenchmarkRound) []domainbenchmark.PromotionDecision {
	items := make([]domainbenchmark.PromotionDecision, 0)
	for _, item := range round.PromotionResults {
		if item.Promoted {
			items = append(items, item)
		}
	}
	return dedupePromotions(items)
}

func expectedStableStatus(scenarioID string) detectionmodel.DetectionStatus {
	switch scenarioID {
	case "commerce-h1-direct-human-purchase",
		"commerce-h2-human-refund-valid-history",
		"commerce-a1-agent-assisted-purchase-valid-controls",
		"commerce-a2-fully-delegated-replenishment-purchase",
		"commerce-a3-agent-assisted-refund-approval-evidence",
		"commerce-clean-agent-assisted-purchase":
		return detectionmodel.DetectionStatusClean
	case "commerce-s1-refund-weak-authorization", "commerce-suspicious-refund-attempt":
		return detectionmodel.DetectionStatusStepUpRequired
	case "commerce-s4-repeated-agent-refund-attempts":
		return detectionmodel.DetectionStatusSuspicious
	default:
		return detectionmodel.DetectionStatusSuspicious
	}
}

func expectedReplayStatus(reason domainbenchmark.PromotionReason) detectionmodel.DetectionStatus {
	if reason == domainbenchmark.PromotionReasonDetectorMiss {
		return detectionmodel.DetectionStatusSuspicious
	}
	if reason == domainbenchmark.PromotionReasonSuspiciousBehaviorTooLow || reason == domainbenchmark.PromotionReasonMeaningfulRegression {
		return detectionmodel.DetectionStatusStepUpRequired
	}
	return detectionmodel.DetectionStatusSuspicious
}

func meetsMinimumStatus(actual detectionmodel.DetectionStatus, expected detectionmodel.DetectionStatus) bool {
	return statusRank(actual) >= statusRank(expected)
}

func statusRank(status detectionmodel.DetectionStatus) int {
	switch status {
	case detectionmodel.DetectionStatusBlocked:
		return 4
	case detectionmodel.DetectionStatusStepUpRequired:
		return 3
	case detectionmodel.DetectionStatusSuspicious:
		return 2
	default:
		return 1
	}
}

func statusScore(status detectionmodel.DetectionStatus) int {
	return statusRank(status) * 10
}

func trustDecisionRefs(items []domaintrust.TrustDecision) []string {
	refs := make([]string, 0, len(items))
	for _, item := range items {
		refs = append(refs, item.ID)
	}
	return refs
}

func memoryRefs(items []scenariosvc.MemoryWriteOutcome) []string {
	refs := make([]string, 0, len(items))
	for _, item := range items {
		refs = append(refs, item.SourceID)
	}
	return refs
}

func scenarioResultID(setKind domainbenchmark.ScenarioSetKind, scenarioID string) string {
	return fmt.Sprintf("sr-%s-%s", setKind, scenarioID)
}

func promotionID(roundID string, scenarioID string, reason string) string {
	if reason == "" {
		reason = "none"
	}
	return fmt.Sprintf("promo-%s-%s-%s", roundID, scenarioID, reason)
}

func firstRef(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[0]
}

func deltaKey(item domainbenchmark.ScenarioResult) string {
	return string(item.SetKind) + ":" + item.ScenarioID
}

func diff(current []string, previous []string) []string {
	previousSet := map[string]struct{}{}
	for _, item := range previous {
		previousSet[item] = struct{}{}
	}
	out := make([]string, 0)
	for _, item := range current {
		if _, ok := previousSet[item]; !ok {
			out = append(out, item)
		}
	}
	sort.Strings(out)
	return out
}

func isNovelVariant(variantID string, previous domainbenchmark.BenchmarkRound) bool {
	for _, item := range previous.PromotionResults {
		if item.ChallengerVariantID == variantID {
			return false
		}
	}
	return true
}

func generateRecommendations(roundID string, results []domainbenchmark.ScenarioResult, promotions []domainbenchmark.PromotionDecision) []domainbenchmark.Recommendation {
	items := []domainbenchmark.Recommendation{
		{
			ID:                             "rec-" + roundID + "-shadow",
			Type:                           domainbenchmark.RecommendationTypeMonitorInShadowMode,
			Rationale:                      "This round is best used as a recommendation-only sidecar beside the incumbent fraud stack so the team can compare benchmark findings against existing review and decision outcomes.",
			Priority:                       domainbenchmark.RecommendationPriorityLow,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              stableAndLivingScenarioIDs(results),
			LinkedPromotionIDs:             promotionIDs(promotions),
			SuggestedAction:                "Keep the harness in shadow mode, compare its outputs with current queue and policy outcomes, and only trial control changes after replay confirms the improvement.",
			ExistingControlIntegrationNote: "Designed to run beside existing fraud rules, queueing, and PSP controls without blocking live traffic.",
		},
	}

	if len(promotions) > 0 {
		scenarioIDs := promotionScenarioIDs(promotions)
		items = append(items, domainbenchmark.Recommendation{
			ID:                             "rec-" + roundID + "-replay",
			Type:                           domainbenchmark.RecommendationTypeAddToReplayStableSet,
			Rationale:                      "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
			Priority:                       domainbenchmark.RecommendationPriorityHigh,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              scenarioIDs,
			LinkedPromotionIDs:             promotionIDs(promotions),
			SuggestedAction:                "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
			ExistingControlIntegrationNote: "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning.",
		})
	}

	if hasRule(results, "refund_weak_authorization") {
		scenarioIDs := scenarioIDsByRule(results, "refund_weak_authorization")
		items = append(items, domainbenchmark.Recommendation{
			ID:                             "rec-" + roundID + "-refund-review",
			Type:                           domainbenchmark.RecommendationTypeTightenRefundReviewRule,
			Rationale:                      "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
			Priority:                       domainbenchmark.RecommendationPriorityHigh,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              scenarioIDs,
			LinkedPromotionIDs:             linkedPromotionIDs(promotions, scenarioIDs),
			SupportingRuleIDs:              []string{"refund_weak_authorization"},
			SuggestedAction:                "Tighten refund review thresholds for weak or missing authority signals and compare the sidecar recommendation rate with current manual-review outcomes.",
			ExistingControlIntegrationNote: "Fits naturally beside existing refund review queues and PSP-side refund controls.",
		})
	}

	if hasRule(results, "agent_refund_without_approval") {
		scenarioIDs := scenarioIDsByRule(results, "agent_refund_without_approval")
		items = append(items, domainbenchmark.Recommendation{
			ID:                             "rec-" + roundID + "-delegated-refund-step-up",
			Type:                           domainbenchmark.RecommendationTypeRequireStepUpForDelegatedRefunds,
			Rationale:                      "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
			Priority:                       domainbenchmark.RecommendationPriorityHigh,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              scenarioIDs,
			LinkedPromotionIDs:             linkedPromotionIDs(promotions, scenarioIDs),
			SupportingRuleIDs:              []string{"agent_refund_without_approval"},
			SuggestedAction:                "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
			ExistingControlIntegrationNote: "Designed to augment incumbent delegated-refund controls rather than replace them.",
		})
	}

	if hasRule(results, "missing_provenance_sensitive_action") {
		scenarioIDs := scenarioIDsByRule(results, "missing_provenance_sensitive_action")
		items = append(items, domainbenchmark.Recommendation{
			ID:                             "rec-" + roundID + "-delegated-provenance",
			Type:                           domainbenchmark.RecommendationTypeRequireProvenanceForDelegatedBuys,
			Rationale:                      "Delegated purchase paths with weak or missing provenance should not be treated as equivalent to ordinary human commerce, especially when they drift into new behavior patterns.",
			Priority:                       domainbenchmark.RecommendationPriorityModerate,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              scenarioIDs,
			LinkedPromotionIDs:             linkedPromotionIDs(promotions, scenarioIDs),
			SupportingRuleIDs:              []string{"missing_provenance_sensitive_action"},
			SuggestedAction:                "Require provenance for delegated purchases or keep them in recommendation-only shadow review until the team is comfortable tightening incumbent purchase controls.",
			ExistingControlIntegrationNote: "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks.",
		})
	}

	if hasRule(results, "repeat_suspicious_context") {
		scenarioIDs := scenarioIDsByRule(results, "repeat_suspicious_context")
		items = append(items, domainbenchmark.Recommendation{
			ID:                             "rec-" + roundID + "-repeat-refund",
			Type:                           domainbenchmark.RecommendationTypeInvestigateRepeatRefundPattern,
			Rationale:                      "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
			Priority:                       domainbenchmark.RecommendationPriorityModerate,
			LinkedRoundID:                  roundID,
			LinkedScenarioIDs:              scenarioIDs,
			LinkedPromotionIDs:             linkedPromotionIDs(promotions, scenarioIDs),
			SupportingRuleIDs:              []string{"repeat_suspicious_context"},
			SuggestedAction:                "Investigate repeat refund patterns, compare them with incumbent case outcomes, and tune queueing logic in shadow mode before any blocking change.",
			ExistingControlIntegrationNote: "Best used as an investigative sidecar signal that feeds existing fraud-review workflows.",
		})
	}

	return items
}

func recommendedFollowUp(promotions []domainbenchmark.PromotionDecision, recommendations []domainbenchmark.Recommendation) string {
	switch {
	case len(promotions) > 0:
		return "Review the promoted cases with the fraud team, add the strongest ones to replay, and compare the sidecar recommendation output against incumbent decisions before changing production policy."
	case len(recommendations) > 1:
		return "Review the recommendation list with the fraud team and trial the suggested adjustments in shadow mode beside the incumbent stack."
	default:
		return "Continue running sidecar rounds to build replay depth, compare results against current controls, and watch recommendation stability over time."
	}
}

func tierCUsageCount(results []domainbenchmark.ScenarioResult) int {
	total := 0
	for _, item := range results {
		for _, note := range item.Notes {
			if strings.Contains(note, "tier_c_used") {
				total++
				break
			}
		}
	}
	return total
}

func hasRule(results []domainbenchmark.ScenarioResult, ruleID string) bool {
	return len(scenarioIDsByRule(results, ruleID)) > 0
}

func scenarioIDsByRule(results []domainbenchmark.ScenarioResult, ruleID string) []string {
	ids := make([]string, 0)
	for _, item := range results {
		for _, triggered := range item.TriggeredRuleIDs {
			if triggered == ruleID {
				ids = append(ids, item.ScenarioID)
				break
			}
		}
	}
	sort.Strings(ids)
	return ids
}

func promotionIDs(items []domainbenchmark.PromotionDecision) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return uniqueStrings(ids)
}

func promotionScenarioIDs(items []domainbenchmark.PromotionDecision) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ScenarioID)
	}
	return uniqueStrings(ids)
}

func stableAndLivingScenarioIDs(results []domainbenchmark.ScenarioResult) []string {
	ids := make([]string, 0, len(results))
	for _, item := range results {
		if item.SetKind == domainbenchmark.ScenarioSetStable || item.SetKind == domainbenchmark.ScenarioSetLiving {
			ids = append(ids, item.ScenarioID)
		}
	}
	return uniqueStrings(ids)
}

func dedupePromotions(items []domainbenchmark.PromotionDecision) []domainbenchmark.PromotionDecision {
	if len(items) == 0 {
		return nil
	}

	bestByKey := make(map[string]domainbenchmark.PromotionDecision, len(items))
	for _, item := range items {
		key := item.ScenarioID
		if item.ChallengerVariantID != "" {
			key += "|" + item.ChallengerVariantID
		}

		existing, ok := bestByKey[key]
		if !ok || promotionPriority(item) > promotionPriority(existing) || (promotionPriority(item) == promotionPriority(existing) && item.CreatedAt.After(existing.CreatedAt)) {
			bestByKey[key] = item
		}
	}

	out := make([]domainbenchmark.PromotionDecision, 0, len(bestByKey))
	for _, item := range bestByKey {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].ScenarioID == out[j].ScenarioID {
			if out[i].ChallengerVariantID == out[j].ChallengerVariantID {
				return out[i].ID < out[j].ID
			}
			return out[i].ChallengerVariantID < out[j].ChallengerVariantID
		}
		return out[i].ScenarioID < out[j].ScenarioID
	})
	return out
}

func promotionPriority(item domainbenchmark.PromotionDecision) int {
	switch item.PromotionReason {
	case domainbenchmark.PromotionReasonMeaningfulRegression:
		return 5
	case domainbenchmark.PromotionReasonDetectorMiss:
		return 4
	case domainbenchmark.PromotionReasonSuspiciousBehaviorTooLow:
		return 3
	case domainbenchmark.PromotionReasonNewTrustGapPattern:
		return 2
	case domainbenchmark.PromotionReasonNovelEvasiveVariation:
		return 1
	default:
		return 0
	}
}

func uniqueStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}

func linkedPromotionIDs(promotions []domainbenchmark.PromotionDecision, scenarioIDs []string) []string {
	if len(promotions) == 0 || len(scenarioIDs) == 0 {
		return nil
	}
	scenarioSet := make(map[string]struct{}, len(scenarioIDs))
	for _, scenarioID := range scenarioIDs {
		scenarioSet[scenarioID] = struct{}{}
	}
	ids := make([]string, 0)
	for _, item := range promotions {
		if _, ok := scenarioSet[item.ScenarioID]; ok {
			ids = append(ids, item.ID)
		}
	}
	return uniqueStrings(ids)
}

func recommendationPriorityRank(priority domainbenchmark.RecommendationPriority) int {
	switch priority {
	case domainbenchmark.RecommendationPriorityHigh:
		return 3
	case domainbenchmark.RecommendationPriorityModerate:
		return 2
	case domainbenchmark.RecommendationPriorityLow:
		return 1
	default:
		return 0
	}
}
