package main

import (
	controlplane "Globular/control-plane"
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	yaml "gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
	//"github.com/globulario/services/golang/authentication/authentication_client"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/persistence/persistence_client"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	service_manager_client "github.com/globulario/services/golang/services_manager/services_manager_client"
	Utility "github.com/globulario/utility"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gookit/color"
	"github.com/kardianos/service"

	//"github.com/slayer/autorestart"

	"github.com/txn2/txeh"

	// Interceptor for authentication, event, log...

	// Client services.

	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/go-acme/lego/v4/challenge/dns01"
)

// Global variable.
var (
	globule *Globule
)

// Globule represents the main server instance.
type Globule struct {
	// The share part of the service.
	Name string // the hostname.
	Mac  string // The Mac addresse

	// Where services can be found.
	ServicesRoot string

	// can be https or http.
	Protocol     string
	PortHTTP     int    // The port of the http file server.
	PortHTTPS    int    // The secure port
	PortsRange   string // The range of grpc ports.
	BackendPort  int    // This is backend resource port (mongodb port)
	BackendStore int

	// Cors policy
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

	Domain            string        // The principale domain
	AlternateDomains  []interface{} // Alternate domain for multiple domains
	IndexApplication  string        // If defined It will be use as the entry point where not application path was given in the url.
	LocalIpAddress    string        // The local ip address of the server.
	ExternalIpAddress string        // The public ip address of the server.

	// Certificate generation variables.
	CertExpirationDelay int
	CertPassword        string
	Country             string // tow letter.
	State               string // Full state name
	City                string
	Organization        string

	// https certificate info.
	Certificate                string
	CertificateAuthorityBundle string
	CertURL                    string
	CertStableURL              string

	// Keep the version number.
	Version  string
	Build    int
	Platform string

	// Admin informations.
	AdminEmail   string
	RootPassword string

	SessionTimeout int // The time before session expire.

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep globular up to date.

	// Update delay in second...
	WatchUpdateDelay int64

	// DNS stuff.
	DNS              string        // DNS.
	NS               []interface{} // Name server.
	DnsUpdateIpInfos []interface{} // The internet provader SetA info to keep ip up to date.

	// OAuth2 configuration.
	OAuth2_ClientId     string
	OAuth2_ClientSecret string
	OAuth2_RedirectUri  string

	// Reverse proxy will conain a list of addrees where to forward request and a route to forward request to.
	// ex : ["http://localhost:9100/metrics | /metric_01", "http://localhost:8080/metrics | /metric_02"]
	ReverseProxies []interface{}

	// Directories.
	Path string // The path of the exec...

	// Get it from
	webRoot string // The root of the http file server.
	data    string // the data directory
	creds   string // tls certificates
	config  string // configuration directory

	users        string // the users files directory
	applications string // The applications
	templates    string // the html/css templates
	projects     string // the web projects

	// ACME protocol registration
	registration *registration.Resource

	// exit channel.
	exit  chan bool
	exit_ bool

	// The http server
	http_server  *http.Server
	https_server *http.Server

	// List of peers
	peers *sync.Map // []*resourcepb.Peer

	// Use to save in configuration file, use peers in code...
	Peers []interface{}

	// Keep track of the strart time...
	startTime time.Time

	// This is use to display information to external service manager.
	logger service.Logger
}

// NewGlobule initializes and returns a new Globule instance with default configuration values.
// It sets up various properties such as version, build, platform, network settings, organization info,
// DNS and name server details, allowed CORS origins/methods/headers, peer map, protocol, certificate settings,
// admin credentials, session timeout, and update delay. It also registers HTTP handlers for various endpoints
// including configuration, certificate management, file uploads, media metadata, OAuth2, and server statistics.
// If no configuration file exists, it creates the necessary directories and saves the default configuration.
// Returns a pointer to the newly created Globule.
func NewGlobule() *Globule {

	// Here I will initialyse configuration.
	g := new(Globule)
	g.startTime = time.Now()
	g.exit_ = false
	g.exit = make(chan bool)
	g.Version = "1.0.0" // Automate version...
	g.Build = 0
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.IndexApplication = ""      // I will use the installer as defaut.
	g.PortHTTP = 80              // The default http port 80 is almost already use by other http server...
	g.PortHTTPS = 443            // The default https port number
	g.PortsRange = "10000-10100" // The default port range.
	g.ServicesRoot = config.GetServicesRoot()
	g.ExternalIpAddress = Utility.MyIP() // The public ip address of the server.

	// Set the default mac.
	g.Mac, _ = config.GetMacAddress()
	g.LocalIpAddress, _ = Utility.MyLocalIP(g.Mac) // The local ip address of the server.

	// THOSE values must be change by the user...
	g.Organization = "GLOBULARIO"
	g.Country = "CA"
	g.State = "QC"
	g.City = "MTL"

	// DNS info.
	g.DNS = globule.getLocalDomain() // The dns server.

	// The name server.
	g.NS = make([]interface{}, 0)

	if g.AllowedOrigins == nil {
		g.AllowedOrigins = []string{"*"}
	}

	if g.AllowedMethods == nil {
		g.AllowedMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}
	}

	if g.AllowedHeaders == nil {
		g.AllowedHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "domain", "application", "token", "video-path", "index-path", "routing"}
	}

	// the map of peers.
	g.peers = new(sync.Map)

	// Set the default checksum...
	g.Protocol = "http"
	g.Name, _ = config.GetName()

	// set default values.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"
	g.AdminEmail = "sa@globular.cloud"
	g.RootPassword = "adminadmin"

	// keep up to date by default.
	g.WatchUpdateDelay = 30 // seconds...
	g.SessionTimeout = 15   // in minutes

	// Keep in global var to by http handlers.
	globule = g

	// Set the list of http handler.

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// return the service descrition (.proto file content as json)
	http.HandleFunc("/get_service_descriptor", getServiceDescriptorHanldler)

	// return the service permissions
	http.HandleFunc("/get_service_permissions", getServicePermissionsHanldler)

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

	// The save configuration handler.
	http.HandleFunc("/save_config", saveConfigHanldler)

	// The checksum handler.
	http.HandleFunc("/checksum", getChecksumHanldler)

	// Handle the get ca certificate function
	http.HandleFunc("/get_ca_certificate", getCaCertificateHanldler)

	// Return info about the server
	http.HandleFunc("/stats", getHardwareData)

	// Return the plublic key
	http.HandleFunc("/public_key", getPublicKeyHanldler)

	// Return the certificate (signed by the let's encrypt CA)
	http.HandleFunc("/get_certificate", getCertificateHanldler)

	// Return the issuer certificate...
	http.HandleFunc("/get_issuer_certificate", getIssuerCertificateHandler)

	// Return the san server configuration.
	http.HandleFunc("/get_san_conf", getSanConfigurationHandler)

	// Handle the signing certificate function (sign by the local CA)
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

	// Return the list of images in the given directory.
	http.HandleFunc("/get_images", GetImagesHandler)

	// Create the video cover if it not already exist and return it as data url
	http.HandleFunc("/get_video_cover_data_url", GetCoverDataUrl)

	// Imdb movie api...
	http.HandleFunc("/imdb_titles", getImdbTitlesHanldler)
	http.HandleFunc("/imdb_title", getImdbTitleHanldler)

	// Get the file size at a given url.
	http.HandleFunc("/file_size", GetFileSizeAtUrl)

	// The OAuth2 handler, google login.
	http.HandleFunc("/oauth2callback", handleGoogleCallback)
	http.HandleFunc("/refresh_google_token", handleTokenRefresh)

	g.Path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	g.Path = strings.ReplaceAll(g.Path, "\\", "/")

	// If no configuration exist I will create it before initialyse directories and start services.
	configPath := config.GetConfigDir() + "/config.json"
	if !Utility.Exists(configPath) {
		err := Utility.CreateDirIfNotExist(config.GetConfigDir())
		if err != nil {
			fmt.Println("fail to create config directory with error", err)
		}

		globule.config = config.GetConfigDir()
		err = globule.saveConfig()
		if err != nil {
			fmt.Println("fail to save local configuration with error", err)
		}
	}

	return g
}

func (globule *Globule) cleanup() {

	p := make(map[string]interface{})

	p["domain"], _ = config.GetDomain()
	p["address"] = globule.getAddress()
	p["hostname"] = globule.Name
	p["mac"] = globule.Mac
	p["PortHTTP"] = globule.PortHTTP
	p["PortHTTPS"] = globule.PortHTTPS

	jsonStr, _ := json.Marshal(&p)

	// set services configuration values
	err := globule.publish("stop_peer_evt", jsonStr)
	if err != nil {
		fmt.Println("fail to publish stop_peer_evt", err)
	}

	// give time to stop peer evt to be publish
	time.Sleep(500 * time.Millisecond)

	// Close all services.
	err = globule.stopServices()
	if err != nil {
		fmt.Println("fail to stop services", err)
	}

	// reset firewall rules.
	err = resetRules()
	if err != nil {
		fmt.Println("fail to reset firewall rules", err)
	}

	err = globule.saveConfig()
	if err != nil {
		fmt.Println("fail to save config", err)
	}

	fmt.Println("bye bye!")
}

// Test if the domain has changed.
func domainHasChanged(domain string) bool {
	// if the domain has chaged that mean the sa@domain does not exist.
	return !Utility.Exists(config.GetDataDir() + "/files/users/sa@" + domain)
}

func (globule *Globule) registerAdminAccount() error {

	if len(globule.Domain) == 0 {
		return errors.New("domain is not set")
	}

	// this will return the first resource service with name resource.ResourceService
	resourceConfig, err := config.GetServiceConfigurationById("resource.ResourceService")
	if err != nil {
		return err
	}

	globule.BackendPort = Utility.ToInt(resourceConfig["Backend_port"])
	globule.BackendStore = 1 // default to SQLITE3

	// Set the backend store.
	if resourceConfig["Backend_type"].(string) == "MONGO" {
		globule.BackendStore = 0
	} else if resourceConfig["Backend_type"].(string) == "SQL" {
		globule.BackendStore = 1
	} else if resourceConfig["Backend_type"].(string) == "SCYLLA" {
		globule.BackendStore = 2
	}

	// get the resource client
	address, _ := config.GetAddress()
	resource_client_, err := getResourceClient(address)
	if err != nil {
		fmt.Println("fail to get resource client ", err)
		return err
	}

	// Create the admin account.
	results, _ := resource_client_.GetAccounts(`{"_id":"sa"}`)
	if len(results) == 0 {
		fmt.Println("fail to get admin account sa", err)
		fmt.Println("create admin account sa for domain ", globule.Domain)

		err := resource_client_.RegisterAccount(globule.Domain, "sa", "sa", globule.AdminEmail, globule.RootPassword, globule.RootPassword)
		if err != nil {
			return err
		}

		// Admin is created
		err = globule.createAdminRole()
		if err != nil {
			if !strings.Contains(err.Error(), "already exist") {
				return err
			}
		}

		path := config.GetDataDir() + "/files/users/sa@" + globule.Domain
		if !Utility.Exists(path) {

			// Set admin role to that account.
			err = resource_client_.AddAccountRole("sa", "admin")
			if err != nil {
				return err
			}

			err = Utility.CreateDirIfNotExist(path)
			if err == nil {
				err = globule.addResourceOwner("/users/sa@"+globule.Domain, "file", "sa@"+globule.Domain, rbacpb.SubjectType_ACCOUNT)
				if err != nil {
					fmt.Println("fail to add resource owner for sa@"+globule.Domain, err)
					return err
				}
			} else {
				fmt.Println("fail to create dir for sa@"+globule.Domain, err)
				return err
			}
		}

	} else {

		if domainHasChanged(globule.Domain) {

			// Alway update the sa domain...
			token, _ := security.GetLocalToken(globule.Mac)

			_, err := security.ValidateToken(token)
			if err != nil {
				fmt.Println("local token is not valid! ", err)
			}

			roles, err := resource_client_.GetRoles("")
			if err == nil {
				for i := 0; i < len(roles); i++ {
					if roles[i].Domain != globule.Domain {
						roles[i].Domain = globule.Domain
						err := resource_client_.UpdateRole(token, roles[i])
						if err != nil {
							fmt.Println("fail to update role with error: ", err)
						}
					}
				}
			}

			accounts, err := resource_client_.GetAccounts("")
			if err == nil {
				for i := 0; i < len(accounts); i++ {
					if accounts[i].Domain != globule.Domain {
						// I will update the account dir name
						err = os.Rename(config.GetDataDir()+"/files/users/"+accounts[i].Id+"@"+accounts[i].Domain, config.GetDataDir()+"/files/users/"+accounts[i].Id+"@"+globule.Domain)
						if err != nil {
							fmt.Println("fail to update account dir name with error: ", err)
							return err
						}

						// I will update the account domain
						accounts[i].Domain = globule.Domain
						err = resource_client_.SetAccount(token, accounts[i])
						if err != nil {
							fmt.Println("fail to update account with error: ", err)
						}
					}
				}
			}

			applications, err := resource_client_.GetApplications("")
			if err == nil {
				for i := 0; i < len(applications); i++ {
					if applications[i].Domain != globule.Domain {
						applications[i].Domain = globule.Domain
						err := resource_client_.UpdateApplication(token, applications[i])
						if err != nil {
							fmt.Println("fail to update application with error: ", err)
						}
					}
				}
			}

			groups, err := resource_client_.GetGroups("")
			if err == nil {
				for i := 0; i < len(groups); i++ {
					if groups[i].Domain != globule.Domain {
						groups[i].Domain = globule.Domain
						err := resource_client_.UpdateGroup(token, groups[i])
						if err != nil {
							fmt.Println("fail to update group with error: ", err)
						}
					}
				}
			}

			organisations, err := resource_client_.GetOrganizations("")
			if err == nil {
				for i := 0; i < len(organisations); i++ {
					if organisations[i].Domain != globule.Domain {
						organisations[i].Domain = globule.Domain
						err := resource_client_.UpdateOrganization(token, organisations[i])
						if err != nil {
							fmt.Println("fail to update organization with error: ", err)
						}
					}
				}
			}

		}
	}

	/* TODO create user connection*/

	// The user console
	return nil

}

