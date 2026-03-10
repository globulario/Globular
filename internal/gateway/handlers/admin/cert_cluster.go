package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/globulario/Globular/internal/controllerclient"
)

// ── Cluster certificate types ───────────────────────────────────────────────

// ClusterCertSummary is the rollup across all nodes.
type ClusterCertSummary struct {
	TotalNodes        int `json:"totalNodes"`
	HealthyNodes      int `json:"healthyNodes"`
	WarningNodes      int `json:"warningNodes"`
	ErrorNodes        int `json:"errorNodes"`
	UnreachableNodes  int `json:"unreachableNodes"`
	ExpiringSoonCount int `json:"expiringSoonCount"`
	ExpiredCount      int `json:"expiredCount"`
}

// ClusterTrustDrift captures cluster-wide inconsistencies.
type ClusterTrustDrift struct {
	InternalSANMismatch  bool     `json:"internalSANMismatch"`
	PublicDomainMismatch bool     `json:"publicDomainMismatch"`
	EnvoyTLSDrift        bool     `json:"envoyTLSDrift"`
	NodesOutOfPolicy     []string `json:"nodesOutOfPolicy"`
}

// NodeCertPKISummary is the per-node internal PKI state.
type NodeCertPKISummary struct {
	CAStatus          string `json:"caStatus"`
	ServiceCertStatus string `json:"serviceCertStatus"`
	DaysUntilExpiry   *int   `json:"daysUntilExpiry,omitempty"`
}

// NodeCertPublicSummary is the per-node public TLS state.
type NodeCertPublicSummary struct {
	Enabled         bool   `json:"enabled"`
	CertStatus      string `json:"certStatus"`
	DaysUntilExpiry *int   `json:"daysUntilExpiry,omitempty"`
	Domain          string `json:"domain,omitempty"`
}

// NodeCertEnvoySummary is the per-node Envoy TLS state.
type NodeCertEnvoySummary struct {
	Status         string `json:"status"`
	ListenerIssues int    `json:"listenerIssues"`
	UpstreamIssues int    `json:"upstreamIssues"`
}

// ClusterNodeCertStatus is the full per-node certificate status.
type ClusterNodeCertStatus struct {
	NodeID      string                 `json:"nodeId"`
	Address     string                 `json:"address"`
	Status      string                 `json:"status"` // healthy | warning | error | unreachable
	InternalPKI *NodeCertPKISummary    `json:"internalPKI,omitempty"`
	PublicTLS   *NodeCertPublicSummary `json:"publicTLS,omitempty"`
	Envoy       *NodeCertEnvoySummary  `json:"envoy,omitempty"`
	Warnings    []Warning              `json:"warnings"`
}

// ClusterCertOverview is the top-level response for GET /admin/certificates/cluster.
type ClusterCertOverview struct {
	Summary ClusterCertSummary      `json:"summary"`
	Drift   ClusterTrustDrift       `json:"drift"`
	Nodes   []ClusterNodeCertStatus `json:"nodes"`
}

// ── Handler ─────────────────────────────────────────────────────────────────

