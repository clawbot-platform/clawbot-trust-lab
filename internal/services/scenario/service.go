package scenario

import (
	"context"
	"fmt"
	"strings"
	"time"

	"clawbot-trust-lab/internal/domain/actors"
	"clawbot-trust-lab/internal/domain/agents"
	"clawbot-trust-lab/internal/domain/commerce"
	domainevents "clawbot-trust-lab/internal/domain/events"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	commerceSvc "clawbot-trust-lab/internal/services/commerce"
	eventsvc "clawbot-trust-lab/internal/services/events"
	trustsvc "clawbot-trust-lab/internal/services/trust"
)

type ScenarioCatalog interface {
	ListScenarios() []domainscenario.Scenario
	GetScenario(string) (domainscenario.Scenario, error)
}

type TrustArtifactWriter interface {
	CreateArtifact(context.Context, domaintrust.CreateArtifactInput) (domaintrust.TrustArtifact, error)
}

type ReplayWriter interface {
	CreateCase(context.Context, domainreplay.CreateCaseInput) (domainreplay.ReplayCase, error)
}

type ExecutionResult struct {
	Scenario       domainscenario.Scenario     `json:"scenario"`
	Entities       EntityRefs                  `json:"entities"`
	TrustDecisions []domaintrust.TrustDecision `json:"trust_decisions"`
	ReplayCaseRefs []string                    `json:"replay_case_refs"`
	MemoryWrites   []MemoryWriteOutcome        `json:"memory_write_outcomes"`
	EventRefs      []string                    `json:"event_refs"`
	SignalContext  map[string]any              `json:"signal_context"`
}

type EntityRefs struct {
	BuyerRefs         []string `json:"buyer_refs"`
	MerchantRefs      []string `json:"merchant_refs"`
	ProductRefs       []string `json:"product_refs"`
	OrderRefs         []string `json:"order_refs"`
	PaymentRefs       []string `json:"payment_refs"`
	RefundRefs        []string `json:"refund_refs"`
	TrustArtifactRefs []string `json:"trust_artifact_refs"`
}

type MemoryWriteOutcome struct {
	Kind     string `json:"kind"`
	SourceID string `json:"source_id"`
	Status   string `json:"status"`
}

type Service struct {
	scenarios ScenarioCatalog
	commerce  *commerceSvc.Service
	events    *eventsvc.Service
	trust     *trustsvc.Service
	artifacts TrustArtifactWriter
	replay    ReplayWriter
	results   map[string]ExecutionResult
	orderRefs map[string]string
}

const (
	flowKindPurchase = "purchase"
	flowKindRefund   = "refund"
	actorKindHuman   = "human"
	actorKindAgent   = "agent"
	entityTypeOrder  = "order"
	entityTypeRefund = "refund"
)

func NewService(scenarios ScenarioCatalog, commerce *commerceSvc.Service, events *eventsvc.Service, trust *trustsvc.Service, artifacts TrustArtifactWriter, replay ReplayWriter) *Service {
	return &Service{
		scenarios: scenarios,
		commerce:  commerce,
		events:    events,
		trust:     trust,
		artifacts: artifacts,
		replay:    replay,
		results:   map[string]ExecutionResult{},
		orderRefs: map[string]string{},
	}
}

func (s *Service) ListScenarios() []domainscenario.Scenario {
	return s.scenarios.ListScenarios()
}

func (s *Service) Execute(ctx context.Context, scenarioID string) (ExecutionResult, error) {
	item, err := s.scenarios.GetScenario(scenarioID)
	if err != nil {
		return ExecutionResult{}, err
	}

	blueprint, err := blueprintForScenario(item)
	if err != nil {
		return ExecutionResult{}, err
	}

	var result ExecutionResult
	switch blueprint.FlowKind {
	case flowKindPurchase:
		result, err = s.executePurchase(ctx, item, blueprint)
	case flowKindRefund:
		result, err = s.executeRefund(ctx, item, blueprint)
	default:
		err = fmt.Errorf("scenario %s has unsupported flow kind %s", item.ID, blueprint.FlowKind)
	}
	if err != nil {
		return ExecutionResult{}, err
	}

	s.storeResult(result)
	return result, nil
}

func (s *Service) GetExecutionResult(scenarioID string) (ExecutionResult, error) {
	result, ok := s.results[scenarioID]
	if !ok {
		return ExecutionResult{}, fmt.Errorf("scenario %s has not been executed", scenarioID)
	}
	return result, nil
}

func (s *Service) GetExecutionResultByOrderID(orderID string) (ExecutionResult, error) {
	scenarioID, ok := s.orderRefs[orderID]
	if !ok {
		return ExecutionResult{}, fmt.Errorf("order %s has no recorded scenario execution", orderID)
	}
	return s.GetExecutionResult(scenarioID)
}

type scenarioBlueprint struct {
	FlowKind                string
	Merchant                commerce.Merchant
	Product                 commerce.Product
	OrderSubmittedBy        string
	RefundRequestedBy       string
	DelegationMode          actors.DelegationMode
	Mandate                 *mandateBlueprint
	Provenance              *provenanceBlueprint
	Approval                *approvalBlueprint
	OrderStatus             commerce.OrderStatus
	RefundStatus            commerce.RefundStatus
	DecisionOutcome         string
	DecisionReasonCodes     []string
	StepUpRequired          bool
	PaymentMethod           string
	OrderCreatedAt          time.Time
	PaymentAuthorizedAt     time.Time
	RefundRequestedAt       time.Time
	OutcomeSummary          string
	ReplayRecommendation    string
	ReplayReason            string
	SignalContext           map[string]any
	ApprovalActionType      string
	DecisionEntityType      string
	OrderReason             string
	TrustDecisionRecordedAt time.Time
}

type mandateBlueprint struct {
	Status         string
	AllowedActions []string
	SpendingLimit  int64
	ExpiresAt      time.Time
}

