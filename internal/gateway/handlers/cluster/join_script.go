package cluster

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
)

// NewJoinScriptHandler serves a self-contained shell script that bootstraps
// a new node to join the cluster. The script includes the CA certificate
// and downloads the node-agent + globularcli binaries from the gateway.
//
// Usage from a new node:
//
//	curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <join-token>
func NewJoinScriptHandler(controllerAddr string, gatewayPort int) http.Handler {
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
		// Use the Host header which contains the actual address:port the
		// client used (e.g. "10.0.0.63:443"). This ensures the join script
		// downloads binaries from the same endpoint that served the script.
		gatewayHost := r.Host
		if gatewayHost == "" {
			gatewayHost = r.URL.Host
		}

		// Extract host and port from the request.
		reqHost, reqPort, err := net.SplitHostPort(gatewayHost)
		if err != nil {
			// No port in Host header — use the host as-is with our fallback port.
			reqHost = gatewayHost
			reqPort = fmt.Sprintf("%d", gatewayPort)
		}
		if reqHost == "" {
			reqHost = "127.0.0.1"
		}

		gatewayAddr := net.JoinHostPort(reqHost, reqPort)

		// Controller address for the join command.
		// Use the request host IP (not the full host:port) with controller port.
		ctrlAddr := controllerAddr
		if ctrlAddr == "" || ctrlAddr == "127.0.0.1:12000" || ctrlAddr == "localhost:12000" {
			ctrlAddr = fmt.Sprintf("%s:12000", reqHost)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, joinScript, caB64, gatewayAddr, ctrlAddr)
	})
}

