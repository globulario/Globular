package globule

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/globulario/Globular/internal/logsink"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/pki"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	"github.com/gookit/color"
)

// Supervisor in main.go runs HTTP/HTTPS. Globule now focuses on:
// - Directories/config and desired-state seeds
// - DNS/IP registration and PKI/TLS bootstrapping
// - Peers & events (NodeAgent/Controller own lifecycle)

type Globule struct {
	// Identity / network
	Name              string
	Mac               string
	Domain            string
	AlternateDomains  []interface{}
	Protocol          string
	PortHTTP          int
	PortHTTPS         int
	PortsRange        string
	localIPAddress    string
	ExternalIPAddress string
	IndexApplication  string

	// TLS/ACME
	CertExpirationDelay        int
	CertPassword               string
	Country, State, City       string
	Organization               string
	Certificate                string
	CertificateAuthorityBundle string
	CertURL, CertStableURL     string
	creds                      string // config/tls/<domain>

	// Admin / auth
	AdminEmail     string
	RootPassword   string
	SessionTimeout int

	// Services / nodes
	ServicesRoot   string
	BackendPort    int
	BackendStore   int
	ReverseProxies []interface{}
	nodes          *sync.Map

	// Discovery / DNS
	DNS              string
	NS               []interface{}
	DNSUpdateIPInfos []interface{}
	SkipLocalDNS     bool
	MutateHostsFile  bool
	MutateResolvConf bool
	dnsRetryCancel   context.CancelFunc

	EnableConsoleLogs bool
	EnablePeerUpserts bool

	// OAuth2 configuration.
	OAuth2ClientID     string
	OAuth2ClientSecret string
	OAuth2RedirectURI  string

	// CORS
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

	// Versioning / updates
	Version          string
	Build            int
	Platform         string
	WatchUpdateDelay int64

	// Directories
	Path, webRoot, data, configDir string

	// lifecycle
	startTime time.Time

	log *slog.Logger

	stopConsole func()

	// Use Envoy proxy for service communication
	UseEnvoy bool
}

// New creates a minimally initialized Globule.
// HTTP handlers/mux are wired elsewhere; we don’t touch net/http servers here.
func New(logger *slog.Logger) *Globule {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	}

	g := &Globule{
		log:                 logger,
		startTime:           time.Now(),
		Version:             "1.0.0",
		Build:               0,
		Platform:            runtime.GOOS + ":" + runtime.GOARCH,
		PortHTTP:            8080,
		PortHTTPS:           8181,
		PortsRange:          "10000-10100",
		AllowedOrigins:      []string{"*"},
		AllowedMethods:      []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:      []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "domain", "application", "token", "video-path", "index-path", "routing"},
		Protocol:            "http",
		CertExpirationDelay: 365,
		CertPassword:        "1111",
		AdminEmail:          "sa@globular.cloud",
		RootPassword:        "adminadmin",
		SessionTimeout:      15,
		WatchUpdateDelay:    30,
		ServicesRoot:        config.GetServicesRoot(),
		nodes:               &sync.Map{},
		configDir:           config.GetRuntimeConfigDir(),
		UseEnvoy:            false,
	}

	configPaths := []string{
		filepath.Join(g.configDir, "config.json"),
		filepath.Join(config.GetConfigDir(), "config.json"),
	}
	for _, cfgPath := range configPaths {
		if g.loadConfigFile(cfgPath) {
			break
		}
	}

	// Identity and addresses
	g.Name, _ = config.GetName()
	g.Domain, _ = config.GetDomain()
	if g.Domain == "" {
		g.Domain = "localhost"
	}
	if g.DNS == "" {
		g.DNS = g.localDomain()
	}

	// Runtime only
	g.Mac, _ = config.GetMacAddress()
	g.localIPAddress, _ = Utility.MyLocalIP(g.Mac)
	g.ExternalIPAddress = Utility.MyIP()

	// Banner
	PrintBanner(g.Version, g.Build)

	return g
}

func (g *Globule) loadConfigFile(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	if err := json.Unmarshal(data, g); err != nil {
		g.log.Warn("failed to load config", "path", path, "err", err)
		return false
	}
	return true
}

