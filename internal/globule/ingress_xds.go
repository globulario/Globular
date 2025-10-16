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

// Public HTTPS ingress on :443 forwards everything to Globular's internal HTTPS on 127.0.0.1:8181.
// SetSnapshot: single public HTTPS listener on :443.
// - One CLUSTER per backend service (host:port) using your upstream TLS material
// - One ROUTE per service: "/<Pkg.Service>/" -> its cluster
// - Final catch-all "/" -> gateway cluster (127.0.0.1:8181)
// - No per-service listeners, no proxy ports.
func (globule *Globule) SetSnapshot() error {
	services, err := config.GetServicesConfigurations()
	if err != nil {
		return fmt.Errorf("unable to read local services configuration: %w", err)
	}

	// helpers
	pathIfExists := func(p string) string {
		s := strings.TrimSpace(fmt.Sprintf("%v", p))
		if s == "" {
			return ""
		}
		if Utility.Exists(s) {
			return s
		}
		return ""
	}
	require := func(label, p string) (string, error) {
		if q := pathIfExists(p); q != "" {
			return q, nil
		}
		return "", fmt.Errorf("%s path missing or empty: %q", label, p)
	}

	const (
		ingressVHost   = "ingress_routes"
		ingressLdsName = "ingress_listener_443"
		gatewayCluster = "globular_https"
	)

	// ---------------- Downstream TLS (public :443) ----------------
	downCert, err := require("downstream certificate", config.GetLocalClientCertificatePath())
	if err != nil {
		return err
	}
	downKey, err := require("downstream private key", config.GetLocalClientKeyPath())
	if err != nil {
		return err
	}
	downIssuer := pathIfExists(config.GetLocalCACertificate()) // optional

	// ---------------- Upstream TLS (Envoy -> services) ----------------
	// Use EXACTLY your getters, no improvisation.
	upCA := pathIfExists(config.GetLocalCertificateAuthorityBundle())
	upCert := pathIfExists(config.GetLocalCertificate())  // client cert (only if mTLS required)
	upKey := pathIfExists(config.GetLocalServerKeyPath()) // client key  (only if mTLS required)

	var snapshots []controlplane.Snapshot
	var routes []controlplane.IngressRoute
	added := map[string]bool{}
	var allClusterNames []string

	// (A) Per-service clusters + routes
	for i := range services {
		svc := services[i]

		name := fmt.Sprintf("%s", svc["Name"]) // e.g. "rbac.RbacService"
		host := strings.Split(fmt.Sprintf("%s", svc["Address"]), ":")[0]
		port := uint32(Utility.ToInt(svc["Port"])) // service's real port
		if name == "" || host == "" || port == 0 {
			fmt.Printf("[xds] skip service: name=%q host=%q port=%d\n", name, host, port)
			continue
		}

		safe := strings.ReplaceAll(name, ".", "_")
		clusterName := safe + "_cluster"

		if !added[clusterName] {
			snapshots = append(snapshots, controlplane.Snapshot{
				ClusterName: clusterName,
				EndPoints:   []controlplane.EndPoint{{Host: host, Port: port, Priority: 0}},

				// Upstream TLS:
				// - CA:     config.GetLocalCACertificate()
				// - mTLS:   config.GetLocalServerCertificatePath()/GetLocalServerKeyPath() ONLY if your backend requires client-auth
				// - SNI:    backend host
				CAFilePath:     downIssuer,
				ServerCertPath: downCert, // leave as-is; backend will ignore if not enforcing mTLS
				KeyFilePath:    downKey,  // leave as-is; backend will ignore if not enforcing mTLS
				SNI:            host,
			})
			added[clusterName] = true
			allClusterNames = append(allClusterNames, clusterName)
		}

		// Route "/Pkg.Service/" -> cluster (order matters: service routes first)
		routes = append(routes, controlplane.IngressRoute{
			Prefix:  "/" + name + "/",
			Cluster: clusterName,
		})
	}

	// (A.1) Add a DEBUG health route so grpcurl via ingress hits a real gRPC cluster
	// Prefer the authentication service if present, else first service cluster if any.
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
		// Ensure the health route is evaluated BEFORE the catch-all
		routes = append([]controlplane.IngressRoute{{
			Prefix:  "/grpc.health.v1.Health/",
			Cluster: healthTarget,
		}}, routes...)
	}

	// (B) Fallback gateway cluster for "/"
	//     Uses the SAME upstream TLS sources you identified.
	address, _ := config.GetAddress()
	hostname := strings.Split(address, ":")[0]
	port := Utility.ToInt(strings.Split(address, ":")[1])
	snapshots = append(snapshots, controlplane.Snapshot{
		ClusterName:    gatewayCluster,
		EndPoints:      []controlplane.EndPoint{{Host: hostname, Port: uint32(port), Priority: 0}},
		CAFilePath:     upCA,
		ServerCertPath: upCert, // if your gateway enforces mTLS from Envoy, this is already correct
		KeyFilePath:    upKey,
		SNI:            hostname, // the hostname on your 8181 cert
	})
	routes = append(routes, controlplane.IngressRoute{Prefix: "/", Cluster: gatewayCluster}) // LAST

	// (C) Single public listener on :443 with all routes
	snapshots = append(snapshots, controlplane.Snapshot{
		RouteName:      ingressVHost,
		ListenerName:   ingressLdsName,
		ListenerHost:   "0.0.0.0",
		ListenerPort:   443,
		CertFilePath:   upCert,
		KeyFilePath:    upKey,
		IssuerFilePath: upCA,
		IngressRoutes:  routes,
	})

	// push (bump version)
	return controlplane.AddSnapshot(nodeID, "5", snapshots)
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
