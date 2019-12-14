package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
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

	"github.com/davecourtois/Globular/api"

	// Admin service
	"github.com/davecourtois/Globular/admin"
	// Ressource service
	"github.com/davecourtois/Globular/ressource"
	// Services management service
	"github.com/davecourtois/Globular/services"
	// Certificate authority
	"github.com/davecourtois/Globular/ca"

	// Interceptor for authentication, event, log...
	Interceptors "github.com/davecourtois/Globular/Interceptors/Authenticate"
	Interceptors_ "github.com/davecourtois/Globular/Interceptors/server"

	// Client services.
	"context"

	"crypto"

	"github.com/davecourtois/Globular/catalog/catalog_client"
	"github.com/davecourtois/Globular/dns/dns_client"
	"github.com/davecourtois/Globular/echo/echo_client"
	"github.com/davecourtois/Globular/event/event_client"
	"github.com/davecourtois/Globular/file/file_client"
	"github.com/davecourtois/Globular/ldap/ldap_client"
	"github.com/davecourtois/Globular/monitoring/monitoring_client"
	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/plc/plc_client"
	"github.com/davecourtois/Globular/smtp/smtp_client"
	"github.com/davecourtois/Globular/spc/spc_client"
	"github.com/davecourtois/Globular/sql/sql_client"
	"github.com/davecourtois/Globular/storage/storage_client"
	"github.com/davecourtois/Utility"
	"github.com/emicklei/proto"
	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/challenge/http01"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/registration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
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
	ExternalApplications       map[string]ExternalApplication
	Domain                     string   // The domain (subdomain) name of your application
	DNS                        []string // Contain a list of domain name server where that computer use as sub-domain
	ReadTimeout                time.Duration
	WriteTimeout               time.Duration
	IdleTimeout                time.Duration
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

	// Directories.
	path    string // The path of the exec...
	webRoot string // The root of the http file server.
	data    string // the data directory
	creds   string // gRpc certificate
	certs   string // https certificates
	config  string // configuration directory

	// The map of client...
	clients map[string]api.Client

	// Create the JWT key used to create the signature
	jwtKey       []byte
	RootPassword string
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
	g.IdleTimeout = 120
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.ReadTimeout = 5
	g.WriteTimeout = 5
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"
	g.DNS = make([]string, 0)

	// Set the list of discorvery service avalaible...
	g.Discoveries = make([]string, 0)
	g.discorveriesEventHub = make(map[string]*event_client.Event_Client, 0)

	// Set the share service info...
	g.Services = make(map[string]interface{}, 0)

	// Set external map services.
	g.ExternalApplications = make(map[string]ExternalApplication, 0)

	// Set the map of client.
	g.clients = make(map[string]api.Client, 0)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	g.path = dir // keep the installation path.

	// if globular is found.
	g.webRoot = dir + string(os.PathSeparator) + "webroot" // The default directory to server.

	// keep the root in global variable for the file handler.
	root = g.webRoot
	Utility.CreateDirIfNotExist(g.webRoot) // Create the directory if it not exist.

	if !Utility.Exists(g.webRoot + string(os.PathSeparator) + "index.html") {

		// in that case I will create a new index.html file.
		ioutil.WriteFile(g.webRoot+string(os.PathSeparator)+"index.html", []byte(
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
	g.data = dir + string(os.PathSeparator) + "data"
	Utility.CreateDirIfNotExist(g.data)

	// Configuration directory
	g.config = dir + string(os.PathSeparator) + "config"
	Utility.CreateDirIfNotExist(g.config)

	// Create the creds directory if it not already exist.
	g.creds = g.config + string(os.PathSeparator) + "grpc_tls"
	Utility.CreateDirIfNotExist(g.creds)

	// https certificates.
	g.certs = g.config + string(os.PathSeparator) + "http_tls"
	Utility.CreateDirIfNotExist(g.certs)

	// Initialyse globular from it configuration file.
	file, err := ioutil.ReadFile(g.config + string(os.PathSeparator) + "config.json")

	// Init the service with the default port address
	if err == nil {
		json.Unmarshal(file, &g)
	}

	// Keep in global var to by http handlers.
	globule = g

	// The configuration handler.
	http.HandleFunc("/client_config", getClientConfigHanldler)
	http.HandleFunc("/config", getConfigHanldler)

	// Configuration must be reachable before services initialysation
	go func() {
		http.ListenAndServe(":10000", nil)
	}()

	return g
}

/**
 * Serve
 */
func (self *Globule) Serve() {

	// Here I will kill proxies if there are running.
	Utility.KillProcessByName("grpcwebproxy")

	// Here it suppose to be only one server instance per computer.
	self.jwtKey = []byte(Utility.RandomUUID())
	err := ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_key", []byte(self.jwtKey), 0644)
	if err != nil {
		log.Panicln(err)
	}

	// The token that identify the server with other services
	token, _ := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, "sa")
	err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"localhost_token", []byte(token), 0644)
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
				err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"localhost_token", []byte(token), 0644)
				if err != nil {
					log.Panicln(err)
				}
			}
		}
	}()

	self.initClients()

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

	// set the services.
	self.initServices()

	// Here I will save the server attribute
	self.saveConfig()

	// Here i will connect the service listener.
	time.Sleep(5 * time.Second) // wait for services to start...

	// lisen
	self.Listen()
}

/**
 * Return the server configuration
 */
