package controlplane

import (
	"context"
	"fmt"
	"sync"

	//"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	//"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

var (
	l      Logger
	mu     sync.Mutex
	cache_ = cache.NewSnapshotCache(false, cache.IDHash{}, l)
)

/**
 * StartControlPlane starts an xDS server at the given port, serving the given
 * snapshot. It blocks until the server terminates or until the context is canceled.
 */
func StartControlPlane(ctx context.Context, port uint, exit chan bool) error {

	// I will test a simple snapshot with a single cluster, route, and listener
	snapshot := GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		l.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		return err
	}

	l.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := cache_.SetSnapshot(context.Background(), "test-id", snapshot); err != nil {
		l.Errorf("snapshot error %q for %+v", err, snapshot)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Run the xDS server in a goroutine
	go func() {
		defer wg.Done()
		cb := &test.Callbacks{Debug: l.Debug}
		srv := server.NewServer(ctx, cache_, cb)
		fmt.Println("Starting the xDS server...")
		RunServer(srv, port)
		fmt.Println("xDS server terminated")
	}()

	// Wait for either the server to finish or the context to be canceled
	select {
	case <-ctx.Done():
		// Context canceled, wait for the goroutine to finish
		wg.Wait()
	case <-exit: // Assuming you have a channel to signal server termination
		// Server terminated, nothing to do here
	}

	return nil
}


// Snapshot represents the configuration snapshot.
type Snapshot struct {
	ClusterName    string
	RouteName      string
	ListenerName   string
	ListenerPort   uint32
	UpstreamHost   string
	UpstreamPort   uint32
	Grpc           bool
	CertFilePath   string
	KeyFilePath    string
	CAFilePath     string
}

// AddSnapshot adds a snapshot to the cache.
/*func AddSnapshot(nodeID string, values []Snapshot) error{
	mu.Lock()
	defer mu.Unlock()

	// Create a map to store resources by type
	resources := make(map[resource.Type][]types.Resource)

	// Iterate over the provided values and create clusters, routes, and listeners
	for _, value := range values {
		cluster := MakeCluster(value.ClusterName, value.CertFilePath, value.KeyFilePath, value.CAFilePath, value.UpstreamHost, value.UpstreamPort)
		route := MakeRoute(value.RouteName, value.ClusterName, value.UpstreamHost)
		listener := MakeHTTPListener(value.ListenerName, value.ClusterName, value.RouteName, value.CertFilePath, value.KeyFilePath, value.CAFilePath, value.ListenerPort, value.Grpc)

		// Add the resources to the map
		resources[resource.ClusterType] = append(resources[resource.ClusterType], cluster)
		resources[resource.RouteType] = append(resources[resource.RouteType], route)
		resources[resource.ListenerType] = append(resources[resource.ListenerType], listener)
	}

	// Create a new snapshot and set it in the cache
	snapshot, err := cache.NewSnapshot(nodeID, resources)
	if err != nil {
		return err
	}

	// Set the snapshot in the cache
	return cache_.SetSnapshot(context.Background(), nodeID, snapshot)
}*/

// RemoveSnapshot removes a snapshot from the cache.
func RemoveSnapshot(nodeID string) {
	mu.Lock()
	defer mu.Unlock()
	cache_.ClearSnapshot(nodeID)
}

// GetSnapshot returns the current snapshot for a given node ID.
func GetSnapshot(nodeID string) (cache.Snapshot, error){
	mu.Lock()
	defer mu.Unlock()
	return cache_.GetSnapshot(nodeID)
}