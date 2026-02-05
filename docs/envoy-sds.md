# Envoy SDS (Secret Discovery Service) for Globular

## Overview

Globular implements Envoy's Secret Discovery Service (SDS) to enable **hot certificate rotation** without Envoy restarts. Certificates are delivered dynamically over gRPC, allowing seamless updates when the internal CA or ACME certificates rotate.

## Architecture

### Components

1. **SDS Secret Builders** (`internal/controlplane/secret.go`)
   - `MakeSecret()`: Builds xDS Secret resources from certificate files
   - Reads PEM files and embeds them as inline bytes
   - Supports internal certs (cluster CA) and public certs (ACME)
   - Content-based versioning using SHA256 hashing

2. **Certificate Rotation Detection** (`internal/xds/watchers/watcher.go`)
   - Polls etcd for certificate generation changes in sync loop
   - Detects when PKI manager rotates certificates
   - Triggers snapshot rebuild when generation changes
   - Pushes updated snapshot to Envoy via unified xDS cache

3. **xDS Integration** (`internal/xds/builder/builder.go`, `internal/controlplane/resource.go`)
   - SDS-aware TLS configuration functions
   - Secrets included in unified xDS snapshot (alongside CDS, LDS, RDS, EDS)
   - Envoy fetches secrets via ADS (Aggregated Discovery Service)
   - Single snapshot cache for all xDS resources

### Secret Types

| Secret Name | Purpose | Source |
|------------|---------|--------|
| `internal-server-cert` | Server cert for `*.globular.internal` | Cluster local CA |
| `internal-ca-bundle` | CA bundle for validating internal services | Cluster local CA |
| `internal-client-cert` | (Future) mTLS client identity | Cluster local CA |
| `public-server-cert` | ACME cert for public domain | Let's Encrypt/ACME |

## Certificate Rotation Flow

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐     ┌───────┐
│ PKI Manager │────>│ etcd (gen++) │────>│ xDS Watcher  │────>│ Envoy │
│  (renews)   │     │  /pki/bundles│     │ (snapshot++)│     │ (swap)│
└─────────────┘     └──────────────┘     └──────────────┘     └───────┘
                                                 │
                                                 v
                                          ┌──────────────┐
                                          │ xDS Snapshot │
                                          │ (CDS/LDS/RDS │
                                          │  + Secrets)  │
                                          └──────────────┘
```

1. **PKI Manager** rotates certificate (local CA or ACME renewal)
2. **etcd** generation counter increments (`/globular/pki/bundles/{domain}`)
3. **xDS Watcher** detects change via `checkCertificateGeneration()` polling
4. **xDS Watcher** rebuilds complete snapshot with updated secrets
5. **xDS Snapshot** pushed to unified cache via `SetSnapshot()`
6. **Envoy** receives ADS push notification and hot-swaps certificates
7. **Zero downtime**: existing connections continue, new connections use new cert

## Usage

### Enabling SDS

SDS is automatically enabled when TLS certificates are configured. The xDS watcher (`internal/xds/watchers/watcher.go`) detects TLS configuration and sets `EnableSDS: true` in the builder input:

```go
// In buildDynamicInput() - automatic SDS enablement
if listener.CertFile != "" && listener.KeyFile != "" {
    input.EnableSDS = true
    input.SDSSecrets = []builder.Secret{
        {
            Name:     "internal-server-cert",
            CertPath: listener.CertFile,
            KeyPath:  listener.KeyFile,
        },
    }
    if listener.IssuerFile != "" {
        input.SDSSecrets = append(input.SDSSecrets, builder.Secret{
            Name:   "internal-ca-bundle",
            CAPath: listener.IssuerFile,
        })
    }
}
```

The watcher also monitors certificate generation changes and triggers snapshot updates:

```go
// In sync() loop - automatic rotation detection
certChanged := w.checkCertificateGeneration(ctx)
if certChanged {
    w.logger.Info("certificate rotation detected - forcing snapshot update")
    // Snapshot rebuild triggered, includes updated secrets
}
```

No manual watcher setup is required - rotation detection is built into the xDS watcher.

## Configuration

### Feature Flag

- **`EnableSDS`**: Enable SDS for TLS certificates (default: `false`)
  - `true`: Use SDS (hot rotation, dynamic delivery)
  - `false`: Use file-based TLS (legacy, requires Envoy restart)

### Certificate Paths

Canonical PKI paths (stable across certificate rotations):

- **Server Cert**: `/var/lib/globular/pki/server.pem`
- **Server Key**: `/var/lib/globular/pki/server-key.pem`
- **CA Bundle**: `/var/lib/globular/pki/ca.pem`

## Envoy Configuration

### Downstream TLS (Listener)

**Without SDS** (file-based):
```yaml
transport_socket:
  name: envoy.transport_sockets.tls
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
    common_tls_context:
      tls_certificates:
        - certificate_chain:
            filename: /var/lib/globular/pki/server.pem
          private_key:
            filename: /var/lib/globular/pki/server-key.pem
```

**With SDS** (dynamic):
```yaml
transport_socket:
  name: envoy.transport_sockets.tls
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext
    common_tls_context:
      tls_certificate_sds_secret_configs:
        - name: internal-server-cert
          sds_config:
            resource_api_version: V3
            ads: {}
```

### Upstream TLS (Cluster)

**Without SDS** (file-based):
```yaml
transport_socket:
  name: envoy.transport_sockets.tls
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
    common_tls_context:
      validation_context:
        trusted_ca:
          filename: /var/lib/globular/pki/ca.pem
