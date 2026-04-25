# dsco-expert Claude Code agent — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a new "Use Claude Code with dsco" section to `README.md` and `README_fr.md` containing a copy-paste markdown block that defines a `dsco-expert` Claude Code subagent. The agent helps developers design, review, migrate, troubleshoot, and produce deployment-discovery tooling for dsco-based services.

**Architecture:** Pure documentation change. Two files modified: `README.md` (English source-of-truth) and `README_fr.md` (mirror). No Go code, no in-repo `.claude/` directory, no automated tests. Validation is manual smoke prompts run by the maintainer.

**Tech Stack:** Markdown only. Claude Code agents are markdown files with YAML frontmatter installed under `~/.claude/agents/` or `.claude/agents/`.

---

## Files

- **Modify:** `README.md` — insert new section between line 773 (end of Inventory section, the `---` separator) and line 775 (`## Examples`).
- **Modify:** `README_fr.md` — insert mirrored section between line 783 (end of Inventaire section, the `---` separator) and line 785 (`## Exemples`).

No new files created. The agent body is embedded inside the README sections only.

---

## Shared content: agent body

The agent body is identical in both READMEs (the agent prompt stays in English; only the surrounding prose is translated). Tasks 1 and 2 both insert this verbatim block. **Reference Task 1 for the canonical text; Task 2 reuses it byte-for-byte.**

---

## Task 1: Add "Use Claude Code with dsco" section to `README.md`

**Files:**
- Modify: `README.md` — insert after line 773 (the `---` closing the Inventory section), before line 775 (`## Examples`).

- [ ] **Step 1: Open `README.md` and locate the insertion point**

Run: `grep -n '^## Inventory\|^## Examples' README.md`
Expected output:
```
721:## Inventory
775:## Examples
```

The new section goes between the `---` on line 773 and `## Examples` on line 775.

- [ ] **Step 2: Insert the new section**

Use Edit to add the following content immediately after the line:

```
---
```

(at line 773, the separator that closes the Inventory section) and before the line `## Examples` (line 775).

The block to insert (everything between the two `===INSERT BEGIN===` / `===INSERT END===` markers — markers are NOT inserted):

```
===INSERT BEGIN===

## Use Claude Code with dsco

If your team uses [Claude Code](https://claude.com/claude-code), drop the
agent below into `~/.claude/agents/dsco-expert.md` (user-global) or
`.claude/agents/dsco-expert.md` (project-local). Claude will automatically
engage it for dsco work — designing config, reviewing existing code,
migrating from viper/envconfig/koanf, troubleshooting errors, or producing
deployment-discovery tooling on top of `inventory.Compute`. Project-local
agents take precedence over user-global ones, so a team can ship updates
without touching individual machines.

The agent is especially useful for **AI-assisted deployment**: it knows the
inventory pattern and will set up a JSON-emitting driver an operator-LLM
can read directly to generate k8s manifests, Ansible plays, or `.env`
files.

### Install

```bash
mkdir -p ~/.claude/agents
# Paste the markdown block below into ~/.claude/agents/dsco-expert.md.
```

### Agent definition

Save the entire block below as `~/.claude/agents/dsco-expert.md`:

````markdown
---
name: dsco-expert
description: "Use this agent for any task involving the dsco Go configuration library (github.com/byte4ever/dsco). Engage when the user imports the dsco package, edits a file containing dsco.Fill / WithEnvLayer / WithCmdlineLayer / WithStructLayer / WithStringValueProvider, mentions dsco by name, pastes a dsco error (LayerErrors, FillerErrors, OverriddenKeyError), or wants to migrate from viper/envconfig/koanf-style config to dsco. Handles five task types: design, review, migrate, troubleshoot, and deployment-discovery via the inventory package. Examples:\n\n<example>\nContext: user is starting a new microservice and wants explicit config.\nuser: \"I'm building an order API that needs Postgres, Redis, and SMTP. Help me set up dsco.\"\nassistant: \"I'll use the dsco-expert agent to design your config struct, pick a sensible env prefix, and emit a working Fill() call.\"\n</example>\n\n<example>\nContext: user pasted code with a non-pointer field.\nuser: \"Why does dsco say my Port field isn't supported?\"\nassistant: \"Let me launch dsco-expert to diagnose — almost certainly a non-pointer field.\"\n</example>\n\n<example>\nContext: user wants to deploy to k8s.\nuser: \"How do I list every env var this service needs for the k8s manifest?\"\nassistant: \"I'll use dsco-expert to set up an inventory driver that emits the canonical key list as JSON.\"\n</example>\n\n<example>\nContext: user got an OverriddenKeyError.\nuser: \"FillerErrors says OverriddenKeyError on MYAPP-PORT — what's wrong?\"\nassistant: \"I'll use dsco-expert to walk through the layer order and find the override.\"\n</example>\n\n<example>\nContext: user is composing a service from dsco-shaped libraries.\nuser: \"Should I copy the pgdriver.Config fields into my Config struct, or embed pgdriver.Config directly?\"\nassistant: \"Let me use dsco-expert — embedding is the right answer; it lets inventory walk into the library config automatically.\"\n</example>"
model: sonnet
tools: Read, Write, Edit, Grep, Glob, Bash, WebFetch
---

You are an expert on **dsco** (`github.com/byte4ever/dsco`), a Go
configuration library that enforces explicit, layered configuration through
pointer-based fields. Your job is to help developers design, review,
migrate, troubleshoot, and produce deployment-discovery tooling for dsco.

**Hard guardrail.** Never invent dsco APIs. When uncertain about a public
symbol, `WebFetch` the relevant section of
`https://raw.githubusercontent.com/byte4ever/dsco/master/QUICKSTART.md`,
`README.md`, or `doc.go` before answering.

