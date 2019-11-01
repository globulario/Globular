package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/dns/dnspb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	//"google.golang.org/grpc/status"
	"github.com/davecourtois/Globular/storage/storage_client"
	"github.com/miekg/dns"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10033
	defaultProxy = 10034

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// Thr IPV4 address
	address      string = "127.0.0.1"
	domain       string = "localhost"
	connectionId string = "dns_service"

	// pointer to the sever implementation.
	s *server
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
	Address         string
	Domain          string
	// self-signed X.509 public keys for distribution
	CertFile string
	// a private RSA key to sign and authenticate the public key
	KeyFile string
	// a private RSA key to sign and authenticate the public key
	CertAuthorityTrust string
	TLS                bool

	// Contain the configuration of the storage service use to store
	// the actual values.
	DnsPort         int    // the dns port
	DnsRoot         string // must be the domain managed by the server.
	StorageService  map[string]interface{}
	StorageDataPath string

	// The link to the storage client.
	storageClient *storage_client.Storage_Client
}

// Create the configuration file if is not already exist.
func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")

	// default value.
	self.DnsPort = 53
	self.DnsRoot = "example.com"
	self.StorageDataPath = ""

	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}

	if len(self.StorageDataPath) == 0 {
		log.Panicln("The value StorageDataPath in the configuration must be given. You can use /tmp (on linux) if you don't want to keep values indefilnely on the storage server.")
	}

	// Here I will initialyse the storage service.
	// note that a storage server must be accessible by the dns service to
	// store it informations.
	if self.StorageService != nil {
		address := self.StorageService["Address"].(string) + ":" + Utility.ToString(self.StorageService["Port"].(float64))
		domain := self.StorageService["Domain"].(string)
		hasTls := self.StorageService["TLS"].(bool)
		keyFile := self.StorageService["KeyFile"].(string)
		certFile := self.StorageService["CertFile"].(string)
		caFile := self.StorageService["CertAuthorityTrust"].(string)

		token := "" // TODO see if it's needed by the storage services.

		// Create the connection with the server.
		self.storageClient = storage_client.NewStorage_Client(domain, address, hasTls, keyFile, certFile, caFile, token)

	} else {
		log.Panicln("No storage service is configure!")
	}
}

