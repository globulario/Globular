package watchers

import (
	"context"
	"log/slog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// routingReconciler watches etcd for routing refresh generation changes
// triggered by leader elections. On change it signals a buffered channel
// so the main watcher loop can rebuild the xDS snapshot.
type routingReconciler struct {
	logger     *slog.Logger
	etcdClient *clientv3.Client

	// Event channel
	routingChanged chan struct{}

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// newRoutingReconciler creates a routing reconciler that watches etcd for
// routing refresh generation changes.
func newRoutingReconciler(logger *slog.Logger, etcdClient *clientv3.Client) *routingReconciler {
	if logger == nil {
		logger = slog.Default()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &routingReconciler{
		logger:         logger,
		etcdClient:     etcdClient,
		routingChanged: make(chan struct{}, 1),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start begins watching etcd for routing refresh generation changes.
func (r *routingReconciler) Start(ctx context.Context) {
	go r.watchRoutingRefresh()
}

// RoutingChangedChan returns a channel that signals routing refresh changes.
func (r *routingReconciler) RoutingChangedChan() <-chan struct{} {
	return r.routingChanged
}

// watchRoutingRefresh watches etcd key /globular/routing/refresh-generation
// for changes and signals the routing changed channel.
func (r *routingReconciler) watchRoutingRefresh() {
	const key = "/globular/routing/refresh-generation"

	r.logger.Info("starting etcd routing refresh watcher", "key", key)

	watchChan := r.etcdClient.Watch(r.ctx, key)

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("etcd routing refresh watcher stopped")
			return

		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				r.logger.Error("etcd routing watch error", "err", watchResp.Err())
				// Retry with backoff
				time.Sleep(5 * time.Second)
				watchChan = r.etcdClient.Watch(r.ctx, key)
				continue
			}

			for _, event := range watchResp.Events {
				if event.Type == clientv3.EventTypePut {
					r.logger.Info("routing refresh generation changed", "key", key)
					// Signal routing change (non-blocking)
					select {
					case r.routingChanged <- struct{}{}:
					default:
						// Already pending, skip
					}
				}
			}
		}
	}
}