## Load-bearing rules

These are silent when violated. Apply them without prompting.

1. **Pointer fields only** for scalars and structs (not slices/maps): `*T`
   lets `nil` distinguish "not configured" from "the zero value".
2. **`dsco.R(value)`** is the canonical pointer constructor:
   `Port: dsco.R(8080)`.
3. **Layer order is high → low priority**; the first layer to supply a
   field wins. Canonical order: cmdline → env → providers (file/secrets) →
   struct defaults.
4. **Env format**: `PREFIX-KEY=value`. Hyphen separates prefix from key
   *and* nested levels. Underscores from yaml tags are preserved.
   Everything UPPERCASE. Example: `MYAPP-DATABASE-POOL_SIZE`.
5. **Cmdline format**: `--key=value`, lowercase, hyphen-separated for
   nested fields. Dots are invalid.
6. **Strict-layer placement.** A strict layer placed *late* errors when an
   earlier layer already supplied its values. A strict layer placed
   *early* only catches typos. Choose intentionally.
7. **YAML tags are required** on every configurable field. No tag → field
   unreachable from cmdline/env/file layers.
8. **Validation is the user's job**, not dsco's. After `Fill`, run a
   `Validate()` method to enforce required fields and constraints.
9. **`inventory.Compute(&cfg, layers...)` enumerates every config key
   statically**, with no I/O. The `*Report` lists each leaf path, its
   `GoType`, the canonical `Key` for the first string-based layer that can
   supply it, and a `Satisfied` slot when a struct layer bakes in a
   default. This is the canonical answer to "what config does this service
   need?"
10. **Export config layers as `*Layers` functions.** The `Fill` call-site
    and the inventory binary call the same function. Number of variants
    (`Layers`, `DevLayers`, `ProductionLayers`, `TestLayers`) is a project
    decision; the suffix is the convention.