// InitFS initializes required directories and reads/saves config.json.
func (g *Globule) InitFS() error {
	// Data / webroot / config roots
	g.data = config.GetDataDir()
	g.webRoot = config.GetWebRootDir()
	g.configDir = config.GetRuntimeConfigDir()
	g.Path, _ = os.Executable()
	stateRoot := config.GetStateRootDir()
	tokensDir := config.GetTokensDir()
	keysDir := config.GetKeysDir()
	tlsDir := config.GetRuntimeTLSDir()
	g.log.Info("runtime paths",
		"stateRoot", stateRoot,
		"runtimeConfig", g.configDir,
		"tokensDir", tokensDir,
		"keysDir", keysDir,
		"tlsDir", tlsDir,
		"adminConfigDir", config.GetConfigDir(),
	)

	for _, d := range []string{
		g.data, g.webRoot,
		g.configDir,
		tokensDir,
		keysDir,
		tlsDir,
	} {
		if err := config.EnsureRuntimeDir(d); err != nil {
			return err
		}
	}

	// TLS creds dir: runtime tls/<domain>
	g.creds = filepath.Join(tlsDir, g.localDomain())

	// Persist defaults on first run if config missing (single save later).
	cfgPath := filepath.Join(g.configDir, "config.json")
	needsSave := !Utility.Exists(cfgPath)

	// Ensure we have a stable identity before generating keys.
	changedID, err := g.ensureNodeID()
	if err != nil {
		g.log.Error("ensure node id", "err", err)
		return err
	}
	if changedID {
		needsSave = true
	}

	// Generate peer keys for this server if missing
	if err := security.GeneratePeerKeys(g.Mac); err != nil {
		g.log.Error("fail to generate peer keys", "err", err)
		return err
	}
	if needsSave {
		if err := g.SaveConfig(); err != nil {
			g.log.Error("fail to save config file", "err", err)
			return err
		}
	}

	return nil
}

func (g *Globule) ensureNodeID() (bool, error) {
	if id := strings.TrimSpace(g.Mac); id != "" {
		g.Mac = id
		return false, nil
	}

	if env := strings.TrimSpace(os.Getenv("GLOBULAR_NODE_ID")); env != "" {
		g.log.Info("using GLOBULAR_NODE_ID fallback", "id", env)
		g.Mac = env
		return true, nil
	}

	if hostname, err := os.Hostname(); err == nil {
		if hostID := strings.TrimSpace(hostname); hostID != "" {
			g.log.Info("using hostname fallback", "id", hostID)
			g.Mac = hostID
			return true, nil
		}
	} else {
		g.log.Warn("hostname fallback failed", "err", err)
	}

	if machineID, err := machineIDDerivedID(); err == nil {
		g.log.Info("using machine-id fallback", "id", machineID)
		g.Mac = machineID
		return true, nil
	} else {
		g.log.Warn("machine-id fallback failed", "err", err)
	}

	return false, errors.New("node id is empty after evaluating fallback sources")
}

func machineIDDerivedID() (string, error) {
	data, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}
	machineID := strings.TrimSpace(string(data))
	if machineID == "" {
		return "", errors.New("machine-id is empty")
	}
	sum := sha256.Sum256([]byte(machineID))
	return fmt.Sprintf("%x", sum[:8]), nil
}

// RegisterIPToDNS is kept public so you can call it on a cron/loop if you want.
func (g *Globule) RegisterIPToDNS(ctx context.Context) error {
	_, err := g.registerIPToDNS(ctx)
	return err
}

// NotifyNodeAgentReconcile records that desired configs changed and NodeAgent should reconcile.
func (g *Globule) NotifyNodeAgentReconcile(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	g.log.Info("desired configs updated; reconcile requested", "node_agent", NodeAgentAddress())
	// TODO: call NodeAgent reconcile RPC when available.
	_ = ctx
}

// WatchConfig hot-reloads the process-level config.json into Globule.
func (g *Globule) WatchConfig() { g.watchConfig() }

// InitNodes initializes node identity map and subscribes to updates.
func (g *Globule) InitNodes() error { return g.initNodes() }

// Publish event convenience
func (g *Globule) Publish(evt string, data []byte) error { return g.publish(evt, data) }

// Helpers
func (g *Globule) localDomain() string {
	addr, _ := config.GetAddress()
	if host := config.HostOnly(addr); host != "" {
		return host
	}
	if host, _, err := net.SplitHostPort(addr); err == nil && host != "" {
		return host
	}
	return addr
}
func (g *Globule) getAddress() string {
	addr, _ := config.GetAddress()
	return addr
}
func (g *Globule) LocalDomain() string { return g.localDomain() }

