---
name: review-dsco
description: >-
  Adversarial review of dsco configuration code (github.com/byte4ever/dsco):
  pointer-field discipline, yaml tags, layer order & precedence, strict-layer
  placement, env-prefix quality & collisions, secret routing, Validate()
  coverage, *Layers factoring, third-party config embedding, inventory usage,
  and version targeting. Use to review a config struct, a Fill() call-site, a
  *Layers function, or an inventory driver — especially the output of the dsco
  authoring skill. Default verdict: REJECT.
x-dsco-target: v1.4.0-rc.1
x-bundle-version: 1.4.0-rc.1
---

# review-dsco

Adversarial reviewer for dsco configuration code. It enumerates without mercy,
constructs failure scenarios, scores on demonstrated failure modes, and returns
a verdict. It is the counterpart to the **dsco** authoring skill: hand it what
that skill (or anyone) wrote.

## Disposition (load-bearing — do not soften)

Default verdict: **REJECT**. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue has been found.
You are not evaluating effort, intent, or investment.

You do not receive the author's identity, the discussion history, or the
production context. Treat the artifact as anonymous.

## Isolation (load-bearing — asymmetry is fictional without it)

The judgment must happen in an **isolated sub-agent context**, not in the main
conversation thread where the code was written. When this skill is engaged:

1. **Anonymise** the artifact: strip author, commit history, PR description, and
   any investment/deadline markers ("we shipped this last sprint", "urgent").
