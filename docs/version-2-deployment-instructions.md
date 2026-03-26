# Clawbot Trust Lab Version 2
## Enterprise Sidekick for Fraud Teams

## Purpose

Clawbot Trust Lab Version 2 is the enterprise-oriented evolution of the platform.

Its purpose is to let incumbent fraud teams use **their own data, scenarios, features, and control assumptions** to pressure-test their current fraud stack in shadow mode.

Version 2 is not the same thing as Version 1.

Version 1 proves the concept works.

Version 2 turns that concept into an **enterprise sidekick** that can run in a fraud team’s internal environment and provide practical value without replacing the incumbent fraud solution.

---

## What Version 2 is

Version 2 is:

- a sidecar to an existing fraud stack
- a shadow-mode recommendation system
- a replay-driven regression harness using incumbent data
- a lab companion for fraud teams evaluating agentic-commerce risks
- a supportable enterprise offering that can be used in internal labs or pilot environments

---

## What Version 2 is not

Version 2 is not:

- a replacement for the incumbent fraud engine
- a mandatory production blocker
- a requirement to rebuild the current fraud stack
- a demand for exotic agent-only telemetry
- a full enterprise fraud platform

---

## Core idea

Version 2 allows an enterprise to use its own:

- transaction data
- refund data
- review outcomes
- detection signals
- internal scenarios
- risk assumptions
- optional delegated/agentic overlay signals

to drive value from Clawbot Trust Lab.

Instead of relying only on the built-in scenario set, the enterprise can:
- map its own scenarios into the harness
- map its own features into the harness
- compare its incumbent decisions vs the harness recommendations
- preserve blind spots in replay
- use the platform to improve fraud controls over time

---

## How an enterprise would use Version 2

### 1. Keep the incumbent fraud solution
The existing fraud stack remains the primary decision-maker.

### 2. Run Clawbot Trust Lab beside it
Version 2 runs in:
- `evaluation_mode = shadow`
- `blocking_mode = recommendation_only`

### 3. Send internal data and scenarios
The enterprise provides:
- normalized order / payment / refund / review data
- optional internal scenarios or challenger cases
- optional agentic overlays such as:
  - actor type
  - delegated action flag
  - approval evidence
  - mandate or provenance fields

### 4. Review outputs
The fraud team uses the harness to review:
- blind spots
- regressions
- promoted replay cases
- recommendations
- trend summaries

---

## How enterprises would send data

Version 2 should support enterprise onboarding through practical methods.

### Recommended first path — batch ingestion
The enterprise exports normalized data as:
- JSON
- JSONL
- CSV

and provides it through:
- mounted internal volume
- object storage
- secure internal file transfer

This is the easiest first enterprise path.

### Recommended next path — ingestion APIs
Version 2 should also support API-based intake for:
- transaction events
- refund events
- review decisions
- detection decisions
- batch scenario inputs

### Longer-term path
Later, Version 2 can support connectors for:
- internal event streams
- message buses
- fraud data pipelines

But that is not required for the initial enterprise lab offering.

---

## Feature model

Version 2 should continue using the same tier model:

### Tier A — commonly available now
Typical merchant / PSP / fraud-stack data:
- amount
- currency
- merchant/category
- order/refund status
- account id
- email
- IP
- device signal
- auth result
- 3DS result
- manual review outcomes

### Tier B — derivable with light engineering
- account age
- prior order count
- recent attempts
- recent refunds
- merchant drift
- prior review count
- prior replay/promotion count

### Tier C — optional enterprise agentic overlay
- actor type
- delegated action flag
- approval present
- mandate present
- provenance present

Version 2 must still work with Tier A + Tier B alone.
Tier C is what lets the enterprise get more differentiated value over time.

---

## Two deployment modes for Version 2

### Recommended enterprise trial mode
Use **Option A**:

Publish four images:

- `clawbot-server`
- `clawmem`
- `clawbot-trust-lab`
- `clawbot-trust-lab-ui`

This is the cleanest enterprise trial model because it separates:
- control-plane/runtime support
- memory/replay support
- trust lab execution
- operator UI

