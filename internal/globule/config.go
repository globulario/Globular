package globule

import (
	"encoding/json"
	"log"

	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"

	"github.com/globulario/services/golang/config"
	Utility "github.com/globulario/utility"
	"github.com/txn2/txeh"
)

// setConfig applies a subset of runtime-reloadable settings and persists them.
// It restarts the process if a change requires it (protocol/domain).
func (g *Globule) SetConfig(m map[string]interface{}) error {
	needRestart := false

	// Domain
	if v, ok := m["Domain"].(string); ok && v != "" && v != g.Domain {
		g.Domain = v
		needRestart = true
	}

	// Protocol
	if v, ok := m["Protocol"].(string); ok && v != "" && v != g.Protocol {
		g.Protocol = v
		needRestart = true
	}

	// Ports
	if v, ok := asInt(m["PortHTTP"]); ok {
		g.PortHTTP = v
	}
	if v, ok := asInt(m["PortHTTPS"]); ok {
		g.PortHTTPS = v
	}

	// CORS
	if v, ok := asStrings(m["AllowedOrigins"]); ok {
		g.AllowedOrigins = v
	}
	if v, ok := asStrings(m["AllowedMethods"]); ok {
		g.AllowedMethods = v
	}
	if v, ok := asStrings(m["AllowedHeaders"]); ok {
		g.AllowedHeaders = v
	}

	// Session timeout (used by token refresher)
	if v, ok := asInt(m["SessionTimeout"]); ok && v > 0 {
		g.SessionTimeout = v
	}

	// Persist
	if err := g.SaveConfig(); err != nil {
		return fmt.Errorf("setConfig: saveConfig: %w", err)
	}

	if needRestart {
		return g.restart()
	}
	return nil
}

// GetConfig returns the current config as a map[string]interface{}.

/**
 * Return globular configuration.
 */
func (globule *Globule) GetConfig() map[string]interface{} {

	// TODO filter unwanted attributes...
	localConfig, _ := Utility.ToMap(globule)
	localConfig["Domain"] = globule.Domain
	localConfig["Name"] = globule.Name
	localConfig["OAuth2ClientID"] = globule.OAuth2ClientID

	services, _ := config.GetServicesConfigurations()

	// Get the array of service and set it back in the configurations.
	localConfig["Services"] = make(map[string]interface{})

	// Here I will set in a map and put in the Services key
	for i := range services {
		s := make(map[string]interface{})
		s["AllowAllOrigins"] = services[i]["AllowAllOrigins"]
		s["AllowedOrigins"] = services[i]["AllowedOrigins"]
		s["Description"] = services[i]["Description"]
		s["Discoveries"] = services[i]["Discoveries"]
		s["Domain"] = services[i]["Domain"]
		s["Address"] = services[i]["Address"]
		s["Id"] = services[i]["Id"]
		s["Keywords"] = services[i]["Keywords"]
		s["Name"] = services[i]["Name"]
		s["Mac"] = services[i]["Mac"]
		s["Port"] = services[i]["Port"]
		s["Proxy"] = services[i]["Proxy"]
		s["PublisherID"] = services[i]["PublisherID"]
		s["State"] = services[i]["State"]
		s["TLS"] = services[i]["TLS"]
		s["Dependencys"] = services[i]["Dependencys"]
		s["Version"] = services[i]["Version"]
		s["CertAuthorityTrust"] = services[i]["CertAuthorityTrust"]
		s["CertFile"] = services[i]["CertFile"]
		s["KeyFile"] = services[i]["KeyFile"]
		s["ConfigPath"] = services[i]["ConfigPath"]
		s["KeepAlive"] = services[i]["KeepAlive"]
		s["KeepUpToDate"] = services[i]["KeepUpToDate"]
		s["Pid"] = services[i]["Process"]

		if services[i]["Name"] == "file.FileService" {
			s["MaximumVideoConversionDelay"] = services[i]["MaximumVideoConversionDelay"]
			s["HasEnableGPU"] = services[i]["HasEnableGPU"]
			s["AutomaticStreamConversion"] = services[i]["AutomaticStreamConversion"]
			s["AutomaticVideoConversion"] = services[i]["AutomaticVideoConversion"]
			s["StartVideoConversionHour"] = services[i]["StartVideoConversionHour"]
		}

		// specific configuration values...
		if services[i]["Root"] != nil {
			s["Root"] = services[i]["Root"]
		}

		localConfig["Services"].(map[string]interface{})[s["Id"].(string)] = s
	}

	return localConfig
}

func asInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	}
	return 0, false
}

func asStrings(v interface{}) ([]string, bool) {
	a, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(a))
	for _, e := range a {
		if s, ok := e.(string); ok {
			out = append(out, s)
		}
	}
	return out, true
}

func (g *Globule) watchConfig() {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed:", err)
	}
	defer func() { _ = w.Close() }()

	cfg := g.configDir + "/config.json"
	_ = w.Add(cfg)

	go func() {
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if ev.Op == fsnotify.Write {
					b, _ := os.ReadFile(cfg)
					m := map[string]interface{}{}
					if err := json.Unmarshal(b, &m); err == nil {
						if err := g.SetConfig(m); err != nil {
							fmt.Println("invalid config on write:", err)
							os.Exit(1)
						}
					}
				}
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				fmt.Println("watch error:", err)
			}
		}
	}()
}

// saveConfig writes the current Globule configuration to disk and updates the hosts file.
// - Writes <configDir>/config.json with 0600 permissions
// - Adds local host mappings for Domain, Name.Domain, and AlternateDomains
func (g *Globule) SaveConfig() error {

	// Ensure config dir exists and remember it on the struct
	cfgDir := config.GetConfigDir()
	if err := Utility.CreateDirIfNotExist(cfgDir); err != nil {
		return fmt.Errorf("saveConfig: create config dir: %w", err)
	}

	// Persist the absolute executable path for reference
	if ex, err := os.Executable(); err == nil {
		if p, err := filepath.Abs(ex); err == nil {
			g.Path = strings.ReplaceAll(p, "\\", "/")
		}
	}

	// Marshal the struct (only exported fields get serialized; unexported server pointers are ignored)
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("saveConfig: marshal: %w", err)
	}

	cfgPath := filepath.Join(cfgDir, "config.json")
	if err := os.WriteFile(cfgPath, data, 0o600); err != nil {
		return fmt.Errorf("saveConfig: write %s: %w", cfgPath, err)
	}

	// Update hosts file entries for local resolution convenience
	if err := g.updateHostsFile(); err != nil {
		// Not fatal to saving config; log and continue
		fmt.Println("saveConfig: warning: update hosts file:", err)
	}

	fmt.Println("globular configuration saved at", cfgPath)
	return nil
}

// updateHostsFile ensures local DNS shortcuts are present in the system hosts file.
func (g *Globule) updateHostsFile() error {
	h, err := txeh.NewHostsDefault()
	if err != nil {
		return fmt.Errorf("hosts init: %w", err)
	}

	// Resolve local IP (best effort)
	localIP, _ := Utility.MyLocalIP(g.Mac)
	if localIP == "" {
		localIP = "127.0.0.1"
	}

	// Main domain entries
	if g.Domain != "" {
		h.AddHost(localIP, g.Domain)
		if g.Name != "" {
			h.AddHost(localIP, g.Name+"."+g.Domain)
		}
	}

	// Alternate domains (handle either []string or []interface{})
	for _, alt := range asStringSlice(g.AlternateDomains) {
		alt = strings.TrimSpace(strings.TrimPrefix(alt, "*."))
		if alt == "" {
			continue
		}
		h.AddHost(localIP, alt)
		if g.Name != "" {
			h.AddHost(localIP, g.Name+"."+alt)
		}
	}

	if err := h.Save(); err != nil {
		return fmt.Errorf("hosts save: %w", err)
	}
	return nil
}

// asStringSlice converts an interface{} slice (e.g., []interface{}) to []string safely.
func asStringSlice(v any) []string {
	switch t := v.(type) {
	case nil:
		return nil
	case []string:
		return t
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, e := range t {
			if s, ok := e.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}
