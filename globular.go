package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
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
	"os/signal"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecourtois/Globular/event/eventpb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"

	// Admin service
	"github.com/davecourtois/Globular/admin"
	// Ressource service
	"github.com/davecourtois/Globular/ressource"
	// Services management service
	"github.com/davecourtois/Globular/services"
	// Certificate authority
	"github.com/davecourtois/Globular/ca"

	// Interceptor for authentication, event, log...
	"github.com/davecourtois/Globular/Interceptors"

	// Client services.
	"context"
	"crypto"

	"github.com/davecourtois/Globular/dns/dns_client"
	"github.com/davecourtois/Globular/event/event_client"
	"github.com/davecourtois/Globular/ldap/ldap_client"
	"github.com/davecourtois/Globular/monitoring/monitoring_client"
	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/storage/storage_store"
	"github.com/davecourtois/Utility"
	"github.com/emicklei/proto"
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	//"google.golang.org/grpc/grpclog"
	"github.com/davecourtois/Globular/security"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

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
	Name                       string // The service name
	PortHttp                   int    // The port of the http file server.
	PortHttps                  int    // The secure port
	AdminPort                  int    // The admin port
	AdminProxy                 int    // The admin proxy port.
	AdminEmail                 string // The admin email
	RessourcePort              int    // The ressource management service port
	RessourceProxy             int    // The ressource management proxy port
	CertificateAuthorityPort   int    // The certificate authority port
	CertificateAuthorityProxy  int    // The certificate authority proxy port
	ServicesDiscoveryPort      int    // The services discovery port
	ServicesDiscoveryProxy     int    // The ressource management proxy port
	ServicesRepositoryPort     int    // The services discovery port
	ServicesRepositoryProxy    int    // The ressource management proxy port
	Protocol                   string // The protocol of the service.
	Services                   map[string]interface{}
	LdapSyncInfos              map[string]interface{} // Contain LdapSyncInfos...
	ExternalApplications       map[string]ExternalApplication
	Domain                     string   // The domain (subdomain) name of your application
	DNS                        []string // Domain name server use to located the server.
	SessionTimeout             time.Duration
	CertExpirationDelay        int
	CertPassword               string
	Certificate                string
	CertificateAuthorityBundle string
	CertURL                    string
	CertStableURL              string
	Version                    string
	registration               *registration.Resource
	Discoveries                []string // Contain the list of discovery service use to keep service up to date.
	discorveriesEventHub       map[string]*event_client.Event_Client

	// The list of method supported by this server.
	methods []string

	// The prometheus logging informations.
	methodsCounterLog *prometheus.CounterVec

	// prometheus.CounterVec

	// Directories.
	path    string // The path of the exec...
	webRoot string // The root of the http file server.
	data    string // the data directory
	creds   string // gRpc certificate
	certs   string // https certificates
	config  string // configuration directory

	// Log store.
	logs *storage_store.LevelDB_store

	// Create the JWT key used to create the signature
	jwtKey       []byte
	RootPassword string

	// client reference...
	persistence_client_ *persistence_client.Persistence_Client

	ldap_client_  *ldap_client.LDAP_Client
	event_client_ *event_client.Event_Client
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {
	// Here I will initialyse configuration.
	g := new(Globule)
	g.Version = "1.0.0" // Automate version...
	g.RootPassword = "adminadmin"

	g.PortHttp = 8080  // The default http port
	g.PortHttps = 8181 // The default https port number.

	g.Name = strings.Replace(Utility.GetExecName(os.Args[0]), ".exe", "", -1)

	g.Protocol = "http"
	g.Domain = "localhost"
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
		http.ListenAndServe(":10000", nil)
	}()

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Initialyse globular from it configuration file.
	g.config = dir + string(os.PathSeparator) + "config"
	file, err := ioutil.ReadFile(g.config + string(os.PathSeparator) + "config.json")
	// Init the service with the default port address
	if err == nil {
		json.Unmarshal(file, &g)
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
 * Initialize the server directories.
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
    <p>Welcome to Globular 1.0</p>
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
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	config := globule.getConfig()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	crt, err := ioutil.ReadFile(globule.creds + string(os.PathSeparator) + "ca.crt")
	if err != nil {
		http.Error(w, "Client ca cert not found!", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(crt))
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	// sign the certificate.
	csr_str := r.URL.Query().Get("csr") // the csr in base64
	csr, err := base64.StdEncoding.DecodeString(csr_str)

	if err != nil {
		http.Error(w, "Fail to decode csr base64 string", http.StatusBadRequest)
		return
	}

	// Now I will sign the certificate.
	crt, err := globule.signCertificate(string(csr))

	if err != nil {
		http.Error(w, "fail to sign certificate!", http.StatusBadRequest)
		return
	}

	// Return the result as text string.
	fmt.Fprint(w, crt)
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
	token, _ := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa")
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
				token, _ := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa")
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
			log.Println("fail to start external service: ", externalServiceId)
		} else {
			log.Println("external service", externalServiceId, "is started with process id ", pid)
		}
	}

	// Here I will set an environement varibale and I named GLOBULAR_ROOT
	// that will give the root of globular.
	log.Println("Set GLOBULAR_ROOT with value' " + self.path + "'")

	// I will save the variable in a tmp file to be sure I can get it outside
	ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"GLOBULAR_ROOT", []byte(self.path), 0644)

	// set the services.
	self.initServices()

	// Here I will save the server attribute
	self.saveConfig()

	// Here i will connect the service listener.
	time.Sleep(5 * time.Second) // wait for services to start...

	// lisen
	self.Listen()
}

func (self *Globule) createApplicationConnection() error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	values, _ := p.Find("local_ressource", "local_ressource", "Applications", "{}", "")

	applications := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(values), &applications)
	if err != nil {
		return err
	}

	for i := 0; i < len(applications); i++ {
		// Open the user database connection.
		err = p.CreateConnection(applications[i]["_id"].(string)+"_db", applications[i]["_id"].(string)+"_db", "0.0.0.0", 27017, 0, applications[i]["_id"].(string), applications[i]["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Return the domain of the Globule. The name can be empty. If the name is empty
 * it will mean that the domain is entirely control by the globule so it must be
 * able to do it own validation, other wise the domain validation will be done by
 * the globule asscosiate with that domain.
 */
func (self *Globule) getDomain() string {
	domain := self.Domain
	if len(self.Name) > 0 {
		domain = self.Name + "." + domain

	}
	domain = strings.ToLower(domain)
	return domain
}

/**
 * Set the ip for a given sub-domain compose of Name + DNS domain.
 */
func (self *Globule) registerIpToDns() error {
	if self.DNS != nil {
		if len(self.DNS) > 0 {
			for i := 0; i < len(self.DNS); i++ {
				log.Println("register domain to dns:", self.DNS[i])

				client, err := dns_client.NewDns_Client(self.DNS[i], "dns_server")
				if err != nil {
					log.Println("fail to connecto to dns server ", err)
					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = client.SetA(self.Domain, self.getDomain(), Utility.MyIP(), 60)

				if err != nil {
					log.Println("fail to bind address ", Utility.MyIP(), "with domain", self.getDomain(), err)
					return err
				}

				// TODO also register the ipv6 here...
				client.Close()
			}
		}
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
		log.Println("Fail to start grpcwebproxy at port ", proxy, " with error ", err)
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
 * Start a local service.
 */
func (self *Globule) startService(s map[string]interface{}) (int, int, error) {
	var err error

	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")

	if !Utility.IsLocal(s["Domain"].(string)) && string(root) != self.path {
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

	servicePath := self.path + s["servicePath"].(string)
	servicePath = strings.ReplaceAll(strings.ReplaceAll(servicePath, "\\", "/"), "/", string(os.PathSeparator))
	if s["Protocol"].(string) == "grpc" {

		hasTls := Utility.ToBool(s["TLS"])
		if hasTls {
			log.Println("start secure service: ", srv.(map[string]interface{})["Name"])
			// Set TLS local services configuration here.
			s["CertAuthorityTrust"] = self.creds + string(os.PathSeparator) + "ca.crt"
			s["CertFile"] = self.creds + string(os.PathSeparator) + "server.crt"
			s["KeyFile"] = self.creds + string(os.PathSeparator) + "server.pem"
		} else {
			log.Println("start service: ", srv.(map[string]interface{})["Name"])
			// not secure services.
			s["CertAuthorityTrust"] = ""
			s["CertFile"] = ""
			s["KeyFile"] = ""
		}

		// Kill previous instance of the program...
		Utility.KillProcessByName(s["Name"].(string))

		// Start the service process.
		log.Println("try to start process ", s["Name"].(string))

		if string(os.PathSeparator) == "\\" {
			servicePath += ".exe" // in case of windows.
		}

		if s["Name"].(string) == "file_server" {
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

		s["State"] = "running"
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
		log.Println(err)
		return -1, -1, err
	}

	// Return the pid of the service.
	if s["ProxyProcess"] != nil {
		return s["Process"].(*exec.Cmd).Process.Pid, s["ProxyProcess"].(*exec.Cmd).Process.Pid, nil
	}

	return s["Process"].(*exec.Cmd).Process.Pid, -1, nil
}

/**
 * Call once when the server start.
 */
func (self *Globule) initService(s map[string]interface{}) error {
	if s["Protocol"] == nil {
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
	security.GenerateServicesCertificates("1111", self.CertExpirationDelay, self.getDomain(), self.creds)

	// That will contain all method path from the proto files.
	self.methods = make([]string, 0)

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
							path_ := path[:strings.LastIndex(path, string(os.PathSeparator))]
							if s["Name"] == nil {
								log.Println("---> no 'Name' attribute found in service configuration in file config ", path)
							} else {
								s["Id"] = s["Name"].(string)

								s["servicePath"] = strings.Replace(strings.Replace(path_+string(os.PathSeparator)+s["Name"].(string), self.path, "", -1), "\\", "/", -1)
								s["configPath"] = strings.Replace(strings.Replace(path, self.path, "", -1), "\\", "/", -1)
								self.Services[s["Name"].(string)] = s
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
			s := self.Services[name]
			if s == nil {
				s = self.Services[name+"_server"]
			}

			if s != nil {
				s.(map[string]interface{})["protoPath"] = strings.Replace(strings.Replace(path, self.path, "", -1), "\\", "/", -1)
			}

			// here I will parse the service defintion file to extract the
			// service difinition.
			reader, _ := os.Open(path)
			//log.Println("--> proto file: ", name)
			defer reader.Close()

			parser := proto.NewParser(reader)
			definition, _ := parser.Parse()

			// Stack values from walking tree
			stack := make([]interface{}, 0)

			handlePackage := func(stack *[]interface{}) func(*proto.Package) {
				return func(p *proto.Package) {
					*stack = append(*stack, p)
				}
			}(&stack)

			handleService := func(stack *[]interface{}) func(*proto.Service) {
				return func(s *proto.Service) {
					*stack = append(*stack, s)
				}
			}(&stack)

			handleRpc := func(stack *[]interface{}) func(*proto.RPC) {
				return func(r *proto.RPC) {
					*stack = append(*stack, r)
				}
			}(&stack)

			// Walk this way
			proto.Walk(definition,
				proto.WithPackage(handlePackage),
				proto.WithService(handleService),
				proto.WithRPC(handleRpc))

			var packageName string
			var serviceName string
			var methodName string

			for len(stack) > 0 {
				var x interface{}
				x, stack = stack[0], stack[1:]
				switch v := x.(type) {
				case *proto.Package:
					packageName = v.Name
				case *proto.Service:
					serviceName = v.Name
				case *proto.RPC:
					methodName = v.Name
					path := "/" + packageName + "." + serviceName + "/" + methodName
					//log.Println("---> method: ", path)
					// So here I will register the method into the backend.
					self.methods = append(self.methods, path)
				}
			}
		}
		return nil
	})

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

// Method must be register in order to be assign to role.
func (self *Globule) registerMethods() error {
	// Here I will create the sa role if it dosen't exist.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Here I will persit the sa role if it dosent already exist.
	count, err := p.Count("local_ressource", "local_ressource", "Roles", `{ "_id":"sa"}`, "")
	admin := make(map[string]interface{})
	if err != nil {
		return err
	} else if count == 0 {
		log.Println("need to create admin roles...")
		admin["_id"] = "sa"
		admin["name"] = "sa"
		admin["actions"] = self.methods
		jsonStr, _ := Utility.ToJson(admin)
		id, err := p.InsertOne("local_ressource", "local_ressource", "Roles", jsonStr, "")
		if err != nil {
			return err
		}
		log.Println("role with id", id, "was created!")
	} else {
		admin["_id"] = "sa"
		admin["name"] = "sa"
		admin["actions"] = self.methods
		jsonStr, _ := Utility.ToJson(admin)
		// I will set the role actions...
		err = p.ReplaceOne("local_ressource", "local_ressource", "Roles", `{"_id":"sa"}`, jsonStr, "")
		if err != nil {
			return err
		}
		log.Println("role sa with was updated!")
	}

	// I will also create the guest role, the basic one
	count, err = p.Count("local_ressource", "local_ressource", "Roles", `{ "_id":"guest"}`, "")
	guest := make(map[string]interface{})
	if err != nil {
		return err
	} else if count == 0 {
		log.Println("need to create roles guest...")
		guest["_id"] = "guest"
		guest["name"] = "guest"
		guest["actions"] = []string{
			"/admin.AdminService/GetConfig",
			"/ressource.RessourceService/RegisterAccount",
			"/ressource.RessourceService/Authenticate",
			"/ressource.RessourceService/RefreshToken",
			"/ressource.RessourceService/GetPermissions",
			"/ressource.RessourceService/GetAllFilesInfo",
			"/ressource.RessourceService/GetAllApplicationsInfo",
			"/ressource.RessourceService/GetRessourceOwners",
			"/ressource.RessourceService/ValidateUserRessourceAccess",
			"/ressource.RessourceService/ValidateApplicationRessourceAccess",
			"/ressource.RessourceService/ValidateUserRessourceAccess",
			"/ressource.RessourceService/ValidateApplicationAccess",
			"/event.EventService/Subscribe",
			"/event.EventService/UnSubscribe", "/event.EventService/OnEvent",
			"/event.EventService/Quit",
			"/event.EventService/Publish",
			"/services.ServiceDiscovery/FindServices",
			"/services.ServiceDiscovery/GetServiceDescriptor",
			"/services.ServiceDiscovery/GetServicesDescriptor",
			"/services.ServiceRepository/downloadBundle",
			"/persistence.PersistenceService/Find",
			"/persistence.PersistenceService/FindOne",
			"/persistence.PersistenceService/Count",
			"/ressource.RessourceService/GetAllActions"}
		jsonStr, _ := Utility.ToJson(guest)
		_, err := p.InsertOne("local_ressource", "local_ressource", "Roles", jsonStr, "")
		if err != nil {
			return err
		}
		log.Println("role guest was created!")
	}

	// Create connection application.
	self.createApplicationConnection()

	return nil
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
 * This code is use to upload a file into the tmp directory of the server
 * via http request.
 */
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	// I will
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		log.Println(w, err)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames
	var path string                         // grab the filenames

	// Get the path where to upload the file.
	path = r.FormValue("path")
	if strings.HasPrefix(path, "/") {
		path = globule.webRoot + path
		// create the dir if not already exist.
		Utility.CreateDirIfNotExist(path)
	}

	// If application is defined.
	application := r.Header.Get("application")
	token := r.Header.Get("token")
	domain := r.Header.Get("domain")
	hasPermission := false

	if len(application) != 0 {
		err := Interceptors.ValidateApplicationRessourceAccess(domain, application, "/file.FileService/FileUploadHandler", path, 2)
		if err != nil && len(token) == 0 {
			log.Println("Fail to upload the file with error ", err.Error())
			return
		}
		hasPermission = err == nil
	}

	if len(token) != 0 && !hasPermission {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		err := Interceptors.ValidateUserRessourceAccess(domain, token, "/file.FileService/FileUploadHandler", path, 2)
		if err != nil {
			log.Println("Fail to upload the file with error ", err.Error())
			return
		}
	}

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(w, err)
			return
		}

		// Create the file.
		out, err := os.Create(path + "/" + files[i].Filename)

		defer out.Close()
		if err != nil {
			log.Println(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}
		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			log.Println(w, err)
			return
		}
	}
}

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	//if empty, set current directory
	dir := string(root)
	if dir == "" {
		dir = "."
	}

	//add prefix and clean
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	upath = path.Clean(upath)

	//path to file
	name := path.Join(dir, upath)

	// this is the ca certificate use to sign client certificate.
	if upath == "/ca.crt" {
		name = globule.creds + upath
	}

	// Now I will test if a token is given in the header and manage it file access.
	application := r.Header.Get("application")
	token := r.Header.Get("token")
	domain := r.Header.Get("domain")
	hasPermission := false

	if len(application) != 0 {
		err := Interceptors.ValidateApplicationRessourceAccess(domain, application, "/file.FileService/ServeFileHandler", name, 4)
		if err != nil && len(token) == 0 {
			log.Println("Fail to download the file with error ", err.Error())
			return
		}
		hasPermission = err == nil
	}

	if len(token) != 0 && !hasPermission {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		err := Interceptors.ValidateUserRessourceAccess(domain, token, "/file.FileService/ServeFileHandler", name, 4)
		if err != nil {
			log.Println("----> 1108")
			log.Println("Fail to dowload the file with error ", err.Error())
			return
		}
	}

	//check if file exists
	f, err := os.Open(name)

	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File "+upath+" not found!", http.StatusBadRequest)
			return
		}
	}
	defer f.Close()

	// If the file is a javascript file...
	//log.Println("Serve file name: ", name)
	var code string
	hasChange := false
	if strings.HasSuffix(name, ".js") {
		w.Header().Add("Content-Type", "application/javascript")
		if err == nil {
			//hasChange = true
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "import") {
					if strings.Index(line, `'@`) > -1 {
						path_, err := resolveImportPath(upath, line)
						if err == nil {
							line = line[0:strings.Index(line, `'@`)] + `'` + path_ + `'`
							hasChange = true
						}
					}
				}
				code += line + "\n"
			}
		}

	} else if strings.HasSuffix(name, ".css") {
		w.Header().Add("Content-Type", "text/css")
	} else if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".htm") {
		w.Header().Add("Content-Type", "text/html")
	}

	// if the file has change...
	if !hasChange {
		http.ServeFile(w, r, name)
	} else {
		// log.Println(code)
		http.ServeContent(w, r, name, time.Now(), strings.NewReader(code))
	}
}

