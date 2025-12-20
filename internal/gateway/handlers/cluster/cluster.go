package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	agentclient "github.com/globulario/Globular/internal/agentclient"
	"github.com/globulario/Globular/internal/controllerclient"
)

// HandlerDeps contains dependencies that cluster handlers need.
type HandlerDeps struct {
	Controller    *controllerclient.Client
	NodeAgentAddr string
}

// Deps groups HTTP handlers to register.
type Deps struct {
	JoinToken   http.Handler
	Nodes       http.Handler
	NodeActions http.Handler
}

// Mount registers the cluster-related routes.
func Mount(mux *http.ServeMux, d Deps) {
	if d.JoinToken != nil {
		mux.Handle("/api/cluster/join-token", d.JoinToken)
	}
	if d.Nodes != nil {
		mux.Handle("/api/cluster/nodes", d.Nodes)
	}
	if d.NodeActions != nil {
		mux.Handle("/api/cluster/nodes/", d.NodeActions)
	}
}

// NewJoinTokenHandler creates the handler that returns a cluster join token.
func NewJoinTokenHandler(deps HandlerDeps) http.Handler {
	const ttl = time.Hour
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if deps.Controller == nil || deps.Controller.Address() == "" {
			http.Error(w, "cluster controller not configured", http.StatusServiceUnavailable)
			return
		}
		resp, err := deps.Controller.CreateJoinToken(r.Context(), ttl)
		if err != nil {
			http.Error(w, fmt.Sprintf("create join token: %v", err), http.StatusServiceUnavailable)
			return
		}
		respondJSON(w, http.StatusOK, resp)
	})
}

// NewNodesHandler lists registered nodes.
func NewNodesHandler(deps HandlerDeps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if deps.Controller == nil || deps.Controller.Address() == "" {
			http.Error(w, "cluster controller not configured", http.StatusServiceUnavailable)
			return
		}
		resp, err := deps.Controller.ListNodes(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("list nodes: %v", err), http.StatusServiceUnavailable)
			return
		}
		respondJSON(w, http.StatusOK, resp)
	})
}

// NewNodeActionsHandler handles node-specific subpaths (profiles, plan apply).
func NewNodeActionsHandler(deps HandlerDeps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if deps.Controller == nil || deps.Controller.Address() == "" {
			http.Error(w, "cluster controller not configured", http.StatusServiceUnavailable)
			return
		}

		nodeID, action := parseNodeRoute(r.URL.Path)
		if nodeID == "" || action == "" {
			http.NotFound(w, r)
			return
		}

		switch action {
		case "profiles":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handleSetProfiles(w, r, deps.Controller, nodeID)
			return
		case "plan/apply":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			handleApplyPlan(w, r, deps, nodeID)
			return
		default:
			http.NotFound(w, r)
		}
	})
}

func handleSetProfiles(w http.ResponseWriter, r *http.Request, controller *controllerclient.Client, nodeID string) {
	var req struct {
		Profiles []string `json:"profiles"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("decode request: %v", err), http.StatusBadRequest)
		return
	}
	if len(req.Profiles) == 0 {
		http.Error(w, "profiles list required", http.StatusBadRequest)
		return
	}
	resp, err := controller.SetNodeProfiles(r.Context(), nodeID, req.Profiles)
	if err != nil {
		http.Error(w, fmt.Sprintf("set profiles: %v", err), http.StatusServiceUnavailable)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func handleApplyPlan(w http.ResponseWriter, r *http.Request, deps HandlerDeps, nodeID string) {
	ctx := r.Context()
	plan, err := deps.Controller.GetNodePlan(ctx, nodeID)
	if err != nil {
		http.Error(w, fmt.Sprintf("get plan: %v", err), http.StatusServiceUnavailable)
		return
	}
	if plan == nil || len(plan.GetUnitActions()) == 0 {
		respondJSON(w, http.StatusOK, map[string]string{"status": "plan empty"})
		return
	}
	if strings.TrimSpace(deps.NodeAgentAddr) == "" {
		http.Error(w, "node agent address not configured", http.StatusServiceUnavailable)
		return
	}
	if err := agentclient.ApplyPlan(ctx, deps.NodeAgentAddr, plan); err != nil {
		http.Error(w, fmt.Sprintf("apply plan: %v", err), http.StatusServiceUnavailable)
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "plan applied"})
}

func parseNodeRoute(path string) (string, string) {
	const prefix = "/api/cluster/nodes/"
	if !strings.HasPrefix(path, prefix) {
		return "", ""
	}
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return "", ""
	}
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], strings.Join(parts[1:], "/")
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
