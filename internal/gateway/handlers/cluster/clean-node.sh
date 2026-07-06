#!/usr/bin/env bash
set -euo pipefail

# ── Globular Node Cleanup ─────────────────────────────────────────────────────
#
# Removes this node cleanly from the cluster, then wipes all local state so the
# node is ready for a fresh Day-1 join.
#
# Phase 0 (NEW): Cluster-level detachment — runs while services are still UP:
#   a. Removes the node from the cluster controller via gateway HTTP API
#      (cascades to envoy/xDS endpoint removal and MinIO pool eviction).
#   b. Decommissions the ScyllaDB node (streams data to peers before shutdown).
#   c. Removes the etcd member (prevents quorum breakage on remaining peers).
#
# Detachment is BEST-EFFORT: any failure warns loudly (with the manual
# remediation command) and continues — a node that is already isolated from the
# cluster is exactly when a force-clean is needed, so detachment must never
# block the local wipe. Use --local-only to skip cluster detachment entirely.
#
# Phases 1–5: Local cleanup — stops services and wipes state.
#
# Usage:
#   sudo bash clean-node.sh              # interactive (asks before wiping)
#   sudo bash clean-node.sh --force      # non-interactive (no prompts)
#   sudo bash clean-node.sh --local-only # skip cluster detachment, wipe locally
#
# Can be run remotely:
#   ssh user@node "sudo bash -s" < clean-node.sh

FORCE=0
LOCAL_ONLY=0
TEST_ALLOW_NON_ROOT="${GLOBULAR_CLEAN_NODE_TEST_ALLOW_NON_ROOT:-0}"
TEST_STOP_AFTER="${GLOBULAR_CLEAN_NODE_TEST_STOP_AFTER:-}"
STATE_DIR_ROOT="${GLOBULAR_STATE_DIR_OVERRIDE:-/var/lib/globular}"

die() { echo "  ✗ ERROR: $*" >&2; exit 1; }
log_info() { echo "  → $*"; }
log_success() { echo "  ✓ $*"; }
log_warn() { echo "  ⚠ $*"; }
log_step() { echo ""; echo "━━━ $* ━━━"; echo ""; }

# Parse args after the log/die helpers are defined — the unknown-argument case
# calls die(), so the helpers must exist before this loop runs.
while [[ $# -gt 0 ]]; do
    case "${1:-}" in
        --force)
            FORCE=1
            ;;
        --local-only)
            LOCAL_ONLY=1
            ;;
        *)
            die "unknown argument: ${1}"
            ;;
    esac
    shift
done
phase_marker() {
    local marker_id="${1:-}"
    local marker_text="${2:-}"
    [[ -n "${marker_text}" ]] && log_info "${marker_text}"
    if [[ -n "${TEST_STOP_AFTER}" && "${TEST_STOP_AFTER}" == "${marker_id}" ]]; then
        log_success "Test stop after phase marker: ${marker_id}"
        exit 0
    fi
}

count_scylla_up_nodes() {
    local status_output=""
    status_output="$(nodetool status 2>/dev/null || true)"
    if [[ -z "${status_output}" ]]; then
        log_warn "nodetool status unavailable — treating ScyllaDB peer count as 0" >&2
        printf '0\n'
        return 0
    fi
    printf '%s\n' "${status_output}" | awk '/^U[NL] / {n++} END {print n+0}'
}

resolve_node_agent_state_path() {
    local canonical="${STATE_DIR_ROOT}/node-agent/state.json"
    local legacy="${STATE_DIR_ROOT}/nodeagent/state.json"
    if [[ -f "${canonical}" ]]; then
        printf '%s\n' "${canonical}"
        return 0
    fi
    if [[ -f "${legacy}" ]]; then
        printf '%s\n' "${legacy}"
        return 0
    fi
    printf '%s\n' "${canonical}"
}

