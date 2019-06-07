package echo_client

import (
	"context"
	"log"

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

	// The ipv4 address
	addresse string

	// The client domain
	domain string

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
func NewEcho_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *Echo_Client {
	client := new(Echo_Client)

	client.addresse = addresse
	client.domain = domain
	client.name = "persistence"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
	client.c = echopb.NewEchoServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *Echo_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *Echo_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Echo_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Echo_Client) Close() {
	self.cc.Close()
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

func (self *Echo_Client) Echo(msg interface{}) (string, error) {
	log.Println("echo service call: ", msg)
	rqst := &echopb.EchoRequest{
		Message: Utility.ToString(msg),
	}

	rsp, err := self.c.Echo(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Message, nil
}
