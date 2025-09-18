package globule

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
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
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
	"github.com/gookit/color"
	"github.com/kardianos/service"
)

// Supervisor in main.go runs HTTP/HTTPS. Globule now focuses on:
// - Directories/config
// - DNS/IP registration
// - ACME/lego certificate management
// - Microservice lifecycle
// - Peers & events

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

	// Services / peers
	ServicesRoot   string
	BackendPort    int
	BackendStore   int
	ReverseProxies []interface{}
	peers          *sync.Map
	Peers          []interface{}

	// Discovery / DNS
	DNS              string
	NS               []interface{}
	DNSUpdateIPInfos []interface{}

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
	Discoveries      []string
	WatchUpdateDelay int64

	// Directories
	Path, webRoot, data, configDir string
	users, applications            string

	// lifecycle
	logger    service.Logger
	startTime time.Time

	// log/signal
	exit   chan bool
	isExit bool

	log *slog.Logger

	stopConsole func()
}

// registrationResource kept only to avoid breaking struct fields;
// lego handles the real registration inside the PKI manager.
type registrationResource struct{}

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
		PortHTTP:            80,
		PortHTTPS:           443,
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
		peers:               &sync.Map{},
		configDir:           config.GetConfigDir(),
	}

	// Load existing config (if present)
	if data, err := os.ReadFile(g.configDir + "/config.json"); err == nil {
		if err := json.Unmarshal(data, &g); err != nil {
			fmt.Println("failed to load config:", err)
			os.Exit(1)
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

	// Ensure sensible defaults
	g.SaveConfig()

	// Runtime only
	g.Mac, _ = config.GetMacAddress()
	g.localIPAddress, _ = Utility.MyLocalIP(g.Mac)
	g.ExternalIPAddress = Utility.MyIP()

	// Banner
	PrintBanner(g.Version, g.Build)

	return g
}

// InitFS initializes required directories and reads/saves config.json.
func (g *Globule) InitFS() error {
	// Data / webroot / config roots
	g.data = config.GetDataDir()
	g.webRoot = config.GetWebRootDir()
	g.configDir = config.GetConfigDir()
	g.Path, _ = os.Executable()

	for _, d := range []string{
		g.data, g.webRoot,
		filepath.Join(g.data, "files", "templates"),
		filepath.Join(g.data, "files", "projects"),
		g.configDir,
		filepath.Join(g.configDir, "tokens"),
	} {
		if err := Utility.CreateDirIfNotExist(d); err != nil {
			return err
		}
	}

	// TLS creds dir: config/tls/<domain>
	g.creds = filepath.Join(g.configDir, "tls", g.localDomain())
	if err := Utility.CreateDirIfNotExist(g.creds); err != nil {
		return err
	}

	// Users / apps
	g.users = filepath.Join(g.data, "files", "users")
	g.applications = filepath.Join(g.data, "files", "applications")
	if err := Utility.CreateDirIfNotExist(g.users); err != nil {
		return err
	}
	if err := Utility.CreateDirIfNotExist(g.applications); err != nil {
		return err
	}

	// Persist defaults on first run
	cfgPath := filepath.Join(g.configDir, "config.json")
	if !Utility.Exists(cfgPath) {
		if err := g.SaveConfig(); err != nil {
			return err
		}
	}

	// Generate peer keys for this server if missing
	if err := security.GeneratePeerKeys(g.Mac); err != nil {
		g.log.Error("fail to generate peer keys", "err", err)
		return err
	}
	if err := g.SaveConfig(); err != nil {
		g.log.Error("fail to save config file", "err", err)
		return err
	}

	// Start the etcd server before starting other services
	if err := process.StartEtcdServer(); err != nil {
		g.log.Error("fail to start etcd server", "err", err)
		return err
	}
	return nil
}

// RegisterIPToDNS is kept public so you can call it on a cron/loop if you want.
func (g *Globule) RegisterIPToDNS() error { return g.registerIPToDNS() }

// StartServices is the public entry used by main.go.
func (g *Globule) StartServices(ctx context.Context) {

	defer g.startConsoleLogs()
	if err := g.startServicesEtcd(ctx); err != nil {
		g.log.Error("StartServices failed", "err", err)
		os.Exit(1)
	}
}

// StopServices is the public shutdown entry used by main.go.
func (g *Globule) StopServices() {
	defer g.stopConsoleLogs()
	if err := g.stopServicesEtcd(); err != nil {
		g.log.Error("StopServices failed", "err", err)
	}
}

// WatchConfig hot-reloads the process-level config.json into Globule.
func (g *Globule) WatchConfig() { g.watchConfig() }

// InitPeers initializes peer map and subscribes to updates.
func (g *Globule) InitPeers() error { return g.initPeers() }

// Publish event convenience
func (g *Globule) Publish(evt string, data []byte) error { return g.publish(evt, data) }

// Helpers
func (g *Globule) localDomain() string {
	addr, _ := config.GetAddress()
	return strings.Split(addr, ":")[0]
}
func (g *Globule) getAddress() string {
	addr, _ := config.GetAddress()
	return addr
}
func (g *Globule) LocalDomain() string { return g.localDomain() }

// ensureAccountKeyAndCSR creates/loads the legacy account key (client.pem) and a server CSR.
func (g *Globule) ensureAccountKeyAndCSR() error {
	creds := filepath.Join(config.GetConfigDir(), "tls", g.localDomain())
	if err := Utility.CreateDirIfNotExist(creds); err != nil {
		return err
	}

	// Account key (client.pem) — ensure via PKI
	acct := filepath.Join(creds, "client.pem")
	m := g.newPKIManager()
	if !Utility.Exists(acct) {
		_, _, _, _ = m.EnsureClientCert(creds, "acme-account", nil, 24*time.Hour) // local CA client auth cert (account key managed by PKI) :contentReference[oaicite:2]{index=2}
	}

	// Maintain historical CSR files for compatibility
	sk := filepath.Join(creds, "server.key")
	csr := filepath.Join(creds, "server.csr")
	if !Utility.Exists(sk) || !Utility.Exists(csr) {
		// DNS SAN CSR (legacy filenames)
		if err := m.EnsureServerKeyAndCSR(creds, g.Domain, g.Country, g.State, g.City, g.Organization, g.allDNS()); err != nil { // :contentReference[oaicite:3]{index=3}
			return fmt.Errorf("ensure server csr: %w", err)
		}
	}
	g.creds = creds
	return nil
}

// obtainInternalServerCert ensures server.crt is signed by the local CA only (mTLS for etcd). // <<< NEW
func (g *Globule) obtainInternalServerCert(ctx context.Context) error { // <<< NEW
	m := g.newPKIManager()
	dir := filepath.Join(config.GetConfigDir(), "tls", g.localDomain())

	// Always local CA for server.crt (internal mTLS). The PKI EnsureServerCert now uses local CA here. :contentReference[oaicite:4]{index=4}
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
	dir := filepath.Join(config.GetConfigDir(), "tls", g.localDomain())

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
		os.Exit(1)
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
	fullchain := filepath.Join(config.GetConfigDir(), "tls", g.localDomain(), "fullchain.pem")
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
			full := filepath.Join(config.GetConfigDir(), "tls", g.localDomain(), "fullchain.pem")
			if na, err := readCertNotAfter(full); err == nil {
				g.scheduleCertRenewal(context.Background(), full, na)
			}
		}
	}()
}

// buildFullChainIfPossible concatenates <domain>.crt + <domain>.issuer.crt into fullchain.pem.
func (g *Globule) buildFullChainIfPossible() error {
	dir := filepath.Join(config.GetConfigDir(), "tls", g.localDomain())
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
		MinLevel:            logpb.LogLevel_INFO_MESSAGE,
		ShowFields:          true,
		BackfillSince:       3 * time.Minute,
		BackfillPerAppLimit: 250,
		BackfillApps:        []string{"dns.DnsService"},
		Apps:                map[string]bool{"dns.DnsService": true},
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
