---
name: review-dsco-validation
description: >-
  Adversarial review of the VALIDATION lane of dsco config: presence and
  coverage of a Validate() method (required fields, ranges, cross-field
  invariants), whether Fill success is wrongly treated as "config valid",
  defaults sourced via WithStructLayer (attributable), and error handling
  (Fill's error checked, typed errors inspected). Usually run by the review-dsco
  orchestrator. Default verdict: REJECT.
x-dsco-target: v1.4.0
x-bundle-version: 1.4.0
---

# review-dsco-validation

Per-aspect reviewer for the **validation** lane: does the config fail closed?
dsco fills; it does not validate. A required nil-able field with no guard, or a
`Validate()` that only checks the easy cases, ships broken config to runtime.

## Disposition (load-bearing — do not soften)

Default verdict: REJECT. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue found. You are
not evaluating effort, intent, or investment. Treat the artifact as anonymous;
you do not receive the author's identity, history, or production context. If you
received this as an isolated sub-agent with an anonymised artifact, execute the
procedure directly.

## Activation

Applies to any dsco config. Signals: a `Validate()` method (or its absence), a
`Fill` call whose error is/ isn't checked, required-looking fields with no
default.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no scoring)

1. **Validate() presence** — a `Validate()` method exists and is called after
   `Fill`.
2. **Required coverage** — every required field (nil-able, no default) is
   enforced by `Validate()`.
3. **Range/format coverage** — value ranges (port 1–65535), positive durations,
   parseable addresses, and cross-field invariants are checked, not just
   non-nil.
4. **Fill≠valid** — `Fill` success is not treated as "config is valid"; `Fill`
   only fills.
5. **Unit traps** — numeric/duration fields where a plausible operator input
   (`SHUTDOWN_WAIT=15` → 15ns) passes a naive `> 0` check.
6. **Defaults source** — defaults live in `WithStructLayer`, not caller-computed
   (so origin shows in the location map).
7. **Error handling** — `Fill`'s returned error is checked.
8. **Typed errors** — where relevant, typed errors (`LayerErrors`,
   `FillerErrors`, `OverriddenKeyError`, `InvalidInputError`) are inspected
   (`errors.As`).

**Minimum threshold: 6 observations.** If fewer, justify triviality in writing.

Detailed checks: [references/checklist-validation.md](references/checklist-validation.md).

### Phase 2 — Adversarial scenarios

**Minimum 3 scenarios.** Trigger / Propagation / Symptom / Detectability. Cover
at least: a required field left unset, and a value that passes a naive check but
is wrong (unit trap or out-of-range).

### Phase 3 — Scoring

- **BLOCKING** — a required nil-able field with no default and no `Validate()`
  guard (nil-deref / empty-credential in production); no `Validate()` at all
  where required fields exist.
- **IMPORTANT** — a `Validate()` present but incomplete (misses a bound like the
  duration unit trap, or a required field); `Fill` error unchecked.
- **NOTED** — `locations` map discarded; typed-error inspection omitted where
  there's no strict layer to inspect.

Calibration: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md).

### Phase 4 — Verdict

≥1 BLOCKING → REJECT. Only IMPORTANT/NOTED → CONDITIONAL APPROVE. APPROVE only
if every required field fails closed, ranges/units are bounded, and `Fill`'s
error is handled — with no surviving scenario.

### Phase 5 — Meta-critique (mandatory)

1. What is the most likely way I am being too lenient here?
2. Which required field or value bound did I not check that I should have?
3. If APPROVE/CONDITIONAL, what would I say to a reviewer who reached REJECT?

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden)

- "Validation looks adequate" without listing each required field and its guard.
- Treating a present-but-thin `Validate()` as complete.
- "Best practice is to validate" without naming the fail-open runtime symptom.
- Praise of the artifact or the author.
- Conditional suggestions instead of explicit issue/remediation pairs.

## References

- [references/checklist-validation.md](references/checklist-validation.md)
- Shared: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md), [../review-dsco/references/good-reviews.md](../review-dsco/references/good-reviews.md)
- Canonical rules: [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
