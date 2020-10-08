package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	globular "github.com/davecourtois/Globular/services/golang/globular_service"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/persistence/persistence_client"
	"github.com/davecourtois/Globular/services/golang/persistence/persistence_store"
	"github.com/davecourtois/Globular/services/golang/persistence/persistencepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	//"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	defaultPort  = 10035
	defaultProxy = 10036

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default domain
	domain string = "localhost"

	// The grpc server.
	grpcServer *grpc.Server
)

// This is the connction to a datastore.
type connection struct {
	Id       string
	Name     string
	Host     string
	Store    persistencepb.StoreType
	User     string
	Port     int32
	Timeout  int32
	Options  string
	password string
}

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id                 string
	Name               string
	Path               string
	Port               int
	Proto              string
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

	// The grpc server.
	grpcServer *grpc.Server

	// saved connections
	Connections map[string]connection

	// unsaved connections
	connections map[string]connection

	// The map of store (also connections...)
	stores map[string]persistence_store.Store
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
	Utility.RegisterFunction("NewPersistenceService_Client", persistence_client.NewPersistenceService_Client)

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

	// init the connections map.
	self.Connections = make(map[string]connection)
	self.connections = make(map[string]connection)

	// initialyse store connection here.
	self.stores = make(map[string]persistence_store.Store)
	// Here I will initialyse the connection.
	for _, c := range self.Connections {
		if c.Store == persistencepb.StoreType_MONGO {
			// here I will create a new mongo data store.
			s := new(persistence_store.MongoStore)

			// Now I will try to connect...
			err := s.Connect(c.Id, c.Host, c.Port, c.User, c.password, c.Name, c.Timeout, c.Options)
			// keep the store for futur call...
			if err == nil {
				self.stores[c.Id] = s
			} else {
				return err
			}
		}
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
	return globular.StopService(self, self.grpcServer)
}

////////////////////////// Persistence specific functions //////////////////////

