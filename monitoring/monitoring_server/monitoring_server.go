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

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/monitoring/monitoring_store"
	"github.com/davecourtois/Globular/monitoring/monitoringpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
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

// This is the connction to a datastore.
type connection struct {
	Id    string
	Host  string
	Store monitoringpb.StoreType
	Port  int32
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

	// saved connections
	Connections map[string]*connection

	// The map of store (also connections...)
	stores map[string]monitoring_store.Store
}

// Create the configuration file if is not already exist.
func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}

	// open the connection with store.
	for _, c := range self.Connections {
		if c.Store == monitoringpb.StoreType_PROMETHEUS {
			s, _ := monitoring_store.NewPrometheusStore(c.Host + ":" + Utility.ToString(c.Port))
			self.stores[c.Id] = s
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
	fmt.Println("Try to create a new connection")

	var err error
	c := new(connection)
	c.Host = rqst.Connection.Host
	c.Id = rqst.Connection.Id
	c.Port = rqst.Connection.Port
	c.Store = rqst.Connection.Store

	// set the connection into the map.
	self.Connections[c.Id] = c

	if self.stores[c.Id] == nil {
		if c.Store == monitoringpb.StoreType_PROMETHEUS {
			s, err := monitoring_store.NewPrometheusStore(c.Host + ":" + Utility.ToString(c.Port))
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

			// keep the store in the map.
			self.stores[c.Id] = s

		} else {
			return nil, status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("Invlaid monitoring store type!")))
		}
	}
	err = self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.CreateConnectionRsp{
		Result: true,
	}, nil

}

// Delete a connection.
func (self *server) DeleteConnection(ctx context.Context, rqst *monitoringpb.DeleteConnectionRqst) (*monitoringpb.DeleteConnectionRsp, error) {
	if self.Connections[rqst.Id] == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.Id)))
	}

	delete(self.Connections, rqst.Id)
	if self.stores[rqst.Id] != nil {
		delete(self.stores, rqst.Id)
	}

	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.DeleteConnectionRsp{
		Result: true,
	}, nil
}

