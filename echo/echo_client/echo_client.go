package echo_client

import (
	"log"
	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/echo/echopb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Echo_Client struct {
	cc *grpc.ClientConn
	c  echopb.EchoServiceClient

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
func NewEcho_Client(address string, name string) (*Echo_Client, error) {
	client := new(Echo_Client)
	err := api.InitClient(client, address, name)
	if err != nil {
		return nil, err
	}
	client.cc = api.GetClientConnection(client)
	client.c = echopb.NewEchoServiceClient(client.cc)

	return client, nil
}

// Return the domain
func (self *Echo_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Echo_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *Echo_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Echo_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Echo_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Echo_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Echo_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Echo_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Echo_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Echo_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Echo_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Echo_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Echo_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Echo_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Echo_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

////////////////// Api //////////////////////
func (self *Echo_Client) Echo(msg interface{}) (string, error) {
	log.Println("echo service call: ", msg)
	rqst := &echopb.EchoRequest{
		Message: Utility.ToString(msg),
	}
	ctx := api.GetClientContext(self)
	rsp, err := self.c.Echo(ctx, rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}