func getClientConfigHanldler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address") // parameter address
	name := r.URL.Query().Get("name")       // parameter name
	config, err := getClientConfig(address, name)
	if err != nil {
		http.Error(w, "Client configuration "+name+" not found!", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
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
 * Set the ip for a given sub-domain compose of Name + DNS domain.
 */
func (self *Globule) registerIpToDns() {
	if self.DNS != nil {
		if len(self.DNS) > 0 {
			for i := 0; i < len(self.DNS); i++ {
				log.Println("register domain to dns:", self.DNS[i])
				client := dns_client.NewDns_Client(self.DNS[i], "dns_server")

				domain, err := client.SetA(strings.ToLower(self.Name), Utility.MyIP(), 60)

				if err != nil {
					log.Println(err)
				} else {
					log.Println("---> register ip ", Utility.MyIP(), "with", domain)
				}

				// TODO also register the ipv6 here...
				client.Close()
			}
		}
	}
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

	proxyBackendAddress := self.Domain + ":" + strconv.Itoa(port)
	proxyAllowAllOrgins := "true"
	proxyArgs := make([]string, 0)

	// Use in a local network or in test.
	proxyArgs = append(proxyArgs, "--backend_addr="+proxyBackendAddress)
	proxyArgs = append(proxyArgs, "--allow_all_origins="+proxyAllowAllOrgins)

	if srv.(map[string]interface{})["TLS"].(bool) == true {
		log.Println("---> start secure service: ", srv.(map[string]interface{})["Name"])
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
		log.Println("---> start non secure service: ", srv.(map[string]interface{})["Name"])

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
	// In case the service must not be kept alive.
	if !s["KeepAlive"].(bool) {
		return
	}

	s["Process"].(*exec.Cmd).Wait()
	time.Sleep(time.Second * 10)
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
	// test if the service is distant or not.
	if !Utility.IsLocal(s["Domain"].(string)) {
		return -1, -1, errors.New("Can not start a distant service localy!")
	}

	// set the domain of the service.
	s["Domain"] = self.Domain

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

	if s["Protocol"].(string) == "grpc" {

		// Stop the previous client if there one.
		if self.clients[s["Name"].(string)+"_service"] != nil {
			self.clients[s["Name"].(string)+"_service"].Close()
		}

		if s["TLS"].(bool) {
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

		// Kill previous instance of the program...
		Utility.KillProcessByName(s["Name"].(string))

		// Start the service process.
		log.Println("try to start process ", s["Name"].(string))

		servicePath := self.path + s["servicePath"].(string)
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

		err = s["Process"].(*exec.Cmd).Start()
		go func() {
			self.keepServiceAlive(s)
		}()

		if err != nil {
			log.Panicln("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
			return -1, -1, err
		}

		// Save configuration stuff.
		self.Services[s["Id"].(string)] = s

		// Start the proxy.
		err = self.startProxy(s["Id"].(string), int(s["Port"].(float64)), int(s["Proxy"].(float64)))
		if err != nil {
			return -1, -1, err
		}

		// get back the service info with the proxy process in it
		s = self.Services[s["Id"].(string)].(map[string]interface{})

		// Init it configuration.
		self.initClient(s["Name"].(string), s["Name"].(string))

		// save it to the config.
		self.saveConfig()

	} else if s["Protocol"].(string) == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(s["Name"].(string), "Globular") {
			// Kill previous instance of the program.
			Utility.KillProcessByName(s["Name"].(string))
			log.Println("try to start process ", s["Name"].(string))
			s["Process"] = exec.Command(s["servicePath"].(string), Utility.ToString(s["Port"]))

			err = s["Process"].(*exec.Cmd).Start()
			if err != nil {
				log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
				return -1, -1, err
			}
			self.Services[s["Id"].(string)] = s

			return s["Process"].(*exec.Cmd).Process.Pid, -1, nil
		}
	}

	if s["Process"].(*exec.Cmd).Process == nil {
		return -1, -1, errors.New("Fail to start process " + s["Name"].(string))
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
		if s["TLS"].(bool) {
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

	// If the protocol is https I will generate the TLS certificate.
	self.GenerateServicesCertificates("1111", self.CertExpirationDelay)

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
				log.Println("--> config found at ", path)
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
							s["Id"] = s["Name"].(string)

							s["servicePath"] = strings.Replace(strings.Replace(path_+string(os.PathSeparator)+s["Name"].(string), self.path, "", -1), "\\", "/", -1)
							s["configPath"] = strings.Replace(strings.Replace(path, self.path, "", -1), "\\", "/", -1)
							s["schemaPath"] = strings.Replace(strings.Replace(path_+string(os.PathSeparator)+"schema.json", self.path, "", -1), "\\", "/", -1)

							self.Services[s["Name"].(string)] = s

						} else {
							log.Println("--> no protocol found")
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
					// So here I will register the method into the backend.
					self.methods = append(self.methods, path)
				}
			}
		}
		return nil
	})

	// Init services.
	for id, s := range self.Services {
		// Remove existing process information.
		delete(s.(map[string]interface{}), "Process")
		delete(s.(map[string]interface{}), "ProxyProcess")
		log.Println("--> init service ", id)
		err := self.initService(s.(map[string]interface{}))
		if err != nil {
			log.Println(err)
		}
	}

	// if a dns service exist I will set the name of that globule on the server.
	if self.clients["dns_service"] != nil {
		// Set the server
		domain, err := self.clients["dns_service"].(*dns_client.DNS_Client).SetA(strings.ToLower(self.Name), Utility.MyIP(), 60)
		if err == nil {
			log.Println("---> set domain ", domain, "with ip", Utility.MyIP())
		} else {
			log.Println("---> fail to register ip with dns", err)
		}
	}

	// set the ip into the DNS servers.
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		ip := Utility.MyIP()
		self.registerIpToDns()
		for {
			select {
			case <-ticker.C:
				/** If the ip change I will update the domain. **/
				if ip != Utility.MyIP() {
					self.registerIpToDns()
				}
			}
		}
	}()
}

