# Globular v1 Conformance-Driven Refactor Analysis

**Date:** 2026-02-06
**Scope:** Complete codebase audit against v1 invariants
**Approach:** Strict enforcement, zero backward compatibility

---

## Part 1: Invariant Violations

### INV-1: Identity MUST NOT Depend on Domain/Host/DNS

**Severity:** CRITICAL - Core system architecture violation

#### INV-1.1: User Identity Constructed from Domain Claims
- **File:** `internal/gateway/handlers/providers.go:121`
- **Pattern:** `claims.ID + "@" + claims.UserDomain`
- **Impact:** User identity is `"user123@globular.io"` - domain is part of identity string
- **Used In:**
  - `serve.go:272` - Access control checks
  - `upload.go:121-205` - Resource ownership
  - All authorization decisions
- **Violation:** Domain change breaks all identity, ownership, and ACL lookups

#### INV-1.2: Host Header Determines Filesystem Root
- **File:** `internal/gateway/handlers/files/serve.go:165-169`
- **Pattern:** `filepath.Join(dir, r.Host)`
- **Impact:** HTTP Host header directly selects filesystem directory
- **Attack:** Host header injection accesses arbitrary virtual host directories
- **Violation:** Network identity (untrusted header) controls file access context

#### INV-1.3: Domain-Based Object Storage Prefixes
- **Files:**
  - `internal/gateway/handlers/providers.go:593,596`
  - `internal/gateway/handlers/files/serve.go:510,566`
- **Pattern:** `path.Join(domain, "users")` → `/localhost/users/` vs `/prod.com/users/`
- **Impact:** MinIO storage paths differ per domain
- **Violation:** Data partitioning by domain - migration required on domain change

#### INV-1.4: Host Header in Object Storage Keys
- **File:** `internal/gateway/handlers/files/serve.go:566`
- **Pattern:** `base := path.Join(host, "webroot")`
- **Impact:** Request Host header constructs storage key
- **Attack:** Spoofed Host header accesses different tenant data
- **Violation:** Storage namespace determined by untrusted network identifier

#### INV-1.5: TLS Certificate Directories Use Domain
- **Files:**
  - `internal/gateway/handlers/config/ca_handlers.go:40`
  - `internal/gateway/httpserver/server.go:38`
- **Pattern:** `filepath.Join(config, "tls", domain)`
- **Impact:** Certificates stored in domain-specific directories
- **Violation:** Cryptographic identity coupled to domain configuration

#### INV-1.6: Resource Ownership Uses Domain-Embedded IDs
- **File:** `internal/gateway/handlers/files/upload.go:121-205`
- **Pattern:** Owner set to `"user@domain"` from INV-1.1
- **Impact:** Ownership lookups fail if domain portion mismatches
- **Violation:** ACL depends on domain string matching

#### INV-1.7: Domain in CORS Allowed Headers
- **Files:**
  - `internal/globule/globule.go:116`
  - `internal/controlplane/resource.go:503`
- **Pattern:** `AllowedHeaders: []string{"domain", ...}`
- **Impact:** Domain can be passed as HTTP header by clients
- **Violation:** Untrusted client-supplied header in auth/routing context

#### INV-1.8: Domain as Persistent State Key
- **File:** `internal/xds/watchers/watcher.go:1544-1555`
- **Pattern:** `key := fmt.Sprintf("/globular/pki/bundles/%s", domain)`
- **Impact:** etcd keys include domain - certificate tracking partitioned by domain
- **Violation:** Persistent state identity depends on domain value

#### INV-1.9: Hardcoded Internal Domain
- **File:** `internal/xds/watchers/watcher.go:1551`
- **Pattern:** `domain = "globular.internal"` (fallback)
- **Impact:** System-wide assumption of specific internal domain
- **Violation:** Internal DNS name hardcoded as identity

---

### INV-2: Token Verification MUST Use Asymmetric Crypto

**Status:** ✅ PASS with warnings

#### Positive: Ed25519 Only
- **File:** `services/golang/security/jwt.go:201,195,242`
- **Status:** ✅ Correct - EdDSA/Ed25519 asymmetric signing
- **Status:** ✅ No HMAC or symmetric algorithms found

