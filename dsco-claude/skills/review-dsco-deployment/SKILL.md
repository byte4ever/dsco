---
name: review-dsco-deployment
description: >-
  Adversarial review of the DEPLOYMENT lane of dsco config: inventory usage
  (inventory.Compute(&cfg), text/JSON/preflight flavours, hand-maintained key
  lists that should be generated), *Layers exported for cross-package reuse,
  version targeting (the go.mod pin actually provides the APIs used, esp.
  inventory ≥ v1.4.0), and third-party dsco-config embedding vs local
  redefinition. Usually run by the review-dsco orchestrator. Default: REJECT.
x-dsco-target: v1.4.0
x-bundle-version: 1.4.0
---

# review-dsco-deployment

Per-aspect reviewer for the **deployment** lane: operability and versioning.
Can operators (human or LLM) discover the required keys, does the tooling reuse
one layer definition, does the code compile against the pinned dsco, and does it
compose dependency configs so their keys surface?

## Disposition (load-bearing — do not soften)

Default verdict: REJECT. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue found. You are
not evaluating effort, intent, or investment. Treat the artifact as anonymous;
you do not receive the author's identity, history, or production context. If you
received this as an isolated sub-agent with an anonymised artifact, execute the
procedure directly.

## Activation

Applies when the artifact uses `inventory`, exposes a `*Layers` function,
embeds a third-party config, or targets deployment. Signals: `inventory.`,
`func .*Layers(`, `*<pkg>.Config`, a version-gated API. Under uncertainty the
orchestrator includes this lane.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no scoring)

1. **Version targeting** — the dsco version in `go.mod` provides every API used
   here (esp. `inventory`, min v1.4.0). If it does not, that is a finding:
   the code won't compile / advice is ahead of the pin.
2. **Inventory target** — `inventory.Compute(&cfg, ...)` (the `**T` rule).
3. **Inventory flavour** — output flavour matches the use: text (human), JSON
   (tooling / operator-LLM), preflight (exit 2 in CI/init).
4. **Generated vs hand-maintained** — no hand-maintained required-env list (in a
   README, k8s manifest, `.env.example`) that `inventory` should generate.
5. **`*Layers` export** — the layer function is EXPORTED (`Layers()`), so an
   inventory binary or a test in another package can call it. Unexported
   `layers()` blocks cross-package reuse and invites drift.
6. **Shared definition** — the `Fill` site and the inventory driver call the
   SAME `*Layers` function (no re-declaration).
7. **Third-party embedding** — a dependency's exported dsco-shaped `Config` is
   embedded (`*pgdriver.Config`), so inventory walks into it and operators see
   the full key surface — not re-declared field-for-field locally.

**Minimum threshold: 6 observations.** If fewer, justify triviality in writing.

Detailed checks: [references/checklist-deployment.md](references/checklist-deployment.md).

### Phase 2 — Adversarial scenarios

**Minimum 3 scenarios.** Trigger / Propagation / Symptom / Detectability. Cover
at least: a version-gated API against an older pin, and a drift between a
hand-maintained list (or a re-declared layer set) and the code.

### Phase 3 — Scoring

- **BLOCKING** — an inventory driver against dsco < v1.4.0 (does not
  compile); `inventory.Compute(cfg)` single pointer.
- **IMPORTANT** — `*Layers` unexported while an inventory binary/test needs it
  (drift risk); a dependency config re-declared locally so inventory misses its
  keys; a hand-maintained key list that can drift from the code.
- **NOTED** — no inventory used but the shape is inventory-ready (out of scope);
  a version pin that is fine because only core APIs are used.

Calibration: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md).

### Phase 4 — Verdict

≥1 BLOCKING → REJECT. Only IMPORTANT/NOTED → CONDITIONAL APPROVE. APPROVE only
if version-gated APIs are available at the pin, inventory (if used) is correct,
the layer function is reusable, and no scenario survives.

### Phase 5 — Meta-critique (mandatory)

1. What is the most likely way I am being too lenient here?
2. Which version gate or drift path did I not check that I should have?
3. If APPROVE/CONDITIONAL, what would I say to a reviewer who reached REJECT?

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden)

- "Deployment story looks fine" without checking the `go.mod` pin against the
  APIs used.
- Assuming `inventory` exists without confirming dsco ≥ v1.4.0.
- "Best practice is an inventory driver" without naming the drift failure.
- Praise of the artifact or the author.
- Conditional suggestions instead of explicit issue/remediation pairs.

## References

- [references/checklist-deployment.md](references/checklist-deployment.md)
- Shared: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md), [../review-dsco/references/good-reviews.md](../review-dsco/references/good-reviews.md)
- Canonical rules: [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