# hard_stop_scylla — kills ScyllaDB completely before any wipe.
# Fails closed (exits non-zero) if Scylla cannot be killed within 10s, because
# a live Scylla process can recreate /var/lib/scylla state during the wipe.
hard_stop_scylla() {
    log_info "Hard-stopping ScyllaDB before wipe..."

    # Stop and disable all Scylla systemd units.
    for unit in scylla-server.service scylla-node-exporter.service scylla-tune-sched.service \
                scylla-manager.service scylla-manager-agent.service; do
        systemctl stop "${unit}" 2>/dev/null || true
        systemctl disable "${unit}" 2>/dev/null || true
        systemctl kill -s SIGKILL "${unit}" 2>/dev/null || true
    done

    # Stop any Scylla timers.
    for timer in $(systemctl list-timers 'scylla-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $NF}'); do
        systemctl stop "${timer}" 2>/dev/null || true
    done

    # Kill by exact process name (comm field).
    pkill -9 -x scylla 2>/dev/null || true
    pkill -9 -x scylla-manager 2>/dev/null || true
    pkill -9 -x scylla-manager-agent 2>/dev/null || true

    # Wait up to 10 s for all Scylla processes to exit.
    # Match ONLY real ScyllaDB binaries by exact process name (-x), never the
    # full command line (-f): a caller whose argv merely *contains* "scylla" —
    # this cleanup script itself, or `systemctl reset-failed 'scylla-*'`, or a
    # log tail — must not be mistaken for a live database and block the wipe.
    local scylla_name_re='scylla|scylla-manager|scylla-manager-agent'
    for i in $(seq 1 10); do
        if ! pgrep -x "${scylla_name_re}" >/dev/null 2>&1; then
            log_success "No ScyllaDB process remains"
            return 0
        fi
        sleep 1
    done

    log_warn "ScyllaDB processes still alive after hard stop:"
    pgrep -ax "${scylla_name_re}" || true
    die "Refusing to wipe /var/lib/scylla while ScyllaDB may still be running. Kill the process manually and rerun."
}

# assert_scylla_wiped — verifies all Scylla on-disk state was removed.
# Fails closed if any path still exists, preventing a false "ready for Day-1 join" message.
assert_scylla_wiped() {
    local failed=0
    for path in /var/lib/scylla /etc/scylla /etc/scylla.d; do
        if [[ -e "${path}" ]]; then
            log_warn "Scylla path still exists after wipe: ${path}"
            failed=1
        fi
    done

    if [[ "${failed}" -eq 1 ]]; then
        die "ScyllaDB wipe incomplete; refusing to mark node ready for Day-1 join"
    fi

    log_success "ScyllaDB local state fully removed"
}

# Must be root
[[ $EUID -eq 0 || "${TEST_ALLOW_NON_ROOT}" == "1" ]] || die "This script must be run as root (use sudo)"

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          GLOBULAR NODE CLEANUP                                 ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "  Host: $(hostname)"
echo "  Date: $(date)"
echo ""

if [[ $FORCE -eq 0 ]] && [[ -t 0 ]]; then
  echo "  This will remove this node from the cluster and wipe all local data."
  echo "  Press Enter to continue, or Ctrl+C to abort..."
  read -r
fi

# ── Phase 0a: Preserve AI memory before wipe ─────────────────────────────────
#
# ai-memory and behavioral-memory live in ScyllaDB keyspaces (ai_memory,
# behavioral_memory) that this script is about to destroy. Export them as a
# logical CQL dump to /var/backups/globular/ai-memory-snapshot — a path the
# wipe below does NOT touch — so install-day0.sh can restore them after a fresh
# bootstrap (clean → day-0 round-trip).
#
# Self-contained on purpose: uses cqlsh directly, independent of the
# backup-manager service (which may be old/unhealthy during a teardown) and of
# scylla-manager. Must run BEFORE Phase 0 below, because ScyllaDB decommission
# streams this node's data away. Best-effort: NEVER blocks the teardown, but
# warns loudly so the operator knows whether ai-memory was preserved.
# Opt out with: GLOBULAR_SKIP_AI_BACKUP=1

AI_BACKUP_DIR="/var/backups/globular/ai-memory-snapshot"

resolve_scylla_cql_host() {
  local h=""
  if [[ -f /etc/scylla/scylla.yaml ]]; then
    h=$(grep -E '^rpc_address:' /etc/scylla/scylla.yaml 2>/dev/null | awk '{print $2}' | tr -d "'\"" || true)
    [[ -z "$h" ]] && h=$(grep -E '^listen_address:' /etc/scylla/scylla.yaml 2>/dev/null | awk '{print $2}' | tr -d "'\"" || true)
  fi
  case "$h" in ""|"localhost"|"0.0.0.0"|127.*|"::1") h="" ;; esac
  [[ -z "$h" ]] && h=$(hostname -I 2>/dev/null | tr ' ' '\n' | grep -vE '^$|^127\.|^::' | head -1 || true)
  echo "${h:-127.0.0.1}"
}

