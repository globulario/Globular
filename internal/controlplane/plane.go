// internal/controlplane/plane.go
package controlplane

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache_v3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	resource_v3 "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	server_v3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	test_v3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

var (
	l     Logger
	mu    sync.Mutex
	cache = cache_v3.NewSnapshotCache(false, cache_v3.IDHash{}, l)
)

// StartControlPlane runs the ADS/xDS management server.
func StartControlPlane(ctx context.Context, port uint, exit chan bool) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		cb := &test_v3.Callbacks{Debug: l.Debug}
		srv := server_v3.NewServer(ctx, cache, cb)
		fmt.Println("Starting the xDS server...")
		RunServer(srv, port)
		fmt.Println("xDS server terminated")
	}()

	select {
	case <-ctx.Done():
		wg.Wait()
	case <-exit:
	}
	return nil
}

// EndPoint represents an upstream endpoint.
type EndPoint struct {
	Host     string
	Port     uint32
	Priority uint32
}

// Snapshot can describe either:
//   - a per-service sidecar (Cluster + single Route + Listener), or
//   - a shared ingress listener with multiple routes (no cluster/endpoints needed here).
type Snapshot struct {
	// Sidecar fields
	ClusterName  string
	RouteName    string
	ListenerName string
	ListenerHost string // used for host-rewrite toward upstream
	ListenerPort uint32
	EndPoints    []EndPoint

	// Upstream (Envoy → service) mTLS
	ServerCertPath string
	KeyFilePath    string
	CAFilePath     string
	SNI            string // Server Name Indication for upstream TLS

	// Downstream (client → Envoy) TLS
	CertFilePath   string
	IssuerFilePath string

	// Ingress: if non-empty, build one listener with many routes (clusters must exist).
	IngressRoutes []IngressRoute

	HostRewrite string // e.g., backend service host
}

// AddSnapshot supports both sidecars and a shared ingress listener.
func AddSnapshot(id, version string, values []Snapshot) error {
	mu.Lock()
	defer mu.Unlock()

	resources := make(map[string][]types.Resource)

	for _, v := range values {
		// ---------------------------
		// Ingress (shared listener)
		// ---------------------------
		if len(v.IngressRoutes) > 0 {
			host := strings.TrimSpace(v.ListenerHost)
			if host == "" {
				host = "0.0.0.0"
			}
			port := v.ListenerPort
			if port == 0 {
				port = defaultIngressPort(host)
			}
			listenerName := v.ListenerName
			if listenerName == "" {
				listenerName = fmt.Sprintf("ingress_listener_%d", port)
			}

			rc := MakeRoutes(v.RouteName, v.IngressRoutes)
			ln := MakeHTTPListener(
				host, port,
				listenerName, v.RouteName,
				v.CertFilePath, v.KeyFilePath, v.IssuerFilePath,
			)

			resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], rc)
			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], ln)
			continue
		}

		// ---------------------------
		// Cluster-only or Sidecar
		// ---------------------------

		// Always register the cluster if provided.
		if v.ClusterName != "" && len(v.EndPoints) > 0 {
			cluster := MakeCluster(v.ClusterName, v.ServerCertPath, v.KeyFilePath, v.CAFilePath, v.SNI, v.EndPoints)
			resources[resource_v3.ClusterType] = append(resources[resource_v3.ClusterType], cluster)
		}

		// Only build route+listener when explicitly requested with a real port.
		if v.ListenerPort > 0 && v.ListenerName != "" && v.RouteName != "" {
			host := strings.TrimSpace(v.ListenerHost)
			if host == "" {
				host = "0.0.0.0"
			}

			route := MakeRoute(v.RouteName, v.ClusterName, v.HostRewrite)
			listener := MakeHTTPListener(
				host, v.ListenerPort,
				v.ListenerName, v.RouteName,
				v.CertFilePath, v.KeyFilePath, v.IssuerFilePath,
			)

			resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], route)
			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], listener)
		}
	}

	snap, err := cache_v3.NewSnapshot(version, resources)
	if err != nil {
		return err
	}
	l.Infof("xds: pushing snapshot %s (L:%d R:%d C:%d)",
		version,
		len(resources[resource_v3.ListenerType]),
		len(resources[resource_v3.RouteType]),
		len(resources[resource_v3.ClusterType]),
	)
	return cache.SetSnapshot(context.Background(), id, snap)
}

func defaultIngressPort(host string) uint32 {
	if isLocalhostHost(host) {
		return 8443
	}
	return 443
}

func isLocalhostHost(host string) bool {
	switch strings.ToLower(strings.TrimSpace(host)) {
	case "", "0.0.0.0", "127.0.0.1", "localhost":
		return true
	default:
		return false
	}
}

// RemoveSnapshot clears cached resources for the given node.
func RemoveSnapshot(nodeID string) {
	mu.Lock()
	defer mu.Unlock()
	cache.ClearSnapshot(nodeID)
}

// GetSnapshot fetches the current snapshot for a node.
func GetSnapshot(nodeID string) (cache_v3.ResourceSnapshot, error) {
	mu.Lock()
	defer mu.Unlock()
	return cache.GetSnapshot(nodeID)
}