// Method must be register in order to be assign to role.
func (self *Globule) registerMethods() error {
	// Here I will create the sa role if it dosen't exist.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		log.Println("---> fail to get local_ressource connection", err)
		return err
	}

	// Here I will persit the sa role if it dosent already exist.
	count, err := p.Count("local_ressource", "local_ressource", "Roles", `{ "_id":"sa"}`, "")
	admin := make(map[string]interface{})
	if err != nil {
		log.Println("---> fail to count local ressource.", err)
		return err
	} else if count == 0 {
		log.Println("need to create admin roles...")
		admin["_id"] = "sa"
		admin["actions"] = self.methods
		jsonStr, _ := Utility.ToJson(admin)
		id, err := p.InsertOne("local_ressource", "local_ressource", "Roles", jsonStr, "")
		if err != nil {
			return err
		}
		log.Println("role with id", id, "was created!")
	} else {
		admin["_id"] = "sa"
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
		guest["actions"] = []string{"/admin.AdminService/GetConfig", "/ressource.RessourceService/RegisterAccount", "/ressource.RessourceService/Authenticate", "/event.EventService/Subscribe", "/event.EventService/UnSubscribe", "/event.EventService/Publish", "/services.ServiceDiscovery/FindServices",
			"/services.ServiceDiscovery/GetServiceDescriptor", "/services.ServiceDiscovery/GetServicesDescriptor", "/services.ServiceRepository/downloadBundle"}
		jsonStr, _ := Utility.ToJson(guest)
		_, err := p.InsertOne("local_ressource", "local_ressource", "Roles", jsonStr, "")
		if err != nil {
			return err
		}
		log.Println("role guest was created!")
	}

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
 * Init client side connection to service.
 */
func (self *Globule) initClient(id string, name string) {
	if self.Services[id] == nil {
		return
	}

	serviceName := name
	name = strings.Split(name, "_")[0]

	fct := "New" + strings.ToUpper(name[0:1]) + name[1:] + "_Client"
	results, err := Utility.CallFunction(fct, self.Services[id].(map[string]interface{})["Domain"].(string), serviceName)
	if err == nil {
		if self.clients[name+"_service"] != nil {
			self.clients[name+"_service"].Close()
		}
		self.clients[name+"_service"] = results[0].Interface().(api.Client)
		log.Println("--> client ", name+"_service", "is now initialysed!")
	} else {
		log.Panicln(err)
	}
}

/**
 * Init the service client.
 * Keep the service constructor for further call. This is not fully generic,
 * maybe reflection will be use in futur implementation.
 */
func (self *Globule) initClients() {

	// Register service constructor function here.
	// The name of the contructor must follow the same pattern
	Utility.RegisterFunction("NewPersistence_Client", persistence_client.NewPersistence_Client)
	Utility.RegisterFunction("NewEcho_Client", echo_client.NewEcho_Client)
	Utility.RegisterFunction("NewSql_Client", sql_client.NewSql_Client)
	Utility.RegisterFunction("NewFile_Client", file_client.NewFile_Client)
	Utility.RegisterFunction("NewSmtp_Client", smtp_client.NewSmtp_Client)
	Utility.RegisterFunction("NewLdap_Client", ldap_client.NewLdap_Client)
	Utility.RegisterFunction("NewStorage_Client", storage_client.NewStorage_Client)
	Utility.RegisterFunction("NewEvent_Client", event_client.NewEvent_Client)
	Utility.RegisterFunction("NewCatalog_Client", catalog_client.NewCatalog_Client)
	Utility.RegisterFunction("NewMonitoring_Client", monitoring_client.NewMonitoring_Client)
	Utility.RegisterFunction("NewDns_Client", dns_client.NewDns_Client)

	// That service is program in c++
	Utility.RegisterFunction("NewSpc_Client", spc_client.NewSpc_Client)
	Utility.RegisterFunction("NewPlc_Client", plc_client.NewPlc_Client)
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
	config["ServicesRepositoryProxy"] = self.ServicesRepositoryProxy
	config["Discoveries"] = self.Discoveries
	config["DNS"] = self.DNS
	config["Protocol"] = self.Protocol
	config["Domain"] = self.Domain
	config["ReadTimeout"] = self.ReadTimeout
	config["WriteTimeout"] = self.WriteTimeout
	config["IdleTimeout"] = self.IdleTimeout
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

	return nil
}

/**
 * Start the monitoring service with prometheus.
 */
