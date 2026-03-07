package admin

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ── JSON response types ─────────────────────────────────────────────────────

// ServicesResponse is the top-level response for GET /admin/metrics/services.
type ServicesResponse struct {
	NowUnixMs  int64                   `json:"now_unix_ms"`
	Range      string                  `json:"range"`
	Prometheus PromStatus              `json:"prometheus"`
	Thresholds SvcThresholds           `json:"thresholds"`
	Groups     []ServiceGroup          `json:"groups"`
	Summary    ServicesSummary         `json:"summary"`
	Infra      map[string]*InfraDetail `json:"infra,omitempty"`
}

// PromStatus reports whether Prometheus is reachable.
type PromStatus struct {
	Connected bool   `json:"connected"`
	Addr      string `json:"addr"`
}

// SvcThresholds are the thresholds used to derive health.
type SvcThresholds struct {
	CPUWarn float64 `json:"cpu_warn_pct"`
	CPUCrit float64 `json:"cpu_crit_pct"`
	MemWarn float64 `json:"mem_warn_pct"`
	MemCrit float64 `json:"mem_crit_pct"`
}

// ServiceGroup is a named category of services.
type ServiceGroup struct {
	Category string            `json:"category"`
	Services []ServiceInstance `json:"services"`
}

// ServiceInstance is a single service with runtime metrics and derived health.
type ServiceInstance struct {
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	ID            string      `json:"id"`
	Version       string      `json:"version"`
	State         string      `json:"state"`
	Port          int         `json:"port"`
	Category      string      `json:"category"`
	Node          string      `json:"node"`
	DerivedStatus string      `json:"derived_status"` // healthy | degraded | critical | unknown
	Reasons       []string    `json:"reasons"`
	Runtime       *SvcRuntime `json:"runtime,omitempty"`
	GRPCHealth    *GRPCHealth `json:"grpc_health"`
}

// SvcRuntime holds per-service Prometheus metrics.
type SvcRuntime struct {
	CPUPct       float64 `json:"cpu_pct"`
	MemoryBytes  float64 `json:"memory_bytes"`
	UptimeSec    float64 `json:"uptime_sec"`
	ReqRate      float64 `json:"req_rate"`
	ErrRate      float64 `json:"err_rate"`
	LatencyP50Ms float64 `json:"latency_p50_ms"`
	LatencyP95Ms float64 `json:"latency_p95_ms"`
	Goroutines   float64 `json:"goroutines"`
	HeapBytes    float64 `json:"heap_bytes"`
	OpenFDs      float64 `json:"open_fds"`
	MaxFDs       float64 `json:"max_fds"`
	MsgRecvRate  float64 `json:"msg_recv_rate"`
	MsgSentRate  float64 `json:"msg_sent_rate"`
}

// GRPCHealth is a placeholder for Phase 2 gRPC health checks.
type GRPCHealth struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

// InfraDetail holds infrastructure-specific metrics (etcd, envoy, node).
type InfraDetail struct {
	// etcd
	EtcdIsLeader    *bool    `json:"etcd_is_leader,omitempty"`
	EtcdDBSizeBytes *float64 `json:"etcd_db_size_bytes,omitempty"`
	EtcdTotalKeys   *float64 `json:"etcd_total_keys,omitempty"`
	// envoy
	EnvoyActiveConns *float64 `json:"envoy_active_conns,omitempty"`
	EnvoyRPS         *float64 `json:"envoy_rps,omitempty"`
	EnvoyHTTP5xx     *float64 `json:"envoy_http_5xx,omitempty"`
	// node
	NodeLoad1          *float64 `json:"node_load1,omitempty"`
	NodeLoad5          *float64 `json:"node_load5,omitempty"`
	NodeMemAvailBytes  *float64 `json:"node_mem_avail_bytes,omitempty"`
	NodeMemTotalBytes  *float64 `json:"node_mem_total_bytes,omitempty"`
	NodeNetRxBytesRate *float64 `json:"node_net_rx_rate,omitempty"`
	NodeNetTxBytesRate *float64 `json:"node_net_tx_rate,omitempty"`
}

