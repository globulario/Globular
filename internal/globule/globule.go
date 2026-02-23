package globule

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
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
	"github.com/gookit/color"
)

// Supervisor in main.go runs HTTP/HTTPS. Globule now focuses on:
// - Directories/config and desired-state seeds
// - Gateway configuration surface (no TLS/DNS ownership)
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
		log:            logger,
		startTime:      time.Now(),
		Version:        "1.0.0",
		Build:          0,
		Platform:       runtime.GOOS + ":" + runtime.GOARCH,
		PortHTTP:       8080,
		PortHTTPS:      8443,
		PortsRange:     "10000-10100",
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		// v1 Conformance: Remove "domain" header (security violation INV-1.7)
		// Domain is routing configuration, not authentication/authorization data
		// Allowing client-supplied "domain" header could influence identity decisions
		// REMOVED: "domain" from allowed headers list
		AllowedHeaders:      []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "application", "token", "video-path", "index-path", "routing"},
		Protocol:            "https", // Secure by default: HTTPS with TLS termination at Envoy
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
	if os.Getenv("GLOBULAR_SKIP_PEER_KEYS") != "1" {
		if err := security.GeneratePeerKeys(g.Mac); err != nil {
			g.log.Error("fail to generate peer keys", "err", err)
			return err
		}
	}

	// Ensure the local service token exists so that internal RBAC calls
	// (e.g. AddResourceOwner) can authenticate without requiring an
	// external token.  Use a one-year lifetime; GetLocalToken refreshes
	// automatically when the token is close to expiry.
	const localTokenLifetimeMinutes = 365 * 24 * 60
	if err := security.SetLocalToken(g.Mac, "sa", "sa", "", localTokenLifetimeMinutes); err != nil {
		// Non-fatal: log and continue.  Internal RBAC calls will fail until
		// the token is available, but the gateway can still serve traffic.
		g.log.Warn("failed to set local service token", "mac", g.Mac, "err", err)
	} else {
		g.log.Info("local service token refreshed", "mac", g.Mac)
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

// WatchConfig hot-reloads the process-level config.json into Globule.
func (g *Globule) WatchConfig() { g.watchConfig() }

// Publish event convenience
func (g *Globule) Publish(evt string, data []byte) error { return g.publish(evt, data) }

// BootstrapTLSAndDNS ensures runtime directories exist while TLS/DNS is handled elsewhere.
func (g *Globule) BootstrapTLSAndDNS(ctx context.Context) error {
	// Ensure dirs/config
	if err := g.InitFS(); err != nil {
		return err
	}

	g.log.Info("tls and dns bootstrap is handled by NodeAgent/controller; gateway makes no changes")
	return nil
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
