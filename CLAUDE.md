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