backup_ai_memory() {
  if [[ "${GLOBULAR_SKIP_AI_BACKUP:-0}" == "1" ]]; then
    log_info "GLOBULAR_SKIP_AI_BACKUP=1 — skipping ai-memory preservation"
    return 0
  fi
  if ! command -v cqlsh >/dev/null 2>&1; then
    log_warn "cqlsh not found — CANNOT preserve ai-memory before wipe"
    return 1
  fi
  if ! systemctl is-active --quiet scylla-server.service 2>/dev/null; then
    log_warn "ScyllaDB not running — CANNOT preserve ai-memory before wipe"
    return 1
  fi
  local host; host="$(resolve_scylla_cql_host)"
  if ! cqlsh "$host" 9042 -e "SELECT now() FROM system.local" >/dev/null 2>&1; then
    log_warn "ScyllaDB CQL not reachable at ${host}:9042 — CANNOT preserve ai-memory"
    return 1
  fi

  rm -rf "$AI_BACKUP_DIR"
  mkdir -p "$AI_BACKUP_DIR"
  local ks tables tbl exported=0
  for ks in ai_memory behavioral_memory; do
    if ! cqlsh "$host" 9042 -e "DESCRIBE KEYSPACE ${ks}" >/dev/null 2>&1; then
      log_info "keyspace ${ks} not present — nothing to preserve"
      continue
    fi
    mkdir -p "${AI_BACKUP_DIR}/${ks}"
    # Capture schema as a restore safety net (service normally recreates it).
    cqlsh "$host" 9042 -e "DESCRIBE KEYSPACE ${ks}" > "${AI_BACKUP_DIR}/${ks}/schema.cql" 2>/dev/null || true
    # Derive base-table names from the captured schema (robust vs cqlsh SELECT
    # output decoration). Indexes / materialized views are rebuilt on restore.
    tables=$(grep -oiE "CREATE TABLE ${ks}\.[A-Za-z0-9_]+" "${AI_BACKUP_DIR}/${ks}/schema.cql" 2>/dev/null \
             | awk '{print $3}' | sed "s/^${ks}\.//" | sort -u || true)
    for tbl in $tables; do
      [[ -z "$tbl" ]] && continue
      if cqlsh "$host" 9042 -e "COPY ${ks}.${tbl} TO '${AI_BACKUP_DIR}/${ks}/${tbl}.csv' WITH HEADER=true;" >/dev/null 2>&1; then
        log_success "Exported ${ks}.${tbl}"
        exported=$((exported + 1))
      else
        log_warn "Could not export ${ks}.${tbl} (continuing)"
      fi
    done
  done

  if [[ "$exported" -gt 0 ]]; then
    date -u +%Y-%m-%dT%H:%M:%SZ > "${AI_BACKUP_DIR}/.saved_at"
    hostname > "${AI_BACKUP_DIR}/.source_host" 2>/dev/null || true
    log_success "AI memory preserved (${exported} table(s)) at ${AI_BACKUP_DIR} — survives the wipe"
    return 0
  fi
  log_warn "No ai-memory tables were exported"
  return 1
}

log_step "Preserving AI Memory (pre-wipe)"
if backup_ai_memory; then
  :
else
  log_warn "ai-memory was NOT preserved — a fresh day-0 will start with seeded knowledge only."
  if [[ $FORCE -eq 0 ]] && [[ -t 0 ]]; then
    echo "  Continue with the wipe anyway? Press Enter to proceed, or Ctrl+C to abort..."
    read -r
  fi
fi

# ── Phase 0: Cluster-level detachment ────────────────────────────────────────
#
# Must run BEFORE stopping services: ScyllaDB decommission and etcd member
# remove both require the respective service to be running. Controller removal
# triggers xDS endpoint pruning and MinIO pool eviction automatically.

log_step "Detaching from Cluster (before local wipe)"

