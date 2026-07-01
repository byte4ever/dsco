# review-dsco — worked reviews (few-shot)

Three complete reviews of dsco code at known, varying quality: one poor
(REJECT), one mediocre (CONDITIONAL APPROVE), one good (APPROVE). They set the
register: enumerate before judging, score on demonstrated failure modes, run the
meta-critique, and let the worst finding govern the verdict. Enumeration lists
observations — including conforming ones — not only defects.

---

## Review A — REJECT

**Artifact (anonymised):**

```go
package config

import "github.com/byte4ever/dsco"

type Config struct {
    Host     string
    Port     int
    Password string `yaml:"password"`
    Debug    bool
}

func Load() (*Config, error) {
    config := &Config{
        Host:     "localhost",
        Port:     8080,
        Password: "s3cr3t",
        Debug:    false,
    }
    _, err := dsco.Fill(config,
        dsco.WithStructLayer(config, "defaults"),
        dsco.WithEnvLayer("APP"),
        dsco.WithCmdlineLayer(),
    )
    return config, err
}
```

### Phase 1 — Enumeration

1. `Host` is `string`, not `*string` (pointer discipline).
2. `Port` is `int`, not `*int`.
3. `Debug` is `bool`, not `*bool`.
4. `Password` is `string`, not `*string`.
5. `Host` has no `yaml` tag.
6. `Port` has no `yaml` tag.
7. `Debug` has no `yaml` tag.
8. `dsco.Fill(config, ...)` passes `*Config`, not `&config` (`**Config`).
9. The same `config` value is both the `Fill` target and the `WithStructLayer`
   input (aliasing target and source).
10. Layer order: `WithStructLayer("defaults")` is listed FIRST, so it supplies
    every field before env/cmdline are consulted (first-to-supply wins).
11. `Password: "s3cr3t"` is a secret literal in a struct committed to VCS.
12. Env prefix `APP` is generic.
13. No `Validate()` method anywhere.
14. Secret is not routed through a provider or env.
15. `Fill`'s error is returned but the `locations` map is discarded.
16. No `*Layers` function; layers are inlined at the call-site.
17. `Debug bool` with a `false` default cannot be forced to `false` explicitly
    nor distinguished from unset — moot here since it's non-pointer anyway.

### Phase 2 — Adversarial scenarios

1. **Trigger:** any deploy. **Propagation:** `Fill(config, ...)` gets `*Config`,
   not `**Config`. **Symptom:** `InvalidInputError`; the service never starts.
   **Detectability:** immediate, every run.
2. **Trigger:** operator sets `APP-PORT=9000`. **Propagation:** the defaults
   layer is first and already supplied `Port`, so env is ignored.
   **Symptom:** override silently has no effect; port stays 8080.
   **Detectability:** only by noticing the wrong port at runtime.
3. **Trigger:** repo is cloned. **Propagation:** `s3cr3t` is in source.
   **Symptom:** production credential leaked to everyone with repo access.
   **Detectability:** never, until rotated after a breach.
4. **Trigger:** a second container joins the pod. **Propagation:** it also reads
   `APP-*`. **Symptom:** cross-container config bleed. **Detectability:** only
   via a confusing production incident.
5. **Trigger:** operator forgets to set a real password in a fixed build.
   **Propagation:** `Password` defaults to `s3cr3t`; no `Validate()`.
   **Symptom:** service connects with the baked-in credential.
   **Detectability:** none pre-runtime.

### Phase 3 — Scoring

- BLOCKING: #8 (`InvalidInputError`, never boots), #11/#3 (secret in VCS),
  #10/#2 (overrides dead — defaults-first), #13 (no `Validate()` guarding the
  hardcoded secret).
- IMPORTANT: #1–#4 (non-pointer fields defeat dsco's core guarantee), #12
  (generic prefix), #9 (target/source aliasing).
- NOTED: #5–#7 (missing tags, subsumed by the pointer fix), #15 (locations
  discarded), #16 (no factoring).

### Phase 4 — Verdict

Multiple BLOCKING → REJECT. Blocking list: `Fill(&config)` fix; remove the
secret literal and route it via env/provider; reorder layers high→low so
env/cmdline win; add a `Validate()` that requires `Password`.

### Phase 5 — Meta-critique

1. *Too lenient?* I might treat the non-pointer fields as "just IMPORTANT", but
   with defaults-first ordering they also make every field unoverridable — the
   pattern is broken end to end.
2. *Category not examined?* Type conversion — moot, since manual literals are
   used, not layer strings.
3. *To a REJECT reviewer:* nothing to answer; I also REJECT.

