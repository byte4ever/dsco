# review-dsco — severity rubric

Score on a **demonstrated failure mode**, never on aesthetic preference or a
"best practice" invoked without a concrete case. Use these calibrated examples
to place a finding. When in doubt between two levels, write the failing scenario
first, then let the impact decide.

## BLOCKING — demonstrated failure, fix before approval

- **Required secret no layer supplies, and `Validate()` doesn't catch it.**
  `Password *string` stays nil after `Fill`; `*config.Password` panics at the
  first DB connect in production. Nothing fails at startup.
- **Strict layer guaranteed to error at startup.** `WithStrictEnvLayer("SVC")`
  placed after `WithCmdlineLayer()` where both routinely supply `--port` /
  `SVC-PORT`; every deploy that sets both dies with `OverriddenKeyError`.
- **Secret on a cmdline flag.** `--db-password=...` is visible in `ps`, the
  shell history, and process listings; a real credential leak.
- **`dsco.Fill(config, ...)`** (single pointer) → `InvalidInputError` at
  startup; the service never boots.
- **Inventory driver against dsco < v1.4.0.** `import ".../inventory"` does
  not compile; the build is broken for anyone on the pinned version.
- **Shared package-level struct across named variants.** `DevLayers` and
  `ProductionLayers` both return the same `&defaults`; pointer-dedup drops one,
  so one environment silently loses its defaults.

## IMPORTANT — justified concern, accept only with the risk documented

- **Generic env prefix in a currently single-container deployment.** `APP-PORT`
  works today but collides the day a second container joins the pod; no failure
  yet, real latent risk.
- **Missing `*Layers` factoring when an inventory binary is planned.** Layers
  are inlined at the `Fill` site; the inventory driver will re-declare them and
  can drift. No failure until the two lists disagree.
- **Non-pointer bool with a `true` default that ops may need to disable.**
  `Verbose bool` can't be set back to `false` via a layer; a real limitation,
  but only bites if that override is ever needed.
- **`Validate()` present but incomplete.** It checks the port range but not the
  required DB host; narrows, doesn't close, the missing-config gap.
- **Local re-declaration of a dependency's exported dsco `Config`.** Works, but
  inventory won't show the dependency's keys and the two structs drift.

## NOTED — awareness only, no action required

- **`locations` map not captured for debugging.** `_, err := dsco.Fill(...)`
  discards the source map; a lost debugging aid, not a defect.
- **Category out of scope for this artifact.** No cmdline layer in a
  server-only config, so cmdline-format checks don't apply — record the skip.
- **Stylistic prefix choice within the role-specific convention.** `ORDERAPI`
  vs `ORDER_API`; consistent and unambiguous, no impact.
- **Defaults constructor named `DefaultConfig()` vs `Defaults()`.** Naming
  preference, no failure mode.

## Boundary calls

- A generic prefix is **IMPORTANT** in a single-container service but
  **BLOCKING** once a concrete sibling container shares the pod and reads the
  same var — the demonstrated collision moves it up.
- A missing `Validate()` is **IMPORTANT** when every required field is supplied
  by a strict layer that already errors on absence, but **BLOCKING** when a
  required nil-able field has no other guard.