if [[ "$LOCAL_ONLY" == "1" ]]; then
  log_warn "--local-only: skipping cluster detachment (controller / ScyllaDB ring / etcd member / MinIO pool)."
  log_warn "  This node's peer-membership in each of those subsystems is NOT cleaned by a local wipe."
  log_warn "  On a multi-node cluster, remove it from a controller/admin host after the wipe —"
  log_warn "  this single command cascades to xDS endpoint pruning AND MinIO pool eviction:"
  log_warn "    globular cluster nodes remove <node-id> --force --drain=false"
fi

_STATE_DIR="${STATE_DIR_ROOT}"
_PKI_DIR="${_STATE_DIR}/pki"
_STATE_FILE="$(resolve_node_agent_state_path)"
_ETCD_CACERT="${_PKI_DIR}/ca.crt"
_ETCD_CERT="${_PKI_DIR}/issued/etcd/client.crt"
_ETCD_KEY="${_PKI_DIR}/issued/etcd/client.key"
_NODE_IP=$(hostname -I | awk '{print $1}')
_ETCD_ENDPOINT="https://${_NODE_IP}:2379"

# Locate globular CLI binary
_GLOBULAR_BIN=$(command -v globular 2>/dev/null || true)
[[ -z "$_GLOBULAR_BIN" ]] && [[ -x "${_STATE_DIR}/bin/globularcli" ]] && _GLOBULAR_BIN="${_STATE_DIR}/bin/globularcli"
# Locate etcdctl
_ETCDCTL_BIN=$(command -v etcdctl 2>/dev/null || true)
[[ -z "$_ETCDCTL_BIN" ]] && [[ -x "${_STATE_DIR}/bin/etcdctl" ]] && _ETCDCTL_BIN="${_STATE_DIR}/bin/etcdctl"

