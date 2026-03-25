package benchmark

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
	registrar Registrar
	scenario  ScenarioExecutor
	detection DetectionEvaluator
	replay    ReplayReader
	store     RoundStore
	reporter  ReportWriter
	now       func() time.Time
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
	}
}

func (s *Service) RegisterRound(ctx context.Context, request domainbenchmark.RegistrationRequest) (domainbenchmark.RegistrationResult, error) {
	return s.registrar.RegisterRound(ctx, request)
}

func (s *Service) Status() map[string]any {
	status := s.registrar.Status()
	status["rounds"] = len(s.store.List())
	return status
}

func (s *Service) RunRound(ctx context.Context, input domainbenchmark.RunInput) (domainbenchmark.BenchmarkRound, error) {
	scenarioFamily := strings.TrimSpace(input.ScenarioFamily)
	if scenarioFamily == "" {
		scenarioFamily = "commerce"
	}
	if scenarioFamily != "commerce" {
		return domainbenchmark.BenchmarkRound{}, fmt.Errorf("scenario family %s is not supported in Phase 7", scenarioFamily)
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
	round.PromotionResults = promotions
	round.StableSet = stableSetResult(stableResults)
	round.LivingSet = livingSetResult(livingResults)
	round.Delta = detectionDelta(round.ScenarioResults, previousRound, hasPrevious)
	round.Summary = roundSummary(round, replayResults, replayPassCount, promotions)
	round.CompletedAt = s.now()
	round.RoundStatus = domainbenchmark.RoundStatusCompleted

	reports, err := s.reporter.Generate(round)
	if err != nil {
		return domainbenchmark.BenchmarkRound{}, err
	}
	round.Reports = reports
	round.ReportDir = reports.Directory
	s.store.Put(round)

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
	for _, item := range s.scenario.ListScenarios() {
		if item.PackID == "commerce-pack" {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *Service) livingVariants() []domainbenchmark.ChallengerVariant {
	variants := []domainbenchmark.ChallengerVariant{
		{
			ID:                     "variant-weakened-provenance",
			ScenarioID:             "commerce-challenger-weakened-provenance-purchase",
			Title:                  "Weakened provenance",
			Description:            "Delegated purchase keeps provenance attached but materially weakens its confidence.",
			ChangeSet:              []string{"weakened_provenance", "low_confidence_context"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-expired-mandate",
			ScenarioID:             "commerce-challenger-expired-mandate-purchase",
			Title:                  "Expired mandate",
			Description:            "Delegated purchase proceeds with expired authority.",
			ChangeSet:              []string{"expired_mandate"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusSuspicious,
			ExpectedRecommendation: detectionmodel.RecommendationReview,
		},
		{
			ID:                     "variant-approval-removed",
			ScenarioID:             "commerce-challenger-approval-removed-refund",
			Title:                  "Approval removed",
			Description:            "Agent-driven refund proceeds without approval evidence.",
			ChangeSet:              []string{"approval_removed", "agent_refund"},
			ExpectedMinimumStatus:  detectionmodel.DetectionStatusStepUpRequired,
			ExpectedRecommendation: detectionmodel.RecommendationStepUp,
		},
	}
	return variants
}

func scenarioResultForStable(item domainscenario.Scenario, execution scenariosvc.ExecutionResult, detection detectionmodel.DetectionResult) domainbenchmark.ScenarioResult {
	expected := expectedStableStatus(item.ID)
	result := baseScenarioResult(item.ID, "commerce", domainbenchmark.ScenarioSetStable, "", execution, detection)
	result.ExpectedMinimumStatus = expected
	result.Passed = meetsMinimumStatus(detection.Status, expected)
	if item.ID == "commerce-clean-agent-assisted-purchase" {
		result.Passed = detection.Status == detectionmodel.DetectionStatusClean
	}
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
	return domainbenchmark.ScenarioResult{
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
	}
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
	return items
}

func expectedStableStatus(scenarioID string) detectionmodel.DetectionStatus {
	if scenarioID == "commerce-clean-agent-assisted-purchase" {
		return detectionmodel.DetectionStatusClean
	}
	return detectionmodel.DetectionStatusSuspicious
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
