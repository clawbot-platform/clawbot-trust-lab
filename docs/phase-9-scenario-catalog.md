# Phase 9 Scenario Catalog
## Clawbot Trust Lab for Agentic Commerce
## Replay-Driven Adversarial Regression Harness for Fraud Controls

## Purpose

This document defines the Phase 9 scenario catalog for the Red Queen demo hardening phase.

The goal of this catalog is to ensure that:
- the demo uses scenarios that map to realistic e-commerce and fraud data
- feature requirements stay aligned with what merchants and existing fraud stacks can actually provide
- human and agent behavior can be distinguished without requiring exotic telemetry
- the Red Queen loop can demonstrate real operational value over a week-long shadow run

This catalog groups features into:

- **Tier A** — commonly available today in merchant, PSP, or fraud-platform data
- **Tier B** — derivable with light engineering from existing history or event data
- **Tier C** — optional agentic overlay fields that increase differentiation but should not be required for the harness to be useful

---

## Tier Definitions

## Tier A — Required / Commonly Available

These are fields that a typical e-commerce business, PSP, or standard fraud stack usually has or can already send:
- amount
- currency
- merchant/category
- item count
- customer/account id
- email
- billing/shipping address
- shopper IP
- device fingerprint or equivalent device signal
- payment authorization result
- 3DS/authentication result
- order status
- refund requested / refund status
- manual review outcome if available

## Tier B — Recommended / Derived from Existing History

These are realistic but usually require simple joins or history windows:
- account age
- prior order count
- recent attempt count
- recent refund count
- merchant/category drift count
- time since last related action
- payment instrument reuse
- prior review/step-up count
- prior replay/promoted case count

## Tier C — Optional Agentic Overlay

These make the harness more differentiated for agentic commerce but should not be mandatory:
- actor_type = human / agent / unknown
- delegated_action_flag
- delegation_mode
- mandate_present
- mandate_active
- mandate_scope_match
- approval_present
- approval_actor_match
- provenance_present
- provenance_confidence
- delegation_depth

---

## Scenario Design Principles

Each scenario should:
- be explainable to a fraud or trust operator
- map to realistic commerce or review data
- demonstrate a business-relevant distinction
- contribute either to:
  - stable-set benchmarking
  - living-set challenger generation
  - replay regression value

For Phase 9, the system should work with **Tier A + Tier B** alone, and become more differentiated when **Tier C** is present.

---

# Scenario Catalog

## H1 — Direct Human Purchase

### Summary
A normal direct purchase initiated by a human customer with no delegated actor and no unusual trust gaps.

### Scenario role
- stable set
- clean baseline
- regression safety check

### Tier A fields
- amount
- currency
- merchant_id
- merchant_category
- item_count
- account_id
- email
- billing_address
- shipping_address
- shopper_ip
- device_signal
- payment_auth_result
- three_ds_result
- order_status

### Tier B fields
- account_age
- prior_order_count
- recent_attempt_count
- time_since_last_related_action

### Tier C fields
- actor_type = human (optional)

### Expected baseline outcome
- clean
- low risk
- allow

### Business value demonstrated
- proves the system is not simply hostile to activity or change
- provides a stable human baseline for comparing agentic and adversarial behavior
- shows that the harness can distinguish normal transaction behavior from suspicious variants

---

## H2 — Human Refund with Valid History

### Summary
A human customer requests a refund on a legitimate order with consistent purchase history and no unusual gaps.

### Scenario role
- stable set
- clean post-purchase baseline
- regression safety check

### Tier A fields
- account_id
- order_id
- amount
- currency
- order_status
- refund_requested
- refund_status
- payment_auth_result
- manual_review_outcome (if present)

### Tier B fields
- prior_order_count
- recent_refund_count
- time_since_purchase
- account_age

### Tier C fields
- actor_type = human (optional)

### Expected baseline outcome
- clean or low-risk
- allow or low-friction review

