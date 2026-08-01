package cluster

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	servicesconfig "github.com/globulario/services/golang/config"
)

// NewJoinScriptHandler serves a self-contained shell script that bootstraps
// a new node to join the cluster. The script follows a strict Day-1 protocol:
//
//	Phase 0: Preflight       — token, identity, lifecycle classification;
//	                           every refusal fires BEFORE any service stop, and
//	                           a compensation trap restores previously-active
//	                           units if the join aborts mid-flight
//	Phase 1: Directories     — create layout and globular user (no state wipe)
//	Phase 2: PKI/Identity    — CA cert, JoinPlan authorization gate, service cert (join_id-gated), TLS verification
//	Phase 3: Connectivity    — DNS, /etc/hosts, libnss, journald, config.json, node-agent state
//	Phase 4: Download        — release tarball (GitHub first, controller-local fallback): installer, packages, workflows
//	Phase 6: ScyllaDB        — install and verify local storage join (only when the plan assigns 'storage')
//	Phase 7: control plane   — xds/envoy/gateway always; workflow/repository/controller only when the plan assigns 'control-plane'
//	Phase 8: node-agent      — install and start; auto-initiates join RPC via state file
//	Phase 9: MinIO hold      — enforce objectstore topology contract
//
// The JoinPlan authorization runs BEFORE certificate signing: the signed cert
// is an mTLS credential (etcd trusts the same CA), so the gateway only signs
// CSRs that present a controller-issued join_id. The plan's assigned_profiles
// drive which local packages Phases 6–7 install — the controller owns
// placement; the script must not install storage or control-plane services
// on a node whose plan does not assign those profiles.
//
// Phase 4 stages the active BOM packages and installs the etcd binary through
// globular-installer. Phase 5 performs the cluster membership operation
// explicitly: ghost cleanup, member add, etcd.yaml, start, and anti-fork
// membership verification.
//
// Repair mode (etcd WAL damage):
//
//	curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <TOKEN> --repair-etcd
//
// Usage (fresh Day-1 join):
//
//	curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <join-token>
//
// activeReleaseAnchorKey is the etcd key holding the cluster's authoritative
// active platform release, written by the controller (ActivatePlatformRelease /
// convergence auto-advance; see services active_release.go). It is the single
// source of truth for "what release is the cluster actually running".
const activeReleaseAnchorKey = "/globular/platform/active_release"

// activePlatformVersion resolves the platform version a joining (Day-1) node
// should bootstrap, in strict priority order:
//
//  1. the etcd active_release anchor — cluster authority. Day-1 must bootstrap
//     the version the cluster is actually running.
//  2. the local release-index.json projection (written by 'globular repo sync')
//     — a node-local copy that can be stale (e.g. the gateway node has not
//     re-synced since the last platform activation). Trusting it previously made
//     a joining node bootstrap an OLD node-agent that never converged.
//  3. the gateway binary's own version — Day-0 bootstrap, before any anchor or
//     repo sync exists.
//
// The anchor is the authority; the on-disk index is only a projection of it, so
// the projection must not be trusted while the authority is reachable.
func activePlatformVersion(binaryVersion string) string {
	return resolvePlatformVersion(binaryVersion, platformReleaseFromAnchor, os.ReadFile)
}

// resolvePlatformVersion is the pure precedence logic, with the etcd read and
// the file read injected so it is unit-testable without cluster infrastructure.
func resolvePlatformVersion(binaryVersion string, anchorFn func() string, readFile func(string) ([]byte, error)) string {
	// 1. etcd active_release anchor (authority).
	if v := anchorFn(); v != "" {
		return v
	}
	// 2. local release-index.json projection (may be stale).
	if data, err := readFile("/var/lib/globular/release-index.json"); err == nil {
		var idx struct {
			PlatformRelease string `json:"platform_release"`
		}
		if jsonErr := json.Unmarshal(data, &idx); jsonErr == nil && idx.PlatformRelease != "" {
			return idx.PlatformRelease
		}
	}
	// 3. the gateway binary's own version (Day-0).
	return binaryVersion
}

// platformReleaseFromAnchor reads platform_release from the etcd active_release
// anchor. It returns "" on any error (etcd unreachable, key absent, corrupt
// JSON) so the caller falls back to the on-disk projection / binary version —
// Day-0 bootstrap (no etcd yet) must still work. The shared etcd client is a
// cached singleton owned by the config package and must not be closed here.
func platformReleaseFromAnchor() string {
	client, err := servicesconfig.GetEtcdClient()
	if err != nil || client == nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := client.Get(ctx, activeReleaseAnchorKey)
	if err != nil || len(resp.Kvs) == 0 {
		return ""
	}
	var a struct {
		PlatformRelease string `json:"platform_release"`
	}
	if jsonErr := json.Unmarshal(resp.Kvs[0].Value, &a); jsonErr != nil {
		return ""
	}
	return a.PlatformRelease
}

func NewJoinScriptHandler(controllerAddr string, gatewayPort int, platformVersion string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read CA certificate to embed in the script.
		caCert, err := os.ReadFile("/var/lib/globular/pki/ca.crt")
		if err != nil {
			http.Error(w, "CA certificate not found", http.StatusInternalServerError)
			return
		}
		caB64 := base64.StdEncoding.EncodeToString(caCert)

		// Determine the gateway address the client connected to.
		gatewayHost := r.Host
		if gatewayHost == "" {
			gatewayHost = r.URL.Host
		}

		reqHost, reqPort, err := net.SplitHostPort(gatewayHost)
		if err != nil {
			reqHost = gatewayHost
			reqPort = fmt.Sprintf("%d", gatewayPort)
		}
		if reqHost == "" || reqHost == "127.0.0.1" || reqHost == "::1" || reqHost == "localhost" {
			// Never embed a loopback address — the joining node is remote.
			if h, _, splitErr := net.SplitHostPort(r.RemoteAddr); splitErr == nil {
				reqHost = h
			}
		}

		gatewayAddr := net.JoinHostPort(reqHost, reqPort)

		ctrlAddr := controllerAddr
		ctrlHost, _, _ := net.SplitHostPort(ctrlAddr)
		if ctrlAddr == "" || ctrlHost == "127.0.0.1" || ctrlHost == "localhost" || ctrlHost == "::1" {
			ctrlAddr = fmt.Sprintf("%s:12000", reqHost)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, joinScript, caB64, gatewayAddr, ctrlAddr, activePlatformVersion(platformVersion), servicesconfig.EtcdClusterToken)
	})
}

