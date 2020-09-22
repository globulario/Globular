package spc_client

import (
	"context"
	// "log"
	"strconv"

	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Globular/services/golang/spc/spcpb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// SPC Client Service
////////////////////////////////////////////////////////////////////////////////
type SPC_Client struct {
	cc *grpc.ClientConn
	c  spcpb.SpcServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The domain
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
func NewSpc_Client(address string, id string) (*SPC_Client, error) {
	client := new(SPC_Client)
	err := globular.InitClient(client, address, id)
	if err != nil {
		return nil, err
	}
	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = spcpb.NewSpcServiceClient(client.cc)

	return client, nil
}

func (self *SPC_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *SPC_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *SPC_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *SPC_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *SPC_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *SPC_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *SPC_Client) SetPort(port int) {
	self.port = port
}

// Set the service instance id
func (self *SPC_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *SPC_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *SPC_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *SPC_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *SPC_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *SPC_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *SPC_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *SPC_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *SPC_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *SPC_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *SPC_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}
