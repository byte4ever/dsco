# dsco-claude

Claude Code tooling for [dsco](https://github.com/byte4ever/dsco) — the
`dsco-expert` agent and any dsco-specific skills.

This directory lives in the dsco repository but is kept **separate from the
library code** and installed on its own. Because it ships with the repo, its
version always matches the dsco tag it was released with (bundle `vX.Y.Z` ==
dsco `vX.Y.Z`).

- **This release targets dsco `v1.4.0-rc.1`** (see [`VERSION`](VERSION)).

## Why a dedicated directory

The agent and skills used to live embedded in the dsco README and in
`~/.claude`. Keeping them here instead means:

- The AI tooling is decoupled from the Go package: you install it with a
  symlink, not with `go get`.
- A single canonical source per artifact, so the README copy and the
  installed copy can't drift.
- The tooling versions in lockstep with the library it targets, guaranteed by
  riding the same repo tags.

## Contents

```
dsco-claude/
  agents/
    dsco-expert.md      # the dsco-expert sub-agent (targets dsco v1.4.0-rc.1)
  skills/               # dsco-specific skills (none yet, ready to hold them)
  VERSION               # bundle version == targeted dsco version
  CHANGELOG.md
  README.md
```

## Version targeting

Every artifact here (the agent and any future skill) declares the dsco version
it targets and follows one rule:

> **Before giving version-gated advice, check the version the user actually
> depends on (their `go.mod`). If a feature needs a newer dsco than they have
> pinned, say so and offer the upgrade instead of assuming the API exists.**

Concretely, the agent will propose:

```bash
go get github.com/byte4ever/dsco@v1.4.0-rc.1
go mod tidy
```

when the requested advice relies on a feature the user's pinned version does
not have.

### Feature minimums (current target)

| Feature / API | Minimum dsco |
|---|---|
| Core (`Fill`, `WithEnvLayer`, `WithCmdlineLayer`, `WithStructLayer`, `WithStrictEnvLayer`, `WithStringValueProvider`, `dsco.R`) | `v1.0.0-beta` |
| `inventory.Compute`, `*Report`, `WriteText` / `WriteJSON` / `WriteYAML` | `v1.4.0-rc.1` |

Skills added later must carry an `x-dsco-target` field in their frontmatter
and the same version-gate behavior. See
[`skills/README.md`](skills/README.md).

## Install

Point Claude Code at the artifacts. The recommended install symlinks the
agent so the repo copy stays the single source of truth. Run from your dsco
checkout:

```bash
# user-global
mkdir -p ~/.claude/agents
ln -sf "$(pwd)/dsco-claude/agents/dsco-expert.md" ~/.claude/agents/dsco-expert.md
```

or, project-local in another project (takes precedence over user-global):

```bash
mkdir -p .claude/agents
ln -sf /abs/path/to/dsco/dsco-claude/agents/dsco-expert.md .claude/agents/dsco-expert.md
```

To pin a specific tool version, check out the matching dsco tag:

```bash
git checkout v1.4.0-rc.1
```

## Releasing

The bundle rides dsco's tags, so a dsco release is a bundle release. Whenever
the targeted dsco version changes or the tooling itself changes:

1. Update [`VERSION`](VERSION) and each artifact's `x-dsco-target` /
   `x-bundle-version`.
2. Update the feature-minimums table if a new API landed.
3. Add a [`CHANGELOG.md`](CHANGELOG.md) entry.
4. Tag the dsco repo as usual (`git tag vX.Y.Z`).

## License

MIT — same as dsco.
