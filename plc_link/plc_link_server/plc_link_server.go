package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	Interceptors_ "github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/plc/plc_client"
	"github.com/davecourtois/Globular/plc_link/plc_link_client"
	"github.com/davecourtois/Globular/plc_link/plc_linkpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10027
	defaultProxy = 10028

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// Thr IPV4 ServiceId
	ServiceId string = "127.0.0.1"

	domain string = "localhost"
)

// Contain necessary information to connect with a plc_server.
type connection struct {
	Id        string
	ServiceId string // 127.0.0.1:1234
	Domain    string // localhost
}

// Tag type.
type Tag struct {
	ServiceId    string
	ConnectionId string // Must be already define within the plc_server.
	Domain       string
	Name         string // The name of the tag
	TypeName     string // Can be BOOL, SINT, INT, DINT, REAL.
	Offset       int32  // The offset where to begin to read the tag.
	Length       int32  // The size of the array.
	Unsigned     bool   // If true the tag is read as unsigned value
}

type Link struct {
	Id        string    // The link id
	Frequency int32     // The refresh rate in milisecond.
	Source    Tag       // The source tag.
	Target    Tag       // The target tag.
	exit      chan bool // stop synchronization.
}

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id                 string
	Name               string
	Path               string
	Proto              string
	Port               int
	Proxy              int
	Protocol           string
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Domain             string
	CertAuthorityTrust string
	CertFile           string
	KeyFile            string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	//clients map[string]*plc_client.Plc_Client
	clients *sync.Map

	// The list of link to keep up to date.
	Links map[string]Link
}