// ensureAccountKeyAndCSR creates/loads the legacy account key (client.pem) and a server CSR.
func (g *Globule) ensureAccountKeyAndCSR() error {
	creds := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain())
	if err := config.EnsureRuntimeDir(creds); err != nil {
		return err
	}

	// Account key (client.pem) — ensure via PKI
	acct := filepath.Join(creds, "client.pem")
	m := g.newPKIManager()
	if !Utility.Exists(acct) {
		_, _, _, _ = m.EnsureClientCert(creds, "acme-account", nil, 24*time.Hour) // local CA client auth cert (account key managed by PKI)
	}

	// Maintain historical CSR files for compatibility
	sk := filepath.Join(creds, "server.key")
	csr := filepath.Join(creds, "server.csr")
	if !Utility.Exists(sk) || !Utility.Exists(csr) {
		// DNS SAN CSR (legacy filenames)
		if err := m.EnsureServerKeyAndCSR(creds, g.Domain, g.Country, g.State, g.City, g.Organization, g.allDNS()); err != nil {
			return fmt.Errorf("ensure server csr: %w", err)
		}
	}
	g.creds = creds
	return nil
}

// obtainInternalServerCert ensures server.crt is signed by the local CA only (mTLS for etcd). // <<< NEW
func (g *Globule) obtainInternalServerCert(ctx context.Context) error { // <<< NEW
	m := g.newPKIManager()
	dir := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain())
	if err := config.EnsureRuntimeDir(dir); err != nil {
		return err
	}

	// Always local CA for server.crt (internal mTLS). The PKI EnsureServerCert now uses local CA here.
	_, leaf, ca, err := m.EnsureServerCert(dir, g.Domain, g.allDNS(), 0)
	if err != nil {
		return err
	}
	_ = ca

	// Normalize / enforce canonical path names (server.crt stays where PKI wrote it).
	if !Utility.Exists(leaf) {
		return fmt.Errorf("internal server.crt not found after issuance")
	}
	g.log.Info("tls: internal server.crt (local CA) ready", "path", leaf)
	return nil
}

// obtainPublicACMECert writes <domain>.crt, <domain>.issuer.crt, and fullchain.pem (ACME). // <<< NEW
func (g *Globule) obtainPublicACMECert(ctx context.Context) error { // <<< NEW
	m := g.newPKIManager()
	dir := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain())
	if err := config.EnsureRuntimeDir(dir); err != nil {
		return err
	}

	// Issue public cert to domain-named files; does NOT touch server.crt. (Requires PKI EnsurePublicACMECert impl.)
	key, leaf, issuer, full, err := m.EnsurePublicACMECert(
		dir, g.Domain, g.Domain, g.allDNS(), 0,
	)
	if err != nil {
		return err
	}
	_ = key

	// If EnsurePublicACMECert didn’t build fullchain, build it here (leaf + issuer).
	if !Utility.Exists(full) && Utility.Exists(leaf) && Utility.Exists(issuer) {
		if err := g.buildFullChainIfPossible(); err != nil {
			g.log.Warn("tls: failed to build fullchain.pem", "err", err)
		}
	}

	// Schedule auto-renew for the public chain.
	fullchain := filepath.Join(dir, "fullchain.pem")
	if Utility.Exists(fullchain) {
		if na, err := readCertNotAfter(fullchain); err == nil {
			g.scheduleCertRenewal(ctx, fullchain, na)
		}
	}

	g.log.Info("tls: public ACME certificate ready", "leaf", leaf)
	return nil
}

// newPKIManager builds a PKI manager configured from globule settings.
func (g *Globule) newPKIManager() pki.Manager {
	opts := pki.Options{
		Storage: pki.FileStorage{},
		ACME: pki.ACMEConfig{
			Enabled:   true,
			Email:     g.AdminEmail,
			Directory: "",         // prod
			Provider:  "globular", // use your DNS service
			DNS:       g.DNS,      // e.g. "10.0.0.63:10033"
			Domain:    g.Domain,
			Timeout:   2 * time.Minute,
		},
		LocalCA: pki.LocalCAConfig{
			Enabled:   true,
			Password:  g.CertPassword,
			Country:   g.Country,
			State:     g.State,
			City:      g.City,
			Org:       g.Organization,
			ValidDays: g.CertExpirationDelay,
		},
		Logger: g.log,

		// Called by the PKI manager when the "globular" DNS-01 provider needs an auth token.
		TokenSource: func(ctx context.Context, dnsAddr string) (string, error) {
			cli, err := dns_client.NewDnsService_Client(dnsAddr, "dns.DnsService")
			if err != nil {
				return "", err
			}
			defer cli.Close()
			mac := cli.GetMac()
			return security.GenerateToken(g.SessionTimeout, mac, "sa", "", g.AdminEmail, g.Domain)
		},
	}
	return pki.NewFileManager(opts)
}

