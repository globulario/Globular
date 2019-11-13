package opc_ua_client

import (
	/*"context"*/
	/* "log" */
	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/echo/echopb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Opc_ua_Client struct {
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
func NewOpc_ua_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string, token string) *Opc_ua_Client {
	client := new(Opc_ua_Client)
	client.domain = domain
	client.name = "echo"
	client.port = port
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
	client.c = echopb.NewEchoServiceClient(client.cc)

	return client
}

// Return the domain
func (self *Opc_ua_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Opc_ua_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *Opc_ua_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Opc_ua_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Opc_ua_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Opc_ua_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Opc_ua_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Opc_ua_Client) GetCaFile() string {
	return self.caFile
}

////////////////////////////////	 Api 	////////////////////////////////
