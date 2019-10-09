package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

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
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Utility"

	// Admin service
	"github.com/davecourtois/Globular/admin"
	"github.com/davecourtois/Globular/ressource"

	// Interceptor for authentication, event, log...
	"github.com/davecourtois/Globular/Interceptors/Authenticate"
	Interceptors_ "github.com/davecourtois/Globular/Interceptors/server"

	// Client services.
	"context"

	"github.com/davecourtois/Globular/catalog/catalog_client"
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
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	root    string
	globule *Globule
)

const serviceStartDelay = 2 // wait tow second.

type ExternalService struct {
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
	Name                string // The service name
	PortHttp            int    // The port of the http file server.
	PortHttps           int    // The secure port
	AdminPort           int    // The admin port
	AdminProxy          int    // The admin proxy port.
	RessourcePort       int    // The ressource management service port
	RessourceProxy      int    // The ressource management proxy port
	Protocol            string // The protocol of the service.
	IP                  string // The local address...
	Services            map[string]interface{}
	ExternalServices    map[string]ExternalService // MongoDB, Prometheus, etc...
	Domain              string                     // The domain name of your application
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	SessionTimeout      time.Duration
	CertExpirationDelay int
	RootPassword        string

	// The list of method supported by this server.
	methods []string

	// Local info.
	webRoot string // The root of the http file server.

	creds string // gRpc credential
	certs string // https certificates

	path string // The path of the exec...

	// The list of avalaible services.
	services map[string]interface{}

	// The map of client...
	clients map[string]api.Client

	// Create the JWT key used to create the signature
	jwtKey []byte

	sync.RWMutex
}

/**
 * Globule constructor.
 */
func NewGlobule(port int) *Globule {
	// Here I will initialyse configuration.
	g := new(Globule)
	g.RootPassword = "adminadmin"
	g.PortHttp = 80
	g.PortHttps = 443 // The default port number.
	g.Name = strings.Replace(Utility.GetExecName(os.Args[0]), ".exe", "", -1)

	g.Protocol = "http"
	var err error
	g.Domain, err = os.Hostname()
	if err != nil {
		g.Domain = "localhost"
	}

	g.IP = Utility.MyIP()
	g.AdminPort = 10001
	g.AdminProxy = 10002
	g.RessourcePort = 10003
	g.RessourceProxy = 10004

	// set default values.
	g.IdleTimeout = 120
	g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.ReadTimeout = 5
	g.WriteTimeout = 5
	g.CertExpirationDelay = 365
	// Set the service map.
	g.services = make(map[string]interface{}, 0)

	// Set the share service info...
	g.Services = make(map[string]interface{}, 0)

	// Set external map services.
	g.ExternalServices = make(map[string]ExternalService, 0)

	// Set the map of client.
	g.clients = make(map[string]api.Client, 0)
	g.initClients()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	g.path = dir // keep the installation path.

	if err == nil {
		g.webRoot = dir + string(os.PathSeparator) + "WebRoot" // The default directory to server.
		Utility.CreateDirIfNotExist(g.webRoot)                 // Create the directory if it not exist.
		file, err := ioutil.ReadFile(g.webRoot + string(os.PathSeparator) + "config.json")
		// Init the servce with the default port address
		if err == nil {
			json.Unmarshal([]byte(file), g)
		}
	}

	// Create the creds directory if it not already exist.
	g.creds = dir + string(os.PathSeparator) + "creds"
	Utility.CreateDirIfNotExist(g.creds)

	// https certificates.
	g.certs = dir + string(os.PathSeparator) + "certs"
	Utility.CreateDirIfNotExist(g.certs)

	// keep the root in global variable for the file handler.
	root = g.webRoot
	globule = g

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
				if err != nil {
					log.Panicln(err)
				}

				// close existing client and re-init it
				for name, s := range g.services {
					if s.(map[string]interface{})["TLS"] != nil && s.(map[string]interface{})["Protocol"] != nil {
						if s.(map[string]interface{})["TLS"].(bool) == true {
							// reconnect the client with it new token value.
							g.initClient(name, token)
						}
					}
				}

			}
		}
	}()

	return g
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
func (self *Globule) startProxy(name string, port int, proxy int) error {
	srv := self.services[name]
	if srv.(map[string]interface{})["ProxyProcess"] != nil {
		srv.(map[string]interface{})["ProxyProcess"].(*exec.Cmd).Process.Kill()
	}

	// Now I will start the proxy that will be use by javascript client.
	proxyPath := self.path + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
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

		/** Self signed certificate for self signed server domain **/
		proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.certs+string(os.PathSeparator)+"server_cert.pem")
		proxyArgs = append(proxyArgs, "--server_tls_key_file="+self.certs+string(os.PathSeparator)+"server_key.pem")

		/* Use those files for public domain server **/
		//proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+self.path+"/sslforfree/ca_bundle.crt")
		//proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.path+"/sslforfree/certificate.crt")
		//proxyArgs = append(proxyArgs, "--server_tls_key_file="+self.path+"/sslforfree/private.key")

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
	proxyProcess := exec.Command(proxyPath, proxyArgs...)
	err := proxyProcess.Start()

	if err != nil {
		log.Println("Fail to start grpcwebproxy at port ", proxy, " with error ", err)
		return err
	}

	// save service configuration.
	srv.(map[string]interface{})["ProxyProcess"] = proxyProcess
	self.services[name] = srv

	return nil
}

