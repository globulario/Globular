package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"

	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"time"

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/event/eventpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10015
	defaultProxy = 10016

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// the default domain.
	domain string = "localhost"
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Name            string
	Port            int
	Proxy           int
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.
	Protocol        string
	Domain          string
	// self-signed X.509 public keys for distribution
	CertFile string
	// a private RSA key to sign and authenticate the public key
	KeyFile string
	// a private RSA key to sign and authenticate the public key
	CertAuthorityTrust string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool

	// Use to sync event channel manipulation.
	actions chan map[string]interface{}
}

// Create the configuration file if is not already exist.
func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}
}

// Save the configuration values.
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

//////////////////////////////////////////////////////////////////////////////
//	Services implementation.
//////////////////////////////////////////////////////////////////////////////

// That function process channel operation and run in it own go routine.
func (self *server) run() {

	channels := make(map[string][]string)
	streams := make(map[string]eventpb.EventService_OnEventServer)
	quits := make(map[string]chan bool)

	// Here will create the action channel.
	self.actions = make(chan map[string]interface{})

	for {
		select {
		case a := <-self.actions:
			action := a["action"].(string)
			if action == "onevent" {
				streams[a["uuid"].(string)] = a["stream"].(eventpb.EventService_OnEventServer)
				quits[a["uuid"].(string)] = a["quit"].(chan bool)
			} else if action == "subscribe" {
				if channels[a["name"].(string)] == nil {
					channels[a["name"].(string)] = make([]string, 0)
				}
				if !Utility.Contains(channels[a["name"].(string)], a["uuid"].(string)) {
					channels[a["name"].(string)] = append(channels[a["name"].(string)], a["uuid"].(string))
				}
			} else if action == "publish" {
				if channels[a["name"].(string)] != nil {
					toDelete := make([]string, 0)
					for i := 0; i < len(channels[a["name"].(string)]); i++ {
						uuid := channels[a["name"].(string)][i]
						stream := streams[uuid]
						if stream != nil {
							// Here I will send data to stream.
							err := stream.Send(&eventpb.OnEventResponse{
								Evt: &eventpb.Event{
									Name: a["name"].(string),
									Data: a["data"].([]byte),
								},
							})

							// In case of error I will remove the subscriber
							// from the list.
							if err != nil {
								// append to channle list to be close.
								toDelete = append(toDelete, uuid)
								log.Println("error publish event to ", uuid, err)
							}
						}
					}

					// remove closed channel
					for i := 0; i < len(toDelete); i++ {
						uuid := toDelete[i]
						// remove uuid from all channels.
						for name, channel := range channels {
							uuids := make([]string, 0)
							for i := 0; i < len(channel); i++ {
								if uuid != channel[i] {
									uuids = append(uuids, channel[i])
								}
							}
							channels[name] = uuids
						}
						// return from OnEvent
						quits[uuid] <- true
						// remove the channel from the map.
						delete(quits, uuid)
					}
				}
			} else if action == "unsubscribe" {
				uuids := make([]string, 0)
				for i := 0; i < len(channels[a["name"].(string)]); i++ {
					if a["uuid"].(string) != channels[a["name"].(string)][i] {
						uuids = append(uuids, channels[a["name"].(string)][i])
					}
				}
				channels[a["name"].(string)] = uuids
			} else if action == "quit" {
				// remove uuid from all channels.
				for name, channel := range channels {
					uuids := make([]string, 0)
					for i := 0; i < len(channel); i++ {
						if a["uuid"].(string) != channel[i] {
							uuids = append(uuids, channel[i])
						}
					}
					channels[name] = uuids
				}
				// return from OnEvent
				quits[a["uuid"].(string)] <- true
				// remove the channel from the map.
				delete(quits, a["uuid"].(string))
			}
		}
	}
}

// Connect to an event channel or create it if it not already exist
// and stay in that function until UnSubscribe is call.
func (self *server) Quit(ctx context.Context, rqst *eventpb.QuitRequest) (*eventpb.QuitResponse, error) {
	quit := make(map[string]interface{})
	quit["action"] = "quit"
	quit["uuid"] = rqst.Uuid

	self.actions <- quit

	return &eventpb.QuitResponse{
		Result: true,
	}, nil
}

// Connect to an event channel or create it if it not already exist
// and stay in that function until UnSubscribe is call.
func (self *server) OnEvent(rqst *eventpb.OnEventRequest, stream eventpb.EventService_OnEventServer) error {
	onevent := make(map[string]interface{})
	onevent["action"] = "onevent"
	onevent["stream"] = stream
	onevent["uuid"] = rqst.Uuid
	onevent["quit"] = make(chan bool)

	self.actions <- onevent

	// wait util unsbscribe or connection is close.
	<-onevent["quit"].(chan bool)

	return nil
}

func (self *server) Subscribe(ctx context.Context, rqst *eventpb.SubscribeRequest) (*eventpb.SubscribeResponse, error) {
	subscribe := make(map[string]interface{})
	subscribe["action"] = "subscribe"
	subscribe["name"] = rqst.Name
	subscribe["uuid"] = rqst.Uuid
	self.actions <- subscribe

	return &eventpb.SubscribeResponse{
		Result: true,
	}, nil
}

// Disconnect to an event channel.(Return from Subscribe)
func (self *server) UnSubscribe(ctx context.Context, rqst *eventpb.UnSubscribeRequest) (*eventpb.UnSubscribeResponse, error) {
	unsubscribe := make(map[string]interface{})
	unsubscribe["action"] = "unsubscribe"
	unsubscribe["name"] = rqst.Name
	unsubscribe["uuid"] = rqst.Uuid

	self.actions <- unsubscribe

	return &eventpb.UnSubscribeResponse{
		Result: true,
	}, nil
}

// Publish event on channel.
func (self *server) Publish(ctx context.Context, rqst *eventpb.PublishRequest) (*eventpb.PublishResponse, error) {
	publish := make(map[string]interface{})
	publish["action"] = "publish"
	publish["name"] = rqst.Evt.Name
	publish["data"] = rqst.Evt.Data

	// publish the data.
	self.actions <- publish
	log.Println("----> publish event ", rqst.Evt.Name)
	return &eventpb.PublishResponse{
		Result: true,
	}, nil

	return nil, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "event_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = Utility.GetExecName(os.Args[0])
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = "localhost"

	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	// First of all I will creat a listener.
	// Create the channel to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("could not list on %s: %s", s_impl.Domain, err)
		return
	}

	var grpcServer *grpc.Server
	if s_impl.TLS {
		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(s_impl.CertFile, s_impl.KeyFile)
		if err != nil {
			log.Fatalf("could not load server key pair: %s", err)
			return
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(s_impl.CertAuthorityTrust)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
			return
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append client certs")
			return
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		// Create the gRPC server with the credentials
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.UnaryAuthInterceptor), grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute, // <--- This fixes it!
		})}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer()
	}

	eventpb.RegisterEventServiceServer(grpcServer, s_impl)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")
		go s_impl.run()

		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
		}
		log.Println(s_impl.Name + " grpc service is closed")
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