### Business value demonstrated
- shows that legitimate refund activity is not automatically treated as suspicious
- provides contrast with suspicious delegated refund scenarios
- helps quantify false-positive-like behavior over time

---

## A1 — Agent-Assisted Purchase with Valid Controls

### Summary
An agent assists with purchase execution, but the transaction remains within normal constraints and the request has clear authorization/provenance signals when available.

### Scenario role
- stable set
- benign agentic baseline
- regression safety check

### Tier A fields
- amount
- currency
- merchant_id
- merchant_category
- item_count
- account_id
- email
- shopper_ip
- device_signal
- payment_auth_result
- three_ds_result
- order_status

### Tier B fields
- account_age
- prior_order_count
- recent_attempt_count
- merchant_category_history

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- approval_present = true
- mandate_present = true
- provenance_present = true

### Expected baseline outcome
- clean
- low risk
- allow

### Business value demonstrated
- proves the system is not “agent = suspicious”
- shows the harness can distinguish benign automation from risky delegation
- provides a clean delegated-action baseline for later challenger variants

---

## A2 — Fully Delegated Replenishment Purchase

### Summary
A recurring or replenishment purchase is executed automatically under a known pattern and low-risk transaction context.

### Scenario role
- stable set
- benign automation baseline
- regression safety check

### Tier A fields
- amount
- currency
- merchant_id
- merchant_category
- account_id
- item_count
- payment_auth_result
- order_status

### Tier B fields
- prior_order_count
- repeat_purchase_pattern
- recent_attempt_count
- time_since_last_purchase

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- delegation_mode = fully_delegated
- mandate_present = true

### Expected baseline outcome
- clean
- low risk
- allow

### Business value demonstrated
- shows value for routine automated commerce
- demonstrates that the Red Queen system does not punish recurring delegated activity by default
- helps validate that controls preserve low-friction paths for safe automation

---

## A3 — Agent-Assisted Refund with Approval Evidence

### Summary
An agent requests a refund, but approval evidence or equivalent authorization is present, making the action traceable and explainable.

### Scenario role
- stable set
- benign delegated refund baseline
- contrast case for suspicious delegated refunds

### Tier A fields
- order_id
- refund_requested
- refund_status
- amount
- currency
- order_status
- payment_auth_result
- manual_review_outcome (if present)

### Tier B fields
- recent_refund_count
- prior_order_count
- time_since_purchase
- prior_review_count

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- approval_present = true
- approval_actor_match = true
- provenance_present = true

### Expected baseline outcome
- clean or low-risk
- allow or review-lite

### Business value demonstrated
- shows that delegated refund flows can be safe when evidence is present
- gives a realistic enterprise story for low-friction automation with review controls
- makes later suspicious refund variants more credible

---

## S1 — Refund Attempt with Weak or Missing Authorization

### Summary
A refund is requested under weak, missing, or inconsistent authorization conditions.

### Scenario role
- stable suspicious baseline
- replay candidate
- Red Queen anchor scenario

### Tier A fields
- order_id
- refund_requested
- refund_status
- amount
- currency
- order_status
- payment_auth_result
- manual_review_outcome (if present)

### Tier B fields
- recent_refund_count
- prior_review_count
- time_since_purchase
- prior_step_up_count
- prior_replay_count

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- approval_present = false
- mandate_present = false or inactive
- provenance_present = weak or missing

### Expected baseline outcome
- suspicious or step_up_required
- high risk
- review or step-up

### Business value demonstrated
- directly maps to real refund abuse concerns
- can be implemented with ordinary refund and approval data
- shows how the harness can sit beside existing fraud controls in shadow mode

---

## S2 — Delegated Purchase with Weak Provenance

### Summary
A delegated purchase is submitted without strong provenance or sufficient request-chain evidence.

### Scenario role
- stable suspicious baseline
- challenger family seed
- replay candidate

