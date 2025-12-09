package globule

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/globulario/Globular/internal/logsink"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

			if state == "starting" && g.recentlyStarting(int64(updated), 15*time.Second) {
				continue
			}

			if g.isActuallyRunning(name, port) {
				_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
				continue
			}

			go func(id string, delay time.Duration) {
				time.Sleep(delay)
				g.tryRestartWithBackoff(ctx, guard, id)
			}(id, time.Duration(50+rand.Intn(250))*time.Millisecond)
		}
	}

	// 2) Live watch
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

			if state == "starting" && g.recentlyStarting(int64(updated), 15*time.Second) {
				return
			}

			if g.isActuallyRunning(name, port) {
				_ = config.PutRuntime(ev.ID, map[string]any{"State": "running", "LastError": ""})
				return
			}

			if proc == -1 || state == "failed" || state == "stopped" {
				go g.tryRestartWithBackoff(ctx, guard, ev.ID)
			}
		})
	}()
}

var (
	envoySnapMu    sync.Mutex
	envoySnapTimer *time.Timer
)

func (g *Globule) refreshEnvoySnapshotDebounced() {
	envoySnapMu.Lock()
	defer envoySnapMu.Unlock()
	if envoySnapTimer != nil {
		envoySnapTimer.Stop()
	}
	envoySnapTimer = time.AfterFunc(500*time.Millisecond, func() {
		if err := g.SetSnapshot(); err != nil {
			g.log.Warn("envoy snapshot refresh failed", "err", err)
		} else {
			g.log.Info("envoy snapshot updated")
		}
	})
}

// bootstrapServiceConfigsFromDisk syncs service definition files into etcd
// without clobbering existing entries; only missing services are restored.
func (g *Globule) bootstrapServiceConfigsFromDisk() error {
	dir := config.GetServicesConfigDir()
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("bootstrap services: stat %s: %w", dir, err)
	}
	if !fi.IsDir() {
		return nil
	}

	existing := map[string]struct{}{}
	if svcs, err := config.GetServicesConfigurations(); err == nil {
		for _, svc := range svcs {
			id := strings.TrimSpace(Utility.ToString(svc["Id"]))
			if id != "" {
				existing[id] = struct{}{}
			}
		}
	} else {
		g.log.Warn("bootstrap services: list existing configs failed", "err", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("bootstrap services: readdir %s: %w", dir, err)
	}

	var loaded int
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}

		path := filepath.Join(dir, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			g.log.Warn("bootstrap services: read config failed", "file", path, "err", err)
			continue
		}

		var desired map[string]interface{}
		if err := json.Unmarshal(raw, &desired); err != nil {
			g.log.Warn("bootstrap services: parse config failed", "file", path, "err", err)
			continue
		}

		id := strings.TrimSpace(Utility.ToString(desired["Id"]))
		if id == "" {
			base := strings.TrimSuffix(name, filepath.Ext(name))
			id = base
			desired["Id"] = id
		}

		if id == "" {
			g.log.Warn("bootstrap services: missing service id", "file", path)
			continue
		}

		if _, ok := existing[id]; ok {
			continue
		}

		if err := config.SaveServiceConfiguration(desired); err != nil {
			g.log.Warn("bootstrap services: save config failed", "id", id, "file", path, "err", err)
			continue
		}

		existing[id] = struct{}{}
		loaded++
	}

	if loaded > 0 {
		g.log.Info("bootstrap services: registered missing service configs", "count", loaded)
	}
	return nil
}

