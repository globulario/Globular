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
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	//	"github.com/davecourtois/Globular/services/golang/globular_client"
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

	"sync"

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

	// The list of install services.
	Services map[string]interface{}
	services chan map[string]interface{}

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

	g.PortHttp = 8080  // The default http port
	g.PortHttps = 8181 // The default https port number

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
	// open the channel to get services map.
	g.services = make(chan map[string]interface{}, 0)

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

	return g
}

func (self *Globule) toMap() map[string]interface{} {
	action := make(map[string]interface{})
	action["result"] = make(chan map[string]interface{})
	action["name"] = "toMap"
	self.services <- action
	return <-action["result"].(chan map[string]interface{})
}

func (self *Globule) getServices() []map[string]interface{} {
	action := make(map[string]interface{})
	action["result"] = make(chan []map[string]interface{})
	action["name"] = "getServices"
	self.services <- action
	return <-action["result"].(chan []map[string]interface{})
}

func (self *Globule) setService(service map[string]interface{}) {
	action := make(map[string]interface{})
	action["result"] = make(chan bool)
	action["service"] = service
	action["name"] = "setService"
	self.services <- action
	<-action["result"].(chan bool)
	return
}

func (self *Globule) getService(id string) map[string]interface{} {
	action := make(map[string]interface{})
	action["id"] = id
	action["result"] = make(chan map[string]interface{})
	action["name"] = "getService"
	self.services <- action
	return <-action["result"].(chan map[string]interface{})
}

func (self *Globule) deleteService(id string) {
	action := make(map[string]interface{})
	action["id"] = id
	action["result"] = make(chan bool)
	action["name"] = "deleteService"
	self.services <- action
	<-action["result"].(chan bool)
}

func (self *Globule) getPortsInUse() []int {
	action := make(map[string]interface{})
	action["result"] = make(chan []int)
	action["name"] = "getPortsInUse"
	self.services <- action
	return <-action["result"].(chan []int)
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

	// Here I will keep values in a synmap.
	services := new(sync.Map)

	for k, v := range self.Services {
		services.Store(k, v)
	}

	go func() {
		for {
			select {
			case action := <-self.services:
				if action["name"] == "getServices" {
					_services_ := make([]map[string]interface{}, 0)

					// Append services into the array.
					services.Range(func(key, value interface{}) bool {
						s := make(map[string]interface{})
						for k, v := range value.(map[string]interface{}) {

							s[k] = v

						}
						_services_ = append(_services_, s)
						return true
					})

					action["result"].(chan []map[string]interface{}) <- _services_

				} else if action["name"] == "getService" {

					id := action["id"].(string)
					value, ok := services.Load(id)
					if ok {
						s := make(map[string]interface{})
						for k, v := range value.(map[string]interface{}) {
							s[k] = v
						}
						action["result"].(chan map[string]interface{}) <- s
					} else {
						action["result"].(chan map[string]interface{}) <- nil
					}

				} else if action["name"] == "deleteService" {

					id := action["id"].(string)
					services.Delete(id)
					action["result"].(chan bool) <- true

				} else if action["name"] == "setService" {

					id := action["service"].(map[string]interface{})["Id"].(string)
					services.Store(id, action["service"])
					action["result"].(chan bool) <- true

				} else if action["name"] == "toMap" {

					_map_, _ := Utility.ToMap(self)
					_services_ := make(map[string]interface{})

					services.Range(func(key, value interface{}) bool {
						s := make(map[string]interface{})
						for k, v := range value.(map[string]interface{}) {
							if k != "Process" && k != "ProxyProcess" {
								s[k] = v
							}
						}
						_services_[key.(string)] = s
						return true
					})
					_map_["Services"] = _services_
					action["result"].(chan map[string]interface{}) <- _map_

				} else if action["name"] == "getPortsInUse" {

					portsInUse := make([]int, 0)
					// I will test if the port is already taken by e services.
					services.Range(func(key, value interface{}) bool {
						s := value.(map[string]interface{})
						if s["Process"] != nil {
							portsInUse = append(portsInUse, Utility.ToInt(s["Port"]))
						}
						if s["ProxyProcess"] != nil {
							portsInUse = append(portsInUse, Utility.ToInt(s["Proxy"]))
						}
						return true
					})

					action["result"].(chan []int) <- portsInUse

				}
			case <-self.exit:
				return
			}
		}
	}()

}

/**
 * Close the server.
 */
