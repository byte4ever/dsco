# review-dsco-validation — checklist

Detailed checks for the validation lane. A failed check is an observation;
score it in Phase 3 against a demonstrated failure mode.

## Validate() presence & wiring
- [ ] A `Validate()` method exists and is called after `Fill` (and its error is
      returned/handled).
- [ ] `Fill` success is not treated as "config valid" — `Fill` only fills.

## Coverage
- [ ] Every required field (nil-able, no default) is enforced (non-nil + not the
      empty/zero sentinel where relevant).
- [ ] Numeric ranges are bounded (e.g. port 1–65535).
- [ ] Durations are bounded sanely, not merely `> 0` (guard the
      `SHUTDOWN_WAIT=15` → 15ns unit trap with a real minimum).
- [ ] Addresses / URLs are format-checked if a malformed value would fail later
      at bind/dial time.
- [ ] Cross-field invariants are checked (e.g. TLS enabled ⇒ cert+key set).

## Defaults & errors
- [ ] Defaults live in a `WithStructLayer`, not computed in caller code (origin
      shows in the location map).
- [ ] `Fill`'s returned error is checked.
- [ ] Where a strict layer or specific failure matters, typed errors are
      inspected with `errors.As` (`LayerErrors`, `FillerErrors`,
      `OverriddenKeyError`, `InvalidInputError`).

## What failure looks like
- A required nil-able field with no default and no guard → nil-deref or
  empty-credential connection failure at first use in production.
- A `Validate()` that checks non-nil but not units/ranges → a plausible operator
  typo (`=15` seconds → 15ns) passes and degrades behavior silently.
