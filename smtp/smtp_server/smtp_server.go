package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"errors"

	"github.com/davecourtois/Globular/api"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/api/client"
	"github.com/davecourtois/Globular/smtp/smtppb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	//"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	gomail "gopkg.in/gomail.v1"
)

var (
	defaultPort  = 10037
	defaultProxy = 10038

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// the domain of the server.
	domain string = "localhost"
)

// Keep connection information here.
type connection struct {
	Id       string // The connection id
	Host     string // can also be ipv4 addresse.
	User     string
	Password string
	Port     int32
}

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
	CertFile           string
	KeyFile            string
	CertAuthorityTrust string
	Version            string
	TLS                bool
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	// The grpc server.
	grpcServer *grpc.Server

	// The map of connection...
	Connections map[string]connection
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
	Utility.RegisterFunction("NewSmtp_Client", client.NewSmtp_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := api.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = api.InitGrpcServer(self, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err != nil {
		return err
	}

	return nil

}

// Save the configuration values.
func (self *server) Save() error {
	// Create the file...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return api.SaveService(dir+"/config.json", self)
}

func (self *server) Start() error {
	return api.StartService(self, self.grpcServer)
}

func (self *server) Stop() error {
	return api.StopService(self)
}

//////////////////////////// SMPT specific functions ///////////////////////////

// Create a new connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *smtppb.CreateConnectionRqst) (*smtppb.CreateConnectionRsp, error) {

	fmt.Println("Try to create a new connection")
	var c connection
	var err error

	// Set the connection info from the request.
	c.Id = rsqt.Connection.Id
	c.Host = rsqt.Connection.Host
	c.Port = rsqt.Connection.Port
	c.User = rsqt.Connection.User
	c.Password = rsqt.Connection.Password

	// set or update the connection and save it in json file.
	self.Connections[c.Id] = c

	// In that case I will save it in file.
	err = self.Save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	api.UpdateServiceConfig(self)

	// test if the connection is reacheable.
	// _, err = self.ping(ctx, c.Id)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &smtppb.CreateConnectionRsp{
		Result: true,
	}, nil
}

// Remove a connection from the map and the file.
func (self *server) DeleteConnection(ctx context.Context, rqst *smtppb.DeleteConnectionRqst) (*smtppb.DeleteConnectionRsp, error) {
	id := rqst.GetId()
	if _, ok := self.Connections[id]; !ok {
		return &smtppb.DeleteConnectionRsp{
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
	return &smtppb.DeleteConnectionRsp{
		Result: true,
	}, nil
}

/**
 * Carbon copy list...
 */
type CarbonCopy struct {
	EMail string
	Name  string
}

/**
 * Attachment file, if the data is empty or nil
 * that means the file is on the server a the given path.
 */
type Attachment struct {
	FileName string
	FileData []byte
}

/**
 * Send mail... The server id is the authentification id...
 */
func (self *server) sendEmail(id string, from string, to []string, cc []*CarbonCopy, subject string, body string, attachs []*Attachment, bodyType string) error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)

	msg.SetHeader("To", to...)

	// Attach the multiple carbon copy...
	var cc_ []string
	for i := 0; i < len(cc); i++ {
		cc_ = append(cc_, msg.FormatAddress(cc[i].EMail, cc[i].Name))
	}

	if len(cc_) > 0 {
		msg.SetHeader("Cc", cc_...)
	}

	msg.SetHeader("Subject", subject)
	msg.SetBody(bodyType, body)

	for i := 0; i < len(attachs); i++ {
		f := gomail.CreateFile(attachs[i].FileName, attachs[i].FileData)
		msg.Attach(f)
	}

	config := self.Connections[id]

	mailer := gomail.NewMailer(config.Host, config.User, config.Password, int(config.Port))

	if err := mailer.Send(msg); err != nil {
		return err
	}
	return nil
}

