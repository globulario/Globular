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
	"time"

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/ldap/ldappb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	LDAP "gopkg.in/ldap.v3"
)

var (
	defaultPort  = 10031
	defaultProxy = 10032

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default domain
	domain string = "localhost"
)

// Keep connection information here.
type connection struct {
	Id       string // The connection id
	Host     string // can also be ipv4 addresse.
	User     string
	Password string
	Port     int32
	conn     *LDAP.Conn
}

type server struct {

	// The global attribute of the services.
	Name               string
	Port               int
	Proxy              int
	Protocol           string
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Domain             string
	CertAuthorityTrust string
	CertFile           string
	KeyFile            string
	Version            string
	TLS                bool
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool

	// The map of connection...
	Connections map[string]connection
}

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

/**
 * Connect to a ldap server...
 */
func (self *server) connect(id string, userId string, pwd string) (*LDAP.Conn, error) {

	// The info must be set before that function is call.
	info := self.Connections[id]

	conn, err := LDAP.Dial("tcp", fmt.Sprintf("%s:%d", info.Host, info.Port))
	if err != nil {
		// handle error
		return nil, err
	}

	conn.SetTimeout(time.Duration(3 * time.Second))

	// Connect with the default user...
	if len(userId) > 0 {
		if len(pwd) > 0 {
			err = conn.Bind(userId, pwd)
		} else {
			err = conn.UnauthenticatedBind(userId)
		}
		if err != nil {
			return nil, err
		}
	} else {
		if len(info.Password) > 0 {
			err = conn.Bind(info.User, info.Password)
		} else {
			err = conn.UnauthenticatedBind(info.User)
		}
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

// Authenticate a user with LDAP server.
func (self *server) Authenticate(ctx context.Context, rqst *ldappb.AuthenticateRqst) (*ldappb.AuthenticateRsp, error) {
	id := rqst.Id
	login := rqst.Login
	pwd := rqst.Pwd

	if len(id) > 0 {
		// I will made use of bind to authenticate the user.
		_, err := self.connect(id, login, pwd)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}
	} else {
		for id, _ := range self.Connections {
			_, err := self.connect(id, login, pwd)
			if err == nil {
				return &ldappb.AuthenticateRsp{
					Result: true,
				}, nil
			}
		}
		// fail to authenticate.
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Authentication fail for user "+rqst.Login)))
	}

	return &ldappb.AuthenticateRsp{
		Result: true,
	}, nil
}

// Create a new SQL connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *ldappb.CreateConnectionRqst) (*ldappb.CreateConnectionRsp, error) {
	fmt.Println("Try to create a new connection")
	// sqlpb
	fmt.Println("Try to create a new connection")
	var c connection
	var err error

	// Set the connection info from the request.
	c.Id = rsqt.Connection.Id
	c.Host = rsqt.Connection.Host
	c.Port = rsqt.Connection.Port
	c.User = rsqt.Connection.User
	c.Password = rsqt.Connection.Password

	// set or update the connection and save it in json file.
	self.Connections[c.Id] = c

	// So here I will create the new ldap connection.
	log.Println("try to connect: ", c.Id, c.Host, c.Password)

	c.conn, err = self.connect(c.Id, c.User, c.Password)
	defer c.conn.Close()

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// In that case I will save it in file.
	err = self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// test if the connection is reacheable.
	//_, err = self.ping(ctx, c.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print the success message here.
	log.Println("Connection " + c.Id + " was created with success!")

	return &ldappb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Remove a connection from the map and the file.
func (self *server) DeleteConnection(ctx context.Context, rqst *ldappb.DeleteConnectionRqst) (*ldappb.DeleteConnectionRsp, error) {

	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return &ldappb.DeleteConnectionRsp{
			Result: true,
		}, nil
	}

	if self.Connections[id].conn != nil {
		// Close the connection.
		self.Connections[id].conn.Close()
	}

	delete(self.Connections, id)

	// In that case I will save it in file.
	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return success.
	return &ldappb.DeleteConnectionRsp{
		Result: true,
	}, nil

}

// Close connection.
func (self *server) Close(ctx context.Context, rqst *ldappb.CloseRqst) (*ldappb.CloseRsp, error) {
	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Connection "+id+" dosent exist!")))
	}

	self.Connections[id].conn.Close()

	// return success.
	return &ldappb.CloseRsp{
		Result: true,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "ldap_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Connections = make(map[string]connection)
	s_impl.Name = Utility.GetExecName(os.Args[0])
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = "localhost"

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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.ServerUnaryAuthInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamAuthInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(Interceptors.ServerUnaryAuthInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamAuthInterceptor))
	}

	ldappb.RegisterLdapServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)
	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")

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

/**
 * Search for a list of value over the ldap server. if the base_dn is
 * not specify the default base is use. It return a list of values. This can
 * be interpret as a tow dimensional array.
 */
func (self *server) search(id string, base_dn string, filter string, attributes []string) ([][]interface{}, error) {

	if _, ok := self.Connections[id]; !ok {
		return nil, errors.New("Connection " + id + " dosent exist!")
	}

	// create the connection.
	c := self.Connections[id]
	conn, err := self.connect(id, self.Connections[id].User, self.Connections[id].Password)
	if err != nil {
		return nil, err
	}

	c.conn = conn
	self.Connections[id] = c

	// close connection after search.
	defer c.conn.Close()

	//Now I will execute the query...
	search_request := LDAP.NewSearchRequest(
		base_dn,
		LDAP.ScopeWholeSubtree, LDAP.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)

	// Create simple search.
	sr, err := self.Connections[id].conn.Search(search_request)

	if err != nil {
		return nil, err
	}

	// Store the founded values in results...
	var results [][]interface{}
	for i := 0; i < len(sr.Entries); i++ {
		entry := sr.Entries[i]
		var row []interface{}
		for j := 0; j < len(attributes); j++ {
			attributeName := attributes[j]
			attributeValues := entry.GetAttributeValues(attributeName)
			row = append(row, attributeValues)
		}
		results = append(results, row)
	}

	return results, nil
}

// Search over LDAP server.
func (self *server) Search(ctx context.Context, rqst *ldappb.SearchRqst) (*ldappb.SearchResp, error) {
	id := rqst.Search.GetId()
	if _, ok := self.Connections[id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Connection "+id+" dosent exist!")))
	}

	results, err := self.search(id, rqst.Search.BaseDN, rqst.Search.Filter, rqst.Search.Attributes)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I got the results.
	str, err := json.Marshal(results)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &ldappb.SearchResp{
		Result: string(str),
	}, nil
}
