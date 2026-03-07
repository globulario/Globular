package admin

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"
)

// ── JSON response types ─────────────────────────────────────────────────────

// EnvoyResponse is the top-level response for GET /admin/metrics/envoy.
type EnvoyResponse struct {
	NowUnixMs  int64           `json:"now_unix_ms"`
	Healthy    bool            `json:"healthy"`
	Prometheus PromStatus      `json:"prometheus"`
	Server     EnvoyServer     `json:"server"`
	Downstream EnvoyDownstream `json:"downstream"`
	Clusters   []EnvoyCluster  `json:"clusters"`
	Listeners  []EnvoyListener `json:"listeners"`
	XDS        EnvoyXDS        `json:"xds"`
}

type EnvoyServer struct {
	State           string  `json:"state"`
	UptimeSec       float64 `json:"uptime_sec"`
	Connections     float64 `json:"connections"`
	MemAllocated    float64 `json:"mem_allocated_bytes"`
	TotalConnsLife  float64 `json:"total_connections_lifetime"`
	HotRestartEpoch float64 `json:"hot_restart_epoch"`
	Version         string  `json:"version"`
}

type EnvoyDownstream struct {
	ActiveConns   float64 `json:"active_conns"`
	RPS           float64 `json:"rps"`
	HTTP2xxRate   float64 `json:"http_2xx_rate"`
	HTTP4xxRate   float64 `json:"http_4xx_rate"`
	HTTP5xxRate   float64 `json:"http_5xx_rate"`
	RxBytesRate   float64 `json:"rx_bytes_rate"`
	TxBytesRate   float64 `json:"tx_bytes_rate"`
	SSLConns      float64 `json:"ssl_conns"`
	SSLHandshakes float64 `json:"ssl_handshake_rate"`
	P50Ms         float64 `json:"p50_ms"`
	P95Ms         float64 `json:"p95_ms"`
	P99Ms         float64 `json:"p99_ms"`
	SSLErrorRate  float64 `json:"ssl_error_rate"`
	DaysUntilCert float64 `json:"days_until_cert_expiry"`
}

type EnvoyCluster struct {
	Name        string  `json:"name"`
	Healthy     int     `json:"healthy"`
	Degraded    int     `json:"degraded"`
	Unhealthy   int     `json:"unhealthy"`
	RPS         float64 `json:"rps"`
	ErrRate     float64 `json:"err_rate"`
	P50Ms       float64 `json:"p50_ms"`
	P99Ms       float64 `json:"p99_ms"`
	ActiveConns float64 `json:"active_conns"`
	RxBytesRate float64 `json:"rx_bytes_rate"`
	TxBytesRate float64 `json:"tx_bytes_rate"`
	RetryRate   float64 `json:"retry_rate"`
	TimeoutRate float64 `json:"timeout_rate"`
	RxResetRate float64 `json:"rx_reset_rate"`
	CBOpen      float64 `json:"circuit_breaker_open"`
}

type EnvoyListener struct {
	Address       string  `json:"address"`
	ActiveConns   float64 `json:"active_conns"`
	RPS           float64 `json:"rps"`
	HTTP4xxRate   float64 `json:"http_4xx_rate"`
	HTTP5xxRate   float64 `json:"http_5xx_rate"`
	SSLHandshakes float64 `json:"ssl_handshake_rate"`
	SSLErrors     float64 `json:"ssl_error_rate"`
}

type EnvoyXDS struct {
	ActiveClusters  float64    `json:"active_clusters"`
	ActiveListeners float64    `json:"active_listeners"`
	CDSSuccess      float64    `json:"cds_update_success"`
	CDSFailure      float64    `json:"cds_update_failure"`
	LDSSuccess      float64    `json:"lds_update_success"`
	LDSFailure      float64    `json:"lds_update_failure"`
	Routes          []RDSRoute `json:"routes"`
}

type RDSRoute struct {
	Name          string  `json:"name"`
	Connected     float64 `json:"connected"`
	UpdateSuccess float64 `json:"update_success"`
	UpdateFailure float64 `json:"update_failure"`
}

// envoy server state enum → string
var envoyStateNames = map[int]string{
	0: "LIVE",
	1: "DRAINING",
	2: "PRE_INIT",
	3: "INIT",
}

// ── Handler ─────────────────────────────────────────────────────────────────