func (g *Globule) tryRestartWithBackoff(ctx context.Context, guard *restartGuard, id string) {
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
		address = strings.Split(address, ":")[0]
		addr := net.JoinHostPort(address, fmt.Sprint(port))

		if g.isActuallyRunning(name, port) {
			_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
			return
		}

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

			if !g.UseEnvoy {
				if _, perr := process.StartServiceProxyProcess(cfg, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); perr != nil {
					g.log.Warn("keepalive: proxy start failed", "name", name, "err", perr)
				}
			}

			if g.waitServiceReady(name, addr, 8*time.Second) {
				_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
				if g.UseEnvoy {
					g.refreshEnvoySnapshotDebounced()
				}
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
	addr = strings.Split(addr, ":")[0]
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
	nowMs := time.Now().UnixMilli()
	tsMs := updatedAt
	if updatedAt < 2_000_000_000 {
		tsMs = updatedAt * 1000
	}
	return nowMs-tsMs < grace.Milliseconds()
}

// Use the strongly-typed descriptor from config.
type serviceDesc = config.ServiceDesc

func shortName(full string) string {
	if full == "" {
		return ""
	}
	return strings.ToLower(strings.Split(full, ".")[0])
}

func isListening(addr string, timeout time.Duration) bool {
	d := net.Dialer{Timeout: timeout}
	c, err := d.Dial("tcp", addr)
	if err != nil {
		return false
	}
	_ = c.Close()
	return true
}

// EnsureRolesWithResource takes role definitions as []map[string]any and upserts
// them in the Resource service, keeping actions in sync (add missing, remove extras).
// address is the peer address (e.g., "<host>:<port>") hosting the Resource service.
func EnsureRolesWithResource(roleMaps []any) error {
	if len(roleMaps) == 0 {
		return nil
	}

	// Get resource service address
	address, _ := config.GetAddress()

	// Connect client
	rc, err := resource_client.NewResourceService_Client(address, "resource.ResourceService")
	if err != nil {
		return fmt.Errorf("resource client: %w", err)
	}
	defer rc.Close()

	// Local token for admin ops
	mac, _ := config.GetMacAddress()
	token, err := security.GetLocalToken(mac)
	if err != nil {
		return fmt.Errorf("get local token: %w", err)
	}

	// Fetch existing roles once
	existing, err := rc.GetRoles("") // fetch all; filter locally by Id
	if err != nil {
		return fmt.Errorf("get roles: %w", err)
	}
	exByID := make(map[string]*resourcepb.Role, len(existing))
	for _, r := range existing {
		exByID[r.Id] = r
	}

	for _, m := range roleMaps {
		role, e := mapToRole(m.(map[string]any))
		if e != nil {
			return e
		}
		if role.Id == "" {
			return errors.New("role missing Id")
		}

		// Normalize unique action list in target
		role.Actions = unique(role.Actions)

		// Does it exist already?
		if cur, ok := exByID[role.Id]; ok {
			// 1) Update metadata (name/description/domain)
			if err := rc.UpdateRole(string(token), &resourcepb.Role{
				Id:          role.Id,
				Name:        role.Name,
				Description: role.Description,
				Domain:      role.Domain,
			}); err != nil {
				return fmt.Errorf("update role %s: %w", role.Id, err)
			}

			// 2) Sync actions
			toAdd, toRemove := diffActions(cur.Actions, role.Actions)
			if len(toAdd) > 0 {
				if err := rc.AddRoleActions(token, cur.Id, toAdd); err != nil {
					return fmt.Errorf("add role actions %s: %w", cur.Id, err)
				}
			}
			for _, a := range toRemove {
				if err := rc.RemoveRoleAction(token, cur.Id, a); err != nil {
					return fmt.Errorf("remove role action %s %s: %w", cur.Id, a, err)
				}
			}
		} else {
			// Create with actions
			if err := rc.CreateRole(string(token), role.Id, role.Name, role.Actions); err != nil {
				return fmt.Errorf("create role %s: %w", role.Id, err)
			}
			// Then set description/domain if provided (UpdateRole updates metadata fields)
			if role.Description != "" || role.Domain != "" {
				if err := rc.UpdateRole(string(token), &resourcepb.Role{
					Id:          role.Id,
					Name:        role.Name,
					Description: role.Description,
					Domain:      role.Domain,
				}); err != nil {
					return fmt.Errorf("post-create update role %s: %w", role.Id, err)
				}
			}
		}

		// Now i will AddAccountRole to the admin user
		if err := rc.AddAccountRole(string(token), "sa", role.Id); err != nil {
			return fmt.Errorf("add account role %s: %w", role.Id, err)
		}
	}

	return nil
}

func mapToRole(m map[string]any) (*resourcepb.Role, error) {
	if m == nil {
		return nil, errors.New("nil role map")
	}
	role := &resourcepb.Role{
		Id:          Utility.ToString(m["id"]),
		Name:        Utility.ToString(m["name"]),
		Description: Utility.ToString(m["description"]),
		Domain:      Utility.ToString(m["domain"]),
		TypeName:    firstNonEmpty(Utility.ToString(m["typeName"]), "resource.Role"),
	}
	// actions can be []string or []any
	switch v := m["actions"].(type) {
	case []string:
		role.Actions = append([]string(nil), v...)
	case []interface{}:
		role.Actions = make([]string, 0, len(v))
		for _, it := range v {
			role.Actions = append(role.Actions, Utility.ToString(it))
		}
	case nil:
		// ok
	default:
		// try to stringify fallback
		s := Utility.ToString(v)
		if strings.TrimSpace(s) != "" {
			role.Actions = []string{s}
		}
	}
	return role, nil
}

// diffActions returns (toAdd, toRemove) to transform current -> target
func diffActions(current, target []string) (add []string, remove []string) {
	cur := make(map[string]bool, len(current))
	tgt := make(map[string]bool, len(target))
	for _, a := range current {
		cur[a] = true
	}
	for _, a := range target {
		tgt[a] = true
	}
	for a := range tgt {
		if !cur[a] {
			add = append(add, a)
		}
	}
	for a := range cur {
		if !tgt[a] {
			remove = append(remove, a)
		}
	}
	sort.Strings(add)
	sort.Strings(remove)
	return
}

func unique(in []string) []string {
	if len(in) == 0 {
		return in
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
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

	// 2) Discover executables
	bins, err := config.DiscoverExecutables(config.GetServicesRoot())
	if err != nil {
		return err
	}
	if len(bins) == 0 {
		return errors.New("no service executables found; check GetServicesRoot()")
	}

	// 3) Allocate ports
	alloc, err := config.NewPortAllocator(g.PortsRange)
	if err != nil {
		return err
	}
	if existing, _ := config.GetServicesConfigurations(); len(existing) > 0 {
		alloc.ReserveExisting(existing)
	}

	// 4) Describe -> desired
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

		if g.Protocol == "https" {
			// Ensure service has TLS certs if running in https mode
			if _, err := os.Stat(config.GetLocalServerKeyPath()); os.IsNotExist(err) {
				g.log.Error("missing TLS certificate", "service", desc.Name)
			} else {
				// The certificate private key
				m["KeyFile"] = config.GetLocalServerKeyPath()
				// The certificate chain to present to clients
				m["CertFile"] = config.GetLocalServerCertificatePath()
				// The CA bundle to validate client certs (if mTLS)
				m["CertAuthorityTrust"] = config.GetLocalCACertificate()
			}
		}

		if err := config.SaveServiceConfiguration(m); err != nil {
			g.log.Error("save desired failed", "service", desc.Name, "id", m["Id"], "err", err)
			continue
		}

		desiredByID[Utility.ToString(m["Id"])] = m
	}

	// --- Start MinIO (if enabled) and inject its settings into FileService, MediaService, and TorrentService ---
	minioCfg, err := g.startMinioIfNeeded(ctx, g.log)
	if err != nil {
		g.log.Error("failed to start MinIO", "err", err)
	} else if minioCfg != nil {
		for id, m := range desiredByID {
			name := Utility.ToString(m["Name"])
			switch {
			case strings.EqualFold(name, "file.FileService"):
				g.applyMinioConfigToService(id, m, minioCfg, "FileService", desiredByID)
			case strings.EqualFold(name, "media.MediaService"):
				g.applyMinioConfigToService(id, m, minioCfg, "MediaService", desiredByID)
			case strings.EqualFold(name, "torrent.TorrentService"):
				g.applyMinioConfigToService(id, m, minioCfg, "TorrentService", desiredByID)
			}
		}
	}

	// 5) Order by deps
	ordered, err := g.topoOrder(desiredByID)
	if err != nil {
		return err
	}
	if len(ordered) == 0 {
		return errors.New("no services to start after desired authoring")
	}

	// 6) Spawn services (no per-service proxy when UseEnvoy)
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

		// Skip if already running
		addr := address + ":" + Utility.ToString(m["Port"])
		if isListening(addr, 300*time.Millisecond) {
			g.log.Info("service already running; skipping start", "name", name, "addr", addr)
			_ = config.PutRuntime(id, map[string]any{"State": "running", "LastError": ""})
			continue
		}

		outW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stdout", logpb.LogLevel_INFO_MESSAGE, os.Stdout)
		errW := logsink.NewServiceLogWriter(address, name, "sa", "/"+name+"/stderr", logpb.LogLevel_ERROR_MESSAGE, os.Stderr)

		// ---- Inject bootstrap env just for the child, then restore ----
		restore := g.pushChildEnv(map[string]string{
			"GLOBULAR_DOMAIN":     strings.ToLower(Utility.ToString(m["Domain"])),
			"GLOBULAR_ADDRESS":    address,
			"GLOBULAR_SERVICE_ID": Utility.ToString(m["Id"]),
			"GLOBULAR_PORT":       fmt.Sprint(port),
			"GLOBULAR_PROXY":      fmt.Sprint(proxy),
		})

		g.log.Info("starting service", "name", name, "id", id, "port", port, "proxy", proxy)
		pid, err := process.StartServiceProcessWithWriters(m, port, outW, errW)
		restore() // always restore parent env

		if err != nil {
			g.log.Warn("service start failed", "name", name, "err", err)
			_ = config.PutRuntime(id, map[string]any{"Process": -1, "State": "failed", "LastError": err.Error()})
			continue
		}

		_ = config.PutRuntime(id, map[string]any{"Process": pid, "State": "starting", "LastError": ""})
		m["Process"] = pid

		if !g.UseEnvoy {
			// Legacy per-service proxy path
			if _, err := process.StartServiceProxyProcess(m, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate()); err != nil {
				g.log.Warn("proxy start failed", "name", name, "err", err)
			}
		}
	}

	// 7) Readiness wait
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

	// 8) If using Envoy, publish a fresh snapshot now that services are up.
	if g.UseEnvoy {
		if err := g.SetSnapshot(); err != nil {
			g.log.Warn("initial envoy snapshot failed", "err", err)
		}
	}

	// 9) Ensure default roles (if any) are present in Resource service
	// (do this after services are running so Resource is available)
	for _, id := range ordered {

		desc, _ := config.GetServiceConfigurationById(id)
		if desc["RolesDefault"] != nil {
			if err := EnsureRolesWithResource(desc["RolesDefault"].([]interface{})); err != nil {
				g.log.Warn("ensure roles failed", "service", desc["Name"], "err", err)
			}
		}

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

// normalizeAddr ensures "host:port", strips schemes, fixes IPv6 literals, and
// adds a sensible default port if missing (proxy vs direct).
func (g *Globule) normalizeAddr(addr string) (host, hostPort string) {
	a := strings.TrimSpace(addr)
	a = strings.TrimPrefix(a, "http://")
	a = strings.TrimPrefix(a, "https://")
	// If it's a bare IPv6 without brackets, add them later via netip parsing.
	// Add default port if missing.
	if !strings.Contains(a, ":") || (strings.Count(a, ":") > 1 && !strings.Contains(a, "]")) {
		// no explicit port, choose based on protocol
		defPort := "443"
		if !strings.EqualFold(g.Protocol, "https") {
			defPort = "80"
		}
		// If it *is* IPv6, we'll bracket it below.
		a = a + ":" + defPort
	}

	// Split host/port reliably; if it fails, try netip.
	if h, p, err := net.SplitHostPort(a); err == nil {
		host = h
		hostPort = net.JoinHostPort(h, p)
	} else {
		// Try to parse host as IP (IPv6 literal without brackets).
		// Take last colon as port sep.
		last := strings.LastIndexByte(a, ':')
		h := a[:last]
		p := a[last+1:]
		if ip, perr := netip.ParseAddr(h); perr == nil && ip.Is6() {
			host = ip.String()
			hostPort = net.JoinHostPort(host, p) // adds brackets
		} else {
			// Fallback: best effort
			host = h
			hostPort = h + ":" + p
		}
	}
	return host, hostPort
}

func (g *Globule) waitServiceReady(name, addr string, total time.Duration) bool {
	if total <= 0 {
		total = 20 * time.Second
	}
	deadline := time.Now().Add(total)

	backoff := []time.Duration{150 * time.Millisecond, 300 * time.Millisecond, 600 * time.Millisecond, 1200 * time.Millisecond, 2000 * time.Millisecond}
	i := 0

	host, hostPort := g.normalizeAddr(addr)
	dialer := &net.Dialer{Timeout: 800 * time.Millisecond}

	for time.Now().Before(deadline) {
		conn, err := dialer.Dial("tcp", hostPort)
		if err == nil {
			_ = conn.Close()

			if g.grpcHealthOK(hostPort) {
				return true
			}

			if strings.EqualFold(g.Protocol, "https") {
				// IMPORTANT: ensure SNI is a hostname even if we dial an IP
				sni := g.pickSNI(host)
				tlsCfg, terr := g.tlsConfigFor(sni)
				if terr == nil {
					if tconn, herr := tls.DialWithDialer(dialer, "tcp", hostPort, tlsCfg); herr == nil {
						_ = tconn.Close()
						if g.binHealthOK(name) {
							return true
						}
					}
				}
			} else {
				if g.binHealthOK(name) {
					return true
				}
			}
		}

		wait := backoff[i]
		if i < len(backoff)-1 {
			i++
		}
		time.Sleep(wait)
	}
	return false
}

func (g *Globule) binHealthOK(serviceName string) bool {
	root := config.GetServicesRoot()
	if root == "" {
		return false
	}
	bin, err := config.FindServiceBinary(root, shortName(serviceName))
	if err != nil || strings.TrimSpace(bin) == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "--health")
	address, _ := config.GetAddress()
	address = strings.Split(address, ":")[0]
	cmd.Env = append(os.Environ(),
		"GLOBULAR_DOMAIN="+g.Domain,
		"GLOBULAR_ADDRESS="+address,
	)
	return cmd.Run() == nil
}

// utilities toAnySlice, toAnyAnySlice, coalesce unchanged (omitted) ...
func looksLikeIP(h string) bool {
	if h == "" {
		return false
	}
	if ip := net.ParseIP(strings.Trim(h, "[]")); ip != nil {
		return true
	}
	// netip is stricter and faster for v6 literals
	if ip, err := netip.ParseAddr(strings.Trim(h, "[]")); err == nil && (ip.Is4() || ip.Is6()) {
		return true
	}
	return false
}

func (g *Globule) defaultSNI() string {
	// prefer node FQDN like "globule-ryzen.globular.io"
	if addr, _ := config.GetAddress(); addr != "" {
		host := strings.Split(addr, ":")[0]
		if host != "" && !looksLikeIP(host) && host != "localhost" {
			return host
		}
	}
	// try service FQDN
	if g.Domain != "" && g.Name != "" {
		return strings.ToLower(g.Name + "." + g.Domain)
	}
	return ""
}

func (g *Globule) pickSNI(serverNameHint string) string {
	// If caller passed a hostname (not an IP), use it.
	if serverNameHint != "" && !looksLikeIP(serverNameHint) && serverNameHint != "localhost" {
		return serverNameHint
	}
	// Otherwise use our best FQDN guess.
	if sni := g.defaultSNI(); sni != "" {
		return sni
	}
	// Last resort: whatever we got (may be IP; TLS will likely fail hostname verification).
	return serverNameHint
}

// ---------- gRPC health with proper SNI ----------
func (g *Globule) grpcHealthOK(addr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	host := addr
	if i := strings.IndexByte(addr, ':'); i > 0 {
		host = addr[:i]
	}

	var dialOpt grpc.DialOption
	if strings.EqualFold(g.Protocol, "https") {
		sni := g.pickSNI(host)
		tlsCfg, err := g.tlsConfigFor(sni)
		if err != nil {
			return false
		}
		dialOpt = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
	} else {
		dialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	cc, err := grpc.DialContext(ctx, addr, dialOpt, grpc.WithBlock())
	if err != nil {
		return false
	}
	defer cc.Close()

	hc := healthpb.NewHealthClient(cc)
	_, err = hc.Check(ctx, &healthpb.HealthCheckRequest{Service: ""})
	return err == nil
}

func (g *Globule) tlsConfigFor(serverName string) (*tls.Config, error) {
	roots, _ := x509.SystemCertPool()
	if roots == nil {
		roots = x509.NewCertPool()
	}
	caPath := config.GetLocalCACertificate()
	if data, err := os.ReadFile(caPath); err == nil {
		_ = roots.AppendCertsFromPEM(data)
	}

	tlsCfg := &tls.Config{
		ServerName: g.pickSNI(serverName),
		RootCAs:    roots,
		MinVersion: tls.VersionTLS12,
	}

	// Optional mTLS
	clientCert := g.creds + "/client.crt"
	clientKey := g.creds + "/client.pem"
	if _, err1 := os.Stat(clientCert); err1 == nil {
		if _, err2 := os.Stat(clientKey); err2 == nil {
			if cert, err := tls.LoadX509KeyPair(clientCert, clientKey); err == nil {
				tlsCfg.Certificates = []tls.Certificate{cert}
			}
		}
	}

	return tlsCfg, nil
}

// pushChildEnv temporarily sets a few env vars for the child process and returns a restore func.
func (g *Globule) pushChildEnv(kv map[string]string) func() {
	type prev struct {
		val string
		ok  bool
	}
	saved := map[string]prev{}
	for k, v := range kv {
		pv, ok := os.LookupEnv(k)
		saved[k] = prev{val: pv, ok: ok}
		_ = os.Setenv(k, v)
	}
	return func() {
		for k, p := range saved {
			if p.ok {
				_ = os.Setenv(k, p.val)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}
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

func (g *Globule) applyMinioConfigToService(id string, m map[string]interface{}, cfg *MinioRuntime, label string, desired map[string]map[string]interface{}) {
	if cfg == nil || m == nil {
		return
	}
	g.log.Info("configuring "+label+" to use MinIO",
		"serviceId", id,
		"endpoint", cfg.Endpoint,
		"bucket", cfg.Bucket,
		"prefix", cfg.Prefix,
	)

	m["UseMinio"] = true
	m["MinioEndpoint"] = cfg.Endpoint
	m["MinioAccessKey"] = cfg.AccessKey
	m["MinioSecretKey"] = cfg.SecretKey
	m["MinioBucket"] = cfg.Bucket
	m["MinioPrefix"] = cfg.Prefix
	m["MinioUseSSL"] = cfg.UseSSL

	if err := config.SaveServiceConfiguration(m); err != nil {
		g.log.Error("failed to save "+label+" MinIO config", "id", id, "err", err)
		return
	}
	if desired != nil {
		desired[id] = m
	}
}