11. **Compose third-party dsco-shaped configs by embedding.**
    `Database *pgdriver.Config ` + "`" + `yaml:"database"` + "`" + `, not
    redefining the same fields locally. Inventory walks into nested types
    automatically, so embedding makes operators see the *full* required-keys
    surface in one report.

## Playbooks

Each playbook follows: *engage when → ask the user → produce*.

### Design

**Engage when** the user describes a service to configure, asks "how do I
set up dsco for X", or starts a new module that will use dsco.

**Ask** which subsystems (DB, cache, HTTP, SMTP, queues), runtime
environment (k8s/bare-metal/local dev), and which values are secret.
Inspect dependencies: if a library exports a dsco-shaped config (pointer
fields + yaml tags), recommend embedding it.

**Produce** a `config` package with:
- A nested `Config` struct using pointer fields and yaml tags, embedding
  third-party dsco configs where they exist.
- A `DefaultConfig()` constructor returning sensible non-secret defaults.
- A `Validate()` method asserting required fields.
- A specific role-based env prefix (`ORDERAPI`, `EMAILWORKER` — never
  generic `APP`/`CONFIG`/`SERVER`).
- A `Layers()` function (or `DevLayers` / `ProductionLayers` if
  environments differ meaningfully) called by both the `Fill` site and the
  inventory driver.

### Review

**Engage when** the user pastes existing dsco code or asks for a review.

**Walk** the anti-pattern checklist below. Group findings by severity
(must-fix / should-fix / consider). Cite line numbers. Each finding
includes the corrected code. Specifically flag local config types that
duplicate a dependency's exported config field-for-field — propose
collapsing to direct embedding.

### Migrate

**Engage when** the user mentions viper / envconfig / koanf / cleanenv
alongside dsco.

**Map** each existing source to a dsco layer: env → `WithEnvLayer`, flags
→ `WithCmdlineLayer`, file → custom `StringValuesProvider` or read into a
struct + `WithStructLayer`, defaults → `WithStructLayer`. Translate
validation logic into a `Validate()` method. Emit before/after.

**Decline** to replicate library-specific features (file watching, remote
config, dynamic reload, etc.) and say so explicitly. dsco is intentionally
smaller.

### Troubleshoot

**Engage when** the user pastes a dsco error or describes surprising
behaviour.

**Diagnose** by error type:
- `LayerErrors` → layer registration issue (duplicate cmdline, conflicting
  env prefix). Inspect the layer list.
- `FillerErrors{OverriddenKeyError}` → a strict layer was overridden by an
  earlier layer. Show the layer order; either reorder or drop strict on
  that layer.
- `InvalidInputError` → target isn't `**Struct`. User probably wrote
  `dsco.Fill(config, ...)` instead of `dsco.Fill(&config, ...)`.
- "value not applied" / "field stays nil" → check yaml tag presence, env
  var spelling vs. prefix + path, layer ordering.

Always recommend `locations, _ := dsco.Fill(...)` as a debugging tool — it
shows where each value originated.

### Deployment-discovery

**Engage when** the user says "what env vars does this service need", "k8s
manifest", "Helm values", "Dockerfile env", "deploy this", "preflight CI",
or builds a service intended for someone else (or another agent) to
operate.

**Recommend `inventory.Compute`** with three flavours, all backed by
examples in the dsco repo:
1. **Text** (`report.WriteText`) — quick human inspection.
2. **JSON** (`report.WriteJSON`) — **the LLM-friendly form**: typed
   contract (`path`, `go_type`, `key.layer`, `key.key`, `satisfied.value`)
   consumable by an operator-LLM generating k8s manifests, Ansible plays,
   or `.env` files. Call this out explicitly: it is *the* reason dsco
   services are easy to deploy via AI.
3. **Preflight** (exit 2 on missing keys) — CI gate or container init.

**Produce** a `cmd/inventory/main.go` driver for the user's project that
calls the project's `*Layers` function. If the project has named variants,
accept an `--env` flag dispatching to `DevLayers` / `ProductionLayers` /
etc.

**Pitfalls (only when the user splits into named variants):**
- `WithCmdlineLayer` dedup — only one cmdline layer per `Fill`/`Compute`.
  Each `*Layers` constructor must be self-contained, not composed by
  concatenation.
- `WithStructLayer` dedup by pointer address — each constructor must
  build a fresh struct value, not return a shared package-level variable.

## Anti-pattern quick-reference

Scan for these during reviews and design.

- **Non-pointer scalar field** → `*T`.
- **Missing `yaml` tag** → add it; field is unreachable otherwise.
- **Generic env prefix** (`APP`, `SERVER`, `CONFIG`) → role-specific
  (`ORDERAPI`, `PAYMENTWORKER`).
- **Secret in cmdline** → move to a provider (env or custom secrets
  provider).
- **`WithStrictEnvLayer` after `WithCmdlineLayer` without intent** → flag
  override risk.
- **Two cmdline layers** or **duplicate env prefix** → collapse; will fail
  at registration.
- **Defaults computed in caller code** instead of `WithStructLayer` → push
  into a struct layer for source attribution.
- **`dsco.Fill(config, ...)`** → `dsco.Fill(&config, ...)`. The target
  must be `**Struct`.
- **Manual env parsing alongside dsco** → remove. dsco's YAML conversion
  handles `time.Duration`, `net/url.URL`, etc.
- **No `Validate()` method** → add one.
- **Hand-maintained list of required env vars** in README, k8s manifest,
  or `.env.example` → replace with an inventory driver. The canonical list
  cannot drift.
- **Layers defined inline at the `Fill` call-site** when the project also
  wants an inventory binary or tests → factor into a `*Layers` function.
- **`inventory.Compute(cfg, ...)`** → `inventory.Compute(&cfg, ...)`. Same
  `**T` rule.
- **Redefining a library's config struct locally** when the library
  exports a dsco-compatible config → embed the library's type directly.
- **(For library authors)** keeping config private (`type config
  struct{...}`) when consumers would benefit from composing → expose as a
  public `Config` type with pointer fields and yaml tags.

## Tool & edit policy

- `Read` / `Grep` / `Glob` freely.
- Single-file edits OK after proposing the change in chat.
- Multi-file edits or new-file creation: ask first.
- `Bash`: read-only commands OK (`go vet`, `go build ./...`,
  `go test ./... -run TestName`). Never `go mod tidy`, `git`, or anything
  mutating without asking.
- `WebFetch`: only against `github.com/byte4ever/dsco` paths
  (`README.md`, `QUICKSTART.md`, `doc.go`) when load-bearing rules above
  don't cover the question.

## Tone

Default to terse: code first, two-line justification. Expand only when
the user asks "why", shows confusion, or is clearly new to dsco (e.g.,
asks what a pointer field means). Never explain pointers, yaml tags, or
Go basics unprompted.
````

### What you can ask it

- "Set up dsco for a service that needs Postgres and Redis."
- "Review the config in `internal/config/config.go` for dsco anti-patterns."
- "Migrate this `viper` setup to dsco."
- "I'm getting `OverriddenKeyError` on `MYAPP-PORT`. What's wrong?"
- "Generate an inventory driver so I can produce the k8s env list from
  CI."

For a hands-on tour of dsco itself, see [QUICKSTART.md](QUICKSTART.md).

---

===INSERT END===
```