```

**With SDS** (dynamic):
```yaml
transport_socket:
  name: envoy.transport_sockets.tls
  typed_config:
    "@type": type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext
    common_tls_context:
      validation_context_sds_secret_config:
        name: internal-ca-bundle
        sds_config:
          resource_api_version: V3
          ads: {}
```

## Testing

### Verify SDS is Working

1. **Check Envoy config dump** to see SDS secrets:
   ```bash
   curl localhost:9901/config_dump | jq '.configs[] | select(.["@type"] | contains("SecretConfigDump"))'
   ```

2. **Check SDS stats** in Envoy:
   ```bash
   curl localhost:9901/stats | grep sds
   # Look for: sds.internal-server-cert.update_success
   ```

3. **Test certificate rotation**:
   ```bash
   # Initial request
   curl https://gateway.globular.internal

   # Rotate certificate (triggers etcd generation increment)
   globularcli cert rotate --domain globular.internal

   # Wait ~10 seconds for watcher to detect

   # Verify new cert served (no downtime)
   curl https://gateway.globular.internal
   ```

### Load Test During Rotation

```bash
# Start load test
hey -z 60s https://gateway.globular.internal/health &

# Rotate cert mid-test
sleep 10
globularcli cert rotate --domain globular.internal

# Check hey output - all requests should succeed (zero dropped)
```

## Migration Guide

### From File-Based TLS to SDS

1. **Day 0 (New Clusters)**: SDS enabled by default
   - `EnableSDS: true` in configuration
   - No file watching needed
   - Hot rotation works immediately

2. **Existing Clusters**: Gradual migration
   ```bash
   # Step 1: Verify certificates in canonical paths
   ls /var/lib/globular/pki/{server.pem,server-key.pem,ca.pem}

   # Step 2: Update xDS config
   # Set EnableSDS: true in xDS builder input

   # Step 3: Push new snapshot
   # Envoy hot-reloads config

   # Step 4: Verify Envoy switched to SDS
   curl localhost:9901/stats | grep sds.internal-server-cert

   # Step 5: Test rotation
   globularcli cert rotate --domain globular.internal
   # Wait 10s, verify cert updates without restart
   ```

3. **Rollback**: If issues arise
   ```bash
   # Set EnableSDS: false
   # Push snapshot
   # Envoy reverts to file-based TLS
   # No restart required
   ```

## Troubleshooting

### Symptoms and Fixes

| Symptom | Diagnosis | Fix |
|---------|-----------|-----|
| Envoy rejects SDS secret | Missing secret in snapshot | Check `SDSSecrets` in Input |
| Certificate not updating | Watcher not running | Verify watcher goroutine started |
| Certificate not updating | etcd generation not incrementing | Check PKI manager rotation logic |
| Old cert still served | Watcher polling too slow | Decrease `PollInterval` (default 10s) |
| SDS stats show errors | Invalid certificate file | Verify PEM format, check logs |

### Debug Commands

```bash
# Check xDS server logs
journalctl -u globular-xds -f

# Check Envoy logs
docker logs envoy -f

# Check certificate generation in etcd
etcdctl get /globular/pki/bundles/globular.internal

# Check current SDS version
curl localhost:9901/config_dump | jq '.configs[] | select(.["@type"] | contains("SecretConfigDump")) | .dynamic_active_secrets[].version_info'
```

## Security Considerations

### xDS Transport Security (Issue D)

The xDS gRPC server (which delivers secrets to Envoy) supports both insecure and mTLS modes:

**Insecure Mode** (default, localhost only):
- No TLS encryption
- Only safe when bound to `127.0.0.1` or `localhost`
- **WARNING**: Do not expose to network in this mode

**mTLS Mode** (recommended for production):
- Configure via `/var/lib/globular/xds/xds.yaml`:
  ```yaml
  tls:
    server_cert: /var/lib/globular/config/tls/fullchain.pem
    server_key: /var/lib/globular/config/tls/privkey.pem
    client_ca: /var/lib/globular/config/tls/ca.pem  # Requires client cert auth
  ```
- Uses cluster internal CA for server certificate
- Requires Envoy to present valid client certificate
- Secrets encrypted in transit over gRPC
- Protects against man-in-the-middle attacks

**Multi-Node Clusters**:
- **MUST** use mTLS mode
- xDS server must be reachable from other nodes
- Envoy bootstrap must reference correct xDS host/port
- Client certificates validated against internal CA

### Secret Protection

1. **Secrets in Transit**: Encrypted via TLS when mTLS mode enabled
2. **Secrets at Rest**: Standard file permissions on PKI directory (`/var/lib/globular/pki/`)
3. **No Insecure Fallback**: If SDS fails, Envoy rejects connections (fail-closed)
4. **Version Tracking**: Content-based hashing prevents replay attacks
5. **Client Authentication**: mTLS ensures only authorized Envoys can request secrets

## Performance

- **Polling Interval**: 10 seconds (configurable)
- **Rotation Latency**: ~10-15 seconds from cert renewal to Envoy swap
- **Connection Impact**: Zero dropped connections during rotation
- **CPU Overhead**: Negligible (polling + hashing)
- **Memory Overhead**: ~1MB per secret (cert + key in memory)

## References

- **Envoy SDS Documentation**: https://www.envoyproxy.io/docs/envoy/latest/configuration/security/secret
- **go-control-plane**: https://github.com/envoyproxy/go-control-plane
- **Globular PKI Implementation**: `services/golang/pki/`

## Future Enhancements

1. **SPIFFE Integration**: Use SPIFFE IDs for workload identity
2. **External Secret Stores**: Vault, AWS Secrets Manager
3. **Certificate Pinning**: Pin specific CA fingerprints
4. **CRL Support**: Certificate Revocation Lists
5. **OCSP Stapling**: Online Certificate Status Protocol
