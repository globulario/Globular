package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	ps "github.com/mitchellh/go-ps"

	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/file/file_client"
	"github.com/globulario/services/golang/interceptors"
	"github.com/globulario/services/golang/ldap/ldap_client"
	"github.com/globulario/services/golang/persistence/persistence_client"
	"github.com/struCoder/pidusage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/prometheus/client_golang/prometheus"

	// Interceptor for authentication, event, log...

	// Client services.
	"crypto"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/storage/storage_store"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"

	globular "github.com/globulario/services/golang/globular_service"
	"github.com/globulario/services/golang/lb/lbpb"
	"github.com/globulario/services/golang/persistence/persistence_store"
	"github.com/globulario/services/golang/security"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// OAuth2
	/*
		errors_ "github.com/go-oauth2/oauth2/errors"
		"github.com/go-oauth2/oauth2/generates"
		"github.com/go-oauth2/oauth2/manage"
		"github.com/go-oauth2/oauth2/models"
		"github.com/go-oauth2/oauth2/server"
		"github.com/go-oauth2/oauth2/store"
	*/)

// Global variable.
var (
	globule *Globule
)

type ExternalApplication struct {
	Id   string
	Path string
	Args []string

	// keep the actual srvice command here.
	srv *exec.Cmd
}

/**
 * The web server.
 */