/**
 * Start a given service.
 */
func (self *Globule) startService(s map[string]interface{}) (int, int, error) {
	var err error

	// if the service already exist.
	srv := self.services[s["Name"].(string)]
	if srv != nil {
		if srv.(map[string]interface{})["Process"] != nil {
			srv.(map[string]interface{})["Process"].(*exec.Cmd).Process.Kill()
		}
	}

	if s["Protocol"].(string) == "grpc" {

		// Stop the previous client if there one.
		if self.clients[s["Name"].(string)+"_service"] != nil {
			self.clients[s["Name"].(string)+"_service"].Close()
		}

		// The domain must be set in the sever configuration and not change after that.
		s["Domain"] = self.Domain // local services.
		s["Address"] = self.IP    // local services.

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

		if s["Name"].(string) == "file_server" {
			// File service need root...
			s["Root"] = globule.webRoot
			s["Process"] = exec.Command(s["servicePath"].(string), Utility.ToString(s["Port"]), globule.webRoot)
		} else {
			s["Process"] = exec.Command(s["servicePath"].(string), Utility.ToString(s["Port"]))
		}

		err = s["Process"].(*exec.Cmd).Start()
		if err != nil {
			log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
			return -1, -1, err
		}
		// Save configuration stuff.
		self.services[s["Name"].(string)] = s

		// Start the proxy.
		err = self.startProxy(s["Name"].(string), int(s["Port"].(float64)), int(s["Proxy"].(float64)))
		if err != nil {
			return -1, -1, err
		}

		// get back the service info with the proxy process in it
		s = self.services[s["Name"].(string)].(map[string]interface{})

		// public configuration info.
		s_ := make(map[string]interface{})

		// export public service values.
		s_["Address"] = s["Address"]
		s_["Domain"] = s["Domain"]
		s_["Proxy"] = s["Proxy"]
		s_["Port"] = s["Port"]
		s_["TLS"] = s["TLS"]
		self.Services[s["Name"].(string)] = s_

		// get the token from file
		token, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "globular_token")

		// Init it configuration.
		self.initClient(s["Name"].(string), string(token))

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
			}
			self.services[s["Name"].(string)] = s
		}
	}

	log.Println("Service ", s["Name"].(string), "is running at port", s["Port"], "it's proxy port is", s["Proxy"])

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

	if s["Protocol"].(string) == "grpc" {
		// The domain must be set in the sever configuration and not change after that.
		s["Domain"] = self.Domain // local services.
		s["Address"] = self.IP    // local services.

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
		// Now I will save the file with those new information in it.
		hasChange := self.saveServiceConfig(s)

		if hasChange || s["Process"] == nil {

			self.services[s["Name"].(string)] = s

			// start the service
			self.startService(s)

			// public configuration info.
			s_ := make(map[string]interface{})

			// export public service values.
			s_["Address"] = s["Address"]
			s_["Domain"] = s["Domain"]
			s_["Proxy"] = s["Proxy"]
			s_["Port"] = s["Port"]
			s_["TLS"] = s["TLS"]
			self.Services[s["Name"].(string)] = s_
			self.saveConfig()
		}

	} else if s["Protocol"].(string) == "http" {
		// any other http server except this one...
		if !strings.HasPrefix(s["Name"].(string), "Globular") {
			self.services[s["Name"].(string)] = s

			hasChange := self.saveServiceConfig(s)

			// Kill previous instance of the program.
			if hasChange || s["Process"] == nil {
				self.startService(s)
			}
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
	// Each service contain a file name config.json that describe service.
	// I will keep services info in services map and also it running process.
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == "config.json" {
			// println(path, info.Name())
			// So here I will read the content of the file.
			s := make(map[string]interface{})
			config, err := ioutil.ReadFile(path)
			log.Println("read config file: ", path)
			if err == nil {
				// Read the config file.
				json.Unmarshal(config, &s)
				if s["Protocol"] != nil {
					path_ := path[:strings.LastIndex(path, string(os.PathSeparator))]
					servicePath := path_ + string(os.PathSeparator) + s["Name"].(string)
					if string(os.PathSeparator) == "\\" {
						servicePath += ".exe" // in case of windows.
					}

					s["servicePath"] = servicePath

					// Now I will start the proxy that will be use by javascript client.
					proxyPath := self.path + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
					if string(os.PathSeparator) == "\\" {
						proxyPath += ".exe" // in case of windows.
					}

					// The proxy path
					s["proxyPath"] = proxyPath

					self.initService(s)
				}
			}
		} else if strings.HasSuffix(info.Name(), ".proto") {
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

	// Those information are needed by interceptors so I will wrote in
	// in a temporary file.
	persistenceAddress := self.services["persistence_server"].(map[string]interface{})["Domain"].(string)
	persistenceAddress = persistenceAddress + ":" + Utility.ToString(self.services["persistence_server"].(map[string]interface{})["Port"].(float64))
	certAuthorityTrust := self.creds + string(os.PathSeparator) + "ca.crt"
	certFile := self.creds + string(os.PathSeparator) + "server.crt"
	keyFile := self.creds + string(os.PathSeparator) + "server.pem"

	// I will wrote the info inside a stucture.
	infos := map[string]string{"address": persistenceAddress, "certAuthorityTrust": certAuthorityTrust, "certFile": certFile, "keyFile": keyFile, "pwd": self.RootPassword}

	infosStr, _ := Utility.ToJson(infos)
	err := ioutil.WriteFile(os.TempDir()+string(os.PathSeparator)+"globular_sa", []byte(infosStr), 0644)
	if err != nil {
		log.Panicln(err)
	}
}

// Method must be register in order to be assign to role.
func (self *Globule) registerMethods() error {
	// Here I will create the sa role if it dosen't exist.
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		return errors.New("No persistence service are available to store ressource information.")
	}

	// Cast-it to the persistence client.
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)
	err := p.Connect("local_ressource", self.RootPassword)

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
	} else if strings.HasSuffix(name, "config.json") {
		b, err := ioutil.ReadAll(f) // b has type []byte
		if err != nil {
			log.Fatal(err)
		}
		// set the global variable here.
		code = "window.globularConfig = " + string(b)
		hasChange = true
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
	str, err := Utility.ToJson(self)
	if err == nil {
		ioutil.WriteFile(self.webRoot+string(os.PathSeparator)+"config.json", []byte(str), 0644)
	}
}

