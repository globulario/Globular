package watchers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/globulario/Globular/internal/controllerclient"
	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/dnscache"
	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/secrets"
	"github.com/globulario/Globular/internal/xds/server"
	clustercontrollerpb "github.com/globulario/services/golang/clustercontroller/clustercontrollerpb"
	"github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type DownstreamTLSMode string

const (
	defaultNodeID    = "globular-xds"
	defaultInterval  = 5 * time.Second
	defaultRouteName = "ingress_routes"
	listenerNameBase = "ingress_listener"
)

const (
	DownstreamTLSDisabled DownstreamTLSMode = "disabled"
	DownstreamTLSOptional DownstreamTLSMode = "optional"
	DownstreamTLSRequired DownstreamTLSMode = "required"
)

// ParseDownstreamTLSMode returns a normalized downstream TLS mode value.
func ParseDownstreamTLSMode(value string) DownstreamTLSMode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(DownstreamTLSDisabled):
		return DownstreamTLSDisabled
	case string(DownstreamTLSRequired):
		return DownstreamTLSRequired
	default:
		return DownstreamTLSOptional
	}
}

var errNoChange = errors.New("config unchanged")

// EndpointSource indicates how an endpoint was resolved (PR5)
type EndpointSource int

const (
	EndpointSourceUnknown   EndpointSource = iota
	EndpointSourceSRV                      // Resolved via SRV record
	EndpointSourceA                        // Resolved via A/AAAA record
	EndpointSourceRegistry                 // From service registry (IP)
	EndpointSourceLocalhost                // Localhost fallback
)

func (s EndpointSource) String() string {
	switch s {
	case EndpointSourceSRV:
		return "SRV"
	case EndpointSourceA:
		return "A"
	case EndpointSourceRegistry:
		return "REGISTRY"
	case EndpointSourceLocalhost:
		return "LOCALHOST"
	default:
		return "UNKNOWN"
	}
}

// EndpointIdentity represents a canonical service endpoint identity (PR5)
type EndpointIdentity struct {
	ServiceDNSLabel string         // Normalized service name for DNS (e.g., "echo-echoservice")
	TargetFQDN      string         // Preferred FQDN (e.g., "node-01.cluster.local")
	TargetIP        string         // Fallback IP address
	Port            int            // Service port
	Source          EndpointSource // How this endpoint was resolved
}

// SortKey returns a stable sorting key for deterministic ordering (PR5)
// Ordering: FQDN → IP → Port
func (e EndpointIdentity) SortKey() string {
	if e.TargetFQDN != "" {
		return fmt.Sprintf("fqdn:%s:%d", e.TargetFQDN, e.Port)
	}
	return fmt.Sprintf("ip:%s:%d", e.TargetIP, e.Port)
}

// Host returns the preferred host identifier (FQDN if available, else IP)
func (e EndpointIdentity) Host() string {
	if e.TargetFQDN != "" {
		return e.TargetFQDN
	}
	return e.TargetIP
}

// logRateLimiter provides rate-limited logging per service per event type (PR5)
type logRateLimiter struct {
	mu       sync.RWMutex
	lastLogs map[string]time.Time // key format: "service:event_type"
	interval time.Duration
}

func newLogRateLimiter(interval time.Duration) *logRateLimiter {
	return &logRateLimiter{
		lastLogs: make(map[string]time.Time),
		interval: interval,
	}
}

func (l *logRateLimiter) shouldLog(service, eventType string) bool {
	key := service + ":" + eventType
	l.mu.RLock()
	lastLog, exists := l.lastLogs[key]
	l.mu.RUnlock()

	if !exists || time.Since(lastLog) >= l.interval {
		l.mu.Lock()
		l.lastLogs[key] = time.Now()
		l.mu.Unlock()
		return true
	}
	return false
}

// Watcher reloads xDS configuration from file or service registry and pushes snapshots to the server.
type Watcher struct {
	logger              *slog.Logger
	server              *server.XDSServer
	configPath          string
	nodeID              string
	interval            time.Duration
	downstreamTLSMode   DownstreamTLSMode
	downstreamTLSWarned bool

	protocol                    string
	gatewayTLSWarned            bool
	downstreamTLSRequiredWarned bool
	ingressTLSWarned            bool
	ingressMTLSWarned           bool
	ingressPortCollisionWarned  bool
	configCached                *XDSConfig
	configMod                   time.Time

	// DNS-first routing (PR4)
	controllerAddr string
	clusterNetwork *clustercontrollerpb.ClusterNetwork
	dnsCache       *dnscache.Cache

	// Churn control (PR5)
	endpointMu         sync.RWMutex
	lastGoodEndpoints  map[string]EndpointIdentity // service name -> last successful endpoint
	lastDNSFailure     time.Time                   // timestamp of last DNS resolution failure
	dnsFailureCooldown time.Duration               // how long to use last-good on repeated failures

	// Observability (PR5)
	logLimiter         *logRateLimiter // rate limiter for routing logs
	snapshotRegenTotal uint64          // total snapshot regenerations
	snapshotNoopTotal  uint64          // snapshots skipped (no changes)
	lastMetricsLog     time.Time       // timestamp of last metrics log

	// Certificate rotation tracking (SDS)
	etcdClient            *clientv3.Client // etcd client for cert generation tracking
	lastCertGeneration    uint64           // DEPRECATED: use certReconciler
	certGenerationChecked bool             // DEPRECATED: use certReconciler

	// ACME certificate rotation tracking (file-based hash polling)
	acmeCertPath     string // Path to ACME fullchain.pem
	acmeKeyPath      string // Path to ACME privkey.pem
	lastACMECertHash string // DEPRECATED: use certReconciler
	lastACMEKeyHash  string // DEPRECATED: use certReconciler

	// v1 Conformance (INV-6): Event-driven certificate reconciliation
	certReconciler  *CertReconciler
	certInitPending bool
	certStarted     bool
}

// New creates a watcher bound to the given server.
// controllerAddr is optional - if provided, enables DNS-based service routing in cluster mode.
func New(logger *slog.Logger, srv *server.XDSServer, configPath, nodeID string, interval time.Duration, downstreamMode DownstreamTLSMode, controllerAddr string) *Watcher {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	if srv == nil {
		logger.Warn("xDS watcher created without server")
	}
	if strings.TrimSpace(nodeID) == "" {
		nodeID = defaultNodeID
	}
	if interval <= 0 {
		interval = defaultInterval
	}
	if downstreamMode == "" {
		downstreamMode = DownstreamTLSOptional
	}
	protocol := detectLocalProtocol()
	return &Watcher{
		logger:             logger,
		server:             srv,
		configPath:         strings.TrimSpace(configPath),
		nodeID:             nodeID,
		interval:           interval,
		downstreamTLSMode:  downstreamMode,
		protocol:           protocol,
		controllerAddr:     strings.TrimSpace(controllerAddr),
		dnsCache:           dnscache.New(30 * time.Second), // Default 30s TTL
		lastGoodEndpoints:  make(map[string]EndpointIdentity),
		dnsFailureCooldown: 30 * time.Second,                    // PR5: Reuse last-good for 30s on DNS failure
		logLimiter:         newLogRateLimiter(60 * time.Second), // PR5: Rate-limit logs to once per minute
	}
}

// SetEtcdClient configures the etcd client for certificate generation tracking.
// This enables hot certificate rotation via SDS.
// v1 Conformance (INV-6): Initializes CertReconciler for event-driven rotation.
func (w *Watcher) SetEtcdClient(client *clientv3.Client) {
	w.etcdClient = client
	if w.logger != nil && client != nil {
		w.logger.Info("etcd client configured for certificate rotation tracking")
	}

	// Initialize certificate reconciler if we have enough information
	w.initializeCertReconciler()
}

// SetACMEPaths configures ACME certificate paths for rotation tracking.
func (w *Watcher) SetACMEPaths(certPath, keyPath string) {
	w.acmeCertPath = certPath
	w.acmeKeyPath = keyPath

	// Initialize certificate reconciler if we have enough information
	w.initializeCertReconciler()
}

