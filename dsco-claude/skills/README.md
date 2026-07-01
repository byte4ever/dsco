# dsco skills

An authoring skill plus a review fleet, modeled on the `go` / `review-go` pair.

- **[`dsco/`](dsco/)** — authoring skill. Rules, playbooks, and a pitfalls
  reference for writing/designing idiomatic dsco config.
- **[`review-dsco/`](review-dsco/)** — the review **orchestrator**. Anonymizes
  the artifact, selects the applicable per-aspect reviewers by signal, fans them
  out concurrently as isolated sub-agents, and arbitrates one global verdict
  (worst-verdict-wins). Holds the shared `references/` (severity rubric + worked
  reviews).
- **Per-aspect reviewers** (each REJECT by default, own `references/checklist`):
  - [`review-dsco-typing/`](review-dsco-typing/) — pointer fields, yaml tags.
  - [`review-dsco-layers/`](review-dsco-layers/) — Fill target, order,
    precedence, strict placement, dedup, factoring.
  - [`review-dsco-secrets/`](review-dsco-secrets/) — env prefixes, secret
    routing, cmdline surface.
  - [`review-dsco-validation/`](review-dsco-validation/) — Validate() coverage,
    required fields, error handling.
  - [`review-dsco-deployment/`](review-dsco-deployment/) — inventory, *Layers
    export, version targeting, third-party embedding.

For one aspect, invoke its reviewer directly; for a full review, invoke
`review-dsco` and let it fan out.

Any skill added here must follow the bundle's version-targeting convention:

1. Declare the dsco version it targets in its `SKILL.md` frontmatter:
   ```yaml
   x-dsco-target: v1.4.0-rc.1
   x-bundle-version: 1.4.0-rc.1
   ```
2. Gate version-specific advice: before using an API, check the user's pinned
   dsco version (`go.mod`) and offer the upgrade if the feature needs a newer
   version. See the `dsco` skill's *Version targeting* section.
3. Bump alongside the bundle `VERSION` when the targeted dsco version changes.
