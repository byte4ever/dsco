# Changelog

All notable changes to the `dsco-claude` bundle. Versions track the dsco
library they target (bundle `vX.Y.Z` → dsco `vX.Y.Z`).

## v1.4.0-rc.1

Initial extraction of the dsco Claude tooling into a dedicated `dsco-claude/`
directory in the dsco repo, kept separate from the library code and installed
on its own.

### Added
- `skills/dsco/`: authoring skill for writing idiomatic dsco config (load-bearing
  rules, design/migrate/troubleshoot/deployment playbooks, and a
  `references/pitfalls.md` anti-pattern catalog). Modeled on the `go` skill.
  Replaces the earlier `dsco-expert` agent.
- `skills/review-dsco/`: the review **orchestrator**, modeled on `review-go`.
  Anonymises the artifact, selects the applicable per-aspect reviewers by
  signal, fans them out concurrently as isolated sub-agents, and arbitrates one
  global verdict (worst-verdict-wins). Holds shared `references/` (severity
  rubric + three worked few-shot reviews).
- `skills/review-dsco-{typing,layers,secrets,validation,deployment}/`: five
  per-aspect reviewers, each built to the reviewer-agent spec — REJECT by
  default, anonymous artifact, isolated sub-agent, phases enumerate → scenarios
  → score → verdict → meta-critique, each with its own `references/checklist`.
- Version targeting: both skills declare the dsco version they target
  (`x-dsco-target: v1.4.0-rc.1`) and check the user's `go.mod` before
  version-gated advice, offering `go get github.com/byte4ever/dsco@v1.4.0-rc.1`
  instead of assuming an API (e.g. `inventory`, min `v1.4.0-rc.1`) exists.
- Feature-minimums table gating `inventory.Compute` and friends at `v1.4.0-rc.1`.
- Cross-platform install/update scripts: `install.sh` (Linux, macOS, WSL, Git
  Bash) and `install.ps1` (Windows PowerShell), with `install` / `update` /
  `uninstall` / `status`, user-global or `--project` scope, symlink by default
  and a copy fallback for filesystems without symlinks.
- "Keeping in sync with dsco" convention: each new library feature ships a
  matching update to both skills in the same change.
