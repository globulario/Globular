package event_client

import (
	"context"
	"log"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/event/eventpb"

	//	"github.com/davecourtois/Utility"
	"io"

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
func NewEvent_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string, token string) *Event_Client {
	client := new(Event_Client)

	client.addresse = addresse
	client.domain = domain
	client.name = "event"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client, token)
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

// Subscribe to an event it return it subscriber uuid. The uuid must be use
// to unsubscribe from the channel. data_channel is use to get event data.
func (self *Event_Client) Subscribe(name string, data_channel chan []byte) (string, error) {
	rqst := &eventpb.SubscribeRequest{
		Name: name,
	}

	stream, err := self.c.Subscribe(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	uuid_channel := make(chan string)

	// Run in it own goroutine.
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// end of stream...
				break
			}

			if err != nil {
				break
			}

			// Get the result...
			switch v := msg.Result.(type) {
			case *eventpb.SubscribeResponse_Uuid:
				uuid_channel <- v.Uuid
			case *eventpb.SubscribeResponse_Evt:
				data_channel <- v.Evt.Data
			}
		}
	}()

	uuid := <-uuid_channel
	// Wait for subscriber uuid and return it to the function caller.
	return uuid, nil
}

// Exit event channel.
func (self *Event_Client) UnSubscribe(name string, uuid string) error {
	log.Println("UnSubscribe event ", name)
	rqst := &eventpb.UnSubscribeRequest{
		Name: name,
		Uuid: uuid,
	}

	_, err := self.c.UnSubscribe(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}
