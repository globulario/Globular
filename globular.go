package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"

	// "github.com/gorilla/mux"
	ps "github.com/mitchellh/go-ps"

	"crypto/tls"
	"crypto/x509"
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
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Utility"

	// Admin service
	"github.com/davecourtois/Globular/admin"
	// Ressource service
	"github.com/davecourtois/Globular/ressource"
	// Services management service
	"github.com/davecourtois/Globular/services"

	// Interceptor for authentication, event, log...
	Interceptors "github.com/davecourtois/Globular/Interceptors/Authenticate"
	Interceptors_ "github.com/davecourtois/Globular/Interceptors/server"

	// Client services.
	"context"

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
	"github.com/emicklei/proto"

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
	ServicesDiscoveryPort      int    // The services discovery port
	ServicesDiscoveryProxy     int    // The ressource management proxy port
	ServicesRepositoryPort     int    // The services discovery port
	ServicesRepositoryProxy    int    // The ressource management proxy port
	Protocol                   string // The protocol of the service.
	Services                   map[string]interface{}
	ExternalApplications       map[string]ExternalApplication
	Domain                     string                   // The domain (subdomain) name of your application
	DNS                        []map[string]interface{} // Contain a list of domain name server where that computer use as sub-domain
	ReadTimeout                time.Duration
	WriteTimeout               time.Duration
	IdleTimeout                time.Duration
	SessionTimeout             time.Duration
	CertExpirationDelay        int
	Certificate                string
	CertificateAuthorityBundle string

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

	// set default values.
	g.IdleTimeout = 120
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.ReadTimeout = 5
	g.WriteTimeout = 5
	g.CertExpirationDelay = 365

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

	// Here I will kill proxies if there are running.
	killProcessByName("grpcwebproxy")

	// Here it suppose to be only one server instance per computer.
	g.jwtKey = []byte(Utility.RandomUUID())
	err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_key", []byte(g.jwtKey), 0644)
	if err != nil {
		log.Panicln(err)
	}

	// The token that identify the server with other services
	token, _ := Interceptors.GenerateToken(g.jwtKey, g.SessionTimeout, "sa")
	err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_token", []byte(token), 0644)
	if err != nil {
		log.Panicln(err)
	}

	// Here I will start the refresh token loop to refresh the server token.
	// the token will be refresh 10 milliseconds before expiration.
	ticker := time.NewTicker((g.SessionTimeout - 10) * time.Millisecond)
	go func() {

		for {
			select {
			case <-ticker.C:
				token, _ := Interceptors.GenerateToken(g.jwtKey, g.SessionTimeout, "sa")
				err = ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_token", []byte(token), 0644)
				log.Println("new sa token generated: ", token)
				if err != nil {
					log.Panicln(err)
				}

				// close existing client and re-init it
				for id, s := range g.Services {
					if s.(map[string]interface{})["TLS"] != nil && s.(map[string]interface{})["Protocol"] != nil {
						if s.(map[string]interface{})["TLS"].(bool) == true {
							// reconnect the client with it new token value.
							name := s.(map[string]interface{})["Name"].(string)
							g.initClient(id, name, token)
						}
					}
				}
			}
		}
	}()

	g.initClients()

	// Keep in global var to by http handlers.
	globule = g

	return g
}

/**
 * Set the ip for a given sub-domain compose of Name + DNS domain.
 */
