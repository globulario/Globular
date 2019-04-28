package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	//	"reflect"
	"strconv"

	"github.com/davecourtois/Utility"
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

	// The list of avalaible services.
	services map[string]interface{}
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {
	// Here I will initialyse configuration.
	g := new(Globule)
	g.Port = 8080       // The default port number.
	g.Path = os.Args[0] // the serive name
	g.Name = "Globule"
	g.Protocol = "http"

	// Set the service map.
	g.services = make(map[string]interface{}, 0)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		g.WebRoot = dir + "/WebRoot"           // The default directory to server.
		Utility.CreateDirIfNotExist(g.WebRoot) // Create the directory if it not exist.
		file, err := ioutil.ReadFile(g.WebRoot + "/config.json")

		if err == nil {
			json.Unmarshal([]byte(file), g)
		}
	}

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
				json.Unmarshal(config, &s)
				if s["Protocol"].(string) == "grpc" {
					log.Println("start rpc service: ", s["Name"].(string))
					// Start the process.
					s["Process"] = exec.Command(s["Path"].(string), Utility.ToString(s["Port"]))
					err = s["Process"].(*exec.Cmd).Start()
					if err != nil {
						log.Println("Fail to start service: ", s["Name"].(string), " at port ", s["Port"], " with error ", err)
					}
				}
				self.services[s["Name"].(string)] = s
			}
		}
		return nil
	})
}

/**
 * Listen for new connection.
 */
func (self *Globule) Listen() {

	log.Println("Start Globule at port ", self.Port)

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// set the services.
	self.initServices()

	// Start listen for http request.
	fs := http.FileServer(http.Dir(self.WebRoot))
	http.Handle("/", fs)

	// Here I will save the server attribute
	str, err := Utility.ToJson(self)
	if err == nil {
		ioutil.WriteFile(self.WebRoot+"/config.json", []byte(str), 0644)
	}

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println("Listening...")
		http.ListenAndServe(":"+strconv.Itoa(self.Port), nil)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch

	// Here the server stop running,
	// so I will close the services.
	for key, value := range self.services {
		log.Println("Stop service ", key)
		if value.(map[string]interface{})["Process"] != nil {
			value.(map[string]interface{})["Process"].(*exec.Cmd).Process.Kill()
		}
	}
}