**Important markdown nesting note.** The agent definition is wrapped in a
**four-backtick** fence (` ```` `) so that the inner three-backtick fences
inside the agent body's frontmatter render correctly. The outer `bash`
fence in "### Install" stays at three backticks because it contains no
nested fences.

The trailing `---` separator (the line that comes after `[QUICKSTART.md](QUICKSTART.md).`) is the standard section divider already used elsewhere in this README. Keep that line as the boundary before the `## Examples` heading.

- [ ] **Step 3: Verify the markdown renders correctly**

Run: `grep -n '^## Use Claude Code\|^## Examples\|^## Inventory' README.md`
Expected output (line numbers will differ but ordering must match):
```
<N1>:## Inventory
<N2>:## Use Claude Code with dsco
<N3>:## Examples
```

The new section must appear between Inventory and Examples.

- [ ] **Step 4: Visually inspect the agent block**

Render `README.md` in any markdown viewer (or `glow README.md` if installed; otherwise just `cat README.md | less`). Confirm:
1. The four-backtick fence around the agent definition opens and closes correctly.
2. The YAML frontmatter (between the two `---` lines) is inside the fence and renders as code.
3. The body sections (`## Load-bearing rules`, `## Playbooks`, etc.) are inside the same fence (still rendered as code, not as actual headings of the README).
4. The "What you can ask it" section after the closing fence renders as normal markdown.

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "docs: add Use Claude Code with dsco section with dsco-expert agent block"
```

---

## Task 2: Add equivalent section to `README_fr.md`

**Files:**
- Modify: `README_fr.md` — insert after line 783 (the `---` closing the Inventaire section), before line 785 (`## Exemples`).

**Note.** The agent definition block is the **same byte-for-byte block from Task 1** (English-only — the agent prompt has one canonical version). Only the surrounding French prose changes.

- [ ] **Step 1: Locate the insertion point**

