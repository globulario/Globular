package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/Globular/services/golang/catalog/catalog_client"
	"github.com/globulario/Globular/services/golang/catalog/catalogpb"
	"github.com/globulario/Globular/services/golang/event/event_client"
	globular "github.com/globulario/Globular/services/golang/globular_service"
	"github.com/globulario/Globular/services/golang/persistence/persistence_client"
	"github.com/davecourtois/Utility"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10017
	defaultProxy = 10018

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default address.
	domain string = "localhost"
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id              string
	Name            string
	Port            int
	Proxy           int
	Path            string
	Proto           string
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.
	Protocol        string
	Domain          string
	Description     string
	Keywords        []string
	Repositories    []string
	Discoveries     []string

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

	// Contain the list of service use by the catalog server.
	Services    map[string]interface{}
	Permissions []interface{}

	// Here I will create client to services use by the catalog server.
	persistenceClient *persistence_client.Persistence_Client
	eventClient       *event_client.Event_Client
	// The grpc server.
	grpcServer *grpc.Server
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

// The description of the service
func (self *server) GetDescription() string {
	return self.Description
}
func (self *server) SetDescription(description string) {
	self.Description = description
}

// The list of keywords of the services.
func (self *server) GetKeywords() []string {
	return self.Keywords
}
func (self *server) SetKeywords(keywords []string) {
	self.Keywords = keywords
}

// Dist
func (self *server) Dist(path string) error {

	return globular.Dist(path, self)
}

func (self *server) GetPlatform() string {
	return globular.GetPlatform()
}

func (self *server) PublishService(address string, user string, password string) error {
	return globular.PublishService(address, user, password, self)
}

func (self *server) GetRepositories() []string {
	return self.Repositories
}
func (self *server) SetRepositories(repositories []string) {
	self.Repositories = repositories
}