// Create the configuration file if is not already exist.
func (self *server) init() {

	self.clients = new(sync.Map)

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewPlcLink_Client", plc_link_client.NewPlcLink_Client)

	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		if len(self.Id) == 0 {
			// Generate random id for the server instance.
			self.Id = Utility.RandomUUID()
		}
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

///////////////////// API //////////////////////////////

// Test a connection exist for a given tag.
func (self *server) setTagConnection(tag *Tag) (*plc_client.Plc_Client, error) {
	client, ok := self.clients.Load(tag.ConnectionId)
	if !ok {
		// Open connection with the client.
		var err error
		client, err = plc_client.NewPlc_Client(tag.Domain, tag.ServiceId)
		if err != nil {
			return nil, err
		}
		self.clients.Store(tag.ConnectionId, client)
	}
	return client.(*plc_client.Plc_Client), nil
}

func (self *server) stopSynchronize(lnk *Link) {
	lnk.exit <- true
}

// Return the type name as float number.
func getTypeNameNumber(typeName string) float64 {

	if typeName == "BOOL" {
		return 0.0
	} else if typeName == "SINT" {
		return 1.0
	} else if typeName == "INT" {
		return 2.0
	} else if typeName == "DINT" {
		return 3.0
	} else if typeName == "REAL" {
		return 4.0
	}
	return -1.0
}

func (self *server) startSynchronize(lnk Link) {
	if lnk.exit == nil {
		lnk.exit = make(chan bool)
	}

	// The main sychronization function.
	nbTry := 10 // number of try to connect with the server before given up...
	go func(lnk Link) {
		// So here i will made use of ticker to synchronice the tag.
		tickChan := time.NewTicker(time.Millisecond * time.Duration(lnk.Frequency))
		for {
			select {
			case <-tickChan.C:
				// So here I will get the value from the source and put it in the
				// taget.
				source, err := self.setTagConnection(&lnk.Source)

				if err == nil {
					values, err := source.ReadTag(lnk.Source.ConnectionId, lnk.Source.Name, getTypeNameNumber(lnk.Source.TypeName), lnk.Source.Offset, lnk.Source.Length, lnk.Source.Unsigned)
					if err != nil {
						time.Sleep(1000)
						if nbTry == 0 {
							lnk.exit <- true // stop synchronize on error.
							break
						}
						nbTry--
					}

					// so here I will write the value to the client.
					target, err := self.setTagConnection(&lnk.Target)
					if err == nil {
						err = target.WriteTag(lnk.Target.ConnectionId, lnk.Target.Name, getTypeNameNumber(lnk.Target.TypeName), values, lnk.Target.Offset, lnk.Target.Length, lnk.Target.Unsigned)
						if err != nil {
							time.Sleep(1000)
							if nbTry == 0 {
								lnk.exit <- true // stop synchronize on error.
								break
							}
							nbTry--
						}
					}
				}

			case <-lnk.exit:
				fmt.Println("Done whit tag synchronization for link ", lnk.Id)
				return
			}
		}
	}(lnk)
}

// Link Tow tag together and make it refresh at a given frequency in milisecond.
func (self *server) Link(ctx context.Context, rqst *plc_link_pb.LinkRqst) (*plc_link_pb.LinkRsp, error) {
	// first of all I will test if connection for link exist in the map of connections.
	var src, trg Tag

	// set the source
	src.ServiceId = rqst.Lnk.Source.ServiceId
	src.ConnectionId = rqst.Lnk.Source.ConnectionId
	src.Name = rqst.Lnk.Source.Name
	src.TypeName = rqst.Lnk.Source.TypeName
	src.Offset = rqst.Lnk.Source.Offset
	src.Length = rqst.Lnk.Source.Length

	// Set the target
	trg.ServiceId = rqst.Lnk.Target.ServiceId
	trg.ConnectionId = rqst.Lnk.Target.ConnectionId
	trg.Name = rqst.Lnk.Target.Name
	trg.TypeName = rqst.Lnk.Target.TypeName
	trg.Offset = rqst.Lnk.Target.Offset
	trg.Length = rqst.Lnk.Target.Length

	// So here the connection are found so I will create the link.
	lnk := Link{Id: rqst.Lnk.Id, Frequency: rqst.Lnk.Frequency, Source: src, Target: trg, exit: make(chan bool)}

	// set the lnk...
	self.Links[rqst.Lnk.Id] = lnk

	// In that case I will save it in file.
	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.startSynchronize(lnk)

	return &plc_link_pb.LinkRsp{
		Result: true,
	}, nil
}

// Remove link from tow task.
func (self *server) UnLink(ctx context.Context, rqst *plc_link_pb.UnLinkRqst) (*plc_link_pb.UnLinkRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// Stop the synchronization
	lnk := self.Links[rqst.Id]
	self.stopSynchronize(&lnk)

	// Remove the link from the map.
	delete(self.Links, rqst.Id)

	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &plc_link_pb.UnLinkRsp{
		Result: true,
	}, nil
}

// Suspend the synchronization of tow tags.
func (self *server) Suspend(ctx context.Context, rqst *plc_link_pb.SuspendRqst) (*plc_link_pb.SuspendRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// Stop the synchronization
	lnk := self.Links[rqst.Id]
	self.stopSynchronize(&lnk)

	return &plc_link_pb.SuspendRsp{
		Result: true,
	}, nil
}

// Resume the synchronization of tow tags.
func (self *server) Resume(ctx context.Context, rqst *plc_link_pb.ResumeRqst) (*plc_link_pb.ResumeRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// reset the chan bool...
	lnk := self.Links[rqst.Id]
	lnk.exit = make(chan bool)
	self.Links[rqst.Id] = lnk

	// Start the synchronization.
	self.startSynchronize(lnk)
	return &plc_link_pb.ResumeRsp{
		Result: true,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "plc_link_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(plc_link_pb.File_plc_link_plc_linkpb_plc_link_proto.Services().Get(0).FullName())
	s_impl.Path, _ = os.Executable()
	package_ := string(plc_link_pb.File_plc_link_plc_linkpb_plc_link_proto.Package().Name())
	s_impl.Path = s_impl.Path[strings.Index(s_impl.Path, package_):]
	s_impl.Proto = plc_link_pb.File_plc_link_plc_linkpb_plc_link_proto.Path()
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = domain
	s_impl.Links = make(map[string]Link, 0)
	s_impl.Permissions = make([]interface{}, 0)

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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors_.ServerUnaryInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer()
	}

	plc_link_pb.RegisterPlcLinkServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

	// Start link sychronization.
	for _, lnk := range s_impl.Links {
		s_impl.startSynchronize(lnk)
	}

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {

			if err.Error() == "signal: killed" {
				fmt.Println("service ", s_impl.Name, " was stop!")
			}
		}
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}
