package detection

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/commerce"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	domainevents "clawbot-trust-lab/internal/domain/events"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	executionsvc "clawbot-trust-lab/internal/services/scenario"
)

type WorldReader interface {
	GetOrder(string) (commerce.Order, error)
	ListRefunds() []commerce.Refund
	ListEvents() []domainevents.Record
	ListTrustDecisions() []domaintrust.TrustDecision
	ListApprovals() []domaintrust.ApprovalRecord
	GetMandate(string) (domaintrust.Mandate, error)
	GetProvenance(string) (domaintrust.ProvenanceRecord, error)
}

type ExecutionReader interface {
	GetExecutionResult(string) (executionsvc.ExecutionResult, error)
	GetExecutionResultByOrderID(string) (executionsvc.ExecutionResult, error)
}

type ReplayReader interface {
	ListCases() []domainreplay.ReplayCase
}

type ResultStore interface {
	Put(detectionmodel.DetectionResult)
	List() []detectionmodel.DetectionResult
	Get(string) (detectionmodel.DetectionResult, error)
	Summary() detectionmodel.DetectionRunSummary
}

type Service struct {
	world  WorldReader
	exec   ExecutionReader
	replay ReplayReader
	memory memory.Client
	store  ResultStore
	now    func() time.Time
}

const (
	ruleMissingMandateDelegatedAction = "missing_mandate_delegated_action"
	ruleExpiredInactiveMandate        = "expired_inactive_mandate"
	ruleMissingProvenanceSensitive    = "missing_provenance_sensitive_action"
	ruleRefundWeakAuthorization       = "refund_weak_authorization"
	ruleAgentRefundWithoutApproval    = "agent_refund_without_approval"
	ruleApprovalRemovedAuthorization  = "approval_removed_after_authorization"
	rulePriorStepUpDecision           = "prior_step_up_decision"
	ruleRepeatSuspiciousContext       = "repeat_suspicious_context"
	ruleMerchantScopeDriftDelegated   = "merchant_scope_drift_delegated_action"
	ruleHighValueDelegatedPurchase    = "high_value_delegated_purchase"
	ruleActorSwitchSensitiveAction    = "actor_switch_sensitive_action"
)

type ruleSpec struct {
	definition detectionmodel.RuleDefinition
	reason     string
}

var ruleCatalog = []ruleSpec{
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleMissingMandateDelegatedAction,
			Title:       "Missing or expired mandate on delegated action",
			Description: "Delegated commerce actions should not proceed without a valid mandate.",
			Severity:    20,
		},
		reason: "delegated action did not carry an active mandate",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleExpiredInactiveMandate,
			Title:       "Expired or inactive mandate on delegated action",
			Description: "Expired or inactive mandate coverage should force a stronger minimum posture for delegated actions.",
			Severity:    25,
		},
		reason: "delegated action relied on expired or inactive mandate coverage",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleMissingProvenanceSensitive,
			Title:       "Missing or weak provenance on sensitive action",
			Description: "Sensitive delegated actions should carry provenance evidence, and weak provenance should still increase concern.",
			Severity:    15,
		},
		reason: "sensitive delegated action lacked strong provenance evidence",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleRefundWeakAuthorization,
			Title:       "Refund with weak authorization",
			Description: "Refunds should not proceed with weak or expired authority.",
			Severity:    25,
		},
		reason: "refund was requested without strong authorization coverage",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleAgentRefundWithoutApproval,
			Title:       "Agent refund with no approval evidence",
			Description: "Agent-driven refund flows should carry explicit approval evidence.",
			Severity:    20,
		},
		reason: "agent-driven refund did not carry approval evidence",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleApprovalRemovedAuthorization,
			Title:       "Approval was removed before agent refund execution",
			Description: "Agent-driven refunds should step up when prior approval evidence was explicitly removed after authorization.",
			Severity:    25,
		},
		reason: "refund approval evidence was removed after authorization and before execution",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          rulePriorStepUpDecision,
			Title:       "Prior trust decision required step-up",
			Description: "Existing step-up decisions should increase downstream concern.",
			Severity:    10,
		},
		reason: "trust surface already recorded a step-up decision",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleRepeatSuspiciousContext,
			Title:       "Repeat suspicious refund behavior",
			Description: "Escalating repeat refund attempts should increase concern even before deeper replay history accumulates.",
			Severity:    15,
		},
		reason: "repeat refund history crossed the local escalation threshold or remained suspicious across replay and memory context",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleMerchantScopeDriftDelegated,
			Title:       "Merchant or category scope drift under delegated action",
			Description: "Delegated purchases drifting outside prior merchant or category scope should be reviewed.",
			Severity:    15,
		},
		reason: "delegated action moved outside prior merchant or category scope",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleHighValueDelegatedPurchase,
			Title:       "High-value delegated purchase above baseline",
			Description: "Delegated purchases that materially exceed the buyer's prior spend baseline should be reviewed.",
			Severity:    15,
		},
		reason: "delegated purchase materially exceeded the buyer's prior spend baseline",
	},
	{
		definition: detectionmodel.RuleDefinition{
			ID:          ruleActorSwitchSensitiveAction,
			Title:       "Sensitive action switched from human to agent",
			Description: "Sensitive flows that switch from human to agent without stronger controls should increase concern.",
			Severity:    10,
		},
		reason: "sensitive action switched to an agent actor without stronger controls",
	},
}