#### Positive: Required Claims Present
- **File:** `services/golang/security/jwt.go:178-183`
- **Status:** ✅ `iss`, `aud`, `exp`, `sub` all included
- **Status:** ✅ Expiration enforced at parse time

---

### INV-3: Tokens MUST Be Header-Only (Not Query/Form)

**Severity:** MEDIUM - Token exposure in logs

#### INV-3.1: Query Parameter Token Extraction
- **Files:**
  - `internal/gateway/handlers/files/serve.go:172`
  - `internal/gateway/handlers/files/upload.go:70`
  - `internal/gateway/handlers/config/save_config.go:30-32`
- **Pattern:** `r.URL.Query().Get("token")` as fallback
- **Impact:** Tokens logged in access logs, proxy logs, browser history
- **Attack:** Referer header leaks, URL sharing exposes tokens
- **Violation:** Bearer tokens in URLs are fundamentally insecure

#### INV-3.2: Form Value Token Extraction
- **File:** `internal/gateway/handlers/files/upload.go:70`
- **Pattern:** `r.FormValue("token")`
- **Impact:** Tokens in request bodies, form caching
- **Violation:** Non-standard and less secure than Authorization header

---

### INV-4: Token Audience MUST Be Validated

**Severity:** MEDIUM - Token replay across services

#### INV-4.1: Audience Validation Disabled at Parse
- **File:** `services/golang/security/jwt.go:251`
- **Pattern:** `jwt.WithAudience("")` - explicit disable
- **Comment:** "we'll enforce aud at the service/router layer"
- **Impact:** Parser accepts tokens with any audience

#### INV-4.2: Service-Layer Audience Validation Missing
- **Files:** All gateway handlers in `internal/gateway/handlers/`
- **Pattern:** Token validated but audience never checked
- **Impact:** Token for peer A can be replayed to peer B
- **Attack:** Cross-service token replay
- **Violation:** Tokens not scoped to recipient

---

### INV-5: Domain MUST Be Routing Label Only (Not Identity)

**Severity:** CRITICAL - Architectural coupling

#### INV-5.1: Domain in Route Virtual Host Assignment
- **File:** `internal/xds/watchers/watcher.go:861-866`
- **Pattern:** `spec.Routes[i].Domains = []string{domain}`
- **Impact:** Route definitions injected with global domain config
- **Violation:** Route identity depends on domain configuration

#### INV-5.2: Domain as SNI in TLS Handshake
- **File:** `internal/xds/watchers/watcher.go:677,1184`
- **Pattern:** `cl.SNI = domain`
- **Impact:** Domain sent to remote servers during TLS
- **Violation:** Domain becomes TLS authentication identity

#### INV-5.3: Domain in Certificate DNSNames
- **File:** `internal/config/xds_paths.go:104-148`
- **Pattern:** `fmt.Sprintf("xds.%s", domain)` in cert DNSNames
- **Impact:** Certificate identity bound to domain
- **Violation:** Cannot change domain without regen certs

#### INV-5.4: Domain for Service Discovery (SRV Records)
- **File:** `internal/xds/watchers/watcher.go:1395-1405`
- **Pattern:** `dnsCache.LookupSRV(ctx, service, "tcp", clusterDomain)`
- **Impact:** Service lookup tied to cluster domain value
- **Violation:** Discovery depends on domain as identity

#### INV-5.5: Domain-Based Virtual Host Grouping
- **File:** `internal/controlplane/ingress.go:44,81-87`
- **Pattern:** Routes grouped by `domainsByKey[key]`
- **Impact:** Route organization keyed by domain
- **Violation:** Domain is primary partitioning dimension

---

### INV-6: DNS Operations MUST Use Reconciliation (Not Imperative)

**Severity:** MEDIUM - Operational pattern violation

#### INV-6.1: Certificate Rotation via Polling
- **File:** `internal/xds/watchers/watcher.go:1537-1601`
- **Pattern:** `checkCertificateGeneration()` polls etcd every 10s
- **Impact:** Imperative detection loop, not declarative reconciliation
- **Violation:** No desired-state vs actual-state model

#### INV-6.2: ACME Rotation via File Hash Polling
- **File:** `internal/xds/watchers/watcher.go:1595-1659`
- **Pattern:** `computeFileHash()` polls files every 10s
- **Impact:** Imperative file monitoring, not event-driven
- **Violation:** State changes detected via polling, not reconciliation

