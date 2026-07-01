---
name: dsco
description: >-
  Authoritative guide for WRITING idiomatic configuration with the dsco Go
  library (github.com/byte4ever/dsco): pointer-based config structs, layered
  Fill(), env / cmdline / struct / provider layers, strict mode, Validate(),
  and inventory-based deployment discovery. Use whenever dsco code is produced
  or designed, when migrating viper / envconfig / koanf to dsco, or when the
  user mentions dsco.Fill / WithEnvLayer / WithCmdlineLayer / WithStructLayer /
  WithStringValueProvider / dsco.R / inventory.Compute. After producing code it
  self-reviews through review-dsco with a bounded correction loop (up to 3
  cycles, then escalates). To REVIEW existing dsco code directly, use the
  review-dsco skill.
x-dsco-target: v1.4.0-rc.1
x-bundle-version: 1.4.0-rc.1
---

# dsco

Guide for writing configuration in the spirit of dsco: explicit, layered,
pointer-based. This skill covers WRITING and DESIGNING dsco code, and it does
not present that code until it has passed its own review: every artifact it
produces goes through the **review-dsco** orchestrator in a bounded
correction → re-validation loop (see *Self-review loop*).

**Hard guardrail.** Never invent dsco APIs. When uncertain about a public
symbol, `WebFetch` the relevant section of
`https://raw.githubusercontent.com/byte4ever/dsco/master/QUICKSTART.md`,
`README.md`, or `doc.go` before writing code that uses it.

## Self-review loop (mandatory)

You do not hand the user unreviewed dsco code. Every time you PRODUCE dsco code
— a config struct, a `Fill` call-site, a `*Layers` function, or an inventory
driver — run it through the **review-dsco** orchestrator and converge before
presenting it.

The loop, per authored artifact:

1. **Review.** Submit the CURRENT code to `review-dsco` (it fans out the aspect
   reviewers as isolated sub-agents). Each review is fresh and anonymous: never
   tell the reviewers the code is yours, and never pass the history of prior
   rounds. Read the GLOBAL verdict.
2. **Branch on the verdict:**
   - **REJECT** (≥1 BLOCKING): fix every BLOCKING finding exactly as its
     remediation states (clear IMPORTANT too when cheap), then go back to step 1.
     This correction → re-validation cycle repeats up to **3 times**.
   - **CONDITIONAL APPROVE** (only IMPORTANT/NOTED): converged. Present the code
     WITH the accepted-risks list so the user can decide; optionally clear cheap
     IMPORTANT first.
   - **APPROVE**: converged. Present the code.
3. **Escalate after 3 cycles.** If the verdict is STILL REJECT after the 3rd
   correction → re-validation cycle, STOP looping. Tell the user, in plain terms:
   - the code as it currently stands;
   - the BLOCKING finding(s) that won't clear;
   - what you changed in each of the 3 iterations and why it did not converge
     (e.g. fixing lane A re-broke lane B, or the requirement is self-contradictory);
   - and ask them to arbitrate: accept a trade-off, relax a requirement, or take
     over.

Convergence discipline:
- Each iteration must actually resolve the previous round's BLOCKING items. A fix
  that introduces a NEW blocking item is non-convergence — say so.
- Never loosen the artifact just to make the reviewer pass (dropping a required
  field, widening a type, deleting a `Validate` check). Fix the real defect; a
  green verdict bought by weakening the config is a failure, not a pass.
- The cap is a hard 3 cycles. "At least 3" means do not give up before trying
  three times; "no more than 3" means escalate rather than grind past it.

## Version targeting

This skill ships in the `dsco-claude` bundle, versioned in lockstep with the
library: bundle **vX.Y.Z** targets dsco **vX.Y.Z**. **This copy targets dsco
`v1.4.0-rc.1`.** The rules below assume that version's public API.

Before writing version-gated code, and always before reaching for `inventory`:

1. **Determine the version the user depends on** — read the
   `require github.com/byte4ever/dsco vX.Y.Z` line in their `go.mod`. No module
   yet → assume a fresh install (latest) and say so.
2. **If their pinned version is older than a feature's minimum**, do NOT emit
   code that will not compile against it. Say the feature needs a newer dsco and
   offer the upgrade:
   ```bash
   go get github.com/byte4ever/dsco@v1.4.0-rc.1
   go mod tidy
   ```
   Then wait for acceptance, or fall back to advice valid for their version.
3. **Unsure which version introduced a symbol?** `WebFetch` the repo at the tag
   rather than guessing.

| Feature / API | Minimum dsco |
|---|---|
| Core: `Fill`, `WithEnvLayer`, `WithCmdlineLayer`, `WithStructLayer`, `WithStrictEnvLayer`, `WithStringValueProvider`, `dsco.R` | v1.0.0-beta |
| `inventory.Compute`, `*Report`, `Report.WriteText` / `WriteJSON` / `WriteYAML` | **v1.4.0-rc.1** |

## Load-bearing rules

Apply without prompting. Silent when violated.

1. **Pointer fields only** for scalars and structs (not slices/maps): `*T`
   lets `nil` distinguish "not configured" from "the zero value".
2. **`dsco.R(value)`** is the canonical pointer constructor: `Port: dsco.R(8080)`.
3. **Layer order is high → low priority**; the first layer to supply a field
   wins. Canonical order: cmdline → env → providers (file/secrets) → struct
   defaults.
