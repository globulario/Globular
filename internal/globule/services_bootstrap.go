package globule

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/globulario/Globular/internal/logsink"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// at top of file
type restartGuard struct {
	mu   sync.Mutex
	busy map[string]struct{}
}

func (g *Globule) startKeepAliveSupervisor(ctx context.Context) {
	guard := &restartGuard{busy: make(map[string]struct{})}

	// 1) Initial reconcile — use RUNTIME truth, not config snapshot
	if svcs, err := config.GetServicesConfigurations(); err == nil {
		for _, s := range svcs {
			id := Utility.ToString(s["Id"])
			if id == "" || !Utility.ToBool(s["KeepAlive"]) {
				continue
			}
			name := Utility.ToString(s["Name"])
			port := Utility.ToInt(s["Port"])
			rt, err := config.GetRuntime(id)
			if err != nil {
				g.log.Warn("keepalive: runtime not found; abort", "id", id, "err", err)
				continue
			}
			state := strings.ToLower(Utility.ToString(rt["State"]))
			updated := Utility.ToInt(rt["UpdatedAt"])

			// If we recently moved to "starting", give it time.
			if state == "starting" && g.recentlyStarting(int64(updated), 15*time.Second) {
				continue
			}

			// If it’s actually alive on the network, correct etcd and skip.
			if g.isActuallyRunning(name, port) {
				_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
				continue
			}

			// Otherwise, try to (re)start.
			go func(id string, delay time.Duration) {
				time.Sleep(delay) // 0–300ms jitter
				g.tryRestartWithBackoff(ctx, guard, id)
			}(id, time.Duration(50+rand.Intn(250))*time.Millisecond)
		}
	}

	// 2) Live watch — act only on *real* down states, with liveness re-check
	go func() {
		_ = config.WatchRuntimes(ctx, func(ev config.RuntimeEvent) {
			cfg, err := config.GetServiceConfigurationById(ev.ID)
			if err != nil || cfg == nil || !Utility.ToBool(cfg["KeepAlive"]) {
				return
			}

			state := strings.ToLower(Utility.ToString(ev.Runtime["State"]))
			proc := Utility.ToInt(ev.Runtime["Process"])
			name := Utility.ToString(cfg["Name"])
			port := Utility.ToInt(cfg["Port"])
			updated := Utility.ToInt(ev.Runtime["UpdatedAt"])

			// Ignore short "starting" to "running" churn
			if state == "starting" && g.recentlyStarting(int64(updated), 15*time.Second) {
				return
			}

			// If the service is actually responding, normalize runtime and skip.
			if g.isActuallyRunning(name, port) {
				_ = config.PutRuntime(ev.ID, map[string]any{"State": "running", "LastError": ""})
				return
			}

			// Only restart on true-down
			if proc == -1 || state == "failed" || state == "stopped" {
				go g.tryRestartWithBackoff(ctx, guard, ev.ID)
			}
		})
	}()
}

func (g *Globule) tryRestartWithBackoff(ctx context.Context, guard *restartGuard, id string) {
	// per-id dedupe
	guard.mu.Lock()
	if _, ok := guard.busy[id]; ok {
		guard.mu.Unlock()
		return
	}
	guard.busy[id] = struct{}{}
	guard.mu.Unlock()
	defer func() { guard.mu.Lock(); delete(guard.busy, id); guard.mu.Unlock() }()

	backoff := time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		cfg, err := config.GetServiceConfigurationById(id)
		if err != nil || cfg == nil {
			g.log.Warn("keepalive: config not found; abort", "id", id, "err", err)
			return
		}
		if !Utility.ToBool(cfg["KeepAlive"]) {
			g.log.Info("keepalive: disabled; abort", "id", id)
			return
		}

		name := Utility.ToString(cfg["Name"])
		port := Utility.ToInt(cfg["Port"])
		address, _ := config.GetAddress()
		address = strings.Split(address, ":")[0] // host only
		addr := net.JoinHostPort(address, fmt.Sprint(port))

		// FINAL sanity: if already up, just mark running and stop.
		if g.isActuallyRunning(name, port) {
			_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
			return
		}

		// Start service (log sinks)

		outW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stdout", logpb.LogLevel_INFO_MESSAGE, os.Stdout)
		errW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stderr", logpb.LogLevel_ERROR_MESSAGE, os.Stderr)

		if port <= 0 {
			if p, perr := config.AllocatePortForService(id); perr == nil {
				port = p
				cfg["Port"] = port
				cfg["Proxy"] = port + 1
				_ = config.SaveServiceConfiguration(cfg)
				addr = net.JoinHostPort(address, fmt.Sprint(port))
			}
		}

		g.log.Info("keepalive: restarting service", "name", name, "id", id, "port", port)
		pid, startErr := process.StartServiceProcessWithWriters(cfg, port, outW, errW)
		if startErr == nil {
			_ = config.PutRuntime(id, map[string]any{"Process": pid, "State": "starting", "LastError": ""})
			cfg["Process"] = pid

			if _, perr := process.StartServiceProxyProcess(cfg, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); perr != nil {
				g.log.Warn("keepalive: proxy start failed", "name", name, "err", perr)
			}

			// short readiness wait
			if g.waitServiceReady(name, addr, 8*time.Second) {
				_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
			} else {
				_ = config.PutRuntime(id, map[string]any{"State": "failed", "LastError": "keepalive startup timeout"})
			}
			return
		}

		g.log.Warn("keepalive: restart failed; retrying", "id", id, "err", startErr, "backoff", backoff)
		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return
		}
		if backoff < 30*time.Second {
			backoff *= 2
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}
		}
	}
}

