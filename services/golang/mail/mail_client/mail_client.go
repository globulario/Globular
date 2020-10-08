package mail_client

import (
	"context"
	"log"

	"fmt"
	"io"
	"os"
	"strconv"

	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Globular/services/golang/mail/mailpb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// mail Client Service
////////////////////////////////////////////////////////////////////////////////
type Mail_Client struct {
	cc *grpc.ClientConn
	c  mailpb.MailServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewMailService_Client(address string, id string) (*Mail_Client, error) {
	client := new(Mail_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = mailpb.NewMailServiceClient(client.cc)

	return client, nil
}

func (self *Mail_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Mail_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Mail_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Mail_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *Mail_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Mail_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Mail_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Mail_Client) SetName(name string) {
	self.name = name
}

// Set the service instance id
func (self *Mail_Client) SetId(id string) {
	self.id = id
}

// Set the domain.
func (self *Mail_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Mail_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Mail_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Mail_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Mail_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Mail_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Mail_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Mail_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Mail_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

//////////////////////////////// Api ////////////////////////////////

// Stop the service.
func (self *Mail_Client) StopService() {
	self.c.Stop(globular.GetClientContext(self), &mailpb.StopRequest{})
}

/**
 * Create a connection with a mail server.
 */
func (self *Mail_Client) CreateConnection(id string, user string, pwd string, port int, host string) error {
	rqst := &mailpb.CreateConnectionRqst{
		Connection: &mailpb.Connection{
			Id:       id,
			User:     user,
			Password: pwd,
			Port:     int32(port),
			Host:     host,
		},
	}

	_, err := self.c.CreateConnection(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Delete a connection with a mail server.
 */
func (self *Mail_Client) DeleteConnection(id string) error {

	rqst := &mailpb.DeleteConnectionRqst{
		Id: id,
	}
	_, err := self.c.DeleteConnection(globular.GetClientContext(self), rqst)
	return err
}

/**
 * Send email whiout files.
 */
func (self *Mail_Client) SendEmail(id string, from string, to []string, cc []*mailpb.CarbonCopy, subject string, body string, bodyType int32) error {

	rqst := &mailpb.SendEmailRqst{
		Id: id,
		Email: &mailpb.Email{
			From:     from,
			To:       to,
			Cc:       cc,
			Subject:  subject,
			Body:     body,
			BodyType: mailpb.BodyType(bodyType),
		},
	}

	_, err := self.c.SendEmail(globular.GetClientContext(self), rqst)

	return err
}

/**
 * Send email with files.
 */
/**
 * Here I will make use of a stream
 */
func sendFile(id string, path string, stream mailpb.MailService_SendEmailWithAttachementsClient) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Fail to open file "+path+" with error: %v", err)
	}

	// close the file when done.
	defer file.Close()

	const BufferSize = 1024 * 5 // the chunck size.

	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if bytesread > 0 {
			rqst := &mailpb.SendEmailWithAttachementsRqst{
				Id: id,
				Data: &mailpb.SendEmailWithAttachementsRqst_Attachements{
					Attachements: &mailpb.Attachement{
						FileName: path,
						FileData: buffer[:bytesread],
					},
				},
			}
			err = stream.Send(rqst)
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
}

/**
 * Test send email whit attachements.
 */
func (self *Mail_Client) SendEmailWithAttachements(id string, from string, to []string, cc []*mailpb.CarbonCopy, subject string, body string, bodyType int32, files []string) error {

	// Open the stream...
	stream, err := self.c.SendEmailWithAttachements(globular.GetClientContext(self))
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	// Send file attachment as a stream, not need to be send first.
	for i := 0; i < len(files); i++ {
		sendFile(id, files[i], stream)
	}

	// Send the email message...
	rqst := &mailpb.SendEmailWithAttachementsRqst{
		Id: id,
		Data: &mailpb.SendEmailWithAttachementsRqst_Email{
			Email: &mailpb.Email{
				From:     from,
				To:       to,
				Cc:       cc,
				Subject:  subject,
				Body:     body,
				BodyType: mailpb.BodyType(bodyType),
			},
		},
	}

	err = stream.Send(rqst)
	if err != nil {
		return err
	}

	_, err = stream.CloseAndRecv()

	return err

}