// createAdminRole creates the "admin" role with all actions.
func (globule *Globule) createAdminRole() error {
	address, err := config.GetAddress()
	if err != nil {
		return err
	}

	resourceClient, err := getResourceClient(address)
	if err != nil {
		return err
	}

	// Normalize and validate MAC -> safe filename component.
	mac := strings.ReplaceAll(globule.Mac, ":", "_")
	if mac == "" || strings.Contains(mac, "..") {
		return fmt.Errorf("invalid mac address for token file: %q", mac)
	}
	for _, r := range mac {
		if !(r == '-' || r == '_' ||
			r >= '0' && r <= '9' ||
			r >= 'A' && r <= 'Z' ||
			r >= 'a' && r <= 'z') {
			return fmt.Errorf("invalid mac address for token file: %q", mac)
		}
	}

	// Constrain to <config>/tokens and resolve base (even if it's a symlink).
	base := filepath.Join(config.GetConfigDir(), "tokens")
	realBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return fmt.Errorf("resolve tokens dir: %w", err)
	}

	name := mac + "_token"
	if name != filepath.Base(name) {
		return fmt.Errorf("invalid token filename: %q", name)
	}

	tokenPath := filepath.Join(realBase, name)

	// Optional hardening: forbid the token file itself from being a symlink.
	if fi, err := os.Lstat(tokenPath); err == nil && (fi.Mode()&os.ModeSymlink) != 0 {
		return fmt.Errorf("token file is a symlink: %s", tokenPath)
	}

	// Read the token; path is validated and constrained to realBase.
	token, err := os.ReadFile(tokenPath) // #nosec G304 -- tokenPath validated & constrained
	if err != nil {
		return err
	}

	servicesManager, err := GetServiceManagerClient(address)
	if err != nil {
		return err
	}

	actions, err := servicesManager.GetAllActions()
	if err != nil {
		return err
	}

	// Create the admin role with all actions.
	if err := resourceClient.CreateRole(strings.TrimSpace(string(token)), "admin", "admin", actions); err != nil {
		fmt.Println("fail to create admin role:", err)
		return err
	}

	return nil
}

/**
 * Return globular configuration.
 */
func (globule *Globule) getConfig() map[string]interface{} {

	// TODO filter unwanted attributes...
	config_, _ := Utility.ToMap(globule)
	config_["Domain"] = globule.Domain
	config_["Name"] = globule.Name
	config_["OAuth2_ClientId"] = globule.OAuth2_ClientId

	services, _ := config.GetServicesConfigurations()

	// Get the array of service and set it back in the configurations.
	config_["Services"] = make(map[string]interface{})

	// Here I will set in a map and put in the Services key
	for i := 0; i < len(services); i++ {
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

		config_["Services"].(map[string]interface{})[s["Id"].(string)] = s
	}

	return config_
}

func (globule *Globule) savePeers() error {
	// Keep peers information here...
	globule.Peers = make([]interface{}, 0)
	globule.peers.Range(func(key, value interface{}) bool {
		p := value.(*resourcepb.Peer)
		port := p.PortHTTP
		if p.Protocol == "https" {
			port = p.PortHTTPS
		}
		globule.Peers = append(globule.Peers, map[string]interface{}{"Hostname": p.Hostname, "Domain": p.Domain, "Mac": p.Mac, "Port": port})
		return true
	})

	return globule.saveConfig()
}

/**
 * Set the Globule configuration with the given configuration, and the globule configuration.
 */
func (globule *Globule) setConfig(config map[string]interface{}) error {

	needRestart := false

	// Set the configuration.
	if config["AllowedOrigins"] != nil {
		globule.AllowedOrigins = make([]string, len(config["AllowedOrigins"].([]interface{})))
		for i := 0; i < len(config["AllowedOrigins"].([]interface{})); i++ {
			globule.AllowedOrigins[i] = config["AllowedOrigins"].([]interface{})[i].(string)
		}
	}

	// Set the allowed methods.
	if config["AllowedMethods"] != nil {
		globule.AllowedMethods = make([]string, len(config["AllowedMethods"].([]interface{})))
		for i := 0; i < len(config["AllowedMethods"].([]interface{})); i++ {
			globule.AllowedMethods[i] = config["AllowedMethods"].([]interface{})[i].(string)
		}
	}

	// Set the allowed headers.
	if config["AllowedHeaders"] != nil {
		globule.AllowedHeaders = make([]string, len(config["AllowedHeaders"].([]interface{})))
		for i := 0; i < len(config["AllowedHeaders"].([]interface{})); i++ {
			globule.AllowedHeaders[i] = config["AllowedHeaders"].([]interface{})[i].(string)
		}
	}

	// Set the domain.
	if config["Domain"] != nil {
		if len(config["Domain"].(string)) > 0 {
			// if the domain has changed I will need to restart the server.
			if globule.Domain != config["Domain"].(string) {
				needRestart = true
			}
			globule.Domain = config["Domain"].(string)
		}
	}

	// Set the AlternateDomains.
	if config["AlternateDomains"] != nil {
		globule.AlternateDomains = config["AlternateDomains"].([]interface{})
	}

	// Set the protocol.
	if config["Protocol"] != nil {
		if len(config["Protocol"].(string)) > 0 {
			// if the protocol has changed I will need to restart the server.
			if globule.Protocol != config["Protocol"].(string) {
				needRestart = true
				globule.Protocol = config["Protocol"].(string)
			}
		}
	}

	// Set the port.
	if config["PortHTTP"] != nil {
		globule.PortHTTP = Utility.ToInt(config["PortHTTP"])
	}

	// Set the port.
	if config["PortHTTPS"] != nil {
		globule.PortHTTPS = Utility.ToInt(config["PortHTTPS"])
	}

	// Set the ports range.
	if config["PortsRange"] != nil {
		if len(config["PortsRange"].(string)) > 0 {
			globule.PortsRange = config["PortsRange"].(string)
		}
	}

	// Set the backend port.
	if config["BackendPort"] != nil {
		globule.BackendPort = Utility.ToInt(config["BackendPort"])
	}

	// Set the backend store.
	if config["BackendStore"] != nil {

		globule.BackendStore = Utility.ToInt(config["BackendStore"])

	}

	// Set the certificate expiration delay.
	if config["CertExpirationDelay"] != nil {

		globule.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"])

	}

	// Set the certificate password.
	if config["CertPassword"] != nil {
		if len(config["CertPassword"].(string)) > 0 {
			globule.CertPassword = config["CertPassword"].(string)
		}
	}

	// Set the country.
	if config["Country"] != nil {
		if len(config["Country"].(string)) > 0 {
			globule.Country = config["Country"].(string)
		}
	}

	// Set the state.
	if config["State"] != nil {
		if len(config["State"].(string)) > 0 {
			globule.State = config["State"].(string)
		}
	}

	// Set the city.
	if config["City"] != nil {
		if len(config["City"].(string)) > 0 {
			globule.City = config["City"].(string)
		}
	}

	// Set the organization.
	if config["Organization"] != nil {
		if len(config["Organization"].(string)) > 0 {
			globule.Organization = config["Organization"].(string)
		}
	}

	// Set the admin email.
	if config["AdminEmail"] != nil {
		if len(config["AdminEmail"].(string)) > 0 {
			globule.AdminEmail = config["AdminEmail"].(string)
		}
	}

	// Set the root password.
	if config["RootPassword"] != nil {
		if len(config["RootPassword"].(string)) > 0 {
			globule.RootPassword = config["RootPassword"].(string)
		}
	}

	// Set the session timeout.
	if config["SessionTimeout"] != nil {
		globule.SessionTimeout = Utility.ToInt(config["SessionTimeout"])
	}

	// Set the watch update delay.
	if config["WatchUpdateDelay"] != nil {
		globule.WatchUpdateDelay = int64(Utility.ToInt(config["WatchUpdateDelay"]))
	}

	// Set the ns.
	if config["NS"] != nil {
		globule.NS = config["NS"].([]interface{})
	}

	// Set the dns.
	if config["DNS"] != nil {
		if len(config["DNS"].(string)) > 0 {
			if globule.DNS != config["DNS"].(string) {

				// set the dns
				globule.DNS = config["DNS"].(string)

				// register the ip to dns.
				err := globule.registerIpToDns()
				if err != nil {
					fmt.Println("fail to register ip to dns with error: ", err)
				}

			}
		}
	}

	// Set the dns update ip infos.
	if config["DnsUpdateIpInfos"] != nil {
		if len(config["DnsUpdateIpInfos"].(string)) > 0 {
			globule.DnsUpdateIpInfos = config["DnsUpdateIpInfos"].([]interface{})
		}
	}

	// Set the reverse proxies.
	if config["ReverseProxies"] != nil {
		if len(config["ReverseProxies"].(string)) > 0 {
			globule.ReverseProxies = config["ReverseProxies"].([]interface{})
		}
	}

	// Set the discoveries.
	if config["Discoveries"] != nil {
		globule.Discoveries = make([]string, len(config["Discoveries"].([]interface{})))
		for i := 0; i < len(config["Discoveries"].([]interface{})); i++ {
			globule.Discoveries[i] = config["Discoveries"].([]interface{})[i].(string)
		}
	}

	// Set the peers.
	if config["Peers"] != nil {
		globule.Peers = make([]interface{}, len(config["Peers"].([]interface{})))
		for i := 0; i < len(config["Peers"].([]interface{})); i++ {
			globule.Peers[i] = config["Peers"].([]interface{})[i].(map[string]interface{})
		}
	}

	// Set the index application.
	if config["IndexApplication"] != nil {
		if len(config["IndexApplication"].(string)) > 0 {
			globule.IndexApplication = config["IndexApplication"].(string)
		}
	}

	if needRestart {
		// I will stop the services...
		err := globule.restart()
		if err != nil {
			fmt.Println("fail to restart globule with error: ", err)
			return err
		}
	}

	return nil
}

/**
 * Save the configuration
 */
func (globule *Globule) saveConfig() error {

	// set the path
	globule.Path, _ = os.Executable()

	jsonStr, err := Utility.ToJson(globule)
	if err != nil {
		fmt.Println("fail to save configuration with error: ", err)
		return err
	}

	err = Utility.CreateDirIfNotExist(globule.config)
	if err != nil {
		fmt.Println("fail to create config directory with error: ", err)
	}

	configPath := globule.config + "/config.json"

	err = os.WriteFile(configPath, []byte(jsonStr), 0600)
	if err != nil {
		fmt.Println("fail to save configuration with error: ", err)
		return err
	}

	// Here I will set the hosts file with the domain and alternate domain.
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		fmt.Println("fail to get hosts file with error: ", err)
	}

	// Set The hosts file with the domain and alternate domain.
	localIp, _ := Utility.MyLocalIP(globule.Mac)
	hosts.AddHost(localIp, globule.Domain)
	hosts.AddHost(localIp, globule.Name+"."+globule.Domain)

	// Set the alternate domain.
	for i := 0; i < len(globule.AlternateDomains); i++ {
		alternateDomain := strings.TrimPrefix(globule.AlternateDomains[i].(string), "*.") // remove the * if exist
		hosts.AddHost(localIp, alternateDomain)
		hosts.AddHost(localIp, globule.Name+"."+alternateDomain)
	}

	// Save the hosts file.
	err = hosts.Save()
	if err != nil {
		fmt.Println("fail to save hosts file with error: ", err)
		return err
	}

	fmt.Println("globular configuration was save at ", configPath)

	return nil
}

// GetEmail returns the administrator email address associated with the Globule instance.
func (globule *Globule) GetEmail() string {
	return globule.AdminEmail
}

// GetRegistration returns the registration resource associated with the Globule instance.
func (globule *Globule) GetRegistration() *registration.Resource {
	return globule.registration
}

// GetPrivateKey returns the private key associated with the Globule instance.
func (globule *Globule) GetPrivateKey() crypto.PrivateKey {
	keyPem, err := os.ReadFile(globule.creds + "/client.pem")
	if err != nil {
		return nil
	}

	keyBlock, _ := pem.Decode(keyPem)
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil
	}
	return privateKey
}

// DNSProviderGlobularDNS implements the Let's Encrypt DNS challenge.
type DNSProviderGlobularDNS struct {
	apiAuthToken string
}

// NewDNSProviderGlobularDNS creates a new instance of DNSProviderGlobularDNS using the provided API authentication token.
// It returns a pointer to the DNSProviderGlobularDNS and an error if the creation fails.
//
// Parameters:
//
//	apiAuthToken - The API authentication token required to interact with the Globular DNS provider.
//
// Returns:
//
//	*DNSProviderGlobularDNS - A pointer to the newly created DNSProviderGlobularDNS instance.
//	error - An error value if the creation fails, otherwise nil.
func NewDNSProviderGlobularDNS(apiAuthToken string) (*DNSProviderGlobularDNS, error) {
	return &DNSProviderGlobularDNS{apiAuthToken: apiAuthToken}, nil
}

// Present creates a DNS TXT record to fulfill the ACME DNS-01 challenge for the specified domain.
// It connects to the Globular DNS service, generates an authentication token, and sets the required
// key-value pair in the DNS server. If any step fails, an error is returned.
//
// Parameters:
//
//	domain  - The domain for which the DNS challenge is being performed.
//	token   - The challenge token provided by the ACME server.
//	keyAuth - The key authorization string.
//
// Returns:
//
//	error - An error if the DNS record could not be set, or nil on success.
func (d *DNSProviderGlobularDNS) Present(domain, token, keyAuth string) error {
	key, value := dns01.GetRecord(domain, keyAuth)

	if len(globule.DNS) > 0 {
		fmt.Println("Let's encrypt dns challenge...")
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS, "dns.DnsService")
		if err != nil {
			return err
		}

		// generate a token for the dns service.
		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail, globule.Domain)

		if err != nil {
			fmt.Println("fail to connect with the dns server")
			return err
		}

		// set the key value pair in the dns server.
		err = dns_client_.SetText(token, key, []string{value}, 30)

		if err != nil {
			fmt.Println("fail to set text with error ", err)
			return err
		}
	}

	return nil
}

