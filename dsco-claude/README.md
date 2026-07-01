# dsco-claude

Claude Code tooling for [dsco](https://github.com/byte4ever/dsco) — two skills:
`dsco` (write idiomatic dsco config) and `review-dsco` (review it
adversarially).

This directory lives in the dsco repository but is kept **separate from the
library code** and installed on its own. Because it ships with the repo, its
version always matches the dsco tag it was released with (bundle `vX.Y.Z` ==
dsco `vX.Y.Z`).

- **This release targets dsco `v1.4.0`** (see [`VERSION`](VERSION)).

## Skills

- **`dsco`** — the authoring skill. Distils the library's best practices and
  pitfalls into rules and playbooks for writing/designing dsco config. Modeled
  on the `go` skill. It does not present code until it passes review: every
  artifact it produces runs through `review-dsco` in a bounded
  correction → re-validation loop (up to 3 cycles on a REJECT, then it escalates
  the unresolved BLOCKING findings to the user to arbitrate).
- **`review-dsco`** — the review **orchestrator** for what `dsco` (or anyone)
  writes. It anonymises the artifact, selects the applicable per-aspect
  reviewers by signal, fans them out concurrently as isolated sub-agents, and
  arbitrates one global verdict (worst-verdict-wins). Modeled on `review-go`.
- **`review-dsco-{typing,layers,secrets,validation,deployment}`** — the five
  per-aspect reviewers. Each is REJECT by default, anonymous artifact, isolated
  sub-agent, phases enumerate → scenarios → score → verdict → meta-critique,
  with its own `references/checklist`. Built to the team's reviewer-agent spec.

## Why a dedicated directory

The tooling used to live embedded in the dsco README and in `~/.claude` (as a
single `dsco-expert` agent, now split into the two skills above). Keeping it
here instead means:

- The AI tooling is decoupled from the Go package: you install it with a
  symlink, not with `go get`.
- A single canonical source per artifact, so the README copy and the
  installed copy can't drift.
- The tooling versions in lockstep with the library it targets, guaranteed by
  riding the same repo tags.

## Contents

```
dsco-claude/
  skills/
    dsco/                     # authoring skill (SKILL.md + references/pitfalls.md)
    review-dsco/              # review orchestrator (+ shared references/)
    review-dsco-typing/       # per-aspect reviewer (SKILL.md + references/)
    review-dsco-layers/       #   "
    review-dsco-secrets/      #   "
    review-dsco-validation/   #   "
    review-dsco-deployment/   #   "
  bootstrap.sh          # one-line curl/wget install (POSIX)
  bootstrap.ps1         # one-line install (Windows PowerShell)
  install.sh            # installer/updater (Linux, macOS, WSL, Git Bash)
  install.ps1           # installer/updater (Windows PowerShell)
  VERSION               # bundle version == targeted dsco version
  CHANGELOG.md
  README.md
```

## Version targeting

Both skills declare the dsco version they target (`x-dsco-target`) and follow
one rule:

> **Before giving version-gated advice, check the version the user actually
> depends on (their `go.mod`). If a feature needs a newer dsco than they have
> pinned, say so and offer the upgrade instead of assuming the API exists.**

Concretely, the skills will propose:

```bash
go get github.com/byte4ever/dsco@v1.4.0
go mod tidy
```

when the requested advice relies on a feature the user's pinned version does
not have.

### Feature minimums (current target)

| Feature / API | Minimum dsco |
|---|---|
| Core (`Fill`, `WithEnvLayer`, `WithCmdlineLayer`, `WithStructLayer`, `WithStrictEnvLayer`, `WithStringValueProvider`, `dsco.R`) | `v1.0.0-beta` |
| `inventory.Compute`, `*Report`, `WriteText` / `WriteJSON` / `WriteYAML` | `v1.4.0` |

Skills added later must carry an `x-dsco-target` field in their frontmatter
and the same version-gate behavior. See
[`skills/README.md`](skills/README.md).

## Quick install (curl / wget)

No checkout needed. The bootstrap downloads the bundle and installs the skills
in one line.

Linux, macOS, WSL, Git Bash:

```bash
curl -fsSL https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.sh | sh
# or with wget:
wget -qO- https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.ps1 | iex
```

It places the bundle in `~/.dsco-claude` (override with `DSCO_CLAUDE_HOME`) and
symlinks the skills into `~/.claude/skills`. Re-run it to update.

Pin a version or pass install options (note the `-s --`):

```bash
curl -fsSL <bootstrap.sh url> | sh -s -- --ref v1.4.0 --copy
# env form: curl -fsSL <bootstrap.sh url> | DSCO_CLAUDE_REF=v1.4.0 sh
```

The default ref is `master`; pin a release tag once one ships the bundle. Add
`--copy` (`-Copy`) for filesystems without symlinks, or `--project .` to install
into the current project's `.claude`.

## Install from a checkout

If you already have the dsco repo, run the bundled installer directly. It
symlinks the skills (`dsco`, `review-dsco`, ...) into Claude Code so the repo
copy stays the single source of truth, and falls back to copying on filesystems
without symlinks. Run it from your dsco checkout.

Linux, macOS, WSL, Git Bash:

```bash
dsco-claude/install.sh              # install into ~/.claude
dsco-claude/install.sh update       # refresh after checking out a new version
dsco-claude/install.sh status       # show what is installed
dsco-claude/install.sh --project .  # install into ./.claude of another project
dsco-claude/install.sh uninstall
```

Windows (PowerShell):

```powershell
dsco-claude\install.ps1             # install into ~\.claude
dsco-claude\install.ps1 update
dsco-claude\install.ps1 status
dsco-claude\install.ps1 -Project . # install into .\.claude
dsco-claude\install.ps1 uninstall
```

Add `--copy` (`-Copy` on PowerShell) to copy the files instead of symlinking.
A project-local install takes precedence over the user-global one.

To pin a specific tool version, check out the matching dsco tag before running
the installer:

```bash
git checkout v1.4.0
```

## Keeping in sync with dsco

Every new dsco feature should land a matching update here, in the same change:

- Extend the `dsco` skill (a new load-bearing rule or a `references/pitfalls.md`
  entry) so the authoring guidance covers the feature.
- Extend `review-dsco` (a Phase-1 category, a `checklist-dsco.md` check, a
  `severity-rubric.md` example) so the reviewer catches its misuse.
- Update the feature-minimums table with the version that introduced it.
- Add a `CHANGELOG.md` entry.

Treat the tooling as part of shipping the feature, not an afterthought.

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
