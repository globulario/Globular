// internal/controlplane/plane.go
package controlplane

import (
	"context"
	"fmt"
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
		// Ingress block: only builds RouteConfiguration + Listener.
		if len(v.IngressRoutes) > 0 {
			ingressRC := MakeRoutes(v.RouteName, v.IngressRoutes)
			ingressL := MakeHTTPListener(
				v.ListenerHost, v.ListenerPort, v.ListenerName, v.RouteName,
				v.CertFilePath, v.KeyFilePath, v.IssuerFilePath,
			)
			resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], ingressRC)
			resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], ingressL)
			continue
		}

		// Regular sidecar (cluster + single route + listener)
		cluster := MakeCluster(v.ClusterName, v.ServerCertPath, v.KeyFilePath, v.CAFilePath, v.EndPoints)
		route := MakeRoute(v.RouteName, v.ClusterName, v.ListenerHost)
		listener := MakeHTTPListener(
			v.ListenerHost, v.ListenerPort, v.ListenerName, v.RouteName,
			v.CertFilePath, v.KeyFilePath, v.IssuerFilePath,
		)

		resources[resource_v3.ClusterType] = append(resources[resource_v3.ClusterType], cluster)
		resources[resource_v3.RouteType] = append(resources[resource_v3.RouteType], route)
		resources[resource_v3.ListenerType] = append(resources[resource_v3.ListenerType], listener)
	}

	snap, err := cache_v3.NewSnapshot(version, resources)
	if err != nil {
		return err
	}
	l.Debugf("will serve snapshot %+v", snap)
	return cache.SetSnapshot(context.Background(), id, snap)
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