// CleanUp removes any DNS state created during the ACME DNS-01 challenge,
// specifically the TXT record associated with the provided domain and keyAuth.
// It connects to the DNS service, generates an authentication token, and
// attempts to remove the TXT record. Returns an error if any step fails.
func (d *DNSProviderGlobularDNS) CleanUp(domain, token, keyAuth string) error {
	// clean up any state you created in Present, like removing the TXT record
	key, _ := dns01.GetRecord(domain, keyAuth)

	if len(globule.DNS) > 0 {
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS, "dns.DnsService")
		if err != nil {
			return err
		}

		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail, globule.Domain)

		if err != nil {

			return err
		}

		err = dns_client_.RemoveText(token, key)
		if err != nil {
			fmt.Println("fail to remove challenge key with error ", err)
			return err
		}

	}
	return nil
}

/**
 * That function work correctly, but the DNS fail time to time to give the
 * IP address that result in a fail request... The DNS must be fix!
 */
func (globule *Globule) obtainCertificateForCsr() error {

	// I will use the dns challenge to get the certificate.
	config_ := lego.NewConfig(globule)
	config_.Certificate.KeyType = certcrypto.RSA2048
	client, err := lego.NewClient(config_)
	if err != nil {
		fmt.Println("fail to create new lego client with error: ", err)
		return err
	}

	// Dns registration will be use in case dns service are available.
	// TODO dns challenge give JWS has invalid anti-replay nonce error... at the moment
	// http chanllenge do the job but wildcald domain name are not allowed...
	if len(globule.DNS) > 0 {

		// Get the local token.
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS, "dns.DnsService")
		if err != nil {
			fmt.Println("fail to create new Dns client with error: ", err)
			return err
		}

		defer dns_client_.Close()

		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail, globule.Domain)
		if err != nil {
			fmt.Println("fail to generate token with error: ", err)
			return err
		}

		globularDNS, err := NewDNSProviderGlobularDNS(token)
		if err != nil {
			fmt.Println("fail to create new dns provider with error: ", err)
			return err
		}

		fmt.Println("use dns challenge")
		err = client.Challenge.SetDNS01Provider(globularDNS)
		if err != nil {
			fmt.Println("fail to set dns provider with error: ", err)
			return err
		}

	} else {
		provider := http01.NewProviderServer("", strconv.Itoa(globule.PortHTTP))
		err = client.Challenge.SetHTTP01Provider(provider)
		if err != nil {
			fmt.Println("fail to set http provider with error: ", err)
			return err
		}
	}

	if err != nil {
		fmt.Println("fail to create new client with error: ", err)
		return err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	globule.registration = reg
	if err != nil {
		fmt.Println("fail to register with error: ", err)
		return err
	}

	csrPem, err := os.ReadFile(globule.creds + "/server.csr")
	if err != nil {
		fmt.Println("fail to read certificate request with error: ", err)
		return err
	}

	csrBlock, _ := pem.Decode(csrPem)
	rqstForCsr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		fmt.Println("fail to parse certificate request with error:  ", err)
		return err
	}

	resource, err := client.Certificate.ObtainForCSR(*rqstForCsr, true)
	if err != nil {
		fmt.Println("fail to obtain certificate with error: ", err)
		return err
	}

	// Keep certificates url in the config.
	globule.CertURL = resource.CertURL
	globule.CertStableURL = resource.CertStableURL

	// Set the certificates paths...
	globule.Certificate = globule.Domain + ".crt"
	globule.CertificateAuthorityBundle = globule.Domain + ".issuer.crt"

	// Save the certificate in the cerst folder.
	err = os.WriteFile(globule.creds+"/"+globule.Certificate, resource.Certificate, 0400)
	if err != nil {
		fmt.Println("fail to save certificate with error: ", err)
		return err
	}
	err = os.WriteFile(globule.creds+"/"+globule.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)
	if err != nil {
		fmt.Println("fail to save certificate authority bundle with error: ", err)
		return err
	}

	// save the config with the values.
	return globule.saveConfig()
}

func (globule *Globule) signCertificate(clientCSR string) (string, error) {
	// --- Parse CSR ---
	csrBlock, _ := pem.Decode([]byte(clientCSR))
	if csrBlock == nil {
		return "", errors.New("invalid CSR: not PEM")
	}
	if csrBlock.Type != "CERTIFICATE REQUEST" && csrBlock.Type != "NEW CERTIFICATE REQUEST" {
		return "", fmt.Errorf("invalid CSR: unexpected PEM type %q", csrBlock.Type)
	}
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse CSR: %w", err)
	}
	if err := csr.CheckSignature(); err != nil {
		return "", fmt.Errorf("CSR signature invalid: %w", err)
	}

	// --- Load CA cert ---
	caCrtPEM, err := os.ReadFile(filepath.Join(globule.creds, "ca.crt"))
	if err != nil {
		return "", fmt.Errorf("read ca.crt: %w", err)
	}
	var caCert *x509.Certificate
	rest := caCrtPEM
	for {
		var b *pem.Block
		b, rest = pem.Decode(rest)
		if b == nil {
			break
		}
		if b.Type == "CERTIFICATE" {
			caCert, err = x509.ParseCertificate(b.Bytes)
			if err != nil {
				return "", fmt.Errorf("parse ca.crt: %w", err)
			}
			break
		}
	}
	if caCert == nil {
		return "", errors.New("ca.crt: no CERTIFICATE block found")
	}

	// --- Load CA key (supports unencrypted, legacy-encrypted, and PKCS#8 unencrypted) ---
	caKeyPEM, err := os.ReadFile(filepath.Join(globule.creds, "ca.key"))
	if err != nil {
		return "", fmt.Errorf("read ca.key: %w", err)
	}
	keyBlock, _ := pem.Decode(caKeyPEM)
	if keyBlock == nil {
		return "", errors.New("ca.key: invalid PEM")
	}

	var keyDER []byte
	if x509.IsEncryptedPEMBlock(keyBlock) {
		if globule.CertPassword == "" {
			return "", errors.New("ca.key is encrypted but CertPassword is empty")
		}
		keyDER, err = x509.DecryptPEMBlock(keyBlock, []byte(globule.CertPassword))
		if err != nil {
			return "", fmt.Errorf("decrypt ca.key: %w", err)
		}
	} else {
		keyDER = keyBlock.Bytes
	}

	var signer crypto.Signer
	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		k, err := x509.ParsePKCS1PrivateKey(keyDER)
		if err != nil {
			return "", fmt.Errorf("parse RSA key: %w", err)
		}
		signer = k
	case "EC PRIVATE KEY":
		k, err := x509.ParseECPrivateKey(keyDER)
		if err != nil {
			return "", fmt.Errorf("parse EC key: %w", err)
		}
		signer = k
	case "PRIVATE KEY": // PKCS#8 (unencrypted)
		kAny, err := x509.ParsePKCS8PrivateKey(keyDER)
		if err != nil {
			return "", fmt.Errorf("parse PKCS#8 key: %w", err)
		}
		switch k := kAny.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			signer = k.(crypto.Signer)
		default:
			return "", fmt.Errorf("unsupported PKCS#8 key type %T", k)
		}
	default:
		// If your ca.key is "ENCRYPTED PRIVATE KEY" (PKCS#8 encrypted), convert it to legacy-encrypted or unencrypted,
		// or use a pkcs8 decrypter library.
		return "", fmt.Errorf("unsupported key PEM type %q", keyBlock.Type)
	}

	// --- Build certificate template ---
	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return "", fmt.Errorf("serial: %w", err)
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.Add(time.Duration(globule.CertExpirationDelay) * 24 * time.Hour)

	tpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               csr.Subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		// Copy SANs & emails from the CSR
		DNSNames:       csr.DNSNames,
		IPAddresses:    csr.IPAddresses,
		URIs:           csr.URIs,
		EmailAddresses: csr.EmailAddresses,
		// Helps chain building if CA has a SubjectKeyId
		AuthorityKeyId: caCert.SubjectKeyId,
	}
	// Respect requested CSR extensions (e.g., custom OIDs)
	tpl.ExtraExtensions = append(tpl.ExtraExtensions, csr.Extensions...)

	// --- Sign ---
	der, err := x509.CreateCertificate(rand.Reader, tpl, caCert, csr.PublicKey, signer)
	if err != nil {
		return "", fmt.Errorf("create certificate: %w", err)
	}

	var pemOut bytes.Buffer
	if err := pem.Encode(&pemOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return "", fmt.Errorf("PEM encode: %w", err)
	}
	return pemOut.String(), nil
}

/**
 * Initialize the server directories config, data, webroot...
 */
func (globule *Globule) initDirectories() error {

	fmt.Println("init directories")

	// initilayse configurations...
	// it must be call here in order to initialyse a sync map...
	_, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	// The dns update ip info.
	// for example:
	// {
	//	"Key": "your key generated by your domain provider",
	//	"Secret": "the secret generated by your domain name provider",
	//	"SetA": "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
	// }
	globule.DnsUpdateIpInfos = make([]interface{}, 0)

	// Set the list of discorvery service avalaible...
	globule.Discoveries = make([]string, 0)

	//////////////////////////////////////////////////////////////////////////////////////
	// There is the default directory initialisation...
	//////////////////////////////////////////////////////////////////////////////////////

	// Create the directory if is not exist.
	globule.data = config.GetDataDir()

	err = Utility.CreateDirIfNotExist(globule.data)
	if err != nil {
		return err
	}

	globule.webRoot = config.GetWebRootDir()

	err = Utility.CreateDirIfNotExist(globule.webRoot)
	if err != nil {
		return err
	}

	globule.templates = globule.data + "/files/templates"

	err = Utility.CreateDirIfNotExist(globule.templates)
	if err != nil {
		return err
	}

	globule.projects = globule.data + "/files/projects"
	err = Utility.CreateDirIfNotExist(globule.projects)
	if err != nil {
		return err
	}

	globule.config = config.GetConfigDir()
	err = Utility.CreateDirIfNotExist(globule.config)

	if err != nil {
		fmt.Println("fail to create configuration directory  with error", err)
		return err
	}

	// Create the tokens directory
	err = Utility.CreateDirIfNotExist(globule.config + "/tokens")
	if err != nil {
		fmt.Println("fail to create tokens directory  with error", err)
		return err
	}

	// Create the creds directory if it not already exist.
	globule.creds = globule.config + "/tls/" + globule.getLocalDomain()
	err = Utility.CreateDirIfNotExist(globule.creds)
	if err != nil {
		fmt.Println("fail to create creds directory with error", err)
		return err
	}

	// Files directorie that contain user's directories and application's directory
	globule.users = globule.data + "/files/users"
	err = Utility.CreateDirIfNotExist(globule.users)
	if err != nil {
		fmt.Println("fail to create users directory with error", err)
		return err
	}

	// Contain the application directory.
	globule.applications = globule.data + "/files/applications"
	err = Utility.CreateDirIfNotExist(globule.applications)
	if err != nil {
		fmt.Println("fail to create applications directory with error", err)
		return err
	}

	// Initialyse globular from it configuration file.
	file, err := os.ReadFile(globule.config + "/config.json")

	// Init the service with the default port address
	if err == nil {

		err := json.Unmarshal(file, &globule)
		if err != nil {
			fmt.Println("fail to init configuation with error ", err)
			return err
		}

	} else {
		jsonStr, err := Utility.ToJson(&globule)
		if err == nil {
			err := os.WriteFile(globule.config+"/config.json", []byte(jsonStr), 0600)
			if err != nil {
				return err
			}
		}
	}

	// I will put the domain into the
	if globule.AlternateDomains == nil && len(globule.Domain) > 0 && globule.Domain != "localhost" {
		globule.AlternateDomains = make([]interface{}, 0)
	}

	// Set the default domain.
	if len(globule.Domain) == 0 {
		globule.Domain = "localhost"
	}

	if len(globule.Mac) == 0 {
		globule.Mac, err = config.GetMacAddress()
		if err != nil {
			return err
		}
	}

	// save config...
	err = globule.saveConfig()
	if err != nil {
		return err
	}

	if !Utility.Exists(globule.webRoot + "/index.html") {

		// in that case I will create a new index.html file.
		err = os.WriteFile(globule.webRoot+"/index.html", []byte(
			`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
		<html lang="en">
		
			<head>
				<meta http-equiv="content-type" content="text/html; charset=utf-8">
				<title>Title Goes Here</title>
			</head>
		
			<body>
				<p>Welcome to Globular `+globule.Version+`</p>
			</body>
		
		</html>`), 0600)

		if err != nil {
			fmt.Println("fail to create index.html with error", err)
			return err
		}
	}

	return nil
}

