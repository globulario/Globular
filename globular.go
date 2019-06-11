package main

import (
	"bufio"
	//	"context"
	//	"crypto/tls"
	"encoding/json"
	"errors"

	//"fmt"
	"io"
	"io/ioutil"
	"log"
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

	// Client services.
	"github.com/davecourtois/Globular/echo/echo_client"
	"github.com/davecourtois/Globular/file/file_client"
	"github.com/davecourtois/Globular/ldap/ldap_client"
	"github.com/davecourtois/Globular/persistence/persistence_client"
	"github.com/davecourtois/Globular/smtp/smtp_client"
	"github.com/davecourtois/Globular/sql/sql_client"
	"github.com/davecourtois/Globular/storage/storage_client"
	// "github.com/davecourtois/Globular/spc/spc_client"
)

var (
	root    string
	globule *Globule
)

/**
 * The web server.
 */
type Globule struct {
	// The share part of the service.
	Name                string // The service name
	PortHttp            int    // The port of the http file server.
	PortHttps           int    // The secure port
	Protocol            string // The protocol of the service.
	IP                  string // The local address...
	Services            map[string]interface{}
	Domain              string // The domain name of your application
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	CertExpirationDelay int

	// Local info.
	webRoot string // The root of the http file server.
	creds   string // where the key will be store.
	path    string // The path of the exec...

	// The list of avalaible services.
	services map[string]interface{}

	// The map of client...
	clients map[string]api.Client
}

/**
 * Globule constructor.
 */
func NewGlobule(port int) *Globule {
	// Here I will initialyse configuration.
	g := new(Globule)
	g.PortHttp = port
	g.PortHttps = port // The default port number.
	g.Name = Utility.GetExecName(os.Args[0])
	g.Protocol = "http"
	g.Domain = "localhost"
	g.IP = Utility.MyIP()

	// set default values.
	g.IdleTimeout = 120
	g.ReadTimeout = 5
	g.WriteTimeout = 5
	g.CertExpirationDelay = 365

	// Set the service map.
	g.services = make(map[string]interface{}, 0)

	// Set the share service info...
	g.Services = make(map[string]interface{}, 0)

	// Set the map of client.
	g.clients = make(map[string]api.Client, 0)

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

	// keep the root in global variable for the file handler.
	root = g.webRoot
	globule = g

	return g
}

/**
 * Here I will set services
 */