func (self *Globule) startMonitoring() error {
	Utility.KillProcessByName("prometheus")
	var m *monitoring_client.Monitoring_Client
	if self.clients["monitoring_service"] == nil {
		log.Println("---> no monitoring service is configure.")
		return errors.New("No monitoring service are available to store monitoring information.")
	}

	// Cast-it to the persistence client.
	m = self.clients["monitoring_service"].(*monitoring_client.Monitoring_Client)

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
  
  - job_name: 'node_exporter_metrics'
    scrape_interval: 5s
    static_configs:
    - targets: ['localhost:9100']
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

	log.Println("---> start prometheus alert manager")
	alertmanager := exec.Command("alertmanager", "--config.file", self.config+string(os.PathSeparator)+"alertmanager.yml")
	err := alertmanager.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	log.Println("---> start prometheus node exporter on port 9100")
	node_exporter := exec.Command("node_exporter")
	err = node_exporter.Start()
	if err != nil {
		log.Println("fail to start prometheus node exporter", err)
		// do not return here in that case simply continue without node exporter metrics.
	}

	log.Println("---> start prometheus on port 9090")
	prometheus := exec.Command("prometheus", "--web.listen-address", "0.0.0.0:9090", "--config.file", self.config+string(os.PathSeparator)+"prometheus.yml", "--storage.tsdb.path", dataPath)
	err = prometheus.Start()
	if err != nil {
		log.Println("fail to start monitoring with prometheus", err)
		return err
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
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		log.Println("---> no persistence service is configure.")
		return nil, errors.New("No persistence service are available to store ressource information.")
	}

	// Cast-it to the persistence client.
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)

	// Connect to the database here.
	err := p.CreateConnection("local_ressource", "local_ressource", "0.0.0.0", 27017, 0, "sa", self.RootPassword, 5000, "", false)
	if err != nil {
		return nil, err
	}

	return p, nil
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
		"db=db.getSiblingDB('%s_db');db=db.getSiblingDB('admin');db.changeUserPassword('%s','%s');",
		"sa", rqst.NewPassword)

	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", rqst.OldPassword, changeRootPasswordScript)
	if err != nil {
		log.Println("---> fail to run script: ")
		log.Println(changeRootPasswordScript)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.saveConfig()

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + self.Domain + string(os.PathSeparator) + "_token")
	return &admin.SetRootPasswordResponse{
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
	services_discovery := services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")
	if services_discovery == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to "+rqst.DicorveryId)))
	}

	// Connect to the repository services.
	services_repository := services.NewServicesRepository_Client(rqst.RepositoryId, "services_repository")
	if services_repository == nil {
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

	err := services_discovery.PublishServiceDescriptor(serviceDescriptor)
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
	event, err := self.getEventHub()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = os.Remove(rqst.Path)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will send a event to be sure all server will update...
	data, _ := json.Marshal(serviceDescriptor)

	// Here I will send an event that the service has a new version...
	event.Publish(serviceDescriptor.PublisherId+":"+serviceDescriptor.Id+":SERVICE_PUBLISH_EVENT", data)

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

	var services_repository *services.ServicesRepository_Client
	for i := 0; i < len(descriptor.Repositories); i++ {
		services_repository = services.NewServicesRepository_Client(descriptor.Repositories[i], "services_repository")
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
			config["schemaPath"] = strings.ReplaceAll(string(os.PathSeparator)+dest+string(os.PathSeparator)+"schema.json", string(os.PathSeparator), "/")

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
	services_discovery = services.NewServicesDiscovery_Client(rqst.DicorveryId, "services_discovery")

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
				log.Println("---> remove service: ", id)
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

	if !Utility.IsLocal(config["Domain"].(string)) {
		return false
	}

	// set the domain of the service.
	config["Domain"] = self.Domain

	// get the config path.
	var process interface{}
	var proxyProcess interface{}

	process = config["Process"]
	proxyProcess = config["ProxyProcess"]

	// remove unused information...
	delete(config, "Process")
	delete(config, "ProxyProcess")

	// In case of persistence_server information must be save in a temp
	// file to be use by the interceptor for token validation.
	if config["Name"] == "persistence_server" {

		// I will wrote the info inside a stucture.
		infos := map[string]interface{}{"address": self.Domain, "name": "persistence_server", "pwd": self.RootPassword}
		infosStr, _ := Utility.ToJson(infos)

		err := ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_sa", []byte(infosStr), 0644)
		if err != nil {
			log.Panicln(err)
		}
	}

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
		log.Panicln("fail to save config file: ", err)
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
		self.ReadTimeout = time.Duration(Utility.ToInt(config["ReadTimeout"].(float64)))
		self.WriteTimeout = time.Duration(Utility.ToInt(config["WriteTimeout"].(float64)))
		self.IdleTimeout = time.Duration(Utility.ToInt(config["IdleTimeout"].(float64)))
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
	s := self.Services[serviceId].(map[string]interface{})
	if s == nil {
		return errors.New("No service found with id " + serviceId)
	}

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

	time.Sleep(2 * time.Second)

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
func (self *Globule) startInternalService(id string, port int, proxy int, hasTls bool) (*grpc.Server, error) {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, name+" service: ", log.LstdFlags))

	// Set the log information in case of crash...
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	var grpcServer *grpc.Server
	if hasTls {
		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
		certFile := self.creds + string(os.PathSeparator) + "server.crt"
		keyFile := self.creds + string(os.PathSeparator) + "server.pem"

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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors_.UnaryAuthInterceptor)}

		// Create the gRPC server with the credentials
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer()
	}

	reflection.Register(grpcServer)

	// Here I will create the service configuration object.
	s := make(map[string]interface{}, 0)
	s["Domain"] = self.Domain
	s["Port"] = port
	s["Proxy"] = proxy
	s["TLS"] = self.Protocol == "https"
	self.Services[id] = s

	// start the proxy
	err := self.startProxy(id, port, proxy)
	if err != nil {
		return nil, err
	}

	return grpcServer, nil
}

/** Stop mongod process **/
func (self *Globule) stopMongod() error {
	log.Println("---> stop mongo db")
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
		log.Println("---> start mongod without auth ", dataPath)
		err = mongod.Start()
		if err != nil {
			log.Println("fail to start mongo db", err)
			return err
		}

		self.waitForMongo(60, false)

		// Now I will create a new user name sa and give it all admin write.
		log.Println("---> create sa user in mongo db")
		createSaScript := fmt.Sprintf(
			`db=db.getSiblingDB('admin');db.createUser({ user: '%s', pwd: '%s', roles: ['userAdminAnyDatabase','userAdmin','readWrite','dbAdmin','clusterAdmin','readWriteAnyDatabase','dbAdminAnyDatabase']});`, "sa", self.RootPassword) // must be change...

		createSaCmd := exec.Command("mongo", "--eval", createSaScript)
		err = createSaCmd.Run()
		if err != nil {
			log.Println("---> fail to run script ", err)
			// remove the mongodb-data
			os.RemoveAll(dataPath)
			log.Println(createSaScript)
			return err
		}
		self.stopMongod()
	}

	// Now I will start mongod with auth available.
	log.Println("---> start mongo db whith auth ", dataPath)
	mongod := exec.Command("mongod", "--auth", "--port", "27017", "--dbpath", dataPath)
	err := mongod.Start()
	if err != nil {
		return err
	}

	// wait 15 seconds that the server restart.
	self.waitForMongo(60, true)

	// Get the list of all services method.
	return self.registerMethods()
}