/**
 * Init client side connection to service.
 */
func (self *Globule) initClient(name string, token string) {
	log.Println("try to connecto to ", name)
	if self.Services[name] == nil {
		return
	}

	// Set the parameters.
	domain := self.Services[name].(map[string]interface{})["Domain"].(string)
	address := self.Services[name].(map[string]interface{})["Address"].(string)

	// ?? does address must contain the domain instead of the ipv4 part...
	address = address + ":" + Utility.ToString(int(self.Services[name].(map[string]interface{})["Port"].(float64)))

	hasTLS := self.Services[name].(map[string]interface{})["TLS"].(bool)

	name = strings.Split(name, "_")[0]
	fct := "New" + strings.ToUpper(name[0:1]) + name[1:] + "_Client"

	// Set the files.
	keyFile := self.creds + string(os.PathSeparator) + "client.crt"
	certFile := self.creds + string(os.PathSeparator) + "client.key"
	caFile := self.creds + string(os.PathSeparator) + "ca.crt"

	results, err := Utility.CallFunction(fct, domain, address, hasTLS, certFile, keyFile, caFile, token)
	if err == nil {
		if self.clients[name+"_service"] != nil {
			log.Println("--> close actual client connection ", name)
			self.clients[name+"_service"].Close()
		}
		self.clients[name+"_service"] = results[0].Interface().(api.Client)
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
	// The name of the contructor must follow the same pattern.
	Utility.RegisterFunction("NewEcho_Client", echo_client.NewEcho_Client)
	Utility.RegisterFunction("NewSql_Client", sql_client.NewSql_Client)
	Utility.RegisterFunction("NewFile_Client", file_client.NewFile_Client)
	Utility.RegisterFunction("NewPersistence_Client", persistence_client.NewPersistence_Client)
	Utility.RegisterFunction("NewSmtp_Client", smtp_client.NewSmtp_Client)
	Utility.RegisterFunction("NewLdap_Client", ldap_client.NewLdap_Client)
	Utility.RegisterFunction("NewStorage_Client", storage_client.NewStorage_Client)
	Utility.RegisterFunction("NewEvent_Client", event_client.NewEvent_Client)
	Utility.RegisterFunction("NewCatalog_Client", catalog_client.NewCatalog_Client)
	Utility.RegisterFunction("NewMonitoring_Client", monitoring_client.NewMonitoring_Client)

	// That service is program in c++
	Utility.RegisterFunction("NewSpc_Client", spc_client.NewSpc_Client)
	Utility.RegisterFunction("NewPlc_Client", plc_client.NewPlc_Client)
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

////////////////////// HTTPS Server Key's /////////////////////////////////////

// To generate a self-signed certificate (in our case, without encryption):
// -Create a new 4096bit RSA key and save it to server_key.pem, without DES encryption (-newkey, -keyout and -nodes)
// -Create a Certificate Signing Request for a given subject, valid for 365 days (-days, -subj)
// -Sign the CSR using the server key, and save it to server_cert.pem as an X.509 certificate (-x509, -out)
func (self *Globule) GenerateSelfSignedCertificate(path string, expiration_delay int, domain string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "server_key.pem") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-x509")
	args = append(args, "-newkey")
	args = append(args, "rsa:4096")
	args = append(args, "-keyout")
	args = append(args, path+string(os.PathSeparator)+"server_key.pem")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server_cert.pem")
	args = append(args, "-nodes")
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-subj")
	args = append(args, "/C=CA/ST=Montreal/O=Globular Application Server/CN="+domain)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}

	return nil
}