func (globule *Globule) refreshLocalToken() error {

	// set the local token.
	return security.SetLocalToken(globule.Mac, globule.Domain, "sa", "sa", globule.AdminEmail, globule.SessionTimeout)
}
func enablePorts(ruleName, portsRange string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	if err := validateRuleName(ruleName); err != nil {
		return fmt.Errorf("invalid rule name: %w", err)
	}
	if err := validatePortsRange(portsRange); err != nil {
		return fmt.Errorf("invalid ports range: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete any existing rules with that name (in/out)
	// #nosec G204 -- see justification above; name validated in addFirewallRule
	_ = exec.CommandContext(ctx,
		"netsh", "advfirewall", "firewall", "delete", "rule", "name="+ruleName,
	).Run() // ignore error if rule didn't exist

	// Inbound allow
	argsIn := []string{
		"advfirewall", "firewall", "add", "rule",
		"name=" + ruleName,
		"dir=in",
		"action=allow",
		"protocol=TCP",
		"localport=" + portsRange,
	}

	// #nosec G204 -- Executable is a constant ("netsh"); inputs are validated (name, portsRange).
	if out, err := exec.CommandContext(ctx, "netsh", argsIn...).CombinedOutput(); err != nil {
		return fmt.Errorf("netsh inbound failed: %v, output: %s", err, strings.TrimSpace(string(out)))
	}

	// Outbound allow
	argsOut := []string{
		"advfirewall", "firewall", "add", "rule",
		"name=" + ruleName,
		"dir=out",
		"action=allow",
		"protocol=TCP",
		"localport=" + portsRange,
	}
	// #nosec G204 -- Executable is a constant ("netsh"); inputs are validated (name, portsRange).
	if out, err := exec.CommandContext(ctx, "netsh", argsOut...).CombinedOutput(); err != nil {
		return fmt.Errorf("netsh outbound failed: %v, output: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

var ruleNameRX = regexp.MustCompile(`^[A-Za-z0-9 _.-]{1,64}$`)

// Validate a friendly name (no quotes, no shell metachars; keep it simple & short)
func validateRuleName(s string) error {
	if !ruleNameRX.MatchString(s) {
		return errors.New(`rule name must match [A-Za-z0-9 _.-], max 64 chars`)
	}
	return nil
}

// Accept forms like: "80", "80-90", "80,443,10000-10100"
func validatePortsRange(s string) error {
	if s == "" {
		return errors.New("empty")
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return errors.New("empty port segment")
		}
		if strings.Contains(p, "-") {
			b := strings.SplitN(p, "-", 2)
			if len(b) != 2 {
				return fmt.Errorf("bad range: %q", p)
			}
			a1, err1 := parsePort(b[0])
			a2, err2 := parsePort(b[1])
			if err1 != nil || err2 != nil || a1 > a2 {
				return fmt.Errorf("bad range: %q", p)
			}
		} else {
			if _, err := parsePort(p); err != nil {
				return err
			}
		}
	}
	return nil
}

func parsePort(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 65535 {
		return 0, fmt.Errorf("invalid port: %q", s)
	}
	return n, nil
}

// enableProgramFwMgr adds inbound and outbound firewall rules for the specified program on Windows systems.
// It uses the Windows 'netsh advfirewall' command to allow network traffic for the given application executable.
// The function takes the rule name and the application path as arguments.
// Returns an error if the firewall rule addition fails, or nil on success.
// On non-Windows systems, the function does nothing and returns nil.
func enableProgramFwMgr(name, appname string) error {

	if runtime.GOOS == "windows" {
		fmt.Println("enable program: ", name, appname)
		// netsh advfirewall firewall add rule name="MongoDB Database Server" dir=in action=allow program="C:\Program Files\Globular\Dependencys\mongodb-win32-x86_64-windows-5.0.5\bin\mongod.exe" enable=yes
		appname = strings.ReplaceAll(appname, "/", "\\")
		// #nosec G204 -- see justification above; name validated in addFirewallRule
		inboundRule := exec.Command("cmd", "/C", fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=in action=allow program="%s" enable=yes`, name, appname))
		inboundRule.Dir = os.TempDir()
		err := inboundRule.Run()
		if err != nil {
			return err
		}
		// #nosec G204 -- see justification above; name validated in addFirewallRule
		outboundRule := exec.Command("cmd", "/C", fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow program="%s" enable=yes`, name, appname))
		outboundRule.Dir = os.TempDir()
		err = outboundRule.Run()
		if err != nil {
			return err
		}

	}
	return nil
}

func deleteRule(name string) error {
	if runtime.GOOS == "windows" {
		// netsh advfirewall firewall delete rule name= rule "Globular-Services"
		// #nosec G204 -- see justification above; name validated in addFirewallRule
		cmd := exec.Command("cmd", "/C", fmt.Sprintf(`netsh advfirewall firewall delete rule name="%s"`, name))
		cmd.Dir = os.TempDir()
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return nil
}

func resetRules() error {

	services, err := config.GetOrderedServicesConfigurations()
	if err != nil {
		return err
	}

	// set rules for services contain in bin folder.
	err = deleteRule("alertmanager")
	if err != nil {
		fmt.Println("fail to delete rule: ", "alertmanager", " with error: ", err)
	}

	err = deleteRule("mongod")
	if err != nil {
		fmt.Println("fail to delete rule: ", "mongod", " with error: ", err)
	}

	err = deleteRule("prometheus")
	if err != nil {
		fmt.Println("fail to delete rule: ", "prometheus", " with error: ", err)
	}

	err = deleteRule("torrent")
	if err != nil {
		fmt.Println("fail to delete rule: ", "torrent", " with error: ", err)
	}

	err = deleteRule("yt-dlp")
	if err != nil {
		fmt.Println("fail to delete rule: ", "yt-dlp", " with error: ", err)
	}

	// other rules.
	err = deleteRule("Globular")
	if err != nil {
		fmt.Println("fail to delete rule: ", "Globular", " with error: ", err)
	}

	err = deleteRule("Globular-http")
	if err != nil {
		fmt.Println("fail to delete rule: ", "Globular-http", " with error: ", err)
	}

	err = deleteRule("Globular-https")
	if err != nil {
		fmt.Println("fail to delete rule: ", "Globular-https", " with error: ", err)
	}

	err = deleteRule("Globular-Services")
	if err != nil {
		fmt.Println("fail to delete rule: ", "Globular-Services", " with error: ", err)
	}

	for i := range services {
		// Create the service process.
		err = deleteRule(services[i]["Name"].(string) + "-" + services[i]["Id"].(string))
		if err != nil {
			fmt.Println("fail to delete rule: ", services[i]["Name"].(string), "-", services[i]["Id"].(string), " with error: ", err)
		}
	}

	return nil
}

func resetSystemPath() error {

	if runtime.GOOS == "windows" {

		err := Utility.UnsetWindowsEnvironmentVariable("OPENSSL_CONF")
		if err != nil {
			fmt.Println("fail to unset environment variable OPENSSL_CONF with error: ", err)
		}

		systemPath, err := Utility.GetWindowsEnvironmentVariable("Path")
		if err != nil {
			fmt.Println("fail to get environment variable Path with error: ", err)
			return err
		}

		// convert to \ to /
		systemPath = strings.ReplaceAll(systemPath, "\\", "/")

		// needed to retreive Globular.exe
		if strings.Contains(systemPath, config.GetRootDir()) {
			systemPath = strings.Replace(systemPath, ";"+config.GetRootDir(), "", 1)
		}

		// needed to retreive various command...
		if strings.Contains(systemPath, config.GetRootDir()+"/bin") {
			systemPath = strings.Replace(systemPath, ";"+config.GetRootDir()+"/bin", "", 1)
		}

		// remove path...
		execs := Utility.GetFilePathsByExtension(config.GetRootDir()+"/Dependencys", ".exe")
		for i := range execs {
			exec := strings.ReplaceAll(execs[i], "\\", "/")
			exec = exec[:strings.LastIndex(exec, "/")]
			if strings.Contains(systemPath, exec) {
				systemPath = strings.Replace(systemPath, ";"+exec, "", 1)
			}
		}

		// set system path...
		err = Utility.SetWindowsEnvironmentVariable("Path", strings.ReplaceAll(systemPath, "/", "\\"))
		if err != nil {
			return err
		}

	}

	return nil
}

// Set all required path.
func setSystemPath() error {
	// so here I will append
	switch runtime.GOOS {
	case "windows":
		// remove previous rules...
		err := resetRules()
		if err != nil {
			return err
		}

		ex, err := os.Executable()
		if err != nil {
			return err
		}

		// set globular firewall run...
		err = enableProgramFwMgr("Globular", ex)
		if err != nil {
			fmt.Println("fail to set rule for Globular with error: ", err)
		}

		// Enable ports
		err = enablePorts("Globular-Services", globule.PortsRange)
		if err != nil {
			fmt.Println("fail to set rule for Globular-Services with error: ", err)
		}

		err = enablePorts("Globular-http", strconv.Itoa(globule.PortHTTP))
		if err != nil {
			fmt.Println("fail to set rule for Globular-http with error: ", err)
		}

		err = enablePorts("Globular-https", strconv.Itoa(globule.PortHTTPS))
		if err != nil {
			fmt.Println("fail to set rule for Globular-https with error: ", err)
		}

		systemPath, err := Utility.GetWindowsEnvironmentVariable("Path")

		if err != nil {
			fmt.Println("fail to get environnement %Path with error", err)
			return err
		}

		// convert to \ to /
		systemPath = strings.ReplaceAll(systemPath, "\\", "/")

		// needed to retreive Globular.exe
		if !strings.Contains(systemPath, config.GetRootDir()) {
			systemPath += ";" + config.GetRootDir()
		}

		// needed to retreive various command...
		if !strings.Contains(systemPath, config.GetRootDir()+"/bin") {
			systemPath += ";" + config.GetRootDir() + "/bin"
		}

		// set rules for services contain in Dependencys folder.
		execs := Utility.GetFilePathsByExtension(config.GetRootDir()+"/Dependencys", ".exe")
		for i := range execs {
			exec := strings.ReplaceAll(execs[i], "\\", "/")

			if strings.HasSuffix(exec, "prometheus.exe") {
				err := enableProgramFwMgr("prometheus", exec)
				if err != nil {
					fmt.Println("fail to set rule for prometheus.exe with error", err)
				}
			}

			if strings.HasSuffix(exec, "mongod.exe") {
				err := enableProgramFwMgr("mongo", exec)
				if err != nil {
					fmt.Println("fail to set rule for mongod.exe with error", err)
				}
			}

			if strings.HasSuffix(exec, "alertmanager.exe") {
				err := enableProgramFwMgr("alertmanager", exec)
				if err != nil {
					fmt.Println("fail to set rule for alertmanager.exe with error", err)
				}
			}

			if strings.HasSuffix(exec, "torrent.exe") {
				err := enableProgramFwMgr("torrent", exec)
				if err != nil {
					fmt.Println("fail to set rule for torrent.exe with error", err)
				}
			}

			if strings.HasSuffix(exec, "yt-dlp.exe") {
				err := enableProgramFwMgr("yt-dlp", exec)
				if err != nil {
					fmt.Println("fail to set rule for yt-dlp.exe with error", err)
				}
			}

			exec = exec[:strings.LastIndex(exec, "/")]
			if !strings.Contains(systemPath, exec) {
				systemPath += ";" + exec
			}
		}

		// now the services rules
		services, err := config.GetServicesConfigurations()
		if err != nil {
			return err
		}

		for i := range services {
			service := services[i]
			id := service["Id"].(string)
			path := service["Path"].(string)
			name := service["Name"].(string)

			// Create the service process.
			err := enableProgramFwMgr(name+"-"+id, path)
			if err != nil {
				fmt.Println("fail to set rule for", name+"-"+id, "with error:", err)
			}
		}

		// Openssl conf require...
		path := strings.ReplaceAll(config.GetRootDir(), "/", "\\") + `\Dependencys\openssl.cnf`

		if Utility.Exists(`C:\Program Files\Globular\Dependencys\openssl.cnf`) {
			err := Utility.SetWindowsEnvironmentVariable("OPENSSL_CONF", path)
			if err != nil {
				fmt.Println("fail to set environment variable OPENSSL_CONF with error:", err)
			}
		} else {
			fmt.Println("Open SSL configuration file ", path, "not found. Require to create environnement variable OPENSSL_CONF.")
		}
		err = Utility.SetWindowsEnvironmentVariable("Path", strings.ReplaceAll(systemPath, "/", "\\"))

		return err
	case "darwin":
		// Fix the path /usr/local/bin is not set by default...
		if Utility.Exists("/Library/LaunchDaemons/Globular.plist") {
			config, err := os.ReadFile("/Library/LaunchDaemons/Globular.plist")
			if err == nil {
				config_ := string(config)
				if !strings.Contains(config_, "<key>PATH</key>") {
					config_ = strings.ReplaceAll(config_, "</dict>",
						`
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/sbin</string>
	</dict>
	</dict>`)

					err = os.WriteFile("/Library/LaunchDaemons/Globular.plist", []byte(config_), 0600)
					if err != nil {
						fmt.Println("fail to update Globular.plist with error", err)
					}
				}
			}
		}
	}

	return nil
}

func refreshTokenPeriodically(ctx context.Context, globule *Globule) {
	ticker := time.NewTicker(time.Duration(globule.SessionTimeout)*time.Minute - 10*time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// The context is done, exit the goroutine
			return
		case <-ticker.C:
			// Refresh the token.
			err := globule.refreshLocalToken()
			if err != nil {
				fmt.Println("fail to refresh local token with error:", err)
			}
		}
	}
}

func (globule *Globule) restart() error {

	// stop watching for update.
	globule.exit_ = true

	// stop listening
	globule.exit <- true

	// stop the services and wait for them to be stopped.
	globule.cleanup()

	// force the restart the process.
	//autorestart.RestartByExec()
	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}

	// Get the current process arguments
	args := append([]string{execPath}, "restarted")
	args = append(args, os.Args[1:]...)

	// Get the current environment variables
	env := os.Environ()

	// Create a new command to execute the current process
	// #nosec G204 -- it always be the same executable
	cmd := exec.Command(execPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = env

	// Start the new process
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting new process: %v\n", err)
		os.Exit(1)
	}

	// Exit the current process
	fmt.Println("Exiting the current process...")
	os.Exit(0)
	return nil
}

/**
 * Here I will start the services manager who will start all microservices
 * installed on that computer.
 */
func (globule *Globule) startServices() error {

	fmt.Println("start services")

	// Here I will generate the keys for this server if not already exist.
	err := security.GeneratePeerKeys(globule.Mac)
	if err != nil {
		return err
	}

	// This is the local token...
	err = globule.refreshLocalToken()
	if err != nil {
		return err
	}

	// Retreive all configurations
	services, err := config.GetOrderedServicesConfigurations()
	if err != nil {
		return err
	}

	// start refresh local token...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Use context.WithTimeout if you want the goroutine to stop after a certain duration.
	// ctx, cancel := context.WithTimeout(context.Background(), someDuration)
	// defer cancel()
	go refreshTokenPeriodically(ctx, globule)

	start_port := Utility.ToInt(strings.Split(globule.PortsRange, "-")[0])
	end_port := Utility.ToInt(strings.Split(globule.PortsRange, "-")[1])

	// I will try to get the services manager configuration from the
	// services configurations list.

	for i := 0; i < len(services); i++ {

		if start_port >= end_port {
			return errors.New("no more available ports")
		}

		if err != nil {
			fmt.Println("fail to save service configuration with error ", err)
		} else if (len(globule.Certificate) > 0 && globule.Protocol == "https") || (globule.Protocol == "http") {

			service := services[i]

			service["State"] = "starting"
			name := service["Name"].(string)
			service["ProxyProcess"] = -1

			port := start_port + (i * 2)

			fmt.Println("try to start service ", name, " on port ", port, " and proxy port ", port+1)
			pid, err := process.StartServiceProcess(service, port)
			if err != nil {
				fmt.Println("fail to start service ", name, err)
			} else {
				service["Process"] = pid
				service["ProxyProcess"] = -1
				_, err = process.StartServiceProxyProcess(service, config.GetLocalCertificateAuthorityBundle(), config.GetLocalCertificate())
				if err != nil {
					fmt.Println("fail to start proxy for service ", name, err)
				}
			}
		}
	}

	// Here I will listen for logger event...
	go func() {

		// wait 2 second before register resource permissions and subscribe to log events.
		time.Sleep(2 * time.Second)

		for i := 0; i < len(services); i++ {
			s := services[i]
			if s["Permissions"] != nil {
				permissions := s["Permissions"].([]interface{})
				for j := 0; j < len(permissions); j++ {
					if permissions[j] != nil {
						err := globule.setActionResourcesPermissions(permissions[j].(map[string]interface{}))
						if err != nil {
							fmt.Println(" fail to register resource permission ", err)
						}
					}
				}
			}
		}

		// subscribe to log events
		err := globule.subscribe("new_log_evt", logListener(globule))
		if err != nil {
			fmt.Println("fail to subscribe to log events ", err)
		}

		// So here I will authenticate the root if the password is "adminadmin" that will
		// reset the password in the backend if it was manualy set in the config file.
		/*config_, err := config.GetLocalConfig(true)
		if err == nil {
			if config_["RootPassword"].(string) == "adminadmin" {

				address, _ := config.GetAddress()

				// Authenticate the user in order to get the token
				authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
				if err != nil {
					log.Println("fail to access resource service at "+address+" with error ", err)
					return
				}

				log.Println("authenticate user ", "sa", " at address ", address)
				_, err = authentication_client.Authenticate("sa", "adminadmin")

				if err != nil {
					log.Println("fail to authenticate user ", err)
					return
				}
			}
		}*/

	}()

	// wait for all services to be started.
	allServicesStarted := false
	nbTry := 20
	for !allServicesStarted {

		// Refresh the services list.
		services, err = config.GetServicesConfigurations()
		if err != nil {
			return err
		}

		// set the state of all services to running.
		allServicesStarted = true
		for i := 0; i < len(services); i++ {
			service := services[i]
			if service["State"].(string) != "running" {
				allServicesStarted = false
				nbTry--
				time.Sleep(1 * time.Second)
				break
			}
		}

		// if all services are not started after 20 second I will return an error.
		if nbTry == 0 {
			return errors.New("fail to start all services ")
		}
	}

	// Try to register to DNS...
	nbTry = 20
	for nbTry > 0 {
		err := globule.registerIpToDns()
		if err == nil {
			break
		}
		nbTry--
		time.Sleep(1 * time.Second)
	}

	return nil
}

/**
 * Update peers list.
 */
func updatePeersEvent(evt *eventpb.Event) {

	fmt.Println("update peers event received...", string(evt.Data))

	p := new(resourcepb.Peer)
	p_ := make(map[string]interface{}, 0)
	err := json.Unmarshal(evt.Data, &p_)
	if err != nil {
		fmt.Println("fail to update peer: ", p)
		return
	}

	p.Domain = p_["domain"].(string)
	p.Hostname = p_["hostname"].(string)
	p.Mac = p_["mac"].(string)
	p.Protocol = p_["protocol"].(string)

	if p_["local_ip_address"] != nil {
		p.LocalIpAddress = p_["local_ip_address"].(string)
	} else if p_["localIpAddress"] != nil {
		p.LocalIpAddress = p_["localIpAddress"].(string)
	}

	if p_["external_ip_address"] != nil {
		p.ExternalIpAddress = p_["external_ip_address"].(string)
	} else if p_["externalIpAddress"] != nil {
		p.ExternalIpAddress = p_["externalIpAddress"].(string)
	}

	state := Utility.ToInt(p_["state"])

	switch state {
	case 0:
		p.State = resourcepb.PeerApprovalState_PEER_PENDING
	case 1:
		p.State = resourcepb.PeerApprovalState_PEER_ACCETEP
	case 2:
		p.State = resourcepb.PeerApprovalState_PEER_REJECTED
	}

	httpPort := Utility.ToInt(p_["PortHTTP"])
	if httpPort > math.MaxInt32 || httpPort < math.MinInt32 {
		fmt.Printf("port value %d out of int32 range\n", httpPort)
		return
	}
	// #nosec G115 -- value has been validated above
	p.PortHTTP = int32(httpPort)

	httpsPort := Utility.ToInt(p_["PortHTTPS"])
	if httpsPort > math.MaxInt32 || httpsPort < math.MinInt32 {
		fmt.Printf("port value %d out of int32 range\n", httpsPort)
		return
	}
	// #nosec G115 -- value has been validated above
	p.PortHTTPS = int32(httpsPort)

	if p_["actions"] != nil {
		p.Actions = make([]string, len(p_["actions"].([]interface{})))

		for i := 0; i < len(p_["actions"].([]interface{})); i++ {
			p.Actions[i] = p_["actions"].([]interface{})[i].(string)
		}
	} else {
		p.Actions = make([]string, 0)
	}
	// #nosec G115 -- Allowing assignment of port number from trusted source
	p.PortHTTPS = int32(httpsPort)

	if p_["actions"] != nil {
		p.Actions = make([]string, len(p_["actions"].([]interface{})))

		for i := 0; i < len(p_["actions"].([]interface{})); i++ {
			p.Actions[i] = p_["actions"].([]interface{})[i].(string)
		}
	} else {
		p.Actions = make([]string, 0)
	}

	globule.peers.Store(p.Mac, p)

	err = globule.savePeers()
	if err != nil {
		fmt.Println("fail to save peers with error:", err)
	}

	// set the peer ip in the /etc/hosts file.
	if Utility.MyIP() == p.ExternalIpAddress {
		err := globule.setHost(p.LocalIpAddress, p.Hostname+"."+p.Domain)
		if err != nil {
			fmt.Println("fail to set host with error:", err)
		}
	}

}

func deletePeersEvent(evt *eventpb.Event) {
	globule.peers.Delete(string(evt.Data))
	err := globule.savePeers()
	if err != nil {
		fmt.Println("fail to save peers with error:", err)
	}
}

func (globule *Globule) initPeer(p *resourcepb.Peer) error {

	// Here I will try to set the peer ip...
	address := p.Hostname
	if p.Domain != "localhost" {
		address += "." + p.Domain
	} else if globule.Domain != "localhost" && p.Protocol == "https" {
		// in that case I will use the globule domain, the peer must be in the same domain anyway...
		address += "." + globule.Domain
	}

	// set the peer ip in the /etc/hosts file.
	if Utility.MyIP() == p.ExternalIpAddress {
		err := globule.setHost(p.LocalIpAddress, address)
		if err != nil {
			fmt.Println("fail to set host with error:", err)
		}
	}

	if p.Protocol == "https" {
		address += ":" + Utility.ToString(p.PortHTTPS)
	} else {
		address += ":" + Utility.ToString(p.PortHTTP)
	}

	// Here I will get the peer public key if not already exist.
	if !Utility.Exists(globule.config + "/keys/" + strings.ReplaceAll(p.Mac, ":", "_") + "_public") {

		// get the peer public key.
		rqst := p.Protocol + "://" + address + "/public_key"
		// #nosec G107 -- Ok
		resp, err := http.Get(rqst)
		if err == nil {

			defer func() {
				if cerr := resp.Body.Close(); cerr != nil {
					fmt.Printf("warning: failed to close response body: %v\n", cerr)
				}
			}()

			body, err := io.ReadAll(resp.Body)
			if err == nil {
				// save the peer public key.
				err = os.WriteFile(globule.config+"/keys/"+strings.ReplaceAll(p.Mac, ":", "_")+"_public", body, 0600)
				if err != nil {
					fmt.Println("fail to save peer public key with error: ", err)
					return err
				}
			} else {
				fmt.Println("fail to read peer public key with error: ", err)
				return err
			}
		} else {
			fmt.Println("fail to get peer public key with error: ", err)
			return err
		}
	}

	// Now I will keep it in the peers list.
	globule.peers.Store(p.Mac, p)

	// Here I will try to update
	token, err := security.GenerateToken(globule.SessionTimeout, p.GetMac(), "sa", "", globule.AdminEmail, globule.Domain)
	if err != nil {
		return err
	}
	// no wait here...

	// update local peer info for each peer...
	resource_client__, err := getResourceClient(address)
	if err != nil {
		return err
	}

	// retreive the local peer infos
	peers_, err := resource_client__.GetPeers(`{"mac":"` + globule.Mac + `"}`)
	if err != nil {
		return err
	}

	if len(peers_) > 0 {
		// set mutable values...
		peer_ := peers_[0]
		peer_.Protocol = globule.Protocol
		peer_.LocalIpAddress = config.GetLocalIP()
		peer_.ExternalIpAddress = Utility.MyIP()

		// #nosec G115 -- Allowing assignment of port number from trusted source
		peer_.PortHTTP = int32(globule.PortHTTP)

		// #nosec G115 -- Allowing assignment of port number from trusted source
		peer_.PortHTTPS = int32(globule.PortHTTPS)

		peer_.Domain = globule.Domain
		err := resource_client__.UpdatePeer(token, peer_)
		if err != nil {
			return err
		}
	} else {
		return errors.New("fail to retreive local peer info " + globule.Mac)
	}

	return nil
}

func (globule *Globule) getPeers() ([]*resourcepb.Peer, error) {
	peers := make([]*resourcepb.Peer, 0)
	address, _ := config.GetAddress()

	resource_client_, err := getResourceClient(address)
	if err != nil {
		return nil, err
	}

	// Return the registered peers
	nbTry := 10
	for i := 0; i < nbTry; i++ {
		peers, err = resource_client_.GetPeers("")
		if err != nil {
			fmt.Println("fail to get peers with error ", err)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	return peers, nil
}

/**
 * Here I will init the list of peers.
 */
func (globule *Globule) initPeers() error {

	// Return the registered peers
	peers, err := globule.getPeers()
	if err != nil {
		return err
	}

	// Now I will set peers in the host file.
	for i := range peers {
		p := peers[i]

		// Set existing value...
		globule.peers.Store(p.Mac, p)

		// Try to update with updated infos...
		go func(p *resourcepb.Peer) {
			err := globule.initPeer(p)
			if err != nil {
				globule.peers.Delete(p.Mac) // remove the peer from the list.
				err := globule.savePeers()  // save the peers list.
				if err != nil {
					fmt.Println("fail to save peers with error:", err)
				}
			}
		}(p)
	}

	// Subscribe to new peers event...
	err = globule.subscribe("update_peers_evt", updatePeersEvent)
	if err != nil {
		fmt.Println("fail to subscribe to update_peers_evt with error:", err)
	}

	err = globule.subscribe("delete_peer_evt", deletePeersEvent)
	if err != nil {
		fmt.Println("fail to subscribe to delete_peer_evt with error:", err)
	}

	// Now I will set the local peer info...
	err = globule.savePeers()
	if err != nil {
		fmt.Println("fail to save peers with error:", err)
	}

	return nil // here if some errors occurred the peers list may be inconsistent.
}

// func (globule *Globule) getHttpClient

/**
 * Stop all services.
 */
func (globule *Globule) stopServices() error {

	// Now I will set configuration values
	services_configs, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	for i := range services_configs {
		fmt.Println("stop service ", services_configs[i]["Name"], " with id ", services_configs[i]["Id"])
		pid := Utility.ToInt(services_configs[i]["Process"])
		services_configs[i]["State"] = "killed"
		proxyPid := Utility.ToInt(services_configs[i]["ProxyProcess"])
		services_configs[i]["ProxyProcess"] = -1

		// save config...
		err := config.SaveServiceConfiguration(services_configs[i])
		if err == nil {
			if pid > 0 {
				// Kill the process.
				process, err := os.FindProcess(pid)
				if err == nil {
					err = process.Signal(syscall.SIGTERM) // make the process stop gracefully.
					if err != nil {
						fmt.Println("Error sending signal:", err)
					}
				}

			}

			// Kill the proxy process.
			if proxyPid > 0 {
				process, err := os.FindProcess(proxyPid)
				if err == nil {
					err = process.Signal(syscall.SIGTERM) // make the process stop gracefully.
					if err != nil {
						fmt.Println("Error sending signal:", err)
					}
				}
			}
		}

	}

	return nil
}

// Start http/https server...
func (globule *Globule) serve() error {

	// Create the admin account.
	err := globule.registerAdminAccount()
	if err != nil {
		fmt.Println("fail to register admin account with error:", err)
	}

	url := globule.Protocol + "://" + globule.getAddress()
	switch globule.Protocol {
	case "https":
		if globule.PortHTTPS != 443 {
			url += ":" + Utility.ToString(globule.PortHTTPS)
		}
	case "http":
		if globule.PortHTTP != 80 {
			url += ":" + Utility.ToString(globule.PortHTTP)
		}
	}

	elapsed := time.Since(globule.startTime)

	fmt.Println("globular version " + globule.Version + " build " + Utility.ToString(globule.Build) + " listen at address " + url)
	fmt.Printf("startup took %s\n", elapsed)

	// create applications connections
	err = globule.createApplicationConnections()
	if err != nil {
		return err
	}

	return nil

}

/**
 * Start the control plane to manage the cluster configuration.
 */
func (globule *Globule) initControlPlane() {

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Set up a signal handler to cancel the context on interrupt signals
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Received interrupt signal. Shutting down...")
		cancel()
	}()

	// Use a WaitGroup to wait for graceful shutdown
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the control plane in a goroutine
	go func() {
		defer wg.Done()

		// so here I will read the envoy yaml configuration file and set it to the control plane.
		configPath := config.GetConfigDir() + "/envoy.yml"
		// #nosec G304 -- Ok
		data, err := os.ReadFile(configPath)
		if err != nil {
			if !Utility.Exists(configPath) {
				// Here I will create the config file for envoy.

				config_ := `
node:
    cluster: globular-cluster
    id: globular-xds

dynamic_resources:
    ads_config:
      api_type: GRPC
      transport_api_version: V3
      grpc_services:
      - envoy_grpc:
          cluster_name: xds_cluster
    cds_config:
        resource_api_version: V3
        ads: {}
    lds_config:
        resource_api_version: V3
        ads: {}

static_resources:
    clusters:
      - type: STRICT_DNS
        typed_extension_protocol_options:
          envoy.extensions.upstreams.http.v3.HttpProtocolOptions:
            "@type": type.googleapis.com/envoy.extensions.upstreams.http.v3.HttpProtocolOptions
            explicit_http_config:
               http2_protocol_options: {}
        name: xds_cluster
        load_assignment:
            cluster_name: xds_cluster
            endpoints:
            - lb_endpoints:
              - endpoint:
                  address:
                    socket_address:
                       address: 0.0.0.0
                       port_value: 9900

admin:
    access_log_path: /dev/null
    address:
      socket_address:
        address: 0.0.0.0
        port_value: 9901
`

				// Read the content of the YAML file
				err := os.WriteFile(configPath, []byte(config_), 0600)
				if err != nil {
					fmt.Println("fail to create envoy configuration file with error ", err)
					os.Exit(1)
				}

				data = []byte(config_)
			}

		}

		config := make(map[string]interface{})
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Println("fail to unmarshal envoy configuration file with error ", err)
			return
		}

		// Get the port from the envoy configuration file.
		port := config["static_resources"].(map[string]interface{})["clusters"].([]interface{})[0].(map[string]interface{})["load_assignment"].(map[string]interface{})["endpoints"].([]interface{})[0].(map[string]interface{})["lb_endpoints"].([]interface{})[0].(map[string]interface{})["endpoint"].(map[string]interface{})["address"].(map[string]interface{})["socket_address"].(map[string]interface{})["port_value"].(int)

		// TODO: Make the control plane port configurable, it must set also in envoy.yml
		// #nosec G115 -- Ok
		if err := controlplane.StartControlPlane(ctx, uint(port), globule.exit); err != nil {
			fmt.Printf("Error starting control plane: %v\n", err)
		}
	}()

	// Wait for either the control plane to finish or the context to be canceled
	<-ctx.Done()
	// Context canceled, wait for the control plane to finish
	wg.Wait()
	fmt.Println("Graceful shutdown complete.")

	fmt.Println("Exiting...")
}

func SetSnapshot() error {
	// Start envoy proxy.
	// Now I will add services to envoy configuration.
	services, _ := config.GetServicesConfigurations()

	spnapShots := make([]controlplane.Snapshot, 0)

	// Add services to envoy configuration.
	proxies := make(map[uint32]bool, 0)

	domain, _ := config.GetAddress()
	domain = strings.Split(domain, ":")[0]

	// TODO: Add peers services on the endpoint.
	// I will add it if the proxy is not already set.
	for i := 0; i < len(services); i++ {
		service := services[i]

		host := strings.Split(service["Address"].(string), ":")[0]
		// #nosec G115 -- Ok
		proxy := uint32((Utility.ToInt(service["Proxy"])))

		if _, ok := proxies[proxy]; !ok {
			snapshot := controlplane.Snapshot{

				ClusterName:  strings.ReplaceAll(service["Name"].(string), ".", "_") + "_cluster",
				RouteName:    strings.ReplaceAll(service["Name"].(string), ".", "_") + "_route",
				ListenerName: strings.ReplaceAll(service["Name"].(string), ".", "_") + "_listener",
				ListenerPort: proxy,
				ListenerHost: "0.0.0.0", // local address.

				// #nosec G115 -- Ok
				EndPoints: []controlplane.EndPoint{{Host: host, Port: uint32(Utility.ToInt(service["Port"])), Priority: 100}},

				// grpc certificate...
				ServerCertPath: service["CertFile"].(string),
				KeyFilePath:    service["KeyFile"].(string),
				CAFilePath:     service["CertAuthorityTrust"].(string),

				// Let's encrypt certificate...
				CertFilePath:   config.GetConfigDir() + "/tls/" + domain + "/" + config.GetLocalCertificate(),
				IssuerFilePath: config.GetConfigDir() + "/tls/" + domain + "/" + config.GetLocalCertificateAuthorityBundle(),
			}

			// Certificates are generated by
			// TODO : Add endpoint for service that can be run on other peers.
			// TEST only for resource services with SCYLLA DB
			/*if services[i]["Name"].(string) == "echo.EchoService" ||
				services[i]["Name"].(string) == "resource.ResourceService" {
				for _, p := range globule.Peers {

					peer := p.(map[string]interface{})

					address := peer["Hostname"].(string)
					if peer["Domain"].(string) != "localhost" {
						address += "." + peer["Domain"].(string)
					} else if globule.Domain != "localhost" {
						address += "." + globule.Domain
					}

					port := Utility.ToInt(peer["Port"])

					remoteService, err := config.GetRemoteServiceConfig(address, port, service["Name"].(string))
					if err == nil {
						endpoint := controlplane.EndPoint{Host: address, Port: uint32(Utility.ToInt(remoteService["Port"])), Priority: 80}
						// I will add the endpoint only if the domain is the same.
						if strings.HasSuffix(domain, remoteService["Domain"].(string)) {
							snapshot.EndPoints = append(snapshot.EndPoints, endpoint)
						}
					} else {
						if strings.Contains(err.Error(), "context deadline exceeded") {
							// I will remove the peer from the list...
							fmt.Println("remove peer ", peer["Mac"].(string), " from the list because it is not reachable")
							globule.peers.Delete(peer["Mac"].(string))
						}
					}
				}
			}*/

			proxies[proxy] = true
			spnapShots = append(spnapShots, snapshot)
		} else {
			fmt.Println("proxy ", proxy, " already set for service ", service["Name"].(string))
		}
	}

	return controlplane.AddSnapshot("globular-xds", "1", spnapShots)
}

// Start envoy as a proxy.
func startEnvoyProxy() {

	go func() {
		err := SetSnapshot()
		if err != nil {
			fmt.Println("fail to generate envoy dynamic configuration with error", err)
			//return err
		}

		// Now I will start the envoy proxy.
		err = process.StartEnvoyProxy()
		if err != nil {
			fmt.Println("fail to start envoy proxy with error ", err)
			time.Sleep(5 * time.Second) // wait 5 second before retrying...
			startEnvoyProxy()
		}
	}()
}

/**
 * Start serving the content.
 */
func (globule *Globule) Serve() error {

	// So here if another instance of the server exist I will kill it.
	pids, err := Utility.GetProcessIdsByName("Globular")
	if err == nil {

		for i := range pids {
			if pids[i] != os.Getpid() {
				err := Utility.TerminateProcess(pids[i], 0)
				if err != nil {
					fmt.Println("fail to terminate process with error:", err)
				}
			}
		}
	}

	// Initialyse directories.
	err = globule.initDirectories()
	if err != nil {
		fmt.Println("fail to initialize directories with error:", err)
		return err
	}

	/*
		// I will now start etcd server.
		go func() {
			err = process.StartEtcdServer()
			if err != nil {
				fmt.Println("fail to start etcd kv store ", err)
				os.Exit(1) // exit with error...
			}
		}()
	*/

	// start listen to http(s)
	// service must be able to get their configuration via http...
	err = globule.Listen()
	if err != nil {
		fmt.Println("fail to start http server ", err)
		os.Exit(1) // exit with error...
		return err
	}

	// Start microservice manager.
	err = globule.startServices()
	if err != nil {
		fmt.Println("fail to start services with error:", err)
		return err
	}

	// Start process monitoring with prometheus.
	err = process.StartProcessMonitoring(globule.Protocol, globule.PortHTTP, globule.exit)
	if err != nil {
		fmt.Println("fail to start process monitoring with error:", err)
		return err
	}

	// Watch config.
	globule.watchConfig()

	// Set the fmt information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// First of all i will set the local host found...
	go func() {
		hosts := Utility.GetHostnameIPMap(globule.LocalIpAddress)
		for k, v := range hosts {
			err := globule.setHost(k, v)
			if err != nil {
				fmt.Println("fail to set host with error:", err)
			}
		}

		// Try with ip address...
		ips, err := Utility.ScanIPs()
		if err == nil {
			for i := range ips {
				config_, err := config.GetRemoteConfig(ips[i], 80)
				if err == nil {
					hostname := config_["Name"].(string)
					if config_["Domain"] != nil {
						if config_["Domain"].(string) != "localhost" {
							hostname += "." + config_["Domain"].(string)
						}
					}

					err = globule.setHost(ips[i], hostname)
					if err != nil {
						fmt.Println("fail to set host with error:", err)
					}

				}
			}
		}
	}()

	// Initialize peers
	errCh := make(chan error, 1)

	go func() {
		errCh <- globule.initPeers()
	}()

	if err := <-errCh; err != nil {
		fmt.Println("initPeers error:", err)
	}

	// Now I will initialize the control plane.
	//go globule.initControlPlane()

	// Start envoy proxy.
	// startEnvoyProxy()
	//StartImprobableProxy()

	p := make(map[string]interface{})

	p["address"] = globule.getAddress()
	p["domain"], _ = config.GetDomain()
	p["hostname"] = globule.Name
	p["mac"] = globule.Mac
	p["PortHTTP"] = globule.PortHTTP
	p["PortHTTPS"] = globule.PortHTTPS

	jsonStr, _ := json.Marshal(&p)

	// set services configuration values
	err = globule.publish("start_peer_evt", jsonStr)
	if err != nil {
		fmt.Println("fail to publish start_peer_evt with error:", err)
		return err
	}

	err = globule.serve()
	if err != nil {
		return err
	}

	return nil
}

func (globule *Globule) watchConfig() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}

	defer func() {
		if err := watcher.Close(); err != nil {
			fmt.Println("failed to close watcher:", err)
		}
	}()

	go func() {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op == fsnotify.Write {
				// renit the service...
				file, _ := os.ReadFile(globule.config + "/config.json")
				config := make(map[string]interface{})
				err := json.Unmarshal(file, &config)
				if err == nil {
					err = globule.setConfig(config)
					if err != nil {
						fmt.Println("fail to set config with error:", err)
						os.Exit(1) // configuration is invalid
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		}

	}()

	// watch for configuration change
	err = watcher.Add(globule.config + "/config.json")
	if err != nil {
		log.Fatal("Add failed:", err)
	}
}

/**
 * Return the globule address.
 */
func (globule *Globule) getAddress() string {
	address, _ := config.GetAddress()

	return address
}

/**
 * The local server address.
 */
func (globule *Globule) getLocalDomain() string {
	address, _ := config.GetAddress()
	domain := strings.Split(address, ":")[0]
	return domain
}

// setHost sets the mapping between an IPv4 address and a domain name in the system's hosts file.
// It handles special cases for "localhost" and ".localhost" addresses, ensuring that local addresses
// are not overwritten by non-local addresses. If the address already exists and is local, it prevents
// replacement by a non-local address and returns an error. The function uses the txeh library to
// manipulate the hosts file and returns an error if any operation fails.
//
// Parameters:
//
//	ipv4   - The IPv4 address to associate with the domain.
//	address - The domain name to map to the IPv4 address.
//
// Returns:
//
//	error - An error if the mapping could not be set or saved, or if invalid parameters are provided.
func (globule *Globule) setHost(ipv4, address string) error {
	if strings.HasSuffix(address, ".localhost") {
		return nil
	}

	if address == "localhost" {
		ipv4 = "127.0.0.1" // force ipv4 address for localhost
	}

	if len(ipv4) == 0 {
		return errors.New("no ipv4 address to set")
	}

	if len(address) == 0 {
		return errors.New("no domain to set")
	}

	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	// Here I will test if the previous address is a local address...
	exist, address_, _ := hosts.HostAddressLookup(address, txeh.IPFamilyV4)
	if exist {
		if Utility.IsLocal(address_) && !Utility.IsLocal(ipv4) {
			// If the previous address was a local address I will not replace it by a non local address...
			// The hosts file must be edited manually.
			return errors.New("previous address was a local address, cannot be replace by a non local address")
		}
	}

	hosts.AddHost(ipv4, address)
	err = hosts.Save()
	if err != nil {
		fmt.Println("fail to save hosts ", ipv4, address, " with error ", err)
	}

	return err
}

/**
 * Set the ip for a given domain or sub-domain.
 * The domain must be manage by the dns provider directly.
 */
func (globule *Globule) registerIpToDns() error {

	// Globular DNS is use to create sub-domain.
	// ex: globular1.globular.io here globular.io is the domain and globular1 is
	// the sub-domain. Domain must be manage by dns provider directly, by using
	// the DnsSetA (set ip api call)... see the next part of that function
	// for more information.
	if len(globule.DNS) > 0 {
		// Here I will set dns in the resolv.conf file
		resolv_conf := "# That file was generated by globular at server startup. To reset to it original move the file resolv.conf_ to resolv.conf\n"
		resolv_conf += "nameserver 8.8.8.8\n"
		resolv_conf += "nameserver 1.1.1.1\n"

		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS, "dns.DnsService")
		if err != nil {
			fmt.Println("fail to create dns client with error ", err)
			return err
		}

		// if the dns server is running...
		if globule.DNS == globule.Name+"."+globule.Domain {
			dns_server_is_running := false
			nbTry := 20
			for !dns_server_is_running {
				// I will get the service configuration.
				dns_server_config, err := config.GetServiceConfigurationById(dns_client_.GetId())
				if err != nil {
					fmt.Println("fail to get dns server configuration with error ", err)
					return err
				}

				if dns_server_config["State"].(string) == "running" {
					dns_server_is_running = true
				} else {
					time.Sleep(1 * time.Second)
				}

				nbTry--
				if nbTry == 0 {
					fmt.Println("fail to get dns server configuration with error ", err)
					return err
				}
			}
		}

		defer dns_client_.Close()

		// Here the token must be generated for the dns server...
		// That peer must be register on the dns to be able to generate a valid token.
		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail, globule.Domain)
		if err != nil {
			fmt.Println("fail to generate token for dns server with error ", err)
			return err
		}

		// try to set the ipv6 address...
		ipv6, err := Utility.MyIPv6()
		if err == nil {
			_, err = dns_client_.SetAAAA(token, globule.getLocalDomain(), ipv6, 60)
			if err != nil {
				fmt.Println("fail to set AAAA  domain ", globule.getLocalDomain(), " with error ", err)
				return err
			}
			fmt.Println("set AAAA record for domain ", globule.getLocalDomain(), " with success")
		}

		// I will set alternate domain only if the globule is the master.
		if globule.DNS == globule.getLocalDomain() {

			// Here I will set the A record for the globular domain.
			err = dns_client_.RemoveA(token, globule.getLocalDomain())
			if err != nil {
				fmt.Println("fail to remove A record for domain ", globule.getLocalDomain(), " with error ", err)
			}

			_, err = dns_client_.SetA(token, globule.getLocalDomain(), Utility.MyIP(), 60)
			if err != nil {
				fmt.Println("fail to set A record for alternate domain ", globule.getLocalDomain(), " with error ", err)
				return err
			}

			fmt.Println("set A record for alternate domain ", globule.getLocalDomain(), Utility.MyIP(), " with success")
			for j := 0; j < len(globule.AlternateDomains); j++ {

				// Here I will set the A record for the alternate domain.
				alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
				_, err = dns_client_.SetA(token, alternateDomain, Utility.MyIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					return err
				}
				fmt.Println("set A record for alternate domain ", alternateDomain, " with success")

				/*_, err = dns_client_.SetA(token, alternateDomain, config.GetLocalIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					continue
				} else {
					fmt.Println("set A record for alternate domain ", alternateDomain, " with success")
				}*/

				_, err = dns_client_.SetA(token, alternateDomain, Utility.MyIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					return err
				}
				fmt.Println("set A record for alternate domain ", alternateDomain, Utility.MyIP(), " with success")

				_, err = dns_client_.SetAAAA(token, alternateDomain, ipv6, 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					return err
				}
				fmt.Println("set AAAA record for alternate domain ", alternateDomain, " with success")

			}
		}

		// I will publish the private ip address only
		_, err = dns_client_.SetA(token, "mail."+globule.Domain, Utility.MyIP(), 60)
		if err != nil {
			fmt.Println("fail to set A record for domain ", "mail."+globule.Domain, " with error ", err)
			return err
		}

		// Now the mx record.
		err = dns_client_.SetMx(token, globule.Domain, 10, "mail."+globule.Domain, 60)
		if err != nil {
			fmt.Println("fail to set MX record for domain ", globule.Domain, " with error ", err)
		}

		// SPF record
		err = dns_client_.RemoveText(token, globule.Domain+".")
		if err != nil {
			fmt.Println("fail to remove TXT record for domain ", globule.Domain, " with error ", err)
		}

		spf := fmt.Sprintf(`v=spf1 mx ip4:%s include:_spf.google.com ~all`, Utility.MyIP())
		err = dns_client_.SetText(token, globule.Domain+".", []string{spf}, 60)
		if err != nil {
			fmt.Println("fail to set TXT record for domain ", globule.Domain, " with error ", err)
			return err
		}

		// DMARC record
		dmarc_policy := fmt.Sprintf(`v=DMARC1;p=quarantine;rua=mailto:%s;ruf=mailto:%s;adkim=r;aspf=r;pct=100`, globule.AdminEmail, globule.AdminEmail)
		err = dns_client_.RemoveText(token, "_dmarc."+globule.Domain+".")
		if err != nil {
			fmt.Println("fail to remove TXT record for domain ", "_dmarc."+globule.Domain, " with error ", err)
		}

		err = dns_client_.SetText(token, "_dmarc."+globule.Domain+".", []string{dmarc_policy}, 60)
		if err != nil {
			fmt.Println("fail to set TXT record for domain ", "_mta-sts."+globule.Domain, " with error ", err)
			return err
		}

		// now the  MTA-STS policy
		if !Utility.Exists(config.GetConfigDir() + "/tls/" + globule.Name + "." + globule.Domain + "/mta-sts.txt") {

			mta_sts_policy := fmt.Sprintf(`version: STSv1
mode: enforce
mx: %s
ttl: 86400
		`, globule.Domain)

			err = os.WriteFile(config.GetConfigDir()+"/tls/"+globule.Name+"."+globule.Domain+"/mta-sts.txt", []byte(mta_sts_policy), 0600)
			if err != nil {
				fmt.Println("fail to write mta-sts policy with error ", err)
				return err
			}
		}

		// endpoints to retrieve the policy
		_, err = dns_client_.SetA(token, "mta-sts."+globule.Domain, Utility.MyIP(), 60)
		if err != nil {
			fmt.Println("fail to set A record for domain ", "mta-sts."+globule.Domain, " with error ", err)
			return err
		}

		err = dns_client_.RemoveText(token, "_mta-sts."+globule.Domain+".")
		if err != nil {
			fmt.Println("fail to remove TXT record for domain ", "_mta-sts."+globule.Domain, " with error ", err)
		}

		err = dns_client_.SetText(token, "_mta-sts."+globule.Domain+".", []string{"v=STSv1; id=cd1e8e2f-311c-3c55-bb5a-cc1eedee398e;"}, 60)
		if err != nil {
			fmt.Println("fail to set TXT record for domain ", "_mta-sts."+globule.Domain, " with error ", err)
			return err
		}

		_, err = dns_client_.SetAAAA(token, "mail."+globule.Domain, ipv6, 60)
		if err != nil {
			fmt.Println("fail to set AAAA record for domain ", "mail."+globule.Domain, " with error ", err)
			return err
		}

		// if the NS server are local I will set the local ip address.
		for j := 0; j < len(globule.NS); j++ {

			// Set it nameservers are part of the domain.
			ns := globule.NS[j].(string)

			//if strings.HasSuffix(ns, globule.Domain) {
			_, err = dns_client_.SetA(token, globule.NS[j].(string), Utility.MyIP(), 60)
			if err != nil {
				fmt.Println("fail to set A record for NS server ", ns, " with error ", err)
				return err
			}

			_, err = dns_client_.SetAAAA(token, globule.NS[j].(string), ipv6, 60)
			if err != nil {
				fmt.Println("fail to set AAAA record for domain ", ns, " with error ", err)
				return err
			}
		}

		// Now the SOA record.
		serial := uint32(1)
		refresh := uint32(86400)
		retry := uint32(7200)
		expire := uint32(4000000)
		ttl := uint32(11200)

		// Now i will set the NS record.
		for j := 0; j < len(globule.AlternateDomains); j++ {
			for k := 0; k < len(globule.NS); k++ {
				alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
				if !strings.HasSuffix(alternateDomain, ".") {
					alternateDomain += "."
				}

				ns := globule.NS[k].(string)
				if !strings.HasSuffix(ns, ".") {
					ns += "."
				}

				err = dns_client_.SetNs(token, alternateDomain, ns, 60)
				if err != nil {
					fmt.Println("fail to set NS record for alternate domain ", alternateDomain, " with error ", err)
					return err
				}
			}
		}

		for j := 0; j < len(globule.AlternateDomains); j++ {
			for k := 0; k < len(globule.NS); k++ {
				alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
				if !strings.HasSuffix(alternateDomain, ".") {
					alternateDomain += "."
				}

				ns := globule.NS[k].(string)
				if !strings.HasSuffix(ns, ".") {
					ns += "."
				}

				// Be sure the email is valid...
				email := globule.AdminEmail
				if len(email) == 0 {
					email = "admin@" + globule.Domain
				}

				if !strings.Contains(email, "@") {
					email = "admin@" + globule.Domain
				}

				if !strings.HasSuffix(email, ".") {
					email += "."
				}

				// Now I will set the SOA record.
				err := dns_client_.SetSoa(token, alternateDomain, ns, email, serial, refresh, retry, expire, ttl, ttl)
				if err != nil {
					fmt.Println("fail to set NS record for alternate domain ", alternateDomain, " with error ", err)
					return err
				}

			}
		}

		for j := 0; j < len(globule.AlternateDomains); j++ {
			alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
			// Now I will set the CAA record.
			err = dns_client_.SetCaa(token, alternateDomain+".", 0, "issue", "letsencrypt.org", 60)
			if err != nil {
				fmt.Println("fail to set CAA record for alternate domain ", globule.AlternateDomains[j], " with error ", err)
				return err
			}
		}

		// save the file to /etc/resolv.conf
		if Utility.Exists("/etc/resolv.conf") {
			err := Utility.MoveFile("/etc/resolv.conf", "/etc/resolv.conf_")
			if err != nil {
				fmt.Println("fail to move /etc/resolv.conf to /etc/resolv.conf_ with error ", err)
				return err
			}

			err = Utility.WriteStringToFile("/etc/resolv.conf", resolv_conf)
			if err != nil {
				fmt.Println("fail to write to /etc/resolv.conf with error ", err)
				return err
			}
		}
	}

	// Here If the DNS provides has api to update the ip address I will use it.
	for i := 0; i < len(globule.DnsUpdateIpInfos); i++ {

		// the api call "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
		// example,
		// {"SetA":"https://api.godaddy.com/v1/domains/globular.io/records/A/@", "Key":"", "Secret":""}
		setA := globule.DnsUpdateIpInfos[i].(map[string]interface{})["SetA"].(string)
		key := globule.DnsUpdateIpInfos[i].(map[string]interface{})["Key"].(string)
		secret := globule.DnsUpdateIpInfos[i].(map[string]interface{})["Secret"].(string)

		// set the data to the actual ip address.
		data := `[{"data":"` + Utility.MyIP() + `"}]`

		// initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, setA, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return err
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "sso-key "+key+":"+secret)

		// execute the request.
		_, err = client.Do(req)
		if err != nil {
			return (err)
		}
	}

	return globule.setHost(config.GetLocalIP(), globule.getLocalDomain())
}

