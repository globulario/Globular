# Envoy SDS (Secret Discovery Service) for Globular

## Overview

Globular implements Envoy's Secret Discovery Service (SDS) to enable **hot certificate rotation** without Envoy restarts. Certificates are delivered dynamically over gRPC, allowing seamless updates when the internal CA or ACME certificates rotate.

## Architecture

### Components

1. **SDS Secret Builders** (`internal/xds/sds/`)
   - `secrets.go`: Builds xDS Secret resources from certificate files
   - Reads from canonical PKI paths (`/var/lib/globular/pki/`)
   - Supports internal certs (cluster CA) and public certs (ACME)

2. **Certificate Rotation Watcher** (`internal/xds/sds/watcher.go`)
   - Polls etcd every 10 seconds for certificate generation changes
   - Detects when PKI manager rotates certificates
   - Reloads secrets from disk and pushes to Envoy via xDS snapshot

3. **xDS Integration** (`internal/controlplane/resource.go`, `internal/xds/builder/`)
   - SDS-aware TLS configuration functions
   - Secrets included in unified xDS snapshot (alongside CDS, LDS, RDS)
   - Envoy fetches secrets via ADS (Aggregated Discovery Service)

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
│ PKI Manager │────>│ etcd (gen++) │────>│  SDS Watcher │────>│ Envoy │
│  (renews)   │     │  /pki/bundles│     │  (reloads)   │     │ (swap)│
└─────────────┘     └──────────────┘     └──────────────┘     └───────┘
                                                 │
                                                 v
                                          ┌──────────────┐
                                          │ xDS Snapshot │
                                          │   (secrets)  │
                                          └──────────────┘
```

1. **PKI Manager** rotates certificate (local CA or ACME renewal)
2. **etcd** generation counter increments (`/globular/pki/bundles/{domain}`)
3. **SDS Watcher** detects change via polling
4. **SDS Watcher** reloads secrets from disk (`/var/lib/globular/pki/`)
5. **xDS Snapshot** updated with new secret versions
6. **Envoy** receives push notification and hot-swaps certificates
7. **Zero downtime**: existing connections continue, new connections use new cert

## Usage

### Enabling SDS

SDS is controlled by the `EnableSDS` flag in the xDS builder input:

```go
import (
    "github.com/globulario/Globular/internal/xds/builder"
    "github.com/globulario/Globular/internal/xds/sds"
)

// Build snapshot with SDS enabled
input := builder.Input{
    NodeID:    "envoy-node",
    EnableSDS: true,  // Enable SDS
    SDSSecrets: []builder.Secret{
        {
            Name:     "internal-server-cert",
            CertPath: "/var/lib/globular/pki/server.pem",
            KeyPath:  "/var/lib/globular/pki/server-key.pem",
        },
        {
            Name:     "internal-ca-bundle",
            CAPath:   "/var/lib/globular/pki/ca.pem",
        },
    },
    Listener: builder.Listener{
        Name: "https_listener",
        Port: 443,
        // ... other config
    },
    // ... clusters, routes, etc.
}

snapshot, err := builder.BuildSnapshot(input, "v1")
```

### Starting the Rotation Watcher

```go
import (
    "context"
    "github.com/globulario/Globular/internal/xds/sds"
    clientv3 "go.etcd.io/etcd/client/v3"
)

// Create etcd client
etcdClient, _ := clientv3.New(clientv3.Config{
    Endpoints: []string{"localhost:2379"},
})

// Create SDS watcher
watcher, err := sds.NewWatcher(sds.WatcherConfig{
    EtcdClient:       etcdClient,
    SDSServer:        sdsServer,  // Your SDS server instance
    Domain:           "globular.internal",
    InternalCertPath: "/var/lib/globular/pki/server.pem",
    InternalKeyPath:  "/var/lib/globular/pki/server-key.pem",
    CAPath:           "/var/lib/globular/pki/ca.pem",
    PollInterval:     10 * time.Second,
})

// Start watcher (runs until context cancelled)
ctx := context.Background()
go watcher.Run(ctx)
```

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

1. **mTLS Required**: SDS gRPC uses same mTLS as xDS control plane
2. **No Insecure Fallback**: If SDS fails, Envoy rejects connections (fail-closed)
3. **Secrets in Transit**: Encrypted via TLS (control plane → Envoy)
4. **Secrets at Rest**: Standard file permissions on PKI directory
5. **Version Tracking**: Content-based hashing prevents replay attacks

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
