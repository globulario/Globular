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
# Phases 1–5: Local cleanup — stops services and wipes state.
#
# Usage:
#   sudo bash clean-node.sh              # interactive (asks before wiping)
#   sudo bash clean-node.sh --force      # non-interactive (no prompts)
#
# Can be run remotely:
#   ssh user@node "sudo bash -s" < clean-node.sh

FORCE=0
[[ "${1:-}" == "--force" ]] && FORCE=1

die() { echo "  ✗ ERROR: $*" >&2; exit 1; }
log_info() { echo "  → $*"; }
log_success() { echo "  ✓ $*"; }
log_warn() { echo "  ⚠ $*"; }
log_step() { echo ""; echo "━━━ $* ━━━"; echo ""; }

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
    for i in $(seq 1 10); do
        if ! pgrep -af 'scylla' >/dev/null 2>&1; then
            log_success "No ScyllaDB process remains"
            return 0
        fi
        sleep 1
    done

    log_warn "ScyllaDB processes still alive after hard stop:"
    pgrep -af 'scylla' || true
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
[[ $EUID -eq 0 ]] || die "This script must be run as root (use sudo)"

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

# ── Phase 0: Cluster-level detachment ────────────────────────────────────────
#
# Must run BEFORE stopping services: ScyllaDB decommission and etcd member
# remove both require the respective service to be running. Controller removal
# triggers xDS endpoint pruning and MinIO pool eviction automatically.

log_step "Detaching from Cluster (before local wipe)"

_STATE_DIR="/var/lib/globular"
_PKI_DIR="${_STATE_DIR}/pki"
_STATE_FILE="${_STATE_DIR}/nodeagent/state.json"
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
    print(d.get('NodeID', '').strip())
except Exception:
    pass
" 2>/dev/null || true)
fi

# ── 0.1 Remove from cluster controller ───────────────────────────────────────
# Primary: gateway HTTP API (DELETE /api/cluster/nodes/<id>). The gateway uses
# its own controller auth — no user token is required on the cleaning node.
# Fallback: globular CLI (needs a cached token at ~/.config/globular/token).

if [[ -n "$_NODE_ID" ]]; then
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
  _HTTP_STATUS=$(curl -sf -o /dev/null -w "%{http_code}" \
    -X DELETE "https://${_GATEWAY_HOST}:8443/api/cluster/nodes/${_NODE_ID}" \
    -k -H "Content-Type: application/json" \
    -d '{"force":true,"drain":false}' 2>/dev/null || echo "000")

  if [[ "$_HTTP_STATUS" == "200" ]]; then
    log_success "Node removed from cluster controller"
  else
    log_warn "Gateway API returned HTTP ${_HTTP_STATUS} — falling back to globular CLI..."
    if [[ -n "$_GLOBULAR_BIN" ]]; then
      _REMOVE_ERR=$("$_GLOBULAR_BIN" cluster nodes remove "$_NODE_ID" --force --drain=false 2>&1 || true)
      if echo "$_REMOVE_ERR" | grep -q "^message:"; then
        log_success "Node removed from cluster controller (via CLI)"
      else
        log_warn "CLI removal also failed: ${_REMOVE_ERR}"
        log_warn "Run manually after cleanup: globular cluster nodes remove ${_NODE_ID} --force --drain=false"
      fi
    else
      log_warn "globular CLI not found — run manually after cleanup:"
      log_warn "  curl -X DELETE https://${_GATEWAY_HOST}:8443/api/cluster/nodes/${_NODE_ID} -k -d '{\"force\":true,\"drain\":false}'"
    fi
  fi
elif [[ -z "$_NODE_ID" ]]; then
  log_warn "No node ID in ${_STATE_FILE} — skipping controller removal (node may not be registered)"
fi

# ── 0.2 ScyllaDB: decommission before shutdown ───────────────────────────────
# Streams data to remaining peers; must run while scylla-server is still active.
# Skip when this is the only ScyllaDB node (nothing to stream to).
if systemctl is-active --quiet scylla-server.service 2>/dev/null; then
  if command -v nodetool >/dev/null 2>&1; then
    _SCYLLA_UP=$(nodetool status 2>/dev/null | grep -cE "^U[NL] " || echo "0")
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
if systemctl is-active --quiet globular-etcd.service 2>/dev/null \
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
elif systemctl is-active --quiet globular-etcd.service 2>/dev/null; then
  log_warn "etcdctl or TLS certs missing — skipping etcd member removal"
  log_warn "  Manual fix: etcdctl member list → etcdctl member remove <id>"
fi

# ── Phase 1: Stop services ────────────────────────────────────────────────────

log_step "Stopping Services"

# Stop all globular services
for unit in $(systemctl list-units 'globular-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $1}'); do
  log_info "Stopping $unit"
  systemctl stop "$unit" 2>/dev/null || true
  systemctl disable "$unit" 2>/dev/null || true
done

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
for timer in $(systemctl list-timers 'scylla-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $NF}'); do
  log_info "Stopping timer $timer"
  systemctl stop "$timer" 2>/dev/null || true
  systemctl disable "$timer" 2>/dev/null || true
done

# Hard-kill ScyllaDB — must succeed before any wipe begins.
# This is a non-negotiable gate: a live Scylla process can recreate system
# table state in /var/lib/scylla even while the directory is being wiped.
hard_stop_scylla

# ── Phase 2: Force-kill survivors ─────────────────────────────────────────────

log_step "Force-Killing Surviving Processes"

# Kill all globular server processes
for proc in $(ps aux 2>/dev/null | grep -E '_server|globularcli|mcp|gateway|xds_server|envoy' | grep -v grep | awk '{print $2}'); do
  cmd=$(ps -p "$proc" -o comm= 2>/dev/null || true)
  log_warn "Killing PID $proc ($cmd)"
  kill -9 "$proc" 2>/dev/null || true
done

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

log_step "Wiping State"

# Globular state — unconditional rm -rf (safe on missing dirs, avoids
# permission-race with the globular user that was just removed)
# Remove stale Globular wrapper scripts from /usr/local/bin that point
# into /usr/lib/globular/bin (which gets removed below). Without this
# they break system commands like sha256sum after the wipe.
for wrapper in /usr/local/bin/claude /usr/local/bin/codex /usr/local/bin/ffmpeg /usr/local/bin/sctool \
               /usr/local/bin/mc /usr/local/bin/etcdctl /usr/local/bin/globular \
               /usr/local/bin/globularcli /usr/local/bin/restic /usr/local/bin/rclone \
               /usr/local/bin/yt-dlp /usr/local/bin/sha256sum; do
  if [[ -f "$wrapper" ]] && grep -q "usr/lib/globular" "$wrapper" 2>/dev/null; then
    rm -f "$wrapper"
    log_success "Removed stale wrapper $(basename "$wrapper")"
  fi
done

for dir in /var/lib/globular /etc/globular /usr/lib/globular /usr/local/lib/codex; do
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

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║     ✓ NODE CLEANUP COMPLETE                                    ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "  Node $(hostname) is ready for Day-1 join."
echo "  Removed $REMOVED unit file(s)."
echo ""