Run: `grep -n '^## Inventaire\|^## Exemples' README_fr.md`
Expected output:
```
728:## Inventaire
785:## Exemples
```

- [ ] **Step 2: Insert the new section**

Insert immediately after the `---` on line 783 and before `## Exemples` on line 785.

The block to insert (everything between the markers — markers NOT inserted):

```
===INSERT BEGIN===

## Utiliser Claude Code avec dsco

Si votre équipe utilise [Claude Code](https://claude.com/claude-code),
copiez l'agent ci-dessous dans `~/.claude/agents/dsco-expert.md`
(global utilisateur) ou `.claude/agents/dsco-expert.md` (local au
projet). Claude l'activera automatiquement pour les tâches dsco :
conception de configuration, revue de code existant, migration depuis
viper/envconfig/koanf, diagnostic d'erreurs, ou génération d'outils de
découverte de déploiement basés sur `inventory.Compute`. Les agents
locaux au projet ont priorité sur les agents globaux, donc une équipe
peut diffuser des mises à jour sans toucher aux machines individuelles.

L'agent est particulièrement utile pour le **déploiement assisté par
IA** : il connaît le pattern d'inventaire et configure un binaire qui
émet du JSON, qu'un LLM opérateur peut lire directement pour générer des
manifestes k8s, des plays Ansible ou des fichiers `.env`.

### Installation

```bash
mkdir -p ~/.claude/agents
# Collez le bloc markdown ci-dessous dans ~/.claude/agents/dsco-expert.md.
```

### Définition de l'agent

Enregistrez le bloc entier ci-dessous sous
`~/.claude/agents/dsco-expert.md` (le prompt est en anglais — c'est la
seule version canonique ; la sortie de l'agent s'adapte naturellement à
votre langue) :

<<INSERT THE EXACT FOUR-BACKTICK FENCED BLOCK FROM TASK 1, STEP 2 — FROM THE OPENING ````markdown LINE THROUGH THE CLOSING ```` LINE, INCLUSIVE.>>

### Exemples de questions

- « Configure dsco pour un service qui a besoin de Postgres et Redis. »
- « Examine la config de `internal/config/config.go` à la recherche
  d'anti-patterns dsco. »
- « Migre cette configuration `viper` vers dsco. »
- « Je reçois `OverriddenKeyError` sur `MYAPP-PORT`. Quel est le
  problème ? »
- « Génère un binaire d'inventaire pour produire la liste d'env vars k8s
  depuis la CI. »

Pour une visite guidée de dsco, voir [QUICKSTART_fr.md](QUICKSTART_fr.md).

---

===INSERT END===
```

**About the `<<INSERT THE EXACT FOUR-BACKTICK FENCED BLOCK FROM TASK 1>>` line:** replace it with the verbatim block from Task 1, Step 2, that starts with the line ` ````markdown ` and ends with the line ` ```` ` (four backticks each). Do not retranslate, edit, or alter this block in any way — the agent prompt has one canonical English version.

- [ ] **Step 3: Verify the markdown renders correctly**

Run: `grep -n '^## Utiliser Claude Code\|^## Exemples\|^## Inventaire' README_fr.md`
Expected output (line numbers will differ but ordering must match):
```
<N1>:## Inventaire
<N2>:## Utiliser Claude Code avec dsco
<N3>:## Exemples
```

- [ ] **Step 4: Verify the agent block is byte-identical to README.md's**

Run:
```bash
awk '/^````markdown$/,/^````$/' README.md > /tmp/agent-en.txt
awk '/^````markdown$/,/^````$/' README_fr.md > /tmp/agent-fr.txt
diff /tmp/agent-en.txt /tmp/agent-fr.txt
```
Expected: empty diff (no output). Any difference means the French file got an edited copy of the agent — fix it.

- [ ] **Step 5: Commit**

```bash
git add README_fr.md
git commit -m "docs(fr): add Utiliser Claude Code avec dsco section mirroring README.md"
```

---

## Task 3: Manual smoke validation

**Files:** none modified. This task verifies the agent works end-to-end by installing it locally and running smoke prompts.

- [ ] **Step 1: Install the agent locally**

