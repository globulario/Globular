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
	"path/filepath"
	"syscall"
	"time"

	"reflect"
	"strconv"

	"path"
	"strings"

	"github.com/davecourtois/Utility"
)

var (
	root string
)

/**
 * The web server.
 */
type Globule struct {
	Name     string          // The service name
	Path     string          // The service path
	Port     int             // The port of the http file server.
	Protocol string          // The protocol of the service.
	WebRoot  string          // The root of the http file server.
	Address  *Utility.IPInfo // contain the ipv4 address.

	// The list of avalaible services.
	Services map[string]interface{}
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
	ip, _ := Utility.MyIP()
	g.Address = ip

	// Set the service map.
	g.Services = make(map[string]interface{}, 0)

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

					self.Services[s["Name"].(string)] = s
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
	str, err := Utility.ToJson(self)
	if err == nil {
		ioutil.WriteFile(self.WebRoot+string(os.PathSeparator)+"config.json", []byte(str), 0644)
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

	r := http.NewServeMux()

	// Start listen for http request.
	r.HandleFunc("/", ServeFileHandler)

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

		for key, value := range self.Services {
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

		// exit cleanly
		os.Exit(0)

	}()

	log.Println("Listening...")
	err := http.ListenAndServe(":"+strconv.Itoa(self.Port), r)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