type provenanceBlueprint struct {
	SourceType string
	SourceRef  string
	Confidence float64
	CreatedAt  time.Time
}

type approvalBlueprint struct {
	Outcome    string
	ApproverID string
	CreatedAt  time.Time
}

func (s *Service) executePurchase(ctx context.Context, item domainscenario.Scenario, plan scenarioBlueprint) (ExecutionResult, error) {
	world := seedWorld()
	world.merchant = plan.Merchant
	world.product = plan.Product
	world.delegation = plan.DelegationMode
	s.seedParticipants(world)

	orderID := deterministicID("order", item.ID)
	paymentID := deterministicID("payment", item.ID)
	decisionID := deterministicID("decision", item.ID)

	mandateRef := ""
	if plan.Mandate != nil {
		mandate := s.trust.RecordMandate(domaintrust.Mandate{
			ID:              deterministicID("mandate", item.ID),
			PrincipalID:     world.principal.PrincipalID,
			DelegateActorID: world.agent.ID,
			AllowedActions:  append([]string(nil), plan.Mandate.AllowedActions...),
			SpendingLimit:   plan.Mandate.SpendingLimit,
			ExpiresAt:       plan.Mandate.ExpiresAt,
			Status:          plan.Mandate.Status,
		})
		mandateRef = mandate.ID
	}

	provenanceRef := ""
	if plan.Provenance != nil {
		provenance := s.trust.RecordProvenance(domaintrust.ProvenanceRecord{
			ID:          deterministicID("prov", item.ID),
			ActorID:     actorIDFromKind(world, plan.OrderSubmittedBy),
			PrincipalID: world.principal.PrincipalID,
			SourceType:  plan.Provenance.SourceType,
			SourceRef:   plan.Provenance.SourceRef,
			Confidence:  plan.Provenance.Confidence,
			CreatedAt:   plan.Provenance.CreatedAt,
		})
		provenanceRef = provenance.ID
	}

	order, err := s.commerce.CreateOrder(commerceSvc.CreateOrderInput{
		ID:                 orderID,
		BuyerID:            world.buyer.ID,
		MerchantID:         world.merchant.ID,
		ProductIDs:         []string{world.product.ID},
		SubmittedByActorID: actorIDFromKind(world, plan.OrderSubmittedBy),
		DelegationMode:     plan.DelegationMode,
		MandateRef:         mandateRef,
		ProvenanceRef:      provenanceRef,
		Status:             plan.OrderStatus,
		CreatedAt:          plan.OrderCreatedAt,
		UpdatedAt:          plan.PaymentAuthorizedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	payment := s.commerce.CreatePayment(commerceSvc.CreatePaymentInput{
		ID:           paymentID,
		OrderID:      order.ID,
		Amount:       order.TotalAmount,
		Currency:     order.Currency,
		Status:       commerce.PaymentStatusAuthorized,
		Method:       plan.PaymentMethod,
		AuthorizedAt: plan.PaymentAuthorizedAt,
	})

	decision, err := s.trust.RecordDecision(domaintrust.TrustDecision{
		ID:             decisionID,
		EntityType:     entityTypeOrder,
		EntityID:       order.ID,
		Outcome:        plan.DecisionOutcome,
		ReasonCodes:    append([]string(nil), plan.DecisionReasonCodes...),
		MandateRef:     mandateRef,
		ProvenanceRef:  provenanceRef,
		StepUpRequired: plan.StepUpRequired,
		RecordedAt:     plan.TrustDecisionRecordedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	eventRefs := make([]string, 0, 6)
	if plan.DelegationMode != actors.DelegationModeDirectHuman || plan.OrderSubmittedBy == actorKindAgent {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-order-submitted", item.ID),
			domainevents.TrustEventOrderSubmittedByAgent,
			"order",
			order.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.OrderCreatedAt,
			map[string]any{"delegation_mode": plan.DelegationMode},
		).ID)
	}
	if mandateRef != "" {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-mandate-checked", item.ID),
			domainevents.TrustEventMandateChecked,
			"mandate",
			mandateRef,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.OrderCreatedAt.Add(15*time.Second),
			map[string]any{"status": plan.Mandate.Status, "spending_limit": plan.Mandate.SpendingLimit},
		).ID)
	}
	if provenanceRef != "" {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-provenance-attached", item.ID),
			domainevents.TrustEventProvenanceAttached,
			"provenance_record",
			provenanceRef,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.Provenance.CreatedAt,
			map[string]any{"confidence": plan.Provenance.Confidence},
		).ID)
	}
	eventRefs = append(eventRefs,
		s.events.RecordTransaction(
			deterministicID("evt-order-created", item.ID),
			domainevents.TransactionEventOrderCreated,
			"order",
			order.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.OrderCreatedAt,
			map[string]any{
				"total_amount":       order.TotalAmount,
				"currency":           order.Currency,
				"merchant_category":  world.merchant.Category,
				"delegation_mode":    plan.DelegationMode,
				"submitted_by_actor": plan.OrderSubmittedBy,
			},
		).ID,
		s.events.RecordTransaction(
			deterministicID("evt-payment-authorized", item.ID),
			domainevents.TransactionEventPaymentAuthorized,
			"payment",
			payment.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.PaymentAuthorizedAt,
			map[string]any{"order_id": order.ID, "method": payment.Method},
		).ID,
		s.events.RecordTrust(
			deterministicID("evt-trust-decision", item.ID),
			domainevents.TrustEventTrustDecisionRecorded,
			"trust_decision",
			decision.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			decision.RecordedAt,
			map[string]any{"outcome": decision.Outcome, "step_up_required": decision.StepUpRequired},
		).ID,
	)

	return s.completeScenarioResult(ctx, item, world, plan, EntityRefs{
		BuyerRefs:    []string{world.buyer.ID},
		MerchantRefs: []string{world.merchant.ID},
		ProductRefs:  []string{world.product.ID},
		OrderRefs:    []string{order.ID},
		PaymentRefs:  []string{payment.ID},
	}, []domaintrust.TrustDecision{decision}, eventRefs)
}