/* Register a new Account */
func (self *Globule) RegisterAccount(ctx context.Context, rqst *ressource.RegisterAccountRqst) (*ressource.RegisterAccountRsp, error) {
	if rqst.ConfirmPassword != rqst.Password {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Password dosen't match!")))
	}

	// encode the password and keep it in the account itself.
	rqst.Account.Password = Utility.GenerateUUID(rqst.Password)

	// That service made user of persistence service.
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// first of all the Persistence service must be active.
	count, err := p.Count("local_ressource", "local_ressource", "Accounts", `{"name":"`+rqst.Account.Name+`"}`, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// one account already exist for the name.
	if count == 1 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("account with name "+rqst.Account.Name+" already exist!")))
	}

	// set the account object and set it basic roles.
	account := make(map[string]interface{})
	account["name"] = rqst.Account.Name
	account["email"] = rqst.Account.Email
	account["password"] = rqst.Account.Password

	// reference the guest role.
	guest := make(map[string]interface{}, 0)
	guest["$id"] = "guest"
	guest["$ref"] = "Roles"
	guest["$db"] = "local_ressource"
	account["roles"] = []map[string]interface{}{guest}

	// serialyse the account and save it.
	accountStr, err := Utility.ToJson(account)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will insert the account in the database.
	id, err := p.InsertOne("local_ressource", "local_ressource", "Accounts", accountStr, "")

	// Each account will have their own database and a use that can read and write
	// into it.
	// Here I will wrote the script for mongoDB...
	createUserScript := fmt.Sprintf(
		"db=db.getSiblingDB('%s_db');db.createCollection('user_data');db=db.getSiblingDB('admin');db.createUser({user: '%s', pwd: '%s',roles: [{ role: 'readWrite', db: '%s_db' }]});",
		rqst.Account.Name, rqst.Account.Name, rqst.Password, rqst.Account.Name)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, createUserScript)
	if err != nil {
		log.Println("---> fail to run script: ")
		log.Println(createUserScript)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = p.CreateConnection(rqst.Account.Name+"_db", rqst.Account.Name+"_db", "localhost", 27017, 0, rqst.Account.Name, rqst.Password, 5000, "", false)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))

	}

	// Now I will
	return &ressource.RegisterAccountRsp{
		Result: id,
	}, nil
}

