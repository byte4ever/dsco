# review-dsco — per-category checklist

Detailed checks behind each Phase-1 category. A check that fails is an
*observation*, not yet a severity — score it in Phase 3 against a demonstrated
failure mode.

## 1. Pointer discipline
- [ ] Every configurable scalar field is `*T` (`*int`, `*string`, `*bool`,
      `*time.Duration`, ...).
- [ ] Nested config structs are held by pointer (`*DatabaseConfig`).
- [ ] Slices and maps are NOT pointers (they are already nilable).
- [ ] No `bool` field where the caller must be able to override a `true` default
      to `false` (a non-pointer bool cannot express that).

## 2. yaml tags
- [ ] Every field reachable from a layer has a `yaml:"..."` tag.
- [ ] Tag names match the intended env/cmdline keys (underscores preserved).
- [ ] Embedded third-party configs carry their own tags (not re-tagged locally).

## 3. Fill target
- [ ] `dsco.Fill(&config, ...)` — target is `**Struct`.
- [ ] `config` is a `*Config` (or `var config *Config`), so `&config` is `**Config`.

## 4. Layer set
- [ ] At most one cmdline layer per `Fill`/`Compute`.
- [ ] No two env layers share a prefix.
- [ ] Each `*Layers` constructor is self-contained (not concatenated from
      another, which would duplicate the cmdline layer).

## 5. Layer order
- [ ] Layers listed high → low priority (cmdline → env → providers → defaults).
- [ ] No struct-default layer placed above the layers meant to override it.
- [ ] The first layer able to supply a field is the intended winner.

## 6. Strict-layer placement
- [ ] For each `WithStrict*Layer`, list every other layer that can supply the
      same field.
- [ ] A strict layer placed AFTER such a layer will `OverriddenKeyError` when
      both supply the field — confirm that is intended (forbidding overrides),
      not accidental.
- [ ] A strict layer placed EARLY only catches typos/unused keys — confirm the
      author isn't expecting it to forbid overrides.

## 7. Env-prefix quality
- [ ] Prefix is role-specific (`ORDERAPI`, `PAYMENTWORKER`), not `APP`/`SERVER`/
      `CONFIG`/`SVC`.
- [ ] In a multi-container pod, the prefix cannot collide with a sibling
      container's vars.
- [ ] The prefix cannot collide with system vars (`PATH`, `HOME`, `HTTP_PROXY`).

## 8. Cmdline
- [ ] Flags are `--key=value`, lowercase, hyphen-separated for nested fields,
      no dots.
- [ ] No secret is passed as a flag (visible in `ps`, shell history).
- [ ] Aliases (`WithAliases`) map to real field paths.

## 9. Secret routing
- [ ] Secrets come from a provider (`WithStringValueProvider`) or env, never a
      cmdline flag.
- [ ] No secret literal is baked into a `WithStructLayer` and committed to VCS.
- [ ] Secret fields are left nil in the defaults layer so a higher layer MUST
      supply them.

## 10. Validate()
- [ ] A `Validate()` method exists and is called after `Fill`.
- [ ] It enforces required fields (the nil-able secrets/hosts), value ranges
      (port 1–65535), and cross-field invariants.
- [ ] `Fill` success is not treated as "config is valid" — `Fill` only fills.

## 11. Defaults source
- [ ] Defaults live in a `WithStructLayer`, not computed in caller code, so
      their origin shows up in `Fill`'s location map.

## 12. Manual env parsing
- [ ] No `os.Getenv` + `strconv`/manual parsing sits alongside dsco for config
      that dsco could carry (dsco's YAML conversion handles durations, URLs,
      ints, bools).

## 13. *Layers factoring
- [ ] Layer list is a `Layers()` (or `DevLayers`/`ProductionLayers`) function
      when the project also has an inventory binary or tests, not duplicated at
      the `Fill` site.

## 14. Struct-layer freshness
- [ ] Each named-variant constructor builds a fresh struct value; none returns a
      shared package-level variable (pointer-address dedup would drop it).

## 15. Third-party config composition
- [ ] A dependency's exported dsco-shaped `Config` is embedded
      (`*pgdriver.Config`), not re-declared field-for-field locally.

## 16. Inventory usage (dsco ≥ v1.4.0-rc.1)
- [ ] `inventory.Compute(&cfg, ...)` (the `**T` rule).
- [ ] Output flavour matches the use (text human / JSON tooling / preflight CI).
- [ ] No hand-maintained required-env list that inventory should generate.

## 17. Version targeting
- [ ] The dsco version in `go.mod` actually provides every API used here
      (esp. `inventory`, min v1.4.0-rc.1).
- [ ] If it does not, that is a finding: the code won't compile / the advice is
      ahead of the pinned version.

## 18. Error handling
- [ ] `Fill`'s returned error is checked.
- [ ] Where relevant, typed errors are inspected (`errors.As` on `LayerErrors`,
      `FillerErrors`, `OverriddenKeyError`, `InvalidInputError`).