---

### INV-7: Production MUST NOT Allow Insecure Defaults

**Status:** ✅ PASS - Already enforced

- **File:** `cmd/globular-xds/main.go:415-489`
- **Pattern:** Fails fast if TLS certs missing (unless `GLOBULAR_XDS_INSECURE=1`)
- **Status:** ✅ Secure by default

---

## Part 2: Target v1 Architecture

### 2.1 Identity Model

**Principle:** Stable, opaque identifiers independent of network configuration

#### User Identity
```go
// ✅ CORRECT v1 Model
type PrincipalID string  // Opaque, stable UUID

// User struct (NO domain)
type User struct {
    ID          PrincipalID  `json:"principal_id"`  // "usr_7f9a3b2c"
    Username    string       `json:"username"`       // "alice" (display only)
    Email       string       `json:"email"`          // Contact, not identity
    CreatedAt   time.Time
    // NO UserDomain field
}

// ❌ WRONG (current)
userID := claims.ID + "@" + claims.UserDomain  // "alice@globular.io"
```

#### Cluster Identity
```go
// ✅ CORRECT v1 Model
type ClusterID string  // Opaque, stable UUID

type ClusterIdentity struct {
    ID                ClusterID   `json:"cluster_id"`      // "cls_a1b2c3d4"
    InternalDomain    string      `json:"cluster_domain"`  // DNS only, not identity
    PublicDomains     []string    `json:"ingress_domains"` // Routing only
}

// Lookups NEVER use domain
func GetCluster(clusterID ClusterID) (*Cluster, error)
// NOT: func GetClusterByDomain(domain string) - FORBIDDEN
```

#### Service Identity
```go
// ✅ CORRECT v1 Model
type ServiceID string  // "svc_echo_7a3f2b"

type Service struct {
    ID              ServiceID
    Name            string     // "echo.EchoService" (discovery label)
    ClusterID       ClusterID
    // NO domain field - domain is in ClusterIdentity only
}
```

---

### 2.2 Token Model

**Principle:** Asymmetric, stateless, audience-scoped

#### Token Structure
```go
// ✅ CORRECT v1 Model
type TokenClaims struct {
    // Standard JWT claims
    Issuer     string    `json:"iss"`  // Issuer cluster ID
    Audience   []string  `json:"aud"`  // Target service/cluster IDs
    Subject    string    `json:"sub"`  // Principal ID (NOT user@domain)
    ExpiresAt  int64     `json:"exp"`
    IssuedAt   int64     `json:"iat"`
    NotBefore  int64     `json:"nbf"`

    // Custom claims
    PrincipalID   string   `json:"principal_id"`  // "usr_7f9a3b2c"
    Scopes        []string `json:"scopes"`        // ["read:files", "write:config"]

    // NO UserDomain field
}
```

#### Token Validation Flow
```
1. Extract from Authorization header ONLY (no query/form)
2. Verify Ed25519 signature with issuer's public key
3. Check expiration (exp)
4. Validate audience matches current service/cluster ID
5. Extract principal_id (opaque, no domain parsing)
6. Load ACL for principal_id (NOT user@domain)
```

---

### 2.3 Domain Responsibilities

**Principle:** Strict separation of concerns

#### Cluster Domain (Internal DNS)
```yaml
cluster_domain: "cluster.local"  # RFC 6761 reserved
responsibilities:
  - Service discovery (SRV records)
  - Internal routing only
  - NOT used for identity
  - NOT in certificates (use SANs with service names)
  - NOT in storage paths
  - NOT in ACLs
```

#### Ingress Domains (Public Routing)
```yaml
ingress_domains:
  - "example.com"
  - "www.example.com"
responsibilities:
  - HTTP Host matching
  - ACME certificate subjects
  - Public DNS A/AAAA records
  - Virtual host routing
  - NOT used for identity
  - NOT in storage paths
  - NOT in user IDs
```