// Alerts returns a list of all active alerts.
func (self *server) Alerts(ctx context.Context, rqst *monitoringpb.AlertsRequest) (*monitoringpb.AlertsResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	results, err := s.Alerts(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.AlertsResponse{
		Results: results,
	}, nil
}

// AlertManagers returns an overview of the current state of the Prometheus alert manager discovery.
func (self *server) AlertManagers(ctx context.Context, rqst *monitoringpb.AlertManagersRequest) (*monitoringpb.AlertManagersResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	results, err := s.AlertManagers(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.AlertManagersResponse{
		Results: results,
	}, nil
}

// CleanTombstones removes the deleted data from disk and cleans up the existing tombstones.
func (self *server) CleanTombstones(ctx context.Context, rqst *monitoringpb.CleanTombstonesRequest) (*monitoringpb.CleanTombstonesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	err := s.CleanTombstones(ctx)
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
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	results, err := s.Config(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.ConfigResponse{
		Results: results,
	}, nil
}

// DeleteSeries deletes data for a selection of series in a time range.
func (self *server) DeleteSeries(ctx context.Context, rqst *monitoringpb.DeleteSeriesRequest) (*monitoringpb.DeleteSeriesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}
	startTime := time.Unix(int64(Utility.ToInt(rqst.StartTime)), 0)
	endTime := time.Unix(int64(Utility.ToInt(rqst.EndTime)), 0)
	err := s.DeleteSeries(ctx, rqst.Matches, startTime, endTime)
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
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	results, err := s.Flags(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.FlagsResponse{
		Results: results,
	}, nil
}

// LabelNames returns all the unique label names present in the block in sorted order.
func (self *server) LabelNames(ctx context.Context, rqst *monitoringpb.LabelNamesRequest) (*monitoringpb.LabelNamesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	labels, warnings, err := s.LabelNames(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.LabelNamesResponse{
		Labels:   labels,
		Warnings: warnings,
	}, nil
}

// LabelValues performs a query for the values of the given label.
func (self *server) LabelValues(ctx context.Context, rqst *monitoringpb.LabelValuesRequest) (*monitoringpb.LabelValuesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	labelValues, warnings, err := s.LabelValues(ctx, rqst.Label)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.LabelValuesResponse{
		LabelValues: labelValues,
		Warnings:    warnings,
	}, nil
}

// Query performs a query for the given time.
func (self *server) Query(ctx context.Context, rqst *monitoringpb.QueryRequest) (*monitoringpb.QueryResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	ts := time.Unix(int64(rqst.Ts), 0)

	results, warnings, err := s.Query(ctx, rqst.Query, ts)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.QueryResponse{
		Value:    results,
		Warnings: warnings,
	}, nil
}

// QueryRange performs a query for the given range.
func (self *server) QueryRange(rqst *monitoringpb.QueryRangeRequest, stream monitoringpb.MonitoringService_QueryRangeServer) error {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	// So here I will return the results.
	startTime := time.Unix(int64(Utility.ToInt(rqst.StartTime)), 0)
	endTime := time.Unix(int64(Utility.ToInt(rqst.EndTime)), 0)
	results, warnings, err := s.QueryRange(stream.Context(), rqst.Query, startTime, endTime, rqst.Step)

	// I will now stream back the results.
	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	size := 1024 // send 1024 bytes at time
	for i := 0; i < len(results); i += size {
		// Send the last bytes as necessary.
		if len(results)-size <= i {
			stream.Send(
				&monitoringpb.QueryRangeResponse{
					Value:    results[i:], // character at time.
					Warnings: warnings,
				},
			)
			return nil
		}

		stream.Send(
			&monitoringpb.QueryRangeResponse{
				Value:    results[i : i+size], // character at time.
				Warnings: warnings,
			},
		)

	}

	return nil

}

// Series finds series by label matchers.
func (self *server) Series(ctx context.Context, rqst *monitoringpb.SeriesRequest) (*monitoringpb.SeriesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	startTime := time.Unix(int64(Utility.ToInt(rqst.StartTime)), 0)
	endTime := time.Unix(int64(Utility.ToInt(rqst.EndTime)), 0)
	results, warnings, err := s.Series(ctx, rqst.Matches, startTime, endTime)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.SeriesResponse{
		LabelSet: results,
		Warnings: warnings,
	}, nil
}

// Snapshot creates a snapshot of all current data into snapshots/<datetime>-<rand>
// under the TSDB's data directory and returns the directory as response.
func (self *server) Snapshot(ctx context.Context, rqst *monitoringpb.SnapshotRequest) (*monitoringpb.SnapshotResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	result, err := s.Snapshot(ctx, rqst.SkipHead)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.SnapshotResponse{
		Result: result,
	}, nil
}

// Rules returns a list of alerting and recording rules that are currently loaded.
func (self *server) Rules(ctx context.Context, rqst *monitoringpb.RulesRequest) (*monitoringpb.RulesResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	result, err := s.Rules(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.RulesResponse{
		Result: result,
	}, nil
}

// Targets returns an overview of the current state of the Prometheus target discovery.
func (self *server) Targets(ctx context.Context, rqst *monitoringpb.TargetsRequest) (*monitoringpb.TargetsResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	result, err := s.Targets(ctx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.TargetsResponse{
		Result: result,
	}, nil
}

// TargetsMetadata returns metadata about metrics currently scraped by the target.
func (self *server) TargetsMetadata(ctx context.Context, rqst *monitoringpb.TargetsMetadataRequest) (*monitoringpb.TargetsMetadataResponse, error) {
	s := self.stores[rqst.ConnectionId]
	if s == nil {
		return nil,
			status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No connection found with id "+rqst.ConnectionId)))
	}

	result, err := s.TargetsMetadata(ctx, rqst.MatchTarget, rqst.Metric, rqst.Limit)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &monitoringpb.TargetsMetadataResponse{
		Result: result,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "monitoring_service: ", log.LstdFlags))

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
	s_impl.PublisherId = "localhost"
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

	monitoringpb.RegisterMonitoringServiceServer(grpcServer, s_impl)
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