type Globule struct {
	// The share part of the service.
	Name string // The service name

	// Globualr specifics ports.
	PortHttp   int    // The port of the http file server.
	PortHttps  int    // The secure port
	AdminEmail string // The admin email
	PortsRange string // The range of port to be use for the service. ex 10000-10200
	DbIpV4     string // The address of the database ex 0.0.0.0:27017

	// can be https or http.
	Protocol string // The protocol of the service.

	// Use to store services informations
	Services map[string]interface{}

	// The list of install services.
	services *sync.Map

	LdapSyncInfos map[string]interface{} // Contain LdapSyncInfos...

	// List of application need to be start by the server.
	ExternalApplications map[string]ExternalApplication

	Domain           string        // The principale domain
	AlternateDomains []interface{} // Alternate domain for multiple domains
	IndexApplication string        // If defined It will be use as the entry point where not application path was given in the url.

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
	Version        string
	Build          int64
	Platform       string
	SessionTimeout time.Duration

	// There's are Directory
	// The user's directory
	UsersDirectory string
	// The application directory
	ApplicationDirectory string

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep service up to date.

	// If it set to true all services will be updated automaticaly.
	KeepAllServicesUpToDate bool

	// If set to true services will be keept alive by default.
	KeepAllServicesAlive bool

	// if set to true the server will be updated automaticaly.
	KeepUpToDate bool

	// DNS stuff.
	DNS []interface{} // Domain name server use to located the server.

	DnsUpdateIpInfos []interface{} // The internet provader SetA info to keep ip up to date.

	discorveriesEventHub map[string]*event_client.Event_Client

	// The list of method supported by this server.
	methods []string

	// The prometheus logging informations.
	methodsCounterLog *prometheus.CounterVec

	// Monitor the cpu usage of process.
	servicesCpuUsage    *prometheus.GaugeVec
	servicesMemoryUsage *prometheus.GaugeVec

	// Directories.
	path         string // The path of the exec...
	webRoot      string // The root of the http file server.
	data         string // the data directory
	creds        string // tls certificates
	config       string // configuration directory
	users        string // the users files directory
	applications string // The applications

	// Log store.
	logs *storage_store.LevelDB_store

	// RBAC store.
	permissions *storage_store.LevelDB_store

	// Create the JWT key used to create the signature
	jwtKey       []byte // This is the client secret.
	RootPassword string

	// local store.
	store persistence_store.Store

	// client reference...
	persistence_client_ *persistence_client.Persistence_Client
	ldap_client_        *ldap_client.LDAP_Client
	event_client_       *event_client.Event_Client
	file_clients_       *sync.Map

	// ACME protocol registration
	registration *registration.Resource

	// load balancing action channel.
	lb_load_info_channel             chan *lbpb.LoadInfo
	lb_remove_candidate_info_channel chan *lbpb.ServerInfo
	lb_get_candidates_info_channel   chan map[string]interface{}
	lb_stop_channel                  chan bool

	// exit channel.
	exit            chan struct{}
	exit_           bool
	inernalServices []*grpc.Server

	// keep track of services updates from external sources.
	subscribers map[string]map[string][]string

	// The http server
	http_server  *http.Server
	https_server *http.Server

	// temporary array to be use to get next available port.
	portsInUse []int
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {

	// Here I will initialyse configuration.
	g := new(Globule)

	g.Version = "1.0.0" // Automate version...
	g.Build = 0
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.RootPassword = "adminadmin"
	g.IndexApplication = "globular_installer" // I will use the installer as defaut.

	g.PortHttp = 80   // The default http port
	g.PortHttps = 443 // The default https port number
	execPath := Utility.GetExecName(os.Args[0])
	g.Name = strings.Replace(execPath, ".exe", "", -1)

	// Set the default checksum...
	g.Protocol = "http"
	g.Domain = "localhost"

	// Set default values.
	g.PortsRange = "10000-10100"
	g.DbIpV4 = "0.0.0.0:27017"

	// set default values.
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"

	// keep up to date by default.
	g.KeepAllServicesUpToDate = true
	g.KeepAllServicesAlive = true

	// No update globular by default.
	g.KeepUpToDate = false

	// Keep the
	g.Services = make(map[string]interface{})
	g.services = new(sync.Map)
	g.inernalServices = make([]*grpc.Server, 0)

	// Contain the list of ldap syncronization info.
	g.LdapSyncInfos = make(map[string]interface{})

	// Configuration must be reachable before services initialysation

	// Promometheus metrics for services.
	http.Handle("/metrics", promhttp.Handler())

	// Keep in global var to by http handlers.
	globule = g

	// Set the list of http handler.

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

	// The checksum handler.
	http.HandleFunc("/checksum", getChecksumHanldler)

	// Handle the get ca certificate function
	http.HandleFunc("/get_ca_certificate", getCaCertificateHanldler)

	// Return the san server configuration.
	http.HandleFunc("/get_san_conf", getSanConfigurationHandler)

	// Handle the signing certificate function.
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

	g.path, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	if Utility.Exists(g.path+"/bin/grpcwebproxy") || Utility.Exists(g.path+"/bin/grpcwebproxy.exe") {
		// TODO test restart with initDirectories
		g.initDirectories()
	}

	return g
}

/**
 * Send a application notification.
 * That function will send notification to all connected user of that application.
 */
func (globule *Globule) sendApplicationNotification(application string, message string) error {

	// That service made user of persistence service.
	p, err := globule.getPersistenceStore()
	if err != nil {
		return err
	}

	/** The notification object. */
	notification := make(map[string]interface{})
	id := time.Now().Unix()
	notification["_id"] = id
	notification["_type"] = 1
	notification["_text"] = message
	notification["_recipient"] = application
	notification["_date"] = id

	// Now I will retreive the application icon...
	data, err := p.FindOne(context.Background(), "local_resource", "local_resource", "Applications", `{"_id":"`+application+`"}`, `[{"Projection":{"icon":1}}]`)
	if err != nil {
		return err
	}

	jsonStr, err := Utility.ToJson(data)
	if err != nil {
		return err
	}
	notification["_sender"] = jsonStr

	_, err = p.InsertOne(context.Background(), "local_resource", application+"_db", "Notifications", notification, "")
	if err != nil {
		return err
	}

	jsonStr, err = Utility.ToJson(notification)
	if err != nil {
		return err
	}

	eventHub, err := globule.getEventHub()
	if err != nil {
		return err
	}

	return eventHub.Publish(application+"_notification_event", []byte(jsonStr))
}

// Little shortcut to get access to map value in one step.
func setValues(m *sync.Map, values map[string]interface{}) {
	if m == nil {
		m = new(sync.Map)
	}
	for k, v := range values {
		m.Store(k, v)
	}

}

func getStringVal(m *sync.Map, k string) string {
	v, ok := m.Load(k)
	if !ok {
		return ""
	}

	return Utility.ToString(v)
}

func getIntVal(m *sync.Map, k string) int {
	v, ok := m.Load(k)
	if !ok {
		return 0
	}

	return Utility.ToInt(v)
}

func getBoolVal(m *sync.Map, k string) bool {
	v, ok := m.Load(k)
	if !ok {
		return false
	}

	return Utility.ToBool(v)
}

func getVal(m *sync.Map, k string) interface{} {
	v, ok := m.Load(k)
	if !ok {
		return nil
	}
	return v
}

func (globule *Globule) isInternalService(s *sync.Map) bool {
	name := getStringVal(s, "Name")

	return name == "packages.PackageRepository" ||
		name == "ca.CertificateAuthority" ||
		name == "lb.LoadBalancingService" ||
		name == "log.LogService" ||
		name == "resource.ResourceService" ||
		name == "packages.PackageDiscovery" ||
		name == "rbac.RbacService" ||
		name == "admin.AdminService"

}

func (globule *Globule) getServices() []*sync.Map {
	_services_ := make([]*sync.Map, 0)
	//Append services into the array.
	globule.services.Range(func(key, s interface{}) bool {
		// I will remove unfounded service from the map...
		servicePath := getStringVal(s.(*sync.Map), "Path")
		if !Utility.Exists(servicePath) && !globule.isInternalService(s.(*sync.Map)) {
			// Here I will set various base on the standard dist directory structure.
			path := globule.path + "/services/" + getStringVal(s.(*sync.Map), "PublisherId") + "/" + getStringVal(s.(*sync.Map), "Name") + "/" + getStringVal(s.(*sync.Map), "Version") + "/" + getStringVal(s.(*sync.Map), "Id")
			execName := servicePath[strings.LastIndex(servicePath, "/")+1:]
			servicePath = path + "/" + execName
			if Utility.Exists(servicePath) {
				s.(*sync.Map).Store("Path", servicePath)
				globule.setService(s.(*sync.Map))
				_services_ = append(_services_, s.(*sync.Map))
			} else {
				log.Println("No executable path was found for path ", servicePath)
				globule.deleteService(getStringVal(s.(*sync.Map), "Id"))
			}
		} else {
			// Here the service exec is found.
			_services_ = append(_services_, s.(*sync.Map))
		}

		//_services_ = append(_services_, s.(*sync.Map))
		return true
	})

	return _services_

}

func (globule *Globule) setService(s *sync.Map) {
	id, _ := s.Load("Id")
	// I will not set the services if it
	if getStringVal(s, "State") != "deleted" {
		globule.services.Store(id.(string), s)
	}
}

func (globule *Globule) getService(id string) *sync.Map {
	s, ok := globule.services.Load(id)
	if ok {
		return s.(*sync.Map)
	} else {
		return nil
	}
}

func (globule *Globule) deleteService(id string) {

	s, exist := globule.services.LoadAndDelete(id)
	if exist {
		log.Println("service", getStringVal(s.(*sync.Map), "Name"), getStringVal(s.(*sync.Map), "Id"), "was remove from the map!")
	}
}

func (globule *Globule) toMap() map[string]interface{} {
	_map_, _ := Utility.ToMap(globule)
	_services_ := make(map[string]interface{})

	globule.services.Range(func(key, value interface{}) bool {
		s := make(map[string]interface{})
		value.(*sync.Map).Range(func(key, value interface{}) bool {
			s[key.(string)] = value
			return true
		})
		_services_[key.(string)] = s
		return true
	})
	_map_["Services"] = _services_
	return _map_
}

func processIsRuning(pid int) bool {
	_, err := os.FindProcess(int(pid))
	return err == nil
}

func (globule *Globule) getPortsInUse() []int {
	portsInUse := globule.portsInUse

	// I will test if the port is already taken by e services.
	globule.services.Range(func(key, value interface{}) bool {
		m := value.(*sync.Map)
		pid_, hasProcess := m.Load("Process")

		if hasProcess {
			pid := Utility.ToInt(pid_)
			if pid != -1 {
				if processIsRuning(pid) {
					p, _ := m.Load("Port")
					portsInUse = append(portsInUse, Utility.ToInt(p))
				}
			}
		}

		proxyPid_, hasProxyProcess := m.Load("ProxyProcess")
		if hasProxyProcess {
			proxyPid := Utility.ToInt(proxyPid_)
			if proxyPid != -1 {
				if processIsRuning(proxyPid) {
					p, _ := m.Load("Proxy")
					portsInUse = append(portsInUse, Utility.ToInt(p))
				}
			}
		}
		return true
	})

	return portsInUse
}

/**
 * test if a given port is avalaible.
 */
func (globule *Globule) isPortAvailable(port int) bool {
	portRange := strings.Split(globule.PortsRange, "-")
	start := Utility.ToInt(portRange[0])
	end := Utility.ToInt(portRange[1])

	if port < start || port > end {
		return false
	}

	portsInUse := globule.getPortsInUse()
	for i := 0; i < len(portsInUse); i++ {
		if portsInUse[i] == port {
			return false
		}
	}

	// wait before interogate the next port
	time.Sleep(100 * time.Millisecond)
	l, err := net.Listen("tcp", "0.0.0.0:"+Utility.ToString(port))
	if err == nil {
		defer l.Close()
		return true
	}

	return false
}

/**
 * Return the next available port.
 **/
func (globule *Globule) getNextAvailablePort() (int, error) {
	portRange := strings.Split(globule.PortsRange, "-")
	start := Utility.ToInt(portRange[0]) + 1 // The first port of the range will be reserve to http configuration handler.
	end := Utility.ToInt(portRange[1])

	for i := start; i < end; i++ {
		if globule.isPortAvailable(i) {
			globule.portsInUse = append(globule.portsInUse, i)
			return i, nil
		}
	}

	return -1, errors.New("No port are available in the range " + globule.PortsRange)

}

/**
 * Initialize the server directories config, data, webroot...
 */
func (globule *Globule) initDirectories() {

	// DNS info.
	globule.DNS = make([]interface{}, 0)
	globule.DnsUpdateIpInfos = make([]interface{}, 0)

	// Set the list of discorvery service avalaible...
	globule.Discoveries = make([]string, 0)
	globule.discorveriesEventHub = make(map[string]*event_client.Event_Client)

	// Set the share service info...
	globule.Services = make(map[string]interface{})

	// Set external map services.
	globule.ExternalApplications = make(map[string]ExternalApplication)

	//////////////////////////////////////////////////////////////////////////////////////
	// There is the default directory initialisation...
	//////////////////////////////////////////////////////////////////////////////////////

	// Create the directory if is not exist.
	if globule.path == "/usr/local/share/globular" {
		// Here we have a linux standard installation.
		globule.data = "/var/globular/data"
		globule.webRoot = "/var/globular/webroot"
		globule.config = "/etc/globular/config"
	} else {

		globule.data = globule.path + "/data"
		Utility.CreateDirIfNotExist(globule.data)

		// Configuration directory
		globule.config = globule.path + "/config"
		Utility.CreateDirIfNotExist(globule.config)

		globule.webRoot = globule.path + "/webroot" // The default directory to server.
		// keep the root in global variable for the file handler.
		Utility.CreateDirIfNotExist(globule.webRoot) // Create the directory if it not exist.
	}

	// Create the creds directory if it not already exist.
	globule.creds = globule.config + "/tls"
	Utility.CreateDirIfNotExist(globule.creds)

	// Files directorie that contain user's directories and application's directory
	globule.users = globule.data + "/files/users"
	Utility.CreateDirIfNotExist(globule.users)

	// Contain the application directory.
	globule.applications = globule.data + "/files/applications"
	Utility.CreateDirIfNotExist(globule.applications)

	// Initialyse globular from it configuration file.
	file, err := ioutil.ReadFile(globule.config + "/config.json")

	// Init the service with the default port address
	if err == nil {

		err := json.Unmarshal(file, &globule)
		if err != nil {
			log.Println("fail to read init from ", globule.config+"/config.json", err)
		}

		// Now I will initialyse sync services map.
		for _, v := range globule.Services {
			s := v.(map[string]interface{})
			s_ := new(sync.Map)
			for k_, v_ := range s {
				s_.Store(k_, v_)
			}
			globule.setService(s_)
		}

	} else {
		log.Println("fail to read configuration ", globule.config+"/config.json", err)
	}

	if !Utility.Exists(globule.webRoot + "/index.html") {

		// in that case I will create a new index.html file.
		ioutil.WriteFile(globule.webRoot+"/index.html", []byte(
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

	// Convert video file if there some to be convert.
	go func() {
		convertVideo()
	}()
	log.Println("Globular is running!")
}

/**
 * Start serving the content.
 */
func (globule *Globule) Serve() {

	//globule.initDirectories()

	// Reset previous connections.
	globule.store = nil
	globule.persistence_client_ = nil
	globule.ldap_client_ = nil
	globule.event_client_ = nil

	// Open logs db.
	if globule.logs == nil {

		// The logs storage.
		globule.logs = storage_store.NewLevelDB_store()
		err := globule.logs.Open(`{"path":"` + globule.data + `", "name":"logs"}`)
		if err != nil {
			log.Println(err)
		}

		// The rbac storage.
		globule.permissions = storage_store.NewLevelDB_store()
		err = globule.permissions.Open(`{"path":"` + globule.data + `", "name":"permissions"}`)
		if err != nil {
			log.Println(err)
		}

		// Here it suppose to be only one server instance per computer.
		err = globule.setKey()
		if err != nil {
			log.Panicln(err)
		}

		// The token that identify the server with other services
		err = globule.setToken()
		if err != nil {
			log.Panicln(err)
		}

		// Here I will start the refresh token loop to refresh the server token.
		// the token will be refresh 10 milliseconds before expiration.
		ticker := time.NewTicker((globule.SessionTimeout - 10) * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					err = globule.setToken()
					if err != nil {
						log.Println(err)
					}
				case <-globule.exit:
					return
				}
			}
		}()

		// Start the monitoring service with prometheus.
		globule.startPrometheus()
	}

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// First of all I will start external services.
	for externalServiceId, _ := range globule.ExternalApplications {
		pid, err := globule.startExternalApplication(externalServiceId)
		if err != nil {
			log.Println("fail to start external service: ", externalServiceId, " pid ", pid)
		}
	}

	// I will save the variable in a tmp file to be sure I can get it outside
	ioutil.WriteFile(os.TempDir()+"/GLOBULAR_ROOT", []byte(globule.path+":"+Utility.ToString(globule.PortHttp)), 0644)

	// set the services.
	globule.initServices()

	// start internal services. (need persistence service to manage permissions)
	globule.startInternalServices()

	// lisen
	err := globule.Listen()

	url := globule.Protocol + "://" + globule.Domain

	if globule.Protocol == "https" {
		if globule.PortHttps != 443 {
			url += ":" + Utility.ToString(globule.PortHttps)
		}
	} else if globule.Protocol == "http" {
		if globule.PortHttp != 80 {
			url += ":" + Utility.ToString(globule.PortHttp)
		}
	}

	log.Println("Globular is running at address " + url)

	// Keep watching if the config file was modify by external agent.
	globule.watchConfigFile()

	if err != nil {
		log.Println(err)
	}
}

/**
 * Return the domain of the Globule. The name can be empty. If the name is empty
 * it will mean that the domain is entirely control by the globule so it must be
 * able to do it own validation, other wise the domain validation will be done by
 * the globule asscosiate with that domain.
 */
func (globule *Globule) getDomain() string {
	domain := globule.Domain
	if len(globule.Name) > 0 && domain != "localhost" {
		domain = /*globule.Name + "." +*/ domain
	}
	return domain
}

/**
 * Return the local token.
 */
func (globule *Globule) getToken() (string, error) {
	token, err := ioutil.ReadFile(os.TempDir() + "/" + globule.getDomain() + "_token")
	if err != nil {
		return "", err
	}
	return string(token), nil
}

/**
 * Remove token (for each domain/alternate domains)
 */
func (globule *Globule) deleteToken() {
	os.Remove(os.TempDir() + "/" + globule.getDomain() + "_token")

	// I will also generate the token for the
	for i := 0; i < len(globule.AlternateDomains); i++ {
		os.Remove(os.TempDir() + "/" + globule.AlternateDomains[i].(string) + "_token")
	}
}

/**
 * Generate a new local token that be use to communicate between local service.
 */
func (globule *Globule) setToken() error {

	// remove previous token if some exist.
	globule.deleteToken()

	// create the token for the main domain.
	token, _ := interceptors.GenerateToken(globule.jwtKey, globule.SessionTimeout, "sa", "sa", globule.AdminEmail)

	err := ioutil.WriteFile(os.TempDir()+"/"+globule.getDomain()+"_token", []byte(token), 0400)
	if err != nil {
		return err
	}

	// I will also generate the token for the
	for i := 0; i < len(globule.AlternateDomains); i++ {
		err := ioutil.WriteFile(os.TempDir()+"/"+globule.AlternateDomains[i].(string)+"_token", []byte(token), 0400)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Set the secret key that will be use to validate token. That key will be generate each time the server will be
 * restarted and all token generated with previous key will be automatically invalidated...
 */
func (globule *Globule) setKey() error {
	globule.jwtKey = []byte(Utility.RandomUUID())
	return ioutil.WriteFile(os.TempDir()+"/globular_key", []byte(globule.jwtKey), 0400)
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
	if globule.DNS != nil {
		if len(globule.DNS) > 0 {
			for i := 0; i < len(globule.DNS); i++ {
				dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[i].(string), "dns.DnsService")
				if err != nil {
					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = dns_client_.SetA(globule.Domain, globule.getDomain(), Utility.MyIP(), 60)

				if err != nil {
					// return the setA error
					return err
				}

				// TODO also register the ipv6 here...
				dns_client_.Close()
			}
		}
	}

	// Here If the DNS provides has api to update the ip address I will use it.
	// TODO test it for different internet provider's

	for i := 0; i < len(globule.DnsUpdateIpInfos); i++ {
		// the api call "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
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
	domains := globule.AlternateDomains

	for i := 0; i < len(domains); i++ {
		if !testDomainIp(domains[i].(string), Utility.MyIP(), 3) {
			return errors.New("The domain " + domains[i].(string) + "is not associated with ip " + Utility.MyIP())
		}
	}

	return nil
}

// Test if a domain is asscociated with a given ip.
func testDomainIp(domain string, ip string, try int) bool {
	if try == 0 {
		return false
	}

	if Utility.DomainHasIp(domain, Utility.MyIP()) {
		return true
	} else {
		time.Sleep(5 * time.Second)
		try--
		return testDomainIp(domain, ip, try)
	}

}

/**
 * Start the grpc proxy.
 */
func (globule *Globule) startProxy(s *sync.Map, port int, proxy int) (int, error) {
	_, hasProxyProcess := s.Load("ProxyProcess")
	if !hasProxyProcess {
		s.Store("ProxyProcess", -1)
	}
	pid := getIntVal(s, "ProxyProcess")
	if pid != -1 {
		Utility.TerminateProcess(pid, 0)
	}

	// Now I will start the proxy that will be use by javascript client.
	proxyPath := "/bin/grpcwebproxy"
	if !strings.HasSuffix(proxyPath, ".exe") && runtime.GOOS == "windows" {
		proxyPath += ".exe" // in case of windows.
	}

	proxyBackendAddress := globule.getDomain() + ":" + strconv.Itoa(port)
	proxyAllowAllOrgins := "true"
	proxyArgs := make([]string, 0)

	// Use in a local network or in test.
	proxyArgs = append(proxyArgs, "--backend_addr="+proxyBackendAddress)
	proxyArgs = append(proxyArgs, "--allow_all_origins="+proxyAllowAllOrgins)
	hasTls := getBoolVal(s, "TLS")
	if hasTls {
		certAuthorityTrust := globule.creds + "/ca.crt"

		/* Services gRpc backend. */
		proxyArgs = append(proxyArgs, "--backend_tls=true")
		proxyArgs = append(proxyArgs, "--backend_tls_ca_files="+certAuthorityTrust)
		proxyArgs = append(proxyArgs, "--backend_client_tls_cert_file="+globule.creds+"/client.crt")
		proxyArgs = append(proxyArgs, "--backend_client_tls_key_file="+globule.creds+"/client.pem")

		/* http2 parameters between the browser and the proxy.*/
		proxyArgs = append(proxyArgs, "--run_http_server=false")
		proxyArgs = append(proxyArgs, "--run_tls_server=true")
		proxyArgs = append(proxyArgs, "--server_http_tls_port="+strconv.Itoa(proxy))

		/* in case of public domain server files **/
		proxyArgs = append(proxyArgs, "--server_tls_key_file="+globule.creds+"/server.pem")

		proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+globule.creds+"/"+globule.CertificateAuthorityBundle)
		proxyArgs = append(proxyArgs, "--server_tls_cert_file="+globule.creds+"/"+globule.Certificate)

	} else {
		// Now I will save the file with those new information in it.
		proxyArgs = append(proxyArgs, "--run_http_server=true")
		proxyArgs = append(proxyArgs, "--run_tls_server=false")
		proxyArgs = append(proxyArgs, "--server_http_debug_port="+strconv.Itoa(proxy))
		proxyArgs = append(proxyArgs, "--backend_tls=false")
	}

	// Keep connection open for longer exchange between client/service. Event Subscribe function
	// is a good example of long lasting connection. (48 hours) seam to be more than enought for
	// browser client connection maximum life.
	proxyArgs = append(proxyArgs, "--server_http_max_read_timeout=48h")
	proxyArgs = append(proxyArgs, "--server_http_max_write_timeout=48h")

	// start the proxy service one time
	proxyProcess := exec.Command(globule.path+proxyPath, proxyArgs...)
	proxyProcess.SysProcAttr = &syscall.SysProcAttr{
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err := proxyProcess.Start()

	if err != nil {
		return -1, err
	}

	// save service configuration.
	s.Store("ProxyProcess", proxyProcess.Process.Pid)

	return proxyProcess.Process.Pid, nil
}

/**
 * That function will
 */
func (globule *Globule) keepServiceAlive(s *sync.Map) {

	if globule.exit_ {
		return
	}

	// In case the service must not be kept alive.
	_, HasKeepAlive := s.Load("KeepAlive")
	if !HasKeepAlive {
		return
	}
	// In case the service must not be kept alive.
	keepAlive := getBoolVal(s, "KeepAlive")
	if !keepAlive {
		return
	}


	// In case the service must not be kept alive.
	if getStringVal(s, "State") == "terminated" || getStringVal(s, "State") == "deleted" {
		return
	}

	pid := getIntVal(s, "Process")
	p, err := os.FindProcess(pid)
	if err != nil {
		return
	}

	// Wait for process to return.
	p.Wait()
	if globule.exit_ {
		return
	}
	_, _, err = globule.startService(s)
	if err != nil {
		return
	}
}

/**
 * Start internal service admin and resource are use that function.
 */
func (globule *Globule) startInternalService(id string, proto string, hasTls bool, unaryInterceptor grpc.UnaryServerInterceptor, streamInterceptor grpc.StreamServerInterceptor) (*grpc.Server, int, error) {
	log.Println("Start internal service ", id)

	s := globule.getService(id)
	if s == nil {
		s = new(sync.Map)
	}

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var grpcServer *grpc.Server
	if hasTls {
		certAuthorityTrust := globule.creds + "/ca.crt"
		certFile := globule.creds + "/server.crt"
		keyFile := globule.creds + "/server.pem"

		s.Store("CertFile", certFile)
		s.Store("KeyFile", keyFile)
		s.Store("CertAuthorityTrust", certAuthorityTrust)

		// Create the TLS credentials
		creds := credentials.NewTLS(globular.GetTLSConfig(keyFile, certFile, certAuthorityTrust))

		// Create the gRPC server with the credentials
		opts := []grpc.ServerOption{grpc.Creds(creds),
			grpc.UnaryInterceptor(unaryInterceptor),
			grpc.StreamInterceptor(streamInterceptor)}

		// Create the gRPC server with the credentials
		grpcServer = grpc.NewServer(opts...)

	} else {
		s.Store("CertFile", "")
		s.Store("KeyFile", "")
		s.Store("CertAuthorityTrust", "")

		grpcServer = grpc.NewServer([]grpc.ServerOption{
			grpc.UnaryInterceptor(unaryInterceptor),
			grpc.StreamInterceptor(streamInterceptor)}...)
	}

	reflection.Register(grpcServer)

	// Here I will create the service configuration object.
	s.Store("Domain", globule.getDomain())
	s.Store("Name", id)
	s.Store("Id", id)
	s.Store("Proto", proto)
	s.Store("Port", 0)
	s.Store("Proxy", 0)
	s.Store("TLS", hasTls)
	s.Store("ProxyProcess", -1) // must be use to reserve the port...
	s.Store("Process", -1)

	globule.portsInUse = make([]int, 0)

	// Todo get next available ports.
	port, err := globule.getNextAvailablePort()

	s.Store("Port", port)

	if err != nil {
		return nil, -1, err
	}

	proxy, err := globule.getNextAvailablePort()
	s.Store("Proxy", proxy)

	globule.setService(s)

	if err != nil {
		return nil, -1, err
	}

	// start the proxy
	_, err = globule.startProxy(s, port, proxy)
	if err != nil {
		return nil, -1, err
	}

	globule.inernalServices = append(globule.inernalServices, grpcServer)

	return grpcServer, port, nil
}

/**
 * Stop external services.
 */
func (globule *Globule) stopServices() {
	log.Println("stop services...")
	// not keep services alive because the server must exist.
	globule.exit_ = true

	// Here I will disconnect service update event.
	for id, _ := range globule.subscribers {
		eventHub := globule.discorveriesEventHub[id]
		if eventHub != nil {
			eventHub.Close()
		}
	}

	eventClient, err := globule.getEventHub()
	if err == nil && eventClient != nil {
		eventClient.Close()
	}

	// stop external service.
	for externalServiceId, _ := range globule.ExternalApplications {
		globule.stopExternalApplication(externalServiceId)
	}
	services := globule.getServices()
	for i := 0; i < len(services); i++ {
		s := services[i]
		if s != nil {
			// I will also try to keep a client connection in to communicate with the service.
			log.Println("stop service: ", getStringVal(s, "Name"))
			globule.stopService(s)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if globule.https_server != nil {
		if err := globule.https_server.Shutdown(ctx); err != nil {
			// handle err
			log.Println(err)
		}
		log.Println("stop listen(https) at port ", globule.PortHttps)
	}

	if globule.http_server != nil {
		if err := globule.http_server.Shutdown(ctx); err != nil {
			// handle err
			log.Println(err)
		}
		log.Println("stop listen(http) at port ", globule.PortHttp)
	}

	// Double check that all process are terminated...
	for i := 0; i < len(services); i++ {
		s := services[i]
		processPid := getIntVal(s, "Process")
		if processPid != -1 {
			globule.killServiceProcess(s, processPid)
		}
	}

	globule.saveConfig()
}

/**
 * Start services define in the configuration.
 */
func (globule *Globule) startService(s *sync.Map) (int, int, error) {

	var err error

	root, _ := ioutil.ReadFile(os.TempDir() + "/GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(getStringVal(s, "Domain")) && root_ != globule.path {
		return -1, -1, errors.New("can not start a distant service localy")
	}

	// set the domain of the service.
	s.Store("Domain", globule.getDomain())
	s.Store("TLS", globule.Protocol == "https")

	// if the service already exist.
	_, hasProcess := s.Load("Process")
	if !hasProcess {
		s.Store("Process", -1)
	}

	pid := getIntVal(s, "Process")
	if pid != -1 {
		if runtime.GOOS == "windows" {
			// Program written with dotnet on window need this command to stop...
			kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
			kill.Stderr = os.Stderr
			kill.Stdout = os.Stdout
			kill.Run()
		} else {
			Utility.TerminateProcess(pid, 0)
		}
	}

	// save the process -1 in the map.
	s.Store("Process", -1)
	globule.setService(s)

	servicePath := getStringVal(s, "Path")
	serviceName := getStringVal(s, "Name")
	proxyPid := getIntVal(s, "ProxyProcess")

	if getStringVal(s, "Protocol") == "grpc" && serviceName != "ResourceReesourcervice" && serviceName != "admin.AdminService" && serviceName != "ca.CertificateAuthority" && serviceName != "packages.PackageDiscovery" {
		// I will test if the service is find if not I will try to set path
		// to standard dist directory structure.
		if !Utility.Exists(servicePath) {
			log.Println("No executable path was found for path ", servicePath)
			// Here I will set various base on the standard dist directory structure.
			path := globule.path + "/services/" + getStringVal(s, "PublisherId") + "/" + getStringVal(s, "Name") + "/" + getStringVal(s, "Version") + "/" + getStringVal(s, "Id")
			execName := servicePath[strings.LastIndex(servicePath, "/")+1:]
			servicePath = path + "/" + execName

			if !Utility.Exists(servicePath) {
				// If the service is running...
				if pid != -1 {
					globule.killServiceProcess(s, pid)
				}

				globule.deleteService(getStringVal(s, "Id"))

				defer globule.saveConfig()

				return -1, -1, errors.New("No executable was found for service " + getStringVal(s, "Name") + servicePath)
			}

			s.Store("Path", path+"/"+execName)
			_, exist := s.Load("Path")
			if !exist {
				return -1, -1, errors.New("Fail to retreive exe path " + servicePath)
			}

			// Try to get the prototype from the standard deployement path.
			path_ := globule.path + "/services/" + getStringVal(s, "PublisherId") + "/" + getStringVal(s, "Name") + "/" + getStringVal(s, "Version")
			files, err := Utility.FindFileByName(path_, ".proto")
			if err != nil {
				return -1, -1, errors.New("No prototype file was found for path '" + path_)
			}

			s.Store("Proto", files[0])
		}

		hasTls := getBoolVal(s, "TLS")
		log.Println("Has TLS ", hasTls, getStringVal(s, "Name"))
		if hasTls {
			// Set TLS local services configuration here.
			s.Store("CertAuthorityTrust", globule.creds+"/ca.crt")
			s.Store("CertFile", globule.creds+"/server.crt")
			s.Store("KeyFile", globule.creds+"/server.pem")
		} else {
			// not secure services.
			s.Store("CertAuthorityTrust", "")
			s.Store("CertFile", "")
			s.Store("KeyFile", "")
		}

		// Reset the list of port in user.
		globule.portsInUse = make([]int, 0)

		// Get the next available port.
		port := getIntVal(s, "Port")

		if !globule.isPortAvailable(port) {
			port, err = globule.getNextAvailablePort()
			if err != nil {
				return -1, -1, err
			}
			s.Store("Port", port)
			globule.setService(s)

		}

		err = os.Chmod(servicePath, 0755)
		if err != nil {
			log.Println(err)
		}

		p := exec.Command(servicePath, Utility.ToString(port))
		var errb bytes.Buffer
		pipe, _ := p.StdoutPipe()
		p.Stderr = &errb

		// Here I will set the command dir.
		p.Dir = servicePath[:strings.LastIndex(servicePath, "/")]
		p.SysProcAttr = &syscall.SysProcAttr{
			//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}

		err = p.Start()
		if err != nil {
			s.Store("State", "fail")
			s.Store("Process", -1)
			log.Println("Fail to start service: ", getStringVal(s, "Name"), " at port ", port, " with error ", err)
			return -1, -1, err
		} else {
			pid = p.Process.Pid
			s.Store("Process", p.Process.Pid)
			log.Println("Service ", getStringVal(s, "Name")+":"+getStringVal(s, "Id"), "started with pid:", getIntVal(s, "Process"))
		}

		// Now I specific services necessary actions...
		if getStringVal(s, "Name") == "persistence.PersistenceService" {
			globule.persistence_client_ = nil
		}

		// save the services in the map.
		go func(s *sync.Map, p *exec.Cmd) {

			s.Store("State", "running")
			globule.keepServiceAlive(s)

			output := make(chan string)
			done := make(chan bool)

			// Process message util the command is done.
			go func() {
				for {
					select {
					case <-done:
						return

					case line := <-output:
						log.Println(line)
						globule.logServiceInfo(getStringVal(s, "Name"), line)
					}
				}

			}()

			// Start reading the output
			go ReadOutput(output, pipe)

			// if the process is not define.
			err = p.Wait() // wait for the program to return
			done <- true
			pipe.Close()

			if err != nil {
				// I will log the program error into the admin logger.
				globule.logServiceError(getStringVal(s, "Name"), err.Error())
			}

			// Print the error
			if len(errb.String()) > 0 {
				fmt.Println("service", getStringVal(s, "Name"), "err:", errb.String())
				globule.logServiceError(getStringVal(s, "Name"), errb.String())
			}

			if !getBoolVal(s, "KeepAlive") || getStringVal(s, "State") == "terminated"  || getStringVal(s, "State") == "deleted" {
				// Terminate it proxy process if not keep alive.
				Utility.TerminateProcess(proxyPid, 0)
				s.Store("ProxyProcess", -1)
			}

			globule.logServiceInfo(getStringVal(s, "Name"), "Service stop.")
			s.Store("Process", -1)
			globule.setService(s)

		}(s, p)

		// get another port.
		if proxyPid == -1 {
			proxy := getIntVal(s, "Proxy")
			if !globule.isPortAvailable(proxy) {
				globule.setService(s)
				proxy, err = globule.getNextAvailablePort()
				if err != nil {
					s.Store("Proxy", -1)

					return -1, -1, err
				}
				// Set back the process
				s.Store("Proxy", proxy)
				globule.setService(s)
			}

			// Start the proxy.
			proxyPid, err = globule.startProxy(s, port, proxy)
			if err != nil {
				return -1, -1, err
			}
		}

		// save service config.
		globule.saveServiceConfig(s)

		// save it to the config because pid and proxy pid have change.
		globule.saveConfig()
		proxy := getIntVal(s, "Proxy")
		log.Println("Service "+getStringVal(s, "Name")+":"+getStringVal(s, "Id")+" is up and running at port ", port, " and proxy ", proxy)

	} else if getStringVal(s, "Protocol") == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(getStringVal(s, "Name"), "Globular") {
			p := exec.Command(servicePath, getStringVal(s, "Port"))

			var errb bytes.Buffer
			pipe, _ := p.StdoutPipe()
			p.Stderr = &errb

			// Here I will set the command dir.
			p.Dir = servicePath[:strings.LastIndex(servicePath, "/")]
			p.SysProcAttr = &syscall.SysProcAttr{
				//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			}

			err = p.Start()

			if err != nil {
				// The process already exist so I will not throw an error and I will use existing process instead. I will make the
				if err.Error() != "exec: already started" {
					s.Store("Process", -1)
					s.Store("State", "fail")
					log.Println("Fail to start service: ", getStringVal(s, "Name"), " at port ", getStringVal(s, "Port"), " with error ", err)
					return -1, -1, err
				}
			}
			pid = p.Process.Pid
			s.Store("Process", p.Process.Pid)
			s.Store("State", "running")

			if err == nil {
				go func(s *sync.Map) {

					globule.keepServiceAlive(s)

					// display the message in the console.
					reader := bufio.NewReader(pipe)
					line, err := reader.ReadString('\n')
					for err == nil {
						name := getStringVal(s, "Name")
						log.Println(name, ":", line)
						line, err = reader.ReadString('\n')
					}

					// if the process is not define.
					pid := getIntVal(s, "Process")
					if pid == -1 {
						log.Println("No process found for service", getStringVal(s, "Name"))
					}
					p, err := os.FindProcess(pid)
					if err == nil {
						_, err := p.Wait()
						if err != nil {
							// I will log the program error into the admin logger.
							globule.logServiceInfo(getStringVal(s, "Name"), errb.String())
						}
					}
				}(s)
			}

			// Save configuration stuff.
			globule.setService(s)
		}
	}

	if pid == -1 {
		s.Store("State", "fail")
		globule.setService(s)
		err := errors.New("Fail to start process " + getStringVal(s, "Name"))
		return -1, -1, err
	}

	// Return the pid of the service.
	if proxyPid != -1 {
		s.Store("State", "running")
		globule.setService(s)
		return pid, proxyPid, nil
	}

	return pid, -1, nil
}

/**
 * Init services configuration.
 */
func (globule *Globule) initService(s *sync.Map) error {
	_, hasProtocol := s.Load("Protocol")
	if !hasProtocol {
		// internal service dosent has Protocol define.
		return nil
	}

	if getStringVal(s, "Protocol") == "grpc" {
		// The domain must be set in the sever configuration and not change after that.
		hasTls := globule.Protocol == "https" //Utility.ToBool(s["TLS"])
		s.Store("TLS", hasTls)                // set the tls...
		if hasTls {
			// Set TLS local services configuration here.
			s.Store("CertAuthorityTrust", globule.creds+"/ca.crt")
			s.Store("CertFile", globule.creds+"/server.crt")
			s.Store("KeyFile", globule.creds+"/server.pem")
		} else {
			// not secure services.
			s.Store("CertAuthorityTrust", "")
			s.Store("CertFile", "")
			s.Store("KeyFile", "")
		}

		// Set the default server value.
		if globule.KeepAllServicesAlive {
			s.Store("KeepAlive", true)
		}

		// Keep service up to date.
		if globule.KeepAllServicesUpToDate {
			s.Store("KeepUpToDate", true)
		}
	}

	// any other http server except this one...
	if !strings.HasPrefix(getStringVal(s, "Name"), "Globular") {
		hasChange := globule.saveServiceConfig(s)
		state := getStringVal(s, "State")
		if hasChange || state == "stopped" {
			if state == "stop" {
				globule.stopService(s)
			} else {
				// TODO watch here wath to do if other conditio are set.
				// here the service will try to restart.
				_, _, err := globule.startService(s)
				if err != nil {
					s.Store("State", "failed")
					return err
				}
			}
			globule.setService(s)
		}
	}

	return nil
}

func (globule *Globule) getBasePath() string {
	// Each service contain a file name config.json that describe service.
	// I will keep services info in services map and also it running process.
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Start from development environnement.
	if Utility.Exists("README.md") {
		// GLOBULAR_SERVICES_ROOT is the path of the globular service executables.
		// if not set the services must be in the same folder as Globurar executable.
		globularServicesRoot := os.Getenv("GLOBULAR_SERVICES_ROOT")
		log.Println("GLOBULAR_SERVICES_ROOT ", globularServicesRoot)
		if len(globularServicesRoot) > 0 {
			basePath = globularServicesRoot
		}
	}
	return basePath
}

func (globule *Globule) killServiceProcess(s *sync.Map, pid int) {

	// Here I will set a variable that tell globular to not keep the service alive...
	s.Store("State", "terminated")

	// also kill it proxy process if exist in that case.
	_, hasProxyProcess := s.Load("ProxyProcess")
	if hasProxyProcess {
		proxyProcessPid := getIntVal(s, "ProxyProcess")
		proxyProcess, err := os.FindProcess(proxyProcessPid)
		if err == nil {
			proxyProcess.Kill()
			s.Store("ProxyProcess", -1)
		}
	}

	// kill it in the name of...
	process, err := os.FindProcess(pid)
	if err == nil {
		err := process.Kill()
		if err == nil {
			s.Store("Process", -1)
			s.Store("State", "stopped")
		} else {
			s.Store("State", "failed")
		}
	}

	// Set the service
	globule.setService(s)
}

/**
 * Call once when the server start.
 */
func (globule *Globule) initServices() {

	log.Println("Initialyse services")
	log.Println("local ip ", Utility.MyLocalIP())
	log.Println("external ip ", Utility.MyIP())

	// If the protocol is https I will generate the TLS certificate.
	if globule.Protocol == "https" {
		// security.GenerateServicesCertificates(globule.CertPassword, globule.CertExpirationDelay, globule.getDomain(), globule.creds)
		if len(globule.Certificate) == 0 {
			globule.registerIpToDns()
			log.Println(" Now let's encrypts!")
			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/tls/server.csr run
			// I need to remove the gRPC certificate and recreate it.

			Utility.RemoveDirContents(globule.creds)

			// recreate the certificates.
			err := security.GenerateServicesCertificates(globule.CertPassword, globule.CertExpirationDelay, globule.getDomain(), globule.creds, globule.Country, globule.State, globule.City, globule.Organization, globule.AlternateDomains)
			if err != nil {
				log.Println(err)
				return
			}

			err = globule.obtainCertificateForCsr()
			if err != nil {
				log.Println(err)
				return
			}
		}

		// Here I will read the certificate
		r, _ := ioutil.ReadFile(globule.creds + "/" + globule.Certificate)
		block, _ := pem.Decode(r)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err == nil {
			delay := cert.NotAfter.Sub(time.Now()) - time.Duration(15*time.Minute)
			timeout := time.NewTimer(delay)
			go func() {
				// Wait to restart the server to regenerate new certificates...
				<-timeout.C
				globule.Certificate = ""
				globule.restartServices()
			}()
		}
	}

	// That will contain all method path from the proto files.
	globule.methods = make([]string, 0)
	globule.methods = append(globule.methods, "/file.FileService/FileUploadHandler")

	// Set local action permission
	globule.setActionResourcesPermissions(map[string]interface{}{"action": "/resource.ResourceService/DeletePermissions", "resources": []interface{}{map[string]interface{}{"index": 0, "permission": "delete"}}})
	globule.setActionResourcesPermissions(map[string]interface{}{"action": "/resource.ResourceService/SetResourceOwner", "resources": []interface{}{map[string]interface{}{"index": 0, "permission": "write"}}})
	globule.setActionResourcesPermissions(map[string]interface{}{"action": "/resource.ResourceService/DeleteResourceOwner", "resources": []interface{}{map[string]interface{}{"index": 0, "permission": "write"}}})
	globule.setActionResourcesPermissions(map[string]interface{}{"action": "/admin.AdminService/DeployApplication", "resources": []interface{}{map[string]interface{}{"index": 0, "permission": "write"}}})
	globule.setActionResourcesPermissions(map[string]interface{}{"action": "/admin.AdminService/PublishService", "resources": []interface{}{map[string]interface{}{"index": 0, "permission": "write"}}})

	// It will be execute the first time only...
	configPath := globule.config + "/config.json"
	if !Utility.Exists(configPath) {
		filepath.Walk(globule.getBasePath(), func(path string, info os.FileInfo, err error) error {
			path = strings.ReplaceAll(path, "\\", "/")
			if info == nil {
				return nil
			}

			if err == nil && info.Name() == "config.json" {
				// So here I will read the content of the file.
				s := make(map[string]interface{})
				config, err := ioutil.ReadFile(path)
				if err == nil {
					// Read the config file.
					err := json.Unmarshal(config, &s)
					if err == nil {
						if s["Protocol"] != nil {
							// If a configuration file exist It will be use to start services,
							// otherwise the service configuration file will be use.
							if s["Name"] == nil {
								log.Println("---> no 'Name' attribute found in service configuration in file config ", path)
							} else {

								// if no id was given I will generate a uuid.
								if s["Id"] == nil {
									s["Id"] = Utility.RandomUUID()
								}

								s_ := new(sync.Map)
								for k, v := range s {
									s_.Store(k, v)
								}

								globule.setService(s_)
							}
						}
					} else {
						log.Println("fail to unmarshal configuration ", err)
					}
				} else {
					log.Println("Fail to read config file ", path, err)
				}
			}
			return nil
		})
	}

	// Set service methods.
	filepath.Walk(globule.getBasePath(), func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		if info == nil {
			return nil
		}
		if err == nil && strings.HasSuffix(info.Name(), ".proto") {
			name := info.Name()[0:strings.Index(info.Name(), ".")]
			globule.setServiceMethods(name, path)
		}
		return nil
	})

	// Set the certificate keys...
	services := globule.getServices()
	for _, s := range services {
		log.Println("init service ", getStringVal(s, "Name"))
		if getStringVal(s, "Protocol") == "grpc" {

			// The domain must be set in the sever configuration and not change after that.
			hasTls := globule.Protocol == "https" //Utility.ToBool(s["TLS"])
			s.Store("TLS", hasTls)                // set the tls...
			if hasTls {
				// Set TLS local services configuration here.
				s.Store("CertAuthorityTrust", globule.creds+"/ca.crt")
				s.Store("CertFile", globule.creds+"/server.crt")
				s.Store("KeyFile", globule.creds+"/server.pem")
			} else {
				// not secure services.
				s.Store("CertAuthorityTrust", "")
				s.Store("CertFile", "")
				s.Store("KeyFile", "")
			}
		}
	}

	// Start the load balancer.
	err := globule.startLoadBalancingService()
	if err != nil {
		log.Println(err)
	}

	log.Println("Init external services: ")
	servicesByName := make(map[string][]int, 0)

	// Initialyse service

	for _, s := range services {
		name := getStringVal(s, "Name")
		// Get existion process information.
		_, hasProcess := s.Load("Process")
		processPid := -1
		if hasProcess {
			processPid = getIntVal(s, "Process")
			// Now I will find if the process is running
			if processPid != -1 {
				_, err := Utility.GetProcessRunningStatus(processPid)
				log.Println("find process ", name, ":", processPid)
				if err != nil {
					globule.killServiceProcess(s, processPid)
					processPid = -1
				} else {
					p, err := ps.FindProcess(processPid)
					if err != nil {
						log.Println("Process ", processPid, " dosent exist anymore...")
						globule.killServiceProcess(s, processPid)
						processPid = -1
					} else {
						processExist := p.Executable() == getStringVal(s, "Path")
						if !processExist {
							log.Println("Process ", processPid, " dosent exist anymore...")
							globule.killServiceProcess(s, processPid)
							processPid = -1
						}
					}
				}
			}
		}

		// Here I will keep track of the services process by it name...
		if _, ok := servicesByName[name]; !ok {
			//do something here
			servicesByName[name] = make([]int, 0)
		}

		// Keep the pid.
		servicesByName[name] = append(servicesByName[name], processPid)

		// If no process already exist I will create one.
		if processPid == -1 {
			_, hasProxyProcess := s.Load("ProxyProcess")
			if hasProxyProcess {
				proxyProcessPid := getIntVal(s, "ProxyProcess")
				proxyProcess, err := os.FindProcess(proxyProcessPid)
				if err == nil {
					proxyProcess.Kill()
				}
			}

			// The service name.
			log.Println("Init service: ", name)
			if name == "file.FileService" {
				s.Store("Root", globule.data+"/files")
			} else if name == "conversation.ConversationService" {
				s.Store("Root", globule.data)
			}

			err := globule.initService(s)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("Process exist for service: ", name)
		}
	}

	// Now I will kill all services who are not listed in the map.
	for serviceName, pids := range servicesByName {
		pids_, err := Utility.GetProcessIdsByName(serviceName)
		if err == nil {
			exist := false
			for i := 0; i < len(pids_); i++ {
				for j := 0; j < len(pids); j++ {
					if pids_[i] == pids[j] {
						exist = true
						break
					}
				}
				// If the pid is not found in the list of pids from the configuration
				// I will remove it.
				if !exist {
					p, err := os.FindProcess(pids_[i])
					if err == nil {
						p.Kill()
					}
				}
			}
		}
	}
}

// That function resolve import path.
func resolveImportPath(path string, importPath string) (string, error) {

	// firt of all i will keep only the path part of the import...
	startIndex := strings.Index(importPath, `'@`) + 1
	endIndex := strings.LastIndex(importPath, `'`)
	importPath_ := importPath[startIndex:endIndex]

	filepath.Walk(globule.webRoot+path[0:strings.Index(path, "/")],
		func(path string, info os.FileInfo, err error) error {
			path = strings.ReplaceAll(path, "\\", "/")
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, importPath_) {
				importPath_ = path
				return io.EOF
			}

			return nil
		})

	importPath_ = strings.Replace(importPath_, strings.Replace(globule.webRoot, "\\", "/", -1), "", -1)

	// Now i will make the path relative.
	importPath__ := strings.Split(importPath_, "/")
	path__ := strings.Split(path, "/")

	var index int
	for ; importPath__[index] == path__[index]; index++ {
	}

	importPath_ = ""

	// move up part..
	for i := index; i < len(path__)-1; i++ {
		importPath_ += "../"
	}

	// go down to the file.
	for i := index; i < len(importPath__); i++ {
		importPath_ += importPath__[i]
		if i < len(importPath__)-1 {
			importPath_ += "/"
		}
	}

	// remove the
	importPath_ = strings.Replace(importPath_, globule.webRoot, "", 1)

	// remove the root path part and the leading / caracter.
	return importPath_, nil
}

/**
 * Start prometheus.
 */
func (globule *Globule) startPrometheus() error {

	var err error

	// Here I will start promethus.
	dataPath := globule.data + "/prometheus-data"
	Utility.CreateDirIfNotExist(dataPath)
	if !Utility.Exists(globule.config + "/prometheus.yml") {
		config := `# my global config
global:
  scrape_interval:     15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
  evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
  # scrape_timeout is set to the global default (10s).

# Alertmanager configuration
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      # - alertmanager:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
rule_files:
  # - "first_rules.yml"
  # - "second_rules.yml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  - job_name: 'prometheus'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['localhost:9090']
  
  - job_name: 'globular_internal_services_metrics'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:` + Utility.ToString(globule.PortHttp) + `']
    
  - job_name: 'node_exporter_metrics'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:9100']
    
  - job_name: 'plc_exporter'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:2112']
`
		err := ioutil.WriteFile(globule.config+"/prometheus.yml", []byte(config), 0644)
		if err != nil {
			return err
		}
	}

	if !Utility.Exists(globule.config + "/alertmanager.yml") {
		config := `global:
  resolve_timeout: 5m

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'
receivers:
- name: 'web.hook'
  webhook_configs:
  - url: 'http://127.0.0.1:5001/'
inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'dev', 'instance']
`
		err := ioutil.WriteFile(globule.config+"/alertmanager.yml", []byte(config), 0644)
		if err != nil {
			return err
		}
	}

	prometheusCmd := exec.Command("prometheus", "--web.listen-address", "0.0.0.0:9090", "--config.file", globule.config+"/prometheus.yml", "--storage.tsdb.path", dataPath)
	err = prometheusCmd.Start()
	prometheusCmd.SysProcAttr = &syscall.SysProcAttr{
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	if err != nil {
		log.Println("fail to start prometheus ", err)
		return err
	}

	// Here I will register various metric that I would like to have for the dashboard.

	// Prometheus logging informations.
	globule.methodsCounterLog = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "globular_methods_counter",
		Help: "Globular services methods usage.",
	},
		[]string{
			"application",
			"user",
			"method"},
	)
	prometheus.MustRegister(globule.methodsCounterLog)

	// Here I will monitor the cpu usage of each services
	globule.servicesCpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "globular_services_cpu_usage_counter",
		Help: "Monitor the cpu usage of each services.",
	},
		[]string{
			"id",
			"name"},
	)

	globule.servicesMemoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "globular_services_memory_usage_counter",
		Help: "Monitor the memory usage of each services.",
	},
		[]string{
			"id",
			"name"},
	)

	// Set the function into prometheus.
	prometheus.MustRegister(globule.servicesCpuUsage)
	prometheus.MustRegister(globule.servicesMemoryUsage)

	// Start feeding the time series...
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				globule.services.Range(func(key, s interface{}) bool {
					pids, err := Utility.GetProcessIdsByName("Globular")
					if err == nil {
						for i := 0; i < len(pids); i++ {
							sysInfo, err := pidusage.GetStat(pids[i])
							if err == nil {
								//log.Println("---> set cpu for process ", pid, getStringVal(s.(*sync.Map), "Name"), sysInfo.CPU)
								globule.servicesCpuUsage.WithLabelValues("Globular", "Globular").Set(sysInfo.CPU)
								globule.servicesMemoryUsage.WithLabelValues("Globular", "Globular").Set(sysInfo.Memory)
							}
						}
					}

					pid := getIntVal(s.(*sync.Map), "Process")
					if pid > 0 {
						sysInfo, err := pidusage.GetStat(pid)
						if err == nil {
							//log.Println("---> set cpu for process ", pid, getStringVal(s.(*sync.Map), "Name"), sysInfo.CPU)
							globule.servicesCpuUsage.WithLabelValues(getStringVal(s.(*sync.Map), "Id"), getStringVal(s.(*sync.Map), "Name")).Set(sysInfo.CPU)
							globule.servicesMemoryUsage.WithLabelValues(getStringVal(s.(*sync.Map), "Id"), getStringVal(s.(*sync.Map), "Name")).Set(sysInfo.Memory)
						}
					} else {
						path := getStringVal(s.(*sync.Map), "Path")
						if len(path) > 0 {
							globule.servicesCpuUsage.WithLabelValues(getStringVal(s.(*sync.Map), "Id"), getStringVal(s.(*sync.Map), "Name")).Set(0)
							globule.servicesMemoryUsage.WithLabelValues(getStringVal(s.(*sync.Map), "Id"), getStringVal(s.(*sync.Map), "Name")).Set(0)
							//log.Println("----> process is close for ", getStringVal(s.(*sync.Map), "Name"))
						}

					}
					return true
				})
			case <-globule.exit:
				return
			}
		}

	}()

	alertmanager := exec.Command("alertmanager", "--config.file", globule.config+"/alertmanager.yml")
	alertmanager.SysProcAttr = &syscall.SysProcAttr{
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = alertmanager.Start()
	if err != nil {
		log.Println("fail to start prometheus alert manager", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	node_exporter := exec.Command("node_exporter")
	node_exporter.SysProcAttr = &syscall.SysProcAttr{
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = node_exporter.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	return nil
}

/**
 * Connection with local persistence grpc service
 */
func (globule *Globule) getPersistenceSaConnection() (*persistence_client.Persistence_Client, error) {
	// That service made user of persistence service.
	if globule.persistence_client_ != nil {
		// Here I will also test the connection...
		err := globule.persistence_client_.Ping("local_resource")
		if err == nil {
			return globule.persistence_client_, nil
		}
		globule.persistence_client_ = nil // set back value to nil.
	}

	configs := globule.getServiceConfigByName("persistence.PersistenceService")
	if len(configs) == 0 {
		err := errors.New("no persistence service configuration was found on that server")
		return nil, err
	}

	var err error
	s := configs[0]

	// Cast-it to the persistence client.
	globule.persistence_client_, err = persistence_client.NewPersistenceService_Client(s["Domain"].(string)+":"+Utility.ToString(globule.PortHttp), s["Id"].(string))
	if err != nil {
		return nil, err
	}

	domain, port := globule.getBackendAddress()

	// Connect to the database here.
	err = globule.persistence_client_.CreateConnection("local_resource", "local_resource", domain, Utility.ToNumeric(port), 0, "sa", globule.RootPassword, 5000, "", false)
	if err != nil {
		return nil, err
	}

	return globule.persistence_client_, nil
}

func (globule *Globule) getBackendAddress() (string, int32) {
	values := strings.Split(globule.DbIpV4, ":")
	return values[0], int32(Utility.ToInt(values[1]))
}

/**
 * Connection to mongo db local store.
 */
func (globule *Globule) getPersistenceStore() (persistence_store.Store, error) {
	// That service made user of persistence service.
	if globule.store == nil {
		globule.store = new(persistence_store.MongoStore)
		domain, port := globule.getBackendAddress()
		err := globule.store.Connect("local_resource", domain, port, "sa", globule.RootPassword, "local_resource", 5000, "")
		if err != nil {
			return nil, err
		}
	}

	return globule.store, nil
}

/** Stop mongod process **/
func (globule *Globule) stopMongod() error {
	closeCmd := exec.Command("mongo", "--eval", "db=db.getSiblingDB('admin');db.adminCommand( { shutdown: 1 } );")
	err := closeCmd.Run()
	time.Sleep(1 * time.Second)
	return err
}

func (globule *Globule) waitForMongo(timeout int, withAuth bool) error {
	logger.Info("Wait for starting mongo db!")
	time.Sleep(1 * time.Second)
	args := make([]string, 0)
	if withAuth {
		args = append(args, "-u")
		args = append(args, "sa")
		args = append(args, "-p")
		args = append(args, globule.RootPassword)
		args = append(args, "--authenticationDatabase")
		args = append(args, "admin")
	}
	args = append(args, "--eval")
	args = append(args, "db=db.getSiblingDB('admin');db.getMongo().getDBNames()")

	script := exec.Command("mongo", args...)
	err := script.Run()
	if err != nil {
		log.Println("wait for mongo...", timeout, "s")
		logger.Info("Fail to start mongod ", err)
		if timeout == 0 {
			return errors.New("mongod is not responding")
		}
		// call again.
		timeout -= 1

		return globule.waitForMongo(timeout, withAuth)
	}

	// Now I will initialyse the application connections...
	globule.createApplicationConnection()

	return nil
}

func (globule *Globule) getLdapClient() (*ldap_client.LDAP_Client, error) {

	configs := globule.getServiceConfigByName("ldap.LdapService")
	if len(configs) == 0 {
		return nil, errors.New("no event service was configure on that globule")
	}

	var err error

	s := configs[0]

	if globule.ldap_client_ == nil {
		globule.ldap_client_, err = ldap_client.NewLdapService_Client(s["Domain"].(string)+":"+Utility.ToString(globule.PortHttp), "ldap.LdapService")
	}

	return globule.ldap_client_, err
}

/**
 * Get access to the event services.
 */
func (globule *Globule) getEventHub() (*event_client.Event_Client, error) {

	// Here I will get a look into the list of initialyse process before trying to connect
	// to event service.
	configs := globule.getServiceConfigByName("event.EventService")
	if len(configs) == 0 {
		return nil, errors.New("no event service was configure on that globule")
	}
	s := configs[0]

	var err error
	if globule.event_client_ == nil {
		globule.event_client_, err = event_client.NewEventService_Client(s["Domain"].(string), s["Id"].(string))
		if err == nil {
			// Here I need to publish a fake event message to be sure the event service is listen.
			err := globule.event_client_.Publish("__init__", []byte("is there anybody out there...?"))
			if err != nil {
				globule.event_client_ = nil
				return nil, err
			}
		}
	}
	log.Println("connection to local event hub succed!")
	return globule.event_client_, err

}

/**
 * The file client is use to access file directories where users and application
 * upload their file. File upload and download are manage by the file service and
 * not by http handler.
 */
func (globule *Globule) GetFileClient(id string) (*file_client.File_Client, error) {

	if globule.file_clients_ == nil {
		globule.file_clients_ = new(sync.Map)
	} else {
		c_, ok := globule.file_clients_.Load(id)
		if ok {
			return c_.(*file_client.File_Client), nil
		}
	}

	return nil, errors.New("No file client found on the server with id '" + id + "'")
}

func (globule *Globule) GetAbsolutePath(path string) string {

	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-2]
	}

	if len(path) > 1 {
		if strings.HasPrefix(path, "/") {
			path = globule.webRoot + path
		} else if !strings.HasSuffix(path, "/") {
			path = globule.webRoot + "/" + path
		} else {
			path = globule.webRoot + path
		}
	} else {
		path = globule.webRoot
	}

	return path

}

func (globule *Globule) startInternalServices() error {

	// Start internal services.

	// Admin service
	err := globule.startAdminService()
	if err != nil {
		return err
	}

	// Log service
	err = globule.startLogService()
	if err != nil {
		return err
	}

	// Resource service
	err = globule.startResourceService()
	if err != nil {
		return err
	}

	// Start Role Based Access Control (RBAC) service.
	err = globule.startRbacService()
	if err != nil {
		return err
	}

	// Directorie service
	err = globule.startPackagesDiscoveryService()
	if err != nil {
		return err
	}

	// Repository service
	err = globule.startPackagesRepositoryService()
	if err != nil {
		return err
	}

	// Certificate autority service.
	err = globule.startCertificateAuthorityService()
	if err != nil {
		return err
	}

	// save the config.
	globule.saveConfig()

	return nil
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

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(bodyBytes), nil
	}

	return "", errors.New("fail to retreive checksum with error " + Utility.ToString(resp.StatusCode))
}

