package main

import (
	"context"
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

	//	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	"github.com/davecourtois/Globular/storage/storage_store"
	"github.com/davecourtois/Globular/storage/storagepb"
	"github.com/davecourtois/Utility"
)

var (
	defaultPort  = 10013
	defaultProxy = 10014

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""
)

// Keep connection information here.
type connection struct {
	Id    string // The connection id
	Name  string // The kv store name
	Type  storagepb.StoreType
	store storage_store.Store
}

type server struct {

	// The global attribute of the services.
	Name            string
	Port            int
	Proxy           int
	Protocol        string
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.

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

// Create a new KV connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *storagepb.CreateConnectionRqst) (*storagepb.CreateConnectionRsp, error) {
	fmt.Println("Try to create a new connection")

	fmt.Println("Try to create a new connection")
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
	err := self.save()
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

	conn.store = store

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

	self.Connections[rqst.GetId()] = conn

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

	store := self.Connections[rqst.GetId()].store

	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

	err := store.Close()
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

	store := self.Connections[rqst.GetId()].store
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

	store := self.Connections[rqst.GetId()].store
	if store == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("no store found for connection with id "+rqst.GetId())))
	}

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

	store := self.Connections[rqst.GetId()].store
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

	store := self.Connections[rqst.GetId()].store
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

	store := self.Connections[rqst.GetId()].store
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
	grpclog.SetLogger(log.New(os.Stdout, "storage_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// First of all I will creat a listener.
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Connections = make(map[string]connection)
	s_impl.Name = Utility.GetExecName(os.Args[0])
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"

	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(grpcServer, s_impl)

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