// initializeCertReconciler creates the CertReconciler if all required config is available.
func (w *Watcher) initializeCertReconciler() {
	// Need at least one certificate source to watch
	hasEtcd := w.etcdClient != nil
	hasACME := w.acmeCertPath != "" && w.acmeKeyPath != ""

	if !hasEtcd && !hasACME {
		return
	}

	// Determine domain for etcd key
	// VIOLATION INV-1.8: Domain-based key (future: use cert ID)
	domain := ""
	if w.clusterNetwork != nil && w.clusterNetwork.Spec != nil {
		domain = w.clusterNetwork.Spec.ClusterDomain
	}
	if domain == "" && w.controllerAddr != "" && !w.clusterNetworkReady() {
		// Wait until cluster network is fetched to avoid nil deref / incorrect domain
		w.certInitPending = true
		return
	}
	if domain == "" {
		domain = "globular.internal" // Hardcoded fallback
	}

	// Create reconciler if not already created
	if w.certReconciler == nil {
		w.certReconciler = NewCertReconciler(
			w.logger,
			w.etcdClient,
			domain,
			w.acmeCertPath,
			w.acmeKeyPath,
		)
		if w.logger != nil {
			w.logger.Info("certificate reconciler initialized",
				"has_etcd", hasEtcd,
				"has_acme", hasACME)
		}
	}

	w.certInitPending = false
}

// Run starts the watcher loop and blocks until the context is canceled.
// v1 Conformance (INV-6): Event-driven with certificate reconciliation.
func (w *Watcher) Run(ctx context.Context) error {
	if w.server == nil {
		return fmt.Errorf("server is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var certChangedChan, acmeChangedChan <-chan struct{}

	startReconciler := func() {
		if w.certReconciler != nil && !w.certStarted {
			if err := w.certReconciler.Start(); err != nil {
				w.logger.Warn("failed to start certificate reconciler", "err", err)
				return
			}
			w.certStarted = true
			w.logger.Info("certificate reconciler started")
			certChangedChan = w.certReconciler.CertChangedChan()
			acmeChangedChan = w.certReconciler.ACMEChangedChan()
		}
	}

	if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
		w.logger.Warn("initial xDS sync failed", "err", err)
	}
	if w.certInitPending && w.clusterNetworkReady() {
		w.initializeCertReconciler()
	}
	startReconciler()

	interval := w.interval
	if interval <= 0 {
		interval = defaultInterval
		w.interval = defaultInterval
	}
	timer := time.NewTimer(interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-timer.C:
			// Initialize/start reconciler once cluster network is ready
			if w.certInitPending && w.clusterNetworkReady() {
				w.initializeCertReconciler()
				startReconciler()
				if w.certReconciler != nil {
					certChangedChan = w.certReconciler.CertChangedChan()
					acmeChangedChan = w.certReconciler.ACMEChangedChan()
				}
			}

			if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
				w.logger.Warn("xDS sync failed", "err", err)
			}
			next := w.interval
			if next <= 0 {
				next = defaultInterval
			}
			timer.Reset(next)

		case <-certChangedChan:
			// v1 Conformance (INV-6.1): Event-driven internal certificate rotation
			w.logger.Info("internal certificate changed - rebuilding snapshot")
			if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
				w.logger.Warn("xDS sync after cert change failed", "err", err)
			}

		case <-acmeChangedChan:
			// v1 Conformance (INV-6.2): Event-driven ACME certificate rotation
			w.logger.Info("ACME certificate changed - rebuilding snapshot")
			if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
				w.logger.Warn("xDS sync after ACME change failed", "err", err)
			}
		}
	}
}

func (w *Watcher) sync(ctx context.Context) error {
	// Fetch cluster network config if controller address is configured (PR4)
	if w.controllerAddr != "" && w.clusterNetwork == nil {
		if err := w.fetchClusterNetwork(ctx); err != nil {
			w.logger.Warn("failed to fetch cluster network config", "err", err, "controller_addr", w.controllerAddr)
			// Continue without cluster config - fall back to IP-based routing
		} else if w.clusterNetwork != nil && w.clusterNetwork.Spec != nil {
			ttl := time.Duration(w.clusterNetwork.Spec.DnsTtl) * time.Second
			if ttl <= 0 {
				ttl = 30 * time.Second
			}

			// PR7: Configure DNS cache with nameservers for high availability
			nameservers := w.clusterNetwork.Spec.DnsNameservers
			if len(nameservers) > 0 {
				w.dnsCache = dnscache.New(ttl, nameservers...)
				w.logger.Info("DNS cache configured", "ttl", ttl, "cluster_domain", w.clusterNetwork.Spec.ClusterDomain, "nameservers", nameservers)
			} else if ttl != 30*time.Second {
				w.dnsCache = dnscache.New(ttl)
				w.logger.Info("DNS cache configured", "ttl", ttl, "cluster_domain", w.clusterNetwork.Spec.ClusterDomain)
			}
		}
	}

	prevProtocol := w.protocol
	w.protocol = detectLocalProtocol()
	if w.logger != nil && prevProtocol != "" && prevProtocol != w.protocol {
		w.logger.Info("protocol changed", "from", prevProtocol, "to", w.protocol)
	}

	// v1 Conformance (INV-6): Certificate rotation now handled by CertReconciler
	// Events are delivered via channels in Run loop (event-driven, not polling)
	// Old polling functions (checkCertificateGeneration, checkACMECertificateRotation)
	// are deprecated and only called if reconciler is not configured (backward compat)
	if w.certReconciler == nil {
		// Fallback to polling if reconciler not configured
		certChanged := w.checkCertificateGeneration(ctx)
		if certChanged && w.logger != nil {
			w.logger.Info("certificate rotation detected - forcing snapshot update")
		}

		acmeChanged := w.checkACMECertificateRotation(w.acmeCertPath, w.acmeKeyPath)
		if acmeChanged && w.logger != nil {
			w.logger.Info("ACME certificate rotation detected - forcing snapshot update")
		}
	}

	input, version, err := w.buildInput(ctx)
	if err != nil {
		return err
	}

	if strings.TrimSpace(input.NodeID) == "" {
		input.NodeID = w.nodeID
	}
	if version == "" {
		version = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	input.Version = version

	snap, err := builder.BuildSnapshot(input, version)
	if err != nil {
		return fmt.Errorf("build snapshot: %w", err)
	}
	if err := w.server.SetSnapshot(input.NodeID, snap); err != nil {
		return fmt.Errorf("push snapshot: %w", err)
	}
	atomic.AddUint64(&w.snapshotRegenTotal, 1) // PR5: Track snapshot regen
	w.logger.Info("xDS snapshot pushed", "node_id", input.NodeID, "version", version)

	// PR5: Log metrics periodically (every 60s)
	w.logMetricsPeriodically()

	return nil
}

// logMetricsPeriodically logs DNS cache and snapshot metrics once per minute (PR5)
func (w *Watcher) logMetricsPeriodically() {
	now := time.Now()
	if now.Sub(w.lastMetricsLog) < 60*time.Second {
		return // Not time yet
	}
	w.lastMetricsLog = now

	if w.dnsCache == nil {
		return // No cache to report
	}

	cacheStats := w.dnsCache.Stats()
	snapshotRegen := atomic.LoadUint64(&w.snapshotRegenTotal)
	snapshotNoop := atomic.LoadUint64(&w.snapshotNoopTotal)

	w.logger.Info("metrics summary",
		"dnscache_a_hit", cacheStats.AHit,
		"dnscache_a_miss", cacheStats.AMiss,
		"dnscache_srv_hit", cacheStats.SRVHit,
		"dnscache_srv_miss", cacheStats.SRVMiss,
		"xds_snapshot_regen_total", snapshotRegen,
		"xds_snapshot_noop_total", snapshotNoop)
}

// fetchClusterNetwork retrieves cluster network configuration from the controller (PR4)
func (w *Watcher) fetchClusterNetwork(ctx context.Context) error {
	if w.controllerAddr == "" {
		return fmt.Errorf("controller address not configured")
	}

	client := controllerclient.New(w.controllerAddr)
	clusterNet, err := client.GetClusterNetwork(ctx)
	if err != nil {
		return fmt.Errorf("get cluster network: %w", err)
	}

	w.clusterNetwork = clusterNet
	w.logger.Info("fetched cluster network config",
		"cluster_domain", clusterNet.Spec.GetClusterDomain(),
		"gateway_fqdn", clusterNet.Spec.GetGatewayFqdn())
	return nil
}

func (w *Watcher) buildInput(ctx context.Context) (builder.Input, string, error) {
	var cfg *XDSConfig
	if w.configPath != "" {
		fi, err := os.Stat(w.configPath)
		if err != nil {
			if os.IsNotExist(err) {
				return w.buildDynamicInput(ctx, nil)
			}
			return builder.Input{}, "", err
		}
		cfg, err = w.loadXDSConfig(fi)
		if err != nil {
			return builder.Input{}, "", err
		}
	}

	return w.buildDynamicInput(ctx, cfg)
}

func (w *Watcher) loadXDSConfig(fi os.FileInfo) (*XDSConfig, error) {
	if w.configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}
	if w.configCached != nil && fi.ModTime().Equal(w.configMod) {
		return w.configCached, nil
	}
	data, err := os.ReadFile(w.configPath)
	if err != nil {
		return nil, err
	}
	var cfg XDSConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.normalize()
	if cfg.SyncIntervalSeconds > 0 {
		w.interval = time.Duration(cfg.SyncIntervalSeconds) * time.Second
	}
	w.configCached = &cfg
	w.configMod = fi.ModTime()
	return &cfg, nil
}

