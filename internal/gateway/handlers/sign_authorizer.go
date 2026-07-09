package handlers

// sign_authorizer.go — issuance gate for /sign_ca_certificate.
//
// A cluster-CA-signed certificate is an mTLS credential: etcd runs
// client-cert-auth against the same CA, so whoever obtains a cert gets full
// read/write access to the cluster's source of truth. Signing therefore may
// not be open to anyone who can reach the gateway. Three requester classes
// are legitimate:
//
//  1. the local node itself (Day-0 bootstrap, local renewals) — identified by
//     a direct connection from a loopback or local-interface address;
//  2. a Day-1 joining node — identified by a controller-issued join_id
//     (obtained from POST /join/authorize, which verified the join token);
//  3. an existing cluster member (cert renewal) — identified by a direct
//     connection from an IP the controller lists for a registered node.
//
// Proxied requests (X-Forwarded-For / Forwarded present) never qualify for
// the IP-based classes: the TCP peer is the proxy, not the requester. They
// must present a join_id.
//
// Everything else is denied. Controller unreachable → classes 2 and 3 fail
// closed; class 1 keeps working so a degraded cluster can still self-issue
// locally.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/globulario/Globular/internal/controllerclient"
)

// signControllerClient is the subset of controllerclient.Client the
// authorizer needs; narrowed for testability.
type signControllerClient interface {
	GetJoinRequestStatusPlanJSON(ctx context.Context, joinID string) ([]byte, error)
	ListNodeIPs(ctx context.Context) ([]string, error)
}

// controllerSignBackend adapts controllerclient.Client to signControllerClient.
type controllerSignBackend struct {
	client *controllerclient.Client
}

func (b controllerSignBackend) GetJoinRequestStatusPlanJSON(ctx context.Context, joinID string) ([]byte, error) {
	resp, err := b.client.GetJoinRequestStatus(ctx, joinID)
	if err != nil {
		return nil, err
	}
	return resp.GetPlanJson(), nil
}

func (b controllerSignBackend) ListNodeIPs(ctx context.Context) ([]string, error) {
	resp, err := b.client.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, n := range resp.GetNodes() {
		ips = append(ips, n.GetIdentity().GetIps()...)
		if ep := n.GetAgentEndpoint(); ep != "" {
			if host, _, err := net.SplitHostPort(ep); err == nil {
				ips = append(ips, host)
			} else {
				ips = append(ips, ep)
			}
		}
	}
	return ips, nil
}

// signAuthorizer implements cfgHandlers.SignAuthorizer.
type signAuthorizer struct {
	controller signControllerClient
	// localAddrs returns the host's own interface IPs; injected for tests.
	localAddrs func() []string
}

func newSignAuthorizer(controller *controllerclient.Client) signAuthorizer {
	return signAuthorizer{
		controller: controllerSignBackend{client: controller},
		localAddrs: hostInterfaceIPs,
	}
}

const signAuthzTimeout = 5 * time.Second

func (a signAuthorizer) AuthorizeSign(r *http.Request) error {
	peerIP := peerIPOf(r)

	// Class 2 — Day-1 joiner with a controller-issued join_id. Checked first:
	// valid for proxied and direct requests alike, since the join_id itself is
	// the credential (issued only after join-token verification).
	if joinID := strings.TrimSpace(r.URL.Query().Get("join_id")); joinID != "" {
		ctx, cancel := context.WithTimeout(r.Context(), signAuthzTimeout)
		defer cancel()
		planJSON, err := a.controller.GetJoinRequestStatusPlanJSON(ctx, joinID)
		if err != nil {
			return fmt.Errorf("join_id not accepted by controller: %v", err)
		}
		if len(planJSON) == 0 {
			// A join record without a stored plan was never authorized
			// (denied at preflight, or pre-v2 record) — not a credential.
			return errors.New("join_id has no authorized JoinPlan")
		}
		return nil
	}

	// IP-based classes require the TCP peer to be the real requester.
	if isProxiedRequest(r) {
		return errors.New("proxied requests must present a join_id")
	}
	if peerIP == "" {
		return errors.New("requester address unavailable")
	}

	// Class 1 — the local node itself.
	ip := net.ParseIP(peerIP)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	for _, local := range a.localAddrs() {
		if peerIP == local {
			return nil
		}
	}

	// Class 3 — a registered cluster member (renewal path).
	if a.controller == nil {
		return errors.New("controller unavailable for membership check")
	}
	ctx, cancel := context.WithTimeout(r.Context(), signAuthzTimeout)
	defer cancel()
	memberIPs, err := a.controller.ListNodeIPs(ctx)
	if err != nil {
		return fmt.Errorf("membership check failed (fail closed): %v", err)
	}
	for _, member := range memberIPs {
		if peerIP == member {
			return nil
		}
	}
	return fmt.Errorf("requester %s is not a cluster member and presented no join_id", peerIP)
}

// peerIPOf extracts the TCP peer IP from the request. It deliberately ignores
// X-Forwarded-For and similar headers — they are caller-controlled.
func peerIPOf(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}

// isProxiedRequest reports whether the request passed through an HTTP proxy
// (Envoy, a load balancer, …). For such requests the TCP peer address
// identifies the proxy, not the requester, so IP-based trust must not apply.
func isProxiedRequest(r *http.Request) bool {
	return r.Header.Get("X-Forwarded-For") != "" ||
		r.Header.Get("Forwarded") != "" ||
		r.Header.Get("X-Envoy-External-Address") != ""
}

// hostInterfaceIPs returns the IP addresses of all local interfaces.
func hostInterfaceIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var ips []string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips
}
