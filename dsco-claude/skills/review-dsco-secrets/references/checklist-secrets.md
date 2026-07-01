# review-dsco-secrets — checklist

Detailed checks for the secrets/env lane. A failed check is an observation;
score it in Phase 3 against a demonstrated failure mode.

## Env-prefix quality
- [ ] Prefix is role-specific (`ORDERAPI`, `PAYMENTWORKER`), not `APP`/`SERVER`/
      `CONFIG`/`SVC`/`SERVICE`.
- [ ] In a multi-container pod, the prefix cannot collide with a sibling
      container's vars.
- [ ] The prefix cannot collide with system vars (`PATH`, `HOME`, `HTTP_PROXY`)
      or a platform-injected `PREFIX_*` namespace.
- [ ] The prefix isn't confusable with a well-known tool's env namespace (note
      dsco uses hyphen separators, `PREFIX-KEY`).

## Secret routing
- [ ] Every secret-looking field (`password`, `token`, `api_key`, `dsn`,
      `database_url`, ...) is identified and its intended channel stated.
- [ ] Secrets arrive via env or a `WithStringValueProvider`, never a cmdline
      flag (flags land in `ps`, `/proc/<pid>/cmdline`, shell history).
- [ ] No secret literal is baked into a `WithStructLayer` and committed to VCS.
- [ ] Secret fields are left nil in the defaults layer so a higher layer MUST
      supply them.
- [ ] For high-value secrets, a provider (vault/secrets manager) is used rather
      than plain env, where warranted.

## Cmdline format
- [ ] Flags are `--key=value`, lowercase, hyphen-separated for nested fields,
      no dots.
- [ ] Aliases (`WithAliases`) map to real field paths.

## What failure looks like
- A secret literal in a struct layer → credential in VCS history for everyone
  with repo access.
- A generic prefix in a shared pod → a sibling container reads the same
  `PREFIX-*` var; cross-container config bleed.
- A cmdline layer above a secret field → the secret can be passed as a flag and
  leak in process listings.