// setupSDSSecrets configures SDS (Secret Discovery Service) for TLS certificate delivery.
// It builds the secret resources from the listener's certificate configuration and validates
// that SDS is only used with TLS-secured xDS control plane connections.
//
// Returns:
//   - enableSDS: whether SDS should be enabled
//   - secrets: list of SDS secrets to provide to Envoy
//   - error: validation error if security constraints are violated
func (w *Watcher) setupSDSSecrets(listener builder.Listener) (bool, []builder.Secret, error) {
	// SDS requires certificate files to be configured
	if listener.CertFile == "" || listener.KeyFile == "" {
		return false, nil, nil
	}

	// Build SDS secrets using canonical TLS paths
	sdsSecrets := []builder.Secret{
		{
			Name:     secrets.InternalServerCert,
			CertPath: listener.CertFile,
			KeyPath:  listener.KeyFile,
		},
	}

	// Add CA bundle if present
	if listener.IssuerFile != "" {
		sdsSecrets = append(sdsSecrets, builder.Secret{
			Name:   secrets.InternalCABundle,
			CAPath: listener.IssuerFile,
		})
	}

	// Security invariant: SDS implies TLS-secured xDS
	// Secrets contain private keys and MUST NOT traverse plaintext gRPC
	if os.Getenv("GLOBULAR_XDS_INSECURE") == "1" {
		return false, nil, fmt.Errorf(
			"security violation: SDS enabled but xDS running insecure (GLOBULAR_XDS_INSECURE=1) - " +
				"secrets would be transmitted over plaintext")
	}

	// Add public ACME certificate if available (for public HTTPS ingress)
	// ACME certificates use fullchain.pem naming convention
	if strings.Contains(listener.CertFile, "fullchain.pem") {
		sdsSecrets = append(sdsSecrets, builder.Secret{
			Name:     secrets.PublicIngressCert,
			CertPath: listener.CertFile,
			KeyPath:  listener.KeyFile,
		})
		// Store ACME cert paths for rotation detection
		w.acmeCertPath = listener.CertFile
		w.acmeKeyPath = listener.KeyFile
	}

	if w.logger != nil {
		w.logger.Info("SDS enabled", "secrets", len(sdsSecrets))
	}

	return true, sdsSecrets, nil
}

// modeResources holds the resources and configuration produced by mode selection (ingress vs legacy).
type modeResources struct {
	listener           builder.Listener
	clusters           []builder.Cluster
	routes             []builder.Route
	gatewayPort        uint32
	ingressHTTPPort    uint32
	enableHTTPRedirect bool
}

// configureModeResources selects between ingress mode and legacy mode, applies the appropriate
// configuration, and returns normalized resources ready for snapshot building.
//
// Mode selection:
//   - If ingressSpec is non-nil: Use ingress mode (modern configuration)
//   - If ingressSpec is nil: Use legacy mode (fallback for old configs)
//
// The function applies defaults and normalization to ensure the listener has valid
// Host, RouteName, and Name fields regardless of which mode is selected.
func (w *Watcher) configureModeResources(ingressSpec *IngressSpec, cfg *XDSConfig) (modeResources, error) {
	var result modeResources

	if ingressSpec != nil {
		// Ingress mode: Use modern ingress configuration
		result.listener = ingressSpec.Listener
		result.clusters = ingressSpec.Clusters
		result.routes = ingressSpec.Routes

		// Normalize gateway port
		normalizedGatewayPort := normalizeIngressGateway(ingressSpec, w, cfg)
		ingressSpec.GatewayPort = normalizedGatewayPort
		result.gatewayPort = ingressSpec.GatewayPort
		result.ingressHTTPPort = ingressSpec.HTTPPort
		result.enableHTTPRedirect = ingressSpec.EnableHTTPRedirect

		// Apply listener defaults for ingress mode
		if result.listener.Host == "" {
			result.listener.Host = "0.0.0.0"
		}
		if result.listener.RouteName == "" {
			result.listener.RouteName = defaultRouteName
		}
		if result.listener.Name == "" {
			port := result.listener.Port
			if port == 0 {
				port = controlplane.DefaultIngressPort(result.listener.Host)
			}
			result.listener.Name = fmt.Sprintf("%s_%d", listenerNameBase, port)
		}
	} else {
		// Legacy mode: Use legacy gateway resources (fallback)
		legacyClusters, legacyRoutes, legacyListener, legacyGatewayPort, legacyTLSEnabled, err := w.buildLegacyGatewayResources(cfg)
		if err != nil {
			return modeResources{}, err
		}

		result.listener = legacyListener
		result.clusters = legacyClusters
		result.routes = legacyRoutes
		result.gatewayPort = legacyGatewayPort

		// Configure HTTP redirect if TLS is enabled in legacy mode
		if cfg != nil && legacyTLSEnabled && cfg.ingressRedirectEnabled() {
			result.enableHTTPRedirect = cfg.ingressRedirectEnabled()
			result.ingressHTTPPort = cfg.Ingress.HTTPPort
		}
	}

	// Final normalization: Ensure RouteName is set regardless of mode
	if result.listener.RouteName == "" {
		result.listener.RouteName = defaultRouteName
	}

	return result, nil
}

func (w *Watcher) buildDynamicInput(ctx context.Context, cfg *XDSConfig) (builder.Input, string, error) {
	clusters, routes, err := w.buildServiceResources(ctx, cfg)
	if err != nil {
		return builder.Input{}, "", err
	}

	ingressSpec, err := w.buildIngressSpec(ctx, cfg)
	if err != nil {
		return builder.Input{}, "", err
	}
	// Log mode selection
	if w.logger != nil {
		if ingressSpec != nil {
			hasGatewayHTTP := findClusterByName(ingressSpec.Clusters, "gateway_http") != nil
			hasGatewayHTTPS := findClusterByName(ingressSpec.Clusters, "globular_https") != nil
			w.logger.Info("ingress path active",
				"https_port", ingressSpec.Listener.Port,
				"http_port", ingressSpec.HTTPPort,
				"redirect", ingressSpec.EnableHTTPRedirect,
				"has_gateway_http", hasGatewayHTTP,
				"has_globular_https", hasGatewayHTTPS,
			)
		} else {
			w.logger.Info("legacy path active")
		}
	}

	// Configure mode-specific resources (ingress vs legacy)
	modeRes, err := w.configureModeResources(ingressSpec, cfg)
	if err != nil {
		return builder.Input{}, "", err
	}

	// Merge mode-specific resources with service resources
	listener := modeRes.listener
	clusters = append(clusters, modeRes.clusters...)
	routes = append(routes, modeRes.routes...)

	version := fmt.Sprintf("%d", time.Now().UnixNano())
	if w.logger != nil {
		if gw := findClusterByName(clusters, "gateway_http"); gw != nil {
			w.logger.Info("xDS clusters resolved", "gateway_http", gw.Endpoints)
		} else {
			w.logger.Info("xDS clusters resolved", "gateway_http", "missing")
		}
	}

	// Configure SDS for TLS certificate delivery (hot rotation)
	enableSDS, sdsSecrets, err := w.setupSDSSecrets(listener)
	if err != nil {
		return builder.Input{}, "", err
	}

	input := builder.Input{
		NodeID:             w.nodeID,
		Listener:           listener,
		Routes:             routes,
		Clusters:           clusters,
		IngressHTTPPort:    modeRes.ingressHTTPPort,
		EnableHTTPRedirect: modeRes.enableHTTPRedirect,
		GatewayPort:        modeRes.gatewayPort,
		Version:            version,
		EnableSDS:          enableSDS,
		SDSSecrets:         sdsSecrets,
	}
	return input, version, nil
}