func (s *Service) executeRefund(ctx context.Context, item domainscenario.Scenario, plan scenarioBlueprint) (ExecutionResult, error) {
	world := seedWorld()
	world.merchant = plan.Merchant
	world.product = plan.Product
	world.delegation = plan.DelegationMode
	s.seedParticipants(world)

	orderID := deterministicID("order", item.ID)
	paymentID := deterministicID("payment", item.ID)
	refundID := deterministicID("refund", item.ID)
	decisionID := deterministicID("decision", item.ID)

	mandateRef := ""
	if plan.Mandate != nil {
		mandate := s.trust.RecordMandate(domaintrust.Mandate{
			ID:              deterministicID("mandate", item.ID),
			PrincipalID:     world.principal.PrincipalID,
			DelegateActorID: world.agent.ID,
			AllowedActions:  append([]string(nil), plan.Mandate.AllowedActions...),
			SpendingLimit:   plan.Mandate.SpendingLimit,
			ExpiresAt:       plan.Mandate.ExpiresAt,
			Status:          plan.Mandate.Status,
		})
		mandateRef = mandate.ID
	}

	provenanceRef := ""
	if plan.Provenance != nil {
		provenance := s.trust.RecordProvenance(domaintrust.ProvenanceRecord{
			ID:          deterministicID("prov", item.ID),
			ActorID:     actorIDFromKind(world, plan.RefundRequestedBy),
			PrincipalID: world.principal.PrincipalID,
			SourceType:  plan.Provenance.SourceType,
			SourceRef:   plan.Provenance.SourceRef,
			Confidence:  plan.Provenance.Confidence,
			CreatedAt:   plan.Provenance.CreatedAt,
		})
		provenanceRef = provenance.ID
	}

	order, err := s.commerce.CreateOrder(commerceSvc.CreateOrderInput{
		ID:                 orderID,
		BuyerID:            world.buyer.ID,
		MerchantID:         world.merchant.ID,
		ProductIDs:         []string{world.product.ID},
		SubmittedByActorID: actorIDFromKind(world, plan.OrderSubmittedBy),
		DelegationMode:     plan.DelegationMode,
		MandateRef:         mandateRef,
		ProvenanceRef:      provenanceRef,
		Status:             plan.OrderStatus,
		CreatedAt:          plan.OrderCreatedAt,
		UpdatedAt:          plan.RefundRequestedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	payment := s.commerce.CreatePayment(commerceSvc.CreatePaymentInput{
		ID:           paymentID,
		OrderID:      order.ID,
		Amount:       order.TotalAmount,
		Currency:     order.Currency,
		Status:       commerce.PaymentStatusAuthorized,
		Method:       plan.PaymentMethod,
		AuthorizedAt: plan.PaymentAuthorizedAt,
	})

	if plan.Approval != nil {
		s.trust.RecordApproval(domaintrust.ApprovalRecord{
			ID:         deterministicID("approval", item.ID),
			OrderID:    order.ID,
			ActionType: plan.ApprovalActionType,
			ApproverID: plan.Approval.ApproverID,
			Outcome:    plan.Approval.Outcome,
			CreatedAt:  plan.Approval.CreatedAt,
		})
	}

	refund := s.commerce.CreateRefund(commerceSvc.CreateRefundInput{
		ID:                 refundID,
		OrderID:            order.ID,
		Amount:             order.TotalAmount,
		Status:             plan.RefundStatus,
		RequestedByActorID: actorIDFromKind(world, plan.RefundRequestedBy),
		Reason:             plan.OrderReason,
		CreatedAt:          plan.RefundRequestedAt,
	})

	decision, err := s.trust.RecordDecision(domaintrust.TrustDecision{
		ID:             decisionID,
		EntityType:     plan.DecisionEntityType,
		EntityID:       refund.ID,
		Outcome:        plan.DecisionOutcome,
		ReasonCodes:    append([]string(nil), plan.DecisionReasonCodes...),
		MandateRef:     mandateRef,
		ProvenanceRef:  provenanceRef,
		StepUpRequired: plan.StepUpRequired,
		RecordedAt:     plan.TrustDecisionRecordedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	eventRefs := []string{
		s.events.RecordTransaction(
			deterministicID("evt-order-created", item.ID),
			domainevents.TransactionEventOrderCreated,
			"order",
			order.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.OrderCreatedAt,
			map[string]any{"total_amount": order.TotalAmount, "currency": order.Currency},
		).ID,
		s.events.RecordTransaction(
			deterministicID("evt-payment-authorized", item.ID),
			domainevents.TransactionEventPaymentAuthorized,
			"payment",
			payment.ID,
			item.ID,
			actorIDFromKind(world, plan.OrderSubmittedBy),
			plan.PaymentAuthorizedAt,
			map[string]any{"order_id": order.ID},
		).ID,
		s.events.RecordTransaction(
			deterministicID("evt-refund-requested", item.ID),
			domainevents.TransactionEventRefundRequested,
			"refund",
			refund.ID,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			plan.RefundRequestedAt,
			map[string]any{"order_id": order.ID, "reason": refund.Reason},
		).ID,
	}

	if plan.DelegationMode != actors.DelegationModeDirectHuman || plan.RefundRequestedBy == "agent" {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-order-submitted", item.ID),
			domainevents.TrustEventOrderSubmittedByAgent,
			plan.DecisionEntityType,
			refund.ID,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			plan.RefundRequestedAt,
			map[string]any{"delegation_mode": plan.DelegationMode},
		).ID)
	}
	if mandateRef != "" {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-mandate-checked", item.ID),
			domainevents.TrustEventMandateChecked,
			"mandate",
			mandateRef,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			plan.RefundRequestedAt.Add(5*time.Second),
			map[string]any{"status": plan.Mandate.Status},
		).ID)
	}
	if provenanceRef != "" {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-provenance-attached", item.ID),
			domainevents.TrustEventProvenanceAttached,
			"provenance_record",
			provenanceRef,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			plan.Provenance.CreatedAt,
			map[string]any{"confidence": plan.Provenance.Confidence},
		).ID)
	}
	if plan.Approval != nil {
		eventRefs = append(eventRefs, s.events.RecordTrust(
			deterministicID("evt-approval-recorded", item.ID),
			domainevents.TrustEventApprovalRecorded,
			"approval_record",
			deterministicID("approval", item.ID),
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			plan.Approval.CreatedAt,
			map[string]any{"outcome": plan.Approval.Outcome},
		).ID)
	}
	eventRefs = append(eventRefs,
		s.events.RecordTrust(
			deterministicID("evt-trust-decision", item.ID),
			domainevents.TrustEventTrustDecisionRecorded,
			"trust_decision",
			decision.ID,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			decision.RecordedAt,
			map[string]any{"outcome": decision.Outcome, "step_up_required": decision.StepUpRequired},
		).ID,
		s.events.RecordTransaction(
			deterministicID("evt-refund-decision", item.ID),
			domainevents.TransactionEventRefundDecisionRecorded,
			"refund",
			refund.ID,
			item.ID,
			actorIDFromKind(world, plan.RefundRequestedBy),
			decision.RecordedAt.Add(5*time.Second),
			map[string]any{"status": refund.Status},
		).ID,
	)

	return s.completeScenarioResult(ctx, item, world, plan, EntityRefs{
		BuyerRefs:    []string{world.buyer.ID},
		MerchantRefs: []string{world.merchant.ID},
		ProductRefs:  []string{world.product.ID},
		OrderRefs:    []string{order.ID},
		PaymentRefs:  []string{payment.ID},
		RefundRefs:   []string{refund.ID},
	}, []domaintrust.TrustDecision{decision}, eventRefs)
}

