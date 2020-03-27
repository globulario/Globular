package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/monitoring/monitoring_store"
	"github.com/davecourtois/Globular/monitoring/monitoringpb"

	"time"

	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

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
	KeepAlive          bool

	// That map contain the list of active connections.
	Connections map[string]connection
	stores      map[string]monitoring_store.Store
}

// Create the configuration file if is not already exist.
func (self *server) init() {
	// Initialyse connection maps.
	self.Connections = make(map[string]connection, 0)

	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}

	// Initialyse stores.
	self.stores = make(map[string]monitoring_store.Store, 0)
	for _, c := range self.Connections {
		if c.Type == monitoringpb.StoreType_PROMETHEUS {
			address := "http://" + c.Host + ":" + Utility.ToString(c.Port)
			store, err := monitoring_store.NewPrometheusStore(address)
			if err == nil {
				self.stores[c.Id] = store
			} else {
				log.Println("fail to connect to "+address, err)
			}
		}
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
	err = self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Print the success message here.
	log.Println("Connection " + c.Id + " was created with success!")

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
	err := self.save()
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

	resultStr, warnings, err := store.LabelValues(ctx, rqst.Label)

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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.ServerUnaryInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(Interceptors.ServerUnaryInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamInterceptor))
	}

	monitoringpb.RegisterMonitoringServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

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