// joinScript is the shell script template served to joining nodes.
// Format args (all positional to avoid ambiguity with shell ${} syntax):
//
//	%[1]s = base64 CA cert
//	%[2]s = gateway address (host:port)
//	%[3]s = controller address (host:port)
//	%[4]s = platform version (e.g. "1.2.88") — used to build the release tarball URL(s)
//	%[5]s = etcd initial-cluster-token (config.EtcdClusterToken, the fixed bootstrap-immutable constant)
//
// IMPORTANT: any literal %% in this string becomes % in the output script.
// Use %% wherever the shell script needs a literal % (e.g. date format strings).
const joinScript = `#!/usr/bin/env bash
set -euo pipefail

# ============================================================================
# Globular Day-1 Node Join Script
# Generated by the bootstrap gateway — do not edit manually.
#
# Phases:
#   0 — Preflight:   token required, identity detection, lifecycle
#                     classification — every refusal fires BEFORE any service
#                     stop; a compensation trap restores previously-active
#                     units and writes a receipt if the join aborts
#   1 — Directories: create layout, globular user (no state wipe)
#   2 — PKI/Identity: CA cert, JoinPlan authorization gate, service cert
#                     (signed only with a valid join_id), TLS smoke-test
#   3 — Connectivity: DNS, /etc/hosts MinIO fallback, libnss, journald,
#                     etcd_endpoints seed, node-agent state.json (join_id)
#   4 — Download:   release tarball (GitHub first, controller fallback) → globular-installer, packages, workflows
#   5 — Etcd join:  ghost cleanup, member add, etcd.yaml, start, health gate, fork-detection
#   6 — ScyllaDB:   install from package, join local storage, wait for CQL
#                   (only when the JoinPlan assigns the 'storage' profile)
#   7 — Control plane: xds/envoy/gateway always (mesh data plane); workflow/
#                   repository/controller only with the 'control-plane' profile
#   8 — node-agent: write unit (After=globular-etcd), start; uses v2 path (join_id)
#   9 — MinIO hold: enforce objectstore topology contract
#
# Phase 2.2 is the JoinPlan authorization gate. The gateway is a courier only —
# the controller assigns profiles, etcd intent, node_id, and signs the plan.
# No cluster-affecting step runs before it succeeds — including certificate
# signing (Phase 2.3): the CA-signed cert is an etcd mTLS credential, so the
# gateway refuses to sign a CSR without a controller-issued join_id.
# The plan's assigned_profiles drive Phases 6–7: the controller owns placement.
#
# The node-agent reads join_token + controller_endpoint + join_id from state
# and auto-initiates the v2 join path on startup (no RequestJoin call needed).
# After the join is approved the controller drives all further convergence:
#   admitted → infra_preparing → etcd_ready → xds_ready →
#   envoy_ready → storage_joining → workload_ready
#
# Repair mode (etcd WAL corrupt or permanently removed from cluster):
#   ... | sudo bash -s -- --token <TOKEN> --repair-etcd
#   This mode REQUIRES explicit operator intent. It backs up and wipes
#   /var/lib/globular/etcd before rejoining. Never wiped without this flag.
# ============================================================================

# ── Embedded cluster parameters ─────────────────────────────────────────────
CA_B64="%[1]s"
GATEWAY="%[2]s"
CONTROLLER="%[3]s"
ETCD_CLUSTER_TOKEN="%[5]s"

# ── Runtime defaults ─────────────────────────────────────────────────────────
JOIN_TOKEN=""
REPAIR_ETCD=0
PROFILES=""
INSTALL_DIR="/usr/lib/globular/bin"
STATE_DIR="/var/lib/globular"
SYSTEMD_DIR="/etc/systemd/system"
PACKAGES_DIR="${STATE_DIR}/packages"

# ── Compensation state ───────────────────────────────────────────────────────
# Armed between the first service stop (Phase 0.1) and proof of the new
# cluster membership (Phase 5.7 anti-fork check). While armed, any abort
# restores the units that were active when the script started and writes a
# durable interrupted-operation receipt. See compensate_on_abort in Phase 0.
JOIN_PHASE_MARKER="startup"
COMPENSATION_ARMED=0
RESTORE_UNITS_ON_ABORT=1
ETCD_WAS_ACTIVE=0
AGENT_WAS_ACTIVE=0

# ── Argument parsing ─────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    --token)       JOIN_TOKEN="$2"; shift 2 ;;
    --controller)  CONTROLLER="$2"; shift 2 ;;
    --repair-etcd) REPAIR_ETCD=1; shift ;;
    # --profiles is an operator preference (comma-separated, e.g.
    # core,control-plane,storage,dns). The controller VALIDATES and decides;
    # empty uses suggested/inherited defaults.
    --profiles)    PROFILES="$2"; shift 2 ;;
    *) echo "Unknown argument: $1" >&2; exit 1 ;;
  esac
done

[[ -n "${JOIN_TOKEN}" ]] || {
  echo "Error: --token is required" >&2
  echo "Usage: curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <TOKEN>" >&2
  exit 1
}

# ── Helpers ──────────────────────────────────────────────────────────────────
die()       { echo "  FAIL: $*" >&2; exit 1; }
log_ok()    { echo "  OK:   $*"; }
log_info()  { echo "  ...   $*"; }
log_warn()  { echo "  WARN: $*" >&2; }
log_phase() { JOIN_PHASE_MARKER="$*"; echo ""; echo "══════════════════════════════════════════════════════"; echo "PHASE $*"; echo "══════════════════════════════════════════════════════"; echo ""; }

find_pkg() {
  local name="$1"
  find "${PACKAGES_DIR}" -maxdepth 1 -name "${name}_*_linux_amd64.tgz" 2>/dev/null | sort | tail -1
}

install_pkg() {
  local name="$1"
  local pkg
  pkg=$(find_pkg "${name}")
  [[ -n "${pkg}" ]] || die "package ${name} not found in ${PACKAGES_DIR}"
  log_info "Installing ${name} via globular-installer..."
  "${INSTALL_DIR}/globular-installer" install \
    --non-interactive --verbose \
    "${pkg}" \
    || die "installer: ${name} install failed"
  log_ok "${name} installed"
}

wait_unit_active() {
  local unit="$1"
  local seconds="$2"
  local label="$3"
  log_info "Waiting for ${label} (${unit}) to become active (max ${seconds}s)..."
  for i in $(seq 1 "${seconds}"); do
    if systemctl is-active --quiet "${unit}" 2>/dev/null; then
      log_ok "${label} active"
      return 0
    fi
    sleep 1
  done
  return 1
}

wait_tcp_port() {
  local host="$1"
  local port="$2"
  local seconds="$3"
  local label="$4"
  log_info "Waiting for ${label} on ${host}:${port} (max ${seconds}s)..."
  for i in $(seq 1 "${seconds}"); do
    if timeout 1 bash -c "</dev/tcp/${host}/${port}" 2>/dev/null; then
      log_ok "${label} reachable on ${host}:${port}"
      return 0
    fi
    sleep 1
  done
  return 1
}

[[ "$(id -u)" -eq 0 ]] || die "This script must be run as root (sudo)"

BOOTSTRAP_HOST=$(echo "${GATEWAY}" | cut -d: -f1)
BOOTSTRAP_ETCD="https://${BOOTSTRAP_HOST}:2379"
NODE_HOSTNAME=$(hostname)
NODE_IP=$(hostname -I | awk '{print $1}')
ETCD_NAME="${NODE_HOSTNAME}"
PKI_DIR="${STATE_DIR}/pki/issued/services"
ETCD_CACERT="${STATE_DIR}/pki/ca.crt"
ETCD_CERT="${PKI_DIR}/service.crt"
ETCD_KEY="${PKI_DIR}/service.key"

echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║          GLOBULAR DAY-1 NODE JOIN                               ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo "  Controller: ${CONTROLLER}"
echo "  Gateway:    ${GATEWAY}"
echo "  Bootstrap:  ${BOOTSTRAP_HOST}"
echo "  This node:  ${NODE_HOSTNAME} (${NODE_IP})"
if [[ ${REPAIR_ETCD} -eq 1 ]]; then
  echo "  Mode:       REPAIR-ETCD (WAL backup + wipe authorized)"
else
  echo "  Mode:       fresh Day-1 join"
fi
echo ""

# ============================================================================
log_phase "0 — Preflight"
# ============================================================================
# LIFECYCLE-FIRST, MUTATION-LAST. Every check that can refuse this join runs
# BEFORE any service is stopped: an aborted preflight must leave the node
# exactly as it found it. (Scar 2026-07-10: re-running this script on an
# already-joined node stopped its etcd member during preflight, THEN refused
# on the existing data dir — leaving the member down and collapsing quorum
# on a 2-node cluster.)
# Stop ONLY the services that must be reset for this join sequence.
# Do not wildcard-stop all globular-* services.

log_info "[0.0] Classifying node lifecycle state (read-only)..."
systemctl is-active --quiet globular-etcd.service 2>/dev/null && ETCD_WAS_ACTIVE=1
systemctl is-active --quiet globular-node-agent.service 2>/dev/null && AGENT_WAS_ACTIVE=1

ETCD_DATA_PRESENT=0
if [[ -d "${STATE_DIR}/etcd" ]] && [[ -n "$(ls -A "${STATE_DIR}/etcd" 2>/dev/null)" ]]; then
  ETCD_DATA_PRESENT=1
fi

ADMITTED_NODE_ID=""
if [[ -f "${STATE_DIR}/nodeagent/state.json" ]]; then
  ADMITTED_NODE_ID=$(grep -oE '"node_id"[[:space:]]*:[[:space:]]*"[^"]+"' "${STATE_DIR}/nodeagent/state.json" 2>/dev/null \
    | head -1 | sed 's/.*:[[:space:]]*"\([^"]*\)"/\1/' || true)
fi

LIFECYCLE="PRISTINE"
if [[ ${ETCD_DATA_PRESENT} -eq 1 || -n "${ADMITTED_NODE_ID}" ]]; then
  if [[ -n "${ADMITTED_NODE_ID}" && ${ETCD_WAS_ACTIVE} -eq 1 && ${AGENT_WAS_ACTIVE} -eq 1 ]]; then
    LIFECYCLE="JOINED"
  else
    LIFECYCLE="PARTIALLY_JOINED"
  fi
fi
log_info "lifecycle=${LIFECYCLE} etcd_data=${ETCD_DATA_PRESENT} node_id=${ADMITTED_NODE_ID:-<none>} etcd_active=${ETCD_WAS_ACTIVE} agent_active=${AGENT_WAS_ACTIVE}"

# Fresh join is ONLY a fresh-join tool — it is not a retry mechanism.
# Refuse a non-pristine node BEFORE touching any service.
if [[ ${REPAIR_ETCD} -eq 0 && "${LIFECYCLE}" != "PRISTINE" ]]; then
  die "REFUSED: node is ${LIFECYCLE} — it already carries cluster state
  (etcd data directory exists at ${STATE_DIR}/etcd and --repair-etcd was not given,
  admitted node_id: ${ADMITTED_NODE_ID:-<none>}). No service was stopped.
  Fresh join is not a retry mechanism. Use instead:
    - node joined but degraded or its units stopped: restart the substrate
        systemctl start globular-etcd globular-node-agent
      and let the controller reconcile ('globular cluster nodes list')
    - prior join failed partway: wipe first with the clean script
        curl -sfL https://<gateway>:8443/clean -k | sudo bash -s -- --force
      then re-run this join on the pristine node
    - etcd WAL damage on THIS node only: re-run with --repair-etcd
      (authorizes a backup + wipe of the local etcd data dir)"
fi

# Repair mode presupposes a broken local member: wiping a LIVE member's data
# is the wrong operation regardless of intent.
if [[ ${REPAIR_ETCD} -eq 1 && ${ETCD_WAS_ACTIVE} -eq 1 ]]; then
  die "REFUSED: --repair-etcd given but globular-etcd.service is ACTIVE.
  Repair mode backs up and wipes the local etcd data dir — running it against
  a live member destroys healthy cluster state. If the member is genuinely
  broken, stop it explicitly first (systemctl stop globular-etcd), then
  re-run with --repair-etcd. No service was stopped."
fi

# From the first service stop until the new membership is proven (Phase 5.7
# anti-fork check), any abort restores the units that were active at start
# and writes a durable interrupted-operation receipt. After the repair-mode
# wipe the pre-join etcd state no longer exists, so unit restore is disarmed
# at that point (the receipt is still written).
restore_unit_after_abort() {
  local unit="$1"
  if systemctl start "${unit}" 2>/dev/null; then
    log_warn "compensation: ${unit} restarted (was active before this join)"
  else
    log_warn "compensation FAILED: could not restart ${unit} — start it manually"
  fi
}

compensate_on_abort() {
  local rc=$?
  [[ ${rc} -ne 0 && ${COMPENSATION_ARMED} -eq 1 ]] || return 0
  log_warn "join aborted during phase '${JOIN_PHASE_MARKER}' (exit ${rc}) — compensating"
  mkdir -p "${STATE_DIR}/nodeagent" 2>/dev/null || true
  cat > "${STATE_DIR}/nodeagent/join_interrupted.json" <<RECEIPT
{
  "operation": "day1-join",
  "aborted_phase": "${JOIN_PHASE_MARKER}",
  "exit_code": ${rc},
  "repair_mode": ${REPAIR_ETCD},
  "etcd_was_active": ${ETCD_WAS_ACTIVE},
  "node_agent_was_active": ${AGENT_WAS_ACTIVE},
  "units_restore_attempted": ${RESTORE_UNITS_ON_ABORT},
  "timestamp": "$(date -Is)"
}
RECEIPT
  if [[ ${RESTORE_UNITS_ON_ABORT} -eq 1 ]]; then
    # Restore via a parameterized helper so the literal unit-start commands
    # exist only in Phase 8 (ordering contracts assert on those literals).
    if [[ ${ETCD_WAS_ACTIVE} -eq 1 ]]; then
      restore_unit_after_abort "globular-etcd.service"
    fi
    if [[ ${AGENT_WAS_ACTIVE} -eq 1 ]]; then
      restore_unit_after_abort "globular-node-agent.service"
    fi
  else
    log_warn "compensation: unit restore not applicable at this phase — receipt at ${STATE_DIR}/nodeagent/join_interrupted.json"
  fi
}
trap compensate_on_abort EXIT
COMPENSATION_ARMED=1

log_info "[0.1] Stopping globular-etcd and globular-node-agent if running..."
for unit in globular-etcd.service globular-node-agent.service; do
  if systemctl is-active --quiet "${unit}" 2>/dev/null; then
    log_info "Stopping ${unit}..."
    systemctl stop "${unit}" 2>/dev/null || true
  fi
done

# Repair mode: back up etcd WAL then wipe — ONLY with explicit --repair-etcd.
# The no-flag case was already refused in [0.0] BEFORE any service stop.
if [[ ${REPAIR_ETCD} -eq 1 ]]; then
  log_info "[0.2] Repair mode: backing up etcd data before wipe..."
  if [[ -d "${STATE_DIR}/etcd" ]]; then
    ETCD_BACKUP="${STATE_DIR}/etcd.bak.$(date +%%Y%%m%%dT%%H%%M%%S)"
    cp -a "${STATE_DIR}/etcd" "${ETCD_BACKUP}" \
      || die "etcd backup to ${ETCD_BACKUP} failed — refusing to wipe without a successful backup"
    log_ok "etcd backup written: ${ETCD_BACKUP}"
    rm -rf "${STATE_DIR}/etcd"
    # Point of no return: the pre-join etcd state is gone. Restarting the old
    # unit on an empty data dir could fork a new cluster — disarm unit restore
    # (the interrupted-operation receipt is still written on abort).
    RESTORE_UNITS_ON_ABORT=0
    log_ok "etcd data dir wiped (backup at ${ETCD_BACKUP})"
  else
    log_info "no etcd data dir found — nothing to wipe"
  fi
fi

log_ok "Preflight complete"

# ============================================================================
log_phase "1 — Directories and globular user"
# ============================================================================
# Create directory layout and the globular system user.
# No etcd.yaml, no PKI material, no node identity is written here.
# State wipe already handled in Phase 0 (repair only) or by clean-node.sh.

log_info "[1.1] Creating directories..."
mkdir -p "${INSTALL_DIR}" \
         "${STATE_DIR}/pki/issued/services" \
         "${STATE_DIR}/services" \
         "${STATE_DIR}/node_agent" \
         "${STATE_DIR}/config" \
         "${STATE_DIR}/workflows" \
         "${STATE_DIR}/keys"

if ! id -u globular &>/dev/null; then
  useradd --system --no-create-home --home-dir "${STATE_DIR}" --shell /usr/sbin/nologin globular
fi
log_ok "Directories and user ready"

# ============================================================================
log_phase "2 — PKI / Identity"
# ============================================================================

log_info "[2.1] Installing cluster CA certificate..."
echo "${CA_B64}" | base64 -d > "${STATE_DIR}/pki/ca.crt"
cp "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
chmod 644 "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"

# Install into system trust store so curl and Go TLS work with the cluster CA.
cp "${STATE_DIR}/pki/ca.crt" /usr/local/share/ca-certificates/globular-ca.crt
update-ca-certificates >/dev/null 2>&1 || true
log_ok "CA certificate installed and trusted"

log_info "[2.2] Obtaining signed JoinPlan from controller (v2 authorization gate)..."
# ────────────────────────────────────────────────────────────────────────────
# The gateway is a courier. The controller is the authority.
# This call must succeed before any cluster-affecting step — including the
# service-certificate signing below: a CA-signed cert is an etcd mTLS
# credential, so the gateway only signs CSRs presenting a valid join_id.
# Cluster-affecting steps: certificate issuance, etcd member add, writing
# cluster membership config, starting node-agent as a joined node, applying
# assigned profiles.
# ────────────────────────────────────────────────────────────────────────────
mkdir -p "${STATE_DIR}/nodeagent"

AUTHORIZE_RESPONSE="${STATE_DIR}/nodeagent/join_authorize.json"

# Build requested_profiles JSON array from the comma-separated --profiles value
# (operator preference; the controller validates and decides). Empty when unset.
REQUESTED_PROFILES_JSON=""
if [[ -n "${PROFILES}" ]]; then
  _pjson=""
  IFS=',' read -ra _parr <<< "${PROFILES}"
  for _p in "${_parr[@]}"; do
    _p=$(echo "${_p}" | xargs)
    [[ -z "${_p}" ]] && continue
    [[ -n "${_pjson}" ]] && _pjson="${_pjson},"
    _pjson="${_pjson}\"${_p}\""
  done
  [[ -n "${_pjson}" ]] && REQUESTED_PROFILES_JSON=",\"requested_profiles\":[${_pjson}]"
fi

CPU_COUNT=$(nproc 2>/dev/null || echo 0)
RAM_BYTES=$(awk '/MemTotal:/ {print $2 * 1024}' /proc/meminfo 2>/dev/null || echo 0)
DISK_BYTES=$(df -B1 / 2>/dev/null | awk 'NR==2 {print $2}')
DISK_FREE_BYTES=$(df -B1 / 2>/dev/null | awk 'NR==2 {print $4}')
CPU_COUNT=${CPU_COUNT:-0}
RAM_BYTES=${RAM_BYTES:-0}
DISK_BYTES=${DISK_BYTES:-0}
DISK_FREE_BYTES=${DISK_FREE_BYTES:-0}
CAPABILITIES_JSON=$(printf ',"capabilities":{"cpu_count":%%s,"ram_bytes":%%s,"disk_bytes":%%s,"disk_free_bytes":%%s,"can_apply_privileged":true,"privilege_reason":"join script running as root"}' \
  "${CPU_COUNT}" "${RAM_BYTES}" "${DISK_BYTES}" "${DISK_FREE_BYTES}")

# Build the authorization request JSON (no bash arrays; use printf for safety)
AUTHORIZE_BODY=$(printf '{"join_token":"%%s","identity":{"hostname":"%%s","ips":["%%s"]}%%s%%s}' \
  "${JOIN_TOKEN}" "${NODE_HOSTNAME}" "${NODE_IP}" "${REQUESTED_PROFILES_JSON}" "${CAPABILITIES_JSON}")

# Use -w to capture HTTP status separately; do NOT use -f so the response body
# is always written to AUTHORIZE_RESPONSE regardless of HTTP status code.
# This lets us extract denied_reason from a 403 instead of hiding it.
AUTHORIZE_HTTP_STATUS=$(curl -sL --cacert "${STATE_DIR}/pki/ca.crt" \
  -X POST "https://${GATEWAY}/join/authorize" \
  -H "Content-Type: application/json" \
  -d "${AUTHORIZE_BODY}" \
  -o "${AUTHORIZE_RESPONSE}" \
  -w "%%{http_code}" 2>/dev/null) \
  || die "join blocked: network error reaching gateway /join/authorize"

if [[ "${AUTHORIZE_HTTP_STATUS}" -ge 500 ]]; then
  GATEWAY_ERR=$(head -1 "${AUTHORIZE_RESPONSE}" 2>/dev/null || true)
  die "join blocked: gateway error (HTTP ${AUTHORIZE_HTTP_STATUS}) — ${GATEWAY_ERR}"
fi

# Parse allowed field: handles both "allowed":true and "allowed": true
JOIN_PLAN_ALLOWED=$(grep -oE '"allowed"\s*:\s*(true|false)' "${AUTHORIZE_RESPONSE}" \
  | head -1 | grep -oE '(true|false)$' || echo "false")

# Extract join_id and denied_reason from response JSON
JOIN_ID=$(grep -oE '"join_id"\s*:\s*"[^"]*"' "${AUTHORIZE_RESPONSE}" \
  | head -1 | sed 's/.*"join_id"\s*:\s*"\([^"]*\)".*/\1/' || echo "")

DENIED_REASON=$(grep -oE '"denied_reason"\s*:\s*"[^"]*"' "${AUTHORIZE_RESPONSE}" \
  | head -1 | sed 's/.*"denied_reason"\s*:\s*"\([^"]*\)".*/\1/' || echo "unknown")

[[ "${JOIN_PLAN_ALLOWED}" == "true" ]] \
  || die "join blocked: controller denied — ${DENIED_REASON}"

[[ -n "${JOIN_ID}" ]] \
  || die "join blocked: signed JoinPlan missing — controller returned empty join_id"

# Extract the controller-assigned profiles from the signed plan. These drive
# Phases 6–7: the controller owns placement, the script only executes it.
# Comma-separated, e.g. "control-plane,core,storage". Empty when the plan
# carries no profiles (older controller) — Phases 6–7 then fall back to the
# legacy full local install so upgrades stay safe.
ASSIGNED_PROFILES=$(grep -oE '"assigned_profiles"\s*:\s*\[[^]]*\]' "${AUTHORIZE_RESPONSE}" \
  | head -1 | sed 's/.*\[\(.*\)\]/\1/' | tr -d '" ' || echo "")

if [[ -z "${ASSIGNED_PROFILES}" ]]; then
  log_warn "JoinPlan carries no assigned_profiles — falling back to full local install (storage + control plane)"
else
  log_info "controller-assigned profiles: ${ASSIGNED_PROFILES}"
fi

# should_install_role <profile> — true when the plan assigns the profile, or
# when the plan carries no profiles at all (legacy fallback: install everything).
should_install_role() {
  [[ -z "${ASSIGNED_PROFILES}" ]] && return 0
  [[ ",${ASSIGNED_PROFILES}," == *",$1,"* ]]
}

log_ok "join plan accepted: join_id=${JOIN_ID} profiles=${ASSIGNED_PROFILES:-<legacy: all>}"

log_info "[2.3] Generating service certificate..."
openssl genrsa -out "${PKI_DIR}/service.key" 2048 2>/dev/null

cat > /tmp/globular-san.cnf <<SANEOF
[req]
req_extensions     = v3_req
distinguished_name = req_dn
prompt             = no
[req_dn]
CN = ${NODE_HOSTNAME}
O  = globular.internal
[v3_req]
subjectAltName = DNS:${NODE_HOSTNAME},DNS:${NODE_HOSTNAME}.globular.internal,DNS:localhost,IP:${NODE_IP},IP:127.0.0.1
SANEOF

openssl req -new -key "${PKI_DIR}/service.key" \
  -config /tmp/globular-san.cnf \
  -out /tmp/globular-join.csr 2>/dev/null

CSR_B64=$(base64 -w0 /tmp/globular-join.csr)
# join_id authenticates this request: the gateway refuses to sign CSRs from
# requesters that are neither cluster members nor authorized Day-1 joiners.
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/sign_ca_certificate?csr=${CSR_B64}&join_id=${JOIN_ID}" \
  -o "${PKI_DIR}/service.crt"
rm -f /tmp/globular-join.csr /tmp/globular-san.cnf

[[ -s "${PKI_DIR}/service.crt" ]] \
  || die "Certificate signing failed — controller returned empty certificate"

chmod 600 "${PKI_DIR}/service.key"
chmod 644 "${PKI_DIR}/service.crt"
log_ok "Service certificate issued for ${NODE_HOSTNAME} (${NODE_IP})"

log_info "[2.4] Verifying TLS access to bootstrap gateway..."
if curl -sf --cacert "${STATE_DIR}/pki/ca.crt" -o /dev/null \
  "https://${GATEWAY}/get_ca_certificate" 2>/dev/null; then
  log_ok "Bootstrap gateway reachable via TLS"
else
  log_warn "Gateway TLS check failed — DNS may not be set yet; continuing"
fi

# ============================================================================
log_phase "3 — Connectivity"
# ============================================================================

log_info "[3.1] Configuring DNS resolver (*.globular.internal → bootstrap)..."
mkdir -p /etc/systemd/resolved.conf.d
cat > /etc/systemd/resolved.conf.d/globular.conf <<DNSCONF
[Resolve]
DNS=${BOOTSTRAP_HOST}
Domains=~globular.internal
DNSCONF
systemctl restart systemd-resolved 2>/dev/null || true
log_ok "DNS resolver → ${BOOTSTRAP_HOST} for globular.internal"

# -- MinIO /etc/hosts fallback (objectstore topology contract) ---------------
#
# This node is NOT automatically a MinIO/objectstore member on Day-1.
# Objectstore membership is governed solely by ObjectStoreDesiredState in
# etcd, written by the controller after explicit 'globular objectstore
# apply-topology'. No topology mutation happens here.
#
# The bootstrap node is written as the ONLY /etc/hosts MinIO entry so that
# artifact fetches resolve before DNS propagates. This node's own IP is
# never written as a MinIO endpoint by the join script.
log_info "[3.2] Adding MinIO /etc/hosts fallback (bootstrap objectstore only)..."
sed -i '/minio\.globular\.internal/d' /etc/hosts
echo "${BOOTSTRAP_HOST} minio.globular.internal  # MinIO bootstrap fallback — ObjectStoreDesiredState governs topology" >> /etc/hosts
log_ok "minio.globular.internal → ${BOOTSTRAP_HOST} (bootstrap; topology contract governs membership)"

log_info "[3.3] Ensuring libnss-resolve (systemd-resolved bridge for Go binaries)..."
if ! dpkg -l libnss-resolve 2>/dev/null | grep -q '^ii'; then
  apt-get install -y libnss-resolve >/dev/null 2>&1 || true
fi
if grep -q 'mdns4_minimal.*NOTFOUND.*return' /etc/nsswitch.conf; then
  sed -i 's/^hosts:.*/hosts:          files resolve dns myhostname/' /etc/nsswitch.conf
fi
log_ok "libnss-resolve configured"

log_info "[3.4] Configuring log rate limits (prevent syslog/journal disk fill)..."
mkdir -p /etc/rsyslog.d
cat > /etc/rsyslog.d/50-globular-rate-limit.conf <<'RSYSLOG'
$SystemLogRateLimitInterval 10
$SystemLogRateLimitBurst 500
RSYSLOG
systemctl restart rsyslog 2>/dev/null || true

mkdir -p /etc/systemd/journald.conf.d
cat > /etc/systemd/journald.conf.d/50-globular.conf <<'JOURNALD'
[Journal]
SystemMaxUse=2G
SystemKeepFree=10G
MaxRetentionSec=7day
JOURNALD
systemctl restart systemd-journald 2>/dev/null || true
log_ok "rsyslog rate-limited, journald capped at 2GB/7days"

log_info "[3.5] Writing node-agent seed config.json..."
cat > "${STATE_DIR}/config.json" <<CFGJSON
{
  "Address": "${NODE_IP}",
  "AlternateDomains": [],
  "Domain": "globular.internal",
  "EnablePeerUpserts": false,
  "EtcdClientPort": "2379",
  "EtcdConfigPath": "${STATE_DIR}/config/etcd.yaml",
  "EtcdDataDir": "${STATE_DIR}/etcd",
  "EtcdEnabled": true,
  "EtcdMode": "existing",
  "EtcdName": "${NODE_HOSTNAME}",
  "EtcdPeerPort": "2380",
  "MutateHostsFile": false,
  "MutateResolvConf": false,
  "Name": "${NODE_HOSTNAME}",
  "Peers": [],
  "PortHTTP": 8080,
  "PortHTTPS": 8443,
  "Protocol": "https",
  "Services": {}
}
CFGJSON

chown -R globular:globular "${STATE_DIR}"
chown globular:globular "${STATE_DIR}/keys"

log_info "[3.6] Writing etcd endpoints file (seed for ScyllaDB post-install)..."
# ScyllaDB's post-install reads this file to discover seed nodes.
# Without it the file is absent and ScyllaDB falls back to the local IP only,
# bootstrapping a standalone cluster instead of joining the existing ring.
# The node-agent's recovery path only writes this file after 3 heartbeat
# failures, which is too late — it must exist before storage_joining triggers.
echo "https://${BOOTSTRAP_HOST}:2379" > "${STATE_DIR}/config/etcd_endpoints"
chown globular:globular "${STATE_DIR}/config/etcd_endpoints"
log_ok "etcd_endpoints seeded with bootstrap host (${BOOTSTRAP_HOST})"

log_info "[3.7] Writing node-agent state (join token + controller endpoint + join_id)..."
# The join_id is included so the node-agent uses the v2 path (no RequestJoin call):
# it polls GetJoinRequestStatus(join_id) directly without consuming the token again.
cat > "${STATE_DIR}/nodeagent/state.json" <<STATEJSON
{
  "controller_endpoint": "${CONTROLLER}",
  "controller_insecure": false,
  "request_id": "",
  "node_id": "",
  "join_token": "${JOIN_TOKEN}",
  "join_id": "${JOIN_ID}",
  "network_generation": 0,
  "cluster_domain": "globular.internal",
  "protocol": "https"
}
STATEJSON
chown -R globular:globular "${STATE_DIR}/nodeagent"
log_ok "node-agent state seeded (join_token set, join_id=${JOIN_ID}, controller=${CONTROLLER})"

log_ok "Connectivity configured"

# ============================================================================
log_phase "4 — Download Release Tarball"
# ============================================================================
# Use the same release tarball shape as Day-0. It contains globular-installer,
# all packages, and all workflow definitions.
#
# Normal path: download the published tarball from GitHub.
# Dev/local-build fallback: if GitHub does not have v${PLATFORM_VERSION}, pull a
# controller-local bootstrap bundle from the gateway. The controller's active
# BOM remains the authority for the version; the fallback only changes the
# transport.

PLATFORM_VERSION="%[4]s"

if [[ -z "${PLATFORM_VERSION}" || "${PLATFORM_VERSION}" == "0.0.0-dev" || "${PLATFORM_VERSION}" == "0.0.0-detect" ]]; then
  die "cannot download release tarball: no valid platform version (version=${PLATFORM_VERSION})"
fi

TARBALL_NAME="globular-${PLATFORM_VERSION}-linux-amd64.tar.gz"
GITHUB_TARBALL_URL="https://github.com/globulario/services/releases/download/v${PLATFORM_VERSION}/${TARBALL_NAME}"
GATEWAY_TARBALL_URL="https://${GATEWAY}/join/bundle/${TARBALL_NAME}"
TARBALL_TMP="/tmp/${TARBALL_NAME}"
TARBALL_DIR="/tmp/globular-${PLATFORM_VERSION}"

log_info "[4.1] Downloading ${TARBALL_NAME}..."
DOWNLOAD_SOURCE="github"
if curl -sfL "${GITHUB_TARBALL_URL}" -o "${TARBALL_TMP}" \
  && curl -sfL "${GITHUB_TARBALL_URL}.sha256" -o "${TARBALL_TMP}.sha256"; then
  EXPECTED_SUM=$(awk '{print $1}' "${TARBALL_TMP}.sha256")
  ACTUAL_SUM=$(sha256sum "${TARBALL_TMP}" | awk '{print $1}')
  if [[ "${EXPECTED_SUM}" != "${ACTUAL_SUM}" ]]; then
    die "tarball checksum mismatch (expected ${EXPECTED_SUM}, got ${ACTUAL_SUM})"
  fi
  rm -f "${TARBALL_TMP}.sha256"
  log_ok "tarball downloaded from GitHub (sha256 ok)"
else
  rm -f "${TARBALL_TMP}" "${TARBALL_TMP}.sha256"
  log_warn "GitHub release tarball unavailable for v${PLATFORM_VERSION} — falling back to controller bundle"
  # Fetch the fallback bundle AND its checksum, then verify — mirroring the
  # GitHub path. The checksum catches corruption/truncation. It is not, by
  # itself, MITM protection over the curl -k transport; pinning the cluster CA
  # here is a tracked follow-up.
  if curl -sfL -k "${GATEWAY_TARBALL_URL}" -o "${TARBALL_TMP}" \
    && curl -sfL -k "${GATEWAY_TARBALL_URL}.sha256" -o "${TARBALL_TMP}.sha256"; then
    EXPECTED_SUM=$(awk '{print $1}' "${TARBALL_TMP}.sha256")
    ACTUAL_SUM=$(sha256sum "${TARBALL_TMP}" | awk '{print $1}')
    if [[ "${EXPECTED_SUM}" != "${ACTUAL_SUM}" ]]; then
      die "controller fallback bundle checksum mismatch (expected ${EXPECTED_SUM}, got ${ACTUAL_SUM})"
    fi
    rm -f "${TARBALL_TMP}.sha256"
    DOWNLOAD_SOURCE="gateway"
    log_ok "tarball downloaded from controller fallback bundle (sha256 ok)"
  else
    rm -f "${TARBALL_TMP}" "${TARBALL_TMP}.sha256"
    die "could not download ${TARBALL_NAME} (+.sha256) from GitHub or controller fallback"
  fi
fi

log_info "[4.2] Extracting tarball..."
mkdir -p "${TARBALL_DIR}"
tar -xzf "${TARBALL_TMP}" -C "${TARBALL_DIR}" --strip-components=1
rm -f "${TARBALL_TMP}"
log_ok "tarball extracted (source=${DOWNLOAD_SOURCE})"

[[ -x "${TARBALL_DIR}/globular-installer" ]] || die "globular-installer not found in tarball"
[[ -d "${TARBALL_DIR}/packages" ]]           || die "packages/ not found in tarball"

cp "${TARBALL_DIR}/globular-installer" "${INSTALL_DIR}/globular-installer"
chmod +x "${INSTALL_DIR}/globular-installer"

log_info "[4.3] Staging packages to ${STATE_DIR}/packages/ ..."
# Move all packages to their canonical location so the node-agent can install
# them locally during convergence without re-fetching from the repository service.
mkdir -p "${STATE_DIR}/packages"
mv "${TARBALL_DIR}/packages/"* "${STATE_DIR}/packages/"
chown -R globular:globular "${STATE_DIR}/packages"
PACKAGES_DIR="${STATE_DIR}/packages"

log_info "[4.4] Copying workflows to state dir..."
cp -r "${TARBALL_DIR}/workflows/." "${STATE_DIR}/workflows/"
chown -R globular:globular "${STATE_DIR}/workflows"
log_ok "globular-installer, packages, and workflows ready"

log_info "[4.5] Installing etcdctl (required by etcd cluster join step)..."
ETCDCTL_PKG=$(find "${PACKAGES_DIR}" -maxdepth 1 -name "etcdctl_*_linux_amd64.tgz" 2>/dev/null | head -1)
if [[ -n "${ETCDCTL_PKG}" ]]; then
  "${INSTALL_DIR}/globular-installer" install \
    --non-interactive \
    "${ETCDCTL_PKG}" \
    || log_warn "etcdctl install failed (non-fatal, will retry via etcd package)"
  ln -sf "${INSTALL_DIR}/etcdctl" /usr/local/bin/etcdctl
  log_ok "etcdctl installed"
else
  log_warn "etcdctl package not found — etcd join may fail if etcdctl is unavailable"
fi

log_info "[4.6] Installing etcd binary and service unit (cluster join in Phase 5)..."
# join_id=${JOIN_ID} — Phase 2.2 validated the JoinPlan. Cluster-affecting steps start in Phase 5.
ETCD_PKG=$(find "${PACKAGES_DIR}" -maxdepth 1 -name "etcd_*_linux_amd64.tgz" 2>/dev/null | head -1)
if [[ -z "${ETCD_PKG}" ]]; then
  die "etcd package not found in ${PACKAGES_DIR} — Phase 4.2 download may have failed"
fi
mkdir -p "${STATE_DIR}/etcd"
chown -R globular:globular "${STATE_DIR}/config" "${STATE_DIR}/etcd"
"${INSTALL_DIR}/globular-installer" install \
  --non-interactive --verbose \
  "${ETCD_PKG}" \
  || die "installer: etcd binary installation failed"
ln -sf "${INSTALL_DIR}/etcdctl" /usr/local/bin/etcdctl
log_ok "etcd binary and service unit installed"

# The installer may have started etcd immediately with a default single-node
# bootstrap config (initial-cluster-state=new). Stop it and wipe any data it
# created so Phase 5 can write the proper join config and start clean.
log_info "[4.6.1] Stopping installer-started etcd and wiping transient bootstrap data..."
systemctl stop globular-etcd.service 2>/dev/null || true
rm -rf "${STATE_DIR}/etcd"
mkdir -p "${STATE_DIR}/etcd"
chown globular:globular "${STATE_DIR}/etcd"
log_ok "installer-started etcd stopped; data dir reset for cluster join"

# ============================================================================
log_phase "5 — Etcd Cluster Join"
# ============================================================================
# All cluster-affecting etcd operations are explicit and auditable here.
# Phase 2.2 (JoinPlan) already authorized this join. join_id=${JOIN_ID}

log_info "[5.1] Checking for ghost etcd member (peer URL from prior failed join)..."
GHOST_ID=$(etcdctl \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
  member list 2>/dev/null \
  | awk -F', ' -v peer="https://${NODE_IP}:2380" '$4==peer {print $1}' | head -1 || true)
if [[ -n "${GHOST_ID}" ]]; then
  log_warn "ghost member detected (peer URL match): ${GHOST_ID} — removing before member add"
  etcdctl \
    --endpoints="${BOOTSTRAP_ETCD}" \
    --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
    member remove "${GHOST_ID}" \
    || log_warn "ghost removal failed — member add may still succeed if peer URL was released"
  log_ok "ghost member removed"
fi

log_info "[5.2] Fetching current cluster member list from ${BOOTSTRAP_ETCD}..."
INITIAL_CLUSTER_RAW=$(etcdctl \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
  member list 2>/dev/null \
  | awk -F', ' 'NF>=4 && $3!="" {print $3"="$4}' | paste -sd, || true)

log_info "[5.3] Normalizing loopback peer URLs (bootstrap may have 127.0.0.1)..."
INITIAL_CLUSTER="${INITIAL_CLUSTER_RAW//127.0.0.1/${BOOTSTRAP_HOST}}"
log_ok "initial-cluster: ${INITIAL_CLUSTER}"

log_info "[5.4] Writing etcd.yaml (initial-cluster-state: existing)..."
mkdir -p "${STATE_DIR}/config"
cat > "${STATE_DIR}/config/etcd.yaml" <<ETCDCFG
name: ${ETCD_NAME}
data-dir: ${STATE_DIR}/etcd

initial-advertise-peer-urls: https://${NODE_IP}:2380
listen-peer-urls: https://0.0.0.0:2380
advertise-client-urls: https://${NODE_IP}:2379
listen-client-urls: https://0.0.0.0:2379

initial-cluster: ${INITIAL_CLUSTER},${ETCD_NAME}=https://${NODE_IP}:2380
initial-cluster-state: "existing"
initial-cluster-token: "${ETCD_CLUSTER_TOKEN}"

client-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}
  trusted-ca-file: ${ETCD_CACERT}
  client-cert-auth: true

peer-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}
  trusted-ca-file: ${ETCD_CACERT}
  client-cert-auth: true
ETCDCFG
chown globular:globular "${STATE_DIR}/config/etcd.yaml"
chmod 640 "${STATE_DIR}/config/etcd.yaml"
log_ok "etcd.yaml written (initial-cluster-state=existing, no loopback peer URLs)"

log_info "[5.5] Adding ${ETCD_NAME} to cluster as a non-voting LEARNER (quorum-safe join)..."
# Learner-first join (Policy A'): the new node joins as a NON-VOTING learner so a
# slow or failed join can never break the existing cluster's quorum. A full-voter
# 1->2 add raises quorum to 2/2 and strands the founder if the joiner dies; a
# learner never changes quorum. The controller promotes caught-up learners to
# voters toward the HA target (>=3 voters) — see cluster_controller
# etcd_learner_promotion.go.
#
# etcd permits only ONE learner at a time (max-learners=1) and serializes raft
# configuration changes. Waiting for that slot and TAKING it must therefore share
# one budget: observing a free slot does not reserve it, and concurrent joiners
# routinely observe the same free slot in the same instant. The loser of that
# race is rejected by the ADD, not by the wait.
#
# Observed 2026-08-01 joining four nodes concurrently: node-3 added instantly;
# node-2 waited 25s for node-3's promotion and added cleanly; node-4 and node-5
# both left the wait, fired "member add" into the contended slot, and died at 30s
# on "etcdserver: request timed out" — having consumed NONE of their 300s budget,
# because the budget covered only the wait while the add was a single attempt
# outside it.
#
# So the add retries until the budget is spent, and errors are classified:
# contention is retryable, a genuine refusal still dies at once, and an exhausted
# budget is fatal. --command-timeout bounds each attempt so one blocked add
# cannot swallow the whole budget in a single call.
# The budget must cover every joiner AHEAD of this one, not just this one's own
# attempt. etcd allows a single learner, so N concurrent joiners serialize: the
# last one waits for N-1 promotions, and each promotion costs a controller
# reconcile cycle plus the new member's raft catch-up.
#
# 300s was too tight. Observed 2026-08-01 with four joiners: promotions landed
# at 18:56:32 and 18:57:02, then the controller stalled and the next did not land
# until 19:01:16 — a four-minute gap. node-4, last in line, waited its full 300s
# and failed with "too many learner members in cluster". The join logic was
# correct; the constant was a guess that did not account for a slow promoter.
#
# 900s is sized from that shape rather than a round number: the observed
# worst-case promotion interval (~4 min) times the queue depth a 5-node bring-up
# produces (3 joiners ahead of the last), with headroom. A joiner that waits this
# long has still cost nothing — it holds no slot while waiting — whereas failing
# early destroys a node that would have joined fine seconds later.
ETCD_ADD_BUDGET_S=900
log_info "[5.5.0] Waiting for the learner slot and adding ${ETCD_NAME} (etcd max-learners=1, up to ${ETCD_ADD_BUDGET_S}s)..."
ETCD_ADD_DEADLINE=$(( $(date +%%s) + ETCD_ADD_BUDGET_S ))
MEMBER_ADD_OUT=""
MEMBER_ADD_OK=0
ETCD_ADD_TRIES=0
while [[ $(date +%%s) -lt ${ETCD_ADD_DEADLINE} ]]; do
  ETCD_ADD_TRIES=$((ETCD_ADD_TRIES + 1))
  # grep -c prints the count and exits 1 on zero matches; under set -euo
  # pipefail that non-zero would abort the whole join at the common 0-learner
  # case, so swallow it with '|| true' (the count "0" is still captured).
  PENDING_LEARNERS=$(etcdctl \
    --endpoints="${BOOTSTRAP_ETCD}" \
    --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
    member list -w json 2>/dev/null | grep -c '"isLearner":true' || true)
  if [[ "${PENDING_LEARNERS:-0}" -ne 0 ]]; then
    if [[ $((ETCD_ADD_TRIES %% 15)) -eq 0 ]]; then
      log_info "still waiting: ${PENDING_LEARNERS} un-promoted learner(s) present..."
    fi
    sleep 2
    continue
  fi

  # The slot LOOKED free. Taking it may still lose to another joiner.
  if MEMBER_ADD_OUT=$(etcdctl \
    --endpoints="${BOOTSTRAP_ETCD}" \
    --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
    --command-timeout=15s \
    member add "${ETCD_NAME}" --learner \
    --peer-urls="https://${NODE_IP}:2380" 2>&1); then
    MEMBER_ADD_OK=1
    break
  fi

  # Contention, not refusal: another joiner holds the slot, or the raft
  # configuration change has not committed yet. Both clear on their own.
  case "${MEMBER_ADD_OUT}" in
    *"too many learner members"*|*"request timed out"*|*"leader changed"*|\
    *"no leader"*|*"unhealthy cluster"*|*"context deadline exceeded"*)
      if [[ $((ETCD_ADD_TRIES %% 15)) -eq 0 ]]; then
        log_info "learner slot contended, retrying: $(echo "${MEMBER_ADD_OUT}" | tail -n1)"
      fi
      sleep 2
      continue
      ;;
  esac

  # Anything else is the bootstrap cluster genuinely refusing this node.
  echo "${MEMBER_ADD_OUT}" >&2
  die "etcd join failed: learner member add rejected by bootstrap cluster"
done

if [[ "${MEMBER_ADD_OK}" -ne 1 ]]; then
  echo "${MEMBER_ADD_OUT}" >&2
  die "etcd join failed: learner member add did not succeed within ${ETCD_ADD_BUDGET_S}s (learner slot never became available)"
fi
echo "${MEMBER_ADD_OUT}"
log_ok "learner member add accepted: ${ETCD_NAME} (non-voting; controller promotes to voter when caught up)"

# [5.5.1] Re-render initial-cluster from the value the member add command just returned.
#
# Phase 5.4 above rendered initial-cluster from a membership snapshot taken
# BEFORE this add. If another node joined inside that window, the snapshot is
# already stale and etcd refuses to start:
#
#   fatal: discovery failed: error validating peerURLs ... member count is unequal
#
# The blast radius is the whole cluster, not this node: etcd allows one learner
# (max-learners=1), so an etcd that never starts never catches up, is never
# promoted, and wedges the learner slot — every subsequent join then dies with
# "too many learner members in cluster". Observed 2026-07-31 joining four nodes
# concurrently: node-2 succeeded, node-3 wedged, node-4 and node-5 both failed
# after waiting the full 300s in phase 5.5.0.
#
# The member add command prints the authoritative ETCD_INITIAL_CLUSTER so the
# caller can write it. Use it rather than the pre-add snapshot.
AUTH_INITIAL_CLUSTER=$(echo "${MEMBER_ADD_OUT}" \
  | sed -n 's/^ETCD_INITIAL_CLUSTER="\(.*\)"$/\1/p' | tail -n1)
if [[ -n "${AUTH_INITIAL_CLUSTER}" ]]; then
  # Same loopback rewrite phase 5.4 applies: a bootstrap founder may advertise
  # 127.0.0.1, which is meaningless from this node.
  AUTH_INITIAL_CLUSTER="${AUTH_INITIAL_CLUSTER//127.0.0.1/${BOOTSTRAP_HOST}}"
  sed -i "s|^initial-cluster: .*|initial-cluster: ${AUTH_INITIAL_CLUSTER}|" \
    "${STATE_DIR}/config/etcd.yaml"
  log_ok "initial-cluster re-rendered from member add: ${AUTH_INITIAL_CLUSTER}"
else
  # Fail closed rather than start etcd against a snapshot we know may be stale.
  die "etcd join failed: could not parse ETCD_INITIAL_CLUSTER from member add output — refusing to start etcd on a possibly stale initial-cluster"
fi

log_info "[5.6] Starting globular-etcd.service..."
systemctl daemon-reload
systemctl enable globular-etcd.service 2>/dev/null || true
systemctl start globular-etcd.service \
  || die "etcd join failed: systemctl start globular-etcd.service returned non-zero"

log_info "[5.6.1] Waiting for local etcd (learner) to come up (max 60s)..."
# A learner REJECTS client read/health RPCs ("etcdserver: rpc not supported for
# learner"), so 'endpoint health' (a linearizable read) cannot be used here.
# 'endpoint status' uses the Status RPC, which learners DO answer — a successful
# status means etcd started and is reachable. Actual cluster membership is proven
# by the anti-fork check below (member list from the bootstrap endpoint).
ETCD_OK=0
for i in $(seq 1 60); do
  if etcdctl \
    --endpoints="https://${NODE_IP}:2379" \
    --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
    endpoint status >/dev/null 2>&1; then
    ETCD_OK=1
    break
  fi
  sleep 1
done
if [[ ${ETCD_OK} -ne 1 ]]; then
  die "etcd did not come up after 60s — etcd join failed; check: journalctl -u globular-etcd"
fi
log_ok "etcd (learner) up on ${NODE_IP}:2379"

log_info "[5.7] Verifying named member presence in cluster (anti-fork check)..."
# Query BOOTSTRAP_ETCD — not local NODE_IP. A forked standalone etcd passes a local check;
# only the bootstrap view proves the node actually joined the existing cluster.
VISIBLE=$(etcdctl \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" --cert="${ETCD_CERT}" --key="${ETCD_KEY}" \
  member list 2>/dev/null | awk -F', ' 'NF>=4 {print $3}' \
  | grep -c "^${ETCD_NAME}$" || echo 0)
if [[ "${VISIBLE}" -eq 0 ]]; then
  die "etcd join failed: ${ETCD_NAME} is not visible in cluster member list from ${BOOTSTRAP_ETCD} — \
this node may have forked its own cluster instead of joining the existing one; \
wipe etcd data and retry with --repair-etcd"
fi
log_ok "member ${ETCD_NAME} confirmed visible in cluster (anti-fork verified)"

# The new membership is proven: restoring pre-join unit state is no longer
# meaningful (the local units now belong to the new cluster state). A failure
# past this point leaves the new units as-is and writes the receipt only.
RESTORE_UNITS_ON_ABORT=0

# ============================================================================
log_phase "6 — ScyllaDB Local Join"
# ============================================================================
# After etcd membership is visible from the bootstrap cluster, install ScyllaDB
# locally from the active BOM package. The scylladb package owns the controlled
# first-start recipe: prevent dpkg autostart, render config from etcd seeds, then
# start scylla-server.service. This keeps Day-1 infra order explicit:
# etcd -> scylla -> local control plane -> node-agent admission/convergence.
#
# Profile gate: joining the ScyllaDB ring changes storage topology (token
# ownership, streaming). It happens ONLY when the controller's JoinPlan
# assigned the 'storage' profile — mirroring the MinIO hold in Phase 9.
# Non-storage nodes reach ScyllaDB over the mesh; convergence owns any later
# storage promotion.

if should_install_role "storage"; then
  log_info "[6.1] Installing ScyllaDB from staged package (join_id=${JOIN_ID})..."
  install_pkg "scylladb"

  wait_unit_active "scylla-server.service" 300 "ScyllaDB" \
    || die "scylla join failed: scylla-server.service did not become active"

  wait_tcp_port "${NODE_IP}" 9042 300 "ScyllaDB CQL" \
    || die "scylla join failed: CQL port 9042 did not become reachable"

  log_info "[6.2] Installing ScyllaDB manager dependencies..."
  install_pkg "scylla-manager"
  install_pkg "scylla-manager-agent"
  log_ok "ScyllaDB local storage stack ready"
else
  log_info "profile gate: 'storage' not in assigned profiles (${ASSIGNED_PROFILES}) — ScyllaDB NOT installed"
  log_ok "ScyllaDB held — controller owns storage placement"
fi

# ============================================================================
log_phase "7 — Local Control Plane"
# ============================================================================
# These packages are installed from the active BOM staged in Phase 4. The
# controller remains the desired-state authority; this phase only ensures the
# joining node has its own local control-plane processes available before the
# node-agent asks the cluster to converge the rest of the profile.
#
# xds/envoy/gateway are the mesh data plane and node HTTP surface — every node
# needs them regardless of profile. workflow/repository/cluster-controller are
# control-plane replicas (leader-elected via etcd) and are installed ONLY when
# the JoinPlan assigned the 'control-plane' profile.

log_info "[7.1] Installing mesh data plane (xds -> envoy -> gateway, join_id=${JOIN_ID})..."
install_pkg "xds"
wait_unit_active "globular-xds.service" 120 "xDS" \
  || die "local control plane failed: globular-xds.service did not become active"

install_pkg "envoy"
wait_unit_active "globular-envoy.service" 120 "Envoy" \
  || die "local control plane failed: globular-envoy.service did not become active"

install_pkg "gateway"
wait_unit_active "globular-gateway.service" 120 "Gateway" \
  || die "local control plane failed: globular-gateway.service did not become active"

if should_install_role "control-plane"; then
  log_info "[7.2] Installing workflow and repository services (join_id=${JOIN_ID})..."
  install_pkg "workflow"
  wait_unit_active "globular-workflow.service" 180 "Workflow" \
    || die "local control plane failed: globular-workflow.service did not become active"

  install_pkg "repository"
  wait_unit_active "globular-repository.service" 180 "Repository" \
    || die "local control plane failed: globular-repository.service did not become active"

  log_info "[7.3] Installing cluster-controller (join_id=${JOIN_ID})..."
  install_pkg "cluster-controller"
  wait_unit_active "globular-cluster-controller.service" 180 "Cluster Controller" \
    || die "local control plane failed: globular-cluster-controller.service did not become active"

  log_ok "Local control-plane stack ready"
else
  log_info "profile gate: 'control-plane' not in assigned profiles (${ASSIGNED_PROFILES}) — workflow/repository/cluster-controller NOT installed"
  log_ok "Mesh data plane ready (control-plane replicas held — controller owns placement)"
fi

# ============================================================================
log_phase "8 — node-agent"
# ============================================================================
# Install the node-agent binary, write a systemd unit that declares an
# ordering dependency on globular-etcd.service, then start it.
# The node-agent reads join_token + controller_endpoint + join_id from the
# state file written in Phase 3.7 and uses the v2 path at startup:
# it polls GetJoinRequestStatus(join_id) without calling RequestJoin again.

log_info "[8.1] Installing globular-node-agent binary via installer engine (join_id=${JOIN_ID})..."
NODE_AGENT_PKG=$(find "${PACKAGES_DIR}" -maxdepth 1 \( -name "node-agent_*_linux_amd64.tgz" -o -name "node_agent_*_linux_amd64.tgz" \) 2>/dev/null | head -1)
if [[ -z "${NODE_AGENT_PKG}" ]]; then
  die "node-agent package not found in ${PACKAGES_DIR} — Phase 4.2 download may have failed"
fi
"${INSTALL_DIR}/globular-installer" install \
  --non-interactive --verbose \
  "${NODE_AGENT_PKG}" \
  || die "installer: node-agent install failed"

log_info "[8.2] Writing node-agent systemd unit (After=globular-etcd.service)..."
NODE_AGENT_BIN=$(find "${INSTALL_DIR}" -maxdepth 1 \( -name "node_agent_server" -o -name "globular-node-agent" \) 2>/dev/null | head -1 || echo "${INSTALL_DIR}/node_agent_server")
cat > "${SYSTEMD_DIR}/globular-node-agent.service" <<NAUNIT
[Unit]
Description=Globular Node Agent
After=network-online.target globular-etcd.service
Requires=network-online.target
Requires=globular-etcd.service

[Service]
User=root
Group=root
ExecStart=${NODE_AGENT_BIN} ${STATE_DIR}
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
NAUNIT
systemctl daemon-reload
systemctl enable globular-node-agent.service 2>/dev/null || true
log_ok "node-agent unit written (Requires=globular-etcd.service)"

log_info "[8.3] Starting globular-node-agent.service..."
systemctl start globular-node-agent.service \
  || log_warn "node-agent start returned non-zero — controller will retry join RPC"

log_info "[8.4] Waiting for node-agent to listen on :11000..."
AGENT_UP=0
for i in $(seq 1 30); do
  if ss -tlnp 2>/dev/null | grep -q ":11000 "; then
    AGENT_UP=1
    break
  fi
  sleep 1
done
if [[ ${AGENT_UP} -ne 1 ]]; then
  log_warn "node-agent not yet listening on :11000 — controller will retry join RPC"
fi
log_ok "node-agent running"

# ============================================================================
log_phase "8.5 — Cleanup installation tarball"
# ============================================================================
log_info "[8.5] Removing extracted tarball (no longer needed)..."
rm -rf "${TARBALL_DIR}"
log_ok "Tarball cleaned up"

# ============================================================================
log_phase "9 — MinIO / Objectstore Hold Enforcement"
# ============================================================================
# The join script MUST NOT start globular-minio.service.
# Objectstore topology is exclusively governed by ObjectStoreDesiredState in
# etcd via 'globular objectstore apply-topology'.
# The controller and node-agent start MinIO ONLY when this node is admitted
# as an objectstore member by an approved TopologyTransition.

log_info "[9.1] Enforcing MinIO hold (objectstore topology contract)..."
if systemctl is-active --quiet globular-minio.service 2>/dev/null; then
  log_warn "globular-minio.service is active — stopping (topology contract violation)"
  systemctl stop globular-minio.service 2>/dev/null || true
fi
if systemctl is-enabled --quiet globular-minio.service 2>/dev/null; then
  systemctl disable globular-minio.service 2>/dev/null || true
  log_info "globular-minio.service disabled (controller drives objectstore topology)"
fi
log_ok "MinIO held — controller owns objectstore topology"

# Join complete — the compensation contract is fulfilled. Disarm the abort
# trap and clear any stale interrupted-join receipt from prior attempts.
trap - EXIT
COMPENSATION_ARMED=0
rm -f "${STATE_DIR}/nodeagent/join_interrupted.json"

# ============================================================================
echo ""
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║          NODE JOINED SUCCESSFULLY                               ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo ""
echo "  Identity:     ${NODE_HOSTNAME} (${NODE_IP})"
echo "  Profiles:     ${ASSIGNED_PROFILES:-<legacy: all roles installed>}"
echo "  etcd:         healthy, member '${ETCD_NAME}'"
if should_install_role "storage"; then
  echo "  scylla:       active locally, CQL reachable"
else
  echo "  scylla:       held (no 'storage' profile — convergence owns storage placement)"
fi
if should_install_role "control-plane"; then
  echo "  controlplane: xds/envoy/gateway/workflow/repository/controller active"
else
  echo "  controlplane: xds/envoy/gateway active (control-plane replicas held)"
fi
echo "  node-agent:   running, auto-initiating join RPC"
echo "  MinIO:        held (not a topology member)"
echo ""
echo "  The node-agent is submitting the join request to the controller."
echo "  Once approved, the controller drives bootstrap convergence:"
echo "    admitted → infra_preparing → etcd_ready →"
echo "    xds_ready → envoy_ready → storage_joining → workload_ready"
echo ""
echo "  Monitor from any cluster node: globular cluster nodes list"
echo ""
`
