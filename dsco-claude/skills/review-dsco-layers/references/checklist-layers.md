# review-dsco-layers — checklist

Detailed checks for the layers lane. A failed check is an observation; score it
in Phase 3 against a demonstrated failure mode.

## Fill target
- [ ] `dsco.Fill(&cfg, ...)` — target is `**Struct`.
- [ ] `cfg` is a `*Config` (or `var cfg *Config`), so `&cfg` is `**Config`.
- [ ] `inventory.Compute(&cfg, ...)` follows the same `**T` rule.

## Layer set
- [ ] At most one cmdline layer per `Fill`/`Compute`.
- [ ] No two env layers share a prefix.
- [ ] Each `*Layers` constructor is self-contained, not concatenated from
      another (which would duplicate the cmdline layer).

## Layer order & precedence
- [ ] Layers listed high → low priority (cmdline → env → providers → defaults).
- [ ] No struct-default layer placed above the layers meant to override it.
- [ ] For each field, the first layer able to supply it is the intended winner.

## Strict-layer placement
- [ ] For each `WithStrict*Layer`, list every other layer that can supply the
      same field.
- [ ] A strict layer AFTER such a layer will `OverriddenKeyError` when both
      supply the field — confirm that is intended (forbidding overrides).
- [ ] A strict layer placed EARLY only catches typos/unused keys — confirm the
      author isn't expecting it to forbid overrides.

## Factoring & freshness
- [ ] Layer list is a `Layers()` function (not inlined) when the project also
      has an inventory binary or tests.
- [ ] Each named-variant constructor builds a fresh struct value; none returns a
      shared package-level variable (pointer-address dedup would drop it).

## Bypass
- [ ] No `os.Getenv`/`strconv` manual parsing sits alongside dsco for config
      that a layer could carry.

## What failure looks like
- `Fill(cfg)` single pointer → `InvalidInputError`; never boots.
- Struct defaults listed first → every env/cmdline override is silently dead.
- Strict layer after a supplier of the same field → `OverriddenKeyError` at
  startup in the environment where both are set.
