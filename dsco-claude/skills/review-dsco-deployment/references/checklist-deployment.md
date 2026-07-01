# review-dsco-deployment — checklist

Detailed checks for the deployment lane. A failed check is an observation;
score it in Phase 3 against a demonstrated failure mode.

## Version targeting
- [ ] The dsco version in `go.mod` provides every API used here.
- [ ] Any version-gated API (esp. `inventory.Compute`, `*Report`, the `Write*`
      methods, min v1.4.0-rc.1) is available at the pinned version.
- [ ] If it is not, that is a finding: the code won't compile / advice is ahead
      of the pin — offer `go get github.com/byte4ever/dsco@v1.4.0-rc.1`.

## Inventory
- [ ] `inventory.Compute(&cfg, ...)` (the `**T` rule).
- [ ] Output flavour matches the use: text (human), JSON (tooling / operator-LLM
      contract: `path`, `go_type`, `key.layer`, `key.key`, `satisfied.value`),
      preflight (exit 2 on missing keys, CI/init gate).
- [ ] No hand-maintained required-env list (README, k8s manifest, `.env.example`)
      that inventory should generate and keep from drifting.

## Reuse & composition
- [ ] The `*Layers` function is EXPORTED so an inventory binary / another
      package's test can call it.
- [ ] The `Fill` site and the inventory driver call the SAME `*Layers` function
      (no re-declaration that can drift).
- [ ] Each named-variant constructor is self-contained (one cmdline layer each;
      fresh struct per constructor).
- [ ] A dependency's exported dsco-shaped `Config` is embedded (`*pgdriver.Config`)
      so inventory walks into it, not re-declared field-for-field locally.

## What failure looks like
- `import ".../inventory"` against dsco < v1.4.0-rc.1 → build break for everyone
  on the pinned version.
- Unexported `layers()` → the inventory binary re-declares the layer list; the
  generated key list drifts from what the service actually reads.
- A dependency config re-declared locally → inventory misses the dependency's
  keys; operators under-provision.