func (globule *Globule) watchForUpdate() {
	for {

		// stop watching if exit was call...
		if globule.exit_ {
			return
		}

		if len(globule.Discoveries) > 0 {
			// Here I will retreive the checksum information from it parent.
			discovery := globule.Discoveries[0]
			address := strings.Split(discovery, ":")[0]
			port := 80
			if strings.Contains(discovery, ":") {
				port = Utility.ToInt(strings.Split(discovery, ":")[1])
			}

			// Here I will test if the checksum has change...
			checksum, err := getChecksum(address, port)
			execPath := Utility.GetExecName(os.Args[0])
			if Utility.Exists("/usr/local/share/globular/Globular") {
				execPath = "/usr/local/share/globular/Globular"
			}
			if err == nil {
				if checksum != Utility.CreateFileChecksum(execPath) {

					if globule.Domain != address && globule.KeepUpToDate {
						err := update_globular_from(globule, discovery, globule.Domain, "sa", globule.RootPassword, runtime.GOOS+":"+runtime.GOARCH)
						if err != nil {
							log.Println("fail to update globular from " + discovery + " with error " + err.Error())
						} else {
							log.Println("update globular checksum is ", checksum)
						}
					}
				}
			}
		}

		// The time here can be set to higher value.
		time.Sleep(30 * time.Second)
	}
}