// joinScript is the shell script template served to joining nodes.
// Args: %[1]s = base64 CA cert, %[2]s = gateway address, %[3]s = controller address
const joinScript = `#!/bin/bash
set -eu

# ── Globular Node Join Script ──────────────────────────────────────────────
# This script bootstraps a new node to join an existing Globular cluster.
# It installs the minimum required components (node-agent + CLI), sets up
# TLS trust, and runs the cluster join command.
#
# Usage:
#   curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <join-token>
# ───────────────────────────────────────────────────────────────────────────

CA_B64="%[1]s"
GATEWAY="%[2]s"
CONTROLLER="%[3]s"
JOIN_TOKEN=""
INSTALL_DIR="/usr/lib/globular/bin"
STATE_DIR="/var/lib/globular"
SYSTEMD_DIR="/etc/systemd/system"

# ── Parse arguments ───────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    --token) JOIN_TOKEN="$2"; shift 2 ;;
    --controller) CONTROLLER="$2"; shift 2 ;;
    *) echo "Unknown argument: $1"; exit 1 ;;
  esac
done

if [[ -z "${JOIN_TOKEN}" ]]; then
  echo "Error: --token is required"
  echo "Usage: curl -sfL https://<gateway>:8443/join -k | sudo bash -s -- --token <join-token>"
  exit 1
fi

BOOTSTRAP_HOST=$(echo "${GATEWAY}" | cut -d: -f1)

echo "━━━ Globular Node Join ━━━"
echo "  Controller: ${CONTROLLER}"
echo "  Gateway:    ${GATEWAY}"
echo "  Bootstrap:  ${BOOTSTRAP_HOST}"
echo ""

# ── Check root ────────────────────────────────────────────────────────────
if [[ "$(id -u)" -ne 0 ]]; then
  echo "Error: this script must be run as root (sudo)"
  exit 1
fi

# ── Stop existing services and clean state ────────────────────────────────
# Critical: old node-agent state causes auto-registration with stale node ID
# instead of going through the proper join flow.
echo "  → Cleaning previous installation state..."
systemctl stop 'globular-*.service' 2>/dev/null || true
systemctl disable 'globular-*.service' 2>/dev/null || true

# Remove any ghost etcd member for this node from the cluster BEFORE wiping
# local state. Without this, a failed join leaves a registered member that
# the bootstrap node keeps trying to reach, causing quorum loss (needs 2/2).
NODE_IP_CLEAN=$(hostname -I | awk '{print $1}')
if [[ -n "${NODE_IP_CLEAN}" ]]; then
  GHOST_ID=$(ETCDCTL_API=3 etcdctl \
    --endpoints="https://${BOOTSTRAP_HOST}:2379" \
    --cacert="${STATE_DIR}/pki/ca.crt" \
    member list -w simple 2>/dev/null \
    | grep "${NODE_IP_CLEAN}" | cut -d',' -f1 || true)
  if [[ -n "${GHOST_ID}" ]]; then
    echo "  → Removing stale etcd member ${GHOST_ID} for ${NODE_IP_CLEAN}..."
    ETCDCTL_API=3 etcdctl \
      --endpoints="https://${BOOTSTRAP_HOST}:2379" \
      --cacert="${STATE_DIR}/pki/ca.crt" \
      member remove "${GHOST_ID}" 2>/dev/null || true
    echo "  ✓ Stale etcd member removed"
  fi
fi

rm -rf "${STATE_DIR}/node_agent" "${STATE_DIR}/nodeagent" "${STATE_DIR}/etcd" "${STATE_DIR}/config"
rm -f /etc/systemd/system/globular-*.service
systemctl daemon-reload 2>/dev/null || true

# ── Create user and directories ───────────────────────────────────────────
echo "  → Creating globular user and directories..."
if ! id -u globular &>/dev/null; then
  useradd --system --no-create-home --home-dir "${STATE_DIR}" --shell /usr/sbin/nologin globular
fi
mkdir -p "${INSTALL_DIR}" "${STATE_DIR}/pki/issued/services" "${STATE_DIR}/services" "${STATE_DIR}/node_agent"

# ── Install CA certificate ────────────────────────────────────────────────
echo "  → Installing cluster CA certificate..."
echo "${CA_B64}" | base64 -d > "${STATE_DIR}/pki/ca.crt"
cp "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
chmod 644 "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
echo "  ✓ CA certificate installed"

# ── Generate TLS certificate ─────────────────────────────────────────────
echo "  → Generating service certificate..."
PKI_DIR="${STATE_DIR}/pki/issued/services"
NODE_HOSTNAME=$(hostname)
NODE_IP=$(hostname -I | awk '{print $1}')

openssl genrsa -out "${PKI_DIR}/service.key" 2048 2>/dev/null

cat > /tmp/globular-san.cnf <<SANEOF
[req]
req_extensions = v3_req
distinguished_name = req_dn
prompt = no
[req_dn]
CN = ${NODE_HOSTNAME}
O = globular.internal
[v3_req]
subjectAltName = DNS:${NODE_HOSTNAME},DNS:localhost,IP:${NODE_IP},IP:127.0.0.1
SANEOF

openssl req -new -key "${PKI_DIR}/service.key" \
  -config /tmp/globular-san.cnf \
  -out /tmp/globular-join.csr 2>/dev/null

CSR_B64=$(base64 -w0 /tmp/globular-join.csr)
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/sign_ca_certificate?csr=${CSR_B64}" \
  -o "${PKI_DIR}/service.crt"
rm -f /tmp/globular-join.csr /tmp/globular-san.cnf

if [[ ! -s "${PKI_DIR}/service.crt" ]]; then
  echo "  ✗ Certificate signing failed"
  exit 1
fi

chmod 600 "${PKI_DIR}/service.key"
chmod 644 "${PKI_DIR}/service.crt"
echo "  ✓ Service certificate issued for ${NODE_HOSTNAME} (${NODE_IP})"

# ── Download binaries ─────────────────────────────────────────────────────
echo "  → Downloading node-agent..."
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/join/bin/node_agent_server" \
  -o "${INSTALL_DIR}/node_agent_server"
chmod +x "${INSTALL_DIR}/node_agent_server"
echo "  ✓ node_agent_server installed"

echo "  → Downloading globular CLI..."
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/join/bin/globularcli" \
  -o "${INSTALL_DIR}/globularcli"
chmod +x "${INSTALL_DIR}/globularcli"
ln -sf "${INSTALL_DIR}/globularcli" /usr/local/bin/globular
echo "  ✓ globularcli installed"

echo "  → Downloading etcd..."
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/join/bin/etcd" \
  -o "${INSTALL_DIR}/etcd"
chmod +x "${INSTALL_DIR}/etcd"
echo "  ✓ etcd installed"

echo "  → Downloading etcdctl..."
curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
  "https://${GATEWAY}/join/bin/etcdctl" \
  -o "${INSTALL_DIR}/etcdctl"
chmod +x "${INSTALL_DIR}/etcdctl"
ln -sf "${INSTALL_DIR}/etcdctl" /usr/local/bin/etcdctl
echo "  ✓ etcdctl installed"

# ── Set ownership ─────────────────────────────────────────────────────────
chown -R globular:globular "${STATE_DIR}"

# ── Join etcd cluster ─────────────────────────────────────────────────────
# etcd MemberAdd + start MUST happen in the join script (not controller)
# because MemberAdd on a single-node cluster immediately requires 2/2 for
# quorum. If the new node's etcd doesn't start fast enough, the entire
# cluster loses quorum and the controller can't even rollback.
#
# The script calls MemberAdd, writes config, and starts etcd IMMEDIATELY.

ETCD_CACERT="${STATE_DIR}/pki/ca.crt"
ETCD_CERT="${PKI_DIR}/service.crt"
ETCD_KEY="${PKI_DIR}/service.key"
BOOTSTRAP_ETCD="https://${BOOTSTRAP_HOST}:2379"
ETCD_NAME=$(echo "${NODE_HOSTNAME}" | sed 's/[^a-zA-Z0-9_-]/-/g; s/^-//; s/-$//')
if [[ -z "${ETCD_NAME}" ]]; then ETCD_NAME="node"; fi

echo "  → Joining etcd cluster..."
mkdir -p "${STATE_DIR}/config" "${STATE_DIR}/etcd"
chown -R globular:globular "${STATE_DIR}/config" "${STATE_DIR}/etcd"

# Get current members to build initial-cluster string.
MEMBER_LIST=$(ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" \
  member list -w simple 2>/dev/null || true)

INITIAL_CLUSTER=""
while IFS=',' read -r id status name peerURLs clientURLs isLearner; do
  name=$(echo "$name" | xargs)
  peerURLs=$(echo "$peerURLs" | xargs)
  if [[ -n "$name" && -n "$peerURLs" ]]; then
    if [[ -n "$INITIAL_CLUSTER" ]]; then
      INITIAL_CLUSTER="${INITIAL_CLUSTER},"
    fi
    INITIAL_CLUSTER="${INITIAL_CLUSTER}${name}=${peerURLs}"
  fi
done <<< "${MEMBER_LIST}"

# Call MemberAdd.
ADD_OUTPUT=$(ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" \
  member add "${ETCD_NAME}" --peer-urls="https://${NODE_IP}:2380" 2>&1) || {
  if echo "${ADD_OUTPUT}" | grep -qi "already exists\|Peer URLs already exists"; then
    echo "  ⚠ Node already registered as etcd member (re-joining)"
  else
    echo "  ✗ Failed to add etcd member: ${ADD_OUTPUT}"
    exit 1
  fi
}

# Append ourselves to initial-cluster.
if [[ -n "$INITIAL_CLUSTER" ]]; then
  INITIAL_CLUSTER="${INITIAL_CLUSTER},${ETCD_NAME}=https://${NODE_IP}:2380"
else
  INITIAL_CLUSTER="${ETCD_NAME}=https://${NODE_IP}:2380"
fi
echo "  ✓ Registered as etcd member: ${ETCD_NAME}"

# Write etcd config — IMMEDIATELY, before quorum times out.
cat > "${STATE_DIR}/config/etcd.yaml" <<ETCDCFG
name: "${ETCD_NAME}"
data-dir: "${STATE_DIR}/etcd"
listen-client-urls: "https://${NODE_IP}:2379,https://127.0.0.1:2379"
advertise-client-urls: "https://${NODE_IP}:2379"
listen-peer-urls: "https://${NODE_IP}:2380"
initial-advertise-peer-urls: "https://${NODE_IP}:2380"
initial-cluster: "${INITIAL_CLUSTER}"
initial-cluster-state: "existing"
initial-cluster-token: "globular-etcd-cluster"

client-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}

peer-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}
  trusted-ca-file: ${ETCD_CACERT}
ETCDCFG
chown globular:globular "${STATE_DIR}/config/etcd.yaml"

# Write etcd endpoints (bootstrap + self).
echo "https://${BOOTSTRAP_HOST}:2379" > "${STATE_DIR}/config/etcd_endpoints"
echo "https://${NODE_IP}:2379" >> "${STATE_DIR}/config/etcd_endpoints"
chown globular:globular "${STATE_DIR}/config/etcd_endpoints"

# Create and START etcd systemd unit IMMEDIATELY after MemberAdd.
# This is time-critical: etcd quorum requires both members running.
cat > "${SYSTEMD_DIR}/globular-etcd.service" <<ETCDUNIT
[Unit]
Description=Globular etcd
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=globular
Group=globular
ExecStartPre=/usr/bin/mkdir -p ${STATE_DIR}/etcd
ExecStartPre=/usr/bin/chown globular:globular ${STATE_DIR}/etcd
ExecStartPre=/usr/bin/chmod 0750 ${STATE_DIR}/etcd
ExecStart=${INSTALL_DIR}/etcd --config-file ${STATE_DIR}/config/etcd.yaml
Restart=on-failure
RestartSec=5
LimitNOFILE=524288

[Install]
WantedBy=multi-user.target
ETCDUNIT

systemctl daemon-reload
systemctl enable globular-etcd.service
systemctl start globular-etcd.service

# Wait for etcd to be healthy (both members must form quorum).
echo "  → Waiting for etcd to join cluster..."
ETCD_OK=0
for i in $(seq 1 30); do
  if ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
    --endpoints="https://127.0.0.1:2379" \
    --cacert="${ETCD_CACERT}" \
    endpoint health 2>/dev/null | grep -q "is healthy"; then
    ETCD_OK=1
    break
  fi
  sleep 2
done

if [[ $ETCD_OK -eq 1 ]]; then
  echo "  ✓ etcd running and healthy (member of cluster)"
else
  echo "  ✗ etcd failed to join — check: journalctl -u globular-etcd.service"
  echo "  ⚠ Continuing with node-agent setup (controller may recover)"
fi

# ── Create node-agent systemd unit ────────────────────────────────────────
echo "  → Creating node-agent systemd unit..."
cat > "${SYSTEMD_DIR}/globular-node-agent.service" <<UNIT
[Unit]
Description=Globular node_agent
After=network-online.target globular-etcd.service
Wants=network-online.target
Requires=globular-etcd.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=${STATE_DIR}/node_agent
Environment=GLOBULAR_SERVICES_DIR=${STATE_DIR}/services
ExecStartPre=/bin/sh -c 'mkdir -p ${STATE_DIR}/node_agent'
ExecStartPre=/bin/sh -c 'for i in \$(seq 1 60); do [ -f ${PKI_DIR}/service.crt ] && exit 0; sleep 1; done; echo "TLS cert not ready"; exit 1'
ExecStart=${INSTALL_DIR}/node_agent_server
Restart=on-failure
RestartSec=2
LimitNOFILE=524288

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable globular-node-agent.service
systemctl start globular-node-agent.service
echo "  ✓ node-agent started"

# Wait for node-agent to be ready (check gRPC port).
echo "  → Waiting for node-agent..."
for i in $(seq 1 30); do
  if ss -tlnp 2>/dev/null | grep -q ":11000 "; then
    break
  fi
  sleep 1
done
sleep 2
echo "  ✓ node-agent ready"

# ── Join the cluster ──────────────────────────────────────────────────────
echo "  → Joining cluster..."
globular cluster join \
  --join-token "${JOIN_TOKEN}" \
  --controller "${CONTROLLER}" \
  --node "localhost:11000" \
  --insecure

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          ✓ NODE JOINED SUCCESSFULLY                          ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "  etcd:       running (member of cluster)"
echo "  node-agent: running"
echo ""
echo "  The controller will now drive this node through bootstrap phases:"
echo "    admitted → infra_preparing → etcd_joining → etcd_ready →"
echo "    xds_ready → envoy_ready → storage_joining → workload_ready"
echo ""
echo "  Monitor progress: globular cluster nodes list"
echo ""
`