// To create a key and a Certificate Signing Request for client...
func (self *Globule) GenerateClientCertificate(path string, expiration_delay int, clientId string) error {
	// openssl req -newkey rsa:4096 -keyout alice_key.pem -out alice_csr.pem -nodes -days 365 -subj "/CN=Alice"
	if Utility.Exists(path + string(os.PathSeparator) + clientId + "_key.pem") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-newkey")
	args = append(args, "rsa:4096")
	args = append(args, "-keyout")
	args = append(args, path+string(os.PathSeparator)+clientId+"_key.pem")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+clientId+"_csr.pem")
	args = append(args, "-nodes")
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-subj")
	args = append(args, "/CN="+clientId)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}

	return nil
}

// Here, we act as a Certificate Authority, so we supply our certificate and key via the -CA parameters:
func (self *Globule) SignedClientCertificate(path string, expiration_delay int, clientId string) error {
	// openssl x509 -req -in alice_csr.pem -CA server_cert.pem -CAkey server_key.pem -out alice_cert.pem -set_serial 01 -days 365
	if Utility.Exists(path + string(os.PathSeparator) + clientId + "_cert.pem") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+clientId+"_csr.pem")
	args = append(args, "-CA")
	args = append(args, path+string(os.PathSeparator)+"server_cert.pem")
	args = append(args, "-CAkey")
	args = append(args, path+string(os.PathSeparator)+"server_key.pem")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+clientId+"_cert.pem")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}
	os.Remove(path + string(os.PathSeparator) + clientId + "_csr.pem")
	return nil
}

