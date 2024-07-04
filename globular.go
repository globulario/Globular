package main

import (
	controlplane "Globular/control-plane"
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
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
	"github.com/golang/protobuf/jsonpb"
	"github.com/gookit/color"
	"github.com/kardianos/service"
	//"github.com/slayer/autorestart"

	"github.com/txn2/txeh"

	// Interceptor for authentication, event, log...

	// Client services.
	"github.com/davecourtois/Utility"
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

/**
 * The web server.
 */
type Globule struct {
	// The share part of the service.
	Name string // the hostname.
	Mac  string // The Mac addresse

	// Where services can be found.
	ServicesRoot string

	// can be https or http.
	Protocol     string
	PortHttp     int    // The port of the http file server.
	PortHttps    int    // The secure port
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

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {

	fmt.Println("Create new Globular instance")
	// Here I will initialyse configuration.
	g := new(Globule)
	g.startTime = time.Now()
	g.exit_ = false
	g.exit = make(chan bool)
	g.Version = "1.0.0" // Automate version...
	g.Build = 0
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.IndexApplication = ""      // I will use the installer as defaut.
	g.PortHttp = 80              // The default http port 80 is almost already use by other http server...
	g.PortHttps = 443            // The default https port number
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
	g.AdminEmail = "root@globular.app"
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

	g.Path, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	g.Path = strings.ReplaceAll(g.Path, "\\", "/")

	// If no configuration exist I will create it before initialyse directories and start services.
	configPath := config.GetConfigDir() + "/config.json"
	if !Utility.Exists(configPath) {
		Utility.CreateDirIfNotExist(config.GetConfigDir())
		globule.config = config.GetConfigDir()
		err := globule.saveConfig()
		if err != nil {
			fmt.Println("fail to save local configuration with error", err)
		} else {
			fmt.Println("local configuration was saved ", Utility.Exists(configPath))
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
	p["portHttp"] = globule.PortHttp
	p["portHttps"] = globule.PortHttps

	jsonStr, _ := json.Marshal(&p)

	// set services configuration values
	globule.publish("stop_peer_evt", jsonStr)

	// give time to stop peer evt to be publish
	time.Sleep(500 * time.Millisecond)

	// Close all services.
	globule.stopServices()

	// reset firewall rules.
	resetRules()

	globule.saveConfig()

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
			fmt.Println("fail to register admin account sa", err)
			return err
		}

	
		// Admin is created
		err = globule.createAdminRole()
		if err != nil {
			if !strings.Contains(err.Error(), "already exist") {
				return err
			}
		}

		// Set admin role to that account.
		resource_client_.AddAccountRole("sa", "admin")

		path := config.GetDataDir() + "/files/users/sa@" + globule.Domain
		if !Utility.Exists(path) {
			err := Utility.CreateDirIfNotExist(path)
			if err == nil {
				globule.addResourceOwner("/users/sa@"+globule.Domain, "file", "sa@"+globule.Domain, rbacpb.SubjectType_ACCOUNT)
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
						os.Rename(config.GetDataDir()+"/files/users/"+accounts[i].Id+"@"+accounts[i].Domain, config.GetDataDir()+"/files/users/"+accounts[i].Id+"@"+globule.Domain)

						// I will update the account domain
						accounts[i].Domain = globule.Domain
						err := resource_client_.SetAccount(token, accounts[i])
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

/**
 * The admin group contain all action...
 */
func (globule *Globule) createAdminRole() error {
	address, _ := config.GetAddress()
	resource_client_, err := getResourceClient(address)
	if err != nil {
		return err
	}

	mac := strings.ReplaceAll(globule.Mac, ":", "_")

	token, err := os.ReadFile(config.GetConfigDir() + "/tokens/" + mac + "_token")
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

	// Create the admin account.
	err = resource_client_.CreateRole(string(token), "admin", "admin", actions)
	if err != nil {
		fmt.Println("fail to create admin role", err)
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
		s["PublisherId"] = services[i]["PublisherId"]
		s["State"] = services[i]["State"]
		s["TLS"] = services[i]["TLS"]
		s["Dependencies"] = services[i]["Dependencies"]
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
		port := p.PortHttp
		if p.Protocol == "https" {
			port = p.PortHttps
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
	if config["PortHttp"] != nil {
		globule.PortHttp = Utility.ToInt(config["PortHttp"])
	}

	// Set the port.
	if config["PortHttps"] != nil {
		globule.PortHttps = Utility.ToInt(config["PortHttps"])
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

	// Set the ns.
	if config["NS"] != nil {

		globule.NS = config["NS"].([]interface{})

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
		globule.restart()

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

	Utility.CreateDirIfNotExist(globule.config)
	configPath := globule.config + "/config.json"

	err = os.WriteFile(configPath, []byte(jsonStr), 0644)
	if err != nil {
		fmt.Println("fail to save configuration with error: ", err)
		return err
	}

	fmt.Println("globular configuration was save at ", configPath)

	return nil
}

/**
 * Return the admin email.
 */
func (globule *Globule) GetEmail() string {
	return globule.AdminEmail
}

/**
 * Use the time of registration... Nil other wise.
 */
func (globule *Globule) GetRegistration() *registration.Resource {
	return globule.registration
}

/**
 * That function must be use to generate public
 */

/**
 * I will reuse the client public key here as key instead of generate another key
 * and manage it...
 */
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

/** That interface withe the let's encrypt DNS chanlenge. **/
type DNSProviderGlobularDNS struct {
	apiAuthToken string
}

func NewDNSProviderGlobularDNS(apiAuthToken string) (*DNSProviderGlobularDNS, error) {
	return &DNSProviderGlobularDNS{apiAuthToken: apiAuthToken}, nil
}

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
		provider := http01.NewProviderServer("", strconv.Itoa(globule.PortHttp))
		err = client.Challenge.SetHTTP01Provider(provider)
		if err != nil {
			fmt.Println("fail to set http provider with error: ", err)
			return err
		}
	}

	if err != nil {
		fmt.Println("1031 fail to create new client with error: ", err)
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
	os.WriteFile(globule.creds+"/"+globule.Certificate, resource.Certificate, 0400)
	os.WriteFile(globule.creds+"/"+globule.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	return globule.saveConfig()
}

func (globule *Globule) signCertificate(client_csr string) (string, error) {

	// first of all I will save the incomming file into a temporary file...
	client_csr_path := os.TempDir() + "/" + Utility.RandomUUID()
	err := os.WriteFile(client_csr_path, []byte(client_csr), 0644)
	if err != nil {
		return "", err

	}

	client_crt_path := os.TempDir() + "/" + Utility.RandomUUID()

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+globule.CertPassword)
	args = append(args, "-days")
	args = append(args, Utility.ToString(globule.CertExpirationDelay))
	args = append(args, "-in")
	args = append(args, client_csr_path)
	args = append(args, "-CA")
	args = append(args, globule.creds+"/ca.crt") // use certificate
	args = append(args, "-CAkey")
	args = append(args, globule.creds+"/ca.key") // and private key to sign the incommin csr
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, client_crt_path)
	args = append(args, "-extfile")
	args = append(args, globule.creds+"/san.conf")
	args = append(args, "-extensions")
	args = append(args, "v3_req")
	cmd_ := exec.Command(cmd, args...)
	cmd_.Dir = os.TempDir()
	err = cmd_.Run()
	if err != nil {
		return "", err
	}

	// I will read back the crt file.
	client_crt, err := os.ReadFile(client_crt_path)

	// remove the tow temporary files.
	defer os.Remove(client_crt_path)
	defer os.Remove(client_csr_path)

	if err != nil {
		return "", err
	}

	return string(client_crt), nil
}

/**
 * Initialize the server directories config, data, webroot...
 */
func (globule *Globule) initDirectories() error {

	fmt.Println("init directories")

	// initilayse configurations...
	// it must be call here in order to initialyse a sync map...
	config.GetServicesConfigurations()

	// DNS info.
	globule.DNS = globule.getLocalDomain() // The dns server.

	// The name server.
	globule.NS = make([]interface{}, 0)

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

	err := Utility.CreateDirIfNotExist(globule.data)
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
	Utility.CreateDirIfNotExist(globule.creds)

	// Files directorie that contain user's directories and application's directory
	globule.users = globule.data + "/files/users"
	Utility.CreateDirIfNotExist(globule.users)

	// Contain the application directory.
	globule.applications = globule.data + "/files/applications"
	Utility.CreateDirIfNotExist(globule.applications)

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
			err := os.WriteFile(globule.config+"/config.json", []byte(jsonStr), 0644)
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
		globule.Mac, _ = config.GetMacAddress()
	}

	// save config...
	globule.saveConfig()

	if !Utility.Exists(globule.webRoot + "/index.html") {

		// in that case I will create a new index.html file.
		os.WriteFile(globule.webRoot+"/index.html", []byte(
			`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
		<html lang="en">
		
			<head>
				<meta http-equiv="content-type" content="text/html; charset=utf-8">
				<title>Title Goes Here</title>
			</head>
		
			<body>
				<p>Welcome to Globular `+globule.Version+`</p>
			</body>
		
		</html>`), 0644)
	}

	return nil
}

func (globule *Globule) refreshLocalToken() error {

	// set the local token.
	return security.SetLocalToken(globule.Mac, globule.Domain, "sa", "sa", globule.AdminEmail, globule.SessionTimeout)
}

// Enable port from window firewall
func enablePorts(ruleName, portsRange string) error {

	if runtime.GOOS == "windows" {
		deleteRule(ruleName)

		inboundRule_ := fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=in action=allow protocol=TCP localport=%s`, ruleName, portsRange)
		fmt.Println(inboundRule_)

		// netsh advfirewall firewall add rule name="Globular-Services" dir=in action=allow protocol=TCP localport=10000-10100
		inboundRule := exec.Command("cmd", "/C", inboundRule_)
		inboundRule.Dir = os.TempDir()
		err := inboundRule.Run()
		if err != nil {
			fmt.Println("fail to add inbound rule: ", ruleName, "with error: ", err)
			return nil
		}

		outboundRule_ := fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=out action=allow protocol=TCP localport=%s`, ruleName, portsRange)
		fmt.Println(outboundRule_)
		outboundRule := exec.Command("cmd", "/C", outboundRule_)
		outboundRule.Dir = os.TempDir()
		err = outboundRule.Run()
		if err != nil {
			fmt.Println("fail to add outbound rule: ", ruleName, "with error: ", err)
			return nil
		}

	}

	return nil
}

func enableProgramFwMgr(name, appname string) error {

	if runtime.GOOS == "windows" {
		fmt.Println("enable program: ", name, appname)
		// netsh advfirewall firewall add rule name="MongoDB Database Server" dir=in action=allow program="C:\Program Files\Globular\dependencies\mongodb-win32-x86_64-windows-5.0.5\bin\mongod.exe" enable=yes
		appname = strings.ReplaceAll(appname, "/", "\\")
		inboundRule := exec.Command("cmd", "/C", fmt.Sprintf(`netsh advfirewall firewall add rule name="%s" dir=in action=allow program="%s" enable=yes`, name, appname))
		inboundRule.Dir = os.TempDir()
		err := inboundRule.Run()
		if err != nil {
			return err
		}

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
	deleteRule("alertmanager")
	deleteRule("mongod")
	deleteRule("prometheus")
	deleteRule("torrent")
	deleteRule("yt-dlp")

	// other rules.
	deleteRule("Globular")
	deleteRule("Globular-http")
	deleteRule("Globular-https")
	deleteRule("Globular-Services")

	for i := 0; i < len(services); i++ {
		// Create the service process.
		deleteRule(services[i]["Name"].(string) + "-" + services[i]["Id"].(string))
	}

	return nil
}

func resetSystemPath() error {

	if runtime.GOOS == "windows" {

		Utility.UnsetWindowsEnvironmentVariable("OPENSSL_CONF")

		systemPath, err := Utility.GetWindowsEnvironmentVariable("Path")

		if err != nil {
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
		execs := Utility.GetFilePathsByExtension(config.GetRootDir()+"/dependencies", ".exe")
		for i := 0; i < len(execs); i++ {
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
	if runtime.GOOS == "windows" {
		// remove previous rules...
		resetRules()

		ex, _ := os.Executable()
		// set globular firewall run...
		enableProgramFwMgr("Globular", ex)

		// Enable ports
		enablePorts("Globular-Services", globule.PortsRange)
		enablePorts("Globular-http", strconv.Itoa(globule.PortHttp))
		enablePorts("Globular-https", strconv.Itoa(globule.PortHttps))

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

		// set rules for services contain in dependencies folder.
		execs := Utility.GetFilePathsByExtension(config.GetRootDir()+"/dependencies", ".exe")
		for i := 0; i < len(execs); i++ {
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
				enableProgramFwMgr("yt-dlp", exec)
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

		for i := 0; i < len(services); i++ {
			service := services[i]
			id := service["Id"].(string)
			path := service["Path"].(string)
			name := service["Name"].(string)

			// Create the service process.
			enableProgramFwMgr(name+"-"+id, path)
		}

		// Openssl conf require...
		path := strings.ReplaceAll(config.GetRootDir(), "/", "\\") + `\dependencies\openssl.cnf`
		
		if Utility.Exists(`C:\Program Files\Globular\dependencies\openssl.cnf`) {
			Utility.SetWindowsEnvironmentVariable("OPENSSL_CONF", path)
		} else {
			fmt.Println("Open SSL configuration file ", path, "not found. Require to create environnement variable OPENSSL_CONF.")
		}
		err = Utility.SetWindowsEnvironmentVariable("Path", strings.ReplaceAll(systemPath, "/", "\\"))

		return err
	} else if runtime.GOOS == "darwin" {
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

					os.WriteFile("/Library/LaunchDaemons/Globular.plist", []byte(config_), 0644)
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
			globule.refreshLocalToken()
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
	security.GeneratePeerKeys(globule.Mac)

	// This is the local token...
	err := globule.refreshLocalToken()
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
			proxyPort := port + 1

			fmt.Println("try to start service ", name, " on port ", port, " and proxy port ", proxyPort)
			_, err := process.StartServiceProcess(service, port)
			if err != nil {
				fmt.Println("fail to start service ", name, err)
			}
		}
	}

	// Here I will listen for logger event...
	go func() {

		// wait 15 second before register resource permissions and subscribe to log events.
		time.Sleep(15 * time.Second)

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
		globule.subscribe("new_log_evt", logListener(globule))

		// So here I will authenticate the root if the password is "adminadmin" that will
		// reset the password in the backend if it was manualy set in the config file.
		/*config_, err := config.GetConfig(true)
		if err == nil {
			if config_["RootPassword"].(string) == "adminadmin" {

				address, _ := config.GetAddress()

				// Authenticate the user in order to get the token
				authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
				if err != nil {
					log.Println("fail to access resource service at "+address+" with error ", err)
					return
				}

				log.Println("authenticate user ", "sa", " at adress ", address)
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

	globule.registerIpToDns()

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
	if state == 0 {
		p.State = resourcepb.PeerApprovalState_PEER_PENDING
	} else if state == 1 {
		p.State = resourcepb.PeerApprovalState_PEER_ACCETEP
	} else if state == 2 {
		p.State = resourcepb.PeerApprovalState_PEER_REJECTED
	}

	p.PortHttp = int32(Utility.ToInt(p_["portHttp"]))
	p.PortHttps = int32(Utility.ToInt(p_["portHttps"]))
	if p_["actions"] != nil {
		p.Actions = make([]string, len(p_["actions"].([]interface{})))

		for i := 0; i < len(p_["actions"].([]interface{})); i++ {
			p.Actions[i] = p_["actions"].([]interface{})[i].(string)
		}
	} else {
		p.Actions = make([]string, 0)
	}

	globule.peers.Store(p.Mac, p)
	globule.savePeers()

	// Here I will try to set the peer ip...

	// set the peer ip in the /etc/hosts file.
	if Utility.MyIP() == p.ExternalIpAddress {
		globule.setHost(p.LocalIpAddress, p.Hostname+"."+p.Domain)
	}

}

func deletePeersEvent(evt *eventpb.Event) {
	globule.peers.Delete(string(evt.Data))
	globule.savePeers()
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
		globule.setHost(p.LocalIpAddress, address)
	}

	if p.Protocol == "https" {
		address += ":" + Utility.ToString(p.PortHttps)
	} else {
		address += ":" + Utility.ToString(p.PortHttp)
	}

	// Here I will get the peer public key if not already exist.
	if !Utility.Exists(globule.config + "/keys/" + strings.ReplaceAll(p.Mac, ":", "_") + "_public") {

		// get the peer public key.
		rqst := p.Protocol + "://" + address + "/public_key"
		resp, err := http.Get(rqst)
		if err == nil {
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err == nil {
				// save the peer public key.
				os.WriteFile(globule.config+"/keys/"+strings.ReplaceAll(p.Mac, ":", "_")+"_public", body, 0644)
			}
		} else {
			fmt.Println("fail to get peer public key with error: ", err)
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
		peer_.PortHttp = int32(globule.PortHttp)
		peer_.PortHttps = int32(globule.PortHttps)
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
	for i := 0; i < len(peers); i++ {
		p := peers[i]

		// Set existing value...
		globule.peers.Store(p.Mac, p)

		// Try to update with updated infos...
		go func(p *resourcepb.Peer) {
			err := globule.initPeer(p)
			if err != nil {
				globule.peers.Delete(p.Mac) // remove the peer from the list.
				globule.savePeers()         // save the peers list.
			}
		}(p)
	}

	// Subscribe to new peers event...
	globule.subscribe("update_peers_evt", updatePeersEvent)
	globule.subscribe("delete_peer_evt", deletePeersEvent)

	// Now I will set the local peer info...
	globule.savePeers()

	return nil
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

	if err == nil {
		for i := 0; i < len(services_configs); i++ {
			pid := Utility.ToInt(services_configs[i]["Process"])
			services_configs[i]["State"] = "killed"
			services_configs[i]["ProxyProcess"] = -1

			// save config...
			err := config.SaveServiceConfiguration(services_configs[i])
			if err == nil {
				if pid > 0 {
					process, err := os.FindProcess(pid)
					if err == nil {
						err = process.Signal(syscall.SIGTERM) // make the process stop gracefully.
						if err != nil {
							fmt.Println("Error sending signal:", err)
						}
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
	globule.registerAdminAccount()

	url := globule.Protocol + "://" + globule.getAddress()
	if globule.Protocol == "https" {
		if globule.PortHttps != 443 {
			url += ":" + Utility.ToString(globule.PortHttps)
		}
	} else if globule.Protocol == "http" {
		if globule.PortHttp != 80 {
			url += ":" + Utility.ToString(globule.PortHttp)
		}
	}

	elapsed := time.Since(globule.startTime)

	fmt.Println("globular version " + globule.Version + " build " + Utility.ToString(globule.Build) + " listen at address " + url)
	fmt.Printf("startup took %s\n", elapsed)

	// create applications connections
	err := globule.createApplicationConnections()
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
				err := os.WriteFile(configPath, []byte(config_), 0644)
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
		proxy := uint32((Utility.ToInt(service["Proxy"])))

		if _, ok := proxies[proxy]; !ok {
			snapshot := controlplane.Snapshot{

				ClusterName:  strings.ReplaceAll(service["Name"].(string), ".", "_") + "_cluster",
				RouteName:    strings.ReplaceAll(service["Name"].(string), ".", "_") + "_route",
				ListenerName: strings.ReplaceAll(service["Name"].(string), ".", "_") + "_listener",
				ListenerPort: proxy,
				ListenerHost: "0.0.0.0", // local address.

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

		for i := 0; i < len(pids); i++ {
			if pids[i] != os.Getpid() {
				Utility.TerminateProcess(pids[i], 0)
			}
		}
	}

	// Initialyse directories.
	globule.initDirectories()
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
	globule.startServices()

	// Start process monitoring with prometheus.
	process.StartProcessMonitoring(globule.Protocol, globule.PortHttp, globule.exit)

	// Watch config.
	globule.watchConfig()

	// Set the fmt information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// First of all i will set the local host found...
	go func() {
		hosts := Utility.GetHostnameIPMap(globule.LocalIpAddress)
		for k, v := range hosts {
			globule.setHost(k, v)
		}

		// Try with ip address...
		ips, err := Utility.ScanIPs()
		if err == nil {
			for i := 0; i < len(ips); i++ {
				config_, err := config.GetRemoteConfig(ips[i], 80)
				if err == nil {
					hostname := config_["Name"].(string)
					if config_["Domain"] != nil {
						if config_["Domain"].(string) != "localhost" {
							hostname += "." + config_["Domain"].(string)
						}
					}

					globule.setHost(ips[i], hostname)

				}
			}
		}
	}()

	// Init peers
	go globule.initPeers()

	// Now I will initialize the control plane.
	go globule.initControlPlane()

	// Start envoy proxy.
	startEnvoyProxy()

	p := make(map[string]interface{})

	p["address"] = globule.getAddress()
	p["domain"], _ = config.GetDomain()
	p["hostname"] = globule.Name
	p["mac"] = globule.Mac
	p["portHttp"] = globule.PortHttp
	p["portHttps"] = globule.PortHttps

	jsonStr, _ := json.Marshal(&p)

	// set services configuration values
	globule.publish("start_peer_evt", jsonStr)

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

	defer watcher.Close()
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
					globule.setConfig(config)
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
 * Set the ip for a given domain or sub-domain
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
		}

		// I will publish the private ip address only
		_, err = dns_client_.SetA(token, globule.getLocalDomain(), config.GetLocalIP(), 60)
		if err != nil {
			fmt.Println("fail to set A record for domain ", globule.getLocalDomain(), " with error ", err)
		} else {
			fmt.Println("set A record for domain ", globule.getLocalDomain(), " with success")
		}

		// try to set the ipv6 address...
		ipv6, err := Utility.MyIPv6()
		if err == nil {
			_, err = dns_client_.SetAAAA(token, globule.getLocalDomain(), ipv6, 60)
			if err != nil {
				fmt.Println("fail to set AAAA  domain ", globule.getLocalDomain(), " with error ", err)
			} else {
				fmt.Println("set AAAA record for domain ", globule.getLocalDomain(), " with success")
			}
		}

		// I will set alternate domain only if the globule is the master.
		if globule.DNS == globule.getLocalDomain() {

			// Here I will set the A record for the globular domain.
			dns_client_.RemoveA(token, globule.getLocalDomain())

			_, err = dns_client_.SetA(token, globule.getLocalDomain(), Utility.MyIP(), 60)
			if err != nil {
				fmt.Println("fail to set A record for alternate domain ", globule.getLocalDomain(), " with error ", err)
			} else {
				fmt.Println("set A record for alternate domain ", globule.getLocalDomain(), Utility.MyIP(), " with success")
			}

			for j := 0; j < len(globule.AlternateDomains); j++ {

				// Here I will set the A record for the alternate domain.
				alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
				_, err = dns_client_.SetA(token, alternateDomain, Utility.MyIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					continue
				} else {
					fmt.Println("set A record for alternate domain ", alternateDomain, " with success")
				}

				_, err = dns_client_.SetA(token, alternateDomain, config.GetLocalIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					continue
				} else {
					fmt.Println("set A record for alternate domain ", alternateDomain, " with success")
				}

				_, err = dns_client_.SetA(token, alternateDomain, Utility.MyIP(), 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					continue
				} else {
					fmt.Println("set A record for alternate domain ", alternateDomain, Utility.MyIP(), " with success")
				}

				_, err = dns_client_.SetAAAA(token, alternateDomain, ipv6, 60)
				if err != nil {
					fmt.Println("fail to set A record for alternate domain ", alternateDomain, " with error ", err)
					continue
				} else {
					fmt.Println("set AAAA record for alternate domain ", alternateDomain, " with success")
				}
			}
		}

		// if the NS server are local I will set the local ip address.
		for j := 0; j < len(globule.NS); j++ {

			// Set it nameservers are part of the domain.
			ns := globule.NS[j].(string)

			//if strings.HasSuffix(ns, globule.Domain) {
			_, err = dns_client_.SetA(token, globule.NS[j].(string), Utility.MyIP(), 60)
			if err != nil {
				fmt.Println("fail to set A record for NS server ", ns, " with error ", err)
				continue
			} else {
				fmt.Println("set A record for NS server ", ns, " with success")
			}

			_, err = dns_client_.SetAAAA(token, globule.NS[j].(string), ipv6, 60)
			if err != nil {
				fmt.Println("fail to set AAAA record for domain ", ns, " with error ", err)
			} else {
				fmt.Println("set AAAA record for domain ", ns, " with success")
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
					continue
				} else {
					fmt.Println("set NS record for alternate domain ", alternateDomain, ns, " with success")
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
					continue
				} else {
					fmt.Println("set SOA record for alternate domain ", alternateDomain, ns, " with success")
				}

			}
		}

		for j := 0; j < len(globule.AlternateDomains); j++ {
			alternateDomain := strings.TrimPrefix(globule.AlternateDomains[j].(string), "*.")
			// Now I will set the CAA record.
			err = dns_client_.SetCaa(token, alternateDomain+".", 0, "issue", "letsencrypt.org", 60)
			if err != nil {
				fmt.Println("fail to set CAA record for alternate domain ", globule.AlternateDomains[j], " with error ", err)
				continue
			} else {
				fmt.Println("set CAA record for alternate domain ", globule.AlternateDomains[j], " with success")
			}
		}

		// save the file to /etc/resolv.conf
		if Utility.Exists("/etc/resolv.conf") {
			Utility.MoveFile("/etc/resolv.conf", "/etc/resolv.conf_")
			Utility.WriteStringToFile("/etc/resolv.conf", resolv_conf)
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

	resp, err = http.Get(url)

	if err != nil || resp.StatusCode != http.StatusCreated {
		url = "https://" + address + ":" + Utility.ToString(port) + "/checksum"
		resp, err = http.Get(url)
		if err != nil {
			return "", err
		}
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
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

						err := update_globular_from(globule, discovery, globule.getLocalDomain(), "sa", globule.RootPassword, runtime.GOOS+":"+runtime.GOARCH)
						if err != nil {
							fmt.Println("fail to update globular from " + discovery + " with error " + err.Error())
						}

						// Here I will restart the server.
						globule.restart()

					}
				} else {
					fmt.Println("fail to get checksum from : ", discovery, " error: ", err)
				}

			}

			services, err := config.GetServicesConfigurations()
			if err == nil {
				// get the resource client
				for i := 0; i < len(services); i++ {
					s := services[i]
					values := strings.Split(s["PublisherId"].(string), "@")
					if len(values) == 2 {
						resource_client_, err := getResourceClient(values[1])
						if err == nil {
							// test if the service need's to be updated, test if the part is part of installed instance and not developement environement.
							if s["KeepUpToDate"].(bool) && strings.Contains(s["Path"].(string), "/globular/services/") {
								// Here I will get the last version of the package...
								descriptor, err := resource_client_.GetPackageDescriptor(s["Id"].(string), s["PublisherId"].(string), "")

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
												if servicesManager.UninstallService(token, s["Domain"].(string), s["PublisherId"].(string), s["Id"].(string), s["Version"].(string)) == nil {
													servicesManager.InstallService(token, s["Domain"].(string), s["PublisherId"].(string), s["Id"].(string), descriptor.Version)
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
			// So here Will display message

			header := info.Application

			// Set the occurence date.
			if info.Level == logpb.LogLevel_ERROR_MESSAGE || info.Level == logpb.LogLevel_FATAL_MESSAGE {
				color.Error.Println(header)
			} else if info.Level == logpb.LogLevel_DEBUG_MESSAGE || info.Level == logpb.LogLevel_TRACE_MESSAGE || info.Level == logpb.LogLevel_INFO_MESSAGE {
				color.Info.Println(header)
			} else if info.Level == logpb.LogLevel_WARN_MESSAGE {
				color.Warn.Println(header)
			} else {
				color.Comment.Println(header)
			}

			if len(info.Message) > 0 {
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
					if info.Level == logpb.LogLevel_ERROR_MESSAGE || info.Level == logpb.LogLevel_FATAL_MESSAGE {
						g.logger.Error(msg)
					} else if info.Level == logpb.LogLevel_WARN_MESSAGE {
						g.logger.Warning(msg)
					} else {
						g.logger.Info(msg)
					}
				}
			}
		}
	}
}

/**
 * Listen for new connection.
 */
func (globule *Globule) Listen() error {

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

		fmt.Println("generate certificates...")
		// Here is the command to be execute in order to ge the certificates.
		// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/tls/server.csr run
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
			proxyPort := Utility.ToInt(dns_config["Proxy"])

			fmt.Println("start service ", name, " on port ", port, " and proxy port ", proxyPort)
			_, err := process.StartServiceProcess(dns_config, port)
			if err != nil {
				fmt.Println("fail to start service ", name, err)
			}
		}

		// Register that peer with the dns.
		err = globule.registerIpToDns()
		if err != nil {
			return err
		}

		err = globule.obtainCertificateForCsr()
		if err != nil {
			return err
		}

		globule.restart()
	}

	// Must be started before other services.
	go func() {

		address := "0.0.0.0" // Utility.MyLocalIP()

		// local - non secure connection.
		globule.http_server = &http.Server{
			Addr: address + ":" + strconv.Itoa(globule.PortHttp),
		}

		err = globule.http_server.ListenAndServe()
		if err != nil {
			fmt.Println("fail to listen with err ", err)
		}
	}()

	// Start the http server.
	if globule.Protocol == "https" {
		address := "0.0.0.0" // Utility.MyLocalIP()

		globule.https_server = &http.Server{
			Addr: address + ":" + strconv.Itoa(globule.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: globule.getLocalDomain(),
			},
		}

		// get the value from the configuration files.
		go func() {
			err = globule.https_server.ListenAndServeTLS(globule.creds+"/"+globule.Certificate, globule.creds+"/server.pem")
			if err != nil {
				fmt.Println("fail to listen with err ", err)
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
		return err
	}

	err = eventClient.Subscribe(evt, globule.Name, listener)
	if err != nil {
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
	if resourceServiceConfig["Backend_type"].(string) == "SQL" {
		storeType = 1.0
	} else if resourceServiceConfig["Backend_type"].(string) == "MONGO" {
		storeType = 0.0
	} else if resourceServiceConfig["Backend_type"].(string) == "SCYLLA" {
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