//* Authenticate a account by it name or email.
// That function test if the password is the correct one for a given user
// if it is a token is generate and that token will be use by other service
// to validate permission over the requested ressource.
func (self *Globule) Authenticate(ctx context.Context, rqst *ressource.AuthenticateRqst) (*ressource.AuthenticateRsp, error) {
	log.Println("---> authenticate: ", rqst.Name)
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

	query := `{"name":"` + rqst.Name + `"}`

	// Can also be an email.
	if Utility.IsEmail(rqst.Name) {
		query = `{"email":"` + rqst.Name + `"}`
	}

	values, err := p.Find("local_ressource", "local_ressource", "Accounts", query, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	objects := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(values), &objects)

	if len(objects) == 0 {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("fail to retreive "+rqst.Name+" informations.")))
	}

	if objects[0]["password"].(string) != Utility.GenerateUUID(rqst.Password) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("wrong password for account "+objects[0]["name"].(string))))
	}

	// Generate a token to identify the user.
	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, objects[0]["name"].(string))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Open the user database connection.
	err = p.CreateConnection(objects[0]["name"].(string)+"_db", objects[0]["name"].(string)+"_db", "localhost", 27017, 0, objects[0]["name"].(string), rqst.Password, 5000, "", false)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))
	}

	// save the newly create token into the database.
	name, expireAt, _ := Interceptors_.ValidateToken(tokenString)
	p.DeleteOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, "")
	_, err = p.InsertOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
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
	name, expireAt, err := Interceptors_.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// If the token is older than seven day without being refresh then I retrun an error.
	if time.Now().Sub(time.Unix(expireAt, 0)) > (7 * 24 * time.Hour) {
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

		// That mean a newer token was already refresh.
		if lastTokenInfo["expireAt"].(int64) > expireAt {
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
	name, expireAt, _ = Interceptors_.ValidateToken(tokenString)

	p.DeleteOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`"}`, "")
	_, err = p.InsertOne("local_ressource", "local_ressource", "Tokens", `{"_id":"`+name+`","expireAt":`+Utility.ToString(expireAt)+`}`, "")
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

	// Try to delete the account...
	err = p.DeleteOne("local_ressource", "local_ressource", "Accounts", `{"name":"`+rqst.Name+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will drop the db user.
	dropUserScript := fmt.Sprintf(
		`db=db.getSiblingDB('admin');db.dropUser('%s', {w: 'majority', wtimeout: 4000})`,
		rqst.Name)

	// I will execute the sript with the admin function.
	err = p.RunAdminCmd("local_ressource", "sa", self.RootPassword, dropUserScript)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = p.DeleteConnection(rqst.Name + "_db")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ressource.DeleteAccountRsp{
		Result: rqst.Name,
	}, nil
}

//////////////////////////////// Services management  //////////////////////////

/**
 * Get access to the event services.
 */
func (self *Globule) getEventHub() (*event_client.Event_Client, error) {
	if self.clients["event_service"] == nil {
		return nil, errors.New("No event service was found on the server.")
	}

	return self.clients["event_service"].(*event_client.Event_Client), nil
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
	log.Println("----------> 2460 retreived services: ", data)
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
	if !Utility.Contains(rqst.Descriptor_.Discoveries, self.Domain) {
		rqst.Descriptor_.Discoveries = append(rqst.Descriptor_.Discoveries, self.Domain)
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
		log.Println("---> fail to get local_ressource connection", err)
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

	repositoryId := self.Domain
	// Now I will append the address of the repository into the service descriptor.
	if !Utility.Contains(bundle.Descriptor_.Repositories, repositoryId) {
		bundle.Descriptor_.Repositories = append(bundle.Descriptor_.Repositories, repositoryId)
		// Publish change into discoveries...
		for i := 0; i < len(bundle.Descriptor_.Discoveries); i++ {
			discoveryId := bundle.Descriptor_.Discoveries[i]
			discoveryService := services.NewServicesDiscovery_Client(discoveryId, "services_discovery")
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
		log.Println("---> fail to get local_ressource connection", err)
		return err
	}

	_, err = p.InsertOne("local_ressource", "local_ressource", "ServiceBundle", `{"_id":"`+id+`","checksum":"`+checksum+`"}`, "")

	return err
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

	// append itself to service discoveries...
	if !Utility.Contains(self.Discoveries, self.Domain) {
		self.Discoveries = append(self.Discoveries, self.Domain)
	}

	//log.Panicln(self.Domain)

	// hub --> channels --> subscriber(list of uuid's)
	subscribers := make(map[string]map[string][]string, 0)

	// Connect to service update events...
	for i := 0; i < len(self.Discoveries); i++ {
		eventHub := event_client.NewEvent_Client(self.Discoveries[i], "event_server")
		data_chan := make(chan []byte)
		subscribers[self.Discoveries[i]] = make(map[string][]string)
		for _, s := range self.Services {
			if s.(map[string]interface{})["PublisherId"] != nil {
				id := s.(map[string]interface{})["PublisherId"].(string) + ":" + s.(map[string]interface{})["Name"].(string) + ":SERVICE_PUBLISH_EVENT"
				if subscribers[self.Discoveries[i]][id] == nil {
					subscribers[self.Discoveries[i]][id] = make([]string, 0)
				}
				// each channel has it event...
				uuid, err := eventHub.Subscribe(id, data_chan)
				subscribers[self.Discoveries[i]][id] = append(subscribers[self.Discoveries[i]][id], uuid)
				if err == nil {
					go func() {
						for {
							select {
							case msg := <-data_chan:
								descriptor := new(services.ServiceDescriptor)
								json.Unmarshal(msg, descriptor)
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
																err = self.installService(descriptor)
																if err != nil {
																	log.Println("---> fail to install service ", err)
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
						}
					}()
				} else {
					log.Println("--> fail to connect to event channel", id, err)
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
					log.Println("---> disconnect ", id, channelId, uuids[i])
					eventHub.UnSubscribe(channelId, uuids[i])
				}
			}
		}

		// stop external service.
		for externalServiceId, _ := range self.ExternalApplications {
			self.stopExternalApplication(externalServiceId)
		}

		for _, value := range self.clients {
			value.Close()
		}

		// exit cleanly
		os.Exit(0)

	}()

	// Start the admin service to give access to server functionality from
	// client side.
	admin_server, err := self.startInternalService("Admin", self.AdminPort, self.AdminProxy, self.Protocol == "https") // must be accessible to all clients...
	if err == nil {
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.AdminPort))
		if err != nil {
			log.Fatalf("could not start admin service %s: %s", self.Domain, err)
		}

		admin.RegisterAdminServiceServer(admin_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start admin service!")
			go func() {
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
	}

	ressource_server, err := self.startInternalService("Ressource", self.RessourcePort, self.RessourceProxy, self.Protocol == "https")
	if err == nil {

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.RessourcePort))
		if err != nil {
			log.Fatalf("could not start ressource service %s: %s", self.Domain, err)
		}

		ressource.RegisterRessourceServiceServer(ressource_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start ressource service!")
			go func() {

				// no web-rpc server.
				if err := ressource_server.Serve(lis); err != nil {
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
		}()
	}

	// The service discovery.
	services_discovery_server, err := self.startInternalService("ServicesDiscovery", self.ServicesDiscoveryPort, self.ServicesDiscoveryProxy, self.Protocol == "https")
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesDiscoveryPort))
		if err != nil {
			log.Fatalf("could not start services discovery service %s: %s", self.Domain, err)
		}

		services.RegisterServiceDiscoveryServer(services_discovery_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start services discovery service!")
			go func() {

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
	}

	// The service repository
	services_repository_server, err := self.startInternalService("ServicesRepository", self.ServicesRepositoryPort, self.ServicesRepositoryProxy, self.Protocol == "https")
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.ServicesRepositoryPort))
		if err != nil {
			log.Fatalf("could not start services repository service %s: %s", self.Domain, err)
		}

		services.RegisterServiceRepositoryServer(services_repository_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start services repository service!")
			go func() {

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
	}

	// The Certificate Authority
	certificate_authority_server, err := self.startInternalService("CertificateAuthority", self.CertificateAuthorityPort, self.CertificateAuthorityProxy, false)
	if err == nil {
		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.CertificateAuthorityPort))
		if err != nil {
			log.Fatalf("could not certificate authority signing  service %s: %s", self.Domain, err)
		}

		ca.RegisterCertificateAuthorityServer(certificate_authority_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start certificate authority signing service!")
			go func() {

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
		if len(self.Certificate) == 0 {

			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/grpc_tls/server.csr run
			if !Utility.IsLocal(self.Domain) {
				// I need to remove the gRPC certificate and recreate it.
				Utility.RemoveDirContents(self.creds)
				self.GenerateServicesCertificates(self.CertPassword, self.CertExpirationDelay)
				self.initServices() // must restart the services with new certificates.
				err := self.ObtainCertificateForCsr()
				if err != nil {
					log.Panicln(err)
				}
			}
		}

		// Start https server.
		server := &http.Server{
			Addr: ":" + strconv.Itoa(self.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: self.Domain,
			},
		}

		// get the value from the configuration files.
		server.ListenAndServeTLS(self.certs+string(os.PathSeparator)+self.Certificate, self.creds+string(os.PathSeparator)+"server.pem")
	} else {

		// local - non secure connection.
		http.ListenAndServe(":"+strconv.Itoa(self.PortHttp), nil)
	}

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

/////////////////////////// Security stuff //////////////////////////////////

// That function will be access via http so event server or client will be able
// to get particular service configuration.
func getClientConfig(address string, name string) (map[string]interface{}, error) {
	config := make(map[string]interface{})
	config["Name"] = name
	config["Domain"] = address

	// First I will retreive the server configuration.
	serverConfig, err := getRemoteConfig(address)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	config["TLS"] = serverConfig["Protocol"].(string) == "https"
	if name == "services_discovery" {
		config["Port"] = Utility.ToInt(serverConfig["ServicesDiscoveryPort"])
	} else if name == "services_repository" {
		config["Port"] = Utility.ToInt(serverConfig["ServicesRepositoryPort"])
	} else if name == "admin" {
		config["Port"] = Utility.ToInt(serverConfig["AdminPort"])
	} else if name == "certificate_authority" {
		config["Port"] = Utility.ToInt(serverConfig["CertificateAuthorityPort"])
		config["TLS"] = false
	} else if name == "ressource" {
		config["Port"] = Utility.ToInt(serverConfig["RessourcePort"])
	} else if serverConfig["Services"].(map[string]interface{})[name] != nil {
		// get the service with the id egal to the given name.
		config["Port"] = Utility.ToInt(serverConfig["Services"].(map[string]interface{})[name].(map[string]interface{})["Port"])
		config["TLS"] = Utility.ToBool(serverConfig["Services"].(map[string]interface{})[name].(map[string]interface{})["TLS"])
	} else {
		return nil, errors.New("No service found whit name " + name + " exist on the server.")
	}

	// get / init credential values.
	if config["TLS"] == false {
		// set the credential function here
		config["KeyFile"] = ""
		config["CertFile"] = ""
		config["CertAuthorityTrust"] = ""
	} else {
		keyPath, certPath, caPath, err := getCredentialConfig(address)
		if err != nil {
			return nil, err
		}
		// set the credential function here
		config["KeyFile"] = keyPath
		config["CertFile"] = certPath
		config["CertAuthorityTrust"] = caPath
	}

	return config, nil
}

/**
 * Get the remote client configuration.
 */
func getRemoteConfig(address string) (map[string]interface{}, error) {

	if Utility.IsLocal(address) {
		return globule.getConfig(), nil
	}

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	resp, err = http.Get("http://" + address + ":10000/config")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var config map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getCredentialConfig(address string) (keyPath string, certPath string, caPath string, err error) {
	if Utility.IsLocal(address) {
		keyPath = globule.creds + string(os.PathSeparator) + "client.pem"
		certPath = globule.creds + string(os.PathSeparator) + "client.crt"
		caPath = globule.creds + string(os.PathSeparator) + "ca.crt"
		return
	}

	creds := globule.creds + string(os.PathSeparator) + address
	Utility.CreateDirIfNotExist(creds)

	// I will connect to the certificate authority of the server where the application must
	// be deployed. Certificate autority run wihtout tls.
	ca_client := ca.NewCa_Client(address, "certificate_authority")

	// Get the ca.crt certificate.
	ca_crt, err := ca_client.GetCaCertificate()
	if err != nil {
		log.Println(err)
		return
	}

	// Write the ca.crt file on the disk
	err = ioutil.WriteFile(creds+string(os.PathSeparator)+"ca.crt", []byte(ca_crt), 0400)
	if err != nil {
		log.Println(err)
		return
	}

	// Now I will generate the certificate for the client...
	// Step 1: Generate client private key.
	err = globule.GenerateClientPrivateKey(creds, globule.CertPassword)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 2: Generate the client signing request.
	err = globule.GenerateClientCertificateSigningRequest(creds, globule.CertPassword, address)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 3: Generate client signed certificate.
	client_csr, err := ioutil.ReadFile(creds + string(os.PathSeparator) + "client.csr")
	if err != nil {
		log.Println(err)
		return
	}

	// Sign the certificate from the server ca...
	client_crt, err := ca_client.SignCertificate(string(client_csr))
	if err != nil {
		log.Println(err)
		return
	}

	// Write bact the client certificate in file on the disk
	err = ioutil.WriteFile(creds+string(os.PathSeparator)+"client.crt", []byte(client_crt), 0400)
	if err != nil {
		log.Println(err)
		return
	}

	// Now ask the ca to sign the certificate.

	// Step 4: Convert to pem format.
	err = globule.KeyToPem("client", creds, globule.CertPassword)
	if err != nil {
		log.Println(err)
		return
	}

	// set the credential paths.
	keyPath = creds + string(os.PathSeparator) + "client.pem"
	certPath = creds + string(os.PathSeparator) + "client.crt"
	caPath = creds + string(os.PathSeparator) + "ca.crt"

	return
}

/////////////////////////////////////////////////////////////////////////////
// Certificate Authority Service
/////////////////////////////////////////////////////////////////////////////

// Signed certificate request (CSR)
func (self *Globule) SignCertificate(ctx context.Context, rqst *ca.SignCertificateRequest) (*ca.SignCertificateResponse, error) {

	// first of all I will save the incomming file into a temporary file...
	client_csr_path := os.TempDir() + string(os.PathSeparator) + Utility.RandomUUID()
	err := ioutil.WriteFile(client_csr_path, []byte(rqst.Csr), 0644)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

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

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	}

	// I will read back the crt file.
	client_crt, err := ioutil.ReadFile(client_crt_path)

	// remove the tow temporary files.
	os.Remove(client_crt_path)
	os.Remove(client_csr_path)

	return &ca.SignCertificateResponse{
		Crt: string(client_crt),
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

//////////////////////////////// Certificate Authority ///////////////////////
// Generate the Certificate Authority private key file (this shouldn't be shared in real life)
func (self *Globule) GenerateAuthorityPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "ca.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the Authority private key")
	}
	return nil
}

// Certificate Authority trust certificate (this should be shared whit users)
func (self *Globule) GenerateAuthorityTrustCertificate(path string, pwd string, expiration_delay int, domain string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "ca.crt") {
		return nil
	}
	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-new")
	args = append(args, "-x509")
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-key")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}

	return nil
}

/////////////////////// Server Keys //////////////////////////////////////////

// Server private key, password protected (this shoudn't be shared)
func (self *Globule) GenerateSeverPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "server.key") {
		return nil
	}
	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server private key")
	}
	return nil
}

