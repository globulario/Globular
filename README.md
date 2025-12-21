<img src="https://globular.io/images/globular_logo.svg" alt="Globular Logo" width="300">

## [Visit globular.io](https://globular.io)

Globular is a **self-hosted application and service platform** designed to run, manage, and expose distributed services in a consistent, secure, and scalable way.

At its core, Globular is a **service management layer** built on gRPC, Envoy (xDS), and strong identity primitives. It provides everything required to run a personal or small-scale cloud: service lifecycle management, routing, authentication, authorization, storage, and media delivery—without relying on external SaaS platforms.

Globular is not a framework you embed into your app.  
It is an **operating environment** your applications run on.

---

## What Globular Is

Globular acts as a **control plane + runtime** for services:

- A **service registry and lifecycle manager** (start, stop, restart, monitor)
- A **secure gRPC-first execution model**
- An **Envoy / xDS–driven networking layer**
- A **resource and identity system** with RBAC baked in
- A **storage abstraction** spanning filesystem, object storage, databases, and torrents
- A **media-capable platform** (files, streaming, metadata, watching state)

You deploy services once, and Globular handles:
- discovery
- routing
- TLS
- access control
- persistence
- restarts
- exposure to web or native clients

---

## What Problems Globular Solves

### 1. Service Management (the missing piece of gRPC)
gRPC excels at defining APIs, but it does not solve:
- how services are discovered
- how versions coexist
- how routing changes dynamically
- how access is enforced
- how failures are handled

Globular exists specifically to solve **those operational problems**.

### 2. Identity and Access Control
Globular includes a first-class **resource model** with:
- users
- groups
- ownership
- permissions
- role-based access control (RBAC)

Access rules are enforced consistently across:
- gRPC methods
- HTTP gateways
- media streams
- filesystem paths
- APIs and UI components

### 3. Secure Networking by Default
Globular integrates:
- TLS and certificate management
- authenticated gRPC
- Envoy as the edge and internal router
- dynamic xDS configuration instead of hard-coded ports

Services do not need to know where other services live.  
They are **addressed by identity**, not by IP.

### 4. Storage Without Lock-In
Globular supports multiple storage backends transparently:
- local filesystem
- object storage (MinIO / S3-compatible)
- SQL databases (SQLite, MySQL, MariaDB, SQL Server)
- NoSQL / KV stores (MongoDB, ScyllaDB, Badger, LevelDB)
- time-series and metrics (Prometheus)
- torrents as mounted resources

Applications interact with **resources**, not vendors.

### 5. Media as a First-Class Citizen
Globular includes native support for:
- large media files
- streaming (MP4 / HLS)
- previews and thumbnails
- metadata extraction
- watching state and progress
- media organization and indexing

This makes Globular suitable for **personal media clouds**, not just APIs.

---

## Core Properties

Globular services are designed around the following guarantees:

- **Identifiable**  
  Every service instance has a unique identity bound to a domain and node.

- **Nameable**  
  Multiple instances can share a logical service name.

- **Versioned**  
  Clients can target specific service versions without breaking others.

- **Maintainable**  
  Services can be updated, restarted, or replaced without manual wiring.

- **Available**  
  Crashed services are automatically restarted.

- **Reachable**  
  Routing and addressing are handled by the platform, not hard-coded ports.

- **Trustable**  
  Authentication, authorization, and TLS are enforced by default.

- **Scalable**  
  Nodes can form a Globular Cluster, enabling a self-hosted cloud.

---

## What Globular Is *Not*

- Not a PaaS like Heroku
- Not Kubernetes (and not trying to be)
- Not a frontend framework
- Not a hosted SaaS

Globular favors **clarity, control, and self-hosting** over abstraction and vendor lock-in.

---

## Typical Use Cases

- Personal cloud (files, media, apps)
- Self-hosted alternatives to SaaS platforms
- Distributed applications with strong access control
- Media servers with identity-aware streaming
- Developer platforms where services evolve over time
- Edge or home-lab clusters without Kubernetes overhead

---

## Installation & Resources

- [General overview of Globular as a personal cloud](https://medium.com/@dave.courtois60/here-comes-globular-5dee34eb52f8)
- [Server installation and configuration guide](https://medium.com/@dave.courtois60/in-this-article-i-will-guide-you-through-the-installation-and-configuration-of-your-personal-cloud-f8bdce33d33a)
- [Installing Globular using Docker](https://medium.com/@dave.courtois60/installing-globular-using-docker-fabd4f96b095)
- [Docker image on Docker Hub](https://hub.docker.com/r/globular/globular)

---

## Project Status

**Globular v1.0 (beta)** is available.

The project is actively evolving toward:
- stronger clustering
- improved observability
- refined packaging and deployment
- richer developer tooling

Globular is built to grow incrementally—without breaking what already works.

---

## Philosophy

Globular is about **owning your infrastructure**.

It favors:
- explicit systems over hidden magic
- identity over addresses
- resources over vendors
- evolution over rewrites

If you want to run your own cloud—on your terms—Globular is designed for you.
