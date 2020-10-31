package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/globulario/Globular/Interceptors"
	globular "github.com/globulario/Globular/services/golang/globular_service"
	"github.com/globulario/Globular/services/golang/monitoring/monitoring_client"
	"github.com/globulario/Globular/services/golang/monitoring/monitoring_store"
	"github.com/globulario/Globular/services/golang/monitoring/monitoringpb"

	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	//"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10019
	defaultProxy = 10020

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	domain string = "localhost"
)

// Keep connection information here.
type connection struct {
	Id   string // The connection id
	Host string // can also be ipv4 addresse.
	Port int32
	Type monitoringpb.StoreType // Only Prometheus at this time.
}

// Value need by Globular to start the services...
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
	Permissions        []interface{} // contains the action permission for the services.
	// The grpc server.
	grpcServer *grpc.Server

	// That map contain the list of active connections.
	Connections map[string]connection
	stores      map[string]monitoring_store.Store
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

func (self *server) GetPlatform() string {
	return globular.GetPlatform()
}

func (self *server) PublishService(address string, user string, password string) error {
	return globular.PublishService(address, user, password, self)
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
	Utility.RegisterFunction("NewMonitoringService_Client", monitoring_client.NewMonitoringService_Client)

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

	self.stores = make(map[string]monitoring_store.Store)
	self.Connections = make(map[string]connection)

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

func (self *server) Stop(context.Context, *monitoringpb.StopRequest) (*monitoringpb.StopResponse, error) {
	return &monitoringpb.StopResponse{}, self.StopService()
}

///////////////////// Monitoring specific functions ////////////////////////////

// Create a connection.
func (self *server) CreateConnection(ctx context.Context, rqst *monitoringpb.CreateConnectionRqst) (*monitoringpb.CreateConnectionRsp, error) {
	var c connection

	// Set the connection info from the request.
	c.Id = rqst.Connection.Id
	c.Host = rqst.Connection.Host
	c.Port = rqst.Connection.Port
	c.Type = rqst.Connection.Store

	if self.Connections == nil {
		self.Connections = make(map[string]connection, 0)
	}

	// set or update the connection and save it in json file.
	self.Connections[c.Id] = c

	var store monitoring_store.Store
	var err error
	address := "http://" + c.Host + ":" + Utility.ToString(c.Port)

	if c.Type == monitoringpb.StoreType_PROMETHEUS {
		store, err = monitoring_store.NewPrometheusStore(address)
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Fail to connect to store!")))
	}

	// Keep the ref to the store.
	self.stores[c.Id] = store

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
	return &monitoringpb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Delete a connection.
func (self *server) DeleteConnection(ctx context.Context, rqst *monitoringpb.DeleteConnectionRqst) (*monitoringpb.DeleteConnectionRsp, error) {
	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return &monitoringpb.DeleteConnectionRsp{
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
	return &monitoringpb.DeleteConnectionRsp{
		Result: true,
	}, nil

}

// Alerts returns a list of all active alerts.
func (self *server) Alerts(ctx context.Context, rqst *monitoringpb.AlertsRequest) (*monitoringpb.AlertsResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	str, err := store.Alerts(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.AlertsResponse{
		Results: str,
	}, nil
}

// AlertManagers returns an overview of the current state of the Prometheus alert manager discovery.
func (self *server) AlertManagers(ctx context.Context, rqst *monitoringpb.AlertManagersRequest) (*monitoringpb.AlertManagersResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	str, err := store.AlertManagers(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.AlertManagersResponse{
		Results: str,
	}, nil
}

// CleanTombstones removes the deleted data from disk and cleans up the existing tombstones.
func (self *server) CleanTombstones(ctx context.Context, rqst *monitoringpb.CleanTombstonesRequest) (*monitoringpb.CleanTombstonesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err := store.CleanTombstones(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.CleanTombstonesResponse{
		Result: true,
	}, nil
}

// Config returns the current Prometheus configuration.
func (self *server) Config(ctx context.Context, rqst *monitoringpb.ConfigRequest) (*monitoringpb.ConfigResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	configStr, err := store.Config(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.ConfigResponse{
		Results: configStr,
	}, nil
}

// DeleteSeries deletes data for a selection of series in a time range.
func (self *server) DeleteSeries(ctx context.Context, rqst *monitoringpb.DeleteSeriesRequest) (*monitoringpb.DeleteSeriesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// convert input arguments...
	startTime := time.Unix(int64(rqst.GetStartTime()), 0)
	endTime := time.Unix(int64(rqst.GetEndTime()), 0)

	err := store.DeleteSeries(ctx, rqst.GetMatches(), startTime, endTime)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.DeleteSeriesResponse{
		Result: true,
	}, nil
}

// Flags returns the flag values that Prometheus was launched with.
func (self *server) Flags(ctx context.Context, rqst *monitoringpb.FlagsRequest) (*monitoringpb.FlagsResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	str, err := store.Flags(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.FlagsResponse{
		Results: str,
	}, nil
}

// LabelNames returns all the unique label names present in the block in sorted order.
func (self *server) LabelNames(ctx context.Context, rqst *monitoringpb.LabelNamesRequest) (*monitoringpb.LabelNamesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	strs, str, err := store.LabelNames(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.LabelNamesResponse{
		Labels:   strs,
		Warnings: str,
	}, nil
}

// LabelValues performs a query for the values of the given label.
func (self *server) LabelValues(ctx context.Context, rqst *monitoringpb.LabelValuesRequest) (*monitoringpb.LabelValuesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	resultStr, warnings, err := store.LabelValues(ctx, rqst.Label, rqst.StartTime, rqst.EndTime)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.LabelValuesResponse{
		LabelValues: resultStr,
		Warnings:    warnings,
	}, nil
}

// Query performs a query for the given time.
func (self *server) Query(ctx context.Context, rqst *monitoringpb.QueryRequest) (*monitoringpb.QueryResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	ts := time.Unix(int64(rqst.GetTs()), 0)
	resultStr, warnings, err := store.Query(ctx, rqst.Query, ts)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.QueryResponse{
		Value:    resultStr,
		Warnings: warnings,
	}, nil
}

// QueryRange performs a query for the given range.
func (self *server) QueryRange(rqst *monitoringpb.QueryRangeRequest, stream monitoringpb.MonitoringService_QueryRangeServer) error {

	store := self.stores[rqst.ConnectionId]
	ctx := stream.Context()

	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	startTime := time.Unix(int64(rqst.GetStartTime()), 0)
	endTime := time.Unix(int64(rqst.GetEndTime()), 0)
	step := rqst.Step

	resultStr, warnings, err := store.QueryRange(ctx, rqst.GetQuery(), startTime, endTime, step)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	maxSize := 2000
	for i := 0; i < len(resultStr); i += maxSize {
		rsp := new(monitoringpb.QueryRangeResponse)
		rsp.Warnings = warnings
		if i+maxSize < len(resultStr) {
			rsp.Value = resultStr[i : i+maxSize]
		} else {
			rsp.Value = resultStr[i:]
		}
		stream.Send(rsp)
	}

	return nil
}

// Series finds series by label matchers.
func (self *server) Series(ctx context.Context, rqst *monitoringpb.SeriesRequest) (*monitoringpb.SeriesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	startTime := time.Unix(int64(rqst.GetStartTime()), 0)
	endTime := time.Unix(int64(rqst.GetEndTime()), 0)

	resultStr, warnings, err := store.Series(ctx, rqst.GetMatches(), startTime, endTime)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.SeriesResponse{
		LabelSet: resultStr,
		Warnings: warnings,
	}, nil
}

// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
// under the TSDB's data directory and returns the directory as response.
func (self *server) Snapshot(ctx context.Context, rqst *monitoringpb.SnapshotRequest) (*monitoringpb.SnapshotResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	resultStr, err := store.Snapshot(ctx, rqst.GetSkipHead())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.SnapshotResponse{
		Result: resultStr,
	}, nil
}

// Rules returns a list of alerting and recording rules that are currently loaded.
func (self *server) Rules(ctx context.Context, rqst *monitoringpb.RulesRequest) (*monitoringpb.RulesResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	resultStr, err := store.Rules(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.RulesResponse{
		Result: resultStr,
	}, nil
}

// Targets returns an overview of the current state of the Prometheus target discovery.
func (self *server) Targets(ctx context.Context, rqst *monitoringpb.TargetsRequest) (*monitoringpb.TargetsResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	resultStr, err := store.Targets(ctx)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.TargetsResponse{
		Result: resultStr,
	}, nil
}

// TargetsMetadata returns metadata about metrics currently scraped by the target.
func (self *server) TargetsMetadata(ctx context.Context, rqst *monitoringpb.TargetsMetadataRequest) (*monitoringpb.TargetsMetadataResponse, error) {
	store := self.stores[rqst.ConnectionId]
	if store == nil {
		err := errors.New("No store connection exist for id " + rqst.ConnectionId)
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	resultStr, err := store.TargetsMetadata(ctx, rqst.GetMatchTarget(), rqst.GetMetric(), rqst.GetLimit())

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.TargetsMetadataResponse{
		Result: resultStr,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, "monitoring_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(monitoringpb.File_services_proto_monitoring_proto.Services().Get(0).FullName())
	s_impl.Proto = monitoringpb.File_services_proto_monitoring_proto.Path()
	s_impl.Port = defaultPort
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.PublisherId = "globulario"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.Permissions = make([]interface{}, 0)
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
		monitoringpb.RegisterMonitoringServiceServer(s_impl.grpcServer, s_impl)
		reflection.Register(s_impl.grpcServer)

		// Start the service.
		s_impl.StartService()
	}
}
