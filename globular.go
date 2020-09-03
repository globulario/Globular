package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"

	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecourtois/Globular/lb/lbpb"

	"github.com/prometheus/client_golang/prometheus"

	// Interceptor for authentication, event, log...
	"github.com/davecourtois/Globular/Interceptors"

	// Client services.
	"crypto"

	"github.com/davecourtois/Globular/dns/dns_client"
	"github.com/davecourtois/Globular/event/event_client"
	"github.com/davecourtois/Globular/ldap/ldap_client"
	"github.com/davecourtois/Globular/monitoring/monitoring_client"

	"github.com/davecourtois/Globular/storage/storage_store"
	"github.com/davecourtois/Utility"
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"

	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/persistence/persistence_store"
	"github.com/davecourtois/Globular/security"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global variable.
var (
	root    string
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
	ConfigurationPort         int    // The port use to get the server configuration.
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
	Services      map[string]interface{}
	LdapSyncInfos map[string]interface{} // Contain LdapSyncInfos...

	// List of application need to be start by the server.
	ExternalApplications map[string]ExternalApplication

	Domain                     string // The domain of the globule.
	CertExpirationDelay        int
	CertPassword               string
	Certificate                string
	CertificateAuthorityBundle string
	CertURL                    string
	CertStableURL              string
	Version                    string
	Platform                   string
	SessionTimeout             time.Duration

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep service up to date.

	// DNS stuff.
	DNS []string // Domain name server use to located the server.

	// the api call "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
	DnsSetA string

	// see https://developer.godaddy.com for more detail.
	DnsKey     string
	DnsSecrect string

	discorveriesEventHub map[string]*event_client.Event_Client

	// The list of method supported by this server.
	methods []string

	// Array of action permissions
	actionPermissions []interface{}

	// The prometheus logging informations.
	methodsCounterLog *prometheus.CounterVec

	// prometheus.CounterVec

	// Directories.
	path     string // The path of the exec...
	webRoot  string // The root of the http file server.
	data     string // the data directory
	creds    string // gRpc certificate
	certs    string // https certificates
	config   string // configuration directory
	lastPort int    // The last attributed port number.

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

	g.PortHttp = 8080           // The default http port
	g.PortHttps = 8181          // The default https port number
	g.ConfigurationPort = 10000 // The default configuration port.

	g.Name = strings.Replace(Utility.GetExecName(os.Args[0]), ".exe", "", -1)

	g.Protocol = "http"
	g.Domain = "localhost"

	// Set default values.
	g.PortsRange = "10001-10100"

	g.AdminPort = 10001
	g.AdminProxy = 10002
	g.AdminEmail = "admin@globular.app"
	g.RessourcePort = 10003
	g.RessourceProxy = 10004

	// services management...
	g.ServicesDiscoveryPort = 10005
	g.ServicesDiscoveryProxy = 10006
	g.ServicesRepositoryPort = 10007
	g.ServicesRepositoryProxy = 10008
	g.CertificateAuthorityPort = 10009
	g.CertificateAuthorityProxy = 10010
	g.LoadBalancingServicePort = 10011
	g.LoadBalancingServiceProxy = 10012

	// set default values.
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"

	// Contain the list of ldap syncronization info.
	g.LdapSyncInfos = make(map[string]interface{}, 0)

	// Configuration must be reachable before services initialysation
	go func() {
		// Promometheus metrics for services.
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":"+Utility.ToString(g.ConfigurationPort), nil)
	}()

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Initialyse globular from it configuration file.
	g.config = dir + string(os.PathSeparator) + "config"
	file, err := ioutil.ReadFile(g.config + string(os.PathSeparator) + "config.json")
	// Init the service with the default port address
	if err == nil {
		json.Unmarshal(file, &g)

		// Now here I will set the services ports...
		portsRange := strings.Split(g.PortsRange, "-")

		start := Utility.ToInt(portsRange[0])
		// end :=  Utility.ToInt(portsRange[0])

		g.AdminPort = start
		g.AdminProxy = start + 1
		g.RessourcePort = start + 2
		g.RessourceProxy = start + 3

		// services management...
		g.ServicesDiscoveryPort = start + 4
		g.ServicesDiscoveryProxy = start + 5
		g.ServicesRepositoryPort = start + 6
		g.ServicesRepositoryProxy = start + 7
		g.CertificateAuthorityPort = start + 8
		g.CertificateAuthorityProxy = start + 9
		g.LoadBalancingServicePort = start + 10
		g.LoadBalancingServiceProxy = start + 11

		// save the configuration to set the port number.
		g.saveConfig()

		// Keep the last attributed port number.
		g.lastPort = start + 11

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

/**
 * Initialize the server directories config, data, webroot...
 */
func (self *Globule) initDirectories() {
	// Intialise directories
	self.DNS = make([]string, 0)

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
	root = self.webRoot
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
	self.creds = self.config + string(os.PathSeparator) + "grpc_tls"
	Utility.CreateDirIfNotExist(self.creds)

	// https certificates.
	self.certs = self.config + string(os.PathSeparator) + "http_tls"
	Utility.CreateDirIfNotExist(self.certs)

	// Initialyse globular from it configuration file.
	file, err := ioutil.ReadFile(self.config + string(os.PathSeparator) + "config.json")

	// Init the service with the default port address
	if err == nil {
		json.Unmarshal(file, &self)
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

	// Handle the signing certificate function.
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

	// Here I will kill proxies if there are running.
	Utility.KillProcessByName("grpcwebproxy")

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
	ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"GLOBULAR_ROOT", []byte(self.path+":"+Utility.ToString(self.ConfigurationPort)), 0644)

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
	domain = strings.ToLower(domain)
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
				dns_client, err := dns_client.NewDns_Client(self.DNS[i], "dns.DnsService")
				if err != nil {
					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = dns_client.SetA(self.Domain, self.getDomain(), Utility.MyIP(), 60)

				if err != nil {
					// return the setA error
					return err
				}

				// TODO also register the ipv6 here...
				dns_client.Close()
			}
		}
	}

	// Here If the DNS provides has api to update the ip address I will use it.
	// TODO test it for different internet provider's
	if len(self.DnsSetA) > 0 {
		// set the data to the actual ip address.
		data := `[{"data":"` + Utility.MyIP() + `"}]`

		// initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, self.DnsSetA, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return err
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "sso-key "+self.DnsKey+":"+self.DnsSecrect)

		// execute the request.
		_, err = client.Do(req)
		if err != nil {
			return (err)
		}

		fmt.Println("ip address for domain", self.getDomain(), "was set to", Utility.MyIP())
	}

	return nil
}