func (s *Service) completeScenarioResult(ctx context.Context, item domainscenario.Scenario, world worldSeed, plan scenarioBlueprint, entities EntityRefs, decisions []domaintrust.TrustDecision, eventRefs []string) (ExecutionResult, error) {
	artifact, err := s.artifacts.CreateArtifact(ctx, domaintrust.CreateArtifactInput{
		ScenarioID: item.ID,
		ArtifactID: deterministicID("ta", item.ID),
		CreatedAt:  plan.TrustDecisionRecordedAt.Add(time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	replayCase, err := s.replay.CreateCase(ctx, domainreplay.CreateCaseInput{
		CaseID:                  deterministicID("rc", item.ID),
		ScenarioID:              item.ID,
		TrustArtifactRefs:       []string{artifact.ID},
		BenchmarkRoundRef:       "phase-9-red-queen",
		OutcomeSummary:          plan.OutcomeSummary,
		PromotionRecommendation: plan.ReplayRecommendation,
		PromotionReason:         plan.ReplayReason,
		RecordedAt:              plan.TrustDecisionRecordedAt.Add(2 * time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	entities.TrustArtifactRefs = []string{artifact.ID}

	return ExecutionResult{
		Scenario:       item,
		Entities:       entities,
		TrustDecisions: decisions,
		ReplayCaseRefs: []string{replayCase.ID},
		MemoryWrites: []MemoryWriteOutcome{
			{Kind: "trust_artifact", SourceID: artifact.ID, Status: "written"},
			{Kind: "replay_case", SourceID: replayCase.ID, Status: "written"},
		},
		EventRefs:     eventRefs,
		SignalContext: signalContext(item, world, plan),
	}, nil
}

func blueprintForScenario(item domainscenario.Scenario) (scenarioBlueprint, error) {
	const merchantID = "merchant-orbit-books"
	const humanID = "human-alex"
	switch item.ID {
	case "commerce-h1-direct-human-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:     "human",
			delegationMode:       actors.DelegationModeDirectHuman,
			orderCreatedAt:       time.Date(2026, 3, 25, 9, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 9, 2, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusAccepted,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"direct_human_purchase", "consistent_buyer_history"},
			paymentMethod:        "card_on_file",
			outcomeSummary:       "Direct human purchase completed with ordinary commerce signals and no delegated action.",
			replayRecommendation: "promote",
			replayReason:         "Human baseline remains useful for replay regression.",
			signals: map[string]any{
				"buyer_history_orders":   12,
				"historical_refund_rate": 0,
				"repeat_attempt_count":   0,
			},
		}), nil
	case "commerce-h2-human-refund-valid-history":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "human",
			refundRequestedBy:    "human",
			delegationMode:       actors.DelegationModeDirectHuman,
			orderCreatedAt:       time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 10, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 10, 8, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRequested,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"valid_refund_history", "human_requested_refund", "approval_present"},
			paymentMethod:        "card_on_file",
			reason:               "Human requested refund with valid order history",
			approval:             &approvalBlueprint{Outcome: "approved", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 10, 8, 15, 0, time.UTC)},
			outcomeSummary:       "Human refund request carried valid history and approval evidence.",
			replayRecommendation: "promote",
			replayReason:         "Human refund baseline should remain in replay coverage.",
			signals: map[string]any{
				"historical_refund_rate": 1,
				"repeat_attempt_count":   0,
				"approval_history":       1,
			},
		}), nil
	case "commerce-a1-agent-assisted-purchase-valid-controls", "commerce-clean-agent-assisted-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:     "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 11, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 11, 2, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusAccepted,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"active_mandate", "high_provenance_confidence", "in_policy_merchant_scope"},
			paymentMethod:        "delegated_card_on_file",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 11, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.96, CreatedAt: time.Date(2026, 3, 25, 11, 0, 20, 0, time.UTC)},
			outcomeSummary:       "Agent-assisted purchase completed with valid controls and strong provenance.",
			replayRecommendation: "promote",
			replayReason:         "Valid agent-assisted baseline should remain in stable replay.",
			signals: map[string]any{
				"buyer_history_orders":       16,
				"merchant_scope_match":       true,
				"category_scope_match":       true,
				"historical_refund_rate":     0,
				"repeat_attempt_count":       0,
				"buyer_spend_baseline":       4500,
				"delegated_indicator":        true,
				"tier_c_optional_improves":   true,
				"existing_control_alignment": "mandate_and_provenance",
			},
		}), nil
	case "commerce-a2-fully-delegated-replenishment-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:     "agent",
			delegationMode:       actors.DelegationModeFullyDelegated,
			orderCreatedAt:       time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 12, 2, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusAccepted,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"active_mandate", "recurring_replenishment_pattern", "merchant_scope_match"},
			paymentMethod:        "delegated_card_on_file",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "replenish_order"}, SpendingLimit: 18000, ExpiresAt: time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "task_list", SourceRef: item.ID, Confidence: 0.88, CreatedAt: time.Date(2026, 3, 25, 12, 0, 15, 0, time.UTC)},
			outcomeSummary:       "Fully delegated replenishment purchase remained inside known merchant and category scope.",
			replayRecommendation: "promote",
			replayReason:         "Replenishment baseline should remain part of stable replay coverage.",
			signals: map[string]any{
				"merchant_scope_match": true,
				"category_scope_match": true,
				"repeat_attempt_count": 0,
				"buyer_spend_baseline": 4200,
				"delegated_indicator":  true,
			},
		}), nil
	case "commerce-a3-agent-assisted-refund-approval-evidence":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "agent",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 13, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 13, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 13, 10, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRequested,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"active_mandate", "approval_present", "valid_refund_history"},
			paymentMethod:        "delegated_card_on_file",
			reason:               "Agent-assisted refund with explicit approval evidence",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "request_refund"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 13, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.83, CreatedAt: time.Date(2026, 3, 25, 13, 0, 30, 0, time.UTC)},
			approval:             &approvalBlueprint{Outcome: "approved", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 13, 10, 10, 0, time.UTC)},
			outcomeSummary:       "Agent-assisted refund completed with approval evidence and valid delegated controls.",
			replayRecommendation: "promote",
			replayReason:         "Approved delegated refund baseline should remain replayable.",
			signals: map[string]any{
				"historical_refund_rate": 1,
				"repeat_attempt_count":   0,
				"approval_history":       1,
				"delegated_indicator":    true,
			},
		}), nil
	case "commerce-s1-refund-weak-authorization", "commerce-suspicious-refund-attempt":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "human",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 14, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 14, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 14, 8, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRejected,
			decisionOutcome:      "step_up_required",
			decisionReasons:      []string{"expired_mandate", "weak_provenance", "missing_human_approval"},
			paymentMethod:        "card_on_file",
			reason:               "Agent attempted refund without strong authority",
			mandate:              &mandateBlueprint{Status: "expired", AllowedActions: []string{"submit_order"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 25, 13, 59, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.41, CreatedAt: time.Date(2026, 3, 25, 14, 7, 0, 0, time.UTC)},
			approval:             &approvalBlueprint{Outcome: "missing", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 14, 8, 10, 0, time.UTC)},
			outcomeSummary:       "Refund attempt required step-up because authority, provenance, and approval coverage were insufficient.",
			replayRecommendation: "hold",
			replayReason:         "Suspicious refund path should remain replayable for regression testing.",
			signals: map[string]any{
				"historical_refund_rate": 2,
				"repeat_attempt_count":   1,
				"approval_history":       0,
				"delegated_indicator":    true,
			},
		}), nil
	case "commerce-s4-repeated-agent-refund-attempts":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "agent",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 15, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 15, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 15, 6, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRejected,
			decisionOutcome:      "step_up_required",
			decisionReasons:      []string{"repeat_refund_attempts", "agent_refund", "prior_step_up_history"},
			paymentMethod:        "delegated_card_on_file",
			reason:               "Repeated agent refund attempts crossed the baseline threshold",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "request_refund"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "task_list", SourceRef: item.ID, Confidence: 0.72, CreatedAt: time.Date(2026, 3, 25, 15, 5, 20, 0, time.UTC)},
			approval:             &approvalBlueprint{Outcome: "approved", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 15, 5, 30, 0, time.UTC)},
			outcomeSummary:       "Repeated agent refund attempts remained explainably suspicious because history increased concern.",
			replayRecommendation: "hold",
			replayReason:         "Repeat refund pattern should stay in the stable suspicious set.",
			signals: map[string]any{
				"historical_refund_rate": 4,
				"repeat_attempt_count":   3,
				"approval_history":       1,
				"delegated_indicator":    true,
			},
		}), nil
	case "commerce-s2-delegated-purchase-weak-provenance", "commerce-v1-weakened-provenance", "commerce-challenger-weakened-provenance-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:     "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 16, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 16, 2, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusAccepted,
			decisionOutcome:      "accepted",
			decisionReasons:      []string{"active_mandate", "weak_provenance_context"},
			paymentMethod:        "delegated_card_on_file",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 16, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.19, CreatedAt: time.Date(2026, 3, 25, 16, 0, 20, 0, time.UTC)},
			outcomeSummary:       "Delegated purchase remained accepted even though provenance confidence was materially weak.",
			replayRecommendation: "candidate",
			replayReason:         "Weak provenance on delegated purchase remains a meaningful blind-spot probe.",
			signals: map[string]any{
				"merchant_scope_match":      true,
				"category_scope_match":      true,
				"buyer_spend_baseline":      4500,
				"repeat_attempt_count":      1,
				"delegated_indicator":       true,
				"provenance_low_confidence": true,
			},
		}), nil
	case "commerce-v2-expired-inactive-mandate", "commerce-challenger-expired-mandate-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:     "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 17, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 17, 2, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusAccepted,
			decisionOutcome:      "review_required",
			decisionReasons:      []string{"expired_mandate"},
			paymentMethod:        "delegated_card_on_file",
			mandate:              &mandateBlueprint{Status: "expired", AllowedActions: []string{"submit_order"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 25, 16, 55, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.92, CreatedAt: time.Date(2026, 3, 25, 17, 0, 15, 0, time.UTC)},
			outcomeSummary:       "Delegated purchase attempted with expired mandate and required review.",
			replayRecommendation: "review",
			replayReason:         "Expired mandate challenger should remain part of replay regression.",
			signals: map[string]any{
				"merchant_scope_match": true,
				"category_scope_match": true,
				"buyer_spend_baseline": 4500,
				"repeat_attempt_count": 1,
				"delegated_indicator":  true,
			},
		}), nil
	case "commerce-s3-approval-removed-after-authorization", "commerce-v3-approval-removed", "commerce-challenger-approval-removed-refund":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "agent",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 18, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 18, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 18, 7, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRejected,
			decisionOutcome:      "step_up_required",
			decisionReasons:      []string{"approval_removed", "agent_refund"},
			paymentMethod:        "delegated_card_on_file",
			reason:               "Refund requested after approval evidence was removed",
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "request_refund"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 18, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "conversation", SourceRef: item.ID, Confidence: 0.74, CreatedAt: time.Date(2026, 3, 25, 18, 0, 20, 0, time.UTC)},
			approval:             &approvalBlueprint{Outcome: "removed", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 18, 6, 30, 0, time.UTC)},
			outcomeSummary:       "Agent-driven refund lost approval evidence before execution and required step-up.",
			replayRecommendation: "hold",
			replayReason:         "Approval-removed refund remains useful for replay regression.",
			signals: map[string]any{
				"historical_refund_rate": 2,
				"repeat_attempt_count":   1,
				"approval_history":       0,
				"delegated_indicator":    true,
			},
		}), nil
	case "commerce-s5-merchant-scope-drift-delegated-action", "commerce-v6-merchant-scope-drift":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:    "agent",
			delegationMode:      actors.DelegationModeFullyDelegated,
			orderCreatedAt:      time.Date(2026, 3, 25, 19, 0, 0, 0, time.UTC),
			paymentAuthorizedAt: time.Date(2026, 3, 25, 19, 2, 0, 0, time.UTC),
			orderStatus:         commerce.OrderStatusAccepted,
			decisionOutcome:     "review_required",
			decisionReasons:     []string{"merchant_scope_drift", "category_scope_drift"},
			paymentMethod:       "delegated_card_on_file",
			merchant: commerce.Merchant{
				ID:       "merchant-horizon-electronics",
				Name:     "Horizon Electronics",
				Category: "electronics",
				Tags:     []string{"scope-drift", "challenger"},
			},
			product: commerce.Product{
				ID:         "product-horizon-tablet",
				MerchantID: "merchant-horizon-electronics",
				Name:       "Horizon Tablet",
				Amount:     8900,
				Currency:   "USD",
				Category:   "electronics",
				Tags:       []string{"scope-drift", "challenger"},
			},
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "replenish_order"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 19, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "task_list", SourceRef: item.ID, Confidence: 0.81, CreatedAt: time.Date(2026, 3, 25, 19, 0, 20, 0, time.UTC)},
			outcomeSummary:       "Delegated purchase drifted into a new merchant and category scope and required review.",
			replayRecommendation: "candidate",
			replayReason:         "Merchant scope drift remains a strong blind-spot probe for sidecar evaluation.",
			signals: map[string]any{
				"merchant_scope_match": false,
				"category_scope_match": false,
				"buyer_spend_baseline": 4200,
				"repeat_attempt_count": 1,
				"delegated_indicator":  true,
			},
		}), nil
	case "commerce-v4-actor-switch-human-to-agent":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "human",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 20, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 20, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 20, 5, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRejected,
			decisionOutcome:      "review_required",
			decisionReasons:      []string{"actor_switch_to_agent", "approval_missing"},
			paymentMethod:        "card_on_file",
			reason:               "Refund flow switched from human to agent without strengthening controls",
			approval:             &approvalBlueprint{Outcome: "missing", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 20, 5, 10, 0, time.UTC)},
			outcomeSummary:       "Refund path switched from human to agent without stronger approval controls.",
			replayRecommendation: "candidate",
			replayReason:         "Actor-switch challenger probes low-telemetry delegated action gaps.",
			signals: map[string]any{
				"historical_refund_rate": 1,
				"repeat_attempt_count":   1,
				"approval_history":       0,
				"delegated_indicator":    true,
				"actor_switch_to_agent":  true,
			},
		}), nil
	case "commerce-v5-repeat-attempt-escalation":
		return refundBlueprint(refundProfile{
			orderSubmittedBy:     "agent",
			refundRequestedBy:    "agent",
			delegationMode:       actors.DelegationModeAgentAssisted,
			orderCreatedAt:       time.Date(2026, 3, 25, 21, 0, 0, 0, time.UTC),
			paymentAuthorizedAt:  time.Date(2026, 3, 25, 21, 1, 0, 0, time.UTC),
			refundRequestedAt:    time.Date(2026, 3, 25, 21, 4, 0, 0, time.UTC),
			orderStatus:          commerce.OrderStatusRefundReview,
			refundStatus:         commerce.RefundStatusRejected,
			decisionOutcome:      "step_up_required",
			decisionReasons:      []string{"repeat_refund_attempts", "agent_refund", "approval_missing"},
			paymentMethod:        "delegated_card_on_file",
			reason:               "Repeat refund attempts escalated above the baseline threshold",
			approval:             &approvalBlueprint{Outcome: "missing", ApproverID: humanID, CreatedAt: time.Date(2026, 3, 25, 21, 4, 10, 0, time.UTC)},
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order", "request_refund"}, SpendingLimit: 15000, ExpiresAt: time.Date(2026, 3, 26, 21, 0, 0, 0, time.UTC)},
			outcomeSummary:       "Repeat refund attempts escalated concern even before any exotic agentic overlay fields were needed.",
			replayRecommendation: "candidate",
			replayReason:         "Repeat attempt escalation is a production-adjacent sidecar recommendation case.",
			signals: map[string]any{
				"historical_refund_rate": 5,
				"repeat_attempt_count":   4,
				"approval_history":       0,
				"delegated_indicator":    true,
			},
		}), nil
	case "commerce-v7-high-value-delegated-purchase":
		return purchaseBlueprint(purchaseProfile{
			orderSubmittedBy:    "agent",
			delegationMode:      actors.DelegationModeFullyDelegated,
			orderCreatedAt:      time.Date(2026, 3, 25, 22, 0, 0, 0, time.UTC),
			paymentAuthorizedAt: time.Date(2026, 3, 25, 22, 2, 0, 0, time.UTC),
			orderStatus:         commerce.OrderStatusAccepted,
			decisionOutcome:     "review_required",
			decisionReasons:     []string{"high_value_delegated_purchase", "merchant_scope_match"},
			paymentMethod:       "delegated_card_on_file",
			merchant: commerce.Merchant{
				ID:       merchantID,
				Name:     "Orbit Books",
				Category: "books",
				Tags:     []string{"baseline", "high-value"},
			},
			product: commerce.Product{
				ID:         "product-orbit-book-bundle-premium",
				MerchantID: merchantID,
				Name:       "Orbit Enterprise Reference Bundle",
				Amount:     18900,
				Currency:   "USD",
				Category:   "reference",
				Tags:       []string{"high-value", "challenger"},
			},
			mandate:              &mandateBlueprint{Status: "active", AllowedActions: []string{"submit_order"}, SpendingLimit: 25000, ExpiresAt: time.Date(2026, 3, 26, 22, 0, 0, 0, time.UTC)},
			provenance:           &provenanceBlueprint{SourceType: "task_list", SourceRef: item.ID, Confidence: 0.79, CreatedAt: time.Date(2026, 3, 25, 22, 0, 20, 0, time.UTC)},
			outcomeSummary:       "High-value delegated purchase pushed materially above the buyer's normal spend baseline.",
			replayRecommendation: "candidate",
			replayReason:         "High-value delegated purchase is a strong sidecar recommendation scenario.",
			signals: map[string]any{
				"merchant_scope_match": true,
				"category_scope_match": true,
				"buyer_spend_baseline": 4200,
				"repeat_attempt_count": 1,
				"delegated_indicator":  true,
				"high_value_threshold": 12000,
			},
		}), nil
	default:
		return scenarioBlueprint{}, fmt.Errorf("scenario %s is not executable in Phase 9", item.ID)
	}
}