func (self *Globule) saveConfig() {
	// Here I will save the server attribute
	config, err := Utility.ToMap(self)
	if err == nil {
		services := config["Services"].(map[string]interface{})
		for _, service := range services {
			// remove running information...
			delete(service.(map[string]interface{}), "Process")
			delete(service.(map[string]interface{}), "ProxyProcess")
		}
		str, err := Utility.ToJson(config)
		if err == nil {
			ioutil.WriteFile(self.config+string(os.PathSeparator)+"config.json", []byte(str), 0644)
		}
	}
}

/**
 * Return globular configuration.
 */
func (self *Globule) GetFullConfig(ctx context.Context, rqst *admin.GetConfigRequest) (*admin.GetConfigResponse, error) {

	config, err := Utility.ToMap(self)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	services := config["Services"].(map[string]interface{})
	for _, service := range services {
		// remove running information...
		delete(service.(map[string]interface{}), "Process")
		delete(service.(map[string]interface{}), "ProxyProcess")
	}

	str, err := Utility.ToJson(config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.GetConfigResponse{
		Result: str,
	}, nil

}

func (self *Globule) getConfig() map[string]interface{} {

	config := make(map[string]interface{}, 0)
	config["Name"] = self.Name
	config["PortHttp"] = self.PortHttp
	config["PortHttps"] = self.PortHttps
	config["AdminPort"] = self.AdminPort
	config["AdminProxy"] = self.AdminProxy
	config["AdminEmail"] = self.AdminEmail
	config["RessourcePort"] = self.RessourcePort
	config["RessourceProxy"] = self.RessourceProxy
	config["ServicesDiscoveryPort"] = self.ServicesDiscoveryPort
	config["ServicesDiscoveryProxy"] = self.ServicesDiscoveryProxy
	config["ServicesRepositoryPort"] = self.ServicesRepositoryPort
	config["SessionTimeout"] = self.SessionTimeout
	config["CertificateAuthorityProxy"] = self.CertificateAuthorityProxy
	config["Discoveries"] = self.Discoveries
	config["DNS"] = self.DNS
	config["Protocol"] = self.Protocol
	config["Domain"] = self.getDomain()
	config["CertExpirationDelay"] = self.CertExpirationDelay
	config["ExternalApplications"] = self.ExternalApplications
	config["CertURL"] = self.CertURL
	config["CertStableURL"] = self.CertStableURL
	config["CertificateAuthorityPort"] = self.CertificateAuthorityPort
	config["CertificateAuthorityProxy"] = self.CertificateAuthorityProxy

	// return the full service configuration.
	// Here I will give only the basic services informations and keep
	// all other infromation secret.
	config["Services"] = make(map[string]interface{}) //self.Services
	for name, service_config := range self.Services {
		s := make(map[string]interface{})
		s["Domain"] = service_config.(map[string]interface{})["Domain"]
		s["Port"] = service_config.(map[string]interface{})["Port"]
		s["Proxy"] = service_config.(map[string]interface{})["Proxy"]
		s["TLS"] = service_config.(map[string]interface{})["TLS"]
		s["Version"] = service_config.(map[string]interface{})["Version"]
		s["PublisherId"] = service_config.(map[string]interface{})["PublisherId"]
		s["KeepUpToDate"] = service_config.(map[string]interface{})["KeepUpToDate"]
		s["KeepAlive"] = service_config.(map[string]interface{})["KeepAlive"]
		s["State"] = service_config.(map[string]interface{})["State"]
		s["Id"] = name
		s["Name"] = service_config.(map[string]interface{})["Name"]
		config["Services"].(map[string]interface{})[name] = s
	}

	return config

}

// Return the configuration.
func (self *Globule) GetConfig(ctx context.Context, rqst *admin.GetConfigRequest) (*admin.GetConfigResponse, error) {

	config := self.getConfig()

	str, err := Utility.ToJson(config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.GetConfigResponse{
		Result: str,
	}, nil
}

// Deloyed a web application to a globular node.
func (self *Globule) DeployApplication(stream admin.AdminService_DeployApplicationServer) error {

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer

	var name string
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&admin.DeployApplicationResponse{
				Result: true,
			})
			break
		} else if err != nil {
			return err
		} else {
			name = msg.Name
			buffer.Write(msg.Data)
		}
	}

	// Read bytes and extract it in the current directory.
	r := bytes.NewReader(buffer.Bytes())
	Utility.ExtractTarGz(r)

	// Copy the files to it final destination
	path := self.webRoot + string(os.PathSeparator) + name

	// Remove the existing files.
	if Utility.Exists(path) {
		os.RemoveAll(path)
	}

	// Recreate the dir and move file in it.
	Utility.CreateDirIfNotExist(path)
	Utility.CopyDirContent(Utility.GenerateUUID(name), path)

	// remove temporary files.
	os.RemoveAll(Utility.GenerateUUID(name))

	// Now I will create the application database in the persistence store,
	// and the Application entry in the database.
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	count, err := p.Count("local_ressource", "local_ressource", "Applications", `{"_id":"`+name+`"}`, "")
	application := make(map[string]interface{})
	application["_id"] = name
	application["password"] = Utility.GenerateUUID(name)
	application["path"] = "/" + name                 // The path must be the same as the application name.
	application["last_deployed"] = time.Now().Unix() // save it as unix time.

	if err != nil || count == 0 {

		// create the application database.
		createApplicationUserDbScript := fmt.Sprintf(
			"db=db.getSiblingDB('%s_db');db.createCollection('application_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
			name, name, application["password"].(string), name)

		err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, createApplicationUserDbScript)
		if err != nil {
			log.Println(createApplicationUserDbScript)
			return err
		}

		application["creation_date"] = time.Now().Unix() // save it as unix time.
		data, _ := json.Marshal(&application)
		_, err := p.InsertOne("local_ressource", "local_ressource", "Applications", string(data), "")
		if err != nil {
			return err
		}

		err = p.CreateConnection(name+"_db", name+"_db", "0.0.0.0", 27017, 0, name, application["password"].(string), 5000, "", false)
		if err != nil {
			return err
		}

	} else {

		err := p.UpdateOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+name+`"}`, `{ "$set":{ "last_deployed":`+Utility.ToString(time.Now().Unix())+` }}`, "")
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * Start the monitoring service with prometheus.
 */
func (self *Globule) startMonitoring() error {
	// Cast-it to the persistence client.
	m, err := monitoring_client.NewMonitoring_Client(self.getDomain(), "monitoring_server")
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
    scrape_interval: 500ms
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
 * That function return the sa connection with local mongo db server.
 */
func (self *Globule) getPersistenceSaConnection() (*persistence_client.Persistence_Client, error) {
	// That service made user of persistence service.
	if self.persistence_client_ != nil {
		return self.persistence_client_, nil
	}

	var err error

	// Cast-it to the persistence client.
	self.persistence_client_, err = persistence_client.NewPersistence_Client(self.getDomain(), "persistence_server")
	if err != nil {
		log.Println("fail to connect to persistence server ", err)
		return nil, err
	}

	// Connect to the database here.
	err = self.persistence_client_.CreateConnection("local_ressource", "local_ressource", "0.0.0.0", 27017, 0, "sa", self.RootPassword, 5000, "", false)
	if err != nil {
		return nil, err
	}

	return self.persistence_client_, nil
}

//Set the root password
func (self *Globule) SetRootPassword(ctx context.Context, rqst *admin.SetRootPasswordRequest) (*admin.SetRootPasswordResponse, error) {
	// Here I will set the root password.
	if self.RootPassword != rqst.OldPassword {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong password given!")))
	}

	// Now I will update de sa password.
	self.RootPassword = rqst.NewPassword

	// Now update the sa password in mongo db.
	changeRootPasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", "sa", rqst.NewPassword)

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", rqst.OldPassword, changeRootPasswordScript)
	if err != nil {
		log.Println(changeRootPasswordScript)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.saveConfig()

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")

	return &admin.SetRootPasswordResponse{
		Token: string(token),
	}, nil

}

func (self *Globule) setPassword(accountId string, oldPassword string, newPassword string) error {

	// First of all I will get the user information from the database.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+accountId+`"}`, ``)
	if err != nil {
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return err
	}

	if len(oldPassword) == 0 {
		return errors.New("You must give your old password!")
	}

	// Test the old password.
	if oldPassword != account["password"] {
		if Utility.GenerateUUID(oldPassword) != account["password"] {
			return errors.New("Wrong password given!")
		}
	}

	// Now update the sa password in mongo db.
	name := account["name"].(string)
	name = strings.ReplaceAll(strings.ReplaceAll(name, ".", "_"), "@", "_")

	changePasswordScript := fmt.Sprintf(
		"db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');", name, newPassword)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, changePasswordScript)
	if err != nil {
		return err
	}

	// Here I will update the user information.
	account["password"] = Utility.GenerateUUID(newPassword)

	// Here I will save the role.
	jsonStr := "{"
	jsonStr += `"name":"` + account["name"].(string) + `",`
	jsonStr += `"email":"` + account["email"].(string) + `",`
	jsonStr += `"password":"` + account["password"].(string) + `",`
	jsonStr += `"roles":[`
	for j := 0; j < len(account["roles"].([]interface{})); j++ {
		db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
		db = strings.ReplaceAll(db, "@", "_")
		db = strings.ReplaceAll(db, ".", "_")

		jsonStr += `{`
		jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
		jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
		jsonStr += `"$db":"` + db + `"`
		jsonStr += `}`
		if j < len(account["roles"].([]interface{}))-1 {
			jsonStr += `,`
		}
	}
	jsonStr += `]`
	jsonStr += "}"

	err = p.ReplaceOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
	if err != nil {
		return err
	}

	return nil

}

//Set the root password
func (self *Globule) SetPassword(ctx context.Context, rqst *admin.SetPasswordRequest) (*admin.SetPasswordResponse, error) {

	// First of all I will get the user information from the database.
	err := self.setPassword(rqst.AccountId, rqst.OldPassword, rqst.NewPassword)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
	return &admin.SetPasswordResponse{
		Token: string(token),
	}, nil

}

//Set the root password
func (self *Globule) SetEmail(ctx context.Context, rqst *admin.SetEmailRequest) (*admin.SetEmailResponse, error) {

	// Here I will set the root password.
	// First of all I will get the user information from the database.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if account["email"].(string) != rqst.OldEmail {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given!")))
	}

	account["email"] = rqst.NewEmail

	// Here I will save the role.
	jsonStr := "{"
	jsonStr += `"name":"` + account["name"].(string) + `",`
	jsonStr += `"email":"` + account["email"].(string) + `",`
	jsonStr += `"password":"` + account["password"].(string) + `",`
	jsonStr += `"roles":[`
	for j := 0; j < len(account["roles"].([]interface{})); j++ {
		db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
		db = strings.ReplaceAll(db, "@", "_")
		db = strings.ReplaceAll(db, ".", "_")
		jsonStr += `{`
		jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
		jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
		jsonStr += `"$db":"` + db + `"`
		jsonStr += `}`
		if j < len(account["roles"].([]interface{}))-1 {
			jsonStr += `,`
		}
	}
	jsonStr += `]`
	jsonStr += "}"

	// set the new email.
	account["email"] = rqst.NewEmail

	err = p.ReplaceOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+account["name"].(string)+`"}`, jsonStr, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// read the token.
	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the token.
	return &admin.SetEmailResponse{
		Token: string(token),
	}, nil
}

//Set the root email
func (self *Globule) SetRootEmail(ctx context.Context, rqst *admin.SetRootEmailRequest) (*admin.SetRootEmailResponse, error) {
	// Here I will set the root password.

	if self.AdminEmail != rqst.OldEmail {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Wrong email given!")))
	}

	// Now I will update de sa password.
	self.AdminEmail = rqst.NewEmail

	// save the configuration.
	self.saveConfig()

	// read the token.
	token, err := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.getDomain() + "_token")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the token.
	return &admin.SetRootEmailResponse{
		Token: string(token),
	}, nil
}

// Upload a service package.
func (self *Globule) UploadServicePackage(stream admin.AdminService_UploadServicePackageServer) error {
	// The bundle will cantain the necessary information to install the service.
	path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()
	fo, err := os.Create(path)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	defer fo.Close()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&admin.UploadServicePackageResponse{
				Path: path,
			})
			break
		} else if err != nil {
			return err
		} else {
			fo.Write(msg.Data)
		}
	}

	return nil
}

// Publish a service. The service must be install localy on the server.
func (self *Globule) PublishService(ctx context.Context, rqst *admin.PublishServiceRequest) (*admin.PublishServiceResponse, error) {

	// Connect to the dicovery services
	services_discovery, err := services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	// Connect to the repository services.
	services_repository, err := services.NewServicesRepository_Client(rqst.RepositoryId, "services_repository")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.RepositoryId)))
	}

	// Now I will upload the service to the repository...
	serviceDescriptor := &services.ServiceDescriptor{
		Id:          rqst.ServiceId,
		PublisherId: rqst.PublisherId,
		Version:     rqst.Version,
		Description: rqst.Description,
		Keywords:    rqst.Keywords,
	}

	err = services_discovery.PublishServiceDescriptor(serviceDescriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Upload the service to the repository.
	err = services_repository.UploadBundle(rqst.DicorveryId, serviceDescriptor.Id, serviceDescriptor.PublisherId, int32(rqst.Platform), rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// So here I will send an plublish event...
	err = os.Remove(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will send a event to be sure all server will update...
	data, _ := json.Marshal(serviceDescriptor)

	// Here I will send an event that the service has a new version...
	eventHub, err := self.getEventHub()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	eventHub.Publish(serviceDescriptor.PublisherId+":"+serviceDescriptor.Id+":SERVICE_PUBLISH_EVENT", data)

	return &admin.PublishServiceResponse{
		Result: true,
	}, nil
}

// Install/Update a service on globular instance.
func (self *Globule) installService(descriptor *services.ServiceDescriptor) error {

	// repository must exist...
	if len(descriptor.Repositories) == 0 {
		return errors.New("No service repository was found for service " + descriptor.Id)
	}

	var platform services.Platform
	// The first step will be to create the archive.
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			platform = services.Platform_WIN64
		} else if runtime.GOARCH == "386" {
			platform = services.Platform_WIN32
		}
	} else if runtime.GOOS == "linux" { // also can be specified to FreeBSD
		if runtime.GOARCH == "amd64" {
			platform = services.Platform_LINUX64
		} else if runtime.GOARCH == "386" {
			platform = services.Platform_LINUX32
		}
	} else if runtime.GOOS == "darwin" {
		/** TODO Deploy services on other platforme here... **/
	}

	for i := 0; i < len(descriptor.Repositories); i++ {
		services_repository, err := services.NewServicesRepository_Client(descriptor.Repositories[i], "services_repository")
		if err != nil {
			return err
		}
		bundle, err := services_repository.DownloadBundle(descriptor, platform)
		if err == nil {
			id := descriptor.PublisherId + "%" + descriptor.Id + "%" + descriptor.Version
			if platform == services.Platform_LINUX32 {
				id += "%LINUX32"
			} else if platform == services.Platform_LINUX64 {
				id += "%LINUX64"
			} else if platform == services.Platform_WIN32 {
				id += "%WIN32"
			} else if platform == services.Platform_WIN64 {
				id += "%WIN64"
			}

			// Create the file.
			r := bytes.NewReader(bundle.Binairies)
			Utility.ExtractTarGz(r)

			// I will save the binairy in file...
			dest := "globular_services" + string(os.PathSeparator) + strings.ReplaceAll(id, "%", string(os.PathSeparator))
			Utility.CreateDirIfNotExist(dest)
			Utility.CopyDirContent(self.path+string(os.PathSeparator)+id, self.path+string(os.PathSeparator)+dest)

			// remove the file...
			os.RemoveAll(self.path + string(os.PathSeparator) + id)

			// I will repalce the service configuration with the new one...
			jsonStr, err := ioutil.ReadFile(dest + string(os.PathSeparator) + "config.json")
			if err != nil {
				return err
			}

			config := make(map[string]interface{})
			json.Unmarshal(jsonStr, &config)

			// save the new paths...
			config["servicePath"] = strings.ReplaceAll(string(os.PathSeparator)+dest+string(os.PathSeparator)+config["Name"].(string), string(os.PathSeparator), "/")
			config["protoPath"] = strings.ReplaceAll(string(os.PathSeparator)+dest+string(os.PathSeparator)+config["Name"].(string)+".proto", string(os.PathSeparator), "/")
			config["configPath"] = strings.ReplaceAll(string(os.PathSeparator)+dest+string(os.PathSeparator)+"config.json", string(os.PathSeparator), "/")

			// Here I will append the execute permission to the service file.

			// Set execute permission
			err = os.Chmod(self.path+config["servicePath"].(string), 0755)
			if err != nil {
				fmt.Println(err)
			}

			// initialyse the new service.
			self.initService(config)

			break
		}
	}

	return nil

}

