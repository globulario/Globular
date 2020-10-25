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

	"github.com/davecourtois/Globular/Interceptors"

	globular "github.com/davecourtois/Globular/services/golang/globular_service"
	"github.com/davecourtois/Globular/services/golang/plc/plc_client"
	"github.com/davecourtois/Globular/services/golang/plc_link/plc_link_pb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10027
	defaultProxy = 10028

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// Thr IPV4 ServiceId
	ServiceId string = "127.0.0.1"

	domain string = "localhost"
)

// Contain necessary information to connect with a plc_server.
type connection struct {
	Id        string
	ServiceId string // 127.0.0.1:1234
	Domain    string // localhost
}

// Tag type.
type Tag struct {
	ServiceId    string
	ConnectionId string // Must be already define within the plc_server.
	Domain       string
	Name         string // The name of the tag
	TypeName     string // Can be BOOL, SINT, INT, DINT, REAL.
	Offset       int32  // The offset where to begin to read the tag.
	Length       int32  // The size of the array.
	Unsigned     bool   // If true the tag is read as unsigned value
}

type Link struct {
	Id        string    // The link id
	Frequency int32     // The refresh rate in milisecond.
	Source    Tag       // The source tag.
	Target    Tag       // The target tag.
	exit      chan bool // stop synchronization.
}

// Value need by Globular to start the services...
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
	Description        string
	Keywords           []string
	Repositories       []string
	Discoveries        []string
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

	//clients map[string]*plc_client.Plc_Client
	clients *sync.Map

	// The list of link to keep up to date.
	Links map[string]Link
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
	self.clients = new(sync.Map)

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewPlcLinkService_Client", plc_client.NewPlcService_Client)

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
	return globular.StopService(self, self.grpcServer)
}

func (self *server) Stop(context.Context, *plc_link_pb.StopRequest) (*plc_link_pb.StopResponse, error) {
	return &plc_link_pb.StopResponse{}, self.StopService()
}

///////////////////// API //////////////////////////////

// Test a connection exist for a given tag.
func (self *server) setTagConnection(tag *Tag) (*plc_client.Plc_Client, error) {
	client, ok := self.clients.Load(tag.ConnectionId)
	if !ok {
		// Open connection with the client.
		var err error
		client, err = plc_client.NewPlcService_Client(tag.Domain, tag.ServiceId)
		if err != nil {
			return nil, err
		}
		self.clients.Store(tag.ConnectionId, client)
	}
	return client.(*plc_client.Plc_Client), nil
}

func (self *server) stopSynchronize(lnk *Link) {
	lnk.exit <- true
}

// Return the type name as float number.
func getTypeNameNumber(typeName string) float64 {

	if typeName == "BOOL" {
		return 0.0
	} else if typeName == "SINT" {
		return 1.0
	} else if typeName == "INT" {
		return 2.0
	} else if typeName == "DINT" {
		return 3.0
	} else if typeName == "REAL" {
		return 4.0
	}
	return -1.0
}

func (self *server) startSynchronize(lnk Link) {
	if lnk.exit == nil {
		lnk.exit = make(chan bool)
	}

	// The main sychronization function.
	nbTry := 10 // number of try to connect with the server before given up...
	go func(lnk Link) {
		// So here i will made use of ticker to synchronice the tag.
		tickChan := time.NewTicker(time.Millisecond * time.Duration(lnk.Frequency))
		for {
			select {
			case <-tickChan.C:
				// So here I will get the value from the source and put it in the
				// taget.
				source, err := self.setTagConnection(&lnk.Source)

				if err == nil {
					values, err := source.ReadTag(lnk.Source.ConnectionId, lnk.Source.Name, getTypeNameNumber(lnk.Source.TypeName), lnk.Source.Offset, lnk.Source.Length, lnk.Source.Unsigned)
					if err != nil {
						time.Sleep(1000)
						if nbTry == 0 {
							lnk.exit <- true // stop synchronize on error.
							break
						}
						nbTry--
					}

					// so here I will write the value to the client.
					target, err := self.setTagConnection(&lnk.Target)
					if err == nil {
						err = target.WriteTag(lnk.Target.ConnectionId, lnk.Target.Name, getTypeNameNumber(lnk.Target.TypeName), values, lnk.Target.Offset, lnk.Target.Length, lnk.Target.Unsigned)
						if err != nil {
							time.Sleep(1000)
							if nbTry == 0 {
								lnk.exit <- true // stop synchronize on error.
								break
							}
							nbTry--
						}
					}
				}

			case <-lnk.exit:
				fmt.Println("Done whit tag synchronization for link ", lnk.Id)
				return
			}
		}
	}(lnk)
}

