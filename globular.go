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

	"github.com/davecourtois/Globular/services/golang/lb/lbpb"

	"github.com/davecourtois/Globular/services/golang/dns/dns_client"
	"github.com/davecourtois/Globular/services/golang/event/event_client"
	"github.com/davecourtois/Globular/services/golang/ldap/ldap_client"
	"github.com/davecourtois/Globular/services/golang/monitoring/monitoring_client"
	"github.com/davecourtois/Globular/services/golang/persistence/persistence_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/prometheus/client_golang/prometheus"

	// Interceptor for authentication, event, log...
	"github.com/davecourtois/Globular/Interceptors"

	// Client services.
	"crypto"

	"github.com/davecourtois/Globular/services/golang/storage/storage_store"
	"github.com/davecourtois/Utility"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/davecourtois/Globular/security"
	globular "github.com/davecourtois/Globular/services/golang/globular_service"
	"github.com/davecourtois/Globular/services/golang/persistence/persistence_store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global variable.
var (
	webRoot string
	globule *Globule
)

const serviceStartDelay = 2 // wait tow second.

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
	PortHttp                  int    // The port of the http file server.
	PortHttps                 int    // The secure port
	AdminPort                 int    // The admin port
	AdminProxy                int    // The admin proxy port.
	AdminEmail                string // The admin email
	RessourcePort             int    // The ressource management service port
	RessourceProxy            int    // The ressource management proxy port
	CertificateAuthorityPort  int    // The certificate authority port
	CertificateAuthorityProxy int    // The certificate authority proxy port
	ServicesDiscoveryPort     int    // The services discovery port
	ServicesDiscoveryProxy    int    // The ressource management proxy port
	ServicesRepositoryPort    int    // The services discovery port
	ServicesRepositoryProxy   int    // The ressource management proxy port
	LoadBalancingServicePort  int    // The load balancing service port
	LoadBalancingServiceProxy int    // The load balancing proxy port
	PortsRange                string // The range of port to be use for the service. ex 10000-10200

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

	Version        string
	Platform       string
	SessionTimeout time.Duration

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep service up to date.

	// DNS stuff.
	DNS []string // Domain name server use to located the server.

	DnsUpdateIpInfos []interface{} // The internet provader SetA info to keep ip up to date.

	discorveriesEventHub map[string]*event_client.Event_Client

	// The list of method supported by this server.
	methods []string

	// Array of action permissions
	actionPermissions []interface{}

	// The prometheus logging informations.
	methodsCounterLog *prometheus.CounterVec

	// prometheus.CounterVec

	// Directories.
	path    string // The path of the exec...
	webRoot string // The root of the http file server.
	data    string // the data directory
	creds   string // tls certificates
	config  string // configuration directory

	// Log store.
	logs *storage_store.LevelDB_store

	// Create the JWT key used to create the signature
	jwtKey       []byte
	RootPassword string

	// local store.
	store persistence_store.Store

	// client reference...
	persistence_client_ *persistence_client.Persistence_Client
	ldap_client_        *ldap_client.LDAP_Client
	event_client_       *event_client.Event_Client

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
	subscribers     map[string]map[string][]string

	// The http server
	http_server  *http.Server
	https_server *http.Server
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {

	// Here I will initialyse configuration.
	g := new(Globule)

	g.Version = "1.0.0" // Automate version...
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.RootPassword = "adminadmin"

	g.PortHttp = 80   // The default http port
	g.PortHttps = 443 // The default https port number

	g.Name = strings.Replace(Utility.GetExecName(os.Args[0]), ".exe", "", -1)

	g.Protocol = "http"
	g.Domain = "localhost"

	// Set default values.
	g.PortsRange = "10000-10100"

	// set default values.
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"

	g.Services = make(map[string]interface{}, 0)
	g.services = new(sync.Map)

	g.inernalServices = make([]*grpc.Server, 0)

	// Contain the list of ldap syncronization info.
	g.LdapSyncInfos = make(map[string]interface{}, 0)

	// Configuration must be reachable before services initialysation

	// Promometheus metrics for services.
	http.Handle("/metrics", promhttp.Handler())

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Initialyse globular from it configuration file.
	g.config = dir + string(os.PathSeparator) + "config"
	file, err := ioutil.ReadFile(g.config + string(os.PathSeparator) + "config.json")

	// Init the service with the default port address
	if err == nil {
		// get the existing configuration.
		err := json.Unmarshal(file, &g)
		if err != nil {
			log.Println("fail to initialyse the globule configuration")
		}

		// Now I will initialyse sync services map.
		for _, v := range g.Services {
			s := v.(map[string]interface{})
			s_ := new(sync.Map)
			for k_, v_ := range s {
				s_.Store(k_, v_)
			}
			g.setService(s_)
		}

	} else {
		// save the configuration to set the port number.
		portRange := strings.Split(g.PortsRange, "-")
		start := Utility.ToInt(portRange[0])
		g.AdminPort = start + 1
		g.AdminProxy = start + 2
		g.AdminEmail = "admin@globular.app"
		g.RessourcePort = start + 3
		g.RessourceProxy = start + 4

		// services management...
		g.ServicesDiscoveryPort = start + 5
		g.ServicesDiscoveryProxy = start + 6
		g.ServicesRepositoryPort = start + 7
		g.ServicesRepositoryProxy = start + 8
		g.CertificateAuthorityPort = start + 9
		g.CertificateAuthorityProxy = start + 10
		g.LoadBalancingServicePort = start + 11
		g.LoadBalancingServiceProxy = start + 12
	}

	// Prometheus logging informations.
	g.methodsCounterLog = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "globular_methods_counter",
		Help: "Globular services methods usage.",
	},
		[]string{
			"type",
			"method"},
	)

	// Set the function into prometheus.
	prometheus.MustRegister(g.methodsCounterLog)

	// Keep in global var to by http handlers.
	globule = g

	// Set the list of http handler.

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

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

	g.initDirectories()

	return g
}

