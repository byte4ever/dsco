# review-dsco-typing — checklist

Detailed checks for the typing lane. A failed check is an observation; score it
in Phase 3 against a demonstrated failure mode.

## Pointer discipline
- [ ] Every configurable scalar field is `*T` (`*int`, `*string`, `*bool`,
      `*float64`, `*time.Duration`, ...).
- [ ] Nested config structs are held by pointer (`*DatabaseConfig`).
- [ ] Slices and maps are NOT pointers (`[]string`, `map[string]string`).
- [ ] No `bool` field where the caller must be able to override a `true` default
      to `false` (a non-pointer bool cannot express unset vs `false`).
- [ ] `dsco.R(...)` is used to build the pointer values in the defaults layer,
      not manual `x := v; &x`.

## yaml tags
- [ ] Every field reachable from a layer has a `yaml:"..."` tag.
- [ ] Tag names match the intended env/cmdline keys; underscores in the tag are
      preserved (env `PREFIX-FIELD_NAME`, cmdline `--field-name`... verify the
      author's intent matches).
- [ ] An embedded third-party config carries its own tags; it is not re-tagged
      or shadowed locally.

## What failure looks like
- A non-pointer scalar → dsco's model scanner returns `UnsupportedTypeError`;
  `Fill` fails; the service never boots. This is the strongest typing failure.
- A missing yaml tag → the field is unreachable from cmdline/env/file layers; it
  can only ever hold its struct-layer default (silent, not an error).
