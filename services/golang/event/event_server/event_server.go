package main

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"fmt"

	//	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/event/event_client"
	"github.com/davecourtois/Globular/services/golang/event/eventpb"
	globular "github.com/davecourtois/Globular/services/golang/globular_service"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10050
	defaultProxy = 10051

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// the default domain.
	domain string = "localhost"
)

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

	// Use to sync event channel manipulation.
	actions chan map[string]interface{}

	// stop the processing loop.
	exit chan bool
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
	Utility.RegisterFunction("NewEventService_Client", event_client.NewEventService_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := globular.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = globular.InitGrpcServer(self /*Interceptors.ServerUnaryInterceptor*/, nil /*Interceptors.ServerStreamInterceptor*/, nil)
	if err != nil {
		return err
	}

	self.exit = make(chan bool)

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
	self.grpcServer.Stop()
	return nil //globular.StopService(self, self.grpcServer)
}

func (self *server) Stop(context.Context, *eventpb.StopRequest) (*eventpb.StopResponse, error) {
	self.exit <- true
	fmt.Println(`Stop event server was called.`)
	return &eventpb.StopResponse{}, self.StopService()
}

//////////////////////////////////////////////////////////////////////////////
//	Services implementation.
//////////////////////////////////////////////////////////////////////////////

// That function process channel operation and run in it own go routine.
func (self *server) run() {
	fmt.Println("start event service")
	channels := make(map[string][]string)
	streams := make(map[string]eventpb.EventService_OnEventServer)
	quits := make(map[string]chan bool)

	// Here will create the action channel.
	self.actions = make(chan map[string]interface{})

	for {
		select {
		case <-self.exit:
			fmt.Println("--------> exit from the run loop.")
			break
		case a := <-self.actions:
			action := a["action"].(string)
			if action == "onevent" {
				streams[a["uuid"].(string)] = a["stream"].(eventpb.EventService_OnEventServer)
				quits[a["uuid"].(string)] = a["quit"].(chan bool)
			} else if action == "subscribe" {
				if channels[a["name"].(string)] == nil {
					channels[a["name"].(string)] = make([]string, 0)
				}
				if !Utility.Contains(channels[a["name"].(string)], a["uuid"].(string)) {
					channels[a["name"].(string)] = append(channels[a["name"].(string)], a["uuid"].(string))
				}
			} else if action == "publish" {
				if channels[a["name"].(string)] != nil {
					toDelete := make([]string, 0)
					for i := 0; i < len(channels[a["name"].(string)]); i++ {
						uuid := channels[a["name"].(string)][i]
						stream := streams[uuid]
						if stream != nil {
							// Here I will send data to stream.
							err := stream.Send(&eventpb.OnEventResponse{
								Evt: &eventpb.Event{
									Name: a["name"].(string),
									Data: a["data"].([]byte),
								},
							})

							// In case of error I will remove the subscriber
							// from the list.
							if err != nil {
								// append to channle list to be close.
								toDelete = append(toDelete, uuid)
							}
						}
					}

					// remove closed channel
					for i := 0; i < len(toDelete); i++ {
						uuid := toDelete[i]
						// remove uuid from all channels.
						for name, channel := range channels {
							uuids := make([]string, 0)
							for i := 0; i < len(channel); i++ {
								if uuid != channel[i] {
									uuids = append(uuids, channel[i])
								}
							}
							channels[name] = uuids
						}
						// return from OnEvent
						quits[uuid] <- true
						// remove the channel from the map.
						delete(quits, uuid)
					}
				}
			} else if action == "unsubscribe" {
				uuids := make([]string, 0)
				for i := 0; i < len(channels[a["name"].(string)]); i++ {
					if a["uuid"].(string) != channels[a["name"].(string)][i] {
						uuids = append(uuids, channels[a["name"].(string)][i])
					}
				}
				channels[a["name"].(string)] = uuids
			} else if action == "quit" {
				// remove uuid from all channels.
				for name, channel := range channels {
					uuids := make([]string, 0)
					for i := 0; i < len(channel); i++ {
						if a["uuid"].(string) != channel[i] {
							uuids = append(uuids, channel[i])
						}
					}
					channels[name] = uuids
				}
				// return from OnEvent
				quits[a["uuid"].(string)] <- true
				// remove the channel from the map.
				delete(quits, a["uuid"].(string))
			}
		}
	}
}

// Connect to an event channel or create it if it not already exist
// and stay in that function until UnSubscribe is call.
func (self *server) Quit(ctx context.Context, rqst *eventpb.QuitRequest) (*eventpb.QuitResponse, error) {
	quit := make(map[string]interface{})
	quit["action"] = "quit"
	quit["uuid"] = rqst.Uuid

	self.actions <- quit

	return &eventpb.QuitResponse{
		Result: true,
	}, nil
}

// Connect to an event channel or create it if it not already exist
// and stay in that function until UnSubscribe is call.
func (self *server) OnEvent(rqst *eventpb.OnEventRequest, stream eventpb.EventService_OnEventServer) error {
	onevent := make(map[string]interface{})
	onevent["action"] = "onevent"
	onevent["stream"] = stream
	onevent["uuid"] = rqst.Uuid
	onevent["quit"] = make(chan bool)

	self.actions <- onevent

	// wait util unsbscribe or connection is close.
	<-onevent["quit"].(chan bool)

	return nil
}

func (self *server) Subscribe(ctx context.Context, rqst *eventpb.SubscribeRequest) (*eventpb.SubscribeResponse, error) {
	subscribe := make(map[string]interface{})
	subscribe["action"] = "subscribe"
	subscribe["name"] = rqst.Name
	subscribe["uuid"] = rqst.Uuid
	self.actions <- subscribe

	return &eventpb.SubscribeResponse{
		Result: true,
	}, nil
}

// Disconnect to an event channel.(Return from Subscribe)
func (self *server) UnSubscribe(ctx context.Context, rqst *eventpb.UnSubscribeRequest) (*eventpb.UnSubscribeResponse, error) {
	unsubscribe := make(map[string]interface{})
	unsubscribe["action"] = "unsubscribe"
	unsubscribe["name"] = rqst.Name
	unsubscribe["uuid"] = rqst.Uuid

	self.actions <- unsubscribe

	return &eventpb.UnSubscribeResponse{
		Result: true,
	}, nil
}

// Publish event on channel.
func (self *server) Publish(ctx context.Context, rqst *eventpb.PublishRequest) (*eventpb.PublishResponse, error) {
	publish := make(map[string]interface{})
	publish["action"] = "publish"
	publish["name"] = rqst.Evt.Name
	publish["data"] = rqst.Evt.Data

	// publish the data.
	self.actions <- publish
	return &eventpb.PublishResponse{
		Result: true,
	}, nil

	return nil, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(eventpb.File_services_proto_event_proto.Services().Get(0).FullName())
	s_impl.Proto = eventpb.File_services_proto_event_proto.Path()
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = domain
	s_impl.Permissions = make([]interface{}, 0)

	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		fmt.Println("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
		return
	}

	// Register the echo services
	eventpb.RegisterEventServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go s_impl.run()

	// Start the service.
	err = s_impl.StartService()

	if err != nil {
		fmt.Println("Fail to start service %s: %s", s_impl.Name, s_impl.Id, err)
		return
	}

}