func (w *Watcher) buildServiceResources(ctx context.Context, cfg *XDSConfig) ([]builder.Cluster, []builder.Route, error) {
	var services []map[string]any
	var err error
	if cfg != nil && len(cfg.EtcdEndpoints) > 0 {
		services, err = w.buildServiceResourcesFromEtcd(ctx, cfg)
		if err != nil && w.logger != nil {
			w.logger.Warn("etcd service discovery failed; falling back to local config", "err", err)
		}
	}
	if services == nil {
		services, err = config.GetServicesConfigurations()
		if err != nil {
			return nil, nil, err
		}
	}

	// Service clusters use CA-only validation (no client certificates)
	_, _, downIssuer, err := w.downstreamTLSConfig()
	if err != nil {
		return nil, nil, err
	}

	var (
		clusters        []builder.Cluster
		routes          []builder.Route
		addedClusters   = map[string]struct{}{}
		allClusterNames []string
	)

	for _, svc := range services {
		name := strings.TrimSpace(fmt.Sprint(svc["Name"]))
		endpoint := w.resolveServiceEndpoint(ctx, svc)

		if w.logger != nil && strings.TrimSpace(fmt.Sprint(svc["Address"])) == "" {
			w.logger.Debug("service Address empty; defaulting to local endpoint",
				"service", name,
				"host", endpoint.Host(),
				"port", endpoint.Port,
				"source", endpoint.Source.String())
		}

		if name == "" || endpoint.Host() == "" || endpoint.Port == 0 {
			continue
		}

		clusterName := strings.ReplaceAll(name, ".", "_") + "_cluster"
		if _, ok := addedClusters[clusterName]; ok {
			continue
		}
		addedClusters[clusterName] = struct{}{}

		// PR5: Use FQDN for endpoint identity when available
		endpointHost := endpoint.Host()

		// Service clusters use CA-only validation (no client certificates)
		// Client certificates cause KEY_VALUES_MISMATCH errors in Envoy
		clusters = append(clusters, builder.Cluster{
			Name:      clusterName,
			Endpoints: []builder.Endpoint{{Host: endpointHost, Port: uint32(endpoint.Port)}},
			CAFile:    downIssuer,
			SNI:       endpointHost,
		})
		allClusterNames = append(allClusterNames, clusterName)
		routes = append(routes, builder.Route{Prefix: "/" + name + "/", Cluster: clusterName})
	}

	healthTarget := ""
	for _, cn := range allClusterNames {
		if strings.HasPrefix(cn, "authentication_AuthenticationService_") {
			healthTarget = cn
			break
		}
	}
	if healthTarget == "" && len(allClusterNames) > 0 {
		healthTarget = allClusterNames[0]
	}
	if healthTarget != "" {
		routes = append([]builder.Route{{Prefix: "/grpc.health.v1.Health/", Cluster: healthTarget}}, routes...)
	}

	return clusters, routes, nil
}

func (w *Watcher) buildLegacyGatewayResources(xdsCfg *XDSConfig) ([]builder.Cluster, []builder.Route, builder.Listener, uint32, bool, error) {
	cfg, _ := config.GetLocalConfig(true)
	host, listenPort := readGatewayAddress()
	listenPort = defaultGatewayPort(listenPort)
	portHTTP, portHTTPS := readGatewayPortsFromConfig(cfg)
	upstreamHost := normalizeUpstreamHost(host)
	gatewayCert, gatewayKey, gatewayCA, ok := w.gatewayTLSPaths()
	tlsEnabled := strings.ToLower(strings.TrimSpace(w.protocol)) == "https" && ok
	if !tlsEnabled {
		gatewayCert, gatewayKey, gatewayCA = "", "", ""
	}

	// CRITICAL: Gateway cluster always uses portHTTP (8080) for TLS termination at Envoy
	// Envoy terminates TLS on the listener, then connects to gateway via plain HTTP
	// Using portHTTPS would create a routing loop (Envoy → itself on 8443)
	gatewayCluster := "gateway_http"
	upstreamPort := portHTTP
	if listenPort != defaultGatewayPort(0) {
		upstreamPort = listenPort
	}

	// Gateway cluster uses plain HTTP (no TLS config) for TLS termination architecture
	clusters := []builder.Cluster{{
		Name:      gatewayCluster,
		Endpoints: []builder.Endpoint{{Host: upstreamHost, Port: uint32(upstreamPort)}},
		// No CAFile, ServerCert, KeyFile, or SNI - plain HTTP upstream
	}}
	routes := []builder.Route{{Prefix: "/", Cluster: gatewayCluster}}

	listenerPort := uint32(portHTTPS)
	if !tlsEnabled {
		listenerPort = uint32(portHTTP)
	}

	listener := builder.Listener{
		Name:       fmt.Sprintf("%s_%d", listenerNameBase, listenerPort),
		RouteName:  defaultRouteName,
		Host:       "0.0.0.0",
		Port:       listenerPort,
		CertFile:   gatewayCert,
		KeyFile:    gatewayKey,
		IssuerFile: gatewayCA,
	}
	if !tlsEnabled {
		listener.CertFile, listener.KeyFile, listener.IssuerFile = "", "", ""
	}
	if w.logger != nil {
		w.logger.Info("legacy gateway resolved",
			"protocol", w.protocol,
			"tls_enabled", tlsEnabled,
			"port_http", portHTTP,
			"port_https", portHTTPS,
			"gateway_listen_override", listenPort != defaultGatewayPort(0),
			"cluster", gatewayCluster,
			"upstream_port", upstreamPort,
		)
	}
	return clusters, routes, listener, uint32(upstreamPort), tlsEnabled, nil
}

func (w *Watcher) buildIngressSpec(ctx context.Context, cfg *XDSConfig) (*IngressSpec, error) {
	if cfg == nil {
		return nil, nil
	}
	if len(cfg.EtcdEndpoints) > 0 {
		spec, err := w.ingressFromEtcd(ctx, cfg)
		if err != nil {
			if w.logger != nil {
				w.logger.Warn("etcd ingress lookup failed", "err", err)
			}
		} else if spec != nil {
			w.applyIngressSettings(spec, cfg)
			return spec, nil
		}
	}
	if cfg.Fallback != nil && cfg.Fallback.Enabled {
		if spec := w.ingressFromFallback(cfg.Fallback); spec != nil {
			if w.logger != nil {
				w.logger.Info("using fallback ingress configuration", "routes", len(spec.Routes), "clusters", len(spec.Clusters))
			}
			w.applyIngressSettings(spec, cfg)
			return spec, nil
		}
	}
	return nil, nil
}