func (self *Globule) registerIpToDns() {
	if self.DNS != nil {
		if len(self.DNS) > 0 {
			for i := 0; i < len(self.DNS); i++ {
				config := self.DNS[i]
				domain := config["Domain"].(string)
				port := int(config["Port"].(float64))
				hasTls := config["TLS"].(bool)
				crt := config["CertFile"].(string)
				key := config["KeyFile"].(string)
				ca := config["CertAuthorityTrust"].(string)
				token := "" // The token...
				client := dns_client.NewDns_Client(domain, port, hasTls, key, crt, ca, token)
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
 * Get the list of process id by it name.
 */
func getProcessIdsByName(name string) ([]int, error) {
	processList, err := ps.Processes()
	if err != nil {
		return nil, errors.New("ps.Processes() Failed, are you using windows?")
	}

	pids := make([]int, 0)

	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		if strings.HasPrefix(process.Executable(), name) {
			pids = append(pids, process.Pid())
		}
	}

	return pids, nil
}

/**
 * Kill a process with a given name.
 */
func killProcessByName(name string) error {
	pids, err := getProcessIdsByName(name)
	if err != nil {
		return err
	}

	for i := 0; i < len(pids); i++ {
		proc, err := os.FindProcess(pids[i])

		if err != nil {
			log.Println(err)
		}
		log.Println("Kill ", name, " pid ", pids[i])
		// Kill the process
		if !strings.HasPrefix(name, "Globular") {
			proc.Kill()
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
	proxyPath := "/bin/grpcwebproxy"
	if string(os.PathSeparator) == "\\" {
		proxyPath += ".exe" // in case of windows.
	}

	proxyBackendAddress := self.Domain + ":" + Utility.ToString(port)
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
		proxyArgs = append(proxyArgs, "--server_http_tls_port="+Utility.ToString(proxy))

		/* in case of public domain server files **/
		proxyArgs = append(proxyArgs, "--server_tls_key_file="+self.creds+string(os.PathSeparator)+"server.pem")

		proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+self.certs+string(os.PathSeparator)+self.CertificateAuthorityBundle)
		proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.certs+string(os.PathSeparator)+self.Certificate)

	} else {
		log.Println("---> start non secure service: ", srv.(map[string]interface{})["Name"])

		// Now I will save the file with those new information in it.
		proxyArgs = append(proxyArgs, "--run_http_server=true")
		proxyArgs = append(proxyArgs, "--run_tls_server=false")
		proxyArgs = append(proxyArgs, "--server_http_debug_port="+Utility.ToString(proxy))
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
 * Start a local service.
 */
func (self *Globule) startService(s map[string]interface{}) (int, int, error) {
	var err error
	// test if the service is distant or not.
	if s["Domain"].(string) != self.Domain {
		if s["Domain"].(string) != "localhost" {
			return -1, -1, errors.New("Can not start a distant service localy!")
		}
	}

	// set the domain of the service.
	s["Domain"] = self.Domain

	// if the service already exist.
	srv := self.Services[s["Id"].(string)]
	if srv != nil {
		if srv.(map[string]interface{})["Process"] != nil {
			if reflect.TypeOf(srv.(map[string]interface{})["Process"]).String() == "*exec.Cmd" {
				srv.(map[string]interface{})["Process"].(*exec.Cmd).Process.Kill()
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
		killProcessByName(s["Name"].(string))

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
		if err != nil {
			log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
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

		// get the token from file
		token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_token")

		// Init it configuration.
		self.initClient(s["Name"].(string), s["Name"].(string), string(token))

		// save it to the config.
		self.saveConfig()

	} else if s["Protocol"].(string) == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(s["Name"].(string), "Globular") {
			// Kill previous instance of the program.
			killProcessByName(s["Name"].(string))
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

	// Return the pid of the service.
	if s["ProxyProcess"] != nil {
		return s["Process"].(*exec.Cmd).Process.Pid, s["ProxyProcess"].(*exec.Cmd).Process.Pid, nil
	}

	return s["Process"].(*exec.Cmd).Process.Pid, -1, nil
}

/**
 * Call once when the server start.
 */
func (self *Globule) initService(s map[string]interface{}) {
	if s["Protocol"] == nil {
		return
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
			self.startService(s)
			self.saveConfig()
		}
	}
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
					json.Unmarshal(config, &s)
					if s["Protocol"] != nil {
						// If a configuration file exist It will be use to start services,
						// otherwise the service configuration file will be use.
						path_ := path[:strings.LastIndex(path, string(os.PathSeparator))]
						s["Id"] = s["Name"].(string)
						servicePath := path_ + string(os.PathSeparator) + s["Name"].(string)
						strings.Index(servicePath, self.path)

						s["servicePath"] = strings.Replace(strings.Replace(servicePath, self.path, "", -1), "\\", "/", -1)

						s["configPath"] = strings.Replace(strings.Replace(path, self.path, "", -1), "\\", "/", -1)

						//self.initService(s)
						self.Services[s["Name"].(string)] = s

					}
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
			name := info.Name()[0:strings.Index(info.Name(), ".")] + "_server"
			if self.Services[name] != nil {
				s := self.Services[name].(map[string]interface{})
				s["protoPath"] = strings.Replace(strings.Replace(path, self.path, "", -1), "\\", "/", -1)
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
	for _, s := range self.Services {
		// Remove existing process information.
		delete(s.(map[string]interface{}), "Process")
		delete(s.(map[string]interface{}), "ProxyProcess")

		self.initService(s.(map[string]interface{}))
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
	} else {
		log.Println("--> no dns service found!")
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
		guest["actions"] = []string{"/admin.AdminService/GetConfig", "/ressource.RessourceService/RegisterAccount", "/ressource.RessourceService/Authenticate"}
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
func (self *Globule) initClient(id string, name string, token string) {

	if self.Services[id] == nil {
		return
	}

	// The connection address.
	// Create the address with the ip address.
	hasTLS := self.Services[id].(map[string]interface{})["TLS"].(bool)
	domain := self.Services[id].(map[string]interface{})["Domain"].(string)
	port := int(self.Services[id].(map[string]interface{})["Port"].(float64))
	name = strings.Split(name, "_")[0]

	fct := "New" + strings.ToUpper(name[0:1]) + name[1:] + "_Client"

	// Set the files.
	keyFile := self.creds + string(os.PathSeparator) + "client.crt"
	certFile := self.creds + string(os.PathSeparator) + "client.key"
	caFile := self.creds + string(os.PathSeparator) + "ca.crt"

	results, err := Utility.CallFunction(fct, domain, port, hasTLS, certFile, keyFile, caFile, token)
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
	config := make(map[string]interface{}, 0)
	config["Name"] = self.Name
	config["PortHttp"] = self.PortHttp
	config["PortHttps"] = self.PortHttps
	config["AdminPort"] = self.AdminPort
	config["AdminProxy"] = self.AdminProxy
	config["RessourcePort"] = self.RessourcePort
	config["RessourceProxy"] = self.RessourceProxy
	config["ServicesDiscoveryPort"] = self.ServicesDiscoveryPort
	config["ServicesDiscoveryProxy"] = self.ServicesDiscoveryProxy
	config["ServicesRepositoryPort"] = self.ServicesRepositoryPort
	config["ServicesRepositoryProxy"] = self.ServicesRepositoryProxy
	config["Protocol"] = self.Protocol
	config["Domain"] = self.Domain
	config["ReadTimeout"] = self.ReadTimeout
	config["WriteTimeout"] = self.WriteTimeout
	config["IdleTimeout"] = self.IdleTimeout
	config["SessionTimeout"] = self.SessionTimeout
	config["CertExpirationDelay"] = self.CertExpirationDelay

	// return the full service configuration.
	config["Services"] = self.Services

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

func (self *Globule) GetConfig(ctx context.Context, rqst *admin.GetConfigRequest) (*admin.GetConfigResponse, error) {
	config := make(map[string]interface{}, 0)
	config["Name"] = self.Name
	config["PortHttp"] = self.PortHttp
	config["PortHttps"] = self.PortHttps
	config["AdminPort"] = self.AdminPort
	config["AdminProxy"] = self.AdminProxy
	config["RessourcePort"] = self.RessourcePort
	config["RessourceProxy"] = self.RessourceProxy
	config["ServicesDiscoveryPort"] = self.ServicesDiscoveryPort
	config["ServicesDiscoveryProxy"] = self.ServicesDiscoveryProxy
	config["ServicesRepositoryPort"] = self.ServicesRepositoryPort
	config["ServicesRepositoryProxy"] = self.ServicesRepositoryProxy
	config["Protocol"] = self.Protocol
	config["Domain"] = self.Domain
	config["ReadTimeout"] = self.ReadTimeout
	config["WriteTimeout"] = self.WriteTimeout
	config["IdleTimeout"] = self.IdleTimeout
	config["CertExpirationDelay"] = self.CertExpirationDelay

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
		config["Services"].(map[string]interface{})[name] = s
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

/**
 * Start the monitoring service with prometheus.
 */
func (self *Globule) startMonitoring() error {
	killProcessByName("prometheus")
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

	// if not I will create one.
	log.Println("---> local_ressource not exist I will try to create the connection with sa and password", self.RootPassword)

	// Connect to the database here.
	err := p.CreateConnection("local_ressource", "local_ressource", "0.0.0.0", 27017, 0, "sa", self.RootPassword, 5000, "", false)
	if err != nil {
		log.Println(`---> Fail to create  the connection "local_ressource"`)
		return nil, err
	}

	return p, nil
}

//Set the root password
func (self *Globule) SetRootPassword(ctx context.Context, rqst *admin.SetRootPasswordRqst) (*admin.SetRootPasswordRsp, error) {
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

	token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_token")
	return &admin.SetRootPasswordRsp{
		Token: string(token),
	}, nil

}

// return true if the configuation has change.
func (self *Globule) saveServiceConfig(config map[string]interface{}) bool {

	if config["Domain"].(string) != self.Domain {
		if config["Domain"].(string) != "localhost" {
			return false
		}
	}

	// set the domain of the service.
	config["Domain"] = self.Domain

	// get the config path.
	var servicePath string
	var process interface{}
	var proxyProcess interface{}

	servicePath = config["servicePath"].(string)
	process = config["Process"]
	proxyProcess = config["ProxyProcess"]

	var path string
	path = self.path + config["servicePath"].(string)               // the path of the executable.
	path = strings.Replace(path, "/", string(os.PathSeparator), -1) // set the os path separator.
	path = path[:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator) + "config.json"

	// remove unused information...
	delete(config, "Process")
	delete(config, "ProxyProcess")

	// remove this info from the file to be save.
	delete(config, "servicePath")

	// In case of persistence_server information must be save in a temp
	// file to be use by the interceptor for token validation.
	if config["Name"] == "persistence_server" {
		persistenceDomain := config["Domain"].(string)
		persistencePort := int(config["Port"].(float64))

		certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
		certFile := self.creds + string(os.PathSeparator) + "server.crt"
		keyFile := self.creds + string(os.PathSeparator) + "server.pem"

		// I will wrote the info inside a stucture.
		infos := map[string]interface{}{"domain": persistenceDomain, "port": persistencePort, "certAuthorityTrust": certAuthorityTrust, "certFile": certFile, "keyFile": keyFile, "pwd": self.RootPassword}

		infosStr, _ := Utility.ToJson(infos)

		err := ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_sa", []byte(infosStr), 0644)
		if err != nil {
			log.Panicln(err)
		}
	}

	// so here I will get the previous information...
	f, err := os.Open(path)

	if err == nil {
		b, err := ioutil.ReadAll(f)
		if err == nil {
			config_ := make(map[string]interface{})
			json.Unmarshal(b, &config_)
			if reflect.DeepEqual(config_, config) {
				f.Close()
				// set back the path's info.
				config["servicePath"] = servicePath
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
	err = ioutil.WriteFile(path, []byte(jsonStr), 0644)
	if err != nil {
		log.Panicln("fail to save config file: ", err)
	}

	// set back internal infos...
	config["servicePath"] = servicePath
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
		} else if config["Services"] != nil {
			// Here I will save the configuration
			self.Name = config["Name"].(string)
			self.PortHttp = Utility.ToInt(config["PortHttp"].(float64))
			self.PortHttps = Utility.ToInt(config["PortHttps"].(float64))
			self.AdminPort = Utility.ToInt(config["AdminPort"].(float64))
			self.AdminProxy = Utility.ToInt(config["AdminProxy"].(float64))
			self.RessourcePort = Utility.ToInt(config["RessourcePort"].(float64))
			self.RessourceProxy = Utility.ToInt(config["RessourceProxy"].(float64))
			self.ServicesDiscoveryPort = Utility.ToInt(config["ServicesDiscoveryPort"].(float64))
			self.ServicesDiscoveryProxy = Utility.ToInt(config["ServicesDiscoveryProxy"].(float64))
			self.ServicesRepositoryPort = Utility.ToInt(config["ServicesRepositoryPort"].(float64))
			self.ServicesRepositoryProxy = Utility.ToInt(config["ServicesRepositoryProxy"].(float64))
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
			// save the application server.
			self.saveConfig()
		}
	}

	// return the new configuration file...
	result, _ := Utility.ToJson(config)
	return &admin.SaveConfigResponse{
		Result: result,
	}, nil
}

// Stop a service
func (self *Globule) StopService(ctx context.Context, rqst *admin.StopServiceRequest) (*admin.StopServiceResponse, error) {
	s := self.Services[rqst.ServiceId]
	if s == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No service found with id "+rqst.ServiceId)))
	}
	err := s.(map[string]interface{})["Process"].(*exec.Cmd).Process.Kill()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if s.(map[string]interface{})["ProxyProcess"] != nil {
		err := s.(map[string]interface{})["ProxyProcess"].(*exec.Cmd).Process.Kill()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
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
func (self *Globule) startInternalService(id string, port int, proxy int) (*grpc.Server, error) {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, name+" service: ", log.LstdFlags))

	// Set the log information in case of crash...
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	var grpcServer *grpc.Server
	if self.Protocol == "https" {
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

	// Get the persistence connection
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, err
	}

	// in case of sa user.
	if rqst.Password == self.RootPassword && rqst.Name == "sa" {
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
func (self *Globule) PublishService(ctx context.Context, rqst *services.PublishServiceRequest) (*services.PublishServiceResponse, error) {

	// Here I will save the descriptor inside the storage...
	p, err := self.getPersistenceSaConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append the self domain to the list of discoveries where the services can
	if !Utility.Contains(rqst.Descriptor_.Discorveries, self.Domain) {
		rqst.Descriptor_.Discorveries = append(rqst.Descriptor_.Discorveries, self.Domain+":"+strconv.Itoa(self.ServicesDiscoveryPort))
	}

	// Here I will test if the services already exist...
	_, err = p.FindOne("local_ressource", "local_ressource", "Services", `{"id":"`+rqst.Descriptor_.Id+`", "publisherId":"`+rqst.Descriptor_.PublisherId+`", "version":"`+rqst.Descriptor_.Version+`"}`, "")
	if err == nil {
		// Update existing descriptor.

		// The list of discoveries...
		discoveries, err := Utility.ToJson(rqst.Descriptor_.Discorveries)
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

		// So here I will send an plublish event...
		event, err := self.getEventHub()
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// Here I will send an event that the service has a new version...
		event.Publish(rqst.Descriptor_.PublisherId+":"+rqst.Descriptor_.Id+":NEW_SERVICE_VERSION_EVENT", data)

	}

	return &services.PublishServiceResponse{
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
	if bundle.Plaform == services.Platform_LINUX32 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%LINUX32"
	} else if bundle.Plaform == services.Platform_LINUX64 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%LINUX64"
	} else if bundle.Plaform == services.Platform_WIN32 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%WIN32"
	} else if bundle.Plaform == services.Platform_WIN64 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%WIN64"
	}

	path := self.data + string(os.PathSeparator) + "service-repository"

	var err error
	// the file must be a zipped archive that contain a .proto, .config and executable.
	bundle.Binairies, err = ioutil.ReadFile(path + string(os.PathSeparator) + id + ".tar.gz")

	if err != nil {
		return err
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
	var id string
	if bundle.Plaform == services.Platform_LINUX32 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%LINUX32"
	} else if bundle.Plaform == services.Platform_LINUX64 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%LINUX64"
	} else if bundle.Plaform == services.Platform_WIN32 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%WIN32"
	} else if bundle.Plaform == services.Platform_WIN64 {
		id = bundle.Descriptor_.Id + "%" + bundle.Descriptor_.PublisherId + "%" + bundle.Descriptor_.Version + "%WIN64"
	}

	repositoryId := self.Domain + ":" + strconv.Itoa(self.ServicesRepositoryPort)
	// Now I will append the address of the repository into the service descriptor.
	if !Utility.Contains(bundle.Descriptor_.Repositories, repositoryId) {
		bundle.Descriptor_.Repositories = append(bundle.Descriptor_.Repositories, repositoryId)
		// Publish change into discoveries...
		for i := 0; i < len(bundle.Descriptor_.Discorveries); i++ {
			discoveryId := bundle.Descriptor_.Discorveries[i]

			// discovery id is compose of the domain:port
			domain := strings.Split(discoveryId, ":")[0]
			port := Utility.ToInt(strings.Split(discoveryId, ":")[1])

			// TODO set certificates and has tls value here...
			discoveryService := services.NewServicesDiscovery_Client(domain, port, false, "", "", "", "")
			discoveryService.PublishService(bundle.Descriptor_)
		}
	}

	path := self.data + string(os.PathSeparator) + "service-repository"
	Utility.CreateDirIfNotExist(path)

	// the file must be a zipped archive that contain a .proto, .config and executable.
	err = ioutil.WriteFile(path+string(os.PathSeparator)+id+".tar.gz", bundle.Binairies, 777)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

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
	admin_server, err := self.startInternalService("Admin", self.AdminPort, self.AdminProxy)
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
					log.Fatalf("failed to serve: %v", err)
				}
				log.Println("Adim grpc service is closed")
			}()
			// Wait for signal to stop.
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			killProcessByName("mongod")
			killProcessByName("prometheus")
		}()
	}

	ressource_server, err := self.startInternalService("Ressource", self.RessourcePort, self.RessourceProxy)
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
					log.Fatalf("failed to serve: %v", err)
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
	services_discovery_server, err := self.startInternalService("ServicesDiscovery", self.ServicesDiscoveryPort, self.ServicesDiscoveryProxy)
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
					log.Fatalf("failed to serve: %v", err)
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
	services_repository_server, err := self.startInternalService("ServicesRepository", self.ServicesRepositoryPort, self.ServicesRepositoryProxy)
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
					log.Fatalf("failed to serve: %v", err)
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
		if !Utility.Exists(self.certs + string(os.PathSeparator) + "certificates") {

			// Here is the command to be execute in order to ge the certificates.
			// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/grpc_tls/server.csr run
			args := make([]string, 0)
			// Now I will get the new certificate.
			legoPath := self.path + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "lego"
			if string(os.PathSeparator) == "\\" {
				legoPath += ".exe" // in case of windows.
			}

			args = append(args, `--email="`+self.AdminEmail+`"`)
			args = append(args, `--accept-tos`)
			args = append(args, `--key-type=rsa4096`)
			args = append(args, `--path=`+self.certs)
			args = append(args, `--http`)
			args = append(args, `--csr=`+self.creds+string(os.PathSeparator)+"server.csr")
			args = append(args, `run`)

			cmdStr := legoPath
			for i := 0; i < len(args); i++ {
				cmdStr += " " + args[i]
			}

			log.Println(cmdStr)

			// Now I will run the command...
			legoCmd := exec.Command(legoPath, args...)
			err := legoCmd.Run()

			if err != nil {
				log.Panicln("Fail to generate certificate with lego", err)
			}

		}

		// Set the certificates paths...
		self.Certificate = "certificates" + string(os.PathSeparator) + self.Domain + ".crt"
		self.CertificateAuthorityBundle = "certificates" + string(os.PathSeparator) + self.Domain + ".issuer.crt"

		// save the config with the values.
		self.saveConfig()

		// Start https server.
		server := &http.Server{
			Addr: ":" + Utility.ToString(self.PortHttps),
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
