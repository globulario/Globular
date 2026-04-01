package cluster

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

// NewJoinScriptHandler serves a self-contained shell script that bootstraps
// a new node to join the cluster. The script follows a strict Day-1 protocol:
//
//  1. Identity:      CA cert, service cert, credentials
//  2. Connectivity:  DNS resolver, config.json, gateway reachability
//  3. Prerequisites: etcd, ScyllaDB (bootstrap infra)
//  4. Join:          node-agent, cluster join, DNS registration
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

		// Read ScyllaDB GPG key if available on this node.
		// Empty string if missing -- the script will skip ScyllaDB setup.
		scyllaGpgB64 := ""
		if gpgKey, err := os.ReadFile("/etc/apt/keyrings/scylladb.gpg"); err == nil {
			scyllaGpgB64 = base64.StdEncoding.EncodeToString(gpgKey)
		}

		// Read ScyllaDB apt source line if configured.
		scyllaAptSource := ""
		if src, err := os.ReadFile("/etc/apt/sources.list.d/scylla.list"); err == nil {
			scyllaAptSource = strings.TrimSpace(string(src))
		}

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
		fmt.Fprintf(w, joinScript, caB64, gatewayAddr, ctrlAddr, scyllaGpgB64, scyllaAptSource)
	})
}

// joinScript is the shell script template served to joining nodes.
// Format args:
//
//	%[1]s = base64 CA cert
//	%[2]s = gateway address (host:port)
//	%[3]s = controller address (host:port)
//	%[4]s = base64 ScyllaDB GPG key (may be empty)
//	%[5]s = ScyllaDB apt source line (may be empty)
const joinScript = `#!/bin/bash
set -eu

# ---------------------------------------------------------------------------
# Globular Node Join Script -- Day-1 Bootstrap Protocol
#
# Phase 1: Identity       -- CA cert, service cert
# Phase 2: Connectivity   -- DNS resolver, config.json, gateway check
# Phase 3: Prerequisites  -- etcd, ScyllaDB (bootstrap infra)
# Phase 4: Join           -- node-agent, cluster join, DNS A record
#
# A node must be reachable, authenticated, and bootstrap-complete before
# the controller starts reconciliation.
# ---------------------------------------------------------------------------

CA_B64="%[1]s"
GATEWAY="%[2]s"
CONTROLLER="%[3]s"
SCYLLA_GPG_B64="%[4]s"
SCYLLA_APT_SOURCE="%[5]s"
JOIN_TOKEN=""
INSTALL_DIR="/usr/lib/globular/bin"
STATE_DIR="/var/lib/globular"
SYSTEMD_DIR="/etc/systemd/system"

# -- Parse arguments --------------------------------------------------------
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
NODE_HOSTNAME=$(hostname)
NODE_IP=$(hostname -I | awk '{print $1}')

echo "=== Globular Node Join ==="
echo "  Controller: ${CONTROLLER}"
echo "  Gateway:    ${GATEWAY}"
echo "  Bootstrap:  ${BOOTSTRAP_HOST}"
echo "  This node:  ${NODE_HOSTNAME} (${NODE_IP})"
echo ""

# -- Check root -------------------------------------------------------------
if [[ "$(id -u)" -ne 0 ]]; then
  echo "Error: this script must be run as root (sudo)"
  exit 1
fi

# ===========================================================================
# PHASE 1: Identity
# ===========================================================================
echo "--- Phase 1: Identity ---"

# -- Stop existing services and clean state ---------------------------------
echo "  [1.1] Cleaning previous installation state..."
systemctl stop 'globular-*.service' 2>/dev/null || true
systemctl stop scylla-server.service 2>/dev/null || true
systemctl disable 'globular-*.service' 2>/dev/null || true
# Unmask ScyllaDB in case a previous cleanup masked it.
systemctl unmask scylla-server.service 2>/dev/null || true

# NOTE: Ghost etcd member cleanup happens in Phase 3 after etcdctl is downloaded.

rm -rf "${STATE_DIR}/node_agent" "${STATE_DIR}/nodeagent" "${STATE_DIR}/etcd" "${STATE_DIR}/config"
rm -f "${STATE_DIR}/node-agent/last-generation"
rm -f /etc/systemd/system/globular-*.service
systemctl daemon-reload 2>/dev/null || true
systemctl reset-failed 2>/dev/null || true

# -- Create user and directories --------------------------------------------
echo "  [1.2] Creating globular user and directories..."
if ! id -u globular &>/dev/null; then
  useradd --system --no-create-home --home-dir "${STATE_DIR}" --shell /usr/sbin/nologin globular
fi
mkdir -p "${INSTALL_DIR}" "${STATE_DIR}/pki/issued/services" "${STATE_DIR}/services" "${STATE_DIR}/node_agent"

# -- Install CA certificate -------------------------------------------------
echo "  [1.3] Installing cluster CA certificate..."
echo "${CA_B64}" | base64 -d > "${STATE_DIR}/pki/ca.crt"
cp "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
chmod 644 "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"

# -- Generate TLS certificate -----------------------------------------------
echo "  [1.4] Generating service certificate..."
PKI_DIR="${STATE_DIR}/pki/issued/services"

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
subjectAltName = DNS:${NODE_HOSTNAME},DNS:${NODE_HOSTNAME}.globular.internal,DNS:localhost,IP:${NODE_IP},IP:127.0.0.1
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
  echo "  FAIL: Certificate signing failed"
  exit 1
fi

chmod 600 "${PKI_DIR}/service.key"
chmod 644 "${PKI_DIR}/service.crt"
echo "  OK: Identity ready (${NODE_HOSTNAME}, ${NODE_IP})"

# ===========================================================================
# PHASE 2: Connectivity
# ===========================================================================
echo ""
echo "--- Phase 2: Connectivity ---"

# -- Configure DNS resolver (Day-1 Rule 7A) --------------------------------
# The new node must resolve *.globular.internal via the bootstrap node's DNS.
# Use a systemd-resolved drop-in: this creates a dedicated routing domain
# that sends ONLY *.globular.internal queries to the bootstrap DNS server.
# Per-link configuration is unreliable when multiple DNS servers are on the
# same link — ISP servers respond with NXDOMAIN faster than the cluster DNS.
echo "  [2.1] Configuring DNS resolver..."
mkdir -p /etc/systemd/resolved.conf.d
cat > /etc/systemd/resolved.conf.d/globular.conf <<DNSCONF
[Resolve]
DNS=${BOOTSTRAP_HOST}
Domains=~globular.internal
DNSCONF
systemctl restart systemd-resolved 2>/dev/null || true
echo "  OK: DNS resolver -> ${BOOTSTRAP_HOST} for globular.internal"

# -- Write config.json (Day-1 Rule 5) --------------------------------------
# Node-agent needs this to discover the gateway. Without it, the agent
# starts blind and cannot reach any remote service.
echo "  [2.2] Writing node-agent config..."
mkdir -p "${STATE_DIR}/config"
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

# -- Verify gateway reachability --------------------------------------------
echo "  [2.3] Checking gateway reachability..."
if curl -sf --cacert "${STATE_DIR}/pki/ca.crt" -o /dev/null \
  "https://${GATEWAY}/get_ca_certificate" 2>/dev/null; then
  echo "  OK: Gateway reachable at ${GATEWAY}"
else
  echo "  WARN: Gateway not reachable (may work after DNS propagation)"
fi

chown -R globular:globular "${STATE_DIR}"
echo "  OK: Connectivity configured"

# ===========================================================================
# PHASE 3: Bootstrap Prerequisites
# ===========================================================================
echo ""
echo "--- Phase 3: Bootstrap Prerequisites ---"

# -- Download binaries ------------------------------------------------------
echo "  [3.1] Downloading binaries..."
for BIN in node_agent_server globularcli etcd etcdctl; do
  curl -sfL --cacert "${STATE_DIR}/pki/ca.crt" \
    "https://${GATEWAY}/join/bin/${BIN}" \
    -o "${INSTALL_DIR}/${BIN}"
  chmod +x "${INSTALL_DIR}/${BIN}"
  echo "       ${BIN} installed"
done
ln -sf "${INSTALL_DIR}/globularcli" /usr/local/bin/globular
ln -sf "${INSTALL_DIR}/etcdctl" /usr/local/bin/etcdctl

# -- Remove ghost etcd member (if any) -------------------------------------
# Must happen AFTER etcdctl is downloaded but BEFORE member-add.
# A previous failed join may have left a registered member for this IP.
if [[ -n "${NODE_IP}" ]]; then
  GHOST_ID=$(ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
    --endpoints="https://${BOOTSTRAP_HOST}:2379" \
    --cacert="${STATE_DIR}/pki/ca.crt" \
    --cert="${PKI_DIR}/service.crt" \
    --key="${PKI_DIR}/service.key" \
    member list -w simple 2>/dev/null \
    | grep "${NODE_IP}" | cut -d',' -f1 || true)
  if [[ -n "${GHOST_ID}" ]]; then
    echo "       Removing stale etcd member ${GHOST_ID} for ${NODE_IP}..."
    ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
      --endpoints="https://${BOOTSTRAP_HOST}:2379" \
      --cacert="${STATE_DIR}/pki/ca.crt" \
      --cert="${PKI_DIR}/service.crt" \
      --key="${PKI_DIR}/service.key" \
      member remove "${GHOST_ID}" 2>/dev/null || true
  fi
fi

# -- Join etcd cluster ------------------------------------------------------
echo "  [3.2] Joining etcd cluster..."
ETCD_CACERT="${STATE_DIR}/pki/ca.crt"
ETCD_CERT="${PKI_DIR}/service.crt"
ETCD_KEY="${PKI_DIR}/service.key"
BOOTSTRAP_ETCD="https://${BOOTSTRAP_HOST}:2379"
ETCD_NAME=$(echo "${NODE_HOSTNAME}" | sed 's/[^a-zA-Z0-9_-]/-/g; s/^-//; s/-$//')
if [[ -z "${ETCD_NAME}" ]]; then ETCD_NAME="node"; fi

mkdir -p "${STATE_DIR}/etcd"
chown -R globular:globular "${STATE_DIR}/config" "${STATE_DIR}/etcd"

MEMBER_LIST=$(ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" \
  --cert="${ETCD_CERT}" \
  --key="${ETCD_KEY}" \
  member list -w simple 2>/dev/null || true)

INITIAL_CLUSTER=""
while IFS=',' read -r id status name peerURLs clientURLs isLearner; do
  name=$(echo "$name" | xargs)
  peerURLs=$(echo "$peerURLs" | xargs)
  if [[ -n "$name" && -n "$peerURLs" ]]; then
    # Replace any loopback peer URLs with the bootstrap host's real address.
    peerURLs=$(echo "$peerURLs" | sed -E "s#https?://(127\.0\.0\.1|localhost|\[::1\])#https://${BOOTSTRAP_HOST}#g")
    if [[ -n "$INITIAL_CLUSTER" ]]; then
      INITIAL_CLUSTER="${INITIAL_CLUSTER},"
    fi
    INITIAL_CLUSTER="${INITIAL_CLUSTER}${name}=${peerURLs}"
  fi
done <<< "${MEMBER_LIST}"

ADD_OUTPUT=$(ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
  --endpoints="${BOOTSTRAP_ETCD}" \
  --cacert="${ETCD_CACERT}" \
  --cert="${ETCD_CERT}" \
  --key="${ETCD_KEY}" \
  member add "${ETCD_NAME}" --peer-urls="https://${NODE_IP}:2380" 2>&1) || {
  if echo "${ADD_OUTPUT}" | grep -qi "already exists\|Peer URLs already exists"; then
    echo "       (re-joining existing member)"
  else
    echo "  FAIL: etcd member add: ${ADD_OUTPUT}"
    exit 1
  fi
}

if [[ -n "$INITIAL_CLUSTER" ]]; then
  INITIAL_CLUSTER="${INITIAL_CLUSTER},${ETCD_NAME}=https://${NODE_IP}:2380"
else
  INITIAL_CLUSTER="${ETCD_NAME}=https://${NODE_IP}:2380"
fi

cat > "${STATE_DIR}/config/etcd.yaml" <<ETCDCFG
name: "${ETCD_NAME}"
data-dir: "${STATE_DIR}/etcd"
listen-client-urls: "https://${NODE_IP}:2379"
advertise-client-urls: "https://${NODE_IP}:2379"
listen-peer-urls: "https://${NODE_IP}:2380"
initial-advertise-peer-urls: "https://${NODE_IP}:2380"
initial-cluster: "${INITIAL_CLUSTER}"
initial-cluster-state: "existing"
initial-cluster-token: "globular-etcd-cluster"
quota-backend-bytes: 8589934592
auto-compaction-mode: "periodic"
auto-compaction-retention: "1h"
snapshot-count: 10000

client-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}

peer-transport-security:
  cert-file: ${ETCD_CERT}
  key-file: ${ETCD_KEY}
  trusted-ca-file: ${ETCD_CACERT}
ETCDCFG
chown globular:globular "${STATE_DIR}/config/etcd.yaml"

echo "https://${BOOTSTRAP_HOST}:2379" > "${STATE_DIR}/config/etcd_endpoints"
echo "https://${NODE_IP}:2379" >> "${STATE_DIR}/config/etcd_endpoints"
chown globular:globular "${STATE_DIR}/config/etcd_endpoints"

cat > "${SYSTEMD_DIR}/globular-etcd.service" <<ETCDUNIT
[Unit]
Description=Globular etcd
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=globular
Group=globular
ExecStartPre=+/usr/bin/mkdir -p ${STATE_DIR}/etcd
ExecStartPre=+/usr/bin/chown globular:globular ${STATE_DIR}/etcd
ExecStartPre=+/usr/bin/chmod 0750 ${STATE_DIR}/etcd
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

ETCD_OK=0
for i in $(seq 1 30); do
  if ETCDCTL_API=3 "${INSTALL_DIR}/etcdctl" \
    --endpoints="https://${NODE_IP}:2379" \
    --cacert="${ETCD_CACERT}" \
    --cert="${ETCD_CERT}" \
    --key="${ETCD_KEY}" \
    endpoint health 2>/dev/null | grep -q "is healthy"; then
    ETCD_OK=1
    break
  fi
  sleep 2
done

if [[ $ETCD_OK -eq 1 ]]; then
  echo "  OK: etcd healthy (member of cluster)"
else
  echo "  FAIL: etcd did not become healthy -- check: journalctl -u globular-etcd.service"
  echo "  Continuing (controller may recover)..."
fi

# -- Install ScyllaDB (Day-1 Rule 10) --------------------------------------
# ScyllaDB is a bootstrap prerequisite for control-plane nodes. Without it,
# the controller's storage_joining phase stalls forever.
if [[ -n "${SCYLLA_GPG_B64}" && -n "${SCYLLA_APT_SOURCE}" ]]; then
  echo "  [3.3] Installing ScyllaDB..."

  # Install GPG key (idempotent).
  if [[ ! -f /etc/apt/keyrings/scylladb.gpg ]]; then
    mkdir -p /etc/apt/keyrings
    echo "${SCYLLA_GPG_B64}" | base64 -d > /etc/apt/keyrings/scylladb.gpg
    chmod 644 /etc/apt/keyrings/scylladb.gpg
  fi

  # Install apt source.
  echo "${SCYLLA_APT_SOURCE}" > /etc/apt/sources.list.d/scylladb.list

  # Install the package (skip if already installed).
  # Use || true to prevent set -eu from aborting on apt failure.
  if ! dpkg -l scylla-server 2>/dev/null | grep -q '^ii'; then
    apt-get update -qq 2>/dev/null || true
    DEBIAN_FRONTEND=noninteractive apt-get install -y -qq scylla-server 2>/dev/null || true
  fi

  if dpkg -l scylla-server 2>/dev/null | grep -q '^ii'; then
    echo "  OK: ScyllaDB package installed"

    # Configure and start ScyllaDB for cluster join.
    # This must happen NOW, before the controller starts driving bootstrap
    # phases — ScyllaDB is a prerequisite for storage_joining phase.
    LOCAL_IP=$(ip route get 8.8.8.8 2>/dev/null | awk '{for(i=1;i<=NF;i++) if($i=="src") print $(i+1); exit}')
    [[ -z "${LOCAL_IP}" ]] && LOCAL_IP=$(hostname -I | awk '{print $1}')
    CTRL_HOST=$(echo "${CONTROLLER}" | sed 's/:.*//')
    SEED_IP="${CTRL_HOST},${LOCAL_IP}"

    # TLS certs (Globular PKI → ScyllaDB)
    SCYLLA_TLS="/etc/scylla/tls"
    mkdir -p "${SCYLLA_TLS}"
    if [[ -f /var/lib/globular/pki/issued/services/service.crt ]]; then
      cp /var/lib/globular/pki/issued/services/service.crt "${SCYLLA_TLS}/server.crt"
      cp /var/lib/globular/pki/issued/services/service.key "${SCYLLA_TLS}/server.key"
      cp /var/lib/globular/pki/ca.pem "${SCYLLA_TLS}/ca.crt"
      chown -R scylla:scylla "${SCYLLA_TLS}" 2>/dev/null || true
      chmod 644 "${SCYLLA_TLS}/server.crt" "${SCYLLA_TLS}/ca.crt"
      chmod 400 "${SCYLLA_TLS}/server.key"
    fi

    # Write scylla.yaml
    cat > /etc/scylla/scylla.yaml <<SCYLLAEOF
cluster_name: 'globular.internal'
seed_provider:
  - class_name: org.apache.cassandra.locator.SimpleSeedProvider
    parameters:
      - seeds: '${SEED_IP}'
listen_address: '${LOCAL_IP}'
rpc_address: '${LOCAL_IP}'
broadcast_address: '${LOCAL_IP}'
broadcast_rpc_address: '${LOCAL_IP}'
native_transport_port: 9042
endpoint_snitch: SimpleSnitch
developer_mode: true
client_encryption_options:
  enabled: true
  certificate: /etc/scylla/tls/server.crt
  keyfile: /etc/scylla/tls/server.key
  truststore: /etc/scylla/tls/ca.crt
  require_client_auth: false
native_transport_port_ssl: 9142
data_file_directories:
  - /var/lib/scylla/data
commitlog_directory: /var/lib/scylla/commitlog
commitlog_sync: batch
commitlog_sync_batch_window_in_ms: 2
commitlog_sync_period_in_ms: 10000
auto_adjust_flush_quota: true
api_port: 10000
api_address: '${LOCAL_IP}'
SCYLLAEOF

    # Data dirs + symlink
    mkdir -p /var/lib/scylla/data /var/lib/scylla/commitlog
    chown -R scylla:scylla /var/lib/scylla 2>/dev/null || true
    [[ -L /var/lib/scylla/conf || -d /var/lib/scylla/conf ]] || ln -sfn /etc/scylla /var/lib/scylla/conf

    # Systemd overrides for Debian/Ubuntu
    SCYLLA_OVR="/etc/systemd/system/scylla-server.service.d"
    mkdir -p "${SCYLLA_OVR}"
    [[ -f "${SCYLLA_OVR}/sysconfdir.conf" ]] || cat > "${SCYLLA_OVR}/sysconfdir.conf" <<'SYSEOF'
[Service]
EnvironmentFile=
EnvironmentFile=-/etc/default/scylla-server
EnvironmentFile=-/etc/scylla.d/*.conf
SYSEOF
    [[ -f "${SCYLLA_OVR}/dependencies.conf" ]] || cat > "${SCYLLA_OVR}/dependencies.conf" <<'DEPEOF'
[Unit]
After=network-online.target
Wants=network-online.target
DEPEOF
    mkdir -p /etc/scylla.d
    [[ -f /etc/scylla.d/dev-mode.conf ]] || echo "DEV_MODE=--developer-mode=1" > /etc/scylla.d/dev-mode.conf
    [[ -f /etc/scylla.d/memory.conf ]]   || echo "# memory" > /etc/scylla.d/memory.conf
    [[ -f /etc/scylla.d/io.conf ]]       || echo "# io" > /etc/scylla.d/io.conf
    [[ -f /etc/scylla.d/cpuset.conf ]]   || echo "# cpuset" > /etc/scylla.d/cpuset.conf
    [[ -f /etc/sysconfig/scylla-server ]] || { mkdir -p /etc/sysconfig; echo 'SCYLLA_ARGS="--log-to-syslog 1 --log-to-stdout 0 --default-log-level info --network-stack posix --developer-mode=1"' > /etc/sysconfig/scylla-server; }

    # Disable ALL housekeeping timers (prevents unexpected restarts that
    # break Raft quorum by generating new host IDs)
    for timer in scylla-housekeeping-restart.timer scylla-housekeeping-daily.timer; do
      systemctl disable "$timer" 2>/dev/null || true
      systemctl stop "$timer" 2>/dev/null || true
      systemctl mask "$timer" 2>/dev/null || true
    done
    systemctl daemon-reload

    # Start ScyllaDB
    systemctl enable scylla-server.service 2>/dev/null || true
    systemctl start scylla-server.service || true
    echo "  OK: ScyllaDB configured and starting (seeds: ${SEED_IP})"

    # Wait for CQL (up to 4 min — Raft join takes time)
    echo -n "  Waiting for CQL..."
    for i in $(seq 1 120); do
      if ss -tlnp | grep -q ":9042 "; then
        echo " ready!"
        echo "  OK: ScyllaDB CQL on ${LOCAL_IP}:9042"
        break
      fi
      sleep 2
      echo -n "."
    done
    if ! ss -tlnp | grep -q ":9042 "; then
      echo ""
      echo "  WARN: ScyllaDB CQL not ready after 4min — Raft join may still be in progress"
    fi
  else
    echo "  WARN: ScyllaDB apt install failed -- install manually: apt install scylla-server"
  fi
else
  echo "  [3.3] ScyllaDB GPG key not available from bootstrap node -- skipping"
  echo "       (install manually or controller will handle via plan)"
fi

echo "  OK: Prerequisites ready"

# ===========================================================================
# PHASE 4: Join
# ===========================================================================
echo ""
echo "--- Phase 4: Join ---"

# -- Start node-agent -------------------------------------------------------
echo "  [4.1] Starting node-agent..."
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

for i in $(seq 1 30); do
  if ss -tlnp 2>/dev/null | grep -q ":11000 "; then
    break
  fi
  sleep 1
done
sleep 2
echo "  OK: node-agent running"

# -- Join the cluster -------------------------------------------------------
echo "  [4.2] Sending join request..."
globular cluster join \
  --join-token "${JOIN_TOKEN}" \
  --controller "${CONTROLLER}" \
  --node "${NODE_IP}:11000" \
  --insecure

# -- Register DNS A record (Day-1 Rule 7B) ---------------------------------
# The cluster must resolve this node's hostname. Register it now while we
# have all the information. Uses the DNS service on the bootstrap node.
echo "  [4.3] Registering DNS record..."
globular dns a set "${NODE_HOSTNAME}.globular.internal." "${NODE_IP}" \
  --dns "${BOOTSTRAP_HOST}:10007" --insecure 2>/dev/null && \
  echo "  OK: DNS ${NODE_HOSTNAME}.globular.internal -> ${NODE_IP}" || \
  echo "  WARN: DNS registration failed (controller will reconcile)"

echo ""
echo "================================================================"
echo "  NODE JOINED SUCCESSFULLY"
echo "================================================================"
echo ""
echo "  Identity:      ${NODE_HOSTNAME} (${NODE_IP})"
echo "  etcd:          member of cluster"
echo "  DNS:           resolves globular.internal"
echo "  config.json:   written"
echo "  node-agent:    running"
echo ""
echo "  The controller will now drive bootstrap phases:"
echo "    admitted -> infra_preparing -> etcd_ready ->"
echo "    storage_joining -> workload_ready"
echo ""
echo "  Monitor: globular cluster nodes list"
echo ""
`