// ---------- helpers ----------
func (g *Globule) isActuallyRunning(name string, port int) bool {
	if port <= 0 {
		return false
	}

	addr, _ := config.GetAddress()
	addr = strings.Split(addr, ":")[0] // host only
	addr = net.JoinHostPort(addr, fmt.Sprint(port))

	// TCP connect
	if c, err := net.DialTimeout("tcp", addr, 600*time.Millisecond); err == nil {
		_ = c.Close()
		// gRPC health
		if g.grpcHealthOK(addr) {
			return true
		}
		// fallback to --health
		if g.binHealthOK(name) {
			return true
		}
	}

	return false
}

func (g *Globule) recentlyStarting(updatedAt int64, grace time.Duration) bool {
	if updatedAt <= 0 {
		return false
	}
	// updatedAt is ms since epoch in your etcd docs; tolerate both seconds/ms
	nowMs := time.Now().UnixMilli()
	tsMs := updatedAt
	if updatedAt < 2_000_000_000 { // seconds?
		tsMs = updatedAt * 1000
	}
	return nowMs-tsMs < grace.Milliseconds()
}

// Use the strongly-typed descriptor from config.
type serviceDesc = config.ServiceDesc

// ------------------------------
// Public entry points
// ------------------------------

func shortName(full string) string {
	if full == "" {
		return ""
	}
	return strings.ToLower(strings.Split(full, ".")[0])
}

