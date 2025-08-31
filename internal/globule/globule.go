package globule

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/security"
	Utility "github.com/globulario/utility"
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
	IndexApplication  string // If defined It will be use as the entry point where not application path was given in the url.

	// TLS/ACME
	CertExpirationDelay        int
	CertPassword               string
	Country, State, City       string
	Organization               string
	Certificate                string
	CertificateAuthorityBundle string
	CertURL, CertStableURL     string
	registration               *registrationResource // wrapper defined in acme.go
	creds                      string                // config/tls/<domain>

	// Admin / auth
	AdminEmail     string
	RootPassword   string
	SessionTimeout int

	// Services / peers
	ServicesRoot   string
	BackendPort    int
	BackendStore   int
	ReverseProxies []interface{}
	peers          *sync.Map // runtime set of peers
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
}

// New creates a minimally initialized Globule.
// HTTP handlers/mux are wired elsewhere; we don’t touch net/http servers here.
func New(logger *slog.Logger) *Globule {

	if logger == nil {
		logger = slog.Default()
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

	// get the actual config and
	// Initialyse globular from it configuration file.
	file, err := os.ReadFile(g.configDir + "/config.json")

	// Init the service with the default port addres
	if err == nil {
		err := json.Unmarshal(file, &g)
		if err != nil {
			fmt.Println("fail to init configuation with error ", err)
			os.Exit(1)
		}
	}

	// Identity and addresses
	g.Mac, _ = config.GetMacAddress()
	g.localIPAddress, _ = Utility.MyLocalIP(g.Mac)
	g.ExternalIPAddress = Utility.MyIP()
	g.Name, _ = config.GetName()
	g.Domain, _ = config.GetDomain()
	if g.Domain == "" {
		g.Domain = "localhost"
	}
	if g.DNS == "" {
		g.DNS = g.localDomain()
	}

	return g
}

// InitFS initializes required directories and reads/saves config.json.
// Idempotent: safe to call at startup.
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

	// Here I will generate the keys for this server if not already exist.
	err := security.GeneratePeerKeys(g.Mac)
	if err != nil {
		fmt.Println("fail to generate peer keys with error ", err)
		return err
	}

	return g.SaveConfig()
}

// RegisterIPToDNS is kept public so you can call it on a cron/loop if you want.
func (g *Globule) RegisterIPToDNS() error { return g.registerIPToDNS() }

// StartServices boots microservices (ports from PortsRange) + their proxies.
func (g *Globule) StartServices(ctx context.Context) error { return g.startServices(ctx) }

// StopServices stops microservices and proxies.
func (g *Globule) StopServices() error { return g.stopServices() }

// WatchConfig hot-reloads the process-level config.json into Globule.
func (g *Globule) WatchConfig() { g.watchConfig() }

// InitPeers initializes peer map and subscribes to updates.
func (g *Globule) InitPeers() error { return g.initPeers() }

// Publish event convenience
func (g *Globule) Publish(evt string, data []byte) error {
	return g.publish(evt, data)
}

// Helpers
func (g *Globule) localDomain() string {
	addr, _ := config.GetAddress()
	return strings.Split(addr, ":")[0]
}

func (g *Globule) LocalDomain() string {
	return g.localDomain()
}

// Add near other methods

// BootstrapTLSAndDNS initializes FS, brings up DNS (if local), registers IPs,
// ensures CSR/key exist, and obtains/renews the certificate.
func (g *Globule) BootstrapTLSAndDNS(ctx context.Context) error {
	// 1) Ensure dirs/config
	if err := g.InitFS(); err != nil {
		return err
	}

	// 2) Start DNS locally if present + register A/AAAA/MX/TXT
	if err := g.maybeStartDNSAndRegister(ctx); err != nil {
		g.log.Warn("tls/dns bootstrap: dns bootstrap failed", "err", err)
	}

	// 3) Make sure we have client key + server CSR on disk
	if err := g.ensureAccountKeyAndCSR(); err != nil {
		return err
	}

	// 3.5) If a valid fullchain already exists, skip obtain and just schedule renewal
	fullchain := filepath.Join(config.GetConfigDir(), "tls", g.localDomain(), "fullchain.pem")
	if Utility.Exists(fullchain) {
		if notAfter, err := readCertNotAfter(fullchain); err == nil {
			now := time.Now()
			if now.Before(notAfter) {
				g.log.Info("tls: existing certificate found", "path", fullchain, "expires", notAfter.Local())
				// Schedule renewal before expiry (see helper for exact lead time)
				g.scheduleCertRenewal(ctx, fullchain, notAfter)
				return nil
			}
			g.log.Warn("tls: existing certificate is expired; will renew now", "expiredAt", notAfter.Local())
		} else {
			g.log.Warn("tls: could not parse existing fullchain; will attempt renew", "err", err)
		}
	}

	// 4) Obtain/renew cert (DNS-01 preferred if g.DNS != "")
	if err := g.obtainCertificateForCSR(ctx); err != nil {
		return err
	}

	// 4.1) Build/refresh fullchain.pem from issued files (best effort)
	if err := g.buildFullChainIfPossible(); err != nil {
		g.log.Warn("tls: failed to build fullchain.pem", "err", err)
	} else {
		// Schedule renewal based on the new fullchain
		if notAfter, err := readCertNotAfter(fullchain); err == nil {
			g.scheduleCertRenewal(ctx, fullchain, notAfter)
		}
	}

	return nil
}

// readCertNotAfter returns the NotAfter time from the first CERTIFICATE in a PEM file.
func readCertNotAfter(pemPath string) (time.Time, error) {
	b, err := os.ReadFile(pemPath)
	if err != nil {
		return time.Time{}, err
	}
	var block *pem.Block
	for {
		block, b = pem.Decode(b)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return time.Time{}, err
			}
			return cert.NotAfter, nil
		}
	}
	return time.Time{}, fmt.Errorf("no CERTIFICATE block in %s", pemPath)
}

// scheduleCertRenewal sets a background timer to renew the certificate shortly before it expires.
// It chooses a lead time of min(30d, max(1d, 10%% of lifetime)).
func (g *Globule) scheduleCertRenewal(ctx context.Context, fullchain string, notAfter time.Time) {
	now := time.Now()
	lifetime := notAfter.Sub(now)
	if lifetime <= 0 {
		// already expired; trigger now in a new goroutine
		go func() {
			if err := g.obtainCertificateForCSR(context.Background()); err != nil {
				g.log.Error("tls: auto-renew failed", "err", err)
				return
			}
			_ = g.buildFullChainIfPossible()
		}()
		return
	}

	// compute lead time
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
		fireAt = now.Add(30 * time.Second) // don't schedule in the past / too soon
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
			if err := g.obtainCertificateForCSR(context.Background()); err != nil {
				g.log.Error("tls: auto-renew failed", "err", err)
				return
			}
			if err := g.buildFullChainIfPossible(); err != nil {
				g.log.Warn("tls: post-renew fullchain build failed", "err", err)
			}
			// reschedule with the new expiry
			full := filepath.Join(config.GetConfigDir(), "tls", g.localDomain(), "fullchain.pem")
			if na, err := readCertNotAfter(full); err == nil {
				g.scheduleCertRenewal(context.Background(), full, na)
			}
		}
	}()
}

// buildFullChainIfPossible concatenates <domain>.crt + <domain>.issuer.crt into fullchain.pem (best effort).
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
	if err := os.WriteFile(full, append(leafB, issuerB...), 0600); err != nil {
		return fmt.Errorf("write fullchain.pem: %w", err)
	}
	g.log.Info("tls: fullchain.pem updated", "path", full)
	return nil
}