### Tier A fields
- amount
- currency
- merchant_id
- merchant_category
- account_id
- payment_auth_result
- order_status

### Tier B fields
- merchant_category_history
- account_age
- prior_order_count
- recent_attempt_count

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- provenance_present = false or weak
- mandate_present = true or ambiguous

### Expected baseline outcome
- suspicious or step_up_required

### Business value demonstrated
- highlights a realistic gap in delegated commerce flows
- shows how the system can distinguish “delegated but explainable” from “delegated but weakly evidenced”
- creates a strong replay/promotion candidate if missed

---

## S3 — Approval Removed After Initial Authorization

### Summary
A transaction begins with apparently valid authorization, but a later sensitive action occurs after approval evidence is removed, missing, or stale.

### Scenario role
- living set / challenger seed
- step-up candidate
- regression probe

### Tier A fields
- order_id
- refund_requested or change_requested
- amount
- currency
- order_status
- payment_auth_result

### Tier B fields
- time_since_last_related_action
- prior_review_count
- sequence_position_in_flow

### Tier C fields
- approval_present = false at final step
- actor_type = agent
- delegated_action_flag = true
- provenance_present = partial

### Expected baseline outcome
- suspicious or step_up_required

### Business value demonstrated
- shows multi-step workflow drift rather than single-event fraud
- demonstrates that the harness can test state transitions, not just static cases
- resembles real-world “authorization looked fine at step 1, not at step 3” problems

---

## S4 — Repeated Agent Refund Attempts

### Summary
The same delegated or agent-driven context produces multiple refund attempts inside a short time window.

### Scenario role
- stable suspicious baseline
- replay candidate
- strong week-long trend candidate

### Tier A fields
- account_id
- order_id
- refund_requested
- refund_status
- amount
- currency
- payment_auth_result

### Tier B fields
- recent_refund_count
- recent_attempt_count
- time_since_last_related_action
- prior_review_count
- prior_replay_count

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- approval_present = false

### Expected baseline outcome
- suspicious or step_up_required
- higher risk than one-off refund attempt

### Business value demonstrated
- highly realistic with standard fraud data
- shows why replay and round-to-round tracking matter
- good candidate for week-long trend reporting and production recommendations

---

## S5 — Merchant or Category Scope Drift Under Delegated Action

### Summary
A delegated action occurs outside the merchant or category pattern expected for that account or delegated policy context.

### Scenario role
- living set / challenger seed
- policy-drift candidate
- replay candidate

### Tier A fields
- merchant_id
- merchant_category
- amount
- currency
- account_id
- order_status
- payment_auth_result

### Tier B fields
- merchant_switch_count
- merchant_category_history
- prior_order_count
- recent_attempt_count

### Tier C fields
- actor_type = agent
- delegated_action_flag = true
- mandate_scope_match = false
- mandate_present = true or ambiguous

### Expected baseline outcome
- suspicious or step_up_required

### Business value demonstrated
- maps to real existing merchant/category drift logic used in fraud controls
- easy to explain to fraud operators
- demonstrates that policy-aware delegated activity is different from unconstrained delegated activity

---

## Challenger Variant Catalog

These are transformations that mutate baseline scenarios into living-set adversarial tests.

## V1 — Weakened Provenance
### Applied to
- A1
- A3
- S2

### Tier usage
- Tier A: unchanged
- Tier B: unchanged
- Tier C: `provenance_present` flips from strong to weak/missing

### Value demonstrated
- tests whether delegated-action evidence actually matters
- strong source of blind spots and promotions

## V2 — Expired or Inactive Mandate
### Applied to
- A1
- A2
- S1

### Tier usage
- Tier A: unchanged
- Tier B: unchanged
- Tier C: `mandate_present=true`, `mandate_active=false`

### Value demonstrated
- tests whether “mandate exists” is treated differently from “mandate is valid”
- realistic production recommendation candidate

## V3 — Approval Removed
### Applied to
- A3
- S1
- S3

