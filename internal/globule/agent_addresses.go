package globule

import (
	"os"
	"strings"
)

const (
	defaultNodeAgentAddr  = "127.0.0.1:11000"
	defaultControllerAddr = "127.0.0.1:12000"
	nodeAgentEnvKey       = "NODE_AGENT_ADDR"
	controllerAgentEnvKey = "CLUSTER_CONTROLLER_ADDR"
)

// NodeAgentAddress returns the configured node agent endpoint.
func NodeAgentAddress() string {
	if v := strings.TrimSpace(os.Getenv(nodeAgentEnvKey)); v != "" {
		return v
	}
	return defaultNodeAgentAddr
}

// ControllerAddress returns the configured cluster controller endpoint.
func ControllerAddress() string {
	if v := strings.TrimSpace(os.Getenv(controllerAgentEnvKey)); v != "" {
		return v
	}
	return defaultControllerAddr
}