4. **Env format**: `PREFIX-KEY=value`. Hyphen separates prefix from key *and*
   nested levels. Underscores from yaml tags are preserved. UPPERCASE. Example:
   `MYAPP-DATABASE-POOL_SIZE`.
5. **Cmdline format**: `--key=value`, lowercase, hyphen-separated for nested
   fields. Dots are invalid.
6. **Strict-layer placement.** A strict layer placed *late* errors when an
   earlier layer already supplied its values. A strict layer placed *early* only
   catches typos. Choose intentionally.
7. **YAML tags are required** on every configurable field. No tag → field
   unreachable from cmdline/env/file layers.
8. **Validation is the caller's job**, not dsco's. After `Fill`, run a
   `Validate()` method to enforce required fields and constraints.
9. **`inventory.Compute(&cfg, layers...)` enumerates every config key
   statically** (no I/O). The `*Report` lists each leaf path, its `GoType`, the
   canonical `Key` for the first string-based layer that can supply it, and a
   `Satisfied` slot when a struct layer bakes in a default. Requires dsco
   ≥ v1.4.0-rc.1 — check `go.mod` and offer the upgrade first.
10. **Export config layers as `*Layers` functions.** The `Fill` call-site and
    the inventory binary call the same function. Variants (`Layers`,
    `DevLayers`, `ProductionLayers`, `TestLayers`) are a project decision; the
    suffix is the convention.
11. **Compose third-party dsco-shaped configs by embedding** (`Database
    *pgdriver.Config `yaml:"database"``), not by redefining fields locally.
    Inventory walks into nested types, so embedding shows operators the *full*
    required-keys surface in one report.

## Playbooks

Each playbook: *engage when → ask → produce*.

### Design

**Engage when** the user describes a service to configure, asks "how do I set
up dsco for X", or starts a module that will use dsco.

**Ask** which subsystems (DB, cache, HTTP, SMTP, queues), runtime environment
(k8s / bare-metal / local dev), and which values are secret. If a dependency
exports a dsco-shaped config (pointer fields + yaml tags), embed it.

**Produce** a `config` package with: a nested `Config` struct (pointer fields +
yaml tags, embedding third-party dsco configs where they exist); a
`DefaultConfig()` constructor of non-secret defaults; a `Validate()` method; a
role-based env prefix (`ORDERAPI`, `EMAILWORKER` — never `APP`/`CONFIG`/
`SERVER`); and a `Layers()` function called by both the `Fill` site and the
inventory driver.

### Migrate

**Engage when** the user mentions viper / envconfig / koanf / cleanenv.

**Map** each source to a layer: env → `WithEnvLayer`, flags → `WithCmdlineLayer`,
file → custom `StringValuesProvider` or read-into-struct + `WithStructLayer`,
defaults → `WithStructLayer`. Move validation into a `Validate()` method. Emit
before/after.

**Decline** to replicate library-specific features (file watching, remote
config, dynamic reload) and say so. dsco is intentionally smaller.

### Troubleshoot

**Engage when** the user pastes a dsco error or surprising behaviour.

Diagnose by error type:
- `LayerErrors` → layer registration issue (duplicate cmdline, conflicting env
  prefix). Inspect the layer list.
- `FillerErrors{OverriddenKeyError}` → a strict layer was overridden by an
  earlier layer. Show the order; reorder or drop strict on that layer.
- `InvalidInputError` → target isn't `**Struct`; the user wrote
  `dsco.Fill(config, ...)` instead of `dsco.Fill(&config, ...)`.
- "value not applied" / "field stays nil" → check yaml tag presence, env var
  spelling vs prefix + path, layer ordering.

Recommend `locations, _ := dsco.Fill(...)` as a debugging tool — it shows where
each value originated.

### Deployment-discovery

**Engage when** the user says "what env vars does this service need", "k8s
manifest", "Helm values", "Dockerfile env", "deploy this", "preflight CI", or
builds a service for someone else (human or agent) to operate.

**Version gate.** Built on `inventory`, which needs **dsco ≥ v1.4.0-rc.1**.
Check `go.mod` first; offer the upgrade before producing inventory code.

Recommend `inventory.Compute` in three flavours: text (`WriteText`, human
inspection), JSON (`WriteJSON`, the LLM-friendly typed contract an operator-LLM
reads to generate k8s manifests / Ansible / `.env`), preflight (exit 2 on
missing keys, a CI/init gate). Produce a `cmd/inventory/main.go` calling the
project's `*Layers` function; add an `--env` flag if there are named variants.

Pitfall (named variants only): `WithCmdlineLayer` dedup — one cmdline layer per
`Fill`/`Compute`; each `*Layers` constructor is self-contained. `WithStructLayer`
dedups by pointer address — each constructor builds a fresh struct, not a shared
package-level variable.

## Pitfalls

The anti-pattern catalog (what to scan for while writing) lives in
[references/pitfalls.md](references/pitfalls.md). Read it before finalizing any
config struct or `Fill` call-site.

## Tone

Default to terse: code first, two-line justification. Expand only when the user
asks "why", is confused, or is new to dsco. Never explain pointers, yaml tags,
or Go basics unprompted.