// NewEnvoyHandler returns a GET-only handler for /admin/metrics/envoy.
func NewEnvoyHandler() http.Handler {
	prom := newPromClient("http://localhost:9090", 8*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()

		resp := EnvoyResponse{NowUnixMs: time.Now().UnixMilli()}
		resp.Prometheus = PromStatus{Addr: prom.addr}

		if !prom.reachable(ctx) {
			resp.Prometheus.Connected = false
			w.Header().Set("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			_ = enc.Encode(resp)
			return
		}
		resp.Prometheus.Connected = true

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

		// vectorByCluster runs a per-cluster query and calls apply for each result.
		vectorByCluster := func(expr string, apply func(name string, v float64)) {
			results, err := prom.query(ctx, expr)
			if err != nil {
				return
			}
			for _, r := range results {
				name := r.Metric["envoy_cluster_name"]
				if name == "" {
					continue
				}
				v := parsePromValue(r.Value[1])
				if !math.IsNaN(v) && !math.IsInf(v, 0) {
					apply(name, v)
				}
			}
		}

		// vectorByLabel runs a per-label query and calls apply for each result.
		vectorByLabel := func(expr, label string, apply func(name string, v float64)) {
			results, err := prom.query(ctx, expr)
			if err != nil {
				return
			}
			for _, r := range results {
				name := r.Metric[label]
				if name == "" {
					continue
				}
				v := parsePromValue(r.Value[1])
				if !math.IsNaN(v) && !math.IsInf(v, 0) {
					apply(name, v)
				}
			}
		}

		var wg sync.WaitGroup

		// ── Server metrics ──────────────────────────────────────────────
		wg.Add(1)
		go func() {
			defer wg.Done()
			scalar("envoy_server_state", func(v float64) {
				if name, ok := envoyStateNames[int(v)]; ok {
					resp.Server.State = name
				}
			})
			scalar("envoy_server_uptime", func(v float64) {
				resp.Server.UptimeSec = v
			})
			scalar("envoy_server_total_connections", func(v float64) {
				resp.Server.Connections = v
			})
			scalar("envoy_server_memory_allocated", func(v float64) {
				resp.Server.MemAllocated = v
			})
			scalar("envoy_server_hot_restart_epoch", func(v float64) {
				resp.Server.HotRestartEpoch = v
			})
			// Extract version from any envoy metric's label
			results, err := prom.query(ctx, "envoy_server_uptime")
			if err == nil && len(results) > 0 {
				if ver, ok := results[0].Metric["envoy_version"]; ok {
					resp.Server.Version = ver
				}
			}
		}()

		// ── Downstream metrics ──────────────────────────────────────────
		wg.Add(1)
		go func() {
			defer wg.Done()
			scalar("sum(envoy_http_downstream_cx_active)", func(v float64) {
				resp.Downstream.ActiveConns = v
			})
			scalar(`sum(rate(envoy_http_downstream_rq_total[1m]))`, func(v float64) {
				resp.Downstream.RPS = v
			})
			scalar(`sum(rate(envoy_http_downstream_rq_xx{envoy_response_code_class="2"}[1m]))`, func(v float64) {
				resp.Downstream.HTTP2xxRate = v
			})
			scalar(`sum(rate(envoy_http_downstream_rq_xx{envoy_response_code_class="4"}[1m]))`, func(v float64) {
				resp.Downstream.HTTP4xxRate = v
			})
			scalar(`sum(rate(envoy_http_downstream_rq_xx{envoy_response_code_class="5"}[1m]))`, func(v float64) {
				resp.Downstream.HTTP5xxRate = v
			})
			scalar(`sum(rate(envoy_http_downstream_cx_rx_bytes_total[1m]))`, func(v float64) {
				resp.Downstream.RxBytesRate = v
			})
			scalar(`sum(rate(envoy_http_downstream_cx_tx_bytes_total[1m]))`, func(v float64) {
				resp.Downstream.TxBytesRate = v
			})
		}()

		// ── SSL metrics (fixed: use correct metric names) ───────────────
		wg.Add(1)
		go func() {
			defer wg.Done()
			scalar(`sum(envoy_http_downstream_cx_ssl_active{envoy_http_conn_manager_prefix="http"})`, func(v float64) {
				resp.Downstream.SSLConns = v
			})
			scalar(`sum(rate(envoy_listener_ssl_handshake[1m]))`, func(v float64) {
				resp.Downstream.SSLHandshakes = v
			})
			scalar(`sum(rate(envoy_listener_ssl_connection_error[1m]))`, func(v float64) {
				resp.Downstream.SSLErrorRate = v
			})
			scalar("envoy_server_days_until_first_cert_expiring", func(v float64) {
				resp.Downstream.DaysUntilCert = v
			})
		}()

		// ── Downstream latency histograms ───────────────────────────────
		wg.Add(1)
		go func() {
			defer wg.Done()
			scalar(`histogram_quantile(0.5, sum by (le)(rate(envoy_http_downstream_rq_time_bucket{envoy_http_conn_manager_prefix="http"}[5m])))`, func(v float64) {
				resp.Downstream.P50Ms = v
			})
			scalar(`histogram_quantile(0.95, sum by (le)(rate(envoy_http_downstream_rq_time_bucket{envoy_http_conn_manager_prefix="http"}[5m])))`, func(v float64) {
				resp.Downstream.P95Ms = v
			})
			scalar(`histogram_quantile(0.99, sum by (le)(rate(envoy_http_downstream_rq_time_bucket{envoy_http_conn_manager_prefix="http"}[5m])))`, func(v float64) {
				resp.Downstream.P99Ms = v
			})
		}()

		// ── Per-listener metrics ────────────────────────────────────────
		type listenerData struct {
			mu sync.Mutex
			m  map[string]*EnvoyListener
		}
		ld := &listenerData{m: make(map[string]*EnvoyListener)}
		ensureListener := func(addr string) *EnvoyListener {
			ld.mu.Lock()
			defer ld.mu.Unlock()
			if l, ok := ld.m[addr]; ok {
				return l
			}
			l := &EnvoyListener{Address: addr}
			ld.m[addr] = l
			return l
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			vectorByLabel("envoy_listener_downstream_cx_active", "envoy_listener_address", func(addr string, v float64) {
				ensureListener(addr).ActiveConns = v
			})
			vectorByLabel(`sum by (envoy_listener_address)(rate(envoy_listener_downstream_cx_total[1m]))`, "envoy_listener_address", func(addr string, v float64) {
				ensureListener(addr).RPS = v
			})
			vectorByLabel(`sum by (envoy_listener_address)(rate(envoy_listener_ssl_handshake[1m]))`, "envoy_listener_address", func(addr string, v float64) {
				ensureListener(addr).SSLHandshakes = v
			})
			vectorByLabel(`sum by (envoy_listener_address)(rate(envoy_listener_ssl_connection_error[1m]))`, "envoy_listener_address", func(addr string, v float64) {
				ensureListener(addr).SSLErrors = v
			})
		}()

		// ── xDS control plane metrics ───────────────────────────────────
		wg.Add(1)
		go func() {
			defer wg.Done()
			scalar("envoy_cluster_manager_active_clusters", func(v float64) {
				resp.XDS.ActiveClusters = v
			})
			scalar("envoy_listener_manager_total_listeners_active", func(v float64) {
				resp.XDS.ActiveListeners = v
			})
			scalar("envoy_cluster_manager_cds_update_success", func(v float64) {
				resp.XDS.CDSSuccess = v
			})
			scalar("envoy_cluster_manager_cds_update_failure", func(v float64) {
				resp.XDS.CDSFailure = v
			})
			scalar("envoy_listener_manager_lds_update_success", func(v float64) {
				resp.XDS.LDSSuccess = v
			})
			scalar("envoy_listener_manager_lds_update_failure", func(v float64) {
				resp.XDS.LDSFailure = v
			})

			// RDS routes
			type rdsData struct {
				mu sync.Mutex
				m  map[string]*RDSRoute
			}
			rd := &rdsData{m: make(map[string]*RDSRoute)}
			ensureRoute := func(name string) *RDSRoute {
				rd.mu.Lock()
				defer rd.mu.Unlock()
				if r, ok := rd.m[name]; ok {
					return r
				}
				r := &RDSRoute{Name: name}
				rd.m[name] = r
				return r
			}
			vectorByLabel("envoy_http_rds_connected_state", "envoy_rds_route_config", func(name string, v float64) {
				ensureRoute(name).Connected = v
			})
			vectorByLabel("envoy_http_rds_update_success", "envoy_rds_route_config", func(name string, v float64) {
				ensureRoute(name).UpdateSuccess = v
			})
			vectorByLabel("envoy_http_rds_update_failure", "envoy_rds_route_config", func(name string, v float64) {
				ensureRoute(name).UpdateFailure = v
			})

			routes := make([]RDSRoute, 0, len(rd.m))
			for _, r := range rd.m {
				routes = append(routes, *r)
			}
			sort.Slice(routes, func(i, j int) bool {
				return routes[i].Name < routes[j].Name
			})
			resp.XDS.Routes = routes
		}()

		// ── Upstream cluster metrics ────────────────────────────────────
		type clusterData struct {
			mu sync.Mutex
			m  map[string]*EnvoyCluster
		}
		cd := &clusterData{m: make(map[string]*EnvoyCluster)}
		ensure := func(name string) *EnvoyCluster {
			cd.mu.Lock()
			defer cd.mu.Unlock()
			if c, ok := cd.m[name]; ok {
				return c
			}
			c := &EnvoyCluster{Name: name}
			cd.m[name] = c
			return c
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			vectorByCluster("envoy_cluster_membership_healthy", func(name string, v float64) {
				ensure(name).Healthy = int(v)
			})
			vectorByCluster("envoy_cluster_membership_degraded", func(name string, v float64) {
				ensure(name).Degraded = int(v)
			})
			vectorByCluster("envoy_cluster_membership_unhealthy", func(name string, v float64) {
				ensure(name).Unhealthy = int(v)
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_rq_total[1m]))`, func(name string, v float64) {
				ensure(name).RPS = v
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_rq_xx{envoy_response_code_class="5"}[1m]))`, func(name string, v float64) {
				ensure(name).ErrRate = v
			})
			vectorByCluster("envoy_cluster_upstream_cx_active", func(name string, v float64) {
				ensure(name).ActiveConns = v
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_cx_rx_bytes_total[1m]))`, func(name string, v float64) {
				ensure(name).RxBytesRate = v
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_cx_tx_bytes_total[1m]))`, func(name string, v float64) {
				ensure(name).TxBytesRate = v
			})
			// Per-cluster extras
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_rq_retry[1m]))`, func(name string, v float64) {
				ensure(name).RetryRate = v
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_rq_timeout[1m]))`, func(name string, v float64) {
				ensure(name).TimeoutRate = v
			})
			vectorByCluster(`sum by (envoy_cluster_name)(rate(envoy_cluster_upstream_rq_rx_reset[1m]))`, func(name string, v float64) {
				ensure(name).RxResetRate = v
			})
			vectorByCluster("envoy_cluster_circuit_breakers_default_rq_open", func(name string, v float64) {
				ensure(name).CBOpen = v
			})
		}()

		// ── Latency histograms (separate goroutine — heavier queries) ───
		wg.Add(1)
		go func() {
			defer wg.Done()
			vectorByCluster(`histogram_quantile(0.5, sum by (envoy_cluster_name, le)(rate(envoy_cluster_upstream_rq_time_bucket[5m])))`, func(name string, v float64) {
				ensure(name).P50Ms = v
			})
			vectorByCluster(`histogram_quantile(0.99, sum by (envoy_cluster_name, le)(rate(envoy_cluster_upstream_rq_time_bucket[5m])))`, func(name string, v float64) {
				ensure(name).P99Ms = v
			})
		}()

		wg.Wait()

		// Build sorted clusters slice
		clusters := make([]EnvoyCluster, 0, len(cd.m))
		for _, c := range cd.m {
			clusters = append(clusters, *c)
		}
		sort.Slice(clusters, func(i, j int) bool {
			return clusters[i].Name < clusters[j].Name
		})
		resp.Clusters = clusters

		// Build sorted listeners slice
		listeners := make([]EnvoyListener, 0, len(ld.m))
		for _, l := range ld.m {
			listeners = append(listeners, *l)
		}
		sort.Slice(listeners, func(i, j int) bool {
			return listeners[i].Address < listeners[j].Address
		})
		resp.Listeners = listeners

		// Determine overall health
		resp.Healthy = resp.Server.State == "LIVE" || resp.Server.State == ""
		if resp.Healthy {
			for _, c := range clusters {
				if c.Unhealthy > 0 {
					resp.Healthy = false
					break
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		_ = enc.Encode(resp)
	})
}
