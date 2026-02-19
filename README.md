<img src="https://github.com/globulario/Globular/blob/master/starter_page/img/logo.png?raw=true" alt="Globular Logo" width="300">

# Globular

**A self-hosted application and service platform for distributed systems**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0--beta-orange.svg)](https://github.com/globulario/Globular/releases)

[Website](https://www.globular.io) • [Documentation](https://github.com/globulario/Globular/wiki) • [Getting Started](#getting-started)

---

## What is Globular?

Globular is a **service management platform** that provides everything you need to run a personal or small-scale cloud infrastructure. It's an **operating environment** for distributed applications—not a framework you embed, but a runtime your services run on.

Built on **gRPC**, **Envoy (xDS)**, and **strong identity primitives**, Globular handles service lifecycle, routing, security, storage, and media delivery without relying on external SaaS platforms.

### Core Capabilities

- **Service Registry & Lifecycle** - Deploy, discover, and manage microservices
- **Secure gRPC Execution** - TLS and mTLS enforced everywhere
- **Dynamic Routing** - Envoy proxy with xDS control plane
- **Identity & Access Control** - Built-in RBAC and resource ownership
- **Storage Abstraction** - Filesystem, S3, SQL, NoSQL, and more
- **Domain Management** - External DNS and ACME certificate automation
- **Media Streaming** - First-class support for files, video, and metadata

---

## Why Globular?

### The Missing Piece of gRPC

gRPC excels at defining APIs but doesn't solve operational challenges:
- ❌ How services discover each other
- ❌ How routing adapts dynamically
- ❌ How access is enforced
- ❌ How certificates are managed
- ❌ How failures are handled

**Globular exists to solve these problems.**

### Key Differentiators

| Feature | Globular | Traditional Approach |
|---------|----------|---------------------|
| **Service Discovery** | Automatic via etcd | Manual configuration |
| **Routing** | Dynamic (xDS) | Static ports/IPs |
| **Security** | TLS/mTLS mandatory | Optional, manual setup |
| **Identity** | First-class resources | External auth systems |
| **Domains** | ACME + DNS automation | Manual Let's Encrypt |
| **Storage** | Multi-backend abstraction | Vendor lock-in |

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│  External Traffic                                                │
│  (HTTPS on port 443)                                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ↓
              ┌──────────────────────┐
              │   Envoy Proxy        │  ← TLS Termination
              │   (xDS-configured)   │  ← Let's Encrypt Certs
              └──────────┬───────────┘
                         │
                         ↓
              ┌──────────────────────┐
              │   Gateway Service    │  ← HTTP/2 + gRPC
              │   (HTTPS:8443)       │  ← Internal TLS
              └──────────┬───────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ↓               ↓               ↓
   ┌─────────┐    ┌─────────┐    ┌─────────┐
   │ Service │    │ Service │    │ Service │  ← Microservices
   │   DNS   │    │  RBAC   │    │Resource │  ← mTLS Auth
   └─────────┘    └─────────┘    └─────────┘  ← etcd Discovery
         │               │               │
         └───────────────┼───────────────┘
                         │
                         ↓
         ┌───────────────────────────────────┐
         │  Infrastructure Layer              │
         │  • etcd (service registry)         │
         │  • MinIO (object storage)          │
         │  • PKI (certificate authority)     │
         └────────────────────────────────────┘
```

### Key Components

#### 1. **Envoy Proxy + xDS Control Plane**
- **Envoy** provides L7 load balancing and TLS termination
- **xDS** (globular-xds service) dynamically configures Envoy routes
- **Per-service subdomain routing**: `resource.globular.app` → Resource Service
- **Automatic endpoint discovery**: Services register in etcd, xDS builds clusters

#### 2. **Service Discovery (etcd)**
- All services register in etcd with health status
- Clients discover services by name, not IP
- Multi-instance support with automatic load balancing
- Cluster membership and configuration storage

#### 3. **Gateway Service**
- Unified entry point for HTTP/gRPC traffic
- Protocol translation (HTTP → gRPC, REST → gRPC)
- WebSocket support for streaming
- Request routing based on service name

#### 4. **Domain Reconciler**
- Manages external domains (e.g., `globular.app`)
- Publishes DNS records to CloudFlare, Route53, etc.
- Obtains ACME certificates from Let's Encrypt
- Supports wildcard certificates (`*.globular.app`)
- Automatic renewal 30 days before expiration

#### 5. **PKI Infrastructure**
- Self-signed CA for internal TLS
- Client certificates for user authentication
- Service certificates for mTLS
- Automatic certificate rotation

---

## Security Model

Globular implements **production-grade security** with multiple layers of defense:

### 1. **Mandatory TLS**
- No HTTP fallback - TLS required everywhere
- mTLS (mutual TLS) for service-to-service communication
- Separate certificate chains for internal and external traffic

### 2. **Bootstrap Mode**
- 30-minute window after installation
- Loopback-only access (127.0.0.1)
- Restricted to initial admin account creation
- Automatically disabled after window expires

### 3. **ClusterID Enforcement**
- Every cluster has a unique identifier
- Tokens validated against cluster ID
- Prevents cross-cluster token replay attacks
- Enforced on both unary and streaming RPCs

### 4. **Role-Based Access Control (RBAC)**
- Fine-grained permissions per service method
- Resource ownership (users, groups, applications)
- Wildcard admin role for system operations
- No hardcoded bypasses (production-ready)

### 5. **Authentication & Authorization**
- JWT tokens with short expiration
- Principal ID as canonical identity
- AuthContext as single source of truth
- Audit logging for all authorization decisions

### Security Best Practices

```bash
# ✓ Services authenticate via mTLS
# ✓ Users authenticate via client certificates
# ✓ Tokens validated against cluster ID
# ✓ Fail-closed security (deny by default)
# ✓ Bootstrap gate prevents unauthorized setup
```

---

## Service Catalog

### Core Infrastructure Services

| Service | Purpose | Port |
|---------|---------|------|
| **etcd** | Service registry & KV store | 2379 |
| **MinIO** | Object storage (S3-compatible) | 9000 |
| **Envoy** | Ingress proxy & load balancer | 443/8443 |
| **xDS** | Dynamic Envoy configuration | 18000 |

### Control Plane

| Service | Purpose | Port |
|---------|---------|------|
| **Gateway** | HTTP/gRPC gateway | 8080/8443 |
| **Cluster Controller** | Orchestration & reconciliation | 10000 |
| **Node Agent** | Per-node management | 11000 |
| **Discovery** | Service registry API | 10002 |

### Identity & Security

| Service | Purpose | Port |
|---------|---------|------|
| **Authentication** | User/service authentication | 10001 |
| **RBAC** | Role-based access control | 10003 |
| **Resource** | Resource ownership & permissions | 10007 |

### Data Services

| Service | Purpose | Port |
|---------|---------|------|
| **DNS** | Internal DNS (port 53 + gRPC) | 10006 |
| **Persistence** | Database abstraction | 10017 |
| **File** | File management & streaming | 10013 |
| **Log** | Centralized logging | 10012 |
| **Event** | Event bus & pub/sub | 10011 |

### Application Services

| Service | Purpose | Port |
|---------|---------|------|
| **Blog** | Blogging platform | 10018 |
| **Catalog** | Media cataloging | 10019 |
| **Conversation** | Messaging & chat | 10020 |
| **Media** | Media processing & streaming | 10023 |
| **Search** | Full-text search | 10025 |
| **Storage** | Storage backend abstraction | 10028 |

---

## Domain Management

Globular includes **automated external domain management** for production deployments.

### Features

- **External DNS Integration**: CloudFlare, Route53, DigitalOcean, etc.
- **ACME Certificate Acquisition**: Let's Encrypt (production & staging)
- **Wildcard Certificates**: `*.globular.app` + apex domain
- **DNS-01 Challenge**: Automated TXT record management
- **Auto-Renewal**: Certificates renewed 30 days before expiration
- **Ingress Configuration**: Automatic Envoy routing updates

### Example: Register External Domain

```bash
# Add DNS provider (CloudFlare)
globular domain provider add \
  --name cloudflare-prod \
  --type cloudflare \
  --zone globular.app \
  --api-token $CLOUDFLARE_TOKEN

# Register domain with wildcard cert
globular domain add \
  --fqdn globular.app \
  --zone globular.app \
  --provider cloudflare-prod \
  --target-ip auto \
  --publish-external \
  --use-wildcard-cert \
  --enable-acme \
  --acme-email admin@globular.app \
  --enable-ingress \
  --ingress-service gateway \
  --ingress-port 8443

# Check status
globular domain status --fqdn globular.app
```

### Domain Lifecycle

```
1. Domain Spec Created (etcd)
   ↓
2. DNS Record Published (CloudFlare/Route53)
   ↓
3. DNS-01 Challenge (TXT record)
   ↓
4. ACME Certificate Obtained (Let's Encrypt)
   ↓
5. Ingress Route Created (Envoy via xDS)
   ↓
6. Domain Ready (HTTPS:443 → Gateway:8443)
```

**Important**: External DNS records persist after uninstall. Always use `globular domain remove --cleanup-dns` before uninstalling.

---

## Installation

### Prerequisites

- **Linux**: Ubuntu 22.04+, Debian 12+, RHEL 9+
- **Architecture**: x86_64 (amd64)
- **RAM**: 2GB minimum, 4GB recommended
- **Disk**: 10GB minimum, 50GB+ for media storage
- **Ports**: 53 (DNS), 443 (HTTPS), 2379 (etcd), 8080/8443 (Gateway)

### Quick Start (Day-0 Installation)

```bash
# Clone installer repository
git clone https://github.com/globulario/globular-installer.git
cd globular-installer

# Build installer binary
./scripts/build.sh

# Install all services (requires sudo)
sudo ./scripts/install-day0.sh

# Verify installation
systemctl status globular-etcd
systemctl status globular-gateway
systemctl status globular-envoy

# Check cluster status
globular cluster status
```

### What Gets Installed

Day-0 installation deploys:

1. **Infrastructure**: etcd, MinIO, Envoy, xDS
2. **Control Plane**: Gateway, Cluster Controller, Node Agent
3. **Core Services**: DNS, Authentication, RBAC, Resource
4. **Data Services**: Persistence, File, Log, Event
5. **PKI**: CA certificate, service certificates, client certificates
6. **CLI**: `globular` command-line tool

### Post-Installation

```bash
# Create admin user (during 30-min bootstrap window)
globular admin create \
  --username admin \
  --password <secure-password> \
  --email admin@localhost

# Register your first external domain (optional)
globular domain add \
  --fqdn mycluster.example.com \
  --zone example.com \
  --provider cloudflare-prod \
  --enable-acme

# Access web interface
open https://localhost:8443
```

### Uninstallation

```bash
# Remove domains with DNS cleanup (IMPORTANT!)
globular domain status
globular domain remove --fqdn <domain> --cleanup-dns

# Uninstall all services
cd globular-installer
sudo ./scripts/uninstall-day0.sh
```

**⚠️ Warning**: Uninstall removes local data but does NOT clean up external DNS records. Always remove domains first!

---

## Configuration

### Cluster Configuration

Primary configuration file: `/var/lib/globular/config/config.json`

```json
{
  "name": "globular-cluster",
  "protocol": "https",
  "domain": "globular.internal",
  "etcd": {
    "endpoints": ["https://127.0.0.1:2379"],
    "ca": "/var/lib/globular/pki/ca.pem",
    "cert": "/var/lib/globular/pki/issued/etcd/client.crt",
    "key": "/var/lib/globular/pki/issued/etcd/client.key"
  }
}
```

### Service Discovery

Services automatically register in etcd:

```bash
# List all services
etcdctl get /services/ --prefix

# Discover service endpoint
globular service resolve dns.DnsService
```

### TLS Certificates

Certificate locations:

- **CA Certificate**: `/var/lib/globular/pki/ca.pem`
- **Service Certs**: `/var/lib/globular/pki/issued/<service>/`
- **Client Certs**: `~/.config/globular/tls/localhost/`
- **Domain Certs**: `/var/lib/globular/domains/<fqdn>/`

---

## Development

### Building Services

```bash
# Clone services repository
git clone https://github.com/globulario/services.git
cd services

# Build a specific service
cd golang/dns/dns_server
go build -o dns_server

# Generate package
cd ../../../
./scripts/package-service.sh dns
```

### Running Tests

```bash
# Security model tests
cd services/golang/security
go test -v

# RBAC tests
cd services/golang/rbac/rbac_server
go test -v

# Integration tests
cd services/golang/interceptors
go test -v
```

### Package Generation

Globular uses **declarative package specs** for deployment:

```yaml
# Example: DNS service spec (specs/dns_service.yaml)
metadata:
  name: dns
  version: 0.0.1
  description: DNS service with port 53 binding

steps:
  - type: ensure_user_group
    user: globular
    group: globular

  - type: install_package_payload
    install_binaries: true
    install_config: false

  - type: install_services
    services: [globular-dns.service]

  - type: enable_services
    services: [globular-dns.service]
```

---

## Storage Backends

Globular abstracts storage across multiple backends:

### Supported Storage Types

| Backend | Use Case | Configuration |
|---------|----------|---------------|
| **Filesystem** | Local files | Default, no config needed |
| **MinIO/S3** | Object storage | S3-compatible API |
| **SQLite** | Embedded DB | Single-file database |
| **MySQL/MariaDB** | Relational DB | Connection string |
| **MongoDB** | Document store | MongoDB URI |
| **ScyllaDB** | Wide-column store | Cassandra-compatible |
| **Badger/LevelDB** | Embedded KV | Local key-value stores |
| **Torrents** | Distributed files | Torrent infohash |

### Storage Configuration

```bash
# Configure S3 backend
globular storage add \
  --name s3-primary \
  --type s3 \
  --endpoint https://s3.amazonaws.com \
  --access-key $AWS_ACCESS_KEY \
  --secret-key $AWS_SECRET_KEY \
  --bucket globular-data

# Configure MongoDB
globular persistence add \
  --name mongo-primary \
  --type mongodb \
  --uri mongodb://localhost:27017/globular
```

---

## Use Cases

### Personal Cloud
Host your own files, photos, and media with complete ownership and privacy.

```bash
# Upload files
curl -X POST https://your-cluster.com/file-upload \
  -H "Authorization: Bearer $TOKEN" \
  -F file=@photo.jpg

# Stream media
curl https://your-cluster.com/media/stream/movie.mp4
```

### Self-Hosted SaaS Alternative
Replace Dropbox, Google Drive, Notion, or other SaaS platforms.

### Distributed Applications
Build microservices with automatic discovery, routing, and security.

### Media Server
Stream video/audio with identity-aware access control and progress tracking.

### Edge Computing
Deploy at home, office, or edge locations without Kubernetes complexity.

### Developer Platform
Prototype and deploy APIs with built-in versioning and lifecycle management.

---

## Clustering (Beta)

Globular supports **multi-node clusters** for high availability and horizontal scaling.

### Cluster Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Node 1    │────▶│   Node 2    │────▶│   Node 3    │
│ (Primary)   │     │ (Replica)   │     │ (Replica)   │
└─────────────┘     └─────────────┘     └─────────────┘
      │                   │                   │
      └───────────────────┴───────────────────┘
                      etcd Raft
```

### Joining a Cluster

```bash
# On primary node
globular cluster init --name production

# On additional nodes
globular cluster join \
  --primary https://node1.globular.internal:2379 \
  --token $CLUSTER_TOKEN
```

---

## Monitoring & Observability

### Health Checks

```bash
# Check service health
globular health check --service dns.DnsService

# View cluster health
globular cluster health

# Service status
systemctl status globular-*
```

### Logs

```bash
# View service logs
journalctl -u globular-gateway -f

# Domain reconciler logs
journalctl -u globular-cluster-controller -f | grep domain

# Envoy access logs
journalctl -u globular-envoy -f
```

### Metrics (Coming Soon)
- Prometheus integration
- Grafana dashboards
- Service-level metrics

---

## Roadmap

### Version 1.1 (Q2 2026)
- [ ] Enhanced clustering with automatic failover
- [ ] Web-based admin console
- [ ] Metrics and monitoring (Prometheus + Grafana)
- [ ] Backup and restore utilities

### Version 1.2 (Q3 2026)
- [ ] Multi-region support
- [ ] Database replication
- [ ] CDN integration
- [ ] Mobile SDK (iOS/Android)

### Version 2.0 (Q4 2026)
- [ ] Plugin system for custom services
- [ ] Marketplace for community services
- [ ] Enhanced developer tooling
- [ ] Kubernetes integration (optional)

---

## Contributing

Contributions are welcome! Please see:

- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Development Setup](docs/DEVELOPMENT.md)

### Quick Contribution Guide

```bash
# Fork and clone
git clone https://github.com/YOUR_USERNAME/Globular.git
cd Globular

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
go test ./...

# Commit with conventional commits
git commit -m "feat(dns): add IPv6 support"

# Push and create PR
git push origin feature/amazing-feature
```

---

## Community

- **Website**: [globular.io](https://globular.io)
- **GitHub**: [github.com/globulario/Globular](https://github.com/globulario/Globular)
- **Issues**: [Bug Reports & Feature Requests](https://github.com/globulario/Globular/issues)
- **Discussions**: [GitHub Discussions](https://github.com/globulario/Globular/discussions)

---

## Philosophy

Globular is about **owning your infrastructure**.

We favor:
- **Explicit systems** over hidden magic
- **Identity** over IP addresses
- **Resources** over vendor lock-in
- **Evolution** over rewrites
- **Self-hosting** over SaaS dependency

If you want to run your own cloud—on your terms—Globular is designed for you.

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Acknowledgments

Built with:
- [gRPC](https://grpc.io/) - High-performance RPC framework
- [Envoy](https://www.envoyproxy.io/) - Cloud-native proxy
- [etcd](https://etcd.io/) - Distributed key-value store
- [MinIO](https://min.io/) - S3-compatible object storage
- [Let's Encrypt](https://letsencrypt.org/) - Free TLS certificates

---

<p align="center">
  <strong>Made with ❤️ by the Globular Community</strong><br>
  <sub>Self-hosted. Open Source. Your Cloud.</sub>
</p>
