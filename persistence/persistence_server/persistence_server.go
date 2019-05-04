package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"github.com/davecourtois/Globular/persistence/store"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Name     string
	Path     string
	Port     int
	Protocol string

	// The data store.
	s *store.Store
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

	// initialyse the data store.
	self.s = store.NewStore()
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

// Implementation of the Persistence method.
func (self *server) PersistEntity(ctx context.Context, rqst *persistencepb.PersistEntityRqst) (*persistencepb.PersistEntityRsp, error) {
	fmt.Println("Persist a value")

	// In that case I will save it in file.
	err := self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	entity := rqst.Entity

	// Here I will compose the entity url.
	address, _ := Utility.MyIP() // Here I will set the url to retreive the entity over the network...
	entity.Url = address.Hostname + ":" + strconv.Itoa(self.Port) + "/?Typename=" + entity.Typename + "&UUID=" + entity.UUID

	err = self.s.PersistEntity(entity)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &persistencepb.PersistEntityRsp{
		Result: entity.Url,
	}, nil
}

// Retreive entity by it uuid.
func (self *server) GetEntityByUuid(ctx context.Context, rqst *persistencepb.GetEntityByUuidRqst) (*persistencepb.GetEntityByUuidRsp, error) {

	entity, err := self.s.GetEntityByUuid(rqst.Typename, rqst.Uuid)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return the entity retreive from the store.
	return &persistencepb.GetEntityByUuidRsp{
		Entity: entity,
	}, nil
}

// Return all entities for a given typeName.
func (self *server) GetEntitiesByTypename(rqst *persistencepb.GetEntitiesByTypenameRqst, stream persistencepb.PersistenceService_GetEntitiesByTypenameServer) error {

	entities, err := self.s.GetEntitiesByTypename(rqst.Typename)

	if err != nil {
		return status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// return one entity at time to the stream...
	for i := 0; i < len(entities); i++ {
		stream.Send(&persistencepb.GetEntitiesByTypenameRsp{
			Entity: entities[i],
		})
	}

	return nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {
	log.Println("Persistence grpc service is starting")

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "persistence_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := 50051 // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// First of all I will creat a listener.
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = "persistence_server"
	s_impl.Port = port
	s_impl.Protocol = "grpc"
	s_impl.Path = os.Args[0] // keep the execution path here...

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	persistencepb.RegisterPersistenceServiceServer(grpcServer, s_impl)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	log.Println("Persistence grpc service is closed")
}