// helper (put near Globule methods)
func isListening(addr string, timeout time.Duration) bool {
	d := net.Dialer{Timeout: timeout}
	c, err := d.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

// startServicesEtcd discovers binaries, authors desired in etcd, orders by dependencies,
// spawns, then waits for health readiness with bounded timeout.
func (g *Globule) startServicesEtcd(ctx context.Context) error {
	// 1) Peer keys & local token
	if err := security.GeneratePeerKeys(g.Mac); err != nil {
		return err
	}
	if err := security.SetLocalToken(g.Mac, g.Domain, "sa", "sa", g.AdminEmail, g.SessionTimeout); err != nil {
		return err
	}

	// 2) Discover executables (shared in config)
	bins, err := config.DiscoverExecutables(config.GetServicesRoot())
	if err != nil {
		return err
	}
	if len(bins) == 0 {
		return errors.New("no service executables found; check GetServicesRoot()")
	}

	// 3) Allocate ports with conflict checks against existing desired
	alloc, err := config.NewPortAllocator(g.PortsRange)
	if err != nil {
		return err
	}
	if existing, _ := config.GetServicesConfigurations(); len(existing) > 0 {
		alloc.ReserveExisting(existing)
	}

	// 4) Describe -> normalize -> merge/author desired
	address, _ := config.GetAddress()
	address = strings.Split(address, ":")[0] // host only

	desiredByID := map[string]map[string]interface{}{}
	for _, bin := range bins {
		desc, err := config.RunDescribe(bin, 2500*time.Millisecond, map[string]string{
			"GLOBULAR_DOMAIN":  g.Domain,
			"GLOBULAR_ADDRESS": g.localIPAddress,
		})
		if err != nil {
			g.log.Warn("describe failed", "bin", bin, "err", err)
			continue
		}
		g.normalizeDescriptor(&desc)

		m, err := g.mergeOrCreateDesired(desc, alloc)
		if err != nil {
			g.log.Error("merge/create desired failed", "name", desc.Name, "err", err)
			continue
		}
		if err := config.SaveServiceConfiguration(m); err != nil {
			g.log.Error("save desired failed", "service", desc.Name, "id", m["Id"], "err", err)
			continue
		}
		desiredByID[Utility.ToString(m["Id"])] = m
	}

	// 5) Order by dependencies
	ordered, err := g.topoOrder(desiredByID)
	if err != nil {
		return err
	}
	if len(ordered) == 0 {
		return errors.New("no services to start after desired authoring")
	}

	// 6) Spawn with log sinks + proxies
	for _, id := range ordered {
		m := desiredByID[id]
		name := Utility.ToString(m["Name"])
		port := Utility.ToInt(m["Port"])
		proxy := Utility.ToInt(m["Proxy"])
		bin := Utility.ToString(m["Path"])
		if bin == "" {
			g.log.Error("missing Path for service", "name", name, "id", id)
			continue
		}

		// Skip if already running (per runtime state or port check)
		addr := address + ":" + Utility.ToString(m["Port"])
		if isListening(addr, 300*time.Millisecond) {
			g.log.Info("service already running; skipping start", "name", name, "addr", addr)
			_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
			continue
		}

		outW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stdout", logpb.LogLevel_INFO_MESSAGE, os.Stdout)
		errW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stderr", logpb.LogLevel_ERROR_MESSAGE, os.Stderr)

		g.log.Info("starting service", "name", name, "id", id, "port", port, "proxy", proxy)
		pid, err := process.StartServiceProcessWithWriters(m, port, outW, errW)
		if err != nil {
			g.log.Warn("service start failed", "name", name, "err", err)
			_ = config.PutRuntime(id, map[string]any{"Process": -1, "State": "failed", "LastError": err.Error()})
			continue
		}

		// IMPORTANT: update in-memory map so proxy code sees the pid
		_ = config.PutRuntime(id, map[string]any{"Process": pid, "State": "starting", "LastError": ""})
		m["Process"] = pid

		if _, err := process.StartServiceProxyProcess(m, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); err != nil {
			g.log.Warn("proxy start failed", "name", name, "err", err)
		}
	}

	// 7) Readiness wait (bounded)
	deadline := time.Now().Add(20 * time.Second)
	for _, id := range ordered {
		m := desiredByID[id]
		name := Utility.ToString(m["Name"])
		addr := address + ":" + Utility.ToString(m["Port"])

		ok := g.waitServiceReady(name, addr, deadline.Sub(time.Now()))
		if !ok {
			_ = config.PutRuntime(id, map[string]any{"State": "failed", "LastError": "startup timeout"})
			g.log.Warn("service failed to become ready", "name", name)
			continue
		}
		_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
		g.log.Info("service ready", "name", name, "addr", addr)
	}

	go refreshTokenPeriodically(ctx, g)

	go g.startKeepAliveSupervisor(ctx)

	return nil
}

// stopServicesEtcd stops proxies first, then services (same behavior as before).
func (g *Globule) stopServicesEtcd() error {

	svcs, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	// 1) Stop proxies first (unchanged)
	for _, s := range svcs {
		_ = process.KillServiceProxyProcess(s)
	}

	// 2) Ask each service to close cooperatively by setting desired State="closing"
	for _, s := range svcs {
		id := Utility.ToString(s["Id"])
		// reload desired to avoid clobbering
		if cur, err := config.GetServiceConfigurationById(id); err == nil && cur != nil {
			cur["State"] = "closing"
			g.log.Info("setting desired State=closing for service %s (%s)", Utility.ToString(s["Name"]), id)
			_ = config.SaveServiceConfiguration(cur) // triggers watchDesiredConfig in the service
		} else {
			g.log.Warn("cannot set desired State=closing; config not found", "id", id)
		}
	}

	// 3) Wait for runtime State="closed" with a timeout; fallback to hard kill if needed
	deadline := time.Now().Add(30 * time.Second)
	closed := map[string]bool{}
	for time.Now().Before(deadline) {
		allClosed := true
		for _, s := range svcs {
			id := Utility.ToString(s["Id"])
			if closed[id] {
				continue
			}
			// Fetch runtime (process/state) — if your config pkg has a getter use it.
			rt, _ := config.GetRuntime(id) // implement if not present; or keep runtime under /globular/services/<id>/runtime
			if Utility.ToString(rt["State"]) == "closed" {
				closed[id] = true
				g.log.Info("service closed", "id", id, "name", Utility.ToString(s["Name"]))
				continue
			}
			allClosed = false
		}
		if allClosed {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// 4) Hard-kill any stragglers and mark them failed, so issues are visible
	for _, s := range svcs {
		id := Utility.ToString(s["Id"])
		if closed[id] {
			continue
		}
		g.log.Warn("service did not close in time; forcing kill", "id", id, "name", Utility.ToString(s["Name"]))
		_ = process.KillServiceProcess(s)
		_ = config.PutRuntime(id, map[string]any{"Process": -1, "State": "failed", "LastError": "forced kill on shutdown"})
	}

	return nil
}

// ------------------------------
// Helpers
// ------------------------------
func (g *Globule) normalizeDescriptor(d *serviceDesc) {
	// Domain/address preference: trust process env or fall back to globule’s values
	if strings.TrimSpace(d.Domain) == "" {
		d.Domain = strings.ToLower(g.Domain)
	}
	if strings.TrimSpace(d.Address) == "" {
		host := d.Domain
		if host == "" {
			host = "localhost"
		}
		d.Address = host
	} else {
		d.Address = config.HostOnly(d.Address)
	}

	// Ensure Path is absolute and normalized (last resort; prefer Path from describe)
	if d.Path == "" {
		d.Path = strings.ReplaceAll(binDirOf(d.Name), "\\", "/")
	}

	// make sure core deps are present unless the service *is* the core
	if !strings.EqualFold(d.Name, "log.LogService") && !strings.EqualFold(d.Name, "rbac.RbacService") {
		need := map[string]bool{"log.LogService": true, "rbac.RbacService": true}
		for _, x := range d.Dependencies {
			delete(need, x)
		}
		for x := range need {
			d.Dependencies = append(d.Dependencies, x)
		}
	}
}

func binDirOf(name string) string {
	// last resort; prefer Path from describe
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return p
}

// mergeOrCreateDesired returns the authoritative desired document for this service id.
func (g *Globule) mergeOrCreateDesired(desc serviceDesc, alloc *config.PortAllocator) (map[string]interface{}, error) {
	var desired map[string]interface{}

	// Try to load existing desired by Id first.
	if cfg, err := config.GetServiceConfigurationById(desc.Id); err == nil && cfg != nil {
		desired = cfg
		if Utility.ToString(desired["Path"]) == "" && desc.Path != "" {
			desired["Path"] = desc.Path
		}
		if Utility.ToString(desired["Checksum"]) == "" && desc.Checksum != "" {
			desired["Checksum"] = desc.Checksum
		}
	} else {
		// New desired
		if desc.Id == "" {
			if desc.Name == "" || desc.Address == "" {
				return nil, errors.New("missing Name or Address in service descriptor; cannot generate stable id")
			}
			desc.Id = Utility.GenerateUUID(desc.Name + ":" + desc.Address)
		}

		// Decide port/proxy:
		var (
			port  = desc.Port
			proxy = desc.Proxy
			err   error
		)
		if port == 0 || proxy == 0 {
			port, proxy, err = alloc.NextPair(desc.Id)
			if err != nil {
				return nil, err
			}
		} else {
			if err = alloc.ClaimPair(desc.Id, port, proxy); err != nil {
				port, proxy, err = alloc.NextPair(desc.Id)
				if err != nil {
					return nil, err
				}
			}
		}

		desired = map[string]any{
			"Id":                 desc.Id,
			"Name":               desc.Name,
			"Description":        desc.Description,
			"PublisherID":        desc.PublisherID,
			"Version":            desc.Version,
			"Proto":              desc.Proto,
			"Path":               desc.Path,
			"Checksum":           desc.Checksum,
			"Keywords":           toAnySlice(desc.Keywords),
			"Repositories":       toAnySlice(desc.Repositories),
			"Discoveries":        toAnySlice(desc.Discoveries),
			"Permissions":        toAnyAnySlice(desc.Permissions),
			"Dependencies":       toAnySlice(desc.Dependencies),
			"Domain":             strings.ToLower(desc.Domain),
			"Address":            strings.ToLower(desc.Address),
			"Protocol":           coalesce(desc.Protocol, "grpc"),
			"Port":               port,
			"Proxy":              proxy,
			"TLS":                desc.TLS,
			"CertAuthorityTrust": desc.CertAuthorityTrust,
			"CertFile":           desc.CertFile,
			"KeyFile":            desc.KeyFile,
			"AllowAllOrigins":    desc.AllowAllOrigins,
			"AllowedOrigins":     desc.AllowedOrigins,
			"KeepAlive":          true,
			"KeepUpToDate":       true,
		}
	}

	// Enforce domain/address from Globule
	if d := strings.ToLower(g.Domain); d != "" {
		desired["Domain"] = d
	}
	addr := config.HostOnly(Utility.ToString(desired["Address"]))
	if addr == "" {
		addr = g.localIPAddress
		if addr == "" {
			addr = "127.0.0.1"
		}
	}
	desired["Address"] = addr

	// Ensure proxy defaults if not set
	if Utility.ToInt(desired["Proxy"]) == 0 {
		desired["Proxy"] = Utility.ToInt(desired["Port"]) + 1
	}
	return desired, nil
}

func (g *Globule) topoOrder(desiredByID map[string]map[string]interface{}) ([]string, error) {
	// Reuse your existing OrderDependencys after rebuilding the slice
	svcs := make([]map[string]interface{}, 0, len(desiredByID))
	for _, v := range desiredByID {
		svcs = append(svcs, v)
	}
	order, err := config.OrderDependencies(svcs)
	if err != nil {
		return nil, err
	}
	// Map names -> ids
	nameToID := map[string]string{}
	for id, s := range desiredByID {
		nameToID[Utility.ToString(s["Name"])] = id
	}
	out := make([]string, 0, len(order))
	for _, name := range order {
		if id := nameToID[name]; id != "" {
			out = append(out, id)
		}
	}
	return out, nil
}

func (g *Globule) waitServiceReady(name, addr string, total time.Duration) bool {
	if total <= 0 {
		total = 20 * time.Second
	}
	deadline := time.Now().Add(total)

	dials := []time.Duration{150 * time.Millisecond, 300 * time.Millisecond, 600 * time.Millisecond, 1200 * time.Millisecond, 2000 * time.Millisecond}
	i := 0

	for time.Now().Before(deadline) {
		// 1) TCP check
		if conn, err := net.DialTimeout("tcp", addr, 800*time.Millisecond); err == nil {
			_ = conn.Close()
			// 2) gRPC health
			if g.grpcHealthOK(addr) {
				return true
			}
			// 3) fallback to --health (best-effort)
			if g.binHealthOK(name) {
				return true
			}
		}
		// backoff
		wait := dials[i]
		if i < len(dials)-1 {
			i++
		}
		time.Sleep(wait)
	}
	return false
}

func (g *Globule) grpcHealthOK(addr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return false
	}
	defer cc.Close()
	hc := healthpb.NewHealthClient(cc)
	_, err = hc.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	return err == nil
}

func (g *Globule) binHealthOK(serviceName string) bool {

	root := config.GetServicesRoot()
	if root == "" {
		return false
	}
	// Use shared finder instead of ad-hoc walk.
	bin, err := config.FindServiceBinary(root, shortName(serviceName))
	if err != nil || strings.TrimSpace(bin) == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "--health")
	address, _ := config.GetAddress()
	address = strings.Split(address, ":")[0] // host only
	cmd.Env = append(os.Environ(),
		"GLOBULAR_DOMAIN="+g.Domain,
		"GLOBULAR_ADDRESS="+address,
	)
	return cmd.Run() == nil
}

// utilities
func toAnySlice(ss []string) []any {
	out := make([]any, 0, len(ss))
	for _, s := range ss {
		out = append(out, s)
	}
	return out
}
func toAnyAnySlice(v []interface{}) []any {
	out := make([]any, 0, len(v))
	for _, x := range v {
		out = append(out, x)
	}
	return out
}
func coalesce(s, def string) string {
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}