// Install/Update a service on globular instance.
func (self *Globule) InstallService(ctx context.Context, rqst *admin.InstallServiceRequest) (*admin.InstallServiceResponse, error) {

	// Connect to the dicovery services
	var services_discovery *services.ServicesDiscovery_Client
	var err error
	services_discovery, err = services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")

	if services_discovery == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	descriptors, err := services_discovery.GetServiceDescriptor(rqst.ServiceId, rqst.PublisherId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The first element in the array is the most recent descriptor
	// so if no version is given the most recent will be taken.
	descriptor := descriptors[0]
	for i := 0; i < len(descriptors); i++ {
		if descriptors[i].Version == rqst.Version {
			descriptor = descriptors[i]
			break
		}
	}

	err = self.installService(descriptor)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.InstallServiceResponse{
		Result: true,
	}, nil

}

// Uninstall a service...
func (self *Globule) UninstallService(ctx context.Context, rqst *admin.UninstallServiceRequest) (*admin.UninstallServiceResponse, error) {

	// First of all I will stop the running service(s) instance.
	for id, service := range self.Services {
		// Stop the instance of the service.
		s := service.(map[string]interface{})
		if s["Name"] != nil {
			if s["PublisherId"].(string) == rqst.PublisherId && s["Name"].(string) == rqst.ServiceId && s["Version"].(string) == rqst.Version {
				self.stopService(s["Id"].(string))
				delete(self.Services, id)
			}
		}
	}

	// Now I will remove the service.
	// Service are located into the globular_services...
	path := self.path + string(os.PathSeparator) + "globular_services" + string(os.PathSeparator) + rqst.PublisherId + string(os.PathSeparator) + rqst.ServiceId + string(os.PathSeparator) + rqst.Version

	// remove directory and sub-directory.
	err := os.RemoveAll(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// save the config.
	self.saveConfig()

	return &admin.UninstallServiceResponse{
		Result: true,
	}, nil
}

// return true if the configuation has change.
func (self *Globule) saveServiceConfig(config map[string]interface{}) bool {
	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	if !Utility.IsLocal(config["Domain"].(string)) && string(root) != self.path {
		return false
	}

	// set the domain of the service.
	config["Domain"] = self.getDomain()

	// get the config path.
	var process interface{}
	var proxyProcess interface{}

	process = config["Process"]
	proxyProcess = config["ProxyProcess"]

	// remove unused information...
	delete(config, "Process")
	delete(config, "ProxyProcess")

	// so here I will get the previous information...
	f, err := os.Open(self.path + string(os.PathSeparator) + config["configPath"].(string))

	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err == nil {
			config_ := make(map[string]interface{})
			json.Unmarshal(b, &config_)
			if reflect.DeepEqual(config_, config) {
				f.Close()
				// set back the path's info.
				config["Process"] = process
				config["ProxyProcess"] = proxyProcess
				return false
			}
		}
	}
	f.Close()

	// sync the data/config file with the service file.
	jsonStr, _ := Utility.ToJson(config)

	// here I will write the file
	err = ioutil.WriteFile(self.path+string(os.PathSeparator)+config["configPath"].(string), []byte(jsonStr), 0644)
	if err != nil {
		log.Println("fail to save config file: ", err)
	}

	// set back internal infos...
	config["Process"] = process
	config["ProxyProcess"] = proxyProcess

	return true
}

// Save a service configuration
func (self *Globule) SaveConfig(ctx context.Context, rqst *admin.SaveConfigRequest) (*admin.SaveConfigResponse, error) {

	config := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(rqst.Config), &config)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// if the configuration is one of services...
	if config["Id"] != nil {
		srv := self.Services[config["Id"].(string)]
		if srv != nil {
			// Attach the actual process and proxy process to the configuration object.
			config["Process"] = srv.(map[string]interface{})["Process"]
			config["ProxyProcess"] = srv.(map[string]interface{})["ProxyProcess"]
			self.initService(config)
		}
	} else if config["Services"] != nil {
		// Here I will save the configuration
		self.Name = config["Name"].(string)
		self.PortHttp = Utility.ToInt(config["PortHttp"].(float64))
		self.PortHttps = Utility.ToInt(config["PortHttps"].(float64))
		self.AdminEmail = config["AdminEmail"].(string)
		self.AdminPort = Utility.ToInt(config["AdminPort"].(float64))
		self.AdminProxy = Utility.ToInt(config["AdminProxy"].(float64))
		self.RessourcePort = Utility.ToInt(config["RessourcePort"].(float64))
		self.RessourceProxy = Utility.ToInt(config["RessourceProxy"].(float64))
		self.ServicesDiscoveryPort = Utility.ToInt(config["ServicesDiscoveryPort"].(float64))
		self.ServicesDiscoveryProxy = Utility.ToInt(config["ServicesDiscoveryProxy"].(float64))
		self.ServicesRepositoryPort = Utility.ToInt(config["ServicesRepositoryPort"].(float64))
		self.ServicesRepositoryProxy = Utility.ToInt(config["ServicesRepositoryProxy"].(float64))
		self.CertificateAuthorityPort = Utility.ToInt(config["CertificateAuthorityPort"].(float64))
		self.CertificateAuthorityProxy = Utility.ToInt(config["CertificateAuthorityProxy"].(float64))

		self.Protocol = config["Protocol"].(string)
		self.Domain = config["Domain"].(string)
		self.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))

		// That will save the services if they have changed.
		for n, s := range config["Services"].(map[string]interface{}) {
			// Attach the actual process and proxy process to the configuration object.
			s.(map[string]interface{})["Process"] = self.Services[n].(map[string]interface{})["Process"]
			s.(map[string]interface{})["ProxyProcess"] = self.Services[n].(map[string]interface{})["ProxyProcess"]
			self.initService(s.(map[string]interface{}))
		}

		// Save Discoveries.
		self.Discoveries = make([]string, 0)
		for i := 0; i < len(config["Discoveries"].([]interface{})); i++ {
			self.Discoveries = append(self.Discoveries, config["Discoveries"].([]interface{})[i].(string))
		}

		// Save DNS
		self.DNS = make([]string, 0)
		for i := 0; i < len(config["DNS"].([]interface{})); i++ {
			self.DNS = append(self.DNS, config["DNS"].([]interface{})[i].(string))
		}

		// save the application server.
		self.saveConfig()
	}

	// return the new configuration file...
	result, _ := Utility.ToJson(config)
	return &admin.SaveConfigResponse{
		Result: result,
	}, nil
}

func (self *Globule) stopService(serviceId string) error {
	if (self.Services[serviceId]) == nil {
		return errors.New("no service with id " + serviceId + " is define on the server!")
	}
	s := self.Services[serviceId].(map[string]interface{})
	if s == nil {
		return errors.New("No service found with id " + serviceId)
	}

	// Set keep alive to false...
	s["KeepAlive"] = false

	if s["Process"] == nil {
		return errors.New("No process running")
	}

	if s["Process"].(*exec.Cmd).Process == nil {
		return errors.New("No process running")
	}

	err := s["Process"].(*exec.Cmd).Process.Kill()
	if err != nil {
		return err
	}

	if s["ProxyProcess"] != nil {
		err := s["ProxyProcess"].(*exec.Cmd).Process.Kill()
		if err != nil {
			return err
		}
	}

	s["State"] = "stopped"

	self.logServiceInfo(s["Name"].(string), time.Now().String()+"Service "+s["Name"].(string)+" was stopped!")

	log.Println("stop service", s["Name"])

	return nil
}

// Stop a service
func (self *Globule) StopService(ctx context.Context, rqst *admin.StopServiceRequest) (*admin.StopServiceResponse, error) {
	err := self.stopService(rqst.ServiceId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.StopServiceResponse{
		Result: true,
	}, nil
}

// Start a service
func (self *Globule) StartService(ctx context.Context, rqst *admin.StartServiceRequest) (*admin.StartServiceResponse, error) {

	s := self.Services[rqst.ServiceId]
	if s == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+rqst.ServiceId)))
	}

	service_pid, proxy_pid, err := self.startService(s.(map[string]interface{}))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &admin.StartServiceResponse{
		ProxyPid:   int64(proxy_pid),
		ServicePid: int64(service_pid),
	}, nil
}

// Start an external service here.
func (self *Globule) startExternalApplication(serviceId string) (int, error) {

	if service, ok := self.ExternalApplications[serviceId]; !ok {
		return -1, errors.New("No external service found with name " + serviceId)
	} else {

		service.srv = exec.Command(service.Path, service.Args...)

		err := service.srv.Start()
		if err != nil {
			return -1, err
		}

		// save back the service in the map.
		self.ExternalApplications[serviceId] = service

		return service.srv.Process.Pid, nil
	}

}

// Stop external service.
func (self *Globule) stopExternalApplication(serviceId string) error {
	if _, ok := self.ExternalApplications[serviceId]; !ok {
		return errors.New("No external service found with name " + serviceId)
	}

	// if no command was created
	if self.ExternalApplications[serviceId].srv == nil {
		return nil
	}

	// if no process running
	if self.ExternalApplications[serviceId].srv.Process == nil {
		return nil
	}

	// kill the process.
	return self.ExternalApplications[serviceId].srv.Process.Kill()
}

// Register external service to be start by Globular in order to run
func (self *Globule) RegisterExternalApplication(ctx context.Context, rqst *admin.RegisterExternalApplicationRequest) (*admin.RegisterExternalApplicationResponse, error) {

	// Here I will get the command path.
	externalCmd := ExternalApplication{
		Id:   rqst.ServiceId,
		Path: rqst.Path,
		Args: rqst.Args,
	}

	self.ExternalApplications[externalCmd.Id] = externalCmd

	// save the config.
	self.saveConfig()

	// start the external service.
	pid, err := self.startExternalApplication(externalCmd.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.RegisterExternalApplicationResponse{
		ServicePid: int64(pid),
	}, nil
}

/**
 * Start internal service admin and ressource are use that function.
 */
func (self *Globule) startInternalService(id string, port int, proxy int, hasTls bool, unaryInterceptor grpc.UnaryServerInterceptor, streamInterceptor grpc.StreamServerInterceptor) (*grpc.Server, error) {

	if self.Services[id] != nil {
		hasTls = self.Services[id].(map[string]interface{})["TLS"].(bool)
	}

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, id+" service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := make(map[string]interface{}, 0)
	var grpcServer *grpc.Server
	if hasTls {
		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
		certFile := self.creds + string(os.PathSeparator) + "server.crt"
		keyFile := self.creds + string(os.PathSeparator) + "server.pem"

		s["CertFile"] = certFile
		s["KeyFile"] = keyFile
		s["CertAuthorityTrust"] = certAuthorityTrust

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("could not load server key pair: %s", err)
			return nil, err
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(certAuthorityTrust)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
			return nil, err
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append client certs")
			return nil, err
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

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
	s["Port"] = port
	s["Proxy"] = proxy
	s["TLS"] = hasTls

	self.Services[id] = s

	// save the config.
	self.saveConfig()

	// start the proxy
	err := self.startProxy(id, port, proxy)
	if err != nil {
		return nil, err
	}

	return grpcServer, nil
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
			log.Println("mongo fail to execute the script.")
			return errors.New("mongod is not responding!")
		}
		// call again.
		timeout -= 1
		return self.waitForMongo(timeout, withAuth)
	}
	return nil
}

/** Create the super administrator in the db. **/
func (self *Globule) registerSa() error {

	// Here I will create super admin if it not already exist.
	dataPath := self.data + string(os.PathSeparator) + "mongodb-data"

	if !Utility.Exists(dataPath) {
		// Kill mongo db server if the process already run...
		self.stopMongod()

		// Here I will create the directory
		err := os.MkdirAll(dataPath, os.ModeDir)
		if err != nil {
			log.Println("fail to create dir", err)
			return err
		}

		// Now I will start the command
		mongod := exec.Command("mongod", "--port", "27017", "--dbpath", dataPath)
		err = mongod.Start()
		if err != nil {
			log.Println("fail to start mongo db", err)
			return err
		}

		self.waitForMongo(60, false)

		// Now I will create a new user name sa and give it all admin write.
		createSaScript := fmt.Sprintf(
			`db=db.getSiblingDB('admin');db.createUser({ user: '%s', pwd: '%s', roles: ['userAdminAnyDatabase','userAdmin','readWrite','dbAdmin','clusterAdmin','readWriteAnyDatabase','dbAdminAnyDatabase']});`, "sa", self.RootPassword) // must be change...

		createSaCmd := exec.Command("mongo", "--eval", createSaScript)
		err = createSaCmd.Run()
		if err != nil {
			// remove the mongodb-data
			os.RemoveAll(dataPath)
			log.Println(createSaScript)
			return err
		}
		self.stopMongod()
	}

	// Now I will start mongod with auth available.
	mongod := exec.Command("mongod", "--auth", "--port", "27017", "--bind_ip", "0.0.0.0", "--dbpath", dataPath)
	err := mongod.Start()
	if err != nil {
		return err
	}

	// wait 15 seconds that the server restart.
	self.waitForMongo(60, true)

	// Get the list of all services method.
	return self.registerMethods()
}

func (self *Globule) getLdapClient() (*ldap_client.LDAP_Client, error) {
	var err error
	if self.ldap_client_ == nil {
		self.ldap_client_, err = ldap_client.NewLdap_Client(self.getDomain(), "ldap_server")
	}
	if err != nil {
		log.Println("fail to connect to ldap server ", err)
	}
	return self.ldap_client_, err
}

/** Append new LDAP synchronization informations. **/
func (self *Globule) SynchronizeLdap(ctx context.Context, rqst *ressource.SynchronizeLdapRqst) (*ressource.SynchronizeLdapRsp, error) {

	if rqst.SyncInfo == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync infos was given!")))
	}

	if rqst.SyncInfo.UserSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync users infos was given!")))
	}

	syncInfo, err := Utility.ToMap(rqst.SyncInfo)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] == nil {
		self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = make([]interface{}, 0)
		self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = append(self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{}), syncInfo)
	} else {
		syncInfos := self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{})
		exist := false
		for i := 0; i < len(syncInfos); i++ {
			if syncInfos[i].(map[string]interface{})["ldapSeriveId"] == syncInfo["ldapSeriveId"] {
				if syncInfos[i].(map[string]interface{})["connectionId"] == syncInfo["connectionId"] {
					// set the connection info.
					syncInfos[i] = syncInfo
					exist = true
					// save the config.
					self.saveConfig()

					break
				}
			}
		}

		if !exist {
			self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)] = append(self.LdapSyncInfos[syncInfo["ldapSeriveId"].(string)].([]interface{}), syncInfo)
			// save the config.
			self.saveConfig()

		}
	}

	// Cast the the correct type.

	// Searh for roles.
	ldap_, err := self.getLdapClient()
	if err != nil {
		return nil, err
	}
	rolesInfo, err := ldap_.Search(rqst.SyncInfo.ConnectionId, rqst.SyncInfo.GroupSyncInfos.Base, rqst.SyncInfo.GroupSyncInfos.Query, []string{rqst.SyncInfo.GroupSyncInfos.Id, "distinguishedName"})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print role info.
	for i := 0; i < len(rolesInfo); i++ {
		name := rolesInfo[i][0].([]interface{})[0].(string)
		id := Utility.GenerateUUID(rolesInfo[i][1].([]interface{})[0].(string))
		self.createRole(id, name, []string{})
	}

	// Synchronize account and user info...
	accountsInfo, err := ldap_.Search(rqst.SyncInfo.ConnectionId, rqst.SyncInfo.UserSyncInfos.Base, rqst.SyncInfo.UserSyncInfos.Query, []string{rqst.SyncInfo.UserSyncInfos.Id, rqst.SyncInfo.UserSyncInfos.Email, "distinguishedName", "memberOf"})

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(accountsInfo); i++ {
		// Print the list of account...
		// I will not set the password...
		name := strings.ToLower(accountsInfo[i][0].([]interface{})[0].(string))

		if len(accountsInfo[i][1].([]interface{})) > 0 {
			email := strings.ToLower(accountsInfo[i][1].([]interface{})[0].(string))

			if len(email) > 0 {

				id := Utility.GenerateUUID(strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string)))
				if len(id) > 0 {
					log.Println("---> register account: ", accountsInfo[i])

					roles := make([]interface{}, 0)
					roles = append(roles, "guest")
					// Here I will set the roles of the user.
					if len(accountsInfo[i][3].([]interface{})) > 0 {
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							roles = append(roles, Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)))
						}
					}

					// Try to create account...
					err := self.registerAccount(id, name, email, id, roles)
					if err == nil {
						log.Println("register account ", id)
					} else {
						rolesStr := `[{"$ref":"Roles","$id":"guest","$db":"local_ressource"}`
						for j := 0; j < len(accountsInfo[i][3].([]interface{})); j++ {
							rolesStr += `,{"$ref":"Roles","$id":"` + Utility.GenerateUUID(accountsInfo[i][3].([]interface{})[j].(string)) + `","$db":"local_ressource"}`
						}

						rolesStr += `]`
						err := p.UpdateOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+id+`"}`, `{ "$set":{"roles":`+rolesStr+`}}`, "")
						if err == nil {
							log.Println("account ", id, " was update!")
						} else {
							log.Println("fail to update account ", id, err)
						}

					}
				}
			} else {
				log.Println("account " + strings.ToLower(accountsInfo[i][2].([]interface{})[0].(string)) + " has no email configured! ")
			}
		}
	}

	if rqst.SyncInfo.GroupSyncInfos == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No LDAP sync groups infos was given!")))
	}

	return &ressource.SynchronizeLdapRsp{
		Result: true,
	}, nil
}

