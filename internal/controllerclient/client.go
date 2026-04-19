package controllerclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	cluster_controllerpb "github.com/globulario/services/golang/cluster_controller/cluster_controllerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Well-known CA cert locations searched in order.
var caCertPaths = []string{
	"/var/lib/globular/pki/ca.crt",
	"/var/lib/globular/pki/ca.pem",
}

// Client wraps the ClusterController gRPC service.
type Client struct {
	addr string
}

// New returns a client that targets the specified controller address.
func New(addr string) *Client {
	return &Client{addr: strings.TrimSpace(addr)}
}

func (c *Client) dial(ctx context.Context) (*grpc.ClientConn, func(), error) {
	if c == nil || c.addr == "" {
		return nil, nil, fmt.Errorf("cluster controller address is empty")
	}
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

	// Try TLS with the Globular CA cert. Fall back to plaintext if no CA is
	// found (pre-TLS installations or test environments).
	// For loopback addresses, use TLS with InsecureSkipVerify since the
	// service cert may not include 127.0.0.1 in its SANs.
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	host, _, _ := strings.Cut(c.addr, ":")
	isLoopback := host == "127.0.0.1" || host == "::1" || host == "localhost"
	if isLoopback {
		creds = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // loopback only
		}))
	} else if tc := loadTLSCreds(); tc != nil {
		creds = grpc.WithTransportCredentials(tc)
	}

	conn, err := grpc.DialContext(dialCtx, c.addr, creds)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("dial %s: %w", c.addr, err)
	}
	return conn, func() {
		conn.Close()
		cancel()
	}, nil
}

// loadTLSCreds attempts to load the Globular CA cert for TLS connections.
// Returns nil if no CA cert is found (caller should fall back to plaintext).
func loadTLSCreds() credentials.TransportCredentials {
	for _, path := range caCertPaths {
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

// CreateJoinToken requests a join token from the controller.
func (c *Client) CreateJoinToken(ctx context.Context, expiresIn time.Duration) (*cluster_controllerpb.CreateJoinTokenResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	req := &cluster_controllerpb.CreateJoinTokenRequest{}
	if expiresIn > 0 {
		req.ExpiresAt = timestamppb.New(time.Now().Add(expiresIn))
	}
	return client.CreateJoinToken(ctx, req)
}

// ListJoinRequests returns pending join requests.
func (c *Client) ListJoinRequests(ctx context.Context) (*cluster_controllerpb.ListJoinRequestsResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	return client.ListJoinRequests(ctx, &cluster_controllerpb.ListJoinRequestsRequest{})
}

// ApproveJoin approves a join request with the provided profiles/metadata.
func (c *Client) ApproveJoin(ctx context.Context, nodeID string, profiles []string, metadata map[string]string) (*cluster_controllerpb.ApproveJoinResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	req := &cluster_controllerpb.ApproveJoinRequest{
		NodeId:   strings.TrimSpace(nodeID),
		Profiles: profiles,
		Metadata: metadata,
	}
	return client.ApproveJoin(ctx, req)
}

// ListNodes returns the nodes registered to the cluster.
func (c *Client) ListNodes(ctx context.Context) (*cluster_controllerpb.ListNodesResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	return client.ListNodes(ctx, &cluster_controllerpb.ListNodesRequest{})
}

// SetNodeProfiles updates the profiles assigned to a node.
func (c *Client) SetNodeProfiles(ctx context.Context, nodeID string, profiles []string) (*cluster_controllerpb.SetNodeProfilesResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	req := &cluster_controllerpb.SetNodeProfilesRequest{
		NodeId:   strings.TrimSpace(nodeID),
		Profiles: profiles,
	}
	return client.SetNodeProfiles(ctx, req)
}

func (c *Client) UpdateClusterNetwork(ctx context.Context, spec *cluster_controllerpb.ClusterNetworkSpec) (*cluster_controllerpb.UpdateClusterNetworkResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	return client.UpdateClusterNetwork(ctx, &cluster_controllerpb.UpdateClusterNetworkRequest{Spec: spec})
}

// RemoveNode removes a node from the cluster.
func (c *Client) RemoveNode(ctx context.Context, nodeID string, force, drain bool) (*cluster_controllerpb.RemoveNodeResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	req := &cluster_controllerpb.RemoveNodeRequest{
		NodeId: strings.TrimSpace(nodeID),
		Force:  force,
		Drain:  drain,
	}
	return client.RemoveNode(ctx, req)
}

// GetClusterHealth returns the overall health status of the cluster.
func (c *Client) GetClusterHealth(ctx context.Context) (*cluster_controllerpb.GetClusterHealthResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	return client.GetClusterHealth(ctx, &cluster_controllerpb.GetClusterHealthRequest{})
}

// GetClusterNetwork returns the cluster network configuration.
func (c *Client) GetClusterNetwork(ctx context.Context) (*cluster_controllerpb.ClusterNetwork, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewResourcesServiceClient(conn)
	return client.GetClusterNetwork(ctx, &cluster_controllerpb.GetClusterNetworkRequest{})
}

// Address returns the configured controller address.
func (c *Client) Address() string {
	if c == nil {
		return ""
	}
	return c.addr
}