// Generate client private key and certificate.
func (self *Globule) GenerateClientPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "client.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.pass.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate client private key " + err.Error())
	}

	args = make([]string, 0)
	args = append(args, "rsa")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+"client.pass.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.key")

	err = exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate client private key " + err.Error())
	}

	// Remove the file.
	os.Remove(path + string(os.PathSeparator) + "client.pass.key")
	return nil
}

func (self *Globule) GenerateClientCertificateSigningRequest(path string, pwd string, domain string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "client.csr") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+string(os.PathSeparator)+"client.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println(args)
		return errors.New("Fail to generate client certificate signing request.")
	}

	return nil
}

func (self *Globule) GenerateSignedClientCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + string(os.PathSeparator) + "client.crt") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+"client.csr")
	args = append(args, "-CA")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.crt")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println("fail to get the signed server certificate")
	}

	return nil
}

// Server certificate signing request (this should be shared with the CA owner)
func (self *Globule) GenerateServerCertificateSigningRequest(path string, pwd string, domain string) error {

	if Utility.Exists(path + string(os.PathSeparator) + "sever.crs") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+string(os.PathSeparator)+"server.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server certificate signing request.")
	}
	return nil
}

// Server certificate signed by the CA (this would be sent back to the client by the CA owner)
func (self *Globule) GenerateSignedServerCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + string(os.PathSeparator) + "sever.crt") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+"server.csr")
	args = append(args, "-CA")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.crt")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println("fail to get the signed server certificate")
	}

	return nil
}

