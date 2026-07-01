---
name: review-dsco-layers
description: >-
  Adversarial review of the LAYERS lane of dsco config: the Fill target shape
  (**Struct), layer set (duplicate cmdline / duplicate env prefix), layer order
  & precedence (high→low, first-to-supply wins), strict-layer placement,
  *Layers factoring vs inlining, struct-layer pointer freshness, and manual env
  parsing that bypasses dsco. Usually run by the review-dsco orchestrator.
  Default verdict: REJECT.
x-dsco-target: v1.4.0-rc.1
x-bundle-version: 1.4.0-rc.1
---

# review-dsco-layers

Per-aspect reviewer for the **layers** lane: how the `Fill`/`Compute` call is
wired. Gets the target shape, the layer set, their order, and strict placement
right — or overrides go dead, registration fails, or a strict layer errors at
startup.

## Disposition (load-bearing — do not soften)

Default verdict: REJECT. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue found. You are
not evaluating effort, intent, or investment. Treat the artifact as anonymous;
you do not receive the author's identity, history, or production context. If you
received this as an isolated sub-agent with an anonymised artifact, execute the
procedure directly.

## Activation

Applies to any `dsco.Fill` / `inventory.Compute` call-site or `*Layers`
function. Signals: `dsco.Fill(`, `With*Layer(`, `inventory.Compute(`.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no scoring)

1. **Fill target** — `dsco.Fill(&cfg, ...)` where `cfg` is `*Config`, so the
   target is `**Config`. `Fill(cfg, ...)` (single pointer) → `InvalidInputError`.
2. **Layer set** — count cmdline layers (must be ≤1) and env prefixes (no two
   the same); a duplicate fails at registration (`LayerErrors` /
   `CmdlineAlreadyUsedError`).
3. **Layer order** — layers listed high → low priority (cmdline → env →
   providers → struct defaults). The first layer to supply a field wins.
4. **Shadowed overrides** — no struct-default layer placed ABOVE the env/cmdline
   layers meant to override it (that makes the overrides dead).
5. **Strict-layer placement** — for each `WithStrict*Layer`, list every other
   layer that can supply the same field; a strict layer placed after such a
   layer errors with `OverriddenKeyError`.
6. **`*Layers` factoring** — is the layer list a function (`Layers()`), or
   inlined at the `Fill` site where an inventory binary/test can't reuse it?
7. **Struct-layer freshness** — across named variants, each constructor builds a
   fresh struct value (dedup is by pointer address); none returns a shared
   package-level variable.
8. **Manual env parsing** — no `os.Getenv`/`strconv` alongside dsco for config
   dsco could carry.

**Minimum threshold: 6 observations.** If fewer, justify triviality in writing.

Detailed checks: [references/checklist-layers.md](references/checklist-layers.md).

### Phase 2 — Adversarial scenarios

**Minimum 3 scenarios.** Trigger / Propagation (which layer wins, which error
fires) / Symptom / Detectability (registration error, startup error, or silent
dead override?).

### Phase 3 — Scoring

- **BLOCKING** — `Fill(cfg)` single-pointer (`InvalidInputError`, never boots);
  struct defaults first so all overrides are dead; a strict layer guaranteed to
  `OverriddenKeyError` at startup; a shared package-struct across named variants
  that pointer-dedup drops.
- **IMPORTANT** — layers inlined when an inventory binary/test is planned (drift
  risk); a strict layer whose placement is defensible but undocumented.
- **NOTED** — a redundant non-nil target pre-alloc; an out-of-scope category.

Calibration: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md).

### Phase 4 — Verdict

≥1 BLOCKING → REJECT. Only IMPORTANT/NOTED → CONDITIONAL APPROVE. APPROVE only
if the target is `**Struct`, the order is high→low with live overrides, no
duplicate/strict-misplacement, and no scenario survives.

### Phase 5 — Meta-critique (mandatory)

1. What is the most likely way I am being too lenient here?
2. Which layer interaction did I not trace that I should have?
3. If APPROVE/CONDITIONAL, what would I say to a reviewer who reached REJECT?

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden)

- "The layers look right" without tracing first-to-supply order per field.
- Accepting defaults-first order because "it still fills".
- "Best practice is cmdline first" without naming the dead-override failure.
- Praise of the artifact or the author.
- Conditional suggestions instead of explicit issue/remediation pairs.

## References

- [references/checklist-layers.md](references/checklist-layers.md)
- Shared: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md), [../review-dsco/references/good-reviews.md](../review-dsco/references/good-reviews.md)
- Canonical rules: [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
