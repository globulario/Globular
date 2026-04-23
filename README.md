<img src="https://github.com/globulario/Globular/blob/master/starter_page/img/logo.png?raw=true" alt="Globular Logo" width="300">

# Globular

**A self-hosted platform for running, operating, and evolving distributed services**

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Website](https://img.shields.io/badge/website-globular.io-black.svg)](https://globular.io)

[Website](https://globular.io) • [Organization](https://github.com/globulario) • [Core Services](https://github.com/globulario/services) • [Installer](https://github.com/globulario/globular-installer) • [Quickstart](https://github.com/globulario/globular-quickstart) • [Admin Console](https://github.com/globulario/globular-admin)

---

## What is Globular?

Globular is a **self-hosted distributed platform** for running native services with **explicit workflows**, **observable convergence**, **strong identity and security**, and **dynamic routing**.

It is not a library you embed into an application. It is an **operating environment** for services and distributed applications that need lifecycle management, discovery, routing, storage, policy, and day-0 to day-2 operational tooling.

Globular is built around a few core ideas:

- **Native services first**: run services as native binaries under systemd
- **etcd as the single source of truth**: cluster state lives in one authoritative place
- **Workflow-driven convergence**: orchestration is explicit and inspectable
- **4-layer state model**: Repository → Desired → Installed → Runtime
- **Strong identity everywhere**: TLS, mTLS, RBAC, and resource ownership are part of the platform, not bolted on later

In practice, Globular gives you a platform to bootstrap a cluster, publish and install packages, route traffic through Envoy/xDS, operate services across nodes, and manage the system through workflows and admin tools.

---

## Repository Map

Globular is split across several repositories, each with a clear role:

| Repository | Role |
|------------|------|
| **[Globular](https://github.com/globulario/Globular)** | Platform entry point, top-level architecture, gateway/xDS layer, project overview |
| **[services](https://github.com/globulario/services)** | Core backend and control-plane services, proto definitions, package build pipeline, main documentation |
| **[globular-installer](https://github.com/globulario/globular-installer)** | Internal day-0 bootstrap and installation logic used by the release installer (`install.sh`); mainly relevant to contributors and packaging work |
| **[globular-quickstart](https://github.com/globulario/globular-quickstart)** | Docker-based cluster simulation, testing, and operational scenarios |
| **[globular-admin](https://github.com/globulario/globular-admin)** | Official admin console, web components, TypeScript SDK layer, and management UI |

If you are new to the project:

- Start here to understand **what Globular is**
- Go to **services** to explore the core platform implementation
- Use **[services releases](https://github.com/globulario/services/releases)** to download installable release artifacts
- Use **globular-quickstart** to simulate and validate cluster behavior
- Use **globular-admin** to manage and operate the platform from the UI

---

## Why Globular?

Modern distributed systems often force you into one of two awkward worlds:

- a pile of independent services with ad hoc scripts, scattered state, and fragile operations
- or a large orchestration stack that hides too much behavior behind implicit control loops

Globular takes a different path.

It is designed for people who want:

- **self-hosted infrastructure** without heavy platform dependency
- **distributed services** with first-class identity, routing, and lifecycle management
- **explicit operational models** instead of invisible automation
- **auditability and debuggability** when something drifts, fails, or needs repair
- **local-first development** with a path to clustered deployment

### What makes it different

| Concern | Globular Approach |
|---------|-------------------|
| **State management** | etcd as the single source of truth |
| **Orchestration** | Workflow-driven, explicit, inspectable |
| **Routing** | Envoy with xDS dynamic configuration |
| **Service execution** | Native binaries under systemd |
| **Security** | TLS, mTLS, RBAC, resource ownership |
| **Package lifecycle** | Repository-backed artifacts and desired/installed/runtime separation |
| **Operations** | Day-0 bootstrap, day-1 rollout, day-2 repair and observation |

Globular is especially relevant when you want a platform for **native services on bare metal or VMs** and do not want container orchestration to be the only story in town.

---

## High-Level Architecture

```text
External Traffic
    ↓
Envoy Gateway / Ingress
    ↓
Gateway + xDS-Controlled Routing
    ↓
Core Platform Services
    ├── Cluster Controller
    ├── Node Agent
    ├── Workflow Service
    ├── Repository
    ├── Authentication
    ├── RBAC
    ├── Resource
    ├── DNS
    ├── Monitoring
    └── Additional platform and application services
    ↓
Infrastructure Layer
    ├── etcd        (source of truth)
    ├── MinIO       (artifacts / objects / backups)
    ├── ScyllaDB    (high-throughput data)
    └── systemd     (native service execution)
```

### Core ideas behind the architecture

#### 1. Explicit workflows
Operational changes are represented as workflows instead of being buried inside opaque automation. This makes rollouts, repair, recovery, and convergence easier to inspect and reason about.

#### 2. 4-layer state model
Globular distinguishes between:

- **Repository**: what artifacts exist
- **Desired**: what the cluster wants
- **Installed**: what nodes have actually installed
- **Runtime**: what is currently healthy and running

That separation is a major part of how the platform stays understandable.

#### 3. Dynamic routing with Envoy and xDS
Traffic enters through Envoy and is routed dynamically based on platform state, rather than static host files and manual port wiring.

#### 4. Native service platform
Globular manages services as native binaries under systemd. It is designed for environments where that model is a feature, not a compromise.

---

## What Globular Does

Globular provides a platform layer across the full lifecycle of distributed services.

### Day-0: Bootstrap
- install the minimum infrastructure and control plane
- initialize identity and certificates
- establish the initial state foundation

### Day-1: Publish and deploy
- build and package services
- publish artifacts to the repository
- install and activate services on nodes
- route traffic dynamically through xDS and gateway components

### Day-2: Operate and evolve
- inspect workflows and state transitions
- detect drift and health issues
- run repair or reseed flows
- observe service/runtime status through admin tools and simulations

---

## Main Building Blocks

### Platform shell and routing
This repository contains the **gateway/xDS-facing project shell** and serves as the umbrella entry point for the platform.

### Core services
The **[services](https://github.com/globulario/services)** repository contains the backend and control-plane implementation, including:

- cluster controller
- node agent
- workflow engine
- repository service
- authentication and RBAC
- DNS and resource services
- monitoring, backup, and AI-assisted operational services

### Install

Installable releases are published from the **[services releases page](https://github.com/globulario/services/releases)**.

Example install flow on Linux:

```bash
VERSION="1.0.56"

curl -LO "https://github.com/globulario/services/releases/download/v${VERSION}/globular-${VERSION}-linux-amd64.tar.gz"
curl -LO "https://github.com/globulario/services/releases/download/v${VERSION}/globular-${VERSION}-linux-amd64.tar.gz.sha256"
/usr/bin/sha256sum -c "globular-${VERSION}-linux-amd64.tar.gz.sha256"

tar xzf "globular-${VERSION}-linux-amd64.tar.gz"
cd "globular-${VERSION}-linux-amd64"
sudo bash install.sh
```

The release archive is the user-facing install entry point. The `install.sh` script uses the logic from **[globular-installer](https://github.com/globulario/globular-installer)** behind the scenes, but most end users should start from the published release artifacts, not from the installer repository directly.

### Quickstart and simulation
The **[globular-quickstart](https://github.com/globulario/globular-quickstart)** repository provides a **Docker-based simulation environment** for validating platform behavior with production-style binaries and realistic scenarios.

### Admin console
The **[globular-admin](https://github.com/globulario/globular-admin)** repository contains the **official management UI**, along with frontend building blocks, TypeScript SDK/components, and the web admin application. It can later be wrapped as a Tauri desktop app.

---

## Start Here

### Understand the platform
- Read this README for the big picture
- Explore the docs in **[services/docs](https://github.com/globulario/services/tree/master/docs)**

### Explore the core implementation
- Open **[services](https://github.com/globulario/services)**
- Review `proto/`, `golang/`, and the architecture/operator docs

### Install a cluster
- Download an installable release from **[services releases](https://github.com/globulario/services/releases)**
- Extract the archive and run `sudo bash install.sh`

### Simulate and test
- Use **[globular-quickstart](https://github.com/globulario/globular-quickstart)** to run scenario-based simulations and validate behavior before or alongside real deployments

### Operate through the UI
- Use **[globular-admin](https://github.com/globulario/globular-admin)** for administration, management, and the official operator experience

### Work on installer internals
- Use **[globular-installer](https://github.com/globulario/globular-installer)** only if you are contributing to the bootstrap/install pipeline itself

---

## Who Globular Is For

Globular is aimed at people who want to operate distributed systems with more control and less fog.

It is a strong fit for:

- self-hosters and infrastructure builders
- operators who want **explicit workflows** and **visible state transitions**
- teams deploying **native services** on bare metal or VMs
- projects that need identity, routing, storage, and lifecycle management in one platform layer
- developers building local-first systems that can grow into clustered deployments

---

## Philosophy

Globular is about **owning your infrastructure**.

We favor:

- **Explicit systems** over hidden magic
- **Observable operations** over black-box orchestration
- **Identity** over ad hoc trust
- **Clear state boundaries** over collapsed abstractions
- **Self-hosting** over SaaS dependency
- **Evolution** over rewrite culture

The goal is not to hide distributed systems behind a curtain.

The goal is to make them **operable, understandable, and yours**.

---

## Project Status

Globular is an actively evolving platform. The project now spans multiple repositories with a clearer separation between platform shell, core services, installer, simulation environment, and admin tooling.

The center of gravity for the core backend implementation and documentation is currently the **[services](https://github.com/globulario/services)** repository.

---

## Community and Links

- **Website**: [globular.io](https://globular.io)
- **Organization**: [github.com/globulario](https://github.com/globulario)
- **Main Platform Entry**: [Globular](https://github.com/globulario/Globular)
- **Core Services**: [services](https://github.com/globulario/services)
- **Installer**: [globular-installer](https://github.com/globulario/globular-installer)
- **Quickstart**: [globular-quickstart](https://github.com/globulario/globular-quickstart)
- **Admin Console**: [globular-admin](https://github.com/globulario/globular-admin)
- **Medium Article**: [Meet Globular](https://medium.com/@dave.courtois60/meet-globular-1e21424ebb45)

---

## License

MIT License. See [LICENSE](LICENSE).

---

<p align="center">
  <strong>Self-hosted. Explicit. Operable.</strong><br>
  <sub>Globular is a platform for running distributed services on your terms.</sub>
</p>