/**
 * Start the grpc proxy.
 */
func (self *Globule) startProxy(id string, port int, proxy int) error {
	srv := self.Services[id]
	if srv.(map[string]interface{})["ProxyProcess"] != nil {
		srv.(map[string]interface{})["ProxyProcess"].(*exec.Cmd).Process.Kill()
		// time.Sleep(time.Second * 1)
	}

	// Now I will start the proxy that will be use by javascript client.
	proxyPath := string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
	if string(os.PathSeparator) == "\\" {
		proxyPath += ".exe" // in case of windows.
	}

	proxyBackendAddress := self.getDomain() + ":" + strconv.Itoa(port)
	proxyAllowAllOrgins := "true"
	proxyArgs := make([]string, 0)

	// Use in a local network or in test.
	proxyArgs = append(proxyArgs, "--backend_addr="+proxyBackendAddress)
	proxyArgs = append(proxyArgs, "--allow_all_origins="+proxyAllowAllOrgins)
	hasTls := Utility.ToBool(srv.(map[string]interface{})["TLS"])
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

		proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+self.certs+string(os.PathSeparator)+self.CertificateAuthorityBundle)
		proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.certs+string(os.PathSeparator)+self.Certificate)

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

	err := proxyProcess.Start()

	if err != nil {
		return err
	}

	// save service configuration.
	srv.(map[string]interface{})["ProxyProcess"] = proxyProcess
	self.Services[id] = srv

	return nil
}