// Create a new Store connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rqst *persistencepb.CreateConnectionRqst) (*persistencepb.CreateConnectionRsp, error) {
	// sqlpb
	var c connection
	var err error

	// use existing connection as we can.
	if _, ok := self.connections[rqst.Connection.Id]; ok {
		c = self.connections[rqst.Connection.Id]
	} else if _, ok := self.Connections[rqst.Connection.Id]; ok {
		c = self.Connections[rqst.Connection.Id]
	} else {

		// Set the connection info from the request.
		c.Id = rqst.Connection.Id
		c.Name = rqst.Connection.Name
		c.Host = rqst.Connection.Host
		c.Port = rqst.Connection.Port
		c.User = rqst.Connection.User
		c.password = rqst.Connection.Password // The password must never be store here.
		c.Store = rqst.Connection.Store

		if c.Store == persistencepb.StoreType_MONGO {
			// here I will create a new mongo data store.
			s := new(persistence_store.MongoStore)

			// Now I will try to connect...
			err := s.Connect(c.Id, c.Host, c.Port, c.User, c.password, c.Name, c.Timeout, c.Options)
			if err != nil {
				// codes.
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

			// keep the store for futur call...
			self.stores[c.Id] = s
		}

		// If the connection need to save in the server configuration.
		if rqst.Save == true {
			self.Connections[c.Id] = c
			// In that case I will save it in file.
			err = self.Save()
			if err != nil {
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			globular.UpdateServiceConfig(self)
		} else {
			self.connections[c.Id] = c
		}
	}

	// test if the connection is reacheable.
	err = self.stores[c.Id].Ping(ctx, c.Id)

	if err != nil {
		self.stores[c.Id].Disconnect(c.Id)
		if _, ok := self.connections[rqst.Connection.Id]; ok {
			delete(self.connections, rqst.Connection.Id)
		} else if _, ok := self.Connections[rqst.Connection.Id]; ok {
			delete(self.Connections, rqst.Connection.Id)
		}
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print the success message here.
	return &persistencepb.CreateConnectionRsp{
		Result: true,
	}, nil
}

func (self *server) Connect(ctx context.Context, rqst *persistencepb.ConnectRqst) (*persistencepb.ConnectRsp, error) {
	store := self.stores[rqst.GetConnectionId()]
	if store != nil {
		return &persistencepb.ConnectRsp{
			Result: true,
		}, nil
	}

	if c, ok := self.Connections[rqst.ConnectionId]; ok {
		// So here I will open the connection.
		c.password = rqst.Password
		if c.Store == persistencepb.StoreType_MONGO {
			// here I will create a new mongo data store.
			s := new(persistence_store.MongoStore)

			// Now I will try to connect...
			err := s.Connect(c.Id, c.Host, c.Port, c.User, c.password, c.Name, c.Timeout, c.Options)
			if err != nil {
				// codes.
				return nil, status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

			// keep the store for futur call...
			self.stores[c.Id] = s
		}

		// set or update the connection and save it in json file.
		self.Connections[c.Id] = c

		return &persistencepb.ConnectRsp{
			Result: true,
		}, nil
	} else {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

}

// Close connection.
func (self *server) Disconnect(ctx context.Context, rqst *persistencepb.DisconnectRqst) (*persistencepb.DisconnectRsp, error) {
	store := self.stores[rqst.GetConnectionId()]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.GetConnectionId())
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.Disconnect(rqst.GetConnectionId())
	if err != nil {
		// codes.
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.DisconnectRsp{
		Result: true,
	}, nil
}

// Create a database
func (self *server) CreateDatabase(ctx context.Context, rqst *persistencepb.CreateDatabaseRqst) (*persistencepb.CreateDatabaseRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("CreateDatabase No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.CreateDatabase(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.CreateDatabaseRsp{
		Result: true,
	}, nil
}

// Delete a database
func (self *server) DeleteDatabase(ctx context.Context, rqst *persistencepb.DeleteDatabaseRqst) (*persistencepb.DeleteDatabaseRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("DeleteDatabase No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.DeleteDatabase(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.DeleteDatabaseRsp{
		Result: true,
	}, nil
}

// Create a Collection
func (self *server) CreateCollection(ctx context.Context, rqst *persistencepb.CreateCollectionRqst) (*persistencepb.CreateCollectionRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("CreateCollection No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.CreateCollection(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.OptionsStr)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.CreateCollectionRsp{
		Result: true,
	}, nil
}

// Delete collection
func (self *server) DeleteCollection(ctx context.Context, rqst *persistencepb.DeleteCollectionRqst) (*persistencepb.DeleteCollectionRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("DeleteCollection No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.DeleteCollection(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.DeleteCollectionRsp{
		Result: true,
	}, nil
}

// Ping a sql connection.
func (self *server) Ping(ctx context.Context, rqst *persistencepb.PingConnectionRqst) (*persistencepb.PingConnectionRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("Ping No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.Ping(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.PingConnectionRsp{
		Result: "pong",
	}, nil
}

// Get the number of entry in a collection
func (self *server) Count(ctx context.Context, rqst *persistencepb.CountRqst) (*persistencepb.CountRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("Count No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	count, err := store.Count(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.CountRsp{
		Result: count,
	}, nil
}

// Implementation of the Persistence method.
func (self *server) InsertOne(ctx context.Context, rqst *persistencepb.InsertOneRqst) (*persistencepb.InsertOneRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("InsertOne No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// In that case I will save it in file.
	err := self.Save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	entity := make(map[string]interface{})
	err = json.Unmarshal([]byte(rqst.JsonStr), &entity)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	var id interface{}
	id, err = store.InsertOne(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, entity, rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	jsonStr, err := json.Marshal(id)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.InsertOneRsp{
		Id: string(jsonStr),
	}, nil
}

func (self *server) InsertMany(stream persistencepb.PersistenceService_InsertManyServer) error {
	ids := make([]interface{}, 0)

	// In that case I will save it in file.
	err := self.Save()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	for {
		rqst, err := stream.Recv()

		// end of the stream.
		if err == io.EOF {
			jsonStr, err := json.Marshal(ids)
			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}
			// Close the stream...
			stream.SendAndClose(&persistencepb.InsertManyRsp{
				Ids: string(jsonStr),
			})

			return nil
		}

		entities := make([]interface{}, 0)
		err = json.Unmarshal([]byte(rqst.JsonStr), &entities)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		var results []interface{}
		results, err = self.stores[rqst.Id].InsertMany(stream.Context(), strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, entities, rqst.Options)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		// append to the list of ids.
		ids = append(ids, results...)

	}
}

// Find many
func (self *server) Find(rqst *persistencepb.FindRqst, stream persistencepb.PersistenceService_FindServer) error {

	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]

	if store == nil {
		err := errors.New("Find No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	results, err := store.Find(stream.Context(), strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Options)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 100
	values := make([]interface{}, 0)
	for i := 0; i < len(results); i++ {
		values = append(values, results[i])
		if len(values) >= maxSize {
			jsonStr, err := json.Marshal(values)
			if err != nil {
				return err
			}
			stream.Send(
				&persistencepb.FindResp{
					JsonStr: string(jsonStr),
				},
			)
			values = make([]interface{}, 0)
		}
	}

	// Send reminding values.
	if len(values) > 0 {
		jsonStr, err := json.Marshal(values)
		if err != nil {
			return err
		}
		stream.Send(
			&persistencepb.FindResp{

				JsonStr: string(jsonStr),
			},
		)
		values = make([]interface{}, 0)
	}

	return nil
}

func (self *server) Aggregate(rqst *persistencepb.AggregateRqst, stream persistencepb.PersistenceService_AggregateServer) error {

	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]

	if store == nil {
		err := errors.New("Aggregate No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	results, err := store.Aggregate(stream.Context(), strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Pipeline, rqst.Options)
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// No I will stream the result over the networks.
	maxSize := 100
	values := make([]interface{}, 0)
	for i := 0; i < len(results); i++ {
		values = append(values, results[i])
		if len(values) >= maxSize {
			jsonStr, err := json.Marshal(values)
			if err != nil {
				return err
			}
			stream.Send(
				&persistencepb.AggregateResp{
					JsonStr: string(jsonStr),
				},
			)
			values = make([]interface{}, 0)
		}
	}

	// Send reminding values.
	if len(values) > 0 {
		jsonStr, err := json.Marshal(values)
		if err != nil {
			return err
		}
		stream.Send(
			&persistencepb.AggregateResp{

				JsonStr: string(jsonStr),
			},
		)
		values = make([]interface{}, 0)
	}

	return nil
}

// Find one
func (self *server) FindOne(ctx context.Context, rqst *persistencepb.FindOneRqst) (*persistencepb.FindOneResp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("FindOne No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	result, err := store.FindOne(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Options)
	if err != nil {
		err = errors.New(rqst.Collection + " " + rqst.Query + " " + err.Error())
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	jsonStr, err := json.Marshal(result)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.FindOneResp{
		JsonStr: string(jsonStr),
	}, nil
}

// Update a single or many value depending of the query
func (self *server) Update(ctx context.Context, rqst *persistencepb.UpdateRqst) (*persistencepb.UpdateRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("Update No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.Update(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Value, rqst.Options)
	if err != nil {
		return nil, err
	}

	return &persistencepb.UpdateRsp{
		Result: true,
	}, nil
}

// Update a single docuemnt value(s)
func (self *server) UpdateOne(ctx context.Context, rqst *persistencepb.UpdateOneRqst) (*persistencepb.UpdateOneRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("UpdateOne No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.UpdateOne(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Value, rqst.Options)
	if err != nil {
		return nil, err
	}

	return &persistencepb.UpdateOneRsp{
		Result: true,
	}, nil
}

// Replace one document by another.
func (self *server) ReplaceOne(ctx context.Context, rqst *persistencepb.ReplaceOneRqst) (*persistencepb.ReplaceOneRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("ReplaceOne No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.ReplaceOne(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Value, rqst.Options)
	if err != nil {
		return nil, err
	}

	return &persistencepb.ReplaceOneRsp{
		Result: true,
	}, nil
}

// Delete many or one.
func (self *server) Delete(ctx context.Context, rqst *persistencepb.DeleteRqst) (*persistencepb.DeleteRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("Delete No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.Delete(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Options)
	if err != nil {
		return nil, err
	}

	return &persistencepb.DeleteRsp{
		Result: true,
	}, nil
}

// Delete one document at time
func (self *server) DeleteOne(ctx context.Context, rqst *persistencepb.DeleteOneRqst) (*persistencepb.DeleteOneRsp, error) {
	store := self.stores[strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")]
	if store == nil {
		err := errors.New("DeleteOne No store connection exist for id " + strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"))
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.DeleteOne(ctx, strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_"), strings.ReplaceAll(strings.ReplaceAll(rqst.Database, "@", "_"), ".", "_"), rqst.Collection, rqst.Query, rqst.Options)
	if err != nil {
		return nil, err
	}

	return &persistencepb.DeleteOneRsp{
		Result: true,
	}, nil
}

// Remove a connection from the map and the file.
func (self *server) DeleteConnection(ctx context.Context, rqst *persistencepb.DeleteConnectionRqst) (*persistencepb.DeleteConnectionRsp, error) {

	id := strings.ReplaceAll(strings.ReplaceAll(rqst.Id, "@", "_"), ".", "_")
	if _, ok := self.Connections[id]; !ok {
		return &persistencepb.DeleteConnectionRsp{
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
	return &persistencepb.DeleteConnectionRsp{
		Result: true,
	}, nil
}

// Create a new user.
func (self *server) RunAdminCmd(ctx context.Context, rqst *persistencepb.RunAdminCmdRqst) (*persistencepb.RunAdminCmdRsp, error) {
	store := self.stores[rqst.GetConnectionId()]
	if store == nil {
		err := errors.New("RunAdminCmd No store connection exist for id " + rqst.GetConnectionId())
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.RunAdminCmd(ctx, rqst.GetConnectionId(), rqst.User, rqst.Password, rqst.Script)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.RunAdminCmdRsp{
		Result: "",
	}, nil
}

func (self *server) Stop(context.Context, *persistencepb.StopRequest) (*persistencepb.StopResponse, error) {
	return &persistencepb.StopResponse{}, self.StopService()
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, "persistence_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(persistencepb.File_services_proto_persistence_proto.Services().Get(0).FullName())
	s_impl.Port = port
	s_impl.Proto = persistencepb.File_services_proto_persistence_proto.Path()
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
	persistencepb.RegisterPersistenceServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)

	// Start the service.
	s_impl.StartService()

}