type purchaseProfile struct {
	orderSubmittedBy     string
	delegationMode       actors.DelegationMode
	orderCreatedAt       time.Time
	paymentAuthorizedAt  time.Time
	orderStatus          commerce.OrderStatus
	decisionOutcome      string
	decisionReasons      []string
	paymentMethod        string
	mandate              *mandateBlueprint
	provenance           *provenanceBlueprint
	merchant             commerce.Merchant
	product              commerce.Product
	outcomeSummary       string
	replayRecommendation string
	replayReason         string
	signals              map[string]any
}

func purchaseBlueprint(profile purchaseProfile) scenarioBlueprint {
	world := seedWorld()
	merchant := world.merchant
	if profile.merchant.ID != "" {
		merchant = profile.merchant
	}
	product := world.product
	if profile.product.ID != "" {
		product = profile.product
	}
	return scenarioBlueprint{
		FlowKind:                flowKindPurchase,
		Merchant:                merchant,
		Product:                 product,
		OrderSubmittedBy:        profile.orderSubmittedBy,
		DelegationMode:          profile.delegationMode,
		Mandate:                 profile.mandate,
		Provenance:              profile.provenance,
		OrderStatus:             profile.orderStatus,
		DecisionOutcome:         profile.decisionOutcome,
		DecisionReasonCodes:     profile.decisionReasons,
		StepUpRequired:          strings.Contains(profile.decisionOutcome, "review") || strings.Contains(profile.decisionOutcome, "step_up"),
		PaymentMethod:           profile.paymentMethod,
		OrderCreatedAt:          profile.orderCreatedAt,
		PaymentAuthorizedAt:     profile.paymentAuthorizedAt,
		OutcomeSummary:          profile.outcomeSummary,
		ReplayRecommendation:    profile.replayRecommendation,
		ReplayReason:            profile.replayReason,
		SignalContext:           profile.signals,
		TrustDecisionRecordedAt: profile.orderCreatedAt.Add(45 * time.Second),
	}
}