#### Separation Enforcement
```go
// ✅ CORRECT - Domain is routing config only
type IngressConfig struct {
    Domains     []string       // ["example.com"]
    CertSecret  string         // Reference to cert by ID, not domain
    Routes      []RouteConfig  // Routes by path, not domain-keyed
}

// ❌ WRONG - Domain influences state
key := fmt.Sprintf("/data/%s/users", domain)  // FORBIDDEN
```

---

### 2.4 Certificate Management

**Principle:** Zero-downtime, identity-independent rotation

#### Certificate Identity
```go
// ✅ CORRECT v1 Model
type Certificate struct {
    ID          string    `json:"cert_id"`      // "cert_7f9a3b"
    Type        CertType  `json:"type"`         // Internal | Public
    ClusterID   ClusterID `json:"cluster_id"`   // Which cluster owns it

    // Certificate content
    Cert        []byte
    Key         []byte
    CA          []byte

    // Metadata
    NotBefore   time.Time
    NotAfter    time.Time
    SANs        []string  // DNS names for matching (NOT identity)

    // NO domain field in ID or key
}

// Storage path
etcdKey := fmt.Sprintf("/globular/pki/certs/%s", cert.ID)
// NOT: /globular/pki/bundles/{domain}
```

#### Rotation Model
```go
// ✅ CORRECT v1 Model - Reconciliation
type CertificateReconciler struct {
    DesiredState  *Certificate  // What should exist
    ActualState   *Certificate  // What currently exists
}

func (r *CertificateReconciler) Reconcile(ctx context.Context) error {
    if r.needsRotation() {
        newCert := r.generateCertificate()
        r.updateSnapshot(newCert)
        r.notifyEnvoy()
    }
    return nil
}

// NOT: Polling-based checkCertificateGeneration()
```

---

### 2.5 Trust Boundaries

**Principle:** Explicit authorization, no implicit trust from network

#### Trust Model
```
┌─────────────────────────────────────────┐
│ UNTRUSTED (external input)             │
├─────────────────────────────────────────┤
│ - HTTP Host header                      │
│ - Query parameters                      │
│ - Form values                           │
│ - Client-supplied domain header         │
│ - DNS reverse lookups                   │
└─────────────────────────────────────────┘
             │
             v
┌─────────────────────────────────────────┐
│ AUTHENTICATION                          │
├─────────────────────────────────────────┤
│ - Extract token from Authorization hdr  │
│ - Verify Ed25519 signature             │
│ - Check expiration                      │
│ - Validate audience                     │
│ → Extract principal_id (opaque)         │
└─────────────────────────────────────────┘
             │
             v
┌─────────────────────────────────────────┐
│ AUTHORIZATION                           │
├─────────────────────────────────────────┤
│ - Load ACL for principal_id             │
│ - Check scopes/permissions              │
│ - Construct resource path (stable)      │
│ - NO domain in decision logic           │
└─────────────────────────────────────────┘
             │
             v
┌─────────────────────────────────────────┐
│ TRUSTED (internal state)                │
├─────────────────────────────────────────┤
│ - principal_id → ACL mapping            │
│ - cluster_id → config mapping           │
│ - cert_id → certificate content         │
│ - Stable storage paths                  │
└─────────────────────────────────────────┘
```

**Key Rules:**
1. Domain NEVER crosses from untrusted → trusted
2. Network identifiers NEVER influence ACL decisions
3. Storage paths NEVER include runtime configuration values
4. Authorization depends ONLY on principal_id + explicit policy

---

## Part 3: Required Refactors

### 3.1 Deletions

#### Remove Domain from Identity
```go
// DELETE these patterns:
- claims.ID + "@" + claims.UserDomain
- filepath.Join(dir, r.Host)
- path.Join(domain, "users")
- path.Join(host, "webroot")
- filepath.Join(config, "tls", domain)
- fmt.Sprintf("/globular/pki/bundles/%s", domain)

// DELETE these fields:
- Claims.UserDomain
- AllowedHeaders: []string{"domain"}
```

#### Remove Query/Form Token Extraction
```go
// DELETE these patterns:
- r.URL.Query().Get("token")
- r.FormValue("token")
- headerOrQuery(r, "token")
- firstNonEmpty(header, query, form)
```

#### Remove Domain from Configuration Keys
```go
// DELETE domain-based etcd keys:
- /globular/pki/bundles/{domain}
- /globular/config/{domain}/*

// REPLACE with ID-based keys:
- /globular/pki/certs/{cert_id}
- /globular/clusters/{cluster_id}/config
```

