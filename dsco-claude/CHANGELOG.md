# Changelog

All notable changes to the `dsco-claude` bundle. Versions track the dsco
library they target (bundle `vX.Y.Z` → dsco `vX.Y.Z`).

## v1.4.0-rc.1

Initial extraction of the dsco Claude tooling into a dedicated `dsco-claude/`
directory in the dsco repo, kept separate from the library code and installed
on its own.

### Added
- `agents/dsco-expert.md`: the `dsco-expert` sub-agent, relocated out of the
  dsco README and `~/.claude`.
- Version targeting: the agent now declares the dsco version it targets
  (`v1.4.0-rc.1`) and checks the user's `go.mod` before giving version-gated
  advice.
- Upgrade prompting: when requested advice relies on a feature newer than the
  user's pinned dsco, the agent offers
  `go get github.com/byte4ever/dsco@v1.4.0-rc.1` instead of assuming the API.
- Feature-minimums table gating `inventory.Compute` and friends at
  `v1.4.0-rc.1`.
- `skills/` scaffold for future dsco-specific skills.
