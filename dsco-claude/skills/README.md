# dsco skills

Two dsco-specific Claude Code skills, modeled on the `go` / `review-go` pair:

- **[`dsco/`](dsco/)** — authoring skill. Rules, playbooks, and a pitfalls
  reference for writing/designing idiomatic dsco config.
- **[`review-dsco/`](review-dsco/)** — adversarial reviewer for that output.
  REJECT by default; runs the review in an isolated sub-agent; `references/`
  holds the checklist, severity rubric, and worked few-shot reviews.

Any skill added here must follow the bundle's version-targeting convention:

1. Declare the dsco version it targets in its `SKILL.md` frontmatter:
   ```yaml
   x-dsco-target: v1.4.0-rc.1
   x-bundle-version: 1.4.0-rc.1
   ```
2. Gate version-specific advice: before using an API, check the user's pinned
   dsco version (`go.mod`) and offer the upgrade if the feature needs a newer
   version. See the `dsco` skill's *Version targeting* section for the exact
   behavior to mirror.
3. Bump alongside the bundle `VERSION` when the targeted dsco version changes.
