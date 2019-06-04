package persistence_client

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/persistence/persistencepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// Persitence Client Service
////////////////////////////////////////////////////////////////////////////////
type Persistence_Client struct {
	cc *grpc.ClientConn
	c  persistencepb.PersistenceServiceClient

	// The name of the service
	name string

	// The ipv4 address
	addresse string

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
func NewPersistence_Client(addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *Persistence_Client {
	client := new(Persistence_Client)
	client.addresse = addresse
	client.name = "persistence"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile

	client.cc = api.GetClientConnection(client)
	client.c = persistencepb.NewPersistenceServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *Persistence_Client) GetAddress() string {
	return self.addresse
}

// Return the name of the service
func (self *Persistence_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Persistence_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Persistence_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Persistence_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Persistence_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Persistence_Client) GetCaFile() string {
	return self.caFile
}

// Test if a connection is found
func (self *Persistence_Client) Ping(connectionId interface{}) (string, error) {

	// Here I will try to ping a non-existing connection.
	rqst := &persistencepb.PingConnectionRqst{
		Id: Utility.ToString(connectionId),
	}

	rsp, err := self.c.Ping(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return rsp.Result, err
}

func (self *Persistence_Client) Find(connectionId string, database string, collection string, query string, fields string, options string) (string, error) {

	// Retreive a single value...
	rqst := &persistencepb.FindRqst{
		Id:         connectionId,
		Database:   database,
		Collection: collection,
		Query:      query,
		Fields:     strings.Split(fields, ","),
		Options:    options,
	}

	stream, err := self.c.Find(context.Background(), rqst)

	if err != nil {
		return "", err
	}

	values := make([]interface{}, 0)
	for {
		results, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}
		if err != nil {
			return "", nil
		}

		values_ := make([]interface{}, 0) // sub array...
		err = json.Unmarshal([]byte(results.JsonStr), &values_)
		if err != nil {
			return "", nil
		}
		values = append(values, values_...)
	}

	valuesStr, err := json.Marshal(values)
	if err != nil {
		return "", nil
	}
	return string(valuesStr), nil
}
