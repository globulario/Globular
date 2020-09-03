package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"errors"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/smtp/smtp_client"
	"github.com/davecourtois/Globular/smtp/smtppb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

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

	// The map of connection...
	Connections map[string]connection
}

func (self *server) init() {

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewSmtp_Client", smtp_client.NewSmtp_Client)

	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		if len(self.Id) == 0 {
			// Generate random id for the server instance.
			self.Id = Utility.RandomUUID()
		}
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

// Create a new SQL connection and store it for futur use. If the connection already
// exist it will be replace by the new one.
func (self *server) CreateConnection(ctx context.Context, rsqt *smtppb.CreateConnectionRqst) (*smtppb.CreateConnectionRsp, error) {
	fmt.Println("Try to create a new connection")
	// sqlpb
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
	err = self.save()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

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
	err := self.save()
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
	s_impl.Path, _ = os.Executable()
	package_ := string(smtppb.File_smtp_smtppb_smtp_proto.Package().Name())
	s_impl.Path = s_impl.Path[strings.Index(s_impl.Path, package_):]
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
		grpcServer = grpc.NewServer([]grpc.ServerOption{grpc.UnaryInterceptor(Interceptors.ServerUnaryInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamInterceptor)}...)
	}

	// Register the smtp service.
	smtppb.RegisterSmtpServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {

		// no web-rpc server.
		fmt.Println(s_impl.Name + " grpc service is starting")
		if err := grpcServer.Serve(lis); err != nil {

			if err.Error() == "signal: killed" {
				fmt.Println("service ", s_impl.Name, " was stop!")
			}
		}

	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}
