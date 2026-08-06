# CLAUDE.md — Globular (gateway + xDS)

Cluster entry point. Two binaries:
- `cmd/gateway` — HTTP/HTTPS reverse proxy, CSR signer, join flow (8 phases), join script server
- `cmd/xds` — Envoy ADS/xDS configuration server with SDS (TLS secrets)

The gateway does NOT write desired state to etcd — that is the controller's job.
PKI lives at `/var/lib/globular/pki/`. mTLS is required in production.

## Build

```bash
go build ./cmd/gateway
go build ./cmd/xds
go test ./... -race
```

## Key paths

- `cmd/gateway/` — gateway entry point
- `cmd/xds/` — xDS server entry point
- `gateway_server/` — join flow phases 1-8, CSR signing, reverse proxy handlers
- `xds/` — Envoy snapshot cache, SDS handler, snapshot builder
- `internal/server/` — etcd member management, shared server primitives
- `docs/awareness/` — awareness knowledge files (authority rules, invariants, failure modes)

---

## AI RULES — Awareness workflow

This project is registered with the awareness system. The graph lives at
`.globular/awareness/graph.json`. The knowledge files are in `docs/awareness/`.

### Required sequence for any non-trivial edit

1. **`awareness session-start`** — open a session before touching files. Records intent and establishes the edit boundary.
2. **`awareness impact <file>`** — before editing a file, check blast radius. Returns affected invariants, rules, and tests.
3. **`awareness scan-violations`** — after editing, scan for invariant violations before committing.

**`NO_MATCH` ≠ safe.** When awareness returns NO_MATCH (no graph nodes matched), it means the graph has no coverage for that file — not that the edit is safe. Always grep `docs/awareness/failure_modes.yaml`, `docs/awareness/invariants.yaml`, and `docs/awareness/forbidden_fixes.yaml` directly on NO_MATCH.

**`UNKNOWN_IMPACT`** — treat as high-risk. Do not proceed without reading the file and understanding the blast radius manually.

### High-risk files — call `awareness decision_context` before editing

- `gateway_server/` — any join phase handler; phase ordering and token validation are critical
- `gateway_server/csr.go` (or equivalent) — CSR signing; must verify token + identity + signature
- `xds/` — snapshot builder and SDS handler; stale snapshots affect the whole cluster
- `internal/server/` — etcd member management; ghost member cleanup must precede member add
- Any path that constructs PKI file paths — must use `/var/lib/globular/pki/` exclusively
- Any path that touches `GLOBULAR_XDS_INSECURE` — dev-only, never set in production unit files

### Awareness token discipline — HARD LIMIT

- **1 preflight per task** — compact (default) unless deep/forensic is justified.
- **Do NOT call `awareness agent_context` in the same turn as `awareness preflight`**.
- **Choose the smallest sufficient mode**: micro → standard → deep → forensic.
- **Never call `awareness session_resume_latest` mid-task** — only at session start if resuming.

### Key invariants enforced

- `pki.ca.canonical.path` — all PKI files under `/var/lib/globular/pki/`; never relative to domain or cwd
- `xds.mtls.required.in.production` — `GLOBULAR_XDS_INSECURE` is dev-only; never production
- `join.token.validated.before.phase` — token validated at Phase 1 and re-checked before Phase 5
- `etcd.ghost.cleared.before.member.add` — remove stale ghost member before adding new one in Phase 5
- `repair.etcd.explicit.flag.only` — WAL wipe requires `--repair-etcd`; never implicit
- `gateway.no.etcd.writes` — gateway must not write desired state or cluster config to etcd
- `csr.signature.verified.before.signing` — verify token + subject + CSR signature before CA signs
- `join.phase.order.enforced` — phases 1-8 are ordered; no skipping or reordering


## Sensei

This project uses Sensei to prevent architectural drift. Before editing files
in protected directories, consult the awareness graph.

Sensei also installs the **Sensei Architect** skill. For architecture-sensitive
planning, implementation, debugging, review, recovery, migration, security, or
state/convergence work, load `.sensei/skills/sensei-architect/SKILL.md` before
planning or editing. Claude Code can also discover the native project skill at
`.claude/skills/sensei-architect/`. Use MCP tools when configured and CLI
fallbacks otherwise. Stay proportional to the risk and preserve durable lessons
with `sensei propose`.

### Behavioral guidelines

General discipline that applies to every change (paraphrased from Andrej
Karpathy's observations on how LLMs fail at coding — the popular
[`andrej-karpathy-skills`](https://github.com/forrestchang/andrej-karpathy-skills)
CLAUDE.md). Sensei adds the *repo-specific* rules below; these are the *general*
ones:

1. **Think before coding.** State your assumptions out loud. If the request is
   ambiguous, ask. If a simpler approach exists, push back. When you are
   confused, stop and name what is unclear — do not pick one interpretation and
   run with it.
2. **Simplicity first.** Write the minimum code that solves the problem. No
   speculative abstractions, no flexibility nobody asked for. Would a senior
   engineer call this overcomplicated?
3. **Surgical changes.** Touch only what the task requires. Do not improve
   neighboring code or refactor what is not broken. Every changed line should
   trace back to the request.
4. **Goal-driven execution.** Turn a vague instruction into a verifiable target
   before writing a line. "Add validation" becomes "write tests for invalid
   inputs, then make them pass."

### Rules

1. **Call briefing before editing high-risk files.** Run `sensei briefing --file <path>`
   before modifying any file listed in `docs/awareness/high_risk_files.yaml`.
   Claude Code hooks enforce this automatically.

2. **Respect forbidden fixes.** The briefing output lists patterns that look
   correct but are known-broken. Do not use them. A PreToolUse hook
   (`edit-check-guard.sh`) also runs your proposed content through
   `sensei edit-check` and blocks an edit that introduces a forbidden-fix shape;
   if that fires, revise rather than force it through.

3. **Run required tests.** The briefing output lists tests that must pass
   when touching protected files. Run them before committing.

4. **End non-trivial code tasks with the AWG summary line:**
   ```
   AWG: briefing(<target>) | invariants: X, Y | uncertainty: Z
   ```

### Commands

```bash
sensei briefing --file <path>     # Context for a file edit
sensei briefing --task "desc"     # Context for a task
sensei preflight --file <path>    # Risk classification
sensei check                      # Validate YAML sources
sensei build                      # Recompile after editing YAML
```

### Adding Knowledge (one typed call)

After fixing a bug, record the scar with a single `sensei propose` call. It
appends the entry to the right YAML source, rebuilds the seed, reloads the
local store, and stages the change — then stops. **You review and commit.**

```bash
# A way the system broke (links the invariant it violated):
sensei propose --kind failure_mode --title "Stale seed served after reload" \
  --related-invariant <inv.id> --evidence "observed stale node" \
  --source-file golang/server/reload.go --required-test path_test.go:TestReloadFresh

# A test that proves the behavior now:
sensei propose --kind required_test --id "path_test.go:TestReloadFresh" \
  --title "Reload serves fresh triples" --related-failure <fm.id>

# A repair that must never be applied again:
sensei propose --kind forbidden_fix --title "Cache the reload path" \
  --related-invariant <inv.id> --description "caused the stale-serve failure"
```

Every entry must answer the contract-first questions: what contract was
violated/clarified, what failure we observed, what test proves it, what future
fix is forbidden, and which invariant/failure_mode it connects to. Vague notes
are rejected. If the contract is unknown, use `--kind contract_unknown` with a
`--proposed-contract` or `--revision-request`.

The `Stop` hook (`sensei feedback-check`) reminds you if a session fixed a durable
error class but added no graph feedback.