type refundProfile struct {
	orderSubmittedBy     string
	refundRequestedBy    string
	delegationMode       actors.DelegationMode
	orderCreatedAt       time.Time
	paymentAuthorizedAt  time.Time
	refundRequestedAt    time.Time
	orderStatus          commerce.OrderStatus
	refundStatus         commerce.RefundStatus
	decisionOutcome      string
	decisionReasons      []string
	paymentMethod        string
	reason               string
	mandate              *mandateBlueprint
	provenance           *provenanceBlueprint
	approval             *approvalBlueprint
	outcomeSummary       string
	replayRecommendation string
	replayReason         string
	signals              map[string]any
}

func refundBlueprint(profile refundProfile) scenarioBlueprint {
	world := seedWorld()
	return scenarioBlueprint{
		FlowKind:                flowKindRefund,
		Merchant:                world.merchant,
		Product:                 world.product,
		OrderSubmittedBy:        profile.orderSubmittedBy,
		RefundRequestedBy:       profile.refundRequestedBy,
		DelegationMode:          profile.delegationMode,
		Mandate:                 profile.mandate,
		Provenance:              profile.provenance,
		Approval:                profile.approval,
		OrderStatus:             profile.orderStatus,
		RefundStatus:            profile.refundStatus,
		DecisionOutcome:         profile.decisionOutcome,
		DecisionReasonCodes:     profile.decisionReasons,
		StepUpRequired:          strings.Contains(profile.decisionOutcome, "review") || strings.Contains(profile.decisionOutcome, "step_up"),
		PaymentMethod:           profile.paymentMethod,
		OrderCreatedAt:          profile.orderCreatedAt,
		PaymentAuthorizedAt:     profile.paymentAuthorizedAt,
		RefundRequestedAt:       profile.refundRequestedAt,
		OutcomeSummary:          profile.outcomeSummary,
		ReplayRecommendation:    profile.replayRecommendation,
		ReplayReason:            profile.replayReason,
		SignalContext:           profile.signals,
		ApprovalActionType:      "refund_request",
		DecisionEntityType:      entityTypeRefund,
		OrderReason:             profile.reason,
		TrustDecisionRecordedAt: profile.refundRequestedAt.Add(30 * time.Second),
	}
}

