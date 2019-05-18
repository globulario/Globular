package main

import (
	"bufio"
	"encoding/json"
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

	"github.com/davecourtois/Utility"
)

var (
	root    string
	globule *Globule
)

/**
 * The web server.
 */
type Globule struct {
	Name     string // The service name
	Path     string // The service path
	Port     int    // The port of the http file server.
	Protocol string // The protocol of the service.
	WebRoot  string // The root of the http file server.
	IP       string // The local address...

	// The list of avalaible services.
	services map[string]interface{}

	// The share part of the service.
	Services map[string]interface{}

	// The map of client...
	clients map[string]Client
}

/**
 * Globule constructor.
 */
func NewGlobule(port int) *Globule {
	// Here I will initialyse configuration.
	g := new(Globule)
	g.Port = port // The default port number.
	g.Name = Utility.GetExecName(os.Args[0])
	g.Protocol = "http"
	g.IP = Utility.MyIP()

	// Set the service map.
	g.services = make(map[string]interface{}, 0)

	// Set the share service info...
	g.Services = make(map[string]interface{}, 0)

	// Set the map of client.
	g.clients = make(map[string]Client, 0)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	g.Path = dir // keep the installation patn.

	if err == nil {
		g.WebRoot = dir + string(os.PathSeparator) + "WebRoot" // The default directory to server.
		Utility.CreateDirIfNotExist(g.WebRoot)                 // Create the directory if it not exist.
		file, err := ioutil.ReadFile(g.WebRoot + string(os.PathSeparator) + "config.json")
		// Init the servce with the default port address
		if err == nil {
			json.Unmarshal([]byte(file), g)
		}
	}

	// keep the root in global variable for the file handler.
	root = g.WebRoot
	globule = g

	return g
}

/**
 * Here I will set services
 */
func (self *Globule) initServices() {
	log.Println("Initialyse services")

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

					// Start the process.
					log.Println("try to start process ", s["Name"].(string))

					s["Process"] = exec.Command(servicePath, Utility.ToString(s["Port"]))
					err = s["Process"].(*exec.Cmd).Start()
					if err != nil {
						log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
					}

					// Now I will start the proxy that will be use by javascript client.
					proxyPath := self.Path + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + "grpcwebproxy"
					if string(os.PathSeparator) == "\\" {
						proxyPath += ".exe" // in case of windows.
					}

					// This is the grpc service to connect with the proxy
					proxyBackendAddress := "localhost:" + Utility.ToString(s["Port"])
					proxyAllowAllOrgins := Utility.ToString(s["AllowAllOrigins"])

					// start the proxy service.
					s["ProxyProcess"] = exec.Command(proxyPath, "--backend_addr="+proxyBackendAddress, "--server_http_debug_port="+Utility.ToString(s["Proxy"]), "--run_tls_server=false", "--allow_all_origins="+proxyAllowAllOrgins)
					err = s["ProxyProcess"].(*exec.Cmd).Start()
					if err != nil {
						log.Println("Fail to start grpcwebproxy: ", s["Name"].(string), " at port ", s["Proxy"], " with error ", err)
					}

					self.services[s["Name"].(string)] = s

					s_ := make(map[string]interface{})

					// export public service values.
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

	// The array of parameters.
	log.Println("-----> ", r.URL.Query())

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

	log.Println(params)

	/** TODO Implement OAUTH2
	var accessTokenId string

	// Try to get access token from the list of parameters.
	accessTokenId = values.Get("access_token")

	if len(accessTokenId) == 0 {
		values := strings.Split(r.Header.Get("Authorization"), " ")
		log.Println("values: ", values)
		if len(values) == 2 {
			if strings.ToLower(values[0]) == "bearer" {
				accessTokenId = values[1]
			}
		}
	}

	if len(accessTokenId) == 0 {
		if len(r.Form["access_token"]) == 1 {
			accessTokenId = r.Form["access_token"][0]
		}
	}

	// The access token variable.
	var accessToken *Config.OAuth2Access

	if len(accessTokenId) > 0 {
		ids := []interface{}{accessTokenId}
		entity, err := GetServer().GetEntityManager().getEntityById("Config.OAuth2Access", "Config", ids)
		if err != nil {
			w.Header().Set("Content-Type", "application/text")
			w.Write([]byte(err.GetBody()))
			return
		}

		// Get the access token here.
		accessToken = entity.(*Config.OAuth2Access)

	} else {
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Access denied!"))
		return
	}

	log.Println("open id token: ", accessToken.GetUserData())
	*/

	// Here I will call the function on the service.
	var err_ interface{}
	var results interface{}
	results, err_ = Utility.CallMethod(service, inputs[1], params)
	if err_ != nil {
		log.Println(results, err_)
		w.Header().Set("Content-Type", "application/text")
		if reflect.TypeOf(err_).Kind() == reflect.String {
			w.Write([]byte(err_.(string)))
		} else {
			w.Write([]byte(err_.(error).Error()))
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
		path = globule.WebRoot + path
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
		ioutil.WriteFile(self.WebRoot+string(os.PathSeparator)+"config.json", []byte(str), 0644)
	}
}

/**
 * Init client side connection to service.
 */
func (self *Globule) initClient(name string) {
	log.Println("connecto to service ", name)
	port := int(self.Services[name+"_server"].(map[string]interface{})["Port"].(float64))
	results, err := Utility.CallFunction("New"+strings.ToUpper(name[0:1])+name[1:]+"_Client", "localhost:"+strconv.Itoa(port))
	if err == nil {
		self.clients[name+"_service"] = results[0].Interface().(Client)
	}
}

/**
 * Init the service client.
 */
func (self *Globule) initClients() {
	// Register service constructor function here.
	// The name of the contructor must follow the same pattern.
	Utility.RegisterFunction("NewEcho_Client", NewEcho_Client)
	Utility.RegisterFunction("NewSql_Client", NewSql_Client)
	Utility.RegisterFunction("NewFile_Client", NewFile_Client)
	Utility.RegisterFunction("NewPersistence_Client", NewPersistence_Client)
	Utility.RegisterFunction("NewSmtp_Client", NewSmtp_Client)
	Utility.RegisterFunction("NewLdap_Client", NewLdap_Client)

	// The echo service
	for k, _ := range self.services {
		name := strings.Split(k, "_")[0]
		self.initClient(name)
	}
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

	log.Println("Start Globular at port ", self.Port)

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
	err := http.ListenAndServe(":"+strconv.Itoa(self.Port), r)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
