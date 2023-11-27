package controlplane

import (
	"context"
	"sync"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
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
	// Use a WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the xDS server in a goroutine
	go func() {
		defer wg.Done()
		cb := &test.Callbacks{Debug: l.Debug}
		srv := server.NewServer(ctx, cache_, cb)
		RunServer(srv, port)
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

// AddSnapshot adds a snapshot to the cache.
//  ex. there is how to create a snapshot:
//
//	snap, _ := cache.NewSnapshot("1",
//	map[resource.Type][]types.Resource{
//		resource.ClusterType: {
//			makeCluster(ClusterName),
//			makeCluster(ClusterName2), // Add another cluster
//		},
//		resource.RouteType: {
//			makeRoute(RouteName, ClusterName),
//			makeRoute(RouteName2, ClusterName2), // Add another route
//		},
//		resource.ListenerType: {
//			makeHTTPListener(ListenerName, RouteName, "/path/to/cert1.pem", "/path/to/key1.pem", "/path/to/ca1.pem"),
//			makeHTTPListener(ListenerName2, RouteName2, "/path/to/cert2.pem", "/path/to/key2.pem", "/path/to/ca2.pem"), // Add another listener
//		},
//	},)
// 

func AddSnapshot(nodeID string, snapshot cache.Snapshot) {
	mu.Lock()
	defer mu.Unlock()
	cache_.SetSnapshot(context.Background(), nodeID, snapshot)
}

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