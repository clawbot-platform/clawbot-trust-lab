package scenario

import (
	"context"
	"fmt"
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
}

func NewService(scenarios ScenarioCatalog, commerce *commerceSvc.Service, events *eventsvc.Service, trust *trustsvc.Service, artifacts TrustArtifactWriter, replay ReplayWriter) *Service {
	return &Service{
		scenarios: scenarios,
		commerce:  commerce,
		events:    events,
		trust:     trust,
		artifacts: artifacts,
		replay:    replay,
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

	switch item.ID {
	case "commerce-clean-agent-assisted-purchase":
		return s.executeCleanPurchase(ctx, item)
	case "commerce-suspicious-refund-attempt":
		return s.executeSuspiciousRefund(ctx, item)
	default:
		return ExecutionResult{}, fmt.Errorf("scenario %s is not executable in Phase 5", item.ID)
	}
}

func (s *Service) executeCleanPurchase(ctx context.Context, item domainscenario.Scenario) (ExecutionResult, error) {
	world := seedWorld()
	s.seedParticipants(world)

	orderCreatedAt := time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC)
	paymentAuthorizedAt := orderCreatedAt.Add(2 * time.Minute)

	mandate := s.trust.RecordMandate(domaintrust.Mandate{
		ID:              "mandate-clean-agent-assisted-purchase",
		PrincipalID:     world.principal.PrincipalID,
		DelegateActorID: world.agent.ID,
		AllowedActions:  []string{"submit_order"},
		SpendingLimit:   15000,
		ExpiresAt:       orderCreatedAt.Add(24 * time.Hour),
		Status:          "active",
	})
	provenance := s.trust.RecordProvenance(domaintrust.ProvenanceRecord{
		ID:          "prov-clean-agent-assisted-purchase",
		ActorID:     world.agent.ID,
		PrincipalID: world.principal.PrincipalID,
		SourceType:  "conversation",
		SourceRef:   "commerce-pack/v1/clean-agent-assisted-purchase",
		Confidence:  0.96,
		CreatedAt:   orderCreatedAt.Add(30 * time.Second),
	})

	order, err := s.commerce.CreateOrder(commerceSvc.CreateOrderInput{
		ID:                 "order-clean-agent-assisted-purchase",
		BuyerID:            world.buyer.ID,
		MerchantID:         world.merchant.ID,
		ProductIDs:         []string{world.product.ID},
		SubmittedByActorID: world.agent.ID,
		DelegationMode:     world.delegation,
		MandateRef:         mandate.ID,
		ProvenanceRef:      provenance.ID,
		Status:             commerce.OrderStatusAccepted,
		CreatedAt:          orderCreatedAt,
		UpdatedAt:          paymentAuthorizedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	payment := s.commerce.CreatePayment(commerceSvc.CreatePaymentInput{
		ID:           "payment-clean-agent-assisted-purchase",
		OrderID:      order.ID,
		Amount:       order.TotalAmount,
		Currency:     order.Currency,
		Status:       commerce.PaymentStatusAuthorized,
		Method:       "delegated_card_on_file",
		AuthorizedAt: paymentAuthorizedAt,
	})

	eventRefs := []string{
		s.events.RecordTrust("evt-clean-order-submitted", domainevents.TrustEventOrderSubmittedByAgent, "order", order.ID, item.ID, world.agent.ID, orderCreatedAt, map[string]any{"delegation_mode": world.delegation}).ID,
		s.events.RecordTrust("evt-clean-mandate-checked", domainevents.TrustEventMandateChecked, "mandate", mandate.ID, item.ID, world.agent.ID, orderCreatedAt.Add(15*time.Second), map[string]any{"status": mandate.Status, "spending_limit": mandate.SpendingLimit}).ID,
		s.events.RecordTrust("evt-clean-provenance-attached", domainevents.TrustEventProvenanceAttached, "provenance_record", provenance.ID, item.ID, world.agent.ID, orderCreatedAt.Add(30*time.Second), map[string]any{"confidence": provenance.Confidence}).ID,
		s.events.RecordTransaction("evt-clean-order-created", domainevents.TransactionEventOrderCreated, "order", order.ID, item.ID, world.agent.ID, orderCreatedAt, map[string]any{"total_amount": order.TotalAmount, "currency": order.Currency}).ID,
		s.events.RecordTransaction("evt-clean-payment-authorized", domainevents.TransactionEventPaymentAuthorized, "payment", payment.ID, item.ID, world.agent.ID, paymentAuthorizedAt, map[string]any{"order_id": order.ID}).ID,
	}

	decision, err := s.trust.RecordDecision(domaintrust.TrustDecision{
		ID:             "decision-clean-agent-assisted-purchase",
		EntityType:     "order",
		EntityID:       order.ID,
		Outcome:        "accepted",
		ReasonCodes:    []string{"active_mandate", "high_provenance_confidence"},
		MandateRef:     mandate.ID,
		ProvenanceRef:  provenance.ID,
		StepUpRequired: false,
		RecordedAt:     orderCreatedAt.Add(45 * time.Second),
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	eventRefs = append(eventRefs, s.events.RecordTrust("evt-clean-trust-decision", domainevents.TrustEventTrustDecisionRecorded, "trust_decision", decision.ID, item.ID, world.agent.ID, decision.RecordedAt, map[string]any{"outcome": decision.Outcome}).ID)

	artifact, err := s.artifacts.CreateArtifact(ctx, domaintrust.CreateArtifactInput{
		ScenarioID: item.ID,
		ArtifactID: "ta-commerce-clean-agent-assisted-purchase",
		CreatedAt:  orderCreatedAt.Add(time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	replayCase, err := s.replay.CreateCase(ctx, domainreplay.CreateCaseInput{
		CaseID:                  "rc-commerce-clean-agent-assisted-purchase",
		ScenarioID:              item.ID,
		TrustArtifactRefs:       []string{artifact.ID},
		BenchmarkRoundRef:       "phase-5-commerce-baseline",
		OutcomeSummary:          "Agent-assisted purchase completed with mandate and provenance coverage.",
		PromotionRecommendation: "promote",
		PromotionReason:         "Clean baseline flow is deterministic and trust-complete.",
		RecordedAt:              paymentAuthorizedAt.Add(time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	return ExecutionResult{
		Scenario: item,
		Entities: EntityRefs{
			BuyerRefs:         []string{world.buyer.ID},
			MerchantRefs:      []string{world.merchant.ID},
			ProductRefs:       []string{world.product.ID},
			OrderRefs:         []string{order.ID},
			PaymentRefs:       []string{payment.ID},
			TrustArtifactRefs: []string{artifact.ID},
		},
		TrustDecisions: []domaintrust.TrustDecision{decision},
		ReplayCaseRefs: []string{replayCase.ID},
		MemoryWrites: []MemoryWriteOutcome{
			{Kind: "trust_artifact", SourceID: artifact.ID, Status: "written"},
			{Kind: "replay_case", SourceID: replayCase.ID, Status: "written"},
		},
		EventRefs: eventRefs,
	}, nil
}

func (s *Service) executeSuspiciousRefund(ctx context.Context, item domainscenario.Scenario) (ExecutionResult, error) {
	world := seedWorld()
	s.seedParticipants(world)

	orderCreatedAt := time.Date(2026, 3, 24, 10, 0, 0, 0, time.UTC)
	refundRequestedAt := orderCreatedAt.Add(8 * time.Minute)

	mandate := s.trust.RecordMandate(domaintrust.Mandate{
		ID:              "mandate-suspicious-refund-attempt",
		PrincipalID:     world.principal.PrincipalID,
		DelegateActorID: world.agent.ID,
		AllowedActions:  []string{"submit_order"},
		SpendingLimit:   15000,
		ExpiresAt:       orderCreatedAt.Add(-time.Minute),
		Status:          "expired",
	})
	provenance := s.trust.RecordProvenance(domaintrust.ProvenanceRecord{
		ID:          "prov-suspicious-refund-attempt",
		ActorID:     world.agent.ID,
		PrincipalID: world.principal.PrincipalID,
		SourceType:  "conversation",
		SourceRef:   "commerce-pack/v1/suspicious-refund-attempt",
		Confidence:  0.41,
		CreatedAt:   refundRequestedAt.Add(-time.Minute),
	})

	order, err := s.commerce.CreateOrder(commerceSvc.CreateOrderInput{
		ID:                 "order-suspicious-refund-attempt",
		BuyerID:            world.buyer.ID,
		MerchantID:         world.merchant.ID,
		ProductIDs:         []string{world.product.ID},
		SubmittedByActorID: world.human.ID,
		DelegationMode:     actors.DelegationModeDirectHuman,
		MandateRef:         mandate.ID,
		ProvenanceRef:      provenance.ID,
		Status:             commerce.OrderStatusRefundReview,
		CreatedAt:          orderCreatedAt,
		UpdatedAt:          refundRequestedAt,
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	payment := s.commerce.CreatePayment(commerceSvc.CreatePaymentInput{
		ID:           "payment-suspicious-refund-attempt",
		OrderID:      order.ID,
		Amount:       order.TotalAmount,
		Currency:     order.Currency,
		Status:       commerce.PaymentStatusAuthorized,
		Method:       "card_on_file",
		AuthorizedAt: orderCreatedAt.Add(time.Minute),
	})

	refund := s.commerce.CreateRefund(commerceSvc.CreateRefundInput{
		ID:                 "refund-suspicious-refund-attempt",
		OrderID:            order.ID,
		Amount:             order.TotalAmount,
		Status:             commerce.RefundStatusRejected,
		RequestedByActorID: world.agent.ID,
		Reason:             "Agent attempted refund without active authority",
		CreatedAt:          refundRequestedAt,
	})

	approval := s.trust.RecordApproval(domaintrust.ApprovalRecord{
		ID:         "approval-suspicious-refund-attempt",
		OrderID:    order.ID,
		ActionType: "refund_request",
		ApproverID: world.human.ID,
		Outcome:    "missing",
		CreatedAt:  refundRequestedAt.Add(15 * time.Second),
	})

	decision, err := s.trust.RecordDecision(domaintrust.TrustDecision{
		ID:             "decision-suspicious-refund-attempt",
		EntityType:     "refund",
		EntityID:       refund.ID,
		Outcome:        "step_up_required",
		ReasonCodes:    []string{"expired_mandate", "low_provenance_confidence", "missing_human_approval"},
		MandateRef:     mandate.ID,
		ProvenanceRef:  provenance.ID,
		StepUpRequired: true,
		RecordedAt:     refundRequestedAt.Add(30 * time.Second),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	eventRefs := []string{
		s.events.RecordTransaction("evt-suspicious-order-created", domainevents.TransactionEventOrderCreated, "order", order.ID, item.ID, world.human.ID, orderCreatedAt, map[string]any{"total_amount": order.TotalAmount, "currency": order.Currency}).ID,
		s.events.RecordTransaction("evt-suspicious-payment-authorized", domainevents.TransactionEventPaymentAuthorized, "payment", payment.ID, item.ID, world.human.ID, payment.AuthorizedAt, map[string]any{"order_id": order.ID}).ID,
		s.events.RecordTransaction("evt-suspicious-refund-requested", domainevents.TransactionEventRefundRequested, "refund", refund.ID, item.ID, world.agent.ID, refundRequestedAt, map[string]any{"order_id": order.ID, "reason": refund.Reason}).ID,
		s.events.RecordTrust("evt-suspicious-order-submitted", domainevents.TrustEventOrderSubmittedByAgent, "refund", refund.ID, item.ID, world.agent.ID, refundRequestedAt, map[string]any{"delegation_mode": world.delegation}).ID,
		s.events.RecordTrust("evt-suspicious-mandate-checked", domainevents.TrustEventMandateChecked, "mandate", mandate.ID, item.ID, world.agent.ID, refundRequestedAt.Add(5*time.Second), map[string]any{"status": mandate.Status}).ID,
		s.events.RecordTrust("evt-suspicious-provenance", domainevents.TrustEventProvenanceAttached, "provenance_record", provenance.ID, item.ID, world.agent.ID, refundRequestedAt.Add(10*time.Second), map[string]any{"confidence": provenance.Confidence}).ID,
		s.events.RecordTrust("evt-suspicious-approval", domainevents.TrustEventApprovalRecorded, "approval_record", approval.ID, item.ID, world.agent.ID, approval.CreatedAt, map[string]any{"outcome": approval.Outcome}).ID,
		s.events.RecordTrust("evt-suspicious-trust-decision", domainevents.TrustEventTrustDecisionRecorded, "trust_decision", decision.ID, item.ID, world.agent.ID, decision.RecordedAt, map[string]any{"outcome": decision.Outcome, "step_up_required": decision.StepUpRequired}).ID,
		s.events.RecordTransaction("evt-suspicious-refund-decision", domainevents.TransactionEventRefundDecisionRecorded, "refund", refund.ID, item.ID, world.agent.ID, decision.RecordedAt.Add(5*time.Second), map[string]any{"status": refund.Status}).ID,
	}

	artifact, err := s.artifacts.CreateArtifact(ctx, domaintrust.CreateArtifactInput{
		ScenarioID: item.ID,
		ArtifactID: "ta-commerce-suspicious-refund-attempt",
		CreatedAt:  decision.RecordedAt.Add(time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	replayCase, err := s.replay.CreateCase(ctx, domainreplay.CreateCaseInput{
		CaseID:                  "rc-commerce-suspicious-refund-attempt",
		ScenarioID:              item.ID,
		TrustArtifactRefs:       []string{artifact.ID},
		BenchmarkRoundRef:       "phase-5-commerce-baseline",
		OutcomeSummary:          "Refund attempt required step-up because mandate, provenance, and approval coverage were insufficient.",
		PromotionRecommendation: "hold",
		PromotionReason:         "Suspicious path should remain available for replay and later detection baselines.",
		RecordedAt:              decision.RecordedAt.Add(2 * time.Minute),
	})
	if err != nil {
		return ExecutionResult{}, err
	}

	return ExecutionResult{
		Scenario: item,
		Entities: EntityRefs{
			BuyerRefs:         []string{world.buyer.ID},
			MerchantRefs:      []string{world.merchant.ID},
			ProductRefs:       []string{world.product.ID},
			OrderRefs:         []string{order.ID},
			PaymentRefs:       []string{payment.ID},
			RefundRefs:        []string{refund.ID},
			TrustArtifactRefs: []string{artifact.ID},
		},
		TrustDecisions: []domaintrust.TrustDecision{decision},
		ReplayCaseRefs: []string{replayCase.ID},
		MemoryWrites: []MemoryWriteOutcome{
			{Kind: "trust_artifact", SourceID: artifact.ID, Status: "written"},
			{Kind: "replay_case", SourceID: replayCase.ID, Status: "written"},
		},
		EventRefs: eventRefs,
	}, nil
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
	return worldSeed{
		buyer: commerce.Buyer{
			ID:       "buyer-alex",
			Name:     "Alex Carter",
			RiskTier: "standard",
			Tags:     []string{"baseline", "trusted"},
		},
		merchant: commerce.Merchant{
			ID:       "merchant-orbit-books",
			Name:     "Orbit Books",
			Category: "books",
			Tags:     []string{"digital", "baseline"},
		},
		product: commerce.Product{
			ID:         "product-orbit-book-1",
			MerchantID: "merchant-orbit-books",
			Name:       "Orbit Operations Handbook",
			Amount:     4200,
			Currency:   "USD",
			Category:   "reference",
			Tags:       []string{"starter", "digital"},
		},
		human: actors.HumanActor{
			ID:        "human-alex",
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
				Version: "phase-5",
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