### Tier usage
- Tier A: unchanged
- Tier B: sequence/timing becomes more important
- Tier C: `approval_present=false`

### Value demonstrated
- probes late-stage workflow gaps
- particularly good for replay-driven regression testing

## V4 — Actor Switch from Human to Agent
### Applied to
- H1
- H2

### Tier usage
- Tier A: unchanged
- Tier B: unchanged
- Tier C: `actor_type` flips to `agent`, `delegated_action_flag=true`

### Value demonstrated
- isolates the effect of automation itself
- helps prove the system does not simply treat all automated activity the same

## V5 — Repeat Attempt Escalation
### Applied to
- S1
- S4

### Tier usage
- Tier A: unchanged
- Tier B: `recent_attempt_count`, `recent_refund_count` increase
- Tier C: optional actor overlay unchanged

### Value demonstrated
- strong week-long demonstration value
- shows how repeated challenger behavior can become increasingly meaningful in replay/history

## V6 — Merchant Scope Drift
### Applied to
- A1
- A2
- S5

### Tier usage
- Tier A: merchant/category changes
- Tier B: merchant/category history and drift count become active
- Tier C: optional `mandate_scope_match=false`

### Value demonstrated
- easy to explain in enterprise demos
- maps well to existing fraud/risk rule thinking

## V7 — High-Value Delegated Purchase
### Applied to
- A1
- A2

### Tier usage
- Tier A: amount increases materially
- Tier B: amount delta vs prior behavior becomes active
- Tier C: delegated-action overlay unchanged

### Value demonstrated
- shows how standard commerce risk signals combine with delegated-action signals
- production-friendly scenario because amount-based controls already exist everywhere

---

## Recommended Phase 9 Stable Set

Use these as the recurring baseline set in every round:
- H1 — Direct Human Purchase
- H2 — Human Refund with Valid History
- A1 — Agent-Assisted Purchase with Valid Controls
- A2 — Fully Delegated Replenishment Purchase
- A3 — Agent-Assisted Refund with Approval Evidence
- S1 — Refund Attempt with Weak or Missing Authorization
- S4 — Repeated Agent Refund Attempts

This gives:
- human baseline
- benign automation baseline
- suspicious delegated baseline

---

## Recommended Phase 9 Living Set

Use these as rotating or challenger-driven round inputs:
- S2 — Delegated Purchase with Weak Provenance
- S3 — Approval Removed After Initial Authorization
- S5 — Merchant or Category Scope Drift Under Delegated Action
- V1 — Weakened Provenance
- V2 — Expired or Inactive Mandate
- V3 — Approval Removed
- V5 — Repeat Attempt Escalation
- V6 — Merchant Scope Drift
- V7 — High-Value Delegated Purchase

---

## Business Value Summary by Scenario Family

### Human baselines
Demonstrate that the harness does not punish normal customer activity.

### Benign agentic scenarios
Demonstrate that the harness can separate safe automation from risky delegation.

### Suspicious delegated scenarios
Demonstrate where existing controls may weaken as agent-driven tactics evolve.

### Challenger variants
Demonstrate why static benchmarking is not enough and why replay-driven adversarial regression matters.

---

## Production-Bridge Message

Phase 9 scenarios should support this enterprise message:

> This system does not require a merchant to replace its existing fraud engine.  
> It can run in shadow mode using ordinary commerce, payment, review, and event data that businesses already collect.  
> As teams begin tagging delegated or agent-driven actions, the harness becomes even more effective at distinguishing safe automation from risky automation.

---

## Acceptance Rule for New Scenarios

Before adding any new Phase 9 scenario, verify:

1. Does it rely primarily on Tier A and Tier B fields?
2. If it uses Tier C, can it still degrade gracefully without them?
3. Does it demonstrate a business-relevant distinction?
4. Does it improve the stable set, living set, or replay-regression value?

If the answer is no to most of those, the scenario probably does not belong in Phase 9.