func (w *Watcher) ingressFromEtcd(ctx context.Context, cfg *XDSConfig) (*IngressSpec, error) {
	if len(cfg.EtcdEndpoints) == 0 {
		return nil, nil
	}

	etcdCfg := clientv3.Config{
		Endpoints:   cfg.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	}

	// TLS is mandatory for etcd connections
	tlsCfg, err := config.GetEtcdTLS()
	if err != nil {
		return nil, fmt.Errorf("etcd TLS configuration required: %w", err)
	}
	etcdCfg.TLS = tlsCfg

	cli, err := clientv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	spec, err := parseEtcdIngress(etcdCtx, clientv3.NewKV(cli))
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func (w *Watcher) ingressFromFallback(fb *FallbackConfig) *IngressSpec {
	if fb == nil || fb.Ingress == nil || !fb.Enabled {
		return nil
	}
	routes := make([]builder.Route, 0, len(fb.Ingress.Routes))
	for _, fr := range fb.Ingress.Routes {
		prefix := strings.TrimSpace(fr.Prefix)
		cluster := strings.TrimSpace(fr.Cluster)
		if prefix == "" || cluster == "" {
			continue
		}
		route := builder.Route{
			Prefix:      prefix,
			Cluster:     cluster,
			HostRewrite: strings.TrimSpace(fr.HostRewrite),
			Authority:   strings.TrimSpace(fr.Authority),
		}
		if dom := strings.TrimSpace(fr.Domains); dom != "" {
			route.Domains = parseDomains(dom)
		}
		routes = append(routes, route)
	}
	if len(routes) == 0 || len(fb.Clusters) == 0 {
		return nil
	}
	// Day-0 Security: Get CA for upstream TLS to internal services
	_, _, caFile, _ := w.downstreamTLSConfig()

	clusters := make([]builder.Cluster, 0, len(fb.Clusters))
	for _, fc := range fb.Clusters {
		name := strings.TrimSpace(fc.Name)
		if name == "" {
			continue
		}
		endpoints := make([]builder.Endpoint, 0, len(fc.Endpoints))
		for _, ep := range fc.Endpoints {
			if ep.Host == "" || ep.Port == 0 {
				continue
			}
			endpoints = append(endpoints, builder.Endpoint{Host: ep.Host, Port: ep.Port, Priority: ep.Priority})
		}
		if len(endpoints) == 0 {
			continue
		}
		// Day-0 Security: Enable upstream TLS with cluster CA when available
		cluster := builder.Cluster{
			Name:      name,
			Endpoints: endpoints,
		}
		if caFile != "" {
			cluster.CAFile = caFile
			// Set SNI to first endpoint hostname if not an IP
			if len(endpoints) > 0 {
				cluster.SNI = endpoints[0].Host
			}
		}
		clusters = append(clusters, cluster)
	}
	if len(clusters) == 0 {
		return nil
	}
	return &IngressSpec{
		Listener: builder.Listener{
			Host:       strings.TrimSpace(fb.Ingress.ListenerHost),
			Port:       fb.Ingress.HTTPSPort,
			RouteName:  defaultRouteName,
			CertFile:   strings.TrimSpace(fb.Ingress.TLS.CertFile),
			KeyFile:    strings.TrimSpace(fb.Ingress.TLS.KeyFile),
			IssuerFile: strings.TrimSpace(fb.Ingress.TLS.IssuerFile),
		},
		Routes:             routes,
		Clusters:           clusters,
		HTTPPort:           fb.Ingress.HTTPPort,
		EnableHTTPRedirect: boolValue(fb.Ingress.EnableHTTPRedirect, true),
		RedirectConfigured: fb.Ingress.EnableHTTPRedirect != nil,
	}
}

func (w *Watcher) applyIngressSettings(spec *IngressSpec, cfg *XDSConfig) {
	if spec == nil || cfg == nil {
		return
	}
	// v1 Conformance (INV-5.1): Use IngressDomains for virtual host matching, not cluster_domain
	// cluster_domain is for internal DNS only (service discovery, internal FQDNs)
	// ingress_domains are for public routing (virtual host matching, certificate SANs)
	if spec.Listener.Host == "" {
		spec.Listener.Host = "0.0.0.0"
	}
	if spec.Listener.Port == 0 {
		spec.Listener.Port = cfg.Ingress.HTTPSPort
	}
	if spec.HTTPPort == 0 {
		spec.HTTPPort = cfg.Ingress.HTTPPort
	}
	if spec.HTTPPort == cfg.gatewayPort() {
		w.warnIngressPortCollision(spec.HTTPPort)
		spec.HTTPPort = controlplane.DefaultIngressHTTPPort(spec.Listener.Host)
	}
	if !spec.RedirectConfigured {
		spec.EnableHTTPRedirect = cfg.ingressRedirectEnabled()
	}
	spec.GatewayPort = cfg.gatewayPort()
	// Assign ingress domains to routes that don't have explicit domains
	if len(cfg.IngressDomains) > 0 {
		for i := range spec.Routes {
			if len(spec.Routes[i].Domains) == 0 {
				// Normalize domains: lowercase and trim whitespace
				normalizedDomains := make([]string, len(cfg.IngressDomains))
				for j, d := range cfg.IngressDomains {
					normalizedDomains[j] = strings.ToLower(strings.TrimSpace(d))
				}
				spec.Routes[i].Domains = normalizedDomains
			}
		}
	}
	w.applyIngressTLS(spec, cfg)
}

func (w *Watcher) applyIngressTLS(spec *IngressSpec, cfg *XDSConfig) {
	if spec == nil || cfg == nil {
		return
	}
	tlsConfig := cfg.Ingress
	if !tlsConfig.tlsEnabled() {
		w.disableIngressTLS(spec)
		return
	}
	certPath := strings.TrimSpace(spec.Listener.CertFile)
	keyPath := strings.TrimSpace(spec.Listener.KeyFile)
	if certPath == "" {
		certPath = tlsConfig.TLS.CertChainPath
	}
	if keyPath == "" {
		keyPath = tlsConfig.TLS.PrivateKeyPath
	}
	cert := pathIfExists(certPath)
	key := pathIfExists(keyPath)
	if cert == "" || key == "" {
		w.warnIngressTLSMissing(certPath, keyPath)
		w.disableIngressTLS(spec)
		return
	}
	spec.Listener.CertFile = cert
	spec.Listener.KeyFile = key
	if tlsConfig.MTLS.Enabled {
		ca := pathIfExists(tlsConfig.MTLS.CAPath)
		if ca == "" {
			w.warnIngressMTLSMissing(tlsConfig.MTLS.CAPath)
			w.disableIngressTLS(spec)
			return
		}
		spec.Listener.IssuerFile = ca
	} else {
		spec.Listener.IssuerFile = ""
	}
}

func (w *Watcher) disableIngressTLS(spec *IngressSpec) {
	spec.Listener.CertFile = ""
	spec.Listener.KeyFile = ""
	spec.Listener.IssuerFile = ""
	spec.EnableHTTPRedirect = false
}

func (w *Watcher) warnIngressTLSMissing(cert, key string) {
	if w.logger == nil || w.ingressTLSWarned {
		return
	}
	w.logger.Warn("ingress TLS enabled but certificate/key missing", "cert", cert, "key", key)
	w.ingressTLSWarned = true
}

func (w *Watcher) warnIngressMTLSMissing(caPath string) {
	if w.logger == nil || w.ingressMTLSWarned {
		return
	}
	w.logger.Warn("ingress mTLS enabled but CA certificate missing", "ca", caPath)
	w.ingressMTLSWarned = true
}

func (w *Watcher) warnIngressPortCollision(port uint32) {
	if w.logger == nil || w.ingressPortCollisionWarned {
		return
	}
	w.logger.Warn("ingress HTTP port collides with gateway listen port; forcing default 80", "port", port)
	w.ingressPortCollisionWarned = true
}

func (w *Watcher) downstreamTLSConfig() (string, string, string, error) {
	certPath := config.GetLocalClientCertificatePath()
	keyPath := config.GetLocalClientKeyPath()
	switch w.downstreamTLSMode {
	case DownstreamTLSDisabled:
		return "", "", "", nil
	case DownstreamTLSRequired:
		if pathIfExists(certPath) == "" || pathIfExists(keyPath) == "" {
			if w.logger != nil && !w.downstreamTLSRequiredWarned {
				w.logger.Warn("downstream TLS required but certificate/key missing; serving without TLS",
					"cert_path", certPath, "key_path", keyPath)
				w.downstreamTLSRequiredWarned = true
			}
			return "", "", "", nil
		}
		downCert, err := require("downstream certificate", certPath)
		if err != nil {
			return "", "", "", err
		}
		downKey, err := require("downstream private key", keyPath)
		if err != nil {
			return "", "", "", err
		}
		return downCert, downKey, pathIfExists(config.GetLocalCACertificate()), nil
	default:
		downCert := pathIfExists(certPath)
		downKey := pathIfExists(keyPath)
		if downCert == "" || downKey == "" {
			if !w.downstreamTLSWarned && w.logger != nil {
				w.logger.Warn("downstream TLS optional but certificate/key missing; serving without TLS",
					"cert_path", certPath, "key_path", keyPath)
				w.downstreamTLSWarned = true
			}
			return "", "", "", nil
		}
		return downCert, downKey, pathIfExists(config.GetLocalCACertificate()), nil
	}
}

func require(label, path string) (string, error) {
	p := strings.TrimSpace(fmt.Sprint(path))
	if p != "" && Utility.Exists(p) {
		return p, nil
	}
	return "", fmt.Errorf("%s path missing or empty: %q", label, path)
}

func pathIfExists(path string) string {
	p := strings.TrimSpace(fmt.Sprint(path))
	if p != "" && Utility.Exists(p) {
		return p
	}
	return ""
}

func parseAddress(address string) (string, int) {
	parts := strings.Split(strings.TrimSpace(address), ":")
	if len(parts) == 0 {
		return "127.0.0.1", 0
	}
	host := parts[0]
	port := 0
	if len(parts) > 1 {
		port = Utility.ToInt(parts[1])
	}
	return host, port
}

func parseGatewayListenValue(raw any) (string, int, bool) {
	value := strings.TrimSpace(Utility.ToString(raw))
	if value == "" {
		return "", 0, false
	}
	host, port := parseAddress(value)
	if port == 0 {
		return "", 0, false
	}
	return host, port, true
}

func toStringAnyMap(raw any) (map[string]any, bool) {
	switch m := raw.(type) {
	case map[string]any:
		return m, true
	case map[interface{}]interface{}:
		result := make(map[string]any, len(m))
		for key, value := range m {
			if keyStr, ok := key.(string); ok {
				result[keyStr] = value
			}
		}
		return result, true
	default:
		return nil, false
	}
}

func readGatewayAddressFrom(cfg map[string]any) (string, int) {
	if cfg == nil {
		return "", 0
	}
	for _, gatewayKey := range []string{"gateway", "Gateway"} {
		if raw, ok := cfg[gatewayKey]; ok {
			if host, port, ok := parseGatewayListenValue(raw); ok {
				return host, port
			}
			if gwMap, ok := toStringAnyMap(raw); ok {
				for _, listenKey := range []string{"listen", "Listen"} {
					if listenRaw, ok := gwMap[listenKey]; ok {
						if host, port, ok := parseGatewayListenValue(listenRaw); ok {
							return host, port
						}
					}
				}
			}
		}
	}
	return "", 0
}

func readGatewayAddressFromConfig(cfg map[string]any) (string, int) {
	host := "127.0.0.1"
	port := 0
	if cfg == nil {
		return host, defaultGatewayPort(port)
	}
	if cfgHost, cfgPort := readGatewayAddressFrom(cfg); cfgHost != "" {
		host = normalizeUpstreamHost(cfgHost)
		if cfgPort != 0 {
			port = cfgPort
		}
	}
	return host, defaultGatewayPort(port)
}

func readGatewayAddress() (string, int) {
	cfg, err := config.GetLocalConfig(true)
	if err != nil || cfg == nil {
		return readGatewayAddressFromConfig(nil)
	}
	return readGatewayAddressFromConfig(cfg)
}

// v1 Conformance (INV-5.4): Removed domain return value
// Domain should not be used for service discovery or SNI - use ClusterDomain from XDSConfig instead
func readGatewayPortsFromConfig(cfg map[string]any) (int, int) {
	portHTTP := defaultGatewayPort(Utility.ToInt(cfg["PortHTTP"]))
	portHTTPS := Utility.ToInt(cfg["PortHTTPS"])
	if portHTTPS == 0 {
		portHTTPS = 8443
	}
	return portHTTP, portHTTPS
}

func defaultGatewayPort(port int) int {
	if port == 0 {
		return 8080
	}
	return port
}

func normalizeGatewayHTTPCluster(spec *IngressSpec, gatewayPort uint32, logger *slog.Logger) {
	if spec == nil || gatewayPort == 0 {
		return
	}
	for ci := range spec.Clusters {
		if strings.TrimSpace(spec.Clusters[ci].Name) != "gateway_http" {
			continue
		}
		for ei := range spec.Clusters[ci].Endpoints {
			endpoint := &spec.Clusters[ci].Endpoints[ei]
			if shouldNormalizeGatewayEndpoint(endpoint.Host) {
				if endpoint.Port == 0 {
					endpoint.Port = gatewayPort
				} else if endpoint.Port != gatewayPort && logger != nil {
					logger.Info("gateway_http endpoint preserved from configuration", "host", endpoint.Host, "port", endpoint.Port, "expected_port", gatewayPort)
				}
			}
		}
	}
}

func shouldNormalizeGatewayEndpoint(host string) bool {
	switch strings.ToLower(strings.TrimSpace(host)) {
	case "", "0.0.0.0", "127.0.0.1", "localhost", "*":
		return true
	}
	return Utility.IsLocal(host)
}

func findClusterByName(clusters []builder.Cluster, name string) *builder.Cluster {
	for i := range clusters {
		if strings.TrimSpace(clusters[i].Name) == name {
			return &clusters[i]
		}
	}
	return nil
}

func normalizeIngressGateway(spec *IngressSpec, w *Watcher, xdsCfg *XDSConfig) uint32 {
	if spec == nil || w == nil {
		return 0
	}
	cfg, _ := config.GetLocalConfig(true)
	portHTTP, _ := readGatewayPortsFromConfig(cfg)
	host, listenPort := readGatewayAddress()
	upstreamHost := normalizeUpstreamHost(host)
	if upstreamHost == "" {
		upstreamHost = "127.0.0.1"
	}
	// CRITICAL: Gateway cluster always uses portHTTP (8080) for TLS termination at Envoy
	// Envoy terminates TLS on the listener, then connects to gateway via plain HTTP
	// Using portHTTPS would create a routing loop (Envoy → itself on 8443)
	targetCluster := "gateway_http"
	targetPort := portHTTP
	if listenPort != defaultGatewayPort(0) {
		targetPort = listenPort
	}

	ensureCluster := func(name string, port uint32) *builder.Cluster {
		cl := findClusterByName(spec.Clusters, name)
		if cl == nil {
			spec.Clusters = append(spec.Clusters, builder.Cluster{Name: name})
			cl = &spec.Clusters[len(spec.Clusters)-1]
		}
		if len(cl.Endpoints) == 0 {
			cl.Endpoints = []builder.Endpoint{{Host: upstreamHost, Port: port}}
		} else {
			for i := range cl.Endpoints {
				if shouldNormalizeGatewayEndpoint(cl.Endpoints[i].Host) {
					cl.Endpoints[i].Port = port
					if cl.Endpoints[i].Host == "" {
						cl.Endpoints[i].Host = upstreamHost
					}
				}
			}
		}
		// Gateway cluster uses plain HTTP - no TLS config
		cl.ServerCert, cl.KeyFile, cl.CAFile, cl.SNI = "", "", "", ""
		return cl
	}

	// Always use gateway_http cluster (plain HTTP to gateway on port 8080)
	ensureCluster(targetCluster, uint32(targetPort))

	rootUpdated := false
	for i := range spec.Routes {
		if strings.TrimSpace(spec.Routes[i].Prefix) == "/" {
			spec.Routes[i].Cluster = targetCluster
			rootUpdated = true
		}
	}
	if !rootUpdated {
		spec.Routes = append([]builder.Route{{Prefix: "/", Cluster: targetCluster}}, spec.Routes...)
	}

	return uint32(targetPort)
}

// resolveExplicitHost handles explicit host/IP addresses from service registry.
// Returns zero value if no explicit host is provided.
func (w *Watcher) resolveExplicitHost(host, serviceDNSLabel string, port int) (EndpointIdentity, bool) {
	if host == "" {
		return EndpointIdentity{}, false
	}

	// Check if host looks like an FQDN (contains cluster domain)
	if w.isClusterMode() && strings.Contains(host, w.clusterNetwork.Spec.ClusterDomain) {
		return EndpointIdentity{
			ServiceDNSLabel: serviceDNSLabel,
			TargetFQDN:      host,
			TargetIP:        "",
			Port:            port,
			Source:          EndpointSourceRegistry,
		}, true
	}

	// Treat as IP address
	return EndpointIdentity{
		ServiceDNSLabel: serviceDNSLabel,
		TargetFQDN:      "",
		TargetIP:        host,
		Port:            port,
		Source:          EndpointSourceRegistry,
	}, true
}

// trySRVResolution attempts DNS SRV record lookup for service discovery.
// Returns endpoint and true if successful, zero value and false otherwise.
func (w *Watcher) trySRVResolution(ctx context.Context, serviceName, serviceDNSLabel string) (EndpointIdentity, bool) {
	if serviceName == "" {
		return EndpointIdentity{}, false
	}

	srvFQDN, srvPort := w.resolveSRV(ctx, serviceName)
	if srvFQDN == "" || srvPort == 0 {
		return EndpointIdentity{}, false
	}

	endpoint := EndpointIdentity{
		ServiceDNSLabel: serviceDNSLabel,
		TargetFQDN:      srvFQDN,
		TargetIP:        "",
		Port:            srvPort,
		Source:          EndpointSourceSRV,
	}

	// Save successful resolution and log SRV hit
	w.endpointMu.Lock()
	w.lastGoodEndpoints[serviceName] = endpoint
	w.endpointMu.Unlock()

	if w.logLimiter.shouldLog(serviceName, "srv_hit") {
		w.logger.Info("dns routing: SRV hit",
			"service", serviceName,
			"target", srvFQDN,
			"port", srvPort)
	}

	return endpoint, true
}

// tryARecordResolution attempts DNS A/AAAA record lookup for service discovery.
// Returns endpoint and true if successful, zero value and false otherwise.
func (w *Watcher) tryARecordResolution(ctx context.Context, svc map[string]any, serviceName, serviceDNSLabel string, port int, logFallback bool) (EndpointIdentity, bool) {
	fqdn := w.tryConstructServiceFQDN(svc)
	if fqdn == "" {
		return EndpointIdentity{}, false
	}

	resolvedIP := w.resolveDNS(ctx, fqdn)
	if resolvedIP == "" {
		return EndpointIdentity{}, false
	}

	endpoint := EndpointIdentity{
		ServiceDNSLabel: serviceDNSLabel,
		TargetFQDN:      fqdn,
		TargetIP:        resolvedIP,
		Port:            port,
		Source:          EndpointSourceA,
	}

	// Save successful resolution and log fallback if needed
	w.endpointMu.Lock()
	w.lastGoodEndpoints[serviceName] = endpoint
	w.endpointMu.Unlock()

	if logFallback && w.logLimiter.shouldLog(serviceName, "fallback_srv_to_a") {
		w.logger.Info("dns routing: fallback",
			"service", serviceName,
			"from", "SRV",
			"to", "A")
	}

	return endpoint, true
}

// handleDNSFailure manages cooldown period and last-good endpoint fallback.
// Returns last-good endpoint and true if in cooldown, zero value and false otherwise.
func (w *Watcher) handleDNSFailure(serviceName string) (EndpointIdentity, bool) {
	w.endpointMu.RLock()
	inCooldown := time.Since(w.lastDNSFailure) < w.dnsFailureCooldown
	lastGood, hasLastGood := w.lastGoodEndpoints[serviceName]
	w.endpointMu.RUnlock()

	if inCooldown && hasLastGood {
		// Reuse last-good endpoint during cooldown
		if w.logLimiter.shouldLog(serviceName, "using_last_good") {
			w.logger.Info("dns routing: using last-good endpoints",
				"service", serviceName,
				"reason", "DNS failure",
				"source", lastGood.Source.String())
		}
		return lastGood, true
	}

	if !inCooldown {
		// Cooldown expired - record new failure timestamp
		w.endpointMu.Lock()
		w.lastDNSFailure = time.Now()
		w.endpointMu.Unlock()
	}

	return EndpointIdentity{}, false
}

// resolveServiceEndpoint resolves a service to a canonical endpoint identity (PR5).
// In cluster mode, prefers FQDN-based routing with stable fallback chain.
// Returns normalized EndpointIdentity with source tracking for observability.
func (w *Watcher) resolveServiceEndpoint(ctx context.Context, svc map[string]any) EndpointIdentity {
	serviceName := strings.TrimSpace(fmt.Sprint(svc["Name"]))
	addr := strings.TrimSpace(fmt.Sprint(svc["Address"]))

	host := ""
	if addr != "" {
		host = strings.TrimSpace(strings.Split(addr, ":")[0])
	}

	port := Utility.ToInt(svc["Port"])
	if port == 0 {
		port = Utility.ToInt(svc["Proxy"])
	}

	serviceDNSLabel := normalizeDNSLabel(serviceName)

	// If explicit host is provided, use it
	if endpoint, ok := w.resolveExplicitHost(host, serviceDNSLabel, port); ok {
		return endpoint
	}

	// No explicit host - try DNS-based routing in cluster mode
	if w.isClusterMode() {
		// PR4.1: Try SRV lookup first for service port discovery
		if endpoint, ok := w.trySRVResolution(ctx, serviceName, serviceDNSLabel); ok {
			return endpoint
		}

		// PR4: Fall back to A/AAAA lookup
		srvAttempted := serviceName != ""
		if endpoint, ok := w.tryARecordResolution(ctx, svc, serviceName, serviceDNSLabel, port, srvAttempted); ok {
			return endpoint
		}

		// PR5: DNS resolution failed - check cooldown period
		if endpoint, ok := w.handleDNSFailure(serviceName); ok {
			return endpoint
		}
	}

	// Fallback to localhost for non-cluster mode or DNS resolution failure
	return EndpointIdentity{
		ServiceDNSLabel: serviceDNSLabel,
		TargetFQDN:      "",
		TargetIP:        "127.0.0.1",
		Port:            port,
		Source:          EndpointSourceLocalhost,
	}
}

// isClusterMode returns true if the watcher is operating in cluster mode
func (w *Watcher) isClusterMode() bool {
	return w.clusterNetwork != nil &&
		w.clusterNetwork.Spec != nil &&
		w.clusterNetwork.Spec.ClusterDomain != ""
}

// tryConstructServiceFQDN attempts to construct an FQDN for a service.
// In practice, services should have their Address field set properly by node-agent.
// This is a fallback to handle edge cases.
func (w *Watcher) tryConstructServiceFQDN(svc map[string]any) string {
	// Extract node information from service if available
	// This is a simplified version - in production, services should have
	// proper Address fields set (e.g., "service.node-01.cluster.local")
	if addr := strings.TrimSpace(fmt.Sprint(svc["Address"])); addr != "" {
		// If Address already looks like an FQDN, use it
		if strings.Contains(addr, ".") && !net.ParseIP(addr).IsLoopback() {
			return strings.Split(addr, ":")[0]
		}
	}

	// Cannot construct FQDN without proper service metadata
	return ""
}

// resolveDNS performs DNS lookup using the cache and returns the first IP as a string.
// Returns empty string on failure (caller should fall back to IP-based routing).
func (w *Watcher) resolveDNS(ctx context.Context, fqdn string) string {
	if w.dnsCache == nil {
		return ""
	}

	ips, err := w.dnsCache.Lookup(ctx, fqdn)
	if err != nil {
		w.logger.Warn("DNS lookup failed, falling back to IP routing",
			"fqdn", fqdn, "err", err)
		return ""
	}

	if len(ips) == 0 {
		return ""
	}

	// Return first IP as string
	return ips[0].String()
}

// resolveSRV performs SRV lookup for a service and returns resolved host and port (PR4.1).
// Returns empty host and 0 port on failure (caller should fall back to A/AAAA).
func (w *Watcher) resolveSRV(ctx context.Context, serviceName string) (string, int) {
	if w.dnsCache == nil || !w.isClusterMode() {
		return "", 0
	}

	// Normalize service name for DNS (e.g., "echo.EchoService" -> "echo-echoservice")
	normalizedName := normalizeDNSLabel(serviceName)
	clusterDomain := w.clusterNetwork.Spec.ClusterDomain

	// Perform SRV lookup: _service._tcp.cluster.local
	records, err := w.dnsCache.LookupSRV(ctx, normalizedName, "tcp", clusterDomain)
	if err != nil {
		// Not an error - SRV records may not exist for all services
		return "", 0
	}

	if len(records) == 0 {
		return "", 0
	}

	// Select best SRV record (lowest priority, highest weight)
	bestRecord := selectBestSRV(records)
	if bestRecord == nil {
		return "", 0
	}

	// Resolve SRV target via A/AAAA lookup
	targetHost := w.resolveDNS(ctx, bestRecord.Target)
	if targetHost == "" {
		w.logger.Warn("SRV target resolution failed",
			"service", serviceName,
			"target", bestRecord.Target)
		return "", 0
	}

	return targetHost, int(bestRecord.Port)
}

// selectBestSRV selects the best SRV record from a list (PR4.1).
// Returns record with lowest priority; if tied, highest weight wins.
func selectBestSRV(records []*dnscache.SRVRecord) *dnscache.SRVRecord {
	if len(records) == 0 {
		return nil
	}

	best := records[0]
	for _, rec := range records[1:] {
		if rec.Priority < best.Priority {
			best = rec
		} else if rec.Priority == best.Priority && rec.Weight > best.Weight {
			best = rec
		}
	}

	return best
}

// normalizeDNSLabel normalizes a service name for use as DNS label (PR4.1).
// Converts dots to hyphens and lowercases (e.g., "echo.EchoService" -> "echo-echoservice").
func normalizeDNSLabel(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, ".", "-")
	return name
}