// To use these certificates in our browser, we need to bundle them in PKCS#12
// format. That will contain both the private key and the certificate, thus the
// browser can use it for encryption. For Alice, we add the -clcerts option, which
// excludes the CA certificate from the bundle. Since we issued the certificate,
// we already have the certificate: we don’t need to include it in Alice’s
// certificate as well. You can also password-protect the certificate.
func (self *Globule) BundleClientCertificate(path string, expiration_delay int, clientId string, pwd string) error {
	// openssl x509 -req -in globular_csr.pem -CA server_cert.pem -CAkey server_key.pem -out globular_cert.pem -set_serial 01 -days 365
	if Utility.Exists(path + string(os.PathSeparator) + clientId + ".p12") {
		return nil
	}
	self.GenerateClientCertificate(self.certs, expiration_delay, "globular")
	self.SignedClientCertificate(self.certs, expiration_delay, "globular")
	// openssl pkcs12 -export -clcerts -in globular_cert.pem -inkey globular_key.pem -out globular.p12
	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "pkcs12")
	args = append(args, "-export")
	args = append(args, "-clcerts")
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+clientId+"_cert.pem")
	args = append(args, "-inkey")
	args = append(args, path+string(os.PathSeparator)+clientId+"_key.pem")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+clientId+".p12")
	args = append(args, "-password")
	args = append(args, "pass:"+pwd)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}

	// Here I will remove intermediate file.
	os.Remove(path + string(os.PathSeparator) + clientId + "_cert.pem")
	os.Remove(path + string(os.PathSeparator) + clientId + "_key.pem")

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

	// Now https certificates.
	self.GenerateSelfSignedCertificate(self.certs, expiration_delay, self.Domain)

	// Generate the client certificate.. this is a simple test.

	self.BundleClientCertificate(self.certs, expiration_delay, "globular", "1234")
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
	config["Protocol"] = self.Protocol
	config["IP"] = self.IP
	config["Domain"] = self.Domain
	config["ReadTimeout"] = self.ReadTimeout
	config["WriteTimeout"] = self.WriteTimeout
	config["IdleTimeout"] = self.IdleTimeout
	config["SessionTimeout"] = self.SessionTimeout
	config["CertExpirationDelay"] = self.CertExpirationDelay

	// return the full service configuration.
	config["Services"] = self.services

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
	config["Protocol"] = self.Protocol
	config["IP"] = self.IP
	config["Domain"] = self.Domain
	config["ReadTimeout"] = self.ReadTimeout
	config["WriteTimeout"] = self.WriteTimeout
	config["IdleTimeout"] = self.IdleTimeout
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