// Conversion of server.key into a format gRpc likes (this shouldn't be shared)
func (self *Globule) KeyToPem(name string, path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + name + ".pem") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "pkcs8")
	args = append(args, "-topk8")
	args = append(args, "-nocrypt")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+name+".key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+name+".pem")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server.pem key from server.key")
	}

	return nil
}

/**
 * That function is use to generate services certificates.
 * Private ca.key, server.key, server.pem, server.crt
 * Share ca.crt (needed by the client), server.csr (needed by the CA)
 */
func (self *Globule) GenerateServicesCertificates(pwd string, expiration_delay int) {
	var domain = self.Domain

	// Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
	err := self.GenerateAuthorityPrivateKey(self.creds, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	err = self.GenerateAuthorityTrustCertificate(self.creds, pwd, expiration_delay, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Setp 2: Generate the server Private Key (server.key)
	err = self.GenerateSeverPrivateKey(self.creds, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Setp 3: Get a certificate signing request from the CA (server.csr)
	err = self.GenerateServerCertificateSigningRequest(self.creds, pwd, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 4: Sign the certificate with the CA we create(it's called self signing) - server.crt
	err = self.GenerateSignedServerCertificate(self.creds, pwd, expiration_delay)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 5: Convert the server Certificate to .pem format (server.pem) - usable by gRpc
	err = self.KeyToPem("server", self.creds, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 6: Generate client private key.
	err = self.GenerateClientPrivateKey(self.creds, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 7: Generate the client signing request.
	err = self.GenerateClientCertificateSigningRequest(self.creds, pwd, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 8: Generate client signed certificate.
	err = self.GenerateSignedClientCertificate(self.creds, pwd, expiration_delay)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 9: Convert to pem format.
	err = self.KeyToPem("client", self.creds, pwd)
	if err != nil {
		log.Println(err)
		return
	}
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
		log.Panicln(err)
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
	self.Certificate = self.Domain + ".crt"
	self.CertificateAuthorityBundle = self.Domain + ".issuer.crt"

	// Save the certificate in the cerst folder.
	ioutil.WriteFile(self.certs+string(os.PathSeparator)+self.Certificate, resource.Certificate, 0400)
	ioutil.WriteFile(self.certs+string(os.PathSeparator)+self.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	self.saveConfig()

	return nil
}