// readCertNotAfter returns NotAfter of first CERTIFICATE block.
func readCertNotAfter(pemPath string) (time.Time, error) {
	b, err := os.ReadFile(pemPath)
	if err != nil {
		return time.Time{}, err
	}
	for {
		var block *pem.Block
		block, b = pem.Decode(b)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			c, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return time.Time{}, err
			}
			return c.NotAfter, nil
		}
	}
	return time.Time{}, errors.New("no CERTIFICATE in " + pemPath)
}

// allDNS returns primary + alternates (stringified), trimmed and deduped.
func (g *Globule) allDNS() []string {
	m := map[string]bool{}
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s != "" {
			m[strings.ToLower(s)] = true
		}
	}
	add(g.Domain)
	for _, v := range g.AlternateDomains {
		add(fmt.Sprint(v))
	}
	out := make([]string, 0, len(m))
	for s := range m {
		out = append(out, s)
	}
	return out
}

// BootstrapTLSAndDNS initializes FS, brings up DNS (if local), registers IPs,
// ensures CSR/key exist, ensures internal server.crt (local CA), and obtains/renews the public ACME cert.
func (g *Globule) BootstrapTLSAndDNS(ctx context.Context) error {
	// 1) Ensure dirs/config
	if err := g.InitFS(); err != nil {
		return err
	}

	// 2) Start DNS locally if present + register A/AAAA/MX/TXT
	if err := g.maybeStartDNSAndRegister(ctx); err != nil {
		g.log.Warn("tls/dns bootstrap: dns bootstrap failed", "err", err)
		return err
	}

	// 3) Make sure we have client key + server CSR on disk
	if err := g.ensureAccountKeyAndCSR(); err != nil {
		return err
	}

	// 3.1) Ensure internal server.crt is LOCAL-CA signed (mTLS). // <<< NEW
	if err := g.obtainInternalServerCert(ctx); err != nil { // <<< NEW
		g.log.Error("tls: internal server.crt issuance failed", "err", err)
		return err
	}

	// 3.5) If a valid fullchain already exists (public), just schedule renewal
	fullchain := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain(), "fullchain.pem")
	if Utility.Exists(fullchain) {
		if notAfter, err := readCertNotAfter(fullchain); err == nil {
			now := time.Now()
			if now.Before(notAfter) {
				g.log.Info("tls: existing public certificate found", "path", fullchain, "expires", notAfter.Local())
				g.scheduleCertRenewal(ctx, fullchain, notAfter)
				return nil
			}
			g.log.Warn("tls: existing public certificate is expired; will renew now", "expiredAt", notAfter.Local())
		} else {
			g.log.Warn("tls: could not parse existing fullchain; will attempt renew", "err", err)
		}
	}

	// 4) Obtain/renew public cert (DNS-01 preferred if g.DNS != "")
	if err := g.obtainPublicACMECert(ctx); err != nil { // <<< NEW
		return err
	}

	// 4.1) Build/refresh fullchain.pem from issued files (best effort)
	if err := g.buildFullChainIfPossible(); err != nil {
		g.log.Warn("tls: failed to build fullchain.pem", "err", err)
	} else {
		if na, err := readCertNotAfter(fullchain); err == nil {
			g.scheduleCertRenewal(ctx, fullchain, na)
		}
	}

	// 5) Promote etcd and proxies to TLS if certs are ready (idempotent)
	if err := g.promoteEtcdAndProxiesToTLS(ctx); err != nil {
		return err
	}

	return nil

}