func (self *Globule) KillProcess() {
	// Here I will kill proxies if there are running.
	Utility.KillProcessByName("grpcwebproxy")

	// Kill previous instance of the program...
	for _, s := range self.getServices() {
		if s["Path"] != nil {
			name := s["Path"].(string)
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

	// initialyse directories.
	self.initDirectories()

	// Open logs db.
	self.logs = storage_store.NewLevelDB_store()
	err := self.logs.Open(`{"path":"` + self.data + `", "name":"logs"}`)
	if err != nil {
		log.Panicln(err)
	}

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

	// Handle the get ca certificate function
	http.HandleFunc("/get_ca_certificate", getCaCertificateHanldler)

	// Return the san server configuration.
	http.HandleFunc("/get_san_conf", getSanConfigurationHandler)

	// Handle the signing certificate function.
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

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
			}
		}
	}()

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
	err = self.Listen()
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

	return nil
}

/**
 * Start the grpc proxy.
 */
func (self *Globule) startProxy(id string, port int, proxy int) error {
	srv := self.getService(id)
	if srv["ProxyProcess"] != nil {
		Utility.TerminateProcess(srv["ProxyProcess"].(*exec.Cmd).Process.Pid)
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
	hasTls := Utility.ToBool(srv["TLS"])
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
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err := proxyProcess.Start()

	if err != nil {
		return err
	}

	// save service configuration.
	srv["ProxyProcess"] = proxyProcess
	srv["Id"] = id
	self.setService(srv)

	return nil
}

/**
 * That function will
 */
func (self *Globule) keepServiceAlive(s map[string]interface{}) {
	if self.exit_ {
		return
	}

	if s["KeepAlive"] == nil {
		return
	}
	// In case the service must not be kept alive.
	keepAlive := Utility.ToBool(s["KeepAlive"])
	if !keepAlive {
		return
	}

	s["Process"].(*exec.Cmd).Wait()

	_, _, err := self.startService(s)
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
		s = make(map[string]interface{}, 0)
	}

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var grpcServer *grpc.Server
	if hasTls {
		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
		certFile := self.creds + string(os.PathSeparator) + "server.crt"
		keyFile := self.creds + string(os.PathSeparator) + "server.pem"

		s["CertFile"] = certFile
		s["KeyFile"] = keyFile
		s["CertAuthorityTrust"] = certAuthorityTrust

		// Create the TLS credentials
		creds := credentials.NewTLS(globular.GetTLSConfig(keyFile, certFile, certAuthorityTrust))

		// Create the gRPC server with the credentials
		opts := []grpc.ServerOption{grpc.Creds(creds),
			grpc.UnaryInterceptor(unaryInterceptor),
			grpc.StreamInterceptor(streamInterceptor)}

		// Create the gRPC server with the credentials
		grpcServer = grpc.NewServer(opts...)

	} else {
		s["CertFile"] = ""
		s["KeyFile"] = ""
		s["CertAuthorityTrust"] = ""

		grpcServer = grpc.NewServer([]grpc.ServerOption{
			grpc.UnaryInterceptor(unaryInterceptor),
			grpc.StreamInterceptor(streamInterceptor)}...)
	}

	reflection.Register(grpcServer)

	// Here I will create the service configuration object.
	s["Domain"] = self.getDomain()
	s["Name"] = id
	s["Id"] = id
	s["Proto"] = proto
	s["Port"] = port
	s["Proxy"] = proxy
	s["TLS"] = hasTls

	self.setService(s)

	// save the config.
	self.saveConfig()

	// start the proxy
	err := self.startProxy(id, port, proxy)
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
		self.inernalServices[i].GracefulStop()
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
			self.stopService(s["Id"].(string))
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
func (self *Globule) startService(s map[string]interface{}) (int, int, error) {
	var err error

	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	root_ := string(root)[0:strings.Index(string(root), ":")]

	if !Utility.IsLocal(s["Domain"].(string)) && root_ != self.path {
		return -1, -1, errors.New("Can not start a distant service localy!")
	}

	// set the domain of the service.
	s["Domain"] = self.getDomain()

	// if the service already exist.
	srv := self.getService(s["Id"].(string))
	if srv != nil {
		if srv["Process"] != nil {
			if reflect.TypeOf(srv["Process"]).String() == "*exec.Cmd" {
				if srv["Process"].(*exec.Cmd).Process != nil {
					Utility.TerminateProcess(srv["Process"].(*exec.Cmd).Process.Pid)
				}
			}
		}
	}

	servicePath := s["Path"].(string)
	if s["Protocol"].(string) == "grpc" {
		hasTls := Utility.ToBool(s["TLS"])
		if hasTls {
			// Set TLS local services configuration here.
			s["CertAuthorityTrust"] = self.creds + string(os.PathSeparator) + "ca.crt"
			s["CertFile"] = self.creds + string(os.PathSeparator) + "server.crt"
			s["KeyFile"] = self.creds + string(os.PathSeparator) + "server.pem"
		} else {
			// not secure services.
			s["CertAuthorityTrust"] = ""
			s["CertFile"] = ""
			s["KeyFile"] = ""
		}

		if !Utility.Exists(servicePath) {
			log.Println("Fail to retreive exe path ", servicePath)
			// Here the service was not retreive so I will try to fix the path...
			root := strings.ReplaceAll(self.path, "\\", "/")
			if strings.Index(servicePath, "services") != -1 {
				// set the service path
				servicePath = root + servicePath[strings.Index(servicePath, "/services"):]
				s["Path"] = servicePath

				// set the proto path
				protoPath := s["Proto"].(string)
				protoPath = root + protoPath[strings.Index(protoPath, "/services"):]
				s["Proto"] = protoPath

				// set the key path
				if s["KeyFile"] != nil {
					if len(s["KeyFile"].(string)) > 0 {
						path := s["KeyFile"].(string)
						if !Utility.Exists(path) {
							path = root + path[strings.Index(path, "/services"):]
							s["KeyFile"] = path
						}
					}
				}

				// set the certificate path
				if s["CertFile"] != nil {
					if len(s["CertFile"].(string)) > 0 {
						path := s["CertFile"].(string)
						if !Utility.Exists(path) {
							path = root + path[strings.Index(path, "/services"):]
							s["CertFile"] = path
						}
					}
				}

				// set the ca path
				if s["CertAuthorityTrust"] != nil {
					if len(s["CertAuthorityTrust"].(string)) > 0 {
						path := s["CertAuthorityTrust"].(string)
						if !Utility.Exists(path) {
							path = root + path[strings.Index(path, "/services"):]
							s["CertAuthorityTrust"] = path
						}
					}
				}

				// here I will keep the configuration path in the global configuration.
				if s["configPath"] != nil {
					if len(s["configPath"].(string)) > 0 {
						configPath := s["configPath"].(string)
						if !Utility.Exists(configPath) {
							configPath = root + configPath[strings.Index(configPath, "/services"):]
							s["configPath"] = configPath
						}
					}
				}
				self.saveServiceConfig(s)

			} else {
				return -1, -1, errors.New("no service found at path " + servicePath)
			}
		}

		// Get the next available port.
		port := Utility.ToInt(s["Port"])
		if !self.isPortAvailable(port) {
			port, err = self.getNextAvailablePort()
			if err != nil {
				return -1, -1, err
			}
		}
		s["Port"] = port

		// File service need root...
		if s["Name"].(string) == "file.FileService" {
			s["Root"] = globule.webRoot
			s["Process"] = exec.Command(servicePath, Utility.ToString(port), globule.webRoot)
		} else {
			s["Process"] = exec.Command(servicePath, Utility.ToString(port))
		}

		var errb bytes.Buffer
		pipe, _ := s["Process"].(*exec.Cmd).StdoutPipe()
		s["Process"].(*exec.Cmd).Stderr = &errb

		// Here I will set the command dir.
		s["Process"].(*exec.Cmd).Dir = servicePath[:strings.LastIndex(servicePath, "/")]
		s["Process"].(*exec.Cmd).SysProcAttr = &syscall.SysProcAttr{
			//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}

		err = s["Process"].(*exec.Cmd).Start()
		if err != nil {
			s["State"] = "fail"
			log.Println("Fail to start service: ", s["Name"].(string), " at port ", port, " with error ", err)
			return -1, -1, err
		}

		go func() {

			// Here I will append the service to the load balancer.
			if port != -1 {
				log.Println("Append ", s["Name"].(string), " to load balancer.")
				load_info := &lbpb.LoadInfo{
					ServerInfo: &lbpb.ServerInfo{
						Id:     s["Id"].(string),
						Name:   s["Name"].(string),
						Domain: s["Domain"].(string),
						Port:   int32(port),
					},
					Load1:  0, // All service will be initialise with a 0 load.
					Load5:  0,
					Load15: 0,
				}

				self.lb_load_info_channel <- load_info
			}
			s["State"] = "running"

			self.keepServiceAlive(s)

			// display the message in the console.
			reader := bufio.NewReader(pipe)
			line, err := reader.ReadString('\n')
			for err == nil {
				line, err = reader.ReadString('\n')
				self.logServiceInfo(s["Name"].(string), line)
			}

			// if the process is not define.
			if s["Process"] == nil {
				log.Println("No process found for service", s["Name"].(string))
				return
			}

			err = s["Process"].(*exec.Cmd).Wait() // wait for the program to return

			if err != nil {
				// I will log the program error into the admin logger.
				self.logServiceInfo(s["Name"].(string), err.Error())

			}

			// Print the error
			if len(errb.String()) > 0 {
				fmt.Println("service", s["Name"].(string), "err:", errb.String())
			}

			// I will remove the service from the load balancer.
			self.lb_remove_candidate_info_channel <- &lbpb.ServerInfo{
				Id:     s["Id"].(string),
				Name:   s["Name"].(string),
				Domain: s["Domain"].(string),
				Port:   int32(port),
			}

		}()

		// get another port.
		proxy := port + 1
		if !self.isPortAvailable(proxy) {
			proxy, err = self.getNextAvailablePort()
			if err != nil {
				return -1, -1, err
			}
		}

		// Start the proxy.
		err = self.startProxy(s["Id"].(string), port, proxy)
		if err != nil {
			return -1, -1, err
		}

		// Save configuration stuff.
		s["Proxy"] = proxy

		self.setService(s)

		// get back the service info with the proxy process in it
		s = self.getService(s["Id"].(string))

		// save it to the config.
		self.saveConfig()
		log.Println("Service "+s["Name"].(string)+":"+s["Id"].(string)+" is up and running at port ", port, " and proxy ", proxy)

	} else if s["Protocol"].(string) == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(s["Name"].(string), "Globular") {

			s["Process"] = exec.Command(servicePath, Utility.ToString(s["Port"]))

			var errb bytes.Buffer
			pipe, _ := s["Process"].(*exec.Cmd).StdoutPipe()
			s["Process"].(*exec.Cmd).Stderr = &errb

			// Here I will set the command dir.
			s["Process"].(*exec.Cmd).Dir = servicePath[:strings.LastIndex(servicePath, string(os.PathSeparator))]
			err = s["Process"].(*exec.Cmd).Start()
			s["Process"].(*exec.Cmd).SysProcAttr = &syscall.SysProcAttr{
				//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			}

			err = s["Process"].(*exec.Cmd).Start()
			if err != nil {
				// The process already exist so I will not throw an error and I will use existing process instead. I will make the
				if err.Error() != "exec: already started" {
					s["State"] = "fail"
					log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
					return -1, -1, err
				}
			}

			s["State"] = "running"
			if err == nil {
				go func() {

					self.keepServiceAlive(s)

					// display the message in the console.
					reader := bufio.NewReader(pipe)
					line, err := reader.ReadString('\n')
					for err == nil {
						log.Println(s["Name"].(string), ":", line)
						line, err = reader.ReadString('\n')
						self.logServiceInfo(s["Name"].(string), line)
					}

					// if the process is not define.
					if s["Process"] == nil {
						log.Println("No process found for service", s["Name"].(string))
					}

					err = s["Process"].(*exec.Cmd).Wait() // wait for the program to resturn

					if err != nil {
						// I will log the program error into the admin logger.
						self.logServiceInfo(s["Name"].(string), errb.String())
					}
				}()
			}

			// Save configuration stuff.
			self.setService(s)
		}
	}

	if s["Process"].(*exec.Cmd).Process == nil {
		s["State"] = "fail"
		err := errors.New("Fail to start process " + s["Name"].(string))
		return -1, -1, err
	}

	// Return the pid of the service.
	if s["ProxyProcess"] != nil {
		return s["Process"].(*exec.Cmd).Process.Pid, s["ProxyProcess"].(*exec.Cmd).Process.Pid, nil
	}

	return s["Process"].(*exec.Cmd).Process.Pid, -1, nil
}

/**
 * Init services configuration.
 */
func (self *Globule) initService(s map[string]interface{}) error {
	if s["Protocol"] == nil {
		// internal service dosent has Protocol define.
		return nil
	}

	if s["Protocol"].(string) == "grpc" {
		// The domain must be set in the sever configuration and not change after that.
		hasTls := Utility.ToBool(s["TLS"])
		if hasTls {
			// Set TLS local services configuration here.
			s["CertAuthorityTrust"] = self.creds + string(os.PathSeparator) + "ca.crt"
			s["CertFile"] = self.creds + string(os.PathSeparator) + "server.crt"
			s["KeyFile"] = self.creds + string(os.PathSeparator) + "server.pem"
		} else {
			// not secure services.
			s["CertAuthorityTrust"] = ""
			s["CertFile"] = ""
			s["KeyFile"] = ""
		}
	}

	// any other http server except this one...
	if !strings.HasPrefix(s["Name"].(string), "Globular") {
		hasChange := self.saveServiceConfig(s)
		if hasChange || s["Process"] == nil {
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
		} else {
			security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.creds, self.Country, self.State, self.City, self.Organization, self.AlternateDomains)
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
								self.setService(s)
								s["configPath"] = path
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
		delete(s, "Process")
		delete(s, "ProxyProcess")
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
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = s["Process"].(*exec.Cmd).Start()
	if err != nil {
		log.Println("fail to start monitoring with prometheus", err)
		return err
	}

	alertmanager := exec.Command("alertmanager", "--config.file", self.config+string(os.PathSeparator)+"alertmanager.yml")
	alertmanager.SysProcAttr = &syscall.SysProcAttr{
		//CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = alertmanager.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
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

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

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
	self.startMonitoring()

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
