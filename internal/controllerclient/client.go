package controllerclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	cluster_controllerpb "github.com/globulario/services/golang/cluster_controller/cluster_controllerpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
	conn, err := grpc.DialContext(dialCtx, c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("dial %s: %w", c.addr, err)
	}
	return conn, func() {
		conn.Close()
		cancel()
	}, nil
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

// GetNodePlan requests the latest plan for the given node.
func (c *Client) GetNodePlan(ctx context.Context, nodeID string) (*cluster_controllerpb.NodePlan, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	resp, err := client.GetNodePlan(ctx, &cluster_controllerpb.GetNodePlanRequest{NodeId: strings.TrimSpace(nodeID)})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("controller returned empty plan")
	}
	return resp.GetPlan(), nil
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

func (c *Client) ApplyNodePlan(ctx context.Context, nodeID string) (*cluster_controllerpb.ApplyNodePlanResponse, error) {
	conn, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	client := cluster_controllerpb.NewClusterControllerServiceClient(conn)
	return client.ApplyNodePlan(ctx, &cluster_controllerpb.ApplyNodePlanRequest{NodeId: strings.TrimSpace(nodeID)})
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
