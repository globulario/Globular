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

	"github.com/davecourtois/Globular/ldap/ldappb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	LDAP "github.com/mavricknz/ldap"
)

var (
	defaultPort  = 10003
	defaultProxy = 10004

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// Thr IPV4 address
	address string = "127.0.0.1"

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

	conn *LDAP.LDAPConnection
}

type server struct {

	// The global attribute of the services.
	Name               string
	Port               int
	Proxy              int
	Protocol           string
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Address            string
	Domain             string
	CertAuthorityTrust string
	CertFile           string
	KeyFile            string
	TLS                bool

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
func (self *server) connect(id string, userId string, pwd string) (*LDAP.LDAPConnection, error) {

	// The info must be set before that function is call.
	info := self.Connections[id]
	conn := LDAP.NewLDAPConnection(info.Host, uint16(info.Port))

	// Try to connect to Ldap, return timeout error after tree second.
	conn.NetworkConnectTimeout = time.Duration(3 * time.Second)
	conn.AbandonMessageOnReadTimeout = true

	err := conn.Connect()
	if err != nil {
		log.Println("106 ----> ", err)
		// handle error
		return nil, err
	}

	// Connect with the default user...
	if len(userId) > 0 {
		err := conn.Bind(userId, pwd)
		if err != nil {
			return nil, err
		}
	} else {
		err := conn.Bind(info.User, info.Password)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

/*
func (this *LdapManager) authenticate(id string, login string, psswd string) bool {

	// Now I will try to make a simple query if it fail that's mean the user
	// does have the permission...
	var base_dn string = "OU=Users," + this.getConfigsInfo()[id].M_searchBase
	var filter string = "(objectClass=user)"

	// Test get some user...
	var attributes []string = []string{"sAMAccountName"}
	_, err := this.search(id, login, psswd, base_dn, filter, attributes)
	if err != nil {
		Utility.Log(Utility.FunctionName(), Utility.FileLine(), "---> ldap authenticate fail: ", login, err)
		return false
	}

	return true
}
*/

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

	err := self.Connections[id].conn.Close()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

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
	s_impl.Address = address
	s_impl.Domain = domain
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

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
		grpcServer = grpc.NewServer(grpc.Creds(creds))

	} else {
		grpcServer = grpc.NewServer()
	}

	ldappb.RegisterLdapServiceServer(grpcServer, s_impl)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")

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
	if self.Connections[id].conn == nil {
		c := self.Connections[id]
		conn, err := self.connect(id, self.Connections[id].User, self.Connections[id].Password)
		if err != nil {
			return nil, err
		}

		c.conn = conn
		self.Connections[id] = c
	}

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
			attributeValue := entry.GetAttributeValue(attributeName)
			row = append(row, attributeValue)
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