---

### 3.2 New Structs

#### Identity Types
```go
// File: internal/identity/types.go
package identity

type PrincipalID string

type Principal struct {
    ID         PrincipalID `json:"principal_id"`
    Username   string      `json:"username"`
    Email      string      `json:"email,omitempty"`
    CreatedAt  time.Time   `json:"created_at"`
}

type ClusterID string

type ClusterIdentity struct {
    ID             ClusterID `json:"cluster_id"`
    ClusterDomain  string    `json:"cluster_domain"`   // DNS only
    IngressDomains []string  `json:"ingress_domains"`  // Routing only
}
```

#### Token Claims (Refactored)
```go
// File: services/golang/security/claims.go
package security

type Claims struct {
    // Standard claims
    Issuer    string   `json:"iss"`  // cluster_id
    Audience  []string `json:"aud"`  // [service_id, cluster_id]
    Subject   string   `json:"sub"`  // principal_id
    ExpiresAt int64    `json:"exp"`
    IssuedAt  int64    `json:"iat"`
    NotBefore int64    `json:"nbf"`

    // Custom claims
    PrincipalID string   `json:"principal_id"`
    Scopes      []string `json:"scopes,omitempty"`

    // REMOVED: UserDomain
}
```

#### Certificate Storage
```go
// File: internal/pki/certificate.go
package pki

type CertificateID string

type Certificate struct {
    ID        CertificateID `json:"cert_id"`
    ClusterID ClusterID     `json:"cluster_id"`
    Type      CertType      `json:"type"`  // Internal | Public

    // Content
    CertPEM []byte `json:"cert_pem"`
    KeyPEM  []byte `json:"key_pem"`
    CAPEM   []byte `json:"ca_pem,omitempty"`

    // Metadata
    NotBefore time.Time `json:"not_before"`
    NotAfter  time.Time `json:"not_after"`
    SANs      []string  `json:"sans"`

    // Version for rotation tracking
    Generation uint64 `json:"generation"`
}

type CertType string
const (
    CertTypeInternal CertType = "internal"
    CertTypePublic   CertType = "public"
)
```

#### Storage Path Builder
```go
// File: internal/storage/paths.go
package storage

// GetPrincipalPath returns stable storage path for a principal
func GetPrincipalPath(principalID PrincipalID, resourceType string) string {
    // ALWAYS stable, NEVER includes domain/host
    return fmt.Sprintf("/principals/%s/%s", principalID, resourceType)
}

// GetCertificatePath returns stable path for certificate
func GetCertificatePath(certID CertificateID) string {
    return fmt.Sprintf("/pki/certs/%s", certID)
}

// FORBIDDEN: GetPathByDomain - must not exist
```

---

### 3.3 New Flows

#### Authentication Flow (Refactored)
```
1. Extract token from Authorization header
   ↓
2. Verify Ed25519 signature
   ↓
3. Check expiration
   ↓
4. Validate audience matches current service
   ↓
5. Extract principal_id (opaque UUID)
   ↓
6. Return authenticated principal (NO domain parsing)
```

#### Authorization Flow (Refactored)
```
1. Receive principal_id from authentication
   ↓
2. Load ACL: /acl/principals/{principal_id}
   ↓
3. Check permissions for requested action
   ↓
4. Construct resource path using stable IDs
   ↓
5. Grant/deny (domain NEVER influences decision)
```

#### Certificate Rotation (Reconciliation Pattern)
```
// Desired state
1. Admin/ACME requests new certificate
   ↓
2. Store desired cert in /pki/desired/{cert_id}
   ↓
3. Reconciler watches desired state
   ↓
4. Reconciler generates/fetches certificate
   ↓
5. Store actual cert in /pki/certs/{cert_id}
   ↓
6. Update snapshot with new cert
   ↓
7. Notify Envoy via xDS push
   ↓
8. Mark desired state as reconciled

// NOT: Poll etcd for generation changes
// NOT: Poll files for hash changes
```

---

## Part 4: Conformance Coverage Table