### Why Option A is recommended
It gives enterprise teams:
- clear deployment boundaries
- simpler internal hosting
- easier network and access control
- better supportability
- cleaner operational ownership

---

## Version 2 deployment model

Version 2 should be deployed as prebuilt containers, not as source checkout.

An incumbent enterprise should not need:
- local Go toolchains
- local Node.js toolchains
- GitHub source build steps

Instead, the deployment model should be:

1. pull approved container images from an internal registry
2. run them on an internal-only network
3. mount persistent volumes for reports and memory
4. mount or connect data-ingestion inputs
5. operate the system in shadow mode

---

## Example container set

Recommended images:

- `clawbot-server`
- `clawmem`
- `clawbot-trust-lab`
- `clawbot-trust-lab-ui`

Recommended volumes:
- reports
- replay/history data
- optional input batch directory

Recommended network posture:
- internal only
- no public exposure required
- recommendation-only mode during evaluation

---

## Enterprise value proposition

Version 2 should answer this question for an incumbent fraud team:

**How do we continuously test whether our current fraud controls still work as delegated and agent-driven commerce behavior evolves?**

The value is:

- continuous regression protection
- early discovery of blind spots
- replay preservation of prior failures
- structured recommendations for control improvement
- trend visibility over time
- a safer path to experimenting with agentic-commerce control logic

---

## What outputs enterprises would get

Version 2 should provide:

- benchmark round summaries
- promotion lists
- replay regression signals
- recommendation reports
- trend summaries
- recommendation APIs
- operator UI for rounds, promotions, reports, and recommendations

This allows a fraud team to answer:
- What suspicious patterns are we still catching?
- What new delegated or agentic behaviors are slipping through?
- What previously known weaknesses regressed?
- What should we tighten next?
- What should remain in shadow mode for observation?

---

## Recommended enterprise workflow

### Step 1 — onboard baseline data
Map normalized commerce and fraud data into the harness.

### Step 2 — run shadow rounds
Run Clawbot Trust Lab on a schedule using enterprise data and scenario inputs.

### Step 3 — review recommendations
Use recommendations to identify:
- refund control gaps
- delegated-action blind spots
- provenance or approval gaps
- replay additions
- step-up triggers worth testing

### Step 4 — tune incumbent controls
Use those findings to adjust:
- rules
- review routing
- step-up policies
- replay coverage
- delegated commerce guardrails

### Step 5 — rerun and compare
Use future rounds to confirm whether:
- controls improved
- regressions remain
- blind spots persist
- recommendations were effective

---

## What makes Version 2 commercially interesting

Version 2 opens the door to a support-based enterprise model.

A practical enterprise offering could be:

### Supported internal lab subscription
The enterprise uses Clawbot Trust Lab in its own lab or shadow environment, and receives:
- deployment support
- scenario onboarding support
- feature mapping support
- review/report interpretation support
- replay and calibration tuning support

This is a much more believable commercial path than claiming full production replacement from day one.

---

## Relationship between Version 1 and Version 2

### Version 1
- self-sufficient
- fixed built-in scenario model
- proves DRQ concept works
- demonstrates value in the commerce fraud domain

### Version 2
- enterprise-facing
- uses incumbent data, features, and scenarios
- drives practical fraud-team value
- can be used as a supportable internal-lab sidekick

Version 1 is the proof.

Version 2 is the enterprise bridge.

---

## Current maturity

Version 2 should be described honestly as:
- the enterprise direction
- partially enabled by current architecture
- not yet fully complete in all onboarding and ingestion areas
- a logical Phase 2 / enterprise onboarding expansion

That honesty improves credibility.

---

## Summary

Clawbot Trust Lab should now be explained in two versions:

### Version 1
A self-sufficient Red Queen / DRQ harness that proves the concept works in e-commerce fraud.

### Version 2
An enterprise sidekick that allows incumbents to use their own data, scenarios, and features to continuously pressure-test existing fraud controls in shadow mode.

This split makes the product story clearer:
- Version 1 proves the idea
- Version 2 is how enterprises eventually adopt it