func (self *Globule) initServices() {
	log.Println("Initialyse services")

	// If the protocol is https I will generate the TLS certificate.
	if self.Protocol == "https" {
		// TODO find a way to save the password somewhere on the server configuration

		// whitout expose it to external world.
		self.GenerateServicesCertificates("1111", self.CertExpirationDelay)
	}

	// Each service contain a file name config.json that describe service.
	// I will keep services info in services map and also it running process.
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == "config.json" {
			// println(path, info.Name())
			// So here I will read the content of the file.
			s := make(map[string]interface{})
			config, err := ioutil.ReadFile(path)
			if err == nil {
				// Read the config file.
				json.Unmarshal(config, &s)

				if s["Protocol"].(string) == "grpc" {

					path_ := path[:strings.LastIndex(path, string(os.PathSeparator))]
					servicePath := path_ + string(os.PathSeparator) + s["Name"].(string)
					if string(os.PathSeparator) == "\\" {
						servicePath += ".exe" // in case of windows.
					}

					// Now I will start the proxy that will be use by javascript client.
					proxyPath := self.path + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
					if string(os.PathSeparator) == "\\" {
						proxyPath += ".exe" // in case of windows.
					}

					// The domain must be set in the sever configuration and not change after that.
					s["Domain"] = self.Domain // local services.
					s["Address"] = self.IP    // local services.
					proxyBackendAddress := s["Domain"].(string) + ":" + Utility.ToString(s["Port"])

					proxyAllowAllOrgins := Utility.ToString(s["AllowAllOrigins"])
					proxyArgs := make([]string, 0)

					// Use in a local network or in test.
					proxyArgs = append(proxyArgs, "--backend_addr="+proxyBackendAddress)
					proxyArgs = append(proxyArgs, "--allow_all_origins="+proxyAllowAllOrgins)

					if self.Protocol == "https" {

						// Set TLS local services configuration here.
						s["TLS"] = true
						s["CertAuthorityTrust"] = self.creds + string(os.PathSeparator) + "ca.crt"
						s["CertFile"] = self.creds + string(os.PathSeparator) + "server.crt"
						s["KeyFile"] = self.creds + string(os.PathSeparator) + "server.pem"

						// Now I will save the file with those new information in it.
						jsonStr, _ := Utility.ToJson(&s)
						ioutil.WriteFile(path, []byte(jsonStr), 0644)

						// Set local client configuration here.

						// Now set the proxy information here.

						/* Services gRpc backend. */
						proxyArgs = append(proxyArgs, "--backend_tls=true")
						proxyArgs = append(proxyArgs, "--backend_tls_ca_files="+self.creds+string(os.PathSeparator)+"ca.crt")
						proxyArgs = append(proxyArgs, "--backend_client_tls_cert_file="+self.creds+string(os.PathSeparator)+"client.crt")
						proxyArgs = append(proxyArgs, "--backend_client_tls_key_file="+self.creds+string(os.PathSeparator)+"client.pem")

						/* http2 parameters between the browser and the proxy.*/
						proxyArgs = append(proxyArgs, "--run_http_server=false")
						proxyArgs = append(proxyArgs, "--run_tls_server=true")
						proxyArgs = append(proxyArgs, "--server_http_tls_port="+Utility.ToString(s["Proxy"]))

						proxyArgs = append(proxyArgs, "--server_tls_client_ca_files="+self.path+"/sslforfree/ca_bundle.crt")
						proxyArgs = append(proxyArgs, "--server_tls_cert_file="+self.path+"/sslforfree/certificate.crt")
						proxyArgs = append(proxyArgs, "--server_tls_key_file="+self.path+"/sslforfree/private.key")

					} else {
						// not secure services.
						s["TLS"] = false
						s["CertAuthorityTrust"] = ""
						s["CertFile"] = ""
						s["KeyFile"] = ""

						// Now I will save the file with those new information in it.
						jsonStr, _ := Utility.ToJson(&s)
						ioutil.WriteFile(path, []byte(jsonStr), 0644)
						proxyArgs = append(proxyArgs, "--run_http_server=true")
						proxyArgs = append(proxyArgs, "--run_tls_server=false")
						proxyArgs = append(proxyArgs, "--server_http_debug_port="+Utility.ToString(s["Proxy"]))
						proxyArgs = append(proxyArgs, "--backend_tls=false")
					}

					// log.Println(proxyPath, proxyArgs)
					// Start the service process.
					log.Println("try to start process ", s["Name"].(string))
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
					}

					// start the proxy service.
					s["ProxyProcess"] = exec.Command(proxyPath, proxyArgs...)

					err = s["ProxyProcess"].(*exec.Cmd).Start()
					if err != nil {
						log.Println("Fail to start grpcwebproxy: ", s["Name"].(string), " at port ", s["Proxy"], " with error ", err)
					}

					self.services[s["Name"].(string)] = s
					s_ := make(map[string]interface{})

					// export public service values.
					s_["Address"] = s["Address"]
					s_["Domain"] = s["Domain"]
					s_["Proxy"] = s["Proxy"]
					s_["Port"] = s["Port"]

					self.Services[s["Name"].(string)] = s_
					self.saveConfig()

					log.Println("Service ", s["Name"].(string), "is running at port", s["Port"], "it's proxy port is", s["Proxy"])
				}
			}
		}
		return nil
	})
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
 * Here here is where services function are call from http.
 */

