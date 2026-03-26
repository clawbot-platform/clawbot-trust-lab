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
	return []detectionmodel.RuleDefinition{
		{ID: "missing_mandate_delegated_action", Title: "Missing or expired mandate on delegated action", Description: "Delegated commerce actions should not proceed without a valid mandate.", Severity: 20},
		{ID: "missing_provenance_sensitive_action", Title: "Missing or weak provenance on sensitive action", Description: "Sensitive delegated actions should carry provenance evidence, and weak provenance should still increase concern.", Severity: 15},
		{ID: "refund_weak_authorization", Title: "Refund with weak authorization", Description: "Refunds should not proceed with weak or expired authority.", Severity: 25},
		{ID: "agent_refund_without_approval", Title: "Agent refund with no approval evidence", Description: "Agent-driven refund flows should carry explicit approval evidence.", Severity: 20},
		{ID: "prior_step_up_decision", Title: "Prior trust decision required step-up", Description: "Existing step-up decisions should increase downstream concern.", Severity: 10},
		{ID: "repeat_suspicious_context", Title: "Repeat suspicious refund behavior", Description: "Escalating repeat refund attempts should increase concern even before deeper replay history accumulates.", Severity: 15},
		{ID: "merchant_scope_drift_delegated_action", Title: "Merchant or category scope drift under delegated action", Description: "Delegated purchases drifting outside prior merchant or category scope should be reviewed.", Severity: 15},
		{ID: "high_value_delegated_purchase", Title: "High-value delegated purchase above baseline", Description: "Delegated purchases that materially exceed the buyer's prior spend baseline should be reviewed.", Severity: 15},
		{ID: "actor_switch_sensitive_action", Title: "Sensitive action switched from human to agent", Description: "Sensitive flows that switch from human to agent without stronger controls should increase concern.", Severity: 10},
	}
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

	order, _ := s.world.GetOrder(orderID)
	var mandate domaintrust.Mandate
	if order.MandateRef != "" {
		mandate, _ = s.world.GetMandate(order.MandateRef)
	}
	var provenance domaintrust.ProvenanceRecord
	if order.ProvenanceRef != "" {
		provenance, _ = s.world.GetProvenance(order.ProvenanceRef)
	}

	approvals := s.world.ListApprovals()
	approvalPresent := false
	for _, approval := range approvals {
		if approval.OrderID == orderID && strings.TrimSpace(approval.Outcome) != "" && !strings.EqualFold(approval.Outcome, "missing") {
			approvalPresent = true
			break
		}
	}

	events := s.world.ListEvents()
	relatedEvents := filterEvents(events, execution.Scenario.ID, append(execution.Entities.OrderRefs, execution.Entities.RefundRefs...))
	trustEventCount := 0
	for _, event := range relatedEvents {
		if event.Category == domainevents.EventCategoryTrust {
			trustEventCount++
		}
	}

	decisions := execution.TrustDecisions
	stepUp := false
	reasonCount := 0
	for _, decision := range decisions {
		reasonCount += len(decision.ReasonCodes)
		if decision.StepUpRequired {
			stepUp = true
		}
	}

	replayCases := s.replay.ListCases()
	replayHistoryCount := 0
	for _, item := range replayCases {
		if item.ScenarioID == execution.Scenario.ID {
			replayHistoryCount++
		}
	}

	memoryContextPresent := false
	memoryStatus := "degraded"
	if response, err := s.memory.LoadScenarioContext(ctx, memory.LoadScenarioContextRequest{ScenarioID: execution.Scenario.ID}); err == nil {
		if count, ok := response.Context["record_count"].(int); ok && count > 0 {
			memoryContextPresent = true
		}
		if count, ok := response.Context["record_count"].(float64); ok && count > 0 {
			memoryContextPresent = true
		}
		memoryStatus = "ok"
	}

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
	tierC, tierCUsed := signalStringSlice(execution.SignalContext, "tier_c_features"), false
	if len(tierC) > 0 && (mandate.ID != "" || provenance.ID != "" || approvalPresent || signalBool(execution.SignalContext, "approval_removed")) {
		tierCUsed = true
	}

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
		TrustDecisionRefs: trustDecisionRefs(decisions),
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
		TierProfile: detectionmodel.TierProfile{
			TierAAvailable: len(tierA) > 0,
			TierBAvailable: len(tierB) > 0,
			TierCAvailable: len(tierC) > 0,
			TierCUsed:      tierCUsed,
			TierANotes:     append([]string(nil), tierA...),
			TierBNotes:     append([]string(nil), tierB...),
			TierCNotes:     append([]string(nil), tierC...),
		},
	}
}