// ServicesSummary is a rollup of service health counts.
type ServicesSummary struct {
	Total    int `json:"total"`
	Healthy  int `json:"healthy"`
	Degraded int `json:"degraded"`
	Critical int `json:"critical"`
	Unknown  int `json:"unknown"`
}

// ── Category sets (match metricsNormalizer.ts) ──────────────────────────────

var (
	coreSet  = newSet("gateway", "discovery", "repository", "resource", "rbac", "authentication", "event", "log", "file")
	infraSet = newSet("etcd", "envoy", "scylla", "minio", "prometheus", "dns")
	mediaSet = newSet("media", "title", "torrent", "search")
)

func categorize(baseName string) string {
	if coreSet[baseName] {
		return "Core"
	}
	if infraSet[baseName] {
		return "Infrastructure"
	}
	if mediaSet[baseName] {
		return "Media"
	}
	return "Other"
}

func newSet(items ...string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[s] = true
	}
	return m
}

// ── Prometheus metric collection ────────────────────────────────────────────

type svcMetrics struct {
	cpuPct      float64
	memoryBytes float64
	startTime   float64
	reqRate     float64 // handled RPCs/sec (all methods, all codes)
	errRate     float64 // non-OK RPCs/sec
	latencyP50  float64 // seconds
	latencyP95  float64 // seconds
	// Go runtime extras
	goroutines  float64
	heapBytes   float64
	openFDs     float64
	maxFDs      float64
	msgRecvRate float64 // gRPC msgs received/sec
	msgSentRate float64 // gRPC msgs sent/sec
}

func fetchPromMetrics(ctx context.Context, prom *promClient) (map[string]*svcMetrics, bool) {
	if !prom.reachable(ctx) {
		return nil, false
	}

	metrics := make(map[string]*svcMetrics)
	ensure := func(key string) *svcMetrics {
		if m, ok := metrics[key]; ok {
			return m
		}
		m := &svcMetrics{}
		metrics[key] = m
		return m
	}

	// CPU: rate(process_cpu_seconds_total[1m]) * 100
	if results, err := prom.query(ctx, "rate(process_cpu_seconds_total[1m])*100"); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			val := parsePromValue(r.Value[1])
			ensure(key).cpuPct = val
		}
	}

	// Memory: process_resident_memory_bytes
	if results, err := prom.query(ctx, "process_resident_memory_bytes"); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			ensure(key).memoryBytes = parsePromValue(r.Value[1])
		}
	}

	// Start time: process_start_time_seconds
	if results, err := prom.query(ctx, "process_start_time_seconds"); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			ensure(key).startTime = parsePromValue(r.Value[1])
		}
	}

	// Request rate: handled RPCs/sec
	if results, err := prom.query(ctx, `sum by (job)(rate(grpc_server_handled_total[1m]))`); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			ensure(key).reqRate = parsePromValue(r.Value[1])
		}
	}

	// Error rate: non-OK RPCs/sec
	if results, err := prom.query(ctx, `sum by (job)(rate(grpc_server_handled_total{grpc_code!="OK"}[1m]))`); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			ensure(key).errRate = parsePromValue(r.Value[1])
		}
	}

	// Latency p50
	if results, err := prom.query(ctx, `histogram_quantile(0.50, sum by (job, le)(rate(grpc_server_handling_seconds_bucket[5m])))`); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			v := parsePromValue(r.Value[1])
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				ensure(key).latencyP50 = v
			}
		}
	}

	// Latency p95
	if results, err := prom.query(ctx, `histogram_quantile(0.95, sum by (job, le)(rate(grpc_server_handling_seconds_bucket[5m])))`); err == nil {
		for _, r := range results {
			key := promJobKey(r.Metric)
			if key == "" {
				continue
			}
			v := parsePromValue(r.Value[1])
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				ensure(key).latencyP95 = v
			}
		}
	}

	// Goroutines
	if results, err := prom.query(ctx, "go_goroutines"); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).goroutines = parsePromValue(r.Value[1])
			}
		}
	}

	// Heap bytes
	if results, err := prom.query(ctx, "go_memstats_heap_alloc_bytes"); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).heapBytes = parsePromValue(r.Value[1])
			}
		}
	}

	// Open FDs
	if results, err := prom.query(ctx, "process_open_fds"); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).openFDs = parsePromValue(r.Value[1])
			}
		}
	}

	// Max FDs
	if results, err := prom.query(ctx, "process_max_fds"); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).maxFDs = parsePromValue(r.Value[1])
			}
		}
	}

	// gRPC msg recv rate
	if results, err := prom.query(ctx, `sum by (job)(rate(grpc_server_msg_received_total[1m]))`); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).msgRecvRate = parsePromValue(r.Value[1])
			}
		}
	}

	// gRPC msg sent rate
	if results, err := prom.query(ctx, `sum by (job)(rate(grpc_server_msg_sent_total[1m]))`); err == nil {
		for _, r := range results {
			if key := promJobKey(r.Metric); key != "" {
				ensure(key).msgSentRate = parsePromValue(r.Value[1])
			}
		}
	}

	return metrics, true
}