func signalContext(item domainscenario.Scenario, world worldSeed, plan scenarioBlueprint) map[string]any {
	context := map[string]any{
		"scenario_code":           item.Code,
		"scenario_family":         item.Family,
		"set_role":                item.SetRole,
		"tier_a_features":         append([]string(nil), item.FeatureModel.TierA...),
		"tier_b_features":         append([]string(nil), item.FeatureModel.TierB...),
		"tier_c_features":         append([]string(nil), item.FeatureModel.TierC...),
		"tier_c_available":        len(item.FeatureModel.TierC) > 0 && (plan.Mandate != nil || plan.Provenance != nil || plan.Approval != nil),
		"delegated_indicator":     plan.DelegationMode != actors.DelegationModeDirectHuman,
		"submitted_by_actor":      plan.OrderSubmittedBy,
		"refund_requested_by":     plan.RefundRequestedBy,
		"merchant_category":       world.merchant.Category,
		"order_amount":            world.product.Amount,
		"approval_present":        plan.Approval != nil && plan.Approval.Outcome == "approved",
		"approval_removed":        plan.Approval != nil && plan.Approval.Outcome == "removed",
		"approval_missing":        plan.Approval == nil || plan.Approval.Outcome == "missing",
		"mandate_status":          "",
		"provenance_confidence":   0.0,
		"tier_c_optional_only":    len(item.FeatureModel.TierC) > 0,
		"existing_control_note":   "Designed to run beside an existing fraud stack in shadow mode.",
		"blocking_mode":           "recommendation_only",
		"evaluation_mode":         "shadow",
		"tier_c_used_in_scenario": len(item.FeatureModel.TierC) > 0,
	}
	if plan.Mandate != nil {
		context["mandate_status"] = plan.Mandate.Status
	}
	if plan.Provenance != nil {
		context["provenance_confidence"] = plan.Provenance.Confidence
	}
	for key, value := range plan.SignalContext {
		context[key] = value
	}
	return context
}