func (s *Service) evaluateRules(contextData detectionmodel.DetectionContext) []detectionmodel.RuleHit {
	hits := []detectionmodel.RuleHit{}
	f := contextData.Features

	if f["delegated_actor_present"] && (f["mandate_missing"] || f["mandate_expired"]) {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "missing_mandate_delegated_action", Title: "Missing or expired mandate on delegated action", Severity: 20, Reason: "delegated action did not carry an active mandate", Metadata: map[string]any{"scenario_id": contextData.ScenarioID}})
	}
	if (f["delegated_actor_present"] || f["order_submitted_by_agent"] || f["refund_requested_by_agent"] || f["fully_delegated_action"]) && (f["provenance_missing"] || f["weak_provenance"]) {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "missing_provenance_sensitive_action", Title: "Missing or weak provenance on sensitive action", Severity: 15, Reason: "sensitive delegated action lacked strong provenance evidence", Metadata: map[string]any{"scenario_id": contextData.ScenarioID, "tier_c_used": contextData.TierProfile.TierCUsed}})
	}
	if f["refund_requested"] && f["refund_without_authorization"] {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "refund_weak_authorization", Title: "Refund with weak authorization", Severity: 25, Reason: "refund was requested without strong authorization coverage", Metadata: map[string]any{"refund_id": contextData.RefundID}})
	}
	if f["refund_requested_by_agent"] && (f["approval_missing"] || f["approval_removed"]) {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "agent_refund_without_approval", Title: "Agent refund with no approval evidence", Severity: 20, Reason: "agent-driven refund did not carry approval evidence", Metadata: map[string]any{"refund_id": contextData.RefundID}})
	}
	if f["trust_decision_step_up"] {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "prior_step_up_decision", Title: "Prior trust decision required step-up", Severity: 10, Reason: "trust surface already recorded a step-up decision", Metadata: map[string]any{"trust_decision_refs": contextData.TrustDecisionRefs}})
	}
	if f["repeat_attempt_escalation"] || ((f["trust_decision_step_up"] || f["refund_requested_by_agent"]) && contextData.ReplayHistoryCount > 1 && contextData.MemoryContextPresent) {
		hits = append(hits, detectionmodel.RuleHit{
			RuleID:   "repeat_suspicious_context",
			Title:    "Repeat suspicious refund behavior",
			Severity: 15,
			Reason:   "repeat refund history crossed the local escalation threshold or remained suspicious across replay and memory context",
			Metadata: map[string]any{
				"replay_history_count": contextData.ReplayHistoryCount,
				"signals":              contextData.Signals,
			},
		})
	}
	if f["delegated_actor_present"] && f["merchant_scope_drift"] {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "merchant_scope_drift_delegated_action", Title: "Merchant or category scope drift under delegated action", Severity: 15, Reason: "delegated action moved outside prior merchant or category scope", Metadata: map[string]any{"signals": contextData.Signals}})
	}
	if f["high_value_delegated_purchase"] {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "high_value_delegated_purchase", Title: "High-value delegated purchase above baseline", Severity: 15, Reason: "delegated purchase materially exceeded the buyer's prior spend baseline", Metadata: map[string]any{"signals": contextData.Signals}})
	}
	if f["actor_switch_to_agent"] {
		hits = append(hits, detectionmodel.RuleHit{RuleID: "actor_switch_sensitive_action", Title: "Sensitive action switched from human to agent", Severity: 10, Reason: "sensitive action switched to an agent actor without stronger controls", Metadata: map[string]any{"scenario_id": contextData.ScenarioID}})
	}

	sort.Slice(hits, func(i, j int) bool { return hits[i].RuleID < hits[j].RuleID })
	return hits
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
