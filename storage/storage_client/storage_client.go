package storage_client

import (
	// "context"
	// "log"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/storage/storagepb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// storage Client Service
////////////////////////////////////////////////////////////////////////////////

type Storage_Client struct {
	cc *grpc.ClientConn
	c  storagepb.StorageServiceClient

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
func NewStorage_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *Storage_Client {

	client := new(Storage_Client)

	client.name = "storage"
	client.addresse = addresse
	client.domain = domain
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile

	client.cc = api.GetClientConnection(client)
	client.c = storagepb.NewStorageServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *Storage_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *Storage_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Storage_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Storage_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Storage_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Storage_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Storage_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Storage_Client) GetCaFile() string {
	return self.caFile
}
