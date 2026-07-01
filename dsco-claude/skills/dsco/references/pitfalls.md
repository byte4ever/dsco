# dsco pitfalls (anti-pattern catalog)

Scan for these while writing or designing dsco config. Each entry is
*symptom → fix*. The review-dsco skill scores the same list against a
demonstrated failure mode.

## Struct shape

- **Non-pointer scalar field** (`Port int`) → `*T`. A zero value can't be told
  apart from "unset", so a missing value is silently accepted.
- **Missing `yaml` tag** → add it. Without a tag the field is unreachable from
  cmdline / env / file layers; it can only ever hold its struct-layer default.
- **Pointer on a slice/map** → don't. Slices and maps are nilable already; keep
  them non-pointer (`Tables []string `yaml:"tables"``).
- **Redefining a library's config struct locally** when the library exports a
  dsco-compatible `Config` → embed the library's type directly, so inventory
  walks it and operators see the full key surface.
- **(Library authors) private config** (`type config struct{...}`) when
  consumers would compose it → expose a public `Config` with pointer fields and
  yaml tags.

## Layers & precedence

- **`dsco.Fill(config, ...)`** → `dsco.Fill(&config, ...)`. The target must be
  `**Struct`; the single-pointer form yields `InvalidInputError`.
- **Two cmdline layers** or a **duplicate env prefix** → collapse; both fail at
  registration (`LayerErrors` / `CmdlineAlreadyUsedError`).
- **Layer order wrong** — remember first-to-supply wins, high → low. Putting
  struct defaults *before* env/cmdline makes the overrides dead.
- **Layers defined inline at the `Fill` call-site** when the project also wants
  an inventory binary or tests → factor into a `*Layers` function so both
  call-sites share one definition.
- **`WithStructLayer` returning a shared package-level struct** across named
  variants → build a fresh value per constructor (dedup is by pointer address).

## Strict mode

- **`WithStrictEnvLayer` / `WithStrictStructLayer` placed late without intent**
  → flag override risk: a strict layer after another layer that supplies the
  same field errors with `OverriddenKeyError`. Place strict early to catch only
  typos, or late deliberately to forbid overrides.

## Env & secrets

- **Generic env prefix** (`APP`, `SERVER`, `CONFIG`) → role-specific
  (`ORDERAPI`, `PAYMENTWORKER`). Generic prefixes collide across containers in a
  shared pod and read ambiguously in manifests.
- **Secret in a cmdline flag** → move it to a provider (env or a custom secrets
  provider). Flags land in shell history and `ps` output.
- **Manual env parsing alongside dsco** (`os.Getenv` + `strconv`) → remove.
  dsco's YAML conversion already handles `time.Duration`, `net/url.URL`, ints,
  bools, etc.

## Validation

- **No `Validate()` method** → add one. dsco fills; it does not enforce ranges,
  required combinations, or cross-field invariants.
- **Defaults computed in caller code** instead of a `WithStructLayer` → push
  them into a struct layer so their source is attributable via `Fill`'s
  location map.

## Inventory (dsco ≥ v1.4.0)

- **`inventory.Compute(cfg, ...)`** → `inventory.Compute(&cfg, ...)`. Same
  `**T` rule as `Fill`.
- **Hand-maintained list of required env vars** in a README, k8s manifest, or
  `.env.example` → replace with an inventory driver. The generated list can't
  drift from the code.
- **Using inventory against dsco < v1.4.0** → the package does not exist
  there. Check `go.mod` and offer the upgrade before writing the driver.
