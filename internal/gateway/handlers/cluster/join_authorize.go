package cluster

import (
	"encoding/json"
	"net/http"
	"strings"

	cluster_controllerpb "github.com/globulario/services/golang/cluster_controller/cluster_controllerpb"
)

// joinAuthorizeRequest is the JSON body accepted by POST /join/authorize.
// It mirrors the controller's JoinAuthorizationRequest for JSON-over-HTTP.
type joinAuthorizeRequest struct {
	JoinToken        string             `json:"join_token"`
	Identity         joinAuthorizeIdent `json:"identity"`
	Labels           map[string]string  `json:"labels,omitempty"`
	CPUCount         uint32             `json:"cpu_count,omitempty"`
	RAMBytes         uint64             `json:"ram_bytes,omitempty"`
	DiskBytes        uint64             `json:"disk_bytes,omitempty"`
	InstallerVersion string             `json:"installer_version,omitempty"`
	ClusterID        string             `json:"cluster_id,omitempty"`
	Nonce            string             `json:"nonce,omitempty"`
}

type joinAuthorizeIdent struct {
	Hostname string   `json:"hostname"`
	IPs      []string `json:"ips,omitempty"`
}

// joinAuthorizeResponse is the JSON body returned by POST /join/authorize.
// It carries the controller's decision and, on success, the raw plan JSON.
type joinAuthorizeResponse struct {
	// Allowed is true when the controller issued a valid JoinPlan.
	Allowed bool `json:"allowed"`
	// DeniedReason is set when Allowed=false.
	DeniedReason string `json:"denied_reason,omitempty"`
	// JoinID is the unique identifier for this authorization.
	JoinID string `json:"join_id,omitempty"`
	// Plan is the raw JoinPlan JSON from the controller. The installer must
	// validate Plan before executing any cluster-affecting step.
	Plan json.RawMessage `json:"plan,omitempty"`
	// ControllerGeneration is the controller state generation at issuance.
	ControllerGeneration int64 `json:"controller_generation,omitempty"`
}

// NewJoinAuthorizeHandler returns an HTTP handler for POST /join/authorize.
//
// The handler is a pure courier: it unmarshals the installer's request, forwards
// it to the controller's RequestJoinAuthorization RPC, and returns the response.
// It does NOT assign profiles, write etcd membership, or make any cluster decision.
//
// The controller is the sole authority for:
//   - profile assignment
//   - etcd join intent
//   - assigned node_id
//   - signing the JoinPlan
func NewJoinAuthorizeHandler(deps HandlerDeps) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed — POST required", http.StatusMethodNotAllowed)
			return
		}
		if deps.Controller == nil || deps.Controller.Address() == "" {
			http.Error(w, "cluster controller not configured", http.StatusServiceUnavailable)
			return
		}

		var req joinAuthorizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		token := strings.TrimSpace(req.JoinToken)
		if token == "" {
			respondJSON(w, http.StatusBadRequest, joinAuthorizeResponse{
				Allowed:      false,
				DeniedReason: "join_token is required",
			})
			return
		}

		// Build proto request — the gateway adds no fields of its own.
		protoReq := &cluster_controllerpb.JoinAuthorizationRequest{
			JoinToken: token,
			Identity: &cluster_controllerpb.NodeIdentity{
				Hostname: strings.TrimSpace(req.Identity.Hostname),
				Ips:      append([]string(nil), req.Identity.IPs...),
			},
			Labels:           req.Labels,
			InstallerVersion: req.InstallerVersion,
			ClusterId:        strings.TrimSpace(req.ClusterID),
			Nonce:            req.Nonce,
		}
		if req.CPUCount > 0 || req.RAMBytes > 0 || req.DiskBytes > 0 {
			protoReq.Capabilities = &cluster_controllerpb.NodeCapabilities{
				CpuCount:  req.CPUCount,
				RamBytes:  req.RAMBytes,
				DiskBytes: req.DiskBytes,
			}
		}

		resp, err := deps.Controller.RequestJoinAuthorization(r.Context(), protoReq)
		if err != nil {
			http.Error(w, "controller error: "+err.Error(), http.StatusBadGateway)
			return
		}

		if !resp.GetAllowed() {
			respondJSON(w, http.StatusForbidden, joinAuthorizeResponse{
				Allowed:      false,
				DeniedReason: resp.GetDeniedReason(),
				JoinID:       resp.GetJoinId(),
			})
			return
		}

		respondJSON(w, http.StatusOK, joinAuthorizeResponse{
			Allowed:              true,
			JoinID:               resp.GetJoinId(),
			Plan:                 json.RawMessage(resp.GetPlanJson()),
			ControllerGeneration: resp.GetControllerGeneration(),
		})
	})
}
