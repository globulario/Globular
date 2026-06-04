// @awareness namespace=globular.platform
// @awareness component=platform_xds.routing_reconciler
// @awareness file_role=xds_routing_refresh_reconciler_via_typed_controller_rpc
// @awareness enforces=globular.platform:invariant.four_layer.truth_read_via_owner_rpc_not_direct_storage
// @awareness risk=medium
package watchers

import (
	"context"
	"log/slog"
	"time"

	"github.com/globulario/Globular/internal/controllerclient"
)

// routingReconciler polls the cluster_controller's
// GetRoutingRefresh typed RPC to detect routing-refresh epoch
// changes (driven by leader elections). On change it signals a
// buffered channel so the main watcher loop can rebuild the xDS
// snapshot.
//
// History: prior to v1.2.177 this reconciler watched the etcd key
// /globular/routing/refresh-generation directly via clientv3. That
// prefix is owned by cluster_controller (writeRoutingRefresh in
// services/golang/cluster_controller/cluster_controller_server/server.go),
// so xDS reading raw etcd violated
// invariant:four_layer.truth_read_via_owner_rpc_not_direct_storage.
// The migration trades sub-second watch latency for ~2s polling
// latency; routing-refresh is a rare cluster event (leader
// elections, weeks/months scale), so the trade-off is acceptable.
type routingReconciler struct {
	logger     *slog.Logger
	controller *controllerclient.Client

	// pollInterval is the cadence at which the reconciler polls
	// GetRoutingRefresh. Configurable for testing; defaults to 2s.
	pollInterval time.Duration

	// Event channel: 1-deep buffer, drop-on-pending semantics.
	routingChanged chan struct{}

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// newRoutingReconciler creates a routing reconciler that polls the
// cluster_controller via the typed GetRoutingRefresh RPC.
//
// controllerAddr must be the controller's gRPC endpoint
// (host:port). A nil/zero-value addr leaves the reconciler in a
// degraded mode where it never signals — the main watcher loop
// continues to drive snapshots on its sync timer.
func newRoutingReconciler(logger *slog.Logger, controllerAddr string) *routingReconciler {
	if logger == nil {
		logger = slog.Default()
	}
	ctx, cancel := context.WithCancel(context.Background())

	var client *controllerclient.Client
	if controllerAddr != "" {
		client = controllerclient.New(controllerAddr)
	}

	return &routingReconciler{
		logger:         logger,
		controller:     client,
		pollInterval:   2 * time.Second,
		routingChanged: make(chan struct{}, 1),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start begins polling the controller for routing-refresh changes.
func (r *routingReconciler) Start(ctx context.Context) {
	go r.pollRoutingRefresh()
}

// RoutingChangedChan returns a channel that signals routing refresh changes.
func (r *routingReconciler) RoutingChangedChan() <-chan struct{} {
	return r.routingChanged
}

// pollRoutingRefresh polls cluster_controller.GetRoutingRefresh on a
// fixed cadence and signals routingChanged whenever the epoch
// advances. On RPC error the loop logs and continues to the next
// tick (no escalating backoff: the controller is on the cluster's
// hot path and either succeeds within the poll interval or surfaces
// the failure via the watcher's broader health signals).
func (r *routingReconciler) pollRoutingRefresh() {
	if r.controller == nil {
		r.logger.Warn("routing reconciler started without a controller client — no routing-change signals will be emitted")
		return
	}
	r.logger.Info("starting controller routing-refresh poller",
		"poll_interval", r.pollInterval,
		"controller_addr", r.controller.Address())

	var lastEpoch uint64
	initialized := false

	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("controller routing-refresh poller stopped")
			return

		case <-ticker.C:
			callCtx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
			resp, err := r.controller.GetRoutingRefresh(callCtx)
			cancel()
			if err != nil {
				r.logger.Warn("controller GetRoutingRefresh poll failed", "err", err)
				continue
			}
			epoch := resp.GetEpoch()
			if !initialized {
				lastEpoch = epoch
				initialized = true
				continue
			}
			if epoch != lastEpoch {
				r.logger.Info("controller routing-refresh epoch advanced",
					"prev", lastEpoch, "next", epoch,
					"leader_addr", resp.GetLeaderAddr())
				lastEpoch = epoch
				select {
				case r.routingChanged <- struct{}{}:
				default:
					// Already pending, skip.
				}
			}
		}
	}
}
