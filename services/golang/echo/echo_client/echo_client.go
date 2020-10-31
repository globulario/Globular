package echo_client

import (
	"strconv"

	"context"

	"github.com/globulario/Globular/services/golang/echo/echopb"
	globular "github.com/globulario/Globular/services/golang/globular_client"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Echo_Client struct {
	cc *grpc.ClientConn
	c  echopb.EchoServiceClient

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
func NewEchoService_Client(address string, id string) (*Echo_Client, error) {
	client := new(Echo_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = echopb.NewEchoServiceClient(client.cc)

	return client, nil
}

func (self *Echo_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *Echo_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Echo_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *Echo_Client) GetId() string {
	return self.id
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

// Set the client instance id.
func (self *Echo_Client) SetId(id string) {
	self.id = id
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
// Stop the service.
func (self *Echo_Client) StopService() {
	self.c.Stop(globular.GetClientContext(self), &echopb.StopRequest{})
}

func (self *Echo_Client) Echo(msg interface{}) (string, error) {

	rqst := &echopb.EchoRequest{
		Message: Utility.ToString(msg),
	}
	ctx := globular.GetClientContext(self)
	rsp, err := self.c.Echo(ctx, rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}