type EvaluateInput struct {
	ScenarioID string `json:"scenario_id,omitempty"`
	OrderID    string `json:"order_id,omitempty"`
}

func NewService(world WorldReader, exec ExecutionReader, replay ReplayReader, memoryClient memory.Client, store ResultStore) *Service {
	return &Service{
		world:  world,
		exec:   exec,
		replay: replay,
		memory: memoryClient,
		store:  store,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Evaluate(ctx context.Context, input EvaluateInput) (detectionmodel.DetectionResult, error) {
	execution, err := s.resolveExecution(input)
	if err != nil {
		return detectionmodel.DetectionResult{}, err
	}

	contextData := s.buildContext(ctx, execution)
	ruleHits := s.evaluateRules(contextData)
	score, status, grade, recommendation := deriveOutcome(ruleHits)

	reasonCodes := make([]string, 0, len(ruleHits))
	for _, hit := range ruleHits {
		reasonCodes = append(reasonCodes, hit.RuleID)
	}

	orderID := firstRef(execution.Entities.OrderRefs)
	refundID := firstRef(execution.Entities.RefundRefs)
	decisionRefs := make([]string, 0, len(execution.TrustDecisions))
	for _, decision := range execution.TrustDecisions {
		decisionRefs = append(decisionRefs, decision.ID)
	}

	result := detectionmodel.DetectionResult{
		ID:                detectionID(execution.Scenario.ID, orderID),
		ScenarioID:        execution.Scenario.ID,
		OrderID:           orderID,
		RefundID:          refundID,
		TrustDecisionRefs: decisionRefs,
		ReplayCaseRefs:    append([]string(nil), execution.ReplayCaseRefs...),
		Status:            status,
		Score:             score,
		Grade:             grade,
		TriggeredRules:    ruleHits,
		ReasonCodes:       reasonCodes,
		Recommendation:    recommendation,
		EvaluatedAt:       s.now(),
		Metadata: map[string]any{
			"context":      contextData,
			"tier_profile": contextData.TierProfile,
		},
	}
	s.store.Put(result)
	return result, nil
}

func (s *Service) ListResults() []detectionmodel.DetectionResult {
	return s.store.List()
}

func (s *Service) GetResult(id string) (detectionmodel.DetectionResult, error) {
	return s.store.Get(id)
}

func (s *Service) Summary() detectionmodel.DetectionRunSummary {
	return s.store.Summary()
}

func (s *Service) Rules() []detectionmodel.RuleDefinition {
	definitions := make([]detectionmodel.RuleDefinition, 0, len(ruleCatalog))
	for _, rule := range ruleCatalog {
		definitions = append(definitions, rule.definition)
	}
	return definitions
}

func (s *Service) resolveExecution(input EvaluateInput) (executionsvc.ExecutionResult, error) {
	if strings.TrimSpace(input.ScenarioID) != "" {
		return s.exec.GetExecutionResult(strings.TrimSpace(input.ScenarioID))
	}
	if strings.TrimSpace(input.OrderID) != "" {
		return s.exec.GetExecutionResultByOrderID(strings.TrimSpace(input.OrderID))
	}
	return executionsvc.ExecutionResult{}, fmt.Errorf("scenario_id or order_id is required")
}

func (s *Service) buildContext(ctx context.Context, execution executionsvc.ExecutionResult) detectionmodel.DetectionContext {
	orderID := firstRef(execution.Entities.OrderRefs)
	refundID := firstRef(execution.Entities.RefundRefs)

	order, mandate, provenance := s.loadOrderTrustContext(orderID)
	approvalPresent := s.hasApproval(orderID)
	relatedEvents, trustEventCount := relatedEventMetrics(s.world.ListEvents(), execution.Scenario.ID, append(execution.Entities.OrderRefs, execution.Entities.RefundRefs...))
	stepUp, reasonCount := trustDecisionMetrics(execution.TrustDecisions)
	replayHistoryCount := s.replayHistoryCount(execution.Scenario.ID)
	memoryContextPresent, memoryStatus := s.loadMemoryContext(ctx, execution.Scenario.ID)

	refundRequestedByAgent := refundID != "" && strings.Contains(firstRefundActor(s.world.ListRefunds(), refundID), "agent")
	delegatedActorPresent := refundRequestedByAgent || strings.Contains(order.SubmittedByActorID, "agent") || order.DelegationMode != "direct_human" || signalBool(execution.SignalContext, "delegated_indicator")
	repeatAttemptCount := signalInt(execution.SignalContext, "repeat_attempt_count")
	merchantScopeMatch := signalBoolDefault(execution.SignalContext, "merchant_scope_match", true)
	categoryScopeMatch := signalBoolDefault(execution.SignalContext, "category_scope_match", true)
	orderAmount := order.TotalAmount
	if signalAmount := signalInt64(execution.SignalContext, "order_amount"); signalAmount > 0 {
		orderAmount = signalAmount
	}
	buyerSpendBaseline := signalInt64(execution.SignalContext, "buyer_spend_baseline")
	highValueThreshold := signalInt64(execution.SignalContext, "high_value_threshold")
	weakProvenance := provenance.ID != "" && provenance.Confidence > 0 && provenance.Confidence < 0.5
	if signalBool(execution.SignalContext, "provenance_low_confidence") {
		weakProvenance = true
	}
	highValueDelegated := delegatedActorPresent && (highValueThreshold > 0 && orderAmount >= highValueThreshold || buyerSpendBaseline > 0 && orderAmount >= buyerSpendBaseline*2)

	tierA := signalStringSlice(execution.SignalContext, "tier_a_features")
	tierB := signalStringSlice(execution.SignalContext, "tier_b_features")
	tierC := signalStringSlice(execution.SignalContext, "tier_c_features")
	tierProfile := buildTierProfile(tierA, tierB, tierC, mandate, provenance, approvalPresent, signalBool(execution.SignalContext, "approval_removed"))

	features := map[string]bool{
		"delegated_actor_present":       delegatedActorPresent,
		"fully_delegated_action":        order.DelegationMode == "fully_delegated",
		"mandate_present":               order.MandateRef != "" && mandate.Status == "active",
		"mandate_missing":               order.MandateRef == "",
		"mandate_expired":               order.MandateRef != "" && mandate.Status != "" && mandate.Status != "active",
		"provenance_present":            order.ProvenanceRef != "" && provenance.ID != "",
		"provenance_missing":            order.ProvenanceRef == "" || provenance.ID == "",
		"weak_provenance":               weakProvenance,
		"approval_present":              approvalPresent,
		"approval_missing":              !approvalPresent,
		"approval_removed":              signalBool(execution.SignalContext, "approval_removed"),
		"refund_requested":              refundID != "",
		"refund_requested_by_agent":     refundRequestedByAgent,
		"refund_without_authorization":  refundID != "" && (!approvalPresent || (delegatedActorPresent && mandate.Status != "active")),
		"order_submitted_by_agent":      strings.Contains(order.SubmittedByActorID, "agent"),
		"trust_decision_step_up":        stepUp,
		"replay_history_present":        replayHistoryCount > 0,
		"memory_context_present":        memoryContextPresent,
		"repeat_attempt_escalation":     repeatAttemptCount >= 3,
		"merchant_scope_drift":          !merchantScopeMatch || !categoryScopeMatch,
		"high_value_delegated_purchase": highValueDelegated,
		"actor_switch_to_agent":         signalBool(execution.SignalContext, "actor_switch_to_agent"),
	}

	return detectionmodel.DetectionContext{
		ScenarioID:        execution.Scenario.ID,
		OrderID:           orderID,
		RefundID:          refundID,
		TrustDecisionRefs: trustDecisionRefs(execution.TrustDecisions),
		ReplayCaseRefs:    append([]string(nil), execution.ReplayCaseRefs...),
		Features:          features,
		Signals: map[string]any{
			"repeat_attempt_count":   repeatAttemptCount,
			"merchant_scope_match":   merchantScopeMatch,
			"category_scope_match":   categoryScopeMatch,
			"historical_refund_rate": signalInt(execution.SignalContext, "historical_refund_rate"),
			"buyer_spend_baseline":   buyerSpendBaseline,
			"order_amount":           orderAmount,
			"tier_a_features":        tierA,
			"tier_b_features":        tierB,
			"tier_c_features":        tierC,
		},
		EventCount:               len(relatedEvents),
		TrustEventCount:          trustEventCount,
		TrustDecisionReasonCount: reasonCount,
		ReplayHistoryCount:       replayHistoryCount,
		MemoryContextPresent:     memoryContextPresent,
		MemoryStatus:             memoryStatus,
		TierProfile:              tierProfile,
	}
}

func (s *Service) evaluateRules(contextData detectionmodel.DetectionContext) []detectionmodel.RuleHit {
	hits := []detectionmodel.RuleHit{}
	f := contextData.Features

	if f["delegated_actor_present"] && (f["mandate_missing"] || f["mandate_expired"]) {
		hits = append(hits, buildRuleHit(ruleMissingMandateDelegatedAction, map[string]any{"scenario_id": contextData.ScenarioID}))
	}
	if !f["refund_requested"] && f["delegated_actor_present"] && f["mandate_expired"] {
		hits = append(hits, buildRuleHit(ruleExpiredInactiveMandate, map[string]any{"scenario_id": contextData.ScenarioID}))
	}
	if (f["delegated_actor_present"] || f["order_submitted_by_agent"] || f["refund_requested_by_agent"] || f["fully_delegated_action"]) && (f["provenance_missing"] || f["weak_provenance"]) {
		hits = append(hits, buildRuleHit(ruleMissingProvenanceSensitive, map[string]any{"scenario_id": contextData.ScenarioID, "tier_c_used": contextData.TierProfile.TierCUsed}))
	}
	if f["refund_requested"] && f["refund_without_authorization"] {
		hits = append(hits, buildRuleHit(ruleRefundWeakAuthorization, map[string]any{"refund_id": contextData.RefundID}))
	}
	if f["refund_requested_by_agent"] && (f["approval_missing"] || f["approval_removed"]) {
		hits = append(hits, buildRuleHit(ruleAgentRefundWithoutApproval, map[string]any{"refund_id": contextData.RefundID}))
	}
	if f["refund_requested_by_agent"] && f["approval_removed"] {
		hits = append(hits, buildRuleHit(ruleApprovalRemovedAuthorization, map[string]any{"refund_id": contextData.RefundID}))
	}
	if f["trust_decision_step_up"] {
		hits = append(hits, buildRuleHit(rulePriorStepUpDecision, map[string]any{"trust_decision_refs": contextData.TrustDecisionRefs}))
	}
	if f["repeat_attempt_escalation"] || ((f["trust_decision_step_up"] || f["refund_requested_by_agent"]) && contextData.ReplayHistoryCount > 1 && contextData.MemoryContextPresent) {
		hits = append(hits, buildRuleHit(ruleRepeatSuspiciousContext, map[string]any{
			"replay_history_count": contextData.ReplayHistoryCount,
			"signals":              contextData.Signals,
		}))
	}
	if f["delegated_actor_present"] && f["merchant_scope_drift"] {
		hits = append(hits, buildRuleHit(ruleMerchantScopeDriftDelegated, map[string]any{"signals": contextData.Signals}))
	}
	if f["high_value_delegated_purchase"] {
		hits = append(hits, buildRuleHit(ruleHighValueDelegatedPurchase, map[string]any{"signals": contextData.Signals}))
	}
	if f["actor_switch_to_agent"] {
		hits = append(hits, buildRuleHit(ruleActorSwitchSensitiveAction, map[string]any{"scenario_id": contextData.ScenarioID}))
	}

	sort.Slice(hits, func(i, j int) bool { return hits[i].RuleID < hits[j].RuleID })
	return hits
}

func (s *Service) loadOrderTrustContext(orderID string) (commerce.Order, domaintrust.Mandate, domaintrust.ProvenanceRecord) {
	order, _ := s.world.GetOrder(orderID)
	var mandate domaintrust.Mandate
	if order.MandateRef != "" {
		mandate, _ = s.world.GetMandate(order.MandateRef)
	}
	var provenance domaintrust.ProvenanceRecord
	if order.ProvenanceRef != "" {
		provenance, _ = s.world.GetProvenance(order.ProvenanceRef)
	}
	return order, mandate, provenance
}

func (s *Service) hasApproval(orderID string) bool {
	for _, approval := range s.world.ListApprovals() {
		if approval.OrderID == orderID && strings.TrimSpace(approval.Outcome) != "" && !strings.EqualFold(approval.Outcome, "missing") {
			return true
		}
	}
	return false
}

func relatedEventMetrics(events []domainevents.Record, scenarioID string, entityIDs []string) ([]domainevents.Record, int) {
	related := filterEvents(events, scenarioID, entityIDs)
	trustEventCount := 0
	for _, event := range related {
		if event.Category == domainevents.EventCategoryTrust {
			trustEventCount++
		}
	}
	return related, trustEventCount
}

func trustDecisionMetrics(decisions []domaintrust.TrustDecision) (bool, int) {
	stepUp := false
	reasonCount := 0
	for _, decision := range decisions {
		reasonCount += len(decision.ReasonCodes)
		if decision.StepUpRequired {
			stepUp = true
		}
	}
	return stepUp, reasonCount
}

func (s *Service) replayHistoryCount(scenarioID string) int {
	count := 0
	for _, item := range s.replay.ListCases() {
		if item.ScenarioID == scenarioID {
			count++
		}
	}
	return count
}

func (s *Service) loadMemoryContext(ctx context.Context, scenarioID string) (bool, string) {
	response, err := s.memory.LoadScenarioContext(ctx, memory.LoadScenarioContextRequest{ScenarioID: scenarioID})
	if err != nil {
		return false, "degraded"
	}
	return contextRecordCount(response.Context) > 0, "ok"
}

func contextRecordCount(context map[string]any) int {
	switch count := context["record_count"].(type) {
	case int:
		return count
	case int64:
		return int(count)
	case float64:
		return int(count)
	default:
		return 0
	}
}

func buildTierProfile(tierA, tierB, tierC []string, mandate domaintrust.Mandate, provenance domaintrust.ProvenanceRecord, approvalPresent, approvalRemoved bool) detectionmodel.TierProfile {
	tierCUsed := len(tierC) > 0 && (mandate.ID != "" || provenance.ID != "" || approvalPresent || approvalRemoved)
	return detectionmodel.TierProfile{
		TierAAvailable: len(tierA) > 0,
		TierBAvailable: len(tierB) > 0,
		TierCAvailable: len(tierC) > 0,
		TierCUsed:      tierCUsed,
		TierANotes:     append([]string(nil), tierA...),
		TierBNotes:     append([]string(nil), tierB...),
		TierCNotes:     append([]string(nil), tierC...),
	}
}

func buildRuleHit(id string, metadata map[string]any) detectionmodel.RuleHit {
	spec := ruleByID(id)
	return detectionmodel.RuleHit{
		RuleID:   spec.definition.ID,
		Title:    spec.definition.Title,
		Severity: spec.definition.Severity,
		Reason:   spec.reason,
		Metadata: metadata,
	}
}

func ruleByID(id string) ruleSpec {
	for _, rule := range ruleCatalog {
		if rule.definition.ID == id {
			return rule
		}
	}
	return ruleSpec{definition: detectionmodel.RuleDefinition{ID: id}}
}

func deriveOutcome(hits []detectionmodel.RuleHit) (int, detectionmodel.DetectionStatus, detectionmodel.RiskGrade, detectionmodel.Recommendation) {
	score := 0
	for _, hit := range hits {
		score += hit.Severity
	}

	switch {
	case score >= 100:
		return score, detectionmodel.DetectionStatusBlocked, detectionmodel.RiskGradeCritical, detectionmodel.RecommendationBlock
	case score >= 40:
		return score, detectionmodel.DetectionStatusStepUpRequired, detectionmodel.RiskGradeHigh, detectionmodel.RecommendationStepUp
	case score >= 15:
		return score, detectionmodel.DetectionStatusSuspicious, detectionmodel.RiskGradeModerate, detectionmodel.RecommendationReview
	default:
		return score, detectionmodel.DetectionStatusClean, detectionmodel.RiskGradeLow, detectionmodel.RecommendationAllow
	}
}

func detectionID(scenarioID string, orderID string) string {
	if orderID != "" {
		return "det-" + orderID
	}
	return "det-" + scenarioID
}

func firstRef(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[0]
}

func signalBool(data map[string]any, key string) bool {
	value, ok := data[key]
	if !ok {
		return false
	}
	parsed, ok := value.(bool)
	return ok && parsed
}

func signalBoolDefault(data map[string]any, key string, fallback bool) bool {
	value, ok := data[key]
	if !ok {
		return fallback
	}
	parsed, ok := value.(bool)
	if !ok {
		return fallback
	}
	return parsed
}

func signalInt(data map[string]any, key string) int {
	value, ok := data[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	default:
		return 0
	}
}

func signalInt64(data map[string]any, key string) int64 {
	value, ok := data[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case float64:
		return int64(typed)
	default:
		return 0
	}
}

func signalStringSlice(data map[string]any, key string) []string {
	value, ok := data[key]
	if !ok {
		return nil
	}
	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...)
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				out = append(out, text)
			}
		}
		return out
	default:
		return nil
	}
}

func trustDecisionRefs(items []domaintrust.TrustDecision) []string {
	refs := make([]string, 0, len(items))
	for _, item := range items {
		refs = append(refs, item.ID)
	}
	return refs
}

func filterEvents(items []domainevents.Record, scenarioID string, entityIDs []string) []domainevents.Record {
	entitySet := map[string]struct{}{}
	for _, id := range entityIDs {
		entitySet[id] = struct{}{}
	}
	filtered := make([]domainevents.Record, 0)
	for _, item := range items {
		if item.ScenarioID != scenarioID {
			continue
		}
		if len(entitySet) > 0 {
			if _, ok := entitySet[item.EntityID]; !ok {
				continue
			}
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func firstRefundActor(refunds []commerce.Refund, refundID string) string {
	for _, item := range refunds {
		if item.ID == refundID {
			return item.RequestedByActorID
		}
	}
	return ""
}
