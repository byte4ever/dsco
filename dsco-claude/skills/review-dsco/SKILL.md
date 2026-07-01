---
name: review-dsco
description: >-
  Orchestrated multi-aspect adversarial review of dsco configuration code
  (github.com/byte4ever/dsco). Selects the applicable per-aspect reviewers
  (typing, layers, secrets, validation, deployment), fans them out
  CONCURRENTLY as isolated sub-agents, collects each verdict, and arbitrates a
  single global verdict (worst-verdict-wins). Use to comprehensively review a
  config struct, a Fill() call-site, a *Layers function, or an inventory driver
  — especially the output of the dsco authoring skill. Default verdict: REJECT.
x-dsco-target: v1.4.0-rc.1
x-bundle-version: 1.4.0-rc.1
---

# review-dsco

Orchestrator for the dsco reviewer fleet. This skill is NOT a per-aspect
reviewer: it does not itself enumerate observations, run scenarios, or judge the
artifact. It SELECTS the applicable reviewers, FANS them out concurrently as
isolated sub-agents, and ARBITRATES a single global verdict.

## Disposition (load-bearing — do not soften)

Default GLOBAL verdict: REJECT. The aggregate APPROVE must be EARNED by EVERY
applicable reviewer independently approving — it is never assumed and never
averaged into existence.

You orchestrate; you do not re-judge the artifact yourself, and you do not
override, downgrade, or drop any reviewer's finding. A reviewer owns its lane;
its verdict and its BLOCKING/IMPORTANT/NOTED items pass through to the aggregate
verbatim and attributed.

The artifact is anonymous. Strip author, commit history, and any
investment/deadline markers BEFORE fan-out. You do not receive (and do not pass
on) the author's identity, the discussion history, or the production context.

You are not evaluating effort, intent, or investment.

## Activation

Trigger when asked to review dsco code COMPREHENSIVELY — "review this dsco
config", "is this Fill setup ready?", or right after the `dsco` authoring skill
produced a config struct / `Fill` call-site / `*Layers` function / inventory
driver.

If the user asks for ONE aspect, do NOT orchestrate — invoke that single
reviewer directly:

| User asks for… | Invoke directly |
|---|---|
| pointer fields / yaml tags / field types | review-dsco-typing |
| layer order / precedence / strict / Fill target / dedup | review-dsco-layers |
| env prefixes / secrets / secret leaks | review-dsco-secrets |
| Validate() / required fields / error handling | review-dsco-validation |
| inventory / deploy / *Layers export / version / embedding | review-dsco-deployment |

## Procedure (mandatory, in order)

1. **Identify & anonymise.** Resolve the artifact(s) under review. Strip author,
   commit history, and investment/deadline markers. Everything downstream
   operates on the anonymised artifact.

2. **Select reviewers** per the Selection table below (by signal). For EACH
   reviewer, record WHY it was selected (which signal fired) or WHY it was
   skipped (which signal was absent). Coverage is never silently dropped.

3. **Fan out CONCURRENTLY.** Dispatch the selected reviewers as PARALLEL
   isolated sub-agents — ALL in ONE batch (multiple Agent/Task calls in a single
   message), one sub-agent per aspect. Each sub-agent's prompt contains ONLY:
   - the anonymised artifact (code / diff);
   - the instruction to act as `review-dsco-<aspect>` by reading
     `../review-dsco-<aspect>/SKILL.md` (its Disposition, Procedure phases 1–5,
     Anti-patterns) and its `references/`;
   - the demand to execute phases 1–5 in order and END with one line:
     `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

   Pass NO author, NO history, NO deadline/investment context. Do not review in
   this orchestrating context — judgment happens inside the sub-agents.

4. **Collect & validate.** Each returned review MUST contain Phase 5
   (meta-critique) AND a `FINAL VERDICT:` line. A review missing either is
   INVALID — re-dispatch that single reviewer. Do not substitute your own
   judgment for a missing one.

5. **Arbitrate** per the Arbitration rules below (worst-verdict-wins).

6. **Produce the consolidated report** per the Report format below.

## Selection (by signal)

Run the greps below over the anonymised artifact, then select by this table.
Record, per reviewer, the signal that fired (selected) or was absent (skipped).

| Signal (in the artifact) | Reviewers added |
|---|---|
| ANY dsco config struct / `dsco.Fill` / `inventory.Compute` (always) | review-dsco-typing + review-dsco-layers + review-dsco-validation |
| an env prefix, a cmdline layer, or a secret-looking field | + review-dsco-secrets |
| `inventory`, a `*Layers` function, an embedded 3rd-party config, or a version-gated API | + review-dsco-deployment |

### Reproducible signal greps

```
# secrets lane (any match → add review-dsco-secrets)
grep -nE 'WithEnvLayer|WithStrictEnvLayer|WithCmdlineLayer' artifact       # env/cmdline surface
grep -niE 'password|secret|token|api[_-]?key|credential|dsn|database_url' artifact  # secret-looking fields

