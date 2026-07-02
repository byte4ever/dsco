# Changelog

All notable changes to dsco are documented here. This project follows
[semantic versioning](https://semver.org); the format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [v1.4.0] - 2026-07-01

First stable 1.4.0 release.

### Added
- **`inventory` sub-package.** Statically enumerate every config key a layered
  setup expects, with no I/O: `inventory.Compute(&cfg, layers...)` returns a
  `*Report` listing each leaf path, its Go type, the canonical key of the first
  string-based layer that can supply it, and whether a struct layer bakes in a
  default. Text, JSON, and YAML output, plus a preflight check that exits
  non-zero when a required key has no default. Runnable examples under
  [`examples/inventory/`](examples/inventory/).
- **`dsco-claude/` tooling bundle.** A Claude Code bundle shipped with the repo
  and versioned in lockstep with the library: a `dsco` authoring skill and a
  `review-dsco` adversarial review fleet, with cross-platform install scripts
  and a one-line `curl` / `wget` / PowerShell bootstrap. See
  [`dsco-claude/README.md`](dsco-claude/README.md).

### Changed
- README and QUICKSTART (English and French) rewritten in a plain technical
  register; the French docs are faithful translations of the cleaned English.

There are no Go API changes since `v1.4.0-rc.1`.

## Earlier releases

For versions before 1.4.0 (v1.3.0, v1.2.0-beta, and earlier), see the
[releases page](https://github.com/byte4ever/dsco/releases) and the
[tags](https://github.com/byte4ever/dsco/tags).

[v1.4.0]: https://github.com/byte4ever/dsco/releases/tag/v1.4.0