/**
* That function handle http query as form of what so called API.
* exemple of use.

  Get all entity prototype from CargoEntities
  ** note the access_token can change over time.
  http://mon176:10000/api/Server/EntityManager/GetEntityPrototypes?storeId=CargoEntities&access_token=C4X_UsRXRCqwqsWfuEdgFA

  Get an entity object with a given uuid.
  * Note because % is in the uuid string it must be escape with %25 so here
  	 the uuid is CargoEntities.Action%7facc2a5-dcb7-4ae7-925a-fb0776a9da00
  http://localhost:10000/api/Server/EntityManager/GetObjectByUuid?p0=CargoEntities.Action%257facc2a5-dcb7-4ae7-925a-fb0776a9da00
*/
func HttpQueryHandler(w http.ResponseWriter, r *http.Request) {

	// So the request will contain...
	// The last tow parameters must be empty because we don't use the websocket
	// here.

	inputs := strings.Split(r.URL.Path[len("/api/"):], "/")

	if len(inputs) < 2 {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("api call error, not enought arguments given!"))
		return
	}

	// Get the client connected to the required service.
	service := globule.clients[inputs[0]]

	if service == nil {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("service " + inputs[0] + " not found"))
		return
	}

	// The parameter values.
	params := make([]interface{}, 0)
	for i := 0; i < len(r.URL.Query()); i++ {
		if r.URL.Query()["p"+strconv.Itoa(i)] != nil {
			params = append(params, r.URL.Query()["p"+strconv.Itoa(i)][0])
		} else {
			w.Header().Set("Content-Type", "application/text")
			w.Write([]byte("p" + strconv.Itoa(i) + " not found!"))
			return
		}
	}

	// Here I will call the function on the service.
	var err_ interface{}
	var results interface{}
	log.Println("call api function: ", inputs[1], params)
	results, err_ = Utility.CallMethod(service, inputs[1], params)
	if err_ != nil {

		w.Header().Set("Content-Type", "application/text")
		switch v := err_.(type) {
		case error:
			w.Write([]byte(v.Error()))
		case string:
			w.Write([]byte(v))
		}

		return
	}

	// Here I will get the res
	var resultStr []byte
	var err error
	resultStr, err = json.Marshal(results)
	if err != nil {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte(err.(error).Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resultStr, _ = Utility.PrettyPrint(resultStr)
	w.Write(resultStr)

}

/**
 * This code is use to upload a file into the tmp directory of the server
 * via http request.
 */
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("FileUploadHandler")

	// I will
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		log.Println(w, err)
		return
	}

	log.Println("FileUploadHandler", 425)
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
func (self *Globule) initClient(name string) {
	log.Println("connecto to service ", name)
	port := int(self.Services[name+"_server"].(map[string]interface{})["Port"].(float64))
	fct := "New" + strings.ToUpper(name[0:1]) + name[1:] + "_Client"

	// Set the parameters.
	address := self.Domain + ":" + strconv.Itoa(port)
	domain := self.Domain
	hasTLS := self.Protocol == "https" // true if the protocol is https.
	keyFile := self.creds + string(os.PathSeparator) + "client.crt"
	certFile := self.creds + string(os.PathSeparator) + "client.key"
	caFile := self.creds + string(os.PathSeparator) + "ca.crt"

	results, err := Utility.CallFunction(fct, domain, address, hasTLS, certFile, keyFile, caFile)
	if err == nil {
		self.clients[name+"_service"] = results[0].Interface().(api.Client)
	}
}

/**
 * Init the service client.
 * Keep the service constructor for further call.
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

	// The echo service
	for k, _ := range self.services {
		name := strings.Split(k, "_")[0]
		self.initClient(name)
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

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// set the services.
	self.initServices()

	// set the client services.
	self.initClients()

	r := http.NewServeMux()

	// Start listen for http request.
	r.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	r.HandleFunc("/uploads", FileUploadHandler)

	// Give access to service.
	r.HandleFunc("/api/", HttpQueryHandler)

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
					log.Println("kill service process ", p.(*exec.Cmd).Process.Pid)
					p.(*exec.Cmd).Process.Kill()
				}
			}

			if value.(map[string]interface{})["ProxyProcess"] != nil {
				p := value.(map[string]interface{})["ProxyProcess"]
				if reflect.TypeOf(p).String() == "*exec.Cmd" {
					log.Println("kill proxy process ", p.(*exec.Cmd).Process.Pid)
					p.(*exec.Cmd).Process.Kill()
				}
			}
		}

		for _, value := range self.clients {
			value.Close()
		}

		// exit cleanly
		os.Exit(0)

	}()

	log.Println("Listening...")
	var err error
	if self.Protocol == "http" {
		err = http.ListenAndServe(":"+strconv.Itoa(self.PortHttp), r)
	} else {
		// Here I will use sslforfree certificate to publish the website.
		err = http.ListenAndServeTLS(":443", self.path+"/sslforfree/certificate.crt", self.path+"/sslforfree/private.key", r)
	}

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
