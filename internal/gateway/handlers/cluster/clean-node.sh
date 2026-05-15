#!/usr/bin/env bash
set -euo pipefail

# ── Globular Node Cleanup ─────────────────────────────────────────────────────
#
# Prepares a node for a fresh Day-1 join by stopping all Globular and ScyllaDB
# services and removing their state. Run this on any node that previously had
# Globular installed before joining it to a new cluster.
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
  echo "  This will stop all Globular/ScyllaDB services and wipe their data."
  echo "  Press Enter to continue, or Ctrl+C to abort..."
  read -r
fi

# ── Phase 1: Stop services ────────────────────────────────────────────────────

log_step "Stopping Services"

# Stop all globular services
for unit in $(systemctl list-units 'globular-*' --no-pager --no-legend --plain 2>/dev/null | awk '{print $1}'); do
  log_info "Stopping $unit"
  systemctl stop "$unit" 2>/dev/null || true
  systemctl disable "$unit" 2>/dev/null || true
done

# Stop ScyllaDB
for unit in scylla-server.service scylla-node-exporter.service scylla-tune-sched.service; do
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
  log_success "Removed $(basename "$unit_file")"
  REMOVED=$((REMOVED + 1))
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
if dpkg -l 'scylla*' 2>/dev/null | grep -q '^ii'; then
  log_info "Removing ScyllaDB packages (node-agent will reinstall on rejoin)"
  DEBIAN_FRONTEND=noninteractive apt-get remove -y --purge 'scylla*' 2>/dev/null || \
    log_warn "apt remove scylla failed — continuing"
  log_success "ScyllaDB packages removed"
fi

# Wipe all ScyllaDB state and data
for dir in /var/lib/scylla /etc/scylla /etc/scylla.d; do
  if [[ -d "$dir" ]]; then
    rm -rf "$dir"
    log_success "Removed $dir"
  fi
done

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