func (self *Globule) registerAccount(id string, name string, email string, password string, roles []interface{}) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// first of all the Persistence service must be active.
	count, err := p.Count("local_ressource", "local_ressource", "Accounts", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	// one account already exist for the name.
	if count == 1 {
		return errors.New("account with name " + name + " already exist!")
	}

	// set the account object and set it basic roles.
	account := make(map[string]interface{})
	account["_id"] = id
	account["name"] = name
	account["email"] = email
	account["password"] = Utility.GenerateUUID(password) // hide the password...

	account["roles"] = make([]map[string]interface{}, 0)
	for i := 0; i < len(roles); i++ {
		role := make(map[string]interface{}, 0)
		role["$id"] = roles[i]
		role["$ref"] = "Roles"
		role["$db"] = "local_ressource"
		account["roles"] = append(account["roles"].([]map[string]interface{}), role)
	}

	// serialyse the account and save it.
	accountStr, err := json.Marshal(account)
	if err != nil {
		return err
	}

	// Here I will insert the account in the database.
	_, err = p.InsertOne("local_ressource", "local_ressource", "Accounts", string(accountStr), "")

	// replace @ and . by _
	name = strings.ReplaceAll(strings.ReplaceAll(name, "@", "_"), ".", "_")

	// Each account will have their own database and a use that can read and write
	// into it.
	// Here I will wrote the script for mongoDB...
	createUserScript := fmt.Sprintf(
		"db=db.getSiblingDB('%s_db');db.createCollection('user_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'dbOwner', db: '%s_db' }]});",
		name, name, password, name)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, createUserScript)
	if err != nil {
		return err
	}

	err = p.CreateConnection(name+"_db", name+"_db", "0.0.0.0", 27017, 0, name, password, 5000, "", false)
	if err != nil {
		return errors.New("No persistence service are available to store ressource information.")
	}

	return nil

}

/* Register a new Account */
func (self *Globule) RegisterAccount(ctx context.Context, rqst *ressource.RegisterAccountRqst) (*ressource.RegisterAccountRsp, error) {
	if rqst.ConfirmPassword != rqst.Password {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to confirm your password!")))

	}

	if rqst.Account == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account information was given!")))

	}

	err := self.registerAccount(rqst.Account.Name, rqst.Account.Name, rqst.Account.Email, rqst.Password, []interface{}{"guest"})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, rqst.Account.Name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	name, expireAt, _ := Interceptors.ValidateToken(tokenString)
	_, err = p.InsertOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will
	return &ressource.RegisterAccountRsp{
		Result: tokenString, // Return the token string.
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Peer's Authorization and Authentication code.
////////////////////////////////////////////////////////////////////////////////

//* Register a new Peer on the network *
func (self *Globule) RegisterPeer(ctx context.Context, rqst *ressource.RegisterPeerRqst) (*ressource.RegisterPeerRsp, error) {
	// A peer want to be part of the network.

	// Get the persistence connection
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		log.Println("authenticate fail to get persistence connection ", err)
		return nil, err
	}

	// Here I will first look if a peer with a same name already exist on the
	// ressources...
	count, _ := p.Count("local_ressource", "local_ressource", "Peers", `{"Name":"`+rqst.Peer.Name+`"}`, "")
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer with name '"+rqst.Peer.Name+"' already exist!")))

	}

	// No authorization exist for that peer I will insert it.

	// Here I will
	return nil, nil
}

//* Return the list of authorized peers *
func (self *Globule) GetPeers(rqst *ressource.GetPeersRqst, stream ressource.RessourceService_GetPeersServer) error {
	return nil
}

//* Remove a peer from the network *
func (self *Globule) DeletePeer(ctx context.Context, rqst *ressource.DeletePeerRqst) (*ressource.DeletePeerRsp, error) {
	return nil, nil
}

//* Add peer action permission *
func (self *Globule) AddPeerAction(ctx context.Context, rqst *ressource.AddPeerActionRqst) (*ressource.AddPeerActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Peers", `{"_id":"`+rqst.PeerId+`"}`, ``)
	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	peer := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &peer)

	needSave := false
	if peer["actions"] == nil {
		peer["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		for i := 0; i < len(peer["actions"].([]interface{})); i++ {
			if peer["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			peer["actions"] = append(peer["actions"].([]interface{}), rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer named "+rqst.PeerId+" already contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(peer)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Peers", `{"_id":"`+rqst.PeerId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.AddPeerActionRsp{Result: true}, nil

}

//* Remove peer action permission *
func (self *Globule) RemovePeerAction(ctx context.Context, rqst *ressource.RemovePeerActionRqst) (*ressource.RemovePeerActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Peers", `{"_id":"`+rqst.PeerId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	peer := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &peer)

	needSave := false
	if peer["actions"] == nil {
		peer["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		for i := 0; i < len(peer["actions"].([]interface{})); i++ {
			if peer["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, peer["actions"].([]interface{})[i])
			}
		}
		if exist {
			peer["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Peer named "+rqst.PeerId+" not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(peer)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Peers", `{"_id":"`+rqst.PeerId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.RemovePeerActionRsp{Result: true}, nil
}

//* Authenticate a account by it name or email.
// That function test if the password is the correct one for a given user
// if it is a token is generate and that token will be use by other service
// to validate permission over the requested ressource.
func (self *Globule) Authenticate(ctx context.Context, rqst *ressource.AuthenticateRqst) (*ressource.AuthenticateRsp, error) {
	// Get the persistence connection
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		log.Println("authenticate fail to get persistence connection ", err)
		return nil, err
	}

	// in case of sa user.(admin)
	if (rqst.Password == self.RootPassword && rqst.Name == "sa") || (rqst.Password == self.RootPassword && rqst.Name == self.AdminEmail) {
		// Generate a token to identify the user.
		tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		/** Return the token only **/
		return &ressource.AuthenticateRsp{
			Token: tokenString,
		}, nil
	}

	values, err := p.Find("local_ressource", "local_ressource", "Accounts", `{"name":"`+rqst.Name+`"}`, "")
	if err != nil || values == "[]" {
		values, err = p.Find("local_ressource", "local_ressource", "Accounts", `{"email":"`+rqst.Name+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	objects := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(values), &objects)

	if len(objects) == 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("fail to retreive "+rqst.Name+" informations.")))
	}

	ldap_, err := self.getLdapClient()
	if err != nil {
		return nil, err
	}

	if objects[0]["password"].(string) != Utility.GenerateUUID(rqst.Password) {
		// Here I will try to made use of ldap if there is a service configure.ldap
		err := ldap_.Authenticate("", objects[0]["name"].(string), rqst.Password)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Set the password whit
		err = self.setPassword(objects[0]["_id"].(string), objects[0]["password"].(string), rqst.Password)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, objects[0]["name"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	name_ := objects[0]["name"].(string)
	name_ = strings.ReplaceAll(strings.ReplaceAll(name_, ".", "_"), "@", "_")

	// Open the user database connection.
	err = p.CreateConnection(name_+"_db", name_+"_db", "0.0.0.0", 27017, 0, name_, rqst.Password, 5000, "", false)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))
	}

	// save the newly create token into the database.
	name, expireAt, _ := Interceptors.ValidateToken(tokenString)
	err = p.ReplaceOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I got the token I will now put it in the cache.
	return &ressource.AuthenticateRsp{
		Token: tokenString,
	}, nil
}

/**
 * Refresh token get a new token.
 */
func (self *Globule) RefreshToken(ctx context.Context, rqst *ressource.RefreshTokenRqst) (*ressource.RefreshTokenRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// first of all I will validate the current token.
	name, expireAt, _ := Interceptors.ValidateToken(rqst.Token)
	// If the token is older than seven day without being refresh then I retrun an error.
	if time.Unix(expireAt, 0).Before(time.Now().AddDate(0, 0, -7)) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token cannot be refresh after 7 day")))
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, _ := p.FindOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `[{"Projection":{"expireAt":1}}]`)
	if len(values) != 0 {
		lastTokenInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(values), &lastTokenInfo)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		savedTokenExpireAt := time.Unix(int64(lastTokenInfo["expireAt"].(float64)), 0)

		// That mean a newer token was already refresh.
		if savedTokenExpireAt.Before(time.Unix(expireAt, 0)) {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("That token cannot not be refresh because a newer one already exist. You need to re-authenticate in order to get a new token.")))
		}
	}

	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, name)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get back the new expireAt
	name, expireAt, _ = Interceptors.ValidateToken(tokenString)

	err = p.ReplaceOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return the token string.
	return &ressource.RefreshTokenRsp{
		Token: tokenString,
	}, nil
}

//* Delete an account *
func (self *Globule) DeleteAccount(ctx context.Context, rqst *ressource.DeleteAccountRqst) (*ressource.DeleteAccountRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	accountStr, _ := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.Id+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	account := make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(accountStr), &account)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	err = p.DeleteOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = p.Delete("local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete the token.
	err = p.DeleteOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+rqst.Id+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	name := account["name"].(string)
	name = strings.ReplaceAll(strings.ReplaceAll(name, ".", "_"), "@", "_")

	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		name)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, dropUserScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Remove the user database.
	err = p.DeleteDatabase("local_ressource", name+"_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = p.DeleteConnection(name + "_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteAccountRsp{
		Result: rqst.Id,
	}, nil
}

func (self *Globule) createRole(id string, name string, actions []string) error {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, _ := p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+id+`"}`, ``)
	if len(values) != 0 {
		return errors.New("Role named " + name + "already exist!")
	}

	// Here will create the new role.
	role := make(map[string]interface{}, 0)
	role["_id"] = id
	role["name"] = name
	role["actions"] = actions

	jsonStr, _ := Utility.ToJson(role)

	_, err = p.InsertOne("local_ressource", "local_ressource", "Roles", jsonStr, "")
	if err != nil {
		return err
	}

	return nil
}

//* Create a role with given action list *
func (self *Globule) CreateRole(ctx context.Context, rqst *ressource.CreateRoleRqst) (*ressource.CreateRoleRsp, error) {
	// That service made user of persistence service.
	err := self.createRole(rqst.Role.Id, rqst.Role.Name, rqst.Role.Actions)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.CreateRoleRsp{Result: true}, nil
}