// Send a simple email whitout file.
func (self *server) SendEmail(ctx context.Context, rqst *smtppb.SendEmailRqst) (*smtppb.SendEmailRsp, error) {
	if rqst.Email == nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No email message was given!")))
	}

	cc := make([]*CarbonCopy, len(rqst.Email.Cc))
	for i := 0; i < len(rqst.Email.Cc); i++ {
		cc[i] = &CarbonCopy{Name: rqst.Email.Cc[i].Name, EMail: rqst.Email.Cc[i].Address}
	}

	bodyType := "text/html"
	if rqst.Email.BodyType != smtppb.BodyType_HTML {
		bodyType = "text/html"
	}

	err := self.sendEmail(rqst.Id, rqst.Email.From, rqst.Email.To, cc, rqst.Email.Subject, rqst.Email.Body, []*Attachment{}, bodyType)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &smtppb.SendEmailRsp{
		Result: true,
	}, nil
}

// Send email with file attachement attachements.
func (self *server) SendEmailWithAttachements(stream smtppb.SmtpService_SendEmailWithAttachementsServer) error {

	// that buffer will contain the file attachement data while data is transfert.
	attachements := make([]*Attachment, 0)
	var bodyType string
	var body string
	var subject string
	var from string
	var to []string
	var cc []*CarbonCopy
	var id string

	// So here I will read the stream until it end...
	for {
		rqst, err := stream.Recv()
		if err == io.EOF {
			// Here all data is read...
			err := self.sendEmail(id, from, to, cc, subject, body, attachements, bodyType)

			if err != nil {
				return status.Errorf(
					codes.Internal,
					Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
			}

			// Close the stream...
			stream.SendAndClose(&smtppb.SendEmailWithAttachementsRsp{
				Result: true,
			})

			return nil
		}

		if err != nil {
			return status.Errorf(
				codes.Internal,
				Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
		}

		id = rqst.Id

		// Receive message informations.
		switch msg := rqst.Data.(type) {
		case *smtppb.SendEmailWithAttachementsRqst_Email:
			cc = make([]*CarbonCopy, len(msg.Email.Cc))

			// The email itself.
			for i := 0; i < len(msg.Email.Cc); i++ {
				cc[i] = &CarbonCopy{Name: msg.Email.Cc[i].Name, EMail: msg.Email.Cc[i].Address}
			}
			bodyType = "text"
			if msg.Email.BodyType == smtppb.BodyType_HTML {
				bodyType = "html"
			}
			from = msg.Email.From
			to = msg.Email.To
			body = msg.Email.Body
			subject = msg.Email.Subject

		case *smtppb.SendEmailWithAttachementsRqst_Attachements:
			var lastAttachement *Attachment
			if len(attachements) > 0 {
				lastAttachement = attachements[len(attachements)-1]
				if lastAttachement.FileName != msg.Attachements.FileName {
					lastAttachement = new(Attachment)
					lastAttachement.FileData = make([]byte, 0)
					lastAttachement.FileName = msg.Attachements.FileName
					attachements = append(attachements, lastAttachement)
				}
			} else {
				lastAttachement = new(Attachment)
				lastAttachement.FileData = make([]byte, 0)
				lastAttachement.FileName = msg.Attachements.FileName
				attachements = append(attachements, lastAttachement)
			}

			// Append the data in the file attachement.
			lastAttachement.FileData = append(lastAttachement.FileData, msg.Attachements.FileData...)
		}

	}

	return nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	//grpclog.SetLogger(log.New(os.Stdout, "smtp_service: ", log.LstdFlags))

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
	s_impl.Name = string(smtppb.File_smtp_smtppb_smtp_proto.Services().Get(0).FullName())
	s_impl.Proto = smtppb.File_smtp_smtppb_smtp_proto.Path()
	s_impl.Port = port
	s_impl.Domain = domain
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
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
	smtppb.RegisterSmtpServiceServer(s_impl.grpcServer, s_impl)
	reflection.Register(s_impl.grpcServer)

	// Start the service.
	s_impl.Start()
}
