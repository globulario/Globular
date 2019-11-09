package storage_client

import (
	"context"
	// "log"
	"strconv"

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
func NewStorage_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string, token string) *Storage_Client {

	client := new(Storage_Client)

	client.name = "storage"
	client.domain = domain
	client.port = port
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile

	client.cc = api.GetClientConnection(client, token)
	client.c = storagepb.NewStorageServiceClient(client.cc)

	return client
}

// Return the domain
func (self *Storage_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Storage_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
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

////////////////// Service functionnality //////////////////////
func (self *Storage_Client) CreateConnection(id string, name string, connectionType float64) error {

	rqst := &storagepb.CreateConnectionRqst{
		Connection: &storagepb.Connection{
			Id:   id,
			Name: name,
			Type: storagepb.StoreType(connectionType), // Disk store (persistent)
		},
	}

	_, err := self.c.CreateConnection(context.Background(), rqst)

	return err
}

func (self *Storage_Client) OpenConnection(id string, options string) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.OpenRqst{
		Id:      id,
		Options: options,
	}

	_, err := self.c.Open(context.Background(), rqst)

	return err
}

func (self *Storage_Client) SetItem(connectionId string, key string, data []byte) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.SetItemRequest{
		Id:    connectionId,
		Key:   key,
		Value: data,
	}

	_, err := self.c.SetItem(context.Background(), rqst)
	return err
}

func (self *Storage_Client) GetItem(connectionId string, key string) ([]byte, error) {
	// I will execute a simple ldap search here...
	rqst := &storagepb.GetItemRequest{
		Id:  connectionId,
		Key: key,
	}

	rsp, err := self.c.GetItem(context.Background(), rqst)
	if err != nil {
		return nil, err
	}
	return rsp.Result, nil
}

func (self *Storage_Client) RemoveItem(connectionId string, key string) error {
	// I will execute a simple ldap search here...
	rqst := &storagepb.RemoveItemRequest{
		Id:  connectionId,
		Key: key,
	}

	_, err := self.c.RemoveItem(context.Background(), rqst)
	return err
}

func (self *Storage_Client) Clear(connectionId string) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.ClearRequest{
		Id: connectionId,
	}

	_, err := self.c.Clear(context.Background(), rqst)
	return err
}

func (self *Storage_Client) Drop(connectionId string) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.DropRequest{
		Id: connectionId,
	}

	_, err := self.c.Drop(context.Background(), rqst)
	return err
}

func (self *Storage_Client) CloseConnection(connectionId string) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.CloseRqst{
		Id: connectionId,
	}

	_, err := self.c.Close(context.Background(), rqst)
	return err
}

func (self *Storage_Client) DeleteConnection(connectionId string) error {

	// I will execute a simple ldap search here...
	rqst := &storagepb.DeleteConnectionRqst{
		Id: connectionId,
	}

	_, err := self.c.DeleteConnection(context.Background(), rqst)
	return err
}
