package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/signal"

	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Globular/plc/plc_client"
	"github.com/davecourtois/Utility"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultPort = int(2112)
)

type Tag struct {
	ServiceId    string // The name of the service itself ex. plc_server_ab or plc_server_siemens
	ConnectionId string // Must be already define within the plc_server.
	Domain       string // The domain where the plc server that contain the Tag run.
	Name         string // The name of the tag
	Label        string // The label is can not contain . []
	Description  string // The tag description
	Unit         string // The tag unit
	TypeName     string // Can be BOOL, SINT, INT, DINT, REAL.
	Offset       int32  // The offset where to begin to read the s.Tags[i].
	Length       int32  // The Length of the array of tag's.
	Unsigned     bool   // If true the tag is read as unsigned value
}

type server struct {

	// The global attribute of the services.
	Id              string
	Name            string
	Port            int
	Proxy           int
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.
	Protocol        string
	Domain          string
	PublisherId     string

	// self-signed X.509 public keys for distribution
	CertFile string

	// a private RSA key to sign and authenticate the public key
	KeyFile string

	// a private RSA key to sign and authenticate the public key
	CertAuthorityTrust string

	TLS          bool
	Version      string
	KeepUpToDate bool
	KeepAlive    bool

	// The list of tag to monitor.
	Tags         []Tag
	Timeout      int
	NbConnection int

	// contain pointer to connection.
	actions chan map[string]interface{}
	done    chan bool
}

// Read plc tag values.
func readPlcTag(c *plc_client.Plc_Client, connectionId string, tag Tag) ([]interface{}, error) {
	var err error
	var val []interface{}

	if tag.TypeName == "BOOL" {
		val, err = c.ReadTag(connectionId, tag.Name, 0.0, tag.Offset, tag.Length, tag.Unsigned)
	} else if tag.TypeName == "SINT" {
		val, err = c.ReadTag(connectionId, tag.Name, 1.0, tag.Offset, tag.Length, tag.Unsigned)
	} else if tag.TypeName == "INT" {
		val, err = c.ReadTag(connectionId, tag.Name, 2.0, tag.Offset, tag.Length, tag.Unsigned)
	} else if tag.TypeName == "DINT" {
		val, err = c.ReadTag(connectionId, tag.Name, 3.0, tag.Offset, tag.Length, tag.Unsigned)
	} else if tag.TypeName == "REAL" {
		val, err = c.ReadTag(connectionId, tag.Name, 4.0, tag.Offset, tag.Length, tag.Unsigned)
	} else {
		return nil, errors.New(tag.TypeName + " is not a valid data type!")
	}
	return val, err
}

// Process values here.
func (self *server) run() {

	self.actions = make(chan map[string]interface{})
	self.done = make(chan bool)

	// Here I will iterate over the tag's and open connection
	// with it plc servers.
	// init container
	tagsTimeSerie := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "plc_tags",
		Help: "Contain the values of various plc tags.",
	},
		[]string{
			"name",
			"label",
			"connectionId",
			"unit",
			"offset",
			"index"},
	)

	prometheus.MustRegister(tagsTimeSerie)

	// Here I will keep connection pool...
	clients := make(map[string]*plc_client.Plc_Client)

	// Contain list of connection that can be use to communicate with the plc.
	done := make(chan bool)

	// tag reading loop.
	for {
		select {

		// we pass here at each timeout...
		case action := <-self.actions:
			tag := action["tag"].(Tag)

			// create a new client if none exist.
			if clients[tag.ServiceId] == nil {
				clients[tag.ServiceId], _ = plc_client.NewPlc_Client(tag.Domain, tag.ServiceId)
			}

			// read the tag value
			values, err := readPlcTag(clients[tag.ServiceId], tag.ConnectionId, tag)

			// diplay error in that case
			if err != nil {
				if strings.HasSuffix(err.Error(), "-32") {
					err = errors.New("Tag reading error " + tag.Name + " at address " + tag.ServiceId + " " + err.Error())
				} else {
					// here I will close the client and reconnect latter.
					clients[tag.ServiceId].Close()
					delete(clients, tag.ServiceId)
				}

				fmt.Println(err)
				time.Sleep(1000)
			} else {
				for i := 0; i < len(values); i++ {
					tagsTimeSerie.WithLabelValues(tag.Name, tag.Label, tag.ConnectionId, tag.Unit, Utility.ToString(tag.Offset), Utility.ToString(i)).Set(Utility.ToNumeric(values[i]))
				}
			}

		case <-self.done:
			// resume the channel...
			done <- true
			<-done
			self.done <- true
			break
		}
	}

}

func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")

	self.Timeout = 500
	self.NbConnection = 5
	self.Tags = make([]Tag, 0)

	if err == nil {
		json.Unmarshal([]byte(file), self)
		// set the default address and domain
		if self.Domain == "" {
			self.Domain = "localhost"
		}
	} else {
		if len(self.Id) == 0 {
			// Generate random id for the server instance.
			self.Id = Utility.RandomUUID()
		}
		self.save()
	}

}

func (self *server) save() error {
	// Create the file...
	str, err := Utility.ToJson(self)
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	ioutil.WriteFile(dir+"/config.json", []byte(str), 0644)
	return nil
}

func main() {
	// The first argument must be the port number to listen to.
	port := defaultPort
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	s := new(server)
	s.Version = "0.0.1"
	s.Domain = "localhost"

	s.init()

	// read tags and keep values in tagsTimeSerie

	// Start the processing loop.
	go func() {
		s.run()
	}()

	done := make(chan bool) // stop tick

	// use ticker to read connection at each interval of time
	go func() {
		fmt.Println("Start plc_exporter...")
		// read tag at each intervals.
		ticker := time.NewTicker(time.Millisecond * time.Duration(s.Timeout))
		go func() {
			for {
				select {
				// we pass here at each timeout...
				case <-ticker.C:
					// fmt.Println("---> read tags...")
					// Read all register values.
					for i := 0; i < len(s.Tags); i++ {
						action := make(map[string]interface{})
						action["tag"] = s.Tags[i]
						// Read tags in separate go routines.
						go func() {
							s.actions <- action
						}()
					}

				case <-done:
					fmt.Println("stop ticker!")
					break
				}

			}
		}()
	}()

	if s.Name == "" {
		s.Name = "plc_exporter"
	}

	if s.Protocol == "" {
		s.Protocol = "http"
	}

	if s.Port == 0 {
		s.Port = port
	}

	s.save()

	// Set the metrics handler.
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":" + strconv.Itoa(s.Port), Handler: nil}

	// Start the http server.
	// TODO set https as needed.
	go func() {
		if err := server.ListenAndServe(); err != nil {
			// handle err
		}
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {

	}

	// close connections loop
	s.done <- true

	// close ticker...
	done <- true

	// wait cleanup
	<-s.done

	// exit correctly.
	fmt.Println("exit cleanly!")
}