/**
 * retreive checksum from the server.
 */
func getChecksum(address string, port int) (string, error) {
	if len(address) == 0 {
		return "", errors.New("no address was given")
	}

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	var url = "http://" + address + ":" + Utility.ToString(port) + "/checksum"
	// #nosec G107 -- Ok
	resp, err = http.Get(url)

	if err != nil || resp.StatusCode != http.StatusCreated {
		url = "https://" + address + ":" + Utility.ToString(port) + "/checksum"
		// #nosec G107 -- Ok
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("warning: failed to close response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	}

	return "", errors.New("fail to retreive checksum with error " + Utility.ToString(resp.StatusCode))
}

/**
 *  Watch if globular need to be update.
 */
func (globule *Globule) watchForUpdate() {
	go func() {

		// Now the service...
		for !globule.exit_ {
			// stop watching if exit was call...
			if len(globule.Discoveries) > 0 {
				// Here I will retreive the checksum information from it parent.
				discovery := globule.Discoveries[0]
				address := strings.Split(discovery, ":")[0]
				port := 80
				if strings.Contains(discovery, ":") {
					port = Utility.ToInt(strings.Split(discovery, ":")[1])
				}

				execPath := Utility.GetExecName(os.Args[0])

				// Here I will test if the checksum has change...
				checksum, err := getChecksum(address, port)
				if Utility.Exists(config.GetRootDir() + "/Globular") {
					execPath = config.GetRootDir() + "/Globular"
				}

				if err == nil {
					if checksum != Utility.CreateFileChecksum(execPath) {

						err := updateGlobularFrom(discovery, globule.getLocalDomain(), "sa", globule.RootPassword, runtime.GOOS+":"+runtime.GOARCH)
						if err != nil {
							fmt.Println("fail to update globular from " + discovery + " with error " + err.Error())
							return
						}

						// Here I will restart the server.
						err = globule.restart()
						if err != nil {
							fmt.Println("fail to restart globular with error ", err)
						}

					}
				} else {
					fmt.Println("fail to get checksum from : ", discovery, " error: ", err)
				}

			}

			services, err := config.GetServicesConfigurations()
			if err == nil {
				// get the resource client
				for i := range services {
					s := services[i]
					values := strings.Split(s["PublisherID"].(string), "@")
					if len(values) == 2 {
						resource_client_, err := getResourceClient(values[1])
						if err == nil {
							// test if the service need's to be updated, test if the part is part of installed instance and not developement environement.
							if s["KeepUpToDate"].(bool) && strings.Contains(s["Path"].(string), "/globular/services/") {
								// Here I will get the last version of the package...
								descriptor, err := resource_client_.GetPackageDescriptor(s["Id"].(string), s["PublisherID"].(string), "")

								if err == nil {
									descriptorVersion := Utility.NewVersion(descriptor.Version)
									serviceVersion := Utility.NewVersion(s["Version"].(string))
									if descriptorVersion.Compare(serviceVersion) == 1 {
										// TODO keep service up to date.
										fmt.Println("service ", s["Name"].(string), s["Id"].(string), "will be updated")
										address, _ := config.GetAddress()
										servicesManager, err := GetServiceManagerClient(address)
										if err == nil {
											if servicesManager.StopServiceInstance(s["Id"].(string)) == nil {
												token, _ := security.GetLocalToken(globule.Mac)
												if servicesManager.UninstallService(token, s["Domain"].(string), s["PublisherID"].(string), s["Id"].(string), s["Version"].(string)) == nil {
													err = servicesManager.InstallService(token, s["Domain"].(string), s["PublisherID"].(string), s["Id"].(string), descriptor.Version)
													if err != nil {
														fmt.Println("fail to install service ", s["Name"].(string), s["Id"].(string), " with error ", err)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}

			} else {
				fmt.Println("fail to get service configurations with err ", err)
			}

			// The time here can be set to higher value.
			time.Sleep(time.Duration(globule.WatchUpdateDelay) * time.Second)
		}
	}()
}

// Try to display application message in a nice way.
func logListener(g *Globule) func(evt *eventpb.Event) {
	return func(evt *eventpb.Event) {

		info := new(logpb.LogInfo)
		err := jsonpb.UnmarshalString(string(evt.Data), info)

		if err == nil {
			// So here Will display message in a nice way...
			// First the header...
			header := info.Application

			// Set the occurence date.
			switch info.Level {
			case logpb.LogLevel_ERROR_MESSAGE, logpb.LogLevel_FATAL_MESSAGE:
				color.Error.Println(header)
			case logpb.LogLevel_DEBUG_MESSAGE, logpb.LogLevel_TRACE_MESSAGE, logpb.LogLevel_INFO_MESSAGE:
				color.Info.Println(header)
			case logpb.LogLevel_WARN_MESSAGE:
				color.Warn.Println(header)
			default:
				color.Comment.Println(header)
			}

			if len(info.Message) > 0 {
				fmt.Println("================= New Log Message ================")
				// Now I will process the message itself...
				msg := info.Message
				// if the message is grpc error I will parse it content and display it content...
				if strings.HasPrefix(msg, "rpc") {
					errorDescription := make(map[string]interface{})
					startIndex := strings.Index(msg, "{")
					endIndex := strings.Index(msg, "}")
					if startIndex >= 0 && endIndex > startIndex {
						jsonStr := msg[startIndex : endIndex+1]
						err := json.Unmarshal([]byte(jsonStr), &errorDescription)
						if err == nil {
							if errorDescription["FileLine"] != nil {
								fmt.Println(errorDescription["FileLine"])
							}
							if errorDescription["ErrorMsg"] != nil {
								color.Comment.Println(errorDescription["ErrorMsg"])
							}
						}
					}
				} else {
					if len(info.Line) > 0 {
						fmt.Println(info.Line)
					}
					color.Comment.Println(msg)
				}

				// I will also display the message in the system logger.
				if g.logger != nil && len(msg) > 0 {
					switch info.Level {
					case logpb.LogLevel_ERROR_MESSAGE, logpb.LogLevel_FATAL_MESSAGE:
						err = g.logger.Error(msg)
					case logpb.LogLevel_WARN_MESSAGE:
						err = g.logger.Warning(msg)
					default:
						err = g.logger.Info(msg)
					}

					if err != nil {
						fmt.Printf("warning: failed to log message: %v\n", err)
					}
				}
				fmt.Println("==================================================")
			}

		}

	}

}

/**
 * Listen for new connection.
 */
func (globule *Globule) Listen() error {

	fmt.Println("globular server is starting...")
	//autorestart.StartWatcher()

	var err error

	// if no certificates are specified I will try to get one from let's encrypts.
	// Start https server.
	if Utility.Exists(globule.creds+"/"+globule.Certificate) && globule.Protocol == "https" && len(globule.Certificate) > 0 {
		err = security.ValidateCertificateExpiration(globule.creds+"/"+globule.Certificate, globule.creds+"/server.pem")
		if err != nil {
			// here I will remove the expired certificates...
			err = Utility.RemoveDirContents(globule.creds)
			if err != nil {
				return err
			}
			globule.Certificate = ""
		}
	}

	if (!Utility.Exists(globule.creds+"/"+globule.Certificate) || len(globule.Certificate) == 0) && globule.Protocol == "https" {

		// Here is the command to be execute in order to ge the certificates.
		// ./lego --email="admin@globular.cloud" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/tls/server.csr run
		// I need to remove the gRPC certificate and recreate it.
		if Utility.Exists(globule.creds) {
			err = Utility.RemoveDirContents(globule.creds)
			if err != nil {
				return err
			}
		}

		// must be generated first.
		err = security.GenerateServicesCertificates(globule.CertPassword, globule.CertExpirationDelay, globule.getLocalDomain(), globule.creds, globule.Country, globule.State, globule.City, globule.Organization, globule.AlternateDomains)
		if err != nil {
			return err
		}

		// I will start the dns service if it is not already started.
		dns_config, err := config.GetServiceConfigurationById("dns.DnsService")
		if err == nil {
			// I will start the dns service...
			dns_config["State"] = "starting"
			name := dns_config["Name"].(string)
			dns_config["ProxyProcess"] = -1

			port := Utility.ToInt(dns_config["Port"])

			fmt.Println("start service ", name, " on port ", port, " and proxy port ", port)
			_, err := process.StartServiceProcess(dns_config, port)
			if err != nil {
				fmt.Println("fail to start service ", name, err)
			}
		} else {
			fmt.Println("fail to get dns service configuration with error ", err)
		}

		// Register that peer with the dns.
		err = globule.registerIpToDns()
		if err != nil {
			fmt.Println("fail to register ip to dns with error ", err)
			return err
		}

		err = globule.obtainCertificateForCsr()
		if err != nil {
			return err
		}

		err = globule.restart()
		if err != nil {
			fmt.Println("fail to restart globule with error ", err)
			return err
		}
	}

	go func() {
		address := "0.0.0.0" // or leave empty to bind all interfaces

		srv := &http.Server{
			Addr: fmt.Sprintf("%s:%d", address, globule.PortHTTP),
			// Defensive timeouts
			ReadHeaderTimeout: 5 * time.Second,   // <- fixes G112
			ReadTimeout:       30 * time.Second,  // full request (headers+body)
			WriteTimeout:      30 * time.Second,  // time to write response
			IdleTimeout:       120 * time.Second, // keep-alive connections
			MaxHeaderBytes:    1 << 20,           // 1 MiB; tune as needed
			// Handler:        yourMux,           // nil uses http.DefaultServeMux
		}
		globule.http_server = srv

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("HTTP server listen error:", err)
		}
	}()

	// Start the http server.
	if globule.Protocol == "https" {
		address := "0.0.0.0" // Utility.MyLocalIP()

		srv := &http.Server{
			Addr: fmt.Sprintf("%s:%d", address, globule.PortHTTPS),

			// --- Security & performance timeouts ---
			ReadHeaderTimeout: 5 * time.Second, // <-- fixes G112
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MiB; adjust as needed

			// Keep your handler (nil => DefaultServeMux)
			// Handler: yourMux,

			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: globule.getLocalDomain(), // optional for servers; safe to keep if you rely on it
				// NextProtos: []string{"h2", "http/1.1"}, // optional
			},
		}
		globule.https_server = srv

		go func() {
			certFile := filepath.Join(globule.creds, globule.Certificate)
			keyFile := filepath.Join(globule.creds, "server.pem")

			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Println("HTTPS server listen error:", err)
			}
		}()
	}

	//////////////////////////////////////////////////////////
	// Keep the server state in sync with configuration and
	// version.
	//////////////////////////////////////////////////////////

	// Keep globular up to date subscription.
	globule.watchForUpdate()

	return err
}

// ////////////////////// Resource Client ////////////////////////////////////////////
func GetServiceManagerClient(address string) (*service_manager_client.Services_Manager_Client, error) {

	Utility.RegisterFunction("NewServicesManagerService_Client", service_manager_client.NewServicesManagerService_Client)
	client, err := globular_client.GetClient(address, "services_manager.ServicesManagerService", "NewServicesManagerService_Client")
	if err != nil {
		return nil, err
	}
	return client.(*service_manager_client.Services_Manager_Client), nil

}

// ////////////////////// Resource Client ////////////////////////////////////////////
func getResourceClient(address string) (*resource_client.Resource_Client, error) {
	Utility.RegisterFunction("NewResourceService_Client", resource_client.NewResourceService_Client)
	client, err := globular_client.GetClient(address, "resource.ResourceService", "NewResourceService_Client")
	if err != nil {
		return nil, err
	}
	return client.(*resource_client.Resource_Client), nil
}

//////////////////////// RBAC function //////////////////////////////////////////////
/**
 * Get the rbac client.
 */
func GetRbacClient(address string) (*rbac_client.Rbac_Client, error) {
	Utility.RegisterFunction("NewRbacService_Client", rbac_client.NewRbacService_Client)
	client, err := globular_client.GetClient(address, "rbac.RbacService", "NewRbacService_Client")
	if err != nil {
		return nil, err
	}
	return client.(*rbac_client.Rbac_Client), nil
}

// Use rbac client here...
func (globule *Globule) addResourceOwner(path, resource_type, subject string, subjectType rbacpb.SubjectType) error {

	address, _ := config.GetAddress()
	rbac_client_, err := GetRbacClient(address)
	if err != nil {
		return err
	}
	return rbac_client_.AddResourceOwner(path, resource_type, subject, subjectType)
}

func (globule *Globule) validateAction(method string, subject string, subjectType rbacpb.SubjectType, infos []*rbacpb.ResourceInfos) (bool, bool, error) {
	address, _ := config.GetAddress()
	rbac_client_, err := GetRbacClient(address)
	if err != nil {
		return false, false, err
	}

	return rbac_client_.ValidateAction(method, subject, subjectType, infos)
}

func (globule *Globule) validateAccess(subject string, subjectType rbacpb.SubjectType, name string, path string) (bool, bool, error) {

	address, _ := config.GetAddress()
	rbac_client_, err := GetRbacClient(address)
	if err != nil {
		return false, false, err
	}
	hasAccess, hasAccessDenied, err := rbac_client_.ValidateAccess(subject, subjectType, name, path)
	return hasAccess, hasAccessDenied, err
}

func ValidateSubjectSpace(subject string, subjectType rbacpb.SubjectType, required_space uint64) (bool, error) {
	address, _ := config.GetAddress()
	rbac_client_, err := GetRbacClient(address)
	if err != nil {
		return false, err
	}
	hasSpace, err := rbac_client_.ValidateSubjectSpace(subject, subjectType, required_space)
	return hasSpace, err
}

func (globule *Globule) setActionResourcesPermissions(permissions map[string]interface{}) error {
	address, _ := config.GetAddress()
	rbac_client_, err := GetRbacClient(address)
	if err != nil {
		return err
	}
	return rbac_client_.SetActionResourcesPermissions(permissions)
}

// /////////////////// event service functions ////////////////////////////////////
func (globule *Globule) getEventClient() (*event_client.Event_Client, error) {

	Utility.RegisterFunction("NewEventService_Client", event_client.NewEventService_Client)
	address, _ := config.GetAddress()
	client, err := globular_client.GetClient(address, "event.EventService", "NewEventService_Client")
	if err != nil {
		fmt.Println("fail to get event client with error ", address, err)
		return nil, err
	}
	return client.(*event_client.Event_Client), nil
}

func (globule *Globule) publish(event string, data []byte) error {
	eventClient, err := globule.getEventClient()
	if err != nil {
		return err
	}
	err = eventClient.Publish(event, data)
	if err != nil {
		fmt.Println("fail to publish event", event, globule.Domain, "with error", err)
	}
	return err
}

func (globule *Globule) subscribe(evt string, listener func(evt *eventpb.Event)) error {

	eventClient, err := globule.getEventClient()
	if err != nil {
		log.Panicln("fail to get event client with error ", err)
		return err
	}

	err = eventClient.Subscribe(evt, globule.Name, listener)
	if err != nil {
		log.Panicln("fail to get event client with error ", err)
		return err
	}

	// register a listener...
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// Persistence service functions
// ////////////////////////////////////////////////////////////////////////////////////////////////
func getPersistenceClient(address string) (*persistence_client.Persistence_Client, error) {
	Utility.RegisterFunction("NewPersistenceService_Client", persistence_client.NewPersistenceService_Client)
	client, err := globular_client.GetClient(address, "persistence.PersistenceService", "NewPersistenceService_Client")
	if err != nil {
		return nil, err
	}
	return client.(*persistence_client.Persistence_Client), nil
}

// Create the application connections in the backend.
func (globule *Globule) createApplicationConnections() error {

	address, _ := config.GetAddress()
	resourceClient, err := getResourceClient(address)
	if err != nil {
		return err
	}

	applications, err := resourceClient.GetApplications("{}")
	if err != nil {
		return err
	}

	for i := 0; i < len(applications); i++ {
		err = globule.createApplicationConnection(applications[i])
		if err != nil {
			fmt.Println("fail to create application connection with error ", err)
		}
	}

	return nil
}

// Create the application connections in the backend.
func (globule *Globule) createApplicationConnection(app *resourcepb.Application) error {

	resourceServiceConfig, err := config.GetServiceConfigurationById("resource.ResourceService")
	if err != nil {
		return err
	}

	var storeType float64
	switch resourceServiceConfig["Backend_type"].(string) {
	case "SQL":
		storeType = 1.0
	case "MONGO":
		storeType = 0.0
	case "SCYLLA":
		storeType = 2.0
	}

	address, _ := config.GetAddress()
	persistenceClient, err := getPersistenceClient(address)

	if err != nil {
		return err
	}

	// Here I will create the connection to the backend.
	err = persistenceClient.CreateConnection(app.Name, app.Name+"_db", address, resourceServiceConfig["Backend_port"].(float64), storeType, resourceServiceConfig["Backend_user"].(string), resourceServiceConfig["Backend_password"].(string), 500, "", false)

	if err != nil {
		fmt.Printf("fail to create connection for %s with error %s", app.Name, err.Error())
	}
	return err
}