// Little shortcut to get access to map value in one step.
func setValues(m *sync.Map, values map[string]interface{}) {
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

func getNumericVal(m *sync.Map, k string) float64 {
	v, ok := m.Load(k)
	if !ok {
		return 0.0
	}

	return Utility.ToNumeric(v)
}

func getVal(m *sync.Map, k string) interface{} {
	v, ok := m.Load(k)
	if !ok {
		return nil
	}
	return v
}

func (self *Globule) getServices() []*sync.Map {
	_services_ := make([]*sync.Map, 0)

	// Append services into the array.
	self.services.Range(func(key, s interface{}) bool {
		_services_ = append(_services_, s.(*sync.Map))
		return true
	})

	return _services_

}

func (self *Globule) setService(s *sync.Map) {
	id, _ := s.Load("Id") //service["Id"].(string)
	self.services.Store(id.(string), s)
}

func (self *Globule) getService(id string) *sync.Map {
	s, ok := self.services.Load(id)
	if ok {
		return s.(*sync.Map)
	} else {
		return nil
	}
}

func (self *Globule) deleteService(id string) {
	self.services.Delete(id)
}

func (self *Globule) toMap() map[string]interface{} {
	_map_, _ := Utility.ToMap(self)
	_services_ := make(map[string]interface{})

	self.services.Range(func(key, value interface{}) bool {
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

func (self *Globule) getPortsInUse() []int {
	portsInUse := make([]int, 0)
	// I will test if the port is already taken by e services.
	self.services.Range(func(key, value interface{}) bool {
		m := value.(*sync.Map)
		_, hasProcess := m.Load("Process")
		if hasProcess {
			p, _ := m.Load("Port")
			portsInUse = append(portsInUse, Utility.ToInt(p))
		}
		_, hasProxyProcess := m.Load("ProxyProcess")
		if hasProxyProcess {
			p, _ := m.Load("Proxy")
			portsInUse = append(portsInUse, Utility.ToInt(p))
		}
		return true
	})

	return portsInUse
}

/**
 * test if a given port is avalaible.
 */
func (self *Globule) isPortAvailable(port int) bool {
	portRange := strings.Split(self.PortsRange, "-")
	start := Utility.ToInt(portRange[0]) + 13 // The first 12 addresse are reserver by internal service...
	end := Utility.ToInt(portRange[1])

	if port < start || port > end {
		return false
	}

	portsInUse := self.getPortsInUse()
	for i := 0; i < len(portsInUse); i++ {
		if portsInUse[i] == port {
			return false
		}
	}

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
func (self *Globule) getNextAvailablePort() (int, error) {
	portRange := strings.Split(self.PortsRange, "-")
	start := Utility.ToInt(portRange[0]) + 13 // The first 13 addresse are reserver by internal service...
	end := Utility.ToInt(portRange[1])

	for i := start; i < end; i++ {
		if self.isPortAvailable(i) {
			return i, nil
		}
	}

	return -1, errors.New("No port are available in the range " + self.PortsRange)

}

/**
 * Initialize the server directories config, data, webroot...
 */
func (self *Globule) initDirectories() {

	// DNS info.
	self.DNS = make([]string, 0)
	self.DnsUpdateIpInfos = make([]interface{}, 0)

	// Set the list of discorvery service avalaible...
	self.Discoveries = make([]string, 0)
	self.discorveriesEventHub = make(map[string]*event_client.Event_Client, 0)

	// Set the share service info...
	self.Services = make(map[string]interface{}, 0)

	// Set external map services.
	self.ExternalApplications = make(map[string]ExternalApplication, 0)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	self.path = dir // keep the installation path.

	// if globular is found.
	self.webRoot = dir + string(os.PathSeparator) + "webroot" // The default directory to server.

	// keep the root in global variable for the file handler.
	webRoot = self.webRoot
	Utility.CreateDirIfNotExist(self.webRoot) // Create the directory if it not exist.

	if !Utility.Exists(self.webRoot + string(os.PathSeparator) + "index.html") {

		// in that case I will create a new index.html file.
		ioutil.WriteFile(self.webRoot+string(os.PathSeparator)+"index.html", []byte(
			`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html lang="en">

	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>Title Goes Here</title>
	</head>

	<body>
		<p>Welcome to Globular `+self.Version+`</p>
	</body>

</html>`), 644)
	}

	// Create the directory if is not exist.
	self.data = dir + string(os.PathSeparator) + "data"
	Utility.CreateDirIfNotExist(self.data)

	// Configuration directory
	self.config = dir + string(os.PathSeparator) + "config"
	Utility.CreateDirIfNotExist(self.config)

	// Create the creds directory if it not already exist.
	self.creds = self.config + string(os.PathSeparator) + "tls"
	Utility.CreateDirIfNotExist(self.creds)

	// Initialyse globular from it configuration file.
	file, err := ioutil.ReadFile(self.config + string(os.PathSeparator) + "config.json")

	// Init the service with the default port address
	if err == nil {
		json.Unmarshal(file, &self)
	}

	log.Println("Globular is running!")
}

/**
 * Close the server.
 */
func (self *Globule) KillProcess() {
	// Here I will kill proxies if there are running.
	Utility.KillProcessByName("grpcwebproxy")

	// Kill previous instance of the program...
	for _, s := range self.getServices() {
		_, ok := s.Load("Path")
		if ok {
			name := getStringVal(s, "Path")
			name = name[strings.LastIndex(name, "/")+1:]
			err := Utility.KillProcessByName(name)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

/**
 * Start serving the content.
 */
func (self *Globule) Serve() {

	// Reset previous connections.
	self.store = nil
	self.persistence_client_ = nil
	self.ldap_client_ = nil
	self.event_client_ = nil

	// Open logs db.
	if self.logs == nil {
		self.logs = storage_store.NewLevelDB_store()
		err := self.logs.Open(`{"path":"` + self.data + `", "name":"logs"}`)
		if err != nil {
			log.Println(err)
		}

		// Here it suppose to be only one server instance per computer.
		self.jwtKey = []byte(Utility.RandomUUID())
		err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_key", []byte(self.jwtKey), 0644)
		if err != nil {
			log.Panicln(err)
		}

		// The token that identify the server with other services
		token, _ := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa", self.AdminEmail)
		err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+self.getDomain()+"_token", []byte(token), 0644)
		if err != nil {
			log.Panicln(err)
		}

		// Here I will start the refresh token loop to refresh the server token.
		// the token will be refresh 10 milliseconds before expiration.
		ticker := time.NewTicker((self.SessionTimeout - 10) * time.Millisecond)
		go func() {
			for {
				select {
				case <-ticker.C:
					token, _ := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa", self.AdminEmail)
					err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+self.getDomain()+"_token", []byte(token), 0644)
					if err != nil {
						log.Println(err)
					}
				case <-self.exit:
					break
				}
			}
		}()
	}
	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// First of all I will start external services.
	for externalServiceId, _ := range self.ExternalApplications {
		pid, err := self.startExternalApplication(externalServiceId)
		if err != nil {
			log.Println("fail to start external service: ", externalServiceId, " pid ", pid)
		}
	}

	// I will save the variable in a tmp file to be sure I can get it outside
	ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"GLOBULAR_ROOT", []byte(self.path+":"+Utility.ToString(self.PortHttp)), 0644)

	// set the services.
	self.initServices()

	// Here I will save the server attribute
	self.saveConfig()

	// lisen
	err := self.Listen()
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
func (self *Globule) getDomain() string {
	domain := self.Domain
	if len(self.Name) > 0 && domain != "localhost" {
		domain = /*self.Name + "." +*/ domain
	}
	return domain
}

/**
 * Set the ip for a given domain or sub-domain
 */
func (self *Globule) registerIpToDns() error {

	// Globular DNS is use to create sub-domain.
	// ex: globular1.globular.io here globular.io is the domain and globular1 is
	// the sub-domain. Domain must be manage by dns provider directly, by using
	// the DnsSetA (set ip api call)... see the next part of that function
	// for more information.
	if self.DNS != nil {
		if len(self.DNS) > 0 {
			for i := 0; i < len(self.DNS); i++ {
				dns_client_, err := dns_client.NewDnsService_Client(self.DNS[i], "dns.DnsService")
				if err != nil {
					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = dns_client_.SetA(self.Domain, self.getDomain(), Utility.MyIP(), 60)

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

	for i := 0; i < len(self.DnsUpdateIpInfos); i++ {
		// the api call "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
		setA := self.DnsUpdateIpInfos[i].(map[string]interface{})["SetA"].(string)
		key := self.DnsUpdateIpInfos[i].(map[string]interface{})["Key"].(string)
		secret := self.DnsUpdateIpInfos[i].(map[string]interface{})["Secret"].(string)

		// set the data to the actual ip address.
		data := `[{"data":"` + Utility.MyIP() + `"}]`

		// initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, setA, bytes.NewBuffer([]byte(data)))
		if err != nil {
			log.Println(err)
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

	domains := self.AlternateDomains

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
func (self *Globule) startProxy(s *sync.Map, port int, proxy int) error {
	_, hasProxyProcess := s.Load("ProxyProcess")
	if !hasProxyProcess {
		s.Store("ProxyProcess", -1)
	}
	pid := getIntVal(s, "ProxyProcess")
	if pid != -1 {
		Utility.TerminateProcess(pid, 0)
	}

	// Now I will start the proxy that will be use by javascript client.
	proxyPath := string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
	if string(os.PathSeparator) == "\\" && !strings.HasSuffix(proxyPath, ".exe") {
		proxyPath += ".exe" // in case of windows.
	}

	proxyBackendAddress := self.getDomain() + ":" + strconv.Itoa(port)
	proxyAllowAllOrgins := "true"
	proxyArgs := make([]string, 0)

	// Use in a local network or in test.
	proxyArgs = append(proxyArgs, "--backend_addr="+proxyBackendAddress)
	proxyArgs = append(proxyArgs, "--allow_all_origins="+proxyAllowAllOrgins)
	hasTls := getBoolVal(s, "TLS")
	if hasTls == true {
		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"

		/* Services gRpc backend. */
		proxyArgs = append(proxyArgs, "--backend_tls=true")
		proxyArgs = append(proxyArgs, "--backend_tls_ca_files="+certAuthorityTrust)
		proxyArgs = append(proxyArgs, "--backend_client_tls_cert_file="+self.creds+string(os.PathSeparator)+"client.crt")
		proxyArgs = append(proxyArgs, "--backend_client_tls_key_file="+self.creds+string(os.PathSeparator)+"client.pem")

		/* http2 parameters between the browser and the proxy.*/
		proxyArgs = append(proxyArgs, "--run_http_server=false")
		proxyArgs = append(proxyArgs, "--run_tls_server=true")
		proxyArgs = append(proxyArgs, "--server_http_tls_port="+strconv.Itoa(proxy))

		/* in case of public domain server files **/
		proxyArgs = append(proxyArgs, "--server_tls_key_file="+self.creds+string(os.PathSeparator)+"server.pem")

		proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+self.creds+string(os.PathSeparator)+self.CertificateAuthorityBundle)
		proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.creds+string(os.PathSeparator)+self.Certificate)

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
	proxyProcess := exec.Command(self.path+proxyPath, proxyArgs...)
	proxyProcess.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err := proxyProcess.Start()

	if err != nil {
		return err
	}

	// save service configuration.
	s.Store("ProxyProcess", proxyProcess.Process.Pid)

	return nil
}

/**
 * That function will
 */
func (self *Globule) keepServiceAlive(s *sync.Map) {

	if self.exit_ {
		return
	}

	_, HasKeepAlive := s.Load("KeepAlive")
	if !HasKeepAlive {
		return
	}
	// In case the service must not be kept alive.
	keepAlive := getBoolVal(s, "KeepAlive")
	if !keepAlive {
		return
	}

	pid := getIntVal(s, "Process")
	p, err := os.FindProcess(pid)
	if err != nil {
		return
	}

	// Wait for process to return.
	p.Wait()

	_, _, err = self.startService(s)
	if err != nil {
		return
	}
}

/**
 * Start internal service admin and ressource are use that function.
 */
func (self *Globule) startInternalService(id string, proto string, port int, proxy int, hasTls bool, unaryInterceptor grpc.UnaryServerInterceptor, streamInterceptor grpc.StreamServerInterceptor) (*grpc.Server, error) {
	log.Println("Start internal service ", id)

	s := self.getService(id)
	if s == nil {
		s = new(sync.Map)
	}

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var grpcServer *grpc.Server
	if hasTls {
		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
		certFile := self.creds + string(os.PathSeparator) + "server.crt"
		keyFile := self.creds + string(os.PathSeparator) + "server.pem"

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
	s.Store("Domain", self.getDomain())
	s.Store("Name", id)
	s.Store("Id", id)
	s.Store("Proto", proto)
	s.Store("Port", port)
	s.Store("Proxy", proxy)
	s.Store("TLS", hasTls)
	s.Store("ProxyProcess", -1)

	self.setService(s)

	// save the config.
	self.saveConfig()

	// start the proxy
	err := self.startProxy(s, port, proxy)
	if err != nil {
		return nil, err
	}

	self.inernalServices = append(self.inernalServices, grpcServer)

	return grpcServer, nil
}

/**
 * Stop internal services ressource admin lb...
 */
func (self *Globule) stopInternalServices() {
	for i := 0; i < len(self.inernalServices); i++ {
		self.inernalServices[i].Stop()

	}
}

/**
 * Stop external services.
 */
func (self *Globule) stopServices() {
	// not keep services alive because the server must exist.
	self.exit_ = true

	// Here I will disconnect service update event.
	for id, subscriber := range self.subscribers {
		eventHub := self.discorveriesEventHub[id]
		for channelId, uuids := range subscriber {
			for i := 0; i < len(uuids); i++ {
				eventHub.UnSubscribe(channelId, uuids[i])
			}
		}
		eventHub.Close()
	}

	// stop external service.
	for externalServiceId, _ := range self.ExternalApplications {
		self.stopExternalApplication(externalServiceId)
	}

	// Stop proxy process...
	for _, s := range self.getServices() {
		if s != nil {
			// I will also try to keep a client connection in to communicate with the service.
			self.stopService(getStringVal(s, "Id"))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if self.https_server != nil {
		if err := self.https_server.Shutdown(ctx); err != nil {
			// handle err
			log.Println(err)
		}
		log.Println("stop listen(https) at port ", self.PortHttps)
	}

	if self.http_server != nil {
		if err := self.http_server.Shutdown(ctx); err != nil {
			// handle err
			log.Println(err)
		}
		log.Println("stop listen(http) at port ", self.PortHttp)
	}
}

/**
 * Start services define in the configuration.
 */
func (self *Globule) startService(s *sync.Map) (int, int, error) {

	var err error

	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(getStringVal(s, "Domain")) && root_ != self.path {
		return -1, -1, errors.New("Can not start a distant service localy!")
	}

	// set the domain of the service.
	s.Store("Domain", self.getDomain())
	s.Store("TLS", self.Protocol == "https")

	// if the service already exist.
	_, hasProcess := s.Load("Process")
	if !hasProcess {
		s.Store("Process", -1)
		s.Store("ProxyProcess", -1)
	}

	pid := getIntVal(s, "Process")
	if pid != -1 {
		Utility.TerminateProcess(pid, 0)
	}

	servicePath := getStringVal(s, "Path")
	if getStringVal(s, "Protocol") == "grpc" {
		// I will test if the service is find if not I will try to set path
		// to standard dist directory structure.

		if !Utility.Exists(servicePath) {
			// Here I will set various base on the standard dist directory structure.
			path := self.path + "/services/" + getStringVal(s, "PublisherId") + "/" + getStringVal(s, "Name") + "/" + getStringVal(s, "Version") + "/" + getStringVal(s, "Id")
			execName := servicePath[strings.LastIndex(servicePath, "/")+1:]

			s.Store("Path", path+"/"+execName)
			_, exist := s.Load("Path")
			if !exist {
				return -1, -1, errors.New("Fail to retreive exe path " + servicePath)
			}

			// Now the proto file.
			files, err := Utility.FindFileByName(self.path+"/services/"+getStringVal(s, "PublisherId")+"/"+getStringVal(s, "Name")+"/"+getStringVal(s, "Version"), ".proto")
			if err != nil {
				return -1, -1, errors.New("No prototype file was found")
			}

			s.Store("Proto", files[0])

		}

		hasTls := getBoolVal(s, "TLS")
		log.Println("Has TLS ", hasTls, getStringVal(s, "Name"))
		if hasTls {
			// Set TLS local services configuration here.
			s.Store("CertAuthorityTrust", self.creds+string(os.PathSeparator)+"ca.crt")
			s.Store("CertFile", self.creds+string(os.PathSeparator)+"server.crt")
			s.Store("KeyFile", self.creds+string(os.PathSeparator)+"server.pem")
		} else {
			// not secure services.
			s.Store("CertAuthorityTrust", "")
			s.Store("CertFile", "")
			s.Store("KeyFile", "")
		}

		// Get the next available port.
		port := getIntVal(s, "Port")

		if !self.isPortAvailable(port) {
			port, err = self.getNextAvailablePort()
			if err != nil {
				return -1, -1, err
			}
			s.Store("Port", port)
			self.setService(s)
		}

		// File service need root...
		if getStringVal(s, "Name") == "file.FileService" {
			// Set it root to the globule root.
			s.Store("Root", globule.webRoot)
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
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}

		err = p.Start()
		if err != nil {
			s.Store("State", "fail")
			s.Store("Process", -1)
			log.Println("Fail to start service: ", getStringVal(s, "Name"), " at port ", port, " with error ", err)
			return -1, -1, err
		} else {
			s.Store("Process", p.Process.Pid)
			log.Println("Service ", getStringVal(s, "Name")+":"+getStringVal(s, "Id"), "started with pid:", getIntVal(s, "Process"))
		}

		// save the services in the map.
		go func(s *sync.Map, p *exec.Cmd) {

			// Here I will append the service to the load balancer.
			if port != -1 {
				log.Println("Append ", getStringVal(s, "Name"), " to load balancer.")
				load_info := &lbpb.LoadInfo{
					ServerInfo: &lbpb.ServerInfo{
						Id:     getStringVal(s, "Id"),
						Name:   getStringVal(s, "Name"),
						Domain: getStringVal(s, "Domain"),
						Port:   int32(port),
					},
					Load1:  0, // All service will be initialise with a 0 load.
					Load5:  0,
					Load15: 0,
				}

				self.lb_load_info_channel <- load_info
			}

			s.Store("State", "running")
			self.keepServiceAlive(s)

			// display the message in the console.
			reader := bufio.NewReader(pipe)
			line, err := reader.ReadString('\n')
			for err == nil {
				line, err = reader.ReadString('\n')
				self.logServiceInfo(getStringVal(s, "Name"), line)
			}

			// if the process is not define.
			err = p.Wait() // wait for the program to return
			if err != nil {
				// I will log the program error into the admin logger.
				self.logServiceInfo(getStringVal(s, "Name"), err.Error())
			}

			// Print the error
			if len(errb.String()) > 0 {
				fmt.Println("service", getStringVal(s, "Name"), "err:", errb.String())
			}

			// I will remove the service from the load balancer.
			self.lb_remove_candidate_info_channel <- &lbpb.ServerInfo{
				Id:     getStringVal(s, "Id"),
				Name:   getStringVal(s, "Name"),
				Domain: getStringVal(s, "Domain"),
				Port:   int32(port),
			}

			s.Store("Process", -1)
			s.Store("ProxyProcess", -1)

			self.saveConfig()

		}(s, p)

		// get another port.
		proxy := getIntVal(s, "Proxy")
		if !self.isPortAvailable(proxy) {
			proxy, err = self.getNextAvailablePort()
			if err != nil {
				s.Store("Proxy", -1)

				return -1, -1, err
			}
			s.Store("Proxy", proxy)
			self.setService(s)
		}

		// Start the proxy.
		err = self.startProxy(s, port, proxy)
		if err != nil {
			return -1, -1, err
		}

		// save service config.
		self.saveServiceConfig(s)

		// save it to the config.
		self.saveConfig()
		log.Println("Service "+getStringVal(s, "Name")+":"+getStringVal(s, "Id")+" is up and running at port ", port, " and proxy ", proxy)

	} else if getStringVal(s, "Protocol") == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(getStringVal(s, "Name"), "Globular") {

			p := exec.Command(servicePath, getStringVal(s, "Port"))

			var errb bytes.Buffer
			pipe, _ := p.StdoutPipe()
			p.Stderr = &errb

			// Here I will set the command dir.
			p.Dir = servicePath[:strings.LastIndex(servicePath, string(os.PathSeparator))]
			p.SysProcAttr = &syscall.SysProcAttr{
				CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
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
			s.Store("Process", p.Process.Pid)
			s.Store("State", "running")

			if err == nil {
				go func(s *sync.Map) {

					self.keepServiceAlive(s)

					// display the message in the console.
					reader := bufio.NewReader(pipe)
					line, err := reader.ReadString('\n')
					for err == nil {
						name := getStringVal(s, "Name")
						log.Println(name, ":", line)
						line, err = reader.ReadString('\n')
						self.logServiceInfo(name, line)
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
							self.logServiceInfo(getStringVal(s, "Name"), errb.String())
						}
					}

				}(s)
			}

			// Save configuration stuff.
			self.setService(s)
		}
	}

	if getIntVal(s, "Process") == -1 {
		s.Store("State", "fail")
		self.setService(s)
		err := errors.New("Fail to start process " + getStringVal(s, "Name"))
		return -1, -1, err
	}

	// Return the pid of the service.
	if getIntVal(s, "ProxyProcess") != -1 {
		s.Store("State", "running")
		self.setService(s)
		return getIntVal(s, "Process"), getIntVal(s, "ProxyProcess"), nil
	}

	return getIntVal(s, "Process"), -1, nil
}

/**
 * Init services configuration.
 */
func (self *Globule) initService(s *sync.Map) error {
	_, hasProtocol := s.Load("Protocol")
	if !hasProtocol {
		// internal service dosent has Protocol define.
		return nil
	}

	if getStringVal(s, "Protocol") == "grpc" {
		// The domain must be set in the sever configuration and not change after that.
		hasTls := self.Protocol == "https" //Utility.ToBool(s["TLS"])
		s.Store("TLS", hasTls)             // set the tls...
		if hasTls {
			// Set TLS local services configuration here.
			s.Store("CertAuthorityTrust", self.creds+string(os.PathSeparator)+"ca.crt")
			s.Store("CertFile", self.creds+string(os.PathSeparator)+"server.crt")
			s.Store("KeyFile", self.creds+string(os.PathSeparator)+"server.pem")
		} else {
			// not secure services.
			s.Store("CertAuthorityTrust", "")
			s.Store("CertFile", "")
			s.Store("KeyFile", "")
		}
	}

	// any other http server except this one...
	if !strings.HasPrefix(getStringVal(s, "Name"), "Globular") {
		hasChange := self.saveServiceConfig(s)
		_, hasProcess := s.Load("Process")
		if !hasProcess {
			s.Store("Process", -1)
		}
		if hasChange || getIntVal(s, "Process") == -1 {
			self.setService(s)
			_, _, err := self.startService(s)
			if err != nil {
				return err
			}
			self.saveConfig()
		}
	}

	return nil
}

/**
 * Call once when the server start.
 */
func (self *Globule) initServices() {

	log.Println("Initialyse services")
	log.Println("local ip ", Utility.MyLocalIP())
	log.Println("external ip ", Utility.MyIP())

	// If the protocol is https I will generate the TLS certificate.
	if self.Protocol == "https" {
		// security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.creds)
		if len(self.Certificate) == 0 {
			self.registerIpToDns()
			log.Println(" Now let's encrypts!")
			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/tls/server.csr run
			// I need to remove the gRPC certificate and recreate it.

			Utility.RemoveDirContents(self.creds)

			// recreate the certificates.
			err := security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.creds, self.Country, self.State, self.City, self.Organization, self.AlternateDomains)
			if err != nil {
				log.Println(err)
				return
			}

			err = self.obtainCertificateForCsr()
			if err != nil {
				log.Println(err)
				return
			}
		}

		// Here I will read the certificate
		r, _ := ioutil.ReadFile(self.creds + "/" + self.Certificate)
		block, _ := pem.Decode(r)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err == nil {
			delay := cert.NotAfter.Sub(time.Now()) - time.Duration(15*time.Minute)
			timeout := time.NewTimer(delay)
			go func() {
				<-timeout.C
				//self.Certificate = ""
				//self.restartServices()
			}()
		}
	}

	// That will contain all method path from the proto files.
	self.methods = make([]string, 0)
	self.methods = append(self.methods, "/file.FileService/FileUploadHandler")
	self.actionPermissions = make([]interface{}, 0)

	// Set local action permission
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/DeletePermissions", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 1}}})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/SetRessourceOwner", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/DeleteRessourceOwner", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/admin.AdminService/DeployApplication", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/admin.AdminService/PublishService", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}})

	// It will be execute the first time only...
	configPath := self.config + string(os.PathSeparator) + "config.json"
	if !Utility.Exists(configPath) {
		// Each service contain a file name config.json that describe service.
		// I will keep services info in services map and also it running process.
		basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
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
								s["configPath"] = path

								s_ := new(sync.Map)
								for k, v := range s {
									s_.Store(k, v)
								}

								self.setService(s_)
								self.saveConfig()
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

	// Rescan the proto file and update the role after.
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		path = strings.ReplaceAll(path, "\\", "/")
		if info == nil {
			return nil
		}
		if err == nil && strings.HasSuffix(info.Name(), ".proto") {
			name := info.Name()[0:strings.Index(info.Name(), ".")]
			self.setServiceMethods(name, path)
		}
		return nil
	})

	// Kill previous instance of the program...
	self.KillProcess()

	// Start the load balancer.
	err := self.startLoadBalancingService()
	if err != nil {
		log.Println(err)
	}

	for _, s := range self.getServices() {
		// Remove existing process information.
		s.Store("Process", -1)
		s.Store("ProxyProcess", -1)
		err := self.initService(s)
		if err != nil {
			log.Println(err)
		}
	}
}

// That function resolve import path.
func resolveImportPath(path string, importPath string) (string, error) {

	// firt of all i will keep only the path part of the import...
	startIndex := strings.Index(importPath, `'@`) + 1
	endIndex := strings.LastIndex(importPath, `'`)
	importPath_ := importPath[startIndex:endIndex]

	filepath.Walk(webRoot+path[0:strings.Index(path, "/")],
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

	importPath_ = strings.Replace(importPath_, strings.Replace(webRoot, "\\", "/", -1), "", -1)

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
	importPath_ = strings.Replace(importPath_, webRoot, "", 1)

	// remove the root path part and the leading / caracter.
	return importPath_, nil
}

/**
 * Start the monitoring service with prometheus.
 */
func (self *Globule) startMonitoring() error {
	if self.getConfig()["Services"].(map[string]interface{})["monitoring.MonitoringService"] == nil {
		return errors.New("No monitoring service configuration was found on that server!")
	}

	var err error
	s := self.getConfig()["Services"].(map[string]interface{})["monitoring.MonitoringService"].(map[string]interface{})

	// Cast-it to the persistence client.
	m, err := monitoring_client.NewMonitoringService_Client(s["Domain"].(string)+":"+Utility.ToString(self.PortHttp), "monitoring.MonitoringService")
	if err != nil {
		return err
	}

	// Here I will start promethus.
	dataPath := self.data + string(os.PathSeparator) + "prometheus-data"
	Utility.CreateDirIfNotExist(dataPath)
	if !Utility.Exists(self.config + string(os.PathSeparator) + "prometheus.yml") {
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
    - targets: ['localhost:10000']
    
  - job_name: 'node_exporter_metrics'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:9100']
    
  - job_name: 'plc_exporter'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:2112']
`
		err := ioutil.WriteFile(self.config+string(os.PathSeparator)+"prometheus.yml", []byte(config), 0644)
		if err != nil {
			return err
		}
	}

	if !Utility.Exists(self.config + string(os.PathSeparator) + "alertmanager.yml") {
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
		err := ioutil.WriteFile(self.config+string(os.PathSeparator)+"alertmanager.yml", []byte(config), 0644)
		if err != nil {
			return err
		}
	}

	prometheus := exec.Command("prometheus", "--web.listen-address", "0.0.0.0:9090", "--config.file", self.config+string(os.PathSeparator)+"prometheus.yml", "--storage.tsdb.path", dataPath)
	err = prometheus.Start()
	prometheus.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = s["Process"].(*exec.Cmd).Start()
	if err != nil {
		log.Println("fail to start monitoring with prometheus", err)
		return err
	}

	alertmanager := exec.Command("alertmanager", "--config.file", self.config+string(os.PathSeparator)+"alertmanager.yml")
	alertmanager.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = alertmanager.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	node_exporter := exec.Command("node_exporter")
	node_exporter.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = node_exporter.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	// Here I will create a new connection.
	err = m.CreateConnection("local_ressource", "localhost", 0, 9090)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Connection with local persistence grpc service
 */
func (self *Globule) getPersistenceSaConnection() (*persistence_client.Persistence_Client, error) {
	// That service made user of persistence service.
	if self.persistence_client_ != nil {
		return self.persistence_client_, nil
	}

	configs := self.getServiceConfigByName("persistence.PersistenceService")
	if len(configs) == 0 {
		err := errors.New("No persistence service configuration was found on that server!")
		return nil, err
	}

	var err error
	s := configs[0]

	// Cast-it to the persistence client.
	self.persistence_client_, err = persistence_client.NewPersistenceService_Client(s["Domain"].(string)+":"+Utility.ToString(self.PortHttp), s["Id"].(string))
	if err != nil {
		return nil, err
	}

	// Connect to the database here.
	err = self.persistence_client_.CreateConnection("local_ressource", "local_ressource", "0.0.0.0", 27017, 0, "sa", self.RootPassword, 5000, "", false)
	if err != nil {
		return nil, err
	}

	return self.persistence_client_, nil
}

/**
 * Connection to mongo db local store.
 */
func (self *Globule) getPersistenceStore() (persistence_store.Store, error) {
	// That service made user of persistence service.
	if self.store == nil {
		self.store = new(persistence_store.MongoStore)
		err := self.store.Connect("local_ressource", "0.0.0.0", 27017, "sa", self.RootPassword, "local_ressource", 5000, "")
		if err != nil {
			return nil, err
		}
	}

	return self.store, nil
}

/** Stop mongod process **/
func (self *Globule) stopMongod() error {
	closeCmd := exec.Command("mongo", "--eval", "db=db.getSiblingDB('admin');db.adminCommand( { shutdown: 1 } );")
	err := closeCmd.Run()
	time.Sleep(1 * time.Second)
	return err
}

func (self *Globule) waitForMongo(timeout int, withAuth bool) error {

	time.Sleep(1 * time.Second)
	args := make([]string, 0)
	if withAuth == true {
		args = append(args, "-u")
		args = append(args, "sa")
		args = append(args, "-p")
		args = append(args, self.RootPassword)
		args = append(args, "--authenticationDatabase")
		args = append(args, "admin")
	}
	args = append(args, "--eval")
	args = append(args, "db=db.getSiblingDB('admin');db.getMongo().getDBNames()")

	script := exec.Command("mongo", args...)
	err := script.Run()
	if err != nil {
		log.Println("wait for mongo...", timeout, "s")
		if timeout == 0 {
			return errors.New("mongod is not responding!")
		}
		// call again.
		timeout -= 1
		return self.waitForMongo(timeout, withAuth)
	}
	return nil
}

func (self *Globule) getLdapClient() (*ldap_client.LDAP_Client, error) {

	configs := self.getServiceConfigByName("ldap.LdapService")
	if len(configs) == 0 {
		return nil, errors.New("No event service was configure on that globule!")
	}

	var err error

	s := configs[0]

	if self.ldap_client_ == nil {
		self.ldap_client_, err = ldap_client.NewLdapService_Client(s["Domain"].(string)+":"+Utility.ToString(self.PortHttp), "ldap.LdapService")
	}

	return self.ldap_client_, err
}

/**
 * Get access to the event services.
 */
func (self *Globule) getEventHub() (*event_client.Event_Client, error) {

	configs := self.getServiceConfigByName("event.EventService")
	if len(configs) == 0 {
		return nil, errors.New("No event service was configure on that globule!")
	}

	s := configs[0]

	var err error
	if self.event_client_ == nil {
		log.Println("Create connection to event hub ", s["Domain"].(string))
		self.event_client_, err = event_client.NewEventService_Client(s["Domain"].(string)+":"+Utility.ToString(self.PortHttp), s["Id"].(string))
	}

	return self.event_client_, err
}

func (self *Globule) GetAbsolutePath(path string) string {

	path = strings.ReplaceAll(path, "\\", "/")
	if strings.HasSuffix(path, "/") {
		path = path[0 : len(path)-2]
	}

	if len(path) > 1 {
		if strings.HasPrefix(path, "/") {
			path = strings.ReplaceAll(self.webRoot, "\\", "/") + path
		} else {
			path = strings.ReplaceAll(self.webRoot, "\\", "/") + "/" + path
		}
	} else {
		path = strings.ReplaceAll(self.webRoot, "\\", "/")
	}

	return path

}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() error {

	// Here I will subscribe to event service to keep then up to date.
	self.subscribers = self.keepServicesUpToDate()

	// Start internal services.

	// Admin service
	err := self.startAdminService()
	if err != nil {
		return err
	}

	// Ressource service
	err = self.startRessourceService()
	if err != nil {
		return err
	}

	// Directorie service
	err = self.startDiscoveryService()
	if err != nil {
		return err
	}

	// Repository service
	err = self.startRepositoryService()
	if err != nil {
		return err
	}

	// Certificate autority service.
	err = self.startCertificateAuthorityService()
	if err != nil {
		return err
	}

	// Must be started before other services.
	go func() {
		// local - non secure connection.
		self.http_server = &http.Server{
			Addr: ":" + strconv.Itoa(self.PortHttp),
		}
		err = self.http_server.ListenAndServe()
	}()

	// Here I will make a signal hook to interrupt to exit cleanly.
	// handle the Interrupt
	// set the register sa user.
	self.registerSa()

	// Start the monitoring service with prometheus.
	//self.startMonitoring()

	// Start the http server.
	if self.Protocol == "https" {

		// if no certificates are specified I will try to get one from let's encrypts.
		// Start https server.
		self.https_server = &http.Server{
			Addr: ":" + strconv.Itoa(self.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: self.getDomain(),
			},
		}

		// get the value from the configuration files.
		go func() {
			err = self.https_server.ListenAndServeTLS(self.creds+string(os.PathSeparator)+self.Certificate, self.creds+string(os.PathSeparator)+"server.pem")
		}()
	}

	log.Println("Globular is running!")
	return err
}

///////// Implement the User Interface. ////////////

/**
 * Return the admin email.
 */
func (self *Globule) GetEmail() string {
	return self.AdminEmail
}

/**
 * Use the time of registration... Nil other wise.
 */
func (self *Globule) GetRegistration() *registration.Resource {
	return self.registration
}

/**
 * I will reuse the client public key here as key instead of generate another key
 * and manage it...
 */
func (self *Globule) GetPrivateKey() crypto.PrivateKey {
	keyPem, err := ioutil.ReadFile(self.creds + string(os.PathSeparator) + "client.pem")
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

///////// End of Implement the User Interface. ////////////

/**
 * That function work correctly, but the DNS fail time to time to give the
 * IP address that result in a fail request... The DNS must be fix!
 */
func (self *Globule) obtainCertificateForCsr() error {

	config := lego.NewConfig(self)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", strconv.Itoa(self.PortHttp)))
	if err != nil {
		log.Fatal(err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	self.registration = reg
	if err != nil {
		return err
	}

	csrPem, err := ioutil.ReadFile(self.creds + string(os.PathSeparator) + "server.csr")
	if err != nil {
		return err
	}

	csrBlock, _ := pem.Decode(csrPem)
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return err
	}

	cert_rqst := certificate.ObtainForCSRRequest{
		CSR:    csr,
		Bundle: true,
	}

	resource, err := client.Certificate.ObtainForCSR(cert_rqst)
	if err != nil {
		return err
	}

	// Keep certificates url in the config.
	self.CertURL = resource.CertURL
	self.CertStableURL = resource.CertStableURL

	// Set the certificates paths...
	self.Certificate = self.getDomain() + ".crt"
	self.CertificateAuthorityBundle = self.getDomain() + ".issuer.crt"

	// Save the certificate in the cerst folder.
	ioutil.WriteFile(self.creds+string(os.PathSeparator)+self.Certificate, resource.Certificate, 0400)
	ioutil.WriteFile(self.creds+string(os.PathSeparator)+self.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	self.saveConfig()

	return nil
}
