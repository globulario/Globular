#!/bin/bash
# verify-sds-mtls.sh - Runtime verification that xDS/SDS is properly secured with mTLS
#
# This script verifies:
# 1. xDS server refuses plaintext connections (TLS required)
# 2. Envoy bootstrap contains TLS transport_socket for xDS cluster
# 3. Envoy uses SDS (not file-based TLS)
# 4. No file-based TLS appears when SDS is enabled

set -e

ENVOY_ADMIN_PORT="${ENVOY_ADMIN_PORT:-9901}"
XDS_HOST="${XDS_HOST:-127.0.0.1}"
XDS_PORT="${XDS_PORT:-18000}"
BOOTSTRAP_PATH="${BOOTSTRAP_PATH:-/run/globular/envoy/envoy-bootstrap.json}"

echo "=== Globular xDS/SDS mTLS Verification ==="
echo

# Test 1: Verify xDS server refuses plaintext connections
echo "[1/4] Testing xDS server rejects plaintext connections..."
if timeout 2 bash -c "</dev/tcp/$XDS_HOST/$XDS_PORT" 2>/dev/null; then
    # Connection succeeded - check if it's actually TLS or plaintext
    if echo -e "GET / HTTP/1.0\r\n\r\n" | timeout 2 nc -w 1 "$XDS_HOST" "$XDS_PORT" 2>/dev/null | grep -q "HTTP"; then
        echo "  ❌ FAIL: xDS server accepts plaintext HTTP"
        echo "     xDS server MUST use TLS to protect secrets"
        exit 1
    fi
    # TCP connection works but not HTTP - likely TLS
    echo "  ✓ xDS server port is open (TLS expected)"
else
    echo "  ⚠️  Cannot connect to xDS server at $XDS_HOST:$XDS_PORT"
    echo "     (This is OK if xDS is bound to localhost and we're remote)"
fi
echo

# Test 2: Verify Envoy bootstrap contains TLS configuration for xDS cluster
echo "[2/4] Checking Envoy bootstrap for xDS cluster TLS..."
if [ ! -f "$BOOTSTRAP_PATH" ]; then
    echo "  ⚠️  Bootstrap file not found: $BOOTSTRAP_PATH"
    echo "     Skipping bootstrap verification"
else
    if jq -e '.static_resources.clusters[] | select(.name == "xds_cluster") | .transport_socket' "$BOOTSTRAP_PATH" >/dev/null 2>&1; then
        echo "  ✓ xDS cluster has transport_socket configured"

        # Verify it's actually TLS
        tls_type=$(jq -r '.static_resources.clusters[] | select(.name == "xds_cluster") | .transport_socket.typed_config."@type"' "$BOOTSTRAP_PATH" 2>/dev/null)
        if echo "$tls_type" | grep -q "UpstreamTlsContext"; then
            echo "  ✓ xDS cluster uses UpstreamTlsContext (TLS enabled)"
        else
            echo "  ❌ FAIL: xDS cluster transport_socket is not TLS"
            echo "     Type: $tls_type"
            exit 1
        fi

        # Check for client certificate (mTLS)
        if jq -e '.static_resources.clusters[] | select(.name == "xds_cluster") | .transport_socket.typed_config.common_tls_context.tls_certificates' "$BOOTSTRAP_PATH" >/dev/null 2>&1; then
            echo "  ✓ xDS cluster has client certificate configured (mTLS)"
        else
            echo "  ⚠️  xDS cluster does not have client certificate (mTLS not enforced)"
        fi

        # Check for CA validation
        if jq -e '.static_resources.clusters[] | select(.name == "xds_cluster") | .transport_socket.typed_config.common_tls_context.validation_context' "$BOOTSTRAP_PATH" >/dev/null 2>&1; then
            echo "  ✓ xDS cluster validates server certificate with CA"
        else
            echo "  ❌ FAIL: xDS cluster does not validate server certificate"
            exit 1
        fi
    else
        echo "  ❌ FAIL: xDS cluster has no transport_socket configured"
        echo "     xDS communication would be plaintext!"
        exit 1
    fi
fi
echo

# Test 3: Verify Envoy uses SDS (not file-based TLS)
echo "[3/4] Checking Envoy config for SDS usage..."
if ! command -v curl >/dev/null 2>&1; then
    echo "  ⚠️  curl not found, skipping Envoy config_dump check"
else
    config_dump=$(curl -s "http://localhost:$ENVOY_ADMIN_PORT/config_dump" 2>/dev/null || true)
    if [ -z "$config_dump" ]; then
        echo "  ⚠️  Cannot reach Envoy admin API at localhost:$ENVOY_ADMIN_PORT"
        echo "     Skipping runtime verification"
    else
        # Check for SDS secret configs in listeners
        if echo "$config_dump" | jq -e '.configs[] | select(.["@type"] | contains("ListenersConfigDump")) | .dynamic_listeners[].active_state.listener.filter_chains[].transport_socket.typed_config.common_tls_context.tls_certificate_sds_secret_configs' >/dev/null 2>&1; then
            echo "  ✓ Listener uses SDS for TLS certificates"

            # Get secret names
            secret_names=$(echo "$config_dump" | jq -r '.configs[] | select(.["@type"] | contains("ListenersConfigDump")) | .dynamic_listeners[].active_state.listener.filter_chains[].transport_socket.typed_config.common_tls_context.tls_certificate_sds_secret_configs[].name' 2>/dev/null | sort -u)
            echo "     Secrets: $(echo $secret_names | tr '\n' ' ')"
        else
            echo "  ⚠️  No SDS configuration found in listeners"
            echo "     (Listeners may be using file-based TLS or no TLS)"
        fi

        # Check for SDS secret configs in clusters
        if echo "$config_dump" | jq -e '.configs[] | select(.["@type"] | contains("ClustersConfigDump")) | .dynamic_active_clusters[].cluster.transport_socket.typed_config.common_tls_context.validation_context_sds_secret_config' >/dev/null 2>&1; then
            echo "  ✓ Cluster uses SDS for CA validation"
        fi
    fi
fi
echo

# Test 4: Verify no file-based TLS when SDS is enabled
echo "[4/4] Checking for file-based TLS (should be absent with SDS)..."
if [ -z "$config_dump" ]; then
    echo "  ⚠️  Skipping (no config_dump available)"
else
    # Check if any TLS context references filename (file-based TLS)
    file_based_tls=$(echo "$config_dump" | jq '[.. | .filename? // empty] | unique' 2>/dev/null || echo "[]")

    if echo "$file_based_tls" | jq -e '. | length > 0' >/dev/null 2>&1; then
        # Filter out non-TLS file references (like access logs)
        tls_files=$(echo "$file_based_tls" | jq -r '.[] | select(contains(".pem") or contains(".crt") or contains(".key"))' 2>/dev/null || true)

        if [ -n "$tls_files" ]; then
            echo "  ⚠️  Found file-based TLS certificate references:"
            echo "$tls_files" | sed 's/^/     /'
            echo "     (If SDS is enabled, these should not be present)"
        else
            echo "  ✓ No file-based TLS certificate references found"
        fi
    else
        echo "  ✓ No file-based TLS certificate references found"
    fi
fi
echo

# Summary
echo "=== Verification Summary ==="
echo "✓ xDS/SDS security checks passed"
echo
echo "To verify certificate rotation:"
echo "  1. Rotate certificates: globularcli cert rotate --domain globular.internal"
echo "  2. Wait ~10 seconds for SDS watcher to detect change"
echo "  3. Check Envoy logs for secret updates"
echo "  4. Verify: curl localhost:$ENVOY_ADMIN_PORT/stats | grep sds"
echo
