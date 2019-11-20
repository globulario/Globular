package smtp_client

import (
	// "context"
	// "log"

	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/smtp/smtppb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// SMTP Client Service
////////////////////////////////////////////////////////////////////////////////
type SMTP_Client struct {
	cc *grpc.ClientConn
	c  smtppb.SmtpServiceClient

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
func NewSmtp_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string) *SMTP_Client {
	client := new(SMTP_Client)

	client.domain = domain
	client.name = "smtp"
	client.port = port
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile

	client.cc = api.GetClientConnection(client)
	client.c = smtppb.NewSmtpServiceClient(client.cc)

	return client
}

// Return the domain
func (self *SMTP_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *SMTP_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *SMTP_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *SMTP_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *SMTP_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *SMTP_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *SMTP_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *SMTP_Client) GetCaFile() string {
	return self.caFile
}
