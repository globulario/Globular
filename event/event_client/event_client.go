package echo_client

import (
	"context"
	"log"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/event/eventpb"

	//	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Event_Client struct {
	cc *grpc.ClientConn
	c  eventpb.EventServiceClient

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
func NewEvent_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *Event_Client {
	client := new(Event_Client)

	client.addresse = addresse
	client.domain = domain
	client.name = "persistence"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
	client.c = eventpb.NewEventServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *Event_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *Event_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Event_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Event_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Event_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Event_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Event_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Event_Client) GetCaFile() string {
	return self.caFile
}

// Publish and event over the network
func (self *Event_Client) Publish(name string, data interface{}) error {
	log.Println("publish event ", name)
	rqst := &eventpb.PublishRequest{
		Evt: &eventpb.Event{
			Name: name,
			Data: data.([]byte),
		},
	}

	_, err := self.c.Publish(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Subscribe to an event. Stay in the function until
// UnSubscribe is call.
func (self *Event_Client) Subscribe(name string) error {
	return nil
}

// Exit event channel.
func (self *Event_Client) UnSubscribe(name string) error {
	return nil
}