// Link Tow tag together and make it refresh at a given frequency in milisecond.
func (self *server) Link(ctx context.Context, rqst *plc_link_pb.LinkRqst) (*plc_link_pb.LinkRsp, error) {
	// first of all I will test if connection for link exist in the map of connections.
	var src, trg Tag

	// set the source
	src.ServiceId = rqst.Lnk.Source.ServiceId
	src.ConnectionId = rqst.Lnk.Source.ConnectionId
	src.Name = rqst.Lnk.Source.Name
	src.TypeName = rqst.Lnk.Source.TypeName
	src.Offset = rqst.Lnk.Source.Offset
	src.Length = rqst.Lnk.Source.Length

	// Set the target
	trg.ServiceId = rqst.Lnk.Target.ServiceId
	trg.ConnectionId = rqst.Lnk.Target.ConnectionId
	trg.Name = rqst.Lnk.Target.Name
	trg.TypeName = rqst.Lnk.Target.TypeName
	trg.Offset = rqst.Lnk.Target.Offset
	trg.Length = rqst.Lnk.Target.Length

	// So here the connection are found so I will create the link.
	lnk := Link{Id: rqst.Lnk.Id, Frequency: rqst.Lnk.Frequency, Source: src, Target: trg, exit: make(chan bool)}

	// set the lnk...
	self.Links[rqst.Lnk.Id] = lnk

	// In that case I will save it in file.
	err := self.Save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.startSynchronize(lnk)

	return &plc_link_pb.LinkRsp{
		Result: true,
	}, nil
}

// Remove link from tow task.
func (self *server) UnLink(ctx context.Context, rqst *plc_link_pb.UnLinkRqst) (*plc_link_pb.UnLinkRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// Stop the synchronization
	lnk := self.Links[rqst.Id]
	self.stopSynchronize(&lnk)

	// Remove the link from the map.
	delete(self.Links, rqst.Id)

	err := self.Save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &plc_link_pb.UnLinkRsp{
		Result: true,
	}, nil
}

// Suspend the synchronization of tow tags.
func (self *server) Suspend(ctx context.Context, rqst *plc_link_pb.SuspendRqst) (*plc_link_pb.SuspendRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// Stop the synchronization
	lnk := self.Links[rqst.Id]
	self.stopSynchronize(&lnk)

	return &plc_link_pb.SuspendRsp{
		Result: true,
	}, nil
}

// Resume the synchronization of tow tags.
func (self *server) Resume(ctx context.Context, rqst *plc_link_pb.ResumeRqst) (*plc_link_pb.ResumeRsp, error) {
	if _, ok := self.Links[rqst.Id]; !ok {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No link found with id "+rqst.Id)))
	}

	// reset the chan bool...
	lnk := self.Links[rqst.Id]
	lnk.exit = make(chan bool)
	self.Links[rqst.Id] = lnk

	// Start the synchronization.
	self.startSynchronize(lnk)
	return &plc_link_pb.ResumeRsp{
		Result: true,
	}, nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "plc_link_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(plc_link_pb.File_services_proto_plc_link_proto.Services().Get(0).FullName())
	s_impl.Proto = plc_link_pb.File_services_proto_plc_link_proto.Path()
	s_impl.Port = defaultPort
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = domain
	s_impl.Links = make(map[string]Link, 0)
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
		plc_link_pb.RegisterPlcLinkServiceServer(s_impl.grpcServer, s_impl)
		reflection.Register(s_impl.grpcServer)

		// Start the service.
		s_impl.StartService()
	}
}