func (self *server) GetDiscoveries() []string {
	return self.Discoveries
}
func (self *server) SetDiscoveries(discoveries []string) {
	self.Discoveries = discoveries
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
	Utility.RegisterFunction("NewCatalogService_Client", catalog_client.NewCatalogService_Client)

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

	self.Services = make(map[string]interface{}, 0)

	// Connect to the persistence service.
	if self.Services["Persistence"] != nil {
		persistence_service := self.Services["Persistence"].(map[string]interface{})
		domain := persistence_service["Domain"].(string)
		self.persistenceClient, err = persistence_client.NewPersistenceService_Client(domain, "persistence.PersistenceService")
		if err != nil {
			log.Println("fail to connect to persistence service ", err)
		}
	} else {
		self.Services["Persistence"] = make(map[string]interface{})
		self.Services["Persistence"].(map[string]interface{})["Domain"] = "localhost"
		self.persistenceClient, err = persistence_client.NewPersistenceService_Client("localhost", "persistence.PersistenceService")
		if err != nil {
			log.Println("fail to connect to persistence service ", err)
		}
	}

	if self.Services["Event"] != nil {
		event_service := self.Services["Event"].(map[string]interface{})
		domain := event_service["Domain"].(string)
		self.eventClient, err = event_client.NewEventService_Client(domain, "event.EventService")
		if err != nil {
			log.Println("fail to connect to event service ", err)
		}
	} else {
		self.Services["Event"] = make(map[string]interface{})
		self.Services["Event"].(map[string]interface{})["Domain"] = "localhost"
		self.eventClient, err = event_client.NewEventService_Client("localhost", "event.EventService")
		if err != nil {
			log.Println("fail to connect to event service ", err)
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

func (self *server) Stop(context.Context, *catalogpb.StopRequest) (*catalogpb.StopResponse, error) {
	return &catalogpb.StopResponse{}, self.StopService()
}

// Create a new connection.
func (self *server) CreateConnection(ctx context.Context, rqst *catalogpb.CreateConnectionRqst) (*catalogpb.CreateConnectionRsp, error) {
	if rqst.Connection == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection information found in the request!")))

	}
	// So here I will call the function on the client.
	if self.persistenceClient != nil {
		persistence := self.Services["Persistence"].(map[string]interface{})
		if persistence["Connections"] == nil {
			persistence["Connections"] = make(map[string]interface{}, 0)
		}

		connections := persistence["Connections"].(map[string]interface{})

		storeType := int32(rqst.Connection.GetStore())
		err := self.persistenceClient.CreateConnection(rqst.Connection.GetId(), rqst.Connection.GetName(), rqst.Connection.GetHost(), Utility.ToNumeric(rqst.Connection.Port), Utility.ToNumeric(storeType), rqst.Connection.GetUser(), rqst.Connection.GetPassword(), Utility.ToNumeric(rqst.Connection.GetTimeout()), rqst.Connection.GetOptions(), true)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		connection := make(map[string]interface{}, 0)
		connection["Id"] = rqst.Connection.GetId()
		connection["Name"] = rqst.Connection.GetName()
		connection["Host"] = rqst.Connection.GetHost()
		connection["Store"] = rqst.Connection.GetStore()
		connection["User"] = rqst.Connection.GetUser()
		connection["Password"] = rqst.Connection.GetPassword()
		connection["Port"] = rqst.Connection.GetPort()
		connection["Timeout"] = rqst.Connection.GetTimeout()
		connection["Options"] = rqst.Connection.GetOptions()

		connections[rqst.Connection.GetId()] = connection

		self.Save()

	}

	return &catalogpb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Delete a connection.
func (self *server) DeleteConnection(ctx context.Context, rqst *catalogpb.DeleteConnectionRqst) (*catalogpb.DeleteConnectionRsp, error) {
	return nil, nil
}

// Create unit of measure exemple inch
func (self *server) SaveUnitOfMeasure(ctx context.Context, rqst *catalogpb.SaveUnitOfMeasureRequest) (*catalogpb.SaveUnitOfMeasureResponse, error) {
	unitOfMeasure := rqst.GetUnitOfMeasure()

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(unitOfMeasure.Id + unitOfMeasure.LanguageCode)
	self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "UnitOfMeasure", `{ "_id" : "`+_id+`" }`, "")

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(unitOfMeasure)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	id, err := self.persistenceClient.InsertOne(connection["Id"].(string), connection["Name"].(string), "UnitOfMeasure", jsonStr, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SaveUnitOfMeasureResponse{
		Id: id,
	}, nil
}

// Create property definition return the uuid of the created property
func (self *server) SavePropertyDefinition(ctx context.Context, rqst *catalogpb.SavePropertyDefinitionRequest) (*catalogpb.SavePropertyDefinitionResponse, error) {
	propertyDefinition := rqst.PropertyDefinition

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(propertyDefinition.Id + propertyDefinition.LanguageCode)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(propertyDefinition)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "PropertyDefinition", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SavePropertyDefinitionResponse{
		Id: _id,
	}, nil
}

// Create item definition.
func (self *server) SaveItemDefinition(ctx context.Context, rqst *catalogpb.SaveItemDefinitionRequest) (*catalogpb.SaveItemDefinitionResponse, error) {
	itemDefinition := rqst.ItemDefinition

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(itemDefinition.Id + itemDefinition.LanguageCode)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(itemDefinition)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Set the db reference.
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "ItemDefinition", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SaveItemDefinitionResponse{
		Id: _id,
	}, nil
}

// Create item response request.
func (self *server) SaveInventory(ctx context.Context, rqst *catalogpb.SaveInventoryRequest) (*catalogpb.SaveInventoryResponse, error) {
	inventory := rqst.Inventory
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(inventory)

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(inventory.LocalisationId + inventory.PacakgeId)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Set the db reference.
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Inventory", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SaveInventoryResponse{
		Id: _id,
	}, nil
}

// Create item response request.
func (self *server) SaveItemInstance(ctx context.Context, rqst *catalogpb.SaveItemInstanceRequest) (*catalogpb.SaveItemInstanceResponse, error) {
	instance := rqst.ItemInstance
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	var marshaler jsonpb.Marshaler

	jsonStr, err := marshaler.MarshalToString(instance)
	if len(instance.Id) == 0 {
		instance.Id = Utility.RandomUUID()
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(instance.Id)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Set the db reference.
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "ItemInstance", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SaveItemInstanceResponse{
		Id: _id,
	}, nil
}

// Save Manufacturer
func (self *server) SaveManufacturer(ctx context.Context, rqst *catalogpb.SaveManufacturerRequest) (*catalogpb.SaveManufacturerResponse, error) {

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	manufacturer := rqst.Manufacturer

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(manufacturer.Id)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(manufacturer)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Manufacturer", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append manafacturer item response.
	return &catalogpb.SaveManufacturerResponse{
		Id: _id,
	}, nil
}

// Save Supplier
func (self *server) SaveSupplier(ctx context.Context, rqst *catalogpb.SaveSupplierRequest) (*catalogpb.SaveSupplierResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	supplier := rqst.Supplier

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(supplier.Id)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(supplier)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Supplier", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append manafacturer item response.
	return &catalogpb.SaveSupplierResponse{
		Id: _id,
	}, nil
}

// Save localisation
func (self *server) SaveLocalisation(ctx context.Context, rqst *catalogpb.SaveLocalisationRequest) (*catalogpb.SaveLocalisationResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	localisation := rqst.Localisation

	// Set the reference correctly.
	if localisation.SubLocalisations != nil {
		for i := 0; i < len(localisation.SubLocalisations.Values); i++ {
			if !Utility.IsUuid(localisation.SubLocalisations.Values[i].GetRefObjId()) {
				localisation.SubLocalisations.Values[i].RefObjId = Utility.GenerateUUID(localisation.SubLocalisations.Values[i].GetRefObjId())
			}
		}
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(localisation.Id + localisation.LanguageCode)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(localisation)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// set the object references...
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Localisation", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append manafacturer item response.
	return &catalogpb.SaveLocalisationResponse{
		Id: _id,
	}, nil
}

// Save Package
func (self *server) SavePackage(ctx context.Context, rqst *catalogpb.SavePackageRequest) (*catalogpb.SavePackageResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	package_ := rqst.Package

	// Set the reference with _id value if not already set.
	for i := 0; i < len(package_.Subpackages); i++ {
		subPackage := package_.Subpackages[i]
		if subPackage.UnitOfMeasure != nil {
			if !Utility.IsUuid(subPackage.UnitOfMeasure.RefObjId) {
				subPackage.UnitOfMeasure.RefObjId = Utility.GenerateUUID(subPackage.UnitOfMeasure.RefObjId)
			}
		}
		if subPackage.Package != nil {
			if !Utility.IsUuid(subPackage.Package.RefObjId) {
				subPackage.Package.RefObjId = Utility.GenerateUUID(subPackage.Package.RefObjId)
			}
		}
	}

	for i := 0; i < len(package_.ItemInstances); i++ {
		itemInstance := package_.ItemInstances[i]
		if itemInstance.UnitOfMeasure != nil {
			if !Utility.IsUuid(itemInstance.UnitOfMeasure.RefObjId) {
				itemInstance.UnitOfMeasure.RefObjId = Utility.GenerateUUID(itemInstance.UnitOfMeasure.RefObjId)
			}
		}
		if !Utility.IsUuid(itemInstance.ItemInstance.RefObjId) {
			itemInstance.ItemInstance.RefObjId = Utility.GenerateUUID(itemInstance.ItemInstance.RefObjId)
		}
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(package_.Id + package_.LanguageCode)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(package_)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// set the object references...
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Package", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append manafacturer item response.
	return &catalogpb.SavePackageResponse{
		Id: _id,
	}, nil

}

// Save Package Supplier
func (self *server) SavePackageSupplier(ctx context.Context, rqst *catalogpb.SavePackageSupplierRequest) (*catalogpb.SavePackageSupplierResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	// Now I will create the PackageSupplier.
	// save the supplier id
	packageSupplier := rqst.PackageSupplier

	if !Utility.IsUuid(packageSupplier.Supplier.RefObjId) {
		packageSupplier.Supplier.RefObjId = Utility.GenerateUUID(packageSupplier.Supplier.RefObjId)
	}

	if !Utility.IsUuid(packageSupplier.Package.RefObjId) {
		packageSupplier.Package.RefObjId = Utility.GenerateUUID(packageSupplier.Package.RefObjId)
	}

	// Test if the pacakge exist
	_, err := self.persistenceClient.FindOne(connection["Id"].(string), rqst.PackageSupplier.Package.RefDbName, rqst.PackageSupplier.Package.RefColId, `{"_id":"`+rqst.PackageSupplier.Package.RefObjId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Test if the supplier exist.
	_, err = self.persistenceClient.FindOne(connection["Id"].(string), rqst.PackageSupplier.Supplier.RefDbName, rqst.PackageSupplier.Supplier.RefColId, `{"_id":"`+rqst.PackageSupplier.Supplier.RefObjId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(packageSupplier.Id)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(packageSupplier)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// set the object references...
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "PackageSupplier", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SavePackageSupplierResponse{
		Id: _id,
	}, nil
}

// Save Item Manufacturer
func (self *server) SaveItemManufacturer(ctx context.Context, rqst *catalogpb.SaveItemManufacturerRequest) (*catalogpb.SaveItemManufacturerResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// Test if the item exist
	_, err := self.persistenceClient.FindOne(connection["Id"].(string), rqst.ItemManafacturer.Item.RefDbName, rqst.ItemManafacturer.Item.RefColId, `{"_id":"`+rqst.ItemManafacturer.Item.RefObjId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Test if the supplier exist.
	_, err = self.persistenceClient.FindOne(connection["Id"].(string), rqst.ItemManafacturer.Manufacturer.RefDbName, rqst.ItemManafacturer.Manufacturer.RefColId, `{"_id":"`+rqst.ItemManafacturer.Manufacturer.RefObjId+`"}`, "")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Now I will create the PackageSupplier.
	// save the supplier id
	itemManafacturer := rqst.ItemManafacturer

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(itemManafacturer.Id)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(itemManafacturer)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// set the object references...
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "ItemManufacturer", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.SaveItemManufacturerResponse{
		Id: _id,
	}, nil
}

// Save Item Category
func (self *server) SaveCategory(ctx context.Context, rqst *catalogpb.SaveCategoryRequest) (*catalogpb.SaveCategoryResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	category := rqst.Category

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(category.Id + category.LanguageCode)

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(category)
	jsonStr = `{ "_id" : "` + _id + `",` + jsonStr[1:]

	// Always create a new
	err = self.persistenceClient.ReplaceOne(connection["Id"].(string), connection["Name"].(string), "Category", `{ "_id" : "`+_id+`"}`, jsonStr, `[{"upsert": true}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Append manafacturer item response.
	return &catalogpb.SaveCategoryResponse{
		Id: _id,
	}, nil
}

// Append a new Item to manufacturer
func (self *server) AppendItemDefinitionCategory(ctx context.Context, rqst *catalogpb.AppendItemDefinitionCategoryRequest) (*catalogpb.AppendItemDefinitionCategoryResponse, error) {

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(rqst.Category)

	// Set the db reference.
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Now I will modify the jsonStr to insert the value in the array.
	jsonStr = `{ "$push": { "categories":` + jsonStr + `}}`

	// Always create a new
	err = self.persistenceClient.UpdateOne(connection["Id"].(string), connection["Name"].(string), rqst.ItemDefinition.RefColId, `{ "_id" : "`+rqst.ItemDefinition.RefObjId+`"}`, jsonStr, `[]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.AppendItemDefinitionCategoryResponse{
		Result: true,
	}, nil
}

// Remove Item from manufacturer
func (self *server) RemoveItemDefinitionCategory(ctx context.Context, rqst *catalogpb.RemoveItemDefinitionCategoryRequest) (*catalogpb.RemoveItemDefinitionCategoryResponse, error) {

	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	var marshaler jsonpb.Marshaler
	jsonStr, err := marshaler.MarshalToString(rqst.Category)

	// Set the db reference.
	jsonStr = strings.Replace(jsonStr, "refObjId", "$id", -1)
	jsonStr = strings.Replace(jsonStr, "refColId", "$ref", -1)
	jsonStr = strings.Replace(jsonStr, "refDbName", "$db", -1)

	// Now I will modify the jsonStr to insert the value in the array.
	jsonStr = `{ "$pull": { "categories":` + jsonStr + `}}` // remove a particular item.

	// Always create a new
	err = self.persistenceClient.UpdateOne(connection["Id"].(string), connection["Name"].(string), rqst.ItemDefinition.RefColId, `{ "_id" : "`+rqst.ItemDefinition.RefObjId+`"}`, jsonStr, `[]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.RemoveItemDefinitionCategoryResponse{
		Result: true,
	}, nil
}

//////////////////////////////// Getter function ///////////////////////////////

// Getter function.

// Getter Item instance.
func (self *server) GetItemInstance(ctx context.Context, rqst *catalogpb.GetItemInstanceRequest) (*catalogpb.GetItemInstanceResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.ItemInstanceId) {
		query = `{"_id":"` + rqst.ItemInstanceId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.ItemInstanceId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "ItemInstance", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	instance := new(catalogpb.ItemInstance)
	err = jsonpb.UnmarshalString(jsonStr, instance)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetItemInstanceResponse{
		ItemInstance: instance,
	}, nil

}

// Get Item Instances.
func (self *server) GetItemInstances(ctx context.Context, rqst *catalogpb.GetItemInstancesRequest) (*catalogpb.GetItemInstancesResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "ItemInstance", rqst.Query, options)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	// set the Packages properties...
	jsonStr = `{ "itemInstances":` + jsonStr + `}`

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// unmarshall object
	instances := new(catalogpb.ItemInstances)
	err = jsonpb.UnmarshalString(string(jsonStr), instances)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetItemInstancesResponse{
		ItemInstances: instances.ItemInstances,
	}, nil
}

// Getter Item defintion.
func (self *server) GetItemDefinition(ctx context.Context, rqst *catalogpb.GetItemDefinitionRequest) (*catalogpb.GetItemDefinitionResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.ItemDefinitionId) {
		query = `{"_id":"` + rqst.ItemDefinitionId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.ItemDefinitionId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "ItemDefinition", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	definition := new(catalogpb.ItemDefinition)
	err = jsonpb.UnmarshalString(jsonStr, definition)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetItemDefinitionResponse{
		ItemDefinition: definition,
	}, nil

}

// Get Inventories.
func (self *server) GetInventories(ctx context.Context, rqst *catalogpb.GetInventoriesRequest) (*catalogpb.GetInventoriesResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}
	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Inventory", rqst.Query, options)

	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	// set the Packages properties...
	jsonStr = `{ "inventories":` + jsonStr + `}`

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// unmarshall object
	inventories := new(catalogpb.Inventories)
	err = jsonpb.UnmarshalString(string(jsonStr), inventories)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetInventoriesResponse{
		Inventories: inventories.Inventories,
	}, nil
}

// Get Item Definitions.
func (self *server) GetItemDefinitions(ctx context.Context, rqst *catalogpb.GetItemDefinitionsRequest) (*catalogpb.GetItemDefinitionsResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "ItemDefinition", rqst.Query, options)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	// set the Packages properties...
	jsonStr = `{ "itemDefinitions":` + jsonStr + `}`

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// unmarshall object
	definitions := new(catalogpb.ItemDefinitions)
	err = jsonpb.UnmarshalString(string(jsonStr), definitions)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetItemDefinitionsResponse{
		ItemDefinitions: definitions.ItemDefinitions,
	}, nil
}

// Getter Supplier.
func (self *server) GetSupplier(ctx context.Context, rqst *catalogpb.GetSupplierRequest) (*catalogpb.GetSupplierResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.SupplierId) {
		query = `{"_id":"` + rqst.SupplierId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.SupplierId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "Supplier", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	supplier := new(catalogpb.Supplier)
	err = jsonpb.UnmarshalString(jsonStr, supplier)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetSupplierResponse{
		Supplier: supplier,
	}, nil

}

// Get Suppliers
func (self *server) GetSuppliers(ctx context.Context, rqst *catalogpb.GetSuppliersRequest) (*catalogpb.GetSuppliersResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}
	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Supplier", rqst.Query, options)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	// set the Packages properties...
	jsonStr = `{ "suppliers":` + jsonStr + `}`

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// unmarshall object
	suppliers := new(catalogpb.Suppliers)
	err = jsonpb.UnmarshalString(string(jsonStr), suppliers)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetSuppliersResponse{
		Suppliers: suppliers.Suppliers,
	}, nil
}

// Get Package supplier.
func (self *server) GetSupplierPackages(ctx context.Context, rqst *catalogpb.GetSupplierPackagesRequest) (*catalogpb.GetSupplierPackagesResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.SupplierId) {
		query = `{"supplier.$id":"` + rqst.SupplierId + `"}`
	} else {
		query = `{"supplier.$id":"` + Utility.GenerateUUID(rqst.SupplierId) + `"}`
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "PackageSupplier", query, `[{"Projection":{"_id":1}}]`)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	results := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(jsonStr), &results)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	packagesSupplier := make([]*catalogpb.PackageSupplier, 0)

	for i := 0; i < len(results); i++ {
		jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "PackageSupplier", `{"_id":"`+results[i]["_id"].(string)+`"}`, `[{"Projection":{"_id":0}}]`)
		// replace the reference.
		jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
		jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
		jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		packageSupplier := new(catalogpb.PackageSupplier)
		err = jsonpb.UnmarshalString(jsonStr, packageSupplier)
		if err != nil {
			log.Println(err)
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		packagesSupplier = append(packagesSupplier, packageSupplier)
	}

	return &catalogpb.GetSupplierPackagesResponse{
		PackagesSupplier: packagesSupplier,
	}, nil

}

// Get Package
func (self *server) GetPackage(ctx context.Context, rqst *catalogpb.GetPackageRequest) (*catalogpb.GetPackageResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.PackageId) {
		query = `{"_id":"` + rqst.PackageId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.PackageId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "Package", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	package_ := new(catalogpb.Package)
	err = jsonpb.UnmarshalString(jsonStr, package_)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetPackageResponse{
		Pacakge: package_,
	}, nil

}

// Get Packages
func (self *server) GetPackages(ctx context.Context, rqst *catalogpb.GetPackagesRequest) (*catalogpb.GetPackagesResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Package", rqst.Query, options)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// set the Packages properties...
	jsonStr = `{ "packages":` + jsonStr + `}`

	// unmarshall object
	packages := new(catalogpb.Packages)
	err = jsonpb.UnmarshalString(jsonStr, packages)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetPackagesResponse{
		Packages: packages.Packages,
	}, nil

}

func (self *server) getLocalisation(localisationId string, connectionId string) (*catalogpb.Localisation, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[connectionId] == nil {
		return nil, errors.New("no connection found with id " + connectionId)
	}

	connection := persistence["Connections"].(map[string]interface{})[connectionId].(map[string]interface{})

	var query string
	if Utility.IsUuid(localisationId) {
		query = `{"_id":"` + localisationId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(localisationId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "Localisation", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, err
	}

	localisation := new(catalogpb.Localisation)
	err = jsonpb.UnmarshalString(jsonStr, localisation)
	if err != nil {
		return nil, err
	}

	return localisation, nil
}

// Get Localisation
func (self *server) GetLocalisation(ctx context.Context, rqst *catalogpb.GetLocalisationRequest) (*catalogpb.GetLocalisationResponse, error) {

	localisation, err := self.getLocalisation(rqst.LocalisationId, rqst.ConnectionId)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetLocalisationResponse{
		Localisation: localisation,
	}, nil

}

func (self *server) getLocalisations(query string, options string, connectionId string) ([]*catalogpb.Localisation, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[connectionId] == nil {
		return nil, errors.New("no connection found with id " + connectionId)
	}

	connection := persistence["Connections"].(map[string]interface{})[connectionId].(map[string]interface{})
	options, err := self.getOptionsString(options)
	if err != nil {
		return nil, err
	}

	if len(query) == 0 {
		query = `{}`
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Localisation", query, options)

	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, err
	}

	// set the Packages properties...
	jsonStr = `{ "localisations":` + jsonStr + `}`

	// unmarshall object
	localisations := new(catalogpb.Localisations)
	err = jsonpb.UnmarshalString(jsonStr, localisations)

	if err != nil {
		return nil, err
	}

	return localisations.Localisations, nil

}

// Get Packages
func (self *server) GetLocalisations(ctx context.Context, rqst *catalogpb.GetLocalisationsRequest) (*catalogpb.GetLocalisationsResponse, error) {

	// unmarshall object
	localisations, err := self.getLocalisations(rqst.Query, rqst.Options, rqst.ConnectionId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetLocalisationsResponse{
		Localisations: localisations,
	}, nil

}

/**
 * Get the category.
 */
func (self *server) getCategory(categoryId string, connectionId string) (*catalogpb.Category, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[connectionId] == nil {
		return nil, errors.New("no connection found with id " + connectionId)
	}

	connection := persistence["Connections"].(map[string]interface{})[connectionId].(map[string]interface{})

	var query string
	if Utility.IsUuid(categoryId) {
		query = `{"_id":"` + categoryId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(categoryId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "Category", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, err
	}

	category := new(catalogpb.Category)
	err = jsonpb.UnmarshalString(jsonStr, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (self *server) GetCategory(ctx context.Context, rqst *catalogpb.GetCategoryRequest) (*catalogpb.GetCategoryResponse, error) {

	category, err := self.getCategory(rqst.CategoryId, rqst.ConnectionId)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetCategoryResponse{
		Category: category,
	}, nil

}

func (self *server) getCategories(query string, options string, connectionId string) ([]*catalogpb.Category, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[connectionId] == nil {
		return nil, errors.New("no connection found with id " + connectionId)
	}

	connection := persistence["Connections"].(map[string]interface{})[connectionId].(map[string]interface{})
	options, err := self.getOptionsString(options)
	if err != nil {
		return nil, err
	}

	if len(query) == 0 {
		query = `{}`
	}

	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Category", query, options)

	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, err
	}

	// set the Packages properties...
	jsonStr = `{ "categories":` + jsonStr + `}`

	// unmarshall object
	categories := new(catalogpb.Categories)
	err = jsonpb.UnmarshalString(jsonStr, categories)

	if err != nil {
		return nil, err
	}

	return categories.Categories, nil

}

// Get Packages
func (self *server) GetCategories(ctx context.Context, rqst *catalogpb.GetCategoriesRequest) (*catalogpb.GetCategoriesResponse, error) {

	// unmarshall object
	categories, err := self.getCategories(rqst.Query, rqst.Options, rqst.ConnectionId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetCategoriesResponse{
		Categories: categories,
	}, nil

}

// Get Manufacturer
func (self *server) GetManufacturer(ctx context.Context, rqst *catalogpb.GetManufacturerRequest) (*catalogpb.GetManufacturerResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.ManufacturerId) {
		query = `{"_id":"` + rqst.ManufacturerId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.ManufacturerId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "Package", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	manufacturer := new(catalogpb.Manufacturer)
	err = jsonpb.UnmarshalString(jsonStr, manufacturer)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetManufacturerResponse{
		Manufacturer: manufacturer,
	}, nil

}

func (self *server) getOptionsString(options string) (string, error) {
	options_ := make([]map[string]interface{}, 0)
	if len(options) > 0 {
		err := json.Unmarshal([]byte(options), &options_)
		if err != nil {
			return "", err
		}

		var projections map[string]interface{}

		for i := 0; i < len(options_); i++ {
			if options_[i]["Projection"] != nil {
				projections = options_[i]["Projection"].(map[string]interface{})
				break
			}
		}

		if projections != nil {
			projections["_id"] = 0
		} else {
			options_ = append(options_, map[string]interface{}{"Projection": map[string]interface{}{"_id": 0}})
		}

	} else {
		options_ = append(options_, map[string]interface{}{"Projection": map[string]interface{}{"_id": 0}})
	}

	optionsStr, err := json.Marshal(options_)
	return string(optionsStr), err
}

// Get Manufacturers
func (self *server) GetManufacturers(ctx context.Context, rqst *catalogpb.GetManufacturersRequest) (*catalogpb.GetManufacturersResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}
	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "Manufacturer", rqst.Query, options)

	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	// set the Packages properties...
	jsonStr = `{ "manufacturers":` + jsonStr + `}`

	// unmarshall object
	manufacturers := new(catalogpb.Manufacturers)
	err = jsonpb.UnmarshalString(jsonStr, manufacturers)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetManufacturersResponse{
		Manufacturers: manufacturers.Manufacturers,
	}, nil

}

// Get Package
func (self *server) GetUnitOfMeasures(ctx context.Context, rqst *catalogpb.GetUnitOfMeasuresRequest) (*catalogpb.GetUnitOfMeasuresResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	options, err := self.getOptionsString(rqst.Options)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	if len(rqst.Query) == 0 {
		rqst.Query = `{}`
	}
	jsonStr, err := self.persistenceClient.Find(connection["Id"].(string), connection["Name"].(string), "UnitOfMeasure", rqst.Query, options)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	// set the Packages properties...
	jsonStr = `{ "unitOfMeasures":` + jsonStr + `}`

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// unmarshall object
	unitOfMeasures := new(catalogpb.UnitOfMeasures)
	err = jsonpb.UnmarshalString(string(jsonStr), unitOfMeasures)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetUnitOfMeasuresResponse{
		UnitOfMeasures: unitOfMeasures.UnitOfMeasures,
	}, nil

}

// Get Unit of measure.
func (self *server) GetUnitOfMeasure(ctx context.Context, rqst *catalogpb.GetUnitOfMeasureRequest) (*catalogpb.GetUnitOfMeasureResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})

	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})
	var query string
	if Utility.IsUuid(rqst.UnitOfMeasureId) {
		query = `{"_id":"` + rqst.UnitOfMeasureId + `"}`
	} else {
		query = `{"_id":"` + Utility.GenerateUUID(rqst.UnitOfMeasureId) + `"}`
	}

	jsonStr, err := self.persistenceClient.FindOne(connection["Id"].(string), connection["Name"].(string), "UnitOfMeasure", query, `[{"Projection":{"_id":0}}]`)
	// replace the reference.
	jsonStr = strings.Replace(jsonStr, "$id", "refObjId", -1)
	jsonStr = strings.Replace(jsonStr, "$ref", "refColId", -1)
	jsonStr = strings.Replace(jsonStr, "$db", "refDbName", -1)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	unitOfMeasure := new(catalogpb.UnitOfMeasure)
	err = jsonpb.UnmarshalString(jsonStr, unitOfMeasure)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.GetUnitOfMeasureResponse{
		UnitOfMeasure: unitOfMeasure,
	}, nil

}

////// Delete function //////

// Delete a package.
func (self *server) DeletePackage(ctx context.Context, rqst *catalogpb.DeletePackageRequest) (*catalogpb.DeletePackageResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})
	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	package_ := rqst.Package

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(package_.Id + package_.LanguageCode)

	err := self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "Package", `{"_id":"`+_id+`"}`, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.DeletePackageResponse{
		Result: true,
	}, nil
}

// Delete a package supplier
func (self *server) DeletePackageSupplier(ctx context.Context, rqst *catalogpb.DeletePackageSupplierRequest) (*catalogpb.DeletePackageSupplierResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})
	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	packageSupplier := rqst.PackageSupplier

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(packageSupplier.Id)

	err := self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "PackageSupplier", `{"_id":"`+_id+`"}`, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.DeletePackageSupplierResponse{
		Result: true,
	}, nil
}