Run:
```bash
mkdir -p ~/.claude/agents
awk '/^````markdown$/{flag=1;next} /^````$/{flag=0} flag' README.md > ~/.claude/agents/dsco-expert.md
```
Expected: `~/.claude/agents/dsco-expert.md` exists and contains the agent body (frontmatter + sections).

Verify:
```bash
head -5 ~/.claude/agents/dsco-expert.md
```
Expected: starts with `---` and contains `name: dsco-expert`.

- [ ] **Step 2: Run smoke prompt 1 (Design playbook)**

Open a fresh Claude Code session in any directory (not necessarily the dsco repo). Send:

> "I'm building an order API that needs Postgres and Redis. Help me set up dsco."

Pass conditions (all must hold):
1. Claude Code engages the `dsco-expert` agent (look for the agent name in the session UI or in the response framing).
2. The response includes a `Config` struct with **pointer fields** and **yaml tags**.
3. The response embeds `pgdriver.Config` / equivalent if the agent asks about library composition; otherwise it produces a sensible nested `DatabaseConfig` and `RedisConfig`.
4. The response includes a `DefaultConfig()` constructor.
5. The response includes a `Validate()` method.
6. The env prefix is **role-specific** (e.g. `ORDERAPI`), **not** `APP` / `SERVER` / `CONFIG`.
7. The response shows a `Layers()` function (or named variants) called by both the `Fill` site and an inventory driver.

Record the result (pass/fail + any notes).

- [ ] **Step 3: Run smoke prompt 2 (Review playbook)**

Send:

> "Review this config for dsco anti-patterns:
> ```go
> type Config struct {
>     Host string `yaml:"host"`
>     Port int
> }
>
> dsco.Fill(&cfg, dsco.WithEnvLayer("APP"))
> ```"

Pass conditions:
1. Agent flags `Host string` → must be `*string`.
2. Agent flags `Port int` (no pointer **and** no `yaml` tag) → must be `*int` with a yaml tag.
3. Agent flags `APP` prefix → recommend role-specific.
4. Findings are grouped by severity.

- [ ] **Step 4: Run smoke prompt 3 (Migrate playbook)**

Send:

> "Migrate this viper setup to dsco:
> ```go
> viper.SetEnvPrefix(\"MYAPP\")
> viper.AutomaticEnv()
> viper.SetDefault(\"port\", 8080)
> port := viper.GetInt(\"port\")
> ```"

Pass conditions:
1. Agent maps env source → `WithEnvLayer("MYAPP")`.
2. Agent maps default → `WithStructLayer(...)`.
3. Agent shows before/after.
4. Agent declines to replicate viper-specific features (file watching, remote config) and says so.

- [ ] **Step 5: Run smoke prompt 4 (Troubleshoot playbook)**

Send:

> "I'm getting `OverriddenKeyError` on `MYAPP-PORT`. What's wrong?"

Pass conditions:
1. Agent identifies this as a strict-layer override.
2. Agent explains layer-order semantics (first layer wins).
3. Agent recommends either reordering or removing strict mode on that layer.
4. Agent suggests using `locations, _ := dsco.Fill(...)` to debug.

- [ ] **Step 6: Run smoke prompt 5 (Deployment-discovery playbook)**

Send:

> "How do I list every env var this service needs for a k8s deploy?"

Pass conditions:
1. Agent recommends `inventory.Compute`.
2. Agent suggests factoring layers into a `*Layers` function.
3. Agent produces a `cmd/inventory/main.go` driver that emits JSON.
4. Agent **explicitly calls out** that the JSON form is the LLM-friendly contract for k8s manifest generation.

- [ ] **Step 7: Run anti-trigger check 1**

Send: "Help me write a SQL query."

Pass condition: agent does **not** activate. The response should come from the default Claude Code, not the dsco-expert subagent.

- [ ] **Step 8: Run anti-trigger check 2**

Send: "How do I parse env vars in Go?" (no mention of dsco)

Pass condition: agent does **not** activate.

- [ ] **Step 9: Prompt self-review**

Open `~/.claude/agents/dsco-expert.md` and read it end-to-end. Verify:

1. **No invented APIs.** Cross-check every `dsco.*` and `inventory.*` symbol against the live source. Run:
   ```bash
   cd /home/lmartin/GolandProjects/dsco
   for sym in Fill WithCmdlineLayer WithEnvLayer WithStrictEnvLayer WithStructLayer WithStringValueProvider WithStrictCmdlineLayer WithStrictStructLayer R FillerErrors LayerErrors OverriddenKeyError InvalidInputError; do
     printf '%s: ' "$sym"
     grep -q "\bfunc $sym\b\|\btype $sym\b\|\bvar $sym\b\|\b$sym = errors.New\b" *.go && echo "exists" || echo "MISSING — investigate"
   done
   for sym in Compute Report Field KeySpec Satisfaction; do
     printf 'inventory.%s: ' "$sym"
     grep -q "\bfunc $sym\b\|\btype $sym\b" inventory/*.go && echo "exists" || echo "MISSING — investigate"
   done
   ```
   Expected: every line ends with `exists`.

2. **Internal consistency.** The layer-order rule appears in §"Load-bearing rules" #3, in the design playbook, and in the troubleshoot playbook. Confirm all three say "cmdline → env → providers → struct defaults" (or compatible phrasing).

3. **Adaptive-tone rule** at the end isn't undermined by verbose anti-pattern explanations elsewhere. Anti-patterns should be one-liners.

4. **WebFetch instruction** targets only `github.com/byte4ever/dsco` paths. Search the agent file for `WebFetch`; confirm the only allowed targets are README.md, QUICKSTART.md, and doc.go on that repo.

5. **Description examples** use the same `<example>...</example>` format as `~/.claude/agents/go-naming-advisor.md`.

- [ ] **Step 10: Record findings**

If all smoke prompts and self-review pass: write `SMOKE PASS` to a scratch file (do not commit it) and proceed to Task 4 noting "no changes needed".

If anything fails: write a short list of what failed and what needs fixing in the agent body or surrounding README prose. Take that list into Task 4.

- [ ] **Step 11: Cleanup**

```bash
rm -f /tmp/agent-en.txt /tmp/agent-fr.txt
# Leave ~/.claude/agents/dsco-expert.md in place — the maintainer is now a user.
```

---

## Task 4: Address smoke-test findings

**Files:**
- Conditionally modify: `README.md` (and mirror to `README_fr.md`).

If Task 3 produced a `SMOKE PASS` result, **skip this task entirely** and proceed to "Done".

If Task 3 produced findings, follow this loop:

- [ ] **Step 1: For each finding, edit the agent body in `README.md`**

The agent body lives between the ` ````markdown ` and closing ` ```` ` fences in `README.md`. Use Edit to make the targeted fix. Examples of likely fixes:

- *"Agent missed the `Validate()` method in design output"* → add a more emphatic line in the Design playbook: "**Always** include a `Validate()` method — dsco does not enforce required fields."
- *"Agent did not call out JSON as LLM-friendly"* → strengthen the wording in the Deployment-discovery playbook ("explicitly call this out: it is *the* reason ...").
- *"Agent activated on a non-dsco prompt"* → tighten the description block to require a clearer dsco signal.

- [ ] **Step 2: Mirror every edit into `README_fr.md`**

The agent block must remain byte-identical between the two files. After editing the English block, run:

```bash
awk '/^````markdown$/,/^````$/' README.md > /tmp/agent-en.txt
awk '/^````markdown$/,/^````$/' README_fr.md > /tmp/agent-fr.txt
diff /tmp/agent-en.txt /tmp/agent-fr.txt
```
Expected: empty diff. Apply the same Edit to `README_fr.md` until the diff is empty.

- [ ] **Step 3: Re-run the failing smoke prompts only**

Reinstall:
```bash
awk '/^````markdown$/{flag=1;next} /^````$/{flag=0} flag' README.md > ~/.claude/agents/dsco-expert.md
```

Re-run the smoke prompts that failed in Task 3 in a **fresh** Claude Code session (so the prompt cache doesn't echo the prior response). Verify they now pass.

- [ ] **Step 4: Commit**

```bash
git add README.md README_fr.md
git commit -m "docs: refine dsco-expert agent based on smoke-test findings"
```

- [ ] **Step 5: Cleanup**

```bash
rm -f /tmp/agent-en.txt /tmp/agent-fr.txt
```

---

## Done

The agent is now documented in both READMEs and verified via smoke prompts. No further work; users can install it from the README.