FINAL VERDICT: REJECT

---

## Review B — CONDITIONAL APPROVE

**Artifact (anonymised):**

```go
package config

import (
    "errors"
    "time"

    "github.com/byte4ever/dsco"
)

type Config struct {
    Host    *string        `yaml:"host"`
    Port    *int           `yaml:"port"`
    Timeout *time.Duration `yaml:"timeout"`
}

func (c *Config) Validate() error {
    if c.Port != nil && (*c.Port < 1 || *c.Port > 65535) {
        return errors.New("port out of range")
    }
    return nil
}

func Load() (*Config, error) {
    var config *Config
    _, err := dsco.Fill(&config,
        dsco.WithCmdlineLayer(),
        dsco.WithEnvLayer("SVC"),
        dsco.WithStructLayer(&Config{
            Host:    dsco.R("localhost"),
            Port:    dsco.R(8080),
            Timeout: dsco.R(30 * time.Second),
        }, "defaults"),
    )
    if err != nil {
        return nil, err
    }
    return config, config.Validate()
}
```

### Phase 1 — Enumeration

1. `Host`, `Port`, `Timeout` are all `*T` — pointer discipline holds.
2. All three fields have `yaml` tags.
3. `dsco.Fill(&config, ...)` — correct `**Config` target.
4. One cmdline layer, one env layer, one struct layer — no dedup issue.
5. Layer order cmdline → env → defaults — correct high→low.
6. No strict layer — no override-error risk.
7. Env prefix `SVC` is generic.
8. Cmdline layer used with default format; no secrets on flags.
9. No secrets in this config at all.
10. `Validate()` exists and is called after `Fill`.
11. `Validate()` checks the port range but nothing else.
12. Defaults live in a `WithStructLayer` — good source attribution.
13. No manual env parsing.
14. Layers are inlined at the `Fill` site; no `Layers()` function.
15. No third-party config to embed here.
16. Not inventory-related; no version-gated API used.
17. `Fill`'s error is checked; `locations` map discarded.

### Phase 2 — Adversarial scenarios

1. **Trigger:** second container joins the pod reading `SVC-*`.
   **Propagation:** prefix collision. **Symptom:** wrong port/host.
   **Detectability:** production incident only. → prefix is generic.
2. **Trigger:** team adds an inventory binary later. **Propagation:** it must
   re-declare the layers; the two lists can diverge. **Symptom:** inventory
   reports keys the service doesn't actually read. **Detectability:** only when
   someone compares them.
3. **Trigger:** operator sets `SVC-PORT=abc`. **Propagation:** YAML conversion
   fails. **Symptom:** `Fill` returns an error; service refuses to start.
   **Detectability:** immediate — acceptable behaviour.
4. **Trigger:** operator sets `SVC-TIMEOUT=30`. **Propagation:** parsed as
   `time.Duration` (30ns). **Symptom:** absurdly short timeout; `Validate()`
   doesn't check it. **Detectability:** runtime only. → thin `Validate()`.
5. **Trigger:** `--port=0`. **Propagation:** caught by the range check.
   **Symptom:** clear validation error. **Detectability:** immediate — fine.

### Phase 3 — Scoring

- BLOCKING: none. No required nil-able field goes unguarded (all have defaults);
  the target and layer order are correct.
- IMPORTANT: #7 generic prefix `SVC` (collision risk); #14 no `Layers()`
  factoring (drift risk if inventory/tests appear); #11/scenario 4 `Validate()`
  doesn't sanity-check `Timeout`.
- NOTED: #17 `locations` discarded.

### Phase 4 — Verdict

No BLOCKING; three IMPORTANT. → CONDITIONAL APPROVE. Accepted risks, to be
documented: generic `SVC` prefix (safe only while single-container); inlined
layers (factor into `Layers()` before adding an inventory binary); `Validate()`
should bound `Timeout`.

### Phase 5 — Meta-critique

1. *Too lenient?* The generic prefix could be argued BLOCKING, but there is no
   concrete sibling container in scope — it stays IMPORTANT until one exists.
2. *Category not examined?* Alias correctness — none used, out of scope.
3. *To a REJECT reviewer:* I'd ask which finding is a *demonstrated* failure in
   the current deployment; each is a latent risk, not a guaranteed break, so
   CONDITIONAL with documented risks is the honest call.

FINAL VERDICT: CONDITIONAL APPROVE

---

## Review C — APPROVE

**Artifact (anonymised):**