# deployment lane (any match → add review-dsco-deployment)
grep -nE 'inventory\.|WriteJSON|WriteText|WriteYAML' artifact              # inventory usage
grep -nE 'func .*Layers\(' artifact                                        # *Layers function
grep -nE '\*[a-z0-9]+\.Config' artifact                                    # embedded 3rd-party config
```

### Inclusion rule under uncertainty

When unsure whether a signal applies, INCLUDE the reviewer. A reviewer with
nothing in its lane returns a "nothing in scope (trivial)" verdict — one cheap
sub-agent, no downstream change. A skipped lane that mattered ships a defect.
Bias toward inclusion. typing, layers, and validation always run.

## Arbitration (worst-verdict-wins — load-bearing)

- **GLOBAL REJECT** if ANY applicable reviewer returns REJECT (≥1 BLOCKING in
  any lane). List every blocking item, attributed to its reviewer.
- **GLOBAL CONDITIONAL APPROVE** if no reviewer REJECTs but ≥1 returns
  CONDITIONAL APPROVE. Accepted-risks list is the UNION of every reviewer's
  IMPORTANT findings, attributed.
- **GLOBAL APPROVE** only if EVERY applicable reviewer returns APPROVE.
- Never let a NOTED or "out of scope" downgrade a BLOCKING. Never average — the
  worst single verdict governs. A split (one REJECT, four APPROVE) is a GLOBAL
  REJECT.

## Report format

```
# Multi-aspect dsco review — <artifact>

## Reviewers run
<aspect> — selected because <signal>   (one line each)
Skipped: <aspect> — <why>

## Per-aspect verdicts
| Aspect | Verdict | BLOCKING | IMPORTANT |
| ... |

## Consolidated findings (by severity, then aspect)
### BLOCKING
- [<aspect>] <finding> — remediation: <...>
### IMPORTANT
- [<aspect>] <finding>
### NOTED
- [<aspect>] <finding>

## GLOBAL VERDICT: <REJECT | CONDITIONAL APPROVE | APPROVE>

## Aggregate meta-check (mandatory)
1. Did I drop, merge-away, or downgrade any reviewer's BLOCKING/IMPORTANT? (must be "no")
2. Did every applicable reviewer run and return a Phase-5 verdict?
3. Is a global APPROVE/CONDITIONAL truly earned, or am I reconciling a split into something milder?
```

## Anti-patterns (forbidden — do not soften the aggregate)

- Overriding or downgrading a reviewer's REJECT/BLOCKING to reconcile a split.
- Averaging verdicts, or declaring APPROVE while any CONDITIONAL/REJECT stands.
- Dropping a reviewer that "probably doesn't apply" without recording the skip.
- Re-judging the artifact yourself instead of relaying the reviewers' findings.
- Praise of the artifact or the author (you do not evaluate effort).
- Passing author / history / deadline context into a sub-agent prompt.
- Proceeding to arbitration with a review missing Phase 5 or its FINAL VERDICT.

## References

- The 5 per-aspect reviewers:
  [review-dsco-typing](../review-dsco-typing/SKILL.md),
  [review-dsco-layers](../review-dsco-layers/SKILL.md),
  [review-dsco-secrets](../review-dsco-secrets/SKILL.md),
  [review-dsco-validation](../review-dsco-validation/SKILL.md),
  [review-dsco-deployment](../review-dsco-deployment/SKILL.md)
- Shared calibration: [references/severity-rubric.md](references/severity-rubric.md), [references/good-reviews.md](references/good-reviews.md)
- Canonical rules the reviewers cite: [../dsco/SKILL.md](../dsco/SKILL.md) and [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