//* Delete a role with a given id *
func (self *Globule) DeleteRole(ctx context.Context, rqst *ressource.DeleteRoleRqst) (*ressource.DeleteRoleRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, _ := p.Find("local_ressource", "local_ressource", "Accounts", `{}`, ``)
	if len(values) != 0 {
		accounts := make([]map[string]interface{}, 0)
		err := json.Unmarshal([]byte(values), &accounts)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		for i := 0; i < len(accounts); i++ {
			if accounts[i]["roles"] != nil {
				roles := accounts[i]["roles"].([]interface{})
				roles_ := make([]interface{}, 0)
				needSave := false
				for j := 0; j < len(roles); j++ {
					// TODO remove the role with name rqst.roleId from the account.
					role := roles[j].(map[string]interface{})
					if role["$id"] == rqst.RoleId {
						needSave = true
					} else {
						roles_ = append(roles_, role)
					}
				}

				// Here I will save the role.
				if needSave {
					accounts[i]["roles"] = roles_
					// Here I will save the role.
					jsonStr := "{"
					jsonStr += `"name":"` + accounts[i]["name"].(string) + `",`
					jsonStr += `"email":"` + accounts[i]["email"].(string) + `",`
					jsonStr += `"password":"` + accounts[i]["password"].(string) + `",`
					jsonStr += `"roles":[`
					for j := 0; j < len(accounts[i]["roles"].([]interface{})); j++ {
						db := accounts[i]["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
						db = strings.ReplaceAll(db, "@", "_")
						db = strings.ReplaceAll(db, ".", "_")
						jsonStr += `{`
						jsonStr += `"$ref":"` + accounts[i]["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
						jsonStr += `"$id":"` + accounts[i]["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
						jsonStr += `"$db":"` + db + `"`
						jsonStr += `}`
						if j < len(accounts[i]["roles"].([]interface{}))-1 {
							jsonStr += `,`
						}
					}
					jsonStr += `]`
					jsonStr += "}"

					err = p.ReplaceOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+accounts[i]["name"].(string)+`"}`, jsonStr, ``)
					if err != nil {
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}
		}
	}

	err = p.DeleteOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Delete permissions
	err = p.Delete("local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.RoleId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteRoleRsp{Result: true}, nil
}

//* Append an action to existing role. *
func (self *Globule) AddRoleAction(ctx context.Context, rqst *ressource.AddRoleActionRqst) (*ressource.AddRoleActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	role := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &role)

	needSave := false
	if role["actions"] == nil {
		role["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		for i := 0; i < len(role["actions"].([]interface{})); i++ {
			if role["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			role["actions"] = append(role["actions"].([]interface{}), rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+"already contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(role)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.AddRoleActionRsp{Result: true}, nil
}

//* Remove an action to existing role. *
func (self *Globule) RemoveRoleAction(ctx context.Context, rqst *ressource.RemoveRoleActionRqst) (*ressource.RemoveRoleActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	role := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &role)

	needSave := false
	if role["actions"] == nil {
		role["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		for i := 0; i < len(role["actions"].([]interface{})); i++ {
			if role["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, role["actions"].([]interface{})[i])
			}
		}
		if exist {
			role["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+"not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(role)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+rqst.RoleId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.RemoveRoleActionRsp{Result: true}, nil
}

//* Add role to a given account *
func (self *Globule) AddAccountRole(ctx context.Context, rqst *ressource.AddAccountRoleRqst) (*ressource.AddAccountRoleRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account named "+rqst.AccountId+" exist!")))
	}

	account := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(values), &account)

	// Now I will test if the account already contain the role.
	if account["roles"] != nil {
		for j := 0; j < len(account["roles"].([]interface{})); j++ {
			if account["roles"].([]interface{})[j].(map[string]interface{})["$id"] == rqst.RoleId {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Role named "+rqst.RoleId+" aleready exist in account "+rqst.AccountId+"!")))
			}
		}

		// append the newly created role.
		account["roles"] = append(account["roles"].([]interface{}), map[string]interface{}{"$ref": "Roles", "$id": rqst.RoleId, "$db": "local_ressource"})

		// Here I will save the role.
		jsonStr := "{"
		jsonStr += `"name":"` + account["name"].(string) + `",`
		jsonStr += `"email":"` + account["email"].(string) + `",`
		jsonStr += `"password":"` + account["password"].(string) + `",`
		jsonStr += `"roles":[`
		for j := 0; j < len(account["roles"].([]interface{})); j++ {
			db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
			db = strings.ReplaceAll(db, "@", "_")
			db = strings.ReplaceAll(db, ".", "_")
			jsonStr += `{`
			jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
			jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
			jsonStr += `"$db":"` + db + `"`
			jsonStr += `}`
			if j < len(account["roles"].([]interface{}))-1 {
				jsonStr += `,`
			}

		}
		jsonStr += `]`
		jsonStr += "}"

		err = p.ReplaceOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &ressource.AddAccountRoleRsp{Result: true}, nil
}

//* Remove a role from a given account *
func (self *Globule) RemoveAccountRole(ctx context.Context, rqst *ressource.RemoveAccountRoleRqst) (*ressource.RemoveAccountRoleRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No account named "+rqst.AccountId+" exist!")))
	}

	account := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(values), &account)

	// Now I will test if the account already contain the role.
	if account["roles"] != nil {
		roles := make([]interface{}, 0)
		needSave := false
		for j := 0; j < len(account["roles"].([]interface{})); j++ {
			if account["roles"].([]interface{})[j].(map[string]interface{})["$id"] == rqst.RoleId {
				needSave = true
			} else {
				roles = append(roles, account["roles"].([]interface{})[j])
			}
		}

		if needSave == false {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Account named "+rqst.AccountId+" does not contain role "+rqst.RoleId+"!")))
		}

		// append the newly created role.
		account["roles"] = roles

		// Here I will save the role.
		jsonStr := "{"
		jsonStr += `"name":"` + account["name"].(string) + `",`
		jsonStr += `"email":"` + account["email"].(string) + `",`
		jsonStr += `"password":"` + account["password"].(string) + `",`
		jsonStr += `"roles":[`
		for j := 0; j < len(account["roles"].([]interface{})); j++ {
			db := account["roles"].([]interface{})[j].(map[string]interface{})["$db"].(string)
			db = strings.ReplaceAll(db, "@", "_")
			db = strings.ReplaceAll(db, ".", "_")
			jsonStr += `{`
			jsonStr += `"$ref":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$ref"].(string) + `",`
			jsonStr += `"$id":"` + account["roles"].([]interface{})[j].(map[string]interface{})["$id"].(string) + `",`
			jsonStr += `"$db":"` + db + `"`
			jsonStr += `}`
			if j < len(account["roles"].([]interface{}))-1 {
				jsonStr += `,`
			}

		}
		jsonStr += `]`
		jsonStr += "}"

		err = p.ReplaceOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+rqst.AccountId+`"}`, jsonStr, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}
	return &ressource.RemoveAccountRoleRsp{Result: true}, nil
}

//* Append an action to existing application. *
func (self *Globule) AddApplicationAction(ctx context.Context, rqst *ressource.AddApplicationActionRqst) (*ressource.AddApplicationActionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	application := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &application)

	needSave := false
	if application["actions"] == nil {
		application["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		for i := 0; i < len(application["actions"].([]interface{})); i++ {
			if application["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
				break
			}
		}
		if !exist {
			application["actions"] = append(application["actions"].([]interface{}), rqst.Action)
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Application named "+rqst.ApplicationId+" already contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(application)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.AddApplicationActionRsp{Result: true}, nil
}

//* Remove an action to existing application. *
func (self *Globule) RemoveApplicationAction(ctx context.Context, rqst *ressource.RemoveApplicationActionRqst) (*ressource.RemoveApplicationActionRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	application := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &application)

	needSave := false
	if application["actions"] == nil {
		application["actions"] = []string{rqst.Action}
		needSave = true
	} else {
		exist := false
		actions := make([]interface{}, 0)
		for i := 0; i < len(application["actions"].([]interface{})); i++ {
			if application["actions"].([]interface{})[i].(string) == rqst.Action {
				exist = true
			} else {
				actions = append(actions, application["actions"].([]interface{})[i])
			}
		}
		if exist {
			application["actions"] = actions
			needSave = true
		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Application named "+rqst.ApplicationId+" not contain actions named "+rqst.Action+"!")))
		}
	}

	if needSave {
		jsonStr, _ := json.Marshal(application)
		err := p.ReplaceOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, string(jsonStr), ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.RemoveApplicationActionRsp{Result: true}, nil
}

//* Delete an application from the server. *
func (self *Globule) DeleteApplication(ctx context.Context, rqst *ressource.DeleteApplicationRqst) (*ressource.DeleteApplicationRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// Here I will test if a newer token exist for that user if it's the case
	// I will not refresh that token.
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, ``)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	application := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonStr), &application)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// First of all I will remove the directory.
	err = os.RemoveAll(application["path"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will remove the database create for the application.
	err = p.DeleteDatabase("local_ressource", rqst.ApplicationId+"_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Finaly I will remove the entry in  the table.
	err = p.DeleteOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+rqst.ApplicationId+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Delete permissions
	err = p.Delete("local_ressource", "local_ressource", "Permissions", `{"owner":"`+rqst.ApplicationId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Drop the application user.
	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		rqst.ApplicationId)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, dropUserScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteApplicationRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteAccountPermissions(name string) error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	err = p.Delete("local_ressource", "local_ressource", "Accounts", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given account *
func (self *Globule) DeleteAccountPermissions(ctx context.Context, rqst *ressource.DeleteAccountPermissionsRqst) (*ressource.DeleteAccountPermissionsRsp, error) {

	err := self.deleteAccountPermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteAccountPermissionsRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteRolePermissions(name string) error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	err = p.Delete("local_ressource", "local_ressource", "Roles", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	return nil
}

//* Delete all permission for a given role *
func (self *Globule) DeleteRolePermissions(ctx context.Context, rqst *ressource.DeleteRolePermissionsRqst) (*ressource.DeleteRolePermissionsRsp, error) {

	err := self.deleteRolePermissions(rqst.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteRolePermissionsRsp{
		Result: true,
	}, nil
}

/**
 * Return the list of all actions avalaible on the server.
 */
func (self *Globule) GetAllActions(ctx context.Context, rqst *ressource.GetAllActionsRqst) (*ressource.GetAllActionsRsp, error) {
	return &ressource.GetAllActionsRsp{Actions: self.methods}, nil
}

/////////////////////// Ressource permission owner /////////////////////////////
func (self *Globule) setRessourceOwner(owner string, path string) error {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// here I if the ressource is a directory I will set the permission on
	// subdirectory and files...
	fileInfo, err := os.Stat(self.GetAbsolutePath(path))
	if err == nil {
		if fileInfo.IsDir() {
			files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
			if err == nil {
				for i := 0; i < len(files); i++ {
					file := files[i]
					self.setRessourceOwner(owner, path+"/"+file.Name())
				}
			} else {
				return err
			}
		}
	}

	// Here I will set ressources whit that path, be sure to have different
	// path than application and webroot path if you dont want permission follow each other.
	ressources, err := self.getRessources(path)
	if err == nil {
		for i := 0; i < len(ressources); i++ {
			if ressources[i].GetPath() != path {
				path_ := ressources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					ressourceOwner := make(map[string]interface{})
					ressourceOwner["owner"] = owner
					ressourceOwner["path"] = path_

					// Here if the
					jsonStr, err := Utility.ToJson(&ressourceOwner)
					if err != nil {
						return err
					}

					err = p.ReplaceOne("local_ressource", "local_ressource", "RessourceOwners", jsonStr, jsonStr, `[{"upsert":true}]`)
					if err != nil {
						return err
					}
				}
			}
			self.setRessourceOwner(owner, ressources[i].GetPath()+"/"+ressources[i].GetName())
		}
	}

	// Here I will set the ressource owner.
	ressourceOwner := make(map[string]interface{})
	ressourceOwner["owner"] = owner
	ressourceOwner["path"] = path

	// Here if the
	jsonStr, err := Utility.ToJson(&ressourceOwner)
	if err != nil {
		return err
	}

	err = p.ReplaceOne("local_ressource", "local_ressource", "RessourceOwners", jsonStr, jsonStr, `[{"upsert":true}]`)
	if err != nil {
		return err
	}

	return nil
}

//* Set Ressource owner *
func (self *Globule) SetRessourceOwner(ctx context.Context, rqst *ressource.SetRessourceOwnerRqst) (*ressource.SetRessourceOwnerRsp, error) {
	// That service made user of persistence service.
	path := rqst.GetPath()

	err := self.setRessourceOwner(rqst.GetOwner(), path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.SetRessourceOwnerRsp{
		Result: true,
	}, nil
}

func (self *Globule) GetAbsolutePath(path string) string {

	path = strings.ReplaceAll(path, "\\", "/")
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

//* Get the ressource owner *
func (self *Globule) GetRessourceOwners(ctx context.Context, rqst *ressource.GetRessourceOwnersRqst) (*ressource.GetRessourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Get the absolute path
	path := rqst.GetPath()

	// find the ressource with it id
	ressourceOwnersStr, err := p.Find("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	ressourceOwners := make([]interface{}, 0)
	err = json.Unmarshal([]byte(ressourceOwnersStr), &ressourceOwners)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	owners := make([]string, 0)
	for i := 0; i < len(ressourceOwners); i++ {
		owners = append(owners, ressourceOwners[i].(map[string]interface{})["owner"].(string))
	}

	return &ressource.GetRessourceOwnersRsp{
		Owners: owners,
	}, nil
}

//* Get the ressource owner *
func (self *Globule) DeleteRessourceOwner(ctx context.Context, rqst *ressource.DeleteRessourceOwnerRqst) (*ressource.DeleteRessourceOwnerRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// Delete the ressource owner for a given path.
	err = p.DeleteOne("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`","owner":"`+rqst.GetOwner()+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteRessourceOwnerRsp{
		Result: true,
	}, nil
}

func (self *Globule) DeleteRessourceOwners(ctx context.Context, rqst *ressource.DeleteRessourceOwnersRqst) (*ressource.DeleteRessourceOwnersRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	path := rqst.GetPath()

	// delete the ressource owners with it path
	err = p.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteRessourceOwnersRsp{
		Result: true,
	}, nil
}

//////////////////////////// Loggin info ///////////////////////////////////////
func (self *Globule) logServiceInfo(service string, message string) error {

	// Here I will use event to publish log information...
	info := new(ressource.LogInfo)
	info.Application = ""
	info.UserId = "globular"
	info.UserName = "globular"
	info.Method = service
	info.Date = time.Now().Unix()
	info.Message = message
	info.Type = ressource.LogType_ERROR // not necessarely errors..
	self.log(info)

	return nil
}

// Log err and info...
func (self *Globule) logInfo(application string, method string, token string, err_ error) error {

	// Remove cyclic calls
	if method == "/ressource.RessourceService/Log" {
		return errors.New("Method " + method + " cannot not be log because it will cause a circular call to itself!")
	}

	// Here I will use event to publish log information...
	info := new(ressource.LogInfo)
	info.Application = application
	info.UserId = token
	info.UserName = token
	info.Method = method
	info.Date = time.Now().Unix()
	if err_ != nil {
		info.Message = err_.Error()
		info.Type = ressource.LogType_ERROR
	} else {
		info.Type = ressource.LogType_INFO
	}

	self.log(info)

	return nil
}

// unaryInterceptor calls authenticateClient with current context
func (self *Globule) unaryRessourceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := info.FullMethod

	// The token and the application id.
	var token string
	var application string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		application = strings.Join(md["application"], "")
		token = strings.Join(md["token"], "")
	}

	hasAccess := false
	var err error

	// Here some method are accessible by default.
	if method == "/ressource.RessourceService/GetAllActions" ||
		method == "/ressource.RessourceService/RegisterAccount" ||
		method == "/ressource.RessourceService/RegisterPeer" ||
		method == "/ressource.RessourceService/Authenticate" ||
		method == "/ressource.RessourceService/RefreshToken" ||
		method == "/ressource.RessourceService/GetPermissions" ||
		method == "/ressource.RessourceService/GetRessourceOwners" ||
		method == "/ressource.RessourceService/GetAllFilesInfo" ||
		method == "/ressource.RessourceService/GetAllApplicationsInfo" ||
		method == "/ressource.RessourceService/GetRessourceOwners" ||
		method == "/ressource.RessourceService/ValidateToken" ||
		method == "/ressource.RessourceService/ValidateUserRessourceAccess" ||
		method == "/ressource.RessourceService/ValidateApplicationRessourceAccess" ||
		method == "/ressource.RessourceService/ValidateUserRessourceAccess" ||
		method == "/ressource.RessourceService/ValidateApplicationAccess" ||
		method == "/ressource.RessourceService/GetActionPermission" ||
		method == "/ressource.RessourceService/Log" ||
		method == "/ressource.RessourceService/GetLog" {
		hasAccess = true
	}

	// Test if the user has access to execute the method
	if len(token) > 0 && !hasAccess {
		clientId, expiredAt, err := Interceptors.ValidateToken(token)

		if err != nil {
			return nil, err
		}

		if expiredAt < time.Now().Unix() {
			return nil, errors.New("The token is expired!")
		}
		if clientId == "sa" {
			hasAccess = true
			// log.Println("run ", method, application, clientId)
		} else {
			// special case that need ownership of the ressource or be sa
			if method == "/ressource.RessourceService/SetPermission" || method == "/ressource.RessourceService/DeletePermissions" ||
				method == "/ressource.RessourceService/SetRessourceOwner" || method == "/ressource.RessourceService/DeleteRessourceOwner" {
				var path string
				if method == "/ressource.RessourceService/SetPermission" {
					rqst := req.(*ressource.SetPermissionRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Permission.GetPath()
				} else if method == "/ressource.RessourceService/DeletePermissions" {
					rqst := req.(*ressource.DeletePermissionsRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				} else if method == "/ressource.RessourceService/SetRessourceOwner" {
					rqst := req.(*ressource.SetRessourceOwnerRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				} else if method == "/ressource.RessourceService/DeleteRessourceOwner" {
					rqst := req.(*ressource.DeleteRessourceOwnerRqst)
					// Here I will validate that the user is the owner.
					path = rqst.Path
				}

				// If the use is the ressource owner he can run the method
				if self.isOwner(clientId, path) {
					hasAccess = true
				}

			} else {
				err = self.validateUserAccess(clientId, method)
				if err == nil {
					hasAccess = true
				}
			}
		}
	}

	// Test if the application has access to execute the method.
	if len(application) > 0 && !hasAccess {
		err = self.validateApplicationAccess(application, method)
		if err == nil {
			hasAccess = true
		}
	}

	if !hasAccess {
		err := errors.New("Permission denied to execute method " + method)
		self.logInfo(application, method, token, err)
		return nil, err
	}

	// Execute the action.
	result, err := handler(ctx, req)
	self.logInfo(application, method, token, err)

	if err == nil {
		// Set permissions in case one of those methode is called.
		if method == "/ressource.RessourceService/DeleteApplication" {
			rqst := req.(*ressource.DeleteApplicationRqst)
			err := self.deleteDirPermissions("/" + rqst.ApplicationId)
			if err != nil {
				log.Println(err)
			}
		} else if method == "/ressource.RessourceService/DeleteRole" {
			rqst := req.(*ressource.DeleteRoleRqst)
			err := self.deleteRolePermissions("/" + rqst.RoleId)
			if err != nil {
				log.Println(err)
			}
		} else if method == "/ressource.RessourceService/DeleteAccount" {
			rqst := req.(*ressource.DeleteAccountRqst)
			err := self.deleteAccountPermissions("/" + rqst.Id)
			if err != nil {
				log.Println(err)
			}
		}
	}

	return result, err

}

// Stream interceptor.
func (self *Globule) streamRessourceInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	err := handler(srv, stream)
	if err != nil {
		return err
	}

	return nil
}

func (self *Globule) getLogInfoKeyValue(info *ressource.LogInfo) (string, string, error) {
	marshaler := new(jsonpb.Marshaler)
	jsonStr, err := marshaler.MarshalToString(info)
	if err != nil {
		return "", "", err
	}

	key := ""
	if info.GetType() == ressource.LogType_INFO {
		// Increnment prometheus counter,
		self.methodsCounterLog.WithLabelValues("INFO", info.Method).Inc()

		// Append the log in leveldb
		key += "/infos/" + info.Method + Utility.ToString(info.Date)

		// Set the application in the path
		if len(info.Application) > 0 {
			key += "/" + info.Application
		}
		// Set the User Name if available.
		if len(info.UserName) > 0 {
			key += "/" + info.UserName
		}

		key += "/" + Utility.GenerateUUID(jsonStr)

	} else {
		// Increnment prometheus counter,
		self.methodsCounterLog.WithLabelValues("ERROR", info.Method).Inc()
		key += "/errors/" + info.Method + Utility.ToString(info.Date)

		// Set the application in the path
		if len(info.Application) > 0 {
			key += "/" + info.Application
		}
		// Set the User Name if available.
		if len(info.UserName) > 0 {
			key += "/" + info.UserName
		}

		key += "/" + Utility.GenerateUUID(jsonStr)

	}
	return key, jsonStr, nil
}

func (self *Globule) log(info *ressource.LogInfo) error {

	// The userId can be a single string or a JWT token.
	if len(info.UserName) > 0 {
		name, _, err := Interceptors.ValidateToken(info.UserName)
		if err == nil {
			info.UserName = name
		}
		info.UserId = info.UserName // keep only the user name
		if info.UserName == "sa" {
			return nil // not log sa activities...
		}
	} else {
		return nil
	}

	key, jsonStr, err := self.getLogInfoKeyValue(info)
	if err != nil {
		return err
	}

	// Append the error in leveldb
	self.logs.SetItem(key, []byte(jsonStr))
	eventHub, err := self.getEventHub()
	if err != nil {
		return err
	}
	eventHub.Publish(info.Method, []byte(jsonStr))

	return nil
}

//* Log error or information into the data base *
func (self *Globule) Log(ctx context.Context, rqst *ressource.LogRqst) (*ressource.LogRsp, error) {
	// Publish event...
	self.log(rqst.Info)

	return &ressource.LogRsp{
		Result: true,
	}, nil
}

//* Log error or information into the data base *
// Retreive log infos (the query must be something like /infos/'date'/'applicationName'/'userName'
func (self *Globule) GetLog(rqst *ressource.GetLogRqst, stream ressource.RessourceService_GetLogServer) error {

	query := rqst.Query
	if len(query) == 0 {
		query = "/*"
	}

	data, err := self.logs.GetItem(query)

	jsonDecoder := json.NewDecoder(strings.NewReader(string(data)))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return err
	}

	infos := make([]*ressource.LogInfo, 0)
	i := 0
	max := 100
	for jsonDecoder.More() {
		info := ressource.LogInfo{}
		err := jsonpb.UnmarshalNext(jsonDecoder, &info)
		if err != nil {
			log.Fatal(err)
		}
		// append the info inside the stream.
		infos = append(infos, &info)
		if i == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &ressource.GetLogRsp{
				Info: infos,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			infos = make([]*ressource.LogInfo, 0)
			i = 0
		}
		i++
	}

	// Send the last infos...
	if len(infos) > 0 {
		rsp := &ressource.GetLogRsp{
			Info: infos,
		}
		err = stream.Send(rsp)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return nil
}

func (self *Globule) deleteLog(query string) error {

	// First of all I will retreive the log info with a given date.
	data, err := self.logs.GetItem(query)

	jsonDecoder := json.NewDecoder(strings.NewReader(string(data)))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return err
	}

	for jsonDecoder.More() {
		info := ressource.LogInfo{}

		err := jsonpb.UnmarshalNext(jsonDecoder, &info)
		if err != nil {
			return err
		}

		key, _, err := self.getLogInfoKeyValue(&info)
		if err != nil {
			return err
		}
		self.logs.RemoveItem(key)

	}

	return nil
}

//* Delete a log info *
func (self *Globule) DeleteLog(ctx context.Context, rqst *ressource.DeleteLogRqst) (*ressource.DeleteLogRsp, error) {

	key, _, _ := self.getLogInfoKeyValue(rqst.Log)
	err := self.logs.RemoveItem(key)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteLogRsp{
		Result: true,
	}, nil
}

//* Clear logs. info or errors *
func (self *Globule) ClearAllLog(ctx context.Context, rqst *ressource.ClearAllLogRqst) (*ressource.ClearAllLogRsp, error) {
	var err error

	if rqst.Type == ressource.LogType_ERROR {
		err = self.deleteLog("/errors/*")
	} else {
		err = self.deleteLog("/infos/*")
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.ClearAllLogRsp{
		Result: true,
	}, nil
}

///////////////////////  ressource management. /////////////////

//* Set a ressource from a client (custom service) to globular
func (self *Globule) SetRessource(ctx context.Context, rqst *ressource.SetRessourceRqst) (*ressource.SetRessourceRsp, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(rqst.Ressource)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(rqst.Ressource.Path + rqst.Ressource.Name)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Always create a new if not already exist.
	err = p.ReplaceOne("local_ressource", "local_ressource", "Ressources", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.SetRessourceRsp{
		Result: true,
	}, nil
}

func (self *Globule) getRessources(path string) ([]*ressource.Ressource, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	data, err := p.Find("local_ressource", "local_ressource", "Ressources", `{}`, `[{"Projection":{"_id":0}}]`)
	if err != nil {
		return nil, err
	}

	jsonDecoder := json.NewDecoder(strings.NewReader(data))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return nil, err
	}

	ressources := make([]*ressource.Ressource, 0)

	for jsonDecoder.More() {
		res := new(ressource.Ressource)
		err := jsonpb.UnmarshalNext(jsonDecoder, res)
		if err != nil {
			return nil, err
		}
		// append the info inside the stream.
		if strings.HasPrefix(res.GetPath(), path) {
			ressources = append(ressources, res)
		}
	}
	return ressources, nil
}

//* Get all ressources
func (self *Globule) GetRessources(rqst *ressource.GetRessourcesRqst, stream ressource.RessourceService_GetRessourcesServer) error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	query := make(map[string]string)
	if len(rqst.Name) > 0 {
		query["name"] = rqst.Name
	}

	if len(rqst.Path) > 0 {
		query["path"] = rqst.Path
	}

	queryStr, _ := Utility.ToJson(query)

	data, err := p.Find("local_ressource", "local_ressource", "Ressources", queryStr, `[{"Projection":{"_id":0}}]`)
	if err != nil {
		return err
	}

	jsonDecoder := json.NewDecoder(strings.NewReader(data))

	// read open bracket
	_, err = jsonDecoder.Token()
	if err != nil {
		return err
	}

	ressources := make([]*ressource.Ressource, 0)
	i := 0
	max := 100
	for jsonDecoder.More() {
		res := new(ressource.Ressource)
		err := jsonpb.UnmarshalNext(jsonDecoder, res)
		if err != nil {
			log.Fatal(err)
		}
		// append the info inside the stream.
		ressources = append(ressources, res)
		if i == max-1 {
			// I will send the stream at each 100 logs...
			rsp := &ressource.GetRessourcesRsp{
				Ressources: ressources,
			}
			// Send the infos
			err = stream.Send(rsp)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			ressources = make([]*ressource.Ressource, 0)
			i = 0
		}
		i++
	}

	// Send the last infos...
	if len(ressources) > 0 {
		rsp := &ressource.GetRessourcesRsp{
			Ressources: ressources,
		}
		err = stream.Send(rsp)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return nil
}

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) RemoveRessource(ctx context.Context, rqst *ressource.RemoveRessourceRqst) (*ressource.RemoveRessourceRsp, error) {

	// Because regex dosent work properly I retreive all the ressources.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// get all ressource with that path.
	ressources, err := self.getRessources(rqst.Ressource.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	toDelete := make([]*ressource.Ressource, 0)
	// Remove ressource that match...
	for i := 0; i < len(ressources); i++ {
		res := ressources[i]
		// In case the ressource is a sub-ressource I will remove it...
		if len(rqst.Ressource.Name) > 0 {
			if rqst.Ressource.Name == res.GetName() {
				toDelete = append(toDelete, res) // mark to be delete.
			}
		} else {
			toDelete = append(toDelete, res) // mark to be delete
		}

	}

	// Now I will delete the ressource.
	for i := 0; i < len(toDelete); i++ {
		id := Utility.GenerateUUID(toDelete[i].Path + toDelete[i].Name)
		err := p.DeleteOne("local_ressource", "local_ressource", "Ressources", `{"_id":"`+id+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Delete the permissions ascosiated permission
		self.deletePermissions(toDelete[i].Path+"/"+toDelete[i].Name, "")
		err = p.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+toDelete[i].Path+"/"+toDelete[i].Name+`"}`, "")
		if err != nil {
			log.Println(err)
		}
	}

	// In that case the
	if len(rqst.Ressource.Name) == 0 {
		self.deletePermissions(rqst.Ressource.Path, "")
		err = p.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+rqst.Ressource.Path+`"}`, "")
		if err != nil {
			log.Println(err)
		}
	}

	return &ressource.RemoveRessourceRsp{
		Result: true,
	}, nil
}

//* Set a ressource from a client (custom service) to globular
func (self *Globule) SetActionPermission(ctx context.Context, rqst *ressource.SetActionPermissionRqst) (*ressource.SetActionPermissionRsp, error) {

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	actionPermission := make(map[string]interface{}, 0)
	actionPermission["action"] = rqst.Action
	actionPermission["permission"] = rqst.Permission
	actionPermission["_id"] = Utility.GenerateUUID(rqst.Action)

	actionPermissionStr, _ := Utility.ToJson(actionPermission)
	err = p.ReplaceOne("local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, actionPermissionStr, `[{"upsert":true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.SetActionPermissionRsp{
		Result: true,
	}, nil
}

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) RemoveActionPermission(ctx context.Context, rqst *ressource.RemoveActionPermissionRqst) (*ressource.RemoveActionPermissionRsp, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	err = p.DeleteOne("local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.RemoveActionPermissionRsp{
		Result: true,
	}, nil
}

//* Remove a ressource from a client (custom service) to globular
func (self *Globule) GetActionPermission(ctx context.Context, rqst *ressource.GetActionPermissionRqst) (*ressource.GetActionPermissionRsp, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Try to delete the account...
	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "ActionPermission", `{"_id":"`+Utility.GenerateUUID(rqst.Action)+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	actionPermission := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonStr), &actionPermission)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.GetActionPermissionRsp{
		Permission: int32(actionPermission["permission"].(float64)),
	}, nil
}

func (self *Globule) savePermission(owner string, path string, permission int32) error {

	// dir cannot be executable....
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Here I will insert one or replcace one depending if permission already exist or not.
	query := `{"owner":"` + owner + `","path":"` + path + `"}`
	jsonStr := `{"owner":"` + owner + `","path":"` + path + `","permission":` + Utility.ToString(permission) + `}`

	count, _ := p.Count("local_ressource", "local_ressource", "Permissions", query, "")

	if count == 0 {
		_, err = p.InsertOne("local_ressource", "local_ressource", "Permissions", jsonStr, "")

	} else {
		err = p.ReplaceOne("local_ressource", "local_ressource", "Permissions", query, jsonStr, "")
	}

	return err
}

// Set directory permission
func (self *Globule) setDirPermission(owner string, path string, permission int32) error {
	log.Println("Set dir permission ", path)

	err := self.savePermission(owner, path, permission)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
	if err != nil {
		return err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			err := self.setDirPermission(owner, path+"/"+file.Name(), permission)
			if err != nil {
				return err
			}
		} else {
			err := self.setRessourcePermission(owner, path+"/"+file.Name(), permission)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Set file permission.
func (self *Globule) setRessourcePermission(owner string, path string, permission int32) error {
	return self.savePermission(owner, path, permission)
}

//* Set a file permission, create new one if not already exist. *
func (self *Globule) SetPermission(ctx context.Context, rqst *ressource.SetPermissionRqst) (*ressource.SetPermissionRsp, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The first thing I will do is test if the file exist.
	path := rqst.GetPermission().GetPath()
	path = strings.ReplaceAll(path, "\\", "/")

	// Now if the permission exist I will read the file info.
	fileInfo, _ := os.Stat(self.GetAbsolutePath(path))

	// Now I will test if the user or the role exist.
	owner := make(map[string]interface{})

	switch v := rqst.Permission.GetOwner().(type) {
	case *ressource.RessourcePermission_User:
		// In that case I will try to find a user with that id
		jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+v.User+`"}`, ``)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		json.Unmarshal([]byte(jsonStr), &owner)
	case *ressource.RessourcePermission_Role:
		// In that case I will try to find a role with that id
		jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+v.Role+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		json.Unmarshal([]byte(jsonStr), &owner)
	case *ressource.RessourcePermission_Application:
		// In that case I will try to find a role with that id
		jsonStr, err := p.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+v.Application+`"}`, "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		json.Unmarshal([]byte(jsonStr), &owner)
	}

	if fileInfo != nil {

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// do directory stuff
			err := self.setDirPermission(owner["_id"].(string), path, rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		case mode.IsRegular():
			// do file stuff
			err := self.setRessourcePermission(owner["_id"].(string), path, rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}
	}

	ressources, err := self.getRessources(path)
	if err == nil {
		for i := 0; i < len(ressources); i++ {
			if ressources[i].GetPath() != path {
				path_ := ressources[i].GetPath()[len(path)+1:]
				paths := strings.Split(path_, "/")
				path_ = path
				// set sub-path...
				for j := 0; j < len(paths); j++ {
					path_ += "/" + paths[j]
					err := self.setRessourcePermission(owner["_id"].(string), path_, rqst.Permission.Number)
					if err != nil {
						return nil, status.Errorf(
							codes.Internal,
							Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
					}
				}
			}

			// create ressource permission
			err := self.setRessourcePermission(owner["_id"].(string), ressources[i].GetPath()+"/"+ressources[i].GetName(), rqst.Permission.Number)
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

		}
		// save ressource path.
		err = self.setRessourcePermission(owner["_id"].(string), path, rqst.Permission.Number)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	}

	return &ressource.SetPermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) setPermissionOwner(owner string, permission *ressource.RessourcePermission) error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Here I will try to find the owner in the user table
	_, err = p.FindOne("local_ressource", "local_ressource", "Accounts", `{"_id":"`+owner+`"}`, ``)
	if err == nil {
		permission.Owner = &ressource.RessourcePermission_User{
			User: owner,
		}
		return nil
	}

	// In the role.
	// In that case I will try to find a role with that id
	_, err = p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &ressource.RessourcePermission_Role{
			Role: owner,
		}
		return nil
	}

	_, err = p.FindOne("local_ressource", "local_ressource", "Applications", `{"_id":"`+owner+`"}`, "")
	if err == nil {
		permission.Owner = &ressource.RessourcePermission_Application{
			Application: owner,
		}
		return nil
	}

	return errors.New("No Role or User found with id " + owner)
}

/**
 *
 */
func (self *Globule) getDirPermissions(path string) ([]*ressource.RessourcePermission, error) {
	if !Utility.Exists(self.GetAbsolutePath(path)) {
		return nil, errors.New("No directory found with path " + self.GetAbsolutePath(path))
	}

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	jsonStr, err := p.Find("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Now I will read the permission values.
	permissions_ := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonStr), &permissions_)
	if err != nil {
		return nil, err
	}

	permissions := make([]*ressource.RessourcePermission, 0)
	for i := 0; i < len(permissions_); i++ {
		permission_ := permissions_[i]
		permission := &ressource.RessourcePermission{Path: permission_["path"].(string), Owner: nil, Number: int32(Utility.ToInt(permission_["permission"].(float64)))}
		err = self.setPermissionOwner(permission_["owner"].(string), permission)
		if err != nil {
			return nil, err
		}

		// append into the permissions.
		permissions = append(permissions, permission)
	}

	// No the recursion.
	files, err := ioutil.ReadDir(self.GetAbsolutePath(path))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(files); i++ {
		file := files[i]
		if file.IsDir() {
			permissions_, err := self.getDirPermissions(path + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			// append to the permissions
			permissions = append(permissions, permissions_...)
		} else {
			permissions_, err := self.getRessourcePermissions(path + "/" + file.Name())
			if err != nil {
				return nil, err
			}
			// append to the permissions
			permissions = append(permissions, permissions_...)
		}
	}

	return permissions, nil
}

func (self *Globule) getRessourcePermissions(path string) ([]*ressource.RessourcePermission, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	jsonStr, err := p.Find("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	permissions_ := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonStr), &permissions_)
	if err != nil {
		return nil, err
	}

	permissions := make([]*ressource.RessourcePermission, 0)

	for i := 0; i < len(permissions_); i++ {
		permission_ := permissions_[i]
		permission := &ressource.RessourcePermission{Path: permission_["path"].(string), Owner: nil, Number: int32(Utility.ToInt(permission_["permission"].(float64)))}
		err = self.setPermissionOwner(permission_["owner"].(string), permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (self *Globule) getPermissions(path string) ([]*ressource.RessourcePermission, error) {

	fileInfo, err := os.Stat(self.GetAbsolutePath(path))
	if err == nil {

		switch mode := fileInfo.Mode(); {
		case mode.IsDir():
			// do directory stuff
			permissions, err := self.getDirPermissions(path)
			if err != nil {
				return nil, err
			}

			return permissions, nil

		case mode.IsRegular():
			// do file stuff
			permissions, err := self.getRessourcePermissions(path)
			if err != nil {
				return nil, err
			}

			return permissions, nil
		}
	} else {
		// do file stuff
		permissions, err := self.getRessourcePermissions(path)
		if err != nil {
			return nil, err
		}

		return permissions, nil
	}

	return nil, nil
}

//* Get All permissions for a given file/dir *
func (self *Globule) GetPermissions(ctx context.Context, rqst *ressource.GetPermissionsRqst) (*ressource.GetPermissionsRsp, error) {

	path := rqst.GetPath()
	permissions, err := self.getPermissions(path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	permissions_ := make([]map[string]interface{}, len(permissions))
	for i := 0; i < len(permissions); i++ {
		permissions_[i] = make(map[string]interface{}, 0)
		// Set the values.
		permissions_[i]["path"] = permissions[i].GetPath()
		permissions_[i]["number"] = permissions[i].GetNumber()
		permissions_[i]["user"] = permissions[i].GetUser()
		permissions_[i]["role"] = permissions[i].GetRole()
		permissions_[i]["application"] = permissions[i].GetApplication()
	}

	jsonStr, err := json.Marshal(&permissions_)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.GetPermissionsRsp{
		Permissions: string(jsonStr),
	}, nil
}

func (self *Globule) deletePermissions(path string, owner string) error {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// First of all I will retreive the permissions for the given path...
	permissions, err := self.getPermissions(path)
	if err != nil {
		return err
	}

	// Get list of all permission with a given path.
	for i := 0; i < len(permissions); i++ {
		permission := permissions[i]
		if len(owner) > 0 {
			switch v := permission.GetOwner().(type) {
			case *ressource.RessourcePermission_User:
				if v.User == owner {
					err := p.DeleteOne("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.User+`"}`, "")
					if err != nil {
						return err
					}
				}
			case *ressource.RessourcePermission_Role:
				if v.Role == owner {
					err := p.DeleteOne("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Role+`"}`, "")
					if err != nil {
						return err
					}
				}

			case *ressource.RessourcePermission_Application:
				if v.Application == owner {
					err := p.DeleteOne("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`","owner":"`+v.Application+`"}`, "")
					if err != nil {
						return err
					}
				}
			}
		} else {
			err := p.DeleteOne("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission.GetPath()+`"}`, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//* Delete a file permission *
func (self *Globule) DeletePermissions(ctx context.Context, rqst *ressource.DeletePermissionsRqst) (*ressource.DeletePermissionsRsp, error) {

	// That service made user of persistence service.
	err := self.deletePermissions(rqst.GetPath(), rqst.GetOwner())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeletePermissionsRsp{
		Result: true,
	}, nil
}

//* Create Permission for a dir (recursive) *
func (self *Globule) CreateDirPermissions(ctx context.Context, rqst *ressource.CreateDirPermissionsRqst) (*ressource.CreateDirPermissionsRsp, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	clientId, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	// A new directory will take the parent permissions by default...
	path := rqst.GetPath()

	permissionsStr, err := p.Find("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Create permission object.
	permissions := make([]interface{}, 0)
	err = json.Unmarshal([]byte(permissionsStr), &permissions)
	if err != nil {
		return nil, err
	}

	// Now I will create the new permission of the created directory.
	for i := 0; i < len(permissions); i++ {
		// Copye the permission.
		permission := permissions[i].(map[string]interface{})
		permission_ := make(map[string]interface{}, 0)
		permission_["owner"] = permission["owner"]
		permission_["path"] = path + "/" + rqst.GetName()
		permission_["permission"] = permission["permission"]
		permissionStr, _ := Utility.ToJson(permission_)
		p.InsertOne("local_ressource", "local_ressource", "Permissions", permissionStr, "")
	}

	ressourceOwnersStr, err := p.Find("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		return nil, err
	}

	// Create permission object.
	ressourceOwners := make([]interface{}, 0)
	err = json.Unmarshal([]byte(ressourceOwnersStr), &ressourceOwners)
	if err != nil {
		return nil, err
	}

	// Now I will create the new permission of the created directory.
	for i := 0; i < len(ressourceOwners); i++ {
		// Copye the permission.
		ressourceOwner := ressourceOwners[i].(map[string]interface{})
		ressourceOwner_ := make(map[string]interface{}, 0)
		ressourceOwner_["owner"] = ressourceOwner["owner"]
		ressourceOwner_["path"] = path + "/" + rqst.GetName()
		ressourceOwnerStr, _ := Utility.ToJson(ressourceOwner_)
		p.InsertOne("local_ressource", "local_ressource", "RessourceOwners", ressourceOwnerStr, "")
	}

	// The user who create a directory will be the owner of the
	// directory.
	if clientId != "sa" && clientId != "guest" {
		ressourceOwner := make(map[string]interface{}, 0)
		ressourceOwner["owner"] = clientId
		ressourceOwner["path"] = path + "/" + rqst.GetName()
		ressourceOwnerStr, _ := Utility.ToJson(ressourceOwner)
		p.ReplaceOne("local_ressource", "local_ressource", "RessourceOwners", ressourceOwnerStr, ressourceOwnerStr, `[{"upsert":true}]`)
	}

	return &ressource.CreateDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Rename file/dir permission *
func (self *Globule) RenameFilePermission(ctx context.Context, rqst *ressource.RenameFilePermissionRqst) (*ressource.RenameFilePermissionRsp, error) {

	path := rqst.GetPath()
	path = strings.ReplaceAll(path, "\\", "/")

	oldPath := rqst.OldName
	newPath := rqst.NewName

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if strings.HasPrefix(path, "/") {
		if len(path) > 1 {
			oldPath = path + "/" + rqst.OldName
			newPath = path + "/" + rqst.NewName

		} else {
			oldPath = "/" + rqst.OldName
			newPath = "/" + rqst.NewName
		}
	}

	// Replace permission path... regex "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"} not work.
	permissionsStr, err := p.Find("local_ressource", "local_ressource", "Permissions", `{}`, "")
	if err == nil {
		permissions := make([]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				err := p.Update("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissionsStr, err = p.Find("local_ressource", "local_ressource", "RessourceOwners", `{}`, "")
	if err == nil {
		permissions := make([]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), oldPath) {
				path := newPath + permission["path"].(string)[len(oldPath):]
				log.Println("rename file permission owner", permission["path"].(string), " with ", newPath)
				err = p.Update("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+permission["path"].(string)+`"}`, `{"$set":{"path":"`+path+`"}}`, "")
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &ressource.RenameFilePermissionRsp{
		Result: true,
	}, nil
}

func (self *Globule) deleteDirPermissions(path string) error {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	path = strings.ReplaceAll(path, "\\", "/")

	// Replace permission path... regex "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"} not work.
	permissionsStr, err := p.Find("local_ressource", "local_ressource", "Permissions", `{}`, "")
	if err == nil {
		permissions := make([]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err := p.Delete("local_ressource", "local_ressource", "Permissions", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	// Replace file owner path... regex not work... "path":{"$regex":"/^`+strings.ReplaceAll(oldPath, "/", "\\/")+`.*/"}
	permissionsStr, err = p.Find("local_ressource", "local_ressource", "RessourceOwners", `{}`, "")
	if err == nil {
		permissions := make([]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		for i := 0; i < len(permissions); i++ { // stringnify and save it...
			permission := permissions[i].(map[string]interface{})
			if strings.HasPrefix(permission["path"].(string), path) {
				err = p.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+permission["path"].(string)+`"}`, "")
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}

//* Delete Permission for a dir (recursive) *
func (self *Globule) DeleteDirPermissions(ctx context.Context, rqst *ressource.DeleteDirPermissionsRqst) (*ressource.DeleteDirPermissionsRsp, error) {
	err := self.deleteDirPermissions(rqst.GetPath())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteDirPermissionsRsp{
		Result: true,
	}, nil
}

//* Delete a single file permission *
func (self *Globule) DeleteFilePermissions(ctx context.Context, rqst *ressource.DeleteFilePermissionsRqst) (*ressource.DeleteFilePermissionsRsp, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	path := rqst.GetPath()

	err = p.Delete("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, "")
	if err != nil {
		log.Println(err)
	}

	err = p.Delete("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`"}`, "")
	if err != nil {
		log.Println(err)
	}

	return &ressource.DeleteFilePermissionsRsp{
		Result: true,
	}, nil
}

//* Validate a token *
func (self *Globule) ValidateToken(ctx context.Context, rqst *ressource.ValidateTokenRqst) (*ressource.ValidateTokenRsp, error) {
	clientId, expireAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return &ressource.ValidateTokenRsp{
		ClientId: clientId,
		Expired:  expireAt,
	}, nil
}

/**
 * Validate application access by role
 */
func (self *Globule) validateApplicationAccess(name string, method string) error {
	//log.Println("-------> validate Application "+name+" for method ", method)
	if len(name) == 0 {
		return errors.New("No application was given to validate method access " + method)
	}

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	values, err := p.FindOne("local_ressource", "local_ressource", "Applications", `{"path":"/`+name+`"}`, ``)
	if err != nil {
		log.Println(err)
		return err
	}

	application := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &application)
	if err != nil {
		log.Println(err)
		return err
	}

	err = errors.New("permission denied! application " + name + " cannot execute methode '" + method + "'")
	if application["actions"] == nil {
		return err
	}

	actions := application["actions"].([]interface{})
	if actions == nil {
		return err
	}

	for i := 0; i < len(actions); i++ {
		if actions[i].(string) == method {
			return nil
		}
	}

	return err
}

/**
 * Validate user access by role
 */
func (self *Globule) validateUserAccess(userName string, method string) error {
	log.Println("---> validate user access ", userName, " for method ", method)
	if len(userName) == 0 {
		return errors.New("No user  name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	// if guest can run the action...
	if self.canRunAction("guest", method) == nil {
		// everybody can run the action in that case.
		return nil
	}

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		log.Println(err)
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		log.Println(err)
		return err
	}

	roles := account["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		if self.canRunAction(role["$id"].(string), method) == nil {
			return nil
		}
	}

	err = errors.New("permission denied! account " + userName + " cannot execute methode '" + method + "'")
	return err
}

// Test if a role can use action.
func (self *Globule) canRunAction(roleName string, method string) error {

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	values, err := p.FindOne("local_ressource", "local_ressource", "Roles", `{"_id":"`+roleName+`"}`, `[{"Projection":{"actions":1}}]`)
	if err != nil {
		return err
	}

	role := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &role)
	if err != nil {
		return err
	}

	// append all action into the actions
	for i := 0; i < len(role["actions"].([]interface{})); i++ {
		if strings.ToLower(role["actions"].([]interface{})[i].(string)) == strings.ToLower(method) {
			return nil
		}
	}

	// Here I will test if the user has write to execute the methode.
	return errors.New("Permission denied!")
}

// authenticateAgent check the client credentials
func (self *Globule) authenticateClient(ctx context.Context) (string, string, int64, error) {
	var userId string
	var applicationId string
	var expired int64
	var err error

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		applicationId = strings.Join(md["application"], "")
		token := strings.Join(md["token"], "")
		// In that case no token was given...
		if len(token) > 0 {
			userId, expired, err = Interceptors.ValidateToken(token)
		}
		return applicationId, userId, expired, err
	}
	return "", "", 0, fmt.Errorf("missing credentials")
}

func (self *Globule) isOwner(name string, path string) bool {

	// get the client...
	client, err := self.getPersistenceSaConnection()
	if err != nil {
		return false
	}

	// Now I will get the user roles and validate if the user can execute the
	// method.
	values, err := client.FindOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+name+`"}`, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return false
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return false
	}

	// If the user is the owner of the ressource it has the permission
	count, err := client.Count("local_ressource", "local_ressource", "RessourceOwners", `{"path":"`+path+`","owner":"`+account["_id"].(string)+`"}`, ``)
	if err == nil {
		if count > 0 {
			return true
		}
	} else {
		log.Println(err)
	}
	return false
}

func (self *Globule) hasPermission(name string, path string, permission int32) (bool, int) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return false, 0
	}

	// If the user is the owner of the ressource it has all permission
	if self.isOwner(name, path) {
		return true, 0
	}

	count, err := p.Count("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		return false, 0
	}

	permissionsStr, err := p.FindOne("local_ressource", "local_ressource", "Permissions", `{"owner":"`+name+`", "path":"`+path+`"}`, ``)
	if err == nil {
		permissions := make(map[string]interface{}, 0)
		json.Unmarshal([]byte(permissionsStr), &permissions)
		if len(permissions) == 0 {
			return false, count
		}

		p := int32(permissions["permission"].(float64))

		// Here the owner have all permissions.
		if p == 7 {
			return true, count
		}

		if permission == p {
			return true, count
		}

		// Delete permission
		if permission == 1 {
			if p == 1 || p == 3 || p == 5 {
				return true, count
			}
		}

		// Write permission
		if permission == 2 {
			if p == 2 || p == 3 || p == 6 {
				return true, count
			}
		}

		// Read permission
		if permission == 4 {
			if p == 4 || p == 5 || p == 6 {
				return true, count
			}
		}

		return false, count
	}

	return false, count
}

/**
 * Validate if a user, a role or an application has write to do operation on a file or a directorty.
 */
func (self *Globule) validateUserRessourceAccess(userName string, method string, path string, permission int32) error {
	log.Println("---> validate user access ", userName, method, path, permission)
	if len(userName) == 0 {
		return errors.New("No user name was given to validate method access " + method)
	}

	// if the user is the super admin no validation is required.
	if userName == "sa" {
		return nil
	}

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	// Find the user role.
	values, err := p.FindOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+userName+`"}`, `[{"Projection":{"roles":1}}]`)
	if err != nil {
		return err
	}

	account := make(map[string]interface{})
	err = json.Unmarshal([]byte(values), &account)
	if err != nil {
		return err
	}

	count := 0
	hasUserPermission, hasUserPermissionCount := self.hasPermission(userName, path, permission)
	if hasUserPermission {
		log.Println("---> user has permission ", userName, permission)
		return nil
	}

	count += hasUserPermissionCount
	roles := account["roles"].([]interface{})
	for i := 0; i < len(roles); i++ {
		role := roles[i].(map[string]interface{})
		hasRolePermission, hasRolePermissionCount := self.hasPermission(role["$id"].(string), path, permission)
		count += hasRolePermissionCount
		if hasRolePermission {
			log.Println("---> role has permission ", role["$id"].(string), permission)
			return nil
		}
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the ressource.
	if count > 0 {
		return errors.New("Permission Denied for " + userName)
	}

	count, err = p.Count("local_ressource", "local_ressource", "Permissions", `{"path":"`+path+`"}`, ``)
	if err != nil {
		if count > 0 {
			return errors.New("Permission Denied for " + userName)
		}
	}
	return nil
}

//* Validate if user can access a given file. *
func (self *Globule) ValidateUserRessourceAccess(ctx context.Context, rqst *ressource.ValidateUserRessourceAccessRqst) (*ressource.ValidateUserRessourceAccessRsp, error) {

	path := rqst.GetPath() // The path of the ressource.

	// first of all I will validate the token.
	clientId, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	err = self.validateUserRessourceAccess(clientId, rqst.Method, path, rqst.Permission)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.ValidateUserRessourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given file. *
func (self *Globule) ValidateApplicationRessourceAccess(ctx context.Context, rqst *ressource.ValidateApplicationRessourceAccessRqst) (*ressource.ValidateApplicationRessourceAccessRsp, error) {

	path := rqst.GetPath()

	hasApplicationPermission, count := self.hasPermission(rqst.Name, path, rqst.Permission)
	if hasApplicationPermission {
		return &ressource.ValidateApplicationRessourceAccessRsp{
			Result: true,
		}, nil
	}

	// If theres permission definied and we are here it's means the user dosent
	// have write to execute the action on the ressource.
	if count > 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Permission Denied for "+rqst.Name)))

	}

	return &ressource.ValidateApplicationRessourceAccessRsp{
		Result: true,
	}, nil
}

//* Validate if a peer can access a given ressource. *
func (self *Globule) ValidatePeerRessourceAccess(ctx context.Context, rqst *ressource.ValidatePeerRessourceAccessRqst) (*ressource.ValidatePeerRessourceAccessRsp, error) {
	return nil, nil
}

//* Validate if a peer can access a given method. *
func (self *Globule) ValidatePeerAccess(ctx context.Context, rqst *ressource.ValidatePeerAccessRqst) (*ressource.ValidatePeerAccessRsp, error) {
	return nil, nil
}

//* Validate if user can access a given method. *
func (self *Globule) ValidateUserAccess(ctx context.Context, rqst *ressource.ValidateUserAccessRqst) (*ressource.ValidateUserAccessRsp, error) {

	// first of all I will validate the token.
	clientID, expiredAt, err := Interceptors.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if expiredAt < time.Now().Unix() {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The token is expired!")))
	}

	// Here I will test if the user can run that function or not...
	err = self.validateUserAccess(clientID, rqst.Method)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.ValidateUserAccessRsp{
		Result: true,
	}, nil
}

//* Validate if application can access a given method. *
func (self *Globule) ValidateApplicationAccess(ctx context.Context, rqst *ressource.ValidateApplicationAccessRqst) (*ressource.ValidateApplicationAccessRsp, error) {
	err := self.validateApplicationAccess(rqst.GetName(), rqst.Method)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.ValidateApplicationAccessRsp{
		Result: true,
	}, nil
}

//* Retrun a json string with all file info *
func (self *Globule) GetAllFilesInfo(ctx context.Context, rqst *ressource.GetAllFilesInfoRqst) (*ressource.GetAllFilesInfoRsp, error) {
	// That map will contain the list of all directories.
	dirs := make(map[string]map[string]interface{})

	err := filepath.Walk(self.webRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				dir := make(map[string]interface{})
				dir["name"] = info.Name()
				dir["size"] = info.Size()
				dir["last_modified"] = info.ModTime().Unix()
				dir["path"] = strings.ReplaceAll(strings.Replace(path, self.path, "", -1), "\\", "/")
				dir["files"] = make([]interface{}, 0)
				dirs[dir["path"].(string)] = dir
				parent := dirs[dir["path"].(string)[0:strings.LastIndex(dir["path"].(string), "/")]]
				if parent != nil {
					parent["files"] = append(parent["files"].([]interface{}), dir)
				}
			} else {
				file := make(map[string]interface{})
				file["name"] = info.Name()
				file["size"] = info.Size()
				file["last_modified"] = info.ModTime().Unix()
				file["path"] = strings.ReplaceAll(strings.Replace(path, self.path, "", -1), "\\", "/")
				dir := dirs[file["path"].(string)[0:strings.LastIndex(file["path"].(string), "/")]]
				dir["files"] = append(dir["files"].([]interface{}), file)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	jsonStr, err := json.Marshal(dirs[strings.ReplaceAll(strings.Replace(self.webRoot, self.path, "", -1), "\\", "/")])
	if err != nil {
		return nil, err
	}
	return &ressource.GetAllFilesInfoRsp{Result: string(jsonStr)}, nil
}

func (self *Globule) GetAllApplicationsInfo(ctx context.Context, rqst *ressource.GetAllApplicationsInfoRqst) (*ressource.GetAllApplicationsInfoRsp, error) {

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// So here I will get the list of retreived permission.
	jsonStr, err := p.Find("local_ressource", "local_ressource", "Applications", `{}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.GetAllApplicationsInfoRsp{
		Result: jsonStr,
	}, nil

}

//////////////////////////////// Services management  //////////////////////////

/**
 * Get access to the event services.
 */
func (self *Globule) getEventHub() (*event_client.Event_Client, error) {
	var err error
	if self.event_client_ == nil {
		self.event_client_, err = event_client.NewEvent_Client(self.getDomain(), "event_server")
	}
	return self.event_client_, err
}

// Discovery
func (self *Globule) FindServices(ctx context.Context, rqst *services.FindServicesDescriptorRequest) (*services.FindServicesDescriptorResponse, error) {
	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	kewordsStr, err := Utility.ToJson(rqst.Keywords)
	if err != nil {
		return nil, err
	}

	// Test...
	query := `{"keywords": { "$all" : ` + kewordsStr + `}}`

	data, err := p.Find("local_ressource", "local_ressource", "Services", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*services.ServiceDescriptor, 0)
	err = json.Unmarshal([]byte(data), &descriptors)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the list of Service Descriptor.
	return &services.FindServicesDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (self *Globule) GetServiceDescriptor(ctx context.Context, rqst *services.GetServiceDescriptorRequest) (*services.GetServiceDescriptorResponse, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	query := `{"id":"` + rqst.ServiceId + `", "publisherId":"` + rqst.PublisherId + `"}`

	data, err := p.Find("local_ressource", "local_ressource", "Services", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*services.ServiceDescriptor, 0)
	err = json.Unmarshal([]byte(data), &descriptors)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	sort.Slice(descriptors[:], func(i, j int) bool {
		return descriptors[i].Version > descriptors[j].Version
	})

	// Return the list of Service Descriptor.
	return &services.GetServiceDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Return the list of all services *
func (self *Globule) GetServicesDescriptor(ctx context.Context, rqst *services.GetServicesDescriptorRequest) (*services.GetServicesDescriptorResponse, error) {
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	data, err := p.Find("local_ressource", "local_ressource", "Services", `{}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	descriptors := make([]*services.ServiceDescriptor, 0)
	err = json.Unmarshal([]byte(data), &descriptors)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Return the list of Service Descriptor.
	return &services.GetServicesDescriptorResponse{
		Results: descriptors,
	}, nil
}

//* Publish a service to service discovery *
func (self *Globule) PublishServiceDescriptor(ctx context.Context, rqst *services.PublishServiceDescriptorRequest) (*services.PublishServiceDescriptorResponse, error) {

	// Here I will save the descriptor inside the storage...
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append the self domain to the list of discoveries where the services can be found.
	if !Utility.Contains(rqst.Descriptor_.Discoveries, self.getDomain()) {
		rqst.Descriptor_.Discoveries = append(rqst.Descriptor_.Discoveries, self.getDomain())
	}

	// Here I will test if the services already exist...
	_, err = p.FindOne("local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherId":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, "")
	if err == nil {
		// Update existing descriptor.

		// The list of discoveries...
		discoveries, err := Utility.ToJson(rqst.Descriptor_.Discoveries)
		if err == nil {
			values := `{"$set":{"discoveries":` + discoveries + `}}`
			err = p.Update("local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherId":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}

		// The list of repositories
		repositories, err := Utility.ToJson(rqst.Descriptor_.Repositories)
		if err == nil {
			values := `{"$set":{"repositories":` + repositories + `}}`
			err = p.Update("local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherId":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, values, "")
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
		}

	} else {
		data, err := json.Marshal(rqst.Descriptor_)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
		// The key will be the descriptor string itself.
		_, err = p.InsertOne("local_ressource", "local_ressource", "Services", string(data), "")
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

	}

	return &services.PublishServiceDescriptorResponse{
		Result: true,
	}, nil
}

// Repository
/** Download a service from a service directory **/
func (self *Globule) DownloadBundle(rqst *services.DownloadBundleRequest, stream services.ServiceRepository_DownloadBundleServer) error {
	bundle := new(services.ServiceBundle)
	bundle.Plaform = rqst.Plaform
	bundle.Descriptor_ = rqst.Descriptor_

	// Generate the bundle id....
	var id string
	id = bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Id + "%" + bundle.Descriptor_.Version
	if bundle.Plaform == services.Platform_LINUX32 {
		id += "%LINUX32"
	} else if bundle.Plaform == services.Platform_LINUX64 {
		id += "%LINUX64"
	} else if bundle.Plaform == services.Platform_WIN32 {
		id += "%WIN32"
	} else if bundle.Plaform == services.Platform_WIN64 {
		id += "%WIN64"
	}

	path := self.data + string(os.PathSeparator) + "service-repository"

	var err error
	// the file must be a zipped archive that contain a .proto, .config and executable.
	bundle.Binairies, err = ioutil.ReadFile(path + string(os.PathSeparator) + id + ".tar.gz")
	if err != nil {
		return err
	}

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	jsonStr, err := p.FindOne("local_ressource", "local_ressource", "ServiceBundle", `{"_id":"`+id+`"}`, "")
	if err != nil {
		return err
	}

	// init the map with json values.
	checksum := make(map[string]interface{}, 0)
	json.Unmarshal([]byte(jsonStr), &checksum)

	// Test if the values change over time.
	if Utility.CreateDataChecksum(bundle.Binairies) != checksum["checksum"].(string) {
		return errors.New("The bundle data cheksum is not valid!")
	}

	const BufferSize = 1024 * 5 // the chunck size.
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer) // Will write to network.
	err = enc.Encode(bundle)
	if err != nil {
		return err
	}

	for {
		var data [BufferSize]byte
		bytesread, err := buffer.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &services.DownloadBundleResponse{
				Data: data[0:bytesread],
			}
			// send the data to the server.
			err = stream.Send(rqst)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

/** Upload a service to a service directory **/
func (self *Globule) UploadBundle(stream services.ServiceRepository_UploadBundleServer) error {

	// The bundle will cantain the necessary information to install the service.
	var buffer bytes.Buffer
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			stream.SendAndClose(&services.UploadBundleResponse{
				Result: true,
			})
			break
		} else if err != nil {
			return err
		} else {
			buffer.Write(msg.Data)
		}
	}

	// The buffer that contain the
	dec := gob.NewDecoder(&buffer)
	bundle := new(services.ServiceBundle)
	err := dec.Decode(bundle)
	if err != nil {
		return err
	}

	// Generate the bundle id....
	id := bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Id + "%" + bundle.Descriptor_.Version
	if bundle.Plaform == services.Platform_LINUX32 {
		id += "%LINUX32"
	} else if bundle.Plaform == services.Platform_LINUX64 {
		id += "%LINUX64"
	} else if bundle.Plaform == services.Platform_WIN32 {
		id += "%WIN32"
	} else if bundle.Plaform == services.Platform_WIN64 {
		id += "%WIN64"
	}

	repositoryId := self.getDomain()
	// Now I will append the address of the repository into the service descriptor.
	if !Utility.Contains(bundle.Descriptor_.Repositories, repositoryId) {
		bundle.Descriptor_.Repositories = append(bundle.Descriptor_.Repositories, repositoryId)
		// Publish change into discoveries...
		for i := 0; i < len(bundle.Descriptor_.Discoveries); i++ {
			discoveryId := bundle.Descriptor_.Discoveries[i]
			discoveryService, err := services.NewServicesDiscovery_Client(discoveryId, "services_discovery")
			if err != nil {
				return err
			}
			discoveryService.PublishServiceDescriptor(bundle.Descriptor_)
		}
	}

	path := self.data + string(os.PathSeparator) + "service-repository"
	Utility.CreateDirIfNotExist(path)

	// the file must be a zipped archive that contain a .proto, .config and executable.
	err = ioutil.WriteFile(path+string(os.PathSeparator)+id+".tar.gz", bundle.Binairies, 777)
	if err != nil {
		return err
	}

	checksum := Utility.CreateDataChecksum(bundle.Binairies)
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return err
	}

	_, err = p.InsertOne("local_ressource", "local_ressource", "ServiceBundle", `{"_id":"`+id+`","checksum":"`+checksum+`"}`, "")

	return err
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() error {
	// append itself to service discoveries...
	if !Utility.Contains(self.Discoveries, self.getDomain()) {
		self.Discoveries = append(self.Discoveries, self.getDomain())
	}

	subscribers := make(map[string]map[string][]string, 0)

	// Connect to service update events...
	for i := 0; i < len(self.Discoveries); i++ {
		eventHub, err := event_client.NewEvent_Client(self.Discoveries[i], "event_server")
		if err == nil {
			subscribers[self.Discoveries[i]] = make(map[string][]string)
			for _, s := range self.Services {
				if s.(map[string]interface{})["PublisherId"] != nil {
					id := s.(map[string]interface{})["PublisherId"].(string) + ":" + s.(map[string]interface{})["Name"].(string) + ":SERVICE_PUBLISH_EVENT"
					if subscribers[self.Discoveries[i]][id] == nil {
						subscribers[self.Discoveries[i]][id] = make([]string, 0)
					}
					// each channel has it event...
					uuid := Utility.RandomUUID()
					fct := func(evt *eventpb.Event) {
						descriptor := new(services.ServiceDescriptor)
						json.Unmarshal(evt.GetData(), descriptor)
						// here I will update the service if it's version is lower
						for _, s := range self.Services {
							service := s.(map[string]interface{})
							if service["PublisherId"] != nil {
								if service["Name"].(string) == descriptor.GetId() && service["PublisherId"].(string) == descriptor.GetPublisherId() {
									if service["KeepUpToDate"] != nil {
										if service["KeepUpToDate"].(bool) {
											// Test if update is needed...
											if Utility.ToInt(strings.Split(service["Version"].(string), ".")[0]) > Utility.ToInt(strings.Split(descriptor.Version, ".")[0]) {
												if Utility.ToInt(strings.Split(service["Version"].(string), ".")[1]) > Utility.ToInt(strings.Split(descriptor.Version, ".")[1]) {
													if Utility.ToInt(strings.Split(service["Version"].(string), ".")[2]) > Utility.ToInt(strings.Split(descriptor.Version, ".")[2]) {
														self.stopService(service["Id"].(string))
														delete(self.Services, service["Id"].(string))
														err := self.installService(descriptor)
														if err != nil {
															log.Println("fail to install service ", err)
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

					// So here I will subscribe to service update event.
					eventHub.Subscribe(id, uuid, fct)
					subscribers[self.Discoveries[i]][id] = append(subscribers[self.Discoveries[i]][id], uuid)
				}
			}
		}
		// keep on memorie...
		self.discorveriesEventHub[self.Discoveries[i]] = eventHub
	}

	// Catch the Ctrl-C and SIGTERM from kill command
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		signalType := <-ch
		signal.Stop(ch)
		log.Println("Exit command received. Exiting...")

		// this is a good place to flush everything to disk
		// before terminating.
		log.Println("Signal type : ", signalType)

		// Here the server stop running,
		// so I will close the services.
		log.Println("Clean ressources.")

		for key, value := range self.Services {
			log.Println("Stop service ", key)
			if value.(map[string]interface{})["Process"] != nil {
				p := value.(map[string]interface{})["Process"]
				if reflect.TypeOf(p).String() == "*exec.Cmd" {
					if p.(*exec.Cmd).Process != nil {
						log.Println("kill service process ", p.(*exec.Cmd).Process.Pid)
						p.(*exec.Cmd).Process.Kill()
					}
				}
			}

			if value.(map[string]interface{})["ProxyProcess"] != nil {
				p := value.(map[string]interface{})["ProxyProcess"]
				if reflect.TypeOf(p).String() == "*exec.Cmd" {
					if p.(*exec.Cmd).Process != nil {
						log.Println("kill proxy process ", p.(*exec.Cmd).Process.Pid)
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

	// Start the admin service to give access to server functionality from
	// client side.
	admin_server, err := self.startInternalService("admin", self.AdminPort, self.AdminProxy, self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor) // must be accessible to all clients...
	if err == nil {
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.AdminPort))
		if err != nil {
			log.Fatalf("could not start admin service %s: %s", self.getDomain(), err)
		}

		admin.RegisterAdminServiceServer(admin_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				log.Println("Admin service is up and running for domain ", self.getDomain())
				// no web-rpc server.
				if err := admin_server.Serve(lis); err != nil {
					f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer f.Close()
				}
				log.Println("Adim grpc service is closed")
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			Utility.KillProcessByName("mongod")
			Utility.KillProcessByName("prometheus")
		}()
	} else {
		log.Println(err)
	}

	ressource_server, err := self.startInternalService("ressource", self.RessourcePort, self.RessourceProxy, self.Protocol == "https", self.unaryRessourceInterceptor, self.streamRessourceInterceptor)
	if err == nil {

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.RessourcePort))
		if err != nil {
			log.Fatalf("could not start ressource service %s: %s", self.getDomain(), err)
		}

		ressource.RegisterRessourceServiceServer(ressource_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				log.Println("Ressource service is up and running for domain ", self.getDomain())
				// no web-rpc server.
				if err = ressource_server.Serve(lis); err != nil {
					f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer f.Close()
				}

				log.Println("ressource grpc service is closed")
			}()

			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch

		}()

		// In order to be able to give permission to a server
		// I must register it to the globule associated
		// with the base domain.

		// set the ip into the DNS servers.
		ticker_ := time.NewTicker(5 * time.Second)
		go func() {
			ip := Utility.MyIP()
			self.registerIpToDns()
			for {
				select {
				case <-ticker_.C:
					if ip != Utility.MyIP() {
						self.registerIpToDns()
					}
				}
			}
		}()

	} else {
		log.Println(err)
	}

	// The service discovery.
	services_discovery_server, err := self.startInternalService("services_discovery", self.ServicesDiscoveryPort, self.ServicesDiscoveryProxy, self.Protocol == "https", Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesDiscoveryPort))
		if err != nil {
			log.Fatalf("could not start services discovery service %s: %s", self.getDomain(), err)
		}

		services.RegisterServiceDiscoveryServer(services_discovery_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				log.Println("Discovery service is up and running for domain ", self.getDomain())

				// no web-rpc server.
				if err := services_discovery_server.Serve(lis); err != nil {
					f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer f.Close()
				}
				log.Println("services discovery grpc service is closed")
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
		}()
	} else {
		log.Println(err)
	}

	// The service repository
	services_repository_server, err := self.startInternalService("services_repository", self.ServicesRepositoryPort, self.ServicesRepositoryProxy,
		self.Protocol == "https",
		Interceptors.ServerUnaryInterceptor,
		Interceptors.ServerStreamInterceptor)

	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesRepositoryPort))
		if err != nil {
			log.Fatalf("could not start services repository service %s: %s", self.getDomain(), err)
		}

		services.RegisterServiceRepositoryServer(services_repository_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				log.Println("Repository service is up and running for domain ", self.getDomain())

				// no web-rpc server.
				if err := services_repository_server.Serve(lis); err != nil {
					f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer f.Close()
				}
				log.Println("services repository grpc service is closed")
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
		}()
	} else {
		log.Println(err)
	}

	// The Certificate Authority
	certificate_authority_server, err := self.startInternalService("certificate_authority", self.CertificateAuthorityPort, self.CertificateAuthorityProxy, false, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.CertificateAuthorityPort))
		if err != nil {
			log.Fatalf("could not certificate authority signing  service %s: %s", self.getDomain(), err)
		}

		ca.RegisterCertificateAuthorityServer(certificate_authority_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			go func() {
				log.Println("Certificate Authority service is up and running for domain ", self.getDomain())

				// no web-rpc server.
				if err := certificate_authority_server.Serve(lis); err != nil {
					f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer f.Close()
				}
				log.Println("services repository grpc service is closed")
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
		}()
	} else {
		log.Println(err)
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

			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/grpc_tls/server.csr run
			// I need to remove the gRPC certificate and recreate it.
			Utility.RemoveDirContents(self.creds)

			security.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay, self.getDomain(), self.certs)

			time.Sleep(15 * time.Second)

			self.initServices() // must restart the services with new certificates.
			err := self.ObtainCertificateForCsr()
			if err != nil {
				log.Println(err)
			}
		}

		log.Println("start https server")
		// get the value from the configuration files.
		err := server.ListenAndServeTLS(self.certs+string(os.PathSeparator)+self.Certificate, self.creds+string(os.PathSeparator)+"server.pem")
		if err != nil {
			log.Println(err)
		}

	} else {
		log.Println("start http server")
		// local - non secure connection.
		http.ListenAndServe(":"+strconv.Itoa(self.PortHttp), nil)
	}

	if err != nil {
		log.Println("ListenAndServe: " + err.Error())
	}
	return nil
}

/////////////////////////// Security stuff //////////////////////////////////

/////////////////////////////////////////////////////////////////////////////
// Certificate Authority Service
/////////////////////////////////////////////////////////////////////////////

func (self *Globule) signCertificate(client_csr string) (string, error) {

	// first of all I will save the incomming file into a temporary file...
	client_csr_path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()
	err := ioutil.WriteFile(client_csr_path, []byte(client_csr), 0644)
	if err != nil {
		return "", err

	}

	client_crt_path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+self.CertPassword)
	args = append(args, "-days")
	args = append(args, Utility.ToString(self.CertExpirationDelay))
	args = append(args, "-in")
	args = append(args, client_csr_path)
	args = append(args, "-CA")
	args = append(args, self.creds+string(os.PathSeparator)+"ca.crt") // use certificate
	args = append(args, "-CAkey")
	args = append(args, self.creds+string(os.PathSeparator)+"ca.key") // and private key to sign the incommin csr
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, client_crt_path)

	err = exec.Command(cmd, args...).Run()
	if err != nil {

		return "", err
	}

	// I will read back the crt file.
	client_crt, err := ioutil.ReadFile(client_crt_path)

	// remove the tow temporary files.
	defer os.Remove(client_crt_path)
	defer os.Remove(client_csr_path)

	if err != nil {
		return "", err
	}

	return string(client_crt), nil

}

// Signed certificate request (CSR)
func (self *Globule) SignCertificate(ctx context.Context, rqst *ca.SignCertificateRequest) (*ca.SignCertificateResponse, error) {

	client_crt, err := self.signCertificate(rqst.Csr)

	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	}

	return &ca.SignCertificateResponse{
		Crt: client_crt,
	}, nil
}

// Return the Authority Trust Certificate. (ca.crt)
func (self *Globule) GetCaCertificate(ctx context.Context, rqst *ca.GetCaCertificateRequest) (*ca.GetCaCertificateResponse, error) {

	ca_crt, err := ioutil.ReadFile(self.creds + string(os.PathSeparator) + "ca.crt")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ca.GetCaCertificateResponse{
		Ca: string(ca_crt),
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Certificate from let's encrytp via lego.
////////////////////////////////////////////////////////////////////////////////

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
		log.Println(err)
	}
	return privateKey
}

///////// End of Implement the User Interface. ////////////

/**
 * That function work correctly, but the DNS fail time to time to give the
 * IP address that result in a fail request... The DNS must be fix!
 */
func (self *Globule) ObtainCertificateForCsr() error {

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