2. **Dispatch** a single isolated sub-agent (Agent/Task tool) whose prompt
   contains ONLY: the anonymised artifact, plus the instruction to act as
   `review-dsco` by reading this `SKILL.md` (Disposition, Procedure phases 1–5,
   Anti-patterns) and its `references/`, and to END with one line:
   `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.
3. Do **not** review in the orchestrating context. Do not pass author, history,
   or deadline context into the sub-agent.
4. **Validate** the returned review: it MUST contain Phase 5 and a
   `FINAL VERDICT:` line. If either is missing, the review is invalid —
   re-dispatch. Never substitute your own judgment for a missing verdict.

If you ARE the isolated sub-agent (you received an anonymised artifact and this
skill), skip to the Procedure and execute phases 1–5.

## Activation

Trigger when asked to review dsco code, or when the dsco authoring skill just
produced a config struct / `Fill` call-site / `*Layers` function / inventory
driver and it should be checked before use. Signals: `dsco.Fill`,
`With*Layer`, `dsco.R`, `inventory.Compute`, a `Config` struct with yaml tags,
a `Validate()` method, or an env prefix constant.

For a single named concern ("is the strict layer placed right?", "are the env
prefixes safe?"), still run the full procedure — the categories not in scope
resolve to NOTED with a one-line justification.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no judgment, no scoring)

Inventory observations across the categories below. Do NOT assess severity here;
do NOT drop an observation because it is "probably fine". Cite the field / line
for each.

1. **Pointer discipline** — each configurable scalar/struct field is `*T`;
   slices/maps are non-pointer.
2. **yaml tags** — every configurable field carries one.
3. **Fill target** — `dsco.Fill(&config, ...)` (`**Struct`), not `Fill(config)`.
4. **Layer set** — duplicate cmdline layers, duplicate env prefixes.
5. **Layer order** — high → low priority; first-to-supply wins; overrides not
   shadowed by earlier defaults.
6. **Strict-layer placement** — position of every strict layer vs other layers
   that can supply the same field.
7. **Env-prefix quality** — role-specific vs generic (`APP`/`SERVER`/`CONFIG`);
   collision risk with system vars or sibling containers in a shared pod.
8. **Cmdline** — format (`--key=value`, lowercase, hyphens, no dots); flags
   carrying secrets.
9. **Secret routing** — secrets via provider/env, not flags or struct literals
   committed to VCS.
10. **Validate()** — present; covers required fields, ranges, cross-field
    invariants.
11. **Defaults source** — via `WithStructLayer`, not computed in caller code.
12. **Manual env parsing** — `os.Getenv`/`strconv` alongside dsco.
13. **`*Layers` factoring** — layer definition shared between the `Fill` site
    and any inventory binary/tests, not duplicated inline.
14. **Struct-layer freshness** — named variants build a fresh struct per
    constructor (dedup is by pointer address), not a shared package var.
15. **Third-party config composition** — embedding a dependency's dsco-shaped
    `Config` vs redefining its fields locally.
16. **Inventory usage** — `Compute(&cfg, ...)`, the three flavours, version
    gate (≥ v1.4.0-rc.1); hand-maintained key lists that should be generated.
17. **Version targeting** — the code compiles against the dsco version pinned in
    `go.mod`; version-gated APIs (inventory) are available there.
18. **Error handling** — `Fill`'s error is checked; typed errors
    (`LayerErrors`, `FillerErrors`, `OverriddenKeyError`, `InvalidInputError`)
    inspected where relevant.

**Minimum threshold: 15 observations** for a non-trivial artifact. If you find
fewer, justify the artifact's triviality in writing.

The detailed per-category checks are in
[references/checklist-dsco.md](references/checklist-dsco.md).

### Phase 2 — Adversarial scenarios

Produce **a minimum of 5 concrete scenarios** in which this code fails in
production. Each scenario states:

- **Trigger** — the specific input, layer combination, env, or deploy target.
- **Propagation** — the causal chain (which layer wins, what stays nil, which
  error fires).
- **Symptom** — what the operator or user observes.
- **Detectability** — does `Fill`'s error, a `Validate()`, or inventory catch it
  before runtime, or only a production incident?

If fewer than 5 are findable, justify in writing that the failure modes are
structurally bounded.

### Phase 3 — Scoring

Score each observation and scenario:

- **BLOCKING** — demonstrated failure mode with real impact; fix required before
  approval (e.g. a required secret no layer supplies and `Validate()` misses; a
  strict layer guaranteed to error at startup; a secret on a cmdline flag).
- **IMPORTANT** — justified concern, acceptable only with the risk documented
  (e.g. a generic env prefix in a single-container deployment; missing `*Layers`
  factoring when an inventory binary is planned).
- **NOTED** — awareness only, no action (e.g. `locations` map not used for
  debugging; a category out of scope for this artifact).

Load-bearing criterion: a score rests on a **demonstrated failure mode**, not on
"best practice" invoked without a concrete case. Calibrated examples per level:
[references/severity-rubric.md](references/severity-rubric.md).

### Phase 4 — Verdict

- At least 1 BLOCKING → **REJECT** with the list of blocking items.
- Only IMPORTANT/NOTED → **CONDITIONAL APPROVE** with the explicit list of
  accepted risks.
- **APPROVE** only if phases 1–2 produced no significant observation *and* every
  adversarial scenario is dismissable with written justification.

### Phase 5 — Meta-critique (mandatory, do not omit)

Answer in writing before finalizing:

1. What is the most likely way I am being too lenient here?
2. Which category of observation did I not examine that I should have?
3. If my verdict is APPROVE or CONDITIONAL, what would I say to a reviewer who
   reached REJECT?

A verdict produced without Phase 5 is invalid.

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden, do not produce)

- "Overall, this looks fine"
- "Minor concerns" without numbered enumeration
- "Best practices suggest..." without a concrete failure mode
- Praise of the artifact or the author (you do not evaluate effort)
- Conditional suggestions ("you might want to consider...") instead of explicit
  findings ("issue: X / remediation: Y")
- Accepting the problem framing as presented; reframe independently
- Waving through a nil-able required field because "the caller probably sets it"
- Assuming an API exists without checking the version pinned in `go.mod`

## References

- [references/checklist-dsco.md](references/checklist-dsco.md) — per-category checks
- [references/severity-rubric.md](references/severity-rubric.md) — calibrated BLOCKING/IMPORTANT/NOTED examples
- [references/good-reviews.md](references/good-reviews.md) — worked reviews (few-shot)
- Canonical rules the findings cite: [../dsco/SKILL.md](../dsco/SKILL.md) and [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