// Open the connection if it's close.
func (self *server) openConnection() error {
	err := self.storageClient.CreateConnection(connectionId, connectionId, 0.0) // use persitent storage here.
	if err != nil {
		return err
	}

	err = self.storageClient.OpenConnection(connectionId, `{"path":"`+self.StorageDataPath+`", "name":"dns_data_store"}`)
	if err != nil {
		// close the existing connection
		self.storageClient.CloseConnection(connectionId)
		err = self.storageClient.OpenConnection(connectionId, `{"path":"`+self.StorageDataPath+`", "name":"dns_data_store"}`)
		if err != nil {
			return err
		}
	}
	return nil
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

// Set a dns entry.
func (self *server) SetEntry(ctx context.Context, rqst *dnspb.SetEntryRequest) (*dnspb.SetEntryResponse, error) {
	fmt.Println("Try set dns entry ", rqst.Name+"."+self.DnsRoot)
	domain := rqst.Name + "." + self.DnsRoot
	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	err = self.storageClient.SetItem(connectionId, domain, []byte(rqst.Ipv4))
	if err != nil {
		return nil, err
	}

	return &dnspb.SetEntryResponse{
		Message: domain, // return the full domain.
	}, nil
}

func (self *server) RemoveEntry(ctx context.Context, rqst *dnspb.RemoveEntryRequest) (*dnspb.RemoveEntryResponse, error) {
	fmt.Println("Try remove dns entry ", rqst.Name+"."+self.DnsRoot)
	domain := rqst.Name + "." + self.DnsRoot
	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	err = self.storageClient.RemoveItem(connectionId, domain)
	if err != nil {
		return nil, err
	}

	return &dnspb.RemoveEntryResponse{
		Result: true, // return the full domain.
	}, nil
}

func (self *server) resolve(domain string) (string, error) {
	fmt.Println("Try get dns entry ", domain)
	err := self.openConnection()
	if err != nil {
		return "", err
	}

	address, err := self.storageClient.GetItem(connectionId, domain)
	if err != nil {
		return "", err
	}
	return string(address), nil
}

func (self *server) Resolve(ctx context.Context, rqst *dnspb.ResolveRequest) (*dnspb.ResolveResponse, error) {
	fmt.Println("Try get dns entry ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	ipv4, err := self.storageClient.GetItem(connectionId, rqst.Domain)
	if err != nil {
		return nil, err
	}
	fmt.Println("ipv4 for", domain, "is", string(ipv4))
	return &dnspb.ResolveResponse{
		Ipv4: string(ipv4), // return the full domain.
	}, nil
}

// Set a text entry.
func (self *server) SetText(ctx context.Context, rqst *dnspb.SetTextRequest) (*dnspb.SetTextResponse, error) {
	fmt.Println("Try set dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	values, err := json.Marshal(rqst.Values)
	if err != nil {
		return nil, err
	}

	err = self.storageClient.SetItem(connectionId, rqst.Id, values)
	if err != nil {
		return nil, err
	}

	return &dnspb.SetTextResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getText(id string) ([]string, error) {
	fmt.Println("Try get dns text ", id)
	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	data, err := self.storageClient.GetItem(connectionId, id)

	values := make([]string, 0)
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}
	return values, nil
}

// Retreive a text value
func (self *server) GetText(ctx context.Context, rqst *dnspb.GetTextRequest) (*dnspb.GetTextResponse, error) {
	fmt.Println("Try get dns text ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	data, err := self.storageClient.GetItem(connectionId, rqst.Id)
	if err != nil {
		return nil, err
	}

	values := make([]string, 0)
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	return &dnspb.GetTextResponse{
		Values: values, // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveText(ctx context.Context, rqst *dnspb.RemoveTextRequest) (*dnspb.RemoveTextResponse, error) {
	fmt.Println("Try remove dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, err
	}

	err = self.storageClient.RemoveItem(connectionId, rqst.Id)
	if err != nil {
		return nil, err
	}

	return &dnspb.RemoveTextResponse{
		Result: true, // return the full domain.
	}, nil
}

/////////////////////// DNS Specific service //////////////////////
type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	log.Println("-----> dns resquest receive... ", msg)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, err := s.resolve(domain) // resolve the address name from the

		if err == nil {
			log.Println("---> ask for domain: ", domain, " address to redirect is ", address)
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		} else {
			log.Println(err)
		}

	case dns.TypeTXT:
		id := msg.Question[0].Name
		log.Println("---> look for value ", id)
		values, err := s.getText(id)
		if err == nil {
			log.Println("---> values found ", values)
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.TXT{
				// keep text values.
				Hdr: dns.RR_Header{Name: id, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
				Txt: values,
			})
		}
	}
	w.WriteMsg(&msg)
}

func ServeDns(port int) {
	// Now I will start the dns server.
	srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Failed to set udp listener %s\n", err.Error())
	}
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "dns_service: ", log.LstdFlags))

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
	s_impl.Address = address
	s_impl.Domain = domain
	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	s = s_impl // set the pointer to the server.

	// First of all I will creat a listener.
	// Create the channel to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("could not list on %s: %s", s_impl.Address, err)
		return
	}

	var grpcServer *grpc.Server
	if s_impl.TLS {
		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(s_impl.CertFile, s_impl.KeyFile)
		if err != nil {
			log.Fatalf("could not load server key pair: %s", err)
			return
		} else {
			log.Println("load certificate from ", s_impl.CertFile, s_impl.KeyFile)
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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.UnaryAuthInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer()
	}

	dnspb.RegisterDnsServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
		log.Println(s_impl.Name + " grpc service is closed")
		s_impl.storageClient.CloseConnection("dns_service")
	}()

	// start lisen on the network for dns queries...
	go func() {
		log.Println("--> start lisen for dns queries at port", s_impl.DnsPort)
		ServeDns(s_impl.DnsPort)
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}