// Delete a supplier
func (self *server) DeleteSupplier(ctx context.Context, rqst *catalogpb.DeleteSupplierRequest) (*catalogpb.DeleteSupplierResponse, error) {
	persistence := self.Services["Persistence"].(map[string]interface{})
	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	supplier := rqst.Supplier

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(supplier.Id)
	err := self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "Supplier", `{"_id":"`+_id+`"}`, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.DeleteSupplierResponse{
		Result: true,
	}, nil
}

// Delete propertie definition
func (self *server) DeletePropertyDefinition(ctx context.Context, rqst *catalogpb.DeletePropertyDefinitionRequest) (*catalogpb.DeletePropertyDefinitionResponse, error) {
	return nil, nil
}

// Delete unit of measure
func (self *server) DeleteUnitOfMeasure(ctx context.Context, rqst *catalogpb.DeleteUnitOfMeasureRequest) (*catalogpb.DeleteUnitOfMeasureResponse, error) {
	return nil, nil
}

// Delete Item Instance
func (self *server) DeleteItemInstance(ctx context.Context, rqst *catalogpb.DeleteItemInstanceRequest) (*catalogpb.DeleteItemInstanceResponse, error) {
	return nil, nil
}

// Delete Manufacturer
func (self *server) DeleteManufacturer(ctx context.Context, rqst *catalogpb.DeleteManufacturerRequest) (*catalogpb.DeleteManufacturerResponse, error) {
	return nil, nil
}

