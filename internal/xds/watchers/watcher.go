package watchers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/globulario/Globular/internal/controlplane"
	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/server"
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
	configCached                *XDSConfig
	configMod                   time.Time
}

// New creates a watcher bound to the given server.
func New(logger *slog.Logger, srv *server.XDSServer, configPath, nodeID string, interval time.Duration, downstreamMode DownstreamTLSMode) *Watcher {
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
		logger:            logger,
		server:            srv,
		configPath:        strings.TrimSpace(configPath),
		nodeID:            nodeID,
		interval:          interval,
		downstreamTLSMode: downstreamMode,
		protocol:          protocol,
	}
}

// Run starts the watcher loop and blocks until the context is canceled.
func (w *Watcher) Run(ctx context.Context) error {
	if w.server == nil {
		return fmt.Errorf("server is nil")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
		w.logger.Warn("initial xDS sync failed", "err", err)
	}

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
			if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
				w.logger.Warn("xDS sync failed", "err", err)
			}
			next := w.interval
			if next <= 0 {
				next = defaultInterval
			}
			timer.Reset(next)
		}
	}
}

func (w *Watcher) sync(ctx context.Context) error {
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
	w.logger.Info("xDS snapshot pushed", "node_id", input.NodeID, "version", version)
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

func (w *Watcher) buildDynamicInput(ctx context.Context, cfg *XDSConfig) (builder.Input, string, error) {
	clusters, routes, err := w.buildServiceResources()
	if err != nil {
		return builder.Input{}, "", err
	}

	var listener builder.Listener
	ingressSpec, err := w.buildIngressSpec(ctx, cfg)
	if err != nil {
		return builder.Input{}, "", err
	}
	var ingressHTTPPort uint32
	var enableHTTPRedirect bool
	var gatewayPort uint32
	if ingressSpec != nil {
		clusters = append(clusters, ingressSpec.Clusters...)
		routes = append(routes, ingressSpec.Routes...)
		listener = ingressSpec.Listener
		ingressHTTPPort = ingressSpec.HTTPPort
		enableHTTPRedirect = ingressSpec.EnableHTTPRedirect
		gatewayPort = ingressSpec.GatewayPort
		if listener.Host == "" {
			listener.Host = "0.0.0.0"
		}
		if listener.RouteName == "" {
			listener.RouteName = defaultRouteName
		}
		if listener.Name == "" {
			port := listener.Port
			if port == 0 {
				port = controlplane.DefaultIngressPort(listener.Host)
			}
			listener.Name = fmt.Sprintf("%s_%d", listenerNameBase, port)
		}
	} else {
		legacyClusters, legacyRoutes, legacyListener, err := w.buildLegacyGatewayResources()
		if err != nil {
			return builder.Input{}, "", err
		}
		clusters = append(clusters, legacyClusters...)
		routes = append(routes, legacyRoutes...)
		listener = legacyListener
	}

	if listener.RouteName == "" {
		listener.RouteName = defaultRouteName
	}

	version := fmt.Sprintf("%d", time.Now().UnixNano())
	input := builder.Input{
		NodeID:             w.nodeID,
		Listener:           listener,
		Routes:             routes,
		Clusters:           clusters,
		IngressHTTPPort:    ingressHTTPPort,
		EnableHTTPRedirect: enableHTTPRedirect,
		GatewayPort:        gatewayPort,
		Version:            version,
	}
	return input, version, nil
}

func (w *Watcher) buildServiceResources() ([]builder.Cluster, []builder.Route, error) {
	services, err := config.GetServicesConfigurations()
	if err != nil {
		return nil, nil, err
	}

	downCert, downKey, downIssuer, err := w.downstreamTLSConfig()
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
		address := strings.TrimSpace(fmt.Sprint(svc["Address"]))
		host := strings.Split(address, ":")[0]
		port := Utility.ToInt(svc["Port"])
		if name == "" || host == "" || port == 0 {
			continue
		}
		clusterName := strings.ReplaceAll(name, ".", "_") + "_cluster"
		if _, ok := addedClusters[clusterName]; ok {
			continue
		}
		addedClusters[clusterName] = struct{}{}

		clusters = append(clusters, builder.Cluster{
			Name:       clusterName,
			Endpoints:  []builder.Endpoint{{Host: host, Port: uint32(port)}},
			ServerCert: downCert,
			KeyFile:    downKey,
			CAFile:     downIssuer,
			SNI:        host,
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

func (w *Watcher) buildLegacyGatewayResources() ([]builder.Cluster, []builder.Route, builder.Listener, error) {
	host, port := readGatewayAddress()
	if port == 0 {
		port = 8181
	}
	upstreamHost := normalizeUpstreamHost(host)
	gatewayCluster := "globular_https"
	gatewayCert, gatewayKey, gatewayCA, _ := w.gatewayTLSPaths()
	clusters := []builder.Cluster{{
		Name:       gatewayCluster,
		Endpoints:  []builder.Endpoint{{Host: upstreamHost, Port: uint32(port)}},
		CAFile:     gatewayCA,
		ServerCert: gatewayCert,
		KeyFile:    gatewayKey,
		SNI:        host,
	}}
	routes := []builder.Route{{Prefix: "/", Cluster: gatewayCluster}}

	listenerPort := uint32(443)
	lowerHost := strings.ToLower(strings.TrimSpace(host))
	if lowerHost == "" || lowerHost == "0.0.0.0" || lowerHost == "127.0.0.1" || lowerHost == "localhost" {
		listenerPort = 8443
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
	return clusters, routes, listener, nil
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
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.EtcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
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
		clusters = append(clusters, builder.Cluster{
			Name:      name,
			Endpoints: endpoints,
		})
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
		Routes:   routes,
		Clusters: clusters,
		HTTPPort: fb.Ingress.HTTPPort,
	}
}

func (w *Watcher) applyIngressSettings(spec *IngressSpec, cfg *XDSConfig) {
	if spec == nil || cfg == nil {
		return
	}
	if spec.Listener.Host == "" {
		spec.Listener.Host = "0.0.0.0"
	}
	if spec.Listener.Port == 0 {
		spec.Listener.Port = cfg.Ingress.HTTPSPort
	}
	if spec.HTTPPort == 0 {
		spec.HTTPPort = cfg.Ingress.HTTPPort
	}
	spec.EnableHTTPRedirect = cfg.ingressRedirectEnabled()
	spec.GatewayPort = cfg.gatewayPort()
	w.applyIngressTLS(spec, cfg)
}

func (w *Watcher) applyIngressTLS(spec *IngressSpec, cfg *XDSConfig) {
	if spec == nil || cfg == nil {
		return
	}
	tlsConfig := cfg.Ingress
	if !tlsConfig.TLS.Enabled {
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

func readGatewayAddress() (string, int) {
	cfg, err := config.GetLocalConfig(true)
	if err == nil && cfg != nil {
		if gwRaw, ok := cfg["gateway"]; ok {
			if gwMap, ok := gwRaw.(map[string]interface{}); ok {
				if listen := strings.TrimSpace(Utility.ToString(gwMap["listen"])); listen != "" {
					host, port := parseAddress(listen)
					if port != 0 {
						return host, port
					}
				}
			}
		}
	}
	address, _ := config.GetAddress()
	return parseAddress(address)
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
		return "http"
	}
	protocol := strings.ToLower(strings.TrimSpace(Utility.ToString(cfg["Protocol"])))
	if protocol == "" {
		return "http"
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
