package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/davecourtois/Globular/api"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/ldap/ldappb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	//"google.golang.org/grpc/grpclog"
	"github.com/davecourtois/Globular/api/client"
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
	Version            string
	TLS                bool
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	// The grpc server.
	grpcServer *grpc.Server

	// The map of connection...
	Connections map[string]connection
}

// Globular services implementation...
// The id of a particular service instance.
func (self *server) GetId() string {
	return self.Id
}
func (self *server) SetId(id string) {
	self.Id = id
}

// The name of a service, must be the gRpc Service name.
func (self *server) GetName() string {
	return self.Name
}
func (self *server) SetName(name string) {
	self.Name = name
}

// The path of the executable.
func (self *server) GetPath() string {
	return self.Path
}
func (self *server) SetPath(path string) {
	self.Path = path
}

// The path of the .proto file.
func (self *server) GetProto() string {
	return self.Proto
}
func (self *server) SetProto(proto string) {
	self.Proto = proto
}

// The gRpc port.
func (self *server) GetPort() int {
	return self.Port
}
func (self *server) SetPort(port int) {
	self.Port = port
}

// The reverse proxy port (use by gRpc Web)
func (self *server) GetProxy() int {
	return self.Proxy
}
func (self *server) SetProxy(proxy int) {
	self.Proxy = proxy
}

// Can be one of http/https/tls
func (self *server) GetProtocol() string {
	return self.Protocol
}
func (self *server) SetProtocol(protocol string) {
	self.Protocol = protocol
}

// Return true if all Origins are allowed to access the mircoservice.
func (self *server) GetAllowAllOrigins() bool {
	return self.AllowAllOrigins
}
func (self *server) SetAllowAllOrigins(allowAllOrigins bool) {
	self.AllowAllOrigins = allowAllOrigins
}

// If AllowAllOrigins is false then AllowedOrigins will contain the
// list of address that can reach the services.
func (self *server) GetAllowedOrigins() string {
	return self.AllowedOrigins
}

func (self *server) SetAllowedOrigins(allowedOrigins string) {
	self.AllowedOrigins = allowedOrigins
}

// Can be a ip address or domain name.
func (self *server) GetDomain() string {
	return self.Domain
}
func (self *server) SetDomain(domain string) {
	self.Domain = domain
}

// TLS section

// If true the service run with TLS. The
func (self *server) GetTls() bool {
	return self.TLS
}
func (self *server) SetTls(hasTls bool) {
	self.TLS = hasTls
}

// The certificate authority file
func (self *server) GetCertAuthorityTrust() string {
	return self.CertAuthorityTrust
}
func (self *server) SetCertAuthorityTrust(ca string) {
	self.CertAuthorityTrust = ca
}

// The certificate file.
func (self *server) GetCertFile() string {
	return self.CertFile
}
func (self *server) SetCertFile(certFile string) {
	self.CertFile = certFile
}

// The key file.
func (self *server) GetKeyFile() string {
	return self.KeyFile
}
func (self *server) SetKeyFile(keyFile string) {
	self.KeyFile = keyFile
}

// The service version
func (self *server) GetVersion() string {
	return self.Version
}
func (self *server) SetVersion(version string) {
	self.Version = version
}

// The publisher id.
func (self *server) GetPublisherId() string {
	return self.PublisherId
}
func (self *server) SetPublisherId(publisherId string) {
	self.PublisherId = publisherId
}

func (self *server) GetKeepUpToDate() bool {
	return self.KeepUpToDate
}
func (self *server) SetKeepUptoDate(val bool) {
	self.KeepUpToDate = val
}

func (self *server) GetKeepAlive() bool {
	return self.KeepAlive
}
func (self *server) SetKeepAlive(val bool) {
	self.KeepAlive = val
}

func (self *server) GetPermissions() []interface{} {
	return self.Permissions
}
func (self *server) SetPermissions(permissions []interface{}) {
	self.Permissions = permissions
}

// Create the configuration file if is not already exist.
func (self *server) Init() error {

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewLdap_Client", client.NewLdap_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := api.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = api.InitGrpcServer(self, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err != nil {
		return err
	}

	return nil

}

// Save the configuration values.
func (self *server) Save() error {
	// Create the file...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return api.SaveService(dir+"/config.json", self)
}

func (self *server) Start() error {
	return api.StartService(self, self.grpcServer)
}

func (self *server) Stop() error {
	return api.StopService(self)
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

	c.conn, err = self.connect(c.Id, c.User, c.Password)
	defer c.conn.Close()

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// In that case I will save it in file.
	err = self.Save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	api.UpdateServiceConfig(self)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

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
	err := self.Save()
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
	//grpclog.SetLogger(log.New(os.Stdout, "ldap_service: ", log.LstdFlags))

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
	s_impl.Name = string(ldappb.File_ldap_ldappb_ldap_proto.Services().Get(0).FullName())
	s_impl.Proto = ldappb.File_ldap_ldappb_ldap_proto.Path()
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = domain
	s_impl.Permissions = make([]interface{}, 0)

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		log.Fatalf("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
	}

	// Register the echo services
	ldappb.RegisterLdapServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)

	// Start the service.
	s_impl.Start()
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