func (w *Watcher) buildServiceResourcesFromEtcd(ctx context.Context, cfg *XDSConfig) ([]map[string]any, error) {
	if cfg == nil || len(cfg.EtcdEndpoints) == 0 {
		return nil, nil
	}

	etcdCfg := clientv3.Config{
		Endpoints:   cfg.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	}

	// TLS is mandatory for etcd connections
	tlsCfg, err := config.GetEtcdTLS()
	if err != nil {
		return nil, fmt.Errorf("etcd TLS configuration required: %w", err)
	}
	etcdCfg.TLS = tlsCfg

	cli, err := clientv3.New(etcdCfg)
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := cli.Get(etcdCtx, "/globular/services/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}

	services := make([]map[string]any, 0, resp.Count)
	for _, kv := range resp.Kvs {
		var svc map[string]any
		if err := json.Unmarshal(kv.Value, &svc); err != nil {
			if w.logger != nil {
				w.logger.Warn("failed to unmarshal service from etcd", "key", string(kv.Key), "err", err)
			}
			continue
		}
		services = append(services, svc)
	}
	return services, nil
}

func normalizeUpstreamHost(host string) string {
	h := strings.TrimSpace(host)
	switch strings.ToLower(h) {
	case "", "0.0.0.0", "*":
		return "127.0.0.1"
	default:
		return h
	}
}

func detectLocalProtocol() string {
	cfg, err := config.GetLocalConfig(true)
	if err != nil {
		// Default to HTTPS if config cannot be read (secure by default)
		return "https"
	}
	protocol := strings.ToLower(strings.TrimSpace(Utility.ToString(cfg["Protocol"])))
	if protocol == "" {
		// Default to HTTPS when not explicitly configured (secure by default)
		return "https"
	}
	return protocol
}

func (w *Watcher) gatewayTLSPaths() (string, string, string, bool) {
	if strings.ToLower(strings.TrimSpace(w.protocol)) != "https" {
		return "", "", "", false
	}
	cert := pathIfExists(config.GetLocalCertificate())
	key := pathIfExists(config.GetLocalServerKeyPath())
	if cert == "" || key == "" {
		if w.logger != nil && !w.gatewayTLSWarned {
			w.logger.Warn("HTTPS configured but gateway certificate/key missing; serving without TLS until certificates are available")
			w.gatewayTLSWarned = true
		}
		return "", "", "", false
	}
	ca := pathIfExists(config.GetLocalCACertificate())
	return cert, key, ca, true
}

func (w *Watcher) clusterNetworkReady() bool {
	return w.clusterNetwork != nil &&
		w.clusterNetwork.Spec != nil &&
		strings.TrimSpace(w.clusterNetwork.Spec.ClusterDomain) != ""
}

// checkCertificateGeneration checks if certificate generation has changed in etcd.
// Returns true if generation changed (snapshot should be rebuilt).
func (w *Watcher) checkCertificateGeneration(ctx context.Context) bool {
	if w.etcdClient == nil {
		return false
	}

	// v1 Conformance: Domain-based etcd keys (violations INV-1.8, INV-1.9)
	// TODO: This violates v1 invariants - certificate tracking should use cert IDs
	// Current implementation uses domain in etcd key: /globular/pki/bundles/{domain}
	// This couples certificate lifecycle to domain configuration
	// v1 Solution: Use /globular/pki/certs/{cert_id} with stable IDs
	//
	// For now, keeping existing behavior but documenting violation
	// Full fix requires PKI service refactor to use certificate IDs
	var domain string
	if w.clusterNetwork != nil && w.clusterNetwork.Spec != nil {
		domain = w.clusterNetwork.Spec.ClusterDomain
	}
	if domain == "" {
		// VIOLATION INV-1.9: Hardcoded internal domain
		// Should use clusterID instead of domain name
		domain = "globular.internal"
	}

	// VIOLATION INV-1.8: Domain-based persistent state key
	// Query etcd for certificate generation (domain-partitioned)
	key := fmt.Sprintf("/globular/pki/bundles/%s", domain)
	resp, err := w.etcdClient.Get(ctx, key)
	if err != nil {
		if w.logger != nil && w.certGenerationChecked {
			w.logger.Debug("failed to check certificate generation", "err", err, "domain", domain)
		}
		return false
	}

	if len(resp.Kvs) == 0 {
		return false
	}

	// Parse generation from bundle JSON
	var payload struct {
		Generation uint64 `json:"generation"`
	}
	if err := json.Unmarshal(resp.Kvs[0].Value, &payload); err != nil {
		if w.logger != nil {
			w.logger.Warn("failed to parse certificate generation", "err", err)
		}
		return false
	}

	// Check if generation changed
	if !w.certGenerationChecked {
		w.lastCertGeneration = payload.Generation
		w.certGenerationChecked = true
		if w.logger != nil {
			w.logger.Info("certificate generation initialized", "generation", payload.Generation, "domain", domain)
		}
		return false
	}

	if payload.Generation != w.lastCertGeneration {
		if w.logger != nil {
			w.logger.Info("certificate generation changed - rebuilding snapshot",
				"old", w.lastCertGeneration,
				"new", payload.Generation,
				"domain", domain)
		}
		w.lastCertGeneration = payload.Generation
		return true
	}

	return false
}

// computeFileHash computes SHA256 hash of file content.
// Returns empty string if file cannot be read.
func computeFileHash(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// checkACMECertificateRotation checks if ACME certificate files have changed.
// Returns true if fullchain.pem or privkey.pem content hash has changed.
// This triggers snapshot rebuild for hot certificate rotation.
func (w *Watcher) checkACMECertificateRotation(acmeCertPath, acmeKeyPath string) bool {
	if acmeCertPath == "" || acmeKeyPath == "" {
		return false
	}

	// Compute current hashes
	currentCertHash := computeFileHash(acmeCertPath)
	currentKeyHash := computeFileHash(acmeKeyPath)

	// If files cannot be read, don't trigger rotation
	if currentCertHash == "" || currentKeyHash == "" {
		if w.logger != nil && (w.lastACMECertHash != "" || w.lastACMEKeyHash != "") {
			w.logger.Warn("ACME certificate files unreadable", "cert", acmeCertPath, "key", acmeKeyPath)
		}
		return false
	}

	// Initialize on first check
	if w.lastACMECertHash == "" && w.lastACMEKeyHash == "" {
		if w.validateACMESecret(acmeCertPath, acmeKeyPath) != nil {
			if w.logger != nil {
				w.logger.Warn("ACME certificate invalid on init; keeping state uninitialized", "cert", acmeCertPath, "key", acmeKeyPath)
			}
			return false
		}
		w.lastACMECertHash = currentCertHash
		w.lastACMEKeyHash = currentKeyHash
		if w.logger != nil {
			w.logger.Info("ACME certificate hashes initialized", "cert", acmeCertPath, "key", acmeKeyPath)
		}
		return false
	}

	// Check if hashes changed
	certChanged := currentCertHash != w.lastACMECertHash
	keyChanged := currentKeyHash != w.lastACMEKeyHash

	if certChanged || keyChanged {
		if err := w.validateACMESecret(acmeCertPath, acmeKeyPath); err != nil {
			if w.logger != nil {
				w.logger.Warn("ACME renewal invalid, keeping last good cert", "err", err, "cert_file", acmeCertPath, "key_file", acmeKeyPath)
			}
			return false
		}
		if w.logger != nil {
			w.logger.Info("ACME certificate rotation detected - rebuilding snapshot",
				"cert_changed", certChanged,
				"key_changed", keyChanged,
				"cert_file", acmeCertPath,
				"key_file", acmeKeyPath)
		}
		w.lastACMECertHash = currentCertHash
		w.lastACMEKeyHash = currentKeyHash
		return true
	}

	return false
}

// validateACMESecret attempts to build the ACME SDS secret to ensure PEM is valid.
// It does not mutate watcher state; callers must update hashes after success.
func (w *Watcher) validateACMESecret(certPath, keyPath string) error {
	_, err := controlplane.MakeSecret(secrets.PublicIngressCert, certPath, keyPath, "")
	return err
}