// fetchInfraMetrics queries Prometheus for infrastructure-specific metrics
// (etcd, envoy, node_exporter) and returns them keyed by service name.
func fetchInfraMetrics(ctx context.Context, prom *promClient) map[string]*InfraDetail {
	if !prom.reachable(ctx) {
		return nil
	}

	infra := make(map[string]*InfraDetail)
	ensure := func(key string) *InfraDetail {
		if d, ok := infra[key]; ok {
			return d
		}
		d := &InfraDetail{}
		infra[key] = d
		return d
	}

	// Helper: run a scalar query and apply the result via callback
	scalar := func(expr string, apply func(float64)) {
		results, err := prom.query(ctx, expr)
		if err != nil || len(results) == 0 {
			return
		}
		v := parsePromValue(results[0].Value[1])
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			apply(v)
		}
	}

	// etcd
	scalar("etcd_server_is_leader", func(v float64) {
		b := v == 1.0
		ensure("etcd").EtcdIsLeader = &b
	})
	scalar("etcd_debugging_mvcc_db_total_size_in_bytes", func(v float64) {
		ensure("etcd").EtcdDBSizeBytes = &v
	})
	scalar("etcd_debugging_mvcc_keys_total", func(v float64) {
		ensure("etcd").EtcdTotalKeys = &v
	})

	// envoy
	scalar("envoy_server_total_connections", func(v float64) {
		ensure("envoy").EnvoyActiveConns = &v
	})
	scalar(`sum(rate(envoy_http_downstream_rq_total[1m]))`, func(v float64) {
		ensure("envoy").EnvoyRPS = &v
	})
	scalar(`sum(rate(envoy_http_downstream_rq_xx{envoy_response_code_class="5"}[1m]))`, func(v float64) {
		ensure("envoy").EnvoyHTTP5xx = &v
	})

	// node
	scalar("node_load1", func(v float64) {
		ensure("node").NodeLoad1 = &v
	})
	scalar("node_load5", func(v float64) {
		ensure("node").NodeLoad5 = &v
	})
	scalar("node_memory_MemAvailable_bytes", func(v float64) {
		ensure("node").NodeMemAvailBytes = &v
	})
	scalar("node_memory_MemTotal_bytes", func(v float64) {
		ensure("node").NodeMemTotalBytes = &v
	})
	scalar(`sum(rate(node_network_receive_bytes_total{device!="lo"}[1m]))`, func(v float64) {
		ensure("node").NodeNetRxBytesRate = &v
	})
	scalar(`sum(rate(node_network_transmit_bytes_total{device!="lo"}[1m]))`, func(v float64) {
		ensure("node").NodeNetTxBytesRate = &v
	})

	return infra
}

