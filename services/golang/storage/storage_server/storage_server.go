package main

import (
	"context"

	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	//	"time"
	globular "github.com/davecourtois/Globular/services/golang/globular_service"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/storage/storage_client"
	"github.com/davecourtois/Globular/services/golang/storage/storage_store"
	"github.com/davecourtois/Globular/services/golang/storage/storagepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	// "google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	defaultPort  = 10013
	defaultProxy = 10014

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The domain
	domain string = "localhost"
)

// Keep connection information here.
type connection struct {
	Id   string // The connection id
	Name string // The kv store name
	Type storagepb.StoreType
}

type server struct {

	// The global attribute of the services.
	Id              string
	Name            string
	Path            string
	Proto           string
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
	Permissions        []interface{} // contains the action permission for the services.
	// The grpc server.
	grpcServer *grpc.Server

	// The map of connection...
	Connections map[string]connection

	// the map of store
	stores map[string]storage_store.Store
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

	self.stores = make(map[string]storage_store.Store)

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewStorage_Client", storage_client.NewStorage_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := globular.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = globular.InitGrpcServer(self, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err != nil {
		return err
	}

	return nil

}

// Save the configuration values.
func (self *server) Save() error {
	// Create the file...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return globular.SaveService(dir+"/config.json", self)
}

func (self *server) StartService() error {
	return globular.StartService(self, self.grpcServer)
}

func (self *server) StopService() error {
	return globular.StopService(self)
}

func (self *server) Stop(context.Context, *storagepb.StopRequest) (*storagepb.StopResponse, error) {
	return &storagepb.StopResponse{}, self.StopService()
}

//////////////////////// Storage specific functions ////////////////////////////

// Create a new KV connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *storagepb.CreateConnectionRqst) (*storagepb.CreateConnectionRsp, error) {
	if rsqt.Connection == nil {
		return nil, errors.New("The request dosent contain connection object!")
	}

	if _, ok := self.Connections[rsqt.Connection.Id]; ok {
		self.stores[rsqt.Connection.Id].Close() // close the previous connection.
	}

	fmt.Println("Try to create a new connection with id: ", rsqt.Connection.Id)
	var c connection
	var err error

	// Set the connection info from the request.
	c.Id = rsqt.Connection.Id
	c.Name = rsqt.Connection.Name

	// set or update the connection and save it in json file.
	self.Connections[c.Id] = c

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

	globular.UpdateServiceConfig(self)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print the success message here.
	fmt.Println("Connection " + c.Id + " was created with success!")

	return &storagepb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Remove a connection from the map and the file.
func (self *server) DeleteConnection(ctx context.Context, rqst *storagepb.DeleteConnectionRqst) (*storagepb.DeleteConnectionRsp, error) {

	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return &storagepb.DeleteConnectionRsp{
			Result: true,
		}, nil
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
	return &storagepb.DeleteConnectionRsp{
		Result: true,
	}, nil

}

// Open the store and set options...
func (self *server) Open(ctx context.Context, rqst *storagepb.OpenRqst) (*storagepb.OpenRsp, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	var store storage_store.Store
	conn := self.Connections[rqst.GetId()]

	// Create the store object.
	if conn.Type == storagepb.StoreType_LEVEL_DB {
		store = storage_store.NewLevelDB_store()
	} else if conn.Type == storagepb.StoreType_BIG_CACHE {
		store = storage_store.NewBigCache_store()
	}

	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.Open(rqst.GetOptions())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.stores[rqst.GetId()] = store

	return &storagepb.OpenRsp{
		Result: true,
	}, nil
}

// Close the data store.
func (self *server) Close(ctx context.Context, rqst *storagepb.CloseRqst) (*storagepb.CloseRsp, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	if self.stores[rqst.GetId()] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := self.stores[rqst.GetId()].Close()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.CloseRsp{
		Result: true,
	}, nil
}

// Save an item in the kv store
func (self *server) SetItem(ctx context.Context, rqst *storagepb.SetItemRequest) (*storagepb.SetItemResponse, error) {

	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	store := self.stores[rqst.GetId()]
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.SetItem(rqst.GetKey(), rqst.GetValue())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.SetItemResponse{
		Result: true,
	}, nil
}

// Retreive a value with a given key
func (self *server) GetItem(ctx context.Context, rqst *storagepb.GetItemRequest) (*storagepb.GetItemResponse, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	store := self.stores[rqst.GetId()]
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}
	log.Println("--> try to find key with value:", rqst.GetKey())
	data, err := store.GetItem(rqst.GetKey())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.GetItemResponse{
		Result: data,
	}, nil

	return nil, nil
}

// Remove an item with a given key
func (self *server) RemoveItem(ctx context.Context, rqst *storagepb.RemoveItemRequest) (*storagepb.RemoveItemResponse, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	store := self.stores[rqst.GetId()]
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.RemoveItem(rqst.GetKey())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.RemoveItemResponse{
		Result: true,
	}, nil

}

// Remove all items
func (self *server) Clear(ctx context.Context, rqst *storagepb.ClearRequest) (*storagepb.ClearResponse, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	store := self.stores[rqst.GetId()]
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.Clear()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.ClearResponse{
		Result: true,
	}, nil
}

// Drop a store
func (self *server) Drop(ctx context.Context, rqst *storagepb.DropRequest) (*storagepb.DropResponse, error) {
	if _, ok := self.Connections[rqst.GetId()]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.GetId())))
	}

	store := self.stores[rqst.GetId()]
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.Drop()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &storagepb.DropResponse{
		Result: true,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, "storage_service: ", log.LstdFlags))

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
	s_impl.Name = string(storagepb.File_services_proto_storage_proto.Services().Get(0).FullName())
	s_impl.Proto = storagepb.File_services_proto_storage_proto.Path()
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = "localhost"

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		log.Fatalf("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
	}

	// Register the echo services
	storagepb.RegisterStorageServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)
	// Start the service.
	s_impl.StartService()

}