/**
 * Keep globular version in sync with the version from it discovery [0]...
 */
func (globule *Globule) keepUpToDate() {
	go func() {
		globule.watchForUpdate()
	}()
}

/**
 * Listen for new connection.
 */
func (globule *Globule) Listen() error {

	// Keep services up to date subscription.
	globule.subscribers = globule.keepServicesUpToDate()

	var err error
	// Must be started before other services.
	go func() {
		// local - non secure connection.
		globule.http_server = &http.Server{
			Addr: ":" + strconv.Itoa(globule.PortHttp),
		}
		err = globule.http_server.ListenAndServe()
	}()

	// Here I will make a signal hook to interrupt to exit cleanly.
	// handle the Interrupt
	// set the register sa user.
	globule.registerSa()

	// Keep globular up to date subscription.
	globule.keepUpToDate()

	// Start the http server.
	if globule.Protocol == "https" {

		// if no certificates are specified I will try to get one from let's encrypts.
		// Start https server.
		globule.https_server = &http.Server{
			Addr: ":" + strconv.Itoa(globule.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: globule.getDomain(),
			},
		}

		// get the value from the configuration files.
		go func() {
			err = globule.https_server.ListenAndServeTLS(globule.creds+"/"+globule.Certificate, globule.creds+"/server.pem")
		}()
	}

	return err
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
 * I will reuse the client public key here as key instead of generate another key
 * and manage it...
 */
func (globule *Globule) GetPrivateKey() crypto.PrivateKey {
	keyPem, err := ioutil.ReadFile(globule.creds + "/client.pem")
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

/**
 * That function work correctly, but the DNS fail time to time to give the
 * IP address that result in a fail request... The DNS must be fix!
 */
func (globule *Globule) obtainCertificateForCsr() error {

	config := lego.NewConfig(globule)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", strconv.Itoa(globule.PortHttp)))
	if err != nil {
		log.Fatal(err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	globule.registration = reg
	if err != nil {
		return err
	}

	csrPem, err := ioutil.ReadFile(globule.creds + "/server.csr")
	if err != nil {
		return err
	}

	csrBlock, _ := pem.Decode(csrPem)
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return err
	}

	resource, err := client.Certificate.ObtainForCSR(*csr, true)
	if err != nil {
		return err
	}

	// Keep certificates url in the config.
	globule.CertURL = resource.CertURL
	globule.CertStableURL = resource.CertStableURL

	// Set the certificates paths...
	globule.Certificate = globule.getDomain() + ".crt"
	globule.CertificateAuthorityBundle = globule.getDomain() + ".issuer.crt"

	// Save the certificate in the cerst folder.
	ioutil.WriteFile(globule.creds+"/"+globule.Certificate, resource.Certificate, 0400)
	ioutil.WriteFile(globule.creds+"/"+globule.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	globule.saveConfig()

	return nil
}