| Invariant | Fix | Test ID | File |
|-----------|-----|---------|------|
| **Identity & Authorization** |
| Identity stable, not domain-based | Remove `@domain` from user IDs | `TestUserIDNoDomain` | providers.go |
| No domain in user/group/role IDs | Refactor Claims struct | `TestClaimsNoDomain` | jwt.go |
| Host header doesn't influence access | Remove `r.Host` from paths | `TestHostHeaderIsolation` | serve.go |
| Domain change doesn't require migration | Remove domain from storage keys | `TestDomainChangeNoMigration` | storage/ |
| **Token & Security** |
| Tokens asymmetrically signed | Already correct (Ed25519) | `TestTokenAsymmetric` | jwt.go |
| Tokens include iss/aud/exp/sub | Already correct | `TestTokenClaims` | jwt.go |
| Tokens header-only | Remove query/form extraction | `TestTokenHeaderOnly` | serve.go, upload.go |
| Audience validated | Add service-layer check | `TestAudienceValidation` | handlers/ |
| **Domain & DNS** |
| Domain is routing label only | Remove from identity contexts | `TestDomainRoutingOnly` | watcher.go |
| Separate cluster_domain/ingress_domains | New ClusterIdentity struct | `TestDomainSeparation` | types.go |
| DNS reconciliation | Replace polling with reconciler | `TestDNSReconciliation` | reconciler.go |
| **Certificates & Ingress** |
| ACME for ingress domains only | Already correct | `TestACMEPublicOnly` | watcher.go |
| Zero-downtime rotation | Already correct (SDS) | `TestHotRotation` | sds_conformance_test.go |
| No internal domains in public certs | Verify SANs | `TestCertSANs` | certificate.go |
| **Configuration & Operations** |
| One canonical config location | Audit config paths | `TestConfigLocation` | config/ |
| Upgrades preserve operator config | Not in scope | - | installer |
| No default passwords in production | Not in scope | - | - |

---

## Part 5: Implementation Priority

### Phase 1: Identity Decoupling (CRITICAL)
1. Create `internal/identity/types.go` with PrincipalID
2. Refactor Claims struct (remove UserDomain)
3. Update ParseUserID to return PrincipalID only
4. Migrate storage paths from domain-based to ID-based
5. **Breaking Change:** All existing user IDs change format

### Phase 2: Token Security (HIGH)
1. Remove query/form token extraction
2. Add service-layer audience validation
3. Enforce Authorization header only
4. **Breaking Change:** Query param tokens no longer work

### Phase 3: Domain Separation (HIGH)
1. Create ClusterIdentity struct with separated domains
2. Remove domain from etcd keys
3. Refactor certificate storage to use cert_id
4. Update xDS watcher to use opaque IDs
5. **Breaking Change:** Certificate storage keys change

### Phase 4: Reconciliation (MEDIUM)
1. Implement certificate reconciler
2. Replace polling with event-driven updates
3. Add desired-state vs actual-state model
4. **Breaking Change:** Operational model changes

---

## Part 6: Migration Strategy (For Reference Only)

**Note:** Per instructions, we are NOT building migration layers. However, for operators' awareness:

### What Will Break
1. All existing tokens with `user@domain` format - **must reissue**
2. All storage paths with domain prefixes - **data migration required**
3. All certificates stored by domain - **must regenerate**
4. Query parameter authentication - **clients must update**
5. Host header-based routing - **config must be explicit**

### No Backward Compatibility
- Old token format: **not supported**
- Domain-based storage: **not readable**
- Query param tokens: **rejected**

This is a **v1.0 breaking release** - clean slate enforcement.

---

## Conclusion

The Globular codebase has **31 identified violations** of v1 invariants across identity, tokens, and domain management. The most critical issues are:

1. **User identity embeds domain** - breaks on domain change
2. **Untrusted headers influence access control** - security vulnerability
3. **Storage paths include domain/host** - requires migration on config change
4. **Token security gaps** - query params, missing audience validation
5. **Domain used as system identity** - architectural coupling

The target v1 architecture enforces:
- Opaque, stable principal IDs
- Domain as routing label only
- Asymmetric token security with audience validation
- Reconciliation-based certificate management
- Zero coupling between network config and identity

All fixes are **breaking changes** by design - this is v1 correctness enforcement, not incremental improvement.