type worldSeed struct {
	buyer      commerce.Buyer
	merchant   commerce.Merchant
	product    commerce.Product
	human      actors.HumanActor
	agent      actors.AgentActor
	principal  actors.PrincipalRef
	delegation actors.DelegationMode
}

func seedWorld() worldSeed {
	principal := actors.PrincipalRef{PrincipalID: "buyer-alex", PrincipalType: "buyer"}
	const merchantID = "merchant-orbit-books"
	const humanID = "human-alex"
	return worldSeed{
		buyer: commerce.Buyer{
			ID:       "buyer-alex",
			Name:     "Alex Carter",
			RiskTier: "standard",
			Tags:     []string{"baseline", "trusted"},
		},
		merchant: commerce.Merchant{
			ID:       merchantID,
			Name:     "Orbit Books",
			Category: "books",
			Tags:     []string{"digital", "baseline"},
		},
		product: commerce.Product{
			ID:         "product-orbit-book-1",
			MerchantID: merchantID,
			Name:       "Orbit Operations Handbook",
			Amount:     4200,
			Currency:   "USD",
			Category:   "reference",
			Tags:       []string{"starter", "digital"},
		},
		human: actors.HumanActor{
			ID:        humanID,
			Name:      "Alex Carter",
			Type:      actors.ActorTypeHuman,
			Principal: principal,
			Tags:      []string{"buyer"},
		},
		agent: actors.AgentActor{
			ID:   "agent-shopping-assistant",
			Name: "Shopping Assistant",
			Type: actors.ActorTypeAgent,
			Role: agents.AgentRoleOperator,
			Runtime: agents.RuntimeRef{
				Runtime: "zeroclaw",
				Version: "phase-9",
				Gateway: "omniroute",
			},
			Principal: principal,
			Tags:      []string{"commerce", "delegate"},
		},
		principal:  principal,
		delegation: actors.DelegationModeAgentAssisted,
	}
}

func (s *Service) seedParticipants(world worldSeed) {
	s.commerce.PutBuyer(world.buyer)
	s.commerce.PutMerchant(world.merchant)
	s.commerce.PutProduct(world.product)
}

func (s *Service) storeResult(result ExecutionResult) {
	s.results[result.Scenario.ID] = result
	for _, orderRef := range result.Entities.OrderRefs {
		s.orderRefs[orderRef] = result.Scenario.ID
	}
}

func actorIDFromKind(world worldSeed, kind string) string {
	if kind == actorKindAgent {
		return world.agent.ID
	}
	return world.human.ID
}

func deterministicID(prefix string, scenarioID string) string {
	replacer := strings.NewReplacer("commerce-", "", "_", "-", " ", "-", "/", "-")
	return prefix + "-" + replacer.Replace(scenarioID)
}