```go
package config

import (
    "errors"

    "github.com/byte4ever/dsco"
    "github.com/example/pgdriver"
)

type Config struct {
    Database *pgdriver.Config `yaml:"database"`
    HTTPAddr *string          `yaml:"http_addr"`
}

func Defaults() *Config {
    return &Config{
        HTTPAddr: dsco.R(":8080"),
        Database: &pgdriver.Config{
            Host: dsco.R("localhost"),
            Port: dsco.R(5432),
            // User / Password intentionally nil: must come from a higher layer.
        },
    }
}

func Layers() []dsco.Layer {
    return []dsco.Layer{
        dsco.WithCmdlineLayer(),
        dsco.WithEnvLayer("ORDERAPI"),
        dsco.WithStructLayer(Defaults(), "defaults"),
    }
}

func (c *Config) Validate() error {
    if c.Database == nil || c.Database.User == nil || c.Database.Password == nil {
        return errors.New("database user and password are required")
    }
    return nil
}

func Load() (*Config, error) {
    var config *Config
    if _, err := dsco.Fill(&config, Layers()...); err != nil {
        return nil, err
    }
    return config, config.Validate()
}
```

### Phase 1 — Enumeration

1. `HTTPAddr` is `*string`; `Database` is `*pgdriver.Config` — pointer
   discipline holds, including the nested config.
2. Both fields have `yaml` tags; the embedded config carries its own.
3. `dsco.Fill(&config, ...)` — correct `**Config`.
4. One cmdline, one env, one struct layer — no dedup issue.
5. Order cmdline → env → defaults — correct.
6. No strict layer — no override-error risk.
7. Prefix `ORDERAPI` is role-specific.
8. No secrets on cmdline flags.
9. `User`/`Password` left nil in defaults so a higher layer MUST supply them.
10. `Validate()` requires `Database.User` and `Database.Password`.
11. Defaults live in a `WithStructLayer` via `Defaults()`.
12. No manual env parsing.
13. `Layers()` is a factored function; `Fill` and any inventory binary/test can
    share it.
14. `Defaults()` builds a fresh struct per call — no shared package var.
15. The dependency's `pgdriver.Config` is embedded, not re-declared, so
    inventory walks into it.
16. No inventory call here, but the shape is inventory-ready (`Layers()` +
    embedding).
17. `Fill`'s error is checked before `Validate()`; `locations` discarded.

### Phase 2 — Adversarial scenarios

1. **Trigger:** deploy without `ORDERAPI-DATABASE-PASSWORD`. **Propagation:**
   `Password` stays nil. **Symptom:** `Validate()` returns "user and password
   required"; the service refuses to start. **Detectability:** immediate — the
   intended behaviour.
2. **Trigger:** operator sets a bad `ORDERAPI-DATABASE-PORT=abc`.
   **Propagation:** YAML conversion fails. **Symptom:** `Fill` error at startup.
   **Detectability:** immediate.
3. **Trigger:** second container in the pod. **Propagation:** it reads its own
   prefix, not `ORDERAPI`. **Symptom:** none — role-specific prefix isolates it.
   **Detectability:** n/a.
4. **Trigger:** team adds an inventory binary. **Propagation:** it calls the same
   `Layers()`; embedding means the DB keys appear too. **Symptom:** none; the
   report is complete and can't drift. **Detectability:** n/a.
5. **Trigger:** `--http-addr=:0`. **Propagation:** accepted as a valid address;
   `Validate()` doesn't constrain it further. **Symptom:** binds an ephemeral
   port. **Detectability:** visible at startup; arguably intended for tests.
   Dismissable — not a config-safety failure.

### Phase 3 — Scoring

- BLOCKING: none.
- IMPORTANT: none. Every required field is guarded by `Validate()`; layers,
  order, target, prefix, and embedding all conform.
- NOTED: #17 `locations` discarded (lost debugging aid); scenario 5 (`http-addr`
  not further constrained) — dismissable, arguably intended.

### Phase 4 — Verdict

Phases 1–2 produced no significant observation; every adversarial scenario is
either intended behaviour or dismissable with justification. → APPROVE.

### Phase 5 — Meta-critique

1. *Too lenient?* I could manufacture a finding on the discarded `locations`
   map, but that is a debugging convenience, not a failure mode — NOTED is
   correct, not IMPORTANT.
2. *Category not examined?* Version targeting — no version-gated API is used, so
   any dsco ≥ v1.0.0-beta compiles this; nothing to gate.
3. *To a REJECT reviewer:* I'd ask them to name a concrete input that produces a
   wrong or unsafe result. Every required field fails closed via `Validate()`,
   and overrides work because defaults are last. I can't construct the break, so
   APPROVE stands.

FINAL VERDICT: APPROVE
