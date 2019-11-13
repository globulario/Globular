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

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/event/eventpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	//"google.golang.org/grpc/status"
	"time"

	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
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

// Channel's use by the server.
type SubscribeEvent struct {
	name string // The name of the event
	uuid string // The subscriber unique
	data chan []byte
	quit chan chan bool
}

type UnSubscribeEvent struct {
	name string // The name of the event
	uuid string // The subscriber unique
	quit chan bool
}

type PublishEvent struct {
	name string // The name of the event
	data []byte
}

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

	// Use to sync event channel manipulation.
	subscribe_events_chan   chan *SubscribeEvent
	unsubscribe_events_chan chan *UnSubscribeEvent
	pulish_events_chan      chan *PublishEvent

	quit chan string
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

	// a -> event name -> subscriber uuid -> channel
	events := make(map[string]map[string]*SubscribeEvent, 0)

	self.subscribe_events_chan = make(chan *SubscribeEvent, 0)
	self.unsubscribe_events_chan = make(chan *UnSubscribeEvent, 0)
	self.pulish_events_chan = make(chan *PublishEvent, 0)

	for {
		select {
		case evt := <-self.subscribe_events_chan:
			log.Println("---> subscribe event receive: ", evt.name)
			if events[evt.name] == nil {
				events[evt.name] = make(map[string]*SubscribeEvent, 0)
			}
			events[evt.name][evt.uuid] = evt

		case evt := <-self.unsubscribe_events_chan:
			log.Println("---> unsubscribe event receive: ", evt.name, evt.uuid)
			if events[evt.name] != nil {
				if events[evt.name][evt.uuid] != nil {
					events[evt.name][evt.uuid].quit <- evt.quit
					delete(events[evt.name], evt.uuid)
				} else {
					log.Println("--> no subscriber found with uuid ", evt.uuid)
				}
			}

		case evt := <-self.pulish_events_chan:
			log.Println("---> publish event receive: ", evt.name)
			if events[evt.name] != nil {
				for _, subscriber := range events[evt.name] {
					// publish the data on the channel.
					subscriber.data <- evt.data
				}
			}

		}
	}

}

// Connect to an event channel or create it if it not already exist
// and stay in that function until UnSubscribe is call.
func (self *server) Subscribe(rqst *eventpb.SubscribeRequest, stream eventpb.EventService_SubscribeServer) error {

	// create a new channel.
	uuid := Utility.RandomUUID()
	evt := new(SubscribeEvent)
	evt.name = rqst.Name
	evt.data = make(chan []byte)
	evt.quit = make(chan chan bool)
	evt.uuid = uuid

	// subscribe to the channel.
	self.subscribe_events_chan <- evt

	// send back the uuid to client for furder unsubscribe request.
	stream.Send(&eventpb.SubscribeResponse{
		Result: &eventpb.SubscribeResponse_Uuid{
			Uuid: evt.uuid,
		},
	})

	for {
		select {
		case data := <-evt.data:

			err := stream.Send(&eventpb.SubscribeResponse{
				Result: &eventpb.SubscribeResponse_Evt{
					Evt: &eventpb.Event{
						Name: rqst.Name,
						Data: data,
					},
				},
			})
			if err != nil {
				// Remove it from the list of subscribers.
				evt_ := new(UnSubscribeEvent)
				evt_.name = rqst.Name
				evt_.uuid = uuid
				evt_.quit = make(chan bool) // no used...
				self.unsubscribe_events_chan <- evt_
				log.Println("error", uuid, err.Error())
				break
			}

		case quit_channel := <-evt.quit:
			// exit
			log.Println("--> quit subscription ", uuid)
			quit_channel <- true
			return nil
		}
	}

	return nil
}

// Disconnect to an event channel.(Return from Subscribe)
func (self *server) UnSubscribe(ctx context.Context, rqst *eventpb.UnSubscribeRequest) (*eventpb.UnSubscribeResponse, error) {
	evt := new(UnSubscribeEvent)
	evt.name = rqst.Name
	evt.uuid = rqst.Uuid
	evt.quit = make(chan bool)
	self.unsubscribe_events_chan <- evt

	//  wait for the subscription loop function to stop
	<-evt.quit

	return &eventpb.UnSubscribeResponse{
		Result: true,
	}, nil

}

// Publish event on channel.
func (self *server) Publish(ctx context.Context, rqst *eventpb.PublishRequest) (*eventpb.PublishResponse, error) {
	evt := new(PublishEvent)
	evt.name = rqst.Evt.Name
	evt.data = rqst.Evt.Data

	// dispatch the evt.
	self.pulish_events_chan <- evt

	return &eventpb.PublishResponse{
		Result: true,
	}, nil
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
	reflection.Register(grpcServer)
	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")
		go s_impl.run()
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		log.Println(s_impl.Name + " grpc service is closed")
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}
