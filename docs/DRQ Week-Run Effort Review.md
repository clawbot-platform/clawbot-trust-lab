# Detailed Report
## DRQ Week-Run Effort Review
### What was used, why it was used, what was learned, and what the next version should be

## 1. Purpose of the effort

This effort was designed to answer a practical question:

**Can a lightweight, homelab-hosted benchmark and replay harness reliably expose detector weaknesses, support controlled tuning, and demonstrate measurable improvement over time?**

The answer from this run is **yes**.

The week run showed that the harness can:
- surface recurring weak cases
- preserve them through replay and reporting
- support a controlled detector-tuning pass
- show measurable before/after improvement over a scheduled run window

## 2. What the effort used

## 2.1 Platform layout

The run used a distributed homelab layout:

- **OptiPlex 1** hosted the control plane (`clawbot-server`)
- **ThinkPad P50** hosted the memory service (`clawmem`)
- **OptiPlex 2** hosted Trust Lab and its UI
- artifacts were written to host-mounted report and replay directories
- images were pulled from GHCR and run directly with Docker

This matters because the effort was not only about detector logic. It also validated a practical operating model for benchmark execution on modest hardware with minimal runtime dependency footprint.

## 2.2 Deployment model

The deployment model used:
- Docker containers only on runtime hosts
- public GHCR images
- host-mounted artifact directories
- no host-side Go builds or local application toolchain on the run nodes

This model was the right choice because it:
- reduced drift between nodes
- simplified cutover between baseline and tuned images
- made the run repeatable
- aligned with the goal that runtime hosts should behave like appliance nodes rather than development workstations

## 2.3 Benchmark structure

The benchmark was built around:
- a **stable set**
- a **living/challenger set**
- replay and report generation
- scheduled execution at fixed intervals
- a mid-run detector tuning cutover

This structure was the right choice because it allowed the run to measure both:
- **stability over time**
- **detector improvement after tuning**

That combination is what made the week run meaningful. A one-time benchmark would have shown misses, but the repeated scheduler window showed that the misses were persistent in Phase A and sustainably corrected in Phase B. :contentReference[oaicite:17]{index=17}

## 2.4 Reporting model

The run generated:
- per-round reports
- daily reports
- management reports
- replay archive artifacts

This was important because the goal was not just engineering verification. The effort also needed to produce human-readable evidence for comparison, review, and decision-making.

## 3. What the effort proved

## 3.1 Baseline detector weaknesses were real and recurring

Phase A repeatedly surfaced the same three weak cases:
- `commerce-v2-expired-inactive-mandate`
- `commerce-v3-approval-removed`
- `commerce-s3-approval-removed-after-authorization`

These were not one-off failures. They repeated across the baseline rounds and produced the same promotion pattern. :contentReference[oaicite:18]{index=18}

## 3.2 Narrow tuning was enough to improve performance materially

The tuned detector did not require a wholesale system redesign. A narrow rule update targeted at the three known weak cases was enough to produce the desired behavior change:
- promotions dropped to zero
- replay pass rate held at one
- the targeted cases now met `step_up_required` rather than missing the floor again :contentReference[oaicite:19]{index=19} :contentReference[oaicite:20]{index=20}

That is a strong result because it shows the system can improve through focused iteration instead of broad and risky change.

## 3.3 The deployment model is viable

The run validated that the platform can be operated in a distributed, image-based model with:
- stable service health
- reproducible image pulls
- usable artifact collection
- clean phase cutover

This matters for future work because it lowers the operational bar for rerunning the benchmark or standing up new detector experiments.

## 4. Why this architecture was appropriate

This effort did **not** use an LLM in the detector path, and that was the right choice for this version.

Why:
- the key objective was controlled replay, benchmark repeatability, and measurable detector tuning
- deterministic or rule-like detector behavior is easier to compare across rounds
- a narrowly scoped benchmark is easier to validate when the detector behavior is explainable
- the targeted weak cases were specific enough that rules and policy posture changes were sufficient

In other words, the goal of this version was **not intelligence expansion**. It was **benchmark discipline**.

That was the right design choice for a first serious DRQ week run.

## 5. What the next version should be

The next version should build on the tuned detector baseline, not restart from the original baseline.

There are two credible paths.

## 5.1 Next version without an LLM

This is the lower-risk and more operationally disciplined path.

### What it would include
- more challenger variants
- more replay-stable cases
- better long-run summaries
- clearer split between discovered blind spots and historical ones
- better grouping of recommendations by root cause
- more formalized phase comparison tooling
- stronger regression gates around replay-stable scenarios

### Why this path makes sense
- easier to validate
- easier to explain
- easier to operationalize
- keeps the benchmark harness deterministic
- strengthens the benchmark foundation before adding probabilistic reasoning

### Best use case
Use this path if the priority is:
- detector quality
- repeatability
- regression control
- disciplined benchmark expansion

## 5.2 Next version with an LLM

This is the higher-ambition path, but it should be introduced carefully.

### Where an LLM would make sense
An LLM should **not** replace the detector core first. It should be introduced in a sidecar role, for example:
- case explanation synthesis
- rationale generation for recommendations
- scenario clustering and grouping
- analyst-facing narrative summaries
- policy suggestion drafts
- triage prioritization assistance

### Why that is the right entry point
It preserves benchmark comparability while still gaining value from LLM capability.

If an LLM is put directly into the scoring path too early:
- determinism drops
- replay comparability gets harder
- debugging gets harder
- before/after benchmark interpretation gets weaker

### Best use case
Use this path if the priority is:
- analyst workflow augmentation
- richer reporting
- narrative explanations
- exploratory pattern grouping beyond the current deterministic rules

## 5.3 What not to do next

The next version should **not**:
- collapse deterministic detection and LLM reasoning into one opaque path
- add broad infrastructure churn unrelated to benchmark value
- lose the clean Phase A / Phase B comparison model
- overfit to one week-run without preserving replay and challenger discipline

## 6. Recommended next roadmap

## Short-term next step
Treat the current tuned detector as the new benchmark reference state.

## Next benchmark milestone
Run a second DRQ cycle with:
- expanded challenger coverage
- tighter report comparison tooling
- clearer replay-stable set governance
- formal summary generation that correctly segments by detector version

## LLM roadmap
If an LLM is introduced, do it in this order:
1. report and explanation sidecar
2. recommendation clustering and analyst narrative support
3. optional shadow-mode advisory scoring
4. only later consider any deeper detector-path role

## 7. Final conclusion

This effort was successful on both engineering and benchmark terms.

It proved that:
- the benchmark harness can run reliably over time
- recurring weaknesses can be identified
- a tuned detector can measurably improve performance
- the deployment model is practical
- a clean phase-based evaluation methodology works

The strongest conclusion is this:

**The next version should preserve the deterministic benchmark core and only add LLM capability in a controlled sidecar role, if at all.**

That gives you the best balance of:
- explainability
- repeatability
- benchmark integrity
- future extensibility

## Final recommendation

Use the tuned Phase B detector as the new reference point, expand the deterministic benchmark harness next, and introduce LLM capability only where it adds explanation or analyst productivity without weakening replay discipline.