// return true if the configuation has change.
func (self *Globule) saveServiceConfig(config map[string]interface{}) bool {

	// get the config path.
	servicePath := config["servicePath"].(string)
	proxyPath := config["proxyPath"].(string)
	process := config["Process"]
	proxyProcess := config["ProxyProcess"]
	path := config["servicePath"].(string)
	path = path[:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator) + "config.json"

	// remove unused information...
	delete(config, "Process")
	delete(config, "ProxyProcess")

	// remove this info from the file to be save.
	delete(config, "servicePath")
	delete(config, "proxyPath")

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
				config["proxyPath"] = proxyPath
				config["Process"] = process
				config["ProxyProcess"] = proxyProcess
				return false
			}
		}
	}
	f.Close()

	// The new config to write.
	jsonStr, _ := Utility.ToJson(config)

	// here I will write the file.
	err = ioutil.WriteFile(path, []byte(jsonStr), 0644)
	if err != nil {
		log.Println("fail to save config file: ", err)
	}

	// set back internal infos...
	config["servicePath"] = servicePath
	config["proxyPath"] = proxyPath
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
	if config["Name"] != nil {
		srv := self.services[config["Name"].(string)]
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
			self.Protocol = config["Protocol"].(string)
			self.IP = config["IP"].(string)
			self.Domain = config["Domain"].(string)
			self.ReadTimeout = time.Duration(Utility.ToInt(config["ReadTimeout"].(float64)))
			self.WriteTimeout = time.Duration(Utility.ToInt(config["WriteTimeout"].(float64)))
			self.IdleTimeout = time.Duration(Utility.ToInt(config["IdleTimeout"].(float64)))
			self.CertExpirationDelay = Utility.ToInt(config["CertExpirationDelay"].(float64))
			// That will save the services if they have changed.
			for n, s := range config["Services"].(map[string]interface{}) {
				// Attach the actual process and proxy process to the configuration object.
				s.(map[string]interface{})["Process"] = self.services[n].(map[string]interface{})["Process"]
				s.(map[string]interface{})["ProxyProcess"] = self.services[n].(map[string]interface{})["ProxyProcess"]
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
	s := self.services[rqst.ServiceId]
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

	s := self.services[rqst.ServiceId]
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
func (self *Globule) startExternalService(serviceId string) (int, error) {

	if service, ok := self.ExternalServices[serviceId]; !ok {
		return -1, errors.New("No external service found with name " + serviceId)
	} else {

		service.srv = exec.Command(service.Path, service.Args...)

		err := service.srv.Start()
		if err != nil {
			return -1, err
		}

		// save back the service in the map.
		self.ExternalServices[serviceId] = service

		return service.srv.Process.Pid, nil
	}

}

// Stop external service.
func (self *Globule) stopExternalService(serviceId string) error {
	if _, ok := self.ExternalServices[serviceId]; !ok {
		return errors.New("No external service found with name " + serviceId)
	}

	// if no command was created
	if self.ExternalServices[serviceId].srv == nil {
		return nil
	}

	// if no process running
	if self.ExternalServices[serviceId].srv.Process == nil {
		return nil
	}

	// kill the process.
	return self.ExternalServices[serviceId].srv.Process.Kill()
}

// Register external service to be start by Globular in order to run
// as exemple MongoDB and Prometheus.
func (self *Globule) RegisterExternalService(ctx context.Context, rqst *admin.RegisterExternalServiceRequest) (*admin.RegisterExternalServiceResponse, error) {

	// Here I will get the command path.
	externalCmd := ExternalService{
		Id:   rqst.ServiceId,
		Path: rqst.Path,
		Args: rqst.Args,
	}

	self.ExternalServices[externalCmd.Id] = externalCmd

	// save the config.
	self.saveConfig()

	// start the external service.
	pid, err := self.startExternalService(externalCmd.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &admin.RegisterExternalServiceResponse{
		ServicePid: int64(pid),
	}, nil
}

/**
 * Start internal service admin and ressource are use that function.
 */
func (self *Globule) startInternalService(name string, port int, proxy int) (*grpc.Server, error) {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, name+" service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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
	s["Address"] = self.IP
	s["Port"] = port
	s["Proxy"] = proxy
	s["Domain"] = self.Domain
	s["TLS"] = self.Protocol == "https"
	self.services[name] = s

	// start the proxy
	err := self.startProxy(name, port, proxy)
	if err != nil {
		return nil, err
	}

	return grpcServer, nil
}

/** Create the super administrator in the db. **/
func (self *Globule) registerSa() error {
	// Get a reference to the persistence service.
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		return errors.New("No persistence service was configure for your globular server.")
	}
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)

	// Test if the connection already exist.
	err := p.Ping("local_ressource")

	// if not I will create one.
	if err != nil {
		err = p.CreateConnection("local_ressource", "local_ressource", "localhost", 27017, 0, "sa", self.RootPassword, 5000, "", true)
		if err != nil {
			return err
		}

		// Now I will create a new user name sa and give it all admin write.
		createSaScript := fmt.Sprintf(
			`db=db.getSiblingDB('admin');db.createUser({ user: '%s', pwd: '%s', roles: ['userAdminAnyDatabase','userAdmin','readWrite','dbAdmin','clusterAdmin','readWriteAnyDatabase','dbAdminAnyDatabase']})`, "sa", self.RootPassword) // must be change...
		log.Println("---> create super admin: ", createSaScript)

		// I will execute the sript with the admin function.
		err = p.RunAdminCmd("local_ressource", "", "", createSaScript)
		if err != nil {
			return err
		}
	}
	return nil

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
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))

	}

	// Cast-it to the persistence client.
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)

	// Test if the connection already exist.
	err := p.Ping("local_ressource")

	// if not I will create one.
	if err != nil {
		// Connect to the database here.
		err = p.CreateConnection("local_ressource", "local_ressource", "localhost", 27017, 0, "sa", self.RootPassword, 5000, "", true)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
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

	// That service made user of persistence service.
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))
	}

	// Cast-it to the persistence client.
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)
	err := p.Connect("local_ressource", self.RootPassword)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
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

	// Here I got the token I will now put it in the cache.
	return &ressource.AuthenticateRsp{
		Token: tokenString,
	}, nil
}

