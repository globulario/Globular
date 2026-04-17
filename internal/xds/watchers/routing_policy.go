package watchers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/globulario/Globular/internal/xds/builder"
	ai_routerpb "github.com/globulario/services/golang/ai_router/ai_routerpb"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Policy staleness thresholds (from AI-Router-Final-Tightenings.md).
const (
	policyFreshThreshold = 15 * time.Second
	policyWarnThreshold  = 45 * time.Second
	policyStaleThreshold = 90 * time.Second
	routerCallTimeout    = 50 * time.Millisecond
)

// routerConn holds a lazy gRPC connection to the AI Router.
type routerConn struct {
	cc     *grpc.ClientConn
	client ai_routerpb.AiRouterServiceClient
	addr   string
}

// getRouterClient returns a connected AI Router client, creating the
// connection lazily on first call. Returns nil if the router is unavailable.
func (w *Watcher) getRouterClient() ai_routerpb.AiRouterServiceClient {
	if w.routerClient != nil && w.routerClient.client != nil {
		return w.routerClient.client
	}

	addr := config.ResolveServiceAddr("ai_router.AiRouterService", "")
	if addr == "" {
		return nil // router not registered
	}

	// Try TLS with Globular CA cert, fall back to plaintext.
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tc := loadRouterTLSCreds(); tc != nil {
		creds = grpc.WithTransportCredentials(tc)
	}

	cc, err := grpc.Dial(addr,
		creds,
		grpc.WithBlock(),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		if w.logger != nil {
			w.logger.Debug("ai_router unavailable", "addr", addr, "err", err)
		}
		return nil
	}

	w.routerClient = &routerConn{
		cc:     cc,
		client: ai_routerpb.NewAiRouterServiceClient(cc),
		addr:   addr,
	}
	if w.logger != nil {
		w.logger.Info("ai_router connected", "addr", addr)
	}
	return w.routerClient.client
}

// applyRoutingPolicy queries the AI Router and applies endpoint weights
// to the cluster list. Safe to call when the router is unavailable — returns
// clusters unchanged.
func (w *Watcher) applyRoutingPolicy(ctx context.Context, clusters []builder.Cluster) []builder.Cluster {
	client := w.getRouterClient()
	if client == nil {
		return clusters // router unavailable, passthrough
	}

	callCtx, cancel := context.WithTimeout(ctx, routerCallTimeout)
	defer cancel()

	// Inject service auth metadata so the ai_router's interceptor accepts
	// the call. Without this, the request arrives as anonymous and is rejected
	// with cluster_id_missing after cluster initialization.
	callCtx = injectServiceAuth(callCtx)

	resp, err := client.GetRoutingPolicy(callCtx, &ai_routerpb.GetRoutingPolicyRequest{})
	if err != nil {
		if w.logger != nil {
			w.logger.Debug("ai_router GetRoutingPolicy failed", "err", err)
		}
		return clusters // error, passthrough
	}

	policy := resp.GetPolicy()
	if policy == nil || len(policy.Services) == 0 {
		return clusters // neutral policy, passthrough
	}

	// Check staleness.
	age := time.Since(time.UnixMilli(policy.ComputedAtMs))
	if age > policyStaleThreshold {
		if w.logger != nil {
			w.logger.Warn("ai_router policy stale, ignoring",
				"age", age.Round(time.Second),
				"generation", policy.Generation)
		}
		return clusters // stale, passthrough
	}
	if age > policyWarnThreshold {
		if w.logger != nil {
			w.logger.Warn("ai_router policy aging",
				"age", age.Round(time.Second),
				"generation", policy.Generation)
		}
	}

	// Check mode — only apply weights in ACTIVE mode.
	if policy.Mode != ai_routerpb.RouterMode_ROUTER_ACTIVE {
		return clusters // observe/neutral mode, passthrough
	}

	// Apply weights to endpoints.
	applied := 0
	for i := range clusters {
		// Map cluster name to service name.
		// Cluster names are like "event_cluster" → service "event.EventService"
		// Try matching by iterating policy services.
		for svcName, sp := range policy.Services {
			if !clusterMatchesService(clusters[i].Name, svcName) {
				continue
			}
			if sp.Confidence < 0.1 {
				continue // low confidence, skip
			}
			for j := range clusters[i].Endpoints {
				ep := &clusters[i].Endpoints[j]
				key := fmt.Sprintf("%s:%d", ep.Host, ep.Port)
				// Also try instance format (127.0.0.1:proxyPort) from Prometheus.
				if w, ok := sp.Weights[key]; ok && w > 0 {
					ep.Weight = w
					applied++
				}
			}
		}
	}

	if applied > 0 && w.logger != nil {
		w.logger.Info("ai_router weights applied",
			"generation", policy.Generation,
			"endpoints_weighted", applied,
			"age", age.Round(time.Second))
	}

	return clusters
}

// clusterMatchesService checks if a cluster name corresponds to a service name.
// Cluster names: "event_cluster", "authentication_cluster", etc.
// Service names: "event.EventService", "authentication.AuthenticationService", etc.
func clusterMatchesService(clusterName, serviceName string) bool {
	// Extract service prefix from cluster name: "event_cluster" → "event"
	prefix := strings.TrimSuffix(clusterName, "_cluster")
	if prefix == clusterName {
		return false // no _cluster suffix
	}
	// Compare with the package part of the service name: "event.EventService" → "event"
	parts := strings.SplitN(serviceName, ".", 2)
	if len(parts) == 0 {
		return false
	}
	return strings.EqualFold(prefix, parts[0])
}

// loadRouterTLSCreds loads the Globular CA for TLS to the AI Router.
// Returns nil if no CA found (caller falls back to plaintext).
func loadRouterTLSCreds() credentials.TransportCredentials {
	caPaths := []string{
		"/var/lib/globular/pki/ca.crt",
		"/var/lib/globular/pki/ca.pem",
	}
	for _, path := range caPaths {
		pem, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(pem) {
			continue
		}
		return credentials.NewTLS(&tls.Config{
			RootCAs: pool,
		})
	}
	return nil
}

// injectServiceAuth generates a short-lived service token and adds it to
// the outgoing gRPC metadata. This allows internal service-to-service calls
// (like xDS → ai_router) to pass the auth interceptor's cluster_id check.
func injectServiceAuth(ctx context.Context) context.Context {
	mac, err := config.GetMacAddress()
	if err != nil {
		return ctx
	}
	token, err := security.GenerateServiceToken(mac)
	if err != nil {
		return ctx
	}
	domain, _ := config.GetDomain()
	md := metadata.New(map[string]string{
		"token":         token,
		"authorization": "Bearer " + token,
		"mac":           mac,
		"domain":        domain,
	})
	return metadata.NewOutgoingContext(ctx, md)
}

// logRoutingPolicy logs the routing policy details at debug level.
func logRoutingPolicy(logger *slog.Logger, policy *ai_routerpb.RoutingPolicy) {
	if logger == nil || policy == nil {
		return
	}
	for svc, sp := range policy.Services {
		logger.Debug("routing_policy",
			"service", svc,
			"confidence", sp.Confidence,
			"endpoints", len(sp.Weights),
			"reasons", sp.Reasons)
	}
}