# Read node ID from node-agent state file
_NODE_ID=""
if [[ -f "$_STATE_FILE" ]] && command -v python3 >/dev/null 2>&1; then
  _NODE_ID=$(python3 -c "
import json
try:
    d = json.load(open('$_STATE_FILE'))
    print(d.get('node_id', '').strip())
except Exception:
    pass
" 2>/dev/null || true)
fi

# ── 0.1 Remove from cluster controller ───────────────────────────────────────
# Primary: gateway HTTP API (DELETE /api/cluster/nodes/<id>). The gateway uses
# its own controller auth — no user token is required on the cleaning node.
# Fallback: globular CLI (needs a cached token at ~/.config/globular/token).

if [[ "$LOCAL_ONLY" != "1" && -n "$_NODE_ID" ]]; then
  phase_marker "controller_removal_start" "controller removal start"
  # Derive gateway host from controller_endpoint in state.json (strip scheme/port).
  _GATEWAY_HOST="globular.internal"
  if [[ -f "$_STATE_FILE" ]] && command -v python3 >/dev/null 2>&1; then
    _GH=$(python3 -c "
import json, re
try:
    d = json.load(open('$_STATE_FILE'))
    ep = d.get('controller_endpoint', '').strip()
    ep = re.sub(r'^https?://', '', ep)
    ep = re.sub(r':\d+$', '', ep)
    if ep: print(ep)
except Exception: pass
" 2>/dev/null || true)
    [[ -n "$_GH" ]] && _GATEWAY_HOST="$_GH"
  fi

  log_info "Removing node ${_NODE_ID} from cluster via gateway API (${_GATEWAY_HOST}:8443)..."
  # NOTE: no -f. Without -f, curl exits 0 on an HTTP error and -w reports the
  # real status (404/500/…); on a connection failure curl reports 000 via -w and
  # exits non-zero. Do NOT append `|| echo 000` — curl already prints 000 on
  # failure, so the extra echo concatenated into a bogus "000000". The `|| true`
  # only guards `set -e`; the `[[ -z ]]` guard covers the no-output edge case.
  _HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -X DELETE "https://${_GATEWAY_HOST}:8443/api/cluster/nodes/${_NODE_ID}" \
    -k -H "Content-Type: application/json" \
    -d '{"force":true,"drain":false}' 2>/dev/null || true)
  [[ -z "$_HTTP_STATUS" ]] && _HTTP_STATUS="000"

  if [[ "$_HTTP_STATUS" == "200" ]]; then
    log_success "Node removed from cluster controller"
    phase_marker "controller_removal_success" "controller removal success"
  else
    phase_marker "controller_removal_failure" "controller removal failure"
    log_warn "Gateway API returned HTTP ${_HTTP_STATUS} — falling back to globular CLI..."
    if [[ -n "$_GLOBULAR_BIN" ]]; then
      _REMOVE_ERR=$("$_GLOBULAR_BIN" cluster nodes remove "$_NODE_ID" --force --drain=false 2>&1 || true)
      if echo "$_REMOVE_ERR" | grep -q "^message:"; then
        log_success "Node removed from cluster controller (via CLI)"
        phase_marker "controller_removal_success" "controller removal success"
      else
        # Best-effort: detachment failure must NOT block the local wipe.
        log_warn "CLI removal also failed: ${_REMOVE_ERR}"
        log_warn "  Controller removal is what cascades to xDS endpoint pruning AND"
        log_warn "  MinIO pool eviction — neither happened. Run this from a controller/admin"
        log_warn "  host after the wipe to complete both:"
        log_warn "    globular cluster nodes remove ${_NODE_ID} --force --drain=false"
      fi
    else
      # No CLI on this node and the gateway API was unreachable. Warn with the
      # manual remediation and continue — the local wipe still proceeds.
      log_warn "globular CLI not found — controller removal could not run on this node."
      log_warn "  xDS endpoint pruning AND MinIO pool eviction cascade from controller removal;"
      log_warn "  neither happened. Complete both from a controller/admin host after the wipe:"
      log_warn "    globular cluster nodes remove ${_NODE_ID} --force --drain=false"
      log_warn "  or: curl -X DELETE https://${_GATEWAY_HOST}:8443/api/cluster/nodes/${_NODE_ID} -k -d '{\"force\":true,\"drain\":false}'"
    fi
  fi
elif [[ "$LOCAL_ONLY" != "1" && -z "$_NODE_ID" ]]; then
  log_warn "No node ID in ${_STATE_FILE} — skipping controller removal (node may not be registered)"
fi

# ── 0.2 ScyllaDB: decommission before shutdown ───────────────────────────────
# Streams data to remaining peers; must run while scylla-server is still active.
# Skip when this is the only ScyllaDB node (nothing to stream to).
if [[ "$LOCAL_ONLY" != "1" ]] && systemctl is-active --quiet scylla-server.service 2>/dev/null; then
  if command -v nodetool >/dev/null 2>&1; then
    _SCYLLA_UP="$(count_scylla_up_nodes)"
    if [[ "$_SCYLLA_UP" -gt 1 ]]; then
      log_info "Decommissioning ScyllaDB node (streaming data to peers — this may take a few minutes)..."
      if nodetool decommission 2>/dev/null; then
        log_success "ScyllaDB node decommissioned cleanly"
      else
        log_warn "ScyllaDB decommission failed — data may be under-replicated."
        log_warn "  From another node: nodetool removenode <host-id>"
      fi
    else
      log_info "Single-node ScyllaDB — skipping decommission"
    fi
  else
    log_warn "nodetool not found — skipping ScyllaDB decommission"
    log_warn "  From another node after this wipe: nodetool removenode <host-id>"
  fi
fi

# ── 0.3 etcd: remove member before data wipe ─────────────────────────────────
# Without this the remaining peers still count this node toward quorum and will
# stall on the next leader election or write if it stays missing.
if [[ "$LOCAL_ONLY" != "1" ]] \
    && systemctl is-active --quiet globular-etcd.service 2>/dev/null \
    && [[ -n "$_ETCDCTL_BIN" ]] \
    && [[ -f "$_ETCD_CACERT" ]] && [[ -f "$_ETCD_CERT" ]] && [[ -f "$_ETCD_KEY" ]]; then

  _MEMBER_ID=$(ETCDCTL_API=3 "$_ETCDCTL_BIN" \
    --endpoints="$_ETCD_ENDPOINT" \
    --cacert="$_ETCD_CACERT" --cert="$_ETCD_CERT" --key="$_ETCD_KEY" \
    member list --write-out=json 2>/dev/null | \
    python3 -c "
import json, sys
try:
    d = json.load(sys.stdin)
    node_ip = '${_NODE_IP}'
    for m in d.get('members', []):
        urls = m.get('peerURLs', []) + m.get('clientURLs', [])
        if any(node_ip in u for u in urls):
            print(m['ID'])
            break
except Exception:
    pass
" 2>/dev/null || true)

  if [[ -n "$_MEMBER_ID" ]]; then
    log_info "Removing etcd member ${_MEMBER_ID} (${_NODE_IP}) from cluster..."
    if ETCDCTL_API=3 "$_ETCDCTL_BIN" \
        --endpoints="$_ETCD_ENDPOINT" \
        --cacert="$_ETCD_CACERT" --cert="$_ETCD_CERT" --key="$_ETCD_KEY" \
        member remove "$_MEMBER_ID" 2>/dev/null; then
      log_success "etcd member removed — remaining peers updated"
    else
      log_warn "etcd member remove failed — remaining peers may have a ghost member."
      log_warn "  From another etcd member: etcdctl member remove ${_MEMBER_ID}"
    fi
  else
    log_warn "This node not found in etcd member list — may already be removed"
  fi
elif [[ "$LOCAL_ONLY" != "1" ]] && systemctl is-active --quiet globular-etcd.service 2>/dev/null; then
  log_warn "etcdctl or TLS certs missing — skipping etcd member removal"
  log_warn "  Manual fix: etcdctl member list → etcdctl member remove <id>"
fi

# ── Phase 1: Stop services ────────────────────────────────────────────────────

phase_marker "service_stop_start" "service stop start"
log_step "Stopping Services"

# Stop all globular services
while IFS= read -r unit; do
  [[ -n "$unit" ]] || continue
  log_info "Stopping $unit"
  systemctl stop "$unit" 2>/dev/null || true
  systemctl disable "$unit" 2>/dev/null || true
done < <(systemctl list-units 'globular-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $1}' || true)

# Stop ScyllaDB (best-effort via systemctl; hard_stop_scylla below does the
# definitive kill and verifies no process remains before we wipe anything).
for unit in scylla-server.service scylla-node-exporter.service scylla-tune-sched.service \
            scylla-manager.service scylla-manager-agent.service; do
  if systemctl is-active --quiet "$unit" 2>/dev/null || systemctl is-enabled --quiet "$unit" 2>/dev/null; then
    log_info "Stopping $unit"
    systemctl stop "$unit" 2>/dev/null || true
    systemctl disable "$unit" 2>/dev/null || true
  fi
done

# Stop ScyllaDB timers
while IFS= read -r timer; do
  [[ -n "$timer" ]] || continue
  log_info "Stopping timer $timer"
  systemctl stop "$timer" 2>/dev/null || true
  systemctl disable "$timer" 2>/dev/null || true
done < <(systemctl list-timers 'scylla-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $NF}' || true)

# Hard-kill ScyllaDB — must succeed before any wipe begins.
# This is a non-negotiable gate: a live Scylla process can recreate system
# table state in /var/lib/scylla even while the directory is being wiped.
hard_stop_scylla

# ── Phase 2: Force-kill survivors ─────────────────────────────────────────────

log_step "Force-Killing Surviving Processes"

# Kill all globular server processes
while IFS= read -r proc; do
  [[ -n "$proc" ]] || continue
  cmd=$(ps -p "$proc" -o comm= 2>/dev/null || true)
  log_warn "Killing PID $proc ($cmd)"
  kill -9 "$proc" 2>/dev/null || true
done < <(ps aux 2>/dev/null | awk '/_server|globularcli|mcp|gateway|xds_server|envoy/ && $0 !~ /awk/ {print $2}' || true)

# Kill etcd if running
pkill -9 -x etcd 2>/dev/null && log_warn "Killed etcd" || true

sleep 1

# ── Phase 3: Remove unit files ───────────────────────────────────────────────

log_step "Removing Unit Files"

REMOVED=0
for unit_file in /etc/systemd/system/globular-*.service; do
  [[ -f "$unit_file" ]] || continue
  rm -f "$unit_file"
  rm -f "${unit_file}.sha256"
  log_success "Removed $(basename "$unit_file")"
  REMOVED=$((REMOVED + 1))
done

# Remove any orphaned sha256 sidecars whose unit file was already gone.
for sha_file in /etc/systemd/system/globular-*.service.sha256; do
  [[ -f "$sha_file" ]] || continue
  rm -f "$sha_file"
  log_success "Removed orphaned $(basename "$sha_file")"
done

# Remove drop-in dirs
for dropin in /etc/systemd/system/globular-*.service.d; do
  [[ -d "$dropin" ]] || continue
  rm -rf "$dropin"
  log_success "Removed $(basename "$dropin")"
done

systemctl daemon-reload 2>/dev/null || true

# ── Phase 4: Wipe state ─────────────────────────────────────────────────────

phase_marker "data_wipe_start" "data wipe start"
log_step "Wiping State"
phase_marker "package_state_wipe_start" "package/state wipe start"

# Globular state — unconditional rm -rf (safe on missing dirs, avoids
# permission-race with the globular user that was just removed)
# Remove stale Globular wrapper scripts from /usr/local/bin that point
# into /usr/lib/globular/bin (which gets removed below). Without this
# they break system commands like sha256sum after the wipe.
for wrapper in /usr/local/bin/claude /usr/local/bin/ffmpeg /usr/local/bin/sctool \
               /usr/local/bin/mc /usr/local/bin/etcdctl /usr/local/bin/globular \
               /usr/local/bin/globularcli /usr/local/bin/restic /usr/local/bin/rclone \
               /usr/local/bin/yt-dlp /usr/local/bin/sha256sum; do
  if [[ -f "$wrapper" ]] && grep -q "usr/lib/globular" "$wrapper" 2>/dev/null; then
    rm -f "$wrapper"
    log_success "Removed stale wrapper $(basename "$wrapper")"
  fi
done

for dir in /var/lib/globular /etc/globular /usr/lib/globular; do
  rm -rf "$dir" && log_success "Removed $dir" || log_warn "Could not fully remove $dir (retrying with -f)"
  rm -rf "$dir" 2>/dev/null || true
done

# MinIO object data (mounted volume — not under /var/lib/globular)
for dir in /mnt/data/minio /var/lib/minio; do
  if [[ -d "$dir" ]]; then
    rm -rf "$dir"
    log_success "Removed $dir"
  fi
done

# Remove ScyllaDB package entirely so the node-agent owns the install from
# scratch on rejoin. Keeping the binary causes a race: systemd auto-starts
# scylla-server before the node-agent can take control, and Scylla hangs on
# SIGTERM while loading system tables requiring a manual SIGKILL.
# Package purge runs BEFORE directory wipe so the package manager's postrm
# scripts cannot restart or recreate Scylla state during the rm -rf below.
if dpkg -l 'scylla*' 2>/dev/null | grep -q '^ii'; then
  log_info "Removing ScyllaDB packages (node-agent will reinstall on rejoin)"
  DEBIAN_FRONTEND=noninteractive apt-get remove -y --purge 'scylla*' 2>/dev/null || \
    log_warn "apt remove scylla failed — continuing"
  log_success "ScyllaDB packages removed"
fi

# Wipe all ScyllaDB state and data.
# hard_stop_scylla already confirmed no Scylla process is alive, so this wipe
# is race-free.
for dir in /var/lib/scylla /etc/scylla /etc/scylla.d; do
  if [[ -d "$dir" ]]; then
    rm -rf "$dir"
    log_success "Removed $dir"
  fi
done

# Assert the wipe is complete before we declare the node ready for Day-1 join.
assert_scylla_wiped

# etcd data
if [[ -d /var/lib/etcd ]]; then
  rm -rf /var/lib/etcd
  log_success "Removed /var/lib/etcd"
fi

# ── keepalived / ingress VIP cleanup ─────────────────────────────────────────
# The keepalived VIP (e.g. the ingress 10.0.0.100) is cluster-managed floating
# state that MUST NOT survive a node wipe. A stale /etc/keepalived/keepalived.conf
# left behind poisons the next install: keepalived raises the old VIP on the
# interface, the node-agent gathers it and — because a fresh day-0 ingress spec
# is "disabled", so lookupIngressVIP() returns "" and the VIP-exclusion guard is
# a no-op — publishes it as the node's PRIMARY identity IP. That IP then fails
# cert SAN coverage (cluster-doctor: security.certs.san_coverage ERROR) and is
# mis-reported as the wired-primary IP. Stop keepalived (releases the VIP) and
# remove its config so the next install starts from a clean, VIP-free identity.
for _ka in keepalived.service globular-keepalived.service; do
  systemctl stop "$_ka" 2>/dev/null || true
  systemctl disable "$_ka" 2>/dev/null || true
done
if [[ -e /etc/keepalived/keepalived.conf ]]; then
  rm -f /etc/keepalived/keepalived.conf
  log_success "Removed /etc/keepalived/keepalived.conf (stale ingress VIP config)"
fi


# ── PKI / Trust store cleanup ─────────────────────────────────────────────────
# Remove all traces of the Globular CA from the system trust store so a
# joining node does not inherit a stale CA. Without this, old CA certs in
# /etc/ssl/certs/ will cause spurious TLS validation failures after CA rotation.

TRUST_CHANGED=0

# /usr/local/share/ca-certificates/ — canonical Debian/Ubuntu location.
# Use a wildcard (not exact filename) because the installer has shipped
# different names over time: globular-ca.crt, globular-root-ca.crt, etc.
# Without the wildcard, update-ca-certificates re-symlinks the leftover
# .crt back into /etc/ssl/certs/ on the next pass.
for cert in /usr/local/share/ca-certificates/*globular* /usr/local/share/ca-certificates/*Globular*; do
  [[ -e "$cert" ]] || continue
  rm -f "$cert"
  TRUST_CHANGED=1
  log_success "Removed $cert"
done

# /etc/ssl/certs/ — symlinks created by update-ca-certificates; also catch any
# manually placed copies. The .0 suffix is the OpenSSL hash-based symlink.
for cert in /etc/ssl/certs/*globular* /etc/ssl/certs/*Globular*; do
  [[ -e "$cert" ]] || continue
  rm -f "$cert"
  TRUST_CHANGED=1
  log_success "Removed $cert"
done

# MinIO TLS artifacts stored outside /var/lib/globular (legacy install paths).
for path in /var/lib/globular/.minio/certs/public.crt \
            /var/lib/globular/.minio/certs/private.key \
            /var/lib/globular/config/tls \
            /var/lib/globular/domains; do
  if [[ -e "$path" ]]; then
    rm -rf "$path"
    log_success "Removed $path"
  fi
done

if [[ $TRUST_CHANGED -eq 1 ]]; then
  update-ca-certificates --fresh >/dev/null 2>&1 || update-ca-certificates >/dev/null 2>&1 || true
  log_success "Rebuilt system CA trust store"
fi

# Remove per-user Globular CA copies and MCP endpoint config so a fresh
# install can regenerate them with the correct new CA and node IP.
for user_home in /root /home/*; do
  [[ -d "$user_home" ]] || continue
  [[ -f "$user_home/.config/globular/ca.crt" ]] && \
    rm -f "$user_home/.config/globular/ca.crt" && \
    log_success "Removed $user_home/.config/globular/ca.crt"
  # Reset MCP endpoint in .mcp.json (remove globular entry, keep others)
  _mcp="$user_home/.claude/.mcp.json"
  if [[ -f "$_mcp" ]] && command -v python3 >/dev/null 2>&1; then
    python3 -c "
import json, sys
try:
    d = json.load(open('$_mcp'))
    d.get('mcpServers', {}).pop('globular', None)
    json.dump(d, open('$_mcp','w'), indent=2)
except Exception:
    pass
"
    log_success "Removed globular MCP entry from $user_home/.claude/.mcp.json"
  fi
done

# User client certificates
for user_home in /home/*; do
  if [[ -d "$user_home/.config/globular" ]]; then
    rm -rf "$user_home/.config/globular"
    log_success "Cleaned certs for $(basename "$user_home")"
  fi
done
[[ -d /root/.config/globular ]] && rm -rf /root/.config/globular && log_success "Cleaned certs for root"

# ── Phase 5: Remove globular user ───────────────────────────────────────────

log_step "Cleanup"

if id globular >/dev/null 2>&1; then
  userdel globular 2>/dev/null || log_warn "Could not remove globular user"
  log_success "Removed globular user"
fi

if getent group globular >/dev/null 2>&1; then
  groupdel globular 2>/dev/null || log_warn "Could not remove globular group"
  log_success "Removed globular group"
fi

# ── Done ──────────────────────────────────────────────────────────────────────

phase_marker "cleanup_complete" "cleanup complete"
echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     ✓ NODE CLEANUP COMPLETE                                    ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "  Node $(hostname) is ready for Day-1 join."
echo "  Removed $REMOVED unit file(s)."
echo ""
