// internal/globule/xds.go
package globule

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/process"
	Utility "github.com/globulario/utility"

	"github.com/globulario/Globular/internal/controlplane"
)

// Choose your management (xDS) and admin ports in ONE place.
const (
	xdsPort   = 18000
	adminPort = 9901
	nodeID    = "globular-xds"
)

// initControlPlane generates Envoy's ADS bootstrap, starts the xDS server with that port,
// waits on context cancel (SIGINT/SIGTERM), and exits cleanly.
func (globule *Globule) initControlPlane() {
	// Cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// OS signal → cancel ctx
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Received interrupt signal. Shutting down xDS...")
		cancel()
	}()

	// 1) Write Envoy ADS bootstrap so Envoy knows how to reach our xDS server.
	//    This avoids parsing envoy.yml later just to rediscover the port.
	envoyBootstrap := config.GetConfigDir() + "/envoy.yml"
	clusterName := "globular-cluster" // logical name; optional cosmetics

	if err := controlplane.WriteBootstrap(envoyBootstrap, controlplane.BootstrapOptions{
		NodeID:                   nodeID,
		Cluster:                  clusterName,
		XDSHost:                  "127.0.0.1",
		XDSPort:                  xdsPort,
		AdminPort:                adminPort,
		MaxActiveDownstreamConns: 50000, // or whatever you prefer
	}); err != nil {
		fmt.Println("failed to write Envoy bootstrap:", err)
		// We continue; Envoy won't connect without it, but the control-plane can still run.
	}
	// (WriteBootstrap emits: node.id, static xds_cluster (HTTP/2), ADS for LDS/CDS/RDS/EDS, admin on :9901)
	// ref: internal/controlplane/bootstrap.go
	// (node.id must match in AddSnapshot below)
	// NOTE: If you want the admin bound elsewhere, change AdminPort above.

	// 2) Start xDS server on xdsPort (must match bootstrap’s xds_cluster endpoint).
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := controlplane.StartControlPlane(ctx, uint(xdsPort), globule.exit); err != nil {
			fmt.Println("Error starting control plane:", err)
		}
	}()

	// 3) Block until context cancelled, then wait goroutine & exit.
	<-ctx.Done()
	wg.Wait()
	fmt.Println("xDS control-plane shutdown complete.")
}

// SetSnapshot builds a per-service Listener/Route/Cluster and pushes a single
// snapshot addressing *all* locally configured services, plus selected peer endpoints.
func (globule *Globule) SetSnapshot() error {
	services, _ := config.GetServicesConfigurations()

	var snapshots []controlplane.Snapshot
	proxies := make(map[uint32]bool) // avoid duplicate listeners on same proxy port

	// Iterate known services and create one Listener/Cluster/Route per *proxy port*.
	for i := range services {
		svc := services[i]
		name := fmt.Sprintf("%s", svc["Name"])
		host := strings.Split(fmt.Sprintf("%s", svc["Address"]), ":")[0]
		port := uint32(Utility.ToInt(svc["Port"]))
		proxy := uint32(Utility.ToInt(svc["Proxy"]))

		if proxies[proxy] {
			fmt.Println("proxy", proxy, "already set for service", name)
			continue
		}

		// Compute resource names (Envoy resources don’t allow dots well → use underscores)
		safe := strings.ReplaceAll(name, ".", "_")
		clusterName := safe + "_cluster"
		routeName := safe + "_route"
		listenerName := safe + "_listener"

		// Build base endpoints with the *local* instance first.
		endpoints := []controlplane.EndPoint{{Host: host, Port: port, Priority: 100}}

		// Build the Snapshot piece for that proxy port.
		// NOTE about host-rewrite:
		//   controlplane.MakeRoute currently *always* sets HostRewriteLiteral(host).
		//   Pass a *meaningful upstream host* here (the primary backend host),
		//   NOT "0.0.0.0", otherwise gRPC auth/certs may break.
		snap := controlplane.Snapshot{
			ClusterName:  clusterName,
			RouteName:    routeName,
			ListenerName: listenerName,
			ListenerHost: "0.0.0.0", // bind address only
			ListenerPort: proxy,

			EndPoints: endpoints,

			// Upstream (Envoy -> service) mTLS

			ServerCertPath: config.GetLocalServerCertificatePath(),
			KeyFilePath:    config.GetLocalServerKeyPath(),
			CAFilePath:     config.GetLocalCACertificate(),

			// Downstream (client -> Envoy) TLS: use your public chain and key.
			// If you don't want client-mTLS, you may leave IssuerFilePath empty.
			CertFilePath:   config.GetLocalCertificate(),
			IssuerFilePath: config.GetLocalCertificateAuthorityBundle(),

			// Ingress: if non-empty, build one listener with many routes (clusters must exist).
			HostRewrite: host, // IMPORTANT: pick a host for host-rewrite. Use the *primary backend host*.

		}

		// IMPORTANT: pick a host for host-rewrite. Use the *primary backend host*.
		// We'll piggyback on ListenerHost field for MakeRoute, but with an upstream host value.
		// Here we choose the first endpoint host:
		snap.HostRewrite = host

		proxies[proxy] = true
		snapshots = append(snapshots, snap)
	}

	// Push one snapshot containing all resources (clusters/routes/listeners).
	// NodeID must match the bootstrap’s node.id.
	return controlplane.AddSnapshot(nodeID, "1", snapshots)
}

// startEnvoyProxy builds the snapshot, then starts Envoy (supervised). On failure,
// it retries after 5s; the xDS server is already up so Envoy will connect via ADS.
func (globule *Globule) startEnvoyProxy() {
	go func() {
		if err := globule.SetSnapshot(); err != nil {
			fmt.Println("fail to generate Envoy dynamic configuration:", err)
		}
		if err := process.StartEnvoyProxy(); err != nil {
			fmt.Println("fail to start Envoy proxy:", err)
			time.Sleep(5 * time.Second)
			globule.startEnvoyProxy()
			return
		}
	}()
}