// Delete Item Manufacturer
func (self *server) DeleteItemManufacturer(ctx context.Context, rqst *catalogpb.DeleteItemManufacturerRequest) (*catalogpb.DeleteItemManufacturerResponse, error) {
	return nil, nil
}

// Delete Category
func (self *server) DeleteCategory(ctx context.Context, rqst *catalogpb.DeleteCategoryRequest) (*catalogpb.DeleteCategoryResponse, error) {
	return nil, nil
}

func (self *server) deleteLocalisation(localisation *catalogpb.Localisation, connectionId string) error {
	persistence := self.Services["Persistence"].(map[string]interface{})
	if persistence["Connections"].(map[string]interface{})[connectionId] == nil {
		return errors.New("no connection found with id " + connectionId)
	}

	connection := persistence["Connections"].(map[string]interface{})[connectionId].(map[string]interface{})

	// I will remove referencing object...
	referenced, err := self.getLocalisations(`{"subLocalisations.values.$id":"`+Utility.GenerateUUID(localisation.GetId()+localisation.GetLanguageCode())+`"}`, "", connectionId)
	if err == nil {
		refStr := `{"$id":"` + Utility.GenerateUUID(localisation.GetId()+localisation.GetLanguageCode()) + `","$ref":"Localisation","$db":"` + connection["Name"].(string) + `"}`
		for i := 0; i < len(referenced); i++ {
			// Now I will modify the jsonStr to insert the value in the array.
			query := `{"$pull":{"subLocalisations.values":` + refStr + `}}` // remove a particular item.
			_id := Utility.GenerateUUID(referenced[i].Id + referenced[i].LanguageCode)
			// Always create a new
			err = self.persistenceClient.UpdateOne(connection["Id"].(string), connection["Name"].(string), "Localisation", `{"_id" : "`+_id+`"}`, query, `[]`)
			if err != nil {
				return err
			}
		}
	}

	// So here I will delete all sub-localisation to...
	if localisation.GetSubLocalisations() != nil {
		for i := 0; i < len(localisation.GetSubLocalisations().GetValues()); i++ {
			subLocalisation, err := self.getLocalisation(localisation.GetSubLocalisations().GetValues()[i].GetRefObjId(), connectionId)
			if err == nil {
				err := self.deleteLocalisation(subLocalisation, connectionId)
				if err != nil {
					return err
				}
			}
		}
	}

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(localisation.Id + localisation.LanguageCode)

	return self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "Localisation", `{"_id":"`+_id+`"}`, "")

}

