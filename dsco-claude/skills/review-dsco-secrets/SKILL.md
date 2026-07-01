---
name: review-dsco-secrets
description: >-
  Adversarial review of the SECRETS/ENV lane of dsco config: env-prefix quality
  and collision risk (generic APP/SERVER vs role-specific, shared-pod and system
  var clashes), secret routing (via env/provider, never cmdline flags or VCS
  literals, left nil in defaults), and cmdline format. Usually run by the
  review-dsco orchestrator. Default verdict: REJECT.
x-dsco-target: v1.4.0-rc.1
x-bundle-version: 1.4.0-rc.1
---

# review-dsco-secrets

Per-aspect reviewer for the **secrets/env** lane: the env-var namespace and how
credentials enter the config. Gets this wrong and a prefix collides across
containers, or a secret leaks through a flag or a committed literal.

## Disposition (load-bearing — do not soften)

Default verdict: REJECT. APPROVE must be justified by explicit demonstration
that the procedure below has been followed and no blocking issue found. You are
not evaluating effort, intent, or investment. Treat the artifact as anonymous;
you do not receive the author's identity, history, or production context. If you
received this as an isolated sub-agent with an anonymised artifact, execute the
procedure directly.

## Activation

Applies when the config has an env layer, a cmdline layer, or a secret-looking
field. Signals: `WithEnvLayer`/`WithStrictEnvLayer`, `WithCmdlineLayer`, and
field names like `password`, `secret`, `token`, `api_key`, `credential`, `dsn`,
`database_url`.

## Procedure (mandatory, in order)

### Phase 1 — Enumeration (no scoring)

1. **Env-prefix quality** — role-specific (`ORDERAPI`, `PAYMENTWORKER`) vs
   generic (`APP`, `SERVER`, `CONFIG`, `SVC`, `SERVICE`).
2. **Sibling collision** — in a multi-container pod, could the prefix collide
   with a sibling container's vars?
3. **System collision** — could the prefix clash with system vars (`PATH`,
   `HOME`, `HTTP_PROXY`) or platform-injected `PREFIX_*`?
4. **Secret inventory** — list every secret-looking field and its intended
   channel.
5. **Secret on a flag** — with a cmdline layer present, each secret field is
   reachable as `--field=...`; flag values land in `ps`, `/proc`, shell history.
6. **VCS literal** — no secret literal is baked into a `WithStructLayer` and
   committed.
7. **Nil-in-defaults** — secret fields are left nil in the defaults layer so a
   higher layer MUST supply them.
8. **Provider use** — for high-value secrets, a `WithStringValueProvider`
   (vault/secrets) rather than plain env, where warranted.
9. **Cmdline format** — `--key=value`, lowercase, hyphens for nesting, no dots.

**Minimum threshold: 6 observations.** If fewer, justify triviality in writing.

Detailed checks: [references/checklist-secrets.md](references/checklist-secrets.md).

### Phase 2 — Adversarial scenarios

**Minimum 3 scenarios.** Trigger / Propagation / Symptom / Detectability. Cover
at least: a prefix collision, and a secret entering via a flag.

### Phase 3 — Scoring

- **BLOCKING** — a secret literal committed to VCS; a generic prefix with a
  concrete demonstrated same-variable collision in a shared pod.
- **IMPORTANT** — a generic prefix with latent collision risk (no concrete
  sibling yet); a cmdline channel left open to a secret field with no provider.
- **NOTED** — a cmdline layer beside a secret env field where env is the clear
  intended path (routing hazard worth stating, not a defect by itself);
  stylistic prefix choice.

Calibration: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md).

### Phase 4 — Verdict

≥1 BLOCKING → REJECT. Only IMPORTANT/NOTED → CONDITIONAL APPROVE. APPROVE only
if the prefix is role-specific, no secret is committed or flag-reachable in a
way that matters, and no scenario survives.

### Phase 5 — Meta-critique (mandatory)

1. What is the most likely way I am being too lenient here?
2. Which secret channel did I not trace that I should have?
3. If APPROVE/CONDITIONAL, what would I say to a reviewer who reached REJECT?

End with: `FINAL VERDICT: <REJECT|CONDITIONAL APPROVE|APPROVE>`.

## Anti-patterns (forbidden)

- "Secrets handling looks fine" without listing each secret field and its channel.
- Accepting a generic prefix because "it works today".
- "Best practice is a role prefix" without naming the collision failure.
- Praise of the artifact or the author.
- Conditional suggestions instead of explicit issue/remediation pairs.

## References

- [references/checklist-secrets.md](references/checklist-secrets.md)
- Shared: [../review-dsco/references/severity-rubric.md](../review-dsco/references/severity-rubric.md), [../review-dsco/references/good-reviews.md](../review-dsco/references/good-reviews.md)
- Canonical rules: [../dsco/references/pitfalls.md](../dsco/references/pitfalls.md)
