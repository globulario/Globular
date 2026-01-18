package controllerclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	clustercontrollerpb "github.com/globulario/services/golang/clustercontroller/clustercontrollerpb"
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

func (c *Client) dial(ctx context.Context) (clustercontrollerpb.ClusterControllerServiceClient, func(), error) {
	if c == nil || c.addr == "" {
		return nil, nil, fmt.Errorf("cluster controller address is empty")
	}
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	conn, err := grpc.DialContext(dialCtx, c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("dial %s: %w", c.addr, err)
	}
	return clustercontrollerpb.NewClusterControllerServiceClient(conn), func() {
		conn.Close()
		cancel()
	}, nil
}

// CreateJoinToken requests a join token from the controller.
func (c *Client) CreateJoinToken(ctx context.Context, expiresIn time.Duration) (*clustercontrollerpb.CreateJoinTokenResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	req := &clustercontrollerpb.CreateJoinTokenRequest{}
	if expiresIn > 0 {
		req.ExpiresAt = timestamppb.New(time.Now().Add(expiresIn))
	}
	return client.CreateJoinToken(ctx, req)
}

// ListJoinRequests returns pending join requests.
func (c *Client) ListJoinRequests(ctx context.Context) (*clustercontrollerpb.ListJoinRequestsResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	return client.ListJoinRequests(ctx, &clustercontrollerpb.ListJoinRequestsRequest{})
}

// ApproveJoin approves a join request with the provided profiles/metadata.
func (c *Client) ApproveJoin(ctx context.Context, nodeID string, profiles []string, metadata map[string]string) (*clustercontrollerpb.ApproveJoinResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	req := &clustercontrollerpb.ApproveJoinRequest{
		NodeId:   strings.TrimSpace(nodeID),
		Profiles: profiles,
		Metadata: metadata,
	}
	return client.ApproveJoin(ctx, req)
}

// ListNodes returns the nodes registered to the cluster.
func (c *Client) ListNodes(ctx context.Context) (*clustercontrollerpb.ListNodesResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	return client.ListNodes(ctx, &clustercontrollerpb.ListNodesRequest{})
}

// SetNodeProfiles updates the profiles assigned to a node.
func (c *Client) SetNodeProfiles(ctx context.Context, nodeID string, profiles []string) (*clustercontrollerpb.SetNodeProfilesResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	req := &clustercontrollerpb.SetNodeProfilesRequest{
		NodeId:   strings.TrimSpace(nodeID),
		Profiles: profiles,
	}
	return client.SetNodeProfiles(ctx, req)
}

// GetNodePlan requests the latest plan for the given node.
func (c *Client) GetNodePlan(ctx context.Context, nodeID string) (*clustercontrollerpb.NodePlan, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	resp, err := client.GetNodePlan(ctx, &clustercontrollerpb.GetNodePlanRequest{NodeId: strings.TrimSpace(nodeID)})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("controller returned empty plan")
	}
	return resp.GetPlan(), nil
}

func (c *Client) UpdateClusterNetwork(ctx context.Context, spec *clustercontrollerpb.ClusterNetworkSpec) (*clustercontrollerpb.UpdateClusterNetworkResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	return client.UpdateClusterNetwork(ctx, &clustercontrollerpb.UpdateClusterNetworkRequest{Spec: spec})
}

func (c *Client) ApplyNodePlan(ctx context.Context, nodeID string) (*clustercontrollerpb.ApplyNodePlanResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	return client.ApplyNodePlan(ctx, &clustercontrollerpb.ApplyNodePlanRequest{NodeId: strings.TrimSpace(nodeID)})
}

// RemoveNode removes a node from the cluster.
func (c *Client) RemoveNode(ctx context.Context, nodeID string, force, drain bool) (*clustercontrollerpb.RemoveNodeResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()

	req := &clustercontrollerpb.RemoveNodeRequest{
		NodeId: strings.TrimSpace(nodeID),
		Force:  force,
		Drain:  drain,
	}
	return client.RemoveNode(ctx, req)
}

// GetClusterHealth returns the overall health status of the cluster.
func (c *Client) GetClusterHealth(ctx context.Context) (*clustercontrollerpb.GetClusterHealthResponse, error) {
	client, closeFn, err := c.dial(ctx)
	if err != nil {
		return nil, err
	}
	defer closeFn()
	return client.GetClusterHealth(ctx, &clustercontrollerpb.GetClusterHealthRequest{})
}

// Address returns the configured controller address.
func (c *Client) Address() string {
	if c == nil {
		return ""
	}
	return c.addr
}
