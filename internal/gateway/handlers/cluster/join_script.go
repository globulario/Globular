package cluster

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// NewJoinScriptHandler serves a self-contained shell script that bootstraps
// a new node to join the cluster. The script includes the CA certificate
// and downloads the node-agent + globularcli binaries from the gateway.
//
// Usage from a new node:
//
//	curl -sfL https://<gateway>:8443/join -k | sh -s -- --token <join-token>
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
		gatewayHost := r.Host
		if gatewayHost == "" {
			gatewayHost = r.URL.Host
		}
		// Strip port if present — we'll use our known port.
		if idx := strings.LastIndex(gatewayHost, ":"); idx > 0 {
			gatewayHost = gatewayHost[:idx]
		}
		if gatewayHost == "" {
			gatewayHost = "127.0.0.1"
		}

		gatewayAddr := fmt.Sprintf("%s:%d", gatewayHost, gatewayPort)

		// Controller address for the join command.
		ctrlAddr := controllerAddr
		if ctrlAddr == "" || ctrlAddr == "127.0.0.1:12000" || ctrlAddr == "localhost:12000" {
			ctrlAddr = fmt.Sprintf("%s:12000", gatewayHost)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, joinScript, caB64, gatewayAddr, ctrlAddr)
	})
}

// joinScript is the shell script template served to joining nodes.
// Args: %[1]s = base64 CA cert, %[2]s = gateway address, %[3]s = controller address
const joinScript = `#!/usr/bin/env bash
set -euo pipefail

# ── Globular Node Join Script ──────────────────────────────────────────────
# This script bootstraps a new node to join an existing Globular cluster.
# It installs the minimum required components (node-agent + CLI), sets up
# TLS trust, and runs the cluster join command.
#
# Usage:
#   curl -sfL https://<gateway>:8443/join -k | sh -s -- --token <join-token>
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
  echo "Usage: curl -sfL https://<gateway>:8443/join -k | sh -s -- --token <join-token>"
  exit 1
fi

echo "━━━ Globular Node Join ━━━"
echo "  Controller: ${CONTROLLER}"
echo "  Gateway:    ${GATEWAY}"
echo ""

# ── Check root ────────────────────────────────────────────────────────────
if [[ "$(id -u)" -ne 0 ]]; then
  echo "Error: this script must be run as root (sudo)"
  exit 1
fi

# ── Create user and directories ───────────────────────────────────────────
echo "  → Creating globular user and directories..."
if ! id -u globular &>/dev/null; then
  useradd --system --no-create-home --home-dir "${STATE_DIR}" --shell /usr/sbin/nologin globular
fi
mkdir -p "${INSTALL_DIR}" "${STATE_DIR}/pki" "${STATE_DIR}/services" "${STATE_DIR}/node_agent"
mkdir -p "${STATE_DIR}/config/tls"

# ── Install CA certificate ────────────────────────────────────────────────
echo "  → Installing cluster CA certificate..."
echo "${CA_B64}" | base64 -d > "${STATE_DIR}/pki/ca.crt"
cp "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
chmod 644 "${STATE_DIR}/pki/ca.crt" "${STATE_DIR}/pki/ca.pem"
echo "  ✓ CA certificate installed"

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
# Symlink for convenience.
ln -sf "${INSTALL_DIR}/globularcli" /usr/local/bin/globular
echo "  ✓ globularcli installed"

# ── Set ownership ─────────────────────────────────────────────────────────
chown -R globular:globular "${STATE_DIR}"

# ── Create node-agent systemd unit ────────────────────────────────────────
echo "  → Creating node-agent systemd unit..."
cat > "${SYSTEMD_DIR}/globular-node-agent.service" <<UNIT
[Unit]
Description=Globular node_agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=${STATE_DIR}/node_agent
Environment=GLOBULAR_SERVICES_DIR=${STATE_DIR}/services
ExecStartPre=/bin/sh -c 'mkdir -p ${STATE_DIR}/node_agent'
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

# Wait for node-agent to be ready.
echo "  → Waiting for node-agent..."
for i in $(seq 1 30); do
  if globular cluster join --help &>/dev/null 2>&1; then
    break
  fi
  sleep 1
done
sleep 3
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
echo "  The controller will now push desired state to this node."
echo "  Services will be installed automatically."
echo "  Monitor progress: globular services status"
echo ""
`