// promJobKey extracts the service base name from metric labels.
// Tries job → service_name → instance, then normalizes.
func promJobKey(labels map[string]string) string {
	key := labels["job"]
	if key == "" {
		key = labels["service_name"]
	}
	if key == "" {
		key = labels["instance"]
	}
	if key == "" {
		return ""
	}
	// Normalize: strip ".service" suffix, take first dot-segment, lowercase
	key = strings.TrimSuffix(key, ".service")
	if idx := strings.IndexByte(key, '.'); idx > 0 {
		key = key[:idx]
	}
	// Strip port from instance-style values like "localhost:8080"
	if idx := strings.IndexByte(key, ':'); idx > 0 {
		key = key[:idx]
	}
	return strings.ToLower(key)
}

func parsePromValue(raw json.RawMessage) float64 {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// ── Health derivation ───────────────────────────────────────────────────────

const (
	cpuWarnPct = 85.0
	cpuCritPct = 95.0
	memWarnPct = 85.0 // percent of total system memory (unused in Phase 1; raw bytes threshold used)
	memCritPct = 95.0
)

func deriveServiceHealth(state string, pm *svcMetrics, promConnected bool) (string, []string) {
	state = strings.ToLower(state)
	var reasons []string

	// State-based rules
	switch state {
	case "failed", "error", "dead":
		reasons = append(reasons, "service state: "+state)
		return "critical", reasons
	case "stopped":
		return "critical", []string{"service state: stopped"}
	case "restarting", "degraded":
		reasons = append(reasons, "service state: "+state)
		// Don't return yet — Prometheus may escalate to critical
	case "running", "active":
		// Cross-validate with Prometheus: if Prometheus is connected but
		// has no process metrics for this service, the etcd state may be
		// stale (e.g. the process was killed without updating etcd).
		if promConnected && pm == nil {
			return "critical", []string{"service state: " + state + " (no process metrics — likely not running)"}
		}
	default:
		if state == "" || state == "unknown" {
			return "unknown", []string{"service state unknown"}
		}
		// Any other unrecognized state
		return "unknown", []string{"service state: " + state}
	}

	// Prometheus-based rules
	if promConnected && pm != nil {
		if pm.cpuPct > cpuCritPct {
			reasons = append(reasons, "CPU "+strconv.FormatFloat(pm.cpuPct, 'f', 1, 64)+"% > "+strconv.FormatFloat(cpuCritPct, 'f', 0, 64)+"%")
			return "critical", reasons
		}
		if pm.cpuPct > cpuWarnPct {
			reasons = append(reasons, "CPU "+strconv.FormatFloat(pm.cpuPct, 'f', 1, 64)+"% > "+strconv.FormatFloat(cpuWarnPct, 'f', 0, 64)+"%")
			if len(reasons) == 1 {
				// Only CPU warning, nothing worse from state
				return "degraded", reasons
			}
		}
	}

	// If we have state-based reasons from restarting/degraded
	if len(reasons) > 0 {
		return "degraded", reasons
	}

	return "healthy", nil
}

// ── Handler ─────────────────────────────────────────────────────────────────

// NewServicesHandler returns a GET-only handler for /admin/metrics/services.
func NewServicesHandler(provider AdminProvider) http.Handler {
	prom := newPromClient("http://localhost:9090", 8*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cfgs, err := provider.AllServiceConfigs()
		if err != nil {
			http.Error(w, "failed to list services: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch Prometheus metrics with timeout context
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		var promMetrics map[string]*svcMetrics
		var promConnected bool
		var infraDetail map[string]*InfraDetail

		var wg sync.WaitGroup
		wg.Add(2)
		go func() { defer wg.Done(); promMetrics, promConnected = fetchPromMetrics(ctx, prom) }()
		go func() { defer wg.Done(); infraDetail = fetchInfraMetrics(ctx, prom) }()
		wg.Wait()

		// Build service instances
		var all []ServiceInstance
		for _, cfg := range cfgs {
			name := mapStr(cfg, "Name")
			if name == "" {
				continue
			}
			base := strings.TrimSuffix(name, ".service")
			if idx := strings.IndexByte(base, '.'); idx > 0 {
				base = base[:idx]
			}
			base = strings.ToLower(base)

			inst := ServiceInstance{
				Name:        name,
				DisplayName: name,
				ID:          mapStr(cfg, "Id"),
				Version:     mapStr(cfg, "Version"),
				State:       mapStr(cfg, "State"),
				Port:        mapInt(cfg, "Port"),
				Category:    categorize(base),
				Node:        provider.Hostname(),
				GRPCHealth:  &GRPCHealth{Enabled: false, Status: "NOT_CHECKED"},
			}

			// Attach Prometheus runtime metrics
			var pm *svcMetrics
			if promMetrics != nil {
				pm = promMetrics[base]
			}
			if pm != nil {
				uptime := 0.0
				if pm.startTime > 0 {
					uptime = float64(time.Now().Unix()) - pm.startTime
					if uptime < 0 {
						uptime = 0
					}
				}
				inst.Runtime = &SvcRuntime{
					CPUPct:       math.Round(pm.cpuPct*10) / 10,
					MemoryBytes:  pm.memoryBytes,
					UptimeSec:    uptime,
					ReqRate:      math.Round(pm.reqRate*100) / 100,
					ErrRate:      math.Round(pm.errRate*100) / 100,
					LatencyP50Ms: math.Round(pm.latencyP50*1000*100) / 100,
					LatencyP95Ms: math.Round(pm.latencyP95*1000*100) / 100,
					Goroutines:   pm.goroutines,
					HeapBytes:    pm.heapBytes,
					OpenFDs:      pm.openFDs,
					MaxFDs:       pm.maxFDs,
					MsgRecvRate:  math.Round(pm.msgRecvRate*100) / 100,
					MsgSentRate:  math.Round(pm.msgSentRate*100) / 100,
				}
			}

			// Derive health status
			status, reasons := deriveServiceHealth(inst.State, pm, promConnected)
			inst.DerivedStatus = status
			inst.Reasons = reasons

			all = append(all, inst)
		}

		// Group by category
		grouped := groupServices(all)

		// Build summary
		summary := ServicesSummary{Total: len(all)}
		for _, s := range all {
			switch s.DerivedStatus {
			case "healthy":
				summary.Healthy++
			case "degraded":
				summary.Degraded++
			case "critical":
				summary.Critical++
			default:
				summary.Unknown++
			}
		}

		resp := ServicesResponse{
			NowUnixMs: time.Now().UnixMilli(),
			Range:     "instant",
			Prometheus: PromStatus{
				Connected: promConnected,
				Addr:      "http://localhost:9090",
			},
			Thresholds: SvcThresholds{
				CPUWarn: cpuWarnPct,
				CPUCrit: cpuCritPct,
				MemWarn: memWarnPct,
				MemCrit: memCritPct,
			},
			Groups:  grouped,
			Summary: summary,
			Infra:   infraDetail,
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}

func groupServices(all []ServiceInstance) []ServiceGroup {
	byCategory := make(map[string][]ServiceInstance)
	for _, s := range all {
		byCategory[s.Category] = append(byCategory[s.Category], s)
	}
	order := []string{"Core", "Infrastructure", "Media", "Other"}
	var groups []ServiceGroup
	for _, cat := range order {
		if svcs, ok := byCategory[cat]; ok && len(svcs) > 0 {
			groups = append(groups, ServiceGroup{Category: cat, Services: svcs})
		}
	}
	return groups
}

// ── map helpers ─────────────────────────────────────────────────────────────

func mapStr(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return ""
}

func mapInt(m map[string]any, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	}
	return 0
}
