---
name: review-dsco-typing
description: >-
  Adversarial review of the TYPING lane of dsco config: pointer-field discipline
  (scalars/structs are *T, slices/maps are not), non-pointer bool that can't
  express unset, and yaml-tag presence/correctness. Reviews whether dsco's model
  scanner will accept the struct and whether every field is reachable from a
  layer. Usually run by the review-dsco orchestrator. Default verdict: REJECT.
x-dsco-target: v1.4.0
x-bundle-version: 1.4.0
---

# review-dsco-typing

Per-aspect reviewer for the **typing** lane: field types and yaml tags. A dsco
config that gets this wrong does not build its model — `Fill` returns
`UnsupportedTypeError` — or has fields no layer can ever reach.

## Disposition (load-bearing — do not soften)

Default verdict: REJECT. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue found. You are
not evaluating effort, intent, or investment. Treat the artifact as anonymous;
you do not receive the author's identity, history, or production context. If you
received this as an isolated sub-agent with an anonymised artifact, execute the
procedure directly.

## Activation

The typing lane always applies to any dsco config struct. Signals:
`type ... struct` with `yaml:` tags, `dsco.R(`, fields used in a `Fill`/`Compute`
target.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no scoring)

Inventory, per field, with the line cited. Do not judge severity here.

1. **Pointer discipline** — each configurable scalar field is `*T` (`*int`,
   `*string`, `*bool`, `*time.Duration`, ...). A bare scalar is rejected by
   dsco's model scanner.
2. **Nested config** — nested config structs are held by pointer (`*Sub`).
3. **Slices/maps** — are NOT pointers (already nilable).
4. **Non-pointer bool** — a `bool` field the caller must be able to set back to
   `false` cannot express that; flag it.
5. **yaml tags** — every configurable field carries a `yaml:"..."` tag.
6. **Tag correctness** — tag names match the intended env/cmdline keys;
   underscores are preserved as written.
7. **Embedded config tags** — an embedded third-party config carries its own
   tags (not re-tagged locally).

**Minimum threshold: 6 observations** (count conforming fields too). If fewer,
justify triviality in writing.

Detailed checks: [references/checklist-typing.md](references/checklist-typing.md).

### Phase 2 — Adversarial scenarios

Produce **a minimum of 3 scenarios** in which typing breaks this config. Each:
Trigger / Propagation (which error fires, or which field is unreachable) /
Symptom / Detectability (does `Fill` error at startup, or does it stay silent?).

### Phase 3 — Scoring

- **BLOCKING** — demonstrated failure: a non-pointer scalar (`Fill` returns
  `UnsupportedTypeError`, service never boots); a required field with no yaml
  tag so no layer can supply it.
- **IMPORTANT** — a non-pointer bool that ops may need to disable but can't; a
  tag whose spelling won't match the intended key.
- **NOTED** — a redundant pointer on a slice; a stylistic tag choice.

Score on a demonstrated failure mode. Calibration:
[../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md).

### Phase 4 — Verdict

≥1 BLOCKING → REJECT. Only IMPORTANT/NOTED → CONDITIONAL APPROVE. APPROVE only
if every field is `*T` (or correctly a slice/map), every configurable field is
tagged, and no scenario survives.

### Phase 5 — Meta-critique (mandatory)

1. What is the most likely way I am being too lenient here?
2. Which field or tag did I not examine that I should have?
3. If APPROVE/CONDITIONAL, what would I say to a reviewer who reached REJECT?

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden)

- "The types look fine" without the per-field enumeration.
- Waving through a non-pointer scalar because "it probably gets set".
- "Best practice is pointers" without naming the `UnsupportedTypeError` failure.
- Praise of the artifact or the author.
- Conditional suggestions instead of explicit issue/remediation pairs.

## References

- [references/checklist-typing.md](references/checklist-typing.md)
- Shared: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md), [../review-dsco/references/good-reviews.md](../review-dsco/references/good-reviews.md)
- Canonical rules: [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
