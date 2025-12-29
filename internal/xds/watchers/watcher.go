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

	"github.com/globulario/Globular/internal/xds/builder"
	"github.com/globulario/Globular/internal/xds/server"
	"github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
)

type DownstreamTLSMode string

const (
	defaultNodeID    = "globular-xds"
	defaultInterval  = 5 * time.Second
	defaultListener  = "ingress_listener_443"
	defaultRouteName = "ingress_routes"
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
	lastMod                     time.Time
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

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.sync(ctx); err != nil && !errors.Is(err, errNoChange) {
				w.logger.Warn("xDS sync failed", "err", err)
			}
		}
	}
}

func (w *Watcher) sync(ctx context.Context) error {
	input, version, err := w.buildInput()
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

func (w *Watcher) buildInput() (builder.Input, string, error) {
	if w.configPath != "" {
		fi, err := os.Stat(w.configPath)
		if err != nil {
			if os.IsNotExist(err) {
				return w.buildDynamicInput()
			}
			return builder.Input{}, "", err
		}
		if !fi.ModTime().After(w.lastMod) {
			return builder.Input{}, "", errNoChange
		}
		w.lastMod = fi.ModTime()

		input, err := w.loadFromFile()
		if err != nil {
			return builder.Input{}, "", err
		}
		version := strings.TrimSpace(input.Version)
		if version == "" {
			version = fmt.Sprintf("%d", fi.ModTime().UnixNano())
		}
		return input, version, nil
	}

	return w.buildDynamicInput()
}

func (w *Watcher) loadFromFile() (builder.Input, error) {
	if w.configPath == "" {
		return builder.Input{}, fmt.Errorf("config path is empty")
	}
	data, err := os.ReadFile(w.configPath)
	if err != nil {
		return builder.Input{}, err
	}
	var input builder.Input
	if err := json.Unmarshal(data, &input); err != nil {
		return builder.Input{}, err
	}
	if strings.TrimSpace(input.NodeID) == "" {
		input.NodeID = w.nodeID
	}
	return input, nil
}

func (w *Watcher) buildDynamicInput() (builder.Input, string, error) {
	services, err := config.GetServicesConfigurations()
	if err != nil {
		return builder.Input{}, "", err
	}

	downCert, downKey, downIssuer, err := w.downstreamTLSConfig()
	if err != nil {
		return builder.Input{}, "", err
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

	address, _ := config.GetAddress()
	host, port := parseAddress(address)
	if port == 0 {
		port = 8181
	}
	gatewayCluster := "globular_https"
	gatewayCert, gatewayKey, gatewayCA, _ := w.gatewayTLSPaths()
	clusters = append(clusters, builder.Cluster{
		Name:       gatewayCluster,
		Endpoints:  []builder.Endpoint{{Host: host, Port: uint32(port)}},
		CAFile:     gatewayCA,
		ServerCert: gatewayCert,
		KeyFile:    gatewayKey,
		SNI:        host,
	})
	routes = append(routes, builder.Route{Prefix: "/", Cluster: gatewayCluster})

	listener := builder.Listener{
		Name:       defaultListener,
		RouteName:  defaultRouteName,
		Host:       "0.0.0.0",
		Port:       443,
		CertFile:   gatewayCert,
		KeyFile:    gatewayKey,
		IssuerFile: gatewayCA,
	}

	version := fmt.Sprintf("%d", time.Now().UnixNano())
	input := builder.Input{
		NodeID:   w.nodeID,
		Listener: listener,
		Routes:   routes,
		Clusters: clusters,
		Version:  version,
	}
	return input, version, nil
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