//* Refresh token get a new token *
func (self *Globule) RefreshToken(ctx context.Context, rqst *ressource.RefreshTokenRqst) (*ressource.RefreshTokenRsp, error) {

	// first of all I will validate the current token.
	name, err := Interceptors_.ValidateToken(rqst.Token)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	tokenString, err := Interceptors.GenerateToken(self.jwtKey, self.SessionTimeout, name)
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
	var p *persistence_client.Persistence_Client
	if self.clients["persistence_service"] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No persistence service are available to store ressource information.")))

	}

	// Cast-it to the persistence client.
	p = self.clients["persistence_service"].(*persistence_client.Persistence_Client)
	err := p.Connect("local_ressource", self.RootPassword)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
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

/**
 * That function is call when it's time to validate the certificate.
 */
func (self *Globule) verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(self.certs + string(os.PathSeparator) + "server_cert.pem")
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", self.certs+string(os.PathSeparator)+"server_cert.pem", err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	opts := x509.VerifyOptions{
		Roots: rootCAs,
	}

	rawCert := rawCerts[0]
	cert, err := x509.ParseCertificate(rawCert)

	// The CommonName contain the name of the user name for who the certificate
	// was generated. So here It's the place to test if the user has the write
	// to access the service.
	log.Println("----->", cert.Subject.CommonName)

	if err != nil {
		return err
	}
	_, err = cert.Verify(opts)

	return err
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// First of all I will start external services.
	for externalServiceId, _ := range self.ExternalServices {
		pid, err := self.startExternalService(externalServiceId)
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

	// Here I will make a signal hook to interrupt to exit cleanly.
	// handle the Interrupt

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

		for key, value := range self.services {
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
		for externalServiceId, _ := range self.ExternalServices {
			self.stopExternalService(externalServiceId)

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
		}()
	}

	ressource_server, err := self.startInternalService("Ressource", self.RessourcePort, self.RessourceProxy)
	if err == nil {
		// First of all I will creat a listener.
		// Create the channel to listen on admin port.

		// Create the channel to listen on admin port.
		lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(self.RessourcePort))
		if err != nil {
			log.Fatalf("could not start admin service %s: %s", self.Domain, err)
		}

		ressource.RegisterRessourceServiceServer(ressource_server, self)

		// Here I will make a signal hook to interrupt to exit cleanly.
		go func() {
			log.Println("---> start ressource service!")
			go func() {
				// Get the list of all services method.
				self.registerMethods()
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

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

	// Start the http server.
	// Start http server.
	go func() {
		log.Println("Globular is listening at port ", self.PortHttp)
		log.Panicln(http.ListenAndServe(":"+strconv.Itoa(self.PortHttp), nil))
	}()

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(self.certs + string(os.PathSeparator) + "server_cert.pem")
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", self.certs+string(os.PathSeparator)+"server_cert.pem", err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// client certificate
	server := &http.Server{
		Addr: ":" + Utility.ToString(self.PortHttps),
		TLSConfig: &tls.Config{
			RootCAs:            rootCAs,
			ClientAuth:         tls.RequestClientCert,
			InsecureSkipVerify: false,
		},
	}

	log.Panicln(server.ListenAndServeTLS(self.certs+string(os.PathSeparator)+"server_cert.pem", self.certs+string(os.PathSeparator)+"server_key.pem"))

	// Use public CA for public website with real domain name.
	// err = http.ListenAndServeTLS(":"+Utility.ToString(self.PortHttps), self.path+"/sslforfree/certificate.crt", self.path+"/sslforfree/private.key", r)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