// NewClusterCertificatesHandler returns a GET handler for /admin/certificates/cluster.
func NewClusterCertificatesHandler(controller *controllerclient.Client, gatewayPort int) http.HandlerFunc {
	if gatewayPort == 0 {
		gatewayPort = 8080
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if controller == nil || controller.Address() == "" {
			writeActionResponse(w, http.StatusServiceUnavailable, false,
				"Cluster controller not configured")
			return
		}

		ctx := r.Context()
		nodes := fetchClusterCertStates(ctx, controller, gatewayPort)

		overview := ClusterCertOverview{
			Summary: computeClusterSummary(nodes),
			Drift:   computeTrustDrift(nodes),
			Nodes:   nodes,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(overview)
	}
}

// ── Cluster data gathering ──────────────────────────────────────────────────

// fetchClusterCertStates queries all known cluster peers in parallel.
func fetchClusterCertStates(ctx context.Context, controller *controllerclient.Client, gatewayPort int) []ClusterNodeCertStatus {
	resp, err := controller.ListNodes(ctx)
	if err != nil {
		return []ClusterNodeCertStatus{{
			NodeID:  "(controller)",
			Address: controller.Address(),
			Status:  "unreachable",
			Warnings: []Warning{{
				Severity: "error",
				Message:  fmt.Sprintf("Cannot list cluster nodes: %v", err),
			}},
		}}
	}

	records := resp.GetNodes()
	if len(records) == 0 {
		return nil
	}

	results := make([]ClusterNodeCertStatus, len(records))
	var wg sync.WaitGroup

	for i, rec := range records {
		wg.Add(1)
		go func(idx int, nodeID, address, fqdn string) {
			defer wg.Done()
			results[idx] = fetchNodeCertStatus(ctx, nodeID, address, fqdn, gatewayPort)
		}(i, rec.GetNodeId(), rec.GetAgentEndpoint(), rec.GetAdvertiseFqdn())
	}

	wg.Wait()
	return results
}

// fetchNodeCertStatus fetches /admin/certificates from a single node.
func fetchNodeCertStatus(ctx context.Context, nodeID, agentEndpoint, fqdn string, gatewayPort int) ClusterNodeCertStatus {
	// Determine the HTTP address to reach this node's gateway
	host := resolveNodeHost(agentEndpoint, fqdn)
	if host == "" {
		return ClusterNodeCertStatus{
			NodeID:  nodeID,
			Address: fqdn,
			Status:  "unreachable",
			Warnings: []Warning{{
				Severity: "error",
				Message:  "No reachable address for node",
			}},
		}
	}

	url := fmt.Sprintf("http://%s:%d/admin/certificates", host, gatewayPort)

	fetchCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(fetchCtx, http.MethodGet, url, nil)
	if err != nil {
		return unreachableNode(nodeID, fqdn, fmt.Sprintf("build request: %v", err))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return unreachableNode(nodeID, fqdn, fmt.Sprintf("fetch failed: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return unreachableNode(nodeID, fqdn, fmt.Sprintf("HTTP %d from node", resp.StatusCode))
	}

	var overview CertOverview
	if err := json.NewDecoder(resp.Body).Decode(&overview); err != nil {
		return unreachableNode(nodeID, fqdn, fmt.Sprintf("decode failed: %v", err))
	}

	return normalizeNodeCertStatus(nodeID, fqdn, &overview)
}

// resolveNodeHost extracts a usable hostname/IP from agent endpoint or FQDN.
func resolveNodeHost(agentEndpoint, fqdn string) string {
	// Try FQDN first (more stable for DNS)
	if fqdn != "" {
		return fqdn
	}
	// Agent endpoint is typically "ip:port" — extract IP
	if agentEndpoint != "" {
		for i := len(agentEndpoint) - 1; i >= 0; i-- {
			if agentEndpoint[i] == ':' {
				return agentEndpoint[:i]
			}
		}
		return agentEndpoint
	}
	return ""
}

func unreachableNode(nodeID, address, msg string) ClusterNodeCertStatus {
	return ClusterNodeCertStatus{
		NodeID:  nodeID,
		Address: address,
		Status:  "unreachable",
		Warnings: []Warning{{
			Severity: "error",
			Message:  msg,
		}},
	}
}

// ── Normalization ───────────────────────────────────────────────────────────

// normalizeNodeCertStatus converts a full CertOverview into a compact node summary.
func normalizeNodeCertStatus(nodeID, address string, overview *CertOverview) ClusterNodeCertStatus {
	node := ClusterNodeCertStatus{
		NodeID:   nodeID,
		Address:  address,
		Warnings: overview.Warnings,
	}

	// Internal PKI summary
	pki := &NodeCertPKISummary{
		CAStatus:          overview.InternalPKI.CA.Status,
		ServiceCertStatus: overview.InternalPKI.ServiceCert.Status,
	}
	if overview.InternalPKI.ServiceCert.DaysUntilExpiry != nil {
		pki.DaysUntilExpiry = overview.InternalPKI.ServiceCert.DaysUntilExpiry
	}
	node.InternalPKI = pki

	// Public TLS summary
	pub := &NodeCertPublicSummary{
		Enabled: overview.PublicTLS.Protocol == "https",
		Domain:  overview.PublicTLS.Domain,
	}
	if overview.PublicTLS.LeafCert != nil {
		pub.CertStatus = overview.PublicTLS.LeafCert.Status
		pub.DaysUntilExpiry = overview.PublicTLS.LeafCert.DaysUntilExpiry
	} else if len(overview.PublicTLS.ExternalDomains) > 0 {
		// Use first external domain cert
		first := overview.PublicTLS.ExternalDomains[0]
		if first.LeafCert != nil {
			pub.CertStatus = first.LeafCert.Status
			pub.DaysUntilExpiry = first.LeafCert.DaysUntilExpiry
			pub.Domain = first.FQDN
		}
	} else if !pub.Enabled {
		pub.CertStatus = "not_applicable"
	} else {
		pub.CertStatus = "missing"
	}
	node.PublicTLS = pub

	// Envoy summary
	envoy := &NodeCertEnvoySummary{Status: "ok"}
	for _, l := range overview.Envoy.Listeners {
		if l.Status != "ok" {
			envoy.ListenerIssues++
		}
	}
	for _, u := range overview.Envoy.Upstreams {
		if u.Status != "ok" && u.Status != "no_tls" {
			envoy.UpstreamIssues++
		}
	}
	if envoy.ListenerIssues > 0 || envoy.UpstreamIssues > 0 {
		envoy.Status = "warning"
	}
	if !overview.Envoy.XDSClient.Exists {
		envoy.Status = "error"
	}
	node.Envoy = envoy

	// Compute overall node status
	node.Status = computeNodeStatus(node)

	return node
}

// computeNodeStatus derives the overall node health from its subsections.
func computeNodeStatus(node ClusterNodeCertStatus) string {
	hasError := false
	hasWarning := false

	if node.InternalPKI != nil {
		switch node.InternalPKI.CAStatus {
		case "expired", "missing", "parse_error":
			hasError = true
		case "expiring":
			hasWarning = true
		}
		switch node.InternalPKI.ServiceCertStatus {
		case "expired", "missing", "parse_error":
			hasError = true
		case "expiring":
			hasWarning = true
		}
	}

	if node.PublicTLS != nil && node.PublicTLS.Enabled {
		switch node.PublicTLS.CertStatus {
		case "expired", "missing", "parse_error":
			hasError = true
		case "expiring":
			hasWarning = true
		}
	}

	if node.Envoy != nil {
		switch node.Envoy.Status {
		case "error":
			hasError = true
		case "warning":
			hasWarning = true
		}
	}

	// Check warnings from the node
	for _, w := range node.Warnings {
		if w.Severity == "error" {
			hasError = true
		} else if w.Severity == "warning" {
			hasWarning = true
		}
	}

	if hasError {
		return "error"
	}
	if hasWarning {
		return "warning"
	}
	return "healthy"
}

// ── Summary computation ─────────────────────────────────────────────────────

func computeClusterSummary(nodes []ClusterNodeCertStatus) ClusterCertSummary {
	s := ClusterCertSummary{TotalNodes: len(nodes)}

	for _, n := range nodes {
		switch n.Status {
		case "healthy":
			s.HealthyNodes++
		case "warning":
			s.WarningNodes++
		case "error":
			s.ErrorNodes++
		case "unreachable":
			s.UnreachableNodes++
		}

		// Count expiring/expired certs
		if n.InternalPKI != nil && n.InternalPKI.DaysUntilExpiry != nil {
			days := *n.InternalPKI.DaysUntilExpiry
			if days < 0 {
				s.ExpiredCount++
			} else if days < expiryWarningDays {
				s.ExpiringSoonCount++
			}
		}
		if n.PublicTLS != nil && n.PublicTLS.DaysUntilExpiry != nil {
			days := *n.PublicTLS.DaysUntilExpiry
			if days < 0 {
				s.ExpiredCount++
			} else if days < expiryWarningDays {
				s.ExpiringSoonCount++
			}
		}
	}

	return s
}

// ── Trust drift detection ───────────────────────────────────────────────────

func computeTrustDrift(nodes []ClusterNodeCertStatus) ClusterTrustDrift {
	drift := ClusterTrustDrift{}

	// Collect baseline values from the first healthy node
	var baselinePKIStatus string
	var baselinePublicDomain string
	var baselineEnvoyStatus string
	baselineSet := false

	var outOfPolicy []string

	for _, n := range nodes {
		if n.Status == "unreachable" {
			continue
		}

		if !baselineSet {
			if n.InternalPKI != nil {
				baselinePKIStatus = n.InternalPKI.ServiceCertStatus
			}
			if n.PublicTLS != nil && n.PublicTLS.Enabled {
				baselinePublicDomain = n.PublicTLS.Domain
			}
			if n.Envoy != nil {
				baselineEnvoyStatus = n.Envoy.Status
			}
			baselineSet = true
			continue
		}

		nodeOutOfPolicy := false

		// Internal PKI drift: different service cert statuses across nodes
		if n.InternalPKI != nil && baselinePKIStatus != "" {
			if n.InternalPKI.ServiceCertStatus != baselinePKIStatus {
				drift.InternalSANMismatch = true
				nodeOutOfPolicy = true
			}
		}

		// Public domain drift: different domains across public-facing nodes
		if n.PublicTLS != nil && n.PublicTLS.Enabled && baselinePublicDomain != "" {
			if n.PublicTLS.Domain != baselinePublicDomain {
				drift.PublicDomainMismatch = true
				nodeOutOfPolicy = true
			}
		}

		// Envoy drift: different TLS health across nodes
		if n.Envoy != nil && baselineEnvoyStatus != "" {
			if n.Envoy.Status != baselineEnvoyStatus {
				drift.EnvoyTLSDrift = true
				nodeOutOfPolicy = true
			}
		}

		if nodeOutOfPolicy {
			outOfPolicy = append(outOfPolicy, n.NodeID)
		}
	}

	drift.NodesOutOfPolicy = outOfPolicy
	return drift
}