/**
 * That function will
 */
func (self *Globule) keepServiceAlive(s map[string]interface{}) {
	if s["KeepAlive"] == nil {
		return
	}
	// In case the service must not be kept alive.
	keepAlive := Utility.ToBool(s["KeepAlive"])
	if !keepAlive {
		return
	}

	s["Process"].(*exec.Cmd).Wait()

	time.Sleep(time.Second * 5)
	_, _, err := self.startService(s)
	if err != nil {
		return
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
	srv := self.Services[s["Id"].(string)]
	if srv != nil {
		if srv.(map[string]interface{})["Process"] != nil {
			if reflect.TypeOf(srv.(map[string]interface{})["Process"]).String() == "*exec.Cmd" {
				if srv.(map[string]interface{})["Process"].(*exec.Cmd).Process != nil {
					srv.(map[string]interface{})["Process"].(*exec.Cmd).Process.Kill()
				}
			}
		}
	}

	servicePath := self.path + string(os.PathSeparator) + s["Path"].(string)
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

		// Start the service process.
		if string(os.PathSeparator) == "\\" {
			servicePath += ".exe" // in case of windows.
		}

		if s["Name"].(string) == "file.FileService" {
			// File service need root...
			s["Root"] = globule.webRoot
			s["Process"] = exec.Command(servicePath, Utility.ToString(s["Port"]), globule.webRoot)
		} else {
			s["Process"] = exec.Command(servicePath, Utility.ToString(s["Port"]))
		}

		var errb bytes.Buffer
		pipe, _ := s["Process"].(*exec.Cmd).StdoutPipe()
		s["Process"].(*exec.Cmd).Stderr = &errb

		// Here I will set the command dir.
		s["Process"].(*exec.Cmd).Dir = servicePath[:strings.LastIndex(servicePath, string(os.PathSeparator))]

		err = s["Process"].(*exec.Cmd).Start()
		if err != nil {
			s["State"] = "fail"
			log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
			return -1, -1, err
		}

		go func() {

			// Here I will append the service to the load balancer.
			if s["Port"] != nil {
				log.Println("Append ", s["Name"].(string), " to load balancer.")
				load_info := &lbpb.LoadInfo{
					ServerInfo: &lbpb.ServerInfo{
						Id:     s["Id"].(string),
						Name:   s["Name"].(string),
						Domain: s["Domain"].(string),
						Port:   int32(s["Port"].(float64)),
					},
					Load1:  0, // All service will be initialise with a 0 load.
					Load5:  0,
					Load15: 0,
				}

				self.lb_load_info_channel <- load_info
			}
			s["State"] = "running"
			log.Println("Service " + s["Name"].(string) + ":" + s["Id"].(string) + " is up and running!")

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

			// Print the
			if len(errb.String()) > 0 {
				fmt.Println("service", s["Name"].(string), "err:", errb.String())
			}

			// I will remove the service from the load balancer.
			self.lb_remove_candidate_info_channel <- &lbpb.ServerInfo{
				Id:     s["Id"].(string),
				Name:   s["Name"].(string),
				Domain: s["Domain"].(string),
				Port:   int32(s["Port"].(float64)),
			}

		}()

		// Save configuration stuff.
		self.Services[s["Id"].(string)] = s

		// Start the proxy.
		err = self.startProxy(s["Id"].(string), int(s["Port"].(float64)), int(s["Proxy"].(float64)))
		if err != nil {
			return -1, -1, err
		}

		// get back the service info with the proxy process in it
		s = self.Services[s["Id"].(string)].(map[string]interface{})

		// save it to the config.
		self.saveConfig()

	} else if s["Protocol"].(string) == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(s["Name"].(string), "Globular") {

			// Kill previous instance of the program...
			if s["Process"] != nil {
				if s["Process"].(*exec.Cmd).Process != nil {
					s["Process"].(*exec.Cmd).Process.Kill()
					// time.Sleep(time.Second * 1)
				}
			}

			s["Process"] = exec.Command(servicePath, Utility.ToString(s["Port"]))

			var errb bytes.Buffer
			pipe, _ := s["Process"].(*exec.Cmd).StdoutPipe()
			s["Process"].(*exec.Cmd).Stderr = &errb

			// Here I will set the command dir.
			s["Process"].(*exec.Cmd).Dir = servicePath[:strings.LastIndex(servicePath, string(os.PathSeparator))]
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
					// Print the
					fmt.Println("service", s["Name"].(string), "err:", errb.String())

				}()
			}

			// Save configuration stuff.
			self.Services[s["Id"].(string)] = s
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

		// Kill previous instance of the program.
		if hasChange || s["Process"] == nil {
			self.Services[s["Id"].(string)] = s
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
		security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.creds)
	}

	// That will contain all method path from the proto files.
	self.methods = make([]string, 0)
	self.methods = append(self.methods, "/file.FileService/FileUploadHandler")

	self.actionPermissions = make([]interface{}, 0)

	// Set local action permission
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/DeletePermissions", "permission": 1})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/SetRessourceOwner", "permission": 2})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/ressource.RessourceService/DeleteRessourceOwner", "permission": 2})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/admin.AdminService/DeployApplication", "permission": 2})
	self.actionPermissions = append(self.actionPermissions, map[string]interface{}{"action": "/admin.AdminService/PublishService", "permission": 2})

	// It will be execute the first time only...
	configPath := self.config + string(os.PathSeparator) + "config.json"
	if !Utility.Exists(configPath) {
		// Each service contain a file name config.json that describe service.
		// I will keep services info in services map and also it running process.
		basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
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
								self.Services[s["Id"].(string)] = s
								// here I will keep the configuration path in the global configuration.
								s["configPath"] = strings.ReplaceAll(path, self.path, "")[1:]
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
	for _, s := range self.Services {
		name := s.(map[string]interface{})["Name"].(string)
		log.Println("Kill service ", name)
		err := Utility.KillProcessByName(name)
		if err != nil {
			log.Println(err)
		}
	}

	// Start the load balancer.
	err := self.startLoadBalancingService()
	if err != nil {
		log.Println(err)
	}

	// Init services.
	for _, s := range self.Services {
		// Remove existing process information.
		delete(s.(map[string]interface{}), "Process")
		delete(s.(map[string]interface{}), "ProxyProcess")
		err := self.initService(s.(map[string]interface{}))
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

	filepath.Walk(root+path[0:strings.Index(path, "/")],
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path = strings.Replace(path, "\\", "/", -1) // Windows back slash replacement here...
			if strings.HasSuffix(path, importPath_) {
				importPath_ = path
				return io.EOF
			}

			return nil
		})

	importPath_ = strings.Replace(importPath_, strings.Replace(root, "\\", "/", -1), "", -1)

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
	importPath_ = strings.Replace(importPath_, root, "", 1)

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
	m, err := monitoring_client.NewMonitoring_Client(s["Domain"].(string), "monitoring.MonitoringService")
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
	if err != nil {
		log.Println("fail to start monitoring with prometheus", err)
		return err
	}

	alertmanager := exec.Command("alertmanager", "--config.file", self.config+string(os.PathSeparator)+"alertmanager.yml")
	err = alertmanager.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	node_exporter := exec.Command("node_exporter")
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
	self.persistence_client_, err = persistence_client.NewPersistence_Client(s["Domain"].(string), s["Id"].(string))
	if err != nil {

		return nil, err
	}

	// Connect to the database here.
	err = self.persistence_client_.CreateConnection("local_ressource", "local_ressource", "0.0.0.0", 27017, 0, "sa", self.RootPassword, 5000, "", false)
	if err != nil {
		log.Println("---> fail to create local_ressource connection! ", err)
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

	if self.getConfig()["Services"].(map[string]interface{})["ldap.LdapService"] == nil {
		return nil, errors.New("No ldap service configuration was found on that server!")
	}

	var err error

	s := self.getConfig()["Services"].(map[string]interface{})["ldap.LdapService"].(map[string]interface{})

	if self.ldap_client_ == nil {
		self.ldap_client_, err = ldap_client.NewLdap_Client(s["Domain"].(string), "ldap.LdapService")
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
		self.event_client_, err = event_client.NewEvent_Client(s["Domain"].(string), s["Id"].(string))
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
	subscribers := self.keepServicesUpToDate()

	// Catch the Ctrl-C and SIGTERM from kill command
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)

		// Stop load balancer
		self.lb_stop_channel <- true

		// Here the server stop running,
		// so I will close the services.
		for _, value := range self.Services {
			if value.(map[string]interface{})["Process"] != nil {
				p := value.(map[string]interface{})["Process"]
				if reflect.TypeOf(p).String() == "*exec.Cmd" {
					if p.(*exec.Cmd).Process != nil {
						p.(*exec.Cmd).Process.Kill()

					}
				}
			}

			if value.(map[string]interface{})["ProxyProcess"] != nil {
				p := value.(map[string]interface{})["ProxyProcess"]
				if reflect.TypeOf(p).String() == "*exec.Cmd" {
					if p.(*exec.Cmd).Process != nil {
						p.(*exec.Cmd).Process.Kill()
					}
				}
			}
		}

		// Here I will disconnect service update event.
		for id, subscriber := range subscribers {
			eventHub := self.discorveriesEventHub[id]
			for channelId, uuids := range subscriber {
				for i := 0; i < len(uuids); i++ {
					eventHub.UnSubscribe(channelId, uuids[i])
				}
			}
		}

		// stop external service.
		for externalServiceId, _ := range self.ExternalApplications {
			self.stopExternalApplication(externalServiceId)
		}

		// exit cleanly
		os.Exit(0)

	}()

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
		server := &http.Server{
			Addr: ":" + strconv.Itoa(self.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: self.getDomain(),
			},
		}

		// Here I will generate the certificate if it not already exist.
		// TODO generate self signed certificate for localhost...
		if len(self.Certificate) == 0 {
			log.Println(" Now let's encrypts!")
			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/grpc_tls/server.csr run
			// I need to remove the gRPC certificate and recreate it.
			Utility.RemoveDirContents(self.creds)

			security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.certs)
			time.Sleep(15 * time.Second)

			self.initServices() // must restart the services with new certificates.
			err := self.obtainCertificateForCsr()
			if err != nil {
				log.Println("----------------> 1463 ", err)
				return err
			}
		}

		log.Println("start https server")
		log.Println("Globular is running!")
		// get the value from the configuration files.
		err = server.ListenAndServeTLS(self.certs+string(os.PathSeparator)+self.Certificate, self.creds+string(os.PathSeparator)+"server.pem")

	} else {
		log.Println("start http server")
		log.Println("Globular is running!")
		// local - non secure connection.
		err = http.ListenAndServe(":"+strconv.Itoa(self.PortHttp), nil)
	}

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

	resource, err := client.Certificate.ObtainForCSR(*csr, true)
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
	ioutil.WriteFile(self.certs+string(os.PathSeparator)+self.Certificate, resource.Certificate, 0400)
	ioutil.WriteFile(self.certs+string(os.PathSeparator)+self.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	self.saveConfig()

	return nil
}