// Delete Localisation
func (self *server) DeleteLocalisation(ctx context.Context, rqst *catalogpb.DeleteLocalisationRequest) (*catalogpb.DeleteLocalisationResponse, error) {

	err := self.deleteLocalisation(rqst.Localisation, rqst.ConnectionId)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.DeleteLocalisationResponse{
		Result: true,
	}, nil
}

// Delete Localisation
func (self *server) DeleteInventory(ctx context.Context, rqst *catalogpb.DeleteInventoryRequest) (*catalogpb.DeleteInventoryResponse, error) {

	persistence := self.Services["Persistence"].(map[string]interface{})
	if persistence["Connections"].(map[string]interface{})[rqst.ConnectionId] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no connection found with id "+rqst.ConnectionId)))
	}

	connection := persistence["Connections"].(map[string]interface{})[rqst.ConnectionId].(map[string]interface{})

	// save the supplier id
	inventory := rqst.Inventory

	// Here I will generate the _id key
	_id := Utility.GenerateUUID(inventory.LocalisationId + inventory.PacakgeId)
	err := self.persistenceClient.DeleteOne(connection["Id"].(string), connection["Name"].(string), "Inventory", `{"_id":"`+_id+`"}`, "")

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &catalogpb.DeleteInventoryResponse{
		Result: true,
	}, nil

}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "catalog_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = Utility.GetExecName(os.Args[0])
	s_impl.Port = defaultPort
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.Keywords = make([]string, 0)
	s_impl.Repositories = make([]string, 0)
	s_impl.Discoveries = make([]string, 0)

	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		log.Fatalf("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
	}
	if len(os.Args) == 2 {
		s_impl.Port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	if len(os.Args) > 2 {
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommand_domain := publishCommand.String("a", "", "The address(domain ex. my.domain.com:8080) of your backend (Required)")
		publishCommand_user := publishCommand.String("u", "", "The user (Required)")
		publishCommand_password := publishCommand.String("p", "", "The password (Required)")

		switch os.Args[1] {
		case "publish":
			publishCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		if publishCommand.Parsed() {
			// Required Flags
			if *publishCommand_domain == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_user == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_password == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			err := s_impl.PublishService(*publishCommand_domain, *publishCommand_user, *publishCommand_password)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Your service was publish successfuly!")
			}
		}
	} else {
		// Register the echo services
		catalogpb.RegisterCatalogServiceServer(s_impl.grpcServer, s_impl)
		reflection.Register(s_impl.grpcServer)

		// Start the service.
		s_impl.StartService()
	}

}