// scheduleCertRenewal sets a background timer to renew the certificate shortly before it expires.
func (g *Globule) scheduleCertRenewal(ctx context.Context, fullchain string, notAfter time.Time) {
	now := time.Now()
	lifetime := notAfter.Sub(now)
	if lifetime <= 0 {
		// already expired; trigger now
		go func() {
			// Renew the PUBLIC ACME cert on expiry. // <<< NEW
			if err := g.obtainPublicACMECert(context.Background()); err != nil {
				g.log.Error("tls: auto-renew failed", "err", err)
				return
			}
			_ = g.buildFullChainIfPossible()
		}()
		return
	}

	// compute lead time: min(30d, max(1d, 10% of lifetime))
	tenPercent := time.Duration(float64(lifetime) * 0.10)
	lead := tenPercent
	if lead < 24*time.Hour {
		lead = 24 * time.Hour
	}
	if lead > 30*24*time.Hour {
		lead = 30 * 24 * time.Hour
	}

	fireAt := notAfter.Add(-lead)
	if fireAt.Before(now.Add(30 * time.Second)) {
		fireAt = now.Add(30 * time.Second)
	}

	g.log.Info("tls: scheduling auto-renew", "when", fireAt.Local(), "lead", lead)

	go func() {
		t := time.NewTimer(time.Until(fireAt))
		defer t.Stop()
		select {
		case <-ctx.Done():
			g.log.Info("tls: auto-renew canceled (context done)")
			return
		case <-t.C:
			// Renew the PUBLIC ACME cert on schedule. // <<< NEW
			if err := g.obtainPublicACMECert(context.Background()); err != nil {
				g.log.Error("tls: auto-renew failed", "err", err)
				return
			}
			if err := g.buildFullChainIfPossible(); err != nil {
				g.log.Warn("tls: post-renew fullchain build failed", "err", err)
			}
			full := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain(), "fullchain.pem")
			if na, err := readCertNotAfter(full); err == nil {
				g.scheduleCertRenewal(context.Background(), full, na)
			}
		}
	}()
}

// buildFullChainIfPossible concatenates <domain>.crt + <domain>.issuer.crt into fullchain.pem.
func (g *Globule) buildFullChainIfPossible() error {
	dir := filepath.Join(config.GetRuntimeTLSDir(), g.localDomain())
	if err := config.EnsureRuntimeDir(dir); err != nil {
		return err
	}
	leaf := filepath.Join(dir, g.Domain+".crt")
	issuer := filepath.Join(dir, g.Domain+".issuer.crt")
	full := filepath.Join(dir, "fullchain.pem")

	leafB, err := os.ReadFile(leaf)
	if err != nil {
		return fmt.Errorf("read leaf: %w", err)
	}
	issuerB, err := os.ReadFile(issuer)
	if err != nil {
		return fmt.Errorf("read issuer: %w", err)
	}
	if err := os.WriteFile(full, append(leafB, issuerB...), 0o600); err != nil {
		return fmt.Errorf("write fullchain.pem: %w", err)
	}
	g.log.Info("tls: fullchain.pem updated", "path", full)
	return nil
}

// console logs
func (g *Globule) startConsoleLogs() {
	sink := logsink.NewConsoleSink(g.getAddress(), logsink.Filter{
		MinLevel:   logpb.LogLevel_INFO_MESSAGE,
		ShowFields: true,
	})
	stop, _ := sink.Start()
	g.stopConsole = stop
	fmt.Println("console sink enabled (will retry until EventService is ready)")
}
func (g *Globule) stopConsoleLogs() {
	if g.stopConsole != nil {
		g.stopConsole()
		g.stopConsole = nil
	}
}

// PrintBanner prints the ASCII logo with Globular metadata.
func PrintBanner(version string, build int) {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	now := time.Now().Format(time.RFC1123)

	lines := []string{
		fmt.Sprintf("Globular v%s (build %d)", version, build),
		fmt.Sprintf("Platform: %s", platform),
		fmt.Sprintf("PID: %d", os.Getpid()),
		fmt.Sprintf("Started at: %s", now),
	}
	for _, l := range lines {
		color.Green.Println(l)
	}
	fmt.Println(strings.Repeat("-", 100))
}

// promoteEtcdAndProxiesToTLS restarts etcd on TLS if certs are ready,
// persists Protocol="https", and asks proxies to flip to TLS.
func (g *Globule) promoteEtcdAndProxiesToTLS(ctx context.Context) error {
	g.Protocol = "https"
	if err := g.SaveConfig(); err != nil {
		g.log.Warn("failed to persist Protocol=https", "err", err)
		return err
	}
	g.log.Info("promoteEtcdAndProxiesToTLS: Protocol set to https; NodeAgent should reconcile TLS")
	g.NotifyNodeAgentReconcile(ctx)
	// TODO: notify node-agent that TLS config changed/reconciliation required.
	return nil
}
