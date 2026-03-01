package agentclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	cluster_controllerpb "github.com/globulario/services/golang/cluster_controller/cluster_controllerpb"
	node_agentpb "github.com/globulario/services/golang/node_agent/node_agentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ApplySingleUnitAction sends one unit action plan to the NodeAgent and waits for completion.
func ApplySingleUnitAction(ctx context.Context, nodeAgentAddr, unit, action string) error {
	plan := &cluster_controllerpb.NodePlan{
		UnitActions: []*cluster_controllerpb.UnitAction{{UnitName: strings.TrimSpace(unit), Action: strings.TrimSpace(action)}},
	}
	return applyPlan(ctx, nodeAgentAddr, plan)
}

// ApplyPlan sends the provided NodePlan to the NodeAgent.
func ApplyPlan(ctx context.Context, nodeAgentAddr string, plan *cluster_controllerpb.NodePlan) error {
	return applyPlan(ctx, nodeAgentAddr, plan)
}

func applyPlan(ctx context.Context, nodeAgentAddr string, plan *cluster_controllerpb.NodePlan) error {
	nodeAgentAddr = strings.TrimSpace(nodeAgentAddr)
	if nodeAgentAddr == "" {
		return fmt.Errorf("node-agent address is empty")
	}
	if plan == nil || (len(plan.GetUnitActions()) == 0 && len(plan.GetRenderedConfig()) == 0) {
		return fmt.Errorf("plan must include at least one action or rendered config change")
	}

	dialCtx, dialCancel := context.WithTimeout(ctx, 5*time.Second)
	defer dialCancel()
	conn, err := grpc.DialContext(dialCtx, nodeAgentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial node-agent %s: %w", nodeAgentAddr, err)
	}
	defer conn.Close()

	client := node_agentpb.NewNodeAgentServiceClient(conn)

	respCtx, respCancel := context.WithTimeout(ctx, 10*time.Second)
	defer respCancel()
	resp, err := client.ApplyPlan(respCtx, &node_agentpb.ApplyPlanRequest{Plan: plan})
	if err != nil {
		return fmt.Errorf("apply plan: %w", err)
	}
	opID := strings.TrimSpace(resp.GetOperationId())
	if opID == "" {
		return fmt.Errorf("node-agent returned empty operation_id")
	}

	streamCtx, streamCancel := context.WithTimeout(ctx, 2*time.Minute)
	defer streamCancel()
	stream, err := client.WatchOperation(streamCtx, &node_agentpb.WatchOperationRequest{OperationId: opID})
	if err != nil {
		return fmt.Errorf("watch operation: %w", err)
	}

	for {
		evt, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("watch operation recv: %w", err)
		}
		if evt.GetDone() {
			if msg := strings.TrimSpace(evt.GetError()); msg != "" {
				return fmt.Errorf("operation failed: %s", msg)
			}
			return nil
		}
	}
}
